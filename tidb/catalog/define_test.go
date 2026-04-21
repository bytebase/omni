package catalog

import (
	"errors"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// -- helpers -------------------------------------------------------------

// parseFirst parses sql and returns the single top-level statement.
func parseFirst(t *testing.T, sql string) nodes.Node {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v\nSQL: %s", err, sql)
	}
	if list == nil || len(list.Items) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(list.Items))
	}
	return list.Items[0]
}

func mustParseDatabase(t *testing.T, sql string) *nodes.CreateDatabaseStmt {
	t.Helper()
	stmt, ok := parseFirst(t, sql).(*nodes.CreateDatabaseStmt)
	if !ok {
		t.Fatalf("not a CreateDatabaseStmt: %T", parseFirst(t, sql))
	}
	return stmt
}

func mustParseTable(t *testing.T, sql string) *nodes.CreateTableStmt {
	t.Helper()
	stmt, ok := parseFirst(t, sql).(*nodes.CreateTableStmt)
	if !ok {
		t.Fatalf("not a CreateTableStmt: %T", parseFirst(t, sql))
	}
	return stmt
}

func mustParseView(t *testing.T, sql string) *nodes.CreateViewStmt {
	t.Helper()
	stmt, ok := parseFirst(t, sql).(*nodes.CreateViewStmt)
	if !ok {
		t.Fatalf("not a CreateViewStmt: %T", parseFirst(t, sql))
	}
	return stmt
}

func mustParseIndex(t *testing.T, sql string) *nodes.CreateIndexStmt {
	t.Helper()
	stmt, ok := parseFirst(t, sql).(*nodes.CreateIndexStmt)
	if !ok {
		t.Fatalf("not a CreateIndexStmt: %T", parseFirst(t, sql))
	}
	return stmt
}

func mustParseRoutine(t *testing.T, sql string) *nodes.CreateFunctionStmt {
	t.Helper()
	stmt, ok := parseFirst(t, sql).(*nodes.CreateFunctionStmt)
	if !ok {
		t.Fatalf("not a CreateFunctionStmt: %T", parseFirst(t, sql))
	}
	return stmt
}

func mustParseTrigger(t *testing.T, sql string) *nodes.CreateTriggerStmt {
	t.Helper()
	stmt, ok := parseFirst(t, sql).(*nodes.CreateTriggerStmt)
	if !ok {
		t.Fatalf("not a CreateTriggerStmt: %T", parseFirst(t, sql))
	}
	return stmt
}

func mustParseEvent(t *testing.T, sql string) *nodes.CreateEventStmt {
	t.Helper()
	stmt, ok := parseFirst(t, sql).(*nodes.CreateEventStmt)
	if !ok {
		t.Fatalf("not a CreateEventStmt: %T", parseFirst(t, sql))
	}
	return stmt
}

// assertErrCode fails the test unless err is a *Error with the given code.
func assertErrCode(t *testing.T, err error, code int) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected *Error(code=%d), got nil", code)
	}
	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("expected *Error, got %T: %v", err, err)
	}
	if e.Code != code {
		t.Fatalf("expected code %d, got %d (%s)", code, e.Code, e.Message)
	}
}

// newCatalogWithDB returns a fresh catalog with db created and selected.
func newCatalogWithDB(t *testing.T, name string) *Catalog {
	t.Helper()
	c := New()
	if _, err := c.Exec("CREATE DATABASE "+name+"; USE "+name+";", nil); err != nil {
		t.Fatalf("setup: %v", err)
	}
	return c
}

// -- §6.1 Happy-path per kind -------------------------------------------

func TestDefineDatabase_HappyPath(t *testing.T) {
	c := New()
	if err := c.DefineDatabase(mustParseDatabase(t, "CREATE DATABASE mydb")); err != nil {
		t.Fatalf("DefineDatabase: %v", err)
	}
	if db := c.GetDatabase("mydb"); db == nil || db.Name != "mydb" {
		t.Fatalf("database not registered: %+v", db)
	}
}

func TestDefineTable_HappyPath(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	stmt := mustParseTable(t, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(50))")
	if err := c.DefineTable(stmt); err != nil {
		t.Fatalf("DefineTable: %v", err)
	}
	got := c.ShowCreateTable("mydb", "t")
	if got == "" {
		t.Fatal("ShowCreateTable returned empty")
	}
}

