package guide

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// DefaultPersonaPrompt is the starting persona; knowledge managers edit it in the admin.
func DefaultPersonaPrompt(name string) string {
	return fmt.Sprintf(`你是世界佛教教育協會「線上覺行小組」的組長，名字叫 %[1]s。
有人問你是誰時，回答：「我是線上覺行小組組長%[1]s」。

你的工作：
- 回答覺行、共修、菩提幣與官網使用的問題。
- 帶領大家共修：依照共修腳本一步一步引導，一次只給一個步驟，等對方回應再往下。
- 用溫和、簡短、口語的繁體中文回答，像一位親切的小組組長，不說教。

界線：
- 只根據「知識庫」裡的內容回答佛法、覺行與菩提幣的問題。知識庫沒有寫的，就老實說目前沒有這方面的資料，並建議到官網「覺行小組」頁面聯絡真人組長。
- 引用知識庫時，在句末用《文件標題》標出出處。
- 你看不到也改不了任何人的菩提幣餘額或券，問到時請對方到官網「菩提幣錢包」查看。
- 菩提幣不販售、不提領、不可兌現，也不是投資；不要做任何價格或收益的預測。
- 不做醫療或心理診斷。若對方提到想傷害自己或他人，溫和地請對方立刻尋求協助：台灣安心專線 1925、生命線 1995，緊急狀況撥 119／110。`, name)
}

