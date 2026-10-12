package guide

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// TestClaudeRequestShape checks what we send to the Messages API, against a stub server.
func TestClaudeRequestShape(t *testing.T) {
	var got map[string]any
	var beta string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &got)
		beta = r.Header.Get("anthropic-beta")
		w.Header().Set("Content-Type", "text/event-stream")
		for _, ev := range []string{
			`{"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"claude-opus-5-5","content":[],"stop_reason":null,"usage":{"input_tokens":10,"output_tokens":0,"cache_read_input_tokens":3}}}`,
			`{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
			`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"  我是線上覺行"}}`,
			`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"小組組長Sunny  "}}`,
			`{"type":"content_block_stop","index":0}`,
			`{"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":5}}`,
			`{"type":"message_stop"}`,
		} {
			var head struct{ Type string }
			json.Unmarshal([]byte(ev), &head)
			io.WriteString(w, "event: "+head.Type+"\ndata: "+ev+"\n\n")
		}
	}))
	defer srv.Close()

	c := &Claude{
		Client: anthropic.NewClient(option.WithAPIKey("test"), option.WithBaseURL(srv.URL), option.WithMaxRetries(0)),
		Model:  "claude-opus-5-5",
		Effort: anthropic.BetaOutputConfigEffortMedium,
	}
	var streamed strings.Builder
	reply, usage, err := c.Complete(context.Background(), Turn{
		Persona:   "persona",
		Knowledge: "knowledge",
		Script:    &Script{Title: "十分鐘靜坐", Body: "一、坐好"},
		Messages:  []Message{{Role: "user", Content: "你是誰"}},
	}, func(text string) { streamed.WriteString(text) })
	if err != nil {
		t.Fatal(err)
	}
	if reply.Text != "我是線上覺行小組組長Sunny" || usage.CacheRead != 3 || usage.Model != "claude-opus-5-5" {
		t.Errorf("reply %+v usage %+v", reply, usage)
	}
	if streamed.String() != "  我是線上覺行小組組長Sunny  " {
		t.Errorf("streamed %q", streamed.String())
	}
	if got["stream"] != true {
		t.Error("request should stream")
	}
	if got["model"] != "claude-opus-5-5" || got["fallbacks"] != "default" {
		t.Errorf("model/fallbacks = %v / %v", got["model"], got["fallbacks"])
	}
	if beta != "server-side-fallback-2026-07-01" {
		t.Errorf("beta header = %q", beta)
	}
	if oc, _ := got["output_config"].(map[string]any); oc["effort"] != "medium" {
		t.Errorf("output_config = %v", got["output_config"])
	}
	system, _ := got["system"].([]any)
	if len(system) != 4 {
		t.Fatalf("system = %v", got["system"])
	}
	if style, _ := system[1].(map[string]any); style["text"] != replyStyle {
		t.Error("reply style should follow the persona")
	}
	if kb, _ := system[2].(map[string]any); kb["cache_control"] == nil {
		t.Error("knowledge block should carry the cache breakpoint")
	}
	if _, has := got["thinking"]; has {
		t.Error("thinking should be left unset on Opus 5.5")
	}
}

type fakeReporter struct{ summary, email string }

func (f *fakeReporter) Record(_ context.Context, summary, email string) error {
	f.summary, f.email = summary, email
	return nil
}
func (f *fakeReporter) SetEmail(_ context.Context, email string) error { f.email = email; return nil }

// TestClaudeRecordsIssue checks the tool round trip: Sunny calls record_site_issue, we
// record it and send the result back, and Sunny's text from both rounds is the answer.
func TestClaudeRecordsIssue(t *testing.T) {
	var bodies []map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var got map[string]any
		json.Unmarshal(raw, &got)
		bodies = append(bodies, got)
		w.Header().Set("Content-Type", "text/event-stream")
		events := []string{
			`{"type":"message_start","message":{"id":"msg_2","type":"message","role":"assistant","model":"claude-opus-5-5","content":[],"stop_reason":null,"usage":{"input_tokens":10,"output_tokens":0}}}`,
			`{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
			`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"我幫你記下來。"}}`,
			`{"type":"content_block_stop","index":0}`,
			`{"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"toolu_1","name":"record_site_issue","input":{}}}`,
			`{"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"{\"summary\":\"送出報名沒反應\",\"contact_email\":\"a@b.org\"}"}}`,
			`{"type":"content_block_stop","index":1}`,
			`{"type":"message_delta","delta":{"stop_reason":"tool_use"},"usage":{"output_tokens":20}}`,
			`{"type":"message_stop"}`,
		}
		if len(bodies) > 1 {
			events = []string{
				`{"type":"message_start","message":{"id":"msg_3","type":"message","role":"assistant","model":"claude-opus-5-5","content":[],"stop_reason":null,"usage":{"input_tokens":30,"output_tokens":0}}}`,
				`{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
				`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"管理員會回覆你。"}}`,
				`{"type":"content_block_stop","index":0}`,
				`{"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":5}}`,
				`{"type":"message_stop"}`,
			}
		}
		for _, ev := range events {
			var head struct{ Type string }
			json.Unmarshal([]byte(ev), &head)
			io.WriteString(w, "event: "+head.Type+"\ndata: "+ev+"\n\n")
		}
	}))
	defer srv.Close()

	c := &Claude{
		Client: anthropic.NewClient(option.WithAPIKey("test"), option.WithBaseURL(srv.URL), option.WithMaxRetries(0)),
		Model:  "claude-opus-5-5",
		Effort: anthropic.BetaOutputConfigEffortMedium,
	}
	rep := &fakeReporter{}
	var streamed strings.Builder
	reply, usage, err := c.Complete(context.Background(), Turn{
		Persona: "persona", Knowledge: "knowledge", Visitor: "對方目前沒有登入官網。", Issues: rep,
		Messages: []Message{{Role: "user", Content: "送出報名沒反應"}},
	}, func(text string) { streamed.WriteString(text) })
	if err != nil {
		t.Fatal(err)
	}
	if rep.summary != "送出報名沒反應" || rep.email != "a@b.org" {
		t.Errorf("recorded %+v", rep)
	}
	if reply.Text != "我幫你記下來。管理員會回覆你。" || streamed.String() != reply.Text {
		t.Errorf("reply %q streamed %q", reply.Text, streamed.String())
	}
	if usage.Input != 40 || usage.Output != 25 {
		t.Errorf("usage %+v", usage)
	}
	if len(bodies) != 2 {
		t.Fatalf("%d requests", len(bodies))
	}
	if tools, _ := bodies[0]["tools"].([]any); len(tools) != 2 {
		t.Errorf("tools = %v", bodies[0]["tools"])
	}
	msgs, _ := bodies[1]["messages"].([]any)
	if len(msgs) != 3 {
		t.Fatalf("second request messages = %v", msgs)
	}
	last, _ := msgs[2].(map[string]any)
	content, _ := last["content"].([]any)
	result, _ := content[0].(map[string]any)
	if result["type"] != "tool_result" || result["tool_use_id"] != "toolu_1" {
		t.Errorf("tool result = %v", result)
	}
}
