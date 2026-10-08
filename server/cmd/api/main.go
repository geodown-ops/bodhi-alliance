// Command api serves the official site's core API: login, 覺行小組 and 共好企業
// sign-ups, 場域管理, 共好企業管理, and volunteer accounts with the 志工名冊.
// It also serves 覺行共修活動 and the 菩提幣 wallet ledger they credit, and the 協會會刊,
// 會員行事曆 and 協會通知 for members who join 世界佛教教育協會.
// With BODHI_CHAIN_* set it also runs the 菩提幣 chain worker (每位會員入會得 1 枚，上鏈可查).
package main

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/geodown-ops/bodhi-alliance/server/internal/apply"
	"github.com/geodown-ops/bodhi-alliance/server/internal/association"
	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
	"github.com/geodown-ops/bodhi-alliance/server/internal/chain"
	"github.com/geodown-ops/bodhi-alliance/server/internal/config"
	"github.com/geodown-ops/bodhi-alliance/server/internal/db"
	"github.com/geodown-ops/bodhi-alliance/server/internal/events"
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

	chainSvc := startChain(pool)

	r := NewRouter(cfg, authSvc, chainSvc)
	log.Printf("api listening on %s", cfg.Addr)
	log.Fatal(r.Run(cfg.Addr))
}

// startChain connects to the chain in the background when its keys are configured,
// retrying every minute across the listed nodes, then runs the 菩提幣 worker.
func startChain(pool *pgxpool.Pool) *chain.Holder {
	cc := chain.ConfigFromEnv()
	h := &chain.Holder{}
	if !cc.Enabled() {
		log.Printf("chain: off (set BODHI_CHAIN_OPERATOR_KEY and BODHI_CHAIN_MEMBER_SEED to turn it on)")
		return h
	}
	go func() {
		for {
			for _, url := range cc.RPCs() {
				svc, err := connectChain(pool, cc, url)
				if err != nil {
					log.Printf("chain: %s: %v", url, err)
					continue
				}
				log.Printf("chain: connected to %s", url)
				h.Set(svc)
				svc.Run(context.Background())
				return
			}
			time.Sleep(time.Minute)
		}
	}()
	return h
}

func connectChain(pool *pgxpool.Pool, cc chain.Config, url string) (*chain.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := ethclient.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	svc, err := chain.New(ctx, pool, client, cc)
	if err != nil {
		client.Close()
		return nil, err
	}
	return svc, nil
}

func NewRouter(cfg config.Config, authSvc *auth.Service, chainSvc *chain.Holder) *gin.Engine {
	r := httpx.NewEngine(cfg.AllowedOrigins, cfg.TrustedProxies)
	api := r.Group("/api")
	authSvc.Routes(api.Group("/auth"))
	(&apply.Handler{DB: authSvc.DB, Auth: authSvc}).Routes(api)
	(&org.Handler{DB: authSvc.DB, Auth: authSvc}).Routes(api)
	(&members.Handler{DB: authSvc.DB, Auth: authSvc}).Routes(api)
	(&events.Handler{DB: authSvc.DB, Auth: authSvc}).Routes(api)
	(&association.Handler{DB: authSvc.DB, Auth: authSvc}).Routes(api)
	(&chain.Handler{DB: authSvc.DB, Auth: authSvc, Svc: chainSvc}).Routes(api)
	return r
}
