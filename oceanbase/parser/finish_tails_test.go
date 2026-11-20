package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
)

func TestOceanBaseFinishTails(t *testing.T) {
	cases := []struct {
		sql    string
		accept bool
	}{
		// VIRTUAL_COLUMN_ID index option.
		{"ALTER TABLE a ADD KEY(a) VIRTUAL_COLUMN_ID 1", true},
		{"ALTER TABLE a ADD KEY(a) VIRTUAL_COLUMN_ID = 1", true},
		{"CREATE TABLE a (b INT, KEY(b) VIRTUAL_COLUMN_ID = 1)", true},
		{"CREATE INDEX idx ON a (b) VIRTUAL_COLUMN_ID = 1", true},
		{"ALTER TABLE a ADD KEY(a) VIRTUAL_COLUMN_ID a", false},
		{"ALTER TABLE a ADD KEY(a) VIRTUAL_COLUMN_ID", false},

		// GEOMETRYCOLLECTION and GEOMCOLLECTION.
		{"ALTER TABLE a ADD a GEOMETRYCOLLECTION", true},
		{"ALTER TABLE a ADD a GEOMCOLLECTION", true},
		{"CREATE TABLE t (g GEOMCOLLECTION NOT NULL)", true},

		// INTERVAL expression boundaries.
		{"SELECT + INTERVAL a DAY", false},
		{"SELECT - INTERVAL a DAY", false},
		{"SELECT ~ INTERVAL a DAY", false},
		{"SELECT ! INTERVAL a DAY", false},
		{"SELECT NOT INTERVAL a DAY", false},
		{"SELECT BINARY INTERVAL a DAY", false},
		{"SELECT INTERVAL a DAY", false},
		{"SELECT INTERVAL a DAY a", false},
		{"SELECT(a, INTERVAL a DAY)", false},
		{"SELECT (INTERVAL 1 DAY) + 1", false},
		{"SELECT 1 + (INTERVAL 1 DAY)", false},
		{"SELECT INTERVAL 1 DAY - 1", false},
		{"SELECT INTERVAL 1 DAY + INTERVAL 1 HOUR", false},
		{"SELECT 1 + INTERVAL a DAY", true},
		{"SELECT 1 - INTERVAL a DAY", true},
		{"SELECT INTERVAL 1 DAY + DATE '2020-01-01'", true},
		{"SELECT INTERVAL 1 DAY + INTERVAL 1 HOUR + DATE '2020-01-01'", true},
		{"SELECT + INTERVAL(1, 2)", true},
		{"SELECT - INTERVAL(1, 2)", true},
		{"SELECT(a, INTERVAL(1, 2))", true},

		// ROWID keyword call refusal.
		{"SELECT ROWID(a, a)", false},
		{"SELECT ROWID()", false},
		{"SELECT `ROWID`(a, a)", true},
		{"SELECT d.ROWID(a, a)", true},
		{"SELECT ROWID FROM t", true},

		// Spaced qualified function calls vs id_dot_id.
		{"SELECT a . STATS_AUTO_RECALC()", false},
		{"SELECT a.STATS_AUTO_RECALC()", true},
		{"SELECT a . `STATS_AUTO_RECALC`()", true},
		{"SELECT a . STATISTICS()", true},
		{"SELECT a . CURRENT_USER()", true},
		{"SELECT a . CURRENT_USER() zzz", true},
		{"SELECT a . CURRENT_USER(1)", false},

		// ARRAY_* lambda functions and JSON -> / ->> bounds.
		{"SELECT ARRAY_FILTER(a)", false},
		{"SELECT ARRAY_MAP(a)", false},
		{"SELECT ARRAY_FIRST(a)", false},
		{"SELECT ARRAY_SORTBY(a)", false},
		{"SELECT `ARRAY_FILTER`(a)", true},
		{"SELECT d.ARRAY_FILTER(a)", true},
		{"SELECT ARRAY_FILTER(a -> a, b)", true},
		{"SELECT ARRAY_MAP((a, b) -> a + b, x, y)", true},
		{"SELECT(a -> a, a)", false},
		{"SELECT a -> a", false},
		{"SELECT (a) -> '$.b'", false},
		{"SELECT c.d.t.a -> '$.b'", false},
		{"SELECT a -> '$.b'", true},
		{"SELECT t.a ->> '$.b'", true},
		{"SELECT d.t.a -> '$.b'", true},
		{"SELECT a -> '$.a' 'b'", true},
		{"SELECT a -> N'$.b'", true},
		{"SELECT a -> _utf8mb4'$.b'", true},

		// expr_as_list on function_name / id_dot_id calls.
		{"SELECT STATISTICS(a zzz)", true},
		{"SELECT STATISTICS(a, a zzz)", true},
		{"SELECT STATISTICS(a a)", true},
		{"SELECT STATISTICS(a 's')", true},
		{"SELECT STATISTICS(a AS b)", true},
		{"SELECT STATISTICS(t.* a)", false},
		{"SELECT STATISTICS(t.* AS a)", false},
		{"SELECT STATISTICS(t.*)", true},
		{"SELECT LEFT(a a, 1)", false},
		{"SELECT CALC_PARTITION_ID(a a)", false},
		{"SELECT SPLIT(a a)", false},
		{"SELECT PERCENTILE_CONT(a a)", false},

		// DELETE modifiers and TABLE(...) collection targets.
		{"DELETE QUICK QUICK zzz FROM a", true},
		{"DELETE LOW_PRIORITY QUICK LOW_PRIORITY FROM a", true},
		{"DELETE IGNORE LOW_PRIORITY QUICK FROM a", true},
		{"DELETE IGNORE IGNORE FROM a", false},
		{"DELETE QUICK IGNORE FROM a", false},
		{"DELETE FROM TABLE(a)", true},
		{"DELETE FROM TABLE(a) zzz", true},
		{"DELETE FROM TABLE(a) a", true},
		{"DELETE FROM TABLE(a) AS a", true},
		{"SELECT * FROM TABLE(a) AS b", true},
		{"DELETE FROM TABLE(1 + 1)", false},
		{"DELETE FROM TABLE(a) 's'", false},
		{"DELETE FROM TABLE(a) USING b", false},

		// WITH RECURSIVE as CTE relation_name.
		{"WITH RECURSIVE AS(SELECT a) SELECT a", true},
		{"WITH RECURSIVE (c) AS (SELECT a) SELECT c FROM RECURSIVE", true},
		{"WITH RECURSIVE c AS (SELECT a) SELECT * FROM c", true},

		// CREATE USER CURRENT_USER vs CURRENT_USER() and REQUIRE string runs.
		{"CREATE USER CURRENT_USER() IDENTIFIED BY 's'", false},
		{"CREATE USER CURRENT_USER IDENTIFIED BY PASSWORD 's'", false},
		{"CREATE USER CURRENT_USER() IDENTIFIED AS 's'", false},
		{"CREATE USER CURRENT_USER IDENTIFIED BY 's'", true},
		{"CREATE USER CURRENT_USER() IDENTIFIED BY PASSWORD 's'", true},
		{"ALTER USER CURRENT_USER() IDENTIFIED BY 's'", true},
		{"CREATE USER 's' REQUIRE CIPHER 's' 's'", true},

		// /*HINT+ ... */ (HINT_VALUE) placements.
		{"/*HINT+ a */", false},
		{"SELECT /*HINT+ a */ 1", false},
		{"SELECT 1 WHERE /*HINT+ a */ 1 = 1", true},
		{"UPDATE t SET a = 1 WHERE /*hint+ a */ b = 2", true},
		{"DELETE FROM t WHERE /*HINT+ a */ b = 2", true},
		{"BEGIN /*HINT+ a */ WORK", true},
		{"START /*HINT+ a */ TRANSACTION", true},
		{"COMMIT /*HINT+ a */ WORK", true},
		{"COMMIT /*HINT+ a */ AND CHAIN", false},
		{"ROLLBACK /*HINT+ a */ WORK", true},
		{"ROLLBACK /*HINT+ a */ TO sp", false},

		// KILL [CONNECTION | QUERY] expr vs KILL expr.
		{"KILL CONNECTION", true},
		{"KILL QUERY", true},
		{"KILL CONNECTION 1", true},
		{"KILL QUERY a", true},
		{"KILL CONNECTION + 1", true},
		{"KILL CONNECTION = 1", true},

		// Bare DEFAULT restricted to expr_or_default / set_expr_or_default.
		{"SELECT DEFAULT", false},
		{"SELECT 1 + DEFAULT", false},
		{"SELECT (DEFAULT)", false},
		{"SELECT DEFAULT(a) FROM t", true},
		{"INSERT INTO t VALUES (DEFAULT)", true},
		{"VALUES ROW(DEFAULT)", true},
		{"UPDATE t SET a = DEFAULT", true},
		{"LOAD DATA INFILE 'f' INTO TABLE t SET a = DEFAULT", true},
		{"SET a = DEFAULT", true},
		{"CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p0 VALUES IN (DEFAULT))", true},
	}

	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if got := err == nil; got != tc.accept {
				t.Fatalf("Parse(%q): accept=%v, want %v (err=%v)", tc.sql, got, tc.accept, err)
			}
			if err == nil && list != nil {
				if s := ast.NodeToString(list); s == "" {
					t.Fatalf("NodeToString(%q) returned empty string", tc.sql)
				}
			}
		})
	}
}

