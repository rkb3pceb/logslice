package parser

import (
	"bufio"
	"io"
	"strings"
)

// LogLevel represents a severity level found in a log line.
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
	LevelFatal LogLevel = "FATAL"
	LevelUnknown LogLevel = "UNKNOWN"
)

// knownLevels lists all recognised severity tokens in priority order.
var knownLevels = []LogLevel{
	LevelFatal,
	LevelError,
	LevelWarn,
	LevelInfo,
	LevelDebug,
}

// DetectLevel scans a single log line and returns the first recognised
// severity level token. Returns LevelUnknown when none is found.
func DetectLevel(line string) LogLevel {
	upper := strings.ToUpper(line)
	for _, lvl := range knownLevels {
		if strings.Contains(upper, string(lvl)) {
			return lvl
		}
	}
	return LevelUnknown
}

// FilterByLevel reads lines from r and writes to w only those whose detected
// severity level matches one of the requested levels. An empty levels slice
// passes all lines through.
func FilterByLevel(r io.Reader, w io.Writer, format string, levels []LogLevel) error {
	if format == "" {
		return ErrEmptyFormat
	}

	want := make(map[LogLevel]bool, len(levels))
	for _, l := range levels {
		want[l] = true
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if isBlank(line) {
			continue
		}
		if len(want) == 0 || want[DetectLevel(line)] {
			if _, err := io.WriteString(w, line+"\n"); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}
