package analysis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

// Engine is PostgreSQL's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter            = Engine
	_ analysis.Parser              = Engine
	_ analysis.StatementClassifier = Engine
	_ analysis.StatementRanger     = Engine
	_ analysis.Diagnoser           = Engine
	_ analysis.QueryValidator      = Engine
	_ analysis.Explainer           = Engine
	_ analysis.ResultLimiter       = Engine
)

// Split cuts the script lexically; a segment's text ends at its own
// semicolon.
func (engine) Split(script string) []analysis.Statement {
	segments := pg.Split(script)
	if len(segments) == 0 {
		return nil
	}
	statements := make([]analysis.Statement, 0, len(segments))
	for _, s := range segments {
		statements = append(statements, analysis.Statement{
			Text:  script[s.ByteStart:s.ByteEnd],
			Empty: s.Empty(),
			Start: s.ByteStart,
			End:   s.ByteEnd,
		})
	}
	return statements
}

// Parse parses each non-empty segment and returns one statement per node it
// holds, each carrying its segment. A syntax error is positioned on the
// segment's lines; its column is the segment's own.
func (e engine) Parse(script string) ([]analysis.ParsedStatement, error) {
	var parsed []analysis.ParsedStatement
	for _, segment := range e.Split(script) {
		if segment.Empty {
			continue
		}
		statements, err := pg.Parse(segment.Text)
		if err != nil {
			return nil, segmentSyntaxError(script, segment, err)
		}
		for _, statement := range statements {
			parsed = append(parsed, analysis.ParsedStatement{Statement: segment, Node: statement.AST})
		}
	}
	return parsed, nil
}

func segmentSyntaxError(script string, segment analysis.Statement, err error) error {
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) {
		return err
	}
	position := analysis.PositionAt(segment.Text, parseErr.Position)
	position.Line += analysis.PositionAt(script, segment.Start).Line - 1
	return &analysis.SyntaxError{
		Position:   position,
		Message:    fmt.Sprintf("Syntax error at line %d: %s", position.Line, parseErr.Message),
		RawMessage: parseErr.Message,
	}
}

func (engine) StatementTypes(statement analysis.ParsedStatement) []analysis.StatementType {
	node, ok := statement.Node.(ast.Node)
	if !ok || node == nil {
		return nil
	}
	return StatementTypes(node)
}

func (engine) StatementRanges(script string) []analysis.LSPRange {
	ranges := StatementRanges(script)
	if ranges == nil {
		return nil
	}
	converted := make([]analysis.LSPRange, 0, len(ranges))
	for _, r := range ranges {
		converted = append(converted, analysis.LSPRange{
			Start: analysis.LSPPosition{Line: r.Start.Line, Character: r.Start.Character},
			End:   analysis.LSPPosition{Line: r.End.Line, Character: r.End.Character},
		})
	}
	return converted
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
	if err.ParseError && err.Offset >= 0 {
		position = analysis.PositionAt(err.Text, err.Offset)
	}
	return append(diagnostics, analysis.Diagnostic{
		Start:    position,
		Severity: analysis.SeverityError,
		Message:  err.Message,
	})
}

// ValidateQuery returns the parser's own error for a script it rejects.
func (engine) ValidateQuery(script string) (bool, bool, error) {
	statements, err := pg.Parse(script)
	if err != nil {
		return false, false, err
	}
	nodes := make([]ast.Node, 0, len(statements))
	for _, statement := range statements {
		nodes = append(nodes, statement.AST)
	}
	readOnly, returnsData := ValidateQuery(nodes)
	return readOnly, returnsData, nil
}

func (engine) Explain(statement string, format analysis.ExplainFormat) (string, error) {
	var name string
	switch format {
	case analysis.ExplainFormatJSON:
		name = "json"
	case analysis.ExplainFormatXML:
		name = "xml"
	case analysis.ExplainFormatYAML:
		name = "yaml"
	default:
	}
	return ExplainStatement(statement, name)
}

