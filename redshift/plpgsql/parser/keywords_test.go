package parser

import (
	"testing"
)

func TestPLReservedKeywords(t *testing.T) {
	reserved := []string{
		"ALL", "BEGIN", "BY", "CASE", "DECLARE", "ELSE", "END", "EXECUTE",
		"FOR", "FOREACH", "FROM", "IF", "IN", "INTO", "LOOP", "NOT",
		"NULL", "OR", "STRICT", "THEN", "TO", "USING", "WHEN", "WHILE",
	}
	for _, kw := range reserved {
		cat, ok := LookupPLKeyword(kw)
		if !ok {
			t.Errorf("expected %q to be a PL/pgSQL keyword", kw)
			continue
		}
		if cat != PLReserved {
			t.Errorf("expected %q to be PLReserved, got %d", kw, cat)
		}
	}
}

func TestPLUnreservedKeywords(t *testing.T) {
	unreserved := []string{
		"ABSOLUTE", "ALIAS", "AND", "ARRAY", "ASSERT", "BACKWARD",
		"CALL", "CHAIN", "CLOSE", "COLLATE", "COLUMN", "COLUMN_NAME",
		"COMMIT", "CONSTANT", "CONSTRAINT", "CONSTRAINT_NAME", "CONTINUE",
		"CURRENT", "CURSOR", "DATATYPE", "DEBUG", "DEFAULT", "DETAIL",
		"DIAGNOSTICS", "DO", "DUMP", "ELSIF", "ELSEIF", "ERRCODE", "ERROR",
		"EXCEPTION", "EXIT", "FETCH", "FIRST", "FORWARD", "GET", "HINT",
		"IMPORT", "INFO", "INSERT", "IS", "LAST", "LOG", "MERGE", "MESSAGE",
		"MESSAGE_TEXT", "MOVE", "NEXT", "NO", "NOTICE", "OPEN", "OPTION",
		"PERFORM", "PG_CONTEXT", "PG_DATATYPE_NAME", "PG_EXCEPTION_CONTEXT",
		"PG_EXCEPTION_DETAIL", "PG_EXCEPTION_HINT", "PG_ROUTINE_OID",
		"PRINT_STRICT_PARAMS", "PRIOR", "QUERY", "RAISE", "RELATIVE",
		"RETURN", "RETURNED_SQLSTATE", "REVERSE", "ROLLBACK", "ROW_COUNT",
		"ROWTYPE", "SCHEMA", "SCHEMA_NAME", "SCROLL", "SLICE", "SQLSTATE",
		"STACKED", "TABLE", "TABLE_NAME", "TYPE", "USE_COLUMN",
		"USE_VARIABLE", "VARIABLE_CONFLICT", "WARNING",
	}
	for _, kw := range unreserved {
		cat, ok := LookupPLKeyword(kw)
		if !ok {
			t.Errorf("expected %q to be a PL/pgSQL keyword", kw)
			continue
		}
		if cat != PLUnreserved {
			t.Errorf("expected %q to be PLUnreserved, got %d", kw, cat)
		}
	}
}

func TestLookupPLKeywordCaseInsensitive(t *testing.T) {
	tests := []struct {
		input   string
		wantCat PLKeywordCategory
		wantOK  bool
	}{
		{"begin", PLReserved, true},
		{"BEGIN", PLReserved, true},
		{"Begin", PLReserved, true},
		{"bEgIn", PLReserved, true},
		{"declare", PLReserved, true},
		{"DECLARE", PLReserved, true},
		{"return", PLUnreserved, true},
		{"RETURN", PLUnreserved, true},
		{"Return", PLUnreserved, true},
		{"exception", PLUnreserved, true},
		{"EXCEPTION", PLUnreserved, true},
		{"pg_context", PLUnreserved, true},
		{"PG_CONTEXT", PLUnreserved, true},
		{"variable_conflict", PLUnreserved, true},
		{"VARIABLE_CONFLICT", PLUnreserved, true},
		// Not a PL/pgSQL keyword.
		{"SELECT", 0, false},
		{"foobar", 0, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		cat, ok := LookupPLKeyword(tt.input)
		if ok != tt.wantOK {
			t.Errorf("LookupPLKeyword(%q): got ok=%v, want %v", tt.input, ok, tt.wantOK)
			continue
		}
		if ok && cat != tt.wantCat {
			t.Errorf("LookupPLKeyword(%q): got cat=%d, want %d", tt.input, cat, tt.wantCat)
		}
	}
}

func TestPLReservedNotUsableAsVariableName(t *testing.T) {
	// Reserved keywords should have PLReserved category, meaning they
	// cannot be used as variable names.
	cat, ok := LookupPLKeyword("BEGIN")
	if !ok || cat != PLReserved {
		t.Error("BEGIN should be PLReserved")
	}
	cat, ok = LookupPLKeyword("END")
	if !ok || cat != PLReserved {
		t.Error("END should be PLReserved")
	}
}

func TestPLUnreservedUsableAsVariableName(t *testing.T) {
	// Unreserved keywords can be used as variable names.
	cat, ok := LookupPLKeyword("CURSOR")
	if !ok || cat != PLUnreserved {
		t.Error("CURSOR should be PLUnreserved")
	}
	cat, ok = LookupPLKeyword("RAISE")
	if !ok || cat != PLUnreserved {
		t.Error("RAISE should be PLUnreserved")
	}
}

func TestPLKeywordCounts(t *testing.T) {
	if got := len(PLReservedKeywords); got != 24 {
		t.Errorf("expected 24 reserved keywords, got %d", got)
	}
	if got := len(PLUnreservedKeywords); got != 83 {
		t.Errorf("expected 83 unreserved keywords, got %d", got)
	}
}
