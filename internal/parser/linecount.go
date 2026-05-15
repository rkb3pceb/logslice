package parser

import (
	"bufio"
	"io"
	"strings"
	"time"
)

// CountLines scans r and returns the number of log lines whose timestamps
// fall within [from, to). Passing an empty string for to means no upper bound.
//
// format must be a layout string as returned by DetectFormat.
func CountLines(r io.Reader, format, from, to string) (int, error) {
	start, err := ParseTimestamp(format, from)
	if err != nil {
		return 0, err
	}

	var endTime time.Time
	hasEnd := to != ""
	if hasEnd {
		endTime, err = ParseTimestamp(format, to)
		if err != nil {
			return 0, err
		}
	}

	var count int
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if isBlank(line) {
			continue
		}
		candidate := extractCandidate(line)
		if candidate == "" {
			continue
		}
		t, err := ParseTimestamp(format, candidate)
		if err != nil {
			continue
		}
		if t.Before(start) {
			continue
		}
		if hasEnd && !t.Before(endTime) {
			continue
		}
		count++
	}
	return count, scanner.Err()
}

// isBlank reports whether line contains only whitespace.
func isBlank(line string) bool {
	return strings.TrimSpace(line) == ""
}
