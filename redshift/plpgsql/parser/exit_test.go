package parser

import (
	"testing"

	"github.com/bytebase/omni/redshift/plpgsql/ast"
)

func TestExitBare(t *testing.T) {
	block := parseOK(t, `BEGIN EXIT; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	exit, ok := block.Body[0].(*ast.PLExit)
	if !ok {
		t.Fatalf("expected *ast.PLExit, got %T", block.Body[0])
	}
	if exit.IsContinue {
		t.Error("expected IsContinue=false")
	}
	if exit.Label != "" {
		t.Errorf("expected empty label, got %q", exit.Label)
	}
	if exit.Condition != "" {
		t.Errorf("expected empty condition, got %q", exit.Condition)
	}
}

func TestExitWithLabel(t *testing.T) {
	block := parseOK(t, `BEGIN EXIT myloop; END`)
	exit := block.Body[0].(*ast.PLExit)
	if exit.IsContinue {
		t.Error("expected IsContinue=false")
	}
	if exit.Label != "myloop" {
		t.Errorf("expected label %q, got %q", "myloop", exit.Label)
	}
	if exit.Condition != "" {
		t.Errorf("expected empty condition, got %q", exit.Condition)
	}
}

func TestExitWithWhen(t *testing.T) {
	block := parseOK(t, `BEGIN EXIT WHEN x > 10; END`)
	exit := block.Body[0].(*ast.PLExit)
	if exit.IsContinue {
		t.Error("expected IsContinue=false")
	}
	if exit.Label != "" {
		t.Errorf("expected empty label, got %q", exit.Label)
	}
	if exit.Condition != "x > 10" {
		t.Errorf("expected condition %q, got %q", "x > 10", exit.Condition)
	}
}

func TestExitWithLabelAndWhen(t *testing.T) {
	block := parseOK(t, `BEGIN EXIT myloop WHEN x > 10; END`)
	exit := block.Body[0].(*ast.PLExit)
	if exit.IsContinue {
		t.Error("expected IsContinue=false")
	}
	if exit.Label != "myloop" {
		t.Errorf("expected label %q, got %q", "myloop", exit.Label)
	}
	if exit.Condition != "x > 10" {
		t.Errorf("expected condition %q, got %q", "x > 10", exit.Condition)
	}
}

func TestContinueBare(t *testing.T) {
	block := parseOK(t, `BEGIN CONTINUE; END`)
	exit := block.Body[0].(*ast.PLExit)
	if !exit.IsContinue {
		t.Error("expected IsContinue=true")
	}
	if exit.Label != "" {
		t.Errorf("expected empty label, got %q", exit.Label)
	}
	if exit.Condition != "" {
		t.Errorf("expected empty condition, got %q", exit.Condition)
	}
}

func TestContinueWithLabel(t *testing.T) {
	block := parseOK(t, `BEGIN CONTINUE myloop; END`)
	exit := block.Body[0].(*ast.PLExit)
	if !exit.IsContinue {
		t.Error("expected IsContinue=true")
	}
	if exit.Label != "myloop" {
		t.Errorf("expected label %q, got %q", "myloop", exit.Label)
	}
}

func TestContinueWithWhen(t *testing.T) {
	block := parseOK(t, `BEGIN CONTINUE WHEN x = 0; END`)
	exit := block.Body[0].(*ast.PLExit)
	if !exit.IsContinue {
		t.Error("expected IsContinue=true")
	}
	if exit.Label != "" {
		t.Errorf("expected empty label, got %q", exit.Label)
	}
	if exit.Condition != "x = 0" {
		t.Errorf("expected condition %q, got %q", "x = 0", exit.Condition)
	}
}

func TestContinueWithLabelAndWhen(t *testing.T) {
	block := parseOK(t, `BEGIN CONTINUE myloop WHEN x = 0; END`)
	exit := block.Body[0].(*ast.PLExit)
	if !exit.IsContinue {
		t.Error("expected IsContinue=true")
	}
	if exit.Label != "myloop" {
		t.Errorf("expected label %q, got %q", "myloop", exit.Label)
	}
	if exit.Condition != "x = 0" {
		t.Errorf("expected condition %q, got %q", "x = 0", exit.Condition)
	}
}

func TestExitWhenConditionSpansUntilSemicolon(t *testing.T) {
	block := parseOK(t, `BEGIN EXIT WHEN (x > 10 AND y < 20) OR z = 0; END`)
	exit := block.Body[0].(*ast.PLExit)
	if exit.Condition != "(x > 10 AND y < 20) OR z = 0" {
		t.Errorf("expected condition %q, got %q", "(x > 10 AND y < 20) OR z = 0", exit.Condition)
	}
}
