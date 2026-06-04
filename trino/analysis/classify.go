// Package analysis provides Bytebase-facing query analysis for Trino SQL:
// statement classification (the query_type / validateQuery contract) and query
// span extraction (table/column lineage for masking and access tracking).
//
// It is the omni counterpart of bytebase's legacy plugin/parser/trino
// getQueryType + GetQuerySpan, built on the hand-written trino/parser AST.
package analysis

import (
	"strings"

	"github.com/bytebase/omni/trino/parser"
)

// QueryType classifies a Trino SQL statement. The values mirror bytebase's
// base.QueryType so the import switch (bytebase-switch node) can map them 1:1.
type QueryType int

const (
	// Unknown is returned for empty input or an unrecognised leading token.
	Unknown QueryType = iota
	// Select is a read-only user SELECT query (SELECT/WITH/TABLE/VALUES) that
	// does not touch a system/information schema.
	Select
	// Explain is an EXPLAIN / EXPLAIN ANALYZE statement.
	Explain
	// SelectInfoSchema is a read-only query against system/information_schema
	// objects, or a SHOW / DESCRIBE statement.
	SelectInfoSchema
	// DDL changes schema: CREATE/ALTER/DROP/RENAME, GRANT/REVOKE/DENY, roles,
	// COMMENT, ANALYZE, REFRESH MATERIALIZED VIEW.
	DDL
	// DML changes table data: INSERT/UPDATE/DELETE/MERGE/TRUNCATE/CALL.
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

// firstKeywordType maps a leading-keyword TokenKind to its QueryType. It is
// populated in init() via parser.KeywordToken so we never depend on the numeric
// values of the parser's unexported kw* constants.
//
// SELECT-family leads (SELECT/WITH/TABLE/VALUES) map to Select and may be
// refined to SelectInfoSchema by Classify when the statement text references a
// system schema (matching the legacy containsSystemSchema check). EXPLAIN maps
// to Explain. SHOW/DESCRIBE map to SelectInfoSchema directly.
var firstKeywordType map[parser.TokenKind]QueryType

func init() {
	firstKeywordType = make(map[parser.TokenKind]QueryType)

	register := func(qt QueryType, names ...string) {
		for _, name := range names {
			kind, ok := parser.KeywordToken(name)
			if !ok {
				panic("trino/analysis: unknown keyword: " + name)
			}
			firstKeywordType[kind] = qt
		}
	}

	// SELECT-family (refined to SelectInfoSchema by Classify on system schemas).
	register(Select,
		"SELECT",
		"WITH",
		"TABLE",
		"VALUES",
		// SET/RESET SESSION are read-only session statements; the legacy
		// queryTypeListener classifies them as Select (read-only, no data).
		"SET",
		"RESET",
	)

	// EXPLAIN / EXPLAIN ANALYZE.
	register(Explain, "EXPLAIN")

	// Information-schema reads: SHOW family and DESCRIBE/DESC.
	register(SelectInfoSchema,
		"SHOW",
		"DESCRIBE",
		"DESC",
	)

	// Data manipulation language.
	register(DML,
		"INSERT",
		"UPDATE",
		"DELETE",
		"MERGE",
		"TRUNCATE",
		"CALL",
	)

	// Data definition language and access control.
	register(DDL,
		"CREATE",
		"ALTER",
		"DROP",
		"COMMENT",
		"ANALYZE",
		"GRANT",
		"REVOKE",
		"DENY",
		"REFRESH",
		"USE",
	)
}

// tokEOF is the zero value the lexer returns at end-of-input (matches the
// unexported parser.tokEOF = 0).
const tokEOF parser.TokenKind = 0

// kwANALYZE / kwVERBOSE are the token kinds of the EXPLAIN ANALYZE [VERBOSE]
// prefix keywords, resolved once via parser.KeywordToken (init-checked).
var (
	kwANALYZE = mustKeyword("ANALYZE")
	kwVERBOSE = mustKeyword("VERBOSE")
)

// mustKeyword resolves a keyword name to its token kind, panicking if the
// parser does not know it (a build-time invariant, like the register helper).
func mustKeyword(name string) parser.TokenKind {
	kind, ok := parser.KeywordToken(name)
	if !ok {
		panic("trino/analysis: unknown keyword: " + name)
	}
	return kind
}

// systemSchemaPrefixes are the Trino system/metadata schema qualifiers that
// promote a SELECT to SelectInfoSchema. Lifted from the legacy
// containsSystemSchema check (plugin/parser/trino/query_type.go) so the omni
// classifier reports the same read-only-but-system verdict.
var systemSchemaPrefixes = []string{
	"system.",
	"information_schema.",
	"$system.",
	"catalog.",
	"metadata.",
}

// Classify determines the QueryType for a single Trino statement by lexing its
// first meaningful token. It does NOT fully parse the statement, so it is robust
// even for statement kinds whose parser bodies are not yet implemented.
//
// Whitespace and comments are skipped by the lexer. A leading '(' (a
// parenthesized query) is treated as the SELECT family. SELECT-family results
// are promoted to SelectInfoSchema when the statement text references a Trino
// system/information schema, matching the legacy classifier. Empty input or an
// unrecognised leading token yields Unknown.
func Classify(statement string) QueryType {
	l := parser.NewLexer(statement)
	tok := l.NextToken()
	if tok.Kind == tokEOF {
		return Unknown
	}

	// A parenthesized query — `( SELECT ... )` / `( TABLE ... )` — leads with
	// '(' rather than a keyword; treat it as the SELECT family.
	if tok.Kind == int('(') {
		return refineSelect(statement, Select)
	}

	qt, ok := firstKeywordType[tok.Kind]
	if !ok {
		return Unknown
	}
	if qt == Select {
		return refineSelect(statement, qt)
	}
	if qt == Explain {
		return classifyExplain(statement, l)
	}
	return qt
}

// classifyExplain refines the classification of an EXPLAIN statement. The
// leading EXPLAIN keyword has already been consumed; l is positioned just after
// it.
//
// Plain EXPLAIN does not execute its inner statement, so it is read-only
// regardless of the inner kind. EXPLAIN ANALYZE, however, RUNS the inner
// statement (oracle-confirmed against Trino 481: EXPLAIN ANALYZE INSERT
// performs the insert), so an EXPLAIN ANALYZE of a mutating statement is NOT
// read-only and is classified by the inner statement's kind (DML/DDL). This is
// a deliberate, oracle-backed improvement over the legacy classifier, which
// reports EXPLAIN ANALYZE as read-only Select irrespective of the inner kind —
// see the migration divergence ledger.
//
// For a read-only or non-mutating inner statement the classification stays
// Explain, except that a reference to a system/information schema promotes it to
// SelectInfoSchema (matching the legacy containsSystemSchema rule).
func classifyExplain(statement string, l *parser.Lexer) QueryType {
	analyze := false
	// Skip the optional `ANALYZE [VERBOSE]` and `( option, … )` prefix, then read
	// the inner statement's first keyword.
	for {
		tok := l.NextToken()
		if tok.Kind == tokEOF {
			break
		}
		if tok.Kind == kwANALYZE {
			analyze = true
			continue
		}
		if tok.Kind == kwVERBOSE {
			continue
		}
		if tok.Kind == int('(') {
			// Skip a balanced `( … )` option list.
			skipParenGroup(l)
			continue
		}
		// First non-prefix token is the inner statement's leading keyword.
		if analyze {
			if inner, ok := firstKeywordType[tok.Kind]; ok && (inner == DML || inner == DDL) {
				return inner
			}
		}
		break
	}
	// Non-mutating (or plain) EXPLAIN: read-only. Promote to SelectInfoSchema on
	// a system-schema reference, else Explain.
	if hasSystemSchema(statement) {
		return SelectInfoSchema
	}
	return Explain
}

// skipParenGroup consumes tokens through the ')' that balances an already-read
// '(' (used to skip an EXPLAIN option list).
func skipParenGroup(l *parser.Lexer) {
	depth := 1
	for depth > 0 {
		tok := l.NextToken()
		if tok.Kind == tokEOF {
			return
		}
		switch tok.Kind {
		case int('('):
			depth++
		case int(')'):
			depth--
		}
	}
}

// refineSelect promotes a Select classification to SelectInfoSchema when the
// statement text references a system/information schema. SET/RESET SESSION are
// left as Select (they have no table reference). The check is a case-insensitive
// substring scan over the whole statement, exactly as the legacy
// containsSystemSchema does.
func refineSelect(statement string, base QueryType) QueryType {
	if base != Select {
		return base
	}
	if hasSystemSchema(statement) {
		return SelectInfoSchema
	}
	return base
}

// hasSystemSchema reports whether the statement text references a Trino
// system/information schema, via the same case-insensitive substring scan the
// legacy containsSystemSchema uses.
func hasSystemSchema(statement string) bool {
	lower := strings.ToLower(statement)
	for _, prefix := range systemSchemaPrefixes {
		if strings.Contains(lower, prefix) {
			return true
		}
	}
	return false
}
