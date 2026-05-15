package parser

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

// Summary holds basic statistics about a log file's time range.
type Summary struct {
	// First is the timestamp of the earliest log line found.
	First time.Time
	// Last is the timestamp of the latest log line found.
	Last time.Time
	// LineCount is the total number of lines with a parseable timestamp.
	LineCount int
}

// String returns a human-readable summary.
func (s Summary) String() string {
	return fmt.Sprintf("lines=%d first=%s last=%s",
		s.LineCount, s.First.Format(time.RFC3339), s.Last.Format(time.RFC3339))
}

// Summarize scans r using format and returns a Summary of the log.
// Lines without a parseable timestamp are skipped.
func Summarize(r io.Reader, format string) (Summary, error) {
	var s Summary
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
		if s.LineCount == 0 || t.Before(s.First) {
			s.First = t
		}
		if t.After(s.Last) {
			s.Last = t
		}
		s.LineCount++
	}
	return s, scanner.Err()
}
