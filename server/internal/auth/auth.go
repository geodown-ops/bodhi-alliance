// Package auth handles admin login, bearer-token sessions and role checks.
// Volunteer and merchant accounts arrive in M2; phase A only has back-office users.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/geodown-ops/bodhi-alliance/server/internal/httpx"
)

const (
	RoleAllianceAdmin    = "alliance_admin"
	RoleKnowledgeManager = "knowledge_manager"

	ScopeAlliance = "alliance"
	ScopeGuide    = "guide"

	sessionTTL = 7 * 24 * time.Hour
	userKey    = "auth.user"
)

type Role struct {
	Role  string `json:"role"`
	Scope string `json:"scope"`
}

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Roles       []Role `json:"roles"`
}

// Has reports whether the user holds role in scope. Alliance admins hold every role.
func (u *User) Has(role, scope string) bool {
	for _, r := range u.Roles {
		if r.Role == RoleAllianceAdmin && r.Scope == ScopeAlliance {
			return true
		}
		if r.Role == role && r.Scope == scope {
			return true
		}
	}
	return false
}

type Service struct {
	DB *pgxpool.Pool
}

// CreateUser adds a user with the given roles. It fails if the email is taken.
func (s *Service) CreateUser(ctx context.Context, email, name, password string, roles ...Role) (string, error) {
	if len(password) < 10 {
		return "", errors.New("密碼至少 10 個字元")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	var id string
	err = pgx.BeginFunc(ctx, s.DB, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx,
			`INSERT INTO app_user (email, display_name, password_hash) VALUES (lower($1), $2, $3) RETURNING id`,
			email, name, string(hash)).Scan(&id); err != nil {
			return err
		}
		for _, r := range roles {
			if _, err := tx.Exec(ctx,
				`INSERT INTO role_assignment (user_id, role, scope) VALUES ($1, $2, $3)`, id, r.Role, r.Scope); err != nil {
				return err
			}
		}
		return nil
	})
	return id, err
}

