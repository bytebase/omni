package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/clickhouse/parser"
)

// Engine is ClickHouse's implementation of the omni/analysis capability
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

func (engine) Split(script string) []analysis.Statement {
	segments := parser.Split(script)
	if len(segments) == 0 {
		return nil
	}
	statements := make([]analysis.Statement, 0, len(segments))
	for _, seg := range segments {
		statements = append(statements, analysis.Statement{
			Text:  seg.Text,
			Empty: seg.Empty(),
			Start: seg.ByteStart,
			End:   seg.ByteEnd,
		})
	}
	return statements
}

// Diagnose reports the script's first syntax error. An error the parser did
// not locate is marked at the script's start.
func (engine) Diagnose(script string) []analysis.Diagnostic {
	diagnostics := []analysis.Diagnostic{}
	err := Diagnose(script)
	if err == nil {
		return diagnostics
	}
	position := analysis.Position{Line: 1, Column: 1}
	if err.Offset >= 0 {
		position = analysis.PositionAt(script, err.Offset)
	}
	return append(diagnostics, analysis.Diagnostic{
		Start:    position,
		Severity: analysis.SeverityError,
		Message:  err.Message,
	})
}

// ValidateQuery returns the parser's own error for a script it rejects.
func (engine) ValidateQuery(script string) (bool, bool, error) {
	return ValidateQuery(script)
}

// Explain plans in JSON when asked for it and in text otherwise.
func (engine) Explain(statement string, format analysis.ExplainFormat) (string, error) {
	name := ""
	if format == analysis.ExplainFormatJSON {
		name = "JSON"
	}
	return ExplainStatement(statement, name), nil
}

// ResultLimit wraps the statement in a CTE.
func (engine) ResultLimit(statement string, limit int, _ string) string {
	return WrapResultLimit(statement, limit)
}
