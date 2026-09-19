package health

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type spyHandler struct {
	called bool
	method string
	path   string
	body   string
}

func (s *spyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.called = true
	s.method = r.Method
	s.path = r.URL.Path
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read body", http.StatusBadRequest)
		return
	}
	s.body = string(body)
	w.WriteHeader(http.StatusTeapot)
}

func TestLivenessEndpoint(t *testing.T) {
	spy := &spyHandler{}
	handler := New(spy)

	tests := []struct {
		name     string
		method   string
		path     string
		wantBody string
	}{
		{name: "GET", method: http.MethodGet, path: Path, wantBody: "ok"},
		{name: "HEAD", method: http.MethodHead, path: Path},
		{name: "query ignored", method: http.MethodGet, path: Path + "?term=blocked", wantBody: "ok"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			*spy = spyHandler{}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(test.method, test.path, strings.NewReader("term"))

			handler.ServeHTTP(rec, req)
			assertLivenessResponse(t, rec, test.wantBody)

			if spy.called {
				t.Fatal("inner handler was called for the liveness endpoint")
			}
		})
	}
}

func assertLivenessResponse(t *testing.T, rec *httptest.ResponseRecorder, wantBody string) {
	t.Helper()

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != wantBody {
		t.Fatalf("body = %q, want %q", got, wantBody)
	}
	if got := rec.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Fatalf("Content-Type = %q", got)
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q", got)
	}
}

func TestLivenessEndpointRejectsOtherMethods(t *testing.T) {
	spy := &spyHandler{}
	handler := New(spy)

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions} {
		t.Run(method, func(t *testing.T) {
			*spy = spyHandler{}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, Path, nil)

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
			}
			if got := rec.Header().Get("Allow"); got != "GET, HEAD" {
				t.Fatalf("Allow = %q, want %q", got, "GET, HEAD")
			}
			if spy.called {
				t.Fatal("inner handler was called for the liveness endpoint")
			}
		})
	}
}

func TestHandlerDelegatesNonHealthPaths(t *testing.T) {
	tests := []string{"/", "/healthz/", "/healthzfoo", "/health", "/other"}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			spy := &spyHandler{}
			handler := New(spy)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, path, strings.NewReader("payload"))

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusTeapot {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusTeapot)
			}
			if !spy.called {
				t.Fatal("inner handler was not called")
			}
			if spy.method != http.MethodPost || spy.path != path || spy.body != "payload" {
				t.Fatalf("inner received method=%q path=%q body=%q", spy.method, spy.path, spy.body)
			}
		})
	}
}
