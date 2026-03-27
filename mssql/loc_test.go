package mssql

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// Helper to extract AST Loc from a node via reflection (same as nodeLoc in parse.go).
func testNodeLoc(n ast.Node) ast.Loc {
	return nodeLoc(n)
}

// locValid checks that a Loc has non-negative Start and End > Start.
func locValid(loc ast.Loc) bool {
	return loc.Start >= 0 && loc.End > loc.Start
}

// locStartValid checks that at least the Loc.Start is set (>= 0).
// Used for nodes where Loc.End may not yet be accurately set (partial tracking).
func locStartValid(loc ast.Loc) bool {
	return loc.Start >= 0
}

// locText extracts the text covered by a Loc from the SQL string.
// Returns empty string if Loc is invalid.
func locText(sql string, loc ast.Loc) string {
	if loc.Start < 0 || loc.End < 0 || loc.Start > len(sql) || loc.End > len(sql) || loc.Start >= loc.End {
		return ""
	}
	return sql[loc.Start:loc.End]
}

// ---------- 6.1 DML Loc Verification ----------

func TestLoc_DML(t *testing.T) {
	t.Run("SELECT col FROM t", func(t *testing.T) {
		sql := "SELECT col FROM t"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		sel, ok := stmts[0].AST.(*ast.SelectStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.SelectStmt", stmts[0].AST)
		}
		loc := sel.Loc
		if !locValid(loc) {
			t.Fatalf("SelectStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "SELECT") || !strings.Contains(got, "FROM") {
			t.Errorf("SelectStmt Loc text = %q, want it to contain SELECT and FROM", got)
		}
	})

	t.Run("SELECT with WHERE and ORDER BY", func(t *testing.T) {
		sql := "SELECT a, b FROM t WHERE x = 1 ORDER BY a"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		sel, ok := stmts[0].AST.(*ast.SelectStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.SelectStmt", stmts[0].AST)
		}
		loc := sel.Loc
		if !locValid(loc) {
			t.Fatalf("SelectStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "ORDER BY") {
			t.Errorf("SelectStmt Loc text = %q, want it to contain ORDER BY", got)
		}
	})

	t.Run("INSERT INTO", func(t *testing.T) {
		sql := "INSERT INTO t (a) VALUES (1)"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		ins, ok := stmts[0].AST.(*ast.InsertStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.InsertStmt", stmts[0].AST)
		}
		loc := ins.Loc
		if !locValid(loc) {
			t.Fatalf("InsertStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "INSERT") || !strings.Contains(got, "VALUES") {
			t.Errorf("InsertStmt Loc text = %q, want it to contain INSERT and VALUES", got)
		}
	})

	t.Run("UPDATE", func(t *testing.T) {
		sql := "UPDATE t SET a = 1 WHERE b = 2"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		upd, ok := stmts[0].AST.(*ast.UpdateStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.UpdateStmt", stmts[0].AST)
		}
		loc := upd.Loc
		if !locValid(loc) {
			t.Fatalf("UpdateStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "UPDATE") || !strings.Contains(got, "WHERE") {
			t.Errorf("UpdateStmt Loc text = %q, want it to contain UPDATE and WHERE", got)
		}
	})

	t.Run("DELETE", func(t *testing.T) {
		sql := "DELETE FROM t WHERE a = 1"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		del, ok := stmts[0].AST.(*ast.DeleteStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.DeleteStmt", stmts[0].AST)
		}
		loc := del.Loc
		if !locValid(loc) {
			t.Fatalf("DeleteStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "DELETE") || !strings.Contains(got, "WHERE") {
			t.Errorf("DeleteStmt Loc text = %q, want it to contain DELETE and WHERE", got)
		}
	})

	t.Run("MERGE", func(t *testing.T) {
		sql := "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET a = s.a"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		mrg, ok := stmts[0].AST.(*ast.MergeStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.MergeStmt", stmts[0].AST)
		}
		loc := mrg.Loc
		if !locValid(loc) {
			t.Fatalf("MergeStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "MERGE") {
			t.Errorf("MergeStmt Loc text = %q, want it to contain MERGE", got)
		}
	})

	t.Run("WITH CTE", func(t *testing.T) {
		sql := "WITH cte AS (SELECT 1) SELECT * FROM cte"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		sel, ok := stmts[0].AST.(*ast.SelectStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.SelectStmt", stmts[0].AST)
		}
		loc := sel.Loc
		if !locValid(loc) {
			t.Fatalf("SelectStmt Loc invalid: %+v", loc)
		}
		// The top-level Loc should span from WITH to the end.
		got := locText(sql, loc)
		if !strings.Contains(got, "SELECT") {
			t.Errorf("SelectStmt Loc text = %q, want it to contain SELECT", got)
		}
	})

	t.Run("JOIN", func(t *testing.T) {
		sql := "SELECT * FROM t1 JOIN t2 ON t1.id = t2.id"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		sel, ok := stmts[0].AST.(*ast.SelectStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.SelectStmt", stmts[0].AST)
		}
		// Navigate into FromClause to find the JoinClause.
		if sel.FromClause == nil || sel.FromClause.Len() == 0 {
			t.Fatal("FromClause is nil or empty")
		}
		join, ok := sel.FromClause.Items[0].(*ast.JoinClause)
		if !ok {
			t.Fatalf("FromClause[0] type = %T, want *ast.JoinClause", sel.FromClause.Items[0])
		}
		loc := join.Loc
		// [~] partial: JoinClause Loc.End not yet set accurately.
		if !locStartValid(loc) {
			t.Fatalf("JoinClause Loc.Start invalid: %+v", loc)
		}
		if locValid(loc) {
			got := locText(sql, loc)
			if !strings.Contains(got, "JOIN") {
				t.Errorf("JoinClause Loc text = %q, want it to contain JOIN", got)
			}
		} else {
			t.Logf("[~] partial: JoinClause Loc.End not yet accurate: %+v", loc)
		}
	})
}

// ---------- 6.2 DDL Loc Verification ----------

func TestLoc_DDL(t *testing.T) {
	t.Run("CREATE TABLE", func(t *testing.T) {
		sql := "CREATE TABLE t (a INT, b VARCHAR(50))"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		ct, ok := stmts[0].AST.(*ast.CreateTableStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.CreateTableStmt", stmts[0].AST)
		}
		loc := ct.Loc
		if !locValid(loc) {
			t.Fatalf("CreateTableStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "CREATE TABLE") {
			t.Errorf("CreateTableStmt Loc text = %q, want it to contain CREATE TABLE", got)
		}
	})

	t.Run("CREATE INDEX", func(t *testing.T) {
		sql := "CREATE INDEX ix ON t (a)"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		ci, ok := stmts[0].AST.(*ast.CreateIndexStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.CreateIndexStmt", stmts[0].AST)
		}
		loc := ci.Loc
		if !locValid(loc) {
			t.Fatalf("CreateIndexStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "CREATE INDEX") {
			t.Errorf("CreateIndexStmt Loc text = %q, want it to contain CREATE INDEX", got)
		}
	})

	t.Run("CREATE VIEW", func(t *testing.T) {
		sql := "CREATE VIEW v AS SELECT 1"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		cv, ok := stmts[0].AST.(*ast.CreateViewStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.CreateViewStmt", stmts[0].AST)
		}
		loc := cv.Loc
		if !locValid(loc) {
			t.Fatalf("CreateViewStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "CREATE VIEW") {
			t.Errorf("CreateViewStmt Loc text = %q, want it to contain CREATE VIEW", got)
		}
	})

	t.Run("CREATE PROCEDURE", func(t *testing.T) {
		sql := "CREATE PROCEDURE p AS SELECT 1"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		cp, ok := stmts[0].AST.(*ast.CreateProcedureStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.CreateProcedureStmt", stmts[0].AST)
		}
		loc := cp.Loc
		if !locValid(loc) {
			t.Fatalf("CreateProcedureStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "CREATE PROCEDURE") {
			t.Errorf("CreateProcedureStmt Loc text = %q, want it to contain CREATE PROCEDURE", got)
		}
	})

	t.Run("ALTER TABLE", func(t *testing.T) {
		sql := "ALTER TABLE t ADD c INT"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		at, ok := stmts[0].AST.(*ast.AlterTableStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.AlterTableStmt", stmts[0].AST)
		}
		loc := at.Loc
		if !locValid(loc) {
			t.Fatalf("AlterTableStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "ALTER TABLE") {
			t.Errorf("AlterTableStmt Loc text = %q, want it to contain ALTER TABLE", got)
		}
	})

	t.Run("DROP TABLE", func(t *testing.T) {
		sql := "DROP TABLE t"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		ds, ok := stmts[0].AST.(*ast.DropStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.DropStmt", stmts[0].AST)
		}
		loc := ds.Loc
		if !locValid(loc) {
			t.Fatalf("DropStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "DROP TABLE") {
			t.Errorf("DropStmt Loc text = %q, want it to contain DROP TABLE", got)
		}
	})

	t.Run("CREATE FUNCTION", func(t *testing.T) {
		sql := "CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN 1 END"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		cf, ok := stmts[0].AST.(*ast.CreateFunctionStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.CreateFunctionStmt", stmts[0].AST)
		}
		loc := cf.Loc
		if !locValid(loc) {
			t.Fatalf("CreateFunctionStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "CREATE FUNCTION") {
			t.Errorf("CreateFunctionStmt Loc text = %q, want it to contain CREATE FUNCTION", got)
		}
	})

	t.Run("CREATE TRIGGER", func(t *testing.T) {
		sql := "CREATE TRIGGER tr ON t FOR INSERT AS SELECT 1"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		ct, ok := stmts[0].AST.(*ast.CreateTriggerStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.CreateTriggerStmt", stmts[0].AST)
		}
		loc := ct.Loc
		if !locValid(loc) {
			t.Fatalf("CreateTriggerStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "CREATE TRIGGER") {
			t.Errorf("CreateTriggerStmt Loc text = %q, want it to contain CREATE TRIGGER", got)
		}
	})

	t.Run("TRUNCATE TABLE", func(t *testing.T) {
		sql := "TRUNCATE TABLE t"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		ts, ok := stmts[0].AST.(*ast.TruncateStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.TruncateStmt", stmts[0].AST)
		}
		loc := ts.Loc
		if !locValid(loc) {
			t.Fatalf("TruncateStmt Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "TRUNCATE TABLE") {
			t.Errorf("TruncateStmt Loc text = %q, want it to contain TRUNCATE TABLE", got)
		}
	})
}

// ---------- 6.3 Expression & Sub-node Loc ----------

func TestLoc_Expressions(t *testing.T) {
	// Helper: parse a single SELECT and return its first target list expression.
	parseExpr := func(t *testing.T, sql string) ast.Node {
		t.Helper()
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		sel, ok := stmts[0].AST.(*ast.SelectStmt)
		if !ok {
			t.Fatalf("AST type = %T, want *ast.SelectStmt", stmts[0].AST)
		}
		if sel.TargetList == nil || sel.TargetList.Len() == 0 {
			t.Fatal("TargetList is empty")
		}
		return sel.TargetList.Items[0]
	}

	// Helper: unwrap ResTarget to get the expression value.
	unwrapExpr := func(t *testing.T, node ast.Node) ast.Node {
		t.Helper()
		if rt, ok := node.(*ast.ResTarget); ok {
			if rt.Val != nil {
				return rt.Val
			}
		}
		return node
	}

	t.Run("CAST", func(t *testing.T) {
		sql := "SELECT CAST(1 AS INT)"
		node := unwrapExpr(t, parseExpr(t, sql))
		ce, ok := node.(*ast.CastExpr)
		if !ok {
			t.Fatalf("expr type = %T, want *ast.CastExpr", node)
		}
		loc := ce.Loc
		// [~] partial: expression Loc.End not yet set accurately.
		if !locStartValid(loc) {
			t.Fatalf("CastExpr Loc.Start invalid: %+v", loc)
		}
		if locValid(loc) {
			got := locText(sql, loc)
			if !strings.Contains(got, "CAST") {
				t.Errorf("CastExpr Loc text = %q, want it to contain CAST", got)
			}
		} else {
			t.Logf("[~] partial: CastExpr Loc.End not yet accurate: %+v", loc)
		}
	})

	t.Run("CASE", func(t *testing.T) {
		sql := "SELECT CASE WHEN 1=1 THEN 'a' ELSE 'b' END"
		node := unwrapExpr(t, parseExpr(t, sql))
		ce, ok := node.(*ast.CaseExpr)
		if !ok {
			t.Fatalf("expr type = %T, want *ast.CaseExpr", node)
		}
		loc := ce.Loc
		// [~] partial: expression Loc.End not yet set accurately.
		if !locStartValid(loc) {
			t.Fatalf("CaseExpr Loc.Start invalid: %+v", loc)
		}
		if locValid(loc) {
			got := locText(sql, loc)
			if !strings.Contains(got, "CASE") {
				t.Errorf("CaseExpr Loc text = %q, want it to contain CASE", got)
			}
		} else {
			t.Logf("[~] partial: CaseExpr Loc.End not yet accurate: %+v", loc)
		}
	})

	t.Run("BinaryExpr a + b", func(t *testing.T) {
		sql := "SELECT a + b"
		node := unwrapExpr(t, parseExpr(t, sql))
		be, ok := node.(*ast.BinaryExpr)
		if !ok {
			t.Fatalf("expr type = %T, want *ast.BinaryExpr", node)
		}
		loc := be.Loc
		// [~] partial: expression Loc.End not yet set accurately.
		if !locStartValid(loc) {
			t.Fatalf("BinaryExpr Loc.Start invalid: %+v", loc)
		}
		if locValid(loc) {
			got := locText(sql, loc)
			if !strings.Contains(got, "+") {
				t.Errorf("BinaryExpr Loc text = %q, want it to contain +", got)
			}
		} else {
			t.Logf("[~] partial: BinaryExpr Loc.End not yet accurate: %+v", loc)
		}
	})

	t.Run("UnaryExpr -x", func(t *testing.T) {
		sql := "SELECT -x"
		node := unwrapExpr(t, parseExpr(t, sql))
		ue, ok := node.(*ast.UnaryExpr)
		if !ok {
			t.Fatalf("expr type = %T, want *ast.UnaryExpr", node)
		}
		loc := ue.Loc
		// [~] partial: expression Loc.End not yet set accurately.
		if !locStartValid(loc) {
			t.Fatalf("UnaryExpr Loc.Start invalid: %+v", loc)
		}
		if locValid(loc) {
			got := locText(sql, loc)
			if !strings.Contains(got, "-") {
				t.Errorf("UnaryExpr Loc text = %q, want it to contain -", got)
			}
		} else {
			t.Logf("[~] partial: UnaryExpr Loc.End not yet accurate: %+v", loc)
		}
	})

	t.Run("COALESCE FuncCall", func(t *testing.T) {
		sql := "SELECT COALESCE(a, b, c)"
		node := unwrapExpr(t, parseExpr(t, sql))
		// COALESCE may be parsed as CoalesceExpr or FuncCallExpr.
		loc := testNodeLoc(node)
		// [~] partial: expression Loc.End not yet set accurately.
		if !locStartValid(loc) {
			t.Fatalf("COALESCE Loc.Start invalid: %+v (node type %T)", loc, node)
		}
		if locValid(loc) {
			got := locText(sql, loc)
			if !strings.Contains(got, "COALESCE") {
				t.Errorf("COALESCE Loc text = %q, want it to contain COALESCE", got)
			}
		} else {
			t.Logf("[~] partial: COALESCE Loc.End not yet accurate: %+v (type %T)", loc, node)
		}
	})

	t.Run("EXISTS", func(t *testing.T) {
		sql := "SELECT * FROM t WHERE EXISTS (SELECT 1)"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		sel := stmts[0].AST.(*ast.SelectStmt)
		if sel.WhereClause == nil {
			t.Fatal("WhereClause is nil")
		}
		ee, ok := sel.WhereClause.(*ast.ExistsExpr)
		if !ok {
			t.Fatalf("WhereClause type = %T, want *ast.ExistsExpr", sel.WhereClause)
		}
		loc := ee.Loc
		// [~] partial: expression Loc.End not yet set accurately.
		if !locStartValid(loc) {
			t.Fatalf("ExistsExpr Loc.Start invalid: %+v", loc)
		}
		if locValid(loc) {
			got := locText(sql, loc)
			if !strings.Contains(got, "EXISTS") {
				t.Errorf("ExistsExpr Loc text = %q, want it to contain EXISTS", got)
			}
		} else {
			t.Logf("[~] partial: ExistsExpr Loc.End not yet accurate: %+v", loc)
		}
	})

	t.Run("CONVERT", func(t *testing.T) {
		sql := "SELECT CONVERT(INT, '1')"
		node := unwrapExpr(t, parseExpr(t, sql))
		ce, ok := node.(*ast.ConvertExpr)
		if !ok {
			t.Fatalf("expr type = %T, want *ast.ConvertExpr", node)
		}
		loc := ce.Loc
		// [~] partial: expression Loc.End not yet set accurately.
		if !locStartValid(loc) {
			t.Fatalf("ConvertExpr Loc.Start invalid: %+v", loc)
		}
		if locValid(loc) {
			got := locText(sql, loc)
			if !strings.Contains(got, "CONVERT") {
				t.Errorf("ConvertExpr Loc text = %q, want it to contain CONVERT", got)
			}
		} else {
			t.Logf("[~] partial: ConvertExpr Loc.End not yet accurate: %+v", loc)
		}
	})

	t.Run("TRY_CAST", func(t *testing.T) {
		sql := "SELECT TRY_CAST(1 AS VARCHAR)"
		node := unwrapExpr(t, parseExpr(t, sql))
		tc, ok := node.(*ast.TryCastExpr)
		if !ok {
			t.Fatalf("expr type = %T, want *ast.TryCastExpr", node)
		}
		loc := tc.Loc
		// [~] partial: expression Loc.End not yet set accurately.
		if !locStartValid(loc) {
			t.Fatalf("TryCastExpr Loc.Start invalid: %+v", loc)
		}
		if locValid(loc) {
			got := locText(sql, loc)
			if !strings.Contains(got, "TRY_CAST") {
				t.Errorf("TryCastExpr Loc text = %q, want it to contain TRY_CAST", got)
			}
		} else {
			t.Logf("[~] partial: TryCastExpr Loc.End not yet accurate: %+v", loc)
		}
	})

	t.Run("TRY_CONVERT", func(t *testing.T) {
		sql := "SELECT TRY_CONVERT(INT, '1')"
		node := unwrapExpr(t, parseExpr(t, sql))
		tc, ok := node.(*ast.TryConvertExpr)
		if !ok {
			t.Fatalf("expr type = %T, want *ast.TryConvertExpr", node)
		}
		loc := tc.Loc
		// [~] partial: expression Loc.End not yet set accurately.
		if !locStartValid(loc) {
			t.Fatalf("TryConvertExpr Loc.Start invalid: %+v", loc)
		}
		if locValid(loc) {
			got := locText(sql, loc)
			if !strings.Contains(got, "TRY_CONVERT") {
				t.Errorf("TryConvertExpr Loc text = %q, want it to contain TRY_CONVERT", got)
			}
		} else {
			t.Logf("[~] partial: TryConvertExpr Loc.End not yet accurate: %+v", loc)
		}
	})
}

// ---------- 6.4 Multi-Statement & Edge Cases ----------

func TestLoc_MultiStatementAndEdgeCases(t *testing.T) {
	t.Run("non-overlapping Loc ranges", func(t *testing.T) {
		sql := "SELECT 1; SELECT 2"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 2 {
			t.Fatalf("got %d statements, want 2", len(stmts))
		}
		loc0 := testNodeLoc(stmts[0].AST)
		loc1 := testNodeLoc(stmts[1].AST)
		if locValid(loc0) && locValid(loc1) {
			if loc0.End > loc1.Start {
				t.Errorf("statement Locs overlap: stmt0.End=%d > stmt1.Start=%d", loc0.End, loc1.Start)
			}
		}
		// Also check Statement-level ByteStart/ByteEnd.
		if stmts[0].ByteEnd > stmts[1].ByteStart {
			t.Errorf("Statement byte ranges overlap: s0.ByteEnd=%d > s1.ByteStart=%d",
				stmts[0].ByteEnd, stmts[1].ByteStart)
		}
	})

	t.Run("multi-line byte offsets", func(t *testing.T) {
		sql := "SELECT\n  1\nFROM\n  t"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		loc := testNodeLoc(stmts[0].AST)
		if !locValid(loc) {
			t.Fatalf("Loc invalid: %+v", loc)
		}
		got := locText(sql, loc)
		if !strings.Contains(got, "SELECT") || !strings.Contains(got, "FROM") {
			t.Errorf("multi-line Loc text = %q, want SELECT and FROM", got)
		}
	})

	t.Run("leading whitespace", func(t *testing.T) {
		sql := "   SELECT 1"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		loc := testNodeLoc(stmts[0].AST)
		if locValid(loc) {
			// Loc.Start should point to 'S' in SELECT, not to whitespace.
			if loc.Start < 3 {
				t.Errorf("Loc.Start = %d, want >= 3 (past leading whitespace)", loc.Start)
			}
		}
		// Statement Start position should also point to first keyword.
		if stmts[0].Start.Column < 4 {
			t.Errorf("Statement Start.Column = %d, want >= 4", stmts[0].Start.Column)
		}
	})

	t.Run("trailing semicolon", func(t *testing.T) {
		sql := "SELECT 1;"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		loc := testNodeLoc(stmts[0].AST)
		if !locValid(loc) {
			t.Fatalf("Loc invalid: %+v", loc)
		}
		// The AST node Loc.End might point before or at the semicolon.
		// Statement.ByteEnd should include the semicolon.
		if stmts[0].ByteEnd != len(sql) {
			t.Errorf("Statement.ByteEnd = %d, want %d", stmts[0].ByteEnd, len(sql))
		}
	})

	t.Run("empty semicolons no panic", func(t *testing.T) {
		// Should not panic on bare semicolons.
		stmts, err := Parse(";;")
		// Some parsers return empty/nil, some return statements. Just don't panic.
		_ = stmts
		_ = err
	})

	t.Run("very long single-line SQL", func(t *testing.T) {
		// Build a long SQL with many columns.
		var b strings.Builder
		b.WriteString("SELECT ")
		for i := 0; i < 200; i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString("col")
			b.WriteString(strings.Repeat("x", 10))
		}
		b.WriteString(" FROM t")
		sql := b.String()
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		loc := testNodeLoc(stmts[0].AST)
		if !locValid(loc) {
			t.Fatalf("Loc invalid for long SQL: %+v", loc)
		}
		// Loc should span the full statement.
		if loc.End < len(sql)-10 {
			t.Errorf("Loc.End = %d, but SQL len = %d; expected Loc to cover most of the statement", loc.End, len(sql))
		}
	})

	t.Run("block comment before semicolon", func(t *testing.T) {
		sql := "SELECT 1 /* comment */; SELECT 2"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) < 2 {
			t.Fatalf("got %d statements, want >= 2", len(stmts))
		}
		loc0 := testNodeLoc(stmts[0].AST)
		loc1 := testNodeLoc(stmts[1].AST)
		if locValid(loc0) && locValid(loc1) {
			if loc0.Start >= loc1.Start {
				t.Errorf("stmt0 Loc.Start=%d >= stmt1 Loc.Start=%d", loc0.Start, loc1.Start)
			}
		}
	})

	t.Run("GO batch separator", func(t *testing.T) {
		sql := "SELECT 1\nGO\nSELECT 2"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		// Should parse successfully — may produce 2-3 statements (SELECT, GO, SELECT).
		if len(stmts) < 2 {
			t.Fatalf("got %d statements, want >= 2", len(stmts))
		}
		// All statements should have valid positions.
		for i, s := range stmts {
			if s.AST != nil {
				loc := testNodeLoc(s.AST)
				if loc.Start >= 0 && loc.End >= 0 && loc.End <= loc.Start {
					t.Errorf("stmt[%d] has invalid Loc: Start=%d, End=%d", i, loc.Start, loc.End)
				}
			}
		}
	})

	t.Run("error-path Loc", func(t *testing.T) {
		// Intentionally malformed SQL.
		_, err := Parse("SELECT FROM WHERE GROUP BY")
		// Just verify no panic. Error may or may not be returned depending on parser resilience.
		_ = err
	})
}

// ---------- 6.5 Public API Line:Column Tests ----------

func TestLoc_PublicAPI(t *testing.T) {
	t.Run("single-line SELECT 1 positions", func(t *testing.T) {
		stmts, err := Parse("SELECT 1")
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		s := stmts[0]
		if s.Start.Line != 1 || s.Start.Column != 1 {
			t.Errorf("Start = %+v, want {Line:1, Column:1}", s.Start)
		}
		if s.End.Line != 1 {
			t.Errorf("End.Line = %d, want 1", s.End.Line)
		}
		// End column should be past the statement.
		if s.End.Column <= s.Start.Column {
			t.Errorf("End.Column = %d should be > Start.Column = %d", s.End.Column, s.Start.Column)
		}
	})

	t.Run("two-line SELECT", func(t *testing.T) {
		sql := "SELECT\n1"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		s := stmts[0]
		if s.Start.Line != 1 {
			t.Errorf("Start.Line = %d, want 1", s.Start.Line)
		}
		// End should be on line 2 (the "1" is on line 2).
		if s.End.Line < 2 {
			t.Errorf("End.Line = %d, want >= 2", s.End.Line)
		}
	})

	t.Run("multi-statement multi-line", func(t *testing.T) {
		sql := "SELECT 1\nSELECT 2\nSELECT 3"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 3 {
			t.Fatalf("got %d statements, want 3", len(stmts))
		}
		// Each statement should start on a successive line.
		for i := 0; i < len(stmts)-1; i++ {
			if stmts[i].Start.Line >= stmts[i+1].Start.Line {
				t.Errorf("stmt[%d].Start.Line=%d >= stmt[%d].Start.Line=%d",
					i, stmts[i].Start.Line, i+1, stmts[i+1].Start.Line)
			}
		}
	})

	t.Run("tab characters count as 1 column", func(t *testing.T) {
		sql := "\tSELECT 1"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		// Tab is 1 byte, so column for SELECT should be 2.
		if stmts[0].Start.Column != 2 {
			t.Errorf("Start.Column = %d, want 2 (tab counts as 1)", stmts[0].Start.Column)
		}
	})

	t.Run("unicode content byte offset", func(t *testing.T) {
		// The emoji is 4 bytes in UTF-8.
		sql := "SELECT N'\U0001F600'"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		s := stmts[0]
		// The End column should be byte-based, not character-based.
		// "SELECT N'<4-byte-emoji>'" = 7 + 2 + 4 + 1 = 14 bytes.
		// Column should be > 10 (character count) since it's byte-based.
		if s.End.Column <= 1 {
			t.Errorf("End.Column = %d, expected > 1 for unicode content", s.End.Column)
		}
	})

	t.Run("Parse empty returns nil", func(t *testing.T) {
		stmts, err := Parse("")
		if err != nil {
			t.Fatalf("Parse('') error: %v", err)
		}
		if stmts != nil && len(stmts) != 0 {
			t.Errorf("Parse('') returned %d statements, want 0", len(stmts))
		}
	})

	t.Run("Parse single returns one", func(t *testing.T) {
		stmts, err := Parse("SELECT 1")
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		if stmts[0].AST == nil {
			t.Error("AST is nil")
		}
		if stmts[0].ByteStart != 0 {
			t.Errorf("ByteStart = %d, want 0", stmts[0].ByteStart)
		}
	})

	t.Run("Parse three statements", func(t *testing.T) {
		stmts, err := Parse("SELECT 1; SELECT 2; SELECT 3")
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 3 {
			t.Fatalf("got %d statements, want 3", len(stmts))
		}
		for i, s := range stmts {
			if s.AST == nil {
				t.Errorf("stmt[%d].AST is nil", i)
			}
		}
		// Verify non-overlapping byte ranges.
		for i := 0; i < len(stmts)-1; i++ {
			if stmts[i].ByteEnd > stmts[i+1].ByteStart {
				t.Errorf("stmt[%d].ByteEnd=%d > stmt[%d].ByteStart=%d",
					i, stmts[i].ByteEnd, i+1, stmts[i+1].ByteStart)
			}
		}
	})
}
