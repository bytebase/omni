// Package completion provides auto-complete candidates for PartiQL queries.
package completion

import (
	"sort"
	"strings"

	"github.com/bytebase/omni/partiql/catalog"
)

// Candidate represents one auto-complete suggestion.
type Candidate struct {
	Text string // the completion text
	Kind string // "keyword", "table", "column", etc.
}

// Complete returns auto-complete candidates for the given input at the
// cursor position. Uses the catalog for table names and a built-in
// keyword list for SQL keywords.
//
// The implementation is intentionally simple: it looks at the token
// before the cursor to determine context (after FROM → suggest tables,
// after SELECT/WHERE/etc → suggest columns and keywords, default →
// suggest keywords).
func Complete(input string, pos int, cat *catalog.Catalog) []Candidate {
	// Clamp the cursor position into [0, len(input)] before any slicing.
	// Complete is exported, so a caller may pass a position past the end of
	// the input (e.g. a stale cursor). extractPrefix clamps internally, but
	// the input[:pos-len(prefix)] slice below uses pos directly, so without
	// this clamp an out-of-range pos panics with "slice bounds out of range".
	if pos > len(input) {
		pos = len(input)
	}
	if pos < 0 {
		pos = 0
	}

	// Get the partial word at cursor position
	prefix := extractPrefix(input, pos)

	var candidates []Candidate

	// Context detection: look backwards for FROM keyword.
	// Strip the partial word being typed so we look at the token before it.
	beforeCursor := strings.ToUpper(input[:pos-len(prefix)])
	inFromContext := isInFromContext(beforeCursor)

	if inFromContext {
		// After FROM: suggest table names
		for _, table := range cat.Tables() {
			if matchesPrefix(table, prefix) {
				candidates = append(candidates, Candidate{Text: table, Kind: "table"})
			}
		}
	}

	// Always suggest matching keywords
	for _, kw := range keywords {
		if matchesPrefix(kw, prefix) {
			candidates = append(candidates, Candidate{Text: kw, Kind: "keyword"})
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Kind != candidates[j].Kind {
			return candidates[i].Kind < candidates[j].Kind
		}
		return candidates[i].Text < candidates[j].Text
	})
	return candidates
}

func extractPrefix(input string, pos int) string {
	if pos > len(input) {
		pos = len(input)
	}
	start := pos
	for start > 0 && isIdentChar(input[start-1]) {
		start--
	}
	return input[start:pos]
}

func isIdentChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') || ch == '_' || ch == '$'
}

func isInFromContext(upper string) bool {
	// The cursor is in a table-name position when the last completed clause
	// keyword is FROM / JOIN / INTO. Each must match as a STANDALONE keyword:
	// a dotted path step such as "x.from" (member access) also ends in "FROM",
	// but `from` after a '.' is a path key, not a FROM clause, so it must NOT
	// trigger table suggestions.
	upper = strings.TrimRight(upper, " \t\n\r")
	return endsWithKeyword(upper, "FROM") ||
		endsWithKeyword(upper, "JOIN") ||
		endsWithKeyword(upper, "INTO")
}

// endsWithKeyword reports whether s ends with kw as a standalone keyword: kw is
// a suffix of s, and either kw is the whole string or the character immediately
// before it is whitespace (a real clause boundary). This rejects a keyword that
// is merely the tail of a larger token — most importantly a '.'-prefixed path
// step like "X.FROM" (member access x.from), which is not a FROM clause.
func endsWithKeyword(s, kw string) bool {
	if !strings.HasSuffix(s, kw) {
		return false
	}
	if len(s) == len(kw) {
		return true
	}
	switch s[len(s)-len(kw)-1] {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}

func matchesPrefix(text, prefix string) bool {
	if prefix == "" {
		return true
	}
	return strings.HasPrefix(strings.ToUpper(text), strings.ToUpper(prefix))
}

// keywords is the list of PartiQL keywords suggested for auto-completion.
// This is a subset of the full keyword list — only commonly-used keywords.
var keywords = []string{
	"SELECT", "FROM", "WHERE", "AND", "OR", "NOT",
	"INSERT", "INTO", "VALUE", "VALUES", "UPDATE", "SET",
	"DELETE", "CREATE", "DROP", "TABLE", "INDEX",
	"AS", "AT", "BY", "ON", "IN", "IS", "LIKE", "BETWEEN",
	"ORDER", "GROUP", "HAVING", "LIMIT", "OFFSET",
	"JOIN", "LEFT", "RIGHT", "INNER", "OUTER", "CROSS",
	"UNION", "INTERSECT", "EXCEPT", "DISTINCT", "ALL",
	"NULL", "MISSING", "TRUE", "FALSE",
	"CAST", "CASE", "WHEN", "THEN", "ELSE", "END",
	"EXISTS", "ASC", "DESC",
}
