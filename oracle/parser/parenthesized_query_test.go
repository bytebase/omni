package parser

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParenthesizedQueries checks queries whose operands are parenthesized
// subqueries, and the EXCEPT and ALL set operators, against Oracle 23ai:
// every accepted form parses there, and every rejected one fails with the
// error noted beside it. Rows ref_225 to ref_236, ref_251 to ref_252,
// ref_257, and ref_260 run them against the engine.
func TestParenthesizedQueries(t *testing.T) {
	for _, sql := range []string{
		// A statement, a set operand, and what may follow the parentheses.
		"(SELECT 1 FROM dual)",
		"((SELECT 1 FROM dual))",
		"(SELECT 1 a FROM dual) ORDER BY 1",
		"(SELECT 1 a FROM dual) UNION ALL (SELECT 2 FROM dual) ORDER BY 1",
		"SELECT 1 FROM dual UNION ALL (SELECT 2 FROM dual)",
		"SELECT 1 FROM dual UNION (SELECT 2 FROM dual) ORDER BY 1",
		"(SELECT 1 FROM dual UNION SELECT 2 FROM dual) INTERSECT SELECT 1 FROM dual",
		"(SELECT 1 FROM dual) UNION ALL (SELECT 2 FROM dual) UNION ALL (SELECT 3 FROM dual)",
		"(SELECT 1 FROM dual) UNION ALL ((SELECT 2 FROM dual) MINUS (SELECT 3 FROM dual))",
		"(SELECT 1 FROM dual) EXCEPT (SELECT 2 FROM dual)",
		"(SELECT a FROM t) FOR UPDATE",
		"(SELECT 1 FROM dual) FETCH FIRST 1 ROWS ONLY",
		"(SELECT 1 FROM dual) OFFSET 1 ROWS",
		"(SELECT 1 a FROM dual FETCH FIRST 1 ROWS ONLY) UNION (SELECT 2 FROM dual)",
		"WITH w AS (SELECT 1 a FROM dual) (SELECT a FROM w)",
		"WITH w AS ((SELECT 1 a FROM dual) UNION (SELECT 2 FROM dual)) SELECT a FROM w",
		// EXCEPT and ALL on every operator (23ai).
		"SELECT 1 FROM dual EXCEPT SELECT 2 FROM dual",
		"SELECT 1 FROM dual EXCEPT ALL SELECT 2 FROM dual",
		"SELECT 1 FROM dual MINUS ALL SELECT 2 FROM dual",
		"SELECT 1 FROM dual INTERSECT ALL SELECT 2 FROM dual",
		// In FROM, the parentheses hold a query or a table reference.
		"SELECT * FROM ((SELECT 1 a FROM dual) UNION (SELECT 2 FROM dual)) x",
		"SELECT * FROM ((SELECT 1 a FROM dual) UNION (SELECT 2 FROM dual) ORDER BY 1)",
		"SELECT * FROM ((SELECT 1 a FROM dual ORDER BY 1))",
		"SELECT * FROM ((SELECT 1 a FROM dual) x JOIN (SELECT 1 a FROM dual) y ON x.a = y.a)",
		"SELECT * FROM (WITH w AS (SELECT 1 a FROM dual) SELECT a FROM w)",
		// Redundant parentheses around an inline view, its alias after them.
		"SELECT * FROM ((SELECT 1 a FROM dual)) x",
		"SELECT * FROM (((SELECT 1 a FROM dual))) x",
		"SELECT x.a FROM (((SELECT 1 a FROM dual) UNION (SELECT 2 FROM dual))) x",
		// In expressions, likewise a query or a parenthesized expression.
		"SELECT 1 FROM dual WHERE 1 IN ((SELECT 1 FROM dual) UNION (SELECT 2 FROM dual))",
		"SELECT 1 FROM dual WHERE (1, 2) NOT IN ((SELECT 1, 2 FROM dual) UNION ALL (SELECT 3, 4 FROM dual))",
		"SELECT 1 FROM dual WHERE EXISTS ((SELECT 1 FROM dual) UNION (SELECT 2 FROM dual))",
		"SELECT ((SELECT 1 FROM dual) UNION (SELECT 1 FROM dual)) FROM dual",
		"SELECT 1 FROM dual WHERE ((SELECT 1 FROM dual) + 1) = 2",
		"SELECT 1 FROM dual WHERE 1 IN ((SELECT 1 FROM dual), 2)",
		"SELECT 1 FROM dual WHERE 1 = ANY ((SELECT 1 FROM dual) UNION (SELECT 2 FROM dual))",
		// Where only a query may stand.
		"INSERT INTO t (a) (SELECT 1 FROM dual)",
		"INSERT INTO t (SELECT 1, 'x' FROM dual)",
		"INSERT INTO t (a) ((SELECT 1 FROM dual) UNION (SELECT 2 FROM dual))",
		"INSERT INTO t (a) (SELECT 1 FROM dual) UNION (SELECT 2 FROM dual)",
		"CREATE TABLE t AS (SELECT 1 a FROM dual)",
		"CREATE TABLE t AS (SELECT 1 a FROM dual) UNION (SELECT 2 FROM dual) ORDER BY 1",
		"CREATE VIEW v AS (SELECT 1 a FROM dual) UNION ALL (SELECT 2 FROM dual)",
		"MERGE INTO t USING ((SELECT 1 a FROM dual) UNION (SELECT 2 FROM dual)) s ON (t.a = s.a) WHEN MATCHED THEN UPDATE SET t.b = 'x'",
		"DELETE FROM t WHERE a IN ((SELECT 1 FROM dual) MINUS (SELECT 2 FROM dual))",
		"DECLARE CURSOR c IS (SELECT 1 a FROM dual) UNION (SELECT 2 FROM dual); BEGIN NULL; END;",
		"DECLARE CURSOR c IS ((SELECT 1 a FROM dual)); BEGIN NULL; END;",
		"BEGIN FOR r IN ((SELECT 1 a FROM dual) UNION ALL (SELECT 2 FROM dual)) LOOP NULL; END LOOP; END;",
		"BEGIN FOR r IN (SELECT 1 a FROM dual UNION ALL (SELECT 2 FROM dual)) LOOP NULL; END LOOP; END;",
		"BEGIN FOR r IN ((SELECT 1 a FROM dual) UNION (SELECT 2 FROM dual) ORDER BY 1) LOOP NULL; END LOOP; END;",
		// A multi-column SET takes its arity from the first operand, and a
		// * leaves it to the engine.
		"UPDATE t SET (a, b) = ((SELECT 1, 2 FROM dual) UNION (SELECT 1, 2 FROM dual))",
		"UPDATE t SET (a, b) = (SELECT * FROM t2)",
		"UPDATE t SET (a, b) = (SELECT x.* FROM t2 x)",
	} {
		ParseAndCheck(t, sql)
	}

	for _, sql := range []string{
		// An ORDER BY ends no parenthesized subquery (ORA-00907).
		"(SELECT 1 a FROM dual ORDER BY 1)",
		"(SELECT 1 a FROM dual ORDER BY 1) UNION (SELECT 2 FROM dual)",
		"(SELECT 1 a FROM dual) UNION (SELECT 2 FROM dual ORDER BY 1)",
		"(SELECT 1 FROM dual UNION SELECT 2 FROM dual ORDER BY 1)",
		"((SELECT 1 a FROM dual) ORDER BY 1)",
		"CREATE TABLE t AS (SELECT 1 a FROM dual ORDER BY 1)",
		"INSERT INTO t (a) (SELECT 1 FROM dual ORDER BY 1)",
		"CREATE VIEW v AS (SELECT 1 a FROM dual ORDER BY 1)",
		"DECLARE CURSOR c IS (SELECT 1 a FROM dual ORDER BY 1); BEGIN NULL; END;",
		// Nor does a WITH clause start one (ORA-32034).
		"(WITH w AS (SELECT 1 a FROM dual) SELECT a FROM w)",
		"(WITH w AS (SELECT 1 a FROM dual) SELECT a FROM w) UNION SELECT 1 FROM dual",
		"SELECT 1 FROM dual UNION (WITH w AS (SELECT 1 a FROM dual) SELECT a FROM w)",
		"SELECT 1 FROM dual UNION WITH w AS (SELECT 2 x FROM dual) SELECT x FROM w",
		"SELECT 1 FROM dual MINUS ALL WITH w AS (SELECT 2 x FROM dual) SELECT x FROM w",
		"SELECT 1 FROM dual EXCEPT WITH w AS (SELECT 2 x FROM dual) SELECT x FROM w",
		"INSERT INTO t (a) (WITH w AS (SELECT 1 a FROM dual) SELECT a FROM w)",
		// No set operator follows ORDER BY, FOR UPDATE, OFFSET, or FETCH
		// (ORA-03048).
		"(SELECT 1 FROM dual) ORDER BY 1 UNION SELECT 2 FROM dual",
		"SELECT 1 FROM dual ORDER BY 1 UNION SELECT 2 FROM dual",
		"SELECT 1 FROM dual FETCH FIRST 1 ROWS ONLY UNION SELECT 2 FROM dual",
		"(SELECT 1 FROM dual) OFFSET 1 ROWS UNION SELECT 2 FROM dual",
		"SELECT a FROM t FOR UPDATE UNION SELECT 2 FROM dual",
		// The statement ends at the parentheses (ORA-03049), EXCEPT is a set
		// operator rather than an alias (ORA-00928), and the parentheses
		// hold a query.
		"( SELECT 1 FROM dual ) x",
		"SELECT 1 except FROM dual",
		"(1)",
		// SQL reserves EXCEPT as it does MINUS: no operand, column, table,
		// or alias (ORA-00936, ORA-00903, ORA-03050), in static SQL too.
		"SELECT 1 + except FROM dual",
		"SELECT 1 FROM dual WHERE except = 1",
		"SELECT * FROM except",
		"SELECT 1 FROM sys.except",
		"SELECT except.a FROM t except",
		"CREATE TABLE t (except NUMBER)",
		"DECLARE n NUMBER; BEGIN SELECT 1 INTO n FROM dual WHERE except = 1; END;",
		// A query follows a set operator (ORA-00928).
		"(SELECT 1 FROM dual EXCEPT)",
		"(SELECT 1 FROM dual MINUS ALL)",
		"SELECT 1 FROM dual UNION",
		"SELECT 1 FROM dual UNION (1)",
		// No operand selects * without FROM (ORA-54058).
		"(SELECT *) ORDER BY 1",
		"(SELECT *) UNION SELECT 1 FROM dual",
		"((SELECT *)) UNION SELECT 1 FROM dual",
		"SELECT 1 FROM dual UNION (SELECT *)",
		// A multi-column SET gets as many values as columns (ORA-00947).
		"UPDATE t SET (a, b) = ((SELECT 1 FROM dual) UNION (SELECT 1 FROM dual))",
	} {
		ParseShouldFail(t, sql)
	}
}

