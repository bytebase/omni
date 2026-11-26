package analysis

import (
	"errors"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/doris/parser"
)

// Engine is Doris's implementation of the omni/analysis capability interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter        = Engine
	_ analysis.StatementRanger = Engine
	_ analysis.Diagnoser       = Engine
	_ analysis.QueryValidator  = Engine
	_ analysis.Explainer       = Engine
	_ analysis.ResultLimiter   = Engine
)

func (engine) Split(script string) []analysis.Statement {
	return SplitStatements(script)
}

func (engine) StatementRanges(script string) []analysis.LSPRange {
	ranges := StatementRanges(script)
	if len(ranges) == 0 {
		return nil
	}
	return analysis.LSPRanges(script, ranges)
}

func (engine) Diagnose(script string) []analysis.Diagnostic {
	syntaxErrors := Diagnose(script)
	diagnostics := make([]analysis.Diagnostic, 0, len(syntaxErrors))
	for _, e := range syntaxErrors {
		diagnostics = append(diagnostics, analysis.Diagnostic{
			Start:    analysis.PositionAt(script, e.Offset),
			Severity: analysis.SeverityError,
			Message:  e.Message,
		})
	}
	return diagnostics
}

func (engine) ValidateQuery(script string) (bool, bool, error) {
	readOnly, returnsData, err := ValidateQuery(script)
	if err != nil {
		return false, false, syntaxError(script, err)
	}
	return readOnly, returnsData, nil
}

// Explain plans in Doris's one plan format whatever the caller asks for.
func (engine) Explain(statement string, _ analysis.ExplainFormat) (string, error) {
	planned, err := ExplainStatement(statement)
	if err != nil {
		return "", syntaxError(statement, err)
	}
	return planned, nil
}

func (engine) ResultLimit(statement string, limit int, _ string) string {
	return RewriteResultLimit(statement, limit)
}

// syntaxError positions the first parse error of text; any other error is
// returned as is.
func syntaxError(text string, err error) error {
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) {
		return err
	}
	return &analysis.SyntaxError{
		Position:   analysis.PositionAt(text, parseErr.Position),
		Message:    parseErr.Error(),
		RawMessage: parseErr.Message,
	}
}
