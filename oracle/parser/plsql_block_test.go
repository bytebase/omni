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

// TestPLSQLCollectionElementFieldAssignment tests collection(index).field assignment.
func TestPLSQLCollectionElementFieldAssignment(t *testing.T) {
	sql := `BEGIN
		rows(i).status := 'READY';
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	assign := block.Statements.Items[0].(*ast.PLSQLAssign)
	field, ok := assign.Target.(*ast.FieldAccessExpr)
	if !ok {
		t.Fatalf("expected FieldAccessExpr target, got %T", assign.Target)
	}
	if field.Field != "STATUS" {
		t.Errorf("expected STATUS field, got %q", field.Field)
	}
}

// TestPLSQLCursorAttributes tests cursor%ROWCOUNT and cursor%NOTFOUND expressions.
func TestPLSQLCursorAttributes(t *testing.T) {
	sql := `BEGIN
		count_rows := c%ROWCOUNT;
		EXIT WHEN c%NOTFOUND;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	assign := block.Statements.Items[0].(*ast.PLSQLAssign)
	attr, ok := assign.Value.(*ast.CursorAttributeExpr)
	if !ok {
		t.Fatalf("expected CursorAttributeExpr value, got %T", assign.Value)
	}
	if attr.Attribute != "ROWCOUNT" {
		t.Errorf("expected ROWCOUNT, got %q", attr.Attribute)
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

func TestPLSQLCursorForLoopSubquery(t *testing.T) {
	sql := `BEGIN
		FOR rec IN (SELECT id, name FROM employees WHERE active = 1) LOOP
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
	if loop.Type != ast.LOOP_CURSOR_FOR {
		t.Fatalf("expected LOOP_CURSOR_FOR, got %d", loop.Type)
	}
	if loop.CursorQuery == nil {
		t.Fatal("expected cursor query")
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

func TestPLSQLCursorDeclParenthesizedQuery(t *testing.T) {
	sql := `DECLARE
		CURSOR emp_cur IS (SELECT id, name FROM employees);
	BEGIN
		NULL;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	cursor := block.Declarations.Items[0].(*ast.PLSQLCursorDecl)
	if cursor.Query == nil {
		t.Fatal("expected cursor query")
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

func TestPLSQLQualifiedExceptionOrHandler(t *testing.T) {
	sql := `BEGIN
		NULL;
	EXCEPTION
		WHEN err_pkg.e_transient OR other_pkg.e_retry THEN
			NULL;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	handler := block.Exceptions.Items[0].(*ast.ExceptionHandler)
	if handler.Exceptions.Len() != 2 {
		t.Fatalf("expected 2 exception names, got %d", handler.Exceptions.Len())
	}
	first := handler.Exceptions.Items[0].(*ast.String)
	if first.Str != "ERR_PKG.E_TRANSIENT" {
		t.Fatalf("expected qualified exception ERR_PKG.E_TRANSIENT, got %q", first.Str)
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

// TestPLSQLOpenForDynamicSQLUsing tests OPEN cursor FOR dynamic SQL USING binds.
func TestPLSQLOpenForDynamicSQLUsing(t *testing.T) {
	sql := `BEGIN
		OPEN rc FOR stmt USING emp_id, NVL(emp_name, 'unknown'), salary + 1;
	END;`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	block := raw.Stmt.(*ast.PLSQLBlock)

	openStmt := block.Statements.Items[0].(*ast.PLSQLOpen)
	if openStmt.Cursor != "RC" {
		t.Errorf("expected RC, got %q", openStmt.Cursor)
	}
	if openStmt.ForExpr == nil {
		t.Error("expected dynamic SQL expression")
	}
	if openStmt.UsingArgs == nil || openStmt.UsingArgs.Len() != 3 {
		t.Fatalf("expected 3 USING args, got %#v", openStmt.UsingArgs)
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

// TestPLSQLBulkCollectInto covers bulk_collect_into_clause wherever PL/SQL
// allows it. Oracle 23ai compiles each block.
func TestPLSQLBulkCollectInto(t *testing.T) {
	decl := "DECLARE TYPE nt IS TABLE OF NUMBER; x nt; y nt; "
	tests := []string{
		decl + "BEGIN SELECT a, b BULK COLLECT INTO x, y FROM t; END;",
		decl + "BEGIN SELECT a BULK COLLECT INTO x FROM t WHERE ROWNUM <= 10; END;",
		decl + "BEGIN SELECT a, b BULK COLLECT INTO x, y FROM (SELECT a, b FROM t ORDER BY a) WHERE ROWNUM <= 10; END;",
		decl + "BEGIN DELETE FROM t RETURNING a, b BULK COLLECT INTO x, y; END;",
		decl + "BEGIN UPDATE t SET a = 1 RETURN a BULK COLLECT INTO x; END;",
		decl + "BEGIN FORALL i IN 1..x.COUNT DELETE FROM t WHERE a = x(i) RETURNING b BULK COLLECT INTO y; END;",
		decl + "BEGIN EXECUTE IMMEDIATE 'SELECT a, b FROM t' BULK COLLECT INTO x, y; END;",
		decl + "BEGIN EXECUTE IMMEDIATE 'DELETE FROM t WHERE a = :1 RETURNING b INTO :2' USING 1 RETURNING BULK COLLECT INTO y; END;",
		"DECLARE y NUMBER; BEGIN EXECUTE IMMEDIATE 'DELETE FROM t WHERE a = 1 RETURNING b INTO :1' RETURNING INTO y; END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
			if violations := CheckLocations(t, sql); len(violations) > 0 {
				t.Fatalf("Loc violations: %v", violations)
			}
		})
	}
}

func TestPLSQLBulkCollectIntoAST(t *testing.T) {
	result := ParseAndCheck(t, "DECLARE TYPE nt IS TABLE OF NUMBER; x nt; y nt; BEGIN "+
		"SELECT a, b BULK COLLECT INTO x, y FROM t; "+
		"EXECUTE IMMEDIATE 'SELECT a FROM t' BULK COLLECT INTO x; "+
		"EXECUTE IMMEDIATE 'DELETE FROM t RETURNING a INTO :1' RETURNING BULK COLLECT INTO y; END;")
	block := result.Items[0].(*ast.RawStmt).Stmt.(*ast.PLSQLBlock)
	sel := block.Statements.Items[0].(*ast.SelectStmt)
	if !sel.BulkCollect || sel.IntoVars.Len() != 2 {
		t.Fatalf("SELECT BulkCollect = %v, IntoVars = %d; want true, 2", sel.BulkCollect, sel.IntoVars.Len())
	}
	if sel.TargetList.Len() != 2 {
		t.Fatalf("SELECT targets = %d, want 2", sel.TargetList.Len())
	}
	exec := block.Statements.Items[1].(*ast.PLSQLExecImmediate)
	if !exec.Bulk || exec.Into.Len() != 1 {
		t.Fatalf("EXECUTE IMMEDIATE Bulk = %v, Into = %v", exec.Bulk, exec.Into)
	}
	ret := block.Statements.Items[2].(*ast.PLSQLExecImmediate)
	if !ret.ReturningBulk || ret.ReturningInto.Len() != 1 || ret.Into != nil {
		t.Fatalf("RETURNING ReturningBulk = %v, ReturningInto = %v, Into = %v", ret.ReturningBulk, ret.ReturningInto, ret.Into)
	}
}

// TestPLSQLBulkCollectRequiresInto: Oracle raises ORA-00925 (missing INTO).
func TestPLSQLBulkCollectRequiresInto(t *testing.T) {
	ParseShouldFail(t, "DECLARE TYPE nt IS TABLE OF NUMBER; x nt; BEGIN SELECT a BULK COLLECT x FROM t; END;")
	ParseShouldFail(t, "DECLARE TYPE nt IS TABLE OF NUMBER; x nt; BEGIN EXECUTE IMMEDIATE 'SELECT a FROM t' BULK COLLECT x; END;")
}

// TestSelectAliasBulk keeps BULK usable as a column alias.
func TestSelectAliasBulk(t *testing.T) {
	result := ParseAndCheck(t, "SELECT a bulk FROM t")
	sel := result.Items[0].(*ast.RawStmt).Stmt.(*ast.SelectStmt)
	if rt := sel.TargetList.Items[0].(*ast.ResTarget); rt.Name != "BULK" || sel.BulkCollect {
		t.Fatalf("alias = %q, BulkCollect = %v", rt.Name, sel.BulkCollect)
	}
}

// TestDMLReturningBulkCollectAST: DML RETURNING records BULK COLLECT.
func TestDMLReturningBulkCollectAST(t *testing.T) {
	result := ParseAndCheck(t, "DECLARE TYPE nt IS TABLE OF NUMBER; x nt; v NUMBER; BEGIN "+
		"DELETE FROM t RETURNING a BULK COLLECT INTO x; "+
		"UPDATE t SET a = 1 RETURN a BULK COLLECT INTO x; "+
		"INSERT INTO t (a) VALUES (1) RETURNING a BULK COLLECT INTO x; "+
		"DELETE FROM t WHERE a = 2 RETURNING a INTO v; END;")
	stmts := result.Items[0].(*ast.RawStmt).Stmt.(*ast.PLSQLBlock).Statements.Items
	if !stmts[0].(*ast.DeleteStmt).ReturningBulk || !stmts[1].(*ast.UpdateStmt).ReturningBulk ||
		!stmts[2].(*ast.InsertStmt).ReturningBulk || stmts[3].(*ast.DeleteStmt).ReturningBulk {
		t.Fatalf("ReturningBulk not recorded as written: %s", ast.NodeToString(result))
	}
}

// TestPLSQLTargetListsRejectEmpty: INTO, BULK COLLECT INTO, USING, and
// RETURNING INTO need at least one target and no empty entry. Oracle 23ai
// rejects each case with PLS-00103.
func TestPLSQLTargetListsRejectEmpty(t *testing.T) {
	decl := "DECLARE TYPE nt IS TABLE OF NUMBER; x nt; v NUMBER; a NUMBER; b NUMBER; CURSOR c IS SELECT a, b FROM t; "
	tests := []string{
		decl + "BEGIN EXECUTE IMMEDIATE 'SELECT 1 FROM dual' BULK COLLECT INTO; END;",
		decl + "BEGIN EXECUTE IMMEDIATE 'DELETE FROM t RETURNING a INTO :1' RETURNING BULK COLLECT INTO; END;",
		decl + "BEGIN EXECUTE IMMEDIATE 'DELETE FROM t RETURNING a INTO :1' RETURNING INTO; END;",
		decl + "BEGIN EXECUTE IMMEDIATE 'SELECT 1 FROM dual' INTO; END;",
		decl + "BEGIN EXECUTE IMMEDIATE 'SELECT 1, 2 FROM dual' INTO a,,b; END;",
		decl + "BEGIN EXECUTE IMMEDIATE 'BEGIN NULL; END;' USING; END;",
		decl + "BEGIN EXECUTE IMMEDIATE 'BEGIN NULL; END;' USING IN; END;",
		decl + "BEGIN OPEN c; FETCH c INTO; END;",
		decl + "BEGIN OPEN c; FETCH c INTO a,,b; END;",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseShouldFail(t, sql)
		})
	}
}

// TestPLSQLUsingBindModes: USING bind arguments take IN, OUT, or IN OUT, as
// Oracle 23ai accepts.
func TestPLSQLUsingBindModes(t *testing.T) {
	sql := "DECLARE v NUMBER := 1; w NUMBER; BEGIN " +
		"EXECUTE IMMEDIATE 'BEGIN :1 := :2 + 1; :3 := 0; END;' USING OUT w, IN v, IN OUT v, 5; END;"
	result := ParseAndCheck(t, sql)
	block := result.Items[0].(*ast.RawStmt).Stmt.(*ast.PLSQLBlock)
	exec := block.Statements.Items[0].(*ast.PLSQLExecImmediate)
	want := []string{"OUT", "IN", "IN OUT", ""}
	if exec.Using.Len() != len(want) || len(exec.UsingModes) != len(want) {
		t.Fatalf("Using = %d args, modes %q; want %d", exec.Using.Len(), exec.UsingModes, len(want))
	}
	for i := range want {
		if exec.UsingModes[i] != want[i] {
			t.Fatalf("UsingModes = %q, want %q", exec.UsingModes, want)
		}
	}
	if violations := CheckLocations(t, sql); len(violations) > 0 {
		t.Fatalf("Loc violations: %v", violations)
	}
	ParseAndCheck(t, "DECLARE TYPE nt IS TABLE OF NUMBER; y nt; BEGIN "+
		"EXECUTE IMMEDIATE 'DELETE FROM t WHERE a = :1 RETURNING b INTO :2' USING IN 1 RETURNING BULK COLLECT INTO y; END;")
}
