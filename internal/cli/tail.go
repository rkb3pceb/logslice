package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/user/logslice/internal/parser"
)

// RunTail prints the last N timestamped lines from a log file.
// Args: <file> <n>
func RunTail(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: logslice tail <file> <n>")
	}

	filePath := args[0]
	n, err := strconv.Atoi(args[1])
	if err != nil || n <= 0 {
		return fmt.Errorf("invalid count %q: must be a positive integer", args[1])
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", filePath, err)
	}
	defer f.Close()

	format, err := detectFormat(f)
	if err != nil {
		return fmt.Errorf("detect format: %w", err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("seek: %w", err)
	}

	res, err := parser.TailLines(f, format, n)
	if err != nil {
		return fmt.Errorf("tail: %w", err)
	}

	if len(res.Lines) == 0 {
		fmt.Println("(no timestamped lines found)")
		return nil
	}

	for _, line := range res.Lines {
		fmt.Println(line)
	}
	fmt.Fprintf(os.Stderr, "# %d lines  [%s .. %s]\n",
		len(res.Lines),
		res.First.Format(format),
		res.Last.Format(format),
	)
	return nil
}
