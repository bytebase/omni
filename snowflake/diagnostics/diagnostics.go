// Package diagnostics provides a thin reporting layer over the Snowflake SQL
// parser, shaping raw parse errors into structured diagnostics suitable for
// editor integration (LSP-style), CI linting, and problem-panel display.
//
// The primary entry point is Analyze, which accepts a SQL string and returns
// a slice of Diagnostic values — one per parse or lex error. Each Diagnostic
// carries a source range (line, column, byte offset) and a human-readable
// message, ready for use in language servers, linters, or test assertions.
//
// See docs/superpowers/specs/2026-04-15-snowflake-diagnostics-design.md
// for the design rationale.
package diagnostics

import (
	"github.com/bytebase/omni/snowflake/parser"
)

// Severity classifies the importance of a diagnostic.
//
// Only SeverityError is emitted today (all parse errors are fatal for
// SQL execution). SeverityWarning and SeverityInfo are reserved for future
// semantic / advisory diagnostics.
type Severity int

const (
	// SeverityError indicates a syntax error that prevents SQL execution.
	SeverityError Severity = iota
	// SeverityWarning indicates a non-fatal issue (reserved for future use).
	SeverityWarning
	// SeverityInfo indicates an informational note (reserved for future use).
	SeverityInfo
)

// String returns a human-readable label for a Severity value.
func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	case SeverityInfo:
		return "info"
	default:
		return "unknown"
	}
}

// Position is a point in the source text identified by line, column, and
// byte offset. All three fields are provided so callers can choose their
// preferred coordinate system.
//
// Line and Column are 1-based (the first character of a file is line 1,
// column 1). Column is measured in bytes, not Unicode code points, matching
// the byte-based tokenization used by the Snowflake lexer.
//
// Offset is 0-based and refers to the byte position within the full input
// string passed to Analyze.
type Position struct {
	Line   int // 1-based line number
	Column int // 1-based column (bytes from line start)
	Offset int // 0-based byte offset within the source
}

// Range is a contiguous region of source text. Start is inclusive;
// End is exclusive — mirroring the ast.Loc convention.
//
// A zero-width range (Start == End) indicates a point diagnostic where
// the error boundary could not be determined (e.g. end-of-input errors).
type Range struct {
	Start Position
	End   Position
}

// Diagnostic is a single structured error or warning produced by the parser.
// It is intentionally shaped after the LSP Diagnostic type so consumers can
// map it to editor underlines, problem-panel entries, or CI annotations with
// minimal transformation.
type Diagnostic struct {
	// Severity is always SeverityError for syntax errors reported by the parser.
	Severity Severity
	// Range identifies the source region to underline.
	Range Range
	// Source identifies the tool that produced this diagnostic.
	// Always "snowflake-parser" for diagnostics from this package.
	Source string
	// Message is the human-readable error description from the parser.
	Message string
}

const source = "snowflake-parser"

// Analyze parses sql using the Snowflake best-effort parser and converts
// every parse or lex error into a Diagnostic with a populated source range.
//
// Returns nil when sql is syntactically valid or empty. The returned slice
// is ordered by the byte offset at which each error was detected.
//
// Line and column numbers are 1-based and measured in bytes. Callers that
// need Unicode-aware column numbers should post-process the Offset field.
func Analyze(sql string) []Diagnostic {
	result := parser.ParseBestEffort(sql)
	if len(result.Errors) == 0 {
		return nil
	}

	lt := parser.NewLineTable(sql)
	diags := make([]Diagnostic, 0, len(result.Errors))

	for _, pe := range result.Errors {
		startOff := pe.Loc.Start
		if startOff < 0 {
			startOff = 0
		}
		startLine, startCol := lt.Position(startOff)

		// End offset may be unknown (-1). Fall back to the start offset so
		// we produce a zero-width (point) diagnostic rather than a garbage range.
		endOff := pe.Loc.End
		if endOff < 0 {
			endOff = startOff
		}
		endLine, endCol := lt.Position(endOff)

		diags = append(diags, Diagnostic{
			Severity: SeverityError,
			Range: Range{
				Start: Position{Line: startLine, Column: startCol, Offset: startOff},
				End:   Position{Line: endLine, Column: endCol, Offset: endOff},
			},
			Source:  source,
			Message: pe.Msg,
		})
	}

	return diags
}
