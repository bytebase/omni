package catalog

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
	gomysql "github.com/go-sql-driver/mysql"
)

// tableSampleViewSchema and tableSampleViewBody are a view TiDB v8.5.5 creates:
// REGIONS is the one sampling method its planner takes, and any other answers 8128.
const (
	tableSampleViewSchema = "CREATE TABLE tbl (a INT, b INT, PRIMARY KEY (a))"
	tableSampleViewBody   = "SELECT a FROM tbl TABLESAMPLE REGIONS() REPEATABLE (1) WHERE a > 1"
)

// omniTableSampleView stores the view in a catalog and returns what SHOW CREATE VIEW
// writes for it.
func omniTableSampleView(t *testing.T) string {
	t.Helper()
	c := New()
	results, err := c.Exec("CREATE DATABASE d; USE d; "+tableSampleViewSchema+"; CREATE VIEW v AS "+tableSampleViewBody, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatal(r.Error)
		}
	}
	return c.ShowCreateView("d", "v")
}

// TestTableSampleView checks a stored view keeps its sample. Without it the view
// would still read, and would read the whole table.
func TestTableSampleView(t *testing.T) {
	got := omniTableSampleView(t)
	for _, want := range []string{"tablesample regions ()", "repeatable (1)"} {
		if !strings.Contains(got, want) {
			t.Errorf("SHOW CREATE VIEW lost %q: %s", want, got)
		}
	}
}

// analyzeSampled analyzes one SELECT against tbl, u and the view vs.
func analyzeSampled(t *testing.T, sql string) *Query {
	t.Helper()
	c := New()
	results, err := c.Exec("CREATE DATABASE d; USE d; "+tableSampleViewSchema+"; CREATE TABLE u (a INT); CREATE VIEW vs AS SELECT a FROM tbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatal(r.Error)
		}
	}
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	q, err := c.AnalyzeSelectStmt(list.Items[0].(*ast.SelectStmt))
	if err != nil {
		t.Fatalf("analyzing %q: %v", sql, err)
	}
	return q
}

// TestTableSampleAnalyzedQuery checks the analyzed query keeps the clause on every
// entry that can hold one. A query deparsed without it still reads, and reads the
// whole table.
//
// The names in the last cases resolve to nothing, and that has to stay an accepted
// query: TestTableSampleNotResolvedTiDB holds the server to it.
func TestTableSampleAnalyzedQuery(t *testing.T) {
	for _, tc := range []struct {
		name, sql string
		want      []string
	}{
		{"table", "SELECT a FROM tbl TABLESAMPLE REGIONS(50 ROWS) REPEATABLE (1)", []string{"`d`.`tbl` tablesample regions (50 rows) repeatable (1)"}},
		{"bare", "SELECT a FROM tbl TABLESAMPLE ()", []string{"`d`.`tbl` tablesample ()"}},
		{"alias", "SELECT x.a FROM tbl AS x TABLESAMPLE SYSTEM (50 PERCENT)", []string{"`d`.`tbl` as `x` tablesample system (50 percent)"}},
		{"both sides of a join", "SELECT x.a FROM tbl x TABLESAMPLE (10) JOIN u y TABLESAMPLE (20) ON x.a = y.a", []string{
			"`d`.`tbl` as `x` tablesample (10)",
			"`d`.`u` as `y` tablesample (20)",
		}},
		{"comma", "SELECT tbl.a FROM tbl TABLESAMPLE (10), u TABLESAMPLE (20)", []string{"`d`.`tbl` tablesample (10), `d`.`u` tablesample (20)"}},
		{"cte", "WITH c AS (SELECT 1 AS x) SELECT x FROM c TABLESAMPLE REGIONS()", []string{"`c` tablesample regions ()"}},
		{"cte alias", "WITH c AS (SELECT 1 AS x) SELECT y.x FROM c AS y TABLESAMPLE REGIONS()", []string{"`c` as `y` tablesample regions ()"}},
		{"view", "SELECT a FROM vs TABLESAMPLE REGIONS()", []string{"`d`.`vs` tablesample regions ()"}},
		{"unresolved size", "SELECT a FROM tbl TABLESAMPLE REGIONS(nosuch)", []string{"tablesample regions (`nosuch`)"}},
		{"unresolved seed", "SELECT a FROM tbl TABLESAMPLE REGIONS() REPEATABLE (nosuchfunc(1))", []string{"repeatable (nosuchfunc(1))"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := DeparseQuery(analyzeSampled(t, tc.sql))
			for _, want := range tc.want {
				if !strings.Contains(got, want) {
					t.Errorf("DeparseQuery = %s, want it to hold %s", got, want)
				}
			}
		})
	}

	t.Run("no clause", func(t *testing.T) {
		q := analyzeSampled(t, "SELECT a FROM tbl")
		if s := q.RangeTable[0].Sample; s != "" {
			t.Errorf("Sample = %q, want none", s)
		}
		if got := DeparseQuery(q); strings.Contains(got, "tablesample") {
			t.Errorf("DeparseQuery = %s, wrote a clause the query does not have", got)
		}
	})

	t.Run("text on the entry", func(t *testing.T) {
		q := analyzeSampled(t, "SELECT a FROM tbl TABLESAMPLE BERNOULLI (5 PERCENT)")
		if got, want := q.RangeTable[0].Sample, "tablesample bernoulli (5 percent)"; got != want {
			t.Errorf("Sample = %q, want %q", got, want)
		}
	})
}

