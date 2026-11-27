package analysis

import (
	"errors"
	"fmt"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/redshift"
	"github.com/bytebase/omni/redshift/parser"
)

// Engine is Redshift's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter        = Engine
	_ analysis.StatementRanger = Engine
	_ analysis.Diagnoser       = Engine
	_ analysis.QueryValidator  = Engine
	_ analysis.Explainer       = Engine
	_ analysis.ResultLimiter   = Engine
	_ analysis.ChangeExtractor = Engine
)

// Split cuts the script lexically; a segment's text ends at its own
// semicolon.
func (engine) Split(script string) []analysis.Statement {
	segments := redshift.Split(script)
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

func (engine) StatementRanges(script string) []analysis.LSPRange {
	return analysis.LSPRanges(script, StatementRanges(script))
}

// Diagnose reports the script's first syntax error.
func (engine) Diagnose(script string) []analysis.Diagnostic {
	diagnostics := []analysis.Diagnostic{}
	diag := Diagnose(script)
	if diag == nil {
		return diagnostics
	}
	return append(diagnostics, analysis.Diagnostic{
		Start:    analysis.PositionAt(script, diag.Offset),
		Severity: analysis.SeverityError,
		Message:  diag.Message,
	})
}

func (engine) ValidateQuery(script string) (bool, bool, error) {
	readOnly, returnsData, err := ValidateQuery(script)
	var parseErr *parser.ParseError
	if errors.As(err, &parseErr) {
		position := analysis.PositionAt(script, parseErr.Position)
		err = &analysis.SyntaxError{
			Position:   position,
			Message:    fmt.Sprintf("Syntax error at line %d: %s", position.Line, parseErr.Message),
			RawMessage: parseErr.Message,
		}
	}
	return readOnly, returnsData, err
}

// Explain ignores format: Redshift's EXPLAIN has no machine-readable form.
func (engine) Explain(statement string, _ analysis.ExplainFormat) (string, error) {
	return ExplainStatement(statement)
}

func (engine) ResultLimit(statement string, limit int, _ string) string {
	return StatementWithResultLimit(statement, limit)
}

// ExtractChanges reports what the request's script writes, re-parsing the
// script rather than reading the statements.
func (engine) ExtractChanges(request analysis.ChangeRequest, _ []analysis.ParsedStatement) (*analysis.Changes, error) {
	extracted, err := ExtractChangedResources(request.Script, request.DefaultDatabase, request.DefaultSchema)
	if err != nil {
		return nil, err
	}
	changes := &analysis.Changes{Resources: extracted.Resources, DMLCount: extracted.DMLCount, InsertCount: extracted.InsertCount}
	for _, text := range extracted.DMLStatements {
		changes.DMLStatements = append(changes.DMLStatements, analysis.DMLStatement{Text: text})
	}
	return changes, nil
}
