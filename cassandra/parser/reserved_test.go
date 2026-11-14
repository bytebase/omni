package parser

import "testing"

// TestReservedWords holds Parser.g's reservation: a keyword of Lexer.g that
// unreserved_keyword (Parser.g:1946) leaves out is never a name unless
// quoted, an unreserved one is, and a function may be named TOKEN or COUNT
// (allowedFunctionName, Parser.g:1613) but not TTL or WRITETIME.
func TestReservedWords(t *testing.T) {
	for _, sql := range []string{
		`USE "keyspace"`,
		"USE ascii",
		"USE users",
		"DROP TABLE IF EXISTS ks.t",
		"SELECT ttl, count, json FROM t",
		"DROP FUNCTION token",
		"DROP FUNCTION ks.count",
		"CREATE AGGREGATE a(int) SFUNC token STYPE int",
		"GRANT EXECUTE ON FUNCTION ks.token(int) TO r",
		"ALTER TABLE t ALTER v MASKED WITH ks.token()",
		"ALTER TABLE t ALTER v MASKED WITH f(token(1), 'a', -1, [1], ?)",
		"CREATE TABLE t(k int PRIMARY KEY, v text MASKED WITH count())",
		"SELECT token(k), ks.count(v) FROM t",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{
		"USE keyspace",
		"USE schema",
		"USE table",
		"DROP TABLE IF",
		"TRUNCATE token",
		"SELECT keyspace FROM t",
		"CREATE ROLE keyspace",
		"DROP FUNCTION ttl",
		"DROP FUNCTION ks.writetime",
		"ALTER TABLE t ALTER v MASKED WITH ttl()",
		"ALTER TABLE t ALTER v MASKED WITH f(c)", // a column is no term (Parser.g:805)
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) = nil error, want one", sql)
		}
	}
}
