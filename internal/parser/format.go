package parser

import (
	"fmt"
	"time"
)

// Known log timestamp formats tried in order during detection.
var knownFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05.999",
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
}

// DetectFormat inspects a sample log line and returns the timestamp format
// string that successfully parses the leading timestamp field, or an error
// if no known format matches.
func DetectFormat(sample string) (string, error) {
	if len(sample) == 0 {
		return "", fmt.Errorf("parser: empty sample line")
	}

	// Extract a candidate token — up to the first space after position 10.
	candidate := extractCandidate(sample)

	for _, fmt := range knownFormats {
		if _, err := time.Parse(fmt, candidate); err == nil {
			return fmt, nil
		}
	}

	// Try progressively longer prefixes to handle formats with embedded spaces.
	for end := 10; end <= len(sample) && end <= 35; end++ {
		token := sample[:end]
		for _, f := range knownFormats {
			if _, err := time.Parse(f, token); err == nil {
				return f, nil
			}
		}
	}

	return "", fmt.Errorf("parser: could not detect timestamp format in sample: %q", sample)
}

// extractCandidate returns the first whitespace-delimited token from line,
// assuming the timestamp starts at position 0.
func extractCandidate(line string) string {
	for i, ch := range line {
		if i > 0 && ch == ' ' {
			return line[:i]
		}
	}
	return line
}