// TestParenthesizedQueryLimits checks against Oracle 23ai where a
// parenthesized operand takes an ORDER BY, and where a query block selects
// * without FROM.
func TestParenthesizedQueryLimits(t *testing.T) {
	// An ORDER BY a row-limiting clause follows belongs to it, so the
	// parenthesized operand takes it; a bare one, or one before FOR UPDATE,
	// fails (ORA-00907).
	for _, sql := range []string{
		"(SELECT deptno FROM emp ORDER BY deptno FETCH FIRST 10 ROWS ONLY) UNION ALL SELECT deptno FROM emp",
		"(SELECT deptno FROM emp ORDER BY deptno OFFSET 1 ROWS) UNION ALL SELECT deptno FROM emp",
		"(SELECT 1 a FROM dual UNION SELECT 2 FROM dual ORDER BY 1 FETCH FIRST 1 ROWS ONLY) UNION ALL SELECT 3 FROM dual",
		"SELECT * FROM dual WHERE 1 IN (SELECT 1 FROM dual ORDER BY 1 FETCH FIRST 1 ROWS ONLY)",
	} {
		ParseAndCheck(t, sql)
	}
	for _, sql := range []string{
		"(SELECT deptno FROM emp ORDER BY deptno) UNION ALL SELECT deptno FROM emp",
		"(SELECT 1 a FROM dual UNION SELECT 2 FROM dual ORDER BY 1) UNION ALL SELECT 3 FROM dual",
		"(SELECT 1 a FROM dual ORDER BY 1 FOR UPDATE) UNION ALL SELECT 3 FROM dual",
	} {
		ParseShouldFail(t, sql)
	}
	// Nor does a FOR UPDATE end one, operand or not (ORA-00907 at FOR); a
	// parenthesized query takes one after it.
	for _, sql := range []string{
		"(SELECT a FROM t FOR UPDATE) UNION SELECT a FROM t",
		"(SELECT a FROM t FOR UPDATE)",
		"SELECT * FROM ((SELECT a FROM t FOR UPDATE) UNION SELECT a FROM t)",
	} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) || !strings.HasPrefix(sql[pe.Position:], "FOR UPDATE") {
			t.Errorf("Parse(%q) = %v, want an error at FOR UPDATE", sql, err)
		}
	}
	ParseAndCheck(t, "(SELECT a FROM t) FOR UPDATE")
	// FOR UPDATE follows the row limiting clause (ORA-03049 before it) and
	// locks the rows of FETCH FIRST n ROWS ONLY alone: an OFFSET, a PERCENT,
	// WITH TIES, or an ORDER BY with the row limit fails (ORA-02014), and so
	// do DISTINCT, UNIQUE, GROUP BY, and a CURSOR expression in the select
	// list (ORA-01786), parenthesized or not.
	for _, sql := range []string{
		"SELECT a FROM t FETCH FIRST 1 ROW ONLY FOR UPDATE",
		"SELECT a FROM t FETCH NEXT 2 ROWS ONLY FOR UPDATE SKIP LOCKED",
		"(SELECT a FROM t FETCH FIRST 1 ROW ONLY) FOR UPDATE",
		"(SELECT a FROM t) FETCH FIRST 1 ROW ONLY FOR UPDATE",
		"DECLARE CURSOR c IS SELECT a FROM t FETCH FIRST 10 ROWS ONLY FOR UPDATE SKIP LOCKED; BEGIN NULL; END;",
		"SELECT a FROM t WHERE a IN (SELECT DISTINCT b FROM u GROUP BY b) FOR UPDATE",
	} {
		ParseAndCheck(t, sql)
	}
	for _, sql := range []string{
		"SELECT a FROM t FETCH FIRST 50 PERCENT ROWS ONLY FOR UPDATE",
		"SELECT a FROM t ORDER BY a FETCH FIRST 1 ROW ONLY FOR UPDATE",
		"SELECT a FROM t OFFSET 1 ROWS FOR UPDATE",
		"SELECT a FROM t OFFSET 1 ROWS FETCH NEXT 1 ROW ONLY FOR UPDATE",
		"(SELECT a FROM t ORDER BY a FETCH FIRST 1 ROW ONLY) FOR UPDATE",
		"(SELECT a FROM t) ORDER BY a FETCH FIRST 1 ROW ONLY FOR UPDATE",
		"(SELECT a FROM t OFFSET 1 ROWS) FOR UPDATE",
		"SELECT a FROM t FETCH FIRST 1 ROWS WITH TIES FOR UPDATE",
		"(SELECT a FROM t FETCH FIRST 1 ROWS WITH TIES) FOR UPDATE",
		"SELECT DISTINCT a FROM t FOR UPDATE",
		"SELECT UNIQUE a FROM t FOR UPDATE",
		"(SELECT DISTINCT a FROM t) FOR UPDATE",
		"SELECT a FROM t GROUP BY a FOR UPDATE",
		"(SELECT a FROM t GROUP BY a) FOR UPDATE",
		"SELECT CURSOR(SELECT 1 FROM dual) FROM dual FOR UPDATE",
		"(SELECT CURSOR(SELECT 1 FROM dual) FROM dual) FOR UPDATE",
		"SELECT id, CURSOR(SELECT 1 FROM dual) c FROM e FOR UPDATE",
	} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) || !strings.HasPrefix(sql[pe.Position:], "FOR UPDATE") {
			t.Errorf("Parse(%q) = %v, want an error at FOR UPDATE", sql, err)
		}
	}
	ParseShouldFail(t, "SELECT a FROM t FOR UPDATE FETCH FIRST 1 ROW ONLY")
	// Nor does any subquery but a cursor's lock rows (ORA-00907, ORA-03048
	// at FOR): Oracle 23ai takes FOR UPDATE on a top-level query, a cursor
	// declaration, OPEN FOR, and a FOR loop's query alone.
	for _, sql := range []string{
		"SELECT 1 FROM dual WHERE EXISTS ((SELECT a FROM t) FOR UPDATE)",
		"SELECT 1 FROM dual WHERE EXISTS (SELECT a FROM t FOR UPDATE)",
		"SELECT 1 FROM dual WHERE 1 IN (SELECT a FROM t FOR UPDATE)",
		"SELECT (SELECT a FROM t FOR UPDATE) FROM dual",
		"SELECT * FROM (SELECT a FROM t FOR UPDATE)",
		"SELECT * FROM dual d, LATERAL (SELECT dummy FROM dual FOR UPDATE)",
		"SELECT * FROM dual d CROSS APPLY (SELECT dummy FROM dual FOR UPDATE)",
		"SELECT * FROM dual d OUTER APPLY (SELECT dummy FROM dual FOR UPDATE)",
		"WITH w AS (SELECT a FROM t FOR UPDATE) SELECT * FROM w",
		"SELECT y, s FROM t MODEL REFERENCE r ON (SELECT c, rate FROM u FOR UPDATE) DIMENSION BY (c) MEASURES (rate) MAIN m DIMENSION BY (y) MEASURES (s) RULES (s[2000] = r.rate['A'])",
		"SELECT y, s FROM t MODEL REFERENCE r ON ((SELECT c, rate FROM u) FOR UPDATE) DIMENSION BY (c) MEASURES (rate) MAIN m DIMENSION BY (y) MEASURES (s) RULES (s[2000] = r.rate['A'])",
		"MERGE INTO t USING (SELECT a FROM t FOR UPDATE) s ON (t.a = s.a) WHEN MATCHED THEN UPDATE SET t.a = 1",
		"INSERT INTO t SELECT a FROM t FOR UPDATE",
		"SELECT CURSOR(SELECT a FROM t FOR UPDATE) FROM dual",
		"CREATE VIEW v AS SELECT a FROM t FOR UPDATE",
		"CREATE TABLE c AS SELECT a FROM t FOR UPDATE",
		"CREATE MATERIALIZED ZONEMAP z AS SELECT SYS_OP_ZONE_ID(rowid), MIN(a), MAX(a) FROM t GROUP BY SYS_OP_ZONE_ID(rowid) FOR UPDATE",
		"CREATE MATERIALIZED ZONEMAP z AS (SELECT SYS_OP_ZONE_ID(rowid), MIN(a), MAX(a) FROM t GROUP BY SYS_OP_ZONE_ID(rowid)) FOR UPDATE",
		"CREATE JSON RELATIONAL DUALITY VIEW v AS (SELECT JSON_OBJECT('_id' VALUE e.id) FROM e e) FOR UPDATE",
		"CREATE JSON RELATIONAL DUALITY VIEW v AS SELECT JSON_OBJECT('_id' VALUE e.id) FROM e e FOR UPDATE",
		"UPDATE (SELECT a FROM t FOR UPDATE) SET a = 1",
		"DELETE FROM (SELECT a FROM t FOR UPDATE)",
		"UPDATE t SET a = (SELECT 1 FROM dual FOR UPDATE)",
	} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) || !strings.HasPrefix(sql[pe.Position:], "FOR UPDATE") {
			t.Errorf("Parse(%q) = %v, want an error at FOR UPDATE", sql, err)
		}
	}
	for _, sql := range []string{
		"SELECT * FROM (SELECT a FROM t) FOR UPDATE",
		"SELECT a FROM t WHERE a IN (SELECT a FROM t) FOR UPDATE",
		"DECLARE CURSOR c IS SELECT a FROM t FOR UPDATE; BEGIN NULL; END;",
		"DECLARE c SYS_REFCURSOR; BEGIN OPEN c FOR SELECT a FROM t FOR UPDATE; CLOSE c; END;",
		"DECLARE n NUMBER; BEGIN SELECT a INTO n FROM t FOR UPDATE; END;",
		"BEGIN FOR r IN (SELECT a FROM t FOR UPDATE) LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN VALUES OF (SELECT a FROM t FOR UPDATE) LOOP NULL; END LOOP; END;",
		"BEGIN FOR r IN ((SELECT a FROM t) FOR UPDATE) LOOP NULL; END LOOP; END;",
	} {
		ParseAndCheck(t, sql)
	}
	// OPEN FOR takes a query bare, parenthesized not even whole (PLS-00103
	// at its SELECT on Oracle 23ai), and a parenthesized string as dynamic
	// SQL.
	for _, sql := range []string{
		"DECLARE c SYS_REFCURSOR; BEGIN OPEN c FOR (SELECT q FROM s); CLOSE c; END;",
		"DECLARE c SYS_REFCURSOR; BEGIN OPEN c FOR ((SELECT q FROM s) ORDER BY q); CLOSE c; END;",
		"DECLARE c SYS_REFCURSOR; BEGIN OPEN c FOR ((SELECT q FROM s) UNION (SELECT q FROM u)); CLOSE c; END;",
	} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) || !strings.HasPrefix(sql[pe.Position:], "SELECT q FROM s") {
			t.Errorf("Parse(%q) = %v, want an error at the first SELECT", sql, err)
		}
	}
	ParseAndCheck(t, "DECLARE c SYS_REFCURSOR; s VARCHAR2(100); BEGIN OPEN c FOR (('SELECT 1 FROM dual' || s)); CLOSE c; END;")
	ParseAndCheck(t, "DECLARE c SYS_REFCURSOR; BEGIN OPEN c FOR SELECT q FROM s UNION SELECT q FROM u; CLOSE c; END;")

	// A FOR loop's query stands alone, the whole of its iteration control
	// (PLS-00103 at what follows it).
	for _, tc := range []struct{ sql, at string }{
		{"BEGIN FOR i IN (SELECT 1 FROM dual) + 1 LOOP NULL; END LOOP; END;", "+ 1"},
		{"BEGIN FOR i IN (SELECT 1 FROM dual ORDER BY 1) + 1 LOOP NULL; END LOOP; END;", "+ 1"},
		{"BEGIN FOR i IN (SELECT 1 FROM dual ORDER BY 1) .. 3 LOOP NULL; END LOOP; END;", ".. 3"},
		{"BEGIN FOR i IN VALUES OF (SELECT 1 FROM dual) + 1 LOOP NULL; END LOOP; END;", "+ 1"},
	} {
		_, err := Parse(tc.sql)
		var pe *ParseError
		if !errors.As(err, &pe) || !strings.HasPrefix(tc.sql[pe.Position:], tc.at) {
			t.Errorf("Parse(%q) = %v, want an error at %q", tc.sql, err, tc.at)
		}
	}
	// No set operation takes a FOR UPDATE (ORA-01786 at FOR).
	for _, sql := range []string{
		"SELECT a FROM t EXCEPT SELECT a FROM t FOR UPDATE",
		"SELECT a FROM t UNION SELECT a FROM t FOR UPDATE",
		"SELECT a FROM t MINUS ALL SELECT a FROM t FOR UPDATE OF a",
		"(SELECT a FROM t) UNION (SELECT a FROM t) FOR UPDATE",
		"(SELECT a FROM t UNION SELECT a FROM t) FOR UPDATE",
		"(SELECT a FROM t EXCEPT SELECT a FROM t) FOR UPDATE",
		"((SELECT a FROM t) UNION (SELECT a FROM t)) FOR UPDATE",
	} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) || !strings.HasPrefix(sql[pe.Position:], "FOR UPDATE") {
			t.Errorf("Parse(%q) = %v, want an error at FOR UPDATE", sql, err)
		}
	}
	ParseAndCheck(t, "SELECT a FROM t UNION SELECT a FROM t ORDER BY 1 FETCH FIRST 1 ROWS ONLY")

	// A subquery that stands for a value or a condition takes no plain
	// ORDER BY either, parenthesized or not (ORA-00907, here at its first
	// item); a row-limited one, an inline view, MERGE USING, CURSOR, and the
	// query of a cursor FOR loop do.
	for _, sql := range []string{
		"SELECT 1 FROM dual WHERE 1 IN (SELECT 1 FROM dual ORDER BY 1)",
		"SELECT 1 FROM dual WHERE 1 IN ((SELECT 1 FROM dual) ORDER BY 1)",
		"SELECT 1 FROM dual WHERE 1 IN ((SELECT 1 FROM dual ORDER BY 1))",
		"SELECT 1 FROM dual WHERE 1 IN (SELECT 1 FROM dual UNION SELECT 2 FROM dual ORDER BY 1)",
		"SELECT 1 FROM dual WHERE (1, 1) IN (SELECT 1, 1 FROM dual ORDER BY 1)",
		"SELECT 1 FROM dual WHERE EXISTS ((SELECT 1 FROM dual) ORDER BY 1)",
		"SELECT 1 FROM dual WHERE 1 = ANY (SELECT 1 FROM dual ORDER BY 1)",
		"SELECT ((SELECT 1 FROM dual) ORDER BY 1) FROM dual",
		"SELECT ABS((SELECT 1 FROM dual ORDER BY 1)) FROM dual",
		"SELECT 1 FROM dual WHERE 1 IN (SELECT 1 FROM dual ORDER SIBLINGS BY 1)",
		"UPDATE t SET (a, b) = (SELECT 1, 2 FROM dual ORDER BY 1)",
		"DECLARE n NUMBER; BEGIN SELECT COUNT(*) INTO n FROM dual WHERE 1 IN (SELECT 1 FROM dual ORDER BY 1); END;",
		"BEGIN FOR r IN (SELECT (SELECT 1 FROM dual ORDER BY 1) a FROM dual) LOOP NULL; END LOOP; END;",
	} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) || !strings.HasSuffix(sql[:pe.Position], "BY ") {
			t.Errorf("Parse(%q) = %v, want an error at the ORDER BY", sql, err)
		}
	}
	for _, sql := range []string{
		"SELECT 1 FROM dual WHERE 1 IN ((SELECT 1 FROM dual) ORDER BY 1 FETCH FIRST 1 ROWS ONLY)",
		"SELECT 1 FROM dual WHERE 1 IN (SELECT 1 FROM dual ORDER BY 1 OFFSET 0 ROWS)",
		"SELECT 1 FROM dual WHERE EXISTS (SELECT 1 FROM dual ORDER BY 1 FETCH FIRST 1 ROWS ONLY)",
		"SELECT CURSOR(SELECT 1 FROM dual ORDER BY 1) FROM dual",
		"SELECT * FROM (SELECT 1 a FROM dual ORDER BY 1)",
		"MERGE INTO t USING (SELECT 1 a FROM dual ORDER BY 1) s ON (t.a = s.a) WHEN MATCHED THEN UPDATE SET t.b = 1",
		"BEGIN FOR r IN (SELECT 1 a FROM dual ORDER BY 1) LOOP NULL; END LOOP; END;",
		"BEGIN FOR r IN (SELECT 1 a FROM dual UNION SELECT 2 FROM dual ORDER BY 1) LOOP NULL; END LOOP; END;",
		"BEGIN FOR r IN ((SELECT 1 a FROM dual) ORDER BY 1) LOOP NULL; END LOOP; END;",
		"BEGIN FOR r IN (SELECT 1 a FROM dual WHERE 1 IN (SELECT 1 FROM dual) ORDER BY 1) LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN VALUES OF (SELECT 1 a FROM dual ORDER BY 1) LOOP NULL; END LOOP; END;",
	} {
		ParseAndCheck(t, sql)
	}

	// No query block selects * without FROM, wherever it stands, and the
	// error is at its SELECT (ORA-54058); but for the one block of an
	// EXISTS subquery, whose select list Oracle does not expand unless a
	// row-limiting clause follows it.
	for _, sql := range []string{
		"CREATE TABLE t AS (SELECT *)",
		"CREATE VIEW v AS (SELECT *)",
		"INSERT INTO t (SELECT *)",
		"SELECT (SELECT * ) FROM dual",
		"SELECT 1 FROM dual WHERE 1 IN (SELECT *)",
		"SELECT * FROM (SELECT *)",
		"WITH w AS (SELECT *) SELECT 1 FROM dual",
		"SELECT 1 FROM dual WHERE EXISTS (SELECT * UNION SELECT 1 FROM dual)",
		"SELECT 1 FROM dual WHERE EXISTS (SELECT 1 FROM dual UNION SELECT *)",
		"SELECT 1 FROM dual WHERE EXISTS (SELECT * FROM (SELECT *))",
		"SELECT 1 FROM dual WHERE EXISTS (SELECT * FETCH FIRST 1 ROWS ONLY)",
		"SELECT 1 FROM dual WHERE EXISTS (SELECT * OFFSET 0 ROWS)",
		"SELECT 1 FROM dual WHERE EXISTS ((SELECT *) FETCH FIRST 1 ROWS ONLY)",
		"SELECT 1 FROM dual WHERE EXISTS ((SELECT *) OFFSET 0 ROWS)",
		"MERGE INTO t USING (SELECT *) s ON (1 = 1) WHEN MATCHED THEN UPDATE SET t.a = 1",
		"BEGIN FOR r IN (SELECT *) LOOP NULL; END LOOP; END;",
	} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) = %v, want a ParseError", sql, err)
			continue
		}
		if at := sql[pe.Position:]; !strings.HasPrefix(at, "SELECT *") {
			t.Errorf("Parse(%q) fails at %q, want at SELECT *", sql, at)
		}
	}
	// So does table.*, correlated or not (ORA-00904).
	for _, sql := range []string{
		"CREATE VIEW v AS (SELECT x.*)",
		"SELECT (SELECT d.*) FROM dual d",
		"SELECT 1 FROM dual d WHERE 'X' IN (SELECT d.*)",
		"SELECT 1 FROM dual WHERE EXISTS (SELECT x.* UNION SELECT 1 FROM dual)",
	} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) || !strings.HasPrefix(sql[pe.Position:], "SELECT ") {
			t.Errorf("Parse(%q) = %v, want an error at its SELECT", sql, err)
		}
	}
	// The text may end in a parenthesized block, which is whole: the error
	// is at its SELECT, ';' or not. An unparenthesized one there has yet to
	// get its FROM, and the error is at the end.
	for _, sql := range []string{
		"SELECT 1 FROM dual UNION (SELECT *)",
		"SELECT 1 FROM dual UNION (SELECT *);",
		"(SELECT *) UNION SELECT 1 FROM dual",
		// The first failing block, though the last awaits its FROM.
		"(SELECT *) UNION SELECT *",
		"SELECT * UNION SELECT *",
	} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) || !strings.HasPrefix(sql[pe.Position:], "SELECT *") {
			t.Errorf("Parse(%q) = %v, want an error at SELECT *", sql, err)
		}
	}
	for _, sql := range []string{"SELECT *", "SELECT 1 FROM dual UNION SELECT *"} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Position != len(sql) {
			t.Errorf("Parse(%q) = %v, want an error at the end", sql, err)
		}
	}
	// A block with a WITH clause starts there, where Oracle reports it.
	const withStar = "CREATE VIEW v AS WITH w AS (SELECT 1 FROM dual) SELECT * ORDER BY 1"
	var pe *ParseError
	if _, err := Parse(withStar); !errors.As(err, &pe) || !strings.HasPrefix(withStar[pe.Position:], "WITH w") {
		t.Errorf("Parse(%q) = %v, want an error at WITH", withStar, err)
	}
	for _, sql := range []string{
		"SELECT 1 FROM dual WHERE EXISTS (SELECT *)",
		"SELECT 1 FROM dual WHERE NOT EXISTS (SELECT *)",
		"INSERT INTO t SELECT 1 FROM dual WHERE EXISTS (SELECT *)",
		"DECLARE n NUMBER; BEGIN SELECT 1 INTO n FROM dual WHERE EXISTS (SELECT *); END;",
		"SELECT 1 FROM dual WHERE EXISTS ((SELECT *))",
		"SELECT 1 FROM dual WHERE EXISTS (WITH w AS (SELECT 1 c FROM dual) SELECT *)",
		"SELECT 1 FROM dual WHERE EXISTS (SELECT x.*)",
		"SELECT 1 FROM dual d WHERE EXISTS (SELECT d.*)",
		"SELECT 1 FROM dual WHERE EXISTS (WITH w AS (SELECT 1 c FROM dual) (SELECT *))",
	} {
		ParseAndCheck(t, sql)
	}

	// An inline view whose first operand is parenthesized exposes that
	// operand's columns, as an unparenthesized one does.
	// MULTISET EXCEPT is an operator inside the select list, not a set
	// operator ending it, but for a column named MULTISET a query follows.
	for _, sql := range []string{
		"SELECT x. FROM ((SELECT c1 AS left_name FROM t) UNION (SELECT c2 FROM u)) x",
		"SELECT x. FROM (SELECT c1 AS left_name FROM t UNION SELECT c2 FROM u) x",
		"SELECT x. FROM (SELECT a MULTISET EXCEPT b AS left_name FROM t) x",
		"SELECT x. FROM (SELECT a MULTISET EXCEPT DISTINCT b AS left_name FROM t) x",
		"SELECT x. FROM (SELECT c1 left_name FROM t EXCEPT SELECT c FROM u) x",
		// After a WITH clause too.
		"SELECT x. FROM (WITH w AS (SELECT 1 c FROM dual) (SELECT c AS left_name FROM w)) x",
		"SELECT x. FROM (WITH w (c) AS (SELECT 1 FROM dual), v AS (SELECT 2 d FROM dual) ((SELECT c AS left_name FROM w) UNION (SELECT d FROM v))) x",
		"SELECT x. FROM (WITH w AS (SELECT 1 c FROM dual) SELECT c AS left_name FROM w) x",
		// The first operand, parenthesized, though one at the query's own
		// depth follows.
		"SELECT x. FROM ((SELECT c1 AS left_name FROM t) UNION SELECT c2 FROM u) x",
	} {
		scope := CollectCompletion(sql, len("SELECT x.")).Scope
		if scope == nil || len(scope.LocalReferences) != 1 || len(scope.LocalReferences[0].Columns) != 1 || scope.LocalReferences[0].Columns[0] != "LEFT_NAME" {
			t.Errorf("CollectCompletion(%q) scope = %+v, want X with LEFT_NAME", sql, scope)
		}
	}

	// PL/SQL does not reserve EXCEPT: Oracle 23ai compiles each of these.
	for _, sql := range []string{
		"DECLARE except NUMBER := 1; x NUMBER; BEGIN x := except + 1; except := x; p(except); p(a => except); END;",
		"DECLARE TYPE r IS RECORD (a NUMBER); except r; BEGIN except.a := 1; IF except.a > 0 THEN NULL; END IF; END;",
		"BEGIN FOR except IN 1..3 LOOP NULL; END LOOP; END;",
	} {
		ParseAndCheck(t, sql)
	}
	// Reserved, EXCEPT still follows MULTISET.
	ParseAndCheck(t, "SELECT a MULTISET EXCEPT b FROM t")

	// Any number of wrappers come off in one pass: 32,000 took 6.8 s when
	// each was peeled with a scan of its own.
	deep := "SELECT x. FROM " + strings.Repeat("(", 32000) + "SELECT c FROM t" + strings.Repeat(")", 32000) + " x"
	start := time.Now()
	if scope := CollectCompletion(deep, len("SELECT x.")).Scope; scope == nil || len(scope.LocalReferences) != 1 ||
		len(scope.LocalReferences[0].Columns) != 1 || scope.LocalReferences[0].Columns[0] != "C" {
		t.Errorf("CollectCompletion through 32,000 wrappers: scope = %+v, want X with C", scope)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("CollectCompletion through 32,000 wrappers took %v", elapsed)
	}
	if scope := CollectCompletion("SELECT x. FROM (SELECT multiset EXCEPT SELECT c FROM t) x", len("SELECT x.")).Scope; scope == nil ||
		len(scope.LocalReferences) != 1 || len(scope.LocalReferences[0].Columns) != 1 || scope.LocalReferences[0].Columns[0] != "MULTISET" {
		t.Errorf("CollectCompletion past a column named MULTISET: scope = %+v, want X with MULTISET", scope)
	}

	// Completion takes EXCEPT as the set operator it is, as MINUS.
	for _, op := range []string{"EXCEPT", "MINUS"} {
		before := "SELECT * FROM t " + op + " SELECT u."
		got := CollectCompletion(before+" FROM u", len(before)).Intent
		if got == nil || got.Qualifier.Object != "U" || got.Qualifier.Schema != "" {
			t.Errorf("CollectCompletion after %s = %+v, want the columns of U", op, got)
		}
	}
}

