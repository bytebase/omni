package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// The paren-route tests hold how the parser routes a `(` in a FROM item, an
// IN list or an ARRAY constructor: to a subquery, a joined table, something
// else, or a rejection. Each case's expected route is committed beside it.
//
// The cases come from pg/parser/SCENARIOS-pg-paren-dispatch.md, and their
// verdicts were measured on PostgreSQL 17 when this fork shared pg's oracle.
// Redshift runs no PostgreSQL 17 and CI can start no Redshift, so the
// verdicts are the parser's own now, not a claim about what Redshift
// accepts.

// OmniStatus classifies how the parser routed the paren-enclosed FROM item.
// The classifier inspects the parsed AST of `SELECT * FROM (...)` and reports
// which node was produced.
type OmniStatus int

const (
	// OmniSubquery — FROM item is *nodes.RangeSubselect, i.e. the `(`
	// was routed to select_with_parens.
	OmniSubquery OmniStatus = iota
	// OmniJoinedTable — FROM item is *nodes.JoinExpr, i.e. `(` was
	// routed to `(` joined_table `)`.
	OmniJoinedTable
	// OmniOther — FROM item parsed but is neither subquery nor
	// joined_table (RangeVar, RangeFunction, RangeTableFunc, JsonTable,
	// etc.), or the statement has no FROM item to route.
	OmniOther
	// OmniRejected — the parser returned an error.
	OmniRejected
)

func (s OmniStatus) String() string {
	switch s {
	case OmniSubquery:
		return "subquery"
	case OmniJoinedTable:
		return "joined_table"
	case OmniOther:
		return "other"
	case OmniRejected:
		return "reject"
	default:
		return "unknown"
	}
}

// classifyOmniAt parses sqlStr and classifies the FROM item at the 0-based
// index of its first SELECT. LATERAL cases put the probed node at Items[1],
// after the anchor relation T. A statement that is not a SELECT, or has no
// FROM clause, is OmniOther; an index out of range is OmniRejected.
func classifyOmniAt(sqlStr string, index int) (OmniStatus, string) {
	stmts, err := Parse(sqlStr)
	if err != nil {
		return OmniRejected, err.Error()
	}
	if stmts == nil || len(stmts.Items) == 0 {
		return OmniRejected, "no statements"
	}
	raw, ok := stmts.Items[0].(*nodes.RawStmt)
	if !ok {
		return OmniOther, ""
	}
	sel, ok := raw.Stmt.(*nodes.SelectStmt)
	if !ok {
		return OmniOther, ""
	}
	if sel.FromClause == nil || len(sel.FromClause.Items) == 0 {
		// `SELECT 1` with no FROM — not a paren-routing question.
		return OmniOther, ""
	}
	if index < 0 || index >= len(sel.FromClause.Items) {
		return OmniRejected, "FROM item index out of bounds"
	}
	switch sel.FromClause.Items[index].(type) {
	case *nodes.RangeSubselect:
		return OmniSubquery, ""
	case *nodes.JoinExpr:
		return OmniJoinedTable, ""
	default:
		// RangeVar, RangeFunction, RangeTableFunc, JsonTable, etc.
		return OmniOther, ""
	}
}

// assertParenRoute fails t unless the parser routes the first FROM item of
// sqlStr as expected.
func assertParenRoute(t *testing.T, sqlStr string, expected OmniStatus) {
	t.Helper()
	assertParenRouteAt(t, sqlStr, 0, expected)
}

// assertParenRouteAt is assertParenRoute for the FROM item at index.
func assertParenRouteAt(t *testing.T, sqlStr string, index int, expected OmniStatus) {
	t.Helper()
	if got, msg := classifyOmniAt(sqlStr, index); got != expected {
		t.Errorf("route of %q: got %s (error %q), want %s", summarize(sqlStr), got, msg, expected)
	}
}

// TestParenRouteHarness checks the classifier on one case of each route.
func TestParenRouteHarness(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		expected OmniStatus
	}{
		{"joined_table single-level", `SELECT * FROM (T JOIN U ON TRUE)`, OmniJoinedTable},
		{"subquery single-level", `SELECT * FROM (SELECT 1)`, OmniSubquery},
		{"empty parens rejected", `SELECT * FROM ()`, OmniRejected},
		{"single relation in parens rejected", `SELECT * FROM (T)`, OmniRejected},
		{"BYT-9315 shape", `SELECT * FROM ((T JOIN U ON TRUE) JOIN V ON TRUE)`, OmniJoinedTable},
		{"bare relation (OmniOther)", `SELECT * FROM T`, OmniOther},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertParenRoute(t, tc.sql, tc.expected)
		})
	}
}

// summarize trims long SQL for messages.
func summarize(sqlStr string) string {
	const limit = 60
	s := strings.ReplaceAll(sqlStr, "\n", " ")
	if len(s) <= limit {
		return s
	}
	return s[:limit-1] + "…"
}
