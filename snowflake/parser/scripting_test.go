package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// mustScriptBlock parses input (a single segment) into a *ast.ScriptBlockStmt
// via the public best-effort entry point, asserting no errors and exactly one
// statement.
func mustScriptBlock(t *testing.T, input string) *ast.ScriptBlockStmt {
	t.Helper()
	node := mustParseOne(t, input)
	block, ok := node.(*ast.ScriptBlockStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.ScriptBlockStmt", input, node)
	}
	return block
}

// parseScriptErr parses input and returns the collected errors. It does NOT
// fail the test, so callers can assert that an error occurred (negatives /
// unterminated-block fast-fail).
func parseScriptErr(input string) []ParseError {
	return ParseBestEffort(input).Errors
}

// firstStmtType returns the concrete type name of the first body statement of a
// block (helper for table-driven dispatch assertions).
func firstBodyStmt(t *testing.T, block *ast.ScriptBlockStmt) ast.Node {
	t.Helper()
	if len(block.Body) == 0 {
		t.Fatalf("block has empty body")
	}
	return block.Body[0]
}

// decl returns the i-th declaration of a block as a *ast.ScriptDeclaration.
func decl(t *testing.T, block *ast.ScriptBlockStmt, i int) *ast.ScriptDeclaration {
	t.Helper()
	if i >= len(block.Decls) {
		t.Fatalf("decl index %d out of range (have %d)", i, len(block.Decls))
	}
	d, ok := block.Decls[i].(*ast.ScriptDeclaration)
	if !ok {
		t.Fatalf("decl[%d] = %T, want *ast.ScriptDeclaration", i, block.Decls[i])
	}
	return d
}

// handler returns the i-th exception handler of a block.
func handler(t *testing.T, block *ast.ScriptBlockStmt, i int) *ast.ScriptExceptionHandler {
	t.Helper()
	if i >= len(block.Handlers) {
		t.Fatalf("handler index %d out of range (have %d)", i, len(block.Handlers))
	}
	h, ok := block.Handlers[i].(*ast.ScriptExceptionHandler)
	if !ok {
		t.Fatalf("handler[%d] = %T, want *ast.ScriptExceptionHandler", i, block.Handlers[i])
	}
	return h
}

// ifBranch returns the i-th branch of an IF statement.
func ifBranch(t *testing.T, s *ast.ScriptIfStmt, i int) *ast.ScriptIfBranch {
	t.Helper()
	if i >= len(s.Branches) {
		t.Fatalf("branch index %d out of range (have %d)", i, len(s.Branches))
	}
	b, ok := s.Branches[i].(*ast.ScriptIfBranch)
	if !ok {
		t.Fatalf("branch[%d] = %T, want *ast.ScriptIfBranch", i, s.Branches[i])
	}
	return b
}

// ---------------------------------------------------------------------------
// Block: BEGIN..END / DECLARE..BEGIN..END / EXCEPTION / label
// ---------------------------------------------------------------------------

func TestScript_BareBeginEnd(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN SELECT 1; END")
	if len(block.Decls) != 0 {
		t.Errorf("decls = %d, want 0", len(block.Decls))
	}
	if len(block.Handlers) != 0 {
		t.Errorf("handlers = %d, want 0", len(block.Handlers))
	}
	if !block.Label.IsEmpty() {
		t.Errorf("label = %q, want empty", block.Label.Normalize())
	}
	if len(block.Body) != 1 {
		t.Fatalf("body = %d, want 1", len(block.Body))
	}
	if _, ok := block.Body[0].(*ast.SelectStmt); !ok {
		t.Errorf("body[0] = %T, want *ast.SelectStmt", block.Body[0])
	}
}

func TestScript_MultipleBodyStatements(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN SELECT 1; INSERT INTO t VALUES (1); SELECT 2; END")
	if len(block.Body) != 3 {
		t.Fatalf("body = %d, want 3", len(block.Body))
	}
}

func TestScript_TrailingSemicolonOptionalBeforeEnd(t *testing.T) {
	// No ';' before END (lenient — legacy task_scripting_statement_list SEMI?).
	block := mustScriptBlock(t, "BEGIN SELECT 1; SELECT 2 END")
	if len(block.Body) != 2 {
		t.Fatalf("body = %d, want 2", len(block.Body))
	}
}

func TestScript_DeclareBeginEnd(t *testing.T) {
	block := mustScriptBlock(t, "DECLARE x INT; BEGIN RETURN x; END")
	if len(block.Decls) != 1 {
		t.Fatalf("decls = %d, want 1", len(block.Decls))
	}
	d := decl(t, block, 0)
	if d.Kind != ast.ScriptDeclVar {
		t.Errorf("decl kind = %v, want ScriptDeclVar", d.Kind)
	}
	if d.Name.Name != "x" {
		t.Errorf("decl name = %q, want x", d.Name.Normalize())
	}
	if d.Type == nil {
		t.Errorf("decl type = nil, want INT")
	}
}

func TestScript_EndLabel(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN SELECT 1; END my_label")
	if block.Label.Name != "my_label" {
		t.Errorf("label = %q, want my_label", block.Label.Normalize())
	}
}

