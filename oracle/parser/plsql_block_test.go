package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestPLSQLSimpleBlock tests a simple BEGIN/END block.
func TestPLSQLSimpleBlock(t *testing.T) {
	result := ParseAndCheck(t, "BEGIN NULL; END;")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	block, ok := raw.Stmt.(*ast.PLSQLBlock)
	if !ok {
		t.Fatalf("expected PLSQLBlock, got %T", raw.Stmt)
	}
	if block.Statements == nil || block.Statements.Len() != 1 {
		t.Fatalf("expected 1 statement in block, got %d", block.Statements.Len())
	}
	if _, ok := block.Statements.Items[0].(*ast.PLSQLNull); !ok {
		t.Errorf("expected PLSQLNull, got %T", block.Statements.Items[0])
	}
}

// TestPLSQLDeclareBlock tests DECLARE with variable declarations.
func TestPLSQLDeclareBlock(t *testing.T) {
	sql := `DECLARE
		v_name VARCHAR2(100);
		v_count NUMBER := 0;
		c_max CONSTANT NUMBER := 100;
	BEGIN
		NULL;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	if block.Declarations == nil || block.Declarations.Len() != 3 {
		t.Fatalf("expected 3 declarations, got %d", block.Declarations.Len())
	}

	// First: v_name VARCHAR2(100)
	d1 := block.Declarations.Items[0].(*ast.PLSQLVarDecl)
	if d1.Name != "V_NAME" {
		t.Errorf("expected V_NAME, got %q", d1.Name)
	}

	// Second: v_count NUMBER := 0
	d2 := block.Declarations.Items[1].(*ast.PLSQLVarDecl)
	if d2.Name != "V_COUNT" {
		t.Errorf("expected V_COUNT, got %q", d2.Name)
	}
	if d2.Default == nil {
		t.Error("expected default value for v_count")
	}

	// Third: c_max CONSTANT NUMBER := 100
	d3 := block.Declarations.Items[2].(*ast.PLSQLVarDecl)
	if d3.Name != "C_MAX" {
		t.Errorf("expected C_MAX, got %q", d3.Name)
	}
	if !d3.Constant {
		t.Error("expected CONSTANT to be true")
	}
}

// TestPLSQLIfElsifElse tests IF/ELSIF/ELSE/END IF.
func TestPLSQLIfElsifElse(t *testing.T) {
	sql := `BEGIN
		IF x > 0 THEN
			NULL;
		ELSIF x = 0 THEN
			NULL;
		ELSE
			NULL;
		END IF;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	if block.Statements.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", block.Statements.Len())
	}

	ifStmt, ok := block.Statements.Items[0].(*ast.PLSQLIf)
	if !ok {
		t.Fatalf("expected PLSQLIf, got %T", block.Statements.Items[0])
	}

	if ifStmt.Condition == nil {
		t.Error("expected IF condition")
	}
	if ifStmt.Then == nil || ifStmt.Then.Len() != 1 {
		t.Error("expected 1 THEN statement")
	}
	if ifStmt.ElsIfs == nil || ifStmt.ElsIfs.Len() != 1 {
		t.Errorf("expected 1 ELSIF, got %d", ifStmt.ElsIfs.Len())
	}
	if ifStmt.Else == nil || ifStmt.Else.Len() != 1 {
		t.Error("expected 1 ELSE statement")
	}
}

// TestPLSQLBasicLoop tests a basic LOOP/END LOOP.
func TestPLSQLBasicLoop(t *testing.T) {
	sql := `BEGIN
		LOOP
			NULL;
		END LOOP;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	loop, ok := block.Statements.Items[0].(*ast.PLSQLLoop)
	if !ok {
		t.Fatalf("expected PLSQLLoop, got %T", block.Statements.Items[0])
	}
	if loop.Type != ast.LOOP_BASIC {
		t.Errorf("expected LOOP_BASIC, got %d", loop.Type)
	}
	if loop.Statements.Len() != 1 {
		t.Errorf("expected 1 statement, got %d", loop.Statements.Len())
	}
}

// TestPLSQLWhileLoop tests WHILE LOOP.
func TestPLSQLWhileLoop(t *testing.T) {
	sql := `BEGIN
		WHILE x > 0 LOOP
			NULL;
		END LOOP;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	loop, ok := block.Statements.Items[0].(*ast.PLSQLLoop)
	if !ok {
		t.Fatalf("expected PLSQLLoop, got %T", block.Statements.Items[0])
	}
	if loop.Type != ast.LOOP_WHILE {
		t.Errorf("expected LOOP_WHILE, got %d", loop.Type)
	}
	if loop.Condition == nil {
		t.Error("expected WHILE condition")
	}
}

