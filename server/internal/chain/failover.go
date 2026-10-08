package chain

import (
	"context"
	"errors"
	"log"
	"math/big"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// Failover spreads the worker's calls over several JSON-RPC nodes. Free public Polygon
// nodes often answer 5xx or time out, so when a node fails a call, the call is tried
// on the next node and that node is used from then on.
type Failover struct {
	URLs    []string
	Nodes   []Backend
	Timeout time.Duration // per attempt; 0 means 20s
	cur     atomic.Int32
}

func NewFailover(urls []string, nodes []Backend) *Failover {
	return &Failover{URLs: urls, Nodes: nodes}
}

// Current is the URL of the node calls go to first.
func (f *Failover) Current() string { return f.URLs[f.cur.Load()] }

func try[T any](ctx context.Context, f *Failover, call func(context.Context, Backend) (T, error)) (T, error) {
	timeout := f.Timeout
	if timeout == 0 {
		timeout = 20 * time.Second
	}
	start := int(f.cur.Load())
	var out T
	var err error
	for i := range f.Nodes {
		k := (start + i) % len(f.Nodes)
		actx, cancel := context.WithTimeout(ctx, timeout)
		out, err = call(actx, f.Nodes[k])
		cancel()
		if !nodeFault(err) || ctx.Err() != nil {
			if i > 0 {
				f.cur.Store(int32(k))
				log.Printf("chain: switched to %s", f.URLs[k])
			}
			return out, err
		}
	}
	return out, err
}

// nodeFault says whether an error is the node's problem (unreachable, HTTP error,
// overloaded) rather than an answer about the request, such as "not found",
// "insufficient funds" or a revert, which another node would give too.
func nodeFault(err error) bool {
	if err == nil || errors.Is(err, ethereum.NotFound) {
		return false
	}
	var rpcErr rpc.Error
	if errors.As(err, &rpcErr) {
		switch rpcErr.ErrorCode() {
		case -32000, 3: // geth's "invalid transaction" family, and execution reverted
			return false
		}
	}
	return true
}

func (f *Failover) ChainID(ctx context.Context) (*big.Int, error) {
	return try(ctx, f, func(ctx context.Context, b Backend) (*big.Int, error) { return b.ChainID(ctx) })
}

func (f *Failover) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return try(ctx, f, func(ctx context.Context, b Backend) (uint64, error) { return b.PendingNonceAt(ctx, account) })
}

func (f *Failover) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return try(ctx, f, func(ctx context.Context, b Backend) (*big.Int, error) { return b.SuggestGasTipCap(ctx) })
}

func (f *Failover) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return try(ctx, f, func(ctx context.Context, b Backend) (*types.Header, error) { return b.HeaderByNumber(ctx, number) })
}

func (f *Failover) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return try(ctx, f, func(ctx context.Context, b Backend) (uint64, error) { return b.EstimateGas(ctx, msg) })
}

// SendTransaction treats "already known" as sent: a node that failed may still have
// passed the transaction on before the retry reached the next one.
func (f *Failover) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	_, err := try(ctx, f, func(ctx context.Context, b Backend) (struct{}, error) { return struct{}{}, b.SendTransaction(ctx, tx) })
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "already known") {
		return nil
	}
	return err
}

func (f *Failover) TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	return try(ctx, f, func(ctx context.Context, b Backend) (*types.Receipt, error) { return b.TransactionReceipt(ctx, hash) })
}

func (f *Failover) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	type res struct {
		tx      *types.Transaction
		pending bool
	}
	r, err := try(ctx, f, func(ctx context.Context, b Backend) (res, error) {
		tx, pending, err := b.TransactionByHash(ctx, hash)
		return res{tx, pending}, err
	})
	return r.tx, r.pending, err
}

func (f *Failover) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return try(ctx, f, func(ctx context.Context, b Backend) ([]byte, error) { return b.CallContract(ctx, msg, blockNumber) })
}

func (f *Failover) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return try(ctx, f, func(ctx context.Context, b Backend) ([]byte, error) { return b.CodeAt(ctx, account, blockNumber) })
}
