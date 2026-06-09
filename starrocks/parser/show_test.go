package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// SHOW TABLES
// ---------------------------------------------------------------------------

func TestShowTables(t *testing.T) {
	file, errs := Parse("SHOW TABLES")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.ShowStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowStmt, got %T", file.Stmts[0])
	}
	if n.Type != "TABLES" {
		t.Errorf("Type = %q, want TABLES", n.Type)
	}
	if n.Full {
		t.Error("Full should be false")
	}
}

func TestShowTablesLike(t *testing.T) {
	file, errs := Parse("SHOW TABLES LIKE '%cm%'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "TABLES" {
		t.Errorf("Type = %q, want TABLES", n.Type)
	}
	if n.Like != "%cm%" {
		t.Errorf("Like = %q, want %%cm%%", n.Like)
	}
}

func TestShowFullTables(t *testing.T) {
	file, errs := Parse("SHOW FULL TABLES")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if !n.Full {
		t.Error("Full should be true")
	}
	if n.Type != "TABLES" {
		t.Errorf("Type = %q, want TABLES", n.Type)
	}
}

func TestShowTablesFromDB(t *testing.T) {
	file, errs := Parse("SHOW TABLES FROM mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.From != "mydb" {
		t.Errorf("From = %q, want mydb", n.From)
	}
}

// ---------------------------------------------------------------------------
// SHOW DATABASES
// ---------------------------------------------------------------------------

func TestShowDatabases(t *testing.T) {
	file, errs := Parse("SHOW DATABASES")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "DATABASES" {
		t.Errorf("Type = %q, want DATABASES", n.Type)
	}
}

func TestShowDatabasesLike(t *testing.T) {
	file, errs := Parse("SHOW DATABASES LIKE 'test%'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Like != "test%" {
		t.Errorf("Like = %q, want test%%", n.Like)
	}
}

// ---------------------------------------------------------------------------
// SHOW COLUMNS
// ---------------------------------------------------------------------------

func TestShowColumns(t *testing.T) {
	file, errs := Parse("SHOW COLUMNS FROM t_agg")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "COLUMNS" {
		t.Errorf("Type = %q, want COLUMNS", n.Type)
	}
	if n.Target == nil || len(n.Target.Parts) == 0 || n.Target.Parts[0] != "t_agg" {
		t.Errorf("Target = %v, want t_agg", n.Target)
	}
}

func TestShowFullColumns(t *testing.T) {
	file, errs := Parse("SHOW FULL COLUMNS FROM t_agg")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if !n.Full {
		t.Error("Full should be true")
	}
	if n.Type != "COLUMNS" {
		t.Errorf("Type = %q, want COLUMNS", n.Type)
	}
}

// ---------------------------------------------------------------------------
// SHOW CREATE TABLE / VIEW / DATABASE / CATALOG
// ---------------------------------------------------------------------------

func TestShowCreateTable(t *testing.T) {
	file, errs := Parse("SHOW CREATE TABLE demo.test_table")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "CREATE TABLE" {
		t.Errorf("Type = %q, want CREATE TABLE", n.Type)
	}
	if n.Target == nil {
		t.Fatal("Target is nil")
	}
	if len(n.Target.Parts) != 2 {
		t.Errorf("Target.Parts = %v, want [demo test_table]", n.Target.Parts)
	}
}

func TestShowBriefCreateTable(t *testing.T) {
	file, errs := Parse("SHOW BRIEF CREATE TABLE demo.test_table")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "BRIEF CREATE TABLE" {
		t.Errorf("Type = %q, want BRIEF CREATE TABLE", n.Type)
	}
}

func TestShowCreateView(t *testing.T) {
	file, errs := Parse("SHOW CREATE VIEW my_view")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "CREATE VIEW" {
		t.Errorf("Type = %q, want CREATE VIEW", n.Type)
	}
}

func TestShowCreateDatabase(t *testing.T) {
	file, errs := Parse("SHOW CREATE DATABASE mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "CREATE DATABASE" {
		t.Errorf("Type = %q, want CREATE DATABASE", n.Type)
	}
}

func TestShowCreateCatalog(t *testing.T) {
	file, errs := Parse("SHOW CREATE CATALOG oracle")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "CREATE CATALOG" {
		t.Errorf("Type = %q, want CREATE CATALOG", n.Type)
	}
	if n.Target == nil || n.Target.Parts[0] != "oracle" {
		t.Errorf("Target = %v, want oracle", n.Target)
	}
}

// ---------------------------------------------------------------------------
// SHOW VARIABLES
// ---------------------------------------------------------------------------

func TestShowVariablesLike(t *testing.T) {
	file, errs := Parse("SHOW VARIABLES LIKE 'max_connections'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "VARIABLES" {
		t.Errorf("Type = %q, want VARIABLES", n.Type)
	}
	if n.Like != "max_connections" {
		t.Errorf("Like = %q, want max_connections", n.Like)
	}
}

func TestShowVariablesWhere(t *testing.T) {
	file, errs := Parse("SHOW VARIABLES WHERE variable_name = 'version'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Where == nil {
		t.Error("Where should not be nil")
	}
}

// ---------------------------------------------------------------------------
// SHOW PARTITIONS
// ---------------------------------------------------------------------------

func TestShowPartitions(t *testing.T) {
	file, errs := Parse("SHOW PARTITIONS FROM t_agg")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "PARTITIONS" {
		t.Errorf("Type = %q, want PARTITIONS", n.Type)
	}
	if n.Target == nil || n.Target.Parts[0] != "t_agg" {
		t.Errorf("Target = %v, want t_agg", n.Target)
	}
}

