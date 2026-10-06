// Package association serves what members of 世界佛教教育協會 get: the 協會會刊, the
// 會員行事曆 and 協會通知. Members see them on their own page; alliance admins publish
// them in the back office.
package association

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

type Issue struct {
	ID       string `json:"id"`
	Title    string `json:"title" binding:"required,max=200"`
	IssuedOn string `json:"issued_on" binding:"required,datetime=2006-01-02"`
	Summary  string `json:"summary" binding:"max=4000"`
	URL      string `json:"url" binding:"omitempty,url,max=1000"`
}

type Event struct {
	ID          string     `json:"id"`
	Title       string     `json:"title" binding:"required,max=200"`
	StartsAt    time.Time  `json:"starts_at" binding:"required"`
	EndsAt      *time.Time `json:"ends_at"`
	Location    string     `json:"location" binding:"max=500"`
	Description string     `json:"description" binding:"max=4000"`
}

type Notice struct {
	ID        string    `json:"id"`
	Title     string    `json:"title" binding:"required,max=200"`
	Body      string    `json:"body" binding:"max=8000"`
	CreatedAt time.Time `json:"created_at"`
}

type Handler struct {
	DB   *pgxpool.Pool
	Auth *auth.Service
	// Now is replaceable in tests.
	Now func() time.Time
}

func (h *Handler) now() time.Time {
	if h.Now != nil {
		return h.Now()
	}
	return time.Now()
}

func (h *Handler) Routes(r *gin.RouterGroup) {
	r.GET("/me/association", h.Auth.RequireUser(), h.mine)

	admin := r.Group("/admin/association", h.Auth.Require(auth.RoleAllianceAdmin, auth.ScopeAlliance))
	admin.GET("/issues", h.listIssues)
	admin.POST("/issues", h.saveIssue)
	admin.PUT("/issues/:id", h.saveIssue)
	admin.DELETE("/issues/:id", h.deleteRow("association_issue"))
	admin.GET("/events", h.listEvents)
	admin.POST("/events", h.saveEvent)
	admin.PUT("/events/:id", h.saveEvent)
	admin.DELETE("/events/:id", h.deleteRow("association_event"))
	admin.GET("/notices", h.listNotices)
	admin.POST("/notices", h.saveNotice)
	admin.PUT("/notices/:id", h.saveNotice)
	admin.DELETE("/notices/:id", h.deleteRow("association_notice"))
	admin.GET("/members", h.members)
}

func fail(c *gin.Context, err error) {
	var pg *pgconn.PgError
	switch {
	case errors.Is(err, pgx.ErrNoRows), errors.As(err, &pg) && pg.Code == "22P02":
		httpx.Error(c, http.StatusNotFound, "找不到這筆資料")
	case errors.As(err, &pg) && pg.Code == "23514":
		httpx.Error(c, http.StatusBadRequest, "結束時間要晚於開始時間")
	default:
		httpx.Error(c, http.StatusInternalServerError, "處理失敗，請稍後再試")
	}
}

func bad(c *gin.Context, msg string) { httpx.Error(c, http.StatusBadRequest, msg) }

const (
	issueSelect  = `SELECT id, title, to_char(issued_on, 'YYYY-MM-DD'), summary, url FROM association_issue`
	eventSelect  = `SELECT id, title, starts_at, ends_at, location, description FROM association_event`
	noticeSelect = `SELECT id, title, body, created_at FROM association_notice`
)

func collect[T any](c *gin.Context, rows pgx.Rows, err error) ([]T, bool) {
	if err != nil {
		fail(c, err)
		return nil, false
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[T])
	if err != nil {
		fail(c, err)
		return nil, false
	}
	return list, true
}

// mine is what an association member sees on their page: recent issues, the calendar
// from today on, and recent notices.
func (h *Handler) mine(c *gin.Context) {
	var in bool
	err := h.DB.QueryRow(c, `SELECT in_association FROM volunteer WHERE user_id = $1`, auth.CurrentUser(c).ID).Scan(&in)
	if errors.Is(err, pgx.ErrNoRows) || err == nil && !in {
		httpx.Error(c, http.StatusForbidden, "請先在個人頁選擇加入世界佛教教育協會")
		return
	}
	if err != nil {
		fail(c, err)
		return
	}
	rows, err := h.DB.Query(c, issueSelect+` ORDER BY issued_on DESC LIMIT 24`)
	issues, ok := collect[Issue](c, rows, err)
	if !ok {
		return
	}
	rows, err = h.DB.Query(c, eventSelect+` WHERE coalesce(ends_at, starts_at) >= $1 ORDER BY starts_at LIMIT 50`, h.now().Add(-24*time.Hour))
	events, ok := collect[Event](c, rows, err)
	if !ok {
		return
	}
	rows, err = h.DB.Query(c, noticeSelect+` ORDER BY created_at DESC LIMIT 30`)
	notices, ok := collect[Notice](c, rows, err)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"issues": issues, "events": events, "notices": notices})
}

