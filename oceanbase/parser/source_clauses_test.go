package parser

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
)

// sourceClauseCases pairs statements holding VALUES query primaries with
// where their ORDER BY and LIMIT end up, one line per source, set operation,
// or parenthesized query, in walk order. A source line gives the text the
// source's Loc covers. MySQL 8.0 (TestSourceClausesOracle) and OceanBase
// (TestSourceClausesOceanBase) accept every statement.
//
// OceanBase has no TABLE query, so where MySQL 8.0 takes a TABLE source in
// each of these places, OceanBase answers 1064 (query_forms_test.go), and
// every context is covered with a VALUES source.
var sourceClauseCases = []struct {
	sql  string
	want []string
}{
	// A source standing alone holds its own clauses.
	{"VALUES ROW(1), ROW(2) ORDER BY 1 LIMIT 1", []string{
		"VALUES ROW(1), ROW(2) ORDER BY 1 LIMIT 1: order 1, limit",
	}},
	{"EXPLAIN VALUES ROW(1), ROW(2) ORDER BY 1", []string{
		"VALUES ROW(1), ROW(2) ORDER BY 1: order 1",
	}},
	// Inside a query expression the SelectStmt over the source holds them.
	{"SELECT (VALUES ROW(1), ROW(2) ORDER BY 1 LIMIT 1)", []string{
		"SelectStmt over VALUES ROW(1), ROW(2): order 1, limit",
	}},
	{"(VALUES ROW(1), ROW(2) LIMIT 1)", []string{
		"parentheses: none",
		"SelectStmt over VALUES ROW(1), ROW(2): limit",
	}},
	{"SELECT * FROM (VALUES ROW(1), ROW(2) ORDER BY 1 LIMIT 1) x", []string{
		"SelectStmt over VALUES ROW(1), ROW(2): order 1, limit",
	}},
	{"WITH c AS (VALUES ROW(1), ROW(2) ORDER BY 1 LIMIT 1) SELECT * FROM c", []string{
		"SelectStmt over VALUES ROW(1), ROW(2): order 1, limit",
	}},
	{"CREATE TABLE n AS VALUES ROW(1), ROW(2) ORDER BY 1 LIMIT 1", []string{
		"SelectStmt over VALUES ROW(1), ROW(2): order 1, limit",
	}},
	// The ORDER BY and LIMIT after the last operand of a set operation apply
	// to the whole operation.
	{"CREATE VIEW v AS VALUES ROW(1) UNION VALUES ROW(2) ORDER BY 1 LIMIT 2", []string{
		"set operation: order 1, limit",
		"SelectStmt over VALUES ROW(1): none",
		"SelectStmt over VALUES ROW(2): none",
	}},
	{"SELECT 1 UNION VALUES ROW(2) ORDER BY 1 LIMIT 1", []string{
		"set operation: order 1, limit",
		"SelectStmt over VALUES ROW(2): none",
	}},
	{"SELECT 1 UNION (VALUES ROW(1), ROW(2) ORDER BY 1 LIMIT 1) ORDER BY 1 LIMIT 2", []string{
		"set operation: order 1, limit",
		"parentheses: none",
		"SelectStmt over VALUES ROW(1), ROW(2): order 1, limit",
	}},
	{"(VALUES ROW(1), ROW(2) ORDER BY 1) ORDER BY 1", []string{
		"parentheses: order 1",
		"SelectStmt over VALUES ROW(1), ROW(2): order 1",
	}},
	{"VALUES ROW(1) ORDER BY (VALUES ROW(2) ORDER BY (VALUES ROW(3) LIMIT 1))", []string{
		"VALUES ROW(1) ORDER BY (VALUES ROW(2) ORDER BY (VALUES ROW(3) LIMIT 1)): order 1",
		"SelectStmt over VALUES ROW(2): order 1",
		"SelectStmt over VALUES ROW(3): limit",
	}},
}

// TestSourceClausesHaveOneOwner: inside a query expression a VALUES query
// primary is a SelectStmt over a ValuesStmt, and the SelectStmt alone holds the
// ORDER BY and LIMIT, as OceanBase's table_values_clause_with_order_by_and_limit
// does. When the source held them too, every walk visited them twice, and a
// source nested in the ORDER BY of another doubled the walk at each level. A
// source standing alone (a statement, EXPLAIN) holds its own.
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

// sourceClauseStatements are the statements of sourceClauseCases and the
// same VALUES sources nested 30 deep in one another's ORDER BY (MySQL rejects
// 64 nested subqueries with 1473).
func sourceClauseStatements() []string {
	sqls := []string{
		nest(30, "VALUES ROW(1) ORDER BY (", "SELECT 1", ")"),
		"SELECT (" + nest(30, "VALUES ROW(1) ORDER BY (", "SELECT 1", ")") + ")",
	}
	for _, tc := range sourceClauseCases {
		sqls = append(sqls, tc.sql)
	}
	return sqls
}

// TestSourceClausesOracle: MySQL 8.0 accepts every statement of
// sourceClauseStatements.
func TestSourceClausesOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, sql := range sourceClauseStatements() {
		if _, err := Parse(sql); err != nil {
			t.Errorf("omni rejects %q: %v", sql, err)
		}
		if !o.checkSyntax(t, sql) {
			t.Errorf("MySQL 8.0 rejects %q", sql)
		}
	}
}

// TestSourceClausesOceanBase: OceanBase accepts every statement of
// sourceClauseStatements, read through PREPARE (DSN-gated, see
// oceanBaseOracle):
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/parser/ -run OceanBase -count=1
func TestSourceClausesOceanBase(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, sql := range sourceClauseStatements() {
		if !prepareAcceptsOceanBase(t, ctx, db, sql) {
			t.Errorf("OceanBase rejects %q", sql)
		}
	}
}
