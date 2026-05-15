package cli

import (
	"fmt"
	"os"

	"github.com/user/logslice/internal/parser"
)

// RunUnique prints all unique values for a given field key found in the log
// file. It auto-detects the timestamp format from the file.
//
// Usage: logslice unique <file> <fieldKey>
func RunUnique(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: logslice unique <file> <fieldKey>")
	}

	filePath := args[0]
	fieldKey := args[1]

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unique: open file: %w", err)
	}
	defer f.Close()

	fmt, err := detectFormat(f)
	if err != nil {
		return fmt.Errorf("unique: detect format: %w", err)
	}

	// Rewind after format detection.
	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("unique: seek: %w", err)
	}

	values, err := parser.UniqueValues(f, fmt, fieldKey)
	if err != nil {
		return fmt.Errorf("unique: %w", err)
	}

	if len(values) == 0 {
		fmt.Fprintf(os.Stdout, "(no values found for field %q)\n", fieldKey)
		return nil
	}

	for _, v := range values {
		fmt.Fprintln(os.Stdout, v)
	}
	return nil
}
