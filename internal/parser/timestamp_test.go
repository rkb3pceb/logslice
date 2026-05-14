package parser_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

type tsCase struct {
	input    string
	wantYear int
	wantErr  bool
}

func TestParseTimestamp(t *testing.T) {
	cases := []tsCase{
		{"2024-03-15T12:34:56Z", 2024, false},
		{"2024-03-15T12:34:56.123456789Z", 2024, false},
		{"2024-03-15 12:34:56", 2024, false},
		{"2024-03-15 12:34:56.000", 2024, false},
		{"15/Mar/2024:12:34:56 +0000", 2024, false},
		{"not-a-timestamp", 0, true},
		{"", 0, true},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, layout, err := parser.ParseTimestamp(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tc.input, err)
			}
			if got.Year() != tc.wantYear {
				t.Errorf("year: got %d, want %d", got.Year(), tc.wantYear)
			}
			if layout == "" {
				t.Errorf("expected non-empty layout for input %q", tc.input)
			}
		})
	}
}

func TestMustParseTimestamp_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid timestamp, got none")
		}
	}()
	parser.MustParseTimestamp("not-valid")
}

func TestMustParseTimestamp_Valid(t *testing.T) {
	got := parser.MustParseTimestamp("2024-06-01T00:00:00Z")
	if got.IsZero() {
		t.Error("expected non-zero time")
	}
	if got.Year() != 2024 {
		t.Errorf("expected year 2024, got %d", got.Year())
	}
	if got.UTC().Month() != time.June {
		t.Errorf("expected June, got %v", got.Month())
	}
}