func TestOceanBaseFinishASTNodes(t *testing.T) {
	// VIRTUAL_COLUMN_ID retains its IndexOption in NodeToString.
	alterList, err := Parse("ALTER TABLE a ADD KEY(a) VIRTUAL_COLUMN_ID = 7")
	if err != nil {
		t.Fatal(err)
	}
	if s := ast.NodeToString(alterList); !strings.Contains(s, ":name VIRTUAL_COLUMN_ID :val {INT_LIT :val 7") {
		t.Errorf("VIRTUAL_COLUMN_ID missing from NodeToString: %s", s)
	}

	// GEOMCOLLECTION normalizes DataType.Name to GEOMETRYCOLLECTION.
	geomList, err := Parse("ALTER TABLE a ADD a GEOMCOLLECTION")
	if err != nil {
		t.Fatal(err)
	}
	alter := geomList.Items[0].(*ast.AlterTableStmt)
	if got := alter.Commands[0].Column.TypeName.Name; got != "GEOMETRYCOLLECTION" {
		t.Errorf("GEOMCOLLECTION TypeName.Name = %q, want GEOMETRYCOLLECTION", got)
	}

	// TableCollectionExpr and LambdaExpr serialize distinctly in NodeToString.
	delList, err := Parse("DELETE FROM TABLE(a) AS b")
	if err != nil {
		t.Fatal(err)
	}
	if s := ast.NodeToString(delList); !strings.Contains(s, "{TABLE_COLLECTION") || !strings.Contains(s, ":alias b") {
		t.Errorf("TableCollectionExpr missing from NodeToString: %s", s)
	}

	lamList, err := Parse("SELECT ARRAY_MAP((x, y) -> x + y, a, b)")
	if err != nil {
		t.Fatal(err)
	}
	if s := ast.NodeToString(lamList); !strings.Contains(s, "{LAMBDA") || !strings.Contains(s, ":params (x y)") {
		t.Errorf("LambdaExpr missing from NodeToString: %s", s)
	}

	// WITH RECURSIVE AS (...) leaves Recursive=false and names the CTE RECURSIVE.
	cteList, err := Parse("WITH RECURSIVE AS (SELECT 1) SELECT * FROM RECURSIVE")
	if err != nil {
		t.Fatal(err)
	}
	sel := cteList.Items[0].(*ast.SelectStmt)
	if len(sel.CTEs) != 1 || sel.CTEs[0].Recursive || sel.CTEs[0].Name != "RECURSIVE" {
		t.Errorf("unexpected CTE: %+v", sel.CTEs)
	}

	// REQUIRE CIPHER 's' 't' folds adjacent string literals.
	cuList, err := Parse("CREATE USER 'u' REQUIRE CIPHER 's' 't'")
	if err != nil {
		t.Fatal(err)
	}
	cu := cuList.Items[0].(*ast.CreateUserStmt)
	if cu.Require == nil || cu.Require.Cipher != "st" {
		t.Errorf("Require.Cipher = %+v, want \"st\"", cu.Require)
	}

	// Split keeps /*HINT+ ... */ as a non-empty segment.
	segs := Split("/*HINT+ a */")
	if len(segs) != 1 || segs[0].Empty() {
		t.Errorf("Split(/*HINT+ a */) = %+v, want 1 non-empty segment", segs)
	}

	// KILL disambiguation: unary `+ 1` shifts into `KILL {CONNECTION|QUERY} expr`
	// (matching OceanBase 4.3.5.6's 1094 Unknown thread id: 1), while `= 1`
	// reduces CONNECTION as the column of `KILL expr` (1054 Unknown column).
	killPlus, err := Parse("KILL QUERY + 1")
	if err != nil {
		t.Fatal(err)
	}
	if ks := killPlus.Items[0].(*ast.KillStmt); !ks.Query {
		t.Errorf("KILL QUERY + 1: Query = false, want true")
	}
	killEq, err := Parse("KILL QUERY = 1")
	if err != nil {
		t.Fatal(err)
	}
	if ks := killEq.Items[0].(*ast.KillStmt); ks.Query {
		t.Errorf("KILL QUERY = 1: Query = true, want false")
	} else if be, ok := ks.ConnectionID.(*ast.BinaryExpr); !ok || be.Op != ast.BinOpEq {
		t.Errorf("KILL QUERY = 1: ConnectionID = %T, want *ast.BinaryExpr(=)", ks.ConnectionID)
	}
}
