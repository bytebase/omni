// Package analysis provides Bytebase-facing query analysis for GoogleSQL — the
// SQL dialect shared by Google BigQuery and Google Cloud Spanner. It ships two
// concerns the legacy plugin/parser/{bigquery,spanner} packages provide today:
//
//   - statement classification (the legacy queryTypeListener: Select / DML /
//     DDL / Explain / SelectInfoSchema), and
//   - query-span extraction (table/column lineage for masking and SQL-editor
//     access tracking — the legacy querySpanExtractor + accessTableListener).
//
// It is built on the hand-written googlesql/parser AST (one parser serves both
// engines) and is parameterized by Dialect for the three genuine bigquery vs
// spanner seams the contract identifies (contract.md §4): the metadata model
// (BigQuery project.dataset.table vs Spanner db.schema.table), the system-schema
// set (BigQuery INFORMATION_SCHEMA vs Spanner INFORMATION_SCHEMA + SPANNER_SYS),
// and the default unqualified-join semantics. Everything else is shared.
//
// Scope (matching the merged omni peers trino/analysis, doris/analysis,
// snowflake/analysis): the extraction is parser-driven and best-effort. It does
// NOT resolve names against live catalog metadata (the legacy
// GetDatabaseMetadataFunc path) nor expand `SELECT *` against a schema — that
// metadata-aware resolution belongs to the bytebase-switch layer, which maps
// these local types onto base.QuerySpan. This package produces the structural
// lineage the parser alone can determine.
package analysis

import (
	"strings"

	"github.com/bytebase/omni/googlesql/ast"
	"github.com/bytebase/omni/googlesql/parser"
)

// Dialect selects the bigquery vs spanner behavior for the three contract
// divergences. One GoogleSQL parser+AST serves both; this enum carries the
// per-dialect deltas (system schemas, name-part bucketing, default join type)
// down into the shared analysis code.
type Dialect int

const (
	// DialectBigQuery is Google BigQuery. Metadata model: project.dataset.table
	// (no schema layer; a dataset is the database). System schema:
	// INFORMATION_SCHEMA only. Default unqualified join: CROSS.
	DialectBigQuery Dialect = iota
	// DialectSpanner is Google Cloud Spanner. Metadata model: db.schema.table
	// (named schemas under one database). System schemas: INFORMATION_SCHEMA and
	// SPANNER_SYS. Default unqualified join: INNER.
	DialectSpanner
)

// String returns a human-readable dialect name.
func (d Dialect) String() string {
	switch d {
	case DialectBigQuery:
		return "BIGQUERY"
	case DialectSpanner:
		return "SPANNER"
	default:
		return "UNKNOWN"
	}
}

// QueryType classifies a GoogleSQL statement. The values mirror bytebase's
// base.QueryType so the bytebase-switch node can map them 1:1 (the legacy
// queryTypeListener emits exactly this set: Select, DDL, DML, Explain,
// SelectInfoSchema, plus QueryTypeUnknown).
type QueryType int

const (
	// Unknown is returned for empty input, a parse failure, or an unrecognized
	// statement node (the legacy base.QueryTypeUnknown).
	Unknown QueryType = iota
	// Select is a read-only user query (query_statement) that does not read
	// exclusively from a system/information schema.
	Select
	// Explain is an EXPLAIN statement (explain_statement). Reserved for
	// forward-compatibility: the merged parser foundation does not yet build an
	// EXPLAIN node (its dispatch case is stubbed), so ClassifySQL does not emit
	// Explain today — but Classify maps the node type when parser-utility lands.
	Explain
	// SelectInfoSchema is a read-only query that reads EXCLUSIVELY from a
	// system/information schema (the legacy allSystems case of a Select).
	SelectInfoSchema
	// DDL changes schema: CREATE/ALTER/DROP over TABLE/VIEW/INDEX/SCHEMA/DATABASE,
	// plus GRANT/REVOKE (the legacy DDL rule list).
	DDL
	// DML changes table data: INSERT/UPDATE/DELETE/MERGE/TRUNCATE (the legacy DML
	// rule list; CALL is DML in the legacy listener but is not yet a merged node).
	DML
)

// String returns a human-readable name for the QueryType. The spellings match
// bytebase's base.QueryType.String() so logs/tests read identically across the
// legacy and omni stacks.
func (q QueryType) String() string {
	switch q {
	case Select:
		return "SELECT"
	case Explain:
		return "EXPLAIN"
	case SelectInfoSchema:
		return "SELECT_INFO_SCHEMA"
	case DDL:
		return "DDL"
	case DML:
		return "DML"
	default:
		return "UNKNOWN"
	}
}