// ResultLimit rewrites the statement's own LIMIT, or wraps a statement it
// cannot rewrite in a CTE.
func (engine) ResultLimit(statement string, limit int, _ string) string {
	rewritten, err := RewriteResultLimit(statement, limit)
	if err != nil {
		return fmt.Sprintf("WITH result AS (\n%s\n) SELECT * FROM result LIMIT %d;", analysis.TrimStatement(statement), limit)
	}
	return rewritten
}

var _ analysis.QuerySpanResolver = Engine

// QuerySpans resolves each statement against the default database's search
// path, narrowed to the default schema when one is set. A syntax error is
// positioned in its statement's own text.
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
	if request.DefaultSchema != "" {
		searchPath = SearchPathForSelectedSchema(request.DefaultSchema, searchPath)
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

// spanSyntaxError positions a span's syntax error one byte before the parser's
// offset in the statement, or at its start when the parser did not locate it.
func spanSyntaxError(statement string, err *SyntaxError) *analysis.SyntaxError {
	position := analysis.Position{Line: 1, Column: 1}
	var parseErr *parser.ParseError
	if errors.As(err.Err, &parseErr) && parseErr.Position > 0 {
		position = analysis.PositionAt(statement, parseErr.Position-1)
	}
	return &analysis.SyntaxError{Position: position, Message: err.Error()}
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

func (d lookupDatabase) SearchFunctions(searchPath []string, name string) []*metadata.FunctionMetadata {
	_, functions := d.lookup.SearchFunctions(searchPath, name)
	return functions
}

var _ analysis.ChangeExtractor = Engine

// ExtractChanges reports what the statements write, resolving unqualified
// names along the default database's search path, narrowed to the default
// schema when one is set.
func (engine) ExtractChanges(request analysis.ChangeRequest, statements []analysis.ParsedStatement) (*analysis.Changes, error) {
	if request.Database == nil {
		return nil, fmt.Errorf("database metadata for database %q not found (database not synced)", request.DefaultDatabase)
	}
	searchPath := analysis.SearchPath(request.Database.Metadata.GetSearchPath())
	if request.DefaultSchema != "" {
		searchPath = SearchPathForSelectedSchema(request.DefaultSchema, searchPath)
	}
	stmts := make([]ChangeStatement, 0, len(statements))
	for _, statement := range statements {
		node, ok := statement.Node.(ast.Node)
		if !ok || node == nil {
			return nil, errors.New("expected AST for PostgreSQL")
		}
		stmts = append(stmts, ChangeStatement{Node: node, Text: statement.Text})
	}
	changes := ExtractChangedResources(stmts, request.DefaultDatabase, searchPath, changeLookup{request.Database.Lookup(true)})
	return &changes, nil
}

// changeLookup is a synced database's ChangeMetadata, through the caller's
// name matching.
type changeLookup struct{ lookup *analysis.Lookup }

func (m changeLookup) HasSchema(name string) bool { return m.lookup.Schema(name) != nil }

func (m changeLookup) RelationSchema(searchPath []string, name string) string {
	for _, schemaName := range searchPath {
		s := m.lookup.Schema(schemaName)
		if s.Table(name) != nil || s.View(name) != nil || s.MaterializedView(name) != nil || s.ExternalTable(name) != nil {
			return s.Proto().GetName()
		}
	}
	return ""
}

func (m changeLookup) IndexTable(searchPath []string, name string) (string, string, bool) {
	return m.lookup.SearchIndex(searchPath, name)
}

var (
	_ analysis.BackupGenerator  = Engine
	_ analysis.RestoreGenerator = Engine
)

// Backups resolves the script's tables along the default database's search
// path and backs them up into the request's backup schema. A backup starts at
// a 0-based column and ends at a 1-based one.
func (engine) Backups(ctx context.Context, request analysis.BackupRequest, script string) ([]analysis.Backup, error) {
	stmts, err := pg.Parse(script)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transformation: failed to parse statement: %w", err)
	}
	if request.Catalog.GetDatabase == nil {
		return nil, errors.New("failed to prepare transformation: the catalog cannot look databases up")
	}
	database, err := request.Catalog.GetDatabase(ctx, request.DefaultDatabase)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transformation: failed to get database metadata: %w", err)
	}
	if database == nil {
		return nil, fmt.Errorf("failed to prepare transformation: database metadata for database %q not found (database not synced)", request.DefaultDatabase)
	}
	lookup := database.Lookup(true)
	backups, err := TransformDMLToSelect(stmts, script, BackupOptions{
		SearchPath: analysis.SearchPath(database.Metadata.GetSearchPath()),
		SearchRelation: func(searchPath []string, name string) string {
			schema, _ := lookup.SearchRelation(searchPath, name)
			return schema
		},
		TargetSchema: request.BackupDatabase,
		TablePrefix:  request.TablePrefix,
	})
	if err != nil {
		return nil, err
	}
	positions := analysis.NewPositions(script)
	result := make([]analysis.Backup, 0, len(backups))
	for _, backup := range backups {
		start := positions.At(backup.StartOffset)
		start.Column--
		result = append(result, analysis.Backup{
			Statement:    backup.Statement,
			SourceSchema: backup.SourceSchema,
			SourceTable:  backup.SourceTableName,
			BackupTable:  backup.TargetTableName,
			Start:        start,
			End:          positions.At(backup.EndOffset),
		})
	}
	return result, nil
}

