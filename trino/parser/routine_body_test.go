package parser

import (
	"testing"
)

// This file is the structural (Layer 2) gate for the parser-routines node's
// control-flow body language (routine_body.go): RETURN, SET (assignment),
// IF/ELSEIF/ELSE, CASE (simple + searched), LOOP/WHILE/REPEAT (labeled),
// ITERATE/LEAVE, BEGIN/END (compound) with DECLARE. It pins the parse-node shape
// of each controlStatement alternative. The authoritative accept/reject
// differential against the live Trino 481 oracle is in oracle_routines_test.go.

// bodyOf parses `CREATE FUNCTION f() RETURNS int <body>` and returns the
// function's single control-statement body, failing the test if the wrapper or
// parse is wrong. It lets each test focus on one routine statement shape.
func bodyOf(t *testing.T, body string) RoutineStatement {
	t.Helper()
	sql := "CREATE FUNCTION f() RETURNS int " + body
	s, ok := parseOneStmt(t, sql).(*CreateFunctionStmt)
	if !ok {
		t.Fatalf("Parse(%q): not a *CreateFunctionStmt", sql)
	}
	return s.Spec.Body
}

func TestRoutineBody_Return(t *testing.T) {
	s, ok := bodyOf(t, "RETURN 1 + 2").(*ReturnStatement)
	if !ok {
		t.Fatalf("body = %T, want *ReturnStatement", s)
	}
	if s.Value == nil {
		t.Errorf("Value = nil")
	}
}

func TestRoutineBody_Compound(t *testing.T) {
	t.Run("declares_and_statements", func(t *testing.T) {
		s := bodyOf(t, "BEGIN DECLARE a, b int; DECLARE c varchar DEFAULT 'x'; SET a = 1; RETURN a; END").(*CompoundStatement)
		if len(s.Declarations) != 2 {
			t.Fatalf("declarations = %d, want 2", len(s.Declarations))
		}
		if len(s.Declarations[0].Names) != 2 {
			t.Errorf("decl[0] names = %d, want 2", len(s.Declarations[0].Names))
		}
		if s.Declarations[1].Default == nil {
			t.Errorf("decl[1] DEFAULT = nil, want a value")
		}
		if len(s.Body) != 2 {
			t.Errorf("body statements = %d, want 2", len(s.Body))
		}
	})

	t.Run("empty", func(t *testing.T) {
		// R6: an empty BEGIN END is valid (both blocks optional).
		s := bodyOf(t, "BEGIN END").(*CompoundStatement)
		if len(s.Declarations) != 0 || len(s.Body) != 0 {
			t.Errorf("empty compound: decls=%d body=%d, want 0/0", len(s.Declarations), len(s.Body))
		}
	})

	t.Run("declares_only", func(t *testing.T) {
		s := bodyOf(t, "BEGIN DECLARE x int; END").(*CompoundStatement)
		if len(s.Declarations) != 1 {
			t.Errorf("declarations = %d, want 1", len(s.Declarations))
		}
		if len(s.Body) != 0 {
			t.Errorf("body statements = %d, want 0", len(s.Body))
		}
	})

	t.Run("nested", func(t *testing.T) {
		s := bodyOf(t, "BEGIN BEGIN RETURN 1; END; END").(*CompoundStatement)
		if len(s.Body) != 1 {
			t.Fatalf("outer body statements = %d, want 1", len(s.Body))
		}
		if _, ok := s.Body[0].(*CompoundStatement); !ok {
			t.Errorf("inner = %T, want *CompoundStatement", s.Body[0])
		}
	})
}

func TestRoutineBody_Assignment(t *testing.T) {
	s := bodyOf(t, "BEGIN DECLARE x int; SET x = 1 + 2; RETURN x; END").(*CompoundStatement)
	assign, ok := s.Body[0].(*AssignmentStatement)
	if !ok {
		t.Fatalf("body[0] = %T, want *AssignmentStatement", s.Body[0])
	}
	if assign.Target.Normalize() != "x" {
		t.Errorf("target = %q, want x", assign.Target.Normalize())
	}
	if assign.Value == nil {
		t.Errorf("value = nil")
	}
}

func TestRoutineBody_If(t *testing.T) {
	s := bodyOf(t, "BEGIN IF 1 > 0 THEN RETURN 1; ELSEIF 1 < 0 THEN RETURN -1; ELSE RETURN 0; END IF; END").(*CompoundStatement)
	ifs, ok := s.Body[0].(*IfStatement)
	if !ok {
		t.Fatalf("body[0] = %T, want *IfStatement", s.Body[0])
	}
	if ifs.Condition == nil {
		t.Errorf("condition = nil")
	}
	if len(ifs.Then) != 1 {
		t.Errorf("then = %d, want 1", len(ifs.Then))
	}
	if len(ifs.ElseIfs) != 1 {
		t.Errorf("elseifs = %d, want 1", len(ifs.ElseIfs))
	}
	if len(ifs.Else) != 1 {
		t.Errorf("else = %d, want 1", len(ifs.Else))
	}
}

func TestRoutineBody_If_NoElse(t *testing.T) {
	s := bodyOf(t, "BEGIN IF 1 > 0 THEN RETURN 1; END IF; RETURN 0; END").(*CompoundStatement)
	ifs := s.Body[0].(*IfStatement)
	if len(ifs.ElseIfs) != 0 {
		t.Errorf("elseifs = %d, want 0", len(ifs.ElseIfs))
	}
	if ifs.Else != nil {
		t.Errorf("else = %v, want nil", ifs.Else)
	}
}

