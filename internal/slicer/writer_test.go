package slicer

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

const rfcFmt = "2006-01-02T15:04:05Z07:00"

func makeLog() string {
	return strings.Join([]string{
		"2024-01-01T10:00:00Z INFO  startup complete",
		"2024-01-01T10:01:00Z DEBUG request received",
		"2024-01-01T10:02:00Z INFO  processing",
		"2024-01-01T10:03:00Z WARN  slow query",
		"2024-01-01T10:04:00Z ERROR timeout",
	}, "\n")
}

func mustTime(s string) time.Time {
	t, err := time.Parse(rfcFmt, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestWriteSlice_BasicRange(t *testing.T) {
	src := strings.NewReader(makeLog())
	var dst bytes.Buffer

	n, err := WriteSlice(&dst, src, WriteOptions{
		Format: rfcFmt,
		From:   mustTime("2024-01-01T10:01:00Z"),
		To:     mustTime("2024-01-01T10:03:00Z"),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Fatalf("expected 3 lines written, got %d", n)
	}
	if !strings.Contains(dst.String(), "request received") {
		t.Error("expected 'request received' in output")
	}
	if strings.Contains(dst.String(), "startup complete") {
		t.Error("'startup complete' should be excluded")
	}
	if strings.Contains(dst.String(), "timeout") {
		t.Error("'timeout' should be excluded")
	}
}

func TestWriteSlice_NoUpperBound(t *testing.T) {
	src := strings.NewReader(makeLog())
	var dst bytes.Buffer

	n, err := WriteSlice(&dst, src, WriteOptions{
		Format: rfcFmt,
		From:   mustTime("2024-01-01T10:03:00Z"),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected 2 lines, got %d", n)
	}
}

func TestWriteSlice_EmptyFormat(t *testing.T) {
	src := strings.NewReader(makeLog())
	var dst bytes.Buffer

	_, err := WriteSlice(&dst, src, WriteOptions{})
	if err == nil {
		t.Fatal("expected error for empty format, got nil")
	}
}

func TestWriteSlice_NoMatchingLines(t *testing.T) {
	src := strings.NewReader(makeLog())
	var dst bytes.Buffer

	n, err := WriteSlice(&dst, src, WriteOptions{
		Format: rfcFmt,
		From:   mustTime("2024-01-01T12:00:00Z"),
		To:     mustTime("2024-01-01T13:00:00Z"),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Fatalf("expected 0 lines, got %d", n)
	}
	if dst.Len() != 0 {
		t.Error("expected empty output")
	}
}
