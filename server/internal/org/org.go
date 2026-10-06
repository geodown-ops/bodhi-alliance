// Package org is the back-office side of 場域管理 (centers and their venues) and
// 共好企業管理 (merchants), including turning a public registration into a merchant.
package org

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
	"github.com/geodown-ops/bodhi-alliance/server/internal/httpx"
)

type Center struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" binding:"required,max=100"`
	Region      string    `json:"region" binding:"max=50"`
	Address     string    `json:"address" binding:"max=200"`
	ContactName string    `json:"contact_name" binding:"max=100"`
	Email       string    `json:"email" binding:"omitempty,email,max=200"`
	Phone       string    `json:"phone" binding:"max=50"`
	Status      string    `json:"status" binding:"required,oneof=active suspended"`
	Note        string    `json:"note" binding:"max=2000"`
	VenueCount  int       `json:"venue_count"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Venue struct {
	ID            string    `json:"id"`
	CenterID      string    `json:"center_id" binding:"required,uuid"`
	CenterName    string    `json:"center_name"`
	Name          string    `json:"name" binding:"required,max=100"`
	Address       string    `json:"address" binding:"max=200"`
	Description   string    `json:"description" binding:"max=2000"`
	ChargesPublic bool      `json:"charges_public"`
	MerchantID    *string   `json:"merchant_id" binding:"omitempty,uuid"`
	MerchantName  *string   `json:"merchant_name"`
	Status        string    `json:"status" binding:"required,oneof=active closed"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Merchant struct {
	ID                 string    `json:"id"`
	Kind               string    `json:"kind" binding:"required,oneof=alliance_unit sponsor"`
	CenterID           *string   `json:"center_id" binding:"omitempty,uuid"`
	CenterName         *string   `json:"center_name"`
	Name               string    `json:"name" binding:"required,max=200"`
	ContactName        string    `json:"contact_name" binding:"max=100"`
	Email              string    `json:"email" binding:"omitempty,email,max=200"`
	Phone              string    `json:"phone" binding:"max=50"`
	Region             string    `json:"region" binding:"max=50"`
	Address            string    `json:"address" binding:"max=200"`
	Offerings          string    `json:"offerings" binding:"max=2000"`
	ProposedMonthlyCap *int64    `json:"proposed_monthly_cap" binding:"omitempty,min=0"`
	Status             string    `json:"status" binding:"required,oneof=pending active suspended"`
	ApplicationID      *string   `json:"application_id"`
	Note               string    `json:"note" binding:"max=2000"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Handler struct {
	DB   *pgxpool.Pool
	Auth *auth.Service
}

func (h *Handler) Routes(r *gin.RouterGroup) {
	r.GET("/venues", h.publicVenues)

	admin := r.Group("/admin", h.Auth.Require(auth.RoleAllianceAdmin, auth.ScopeAlliance))
	admin.GET("/centers", h.listCenters)
	admin.POST("/centers", h.saveCenter)
	admin.PUT("/centers/:id", h.saveCenter)
	admin.DELETE("/centers/:id", h.deleteRow("center", "這個中心還有場域或聯盟單位，請先移除或改到別的中心"))
	admin.GET("/venues", h.listVenues)
	admin.POST("/venues", h.saveVenue)
	admin.PUT("/venues/:id", h.saveVenue)
	admin.DELETE("/venues/:id", h.deleteRow("venue", ""))
	admin.GET("/merchants", h.listMerchants)
	admin.POST("/merchants", h.saveMerchant)
	admin.PUT("/merchants/:id", h.saveMerchant)
	admin.DELETE("/merchants/:id", h.deleteRow("merchant", ""))
	admin.POST("/merchant-applications/:id/merchant", h.fromApplication)
}

func fail(c *gin.Context, err error) {
	var pg *pgconn.PgError
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		httpx.Error(c, http.StatusNotFound, "找不到這筆資料")
	case errors.As(err, &pg) && pg.Code == "23505":
		httpx.Error(c, http.StatusConflict, "名稱重複了")
	case errors.As(err, &pg) && pg.Code == "23514":
		httpx.Error(c, http.StatusBadRequest, "聯盟單位要選所屬中心；贊助商家不屬於任何中心")
	case errors.As(err, &pg) && pg.Code == "23503":
		httpx.Error(c, http.StatusBadRequest, "選的中心或共好企業不存在")
	case errors.As(err, &pg) && pg.Code == "22P02":
		httpx.Error(c, http.StatusNotFound, "找不到這筆資料")
	default:
		httpx.Error(c, http.StatusInternalServerError, "儲存失敗")
	}
}

func bad(c *gin.Context, msg string) { httpx.Error(c, http.StatusBadRequest, msg) }

func (h *Handler) listCenters(c *gin.Context) {
	rows, err := h.DB.Query(c, `
		SELECT c.id, c.name, c.region, c.address, c.contact_name, c.email, c.phone, c.status, c.note,
		       (SELECT count(*) FROM venue v WHERE v.center_id = c.id), c.updated_at
		FROM center c ORDER BY c.region, c.name`)
	if err != nil {
		fail(c, err)
		return
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Center])
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) saveCenter(c *gin.Context) {
	var v Center
	if err := c.ShouldBindJSON(&v); err != nil {
		bad(c, "請填寫中心名稱，電子郵件格式要正確")
		return
	}
	v.Name = strings.TrimSpace(v.Name)
	var err error
	if id := c.Param("id"); id == "" {
		err = h.DB.QueryRow(c, `
			INSERT INTO center (name, region, address, contact_name, email, phone, status, note)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
			v.Name, v.Region, v.Address, v.ContactName, v.Email, v.Phone, v.Status, v.Note).Scan(&v.ID)
	} else {
		err = h.DB.QueryRow(c, `
			UPDATE center SET name = $2, region = $3, address = $4, contact_name = $5, email = $6, phone = $7,
			       status = $8, note = $9, updated_at = now()
			WHERE id = $1 RETURNING id`,
			id, v.Name, v.Region, v.Address, v.ContactName, v.Email, v.Phone, v.Status, v.Note).Scan(&v.ID)
	}
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *Handler) listVenues(c *gin.Context) {
	rows, err := h.DB.Query(c, `
		SELECT v.id, v.center_id, c.name, v.name, v.address, v.description, v.charges_public,
		       v.merchant_id, m.name, v.status, v.updated_at
		FROM venue v JOIN center c ON c.id = v.center_id LEFT JOIN merchant m ON m.id = v.merchant_id
		WHERE $1 = '' OR v.center_id::text = $1
		ORDER BY c.name, v.name`, c.Query("center_id"))
	if err != nil {
		fail(c, err)
		return
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Venue])
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

