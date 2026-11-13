package v1review

import "testing"

func TestV1ReservedKeywords(t *testing.T) {
	if got := len(v1ReservedKeywords); got != 484 {
		t.Errorf("v1ReservedKeywords has %d entries, want 484", got)
	}
	for name, want := range map[string]bool{
		"OVER":       true, // first entry
		"ANY":        true, // last entry
		"SELECT":     true,
		"select":     true, // the lookup upper-cases
		"Fillfactor": true,
		// Upper-casing is Unicode's: dotless ı and long ſ become I and S.
		"ınsert": true,
		"ſelect": true,
		// ODBC and future reserved words the parser accepts as names.
		"ABSOLUTE":      true,
		"XMLATTRIBUTES": true,
		// The two entries that are not one word.
		"END-EXEC":     true,
		"WITHIN GROUP": true,
		// Keywords to the parser, absent from the V1 table.
		"TRY_CAST":  false,
		"TRY_PARSE": false,
		"ASSEMBLY":  false,
		// Type names, absent too.
		"NVARCHAR": false,
		"BIGINT":   false,
		// A function name, in Bytebase's other table and not in this one.
		"USER_NAME": false,
		// Near misses.
		"":              false,
		"SELECTS":       false,
		"SELEC":         false,
		"SELECT ":       false,
		" SELECT":       false,
		"[SELECT]":      false,
		"#SELECT":       false,
		"END_EXEC":      false,
		"WITHIN  GROUP": false,
		"WITHIN_GROUP":  false,
	} {
		if got := isV1ReservedKeyword(name); got != want {
			t.Errorf("isV1ReservedKeyword(%q) = %v, want %v", name, got, want)
		}
	}
}
