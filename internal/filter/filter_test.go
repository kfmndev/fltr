package filter

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContentFilterRoutesAllowedAndBlockedRequests(t *testing.T) {
	var allowed, blocked string
	filter := ContentFilter{
		BlockRules:  map[string][]string{"credentials": {"password", "secret"}},
		MaxBodySize: 10 * 1024 * 1024,
		AllowProxy: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			allowed = string(body)
			w.WriteHeader(http.StatusAccepted)
		}),
		BlockProxy: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			blocked = string(body)
			w.WriteHeader(http.StatusCreated)
		}),
	}

	allowedRequest := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("ordinary content"))
	allowedResponse := httptest.NewRecorder()
	filter.ServeHTTP(allowedResponse, allowedRequest)
	if allowedResponse.Code != http.StatusAccepted {
		t.Fatalf("allowed status = %d, want %d", allowedResponse.Code, http.StatusAccepted)
	}
	if allowed != "ordinary content" {
		t.Errorf("allowed request = %q, want %q", allowed, "ordinary content")
	}

	blockedRequest := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("contains SECRET"))
	blockedRequest.Header.Set("Title", "PASSWORD")
	blockedResponse := httptest.NewRecorder()
	filter.ServeHTTP(blockedResponse, blockedRequest)
	if blockedResponse.Code != http.StatusCreated {
		t.Fatalf("blocked status = %d, want %d", blockedResponse.Code, http.StatusCreated)
	}
	if blocked != "contains SECRET" {
		t.Errorf("blocked request = %q, want %q", blocked, "contains SECRET")
	}
}

func TestContentFilterMatchesQueryParameters(t *testing.T) {
	var allowed, blocked string
	filter := ContentFilter{
		BlockRules:  map[string][]string{"credentials": {"password", "secret"}},
		MaxBodySize: 10 * 1024 * 1024,
		AllowProxy: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			allowed = string(body)
			w.WriteHeader(http.StatusAccepted)
		}),
		BlockProxy: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			blocked = string(body)
			w.WriteHeader(http.StatusCreated)
		}),
	}

	allowedRequest := httptest.NewRequest(http.MethodPost, "/?title=ordinary&message=content", strings.NewReader(""))
	allowedResponse := httptest.NewRecorder()
	filter.ServeHTTP(allowedResponse, allowedRequest)
	if allowedResponse.Code != http.StatusAccepted {
		t.Fatalf("allowed status = %d, want %d", allowedResponse.Code, http.StatusAccepted)
	}

	blockedRequest := httptest.NewRequest(http.MethodPost, "/?title=secret&message=my+password", strings.NewReader(""))
	blockedResponse := httptest.NewRecorder()
	filter.ServeHTTP(blockedResponse, blockedRequest)
	if blockedResponse.Code != http.StatusCreated {
		t.Fatalf("blocked status = %d, want %d", blockedResponse.Code, http.StatusCreated)
	}

	pctEncodedRequest := httptest.NewRequest(http.MethodPost, "/?title=secret&message=my%20password", strings.NewReader(""))
	pctEncodedResponse := httptest.NewRecorder()
	filter.ServeHTTP(pctEncodedResponse, pctEncodedRequest)
	if pctEncodedResponse.Code != http.StatusCreated {
		t.Fatalf("percent-encoded blocked status = %d, want %d", pctEncodedResponse.Code, http.StatusCreated)
	}

	if allowed != "" || blocked != "" {
		t.Errorf("body proxied = %q/%q, want empty", allowed, blocked)
	}
}

func TestContentFilterDiscardsBlockedRequestWithoutBlockProxy(t *testing.T) {
	filter := ContentFilter{
		BlockRules:  map[string][]string{"secret": {"secret"}},
		MaxBodySize: 10 * 1024 * 1024,
		AllowProxy: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			t.Fatal("allow proxy was called")
		}),
	}

	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("secret"))
	response := httptest.NewRecorder()
	filter.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("discarded status = %d, want %d", response.Code, http.StatusOK)
	}
	if response.Body.String() != "Request blocked, discarded" {
		t.Fatalf("discarded response = %q", response.Body.String())
	}
}

func TestContentFilterRestoresRequestBodyAndContentLength(t *testing.T) {
	var receivedBody string
	var receivedLength int64
	filter := ContentFilter{
		BlockRules:  map[string][]string{"secret": {"secret"}},
		MaxBodySize: 10 * 1024 * 1024,
		AllowProxy: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("read proxied body: %v", err)
			}
			receivedBody = string(body)
			receivedLength = r.ContentLength
			w.WriteHeader(http.StatusNoContent)
		}),
	}

	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("ordinary content"))
	response := httptest.NewRecorder()
	filter.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("proxied status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if receivedBody != "ordinary content" {
		t.Errorf("proxied body = %q, want %q", receivedBody, "ordinary content")
	}
	if receivedLength != int64(len("ordinary content")) {
		t.Errorf("proxied content length = %d, want %d", receivedLength, len("ordinary content"))
	}
}

func TestAllTermsFound(t *testing.T) {
	tests := []struct {
		name          string
		searchCorpus  string
		terms         []string
		caseSensitive bool
		want          bool
	}{
		{
			name:          "all terms present case insensitive",
			searchCorpus:  "an apple request is banana",
			terms:         []string{"APPLE", "BANANA"},
			caseSensitive: false,
			want:          true,
		},
		{
			name:          "all terms present wrong case",
			searchCorpus:  "an apple request is banana",
			terms:         []string{"APPLE", "BANANA"},
			caseSensitive: true,
			want:          false,
		},
		{
			name:          "one term missing",
			searchCorpus:  "apple request",
			terms:         []string{"apple", "banana"},
			caseSensitive: false,
			want:          false,
		},
		{
			name:          "empty terms",
			searchCorpus:  "anything",
			terms:         nil,
			caseSensitive: false,
			want:          false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := AllTermsFound(test.searchCorpus, test.terms, test.caseSensitive); got != test.want {
				t.Fatalf("AllTermsFound(%q, %#v) = %v, want %v", test.searchCorpus, test.terms, got, test.want)
			}
		})
	}
}

func TestMatchRules(t *testing.T) {
	rules := map[string][]string{
		"rule-one": {"apple", "banana", "cherry"},
		"rule-two": {"peach", "banana", "cherry"},
	}

	tests := []struct {
		name          string
		searchCorpus  string
		caseSensitive bool
		want          bool
	}{
		{
			name:          "matches one complete rule",
			searchCorpus:  "PEACH banana CHERRY",
			caseSensitive: false,
			want:          true,
		},
		{
			name:          "matches one complete rule",
			searchCorpus:  "peach banana cherry",
			caseSensitive: true,
			want:          true,
		},
		{
			name:         "only matches incomplete rule",
			searchCorpus: "peach banana fruit",
			want:         false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := MatchRules(test.searchCorpus, rules, test.caseSensitive); got != test.want {
				t.Fatalf("MatchRules(%q, %#v) = %v, want %v", test.searchCorpus, rules, got, test.want)
			}
		})
	}
}

func TestAllTermsFoundDoesNotMutateTerms(t *testing.T) {
	terms := []string{"APPLE", "BANANA"}
	termsCopy := make([]string, len(terms))
	copy(termsCopy, terms)

	_ = AllTermsFound("apple banana cherry", terms, false)

	if terms[0] != termsCopy[0] || terms[1] != termsCopy[1] {
		t.Errorf("terms were mutated: %#v -> %#v", termsCopy, terms)
	}
}
