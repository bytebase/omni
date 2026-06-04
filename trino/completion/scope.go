package completion

import (
	"strings"

	"github.com/bytebase/omni/trino/analysis"
	"github.com/bytebase/omni/trino/catalog"
	"github.com/bytebase/omni/trino/parser"
)

// placeholder is the sentinel identifier inserted at the caret to make an
// in-progress statement parse, so query-span analysis can recover the FROM
// scope and CTE names that surround the caret. It is an unlikely real
// identifier and is filtered out of any candidate it might masquerade as.
const placeholder = "__omni_completion_placeholder__"

// scopeTable is one table/view in the FROM scope of the statement under the
// caret, as reported by query-span analysis. Names are already Trino-normalized
// (the analysis package folds them), so they key directly into the catalog.
type scopeTable struct {
	catalog string // normalized; "" when the source name omitted it
	schema  string // normalized; "" when omitted
	table   string // normalized table/view name
	alias   string // normalized alias; "" when none
}

// analyzeAtCaret splits sql, finds the statement containing the caret (at byte
// offset limit, the caret with its partial word already stripped), inserts the
// placeholder identifier there to keep the statement parseable, and runs
// query-span analysis on just that statement.
//
// Why the whole containing statement and not the prefix: for SELECT-list and
// predicate completion the FROM clause is to the RIGHT of the caret, so
// truncating at the caret would hide the very tables we need to offer columns
// for. Analyzing the full statement (with a placeholder closing the
// half-written clause at the caret) recovers them; the trailing text after the
// caret is real SQL the user already wrote.
func analyzeAtCaret(sql string, limit int) *analysis.QuerySpan {
	if limit > len(sql) {
		limit = len(sql)
	}
	if limit < 0 {
		limit = 0
	}

	start, end := containingStatement(sql, limit)
	stmt := sql[start:end]
	rel := limit - start // caret offset within the statement

	// Insert the placeholder at the caret, padded with spaces so it never fuses
	// with an adjacent token ("o." + ph must read as "o . ph", and "FROMx" must
	// not arise from "FROM" + ph with no gap — though limit is already at a
	// token boundary, the padding is belt-and-suspenders).
	patched := stmt[:rel] + " " + placeholder + " " + stmt[rel:]

	span, err := analysis.GetQuerySpan(patched)
	if err != nil {
		return nil
	}
	// Robustness fallback: a statement can carry a SECOND incomplete clause away
	// from the caret — most commonly an empty SELECT list ("SELECT  FROM t"
	// while the caret edits the WHERE) — which fails the parse so no FROM scope
	// is recovered. If the caret-patched parse found neither tables nor CTEs,
	// retry once with empty SELECT lists filled by a placeholder select item, so
	// the FROM clause becomes reachable.
	if span != nil && len(span.AccessTables) == 0 && len(span.CTEs) == 0 {
		if filled := fillEmptySelectLists(patched); filled != patched {
			if span2, err2 := analysis.GetQuerySpan(filled); err2 == nil && span2 != nil &&
				(len(span2.AccessTables) > 0 || len(span2.CTEs) > 0) {
				return span2
			}
		}
	}
	return span
}

// fillEmptySelectLists inserts a placeholder select item into any "SELECT FROM"
// (empty select list) occurrence so the statement parses. It is a coarse text
// transform used only as a scope-recovery fallback; it keys on a SELECT keyword
// token immediately followed by a FROM keyword token. Returns the input
// unchanged when there is nothing to fill.
func fillEmptySelectLists(stmt string) string {
	toks, _ := parser.Tokenize(stmt)
	// Find the first SELECT immediately followed by FROM and splice a star in
	// between. Doing one is enough to make the common single-query case parse;
	// nested empties are rare enough to leave to the caret placeholder.
	for i := 0; i+1 < len(toks); i++ {
		if isKeyword(toks[i], "select") && isKeyword(toks[i+1], "from") {
			at := toks[i].Loc.End // byte just after SELECT
			return stmt[:at] + " * " + stmt[at:]
		}
	}
	return stmt
}

