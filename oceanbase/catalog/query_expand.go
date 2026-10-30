package catalog

// ExpandMergeViews creates a copy of the Query where RTERelation entries
// for MERGE views have their Subquery field populated from View.AnalyzedQuery.
// TEMPTABLE views remain opaque.
//
// This is a consume-time operation (decision D5): the analyzer keeps views
// opaque; consumers that want lineage transparency call this method.
//
// The expansion is recursive: if a view references another view, the inner
// view's RTE is also expanded.
func (q *Query) ExpandMergeViews(c *Catalog) *Query {
	if q == nil {
		return nil
	}

	// Shallow copy the query; we only need to replace the RangeTable slice.
	expanded := *q
	expanded.RangeTable = make([]*RangeTableEntryQ, len(q.RangeTable))
	for i, rte := range q.RangeTable {
		if rte.IsView && (rte.ViewAlgorithm == ViewAlgMerge || rte.ViewAlgorithm == ViewAlgUndefined) {
			// Look up the view's AnalyzedQuery from the catalog.
			db := c.GetDatabase(rte.DBName)
			if db != nil {
				view := db.Views[toLower(rte.TableName)]
				if view != nil && view.AnalyzedQuery != nil {
					// Copy RTE and set Subquery to the recursively expanded view body.
					rteCopy := *rte
					rteCopy.Subquery = view.AnalyzedQuery.ExpandMergeViews(c)
					expanded.RangeTable[i] = &rteCopy
					continue
				}
			}
		}
		expanded.RangeTable[i] = rte
	}
	return &expanded
}
