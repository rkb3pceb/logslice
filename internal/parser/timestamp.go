package parser

import (
	"fmt"
	"time"
)

// CommonFormats lists timestamp formats commonly found in log files.
var CommonFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
}

// ParseTimestamp attempts to parse a timestamp string using known formats.
// It returns the parsed time and the matched format string, or an error if
// no format matched.
func ParseTimestamp(s string) (time.Time, string, error) {
	for _, layout := range CommonFormats {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, layout, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("parser: unrecognized timestamp format: %q", s)
}

// MustParseTimestamp is like ParseTimestamp but panics on error.
// Intended for use in tests and CLI argument parsing.
func MustParseTimestamp(s string) time.Time {
	t, _, err := ParseTimestamp(s)
	if err != nil {
		panic(err)
	}
	return t
}
