// Package events serves 覺行共修活動 and the 菩提幣 they earn: anyone signed up can
// start an activity of three or more people, others join as participants or helpers
// (協辦志工、減壓教練), and after it ends the organizer submits it for review. When a
// 超級管理員 approves, the organizer and helpers are credited in the wallet ledger.
package events

import (
	"errors"
	"math"
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

// MinPeople is the smallest group that counts as a 覺行小組 activity.
const MinPeople = 3

type Event struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	IsOnline      bool      `json:"is_online"`
	Location      string    `json:"location"`
	StartsAt      time.Time `json:"starts_at"`
	EndsAt        time.Time `json:"ends_at"`
	Capacity      int       `json:"capacity"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	OrganizerName string    `json:"organizer_name"`
	Joined        int       `json:"joined"`
	VenueID       *string   `json:"venue_id"`
	VenueName     string    `json:"venue_name"`
	// Set only on the signed-in user's own list.
	MyRole      string `json:"my_role,omitempty"`
	ClaimStatus string `json:"claim_status,omitempty"`
	ClaimNote   string `json:"claim_note,omitempty"`
}

const eventSelect = `
	SELECT e.id, e.title, e.is_online, e.location, e.starts_at, e.ends_at, e.capacity, e.description, e.status,
	       u.display_name, (SELECT count(*) FROM event_participant p WHERE p.event_id = e.id),
	       e.venue_id, coalesce(v.name, '')
	FROM practice_event e JOIN volunteer o ON o.id = e.organizer_id JOIN app_user u ON u.id = o.user_id
	LEFT JOIN venue v ON v.id = e.venue_id`

func scanEvent(row pgx.CollectableRow) (Event, error) {
	var e Event
	err := row.Scan(&e.ID, &e.Title, &e.IsOnline, &e.Location, &e.StartsAt, &e.EndsAt, &e.Capacity, &e.Description, &e.Status,
		&e.OrganizerName, &e.Joined, &e.VenueID, &e.VenueName)
	return e, err
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
	r.GET("/events", h.publicEvents)

	me := r.Group("/me", h.Auth.RequireUser())
	me.GET("/events", h.myEvents)
	me.POST("/events", h.createEvent)
	me.POST("/events/:id/join", h.joinEvent)
	me.DELETE("/events/:id/join", h.leaveEvent)
	me.POST("/events/:id/cancel", h.cancelEvent)
	me.POST("/events/:id/claim", h.submitClaim)
	me.GET("/wallet", h.myWallet)

	admin := r.Group("/admin", h.Auth.Require(auth.RoleAllianceAdmin, auth.ScopeAlliance))
	admin.GET("/claims", h.listClaims)
	admin.PATCH("/claims/:id", h.reviewClaim)
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

// userError is an error whose message is safe to show as is.
type userError struct {
	status int
	msg    string
}

func (e userError) Error() string { return e.msg }

func respond(c *gin.Context, err error) bool {
	var ue userError
	switch {
	case err == nil:
		return true
	case errors.As(err, &ue):
		httpx.Error(c, ue.status, ue.msg)
	default:
		fail(c, err)
	}
	return false
}

// myVolunteerID is the signed-in member, who must have joined 覺行小組: only they get
// the wallet and activities. Leaving or cancelling an activity uses anyMemberID instead.
func (h *Handler) myVolunteerID(c *gin.Context) (string, bool) {
	return h.memberID(c, true)
}

func (h *Handler) anyMemberID(c *gin.Context) (string, bool) {
	return h.memberID(c, false)
}

func (h *Handler) memberID(c *gin.Context, needGroups bool) (string, bool) {
	var id string
	var inGroups bool
	err := h.DB.QueryRow(c, `SELECT id, in_groups FROM volunteer WHERE user_id = $1`, auth.CurrentUser(c).ID).Scan(&id, &inGroups)
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(c, http.StatusNotFound, "請先在覺行小組頁報名，建立你的帳號")
		return "", false
	}
	if err != nil {
		fail(c, err)
		return "", false
	}
	if needGroups && !inGroups {
		httpx.Error(c, http.StatusForbidden, "請先在個人頁選擇加入覺行小組，才能使用菩提幣錢包與共修活動")
		return "", false
	}
	return id, true
}

// publicEvents lists open activities that have not ended yet, soonest first.
func (h *Handler) publicEvents(c *gin.Context) {
	rows, err := h.DB.Query(c, eventSelect+`
		WHERE e.status = 'open' AND e.ends_at > $1 ORDER BY e.starts_at LIMIT 200`, h.now())
	if err != nil {
		fail(c, err)
		return
	}
	list, err := pgx.CollectRows(rows, scanEvent)
	if err != nil {
		fail(c, err)
		return
	}
	// 線上活動的連結只給報名的人看
	for i := range list {
		if list[i].IsOnline {
			list[i].Location = ""
		}
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) myEvents(c *gin.Context) {
	vid, ok := h.myVolunteerID(c)
	if !ok {
		return
	}
	rows, err := h.DB.Query(c, `
		SELECT e.id, e.title, e.is_online, e.location, e.starts_at, e.ends_at, e.capacity, e.description, e.status,
		       u.display_name, (SELECT count(*) FROM event_participant p WHERE p.event_id = e.id),
		       e.venue_id, coalesce(v.name, ''), me.role, coalesce(cl.status, ''), coalesce(cl.review_note, '')
		FROM event_participant me
		JOIN practice_event e ON e.id = me.event_id
		JOIN volunteer o ON o.id = e.organizer_id JOIN app_user u ON u.id = o.user_id
		LEFT JOIN venue v ON v.id = e.venue_id
		LEFT JOIN coin_claim cl ON cl.event_id = e.id
		WHERE me.volunteer_id = $1
		ORDER BY e.starts_at DESC LIMIT 200`, vid)
	if err != nil {
		fail(c, err)
		return
	}
	list, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Event, error) {
		var e Event
		err := row.Scan(&e.ID, &e.Title, &e.IsOnline, &e.Location, &e.StartsAt, &e.EndsAt, &e.Capacity, &e.Description, &e.Status,
			&e.OrganizerName, &e.Joined, &e.VenueID, &e.VenueName, &e.MyRole, &e.ClaimStatus, &e.ClaimNote)
		return e, err
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) createEvent(c *gin.Context) {
	var req struct {
		Title       string    `json:"title" binding:"required,max=100"`
		IsOnline    bool      `json:"is_online"`
		Location    string    `json:"location" binding:"max=500"`
		StartsAt    time.Time `json:"starts_at" binding:"required"`
		EndsAt      time.Time `json:"ends_at" binding:"required"`
		Capacity    int       `json:"capacity"`
		Description string    `json:"description" binding:"max=2000"`
		VenueID     string    `json:"venue_id" binding:"omitempty,uuid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "請填寫活動名稱、開始與結束時間")
		return
	}
	req.Title, req.Location = strings.TrimSpace(req.Title), strings.TrimSpace(req.Location)
	// 線下活動可以選一個場域；沒填地點就用場域地址
	var venue *string
	if req.VenueID != "" && !req.IsOnline {
		var address string
		err := h.DB.QueryRow(c, `
			SELECT v.address FROM venue v JOIN center ce ON ce.id = v.center_id
			WHERE v.id = $1 AND v.status = 'active' AND ce.status = 'active'`, req.VenueID).Scan(&address)
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.Error(c, http.StatusBadRequest, "選的場域不存在或已停用")
			return
		}
		if err != nil {
			fail(c, err)
			return
		}
		venue = &req.VenueID
		if req.Location == "" {
			req.Location = address
		}
	}
	switch {
	case req.Title == "":
		httpx.Error(c, http.StatusBadRequest, "請填寫活動名稱")
		return
	case req.Location == "":
		httpx.Error(c, http.StatusBadRequest, "請填寫地點；線上活動請寫會議連結或集合方式")
		return
	case !req.EndsAt.After(req.StartsAt):
		httpx.Error(c, http.StatusBadRequest, "結束時間要晚於開始時間")
		return
	case req.EndsAt.Sub(req.StartsAt) > 7*24*time.Hour:
		httpx.Error(c, http.StatusBadRequest, "一次活動最長七天")
		return
	case req.EndsAt.Before(h.now()):
		httpx.Error(c, http.StatusBadRequest, "活動時間已經過了")
		return
	case req.Capacity < MinPeople || req.Capacity > 500:
		httpx.Error(c, http.StatusBadRequest, "開放人數要在 3 到 500 人之間（含發起人）")
		return
	}
	vid, ok := h.myVolunteerID(c)
	if !ok {
		return
	}
	var id string
	err := pgx.BeginFunc(c, h.DB, func(tx pgx.Tx) error {
		if err := tx.QueryRow(c, `
			INSERT INTO practice_event (organizer_id, title, is_online, location, starts_at, ends_at, capacity, description, venue_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
			vid, req.Title, req.IsOnline, req.Location, req.StartsAt, req.EndsAt, req.Capacity, strings.TrimSpace(req.Description), venue).Scan(&id); err != nil {
			return err
		}
		_, err := tx.Exec(c, `INSERT INTO event_participant (event_id, volunteer_id, role) VALUES ($1, $2, 'organizer')`, id, vid)
		return err
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) joinEvent(c *gin.Context) {
	var req struct {
		Role string `json:"role"`
	}
	c.ShouldBindJSON(&req)
	if req.Role == "" {
		req.Role = "participant"
	}
	if req.Role != "participant" && req.Role != "helper" {
		httpx.Error(c, http.StatusBadRequest, "請選擇參加或協辦")
		return
	}
	vid, ok := h.myVolunteerID(c)
	if !ok {
		return
	}
	err := pgx.BeginFunc(c, h.DB, func(tx pgx.Tx) error {
		var status string
		var endsAt time.Time
		var capacity, joined int
		// Lock the event so two last-seat joins cannot both get in.
		err := tx.QueryRow(c, `SELECT status, ends_at, capacity FROM practice_event WHERE id = $1 FOR UPDATE`, c.Param("id")).
			Scan(&status, &endsAt, &capacity)
		if err != nil {
			return err
		}
		var current string
		err = tx.QueryRow(c, `SELECT role FROM event_participant WHERE event_id = $1 AND volunteer_id = $2`, c.Param("id"), vid).Scan(&current)
		switch {
		case err == nil && current == "organizer":
			return userError{http.StatusConflict, "你是這場活動的發起人"}
		case err == nil:
			_, err = tx.Exec(c, `UPDATE event_participant SET role = $3 WHERE event_id = $1 AND volunteer_id = $2`, c.Param("id"), vid, req.Role)
			return err
		case !errors.Is(err, pgx.ErrNoRows):
			return err
		}
		if status != "open" || !endsAt.After(h.now()) {
			return userError{http.StatusConflict, "這場活動已經結束或取消"}
		}
		if err := tx.QueryRow(c, `SELECT count(*) FROM event_participant WHERE event_id = $1`, c.Param("id")).Scan(&joined); err != nil {
			return err
		}
		if joined >= capacity {
			return userError{http.StatusConflict, "這場活動已經額滿"}
		}
		_, err = tx.Exec(c, `INSERT INTO event_participant (event_id, volunteer_id, role) VALUES ($1, $2, $3)`, c.Param("id"), vid, req.Role)
		return err
	})
	if respond(c, err) {
		c.Status(http.StatusNoContent)
	}
}

func (h *Handler) leaveEvent(c *gin.Context) {
	vid, ok := h.anyMemberID(c)
	if !ok {
		return
	}
	// Once the organizer has submitted, the list of helpers is fixed.
	tag, err := h.DB.Exec(c, `
		DELETE FROM event_participant p WHERE p.event_id = $1 AND p.volunteer_id = $2 AND p.role <> 'organizer'
		  AND NOT EXISTS (SELECT 1 FROM coin_claim cl WHERE cl.event_id = p.event_id AND cl.status <> 'rejected')`, c.Param("id"), vid)
	if err != nil {
		fail(c, err)
		return
	}
	if tag.RowsAffected() == 0 {
		httpx.Error(c, http.StatusConflict, "發起人不能退出；已經送審的活動也不能退出")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) cancelEvent(c *gin.Context) {
	vid, ok := h.anyMemberID(c)
	if !ok {
		return
	}
	tag, err := h.DB.Exec(c, `
		UPDATE practice_event e SET status = 'cancelled', updated_at = now()
		WHERE e.id = $1 AND e.organizer_id = $2 AND e.status = 'open'
		  AND NOT EXISTS (SELECT 1 FROM coin_claim cl WHERE cl.event_id = e.id)`, c.Param("id"), vid)
	if err != nil {
		fail(c, err)
		return
	}
	if tag.RowsAffected() == 0 {
		httpx.Error(c, http.StatusConflict, "只有發起人能取消，已送審的活動不能取消")
		return
	}
	c.Status(http.StatusNoContent)
}

// SuggestedAmount follows the initial tier table (菩提幣介紹): 300 per hour under half a
// day, 1,500 from four hours, 3,000 from eight. The reviewer can change it on approval.
func SuggestedAmount(d time.Duration) int64 {
	hours := d.Hours()
	switch {
	case hours >= 8:
		return 3000
	case hours >= 4:
		return 1500
	default:
		return 300 * int64(math.Max(1, math.Ceil(hours)))
	}
}

func (h *Handler) submitClaim(c *gin.Context) {
	var req struct {
		Attendance int    `json:"attendance"`
		Report     string `json:"report" binding:"max=2000"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "請填寫實際出席人數")
		return
	}
	if req.Attendance < MinPeople {
		httpx.Error(c, http.StatusBadRequest, "實際出席要三人以上（含發起人）才能申請菩提幣")
		return
	}
	vid, ok := h.myVolunteerID(c)
	if !ok {
		return
	}
	err := pgx.BeginFunc(c, h.DB, func(tx pgx.Tx) error {
		var organizer, status string
		var startsAt, endsAt time.Time
		err := tx.QueryRow(c, `SELECT organizer_id, status, starts_at, ends_at FROM practice_event WHERE id = $1 FOR UPDATE`, c.Param("id")).
			Scan(&organizer, &status, &startsAt, &endsAt)
		if err != nil {
			return err
		}
		switch {
		case organizer != vid:
			return userError{http.StatusForbidden, "只有發起人能送審"}
		case status != "open":
			return userError{http.StatusConflict, "已取消的活動不能送審"}
		case endsAt.After(h.now()):
			return userError{http.StatusConflict, "活動結束後才能送審"}
		}
		var claimID, claimStatus string
		err = tx.QueryRow(c, `SELECT id, status FROM coin_claim WHERE event_id = $1`, c.Param("id")).Scan(&claimID, &claimStatus)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			if err := tx.QueryRow(c, `
				INSERT INTO coin_claim (event_id, submitted_by, attendance, report) VALUES ($1, $2, $3, $4) RETURNING id`,
				c.Param("id"), vid, req.Attendance, strings.TrimSpace(req.Report)).Scan(&claimID); err != nil {
				return err
			}
		case err != nil:
			return err
		case claimStatus != "rejected":
			return userError{http.StatusConflict, "這場活動已經送審過了"}
		default:
			if _, err := tx.Exec(c, `
				UPDATE coin_claim SET attendance = $2, report = $3, status = 'submitted', review_note = '',
				       reviewed_by = NULL, reviewed_at = NULL, updated_at = now() WHERE id = $1`,
				claimID, req.Attendance, strings.TrimSpace(req.Report)); err != nil {
				return err
			}
			if _, err := tx.Exec(c, `DELETE FROM claim_recipient WHERE claim_id = $1`, claimID); err != nil {
				return err
			}
		}
		_, err = tx.Exec(c, `
			INSERT INTO claim_recipient (claim_id, volunteer_id, role, amount)
			SELECT $1, volunteer_id, role, $3 FROM event_participant WHERE event_id = $2 AND role IN ('organizer', 'helper')`,
			claimID, c.Param("id"), SuggestedAmount(endsAt.Sub(startsAt)))
		return err
	})
	if respond(c, err) {
		c.Status(http.StatusNoContent)
	}
}