// TestScript_NestedEndNoSemicolon guards that an inner block's END does not
// swallow the OUTER block's END as a label when no ';' separates them
// (`BEGIN BEGIN ... END END`). END is a non-reserved keyword, so a naive
// label scan would consume it; tryParseScriptLabel must reject keyword tokens.
func TestScript_NestedEndNoSemicolon(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN BEGIN SELECT 1; END END")
	if !block.Label.IsEmpty() {
		t.Errorf("outer label = %q, want empty (END is not a label)", block.Label.Name)
	}
	if len(block.Body) != 1 {
		t.Fatalf("outer body = %d, want 1 (the inner block)", len(block.Body))
	}
	if _, ok := block.Body[0].(*ast.ScriptBlockStmt); !ok {
		t.Errorf("body[0] = %T, want *ast.ScriptBlockStmt", block.Body[0])
	}
}

func TestScript_NestedBlock(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN BEGIN SELECT 1; END; SELECT 2; END")
	if len(block.Body) != 2 {
		t.Fatalf("outer body = %d, want 2", len(block.Body))
	}
	inner, ok := block.Body[0].(*ast.ScriptBlockStmt)
	if !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptBlockStmt (nested)", block.Body[0])
	}
	if len(inner.Body) != 1 {
		t.Errorf("inner body = %d, want 1", len(inner.Body))
	}
}

func TestScript_BlockLoc(t *testing.T) {
	input := "BEGIN SELECT 1; END"
	block := mustScriptBlock(t, input)
	if block.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", block.Loc.Start)
	}
	if block.Loc.End != len(input) {
		t.Errorf("Loc.End = %d, want %d", block.Loc.End, len(input))
	}
}

// ---------------------------------------------------------------------------
// DECLARE forms
// ---------------------------------------------------------------------------

func TestScript_DeclareVarForms(t *testing.T) {
	tests := []struct {
		name     string
		decl     string // the declaration text (without trailing ';')
		wantType bool
		wantDef  bool
	}{
		{"type only", "x INT", true, false},
		{"type and default kw", "x INT DEFAULT 0", true, true},
		{"type and assign", "x INT := 0", true, true},
		{"typeless assign", "x := 0", false, true},
		{"varchar typed", "name VARCHAR DEFAULT 'a'", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := mustScriptBlock(t, "DECLARE "+tt.decl+"; BEGIN RETURN x; END")
			if len(block.Decls) != 1 {
				t.Fatalf("decls = %d, want 1", len(block.Decls))
			}
			d := decl(t, block, 0)
			if d.Kind != ast.ScriptDeclVar {
				t.Fatalf("kind = %v, want ScriptDeclVar", d.Kind)
			}
			if (d.Type != nil) != tt.wantType {
				t.Errorf("hasType = %v, want %v", d.Type != nil, tt.wantType)
			}
			if (d.Default != nil) != tt.wantDef {
				t.Errorf("hasDefault = %v, want %v", d.Default != nil, tt.wantDef)
			}
		})
	}
}

func TestScript_DeclareCursor(t *testing.T) {
	block := mustScriptBlock(t, "DECLARE c1 CURSOR FOR SELECT price FROM invoices; BEGIN OPEN c1; END")
	d := decl(t, block, 0)
	if d.Kind != ast.ScriptDeclCursor {
		t.Fatalf("kind = %v, want ScriptDeclCursor", d.Kind)
	}
	if d.Query == nil {
		t.Errorf("cursor query = nil, want a SELECT")
	}
	if _, ok := d.Query.(*ast.SelectStmt); !ok {
		t.Errorf("cursor query = %T, want *ast.SelectStmt", d.Query)
	}
}

func TestScript_DeclareResultset(t *testing.T) {
	tests := []struct {
		name      string
		decl      string
		wantQuery bool
		wantAsync bool
	}{
		{"bare", "rs RESULTSET", false, false},
		{"assigned", "rs RESULTSET := (SELECT 1)", true, false},
		{"default", "rs RESULTSET DEFAULT (SELECT 1)", true, false},
		{"async", "rs RESULTSET := ASYNC (SELECT 1)", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := mustScriptBlock(t, "DECLARE "+tt.decl+"; BEGIN RETURN 1; END")
			d := decl(t, block, 0)
			if d.Kind != ast.ScriptDeclResultset {
				t.Fatalf("kind = %v, want ScriptDeclResultset", d.Kind)
			}
			if (d.Query != nil) != tt.wantQuery {
				t.Errorf("hasQuery = %v, want %v", d.Query != nil, tt.wantQuery)
			}
			if d.Async != tt.wantAsync {
				t.Errorf("async = %v, want %v", d.Async, tt.wantAsync)
			}
		})
	}
}

func TestScript_DeclareException(t *testing.T) {
	tests := []struct {
		name    string
		decl    string
		wantArg int
	}{
		{"bare", "e EXCEPTION", 0},
		{"with args", "e EXCEPTION (-20002, 'Raised E.')", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := mustScriptBlock(t, "DECLARE "+tt.decl+"; BEGIN RETURN 1; END")
			d := decl(t, block, 0)
			if d.Kind != ast.ScriptDeclException {
				t.Fatalf("kind = %v, want ScriptDeclException", d.Kind)
			}
			if len(d.ExcArgs) != tt.wantArg {
				t.Errorf("excArgs = %d, want %d", len(d.ExcArgs), tt.wantArg)
			}
		})
	}
}

