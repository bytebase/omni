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
	// Simple heuristic: last SQL keyword before cursor is FROM
	upper = strings.TrimRight(upper, " \t\n\r")
	return strings.HasSuffix(upper, "FROM") ||
		strings.HasSuffix(upper, "JOIN") ||
		strings.HasSuffix(upper, "INTO")
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
	"EXISTS", "NOT", "ASC", "DESC",
}
