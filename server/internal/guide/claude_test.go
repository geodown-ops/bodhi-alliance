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
