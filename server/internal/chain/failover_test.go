package chain

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type fakeNode struct {
	Backend // methods the test does not stub panic
	err     error
	calls   int
}

func (n *fakeNode) ChainID(context.Context) (*big.Int, error) {
	n.calls++
	if n.err != nil {
		return nil, n.err
	}
	return big.NewInt(137), nil
}

func (n *fakeNode) TransactionReceipt(context.Context, common.Hash) (*types.Receipt, error) {
	n.calls++
	return nil, n.err
}

func (n *fakeNode) SendTransaction(context.Context, *types.Transaction) error {
	n.calls++
	return n.err
}

type codeErr struct {
	code int
	msg  string
}

func (e codeErr) Error() string  { return e.msg }
func (e codeErr) ErrorCode() int { return e.code }

func TestFailover(t *testing.T) {
	ctx := context.Background()
	down := &fakeNode{err: rpc.HTTPError{StatusCode: 503, Status: "503 Service Unavailable"}}
	busy := &fakeNode{err: codeErr{-32005, "limit exceeded"}}
	up := &fakeNode{}
	f := NewFailover([]string{"down", "busy", "up"}, []Backend{down, busy, up})

	if id, err := f.ChainID(ctx); err != nil || id.Int64() != 137 {
		t.Fatalf("ChainID = %v, %v", id, err)
	}
	if f.Current() != "up" {
		t.Errorf("current node = %s, want up", f.Current())
	}
	// The working node is used first from now on.
	f.ChainID(ctx)
	if down.calls != 1 || busy.calls != 1 || up.calls != 2 {
		t.Errorf("calls = %d %d %d", down.calls, busy.calls, up.calls)
	}

	// Answers about the request itself are not the node's fault: no switching.
	up.err = ethereum.NotFound
	if _, err := f.TransactionReceipt(ctx, common.Hash{}); !errors.Is(err, ethereum.NotFound) || down.calls != 1 {
		t.Errorf("not found: %v, down called %d", err, down.calls)
	}
	up.err = codeErr{-32000, "insufficient funds for gas * price + value"}
	if err := f.SendTransaction(ctx, nil); err == nil || down.calls != 1 {
		t.Errorf("insufficient funds: %v, down called %d", err, down.calls)
	}
	up.err = codeErr{-32000, "already known"}
	if err := f.SendTransaction(ctx, nil); err != nil {
		t.Errorf("already known should count as sent: %v", err)
	}

	// All nodes failing returns the last error.
	up.err = rpc.HTTPError{StatusCode: 500, Status: "500 Internal Server Error"}
	if _, err := f.ChainID(ctx); err == nil {
		t.Error("all nodes down but no error")
	}
}
