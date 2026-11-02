package catalog

import (
	"strings"
	"testing"
)

// TestAnalyzeKeywordArgRoundTrip covers the analyzer's half of the
// keyword-argument contract. The unit is not a column, so it must contribute
// no lineage, and it must deparse bare: the engine lexes the position as a
// keyword and rejects timestampdiff('SECOND', ...), so a quoted unit would not
// re-parse.
//
// The table deliberately holds a column named after a unit. That is the case
// that used to resolve: while the unit was a column reference the analyzer
// succeeded and attributed a read of `second` that the statement never makes.
func TestAnalyzeKeywordArgRoundTrip(t *testing.T) {
	c := New()
	for _, ddl := range []string{
		"CREATE DATABASE kwdb",
		"USE kwdb",
		"CREATE TABLE t (second INT, a DATETIME, b DATETIME)",
	} {
		results, err := c.Exec(ddl, nil)
		if err != nil {
			t.Fatalf("exec %q: %v", ddl, err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Fatalf("exec %q: %v", ddl, r.Error)
			}
		}
	}

	for _, tc := range []struct {
		sql  string
		bare string
	}{
		{"SELECT TIMESTAMPDIFF(SECOND, a, b) FROM t", "SECOND"},
		{"SELECT TIMESTAMPADD(MINUTE, 5, a) FROM t", "MINUTE"},
		{"SELECT GET_FORMAT(DATE, 'USA') FROM t", "DATE"},
	} {
		// roundTrip fatals if either analysis fails, and the re-parse would
		// fail on a quoted unit.
		_, deparsed, _ := roundTrip(t, c, tc.sql)

		if strings.Contains(deparsed, "'"+tc.bare+"'") {
			t.Errorf("%q deparsed the unit quoted: %s", tc.sql, deparsed)
		}
		if !strings.Contains(deparsed, tc.bare) {
			t.Errorf("%q lost the unit: %s", tc.sql, deparsed)
		}
		// The unit is not a read of the same-named column.
		if strings.Contains(deparsed, "`second`") {
			t.Errorf("%q read the column named after the unit: %s", tc.sql, deparsed)
		}
	}
}

// TestAnalyzeKeywordArgViewLineage pins the consequence, not just the shape.
// A view body the analyzer cannot read is stored with AnalyzedQuery nil and no
// error (define.go documents that), so ExpandMergeViews has nothing to expand
// and a query over the view silently loses the view's lineage. Before the
// analyzer learned KeywordArg, that was the standing behaviour for every view
// over TIMESTAMPDIFF, TIMESTAMPADD or GET_FORMAT.
func TestAnalyzeKeywordArgViewLineage(t *testing.T) {
	c := New()
	for _, ddl := range []string{
		"CREATE DATABASE vdb",
		"USE vdb",
		"CREATE TABLE t (second INT, a DATETIME, b DATETIME)",
		"CREATE VIEW v AS SELECT TIMESTAMPDIFF(SECOND, a, b) AS d FROM t",
	} {
		results, err := c.Exec(ddl, nil)
		if err != nil {
			t.Fatalf("exec %q: %v", ddl, err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Fatalf("exec %q: %v", ddl, r.Error)
			}
		}
	}

	view := c.GetDatabase("vdb").Views["v"]
	if view == nil {
		t.Fatal("view v was not created")
	}
	if view.AnalyzedQuery == nil {
		t.Fatal("view v has no AnalyzedQuery, so its lineage cannot be expanded")
	}

	sel := parseSelectForDeparse(t, "SELECT d FROM v")
	q, err := c.AnalyzeSelectStmt(sel)
	if err != nil {
		t.Fatalf("analyze over the view: %v", err)
	}
	expanded := q.ExpandMergeViews(c)
	for _, rte := range expanded.RangeTable {
		if rte.IsView && rte.Subquery == nil {
			t.Errorf("the view body was not expanded into the query, so its lineage is lost")
		}
	}
}
