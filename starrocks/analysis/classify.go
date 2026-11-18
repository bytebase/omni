package analysis

import (
	"github.com/bytebase/omni/starrocks/ast"
	"github.com/bytebase/omni/starrocks/parser"
)

// QueryType classifies a SQL statement.
type QueryType int

const (
	QueryTypeUnknown          QueryType = iota
	QueryTypeSelect                     // User SELECT query
	QueryTypeSelectInfoSchema           // SELECT from system tables, or SHOW/DESCRIBE
	QueryTypeDML                        // INSERT, UPDATE, DELETE, MERGE, LOAD, EXPORT, TRUNCATE, COPY
	QueryTypeDDL                        // CREATE, ALTER, DROP, plus GRANT, REVOKE, ADMIN, transaction control, KILL, SET, etc.
)

// String returns a human-readable name for the QueryType.
func (q QueryType) String() string {
	switch q {
	case QueryTypeSelect:
		return "SELECT"
	case QueryTypeSelectInfoSchema:
		return "SELECT_INFO_SCHEMA"
	case QueryTypeDML:
		return "DML"
	case QueryTypeDDL:
		return "DDL"
	default:
		return "UNKNOWN"
	}
}

// firstKeywordType maps a leading keyword TokenKind to its QueryType.
// Populated in init() using parser.KeywordToken so we don't depend on the
// numeric values of the unexported kw* constants.
var firstKeywordType map[parser.TokenKind]QueryType

func init() {
	firstKeywordType = make(map[parser.TokenKind]QueryType)

	register := func(qt QueryType, names ...string) {
		for _, name := range names {
			kind, ok := parser.KeywordToken(name)
			if !ok {
				panic("starrocks/analysis: unknown keyword: " + name)
			}
			firstKeywordType[kind] = qt
		}
	}

	// User SELECT queries (may be refined by downstream analysis for CTEs).
	register(QueryTypeSelect,
		"SELECT",
		"WITH",
	)

	// Administrative / info-schema reads.
	register(QueryTypeSelectInfoSchema,
		"SHOW",
		"DESCRIBE",
		"DESC",
		"HELP",
		"EXPLAIN",
	)

	// Data manipulation language.
	register(QueryTypeDML,
		"INSERT",
		"UPDATE",
		"DELETE",
		"MERGE",
		"LOAD",
		"EXPORT",
		"TRUNCATE",
	)

	// Data definition language and session / administrative control.
	register(QueryTypeDDL,
		"CREATE",
		"ALTER",
		"DROP",
		"GRANT",
		"REVOKE",
		"ADMIN",
		"BEGIN",
		"COMMIT",
		"ROLLBACK",
		"KILL",
		"SET",
		"UNSET",
		"REFRESH",
		"CANCEL",
		"RECOVER",
		"CLEAN",
		"ANALYZE",
		"SYNC",
		"INSTALL",
		"UNINSTALL",
		"BACKUP",
		"RESTORE",
		"LOCK",
		"UNLOCK",
		"USE",
		"PAUSE",
		"RESUME",
	)
}

// tokEOF is the zero value returned by the lexer at end-of-input.
// Matches the unexported parser.tokEOF = 0.
const tokEOF parser.TokenKind = 0

// ClassifyStatement classifies a statement by inspecting its first top-level
// AST node, falling back to the keyword-based Classify when the text does not
// parse. AST inspection reads the real operation off the tree — an
// `EXPLAIN DROP TABLE` is DDL and a `USE` is Unknown — rather than trusting the
// leading keyword.
func ClassifyStatement(statement string) QueryType {
	file, err := parser.Parse(statement)
	if err != nil || file == nil || len(file.Stmts) == 0 {
		return Classify(statement)
	}
	if queryType, ok := astQueryType(file.Stmts[0]); ok {
		return queryType
	}
	return Classify(statement)
}

// astQueryType maps a top-level AST node to its QueryType. The bool return is
// false only when the node is nil-shaped (an EXPLAIN with no inner query); the
// caller then falls back to Classify.
//
// EXPLAIN recurses into the statement it plans, so `EXPLAIN DROP TABLE t` is
// DDL rather than a read-only downgrade. Anything not otherwise enumerated is
// DDL, a safe upper bound for permission checks: it never labels a write as a
// read.
func astQueryType(node ast.Node) (QueryType, bool) {
	if node == nil {
		return 0, false
	}
	switch n := node.(type) {
	case *ast.SelectStmt, *ast.SetOpStmt:
		return QueryTypeSelect, true
	case *ast.ParenSelect:
		// (SELECT ...) at the top level wraps a query; classify by the
		// wrapped query.
		return astQueryType(n.Sel)
	case *ast.ShowStmt,
		*ast.ShowRoutineLoadStmt, *ast.ShowRoutineLoadTaskStmt,
		*ast.ShowJobStmt, *ast.ShowJobTaskStmt,
		*ast.ShowConstraintsStmt, *ast.ShowAnalyzeStmt, *ast.ShowStatsStmt,
		*ast.DescribeStmt, *ast.DescStorageVolumeStmt, *ast.HelpStmt:
		return QueryTypeSelectInfoSchema, true
	case *ast.ExplainStmt:
		if n.Query == nil {
			return 0, false
		}
		return astQueryType(n.Query)
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt,
		*ast.MergeStmt, *ast.TruncateTableStmt:
		return QueryTypeDML, true
	case *ast.UseStmt:
		// USE stays Unknown so a permission check treats it as a hard deny
		// rather than permitting it under the DDL grant.
		return QueryTypeUnknown, true
	}
	return QueryTypeDDL, true
}

// Classify determines the QueryType for a SQL statement by lexing the first
// meaningful token. It does NOT fully parse the statement.
//
// Whitespace and comments are skipped automatically by the lexer. The
// function returns QueryTypeUnknown for empty input or unrecognised leading
// tokens.
func Classify(statement string) QueryType {
	l := parser.NewLexer(statement)
	tok := l.NextToken()
	// Query-grouping parentheses are transparent for classification:
	// (SELECT ...) and ((SELECT ...)) are SELECT statements
	// (engine-verified accepts).
	for tok.Kind == parser.TokenKind('(') {
		tok = l.NextToken()
	}
	if tok.Kind == tokEOF {
		return QueryTypeUnknown
	}
	if qt, ok := firstKeywordType[tok.Kind]; ok {
		return qt
	}
	return QueryTypeUnknown
}
