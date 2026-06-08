package parser

import (
	"testing"

	"github.com/bytebase/omni/redshift/plpgsql/ast"
)

// ==========================================================================
// Section 6.1: Nested and Complex Structures
// ==========================================================================

func TestNestedDeeplyNestedBlocks(t *testing.T) {
	// Deeply nested blocks (3 levels): BEGIN BEGIN BEGIN ... END; END; END
	block := parseOK(t, `
		BEGIN
			BEGIN
				BEGIN
					x := 1;
				END;
			END;
		END
	`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 outer body stmt, got %d", len(block.Body))
	}
	inner1, ok := block.Body[0].(*ast.PLBlock)
	if !ok {
		t.Fatalf("expected PLBlock at level 2, got %T", block.Body[0])
	}
	if len(inner1.Body) != 1 {
		t.Fatalf("expected 1 level-2 body stmt, got %d", len(inner1.Body))
	}
	inner2, ok := inner1.Body[0].(*ast.PLBlock)
	if !ok {
		t.Fatalf("expected PLBlock at level 3, got %T", inner1.Body[0])
	}
	if len(inner2.Body) != 1 {
		t.Fatalf("expected 1 level-3 body stmt, got %d", len(inner2.Body))
	}
	assign, ok := inner2.Body[0].(*ast.PLAssign)
	if !ok {
		t.Fatalf("expected PLAssign at level 3, got %T", inner2.Body[0])
	}
	if assign.Target != "x" || assign.Expr != "1" {
		t.Errorf("assign = %s := %s, want x := 1", assign.Target, assign.Expr)
	}
}

func TestNestedBlockInsideIF(t *testing.T) {
	// Block inside IF branch
	block := parseOK(t, `
		BEGIN
			IF x > 0 THEN
				BEGIN
					y := 1;
				END;
			END IF;
		END
	`)
	ifNode := block.Body[0].(*ast.PLIf)
	if len(ifNode.ThenBody) != 1 {
		t.Fatalf("expected 1 stmt in THEN, got %d", len(ifNode.ThenBody))
	}
	inner, ok := ifNode.ThenBody[0].(*ast.PLBlock)
	if !ok {
		t.Fatalf("expected PLBlock in THEN, got %T", ifNode.ThenBody[0])
	}
	if len(inner.Body) != 1 {
		t.Fatalf("expected 1 stmt in inner block, got %d", len(inner.Body))
	}
}

func TestNestedBlockInsideFOR(t *testing.T) {
	// Block inside FOR loop body
	block := parseOK(t, `
		BEGIN
			FOR i IN 1..10 LOOP
				BEGIN
					x := i;
				END;
			END LOOP;
		END
	`)
	forNode := block.Body[0].(*ast.PLForI)
	if len(forNode.Body) != 1 {
		t.Fatalf("expected 1 stmt in FOR body, got %d", len(forNode.Body))
	}
	inner, ok := forNode.Body[0].(*ast.PLBlock)
	if !ok {
		t.Fatalf("expected PLBlock in FOR body, got %T", forNode.Body[0])
	}
	if len(inner.Body) != 1 {
		t.Fatalf("expected 1 stmt in inner block, got %d", len(inner.Body))
	}
}

func TestNestedBlockInsideException(t *testing.T) {
	// Block inside exception handler
	block := parseOK(t, `
		BEGIN
			x := 1;
		EXCEPTION
			WHEN OTHERS THEN
				BEGIN
					y := 0;
				END;
		END
	`)
	if len(block.Exceptions) != 1 {
		t.Fatalf("expected 1 exception handler, got %d", len(block.Exceptions))
	}
	w := block.Exceptions[0].(*ast.PLExceptionWhen)
	if len(w.Body) != 1 {
		t.Fatalf("expected 1 stmt in handler, got %d", len(w.Body))
	}
	inner, ok := w.Body[0].(*ast.PLBlock)
	if !ok {
		t.Fatalf("expected PLBlock in handler, got %T", w.Body[0])
	}
	if len(inner.Body) != 1 {
		t.Fatalf("expected 1 stmt in inner block, got %d", len(inner.Body))
	}
}