func TestDefineView_HappyPath(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (id INT PRIMARY KEY)")); err != nil {
		t.Fatal(err)
	}
	stmt := mustParseView(t, "CREATE VIEW v AS SELECT id FROM t")
	if err := c.DefineView(stmt); err != nil {
		t.Fatalf("DefineView: %v", err)
	}
	db := c.GetDatabase("mydb")
	view := db.Views["v"]
	if view == nil {
		t.Fatal("view not registered")
	}
	if view.AnalyzedQuery == nil {
		t.Fatal("AnalyzedQuery should be non-nil when referenced table is installed")
	}
}

func TestDefineIndex_HappyPath(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (id INT, name VARCHAR(50))")); err != nil {
		t.Fatal(err)
	}
	if err := c.DefineIndex(mustParseIndex(t, "CREATE INDEX idx_name ON t (name)")); err != nil {
		t.Fatalf("DefineIndex: %v", err)
	}
	db := c.GetDatabase("mydb")
	tbl := db.Tables["t"]
	if tbl == nil {
		t.Fatal("table missing")
	}
	found := false
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_name" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("idx_name not found in table indexes")
	}
}

func TestDefineFunction_HappyPath(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	stmt := mustParseRoutine(t, "CREATE FUNCTION f() RETURNS INT RETURN 1")
	if err := c.DefineFunction(stmt); err != nil {
		t.Fatalf("DefineFunction: %v", err)
	}
	if got := c.ShowCreateFunction("mydb", "f"); got == "" {
		t.Fatal("ShowCreateFunction empty")
	}
}

func TestDefineProcedure_HappyPath(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	stmt := mustParseRoutine(t, "CREATE PROCEDURE p() BEGIN END")
	if err := c.DefineProcedure(stmt); err != nil {
		t.Fatalf("DefineProcedure: %v", err)
	}
	if got := c.ShowCreateProcedure("mydb", "p"); got == "" {
		t.Fatal("ShowCreateProcedure empty")
	}
}

func TestDefineRoutine_FunctionRoutesToFunctions(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	stmt := mustParseRoutine(t, "CREATE FUNCTION f() RETURNS INT RETURN 1")
	if stmt.IsProcedure {
		t.Fatalf("expected IsProcedure=false for CREATE FUNCTION")
	}
	if err := c.DefineRoutine(stmt); err != nil {
		t.Fatal(err)
	}
	db := c.GetDatabase("mydb")
	if _, ok := db.Functions["f"]; !ok {
		t.Fatal("f not in db.Functions")
	}
	if _, ok := db.Procedures["f"]; ok {
		t.Fatal("f should not be in db.Procedures")
	}
}

func TestDefineTrigger_HappyPath(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (id INT PRIMARY KEY)")); err != nil {
		t.Fatal(err)
	}
	stmt := mustParseTrigger(t, "CREATE TRIGGER trg_ins BEFORE INSERT ON t FOR EACH ROW BEGIN END")
	if err := c.DefineTrigger(stmt); err != nil {
		t.Fatalf("DefineTrigger: %v", err)
	}
	if got := c.ShowCreateTrigger("mydb", "trg_ins"); got == "" {
		t.Fatal("ShowCreateTrigger empty")
	}
}

func TestDefineEvent_HappyPath(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	stmt := mustParseEvent(t, "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1")
	if err := c.DefineEvent(stmt); err != nil {
		t.Fatalf("DefineEvent: %v", err)
	}
	if got := c.ShowCreateEvent("mydb", "e"); got == "" {
		t.Fatal("ShowCreateEvent empty")
	}
}

// -- §6.2 Duplicate install ----------------------------------------------

func TestDefineDatabase_Duplicate(t *testing.T) {
	c := New()
	if err := c.DefineDatabase(mustParseDatabase(t, "CREATE DATABASE mydb")); err != nil {
		t.Fatal(err)
	}
	err := c.DefineDatabase(mustParseDatabase(t, "CREATE DATABASE mydb"))
	assertErrCode(t, err, ErrDupDatabase)
}

func TestDefineTable_Duplicate(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (id INT)")); err != nil {
		t.Fatal(err)
	}
	err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (id INT)"))
	assertErrCode(t, err, ErrDupTable)
}

