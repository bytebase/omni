package parser

import "testing"

// TestUpstreamLexerQuirks pins forms the wave-21 compat rows and their reject
// fuzz found, with TiDB v8.5.5's answers measured by running each.
func TestUpstreamLexerQuirks(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		accept bool
	}{
		// A @@ no name follows hands the next byte back as a token (bareSysVar).
		{"select tidb_distsql_scan_concurrency @@", true},
		{"select 1 @@", true},
		{"select 1 @@;", true},
		{"select 1 @@; select 2", true},
		{"select 1, 2 @@", true},
		{"do 1 @@", true},
		{"select 1 @@global.", true},
		{"select 1 @@-2", true},
		{"select @@-1", true},
		{"select @@`autocommit`", true},
		{"select @@'autocommit'", true},
		{"select @@global.autocommit", true},
		{"select 1 @@ ", false},
		{"select 1 @@ -- c", false},
		{"select 1 @@@@", false},
		{"select 1 @@`", false},
		{"select 1 @@.", false},
		{"select (1 @@)", false},
		{"select 1 @@x", false},
		{"select @@", false},
		{"select @@ autocommit", false},
		// VALUES() takes a column name, not an expression.
		{"DO VALUES(a)", true},
		{"DO VALUES(a.b.c)", true},
		{"DO VALUES(i*ent)", false},
		{"DO VALUES(a.*)", false},
		// A CHECK's NOT needs its ENFORCED.
		{"ALTER TABLE t ADD CHECK(a) NOT ENFORCED", true},
		{"ALTER TABLE t ADD CHECK(a) NOT", false},
	} {
		_, err := Parse(tc.sql)
		if got := err == nil; got != tc.accept {
			t.Errorf("%q: accepted %v, TiDB v8.5.5 %v (err %v)", tc.sql, got, tc.accept, err)
		}
	}
}