func TestScript_DeclareMultiple(t *testing.T) {
	block := mustScriptBlock(t, `DECLARE
		x INT;
		y FLOAT DEFAULT 1.0;
		c CURSOR FOR SELECT 1;
		e EXCEPTION;
	BEGIN RETURN x; END`)
	if len(block.Decls) != 4 {
		t.Fatalf("decls = %d, want 4", len(block.Decls))
	}
	wantKinds := []ast.ScriptDeclKind{
		ast.ScriptDeclVar, ast.ScriptDeclVar, ast.ScriptDeclCursor, ast.ScriptDeclException,
	}
	for i, k := range wantKinds {
		if got := decl(t, block, i).Kind; got != k {
			t.Errorf("decl[%d] kind = %v, want %v", i, got, k)
		}
	}
}

// ---------------------------------------------------------------------------
// Assignment / LET
// ---------------------------------------------------------------------------

func TestScript_Assignment(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN x := 3; END")
	s, ok := firstBodyStmt(t, block).(*ast.ScriptAssignStmt)
	if !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptAssignStmt", block.Body[0])
	}
	if s.Target.Name != "x" {
		t.Errorf("target = %q, want x", s.Target.Normalize())
	}
	if s.Value == nil {
		t.Errorf("value = nil")
	}
}

func TestScript_AssignmentResultsetSubquery(t *testing.T) {
	// rs := (SELECT ...) — RESULTSET assignment via the plain assignment path.
	block := mustScriptBlock(t, "BEGIN rs := (SELECT 1); END")
	s, ok := firstBodyStmt(t, block).(*ast.ScriptAssignStmt)
	if !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptAssignStmt", block.Body[0])
	}
	if _, ok := s.Value.(*ast.SubqueryExpr); !ok {
		t.Errorf("value = %T, want *ast.SubqueryExpr", s.Value)
	}
}

func TestScript_LetForms(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		wantKind ast.ScriptDeclKind
		wantType bool
		wantQ    bool
	}{
		{"var typed assign", "LET x INT := 1", ast.ScriptDeclVar, true, false},
		{"var typed default", "LET x INT DEFAULT 1", ast.ScriptDeclVar, true, false},
		{"var inferred", "LET x := 1", ast.ScriptDeclVar, false, false},
		{"cursor", "LET c CURSOR FOR SELECT 1", ast.ScriptDeclCursor, false, true},
		{"resultset", "LET rs RESULTSET := (SELECT 1)", ast.ScriptDeclResultset, false, true},
		{"resultset async", "LET rs RESULTSET := ASYNC (SELECT 1)", ast.ScriptDeclResultset, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := mustScriptBlock(t, "BEGIN "+tt.stmt+"; END")
			s, ok := firstBodyStmt(t, block).(*ast.ScriptLetStmt)
			if !ok {
				t.Fatalf("body[0] = %T, want *ast.ScriptLetStmt", block.Body[0])
			}
			if s.Kind != tt.wantKind {
				t.Errorf("kind = %v, want %v", s.Kind, tt.wantKind)
			}
			if (s.Type != nil) != tt.wantType {
				t.Errorf("hasType = %v, want %v", s.Type != nil, tt.wantType)
			}
			if (s.Query != nil) != tt.wantQ {
				t.Errorf("hasQuery = %v, want %v", s.Query != nil, tt.wantQ)
			}
		})
	}
}

func TestScript_LetResultsetAsync(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN LET rs RESULTSET := ASYNC (SELECT 1); END")
	s := firstBodyStmt(t, block).(*ast.ScriptLetStmt)
	if !s.Async {
		t.Errorf("async = false, want true")
	}
}

// ---------------------------------------------------------------------------
// IF
// ---------------------------------------------------------------------------

func TestScript_IfThenEnd(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN IF (x = 1) THEN SELECT 1; END IF; END")
	s, ok := firstBodyStmt(t, block).(*ast.ScriptIfStmt)
	if !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptIfStmt", block.Body[0])
	}
	if len(s.Branches) != 1 {
		t.Fatalf("branches = %d, want 1", len(s.Branches))
	}
	b0 := ifBranch(t, s, 0)
	if b0.Cond == nil {
		t.Errorf("branch cond = nil")
	}
	if len(b0.Body) != 1 {
		t.Errorf("branch body = %d, want 1", len(b0.Body))
	}
	if s.Else != nil {
		t.Errorf("else = %v, want nil", s.Else)
	}
}

func TestScript_IfElseifElse(t *testing.T) {
	block := mustScriptBlock(t, `BEGIN
		IF (x = 1) THEN SELECT 1;
		ELSEIF (x = 2) THEN SELECT 2;
		ELSEIF (x = 3) THEN SELECT 3;
		ELSE SELECT 4;
		END IF;
	END`)
	s := firstBodyStmt(t, block).(*ast.ScriptIfStmt)
	if len(s.Branches) != 3 {
		t.Fatalf("branches = %d, want 3 (IF + 2 ELSEIF)", len(s.Branches))
	}
	if len(s.Else) != 1 {
		t.Errorf("else body = %d, want 1", len(s.Else))
	}
}

func TestScript_IfMultiStmtBranch(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN IF (x) THEN SELECT 1; SELECT 2; SELECT 3; END IF; END")
	s := firstBodyStmt(t, block).(*ast.ScriptIfStmt)
	if len(ifBranch(t, s, 0).Body) != 3 {
		t.Errorf("branch body = %d, want 3", len(ifBranch(t, s, 0).Body))
	}
}

