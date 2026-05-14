package parser

import (
	"bufio"
	"io"
	"strings"
)

// LineCountResult holds the result of counting lines in a log file.
type LineCountResult struct {
	// Total is the total number of lines scanned.
	Total int
	// Matched is the number of lines that contain a parseable timestamp.
	Matched int
	// Blank is the number of blank or whitespace-only lines.
	Blank int
}

// CountLines scans r using the provided format string and returns a
// LineCountResult describing the composition of the log data.
// If format is empty, CountLines returns an error.
func CountLines(r io.Reader, format string) (LineCountResult, error) {
	if format == "" {
		return LineCountResult{}, ErrNoFormat
	}

	var result LineCountResult
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		result.Total++
		if isBlank(line) {
			result.Blank++
			continue
		}
		candidate := extractCandidate(line, format)
		if candidate == "" {
			continue
		}
		_, err := ParseTimestamp(candidate, format)
		if err == nil {
			result.Matched++
		}
	}
	if err := scanner.Err(); err != nil {
		return result, err
	}
	return result, nil
}

// isBlank returns true when line contains only whitespace characters.
func isBlank(line string) bool {
	return strings.TrimSpace(line) == ""
}
