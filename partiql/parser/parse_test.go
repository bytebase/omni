package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// TestParse exercises the public Parse function (script-level entry point).
func TestParse(t *testing.T) {
	type tc struct {
		name    string
		input   string
		wantLen int
		wantStr string // substring that NodeToString of the list must contain; empty means skip
		wantErr string // non-empty: error must contain this substring
	}

	cases := []tc{
		// -----------------------------------------------------------------------
		// Empty input
		// -----------------------------------------------------------------------
		{
			name:    "empty_input",
			input:   "",
			wantLen: 0,
		},
		{
			name:    "whitespace_only",
			input:   "   ",
			wantLen: 0,
		},

		// -----------------------------------------------------------------------
		// Single statements
		// -----------------------------------------------------------------------
		{
			name:    "single_select_star",
			input:   "SELECT * FROM t",
			wantLen: 1,
			wantStr: "SelectStmt{Star:true",
		},
		{
			name:    "single_select_items",
			input:   "SELECT a, b FROM t",
			wantLen: 1,
			wantStr: "SelectStmt{",
		},
		{
			name:    "single_create_table",
			input:   "CREATE TABLE foo",
			wantLen: 1,
			wantStr: "CreateTableStmt{",
		},
		{
			name:    "single_insert",
			input:   "INSERT INTO t VALUE 1",
			wantLen: 1,
			wantStr: "InsertStmt{",
		},

		// -----------------------------------------------------------------------
		// Trailing semicolon
		// -----------------------------------------------------------------------
		{
			name:    "trailing_semicolon",
			input:   "SELECT * FROM t;",
			wantLen: 1,
		},

		// -----------------------------------------------------------------------
		// Multi-statement script
		// -----------------------------------------------------------------------
		{
			name:    "two_selects",
			input:   "SELECT * FROM a; SELECT * FROM b",
			wantLen: 2,
		},
		{
			name:    "three_stmts_trailing_semi",
			input:   "SELECT * FROM a; SELECT * FROM b; SELECT * FROM c;",
			wantLen: 3,
		},

		// -----------------------------------------------------------------------
		// EXPLAIN wrapper
		// -----------------------------------------------------------------------
		{
			name:    "explain_select",
			input:   "EXPLAIN SELECT * FROM t",
			wantLen: 1,
			wantStr: "ExplainStmt{Inner:SelectStmt{",
		},
		{
			name:    "explain_with_options",
			input:   "EXPLAIN (format text) SELECT * FROM t",
			wantLen: 1,
			wantStr: "ExplainStmt{Inner:SelectStmt{",
		},

		// -----------------------------------------------------------------------
		// EXEC command
		// -----------------------------------------------------------------------
		{
			name:    "exec_no_args",
			input:   "EXEC myProc",
			wantLen: 1,
			wantStr: "ExecStmt{Name:VarRef{Name:myProc} Args:[]}",
		},
		{
			name:    "exec_with_args",
			input:   "EXEC myProc 1, 'hello'",
			wantLen: 1,
			wantStr: "ExecStmt{Name:VarRef{Name:myProc} Args:[NumberLit{Val:1} StringLit{Val:\"hello\"}]}",
		},
		{
			name:    "execute_keyword",
			input:   "EXECUTE myProc",
			wantLen: 1,
			wantStr: "ExecStmt{",
		},

		// -----------------------------------------------------------------------
		// EXEC via ParseStatement (single-stmt path)
		// -----------------------------------------------------------------------
		{
			name:    "parse_statement_exec",
			input:   "EXEC myProc",
			wantLen: 1,
			wantStr: "ExecStmt{",
		},

		// -----------------------------------------------------------------------
		// Mixed script with EXPLAIN
		// -----------------------------------------------------------------------
		{
			name:    "explain_then_select",
			input:   "EXPLAIN SELECT * FROM t; SELECT * FROM u",
			wantLen: 2,
			wantStr: "ExplainStmt{",
		},

		// -----------------------------------------------------------------------
		// Error cases
		// -----------------------------------------------------------------------
		{
			name:    "bare_semicolon_only",
			input:   ";",
			wantErr: "unexpected token",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			list, err := Parse(c.input)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("Parse(%q) = nil error, want error containing %q", c.input, c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("Parse(%q) error = %q, want containing %q", c.input, err.Error(), c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", c.input, err)
			}
			if list == nil {
				t.Fatalf("Parse(%q) returned nil list", c.input)
			}
			if list.Len() != c.wantLen {
				t.Errorf("Parse(%q): len = %d, want %d", c.input, list.Len(), c.wantLen)
			}
			if c.wantStr != "" {
				got := ast.NodeToString(list)
				if !strings.Contains(got, c.wantStr) {
					t.Errorf("Parse(%q): NodeToString does not contain %q\ngot: %s", c.input, c.wantStr, got)
				}
			}
		})
	}
}

// TestParseStatement_Exec verifies the single-statement path for EXEC
// through the existing ParseStatement entry point.
func TestParseStatement_Exec(t *testing.T) {
	cases := []struct {
		input   string
		wantStr string
	}{
		{
			input:   "EXEC myProc",
			wantStr: "ExecStmt{Name:VarRef{Name:myProc} Args:[]}",
		},
		{
			input:   "EXEC myProc 1, 2",
			wantStr: "ExecStmt{Name:VarRef{Name:myProc} Args:[NumberLit{Val:1} NumberLit{Val:2}]}",
		},
		{
			input:   "EXECUTE proc123",
			wantStr: "ExecStmt{Name:VarRef{Name:proc123} Args:[]}",
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			p := NewParser(c.input)
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement(%q) error: %v", c.input, err)
			}
			got := ast.NodeToString(stmt)
			if !strings.Contains(got, c.wantStr) {
				t.Errorf("ParseStatement(%q):\n got: %s\nwant (containing): %s", c.input, got, c.wantStr)
			}
		})
	}
}
