package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestRunSummary_Basic(t *testing.T) {
	log := strings.Join([]string{
		"2024-01-01T10:00:00Z INFO starting",
		"2024-01-01T10:01:00Z INFO processing",
		"2024-01-01T10:02:00Z INFO done",
	}, "\n") + "\n"

	path := writeTempLog(t, log)
	var buf bytes.Buffer
	if err := RunSummary([]string{path}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "2024-01-01T10:00:00Z") {
		t.Errorf("expected first timestamp in output, got: %s", out)
	}
	if !strings.Contains(out, "2024-01-01T10:02:00Z") {
		t.Errorf("expected last timestamp in output, got: %s", out)
	}
}

func TestRunSummary_MissingArgs(t *testing.T) {
	var buf bytes.Buffer
	err := RunSummary([]string{}, &buf)
	if err == nil {
		t.Fatal("expected error for missing args, got nil")
	}
	if !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage hint in error, got: %v", err)
	}
}

func TestRunSummary_FileNotFound(t *testing.T) {
	var buf bytes.Buffer
	err := RunSummary([]string{"/nonexistent/path/file.log"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRunSummary_EmptyFile(t *testing.T) {
	path := writeTempLog(t, "")
	var buf bytes.Buffer
	// An empty file should not crash; it may return an error or an empty summary.
	_ = RunSummary([]string{path}, &buf)
}
