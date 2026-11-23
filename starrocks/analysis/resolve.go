package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/starrocks/ast"
	"github.com/bytebase/omni/starrocks/parser"
)

// ColumnResource identifies a table a statement reads. StarRocks has no schema
// level, so the resource key stops at the table.
type ColumnResource = analysis.ColumnResource

// SourceColumnSet is the table-level set of resources a statement reads.
type SourceColumnSet = analysis.ColumnSet

// ResolvedQuerySpan is a query span resolved for the session it runs in: the
// tables carry the default database the caller's session supplies, and Type is
// the statement's real operation rather than its leading keyword.
type ResolvedQuerySpan struct {
	// Type is the statement's classification, promoted to
	// QueryTypeSelectInfoSchema when every table it reads is a system table.
	Type QueryType
	// SourceColumns holds one table-level resource per distinct table read.
	SourceColumns SourceColumnSet
	// Results is the outermost query's output columns. With a catalog each
	// names the table columns it reads; without one, each names the column
	// references its expression holds, qualified as written.
	Results []QuerySpanResult
	// NotFoundError is set when the catalog lacks a database, table or
	// column the query names; Results is then empty.
	NotFoundError *NotFoundError
}

// ResolveOptions carries the session facts the caller resolves before the
// analysis. DefaultDatabase fills in a table reference with no database
// qualifier. IgnoreCaseSensitive folds database names when the caller's
// metadata is case-insensitive. Catalog, when set, resolves a query's
// columns and the views it reads.
type ResolveOptions struct {
	DefaultDatabase     string
	IgnoreCaseSensitive bool
	Catalog             Catalog
}

// ErrMixUserSystemTables reports a statement that reads user and system tables
// at once, which the permission check cannot authorize as either.
var ErrMixUserSystemTables = analysis.ErrMixUserSystemTables

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
			Results:       []QuerySpanResult{},
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
	queryType := Classify(single)
	if classified, ok := astQueryType(node); ok {
		queryType = classified
	}
	span, err := getQuerySpan(queryType, &ast.File{Stmts: []ast.Node{node}})
	if err != nil {
		return nil, err
	}

	defaultDatabase := opts.DefaultDatabase
	if submit, ok := node.(*ast.SubmitTaskStmt); ok {
		// A qualified task name selects the database its statement runs in.
		if taskDB, _ := objectName(submit.Name, ""); taskDB != "" {
			defaultDatabase = taskDB
		}
	}
	sources := make(SourceColumnSet, len(span.AccessTables))
	for _, table := range span.AccessTables {
		database := table.Database
		if database == "" {
			database = defaultDatabase
		}
		sources[ColumnResource{Database: database, Table: table.Table}] = true
	}
	results := syntacticResults(span.Results)
	var notFound *NotFoundError
	if opts.Catalog.GetDatabase != nil && isQuery(node) {
		// The catalog resolves each column to the table columns it reads,
		// and each view to the tables behind it, which the plan reads too.
		r := newSpanResolver(opts, single)
		cols := r.query(node, nil, nil)
		if r.err != nil {
			return nil, r.err
		}
		sources, results, notFound = r.tables, results[:0], r.notFound
		if notFound == nil {
			results = resolvedResults(cols)
		}
	}

	allSystems, mixed := isMixedQuery(sources, opts.IgnoreCaseSensitive)
	if mixed {
		return nil, ErrMixUserSystemTables
	}

	queryType = promoteAllSystemQuery(queryType, allSystems)

	return &ResolvedQuerySpan{
		Type:          queryType,
		SourceColumns: sources,
		Results:       results,
		NotFoundError: notFound,
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

// systemDatabases contains StarRocks system databases: the information schema,
// the read-only `sys` metadatabase, and `_statistics_`. StatisticsMetaManager
// sleeps statistic_manager_sleep_time_sec (60s) before creating `_statistics_`,
// so a fresh server can miss it in SHOW DATABASES (starrocks/testdata/pins.tsv),
// while starrocks/driver/sync.go excludes it and CBO queries join `_statistics_`
// with `information_schema`. StarRocks does not have Doris's `mysql` or
// `__internal_schema`. Reference:
// https://docs.starrocks.io/docs/administration/management/system_catalogs/sys/
var systemDatabases = map[string]bool{
	"information_schema": true,
	"sys":                true,
	"_statistics_":       true,
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