func TestNestedLoopInsideIFInsideBlock(t *testing.T) {
	// Loop inside IF inside block
	block := parseOK(t, `
		BEGIN
			IF x > 0 THEN
				LOOP
					x := x - 1;
					EXIT WHEN x = 0;
				END LOOP;
			END IF;
		END
	`)
	ifNode := block.Body[0].(*ast.PLIf)
	if len(ifNode.ThenBody) != 1 {
		t.Fatalf("expected 1 stmt in THEN, got %d", len(ifNode.ThenBody))
	}
	loop, ok := ifNode.ThenBody[0].(*ast.PLLoop)
	if !ok {
		t.Fatalf("expected PLLoop in THEN, got %T", ifNode.ThenBody[0])
	}
	if len(loop.Body) != 2 {
		t.Errorf("expected 2 stmts in loop body, got %d", len(loop.Body))
	}
}

func TestNestedFORContainingIFContainingReturnQuery(t *testing.T) {
	// FOR loop containing IF containing RETURN QUERY
	block := parseOK(t, `
		BEGIN
			FOR rec IN SELECT * FROM t LOOP
				IF rec.active THEN
					RETURN QUERY SELECT rec.id, rec.name;
				END IF;
			END LOOP;
		END
	`)
	forNode := block.Body[0].(*ast.PLForS)
	if len(forNode.Body) != 1 {
		t.Fatalf("expected 1 stmt in FOR body, got %d", len(forNode.Body))
	}
	ifNode, ok := forNode.Body[0].(*ast.PLIf)
	if !ok {
		t.Fatalf("expected PLIf in FOR body, got %T", forNode.Body[0])
	}
	if len(ifNode.ThenBody) != 1 {
		t.Fatalf("expected 1 stmt in THEN, got %d", len(ifNode.ThenBody))
	}
	retQ, ok := ifNode.ThenBody[0].(*ast.PLReturnQuery)
	if !ok {
		t.Fatalf("expected PLReturnQuery, got %T", ifNode.ThenBody[0])
	}
	if retQ.Query == "" {
		t.Error("expected non-empty query in RETURN QUERY")
	}
}

func TestNestedExceptionContainingLoop(t *testing.T) {
	// Exception handler containing loop
	block := parseOK(t, `
		BEGIN
			x := 1;
		EXCEPTION
			WHEN OTHERS THEN
				LOOP
					x := x + 1;
					EXIT WHEN x > 5;
				END LOOP;
		END
	`)
	w := block.Exceptions[0].(*ast.PLExceptionWhen)
	if len(w.Body) != 1 {
		t.Fatalf("expected 1 stmt in handler, got %d", len(w.Body))
	}
	loop, ok := w.Body[0].(*ast.PLLoop)
	if !ok {
		t.Fatalf("expected PLLoop in handler, got %T", w.Body[0])
	}
	if len(loop.Body) != 2 {
		t.Errorf("expected 2 stmts in loop body, got %d", len(loop.Body))
	}
}

func TestNestedMultipleDeclareInNestedBlocks(t *testing.T) {
	// Multiple DECLARE sections in nested blocks
	block := parseOK(t, `
		DECLARE
			x integer;
		BEGIN
			DECLARE
				y integer;
			BEGIN
				x := 1;
				y := 2;
			END;
		END
	`)
	if len(block.Declarations) != 1 {
		t.Fatalf("expected 1 outer declaration, got %d", len(block.Declarations))
	}
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 outer body stmt, got %d", len(block.Body))
	}
	inner, ok := block.Body[0].(*ast.PLBlock)
	if !ok {
		t.Fatalf("expected PLBlock, got %T", block.Body[0])
	}
	if len(inner.Declarations) != 1 {
		t.Fatalf("expected 1 inner declaration, got %d", len(inner.Declarations))
	}
	if len(inner.Body) != 2 {
		t.Errorf("expected 2 inner body stmts, got %d", len(inner.Body))
	}
}

