package guide

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errNotFound = errors.New("找不到這份文件")

var categories = map[string]string{
	"practice":    "覺行知識",
	"script":      "共修腳本",
	"coin":        "菩提幣規則",
	"association": "協會介紹",
	"faq":         "常見問題",
}

type Document struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Category       string    `json:"category"`
	Status         string    `json:"status"`
	CurrentVersion int       `json:"current_version"`
	SourceName     string    `json:"source_name"`
	Chars          int       `json:"chars"`
	UploadedBy     *string   `json:"uploaded_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	Body           string    `json:"body,omitempty"`
}

type Version struct {
	Version    int       `json:"version"`
	SourceName string    `json:"source_name"`
	Chars      int       `json:"chars"`
	EditedBy   *string   `json:"edited_by"`
	CreatedAt  time.Time `json:"created_at"`
}

type Persona struct {
	Version   int       `json:"version"`
	Name      string    `json:"name"`
	Prompt    string    `json:"prompt"`
	EditedBy  *string   `json:"edited_by"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	DB *pgxpool.Pool
}

func (s *Store) CreateDocument(ctx context.Context, title, category, sourceName, body, userID string) (string, error) {
	var id string
	err := pgx.BeginFunc(ctx, s.DB, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx,
			`INSERT INTO guide.document (title, category, uploaded_by) VALUES ($1, $2, $3) RETURNING id`,
			title, category, userID).Scan(&id); err != nil {
			return err
		}
		_, err := tx.Exec(ctx,
			`INSERT INTO guide.document_version (document_id, version, source_name, body, edited_by) VALUES ($1, 1, $2, $3, $4)`,
			id, sourceName, body, userID)
		return err
	})
	return id, err
}

const docColumns = `
	d.id, d.title, d.category, d.status, d.current_version, v.source_name, char_length(v.body),
	u.display_name, d.updated_at`

func scanDoc(row pgx.Row, withBody bool) (Document, error) {
	var d Document
	dest := []any{&d.ID, &d.Title, &d.Category, &d.Status, &d.CurrentVersion, &d.SourceName, &d.Chars, &d.UploadedBy, &d.UpdatedAt}
	if withBody {
		dest = append(dest, &d.Body)
	}
	err := row.Scan(dest...)
	if errors.Is(err, pgx.ErrNoRows) {
		return d, errNotFound
	}
	return d, err
}

func (s *Store) ListDocuments(ctx context.Context) ([]Document, error) {
	rows, err := s.DB.Query(ctx, `SELECT `+docColumns+`
		FROM guide.document d
		JOIN guide.document_version v ON v.document_id = d.id AND v.version = d.current_version
		LEFT JOIN app_user u ON u.id = d.uploaded_by
		WHERE d.status <> 'archived'
		ORDER BY d.category, d.title`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (Document, error) { return scanDoc(r, false) })
}

func (s *Store) GetDocument(ctx context.Context, id string) (Document, error) {
	return scanDoc(s.DB.QueryRow(ctx, `SELECT `+docColumns+`, v.body
		FROM guide.document d
		JOIN guide.document_version v ON v.document_id = d.id AND v.version = d.current_version
		LEFT JOIN app_user u ON u.id = d.uploaded_by
		WHERE d.id::text = $1`, id), true)
}

func (s *Store) Versions(ctx context.Context, id string) ([]Version, error) {
	rows, err := s.DB.Query(ctx, `
		SELECT v.version, v.source_name, char_length(v.body), u.display_name, v.created_at
		FROM guide.document_version v LEFT JOIN app_user u ON u.id = v.edited_by
		WHERE v.document_id::text = $1 ORDER BY v.version DESC`, id)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Version])
}

