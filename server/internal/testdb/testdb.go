// Package testdb gives each test package its own throwaway, migrated database.
// Tests that need it are skipped unless TEST_DATABASE_URL points at a server
// where the user may create databases.
package testdb

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/url"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/geodown-ops/bodhi-alliance/server/internal/db"
)

func New(t *testing.T) *pgxpool.Pool {
	t.Helper()
	admin := os.Getenv("TEST_DATABASE_URL")
	if admin == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	ctx := context.Background()
	b := make([]byte, 6)
	rand.Read(b)
	name := "bodhi_t_" + hex.EncodeToString(b)

	conn, err := pgx.Connect(ctx, admin)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(ctx, "CREATE DATABASE "+name); err != nil {
		t.Fatal(err)
	}
	conn.Close(ctx)

	u, err := url.Parse(admin)
	if err != nil {
		t.Fatal(err)
	}
	u.Path = "/" + name
	pool, err := db.Open(ctx, u.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		pool.Close()
		if c, err := pgx.Connect(ctx, admin); err == nil {
			c.Exec(ctx, "DROP DATABASE IF EXISTS "+name+" WITH (FORCE)")
			c.Close(ctx)
		}
	})
	if err := db.Migrate(ctx, pool); err != nil {
		t.Fatal(err)
	}
	// Running migrations twice must be a no-op.
	if err := db.Migrate(ctx, pool); err != nil {
		t.Fatal(err)
	}
	return pool
}
