package parser

import (
	"bufio"
	"io"
	"time"
)

// TailResult holds the last N parsed log entries.
type TailResult struct {
	Lines []string
	First time.Time
	Last  time.Time
}

// TailLines reads all lines from r and returns the last n lines that
// contain a parseable timestamp using the given format. Blank lines are
// skipped. If format is empty, TailLines returns an error.
func TailLines(r io.Reader, format string, n int) (TailResult, error) {
	if format == "" {
		return TailResult{}, ErrEmptyFormat
	}
	if n <= 0 {
		return TailResult{}, nil
	}

	// ring buffer of size n
	ring := make([]string, n)
	times := make([]time.Time, n)
	pos := 0
	count := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if isBlank(line) {
			continue
		}
		candidate := extractCandidate(line, format)
		if candidate == "" {
			continue
		}
		t, err := ParseTimestamp(candidate, format)
		if err != nil {
			continue
		}
		ring[pos%n] = line
		times[pos%n] = t
		pos++
		count++
	}
	if err := scanner.Err(); err != nil {
		return TailResult{}, err
	}

	if count == 0 {
		return TailResult{}, nil
	}

	size := count
	if size > n {
		size = n
	}

	result := make([]string, size)
	resultTimes := make([]time.Time, size)
	start := 0
	if count > n {
		start = pos % n
	}
	for i := 0; i < size; i++ {
		result[i] = ring[(start+i)%n]
		resultTimes[i] = times[(start+i)%n]
	}

	return TailResult{
		Lines: result,
		First: resultTimes[0],
		Last:  resultTimes[size-1],
	}, nil
}
