package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeHistogramLog(t *testing.T) string {
	t.Helper()
	content := strings.Join([]string{
		"2024-03-01T08:00:05Z boot complete",
		"2024-03-01T08:00:45Z first request",
		"2024-03-01T08:01:12Z cache warm",
		"2024-03-01T08:01:58Z db query",
		"2024-03-01T08:02:33Z response sent",
	}, "\n")
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp log: %v", err)
	}
	return p
}

func TestRunHistogram_Basic(t *testing.T) {
	path := writeHistogramLog(t)
	if err := RunHistogram([]string{path, "1m"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunHistogram_MissingArgs(t *testing.T) {
	err := RunHistogram([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
	if !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage hint in error, got: %v", err)
	}
}

func TestRunHistogram_InvalidInterval(t *testing.T) {
	path := writeHistogramLog(t)
	err := RunHistogram([]string{path, "notaduration"})
	if err == nil {
		t.Fatal("expected error for invalid interval")
	}
}

func TestRunHistogram_FileNotFound(t *testing.T) {
	err := RunHistogram([]string{"/nonexistent/path.log", "1m"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunHistogram_ZeroInterval(t *testing.T) {
	path := writeHistogramLog(t)
	err := RunHistogram([]string{path, "0s"})
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}
