package completion

import (
	"github.com/bytebase/omni/trino/catalog"
)

// candidatesFor produces the candidate list for a classified context. sql and
// scanLimit identify the caret (with the partial word stripped) so column
// contexts can resolve the in-scope tables via query-span analysis. cat may be
// nil.
func candidatesFor(ctx completionContext, sql string, scanLimit int, cat *catalog.Catalog) []Candidate {
	switch ctx.kind {
	case kindStatementStart:
		return keywordCandidates(statementStartKeywords)
	case kindRelation:
		return relationCandidates(sql, scanLimit, cat)
	case kindColumnOrExpr:
		return append(columnCandidates(sql, scanLimit, cat), keywordCandidates(expressionKeywords)...)
	case kindColumnOnly:
		// Assignment target (UPDATE ... SET): columns only, no expression keywords.
		return columnCandidates(sql, scanLimit, cat)
	case kindDottedName:
		return dottedCandidates(ctx.qualifier, sql, scanLimit, cat)
	default: // kindKeyword
		return keywordCandidates(clauseKeywords)
	}
}

// keywordCandidates wraps a keyword list as Candidates.
func keywordCandidates(words []string) []Candidate {
	result := make([]Candidate, 0, len(words))
	for _, w := range words {
		result = append(result, Candidate{Text: w, Type: CandidateKeyword})
	}
	return result
}

// relationCandidates offers the table/view names visible in the session schema,
// plus schema names in the session catalog and all catalog names, so the user
// can either pick a relation directly or drill down a qualifier. CTE names
// defined earlier in the statement are added as table candidates (a CTE is a
// relation you can select FROM). cat may be nil, in which case only CTE names
// from the statement are returned.
func relationCandidates(sql string, scanLimit int, cat *catalog.Catalog) []Candidate {
	var result []Candidate
	if cat != nil {
		// All catalog names (top-level qualifier drill-down).
		for _, name := range cat.Catalogs() {
			result = append(result, Candidate{Text: QuoteIdentifierIfNeeded(name), Type: CandidateCatalog})
		}
		// Schemas + relations in the session catalog/schema.
		if cur := cat.CurrentCatalog(); cur != "" {
			if db := cat.GetCatalog(cur); db != nil {
				for _, sname := range db.Schemas() {
					result = append(result, Candidate{Text: QuoteIdentifierIfNeeded(sname), Type: CandidateSchema})
				}
				if cs := cat.CurrentSchema(); cs != "" {
					if sch := db.GetSchema(cs); sch != nil {
						result = append(result, relationsOf(sch)...)
					}
				}
			}
		}
	}
	// CTE names declared in this statement are select-able relations.
	for _, name := range cteNames(sql, scanLimit) {
		result = append(result, Candidate{Text: QuoteIdentifierIfNeeded(name), Type: CandidateTable})
	}
	return result
}

// relationsOf returns the table and view candidates of a schema.
func relationsOf(sch *catalog.Schema) []Candidate {
	var result []Candidate
	for _, name := range sch.Tables() {
		result = append(result, Candidate{Text: QuoteIdentifierIfNeeded(name), Type: CandidateTable})
	}
	for _, name := range sch.Views() {
		result = append(result, Candidate{Text: QuoteIdentifierIfNeeded(name), Type: CandidateView})
	}
	return result
}

// columnCandidates offers the columns of every table/view in scope at the
// caret. Scope is derived from the statement's FROM clauses via query-span
// analysis (trino/analysis), which yields the accessed tables; each is resolved
// against the catalog to its column list. When no catalog is available, or no
// table resolves, the result is empty (and the caller still offers expression
// keywords). cat may be nil.
func columnCandidates(sql string, scanLimit int, cat *catalog.Catalog) []Candidate {
	if cat == nil {
		return nil
	}
	var result []Candidate
	for _, names := range inScopeColumnNames(sql, scanLimit, cat) {
		for _, c := range names {
			result = append(result, Candidate{Text: QuoteIdentifierIfNeeded(c), Type: CandidateColumn})
		}
	}
	return result
}

// dottedCandidates resolves a qualifier chain ("<parts> .") to its drill-down
// candidates:
//
//   - 1 part [x]:
//     x as a table/alias in the FROM scope -> its columns;
//     x as a schema in the session catalog -> its relations;
//     x as a catalog                        -> its schemas.
//     All matching interpretations are offered (the caret is ambiguous).
//   - 2 parts [x, y]:
//     x.y as catalog.schema -> relations of that schema;
//     x.y as schema.table   -> columns of that table (session catalog).
//   - 3 parts [x, y, z]: x.y.z as catalog.schema.table -> its columns.
//
// cat may be nil, in which case only the FROM-scope alias interpretation (which
// needs the catalog too) is unavailable and the result is empty.
func dottedCandidates(qualifier []string, sql string, scanLimit int, cat *catalog.Catalog) []Candidate {
	if cat == nil || len(qualifier) == 0 {
		return nil
	}
	var result []Candidate

	switch len(qualifier) {
	case 1:
		x := qualifier[0]
		// (a) table/alias in scope -> columns.
		if names, ok := scopeColumnsFor(sql, scanLimit, cat, x); ok {
			for _, c := range names {
				result = append(result, Candidate{Text: QuoteIdentifierIfNeeded(c), Type: CandidateColumn})
			}
		}
		// (b) schema in the session catalog -> relations.
		if cur := cat.CurrentCatalog(); cur != "" {
			if db := cat.GetCatalog(cur); db != nil {
				if sch := db.GetSchema(x); sch != nil {
					result = append(result, relationsOf(sch)...)
				}
			}
		}
		// (c) catalog -> schemas.
		if db := cat.GetCatalog(x); db != nil {
			for _, sname := range db.Schemas() {
				result = append(result, Candidate{Text: QuoteIdentifierIfNeeded(sname), Type: CandidateSchema})
			}
		}

	case 2:
		// (a) catalog.schema -> relations.
		if db := cat.GetCatalog(qualifier[0]); db != nil {
			if sch := db.GetSchema(qualifier[1]); sch != nil {
				result = append(result, relationsOf(sch)...)
			}
		}
		// (b) schema.table (session catalog) -> columns.
		if cur := cat.CurrentCatalog(); cur != "" {
			result = append(result, columnsOfRelation(cat, []string{qualifier[0], qualifier[1]})...)
		}

	case 3:
		// catalog.schema.table -> columns.
		result = append(result, columnsOfRelation(cat, qualifier)...)
	}

	return result
}

// columnsOfRelation resolves parts to a relation and returns its column
// candidates, or nil if it does not resolve.
func columnsOfRelation(cat *catalog.Catalog, parts []string) []Candidate {
	rr := cat.ResolveRelation(parts)
	if !rr.Found {
		return nil
	}
	names := rr.ColumnNames()
	result := make([]Candidate, 0, len(names))
	for _, c := range names {
		result = append(result, Candidate{Text: QuoteIdentifierIfNeeded(c), Type: CandidateColumn})
	}
	return result
}
