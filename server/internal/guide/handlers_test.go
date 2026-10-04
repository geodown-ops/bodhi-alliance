package guide

import (
	"bytes"
	"context"
	"encoding/json"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"

	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
	"github.com/geodown-ops/bodhi-alliance/server/internal/testdb"
)

type fakeLLM struct{ last Turn }

func (f *fakeLLM) Complete(_ context.Context, t Turn, onText func(string)) (Reply, Usage, error) {
	if onText != nil {
		onText("好")
		onText("的")
	}
	f.last = t
	return Reply{Text: "好的"}, Usage{Model: "claude-opus-5-5", Input: 1000, Output: 100}, nil
}

type env struct {
	t     *testing.T
	r     *gin.Engine
	llm   *fakeLLM
	store *Store
	token string
}

func setup(t *testing.T) *env {
	gin.SetMode(gin.TestMode)
	pool := testdb.New(t)
	ctx := context.Background()
	authSvc := &auth.Service{DB: pool}
	if _, err := authSvc.CreateUser(ctx, "km@example.org", "知識管理員", "correct-horse-battery",
		auth.Role{Role: auth.RoleKnowledgeManager, Scope: auth.ScopeGuide}); err != nil {
		t.Fatal(err)
	}
	token, err := authSvc.Login(ctx, "km@example.org", "correct-horse-battery")
	if err != nil {
		t.Fatal(err)
	}
	store := &Store{DB: pool}
	if err := store.EnsurePersona(ctx, "Sunny"); err != nil {
		t.Fatal(err)
	}
	llm := &fakeLLM{}
	r := gin.New()
	(&Handler{Store: store, Auth: authSvc, LLM: llm}).Routes(r.Group("/guide"))
	return &env{t: t, r: r, llm: llm, store: store, token: token}
}

func (e *env) do(method, path string, body any, authed bool) *httptest.ResponseRecorder {
	e.t.Helper()
	var buf bytes.Buffer
	ct := "application/json"
	switch b := body.(type) {
	case nil:
	case *multipart.Writer:
		e.t.Fatal("use upload")
	case []byte:
		buf.Write(b)
	default:
		json.NewEncoder(&buf).Encode(b)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", ct)
	if authed {
		req.Header.Set("Authorization", "Bearer "+e.token)
	}
	w := httptest.NewRecorder()
	e.r.ServeHTTP(w, req)
	return w
}

func (e *env) upload(title, category, filename, content string) string {
	e.t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	mw.WriteField("title", title)
	mw.WriteField("category", category)
	fw, _ := mw.CreateFormFile("file", filename)
	fw.Write([]byte(content))
	mw.Close()
	req := httptest.NewRequest(http.MethodPost, "/guide/admin/documents", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+e.token)
	w := httptest.NewRecorder()
	e.r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		e.t.Fatalf("upload: %d %s", w.Code, w.Body)
	}
	var out struct{ ID string }
	json.Unmarshal(w.Body.Bytes(), &out)
	return out.ID
}

func TestAdminNeedsLogin(t *testing.T) {
	e := setup(t)
	if w := e.do(http.MethodGet, "/guide/admin/documents", nil, false); w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want 401", w.Code)
	}
}