// UpdateDocument changes title/category, and adds a new version when body changes.
// A published document is re-chunked in the same transaction, so the change is live at once.
func (s *Store) UpdateDocument(ctx context.Context, id, title, category, body, sourceName, userID string) error {
	return pgx.BeginFunc(ctx, s.DB, func(tx pgx.Tx) error {
		var status, current string
		var version int
		err := tx.QueryRow(ctx, `
			SELECT d.status, d.current_version, v.body FROM guide.document d
			JOIN guide.document_version v ON v.document_id = d.id AND v.version = d.current_version
			WHERE d.id::text = $1 FOR UPDATE OF d`, id).Scan(&status, &version, &current)
		if errors.Is(err, pgx.ErrNoRows) {
			return errNotFound
		}
		if err != nil {
			return err
		}
		if body != "" && body != current {
			version++
			if _, err := tx.Exec(ctx, `
				INSERT INTO guide.document_version (document_id, version, source_name, body, edited_by)
				VALUES ($1, $2, $3, $4, $5)`, id, version, sourceName, body, userID); err != nil {
				return err
			}
		}
		if _, err := tx.Exec(ctx, `
			UPDATE guide.document SET title = $2, category = $3, current_version = $4, updated_at = now()
			WHERE id = $1`, id, title, category, version); err != nil {
			return err
		}
		if status == "published" {
			return rechunk(ctx, tx, id)
		}
		return nil
	})
}

// SetStatus publishes (chunks the current version), unpublishes or archives a document.
func (s *Store) SetStatus(ctx context.Context, id, status string) error {
	return pgx.BeginFunc(ctx, s.DB, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `UPDATE guide.document SET status = $2, updated_at = now() WHERE id::text = $1`, id, status)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return errNotFound
		}
		if status == "published" {
			return rechunk(ctx, tx, id)
		}
		_, err = tx.Exec(ctx, `DELETE FROM guide.chunk WHERE document_id::text = $1`, id)
		return err
	})
}

func rechunk(ctx context.Context, tx pgx.Tx, id string) error {
	var body string
	if err := tx.QueryRow(ctx, `
		SELECT v.body FROM guide.document d
		JOIN guide.document_version v ON v.document_id = d.id AND v.version = d.current_version
		WHERE d.id::text = $1`, id).Scan(&body); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM guide.chunk WHERE document_id::text = $1`, id); err != nil {
		return err
	}
	for _, c := range SplitChunks(body) {
		if _, err := tx.Exec(ctx, `INSERT INTO guide.chunk (document_id, seq, heading, body) VALUES ($1, $2, $3, $4)`,
			id, c.Seq, c.Heading, c.Body); err != nil {
			return err
		}
	}
	return nil
}

// PublishedChunks returns every live chunk in a stable order (stable order keeps the prompt cacheable).
func (s *Store) PublishedChunks(ctx context.Context) ([]Chunk, error) {
	rows, err := s.DB.Query(ctx, `
		SELECT c.document_id, d.title, d.category, c.seq, c.heading, c.body
		FROM guide.chunk c JOIN guide.document d ON d.id = c.document_id
		WHERE d.status = 'published' AND d.category <> 'script'
		ORDER BY d.category, d.title, d.id, c.seq`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Chunk])
}

type Script struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"-"`
}

// Scripts lists published 共修腳本; withBody=false for the public picker.
func (s *Store) Scripts(ctx context.Context) ([]Script, error) {
	rows, err := s.DB.Query(ctx, `
		SELECT d.id, d.title, v.body FROM guide.document d
		JOIN guide.document_version v ON v.document_id = d.id AND v.version = d.current_version
		WHERE d.status = 'published' AND d.category = 'script' ORDER BY d.title`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Script])
}

func (s *Store) Script(ctx context.Context, id string) (Script, error) {
	var sc Script
	err := s.DB.QueryRow(ctx, `
		SELECT d.id, d.title, v.body FROM guide.document d
		JOIN guide.document_version v ON v.document_id = d.id AND v.version = d.current_version
		WHERE d.id::text = $1 AND d.status = 'published' AND d.category = 'script'`, id).Scan(&sc.ID, &sc.Title, &sc.Body)
	if errors.Is(err, pgx.ErrNoRows) {
		return sc, errNotFound
	}
	return sc, err
}

func (s *Store) Persona(ctx context.Context) (Persona, error) {
	var p Persona
	err := s.DB.QueryRow(ctx, `
		SELECT p.version, p.name, p.prompt, u.display_name, p.created_at
		FROM guide.persona_version p LEFT JOIN app_user u ON u.id = p.edited_by
		ORDER BY p.version DESC LIMIT 1`).Scan(&p.Version, &p.Name, &p.Prompt, &p.EditedBy, &p.CreatedAt)
	return p, err
}

