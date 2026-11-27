package analysis

import (
	"context"
	"errors"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/spanner/diagnostics"
	"github.com/bytebase/omni/spanner/parser"
)

// Engine is Spanner's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter       = Engine
	_ analysis.Parser         = Engine
	_ analysis.Diagnoser      = Engine
	_ analysis.QueryValidator = Engine
	_ analysis.ResultLimiter  = Engine
)

// Split keeps empty segments. A segment starts where the one before it ends, so
// the whitespace and comments before a statement belong to it, and its text
// ends at the semicolon that follows it.
func (engine) Split(script string) []analysis.Statement {
	segments := parser.Split(script)
	if len(segments) == 0 {
		return nil
	}
	statements := make([]analysis.Statement, 0, len(segments))
	start := 0
	for _, s := range segments {
		end := s.ByteEnd
		if end < len(script) && script[end] == ';' {
			end++
		}
		statements = append(statements, analysis.Statement{
			Text:  script[start:end],
			Empty: s.Empty(),
			Start: start,
			End:   end,
		})
		start = end
	}
	return statements
}

// Parse returns one statement per node of each non-empty segment. It never
// fails: a segment the parser rejects is one statement with a nil Node.
func (e engine) Parse(script string) ([]analysis.ParsedStatement, error) {
	var parsed []analysis.ParsedStatement
	for _, segment := range e.Split(script) {
		if segment.Empty {
			continue
		}
		file, err := parser.Parse(segment.Text)
		if err != nil || file == nil || len(file.Stmts) == 0 {
			parsed = append(parsed, analysis.ParsedStatement{Statement: segment})
			continue
		}
		for _, node := range file.Stmts {
			parsed = append(parsed, analysis.ParsedStatement{Statement: segment, Node: node})
		}
	}
	return parsed, nil
}

// Diagnose reports every syntax error except the parser's "statement parsing is
// not yet supported" stubs, which mark valid statements it does not model.
func (engine) Diagnose(script string) []analysis.Diagnostic {
	positions := analysis.NewPositions(script)
	result := []analysis.Diagnostic{}
	for _, d := range diagnostics.Analyze(script) {
		if strings.HasSuffix(d.Message, "statement parsing is not yet supported") {
			continue
		}
		result = append(result, analysis.Diagnostic{
			Start:    positions.At(d.Range.Start.Offset),
			Severity: analysis.SeverityError,
			Message:  d.Message,
		})
	}
	return result
}

// ValidateQuery returns the parser's own error for a script it rejects.
func (engine) ValidateQuery(script string) (bool, bool, error) {
	return ValidateQuery(script)
}

func (engine) ResultLimit(statement string, limit int, _ string) string {
	return RewriteResultLimit(statement, limit)
}

var _ analysis.QuerySpanResolver = Engine

// QuerySpans resolves each statement in the default database; the default
// schema and IgnoreCaseSensitive are not read. A statement's predicate columns
// are not reported. A syntax error is the parser's own.
func (engine) QuerySpans(ctx context.Context, request analysis.QuerySpanRequest, statements []analysis.QuerySpanStatement) ([]*analysis.ResolvedSpan, error) {
	load := request.Catalog.MetadataLoader()
	options := ResolveOptions{
		DefaultDatabase: request.DefaultDatabase,
		// Named schemas under the one database (a written database is
		// ignored), an empty span for a query reading only system tables, SET
		// classified as a select, and a JOIN USING key named in its metadata case.
		IgnoreWrittenDatabase:  true,
		SchemaScopedMetadata:   true,
		SystemOnlyEmptySpan:    true,
		SetStatementIsSelect:   true,
		CanonicalBaseFieldName: true,
		GetDatabaseMetadata: func(database string) (*metadata.DatabaseSchemaMetadata, error) {
			if load == nil {
				return nil, &ResourceNotFoundError{Database: &database}
			}
			return load(ctx, database)
		},
	}
	spans := make([]*analysis.ResolvedSpan, 0, len(statements))
	for _, statement := range statements {
		span, err := ResolveQuerySpan(statement.Text, options)
		if err != nil {
			return nil, err
		}
		resolved := &analysis.ResolvedSpan{Type: span.Type, SourceColumns: span.SourceColumns, Results: span.Results}
		errors.As(span.NotFoundError, &resolved.NotFoundError)
		spans = append(spans, resolved)
	}
	return spans, nil
}
