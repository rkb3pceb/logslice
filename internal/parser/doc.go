// Package parser provides utilities for detecting and parsing timestamps
// embedded in structured log lines.
//
// Log files produced by different systems use a wide variety of timestamp
// formats. This package centralises format detection so that the rest of
// logslice can work with standard time.Time values regardless of the
// original source format.
//
// Supported formats include:
//
//	- RFC 3339 / ISO 8601  (e.g. 2006-01-02T15:04:05Z)
//	- RFC 3339 with nanoseconds
//	- Space-separated datetime (e.g. 2006-01-02 15:04:05)
//	- Apache / nginx combined log format  (e.g. 02/Jan/2006:15:04:05 -0700)
//	- Syslog-style  (e.g. Jan 02 15:04:05)
//
// To add support for a new format, append its time.Parse layout string to
// CommonFormats before calling ParseTimestamp.
package parser
