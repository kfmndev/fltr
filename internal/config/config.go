package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func LoadBlockRules(path string) (map[string][]string, error) {
	if path == "" {
		return nil, errors.New("block rules path is empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rules map[string][]string
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}

	for rule, terms := range rules {
		if len(terms) == 0 {
			return nil, fmt.Errorf("rule %q has no terms", rule)
		}

		for i, term := range terms {
			term = strings.TrimSpace(term)
			if term == "" {
				return nil, fmt.Errorf("rule %q has an empty term at index %d", rule, i)
			}
			rules[rule][i] = term
		}
	}

	return rules, nil
}

func EnvOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func GetLogLevel() log.Level {
	envLevel := os.Getenv("LOG_LEVEL")
	level, err := log.ParseLevel(envLevel)
	if err != nil {
		return log.InfoLevel
	}
	return level
}
