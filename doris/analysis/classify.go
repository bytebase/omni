package analysis

import (
	"github.com/bytebase/omni/doris/parser"
)

// QueryType classifies a SQL statement.
type QueryType int

const (
	QueryTypeUnknown         QueryType = iota
	QueryTypeSelect                    // User SELECT query
	QueryTypeSelectInfoSchema          // SELECT from system tables, or SHOW/DESCRIBE
	QueryTypeDML                       // INSERT, UPDATE, DELETE, MERGE, LOAD, EXPORT, TRUNCATE, COPY
	QueryTypeDDL                       // CREATE, ALTER, DROP, plus GRANT, REVOKE, ADMIN, transaction control, KILL, SET, etc.
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
				panic("doris/analysis: unknown keyword: " + name)
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
		"COPY",
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
		"WARM",
		"PAUSE",
		"RESUME",
	)
}

// tokEOF is the zero value returned by the lexer at end-of-input.
// Matches the unexported parser.tokEOF = 0.
const tokEOF parser.TokenKind = 0

// Classify determines the QueryType for a SQL statement by lexing the first
// meaningful token. It does NOT fully parse the statement.
//
// Whitespace and comments are skipped automatically by the lexer. The
// function returns QueryTypeUnknown for empty input or unrecognised leading
// tokens.
func Classify(statement string) QueryType {
	l := parser.NewLexer(statement)
	tok := l.NextToken()
	if tok.Kind == tokEOF {
		return QueryTypeUnknown
	}
	if qt, ok := firstKeywordType[tok.Kind]; ok {
		return qt
	}
	return QueryTypeUnknown
}
