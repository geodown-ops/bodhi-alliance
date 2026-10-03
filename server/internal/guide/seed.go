package guide

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"errors"
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

// SeedDocuments publishes each seed/*.md file the first time the service sees it, and
// carries later changes to the file into a document nobody has edited since. The setting
// "seed:<file name>" holds the document id and a digest of the text the file last put there,
// so edits, unpublishing or archiving in the admin stick across restarts.
func (s *Store) SeedDocuments(ctx context.Context, fsys fs.FS) (added, updated int, err error) {
	files, err := fs.Glob(fsys, "seed/*.md")
	if err != nil {
		return 0, 0, err
	}
	for _, name := range files {
		raw, err := fs.ReadFile(fsys, name)
		if err != nil {
			return added, updated, err
		}
		title, category, body, err := parseSeed(string(raw))
		if err != nil {
			return added, updated, fmt.Errorf("%s: %w", name, err)
		}
		a, u, err := s.seedOne(ctx, "seed:"+path.Base(name), title, category, name, body)
		if err != nil {
			return added, updated, fmt.Errorf("%s: %w", name, err)
		}
		if a {
			added++
		}
		if u {
			updated++
		}
	}
	return added, updated, nil
}

func (s *Store) seedOne(ctx context.Context, key, title, category, sourceName, body string) (added, updated bool, err error) {
	err = pgx.BeginFunc(ctx, s.DB, func(tx pgx.Tx) error {
		// Claiming the key first makes a second instance starting at the same time wait, then skip.
		tag, err := tx.Exec(ctx, `INSERT INTO guide.setting (key, value) VALUES ($1, '') ON CONFLICT (key) DO NOTHING`, key)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			updated, err = refreshSeed(ctx, tx, key, sourceName, body)
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
		if _, err := tx.Exec(ctx, `UPDATE guide.setting SET value = $2 WHERE key = $1`, key, id+" "+digest(body)); err != nil {
			return err
		}
		added = true
		return rechunk(ctx, tx, id)
	})
	return added && err == nil, updated && err == nil, err
}

// refreshSeed adds the file's new text as a version of an already seeded document,
// unless someone changed the document's text in the admin.
func refreshSeed(ctx context.Context, tx pgx.Tx, key, sourceName, body string) (bool, error) {
	var value string
	if err := tx.QueryRow(ctx, `SELECT value FROM guide.setting WHERE key = $1 FOR UPDATE`, key).Scan(&value); err != nil {
		return false, err
	}
	id, sum, _ := strings.Cut(value, " ")
	var status, current string
	var version int
	var edited bool
	err := tx.QueryRow(ctx, `
		SELECT d.status, d.current_version, v.body, v.edited_by IS NOT NULL FROM guide.document d
		JOIN guide.document_version v ON v.document_id = d.id AND v.version = d.current_version
		WHERE d.id::text = $1 FOR UPDATE OF d`, id).Scan(&status, &version, &current, &edited)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if sum == "" {
		// Seeded before the digest was recorded: only the untouched first version is ours.
		if version != 1 || edited {
			return false, nil
		}
	} else if sum != digest(current) {
		return false, nil
	}
	if current == body {
		_, err := tx.Exec(ctx, `UPDATE guide.setting SET value = $2 WHERE key = $1`, key, id+" "+digest(body))
		return false, err
	}
	version++
	if _, err := tx.Exec(ctx,
		`INSERT INTO guide.document_version (document_id, version, source_name, body) VALUES ($1, $2, $3, $4)`,
		id, version, sourceName, body); err != nil {
		return false, err
	}
	if _, err := tx.Exec(ctx, `UPDATE guide.document SET current_version = $2, updated_at = now() WHERE id = $1`, id, version); err != nil {
		return false, err
	}
	if _, err := tx.Exec(ctx, `UPDATE guide.setting SET value = $2 WHERE key = $1`, key, id+" "+digest(body)); err != nil {
		return false, err
	}
	if status == "published" {
		return true, rechunk(ctx, tx, id)
	}
	return true, nil
}

func digest(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
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
