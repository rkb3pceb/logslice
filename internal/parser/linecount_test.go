package parser

import (
	"strings"
	"testing"
)

const rfc3339Layout = "2006-01-02T15:04:05Z07:00"

func makeLogLines(lines []string) string {
	return strings.Join(lines, "\n") + "\n"
}

func TestCountLines_BasicRFC3339(t *testing.T) {
	log := makeLogLines([]string{
		"2024-01-01T00:00:00Z INFO starting",
		"2024-01-01T00:01:00Z INFO running",
		"2024-01-01T00:02:00Z INFO done",
	})
	n, err := CountLines(strings.NewReader(log), rfc3339Layout,
		"2024-01-01T00:00:00Z", "2024-01-01T00:02:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 lines, got %d", n)
	}
}

func TestCountLines_BlankLines(t *testing.T) {
	log := makeLogLines([]string{
		"2024-01-01T00:00:00Z INFO a",
		"",
		"   ",
		"2024-01-01T00:01:00Z INFO b",
	})
	n, err := CountLines(strings.NewReader(log), rfc3339Layout,
		"2024-01-01T00:00:00Z", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 lines, got %d", n)
	}
}

func TestCountLines_EmptyReader(t *testing.T) {
	n, err := CountLines(strings.NewReader(""), rfc3339Layout,
		"2024-01-01T00:00:00Z", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
}

func TestCountLines_NoFormat(t *testing.T) {
	_, err := CountLines(strings.NewReader("anything"), rfc3339Layout,
		"not-a-timestamp", "")
	if err == nil {
		t.Fatal("expected error for invalid from timestamp")
	}
}

func TestCountLines_MixedTimestamps(t *testing.T) {
	log := makeLogLines([]string{
		"2024-01-01T00:00:00Z INFO first",
		"no timestamp here",
		"2024-01-01T00:01:00Z INFO second",
		"2024-01-01T00:03:00Z INFO third",
	})
	n, err := CountLines(strings.NewReader(log), rfc3339Layout,
		"2024-01-01T00:00:00Z", "2024-01-01T00:02:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 lines, got %d", n)
	}
}