// ---------------------------------------------------------------------------
// CASE
// ---------------------------------------------------------------------------

func TestScript_CaseSimple(t *testing.T) {
	block := mustScriptBlock(t, `BEGIN
		CASE (x)
			WHEN 1 THEN SELECT 'one';
			WHEN 2 THEN SELECT 'two';
			ELSE SELECT 'other';
		END CASE;
	END`)
	s, ok := firstBodyStmt(t, block).(*ast.ScriptCaseStmt)
	if !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptCaseStmt", block.Body[0])
	}
	if s.Operand == nil {
		t.Errorf("operand = nil, want simple-CASE operand")
	}
	if len(s.Whens) != 2 {
		t.Errorf("whens = %d, want 2", len(s.Whens))
	}
	if len(s.Else) != 1 {
		t.Errorf("else = %d, want 1", len(s.Else))
	}
}

func TestScript_CaseSearched(t *testing.T) {
	block := mustScriptBlock(t, `BEGIN
		CASE
			WHEN x = 1 THEN SELECT 1;
			WHEN x = 2 THEN SELECT 2;
		END;
	END`)
	s := firstBodyStmt(t, block).(*ast.ScriptCaseStmt)
	if s.Operand != nil {
		t.Errorf("operand = %v, want nil (searched form)", s.Operand)
	}
	if len(s.Whens) != 2 {
		t.Errorf("whens = %d, want 2", len(s.Whens))
	}
}

func TestScript_CaseEndWithoutCaseKeyword(t *testing.T) {
	// END (no CASE) is valid per docs: END [ CASE ].
	block := mustScriptBlock(t, "BEGIN CASE WHEN x THEN SELECT 1; END; END")
	if _, ok := firstBodyStmt(t, block).(*ast.ScriptCaseStmt); !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptCaseStmt", block.Body[0])
	}
}

// ---------------------------------------------------------------------------
// FOR
// ---------------------------------------------------------------------------

func TestScript_ForCounter(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN FOR i IN 1 TO 10 DO SELECT i; END FOR; END")
	s, ok := firstBodyStmt(t, block).(*ast.ScriptForStmt)
	if !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptForStmt", block.Body[0])
	}
	if s.Kind != ast.ScriptForCounter {
		t.Fatalf("kind = %v, want ScriptForCounter", s.Kind)
	}
	if s.Var.Name != "i" {
		t.Errorf("var = %q, want i", s.Var.Normalize())
	}
	if s.Start == nil || s.End == nil {
		t.Errorf("start/end = %v/%v, want non-nil", s.Start, s.End)
	}
	if s.Reverse {
		t.Errorf("reverse = true, want false")
	}
}

func TestScript_ForCounterReverse(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN FOR i IN REVERSE 1 TO 10 DO SELECT i; END FOR; END")
	s := firstBodyStmt(t, block).(*ast.ScriptForStmt)
	if !s.Reverse {
		t.Errorf("reverse = false, want true")
	}
}

func TestScript_ForCounterLoopKeyword(t *testing.T) {
	// Counter FOR may use LOOP/END LOOP instead of DO/END FOR.
	block := mustScriptBlock(t, "BEGIN FOR i IN 1 TO 3 LOOP SELECT i; END LOOP; END")
	s := firstBodyStmt(t, block).(*ast.ScriptForStmt)
	if s.Kind != ast.ScriptForCounter {
		t.Errorf("kind = %v, want ScriptForCounter", s.Kind)
	}
}

func TestScript_ForCursor(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN FOR rec IN c1 DO SELECT rec.x; END FOR; END")
	s, ok := firstBodyStmt(t, block).(*ast.ScriptForStmt)
	if !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptForStmt", block.Body[0])
	}
	if s.Kind != ast.ScriptForCursor {
		t.Fatalf("kind = %v, want ScriptForCursor", s.Kind)
	}
	if s.Source.Name != "c1" {
		t.Errorf("source = %q, want c1", s.Source.Normalize())
	}
}

func TestScript_ForLabel(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN FOR i IN 1 TO 3 DO SELECT i; END FOR loop_label; END")
	s := firstBodyStmt(t, block).(*ast.ScriptForStmt)
	if s.Label.Name != "loop_label" {
		t.Errorf("label = %q, want loop_label", s.Label.Normalize())
	}
}

// ---------------------------------------------------------------------------
// WHILE / REPEAT / LOOP
// ---------------------------------------------------------------------------

func TestScript_While(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN WHILE (x < 10) DO SELECT x; END WHILE; END")
	s, ok := firstBodyStmt(t, block).(*ast.ScriptWhileStmt)
	if !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptWhileStmt", block.Body[0])
	}
	if s.Cond == nil {
		t.Errorf("cond = nil")
	}
	if len(s.Body) != 1 {
		t.Errorf("body = %d, want 1", len(s.Body))
	}
}

func TestScript_WhileLoopKeyword(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN WHILE (x < 10) LOOP SELECT x; END LOOP; END")
	if _, ok := firstBodyStmt(t, block).(*ast.ScriptWhileStmt); !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptWhileStmt", block.Body[0])
	}
}

