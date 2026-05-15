package parser

import (
	"strings"
	"testing"
)

func makeContextLog() string {
	return strings.Join([]string{
		`2024-01-01T10:00:00Z INFO  startup complete`,
		`2024-01-01T10:00:01Z DEBUG checking config`,
		`2024-01-01T10:00:02Z ERROR connection refused`,
		`2024-01-01T10:00:03Z INFO  retrying connection`,
		`2024-01-01T10:00:04Z INFO  connection established`,
		`2024-01-01T10:00:05Z DEBUG all systems nominal`,
	}, "\n")
}

const contextFormat = "2006-01-02T15:04:05Z"

func TestGrepWithContext_Basic(t *testing.T) {
	r := strings.NewReader(makeContextLog())
	results, err := GrepWithContext(r, contextFormat, "ERROR", 1, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// expect lines at index 1,2,3 (before, match, after)
	if len(results) != 3 {
		t.Fatalf("expected 3 context lines, got %d", len(results))
	}
	if !results[1].Matched {
		t.Errorf("expected middle result to be the matched line")
	}
	if results[0].Matched || results[2].Matched {
		t.Errorf("context lines should not be marked as matched")
	}
	if !strings.Contains(results[1].Line, "ERROR") {
		t.Errorf("matched line should contain ERROR")
	}
}

func TestGrepWithContext_CaseInsensitive(t *testing.T) {
	r := strings.NewReader(makeContextLog())
	results, err := GrepWithContext(r, contextFormat, "error", 0, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Matched {
		t.Errorf("result should be marked as matched")
	}
}

func TestGrepWithContext_NoMatch(t *testing.T) {
	r := strings.NewReader(makeContextLog())
	results, err := GrepWithContext(r, contextFormat, "FATAL", 2, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestGrepWithContext_EmptyFormat(t *testing.T) {
	r := strings.NewReader(makeContextLog())
	_, err := GrepWithContext(r, "", "ERROR", 1, true)
	if err == nil {
		t.Error("expected error for empty format")
	}
}

func TestGrepWithContext_ZeroContext(t *testing.T) {
	r := strings.NewReader(makeContextLog())
	results, err := GrepWithContext(r, contextFormat, "ERROR", 0, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result with zero context, got %d", len(results))
	}
	if !results[0].Matched {
		t.Errorf("result should be marked matched")
	}
}

func TestGrepWithContext_TimestampParsed(t *testing.T) {
	r := strings.NewReader(makeContextLog())
	results, err := GrepWithContext(r, contextFormat, "ERROR", 0, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Time.IsZero() {
		t.Errorf("expected non-zero timestamp on matched line")
	}
}
