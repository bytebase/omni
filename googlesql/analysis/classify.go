// Package analysis provides Bytebase-facing query analysis for GoogleSQL — the
// SQL dialect shared by Google BigQuery and Google Cloud Spanner. It ships the
// concerns the legacy plugin/parser/{bigquery,spanner} packages provide today:
//
//   - statement classification (the legacy queryTypeListener: Select / DML /
//     DDL / Explain / SelectInfoSchema),
//   - fine-grained statement types (the legacy AST-to-storepb switch approval
//     rules and risk levels read),
//   - result-limit rewriting (BigQuery's CTE wrap, Spanner's in-place LIMIT),
//     and
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

	"github.com/bytebase/omni/analysis"
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
type QueryType = analysis.QueryType

const (
	// Unknown is returned for empty input, a parse failure, or an unrecognized
	// statement node (the legacy base.QueryTypeUnknown).
	Unknown = analysis.QueryTypeUnknown
	// Select is a read-only user query (query_statement) that does not read
	// exclusively from a system/information schema.
	Select = analysis.Select
	// Explain is an EXPLAIN statement (explain_statement). Reserved for
	// forward-compatibility: the merged parser foundation does not yet build an
	// EXPLAIN node (its dispatch case is stubbed), so ClassifySQL does not emit
	// Explain today — but Classify maps the node type when parser-utility lands.
	Explain = analysis.Explain
	// SelectInfoSchema is a read-only query that reads EXCLUSIVELY from a
	// system/information schema (the legacy allSystems case of a Select).
	SelectInfoSchema = analysis.SelectInfoSchema
	// DDL changes schema: CREATE/ALTER/DROP over TABLE/VIEW/INDEX/SCHEMA/DATABASE,
	// plus GRANT/REVOKE (the legacy DDL rule list).
	DDL = analysis.DDL
	// DML changes table data: INSERT/UPDATE/DELETE/MERGE/TRUNCATE (the legacy DML
	// rule list; CALL is DML in the legacy listener but is not yet a merged node).
	DML = analysis.DML
)

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
	switch n := node.(type) {
	// Query (query_statement) — Select, possibly promoted to SelectInfoSchema.
	// A query with a pipe suffix is Select unless one of its operators writes
	// (a pipe INSERT is DML, a pipe CREATE TABLE is DDL).
	case *ast.QueryStmt:
		if n.Pipe != nil {
			if t := classifyPipe(n.Pipe); t != Select {
				return t
			}
		}
		if isAllSystemQuery(node, dialect) {
			return SelectInfoSchema
		}
		return Select
	case *ast.SelectStmt, *ast.SetOperation:
		if isAllSystemQuery(node, dialect) {
			return SelectInfoSchema
		}
		return Select

	// A statement that is only a pipe suffix, and a table-returning statement
	// followed by one (sql_statement_body). Their classification is the
	// operators': a read unless an operator writes.
	case *ast.SubpipelineStmt:
		return classifyPipe(n.Pipeline)
	case *ast.PipeStatement:
		// The head and the suffix both decide. A CALL can write from its body,
		// so an unknown head stays Unknown even when the suffix writes: the
		// combined statement is no more classifiable than the head.
		head := Select
		if n.Head != nil {
			head = Classify(n.Head, dialect)
		}
		return mergeQueryTypes(head, classifyPipe(n.Pipeline))

	// DML (dml_statement / merge_statement / truncate_statement).
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt,
		*ast.MergeStmt, *ast.TruncateStmt:
		return DML

	// DDL — CREATE / ALTER / DROP over table-like objects, plus the BigQuery-only
	// object DDL (FUNCTION / TABLE FUNCTION / PROCEDURE, MATERIALIZED|APPROX VIEW,
	// SEARCH|VECTOR INDEX, SNAPSHOT TABLE, generic-entity CAPACITY/RESERVATION/
	// ASSIGNMENT, ROW ACCESS POLICY). The legacy listener groups every CREATE/
	// ALTER/DROP form under DDL, so these classify as DDL too.
	case *ast.CreateTableStmt, *ast.CreateViewStmt, *ast.CreateIndexStmt,
		*ast.CreateSchemaStmt, *ast.CreateDatabaseStmt,
		*ast.AlterStmt, *ast.DropStmt,
		*ast.CreateFunctionStmt, *ast.CreateProcedureStmt,
		*ast.CreateMaterializedViewStmt, *ast.CreateSnapshotStmt,
		*ast.SearchVectorIndexStmt, *ast.CreateRowAccessPolicyStmt,
		*ast.CreateEntityStmt, *ast.BQAlterStmt, *ast.BQDropStmt,
		*ast.DropAllRowAccessPoliciesStmt,
		// Spanner-only object DDL (parser-ddl-spanner node): CHANGE STREAM,
		// SEQUENCE, ROLE, LOCALITY GROUP, PROTO BUNDLE. Same DDL classification as
		// every other CREATE/ALTER/DROP form.
		*ast.CreateChangeStreamStmt, *ast.AlterChangeStreamStmt, *ast.DropChangeStreamStmt,
		*ast.CreateSequenceStmt, *ast.AlterSequenceStmt, *ast.DropSequenceStmt,
		*ast.CreateRoleStmt, *ast.DropRoleStmt,
		*ast.CreateLocalityGroupStmt, *ast.AlterLocalityGroupStmt, *ast.DropLocalityGroupStmt,
		*ast.CreateProtoBundleStmt, *ast.AlterProtoBundleStmt,
		// CREATE PROPERTY GRAPH (parser-gql node) is DDL — the legacy
		// queryTypeListener lists create_property_graph_statement directly among the
		// DDL forms (plugin/parser/{bigquery,spanner}/query_type.go), same as every
		// other CREATE.
		*ast.CreatePropertyGraphStmt:
		return DDL

	// DCL — GRANT / REVOKE are DDL in the legacy listener (its rule list groups
	// privilege statements under DDL alongside CREATE/ALTER/DROP).
	case *ast.GrantStmt, *ast.RevokeStmt:
		return DDL

	// UNDROP restores a dropped schema, which is schema-changing DDL
	// (undrop_statement sits among the DDL alternatives of
	// sql_statement_body_no_pipe_suffix). The legacy listener never saw the form
	// — its AST had no node for it — so there is no parity case to mirror.
	case *ast.UndropStmt:
		return DDL

	// GQL graph query (parser-gql node: `GRAPH <name> <ops>`). The legacy
	// queryTypeListener has NO case for gql_statement, so it falls through its
	// switch to `default: base.QueryTypeUnknown` — a bare GQL graph read is
	// classified Unknown, NOT Select, in BOTH dialects
	// (plugin/parser/{bigquery,spanner}/query_type.go). We mirror that exactly:
	// query-type parity with bytebase means Unknown here (the legacy GetQuerySpan
	// then short-circuits a non-Select to an empty span — see query_span.go's
	// GQLStmt handling, divergence-ledger note). An explicit case (rather than
	// relying on the default) documents this is a deliberate, evidence-backed
	// parity choice, not an oversight.
	case *ast.GQLStmt:
		return Unknown

	default:
		return Unknown
	}
}

