// Package apply serves the public sign-up forms (覺行小組報名、共好企業登記),
// the public 覺行小組 list, and the admin views over them.
package apply

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
	"github.com/geodown-ops/bodhi-alliance/server/internal/httpx"
)

type Group struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" binding:"required,max=100"`
	Region      string  `json:"region" binding:"required,max=50"`
	CenterID    *string `json:"center_id" binding:"omitempty,uuid"`
	CenterName  string  `json:"center_name" binding:"max=100"`
	Schedule    string  `json:"schedule" binding:"max=200"`
	Description string  `json:"description" binding:"max=2000"`
	IsOnline    bool    `json:"is_online"`
	IsListed    bool    `json:"is_listed"`
	SortOrder   int     `json:"sort_order"`
}

type GroupApplication struct {
	ID         string    `json:"id"`
	GroupID    *string   `json:"group_id"`
	GroupName  *string   `json:"group_name"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Region     string    `json:"region"`
	WantsCoach bool      `json:"wants_coach"`
	Message    string    `json:"message"`
	Status     string    `json:"status"`
	AdminNote  string    `json:"admin_note"`
	CreatedAt  time.Time `json:"created_at"`
}

type MerchantApplication struct {
	ID           string    `json:"id"`
	Kind         string    `json:"kind"`
	OrgName      string    `json:"org_name"`
	ContactName  string    `json:"contact_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Region       string    `json:"region"`
	Offerings    string    `json:"offerings"`
	MonthlyScale string    `json:"monthly_scale"`
	Message      string    `json:"message"`
	Status       string    `json:"status"`
	AdminNote    string    `json:"admin_note"`
	CreatedAt    time.Time `json:"created_at"`
}

type Handler struct {
	DB   *pgxpool.Pool
	Auth *auth.Service
}

func (h *Handler) Routes(r *gin.RouterGroup) {
	limit := httpx.RateLimit(5, time.Minute)
	r.GET("/groups", h.listGroups(true))
	r.POST("/group-applications", limit, h.createGroupApplication)
	r.POST("/merchant-applications", limit, h.createMerchantApplication)

	admin := r.Group("/admin", h.Auth.Require(auth.RoleAllianceAdmin, auth.ScopeAlliance))
	admin.GET("/groups", h.listGroups(false))
	admin.POST("/groups", h.saveGroup)
	admin.PUT("/groups/:id", h.saveGroup)
	admin.DELETE("/groups/:id", h.deleteGroup)
	admin.GET("/group-applications", h.listGroupApplications)
	admin.PATCH("/group-applications/:id", h.reviewApplication("group_application"))
	admin.GET("/merchant-applications", h.listMerchantApplications)
	admin.PATCH("/merchant-applications/:id", h.reviewApplication("merchant_application"))
}

func (h *Handler) listGroups(publicOnly bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := h.DB.Query(c, `
			SELECT g.id, g.name, g.region, g.center_id, coalesce(c.name, g.center_name), g.schedule, g.description,
			       g.is_online, g.is_listed, g.sort_order
			FROM practice_group g LEFT JOIN center c ON c.id = g.center_id
			WHERE g.is_listed OR NOT $1
			ORDER BY g.sort_order, g.region, g.name`, publicOnly)
		if err != nil {
			httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
			return
		}
		groups, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Group])
		if err != nil {
			httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
			return
		}
		c.JSON(http.StatusOK, groups)
	}
}

