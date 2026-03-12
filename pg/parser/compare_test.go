package parser

import (
	"testing"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/yacc"
)

// CompareWithYacc parses sql with both the recursive descent parser and the yacc parser,
// then compares the AST output. This is the primary validation mechanism.
func CompareWithYacc(t *testing.T, sql string) {
	t.Helper()

	yaccResult, yaccErr := yacc.Parse(sql)
	rdResult, rdErr := Parse(sql)

	// Both should succeed or both should fail.
	if (yaccErr != nil) != (rdErr != nil) {
		t.Fatalf("Parse(%q): yacc err=%v, rd err=%v", sql, yaccErr, rdErr)
	}
	if yaccErr != nil {
		return // Both failed, OK.
	}

	// Compare item counts.
	yaccLen := 0
	if yaccResult != nil {
		yaccLen = len(yaccResult.Items)
	}
	rdLen := 0
	if rdResult != nil {
		rdLen = len(rdResult.Items)
	}
	if yaccLen != rdLen {
		t.Fatalf("Parse(%q): yacc returned %d items, rd returned %d items", sql, yaccLen, rdLen)
	}

	// Compare each statement's AST string representation.
	for i := 0; i < yaccLen; i++ {
		yaccStr := ast.NodeToString(yaccResult.Items[i])
		rdStr := ast.NodeToString(rdResult.Items[i])
		if yaccStr != rdStr {
			t.Errorf("Parse(%q) stmt[%d] mismatch:\n  yacc: %s\n  rd:   %s", sql, i, yaccStr, rdStr)
		}
	}
}

