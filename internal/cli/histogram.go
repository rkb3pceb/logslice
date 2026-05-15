package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/user/logslice/internal/parser"
)

// RunHistogram prints a time-bucketed histogram of log line frequency.
// Args: [file, interval]
func RunHistogram(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: logslice histogram <file> <interval>")
	}

	filePath := args[0]
	intervalStr := args[1]

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return fmt.Errorf("invalid interval %q: %w", intervalStr, err)
	}
	if interval <= 0 {
		return fmt.Errorf("interval must be positive, got %s", intervalStr)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open file %q: %w", filePath, err)
	}
	defer f.Close()

	fmt.Fprintf(os.Stderr, "detecting format from %s...\n", filePath)
	fmt2, err := detectFormat(filePath)
	if err != nil {
		return fmt.Errorf("format detection failed: %w", err)
	}

	// Re-open for histogram since detectFormat consumed the reader.
	f.Seek(0, 0)

	h, err := parser.BuildHistogram(f, fmt2, interval)
	if err != nil {
		return fmt.Errorf("histogram error: %w", err)
	}

	fmt.Print(h.String())
	return nil
}
