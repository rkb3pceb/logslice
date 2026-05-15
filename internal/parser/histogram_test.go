package parser

import (
	"strings"
	"testing"
	"time"
)

func makeHistogramLog() string {
	return strings.Join([]string{
		"2024-01-01T10:00:05Z app started",
		"2024-01-01T10:00:30Z request received",
		"2024-01-01T10:01:10Z request processed",
		"2024-01-01T10:01:45Z cache miss",
		"2024-01-01T10:02:05Z shutdown initiated",
	}, "\n")
}

func TestBuildHistogram_Basic(t *testing.T) {
	r := strings.NewReader(makeHistogramLog())
	h, err := BuildHistogram(r, time.RFC3339, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h.Buckets) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(h.Buckets))
	}
	if h.Buckets[0].Count != 2 {
		t.Errorf("bucket 0: expected 2 lines, got %d", h.Buckets[0].Count)
	}
	if h.Buckets[1].Count != 2 {
		t.Errorf("bucket 1: expected 2 lines, got %d", h.Buckets[1].Count)
	}
	if h.Buckets[2].Count != 1 {
		t.Errorf("bucket 2: expected 1 line, got %d", h.Buckets[2].Count)
	}
}

func TestBuildHistogram_EmptyFormat(t *testing.T) {
	r := strings.NewReader(makeHistogramLog())
	_, err := BuildHistogram(r, "", time.Minute)
	if err == nil {
		t.Fatal("expected error for empty format")
	}
}

func TestBuildHistogram_InvalidInterval(t *testing.T) {
	r := strings.NewReader(makeHistogramLog())
	_, err := BuildHistogram(r, time.RFC3339, 0)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestBuildHistogram_EmptyReader(t *testing.T) {
	r := strings.NewReader("")
	h, err := BuildHistogram(r, time.RFC3339, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h.Buckets) != 0 {
		t.Errorf("expected 0 buckets, got %d", len(h.Buckets))
	}
}

func TestHistogram_String(t *testing.T) {
	r := strings.NewReader(makeHistogramLog())
	h, err := BuildHistogram(r, time.RFC3339, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := h.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}

func TestHistogram_String_Empty(t *testing.T) {
	h := &Histogram{}
	if h.String() != "(empty histogram)" {
		t.Errorf("unexpected empty string: %q", h.String())
	}
}
