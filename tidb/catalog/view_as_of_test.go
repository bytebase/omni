package catalog

import (
	"errors"
	"strings"
	"testing"

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
	if want := "as of timestamp (now() - interval 1 second)"; !strings.Contains(got, want) {
		t.Errorf("SHOW CREATE VIEW lost %q: %s", want, got)
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
	if !strings.Contains(body, "as of timestamp") {
		t.Fatalf("the catalog's text has no clause: %s", body)
	}
	body = strings.ReplaceAll(body, "`d`.", "`omni_as_of_view`.")
	if _, err := conn.ExecContext(ctr.ctx, body); err != nil {
		t.Errorf("TiDB refused the catalog's text %s: %v", body, err)
	}
}
