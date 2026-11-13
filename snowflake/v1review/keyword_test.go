package v1review

import "testing"

func TestV1Keywords(t *testing.T) {
	if got := len(v1Keywords); got != 91 {
		t.Errorf("v1Keywords has %d entries, want 91", got)
	}
	for name, want := range map[string]bool{
		"ACCOUNT": true, // first entry
		"WITH":    true, // last entry
		"SELECT":  true,
		"select":  true, // the lookup upper-cases
		"Connect": true,
		// In the V1 table though the parser does not reserve them.
		"GSCLUSTER": true,
		"ISSUE":     true,
		"WHENEVER":  true,
		"TRIGGER":   true,
		// Keywords to the parser, absent from the V1 table.
		"INDEX":     false,
		"LIMIT":     false,
		"PRIMARY":   false,
		"TIMESTAMP": false,
		// Near misses.
		"":         false,
		"SELECTS":  false,
		"SELEC":    false,
		"SELECT ":  false,
		`"SELECT"`: false,
		"TRY CAST": false,
	} {
		if got := isV1Keyword(name); got != want {
			t.Errorf("isV1Keyword(%q) = %v, want %v", name, got, want)
		}
	}
}
