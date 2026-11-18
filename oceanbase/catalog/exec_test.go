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
	if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE VIEW v AS SELECT 1; CREATE VIEW v2 AS SELECT 2; CREATE MATERIALIZED VIEW mv AS SELECT 1;", nil); err != nil {
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
	// A wrong kind later in a multi-name drop leaves the earlier view in place.
	results, err := c.Exec("DROP VIEW v2, mv;", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Error == nil {
		t.Fatalf("DROP VIEW v2, mv: error = %v, want a wrong kind", results)
	}
	if c.GetDatabase("test").Views["v2"] == nil {
		t.Error("DROP VIEW v2, mv removed v2 before rejecting mv")
	}
	results, err = c.Exec("DROP MATERIALIZED VIEW mv; DROP VIEW v; DROP VIEW v2;", nil)
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

// TestExecMaterializedViewLogState pins that a log is created once, altered
// and dropped only while it exists, and created again after a drop.
func TestExecMaterializedViewLogState(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE TABLE t (a INT);", nil); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		sql     string
		wantErr bool
	}{
		{"CREATE MATERIALIZED VIEW LOG ON t;", false},
		{"CREATE TABLE IF NOT EXISTS t (a INT);", false},
		{"CREATE MATERIALIZED VIEW LOG ON t;", true},
		{"CREATE OR REPLACE MATERIALIZED VIEW LOG ON t;", false},
		{"ALTER MATERIALIZED VIEW LOG ON t NOPARALLEL;", false},
		{"DROP MATERIALIZED VIEW LOG ON t;", false},
		{"ALTER MATERIALIZED VIEW LOG ON t NOPARALLEL;", true},
		{"DROP MATERIALIZED VIEW LOG ON t;", true},
		{"CREATE MATERIALIZED VIEW LOG ON t;", false},
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

// TestExecRenameMaterializedView pins that RENAME TABLE moves the view a
// materialized view installs, keeping its kind and its log.
func TestExecRenameMaterializedView(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE TABLE t (a INT); CREATE MATERIALIZED VIEW mv AS SELECT a FROM t; CREATE MATERIALIZED VIEW LOG ON t; RENAME TABLE mv TO mv2;", nil); err != nil {
		t.Fatal(err)
	}
	db := c.GetDatabase("test")
	if db.Views["mv"] != nil {
		t.Error("RENAME TABLE left the old materialized view name")
	}
	if mv := db.Views["mv2"]; mv == nil || !mv.Materialized {
		t.Errorf("RENAME TABLE did not move the materialized view: %+v", mv)
	}
	for _, tc := range []struct {
		sql     string
		wantErr bool
	}{
		{"ALTER MATERIALIZED VIEW mv2 NOPARALLEL;", false},
		{"DROP VIEW mv2;", true},
		{"DROP MATERIALIZED VIEW mv2;", false},
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

// TestExecAlterTableRenameLog pins that ALTER TABLE ... RENAME moves the
// table's materialized view log state.
func TestExecAlterTableRenameLog(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE TABLE t (a INT); CREATE MATERIALIZED VIEW LOG ON t; ALTER TABLE t RENAME TO u;", nil); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		sql     string
		wantErr bool
	}{
		{"CREATE MATERIALIZED VIEW LOG ON t;", true},
		{"CREATE MATERIALIZED VIEW LOG ON u;", true},
		{"ALTER MATERIALIZED VIEW LOG ON u NOPARALLEL;", false},
		{"DROP MATERIALIZED VIEW LOG ON u;", false},
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

// TestExecRenameOntoMaterializedView pins that a table cannot be renamed onto
// a cataloged materialized view.
func TestExecRenameOntoMaterializedView(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE TABLE t (a INT); CREATE MATERIALIZED VIEW mv AS SELECT 1 AS a;", nil); err != nil {
		t.Fatal(err)
	}
	for _, sql := range []string{"RENAME TABLE t TO mv;", "ALTER TABLE t RENAME TO mv;"} {
		results, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatalf("%q: %v", sql, err)
		}
		if len(results) != 1 || results[0].Error == nil {
			t.Errorf("%q: error = %v, want a duplicate name", sql, results)
		}
	}
}