func TestScript_Repeat(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN REPEAT SELECT 1; UNTIL (x > 5) END REPEAT; END")
	s, ok := firstBodyStmt(t, block).(*ast.ScriptRepeatStmt)
	if !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptRepeatStmt", block.Body[0])
	}
	if s.Cond == nil {
		t.Errorf("until cond = nil")
	}
	if len(s.Body) != 1 {
		t.Errorf("body = %d, want 1", len(s.Body))
	}
}

func TestScript_Loop(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN LOOP SELECT 1; BREAK; END LOOP; END")
	s, ok := firstBodyStmt(t, block).(*ast.ScriptLoopStmt)
	if !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptLoopStmt", block.Body[0])
	}
	if len(s.Body) != 2 {
		t.Errorf("body = %d, want 2 (SELECT + BREAK)", len(s.Body))
	}
}

// ---------------------------------------------------------------------------
// Jump statements
// ---------------------------------------------------------------------------

func TestScript_BreakContinue(t *testing.T) {
	tests := []struct {
		name  string
		stmt  string
		check func(t *testing.T, n ast.Node)
	}{
		{"break", "BREAK", func(t *testing.T, n ast.Node) {
			if _, ok := n.(*ast.ScriptBreakStmt); !ok {
				t.Errorf("got %T, want *ast.ScriptBreakStmt", n)
			}
		}},
		{"break label", "BREAK lbl", func(t *testing.T, n ast.Node) {
			s := n.(*ast.ScriptBreakStmt)
			if s.Label.Name != "lbl" {
				t.Errorf("label = %q, want lbl", s.Label.Normalize())
			}
		}},
		{"exit alias", "EXIT", func(t *testing.T, n ast.Node) {
			if _, ok := n.(*ast.ScriptBreakStmt); !ok {
				t.Errorf("got %T, want *ast.ScriptBreakStmt", n)
			}
		}},
		{"continue", "CONTINUE", func(t *testing.T, n ast.Node) {
			if _, ok := n.(*ast.ScriptContinueStmt); !ok {
				t.Errorf("got %T, want *ast.ScriptContinueStmt", n)
			}
		}},
		{"continue label", "CONTINUE lbl", func(t *testing.T, n ast.Node) {
			s := n.(*ast.ScriptContinueStmt)
			if s.Label.Name != "lbl" {
				t.Errorf("label = %q, want lbl", s.Label.Normalize())
			}
		}},
		{"iterate alias", "ITERATE", func(t *testing.T, n ast.Node) {
			if _, ok := n.(*ast.ScriptContinueStmt); !ok {
				t.Errorf("got %T, want *ast.ScriptContinueStmt", n)
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := mustScriptBlock(t, "BEGIN LOOP "+tt.stmt+"; END LOOP; END")
			loop := firstBodyStmt(t, block).(*ast.ScriptLoopStmt)
			tt.check(t, loop.Body[0])
		})
	}
}

func TestScript_Return(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		wantExpr bool
	}{
		{"bare", "RETURN", false},
		{"with expr", "RETURN x + 1", true},
		{"with func call", "RETURN pi() * 2", true},
		{"with subquery", "RETURN (SELECT count(*) FROM t)", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := mustScriptBlock(t, "BEGIN "+tt.stmt+"; END")
			s, ok := firstBodyStmt(t, block).(*ast.ScriptReturnStmt)
			if !ok {
				t.Fatalf("body[0] = %T, want *ast.ScriptReturnStmt", block.Body[0])
			}
			if (s.Value != nil) != tt.wantExpr {
				t.Errorf("hasValue = %v, want %v", s.Value != nil, tt.wantExpr)
			}
		})
	}
}

