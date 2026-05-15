package parser

import (
	"strings"
	"testing"
	"time"
)

func makeTailLog(entries []string) string {
	return strings.Join(entries, "\n") + "\n"
}

func TestTailLines_Basic(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z INFO a",
		"2024-01-01T10:00:01Z INFO b",
		"2024-01-01T10:00:02Z INFO c",
		"2024-01-01T10:00:03Z INFO d",
		"2024-01-01T10:00:04Z INFO e",
	}
	r := strings.NewReader(makeTailLog(lines))
	res, err := TailLines(r, time.RFC3339, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(res.Lines))
	}
	if res.Lines[0] != lines[2] {
		t.Errorf("expected first tail line %q, got %q", lines[2], res.Lines[0])
	}
	if res.Lines[2] != lines[4] {
		t.Errorf("expected last tail line %q, got %q", lines[4], res.Lines[2])
	}
}

func TestTailLines_FewerThanN(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z INFO a",
		"2024-01-01T10:00:01Z INFO b",
	}
	r := strings.NewReader(makeTailLog(lines))
	res, err := TailLines(r, time.RFC3339, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(res.Lines))
	}
}

func TestTailLines_EmptyFormat(t *testing.T) {
	r := strings.NewReader("2024-01-01T10:00:00Z INFO a\n")
	_, err := TailLines(r, "", 5)
	if err != ErrEmptyFormat {
		t.Fatalf("expected ErrEmptyFormat, got %v", err)
	}
}

func TestTailLines_EmptyReader(t *testing.T) {
	r := strings.NewReader("")
	res, err := TailLines(r, time.RFC3339, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Lines) != 0 {
		t.Errorf("expected 0 lines, got %d", len(res.Lines))
	}
}

func TestTailLines_SkipsBlanks(t *testing.T) {
	input := "\n2024-01-01T10:00:00Z INFO a\n\n2024-01-01T10:00:01Z INFO b\n"
	r := strings.NewReader(input)
	res, err := TailLines(r, time.RFC3339, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(res.Lines))
	}
}

func TestTailLines_TimestampBounds(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z INFO a",
		"2024-01-01T10:00:01Z INFO b",
		"2024-01-01T10:00:02Z INFO c",
	}
	r := strings.NewReader(makeTailLog(lines))
	res, err := TailLines(r, time.RFC3339, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expFirst, _ := time.Parse(time.RFC3339, "2024-01-01T10:00:01Z")
	expLast, _ := time.Parse(time.RFC3339, "2024-01-01T10:00:02Z")
	if !res.First.Equal(expFirst) {
		t.Errorf("expected First %v, got %v", expFirst, res.First)
	}
	if !res.Last.Equal(expLast) {
		t.Errorf("expected Last %v, got %v", expLast, res.Last)
	}
}
