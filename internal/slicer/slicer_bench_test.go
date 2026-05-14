package slicer_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/slicer"
)

// generateLog creates a synthetic log with n lines spanning one hour.
func generateLog(n int) io.Reader {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		ts := base.Add(time.Duration(i) * time.Second)
		fmt.Fprintf(&buf, "%s INFO  benchmark log line number %d\n",
			ts.Format("2006-01-02T15:04:05"), i)
	}
	return &buf
}

func BenchmarkSlice_10k(b *testing.B) {
	const n = 10_000
	from := time.Date(2024, 1, 1, 0, 10, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 20, 0, 0, time.UTC)

	for b.Loop() {
		r := generateLog(n)
		_, err := slicer.Slice(r, io.Discard, slicer.Options{
			From:   from,
			To:     to,
			Layout: "2006-01-02T15:04:05",
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
