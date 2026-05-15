package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/yourorg/logslice/internal/parser"
)

// RunStats prints aggregate statistics for the given log file to out.
// It detects the timestamp format automatically from the file contents.
func RunStats(path string, out io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()

	format, err := detectFormat(f)
	if err != nil {
		return fmt.Errorf("detect format: %w", err)
	}
	if format == "" {
		return fmt.Errorf("could not detect timestamp format in %q", path)
	}

	// Rewind after format detection.
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek: %w", err)
	}

	s, err := parser.CollectStats(f, format)
	if err != nil {
		return fmt.Errorf("collect stats: %w", err)
	}

	fmt.Fprintf(out, "File:          %s\n", path)
	fmt.Fprintf(out, "Format:        %s\n", format)
	fmt.Fprintf(out, "Total lines:   %d\n", s.TotalLines)
	fmt.Fprintf(out, "Parsed lines:  %d\n", s.ParsedLines)
	fmt.Fprintf(out, "Skipped lines: %d\n", s.SkippedLines)
	if !s.Earliest.IsZero() {
		fmt.Fprintf(out, "Earliest:      %s\n", s.Earliest.Format(format))
		fmt.Fprintf(out, "Latest:        %s\n", s.Latest.Format(format))
		fmt.Fprintf(out, "Duration:      %s\n", s.Duration)
	} else {
		fmt.Fprintf(out, "Earliest:      (none)\n")
		fmt.Fprintf(out, "Latest:        (none)\n")
	}
	return nil
}
