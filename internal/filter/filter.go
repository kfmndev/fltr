package filter

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ContentFilter struct {
	BlockRules    map[string][]string
	CaseSensitive bool
	AllowProxy    http.Handler
	BlockProxy    http.Handler
	MaxBodySize   int64
}

func (f ContentFilter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, f.MaxBodySize)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		if _, ok := err.(*http.MaxBytesError); ok {
			http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	var searchKeys = []string{"Title", "Message"}

	searchText := string(body)

	for _, key := range searchKeys {
		searchText += " " + r.Header.Get(key)
		searchText += " " + r.Header.Get(strings.ToLower(key))
		searchText += " " + r.URL.Query().Get(key)
		searchText += " " + r.URL.Query().Get(strings.ToLower(key))
	}

	r.Body = io.NopCloser(bytes.NewReader(body))
	r.ContentLength = int64(len(body))

	if MatchRules(searchText, f.BlockRules, f.CaseSensitive) {
		if f.BlockProxy != nil {
			log.Debug("Request blocked, proxied")
			f.BlockProxy.ServeHTTP(w, r)
		} else {
			log.Debug("Request blocked, discarded")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("Request blocked, discarded")); err != nil {
				log.Debugf("Failed to write response: %v", err)
			}
		}
		return
	}

	log.Debug("Request allowed")
	f.AllowProxy.ServeHTTP(w, r)
}

func MatchRules(searchText string, blockRules map[string][]string, caseSensitive bool) bool {
	for key, terms := range blockRules {
		if AllTermsFound(searchText, terms, caseSensitive) {
			log.Debugf("Matched terms %#v (rule %q)", terms, key)
			return true
		}
	}
	return false
}

func AllTermsFound(searchText string, terms []string, caseSensitive bool) bool {
	if len(terms) == 0 {
		return false
	}

	if !caseSensitive {
		searchText = strings.ToLower(searchText)
		lowerTerms := make([]string, len(terms))
		for i, term := range terms {
			lowerTerms[i] = strings.ToLower(term)
		}
		terms = lowerTerms
	}

	for _, term := range terms {
		if !strings.Contains(searchText, term) {
			log.Debugf("Term %q not found in search corpus", term)
			return false
		}
	}
	log.Debugf("All terms %v found in search corpus", terms)
	return true
}
