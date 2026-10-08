package chain

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const workerLockID = 724_201_009

// Service owns the operator key and talks to the chain.
type Service struct {
	DB       *pgxpool.Pool
	Backend  Backend
	Explorer string

	ChainID  *big.Int
	operator *ecdsa.PrivateKey
	seed     []byte
	contract common.Address // zero until deployed or configured

	// ReceiptWait bounds how long one transaction is waited on before moving on.
	ReceiptWait time.Duration
	// Poll is how often the worker looks for new members and open transactions.
	Poll time.Duration
}

// New checks the keys and asks the node which chain it is on.
func New(ctx context.Context, db *pgxpool.Pool, b Backend, cfg Config) (*Service, error) {
	op, err := parseKey(cfg.OperatorKey)
	if err != nil {
		return nil, fmt.Errorf("BODHI_CHAIN_OPERATOR_KEY: %w", err)
	}
	seed, err := parseSeed(cfg.MemberSeed)
	if err != nil {
		return nil, err
	}
	id, err := b.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}
	s := &Service{DB: db, Backend: b, Explorer: cfg.ExplorerURL, ChainID: id, operator: op, seed: seed,
		ReceiptWait: 2 * time.Minute, Poll: 20 * time.Second}
	if cfg.Contract != "" {
		if !common.IsHexAddress(cfg.Contract) {
			return nil, errors.New("BODHI_CHAIN_CONTRACT is not an address")
		}
		s.contract = common.HexToAddress(cfg.Contract)
	}
	return s, nil
}

func (s *Service) Operator() common.Address { return crypto.PubkeyToAddress(s.operator.PublicKey) }
func (s *Service) Contract() common.Address { return s.contract }

func (s *Service) settingKey(name string) string { return name + ":" + s.ChainID.String() }

