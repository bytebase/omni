package parser

import "testing"

// TestParenthesizedQueryExpression pins TiDB v8.5.0 parity: a statement may be a
// parenthesized query expression — ( query ), parenthesized set-operation
// operands, nested parens, per-operand or trailing ORDER BY / LIMIT, and
// INTERSECT / EXCEPT. All container-verified accepted on pingcap/tidb:v8.5.0
// (BYT-9650). Before this change omni rejected every leading-'(' statement
// because parseStmt routed it to parseSelectStmt, which consumed the '(' as if
// it were the SELECT keyword.
func TestParenthesizedQueryExpression(t *testing.T) {
	accepted := []string{
		"(SELECT 1)",
		"(SELECT a FROM t)",
		"((SELECT 1))",
		"((SELECT a FROM t))",
		"(SELECT 1) UNION (SELECT 2)",
		"(SELECT 1) UNION SELECT 2",
		"SELECT 1 UNION (SELECT 2)",
		"(SELECT 1) UNION ALL (SELECT 2)",
		"(SELECT a FROM t ORDER BY a) UNION (SELECT b FROM t)",
		"(SELECT a FROM t LIMIT 1) UNION (SELECT b FROM t)",
		"(SELECT 1) UNION (SELECT 2) ORDER BY 1",
		"(SELECT 1) UNION (SELECT 2) LIMIT 1",
		"(SELECT 1 UNION SELECT 2)",
		"((SELECT 1) UNION (SELECT 2))",
		"(SELECT 1) UNION (SELECT 2) UNION (SELECT 3)",
		"(SELECT 1) INTERSECT (SELECT 2)",
		"(SELECT 1) EXCEPT (SELECT 2)",
		// Trailing ORDER BY / LIMIT on a bare parenthesized query apply to the
		// parenthesized result (TiDB accepts; a trailing locking clause does not
		// — see TestParenthesizedQueryRejected).
		"(SELECT a FROM t) ORDER BY a",
		"(SELECT a FROM t) LIMIT 1",
		"(SELECT a FROM t) ORDER BY a LIMIT 1",
		"(SELECT 1) ORDER BY 1",
	}
	for _, sql := range accepted {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParenthesizedQueryDelegatedContexts covers contexts that reach the
// statement query parser via parseSelectStmt — INSERT source, CREATE VIEW,
// CTE, CREATE TABLE ... AS, and a routine body — so a parenthesized set-op
// query is accepted there too (the handoff repro is the procedure-body case).
// All container-verified accepted on pingcap/tidb:v8.5.0.
//
// The derived-table case `FROM ((SELECT 1) UNION (SELECT 2)) x` is now handled
// via a backtracking classifier in parseTableFactor (see TestParenFromDerived).
// Still not covered (subquery-expression production, tracked as PR3): a subquery
// expression `IN ((SELECT 1) UNION (SELECT 2))`, parsed by parseParenExpr.
func TestParenthesizedQueryDelegatedContexts(t *testing.T) {
	accepted := []string{
		"INSERT INTO t (SELECT 1) UNION (SELECT 2)",
		"INSERT INTO t (SELECT 1 UNION SELECT 2)",
		"CREATE VIEW v AS (SELECT 1) UNION (SELECT 2)",
		"WITH cte AS ((SELECT 1) UNION (SELECT 2)) SELECT * FROM cte",
		"CREATE TABLE tt AS (SELECT 1) UNION (SELECT 2)",
		"CREATE PROCEDURE p() BEGIN (SELECT 1) UNION (SELECT 2); END",
		// WITH clause followed by a parenthesized query body.
		"WITH cte AS (SELECT 1) (SELECT * FROM cte)",
		"WITH cte AS (SELECT 1) (SELECT * FROM cte) UNION (SELECT 2)",
		"WITH cte AS (SELECT 1) (SELECT * FROM cte ORDER BY 1)",
	}
	for _, sql := range accepted {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParenthesizedQueryInsertSource pins INSERT / REPLACE with a parenthesized
// query source, including a preceding explicit column list and a parenthesized
// set-op source — all container-verified accepted on pingcap/tidb:v8.5.0. A
// column list never starts with SELECT / WITH / '(', so the source is
// unambiguous after a one-token peek.
func TestParenthesizedQueryInsertSource(t *testing.T) {
	for _, sql := range []string{
		"INSERT INTO t ((SELECT 1) UNION (SELECT 2))",
		"INSERT INTO t (a) (SELECT 1)",
		"INSERT INTO t (a) ((SELECT 1) UNION (SELECT 2))",
		"INSERT INTO t (a, b) ((SELECT 1, 2))",
		"INSERT INTO t (a) (((SELECT 1)))",
		"INSERT INTO t (a) (SELECT 1) UNION (SELECT 2)",
		"REPLACE INTO t (a) ((SELECT 1) UNION (SELECT 2))",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// Deferred (separate follow-up, not yet supported): a parenthesized TABLE or
// VALUES table-value constructor as a query primary — "(TABLE t)",
// "(VALUES ROW(1))", "(SELECT 1) UNION (TABLE t)". TiDB v8.5.0 parses these;
// omni only models a SELECT as a parenthesizable primary. Tracked with the
// derived-table / IN-subquery contexts above.

// TestParenthesizedQueryRejected pins malformed parenthesized queries that TiDB
// v8.5.0 rejects (1064), so the new leading-'(' handling does not over-accept.
func TestParenthesizedQueryRejected(t *testing.T) {
	for _, sql := range []string{
		"(SELECT 1",                    // unclosed paren
		"()",                           // empty parens
		"(SELECT 1) UNION",             // dangling set operation
		"(SELECT 1) UNION ()",          // empty right operand
		"(SELECT a FROM t) FOR UPDATE", // TiDB rejects a locking clause after a parenthesized query
		// Duplicate ORDER BY / LIMIT at the outer (wrapper / set-op) level is
		// rejected, just like on a bare SELECT.
		"(SELECT 1) ORDER BY 1 ORDER BY 2",
		"(SELECT 1) LIMIT 1 LIMIT 2",
		"SELECT 1 UNION SELECT 2 ORDER BY 1 ORDER BY 1",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseExpectError(t, sql)
		})
	}
}

// TestParenthesizedQueryTwoScopes pins the inner-vs-outer ORDER BY / LIMIT
// scoping that a parenthesized query introduces. TiDB v8.5.0 accepts an OUTER
// ORDER BY / LIMIT applied to a parenthesized query that ALREADY carries its
// own ORDER BY / LIMIT (unlike PostgreSQL, which rejects "multiple ORDER BY
// clauses") — both scopes are preserved. All container-verified accepted on
// pingcap/tidb:v8.5.0.
func TestParenthesizedQueryTwoScopes(t *testing.T) {
	for _, sql := range []string{
		"(SELECT 1 LIMIT 5) LIMIT 2",
		"(SELECT 1 ORDER BY 1) ORDER BY 2",
		"((SELECT 1 ORDER BY 1)) ORDER BY 2",
		"(SELECT 1) ORDER BY 1 LIMIT 2 OFFSET 1",
		"(SELECT 1) UNION ((SELECT 2) LIMIT 1) LIMIT 5",
		"((SELECT 1) UNION (SELECT 2)) ORDER BY 1",
		"((SELECT 1) UNION (SELECT 2)) LIMIT 1",
		"(SELECT a FROM t LIMIT 1) UNION (SELECT b FROM t LIMIT 2) ORDER BY 1",
		"(SELECT 1) UNION ((SELECT 2) ORDER BY 1)",
		"(SELECT 1 ORDER BY 1 LIMIT 3) UNION (SELECT 2) ORDER BY 1 LIMIT 9",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParenthesizedQueryLocking pins TiDB v8.5.0's locking-placement rule:
// FOR UPDATE binds to a naked SELECT leaf, so it is accepted inside parens and
// directly after an UNPARENTHESIZED set-op operand, but is rejected as a
// trailing clause after a parenthesized query or a parenthesized set-op (it has
// no leaf to bind to past the closing paren). All container-verified on
// pingcap/tidb:v8.5.0 — porting PostgreSQL's outer-level locking would wrongly
// reject the accepted cases below.
func TestParenthesizedQueryLocking(t *testing.T) {
	accepted := []string{
		"(SELECT a FROM t FOR UPDATE)",
		"(SELECT a FROM t FOR UPDATE) UNION (SELECT b FROM t)",
		"SELECT a FROM t UNION SELECT b FROM t FOR UPDATE",
		"SELECT a FROM t UNION SELECT b FROM t ORDER BY 1 FOR UPDATE",
		"SELECT a FROM t FOR UPDATE UNION SELECT b FROM t",
		"(SELECT a FROM t) UNION SELECT b FROM t FOR UPDATE",
		"SELECT a FROM t ORDER BY a LIMIT 1 FOR UPDATE",
		"(SELECT a FROM t ORDER BY a LIMIT 1 FOR UPDATE)",
		"SELECT 1 UNION (SELECT 2 ORDER BY 1)",
	}
	for _, sql := range accepted {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
	rejected := []string{
		"(SELECT a FROM t) UNION (SELECT b FROM t) FOR UPDATE",
		"((SELECT a FROM t) UNION (SELECT b FROM t)) FOR UPDATE",
		"(SELECT a FROM t) UNION (SELECT b FROM t) ORDER BY 1 FOR UPDATE",
	}
	for _, sql := range rejected {
		t.Run(sql, func(t *testing.T) {
			ParseExpectError(t, sql)
		})
	}
}

// PR2: derived tables whose body is a parenthesized / set-op query expression.
// All container-verified accepted on pingcap/tidb:v8.5.0. Before this change
// parseTableFactor consumed '(' then peeked kwSELECT||kwWITH, so a leading
// nested '(' misrouted to parseTableReferenceList and rejected.
func TestParenFromDerived(t *testing.T) {
	accepted := []string{
		"SELECT * FROM ((SELECT 1) UNION (SELECT 2)) x",
		"SELECT * FROM ((SELECT 1)) x",
		"SELECT * FROM (((SELECT 1))) x",
		"SELECT * FROM ((SELECT 1) UNION (SELECT 2) ORDER BY 1) x",
		"SELECT * FROM ((SELECT 1 UNION SELECT 2)) x",
		"SELECT * FROM ((SELECT a FROM t) UNION (SELECT b FROM t)) x",
	}
	for _, sql := range accepted {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestParenFromDerivedOuterClause covers a derived-table body that is a
// parenthesized query with an OUTER ORDER BY / LIMIT (no set op). TiDB v8.5.0
// accepts all of these (container-verified); the classifier must inherit the
// nested verdict across a trailing ORDER/LIMIT, not just across a set-op
// keyword. This is the bare sibling of the already-passing
// `((SELECT 1) UNION (SELECT 2) ORDER BY 1) x`.
func TestParenFromDerivedOuterClause(t *testing.T) {
	accepted := []string{
		"SELECT * FROM ((SELECT 1) ORDER BY 1) x",
		"SELECT * FROM ((SELECT 1) LIMIT 1) x",
		"SELECT * FROM (((SELECT 1) LIMIT 1)) x",
		"SELECT * FROM ((SELECT 1) ORDER BY 1 LIMIT 1) x",
	}
	for _, sql := range accepted {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// Regression: the classifier must preserve every shape that parses today.
func TestParenFromRegression(t *testing.T) {
	accepted := []string{
		"SELECT * FROM (SELECT 1) x",              // simple derived table
		"SELECT * FROM (t1)",                      // parenthesized single table
		"SELECT * FROM (t1 JOIN t2 ON t1.a=t2.b)", // parenthesized join
		"SELECT * FROM (t1, t2)",                  // parenthesized comma cross-join
		"SELECT * FROM t1, LATERAL (SELECT 1) x",  // LATERAL derived table
		"(SELECT 1) UNION (SELECT 2)",             // top-level set-op (statement)
	}
	for _, sql := range accepted {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestParenFromRobustness exercises the trickier paths: a conditional comment
// spliced inside the body during the classifier scan (validates the snapshot's
// input capture), a set-op derived table followed by a JOIN, and a set-op body
// with no alias. All container-verified parse on pingcap/tidb:v8.5.0.
func TestParenFromRobustness(t *testing.T) {
	accepted := []string{
		"SELECT * FROM ((SELECT 1) UNION (/*!40000 SELECT 2 */)) x",
		"SELECT * FROM ((SELECT 1) UNION (SELECT 2)) x JOIN t3 ON x.a=t3.a",
		"SELECT * FROM ((SELECT 1) UNION (SELECT 2))",
	}
	for _, sql := range accepted {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestParenFromRejected: malformed paren bodies must still be rejected. Both are
// 1064 syntax errors on the TiDB v8.5.0 container.
func TestParenFromRejected(t *testing.T) {
	rejected := []string{
		"SELECT * FROM ((SELECT 1) UNION (SELECT 2)", // unbalanced parens
		"SELECT * FROM ((SELECT 1) (SELECT 2)) x",    // two selects, no set-op
		"SELECT * FROM ((SELECT 1) FOR UPDATE) x",    // FOR UPDATE binds to the leaf, not the outer level (TiDB rejects)
	}
	for _, sql := range rejected {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}
