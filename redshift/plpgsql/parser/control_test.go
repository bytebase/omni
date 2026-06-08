package parser

import (
	"testing"

	"github.com/bytebase/omni/redshift/plpgsql/ast"
)

// --- Section 2.1: IF / ELSIF / ELSE ---

func TestIfSimple(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			"simple IF-THEN-END IF",
			"BEGIN IF x > 0 THEN y := 1; END IF; END",
		},
		{
			"IF-THEN-ELSE",
			"BEGIN IF x > 0 THEN y := 1; ELSE y := 0; END IF; END",
		},
		{
			"IF-THEN-ELSIF-ELSE",
			"BEGIN IF x > 0 THEN y := 1; ELSIF x = 0 THEN y := 0; ELSE y := -1; END IF; END",
		},
		{
			"multiple ELSIF branches",
			"BEGIN IF a THEN x := 1; ELSIF b THEN x := 2; ELSIF c THEN x := 3; ELSIF d THEN x := 4; END IF; END",
		},
		{
			"ELSEIF synonym",
			"BEGIN IF x THEN y := 1; ELSEIF z THEN y := 2; END IF; END",
		},
		{
			"nested IF inside IF",
			"BEGIN IF a THEN IF b THEN c := 1; END IF; END IF; END",
		},
		{
			"condition expression spans until THEN",
			"BEGIN IF x > 0 AND y < 10 THEN z := 1; END IF; END",
		},
		{
			"empty THEN body",
			"BEGIN IF x THEN END IF; END",
		},
		{
			"multiple statements in THEN body",
			"BEGIN IF x THEN a := 1; b := 2; c := 3; END IF; END",
		},
		{
			"multiple statements in ELSE body",
			"BEGIN IF x THEN a := 1; ELSE b := 2; c := 3; END IF; END",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := parseOK(t, tt.body)
			if len(block.Body) == 0 {
				t.Fatal("expected at least one statement in block body")
			}
			ifNode, ok := block.Body[0].(*ast.PLIf)
			if !ok {
				t.Fatalf("expected PLIf, got %T", block.Body[0])
			}
			if ifNode.Condition == "" {
				t.Fatal("expected non-empty condition")
			}
		})
	}
}

func TestIfStructure(t *testing.T) {
	t.Run("IF-THEN-ELSIF-ELSE structure", func(t *testing.T) {
		block := parseOK(t, "BEGIN IF x > 0 THEN y := 1; ELSIF x = 0 THEN y := 0; ELSE y := -1; END IF; END")
		ifNode := block.Body[0].(*ast.PLIf)

		if ifNode.Condition != "x > 0" {
			t.Errorf("condition = %q, want %q", ifNode.Condition, "x > 0")
		}
		if len(ifNode.ThenBody) != 1 {
			t.Errorf("ThenBody len = %d, want 1", len(ifNode.ThenBody))
		}
		if len(ifNode.ElsIfs) != 1 {
			t.Errorf("ElsIfs len = %d, want 1", len(ifNode.ElsIfs))
		}
		if len(ifNode.ElseBody) != 1 {
			t.Errorf("ElseBody len = %d, want 1", len(ifNode.ElseBody))
		}

		elsif := ifNode.ElsIfs[0].(*ast.PLElsIf)
		if elsif.Condition != "x = 0" {
			t.Errorf("elsif condition = %q, want %q", elsif.Condition, "x = 0")
		}
	})

	t.Run("multiple ELSIF branches", func(t *testing.T) {
		block := parseOK(t, "BEGIN IF a THEN x := 1; ELSIF b THEN x := 2; ELSIF c THEN x := 3; ELSIF d THEN x := 4; END IF; END")
		ifNode := block.Body[0].(*ast.PLIf)
		if len(ifNode.ElsIfs) != 3 {
			t.Errorf("ElsIfs len = %d, want 3", len(ifNode.ElsIfs))
		}
	})

	t.Run("ELSEIF synonym", func(t *testing.T) {
		block := parseOK(t, "BEGIN IF x THEN y := 1; ELSEIF z THEN y := 2; END IF; END")
		ifNode := block.Body[0].(*ast.PLIf)
		if len(ifNode.ElsIfs) != 1 {
			t.Errorf("ElsIfs len = %d, want 1", len(ifNode.ElsIfs))
		}
	})

	t.Run("nested IF", func(t *testing.T) {
		block := parseOK(t, "BEGIN IF a THEN IF b THEN c := 1; END IF; END IF; END")
		ifNode := block.Body[0].(*ast.PLIf)
		if len(ifNode.ThenBody) != 1 {
			t.Fatalf("ThenBody len = %d, want 1", len(ifNode.ThenBody))
		}
		innerIf, ok := ifNode.ThenBody[0].(*ast.PLIf)
		if !ok {
			t.Fatalf("expected nested PLIf, got %T", ifNode.ThenBody[0])
		}
		if innerIf.Condition != "b" {
			t.Errorf("inner condition = %q, want %q", innerIf.Condition, "b")
		}
	})

	t.Run("condition spans until THEN", func(t *testing.T) {
		block := parseOK(t, "BEGIN IF x > 0 AND y < 10 THEN z := 1; END IF; END")
		ifNode := block.Body[0].(*ast.PLIf)
		if ifNode.Condition != "x > 0 AND y < 10" {
			t.Errorf("condition = %q, want %q", ifNode.Condition, "x > 0 AND y < 10")
		}
	})

	t.Run("empty THEN body", func(t *testing.T) {
		block := parseOK(t, "BEGIN IF x THEN END IF; END")
		ifNode := block.Body[0].(*ast.PLIf)
		if len(ifNode.ThenBody) != 0 {
			t.Errorf("ThenBody len = %d, want 0", len(ifNode.ThenBody))
		}
	})
}