// PublicVenue is an active venue as the 覺行小組 page lists it, with how many open
// activities there have not ended yet.
type PublicVenue struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CenterName  string `json:"center_name"`
	Region      string `json:"region"`
	Address     string `json:"address"`
	Description string `json:"description"`
	Upcoming    int    `json:"upcoming"`
}

func (h *Handler) publicVenues(c *gin.Context) {
	rows, err := h.DB.Query(c, `
		SELECT v.id, v.name, c.name, c.region, v.address, v.description,
		       (SELECT count(*) FROM practice_event e WHERE e.venue_id = v.id AND e.status = 'open' AND e.ends_at > now())
		FROM venue v JOIN center c ON c.id = v.center_id
		WHERE v.status = 'active' AND c.status = 'active'
		ORDER BY c.region, c.name, v.name`)
	if err != nil {
		fail(c, err)
		return
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[PublicVenue])
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) saveVenue(c *gin.Context) {
	var v Venue
	if err := c.ShouldBindJSON(&v); err != nil {
		bad(c, "請選所屬中心並填寫場域名稱")
		return
	}
	v.Name = strings.TrimSpace(v.Name)
	var err error
	if id := c.Param("id"); id == "" {
		err = h.DB.QueryRow(c, `
			INSERT INTO venue (center_id, name, address, description, charges_public, merchant_id, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			v.CenterID, v.Name, v.Address, v.Description, v.ChargesPublic, v.MerchantID, v.Status).Scan(&v.ID)
	} else {
		err = h.DB.QueryRow(c, `
			UPDATE venue SET center_id = $2, name = $3, address = $4, description = $5, charges_public = $6,
			       merchant_id = $7, status = $8, updated_at = now()
			WHERE id = $1 RETURNING id`,
			id, v.CenterID, v.Name, v.Address, v.Description, v.ChargesPublic, v.MerchantID, v.Status).Scan(&v.ID)
	}
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *Handler) listMerchants(c *gin.Context) {
	rows, err := h.DB.Query(c, `
		SELECT m.id, m.kind, m.center_id, c.name, m.name, m.contact_name, m.email, m.phone, m.region, m.address,
		       m.offerings, m.proposed_monthly_cap, m.status, m.application_id, m.note, m.updated_at
		FROM merchant m LEFT JOIN center c ON c.id = m.center_id
		WHERE $1 = '' OR m.status = $1
		ORDER BY m.status, m.name`, c.Query("status"))
	if err != nil {
		fail(c, err)
		return
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Merchant])
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) saveMerchant(c *gin.Context) {
	var m Merchant
	if err := c.ShouldBindJSON(&m); err != nil {
		bad(c, "請填寫企業名稱與類型，電子郵件格式要正確")
		return
	}
	m.Name = strings.TrimSpace(m.Name)
	if m.Kind == "sponsor" {
		m.CenterID = nil
	}
	var err error
	if id := c.Param("id"); id == "" {
		err = h.DB.QueryRow(c, `
			INSERT INTO merchant (kind, center_id, name, contact_name, email, phone, region, address, offerings,
			                      proposed_monthly_cap, status, note)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
			m.Kind, m.CenterID, m.Name, m.ContactName, m.Email, m.Phone, m.Region, m.Address, m.Offerings,
			m.ProposedMonthlyCap, m.Status, m.Note).Scan(&m.ID)
	} else {
		err = h.DB.QueryRow(c, `
			UPDATE merchant SET kind = $2, center_id = $3, name = $4, contact_name = $5, email = $6, phone = $7,
			       region = $8, address = $9, offerings = $10, proposed_monthly_cap = $11, status = $12, note = $13,
			       updated_at = now()
			WHERE id = $1 RETURNING id`,
			id, m.Kind, m.CenterID, m.Name, m.ContactName, m.Email, m.Phone, m.Region, m.Address, m.Offerings,
			m.ProposedMonthlyCap, m.Status, m.Note).Scan(&m.ID)
	}
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, m)
}

