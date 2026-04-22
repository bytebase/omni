package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestParseParenMultiJoin covers the parenthesized joined_table shape
// that pg_get_viewdef() emits, e.g.:
//
//	SELECT * FROM ((a JOIN b ON TRUE) JOIN c ON TRUE)
//
// A single-token peek after '(' cannot distinguish
//
//	'(' '(' SELECT ... ')' ')'         — nested subquery
//	'(' '(' a JOIN b ')' JOIN c ')'    — parenthesized joined_table
//
// and the previous implementation always routed '((' to the subquery
// branch, so every pg_get_viewdef() output with 2+ parenthesized joins
// failed to parse — which cascaded through catalog.Exec and killed
// SQL review for every statement against the instance (BYT-9308).
func TestParseParenMultiJoin(t *testing.T) {
	cases := []string{
		// Plain single-level parenthesized joined_table (no alias).
		`SELECT * FROM (a JOIN b ON TRUE)`,
		// Plain parenthesized joined_table with alias clause.
		`SELECT * FROM (a JOIN b ON TRUE) AS jt`,
		// Plain parenthesized joined_table with column-list alias.
		`SELECT * FROM (a JOIN b ON TRUE) jt(x, y)`,
		// Exact shape from BYT-9315 description.
		`SELECT * FROM ((a JOIN b ON TRUE) JOIN c ON TRUE)`,
		// Deeper nesting.
		`SELECT * FROM (((a JOIN b ON TRUE) JOIN c ON TRUE) JOIN d ON TRUE)`,
		// Trivially wrapped joined_table.
		`SELECT * FROM ((a JOIN b ON TRUE))`,
		// Subquery as a join operand inside the outer parens — '((' but
		// outer is joined_table, not subquery.
		`SELECT * FROM ((SELECT 1) x JOIN (SELECT 2) y ON TRUE)`,
		`SELECT * FROM ((SELECT 1) x CROSS JOIN b)`,
		// Alias on the outer parenthesized joined_table.
		`SELECT * FROM ((a JOIN b ON TRUE) JOIN c ON TRUE) AS jt`,
		// Column-list alias on the outer parenthesized joined_table.
		`SELECT * FROM ((a JOIN b ON TRUE) JOIN c ON TRUE) jt(x, y, z)`,
		// Real pg_get_viewdef() output shape (double parens on join qual).
		`CREATE VIEW v AS SELECT a.*, b.*, c.* FROM ((a JOIN b ON ((a.id = b.aid))) JOIN c ON ((b.id = c.bid)))`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("parse failed: %v", err)
			}
		})
	}
}

// TestParseParenMultiJoinAST verifies the JoinExpr nesting is built
// left-associatively for the pg_get_viewdef() shape, matching
// PostgreSQL's own parse tree.
func TestParseParenMultiJoinAST(t *testing.T) {
	stmts, err := Parse(`SELECT * FROM ((a JOIN b ON TRUE) JOIN c ON TRUE)`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	if sel.FromClause == nil || len(sel.FromClause.Items) != 1 {
		t.Fatalf("expected 1 FROM item, got %v", sel.FromClause)
	}
	outer, ok := sel.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected outer JoinExpr, got %T", sel.FromClause.Items[0])
	}
	inner, ok := outer.Larg.(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected inner JoinExpr as outer.Larg, got %T", outer.Larg)
	}
	a, ok := inner.Larg.(*nodes.RangeVar)
	if !ok || a.Relname != "a" {
		t.Fatalf("expected RangeVar a as inner.Larg, got %T %+v", inner.Larg, inner.Larg)
	}
	b, ok := inner.Rarg.(*nodes.RangeVar)
	if !ok || b.Relname != "b" {
		t.Fatalf("expected RangeVar b as inner.Rarg, got %T %+v", inner.Rarg, inner.Rarg)
	}
	c, ok := outer.Rarg.(*nodes.RangeVar)
	if !ok || c.Relname != "c" {
		t.Fatalf("expected RangeVar c as outer.Rarg, got %T %+v", outer.Rarg, outer.Rarg)
	}
}

// TestParseParenSetOpInSubquery covers the '(' select_no_parens ')' where
// the select_no_parens is a set-op with parenthesized operands:
//
//	SELECT * FROM ((SELECT 1) UNION (SELECT 2))
//
// The prior parseSelectWithParens short-circuited on cur == '(' and
// recursed as a nested select_with_parens, consuming the inner
// '(SELECT 1)' as a self-contained subquery and then failing to find
// ')' where UNION appeared. Left-factoring to go through
// parseSelectNoParens lets parseSelectClause's set-op precedence
// climbing handle the inner '(' as the first operand of a UNION.
func TestParseParenSetOpInSubquery(t *testing.T) {
	cases := []string{
		`SELECT * FROM ((SELECT 1) UNION (SELECT 2))`,
		`SELECT * FROM ((SELECT 1) UNION ALL (SELECT 2))`,
		`SELECT * FROM ((SELECT 1) INTERSECT (SELECT 2))`,
		`SELECT * FROM ((SELECT 1) EXCEPT (SELECT 2))`,
		// Nested set-ops
		`SELECT * FROM (((SELECT 1) UNION (SELECT 2)) UNION (SELECT 3))`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("parse failed: %v", err)
			}
		})
	}
}

// TestParseParenTableRefRejectsNonJoined verifies that '(' X ')' in a
// FROM clause rejects X that is not a joined_table (i.e., has no JOIN).
// PG's grammar only allows '(' joined_table ')' alias_clause as the
// parenthesized-joined-table production of table_ref; a single-relation
// or aliased-subquery wrapped in extra parens is not a valid table_ref.
func TestParseParenTableRefRejectsNonJoined(t *testing.T) {
	cases := []string{
		// Single relation wrapped in parens — not a valid table_ref.
		`SELECT * FROM (a)`,
		`SELECT * FROM ((a))`,
		// Subquery with alias wrapped in extra parens — PG rejects
		// because the outer '(' ... ')' must contain a joined_table.
		`SELECT * FROM ((SELECT 1) x)`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("expected syntax error, got nil")
			}
		})
	}
}

// TestParseParenSubqueryStillWorks guards the non-regression direction:
// the disambiguation must not break genuine select_with_parens subqueries.
func TestParseParenSubqueryStillWorks(t *testing.T) {
	cases := []string{
		`SELECT * FROM (SELECT 1)`,
		`SELECT * FROM ((SELECT 1))`,
		`SELECT * FROM (((SELECT 1)))`,
		`SELECT * FROM (VALUES (1), (2)) AS v`,
		`SELECT * FROM (WITH cte AS (SELECT 1) SELECT * FROM cte)`,
		`SELECT * FROM (TABLE foo)`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			stmts, err := Parse(sql)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
			if _, ok := sel.FromClause.Items[0].(*nodes.RangeSubselect); !ok {
				t.Fatalf("expected RangeSubselect, got %T", sel.FromClause.Items[0])
			}
		})
	}
}
