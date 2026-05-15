package parser

import (
	"strings"
	"testing"
	"time"
)

func makeStatsLog(lines []string) *strings.Reader {
	return strings.NewReader(strings.Join(lines, "\n"))
}

func TestCollectStats_Basic(t *testing.T) {
	format := time.RFC3339
	logLines := []string{
		"2024-01-01T10:00:00Z level=info msg=start",
		"2024-01-01T10:05:00Z level=warn msg=middle",
		"2024-01-01T10:10:00Z level=info msg=end",
	}

	s, err := CollectStats(makeStatsLog(logLines), format)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalLines != 3 {
		t.Errorf("TotalLines: got %d, want 3", s.TotalLines)
	}
	if s.ParsedLines != 3 {
		t.Errorf("ParsedLines: got %d, want 3", s.ParsedLines)
	}
	if s.SkippedLines != 0 {
		t.Errorf("SkippedLines: got %d, want 0", s.SkippedLines)
	}
	wantDur := 10 * time.Minute
	if s.Duration != wantDur {
		t.Errorf("Duration: got %v, want %v", s.Duration, wantDur)
	}
}

func TestCollectStats_SkipsUnparseable(t *testing.T) {
	format := time.RFC3339
	logLines := []string{
		"2024-01-01T08:00:00Z level=info msg=ok",
		"not a timestamp at all",
		"2024-01-01T09:00:00Z level=info msg=ok",
	}

	s, err := CollectStats(makeStatsLog(logLines), format)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalLines != 3 {
		t.Errorf("TotalLines: got %d, want 3", s.TotalLines)
	}
	if s.ParsedLines != 2 {
		t.Errorf("ParsedLines: got %d, want 2", s.ParsedLines)
	}
	if s.SkippedLines != 1 {
		t.Errorf("SkippedLines: got %d, want 1", s.SkippedLines)
	}
}

func TestCollectStats_EmptyReader(t *testing.T) {
	s, err := CollectStats(strings.NewReader(""), time.RFC3339)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalLines != 0 {
		t.Errorf("TotalLines: got %d, want 0", s.TotalLines)
	}
	if !s.Earliest.IsZero() {
		t.Errorf("Earliest should be zero")
	}
	if s.Duration != 0 {
		t.Errorf("Duration should be zero")
	}
}

func TestCollectStats_EmptyFormat(t *testing.T) {
	logLines := []string{"2024-01-01T10:00:00Z msg=hello"}
	s, err := CollectStats(makeStatsLog(logLines), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalLines != 0 {
		t.Errorf("expected zero stats for empty format, got TotalLines=%d", s.TotalLines)
	}
}
