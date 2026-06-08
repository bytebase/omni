package parser

import (
	"testing"

	"github.com/bytebase/omni/redshift/plpgsql/ast"
)

// --------------------------------------------------------------------------
// Section 5.2: ASSERT, CALL, Transaction Control, NULL
// --------------------------------------------------------------------------

func TestAssert(t *testing.T) {
	t.Run("simple assert", func(t *testing.T) {
		block := parseOK(t, `BEGIN ASSERT x > 0; END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		node, ok := block.Body[0].(*ast.PLAssert)
		if !ok {
			t.Fatalf("expected PLAssert, got %T", block.Body[0])
		}
		if node.Condition != "x > 0" {
			t.Errorf("Condition = %q, want %q", node.Condition, "x > 0")
		}
		if node.Message != "" {
			t.Errorf("Message = %q, want empty", node.Message)
		}
	})

	t.Run("assert with message", func(t *testing.T) {
		block := parseOK(t, `BEGIN ASSERT x > 0, 'x must be positive'; END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		node, ok := block.Body[0].(*ast.PLAssert)
		if !ok {
			t.Fatalf("expected PLAssert, got %T", block.Body[0])
		}
		if node.Condition != "x > 0" {
			t.Errorf("Condition = %q, want %q", node.Condition, "x > 0")
		}
		if node.Message != "'x must be positive'" {
			t.Errorf("Message = %q, want %q", node.Message, "'x must be positive'")
		}
	})
}

func TestCallStmt(t *testing.T) {
	t.Run("call procedure", func(t *testing.T) {
		block := parseOK(t, `BEGIN CALL my_proc(1, 2); END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		node, ok := block.Body[0].(*ast.PLCall)
		if !ok {
			t.Fatalf("expected PLCall, got %T", block.Body[0])
		}
		if node.SQLText != "CALL my_proc(1, 2)" {
			t.Errorf("SQLText = %q, want %q", node.SQLText, "CALL my_proc(1, 2)")
		}
	})

	t.Run("call with schema", func(t *testing.T) {
		block := parseOK(t, `BEGIN CALL public.my_proc(); END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		node, ok := block.Body[0].(*ast.PLCall)
		if !ok {
			t.Fatalf("expected PLCall, got %T", block.Body[0])
		}
		if node.SQLText != "CALL public.my_proc()" {
			t.Errorf("SQLText = %q, want %q", node.SQLText, "CALL public.my_proc()")
		}
	})

	t.Run("do block", func(t *testing.T) {
		block := parseOK(t, `BEGIN DO $$ BEGIN RAISE NOTICE 'hi'; END $$; END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		node, ok := block.Body[0].(*ast.PLCall)
		if !ok {
			t.Fatalf("expected PLCall, got %T", block.Body[0])
		}
		if node.SQLText != "DO $$ BEGIN RAISE NOTICE 'hi'; END $$" {
			t.Errorf("SQLText = %q, want %q", node.SQLText, "DO $$ BEGIN RAISE NOTICE 'hi'; END $$")
		}
	})
}

func TestCommit(t *testing.T) {
	t.Run("bare commit", func(t *testing.T) {
		block := parseOK(t, `BEGIN COMMIT; END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		node, ok := block.Body[0].(*ast.PLCommit)
		if !ok {
			t.Fatalf("expected PLCommit, got %T", block.Body[0])
		}
		if node.Chain {
			t.Error("Chain = true, want false")
		}
	})

	t.Run("commit and chain", func(t *testing.T) {
		block := parseOK(t, `BEGIN COMMIT AND CHAIN; END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		node, ok := block.Body[0].(*ast.PLCommit)
		if !ok {
			t.Fatalf("expected PLCommit, got %T", block.Body[0])
		}
		if !node.Chain {
			t.Error("Chain = false, want true")
		}
	})

	t.Run("commit and no chain", func(t *testing.T) {
		block := parseOK(t, `BEGIN COMMIT AND NO CHAIN; END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		node, ok := block.Body[0].(*ast.PLCommit)
		if !ok {
			t.Fatalf("expected PLCommit, got %T", block.Body[0])
		}
		if node.Chain {
			t.Error("Chain = true, want false")
		}
	})
}

func TestRollback(t *testing.T) {
	t.Run("bare rollback", func(t *testing.T) {
		block := parseOK(t, `BEGIN ROLLBACK; END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		node, ok := block.Body[0].(*ast.PLRollback)
		if !ok {
			t.Fatalf("expected PLRollback, got %T", block.Body[0])
		}
		if node.Chain {
			t.Error("Chain = true, want false")
		}
	})

	t.Run("rollback and chain", func(t *testing.T) {
		block := parseOK(t, `BEGIN ROLLBACK AND CHAIN; END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		node, ok := block.Body[0].(*ast.PLRollback)
		if !ok {
			t.Fatalf("expected PLRollback, got %T", block.Body[0])
		}
		if !node.Chain {
			t.Error("Chain = false, want true")
		}
	})

	t.Run("rollback and no chain", func(t *testing.T) {
		block := parseOK(t, `BEGIN ROLLBACK AND NO CHAIN; END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		node, ok := block.Body[0].(*ast.PLRollback)
		if !ok {
			t.Fatalf("expected PLRollback, got %T", block.Body[0])
		}
		if node.Chain {
			t.Error("Chain = true, want false")
		}
	})
}

func TestNullStmt(t *testing.T) {
	t.Run("bare null", func(t *testing.T) {
		block := parseOK(t, `BEGIN NULL; END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		_, ok := block.Body[0].(*ast.PLNull)
		if !ok {
			t.Fatalf("expected PLNull, got %T", block.Body[0])
		}
	})

	t.Run("null in empty branch", func(t *testing.T) {
		block := parseOK(t, `BEGIN IF x THEN NULL; ELSE do_something := 1; END IF; END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		ifNode, ok := block.Body[0].(*ast.PLIf)
		if !ok {
			t.Fatalf("expected PLIf, got %T", block.Body[0])
		}
		if len(ifNode.ThenBody) != 1 {
			t.Fatalf("expected 1 THEN statement, got %d", len(ifNode.ThenBody))
		}
		_, ok = ifNode.ThenBody[0].(*ast.PLNull)
		if !ok {
			t.Fatalf("expected PLNull in THEN body, got %T", ifNode.ThenBody[0])
		}
		if len(ifNode.ElseBody) != 1 {
			t.Fatalf("expected 1 ELSE statement, got %d", len(ifNode.ElseBody))
		}
	})
}