// TestParenthesizedQueryAST checks how a parenthesized subquery is
// recorded: alone it is the query itself; with a set operator, a clause
// after it, or a WITH clause before it, a query around it holds it as Larg.
func TestParenthesizedQueryAST(t *testing.T) {
	const bare = "((SELECT a FROM t))"
	plain := parseOne[*ast.SelectStmt](t, bare)
	if plain.Larg != nil || plain.TargetList.Len() != 1 {
		t.Errorf("%s = Larg %v, %d targets; want the query itself", bare, plain.Larg, plain.TargetList.Len())
	}
	if plain.Loc.Start != 0 || plain.Loc.End != len(bare) {
		t.Errorf("%s spans %d:%d, want its parentheses, 0:%d", bare, plain.Loc.Start, plain.Loc.End, len(bare))
	}

	union := parseOne[*ast.SelectStmt](t, "(SELECT a FROM t) UNION ALL (SELECT b FROM u)")
	if union.Op != ast.SETOP_UNION || !union.SetAll || union.Larg == nil || union.Rarg == nil ||
		union.Larg.TargetList.Len() != 1 || union.Rarg.TargetList.Len() != 1 {
		t.Fatalf("union = op %v all %v larg %v rarg %v", union.Op, union.SetAll, union.Larg, union.Rarg)
	}

	ordered := parseOne[*ast.SelectStmt](t, "(SELECT a FROM t) ORDER BY 1")
	if ordered.Larg == nil || ordered.OrderBy.Len() != 1 || ordered.Op != ast.SETOP_NONE {
		t.Errorf("ordered = larg %v order %v op %v", ordered.Larg, ordered.OrderBy, ordered.Op)
	}

	except := parseOne[*ast.SelectStmt](t, "SELECT a FROM t EXCEPT ALL SELECT b FROM u")
	if except.Op != ast.SETOP_EXCEPT || !except.SetAll {
		t.Errorf("EXCEPT ALL = op %v all %v", except.Op, except.SetAll)
	}

	insert := parseOne[*ast.InsertStmt](t, "INSERT INTO t (SELECT 1 FROM dual)")
	if insert.Columns != nil || insert.Select == nil {
		t.Errorf("INSERT INTO t (SELECT ...) = columns %v, select %v", insert.Columns, insert.Select)
	}
}
