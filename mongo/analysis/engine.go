package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/mongo"
)

// Engine is MongoDB's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.StatementRanger = Engine
	_ analysis.Diagnoser       = Engine
	_ analysis.Explainer       = Engine
)

// StatementRanges parses best effort, so the statements before and after one
// the user is still typing keep their ranges.
func (engine) StatementRanges(script string) []analysis.LSPRange {
	var ranges []analysis.ByteRange
	for _, statement := range mongo.ParseBestEffort(script).Statements {
		if !statement.Empty() {
			ranges = append(ranges, analysis.ByteRange{Start: statement.ByteStart, End: statement.ByteEnd})
		}
	}
	return analysis.LSPRanges(script, ranges)
}

// Diagnose reports every syntax error a best-effort parse recovers from, not
// only the first.
func (engine) Diagnose(script string) []analysis.Diagnostic {
	var diagnostics []analysis.Diagnostic
	for _, parseErr := range mongo.ParseBestEffort(script).Errors {
		diagnostics = append(diagnostics, analysis.Diagnostic{
			Start:    analysis.Position{Line: parseErr.Line, Column: parseErr.Column},
			Severity: analysis.SeverityError,
			Message:  parseErr.Message,
		})
	}
	return diagnostics
}

// Explain ignores format: MongoDB states a plan as one document.
func (engine) Explain(statement string, _ analysis.ExplainFormat) (string, error) {
	return ExplainStatement(statement)
}
