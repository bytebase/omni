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

// TestParse_TrailingLexerError verifies that the script-level entry point
// (Parse -> parseScript) surfaces trailing garbage that lexes to a LEXER error
// rather than silently dropping it. parseScript calls checkLexerErr() at entry
// but, before the fix, never after: once the trailing '#' / unterminated '/*'
// set p.lexer.Err, advance() yields tokEOF, the post-parse `cur==tokEOF` check
// passes, and the pending lexer error was lost. Oracle (antlr_fallback,
// PartiQLLexer.g4): '#' matches only UNRECOGNIZED (default channel) and an
// unterminated '/*' cannot match COMMENT_BLOCK, so ANTLR rejects both; closed
// comments and whitespace are HIDDEN-channel and must still be accepted.
func TestParse_TrailingLexerError(t *testing.T) {
	bases := []string{
		"SELECT a FROM t",
		"CREATE TABLE t",
		"DROP TABLE t",
		"INSERT INTO t VALUE 1",
		"UPDATE t SET x = 1",
		"DELETE FROM t",
	}
	for _, base := range bases {
		t.Run(base+" #", func(t *testing.T) {
			input := base + " #"
			_, err := Parse(input)
			if err == nil {
				t.Fatalf("Parse(%q) = nil error, want trailing-lexer-error rejection", input)
			}
			if !strings.Contains(err.Error(), "unexpected character") {
				t.Errorf("Parse(%q) error = %q, want to contain %q", input, err.Error(), "unexpected character")
			}
		})
		t.Run(base+" /*", func(t *testing.T) {
			input := base + " /*"
			_, err := Parse(input)
			if err == nil {
				t.Fatalf("Parse(%q) = nil error, want trailing-lexer-error rejection", input)
			}
			if !strings.Contains(err.Error(), "unterminated block comment") {
				t.Errorf("Parse(%q) error = %q, want to contain %q", input, err.Error(), "unterminated block comment")
			}
		})
	}

	// A lexer error on a trailing statement AFTER a semicolon must also surface
	// (the multi-root loop path), not just on the first statement.
	t.Run("trailing_stmt_after_semicolon", func(t *testing.T) {
		input := "SELECT a FROM t; SELECT b FROM u #"
		_, err := Parse(input)
		if err == nil {
			t.Fatalf("Parse(%q) = nil error, want trailing-lexer-error rejection", input)
		}
		if !strings.Contains(err.Error(), "unexpected character") {
			t.Errorf("Parse(%q) error = %q, want to contain %q", input, err.Error(), "unexpected character")
		}
	})

	// Garbage IMMEDIATELY after a semicolon exercises the loop's
	// trailing-semicolon `break` path (the lexer error is set while scanning
	// for the next root, which yields tokEOF and breaks the loop): the
	// post-loop checkLexerErr must still surface it.
	t.Run("garbage_immediately_after_semicolon", func(t *testing.T) {
		for _, tc := range []struct{ input, want string }{
			{"SELECT a FROM t ; #", "unexpected character"},
			{"SELECT a FROM t ; /*", "unterminated block comment"},
		} {
			_, err := Parse(tc.input)
			if err == nil {
				t.Fatalf("Parse(%q) = nil error, want trailing-lexer-error rejection", tc.input)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("Parse(%q) error = %q, want to contain %q", tc.input, err.Error(), tc.want)
			}
		}
	})

	// Valid trailing whitespace / closed comments and a trailing semicolon must
	// STILL parse cleanly (HIDDEN-channel content must not regress).
	t.Run("valid_trailing_accepted", func(t *testing.T) {
		for _, input := range []string{
			"SELECT a FROM t -- ok",
			"SELECT a FROM t /* closed */",
			"SELECT a FROM t;",
			"SELECT a FROM t; -- ok",
			"SELECT a FROM t   ",
		} {
			if _, err := Parse(input); err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", input, err)
			}
		}
	})
}
