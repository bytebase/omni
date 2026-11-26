package parser

import (
	"testing"

	"github.com/bytebase/omni/bigquery/ast"
)

// TestEverySiteForms holds the forms BigQuery dry-run accepted at sites the
// second every-alternative pass reached (testdata/compat.tsv), and the forms
// next to them it refused.
func TestEverySiteForms(t *testing.T) {
	accept := []string{
		// string_literal is one or more adjacent components at every site
		// (googlesql.tm:11369), bytes_literal likewise (11431).
		"CREATE FUNCTION a() AS 'a' 'b'",
		"CREATE FUNCTION a() OPTIONS() AS 'a' 'b'",
		"CREATE PROCEDURE a() LANGUAGE a AS 'a' 'b'",
		"SELECT 1 |> UNPIVOT(a FOR a IN(a 'a' 'b'))",
		"DECLARE a STRING COLLATE 'und' ':ci'",
		"DECLARE a a('a' 'b')",
		"DECLARE a a(1, b'a' b'b')",
		"SELECT DATE '2020' '-01-01'",
		"SELECT NUMERIC '1' '2', BIGNUMERIC '1' '2', JSON '{}' ''",
		"SELECT RANGE<DATE> '[2020-01-01,' ' 2021-01-01)'",
		// templated_parameter_kind (googlesql.tm:10400).
		"CREATE FUNCTION a(x ANY PROTO, y ANY ENUM, z ANY STRUCT, w ANY ARRAY) AS (1)",
		"CREATE FUNCTION a(x ANY type, y ANY `TABLE`) AS (1)",
		"CREATE PROCEDURE a(a ANY PROTO) BEGIN END",
		// limit_all (googlesql.tm:8504), in a query and in a call.
		"SELECT 1 LIMIT ALL",
		"SELECT 1 LIMIT ALL OFFSET 1",
		"SELECT ARRAY_AGG(1 LIMIT ALL OFFSET 1)",
		// named_argument's input_table_argument (googlesql.tm:10907).
		"SELECT * FROM UNNEST([1], a => INPUT TABLE)",
		"SELECT f(1, a => INPUT TABLE)",
	}
	for _, sql := range accept {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q) = %v, want accept", sql, err)
		}
	}
	reject := []string{
		// A string and bytes do not concatenate.
		"SELECT DATE '2020' b'a'",
		"CREATE FUNCTION a() AS 'a' b'a'",
		// BigQuery: "Syntax error: unexpected ANY template type".
		"CREATE FUNCTION a(x ANY FOO) AS (1)",
		"CREATE FUNCTION a(x ANY `foo`) AS (1)",
		// Only a named argument takes INPUT TABLE.
		"SELECT f(INPUT TABLE)",
		"SELECT 1 LIMIT ALL zzz",
	}
	for _, sql := range reject {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted, want reject", sql)
		}
	}
}

func TestEverySiteFormsTree(t *testing.T) {
	stmts, err := Parse("SELECT DATE '2020' '-01-01'")
	if err != nil {
		t.Fatal(err)
	}
	var lit *ast.TypedLiteral
	ast.Inspect(stmts.Stmts[0], func(n ast.Node) bool {
		if l, ok := n.(*ast.TypedLiteral); ok {
			lit = l
		}
		return true
	})
	if lit == nil || lit.Value != "2020-01-01" {
		t.Fatalf("typed literal = %+v, want value 2020-01-01", lit)
	}

	stmts, err = Parse("SELECT 1 LIMIT ALL OFFSET 2")
	if err != nil {
		t.Fatal(err)
	}
	q, ok := stmts.Stmts[0].(*ast.QueryStmt)
	if !ok || !q.LimitAll || q.Limit != nil || q.Offset == nil {
		t.Fatalf("query = %+v, want LIMIT ALL with an OFFSET", stmts.Stmts[0])
	}
}
