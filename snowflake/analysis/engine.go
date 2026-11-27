package analysis

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/snowflake/diagnostics"
	"github.com/bytebase/omni/snowflake/parser"
)

// Engine is Snowflake's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter          = Engine
	_ analysis.Parser            = Engine
	_ analysis.StatementRanger   = Engine
	_ analysis.Diagnoser         = Engine
	_ analysis.QueryValidator    = Engine
	_ analysis.Explainer         = Engine
	_ analysis.ResultLimiter     = Engine
	_ analysis.QuerySpanResolver = Engine
)

// Split returns SplitStatements: a segment's text ends at the semicolon that
// follows it.
func (engine) Split(script string) []analysis.Statement {
	return SplitStatements(script)
}

// Parse parses each non-empty segment Split returns into the segment's first
// node. A syntax error is positioned in the script, and its message names that
// position.
func (e engine) Parse(script string) ([]analysis.ParsedStatement, error) {
	parsed := []analysis.ParsedStatement{}
	for _, segment := range e.Split(script) {
		if segment.Empty {
			continue
		}
		file, err := parser.Parse(segment.Text)
		if err != nil {
			return nil, segmentSyntaxError(script, segment, err)
		}
		if file == nil || len(file.Stmts) == 0 {
			return nil, errors.New("statement yielded no parse tree")
		}
		parsed = append(parsed, analysis.ParsedStatement{Statement: segment, Node: file.Stmts[0]})
	}
	return parsed, nil
}

func segmentSyntaxError(script string, segment analysis.Statement, err error) error {
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) {
		return err
	}
	position := analysis.PositionAt(segment.Text, parseErr.Position).Shift(analysis.PositionAt(script, segment.Start))
	return &analysis.SyntaxError{
		Position:   position,
		Message:    fmt.Sprintf("Syntax error at line %d:%d \n%s", position.Line, position.Column, parseErr.Message),
		RawMessage: parseErr.Message,
	}
}

func (engine) StatementRanges(script string) []analysis.LSPRange {
	return analysis.LSPRanges(script, StatementRanges(script))
}

func (engine) Diagnose(script string) []analysis.Diagnostic {
	positions := analysis.NewPositions(script)
	found := diagnostics.Analyze(script)
	result := make([]analysis.Diagnostic, 0, len(found))
	for _, d := range found {
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

// Explain writes EXPLAIN USING JSON for the JSON format and the tabular EXPLAIN
// for any other. Its errors are the parser's own.
func (engine) Explain(statement string, format analysis.ExplainFormat) (string, error) {
	return ExplainStatement(statement, format == analysis.ExplainFormatJSON)
}

// ResultLimit rewrites the statement's own LIMIT, or wraps a statement it
// cannot rewrite in a subquery.
func (engine) ResultLimit(statement string, limit int, _ string) string {
	rewritten, err := RewriteResultLimit(statement, limit)
	if err != nil {
		return fmt.Sprintf("SELECT * FROM (%s) LIMIT %d", analysis.TrimStatement(statement), limit)
	}
	return rewritten
}

// QuerySpans resolves each statement in the default database and schema. A
// table missing from the catalog sets the span's NotFoundError, with the
// tables the statement accesses as its SourceColumns.
func (engine) QuerySpans(ctx context.Context, request analysis.QuerySpanRequest, statements []analysis.QuerySpanStatement) ([]*analysis.ResolvedSpan, error) {
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

func querySpan(ctx context.Context, request analysis.QuerySpanRequest, statement string) (*ResolvedSpan, error) {
	span, err := ResolveQuerySpan(statement, ResolveOptions{
		DefaultDatabase:     request.DefaultDatabase,
		DefaultSchema:       request.DefaultSchema,
		IgnoreCaseSensitive: request.IgnoreCaseSensitive,
		ResolveTable: func(database, schema, table string) (*Table, error) {
			return lookupTable(ctx, request.Catalog, database, schema, table)
		},
	})
	if err == nil || errors.Is(err, ErrMixUserSystemTables) {
		return span, err
	}
	var resolveErr *ResolveError
	if !errors.As(err, &resolveErr) {
		return nil, err
	}
	var notFound *analysis.ResourceNotFoundError
	if errors.As(resolveErr.Err, &notFound) {
		return &ResolvedSpan{SourceColumns: resolveErr.SourceColumns, Results: []QuerySpanResult{}, NotFoundError: notFound}, nil
	}
	return nil, resolveErr.Err
}

// lookupTable finds a table by its folded names, which must match the
// catalog's spelling exactly. An empty database or schema name matches any
// listed one but looks up the empty name itself.
func lookupTable(ctx context.Context, catalog analysis.Catalog, databaseName, schemaName, tableName string) (*Table, error) {
	if catalog.ListDatabaseNames == nil || catalog.GetDatabase == nil {
		return nil, errors.New("the catalog cannot look tables up")
	}
	databaseNames, err := catalog.ListDatabaseNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all databases: %w", err)
	}
	notFound := fmt.Errorf("table %s.%s.%s is not found", databaseName, schemaName, tableName)
	if !listed(databaseNames, databaseName) {
		return nil, notFound
	}
	database, err := catalog.GetDatabase(ctx, databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to get database %s: %w", databaseName, err)
	}
	if database == nil {
		return nil, fmt.Errorf("database metadata for database %q not found (database not synced)", databaseName)
	}
	lookup := database.Lookup(true)
	if !listed(lookup.ListSchemaNames(), schemaName) {
		return nil, notFound
	}
	schema := lookup.Schema(schemaName)
	if schema == nil {
		return nil, fmt.Errorf("schema %s.%s is not found", databaseName, schemaName)
	}
	if !slices.Contains(schema.ListTableNames(), tableName) {
		return nil, notFound
	}
	table := schema.Table(tableName)
	var columns []string
	for _, column := range table.GetColumns() {
		columns = append(columns, column.GetName())
	}
	return &Table{Name: table.GetName(), Columns: columns}, nil
}

// listed reports whether names holds name, or any name when name is empty.
func listed(names []string, name string) bool {
	return name == "" && len(names) > 0 || slices.Contains(names, name)
}
