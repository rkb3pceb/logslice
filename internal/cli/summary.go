package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/yourorg/logslice/internal/parser"
)

// RunSummary reads the log file specified by args and prints a human-readable
// summary (first timestamp, last timestamp, total line count) to out.
// It accepts the following args:
//
//	args[0] — path to the log file
func RunSummary(args []string, out io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: logslice summary <file>")
	}

	filePath := args[0]

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open %q: %w", filePath, err)
	}
	defer f.Close()

	fmt, err := detectFormat(f)
	if err != nil {
		return fmt.Errorf("detect format: %w", err)
	}

	// Rewind after format detection.
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek: %w", err)
	}

	summary, err := parser.Summarize(f, fmt)
	if err != nil {
		return fmt.Errorf("summarize: %w", err)
	}

	_, err = fmt.Fprintln(out, summary)
	return err
}
