package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func parseCreateMTMVStmt(t *testing.T, sql string) *ast.CreateMTMVStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateMTMVStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateMTMVStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseAlterMTMVStmt(t *testing.T, sql string) *ast.AlterMTMVStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterMTMVStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterMTMVStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseDropMTMVStmt(t *testing.T, sql string) *ast.DropMTMVStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropMTMVStmt)
	if !ok {
		t.Fatalf("expected *ast.DropMTMVStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseRefreshMTMVStmt(t *testing.T, sql string) *ast.RefreshMTMVStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
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

func TestCreateMTMV_QualifiedName(t *testing.T) {
	stmt := parseCreateMTMVStmt(t, "CREATE MATERIALIZED VIEW db1.mv1 AS SELECT * FROM t")
	if stmt.Name.String() != "db1.mv1" {
		t.Errorf("Name = %q, want db1.mv1", stmt.Name.String())
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
	// REFRESH COMPLETE is Doris's; StarRocks.g4 refreshSchemeDesc takes
	// (IMMEDIATE | DEFERRED)? ASYNC [...] | INCREMENTAL | MANUAL.
	if _, errs := parseForTest("ALTER MATERIALIZED VIEW partition_mv REFRESH COMPLETE"); len(errs) == 0 {
		t.Error("REFRESH COMPLETE parsed; want a syntax error")
	}
	for sql, want := range map[string]string{
		"ALTER MATERIALIZED VIEW mv REFRESH ASYNC":                                             "ASYNC",
		"ALTER MATERIALIZED VIEW mv REFRESH DEFERRED ASYNC":                                    "DEFERRED ASYNC",
		"ALTER MATERIALIZED VIEW mv REFRESH ASYNC START ('2024-01-01') EVERY (INTERVAL 1 DAY)": "ASYNC",
		"ALTER MATERIALIZED VIEW mv REFRESH MANUAL":                                            "MANUAL",
	} {
		stmt := parseAlterMTMVStmt(t, sql)
		if stmt.RefreshMethod != want {
			t.Errorf("%q: RefreshMethod = %q, want %q", sql, stmt.RefreshMethod, want)
		}
	}
}

func TestAlterMTMV_RefreshAsyncSchedule(t *testing.T) {
	stmt := parseAlterMTMVStmt(t, "ALTER MATERIALIZED VIEW mv REFRESH ASYNC START ('2024-01-01') EVERY (INTERVAL 1 DAY)")
	tr := stmt.RefreshTrigger
	if tr == nil {
		t.Fatal("RefreshTrigger is nil")
	}
	if !tr.OnSchedule || tr.StartsAt != "2024-01-01" || tr.Interval != "1 DAY" {
		t.Errorf("RefreshTrigger = %+v, want schedule from 2024-01-01 every 1 DAY", tr)
	}
	if stmt := parseAlterMTMVStmt(t, "ALTER MATERIALIZED VIEW mv REFRESH ASYNC"); stmt.RefreshTrigger != nil {
		t.Errorf("REFRESH ASYNC: RefreshTrigger = %+v, want nil", stmt.RefreshTrigger)
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
	// REPLACE WITH MATERIALIZED VIEW is Doris's; StarRocks swaps with
	// SWAP WITH (StarRocks.g4 swapTableClause).
	if _, errs := parseForTest("ALTER MATERIALIZED VIEW mv REPLACE WITH MATERIALIZED VIEW mv2"); len(errs) == 0 {
		t.Error("REPLACE WITH parsed; want a syntax error")
	}
	stmt := parseAlterMTMVStmt(t, "ALTER MATERIALIZED VIEW mv SWAP WITH mv2")
	if stmt.SwapWith == nil || stmt.SwapWith.String() != "mv2" {
		t.Errorf("SwapWith = %v, want mv2", stmt.SwapWith)
	}
	stmt = parseAlterMTMVStmt(t, "ALTER MATERIALIZED VIEW mv INACTIVE")
	if stmt.Status != "INACTIVE" {
		t.Errorf("Status = %q, want INACTIVE", stmt.Status)
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

// TestMTMVGrammar holds CREATE and REFRESH MATERIALIZED VIEW to StarRocks.g4
// (:646-686); the rejected forms are Doris's, which this file once pinned.
func TestMTMVGrammar(t *testing.T) {
	cases := []struct {
		sql string
		ok  bool
	}{
		{"CREATE MATERIALIZED VIEW mv (a COMMENT 'c', b, INDEX i (a) USING BITMAP COMMENT 'x') COMMENT 'm' DISTRIBUTED BY HASH(a) BUCKETS 2 ORDER BY (a) REFRESH DEFERRED ASYNC START('2024-12-01 20:30:00') EVERY (INTERVAL 1 DAY) PARTITION BY date_trunc('day', a) PROPERTIES ('k' = 'v') AS SELECT a, b FROM t", true},
		{"CREATE MATERIALIZED VIEW mv REFRESH MANUAL PARTITION BY (a, b) AS SELECT a, b FROM t", true},
		{"CREATE MATERIALIZED VIEW mv (INDEX i (a)) AS SELECT a FROM t", false},
		{"CREATE MATERIALIZED VIEW mv (INDEX i (a), b) AS SELECT a FROM t", false},
		{"CREATE MATERIALIZED VIEW mv REFRESH ASYNC zzz AS SELECT 1", false},
		{"CREATE MATERIALIZED VIEW mv ASYNC AS SELECT 1", false},
		{"CREATE MATERIALIZED VIEW mv PARTITION BY CURRENT_USER AS SELECT 1", false},
		{"CREATE MATERIALIZED VIEW mv SELECT 1", false},
		{"CREATE MATERIALIZED VIEW mv BUILD IMMEDIATE REFRESH AUTO AS SELECT id FROM t", false},
		{"CREATE MATERIALIZED VIEW mv REFRESH COMPLETE ON SCHEDULE EVERY 1 DAY AS SELECT * FROM t", false},
		{"REFRESH MATERIALIZED VIEW mv PARTITION START('a') END('b') FORCE WITH SYNC MODE WITH PRIORITY 1", true},
		{"REFRESH MATERIALIZED VIEW mv PARTITION (('a', NULL), ('b'))", true},
		{"REFRESH MATERIALIZED VIEW mv PARTITION ((('a')))", false},
		{"REFRESH MATERIALIZED VIEW mv WITH ASYNC MODE", true},
		{"REFRESH MATERIALIZED VIEW mv1 AUTO", false},
		{"REFRESH MATERIALIZED VIEW mv1 PARTITIONS(p1, p2)", false},
	}
	for _, c := range cases {
		_, err := Parse(c.sql)
		if (err == nil) != c.ok {
			t.Errorf("Parse(%q): err = %v, want accepted=%v", c.sql, err, c.ok)
		}
	}
}