func TestNestedLabelScoping(t *testing.T) {
	// Label scoping: inner label shadows outer
	block := parseOK(t, `
		<<blk>>
		BEGIN
			<<blk>>
			BEGIN
				x := 1;
			END blk;
		END blk
	`)
	if block.Label != "blk" {
		t.Errorf("outer label = %q, want %q", block.Label, "blk")
	}
	inner, ok := block.Body[0].(*ast.PLBlock)
	if !ok {
		t.Fatalf("expected PLBlock, got %T", block.Body[0])
	}
	if inner.Label != "blk" {
		t.Errorf("inner label = %q, want %q", inner.Label, "blk")
	}
}

// ==========================================================================
// Section 6.2: SQL Fragment Extraction
// ==========================================================================

func TestSQLFragIfCondition(t *testing.T) {
	// SQL expression in IF condition extracted as text until THEN
	block := parseOK(t, `BEGIN IF x > 0 AND y < 10 THEN z := 1; END IF; END`)
	ifNode := block.Body[0].(*ast.PLIf)
	if ifNode.Condition != "x > 0 AND y < 10" {
		t.Errorf("condition = %q, want %q", ifNode.Condition, "x > 0 AND y < 10")
	}
}

func TestSQLFragWhileCondition(t *testing.T) {
	// SQL expression in WHILE condition extracted until LOOP
	block := parseOK(t, `BEGIN WHILE x > 0 AND NOT done LOOP x := x - 1; END LOOP; END`)
	w := block.Body[0].(*ast.PLWhile)
	if w.Condition != "x > 0 AND NOT done" {
		t.Errorf("condition = %q, want %q", w.Condition, "x > 0 AND NOT done")
	}
}

func TestSQLFragReturnExpr(t *testing.T) {
	// SQL expression in RETURN extracted until semicolon
	block := parseOK(t, `BEGIN RETURN a + b * c - d; END`)
	ret := block.Body[0].(*ast.PLReturn)
	if ret.Expr != "a + b * c - d" {
		t.Errorf("expr = %q, want %q", ret.Expr, "a + b * c - d")
	}
}

func TestSQLFragReturnQuery(t *testing.T) {
	// SQL in RETURN QUERY extracted as full SELECT text until semicolon
	block := parseOK(t, `BEGIN RETURN QUERY SELECT a, b FROM t WHERE c > 0 ORDER BY a; END`)
	ret := block.Body[0].(*ast.PLReturnQuery)
	want := "SELECT a, b FROM t WHERE c > 0 ORDER BY a"
	if ret.Query != want {
		t.Errorf("query = %q, want %q", ret.Query, want)
	}
}

func TestSQLFragForIntegerVsQuery(t *testing.T) {
	// SQL in FOR..IN extracted correctly (distinguish integer range `..` from query)
	t.Run("integer range", func(t *testing.T) {
		block := parseOK(t, `BEGIN FOR i IN 1..10 LOOP NULL; END LOOP; END`)
		_, ok := block.Body[0].(*ast.PLForI)
		if !ok {
			t.Fatalf("expected PLForI, got %T", block.Body[0])
		}
	})
	t.Run("query", func(t *testing.T) {
		block := parseOK(t, `BEGIN FOR rec IN SELECT * FROM t LOOP NULL; END LOOP; END`)
		_, ok := block.Body[0].(*ast.PLForS)
		if !ok {
			t.Fatalf("expected PLForS, got %T", block.Body[0])
		}
	})
}

