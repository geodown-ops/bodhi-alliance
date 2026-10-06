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

func TestVenueEvents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testdb.New(t)
	authSvc := &auth.Service{DB: pool}
	if err := authSvc.Bootstrap(context.Background(), "admin@example.org", "a-long-admin-password"); err != nil {
		t.Fatal(err)
	}
	r := NewRouter(config.Config{}, authSvc)
	admin := login(t, r, "admin@example.org", "a-long-admin-password")

	w := call(t, r, http.MethodPost, "/api/admin/centers", admin, map[string]any{"name": "台中禪修中心", "region": "台中", "status": "active"})
	center := idOf(t, w.Body.Bytes())
	w = call(t, r, http.MethodPost, "/api/admin/venues", admin, map[string]any{"center_id": center, "name": "林間禪堂", "address": "台中市北屯區山邊路 1 號", "status": "active"})
	venue := idOf(t, w.Body.Bytes())
	w = call(t, r, http.MethodPost, "/api/admin/venues", admin, map[string]any{"center_id": center, "name": "舊教室", "status": "closed"})
	closed := idOf(t, w.Body.Bytes())

	// The public list shows only active venues.
	w = call(t, r, http.MethodGet, "/api/venues", "", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "林間禪堂") || strings.Contains(w.Body.String(), "舊教室") ||
		!strings.Contains(w.Body.String(), `"upcoming":0`) {
		t.Fatalf("venues = %d %s", w.Code, w.Body)
	}

	w = call(t, r, http.MethodPost, "/api/volunteers", "", map[string]any{"email": "org@example.org", "password": "lotus-pond-evening", "legal_name": "王小明"})
	var reg struct{ Token string }
	json.Unmarshal(w.Body.Bytes(), &reg)

	start := time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second)
	event := func(venueID string, online bool, location string) map[string]any {
		return map[string]any{"title": "週六晨間正念", "is_online": online, "location": location, "capacity": 6,
			"starts_at": start, "ends_at": start.Add(2 * time.Hour), "venue_id": venueID}
	}
	if w := call(t, r, http.MethodPost, "/api/me/events", reg.Token, event(closed, false, "")); w.Code != http.StatusBadRequest {
		t.Errorf("closed venue: %d %s", w.Code, w.Body)
	}
	// With a venue and no address typed, the venue's address is used.
	if w := call(t, r, http.MethodPost, "/api/me/events", reg.Token, event(venue, false, "")); w.Code != http.StatusCreated {
		t.Fatalf("venue event: %d %s", w.Code, w.Body)
	}
	// An online activity never belongs to a venue.
	if w := call(t, r, http.MethodPost, "/api/me/events", reg.Token, event(venue, true, "https://meet.example.org/x")); w.Code != http.StatusCreated {
		t.Fatalf("online event: %d %s", w.Code, w.Body)
	}

	w = call(t, r, http.MethodGet, "/api/events", "", nil)
	var events []struct {
		IsOnline  bool    `json:"is_online"`
		Location  string  `json:"location"`
		VenueID   *string `json:"venue_id"`
		VenueName string  `json:"venue_name"`
	}
	json.Unmarshal(w.Body.Bytes(), &events)
	if len(events) != 2 {
		t.Fatalf("events = %s", w.Body)
	}
	for _, e := range events {
		switch {
		case e.IsOnline && e.VenueID != nil:
			t.Errorf("online event has a venue: %+v", e)
		case !e.IsOnline && (e.VenueID == nil || *e.VenueID != venue || e.VenueName != "林間禪堂" || e.Location != "台中市北屯區山邊路 1 號"):
			t.Errorf("venue event = %+v", e)
		}
	}
	if w := call(t, r, http.MethodGet, "/api/venues", "", nil); !strings.Contains(w.Body.String(), `"upcoming":1`) {
		t.Errorf("venues after event = %s", w.Body)
	}
}
