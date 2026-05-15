package parser

import (
	"strings"
	"testing"
)

func makeGrepLog() string {
	return strings.Join([]string{
		"2024-01-01T10:00:00Z INFO  server started",
		"2024-01-01T10:01:00Z ERROR connection refused",
		"2024-01-01T10:02:00Z INFO  request received",
		"2024-01-01T10:03:00Z WARN  disk usage high",
		"2024-01-01T10:04:00Z ERROR timeout exceeded",
		"",
		"not-a-timestamp garbage line",
	}, "\n")
}

const grepFmt = "2006-01-02T15:04:05Z"

func TestGrepLines_Basic(t *testing.T) {
	r := strings.NewReader(makeGrepLog())
	results, err := GrepLines(r, GrepOptions{Pattern: "error", Format: grepFmt})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !strings.Contains(results[0].Line, "connection refused") {
		t.Errorf("unexpected first match: %s", results[0].Line)
	}
}

func TestGrepLines_CaseSensitive(t *testing.T) {
	r := strings.NewReader(makeGrepLog())
	results, err := GrepLines(r, GrepOptions{Pattern: "ERROR", Format: grepFmt, CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestGrepLines_CaseSensitive_NoMatch(t *testing.T) {
	r := strings.NewReader(makeGrepLog())
	results, err := GrepLines(r, GrepOptions{Pattern: "error", Format: grepFmt, CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestGrepLines_EmptyFormat(t *testing.T) {
	r := strings.NewReader(makeGrepLog())
	_, err := GrepLines(r, GrepOptions{Pattern: "error", Format: ""})
	if err == nil {
		t.Fatal("expected error for empty format")
	}
}

func TestGrepLines_SkipsUnparseable(t *testing.T) {
	r := strings.NewReader(makeGrepLog())
	// "garbage line" contains "garbage" but has no timestamp — must be skipped.
	results, err := GrepLines(r, GrepOptions{Pattern: "garbage", Format: grepFmt})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestGrepLines_LineNumbers(t *testing.T) {
	r := strings.NewReader(makeGrepLog())
	results, err := GrepLines(r, GrepOptions{Pattern: "info", Format: grepFmt})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].LineNum != 1 {
		t.Errorf("expected line 1, got %d", results[0].LineNum)
	}
	if results[1].LineNum != 3 {
		t.Errorf("expected line 3, got %d", results[1].LineNum)
	}
}
