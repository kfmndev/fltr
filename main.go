package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"fltr/internal/config"
	"fltr/internal/filter"
	"fltr/internal/health"
	"fltr/internal/proxy"

	"github.com/docker/go-units"
	log "github.com/sirupsen/logrus"
)

// BEGIN-NOSCAN
func main() {
	// Configure logging
	config.SetupLogging()

	if err := run(func(server *http.Server) error {
		return server.ListenAndServe()
	}); err != nil {
		log.Fatal(err)
	}
}

func run(serve func(*http.Server) error) error {
	server, err := newServer()
	if err != nil {
		return err
	}

	if err := serve(server); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server stopped: %w", err)
	}
	return nil
}

// END-NOSCAN

func newServer() (*http.Server, error) {
	// Required allow upstream
	allowUpstream := os.Getenv("FLTR_ALLOW_UPSTREAM")
	if allowUpstream == "" {
		return nil, fmt.Errorf("FLTR_ALLOW_UPSTREAM is required")
	}

	allowProxy, err := proxy.NewProxy(allowUpstream)
	if err != nil {
		return nil, fmt.Errorf("invalid FLTR_ALLOW_UPSTREAM: %w", err)
	}
	if err := proxy.CheckUpstream(allowUpstream); err != nil {
		return nil, fmt.Errorf("FLTR_ALLOW_UPSTREAM is unreachable: %w", err)
	}

	// Optional block upstream
	blockUpstream := os.Getenv("FLTR_BLOCK_UPSTREAM")
	var blockProxy http.Handler
	if blockUpstream == "" {
		log.Info("FLTR_BLOCK_UPSTREAM is not set")
	} else {
		blockProxy, err = proxy.NewProxy(blockUpstream)
		if err != nil {
			return nil, fmt.Errorf("invalid FLTR_BLOCK_UPSTREAM: %w", err)
		}
		if err := proxy.CheckUpstream(blockUpstream); err != nil {
			return nil, fmt.Errorf("FLTR_BLOCK_UPSTREAM is unreachable: %w", err)
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
		return nil, fmt.Errorf("could not load block rules: %w", err)
	}

	log.Infof("Loaded %d block rules", len(blockRules))
	log.Debugf("Block rules: %#v", blockRules)

	// Maximum request body size
	maxBodySize, err := units.FromHumanSize(config.EnvOr("FLTR_MAX_BODY_SIZE", "10 MB"))
	if err != nil {
		return nil, fmt.Errorf("invalid FLTR_MAX_BODY_SIZE: %w", err)
	}

	log.Infof("Request bodies >%d bytes (%s) are rejected", maxBodySize, units.HumanSize(float64(maxBodySize)))

	// Initialize and start the HTTP server
	server := &http.Server{
		Addr:              config.EnvOr("FLTR_ADDR", ":8080"),
		Handler:           health.New(filter.ContentFilter{BlockRules: blockRules, CaseSensitive: caseSensitive, AllowProxy: allowProxy, BlockProxy: blockProxy, MaxBodySize: maxBodySize}),
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

	return server, nil
}
