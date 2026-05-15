package parser

import (
	"io"
	"time"
)

// Stats holds aggregate statistics computed from a log file.
type Stats struct {
	// TotalLines is the number of non-blank lines scanned.
	TotalLines int
	// ParsedLines is the number of lines with a successfully parsed timestamp.
	ParsedLines int
	// SkippedLines is the number of lines whose timestamp could not be parsed.
	SkippedLines int
	// Earliest is the earliest timestamp found, zero if none.
	Earliest time.Time
	// Latest is the latest timestamp found, zero if none.
	Latest time.Time
	// Duration is the span between Earliest and Latest.
	Duration time.Duration
}

// CollectStats reads all lines from r, parses timestamps using format, and
// returns aggregate Stats. Blank lines are ignored. The reader is consumed
// entirely; callers should pass a fresh reader or reset before calling.
func CollectStats(r io.Reader, format string) (Stats, error) {
	if format == "" {
		return Stats{}, nil
	}

	lines, err := ScanLines(r)
	if err != nil {
		return Stats{}, err
	}

	var s Stats
	for _, line := range lines {
		if isBlank(line) {
			continue
		}
		s.TotalLines++

		candidate := extractCandidate(line, format)
		t, err := ParseTimestamp(candidate, format)
		if err != nil {
			s.SkippedLines++
			continue
		}
		s.ParsedLines++

		if s.Earliest.IsZero() || t.Before(s.Earliest) {
			s.Earliest = t
		}
		if s.Latest.IsZero() || t.After(s.Latest) {
			s.Latest = t
		}
	}

	if !s.Earliest.IsZero() && !s.Latest.IsZero() {
		s.Duration = s.Latest.Sub(s.Earliest)
	}

	return s, nil
}