// TestPLSQLForLoop tests FOR LOOP with numeric range.
func TestPLSQLForLoop(t *testing.T) {
	sql := `BEGIN
		FOR i IN 1..10 LOOP
			NULL;
		END LOOP;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	loop, ok := block.Statements.Items[0].(*ast.PLSQLLoop)
	if !ok {
		t.Fatalf("expected PLSQLLoop, got %T", block.Statements.Items[0])
	}
	if loop.Type != ast.LOOP_FOR {
		t.Errorf("expected LOOP_FOR, got %d", loop.Type)
	}
	if loop.Iterator != "I" {
		t.Errorf("expected I, got %q", loop.Iterator)
	}
	if loop.LowerBound == nil || loop.UpperBound == nil {
		t.Error("expected lower and upper bounds")
	}
}

// TestPLSQLForLoopReverse tests FOR LOOP with REVERSE.
func TestPLSQLForLoopReverse(t *testing.T) {
	sql := `BEGIN
		FOR i IN REVERSE 1..10 LOOP
			NULL;
		END LOOP;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	loop := block.Statements.Items[0].(*ast.PLSQLLoop)
	if !loop.Reverse {
		t.Error("expected REVERSE to be true")
	}
}

// TestPLSQLAssignment tests assignment statements.
func TestPLSQLAssignment(t *testing.T) {
	sql := `BEGIN
		x := 42;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	assign, ok := block.Statements.Items[0].(*ast.PLSQLAssign)
	if !ok {
		t.Fatalf("expected PLSQLAssign, got %T", block.Statements.Items[0])
	}
	if assign.Target == nil || assign.Value == nil {
		t.Error("expected target and value in assignment")
	}
}

// TestPLSQLReturn tests RETURN statements.
func TestPLSQLReturn(t *testing.T) {
	sql := `BEGIN RETURN 42; END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	ret, ok := block.Statements.Items[0].(*ast.PLSQLReturn)
	if !ok {
		t.Fatalf("expected PLSQLReturn, got %T", block.Statements.Items[0])
	}
	if ret.Expr == nil {
		t.Error("expected return expression")
	}
}

// TestPLSQLReturnVoid tests RETURN without expression.
func TestPLSQLReturnVoid(t *testing.T) {
	sql := `BEGIN RETURN; END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	ret, ok := block.Statements.Items[0].(*ast.PLSQLReturn)
	if !ok {
		t.Fatalf("expected PLSQLReturn, got %T", block.Statements.Items[0])
	}
	if ret.Expr != nil {
		t.Error("expected nil return expression")
	}
}

// TestPLSQLNull tests NULL statement.
func TestPLSQLNullStmt(t *testing.T) {
	sql := `BEGIN NULL; END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	if _, ok := block.Statements.Items[0].(*ast.PLSQLNull); !ok {
		t.Fatalf("expected PLSQLNull, got %T", block.Statements.Items[0])
	}
}

// TestPLSQLRaise tests RAISE statement.
func TestPLSQLRaise(t *testing.T) {
	sql := `BEGIN RAISE NO_DATA_FOUND; END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	raise, ok := block.Statements.Items[0].(*ast.PLSQLRaise)
	if !ok {
		t.Fatalf("expected PLSQLRaise, got %T", block.Statements.Items[0])
	}
	if raise.Exception != "NO_DATA_FOUND" {
		t.Errorf("expected NO_DATA_FOUND, got %q", raise.Exception)
	}
}

// TestPLSQLRaiseEmpty tests RAISE without exception name (re-raise).
func TestPLSQLRaiseEmpty(t *testing.T) {
	sql := `BEGIN RAISE; END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	raise := block.Statements.Items[0].(*ast.PLSQLRaise)
	if raise.Exception != "" {
		t.Errorf("expected empty exception, got %q", raise.Exception)
	}
}

