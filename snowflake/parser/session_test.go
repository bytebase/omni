package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// mustParseOne parses input as a single statement, fails the test if any parse
// error occurred or if the statement count is not exactly 1, and returns the
// lone statement node. Shared by the T6.3 utility-statement tests
// (session/show/comment_truncate).
func mustParseOne(t *testing.T, input string) ast.Node {
	t.Helper()
	result := ParseBestEffort(input)
	if len(result.Errors) > 0 {
		t.Fatalf("parse %q: %d error(s): %v", input, len(result.Errors), result.Errors)
	}
	if len(result.File.Stmts) != 1 {
		t.Fatalf("parse %q: got %d statements, want 1", input, len(result.File.Stmts))
	}
	return result.File.Stmts[0]
}

// ---------------------------------------------------------------------------
// USE
// ---------------------------------------------------------------------------

func TestParseUse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantKind ast.UseTargetKind
		wantName string // ObjectName.String(); "" when SecondaryRoles
		wantSec  string // SecondaryRoles value; "" otherwise
	}{
		{"use database", "USE DATABASE obj", ast.UseDatabase, "obj", ""},
		{"use role", "USE ROLE obj", ast.UseRole, "obj", ""},
		{"use schema", "USE SCHEMA obj", ast.UseSchema, "obj", ""},
		{"use schema qualified", "USE SCHEMA sch.obj", ast.UseSchema, "sch.obj", ""},
		{"use bare schema", "USE sch.obj", ast.UseDefault, "sch.obj", ""},
		{"use bare db", "USE mydb", ast.UseDefault, "mydb", ""},
		{"use warehouse", "USE WAREHOUSE obj", ast.UseWarehouse, "obj", ""},
		{"use secondary roles all", "USE SECONDARY ROLES ALL", ast.UseSecondaryRoles, "", "ALL"},
		{"use secondary roles none", "USE SECONDARY ROLES NONE", ast.UseSecondaryRoles, "", "NONE"},
		{"use quoted db", `USE DATABASE "My DB"`, ast.UseDatabase, `"My DB"`, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := mustParseOne(t, tt.input)
			stmt, ok := node.(*ast.UseStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.UseStmt", node)
			}
			if stmt.Kind != tt.wantKind {
				t.Errorf("Kind = %v, want %v", stmt.Kind, tt.wantKind)
			}
			if tt.wantName != "" {
				if stmt.Name == nil {
					t.Fatalf("Name is nil, want %q", tt.wantName)
				}
				if got := stmt.Name.String(); got != tt.wantName {
					t.Errorf("Name = %q, want %q", got, tt.wantName)
				}
			}
			if stmt.SecondaryRoles != tt.wantSec {
				t.Errorf("SecondaryRoles = %q, want %q", stmt.SecondaryRoles, tt.wantSec)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// SET
// ---------------------------------------------------------------------------

func TestParseSet(t *testing.T) {
	t.Run("single scalar", func(t *testing.T) {
		node := mustParseOne(t, "SET V1 = 10")
		stmt, ok := node.(*ast.SetStmt)
		if !ok {
			t.Fatalf("got %T, want *ast.SetStmt", node)
		}
		if stmt.Paren {
			t.Errorf("Paren = true, want false")
		}
		if len(stmt.Vars) != 1 {
			t.Fatalf("Vars len = %d, want 1", len(stmt.Vars))
		}
		if stmt.Vars[0].Name.Name != "V1" {
			t.Errorf("Vars[0].Name = %q, want V1", stmt.Vars[0].Name.Name)
		}
		lit, ok := stmt.Vars[0].Value.(*ast.Literal)
		if !ok {
			t.Fatalf("Vars[0].Value = %T, want *ast.Literal", stmt.Vars[0].Value)
		}
		if lit.Kind != ast.LitInt || lit.Ival != 10 {
			t.Errorf("Value lit = %+v, want LitInt 10", lit)
		}
	})

	t.Run("single string", func(t *testing.T) {
		node := mustParseOne(t, "SET V2 = 'example'")
		stmt := node.(*ast.SetStmt)
		if len(stmt.Vars) != 1 {
			t.Fatalf("Vars len = %d, want 1", len(stmt.Vars))
		}
		lit := stmt.Vars[0].Value.(*ast.Literal)
		if lit.Kind != ast.LitString {
			t.Errorf("Value kind = %v, want LitString", lit.Kind)
		}
	})

	t.Run("single subquery", func(t *testing.T) {
		node := mustParseOne(t, "SET id_threshold = (SELECT COUNT(*)/2 FROM table1)")
		stmt := node.(*ast.SetStmt)
		if len(stmt.Vars) != 1 {
			t.Fatalf("Vars len = %d, want 1", len(stmt.Vars))
		}
		if stmt.Vars[0].Value == nil {
			t.Errorf("Value is nil, want a parenthesized subquery expression")
		}
	})

	t.Run("paren multi", func(t *testing.T) {
		node := mustParseOne(t, "SET (V1, V2) = (10, 'example')")
		stmt := node.(*ast.SetStmt)
		if !stmt.Paren {
			t.Errorf("Paren = false, want true")
		}
		if len(stmt.Vars) != 2 {
			t.Fatalf("Vars len = %d, want 2", len(stmt.Vars))
		}
		if stmt.Vars[0].Name.Name != "V1" || stmt.Vars[1].Name.Name != "V2" {
			t.Errorf("names = %q,%q want V1,V2", stmt.Vars[0].Name.Name, stmt.Vars[1].Name.Name)
		}
	})

	t.Run("paren expr rhs", func(t *testing.T) {
		// Each value may be a full expression. ($var session-variable refs are a
		// separate matter handled by the expression parser; see the divergence
		// note below.)
		node := mustParseOne(t, "SET (min, max) = (40, 2 * 35)")
		stmt := node.(*ast.SetStmt)
		if len(stmt.Vars) != 2 {
			t.Fatalf("Vars len = %d, want 2", len(stmt.Vars))
		}
	})
}

// TestParseSet_DollarVarLimitation documents that the SET value (an expression)
// cannot yet contain a $var session-variable reference, because the shared
// expression parser (T3) does not implement the tokVariable primary. The SET
// parser itself is correct — it delegates value parsing to parseExpr — so when
// the expression parser gains $var support, these forms will parse with no
// change to session.go. Flagged divergence: official set/example_06.sql
// ("SET (min, max) = (50, 2 * $min)") currently errors for this reason, NOT
// because of the SET grammar.
func TestParseSet_DollarVarLimitation(t *testing.T) {
	result := ParseBestEffort("SET (min, max) = (50, 2 * $min)")
	if len(result.Errors) == 0 {
		t.Skip("expression parser now supports $var — update the SET corpus filter to include set/example_06.sql")
	}
}

// TestParseSession_Loc verifies the statement Loc spans the full statement
// text. Loc is load-bearing for downstream query-span and diagnostics.
func TestParseSession_Loc(t *testing.T) {
	cases := []struct {
		input string
		loc   func(ast.Node) ast.Loc
	}{
		{"USE WAREHOUSE wh", func(n ast.Node) ast.Loc { return n.(*ast.UseStmt).Loc }},
		{"SET (V1, V2) = (10, 'example')", func(n ast.Node) ast.Loc { return n.(*ast.SetStmt).Loc }},
		{"UNSET (V1, V2)", func(n ast.Node) ast.Loc { return n.(*ast.UnsetStmt).Loc }},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			loc := c.loc(mustParseOne(t, c.input))
			if loc.Start != 0 || loc.End != len(c.input) {
				t.Errorf("Loc = {%d,%d}, want {0,%d}", loc.Start, loc.End, len(c.input))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// UNSET
// ---------------------------------------------------------------------------

func TestParseUnset(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		node := mustParseOne(t, "UNSET V1")
		stmt, ok := node.(*ast.UnsetStmt)
		if !ok {
			t.Fatalf("got %T, want *ast.UnsetStmt", node)
		}
		if stmt.Paren {
			t.Errorf("Paren = true, want false")
		}
		if len(stmt.Names) != 1 || stmt.Names[0].Name != "V1" {
			t.Errorf("Names = %+v, want [V1]", stmt.Names)
		}
	})

	t.Run("paren multi", func(t *testing.T) {
		node := mustParseOne(t, "UNSET (V1, V2)")
		stmt := node.(*ast.UnsetStmt)
		if !stmt.Paren {
			t.Errorf("Paren = false, want true")
		}
		if len(stmt.Names) != 2 {
			t.Fatalf("Names len = %d, want 2", len(stmt.Names))
		}
		if stmt.Names[0].Name != "V1" || stmt.Names[1].Name != "V2" {
			t.Errorf("Names = %+v, want [V1 V2]", stmt.Names)
		}
	})
}

func TestParseSession_Errors(t *testing.T) {
	cases := []string{
		"USE",                 // nothing after USE
		"USE ROLE",            // ROLE with no name
		"USE SECONDARY ROLES", // missing ALL/NONE
		"SET",                 // nothing after SET
		"SET V1",              // missing = expr
		"SET V1 =",            // missing expr
		"SET (V1, V2) = (10)", // arity mismatch — 2 vars, 1 value
		"SET (V1) =",          // missing value list
		"UNSET",               // nothing after UNSET
		"UNSET ()",            // empty paren list
	}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			result := ParseBestEffort(c)
			if len(result.Errors) == 0 {
				t.Errorf("expected parse error for %q, got none (stmts=%d)", c, len(result.File.Stmts))
			}
		})
	}
}
