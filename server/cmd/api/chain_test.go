package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/gin-gonic/gin"

	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
	"github.com/geodown-ops/bodhi-alliance/server/internal/chain"
	"github.com/geodown-ops/bodhi-alliance/server/internal/config"
	"github.com/geodown-ops/bodhi-alliance/server/internal/testdb"
)

// Every member gets an address and 1 枚 菩提幣 on a simulated chain; members cannot
// pass coins to each other.
func TestJoinGiftOnChain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.New(t)
	ctx := context.Background()
	authSvc := &auth.Service{DB: pool}

	opKey, _ := crypto.GenerateKey()
	op := crypto.PubkeyToAddress(opKey.PublicKey)
	sim := simulated.NewBackend(types.GenesisAlloc{op: {Balance: new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil)}})
	defer sim.Close()
	stop := make(chan struct{})
	defer close(stop)
	go func() { // mine a block every 50ms
		for {
			select {
			case <-stop:
				return
			case <-time.After(50 * time.Millisecond):
				sim.Commit()
			}
		}
	}()

	seed := strings.Repeat("ab", 32)
	svc, err := chain.New(ctx, pool, sim.Client(), chain.Config{
		OperatorKey: hex.EncodeToString(crypto.FromECDSA(opKey)), MemberSeed: seed, ExplorerURL: "https://explorer.test"})
	if err != nil {
		t.Fatal(err)
	}
	r := NewRouter(config.Config{}, authSvc, chain.NewHolder(svc))

	signup := func(email string) string {
		w := call(t, r, http.MethodPost, "/api/volunteers", "", map[string]any{"email": email, "password": "lotus-pond-evening", "legal_name": "陳大華"})
		if w.Code != http.StatusCreated {
			t.Fatalf("register: %d %s", w.Code, w.Body)
		}
		var reg struct{ Token string }
		json.Unmarshal(w.Body.Bytes(), &reg)
		return reg.Token
	}
	mei, hua := signup("mei@example.org"), signup("hua@example.org")

	type view struct {
		Enabled, Address, AddressURL, Balance, Contract string
		Grants                                          []struct{ Kind, Amount, Status, TxHash, TxURL string }
	}
	get := func(token string) map[string]any {
		w := call(t, r, http.MethodGet, "/api/me/chain", token, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("me/chain: %d %s", w.Code, w.Body)
		}
		var v map[string]any
		json.Unmarshal(w.Body.Bytes(), &v)
		return v
	}
	if v := get(mei); v["enabled"] != true || v["address"] != "" {
		t.Fatalf("before the worker runs: %v", v)
	}

	if err := svc.Step(ctx); err != nil {
		t.Fatal(err)
	}
	if svc.Contract() == (common.Address{}) {
		t.Fatal("token was not deployed")
	}
	if w := call(t, r, http.MethodGet, "/api/chain", "", nil); w.Code != http.StatusOK ||
		!strings.Contains(w.Body.String(), `"contract_url":"https://explorer.test/token/`+svc.Contract().Hex()+`"`) {
		t.Errorf("public chain info: %d %s", w.Code, w.Body)
	}
	// A second run neither redeploys nor gives twice.
	deployed := svc.Contract()
	if err := svc.Step(ctx); err != nil {
		t.Fatal(err)
	}
	if svc.Contract() != deployed {
		t.Error("token redeployed")
	}

	var addrs []common.Address
	for _, tok := range []string{mei, hua} {
		v := get(tok)
		grants, _ := v["grants"].([]any)
		if len(grants) != 1 {
			t.Fatalf("grants: %v", v)
		}
		g := grants[0].(map[string]any)
		if g["status"] != "confirmed" || g["amount"] != "1.00" || !strings.HasPrefix(g["tx_url"].(string), "https://explorer.test/tx/0x") {
			t.Errorf("grant: %v", g)
		}
		if v["balance"] != "1.00" || !strings.HasPrefix(v["address_url"].(string), "https://explorer.test/address/0x") {
			t.Errorf("view: %v", v)
		}
		addrs = append(addrs, common.HexToAddress(v["address"].(string)))
	}
	if addrs[0] == addrs[1] {
		t.Error("members share an address")
	}
	var meiID string
	pool.QueryRow(ctx, `SELECT v.id FROM volunteer v JOIN app_user u ON u.id = v.user_id WHERE u.email = 'mei@example.org'`).Scan(&meiID)
	seedBytes, _ := hex.DecodeString(seed)
	if want, _ := chain.MemberAddress(seedBytes, meiID); want != addrs[0] {
		t.Errorf("address %s is not derived from the seed (%s)", addrs[0], want)
	}

	treasury, _ := svc.BalanceOf(ctx, op)
	if want := new(big.Int).Sub(big.NewInt(500_000_000*100), big.NewInt(2*chain.JoinGift)); treasury.Cmp(want) != 0 {
		t.Errorf("treasury = %s, want %s", treasury, want)
	}

	// Member to member is refused by the contract.
	tokenABI, _ := abi.JSON(strings.NewReader(`[{"type":"function","name":"transfer","inputs":[{"type":"address"},{"type":"uint256"}],"outputs":[{"type":"bool"}]}]`))
	data, _ := tokenABI.Pack("transfer", addrs[1], big.NewInt(1))
	if _, err := sim.Client().EstimateGas(ctx, ethereum.CallMsg{From: addrs[0], To: &deployed, Data: data}); err == nil {
		t.Error("member-to-member transfer was allowed")
	}
	// The treasury can still pay a member.
	if _, err := sim.Client().EstimateGas(ctx, ethereum.CallMsg{From: op, To: &deployed, Data: data}); err != nil {
		t.Errorf("treasury transfer: %v", err)
	}
}
