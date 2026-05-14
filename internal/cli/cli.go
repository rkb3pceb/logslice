// Package cli provides the command-line interface for logslice.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/slicer"
)

// Config holds the parsed CLI options.
type Config struct {
	From   time.Time
	To     time.Time
	Input  string
	Format string
}

// Run parses args and executes the slice operation, writing results to out.
func Run(args []string, in io.Reader, out io.Writer, errOut io.Writer) error {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(errOut)

	var fromStr, toStr, inputPath, format string
	fs.StringVar(&fromStr, "from", "", "start timestamp (RFC3339 or detected format)")
	fs.StringVar(&toStr, "to", "", "end timestamp (RFC3339 or detected format), optional")
	fs.StringVar(&inputPath, "input", "", "path to log file (default: stdin)")
	fs.StringVar(&format, "format", "", "timestamp format override")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fromStr == "" {
		return errors.New("--from is required")
	}

	from, err := parser.ParseTimestamp(fromStr)
	if err != nil {
		return fmt.Errorf("invalid --from: %w", err)
	}

	var to time.Time
	if toStr != "" {
		to, err = parser.ParseTimestamp(toStr)
		if err != nil {
			return fmt.Errorf("invalid --to: %w", err)
		}
	}

	var r io.Reader = in
	if inputPath != "" {
		f, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("opening input: %w", err)
		}
		defer f.Close()
		r = f
	}

	if format == "" {
		format, r, err = detectFormat(r)
		if err != nil {
			return fmt.Errorf("detecting format: %w", err)
		}
	}

	return slicer.WriteSlice(r, out, from, to, format)
}

func detectFormat(r io.Reader) (string, io.Reader, error) {
	fmt, buffered, err := parser.DetectFormatFromReader(r)
	if err != nil {
		return "", nil, err
	}
	if fmt == "" {
		return "", nil, errors.New("could not detect timestamp format; use --format")
	}
	return fmt, buffered, nil
}