// containingStatement returns the [start, end) byte range of the top-level
// statement that contains the caret at offset limit. A caret sitting on a ';'
// boundary, or past the last statement, attaches to the nearest statement so
// completion still has context. Falls back to the whole input when Split finds
// nothing.
func containingStatement(sql string, limit int) (start, end int) {
	segs := parser.Split(sql)
	if len(segs) == 0 {
		return 0, len(sql)
	}
	for _, s := range segs {
		// Caret within or at the trailing edge of this segment.
		if limit >= s.ByteStart && limit <= s.ByteEnd {
			return s.ByteStart, s.ByteEnd
		}
	}
	// Caret is in inter-statement whitespace / after the last ';': attach to the
	// last segment that starts at or before the caret, else the first segment.
	chosen := segs[0]
	for _, s := range segs {
		if s.ByteStart <= limit {
			chosen = s
		}
	}
	return chosen.ByteStart, chosen.ByteEnd
}

// fromScope returns the tables in scope for column completion in the statement
// containing the caret: the FROM/JOIN tables query-span analysis reports, plus
// the target relation of a DML statement (UPDATE/DELETE/INSERT/MERGE), which the
// span's AccessTables omits because it tracks tables a query *reads*, not a
// DML target. The placeholder pseudo-table (present when the caret itself is the
// relation being typed, e.g. "FROM |") is excluded.
func fromScope(sql string, limit int) []scopeTable {
	var result []scopeTable
	if span := analyzeAtCaret(sql, limit); span != nil {
		for _, t := range span.AccessTables {
			if t.Table == placeholder {
				continue
			}
			result = append(result, scopeTable{
				catalog: t.Catalog,
				schema:  t.Schema,
				table:   t.Table,
				alias:   t.Alias,
			})
		}
	}
	result = append(result, dmlTargetTables(sql, limit)...)
	return result
}

// dmlTargetTables returns the target relation that forms the column-completion
// scope of a DML statement, which query-span does not surface (it records only
// read tables). The target follows the statement's LEADING keyword:
//
//   - UPDATE <target> / DELETE FROM <target>: the target is the column scope for
//     SET and WHERE. Trino does NOT allow a target alias here (oracle 481:
//     "UPDATE t a SET ..." => SYNTAX_ERROR), so no alias is parsed.
//   - MERGE INTO <target> [AS] alias: the target IS aliasable (oracle 481
//     accepts the alias), and merge-clause columns reference it, so the alias is
//     parsed.
//
// INSERT is intentionally excluded: its column scope is the SOURCE query (which
// query-span already provides) or the explicit "(col, ...)" list; the target
// table's columns are not visible in the source SELECT, so adding them would
// over-offer (e.g. "INSERT INTO t SELECT |" must not suggest t's columns).
// TRUNCATE has no column scope.
//
// Only the leading position is inspected, so a DML keyword appearing
// mid-statement (e.g. the UPDATE inside a MERGE's WHEN MATCHED THEN UPDATE) is
// never mistaken for a fresh target. Names are normalized to the catalog key
// form.
func dmlTargetTables(sql string, limit int) []scopeTable {
	start, end := containingStatement(sql, limit)
	toks, _ := parser.Tokenize(sql[start:end])
	if len(toks) == 0 {
		return nil
	}

	// The target name index and whether an alias is grammatical, from the
	// statement's leading keyword(s).
	j := -1
	allowAlias := false
	switch {
	case isKeywordAt(toks, 0, "update"):
		j = 1
	case isKeywordAt(toks, 0, "delete") && isKeywordAt(toks, 1, "from"):
		j = 2
	case isKeywordAt(toks, 0, "merge") && isKeywordAt(toks, 1, "into"):
		j, allowAlias = 2, true
	}
	if j < 0 {
		return nil
	}
	if st, _ := readDottedRelation(toks, j, allowAlias); st != nil {
		return []scopeTable{*st}
	}
	return nil
}

