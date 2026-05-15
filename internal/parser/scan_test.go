package parser

import (
	"strings"
	"testing"
	"time"
)

const rfc3339Format = "2006-01-02T15:04:05Z07:00"

func TestScanLines_Basic(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-01T10:00:00Z app started",
		"2024-01-01T10:01:00Z request received",
		"2024-01-01T10:02:00Z request done",
	}, "\n")

	lines, err := ScanLines(strings.NewReader(input), rfc3339Format)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0].LineNumber != 1 {
		t.Errorf("expected LineNumber 1, got %d", lines[0].LineNumber)
	}
	if lines[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp for first line")
	}
}

func TestScanLines_SkipsBlanks(t *testing.T) {
	input := "2024-01-01T10:00:00Z line one\n\n   \n2024-01-01T10:01:00Z line two"
	lines, err := ScanLines(strings.NewReader(input), rfc3339Format)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Errorf("expected 2 lines (blanks skipped), got %d", len(lines))
	}
}

func TestScanLines_EmptyFormat(t *testing.T) {
	input := "some log line without timestamp"
	lines, err := ScanLines(strings.NewReader(input), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if !lines[0].Timestamp.IsZero() {
		t.Error("expected zero timestamp when format is empty")
	}
}

func TestFilterByRange_Basic(t *testing.T) {
	mk := func(s string) time.Time {
		t2, _ := time.Parse(rfc3339Format, s)
		return t2
	}
	lines := []LogLine{
		{Raw: "a", Timestamp: mk("2024-01-01T10:00:00Z"), LineNumber: 1},
		{Raw: "b", Timestamp: mk("2024-01-01T10:05:00Z"), LineNumber: 2},
		{Raw: "c", Timestamp: mk("2024-01-01T10:10:00Z"), LineNumber: 3},
	}

	got := FilterByRange(lines, mk("2024-01-01T10:03:00Z"), mk("2024-01-01T10:07:00Z"))
	if len(got) != 1 || got[0].Raw != "b" {
		t.Errorf("expected only line 'b', got %v", got)
	}
}

func TestFilterByRange_NoUpperBound(t *testing.T) {
	mk := func(s string) time.Time {
		t2, _ := time.Parse(rfc3339Format, s)
		return t2
	}
	lines := []LogLine{
		{Raw: "a", Timestamp: mk("2024-01-01T09:00:00Z"), LineNumber: 1},
		{Raw: "b", Timestamp: mk("2024-01-01T10:00:00Z"), LineNumber: 2},
		{Raw: "c", Timestamp: mk("2024-01-01T11:00:00Z"), LineNumber: 3},
	}

	got := FilterByRange(lines, mk("2024-01-01T10:00:00Z"), time.Time{})
	if len(got) != 2 {
		t.Errorf("expected 2 lines, got %d", len(got))
	}
}

func TestFilterByRange_SkipsZeroTimestamps(t *testing.T) {
	mk := func(s string) time.Time {
		t2, _ := time.Parse(rfc3339Format, s)
		return t2
	}
	lines := []LogLine{
		{Raw: "unparseable", Timestamp: time.Time{}, LineNumber: 1},
		{Raw: "b", Timestamp: mk("2024-01-01T10:00:00Z"), LineNumber: 2},
	}

	got := FilterByRange(lines, mk("2024-01-01T09:00:00Z"), time.Time{})
	if len(got) != 1 || got[0].Raw != "b" {
		t.Errorf("expected only 'b', got %v", got)
	}
}