func (h *Handler) saveGroup(c *gin.Context) {
	var g Group
	if err := c.ShouldBindJSON(&g); err != nil {
		httpx.Error(c, http.StatusBadRequest, "請填寫小組名稱與地區")
		return
	}
	id := c.Param("id")
	var err error
	if id == "" {
		err = h.DB.QueryRow(c, `
			INSERT INTO practice_group (name, region, center_id, center_name, schedule, description, is_online, is_listed, sort_order)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
			g.Name, g.Region, g.CenterID, g.CenterName, g.Schedule, g.Description, g.IsOnline, g.IsListed, g.SortOrder).Scan(&g.ID)
	} else {
		var tag interface{ RowsAffected() int64 }
		tag, err = h.DB.Exec(c, `
			UPDATE practice_group SET name = $2, region = $3, center_id = $4, center_name = $5, schedule = $6,
			       description = $7, is_online = $8, is_listed = $9, sort_order = $10, updated_at = now()
			WHERE id = $1`,
			id, g.Name, g.Region, g.CenterID, g.CenterName, g.Schedule, g.Description, g.IsOnline, g.IsListed, g.SortOrder)
		if err == nil && tag.RowsAffected() == 0 {
			httpx.Error(c, http.StatusNotFound, "找不到這個小組")
			return
		}
		g.ID = id
	}
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "儲存失敗")
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *Handler) deleteGroup(c *gin.Context) {
	if _, err := h.DB.Exec(c, `DELETE FROM practice_group WHERE id = $1`, c.Param("id")); err != nil {
		httpx.Error(c, http.StatusBadRequest, "刪除失敗")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) createGroupApplication(c *gin.Context) {
	var req struct {
		GroupID    string `json:"group_id"`
		Name       string `json:"name" binding:"required,max=100"`
		Email      string `json:"email" binding:"required,email,max=200"`
		Phone      string `json:"phone" binding:"max=50"`
		Region     string `json:"region" binding:"max=50"`
		WantsCoach bool   `json:"wants_coach"`
		Message    string `json:"message" binding:"max=2000"`
		Website    string `json:"website"` // honeypot：真人看不到這個欄位
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "請填寫姓名與正確的電子郵件")
		return
	}
	if req.Website != "" {
		c.JSON(http.StatusCreated, gin.H{"ok": true})
		return
	}
	var groupID *string
	if req.GroupID != "" {
		groupID = &req.GroupID
	}
	_, err := h.DB.Exec(c, `
		INSERT INTO group_application (group_id, name, email, phone, region, wants_coach, message)
		VALUES ((SELECT id FROM practice_group WHERE id::text = $1 AND is_listed), $2, $3, $4, $5, $6, $7)`,
		groupID, strings.TrimSpace(req.Name), strings.TrimSpace(req.Email), req.Phone, req.Region, req.WantsCoach, req.Message)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "送出失敗，請稍後再試")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ok": true})
}

func (h *Handler) createMerchantApplication(c *gin.Context) {
	var req struct {
		Kind         string `json:"kind" binding:"required,oneof=center sponsor other"`
		OrgName      string `json:"org_name" binding:"required,max=200"`
		ContactName  string `json:"contact_name" binding:"required,max=100"`
		Email        string `json:"email" binding:"required,email,max=200"`
		Phone        string `json:"phone" binding:"max=50"`
		Region       string `json:"region" binding:"max=50"`
		Offerings    string `json:"offerings" binding:"max=2000"`
		MonthlyScale string `json:"monthly_scale" binding:"max=200"`
		Message      string `json:"message" binding:"max=2000"`
		Website      string `json:"website"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "請填寫身份、單位名稱、聯絡人與正確的電子郵件")
		return
	}
	if req.Website != "" {
		c.JSON(http.StatusCreated, gin.H{"ok": true})
		return
	}
	_, err := h.DB.Exec(c, `
		INSERT INTO merchant_application (kind, org_name, contact_name, email, phone, region, offerings, monthly_scale, message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		req.Kind, strings.TrimSpace(req.OrgName), strings.TrimSpace(req.ContactName), strings.TrimSpace(req.Email),
		req.Phone, req.Region, req.Offerings, req.MonthlyScale, req.Message)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "送出失敗，請稍後再試")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ok": true})
}

func (h *Handler) listGroupApplications(c *gin.Context) {
	rows, err := h.DB.Query(c, `
		SELECT a.id, a.group_id, g.name, a.name, a.email, a.phone, a.region, a.wants_coach, a.message,
		       a.status, a.admin_note, a.created_at
		FROM group_application a LEFT JOIN practice_group g ON g.id = a.group_id
		WHERE $1 = '' OR a.status = $1
		ORDER BY a.created_at DESC LIMIT 500`, c.Query("status"))
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[GroupApplication])
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) listMerchantApplications(c *gin.Context) {
	rows, err := h.DB.Query(c, `
		SELECT id, kind, org_name, contact_name, email, phone, region, offerings, monthly_scale, message,
		       status, admin_note, created_at
		FROM merchant_application
		WHERE $1 = '' OR status = $1
		ORDER BY created_at DESC LIMIT 500`, c.Query("status"))
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[MerchantApplication])
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	c.JSON(http.StatusOK, list)
}

// reviewApplication updates status and note. table is one of two fixed names, never user input.
func (h *Handler) reviewApplication(table string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Status    string `json:"status" binding:"required,oneof=new contacted accepted declined"`
			AdminNote string `json:"admin_note" binding:"max=2000"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			httpx.Error(c, http.StatusBadRequest, "狀態不正確")
			return
		}
		var id string
		err := h.DB.QueryRow(c, `UPDATE `+table+` SET status = $2, admin_note = $3, updated_at = now()
			WHERE id = $1 RETURNING id`, c.Param("id"), req.Status, req.AdminNote).Scan(&id)
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.Error(c, http.StatusNotFound, "找不到這筆資料")
			return
		}
		if err != nil {
			httpx.Error(c, http.StatusBadRequest, "更新失敗")
			return
		}
		c.Status(http.StatusNoContent)
	}
}
