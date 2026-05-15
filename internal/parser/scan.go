package parser

import (
	"bufio"
	"io"
	"time"
)

// LogLine represents a parsed line from a log file with its associated timestamp.
type LogLine struct {
	// Raw is the original line text.
	Raw string
	// Timestamp is the parsed time from the line, or zero if unparseable.
	Timestamp time.Time
	// LineNumber is the 1-based index of the line in the source.
	LineNumber int
}

// ScanLines reads all lines from r, parsing each timestamp using the given
// format string. Lines that cannot be parsed are included with a zero
// Timestamp. Blank lines are skipped entirely.
func ScanLines(r io.Reader, format string) ([]LogLine, error) {
	scanner := bufio.NewScanner(r)
	var lines []LogLine
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		if isBlank(raw) {
			continue
		}

		ll := LogLine{
			Raw:        raw,
			LineNumber: lineNum,
		}

		if format != "" {
			if candidate := extractCandidate(raw, format); candidate != "" {
				if t, err := time.Parse(format, candidate); err == nil {
					ll.Timestamp = t
				}
			}
		}

		lines = append(lines, ll)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// FilterByRange returns only those LogLines whose Timestamp falls within
// [from, to]. If to is zero, only the lower bound is applied.
func FilterByRange(lines []LogLine, from, to time.Time) []LogLine {
	var result []LogLine
	for _, ll := range lines {
		if ll.Timestamp.IsZero() {
			continue
		}
		if ll.Timestamp.Before(from) {
			continue
		}
		if !to.IsZero() && ll.Timestamp.After(to) {
			continue
		}
		result = append(result, ll)
	}
	return result
}
