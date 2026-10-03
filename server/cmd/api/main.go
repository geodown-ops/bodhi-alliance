// Command api serves the official site's core API: login, 覺行小組 and 共好企業
// sign-ups, 場域管理, 共好企業管理, and volunteer accounts with the 志工名冊.
// Ledger, vouchers and governance join here from M3 on.
package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/geodown-ops/bodhi-alliance/server/internal/apply"
	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
	"github.com/geodown-ops/bodhi-alliance/server/internal/config"
	"github.com/geodown-ops/bodhi-alliance/server/internal/db"
	"github.com/geodown-ops/bodhi-alliance/server/internal/httpx"
	"github.com/geodown-ops/bodhi-alliance/server/internal/members"
	"github.com/geodown-ops/bodhi-alliance/server/internal/org"
)

func main() {
	cfg := config.Load(":8081")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := db.Migrate(ctx, pool); err != nil {
		log.Fatal(err)
	}
	authSvc := &auth.Service{DB: pool}
	if err := authSvc.Bootstrap(ctx, cfg.BootstrapAdminEmail, cfg.BootstrapAdminPassword); err != nil {
		log.Fatalf("bootstrap admin: %v", err)
	}

	r := NewRouter(cfg, authSvc)
	log.Printf("api listening on %s", cfg.Addr)
	log.Fatal(r.Run(cfg.Addr))
}

func NewRouter(cfg config.Config, authSvc *auth.Service) *gin.Engine {
	r := httpx.NewEngine(cfg.AllowedOrigins, cfg.TrustedProxies)
	api := r.Group("/api")
	authSvc.Routes(api.Group("/auth"))
	(&apply.Handler{DB: authSvc.DB, Auth: authSvc}).Routes(api)
	(&org.Handler{DB: authSvc.DB, Auth: authSvc}).Routes(api)
	(&members.Handler{DB: authSvc.DB, Auth: authSvc}).Routes(api)
	return r
}