func TestRoutineBody_CaseSimple(t *testing.T) {
	s := bodyOf(t, "BEGIN CASE 1 WHEN 1 THEN RETURN 1; WHEN 2 THEN RETURN 2; ELSE RETURN 0; END CASE; END").(*CompoundStatement)
	cs, ok := s.Body[0].(*CaseStatement)
	if !ok {
		t.Fatalf("body[0] = %T, want *CaseStatement", s.Body[0])
	}
	if cs.Operand == nil {
		t.Errorf("operand = nil, want a subject expression (simple CASE)")
	}
	if len(cs.Whens) != 2 {
		t.Errorf("whens = %d, want 2", len(cs.Whens))
	}
	if len(cs.Else) != 1 {
		t.Errorf("else = %d, want 1", len(cs.Else))
	}
}

func TestRoutineBody_CaseSearched(t *testing.T) {
	s := bodyOf(t, "BEGIN CASE WHEN 1 > 0 THEN RETURN 1; ELSE RETURN 0; END CASE; END").(*CompoundStatement)
	cs := s.Body[0].(*CaseStatement)
	if cs.Operand != nil {
		t.Errorf("operand = %v, want nil (searched CASE)", cs.Operand)
	}
	if len(cs.Whens) != 1 {
		t.Errorf("whens = %d, want 1", len(cs.Whens))
	}
}

func TestRoutineBody_Loop(t *testing.T) {
	t.Run("labeled_with_leave", func(t *testing.T) {
		s := bodyOf(t, "BEGIN DECLARE i int DEFAULT 0; top: LOOP SET i = i + 1; IF i > 10 THEN LEAVE top; END IF; END LOOP; RETURN i; END").(*CompoundStatement)
		loop, ok := s.Body[0].(*LoopStatement)
		if !ok {
			t.Fatalf("body[0] = %T, want *LoopStatement", s.Body[0])
		}
		if loop.Label == nil || loop.Label.Normalize() != "top" {
			t.Errorf("label = %v, want top", loop.Label)
		}
		// inner: SET, then IF (containing LEAVE)
		if len(loop.Body) != 2 {
			t.Fatalf("loop body = %d, want 2", len(loop.Body))
		}
		innerIf := loop.Body[1].(*IfStatement)
		if _, ok := innerIf.Then[0].(*LeaveStatement); !ok {
			t.Errorf("LEAVE = %T, want *LeaveStatement", innerIf.Then[0])
		}
	})

	t.Run("unlabeled", func(t *testing.T) {
		s := bodyOf(t, "BEGIN LOOP RETURN 1; END LOOP; END").(*CompoundStatement)
		loop := s.Body[0].(*LoopStatement)
		if loop.Label != nil {
			t.Errorf("label = %v, want nil", loop.Label)
		}
	})

	t.Run("label_named_as_control_keyword", func(t *testing.T) {
		// A label may be a non-reserved control keyword (`loop`, `set`, …): the
		// label check must win over the keyword dispatch.
		s := bodyOf(t, "BEGIN loop: LOOP LEAVE loop; END LOOP; RETURN 1; END").(*CompoundStatement)
		loop, ok := s.Body[0].(*LoopStatement)
		if !ok {
			t.Fatalf("body[0] = %T, want *LoopStatement", s.Body[0])
		}
		if loop.Label == nil || loop.Label.Normalize() != "loop" {
			t.Errorf("label = %v, want loop", loop.Label)
		}
		leave := loop.Body[0].(*LeaveStatement)
		if leave.Label.Normalize() != "loop" {
			t.Errorf("LEAVE label = %q, want loop", leave.Label.Normalize())
		}
	})
}

func TestRoutineBody_While(t *testing.T) {
	t.Run("labeled_with_iterate", func(t *testing.T) {
		s := bodyOf(t, "BEGIN DECLARE i int DEFAULT 0; abc: WHILE i < 10 DO SET i = i + 1; ITERATE abc; END WHILE; RETURN i; END").(*CompoundStatement)
		w, ok := s.Body[0].(*WhileStatement)
		if !ok {
			t.Fatalf("body[0] = %T, want *WhileStatement", s.Body[0])
		}
		if w.Label == nil || w.Label.Normalize() != "abc" {
			t.Errorf("label = %v, want abc", w.Label)
		}
		if w.Condition == nil {
			t.Errorf("condition = nil")
		}
		iter, ok := w.Body[1].(*IterateStatement)
		if !ok {
			t.Fatalf("body[1] = %T, want *IterateStatement", w.Body[1])
		}
		if iter.Label.Normalize() != "abc" {
			t.Errorf("ITERATE label = %q, want abc", iter.Label.Normalize())
		}
	})
}

func TestRoutineBody_Repeat(t *testing.T) {
	s := bodyOf(t, "BEGIN DECLARE i int DEFAULT 0; REPEAT SET i = i + 1; UNTIL i >= 10 END REPEAT; RETURN i; END").(*CompoundStatement)
	r, ok := s.Body[0].(*RepeatStatement)
	if !ok {
		t.Fatalf("body[0] = %T, want *RepeatStatement", s.Body[0])
	}
	if len(r.Body) != 1 {
		t.Errorf("repeat body = %d, want 1", len(r.Body))
	}
	if r.Condition == nil {
		t.Errorf("UNTIL condition = nil")
	}
}
