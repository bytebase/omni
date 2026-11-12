package compat

import (
	"strings"
	"testing"
)

func parse(t *testing.T, rows string) *File {
	t.Helper()
	f, err := Parse("compat.tsv", []byte(rows))
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestGuard(t *testing.T) {
	for _, c := range []struct {
		name          string
		before, after string
		pinBumped     bool
		want          string // a substring of the problem, or "" for none
	}{
		{"a closed gap", "reject\taccept\tengine\tr/A\tA\n", "accept\taccept\tengine\tr/A\tA\n", false, ""},
		{"a committed regression", "accept\taccept\tengine\tr/A\tA\n", "reject\taccept\tengine\tr/A\tA\n", false, "a regression"},
		{"a recording fills omni", "?\t?\t-\tissue:X\tA\n", "reject\taccept\tengine\tissue:X\tA\n", false, ""},
		{"an edited expected", "reject\taccept\tengine\tr/A\tA\n", "reject\treject\tengine\tr/A\tA\n", false, "without a pin bump"},
		{"an edited source", "reject\taccept\tengine\tr/A\tA\n", "reject\taccept\tgrammar\tr/A\tA\n", false, "without a pin bump"},
		{"a pin bump", "reject\taccept\tengine\tr/A\tA\n", "reject\treject\tengine\tr/A\tA\n", true, ""},
		{"a new row", "", "reject\taccept\tengine\tr/A\tA\n", false, ""},
		{"a removed row", "reject\taccept\tengine\tr/A\tA\n", "", false, ""},
	} {
		got := strings.Join(Guard(parse(t, c.before), parse(t, c.after), c.pinBumped), "\n")
		if c.want == "" && got != "" || !strings.Contains(got, c.want) {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}
