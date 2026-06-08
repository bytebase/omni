package parser

import (
	"strings"
	"testing"
)

// Negative / coverage backfill for the public Parse entry point (parse.go:
// Parse -> parseScript -> parseRoot -> parseRootStatement / parseExecCommand).
//
// Oracle (truth2-first): the executable ANTLR grammar at
// /Users/h3n4l/OpenSource/parser/partiql/PartiQLParser.g4. The relevant rules:
//
//	script        : root (COLON_SEMI root)* COLON_SEMI? EOF;
//	root          : (EXPLAIN (PAREN_LEFT explainOption (COMMA explainOption)* PAREN_RIGHT)?)? statement;
//	statement     : dql | dml | ddl | execCommand;
//	explainOption : param=IDENTIFIER value=IDENTIFIER;
//	execCommand   : EXEC name=expr (args+=expr (COMMA args+=expr)*)?;
//
// For these rules "reject" (omni Parse returns a non-nil error) is the
// correct behavior, matching ANTLR which cannot derive any of these inputs.
//
// The trailing-lexer-error path (a trailing '#' / unterminated '/*') is
// already covered by TestParse_TrailingLexerError (PR #201) and is not
// repeated here.
//
// All cases below were verified to reject in omni and to be non-derivable in
// the ANTLR grammar (no divergence found).
func TestParse_EntryRejections(t *testing.T) {
	cases := []struct {
		name  string
		input string
		// want is a substring the error message must contain. Empty means the
		// case is asserted to reject but the exact message is not pinned (used
		// where an unrelated requirement, e.g. SELECT's mandatory FROM clause,
		// fires before the structural error under test — both omni and ANTLR
		// still reject, so the suite stays honest).
		want string
	}{
		// --------------------------------------------------------------------
		// EXEC without a name.
		// execCommand: EXEC name=expr ... — `name=expr` is mandatory, so a bare
		// EXEC (followed only by EOF) has no name expression. parseExecCommand
		// calls parseExprTop, which hits EOF and rejects.
		// --------------------------------------------------------------------
		{
			name:  "exec_no_name",
			input: "EXEC",
			want:  "unexpected token",
		},
		{
			name:  "exec_no_name_trailing_ws",
			input: "EXEC   ",
			want:  "unexpected token",
		},
		{
			name:  "execute_no_name",
			input: "EXECUTE",
			want:  "unexpected token",
		},
		{
			name:  "exec_no_name_then_semicolon",
			input: "EXEC;",
			want:  "unexpected token",
		},

		// --------------------------------------------------------------------
		// Malformed EXPLAIN options.
		// explainOption: param=IDENTIFIER value=IDENTIFIER — exactly two
		// identifiers per option. A trailing example using a *valid* inner
		// statement (SELECT a FROM t) is included for each so the rejection is
		// caused by the option list, not by any later clause.
		// --------------------------------------------------------------------
		{
			// `(foo)` supplies only the first identifier; value=IDENTIFIER is
			// missing, so parsing the option's value sees ')'.
			name:  "explain_one_ident_option",
			input: "EXPLAIN (foo) SELECT a FROM t",
			want:  "expected IDENT",
		},
		{
			// `()` supplies no identifiers at all; the first IDENT is missing.
			name:  "explain_empty_option_list",
			input: "EXPLAIN () SELECT a FROM t",
			want:  "expected IDENT",
		},
		{
			// Missing PAREN_RIGHT: after a complete option `foo bar`, the next
			// token is SELECT instead of ',' or ')'.
			name:  "explain_unterminated_option_list",
			input: "EXPLAIN (foo bar SELECT a FROM t",
			want:  "expected PAREN_RIGHT",
		},
		{
			// `(SELECT 1` — the '(' opens an option list, but SELECT is a
			// reserved keyword and cannot be the param=IDENTIFIER, and the ')'
			// is also absent. The option list rejects first.
			name:  "explain_open_paren_then_select",
			input: "EXPLAIN (SELECT 1",
			want:  "expected IDENT",
		},
		{
			// Task-spec literal forms. SELECT (no FROM) is itself non-derivable
			// in both omni and ANTLR (exprSelect#SfwQuery requires fromClause),
			// but here the explainOption rejects strictly before SELECT is even
			// reached, so the structural error is what fires.
			name:  "explain_one_ident_option_bare_select",
			input: "EXPLAIN (foo) SELECT 1",
			want:  "expected IDENT",
		},
		{
			name:  "explain_empty_option_list_bare_select",
			input: "EXPLAIN () SELECT 1",
			want:  "expected IDENT",
		},

		// --------------------------------------------------------------------
		// Nested EXPLAIN.
		// root: (EXPLAIN (...)? )? statement — at most one EXPLAIN, and what
		// follows must be a `statement`. A second EXPLAIN is not a statement;
		// omni reaches the DQL fallback (parseExprTop) on the inner EXPLAIN
		// keyword and rejects. ANTLR likewise cannot derive `EXPLAIN EXPLAIN`.
		// --------------------------------------------------------------------
		{
			name:  "nested_explain_valid_inner",
			input: "EXPLAIN EXPLAIN SELECT a FROM t",
			want:  "unexpected token \"EXPLAIN\" in expression",
		},
		{
			name:  "nested_explain_bare_select",
			input: "EXPLAIN EXPLAIN SELECT 1",
			want:  "unexpected token \"EXPLAIN\" in expression",
		},

		// --------------------------------------------------------------------
		// Missing semicolon between two statements.
		// script: root (COLON_SEMI root)* ... — a COLON_SEMI is required
		// between roots. With two valid statements and no separator, the first
		// root parses to completion and the trailing tokens of the second
		// statement remain, which parseScript reports as trailing garbage.
		// --------------------------------------------------------------------
		{
			name:  "missing_semicolon_valid_stmts",
			input: "SELECT a FROM t SELECT b FROM u",
			want:  "unexpected token after statement",
		},
		{
			// Task-spec literal form. Here bare SELECT's mandatory FROM clause
			// (shared by omni and ANTLR) is what fails first, before the
			// missing-separator condition is observable; both still reject.
			name:  "missing_semicolon_bare_selects",
			input: "SELECT 1 SELECT 2",
			want:  "", // rejects; exact message not pinned (FROM requirement fires first)
		},

		// --------------------------------------------------------------------
		// Empty middle statement.
		// script: root (COLON_SEMI root)* COLON_SEMI? EOF — between two
		// COLON_SEMI there must be a `root`; an empty root is non-derivable.
		// In omni, after the first ';' the loop sees a non-EOF ';' and calls
		// parseRoot again, which dispatches to the DQL fallback on ';'.
		// --------------------------------------------------------------------
		{
			name:  "empty_middle_statement_valid_stmts",
			input: "SELECT a FROM t;;SELECT b FROM u",
			want:  "unexpected token \";\" in expression",
		},
		{
			// Task-spec literal form. Bare SELECT's FROM requirement fails on
			// the first statement before the empty-middle ';;' is reached; both
			// omni and ANTLR reject the input as a whole.
			name:  "empty_middle_statement_bare_selects",
			input: "SELECT 1;;SELECT 2",
			want:  "", // rejects; exact message not pinned (FROM requirement fires first)
		},

		// --------------------------------------------------------------------
		// EXEC with non-comma-separated arguments.
		// execCommand: EXEC name=expr (args+=expr (COMMA args+=expr)*)? — after
		// the first argument, further arguments must each be preceded by COMMA.
		// `EXEC foo a b` parses name=foo, arg=a, then the second arg `b` lacks
		// the required COMMA, so the argument loop stops and `b` is left as
		// trailing garbage. ANTLR cannot derive a second arg without a COMMA
		// either (script's EOF cannot follow the stray token).
		// --------------------------------------------------------------------
		{
			name:  "exec_args_no_comma_idents",
			input: "EXEC foo a b",
			want:  "unexpected token after statement",
		},
		{
			name:  "exec_args_no_comma_numbers",
			input: "EXEC foo 1 2",
			want:  "unexpected token after statement",
		},
		{
			name:  "exec_args_no_comma_after_two",
			input: "EXEC foo a, b c",
			want:  "unexpected token after statement",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Parse(c.input)
			if err == nil {
				t.Fatalf("Parse(%q) = nil error, want rejection", c.input)
			}
			if c.want != "" && !strings.Contains(err.Error(), c.want) {
				t.Errorf("Parse(%q) error = %q, want containing %q", c.input, err.Error(), c.want)
			}
		})
	}
}

// TestParseStatement_EntryRejections covers the single-statement entry point
// (Parser.ParseStatement) for the EXEC-without-name case. ParseStatement
// dispatches tokEXEC/tokEXECUTE to parseExecCommand just like the script entry
// point, so a missing name expression rejects here as well.
func TestParseStatement_EntryRejections(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "exec_no_name",
			input: "EXEC",
			want:  "unexpected token",
		},
		{
			name:  "execute_no_name",
			input: "EXECUTE",
			want:  "unexpected token",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := NewParser(c.input)
			_, err := p.ParseStatement()
			if err == nil {
				t.Fatalf("ParseStatement(%q) = nil error, want rejection", c.input)
			}
			if c.want != "" && !strings.Contains(err.Error(), c.want) {
				t.Errorf("ParseStatement(%q) error = %q, want containing %q", c.input, err.Error(), c.want)
			}
		})
	}
}
