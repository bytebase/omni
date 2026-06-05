package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// This file is a TEST-ONLY negative/coverage backfill for the predicate and
// generic-function-call layers (DAG nodes foundation-builtins /
// parser-builtins-generic-call). Every case below was verified against the
// executable generated ANTLR parser at github.com/bytebase/parser/partiql via
// its top-level Script() rule (the migration oracle). For the reject cases the
// fragment was embedded as `SELECT 1 FROM t WHERE <frag>`; ANTLR reports a
// syntax error for all five, matching omni. The positive case was embedded as
// `SELECT "foo"(x) FROM t`; ANTLR accepts it, matching omni. No divergences
// were found — all cases assert the agreed-upon correct behavior.

// TestPredicateMissingArg_Reject pins the predicate forms whose required
// argument is absent at end-of-input. These exercise the error paths of
// parsePredicate's sub-parsers (parseLikeBody / parseIsBody / parseBetweenBody
// / parseInBody in expr.go) when the trailing operand never arrives.
//
// Oracle (executable ANTLR, via Script() on `SELECT 1 FROM t WHERE <frag>`):
//
//	a LIKE              -> reject: mismatched input '<EOF>' expecting <primary first-set>
//	a LIKE 'x' ESCAPE   -> reject: mismatched input '<EOF>' expecting <expr first-set>
//	a IS                -> reject: mismatched input '<EOF>' expecting <type first-set>
//	a BETWEEN 1         -> reject: mismatched input '<EOF>' expecting 'AND'
//	a IN                -> reject: no viable alternative at input 'IN'
//
// omni rejects all five as well (verified differentially). The asserted
// substrings are omni's own messages; where omni's wording differs from
// ANTLR's (e.g. omni reaches the primary-expression error for `a IN` rather
// than ANTLR's "no viable alternative"), the *verdict* still matches — both
// reject — so this is a faithful, non-weakened negative test, not a divergence.
func TestPredicateMissingArg_Reject(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
		// note documents the omni code path that produces the error, for the
		// next maintainer who greps for these.
	}{
		{
			// parseLikeBody consumes LIKE, then parseMathOp00 -> primary at EOF.
			name:      "like_no_pattern",
			input:     "a LIKE",
			wantErrIn: `unexpected token "" in expression`,
		},
		{
			// parseLikeBody consumes ESCAPE, then parseMathOp00 -> primary at EOF.
			name:      "like_escape_no_arg",
			input:     "a LIKE 'x' ESCAPE",
			wantErrIn: `unexpected token "" in expression`,
		},
		{
			// parseIsBody -> parseType at EOF: "expected type, got ...".
			name:      "is_eof",
			input:     "a IS",
			wantErrIn: "expected type",
		},
		{
			// parseBetweenBody parses lower=1, then expect(AND) fails at EOF.
			name:      "between_no_and",
			input:     "a BETWEEN 1",
			wantErrIn: "expected AND",
		},
		{
			// parseInBody consumes IN; no '(' so it takes the expression form,
			// parseMathOp00 -> primary at EOF.
			name:      "in_no_list",
			input:     "a IN",
			wantErrIn: `unexpected token "" in expression`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			expr, err := p.ParseExpr()
			if err == nil {
				t.Fatalf("ParseExpr(%q) = %s, nil error; want error containing %q",
					tc.input, ast.NodeToString(expr), tc.wantErrIn)
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("ParseExpr(%q) error = %q, want to contain %q",
					tc.input, err.Error(), tc.wantErrIn)
			}
		})
	}
}

// TestGenericFuncCall_QuotedName_Accept covers a generic function call whose
// name is a double-quoted (case-sensitive) identifier: `"foo"(x)`. This flips
// the parser-builtins-generic-call node's quoted-name branch to fully covered.
//
// In the grammar a function call is `name=symbolPrimitive '(' ... ')'`
// (PartiQLParser.g4:615 functionCall#FunctionCallIdent), and symbolPrimitive
// (g4:742) admits IDENTIFIER_QUOTED. omni's parseVarRef detects the trailing
// PAREN_LEFT after parseSymbolPrimitive and emits an *ast.FuncCall. The lexer
// strips the surrounding quotes (scanQuotedIdent), so the FuncCall Name is the
// decoded text `foo`.
//
// Oracle (executable ANTLR, via Script() on `SELECT "foo"(x) FROM t`): accepts
// with zero lexer/parser syntax errors. omni accepts and renders
// `FuncCall{Name:foo Args:[VarRef{Name:x}]}` (verified differentially).
func TestGenericFuncCall_QuotedName_Accept(t *testing.T) {
	const input = `"foo"(x)`
	const want = "FuncCall{Name:foo Args:[VarRef{Name:x}]}"

	p := NewParser(input)
	expr, err := p.ParseExpr()
	if err != nil {
		t.Fatalf("ParseExpr(%q) unexpected error: %v", input, err)
	}
	got := ast.NodeToString(expr)
	if got != want {
		t.Errorf("ParseExpr(%q)\n got: %s\nwant: %s", input, got, want)
	}
}
