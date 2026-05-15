package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// UniqueValues extracts all unique values for a given field key from structured
// log lines. Lines where the key is not found are silently skipped.
//
// fieldKey is matched as a simple substring prefix in the form "key=" or
// "key:", covering common structured log formats.
//
// Results are returned in first-seen order.
func UniqueValues(r io.Reader, format, fieldKey string) ([]string, error) {
	if format == "" {
		return nil, fmt.Errorf("unique: format must not be empty")
	}
	if fieldKey == "" {
		return nil, fmt.Errorf("unique: fieldKey must not be empty")
	}

	seen := make(map[string]struct{})
	var ordered []string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if isBlank(line) {
			continue
		}
		val, ok := extractFieldValue(line, fieldKey)
		if !ok {
			continue
		}
		if _, exists := seen[val]; !exists {
			seen[val] = struct{}{}
			ordered = append(ordered, val)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("unique: scan error: %w", err)
	}
	return ordered, nil
}

// extractFieldValue searches for fieldKey= or fieldKey: in a log line and
// returns the associated value token (space- or comma-delimited).
func extractFieldValue(line, key string) (string, bool) {
	for _, sep := range []string{key + "=", key + ":"} {
		idx := strings.Index(line, sep)
		if idx == -1 {
			continue
		}
		rest := line[idx+len(sep):]
		// strip optional leading quote
		rest = strings.TrimPrefix(rest, "\"")
		// value ends at space, comma, quote, or end of string
		end := strings.IndexAny(rest, " ,\t\"")
		if end == -1 {
			return strings.TrimSpace(rest), true
		}
		return strings.TrimSpace(rest[:end]), true
	}
	return "", false
}