// TestPLSQLGoto tests GOTO statement.
func TestPLSQLGoto(t *testing.T) {
	sql := `BEGIN GOTO my_label; END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	g, ok := block.Statements.Items[0].(*ast.PLSQLGoto)
	if !ok {
		t.Fatalf("expected PLSQLGoto, got %T", block.Statements.Items[0])
	}
	if g.Label != "MY_LABEL" {
		t.Errorf("expected MY_LABEL, got %q", g.Label)
	}
}

// TestPLSQLExceptionHandlers tests EXCEPTION section.
func TestPLSQLExceptionHandlers(t *testing.T) {
	sql := `BEGIN
		NULL;
	EXCEPTION
		WHEN NO_DATA_FOUND THEN
			NULL;
		WHEN OTHERS THEN
			NULL;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	if block.Exceptions == nil || block.Exceptions.Len() != 2 {
		t.Fatalf("expected 2 exception handlers, got %d", block.Exceptions.Len())
	}

	h1 := block.Exceptions.Items[0].(*ast.ExceptionHandler)
	if h1.Exceptions.Len() != 1 {
		t.Errorf("expected 1 exception name, got %d", h1.Exceptions.Len())
	}
	name1 := h1.Exceptions.Items[0].(*ast.String)
	if name1.Str != "NO_DATA_FOUND" {
		t.Errorf("expected NO_DATA_FOUND, got %q", name1.Str)
	}

	h2 := block.Exceptions.Items[1].(*ast.ExceptionHandler)
	name2 := h2.Exceptions.Items[0].(*ast.String)
	if name2.Str != "OTHERS" {
		t.Errorf("expected OTHERS, got %q", name2.Str)
	}
}

// TestPLSQLNestedBlocks tests nested BEGIN/END blocks.
func TestPLSQLNestedBlocks(t *testing.T) {
	sql := `BEGIN
		BEGIN
			NULL;
		END;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	if block.Statements.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", block.Statements.Len())
	}

	inner, ok := block.Statements.Items[0].(*ast.PLSQLBlock)
	if !ok {
		t.Fatalf("expected nested PLSQLBlock, got %T", block.Statements.Items[0])
	}
	if inner.Statements.Len() != 1 {
		t.Errorf("expected 1 statement in inner block, got %d", inner.Statements.Len())
	}
}

// TestPLSQLExecImmediate tests EXECUTE IMMEDIATE.
func TestPLSQLExecImmediate(t *testing.T) {
	sql := `BEGIN
		EXECUTE IMMEDIATE 'DROP TABLE temp_t';
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	exec, ok := block.Statements.Items[0].(*ast.PLSQLExecImmediate)
	if !ok {
		t.Fatalf("expected PLSQLExecImmediate, got %T", block.Statements.Items[0])
	}
	if exec.SQL == nil {
		t.Error("expected SQL expression")
	}
}

// TestPLSQLExecImmediateIntoUsing tests EXECUTE IMMEDIATE with INTO and USING.
func TestPLSQLExecImmediateIntoUsing(t *testing.T) {
	sql := `BEGIN
		EXECUTE IMMEDIATE 'SELECT name FROM emp WHERE id = :1' INTO v_name USING v_id;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	exec := block.Statements.Items[0].(*ast.PLSQLExecImmediate)
	if exec.Into == nil || exec.Into.Len() != 1 {
		t.Errorf("expected 1 INTO var, got %d", exec.Into.Len())
	}
	if exec.Using == nil || exec.Using.Len() != 1 {
		t.Errorf("expected 1 USING var, got %d", exec.Using.Len())
	}
}

// TestPLSQLOpenFetchClose tests cursor operations.
func TestPLSQLOpenFetchClose(t *testing.T) {
	sql := `BEGIN
		OPEN my_cursor;
		FETCH my_cursor INTO v_name;
		CLOSE my_cursor;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	if block.Statements.Len() != 3 {
		t.Fatalf("expected 3 statements, got %d", block.Statements.Len())
	}

	openStmt, ok := block.Statements.Items[0].(*ast.PLSQLOpen)
	if !ok {
		t.Fatalf("expected PLSQLOpen, got %T", block.Statements.Items[0])
	}
	if openStmt.Cursor != "MY_CURSOR" {
		t.Errorf("expected MY_CURSOR, got %q", openStmt.Cursor)
	}

	fetchStmt, ok := block.Statements.Items[1].(*ast.PLSQLFetch)
	if !ok {
		t.Fatalf("expected PLSQLFetch, got %T", block.Statements.Items[1])
	}
	if fetchStmt.Cursor != "MY_CURSOR" {
		t.Errorf("expected MY_CURSOR, got %q", fetchStmt.Cursor)
	}
	if fetchStmt.Into == nil || fetchStmt.Into.Len() != 1 {
		t.Errorf("expected 1 INTO var, got %d", fetchStmt.Into.Len())
	}

	closeStmt, ok := block.Statements.Items[2].(*ast.PLSQLClose)
	if !ok {
		t.Fatalf("expected PLSQLClose, got %T", block.Statements.Items[2])
	}
	if closeStmt.Cursor != "MY_CURSOR" {
		t.Errorf("expected MY_CURSOR, got %q", closeStmt.Cursor)
	}
}

