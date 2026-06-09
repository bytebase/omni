package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func parseBackupStmt(t *testing.T, sql string) *ast.BackupStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.BackupStmt)
	if !ok {
		t.Fatalf("expected *ast.BackupStmt, got %T", n)
	}
	return stmt
}

func parseRestoreStmt(t *testing.T, sql string) *ast.RestoreStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.RestoreStmt)
	if !ok {
		t.Fatalf("expected *ast.RestoreStmt, got %T", n)
	}
	return stmt
}

func parseKillStmt(t *testing.T, sql string) *ast.KillStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.KillStmt)
	if !ok {
		t.Fatalf("expected *ast.KillStmt, got %T", n)
	}
	return stmt
}

func parseLockTablesStmt(t *testing.T, sql string) *ast.LockTablesStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.LockTablesStmt)
	if !ok {
		t.Fatalf("expected *ast.LockTablesStmt, got %T", n)
	}
	return stmt
}

func parseUnlockTablesStmt(t *testing.T, sql string) *ast.UnlockTablesStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.UnlockTablesStmt)
	if !ok {
		t.Fatalf("expected *ast.UnlockTablesStmt, got %T", n)
	}
	return stmt
}

func parseInstallPluginStmt(t *testing.T, sql string) *ast.InstallPluginStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.InstallPluginStmt)
	if !ok {
		t.Fatalf("expected *ast.InstallPluginStmt, got %T", n)
	}
	return stmt
}

func parseUninstallPluginStmt(t *testing.T, sql string) *ast.UninstallPluginStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.UninstallPluginStmt)
	if !ok {
		t.Fatalf("expected *ast.UninstallPluginStmt, got %T", n)
	}
	return stmt
}

func parseWarmUpStmt(t *testing.T, sql string) *ast.WarmUpStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.WarmUpStmt)
	if !ok {
		t.Fatalf("expected *ast.WarmUpStmt, got %T", n)
	}
	return stmt
}

func parseCleanStmt(t *testing.T, sql string) *ast.CleanStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.CleanStmt)
	if !ok {
		t.Fatalf("expected *ast.CleanStmt, got %T", n)
	}
	return stmt
}

func parseCancelStmt(t *testing.T, sql string) *ast.CancelStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.CancelStmt)
	if !ok {
		t.Fatalf("expected *ast.CancelStmt, got %T", n)
	}
	return stmt
}

func parseRecoverStmt(t *testing.T, sql string) *ast.RecoverStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.RecoverStmt)
	if !ok {
		t.Fatalf("expected *ast.RecoverStmt, got %T", n)
	}
	return stmt
}

// ---------------------------------------------------------------------------
// BACKUP tests
// ---------------------------------------------------------------------------

func TestBackupBasic(t *testing.T) {
	stmt := parseBackupStmt(t, `BACKUP SNAPSHOT example_db.snapshot_label1 TO example_repo`)
	if stmt.Label != "example_db.snapshot_label1" {
		t.Errorf("Label: got %q, want %q", stmt.Label, "example_db.snapshot_label1")
	}
	if stmt.Repo != "example_repo" {
		t.Errorf("Repo: got %q, want %q", stmt.Repo, "example_repo")
	}
	if len(stmt.Tables) != 0 {
		t.Errorf("Tables: got %d, want 0", len(stmt.Tables))
	}
}

func TestBackupWithTables(t *testing.T) {
	stmt := parseBackupStmt(t,
		`BACKUP SNAPSHOT mydb.snap1 TO my_repo ON (table1, table2)`)
	if stmt.Label != "mydb.snap1" {
		t.Errorf("Label: got %q", stmt.Label)
	}
	if stmt.Repo != "my_repo" {
		t.Errorf("Repo: got %q", stmt.Repo)
	}
	if len(stmt.Tables) != 2 {
		t.Errorf("Tables: got %d, want 2", len(stmt.Tables))
	}
	if stmt.Tables[0].String() != "table1" {
		t.Errorf("Tables[0]: got %q", stmt.Tables[0].String())
	}
}