func TestUploadPublishAndChat(t *testing.T) {
	e := setup(t)
	id := e.upload("", "practice", "靜坐入門.md", "# 調息\n\n吸氣數一，呼氣數二。")

	// Drafts are not used by the public chat.
	msg := map[string]any{"messages": []Message{{Role: "user", Content: "怎麼調息？"}}}
	if w := e.do(http.MethodPost, "/guide/chat", msg, false); w.Code != http.StatusOK {
		t.Fatalf("chat: %d %s", w.Code, w.Body)
	}
	if strings.Contains(e.llm.last.Knowledge, "吸氣數一") {
		t.Error("draft leaked into public chat")
	}
	if !strings.Contains(e.llm.last.Persona, "我是線上覺行小組組長Sunny") {
		t.Error("persona missing self-introduction")
	}

	// 試問 can include the draft.
	try := map[string]any{"messages": msg["messages"], "include_document_id": id}
	if w := e.do(http.MethodPost, "/guide/admin/try", try, true); w.Code != http.StatusOK {
		t.Fatalf("try: %d %s", w.Code, w.Body)
	}
	if !strings.Contains(e.llm.last.Knowledge, "吸氣數一") || !strings.Contains(e.llm.last.Knowledge, `title="靜坐入門"`) {
		t.Errorf("try knowledge = %q", e.llm.last.Knowledge)
	}

	if w := e.do(http.MethodPost, "/guide/admin/documents/"+id+"/status", map[string]string{"status": "published"}, true); w.Code != http.StatusNoContent {
		t.Fatalf("publish: %d %s", w.Code, w.Body)
	}
	e.do(http.MethodPost, "/guide/chat", msg, false)
	if !strings.Contains(e.llm.last.Knowledge, "## 調息") {
		t.Errorf("published doc not in knowledge: %q", e.llm.last.Knowledge)
	}

	// Editing a published document creates a version and is live immediately.
	edit := map[string]string{"title": "靜坐入門", "category": "practice", "body": "# 調息\n\n吸氣數一，呼氣數二，數到十再從一開始。"}
	if w := e.do(http.MethodPut, "/guide/admin/documents/"+id, edit, true); w.Code != http.StatusNoContent {
		t.Fatalf("edit: %d %s", w.Code, w.Body)
	}
	e.do(http.MethodPost, "/guide/chat", msg, false)
	if !strings.Contains(e.llm.last.Knowledge, "數到十") {
		t.Error("edit not live")
	}
	w := e.do(http.MethodGet, "/guide/admin/documents/"+id+"/versions", nil, true)
	var versions []Version
	json.Unmarshal(w.Body.Bytes(), &versions)
	if len(versions) != 2 || versions[0].Version != 2 {
		t.Errorf("versions = %+v", versions)
	}

	// Unpublishing removes it.
	e.do(http.MethodPost, "/guide/admin/documents/"+id+"/status", map[string]string{"status": "draft"}, true)
	e.do(http.MethodPost, "/guide/chat", msg, false)
	if strings.Contains(e.llm.last.Knowledge, "調息") {
		t.Error("unpublished doc still used")
	}
}

func TestPracticeModeUsesScript(t *testing.T) {
	e := setup(t)
	id := e.upload("十分鐘靜坐", "script", "s.txt", "一、坐好\n\n二、數息五分鐘\n\n三、回向")
	e.do(http.MethodPost, "/guide/admin/documents/"+id+"/status", map[string]string{"status": "published"}, true)

	w := e.do(http.MethodGet, "/guide/scripts", nil, false)
	if !strings.Contains(w.Body.String(), "十分鐘靜坐") || strings.Contains(w.Body.String(), "數息") {
		t.Errorf("scripts list = %s (should list titles only)", w.Body)
	}
	msg := map[string]any{"mode": "practice", "script_id": id, "messages": []Message{{Role: "user", Content: "開始"}}}
	if w := e.do(http.MethodPost, "/guide/chat", msg, false); w.Code != http.StatusOK {
		t.Fatalf("chat: %d %s", w.Code, w.Body)
	}
	if e.llm.last.Script == nil || !strings.Contains(e.llm.last.Script.Body, "數息五分鐘") {
		t.Error("script not passed to model")
	}
	if strings.Contains(e.llm.last.Knowledge, "數息") {
		t.Error("scripts should not be mixed into general knowledge")
	}
}

func TestBudgetStopsChat(t *testing.T) {
	e := setup(t)
	if w := e.do(http.MethodPut, "/guide/admin/settings", map[string]float64{"monthly_budget_usd": 0.01}, true); w.Code != http.StatusNoContent {
		t.Fatalf("settings: %d", w.Code)
	}
	msg := map[string]any{"messages": []Message{{Role: "user", Content: "你好"}}}
	if w := e.do(http.MethodPost, "/guide/chat", msg, false); w.Code != http.StatusOK {
		t.Fatalf("first chat: %d", w.Code)
	}
	// The fake spends $0.006 per call; the second pushes past $0.01.
	e.do(http.MethodPost, "/guide/chat", msg, false)
	if w := e.do(http.MethodPost, "/guide/chat", msg, false); w.Code != http.StatusServiceUnavailable {
		t.Errorf("over budget: got %d, want 503", w.Code)
	}
	w := e.do(http.MethodGet, "/guide/info", nil, false)
	if !strings.Contains(w.Body.String(), `"available":false`) {
		t.Errorf("info = %s", w.Body)
	}
}

