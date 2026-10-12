package guide

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
	"github.com/geodown-ops/bodhi-alliance/server/internal/httpx"
)

// 使用問題：Sunny 答不了的網站使用問題，或對方說遇到的狀況，由 Sunny 記下，
// 超級管理員在後台回答；會員在個人頁看到回答，沒登入的訪客可以留 email 讓管理員聯絡。

const (
	maxIssueRunes        = 1000
	maxIssuesPerChat     = 5 // 一次對話最多記幾筆，避免被拿來灌資料
	issueTranscriptTurns = 6
)

var errIssueLimit = errors.New("這次對話已經記下很多問題了")

// SiteIssue is one recorded usage problem.
type SiteIssue struct {
	ID           string     `json:"id"`
	UserID       *string    `json:"user_id"`
	MemberName   string     `json:"member_name"`
	MemberEmail  string     `json:"member_email"`
	ContactEmail string     `json:"contact_email"`
	Summary      string     `json:"summary"`
	Transcript   []Message  `json:"transcript"`
	Status       string     `json:"status"`
	Answer       string     `json:"answer"`
	AnsweredBy   string     `json:"answered_by"`
	AnsweredAt   *time.Time `json:"answered_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

// IssueReporter records what the visitor in one chat reports; Sunny calls it through tools.
type IssueReporter struct {
	store        *Store
	conversation string
	userID       *string
	transcript   []Message
}

// Record saves a new issue for this chat.
func (r *IssueReporter) Record(ctx context.Context, summary, email string) error {
	summary = clip(strings.TrimSpace(summary), maxIssueRunes)
	if summary == "" {
		return errors.New("請寫下問題摘要")
	}
	email, err := cleanEmail(email)
	if err != nil {
		return err
	}
	tr := r.transcript
	if len(tr) > issueTranscriptTurns {
		tr = tr[len(tr)-issueTranscriptTurns:]
	}
	raw, _ := json.Marshal(tr)
	return pgx.BeginFunc(ctx, r.store.DB, func(tx pgx.Tx) error {
		if r.conversation != "" {
			var n int
			if err := tx.QueryRow(ctx, `SELECT count(*) FROM guide.site_issue WHERE conversation = $1`, r.conversation).Scan(&n); err != nil {
				return err
			}
			if n >= maxIssuesPerChat {
				return errIssueLimit
			}
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO guide.site_issue (conversation, user_id, contact_email, summary, transcript)
			VALUES ($1, $2, $3, $4, $5)`, r.conversation, r.userID, email, summary, raw)
		return err
	})
}

// SetEmail adds the visitor's contact email to the issues already recorded in this chat.
func (r *IssueReporter) SetEmail(ctx context.Context, email string) error {
	email, err := cleanEmail(email)
	if err != nil {
		return err
	}
	if email == "" {
		return errors.New("沒有收到 email")
	}
	if r.conversation == "" {
		return errors.New("這次對話還沒有記下問題")
	}
	tag, err := r.store.DB.Exec(ctx,
		`UPDATE guide.site_issue SET contact_email = $2 WHERE conversation = $1`, r.conversation, email)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("這次對話還沒有記下問題，請先記下問題")
	}
	return nil
}

func cleanEmail(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", nil
	}
	a, err := mail.ParseAddress(s)
	if err != nil || len(a.Address) > 200 || !strings.Contains(a.Address, "@") {
		return "", errors.New("email 格式不正確，請對方再確認一次")
	}
	return strings.ToLower(a.Address), nil
}

func clip(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n])
	}
	return s
}

const issueColumns = `
	i.id, i.user_id, coalesce(u.display_name, ''), coalesce(u.email, ''), i.contact_email, i.summary,
	i.transcript, i.status, i.answer, coalesce(a.display_name, ''), i.answered_at, i.created_at`

func scanIssues(rows pgx.Rows) ([]SiteIssue, error) {
	out := []SiteIssue{}
	for rows.Next() {
		var s SiteIssue
		var raw []byte
		if err := rows.Scan(&s.ID, &s.UserID, &s.MemberName, &s.MemberEmail, &s.ContactEmail, &s.Summary,
			&raw, &s.Status, &s.Answer, &s.AnsweredBy, &s.AnsweredAt, &s.CreatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(raw, &s.Transcript); err != nil || s.Transcript == nil {
			s.Transcript = []Message{}
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// SiteIssues lists issues, unanswered first, newest first.
func (s *Store) SiteIssues(ctx context.Context) ([]SiteIssue, error) {
	rows, err := s.DB.Query(ctx, `SELECT `+issueColumns+`
		FROM guide.site_issue i
		LEFT JOIN app_user u ON u.id = i.user_id
		LEFT JOIN app_user a ON a.id = i.answered_by
		ORDER BY i.status = 'answered', i.created_at DESC
		LIMIT 500`)
	if err != nil {
		return nil, err
	}
	return scanIssues(rows)
}

// MemberIssues lists one member's issues, newest first.
func (s *Store) MemberIssues(ctx context.Context, userID string) ([]SiteIssue, error) {
	rows, err := s.DB.Query(ctx, `SELECT `+issueColumns+`
		FROM guide.site_issue i
		LEFT JOIN app_user u ON u.id = i.user_id
		LEFT JOIN app_user a ON a.id = i.answered_by
		WHERE i.user_id = $1
		ORDER BY i.created_at DESC
		LIMIT 100`, userID)
	if err != nil {
		return nil, err
	}
	return scanIssues(rows)
}

// AnswerIssue saves (or replaces) the admin's answer.
func (s *Store) AnswerIssue(ctx context.Context, id, answer, userID string) error {
	tag, err := s.DB.Exec(ctx, `
		UPDATE guide.site_issue SET answer = $2, status = 'answered', answered_by = $3, answered_at = now()
		WHERE id::text = $1`, id, answer, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errIssueNotFound
	}
	return nil
}

var errIssueNotFound = errors.New("找不到這個問題")

// optionalUser returns the signed-in user when the request carries a valid token, else nil.
func (h *Handler) optionalUser(c *gin.Context) *auth.User {
	t, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
	if !ok || strings.TrimSpace(t) == "" {
		return nil
	}
	u, err := h.Auth.UserForToken(c, strings.TrimSpace(t))
	if err != nil {
		return nil
	}
	return u
}

func (h *Handler) myIssues(c *gin.Context) {
	list, err := h.Store.MemberIssues(c, auth.CurrentUser(c).ID)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) listIssues(c *gin.Context) {
	list, err := h.Store.SiteIssues(c)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) answerIssue(c *gin.Context) {
	var req struct {
		Answer string `json:"answer"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "格式錯誤")
		return
	}
	answer := strings.TrimSpace(req.Answer)
	if answer == "" || utf8.RuneCountInString(answer) > 4000 {
		httpx.Error(c, http.StatusBadRequest, "請填寫回答（4000 字以內）")
		return
	}
	err := h.Store.AnswerIssue(c, c.Param("id"), answer, auth.CurrentUser(c).ID)
	if errors.Is(err, errIssueNotFound) {
		httpx.Error(c, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "儲存失敗")
		return
	}
	c.Status(http.StatusNoContent)
}
