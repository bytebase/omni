package parser

import "testing"

// The forms ANALYZE TABLE takes on this engine, and the forms it does not.
//
// TiDB's ANALYZE is far wider than this one -- a partition list, an index
// list, ALL COLUMNS, PREDICATE COLUMNS, a COLUMNS list, INCREMENTAL, and a
// WITH list of tuning options -- and none of it is here, which this file
// records: tidb/parser owns that grammar and this engine's server answers
// 1064 for every form of it. DROP STATS, LOAD STATS and CREATE and DROP
// STATISTICS are TiDB's alone for the same reason.
//
// The one case this file was written for is `ANALYZE t`: TABLE is not
// optional, and omni used to match it leniently, so that parsed here as an
// ANALYZE of t where the server answers 1064.
//
// Every case below was measured against this package's pinned server through
// PREPARE on 2026-11-04. TestAnalyzeFormsOracle is what keeps it a
// recording: an image upgrade that changes an answer fails there rather than
// leaving a stale expectation behind.

// analyzeFormCase is a statement and whether the engine's parser accepts it.
type analyzeFormCase struct {
	sql    string
	accept bool
}

var analyzeFormCases = []analyzeFormCase{
	// TABLE is required.
	{"ANALYZE TABLE t", true},
	{"ANALYZE TABLE t1, t2", true},
	{"ANALYZE LOCAL TABLE t", true},
	{"ANALYZE NO_WRITE_TO_BINLOG TABLE t", true},
	{"ANALYZE t", false},
	{"ANALYZE LOCAL t", false},
	{"ANALYZE NO_WRITE_TO_BINLOG t", false},
	{"ANALYZE TABLE", false},
	{"ANALYZE", false},
	// None of TiDB's ANALYZE grammar is this engine's.
	{"ANALYZE TABLE t WITH 4 BUCKETS", false},
	{"ANALYZE TABLE t WITH 4 TOPN", false},
	{"ANALYZE TABLE t WITH 0.5 SAMPLERATE", false},
	{"ANALYZE TABLE t INDEX a", false},
	{"ANALYZE TABLE t INDEX", false},
	{"ANALYZE TABLE t ALL COLUMNS", false},
	{"ANALYZE TABLE t PREDICATE COLUMNS", false},
	{"ANALYZE TABLE t COLUMNS c1, c2", false},
	{"ANALYZE TABLE t PARTITION p0", false},
	{"ANALYZE INCREMENTAL TABLE t INDEX a", false},
	// Nor are TiDB's other statistics statements.
	{"DROP STATS t", false},
	{"LOAD STATS '/tmp/stats.json'", false},
	{"CREATE STATISTICS s1 (CARDINALITY) ON t(a)", false},
	{"DROP STATISTICS s1", false},
	// MySQL 8.0's histogram arms, which this engine has. TOPN is no unit of
	// MySQL's WITH, which takes BUCKETS alone.
	{"ANALYZE TABLE t UPDATE HISTOGRAM ON b", true},
	{"ANALYZE TABLE t UPDATE HISTOGRAM ON b WITH 4 BUCKETS", true},
	{"ANALYZE TABLE t UPDATE HISTOGRAM ON b USING DATA '{}'", true},
	{"ANALYZE TABLE t DROP HISTOGRAM ON b", true},
	{"ANALYZE TABLE t UPDATE HISTOGRAM ON b WITH 4 TOPN", false},
}

// TestAnalyzeForms holds this parser to analyzeFormCases. It needs no
// container, so it guards the grammar wherever it runs.
func TestAnalyzeForms(t *testing.T) {
	for _, c := range analyzeFormCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if accepts := err == nil; accepts != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", accepts, c.accept, c.sql, err)
			}
		})
	}
}