func TestScript_Raise(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		wantName string
	}{
		{"bare reraise", "RAISE", ""},
		{"named", "RAISE my_exception", "my_exception"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := mustScriptBlock(t, "BEGIN "+tt.stmt+"; END")
			s, ok := firstBodyStmt(t, block).(*ast.ScriptRaiseStmt)
			if !ok {
				t.Fatalf("body[0] = %T, want *ast.ScriptRaiseStmt", block.Body[0])
			}
			if s.Name.Name != tt.wantName {
				t.Errorf("name = %q, want %q", s.Name.Name, tt.wantName)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Cursor statements
// ---------------------------------------------------------------------------

func TestScript_OpenFetchClose(t *testing.T) {
	block := mustScriptBlock(t, `BEGIN
		OPEN c1;
		FETCH c1 INTO v1;
		CLOSE c1;
	END`)
	if len(block.Body) != 3 {
		t.Fatalf("body = %d, want 3", len(block.Body))
	}
	open, ok := block.Body[0].(*ast.ScriptOpenStmt)
	if !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptOpenStmt", block.Body[0])
	}
	if open.Cursor.Name != "c1" {
		t.Errorf("open cursor = %q, want c1", open.Cursor.Normalize())
	}
	if open.Using != nil {
		t.Errorf("open using = %v, want nil", open.Using)
	}
	fetch, ok := block.Body[1].(*ast.ScriptFetchStmt)
	if !ok {
		t.Fatalf("body[1] = %T, want *ast.ScriptFetchStmt", block.Body[1])
	}
	if len(fetch.Into) != 1 || fetch.Into[0].Name != "v1" {
		t.Errorf("fetch into = %v, want [v1]", fetch.Into)
	}
	if _, ok := block.Body[2].(*ast.ScriptCloseStmt); !ok {
		t.Fatalf("body[2] = %T, want *ast.ScriptCloseStmt", block.Body[2])
	}
}

func TestScript_OpenUsing(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN OPEN c1 USING (min_price, max_price); END")
	open := firstBodyStmt(t, block).(*ast.ScriptOpenStmt)
	if len(open.Using) != 2 {
		t.Errorf("using = %d, want 2", len(open.Using))
	}
}

func TestScript_FetchMultipleInto(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN FETCH c1 INTO a, b, c; END")
	fetch := firstBodyStmt(t, block).(*ast.ScriptFetchStmt)
	if len(fetch.Into) != 3 {
		t.Errorf("into = %d, want 3", len(fetch.Into))
	}
}

// ---------------------------------------------------------------------------
// EXCEPTION handlers
// ---------------------------------------------------------------------------

func TestScript_ExceptionHandlers(t *testing.T) {
	block := mustScriptBlock(t, `BEGIN
		SELECT 1;
	EXCEPTION
		WHEN my_exc THEN SELECT 2;
		WHEN OTHER THEN SELECT 3;
	END`)
	if len(block.Handlers) != 2 {
		t.Fatalf("handlers = %d, want 2", len(block.Handlers))
	}
	h0 := handler(t, block, 0)
	if h0.Other {
		t.Errorf("handler[0].Other = true, want false")
	}
	if len(h0.Names) != 1 || h0.Names[0].Name != "my_exc" {
		t.Errorf("handler[0].Names = %v, want [my_exc]", h0.Names)
	}
	if !handler(t, block, 1).Other {
		t.Errorf("handler[1].Other = false, want true (WHEN OTHER)")
	}
}

func TestScript_ExceptionHandlerOrNames(t *testing.T) {
	block := mustScriptBlock(t, `BEGIN
		SELECT 1;
	EXCEPTION
		WHEN exc_a OR exc_b OR exc_c THEN SELECT 2;
	END`)
	if len(block.Handlers) != 1 {
		t.Fatalf("handlers = %d, want 1", len(block.Handlers))
	}
	if len(handler(t, block, 0).Names) != 3 {
		t.Errorf("names = %d, want 3 (OR-separated)", len(handler(t, block, 0).Names))
	}
}

// ---------------------------------------------------------------------------
// Nested SQL inside a block (reuse parseStmt)
// ---------------------------------------------------------------------------

func TestScript_NestedSQLStatements(t *testing.T) {
	block := mustScriptBlock(t, `BEGIN
		CREATE TEMP TABLE t (id INT);
		INSERT INTO t VALUES (1);
		UPDATE t SET id = 2 WHERE id = 1;
		DELETE FROM t WHERE id = 2;
		MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET t.id = s.id;
	END`)
	if len(block.Body) != 5 {
		t.Fatalf("body = %d, want 5", len(block.Body))
	}
}

func TestScript_ExecuteImmediateInBlock(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN EXECUTE IMMEDIATE 'SELECT 1'; END")
	if _, ok := firstBodyStmt(t, block).(*ast.ExecuteImmediateStmt); !ok {
		t.Fatalf("body[0] = %T, want *ast.ExecuteImmediateStmt", block.Body[0])
	}
}

// ---------------------------------------------------------------------------
// Real-world example from the legacy corpus (task body, area-of-circle)
// ---------------------------------------------------------------------------

func TestScript_LegacyCorpusAreaOfCircle(t *testing.T) {
	block := mustScriptBlock(t, `DECLARE
		radius_of_circle float;
		area_of_circle float;
	BEGIN
		radius_of_circle := 3;
		area_of_circle := pi() * radius_of_circle * radius_of_circle;
		return area_of_circle;
	END`)
	if len(block.Decls) != 2 {
		t.Errorf("decls = %d, want 2", len(block.Decls))
	}
	if len(block.Body) != 3 {
		t.Errorf("body = %d, want 3 (2 assigns + return)", len(block.Body))
	}
	if _, ok := block.Body[0].(*ast.ScriptAssignStmt); !ok {
		t.Errorf("body[0] = %T, want *ast.ScriptAssignStmt", block.Body[0])
	}
	if _, ok := block.Body[2].(*ast.ScriptReturnStmt); !ok {
		t.Errorf("body[2] = %T, want *ast.ScriptReturnStmt", block.Body[2])
	}
}

// ---------------------------------------------------------------------------
// Deeply nested control flow (exercises depth + segmentation end-to-end)
// ---------------------------------------------------------------------------

func TestScript_DeeplyNestedControlFlow(t *testing.T) {
	block := mustScriptBlock(t, `BEGIN
		FOR i IN 1 TO 3 DO
			IF (i = 2) THEN
				WHILE (x > 0) DO
					CASE
						WHEN x = 1 THEN BREAK;
						ELSE CONTINUE;
					END CASE;
				END WHILE;
			END IF;
		END FOR;
	END`)
	forStmt, ok := firstBodyStmt(t, block).(*ast.ScriptForStmt)
	if !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptForStmt", block.Body[0])
	}
	ifStmt, ok := forStmt.Body[0].(*ast.ScriptIfStmt)
	if !ok {
		t.Fatalf("for body[0] = %T, want *ast.ScriptIfStmt", forStmt.Body[0])
	}
	thenBody := ifBranch(t, ifStmt, 0).Body
	whileStmt, ok := thenBody[0].(*ast.ScriptWhileStmt)
	if !ok {
		t.Fatalf("if body[0] = %T, want *ast.ScriptWhileStmt", thenBody[0])
	}
	if _, ok := whileStmt.Body[0].(*ast.ScriptCaseStmt); !ok {
		t.Fatalf("while body[0] = %T, want *ast.ScriptCaseStmt", whileStmt.Body[0])
	}
}

// ---------------------------------------------------------------------------
// Negatives — every malformed form must ERROR (and never hang)
//
// These are the loop-safety guarantees: an unterminated block / construct must
// fail fast with a ParseError rather than spin. The Go test runner's package
// timeout (60s) would catch a hang; these complete in microseconds.
// ---------------------------------------------------------------------------

func TestScript_Negatives(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unterminated block (no END)", "BEGIN SELECT 1;"},
		{"unterminated declare (no BEGIN)", "DECLARE x INT;"},
		{"unterminated if (no END IF)", "BEGIN IF (x) THEN SELECT 1; END"},
		{"if missing THEN", "BEGIN IF (x) SELECT 1; END IF; END"},
		{"if missing parens", "BEGIN IF x THEN SELECT 1; END IF; END"},
		{"unterminated for", "BEGIN FOR i IN 1 TO 3 DO SELECT i; END"},
		{"for missing TO", "BEGIN FOR i IN 1 DO SELECT i; END FOR; END"},
		{"unterminated while", "BEGIN WHILE (x) DO SELECT 1; END"},
		{"unterminated loop", "BEGIN LOOP SELECT 1; END"},
		{"repeat missing until", "BEGIN REPEAT SELECT 1; END REPEAT; END"},
		{"case no when", "BEGIN CASE END CASE; END"},
		{"exception no when", "BEGIN SELECT 1; EXCEPTION END"},
		{"fetch missing into", "BEGIN FETCH c1; END"},
		{"let missing value", "BEGIN LET x INT; END"},
		{"cursor missing for", "DECLARE c CURSOR SELECT 1; BEGIN RETURN 1; END"},
		{"declare missing semicolon", "DECLARE x INT y FLOAT; BEGIN RETURN 1; END"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := parseScriptErr(tt.input)
			if len(errs) == 0 {
				t.Errorf("%q: expected a parse error, got none", tt.input)
			}
		})
	}
}

