package chain

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
	"github.com/geodown-ops/bodhi-alliance/server/internal/httpx"
)

// Handler serves a member's on-chain 菩提幣 for the personal page.
type Handler struct {
	DB   *pgxpool.Pool
	Auth *auth.Service
	Svc  *Holder // empty when the chain keys are not configured or the node is unreachable
}

func (h *Handler) Routes(r *gin.RouterGroup) {
	r.GET("/chain", h.publicInfo)
	r.GET("/me/chain", h.Auth.RequireUser(), h.myChain)
}

// publicInfo tells the 菩提幣介紹 page which chain and contract to point visitors to.
func (h *Handler) publicInfo(c *gin.Context) {
	s := h.Svc.Get()
	if s == nil {
		c.JSON(http.StatusOK, gin.H{"enabled": false})
		return
	}
	out := gin.H{"enabled": true, "chain_id": s.ChainID.Int64(), "network": NetworkName(s.ChainID.Int64()), "explorer": s.Explorer}
	if contract := s.Contract(); contract != (common.Address{}) {
		out["contract"] = contract.Hex()
		out["contract_url"] = s.Explorer + "/token/" + contract.Hex()
	}
	c.JSON(http.StatusOK, out)
}

type grantView struct {
	Kind        string     `json:"kind"`
	Amount      string     `json:"amount"`
	Status      string     `json:"status"`
	TxHash      string     `json:"tx_hash"`
	TxURL       string     `json:"tx_url"`
	ConfirmedAt *time.Time `json:"confirmed_at"`
}

func (h *Handler) myChain(c *gin.Context) {
	s := h.Svc.Get()
	if s == nil {
		c.JSON(http.StatusOK, gin.H{"enabled": false})
		return
	}
	var vid, address string
	err := h.DB.QueryRow(c, `
		SELECT v.id, coalesce(a.address, '') FROM volunteer v
		LEFT JOIN member_chain_account a ON a.volunteer_id = v.id AND a.chain_id = $2
		WHERE v.user_id = $1`, auth.CurrentUser(c).ID, s.ChainID.Int64()).Scan(&vid, &address)
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(c, http.StatusNotFound, "這個帳號還不是會員")
		return
	}
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "處理失敗，請稍後再試")
		return
	}
	out := gin.H{"enabled": true, "chain_id": s.ChainID.Int64(), "network": NetworkName(s.ChainID.Int64()), "address": address, "grants": []grantView{}}
	contract := s.Contract()
	if contract != (common.Address{}) {
		out["contract"] = contract.Hex()
		out["contract_url"] = s.Explorer + "/token/" + contract.Hex()
	}
	if address == "" {
		c.JSON(http.StatusOK, out)
		return
	}
	out["address_url"] = s.Explorer + "/address/" + address
	if contract != (common.Address{}) {
		out["token_url"] = s.Explorer + "/token/" + contract.Hex() + "?a=" + address
	}

	rows, err := h.DB.Query(c, `SELECT kind, amount, status, coalesce(tx_hash, ''), confirmed_at
		FROM chain_grant WHERE volunteer_id = $1 AND chain_id = $2 ORDER BY created_at`, vid, s.ChainID.Int64())
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "處理失敗，請稍後再試")
		return
	}
	var grants []grantView
	for rows.Next() {
		var g grantView
		var amount int64
		if err := rows.Scan(&g.Kind, &amount, &g.Status, &g.TxHash, &g.ConfirmedAt); err != nil {
			rows.Close()
			httpx.Error(c, http.StatusInternalServerError, "處理失敗，請稍後再試")
			return
		}
		g.Amount = FormatAmount(bigInt(amount))
		if g.TxHash != "" {
			g.TxURL = s.Explorer + "/tx/" + g.TxHash
		}
		grants = append(grants, g)
	}
	rows.Close()
	if grants != nil {
		out["grants"] = grants
	}

	// 餘額直接問鏈；鏈上節點慢或連不到時，頁面仍顯示地址與交易
	ctx, cancel := context.WithTimeout(c, 4*time.Second)
	defer cancel()
	if bal, err := s.BalanceOf(ctx, common.HexToAddress(address)); err == nil {
		out["balance"] = FormatAmount(bal)
	}
	c.JSON(http.StatusOK, out)
}

// NetworkName is how the personal page names the chain.
func NetworkName(id int64) string {
	switch id {
	case 80002:
		return "Polygon Amoy 測試鏈"
	case 137:
		return "Polygon"
	default:
		return "測試鏈"
	}
}
