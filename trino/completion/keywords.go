package completion

import "github.com/bytebase/omni/trino/parser"

// Curated keyword candidate lists. The bytebase Trino completer surfaces
// keywords from the c3 follow-set; without c3 we offer hand-picked sets keyed
// to the context. The spellings are upper-case (Trino accepts any case, and
// upper-case is the conventional display form). Sets are intentionally focused
// — a completer that dumps all ~290 keywords at every caret is noise.

// statementStartKeywords are offered at the start of a statement: the
// statement-leading keywords of the Trino grammar's statement rule.
var statementStartKeywords = []string{
	"ALTER",
	"ANALYZE",
	"CALL",
	"COMMENT",
	"COMMIT",
	"CREATE",
	"DEALLOCATE",
	"DELETE",
	"DENY",
	"DESCRIBE",
	"DROP",
	"EXECUTE",
	"EXPLAIN",
	"GRANT",
	"INSERT",
	"MERGE",
	"PREPARE",
	"REFRESH",
	"RESET",
	"REVOKE",
	"ROLLBACK",
	"SELECT",
	"SET",
	"SHOW",
	"START",
	"TABLE",
	"TRUNCATE",
	"UPDATE",
	"USE",
	"VALUES",
	"WITH",
}

// clauseKeywords are offered as a fallback (e.g. just after a complete table
// name): the clause keywords that commonly continue a query.
var clauseKeywords = []string{
	"AS",
	"CROSS",
	"EXCEPT",
	"FETCH",
	"FROM",
	"FULL",
	"GROUP",
	"HAVING",
	"INNER",
	"INTERSECT",
	"JOIN",
	"LEFT",
	"LIMIT",
	"NATURAL",
	"OFFSET",
	"ON",
	"ORDER",
	"RIGHT",
	"UNION",
	"USING",
	"WHERE",
	"WINDOW",
}

// expressionKeywords are offered in a value/predicate position (after SELECT,
// WHERE, ON, AND/OR, a comparison): the operators and constructs that begin or
// continue an expression, plus the literal keywords.
var expressionKeywords = []string{
	"ALL",
	"AND",
	"ANY",
	"ARRAY",
	"BETWEEN",
	"CASE",
	"CAST",
	"CURRENT_DATE",
	"CURRENT_TIME",
	"CURRENT_TIMESTAMP",
	"DISTINCT",
	"EXISTS",
	"FALSE",
	"IN",
	"INTERVAL",
	"IS",
	"LIKE",
	"NOT",
	"NULL",
	"OR",
	"SOME",
	"TRUE",
	"TRY_CAST",
}

// QuoteIdentifierIfNeeded returns name wrapped in double quotes when it would
// otherwise be lexed as something other than a single unquoted identifier, and
// returns it unchanged when it is a clean unquoted identifier.
//
// A name needs quoting when it is empty, does not start with a letter or '_',
// contains a character outside [A-Za-z0-9_], or collides with a reserved Trino
// keyword (an unquoted reserved word cannot stand for an identifier). Embedded
// double quotes are doubled per Trino's quoting rule. This mirrors the legacy
// bytebase quotedIdentifierIfNeeded helper.
//
// Note: a name produced by catalog normalization is already lower-cased for an
// originally-unquoted identifier, or case-preserved for an originally-quoted
// one; we re-quote only when the *characters* require it, so a normalized
// "MyTable" (which came from a quoted source) is re-quoted to preserve its
// case, which is the correct, round-trippable completion.
func QuoteIdentifierIfNeeded(name string) string {
	if needsQuoting(name) {
		return `"` + escapeDoubleQuotes(name) + `"`
	}
	return name
}

// needsQuoting reports whether name cannot be written as a bare unquoted Trino
// identifier.
func needsQuoting(name string) bool {
	if name == "" {
		return true
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case c >= 'a' && c <= 'z', c == '_':
			// A clean lower-case identifier character; always allowed unquoted.
		case c >= 'A' && c <= 'Z':
			// An upper-case letter forces quoting: Trino folds an *unquoted*
			// identifier to lower case, so emitting MyCol unquoted would resolve
			// to mycol. A normalized name retaining upper case came from a quoted
			// source and is case-sensitive, so it must be re-quoted to round-trip.
			return true
		case c >= '0' && c <= '9':
			if i == 0 {
				// A leading digit makes it a DIGIT_IDENTIFIER at best; quote so
				// the completion is an unambiguous identifier.
				return true
			}
		default:
			return true
		}
	}
	// Upper-case the whole word and check against the reserved set: a reserved
	// keyword used unquoted is not a valid identifier.
	return isReservedWord(name)
}

// isReservedWord reports whether name is a reserved Trino keyword (one that
// cannot be used as an unquoted identifier). Non-reserved keywords (e.g. DATA,
// FORMAT) are valid unquoted identifiers and need no quoting.
func isReservedWord(name string) bool {
	kind, ok := parser.KeywordToken(name)
	if !ok {
		return false
	}
	return parser.IsReserved(kind)
}

// escapeDoubleQuotes doubles any embedded double-quote, per Trino's
// "" escape inside a quoted identifier.
func escapeDoubleQuotes(s string) string {
	out := make([]byte, 0, len(s)+2)
	for i := 0; i < len(s); i++ {
		if s[i] == '"' {
			out = append(out, '"', '"')
			continue
		}
		out = append(out, s[i])
	}
	return string(out)
}