// Classify returns the QueryType for a single parsed statement node, in the
// given dialect.
//
// It type-switches over the concrete googlesql/ast statement nodes the merged
// parser builds (query / DML / DDL / DCL). A query statement is reported as
// Select, then promoted to SelectInfoSchema when the statement reads exclusively
// from a system/information schema — exactly the legacy queryTypeListener's
// `allSystems` rule, evaluated here by walking the parsed AST's table paths
// (the legacy accessTableListener) rather than re-scanning the source text.
//
// A nil or unrecognized node yields Unknown.
func Classify(node ast.Node, dialect Dialect) QueryType {
	if node == nil {
		return Unknown
	}
	switch node.(type) {
	// Query (query_statement) — Select, possibly promoted to SelectInfoSchema.
	case *ast.QueryStmt, *ast.SelectStmt, *ast.SetOperation:
		if isAllSystemQuery(node, dialect) {
			return SelectInfoSchema
		}
		return Select

	// DML (dml_statement / merge_statement / truncate_statement).
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt,
		*ast.MergeStmt, *ast.TruncateStmt:
		return DML

	// DDL — CREATE / ALTER / DROP over table-like objects.
	case *ast.CreateTableStmt, *ast.CreateViewStmt, *ast.CreateIndexStmt,
		*ast.CreateSchemaStmt, *ast.CreateDatabaseStmt,
		*ast.AlterStmt, *ast.DropStmt:
		return DDL

	// DCL — GRANT / REVOKE are DDL in the legacy listener (its rule list groups
	// privilege statements under DDL alongside CREATE/ALTER/DROP).
	case *ast.GrantStmt, *ast.RevokeStmt:
		return DDL

	default:
		return Unknown
	}
}

// ClassifySQL parses the first statement in sql (best-effort) and returns its
// QueryType in the given dialect. Parse errors are tolerated — the first
// successfully-parsed statement node is classified; empty/whitespace/comment-only
// input or a total parse failure yields Unknown.
func ClassifySQL(sql string, dialect Dialect) QueryType {
	return ClassifyFromFile(parseFile(sql), dialect)
}

// ClassifyFromFile classifies the first statement of an already-parsed file.
// Returns Unknown for a nil/empty file.
func ClassifyFromFile(file *ast.File, dialect Dialect) QueryType {
	if file == nil || len(file.Stmts) == 0 {
		return Unknown
	}
	return Classify(file.Stmts[0], dialect)
}

// parseFile runs the best-effort GoogleSQL parser and returns the parsed file
// (which always reflects whatever statements parsed, even on error). It is the
// single parse entry point both classification and query-span share, so they
// never disagree about which statement is "first".
func parseFile(sql string) *ast.File {
	file, _ := parser.Parse(sql)
	return file
}

// isAllSystemQuery reports whether a query node reads EXCLUSIVELY from system
// (INFORMATION_SCHEMA / — Spanner — SPANNER_SYS) tables. It is the AST analogue
// of the legacy isMixedQuery()'s allSystems result, evaluated over the access
// tables the query references. A query with no table references at all (e.g.
// `SELECT 1`) is NOT all-system (it has no user tables, but the legacy code only
// returns allSystems when at least one system table is present:
// `!hasUser && hasSystem`).
func isAllSystemQuery(node ast.Node, dialect Dialect) bool {
	tables := collectAccessTables(node, dialect)
	hasSystem, hasUser := false, false
	for _, ta := range tables {
		if ta.IsSystem {
			hasSystem = true
		} else {
			hasUser = true
		}
	}
	return hasSystem && !hasUser
}

// systemSchemas returns the schema names that mark a system/information-schema
// reference for the given dialect. BigQuery: INFORMATION_SCHEMA only. Spanner:
// INFORMATION_SCHEMA and SPANNER_SYS (contract.md §4; oracle.md).
func systemSchemas(dialect Dialect) []string {
	switch dialect {
	case DialectSpanner:
		return []string{"INFORMATION_SCHEMA", "SPANNER_SYS"}
	default:
		return []string{"INFORMATION_SCHEMA"}
	}
}

// isSystemSchemaName reports whether name equals (case-insensitively) one of the
// dialect's system schema names. Matches the legacy isSystemResource (BigQuery:
// EqualFold INFORMATION_SCHEMA; Spanner: INFORMATION_SCHEMA OR SPANNER_SYS).
func isSystemSchemaName(name string, dialect Dialect) bool {
	for _, s := range systemSchemas(dialect) {
		if strings.EqualFold(name, s) {
			return true
		}
	}
	return false
}
