package catalog

import "testing"

func TestExecSkipsDML(t *testing.T) {
	c := New()
	results, err := c.Exec("SELECT 1; INSERT INTO t VALUES (1)", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if !r.Skipped {
			t.Errorf("expected DML to be skipped")
		}
	}
}

func TestExecEmpty(t *testing.T) {
	c := New()
	results, err := c.Exec("", nil)
	if err != nil {
		t.Fatal(err)
	}
	if results != nil {
		t.Errorf("expected nil results for empty SQL")
	}
}

// TestExecSkipsEmptyRows: INSERT and REPLACE with an empty column list or row are
// DML like any other.
func TestExecSkipsEmptyRows(t *testing.T) {
	c := New()
	results, err := c.Exec("INSERT INTO t () VALUES (); REPLACE INTO t VALUES (), (); INSERT INTO t () SELECT 1;", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}
	for _, r := range results {
		if !r.Skipped || r.Error != nil {
			t.Errorf("statement %d: Skipped = %v, Error = %v", r.Index, r.Skipped, r.Error)
		}
	}
}

// TestExecMaterializedView pins that a CREATE MATERIALIZED VIEW installs the
// view its query is, as a snapshot's materialized view does, so a later
// statement can read it; and that the materialized view statements keep their
// own line in the result.
func TestExecMaterializedView(t *testing.T) {
	c := New()
	sql := "CREATE DATABASE test;\nUSE test;\nCREATE TABLE t (a INT);\n" +
		"CREATE MATERIALIZED VIEW mv AS SELECT a FROM t;\n" +
		"CREATE VIEW v AS SELECT a FROM mv;\n" +
		"CREATE MATERIALIZED VIEW LOG ON t;\n" +
		"DROP MATERIALIZED VIEW LOG ON t;"
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec error: %v", err)
	}
	if len(results) != 7 {
		t.Fatalf("got %d results, want 7", len(results))
	}
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("result[%d] = %v", i, r.Error)
		}
		if r.Line != i+1 {
			t.Errorf("result[%d].Line = %d, want %d", i, r.Line, i+1)
		}
	}
	if _, ok := c.GetDatabase("test").Views["mv"]; !ok {
		t.Error("CREATE MATERIALIZED VIEW did not install the view its query is")
	}
}

// TestExecMaterializedViewKind pins that the catalog keeps a materialized view
// apart from an ordinary view while storing both as a view, so one DROP does
// not remove the other.
func TestExecMaterializedViewKind(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE VIEW v AS SELECT 1; CREATE MATERIALIZED VIEW mv AS SELECT 1;", nil); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		sql     string
		wantErr bool
	}{
		{"DROP MATERIALIZED VIEW v;", true},
		{"DROP VIEW mv;", true},
		{"DROP MATERIALIZED VIEW IF EXISTS v;", true},
	} {
		results, err := c.Exec(tc.sql, nil)
		if err != nil {
			t.Fatalf("%q: %v", tc.sql, err)
		}
		if got := len(results) == 1 && results[0].Error != nil; got != tc.wantErr {
			t.Errorf("%q: error = %v, want error %v", tc.sql, results[0].Error, tc.wantErr)
		}
	}
	results, err := c.Exec("DROP MATERIALIZED VIEW mv; DROP VIEW v;", nil)
	if err != nil {
		t.Fatal(err)
	}
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("result[%d] = %v", i, r.Error)
		}
	}
}

// TestExecAlterMaterializedViewTarget pins that ALTER MATERIALIZED VIEW
// requires a materialized view: a missing name and an ordinary view are
// rejected, and a materialized view's option change is accepted.
func TestExecAlterMaterializedViewTarget(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE VIEW v AS SELECT 1; CREATE MATERIALIZED VIEW mv AS SELECT 1;", nil); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		sql     string
		wantErr bool
	}{
		{"ALTER MATERIALIZED VIEW missing NOPARALLEL;", true},
		{"ALTER MATERIALIZED VIEW v NOPARALLEL;", true},
		{"ALTER MATERIALIZED VIEW mv NOPARALLEL;", false},
		{"ALTER MATERIALIZED VIEW mv REFRESH FAST;", false},
	} {
		results, err := c.Exec(tc.sql, nil)
		if err != nil {
			t.Fatalf("%q: %v", tc.sql, err)
		}
		if got := len(results) == 1 && results[0].Error != nil; got != tc.wantErr {
			t.Errorf("%q: error = %v, want error %v", tc.sql, results[0].Error, tc.wantErr)
		}
	}
}

// TestExecMaterializedViewLogTarget pins that a materialized view log names an
// existing base table, not a missing name or a view.
func TestExecMaterializedViewLogTarget(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE TABLE t (a INT); CREATE VIEW v AS SELECT a FROM t;", nil); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		sql     string
		wantErr bool
	}{
		{"CREATE MATERIALIZED VIEW LOG ON missing;", true},
		{"CREATE MATERIALIZED VIEW LOG ON v;", true},
		{"CREATE MATERIALIZED VIEW LOG ON t;", false},
		{"ALTER MATERIALIZED VIEW LOG ON missing NOPARALLEL;", true},
		{"ALTER MATERIALIZED VIEW LOG ON v NOPARALLEL;", true},
		{"ALTER MATERIALIZED VIEW LOG ON t NOPARALLEL;", false},
		{"DROP MATERIALIZED VIEW LOG ON missing;", true},
		{"DROP MATERIALIZED VIEW LOG ON v;", true},
		{"DROP MATERIALIZED VIEW LOG ON t;", false},
	} {
		results, err := c.Exec(tc.sql, nil)
		if err != nil {
			t.Fatalf("%q: %v", tc.sql, err)
		}
		if got := len(results) == 1 && results[0].Error != nil; got != tc.wantErr {
			t.Errorf("%q: error = %v, want error %v", tc.sql, results[0].Error, tc.wantErr)
		}
	}
}
