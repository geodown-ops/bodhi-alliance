package guide

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/jackc/pgx/v5"
)

// SeedFiles is the knowledge the guide ships with, so a fresh deploy can answer
// without anyone uploading files first.
//
//go:embed seed/*.md
var SeedFiles embed.FS

// SeedDocuments publishes each seed/*.md file the first time the service sees it.
// The setting "seed:<file name>" remembers it, so edits, unpublishing or archiving
// in the admin stick across restarts. It returns how many documents it added.
func (s *Store) SeedDocuments(ctx context.Context, fsys fs.FS) (int, error) {
	files, err := fs.Glob(fsys, "seed/*.md")
	if err != nil {
		return 0, err
	}
	added := 0
	for _, name := range files {
		raw, err := fs.ReadFile(fsys, name)
		if err != nil {
			return added, err
		}
		title, category, body, err := parseSeed(string(raw))
		if err != nil {
			return added, fmt.Errorf("%s: %w", name, err)
		}
		ok, err := s.seedOne(ctx, "seed:"+path.Base(name), title, category, name, body)
		if err != nil {
			return added, fmt.Errorf("%s: %w", name, err)
		}
		if ok {
			added++
		}
	}
	return added, nil
}

func (s *Store) seedOne(ctx context.Context, key, title, category, sourceName, body string) (bool, error) {
	added := false
	err := pgx.BeginFunc(ctx, s.DB, func(tx pgx.Tx) error {
		// Claiming the key first makes a second instance starting at the same time wait, then skip.
		tag, err := tx.Exec(ctx, `INSERT INTO guide.setting (key, value) VALUES ($1, '') ON CONFLICT (key) DO NOTHING`, key)
		if err != nil || tag.RowsAffected() == 0 {
			return err
		}
		var id string
		if err := tx.QueryRow(ctx,
			`INSERT INTO guide.document (title, category, status) VALUES ($1, $2, 'published') RETURNING id`,
			title, category).Scan(&id); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO guide.document_version (document_id, version, source_name, body) VALUES ($1, 1, $2, $3)`,
			id, sourceName, body); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE guide.setting SET value = $2 WHERE key = $1`, key, id); err != nil {
			return err
		}
		added = true
		return rechunk(ctx, tx, id)
	})
	return added && err == nil, err
}

// parseSeed reads the "---" header (title, category) in front of a seed file.
func parseSeed(raw string) (title, category, body string, err error) {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	rest, ok := strings.CutPrefix(raw, "---\n")
	if !ok {
		return "", "", "", fmt.Errorf("missing --- header")
	}
	header, body, ok := strings.Cut(rest, "\n---\n")
	if !ok {
		return "", "", "", fmt.Errorf("unterminated --- header")
	}
	for _, line := range strings.Split(header, "\n") {
		k, v, _ := strings.Cut(line, ":")
		switch strings.TrimSpace(k) {
		case "title":
			title = strings.TrimSpace(v)
		case "category":
			category = strings.TrimSpace(v)
		}
	}
	if title == "" {
		return "", "", "", fmt.Errorf("header needs a title")
	}
	if _, ok := categories[category]; !ok {
		return "", "", "", fmt.Errorf("unknown category %q", category)
	}
	return title, category, strings.TrimSpace(body), nil
}
