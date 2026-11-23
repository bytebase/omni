package parser

import "fmt"

// Loc records the [Start, End) byte offsets of an AST node in the source SQL.
type Loc struct {
	Start int
	End   int
}

// ParseError reports a syntax or lexical error with its byte offset.
type ParseError struct {
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("syntax error at position %d: %s", e.Position, e.Message)
}

// ParseResult holds the parsed statements from a Databricks SQL input.
type ParseResult struct {
	Statements []Statement
}

// Statement is implemented by every top-level Databricks SQL statement node.
type Statement interface {
	statementNode()
	Span() Loc
}

// StmtKind categorizes a parsed statement.
type StmtKind string

const (
	StmtQuery            StmtKind = "QUERY"
	StmtInsert           StmtKind = "INSERT"
	StmtMultiInsert      StmtKind = "MULTI_INSERT"
	StmtUpdate           StmtKind = "UPDATE"
	StmtDelete           StmtKind = "DELETE"
	StmtMerge            StmtKind = "MERGE"
	StmtCreateNamespace  StmtKind = "CREATE_NAMESPACE"
	StmtAlterNamespace   StmtKind = "ALTER_NAMESPACE"
	StmtDropNamespace    StmtKind = "DROP_NAMESPACE"
	StmtCreateTable      StmtKind = "CREATE_TABLE"
	StmtCreateTableLike  StmtKind = "CREATE_TABLE_LIKE"
	StmtReplaceTable     StmtKind = "REPLACE_TABLE"
	StmtAlterTable       StmtKind = "ALTER_TABLE"
	StmtDropTable        StmtKind = "DROP_TABLE"
	StmtCreateView       StmtKind = "CREATE_VIEW"
	StmtAlterView        StmtKind = "ALTER_VIEW"
	StmtDropView         StmtKind = "DROP_VIEW"
	StmtCreateFunction   StmtKind = "CREATE_FUNCTION"
	StmtDropFunction     StmtKind = "DROP_FUNCTION"
	StmtCreateIndex      StmtKind = "CREATE_INDEX"
	StmtDropIndex        StmtKind = "DROP_INDEX"
	StmtDeclare          StmtKind = "DECLARE"
	StmtDropVariable     StmtKind = "DROP_VARIABLE"
	StmtUse              StmtKind = "USE"
	StmtSetCatalog       StmtKind = "SET_CATALOG"
	StmtSet              StmtKind = "SET"
	StmtReset            StmtKind = "RESET"
	StmtShow             StmtKind = "SHOW"
	StmtDescribe         StmtKind = "DESCRIBE"
	StmtExplain          StmtKind = "EXPLAIN"
	StmtAnalyze          StmtKind = "ANALYZE"
	StmtComment          StmtKind = "COMMENT"
	StmtRefresh          StmtKind = "REFRESH"
	StmtCache            StmtKind = "CACHE"
	StmtUncache          StmtKind = "UNCACHE"
	StmtClearCache       StmtKind = "CLEAR_CACHE"
	StmtLoadData         StmtKind = "LOAD_DATA"
	StmtTruncate         StmtKind = "TRUNCATE"
	StmtRepair           StmtKind = "REPAIR"
	StmtManageResource   StmtKind = "MANAGE_RESOURCE"
	StmtCall             StmtKind = "CALL"
	StmtExecuteImmediate StmtKind = "EXECUTE_IMMEDIATE"
	StmtCompound         StmtKind = "COMPOUND"
	StmtHiveNative       StmtKind = "HIVE_NATIVE"
)

// TableRef identifies a table or view referenced by a parsed statement.
// Parts holds the unquoted identifier components in source order
// ([table], [schema, table], or [catalog, schema, table]).
type TableRef struct {
	Parts []string
	Loc   Loc
}

// ParsedStmt is a concrete parsed Databricks SQL statement node carrying its
// statement kind, optional target name/text metadata, and source span.
type ParsedStmt struct {
	Kind   StmtKind
	Target string
	// TargetLoc records the source byte span of Target when applicable (e.g.
	// CREATE VIEW / ALTER VIEW target name).
	TargetLoc Loc
	Loc       Loc

	// Materialized is true for CREATE [OR REPLACE] MATERIALIZED VIEW.
	Materialized bool
	// ExplainMode records the EXPLAIN modifier ("FORMATTED", "EXTENDED",
	// "CODEGEN", "COST") when Kind == StmtExplain.
	ExplainMode string
	// Inner holds the explained statement when Kind == StmtExplain.
	Inner Statement
	// HasTransform is true when a query contains a Hive/Spark script transform
	// (MAP, REDUCE, or SELECT TRANSFORM(...) USING ...).
	HasTransform bool
	// HasEmptyPipe is true when a query ends with an empty `|>` pipe stage.
	HasEmptyPipe bool
	// LimitLoc records the byte span of the outermost query's LIMIT expression
	// (or ALL keyword when LimitAll is true). Zero when no top-level LIMIT is
	// present.
	LimitLoc Loc
	LimitAll bool
	// OffsetLoc records the byte span of the outermost query's OFFSET keyword
	// when present.
	OffsetLoc Loc
	// TableRefs holds table and view references encountered in the statement's
	// query/view body (excluding CTE names in scope and dynamic IDENTIFIER(...)).
	TableRefs []TableRef
}

func (s *ParsedStmt) statementNode() {}
func (s *ParsedStmt) Span() Loc      { return s.Loc }
