package parser

import (
	"strings"
	"testing"
)

func TestCountLines_BasicRFC3339(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-01T00:00:00Z starting server",
		"2024-01-01T00:00:01Z request received",
		"2024-01-01T00:00:02Z request handled",
	}, "\n")
	stats, err := CountLines(strings.NewReader(input), "2006-01-02T15:04:05Z07:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Total != 3 {
		t.Errorf("Total: got %d, want 3", stats.Total)
	}
	if stats.Timestamped != 3 {
		t.Errorf("Timestamped: got %d, want 3", stats.Timestamped)
	}
	if stats.Blank != 0 {
		t.Errorf("Blank: got %d, want 0", stats.Blank)
	}
}

func TestCountLines_BlankLines(t *testing.T) {
	input := "2024-01-01T00:00:00Z line one\n\n   \n2024-01-01T00:00:01Z line two"
	stats, err := CountLines(strings.NewReader(input), "2006-01-02T15:04:05Z07:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Total != 4 {
		t.Errorf("Total: got %d, want 4", stats.Total)
	}
	if stats.Blank != 2 {
		t.Errorf("Blank: got %d, want 2", stats.Blank)
	}
	if stats.Timestamped != 2 {
		t.Errorf("Timestamped: got %d, want 2", stats.Timestamped)
	}
}

func TestCountLines_EmptyReader(t *testing.T) {
	stats, err := CountLines(strings.NewReader(""), "2006-01-02T15:04:05Z07:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Total != 0 {
		t.Errorf("Total: got %d, want 0", stats.Total)
	}
}

func TestCountLines_NoFormat(t *testing.T) {
	input := "2024-01-01T00:00:00Z line one\nno timestamp here"
	stats, err := CountLines(strings.NewReader(input), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Timestamped != 0 {
		t.Errorf("Timestamped: got %d, want 0 when format is empty", stats.Timestamped)
	}
	if stats.Total != 2 {
		t.Errorf("Total: got %d, want 2", stats.Total)
	}
}

func TestCountLines_MixedTimestamps(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-01T00:00:00Z valid",
		"not a timestamp line",
		"2024-01-01T00:00:02Z also valid",
		"another plain line",
	}, "\n")
	stats, err := CountLines(strings.NewReader(input), "2006-01-02T15:04:05Z07:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Total != 4 {
		t.Errorf("Total: got %d, want 4", stats.Total)
	}
	if stats.Timestamped != 2 {
		t.Errorf("Timestamped: got %d, want 2", stats.Timestamped)
	}
}
