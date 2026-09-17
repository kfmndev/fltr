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

func resetLogging() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
}

func requireLogLevel(t *testing.T, want log.Level) {
	t.Helper()
	if got := log.GetLevel(); got != want {
		t.Fatalf("log level = %v, want %v", got, want)
	}
}

func requireTextFormatter(t *testing.T, wantFullTimestamp, wantForceColors bool) {
	t.Helper()
	formatter, ok := log.StandardLogger().Formatter.(*log.TextFormatter)
	if !ok {
		t.Fatalf("formatter = %T, want *log.TextFormatter", log.StandardLogger().Formatter)
	}
	if formatter.FullTimestamp != wantFullTimestamp {
		t.Fatalf("FullTimestamp = %v, want %v", formatter.FullTimestamp, wantFullTimestamp)
	}
	if formatter.DisableTimestamp != !wantFullTimestamp {
		t.Fatalf("DisableTimestamp = %v, want %v", formatter.DisableTimestamp, !wantFullTimestamp)
	}
	if formatter.ForceColors != wantForceColors {
		t.Fatalf("ForceColors = %v, want %v", formatter.ForceColors, wantForceColors)
	}
	if got, want := formatter.TimestampFormat, "2006-01-02 15:04:05"; got != want {
		t.Fatalf("TimestampFormat = %q, want %q", got, want)
	}
}

func TestSetupLoggingDefaultsToText(t *testing.T) {
	resetLogging()
	t.Setenv("LOG_FORMAT", "")
	t.Setenv("LOG_LEVEL", "")

	SetupLogging()

	requireTextFormatter(t, true, true)
	requireLogLevel(t, log.InfoLevel)
	if out := log.StandardLogger().Out; out != os.Stderr {
		t.Fatalf("output = %v, want os.Stderr", out)
	}
}

func TestSetupLoggingText(t *testing.T) {
	resetLogging()
	t.Setenv("LOG_FORMAT", "text")
	t.Setenv("LOG_LEVEL", "warn")

	SetupLogging()

	requireTextFormatter(t, true, true)
	requireLogLevel(t, log.WarnLevel)
	if out := log.StandardLogger().Out; out != os.Stderr {
		t.Fatalf("output = %v, want os.Stderr", out)
	}
}

func TestSetupLoggingJSON(t *testing.T) {
	resetLogging()
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_LEVEL", "debug")

	SetupLogging()

	if _, ok := log.StandardLogger().Formatter.(*log.JSONFormatter); !ok {
		t.Fatalf("formatter = %T, want *log.JSONFormatter", log.StandardLogger().Formatter)
	}
	requireLogLevel(t, log.DebugLevel)
	if out := log.StandardLogger().Out; out != os.Stdout {
		t.Fatalf("output = %v, want os.Stdout", out)
	}
}

func TestSetupLoggingUnknownFormatFallsBackToText(t *testing.T) {
	resetLogging()
	t.Setenv("LOG_FORMAT", "yaml")

	SetupLogging()

	requireTextFormatter(t, true, true)
}

func TestGetLogLevelDefaultsToInfoOnInvalidValue(t *testing.T) {
	t.Setenv("LOG_LEVEL", "not-a-level")

	if got := GetLogLevel(); got != log.InfoLevel {
		t.Fatalf("GetLogLevel = %v, want %v", got, log.InfoLevel)
	}
}

func TestGetLogLevelParsesValidValue(t *testing.T) {
	t.Setenv("LOG_LEVEL", "error")

	if got := GetLogLevel(); got != log.ErrorLevel {
		t.Fatalf("GetLogLevel = %v, want %v", got, log.ErrorLevel)
	}
}
