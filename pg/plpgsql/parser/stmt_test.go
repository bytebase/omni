package parser

import (
	"testing"

	"github.com/bytebase/omni/pg/plpgsql/ast"
)

// --------------------------------------------------------------------------
// Section 3.1: Variable Assignment
// --------------------------------------------------------------------------

func TestAssignSimple(t *testing.T) {
	block := parseOK(t, `BEGIN x := 1; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	assign, ok := block.Body[0].(*ast.PLAssign)
	if !ok {
		t.Fatalf("expected PLAssign, got %T", block.Body[0])
	}
	if assign.Target != "x" {
		t.Errorf("target = %q, want %q", assign.Target, "x")
	}
	if assign.Expr != "1" {
		t.Errorf("expr = %q, want %q", assign.Expr, "1")
	}
}

func TestAssignExpression(t *testing.T) {
	block := parseOK(t, `BEGIN x := a + b * c; END`)
	assign := block.Body[0].(*ast.PLAssign)
	if assign.Target != "x" {
		t.Errorf("target = %q, want %q", assign.Target, "x")
	}
	if assign.Expr != "a + b * c" {
		t.Errorf("expr = %q, want %q", assign.Expr, "a + b * c")
	}
}

func TestAssignFunctionCall(t *testing.T) {
	block := parseOK(t, `BEGIN x := my_func(a, b); END`)
	assign := block.Body[0].(*ast.PLAssign)
	if assign.Expr != "my_func(a, b)" {
		t.Errorf("expr = %q, want %q", assign.Expr, "my_func(a, b)")
	}
}

func TestAssignSubquery(t *testing.T) {
	block := parseOK(t, `BEGIN x := (SELECT max(a) FROM t); END`)
	assign := block.Body[0].(*ast.PLAssign)
	if assign.Expr != "(SELECT max(a) FROM t)" {
		t.Errorf("expr = %q, want %q", assign.Expr, "(SELECT max(a) FROM t)")
	}
}

func TestAssignRecordField(t *testing.T) {
	block := parseOK(t, `BEGIN rec.field := 42; END`)
	assign := block.Body[0].(*ast.PLAssign)
	if assign.Target != "rec.field" {
		t.Errorf("target = %q, want %q", assign.Target, "rec.field")
	}
	if assign.Expr != "42" {
		t.Errorf("expr = %q, want %q", assign.Expr, "42")
	}
}

func TestAssignArrayElement(t *testing.T) {
	block := parseOK(t, `BEGIN arr[1] := 'hello'; END`)
	assign := block.Body[0].(*ast.PLAssign)
	if assign.Target != "arr[1]" {
		t.Errorf("target = %q, want %q", assign.Target, "arr[1]")
	}
	if assign.Expr != "'hello'" {
		t.Errorf("expr = %q, want %q", assign.Expr, "'hello'")
	}
}

func TestAssignArraySlice(t *testing.T) {
	block := parseOK(t, `BEGIN arr[1:3] := ARRAY[1,2,3]; END`)
	assign := block.Body[0].(*ast.PLAssign)
	if assign.Target != "arr[1:3]" {
		t.Errorf("target = %q, want %q", assign.Target, "arr[1:3]")
	}
	if assign.Expr != "ARRAY[1,2,3]" {
		t.Errorf("expr = %q, want %q", assign.Expr, "ARRAY[1,2,3]")
	}
}

func TestAssignMultiLevelField(t *testing.T) {
	block := parseOK(t, `BEGIN rec.nested.field := 1; END`)
	assign := block.Body[0].(*ast.PLAssign)
	if assign.Target != "rec.nested.field" {
		t.Errorf("target = %q, want %q", assign.Target, "rec.nested.field")
	}
}

func TestAssignComplexRHS(t *testing.T) {
	block := parseOK(t, `BEGIN x := (a + b) * (c - d) / e; END`)
	assign := block.Body[0].(*ast.PLAssign)
	if assign.Expr != "(a + b) * (c - d) / e" {
		t.Errorf("expr = %q, want %q", assign.Expr, "(a + b) * (c - d) / e")
	}
}

func TestAssignEqualsOperator(t *testing.T) {
	block := parseOK(t, `BEGIN x = 1; END`)
	assign := block.Body[0].(*ast.PLAssign)
	if assign.Target != "x" {
		t.Errorf("target = %q, want %q", assign.Target, "x")
	}
	if assign.Expr != "1" {
		t.Errorf("expr = %q, want %q", assign.Expr, "1")
	}
}

// --------------------------------------------------------------------------
// Section 3.2: RETURN Variants
// --------------------------------------------------------------------------

func TestReturnBare(t *testing.T) {
	block := parseOK(t, `BEGIN RETURN; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	ret, ok := block.Body[0].(*ast.PLReturn)
	if !ok {
		t.Fatalf("expected PLReturn, got %T", block.Body[0])
	}
	if ret.Expr != "" {
		t.Errorf("expr = %q, want empty", ret.Expr)
	}
}

func TestReturnWithExpr(t *testing.T) {
	block := parseOK(t, `BEGIN RETURN x + 1; END`)
	ret := block.Body[0].(*ast.PLReturn)
	if ret.Expr != "x + 1" {
		t.Errorf("expr = %q, want %q", ret.Expr, "x + 1")
	}
}

func TestReturnWithSubquery(t *testing.T) {
	block := parseOK(t, `BEGIN RETURN (SELECT count(*) FROM t); END`)
	ret := block.Body[0].(*ast.PLReturn)
	if ret.Expr != "(SELECT count(*) FROM t)" {
		t.Errorf("expr = %q, want %q", ret.Expr, "(SELECT count(*) FROM t)")
	}
}

func TestReturnNext(t *testing.T) {
	block := parseOK(t, `BEGIN RETURN NEXT x; END`)
	ret, ok := block.Body[0].(*ast.PLReturnNext)
	if !ok {
		t.Fatalf("expected PLReturnNext, got %T", block.Body[0])
	}
	if ret.Expr != "x" {
		t.Errorf("expr = %q, want %q", ret.Expr, "x")
	}
}

func TestReturnNextRecord(t *testing.T) {
	block := parseOK(t, `BEGIN RETURN NEXT rec; END`)
	ret := block.Body[0].(*ast.PLReturnNext)
	if ret.Expr != "rec" {
		t.Errorf("expr = %q, want %q", ret.Expr, "rec")
	}
}

func TestReturnNextBare(t *testing.T) {
	block := parseOK(t, `BEGIN RETURN NEXT; END`)
	ret := block.Body[0].(*ast.PLReturnNext)
	if ret.Expr != "" {
		t.Errorf("expr = %q, want empty", ret.Expr)
	}
}

func TestReturnQueryStatic(t *testing.T) {
	block := parseOK(t, `BEGIN RETURN QUERY SELECT * FROM t; END`)
	ret, ok := block.Body[0].(*ast.PLReturnQuery)
	if !ok {
		t.Fatalf("expected PLReturnQuery, got %T", block.Body[0])
	}
	if ret.Query != "SELECT * FROM t" {
		t.Errorf("query = %q, want %q", ret.Query, "SELECT * FROM t")
	}
	if ret.DynQuery != "" {
		t.Errorf("dynQuery = %q, want empty", ret.DynQuery)
	}
}

func TestReturnQueryWithWhere(t *testing.T) {
	block := parseOK(t, `BEGIN RETURN QUERY SELECT a, b FROM t WHERE c > 0; END`)
	ret := block.Body[0].(*ast.PLReturnQuery)
	if ret.Query != "SELECT a, b FROM t WHERE c > 0" {
		t.Errorf("query = %q, want %q", ret.Query, "SELECT a, b FROM t WHERE c > 0")
	}
}

func TestReturnQueryExecute(t *testing.T) {
	block := parseOK(t, `BEGIN RETURN QUERY EXECUTE 'SELECT * FROM t'; END`)
	ret := block.Body[0].(*ast.PLReturnQuery)
	if ret.DynQuery != "'SELECT * FROM t'" {
		t.Errorf("dynQuery = %q, want %q", ret.DynQuery, "'SELECT * FROM t'")
	}
	if ret.Query != "" {
		t.Errorf("query = %q, want empty", ret.Query)
	}
}

func TestReturnQueryExecuteUsing(t *testing.T) {
	block := parseOK(t, `BEGIN RETURN QUERY EXECUTE 'SELECT * FROM t WHERE id = $1' USING my_id; END`)
	ret := block.Body[0].(*ast.PLReturnQuery)
	if ret.DynQuery != "'SELECT * FROM t WHERE id = $1'" {
		t.Errorf("dynQuery = %q, want %q", ret.DynQuery, "'SELECT * FROM t WHERE id = $1'")
	}
	if len(ret.Params) != 1 || ret.Params[0] != "my_id" {
		t.Errorf("params = %v, want [my_id]", ret.Params)
	}
}

func TestReturnQueryExecuteMultipleUsing(t *testing.T) {
	block := parseOK(t, `BEGIN RETURN QUERY EXECUTE $q$ SELECT $1, $2 $q$ USING a, b; END`)
	ret := block.Body[0].(*ast.PLReturnQuery)
	if len(ret.Params) != 2 {
		t.Fatalf("params len = %d, want 2", len(ret.Params))
	}
	if ret.Params[0] != "a" {
		t.Errorf("params[0] = %q, want %q", ret.Params[0], "a")
	}
	if ret.Params[1] != "b" {
		t.Errorf("params[1] = %q, want %q", ret.Params[1], "b")
	}
}

func TestReturnExprSpansUntilSemicolon(t *testing.T) {
	block := parseOK(t, `BEGIN RETURN a + b * c - d; END`)
	ret := block.Body[0].(*ast.PLReturn)
	if ret.Expr != "a + b * c - d" {
		t.Errorf("expr = %q, want %q", ret.Expr, "a + b * c - d")
	}
}

// --------------------------------------------------------------------------
// Section 3.3: PERFORM and Bare SQL
// --------------------------------------------------------------------------

func TestPerformCall(t *testing.T) {
	block := parseOK(t, `BEGIN PERFORM my_func(1, 2); END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	perf, ok := block.Body[0].(*ast.PLPerform)
	if !ok {
		t.Fatalf("expected PLPerform, got %T", block.Body[0])
	}
	if perf.Expr != "my_func(1, 2)" {
		t.Errorf("expr = %q, want %q", perf.Expr, "my_func(1, 2)")
	}
}

func TestPerformQuery(t *testing.T) {
	block := parseOK(t, `BEGIN PERFORM * FROM t WHERE a > 0; END`)
	perf := block.Body[0].(*ast.PLPerform)
	if perf.Expr != "* FROM t WHERE a > 0" {
		t.Errorf("expr = %q, want %q", perf.Expr, "* FROM t WHERE a > 0")
	}
}

func TestBareSQLInsert(t *testing.T) {
	block := parseOK(t, `BEGIN INSERT INTO t VALUES (1, 2); END`)
	exec, ok := block.Body[0].(*ast.PLExecSQL)
	if !ok {
		t.Fatalf("expected PLExecSQL, got %T", block.Body[0])
	}
	if exec.SQLText != "INSERT INTO t VALUES (1, 2)" {
		t.Errorf("sql = %q, want %q", exec.SQLText, "INSERT INTO t VALUES (1, 2)")
	}
}

func TestBareSQLUpdate(t *testing.T) {
	block := parseOK(t, `BEGIN UPDATE t SET a = 1 WHERE b = 2; END`)
	exec := block.Body[0].(*ast.PLExecSQL)
	if exec.SQLText != "UPDATE t SET a = 1 WHERE b = 2" {
		t.Errorf("sql = %q, want %q", exec.SQLText, "UPDATE t SET a = 1 WHERE b = 2")
	}
}

func TestBareSQLDelete(t *testing.T) {
	block := parseOK(t, `BEGIN DELETE FROM t WHERE a = 1; END`)
	exec := block.Body[0].(*ast.PLExecSQL)
	if exec.SQLText != "DELETE FROM t WHERE a = 1" {
		t.Errorf("sql = %q, want %q", exec.SQLText, "DELETE FROM t WHERE a = 1")
	}
}

func TestBareSQLWithSelectInto(t *testing.T) {
	block := parseOK(t, `BEGIN WITH rows AS (SELECT 1 AS id) SELECT id INTO result FROM rows; END`)
	exec, ok := block.Body[0].(*ast.PLExecSQL)
	if !ok {
		t.Fatalf("expected PLExecSQL, got %T", block.Body[0])
	}
	if len(exec.Into) != 1 || exec.Into[0] != "result" {
		t.Errorf("into = %v, want [result]", exec.Into)
	}
	if exec.SQLText != "WITH rows AS (SELECT 1 AS id) SELECT id FROM rows" {
		t.Errorf("sql = %q, want %q", exec.SQLText, "WITH rows AS (SELECT 1 AS id) SELECT id FROM rows")
	}
}

func TestBareSQLWithInsertDoesNotExtractInsertInto(t *testing.T) {
	block := parseOK(t, `BEGIN WITH rows AS (SELECT 1 AS id) INSERT INTO t SELECT id FROM rows; END`)
	exec, ok := block.Body[0].(*ast.PLExecSQL)
	if !ok {
		t.Fatalf("expected PLExecSQL, got %T", block.Body[0])
	}
	if len(exec.Into) != 0 {
		t.Errorf("into = %v, want empty", exec.Into)
	}
	if exec.SQLText != "WITH rows AS (SELECT 1 AS id) INSERT INTO t SELECT id FROM rows" {
		t.Errorf("sql = %q, want %q", exec.SQLText, "WITH rows AS (SELECT 1 AS id) INSERT INTO t SELECT id FROM rows")
	}
}

func TestBareSQLInsertReturningInto(t *testing.T) {
	block := parseOK(t, `BEGIN INSERT INTO t (a) VALUES (1) RETURNING id INTO new_id; END`)
	exec, ok := block.Body[0].(*ast.PLExecSQL)
	if !ok {
		t.Fatalf("expected PLExecSQL, got %T", block.Body[0])
	}
	if len(exec.Into) != 1 || exec.Into[0] != "new_id" {
		t.Errorf("into = %v, want [new_id]", exec.Into)
	}
	if exec.SQLText != "INSERT INTO t (a) VALUES (1) RETURNING id" {
		t.Errorf("sql = %q, want %q", exec.SQLText, "INSERT INTO t (a) VALUES (1) RETURNING id")
	}
}

func TestBareSQLUpdateReturningIntoStrict(t *testing.T) {
	block := parseOK(t, `BEGIN UPDATE t SET a = 1 WHERE id = 2 RETURNING a INTO STRICT updated_a; END`)
	exec, ok := block.Body[0].(*ast.PLExecSQL)
	if !ok {
		t.Fatalf("expected PLExecSQL, got %T", block.Body[0])
	}
	if len(exec.Into) != 1 || exec.Into[0] != "updated_a" {
		t.Errorf("into = %v, want [updated_a]", exec.Into)
	}
	if !exec.Strict {
		t.Errorf("strict = false, want true")
	}
	if exec.SQLText != "UPDATE t SET a = 1 WHERE id = 2 RETURNING a" {
		t.Errorf("sql = %q, want %q", exec.SQLText, "UPDATE t SET a = 1 WHERE id = 2 RETURNING a")
	}
}

func TestBareSQLUtilityCommands(t *testing.T) {
	tests := []struct {
		body string
		want string
	}{
		{
			body: `BEGIN REFRESH MATERIALIZED VIEW CONCURRENTLY mv; END`,
			want: "REFRESH MATERIALIZED VIEW CONCURRENTLY mv",
		},
		{
			body: `BEGIN GRANT USAGE ON SCHEMA api TO web_anon; END`,
			want: "GRANT USAGE ON SCHEMA api TO web_anon",
		},
		{
			body: `BEGIN CREATE OR REPLACE FUNCTION helper() RETURNS void LANGUAGE sql AS 'SELECT 1'; END`,
			want: "CREATE OR REPLACE FUNCTION helper() RETURNS void LANGUAGE sql AS 'SELECT 1'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			block := parseOK(t, tt.body)
			exec, ok := block.Body[0].(*ast.PLExecSQL)
			if !ok {
				t.Fatalf("expected PLExecSQL, got %T", block.Body[0])
			}
			if exec.SQLText != tt.want {
				t.Errorf("sql = %q, want %q", exec.SQLText, tt.want)
			}
		})
	}
}

func TestSelectInto(t *testing.T) {
	block := parseOK(t, `BEGIN SELECT a, b INTO x, y FROM t WHERE id = 1; END`)
	exec, ok := block.Body[0].(*ast.PLExecSQL)
	if !ok {
		t.Fatalf("expected PLExecSQL, got %T", block.Body[0])
	}
	if len(exec.Into) != 2 || exec.Into[0] != "x" || exec.Into[1] != "y" {
		t.Errorf("into = %v, want [x y]", exec.Into)
	}
	if exec.Strict {
		t.Errorf("strict = true, want false")
	}
	// The SQL text should have INTO clause removed
	if exec.SQLText != "SELECT a, b FROM t WHERE id = 1" {
		t.Errorf("sql = %q, want %q", exec.SQLText, "SELECT a, b FROM t WHERE id = 1")
	}
}

func TestSelectIntoStrict(t *testing.T) {
	block := parseOK(t, `BEGIN SELECT a INTO STRICT x FROM t WHERE id = 1; END`)
	exec := block.Body[0].(*ast.PLExecSQL)
	if len(exec.Into) != 1 || exec.Into[0] != "x" {
		t.Errorf("into = %v, want [x]", exec.Into)
	}
	if !exec.Strict {
		t.Errorf("strict = false, want true")
	}
	if exec.SQLText != "SELECT a FROM t WHERE id = 1" {
		t.Errorf("sql = %q, want %q", exec.SQLText, "SELECT a FROM t WHERE id = 1")
	}
}

func TestBareSQLTextExtracted(t *testing.T) {
	block := parseOK(t, `BEGIN INSERT INTO t (a, b) SELECT x, y FROM s; END`)
	exec := block.Body[0].(*ast.PLExecSQL)
	if exec.SQLText != "INSERT INTO t (a, b) SELECT x, y FROM s" {
		t.Errorf("sql = %q, want %q", exec.SQLText, "INSERT INTO t (a, b) SELECT x, y FROM s")
	}
}

func TestBareSQLKeywords(t *testing.T) {
	// Test that all SQL keywords are recognized as statement starters
	tests := []struct {
		body string
		want string
	}{
		{`BEGIN INSERT INTO t VALUES (1); END`, "INSERT INTO t VALUES (1)"},
		{`BEGIN UPDATE t SET a = 1; END`, "UPDATE t SET a = 1"},
		{`BEGIN DELETE FROM t; END`, "DELETE FROM t"},
		{`BEGIN SELECT 1; END`, "SELECT 1"},
		{`BEGIN MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET a = s.a; END`,
			"MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET a = s.a"},
	}
	for _, tt := range tests {
		t.Run(tt.want[:6], func(t *testing.T) {
			block := parseOK(t, tt.body)
			exec, ok := block.Body[0].(*ast.PLExecSQL)
			if !ok {
				t.Fatalf("expected PLExecSQL, got %T", block.Body[0])
			}
			if exec.SQLText != tt.want {
				t.Errorf("sql = %q, want %q", exec.SQLText, tt.want)
			}
		})
	}
}

func TestBareSQLImport(t *testing.T) {
	block := parseOK(t, `BEGIN IMPORT FOREIGN SCHEMA s FROM SERVER srv INTO public; END`)
	exec, ok := block.Body[0].(*ast.PLExecSQL)
	if !ok {
		t.Fatalf("expected PLExecSQL, got %T", block.Body[0])
	}
	if exec.SQLText != "IMPORT FOREIGN SCHEMA s FROM SERVER srv INTO public" {
		t.Errorf("sql = %q, want %q", exec.SQLText, "IMPORT FOREIGN SCHEMA s FROM SERVER srv INTO public")
	}
}