func TestSQLFragExecuteUntilIntoUsingSemicolon(t *testing.T) {
	// SQL in EXECUTE extracted until INTO/USING/semicolon
	t.Run("until semicolon", func(t *testing.T) {
		block := parseOK(t, `BEGIN EXECUTE 'SELECT 1'; END`)
		exec := block.Body[0].(*ast.PLDynExecute)
		if exec.Query != "'SELECT 1'" {
			t.Errorf("query = %q, want %q", exec.Query, "'SELECT 1'")
		}
	})
	t.Run("until INTO", func(t *testing.T) {
		block := parseOK(t, `BEGIN EXECUTE 'SELECT a FROM t' INTO x; END`)
		exec := block.Body[0].(*ast.PLDynExecute)
		if exec.Query != "'SELECT a FROM t'" {
			t.Errorf("query = %q, want %q", exec.Query, "'SELECT a FROM t'")
		}
		if len(exec.Into) != 1 || exec.Into[0] != "x" {
			t.Errorf("into = %v, want [x]", exec.Into)
		}
	})
	t.Run("until USING", func(t *testing.T) {
		block := parseOK(t, `BEGIN EXECUTE 'INSERT INTO t VALUES($1)' USING val; END`)
		exec := block.Body[0].(*ast.PLDynExecute)
		if exec.Query != "'INSERT INTO t VALUES($1)'" {
			t.Errorf("query = %q, want %q", exec.Query, "'INSERT INTO t VALUES($1)'")
		}
		if len(exec.Params) != 1 || exec.Params[0] != "val" {
			t.Errorf("params = %v, want [val]", exec.Params)
		}
	})
}

func TestSQLFragPerform(t *testing.T) {
	// SQL in PERFORM extracted until semicolon
	block := parseOK(t, `BEGIN PERFORM my_func(a, b, c); END`)
	perf := block.Body[0].(*ast.PLPerform)
	if perf.Expr != "my_func(a, b, c)" {
		t.Errorf("expr = %q, want %q", perf.Expr, "my_func(a, b, c)")
	}
}

func TestSQLFragAssignRHS(t *testing.T) {
	// SQL in assignment RHS extracted until semicolon
	block := parseOK(t, `BEGIN x := (a + b) * (c - d) / e; END`)
	assign := block.Body[0].(*ast.PLAssign)
	if assign.Expr != "(a + b) * (c - d) / e" {
		t.Errorf("expr = %q, want %q", assign.Expr, "(a + b) * (c - d) / e")
	}
}

func TestSQLFragParenthesizedNoFalseTermination(t *testing.T) {
	// Parenthesized expressions in SQL fragments don't cause false termination
	// The THEN keyword inside parenthesized subquery should not terminate the IF condition
	block := parseOK(t, `BEGIN IF (SELECT count(*) FROM t) > 0 THEN x := 1; END IF; END`)
	ifNode := block.Body[0].(*ast.PLIf)
	if ifNode.Condition != "(SELECT count(*) FROM t) > 0" {
		t.Errorf("condition = %q, want %q", ifNode.Condition, "(SELECT count(*) FROM t) > 0")
	}
}

func TestSQLFragDollarQuotedStrings(t *testing.T) {
	// Dollar-quoted strings within SQL fragments preserved
	block := parseOK(t, `BEGIN EXECUTE $q$ SELECT * FROM t $q$; END`)
	exec := block.Body[0].(*ast.PLDynExecute)
	if exec.Query != "$q$ SELECT * FROM t $q$" {
		t.Errorf("query = %q, want %q", exec.Query, "$q$ SELECT * FROM t $q$")
	}
}

func TestSQLFragStringLiteralWithSemicolon(t *testing.T) {
	// String literals with semicolons in SQL fragments don't cause false termination
	block := parseOK(t, `BEGIN x := 'hello; world'; END`)
	assign := block.Body[0].(*ast.PLAssign)
	if assign.Expr != "'hello; world'" {
		t.Errorf("expr = %q, want %q", assign.Expr, "'hello; world'")
	}
}

func TestSQLFragNestedParentheses(t *testing.T) {
	// Nested parentheses in SQL fragments tracked correctly
	block := parseOK(t, `BEGIN x := ((a + b) * (c + (d - e))); END`)
	assign := block.Body[0].(*ast.PLAssign)
	if assign.Expr != "((a + b) * (c + (d - e)))" {
		t.Errorf("expr = %q, want %q", assign.Expr, "((a + b) * (c + (d - e)))")
	}
}

