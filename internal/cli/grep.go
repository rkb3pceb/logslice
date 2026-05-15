package cli

import (
	"fmt"
	"os"

	"github.com/user/logslice/internal/parser"
)

// RunGrep implements the `logslice grep` sub-command.
// args must be: <pattern> <file>
// Flags:
//
//	--case-sensitive  perform a case-sensitive search
func RunGrep(args []string, caseSensitive bool) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: logslice grep <pattern> <file>")
	}

	pattern := args[0]
	filePath := args[1]

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", filePath, err)
	}
	defer f.Close()

	fmt, err := detectFormat(f)
	if err != nil {
		return fmt.Errorf("detect format: %w", err)
	}

	// Rewind after format detection.
	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("seek: %w", err)
	}

	opts := parser.GrepOptions{
		Pattern:       pattern,
		Format:        fmt,
		CaseSensitive: caseSensitive,
	}

	results, err := parser.GrepLines(f, opts)
	if err != nil {
		return fmt.Errorf("grep: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintf(os.Stderr, "no matches found for %q\n", pattern)
		return nil
	}

	for _, res := range results {
		fmt.Printf("%d\t%s\n", res.LineNum, res.Line)
	}

	return nil
}
