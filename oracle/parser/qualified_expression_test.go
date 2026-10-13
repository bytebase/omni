package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// qualifiedDecls declares the types the qualified expressions below build.
const qualifiedDecls = "TYPE r IS RECORD (a NUMBER, b NUMBER); TYPE aa IS TABLE OF BOOLEAN INDEX BY VARCHAR2(30); " +
	"TYPE ai IS TABLE OF NUMBER INDEX BY PLS_INTEGER; TYPE nt IS TABLE OF NUMBER; TYPE rt IS RECORD (x NUMBER); " +
	"TYPE rtt IS TABLE OF rt; s VARCHAR2(100) := 'SELECT 1 x FROM dual'; "

// TestQualifiedExpressions checks the choices of PL/SQL qualified
// expressions and Oracle 21c's FOR iterators against Oracle 23ai: every
// accepted form compiles there, and every rejected one fails with the
// error noted beside it. Rows ref_275 to ref_280 run them against the
// engine.
func TestQualifiedExpressions(t *testing.T) {
	for _, decl := range []string{
		"x r := r(a => 1, b => 2);",
		"x r := r(a => 1, OTHERS => 0);",
		"x aa := aa('a' => TRUE, 'b'=>FALSE);",
		"x aa := aa(UPPER('a') => TRUE, ('b') => FALSE);",
		"x ai := ai(1 => 10, 1+1 => 20, -1 => 30);",
		"x ai := ai(1|2 => 5, 4..6 => 7, 8 .. 9 | 11 => 0);",
		"x ai := ai(1, 2 => 20);",
		"x ai := ai(FOR i IN 1..3 => i*2);",
		"x ai := ai(FOR i IN 1..3 INDEX i+10 => i);",
		"x nt := nt(FOR i IN 1..3 SEQUENCE => i, 9);",
		"x ai := ai(FOR i IN 1..3 => i, FOR j IN 5..6 => j, 10 => 0);",
		"x ai := ai(FOR i IN 1..3, 5..6 => i);",
		"x ai := ai(FOR i IN REVERSE 1..3 BY 2 WHEN i > 1 => i);",
		"x ai := ai(1 => 1, FOR i IN 1..3 => i);",
		"x ai := ai();",
		"x nt := nt(9, FOR i IN 1..3 SEQUENCE => i, 10);",
		"x ai := ai(FOR i, j IN PAIRS OF a INDEX i => j);",
		"x nt := nt(FOR i, j IN PAIRS OF a SEQUENCE => j);",
		"x rtt := rtt(FOR r rt IN (EXECUTE IMMEDIATE s) SEQUENCE => r);",
		`TYPE "TRIM" IS TABLE OF NUMBER INDEX BY PLS_INTEGER; y "TRIM" := "TRIM"(1 => 2);`,
	} {
		ParseAndCheck(t, "DECLARE "+qualifiedDecls+decl+" BEGIN NULL; END;")
	}
	for _, sql := range []string{
		// Oracle 21c iteration controls in FOR loops.
		"BEGIN FOR i IN 1..10 BY 2 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN 1..3, 7..9 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN 1..3 WHILE i < 3 WHEN i > 1 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN REVERSE 1..3, REVERSE 5..6 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN 1, 5, 7 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN 1, REPEAT i*2 WHILE i < 100 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN VALUES OF v LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN REVERSE INDICES OF a WHEN i > 0 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i PLS_INTEGER, j NUMBER IN PAIRS OF a LOOP NULL; END LOOP; END;",
		"BEGIN FOR i MUTABLE PLS_INTEGER IN 1..3 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i NUMBER(5,2) IN 1..3 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN VALUES OF (SELECT 1 FROM dual) LOOP NULL; END LOOP; END;",
		"BEGIN FOR r IN c WHEN r.a > 0 LOOP NULL; END LOOP; END;",
		"BEGIN FOR r IN c2(1) LOOP NULL; END LOOP; END;",
		"BEGIN FOR r IN (SELECT 1 a FROM dual) LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN x + 1 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i, j MUTABLE IN PAIRS OF a LOOP NULL; END LOOP; END;",
		"BEGIN FOR i MUTABLE IN VALUES OF v LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN REVERSE VALUES OF v LOOP NULL; END LOOP; END;",
		// Dynamic SQL in parentheses, with a typed iterand.
		"BEGIN FOR r rt IN (EXECUTE IMMEDIATE s) LOOP NULL; END LOOP; END;",
		"BEGIN FOR r rt IN (EXECUTE IMMEDIATE 'SELECT 1 x FROM dual' USING 1) LOOP NULL; END LOOP; END;",
		"BEGIN FOR r rt IN VALUES OF (EXECUTE IMMEDIATE s USING IN 1, 2) LOOP NULL; END LOOP; END;",
		// The first control types a later dynamic one; PAIRS OF types its
		// value iterand.
		"BEGIN FOR r IN (SELECT 1 x FROM dual), (EXECUTE IMMEDIATE s) LOOP NULL; END LOOP; END;",
		"BEGIN FOR r IN 1..2, (EXECUTE IMMEDIATE s) LOOP NULL; END LOOP; END;",
		"BEGIN FOR i, r rt IN PAIRS OF (EXECUTE IMMEDIATE s) LOOP NULL; END LOOP; END;",
		// A qualified expression in a statement, a RETURN, a call, and a
		// parameter's default.
		"CREATE PROCEDURE p (x pkg.ai := pkg.ai(1 => 2)) IS BEGIN NULL; END;",
		"CREATE PACKAGE k AS PROCEDURE p (x pkg.ai := pkg.ai(1 => 2, 3 => 4)); END;",
		"BEGIN v := nt(FOR i IN VALUES OF w SEQUENCE => i); END;",
		"CREATE FUNCTION f RETURN ai IS BEGIN RETURN ai(1 => 1, 2 => 2); END;",
		"BEGIN p(aa('k' => TRUE)); END;",
	} {
		ParseAndCheck(t, sql)
	}

	for _, sql := range []string{
		// A positional value, a SEQUENCE iterator choice among them, after
		// an explicit choice (PLS-00312); a choice after OTHERS (PLS-00851);
		// PAIRS OF in a basic iterator choice (PLS-00864).
		"DECLARE " + qualifiedDecls + "x ai := ai(1 => 10, 20); BEGIN NULL; END;",
		"DECLARE " + qualifiedDecls + "x ai := ai(FOR i IN 1..3 => i, 9); BEGIN NULL; END;",
		"DECLARE " + qualifiedDecls + "x ai := ai(FOR i IN 1..3 INDEX i => i, 9); BEGIN NULL; END;",
		"DECLARE " + qualifiedDecls + "x nt := nt(1 => 1, FOR i IN 1..3 SEQUENCE => i); BEGIN NULL; END;",
		"DECLARE " + qualifiedDecls + "x r := r(OTHERS => 0, a => 1); BEGIN NULL; END;",
		"DECLARE " + qualifiedDecls + "x ai := ai(FOR i, j IN PAIRS OF a => j); BEGIN NULL; END;",
		// OTHERS stands alone (PLS-00851, PLS-00750).
		"DECLARE " + qualifiedDecls + "x r := r(a | OTHERS => 0); BEGIN NULL; END;",
		"DECLARE " + qualifiedDecls + "x ai := ai(1 | OTHERS => 0); BEGIN NULL; END;",
		// MUTABLE takes no INDICES OF iterand nor PAIRS OF index iterand
		// (PLS-00861, PLS-00854); REVERSE no collection a query returns
		// (PLS-00682); EXECUTE IMMEDIATE is no control (PLS-00103).
		"BEGIN FOR i MUTABLE IN INDICES OF a LOOP NULL; END LOOP; END;",
		"BEGIN FOR i MUTABLE, j IN PAIRS OF a LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN REVERSE VALUES OF (SELECT 1 FROM dual) LOOP NULL; END LOOP; END;",
		"BEGIN FOR i, j IN REVERSE PAIRS OF (SELECT 1, 2 FROM dual) LOOP NULL; END LOOP; END;",
		"BEGIN FOR q rt IN EXECUTE IMMEDIATE s LOOP NULL; END LOOP; END;",
		// Dynamic SQL needs a typed iterand (PLS-00858) and takes no REVERSE
		// (PLS-00682).
		"BEGIN FOR i IN VALUES OF (EXECUTE IMMEDIATE s) LOOP NULL; END LOOP; END;",
		"BEGIN FOR i PLS_INTEGER, r IN PAIRS OF (EXECUTE IMMEDIATE s) LOOP NULL; END LOOP; END;",
		// A query's binds are IN alone (PLS-00254).
		"BEGIN FOR r rt IN (EXECUTE IMMEDIATE s USING OUT v) LOOP NULL; END LOOP; END;",
		"BEGIN FOR r rt IN (EXECUTE IMMEDIATE s USING IN OUT v) LOOP NULL; END LOOP; END;",
		"BEGIN FOR r rt IN REVERSE (EXECUTE IMMEDIATE s) LOOP NULL; END LOOP; END;",
		// INDEX takes an expression, and not with SEQUENCE (PLS-00103).
		"DECLARE " + qualifiedDecls + "x ai := ai(FOR i IN 1..3 SEQUENCE INDEX i => i); BEGIN NULL; END;",
		"DECLARE " + qualifiedDecls + "x ai := ai(FOR i IN 1..3 INDEX => i); BEGIN NULL; END;",
		"DECLARE " + qualifiedDecls + "x ai := ai(1 | => 5); BEGIN NULL; END;",
		"DECLARE " + qualifiedDecls + "x ai := ai(1.. => 5); BEGIN NULL; END;",
		// SQL takes no choices (ORA-01748, ORA-00907).
		"SELECT t(1 => 1) FROM dual",
		"SELECT t('A' => 1) FROM dual",
		"SELECT UPPER(1 => 'x') FROM dual",
		"BEGIN SELECT t(1 => 1) INTO x FROM dual; END;",
		// Two iterands take PAIRS OF, which takes two (PLS-00856,
		// PLS-00855); no third (PLS-00103).
		"BEGIN FOR i, j IN 1..3 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i, j IN VALUES OF v LOOP NULL; END LOOP; END;",
		"BEGIN FOR i, j IN PAIRS OF a, 1..3 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN PAIRS OF a LOOP NULL; END LOOP; END;",
		"BEGIN FOR i, j, k IN PAIRS OF a LOOP NULL; END LOOP; END;",
		// REVERSE takes no single expression, cursor, or REPEAT (PLS-00857,
		// PLS-00682).
		"BEGIN FOR i IN REVERSE x LOOP NULL; END LOOP; END;",
		"BEGIN FOR r IN REVERSE c LOOP NULL; END LOOP; END;",
		"BEGIN FOR r IN REVERSE (SELECT 1 a FROM dual) LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN REVERSE REPEAT i+1 LOOP EXIT; END LOOP; END;",
		// BY after a range, once; WHILE before WHEN, once each; a control
		// after each comma; OF after VALUES (PLS-00103).
		"BEGIN FOR i IN 1 BY 2 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN 1..3 BY 2 BY 3 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN VALUES OF v BY 2 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN 1..3 WHEN i > 1 WHILE i < 3 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN 1..3 WHEN i > 1 WHEN i < 3 LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN 1..3, LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN VALUES v LOOP NULL; END LOOP; END;",
		"BEGIN FOR i IN repeat LOOP NULL; END LOOP; END;",
		"BEGIN FOR i 1..3 LOOP NULL; END LOOP; END;",
	} {
		ParseShouldFail(t, sql)
	}
}

