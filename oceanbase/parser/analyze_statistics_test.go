package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestAnalyzeStatisticsClause pins what ANALYZE TABLE's statistics clause
// parses into. Whether each form parses is held to OceanBase by the
// analyze_stmt and analyze_statistics_clause rows of testdata/compat.tsv.
func TestAnalyzeStatisticsClause(t *testing.T) {
	for _, tc := range []struct{ sql, statistics string }{
		{"ANALYZE TABLE t COMPUTE STATISTICS", "COMPUTE"},
		{"ANALYZE TABLE d.t PARTITION (p0, p1) COMPUTE STATISTICS FOR ALL INDEXED COLUMNS SIZE AUTO", "COMPUTE"},
		{"ANALYZE TABLE t ESTIMATE STATISTICS FOR COLUMNS a SIZE 10, (b, c) SIZE SKEWONLY", "ESTIMATE"},
		{"ANALYZE TABLE t ESTIMATE STATISTICS SAMPLE 5 PERCENTAGE", "ESTIMATE"},
		{"ANALYZE TABLE t ESTIMATE STATISTICS FOR COLUMNS a SAMPLE 5 ROWS", "ESTIMATE"},
		{"ANALYZE TABLE t COMPUTE STATISTICS FOR COLUMNS SIZE 10", "COMPUTE"},
		{"ANALYZE TABLE t COMPUTE STATISTICS FOR COLUMNS size", "COMPUTE"},
		{"ANALYZE TABLE t PARTITION (region = 'us') COMPUTE STATISTICS", "COMPUTE"},
		{"ANALYZE TABLE t ESTIMATE STATISTICS FOR TABLE SAMPLE 100 ROWS", "ESTIMATE"},
		{"ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 8 BUCKETS", ""},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", tc.sql, err)
			continue
		}
		stmt, ok := list.Items[0].(*nodes.AnalyzeTableStmt)
		if len(list.Items) != 1 || !ok || stmt.Statistics != tc.statistics || stmt.Loc.End != len(tc.sql) {
			t.Errorf("Parse(%q) = %#v, want one AnalyzeTableStmt with Statistics %q ending at %d", tc.sql, list.Items[0], tc.statistics, len(tc.sql))
		}
	}
	for _, sql := range []string{
		"ANALYZE TABLE t1, t2 COMPUTE STATISTICS", // the clause names one table
		"ANALYZE TABLE t UPDATE ON a",             // HISTOGRAM is required
		"ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 8",
		"ANALYZE TABLE t DROP HISTOGRAM ON a WITH 8 BUCKETS",
		"ANALYZE TABLE t COMPUTE STATISTICS SAMPLE 5 ROWS", // only ESTIMATE samples
		"ANALYZE TABLE t PARTITION (p0)",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}
