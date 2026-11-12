package catalog

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
	gomysql "github.com/go-sql-driver/mysql"
)

const asOfViewBody = "SELECT a FROM tbl AS OF TIMESTAMP NOW() - INTERVAL 1 SECOND WHERE a > 1"

// omniAsOfView stores the view in a catalog and returns what SHOW CREATE VIEW writes.
func omniAsOfView(t *testing.T) string {
	t.Helper()
	c := New()
	results, err := c.Exec("CREATE DATABASE d; USE d; CREATE TABLE tbl (a INT, b INT, PRIMARY KEY (a)); CREATE VIEW v AS "+asOfViewBody, nil)
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

// TestAsOfView checks a stored view keeps the clause. Without it the view would read
// the current data.
func TestAsOfView(t *testing.T) {
	got := omniAsOfView(t)
	if want := "`d`.`tbl` AS OF TIMESTAMP DATE_SUB(NOW(), INTERVAL 1 SECOND)"; !strings.Contains(got, want) {
		t.Errorf("SHOW CREATE VIEW lost %q: %s", want, got)
	}
}

// TestAsOfAnalyzedQuery checks the analyzed query keeps the clause on every entry that
// can hold one. A query deparsed without it still reads, and reads current data. The
// names in the last case resolve to nothing, and the catalog does not check them:
// TestAsOfNotInQueryScopeTiDB records what the server does.
func TestAsOfAnalyzedQuery(t *testing.T) {
	for _, tc := range []struct {
		name, sql string
		want      []string
	}{
		{"table", "SELECT a FROM tbl AS OF TIMESTAMP '2020-01-01'", []string{"`d`.`tbl` as of timestamp '2020-01-01'"}},
		{"alias", "SELECT x.a FROM tbl AS x AS OF TIMESTAMP NOW()", []string{"`d`.`tbl` as `x` as of timestamp now()"}},
		{"bare alias", "SELECT x.a FROM tbl x AS OF TIMESTAMP NOW()", []string{"`d`.`tbl` as `x` as of timestamp now()"}},
		{"both sides of a join", "SELECT x.a FROM tbl x AS OF TIMESTAMP 1 JOIN u y AS OF TIMESTAMP 2 ON x.a = y.a", []string{
			"`d`.`tbl` as `x` as of timestamp 1",
			"`d`.`u` as `y` as of timestamp 2",
		}},
		{"comma", "SELECT tbl.a FROM tbl AS OF TIMESTAMP 1, u AS OF TIMESTAMP 2", []string{"`d`.`tbl` as of timestamp 1, `d`.`u` as of timestamp 2"}},
		{"before the sample", "SELECT a FROM tbl AS OF TIMESTAMP 1 TABLESAMPLE REGIONS()", []string{"`d`.`tbl` as of timestamp 1 tablesample regions ()"}},
		{"cte", "WITH c AS (SELECT 1 AS x) SELECT x FROM c AS OF TIMESTAMP 1", []string{"`c` as of timestamp 1"}},
		{"cte alias", "WITH c AS (SELECT 1 AS x) SELECT y.x FROM c AS y AS OF TIMESTAMP 1", []string{"`c` as `y` as of timestamp 1"}},
		{"view", "SELECT a FROM vs AS OF TIMESTAMP 1", []string{"`d`.`vs` as of timestamp 1"}},
		{"unresolved", "SELECT a FROM tbl AS OF TIMESTAMP nosuch", []string{"as of timestamp `nosuch`"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := DeparseQuery(analyzeSampled(t, tc.sql))
			for _, want := range tc.want {
				if !strings.Contains(got, want) {
					t.Errorf("DeparseQuery lost %q: %s", want, got)
				}
			}
		})
	}
}

// TestAsOfAnalyzedView checks the analyzed query a stored view holds keeps the clause,
// the same way its text does, and that expanding the view keeps it too.
func TestAsOfAnalyzedView(t *testing.T) {
	c := New()
	results, err := c.Exec("CREATE DATABASE d; USE d; CREATE TABLE tbl (a INT, b INT, PRIMARY KEY (a)); CREATE VIEW v AS "+asOfViewBody, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatal(r.Error)
		}
	}
	view := c.GetDatabase("d").Views["v"]
	if view == nil || view.AnalyzedQuery == nil {
		t.Fatalf("no analyzed query on the view: %+v", view)
	}
	const want = "as of timestamp (now() - interval 1 second)"
	if got := DeparseQuery(view.AnalyzedQuery); !strings.Contains(got, want) {
		t.Errorf("AnalyzedQuery lost %q: %s", want, got)
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

// TestAsOfNotInQueryScopeTiDB records why RangeTableEntryQ.AsOf is text: TiDB evaluates
// the expression with no table in scope, so a column of the table it reads is unknown
// there. Each of these answers 1054 and parses, which the catalog does not check.
func TestAsOfNotInQueryScopeTiDB(t *testing.T) {
	ctr := startTiDBForCatalog(t)
	conn, err := ctr.db.Conn(ctr.ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for _, stmt := range []string{
		"DROP DATABASE IF EXISTS omni_as_of_scope",
		"CREATE DATABASE omni_as_of_scope",
		"USE omni_as_of_scope",
		"CREATE TABLE tbl (a INT, b INT, PRIMARY KEY (a))",
	} {
		if _, err := conn.ExecContext(ctr.ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	defer func() {
		_, _ = conn.ExecContext(ctr.ctx, "DROP DATABASE IF EXISTS omni_as_of_scope")
	}()
	for _, sql := range []string{
		"SELECT a FROM tbl AS OF TIMESTAMP a",
		"SELECT a FROM tbl AS OF TIMESTAMP tbl.a",
		"SELECT a FROM tbl AS OF TIMESTAMP nosuch",
	} {
		t.Run(sql, func(t *testing.T) {
			_, err := conn.ExecContext(ctr.ctx, sql)
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) || myErr.Number != 1054 {
				t.Errorf("err = %v, want error 1054; RangeTableEntryQ.AsOf may need analyzing", err)
			}
			if got := DeparseQuery(analyzeSampled(t, sql)); !strings.Contains(got, "as of timestamp") {
				t.Errorf("no clause in %s", got)
			}
		})
	}
}

// TestAsOfExpressionsReadBack sends every kind of expression through the analyzed query
// and holds the source and the result to a live TiDB, so a kind DeparseQuery cannot
// write fails here. A planner error still counts as parsed, and only a 1064 fails.
func TestAsOfExpressionsReadBack(t *testing.T) {
	ctr := startTiDBForCatalog(t)
	conn, err := ctr.db.Conn(ctr.ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	parses := func(t *testing.T, sql string) bool {
		t.Helper()
		return prepareParses(t, ctr.ctx, conn, sql)
	}
	for _, expr := range tableSampleExpressions {
		t.Run(expr, func(t *testing.T) {
			sql := "SELECT a FROM tbl AS OF TIMESTAMP " + expr
			if !parses(t, sql) {
				t.Fatalf("TiDB refused the source text %s", sql)
			}
			analyzed := DeparseQuery(analyzeSampled(t, sql))
			if strings.Contains(analyzed, "unsupported") {
				t.Fatalf("DeparseQuery could not write it: %s", analyzed)
			}
			if !strings.Contains(analyzed, "as of timestamp") {
				t.Fatalf("no clause in %s", analyzed)
			}
			if !parses(t, analyzed) {
				t.Errorf("TiDB refused the analyzed text %s", analyzed)
			}
			list, err := parser.Parse(analyzed)
			if err != nil {
				t.Fatalf("omni refused its own text %s: %v", analyzed, err)
			}
			if again := DeparseQuery(analyzeSampledStmt(t, list.Items[0].(*ast.SelectStmt))); again != analyzed {
				t.Errorf("reads back as %s, want %s", again, analyzed)
			}
		})
	}
}

// TestAsOfViewTiDBDiffers records what TiDB does with this view and where the catalog
// differs. TiDB v8.5.5 parses the statement and answers 1356 to it, while the catalog
// stores it, because the catalog does not run the planner. The same SELECT runs on its
// own, which tells that answer from a table too young for the timestamp, and the text
// the catalog writes runs on the server.
func TestAsOfViewTiDBDiffers(t *testing.T) {
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
	exec("DROP DATABASE IF EXISTS omni_as_of_view")
	exec("CREATE DATABASE omni_as_of_view")
	defer func() {
		_, _ = conn.ExecContext(ctr.ctx, "DROP DATABASE IF EXISTS omni_as_of_view")
	}()
	exec("USE omni_as_of_view")
	exec("CREATE TABLE tbl (a INT, b INT, PRIMARY KEY (a))")
	// The timestamp has to be later than the table.
	exec("DO SLEEP(3)")
	exec(asOfViewBody)

	_, err = conn.ExecContext(ctr.ctx, "CREATE VIEW v AS "+asOfViewBody)
	var myErr *gomysql.MySQLError
	if !errors.As(err, &myErr) || myErr.Number != 1356 {
		t.Errorf("CREATE VIEW: err = %v, want error 1356; the catalog may need to match it", err)
	}

	omni := omniAsOfView(t)
	_, body, ok := strings.Cut(omni, " AS ")
	if !ok {
		t.Fatalf("no body in %s", omni)
	}
	if !strings.Contains(body, "AS OF TIMESTAMP") {
		t.Fatalf("the catalog's text has no clause: %s", body)
	}
	body = strings.ReplaceAll(body, "`d`.", "`omni_as_of_view`.")
	if _, err := conn.ExecContext(ctr.ctx, body); err != nil {
		t.Errorf("TiDB refused the catalog's text %s: %v", body, err)
	}
}
