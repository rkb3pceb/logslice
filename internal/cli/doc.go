// Package cli implements the command-line interface for logslice.
//
// It parses flags, resolves the input source (file or stdin), auto-detects
// the timestamp format when not explicitly provided, and delegates the actual
// slicing work to the slicer and parser packages.
//
// Usage:
//
//	logslice --from <timestamp> [--to <timestamp>] [--input <file>] [--format <layout>]
//
// Flags:
//
//	--from    Required. Start of the time range to extract.
//	--to      Optional. End of the time range (inclusive). Omit to read to EOF.
//	--input   Optional. Path to the log file. Defaults to stdin.
//	--format  Optional. Go time layout string. Auto-detected from the first
//	          non-blank lines of the input when omitted.
package cli
