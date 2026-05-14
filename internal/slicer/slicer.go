package slicer

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// Options configures the slicing behaviour.
type Options struct {
	// From is the inclusive start of the time range.
	From time.Time
	// To is the inclusive end of the time range.
	To time.Time
	// Layout is the timestamp layout used to parse log lines.
	Layout string
}

// Slice reads log lines from r, writes lines whose timestamps fall within
// [opts.From, opts.To] to w, and returns the number of lines written.
func Slice(r io.Reader, w io.Writer, opts Options) (int, error) {
	if opts.Layout == "" {
		return 0, fmt.Errorf("slicer: layout must not be empty")
	}
	if !opts.To.IsZero() && opts.From.After(opts.To) {
		return 0, fmt.Errorf("slicer: From (%s) is after To (%s)", opts.From, opts.To)
	}

	scanner := bufio.NewScanner(r)
	written := 0

	for scanner.Scan() {
		line := scanner.Text()
		ts, err := parser.ParseTimestamp(line, opts.Layout)
		if err != nil {
			// Skip lines that cannot be parsed.
			continue
		}

		if ts.Before(opts.From) {
			continue
		}
		if !opts.To.IsZero() && ts.After(opts.To) {
			continue
		}

		if _, err := fmt.Fprintln(w, line); err != nil {
			return written, fmt.Errorf("slicer: write error: %w", err)
		}
		written++
	}

	if err := scanner.Err(); err != nil {
		return written, fmt.Errorf("slicer: scan error: %w", err)
	}
	return written, nil
}
