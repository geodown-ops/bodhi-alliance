package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
	"github.com/geodown-ops/bodhi-alliance/server/internal/config"
	"github.com/geodown-ops/bodhi-alliance/server/internal/testdb"
)

func TestVolunteersAndCenterAdmins(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.New(t)
	authSvc := &auth.Service{DB: pool}
	if err := authSvc.Bootstrap(context.Background(), "admin@example.org", "a-long-admin-password"); err != nil {
		t.Fatal(err)
	}
	r := NewRouter(config.Config{}, authSvc, nil)
	admin := login(t, r, "admin@example.org", "a-long-admin-password")

	newCenter := func(name string) string {
		w := call(t, r, http.MethodPost, "/api/admin/centers", admin, map[string]any{"name": name, "status": "active"})
		if w.Code != http.StatusOK {
			t.Fatalf("center: %d %s", w.Code, w.Body)
		}
		return idOf(t, w.Body.Bytes())
	}
	taipei, taichung := newCenter("台北禪修中心"), newCenter("台中禪修中心")
	w := call(t, r, http.MethodPost, "/api/admin/groups", admin, map[string]any{"name": "週三晚間小組", "region": "台北", "center_id": taipei, "is_listed": true})
	if w.Code != http.StatusOK {
		t.Fatalf("group: %d %s", w.Code, w.Body)
	}
	group := idOf(t, w.Body.Bytes())

	if w := call(t, r, http.MethodGet, "/api/centers", "", nil); !strings.Contains(w.Body.String(), "台中禪修中心") {
		t.Errorf("public centers = %s", w.Body)
	}

	// A volunteer signs up, is signed in at once, and starts out pending.
	signup := map[string]any{"email": "Mei@Example.org", "password": "lotus-pond-evening", "legal_name": "林美玲",
		"phone": "0912345678", "home_center_id": taipei, "wants_coach": true}
	w = call(t, r, http.MethodPost, "/api/volunteers", "", signup)
	if w.Code != http.StatusCreated {
		t.Fatalf("register: %d %s", w.Code, w.Body)
	}
	var reg struct{ Token string }
	json.Unmarshal(w.Body.Bytes(), &reg)
	if w := call(t, r, http.MethodPost, "/api/volunteers", "", signup); w.Code != http.StatusConflict {
		t.Errorf("duplicate email: got %d", w.Code)
	}
	vol := reg.Token
	w = call(t, r, http.MethodGet, "/api/me/volunteer", vol, nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"status":"pending"`) || !strings.Contains(w.Body.String(), `"display_name":"林美玲"`) {
		t.Fatalf("me: %d %s", w.Code, w.Body)
	}
	var me struct{ ID string }
	json.Unmarshal(w.Body.Bytes(), &me)

	// Volunteers have no back-office access.
	if w := call(t, r, http.MethodGet, "/api/admin/volunteers", vol, nil); w.Code != http.StatusForbidden {
		t.Errorf("volunteer on admin: got %d", w.Code)
	}

	if w := call(t, r, http.MethodPost, "/api/me/groups/"+group, vol, nil); w.Code != http.StatusNoContent {
		t.Fatalf("join: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodPost, "/api/me/groups/"+group, vol, nil); w.Code != http.StatusNoContent {
		t.Errorf("join twice: %d", w.Code)
	}
	if w := call(t, r, http.MethodPost, "/api/me/groups/not-a-group", vol, nil); w.Code != http.StatusNotFound {
		t.Errorf("join unknown: %d", w.Code)
	}

	// Each center admin sees and verifies only their own center's volunteers.
	newAdmin := func(email, center string) string {
		w := call(t, r, http.MethodPost, "/api/auth/users", admin, map[string]any{"email": email, "display_name": email,
			"password": "center-admin-password", "role": "center_admin", "center_id": center})
		if w.Code != http.StatusCreated {
			t.Fatalf("center admin: %d %s", w.Code, w.Body)
		}
		return login(t, r, email, "center-admin-password")
	}
	tp, tc := newAdmin("tp@example.org", taipei), newAdmin("tc@example.org", taichung)
	if w := call(t, r, http.MethodGet, "/api/admin/volunteers", tc, nil); w.Body.String() != "[]" {
		t.Errorf("other center sees %s", w.Body)
	}
	review := map[string]any{"status": "verified", "is_coach": true}
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, tc, review); w.Code != http.StatusForbidden {
		t.Errorf("other center verifies: got %d", w.Code)
	}
	if w := call(t, r, http.MethodGet, "/api/admin/groups/"+group+"/members", tc, nil); w.Code != http.StatusForbidden {
		t.Errorf("other center lists members: got %d", w.Code)
	}
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, tp, map[string]any{"status": "rejected"}); w.Code != http.StatusBadRequest {
		t.Errorf("reject without reason: got %d", w.Code)
	}
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, tp, review); w.Code != http.StatusNoContent {
		t.Fatalf("verify: %d %s", w.Code, w.Body)
	}
	w = call(t, r, http.MethodGet, "/api/admin/groups/"+group+"/members", tp, nil)
	if !strings.Contains(w.Body.String(), `"legal_name":"林美玲"`) || !strings.Contains(w.Body.String(), `"is_coach":true`) {
		t.Errorf("members = %s", w.Body)
	}
	// 覺行小組長：打勾設定，之後不帶這個欄位的儲存不會把它清掉
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, tp, map[string]any{"status": "verified", "is_coach": true, "is_group_leader": true, "is_committee": true}); w.Code != http.StatusNoContent {
		t.Fatalf("set group leader: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, tp, review); w.Code != http.StatusNoContent {
		t.Fatalf("save without group leader: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodGet, "/api/admin/volunteers", tp, nil); !strings.Contains(w.Body.String(), `"is_group_leader":true`) || !strings.Contains(w.Body.String(), `"is_committee":true`) {
		t.Errorf("group leader lost: %s", w.Body)
	}
	if w := call(t, r, http.MethodPut, "/api/admin/groups/"+group+"/members/"+me.ID, tp, map[string]string{"role": "leader"}); w.Code != http.StatusNoContent {
		t.Errorf("make leader: %d %s", w.Code, w.Body)
	}

	// After verification the legal name is locked; other fields still change.
	w = call(t, r, http.MethodPut, "/api/me/volunteer", vol, map[string]any{"display_name": "美玲", "legal_name": "別人", "phone": "0987654321"})
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"legal_name":"林美玲"`) ||
		!strings.Contains(w.Body.String(), `"display_name":"美玲"`) || !strings.Contains(w.Body.String(), `"role":"leader"`) {
		t.Errorf("update me: %d %s", w.Code, w.Body)
	}

	// An admin account is not a volunteer.
	if w := call(t, r, http.MethodGet, "/api/me/volunteer", admin, nil); w.Code != http.StatusNotFound {
		t.Errorf("admin as volunteer: got %d", w.Code)
	}
	if w := call(t, r, http.MethodPost, "/api/auth/users", admin, map[string]any{"email": "x@example.org", "display_name": "x",
		"password": "center-admin-password", "role": "center_admin"}); w.Code != http.StatusBadRequest {
		t.Errorf("center admin without center: got %d", w.Code)
	}

	// 管理員：只有超級管理員能打勾；打勾後會員用自己的帳密進後台，但不能管帳號、審菩提幣或指派管理員
	makeAdmin := map[string]any{"status": "verified", "is_coach": true, "is_admin": true}
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, tp, makeAdmin); w.Code != http.StatusForbidden {
		t.Errorf("center admin grants manager: got %d", w.Code)
	}
	if w := call(t, r, http.MethodGet, "/api/admin/centers", vol, nil); w.Code != http.StatusForbidden {
		t.Errorf("member before manager: got %d", w.Code)
	}
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, admin, makeAdmin); w.Code != http.StatusNoContent {
		t.Fatalf("grant manager: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodGet, "/api/admin/volunteers", admin, nil); !strings.Contains(w.Body.String(), `"is_admin":true`) {
		t.Errorf("is_admin missing: %s", w.Body)
	}
	mgr := login(t, r, "mei@example.org", "lotus-pond-evening")
	for _, path := range []string{"/api/admin/centers", "/api/admin/groups", "/api/admin/merchant-applications", "/api/admin/association/notices"} {
		if w := call(t, r, http.MethodGet, path, mgr, nil); w.Code != http.StatusOK {
			t.Errorf("manager GET %s: got %d", path, w.Code)
		}
	}
	if w := call(t, r, http.MethodGet, "/api/admin/volunteers", mgr, nil); !strings.Contains(w.Body.String(), `"email":"mei@example.org"`) {
		t.Errorf("manager roster = %s", w.Body)
	}
	for _, path := range []string{"/api/auth/users", "/api/admin/claims"} {
		if w := call(t, r, http.MethodGet, path, mgr, nil); w.Code != http.StatusForbidden {
			t.Errorf("manager GET %s: got %d", path, w.Code)
		}
	}
	// 管理員自己不能取消或指派管理員；照原值儲存其他欄位可以
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, mgr, map[string]any{"status": "verified", "is_admin": false}); w.Code != http.StatusForbidden {
		t.Errorf("manager revokes manager: got %d", w.Code)
	}
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, mgr, makeAdmin); w.Code != http.StatusNoContent {
		t.Errorf("manager saves unchanged: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, admin, map[string]any{"status": "verified", "is_admin": false}); w.Code != http.StatusNoContent {
		t.Fatalf("revoke manager: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodGet, "/api/admin/centers", mgr, nil); w.Code != http.StatusForbidden {
		t.Errorf("after revoke: got %d", w.Code)
	}

	// 超級管理員：一樣只有超級管理員能打勾，打勾後全部功能都能用
	makeSuper := map[string]any{"status": "verified", "is_coach": true, "is_super_admin": true}
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, tp, makeSuper); w.Code != http.StatusForbidden {
		t.Errorf("center admin grants super: got %d", w.Code)
	}
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, admin, makeSuper); w.Code != http.StatusNoContent {
		t.Fatalf("grant super: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodGet, "/api/admin/volunteers", admin, nil); !strings.Contains(w.Body.String(), `"is_super_admin":true`) {
		t.Errorf("is_super_admin missing: %s", w.Body)
	}
	for _, path := range []string{"/api/auth/users", "/api/admin/claims"} {
		if w := call(t, r, http.MethodGet, path, mgr, nil); w.Code != http.StatusOK {
			t.Errorf("super member GET %s: got %d", path, w.Code)
		}
	}
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, mgr, map[string]any{"status": "verified", "is_super_admin": false}); w.Code != http.StatusBadRequest {
		t.Errorf("revoke own super: got %d", w.Code)
	}
	if w := call(t, r, http.MethodPatch, "/api/admin/volunteers/"+me.ID, admin, map[string]any{"status": "verified", "is_super_admin": false}); w.Code != http.StatusNoContent {
		t.Fatalf("revoke super: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodGet, "/api/auth/users", mgr, nil); w.Code != http.StatusForbidden {
		t.Errorf("after super revoke: got %d", w.Code)
	}
}
