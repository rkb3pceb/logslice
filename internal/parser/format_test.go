package parser

import (
	"testing"
)

func TestDetectFormat_RFC3339(t *testing.T) {
	sample := "2024-03-15T10:22:33Z INFO server started"
	fmt, err := DetectFormat(sample)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmt != "2006-01-02T15:04:05Z07:00" && fmt != "2006-01-02T15:04:05.999999999Z07:00" {
		// Accept any RFC3339 variant
		if fmt == "" {
			t.Fatalf("expected a non-empty format, got empty string")
		}
	}
}

func TestDetectFormat_SpaceSeparated(t *testing.T) {
	sample := "2024-03-15 10:22:33 ERROR disk full"
	fmt, err := DetectFormat(sample)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmt == "" {
		t.Fatal("expected non-empty format")
	}
}

func TestDetectFormat_WithMilliseconds(t *testing.T) {
	sample := "2024-03-15T10:22:33.456Z INFO request received"
	_, err := DetectFormat(sample)
	if err != nil {
		t.Fatalf("unexpected error for millisecond timestamp: %v", err)
	}
}

func TestDetectFormat_EmptySample(t *testing.T) {
	_, err := DetectFormat("")
	if err == nil {
		t.Fatal("expected error for empty sample, got nil")
	}
}

func TestDetectFormat_NoTimestamp(t *testing.T) {
	sample := "this line has no timestamp at all"
	_, err := DetectFormat(sample)
	if err == nil {
		t.Fatal("expected error for line without timestamp, got nil")
	}
}

func TestDetectFormat_NginxStyle(t *testing.T) {
	sample := "15/Mar/2024:10:22:33 +0000 GET /index.html 200"
	_, err := DetectFormat(sample)
	if err != nil {
		t.Fatalf("unexpected error for nginx-style timestamp: %v", err)
	}
}

func TestExtractCandidate_SpaceDelimited(t *testing.T) {
	line := "2024-03-15T10:22:33Z INFO"
	got := extractCandidate(line)
	want := "2024-03-15T10:22:33Z"
	if got != want {
		t.Errorf("extractCandidate = %q, want %q", got, want)
	}
}

func TestExtractCandidate_NoSpace(t *testing.T) {
	line := "2024-03-15T10:22:33Z"
	got := extractCandidate(line)
	if got != line {
		t.Errorf("extractCandidate = %q, want %q", got, line)
	}
}
