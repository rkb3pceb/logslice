// Package slicer provides the core log-slicing functionality for logslice.
//
// It reads structured log lines from an [io.Reader], parses each line's
// leading timestamp using the layout supplied in [Options], and writes only
// those lines whose timestamps fall within the requested [From, To] range to
// an [io.Writer].
//
// Lines that do not begin with a parseable timestamp are silently skipped,
// making the slicer tolerant of multi-line stack traces and other non-timestamped
// continuation lines.
//
// Example:
//
//	opts := slicer.Options{
//		From:   from,
//		To:     to,
//		Layout: "2006-01-02T15:04:05",
//	}
//	n, err := slicer.Slice(logFile, os.Stdout, opts)
package slicer
