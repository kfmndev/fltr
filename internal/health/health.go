package health

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

const Path = "/healthz"

type Handler struct {
	Inner http.Handler
}

func New(inner http.Handler) Handler {
	return Handler{Inner: inner}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != Path {
		h.Inner.ServeHTTP(w, r)
		return
	}

	log.Debug("Liveness probe")

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)

	if r.Method == http.MethodHead {
		return
	}

	if _, err := w.Write([]byte("ok")); err != nil {
		log.Debugf("Failed to write health response: %v", err)
	}
}
