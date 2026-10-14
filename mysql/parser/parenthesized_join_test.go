package parser

import (
	"testing"

	ast "github.com/bytebase/omni/mysql/ast"
)

// TestParenthesizedJoinGroup pins the FROM-clause table_reference parser on
// parenthesized join groups — `( table_reference )` where the inner is itself a
// (possibly joined) table reference. MySQL stores a multi-table-join view's FROM
// clause in this canonical left-deep, fully-parenthesized form (the shape emitted
// by SHOW CREATE VIEW / mysqldump), e.g. sakila's film_list:
//
//	from ((((`category` left join `film_category` on(...)) left join `film` on(...))
//	       join `film_actor` on(...)) join `actor` on(...))
//
// A leading run of '(' here must NOT be mistaken for a parenthesized sub-SELECT;
// the deciding token is the first one past the parens (a query primary means a
// derived table, a table reference means a join group).
func TestParenthesizedJoinGroup(t *testing.T) {
	for _, sql := range []string{
		// single parenthesized table
		"SELECT * FROM (t)",
		"SELECT * FROM ((t))",
		"SELECT * FROM (((t)))",
		// 2-table parenthesized join group
		"SELECT * FROM (a JOIN b ON a.id = b.id)",
		"SELECT * FROM (`a` JOIN `b` ON(`a`.`id` = `b`.`id`))",
		// join group whose factors carry aliases (the SHOW CREATE VIEW form)
		"SELECT * FROM (`category` `a` JOIN `film_category` `b` ON(`a`.`category_id` = `b`.`category_id`))",
		// nested left-deep groups (3, 4, 5 tables)
		"SELECT * FROM ((a JOIN b ON a.id = b.id) JOIN c ON b.id = c.id)",
		"SELECT * FROM (((a JOIN b ON a.id=b.id) JOIN c ON b.id=c.id) JOIN d ON c.id=d.id)",
		"SELECT * FROM ((((`category` left join `film_category` on((`category`.`category_id` = `film_category`.`category_id`))) left join `film` on((`film_category`.`film_id` = `film`.`film_id`))) join `film_actor` on((`film`.`film_id` = `film_actor`.`film_id`))) join `actor` on((`film_actor`.`actor_id` = `actor`.`actor_id`)))",
		// right-deep parenthesized group
		"SELECT * FROM (a JOIN (b JOIN c ON b.id = c.id) ON a.id = b.id)",
		// LEFT / INNER / CROSS / USING inside groups
		"SELECT * FROM (a LEFT JOIN b ON a.id = b.id)",
		"SELECT * FROM (a INNER JOIN b USING (id))",
		"SELECT * FROM (a CROSS JOIN b)",
		"SELECT * FROM ((a LEFT JOIN b ON a.id=b.id) JOIN c ON b.id=c.id)",
		// full sakila film_list CREATE VIEW (the real-world repro)
		"CREATE VIEW `film_list` AS select `film`.`film_id` AS `FID` from ((((`category` left join `film_category` on((`category`.`category_id` = `film_category`.`category_id`))) left join `film` on((`film_category`.`film_id` = `film`.`film_id`))) join `film_actor` on((`film`.`film_id` = `film_actor`.`film_id`))) join `actor` on((`film_actor`.`actor_id` = `actor`.`actor_id`))) group by `film`.`film_id`",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParenthesizedJoinCompoundOn pins the residual case #356 did not cover: an
// aliased parenthesized join group whose ON condition is COMPOUND (AND/OR), used
// as the left operand of a further top-level join. MySQL's SHOW CREATE VIEW
// double-wraps a compound ON — `on(((a = b) and (c = d)))` — and a redundant
// single-condition wrap — `on(((a = b)))`. The leading `((` inside such an ON
// must parse as a parenthesized scalar expression, NOT a subquery. This is the
// exact shape of the employees `current_dept_emp` view.
func TestParenthesizedJoinCompoundOn(t *testing.T) {
	for _, sql := range []string{
		// compound (AND) ON inside a parenthesized aliased join group
		"SELECT * FROM (`a` `d` JOIN `b` `l` ON((`d`.`x` = `l`.`x`) AND (`d`.`y` = `l`.`y`)))",
		// double-wrapped compound ON (the exact SHOW CREATE VIEW rendering)
		"SELECT * FROM (`a` `d` JOIN `b` `l` ON(((`d`.`x` = `l`.`x`) AND (`d`.`y` = `l`.`y`))))",
		// redundant-wrapped single-condition ON
		"SELECT * FROM (`a` `d` JOIN `b` `l` ON(((`d`.`x` = `l`.`x`))))",
		// aliased compound-ON paren group as LEFT operand of a further join (employees shape)
		"SELECT * FROM (`a` `d` JOIN `b` `l` ON(((`d`.`x` = `l`.`x`) AND (`d`.`y` = `l`.`y`)))) JOIN `c` `dp` ON((`d`.`z` = `dp`.`z`))",
		// the exact employees current_dept_emp FROM clause
		"SELECT 1 FROM ((`dept_emp` `d` JOIN `dept_emp_latest_date` `l` ON(((`d`.`emp_no` = `l`.`emp_no`) AND (`d`.`from_date` = `l`.`from_date`)))) JOIN `departments` `dp` ON((`d`.`dept_no` = `dp`.`dept_no`)))",
		// the full employees current_dept_emp CREATE VIEW (engine-emitted canonical form)
		"CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `current_dept_emp` AS select `l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`l`.`from_date` AS `from_date`,`l`.`to_date` AS `to_date` from ((`dept_emp` `d` join `dept_emp_latest_date` `l` on(((`d`.`emp_no` = `l`.`emp_no`) and (`d`.`from_date` = `l`.`from_date`)))) join `departments` `dp` on((`d`.`dept_no` = `dp`.`dept_no`)))",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParenthesizedScalarExprNotSubquery guards the disambiguation on the
// expression side: a doubly (or more) parenthesized SCALAR expression `((expr))`
// — anywhere an expression is parsed — must NOT be misread as a subquery. These
// are the redundant-paren forms MySQL emits and that a fixed one-token '(('
// lookahead wrongly classified as a query expression.
func TestParenthesizedScalarExprNotSubquery(t *testing.T) {
	for _, sql := range []string{
		"SELECT ((1))",
		"SELECT (((1 + 2)))",
		"SELECT ((a = b)) FROM t",
		"SELECT * FROM t WHERE ((a = b) AND (c = d))",
		"SELECT * FROM t WHERE (((a = b)))",
		"SELECT * FROM t WHERE ((a > 1))",
		"SELECT ((a)) + ((b)) FROM t",
		"SELECT * FROM t WHERE a IN ((1), (2))",
		"SELECT (SELECT ((1)))",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParenthesizedSubqueryStillSubquery guards the OTHER side of the same
// disambiguation: a '(' run that ultimately opens a query primary must still be
// read as a subquery in expression contexts (scalar subquery, IN, EXISTS,
// quantified), including the doubly-parenthesized set-op forms.
func TestParenthesizedSubqueryStillSubquery(t *testing.T) {
	for _, sql := range []string{
		"SELECT (SELECT 1)",
		"SELECT ((SELECT 1))",
		"SELECT (((SELECT 1)))",
		"SELECT * FROM t WHERE a IN (SELECT x FROM u)",
		"SELECT * FROM t WHERE a IN ((SELECT x FROM u))",
		"SELECT * FROM t WHERE EXISTS (SELECT 1)",
		"SELECT * FROM t WHERE EXISTS ((SELECT 1))",
		"SELECT * FROM t WHERE a = (SELECT MAX(x) FROM u)",
		"SELECT * FROM t WHERE a > ALL (SELECT x FROM u)",
		"SELECT * FROM t WHERE a > ALL ((SELECT x FROM u))",
		"SELECT * FROM t WHERE a = ((SELECT 1) UNION (SELECT 2))",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParenthesizedJoinGroupRegression guards the derived-table side of the
// disambiguation: a '(' run that ultimately opens a query primary must still be
// parsed as a derived-table subquery (with its alias / column list), and plain /
// comma / explicit joins must be unaffected.
func TestParenthesizedJoinGroupRegression(t *testing.T) {
	for _, sql := range []string{
		"SELECT * FROM t",
		"SELECT * FROM a, b",
		"SELECT * FROM a JOIN b ON a.id = b.id",
		"SELECT * FROM a LEFT JOIN b ON a.id = b.id LEFT JOIN c ON b.id = c.id",
		// derived tables
		"SELECT * FROM (SELECT 1) x",
		"SELECT * FROM ((SELECT 1)) x",
		"SELECT * FROM (((SELECT 1))) x",
		"SELECT * FROM (SELECT 1 AS c) AS t(x)",
		"SELECT * FROM ((SELECT 1) UNION (SELECT 2)) u",
		"SELECT * FROM (SELECT * FROM (SELECT 1) y) x",
		// derived table mixed with a join group
		"SELECT * FROM (a JOIN (SELECT 1) d ON a.id = d.x)",
		"SELECT * FROM t JOIN (SELECT x FROM u) d ON t.c = d.x",
		// LATERAL derived table (8.0)
		"SELECT * FROM t, LATERAL (SELECT 1) d",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParenthesizedJoinGroupStructure asserts the parenthesized join group folds
// into the correct left-deep JoinClause tree (right join types, ON conditions
// present, non-degenerate Loc spans) — not merely that it parses.
func TestParenthesizedJoinGroupStructure(t *testing.T) {
	sql := "SELECT * FROM ((((`category` left join `film_category` on(1)) left join `film` on(2)) join `film_actor` on(3)) join `actor` on(4))"
	list := ParseAndCheck(t, sql)
	sel, ok := list.Items[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected *SelectStmt, got %T", list.Items[0])
	}
	if len(sel.From) != 1 {
		t.Fatalf("expected 1 folded FROM element, got %d", len(sel.From))
	}

	// Outer-to-inner spine of the left-deep tree.
	wantSpine := []struct {
		jt    ast.JoinType
		right string
	}{
		{ast.JoinInner, "actor"},
		{ast.JoinInner, "film_actor"},
		{ast.JoinLeft, "film"},
		{ast.JoinLeft, "film_category"},
	}
	cur := sel.From[0]
	for i, w := range wantSpine {
		jc, ok := cur.(*ast.JoinClause)
		if !ok {
			t.Fatalf("spine[%d]: expected *JoinClause, got %T", i, cur)
		}
		if jc.Type != w.jt {
			t.Errorf("spine[%d]: join type = %v, want %v", i, jc.Type, w.jt)
		}
		if jc.Condition == nil {
			t.Errorf("spine[%d]: missing ON/USING condition", i)
		}
		if jc.Loc.End <= jc.Loc.Start {
			t.Errorf("spine[%d]: degenerate Loc span %d..%d", i, jc.Loc.Start, jc.Loc.End)
		}
		rt, ok := jc.Right.(*ast.TableRef)
		if !ok {
			t.Fatalf("spine[%d]: right expected *TableRef, got %T", i, jc.Right)
		}
		if rt.Name != w.right {
			t.Errorf("spine[%d]: right table = %q, want %q", i, rt.Name, w.right)
		}
		cur = jc.Left
	}
	if leaf, ok := cur.(*ast.TableRef); !ok {
		t.Fatalf("innermost left expected *TableRef, got %T", cur)
	} else if leaf.Name != "category" {
		t.Errorf("innermost left = %q, want %q", leaf.Name, "category")
	}
}

// TestParenthesizedDerivedTableStillSubquery verifies derived tables are not
// misrouted to the join-group path by the disambiguation.
func TestParenthesizedDerivedTableStillSubquery(t *testing.T) {
	for _, sql := range []string{
		"SELECT * FROM (SELECT 1) x",
		"SELECT * FROM ((SELECT 1)) y",
		"SELECT * FROM (((SELECT 1 UNION SELECT 2))) z",
	} {
		t.Run(sql, func(t *testing.T) {
			list := ParseAndCheck(t, sql)
			sel := list.Items[0].(*ast.SelectStmt)
			if len(sel.From) != 1 {
				t.Fatalf("expected 1 FROM element, got %d", len(sel.From))
			}
			if _, ok := sel.From[0].(*ast.SubqueryExpr); !ok {
				t.Errorf("expected *SubqueryExpr, got %T", sel.From[0])
			}
		})
	}
}
