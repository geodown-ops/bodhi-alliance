// Package members serves volunteer accounts (志工註冊、我的資料、加入覺行小組) and the
// back-office 志工名冊, where a center's admins verify the volunteers who name it as
// their home center.
package members

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

type Membership struct {
	GroupID  string    `json:"group_id"`
	Name     string    `json:"name"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type Volunteer struct {
	ID           string       `json:"id"`
	UserID       string       `json:"user_id"`
	Email        string       `json:"email"`
	DisplayName  string       `json:"display_name"`
	LegalName    string       `json:"legal_name"`
	Phone        string       `json:"phone"`
	LineID       string       `json:"line_id"`
	HomeCenterID string       `json:"home_center_id"`
	CenterName   string       `json:"center_name"`
	WantsCoach   bool         `json:"wants_coach"`
	IsCoach      bool         `json:"is_coach"`
	Status       string       `json:"status"`
	ReviewNote   string       `json:"review_note"`
	VerifiedAt   *time.Time   `json:"verified_at"`
	Frozen       bool         `json:"frozen"`
	CreatedAt    time.Time    `json:"created_at"`
	// 系統會員可以加入覺行小組、世界佛教教育協會，或兩者都加入
	InGroups            bool         `json:"in_groups"`
	InAssociation       bool         `json:"in_association"`
	AssociationJoinedAt *time.Time   `json:"association_joined_at"`
	Groups              []Membership `json:"groups"`
}

const volunteerSelect = `
	SELECT v.id, v.user_id, u.email, u.display_name, v.legal_name, v.phone, v.line_id, coalesce(v.home_center_id::text, ''), coalesce(c.name, ''),
	       v.wants_coach, v.is_coach, v.status, v.review_note, v.verified_at, v.qr_frozen_at IS NOT NULL, v.created_at,
	       v.in_groups, v.in_association, v.association_joined_at,
	       coalesce((SELECT json_agg(json_build_object('group_id', g.id, 'name', g.name, 'role', m.role, 'joined_at', m.joined_at) ORDER BY m.joined_at)
	                 FROM group_member m JOIN practice_group g ON g.id = m.group_id WHERE m.volunteer_id = v.id), '[]')
	FROM volunteer v JOIN app_user u ON u.id = v.user_id LEFT JOIN center c ON c.id = v.home_center_id`

func scanVolunteer(row pgx.CollectableRow) (Volunteer, error) {
	var v Volunteer
	err := row.Scan(&v.ID, &v.UserID, &v.Email, &v.DisplayName, &v.LegalName, &v.Phone, &v.LineID, &v.HomeCenterID, &v.CenterName,
		&v.WantsCoach, &v.IsCoach, &v.Status, &v.ReviewNote, &v.VerifiedAt, &v.Frozen, &v.CreatedAt,
		&v.InGroups, &v.InAssociation, &v.AssociationJoinedAt, &v.Groups)
	return v, err
}

type Handler struct {
	DB   *pgxpool.Pool
	Auth *auth.Service
}

func (h *Handler) Routes(r *gin.RouterGroup) {
	r.GET("/centers", h.publicCenters)
	r.POST("/volunteers", httpx.RateLimit(5, time.Minute), h.register)

	me := r.Group("/me", h.Auth.RequireUser())
	me.GET("/volunteer", h.myProfile)
	me.PUT("/volunteer", h.updateMyProfile)
	me.PUT("/memberships", h.setMemberships)
	me.POST("/groups/:id", h.joinGroup)
	me.DELETE("/groups/:id", h.leaveGroup)

	staff := r.Group("/admin", h.Auth.RequireCenterStaff())
	staff.GET("/volunteers", h.listVolunteers)
	staff.PATCH("/volunteers/:id", h.reviewVolunteer)
	staff.GET("/groups/:id/members", h.listMembers)
	staff.PUT("/groups/:id/members/:vid", h.setMemberRole)
	staff.DELETE("/groups/:id/members/:vid", h.removeMember)
}

func fail(c *gin.Context, err error) {
	var pg *pgconn.PgError
	switch {
	case errors.Is(err, pgx.ErrNoRows), errors.As(err, &pg) && pg.Code == "22P02":
		httpx.Error(c, http.StatusNotFound, "找不到這筆資料")
	default:
		httpx.Error(c, http.StatusInternalServerError, "處理失敗，請稍後再試")
	}
}

func (h *Handler) publicCenters(c *gin.Context) {
	rows, err := h.DB.Query(c, `SELECT id, name, region FROM center WHERE status = 'active' ORDER BY region, name`)
	if err != nil {
		fail(c, err)
		return
	}
	type center struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Region string `json:"region"`
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[center])
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) register(c *gin.Context) {
	var req struct {
		Email        string `json:"email" binding:"required,email,max=200"`
		Password     string `json:"password" binding:"required,max=200"`
		DisplayName  string `json:"display_name" binding:"max=50"`
		LegalName    string `json:"legal_name" binding:"required,max=100"`
		Phone        string `json:"phone" binding:"max=50"`
		LineID       string `json:"line_id" binding:"max=100"`
		HomeCenterID string `json:"home_center_id" binding:"omitempty,uuid"`
		WantsCoach   bool   `json:"wants_coach"`
		// 沒指定時當作加入覺行小組（舊的報名頁）
		InGroups      *bool  `json:"in_groups"`
		InAssociation bool   `json:"in_association"`
		Website       string `json:"website"` // honeypot
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "請填寫真實姓名、正確的電子郵件與密碼")
		return
	}
	if req.Website != "" {
		c.JSON(http.StatusCreated, gin.H{})
		return
	}
	req.LegalName = strings.TrimSpace(req.LegalName)
	if req.LegalName == "" {
		httpx.Error(c, http.StatusBadRequest, "請填寫真實姓名")
		return
	}
	inGroups := req.InGroups == nil || *req.InGroups
	if !inGroups && !req.InAssociation {
		httpx.Error(c, http.StatusBadRequest, "請至少選擇加入覺行小組或世界佛教教育協會其中一個")
		return
	}
	// 所屬中心可以不選；選了就要是營運中的中心
	var center *string
	if req.HomeCenterID != "" {
		center = &req.HomeCenterID
	}
	if strings.TrimSpace(req.DisplayName) == "" {
		req.DisplayName = req.LegalName
	}
	var userID string
	err := pgx.BeginFunc(c, h.DB, func(tx pgx.Tx) error {
		if center != nil {
			var active bool
			if err := tx.QueryRow(c, `SELECT EXISTS (SELECT 1 FROM center WHERE id = $1 AND status = 'active')`, *center).Scan(&active); err != nil {
				return err
			}
			if !active {
				return errNoCenter
			}
		}
		var err error
		if userID, err = auth.CreateUserTx(c, tx, req.Email, req.DisplayName, req.Password); err != nil {
			return err
		}
		_, err = tx.Exec(c, `
			INSERT INTO volunteer (user_id, home_center_id, legal_name, phone, line_id, wants_coach, in_groups, in_association, association_joined_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CASE WHEN $8 THEN now() END)`,
			userID, center, req.LegalName, strings.TrimSpace(req.Phone), strings.TrimSpace(req.LineID), req.WantsCoach, inGroups, req.InAssociation)
		return err
	})
	var pg *pgconn.PgError
	switch {
	case errors.Is(err, errNoCenter):
		httpx.Error(c, http.StatusBadRequest, "選的中心不存在或已暫停")
		return
	case errors.As(err, &pg) && pg.Code == "23505":
		httpx.Error(c, http.StatusConflict, "這個電子郵件已經註冊過，請直接登入")
		return
	case err != nil && strings.HasPrefix(err.Error(), "密碼"):
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	case err != nil:
		fail(c, err)
		return
	}
	token, err := h.Auth.NewSession(c, userID)
	if err != nil {
		fail(c, err)
		return
	}
	u, err := h.Auth.UserForToken(c, token)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"token": token, "user": u})
}

var errNoCenter = errors.New("no such center")

func (h *Handler) myProfile(c *gin.Context) {
	rows, err := h.DB.Query(c, volunteerSelect+` WHERE v.user_id = $1`, auth.CurrentUser(c).ID)
	if err != nil {
		fail(c, err)
		return
	}
	v, err := pgx.CollectExactlyOneRow(rows, scanVolunteer)
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(c, http.StatusNotFound, "這個帳號還不是志工")
		return
	}
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *Handler) updateMyProfile(c *gin.Context) {
	var req struct {
		DisplayName string `json:"display_name" binding:"required,max=50"`
		LegalName   string `json:"legal_name" binding:"required,max=100"`
		Phone       string `json:"phone" binding:"max=50"`
		LineID      string `json:"line_id" binding:"max=100"`
		WantsCoach  bool   `json:"wants_coach"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "請填寫暱稱與真實姓名")
		return
	}
	uid := auth.CurrentUser(c).ID
	err := pgx.BeginFunc(c, h.DB, func(tx pgx.Tx) error {
		// 實名核可後就不能自己改姓名，要請中心管理員處理，避免核可的人與領幣的人不同。
		tag, err := tx.Exec(c, `
			UPDATE volunteer SET phone = $2, wants_coach = $3, line_id = $5, updated_at = now(),
			       legal_name = CASE WHEN status = 'verified' THEN legal_name ELSE $4 END,
			       status = CASE WHEN status = 'rejected' THEN 'pending' ELSE status END
			WHERE user_id = $1`, uid, strings.TrimSpace(req.Phone), req.WantsCoach, strings.TrimSpace(req.LegalName), strings.TrimSpace(req.LineID))
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
		_, err = tx.Exec(c, `UPDATE app_user SET display_name = $2 WHERE id = $1`, uid, strings.TrimSpace(req.DisplayName))
		return err
	})
	if err != nil {
		fail(c, err)
		return
	}
	h.myProfile(c)
}