// fromApplication creates a pending merchant from a 共好企業登記 and marks the
// registration accepted. A center's registration becomes an alliance unit only once
// an admin picks the center, so it is created as a sponsor unless center_id is given.
func (h *Handler) fromApplication(c *gin.Context) {
	var req struct {
		CenterID string `json:"center_id"`
	}
	_ = c.ShouldBindJSON(&req)
	var id string
	err := pgx.BeginFunc(c, h.DB, func(tx pgx.Tx) error {
		var centerID *string
		kind := "sponsor"
		if req.CenterID != "" {
			centerID, kind = &req.CenterID, "alliance_unit"
		}
		err := tx.QueryRow(c, `
			INSERT INTO merchant (kind, center_id, name, contact_name, email, phone, region, offerings, note, application_id)
			SELECT $2, $3, org_name, contact_name, email, phone, region, offerings,
			       CASE WHEN monthly_scale <> '' THEN '登記時填的每月規模：' || monthly_scale ELSE '' END, id
			FROM merchant_application WHERE id::text = $1
			RETURNING id`, c.Param("id"), kind, centerID).Scan(&id)
		if err != nil {
			return err
		}
		_, err = tx.Exec(c, `UPDATE merchant_application SET status = 'accepted', updated_at = now() WHERE id::text = $1`, c.Param("id"))
		return err
	})
	var pg *pgconn.PgError
	if errors.As(err, &pg) && pg.Code == "23505" {
		httpx.Error(c, http.StatusConflict, "這筆登記已經建立過共好企業")
		return
	}
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// deleteRow deletes by id from one of three fixed table names (never user input).
func (h *Handler) deleteRow(table, inUse string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tag, err := h.DB.Exec(c, `DELETE FROM `+table+` WHERE id::text = $1`, c.Param("id"))
		var pg *pgconn.PgError
		if errors.As(err, &pg) && pg.Code == "23503" && inUse != "" {
			httpx.Error(c, http.StatusConflict, inUse)
			return
		}
		if err != nil {
			fail(c, err)
			return
		}
		if tag.RowsAffected() == 0 {
			httpx.Error(c, http.StatusNotFound, "找不到這筆資料")
			return
		}
		c.Status(http.StatusNoContent)
	}
}
