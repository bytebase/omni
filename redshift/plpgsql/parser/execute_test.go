package parser

import (
	"testing"

	"github.com/bytebase/omni/redshift/plpgsql/ast"
)

func TestExecuteBasic(t *testing.T) {
	// Basic EXECUTE: EXECUTE 'SELECT 1';
	block := parseOK(t, `BEGIN EXECUTE 'SELECT 1'; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	exec, ok := block.Body[0].(*ast.PLDynExecute)
	if !ok {
		t.Fatalf("expected PLDynExecute, got %T", block.Body[0])
	}
	if exec.Query != "'SELECT 1'" {
		t.Errorf("Query = %q, want %q", exec.Query, "'SELECT 1'")
	}
	if exec.Into != nil {
		t.Errorf("Into = %v, want nil", exec.Into)
	}
	if exec.Strict {
		t.Error("Strict = true, want false")
	}
	if exec.Params != nil {
		t.Errorf("Params = %v, want nil", exec.Params)
	}
}

func TestExecuteInto(t *testing.T) {
	// EXECUTE with INTO
	block := parseOK(t, `BEGIN EXECUTE 'SELECT a FROM t' INTO x; END`)
	exec := block.Body[0].(*ast.PLDynExecute)
	if exec.Query != "'SELECT a FROM t'" {
		t.Errorf("Query = %q, want %q", exec.Query, "'SELECT a FROM t'")
	}
	if len(exec.Into) != 1 || exec.Into[0] != "x" {
		t.Errorf("Into = %v, want [x]", exec.Into)
	}
	if exec.Strict {
		t.Error("Strict = true, want false")
	}
}

func TestExecuteIntoStrict(t *testing.T) {
	// EXECUTE with INTO STRICT
	block := parseOK(t, `BEGIN EXECUTE 'SELECT a FROM t' INTO STRICT x; END`)
	exec := block.Body[0].(*ast.PLDynExecute)
	if exec.Query != "'SELECT a FROM t'" {
		t.Errorf("Query = %q, want %q", exec.Query, "'SELECT a FROM t'")
	}
	if len(exec.Into) != 1 || exec.Into[0] != "x" {
		t.Errorf("Into = %v, want [x]", exec.Into)
	}
	if !exec.Strict {
		t.Error("Strict = false, want true")
	}
}

func TestExecuteUsing(t *testing.T) {
	// EXECUTE with USING
	block := parseOK(t, `BEGIN EXECUTE 'INSERT INTO t VALUES($1)' USING val; END`)
	exec := block.Body[0].(*ast.PLDynExecute)
	if exec.Query != "'INSERT INTO t VALUES($1)'" {
		t.Errorf("Query = %q, want %q", exec.Query, "'INSERT INTO t VALUES($1)'")
	}
	if len(exec.Params) != 1 || exec.Params[0] != "val" {
		t.Errorf("Params = %v, want [val]", exec.Params)
	}
}

func TestExecuteMultipleUsing(t *testing.T) {
	// EXECUTE with multiple USING params
	block := parseOK(t, `BEGIN EXECUTE $q$ SELECT $1, $2 $q$ USING a, b; END`)
	exec := block.Body[0].(*ast.PLDynExecute)
	if len(exec.Params) != 2 {
		t.Fatalf("Params length = %d, want 2", len(exec.Params))
	}
	if exec.Params[0] != "a" {
		t.Errorf("Params[0] = %q, want %q", exec.Params[0], "a")
	}
	if exec.Params[1] != "b" {
		t.Errorf("Params[1] = %q, want %q", exec.Params[1], "b")
	}
}

func TestExecuteIntoAndUsing(t *testing.T) {
	// EXECUTE with INTO and USING (INTO first)
	block := parseOK(t, `BEGIN EXECUTE 'SELECT a FROM t WHERE id=$1' INTO x USING my_id; END`)
	exec := block.Body[0].(*ast.PLDynExecute)
	if exec.Query != "'SELECT a FROM t WHERE id=$1'" {
		t.Errorf("Query = %q, want %q", exec.Query, "'SELECT a FROM t WHERE id=$1'")
	}
	if len(exec.Into) != 1 || exec.Into[0] != "x" {
		t.Errorf("Into = %v, want [x]", exec.Into)
	}
	if len(exec.Params) != 1 || exec.Params[0] != "my_id" {
		t.Errorf("Params = %v, want [my_id]", exec.Params)
	}
}

func TestExecuteUsingBeforeInto(t *testing.T) {
	// EXECUTE with USING before INTO (reversed order)
	block := parseOK(t, `BEGIN EXECUTE 'SELECT $1' USING a INTO x; END`)
	exec := block.Body[0].(*ast.PLDynExecute)
	if exec.Query != "'SELECT $1'" {
		t.Errorf("Query = %q, want %q", exec.Query, "'SELECT $1'")
	}
	if len(exec.Into) != 1 || exec.Into[0] != "x" {
		t.Errorf("Into = %v, want [x]", exec.Into)
	}
	if len(exec.Params) != 1 || exec.Params[0] != "a" {
		t.Errorf("Params = %v, want [a]", exec.Params)
	}
}

func TestExecuteExprConcat(t *testing.T) {
	// EXECUTE expression concatenation
	block := parseOK(t, `BEGIN EXECUTE 'SELECT * FROM ' || quote_ident(tbl); END`)
	exec := block.Body[0].(*ast.PLDynExecute)
	if exec.Query != "'SELECT * FROM ' || quote_ident(tbl)" {
		t.Errorf("Query = %q, want %q", exec.Query, "'SELECT * FROM ' || quote_ident(tbl)")
	}
}

func TestExecuteFormat(t *testing.T) {
	// EXECUTE with format()
	block := parseOK(t, `BEGIN EXECUTE format('SELECT * FROM %I', tbl); END`)
	exec := block.Body[0].(*ast.PLDynExecute)
	if exec.Query != "format('SELECT * FROM %I', tbl)" {
		t.Errorf("Query = %q, want %q", exec.Query, "format('SELECT * FROM %I', tbl)")
	}
}

func TestExecuteDuplicateIntoError(t *testing.T) {
	// Duplicate INTO produces error
	parseErr(t, `BEGIN EXECUTE 'SELECT 1' INTO x INTO y; END`, "duplicate INTO")
}

func TestExecuteDuplicateUsingError(t *testing.T) {
	// Duplicate USING produces error
	parseErr(t, `BEGIN EXECUTE 'SELECT $1' USING a USING b; END`, "duplicate USING")
}