func TestIfErrors(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		wantContains string
	}{
		{
			"missing THEN",
			"BEGIN IF x END IF; END",
			"THEN",
		},
		{
			"missing END IF",
			"BEGIN IF x THEN y := 1; END",
			"IF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseErr(t, tt.body, tt.wantContains)
		})
	}
}

// --- Section 2.2: CASE Statement ---

func TestCaseSearched(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			"searched CASE",
			"BEGIN CASE WHEN x > 0 THEN y := 1; WHEN x = 0 THEN y := 0; END CASE; END",
		},
		{
			"searched CASE with ELSE",
			"BEGIN CASE WHEN x > 0 THEN y := 1; ELSE y := -1; END CASE; END",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := parseOK(t, tt.body)
			caseNode, ok := block.Body[0].(*ast.PLCase)
			if !ok {
				t.Fatalf("expected PLCase, got %T", block.Body[0])
			}
			if caseNode.HasTest {
				t.Error("expected searched CASE (HasTest=false)")
			}
		})
	}
}

func TestCaseSimple(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			"simple CASE",
			"BEGIN CASE x WHEN 1 THEN y := 'a'; WHEN 2 THEN y := 'b'; END CASE; END",
		},
		{
			"simple CASE with ELSE",
			"BEGIN CASE x WHEN 1 THEN y := 'a'; ELSE y := 'z'; END CASE; END",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := parseOK(t, tt.body)
			caseNode, ok := block.Body[0].(*ast.PLCase)
			if !ok {
				t.Fatalf("expected PLCase, got %T", block.Body[0])
			}
			if !caseNode.HasTest {
				t.Error("expected simple CASE (HasTest=true)")
			}
			if caseNode.TestExpr != "x" {
				t.Errorf("TestExpr = %q, want %q", caseNode.TestExpr, "x")
			}
		})
	}
}

func TestCaseStructure(t *testing.T) {
	t.Run("multiple WHEN branches", func(t *testing.T) {
		block := parseOK(t, "BEGIN CASE WHEN a THEN x := 1; WHEN b THEN x := 2; WHEN c THEN x := 3; END CASE; END")
		caseNode := block.Body[0].(*ast.PLCase)
		if len(caseNode.Whens) != 3 {
			t.Errorf("Whens len = %d, want 3", len(caseNode.Whens))
		}
	})

	t.Run("multiple statements per WHEN body", func(t *testing.T) {
		block := parseOK(t, "BEGIN CASE WHEN x > 0 THEN a := 1; b := 2; WHEN x = 0 THEN c := 3; END CASE; END")
		caseNode := block.Body[0].(*ast.PLCase)
		when0 := caseNode.Whens[0].(*ast.PLCaseWhen)
		if len(when0.Body) != 2 {
			t.Errorf("first WHEN body len = %d, want 2", len(when0.Body))
		}
	})

	t.Run("nested CASE inside CASE", func(t *testing.T) {
		block := parseOK(t, "BEGIN CASE WHEN a THEN CASE WHEN b THEN x := 1; END CASE; END CASE; END")
		caseNode := block.Body[0].(*ast.PLCase)
		when0 := caseNode.Whens[0].(*ast.PLCaseWhen)
		if len(when0.Body) != 1 {
			t.Fatalf("first WHEN body len = %d, want 1", len(when0.Body))
		}
		innerCase, ok := when0.Body[0].(*ast.PLCase)
		if !ok {
			t.Fatalf("expected nested PLCase, got %T", when0.Body[0])
		}
		if len(innerCase.Whens) != 1 {
			t.Errorf("inner CASE whens = %d, want 1", len(innerCase.Whens))
		}
	})

	t.Run("WHEN expression text", func(t *testing.T) {
		block := parseOK(t, "BEGIN CASE WHEN x > 0 THEN y := 1; END CASE; END")
		caseNode := block.Body[0].(*ast.PLCase)
		when0 := caseNode.Whens[0].(*ast.PLCaseWhen)
		if when0.Expr != "x > 0" {
			t.Errorf("WHEN expr = %q, want %q", when0.Expr, "x > 0")
		}
	})
}

func TestCaseErrors(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		wantContains string
	}{
		{
			"missing END CASE",
			"BEGIN CASE WHEN x THEN y := 1; END",
			"CASE",
		},
		{
			"no WHEN clause",
			"BEGIN CASE END CASE; END",
			"WHEN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseErr(t, tt.body, tt.wantContains)
		})
	}
}