func TestTryDraftScriptRunsPractice(t *testing.T) {
	e := setup(t)
	id := e.upload("行禪", "script", "w.md", "一、站立\n\n二、慢步")
	try := map[string]any{"include_document_id": id, "messages": []Message{{Role: "user", Content: "開始"}}}
	if w := e.do(http.MethodPost, "/guide/admin/try", try, true); w.Code != http.StatusOK {
		t.Fatalf("try: %d %s", w.Code, w.Body)
	}
	if e.llm.last.Script == nil || e.llm.last.Script.Title != "行禪" {
		t.Errorf("draft script not used: %+v", e.llm.last.Script)
	}
	if strings.Contains(e.llm.last.Knowledge, "慢步") {
		t.Error("script leaked into knowledge")
	}
}

func TestChatStreams(t *testing.T) {
	e := setup(t)
	body, _ := json.Marshal(map[string]any{"messages": []Message{{Role: "user", Content: "你好"}}})
	req := httptest.NewRequest(http.MethodPost, "/guide/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	w := httptest.NewRecorder()
	e.r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/event-stream") {
		t.Fatalf("status %d, type %q", w.Code, w.Header().Get("Content-Type"))
	}
	want := "event:delta\ndata:{\"text\":\"好\"}\n\nevent:delta\ndata:{\"text\":\"的\"}\n\nevent:done\ndata:{\"reply\":\"好的\"}\n\n"
	if w.Body.String() != want {
		t.Errorf("body = %q", w.Body.String())
	}
}

func seedFS(body string) fstest.MapFS {
	return fstest.MapFS{"seed/coin.md": {Data: []byte("---\ntitle: 菩提幣與共好企業\ncategory: coin\n---\n" + body)}}
}

func (e *env) seed(fsys fs.FS, wantAdded, wantUpdated int) {
	e.t.Helper()
	if a, u, err := e.store.SeedDocuments(context.Background(), fsys); err != nil || a != wantAdded || u != wantUpdated {
		e.t.Fatalf("seed: added=%d updated=%d err=%v, want %d %d", a, u, err, wantAdded, wantUpdated)
	}
}

func (e *env) ask(q string) string {
	e.t.Helper()
	e.do(http.MethodPost, "/guide/chat", map[string]any{"messages": []Message{{Role: "user", Content: q}}}, false)
	return e.llm.last.Knowledge
}

func TestSeedDocumentsOnce(t *testing.T) {
	e := setup(t)
	ctx := context.Background()
	e.seed(seedFS("# 梯級\n\n半日服務：1,500 幣"), 1, 0)
	if k := e.ask("半日服務幾幣？"); !strings.Contains(k, "1,500 幣") || !strings.Contains(k, `title="菩提幣與共好企業"`) {
		t.Fatalf("seeded doc not live: %q", k)
	}

	// An admin archiving it is not undone by the next start.
	docs, err := e.store.ListDocuments(ctx)
	if err != nil || len(docs) != 1 || docs[0].SourceName != "seed/coin.md" {
		t.Fatalf("docs = %+v, %v", docs, err)
	}
	if err := e.store.SetStatus(ctx, docs[0].ID, "archived"); err != nil {
		t.Fatal(err)
	}
	e.seed(seedFS("# 梯級\n\n半日服務：1,500 幣"), 0, 0)
	if docs, _ := e.store.ListDocuments(ctx); len(docs) != 0 {
		t.Errorf("archived seed came back: %+v", docs)
	}
}

