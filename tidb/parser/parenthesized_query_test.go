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
// Not yet covered (separate productions, tracked as a follow-up): a derived
// table `FROM ((SELECT 1) UNION (SELECT 2)) x` and a subquery expression
// `IN ((SELECT 1) UNION (SELECT 2))`, parsed by parseTableFactor / parseParenExpr.
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
