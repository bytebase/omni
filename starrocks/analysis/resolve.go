package analysis

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
	"github.com/bytebase/omni/starrocks/parser"
)

// ColumnResource identifies a table a statement reads. StarRocks has no schema
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
	stmts := SplitStatements(statement)
	nonEmpty := 0
	var single string
	for _, s := range stmts {
		if s.Empty {
			continue
		}
		nonEmpty++
		single = s.Text
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

	// GetQuerySpan walks only a top-level SELECT/set operation, so EXPLAIN
	// would report no access tables unless the planned statement is unwrapped.
	span, err := GetQuerySpan(unwrapExplainForSpan(single))
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

	return &ResolvedQuerySpan{
		Type:          classifiedQueryType(single, allSystems),
		SourceColumns: sources,
		Results:       span.Results,
	}, nil
}

// unwrapExplainForSpan returns the text of the statement an EXPLAIN plans, or
// the original text when it does not parse to a top-level EXPLAIN.
func unwrapExplainForSpan(statement string) string {
	file, err := parser.Parse(statement)
	if err != nil || file == nil || len(file.Stmts) == 0 {
		return statement
	}
	explain, ok := file.Stmts[0].(*ast.ExplainStmt)
	if !ok || explain.Query == nil {
		return statement
	}
	inner := ast.NodeLoc(explain.Query)
	if inner.Start < 0 || inner.End > len(statement) || inner.Start >= inner.End {
		return statement
	}
	return statement[inner.Start:inner.End]
}

// classifiedQueryType classifies statement and promotes a user SELECT to
// SelectInfoSchema when every table it reads is a system table.
func classifiedQueryType(statement string, allSystems bool) QueryType {
	queryType := ClassifyStatement(statement)
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

// systemDatabases contains StarRocks system databases. StarRocks exposes a
// read-only `sys` metadatabase and does not have Doris's `mysql` or
// `__internal_schema`.
// Reference: https://docs.starrocks.io/docs/administration/management/system_catalogs/sys/
var systemDatabases = map[string]bool{
	"information_schema": true,
	"_statistics_":       true,
	"sys":                true,
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
