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

func TestTableAliasStringLiteral(t *testing.T) {
	cases := []string{
		`SELECT * FROM t AS 'a'`,
		`SELECT * FROM t 'a'`,
		`SELECT * FROM t AS "a" JOIN u AS 'b' ON t.id = b.id`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
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
