package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fltr/internal/filter"
)

func TestNewServerDefaults(t *testing.T) {
	allowUpstream := newUpstream(t)
	t.Setenv("FLTR_ALLOW_UPSTREAM", allowUpstream.URL)
	t.Setenv("FLTR_BLOCK_UPSTREAM", "")
	t.Setenv("FLTR_BLOCKED_FILE", writeRules(t, `{"rule":[" term "]}`))
	t.Setenv("FLTR_ADDR", "127.0.0.1:0")
	t.Setenv("FLTR_CASE_SENSITIVE", "invalid")
	t.Setenv("FLTR_MAX_BODY_SIZE", "")

	server, err := newServer()
	if err != nil {
		t.Fatalf("newServer returned error: %v", err)
	}

	handler, ok := server.Handler.(filter.ContentFilter)
	if !ok {
		t.Fatalf("server handler has type %T, want filter.ContentFilter", server.Handler)
	}
	if handler.BlockProxy != nil || handler.CaseSensitive {
		t.Fatal("default handler has unexpected block proxy or case sensitivity")
	}
	if handler.MaxBodySize != 10_000_000 {
		t.Fatalf("MaxBodySize = %d, want %d", handler.MaxBodySize, 10_000_000)
	}
	if server.ReadHeaderTimeout != 5*time.Second || server.ReadTimeout != 15*time.Second || server.WriteTimeout != 30*time.Second || server.IdleTimeout != 60*time.Second {
		t.Fatal("server timeout configuration does not match defaults")
	}
}

func TestNewServerWithBlockUpstreamAndOptions(t *testing.T) {
	allowUpstream := newUpstream(t)
	blockUpstream := newUpstream(t)
	t.Setenv("FLTR_ALLOW_UPSTREAM", allowUpstream.URL)
	t.Setenv("FLTR_BLOCK_UPSTREAM", blockUpstream.URL)
	t.Setenv("FLTR_BLOCKED_FILE", writeRules(t, `{"rule":["term"]}`))
	t.Setenv("FLTR_CASE_SENSITIVE", "true")
	t.Setenv("FLTR_MAX_BODY_SIZE", "5 KB")
	t.Setenv("FLTR_ADDR", "127.0.0.1:12345")

	server, err := newServer()
	if err != nil {
		t.Fatalf("newServer returned error: %v", err)
	}
	handler := server.Handler.(filter.ContentFilter)
	if handler.BlockProxy == nil || !handler.CaseSensitive || handler.MaxBodySize != 5_000 {
		t.Fatal("configured handler does not contain the requested options")
	}
	if server.Addr != "127.0.0.1:12345" {
		t.Fatalf("server address = %q, want 127.0.0.1:12345", server.Addr)
	}
}

func TestNewServerRejectsInvalidConfiguration(t *testing.T) {
	allowUpstream := newUpstream(t)
	rulesPath := writeRules(t, `{"rule":["term"]}`)
	tests := []struct {
		name    string
		allow   string
		block   string
		rules   string
		maxBody string
		wantErr string
	}{
		{name: "missing allow upstream", wantErr: "FLTR_ALLOW_UPSTREAM is required"},
		{name: "invalid allow upstream", allow: "ftp://example.com", wantErr: "invalid FLTR_ALLOW_UPSTREAM"},
		{name: "unreachable allow upstream", allow: closedUpstreamURL(t), wantErr: "FLTR_ALLOW_UPSTREAM is unreachable"},
		{name: "invalid block upstream", allow: allowUpstream.URL, block: "ftp://example.com", rules: rulesPath, wantErr: "invalid FLTR_BLOCK_UPSTREAM"},
		{name: "unreachable block upstream", allow: allowUpstream.URL, block: closedUpstreamURL(t), rules: rulesPath, wantErr: "FLTR_BLOCK_UPSTREAM is unreachable"},
		{name: "missing rules file", allow: allowUpstream.URL, rules: "/missing/rules.json", wantErr: "could not load block rules"},
		{name: "invalid max body size", allow: allowUpstream.URL, rules: rulesPath, maxBody: "not a size", wantErr: "invalid FLTR_MAX_BODY_SIZE"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("FLTR_ALLOW_UPSTREAM", test.allow)
			t.Setenv("FLTR_BLOCK_UPSTREAM", test.block)
			t.Setenv("FLTR_BLOCKED_FILE", test.rules)
			t.Setenv("FLTR_MAX_BODY_SIZE", test.maxBody)

			_, err := newServer()
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("newServer error = %v, want substring %q", err, test.wantErr)
			}
		})
	}
}

func TestRunHandlesServerResults(t *testing.T) {
	allowUpstream := newUpstream(t)
	t.Setenv("FLTR_ALLOW_UPSTREAM", allowUpstream.URL)
	t.Setenv("FLTR_BLOCKED_FILE", writeRules(t, `{"rule":["term"]}`))

	if err := run(func(*http.Server) error { return http.ErrServerClosed }); err != nil {
		t.Fatalf("run returned error for normal server close: %v", err)
	}

	wantErr := errors.New("listen failed")
	err := run(func(*http.Server) error { return wantErr })
	if !strings.Contains(err.Error(), wantErr.Error()) {
		t.Fatalf("run error = %v, want wrapped %v", err, wantErr)
	}
}

func newUpstream(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	t.Cleanup(server.Close)
	return server
}

func closedUpstreamURL(t *testing.T) string {
	t.Helper()
	server := newUpstream(t)
	url := server.URL
	server.Close()
	return url
}

func writeRules(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "rules.json")
	if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}
