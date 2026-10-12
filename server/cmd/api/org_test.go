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

func idOf(t *testing.T, body []byte) string {
	t.Helper()
	var v struct{ ID string }
	if err := json.Unmarshal(body, &v); err != nil || v.ID == "" {
		t.Fatalf("no id in %s", body)
	}
	return v.ID
}

func TestVenuesAndMerchants(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.New(t)
	authSvc := &auth.Service{DB: pool}
	if err := authSvc.Bootstrap(context.Background(), "admin@example.org", "a-long-admin-password"); err != nil {
		t.Fatal(err)
	}
	r := NewRouter(config.Config{}, authSvc, nil)
	admin := login(t, r, "admin@example.org", "a-long-admin-password")

	w := call(t, r, http.MethodPost, "/api/admin/centers", admin, map[string]any{"name": "台北禪修中心", "region": "台北", "status": "active"})
	if w.Code != http.StatusOK {
		t.Fatalf("center: %d %s", w.Code, w.Body)
	}
	center := idOf(t, w.Body.Bytes())
	if w := call(t, r, http.MethodPost, "/api/admin/centers", admin, map[string]any{"name": "台北禪修中心", "status": "active"}); w.Code != http.StatusConflict {
		t.Errorf("duplicate center: got %d", w.Code)
	}

	// An alliance unit needs a center; a sponsor must not have one (the server drops it).
	if w := call(t, r, http.MethodPost, "/api/admin/merchants", admin, map[string]any{"kind": "alliance_unit", "name": "中心齋堂", "status": "active"}); w.Code != http.StatusBadRequest {
		t.Errorf("alliance unit without center: got %d %s", w.Code, w.Body)
	}
	w = call(t, r, http.MethodPost, "/api/admin/merchants", admin, map[string]any{"kind": "alliance_unit", "center_id": center, "name": "中心齋堂", "status": "active", "proposed_monthly_cap": 30000})
	if w.Code != http.StatusOK {
		t.Fatalf("merchant: %d %s", w.Code, w.Body)
	}
	unit := idOf(t, w.Body.Bytes())
	w = call(t, r, http.MethodPost, "/api/admin/merchants", admin, map[string]any{"kind": "sponsor", "center_id": center, "name": "山林飯店", "status": "pending"})
	if w.Code != http.StatusOK {
		t.Fatalf("sponsor: %d %s", w.Code, w.Body)
	}

	w = call(t, r, http.MethodPost, "/api/admin/venues", admin, map[string]any{"center_id": center, "name": "大禪堂", "status": "active", "merchant_id": unit, "charges_public": true})
	if w.Code != http.StatusOK {
		t.Fatalf("venue: %d %s", w.Code, w.Body)
	}
	venue := idOf(t, w.Body.Bytes())
	w = call(t, r, http.MethodGet, "/api/admin/venues?center_id="+center, admin, nil)
	if !strings.Contains(w.Body.String(), `"merchant_name":"中心齋堂"`) || !strings.Contains(w.Body.String(), `"center_name":"台北禪修中心"`) {
		t.Errorf("venues = %s", w.Body)
	}

	// A center with venues cannot be deleted.
	if w := call(t, r, http.MethodDelete, "/api/admin/centers/"+center, admin, nil); w.Code != http.StatusConflict {
		t.Errorf("delete center in use: got %d", w.Code)
	}
	if w := call(t, r, http.MethodDelete, "/api/admin/venues/"+venue, admin, nil); w.Code != http.StatusNoContent {
		t.Errorf("delete venue: got %d", w.Code)
	}
	if w := call(t, r, http.MethodPut, "/api/admin/venues/not-a-uuid", admin, map[string]any{"center_id": center, "name": "x", "status": "active"}); w.Code != http.StatusNotFound {
		t.Errorf("update bad id: got %d", w.Code)
	}

	// A 共好企業 registration becomes a pending merchant once, and the registration is accepted.
	call(t, r, http.MethodPost, "/api/merchant-applications", "", map[string]any{
		"kind": "sponsor", "org_name": "海邊民宿", "contact_name": "林小姐", "email": "sea@example.org", "monthly_scale": "每月 5 間房晚"})
	w = call(t, r, http.MethodGet, "/api/admin/merchant-applications", admin, nil)
	var apps []struct{ ID string }
	json.Unmarshal(w.Body.Bytes(), &apps)
	if len(apps) != 1 {
		t.Fatalf("applications = %s", w.Body)
	}
	if w := call(t, r, http.MethodPost, "/api/admin/merchant-applications/"+apps[0].ID+"/merchant", admin, nil); w.Code != http.StatusCreated {
		t.Fatalf("convert: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodPost, "/api/admin/merchant-applications/"+apps[0].ID+"/merchant", admin, nil); w.Code != http.StatusConflict {
		t.Errorf("convert twice: got %d", w.Code)
	}
	w = call(t, r, http.MethodGet, "/api/admin/merchants?status=pending", admin, nil)
	if !strings.Contains(w.Body.String(), "海邊民宿") || !strings.Contains(w.Body.String(), "每月 5 間房晚") || !strings.Contains(w.Body.String(), "山林飯店") {
		t.Errorf("pending merchants = %s", w.Body)
	}
	if strings.Contains(w.Body.String(), `"center_name":"台北禪修中心"`) {
		t.Error("sponsor should not keep a center")
	}
	w = call(t, r, http.MethodGet, "/api/admin/merchant-applications?status=accepted", admin, nil)
	if !strings.Contains(w.Body.String(), "海邊民宿") {
		t.Errorf("registration not accepted: %s", w.Body)
	}

	// Groups can point at a center.
	w = call(t, r, http.MethodPost, "/api/admin/groups", admin, map[string]any{"name": "週三共修", "region": "台北", "center_id": center, "is_listed": true})
	if w.Code != http.StatusOK {
		t.Fatalf("group: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodGet, "/api/groups", "", nil); !strings.Contains(w.Body.String(), `"center_name":"台北禪修中心"`) {
		t.Errorf("public groups = %s", w.Body)
	}
}
