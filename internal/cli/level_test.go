package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeLevelLog(t *testing.T) string {
	t.Helper()
	content := strings.Join([]string{
		"2024-01-01T00:00:01Z INFO  app started",
		"2024-01-01T00:00:02Z DEBUG  init done",
		"2024-01-01T00:00:03Z WARN  low disk",
		"2024-01-01T00:00:04Z ERROR  timeout",
		"2024-01-01T00:00:05Z INFO  ready",
	}, "\n") + "\n"

	tmp := filepath.Join(t.TempDir(), "level.log")
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp log: %v", err)
	}
	return tmp
}

func TestRunLevel_Basic(t *testing.T) {
	path := writeLevelLog(t)
	if err := RunLevel([]string{path, "ERROR,WARN"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunLevel_MissingArgs(t *testing.T) {
	err := RunLevel([]string{"only-one-arg"})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunLevel_FileNotFound(t *testing.T) {
	err := RunLevel([]string{"/nonexistent/file.log", "INFO"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunLevel_EmptyFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "empty.log")
	if err := os.WriteFile(tmp, []byte{}, 0o644); err != nil {
		t.Fatalf("write empty file: %v", err)
	}
	err := RunLevel([]string{tmp, "INFO"})
	if err == nil {
		t.Fatal("expected error for empty file (no format detected)")
	}
}