// restoreCommentLimit is how many characters of the backed-up statements the
// restore's comment quotes.
const restoreCommentLimit = 1000

// RestoreSQL restores the rows the UPDATE and DELETE statements in the
// backup's range changed, writing them back into the source table, which an
// unqualified name finds in public.
func (engine) RestoreSQL(ctx context.Context, request analysis.RestoreRequest) (string, error) {
	stmts, err := pg.Parse(request.Script)
	if err != nil {
		return "", fmt.Errorf("failed to extract single SQL: failed to parse statement: %v", err)
	}
	positions := analysis.NewPositions(request.Script)
	var dmls []BackedUpDML
	var sqls []string
	for _, dml := range BackedUpDMLs(stmts) {
		start := positions.At(dml.Start)
		start.Column--
		if !analysis.InRange(start, positions.At(dml.End), request.Start, request.End) {
			continue
		}
		dmls = append(dmls, dml)
		sqls = append(sqls, strings.TrimRight(strings.TrimSpace(dml.Text), ";"))
	}
	originalSQL := strings.Join(sqls, ";\n")
	if originalSQL == "" {
		return "", errors.New("no original SQL")
	}
	if len(dmls) == 0 {
		return "", analysis.StatementNotFoundError(request)
	}
	table, err := analysis.LookUpRestoreTable(ctx, request, "public", true, func(database string) error {
		return fmt.Errorf("database metadata not found for %s", database)
	})
	if err != nil {
		return "", err
	}
	database, err := request.Catalog.GetDatabase(ctx, request.SourceDatabase)
	if err != nil || database == nil {
		return "", fmt.Errorf("failed to get database metadata for %s: %v", request.SourceDatabase, err)
	}
	lookup := database.Lookup(true)
	onTable := tableDMLs(stmts, analysis.SearchPath(database.Metadata.GetSearchPath()), func(searchPath []string, name string) string {
		schema, _ := lookup.SearchRelation(searchPath, name)
		return schema
	}, table.Schema, request.SourceTable)
	var nodes []ast.Node
	for _, dml := range dmls {
		if onTable[dml.Node] {
			nodes = append(nodes, dml.Node)
		}
	}
	if len(nodes) == 0 {
		return "", analysis.StatementNotFoundError(request)
	}
	return RestoreSQL(nodes, RestoreOptions{
		OriginalSchema:    table.Schema,
		OriginalTable:     request.SourceTable,
		BackupSchema:      request.BackupSchema,
		BackupTable:       request.BackupTable,
		Table:             table.Table,
		PrimaryKey:        table.PrimaryKey,
		OriginalSQL:       analysis.TruncateComment(originalSQL, restoreCommentLimit),
		PrependStatements: PrependStatements(stmts),
	})
}
