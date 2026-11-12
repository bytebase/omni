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

const (
	pin1 = "# pin: default e:1 sha256:1 1\n"
	pin2 = "# pin: default e:2 sha256:2 2\n"
)

func TestGuard(t *testing.T) {
	for _, c := range []struct {
		name          string
		before, after string
		want          string // a substring of the problem, or "" for none
	}{
		{"a closed gap", "reject\taccept\tengine\tr/A\tA\n", "accept\taccept\tengine\tr/A\tA\n", ""},
		{"a committed regression", "accept\taccept\tengine\tr/A\tA\n", "reject\taccept\tengine\tr/A\tA\n", "a regression"},
		{"a recording fills omni", "?\t?\t-\tissue:X\tA\n", "reject\taccept\tengine\tissue:X\tA\n", ""},
		{"an edited expected", "reject\taccept\tengine\tr/A\tA\n", "reject\treject\tengine\tr/A\tA\n", "without a pin bump"},
		{"an edited source", "reject\taccept\tengine\tr/A\tA\n", "reject\taccept\tgrammar\tr/A\tA\n", "without a pin bump"},
		{"a pin bump", pin1 + "reject\taccept\tengine\tr/A\tA\n", pin2 + "reject\treject\tengine\tr/A\tA\n", ""},
		{"an edit at the same pin", pin1 + "reject\taccept\tengine\tr/A\tA\n", pin1 + "reject\treject\tengine\tr/A\tA\n", "without a pin bump"},
		{"a new row", "", "reject\taccept\tengine\tr/A\tA\n", ""},
		{"a removed row", "reject\taccept\tengine\tr/A\tA\n", "", ""},
	} {
		got := strings.Join(Guard(parse(t, c.before), parse(t, c.after)), "\n")
		if c.want == "" && got != "" || !strings.Contains(got, c.want) {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}
