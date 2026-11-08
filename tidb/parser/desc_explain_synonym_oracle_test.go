package parser

import "testing"

// TestDescExplainSynonym holds omni to live TiDB v8.5.5 on DESC and DESCRIBE as
// full EXPLAIN synonyms (the §3.4 integrationtest corpus's largest gap, 211
// rows). EXPLAIN, DESCRIBE and DESC are the one ExplainSym in parser.y, so every
// form is measured against the server for all three keywords and omni must agree
// -- the explainable statements that were the gap, the table form that must stay
// valid, and the modifier combinations the server refuses.
func TestDescExplainSynonym(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	tails := []string{
		"SELECT 1",
		"FORMAT = 'brief' SELECT max(a) FROM t",
		"ANALYZE SELECT * FROM t",
		"UPDATE t SET a = 1",
		"DELETE FROM t",
		"INSERT INTO t VALUES (1)",
		"REPLACE INTO t VALUES (1)",
		"TABLE t",
		"WITH cte AS (SELECT 1) SELECT * FROM cte",
		"ALTER TABLE t ADD COLUMN c int",
		"FOR CONNECTION 1",
		"t",                      // table form
		"t c",                    // table + column
		"ANALYZE t",              // modifier rules out the table form
		"FORMAT = 'row' t",       // modifier rules out the table form
		"CREATE TABLE x (a int)", // not an explainable statement
	}
	for _, tail := range tails {
		for _, kw := range []string{"EXPLAIN", "DESC", "DESCRIBE"} {
			sql := kw + " " + tail
			tidbAccepts, _, err := tc.syntaxVerdict(t, sql, "<desc synonym test>")
			if err != nil {
				t.Fatalf("no verdict for %q: %v", sql, err)
			}
			_, perr := Parse(sql)
			if omniAccepts := perr == nil; omniAccepts != tidbAccepts {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB accepts=%v for %q: %v",
					omniAccepts, tidbAccepts, sql, perr)
			}
		}
	}
}
