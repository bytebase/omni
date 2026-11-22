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

// ParsedStmt is a concrete parsed Databricks SQL statement node carrying its
// statement kind, optional target name/text metadata, and source span.
type ParsedStmt struct {
	Kind   StmtKind
	Target string
	Loc    Loc
}

func (s *ParsedStmt) statementNode() {}
func (s *ParsedStmt) Span() Loc      { return s.Loc }
