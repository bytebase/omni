package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// A PIVOT ends the FROM clause and holds its aggregates, FOR columns and
// IN values, each with its alias and position.
func TestPivotClause(t *testing.T) {
	sql := "SELECT * FROM t, u PIVOT (SUM(a) AS s, COUNT(*) FOR (b, c) IN ((1, 'x') AS one, (2, 'y'))) WHERE d = 1"
	file, err := Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	stmt := file.Stmts[0].(*ast.SelectStmt)
	if len(stmt.From) != 2 || stmt.Where == nil {
		t.Fatalf("From %d items, Where %v", len(stmt.From), stmt.Where)
	}
	pv := stmt.Pivot
	if pv == nil {
		t.Fatal("no PIVOT")
	}
	if got := sql[pv.Loc.Start:pv.Loc.End]; got != "PIVOT (SUM(a) AS s, COUNT(*) FOR (b, c) IN ((1, 'x') AS one, (2, 'y')))" {
		t.Errorf("PIVOT text %q", got)
	}
	if len(pv.Aggregations) != 2 || pv.Aggregations[0].Alias != "s" || pv.Aggregations[1].Alias != "" {
		t.Errorf("aggregations %+v", pv.Aggregations)
	}
	if got := sql[pv.Aggregations[0].Loc.Start:pv.Aggregations[0].Loc.End]; got != "SUM(a) AS s" {
		t.Errorf("aggregation text %q", got)
	}
	if len(pv.Columns) != 2 || pv.Columns[0] != "b" || pv.Columns[1] != "c" {
		t.Errorf("columns %v", pv.Columns)
	}
	if len(pv.Values) != 2 || pv.Values[0].Alias != "one" || len(pv.Values[1].Exprs) != 2 {
		t.Errorf("values %+v", pv.Values)
	}
	if got := sql[pv.Values[0].Loc.Start:pv.Values[0].Loc.End]; got != "(1, 'x') AS one" {
		t.Errorf("value text %q", got)
	}
}

// PIVOT is nonReserved: before anything but a pivot body it is an alias.
func TestPivotAlias(t *testing.T) {
	for sql, alias := range map[string]string{
		"SELECT * FROM t pivot":                                            "pivot",
		"SELECT * FROM (SELECT 1 a) pivot (a)":                             "pivot",
		"SELECT * FROM t pivot PIVOT (SUM(a) FOR b IN (1))":                "pivot",
		"SELECT * FROM (SELECT 1 a) pivot (a) PIVOT (MAX(a) FOR a IN (1))": "pivot",
	} {
		file, err := Parse(sql)
		if err != nil {
			t.Errorf("%s: %v", sql, err)
			continue
		}
		ref, ok := file.Stmts[0].(*ast.SelectStmt).From[0].(*ast.TableRef)
		if !ok || ref.Alias != alias {
			t.Errorf("%s: From[0] %+v, want alias %q", sql, file.Stmts[0].(*ast.SelectStmt).From[0], alias)
		}
	}
}