// TestCompareBasic runs basic SELECT comparison tests.
func TestCompareBasic(t *testing.T) {
	tests := []string{
		"SELECT 1",
		"SELECT 1, 2, 3",
		"SELECT 'hello'",
		"SELECT TRUE",
		"SELECT FALSE",
		"SELECT NULL",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareExpressions runs expression comparison tests for batch 3.
func TestCompareExpressions(t *testing.T) {
	tests := []string{
		// Arithmetic
		"SELECT 1 + 2",
		"SELECT 1 - 2",
		"SELECT 2 * 3",
		"SELECT 6 / 2",
		"SELECT 7 % 3",
		"SELECT 2 ^ 3",
		"SELECT -1",
		"SELECT +1",

		// Comparison
		"SELECT 1 = 2",
		"SELECT 1 <> 2",
		"SELECT 1 != 2",
		"SELECT 1 < 2",
		"SELECT 1 > 2",
		"SELECT 1 <= 2",
		"SELECT 1 >= 2",

		// Boolean
		"SELECT TRUE AND FALSE",
		"SELECT TRUE OR FALSE",
		"SELECT NOT TRUE",
		"SELECT 1 = 1 AND 2 = 2",
		"SELECT 1 = 1 OR 2 = 2",

		// IS tests
		"SELECT NULL IS NULL",
		"SELECT 1 IS NOT NULL",
		"SELECT TRUE IS TRUE",
		"SELECT TRUE IS NOT TRUE",
		"SELECT TRUE IS FALSE",
		"SELECT TRUE IS NOT FALSE",
		"SELECT TRUE IS UNKNOWN",
		"SELECT TRUE IS NOT UNKNOWN",
		"SELECT 1 IS DISTINCT FROM 2",
		"SELECT 1 IS NOT DISTINCT FROM 2",

		// BETWEEN
		"SELECT 5 BETWEEN 1 AND 10",
		"SELECT 5 NOT BETWEEN 1 AND 10",
		"SELECT 5 BETWEEN SYMMETRIC 1 AND 10",

		// IN
		"SELECT 1 IN (1, 2, 3)",
		"SELECT 1 NOT IN (1, 2, 3)",

		// LIKE / ILIKE
		"SELECT 'foo' LIKE 'f%'",
		"SELECT 'foo' NOT LIKE 'f%'",
		"SELECT 'foo' ILIKE 'F%'",
		"SELECT 'foo' NOT ILIKE 'F%'",
		"SELECT 'foo' LIKE 'f%' ESCAPE '\\'",

		// SIMILAR TO
		"SELECT 'foo' SIMILAR TO 'f.*'",
		"SELECT 'foo' NOT SIMILAR TO 'f.*'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareCast runs CAST and typecast comparison tests.
func TestCompareCast(t *testing.T) {
	tests := []string{
		"SELECT CAST(1 AS text)",
		"SELECT CAST('123' AS integer)",
		"SELECT 1::text",
		"SELECT '2024-01-01'::date",
		"SELECT 1::int::text",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareCase runs CASE expression comparison tests.
func TestCompareCase(t *testing.T) {
	tests := []string{
		"SELECT CASE WHEN TRUE THEN 1 ELSE 0 END",
		"SELECT CASE WHEN 1 = 1 THEN 'yes' ELSE 'no' END",
		"SELECT CASE WHEN 1 = 1 THEN 'a' WHEN 2 = 2 THEN 'b' ELSE 'c' END",
		"SELECT CASE 1 WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareArray runs ARRAY expression comparison tests.
func TestCompareArray(t *testing.T) {
	tests := []string{
		"SELECT ARRAY[1, 2, 3]",
		"SELECT ARRAY[1, 2, 3][1]",
		"SELECT ARRAY[[1, 2], [3, 4]]",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareRow runs ROW expression comparison tests.
func TestCompareRow(t *testing.T) {
	tests := []string{
		"SELECT ROW(1, 2, 3)",
		"SELECT (1, 2, 3)",
		"SELECT ROW(1, 'a', TRUE)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareSubquery runs subquery expression comparison tests.
func TestCompareSubquery(t *testing.T) {
	tests := []string{
		"SELECT (SELECT 1)",
		"SELECT EXISTS (SELECT 1)",
		"SELECT 1 IN (SELECT 1)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareDistinct runs DISTINCT tests.
func TestCompareDistinct(t *testing.T) {
	tests := []string{
		"SELECT DISTINCT 1",
		"SELECT DISTINCT ON (1) 1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareFuncExprCommon runs func_expr_common_subexpr comparison tests.
func TestCompareFuncExprCommon(t *testing.T) {
	tests := []string{
		"SELECT COALESCE(1, 2, 3)",
		"SELECT NULLIF(1, 2)",
		"SELECT GREATEST(1, 2, 3)",
		"SELECT LEAST(1, 2, 3)",
		"SELECT EXTRACT(YEAR FROM TIMESTAMP '2024-01-01')",
		"SELECT POSITION('b' IN 'abc')",
		"SELECT SUBSTRING('hello' FROM 1 FOR 3)",
		"SELECT TRIM(BOTH ' ' FROM '  hello  ')",
		"SELECT TRIM(LEADING FROM '  hello')",
		"SELECT TRIM(TRAILING FROM 'hello  ')",
		"SELECT OVERLAY('hello' PLACING 'xx' FROM 2 FOR 3)",
		"SELECT CURRENT_DATE",
		"SELECT CURRENT_TIME",
		"SELECT CURRENT_TIMESTAMP",
		"SELECT LOCALTIME",
		"SELECT LOCALTIMESTAMP",
		"SELECT CURRENT_USER",
		"SELECT SESSION_USER",
		"SELECT USER",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareFuncCalls runs function call comparison tests.
func TestCompareFuncCalls(t *testing.T) {
	tests := []string{
		"SELECT count(*)",
		"SELECT count(1)",
		"SELECT count(DISTINCT 1)",
		"SELECT max(1)",
		"SELECT min(1)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareFunctions runs batch 4 function/window/aggregate comparison tests.
func TestCompareFunctions(t *testing.T) {
	tests := []string{
		// Basic function calls (already tested in TestCompareFuncCalls, extended here)
		"SELECT abs(-1)",
		"SELECT length('hello')",
		"SELECT pg_catalog.now()",

		// Aggregate with ORDER BY
		"SELECT array_agg(x ORDER BY x)",
		"SELECT string_agg(x, ',' ORDER BY x DESC)",

		// VARIADIC
		"SELECT concat(VARIADIC ARRAY['a', 'b'])",

		// ALL
		"SELECT count(ALL 1)",

		// Named arguments
		"SELECT make_interval(days => 10)",

		// WITHIN GROUP
		"SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY x)",
		"SELECT percentile_disc(0.5) WITHIN GROUP (ORDER BY x DESC)",

		// FILTER clause
		"SELECT count(*) FILTER (WHERE x > 0)",
		"SELECT sum(x) FILTER (WHERE x > 0)",

		// OVER clause - named window
		"SELECT count(*) OVER w",
		"SELECT sum(x) OVER w",

		// OVER clause - inline window specification
		"SELECT row_number() OVER ()",
		"SELECT row_number() OVER (ORDER BY x)",
		"SELECT rank() OVER (PARTITION BY a ORDER BY b)",
		"SELECT sum(x) OVER (PARTITION BY a)",

		// Window with frame clause
		"SELECT sum(x) OVER (ORDER BY x ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)",
		"SELECT sum(x) OVER (ORDER BY x RANGE BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING)",
		"SELECT sum(x) OVER (ORDER BY x ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING)",
		"SELECT sum(x) OVER (ORDER BY x GROUPS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)",

		// Window with frame exclusion
		"SELECT sum(x) OVER (ORDER BY x ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW EXCLUDE CURRENT ROW)",
		"SELECT sum(x) OVER (ORDER BY x ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW EXCLUDE GROUP)",
		"SELECT sum(x) OVER (ORDER BY x ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW EXCLUDE TIES)",
		"SELECT sum(x) OVER (ORDER BY x ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW EXCLUDE NO OTHERS)",

		// Combined FILTER + OVER
		"SELECT count(*) FILTER (WHERE x > 0) OVER (PARTITION BY a)",

		// WITHIN GROUP + FILTER + OVER
		"SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY x) FILTER (WHERE x > 0)",

		// WINDOW clause in SELECT
		"SELECT sum(x) OVER w FROM t WINDOW w AS (ORDER BY x)",
		"SELECT sum(x) OVER w FROM t WINDOW w AS (PARTITION BY a ORDER BY b)",
		"SELECT sum(x) OVER w1, avg(y) OVER w2 FROM t WINDOW w1 AS (ORDER BY x), w2 AS (ORDER BY y)",

		// Window with referenced window
		"SELECT sum(x) OVER (w ORDER BY x) FROM t WINDOW w AS (PARTITION BY a)",

		// func_expr_common_subexpr extras
		"SELECT GROUPING(a, b)",
		"SELECT CURRENT_ROLE",
		"SELECT CURRENT_CATALOG",
		"SELECT CURRENT_SCHEMA",
		"SELECT CURRENT_TIME(3)",
		"SELECT CURRENT_TIMESTAMP(3)",
		"SELECT LOCALTIME(3)",
		"SELECT LOCALTIMESTAMP(3)",
		"SELECT NORMALIZE('hello')",
		"SELECT NORMALIZE('hello', NFC)",
		"SELECT COLLATION FOR ('hello')",
		"SELECT TREAT('hello' AS text)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareParenthesized runs parenthesized expression tests.
func TestCompareParenthesized(t *testing.T) {
	tests := []string{
		"SELECT (1)",
		"SELECT (1 + 2)",
		"SELECT (1 + 2) * 3",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareColumnRef runs column reference tests.
func TestCompareColumnRef(t *testing.T) {
	tests := []string{
		"SELECT a",
		"SELECT a, b, c",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareParam runs parameter reference tests.
func TestCompareParam(t *testing.T) {
	tests := []string{
		"SELECT $1",
		"SELECT $1, $2",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareXml runs XML expression comparison tests for batch 5.
func TestCompareXml(t *testing.T) {
	tests := []string{
		// XMLCONCAT
		"SELECT XMLCONCAT('<a/>', '<b/>')",

		// XMLELEMENT
		"SELECT XMLELEMENT(NAME foo)",
		"SELECT XMLELEMENT(NAME foo, XMLATTRIBUTES('bar' AS baz))",
		"SELECT XMLELEMENT(NAME foo, 'content')",
		"SELECT XMLELEMENT(NAME foo, XMLATTRIBUTES(1 AS id), 'content')",

		// XMLEXISTS
		"SELECT XMLEXISTS('//foo' PASSING '<foo/>')",
		"SELECT XMLEXISTS('//foo' PASSING BY REF '<foo/>')",

		// XMLFOREST
		"SELECT XMLFOREST(1 AS id, 'hello' AS name)",

		// XMLPARSE
		"SELECT XMLPARSE(DOCUMENT '<doc/>')",
		"SELECT XMLPARSE(CONTENT '<doc/>')",
		"SELECT XMLPARSE(DOCUMENT '<doc/>' PRESERVE WHITESPACE)",
		"SELECT XMLPARSE(DOCUMENT '<doc/>' STRIP WHITESPACE)",

		// XMLPI
		"SELECT XMLPI(NAME php)",
		"SELECT XMLPI(NAME php, 'echo 1')",

		// XMLROOT
		"SELECT XMLROOT('<doc/>', VERSION '1.0')",
		"SELECT XMLROOT('<doc/>', VERSION NO VALUE)",
		"SELECT XMLROOT('<doc/>', VERSION '1.0', STANDALONE YES)",
		"SELECT XMLROOT('<doc/>', VERSION '1.0', STANDALONE NO)",
		"SELECT XMLROOT('<doc/>', VERSION '1.0', STANDALONE NO VALUE)",

		// XMLSERIALIZE
		"SELECT XMLSERIALIZE(DOCUMENT '<doc/>' AS text)",
		"SELECT XMLSERIALIZE(CONTENT '<doc/>' AS text)",
		"SELECT XMLSERIALIZE(DOCUMENT '<doc/>' AS text INDENT)",
		"SELECT XMLSERIALIZE(DOCUMENT '<doc/>' AS text NO INDENT)",

		// IS DOCUMENT (already in expr.go, verify it still works)
		"SELECT '<doc/>' IS DOCUMENT",
		"SELECT '<doc/>' IS NOT DOCUMENT",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareJson runs JSON expression comparison tests for batch 6.
func TestCompareJson(t *testing.T) {
	tests := []string{
		// JSON constructors
		"SELECT JSON_OBJECT('key' VALUE 'value')",
		"SELECT JSON_OBJECT('key' : 'value')",
		"SELECT JSON_OBJECT('a' VALUE 1, 'b' VALUE 2)",
		"SELECT JSON_OBJECT('a' : 1, 'b' : 2)",
		"SELECT JSON_OBJECT()",
		"SELECT JSON_OBJECT(RETURNING text)",

		// JSON_OBJECT null/unique options
		"SELECT JSON_OBJECT('a' VALUE 1 NULL ON NULL)",
		"SELECT JSON_OBJECT('a' VALUE 1 ABSENT ON NULL)",
		"SELECT JSON_OBJECT('a' VALUE 1 WITH UNIQUE KEYS)",

		// JSON_ARRAY
		"SELECT JSON_ARRAY(1, 2, 3)",
		"SELECT JSON_ARRAY()",
		"SELECT JSON_ARRAY(1, 2 ABSENT ON NULL)",
		"SELECT JSON_ARRAY(1, 2 NULL ON NULL)",
		"SELECT JSON_ARRAY(RETURNING text)",

		// JSON()
		"SELECT JSON('{}')",

		// JSON_SCALAR
		"SELECT JSON_SCALAR(42)",
		"SELECT JSON_SCALAR('hello')",

		// JSON_SERIALIZE
		"SELECT JSON_SERIALIZE('{}')",
		"SELECT JSON_SERIALIZE('{}' RETURNING text)",

		// JSON_VALUE
		"SELECT JSON_VALUE('{\"a\":1}', '$.a')",
		"SELECT JSON_VALUE('{\"a\":1}', '$.a' RETURNING int)",
		"SELECT JSON_VALUE('{\"a\":1}', '$.a' NULL ON ERROR)",
		"SELECT JSON_VALUE('{\"a\":1}', '$.a' ERROR ON ERROR)",

		// JSON_QUERY
		"SELECT JSON_QUERY('{\"a\":1}', '$')",
		"SELECT JSON_QUERY('{\"a\":1}', '$' RETURNING text)",
		"SELECT JSON_QUERY('{\"a\":[1,2]}', '$.a' WITH WRAPPER)",
		"SELECT JSON_QUERY('{\"a\":[1,2]}', '$.a' WITHOUT WRAPPER)",
		"SELECT JSON_QUERY('{\"a\":[1,2]}', '$.a' WITH CONDITIONAL WRAPPER)",

		// JSON_EXISTS
		"SELECT JSON_EXISTS('{\"a\":1}', '$.a')",
		"SELECT JSON_EXISTS('{\"a\":1}', '$.a' ERROR ON ERROR)",

		// IS JSON predicates
		"SELECT '{}' IS JSON",
		"SELECT '{}' IS NOT JSON",
		"SELECT '{}' IS JSON VALUE",
		"SELECT '{}' IS JSON OBJECT",
		"SELECT '[]' IS JSON ARRAY",
		"SELECT '1' IS JSON SCALAR",
		"SELECT '{}' IS JSON WITH UNIQUE KEYS",
		"SELECT '{}' IS NOT JSON OBJECT",

		// JSON aggregates
		"SELECT JSON_OBJECTAGG(k VALUE v) FROM t",
		"SELECT JSON_OBJECTAGG(k VALUE v NULL ON NULL) FROM t",
		"SELECT JSON_OBJECTAGG(k VALUE v ABSENT ON NULL) FROM t",
		"SELECT JSON_ARRAYAGG(v) FROM t",
		"SELECT JSON_ARRAYAGG(v ORDER BY v) FROM t",
		"SELECT JSON_ARRAYAGG(v NULL ON NULL) FROM t",
		"SELECT JSON_ARRAYAGG(v ABSENT ON NULL) FROM t",

		// JSON aggregates with FILTER and OVER
		"SELECT JSON_OBJECTAGG(k VALUE v) FILTER (WHERE k IS NOT NULL) FROM t",
		"SELECT JSON_ARRAYAGG(v) OVER (PARTITION BY grp) FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareSelect runs SELECT statement comparison tests for batch 7.
func TestCompareSelect(t *testing.T) {
	tests := []string{
		// Basic FROM
		"SELECT * FROM users",
		"SELECT * FROM public.users",
		"SELECT * FROM users AS u",
		"SELECT * FROM users u",

		// WHERE
		"SELECT * FROM t WHERE id = 1",
		"SELECT * FROM t WHERE id = 1 AND name = 'test'",
		"SELECT * FROM t WHERE a IS NULL",
		"SELECT * FROM t WHERE a IS NOT NULL",
		"SELECT * FROM t WHERE name LIKE '%test%'",
		"SELECT * FROM t WHERE a BETWEEN 1 AND 10",
		"SELECT * FROM t WHERE a IN (1, 2, 3)",

		// ORDER BY
		"SELECT * FROM t ORDER BY id",
		"SELECT * FROM t ORDER BY id DESC",
		"SELECT * FROM t ORDER BY id ASC NULLS FIRST",
		"SELECT * FROM t ORDER BY id DESC NULLS LAST",
		"SELECT * FROM t ORDER BY a, b DESC",

		// LIMIT / OFFSET
		"SELECT * FROM t LIMIT 10",
		"SELECT * FROM t OFFSET 5",
		"SELECT * FROM t LIMIT 10 OFFSET 5",
		"SELECT * FROM t OFFSET 5 LIMIT 10",
		"SELECT * FROM t LIMIT ALL",
		"SELECT * FROM t LIMIT 5 + 5",
		"SELECT * FROM t ORDER BY id LIMIT 10 OFFSET 5",

		// FETCH
		"SELECT * FROM t FETCH FIRST 5 ROWS ONLY",
		"SELECT * FROM t FETCH NEXT ROW ONLY",
		"SELECT * FROM t FETCH FIRST 3 ROWS WITH TIES",

		// DISTINCT
		"SELECT DISTINCT id FROM t",
		"SELECT DISTINCT a, b, c FROM t",
		"SELECT DISTINCT ON (a) a, b FROM t",

		// GROUP BY
		"SELECT a, count(*) FROM t GROUP BY a",
		"SELECT a, count(*) FROM t GROUP BY a HAVING count(*) > 1",
		"SELECT a, b, sum(c) FROM t GROUP BY a, b",
		"SELECT a, count(*) FROM t GROUP BY DISTINCT a",
		"SELECT a, count(*) FROM t GROUP BY ALL a",

		// UNION / INTERSECT / EXCEPT
		"SELECT a FROM t1 UNION SELECT b FROM t2",
		"SELECT a FROM t1 UNION ALL SELECT b FROM t2",
		"SELECT a FROM t1 UNION DISTINCT SELECT b FROM t2",
		"SELECT a FROM t1 INTERSECT SELECT b FROM t2",
		"SELECT a FROM t1 EXCEPT SELECT b FROM t2",
		"SELECT a FROM t1 INTERSECT ALL SELECT b FROM t2",
		"SELECT a FROM t1 EXCEPT ALL SELECT b FROM t2",
		"SELECT a FROM t1 UNION SELECT b FROM t2 UNION SELECT c FROM t3",
		"SELECT a FROM t1 UNION SELECT b FROM t2 ORDER BY a",

		// VALUES
		"VALUES (1, 2, 3)",
		"VALUES (1, 'a'), (2, 'b'), (3, 'c')",

		// TABLE
		"TABLE t",

		// Subquery in FROM
		"SELECT * FROM (SELECT 1) AS sub",

		// JOINs
		"SELECT * FROM t1 CROSS JOIN t2",
		"SELECT * FROM t1 INNER JOIN t2 ON t1.id = t2.id",
		"SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id",
		"SELECT * FROM t1 RIGHT JOIN t2 ON t1.id = t2.id",
		"SELECT * FROM t1 FULL JOIN t2 ON t1.id = t2.id",
		"SELECT * FROM t1 LEFT OUTER JOIN t2 ON t1.id = t2.id",
		"SELECT * FROM t1 RIGHT OUTER JOIN t2 ON t1.id = t2.id",
		"SELECT * FROM t1 FULL OUTER JOIN t2 ON t1.id = t2.id",
		"SELECT * FROM t1 JOIN t2 USING (id)",
		"SELECT * FROM t1 JOIN t2 USING (a, b)",
		"SELECT * FROM t1 NATURAL JOIN t2",
		"SELECT * FROM t1 NATURAL LEFT JOIN t2",
		"SELECT * FROM t1 NATURAL RIGHT JOIN t2",
		"SELECT * FROM t1 NATURAL FULL JOIN t2",
		"SELECT * FROM t1 JOIN t2 ON t1.id = t2.id JOIN t3 ON t2.id = t3.id",
		"SELECT * FROM t1 JOIN t2 ON t1.id = t2.id WHERE t1.active = true",

		// WITH (CTE)
		"WITH cte AS (SELECT 1) SELECT * FROM cte",
		"WITH RECURSIVE cte AS (SELECT 1 UNION ALL SELECT n FROM cte) SELECT * FROM cte",
		"WITH cte1 AS (SELECT 1), cte2 AS (SELECT 2) SELECT * FROM cte1, cte2",
		"WITH cte(a, b) AS (SELECT 1, 2) SELECT * FROM cte",
		"WITH cte AS (SELECT 1 AS n) SELECT * FROM cte ORDER BY n",
		"WITH cte AS (SELECT 1) SELECT * FROM cte LIMIT 5",

		// Subqueries in expressions
		"SELECT (SELECT 1)",
		"SELECT EXISTS (SELECT 1)",
		"SELECT * FROM t WHERE id IN (SELECT id FROM t2)",
		"SELECT * FROM t WHERE id NOT IN (SELECT id FROM t2)",
		"SELECT * FROM t WHERE id = ANY (SELECT id FROM t2)",
		"SELECT * FROM t WHERE id = ALL (SELECT id FROM t2)",
		"SELECT * FROM t WHERE id IN (SELECT id FROM t2 LIMIT 5)",
		"SELECT * FROM t WHERE EXISTS (SELECT 1 FROM t2)",
		"SELECT (SELECT 1) + 2",

		// WINDOW clause
		"SELECT sum(x) OVER w FROM t WINDOW w AS (ORDER BY x)",

		// Combined features
		"SELECT a, count(*) FROM t WHERE a > 0 GROUP BY a HAVING count(*) > 1 ORDER BY a LIMIT 10",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareInsert runs INSERT statement comparison tests for batch 8.
func TestCompareInsert(t *testing.T) {
	tests := []string{
		// Basic INSERT
		"INSERT INTO t VALUES (1, 'hello')",
		"INSERT INTO t (id, name) VALUES (1, 'hello')",
		"INSERT INTO t SELECT * FROM t2",
		"INSERT INTO t DEFAULT VALUES",

		// Multiple VALUES rows
		"INSERT INTO t VALUES (1, 'a'), (2, 'b'), (3, 'c')",

		// RETURNING
		"INSERT INTO t VALUES (1) RETURNING id",
		"INSERT INTO t VALUES (1) RETURNING *",
		"INSERT INTO t VALUES (1, 'a') RETURNING id, name",

		// ON CONFLICT DO NOTHING
		"INSERT INTO t VALUES (1) ON CONFLICT DO NOTHING",
		"INSERT INTO t VALUES (1) ON CONFLICT (id) DO NOTHING",
		"INSERT INTO t VALUES (1) ON CONFLICT ON CONSTRAINT t_pkey DO NOTHING",

		// ON CONFLICT DO UPDATE
		"INSERT INTO t VALUES (1) ON CONFLICT DO UPDATE SET name = 'new'",
		"INSERT INTO t VALUES (1) ON CONFLICT DO UPDATE SET name = 'x' WHERE id > 0",
		"INSERT INTO t VALUES (1) ON CONFLICT (id) DO UPDATE SET name = 'new'",

		// WITH clause
		"WITH cte AS (SELECT 1) INSERT INTO t SELECT * FROM cte",

		// Table alias
		"INSERT INTO t AS target VALUES (1)",

		// Schema-qualified
		"INSERT INTO myschema.t VALUES (1)",

		// Columns + RETURNING
		"INSERT INTO t (id) VALUES (1) RETURNING *",

		// INSERT with SELECT subquery
		"INSERT INTO t (a, b) SELECT x, y FROM t2 WHERE z > 0",

		// Multi-column SET in ON CONFLICT
		"INSERT INTO t VALUES (1) ON CONFLICT DO UPDATE SET a = 1, b = 2",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareUpdate runs UPDATE statement comparison tests for batch 9.
func TestCompareUpdate(t *testing.T) {
	tests := []string{
		// Basic UPDATE
		"UPDATE t SET name = 'new' WHERE id = 1",
		"UPDATE t SET name = 'new', age = 30 WHERE id = 1",
		"UPDATE t SET active = false",

		// FROM clause
		"UPDATE t SET name = t2.name FROM t2 WHERE t.id = t2.id",
		"UPDATE t SET name = t2.name FROM t2, t3 WHERE t.id = t2.id AND t2.id = t3.id",

		// RETURNING
		"UPDATE t SET name = 'new' RETURNING id, name",
		"UPDATE t SET name = 'new' RETURNING *",

		// Alias
		"UPDATE t AS target SET name = 'new'",

		// Schema-qualified
		"UPDATE myschema.t SET name = 'new'",

		// WITH clause
		"WITH cte AS (SELECT 1) UPDATE t SET name = 'new'",

		// Combined
		"UPDATE t SET name = 'new' FROM t2 WHERE t.id = t2.id RETURNING *",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareDelete runs DELETE statement comparison tests for batch 9.
func TestCompareDelete(t *testing.T) {
	tests := []string{
		// Basic DELETE
		"DELETE FROM t WHERE id = 1",
		"DELETE FROM t",

		// USING clause
		"DELETE FROM t USING t2 WHERE t.id = t2.id",
		"DELETE FROM t USING t2, t3 WHERE t.id = t2.id AND t2.id = t3.id",

		// RETURNING
		"DELETE FROM t WHERE id = 1 RETURNING id",
		"DELETE FROM t RETURNING *",

		// Alias
		"DELETE FROM t AS target WHERE target.id = 1",

		// Schema-qualified
		"DELETE FROM myschema.t WHERE id = 1",

		// WITH clause
		"WITH cte AS (SELECT 1) DELETE FROM t WHERE id IN (SELECT * FROM cte)",

		// Combined
		"DELETE FROM t USING t2 WHERE t.id = t2.id RETURNING *",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareMerge runs MERGE statement comparison tests for batch 9.
func TestCompareMerge(t *testing.T) {
	tests := []string{
		// Basic MERGE with UPDATE
		"MERGE INTO target t USING source s ON t.id = s.id WHEN MATCHED THEN UPDATE SET col = s.col",

		// MERGE with INSERT
		"MERGE INTO target USING source ON target.id = source.id WHEN NOT MATCHED THEN INSERT (col) VALUES (source.col)",

		// MERGE with DELETE
		"MERGE INTO target t USING source s ON t.id = s.id WHEN MATCHED THEN DELETE",

		// DO NOTHING
		"MERGE INTO target t USING source s ON t.id = s.id WHEN MATCHED THEN DO NOTHING",
		"MERGE INTO target t USING source s ON t.id = s.id WHEN NOT MATCHED THEN DO NOTHING",

		// Multiple clauses
		"MERGE INTO target t USING source s ON t.id = s.id WHEN MATCHED THEN UPDATE SET col = s.col WHEN NOT MATCHED THEN INSERT (col) VALUES (s.col) WHEN MATCHED THEN DELETE",

		// Subquery source
		"MERGE INTO target t USING (SELECT * FROM source) s ON t.id = s.id WHEN MATCHED THEN UPDATE SET col = s.col",

		// Condition
		"MERGE INTO target t USING source s ON t.id = s.id WHEN MATCHED AND t.col > 0 THEN UPDATE SET col = s.col",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareCreateTable runs CREATE TABLE comparison tests for batch 10.
func TestCompareCreateTable(t *testing.T) {
	tests := []string{
		// Basic
		"CREATE TABLE t (id integer, name varchar(100))",
		"CREATE TABLE IF NOT EXISTS t (id integer)",
		"CREATE TABLE myschema.t (id integer)",

		// Column constraints
		"CREATE TABLE t (id integer NOT NULL)",
		"CREATE TABLE t (id integer NULL)",
		"CREATE TABLE t (id integer DEFAULT 0)",
		"CREATE TABLE t (name text DEFAULT 'hello')",
		"CREATE TABLE t (id integer PRIMARY KEY)",
		"CREATE TABLE t (name varchar UNIQUE)",
		"CREATE TABLE t (age integer CHECK (age > 0))",
		"CREATE TABLE t (user_id integer REFERENCES users(id))",
		"CREATE TABLE t (user_id integer REFERENCES users)",
		"CREATE TABLE t (id integer CONSTRAINT pk_id PRIMARY KEY)",
		"CREATE TABLE t (id integer NOT NULL PRIMARY KEY)",
		"CREATE TABLE t (id integer NOT NULL PRIMARY KEY, name varchar(100) NOT NULL DEFAULT 'unknown')",

		// Table constraints
		"CREATE TABLE t (id integer, PRIMARY KEY (id))",
		"CREATE TABLE t (id integer, UNIQUE (id))",
		"CREATE TABLE t (id integer, CHECK (id > 0))",
		"CREATE TABLE t (id integer, FOREIGN KEY (id) REFERENCES t2(id))",
		"CREATE TABLE t (id integer, FOREIGN KEY (id) REFERENCES t2)",
		"CREATE TABLE t (id integer, CONSTRAINT pk_t PRIMARY KEY (id))",

		// Composite keys
		"CREATE TABLE t (a integer, b integer, PRIMARY KEY (a, b))",
		"CREATE TABLE t (a integer, b integer, UNIQUE (a, b))",
		"CREATE TABLE t (a integer, b integer, FOREIGN KEY (a, b) REFERENCES t2(x, y))",

		// Relpersistence
		"CREATE TEMP TABLE t (id integer)",
		"CREATE TEMPORARY TABLE t (id integer)",
		"CREATE UNLOGGED TABLE t (id integer)",
		"CREATE LOCAL TEMP TABLE t (id integer)",
		"CREATE LOCAL TEMPORARY TABLE t (id integer)",

		// Empty table
		"CREATE TABLE t ()",

		// FK actions
		"CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES p(id) ON DELETE CASCADE)",
		"CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES p(id) ON UPDATE RESTRICT)",
		"CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET NULL)",
		"CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET DEFAULT)",
		"CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES p(id) ON DELETE NO ACTION)",
		"CREATE TABLE t (a int, b int, FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET NULL (a, b))",
		"CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET DEFAULT (a))",
		"CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES p(id) ON UPDATE CASCADE ON DELETE SET NULL)",

		// MATCH
		"CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES p(id) MATCH FULL)",

		// CHECK NOT VALID
		"CREATE TABLE t (x int, CHECK (x > 0) NOT VALID)",
		"CREATE TABLE t (x int CHECK (x > 0))",

		// Column REFERENCES with FK actions
		"CREATE TABLE t (a int REFERENCES p(id) ON DELETE CASCADE ON UPDATE SET NULL)",

		// INHERITS
		"CREATE TABLE t (id int) INHERITS (parent)",
		"CREATE TABLE t (id int) INHERITS (parent1, parent2)",

		// ON COMMIT
		"CREATE TEMP TABLE t (id int) ON COMMIT DROP",
		"CREATE TEMP TABLE t (id int) ON COMMIT DELETE ROWS",
		"CREATE TEMP TABLE t (id int) ON COMMIT PRESERVE ROWS",

		// TABLESPACE
		"CREATE TABLE t (id int) TABLESPACE myts",

		// PARTITION BY
		"CREATE TABLE t (id int, name text) PARTITION BY RANGE (id)",
		"CREATE TABLE t (id int) PARTITION BY LIST (id)",
		"CREATE TABLE t (id int) PARTITION BY HASH (id)",

		// LIKE
		"CREATE TABLE t (LIKE other)",
		"CREATE TABLE t (LIKE other INCLUDING ALL)",
		"CREATE TABLE t (LIKE other INCLUDING DEFAULTS EXCLUDING CONSTRAINTS)",

		// GENERATED columns
		"CREATE TABLE t (id int GENERATED ALWAYS AS IDENTITY)",
		"CREATE TABLE t (id int GENERATED BY DEFAULT AS IDENTITY)",
		"CREATE TABLE t (total int GENERATED ALWAYS AS (price * qty) STORED)",

		// DEFERRABLE constraints
		"CREATE TABLE t (id int, FOREIGN KEY (id) REFERENCES p(id) DEFERRABLE)",
		"CREATE TABLE t (id int, FOREIGN KEY (id) REFERENCES p(id) DEFERRABLE INITIALLY DEFERRED)",

		// NO INHERIT
		"CREATE TABLE t (x int, CHECK (x > 0) NO INHERIT)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareAlterTable runs ALTER TABLE comparison tests for batch 11.
func TestCompareAlterTable(t *testing.T) {
	tests := []string{
		// ADD COLUMN
		"ALTER TABLE t ADD COLUMN name text",
		"ALTER TABLE t ADD name text",
		"ALTER TABLE t ADD COLUMN IF NOT EXISTS name text",
		"ALTER TABLE t ADD IF NOT EXISTS name text",

		// DROP COLUMN
		"ALTER TABLE t DROP COLUMN name",
		"ALTER TABLE t DROP COLUMN name CASCADE",
		"ALTER TABLE t DROP COLUMN name RESTRICT",
		"ALTER TABLE t DROP COLUMN IF EXISTS name",
		"ALTER TABLE t DROP COLUMN IF EXISTS name CASCADE",
		"ALTER TABLE t DROP name",
		"ALTER TABLE t DROP name CASCADE",
		"ALTER TABLE t DROP IF EXISTS name",

		// ALTER COLUMN SET/DROP DEFAULT
		"ALTER TABLE t ALTER COLUMN name SET DEFAULT 'unknown'",
		"ALTER TABLE t ALTER COLUMN id SET DEFAULT 0",
		"ALTER TABLE t ALTER COLUMN name DROP DEFAULT",
		"ALTER TABLE t ALTER name SET DEFAULT 42",
		"ALTER TABLE t ALTER name DROP DEFAULT",

		// ALTER COLUMN SET/DROP NOT NULL
		"ALTER TABLE t ALTER COLUMN name SET NOT NULL",
		"ALTER TABLE t ALTER COLUMN name DROP NOT NULL",
		"ALTER TABLE t ALTER name SET NOT NULL",
		"ALTER TABLE t ALTER name DROP NOT NULL",

		// ALTER COLUMN TYPE
		"ALTER TABLE t ALTER COLUMN name TYPE varchar(200)",
		"ALTER TABLE t ALTER COLUMN id TYPE bigint",
		"ALTER TABLE t ALTER name TYPE integer",
		"ALTER TABLE t ALTER COLUMN name SET DATA TYPE varchar(200)",
		"ALTER TABLE t ALTER name SET DATA TYPE integer",

		// ADD/DROP CONSTRAINT
		"ALTER TABLE t ADD CONSTRAINT pk_id PRIMARY KEY (id)",
		"ALTER TABLE t ADD UNIQUE (name)",
		"ALTER TABLE t ADD CHECK (age > 0)",
		"ALTER TABLE t ADD FOREIGN KEY (user_id) REFERENCES users(id)",
		"ALTER TABLE t DROP CONSTRAINT pk_id",
		"ALTER TABLE t DROP CONSTRAINT pk_id CASCADE",
		"ALTER TABLE t DROP CONSTRAINT IF EXISTS pk_id",
		"ALTER TABLE t DROP CONSTRAINT IF EXISTS pk_id CASCADE",

		// OWNER TO
		"ALTER TABLE t OWNER TO newowner",

		// VALIDATE CONSTRAINT
		"ALTER TABLE t VALIDATE CONSTRAINT check_age",

		// IF EXISTS
		"ALTER TABLE IF EXISTS t ADD COLUMN name text",
		"ALTER TABLE IF EXISTS t DROP COLUMN name",

		// Schema-qualified
		"ALTER TABLE myschema.t ADD COLUMN name text",

		// Multiple commands
		"ALTER TABLE t ADD COLUMN a integer, DROP COLUMN b",
		"ALTER TABLE t ADD COLUMN a integer, DROP COLUMN b, ALTER COLUMN c SET NOT NULL",

		// RENAME
		"ALTER TABLE t RENAME TO t2",
		"ALTER TABLE IF EXISTS t RENAME TO t2",
		"ALTER TABLE t RENAME COLUMN old_col TO new_col",
		"ALTER TABLE t RENAME old_col TO new_col",
		"ALTER TABLE t RENAME CONSTRAINT old_con TO new_con",

		// SET/RESET options
		"ALTER TABLE t SET (fillfactor = 70)",
		"ALTER TABLE t RESET (fillfactor)",

		// CLUSTER/TABLESPACE
		"ALTER TABLE t CLUSTER ON idx_name",
		"ALTER TABLE t SET WITHOUT CLUSTER",
		"ALTER TABLE t SET TABLESPACE new_ts",

		// LOGGED/UNLOGGED
		"ALTER TABLE t SET LOGGED",
		"ALTER TABLE t SET UNLOGGED",

		// ROW LEVEL SECURITY
		"ALTER TABLE t ENABLE ROW LEVEL SECURITY",
		"ALTER TABLE t DISABLE ROW LEVEL SECURITY",
		"ALTER TABLE t FORCE ROW LEVEL SECURITY",
		"ALTER TABLE t NO FORCE ROW LEVEL SECURITY",

		// ENABLE/DISABLE TRIGGER
		"ALTER TABLE t ENABLE TRIGGER my_trigger",
		"ALTER TABLE t DISABLE TRIGGER my_trigger",
		"ALTER TABLE t ENABLE TRIGGER ALL",
		"ALTER TABLE t DISABLE TRIGGER ALL",
		"ALTER TABLE t ENABLE TRIGGER USER",
		"ALTER TABLE t DISABLE TRIGGER USER",
		"ALTER TABLE t ENABLE ALWAYS TRIGGER my_trigger",
		"ALTER TABLE t ENABLE REPLICA TRIGGER my_trigger",

		// ENABLE/DISABLE RULE
		"ALTER TABLE t ENABLE RULE my_rule",
		"ALTER TABLE t DISABLE RULE my_rule",
		"ALTER TABLE t ENABLE ALWAYS RULE my_rule",
		"ALTER TABLE t ENABLE REPLICA RULE my_rule",

		// REPLICA IDENTITY
		"ALTER TABLE t REPLICA IDENTITY DEFAULT",
		"ALTER TABLE t REPLICA IDENTITY FULL",
		"ALTER TABLE t REPLICA IDENTITY NOTHING",
		"ALTER TABLE t REPLICA IDENTITY USING INDEX idx",

		// SET STORAGE
		"ALTER TABLE t ALTER COLUMN col SET STORAGE PLAIN",
		"ALTER TABLE t ALTER COLUMN col SET STORAGE DEFAULT",
		"ALTER TABLE t ALTER col SET STORAGE PLAIN",

		// SET STATISTICS
		"ALTER TABLE t ALTER COLUMN col SET STATISTICS 100",
		"ALTER TABLE t ALTER col SET STATISTICS 100",

		// SET COMPRESSION
		"ALTER TABLE t ALTER COLUMN col SET COMPRESSION pglz",
		"ALTER TABLE t ALTER col SET COMPRESSION lz4",

		// INHERIT / NO INHERIT
		"ALTER TABLE t INHERIT parent",
		"ALTER TABLE t NO INHERIT parent",

		// SET WITH/WITHOUT OIDS
		"ALTER TABLE t SET WITHOUT OIDS",

		// NOT OF
		"ALTER TABLE t NOT OF",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareCreateIndex runs CREATE INDEX comparison tests for batch 12.
func TestCompareCreateIndex(t *testing.T) {
	tests := []string{
		// Basic
		"CREATE INDEX idx ON foo (bar)",
		"CREATE INDEX ON foo (bar)",

		// UNIQUE
		"CREATE UNIQUE INDEX idx ON foo (bar)",

		// CONCURRENTLY
		"CREATE INDEX CONCURRENTLY idx ON foo (bar)",

		// IF NOT EXISTS
		"CREATE INDEX IF NOT EXISTS idx ON foo (bar)",

		// USING
		"CREATE INDEX idx ON foo USING btree (bar)",
		"CREATE INDEX idx ON foo USING hash (bar)",
		"CREATE INDEX idx ON foo USING gin (bar)",
		"CREATE INDEX idx ON foo USING gist (bar)",

		// Multiple columns
		"CREATE INDEX idx ON foo (bar, baz)",
		"CREATE INDEX idx ON foo (a, b, c)",

		// ASC/DESC
		"CREATE INDEX idx ON foo (bar ASC)",
		"CREATE INDEX idx ON foo (bar DESC)",

		// NULLS FIRST/LAST
		"CREATE INDEX idx ON foo (bar NULLS FIRST)",
		"CREATE INDEX idx ON foo (bar NULLS LAST)",
		"CREATE INDEX idx ON foo (bar DESC NULLS FIRST)",
		"CREATE INDEX idx ON foo (bar ASC NULLS LAST)",

		// Expression index
		"CREATE INDEX idx ON foo ((col1 + col2))",
		"CREATE INDEX idx ON foo ((lower(name)))",

		// Partial index (WHERE)
		"CREATE INDEX idx ON foo (bar) WHERE bar > 10",
		"CREATE INDEX idx ON foo (bar) WHERE bar IS NOT NULL",

		// Schema-qualified table
		"CREATE INDEX idx ON myschema.foo (bar)",

		// INCLUDE
		"CREATE INDEX idx ON foo (bar) INCLUDE (baz)",
		"CREATE INDEX idx ON foo (bar) INCLUDE (baz, qux)",

		// WITH options
		"CREATE INDEX idx ON foo (bar) WITH (fillfactor = 70)",

		// TABLESPACE
		"CREATE INDEX idx ON foo (bar) TABLESPACE myts",

		// Combined
		"CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx ON foo USING hash (bar DESC NULLS FIRST) WHERE bar IS NOT NULL",
		"CREATE INDEX idx ON foo (bar) INCLUDE (baz) WITH (fillfactor = 70) TABLESPACE myts WHERE bar > 0",

		// COLLATE
		"CREATE INDEX idx ON foo (bar COLLATE \"C\")",

		// ONLY
		"CREATE INDEX idx ON ONLY foo (bar)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareCreateView runs CREATE VIEW comparison tests for batch 13.
func TestCompareCreateView(t *testing.T) {
	tests := []string{
		// Basic CREATE VIEW
		"CREATE VIEW v AS SELECT 1",
		"CREATE VIEW v AS SELECT * FROM t",
		"CREATE VIEW myschema.v AS SELECT 1",

		// OR REPLACE
		"CREATE OR REPLACE VIEW v AS SELECT 1",
		"CREATE OR REPLACE VIEW myschema.v AS SELECT 1",

		// TEMP
		"CREATE TEMP VIEW v AS SELECT 1",
		"CREATE TEMPORARY VIEW v AS SELECT 1",

		// OR REPLACE + TEMP
		"CREATE OR REPLACE TEMP VIEW v AS SELECT 1",

		// Column list
		"CREATE VIEW v (col1, col2) AS SELECT 1, 2",
		"CREATE VIEW v (a, b, c) AS SELECT 1, 2, 3",

		// WITH CHECK OPTION
		"CREATE VIEW v AS SELECT 1 WITH CHECK OPTION",
		"CREATE VIEW v AS SELECT 1 WITH CASCADED CHECK OPTION",
		"CREATE VIEW v AS SELECT 1 WITH LOCAL CHECK OPTION",

		// RECURSIVE VIEW
		"CREATE RECURSIVE VIEW v (n) AS SELECT 1 UNION ALL SELECT n + 1 FROM v WHERE n < 10",
		"CREATE OR REPLACE RECURSIVE VIEW v (n) AS SELECT 1 UNION ALL SELECT n + 1 FROM v WHERE n < 10",

		// CREATE TABLE AS (CTAS)
		"CREATE TABLE t AS SELECT 1",
		"CREATE TABLE t AS SELECT * FROM other",
		"CREATE TABLE IF NOT EXISTS t AS SELECT 1",
		"CREATE TABLE t AS SELECT 1 WITH DATA",
		"CREATE TABLE t AS SELECT 1 WITH NO DATA",
		"CREATE TABLE t (a, b) AS SELECT 1, 2",
		"CREATE TEMP TABLE t AS SELECT 1",
		"CREATE TEMPORARY TABLE t AS SELECT 1",

		// CREATE MATERIALIZED VIEW
		"CREATE MATERIALIZED VIEW mv AS SELECT 1",
		"CREATE MATERIALIZED VIEW mv AS SELECT * FROM t",
		"CREATE MATERIALIZED VIEW IF NOT EXISTS mv AS SELECT 1",
		"CREATE MATERIALIZED VIEW mv AS SELECT 1 WITH NO DATA",
		"CREATE MATERIALIZED VIEW mv AS SELECT 1 WITH DATA",
		"CREATE UNLOGGED MATERIALIZED VIEW mv AS SELECT 1",
		"CREATE MATERIALIZED VIEW mv (a, b) AS SELECT 1, 2",

		// REFRESH MATERIALIZED VIEW
		"REFRESH MATERIALIZED VIEW mv",
		"REFRESH MATERIALIZED VIEW CONCURRENTLY mv",
		"REFRESH MATERIALIZED VIEW mv WITH DATA",
		"REFRESH MATERIALIZED VIEW mv WITH NO DATA",

		// Ensure regular CREATE TABLE still works
		"CREATE TABLE t (id integer, name varchar(100))",
		"CREATE TABLE IF NOT EXISTS t (id integer)",
		"CREATE TEMP TABLE t (id integer)",
		"CREATE UNLOGGED TABLE t (id integer)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareTransaction runs transaction statement comparison tests for batch 24.
func TestCompareTransaction(t *testing.T) {
	tests := []string{
		// BEGIN
		"BEGIN",
		"BEGIN WORK",
		"BEGIN TRANSACTION",

		// START TRANSACTION
		"START TRANSACTION",

		// COMMIT
		"COMMIT",
		"COMMIT WORK",
		"COMMIT TRANSACTION",

		// END (alias for COMMIT)
		"END",
		"END WORK",
		"END TRANSACTION",

		// ROLLBACK
		"ROLLBACK",
		"ROLLBACK WORK",
		"ROLLBACK TRANSACTION",

		// ABORT (alias for ROLLBACK)
		"ABORT",
		"ABORT WORK",
		"ABORT TRANSACTION",

		// AND CHAIN / AND NO CHAIN
		"COMMIT AND CHAIN",
		"COMMIT AND NO CHAIN",
		"ROLLBACK AND CHAIN",
		"ROLLBACK AND NO CHAIN",
		"END AND CHAIN",
		"ABORT AND CHAIN",
		"COMMIT WORK AND CHAIN",
		"ROLLBACK WORK AND NO CHAIN",

		// SAVEPOINT
		"SAVEPOINT sp1",
		"SAVEPOINT my_savepoint",

		// RELEASE
		"RELEASE sp1",
		"RELEASE SAVEPOINT sp1",

		// ROLLBACK TO
		"ROLLBACK TO sp1",
		"ROLLBACK TO SAVEPOINT sp1",
		"ROLLBACK WORK TO sp1",
		"ROLLBACK WORK TO SAVEPOINT sp1",
		"ROLLBACK TRANSACTION TO sp1",
		"ROLLBACK TRANSACTION TO SAVEPOINT sp1",

		// PREPARE TRANSACTION
		"PREPARE TRANSACTION 'txn_id'",

		// COMMIT PREPARED / ROLLBACK PREPARED
		"COMMIT PREPARED 'txn_id'",
		"ROLLBACK PREPARED 'txn_id'",

		// Transaction modes
		"BEGIN ISOLATION LEVEL SERIALIZABLE",
		"BEGIN ISOLATION LEVEL READ COMMITTED",
		"BEGIN ISOLATION LEVEL READ UNCOMMITTED",
		"BEGIN ISOLATION LEVEL REPEATABLE READ",
		"BEGIN READ ONLY",
		"BEGIN READ WRITE",
		"BEGIN DEFERRABLE",
		"BEGIN NOT DEFERRABLE",
		"START TRANSACTION ISOLATION LEVEL SERIALIZABLE",
		"START TRANSACTION READ ONLY",
		"START TRANSACTION READ WRITE",
		"START TRANSACTION DEFERRABLE",
		"START TRANSACTION NOT DEFERRABLE",

		// Multiple transaction modes
		"BEGIN ISOLATION LEVEL SERIALIZABLE, READ ONLY",
		"BEGIN ISOLATION LEVEL SERIALIZABLE READ ONLY",
		"BEGIN READ ONLY, DEFERRABLE",
		"START TRANSACTION ISOLATION LEVEL SERIALIZABLE, READ ONLY, DEFERRABLE",
		"BEGIN TRANSACTION ISOLATION LEVEL READ COMMITTED, READ WRITE, NOT DEFERRABLE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareSet runs SET/SHOW/RESET comparison tests for batch 23.
func TestCompareSet(t *testing.T) {
	tests := []string{
		// Basic SET with TO
		"SET search_path TO public",
		"SET search_path TO public, pg_catalog",

		// SET with =
		"SET search_path = 'public'",

		// SET to DEFAULT
		"SET search_path TO DEFAULT",
		"SET search_path = DEFAULT",

		// SET LOCAL / SET SESSION
		"SET LOCAL timezone = 'UTC'",
		"SET SESSION timezone = 'UTC'",

		// SET TIME ZONE
		"SET TIME ZONE 'UTC'",
		"SET TIME ZONE LOCAL",
		"SET TIME ZONE DEFAULT",
		"SET TIME ZONE 0",

		// SET NAMES
		"SET NAMES 'UTF8'",

		// SET SCHEMA
		"SET SCHEMA 'myschema'",

		// SET ROLE
		"SET ROLE 'admin'",

		// SET SESSION AUTHORIZATION
		"SET SESSION AUTHORIZATION 'user1'",
		"SET SESSION AUTHORIZATION DEFAULT",

		// SET XML OPTION
		"SET XML OPTION DOCUMENT",
		"SET XML OPTION CONTENT",

		// SET TRANSACTION
		"SET TRANSACTION ISOLATION LEVEL SERIALIZABLE",
		"SET TRANSACTION READ ONLY",
		"SET TRANSACTION READ WRITE",
		"SET TRANSACTION DEFERRABLE",
		"SET TRANSACTION NOT DEFERRABLE",

		// SET CONSTRAINTS
		"SET CONSTRAINTS ALL DEFERRED",
		"SET CONSTRAINTS my_fk IMMEDIATE",

		// Dotted var_name
		"SET myapp.debug TO true",

		// SHOW
		"SHOW search_path",
		"SHOW TIME ZONE",
		"SHOW TRANSACTION ISOLATION LEVEL",
		"SHOW SESSION AUTHORIZATION",
		"SHOW ALL",

		// RESET
		"RESET search_path",
		"RESET TIME ZONE",
		"RESET TRANSACTION ISOLATION LEVEL",
		"RESET SESSION AUTHORIZATION",
		"RESET ALL",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareCreateFunction runs CREATE FUNCTION/PROCEDURE comparison tests for batch 14.
func TestCompareCreateFunction(t *testing.T) {
	tests := []string{
		// Basic
		"CREATE FUNCTION foo() RETURNS void LANGUAGE sql AS 'SELECT 1'",
		"CREATE FUNCTION foo() RETURNS integer LANGUAGE sql AS 'SELECT 1'",

		// OR REPLACE
		"CREATE OR REPLACE FUNCTION foo() RETURNS void LANGUAGE sql AS 'SELECT 1'",

		// Parameters
		"CREATE FUNCTION add(a integer, b integer) RETURNS integer LANGUAGE sql AS 'SELECT a + b'",
		"CREATE FUNCTION foo(integer) RETURNS integer LANGUAGE plpgsql AS 'body'",

		// IN/OUT params
		"CREATE FUNCTION foo(IN a integer, OUT b integer) LANGUAGE sql AS 'body'",
		"CREATE FUNCTION foo(INOUT x integer) RETURNS integer LANGUAGE sql AS 'body'",
		"CREATE FUNCTION foo(VARIADIC args integer[]) RETURNS integer LANGUAGE sql AS 'body'",

		// Default params
		"CREATE FUNCTION foo(x integer DEFAULT 0) RETURNS integer LANGUAGE sql AS 'SELECT x'",
		"CREATE FUNCTION foo(x integer = 0) RETURNS integer LANGUAGE sql AS 'SELECT x'",

		// No RETURNS
		"CREATE FUNCTION foo() LANGUAGE sql AS 'SELECT 1'",

		// RETURNS TABLE
		"CREATE FUNCTION foo() RETURNS TABLE (id integer, name text) LANGUAGE sql AS 'SELECT 1, ''a'''",

		// Options
		"CREATE FUNCTION foo() RETURNS integer LANGUAGE sql IMMUTABLE AS 'SELECT 1'",
		"CREATE FUNCTION foo() RETURNS integer LANGUAGE sql STABLE AS 'SELECT 1'",
		"CREATE FUNCTION foo() RETURNS integer LANGUAGE sql VOLATILE AS 'SELECT 1'",
		"CREATE FUNCTION foo() RETURNS integer LANGUAGE sql STRICT AS 'SELECT 1'",
		"CREATE FUNCTION foo() RETURNS integer CALLED ON NULL INPUT LANGUAGE sql AS 'SELECT 1'",
		"CREATE FUNCTION foo() RETURNS integer RETURNS NULL ON NULL INPUT LANGUAGE sql AS 'SELECT 1'",
		"CREATE FUNCTION foo() RETURNS integer SECURITY DEFINER LANGUAGE sql AS 'SELECT 1'",
		"CREATE FUNCTION foo() RETURNS integer SECURITY INVOKER LANGUAGE sql AS 'SELECT 1'",
		"CREATE FUNCTION foo() RETURNS integer LEAKPROOF LANGUAGE sql AS 'SELECT 1'",
		"CREATE FUNCTION foo() RETURNS integer NOT LEAKPROOF LANGUAGE sql AS 'SELECT 1'",
		"CREATE FUNCTION foo() RETURNS integer LANGUAGE sql COST 100 AS 'SELECT 1'",
		"CREATE FUNCTION foo() RETURNS SETOF integer LANGUAGE sql ROWS 1000 AS 'SELECT 1'",
		"CREATE FUNCTION foo() RETURNS integer LANGUAGE sql PARALLEL SAFE AS 'SELECT 1'",

		// AS with two strings (C functions)
		"CREATE FUNCTION foo() RETURNS integer AS 'obj_file', 'link_symbol' LANGUAGE C",

		// PROCEDURE
		"CREATE PROCEDURE do_something() LANGUAGE sql AS 'SELECT 1'",
		"CREATE OR REPLACE PROCEDURE do_something() LANGUAGE sql AS 'SELECT 1'",
		"CREATE PROCEDURE do_something(x integer) LANGUAGE sql AS 'SELECT 1'",

		// Schema-qualified
		"CREATE FUNCTION myschema.foo() RETURNS integer LANGUAGE sql AS 'SELECT 1'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareDrop tests DROP and TRUNCATE statements against the yacc parser.
func TestCompareDrop(t *testing.T) {
	tests := []string{
		// DROP TABLE
		"DROP TABLE foo",
		"DROP TABLE IF EXISTS foo",
		"DROP TABLE foo CASCADE",
		"DROP TABLE foo RESTRICT",
		"DROP TABLE myschema.foo",
		"DROP TABLE foo, bar",
		"DROP TABLE schema1.t1, schema2.t2",
		"DROP TABLE IF EXISTS foo CASCADE",
		"DROP TABLE IF EXISTS foo, bar CASCADE",

		// DROP VIEW
		"DROP VIEW myview",
		"DROP VIEW IF EXISTS v1",
		"DROP VIEW v1, v2",
		"DROP VIEW IF EXISTS v1 CASCADE",

		// DROP INDEX
		"DROP INDEX myindex",
		"DROP INDEX IF EXISTS idx1 CASCADE",
		"DROP INDEX CONCURRENTLY myindex",
		"DROP INDEX CONCURRENTLY IF EXISTS myindex",
		"DROP INDEX CONCURRENTLY IF EXISTS idx1 CASCADE",

		// DROP SEQUENCE
		"DROP SEQUENCE myseq",
		"DROP SEQUENCE IF EXISTS myseq",

		// DROP MATERIALIZED VIEW
		"DROP MATERIALIZED VIEW mymatview",
		"DROP MATERIALIZED VIEW IF EXISTS mymatview",
		"DROP MATERIALIZED VIEW IF EXISTS mymatview CASCADE",

		// DROP FOREIGN TABLE
		"DROP FOREIGN TABLE ft",
		"DROP FOREIGN TABLE IF EXISTS ft",

		// DROP COLLATION
		"DROP COLLATION mycoll",
		"DROP COLLATION IF EXISTS mycoll",

		// DROP CONVERSION
		"DROP CONVERSION myconv",
		"DROP CONVERSION IF EXISTS myconv",

		// DROP STATISTICS
		"DROP STATISTICS mystats",
		"DROP STATISTICS IF EXISTS mystats",

		// DROP TEXT SEARCH objects
		"DROP TEXT SEARCH PARSER myparser",
		"DROP TEXT SEARCH DICTIONARY mydict",
		"DROP TEXT SEARCH TEMPLATE mytmpl",
		"DROP TEXT SEARCH CONFIGURATION myconfig",
		"DROP TEXT SEARCH PARSER IF EXISTS myparser",

		// DROP ACCESS METHOD
		"DROP ACCESS METHOD myam",
		"DROP ACCESS METHOD IF EXISTS myam",

		// DROP TYPE
		"DROP TYPE mytype",
		"DROP TYPE IF EXISTS mytype",
		"DROP TYPE mytype CASCADE",
		"DROP TYPE myschema.mytype",

		// DROP DOMAIN
		"DROP DOMAIN mydomain",
		"DROP DOMAIN IF EXISTS mydomain CASCADE",

		// DROP SCHEMA
		"DROP SCHEMA myschema",
		"DROP SCHEMA IF EXISTS myschema",
		"DROP SCHEMA myschema CASCADE",
		"DROP SCHEMA s1, s2",

		// DROP EXTENSION
		"DROP EXTENSION myext",
		"DROP EXTENSION IF EXISTS myext",
		"DROP EXTENSION myext CASCADE",

		// DROP FUNCTION
		"DROP FUNCTION myfunc",
		"DROP FUNCTION IF EXISTS myfunc",
		"DROP FUNCTION myfunc CASCADE",
		"DROP FUNCTION myschema.myfunc",

		// DROP PROCEDURE
		"DROP PROCEDURE myproc",
		"DROP PROCEDURE IF EXISTS myproc",

		// DROP ROUTINE
		"DROP ROUTINE myroutine",
		"DROP ROUTINE IF EXISTS myroutine",

		// DROP AGGREGATE
		"DROP AGGREGATE myagg",
		"DROP AGGREGATE IF EXISTS myagg",

		// DROP TRIGGER
		"DROP TRIGGER mytrig ON mytable",
		"DROP TRIGGER IF EXISTS mytrig ON mytable",
		"DROP TRIGGER mytrig ON mytable CASCADE",

		// DROP POLICY
		"DROP POLICY mypol ON mytable",
		"DROP POLICY IF EXISTS mypol ON mytable",

		// DROP RULE
		"DROP RULE myrule ON mytable",
		"DROP RULE IF EXISTS myrule ON mytable",

		// DROP EVENT TRIGGER
		"DROP EVENT TRIGGER myevttrig",
		"DROP EVENT TRIGGER IF EXISTS myevttrig",

		// DROP LANGUAGE
		"DROP LANGUAGE plpgsql",
		"DROP LANGUAGE IF EXISTS plpgsql",

		// DROP PUBLICATION
		"DROP PUBLICATION mypub",
		"DROP PUBLICATION IF EXISTS mypub",

		// DROP SUBSCRIPTION
		"DROP SUBSCRIPTION mysub",
		"DROP SUBSCRIPTION IF EXISTS mysub",

		// DROP SERVER
		"DROP SERVER myserver",
		"DROP SERVER IF EXISTS myserver",

		// DROP FOREIGN DATA WRAPPER
		"DROP FOREIGN DATA WRAPPER myfdw",
		"DROP FOREIGN DATA WRAPPER IF EXISTS myfdw",

		// DROP OWNED BY
		"DROP OWNED BY myrole",
		"DROP OWNED BY role1, role2",

		// TRUNCATE
		"TRUNCATE foo",
		"TRUNCATE TABLE foo",
		"TRUNCATE foo, bar",
		"TRUNCATE TABLE foo, bar",
		"TRUNCATE foo RESTART IDENTITY",
		"TRUNCATE foo CONTINUE IDENTITY",
		"TRUNCATE foo CASCADE",
		"TRUNCATE foo RESTRICT",
		"TRUNCATE foo RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE foo RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE IF EXISTS foo, bar",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareSequence runs CREATE/ALTER SEQUENCE comparison tests for batch 17.
func TestCompareSequence(t *testing.T) {
	tests := []string{
		// Basic CREATE SEQUENCE
		"CREATE SEQUENCE seq",
		"CREATE SEQUENCE myschema.seq",
		"CREATE SEQUENCE IF NOT EXISTS seq",

		// TEMP SEQUENCE
		"CREATE TEMP SEQUENCE seq",
		"CREATE TEMPORARY SEQUENCE seq",

		// Sequence options
		"CREATE SEQUENCE seq INCREMENT 5",
		"CREATE SEQUENCE seq INCREMENT BY 5",
		"CREATE SEQUENCE seq START 100",
		"CREATE SEQUENCE seq START WITH 100",
		"CREATE SEQUENCE seq MINVALUE 1",
		"CREATE SEQUENCE seq MAXVALUE 1000",
		"CREATE SEQUENCE seq NO MINVALUE",
		"CREATE SEQUENCE seq NO MAXVALUE",
		"CREATE SEQUENCE seq CYCLE",
		"CREATE SEQUENCE seq NO CYCLE",
		"CREATE SEQUENCE seq CACHE 20",
		"CREATE SEQUENCE seq OWNED BY t.id",

		// Combined options
		"CREATE SEQUENCE seq INCREMENT BY 5 START WITH 100",
		"CREATE SEQUENCE seq MINVALUE 1 MAXVALUE 1000",
		"CREATE SEQUENCE seq INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 100 CYCLE CACHE 10",
		"CREATE SEQUENCE IF NOT EXISTS seq INCREMENT BY 2 NO MINVALUE NO MAXVALUE NO CYCLE",

		// ALTER SEQUENCE
		"ALTER SEQUENCE seq RESTART",
		"ALTER SEQUENCE seq RESTART WITH 100",
		"ALTER SEQUENCE seq INCREMENT BY 10",
		"ALTER SEQUENCE seq INCREMENT 5",
		"ALTER SEQUENCE IF EXISTS seq INCREMENT 5",
		"ALTER SEQUENCE seq MINVALUE 1 MAXVALUE 1000",
		"ALTER SEQUENCE seq NO MINVALUE NO MAXVALUE",
		"ALTER SEQUENCE seq CACHE 50",
		"ALTER SEQUENCE seq CYCLE",
		"ALTER SEQUENCE seq NO CYCLE",
		"ALTER SEQUENCE seq OWNED BY t.id",
		"ALTER SEQUENCE seq START WITH 1",

		// CREATE TYPE ... AS ENUM
		"CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy')",
		"CREATE TYPE myschema.mood AS ENUM ('sad', 'ok', 'happy')",
		"CREATE TYPE empty_enum AS ENUM ()",
		"CREATE TYPE single AS ENUM ('one')",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareDomain runs CREATE DOMAIN comparison tests for batch 17.
func TestCompareDomain(t *testing.T) {
	tests := []string{
		// Basic CREATE DOMAIN
		"CREATE DOMAIN d AS integer",
		"CREATE DOMAIN d integer",
		"CREATE DOMAIN myschema.d AS text",

		// Domain constraints
		"CREATE DOMAIN d AS integer NOT NULL",
		"CREATE DOMAIN d AS text DEFAULT 'hello'",
		"CREATE DOMAIN d AS integer CHECK (VALUE > 0)",
		"CREATE DOMAIN d AS integer NOT NULL DEFAULT 0",
		"CREATE DOMAIN d AS integer CONSTRAINT pos CHECK (VALUE > 0)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareCursor runs comparison tests for DECLARE CURSOR, FETCH, MOVE, and CLOSE.
func TestCompareCursor(t *testing.T) {
	tests := []string{
		// DECLARE CURSOR
		"DECLARE c CURSOR FOR SELECT 1",
		"DECLARE c BINARY CURSOR FOR SELECT 1",
		"DECLARE c SCROLL CURSOR FOR SELECT 1",
		"DECLARE c NO SCROLL CURSOR FOR SELECT 1",
		"DECLARE c INSENSITIVE CURSOR FOR SELECT 1",
		"DECLARE c ASENSITIVE CURSOR FOR SELECT 1",
		"DECLARE c CURSOR WITH HOLD FOR SELECT 1",
		"DECLARE c CURSOR WITHOUT HOLD FOR SELECT 1",
		"DECLARE c BINARY SCROLL CURSOR WITH HOLD FOR SELECT 1",

		// FETCH
		"FETCH c",
		"FETCH FROM c",
		"FETCH IN c",
		"FETCH NEXT FROM c",
		"FETCH PRIOR FROM c",
		"FETCH FIRST FROM c",
		"FETCH LAST FROM c",
		"FETCH ABSOLUTE 5 FROM c",
		"FETCH ABSOLUTE -3 FROM c",
		"FETCH RELATIVE 2 FROM c",
		"FETCH RELATIVE -2 FROM c",
		"FETCH 10 FROM c",
		"FETCH ALL FROM c",
		"FETCH FORWARD FROM c",
		"FETCH FORWARD 5 FROM c",
		"FETCH FORWARD ALL FROM c",
		"FETCH BACKWARD FROM c",
		"FETCH BACKWARD 5 FROM c",
		"FETCH BACKWARD ALL FROM c",

		// MOVE
		"MOVE c",
		"MOVE NEXT FROM c",
		"MOVE PRIOR FROM c",
		"MOVE FIRST FROM c",
		"MOVE LAST FROM c",
		"MOVE ABSOLUTE 5 FROM c",
		"MOVE RELATIVE -1 FROM c",
		"MOVE 10 FROM c",
		"MOVE ALL FROM c",
		"MOVE FORWARD FROM c",
		"MOVE FORWARD 10 FROM c",
		"MOVE FORWARD ALL FROM c",
		"MOVE BACKWARD FROM c",
		"MOVE BACKWARD 5 FROM c",
		"MOVE BACKWARD ALL FROM c",

		// CLOSE
		"CLOSE c",
		"CLOSE ALL",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestComparePrepare runs comparison tests for PREPARE, EXECUTE, and DEALLOCATE statements.
func TestComparePrepare(t *testing.T) {
	tests := []string{
		// PREPARE
		"PREPARE q AS SELECT 1",
		"PREPARE q AS SELECT * FROM t WHERE id = $1",
		"PREPARE q(integer) AS SELECT * FROM t WHERE id = $1",
		"PREPARE q(integer, text) AS SELECT * FROM t WHERE id = $1 AND name = $2",
		"PREPARE ins AS INSERT INTO t VALUES (1, 2)",
		"PREPARE upd AS UPDATE t SET x = 1 WHERE id = $1",
		"PREPARE del AS DELETE FROM t WHERE id = $1",

		// EXECUTE
		"EXECUTE q",
		"EXECUTE q(1)",
		"EXECUTE q(1, 'hello')",
		"EXECUTE q(1, 2, 3)",

		// DEALLOCATE
		"DEALLOCATE q",
		"DEALLOCATE PREPARE q",
		"DEALLOCATE ALL",
		"DEALLOCATE PREPARE ALL",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareLock tests LOCK TABLE and SET CONSTRAINTS statements.
func TestCompareLock(t *testing.T) {
	tests := []string{
		// Basic LOCK
		"LOCK t",
		"LOCK TABLE t",

		// Multiple tables
		"LOCK t1, t2",
		"LOCK TABLE t1, t2, t3",

		// Schema-qualified
		"LOCK myschema.t",
		"LOCK TABLE myschema.t",

		// Lock modes
		"LOCK t IN ACCESS SHARE MODE",
		"LOCK t IN ROW SHARE MODE",
		"LOCK t IN ROW EXCLUSIVE MODE",
		"LOCK t IN SHARE UPDATE EXCLUSIVE MODE",
		"LOCK t IN SHARE MODE",
		"LOCK t IN SHARE ROW EXCLUSIVE MODE",
		"LOCK t IN EXCLUSIVE MODE",
		"LOCK t IN ACCESS EXCLUSIVE MODE",

		// NOWAIT
		"LOCK t NOWAIT",
		"LOCK t IN SHARE MODE NOWAIT",
		"LOCK TABLE t IN ACCESS EXCLUSIVE MODE NOWAIT",

		// Combined
		"LOCK TABLE t1, t2 IN ROW EXCLUSIVE MODE NOWAIT",

		// SET CONSTRAINTS
		"SET CONSTRAINTS ALL DEFERRED",
		"SET CONSTRAINTS ALL IMMEDIATE",
		"SET CONSTRAINTS myconstraint DEFERRED",
		"SET CONSTRAINTS myconstraint IMMEDIATE",
		"SET CONSTRAINTS c1, c2 DEFERRED",
		"SET CONSTRAINTS myschema.c1 IMMEDIATE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareMaintenance runs VACUUM, ANALYZE, CLUSTER, REINDEX comparison tests for batch 27.
func TestCompareMaintenance(t *testing.T) {
	tests := []string{
		// VACUUM
		"VACUUM",
		"VACUUM t",
		"VACUUM t, t2",
		"VACUUM FULL t",
		"VACUUM FREEZE t",
		"VACUUM VERBOSE t",
		"VACUUM FULL FREEZE VERBOSE ANALYZE t",
		"VACUUM (VERBOSE) t",
		"VACUUM (ANALYZE) t",
		"VACUUM (VERBOSE, ANALYZE) t",
		"VACUUM (FULL, FREEZE, VERBOSE, ANALYZE) t",
		"VACUUM (VERBOSE true)",
		"VACUUM (VERBOSE false)",
		"VACUUM myschema.mytable",

		// ANALYZE
		"ANALYZE",
		"ANALYZE t",
		"ANALYZE t, t2",
		"ANALYZE VERBOSE t",
		"ANALYZE t(col1)",
		"ANALYZE t(col1, col2)",
		"ANALYZE (VERBOSE) t",
		"ANALYSE t",
		"ANALYZE myschema.mytable",

		// CLUSTER
		"CLUSTER t",
		"CLUSTER t USING idx",
		"CLUSTER VERBOSE t",
		"CLUSTER VERBOSE t USING idx",
		"CLUSTER (VERBOSE) t",
		"CLUSTER (VERBOSE) t USING idx",
		"CLUSTER (VERBOSE)",
		"CLUSTER VERBOSE",
		"CLUSTER",
		"CLUSTER idx ON t",
		"CLUSTER VERBOSE idx ON t",

		// REINDEX
		"REINDEX TABLE t",
		"REINDEX INDEX idx",
		"REINDEX SCHEMA myschema",
		"REINDEX DATABASE mydb",
		"REINDEX SYSTEM mydb",
		"REINDEX (VERBOSE) TABLE t",
		"REINDEX (VERBOSE) INDEX idx",
		"REINDEX (VERBOSE) SCHEMA myschema",
		"REINDEX (VERBOSE) DATABASE mydb",
		"REINDEX (VERBOSE) SYSTEM mydb",
		"REINDEX TABLE CONCURRENTLY t",
		"REINDEX INDEX CONCURRENTLY idx",
		"REINDEX SCHEMA CONCURRENTLY myschema",
		"REINDEX DATABASE CONCURRENTLY mydb",
		"REINDEX TABLE myschema.mytable",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareDatabase tests CREATE/ALTER/DROP DATABASE statements.
func TestCompareDatabase(t *testing.T) {
	tests := []string{
		// CREATE DATABASE
		"CREATE DATABASE mydb",
		"CREATE DATABASE mydb OWNER myuser",
		"CREATE DATABASE mydb ENCODING 'UTF8'",
		"CREATE DATABASE mydb TEMPLATE template0",
		"CREATE DATABASE mydb OWNER myuser ENCODING 'UTF8' TEMPLATE template0",
		"CREATE DATABASE mydb CONNECTION LIMIT 100",
		"CREATE DATABASE mydb TABLESPACE mytbs",
		"CREATE DATABASE mydb OWNER = myuser",
		"CREATE DATABASE mydb ENCODING = 'UTF8'",
		"CREATE DATABASE mydb CONNECTION LIMIT = 100",

		// ALTER DATABASE
		"ALTER DATABASE mydb SET TABLESPACE newtbs",
		"ALTER DATABASE mydb SET search_path TO public",
		"ALTER DATABASE mydb SET search_path = public",
		"ALTER DATABASE mydb RESET search_path",
		"ALTER DATABASE mydb RESET ALL",
		"ALTER DATABASE mydb CONNECTION LIMIT 100",

		// DROP DATABASE
		"DROP DATABASE mydb",
		"DROP DATABASE IF EXISTS mydb",
		"DROP DATABASE mydb (FORCE)",
		"DROP DATABASE IF EXISTS mydb (FORCE)",
		"DROP DATABASE mydb WITH (FORCE)",
		"DROP DATABASE IF EXISTS mydb WITH (FORCE)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareCopy tests COPY statements against the yacc parser.
func TestCompareCopy(t *testing.T) {
	tests := []string{
		// Basic COPY FROM/TO
		"COPY foo FROM '/tmp/data.csv'",
		"COPY foo TO '/tmp/data.csv'",

		// COPY with columns
		"COPY foo (col1, col2) FROM '/tmp/data.csv'",

		// COPY FROM STDIN / TO STDOUT
		"COPY foo FROM STDIN",
		"COPY foo TO STDOUT",

		// COPY with new-style options (parenthesized)
		"COPY foo FROM '/tmp/data.csv' WITH (FORMAT csv)",
		"COPY foo FROM '/tmp/data.csv' WITH (FORMAT csv, HEADER true)",
		"COPY foo FROM '/tmp/data.csv' (FORMAT csv)",

		// COPY query form
		"COPY (SELECT * FROM foo) TO '/tmp/data.csv'",
		"COPY (SELECT * FROM foo) TO '/tmp/data.csv' WITH (FORMAT csv)",
		"COPY (SELECT 1) TO '/tmp/out.csv' (FORMAT csv)",

		// COPY with WHERE clause
		"COPY foo FROM '/tmp/data.csv' WHERE x > 0",
		"COPY foo FROM '/tmp/data.csv' WITH (FORMAT csv) WHERE x > 0",

		// COPY with PROGRAM
		"COPY foo FROM PROGRAM '/bin/cat /tmp/data.csv'",
		"COPY foo TO PROGRAM '/bin/gzip > /tmp/data.csv.gz'",

		// COPY with schema-qualified table
		"COPY public.foo FROM '/tmp/data.csv'",
		"COPY public.foo (col1) TO STDOUT",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareExtension tests CREATE/ALTER/DROP EXTENSION, ACCESS METHOD, CAST, TRANSFORM statements.
func TestCompareExtension(t *testing.T) {
	tests := []string{
		// CREATE EXTENSION
		"CREATE EXTENSION pgcrypto",
		"CREATE EXTENSION IF NOT EXISTS pgcrypto",
		"CREATE EXTENSION pgcrypto SCHEMA public",
		"CREATE EXTENSION pgcrypto VERSION '1.3'",
		"CREATE EXTENSION pgcrypto CASCADE",
		"CREATE EXTENSION IF NOT EXISTS pgcrypto SCHEMA myschema VERSION '2.0' CASCADE",
		"CREATE EXTENSION pgcrypto WITH SCHEMA public",

		// ALTER EXTENSION UPDATE
		"ALTER EXTENSION pgcrypto UPDATE",
		"ALTER EXTENSION pgcrypto UPDATE TO '2.0'",

		// ALTER EXTENSION ADD/DROP
		"ALTER EXTENSION pgcrypto ADD TABLE t1",
		"ALTER EXTENSION pgcrypto DROP TABLE t1",
		"ALTER EXTENSION pgcrypto ADD SEQUENCE s1",
		"ALTER EXTENSION pgcrypto ADD FUNCTION myfunc(integer)",
		"ALTER EXTENSION pgcrypto ADD FUNCTION myfunc(integer, text)",
		"ALTER EXTENSION pgcrypto DROP FUNCTION myfunc(integer)",
		"ALTER EXTENSION pgcrypto ADD SCHEMA myschema",
		"ALTER EXTENSION pgcrypto ADD TYPE mytype",
		"ALTER EXTENSION pgcrypto ADD DOMAIN mydomain",

		// CREATE ACCESS METHOD
		"CREATE ACCESS METHOD myam TYPE INDEX HANDLER handler_func",
		"CREATE ACCESS METHOD myam TYPE TABLE HANDLER handler_func",
		"CREATE ACCESS METHOD myam TYPE INDEX HANDLER myschema.handler_func",

		// CREATE CAST
		"CREATE CAST (bigint AS integer) WITH FUNCTION int4(bigint)",
		"CREATE CAST (bigint AS integer) WITHOUT FUNCTION",
		"CREATE CAST (bigint AS integer) WITH INOUT",
		"CREATE CAST (bigint AS integer) WITH FUNCTION int4(bigint) AS IMPLICIT",
		"CREATE CAST (bigint AS integer) WITH FUNCTION int4(bigint) AS ASSIGNMENT",
		"CREATE CAST (bigint AS integer) WITHOUT FUNCTION AS IMPLICIT",
		"CREATE CAST (bigint AS integer) WITH INOUT AS ASSIGNMENT",

		// DROP CAST
		"DROP CAST (bigint AS integer)",
		"DROP CAST IF EXISTS (bigint AS integer)",
		"DROP CAST (bigint AS integer) CASCADE",
		"DROP CAST IF EXISTS (bigint AS integer) CASCADE",

		// CREATE TRANSFORM
		"CREATE TRANSFORM FOR integer LANGUAGE plpythonu (FROM SQL WITH FUNCTION int_to_py(integer), TO SQL WITH FUNCTION py_to_int(internal))",
		"CREATE TRANSFORM FOR integer LANGUAGE plpythonu (FROM SQL WITH FUNCTION int_to_py(integer))",
		"CREATE TRANSFORM FOR integer LANGUAGE plpythonu (TO SQL WITH FUNCTION py_to_int(internal))",
		"CREATE TRANSFORM FOR integer LANGUAGE plpythonu (TO SQL WITH FUNCTION py_to_int(internal), FROM SQL WITH FUNCTION int_to_py(integer))",
		"CREATE OR REPLACE TRANSFORM FOR integer LANGUAGE plpythonu (FROM SQL WITH FUNCTION int_to_py(integer), TO SQL WITH FUNCTION py_to_int(internal))",

		// DROP TRANSFORM
		"DROP TRANSFORM FOR integer LANGUAGE plpythonu",
		"DROP TRANSFORM IF EXISTS FOR integer LANGUAGE plpythonu",
		"DROP TRANSFORM FOR integer LANGUAGE plpythonu CASCADE",
		"DROP TRANSFORM IF EXISTS FOR integer LANGUAGE plpythonu CASCADE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareUtility tests utility statements against the yacc parser.
func TestCompareUtility(t *testing.T) {
	tests := []string{
		// CHECKPOINT
		"CHECKPOINT",

		// DISCARD
		"DISCARD ALL",
		"DISCARD TEMP",
		"DISCARD TEMPORARY",
		"DISCARD PLANS",
		"DISCARD SEQUENCES",

		// LISTEN
		"LISTEN mychannel",

		// UNLISTEN
		"UNLISTEN mychannel",
		"UNLISTEN *",

		// NOTIFY
		"NOTIFY mychannel",
		"NOTIFY mychannel, 'hello world'",

		// LOAD
		"LOAD '/usr/lib/libpq.so'",

		// DO
		"DO 'BEGIN RAISE NOTICE ''hello''; END'",
		"DO LANGUAGE plpgsql 'BEGIN NULL; END'",
		"DO 'BEGIN NULL; END' LANGUAGE plpgsql",

		// REASSIGN OWNED
		"REASSIGN OWNED BY alice TO bob",
		"REASSIGN OWNED BY alice, bob TO charlie",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareExplain tests EXPLAIN statements against the yacc parser.
func TestCompareExplain(t *testing.T) {
	tests := []string{
		"EXPLAIN SELECT 1",
		"EXPLAIN ANALYZE SELECT 1",
		"EXPLAIN VERBOSE SELECT 1",
		"EXPLAIN ANALYZE VERBOSE SELECT 1",
		"EXPLAIN (ANALYZE) SELECT 1",
		"EXPLAIN (VERBOSE) SELECT 1",
		"EXPLAIN (ANALYZE, VERBOSE) SELECT 1",
		"EXPLAIN (FORMAT JSON) SELECT 1",
		"EXPLAIN (FORMAT YAML) SELECT 1",
		"EXPLAIN (FORMAT TEXT) SELECT 1",
		"EXPLAIN (FORMAT XML) SELECT 1",
		"EXPLAIN (COSTS off) SELECT 1",
		"EXPLAIN (BUFFERS) SELECT 1",
		"EXPLAIN (TIMING) SELECT 1",
		"EXPLAIN (ANALYZE true, VERBOSE true) SELECT 1",
		"EXPLAIN (ANALYZE true, VERBOSE true, FORMAT JSON) SELECT 1",
		"EXPLAIN INSERT INTO foo VALUES (1)",
		"EXPLAIN DELETE FROM foo",
		"EXPLAIN UPDATE foo SET x = 1",
		"EXPLAIN ANALYZE INSERT INTO foo VALUES (1)",
		"EXPLAIN ANALYZE DELETE FROM foo",
		"EXPLAIN ANALYZE UPDATE foo SET x = 1",
		"EXPLAIN SELECT * FROM foo WHERE id = 1",
		"EXPLAIN SELECT * FROM foo JOIN bar ON foo.id = bar.id",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareFDW tests foreign data wrapper, server, foreign table, user mapping,
// import foreign schema, and CREATE LANGUAGE statements.
func TestCompareFDW(t *testing.T) {
	tests := []string{
		// CREATE FOREIGN DATA WRAPPER
		"CREATE FOREIGN DATA WRAPPER myfdw",
		"CREATE FOREIGN DATA WRAPPER myfdw HANDLER fdw_handler",
		"CREATE FOREIGN DATA WRAPPER myfdw HANDLER fdw_handler VALIDATOR fdw_validator",
		"CREATE FOREIGN DATA WRAPPER myfdw NO HANDLER NO VALIDATOR",
		"CREATE FOREIGN DATA WRAPPER myfdw HANDLER myschema.fdw_handler",
		"CREATE FOREIGN DATA WRAPPER myfdw OPTIONS (debug 'true')",
		"CREATE FOREIGN DATA WRAPPER myfdw HANDLER fdw_handler OPTIONS (debug 'true')",
		"CREATE FOREIGN DATA WRAPPER myfdw NO HANDLER OPTIONS (debug 'true', mode 'fast')",

		// ALTER FOREIGN DATA WRAPPER
		"ALTER FOREIGN DATA WRAPPER myfdw HANDLER new_handler",
		"ALTER FOREIGN DATA WRAPPER myfdw NO HANDLER",
		"ALTER FOREIGN DATA WRAPPER myfdw VALIDATOR new_validator",
		"ALTER FOREIGN DATA WRAPPER myfdw NO VALIDATOR",
		"ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (ADD debug 'true')",
		"ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (SET debug 'false')",
		"ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (DROP debug)",
		"ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (ADD key1 'val1', SET key2 'val2', DROP key3)",
		"ALTER FOREIGN DATA WRAPPER myfdw HANDLER new_handler OPTIONS (ADD debug 'true')",

		// CREATE SERVER
		"CREATE SERVER myserver FOREIGN DATA WRAPPER myfdw",
		"CREATE SERVER IF NOT EXISTS myserver FOREIGN DATA WRAPPER myfdw",
		"CREATE SERVER myserver TYPE 'postgres' FOREIGN DATA WRAPPER myfdw",
		"CREATE SERVER myserver VERSION '15' FOREIGN DATA WRAPPER myfdw",
		"CREATE SERVER myserver TYPE 'postgres' VERSION '15' FOREIGN DATA WRAPPER myfdw",
		"CREATE SERVER myserver TYPE 'postgres' VERSION '15' FOREIGN DATA WRAPPER myfdw OPTIONS (host 'localhost', port '5432')",

		// ALTER SERVER
		"ALTER SERVER myserver VERSION '16'",
		"ALTER SERVER myserver VERSION NULL",
		"ALTER SERVER myserver OPTIONS (SET host 'newhost')",
		"ALTER SERVER myserver VERSION '16' OPTIONS (SET host 'newhost')",

		// CREATE FOREIGN TABLE
		"CREATE FOREIGN TABLE ft (id integer, name text) SERVER myserver",
		"CREATE FOREIGN TABLE ft () SERVER myserver",
		"CREATE FOREIGN TABLE IF NOT EXISTS ft (id integer) SERVER myserver",
		"CREATE FOREIGN TABLE myschema.ft (id integer) SERVER myserver",
		"CREATE FOREIGN TABLE ft (id integer) SERVER myserver OPTIONS (table_name 'remote_table')",
		"CREATE FOREIGN TABLE IF NOT EXISTS ft (id integer) SERVER myserver OPTIONS (table_name 'remote_table')",

		// CREATE USER MAPPING
		"CREATE USER MAPPING FOR CURRENT_USER SERVER myserver",
		"CREATE USER MAPPING FOR CURRENT_USER SERVER myserver OPTIONS (user 'me', password 'secret')",
		"CREATE USER MAPPING IF NOT EXISTS FOR CURRENT_USER SERVER myserver",
		"CREATE USER MAPPING FOR myuser SERVER myserver",
		"CREATE USER MAPPING FOR USER SERVER myserver",

		// ALTER USER MAPPING
		"ALTER USER MAPPING FOR CURRENT_USER SERVER myserver OPTIONS (SET password 'new')",

		// DROP USER MAPPING
		"DROP USER MAPPING FOR CURRENT_USER SERVER myserver",
		"DROP USER MAPPING IF EXISTS FOR CURRENT_USER SERVER myserver",

		// IMPORT FOREIGN SCHEMA
		"IMPORT FOREIGN SCHEMA remote_schema FROM SERVER myserver INTO local_schema",
		"IMPORT FOREIGN SCHEMA remote_schema LIMIT TO (t1, t2) FROM SERVER myserver INTO local_schema",
		"IMPORT FOREIGN SCHEMA remote_schema EXCEPT (t3) FROM SERVER myserver INTO local_schema",
		"IMPORT FOREIGN SCHEMA remote_schema FROM SERVER myserver INTO local_schema OPTIONS (import_default 'true')",

		// CREATE LANGUAGE
		"CREATE LANGUAGE plpgsql",
		"CREATE TRUSTED LANGUAGE plpgsql",
		"CREATE PROCEDURAL LANGUAGE plpgsql",
		"CREATE TRUSTED PROCEDURAL LANGUAGE plpgsql",
		"CREATE LANGUAGE plpgsql HANDLER plpgsql_call_handler",
		"CREATE LANGUAGE plpgsql HANDLER plpgsql_call_handler INLINE plpgsql_inline_handler",
		"CREATE LANGUAGE plpgsql HANDLER plpgsql_call_handler INLINE plpgsql_inline_handler VALIDATOR plpgsql_validator",
		"CREATE LANGUAGE plpgsql HANDLER plpgsql_call_handler VALIDATOR plpgsql_validator",
		"CREATE LANGUAGE plpgsql HANDLER plpgsql_call_handler NO VALIDATOR",
		"CREATE OR REPLACE LANGUAGE plpgsql",
		"CREATE OR REPLACE TRUSTED LANGUAGE plpgsql HANDLER plpgsql_call_handler",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareTrigger runs trigger comparison tests for batch 21.
func TestCompareTrigger(t *testing.T) {
	tests := []string{
		// CREATE TRIGGER - basic
		"CREATE TRIGGER trig BEFORE INSERT ON t EXECUTE FUNCTION func()",
		"CREATE TRIGGER trig AFTER UPDATE ON t EXECUTE FUNCTION func()",
		"CREATE TRIGGER trig AFTER INSERT OR UPDATE ON t EXECUTE FUNCTION func()",
		"CREATE TRIGGER trig BEFORE DELETE ON t FOR EACH ROW EXECUTE FUNCTION func()",

		// UPDATE OF columns
		"CREATE TRIGGER trig AFTER UPDATE OF col1, col2 ON t EXECUTE FUNCTION func()",

		// INSTEAD OF
		"CREATE TRIGGER trig INSTEAD OF INSERT ON v EXECUTE FUNCTION func()",

		// WHEN clause
		"CREATE TRIGGER trig BEFORE INSERT ON t FOR EACH ROW WHEN (NEW.active) EXECUTE FUNCTION func()",

		// OR REPLACE
		"CREATE OR REPLACE TRIGGER trig BEFORE INSERT ON t EXECUTE FUNCTION func()",

		// EXECUTE PROCEDURE (alias)
		"CREATE TRIGGER trig BEFORE INSERT ON t EXECUTE PROCEDURE func()",

		// Constraint trigger
		"CREATE CONSTRAINT TRIGGER trig AFTER INSERT ON t DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION func()",

		// FOR EACH STATEMENT
		"CREATE TRIGGER trig AFTER INSERT ON t FOR EACH STATEMENT EXECUTE FUNCTION func()",

		// CREATE EVENT TRIGGER
		"CREATE EVENT TRIGGER etrig ON ddl_command_start EXECUTE FUNCTION func()",
		"CREATE EVENT TRIGGER etrig ON ddl_command_start WHEN TAG IN ('CREATE TABLE', 'DROP TABLE') EXECUTE FUNCTION func()",

		// ALTER EVENT TRIGGER
		"ALTER EVENT TRIGGER etrig ENABLE",
		"ALTER EVENT TRIGGER etrig DISABLE",
		"ALTER EVENT TRIGGER etrig ENABLE REPLICA",
		"ALTER EVENT TRIGGER etrig ENABLE ALWAYS",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareSchema runs CREATE SCHEMA comparison tests for batch 19.
func TestCompareSchema(t *testing.T) {
	tests := []string{
		// Basic CREATE SCHEMA
		"CREATE SCHEMA myschema",
		"CREATE SCHEMA IF NOT EXISTS myschema",

		// With AUTHORIZATION
		"CREATE SCHEMA AUTHORIZATION myrole",
		"CREATE SCHEMA myschema AUTHORIZATION myrole",
		"CREATE SCHEMA IF NOT EXISTS AUTHORIZATION myrole",
		"CREATE SCHEMA IF NOT EXISTS myschema AUTHORIZATION myrole",

		// With schema elements
		"CREATE SCHEMA myschema CREATE TABLE t(id int)",
		"CREATE SCHEMA myschema CREATE TABLE t(id int) CREATE VIEW v AS SELECT 1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareComment runs COMMENT ON comparison tests for batch 19.
func TestCompareComment(t *testing.T) {
	tests := []string{
		// COMMENT ON object_type_any_name
		"COMMENT ON TABLE t IS 'my comment'",
		"COMMENT ON TABLE t IS NULL",
		"COMMENT ON VIEW v IS 'view comment'",
		"COMMENT ON INDEX idx IS 'index comment'",
		"COMMENT ON SEQUENCE s IS 'seq comment'",

		// COMMENT ON COLUMN
		"COMMENT ON COLUMN t.c IS 'column comment'",

		// COMMENT ON object_type_name
		"COMMENT ON SCHEMA myschema IS 'schema comment'",
		"COMMENT ON DATABASE mydb IS 'db comment'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareRename runs ALTER ... RENAME comparison tests for batch 19.
func TestCompareRename(t *testing.T) {
	tests := []string{
		// ALTER SCHEMA RENAME
		"ALTER SCHEMA old_schema RENAME TO new_schema",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestComparePublication runs publication comparison tests for batch 33.
func TestComparePublication(t *testing.T) {
	tests := []string{
		"CREATE PUBLICATION pub",
		"CREATE PUBLICATION pub FOR ALL TABLES",
		"CREATE PUBLICATION pub FOR TABLE t1",
		"CREATE PUBLICATION pub FOR TABLE t1, TABLE t2",
		"CREATE PUBLICATION pub FOR TABLE t1 WHERE (col > 0)",
		"CREATE PUBLICATION pub FOR TABLES IN SCHEMA public",
		"CREATE PUBLICATION pub FOR TABLES IN SCHEMA CURRENT_SCHEMA",
		"CREATE PUBLICATION pub FOR ALL TABLES WITH (publish = 'insert')",
		"CREATE PUBLICATION pub WITH (publish = 'insert')",
		"ALTER PUBLICATION pub SET (publish = 'insert,update')",
		"ALTER PUBLICATION pub ADD TABLE t3",
		"ALTER PUBLICATION pub DROP TABLE t1",
		"ALTER PUBLICATION pub SET TABLE t1, TABLE t2",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareSubscription runs subscription comparison tests for batch 33.
func TestCompareSubscription(t *testing.T) {
	tests := []string{
		"CREATE SUBSCRIPTION sub CONNECTION 'host=localhost' PUBLICATION pub1",
		"CREATE SUBSCRIPTION sub CONNECTION 'host=localhost' PUBLICATION pub1, pub2",
		"CREATE SUBSCRIPTION sub CONNECTION 'connstr' PUBLICATION pub WITH (slot_name = 'none')",
		"ALTER SUBSCRIPTION sub CONNECTION 'new_conn'",
		"ALTER SUBSCRIPTION sub SET PUBLICATION pub2",
		"ALTER SUBSCRIPTION sub REFRESH PUBLICATION",
		"ALTER SUBSCRIPTION sub ENABLE",
		"ALTER SUBSCRIPTION sub DISABLE",
		"ALTER SUBSCRIPTION sub SET (slot_name = 'myslot')",
		"ALTER SUBSCRIPTION sub ADD PUBLICATION pub2",
		"ALTER SUBSCRIPTION sub DROP PUBLICATION pub1",
		"DROP SUBSCRIPTION sub",
		"DROP SUBSCRIPTION IF EXISTS sub",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareRule runs rule comparison tests for batch 33.
func TestCompareRule(t *testing.T) {
	tests := []string{
		"CREATE RULE r AS ON INSERT TO t DO NOTHING",
		"CREATE RULE r AS ON INSERT TO t DO INSTEAD NOTHING",
		"CREATE RULE r AS ON INSERT TO t DO ALSO INSERT INTO log DEFAULT VALUES",
		"CREATE RULE r AS ON DELETE TO t DO (INSERT INTO log DEFAULT VALUES)",
		"CREATE OR REPLACE RULE r AS ON UPDATE TO t WHERE OLD.id > 0 DO NOTHING",
		"CREATE RULE r AS ON SELECT TO v DO INSTEAD SELECT * FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareGrant runs GRANT/REVOKE comparison tests.
func TestCompareGrant(t *testing.T) {
	tests := []string{
		// GRANT on tables
		"GRANT SELECT ON TABLE t TO alice",
		"GRANT ALL ON t TO alice",
		"GRANT ALL PRIVILEGES ON TABLE t TO alice",
		"GRANT SELECT, INSERT, UPDATE ON t TO alice, bob",
		"GRANT SELECT ON t TO CURRENT_USER",
		"GRANT SELECT ON t TO SESSION_USER",
		"GRANT SELECT ON t TO alice WITH GRANT OPTION",
		// GRANT on multiple object types
		"GRANT USAGE ON SEQUENCE s TO alice",
		"GRANT CREATE ON SCHEMA myschema TO alice",
		"GRANT CREATE ON DATABASE mydb TO alice",
		"GRANT EXECUTE ON FUNCTION myfunc() TO alice",
		"GRANT EXECUTE ON FUNCTION myfunc(integer, text) TO alice",
		"GRANT ALL ON ALL TABLES IN SCHEMA public TO alice",
		"GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO alice",
		"GRANT ALL ON ALL FUNCTIONS IN SCHEMA public TO alice",
		// REVOKE
		"REVOKE SELECT ON t FROM alice",
		"REVOKE ALL ON t FROM alice CASCADE",
		"REVOKE ALL PRIVILEGES ON TABLE t FROM alice RESTRICT",
		"REVOKE SELECT ON t FROM alice, bob",
		"REVOKE GRANT OPTION FOR SELECT ON t FROM alice",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareRole runs CREATE/ALTER/DROP ROLE comparison tests.
func TestCompareRole(t *testing.T) {
	tests := []string{
		// CREATE ROLE
		"CREATE ROLE myrole",
		"CREATE ROLE myrole WITH LOGIN",
		"CREATE ROLE myrole WITH SUPERUSER",
		"CREATE ROLE myrole WITH PASSWORD 'secret'",
		"CREATE ROLE myrole WITH PASSWORD NULL",
		"CREATE ROLE myrole CONNECTION LIMIT 10",
		"CREATE ROLE myrole VALID UNTIL '2025-01-01'",
		"CREATE USER myuser",
		"CREATE USER myuser WITH PASSWORD 'secret'",
		"CREATE GROUP mygroup",
		// ALTER ROLE
		"ALTER ROLE myrole SUPERUSER",
		"ALTER ROLE myrole NOSUPERUSER",
		"ALTER ROLE myrole CREATEDB",
		"ALTER ROLE myrole NOCREATEDB",
		"ALTER ROLE myrole LOGIN",
		"ALTER ROLE myrole NOLOGIN",
		"ALTER ROLE myrole SET search_path TO public",
		"ALTER ROLE myrole RESET search_path",
		// DROP ROLE
		"DROP ROLE myrole",
		"DROP ROLE IF EXISTS myrole",
		"DROP ROLE myrole, otherrole",
		"DROP USER myuser",
		"DROP USER IF EXISTS myuser",
		"DROP GROUP mygroup",
		// GRANT ROLE
		"GRANT myrole TO alice",
		"GRANT myrole TO alice WITH ADMIN OPTION",
		// REVOKE ROLE
		"REVOKE myrole FROM alice",
		"REVOKE ADMIN OPTION FOR myrole FROM alice",
		// ALTER GROUP
		"ALTER GROUP mygroup ADD USER alice",
		"ALTER GROUP mygroup DROP USER alice",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestComparePolicy runs CREATE/ALTER POLICY comparison tests.
func TestComparePolicy(t *testing.T) {
	tests := []string{
		"CREATE POLICY p ON t USING (true)",
		"CREATE POLICY p ON t FOR SELECT TO PUBLIC USING (true)",
		"CREATE POLICY p ON t FOR INSERT WITH CHECK (true)",
		"CREATE POLICY p ON t FOR ALL TO alice, bob USING (true) WITH CHECK (true)",
		"CREATE POLICY p ON t AS PERMISSIVE FOR ALL TO PUBLIC USING (true)",
		"CREATE POLICY p ON t AS RESTRICTIVE FOR SELECT USING (id > 0)",
		"CREATE POLICY p ON t FOR DELETE USING (owner = CURRENT_USER)",
		"ALTER POLICY p ON t USING (true)",
		"ALTER POLICY p ON t TO alice, bob",
		"ALTER POLICY p ON t USING (id > 0) WITH CHECK (id > 0)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareDefine runs comparison tests for batch 20 (define).
func TestCompareDefine(t *testing.T) {
	tests := []string{
		// CREATE AGGREGATE - new style
		"CREATE AGGREGATE myagg(integer) (SFUNC = int4_avg_accum, STYPE = _int8)",
		"CREATE AGGREGATE myagg(integer, integer) (SFUNC = f, STYPE = int8)",
		"CREATE AGGREGATE myagg(*) (SFUNC = f, STYPE = int8)",
		// CREATE AGGREGATE - old style
		"CREATE AGGREGATE myagg(basetype = integer, sfunc = f, stype = int8)",
		// CREATE OR REPLACE AGGREGATE
		"CREATE OR REPLACE AGGREGATE myagg(integer) (SFUNC = f, STYPE = int8)",
		// CREATE OPERATOR
		"CREATE OPERATOR === (leftarg = text, rightarg = text, procedure = texteq)",
		"CREATE OPERATOR + (leftarg = integer, rightarg = integer, procedure = int4pl)",
		"CREATE OPERATOR myschema.+ (leftarg = integer, rightarg = integer, procedure = int4pl)",
		// CREATE TYPE - shell
		"CREATE TYPE mytype",
		// CREATE TYPE - base
		"CREATE TYPE mytype (INPUT = mytype_in, OUTPUT = mytype_out)",
		"CREATE TYPE myschema.mytype (INPUT = mytype_in, OUTPUT = mytype_out)",
		"CREATE TYPE mytype (PASSEDBYVALUE)",
		"CREATE TYPE mytype (INTERNALLENGTH = 16)",
		// CREATE TYPE - composite
		"CREATE TYPE complex AS (r float8, i float8)",
		"CREATE TYPE empty_type AS ()",
		"CREATE TYPE myschema.complex AS (r float8, i float8)",
		// CREATE TYPE - enum
		"CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy')",
		"CREATE TYPE mood AS ENUM ()",
		// CREATE TYPE - range
		"CREATE TYPE float8_range AS RANGE (SUBTYPE = float8)",
		"CREATE TYPE r AS RANGE (SUBTYPE = float8, SUBTYPE_OPCLASS = float8_ops)",
		// CREATE TEXT SEARCH
		"CREATE TEXT SEARCH PARSER myparser (START = prsd_start, GETTOKEN = prsd_nexttoken, END = prsd_end, LEXTYPES = prsd_lextype)",
		"CREATE TEXT SEARCH DICTIONARY mydict (TEMPLATE = simple, STOPWORDS = english)",
		"CREATE TEXT SEARCH DICTIONARY d (TEMPLATE = simple, STOPWORDS = 'english')",
		"CREATE TEXT SEARCH TEMPLATE mytempl (INIT = dsimple_init, LEXIZE = dsimple_lexize)",
		"CREATE TEXT SEARCH CONFIGURATION myconfig (PARSER = default)",
		// CREATE COLLATION
		"CREATE COLLATION mycoll (LOCALE = 'en_US.utf8')",
		"CREATE COLLATION IF NOT EXISTS mycoll (LOCALE = 'en_US.utf8')",
		`CREATE COLLATION mycoll FROM "C"`,
		`CREATE COLLATION IF NOT EXISTS mycoll FROM "C"`,
		// DefElem variations
		"CREATE TYPE mytype (INTERNALLENGTH = 16)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareOpClass runs comparison tests for batch 20 (opclass).
func TestCompareOpClass(t *testing.T) {
	tests := []string{
		"CREATE OPERATOR CLASS myclass DEFAULT FOR TYPE int4 USING btree AS OPERATOR 1 <, FUNCTION 1 btint4cmp(int4, int4)",
		"CREATE OPERATOR CLASS myclass FOR TYPE int4 USING btree AS OPERATOR 1 <",
		"CREATE OPERATOR FAMILY myfam USING btree",
		"ALTER OPERATOR FAMILY myfam USING btree ADD OPERATOR 1 <(int4, int4)",
		"ALTER OPERATOR FAMILY myfam USING btree DROP OPERATOR 1 (int4, int4)",
		"DROP OPERATOR CLASS myclass USING btree",
		"DROP OPERATOR CLASS IF EXISTS myclass USING btree CASCADE",
		"DROP OPERATOR FAMILY myfam USING btree",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareStats runs comparison tests for batch 20 (stats).
func TestCompareStats(t *testing.T) {
	tests := []string{
		"CREATE STATISTICS mystat ON col1, col2 FROM t",
		"CREATE STATISTICS IF NOT EXISTS mystat (dependencies) ON col1, col2 FROM t",
		"CREATE STATISTICS mystat (ndistinct, dependencies) ON col1, col2 FROM t",
		"ALTER STATISTICS mystat SET STATISTICS 100",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareAlterFunc runs comparison tests for ALTER FUNCTION/PROCEDURE/ROUTINE statements.
func TestCompareAlterFunc(t *testing.T) {
	tests := []string{
		// ALTER FUNCTION with alterfunc_opt_list
		"ALTER FUNCTION myfunc(integer) IMMUTABLE",
		"ALTER FUNCTION myfunc(integer) STABLE",
		"ALTER FUNCTION myfunc(integer) VOLATILE",
		"ALTER FUNCTION myfunc(integer) STRICT",
		"ALTER FUNCTION myfunc(integer) CALLED ON NULL INPUT",
		"ALTER FUNCTION myfunc(integer) SECURITY DEFINER",
		"ALTER FUNCTION myfunc(integer) SECURITY INVOKER",
		"ALTER FUNCTION myfunc(integer) COST 100",
		"ALTER FUNCTION myfunc(integer) ROWS 1000",
		"ALTER FUNCTION myfunc(integer) IMMUTABLE STRICT",
		"ALTER FUNCTION myfunc(integer) STABLE SECURITY DEFINER COST 50",
		"ALTER FUNCTION myfunc(integer) IMMUTABLE RESTRICT",
		"ALTER FUNCTION myfunc IMMUTABLE",
		"ALTER FUNCTION myfunc(integer, text) IMMUTABLE",
		"ALTER FUNCTION myschema.myfunc(integer) IMMUTABLE",
		"ALTER PROCEDURE myproc(integer) SECURITY INVOKER",
		"ALTER PROCEDURE myproc(integer) SECURITY DEFINER",
		"ALTER ROUTINE myroutine(integer) IMMUTABLE",
		"ALTER FUNCTION myfunc(integer) RENAME TO newfunc",
		"ALTER PROCEDURE myproc(integer) RENAME TO newproc",
		"ALTER FUNCTION myfunc(integer) OWNER TO newowner",
		"ALTER PROCEDURE myproc(integer) OWNER TO newowner",
		"ALTER FUNCTION myfunc(integer) SET SCHEMA newschema",
		"ALTER PROCEDURE myproc(integer) SET SCHEMA newschema",
		"ALTER FUNCTION myfunc(integer) DEPENDS ON EXTENSION myext",
		"ALTER FUNCTION myfunc(integer) NO DEPENDS ON EXTENSION myext",
		"ALTER PROCEDURE myproc(integer) DEPENDS ON EXTENSION myext",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}

// TestCompareAlterGeneric runs comparison tests for various ALTER statements (batch 15).
func TestCompareAlterGeneric(t *testing.T) {
	tests := []string{
		// ALTER TYPE - enum
		"ALTER TYPE myenum ADD VALUE 'newval'",
		"ALTER TYPE myenum ADD VALUE IF NOT EXISTS 'newval'",
		"ALTER TYPE myenum ADD VALUE 'newval' BEFORE 'existval'",
		"ALTER TYPE myenum ADD VALUE 'newval' AFTER 'existval'",
		"ALTER TYPE myenum RENAME VALUE 'old' TO 'new'",

		// ALTER TYPE - composite
		"ALTER TYPE mytype ADD ATTRIBUTE col1 integer",
		"ALTER TYPE mytype ADD ATTRIBUTE col1 integer CASCADE",
		"ALTER TYPE mytype DROP ATTRIBUTE col1",
		"ALTER TYPE mytype DROP ATTRIBUTE IF EXISTS col1",
		"ALTER TYPE mytype DROP ATTRIBUTE col1 CASCADE",
		"ALTER TYPE mytype ALTER ATTRIBUTE col1 TYPE integer",
		"ALTER TYPE mytype ALTER ATTRIBUTE col1 SET DATA TYPE text",

		// ALTER TYPE - rename / owner / schema / set
		"ALTER TYPE mytype RENAME TO newtype",
		"ALTER TYPE mytype RENAME ATTRIBUTE col1 TO col2",
		"ALTER TYPE mytype OWNER TO newowner",
		"ALTER TYPE mytype SET SCHEMA newschema",
		"ALTER TYPE mytype SET (RECEIVE = myrecv)",

		// ALTER DOMAIN
		"ALTER DOMAIN mydom SET DEFAULT 42",
		"ALTER DOMAIN mydom DROP DEFAULT",
		"ALTER DOMAIN mydom SET NOT NULL",
		"ALTER DOMAIN mydom DROP NOT NULL",
		"ALTER DOMAIN mydom ADD CHECK (VALUE > 0)",
		"ALTER DOMAIN mydom ADD CONSTRAINT mycheck CHECK (VALUE > 0)",
		"ALTER DOMAIN mydom DROP CONSTRAINT mycheck",
		"ALTER DOMAIN mydom DROP CONSTRAINT IF EXISTS mycheck",
		"ALTER DOMAIN mydom DROP CONSTRAINT mycheck CASCADE",
		"ALTER DOMAIN mydom VALIDATE CONSTRAINT mycheck",
		"ALTER DOMAIN mydom OWNER TO newowner",
		"ALTER DOMAIN mydom RENAME TO newdom",
		"ALTER DOMAIN mydom RENAME CONSTRAINT oldcon TO newcon",
		"ALTER DOMAIN mydom SET SCHEMA newschema",

		// ALTER SCHEMA
		"ALTER SCHEMA myschema OWNER TO newowner",
		"ALTER SCHEMA myschema RENAME TO newschema",

		// ALTER COLLATION
		"ALTER COLLATION mycoll REFRESH VERSION",
		"ALTER COLLATION mycoll RENAME TO newcoll",
		"ALTER COLLATION mycoll OWNER TO newowner",
		"ALTER COLLATION mycoll SET SCHEMA newschema",

		// ALTER CONVERSION
		"ALTER CONVERSION myconv RENAME TO newconv",
		"ALTER CONVERSION myconv OWNER TO newowner",
		"ALTER CONVERSION myconv SET SCHEMA newschema",

		// ALTER AGGREGATE
		"ALTER AGGREGATE myagg(integer) RENAME TO newagg",
		"ALTER AGGREGATE myagg(integer) OWNER TO newowner",
		"ALTER AGGREGATE myagg(integer) SET SCHEMA newschema",

		// ALTER OPERATOR (plain)
		"ALTER OPERATOR +(integer, integer) SET (RESTRICT = myfunc)",
		"ALTER OPERATOR +(integer, integer) OWNER TO newowner",
		"ALTER OPERATOR +(integer, integer) SET SCHEMA newschema",

		// ALTER OPERATOR CLASS / FAMILY
		"ALTER OPERATOR CLASS myclass USING btree SET SCHEMA newschema",
		"ALTER OPERATOR CLASS myclass USING btree OWNER TO newowner",
		"ALTER OPERATOR CLASS myclass USING btree RENAME TO newclass",
		"ALTER OPERATOR FAMILY myfam USING btree SET SCHEMA newschema",
		"ALTER OPERATOR FAMILY myfam USING btree OWNER TO newowner",
		"ALTER OPERATOR FAMILY myfam USING btree RENAME TO newfam",

		// ALTER TEXT SEARCH DICTIONARY
		"ALTER TEXT SEARCH DICTIONARY mydict (STOPWORDS = 'english')",
		"ALTER TEXT SEARCH DICTIONARY mydict RENAME TO newdict",
		"ALTER TEXT SEARCH DICTIONARY mydict OWNER TO newowner",
		"ALTER TEXT SEARCH DICTIONARY mydict SET SCHEMA newschema",

		// ALTER TEXT SEARCH CONFIGURATION
		"ALTER TEXT SEARCH CONFIGURATION myconfig ADD MAPPING FOR word WITH simple",
		"ALTER TEXT SEARCH CONFIGURATION myconfig ALTER MAPPING FOR word WITH simple",
		"ALTER TEXT SEARCH CONFIGURATION myconfig ALTER MAPPING REPLACE simple WITH english",
		"ALTER TEXT SEARCH CONFIGURATION myconfig ALTER MAPPING FOR word REPLACE simple WITH english",
		"ALTER TEXT SEARCH CONFIGURATION myconfig DROP MAPPING FOR word",
		"ALTER TEXT SEARCH CONFIGURATION myconfig DROP MAPPING IF EXISTS FOR word",
		"ALTER TEXT SEARCH CONFIGURATION myconfig RENAME TO newconfig",
		"ALTER TEXT SEARCH CONFIGURATION myconfig OWNER TO newowner",
		"ALTER TEXT SEARCH CONFIGURATION myconfig SET SCHEMA newschema",

		// ALTER TEXT SEARCH PARSER / TEMPLATE
		"ALTER TEXT SEARCH PARSER myparser RENAME TO newparser",
		"ALTER TEXT SEARCH PARSER myparser SET SCHEMA newschema",
		"ALTER TEXT SEARCH TEMPLATE mytemplate RENAME TO newtemplate",
		"ALTER TEXT SEARCH TEMPLATE mytemplate SET SCHEMA newschema",

		// ALTER LANGUAGE
		"ALTER LANGUAGE plpgsql RENAME TO newlang",
		"ALTER LANGUAGE plpgsql OWNER TO newowner",

		// ALTER LARGE OBJECT
		"ALTER LARGE OBJECT 12345 OWNER TO newowner",

		// ALTER EVENT TRIGGER
		"ALTER EVENT TRIGGER mytrig OWNER TO newowner",
		"ALTER EVENT TRIGGER mytrig RENAME TO newtrig",
		"ALTER EVENT TRIGGER mytrig ENABLE",
		"ALTER EVENT TRIGGER mytrig ENABLE REPLICA",
		"ALTER EVENT TRIGGER mytrig ENABLE ALWAYS",
		"ALTER EVENT TRIGGER mytrig DISABLE",

		// ALTER TABLESPACE
		"ALTER TABLESPACE mytbs OWNER TO newowner",
		"ALTER TABLESPACE mytbs RENAME TO newtbs",

		// ALTER TRIGGER ... DEPENDS ON EXTENSION
		"ALTER TRIGGER mytrig ON mytable DEPENDS ON EXTENSION myext",
		"ALTER TRIGGER mytrig ON mytable NO DEPENDS ON EXTENSION myext",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			CompareWithYacc(t, sql)
		})
	}
}