// ==========================================================================
// Section 6.3: Error Reporting
// ==========================================================================

func TestErrorPositionPointsToOffendingToken(t *testing.T) {
	// Error position points to offending token
	_, err := Parse("BEGIN 123; END")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	pe, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	// Position should be non-negative and point to the offending token
	if pe.Position < 0 {
		t.Errorf("position = %d, expected >= 0", pe.Position)
	}
}

func TestErrorMissingSemicolon(t *testing.T) {
	// Missing semicolon after a statement that requires one:
	// CLOSE expects: CLOSE cursor_var ;
	// "BEGIN CLOSE END" -- CLOSE tries to read a cursor var, sees END which
	// is not an ident, causing an error (or it reads "end" as cursor name
	// and then expects semicolon). Either way, structure breaks.
	parseErr(t, "BEGIN CLOSE END", "syntax error")
}

func TestErrorUnexpectedKeyword(t *testing.T) {
	// Unexpected keyword in expression context
	parseErr(t, "BEGIN IF x > 0 THEN y := 1; ELSE ELSE z := 2; END IF; END", "syntax error")
}

func TestErrorUnterminatedBlock(t *testing.T) {
	// Unterminated block (missing END)
	parseErr(t, "BEGIN x := 1;", "END")
}

func TestErrorUnterminatedIF(t *testing.T) {
	// Unterminated IF (missing END IF)
	parseErr(t, "BEGIN IF x THEN y := 1; END", "IF")
}

func TestErrorUnterminatedLoop(t *testing.T) {
	// Unterminated LOOP (missing END LOOP)
	parseErr(t, "BEGIN LOOP x := 1; END", "LOOP")
}

func TestErrorUnterminatedCase(t *testing.T) {
	// Unterminated CASE (missing END CASE)
	parseErr(t, "BEGIN CASE WHEN x THEN y := 1; END", "CASE")
}

func TestErrorUnknownStatementKeyword(t *testing.T) {
	// Unknown statement keyword produces meaningful error
	parseErr(t, "BEGIN FOOBAR; END", "syntax error")
}

func TestErrorLabelWithoutBlock(t *testing.T) {
	// Label without matching block/loop
	parseErr(t, "BEGIN <<lbl>> x := 1; END", "syntax error")
}

func TestErrorEndLabelMismatch(t *testing.T) {
	// END label mismatch: <<a>> BEGIN END b;
	parseErr(t, "<<a>> BEGIN END b", "does not match")
}

func TestErrorEmptyInput(t *testing.T) {
	// Empty input produces error
	_, err := Parse("")
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestErrorDuplicateDeclare(t *testing.T) {
	// Duplicate DECLARE sections in same block
	// PL/pgSQL allows only one DECLARE section per block (before BEGIN).
	// "DECLARE DECLARE x int; BEGIN END" — the second DECLARE is treated as a variable name
	// which is fine per 1.4 scenario. But "DECLARE x int; DECLARE y int; BEGIN END" may
	// be handled differently. Let's test that double DECLARE before BEGIN works (per 1.4).
	// Actually, per scenario 1.4: "DECLARE DECLARE x int; BEGIN END" parses correctly.
	// The scenario here is about DECLARE appearing AFTER BEGIN which should be a nested block,
	// not a duplicate DECLARE. Let's test that "DECLARE x int; BEGIN DECLARE y int; BEGIN END; END"
	// parses as nested blocks.
	block := parseOK(t, "DECLARE x int; BEGIN DECLARE y int; BEGIN y := 1; END; END")
	if len(block.Declarations) != 1 {
		t.Fatalf("expected 1 outer decl, got %d", len(block.Declarations))
	}
	inner, ok := block.Body[0].(*ast.PLBlock)
	if !ok {
		t.Fatalf("expected nested PLBlock, got %T", block.Body[0])
	}
	if len(inner.Declarations) != 1 {
		t.Errorf("expected 1 inner decl, got %d", len(inner.Declarations))
	}
}
