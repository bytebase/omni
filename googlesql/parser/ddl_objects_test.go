package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit tests for the remaining core-DDL objects (parser-ddl node): CREATE VIEW,
// CREATE/ALTER/DROP INDEX, CREATE/DROP SCHEMA, CREATE DATABASE, DROP TABLE/VIEW.
// Structural-gate assertions; oracle-verified accept/reject in the *_oracle_test.go
// files.

// --- CREATE VIEW ---

func viewOf(t *testing.T, sql string) *ast.CreateViewStmt {
	t.Helper()
	n := parseDDL(t, sql)
	v, ok := n.(*ast.CreateViewStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.CreateViewStmt", sql, n)
	}
	return v
}

func TestCreateView_SpannerSQLSecurity(t *testing.T) {
	v := viewOf(t, "CREATE VIEW SingerView SQL SECURITY INVOKER AS SELECT SingerId, FirstName FROM Singers")
	if v.Name.String() != "SingerView" {
		t.Errorf("Name = %q, want SingerView", v.Name.String())
	}
	if v.SQLSecurity != "INVOKER" {
		t.Errorf("SQLSecurity = %q, want INVOKER", v.SQLSecurity)
	}
	if v.AsQuery == nil {
		t.Error("AsQuery = nil, want the SELECT body")
	}
}

func TestCreateView_OrReplaceDefiner(t *testing.T) {
	// SQL SECURITY DEFINER is rejected by the Spanner emulator (INVOKER-only) but
	// valid GoogleSQL (legacy .g4 + BigQuery) — divergence; the union parser
	// accepts it.
	v := viewOf(t, "CREATE OR REPLACE VIEW TopAlbums SQL SECURITY DEFINER AS SELECT a FROM Albums WHERE b > 1")
	if !v.OrReplace {
		t.Error("OrReplace = false, want true")
	}
	if v.SQLSecurity != "DEFINER" {
		t.Errorf("SQLSecurity = %q, want DEFINER", v.SQLSecurity)
	}
}

func TestCreateView_BigQueryColumnsAndOptions(t *testing.T) {
	v := viewOf(t, "CREATE OR REPLACE VIEW mydataset.age_groups(age, count) OPTIONS (description = 'd') AS SELECT age, COUNT(*) FROM people GROUP BY age")
	if len(v.Columns) != 2 || v.Columns[0].Name != "age" || v.Columns[1].Name != "count" {
		t.Errorf("Columns = %+v, want [age count]", v.Columns)
	}
	if len(v.Options) != 1 {
		t.Errorf("Options = %d, want 1", len(v.Options))
	}
}

func TestCreateView_Plain(t *testing.T) {
	v := viewOf(t, "CREATE VIEW v AS SELECT 1")
	if v.SQLSecurity != "" {
		t.Errorf("SQLSecurity = %q, want empty", v.SQLSecurity)
	}
	if v.AsQuery == nil {
		t.Error("AsQuery = nil")
	}
}

func TestCreateView_Recursive(t *testing.T) {
	v := viewOf(t, "CREATE RECURSIVE VIEW v AS SELECT 1")
	if !v.Recursive {
		t.Error("Recursive = false, want true")
	}
}

