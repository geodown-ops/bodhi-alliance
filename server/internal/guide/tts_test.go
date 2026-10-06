package guide

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAzureSpeech(t *testing.T) {
	var gotBody, gotKey, gotFormat string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody, gotKey, gotFormat = string(b), r.Header.Get("Ocp-Apim-Subscription-Key"), r.Header.Get("X-Microsoft-OutputFormat")
		w.Write([]byte("mp3"))
	}))
	defer srv.Close()
	a := NewAzureSpeech("k", "eastasia", "zh-TW-HsiaoChenNeural", "-5%")
	a.Client = srv.Client()
	a.Client.Transport = rewrite{srv.URL}
	audio, err := a.Speak(context.Background(), "你好<&>")
	if err != nil || string(audio) != "mp3" {
		t.Fatalf("Speak = %q, %v", audio, err)
	}
	if gotKey != "k" || !strings.Contains(gotFormat, "mp3") {
		t.Errorf("headers key=%q format=%q", gotKey, gotFormat)
	}
	for _, want := range []string{`xml:lang="zh-TW"`, `name="zh-TW-HsiaoChenNeural"`, `rate="-5%"`, "你好&lt;&amp;&gt;"} {
		if !strings.Contains(gotBody, want) {
			t.Errorf("ssml %q missing %q", gotBody, want)
		}
	}
}

// rewrite sends every request to the test server.
type rewrite struct{ base string }

func (r rewrite) RoundTrip(req *http.Request) (*http.Response, error) {
	u := *req.URL
	u.Scheme, u.Host = "http", strings.TrimPrefix(r.base, "http://")
	req.URL = &u
	return http.DefaultTransport.RoundTrip(req)
}
