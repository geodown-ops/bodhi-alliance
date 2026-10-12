package guide

import (
	"context"
	"encoding/json"
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
- 用溫和、口語的繁體中文回答，像一位親切的小組組長，不說教。回答越精簡越好，並邀請對方一起聊、一起練習。

界線：
- 只根據「知識庫」裡的內容回答佛法、覺行與菩提幣的問題。知識庫沒有寫的，就老實說目前沒有這方面的資料，並建議到官網「覺行小組」頁面聯絡真人組長。
- 對方問起資料來源時，再用《文件標題》說明出處。
- 你看不到也改不了任何人的菩提幣餘額或券，問到時請對方到官網「我的錢包」查看。
- 菩提幣不販售、不提領、不可兌現，也不是投資；不要做任何價格或收益的預測。
- 不做醫療或心理診斷。若對方提到想傷害自己或他人，溫和地請對方立刻尋求協助：台灣安心專線 1925、生命線 1995，緊急狀況撥 119／110。`, name)
}

// replyStyle is how every answer is shaped. It lives in code, apart from the persona
// edited in the admin, so it applies whatever the persona says.
const replyStyle = `回答方式（每一則回答都要遵守）：
- 越精簡越好：一般回答一到三句、八十字以內；對方請你詳細說明時才多說，也盡量不超過一百五十字。
- 先講重點，不鋪陳、不重複對方的問題、不列一長串。
- 回答會被唸出來，所以用說話的口氣：不用標題、條列、粗體、表格或表情符號，也不必在句末標出處。
- 結尾用一句簡短、自然的問句邀請對方互動，例如問對方的經驗或感受、邀請現在一起試一下，或問想多了解哪一部分。每次換個說法，不要每次都問一樣的話。`

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
	// Issues, when set, lets Sunny record usage problems for the admin (public chat only).
	Issues Reporter
	// Visitor says whether the person chatting is a signed-in member.
	Visitor string
}

// Reporter records usage problems Sunny can't solve; *IssueReporter in production.
type Reporter interface {
	Record(ctx context.Context, summary, email string) error
	SetEmail(ctx context.Context, email string) error
}

// LLM is the model behind the guide; tests use a fake. onText receives the answer as
// it is generated and may be nil. On a refusal, text already sent should be discarded
// in favor of the returned Reply.
type LLM interface {
	Complete(ctx context.Context, t Turn, onText func(string)) (Reply, Usage, error)
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

// issueRules tells Sunny when and how to record a usage problem; it is sent only with the tools.
const issueRules = `網站使用問題（官網怎麼用、加入會員、登入、個人頁、錢包、共修活動等）：
- 先依知識庫「官網使用說明」回答。
- 知識庫查不到答案、你不確定，或對方說遇到錯誤、按了沒反應、登不進去、資料不對等狀況時，不要猜，用 record_site_issue 把問題記下來交給管理員，再告訴對方已經記下、管理員會回覆。
- 記錄前不必再三確認；摘要寫清楚對方想做什麼、在哪一頁、發生什麼事，用對方的話。
- 同一個問題只記一次。對方補充細節時不要重複記錄，除非是另一個新問題。
- 對方已登入時，告訴對方管理員的回覆會出現在個人頁的「我的使用問題」。
- 對方沒有登入時，問對方願不願意留 email 讓管理員聯絡（可以不留）；對方給了 email，用 set_issue_email 補上。不要追問，也不要要求電話或其他個資。
- 忘記密碼目前要由管理員協助，也要記下來。`

var issueTools = []anthropic.BetaToolUnionParam{
	{OfTool: &anthropic.BetaToolParam{
		Name:        "record_site_issue",
		Description: anthropic.String("把對方遇到的網站使用問題記到後台「使用問題」，讓超級管理員回答。只在知識庫答不了，或對方回報錯誤、異常時使用。"),
		InputSchema: anthropic.BetaToolInputSchemaParam{
			Properties: map[string]any{
				"summary":       map[string]any{"type": "string", "description": "問題摘要（繁體中文，三百字以內）：對方想做什麼、在哪一頁、發生了什麼。"},
				"contact_email": map[string]any{"type": "string", "description": "對方主動提供的 email；沒有就不要填。"},
			},
			Required: []string{"summary"},
		},
	}},
	{OfTool: &anthropic.BetaToolParam{
		Name:        "set_issue_email",
		Description: anthropic.String("對方在記下問題之後才留 email 時，把 email 補到這次對話記下的問題上。"),
		InputSchema: anthropic.BetaToolInputSchemaParam{
			Properties: map[string]any{
				"email": map[string]any{"type": "string", "description": "對方留下的 email。"},
			},
			Required: []string{"email"},
		},
	}},
}

// maxToolRounds bounds how many times one answer may call tools.
const maxToolRounds = 3

// runTool executes one tool call and returns the text sent back to the model.
func runTool(ctx context.Context, r Reporter, name string, input any) (string, bool) {
	raw, _ := json.Marshal(input)
	var in struct {
		Summary      string `json:"summary"`
		ContactEmail string `json:"contact_email"`
		Email        string `json:"email"`
	}
	_ = json.Unmarshal(raw, &in)
	var err error
	switch name {
	case "record_site_issue":
		err = r.Record(ctx, in.Summary, in.ContactEmail)
	case "set_issue_email":
		err = r.SetEmail(ctx, in.Email)
	default:
		err = errors.New("沒有這個工具")
	}
	if err != nil {
		return "沒有成功：" + err.Error(), true
	}
	return "已記下。", false
}

func (c *Claude) Complete(ctx context.Context, t Turn, onText func(string)) (Reply, Usage, error) {
	// Persona and knowledge are stable between requests, so they are cached; the script
	// (practice mode only) comes after the breakpoint.
	system := []anthropic.BetaTextBlockParam{
		{Text: t.Persona},
		{Text: replyStyle},
		{Text: t.Knowledge, CacheControl: anthropic.NewBetaCacheControlEphemeralParam()},
	}
	if t.Script != nil {
		system = append(system, anthropic.BetaTextBlockParam{Text: practiceInstructions(t.Script)})
	}
	var tools []anthropic.BetaToolUnionParam
	if t.Issues != nil {
		system = append(system, anthropic.BetaTextBlockParam{Text: issueRules})
		tools = issueTools
	}
	if t.Visitor != "" {
		system = append(system, anthropic.BetaTextBlockParam{Text: t.Visitor})
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
	var u Usage
	var b strings.Builder
	for round := 0; ; round++ {
		params := anthropic.BetaMessageNewParams{
			Model:        c.Model,
			MaxTokens:    4096,
			System:       system,
			Messages:     msgs,
			OutputConfig: anthropic.BetaOutputConfigParam{Effort: c.Effort},
			// On a safety refusal the API re-serves the request on a fallback model, on the
			// same stream; text already streamed stays valid and the answer continues.
			Fallbacks: anthropic.BetaFallbacksParamUnion{OfDefault: constant.ValueOf[constant.Default]()},
			Betas:     []anthropic.AnthropicBeta{anthropic.AnthropicBetaServerSideFallback2026_07_01},
		}
		if len(tools) > 0 && round < maxToolRounds {
			params.Tools = tools
		}
		resp, err := c.stream(ctx, params, onText)
		if err != nil {
			return Reply{}, u, err
		}
		u.Model = string(resp.Model)
		u.Input += resp.Usage.InputTokens
		u.Output += resp.Usage.OutputTokens
		u.CacheRead += resp.Usage.CacheReadInputTokens
		u.CacheW += resp.Usage.CacheCreationInputTokens
		if resp.StopReason == anthropic.BetaStopReasonRefusal {
			return Reply{Text: "這個問題我沒辦法回答。若有其他關於覺行或共修的問題，歡迎再問我。", Refused: true}, u, nil
		}
		var results []anthropic.BetaContentBlockParamUnion
		for _, block := range resp.Content {
			switch v := block.AsAny().(type) {
			case anthropic.BetaTextBlock:
				b.WriteString(v.Text)
			case anthropic.BetaToolUseBlock:
				text, isErr := runTool(ctx, t.Issues, v.Name, v.Input)
				results = append(results, anthropic.NewBetaToolResultBlock(v.ID, text, isErr))
			}
		}
		if resp.StopReason != anthropic.BetaStopReasonToolUse || len(results) == 0 || t.Issues == nil {
			break
		}
		msgs = append(msgs, resp.ToParam(), anthropic.BetaMessageParam{
			Role:    anthropic.BetaMessageParamRoleUser,
			Content: results,
		})
	}
	return Reply{Text: strings.TrimSpace(b.String())}, u, nil
}

// stream runs one request, passing text to onText as it arrives.
func (c *Claude) stream(ctx context.Context, params anthropic.BetaMessageNewParams, onText func(string)) (anthropic.BetaMessage, error) {
	stream := c.Client.Beta.Messages.NewStreaming(ctx, params)
	defer stream.Close()
	var resp anthropic.BetaMessage
	for stream.Next() {
		ev := stream.Current()
		if err := resp.Accumulate(ev); err != nil {
			return resp, err
		}
		if d, ok := ev.AsAny().(anthropic.BetaRawContentBlockDeltaEvent); ok && onText != nil {
			if td, ok := d.Delta.AsAny().(anthropic.BetaTextDelta); ok {
				onText(td.Text)
			}
		}
	}
	return resp, stream.Err()
}