// classifyPipe classifies a pipe suffix by what its operators do. It is Select
// unless an operator writes rows (pipe INSERT → DML) or schema (pipe CREATE
// TABLE → DDL). EXPORT DATA writes outside the database and has no
// classification of its own elsewhere in this package, so it stays Unknown.
func classifyPipe(p *ast.Subpipeline) QueryType {
	if p == nil {
		return Select
	}
	t := Select
	for _, op := range p.Ops {
		t = mergeQueryTypes(t, classifyPipeOp(op))
	}
	return t
}

// classifyPipeOp classifies one pipe operator, including the subpipelines a
// branching operator holds.
func classifyPipeOp(op ast.Node) QueryType {
	switch n := op.(type) {
	case *ast.PipeInsert:
		return DML
	case *ast.PipeCreateTable:
		return DDL
	case *ast.PipeExportData:
		return Unknown
	case *ast.PipeLog:
		return classifyPipe(n.Pipeline)
	case *ast.PipeIf:
		t := Select
		for _, c := range n.Cases {
			t = mergeQueryTypes(t, classifyPipe(c.Then))
		}
		return mergeQueryTypes(t, classifyPipe(n.Else))
	case *ast.PipeFork:
		t := Select
		for _, branch := range n.Pipelines {
			t = mergeQueryTypes(t, classifyPipe(branch))
		}
		return t
	case *ast.PipeTee:
		t := Select
		for _, branch := range n.Pipelines {
			t = mergeQueryTypes(t, classifyPipe(branch))
		}
		return t
	}
	return Select
}

// mergeQueryTypes combines two classifications of paths a statement can take.
// Reading one path and writing another still writes (DML or DDL); a DML and a
// DDL path, or a path this package cannot classify, make the whole statement
// Unknown rather than picking one.
func mergeQueryTypes(a, b QueryType) QueryType {
	if a == Unknown || b == Unknown {
		return Unknown
	}
	if a == b {
		return a
	}
	if a == Select {
		return b
	}
	if b == Select {
		return a
	}
	return Unknown
}

// ClassifySQL parses sql and returns the QueryType of its first statement in
// the given dialect. Input that does not parse completely yields Unknown: a
// truncated statement could classify as a read when the engine would execute
// something else, and Unknown is the verdict that grants nothing. Empty,
// whitespace, or comment-only input is Unknown as well.
func ClassifySQL(sql string, dialect Dialect) QueryType {
	file, err := parseFile(sql)
	if err != nil {
		return Unknown
	}
	return ClassifyFromFile(file, dialect)
}

// ClassifyFromFile classifies the first statement of an already-parsed file.
// Returns Unknown for a nil/empty file.
func ClassifyFromFile(file *ast.File, dialect Dialect) QueryType {
	if file == nil || len(file.Stmts) == 0 {
		return Unknown
	}
	return Classify(file.Stmts[0], dialect)
}

// parseFile runs the strict GoogleSQL parser. It is the single parse entry
// point both classification and query-span share, so they never disagree about
// which statement is "first" or whether the input parsed at all.
func parseFile(sql string) (*ast.File, error) {
	return parser.Parse(sql)
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
