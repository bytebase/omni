package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParenthesizedQueries checks queries whose operands are parenthesized
// subqueries, and the EXCEPT and ALL set operators, against Oracle 23ai:
// every accepted form parses there, and every rejected one fails with the
// error noted beside it. Rows ref_225 to ref_236, ref_251 to ref_252, and
// ref_257 run them against the engine.
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
