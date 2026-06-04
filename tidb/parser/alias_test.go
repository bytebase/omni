package parser

import "testing"

// TiDB accepts a quoted string as an alias ONLY in the SELECT list (column
// alias), matching its grammar (select_alias_ident -> ident | TEXT_STRING).
// Unlike MySQL, TiDB rejects string-literal aliases on table refs, derived
// tables, JSON_TABLE, and UPDATE/DELETE targets — verified on TiDB v8.5.0.

func TestSelectAliasStringLiteral(t *testing.T) {
	cases := []string{
		`SELECT 1 AS 'start'`,
		`SELECT 1 AS "start"`,
		`SELECT 1 'start'`,
		`SELECT a AS 'label', b FROM t`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}

// TestNonSelectAliasRejectsStringLiteral pins that string-literal aliases are
// rejected everywhere TiDB rejects them — omni must not parse SQL TiDB fails.
func TestNonSelectAliasRejectsStringLiteral(t *testing.T) {
	cases := []string{
		`SELECT * FROM t AS 'a'`,
		`SELECT * FROM t 'a'`,
		`SELECT * FROM (SELECT 1) AS 'a'`,
		`SELECT * FROM JSON_TABLE('[]', '$' COLUMNS (v INT PATH '$')) AS 'a'`,
		`UPDATE t AS 'a' SET a = 1`,
		`DELETE FROM t AS 'a' WHERE a = 1`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) accepted a string-literal alias, but TiDB rejects it", sql)
			}
		})
	}
}