// TestQualifiedExpressionTree checks the nodes built for choices and
// iterators, and that a classic FOR loop keeps its own fields.
func TestQualifiedExpressionTree(t *testing.T) {
	sql := "BEGIN x := ai('a' | 'b' => 1, 1+1..5 => 2, FOR i IN 1..3 INDEX i => i, OTHERS => 0); END;"
	block := parseOne[*ast.PLSQLBlock](t, sql)
	var call *ast.FuncCallExpr
	ast.Inspect(block, func(n ast.Node) bool {
		if fc, ok := n.(*ast.FuncCallExpr); ok && call == nil {
			call = fc
		}
		return true
	})
	if call == nil || call.Args.Len() != 4 {
		t.Fatalf("call = %+v", call)
	}
	list := call.Args.Items[0].(*ast.ChoiceAssocExpr)
	if list.Choices.Len() != 2 || sql[list.Loc.Start:list.Loc.End] != "'a' | 'b' => 1" {
		t.Errorf("choice list = %+v, spans %q", list, sql[list.Loc.Start:list.Loc.End])
	}
	// A choice starts at its first token, before a binary key's operator.
	rangeChoice := call.Args.Items[1].(*ast.ChoiceAssocExpr)
	rng := rangeChoice.Choices.Items[0].(*ast.ChoiceRange)
	if sql[rng.Loc.Start:rng.Loc.End] != "1+1..5" || sql[rangeChoice.Loc.Start:rangeChoice.Loc.End] != "1+1..5 => 2" {
		t.Errorf("range spans %q in %q", sql[rng.Loc.Start:rng.Loc.End], sql[rangeChoice.Loc.Start:rangeChoice.Loc.End])
	}
	if others := call.Args.Items[3].(*ast.ChoiceAssocExpr); !others.Others || others.Choices != nil {
		t.Errorf("OTHERS choice = %+v", others)
	}
	iter := call.Args.Items[2].(*ast.ChoiceAssocExpr)
	if iter.Iterator == nil || iter.Index == nil || iter.Sequence || sql[iter.Loc.Start:iter.Loc.End] != "FOR i IN 1..3 INDEX i => i" {
		t.Errorf("iterator choice = %+v", iter)
	}

	loop := parseOne[*ast.PLSQLBlock](t, "BEGIN FOR i IN REVERSE 1..3 LOOP NULL; END LOOP; END;").Statements.Items[0].(*ast.PLSQLLoop)
	if loop.Type != ast.LOOP_FOR || loop.Iterator != "I" || !loop.Reverse || loop.LowerBound == nil || loop.UpperBound == nil || loop.Iteration != nil {
		t.Errorf("range loop = %+v", loop)
	}
	// Each node of a loop is walked once, whichever fields hold it.
	calls := 0
	ast.Inspect(parseOne[*ast.PLSQLBlock](t, "BEGIN FOR i IN f()..g() LOOP NULL; END LOOP; END;"), func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncCallExpr); ok {
			calls++
		}
		return true
	})
	if calls != 2 {
		t.Errorf("walk meets %d calls in FOR i IN f()..g(), want 2", calls)
	}
	loop = parseOne[*ast.PLSQLBlock](t, "BEGIN FOR r IN c2(1) LOOP NULL; END LOOP; END;").Statements.Items[0].(*ast.PLSQLLoop)
	if loop.Type != ast.LOOP_CURSOR_FOR || loop.CursorName != "C2" || loop.CursorArgs.Len() != 1 {
		t.Errorf("cursor loop = %+v", loop)
	}
	loop = parseOne[*ast.PLSQLBlock](t, "BEGIN FOR i IN x + 1 LOOP NULL; END LOOP; END;").Statements.Items[0].(*ast.PLSQLLoop)
	if loop.Type != ast.LOOP_FOR || loop.CursorName != "" {
		t.Errorf("single-expression loop = %+v", loop)
	}
	loop = parseOne[*ast.PLSQLBlock](t, "BEGIN FOR i, v IN PAIRS OF t LOOP NULL; END LOOP; END;").Statements.Items[0].(*ast.PLSQLLoop)
	ctl := loop.Iteration.Controls.Items[0].(*ast.PLSQLIterationControl)
	if loop.Iteration.Iterands.Len() != 2 || ctl.Kind != "PAIRS OF" || ctl.Expr == nil || loop.LowerBound != nil {
		t.Errorf("pairs loop = %+v, control %+v", loop, ctl)
	}
}