func (s *Store) PersonaHistory(ctx context.Context) ([]Persona, error) {
	rows, err := s.DB.Query(ctx, `
		SELECT p.version, p.name, p.prompt, u.display_name, p.created_at
		FROM guide.persona_version p LEFT JOIN app_user u ON u.id = p.edited_by
		ORDER BY p.version DESC LIMIT 50`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Persona])
}

func (s *Store) SavePersona(ctx context.Context, name, prompt string, userID *string) error {
	_, err := s.DB.Exec(ctx, `INSERT INTO guide.persona_version (name, prompt, edited_by) VALUES ($1, $2, $3)`, name, prompt, userID)
	return err
}

// EnsurePersona seeds the default persona the first time the service starts.
func (s *Store) EnsurePersona(ctx context.Context, name string) error {
	var n int
	if err := s.DB.QueryRow(ctx, `SELECT count(*) FROM guide.persona_version`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	return s.SavePersona(ctx, name, DefaultPersonaPrompt(name), nil)
}

func (s *Store) Setting(ctx context.Context, key string) (string, error) {
	var v string
	err := s.DB.QueryRow(ctx, `SELECT value FROM guide.setting WHERE key = $1`, key).Scan(&v)
	return v, err
}

func (s *Store) SetSetting(ctx context.Context, key, value string) error {
	_, err := s.DB.Exec(ctx, `
		INSERT INTO guide.setting (key, value) VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET value = excluded.value`, key, value)
	return err
}

func (s *Store) RecordUsage(ctx context.Context, model string, in, out, cacheRead, cacheWrite int64) error {
	_, err := s.DB.Exec(ctx, `
		INSERT INTO guide.usage (model, input_tokens, output_tokens, cache_read_tokens, cache_write_tokens)
		VALUES ($1, $2, $3, $4, $5)`, model, in, out, cacheRead, cacheWrite)
	return err
}

type MonthUsage struct {
	Requests  int64   `json:"requests"`
	CostUSD   float64 `json:"cost_usd"`
	BudgetUSD float64 `json:"budget_usd"`
	Enabled   bool    `json:"enabled"`
}

// Usage sums this calendar month's estimated cost (UTC) against the budget setting.
func (s *Store) Usage(ctx context.Context) (MonthUsage, error) {
	var u MonthUsage
	rows, err := s.DB.Query(ctx, `
		SELECT model, count(*), sum(input_tokens), sum(output_tokens), sum(cache_read_tokens), sum(cache_write_tokens)
		FROM guide.usage WHERE created_at >= date_trunc('month', now() AT TIME ZONE 'UTC') AT TIME ZONE 'UTC'
		GROUP BY model`)
	if err != nil {
		return u, err
	}
	defer rows.Close()
	for rows.Next() {
		var model string
		var n, in, out, cr, cw int64
		if err := rows.Scan(&model, &n, &in, &out, &cr, &cw); err != nil {
			return u, err
		}
		u.Requests += n
		u.CostUSD += estimateCost(model, in, out, cr, cw)
	}
	if err := rows.Err(); err != nil {
		return u, err
	}
	budget, err := s.Setting(ctx, "monthly_budget_usd")
	if err != nil {
		return u, err
	}
	if u.BudgetUSD, err = strconv.ParseFloat(budget, 64); err != nil {
		return u, fmt.Errorf("monthly_budget_usd: %w", err)
	}
	enabled, err := s.Setting(ctx, "chat_enabled")
	u.Enabled = enabled == "true"
	return u, err
}

// Prices in USD per million tokens. Unknown models are priced as Claude Opus 5.5.
var prices = map[string][2]float64{
	"claude-opus-5-5":   {4, 20},
	"claude-sonnet-5-5": {2, 10},
	"claude-haiku-4-5":  {1, 5},
	"claude-fable-5-1":  {10, 50},
}

func estimateCost(model string, in, out, cacheRead, cacheWrite int64) float64 {
	p, ok := prices[model]
	if !ok {
		p = prices["claude-opus-5-5"]
	}
	// Cache reads cost 0.1x input; 5-minute cache writes 1.25x input.
	return (float64(in)*p[0] + float64(out)*p[1] + float64(cacheRead)*p[0]*0.1 + float64(cacheWrite)*p[0]*1.25) / 1e6
}
