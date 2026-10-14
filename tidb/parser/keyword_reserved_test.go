package parser

import "testing"

func TestIsReservedKeyword(t *testing.T) {
	for _, tc := range []struct {
		word string
		want bool
	}{
		{"select", true},
		{"SELECT", true},
		{"table", true},
		{"status", false},
		{"orders", false},
	} {
		if got := IsReservedKeyword(tc.word); got != tc.want {
			t.Errorf("IsReservedKeyword(%q) = %v, want %v", tc.word, got, tc.want)
		}
	}
}
