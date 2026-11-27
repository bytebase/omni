package analysis

import (
	"context"
	"errors"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/doris/ast"
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
	positions := analysis.NewPositions(script)
	syntaxErrors := Diagnose(script)
	diagnostics := make([]analysis.Diagnostic, 0, len(syntaxErrors))
	for _, e := range syntaxErrors {
		diagnostics = append(diagnostics, analysis.Diagnostic{
			Start:    positions.At(e.Offset),
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

var _ analysis.QuerySpanResolver = Engine

// QuerySpans resolves the tables each statement reads, an unqualified one in
// the default database, without the catalog: the span carries no column
// lineage. A syntax error is positioned in the script.
func (engine) QuerySpans(_ context.Context, request analysis.QuerySpanRequest, statements []analysis.QuerySpanStatement) ([]*analysis.ResolvedSpan, error) {
	spans := make([]*analysis.ResolvedSpan, 0, len(statements))
	for _, statement := range statements {
		resolved, err := ResolveQuerySpan(statement.Text, ResolveOptions{
			DefaultDatabase:     request.DefaultDatabase,
			IgnoreCaseSensitive: request.IgnoreCaseSensitive,
		})
		if err != nil {
			err = syntaxError(statement.Text, err)
			if syntaxErr, ok := err.(*analysis.SyntaxError); ok {
				syntaxErr.Position = syntaxErr.Position.Shift(statement.Start)
			}
			return nil, err
		}
		spans = append(spans, &analysis.ResolvedSpan{
			Type:             resolved.Type,
			SourceColumns:    resolved.SourceColumns,
			PredicateColumns: analysis.ColumnSet{},
			Results:          []analysis.QuerySpanResult{},
		})
	}
	return spans, nil
}

var _ analysis.ChangeExtractor = Engine

// ExtractChanges reports what the statements write, an unqualified name in the
// default database. A statement without a node changes nothing, and a resource
// that names no table, such as an INSERT into a table ID, is dropped.
func (engine) ExtractChanges(request analysis.ChangeRequest, statements []analysis.ParsedStatement) (*analysis.Changes, error) {
	changes := &analysis.Changes{}
	for _, statement := range statements {
		if statement.Node == nil {
			continue
		}
		node, ok := statement.Node.(ast.Node)
		if !ok {
			return nil, errors.New("expected astNode for Doris")
		}
		statementChanges := ExtractChangedResources(node, request.DefaultDatabase, statement.Text)
		for _, resource := range statementChanges.Resources {
			if resource.DatabaseOnly || resource.Table != "" {
				changes.Resources = append(changes.Resources, resource)
			}
		}
		for _, text := range statementChanges.DMLStatements {
			changes.DMLStatements = append(changes.DMLStatements, analysis.DMLStatement{Text: text})
		}
		changes.DMLCount += statementChanges.DMLCount
		changes.InsertCount += statementChanges.InsertCount
	}
	return changes, nil
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
