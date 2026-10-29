package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/mysql/ast"
)

// sourceClauseCases pairs statements holding TABLE and VALUES query primaries
// with where their ORDER BY, LIMIT, and INTO end up, one line per source,
// set operation, or parenthesized query, in walk order. A source line gives
// the text the source's Loc covers. MySQL 8.0 accepts every statement
// (TestSourceClausesOracle).
var sourceClauseCases = []struct {
	sql  string
	want []string
}{
	// A source standing alone holds its own clauses.
	{"TABLE t ORDER BY a LIMIT 1", []string{
		"TABLE t ORDER BY a LIMIT 1: order 1, limit",
	}},
	{"TABLE t LIMIT 1 INTO @x", []string{
		"TABLE t LIMIT 1 INTO @x: limit, into",
	}},
	{"VALUES ROW(1), ROW(2) ORDER BY 1 LIMIT 1", []string{
		"VALUES ROW(1), ROW(2) ORDER BY 1 LIMIT 1: order 1, limit",
	}},
	{"INSERT INTO t TABLE u ORDER BY a LIMIT 1", []string{
		"TABLE u ORDER BY a LIMIT 1: order 1, limit",
	}},
	{"EXPLAIN TABLE t ORDER BY a", []string{
		"TABLE t ORDER BY a: order 1",
	}},
	// Inside a query expression the SelectStmt over the source holds them.
	{"SELECT (TABLE t ORDER BY a LIMIT 1)", []string{
		"SelectStmt over TABLE t: order 1, limit",
	}},
	{"(TABLE t LIMIT 1 INTO @x)", []string{
		"parentheses: none",
		"SelectStmt over TABLE t: limit, into",
	}},
	{"SELECT * FROM (VALUES ROW(1), ROW(2) ORDER BY 1 LIMIT 1) x", []string{
		"SelectStmt over VALUES ROW(1), ROW(2): order 1, limit",
	}},
	{"WITH c AS (TABLE t ORDER BY a LIMIT 1) SELECT * FROM c", []string{
		"SelectStmt over TABLE t: order 1, limit",
	}},
	{"CREATE TABLE n AS TABLE t ORDER BY a LIMIT 1", []string{
		"SelectStmt over TABLE t: order 1, limit",
	}},
	// The ORDER BY and LIMIT after the last operand of a set operation apply
	// to the whole operation.
	{"CREATE VIEW v AS TABLE t UNION TABLE u ORDER BY a LIMIT 2", []string{
		"set operation: order 1, limit",
		"SelectStmt over TABLE t: none",
		"SelectStmt over TABLE u: none",
	}},
	{"SELECT 1 UNION VALUES ROW(2) ORDER BY 1 LIMIT 1", []string{
		"set operation: order 1, limit",
		"SelectStmt over VALUES ROW(2): none",
	}},
	{"SELECT 1 UNION (TABLE t ORDER BY a LIMIT 1) ORDER BY 1 LIMIT 2", []string{
		"set operation: order 1, limit",
		"parentheses: none",
		"SelectStmt over TABLE t: order 1, limit",
	}},
	{"(TABLE t ORDER BY a) ORDER BY b", []string{
		"parentheses: order 1",
		"SelectStmt over TABLE t: order 1",
	}},
	{"TABLE t ORDER BY (TABLE u ORDER BY (VALUES ROW(1) LIMIT 1))", []string{
		"TABLE t ORDER BY (TABLE u ORDER BY (VALUES ROW(1) LIMIT 1)): order 1",
		"SelectStmt over TABLE u: order 1",
		"SelectStmt over VALUES ROW(1): limit",
	}},
}