// TestTableSampleAnalyzedView checks the view's analyzed body keeps the clause, as
// its Definition does, and that a query over the view reaches it through expansion.
func TestTableSampleAnalyzedView(t *testing.T) {
	c := New()
	results, err := c.Exec("CREATE DATABASE d; USE d; "+tableSampleViewSchema+"; CREATE VIEW v AS "+tableSampleViewBody, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatal(r.Error)
		}
	}
	const want = "`d`.`tbl` tablesample regions () repeatable (1)"
	view := c.GetDatabase("d").Views["v"]
	if view.AnalyzedQuery == nil {
		t.Fatal("no analyzed body")
	}
	if got := DeparseQuery(view.AnalyzedQuery); !strings.Contains(got, want) {
		t.Errorf("analyzed body = %s, want it to hold %s", got, want)
	}

	list, err := parser.Parse("SELECT a FROM v")
	if err != nil {
		t.Fatal(err)
	}
	q, err := c.AnalyzeSelectStmt(list.Items[0].(*ast.SelectStmt))
	if err != nil {
		t.Fatal(err)
	}
	expanded := q.ExpandMergeViews(c).RangeTable[0].Subquery
	if expanded == nil {
		t.Fatal("the view did not expand")
	}
	if got := DeparseQuery(expanded); !strings.Contains(got, want) {
		t.Errorf("expanded body = %s, want it to hold %s", got, want)
	}
}

// TestTableSampleNotResolvedTiDB records why the analyzed query keeps the clause as
// parsed instead of analyzing it: TiDB v8.5.5 runs these, though every name in them
// resolves to nothing. It also holds the catalog's analyzed text to the server.
func TestTableSampleNotResolvedTiDB(t *testing.T) {
	ctr := startTiDBForCatalog(t)
	conn, err := ctr.db.Conn(ctr.ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for _, stmt := range []string{
		"DROP DATABASE IF EXISTS omni_tablesample_resolve",
		"CREATE DATABASE omni_tablesample_resolve",
		"USE omni_tablesample_resolve",
		tableSampleViewSchema,
	} {
		if _, err := conn.ExecContext(ctr.ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	defer func() {
		_, _ = conn.ExecContext(ctr.ctx, "DROP DATABASE IF EXISTS omni_tablesample_resolve")
	}()
	for _, sql := range []string{
		"SELECT a FROM tbl TABLESAMPLE REGIONS(nosuch)",
		"SELECT a FROM tbl TABLESAMPLE REGIONS(nosuch.col)",
		"SELECT a FROM tbl TABLESAMPLE REGIONS(nosuchfunc(1))",
		"SELECT a FROM tbl TABLESAMPLE REGIONS((SELECT nosuch FROM nosuch))",
		"SELECT a FROM tbl TABLESAMPLE REGIONS() REPEATABLE (nosuch)",
		tableSampleViewBody,
	} {
		t.Run(sql, func(t *testing.T) {
			if _, err := conn.ExecContext(ctr.ctx, sql); err != nil {
				t.Errorf("TiDB refused %q: %v", sql, err)
			}
			analyzed := DeparseQuery(analyzeSampled(t, sql))
			if !strings.Contains(analyzed, "tablesample regions") {
				t.Fatalf("no sample in %s", analyzed)
			}
			analyzed = strings.ReplaceAll(analyzed, "`d`.", "`omni_tablesample_resolve`.")
			if _, err := conn.ExecContext(ctr.ctx, analyzed); err != nil {
				t.Errorf("TiDB refused the analyzed text %s: %v", analyzed, err)
			}
		})
	}
}

// TestTableSampleViewTiDBDiffers records what TiDB does with this view, and why the
// catalog does not copy it. TiDB v8.5.5 stores REGIONS as REGION, which its own parser
// refuses, so the view it just created answers 1105 when read. The catalog writes the
// word the grammar has, and TiDB takes that text.
func TestTableSampleViewTiDBDiffers(t *testing.T) {
	ctr := startTiDBForCatalog(t)
	conn, err := ctr.db.Conn(ctr.ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	exec := func(stmt string) {
		t.Helper()
		if _, err := conn.ExecContext(ctr.ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	exec("DROP DATABASE IF EXISTS omni_tablesample_view")
	exec("CREATE DATABASE omni_tablesample_view")
	defer func() {
		_, _ = conn.ExecContext(ctr.ctx, "DROP DATABASE IF EXISTS omni_tablesample_view")
	}()
	exec("USE omni_tablesample_view")
	exec(tableSampleViewSchema)
	exec("CREATE VIEW v AS " + tableSampleViewBody)

	var name, stored, charset, collation string
	if err := conn.QueryRowContext(ctr.ctx, "SHOW CREATE VIEW v").Scan(&name, &stored, &charset, &collation); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stored, "TABLESAMPLE REGION ()") {
		t.Errorf("TiDB no longer stores REGION; the catalog may be able to match it: %s", stored)
	}

	_, err = conn.ExecContext(ctr.ctx, "SELECT * FROM v")
	var myErr *gomysql.MySQLError
	if !errors.As(err, &myErr) || myErr.Number != 1105 {
		t.Errorf("reading the view: err = %v, want error 1105", err)
	}

	omni := omniTableSampleView(t)
	_, body, ok := strings.Cut(omni, " AS ")
	if !ok {
		t.Fatalf("no body in %s", omni)
	}
	if !strings.Contains(body, "tablesample regions ()") {
		t.Fatalf("the catalog's text has no sample: %s", body)
	}
	body = strings.ReplaceAll(body, "`d`.", "`omni_tablesample_view`.")
	if _, err := conn.ExecContext(ctr.ctx, body); err != nil {
		t.Errorf("TiDB refused the catalog's text %s: %v", body, err)
	}
}
