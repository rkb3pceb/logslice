package parser

import (
	"strings"
	"testing"
)

func TestDetectFormatFromReader_RFC3339(t *testing.T) {
	input := strings.NewReader(
		"2024-01-15T10:00:00Z INFO server started\n" +
			"2024-01-15T10:00:01Z DEBUG request received\n",
	)
	format, sample, err := DetectFormatFromReader(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if format == "" {
		t.Fatal("expected a format to be detected, got empty string")
	}
	if len(sample) != 2 {
		t.Fatalf("expected 2 sample lines, got %d", len(sample))
	}
}

func TestDetectFormatFromReader_Empty(t *testing.T) {
	input := strings.NewReader("")
	format, sample, err := DetectFormatFromReader(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if format != "" {
		t.Fatalf("expected empty format for empty input, got %q", format)
	}
	if len(sample) != 0 {
		t.Fatalf("expected 0 sample lines, got %d", len(sample))
	}
}

func TestDetectFormatFromReader_SkipsBlanks(t *testing.T) {
	input := strings.NewReader(
		"\n" +
			"   \n" +
			"2024-03-01 08:30:00 WARN disk usage high\n",
	)
	_, sample, err := DetectFormatFromReader(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, line := range sample {
		if strings.TrimSpace(line) == "" {
			t.Errorf("sample should not contain blank lines, got %q", line)
		}
	}
}

func TestDetectFormatFromReader_SampleCapped(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 20; i++ {
		sb.WriteString("2024-06-01T12:00:00Z INFO line\n")
	}
	_, sample, err := DetectFormatFromReader(strings.NewReader(sb.String()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sample) > SampleSize {
		t.Fatalf("sample size %d exceeds SampleSize %d", len(sample), SampleSize)
	}
}

func TestDetectFormatFromReader_NoTimestamp(t *testing.T) {
	input := strings.NewReader(
		"just some plain text\n" +
			"no timestamps here\n",
	)
	format, _, err := DetectFormatFromReader(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if format != "" {
		t.Fatalf("expected empty format for unrecognised input, got %q", format)
	}
}