// setMemberships switches the signed-in member in or out of 覺行小組 and the
// association. Leaving 覺行小組 keeps the wallet ledger; it only hides the wallet and
// activities until they join again.
func (h *Handler) setMemberships(c *gin.Context) {
	var req struct {
		InGroups      bool `json:"in_groups"`
		InAssociation bool `json:"in_association"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "請選擇要加入的會員身分")
		return
	}
	if !req.InGroups && !req.InAssociation {
		httpx.Error(c, http.StatusBadRequest, "至少要保留覺行小組或世界佛教教育協會其中一個會員身分")
		return
	}
	tag, err := h.DB.Exec(c, `
		UPDATE volunteer SET in_groups = $2, in_association = $3, updated_at = now(),
		       association_joined_at = CASE WHEN $3 THEN coalesce(association_joined_at, now()) END
		WHERE user_id = $1`, auth.CurrentUser(c).ID, req.InGroups, req.InAssociation)
	if err == nil && tag.RowsAffected() == 0 {
		err = pgx.ErrNoRows
	}
	if err != nil {
		fail(c, err)
		return
	}
	h.myProfile(c)
}

func (h *Handler) myVolunteerID(c *gin.Context) (string, bool) {
	var id string
	err := h.DB.QueryRow(c, `SELECT id FROM volunteer WHERE user_id = $1`, auth.CurrentUser(c).ID).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(c, http.StatusNotFound, "這個帳號還不是志工")
		return "", false
	}
	if err != nil {
		fail(c, err)
		return "", false
	}
	return id, true
}

func (h *Handler) joinGroup(c *gin.Context) {
	vid, ok := h.myVolunteerID(c)
	if !ok {
		return
	}
	tag, err := h.DB.Exec(c, `
		INSERT INTO group_member (group_id, volunteer_id)
		SELECT id, $2 FROM practice_group WHERE id = $1 AND is_listed
		ON CONFLICT DO NOTHING`, c.Param("id"), vid)
	if err != nil {
		fail(c, err)
		return
	}
	if tag.RowsAffected() == 0 {
		// Either already a member, or no such listed group.
		var member bool
		h.DB.QueryRow(c, `SELECT EXISTS (SELECT 1 FROM group_member WHERE group_id::text = $1 AND volunteer_id = $2)`, c.Param("id"), vid).Scan(&member)
		if !member {
			httpx.Error(c, http.StatusNotFound, "找不到這個小組")
			return
		}
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) leaveGroup(c *gin.Context) {
	vid, ok := h.myVolunteerID(c)
	if !ok {
		return
	}
	if _, err := h.DB.Exec(c, `DELETE FROM group_member WHERE group_id::text = $1 AND volunteer_id = $2`, c.Param("id"), vid); err != nil {
		fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// visibleCenters returns nil for an alliance admin (every center) or the centers a
// center admin manages.
func visibleCenters(u *auth.User) []string {
	if u.IsAllianceAdmin() {
		return nil
	}
	return u.CenterIDs()
}

func (h *Handler) listVolunteers(c *gin.Context) {
	centers := visibleCenters(auth.CurrentUser(c))
	rows, err := h.DB.Query(c, volunteerSelect+`
		WHERE ($1::uuid[] IS NULL OR v.home_center_id = ANY($1))
		  AND ($2 = '' OR v.home_center_id::text = $2)
		  AND ($3 = '' OR v.status = $3)
		ORDER BY v.status = 'pending' DESC, v.created_at DESC LIMIT 1000`,
		centers, c.Query("center_id"), c.Query("status"))
	if err != nil {
		fail(c, err)
		return
	}
	list, err := pgx.CollectRows(rows, scanVolunteer)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

// canManage checks that the signed-in staff member may manage volunteers of center.
func canManage(c *gin.Context, center string) bool {
	u := auth.CurrentUser(c)
	if u.IsAllianceAdmin() || u.ManagesCenter(center) {
		return true
	}
	httpx.Error(c, http.StatusForbidden, "只能管理自己中心的志工")
	return false
}

func (h *Handler) reviewVolunteer(c *gin.Context) {
	var req struct {
		Status     string `json:"status" binding:"required,oneof=pending verified rejected"`
		IsCoach    bool   `json:"is_coach"`
		ReviewNote string `json:"review_note" binding:"max=2000"`
		Frozen     bool   `json:"frozen"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "狀態不正確")
		return
	}
	if req.Status == "rejected" && strings.TrimSpace(req.ReviewNote) == "" {
		httpx.Error(c, http.StatusBadRequest, "退回時請寫原因，志工會看到")
		return
	}
	var center string
	if err := h.DB.QueryRow(c, `SELECT coalesce(home_center_id::text, '') FROM volunteer WHERE id = $1`, c.Param("id")).Scan(&center); err != nil {
		fail(c, err)
		return
	}
	if !canManage(c, center) {
		return
	}
	_, err := h.DB.Exec(c, `
		UPDATE volunteer SET status = $2, is_coach = $3, review_note = $4, updated_at = now(),
		       verified_by = CASE WHEN $2 = 'verified' THEN coalesce(CASE WHEN status = 'verified' THEN verified_by END, $5) END,
		       verified_at = CASE WHEN $2 = 'verified' THEN coalesce(CASE WHEN status = 'verified' THEN verified_at END, now()) END,
		       qr_frozen_at = CASE WHEN $6 THEN coalesce(qr_frozen_at, now()) END
		WHERE id = $1`,
		c.Param("id"), req.Status, req.IsCoach, strings.TrimSpace(req.ReviewNote), auth.CurrentUser(c).ID, req.Frozen)
	if err != nil {
		fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// groupCenter returns the group's center, or "" when it has none (alliance admins only).
func (h *Handler) groupCenter(c *gin.Context) (string, bool) {
	var center *string
	if err := h.DB.QueryRow(c, `SELECT center_id FROM practice_group WHERE id = $1`, c.Param("id")).Scan(&center); err != nil {
		fail(c, err)
		return "", false
	}
	if center == nil {
		if !auth.CurrentUser(c).IsAllianceAdmin() {
			httpx.Error(c, http.StatusForbidden, "只能管理自己中心的小組")
			return "", false
		}
		return "", true
	}
	return *center, canManage(c, *center)
}

func (h *Handler) listMembers(c *gin.Context) {
	if _, ok := h.groupCenter(c); !ok {
		return
	}
	rows, err := h.DB.Query(c, volunteerSelect+`
		JOIN group_member gm ON gm.volunteer_id = v.id
		WHERE gm.group_id = $1 ORDER BY gm.role = 'leader' DESC, gm.joined_at`, c.Param("id"))
	if err != nil {
		fail(c, err)
		return
	}
	list, err := pgx.CollectRows(rows, scanVolunteer)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) setMemberRole(c *gin.Context) {
	var req struct {
		Role string `json:"role" binding:"required,oneof=member leader"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "角色不正確")
		return
	}
	if _, ok := h.groupCenter(c); !ok {
		return
	}
	tag, err := h.DB.Exec(c, `UPDATE group_member SET role = $3 WHERE group_id = $1 AND volunteer_id = $2`,
		c.Param("id"), c.Param("vid"), req.Role)
	if err == nil && tag.RowsAffected() == 0 {
		err = pgx.ErrNoRows
	}
	if err != nil {
		fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) removeMember(c *gin.Context) {
	if _, ok := h.groupCenter(c); !ok {
		return
	}
	if _, err := h.DB.Exec(c, `DELETE FROM group_member WHERE group_id = $1 AND volunteer_id = $2`, c.Param("id"), c.Param("vid")); err != nil {
		fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
