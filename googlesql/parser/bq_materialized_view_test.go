package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit tests for the parser-ddl-bigquery node: CREATE MATERIALIZED|APPROX VIEW,
// ALTER MATERIALIZED VIEW SET OPTIONS, DROP MATERIALIZED VIEW. BigQuery-only at
// the union level (Spanner rejects MATERIALIZED VIEW outright, probed
// 2026-06-05); accept/reject verdicts triangulated against the legacy
// GoogleSQLParser.g4 + BigQuery truth1 (DDL-011/012/040/048).

func mvOf(t *testing.T, sql string) *ast.CreateMaterializedViewStmt {
	t.Helper()
	n := parseDDL(t, sql)
	mv, ok := n.(*ast.CreateMaterializedViewStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.CreateMaterializedViewStmt", sql, n)
	}
	return mv
}

func TestCreateMaterializedView_Basic(t *testing.T) {
	// DDL-011.
	mv := mvOf(t, "CREATE MATERIALIZED VIEW myProject.myDataset.myView AS SELECT product, SUM(sales) AS total FROM mydataset.sales GROUP BY product")
	if mv.Kind != ast.ViewMaterialized {
		t.Errorf("Kind = %v, want MATERIALIZED VIEW", mv.Kind)
	}
	if mv.Name.String() != "myProject.myDataset.myView" {
		t.Errorf("Name = %q", mv.Name.String())
	}
	if mv.AsQuery == nil {
		t.Error("AsQuery = nil")
	}
}

func TestCreateMaterializedView_PartitionCluster(t *testing.T) {
	// DDL-011: PARTITION BY DATE(ts) CLUSTER BY product.
	mv := mvOf(t, "CREATE MATERIALIZED VIEW p.d.v PARTITION BY DATE(ts) CLUSTER BY product OPTIONS(enable_refresh=true) AS SELECT product, ts, SUM(sales) AS total FROM mydataset.sales GROUP BY product, ts")
	if len(mv.PartitionBy) != 1 {
		t.Errorf("PartitionBy = %d, want 1", len(mv.PartitionBy))
	}
	if len(mv.ClusterBy) != 1 {
		t.Errorf("ClusterBy = %d, want 1", len(mv.ClusterBy))
	}
	if len(mv.Options) != 1 {
		t.Errorf("Options = %d, want 1", len(mv.Options))
	}
}

func TestCreateMaterializedView_ReplicaOf(t *testing.T) {
	// DDL-012: CREATE MATERIALIZED VIEW … AS REPLICA OF source.
	mv := mvOf(t, "CREATE MATERIALIZED VIEW mydataset.my_replica AS REPLICA OF other_project.other_dataset.my_mv")
	if mv.ReplicaOf == nil || mv.ReplicaOf.String() != "other_project.other_dataset.my_mv" {
		t.Errorf("ReplicaOf = %v, want other_project.other_dataset.my_mv", mv.ReplicaOf)
	}
	if mv.AsQuery != nil {
		t.Error("AsQuery != nil, want nil for a REPLICA OF view")
	}
}

func TestCreateMaterializedView_ReplicaOfOptions(t *testing.T) {
	// Review finding / DDL-012: trailing OPTIONS(...) after AS REPLICA OF source.
	// The legacy .g4 omits this (OPTIONS only before AS), but BigQuery documents it
	// and accepting it is additive (Spanner rejects the whole REPLICA-OF view).
	mv := mvOf(t, "CREATE MATERIALIZED VIEW ds.r AS REPLICA OF ds.src OPTIONS(replication_interval_seconds=300)")
	if mv.ReplicaOf == nil || mv.ReplicaOf.String() != "ds.src" {
		t.Errorf("ReplicaOf = %v, want ds.src", mv.ReplicaOf)
	}
	if len(mv.Options) != 1 {
		t.Errorf("Options = %d, want 1 (trailing OPTIONS after REPLICA OF)", len(mv.Options))
	}
}

func TestCreateApproxView(t *testing.T) {
	mv := mvOf(t, "CREATE APPROX VIEW ds.v OPTIONS(description='d') AS SELECT 1 AS n")
	if mv.Kind != ast.ViewApprox {
		t.Errorf("Kind = %v, want APPROX VIEW", mv.Kind)
	}
}

func TestCreateMaterializedView_OrReplaceColumnsSecurity(t *testing.T) {
	mv := mvOf(t, "CREATE OR REPLACE MATERIALIZED VIEW IF NOT EXISTS ds.v(a, b) SQL SECURITY INVOKER AS SELECT 1 AS a, 2 AS b")
	if !mv.OrReplace || !mv.IfNotExists {
		t.Errorf("OrReplace=%v IfNotExists=%v, want both true", mv.OrReplace, mv.IfNotExists)
	}
	if len(mv.Columns) != 2 {
		t.Errorf("Columns = %d, want 2", len(mv.Columns))
	}
	if mv.SQLSecurity != "INVOKER" {
		t.Errorf("SQLSecurity = %q, want INVOKER", mv.SQLSecurity)
	}
}

func TestCreateMaterializedView_Rejects(t *testing.T) {
	cases := []string{
		"CREATE MATERIALIZED VIEW v",                      // missing AS body
		"CREATE MATERIALIZED v AS SELECT 1",               // missing VIEW keyword
		"CREATE TEMP MATERIALIZED VIEW v AS SELECT 1",     // MATERIALIZED takes no scope
		"CREATE APPROX VIEW v PARTITION BY x AS SELECT 1", // PARTITION BY is materialized-only; rejects on APPROX
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

// --- ALTER MATERIALIZED VIEW ---

func bqAlterOf(t *testing.T, sql string) *ast.BQAlterStmt {
	t.Helper()
	n := parseDDL(t, sql)
	a, ok := n.(*ast.BQAlterStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.BQAlterStmt", sql, n)
	}
	return a
}

func TestAlterMaterializedView_SetOptions(t *testing.T) {
	// DDL-040.
	a := bqAlterOf(t, "ALTER MATERIALIZED VIEW mydataset.myview SET OPTIONS(enable_refresh=false)")
	if a.Object != ast.BQAlterMaterializedView {
		t.Errorf("Object = %v, want MATERIALIZED VIEW", a.Object)
	}
	if len(a.SetOptions) != 1 {
		t.Errorf("SetOptions = %d, want 1", len(a.SetOptions))
	}
}

func TestAlterMaterializedView_IfExists(t *testing.T) {
	a := bqAlterOf(t, "ALTER MATERIALIZED VIEW IF EXISTS ds.v SET OPTIONS(x=1)")
	if !a.IfExists {
		t.Error("IfExists = false, want true")
	}
}

func TestAlterMaterializedView_Rejects(t *testing.T) {
	cases := []string{
		"ALTER MATERIALIZED VIEW v",                    // missing SET OPTIONS
		"ALTER MATERIALIZED VIEW v ADD COLUMN c INT64", // unsupported action on an MV
		"ALTER MATERIALIZED v SET OPTIONS(x=1)",        // missing VIEW
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

// --- DROP MATERIALIZED VIEW ---

func TestDropMaterializedView(t *testing.T) {
	// DDL-048.
	d := bqDropOf(t, "DROP MATERIALIZED VIEW IF EXISTS mydataset.my_mv")
	if d.Object != ast.BQDropMaterializedView {
		t.Errorf("Object = %v, want MATERIALIZED VIEW", d.Object)
	}
	if !d.IfExists {
		t.Error("IfExists = false, want true")
	}
}
