package guide

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
)

func (e *env) login(email string, roles ...auth.Role) string {
	e.t.Helper()
	svc := &auth.Service{DB: e.store.DB}
	ctx := context.Background()
	if _, err := svc.CreateUser(ctx, email, email, "correct-horse-battery", roles...); err != nil {
		e.t.Fatal(err)
	}
	token, err := svc.Login(ctx, email, "correct-horse-battery")
	if err != nil {
		e.t.Fatal(err)
	}
	return token
}

func (e *env) as(token string, fn func()) {
	saved := e.token
	e.token = token
	defer func() { e.token = saved }()
	fn()
}

func (e *env) issues(token string) []SiteIssue {
	e.t.Helper()
	var out []SiteIssue
	e.as(token, func() {
		w := e.do(http.MethodGet, "/guide/admin/issues", nil, true)
		if w.Code != http.StatusOK {
			e.t.Fatalf("list issues: %d %s", w.Code, w.Body)
		}
		json.Unmarshal(w.Body.Bytes(), &out)
	})
	return out
}

func TestVisitorIssueWithEmail(t *testing.T) {
	e := setup(t)
	super := e.login("super@example.org", auth.Role{Role: auth.RoleAllianceAdmin, Scope: auth.ScopeAlliance})

	e.llm.report = "加入會員時按送出報名沒有反應"
	chat := map[string]any{"messages": []Message{{Role: "user", Content: "我按送出報名都沒反應"}}, "conversation": "c1"}
	if w := e.do(http.MethodPost, "/guide/chat", chat, false); w.Code != http.StatusOK {
		t.Fatalf("chat: %d %s", w.Code, w.Body)
	}
	if e.llm.toolErr != nil {
		t.Fatal(e.llm.toolErr)
	}
	if e.llm.last.Visitor != "對方目前沒有登入官網。" {
		t.Errorf("visitor = %q", e.llm.last.Visitor)
	}

	e.llm.report, e.llm.email = "", "Visitor@Example.org"
	if w := e.do(http.MethodPost, "/guide/chat", chat, false); w.Code != http.StatusOK || e.llm.toolErr != nil {
		t.Fatalf("email: %d %v", w.Code, e.llm.toolErr)
	}

	list := e.issues(super)
	if len(list) != 1 || list[0].ContactEmail != "visitor@example.org" || list[0].Summary != "加入會員時按送出報名沒有反應" ||
		list[0].Status != "open" || len(list[0].Transcript) != 1 || list[0].UserID != nil {
		t.Fatalf("issues = %+v", list)
	}

	e.as(super, func() {
		if w := e.do(http.MethodPut, "/guide/admin/issues/"+list[0].ID+"/answer", map[string]string{"answer": "已修好，請再試一次"}, true); w.Code != http.StatusNoContent {
			t.Fatalf("answer: %d %s", w.Code, w.Body)
		}
	})
	if got := e.issues(super)[0]; got.Status != "answered" || got.Answer != "已修好，請再試一次" || got.AnsweredBy == "" {
		t.Errorf("after answer = %+v", got)
	}
}

func TestMemberSeesAnswer(t *testing.T) {
	e := setup(t)
	super := e.login("super@example.org", auth.Role{Role: auth.RoleAllianceAdmin, Scope: auth.ScopeAlliance})
	member := e.login("member@example.org")
	other := e.login("other@example.org")

	e.llm.report = "個人頁看不到錢包餘額"
	e.as(member, func() {
		chat := map[string]any{"messages": []Message{{Role: "user", Content: "我的錢包是空的？"}}}
		if w := e.do(http.MethodPost, "/guide/chat", chat, true); w.Code != http.StatusOK || e.llm.toolErr != nil {
			t.Fatalf("chat: %d %v", w.Code, e.llm.toolErr)
		}
	})
	if e.llm.last.Visitor != "對方已登入官網，是會員「member@example.org」。" {
		t.Errorf("visitor = %q", e.llm.last.Visitor)
	}
	id := e.issues(super)[0].ID
	e.as(super, func() {
		e.do(http.MethodPut, "/guide/admin/issues/"+id+"/answer", map[string]string{"answer": "餘額在活動審核通過後才會入帳"}, true)
	})

	var mine []SiteIssue
	e.as(member, func() {
		w := e.do(http.MethodGet, "/guide/me/issues", nil, true)
		json.Unmarshal(w.Body.Bytes(), &mine)
	})
	if len(mine) != 1 || mine[0].Answer != "餘額在活動審核通過後才會入帳" {
		t.Errorf("member issues = %+v", mine)
	}
	e.as(other, func() {
		w := e.do(http.MethodGet, "/guide/me/issues", nil, true)
		if w.Body.String() != "[]" {
			t.Errorf("other member sees %s", w.Body)
		}
	})
}

func TestIssuesAreSuperAdminOnly(t *testing.T) {
	e := setup(t)
	manager := e.login("manager@example.org", auth.Role{Role: auth.RoleManager, Scope: auth.ScopeAlliance})
	member := e.login("member@example.org")
	for name, token := range map[string]string{"knowledge manager": e.token, "manager": manager, "member": member} {
		e.as(token, func() {
			if w := e.do(http.MethodGet, "/guide/admin/issues", nil, true); w.Code != http.StatusForbidden {
				t.Errorf("%s: got %d, want 403", name, w.Code)
			}
		})
	}
	if w := e.do(http.MethodGet, "/guide/admin/issues", nil, false); w.Code != http.StatusUnauthorized {
		t.Errorf("anonymous: got %d", w.Code)
	}
}

func TestIssueLimitsAndTryHasNoTools(t *testing.T) {
	e := setup(t)
	e.llm.report = "問題"
	chat := map[string]any{"messages": []Message{{Role: "user", Content: "壞了"}}, "conversation": "spam"}
	for i := 0; i < maxIssuesPerChat; i++ {
		e.do(http.MethodPost, "/guide/chat", chat, false)
		if e.llm.toolErr != nil {
			t.Fatalf("issue %d: %v", i, e.llm.toolErr)
		}
	}
	e.do(http.MethodPost, "/guide/chat", chat, false)
	if e.llm.toolErr != errIssueLimit {
		t.Errorf("over limit: %v", e.llm.toolErr)
	}

	e.do(http.MethodPost, "/guide/admin/try", map[string]any{"messages": []Message{{Role: "user", Content: "壞了"}}}, true)
	if e.llm.last.Issues != nil {
		t.Error("admin 試問 should not record issues")
	}

	e.llm.report, e.llm.email = "", "not-an-email"
	e.do(http.MethodPost, "/guide/chat", chat, false)
	if e.llm.toolErr == nil {
		t.Error("bad email accepted")
	}
}
