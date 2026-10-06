package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
	"github.com/geodown-ops/bodhi-alliance/server/internal/config"
	"github.com/geodown-ops/bodhi-alliance/server/internal/testdb"
)

func TestAssociationMembership(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.New(t)
	authSvc := &auth.Service{DB: pool}
	if err := authSvc.Bootstrap(context.Background(), "admin@example.org", "a-long-admin-password"); err != nil {
		t.Fatal(err)
	}
	r := NewRouter(config.Config{}, authSvc)
	admin := login(t, r, "admin@example.org", "a-long-admin-password")

	signup := func(body map[string]any) string {
		body["password"] = "lotus-pond-evening"
		body["legal_name"] = "陳大華"
		w := call(t, r, http.MethodPost, "/api/volunteers", "", body)
		if w.Code != http.StatusCreated {
			t.Fatalf("register %v: %d %s", body, w.Code, w.Body)
		}
		var reg struct{ Token string }
		json.Unmarshal(w.Body.Bytes(), &reg)
		return reg.Token
	}
	if w := call(t, r, http.MethodPost, "/api/volunteers", "", map[string]any{"email": "none@example.org", "password": "lotus-pond-evening",
		"legal_name": "x", "in_groups": false, "in_association": false}); w.Code != http.StatusBadRequest {
		t.Errorf("joining nothing: %d", w.Code)
	}

	// Joining only the association: no wallet, but the association page works.
	assoc := signup(map[string]any{"email": "assoc@example.org", "in_groups": false, "in_association": true})
	w := call(t, r, http.MethodGet, "/api/me/volunteer", assoc, nil)
	if !strings.Contains(w.Body.String(), `"in_groups":false`) || !strings.Contains(w.Body.String(), `"in_association":true`) {
		t.Fatalf("profile = %s", w.Body)
	}
	if w := call(t, r, http.MethodGet, "/api/me/wallet", assoc, nil); w.Code != http.StatusForbidden {
		t.Errorf("wallet without 覺行小組: %d", w.Code)
	}
	start := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second)
	if w := call(t, r, http.MethodPost, "/api/me/events", assoc, map[string]any{"title": "x", "location": "y", "capacity": 3,
		"starts_at": start, "ends_at": start.Add(time.Hour)}); w.Code != http.StatusForbidden {
		t.Errorf("event without 覺行小組: %d", w.Code)
	}

	// The old sign-up (no choice given) joins 覺行小組 only.
	groups := signup(map[string]any{"email": "groups@example.org"})
	if w := call(t, r, http.MethodGet, "/api/me/association", groups, nil); w.Code != http.StatusForbidden {
		t.Errorf("association page without joining: %d", w.Code)
	}
	if w := call(t, r, http.MethodGet, "/api/me/wallet", groups, nil); w.Code != http.StatusOK {
		t.Errorf("wallet: %d", w.Code)
	}

	// Admins publish; members see it.
	for path, body := range map[string]map[string]any{
		"issues":  {"title": "協會會刊第一期", "issued_on": "2026-10-01", "url": "https://example.org/issue1.pdf"},
		"events":  {"title": "秋季會員大會", "starts_at": start, "location": "台北"},
		"notices": {"title": "會員大會報名開始", "body": "請在十月底前回覆"},
	} {
		if w := call(t, r, http.MethodPost, "/api/admin/association/"+path, admin, body); w.Code != http.StatusOK {
			t.Fatalf("publish %s: %d %s", path, w.Code, w.Body)
		}
	}
	if w := call(t, r, http.MethodPost, "/api/admin/association/issues", admin, map[string]any{"title": "x", "issued_on": "10/1"}); w.Code != http.StatusBadRequest {
		t.Errorf("bad date: %d", w.Code)
	}
	if w := call(t, r, http.MethodPost, "/api/admin/association/issues", assoc, map[string]any{"title": "x", "issued_on": "2026-10-01"}); w.Code != http.StatusForbidden && w.Code != http.StatusUnauthorized {
		t.Errorf("member publishing: %d", w.Code)
	}
	w = call(t, r, http.MethodGet, "/api/me/association", assoc, nil)
	for _, want := range []string{"協會會刊第一期", `"issued_on":"2026-10-01"`, "秋季會員大會", "會員大會報名開始"} {
		if !strings.Contains(w.Body.String(), want) {
			t.Errorf("association page missing %s: %s", want, w.Body)
		}
	}

	// Switching: the 覺行小組 member joins the association too; nobody can leave both.
	if w := call(t, r, http.MethodPut, "/api/me/memberships", groups, map[string]any{"in_groups": false, "in_association": false}); w.Code != http.StatusBadRequest {
		t.Errorf("leave both: %d", w.Code)
	}
	w = call(t, r, http.MethodPut, "/api/me/memberships", groups, map[string]any{"in_groups": true, "in_association": true})
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"in_association":true`) {
		t.Fatalf("join association: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodGet, "/api/me/association", groups, nil); w.Code != http.StatusOK {
		t.Errorf("association page after joining: %d", w.Code)
	}
	w = call(t, r, http.MethodGet, "/api/admin/association/members", admin, nil)
	if !strings.Contains(w.Body.String(), "assoc@example.org") || !strings.Contains(w.Body.String(), "groups@example.org") {
		t.Errorf("members = %s", w.Body)
	}
}
