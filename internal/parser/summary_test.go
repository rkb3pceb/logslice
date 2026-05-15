package parser

import (
	"strings"
	"testing"
	"time"
)

func mustParseRFC3339(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSummarize_Basic(t *testing.T) {
	log := strings.Join([]string{
		"2024-03-01T08:00:00Z INFO boot",
		"2024-03-01T08:05:00Z INFO ready",
		"2024-03-01T09:00:00Z WARN high load",
	}, "\n") + "\n"

	s, err := Summarize(strings.NewReader(log), rfc3339Layout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.LineCount != 3 {
		t.Errorf("LineCount: want 3, got %d", s.LineCount)
	}
	if !s.First.Equal(mustParseRFC3339("2024-03-01T08:00:00Z")) {
		t.Errorf("First: want 08:00, got %v", s.First)
	}
	if !s.Last.Equal(mustParseRFC3339("2024-03-01T09:00:00Z")) {
		t.Errorf("Last: want 09:00, got %v", s.Last)
	}
}

func TestSummarize_Empty(t *testing.T) {
	s, err := Summarize(strings.NewReader(""), rfc3339Layout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.LineCount != 0 {
		t.Errorf("expected 0 lines, got %d", s.LineCount)
	}
}

func TestSummarize_SkipsUnparseable(t *testing.T) {
	log := strings.Join([]string{
		"2024-03-01T10:00:00Z INFO start",
		"this line has no timestamp",
		"2024-03-01T10:30:00Z INFO end",
	}, "\n") + "\n"

	s, err := Summarize(strings.NewReader(log), rfc3339Layout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.LineCount != 2 {
		t.Errorf("LineCount: want 2, got %d", s.LineCount)
	}
}

func TestSummary_String(t *testing.T) {
	s := Summary{
		First:     mustParseRFC3339("2024-01-01T00:00:00Z"),
		Last:      mustParseRFC3339("2024-01-01T01:00:00Z"),
		LineCount: 42,
	}
	out := s.String()
	if out == "" {
		t.Error("String() returned empty string")
	}
}
