package parser

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Bucket represents a time bucket with a line count.
type Bucket struct {
	Start time.Time
	Count int
}

// Histogram holds bucketed log line counts over a time range.
type Histogram struct {
	Buckets  []Bucket
	Interval time.Duration
}

// String returns a simple text representation of the histogram.
func (h *Histogram) String() string {
	if len(h.Buckets) == 0 {
		return "(empty histogram)"
	}
	out := fmt.Sprintf("Histogram (interval=%s)\n", h.Interval)
	for _, b := range h.Buckets {
		bar := ""
		for i := 0; i < b.Count && i < 40; i++ {
			bar += "#"
		}
		out += fmt.Sprintf("  %s  %s (%d)\n", b.Start.Format(time.RFC3339), bar, b.Count)
	}
	return out
}

// BuildHistogram reads lines from r using the given format and groups them
// into time buckets of the specified interval.
func BuildHistogram(r io.Reader, format string, interval time.Duration) (*Histogram, error) {
	if format == "" {
		return nil, fmt.Errorf("histogram: format must not be empty")
	}
	if interval <= 0 {
		return nil, fmt.Errorf("histogram: interval must be positive")
	}

	lines, err := ScanLines(r, format)
	if err != nil {
		return nil, fmt.Errorf("histogram: scan error: %w", err)
	}

	counts := make(map[time.Time]int)
	for _, l := range lines {
		bucket := l.Time.Truncate(interval)
		counts[bucket]++
	}

	buckets := make([]Bucket, 0, len(counts))
	for t, c := range counts {
		buckets = append(buckets, Bucket{Start: t, Count: c})
	}
	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].Start.Before(buckets[j].Start)
	})

	return &Histogram{Buckets: buckets, Interval: interval}, nil
}