func TestDefineFunction_Duplicate(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	stmt := mustParseRoutine(t, "CREATE FUNCTION f() RETURNS INT RETURN 1")
	if err := c.DefineFunction(stmt); err != nil {
		t.Fatal(err)
	}
	err := c.DefineFunction(mustParseRoutine(t, "CREATE FUNCTION f() RETURNS INT RETURN 1"))
	assertErrCode(t, err, ErrDupFunction)
}

func TestDefineTrigger_Duplicate(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (id INT)")); err != nil {
		t.Fatal(err)
	}
	stmt := mustParseTrigger(t, "CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW BEGIN END")
	if err := c.DefineTrigger(stmt); err != nil {
		t.Fatal(err)
	}
	err := c.DefineTrigger(mustParseTrigger(t, "CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW BEGIN END"))
	assertErrCode(t, err, ErrDupTrigger)
}

func TestDefineEvent_Duplicate(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	stmt := mustParseEvent(t, "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1")
	if err := c.DefineEvent(stmt); err != nil {
		t.Fatal(err)
	}
	err := c.DefineEvent(mustParseEvent(t, "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1"))
	assertErrCode(t, err, ErrDupEvent)
}

// -- §6.3 Missing currentDB ----------------------------------------------

func TestDefineTable_NoCurrentDatabase(t *testing.T) {
	c := New()
	err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (id INT)"))
	assertErrCode(t, err, ErrNoDatabaseSelected)
}

func TestDefineFunction_NoCurrentDatabase(t *testing.T) {
	c := New()
	err := c.DefineFunction(mustParseRoutine(t, "CREATE FUNCTION f() RETURNS INT RETURN 1"))
	assertErrCode(t, err, ErrNoDatabaseSelected)
}

// -- §6.4 Nil / incomplete stmt ------------------------------------------

func TestDefine_NilStmt(t *testing.T) {
	c := New()
	assertErrCode(t, c.DefineDatabase(nil), ErrWrongArguments)
	assertErrCode(t, c.DefineTable(nil), ErrWrongArguments)
	assertErrCode(t, c.DefineView(nil), ErrWrongArguments)
	assertErrCode(t, c.DefineIndex(nil), ErrWrongArguments)
	assertErrCode(t, c.DefineFunction(nil), ErrWrongArguments)
	assertErrCode(t, c.DefineProcedure(nil), ErrWrongArguments)
	assertErrCode(t, c.DefineRoutine(nil), ErrWrongArguments)
	assertErrCode(t, c.DefineTrigger(nil), ErrWrongArguments)
	assertErrCode(t, c.DefineEvent(nil), ErrWrongArguments)
}

func TestDefine_IncompleteStmt(t *testing.T) {
	c := New()
	assertErrCode(t, c.DefineDatabase(&nodes.CreateDatabaseStmt{}), ErrWrongArguments)
	assertErrCode(t, c.DefineTable(&nodes.CreateTableStmt{}), ErrWrongArguments)
	assertErrCode(t, c.DefineTable(&nodes.CreateTableStmt{Table: &nodes.TableRef{}}), ErrWrongArguments)
	assertErrCode(t, c.DefineView(&nodes.CreateViewStmt{}), ErrWrongArguments)
	assertErrCode(t, c.DefineIndex(&nodes.CreateIndexStmt{}), ErrWrongArguments)
	assertErrCode(t, c.DefineFunction(&nodes.CreateFunctionStmt{}), ErrWrongArguments)
	assertErrCode(t, c.DefineTrigger(&nodes.CreateTriggerStmt{}), ErrWrongArguments)
	assertErrCode(t, c.DefineTrigger(&nodes.CreateTriggerStmt{Name: "x"}), ErrWrongArguments)
	assertErrCode(t, c.DefineEvent(&nodes.CreateEventStmt{}), ErrWrongArguments)
}

// -- §6.5 Cross-database explicit schema --------------------------------

func TestDefineTable_ExplicitSchemaBypassesCurrentDB(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE main_db; CREATE DATABASE other_db; USE main_db;", nil); err != nil {
		t.Fatal(err)
	}
	// currentDB is main_db; install into other_db via qualifier.
	stmt := mustParseTable(t, "CREATE TABLE other_db.t (id INT)")
	if err := c.DefineTable(stmt); err != nil {
		t.Fatalf("DefineTable: %v", err)
	}
	if tbl := c.GetDatabase("other_db").Tables["t"]; tbl == nil {
		t.Fatal("table should be in other_db")
	}
	if tbl := c.GetDatabase("main_db").Tables["t"]; tbl != nil {
		t.Fatal("table should NOT be in main_db")
	}
}

