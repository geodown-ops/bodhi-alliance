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

func TestEventsAndCoinClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.New(t)
	authSvc := &auth.Service{DB: pool}
	if err := authSvc.Bootstrap(context.Background(), "admin@example.org", "a-long-admin-password"); err != nil {
		t.Fatal(err)
	}
	r := NewRouter(config.Config{}, authSvc, nil)
	admin := login(t, r, "admin@example.org", "a-long-admin-password")

	// Signing up needs only a legal name, email and password; nickname and LINE ID are optional.
	signup := func(email, name, line string) (token, volunteerID string) {
		w := call(t, r, http.MethodPost, "/api/volunteers", "", map[string]any{"email": email, "password": "lotus-pond-evening",
			"legal_name": name, "display_name": "", "line_id": line})
		if w.Code != http.StatusCreated {
			t.Fatalf("register %s: %d %s", email, w.Code, w.Body)
		}
		var reg struct{ Token string }
		json.Unmarshal(w.Body.Bytes(), &reg)
		w = call(t, r, http.MethodGet, "/api/me/volunteer", reg.Token, nil)
		var v struct{ ID string }
		json.Unmarshal(w.Body.Bytes(), &v)
		if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"line_id":"`+line+`"`) || !strings.Contains(w.Body.String(), `"center_name":""`) {
			t.Fatalf("me: %d %s", w.Code, w.Body)
		}
		return reg.Token, v.ID
	}
	if w := call(t, r, http.MethodPost, "/api/volunteers", "", map[string]any{"email": "x@example.org", "password": "lotus-pond-evening"}); w.Code != http.StatusBadRequest {
		t.Errorf("signup without legal name: %d", w.Code)
	}
	org, _ := signup("org@example.org", "陳發起", "chen-line")
	helper, helperID := signup("help@example.org", "王協辦", "")
	p1, _ := signup("p1@example.org", "李參加", "")
	p2, _ := signup("p2@example.org", "張參加", "")

	start := time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second)
	event := map[string]any{"title": "週六湖邊正念減壓", "is_online": false, "location": "大安森林公園", "capacity": 3,
		"starts_at": start, "ends_at": start.Add(2 * time.Hour)}
	if w := call(t, r, http.MethodPost, "/api/me/events", org, map[string]any{"title": "兩個人", "location": "x", "capacity": 2,
		"starts_at": start, "ends_at": start.Add(time.Hour)}); w.Code != http.StatusBadRequest {
		t.Errorf("capacity 2: %d", w.Code)
	}
	if w := call(t, r, http.MethodPost, "/api/me/events", admin, event); w.Code != http.StatusNotFound {
		t.Errorf("admin without volunteer record creates event: %d", w.Code)
	}
	w := call(t, r, http.MethodPost, "/api/me/events", org, event)
	if w.Code != http.StatusCreated {
		t.Fatalf("create event: %d %s", w.Code, w.Body)
	}
	id := idOf(t, w.Body.Bytes())

	if w := call(t, r, http.MethodGet, "/api/events", "", nil); !strings.Contains(w.Body.String(), `"joined":1`) || !strings.Contains(w.Body.String(), "週六湖邊正念減壓") {
		t.Errorf("public events = %s", w.Body)
	}
	if w := call(t, r, http.MethodPost, "/api/me/events/"+id+"/join", helper, map[string]string{"role": "helper"}); w.Code != http.StatusNoContent {
		t.Fatalf("helper join: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodPost, "/api/me/events/"+id+"/join", p1, nil); w.Code != http.StatusNoContent {
		t.Fatalf("p1 join: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodPost, "/api/me/events/"+id+"/join", p2, nil); w.Code != http.StatusConflict {
		t.Errorf("join a full event: %d", w.Code)
	}
	if w := call(t, r, http.MethodPost, "/api/me/events/"+id+"/claim", org, map[string]any{"attendance": 3}); w.Code != http.StatusConflict {
		t.Errorf("claim before the event ends: %d %s", w.Code, w.Body)
	}

	// The event takes place.
	if _, err := pool.Exec(context.Background(), `UPDATE practice_event SET starts_at = now() - interval '3 hours', ends_at = now() - interval '1 hour' WHERE id = $1`, id); err != nil {
		t.Fatal(err)
	}
	if w := call(t, r, http.MethodPost, "/api/me/events/"+id+"/claim", helper, map[string]any{"attendance": 3}); w.Code != http.StatusForbidden {
		t.Errorf("helper submits: %d", w.Code)
	}
	if w := call(t, r, http.MethodPost, "/api/me/events/"+id+"/claim", org, map[string]any{"attendance": 2}); w.Code != http.StatusBadRequest {
		t.Errorf("two people attended: %d", w.Code)
	}
	if w := call(t, r, http.MethodPost, "/api/me/events/"+id+"/claim", org, map[string]any{"attendance": 3, "report": "三人到齊"}); w.Code != http.StatusNoContent {
		t.Fatalf("claim: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodPost, "/api/me/events/"+id+"/claim", org, map[string]any{"attendance": 3}); w.Code != http.StatusConflict {
		t.Errorf("claim twice: %d", w.Code)
	}
	if w := call(t, r, http.MethodDelete, "/api/me/events/"+id+"/join", helper, nil); w.Code != http.StatusConflict {
		t.Errorf("helper leaves after submission: %d", w.Code)
	}

	// Only 超級管理員 review; organizer and helper get the suggested 2 h × 300.
	if w := call(t, r, http.MethodGet, "/api/admin/claims", org, nil); w.Code != http.StatusForbidden {
		t.Errorf("volunteer lists claims: %d", w.Code)
	}
	w = call(t, r, http.MethodGet, "/api/admin/claims?status=submitted", admin, nil)
	var claims []struct {
		ID         string
		Recipients []struct {
			VolunteerID string `json:"volunteer_id"`
			Role        string
			Amount      int64
		}
	}
	json.Unmarshal(w.Body.Bytes(), &claims)
	if len(claims) != 1 || len(claims[0].Recipients) != 2 || claims[0].Recipients[0].Amount != 600 {
		t.Fatalf("claims = %s", w.Body)
	}
	claim := claims[0].ID
	if w := call(t, r, http.MethodPatch, "/api/admin/claims/"+claim, admin, map[string]any{"status": "rejected"}); w.Code != http.StatusBadRequest {
		t.Errorf("reject without reason: %d", w.Code)
	}
	if w := call(t, r, http.MethodPatch, "/api/admin/claims/"+claim, admin, map[string]any{"status": "rejected", "review_note": "請補出席名單"}); w.Code != http.StatusNoContent {
		t.Fatalf("reject: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodGet, "/api/me/events", org, nil); !strings.Contains(w.Body.String(), `"claim_status":"rejected"`) || !strings.Contains(w.Body.String(), "請補出席名單") {
		t.Errorf("my events after reject = %s", w.Body)
	}
	if w := call(t, r, http.MethodPost, "/api/me/events/"+id+"/claim", org, map[string]any{"attendance": 3, "report": "李參加、王協辦、陳發起"}); w.Code != http.StatusNoContent {
		t.Fatalf("resubmit: %d %s", w.Code, w.Body)
	}
	approve := map[string]any{"status": "approved", "amounts": map[string]int64{helperID: 900}}
	if w := call(t, r, http.MethodPatch, "/api/admin/claims/"+claim, admin, approve); w.Code != http.StatusNoContent {
		t.Fatalf("approve: %d %s", w.Code, w.Body)
	}
	if w := call(t, r, http.MethodPatch, "/api/admin/claims/"+claim, admin, approve); w.Code != http.StatusConflict {
		t.Errorf("approve twice: %d", w.Code)
	}
	balance := func(token string) int64 {
		w := call(t, r, http.MethodGet, "/api/me/wallet", token, nil)
		var res struct{ Balance int64 }
		json.Unmarshal(w.Body.Bytes(), &res)
		return res.Balance
	}
	if b := balance(org); b != 600 {
		t.Errorf("organizer balance = %d", b)
	}
	if b := balance(helper); b != 900 {
		t.Errorf("helper balance = %d", b)
	}
	if b := balance(p1); b != 0 {
		t.Errorf("participant balance = %d", b)
	}
}
