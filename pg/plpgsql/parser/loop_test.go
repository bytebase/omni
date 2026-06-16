package parser

import (
	"testing"

	"github.com/bytebase/omni/pg/plpgsql/ast"
)

// --------------------------------------------------------------------------
// Section 2.3: Basic Loops — LOOP, WHILE
// --------------------------------------------------------------------------

func TestLoopInfinite(t *testing.T) {
	// Infinite LOOP: LOOP x := x + 1; EXIT WHEN x > 10; END LOOP;
	block := parseOK(t, `BEGIN LOOP x := x + 1; EXIT WHEN x > 10; END LOOP; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	loop, ok := block.Body[0].(*ast.PLLoop)
	if !ok {
		t.Fatalf("expected PLLoop, got %T", block.Body[0])
	}
	if loop.Label != "" {
		t.Errorf("expected empty label, got %q", loop.Label)
	}
	if len(loop.Body) != 2 {
		t.Errorf("expected 2 statements in loop body, got %d", len(loop.Body))
	}
}

func TestLoopLabeled(t *testing.T) {
	// Labeled LOOP: <<myloop>> LOOP EXIT myloop; END LOOP myloop;
	block := parseOK(t, `BEGIN <<myloop>> LOOP EXIT myloop; END LOOP myloop; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	loop, ok := block.Body[0].(*ast.PLLoop)
	if !ok {
		t.Fatalf("expected PLLoop, got %T", block.Body[0])
	}
	if loop.Label != "myloop" {
		t.Errorf("expected label 'myloop', got %q", loop.Label)
	}
}

func TestWhileLoop(t *testing.T) {
	// WHILE loop: WHILE x > 0 LOOP x := x - 1; END LOOP;
	block := parseOK(t, `BEGIN WHILE x > 0 LOOP x := x - 1; END LOOP; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	w, ok := block.Body[0].(*ast.PLWhile)
	if !ok {
		t.Fatalf("expected PLWhile, got %T", block.Body[0])
	}
	if w.Label != "" {
		t.Errorf("expected empty label, got %q", w.Label)
	}
	if w.Condition != "x > 0" {
		t.Errorf("expected condition 'x > 0', got %q", w.Condition)
	}
	if len(w.Body) != 1 {
		t.Errorf("expected 1 statement in while body, got %d", len(w.Body))
	}
}

func TestWhileLabeledLoop(t *testing.T) {
	// Labeled WHILE: <<w>> WHILE x > 0 LOOP x := x - 1; END LOOP w;
	block := parseOK(t, `BEGIN <<w>> WHILE x > 0 LOOP x := x - 1; END LOOP w; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	w, ok := block.Body[0].(*ast.PLWhile)
	if !ok {
		t.Fatalf("expected PLWhile, got %T", block.Body[0])
	}
	if w.Label != "w" {
		t.Errorf("expected label 'w', got %q", w.Label)
	}
	if w.Condition != "x > 0" {
		t.Errorf("expected condition 'x > 0', got %q", w.Condition)
	}
}

func TestWhileConditionSpansUntilLoop(t *testing.T) {
	// Condition expression spans until LOOP keyword
	block := parseOK(t, `BEGIN WHILE x > 0 AND y < 100 LOOP NULL; END LOOP; END`)
	w := block.Body[0].(*ast.PLWhile)
	if w.Condition != "x > 0 AND y < 100" {
		t.Errorf("expected condition 'x > 0 AND y < 100', got %q", w.Condition)
	}
}

func TestLoopNested(t *testing.T) {
	// Nested loops: LOOP inside LOOP
	block := parseOK(t, `BEGIN LOOP LOOP EXIT; END LOOP; EXIT; END LOOP; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	outer, ok := block.Body[0].(*ast.PLLoop)
	if !ok {
		t.Fatalf("expected PLLoop, got %T", block.Body[0])
	}
	if len(outer.Body) != 2 {
		t.Fatalf("expected 2 statements in outer loop, got %d", len(outer.Body))
	}
	inner, ok := outer.Body[0].(*ast.PLLoop)
	if !ok {
		t.Fatalf("expected inner PLLoop, got %T", outer.Body[0])
	}
	if len(inner.Body) != 1 {
		t.Errorf("expected 1 statement in inner loop, got %d", len(inner.Body))
	}
}

func TestLoopEmptyBody(t *testing.T) {
	// Empty loop body: LOOP END LOOP;
	block := parseOK(t, `BEGIN LOOP END LOOP; END`)
	loop := block.Body[0].(*ast.PLLoop)
	if len(loop.Body) != 0 {
		t.Errorf("expected 0 statements in loop body, got %d", len(loop.Body))
	}
}

func TestLoopMissingEndLoop(t *testing.T) {
	// Missing END LOOP produces error
	parseErr(t, `BEGIN LOOP EXIT; END; END`, "expected LOOP after END")
}

func TestWhileMissingLoopKeyword(t *testing.T) {
	// WHILE missing LOOP keyword produces error
	parseErr(t, `BEGIN WHILE x > 0 END LOOP; END`, "expected LOOP")
}

// --------------------------------------------------------------------------
// Section 2.4: FOR Loops — All Variants
// --------------------------------------------------------------------------

func TestForInteger(t *testing.T) {
	// Integer FOR: FOR i IN 1..10 LOOP x := x + i; END LOOP;
	block := parseOK(t, `BEGIN FOR i IN 1..10 LOOP x := x + i; END LOOP; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	f, ok := block.Body[0].(*ast.PLForI)
	if !ok {
		t.Fatalf("expected PLForI, got %T", block.Body[0])
	}
	if f.Var != "i" {
		t.Errorf("expected var 'i', got %q", f.Var)
	}
	if f.Lower != "1" {
		t.Errorf("expected lower '1', got %q", f.Lower)
	}
	if f.Upper != "10" {
		t.Errorf("expected upper '10', got %q", f.Upper)
	}
	if f.Reverse {
		t.Error("expected reverse=false")
	}
	if f.Step != "" {
		t.Errorf("expected empty step, got %q", f.Step)
	}
	if len(f.Body) != 1 {
		t.Errorf("expected 1 statement in body, got %d", len(f.Body))
	}
}

func TestForIntegerReverse(t *testing.T) {
	// Integer FOR REVERSE: FOR i IN REVERSE 10..1 LOOP x := x + i; END LOOP;
	block := parseOK(t, `BEGIN FOR i IN REVERSE 10..1 LOOP x := x + i; END LOOP; END`)
	f := block.Body[0].(*ast.PLForI)
	if f.Var != "i" {
		t.Errorf("expected var 'i', got %q", f.Var)
	}
	if f.Lower != "10" {
		t.Errorf("expected lower '10', got %q", f.Lower)
	}
	if f.Upper != "1" {
		t.Errorf("expected upper '1', got %q", f.Upper)
	}
	if !f.Reverse {
		t.Error("expected reverse=true")
	}
}

func TestForIntegerWithBy(t *testing.T) {
	// Integer FOR with BY: FOR i IN 1..100 BY 5 LOOP x := x + i; END LOOP;
	block := parseOK(t, `BEGIN FOR i IN 1..100 BY 5 LOOP x := x + i; END LOOP; END`)
	f := block.Body[0].(*ast.PLForI)
	if f.Lower != "1" {
		t.Errorf("expected lower '1', got %q", f.Lower)
	}
	if f.Upper != "100" {
		t.Errorf("expected upper '100', got %q", f.Upper)
	}
	if f.Step != "5" {
		t.Errorf("expected step '5', got %q", f.Step)
	}
}

func TestForIntegerReverseWithBy(t *testing.T) {
	// Integer FOR REVERSE with BY: FOR i IN REVERSE 100..1 BY 10 LOOP END LOOP;
	block := parseOK(t, `BEGIN FOR i IN REVERSE 100..1 BY 10 LOOP END LOOP; END`)
	f := block.Body[0].(*ast.PLForI)
	if f.Lower != "100" {
		t.Errorf("expected lower '100', got %q", f.Lower)
	}
	if f.Upper != "1" {
		t.Errorf("expected upper '1', got %q", f.Upper)
	}
	if f.Step != "10" {
		t.Errorf("expected step '10', got %q", f.Step)
	}
	if !f.Reverse {
		t.Error("expected reverse=true")
	}
}

func TestForIntegerExprBounds(t *testing.T) {
	// Integer FOR with expression bounds: FOR i IN a+1..b*2 LOOP END LOOP;
	block := parseOK(t, `BEGIN FOR i IN a+1..b*2 LOOP END LOOP; END`)
	f := block.Body[0].(*ast.PLForI)
	if f.Lower != "a+1" {
		t.Errorf("expected lower 'a+1', got %q", f.Lower)
	}
	if f.Upper != "b*2" {
		t.Errorf("expected upper 'b*2', got %q", f.Upper)
	}
}

func TestForIntegerRejectsTargetList(t *testing.T) {
	parseErr(t, `BEGIN FOR a, b IN 1..10 LOOP NULL; END LOOP; END`, "expected single loop variable")
}

func TestForQuery(t *testing.T) {
	// Query FOR: FOR rec IN SELECT * FROM t LOOP x := rec.a; END LOOP;
	block := parseOK(t, `BEGIN FOR rec IN SELECT * FROM t LOOP x := rec.a; END LOOP; END`)
	f, ok := block.Body[0].(*ast.PLForS)
	if !ok {
		t.Fatalf("expected PLForS, got %T", block.Body[0])
	}
	if f.Var != "rec" {
		t.Errorf("expected var 'rec', got %q", f.Var)
	}
	if f.Query != "SELECT * FROM t" {
		t.Errorf("expected query 'SELECT * FROM t', got %q", f.Query)
	}
	if len(f.Body) != 1 {
		t.Errorf("expected 1 statement in body, got %d", len(f.Body))
	}
}

func TestForQueryComplex(t *testing.T) {
	// Query FOR with complex query
	block := parseOK(t, `BEGIN FOR rec IN SELECT a, b FROM t WHERE c > 0 ORDER BY a LOOP END LOOP; END`)
	f := block.Body[0].(*ast.PLForS)
	if f.Query != "SELECT a, b FROM t WHERE c > 0 ORDER BY a" {
		t.Errorf("unexpected query %q", f.Query)
	}
}

func TestForQueryTargetList(t *testing.T) {
	block := parseOK(t, `BEGIN FOR a, b IN SELECT x, y FROM t LOOP NULL; END LOOP; END`)
	f, ok := block.Body[0].(*ast.PLForS)
	if !ok {
		t.Fatalf("expected PLForS, got %T", block.Body[0])
	}
	if f.Var != "a, b" {
		t.Errorf("expected var 'a, b', got %q", f.Var)
	}
	if f.Query != "SELECT x, y FROM t" {
		t.Errorf("expected query 'SELECT x, y FROM t', got %q", f.Query)
	}
}

func TestForQueryParenthesizedTargetList(t *testing.T) {
	block := parseOK(t, `BEGIN FOR a, b IN (SELECT x, y FROM t) LOOP NULL; END LOOP; END`)
	f, ok := block.Body[0].(*ast.PLForS)
	if !ok {
		t.Fatalf("expected PLForS, got %T", block.Body[0])
	}
	if f.Var != "a, b" {
		t.Errorf("expected var 'a, b', got %q", f.Var)
	}
	if f.Query != "(SELECT x, y FROM t)" {
		t.Errorf("expected query '(SELECT x, y FROM t)', got %q", f.Query)
	}
}

func TestForCursor(t *testing.T) {
	// Cursor FOR: FOR rec IN cur LOOP END LOOP;
	block := parseOK(t, `BEGIN FOR rec IN cur LOOP END LOOP; END`)
	f, ok := block.Body[0].(*ast.PLForC)
	if !ok {
		t.Fatalf("expected PLForC, got %T", block.Body[0])
	}
	if f.Var != "rec" {
		t.Errorf("expected var 'rec', got %q", f.Var)
	}
	if f.CursorVar != "cur" {
		t.Errorf("expected cursor var 'cur', got %q", f.CursorVar)
	}
	if f.ArgQuery != "" {
		t.Errorf("expected empty arg query, got %q", f.ArgQuery)
	}
}

func TestForCursorWithArgs(t *testing.T) {
	// Cursor FOR with arguments: FOR rec IN cur(1, 'a') LOOP END LOOP;
	block := parseOK(t, `BEGIN FOR rec IN cur(1, 'a') LOOP END LOOP; END`)
	f, ok := block.Body[0].(*ast.PLForC)
	if !ok {
		t.Fatalf("expected PLForC, got %T", block.Body[0])
	}
	if f.CursorVar != "cur" {
		t.Errorf("expected cursor var 'cur', got %q", f.CursorVar)
	}
	if f.ArgQuery != "1, 'a'" {
		t.Errorf("expected arg query \"1, 'a'\", got %q", f.ArgQuery)
	}
}

func TestForCursorRejectsTargetList(t *testing.T) {
	parseErr(t, `BEGIN FOR a, b IN cur LOOP NULL; END LOOP; END`, "expected single loop variable")
}

func TestForDynamic(t *testing.T) {
	// Dynamic FOR: FOR rec IN EXECUTE 'SELECT * FROM ' || tbl LOOP END LOOP;
	block := parseOK(t, `BEGIN FOR rec IN EXECUTE 'SELECT * FROM ' || tbl LOOP END LOOP; END`)
	f, ok := block.Body[0].(*ast.PLForDynS)
	if !ok {
		t.Fatalf("expected PLForDynS, got %T", block.Body[0])
	}
	if f.Var != "rec" {
		t.Errorf("expected var 'rec', got %q", f.Var)
	}
	if f.Query != "'SELECT * FROM ' || tbl" {
		t.Errorf("expected query \"'SELECT * FROM ' || tbl\", got %q", f.Query)
	}
	if f.Params != nil {
		t.Errorf("expected nil params, got %v", f.Params)
	}
}

func TestForDynamicWithUsing(t *testing.T) {
	// Dynamic FOR with USING
	block := parseOK(t, `BEGIN FOR rec IN EXECUTE 'SELECT * FROM t WHERE id = $1' USING my_id LOOP END LOOP; END`)
	f := block.Body[0].(*ast.PLForDynS)
	if f.Query != "'SELECT * FROM t WHERE id = $1'" {
		t.Errorf("unexpected query %q", f.Query)
	}
	if len(f.Params) != 1 || f.Params[0] != "my_id" {
		t.Errorf("expected params [my_id], got %v", f.Params)
	}
}

func TestForEachArray(t *testing.T) {
	// FOREACH ARRAY: FOREACH x IN ARRAY arr LOOP END LOOP;
	block := parseOK(t, `BEGIN FOREACH x IN ARRAY arr LOOP END LOOP; END`)
	f, ok := block.Body[0].(*ast.PLForEachA)
	if !ok {
		t.Fatalf("expected PLForEachA, got %T", block.Body[0])
	}
	if f.Var != "x" {
		t.Errorf("expected var 'x', got %q", f.Var)
	}
	if f.SliceDim != 0 {
		t.Errorf("expected slice dim 0, got %d", f.SliceDim)
	}
	if f.ArrayExpr != "arr" {
		t.Errorf("expected array expr 'arr', got %q", f.ArrayExpr)
	}
}

func TestForEachArrayWithSlice(t *testing.T) {
	// FOREACH ARRAY with SLICE: FOREACH x SLICE 1 IN ARRAY arr LOOP END LOOP;
	block := parseOK(t, `BEGIN FOREACH x SLICE 1 IN ARRAY arr LOOP END LOOP; END`)
	f := block.Body[0].(*ast.PLForEachA)
	if f.SliceDim != 1 {
		t.Errorf("expected slice dim 1, got %d", f.SliceDim)
	}
	if f.ArrayExpr != "arr" {
		t.Errorf("expected array expr 'arr', got %q", f.ArrayExpr)
	}
}

func TestForLabeled(t *testing.T) {
	// Labeled FOR: <<fl>> FOR i IN 1..10 LOOP END LOOP fl;
	block := parseOK(t, `BEGIN <<fl>> FOR i IN 1..10 LOOP END LOOP fl; END`)
	f, ok := block.Body[0].(*ast.PLForI)
	if !ok {
		t.Fatalf("expected PLForI, got %T", block.Body[0])
	}
	if f.Label != "fl" {
		t.Errorf("expected label 'fl', got %q", f.Label)
	}
}

func TestForLabelMismatch(t *testing.T) {
	// Label mismatch on FOR loop END produces error
	parseErr(t, `BEGIN <<a>> FOR i IN 1..10 LOOP END LOOP b; END`, "does not match")
}

func TestForNested(t *testing.T) {
	// Nested FOR loops
	block := parseOK(t, `BEGIN FOR i IN 1..10 LOOP FOR j IN 1..5 LOOP NULL; END LOOP; END LOOP; END`)
	outer, ok := block.Body[0].(*ast.PLForI)
	if !ok {
		t.Fatalf("expected PLForI, got %T", block.Body[0])
	}
	if len(outer.Body) != 1 {
		t.Fatalf("expected 1 statement in outer body, got %d", len(outer.Body))
	}
	inner, ok := outer.Body[0].(*ast.PLForI)
	if !ok {
		t.Fatalf("expected inner PLForI, got %T", outer.Body[0])
	}
	if inner.Var != "j" {
		t.Errorf("expected inner var 'j', got %q", inner.Var)
	}
}