// -- §6.6 FK forward reference under foreignKeyChecks=false -------------

func TestDefineTable_FKForwardRefWithChecksOff(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	c.SetForeignKeyChecks(false)
	defer c.SetForeignKeyChecks(true)

	// child references parent which is NOT yet installed.
	child := mustParseTable(t, `
		CREATE TABLE child (
			id INT PRIMARY KEY,
			parent_id INT,
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent (id)
		)`)
	if err := c.DefineTable(child); err != nil {
		t.Fatalf("child install with checks off should succeed, got: %v", err)
	}

	parent := mustParseTable(t, "CREATE TABLE parent (id INT PRIMARY KEY)")
	if err := c.DefineTable(parent); err != nil {
		t.Fatalf("parent install: %v", err)
	}

	db := c.GetDatabase("mydb")
	if db.Tables["child"] == nil || db.Tables["parent"] == nil {
		t.Fatal("both tables should be present")
	}
	// Confirm the FK struct was carried on child even though unvalidated.
	var fkFound bool
	for _, con := range db.Tables["child"].Constraints {
		if con.Type == ConForeignKey {
			fkFound = true
			break
		}
	}
	if !fkFound {
		t.Fatal("child should still carry FK constraint struct (unvalidated but present)")
	}
}

func TestDefineTable_FKForwardRefWithChecksOn(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	// Default is FK checks ON.
	child := mustParseTable(t, `
		CREATE TABLE child (
			id INT PRIMARY KEY,
			parent_id INT,
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent (id)
		)`)
	err := c.DefineTable(child)
	assertErrCode(t, err, ErrFKNoRefTable)
}

// -- §6.7 Trigger ordering (both directions) ----------------------------

func TestDefineTrigger_BeforeTable(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	stmt := mustParseTrigger(t, "CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW BEGIN END")
	err := c.DefineTrigger(stmt)
	assertErrCode(t, err, ErrNoSuchTable)
}

func TestDefineTrigger_AfterTable(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (id INT)")); err != nil {
		t.Fatal(err)
	}
	stmt := mustParseTrigger(t, "CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW BEGIN END")
	if err := c.DefineTrigger(stmt); err != nil {
		t.Fatalf("DefineTrigger after table install should succeed: %v", err)
	}
}

// -- §6.8 View forward reference — loader contract ----------------------

func TestDefineView_ForwardRefDegradesGracefully(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	// View references a table that doesn't exist yet.
	stmt := mustParseView(t, "CREATE VIEW v AS SELECT id FROM missing_table")
	if err := c.DefineView(stmt); err != nil {
		t.Fatalf("DefineView with forward ref should NOT error (loader contract): %v", err)
	}
	view := c.GetDatabase("mydb").Views["v"]
	if view == nil {
		t.Fatal("view should be registered even without resolved deps")
	}
	if view.AnalyzedQuery != nil {
		t.Fatal("AnalyzedQuery should be nil when reference cannot be resolved")
	}
}

func TestDefineView_ResolvedRefHasAnalyzedQuery(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE t (id INT)")); err != nil {
		t.Fatal(err)
	}
	if err := c.DefineView(mustParseView(t, "CREATE VIEW v AS SELECT id FROM t")); err != nil {
		t.Fatal(err)
	}
	view := c.GetDatabase("mydb").Views["v"]
	if view.AnalyzedQuery == nil {
		t.Fatal("AnalyzedQuery should be populated when reference resolves")
	}
}

// -- §6.9 Routine dispatch by IsProcedure -------------------------------

func TestDefineRoutine_ProcedureRoutesToProcedures(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	stmt := mustParseRoutine(t, "CREATE PROCEDURE p() BEGIN END")
	if !stmt.IsProcedure {
		t.Fatalf("expected IsProcedure=true for CREATE PROCEDURE")
	}
	if err := c.DefineRoutine(stmt); err != nil {
		t.Fatal(err)
	}
	db := c.GetDatabase("mydb")
	if _, ok := db.Procedures["p"]; !ok {
		t.Fatal("p not in db.Procedures")
	}
	if _, ok := db.Functions["p"]; ok {
		t.Fatal("p should not be in db.Functions")
	}
}

