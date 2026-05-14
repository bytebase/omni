package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func parseCreateMTMVStmt(t *testing.T, sql string) *ast.CreateMTMVStmt {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateMTMVStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateMTMVStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseAlterMTMVStmt(t *testing.T, sql string) *ast.AlterMTMVStmt {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterMTMVStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterMTMVStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseDropMTMVStmt(t *testing.T, sql string) *ast.DropMTMVStmt {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropMTMVStmt)
	if !ok {
		t.Fatalf("expected *ast.DropMTMVStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseRefreshMTMVStmt(t *testing.T, sql string) *ast.RefreshMTMVStmt {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.RefreshMTMVStmt)
	if !ok {
		t.Fatalf("expected *ast.RefreshMTMVStmt, got %T", file.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// CREATE MATERIALIZED VIEW
// ---------------------------------------------------------------------------

func TestCreateMTMV_Minimal(t *testing.T) {
	stmt := parseCreateMTMVStmt(t, "CREATE MATERIALIZED VIEW mv AS SELECT * FROM t")
	if stmt.Name.String() != "mv" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mv")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if stmt.BuildMode != "" {
		t.Errorf("BuildMode = %q, want empty", stmt.BuildMode)
	}
	if stmt.RefreshMethod != "" {
		t.Errorf("RefreshMethod = %q, want empty", stmt.RefreshMethod)
	}
	if stmt.Query == nil {
		t.Fatal("Query is nil")
	}
}

func TestCreateMTMV_IfNotExists(t *testing.T) {
	stmt := parseCreateMTMVStmt(t, "CREATE MATERIALIZED VIEW IF NOT EXISTS mv AS SELECT id FROM t")
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if stmt.Name.String() != "mv" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mv")
	}
}

func TestCreateMTMV_BuildImmediate_RefreshAuto(t *testing.T) {
	stmt := parseCreateMTMVStmt(t,
		"CREATE MATERIALIZED VIEW IF NOT EXISTS mv BUILD IMMEDIATE REFRESH AUTO AS SELECT id FROM t")
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if stmt.BuildMode != "IMMEDIATE" {
		t.Errorf("BuildMode = %q, want IMMEDIATE", stmt.BuildMode)
	}
	if stmt.RefreshMethod != "AUTO" {
		t.Errorf("RefreshMethod = %q, want AUTO", stmt.RefreshMethod)
	}
	if stmt.Query == nil {
		t.Fatal("Query is nil")
	}
}

func TestCreateMTMV_RefreshComplete_OnSchedule(t *testing.T) {
	sql := "CREATE MATERIALIZED VIEW mv REFRESH COMPLETE ON SCHEDULE EVERY 1 DAY AS SELECT * FROM t"
	stmt := parseCreateMTMVStmt(t, sql)
	if stmt.Name.String() != "mv" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mv")
	}
	if stmt.RefreshMethod != "COMPLETE" {
		t.Errorf("RefreshMethod = %q, want COMPLETE", stmt.RefreshMethod)
	}
	if stmt.RefreshTrigger == nil {
		t.Fatal("RefreshTrigger is nil")
	}
	if !stmt.RefreshTrigger.OnSchedule {
		t.Error("RefreshTrigger.OnSchedule should be true")
	}
	if stmt.RefreshTrigger.Interval != "1 DAY" {
		t.Errorf("Interval = %q, want %q", stmt.RefreshTrigger.Interval, "1 DAY")
	}
}

func TestCreateMTMV_OnSchedule_WithStarts(t *testing.T) {
	sql := `CREATE MATERIALIZED VIEW mv BUILD IMMEDIATE REFRESH AUTO ON SCHEDULE EVERY 1 DAY STARTS '2024-12-01 20:30:00' DISTRIBUTED BY HASH(orderkey) BUCKETS 2 PROPERTIES("replication_num" = "1") AS SELECT id FROM t`
	stmt := parseCreateMTMVStmt(t, sql)
	if stmt.BuildMode != "IMMEDIATE" {
		t.Errorf("BuildMode = %q, want IMMEDIATE", stmt.BuildMode)
	}
	if stmt.RefreshMethod != "AUTO" {
		t.Errorf("RefreshMethod = %q, want AUTO", stmt.RefreshMethod)
	}
	if stmt.RefreshTrigger == nil {
		t.Fatal("RefreshTrigger is nil")
	}
	if stmt.RefreshTrigger.StartsAt != "2024-12-01 20:30:00" {
		t.Errorf("StartsAt = %q, want %q", stmt.RefreshTrigger.StartsAt, "2024-12-01 20:30:00")
	}
	if stmt.DistributedBy == nil {
		t.Fatal("DistributedBy is nil")
	}
	if stmt.DistributedBy.Type != "HASH" {
		t.Errorf("DistributedBy.Type = %q, want HASH", stmt.DistributedBy.Type)
	}
	if stmt.DistributedBy.Buckets != 2 {
		t.Errorf("DistributedBy.Buckets = %d, want 2", stmt.DistributedBy.Buckets)
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("Properties: got %d, want 1", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "replication_num" || stmt.Properties[0].Value != "1" {
		t.Errorf("Properties[0] = {%q: %q}, want {replication_num: 1}",
			stmt.Properties[0].Key, stmt.Properties[0].Value)
	}
}

func TestCreateMTMV_WithColumnList(t *testing.T) {
	sql := `CREATE MATERIALIZED VIEW complete_mv (orderdate COMMENT 'Order date', orderkey COMMENT 'Order key', partkey COMMENT 'Part key') BUILD IMMEDIATE REFRESH AUTO ON SCHEDULE EVERY 1 DAY STARTS '2024-12-01 20:30:00' DISTRIBUTED BY HASH(orderkey) BUCKETS 2 PROPERTIES("replication_num" = "1") AS SELECT o_orderdate, l_orderkey, l_partkey FROM orders`
	stmt := parseCreateMTMVStmt(t, sql)
	if stmt.Name.String() != "complete_mv" {
		t.Errorf("Name = %q, want complete_mv", stmt.Name.String())
	}
	if len(stmt.Columns) != 3 {
		t.Fatalf("Columns: got %d, want 3", len(stmt.Columns))
	}
	if stmt.Columns[0].Name != "orderdate" {
		t.Errorf("Columns[0].Name = %q, want orderdate", stmt.Columns[0].Name)
	}
	if stmt.Columns[0].Comment != "Order date" {
		t.Errorf("Columns[0].Comment = %q, want 'Order date'", stmt.Columns[0].Comment)
	}
}

func TestCreateMTMV_WithPartitionBy(t *testing.T) {
	sql := `CREATE MATERIALIZED VIEW partition_mv BUILD IMMEDIATE REFRESH AUTO ON SCHEDULE EVERY 1 DAY STARTS '2024-12-01 20:30:00' PARTITION BY (o_orderdate) DISTRIBUTED BY HASH(l_orderkey) BUCKETS 2 PROPERTIES("replication_num" = "3") AS SELECT o_orderdate, l_orderkey FROM orders`
	stmt := parseCreateMTMVStmt(t, sql)
	if stmt.Name.String() != "partition_mv" {
		t.Errorf("Name = %q, want partition_mv", stmt.Name.String())
	}
	if stmt.PartitionBy == nil {
		t.Fatal("PartitionBy is nil")
	}
}

func TestCreateMTMV_DistributedByHash_Properties(t *testing.T) {
	sql := `CREATE MATERIALIZED VIEW mv DISTRIBUTED BY HASH(id) BUCKETS 3 PROPERTIES("replication_num"="3") AS SELECT id FROM t`
	stmt := parseCreateMTMVStmt(t, sql)
	if stmt.DistributedBy == nil {
		t.Fatal("DistributedBy is nil")
	}
	if stmt.DistributedBy.Type != "HASH" {
		t.Errorf("DistributedBy.Type = %q, want HASH", stmt.DistributedBy.Type)
	}
	if stmt.DistributedBy.Buckets != 3 {
		t.Errorf("DistributedBy.Buckets = %d, want 3", stmt.DistributedBy.Buckets)
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("Properties: got %d, want 1", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "replication_num" {
		t.Errorf("Properties[0].Key = %q, want replication_num", stmt.Properties[0].Key)
	}
}

func TestCreateMTMV_OnManual(t *testing.T) {
	sql := "CREATE MATERIALIZED VIEW mv REFRESH COMPLETE ON MANUAL AS SELECT * FROM t"
	stmt := parseCreateMTMVStmt(t, sql)
	if stmt.RefreshTrigger == nil {
		t.Fatal("RefreshTrigger is nil")
	}
	if !stmt.RefreshTrigger.OnManual {
		t.Error("RefreshTrigger.OnManual should be true")
	}
	if stmt.RefreshTrigger.OnSchedule {
		t.Error("RefreshTrigger.OnSchedule should be false")
	}
}

func TestCreateMTMV_OnCommit(t *testing.T) {
	sql := "CREATE MATERIALIZED VIEW mv REFRESH AUTO ON COMMIT AS SELECT * FROM t"
	stmt := parseCreateMTMVStmt(t, sql)
	if stmt.RefreshTrigger == nil {
		t.Fatal("RefreshTrigger is nil")
	}
	if !stmt.RefreshTrigger.OnCommit {
		t.Error("RefreshTrigger.OnCommit should be true")
	}
}

func TestCreateMTMV_RefreshNever(t *testing.T) {
	sql := "CREATE MATERIALIZED VIEW mv REFRESH NEVER AS SELECT * FROM t"
	stmt := parseCreateMTMVStmt(t, sql)
	if stmt.RefreshMethod != "NEVER" {
		t.Errorf("RefreshMethod = %q, want NEVER", stmt.RefreshMethod)
	}
}

func TestCreateMTMV_QualifiedName(t *testing.T) {
	stmt := parseCreateMTMVStmt(t, "CREATE MATERIALIZED VIEW db1.mv1 AS SELECT * FROM t")
	if stmt.Name.String() != "db1.mv1" {
		t.Errorf("Name = %q, want db1.mv1", stmt.Name.String())
	}
}

// Legacy corpus test: complete_mv (from materialized_view.sql line 1)
func TestCreateMTMV_Legacy_CompleteMV(t *testing.T) {
	sql := `CREATE MATERIALIZED VIEW complete_mv (orderdate COMMENT 'Order date', orderkey COMMENT 'Order key', partkey COMMENT 'Part key') BUILD IMMEDIATE REFRESH AUTO ON SCHEDULE EVERY 1 DAY STARTS '2024-12-01 20:30:00' DISTRIBUTED BY HASH(orderkey) BUCKETS 2 PROPERTIES("replication_num" = "1") AS SELECT o_orderdate, l_orderkey, l_partkey FROM orders LEFT JOIN lineitem ON l_orderkey = o_orderkey LEFT JOIN partsupp ON ps_partkey = l_partkey AND l_suppkey = ps_suppkey`
	stmt := parseCreateMTMVStmt(t, sql)
	if stmt.Name.String() != "complete_mv" {
		t.Errorf("Name = %q, want complete_mv", stmt.Name.String())
	}
	if stmt.BuildMode != "IMMEDIATE" {
		t.Errorf("BuildMode = %q, want IMMEDIATE", stmt.BuildMode)
	}
	if stmt.RefreshMethod != "AUTO" {
		t.Errorf("RefreshMethod = %q, want AUTO", stmt.RefreshMethod)
	}
	if stmt.RefreshTrigger == nil || !stmt.RefreshTrigger.OnSchedule {
		t.Error("RefreshTrigger.OnSchedule should be true")
	}
	if stmt.RefreshTrigger.Interval != "1 DAY" {
		t.Errorf("Interval = %q, want '1 DAY'", stmt.RefreshTrigger.Interval)
	}
	if stmt.RefreshTrigger.StartsAt != "2024-12-01 20:30:00" {
		t.Errorf("StartsAt = %q, want '2024-12-01 20:30:00'", stmt.RefreshTrigger.StartsAt)
	}
	if stmt.DistributedBy == nil || stmt.DistributedBy.Buckets != 2 {
		t.Errorf("DistributedBy.Buckets = %d, want 2", stmt.DistributedBy.Buckets)
	}
	if stmt.Query == nil {
		t.Fatal("Query is nil")
	}
}

// Legacy corpus test: partition_mv (from async_materialized_view.sql line 16)
func TestCreateMTMV_Legacy_PartitionMV(t *testing.T) {
	sql := `CREATE MATERIALIZED VIEW partition_mv BUILD IMMEDIATE REFRESH AUTO ON SCHEDULE EVERY 1 DAY STARTS '2024-12-01 20:30:00' PARTITION BY (DATE_TRUNC(o_orderdate, 'MONTH')) DISTRIBUTED BY HASH(l_orderkey) BUCKETS 2 PROPERTIES("replication_num" = "3") AS SELECT o_orderdate, l_orderkey, l_partkey FROM orders LEFT JOIN lineitem ON l_orderkey = o_orderkey LEFT JOIN partsupp ON ps_partkey = l_partkey AND l_suppkey = ps_suppkey`
	stmt := parseCreateMTMVStmt(t, sql)
	if stmt.Name.String() != "partition_mv" {
		t.Errorf("Name = %q, want partition_mv", stmt.Name.String())
	}
	if stmt.PartitionBy == nil {
		t.Fatal("PartitionBy is nil")
	}
	if stmt.DistributedBy == nil || stmt.DistributedBy.Buckets != 2 {
		t.Errorf("DistributedBy.Buckets = %d, want 2", stmt.DistributedBy.Buckets)
	}
	if len(stmt.Properties) != 1 || stmt.Properties[0].Key != "replication_num" {
		t.Error("expected replication_num property")
	}
}

// ---------------------------------------------------------------------------
// ALTER MATERIALIZED VIEW
// ---------------------------------------------------------------------------

func TestAlterMTMV_Rename(t *testing.T) {
	stmt := parseAlterMTMVStmt(t, "ALTER MATERIALIZED VIEW mv RENAME new_mv")
	if stmt.Name.String() != "mv" {
		t.Errorf("Name = %q, want mv", stmt.Name.String())
	}
	if stmt.NewName == nil {
		t.Fatal("NewName is nil")
	}
	if stmt.NewName.String() != "new_mv" {
		t.Errorf("NewName = %q, want new_mv", stmt.NewName.String())
	}
}

func TestAlterMTMV_RefreshComplete(t *testing.T) {
	// Legacy corpus: ALTER MATERIALIZED VIEW partition_mv REFRESH COMPLETE;
	stmt := parseAlterMTMVStmt(t, "ALTER MATERIALIZED VIEW partition_mv REFRESH COMPLETE")
	if stmt.Name.String() != "partition_mv" {
		t.Errorf("Name = %q, want partition_mv", stmt.Name.String())
	}
	if stmt.RefreshMethod != "COMPLETE" {
		t.Errorf("RefreshMethod = %q, want COMPLETE", stmt.RefreshMethod)
	}
}

func TestAlterMTMV_SetProperties(t *testing.T) {
	// Legacy corpus: ALTER MATERIALIZED VIEW partition_mv SET ("grace_period" = "10", ...);
	sql := `ALTER MATERIALIZED VIEW partition_mv SET ("grace_period" = "10", "excluded_trigger_tables" = "lineitem,partsupp")`
	stmt := parseAlterMTMVStmt(t, sql)
	if stmt.Name.String() != "partition_mv" {
		t.Errorf("Name = %q, want partition_mv", stmt.Name.String())
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("Properties: got %d, want 2", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "grace_period" || stmt.Properties[0].Value != "10" {
		t.Errorf("Properties[0] = {%q: %q}, want {grace_period: 10}",
			stmt.Properties[0].Key, stmt.Properties[0].Value)
	}
	if stmt.Properties[1].Key != "excluded_trigger_tables" {
		t.Errorf("Properties[1].Key = %q, want excluded_trigger_tables", stmt.Properties[1].Key)
	}
}

func TestAlterMTMV_ReplaceWith(t *testing.T) {
	stmt := parseAlterMTMVStmt(t, "ALTER MATERIALIZED VIEW mv REPLACE WITH MATERIALIZED VIEW mv2")
	if stmt.Name.String() != "mv" {
		t.Errorf("Name = %q, want mv", stmt.Name.String())
	}
	if !stmt.Replace {
		t.Error("Replace should be true")
	}
	if stmt.ReplaceTarget == nil {
		t.Fatal("ReplaceTarget is nil")
	}
	if stmt.ReplaceTarget.String() != "mv2" {
		t.Errorf("ReplaceTarget = %q, want mv2", stmt.ReplaceTarget.String())
	}
}

// ---------------------------------------------------------------------------
// DROP MATERIALIZED VIEW
// ---------------------------------------------------------------------------

func TestDropMTMV_Basic(t *testing.T) {
	stmt := parseDropMTMVStmt(t, "DROP MATERIALIZED VIEW mv1")
	if stmt.Name.String() != "mv1" {
		t.Errorf("Name = %q, want mv1", stmt.Name.String())
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
	if stmt.OnBase != nil {
		t.Error("OnBase should be nil")
	}
}

func TestDropMTMV_IfExists(t *testing.T) {
	stmt := parseDropMTMVStmt(t, "DROP MATERIALIZED VIEW IF EXISTS mv1")
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if stmt.Name.String() != "mv1" {
		t.Errorf("Name = %q, want mv1", stmt.Name.String())
	}
}

func TestDropMTMV_IfExistsQualified(t *testing.T) {
	// Legacy corpus: DROP MATERIALIZED VIEW IF EXISTS db1.mv1
	stmt := parseDropMTMVStmt(t, "DROP MATERIALIZED VIEW IF EXISTS db1.mv1")
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if stmt.Name.String() != "db1.mv1" {
		t.Errorf("Name = %q, want db1.mv1", stmt.Name.String())
	}
}

func TestDropMTMV_OnBaseTable(t *testing.T) {
	// Sync MV: DROP MATERIALIZED VIEW mv ON base_table
	stmt := parseDropMTMVStmt(t, "DROP MATERIALIZED VIEW mv ON base_table")
	if stmt.Name.String() != "mv" {
		t.Errorf("Name = %q, want mv", stmt.Name.String())
	}
	if stmt.OnBase == nil {
		t.Fatal("OnBase is nil")
	}
	if stmt.OnBase.String() != "base_table" {
		t.Errorf("OnBase = %q, want base_table", stmt.OnBase.String())
	}
}

// ---------------------------------------------------------------------------
// REFRESH MATERIALIZED VIEW
// ---------------------------------------------------------------------------

func TestRefreshMTMV_NoMode(t *testing.T) {
	stmt := parseRefreshMTMVStmt(t, "REFRESH MATERIALIZED VIEW mv1")
	if stmt.Name.String() != "mv1" {
		t.Errorf("Name = %q, want mv1", stmt.Name.String())
	}
	if stmt.Mode != "" {
		t.Errorf("Mode = %q, want empty", stmt.Mode)
	}
}

func TestRefreshMTMV_Auto(t *testing.T) {
	// Legacy corpus: REFRESH MATERIALIZED VIEW mv1 AUTO
	stmt := parseRefreshMTMVStmt(t, "REFRESH MATERIALIZED VIEW mv1 AUTO")
	if stmt.Name.String() != "mv1" {
		t.Errorf("Name = %q, want mv1", stmt.Name.String())
	}
	if stmt.Mode != "AUTO" {
		t.Errorf("Mode = %q, want AUTO", stmt.Mode)
	}
}

func TestRefreshMTMV_Complete(t *testing.T) {
	stmt := parseRefreshMTMVStmt(t, "REFRESH MATERIALIZED VIEW mv COMPLETE")
	if stmt.Mode != "COMPLETE" {
		t.Errorf("Mode = %q, want COMPLETE", stmt.Mode)
	}
}

func TestRefreshMTMV_Partitions(t *testing.T) {
	// Legacy corpus: REFRESH MATERIALIZED VIEW mv1 PARTITIONS(p_19950801_19950901, p_19950901_19951001)
	stmt := parseRefreshMTMVStmt(t,
		"REFRESH MATERIALIZED VIEW mv1 PARTITIONS(p_19950801_19950901, p_19950901_19951001)")
	if stmt.Name.String() != "mv1" {
		t.Errorf("Name = %q, want mv1", stmt.Name.String())
	}
	if len(stmt.Partitions) != 2 {
		t.Fatalf("Partitions: got %d, want 2", len(stmt.Partitions))
	}
	if stmt.Partitions[0] != "p_19950801_19950901" {
		t.Errorf("Partitions[0] = %q, want p_19950801_19950901", stmt.Partitions[0])
	}
	if stmt.Partitions[1] != "p_19950901_19951001" {
		t.Errorf("Partitions[1] = %q, want p_19950901_19951001", stmt.Partitions[1])
	}
}

// ---------------------------------------------------------------------------
// PAUSE / RESUME MATERIALIZED VIEW JOB
// ---------------------------------------------------------------------------

func TestPauseMTMVJob(t *testing.T) {
	file, errs := Parse("PAUSE MATERIALIZED VIEW JOB ON mv")
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("got %d stmts, want 1", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.PauseMTMVJobStmt)
	if !ok {
		t.Fatalf("expected *ast.PauseMTMVJobStmt, got %T", file.Stmts[0])
	}
	if stmt.Name.String() != "mv" {
		t.Errorf("Name = %q, want mv", stmt.Name.String())
	}
}

func TestResumeMTMVJob(t *testing.T) {
	file, errs := Parse("RESUME MATERIALIZED VIEW JOB ON mv")
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("got %d stmts, want 1", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.ResumeMTMVJobStmt)
	if !ok {
		t.Fatalf("expected *ast.ResumeMTMVJobStmt, got %T", file.Stmts[0])
	}
	if stmt.Name.String() != "mv" {
		t.Errorf("Name = %q, want mv", stmt.Name.String())
	}
}

// ---------------------------------------------------------------------------
// CANCEL MATERIALIZED VIEW TASK
// ---------------------------------------------------------------------------

func TestCancelMTMVTask(t *testing.T) {
	file, errs := Parse("CANCEL MATERIALIZED VIEW TASK 12345 ON mv")
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("got %d stmts, want 1", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CancelMTMVTaskStmt)
	if !ok {
		t.Fatalf("expected *ast.CancelMTMVTaskStmt, got %T", file.Stmts[0])
	}
	if stmt.TaskID != 12345 {
		t.Errorf("TaskID = %d, want 12345", stmt.TaskID)
	}
	if stmt.Name.String() != "mv" {
		t.Errorf("Name = %q, want mv", stmt.Name.String())
	}
}

// ---------------------------------------------------------------------------
// Node tags
// ---------------------------------------------------------------------------

func TestMTMVNodeTags(t *testing.T) {
	tests := []struct {
		node ast.Node
		want ast.NodeTag
		name string
	}{
		{&ast.CreateMTMVStmt{}, ast.T_CreateMTMVStmt, "CreateMTMVStmt"},
		{&ast.AlterMTMVStmt{}, ast.T_AlterMTMVStmt, "AlterMTMVStmt"},
		{&ast.DropMTMVStmt{}, ast.T_DropMTMVStmt, "DropMTMVStmt"},
		{&ast.RefreshMTMVStmt{}, ast.T_RefreshMTMVStmt, "RefreshMTMVStmt"},
		{&ast.MTMVRefreshTrigger{}, ast.T_MTMVRefreshTrigger, "MTMVRefreshTrigger"},
		{&ast.PauseMTMVJobStmt{}, ast.T_PauseMTMVJobStmt, "PauseMTMVJobStmt"},
		{&ast.ResumeMTMVJobStmt{}, ast.T_ResumeMTMVJobStmt, "ResumeMTMVJobStmt"},
		{&ast.CancelMTMVTaskStmt{}, ast.T_CancelMTMVTaskStmt, "CancelMTMVTaskStmt"},
	}
	for _, tt := range tests {
		if tt.node.Tag() != tt.want {
			t.Errorf("%s.Tag() = %v, want %v", tt.name, tt.node.Tag(), tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Tag string representation
// ---------------------------------------------------------------------------

func TestMTMVNodeTagStrings(t *testing.T) {
	tests := []struct {
		tag  ast.NodeTag
		want string
	}{
		{ast.T_CreateMTMVStmt, "CreateMTMVStmt"},
		{ast.T_AlterMTMVStmt, "AlterMTMVStmt"},
		{ast.T_DropMTMVStmt, "DropMTMVStmt"},
		{ast.T_RefreshMTMVStmt, "RefreshMTMVStmt"},
		{ast.T_MTMVRefreshTrigger, "MTMVRefreshTrigger"},
		{ast.T_PauseMTMVJobStmt, "PauseMTMVJobStmt"},
		{ast.T_ResumeMTMVJobStmt, "ResumeMTMVJobStmt"},
		{ast.T_CancelMTMVTaskStmt, "CancelMTMVTaskStmt"},
	}
	for _, tt := range tests {
		if tt.tag.String() != tt.want {
			t.Errorf("NodeTag(%d).String() = %q, want %q", int(tt.tag), tt.tag.String(), tt.want)
		}
	}
}