func TestCreateView_Rejects(t *testing.T) {
	cases := []string{
		"CREATE VIEW v",                              // missing AS query
		"CREATE VIEW v SQL SECURITY AS SELECT 1",     // missing INVOKER/DEFINER
		"CREATE VIEW v SQL SECURITY BOGUS AS SELECT 1", // bad security kind
		"CREATE VIEW AS SELECT 1",                    // missing name
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

// --- CREATE INDEX ---

func indexOf(t *testing.T, sql string) *ast.CreateIndexStmt {
	t.Helper()
	n := parseDDL(t, sql)
	i, ok := n.(*ast.CreateIndexStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.CreateIndexStmt", sql, n)
	}
	return i
}

func TestCreateIndex_Basic(t *testing.T) {
	i := indexOf(t, "CREATE INDEX SingersByName ON Singers (LastName, FirstName)")
	if i.Name.String() != "SingersByName" {
		t.Errorf("Name = %q, want SingersByName", i.Name.String())
	}
	if i.Table.String() != "Singers" {
		t.Errorf("Table = %q, want Singers", i.Table.String())
	}
	if len(i.Keys) != 2 {
		t.Fatalf("got %d keys, want 2", len(i.Keys))
	}
	if i.Keys[0].Name != "LastName" || i.Keys[1].Name != "FirstName" {
		t.Errorf("keys = %+v, want LastName, FirstName", i.Keys)
	}
}

func TestCreateIndex_UniqueNullFiltered(t *testing.T) {
	u := indexOf(t, "CREATE UNIQUE INDEX UniqueEmail ON Users (Email)")
	if !u.Unique {
		t.Error("Unique = false, want true")
	}
	nf := indexOf(t, "CREATE NULL_FILTERED INDEX AlbumsByTitle ON Albums (Title)")
	if !nf.NullFiltered {
		t.Error("NullFiltered = false, want true")
	}
}

func TestCreateIndex_Direction(t *testing.T) {
	i := indexOf(t, "CREATE INDEX idx ON T (a DESC, b ASC)")
	if i.Keys[0].Direction != "DESC" || i.Keys[1].Direction != "ASC" {
		t.Errorf("key directions = %q/%q, want DESC/ASC", i.Keys[0].Direction, i.Keys[1].Direction)
	}
}

func TestCreateIndex_Storing(t *testing.T) {
	i := indexOf(t, "CREATE INDEX AlbumsByArtist ON Albums (SingerId, Title) STORING (MarketingBudget)")
	if len(i.Storing) != 1 {
		t.Errorf("Storing = %d exprs, want 1", len(i.Storing))
	}
}

func TestCreateIndex_StoringAndInterleave(t *testing.T) {
	i := indexOf(t, "CREATE INDEX SongsBySinger ON Songs (SingerId, AlbumId) STORING (SongName), INTERLEAVE IN Albums")
	if len(i.Storing) != 1 {
		t.Errorf("Storing = %d, want 1", len(i.Storing))
	}
	if i.Interleave == nil || i.Interleave.String() != "Albums" {
		t.Errorf("Interleave = %+v, want Albums", i.Interleave)
	}
}

func TestCreateIndex_IfNotExistsAndWhere(t *testing.T) {
	i := indexOf(t, "CREATE INDEX IF NOT EXISTS idx ON T (a) WHERE a IS NOT NULL")
	if !i.IfNotExists {
		t.Error("IfNotExists = false, want true")
	}
	if i.Where == nil {
		t.Error("Where = nil, want the filter expression")
	}
}

func TestCreateIndex_PartitionByOptions(t *testing.T) {
	// BigQuery index suffix: PARTITION BY + OPTIONS.
	i := indexOf(t, "CREATE INDEX idx ON T (a) PARTITION BY b OPTIONS (x = 1)")
	if len(i.PartitionBy) != 1 {
		t.Errorf("PartitionBy = %d, want 1", len(i.PartitionBy))
	}
	if len(i.Options) != 1 {
		t.Errorf("Options = %d, want 1", len(i.Options))
	}
}

func TestCreateIndex_EmptyKeyListAccepted(t *testing.T) {
	// oracle-confirmed accept: the Spanner emulator accepts `CREATE INDEX … ()`
	// (an empty key list), even though the legacy .g4 index_order_by_and_options
	// mandates >= 1 element. CREATE INDEX is a Spanner-authoritative form, so the
	// live oracle decides — the parser accepts the empty list (shared key-part
	// list helper). Divergence from the .g4, recorded.
	i := indexOf(t, "CREATE INDEX idx ON T ()")
	if len(i.Keys) != 0 {
		t.Errorf("got %d keys, want 0", len(i.Keys))
	}
}

func TestCreateIndex_Rejects(t *testing.T) {
	cases := []string{
		"CREATE INDEX idx",      // missing ON table + keys
		"CREATE INDEX idx ON T", // missing key list
		"CREATE INDEX ON T (a)", // missing index name
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

// --- DROP INDEX ---

func TestDropIndex(t *testing.T) {
	d := dropOf(t, "DROP INDEX SingersByName")
	if d.Object != ast.DropIndex || d.Name.String() != "SingersByName" {
		t.Errorf("drop = %+v, want DROP INDEX SingersByName", d)
	}
	d2 := dropOf(t, "DROP INDEX IF EXISTS idx")
	if !d2.IfExists {
		t.Error("IfExists = false, want true")
	}
}

// --- CREATE / DROP SCHEMA ---

func TestCreateSchema(t *testing.T) {
	n := parseDDL(t, "CREATE SCHEMA myschema")
	s, ok := n.(*ast.CreateSchemaStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.CreateSchemaStmt", n)
	}
	if s.Name.String() != "myschema" {
		t.Errorf("Name = %q, want myschema", s.Name.String())
	}
}

func TestCreateSchema_IfNotExistsCollateOptions(t *testing.T) {
	n := parseDDL(t, "CREATE SCHEMA IF NOT EXISTS mydataset DEFAULT COLLATE 'und:ci' OPTIONS (location = 'us')")
	s := n.(*ast.CreateSchemaStmt)
	if !s.IfNotExists {
		t.Error("IfNotExists = false, want true")
	}
	if s.DefaultCollate == "" {
		t.Error("DefaultCollate empty, want und:ci")
	}
	if len(s.Options) != 1 {
		t.Errorf("Options = %d, want 1", len(s.Options))
	}
}

func TestDropSchema(t *testing.T) {
	d := dropOf(t, "DROP SCHEMA myschema")
	if d.Object != ast.DropSchema {
		t.Errorf("Object = %v, want SCHEMA", d.Object)
	}
	d2 := dropOf(t, "DROP SCHEMA IF EXISTS mydataset CASCADE")
	if !d2.IfExists || d2.DropMode != "CASCADE" {
		t.Errorf("drop = %+v, want IF EXISTS … CASCADE", d2)
	}
}

func TestDropExternalSchema(t *testing.T) {
	d := dropOf(t, "DROP EXTERNAL SCHEMA mydataset")
	if d.Object != ast.DropSchema || !d.External {
		t.Errorf("drop = %+v, want DROP EXTERNAL SCHEMA", d)
	}
}

// --- CREATE DATABASE ---

func TestCreateDatabase(t *testing.T) {
	n := parseDDL(t, "CREATE DATABASE my_database")
	db, ok := n.(*ast.CreateDatabaseStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.CreateDatabaseStmt", n)
	}
	if db.Name.String() != "my_database" {
		t.Errorf("Name = %q, want my_database", db.Name.String())
	}
}

func TestCreateDatabase_Options(t *testing.T) {
	n := parseDDL(t, "CREATE DATABASE db OPTIONS (x = 1)")
	db := n.(*ast.CreateDatabaseStmt)
	if len(db.Options) != 1 {
		t.Errorf("Options = %d, want 1", len(db.Options))
	}
}

func TestCreateDatabase_NoOrReplace(t *testing.T) {
	// CREATE DATABASE has no OR REPLACE production.
	assertReject(t, "CREATE OR REPLACE DATABASE db")
}

// --- DROP TABLE / VIEW ---

func dropOf(t *testing.T, sql string) *ast.DropStmt {
	t.Helper()
	n := parseDDL(t, sql)
	d, ok := n.(*ast.DropStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.DropStmt", sql, n)
	}
	return d
}

func TestDropTable(t *testing.T) {
	d := dropOf(t, "DROP TABLE Singers")
	if d.Object != ast.DropTable || d.Name.String() != "Singers" {
		t.Errorf("drop = %+v, want DROP TABLE Singers", d)
	}
	d2 := dropOf(t, "DROP TABLE IF EXISTS temp_table")
	if !d2.IfExists {
		t.Error("IfExists = false, want true")
	}
}

func TestDropTable_DashedPath(t *testing.T) {
	d := dropOf(t, "DROP TABLE `my-project`.dataset.tbl")
	if d.Name.String() != "my-project.dataset.tbl" {
		t.Errorf("Name = %q, want my-project.dataset.tbl", d.Name.String())
	}
}

func TestDropView(t *testing.T) {
	d := dropOf(t, "DROP VIEW IF EXISTS mydataset.myview")
	if d.Object != ast.DropView || !d.IfExists {
		t.Errorf("drop = %+v, want DROP VIEW IF EXISTS", d)
	}
}

func TestDrop_Rejects(t *testing.T) {
	cases := []string{
		"DROP TABLE",      // missing name
		"DROP INDEX",      // missing name
		"DROP SCHEMA",     // missing name
		"DROP TABLE T junk", // trailing junk
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

// --- dialect-node object kinds route to the unsupported stub (not unknown) ---

func TestCreateAndDrop_DialectObjectsStubbed(t *testing.T) {
	stubbed := []string{
		"CREATE FUNCTION f() RETURNS INT64 AS (1)",
		"CREATE PROCEDURE p() BEGIN SELECT 1; END",
		"CREATE MATERIALIZED VIEW mv AS SELECT 1",
		"CREATE SEARCH INDEX si ON T (c)",
		"CREATE VECTOR INDEX vi ON T (c) OPTIONS (x = 1)",
		"CREATE TABLE FUNCTION tf() AS SELECT 1",
		"CREATE MODEL m OPTIONS (x = 1)",
		"DROP FUNCTION f",
		"DROP PROCEDURE p",
		"DROP MATERIALIZED VIEW mv",
		"DROP TABLE FUNCTION tf",
		"ALTER MATERIALIZED VIEW mv SET OPTIONS (x = 1)",
	}
	for _, sql := range stubbed {
		_, errs := Parse(sql)
		if len(errs) == 0 {
			// Some of these (e.g. an unsupported body that happens to be a no-op)
			// could parse; the load-bearing assertion is that none hits the UNKNOWN
			// branch, checked below.
			continue
		}
		for _, e := range errs {
			if contains(e.Msg, "unknown or unsupported statement") {
				t.Errorf("Parse(%q): hit UNKNOWN branch: %q (should be a 'not yet supported' stub)", sql, e.Msg)
			}
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || indexOfSub(s, sub) >= 0)
}

func indexOfSub(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
