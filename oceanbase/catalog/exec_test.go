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
