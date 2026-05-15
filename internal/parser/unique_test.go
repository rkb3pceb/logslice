package parser

import (
	"strings"
	"testing"
)

func makeUniqueLog() string {
	return strings.Join([]string{
		`2024-01-01T00:00:01Z level=info method=GET path=/health`,
		`2024-01-01T00:00:02Z level=warn method=POST path=/login`,
		`2024-01-01T00:00:03Z level=info method=GET path=/health`,
		`2024-01-01T00:00:04Z level=error method=DELETE path=/item`,
		`2024-01-01T00:00:05Z level=warn method=POST path=/login`,
		``,
		`2024-01-01T00:00:06Z level=info method=GET path=/status`,
	}, "\n")
}

func TestUniqueValues_Basic(t *testing.T) {
	r := strings.NewReader(makeUniqueLog())
	got, err := UniqueValues(r, "2006-01-02T15:04:05Z", "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"info", "warn", "error"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("index %d: got %q, want %q", i, got[i], v)
		}
	}
}

func TestUniqueValues_Method(t *testing.T) {
	r := strings.NewReader(makeUniqueLog())
	got, err := UniqueValues(r, "2006-01-02T15:04:05Z", "method")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 unique methods, got %d: %v", len(got), got)
	}
}

func TestUniqueValues_EmptyFormat(t *testing.T) {
	r := strings.NewReader(makeUniqueLog())
	_, err := UniqueValues(r, "", "level")
	if err == nil {
		t.Fatal("expected error for empty format")
	}
}

func TestUniqueValues_EmptyFieldKey(t *testing.T) {
	r := strings.NewReader(makeUniqueLog())
	_, err := UniqueValues(r, "2006-01-02T15:04:05Z", "")
	if err == nil {
		t.Fatal("expected error for empty fieldKey")
	}
}

func TestUniqueValues_NoMatch(t *testing.T) {
	r := strings.NewReader(makeUniqueLog())
	got, err := UniqueValues(r, "2006-01-02T15:04:05Z", "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected no results, got %v", got)
	}
}

func TestUniqueValues_EmptyReader(t *testing.T) {
	r := strings.NewReader("")
	got, err := UniqueValues(r, "2006-01-02T15:04:05Z", "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}