type LedgerEntry struct {
	Amount    int64     `json:"amount"`
	Kind      string    `json:"kind"`
	Memo      string    `json:"memo"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) myWallet(c *gin.Context) {
	vid, ok := h.myVolunteerID(c)
	if !ok {
		return
	}
	var balance int64
	if err := h.DB.QueryRow(c, `SELECT coalesce(sum(amount), 0) FROM coin_ledger WHERE volunteer_id = $1`, vid).Scan(&balance); err != nil {
		fail(c, err)
		return
	}
	rows, err := h.DB.Query(c, `SELECT amount, kind, memo, created_at FROM coin_ledger WHERE volunteer_id = $1 ORDER BY created_at DESC LIMIT 100`, vid)
	if err != nil {
		fail(c, err)
		return
	}
	entries, err := pgx.CollectRows(rows, pgx.RowToStructByPos[LedgerEntry])
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance, "entries": entries})
}

type Recipient struct {
	VolunteerID string `json:"volunteer_id"`
	LegalName   string `json:"legal_name"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Verified    bool   `json:"verified"`
	Role        string `json:"role"`
	Amount      int64  `json:"amount"`
}

type Claim struct {
	ID            string      `json:"id"`
	Status        string      `json:"status"`
	Attendance    int         `json:"attendance"`
	Report        string      `json:"report"`
	ReviewNote    string      `json:"review_note"`
	CreatedAt     time.Time   `json:"created_at"`
	ReviewedAt    *time.Time  `json:"reviewed_at"`
	Event         Event       `json:"event"`
	Participants  int         `json:"participants"`
	Recipients    []Recipient `json:"recipients"`
	OrganizerName string      `json:"organizer_legal_name"`
}

