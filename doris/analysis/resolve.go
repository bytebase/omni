package analysis

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/doris/ast"
	"github.com/bytebase/omni/doris/parser"
)

// ColumnResource identifies a table a statement reads. Doris has no schema
// level, so the resource key stops at the table.
type ColumnResource struct {
	Database string
	Table    string
}

// SourceColumnSet is the table-level set of resources a statement reads.
type SourceColumnSet map[ColumnResource]bool

// ResolvedQuerySpan is a query span resolved for the session it runs in: the
// tables carry the default database the caller's session supplies, and Type is
// the statement's real operation rather than its leading keyword.
type ResolvedQuerySpan struct {
	// Type is the statement's classification, promoted to
	// QueryTypeSelectInfoSchema when every table it reads is a system table.
	Type QueryType
	// SourceColumns holds one table-level resource per distinct table read.
	SourceColumns SourceColumnSet
	// Results is the outermost SELECT's output columns.
	Results []ColumnInfo
}

// ResolveOptions carries the session facts the caller resolves before the
// analysis. DefaultDatabase fills in a table reference with no database
// qualifier. IgnoreCaseSensitive folds database names when the caller's
// metadata is case-insensitive.
type ResolveOptions struct {
	DefaultDatabase     string
	IgnoreCaseSensitive bool
}

// ErrMixUserSystemTables reports a statement that reads user and system tables
// at once, which the permission check cannot authorize as either.
var ErrMixUserSystemTables = errors.New("cannot access user and system tables at the same time")

// ResolveQuerySpan analyzes a single statement for the session described by
// opts and returns the tables it reads and its classification.
//
// A text with no statement yields an empty span; more than one statement is an
// error, matching the one-statement contract every engine's query span keeps.
// An EXPLAIN is unwrapped to the statement it plans, so table-level access
// checks see the tables the plan reads. The span fails closed on a parse error
// exactly as GetQuerySpan does.
func ResolveQuerySpan(statement string, opts ResolveOptions) (*ResolvedQuerySpan, error) {
	// Validate the whole text before splitting: a buffer that holds only a
	// broken comment ("/*") splits into no statement, and an unterminated
	// comment after a real statement is not part of any segment, so both would
	// otherwise answer a span for text the engine rejects.
	if err := firstLexError(statement); err != nil {
		return nil, err
	}

	stmts := SplitStatements(statement)
	nonEmpty := 0
	var single string
	var singleStart int
	for _, s := range stmts {
		if s.Empty {
			continue
		}
		nonEmpty++
		single = s.Text
		singleStart = s.Start
	}
	if nonEmpty == 0 {
		return &ResolvedQuerySpan{
			SourceColumns: SourceColumnSet{},
			Results:       []ColumnInfo{},
		}, nil
	}
	if nonEmpty > 1 {
		return nil, fmt.Errorf("expecting only one statement to get query span, but got %d", nonEmpty)
	}

	// Parse once: EXPLAIN resolution, the span walk and the classification
	// all read this tree.
	file, err := parser.Parse(single)
	if err != nil {
		// The segment was cut out of the script, so its parse errors are
		// relative to the segment; move them back into script coordinates.
		return nil, offsetParseError(parser.FirstError(err), singleStart)
	}
	var node ast.Node
	if file != nil && len(file.Stmts) > 0 {
		node = file.Stmts[0]
	}
	// GetQuerySpan walks only a top-level SELECT/set operation, so EXPLAIN
	// would report no access tables unless the planned statement is walked.
	if explain, ok := node.(*ast.ExplainStmt); ok && explain.Query != nil {
		node = explain.Query
	}
	queryType := Classify(single)
	if classified, ok := astQueryType(node); ok {
		queryType = classified
	}
	span, err := getQuerySpan(queryType, &ast.File{Stmts: []ast.Node{node}})
	if err != nil {
		return nil, err
	}

	sources := make(SourceColumnSet, len(span.AccessTables))
	for _, table := range span.AccessTables {
		database := table.Database
		if database == "" {
			database = opts.DefaultDatabase
		}
		sources[ColumnResource{Database: database, Table: table.Table}] = true
	}

	allSystems, mixed := isMixedQuery(sources, opts.IgnoreCaseSensitive)
	if mixed {
		return nil, ErrMixUserSystemTables
	}

	queryType = promoteAllSystemQuery(queryType, allSystems)

	return &ResolvedQuerySpan{
		Type:          queryType,
		SourceColumns: sources,
		Results:       span.Results,
	}, nil
}

// promoteAllSystemQuery promotes a user SELECT to SelectInfoSchema when every
// table it reads is a system table.
func promoteAllSystemQuery(queryType QueryType, allSystems bool) QueryType {
	if queryType == QueryTypeSelect && allSystems {
		return QueryTypeSelectInfoSchema
	}
	return queryType
}

// isMixedQuery reports whether the resources mix user and system tables, and
// whether every resource is a system table.
func isMixedQuery(resources SourceColumnSet, ignoreCaseSensitive bool) (bool, bool) {
	hasSystem, hasUser := false, false
	for resource := range resources {
		if isSystemResource(resource, ignoreCaseSensitive) {
			hasSystem = true
		} else {
			hasUser = true
		}
	}
	if hasSystem && hasUser {
		return false, true
	}
	return !hasUser && hasSystem, false
}

// systemDatabases contains Doris system databases: the information schema,
// the MySQL-protocol database, and the internal schema that holds audit and
// statistics tables.
//
// The pinned 3.1.4 server answers SHOW DATABASES with exactly these three
// (doris/testdata/pins.tsv); `_statistics_` is not one of them, even after
// ANALYZE TABLE. Reference:
// https://doris.apache.org/docs/3.x/admin-manual/system-tables/overview
var systemDatabases = map[string]bool{
	"information_schema": true,
	"mysql":              true,
	"__internal_schema":  true,
}

// IsSystemDatabase reports whether name is one of the engine's system
// databases. A case-insensitive caller folds the name first.
func IsSystemDatabase(name string, ignoreCaseSensitive bool) bool {
	if ignoreCaseSensitive {
		name = strings.ToLower(name)
	}
	return systemDatabases[name]
}

func isSystemResource(resource ColumnResource, ignoreCaseSensitive bool) bool {
	return IsSystemDatabase(resource.Database, ignoreCaseSensitive)
}