// TestScript_UnterminatedDoesNotHang specifically exercises the most dangerous
// non-termination shapes (never-closed nested blocks) and asserts they return
// quickly with an error. If the loop guards regressed, this would hang and the
// package timeout would fire.
func TestScript_UnterminatedDoesNotHang(t *testing.T) {
	// Note: a bare "BEGIN" (BEGIN followed by EOF) is a VALID TCL transaction
	// opener, not an unterminated scripting block — so it is intentionally NOT
	// listed here. Every input below opens a scripting block (BEGIN followed by
	// a statement opener, or DECLARE) that is never closed.
	inputs := []string{
		"BEGIN SELECT 1",
		"BEGIN BEGIN BEGIN",
		"BEGIN IF (x) THEN",
		"BEGIN FOR i IN 1 TO 3 DO",
		"BEGIN LOOP",
		"BEGIN WHILE (x) DO",
		"BEGIN REPEAT",
		"DECLARE x INT",
		strings.Repeat("BEGIN SELECT 1; ", 50), // 50 unterminated nested-ish blocks
	}
	for _, in := range inputs {
		errs := parseScriptErr(in)
		if len(errs) == 0 {
			t.Errorf("%q: expected an error (unterminated), got none", in)
		}
	}
}

// TestScript_DepthBound asserts that pathologically deep nesting fails with a
// "too deeply" error rather than overflowing the stack. The depth ceiling is
// maxScriptDepth (200); 1000 nested BEGINs is comfortably past it.
func TestScript_DepthBound(t *testing.T) {
	deep := strings.Repeat("BEGIN ", 1000) + strings.Repeat("END ", 1000)
	errs := parseScriptErr(deep)
	if len(errs) == 0 {
		t.Fatal("expected a depth-bound error for 1000-deep nesting, got none")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Msg, "nested too deeply") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected a 'nested too deeply' error; got: %v", errs)
	}
}

// ---------------------------------------------------------------------------
// Walker integration — the block and its body statements are reachable.
// ---------------------------------------------------------------------------

func TestScript_WalkerVisitsBody(t *testing.T) {
	block := mustScriptBlock(t, "BEGIN SELECT 1; SELECT 2; END")
	var selects int
	ast.Inspect(block, func(n ast.Node) bool {
		if _, ok := n.(*ast.SelectStmt); ok {
			selects++
		}
		return true
	})
	if selects != 2 {
		t.Errorf("walker visited %d SelectStmt, want 2", selects)
	}
}

// ---------------------------------------------------------------------------
// Case-insensitivity (Snowflake keywords are case-insensitive)
// ---------------------------------------------------------------------------

func TestScript_CaseInsensitiveKeywords(t *testing.T) {
	block := mustScriptBlock(t, "begin if (x = 1) then select 1; end if; end")
	if _, ok := firstBodyStmt(t, block).(*ast.ScriptIfStmt); !ok {
		t.Fatalf("body[0] = %T, want *ast.ScriptIfStmt", block.Body[0])
	}
}

