package parser

import "testing"

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

// TestFromAliasIsNoString: a select item's alias may be a string, and an alias in FROM may
// not: table_alias reads an identifier, on a table, a derived table and a JSON_TABLE alike.
func TestFromAliasIsNoString(t *testing.T) {
	for _, sql := range []string{
		`SELECT * FROM t AS 'a'`,
		`SELECT * FROM t 'a'`,
		`SELECT * FROM t AS "a" JOIN u AS 'b' ON t.id = b.id`,
		`SELECT * FROM JSON_TABLE('[]', '$' COLUMNS (v INT PATH '$')) AS 'j'`,
		`SELECT * FROM (SELECT 1) AS 'a'`,
		`SELECT * FROM (SELECT 1) 'a'`,
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) took a string as an alias", sql)
		}
	}
}
