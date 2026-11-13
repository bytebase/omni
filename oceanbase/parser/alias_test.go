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

// A table's alias is a name and never a string: tbl_name's alias is relation_name, where
// MySQL's is ident_or_text. TestTableNameClausesOceanBaseOracle holds the first two to
// OceanBase 4.3.5.6, which answers 1064 for each.
func TestTableAliasStringLiteral(t *testing.T) {
	cases := []string{
		`SELECT * FROM t AS 'a'`,
		`SELECT * FROM t 'a'`,
		`SELECT * FROM t AS "a" JOIN u AS 'b' ON t.id = b.id`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) accepted a string as a table alias", sql)
			}
		})
	}
}

func TestJsonTableAliasStringLiteral(t *testing.T) {
	sql := `SELECT * FROM JSON_TABLE('[]', '$' COLUMNS (v INT PATH '$')) AS 'j'`
	if _, err := Parse(sql); err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
}

func TestDerivedTableAliasStringLiteral(t *testing.T) {
	cases := []string{
		`SELECT * FROM (SELECT 1) AS 'a'`,
		`SELECT * FROM (SELECT 1) 'a'`,
		`SELECT * FROM LATERAL (SELECT 1) AS 'a'`,
		`SELECT * FROM LATERAL (SELECT 1) 'a'`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}
