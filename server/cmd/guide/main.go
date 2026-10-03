// Command guide serves the 線上覺行小組 AI 組長 and its knowledge-base admin.
// It runs apart from the core API so it can be scaled or switched off on its own.
package main

import (
	"context"
	"log"
	"time"

	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
	"github.com/geodown-ops/bodhi-alliance/server/internal/config"
	"github.com/geodown-ops/bodhi-alliance/server/internal/db"
	"github.com/geodown-ops/bodhi-alliance/server/internal/guide"
	"github.com/geodown-ops/bodhi-alliance/server/internal/httpx"
)

func main() {
	cfg := config.Load(":8082")
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
	store := &guide.Store{DB: pool}
	if err := store.EnsurePersona(ctx, cfg.GuideName); err != nil {
		log.Fatalf("seed persona: %v", err)
	}
	if n, err := store.SeedDocuments(ctx, guide.SeedFiles); err != nil {
		log.Fatalf("seed knowledge: %v", err)
	} else if n > 0 {
		log.Printf("published %d bundled knowledge documents", n)
	}

	h := &guide.Handler{Store: store, Auth: &auth.Service{DB: pool}}
	if cfg.AnthropicAPIKey != "" {
		h.LLM = guide.NewClaude(cfg.AnthropicAPIKey, cfg.GuideModel)
	} else {
		log.Print("ANTHROPIC_API_KEY is not set; chat is off, the knowledge-base admin still works")
	}

	r := httpx.NewEngine(cfg.AllowedOrigins, cfg.TrustedProxies)
	r.MaxMultipartMemory = 8 << 20
	h.Routes(r.Group("/guide"))
	log.Printf("guide listening on %s", cfg.Addr)
	log.Fatal(r.Run(cfg.Addr))
}
