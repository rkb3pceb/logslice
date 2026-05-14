package parser

import (
	"strings"
	"testing"
)

func TestCountLines_BasicRFC3339(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-01T10:00:00Z INFO starting server",
		"2024-01-01T10:00:01Z DEBUG listening on :8080",
		"2024-01-01T10:00:02Z ERROR connection refused",
	}, "\n")

	result, err := CountLines(strings.NewReader(input), "2006-01-02T15:04:05Z07:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 3 {
		t.Errorf("Total: got %d, want 3", result.Total)
	}
	if result.Blank != 0 {
		t.Errorf("Blank: got %d, want 0", result.Blank)
	}
	if result.WithTimestamp != 3 {
		t.Errorf("WithTimestamp: got %d, want 3", result.WithTimestamp)
	}
}

func TestCountLines_BlankLines(t *testing.T) {
	input := "2024-01-01T10:00:00Z INFO a\n\n   \n2024-01-01T10:00:01Z INFO b\n"

	result, err := CountLines(strings.NewReader(input), "2006-01-02T15:04:05Z07:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 4 {
		t.Errorf("Total: got %d, want 4", result.Total)
	}
	if result.Blank != 2 {
		t.Errorf("Blank: got %d, want 2", result.Blank)
	}
	if result.WithTimestamp != 2 {
		t.Errorf("WithTimestamp: got %d, want 2", result.WithTimestamp)
	}
}

func TestCountLines_EmptyReader(t *testing.T) {
	result, err := CountLines(strings.NewReader(""), "2006-01-02T15:04:05Z07:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 0 {
		t.Errorf("Total: got %d, want 0", result.Total)
	}
}

func TestCountLines_NoFormat(t *testing.T) {
	input := "2024-01-01T10:00:00Z INFO a\n2024-01-01T10:00:01Z INFO b\n"

	result, err := CountLines(strings.NewReader(input), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 2 {
		t.Errorf("Total: got %d, want 2", result.Total)
	}
	if result.WithTimestamp != 0 {
		t.Errorf("WithTimestamp: got %d, want 0 when format is empty", result.WithTimestamp)
	}
}

func TestCountLines_MixedTimestamps(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-01T10:00:00Z INFO with timestamp",
		"no timestamp here",
		"2024-01-01T10:00:02Z WARN another timestamped line",
	}, "\n")

	result, err := CountLines(strings.NewReader(input), "2006-01-02T15:04:05Z07:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 3 {
		t.Errorf("Total: got %d, want 3", result.Total)
	}
	if result.WithTimestamp != 2 {
		t.Errorf("WithTimestamp: got %d, want 2", result.WithTimestamp)
	}
}
