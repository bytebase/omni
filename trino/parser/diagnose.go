package parser

import "github.com/bytebase/omni/trino/ast"

// Diagnostic represents one syntax diagnostic (error) with its source position.
// It is the parser-facing shape that bytebase's Diagnose handler maps to LSP
// diagnostics, converting Loc byte offsets to line/column against the source.
type Diagnostic struct {
	Msg string
	Loc ast.Loc
}

// Diagnose parses input and returns every syntax error as a Diagnostic. It
// drives the public Parse, so it inherits Parse's best-effort behavior: errors
// from every statement segment are collected (not just the first), and
// LexErrors promoted during parsing are included.
//
// This backs bytebase's RegisterDiagnoseFunc consumer, which runs on every
// editor keystroke and must not emit false positives for valid Trino. While
// statement bodies are stubbed in the parser-foundation node, syntactically
// valid statements still produce a "not yet supported" diagnostic; those
// disappear as later DAG nodes (types, expressions, select, ddl, dml, …)
// implement real statement parsing. Genuine lexer errors (unterminated
// string/identifier/comment) and unknown-statement errors are always reported.
func Diagnose(input string) []Diagnostic {
	_, err := Parse(input)
	errs := AllErrors(err)
	if len(errs) == 0 {
		return nil
	}
	diags := make([]Diagnostic, len(errs))
	for i, e := range errs {
		diags[i] = Diagnostic{Msg: e.Message, Loc: ast.Loc{Start: e.Position, End: e.End}}
	}
	return diags
}
