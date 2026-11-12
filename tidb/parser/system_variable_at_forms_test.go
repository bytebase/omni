package parser

import "testing"

// A system variable's name after its `@@`, every verdict measured with PREPARE on
// pingcap/tidb:v8.5.5 on 2026-11-12.
//
// Upstream's scanner reads `@@`, an optional `global.`, `session.` or `local.`, and
// the name as one token (startWithAt, pkg/parser/lexer.go). A name is an identifier or
// a quoted string, and neither starts with `@`: at a third `@` the scanner gives the
// token up and leaves a lone `@`, which no production takes. omni's lexer scans `@a`
// as a user variable, so parseSysVarName read it as the name.
//
// The reject-fuzz sweep found the first of these, by doubling an `@`.

var systemVariableAtCases = []tableSampleCase{
	{"select @@@tidb_index_lookup_pushdown_policy", false},
	{"select @@@autocommit", false},
	{"select @@@a", false},
	{"select @@@", false},
	{"select @@@a.b", false},
	{"select @@@`autocommit`", false},
	{"select @@@'autocommit'", false},
	{"select @@@a from t", false},
	{"select 1 + @@@a", false},
	{"select @@autocommit, @@@a", false},
	// After a scope, where the name starts again.
	{"select @@global.@autocommit", false},
	{"select @@session.@a", false},
	{"select @@local.@a", false},
	{"select @@global.@`a`", false},
	{"select @@@global.autocommit", false},
	{"select @@global.@@autocommit", false},
	// As an assignment's target.
	{"set @@@autocommit = 1", false},
	{"set @@global.@autocommit = 1", false},
	{"set @@session.@a = 1", false},
	{"set @@@a := 1", false},
	// The names it does take, and the user variable beside it.
	{"select @@autocommit", true},
	{"select @@`autocommit`", true},
	{"select @@'autocommit'", true},
	{"select @@global.autocommit", true},
	{"select @@global.`autocommit`", true},
	{"select @@global.'autocommit'", true},
	{"set @@autocommit = 1", true},
	{"select @a", true},
	{"select @@autocommit, @a", true},
	{"show variables like '@@@a'", true},
}

// TestSystemVariableAtForms holds omni to the recorded verdicts without a server.
func TestSystemVariableAtForms(t *testing.T) {
	for _, c := range systemVariableAtCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}