// Run keeps working until ctx ends; errors are logged and retried on the next tick.
func (s *Service) Run(ctx context.Context) {
	log.Printf("chain: chain %s, operator %s", s.ChainID, s.Operator().Hex())
	for {
		if err := s.Step(ctx); err != nil && ctx.Err() == nil {
			log.Printf("chain: %v", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(s.Poll):
		}
	}
}

// Step does one round: make sure the token exists, give new members an address and
// a pending gift, then send pending gifts one at a time.
func (s *Service) Step(ctx context.Context) error {
	// Only one api instance sends at a time, so the operator's nonces never collide.
	conn, err := s.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	var locked bool
	if err := conn.QueryRow(ctx, `SELECT pg_try_advisory_lock($1)`, workerLockID).Scan(&locked); err != nil || !locked {
		return err
	}
	defer conn.Exec(context.Background(), `SELECT pg_advisory_unlock($1)`, workerLockID)

	if err := s.ensureContract(ctx); err != nil {
		return fmt.Errorf("contract: %w", err)
	}
	if err := s.enrolMembers(ctx); err != nil {
		return fmt.Errorf("enrol: %w", err)
	}
	if err := s.recheckSent(ctx); err != nil {
		return fmt.Errorf("recheck: %w", err)
	}
	return s.sendPending(ctx)
}

// ensureContract loads the token address, or deploys BodhiCoin with the operator as treasury.
func (s *Service) ensureContract(ctx context.Context) error {
	if s.contract != (common.Address{}) {
		return nil
	}
	var addr string
	err := s.DB.QueryRow(ctx, `SELECT value FROM chain_setting WHERE key = $1`, s.settingKey("contract")).Scan(&addr)
	if err == nil {
		s.contract = common.HexToAddress(addr)
		return nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	// A deploy may already be in flight from an earlier run.
	var pending string
	sentAt := time.Now()
	err = s.DB.QueryRow(ctx, `SELECT value, updated_at FROM chain_setting WHERE key = $1`, s.settingKey("deploy_tx")).Scan(&pending, &sentAt)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	var hash common.Hash
	if pending != "" {
		hash = common.HexToHash(pending)
	} else {
		ctorArgs, err := tokenABI.Pack("", s.Operator())
		if err != nil {
			return err
		}
		code := append(common.FromHex(strings.TrimSpace(binHex)), ctorArgs...)
		tx, err := s.send(ctx, nil, code)
		if err != nil {
			return fmt.Errorf("deploy: %w", err)
		}
		hash = tx.Hash()
		if _, err := s.DB.Exec(ctx, `INSERT INTO chain_setting (key, value) VALUES ($1, $2)
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = now()`, s.settingKey("deploy_tx"), hash.Hex()); err != nil {
			return err
		}
		log.Printf("chain: deploying BodhiCoin in %s", hash.Hex())
	}

	rcpt, err := s.waitReceipt(ctx, hash, sentAt)
	if err != nil {
		if errors.Is(err, errDropped) {
			_, err = s.DB.Exec(ctx, `DELETE FROM chain_setting WHERE key = $1`, s.settingKey("deploy_tx"))
		}
		return err
	}
	if rcpt.Status != types.ReceiptStatusSuccessful {
		s.DB.Exec(ctx, `DELETE FROM chain_setting WHERE key = $1`, s.settingKey("deploy_tx"))
		return fmt.Errorf("deploy %s reverted", hash.Hex())
	}
	if _, err := s.DB.Exec(ctx, `INSERT INTO chain_setting (key, value) VALUES ($1, $2)`, s.settingKey("contract"), rcpt.ContractAddress.Hex()); err != nil {
		return err
	}
	s.contract = rcpt.ContractAddress
	log.Printf("chain: BodhiCoin at %s", s.contract.Hex())
	return nil
}

// enrolMembers gives every member without one an address and a pending 入會贈幣.
func (s *Service) enrolMembers(ctx context.Context) error {
	rows, err := s.DB.Query(ctx, `
		SELECT v.id FROM volunteer v
		WHERE NOT EXISTS (SELECT 1 FROM member_chain_account a WHERE a.volunteer_id = v.id)
		ORDER BY v.created_at LIMIT 500`)
	if err != nil {
		return err
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, id := range ids {
		addr, err := MemberAddress(s.seed, id)
		if err != nil {
			return err
		}
		err = pgx.BeginFunc(ctx, s.DB, func(tx pgx.Tx) error {
			if _, err := tx.Exec(ctx, `INSERT INTO member_chain_account (volunteer_id, chain_id, address) VALUES ($1, $2, $3)
				ON CONFLICT DO NOTHING`, id, s.ChainID.Int64(), addr.Hex()); err != nil {
				return err
			}
			_, err := tx.Exec(ctx, `INSERT INTO chain_grant (volunteer_id, chain_id, kind, amount, to_address) VALUES ($1, $2, 'join', $3, $4)
				ON CONFLICT DO NOTHING`, id, s.ChainID.Int64(), JoinGift, addr.Hex())
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// recheckSent settles gifts whose transaction went out on an earlier run.
func (s *Service) recheckSent(ctx context.Context) error {
	rows, err := s.DB.Query(ctx, `SELECT id::text, tx_hash, sent_at FROM chain_grant WHERE status = 'sent' AND chain_id = $1 ORDER BY sent_at`, s.ChainID.Int64())
	if err != nil {
		return err
	}
	type sent struct {
		ID, Hash string
		At       time.Time
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[sent])
	if err != nil {
		return err
	}
	for _, g := range list {
		if err := s.settle(ctx, g.ID, common.HexToHash(g.Hash), g.At); err != nil {
			return err
		}
	}
	return nil
}

// sendPending sends open gifts oldest first, waiting for each before the next so the
// operator's nonces never collide.
func (s *Service) sendPending(ctx context.Context) error {
	for i := 0; i < 50; i++ {
		var id, to string
		var amount int64
		err := s.DB.QueryRow(ctx, `SELECT id::text, to_address, amount FROM chain_grant
			WHERE status = 'pending' AND chain_id = $1 ORDER BY created_at LIMIT 1`, s.ChainID.Int64()).Scan(&id, &to, &amount)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		data, err := tokenABI.Pack("transfer", common.HexToAddress(to), big.NewInt(amount))
		if err != nil {
			return err
		}
		contract := s.contract
		tx, err := s.send(ctx, &contract, data)
		if err != nil {
			s.DB.Exec(ctx, `UPDATE chain_grant SET last_error = $2 WHERE id = $1`, id, truncate(err.Error()))
			return fmt.Errorf("send gift %s: %w", id, err)
		}
		if _, err := s.DB.Exec(ctx, `UPDATE chain_grant SET status = 'sent', tx_hash = $2, sent_at = now(), last_error = '' WHERE id = $1`,
			id, tx.Hash().Hex()); err != nil {
			return err
		}
		if err := s.settle(ctx, id, tx.Hash(), time.Now()); err != nil {
			return err
		}
	}
	return nil
}

// settle waits for a gift's receipt and records the outcome.
func (s *Service) settle(ctx context.Context, id string, hash common.Hash, sentAt time.Time) error {
	rcpt, err := s.waitReceipt(ctx, hash, sentAt)
	switch {
	case errors.Is(err, errDropped):
		_, err = s.DB.Exec(ctx, `UPDATE chain_grant SET status = 'pending', tx_hash = NULL, sent_at = NULL, last_error = 'transaction dropped; resending' WHERE id = $1`, id)
		return err
	case errors.Is(err, errNotYet):
		return nil
	case err != nil:
		return err
	}
	if rcpt.Status != types.ReceiptStatusSuccessful {
		_, err = s.DB.Exec(ctx, `UPDATE chain_grant SET status = 'failed', block_number = $2, last_error = 'transaction reverted' WHERE id = $1`,
			id, rcpt.BlockNumber.Int64())
		return err
	}
	_, err = s.DB.Exec(ctx, `UPDATE chain_grant SET status = 'confirmed', block_number = $2, confirmed_at = now() WHERE id = $1`,
		id, rcpt.BlockNumber.Int64())
	return err
}

var (
	errNotYet   = errors.New("transaction not mined yet")
	errDropped  = errors.New("transaction dropped")
	minTipFloor = map[int64]*big.Int{
		80002: big.NewInt(25_000_000_000), // Polygon Amoy
		137:   big.NewInt(25_000_000_000), // Polygon PoS
	}
)

// DropAfter is how long a transaction the node does not know about is waited on before
// it counts as dropped and is sent again. Public RPCs are load-balanced, so a fresh
// transaction can briefly be unknown to the node that answers.
const DropAfter = 10 * time.Minute

// waitReceipt polls for a receipt up to ReceiptWait.
func (s *Service) waitReceipt(ctx context.Context, hash common.Hash, sentAt time.Time) (*types.Receipt, error) {
	deadline := time.Now().Add(s.ReceiptWait)
	for {
		rcpt, err := s.Backend.TransactionReceipt(ctx, hash)
		if err == nil {
			return rcpt, nil
		}
		if !errors.Is(err, ethereum.NotFound) {
			return nil, err
		}
		if time.Since(sentAt) > DropAfter {
			if _, _, err := s.Backend.TransactionByHash(ctx, hash); errors.Is(err, ethereum.NotFound) {
				return nil, errDropped
			}
		}
		if time.Now().After(deadline) {
			return nil, errNotYet
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
}

// send signs and broadcasts an EIP-1559 transaction from the operator.
// to == nil deploys a contract.
func (s *Service) send(ctx context.Context, to *common.Address, data []byte) (*types.Transaction, error) {
	from := s.Operator()
	nonce, err := s.Backend.PendingNonceAt(ctx, from)
	if err != nil {
		return nil, err
	}
	tip, err := s.Backend.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}
	if floor := minTipFloor[s.ChainID.Int64()]; floor != nil && tip.Cmp(floor) < 0 {
		tip = new(big.Int).Set(floor)
	}
	head, err := s.Backend.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	feeCap := new(big.Int).Add(tip, new(big.Int).Mul(head.BaseFee, big.NewInt(2)))
	gas, err := s.Backend.EstimateGas(ctx, ethereum.CallMsg{From: from, To: to, Data: data, GasTipCap: tip, GasFeeCap: feeCap})
	if err != nil {
		return nil, fmt.Errorf("estimate gas: %w", err)
	}
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID: s.ChainID, Nonce: nonce, GasTipCap: tip, GasFeeCap: feeCap,
		Gas: gas * 12 / 10, To: to, Data: data,
	})
	signed, err := types.SignTx(tx, types.LatestSignerForChainID(s.ChainID), s.operator)
	if err != nil {
		return nil, err
	}
	if err := s.Backend.SendTransaction(ctx, signed); err != nil {
		return nil, err
	}
	return signed, nil
}

// BalanceOf reads an address's 菩提幣 balance (smallest units) straight from the chain.
func (s *Service) BalanceOf(ctx context.Context, addr common.Address) (*big.Int, error) {
	if s.contract == (common.Address{}) {
		return nil, errors.New("token not deployed yet")
	}
	data, err := tokenABI.Pack("balanceOf", addr)
	if err != nil {
		return nil, err
	}
	contract := s.contract
	out, err := s.Backend.CallContract(ctx, ethereum.CallMsg{To: &contract, Data: data}, nil)
	if err != nil {
		return nil, err
	}
	vals, err := tokenABI.Unpack("balanceOf", out)
	if err != nil {
		return nil, err
	}
	return vals[0].(*big.Int), nil
}

func truncate(s string) string {
	if len(s) > 500 {
		return s[:500]
	}
	return s
}
