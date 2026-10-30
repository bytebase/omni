package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// sourceClauseCases pairs statements holding table value constructors with
// where their ORDER BY, LIMIT, and INTO end up, one line per constructor, set
// operation, or parenthesized query, in walk order. A constructor line gives
// the text the constructor's Loc covers. MariaDB 11.8 accepts every statement
// (TestTableValueConstructorOracle runs them).
var sourceClauseCases = []struct {
	sql  string
	want []string
}{
	// A statement, an EXPLAIN target, or an INSERT source is a SelectStmt
	// over the constructor, which holds the clauses.
	{"VALUES (1), (2) ORDER BY 1 LIMIT 1", []string{
		"SelectStmt over VALUES (1), (2): order 1, limit",
	}},
	{"VALUES (1) LIMIT 1 INTO @x", []string{
		"SelectStmt over VALUES (1): limit, into",
	}},
	{"INSERT INTO t VALUES (1, 2) ORDER BY 1 LIMIT 1", []string{
		"SelectStmt over VALUES (1, 2): order 1, limit",
	}},
	{"EXPLAIN VALUES (1) ORDER BY 1", []string{
		"SelectStmt over VALUES (1): order 1",
	}},
	// So is one inside a query expression.
	{"SELECT (VALUES (1) ORDER BY 1 LIMIT 1)", []string{
		"SelectStmt over VALUES (1): order 1, limit",
	}},
	{"(VALUES (1) LIMIT 1) INTO @x", []string{
		"parentheses: into",
		"SelectStmt over VALUES (1): limit",
	}},
	{"SELECT * FROM (VALUES (1), (2) ORDER BY 1 LIMIT 1) x", []string{
		"SelectStmt over VALUES (1), (2): order 1, limit",
	}},
	{"WITH c AS (VALUES (1) ORDER BY 1 LIMIT 1) SELECT * FROM c", []string{
		"SelectStmt over VALUES (1): order 1, limit",
	}},
	{"CREATE TABLE n AS VALUES (1) ORDER BY 1 LIMIT 1", []string{
		"SelectStmt over VALUES (1): order 1, limit",
	}},
	// The ORDER BY and LIMIT after the last operand of a set operation apply
	// to the whole operation, and so does the INTO after it.
	{"CREATE VIEW v AS VALUES (1) UNION VALUES (2) ORDER BY 1 LIMIT 2", []string{
		"set operation: order 1, limit",
		"SelectStmt over VALUES (1): none",
		"SelectStmt over VALUES (2): none",
	}},
	{"SELECT 1 UNION VALUES (2) ORDER BY 1 LIMIT 1", []string{
		"set operation: order 1, limit",
		"SelectStmt over VALUES (2): none",
	}},
	{"VALUES (1) UNION VALUES (2) LIMIT 1 INTO @x", []string{
		"set operation: limit, into",
		"SelectStmt over VALUES (1): none",
		"SelectStmt over VALUES (2): none",
	}},
	{"SELECT 1 UNION (VALUES (1) ORDER BY 1 LIMIT 1) ORDER BY 1 LIMIT 2", []string{
		"set operation: order 1, limit",
		"parentheses: none",
		"SelectStmt over VALUES (1): order 1, limit",
	}},
	{"(VALUES (1) ORDER BY 1) ORDER BY 1", []string{
		"parentheses: order 1",
		"SelectStmt over VALUES (1): order 1",
	}},
	{"VALUES (1) ORDER BY (VALUES (2) ORDER BY (VALUES (3) LIMIT 1))", []string{
		"SelectStmt over VALUES (1): order 1",
		"SelectStmt over VALUES (2): order 1",
		"SelectStmt over VALUES (3): limit",
	}},
}

// TestSourceClausesHaveOneOwner: a table value constructor is always a
// SelectStmt over a ValuesStmt, and the SelectStmt alone holds the ORDER BY,
// LIMIT, and INTO, as MariaDB's query_expression_body_ext and select rules
// do. When the constructor held them too, every walk visited them twice, and
// a constructor nested in the ORDER BY of another doubled the walk at each
// level. MariaDB has no TABLE statement, so no TableStmt appears.
func TestSourceClausesHaveOneOwner(t *testing.T) {
	for _, tc := range sourceClauseCases {
		list := ParseAndCheck(t, tc.sql)
		if len(list.Items) != 1 {
			t.Errorf("%s: %d statements, want 1", tc.sql, len(list.Items))
			continue
		}
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
		ast.Inspect(list.Items[0], func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.SelectStmt:
				switch {
				case n.ValuesSource != nil:
					wrapped[n.ValuesSource] = true
					owner := "SelectStmt over " + tc.sql[n.ValuesSource.Loc.Start:n.ValuesSource.Loc.End]
					got = append(got, describe(owner, n.OrderBy, n.Limit, n.Into))
				case n.SetOp != ast.SetOpNone:
					got = append(got, describe("set operation", n.OrderBy, n.Limit, n.Into))
				case n.ParenSource != nil:
					got = append(got, describe("parentheses", n.OrderBy, n.Limit, n.Into))
				}
			case *ast.ValuesStmt:
				if !wrapped[n] {
					t.Errorf("%s: a ValuesStmt outside a SelectStmt", tc.sql)
				}
				if len(n.OrderBy) > 0 || n.Limit != nil {
					t.Errorf("%s: the ValuesStmt holds clauses its SelectStmt holds", tc.sql)
				}
			case *ast.TableStmt:
				t.Errorf("%s: a TableStmt", tc.sql)
			}
			return true
		})
		if strings.Join(got, "\n") != strings.Join(tc.want, "\n") {
			t.Errorf("%s:\ngot:\n  %s\nwant:\n  %s", tc.sql, strings.Join(got, "\n  "), strings.Join(tc.want, "\n  "))
		}
	}
}
