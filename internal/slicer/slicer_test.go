package slicer_test

import (
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/slicer"
)

const layout = "2006-01-02T15:04:05"

func lines(ll ...string) string {
	return strings.Join(ll, "\n") + "\n"
}

var logData = lines(
	"2024-03-01T08:00:00 INFO  server started",
	"2024-03-01T08:05:00 DEBUG request received",
	"2024-03-01T08:10:00 INFO  processing",
	"2024-03-01T08:15:00 WARN  slow query",
	"2024-03-01T08:20:00 ERROR timeout",
	"not a log line",
	"2024-03-01T08:25:00 INFO  recovered",
)

func ts(s string) time.Time {
	t, err := time.Parse(layout, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSlice_BasicRange(t *testing.T) {
	r := strings.NewReader(logData)
	var w strings.Builder

	n, err := slicer.Slice(r, &w, slicer.Options{
		From:   ts("2024-03-01T08:05:00"),
		To:     ts("2024-03-01T08:15:00"),
		Layout: layout,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 lines, got %d", n)
	}
}

func TestSlice_NoUpperBound(t *testing.T) {
	r := strings.NewReader(logData)
	var w strings.Builder

	n, err := slicer.Slice(r, &w, slicer.Options{
		From:   ts("2024-03-01T08:20:00"),
		Layout: layout,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 lines, got %d", n)
	}
}

func TestSlice_InvalidRange(t *testing.T) {
	r := strings.NewReader(logData)
	var w strings.Builder

	_, err := slicer.Slice(r, &w, slicer.Options{
		From:   ts("2024-03-01T09:00:00"),
		To:     ts("2024-03-01T08:00:00"),
		Layout: layout,
	})
	if err == nil {
		t.Fatal("expected error for inverted range, got nil")
	}
}

func TestSlice_EmptyLayout(t *testing.T) {
	_, err := slicer.Slice(strings.NewReader(""), &strings.Builder{}, slicer.Options{})
	if err == nil {
		t.Fatal("expected error for empty layout")
	}
}
