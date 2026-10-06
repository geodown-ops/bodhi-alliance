// Package config reads service settings from environment variables.
package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL string
	Addr        string
	// AllowedOrigins lists front-end origins allowed to call the API (CORS).
	AllowedOrigins []string
	// TrustedProxies lists proxy IPs/CIDRs whose X-Forwarded-For is believed (for rate limits).
	TrustedProxies []string

	// First alliance admin, created on startup if no user has that email yet.
	BootstrapAdminEmail    string
	BootstrapAdminPassword string

	// AI 組長
	AnthropicAPIKey string
	GuideModel      string
	GuideName       string

	// Sunny 的雲端語音（Azure AI Speech）；沒有金鑰時網站改用瀏覽器內建語音
	AzureSpeechKey    string
	AzureSpeechRegion string
	TTSVoice          string
	TTSRate           string
}

func Load(defaultAddr string) Config {
	return Config{
		DatabaseURL:            env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/bodhi?sslmode=disable"),
		Addr:                   env("ADDR", portAddr(defaultAddr)),
		AllowedOrigins:         split(env("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:5174")),
		TrustedProxies:         split(os.Getenv("TRUSTED_PROXIES")),
		BootstrapAdminEmail:    os.Getenv("BOOTSTRAP_ADMIN_EMAIL"),
		BootstrapAdminPassword: os.Getenv("BOOTSTRAP_ADMIN_PASSWORD"),
		AnthropicAPIKey:        os.Getenv("ANTHROPIC_API_KEY"),
		GuideModel:             env("GUIDE_MODEL", "claude-opus-5-5"),
		GuideName:              env("GUIDE_NAME", "Sunny"),
		AzureSpeechKey:         strings.TrimSpace(os.Getenv("AZURE_SPEECH_KEY")),
		AzureSpeechRegion:      env("AZURE_SPEECH_REGION", "eastasia"),
		TTSVoice:               env("TTS_VOICE", "zh-TW-HsiaoChenNeural"),
		TTSRate:                strings.TrimSpace(os.Getenv("TTS_RATE")),
	}
}

func env(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func split(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// portAddr honours PORT (set by Railway and most hosts) when ADDR is not given.
func portAddr(fallback string) string {
	if p := strings.TrimSpace(os.Getenv("PORT")); p != "" {
		return ":" + p
	}
	return fallback
}
