package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/logslice/internal/parser"
)

// RunLevel filters a log file by one or more severity levels and writes
// matching lines to stdout.
//
// Usage: logslice level <file> <LEVEL[,LEVEL...]>
func RunLevel(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: logslice level <file> <LEVEL[,LEVEL,...]>")
	}

	filePath := args[0]
	levelTokens := strings.Split(strings.ToUpper(args[1]), ",")

	levels := make([]parser.LogLevel, 0, len(levelTokens))
	for _, tok := range levelTokens {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		levels = append(levels, parser.LogLevel(tok))
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

	return parser.FilterByLevel(f, os.Stdout, format, levels)
}