func TestShowTemporaryPartitions(t *testing.T) {
	file, errs := Parse("SHOW TEMPORARY PARTITIONS FROM t_temp")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "PARTITIONS" {
		t.Errorf("Type = %q, want PARTITIONS", n.Type)
	}
	if n.Args != "TEMPORARY" {
		t.Errorf("Args = %q, want TEMPORARY", n.Args)
	}
}

func TestShowPartitionsWhere(t *testing.T) {
	file, errs := Parse(`SHOW PARTITIONS FROM t_agg WHERE PartitionName = "p2024"`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Target == nil {
		t.Fatal("Target is nil")
	}
	if n.Args == "" {
		t.Error("Args should contain WHERE clause")
	}
}

// ---------------------------------------------------------------------------
// SHOW GRANTS / ROLES
// ---------------------------------------------------------------------------

func TestShowGrants(t *testing.T) {
	file, errs := Parse("SHOW GRANTS")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "GRANTS" {
		t.Errorf("Type = %q, want GRANTS", n.Type)
	}
}

func TestShowGrantsForUser(t *testing.T) {
	file, errs := Parse("SHOW GRANTS FOR jack@'%'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "GRANTS" {
		t.Errorf("Type = %q, want GRANTS", n.Type)
	}
	if n.Args == "" {
		t.Error("Args should contain FOR user")
	}
}

func TestShowAllGrants(t *testing.T) {
	file, errs := Parse("SHOW ALL GRANTS")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "GRANTS" {
		t.Errorf("Type = %q, want GRANTS", n.Type)
	}
	if n.Args != "ALL" {
		t.Errorf("Args = %q, want ALL", n.Args)
	}
}

func TestShowRoles(t *testing.T) {
	file, errs := Parse("SHOW ROLES")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "ROLES" {
		t.Errorf("Type = %q, want ROLES", n.Type)
	}
}

// ---------------------------------------------------------------------------
// SHOW CATALOGS
// ---------------------------------------------------------------------------

func TestShowCatalogs(t *testing.T) {
	file, errs := Parse("SHOW CATALOGS")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "CATALOGS" {
		t.Errorf("Type = %q, want CATALOGS", n.Type)
	}
}

func TestShowCatalogsLike(t *testing.T) {
	file, errs := Parse("SHOW CATALOGS LIKE 'hi%'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Like != "hi%" {
		t.Errorf("Like = %q, want hi%%", n.Like)
	}
}

// ---------------------------------------------------------------------------
// SHOW TABLE STATUS
// ---------------------------------------------------------------------------

func TestShowTableStatus(t *testing.T) {
	file, errs := Parse("SHOW TABLE STATUS")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "TABLE STATUS" {
		t.Errorf("Type = %q, want TABLE STATUS", n.Type)
	}
}

func TestShowTableStatusFromDB(t *testing.T) {
	file, errs := Parse(`SHOW TABLE STATUS FROM db LIKE "%test%"`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "TABLE STATUS" {
		t.Errorf("Type = %q, want TABLE STATUS", n.Type)
	}
	if n.From != "db" {
		t.Errorf("From = %q, want db", n.From)
	}
	if n.Like != "%test%" {
		t.Errorf("Like = %q, want %%test%%", n.Like)
	}
}

// ---------------------------------------------------------------------------
// SHOW ALTER TABLE
// ---------------------------------------------------------------------------

func TestShowAlterTableColumn(t *testing.T) {
	file, errs := Parse("SHOW ALTER TABLE COLUMN")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "ALTER TABLE COLUMN" {
		t.Errorf("Type = %q, want ALTER TABLE COLUMN", n.Type)
	}
}

func TestShowAlterTableColumnWhere(t *testing.T) {
	file, errs := Parse(`SHOW ALTER TABLE COLUMN WHERE TableName = "table1" ORDER BY CreateTime DESC LIMIT 1`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "ALTER TABLE COLUMN" {
		t.Errorf("Type = %q, want ALTER TABLE COLUMN", n.Type)
	}
	if n.Args == "" {
		t.Error("Args should contain WHERE/ORDER BY/LIMIT clause")
	}
}

func TestShowAlterTableRollup(t *testing.T) {
	file, errs := Parse("SHOW ALTER TABLE ROLLUP FROM example_db")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "ALTER TABLE ROLLUP" {
		t.Errorf("Type = %q, want ALTER TABLE ROLLUP", n.Type)
	}
	if n.From != "example_db" {
		t.Errorf("From = %q, want example_db", n.From)
	}
}

// ---------------------------------------------------------------------------
// SHOW generic (PROCESSLIST, PROC, WARNINGS, ERRORS, etc.)
// ---------------------------------------------------------------------------

func TestShowProcesslist(t *testing.T) {
	file, errs := Parse("SHOW PROCESSLIST")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type == "" {
		t.Error("Type should not be empty")
	}
}

func TestShowWarnings(t *testing.T) {
	file, errs := Parse("SHOW WARNINGS")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type == "" {
		t.Error("Type should not be empty")
	}
}

// ---------------------------------------------------------------------------
// SHOW — Loc validity
// ---------------------------------------------------------------------------

func TestShowLocValid(t *testing.T) {
	input := "SHOW TABLES"
	file, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if !n.Loc.IsValid() {
		t.Error("ShowStmt.Loc should be valid")
	}
	if n.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", n.Loc.Start)
	}
}
