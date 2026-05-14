package parser

import (
	"strings"
	"testing"
)

func TestCountLines_BasicRFC3339(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-01T00:00:00Z INFO starting",
		"2024-01-01T00:00:01Z DEBUG ready",
		"2024-01-01T00:00:02Z WARN slow query",
	}, "\n")

	res, err := CountLines(strings.NewReader(input), "2006-01-02T15:04:05Z07:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 3 {
		t.Errorf("Total: want 3, got %d", res.Total)
	}
	if res.Matched != 3 {
		t.Errorf("Matched: want 3, got %d", res.Matched)
	}
	if res.Blank != 0 {
		t.Errorf("Blank: want 0, got %d", res.Blank)
	}
}

func TestCountLines_BlankLines(t *testing.T) {
	input := "2024-01-01T00:00:00Z INFO start\n\n   \n2024-01-01T00:00:01Z INFO end"
	res, err := CountLines(strings.NewReader(input), "2006-01-02T15:04:05Z07:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 4 {
		t.Errorf("Total: want 4, got %d", res.Total)
	}
	if res.Blank != 2 {
		t.Errorf("Blank: want 2, got %d", res.Blank)
	}
	if res.Matched != 2 {
		t.Errorf("Matched: want 2, got %d", res.Matched)
	}
}

func TestCountLines_EmptyReader(t *testing.T) {
	res, err := CountLines(strings.NewReader(""), "2006-01-02T15:04:05Z07:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 0 || res.Matched != 0 || res.Blank != 0 {
		t.Errorf("expected all zeros, got %+v", res)
	}
}

func TestCountLines_NoFormat(t *testing.T) {
	_, err := CountLines(strings.NewReader("anything"), "")
	if err == nil {
		t.Fatal("expected error for empty format, got nil")
	}
}

func TestCountLines_MixedTimestamps(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-01T00:00:00Z INFO ok",
		"not a timestamp at all",
		"2024-01-01T00:00:02Z INFO also ok",
	}, "\n")

	res, err := CountLines(strings.NewReader(input), "2006-01-02T15:04:05Z07:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 3 {
		t.Errorf("Total: want 3, got %d", res.Total)
	}
	if res.Matched != 2 {
		t.Errorf("Matched: want 2, got %d", res.Matched)
	}
}