func (h *Handler) listClaims(c *gin.Context) {
	rows, err := h.DB.Query(c, `
		SELECT cl.id, cl.status, cl.attendance, cl.report, cl.review_note, cl.created_at, cl.reviewed_at,
		       e.id, e.title, e.is_online, e.location, e.starts_at, e.ends_at, e.capacity, e.description, e.status, u.display_name,
		       (SELECT count(*) FROM event_participant p WHERE p.event_id = e.id), o.legal_name,
		       coalesce((SELECT json_agg(json_build_object('volunteer_id', v.id, 'legal_name', v.legal_name, 'display_name', ru.display_name,
		                   'email', ru.email, 'verified', v.status = 'verified', 'role', r.role, 'amount', r.amount) ORDER BY r.role DESC, v.legal_name)
		                 FROM claim_recipient r JOIN volunteer v ON v.id = r.volunteer_id JOIN app_user ru ON ru.id = v.user_id
		                 WHERE r.claim_id = cl.id), '[]')
		FROM coin_claim cl JOIN practice_event e ON e.id = cl.event_id
		JOIN volunteer o ON o.id = e.organizer_id JOIN app_user u ON u.id = o.user_id
		WHERE ($1 = '' OR cl.status = $1)
		ORDER BY cl.status = 'submitted' DESC, cl.created_at DESC LIMIT 500`, c.Query("status"))
	if err != nil {
		fail(c, err)
		return
	}
	list, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Claim, error) {
		var cl Claim
		e := &cl.Event
		err := row.Scan(&cl.ID, &cl.Status, &cl.Attendance, &cl.Report, &cl.ReviewNote, &cl.CreatedAt, &cl.ReviewedAt,
			&e.ID, &e.Title, &e.IsOnline, &e.Location, &e.StartsAt, &e.EndsAt, &e.Capacity, &e.Description, &e.Status, &e.OrganizerName,
			&cl.Participants, &cl.OrganizerName, &cl.Recipients)
		e.Joined = cl.Participants
		return cl, err
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) reviewClaim(c *gin.Context) {
	var req struct {
		Status     string           `json:"status" binding:"required,oneof=approved rejected"`
		ReviewNote string           `json:"review_note" binding:"max=2000"`
		Amounts    map[string]int64 `json:"amounts"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "請選擇核准或退回")
		return
	}
	note := strings.TrimSpace(req.ReviewNote)
	if req.Status == "rejected" && note == "" {
		httpx.Error(c, http.StatusBadRequest, "退回時請寫原因，發起人會看到")
		return
	}
	for _, a := range req.Amounts {
		if a < 0 || a > 100_000 {
			httpx.Error(c, http.StatusBadRequest, "每人核發數量要在 0 到 100,000 之間")
			return
		}
	}
	reviewer := auth.CurrentUser(c).ID
	err := pgx.BeginFunc(c, h.DB, func(tx pgx.Tx) error {
		var status, title string
		err := tx.QueryRow(c, `
			SELECT cl.status, e.title FROM coin_claim cl JOIN practice_event e ON e.id = cl.event_id
			WHERE cl.id = $1 FOR UPDATE OF cl`, c.Param("id")).Scan(&status, &title)
		if err != nil {
			return err
		}
		if status != "submitted" {
			return userError{http.StatusConflict, "這筆申請已經審核過了"}
		}
		if _, err := tx.Exec(c, `
			UPDATE coin_claim SET status = $2, review_note = $3, reviewed_by = $4, reviewed_at = now(), updated_at = now()
			WHERE id = $1`, c.Param("id"), req.Status, note, reviewer); err != nil {
			return err
		}
		if req.Status == "rejected" {
			return nil
		}
		for vid, amount := range req.Amounts {
			tag, err := tx.Exec(c, `UPDATE claim_recipient SET amount = $3 WHERE claim_id = $1 AND volunteer_id::text = $2`, c.Param("id"), vid, amount)
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				return userError{http.StatusBadRequest, "核發對象不在這筆申請裡"}
			}
		}
		_, err = tx.Exec(c, `
			INSERT INTO coin_ledger (volunteer_id, amount, kind, claim_id, memo)
			SELECT volunteer_id, amount, 'issue', claim_id, $2 || CASE role WHEN 'organizer' THEN '（發起）' ELSE '（協辦）' END
			FROM claim_recipient WHERE claim_id = $1 AND amount > 0`, c.Param("id"), title)
		return err
	})
	if respond(c, err) {
		c.Status(http.StatusNoContent)
	}
}
