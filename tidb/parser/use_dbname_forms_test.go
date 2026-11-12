package parser

import "testing"

// USE takes a DBName, which is an Identifier (parser.y:4440, 6774) and so
// admits no reserved keyword. TIDB_CURRENT_TSO is reserved on v8.5.5 (its
// ReservedKeyword %token block, parser.y:277) and reaches the grammar only as a
// function name (FunctionNameOptionalBraces, parser.y:8148), so `USE
// TIDB_CURRENT_TSO` is a 1064 while the function forms and a backquoted name
// parse. omni's MySQL-forked keyword table has no token for the word, so the
// lexer scans it as a bare identifier; the DBName gate in parseUseStmt is what
// rejects it. The engine verdict is pinned by the
// suite:cited rows of testdata/compat.tsv, measured against
// pingcap/tidb:v8.5.5.
var useDBNameFormCases = []struct {
	sql    string
	accept bool
}{
	// The reserved word is no bare DBName.
	{"USE TIDB_CURRENT_TSO", false},
	// A backquoted name of the same spelling is an ordinary identifier.
	{"USE `tidb_current_tso`", true},
	// An ordinary name and an unreserved keyword are both DBNames; only the
	// reserved word is refused, so USE keeps reading every real database name.
	{"USE somedb", true},
	{"USE SEND_CREDENTIALS_TO_TIKV", true},
	// The word stays a function name, read through the expression path, with or
	// without braces (FunctionNameOptionalBraces OptionalBraces).
	{"SELECT TIDB_CURRENT_TSO", true},
	{"SELECT TIDB_CURRENT_TSO()", true},
}

func TestUseDBNameForms(t *testing.T) {
	for _, c := range useDBNameFormCases {
		_, err := Parse(c.sql)
		if c.accept && err != nil {
			t.Errorf("Parse(%q) = %v, want accept", c.sql, err)
		}
		if !c.accept && err == nil {
			t.Errorf("Parse(%q) = accept, want reject", c.sql)
		}
	}
}
