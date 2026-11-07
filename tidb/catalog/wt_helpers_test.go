package catalog

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/tidb/parser"
)

// wtSetup creates a Catalog with database "testdb" selected, ready for walk-through tests.
func wtSetup(t *testing.T) *Catalog {
	t.Helper()
	c := New()
	results, err := c.Exec("CREATE DATABASE testdb; USE testdb;", nil)
	if err != nil {
		t.Fatalf("wtSetup parse error: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("wtSetup exec error: %v", r.Error)
		}
	}
	return c
}

// wtExec executes SQL on the catalog and fatals on any error.
func wtExec(t *testing.T, c *Catalog, sql string) {
	t.Helper()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("wtExec parse error: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("wtExec exec error on stmt %d: %v", r.Index, r.Error)
		}
	}
}

// assertRejectedDDL preserves the complete modeled schema on native-rejected DDL.
func assertRejectedDDL(t *testing.T, c *Catalog, sql string, code int) {
	t.Helper()
	before := make(map[string]string)
	for _, db := range c.Databases() {
		for _, tbl := range db.Tables {
			before[db.Name+"."+tbl.Name] = c.ShowCreateTable(db.Name, tbl.Name)
		}
	}
	err := execExpectError(t, c, sql)
	var syntax *parser.ParseError
	if code != 1064 || !errors.As(err, &syntax) {
		assertError(t, err, code)
	}
	after := make(map[string]string)
	for _, db := range c.Databases() {
		for _, tbl := range db.Tables {
			after[db.Name+"."+tbl.Name] = c.ShowCreateTable(db.Name, tbl.Name)
		}
	}
	if len(before) != len(after) {
		t.Fatalf("rejected DDL changed table count: before %v, after %v", before, after)
	}
	for table, ddl := range before {
		if after[table] != ddl {
			t.Errorf("rejected DDL changed %s: before %s; after %s", table, ddl, after[table])
		}
	}
}

// Snapshots can still contain procedure metadata. This fixture deliberately uses
// the metadata API after proving the corresponding SQL execution is unsupported.
func installProcedureMetadataFixture(t *testing.T, c *Catalog, sql string) {
	t.Helper()
	assertRejectedDDL(t, c, sql, 8108)
	if err := c.DefineProcedure(mustParseRoutine(t, sql)); err != nil {
		t.Fatal(err)
	}
}

// assertError asserts that err is a *catalog.Error with the given MySQL error code.
func assertError(t *testing.T, err error, code int) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %d, got nil", code)
	}
	catErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *catalog.Error, got %T: %v", err, err)
	}
	if catErr.Code != code {
		t.Errorf("expected error code %d, got %d: %s", code, catErr.Code, catErr.Message)
	}
}

// assertNoError asserts that err is nil.
func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

// assertShowCreate is SHOW CREATE TABLE compared whole, for text read from the
// server.
func assertShowCreate(t *testing.T, c *Catalog, table, want string) {
	t.Helper()
	if got := c.ShowCreateTable("testdb", table); got != want {
		t.Errorf("SHOW CREATE TABLE %s:\ngot:\n%s\nwant:\n%s", table, got, want)
	}
}
