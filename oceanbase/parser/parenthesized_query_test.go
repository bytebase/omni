package parser

import "testing"

// TestParenthesizedQueryExpression pins MySQL 8.0 parenthesized query
// expressions at statement level: a query expression may start with '(',
// parenthesized set-op operands may carry their own clauses, and trailing
// ORDER BY / LIMIT after the closing ')' bind to the outer query expression.
func TestParenthesizedQueryExpression(t *testing.T) {
	for _, sql := range []string{
		"(SELECT 1)",
		"((SELECT 1))",
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
		"(SELECT a FROM t) ORDER BY a",
		"(SELECT a FROM t) LIMIT 1",
		"(SELECT a FROM t) ORDER BY a LIMIT 1",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParenthesizedQueryDelegatedContexts covers parser contexts that delegate
// to the query-expression parser instead of starting directly at parseStmt.
func TestParenthesizedQueryDelegatedContexts(t *testing.T) {
	for _, sql := range []string{
		"INSERT INTO t (SELECT 1) UNION (SELECT 2)",
		"INSERT INTO t (SELECT 1 UNION SELECT 2)",
		"INSERT INTO t (a) (SELECT 1)",
		"INSERT INTO t (a) ((SELECT 1) UNION (SELECT 2))",
		"REPLACE INTO t (a) ((SELECT 1) UNION (SELECT 2))",
		"CREATE VIEW v AS (SELECT 1) UNION (SELECT 2)",
		"WITH cte AS ((SELECT 1) UNION (SELECT 2)) SELECT * FROM cte",
		"WITH cte AS (SELECT 1) (SELECT * FROM cte)",
		"WITH cte AS (SELECT 1) (SELECT * FROM cte) UNION (SELECT 2)",
		"CREATE TABLE tt AS (SELECT 1) UNION (SELECT 2)",
		"CREATE TABLE tt (SELECT 1) UNION (SELECT 2)",
		"CREATE PROCEDURE p() BEGIN (SELECT 1) UNION (SELECT 2); END",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParenthesizedQueryTwoScopes ensures the AST can preserve both the inner
// clause scope and the outer clause scope. A flat representation would overwrite
// one of these LIMIT / ORDER BY clauses.
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

func TestParenthesizedQueryRejected(t *testing.T) {
	for _, sql := range []string{
		"(SELECT 1",
		"()",
		"(SELECT 1) UNION",
		"(SELECT 1) UNION ()",
		"(SELECT 1) ORDER BY 1 ORDER BY 2",
		"(SELECT 1) LIMIT 1 LIMIT 2",
		"SELECT 1 UNION SELECT 2 ORDER BY 1 ORDER BY 1",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseExpectError(t, sql)
		})
	}
}

// TestParenthesizedQueryLocking pins MySQL 8.0 locking placement. Unlike TiDB,
// MySQL accepts FOR UPDATE after a parenthesized query expression and after a
// parenthesized set operation.
func TestParenthesizedQueryLocking(t *testing.T) {
	for _, sql := range []string{
		"(SELECT a FROM t FOR UPDATE)",
		"(SELECT a FROM t) FOR UPDATE",
		"(SELECT a FROM t) LOCK IN SHARE MODE",
		"(SELECT a FROM t FOR UPDATE) UNION (SELECT a FROM t)",
		"(SELECT a FROM t) UNION (SELECT a FROM t) FOR UPDATE",
		"((SELECT a FROM t) UNION (SELECT a FROM t)) FOR UPDATE",
		"(SELECT a FROM t) UNION (SELECT a FROM t) ORDER BY 1 FOR UPDATE",
		"SELECT a FROM t UNION SELECT a FROM t FOR UPDATE",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParenthesizedQueryMysqlPrimaries tracks MySQL query primaries that TiDB's
// parenthesized-query PR explicitly deferred. These require SelectStmt to model
// non-SELECT primaries as operands, not just SELECT leaves.
func TestParenthesizedQueryMysqlPrimaries(t *testing.T) {
	for _, sql := range []string{
		"(TABLE t)",
		"(VALUES ROW(1))",
		"(SELECT 1) UNION (TABLE t)",
		"(TABLE t) UNION (SELECT 1)",
		"(SELECT 1) UNION (VALUES ROW(2))",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParenthesizedQueryNestedDerivedAndSubqueries tracks contexts that need
// token-stream backtracking: after the first '(' the parser has to distinguish a
// parenthesized query expression from a parenthesized table/value/expression
// list, and one-token lookahead is not enough.
func TestParenthesizedQueryNestedDerivedAndSubqueries(t *testing.T) {
	for _, sql := range []string{
		"SELECT * FROM ((SELECT 1) UNION (SELECT 2)) AS x",
		"SELECT 1 WHERE 1 IN ((SELECT 1) UNION (SELECT 2))",
		"SELECT (SELECT 1 UNION SELECT 2)",
		"SELECT EXISTS ((SELECT 1) UNION (SELECT 2))",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}