// TestPLSQLFetchBulkCollect tests FETCH with BULK COLLECT.
func TestPLSQLFetchBulkCollect(t *testing.T) {
	sql := `BEGIN
		FETCH my_cursor BULK COLLECT INTO v_names LIMIT 100;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	fetchStmt := block.Statements.Items[0].(*ast.PLSQLFetch)
	if !fetchStmt.Bulk {
		t.Error("expected BULK to be true")
	}
	if fetchStmt.Limit == nil {
		t.Error("expected LIMIT expression")
	}
}

// TestPLSQLCursorDecl tests cursor declaration.
func TestPLSQLCursorDecl(t *testing.T) {
	sql := `DECLARE
		CURSOR emp_cur IS SELECT id, name FROM employees;
	BEGIN
		NULL;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	if block.Declarations.Len() != 1 {
		t.Fatalf("expected 1 declaration, got %d", block.Declarations.Len())
	}
	cursor := block.Declarations.Items[0].(*ast.PLSQLCursorDecl)
	if cursor.Name != "EMP_CUR" {
		t.Errorf("expected EMP_CUR, got %q", cursor.Name)
	}
	if cursor.Query == nil {
		t.Error("expected cursor query")
	}
}

// TestPLSQLLabeledBlock tests a labeled block.
func TestPLSQLLabeledBlock(t *testing.T) {
	sql := `<<my_block>> BEGIN NULL; END my_block;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	if block.Label != "MY_BLOCK" {
		t.Errorf("expected MY_BLOCK, got %q", block.Label)
	}
}

// TestPLSQLNotNullDecl tests NOT NULL constraint on variable declaration.
func TestPLSQLNotNullDecl(t *testing.T) {
	sql := `DECLARE
		v_id NUMBER NOT NULL := 1;
	BEGIN
		NULL;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	decl := block.Declarations.Items[0].(*ast.PLSQLVarDecl)
	if !decl.NotNull {
		t.Error("expected NOT NULL to be true")
	}
	if decl.Default == nil {
		t.Error("expected default value")
	}
}

// TestPLSQLExceptionOrHandler tests WHEN exc1 OR exc2 THEN.
func TestPLSQLExceptionOrHandler(t *testing.T) {
	sql := `BEGIN
		NULL;
	EXCEPTION
		WHEN NO_DATA_FOUND OR TOO_MANY_ROWS THEN
			NULL;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	handler := block.Exceptions.Items[0].(*ast.ExceptionHandler)
	if handler.Exceptions.Len() != 2 {
		t.Fatalf("expected 2 exception names, got %d", handler.Exceptions.Len())
	}
}

// TestPLSQLDeclareDefault tests DEFAULT keyword in declarations.
func TestPLSQLDeclareDefault(t *testing.T) {
	sql := `DECLARE
		v_x NUMBER DEFAULT 10;
	BEGIN
		NULL;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	decl := block.Declarations.Items[0].(*ast.PLSQLVarDecl)
	if decl.Default == nil {
		t.Error("expected default value")
	}
}

// TestPLSQLOpenFor tests OPEN cursor FOR query.
func TestPLSQLOpenFor(t *testing.T) {
	sql := `BEGIN
		OPEN rc FOR SELECT id FROM employees;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	openStmt := block.Statements.Items[0].(*ast.PLSQLOpen)
	if openStmt.Cursor != "RC" {
		t.Errorf("expected RC, got %q", openStmt.Cursor)
	}
	if openStmt.ForQuery == nil {
		t.Error("expected FOR query")
	}
}

// TestPLSQLMultipleStatements tests multiple statements in a block.
func TestPLSQLMultipleStatements(t *testing.T) {
	sql := `BEGIN
		x := 1;
		y := 2;
		NULL;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	if block.Statements.Len() != 3 {
		t.Fatalf("expected 3 statements, got %d", block.Statements.Len())
	}
}
