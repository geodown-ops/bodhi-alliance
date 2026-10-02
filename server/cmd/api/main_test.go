package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
	"github.com/geodown-ops/bodhi-alliance/server/internal/config"
	"github.com/geodown-ops/bodhi-alliance/server/internal/testdb"
)

func call(t *testing.T, r *gin.Engine, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "192.0.2.1:1234"
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func login(t *testing.T, r *gin.Engine, email, password string) string {
	t.Helper()
	w := call(t, r, http.MethodPost, "/api/auth/login", "", map[string]string{"email": email, "password": password})
	if w.Code != http.StatusOK {
		t.Fatalf("login %s: %d %s", email, w.Code, w.Body)
	}
	var out struct{ Token string }
	json.Unmarshal(w.Body.Bytes(), &out)
	return out.Token
}

func TestSignUpsAndReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.New(t)
	authSvc := &auth.Service{DB: pool}
	ctx := context.Background()
	if err := authSvc.Bootstrap(ctx, "Admin@Example.org", "a-long-admin-password"); err != nil {
		t.Fatal(err)
	}
	if err := authSvc.Bootstrap(ctx, "admin@example.org", "ignored-second-time"); err != nil {
		t.Fatal(err)
	}
	r := NewRouter(config.Config{}, authSvc)

	if w := call(t, r, http.MethodPost, "/api/auth/login", "", map[string]string{"email": "admin@example.org", "password": "wrong"}); w.Code != http.StatusUnauthorized {
		t.Errorf("bad password: got %d", w.Code)
	}
	admin := login(t, r, "admin@example.org", "a-long-admin-password")

	// Admin lists a 覺行小組; the public sees it.
	w := call(t, r, http.MethodPost, "/api/admin/groups", admin, map[string]any{
		"name": "台北週三共修", "region": "台北", "schedule": "每週三 19:30", "is_listed": true})
	if w.Code != http.StatusOK {
		t.Fatalf("create group: %d %s", w.Code, w.Body)
	}
	var group struct{ ID string }
	json.Unmarshal(w.Body.Bytes(), &group)
	w = call(t, r, http.MethodGet, "/api/groups", "", nil)
	if !strings.Contains(w.Body.String(), "台北週三共修") {
		t.Errorf("public groups = %s", w.Body)
	}

	// Public sign-ups.
	if w := call(t, r, http.MethodPost, "/api/group-applications", "", map[string]any{
		"group_id": group.ID, "name": "王小明", "email": "ming@example.org", "wants_coach": true}); w.Code != http.StatusCreated {
		t.Fatalf("group application: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodPost, "/api/group-applications", "", map[string]any{"name": "x", "email": "not-an-email"}); w.Code != http.StatusBadRequest {
		t.Errorf("bad email: got %d", w.Code)
	}
	if w := call(t, r, http.MethodPost, "/api/merchant-applications", "", map[string]any{
		"kind": "sponsor", "org_name": "山林民宿", "contact_name": "陳老闆", "email": "inn@example.org", "offerings": "雙人房"}); w.Code != http.StatusCreated {
		t.Fatalf("merchant application: %d %s", w.Code, w.Body)
	}
	// Honeypot: accepted silently, not stored.
	call(t, r, http.MethodPost, "/api/merchant-applications", "", map[string]any{
		"kind": "sponsor", "org_name": "spam", "contact_name": "bot", "email": "bot@example.org", "website": "http://spam"})

	if w := call(t, r, http.MethodGet, "/api/admin/group-applications", "", nil); w.Code != http.StatusUnauthorized {
		t.Errorf("anonymous admin access: got %d", w.Code)
	}
	w = call(t, r, http.MethodGet, "/api/admin/group-applications", admin, nil)
	var apps []struct {
		ID        string
		GroupName *string `json:"group_name"`
		Status    string
	}
	json.Unmarshal(w.Body.Bytes(), &apps)
	if len(apps) != 1 || apps[0].GroupName == nil || *apps[0].GroupName != "台北週三共修" || apps[0].Status != "new" {
		t.Fatalf("group applications = %s", w.Body)
	}
	if w := call(t, r, http.MethodPatch, "/api/admin/group-applications/"+apps[0].ID, admin,
		map[string]string{"status": "contacted", "admin_note": "已電話聯絡"}); w.Code != http.StatusNoContent {
		t.Errorf("review: %d %s", w.Code, w.Body)
	}
	w = call(t, r, http.MethodGet, "/api/admin/merchant-applications", admin, nil)
	if !strings.Contains(w.Body.String(), "山林民宿") || strings.Contains(w.Body.String(), "spam") {
		t.Errorf("merchant applications = %s", w.Body)
	}

	// A knowledge manager cannot see sign-ups.
	if w := call(t, r, http.MethodPost, "/api/auth/users", admin, map[string]string{
		"email": "km@example.org", "display_name": "知識管理員", "password": "km-password-123", "role": "knowledge_manager"}); w.Code != http.StatusCreated {
		t.Fatalf("create user: %d %s", w.Code, w.Body)
	}
	km := login(t, r, "km@example.org", "km-password-123")
	if w := call(t, r, http.MethodGet, "/api/admin/group-applications", km, nil); w.Code != http.StatusForbidden {
		t.Errorf("knowledge manager: got %d, want 403", w.Code)
	}

	// Logout ends the session.
	call(t, r, http.MethodPost, "/api/auth/logout", km, nil)
	if w := call(t, r, http.MethodGet, "/api/auth/me", km, nil); w.Code != http.StatusUnauthorized {
		t.Errorf("after logout: got %d", w.Code)
	}
}
