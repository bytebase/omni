package analysis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql"
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/mssql/parser"
)

// Engine is SQL Server's implementation of the omni/analysis capability
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

// Split returns mssql.Split's segments; a GO separator is an empty segment.
func (engine) Split(script string) []analysis.Statement {
	segments := mssql.Split(script)
	if len(segments) == 0 {
		return nil
	}
	statements := make([]analysis.Statement, 0, len(segments))
	for _, s := range segments {
		statements = append(statements, analysis.Statement{
			Text:  script[s.ByteStart:s.ByteEnd],
			Empty: s.Empty,
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
		statements, err := mssql.Parse(segment.Text)
		if err != nil {
			return nil, segmentSyntaxError(script, segment, err)
		}
		for _, statement := range statements {
			if statement.Empty() {
				continue
			}
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
	statements, err := mssql.Parse(script)
	if err != nil {
		return false, false, err
	}
	valid, allQuery := ValidateQuery(statements)
	return valid, allQuery, nil
}

// ResultLimit wraps a statement it cannot rewrite in a CTE.
func (engine) ResultLimit(statement string, limit int, _ string) string {
	rewritten, err := mssql.StatementWithResultLimit(statement, limit)
	if err != nil {
		return fmt.Sprintf("WITH result AS (%s) SELECT TOP %d * FROM result;", analysis.TrimStatement(statement), limit)
	}
	return rewritten
}

var _ analysis.QuerySpanResolver = Engine

// QuerySpans resolves the statements as one batch: the temp tables and plan
// options a statement declares hold for the ones after it. A span leaves out
// Untraced.
func (engine) QuerySpans(ctx context.Context, request analysis.QuerySpanRequest, statements []analysis.QuerySpanStatement) ([]*analysis.ResolvedSpan, error) {
	opts := ResolveOptions{
		DefaultDatabase:     request.DefaultDatabase,
		DefaultSchema:       request.DefaultSchema,
		IgnoreCaseSensitive: request.IgnoreCaseSensitive,
		Catalog:             spanCatalog(ctx, request.Catalog),
		TempTables:          map[string]*TempTable{},
		PlanOutput:          map[string]bool{},
	}
	spans := make([]*analysis.ResolvedSpan, 0, len(statements))
	for _, statement := range statements {
		span, err := ResolveQuerySpan(statement.Text, opts)
		if err != nil {
			return nil, err
		}
		spans = append(spans, &analysis.ResolvedSpan{
			Type:             span.Type,
			ExplainAnalyze:   span.ExplainAnalyze,
			SourceColumns:    span.SourceColumns,
			Results:          span.Results,
			PredicateColumns: span.PredicateColumns,
			NotFoundError:    span.NotFoundError,
		})
	}
	return spans, nil
}

func spanCatalog(ctx context.Context, catalog analysis.Catalog) Catalog {
	var c Catalog
	if list := catalog.ListDatabaseNames; list != nil {
		c.ListDatabaseNames = func() ([]string, error) { return list(ctx) }
	}
	if get := catalog.GetDatabase; get != nil {
		c.GetDatabase = func(name string) (Database, error) {
			database, err := get(ctx, name)
			if err != nil || database == nil {
				return nil, err
			}
			return lookupDatabase{database.Lookup(false)}, nil
		}
	}
	return c
}

// lookupDatabase is a synced database's Database, through the caller's name
// matching.
type lookupDatabase struct{ lookup *analysis.Lookup }

func (d lookupDatabase) ListSchemaNames() []string { return d.lookup.ListSchemaNames() }

func (d lookupDatabase) GetSchema(name string) Schema {
	schema := d.lookup.Schema(name)
	if schema == nil {
		return nil
	}
	return lookupSchema{schema}
}

type lookupSchema struct{ schema *analysis.SchemaLookup }

func (s lookupSchema) ListTableNames() []string { return s.schema.ListTableNames() }

func (s lookupSchema) GetTable(name string) *metadata.TableMetadata { return s.schema.Table(name) }

func (s lookupSchema) ListViewNames() []string { return s.schema.ListViewNames() }

func (s lookupSchema) GetView(name string) *metadata.ViewMetadata { return s.schema.View(name) }

var _ analysis.ChangeExtractor = Engine

// ExtractChanges reports what the statements write, resolving unqualified
// names in the default database and schema. A statement without a node
// changes nothing.
func (engine) ExtractChanges(request analysis.ChangeRequest, statements []analysis.ParsedStatement) (*analysis.Changes, error) {
	changes := &analysis.Changes{}
	for _, statement := range statements {
		if statement.Node == nil {
			continue
		}
		node, ok := statement.Node.(ast.Node)
		if !ok {
			return nil, errors.New("expected AST for MSSQL")
		}
		appendChanges(changes, ExtractChangedResources(node, request.DefaultDatabase, request.DefaultSchema, statement.Text))
	}
	return changes, nil
}

func appendChanges(changes *analysis.Changes, statement StatementChanges) {
	changes.Resources = append(changes.Resources, statement.Resources...)
	for _, text := range statement.DMLStatements {
		changes.DMLStatements = append(changes.DMLStatements, analysis.DMLStatement{Text: text})
	}
	changes.DMLCount += statement.DMLCount
	changes.InsertCount += statement.InsertCount
}

var (
	_ analysis.BackupGenerator  = Engine
	_ analysis.RestoreGenerator = Engine
)

// Backups backs up the script's UPDATE and DELETE statements, whose
// unqualified tables are in the source database, into the backup database.
// A backup's positions are TransformDMLToSelect's.
func (e engine) Backups(_ context.Context, request analysis.BackupRequest, script string) ([]analysis.Backup, error) {
	stmts, err := e.parseForBackup(script)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transformation: failed to parse statement: %w", err)
	}
	backups, err := TransformDMLToSelect(stmts, request.SourceDatabase, request.BackupDatabase, request.TablePrefix, request.CaseSensitive)
	if err != nil {
		return nil, err
	}
	result := make([]analysis.Backup, 0, len(backups))
	for _, backup := range backups {
		result = append(result, analysis.Backup{
			Statement:    backup.Statement,
			SourceSchema: backup.SourceSchema,
			SourceTable:  backup.SourceTableName,
			BackupTable:  backup.TargetTableName,
			Start:        analysis.Position(backup.StartPosition),
			End:          analysis.Position(backup.EndPosition),
		})
	}
	return result, nil
}

// RestoreSQL restores the rows the UPDATE and DELETE statements in the
// backup's range changed, writing them back into the source table, which an
// unqualified name finds in dbo.
func (e engine) RestoreSQL(ctx context.Context, request analysis.RestoreRequest) (string, error) {
	stmts, err := e.parseForBackup(request.Script)
	if err != nil {
		return "", fmt.Errorf("failed to extract single SQL: failed to parse statement: %v", err)
	}
	var dmls []BackedUpDML
	var sqls []string
	for _, dml := range BackedUpDMLs(stmts) {
		if analysis.InRange(analysis.Position(dml.StartPosition), analysis.Position(dml.EndPosition), request.Start, request.End) {
			dmls = append(dmls, dml)
			sqls = append(sqls, dml.SQL)
		}
	}
	if len(dmls) == 0 {
		return "", errors.New("no original SQL")
	}
	table, err := analysis.LookUpRestoreTable(ctx, request, defaultSchema, false, func(string) error {
		return errors.New("database metadata is nil")
	})
	if err != nil {
		return "", err
	}
	return RestoreSQL(dmls[0].Node, RestoreOptions{
		OriginalDatabase: request.SourceDatabase,
		OriginalSchema:   table.Schema,
		OriginalTable:    request.SourceTable,
		BackupDatabase:   request.BackupDatabase,
		BackupTable:      request.BackupTable,
		Table:            table.Table,
		PrimaryKey:       table.PrimaryKey,
		OriginalSQL:      strings.Join(sqls, ";\n") + ";",
	})
}

// parseForBackup is Parse's statements with the line each starts on in the
// script.
func (e engine) parseForBackup(script string) ([]ParsedStatement, error) {
	parsed, err := e.Parse(script)
	if err != nil {
		return nil, err
	}
	positions := analysis.NewPositions(script)
	stmts := make([]ParsedStatement, 0, len(parsed))
	for _, statement := range parsed {
		node, _ := statement.Node.(ast.Node)
		stmts = append(stmts, ParsedStatement{Text: statement.Text, Node: node, StartLine: positions.At(statement.Start).Line})
	}
	return stmts, nil
}
