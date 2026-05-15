package cli

import (
	"os"
	"strings"
	"testing"
)

func writeGrepLog(t *testing.T) string {
	t.Helper()
	content := strings.Join([]string{
		"2024-01-01T10:00:00Z INFO  server started",
		"2024-01-01T10:01:00Z ERROR connection refused",
		"2024-01-01T10:02:00Z INFO  request received",
		"2024-01-01T10:03:00Z ERROR timeout exceeded",
	}, "\n") + "\n"

	f, err := os.CreateTemp("", "grep-*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestRunGrep_Basic(t *testing.T) {
	path := writeGrepLog(t)
	err := RunGrep([]string{"error", path}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunGrep_CaseSensitive(t *testing.T) {
	path := writeGrepLog(t)
	err := RunGrep([]string{"ERROR", path}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunGrep_MissingArgs(t *testing.T) {
	err := RunGrep([]string{"onlyone"}, false)
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunGrep_FileNotFound(t *testing.T) {
	err := RunGrep([]string{"error", "/nonexistent/path/file.log"}, false)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunGrep_NoMatches(t *testing.T) {
	path := writeGrepLog(t)
	// Pattern that won't match any timestamped line.
	err := RunGrep([]string{"zzznomatch", path}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
