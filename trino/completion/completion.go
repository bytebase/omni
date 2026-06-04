// Package completion provides parser-native auto-complete candidates for Trino
// SQL, the omni counterpart of bytebase's plugin/parser/trino Completion.
//
// # Why this is not a C3 port
//
// The bytebase Trino completer is a ~1730-line CodeCompletionCore (c3 /
// follow-set) implementation that walks the ANTLR parser's rule indices. The
// omni Trino parser is a hand-written recursive-descent parser that was not
// built with c3 instrumentation (no rule-candidate collection hooks), and this
// node may only touch trino/completion/*.go — it cannot weave collection points
// through trino/parser. So, like omni's mongo and partiql completers, this
// package is self-contained: it tokenizes the input up to the caret with the
// exported lexer, classifies the completion context from the token sequence,
// and produces candidates from the three-level Trino catalog
// (catalog -> schema -> table/view -> column) plus a curated keyword set.
//
// FROM-clause column candidates are resolved the same way bytebase's completer
// ultimately resolves them — through query-span analysis (trino/analysis):
// GetQuerySpan yields the tables a SELECT reads, which the catalog maps to
// column candidates. This keeps the candidate *set* faithful to the legacy
// completer's semantics even though the mechanism differs.
//
// The package depends only on trino/parser (lexer + identifier helpers),
// trino/catalog (object metadata), and trino/analysis (FROM-scope resolution).
package completion

import (
	"sort"
	"strings"

	"github.com/bytebase/omni/trino/catalog"
)

// CandidateType classifies a completion candidate so a consumer (LSP) can pick
// an icon / sort bucket. The values mirror the candidate kinds the bytebase
// Trino completer emits: keyword, catalog (database), schema, table, view,
// column. Functions are categorized by the legacy completer but its function
// table is never populated, so this package does not emit function candidates.
type CandidateType int

const (
	// CandidateKeyword is a SQL keyword (SELECT, FROM, JOIN, ...).
	CandidateKeyword CandidateType = iota
	// CandidateCatalog is a Trino catalog name (the top namespace level).
	CandidateCatalog
	// CandidateSchema is a schema name within a catalog.
	CandidateSchema
	// CandidateTable is a base-table name.
	CandidateTable
	// CandidateView is a view (or materialized view) name.
	CandidateView
	// CandidateColumn is a column name of an in-scope table/view/CTE.
	CandidateColumn
)

// String returns the lower-case kind name, handy for tests and logging.
func (t CandidateType) String() string {
	switch t {
	case CandidateKeyword:
		return "keyword"
	case CandidateCatalog:
		return "catalog"
	case CandidateSchema:
		return "schema"
	case CandidateTable:
		return "table"
	case CandidateView:
		return "view"
	case CandidateColumn:
		return "column"
	default:
		return "unknown"
	}
}

// Candidate is a single completion suggestion. Text is the completion string
// (already quoted when the identifier needs quoting, via QuoteIdentifierIfNeeded
// at the call sites that emit catalog object names).
type Candidate struct {
	// Text is the suggestion to insert.
	Text string
	// Type classifies the suggestion.
	Type CandidateType
}

// Complete returns completion candidates for sql at byte offset cursorOffset.
// cat may be nil, in which case only keyword candidates (and CTE names found in
// the statement text) are produced — object candidates need catalog metadata.
//
// The cursorOffset is clamped into [0, len(sql)]; Complete never panics on an
// out-of-range offset (it is exported and a stale LSP caret may overrun the
// text). Candidates are returned sorted by (Type, Text) and de-duplicated, and
// filtered to those whose Text matches the partial word under the caret
// (case-insensitively — Trino folds unquoted identifiers and keywords to lower
// case, so a lower-case prefix must match a Candidate regardless of its case).
func Complete(sql string, cursorOffset int, cat *catalog.Catalog) []Candidate {
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}
	if cursorOffset < 0 {
		cursorOffset = 0
	}

	// The partial word under the caret is stripped before context detection so
	// the lexer sees the position *before* the partial text: completing
	// "... FROM ord|" must classify as a FROM context (offering tables), not as
	// "inside an identifier". The stripped word becomes the prefix filter.
	prefix := extractPrefix(sql, cursorOffset)
	scanLimit := cursorOffset - len(prefix)

	ctx := detectContext(sql, scanLimit)
	candidates := candidatesFor(ctx, sql, scanLimit, cat)

	return filterByPrefix(dedup(candidates), prefix)
}

// extractPrefix returns the partial identifier word ending at cursorOffset.
// A Trino unquoted identifier is (LETTER|'_')(LETTER|DIGIT|'_')*; we also accept
// a leading digit so a digit-leading identifier (DIGIT_IDENTIFIER_) being typed
// is treated as a prefix rather than a number. The prefix never includes a '.',
// so "schema." yields an empty prefix (the caret is positioned to complete the
// object after the dot).
func extractPrefix(sql string, cursorOffset int) string {
	i := cursorOffset
	for i > 0 && isIdentByte(sql[i-1]) {
		i--
	}
	return sql[i:cursorOffset]
}

// isIdentByte reports whether c can appear in an unquoted/digit Trino identifier.
func isIdentByte(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') || c == '_'
}

// filterByPrefix keeps candidates whose Text has prefix as a case-insensitive
// prefix. An empty prefix keeps everything. Matching is case-insensitive
// because Trino folds unquoted names to lower case: a user typing "ORD" or
// "ord" should both match a table candidate "orders".
func filterByPrefix(candidates []Candidate, prefix string) []Candidate {
	if prefix == "" {
		return candidates
	}
	upper := strings.ToUpper(prefix)
	result := candidates[:0:0]
	for _, c := range candidates {
		if strings.HasPrefix(strings.ToUpper(c.Text), upper) {
			result = append(result, c)
		}
	}
	return result
}

// dedup removes duplicate (Text, Type) pairs and returns the survivors sorted by
// (Type, Text). De-duplication keys on the EXACT Text, not a case-folded form:
// two quoted identifiers that differ only in case ("Foo" vs "foo") are distinct
// Trino objects and must both survive. Names that genuinely coincide already
// share an identical (normalized, then quoted) Text, so exact-match dedup still
// collapses the catalog-and-CTE-both-surface case.
func dedup(candidates []Candidate) []Candidate {
	type key struct {
		text string
		typ  CandidateType
	}
	seen := make(map[key]struct{}, len(candidates))
	result := make([]Candidate, 0, len(candidates))
	for _, c := range candidates {
		k := key{c.Text, c.Type}
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		result = append(result, c)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Type != result[j].Type {
			return result[i].Type < result[j].Type
		}
		return result[i].Text < result[j].Text
	})
	return result
}
