package parser

import (
	"bufio"
	"io"
	"strings"
	"time"
)

// GrepResult holds a matched log line along with its parsed timestamp.
type GrepResult struct {
	Line      string
	Timestamp time.Time
	LineNum   int
}

// GrepOptions configures the behaviour of GrepLines.
type GrepOptions struct {
	// Pattern is the substring or keyword to search for (case-insensitive when
	// CaseSensitive is false).
	Pattern string

	// CaseSensitive controls whether matching is case-sensitive.
	CaseSensitive bool

	// Format is the timestamp format string used to parse each line.
	Format string
}

// GrepLines reads from r and returns every log line whose text contains
// opts.Pattern. Lines that do not contain a parseable timestamp are skipped.
// An empty format string causes GrepLines to return an error immediately.
func GrepLines(r io.Reader, opts GrepOptions) ([]GrepResult, error) {
	if opts.Format == "" {
		return nil, ErrEmptyFormat
	}

	pattern := opts.Pattern
	if !opts.CaseSensitive {
		pattern = strings.ToLower(pattern)
	}

	var results []GrepResult
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		if strings.TrimSpace(raw) == "" {
			continue
		}

		ts, err := ParseTimestamp(raw, opts.Format)
		if err != nil {
			continue
		}

		haystack := raw
		if !opts.CaseSensitive {
			haystack = strings.ToLower(raw)
		}

		if strings.Contains(haystack, pattern) {
			results = append(results, GrepResult{
				Line:      raw,
				Timestamp: ts,
				LineNum:   lineNum,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
