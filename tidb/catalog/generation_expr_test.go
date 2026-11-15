package catalog

import "testing"

// generationExpressionsTable has a generated column of each expression form
// TiDB restores its own way: parentheses only as written, lower-case
// keywords, string literals with their character set, a quote doubled and a
// backslash escaped, -> as json_extract, and MOD as %.
const generationExpressionsTable = `CREATE TABLE gen_exprs (a int, b int, s varchar(20), j json,
 g1 int AS (a * (b + 1)),
 g2 int AS ((a + b) * 2),
 g3 int AS (a + b * 2),
 g4 varchar(30) AS (concat('it''s', s)),
 g5 varchar(30) AS (CASE WHEN a < 5 THEN 'x' WHEN a IS NULL THEN 'n' ELSE 'y' END) STORED,
 g6 int AS (a IN (1, 2, 3)),
 g7 int AS (s LIKE 'a%'),
 g8 int AS (a BETWEEN 1 AND 5),
 g9 char(10) AS (CAST(a AS char)),
 g10 int AS (IF(a > 0, a, -a)),
 g11 varchar(20) AS (concat('a\\b', s)),
 g12 int AS (NOT (a > 1)),
 g13 int AS (a IS NOT NULL),
 g14 varchar(40) AS (j->>'$.x'),
 g15 int AS (-a + 1),
 g16 int AS (a DIV 2 MOD 3),
 g17 int AS (coalesce(a, b, 0)));`

func TestContainerGenerationExpressionsMatchTiDB(t *testing.T) {
	ctr, done := startContainer(t)
	defer done()
	for _, stmt := range []string{"DROP DATABASE IF EXISTS omni_genexpr", "CREATE DATABASE omni_genexpr", "USE omni_genexpr", generationExpressionsTable} {
		if _, err := ctr.db.ExecContext(ctr.ctx, stmt); err != nil {
			t.Fatalf("TiDB: %q: %v", stmt, err)
		}
	}
	defer func() { _, _ = ctr.db.ExecContext(ctr.ctx, "USE test") }()
	want := showCreateOnTiDB(t, ctr, "table:gen_exprs")
	c := New()
	results, err := c.Exec("CREATE DATABASE omni_genexpr; USE omni_genexpr; "+generationExpressionsTable, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatal(r.Error)
		}
	}
	if got := c.ShowCreateTable("omni_genexpr", "gen_exprs"); got != want {
		t.Errorf("catalog\n%s\nTiDB\n%s", got, want)
	}
}
