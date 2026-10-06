package guide

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"github.com/geodown-ops/bodhi-alliance/server/internal/httpx"
)

// maxSpeechRunes caps one /guide/tts request; the web app sends one sentence at a time.
const maxSpeechRunes = 300

// Speaker turns text into speech (MP3). Tests use a fake.
type Speaker interface {
	Speak(ctx context.Context, text string) ([]byte, error)
}

// AzureSpeech uses Azure AI Speech neural voices; zh-TW-HsiaoChenNeural is a natural
// Taiwanese Mandarin voice. https://learn.microsoft.com/azure/ai-services/speech-service/rest-text-to-speech
type AzureSpeech struct {
	Key, Region, Voice string
	// Rate is an SSML prosody rate such as "-5%"; empty keeps the voice's own pace.
	Rate   string
	Client *http.Client
}

func NewAzureSpeech(key, region, voice, rate string) *AzureSpeech {
	return &AzureSpeech{Key: key, Region: region, Voice: voice, Rate: rate, Client: &http.Client{Timeout: 20 * time.Second}}
}

// SSML builds the request body; text is escaped.
func (a *AzureSpeech) SSML(text string) string {
	var esc bytes.Buffer
	_ = xml.EscapeText(&esc, []byte(text))
	lang := "zh-TW"
	if parts := strings.SplitN(a.Voice, "-", 3); len(parts) == 3 {
		lang = parts[0] + "-" + parts[1]
	}
	body := esc.String()
	if a.Rate != "" {
		body = fmt.Sprintf(`<prosody rate=%q>%s</prosody>`, a.Rate, body)
	}
	return fmt.Sprintf(`<speak version="1.0" xmlns="http://www.w3.org/2001/10/synthesis" xml:lang=%q><voice name=%q>%s</voice></speak>`,
		lang, a.Voice, body)
}

func (a *AzureSpeech) Speak(ctx context.Context, text string) ([]byte, error) {
	url := fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/v1", a.Region)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(a.SSML(text)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", a.Key)
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("X-Microsoft-OutputFormat", "audio-24khz-48kbitrate-mono-mp3")
	req.Header.Set("User-Agent", "sunnylife-guide")
	res, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	audio, err := io.ReadAll(io.LimitReader(res.Body, 8<<20))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("azure tts: %s: %.200s", res.Status, audio)
	}
	return audio, nil
}

// tts reads one sentence aloud with the cloud voice. The web app falls back to the
// browser's own voice when this is off or fails.
func (h *Handler) tts(c *gin.Context) {
	if h.TTS == nil {
		httpx.Error(c, http.StatusNotFound, "沒有設定雲端語音")
		return
	}
	var req struct {
		Text string `json:"text"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "格式錯誤")
		return
	}
	text := strings.TrimSpace(req.Text)
	if text == "" || utf8.RuneCountInString(text) > maxSpeechRunes {
		httpx.Error(c, http.StatusBadRequest, "文字長度不正確")
		return
	}
	u, err := h.Store.Usage(c)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "暫時無法發聲")
		return
	}
	if !h.available(u) {
		httpx.Error(c, http.StatusServiceUnavailable, "線上組長目前休息中")
		return
	}
	audio, err := h.TTS.Speak(c, text)
	if err != nil {
		log.Printf("guide tts: %v", err)
		httpx.Error(c, http.StatusBadGateway, "語音暫時無法使用")
		return
	}
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "audio/mpeg", audio)
}
