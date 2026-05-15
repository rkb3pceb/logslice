package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

// ContextLine holds a matched log line along with its surrounding context lines.
type ContextLine struct {
	Line    string
	Matched bool
	Time    time.Time
}

// GrepWithContext searches log lines for a pattern and returns matched lines
// along with up to contextN lines before and after each match.
// Lines that cannot be parsed with the given format are skipped for timestamp
// purposes but still included in context output.
func GrepWithContext(r io.Reader, format, pattern string, contextN int, caseSensitive bool) ([]ContextLine, error) {
	if format == "" {
		return nil, fmt.Errorf("format must not be empty")
	}
	if contextN < 0 {
		contextN = 0
	}

	var all []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if isBlank(line) {
			continue
		}
		all = append(all, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning: %w", err)
	}

	search := pattern
	if !caseSensitive {
		search = strings.ToLower(pattern)
	}

	// Mark which lines match.
	matched := make([]bool, len(all))
	for i, line := range all {
		candidate := line
		if !caseSensitive {
			candidate = strings.ToLower(line)
		}
		if strings.Contains(candidate, search) {
			matched[i] = true
		}
	}

	// Collect indices to include (match ± contextN).
	include := make([]bool, len(all))
	for i, m := range matched {
		if !m {
			continue
		}
		start := i - contextN
		if start < 0 {
			start = 0
		}
		end := i + contextN
		if end >= len(all) {
			end = len(all) - 1
		}
		for j := start; j <= end; j++ {
			include[j] = true
		}
	}

	var results []ContextLine
	for i, line := range all {
		if !include[i] {
			continue
		}
		t, _ := ParseTimestamp(line, format)
		results = append(results, ContextLine{
			Line:    line,
			Matched: matched[i],
			Time:    t,
		})
	}
	return results, nil
}
