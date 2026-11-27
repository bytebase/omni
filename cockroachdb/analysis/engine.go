package analysis

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

// Engine is CockroachDB's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter          = Engine
	_ analysis.Diagnoser         = Engine
	_ analysis.QueryValidator    = Engine
	_ analysis.Explainer         = Engine
	_ analysis.ResultLimiter     = Engine
	_ analysis.QuerySpanResolver = Engine
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

// QuerySpans resolves each statement against the default database's search
// path; a selected schema other than public goes before public. A syntax error
// is positioned in its statement's own text.
func (engine) QuerySpans(ctx context.Context, request analysis.QuerySpanRequest, statements []analysis.QuerySpanStatement) ([]*analysis.ResolvedSpan, error) {
	if request.Catalog.GetDatabase == nil {
		return nil, errors.New("the catalog cannot look databases up")
	}
	spans := make([]*analysis.ResolvedSpan, 0, len(statements))
	for _, statement := range statements {
		span, err := querySpan(ctx, request, statement.Text)
		if err != nil {
			return nil, err
		}
		spans = append(spans, span)
	}
	return spans, nil
}

func querySpan(ctx context.Context, request analysis.QuerySpanRequest, statement string) (*analysis.ResolvedSpan, error) {
	database, err := request.Catalog.GetDatabase(ctx, request.DefaultDatabase)
	if err != nil {
		return nil, fmt.Errorf("failed to get database metadata for database %q: %w", request.DefaultDatabase, err)
	}
	if database == nil {
		return nil, fmt.Errorf("database metadata for database %q not found (database not synced)", request.DefaultDatabase)
	}
	searchPath := analysis.SearchPath(database.Metadata.GetSearchPath())
	switch request.DefaultSchema {
	case "":
	case "public":
		searchPath = []string{"public"}
	default:
		searchPath = []string{request.DefaultSchema, "public"}
	}
	span, err := ResolveQuerySpan(ctx, statement, ResolveOptions{
		DefaultDatabase: request.DefaultDatabase,
		SearchPath:      searchPath,
		Catalog: Catalog{GetDatabase: func(name string) (Database, error) {
			database, err := request.Catalog.GetDatabase(ctx, name)
			if err != nil || database == nil {
				return nil, err
			}
			return lookupDatabase{database.Lookup(true)}, nil
		}},
	})
	if err != nil {
		var syntaxErr *SyntaxError
		if errors.As(err, &syntaxErr) {
			err = spanSyntaxError(statement, syntaxErr)
		}
		return nil, fmt.Errorf("failed to get query span from statement: %s: %w", statement, err)
	}
	return span, nil
}

// spanSyntaxError positions a span's syntax error at the parser's offset in the
// statement, or at its start when the parser did not locate it.
func spanSyntaxError(statement string, err *SyntaxError) *analysis.SyntaxError {
	if err.Position < 0 {
		return &analysis.SyntaxError{Position: analysis.Position{Line: 1, Column: 1}, Message: err.Error()}
	}
	return &analysis.SyntaxError{Position: analysis.PositionAt(statement, err.Position), Message: err.Error(), RawMessage: err.Error()}
}

// lookupDatabase is a synced database's Database, through the caller's name
// matching.
type lookupDatabase struct{ lookup *analysis.Lookup }

func (d lookupDatabase) Metadata() *metadata.DatabaseSchemaMetadata { return d.lookup.Metadata() }

func (d lookupDatabase) GetTable(schema, table string) *metadata.TableMetadata {
	return d.lookup.Schema(schema).Table(table)
}

func (d lookupDatabase) SearchRelation(searchPath []string, name string) (string, string) {
	return d.lookup.SearchRelation(searchPath, name)
}
