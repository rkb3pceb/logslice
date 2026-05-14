package cli

import (
	"bytes"
	"strings"
	"testing"
)

const sampleLog = `2024-01-15T10:00:00Z INFO starting service
2024-01-15T10:01:00Z INFO request received
2024-01-15T10:02:00Z WARN slow query detected
2024-01-15T10:03:00Z INFO request completed
2024-01-15T10:04:00Z INFO shutting down
`

func TestRun_BasicSlice(t *testing.T) {
	in := strings.NewReader(sampleLog)
	var out bytes.Buffer
	var errOut bytes.Buffer

	err := Run([]string{
		"--from", "2024-01-15T10:01:00Z",
		"--to", "2024-01-15T10:03:00Z",
	}, in, &out, &errOut)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "request received") {
		t.Errorf("expected 'request received' in output, got:\n%s", got)
	}
	if strings.Contains(got, "starting service") {
		t.Errorf("expected 'starting service' to be excluded, got:\n%s", got)
	}
}

func TestRun_MissingFrom(t *testing.T) {
	in := strings.NewReader(sampleLog)
	var out, errOut bytes.Buffer

	err := Run([]string{}, in, &out, &errOut)
	if err == nil {
		t.Fatal("expected error when --from is missing")
	}
	if !strings.Contains(err.Error(), "--from is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRun_InvalidFrom(t *testing.T) {
	in := strings.NewReader(sampleLog)
	var out, errOut bytes.Buffer

	err := Run([]string{"--from", "not-a-time"}, in, &out, &errOut)
	if err == nil {
		t.Fatal("expected error for invalid --from")
	}
	if !strings.Contains(err.Error(), "invalid --from") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRun_NoUpperBound(t *testing.T) {
	in := strings.NewReader(sampleLog)
	var out, errOut bytes.Buffer

	err := Run([]string{
		"--from", "2024-01-15T10:03:00Z",
	}, in, &out, &errOut)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "request completed") {
		t.Errorf("expected lines from 10:03 onward, got:\n%s", got)
	}
	if !strings.Contains(got, "shutting down") {
		t.Errorf("expected 'shutting down' in unbounded output, got:\n%s", got)
	}
}