func (h *Handler) listIssues(c *gin.Context) {
	rows, err := h.DB.Query(c, issueSelect+` ORDER BY issued_on DESC`)
	if list, ok := collect[Issue](c, rows, err); ok {
		c.JSON(http.StatusOK, list)
	}
}

func (h *Handler) saveIssue(c *gin.Context) {
	var v Issue
	if err := c.ShouldBindJSON(&v); err != nil {
		bad(c, "請填寫會刊名稱與出刊日期；連結要是完整網址")
		return
	}
	v.Title = strings.TrimSpace(v.Title)
	var err error
	if id := c.Param("id"); id == "" {
		err = h.DB.QueryRow(c, `INSERT INTO association_issue (title, issued_on, summary, url) VALUES ($1, $2, $3, $4) RETURNING id`,
			v.Title, v.IssuedOn, strings.TrimSpace(v.Summary), strings.TrimSpace(v.URL)).Scan(&v.ID)
	} else {
		err = h.DB.QueryRow(c, `
			UPDATE association_issue SET title = $2, issued_on = $3, summary = $4, url = $5, updated_at = now()
			WHERE id = $1 RETURNING id`, id, v.Title, v.IssuedOn, strings.TrimSpace(v.Summary), strings.TrimSpace(v.URL)).Scan(&v.ID)
	}
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *Handler) listEvents(c *gin.Context) {
	rows, err := h.DB.Query(c, eventSelect+` ORDER BY starts_at DESC`)
	if list, ok := collect[Event](c, rows, err); ok {
		c.JSON(http.StatusOK, list)
	}
}

func (h *Handler) saveEvent(c *gin.Context) {
	var v Event
	if err := c.ShouldBindJSON(&v); err != nil {
		bad(c, "請填寫活動名稱與開始時間")
		return
	}
	v.Title = strings.TrimSpace(v.Title)
	var err error
	if id := c.Param("id"); id == "" {
		err = h.DB.QueryRow(c, `
			INSERT INTO association_event (title, starts_at, ends_at, location, description) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			v.Title, v.StartsAt, v.EndsAt, strings.TrimSpace(v.Location), strings.TrimSpace(v.Description)).Scan(&v.ID)
	} else {
		err = h.DB.QueryRow(c, `
			UPDATE association_event SET title = $2, starts_at = $3, ends_at = $4, location = $5, description = $6, updated_at = now()
			WHERE id = $1 RETURNING id`,
			id, v.Title, v.StartsAt, v.EndsAt, strings.TrimSpace(v.Location), strings.TrimSpace(v.Description)).Scan(&v.ID)
	}
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *Handler) listNotices(c *gin.Context) {
	rows, err := h.DB.Query(c, noticeSelect+` ORDER BY created_at DESC`)
	if list, ok := collect[Notice](c, rows, err); ok {
		c.JSON(http.StatusOK, list)
	}
}

func (h *Handler) saveNotice(c *gin.Context) {
	var v Notice
	if err := c.ShouldBindJSON(&v); err != nil {
		bad(c, "請填寫通知標題")
		return
	}
	v.Title = strings.TrimSpace(v.Title)
	var err error
	if id := c.Param("id"); id == "" {
		err = h.DB.QueryRow(c, `INSERT INTO association_notice (title, body) VALUES ($1, $2) RETURNING id, created_at`,
			v.Title, strings.TrimSpace(v.Body)).Scan(&v.ID, &v.CreatedAt)
	} else {
		err = h.DB.QueryRow(c, `
			UPDATE association_notice SET title = $2, body = $3, updated_at = now() WHERE id = $1 RETURNING id, created_at`,
			id, v.Title, strings.TrimSpace(v.Body)).Scan(&v.ID, &v.CreatedAt)
	}
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, v)
}

// members lists who has joined the association, newest first.
func (h *Handler) members(c *gin.Context) {
	type member struct {
		ID          string     `json:"id"`
		LegalName   string     `json:"legal_name"`
		DisplayName string     `json:"display_name"`
		Email       string     `json:"email"`
		Phone       string     `json:"phone"`
		LineID      string     `json:"line_id"`
		InGroups    bool       `json:"in_groups"`
		JoinedAt    *time.Time `json:"joined_at"`
	}
	rows, err := h.DB.Query(c, `
		SELECT v.id, v.legal_name, u.display_name, u.email, v.phone, v.line_id, v.in_groups, v.association_joined_at
		FROM volunteer v JOIN app_user u ON u.id = v.user_id
		WHERE v.in_association ORDER BY v.association_joined_at DESC NULLS LAST LIMIT 2000`)
	if list, ok := collect[member](c, rows, err); ok {
		c.JSON(http.StatusOK, list)
	}
}

// deleteRow deletes by id from one of three fixed table names (never user input).
func (h *Handler) deleteRow(table string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tag, err := h.DB.Exec(c, `DELETE FROM `+table+` WHERE id::text = $1`, c.Param("id"))
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
