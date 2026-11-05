package catalog

import (
	"errors"
	"strings"
	"testing"

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
