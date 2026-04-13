package parser

import "github.com/bytebase/omni/doris/ast"

// Diagnostic represents a syntax diagnostic (error or warning) with position.
type Diagnostic struct {
	Msg string
	Loc ast.Loc
}

// Diagnose parses the input and returns any syntax errors as diagnostics.
// It calls Parse internally; all ParseErrors (including lex errors promoted
// during parseSingle) are converted to Diagnostics and returned.
//
// Note: while statement-level parsing is stubbed, valid SQL input produces
// "not yet supported" diagnostics. Those disappear as Tier 1+ nodes implement
// real parsing. Lexer errors (unterminated strings, unknown characters) are
// always genuine diagnostics.
func Diagnose(input string) []Diagnostic {
	_, errs := Parse(input)
	diags := make([]Diagnostic, len(errs))
	for i, e := range errs {
		diags[i] = Diagnostic{Msg: e.Msg, Loc: e.Loc}
	}
	return diags
}