// DefineFunction with IsProcedure=true routes by stmt bit (no kind guard).
// This documents the loader philosophy: we do not reject mismatched
// metadata; we route it where the stmt says it belongs.
func TestDefineFunction_WithProcedureStmt_RoutesByBit(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	stmt := mustParseRoutine(t, "CREATE PROCEDURE p() BEGIN END")
	if !stmt.IsProcedure {
		t.Fatal("expected IsProcedure=true")
	}
	if err := c.DefineFunction(stmt); err != nil {
		t.Fatalf("DefineFunction with procedure stmt should not error; loader routes by stmt bit: %v", err)
	}
	db := c.GetDatabase("mydb")
	if _, ok := db.Procedures["p"]; !ok {
		t.Fatal("procedure should land in db.Procedures regardless of entry-point name")
	}
}

// -- §6.10 Partial-load downstream usability ----------------------------

// A loader with a broken table should still produce a catalog where
// queries over healthy tables analyze correctly. This is the end-to-end
// proof that the loader contract delivers isolation.
func TestDefine_PartialLoadPreservesHealthyAnalysis(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	c.SetForeignKeyChecks(false)
	defer c.SetForeignKeyChecks(true)

	// Two healthy tables.
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50))")); err != nil {
		t.Fatal(err)
	}
	if err := c.DefineTable(mustParseTable(t, "CREATE TABLE orders (id INT PRIMARY KEY, user_id INT, amount INT)")); err != nil {
		t.Fatal(err)
	}
	// One broken table: FK to non-existent table. Under checks-off, it still installs.
	broken := mustParseTable(t, `
		CREATE TABLE broken (
			id INT PRIMARY KEY,
			ghost_id INT,
			CONSTRAINT fk_ghost FOREIGN KEY (ghost_id) REFERENCES ghost_table (id)
		)`)
	if err := c.DefineTable(broken); err != nil {
		t.Fatalf("broken table install should succeed under checks-off: %v", err)
	}

	// Analyze a query joining the two healthy tables. Must not be poisoned
	// by the broken table's dangling FK.
	stmts, err := parser.Parse("SELECT u.name, o.amount FROM users u JOIN orders o ON o.user_id = u.id")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	sel, ok := stmts.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("not a SelectStmt")
	}
	q, err := c.AnalyzeSelectStmt(sel)
	if err != nil {
		t.Fatalf("AnalyzeSelectStmt on healthy tables should succeed despite broken peer: %v", err)
	}
	if q == nil {
		t.Fatal("Query should be non-nil")
	}
}

// -- §6.11 Parity with Exec (scoped to CREATE kinds) --------------------

func TestDefineTable_ParityWithExec(t *testing.T) {
	stmt := "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(50))"

	cDefine := newCatalogWithDB(t, "mydb")
	if err := cDefine.DefineTable(mustParseTable(t, stmt)); err != nil {
		t.Fatal(err)
	}

	cExec := newCatalogWithDB(t, "mydb")
	if _, err := cExec.Exec(stmt, nil); err != nil {
		t.Fatal(err)
	}

	if a, b := cDefine.ShowCreateTable("mydb", "t"), cExec.ShowCreateTable("mydb", "t"); a != b {
		t.Fatalf("ShowCreateTable diverges:\nDefine:\n%s\nExec:\n%s", a, b)
	}
}

func TestDefineView_ParityWithExec(t *testing.T) {
	setup := "CREATE TABLE t (id INT);"
	stmt := "CREATE VIEW v AS SELECT id FROM t"

	cDefine := newCatalogWithDB(t, "mydb")
	if _, err := cDefine.Exec(setup, nil); err != nil {
		t.Fatal(err)
	}
	if err := cDefine.DefineView(mustParseView(t, stmt)); err != nil {
		t.Fatal(err)
	}

	cExec := newCatalogWithDB(t, "mydb")
	if _, err := cExec.Exec(setup+stmt, nil); err != nil {
		t.Fatal(err)
	}

	if a, b := cDefine.ShowCreateView("mydb", "v"), cExec.ShowCreateView("mydb", "v"); a != b {
		t.Fatalf("ShowCreateView diverges:\nDefine:\n%s\nExec:\n%s", a, b)
	}
}
