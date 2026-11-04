package catalog

import "testing"

// TiDB's ADMIN statements change no schema object the catalog keeps, so
// Exec takes each of them and leaves the catalog as it was. The one that
// needs saying out loud is ADMIN REPAIR TABLE: it carries a whole CREATE
// TABLE, and a catalog that dispatched on the definition rather than on the
// statement would create or redefine a table on a statement that does
// neither -- outside repair mode TiDB answers 8215, and in repair mode the
// statement rewrites the metadata of a table that is already there.
//
// These tests run unconditionally; the differential tests behind
// TIDB_CATALOG_PARITY are a separate gate.

func TestExecAdminChangesNothing(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, a INT)")

	before := len(c.GetDatabase("testdb").Tables)

	for _, sql := range []string{
		"ADMIN SHOW DDL",
		"ADMIN SHOW DDL JOBS 5 WHERE id = 0",
		"ADMIN SHOW DDL JOB QUERIES LIMIT 2 OFFSET 1",
		"ADMIN SHOW t NEXT_ROW_ID",
		"ADMIN SHOW SLOW TOP ALL 3",
		"ADMIN SHOW BDR ROLE",
		"ADMIN CHECK TABLE t",
		"ADMIN CHECK INDEX t i (1, 2)",
		"ADMIN CHECKSUM TABLE t",
		"ADMIN RECOVER INDEX t i",
		"ADMIN CLEANUP INDEX t i",
		"ADMIN CLEANUP TABLE LOCK t",
		"ADMIN CANCEL DDL JOBS 1",
		"ADMIN PAUSE DDL JOBS 1",
		"ADMIN RESUME DDL JOBS 1",
		"ADMIN ALTER DDL JOBS 1 THREAD = 2",
		"ADMIN RELOAD EXPR_PUSHDOWN_BLACKLIST",
		"ADMIN RELOAD BINDINGS",
		"ADMIN FLUSH BINDINGS",
		"ADMIN CAPTURE BINDINGS",
		"ADMIN EVOLVE BINDINGS",
		"ADMIN FLUSH INSTANCE PLAN_CACHE",
		"ADMIN PLUGINS ENABLE audit",
		"ADMIN SET BDR ROLE PRIMARY",
		"ADMIN UNSET BDR ROLE",
	} {
		mustExec(t, c, sql)
	}

	db := c.GetDatabase("testdb")
	if got := len(db.Tables); got != before {
		t.Errorf("the database holds %d tables after the ADMIN statements, want %d", got, before)
	}
	tbl := db.GetTable("t")
	if tbl == nil {
		t.Fatal("table t is gone")
	}
	if len(tbl.Columns) != 2 {
		t.Errorf("t has %d columns, want the 2 it was created with", len(tbl.Columns))
	}
}

func TestExecAdminRepairTableCreatesNothing(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, a INT)")

	// A definition of a table that is not there leaves it not there.
	mustExec(t, c, "ADMIN REPAIR TABLE t CREATE TABLE fresh (x INT, y INT, z INT)")
	if tbl := c.GetDatabase("testdb").GetTable("fresh"); tbl != nil {
		t.Errorf("ADMIN REPAIR TABLE created table fresh with %d columns", len(tbl.Columns))
	}

	// A definition of the table itself leaves it as it was, columns and
	// all: the catalog does not apply the definition to it.
	mustExec(t, c, "ADMIN REPAIR TABLE t CREATE TABLE t (only_one INT)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table t is gone")
	}
	if len(tbl.Columns) != 2 {
		t.Fatalf("t has %d columns, want the 2 it was created with", len(tbl.Columns))
	}
	if tbl.GetColumn("id") == nil || tbl.GetColumn("a") == nil {
		t.Errorf("t's columns are %v, want id and a", tbl.Columns)
	}
	if tbl.GetColumn("only_one") != nil {
		t.Error("t took the column the REPAIR's definition named")
	}
}