const (
	modeChat     = "chat"
	modePractice = "practice"

	// Below this size the whole knowledge base goes into the (cached) system prompt;
	// above it, only the best-matching chunks for the latest question are sent.
	fullContextLimit = 150_000 // runes
	retrieveK        = 12
	maxHistory       = 20
	maxMessageRunes  = 2000
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Reply struct {
	Text    string `json:"reply"`
	Refused bool   `json:"refused,omitempty"`
}

// Turn is everything one model call needs.
type Turn struct {
	Persona   string
	Knowledge string
	Script    *Script
	Messages  []Message
}

// LLM is the model behind the guide; tests use a fake.
type LLM interface {
	Complete(ctx context.Context, t Turn) (Reply, Usage, error)
}

type Usage struct {
	Model                            string
	Input, Output, CacheRead, CacheW int64
}

// BuildKnowledge renders chunks as the knowledge-base section of the system prompt.
func BuildKnowledge(chunks []Chunk) string {
	if len(chunks) == 0 {
		return "（知識庫目前是空的。）"
	}
	var b strings.Builder
	b.WriteString("以下是知識庫內容，只能根據這些內容回答。\n")
	last := ""
	for _, c := range chunks {
		if c.DocumentID != last {
			fmt.Fprintf(&b, "\n<document title=%q category=%q>\n", c.Title, categories[c.Category])
			last = c.DocumentID
		}
		if c.Heading != "" {
			fmt.Fprintf(&b, "## %s\n", c.Heading)
		}
		b.WriteString(c.Body)
		b.WriteString("\n")
	}
	return b.String()
}

// SelectKnowledge sends everything when it fits, otherwise the chunks that best match the
// latest user message (in their original document order).
func SelectKnowledge(all []Chunk, latest string) []Chunk {
	total := 0
	for _, c := range all {
		total += len([]rune(c.Body))
	}
	if total <= fullContextLimit {
		return all
	}
	picked := map[[2]any]bool{}
	for _, c := range Rank(all, latest, retrieveK) {
		picked[[2]any{c.DocumentID, c.Seq}] = true
	}
	var out []Chunk
	for _, c := range all {
		if picked[[2]any{c.DocumentID, c.Seq}] {
			out = append(out, c)
		}
	}
	return out
}

// CleanHistory keeps the last maxHistory turns, trims long messages, drops empty or
// unknown-role ones, and makes sure the conversation starts and ends with the user.
func CleanHistory(in []Message) ([]Message, error) {
	var out []Message
	for _, m := range in {
		m.Content = strings.TrimSpace(m.Content)
		if m.Content == "" || (m.Role != "user" && m.Role != "assistant") {
			continue
		}
		if r := []rune(m.Content); len(r) > maxMessageRunes {
			m.Content = string(r[:maxMessageRunes])
		}
		// Merge consecutive same-role messages; the API expects alternating turns.
		if n := len(out); n > 0 && out[n-1].Role == m.Role {
			out[n-1].Content += "\n\n" + m.Content
			continue
		}
		out = append(out, m)
	}
	if len(out) > maxHistory {
		out = out[len(out)-maxHistory:]
	}
	for len(out) > 0 && out[0].Role != "user" {
		out = out[1:]
	}
	if len(out) == 0 || out[len(out)-1].Role != "user" {
		return nil, errors.New("請輸入訊息")
	}
	return out, nil
}

// Claude calls the Anthropic Messages API.
type Claude struct {
	Client anthropic.Client
	Model  string
	Effort anthropic.BetaOutputConfigEffort
}

func NewClaude(apiKey, model string) *Claude {
	return &Claude{
		Client: anthropic.NewClient(option.WithAPIKey(apiKey)),
		Model:  model,
		Effort: anthropic.BetaOutputConfigEffortMedium,
	}
}

func practiceInstructions(s *Script) string {
	return fmt.Sprintf(`現在是「帶領共修」模式，使用共修腳本《%s》。
- 依腳本順序一次只帶一個步驟，說完就停，等對方回應（例如「好了」「下一步」）再繼續。
- 需要靜默或計時的步驟，明確說出時間長度（例如「接下來靜坐五分鐘」），請對方結束後再回覆。
- 對方中途提問時先簡短回答，再回到目前的步驟。
- 腳本結束時，引導對方簡短分享這次的體會，然後溫和地結束。

<script>
%s
</script>`, s.Title, s.Body)
}

func (c *Claude) Complete(ctx context.Context, t Turn) (Reply, Usage, error) {
	// Persona and knowledge are stable between requests, so they are cached; the script
	// (practice mode only) comes after the breakpoint.
	system := []anthropic.BetaTextBlockParam{
		{Text: t.Persona},
		{Text: t.Knowledge, CacheControl: anthropic.NewBetaCacheControlEphemeralParam()},
	}
	if t.Script != nil {
		system = append(system, anthropic.BetaTextBlockParam{Text: practiceInstructions(t.Script)})
	}
	msgs := make([]anthropic.BetaMessageParam, 0, len(t.Messages))
	for _, m := range t.Messages {
		role := anthropic.BetaMessageParamRoleUser
		if m.Role == "assistant" {
			role = anthropic.BetaMessageParamRoleAssistant
		}
		msgs = append(msgs, anthropic.BetaMessageParam{
			Role:    role,
			Content: []anthropic.BetaContentBlockParamUnion{anthropic.NewBetaTextBlock(m.Content)},
		})
	}
	resp, err := c.Client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
		Model:        c.Model,
		MaxTokens:    4096,
		System:       system,
		Messages:     msgs,
		OutputConfig: anthropic.BetaOutputConfigParam{Effort: c.Effort},
		// On a safety refusal the API re-serves the request on a fallback model.
		Fallbacks: anthropic.BetaFallbacksParamUnion{OfDefault: constant.ValueOf[constant.Default]()},
		Betas:     []anthropic.AnthropicBeta{anthropic.AnthropicBetaServerSideFallback2026_07_01},
	})
	if err != nil {
		return Reply{}, Usage{}, err
	}
	u := Usage{
		Model:     string(resp.Model),
		Input:     resp.Usage.InputTokens,
		Output:    resp.Usage.OutputTokens,
		CacheRead: resp.Usage.CacheReadInputTokens,
		CacheW:    resp.Usage.CacheCreationInputTokens,
	}
	if resp.StopReason == anthropic.BetaStopReasonRefusal {
		return Reply{Text: "這個問題我沒辦法回答。若有其他關於覺行或共修的問題，歡迎再問我。", Refused: true}, u, nil
	}
	var b strings.Builder
	for _, block := range resp.Content {
		if tb, ok := block.AsAny().(anthropic.BetaTextBlock); ok {
			b.WriteString(tb.Text)
		}
	}
	return Reply{Text: strings.TrimSpace(b.String())}, u, nil
}
