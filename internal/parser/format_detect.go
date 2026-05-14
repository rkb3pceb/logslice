package parser

import (
	"bufio"
	"io"
	"strings"
)

// SampleSize is the number of lines read from a log source to detect its format.
const SampleSize = 10

// DetectFormatFromReader reads up to SampleSize lines from r and attempts to
// detect the timestamp format used in the log file. It returns the detected
// format string (compatible with time.Parse) and the sampled lines so callers
// can replay them without re-reading the source.
//
// If no recognisable timestamp format is found, an empty string is returned.
func DetectFormatFromReader(r io.Reader) (format string, sample []string, err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() && len(sample) < SampleSize {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		sample = append(sample, line)
	}
	if err = scanner.Err(); err != nil {
		return "", sample, err
	}

	// Try each collected line until we get a hit.
	for _, line := range sample {
		if f := DetectFormat(line); f != "" {
			return f, sample, nil
		}
	}
	return "", sample, nil
}

// DetectFormatFromLines attempts to detect the timestamp format from a slice
// of log lines. It iterates through the provided lines and returns the first
// recognisable format string (compatible with time.Parse), or an empty string
// if none is found. This is useful when lines have already been read into
// memory and a separate reader is not available.
func DetectFormatFromLines(lines []string) string {
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if f := DetectFormat(line); f != "" {
			return f
		}
	}
	return ""
}
