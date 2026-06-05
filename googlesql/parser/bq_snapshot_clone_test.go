package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit tests for the parser-ddl-bigquery node: CREATE SNAPSHOT TABLE … CLONE …
// and DROP SNAPSHOT TABLE. BigQuery-only at the union level (Spanner rejects
// SNAPSHOT outright, probed 2026-06-05); verdicts triangulated against the legacy
// GoogleSQLParser.g4 + BigQuery truth1 (DDL-008/045). (CREATE TABLE … CLONE/COPY/
// LIKE is core CREATE TABLE, covered by the parser-ddl node — not here.)

func snapOf(t *testing.T, sql string) *ast.CreateSnapshotStmt {
	t.Helper()
	n := parseDDL(t, sql)
	s, ok := n.(*ast.CreateSnapshotStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.CreateSnapshotStmt", sql, n)
	}
	return s
}

func TestCreateSnapshotTable_Basic(t *testing.T) {
	s := snapOf(t, "CREATE SNAPSHOT TABLE mydataset.my_snapshot CLONE mydataset.mytable")
	if s.Name.String() != "mydataset.my_snapshot" {
		t.Errorf("Name = %q", s.Name.String())
	}
	if s.CloneSource == nil || s.CloneSource.String() != "mydataset.mytable" {
		t.Errorf("CloneSource = %v", s.CloneSource)
	}
}

func TestCreateSnapshotTable_ForSystemTime(t *testing.T) {
	// DDL-008.
	s := snapOf(t, "CREATE SNAPSHOT TABLE mydataset.my_snapshot CLONE mydataset.mytable FOR SYSTEM_TIME AS OF TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 HOUR)")
	if s.ForSystemTime == nil {
		t.Error("ForSystemTime = nil, want the time expression")
	}
}

func TestCreateSnapshotTable_IfNotExistsOptions(t *testing.T) {
	s := snapOf(t, "CREATE SNAPSHOT TABLE IF NOT EXISTS ds.s CLONE ds.t OPTIONS(expiration_timestamp=TIMESTAMP '2025-01-01 00:00:00 UTC')")
	if !s.IfNotExists {
		t.Error("IfNotExists = false, want true")
	}
	if len(s.Options) != 1 {
		t.Errorf("Options = %d, want 1", len(s.Options))
	}
}

func TestCreateSnapshotTable_SpaceSystemTime(t *testing.T) {
	// `FOR SYSTEM TIME AS OF` (two-word spelling) is accepted as well as SYSTEM_TIME.
	s := snapOf(t, "CREATE SNAPSHOT TABLE ds.s CLONE ds.t FOR SYSTEM TIME AS OF CURRENT_TIMESTAMP()")
	if s.ForSystemTime == nil {
		t.Error("ForSystemTime = nil")
	}
}

func TestCreateSnapshotTable_Rejects(t *testing.T) {
	cases := []string{
		"CREATE SNAPSHOT TABLE s",                              // missing CLONE source
		"CREATE SNAPSHOT s CLONE t",                            // missing TABLE keyword
		"CREATE SNAPSHOT TABLE s CLONE t FOR SYSTEM_TIME AS x", // bad FOR SYSTEM_TIME (missing OF)
		// Regression (review finding): create_snapshot_statement has NO
		// opt_create_scope; a leading TEMP/TEMPORARY must reject (not be silently
		// dropped).
		"CREATE TEMP SNAPSHOT TABLE s CLONE t",
		"CREATE TEMPORARY SNAPSHOT TABLE s CLONE t",
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

func TestDropSnapshotTable(t *testing.T) {
	// DDL-045.
	d := bqDropOf(t, "DROP SNAPSHOT TABLE IF EXISTS mydataset.mytablesnapshot")
	if d.Object != ast.BQDropSnapshotTable {
		t.Errorf("Object = %v, want SNAPSHOT TABLE", d.Object)
	}
	if !d.IfExists {
		t.Error("IfExists = false, want true")
	}
}
