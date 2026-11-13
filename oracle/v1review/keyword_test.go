package v1review

import "testing"

func TestIsKeyword(t *testing.T) {
	if got := len(oracleReservedWords); got != 85 {
		t.Errorf("oracleReservedWords has %d entries, want 85", got)
	}
	if got := len(oracleKeywords); got != 252 {
		t.Errorf("oracleKeywords has %d entries, want 252", got)
	}
	for word := range oracleReservedWords {
		if oracleKeywords[word] {
			t.Errorf("%q is in both tables", word)
		}
	}
	for _, tt := range []struct {
		text string
		want bool
	}{
		{"TABLE", true}, // reserved word
		{"table", true},
		{"Table", true},
		{"NAME", true}, // keyword
		{"name", true},
		{"A", true},
		{"AVG", true},
		{"COUNT", true},
		{"LEVEL", true},
		{"TIMESTAMP", true},
		{"", false},
		{"ANALYZE", false},
		{"PRIMARY", false},
		{"ROWNUM", false},
		{"NUMBER", false},
		{"VARCHAR2", false},
		{"TABLE ", false},
		{"TABLES", false},
		{"ID", false},
	} {
		if got := isKeyword(tt.text); got != tt.want {
			t.Errorf("isKeyword(%q) = %v, want %v", tt.text, got, tt.want)
		}
	}
}
