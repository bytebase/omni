package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// SHOW TABLES
// ---------------------------------------------------------------------------

func TestShowTables(t *testing.T) {
	file, errs := parseForTest("SHOW TABLES")
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
	file, errs := parseForTest("SHOW TABLES LIKE '%cm%'")
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
	file, errs := parseForTest("SHOW FULL TABLES")
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
	file, errs := parseForTest("SHOW TABLES FROM mydb")
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
	file, errs := parseForTest("SHOW DATABASES")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "DATABASES" {
		t.Errorf("Type = %q, want DATABASES", n.Type)
	}
}

func TestShowDatabasesLike(t *testing.T) {
	file, errs := parseForTest("SHOW DATABASES LIKE 'test%'")
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
	file, errs := parseForTest("SHOW COLUMNS FROM t_agg")
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
	file, errs := parseForTest("SHOW FULL COLUMNS FROM t_agg")
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
	file, errs := parseForTest("SHOW CREATE TABLE demo.test_table")
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

func TestShowCreateView(t *testing.T) {
	file, errs := parseForTest("SHOW CREATE VIEW my_view")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "CREATE VIEW" {
		t.Errorf("Type = %q, want CREATE VIEW", n.Type)
	}
}

func TestShowCreateDatabase(t *testing.T) {
	file, errs := parseForTest("SHOW CREATE DATABASE mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "CREATE DATABASE" {
		t.Errorf("Type = %q, want CREATE DATABASE", n.Type)
	}
}

func TestShowCreateCatalog(t *testing.T) {
	file, errs := parseForTest("SHOW CREATE CATALOG oracle")
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
	file, errs := parseForTest("SHOW VARIABLES LIKE 'max_connections'")
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
	file, errs := parseForTest("SHOW VARIABLES WHERE variable_name = 'version'")
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
	file, errs := parseForTest("SHOW PARTITIONS FROM t_agg")
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
	file, errs := parseForTest("SHOW TEMPORARY PARTITIONS FROM t_temp")
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
	file, errs := parseForTest(`SHOW PARTITIONS FROM t_agg WHERE PartitionName = "p2024"`)
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
	file, errs := parseForTest("SHOW GRANTS")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "GRANTS" {
		t.Errorf("Type = %q, want GRANTS", n.Type)
	}
}

func TestShowGrantsForUser(t *testing.T) {
	file, errs := parseForTest("SHOW GRANTS FOR jack@'%'")
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

func TestShowRoles(t *testing.T) {
	file, errs := parseForTest("SHOW ROLES")
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
	file, errs := parseForTest("SHOW CATALOGS")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "CATALOGS" {
		t.Errorf("Type = %q, want CATALOGS", n.Type)
	}
}

// ---------------------------------------------------------------------------
// SHOW TABLE STATUS
// ---------------------------------------------------------------------------

func TestShowTableStatus(t *testing.T) {
	file, errs := parseForTest("SHOW TABLE STATUS")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "TABLE STATUS" {
		t.Errorf("Type = %q, want TABLE STATUS", n.Type)
	}
}

func TestShowTableStatusFromDB(t *testing.T) {
	file, errs := parseForTest(`SHOW TABLE STATUS FROM db LIKE "%test%"`)
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
	file, errs := parseForTest("SHOW ALTER TABLE COLUMN")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "ALTER TABLE COLUMN" {
		t.Errorf("Type = %q, want ALTER TABLE COLUMN", n.Type)
	}
}

func TestShowAlterTableColumnWhere(t *testing.T) {
	file, errs := parseForTest(`SHOW ALTER TABLE COLUMN WHERE TableName = "table1" ORDER BY CreateTime DESC LIMIT 1`)
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
	file, errs := parseForTest("SHOW ALTER TABLE ROLLUP FROM example_db")
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
	file, errs := parseForTest("SHOW PROCESSLIST")
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
	file, errs := parseForTest("SHOW WARNINGS")
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
	file, errs := parseForTest(input)
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

func TestShowDatabasesFromCatalog(t *testing.T) {
	for _, sql := range []string{"SHOW DATABASES FROM hive_catalog", "SHOW DATABASES IN hive_catalog"} {
		file, errs := parseForTest(sql)
		if len(errs) != 0 {
			t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
		}
		n := file.Stmts[0].(*ast.ShowStmt)
		if n.Type != "DATABASES" || n.From != "hive_catalog" {
			t.Errorf("parseForTest(%q) = Type %q From %q, want DATABASES hive_catalog", sql, n.Type, n.From)
		}
	}
}

// TestShowStrictRejects holds Parse to the SHOW forms of StarRocks.g4: a form
// cut short, or a Doris form StarRocks does not have, is a syntax error.
func TestShowStrictRejects(t *testing.T) {
	for _, sql := range []string{
		"SHOW",
		"SHOW a",
		"SHOW PROCESSLIST zzz",
		"SHOW CREATE TABLE",
		"SHOW BRIEF CREATE TABLE demo.test_table",       // Doris only
		"SHOW ROUTINE LOAD LIKE 'job%'",                 // StarRocks.g4:1275 has no LIKE
		"SHOW CREATE ROUTINE LOAD FOR a",                // StarRocks.g4:1287 has no FOR
		"SHOW ALL GRANTS",                               // StarRocks.g4:1717
		"SHOW CATALOGS LIKE 'hi%'",                      // StarRocks.g4:773
		"SHOW WORKLOAD GROUPS",                          // Doris only
		"SHOW TABLETS FROM t PARTITION (p) LIMIT 1, ?x", // a cut limit
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) = nil error, want a syntax error", sql)
		}
	}
	for _, sql := range []string{
		"SHOW FULL TABLES FROM db LIKE 'a%'",
		"SHOW LOAD WHERE role = 'a'",
		"SHOW LOAD WHERE Label MATCH !Label",
		"SHOW LOAD WHERE GROUP_CONCAT(Label SEPARATOR Label) = 'a'",
		"SHOW LOAD WHERE !(Label = 'job')",
		"SHOW LOAD WHERE Label MATCH 'job' AND 5 MOD 2 = 1",
		"SHOW LOAD WHERE Label COLLATE utf8_general_ci = 'job'",
		"SHOW RESOURCE GROUP role",
		"SHOW FILE FROM role",
		"SHOW GRANTS FOR ROLE role",
		"SHOW TABLETS FROM t PARTITION (p1, p2) WHERE a = 1 ORDER BY a LIMIT 1, 2",
		"SHOW GRANTS FOR ROLE r",
		"SHOW GRANTS FOR ROLE",
		"SHOW GRANTS FOR 'u'@'%'",
		"SHOW AUTHENTICATION FOR 'u'@['h']",
		"SHOW ROUTINE LOAD TASK FROM db WHERE JobName = 'j'",
		"SHOW PIPES FROM db ORDER BY a LIMIT ?",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}

// TestMatchBindsAsPrimary holds MATCH to StarRocks.g4:2425, between primary
// expressions: it binds before a comparison.
func TestMatchBindsAsPrimary(t *testing.T) {
	file, err := Parse("SELECT Label MATCH 'job' = TRUE FROM t")
	if err != nil {
		t.Fatal(err)
	}
	sel := file.Stmts[0].(*ast.SelectStmt)
	cmp, ok := sel.Items[0].Expr.(*ast.BinaryExpr)
	if !ok || cmp.Op != ast.BinEq {
		t.Fatalf("got %T, want a comparison", sel.Items[0].Expr)
	}
	if _, ok := cmp.Left.(*ast.MatchExpr); !ok {
		t.Errorf("left of = is %T, want *ast.MatchExpr", cmp.Left)
	}
}

// TestShiftWordsBindLooserThanBitOr holds the shift words to StarRocks.g4:
// 2389-2394, where they come after | and so bind looser.
func TestShiftWordsBindLooserThanBitOr(t *testing.T) {
	file, err := Parse("SELECT 1 | 2 BITSHIFTLEFT 3")
	if err != nil {
		t.Fatal(err)
	}
	shift, ok := file.Stmts[0].(*ast.SelectStmt).Items[0].Expr.(*ast.BinaryExpr)
	if !ok || shift.Op != ast.BinShiftLeft {
		t.Fatalf("got %T, want a BITSHIFTLEFT at the top", file.Stmts[0].(*ast.SelectStmt).Items[0].Expr)
	}
	if or, ok := shift.Left.(*ast.BinaryExpr); !ok || or.Op != ast.BinBitOr {
		t.Errorf("left of BITSHIFTLEFT is %T, want 1 | 2", shift.Left)
	}
}

// TestUnaryBindsBeforeMatch holds MATCH below the unary operators of
// StarRocks.g4:2408-2409: !Label MATCH 'job' is (!Label) MATCH 'job'.
func TestUnaryBindsBeforeMatch(t *testing.T) {
	file, err := Parse("SELECT !Label MATCH 'job' FROM t")
	if err != nil {
		t.Fatal(err)
	}
	m, ok := file.Stmts[0].(*ast.SelectStmt).Items[0].Expr.(*ast.MatchExpr)
	if !ok {
		t.Fatalf("got %T, want *ast.MatchExpr at the top", file.Stmts[0].(*ast.SelectStmt).Items[0].Expr)
	}
	if _, ok := m.Expr.(*ast.UnaryExpr); !ok {
		t.Errorf("left of MATCH is %T, want !Label", m.Expr)
	}
}
