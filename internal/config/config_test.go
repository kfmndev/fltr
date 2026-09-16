package config

import (
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestLoadBlockRulesTrimsTerms(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rules.json")
	if err := os.WriteFile(path, []byte(`{"rule":[" Foo ","BAR"]}`), 0600); err != nil {
		t.Fatal(err)
	}

	rules, err := LoadBlockRules(path)
	if err != nil {
		t.Fatalf("LoadBlockRules returned error: %v", err)
	}

	want := []string{"Foo", "BAR"}
	if got := rules["rule"]; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("trimmed terms = %#v, want %#v", got, want)
	}
}

func TestLoadBlockRulesRejectsEmptyTerms(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "blank term",
			json: `{"rule":["   "]}`,
		},
		{
			name: "empty rule",
			json: `{"rule":[]}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "rules.json")
			if err := os.WriteFile(path, []byte(test.json), 0600); err != nil {
				t.Fatal(err)
			}

			if _, err := LoadBlockRules(path); err == nil {
				t.Fatalf("LoadBlockRules(%q) returned nil error", test.json)
			}
		})
	}
}

func TestLoadBlockRulesMalformedJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rules.json")
	if err := os.WriteFile(path, []byte(`{invalid json}`), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := LoadBlockRules(path)
	if err == nil {
		t.Fatal("LoadBlockRules should return error for malformed JSON")
	}
}

func TestLoadBlockRulesMissingFile(t *testing.T) {
	_, err := LoadBlockRules("/nonexistent/path/to/rules.json")
	if err == nil {
		t.Fatal("LoadBlockRules should return error for missing file")
	}
}

func TestLoadBlockRulesEmptyPath(t *testing.T) {
	_, err := LoadBlockRules("")
	if err == nil {
		t.Fatal("LoadBlockRules should return error for empty path")
	}
}

func TestSetupLoggingDefaultsToText(t *testing.T) {
	t.Setenv("LOG_FORMAT", "")
	t.Setenv("LOG_LEVEL", "")

	SetupLogging()

	formatter, ok := log.StandardLogger().Formatter.(*log.TextFormatter)
	if !ok {
		t.Fatalf("formatter = %T, want *log.TextFormatter", log.StandardLogger().Formatter)
	}
	if !formatter.FullTimestamp {
		t.Fatal("text formatter has FullTimestamp disabled, want enabled")
	}
	if formatter.DisableTimestamp {
		t.Fatal("text formatter disables timestamps, want them enabled")
	}
	if got, want := log.GetLevel(), log.InfoLevel; got != want {
		t.Fatalf("log level = %v, want %v", got, want)
	}
}

func TestSetupLoggingJSON(t *testing.T) {
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_LEVEL", "debug")

	SetupLogging()

	if _, ok := log.StandardLogger().Formatter.(*log.JSONFormatter); !ok {
		t.Fatalf("formatter = %T, want *log.JSONFormatter", log.StandardLogger().Formatter)
	}
	if got, want := log.GetLevel(), log.DebugLevel; got != want {
		t.Fatalf("log level = %v, want %v", got, want)
	}
}

func TestSetupLoggingUnknownFormatFallsBackToText(t *testing.T) {
	t.Setenv("LOG_FORMAT", "yaml")

	SetupLogging()

	if _, ok := log.StandardLogger().Formatter.(*log.TextFormatter); !ok {
		t.Fatalf("formatter = %T, want *log.TextFormatter", log.StandardLogger().Formatter)
	}
}
