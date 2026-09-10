package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNewProxyURLValidation(t *testing.T) {
	tests := []struct {
		name     string
		upstream string
		wantErr  bool
	}{
		{
			name:     "http URL",
			upstream: "http://example.com",
			wantErr:  false,
		},
		{
			name:     "http URL with port",
			upstream: "http://localhost:9000",
			wantErr:  false,
		},
		{
			name:     "https URL",
			upstream: "https://example.com",
			wantErr:  false,
		},
		{
			name:     "missing scheme",
			upstream: "localhost:9000",
			wantErr:  true,
		},
		{
			name:     "unsupported scheme",
			upstream: "ftp://example.com",
			wantErr:  true,
		},
		{
			name:     "missing hostname",
			upstream: "http://",
			wantErr:  true,
		},
		{
			name:     "empty hostname",
			upstream: "http://:9000",
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewProxy(test.upstream)

			if (err != nil) != test.wantErr {
				t.Fatalf("NewProxy(%q) error = %v, wantErr %v", test.upstream, err, test.wantErr)
			}
		})
	}
}

func TestNewProxyForwardsToConfiguredURLAndPreservesQuery(t *testing.T) {
	var received struct {
		method string
		path   string
		query  string
		header string
		body   string
	}
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received.method = r.Method
		received.path = r.URL.Path
		received.query = r.URL.RawQuery
		received.header = r.Header.Get("X-Test")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read upstream body: %v", err)
			return
		}
		received.body = string(body)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer upstream.Close()

	proxyHandler, err := NewProxy(upstream.URL + "/configured?fixed=1")
	if err != nil {
		t.Fatalf("NewProxy returned error: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "http://client.test/incoming?incoming=1&repeat=two&repeat=three", strings.NewReader("payload"))
	request.Header.Set("X-Test", "present")
	response := httptest.NewRecorder()
	proxyHandler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("proxy status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if received.method != http.MethodPost {
		t.Errorf("upstream method = %q, want %q", received.method, http.MethodPost)
	}
	if received.path != "/configured" {
		t.Errorf("upstream path = %q, want %q", received.path, "/configured")
	}
	gotQuery, err := url.ParseQuery(received.query)
	if err != nil {
		t.Fatalf("parse upstream query %q: %v", received.query, err)
	}
	wantQuery := url.Values{
		"fixed":    {"1"},
		"incoming": {"1"},
		"repeat":   {"two", "three"},
	}
	if gotQuery.Encode() != wantQuery.Encode() {
		t.Errorf("upstream query = %q, want %q", gotQuery.Encode(), wantQuery.Encode())
	}
	if received.header != "present" {
		t.Errorf("upstream X-Test header = %q, want %q", received.header, "present")
	}
	if received.body != "payload" {
		t.Errorf("upstream body = %q, want %q", received.body, "payload")
	}
}

func TestCheckUpstream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead {
			t.Errorf("request method = %q, want %q", r.Method, http.MethodHead)
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	if err := CheckUpstream(server.URL); err != nil {
		t.Fatalf("CheckUpstream returned error for reachable server: %v", err)
	}

	server.Close()
	if err := CheckUpstream(server.URL); err == nil {
		t.Fatal("CheckUpstream returned nil for unreachable server")
	}
}

func TestCheckUpstreamNetworkError(t *testing.T) {
	err := CheckUpstream("http://invalid-hostname-that-does-not-exist-12345.test:9999")
	if err == nil {
		t.Fatal("CheckUpstream should return error for unreachable host")
	}
}