// TestScript_WalkerReachesNestedSelect confirms the generated walker descends
// through control-flow body slices (FOR -> IF -> SELECT), since those are
// []Node fields that genwalker walks.
func TestScript_WalkerReachesNestedSelect(t *testing.T) {
	block := mustScriptBlock(t, `BEGIN
		FOR i IN 1 TO 3 DO
			IF (i = 1) THEN SELECT 100; END IF;
		END FOR;
	END`)
	var found bool
	ast.Inspect(block, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectStmt); ok {
			_ = sel
			found = true
		}
		return true
	})
	if !found {
		t.Errorf("walker did not reach the nested SELECT inside FOR/IF body")
	}
}

// TestScript_AssignExprForms checks assignment RHS reuses the full expression
// grammar (function calls, casts, semi-structured access, arithmetic).
func TestScript_AssignExprForms(t *testing.T) {
	exprs := []string{
		"x := pi() * 2",
		"x := y::INT",
		"x := obj:field",
		"x := CASE WHEN y > 0 THEN 1 ELSE 0 END",
		"x := (SELECT max(c) FROM t)",
	}
	for _, e := range exprs {
		t.Run(e, func(t *testing.T) {
			block := mustScriptBlock(t, "BEGIN "+e+"; END")
			s, ok := firstBodyStmt(t, block).(*ast.ScriptAssignStmt)
			if !ok {
				t.Fatalf("body[0] = %T, want *ast.ScriptAssignStmt", block.Body[0])
			}
			if s.Value == nil {
				t.Errorf("value = nil for %q", e)
			}
		})
	}
}

// TestScript_FullFormDeclareBodyException exercises a complete block with all
// three sections: DECLARE, BEGIN body, and EXCEPTION.
func TestScript_FullFormDeclareBodyException(t *testing.T) {
	block := mustScriptBlock(t, `DECLARE
		result INT DEFAULT 0;
		my_exc EXCEPTION (-20001, 'boom');
	BEGIN
		result := 1;
		IF (result = 1) THEN
			RAISE my_exc;
		END IF;
		RETURN result;
	EXCEPTION
		WHEN my_exc THEN
			RETURN -1;
		WHEN OTHER THEN
			RETURN -2;
	END`)
	if len(block.Decls) != 2 {
		t.Errorf("decls = %d, want 2", len(block.Decls))
	}
	if len(block.Body) != 3 {
		t.Errorf("body = %d, want 3 (assign, if, return)", len(block.Body))
	}
	if len(block.Handlers) != 2 {
		t.Errorf("handlers = %d, want 2", len(block.Handlers))
	}
}

// TestScript_BeginTransactionStillTCL is a guard that the disambiguation did
// NOT capture the TCL BEGIN forms into the scripting parser.
func TestScript_BeginTransactionStillTCL(t *testing.T) {
	tcl := []string{"BEGIN", "BEGIN;", "BEGIN WORK", "BEGIN TRANSACTION", "BEGIN NAME t1"}
	for _, in := range tcl {
		t.Run(in, func(t *testing.T) {
			result := ParseBestEffort(in)
			if len(result.Errors) != 0 {
				t.Fatalf("%q: unexpected errors %v", in, result.Errors)
			}
			if len(result.File.Stmts) == 0 {
				t.Fatalf("%q: no statement parsed", in)
			}
			if _, ok := result.File.Stmts[0].(*ast.BeginStmt); !ok {
				t.Errorf("%q: got %T, want *ast.BeginStmt (TCL)", in, result.File.Stmts[0])
			}
		})
	}
}

// TestScript_BlockThenStatementSegmentation locks in that a scripting block is
// kept as ONE segment by Split (even when its loop headers embed CASE/IF
// expressions, or it contains bare-END CASE), and that a statement following
// the block is NOT absorbed into it. This guards the Split construct-depth
// bookkeeping (block extent) end-to-end through ParseBestEffort.
func TestScript_BlockThenStatementSegmentation(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantStmts int
	}{
		{
			"case-in-while-header then trailing select",
			"BEGIN WHILE (CASE WHEN x THEN 1 ELSE 0 END > 0) LOOP SELECT 1; END LOOP; END; SELECT 999",
			2,
		},
		{
			"bare-end case in block then trailing select",
			"BEGIN CASE WHEN x THEN SELECT 1; END; END; SELECT 2",
			2,
		},
		{
			"nested standalone loop inside while-do then trailing select",
			"BEGIN WHILE (x) DO LOOP SELECT 1; END LOOP; END WHILE; END; SELECT 7",
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseBestEffort(tt.input)
			if len(result.Errors) != 0 {
				t.Fatalf("%q: unexpected errors %v", tt.input, result.Errors)
			}
			if len(result.File.Stmts) != tt.wantStmts {
				t.Fatalf("%q: got %d statements, want %d", tt.input, len(result.File.Stmts), tt.wantStmts)
			}
			// The first statement must be the scripting block; the last must be
			// the trailing SELECT (not absorbed).
			if _, ok := result.File.Stmts[0].(*ast.ScriptBlockStmt); !ok {
				t.Errorf("stmt[0] = %T, want *ast.ScriptBlockStmt", result.File.Stmts[0])
			}
			last := result.File.Stmts[len(result.File.Stmts)-1]
			if _, ok := last.(*ast.SelectStmt); !ok {
				t.Errorf("last stmt = %T, want *ast.SelectStmt (trailing, not absorbed)", last)
			}
		})
	}
}
