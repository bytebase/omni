package catalog

import "testing"

// TestQuotedBuiltinCallLineage pins the lineage of the calls the parser reads
// as ordinary ones: a backtick-quoted or schema-qualified builtin name, and a
// qualified TRIM or SUBSTRING, which MariaDB reads with either the builtin's
// syntax or an ordinary argument list. Every argument is still an expression
// over the table, so every output column keeps its source.
func TestQuotedBuiltinCallLineage(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b VARCHAR(20))")

	sel := parseSelect(t, "SELECT `trim`(b, 'x') AS t1, `trim`(b) AS t2, "+
		"`substring`(b, 1, 2) AS s, `group_concat`(b) AS g, `count`(a) AS c, "+
		"testdb.trim(b, a) AS q, testdb.trim('x' FROM b) AS qf, "+
		"testdb.substring(b FROM a) AS qs FROM t")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	a := columnLineage{DB: "testdb", Table: "t", Column: "a"}
	b := columnLineage{DB: "testdb", Table: "t", Column: "b"}
	assertLineage(t, collectColumnLineage(c, q), []resultLineage{
		{Name: "t1", SourceCols: []columnLineage{b}},
		{Name: "t2", SourceCols: []columnLineage{b}},
		{Name: "s", SourceCols: []columnLineage{b}},
		{Name: "g", SourceCols: []columnLineage{b}},
		{Name: "c", SourceCols: []columnLineage{a}},
		{Name: "q", SourceCols: []columnLineage{b, a}},
		{Name: "qf", SourceCols: []columnLineage{b}},
		{Name: "qs", SourceCols: []columnLineage{b, a}},
	})
}
