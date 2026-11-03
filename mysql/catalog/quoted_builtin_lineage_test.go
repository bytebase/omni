package catalog

import "testing"

// TestQuotedBuiltinCallLineage pins the lineage of the calls the parser reads
// as ordinary ones: a backtick-quoted or schema-qualified builtin name. Their
// special argument syntax is gone, but each argument is still an expression
// over the table, so every output column keeps its source.
func TestQuotedBuiltinCallLineage(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b VARCHAR(20))")

	sel := parseSelect(t, "SELECT `trim`(b, 'x') AS t1, `trim`(b) AS t2, "+
		"`substring`(b, 1, 2) AS s, `group_concat`(b) AS g, `count`(a) AS c, "+
		"testdb.trim(b, a) AS q FROM t")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	b := []columnLineage{{DB: "testdb", Table: "t", Column: "b"}}
	assertLineage(t, collectColumnLineage(c, q), []resultLineage{
		{Name: "t1", SourceCols: b},
		{Name: "t2", SourceCols: b},
		{Name: "s", SourceCols: b},
		{Name: "g", SourceCols: b},
		{Name: "c", SourceCols: []columnLineage{{DB: "testdb", Table: "t", Column: "a"}}},
		{Name: "q", SourceCols: []columnLineage{
			{DB: "testdb", Table: "t", Column: "b"},
			{DB: "testdb", Table: "t", Column: "a"},
		}},
	})
}
