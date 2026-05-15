// Package parser provides utilities for detecting timestamp formats in log
// files, parsing individual timestamps, scanning log lines into structured
// representations, and summarising log file metadata.
//
// # Format Detection
//
// DetectFormat inspects a sample string and returns the Go time layout that
// best matches the timestamps found within it. DetectFormatFromReader and
// DetectFormatFromLines offer higher-level helpers that sample multiple lines
// before deciding on a format.
//
// # Scanning
//
// ScanLines reads an [io.Reader] line-by-line and returns a slice of LogLine
// values, each carrying the raw text, its 1-based line number, and the parsed
// timestamp (zero if the line could not be parsed). FilterByRange then narrows
// that slice to a caller-supplied time window.
//
// # Counting and Summarising
//
// CountLines counts the log lines that fall within a time range without
// materialising the full slice. Summarize returns a [Summary] describing the
// first timestamp, last timestamp, and total line count of a log stream.
package parser
