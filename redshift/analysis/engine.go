package analysis

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/redshift"
	"github.com/bytebase/omni/redshift/parser"
)

// Engine is Redshift's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter          = Engine
	_ analysis.StatementRanger   = Engine
	_ analysis.Diagnoser         = Engine
	_ analysis.QueryValidator    = Engine
	_ analysis.Explainer         = Engine
	_ analysis.ResultLimiter     = Engine
	_ analysis.ChangeExtractor   = Engine
	_ analysis.QuerySpanResolver = Engine
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

// QuerySpans resolves each statement against the default database's search
// path, or the default schema alone when one is set. A syntax error is
// positioned in the script.
func (engine) QuerySpans(ctx context.Context, request analysis.QuerySpanRequest, statements []analysis.QuerySpanStatement) ([]*analysis.ResolvedSpan, error) {
	if request.Catalog.GetDatabase == nil {
		return nil, errors.New("the catalog cannot look databases up")
	}
	spans := make([]*analysis.ResolvedSpan, 0, len(statements))
	for _, statement := range statements {
		span, err := querySpan(ctx, request, statement)
		if err != nil {
			return nil, err
		}
		spans = append(spans, span)
	}
	return spans, nil
}

func querySpan(ctx context.Context, request analysis.QuerySpanRequest, statement analysis.QuerySpanStatement) (*analysis.ResolvedSpan, error) {
	database, err := request.Catalog.GetDatabase(ctx, request.DefaultDatabase)
	if err != nil {
		return nil, fmt.Errorf("failed to get database metadata for database %q: %w", request.DefaultDatabase, err)
	}
	if database == nil {
		return nil, fmt.Errorf("database metadata for database %q not found (database not synced)", request.DefaultDatabase)
	}
	searchPath := analysis.SearchPath(database.Metadata.GetSearchPath())
	if request.DefaultSchema != "" {
		searchPath = []string{request.DefaultSchema}
	}
	resolved, err := ResolveQuerySpan(statement.Text, ResolveOptions{
		DefaultDatabase: request.DefaultDatabase,
		SearchPath:      searchPath,
		GetDatabaseMetadata: func(name string) (*metadata.DatabaseSchemaMetadata, error) {
			database, err := request.Catalog.GetDatabase(ctx, name)
			if err != nil || database == nil {
				return nil, err
			}
			return database.Metadata, nil
		},
	})
	if errors.Is(err, analysis.ErrMixUserSystemTables) {
		return nil, err
	}
	if err != nil {
		var parseErr *parser.ParseError
		if errors.As(err, &parseErr) {
			position := analysis.PositionAt(statement.Text, parseErr.Position).Shift(statement.Start)
			err = &analysis.SyntaxError{
				Position:   position,
				Message:    fmt.Sprintf("Syntax error at line %d: %s", position.Line, parseErr.Message),
				RawMessage: parseErr.Message,
			}
		}
		return nil, fmt.Errorf("failed to get query span from statement: %s: %w", statement.Text, err)
	}
	span := &analysis.ResolvedSpan{
		Type:             resolved.Type,
		SourceColumns:    resolved.SourceColumns,
		PredicateColumns: resolved.PredicateColumns,
		Results:          resolved.Results,
	}
	errors.As(resolved.FunctionNotSupportedError, &span.FunctionNotSupportedError)
	// A catalog error the resolver reports as not found, such as a missing
	// column, is not a typed not-found error.
	if !errors.As(resolved.NotFoundError, &span.NotFoundError) {
		span.LookupError = resolved.NotFoundError
	}
	return span, nil
}
