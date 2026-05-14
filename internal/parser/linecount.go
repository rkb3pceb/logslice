// Package parser provides utilities for parsing and detecting timestamp
// formats in structured log files.
package parser

import (
	"bufio"
	"io"
	"strings"
)

// LineStats holds summary information about a scanned log file.
type LineStats struct {
	// Total is the total number of lines (including blank lines).
	Total int
	// Timestamped is the number of lines that contain a parseable timestamp.
	Timestamped int
	// Blank is the number of blank or whitespace-only lines.
	Blank int
}

// CountLines scans r and returns statistics about line composition.
// format is the timestamp format string used to test each line; if format is
// empty every non-blank line is counted as non-timestamped.
func CountLines(r io.Reader, format string) (LineStats, error) {
	var stats LineStats
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		stats.Total++
		if isBlank(line) {
			stats.Blank++
			continue
		}
		if format != "" {
			candidate := extractCandidate(line)
			if _, err := ParseTimestamp(format, candidate); err == nil {
				stats.Timestamped++
			}
		}
	}
	return stats, scanner.Err()
}

// isBlank reports whether s is empty or contains only whitespace.
func isBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}