// TestSourceClausesHaveOneOwner: inside a query expression a TABLE or VALUES
// query primary is a SelectStmt over a TableStmt or ValuesStmt, and the
// SelectStmt alone holds the ORDER BY, LIMIT, and INTO, as MySQL's
// query_expression and select_stmt_with_into do. When the source held them
// too, every walk visited them twice, and a source nested in the ORDER BY of
// another doubled the walk at each level. A source standing alone (a
// statement, EXPLAIN, INSERT ... TABLE) holds its own.
func TestSourceClausesHaveOneOwner(t *testing.T) {
	for _, tc := range sourceClauseCases {
		list := ParseAndCheck(t, tc.sql)
		var got []string
		wrapped := make(map[ast.Node]bool)
		describe := func(owner string, orderBy []*ast.OrderByItem, limit *ast.Limit, into *ast.IntoClause) string {
			var parts []string
			if len(orderBy) > 0 {
				parts = append(parts, fmt.Sprintf("order %d", len(orderBy)))
			}
			if limit != nil {
				parts = append(parts, "limit")
			}
			if into != nil {
				parts = append(parts, "into")
			}
			if len(parts) == 0 {
				parts = append(parts, "none")
			}
			return owner + ": " + strings.Join(parts, ", ")
		}
		ownClauses := func(source ast.Node, orderBy []*ast.OrderByItem, limit *ast.Limit, into *ast.IntoClause) {
			if wrapped[source] && (len(orderBy) > 0 || limit != nil || into != nil) {
				t.Errorf("%s: the %T holds clauses its SelectStmt holds", tc.sql, source)
			}
		}
		ast.Inspect(list.Items[0], func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.SelectStmt:
				var source ast.Node
				var sourceLoc ast.Loc
				switch {
				case n.TableSource != nil:
					source, sourceLoc = n.TableSource, n.TableSource.Loc
				case n.ValuesSource != nil:
					source, sourceLoc = n.ValuesSource, n.ValuesSource.Loc
				case n.SetOp != ast.SetOpNone:
					got = append(got, describe("set operation", n.OrderBy, n.Limit, n.Into))
					return true
				case n.ParenSource != nil:
					got = append(got, describe("parentheses", n.OrderBy, n.Limit, n.Into))
					return true
				default:
					return true
				}
				wrapped[source] = true
				owner := "SelectStmt over " + tc.sql[sourceLoc.Start:sourceLoc.End]
				got = append(got, describe(owner, n.OrderBy, n.Limit, n.Into))
			case *ast.TableStmt:
				ownClauses(n, n.OrderBy, n.Limit, n.Into)
				if !wrapped[n] {
					got = append(got, describe(tc.sql[n.Loc.Start:n.Loc.End], n.OrderBy, n.Limit, n.Into))
				}
			case *ast.ValuesStmt:
				ownClauses(n, n.OrderBy, n.Limit, nil)
				if !wrapped[n] {
					got = append(got, describe(tc.sql[n.Loc.Start:n.Loc.End], n.OrderBy, n.Limit, nil))
				}
			}
			return true
		})
		if strings.Join(got, "\n") != strings.Join(tc.want, "\n") {
			t.Errorf("%s:\ngot:\n  %s\nwant:\n  %s", tc.sql, strings.Join(got, "\n  "), strings.Join(tc.want, "\n  "))
		}
	}
}

// TestSourceClausesOracle: MySQL 8.0 accepts every statement of
// sourceClauseCases, and the same TABLE and VALUES sources nested 30 deep in
// one another's ORDER BY (MySQL rejects 64 nested subqueries with 1473).
func TestSourceClausesOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, setup := range []string{"CREATE TABLE t (a INT, b INT)", "CREATE TABLE u (a INT, b INT)"} {
		if _, err := o.db.ExecContext(o.ctx, setup); err != nil {
			t.Fatalf("%s: %v", setup, err)
		}
	}
	sqls := []string{
		nest(30, "TABLE t ORDER BY (", "SELECT 1", ")"),
		nest(30, "VALUES ROW(1) ORDER BY (", "SELECT 1", ")"),
		"SELECT (" + nest(30, "TABLE t ORDER BY (", "SELECT 1", ")") + ")",
		"SELECT (" + nest(30, "VALUES ROW(1) ORDER BY (", "SELECT 1", ")") + ")",
	}
	for _, tc := range sourceClauseCases {
		sqls = append(sqls, tc.sql)
	}
	for _, sql := range sqls {
		if _, err := Parse(sql); err != nil {
			t.Errorf("omni rejects %q: %v", sql, err)
		}
		if !o.checkSyntax(sql) {
			t.Errorf("MySQL 8.0 rejects %q", sql)
		}
	}
}