func TestSeedFileChangesFollowUntilEdited(t *testing.T) {
	e := setup(t)
	ctx := context.Background()
	e.seed(seedFS("# 審核\n\n由主辦審核小組審核。"), 1, 0)

	// A changed file reaches the published document as a new version.
	e.seed(seedFS("# 審核\n\n由菩提幣決策小組審核。"), 0, 1)
	if k := e.ask("誰審核？"); !strings.Contains(k, "菩提幣決策小組") || strings.Contains(k, "主辦審核小組") {
		t.Fatalf("seed change not live: %q", k)
	}
	docs, _ := e.store.ListDocuments(ctx)
	if len(docs) != 1 || docs[0].CurrentVersion != 2 {
		t.Fatalf("docs = %+v", docs)
	}

	// Once someone edits it in the admin, the file no longer touches it.
	var userID string
	if err := e.store.DB.QueryRow(ctx, `SELECT id FROM app_user WHERE email = 'km@example.org'`).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	if err := e.store.UpdateDocument(ctx, docs[0].ID, docs[0].Title, "coin", "# 審核\n\n管理員改過的說明。", "edit", userID); err != nil {
		t.Fatal(err)
	}
	e.seed(seedFS("# 審核\n\n第三版。"), 0, 0)
	if k := e.ask("誰審核？"); !strings.Contains(k, "管理員改過的說明") {
		t.Errorf("admin edit overwritten: %q", k)
	}
}

func TestSeedDigestMissing(t *testing.T) {
	e := setup(t)
	ctx := context.Background()
	e.seed(seedFS("舊版"), 1, 0)
	// Deployments seeded before the digest existed stored only the document id.
	id, _ := e.store.Setting(ctx, "seed:coin.md")
	id, _, _ = strings.Cut(id, " ")
	if err := e.store.SetSetting(ctx, "seed:coin.md", id); err != nil {
		t.Fatal(err)
	}
	e.seed(seedFS("新版"), 0, 1)
	if v, _ := e.store.Setting(ctx, "seed:coin.md"); v != id+" "+digest("新版") {
		t.Errorf("setting = %q", v)
	}
}

func TestBundledSeedPublishes(t *testing.T) {
	e := setup(t)
	files, _ := fs.Glob(SeedFiles, "seed/*.md")
	e.seed(SeedFiles, len(files), 0)
	if k := e.ask("誰審核核發名單？"); !strings.Contains(k, "菩提幣決策小組") {
		t.Error("bundled knowledge not live")
	}
}

func TestPersonaFollowsDefaultUntilEdited(t *testing.T) {
	e := setup(t)
	ctx := context.Background()
	p, _ := e.store.Persona(ctx)
	if p.Prompt != DefaultPersonaPrompt("Sunny") {
		t.Fatalf("fresh persona = %q", p.Prompt)
	}

	// A deployment seeded with an older default (and no digest) moves to the current one.
	if _, err := e.store.DB.Exec(ctx, `DELETE FROM guide.persona_version; DELETE FROM guide.setting WHERE key = 'seed:persona';
		ALTER SEQUENCE guide.persona_version_version_seq RESTART`); err != nil {
		t.Fatal(err)
	}
	old := strings.ReplaceAll(DefaultPersonaPrompt("Sunny"), "「我的錢包」", "「菩提幣錢包」")
	if err := e.store.SavePersona(ctx, "Sunny", old, nil); err != nil {
		t.Fatal(err)
	}
	if err := e.store.EnsurePersona(ctx, "Sunny"); err != nil {
		t.Fatal(err)
	}
	if p, _ := e.store.Persona(ctx); p.Version != 2 || p.Prompt != DefaultPersonaPrompt("Sunny") {
		t.Fatalf("persona not upgraded: v%d %q", p.Version, p.Prompt)
	}

	// A persona saved in the admin stays.
	var userID string
	e.store.DB.QueryRow(ctx, `SELECT id FROM app_user WHERE email = 'km@example.org'`).Scan(&userID)
	if err := e.store.SavePersona(ctx, "Sunny", "管理員寫的設定", &userID); err != nil {
		t.Fatal(err)
	}
	if err := e.store.EnsurePersona(ctx, "Sunny"); err != nil {
		t.Fatal(err)
	}
	if p, _ := e.store.Persona(ctx); p.Prompt != "管理員寫的設定" {
		t.Errorf("admin persona replaced: %q", p.Prompt)
	}
}