// Bootstrap creates the first alliance admin if that email does not exist yet.
func (s *Service) Bootstrap(ctx context.Context, email, password string) error {
	if email == "" || password == "" {
		return nil
	}
	var exists bool
	if err := s.DB.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM app_user WHERE email = lower($1))`, email).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err := s.CreateUser(ctx, email, "聯盟管理員", password, Role{RoleAllianceAdmin, ScopeAlliance})
	return err
}

var (
	errBadLogin  = errors.New("電子郵件或密碼錯誤")
	dummyHash, _ = bcrypt.GenerateFromPassword([]byte("bodhi"), bcrypt.DefaultCost)
)

// Login checks the password and returns a new session token.
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	var id, hash string
	err := s.DB.QueryRow(ctx,
		`SELECT id, password_hash FROM app_user WHERE email = lower($1) AND disabled_at IS NULL`, email).Scan(&id, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		// Spend the same time as a real check so timing does not reveal which emails exist.
		bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
		return "", errBadLogin
	}
	if err != nil {
		return "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return "", errBadLogin
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	_, err = s.DB.Exec(ctx,
		`INSERT INTO user_session (token_hash, user_id, expires_at) VALUES ($1, $2, $3)`,
		hashToken(token), id, time.Now().Add(sessionTTL))
	return token, err
}

func (s *Service) Logout(ctx context.Context, token string) error {
	_, err := s.DB.Exec(ctx, `DELETE FROM user_session WHERE token_hash = $1`, hashToken(token))
	return err
}

// UserForToken returns the session's user, or nil if the token is unknown or expired.
func (s *Service) UserForToken(ctx context.Context, token string) (*User, error) {
	u := &User{}
	err := s.DB.QueryRow(ctx, `
		SELECT u.id, u.email, u.display_name
		FROM user_session s JOIN app_user u ON u.id = s.user_id
		WHERE s.token_hash = $1 AND s.expires_at > now() AND u.disabled_at IS NULL`,
		hashToken(token)).Scan(&u.ID, &u.Email, &u.DisplayName)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rows, err := s.DB.Query(ctx, `SELECT role, scope FROM role_assignment WHERE user_id = $1 ORDER BY role, scope`, u.ID)
	if err != nil {
		return nil, err
	}
	u.Roles, err = pgx.CollectRows(rows, pgx.RowToStructByPos[Role])
	return u, err
}

func hashToken(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}

func bearer(c *gin.Context) string {
	h := c.GetHeader("Authorization")
	if t, ok := strings.CutPrefix(h, "Bearer "); ok {
		return strings.TrimSpace(t)
	}
	return ""
}

// Require lets the request through only if the bearer token's user holds role in scope.
func (s *Service) Require(role, scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := bearer(c)
		if token == "" {
			httpx.Error(c, http.StatusUnauthorized, "請先登入")
			return
		}
		u, err := s.UserForToken(c, token)
		if err != nil {
			httpx.Error(c, http.StatusInternalServerError, "無法驗證登入狀態")
			return
		}
		if u == nil {
			httpx.Error(c, http.StatusUnauthorized, "登入已過期，請重新登入")
			return
		}
		if !u.Has(role, scope) {
			httpx.Error(c, http.StatusForbidden, "沒有這個功能的權限")
			return
		}
		c.Set(userKey, u)
		c.Next()
	}
}

// CurrentUser returns the user set by Require.
func CurrentUser(c *gin.Context) *User {
	u, _ := c.Get(userKey)
	user, _ := u.(*User)
	return user
}

// Routes mounts login, logout and "who am I" under g.
func (s *Service) Routes(g *gin.RouterGroup) {
	g.POST("/login", httpx.RateLimit(10, 30*time.Second), func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			httpx.Error(c, http.StatusBadRequest, "格式錯誤")
			return
		}
		token, err := s.Login(c, req.Email, req.Password)
		if errors.Is(err, errBadLogin) {
			httpx.Error(c, http.StatusUnauthorized, err.Error())
			return
		}
		if err != nil {
			httpx.Error(c, http.StatusInternalServerError, "登入失敗")
			return
		}
		u, err := s.UserForToken(c, token)
		if err != nil {
			httpx.Error(c, http.StatusInternalServerError, "登入失敗")
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token, "user": u})
	})
	g.POST("/logout", func(c *gin.Context) {
		if t := bearer(c); t != "" {
			_ = s.Logout(c, t)
		}
		c.Status(http.StatusNoContent)
	})
	admin := g.Group("/users", s.Require(RoleAllianceAdmin, ScopeAlliance))
	admin.GET("", func(c *gin.Context) {
		rows, err := s.DB.Query(c, `
			SELECT u.id, u.email, u.display_name,
			       coalesce(json_agg(json_build_object('role', r.role, 'scope', r.scope)) FILTER (WHERE r.role IS NOT NULL), '[]')
			FROM app_user u LEFT JOIN role_assignment r ON r.user_id = u.id
			WHERE u.disabled_at IS NULL
			GROUP BY u.id ORDER BY u.created_at`)
		if err != nil {
			httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
			return
		}
		users, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (User, error) {
			var u User
			err := r.Scan(&u.ID, &u.Email, &u.DisplayName, &u.Roles)
			return u, err
		})
		if err != nil {
			httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
			return
		}
		c.JSON(http.StatusOK, users)
	})
	admin.POST("", func(c *gin.Context) {
		var req struct {
			Email       string `json:"email" binding:"required,email"`
			DisplayName string `json:"display_name" binding:"required"`
			Password    string `json:"password" binding:"required"`
			Role        string `json:"role" binding:"required,oneof=alliance_admin knowledge_manager"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			httpx.Error(c, http.StatusBadRequest, "請填寫電子郵件、名稱、密碼與角色")
			return
		}
		scope := ScopeGuide
		if req.Role == RoleAllianceAdmin {
			scope = ScopeAlliance
		}
		id, err := s.CreateUser(c, req.Email, req.DisplayName, req.Password, Role{req.Role, scope})
		if err != nil {
			httpx.Error(c, http.StatusBadRequest, "建立失敗："+userError(err))
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	})
	admin.DELETE("/:id", func(c *gin.Context) {
		if c.Param("id") == CurrentUser(c).ID {
			httpx.Error(c, http.StatusBadRequest, "不能停用自己的帳號")
			return
		}
		_, err := s.DB.Exec(c, `UPDATE app_user SET disabled_at = now() WHERE id = $1`, c.Param("id"))
		if err != nil {
			httpx.Error(c, http.StatusBadRequest, "停用失敗")
			return
		}
		_, _ = s.DB.Exec(c, `DELETE FROM user_session WHERE user_id = $1`, c.Param("id"))
		c.Status(http.StatusNoContent)
	})

	g.GET("/me", func(c *gin.Context) {
		u, err := s.UserForToken(c, bearer(c))
		if err != nil || u == nil {
			httpx.Error(c, http.StatusUnauthorized, "請先登入")
			return
		}
		c.JSON(http.StatusOK, u)
	})
}

func userError(err error) string {
	if strings.Contains(err.Error(), "app_user_email_key") {
		return "這個電子郵件已經有帳號"
	}
	if strings.HasPrefix(err.Error(), "密碼") {
		return err.Error()
	}
	return "請檢查輸入內容"
}