func TestBackupWithProperties(t *testing.T) {
	stmt := parseBackupStmt(t,
		`BACKUP SNAPSHOT db.snap TO repo PROPERTIES("type"="full")`)
	if len(stmt.Properties) != 1 {
		t.Fatalf("Properties: got %d, want 1", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "type" {
		t.Errorf("Properties[0].Key: got %q", stmt.Properties[0].Key)
	}
	if stmt.Properties[0].Value != "full" {
		t.Errorf("Properties[0].Value: got %q", stmt.Properties[0].Value)
	}
}

func TestBackupWithTablesAndPartitions(t *testing.T) {
	stmt := parseBackupStmt(t,
		`BACKUP SNAPSHOT db.snap TO repo ON (t1 PARTITION(p1, p2), t2)`)
	if len(stmt.Tables) != 2 {
		t.Errorf("Tables: got %d, want 2", len(stmt.Tables))
	}
}

func TestBackupLoc(t *testing.T) {
	sql := `BACKUP SNAPSHOT db.snap TO repo`
	stmt := parseBackupStmt(t, sql)
	if !stmt.Loc.IsValid() {
		t.Error("Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// RESTORE tests
// ---------------------------------------------------------------------------

func TestRestoreBasic(t *testing.T) {
	stmt := parseRestoreStmt(t, `RESTORE SNAPSHOT example_db.snapshot1 FROM example_repo`)
	if stmt.Label != "example_db.snapshot1" {
		t.Errorf("Label: got %q", stmt.Label)
	}
	if stmt.Repo != "example_repo" {
		t.Errorf("Repo: got %q", stmt.Repo)
	}
}

func TestRestoreWithTablesAndProperties(t *testing.T) {
	stmt := parseRestoreStmt(t,
		`RESTORE SNAPSHOT db.snap FROM repo ON (tbl) PROPERTIES("timeout"="3600")`)
	if len(stmt.Tables) != 1 {
		t.Errorf("Tables: got %d, want 1", len(stmt.Tables))
	}
	if len(stmt.Properties) != 1 {
		t.Errorf("Properties: got %d, want 1", len(stmt.Properties))
	}
}

func TestRestoreLabelParts(t *testing.T) {
	stmt := parseRestoreStmt(t, `RESTORE SNAPSHOT mysnap FROM repo`)
	if len(stmt.LabelParts) != 1 {
		t.Errorf("LabelParts: got %d, want 1", len(stmt.LabelParts))
	}
	if stmt.LabelParts[0] != "mysnap" {
		t.Errorf("LabelParts[0]: got %q", stmt.LabelParts[0])
	}
}

// ---------------------------------------------------------------------------
// KILL tests
// ---------------------------------------------------------------------------

func TestKillConnectionId(t *testing.T) {
	stmt := parseKillStmt(t, `KILL 12345`)
	if stmt.Kind != "" {
		t.Errorf("Kind: got %q, want empty", stmt.Kind)
	}
	if stmt.Target != "12345" {
		t.Errorf("Target: got %q, want 12345", stmt.Target)
	}
}

func TestKillWithConnectionKeyword(t *testing.T) {
	stmt := parseKillStmt(t, `KILL CONNECTION 99`)
	if stmt.Kind != "CONNECTION" {
		t.Errorf("Kind: got %q, want CONNECTION", stmt.Kind)
	}
	if stmt.Target != "99" {
		t.Errorf("Target: got %q, want 99", stmt.Target)
	}
}

func TestKillQuery(t *testing.T) {
	stmt := parseKillStmt(t, `KILL QUERY "query-id-abc"`)
	if stmt.Kind != "QUERY" {
		t.Errorf("Kind: got %q, want QUERY", stmt.Kind)
	}
	if stmt.Target != "query-id-abc" {
		t.Errorf("Target: got %q", stmt.Target)
	}
}

func TestKillLoc(t *testing.T) {
	stmt := parseKillStmt(t, `KILL 1`)
	if !stmt.Loc.IsValid() {
		t.Error("Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// LOCK TABLES tests
// ---------------------------------------------------------------------------

func TestLockTablesRead(t *testing.T) {
	stmt := parseLockTablesStmt(t, `LOCK TABLES t1 READ`)
	if len(stmt.Items) != 1 {
		t.Fatalf("Items: got %d, want 1", len(stmt.Items))
	}
	if stmt.Items[0].Table.String() != "t1" {
		t.Errorf("Items[0].Table: got %q", stmt.Items[0].Table.String())
	}
	if stmt.Items[0].Mode != "READ" {
		t.Errorf("Items[0].Mode: got %q, want READ", stmt.Items[0].Mode)
	}
}

func TestLockTablesWrite(t *testing.T) {
	stmt := parseLockTablesStmt(t, `LOCK TABLES t1 WRITE`)
	if stmt.Items[0].Mode != "WRITE" {
		t.Errorf("Mode: got %q, want WRITE", stmt.Items[0].Mode)
	}
}

func TestLockTablesReadLocal(t *testing.T) {
	stmt := parseLockTablesStmt(t, `LOCK TABLES t1 READ LOCAL`)
	if stmt.Items[0].Mode != "READ LOCAL" {
		t.Errorf("Mode: got %q, want READ LOCAL", stmt.Items[0].Mode)
	}
}

func TestLockTablesLowPriorityWrite(t *testing.T) {
	stmt := parseLockTablesStmt(t, `LOCK TABLES t1 LOW_PRIORITY WRITE`)
	if stmt.Items[0].Mode != "LOW_PRIORITY WRITE" {
		t.Errorf("Mode: got %q, want LOW_PRIORITY WRITE", stmt.Items[0].Mode)
	}
}

func TestLockTablesMultiple(t *testing.T) {
	stmt := parseLockTablesStmt(t, `LOCK TABLES t1 READ, t2 WRITE`)
	if len(stmt.Items) != 2 {
		t.Fatalf("Items: got %d, want 2", len(stmt.Items))
	}
	if stmt.Items[1].Table.String() != "t2" {
		t.Errorf("Items[1].Table: got %q", stmt.Items[1].Table.String())
	}
}

func TestLockTablesWithAlias(t *testing.T) {
	stmt := parseLockTablesStmt(t, `LOCK TABLES t1 AS t WRITE`)
	if stmt.Items[0].Alias != "t" {
		t.Errorf("Alias: got %q, want t", stmt.Items[0].Alias)
	}
	if stmt.Items[0].Mode != "WRITE" {
		t.Errorf("Mode: got %q, want WRITE", stmt.Items[0].Mode)
	}
}

func TestUnlockTables(t *testing.T) {
	stmt := parseUnlockTablesStmt(t, `UNLOCK TABLES`)
	if !stmt.Loc.IsValid() {
		t.Error("Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// INSTALL / UNINSTALL PLUGIN tests (from legacy corpus)
// ---------------------------------------------------------------------------

func TestInstallPluginFromPath(t *testing.T) {
	stmt := parseInstallPluginStmt(t, `INSTALL PLUGIN FROM "/home/users/doris/auditdemo.zip"`)
	if stmt.Source != "/home/users/doris/auditdemo.zip" {
		t.Errorf("Source: got %q", stmt.Source)
	}
	if stmt.IsSoname {
		t.Error("IsSoname should be false")
	}
}

func TestInstallPluginFromDir(t *testing.T) {
	stmt := parseInstallPluginStmt(t, `INSTALL PLUGIN FROM "/home/users/doris/auditdemo/"`)
	if stmt.Source != "/home/users/doris/auditdemo/" {
		t.Errorf("Source: got %q", stmt.Source)
	}
}

func TestInstallPluginFromURL(t *testing.T) {
	stmt := parseInstallPluginStmt(t, `INSTALL PLUGIN FROM "http://mywebsite.com/plugin.zip"`)
	if stmt.Source != "http://mywebsite.com/plugin.zip" {
		t.Errorf("Source: got %q", stmt.Source)
	}
}

func TestInstallPluginFromURLWithProperties(t *testing.T) {
	stmt := parseInstallPluginStmt(t,
		`INSTALL PLUGIN FROM "http://mywebsite.com/plugin.zip" PROPERTIES("md5sum" = "73877f6029216f4314d712086a146570")`)
	if stmt.Source != "http://mywebsite.com/plugin.zip" {
		t.Errorf("Source: got %q", stmt.Source)
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("Properties: got %d, want 1", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "md5sum" {
		t.Errorf("Properties[0].Key: got %q", stmt.Properties[0].Key)
	}
}

func TestInstallPluginSoname(t *testing.T) {
	stmt := parseInstallPluginStmt(t, `INSTALL PLUGIN FROM SONAME "mylib.so"`)
	if !stmt.IsSoname {
		t.Error("IsSoname should be true")
	}
	if stmt.Source != "mylib.so" {
		t.Errorf("Source: got %q", stmt.Source)
	}
}

func TestUninstallPlugin(t *testing.T) {
	stmt := parseUninstallPluginStmt(t, `UNINSTALL PLUGIN auditdemo`)
	if stmt.Name != "auditdemo" {
		t.Errorf("Name: got %q, want auditdemo", stmt.Name)
	}
}

// ---------------------------------------------------------------------------
// WARM UP tests
// ---------------------------------------------------------------------------

func TestWarmUpCluster(t *testing.T) {
	stmt := parseWarmUpStmt(t, `WARM UP CLUSTER cloud_cluster FROM cloud_cluster2`)
	if stmt.Verb != "CLUSTER" {
		t.Errorf("Verb: got %q, want CLUSTER", stmt.Verb)
	}
}

func TestWarmUpComputeGroup(t *testing.T) {
	stmt := parseWarmUpStmt(t, `WARM UP COMPUTE GROUP cg1 WITH TABLE t1`)
	if stmt.Verb != "COMPUTE GROUP" {
		t.Errorf("Verb: got %q, want COMPUTE GROUP", stmt.Verb)
	}
}

func TestWarmUpLoc(t *testing.T) {
	stmt := parseWarmUpStmt(t, `WARM UP CLUSTER c1 FROM c2`)
	if !stmt.Loc.IsValid() {
		t.Error("Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// CLEAN tests
// ---------------------------------------------------------------------------

func TestCleanAllProfile(t *testing.T) {
	stmt := parseCleanStmt(t, `CLEAN ALL PROFILE`)
	if stmt.Verb != "ALL PROFILE" {
		t.Errorf("Verb: got %q, want ALL PROFILE", stmt.Verb)
	}
}

func TestCleanLabel(t *testing.T) {
	stmt := parseCleanStmt(t, `CLEAN LABEL my_label FROM my_db`)
	if stmt.Verb != "LABEL" {
		t.Errorf("Verb: got %q, want LABEL", stmt.Verb)
	}
}

func TestCleanQueryStats(t *testing.T) {
	stmt := parseCleanStmt(t, `CLEAN QUERY STATS FROM my_db`)
	if stmt.Verb != "QUERY STATS" {
		t.Errorf("Verb: got %q, want QUERY STATS", stmt.Verb)
	}
}

func TestCleanQueryStatsAll(t *testing.T) {
	stmt := parseCleanStmt(t, `CLEAN QUERY STATS ALL`)
	if stmt.Verb != "QUERY STATS" {
		t.Errorf("Verb: got %q, want QUERY STATS", stmt.Verb)
	}
}

func TestCleanLoc(t *testing.T) {
	stmt := parseCleanStmt(t, `CLEAN ALL PROFILE`)
	if !stmt.Loc.IsValid() {
		t.Error("Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// CANCEL (generic) tests
// ---------------------------------------------------------------------------

func TestCancelLoad(t *testing.T) {
	stmt := parseCancelStmt(t, `CANCEL LOAD FROM mydb WHERE label = "my_load"`)
	if stmt.Verb != "LOAD" {
		t.Errorf("Verb: got %q, want LOAD", stmt.Verb)
	}
}

func TestCancelExport(t *testing.T) {
	stmt := parseCancelStmt(t, `CANCEL EXPORT FROM mydb WHERE queryid = "abc"`)
	if stmt.Verb != "EXPORT" {
		t.Errorf("Verb: got %q, want EXPORT", stmt.Verb)
	}
}

func TestCancelBackup(t *testing.T) {
	stmt := parseCancelStmt(t, `CANCEL BACKUP FROM example_db`)
	if stmt.Verb != "BACKUP" {
		t.Errorf("Verb: got %q, want BACKUP", stmt.Verb)
	}
}

func TestCancelRestore(t *testing.T) {
	stmt := parseCancelStmt(t, `CANCEL RESTORE FROM example_db`)
	if stmt.Verb != "RESTORE" {
		t.Errorf("Verb: got %q, want RESTORE", stmt.Verb)
	}
}

func TestCancelAlterTable(t *testing.T) {
	stmt := parseCancelStmt(t, `CANCEL ALTER TABLE COLUMN FROM example_tbl`)
	if stmt.Verb != "ALTER TABLE" {
		t.Errorf("Verb: got %q, want ALTER TABLE", stmt.Verb)
	}
}

func TestCancelAlterTableRollup(t *testing.T) {
	stmt := parseCancelStmt(t, `CANCEL ALTER TABLE ROLLUP FROM example_tbl`)
	if stmt.Verb != "ALTER TABLE" {
		t.Errorf("Verb: got %q, want ALTER TABLE", stmt.Verb)
	}
}

func TestCancelBuildIndex(t *testing.T) {
	stmt := parseCancelStmt(t, `CANCEL BUILD INDEX ON example_tbl`)
	if stmt.Verb != "BUILD INDEX" {
		t.Errorf("Verb: got %q, want BUILD INDEX", stmt.Verb)
	}
}

func TestCancelLoc(t *testing.T) {
	stmt := parseCancelStmt(t, `CANCEL LOAD FROM db WHERE label = "x"`)
	if !stmt.Loc.IsValid() {
		t.Error("Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// RECOVER tests (from legacy corpus)
// ---------------------------------------------------------------------------

func TestRecoverDatabase(t *testing.T) {
	stmt := parseRecoverStmt(t, `RECOVER DATABASE example_db`)
	if stmt.Verb != "DATABASE" {
		t.Errorf("Verb: got %q, want DATABASE", stmt.Verb)
	}
	if stmt.Name == nil || stmt.Name.String() != "example_db" {
		t.Errorf("Name: got %v", stmt.Name)
	}
	if stmt.ID != "" {
		t.Errorf("ID: expected empty, got %q", stmt.ID)
	}
}

func TestRecoverTable(t *testing.T) {
	stmt := parseRecoverStmt(t, `RECOVER TABLE example_db.example_tbl`)
	if stmt.Verb != "TABLE" {
		t.Errorf("Verb: got %q, want TABLE", stmt.Verb)
	}
	if stmt.Name == nil || stmt.Name.String() != "example_db.example_tbl" {
		t.Errorf("Name: got %v", stmt.Name)
	}
}

func TestRecoverPartition(t *testing.T) {
	stmt := parseRecoverStmt(t, `RECOVER PARTITION p1 FROM example_tbl`)
	if stmt.Verb != "PARTITION" {
		t.Errorf("Verb: got %q, want PARTITION", stmt.Verb)
	}
	if stmt.Name == nil || stmt.Name.String() != "p1" {
		t.Errorf("Name: got %v", stmt.Name)
	}
	if stmt.FromTable == nil || stmt.FromTable.String() != "example_tbl" {
		t.Errorf("FromTable: got %v", stmt.FromTable)
	}
}

func TestRecoverDatabaseWithID(t *testing.T) {
	stmt := parseRecoverStmt(t, `RECOVER DATABASE example_db 12345`)
	if stmt.ID != "12345" {
		t.Errorf("ID: got %q, want 12345", stmt.ID)
	}
}

func TestRecoverTableWithID(t *testing.T) {
	stmt := parseRecoverStmt(t, `RECOVER TABLE example_db.example_tbl 12346`)
	if stmt.ID != "12346" {
		t.Errorf("ID: got %q, want 12346", stmt.ID)
	}
}

func TestRecoverPartitionWithID(t *testing.T) {
	stmt := parseRecoverStmt(t, `RECOVER PARTITION p1 12347 FROM example_tbl`)
	if stmt.ID != "12347" {
		t.Errorf("ID: got %q, want 12347", stmt.ID)
	}
}

func TestRecoverDatabaseWithIDAndNewName(t *testing.T) {
	stmt := parseRecoverStmt(t, `RECOVER DATABASE example_db 12345 AS new_example_db`)
	if stmt.ID != "12345" {
		t.Errorf("ID: got %q, want 12345", stmt.ID)
	}
	if stmt.NewName != "new_example_db" {
		t.Errorf("NewName: got %q, want new_example_db", stmt.NewName)
	}
}

func TestRecoverTableWithNewName(t *testing.T) {
	stmt := parseRecoverStmt(t, `RECOVER TABLE example_db.example_tbl AS new_example_tbl`)
	if stmt.NewName != "new_example_tbl" {
		t.Errorf("NewName: got %q, want new_example_tbl", stmt.NewName)
	}
}

func TestRecoverPartitionWithIDAndNewName(t *testing.T) {
	stmt := parseRecoverStmt(t, `RECOVER PARTITION p1 12347 AS new_p1 FROM example_tbl`)
	if stmt.ID != "12347" {
		t.Errorf("ID: got %q, want 12347", stmt.ID)
	}
	if stmt.NewName != "new_p1" {
		t.Errorf("NewName: got %q, want new_p1", stmt.NewName)
	}
	if stmt.FromTable == nil || stmt.FromTable.String() != "example_tbl" {
		t.Errorf("FromTable: got %v", stmt.FromTable)
	}
}

func TestRecoverLoc(t *testing.T) {
	stmt := parseRecoverStmt(t, `RECOVER DATABASE mydb`)
	if !stmt.Loc.IsValid() {
		t.Error("Loc should be valid")
	}
}
