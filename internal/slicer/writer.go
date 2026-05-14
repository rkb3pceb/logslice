// Package slicer provides functionality for extracting time-range segments
// from structured log files.
package slicer

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// WriteOptions configures the behaviour of WriteSlice.
type WriteOptions struct {
	// Format is the timestamp layout used to parse log lines.
	// If empty, DetectFormatFromReader is attempted before slicing.
	Format string

	// From is the inclusive lower bound of the time range.
	From time.Time

	// To is the inclusive upper bound of the time range.
	// A zero value means "no upper bound".
	To time.Time
}

// WriteSlice reads log lines from src, extracts those whose timestamps fall
// within the range defined by opts, and writes matching lines to dst.
// It returns the number of lines written and any error encountered.
func WriteSlice(dst io.Writer, src io.Reader, opts WriteOptions) (int, error) {
	if opts.Format == "" {
		return 0, fmt.Errorf("slicer: WriteOptions.Format must not be empty")
	}

	scanner := bufio.NewScanner(src)
	writer := bufio.NewWriter(dst)

	var written int

	for scanner.Scan() {
		line := scanner.Text()

		ts, err := parser.ParseTimestamp(line, opts.Format)
		if err != nil {
			// Skip lines we cannot parse a timestamp from.
			continue
		}

		if ts.Before(opts.From) {
			continue
		}
		if !opts.To.IsZero() && ts.After(opts.To) {
			break
		}

		if _, err := fmt.Fprintln(writer, line); err != nil {
			return written, fmt.Errorf("slicer: write error: %w", err)
		}
		written++
	}

	if err := scanner.Err(); err != nil {
		return written, fmt.Errorf("slicer: scan error: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return written, fmt.Errorf("slicer: flush error: %w", err)
	}

	return written, nil
}