// readDottedRelation reads a 1–3 part dotted relation name beginning at toks[j],
// returning the scopeTable and the number of tokens consumed past j. Returns nil
// when toks[j] is not a name. When allowAlias is true an optional trailing
// "AS x" / bare-identifier alias is consumed (only MERGE targets are aliasable);
// when false no alias is parsed (UPDATE/DELETE targets reject one).
func readDottedRelation(toks []parser.Token, j int, allowAlias bool) (*scopeTable, int) {
	if j >= len(toks) || !isNameToken(toks[j]) {
		return nil, 0
	}
	var parts []string
	k := j
	for k < len(toks) && isNameToken(toks[k]) {
		parts = append(parts, normalizeQualifierPart(toks[k]))
		k++
		if k < len(toks) && toks[k].Kind == int('.') {
			k++
			continue
		}
		break
	}
	st := &scopeTable{}
	switch len(parts) {
	case 1:
		st.table = parts[0]
	case 2:
		st.schema, st.table = parts[0], parts[1]
	default:
		st.catalog, st.schema, st.table = parts[len(parts)-3], parts[len(parts)-2], parts[len(parts)-1]
	}
	// Optional alias (MERGE only): "AS x" or a bare following identifier (a
	// keyword such as SET that starts the next clause is not an alias).
	if allowAlias && k < len(toks) {
		if isKeyword(toks[k], "as") && k+1 < len(toks) && isPlainIdent(toks[k+1]) {
			st.alias = normalizeQualifierPart(toks[k+1])
			k += 2
		} else if isPlainIdent(toks[k]) {
			st.alias = normalizeQualifierPart(toks[k])
			k++
		}
	}
	return st, k - j
}

// isKeywordAt reports whether toks[i] is the given keyword (bounds-checked).
func isKeywordAt(toks []parser.Token, i int, kw string) bool {
	return i >= 0 && i < len(toks) && isKeyword(toks[i], kw)
}

// isPlainIdent reports whether tok is an (unquoted/quoted) identifier rather
// than a keyword — used so a clause keyword following a DML target is not
// mistaken for its alias.
func isPlainIdent(tok parser.Token) bool {
	switch parser.TokenName(tok.Kind) {
	case "IDENTIFIER", "QUOTED_IDENTIFIER", "BACKQUOTED_IDENTIFIER", "DIGIT_IDENTIFIER":
		return true
	}
	return false
}

// inScopeColumnNames returns the column-name lists of all FROM-scope tables that
// resolve against the catalog, one slice per resolved table. A table that does
// not resolve (unknown to the catalog, or whose session context is unset)
// contributes nothing.
func inScopeColumnNames(sql string, limit int, cat *catalog.Catalog) [][]string {
	var out [][]string
	for _, st := range fromScope(sql, limit) {
		if names := resolveColumns(cat, st); len(names) > 0 {
			out = append(out, names)
		}
	}
	return out
}

// scopeColumnsFor returns the column names of the single FROM-scope table the
// qualifier x selects, where x is either a table's alias or its (normalized)
// table name. ok is false when x matches no in-scope table. x is normalized by
// the caller (collectQualifier) to the catalog's key form.
func scopeColumnsFor(sql string, limit int, cat *catalog.Catalog, x string) (names []string, ok bool) {
	for _, st := range fromScope(sql, limit) {
		// An alias, when present, shadows the table name as the column
		// qualifier (SELECT o.id FROM orders o references via the alias).
		if st.alias != "" {
			if st.alias == x {
				return resolveColumns(cat, st), true
			}
			continue
		}
		if st.table == x {
			return resolveColumns(cat, st), true
		}
	}
	return nil, false
}

// resolveColumns resolves a scope table to its catalog column names, filling an
// omitted catalog/schema from the session context the way Trino does. Returns
// nil when the table is unknown to the catalog or the needed session context is
// unset.
func resolveColumns(cat *catalog.Catalog, st scopeTable) []string {
	if cat == nil {
		return nil
	}
	parts := make([]string, 0, 3)
	if st.catalog != "" {
		parts = append(parts, st.catalog)
	}
	if st.schema != "" {
		parts = append(parts, st.schema)
	}
	parts = append(parts, st.table)
	return cat.ResolveRelation(parts).ColumnNames()
}

// cteNames returns the CTE names declared in the statement under the caret,
// derived from query-span analysis with the caret placeholder. They are
// select-able relations the completer offers after FROM/JOIN. Names are
// normalized; the placeholder is never among them.
func cteNames(sql string, limit int) []string {
	span := analyzeAtCaret(sql, limit)
	if span == nil {
		return nil
	}
	seen := make(map[string]struct{}, len(span.CTEs))
	out := make([]string, 0, len(span.CTEs))
	for _, name := range span.CTEs {
		if name == placeholder {
			continue
		}
		key := strings.ToLower(name)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, name)
	}
	return out
}
