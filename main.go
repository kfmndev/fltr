package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"fltr/internal/config"
	"fltr/internal/filter"
	"fltr/internal/proxy"

	"github.com/docker/go-units"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Configure logging
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	log.SetLevel(config.GetLogLevel())

	// Required allow upstream
	allowUpstream := os.Getenv("FLTR_ALLOW_UPSTREAM")
	if allowUpstream == "" {
		log.Fatal("FLTR_ALLOW_UPSTREAM is required")
	}

	allowProxy, err := proxy.NewProxy(allowUpstream)
	if err != nil {
		log.Fatalf("Invalid FLTR_ALLOW_UPSTREAM: %v", err)
	}
	if err := proxy.CheckUpstream(allowUpstream); err != nil {
		log.Fatalf("FLTR_ALLOW_UPSTREAM is unreachable: %v", err)
	}

	// Optional block upstream
	blockUpstream := os.Getenv("FLTR_BLOCK_UPSTREAM")
	var blockProxy http.Handler
	if blockUpstream == "" {
		log.Info("FLTR_BLOCK_UPSTREAM is not set")
	} else {
		blockProxy, err = proxy.NewProxy(blockUpstream)
		if err != nil {
			log.Fatalf("Invalid FLTR_BLOCK_UPSTREAM: %v", err)
		}
		if err := proxy.CheckUpstream(blockUpstream); err != nil {
			log.Fatalf("FLTR_BLOCK_UPSTREAM is unreachable: %v", err)
		}
	}

	// Determine case sensitivity for block rule matching
	caseSensitive, err := strconv.ParseBool(config.EnvOr("FLTR_CASE_SENSITIVE", "false"))
	if err != nil {
		// Default to case insensitivity
		caseSensitive = false
	}

	if caseSensitive {
		log.Infof("Block rule matching is case sensitive")
	} else {
		log.Infof("Block rule matching is case insensitive (default)")
	}

	// Load block rules from the specified file
	blockRules, err := config.LoadBlockRules(config.EnvOr("FLTR_BLOCK_RULES_FILE", "block_rules.json"))
	if err != nil {
		log.Fatalf("Could not load block rules: %v", err)
	}

	log.Infof("Loaded %d block rules", len(blockRules))
	log.Debugf("Block rules: %#v", blockRules)

	// Maximum request body size
	maxBodySize, err := units.FromHumanSize(config.EnvOr("FLTR_MAX_BODY_SIZE", "10 MB"))
	if err != nil {
		log.Fatalf("Invalid FLTR_MAX_BODY_SIZE: %v", err)
	}

	log.Infof("Request bodies >%d bytes (%s) are rejected", maxBodySize, units.HumanSize(float64(maxBodySize)))

	// Initialize and start the HTTP server
	server := &http.Server{
		Addr:              config.EnvOr("FLTR_ADDR", ":8080"),
		Handler:           filter.ContentFilter{BlockRules: blockRules, CaseSensitive: caseSensitive, AllowProxy: allowProxy, BlockProxy: blockProxy, MaxBodySize: maxBodySize},
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	if blockProxy == nil {
		log.Infof("Allowed requests -> %s", allowUpstream)
		log.Info("Blocked request -> DISCARDED")
	} else {
		log.Infof("Allowed requests -> %s", allowUpstream)
		log.Infof("Blocked requests -> %s", blockUpstream)
	}

	log.Infof("Listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server stopped: %v", err)
	}
}
