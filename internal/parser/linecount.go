package parser

import (
	"bufio"
	"io"
)

// LineCountResult holds statistics about lines scanned in a log file.
type LineCountResult struct {
	// Total is the total number of lines in the file.
	Total int
	// Blank is the number of blank (empty or whitespace-only) lines.
	Blank int
	// WithTimestamp is the number of lines that contain a parseable timestamp
	// using the given format.
	WithTimestamp int
}

// CountLines scans r line by line and returns a LineCountResult.
// format is the timestamp format string (as accepted by ParseTimestamp).
// If format is empty, WithTimestamp will always be zero.
func CountLines(r io.Reader, format string) (LineCountResult, error) {
	var result LineCountResult

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		result.Total++

		if isBlank(line) {
			result.Blank++
			continue
		}

		if format != "" {
			candidate := extractCandidate(line)
			if _, err := ParseTimestamp(candidate, format); err == nil {
				result.WithTimestamp++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}

// isBlank reports whether s contains only whitespace characters.
func isBlank(s string) bool {
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\r' {
			return false
		}
	}
	return true
}
