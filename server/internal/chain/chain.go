// Package chain puts 菩提幣 on-chain: it deploys the BodhiCoin ERC-20 the first
// time it runs, gives every 系統會員 a platform-custodied address, and has the
// treasury send each member 1 枚 as a sign-up gift. A background worker in the api
// service does the sending; members see the address and the transaction, with
// block-explorer links, on their personal page.
//
// Member private keys are never stored: each is derived from BODHI_CHAIN_MEMBER_SEED
// and the member's id, so the platform can sign for a member later if it must.
package chain

import (
	"context"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

//go:embed BodhiCoin.abi.json
var abiJSON string

//go:embed BodhiCoin.bin
var binHex string

// JoinGift is the sign-up gift in the token's smallest unit (2 decimals): 1 枚.
const JoinGift = 100

// Backend is the slice of an Ethereum JSON-RPC client the worker needs;
// *ethclient.Client and the simulated backend's client both satisfy it.
type Backend interface {
	ChainID(ctx context.Context) (*big.Int, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
}

// Config comes from environment variables; the feature is off unless both keys are set.
type Config struct {
	Network     string // amoy (測試鏈，預設) or polygon (主網)
	ChainID     int64  // the chain the network preset expects; 0 skips the check
	RPCURL      string // one or more JSON-RPC URLs, comma separated
	OperatorKey string // hex private key of the platform operator (also the testnet treasury)
	MemberSeed  string // hex secret that member addresses are derived from
	Contract    string // optional: use an existing BodhiCoin instead of deploying one
	ExplorerURL string
}

// Network presets: BODHI_CHAIN_NETWORK picks the chain and its default nodes and explorer.
type Network struct {
	ChainID  int64
	RPCs     string
	Explorer string
}

var Networks = map[string]Network{
	"amoy":    {80002, "https://rpc-amoy.polygon.technology,https://polygon-amoy-bor-rpc.publicnode.com,https://polygon-amoy.drpc.org", "https://amoy.polygonscan.com"},
	"polygon": {137, "https://polygon-rpc.com,https://polygon-bor-rpc.publicnode.com,https://polygon.drpc.org", "https://polygonscan.com"},
}

func ConfigFromEnv() Config {
	name := strings.ToLower(strings.TrimSpace(os.Getenv("BODHI_CHAIN_NETWORK")))
	if name == "" {
		name = "amoy"
	}
	net, ok := Networks[name]
	if !ok {
		net = Network{RPCs: "", Explorer: ""}
	}
	get := func(k, def string) string {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return v
		}
		return def
	}
	return Config{
		Network:     name,
		ChainID:     net.ChainID,
		RPCURL:      get("BODHI_CHAIN_RPC", net.RPCs),
		OperatorKey: get("BODHI_CHAIN_OPERATOR_KEY", ""),
		MemberSeed:  get("BODHI_CHAIN_MEMBER_SEED", ""),
		Contract:    get("BODHI_CHAIN_CONTRACT", ""),
		ExplorerURL: strings.TrimRight(get("BODHI_CHAIN_EXPLORER", net.Explorer), "/"),
	}
}

func (c Config) Enabled() bool { return c.OperatorKey != "" && c.MemberSeed != "" }

var tokenABI = func() abi.ABI {
	a, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		panic(err)
	}
	return a
}()

func parseKey(s string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(strings.TrimPrefix(strings.TrimSpace(s), "0x"))
}

func parseSeed(s string) ([]byte, error) {
	b, err := hex.DecodeString(strings.TrimPrefix(strings.TrimSpace(s), "0x"))
	if err != nil {
		return nil, fmt.Errorf("BODHI_CHAIN_MEMBER_SEED: %w", err)
	}
	if len(b) < 32 {
		return nil, errors.New("BODHI_CHAIN_MEMBER_SEED must be at least 32 bytes of hex")
	}
	return b, nil
}

// MemberKey derives a member's private key from the seed and their volunteer id.
func MemberKey(seed []byte, volunteerID string) (*ecdsa.PrivateKey, error) {
	for i := 0; i < 8; i++ {
		m := hmac.New(sha256.New, seed)
		m.Write([]byte("bodhi-member:" + volunteerID + ":" + strconv.Itoa(i)))
		if k, err := crypto.ToECDSA(m.Sum(nil)); err == nil {
			return k, nil
		}
	}
	return nil, errors.New("could not derive a member key")
}

// MemberAddress is the public address for MemberKey.
func MemberAddress(seed []byte, volunteerID string) (common.Address, error) {
	k, err := MemberKey(seed, volunteerID)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(k.PublicKey), nil
}

// FormatAmount renders smallest units as 菩提幣 with 2 decimals, e.g. 100 -> "1.00".
func FormatAmount(units *big.Int) string {
	if units == nil {
		return "0.00"
	}
	q, r := new(big.Int).QuoRem(units, big.NewInt(100), new(big.Int))
	return fmt.Sprintf("%s.%02d", q.String(), r.Abs(r).Int64())
}

func bigInt(v int64) *big.Int { return big.NewInt(v) }

// Holder hands the handler a Service that may only appear once the node is reachable.
type Holder struct{ p atomic.Pointer[Service] }

func NewHolder(s *Service) *Holder {
	h := &Holder{}
	h.p.Store(s)
	return h
}

// Get is nil when the holder is nil or the chain is not connected yet.
func (h *Holder) Get() *Service {
	if h == nil {
		return nil
	}
	return h.p.Load()
}

func (h *Holder) Set(s *Service) { h.p.Store(s) }

// RPCs splits BODHI_CHAIN_RPC, which may list several nodes separated by commas.
func (c Config) RPCs() []string {
	var out []string
	for _, u := range strings.Split(c.RPCURL, ",") {
		if u = strings.TrimSpace(u); u != "" {
			out = append(out, u)
		}
	}
	return out
}
