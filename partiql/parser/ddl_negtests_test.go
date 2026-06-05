package parser

import (
	"strings"
	"testing"
)

// TestParser_DDLNameStructureReject is a negative/coverage backfill for the
// DDL name and structure contracts of createCommand/dropCommand. Each input
// violates the grammar and MUST be rejected.
//
// Grammar (PartiQLParser.g4):
//
//	createCommand
//	  : CREATE TABLE symbolPrimitive                                                            # CreateTable
//	  | CREATE INDEX ON symbolPrimitive PAREN_LEFT pathSimple ( COMMA pathSimple )* PAREN_RIGHT # CreateIndex
//	dropCommand
//	  : DROP TABLE symbolPrimitive                       # DropTable
//	  | DROP INDEX target ON on                          # DropIndex
//	symbolPrimitive : ident=( IDENTIFIER | IDENTIFIER_QUOTED ) ;
//	pathSimple      : symbolPrimitive pathSimpleSteps* ;
//	pathSimpleSteps : BRACKET_LEFT key=literal BRACKET_RIGHT          # PathSimpleLiteral
//	                | BRACKET_LEFT key=symbolPrimitive BRACKET_RIGHT  # PathSimpleSymbol
//	                | PERIOD key=symbolPrimitive                      # PathSimpleDotSymbol ;
//	script          : root (COLON_SEMI root)* COLON_SEMI? EOF ;
//
// Oracle: the executable generated ANTLR parser (bytebase/parser/partiql),
// driven through the full Script() entrypoint, REJECTS every input below
// (verified differentially in a throwaway /tmp probe; all five families plus
// variants agreed reject-vs-reject with omni's Parse()/ParseStatement(), and
// the corresponding positive forms agreed accept-vs-accept). No divergences.
//
// "reject" = ParseStatement() returns a non-nil error. ParseStatement is the
// single-statement entrypoint and asserts EOF, so it is the closest analog of
// the ANTLR Script() rule on a one-statement input — it is what surfaces the
// trailing-token failures below. The wantErrIn substrings are stable fragments
// of the current omni error text; they pin the failure to the intended cause
// (e.g. the identifier contract vs. an unrelated downstream error).
func TestParser_DDLNameStructureReject(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		// --- CREATE TABLE name contracts: symbolPrimitive is a SINGLE
		//     IDENTIFIER / IDENTIFIER_QUOTED. A dotted/qualified name is not a
		//     symbolPrimitive; the trailing `.b` cannot extend the statement,
		//     so Script() (and ParseStatement, which asserts EOF) rejects. ---
		{
			// CREATE TABLE a.b — qualified name not allowed: `a` is the table
			// name, then `.b` is unconsumed input after a complete statement.
			name:      "create_table_qualified_dotted",
			input:     "CREATE TABLE a.b",
			wantErrIn: "unexpected token",
		},
		{
			// CREATE TABLE a.b.c — three-part qualified name, same contract.
			name:      "create_table_qualified_three_part",
			input:     "CREATE TABLE a.b.c",
			wantErrIn: "unexpected token",
		},
		{
			// CREATE TABLE "a"."b" — even quoted identifiers cannot be dotted
			// in a symbolPrimitive; the `.` after `"a"` is trailing garbage.
			name:      "create_table_qualified_quoted",
			input:     `CREATE TABLE "a"."b"`,
			wantErrIn: "unexpected token",
		},
		{
			// CREATE TABLE 'foo' — a string literal is not an IDENTIFIER /
			// IDENTIFIER_QUOTED, so it is not a valid symbolPrimitive.
			name:      "create_table_string_name",
			input:     "CREATE TABLE 'foo'",
			wantErrIn: "expected identifier",
		},

		// --- DROP TABLE shares the same symbolPrimitive name contract. ---
		{
			// DROP TABLE 'foo' — string literal is not a valid table name.
			name:      "drop_table_string_name",
			input:     "DROP TABLE 'foo'",
			wantErrIn: "expected identifier",
		},

		// --- CREATE INDEX path-list structure. ---
		{
			// CREATE INDEX ON t () — the path list requires at least one
			// pathSimple (grammar: PAREN_LEFT pathSimple ( COMMA pathSimple )*
			// PAREN_RIGHT). An empty list has no leading symbolPrimitive, so
			// the very first pathSimple fails on the `)`.
			name:      "create_index_empty_path_list",
			input:     "CREATE INDEX ON t ()",
			wantErrIn: "expected identifier",
		},
		{
			// CREATE INDEX ON t (a[) — `a` parses as the pathSimple root, then
			// `[` opens a pathSimpleSteps that requires key=literal or
			// key=symbolPrimitive before BRACKET_RIGHT. The `)` is neither a
			// literal nor an identifier, so the bracket key (and thus the whole
			// statement) is rejected — an unclosed bracket.
			name:      "create_index_unclosed_bracket",
			input:     "CREATE INDEX ON t (a[)",
			wantErrIn: "expected identifier",
		},
		{
			// CREATE INDEX ON 'foo' (a) — the indexed table name is a
			// symbolPrimitive too; a string literal is rejected there.
			name:      "create_index_string_table_name",
			input:     "CREATE INDEX ON 'foo' (a)",
			wantErrIn: "expected identifier",
		},

		// --- Trailing token after an otherwise-complete DDL statement. The
		//     script rule terminates each statement at EOF (or COLON_SEMI);
		//     a bare extra token is not a valid continuation. ---
		{
			// CREATE TABLE t foo — `CREATE TABLE t` is complete; `foo` is an
			// extra identifier after the statement with no separating `;`.
			name:      "trailing_token_after_create_table",
			input:     "CREATE TABLE t foo",
			wantErrIn: "unexpected token",
		},
		{
			// DROP TABLE t foo — same trailing-token contract for DROP.
			name:      "trailing_token_after_drop_table",
			input:     "DROP TABLE t foo",
			wantErrIn: "unexpected token",
		},
		{
			// CREATE INDEX ON t (a) foo — complete CREATE INDEX followed by a
			// stray identifier.
			name:      "trailing_token_after_create_index",
			input:     "CREATE INDEX ON t (a) foo",
			wantErrIn: "unexpected token",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseStatement()
			if err == nil {
				t.Fatalf("expected error for %q, got nil", tc.input)
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}
