# logslice

Fast log file slicer that extracts time-range segments from large structured log files.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

## Usage

Extract log entries between two timestamps:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log
```

Pipe output to a file:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log > slice.log
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to the log file | required |
| `--from` | Start timestamp (RFC3339) | required |
| `--to` | End timestamp (RFC3339) | required |
| `--format` | Timestamp format in logs | `rfc3339` |
| `--field` | Timestamp field name | `time` |

### Example

```bash
# Extract the last hour of errors from a JSON log file
logslice \
  --file /var/log/app.log \
  --from "2024-01-15T12:00:00Z" \
  --to "2024-01-15T13:00:00Z" \
  --field timestamp
```

## How It Works

`logslice` uses binary search to efficiently locate the start of the requested time range, avoiding the need to scan the entire file. This makes it significantly faster than `grep` on large log files.

## License

MIT © 2024