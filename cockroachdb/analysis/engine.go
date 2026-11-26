package analysis

import (
	"fmt"

	"github.com/bytebase/omni/analysis"
)

// Engine is CockroachDB's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter       = Engine
	_ analysis.Diagnoser      = Engine
	_ analysis.QueryValidator = Engine
	_ analysis.Explainer      = Engine
	_ analysis.ResultLimiter  = Engine
)

// Split cuts the script lexically; a segment's text ends at its own
// semicolon.
func (engine) Split(script string) []analysis.Statement {
	segments := SplitSegments(script)
	if len(segments) == 0 {
		return nil
	}
	statements := make([]analysis.Statement, 0, len(segments))
	for _, s := range segments {
		statements = append(statements, analysis.Statement{
			Text:  s.Text,
			Empty: s.Empty(),
			Start: s.ByteStart,
			End:   s.ByteEnd,
		})
	}
	return statements
}

// Diagnose reports the script's first syntax error.
func (engine) Diagnose(script string) []analysis.Diagnostic {
	diagnostics := []analysis.Diagnostic{}
	err := Diagnose(script)
	if err == nil {
		return diagnostics
	}
	return append(diagnostics, analysis.Diagnostic{
		Start:    analysis.PositionAt(script, err.Offset),
		Severity: analysis.SeverityError,
		Message:  err.Message,
	})
}

// ValidateQuery returns the parser's own error for a script it rejects.
func (engine) ValidateQuery(script string) (bool, bool, error) {
	return ValidateQuery(script)
}

// Explain plans with EXPLAIN, or EXPLAIN (DISTSQL) for JSON, whose diagram link
// is the one machine-readable plan CockroachDB returns.
func (engine) Explain(statement string, format analysis.ExplainFormat) (string, error) {
	return ExplainStatement(statement, format == analysis.ExplainFormatJSON)
}

// ResultLimit wraps in a CTE only what RewriteResultLimit refuses: an empty
// script, one it cannot parse, or several statements.
func (engine) ResultLimit(statement string, limit int, _ string) string {
	rewritten, err := RewriteResultLimit(statement, limit)
	if err != nil {
		return fmt.Sprintf("WITH result AS (\n%s\n) SELECT * FROM result LIMIT %d;", analysis.TrimStatement(statement), limit)
	}
	return rewritten
}
