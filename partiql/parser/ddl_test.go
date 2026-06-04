package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// TestParser_CreateIndexBracketKeyLiteral verifies that every literal form
// permitted by the grammar's bracket-key path step is accepted in a
// CREATE INDEX path list.
//
// Grammar (PartiQLParser.g4 line 114):
//
//	pathSimpleSteps : BRACKET_LEFT key=literal BRACKET_RIGHT  # PathSimpleLiteral
//
// where `literal` (g4 lines 661-672) is NULL | MISSING | TRUE | FALSE |
// LITERAL_STRING | LITERAL_INTEGER | LITERAL_DECIMAL | ION_CLOSURE |
// DATE LITERAL_STRING | TIME (...) LITERAL_STRING — i.e. far more than the
// string/integer/decimal subset the parser previously routed to
// parseLiteral.
//
// Oracle: the executable generated ANTLR parser (bytebase/parser/partiql)
// ACCEPTS every accept-case below and REJECTS every reject-case (verified
// differentially via the full Script() entrypoint). These were the exact
// forms the omni parser wrongly rejected with "expected identifier, got
// null" because the bracket-key gate only admitted SCONST/ICONST/FCONST.
func TestParser_CreateIndexBracketKeyLiteral(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		// --- newly-accepted non-numeric/non-string literal keys ---
		{
			name:  "null_key",
			input: "CREATE INDEX ON t (a[null])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:NullLit{}}]}]}",
		},
		{
			name:  "true_key",
			input: "CREATE INDEX ON t (a[true])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:BoolLit{Val:true}}]}]}",
		},
		{
			name:  "false_key",
			input: "CREATE INDEX ON t (a[false])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:BoolLit{Val:false}}]}]}",
		},
		{
			name:  "missing_key",
			input: "CREATE INDEX ON t (a[missing])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:MissingLit{}}]}]}",
		},
		{
			name:  "ion_int_key",
			input: "CREATE INDEX ON t (a[`42`])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:IonLit{Text:\"42\"}}]}]}",
		},
		{
			name:  "ion_struct_key",
			input: "CREATE INDEX ON t (a[`{x:1}`])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:IonLit{Text:\"{x:1}\"}}]}]}",
		},
		{
			name:  "null_key_then_dot",
			input: "CREATE INDEX ON t (a[null].b)",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:NullLit{}} DotStep{Field:b}]}]}",
		},
		{
			name:  "dot_then_null_key",
			input: "CREATE INDEX ON t (a.b[null])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[DotStep{Field:b} IndexStep{Index:NullLit{}}]}]}",
		},
		// --- DATE / TIME literal keys (grammar permits; oracle accepts) ---
		{
			name:  "date_key",
			input: "CREATE INDEX ON t (a[DATE '2020-01-01'])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:DateLit{Val:2020-01-01}}]}]}",
		},
		{
			name:  "time_key",
			input: "CREATE INDEX ON t (a[TIME '12:00:00'])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:TimeLit{Val:12:00:00}}]}]}",
		},
		// --- regression guards: forms already accepted pre-fix ---
		{
			name:  "string_key",
			input: "CREATE INDEX ON t (a['k'])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:StringLit{Val:\"k\"}}]}]}",
		},
		{
			name:  "int_key",
			input: "CREATE INDEX ON t (a[42])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:NumberLit{Val:42}}]}]}",
		},
		{
			name:  "decimal_key",
			input: "CREATE INDEX ON t (a[3.14])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:NumberLit{Val:3.14}}]}]}",
		},
		{
			name:  "symbol_key",
			input: "CREATE INDEX ON t (a[b])",
			want:  "CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[IndexStep{Index:VarRef{Name:b}}]}]}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tc.want {
				t.Errorf("AST mismatch\ngot:  %s\nwant: %s", got, tc.want)
			}
		})
	}
}

// TestParser_CreateIndexBracketKeyReject verifies that the bracket-key path
// step still rejects everything the grammar forbids: a non-literal/
// non-symbol bracket key is not a valid `pathSimpleSteps`. Oracle (the
// generated ANTLR parser) rejects every input below.
//
// Negative coverage is required: routing the full literal set through
// parseLiteral must NOT loosen the gate into accepting general expressions
// (e.g. `a[1+1]`, `a[(null)]`), an empty bracket, or a TIMESTAMP literal
// (which is not a member of the `literal` rule).
func TestParser_CreateIndexBracketKeyReject(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{
			name:      "empty_bracket",
			input:     "CREATE INDEX ON t (a[])",
			wantErrIn: "expected",
		},
		{
			name:      "missing_rbracket",
			input:     "CREATE INDEX ON t (a[null)",
			wantErrIn: "BRACKET_RIGHT",
		},
		{
			name:      "expr_not_literal",
			input:     "CREATE INDEX ON t (a[1+1])",
			wantErrIn: "BRACKET_RIGHT",
		},
		{
			name:      "paren_literal",
			input:     "CREATE INDEX ON t (a[(null)])",
			wantErrIn: "expected",
		},
		{
			// TIMESTAMP is intentionally NOT a member of the literal rule;
			// parseLiteral rejects it with a precise message. Oracle also
			// rejects `a[TIMESTAMP '...']` (no viable alternative).
			name:      "timestamp_key",
			input:     "CREATE INDEX ON t (a[TIMESTAMP '2020-01-01 00:00:00'])",
			wantErrIn: "TIMESTAMP literal is not supported",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseStatement()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}
