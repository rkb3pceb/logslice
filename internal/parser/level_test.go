package parser

import (
	"bytes"
	"strings"
	"testing"
)

func TestDetectLevel_Known(t *testing.T) {
	cases := []struct {
		line  string
		want  LogLevel
	}{
		{"2024-01-01T00:00:00Z INFO  server started", LevelInfo},
		{"2024-01-01T00:00:00Z DEBUG  connecting", LevelDebug},
		{"2024-01-01T00:00:00Z WARN  high memory", LevelWarn},
		{"2024-01-01T00:00:00Z ERROR  disk full", LevelError},
		{"2024-01-01T00:00:00Z FATAL  panic", LevelFatal},
		{"2024-01-01T00:00:00Z no level here", LevelUnknown},
	}
	for _, tc := range cases {
		got := DetectLevel(tc.line)
		if got != tc.want {
			t.Errorf("DetectLevel(%q) = %q; want %q", tc.line, got, tc.want)
		}
	}
}

func TestDetectLevel_CaseInsensitive(t *testing.T) {
	if got := DetectLevel("timestamp warn something"); got != LevelWarn {
		t.Errorf("expected WARN, got %q", got)
	}
}

func makeLevelLog() string {
	return strings.Join([]string{
		"2024-01-01T00:00:01Z INFO  app started",
		"2024-01-01T00:00:02Z DEBUG  init done",
		"2024-01-01T00:00:03Z WARN  low disk",
		"2024-01-01T00:00:04Z ERROR  timeout",
		"2024-01-01T00:00:05Z INFO  ready",
	}, "\n") + "\n"
}

func TestFilterByLevel_Basic(t *testing.T) {
	r := strings.NewReader(makeLevelLog())
	var buf bytes.Buffer
	err := FilterByLevel(r, &buf, "2006-01-02T15:04:05Z", []LogLevel{LevelError, LevelWarn})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(lines), lines)
	}
}

func TestFilterByLevel_All(t *testing.T) {
	r := strings.NewReader(makeLevelLog())
	var buf bytes.Buffer
	err := FilterByLevel(r, &buf, "2006-01-02T15:04:05Z", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(lines))
	}
}

func TestFilterByLevel_EmptyFormat(t *testing.T) {
	r := strings.NewReader(makeLevelLog())
	var buf bytes.Buffer
	err := FilterByLevel(r, &buf, "", []LogLevel{LevelInfo})
	if err != ErrEmptyFormat {
		t.Errorf("expected ErrEmptyFormat, got %v", err)
	}
}
