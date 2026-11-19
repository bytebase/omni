package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// CREATE CATALOG
// ---------------------------------------------------------------------------

func TestCreateCatalog_Basic(t *testing.T) {
	file, errs := parseForTest(`CREATE EXTERNAL CATALOG hive_catalog PROPERTIES("type"="hms", "hive.metastore.uris"="thrift://127.0.0.1:7004")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateCatalogStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateCatalogStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "hive_catalog" {
		t.Errorf("Name = %q, want %q", stmt.Name, "hive_catalog")
	}
	if !stmt.External {
		t.Error("External should be true")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "type" {
		t.Errorf("Properties[0].Key = %q, want %q", stmt.Properties[0].Key, "type")
	}
	if stmt.Properties[0].Value != "hms" {
		t.Errorf("Properties[0].Value = %q, want %q", stmt.Properties[0].Value, "hms")
	}
}

func TestCreateCatalog_External(t *testing.T) {
	file, errs := parseForTest(`CREATE EXTERNAL CATALOG iceberg_catalog PROPERTIES("type"="iceberg")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.CreateCatalogStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateCatalogStmt, got %T", file.Stmts[0])
	}
	if !stmt.External {
		t.Error("External should be true")
	}
	if stmt.Name != "iceberg_catalog" {
		t.Errorf("Name = %q, want %q", stmt.Name, "iceberg_catalog")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "type" || stmt.Properties[0].Value != "iceberg" {
		t.Errorf("unexpected property: key=%q val=%q", stmt.Properties[0].Key, stmt.Properties[0].Value)
	}
}

func TestCreateCatalog_IfNotExists(t *testing.T) {
	file, errs := parseForTest(`CREATE EXTERNAL CATALOG IF NOT EXISTS my_catalog PROPERTIES("type"="hms")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateCatalogStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if stmt.Name != "my_catalog" {
		t.Errorf("Name = %q, want %q", stmt.Name, "my_catalog")
	}
}

func TestCreateCatalog_WithComment(t *testing.T) {
	file, errs := parseForTest(`CREATE EXTERNAL CATALOG hive COMMENT 'hive catalog' PROPERTIES('type'='hms', 'hive.metastore.uris'='thrift://127.0.0.1:7004')`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateCatalogStmt)
	if stmt.Name != "hive" {
		t.Errorf("Name = %q, want %q", stmt.Name, "hive")
	}
	if stmt.Comment != "hive catalog" {
		t.Errorf("Comment = %q, want %q", stmt.Comment, "hive catalog")
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateCatalog_Refused(t *testing.T) {
	// createExternalCatalogStatement (StarRocks.g4:760) requires EXTERNAL and
	// a properties list, and has no WITH RESOURCE; the pinned server refuses
	// each of these.
	for _, sql := range []string{
		`CREATE CATALOG my_catalog PROPERTIES("type"="hms")`,
		`CREATE EXTERNAL CATALOG my_catalog`,
		`CREATE EXTERNAL CATALOG hms_catalog WITH RESOURCE hms_resource`,
		`CREATE EXTERNAL CATALOG hms_catalog PROPERTIES("type"="hms") WITH RESOURCE hms_resource`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestCreateCatalog_ES(t *testing.T) {
	file, errs := parseForTest(`CREATE EXTERNAL CATALOG es PROPERTIES ("type"="es", "hosts"="http://127.0.0.1:9200")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateCatalogStmt)
	if stmt.Name != "es" {
		t.Errorf("Name = %q, want %q", stmt.Name, "es")
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateCatalog_JDBC(t *testing.T) {
	input := `CREATE EXTERNAL CATALOG jdbc PROPERTIES (
		"type"="jdbc",
		"user"="root",
		"password"="123456",
		"jdbc_url" = "jdbc:mysql://127.0.0.1:3316/doris_test?useSSL=false",
		"driver_url" = "https://example.com/mysql-connector-java-8.0.25.jar",
		"driver_class" = "com.mysql.cj.jdbc.Driver"
	)`
	file, errs := parseForTest(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateCatalogStmt)
	if stmt.Name != "jdbc" {
		t.Errorf("Name = %q, want %q", stmt.Name, "jdbc")
	}
	if len(stmt.Properties) != 6 {
		t.Fatalf("expected 6 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateCatalog_Tag(t *testing.T) {
	file, errs := parseForTest(`CREATE EXTERNAL CATALOG c PROPERTIES("type"="hms")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateCatalogStmt {
		t.Errorf("Tag() = %v, want T_CreateCatalogStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// ALTER CATALOG
// ---------------------------------------------------------------------------

func TestAlterCatalog_Set(t *testing.T) {
	file, errs := parseForTest(`ALTER CATALOG hive SET ('hive.metastore.uris'='thrift://172.21.0.1:9083')`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.AlterCatalogStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterCatalogStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "hive" {
		t.Errorf("Name = %q, want %q", stmt.Name, "hive")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "hive.metastore.uris" {
		t.Errorf("Property key = %q, want %q", stmt.Properties[0].Key, "hive.metastore.uris")
	}
}

func TestAlterCatalog_Refused(t *testing.T) {
	// alterCatalogStatement (StarRocks.g4:2152) has one clause, SET with a
	// property list; the fork's RENAME, SET PROPERTIES, MODIFY COMMENT and
	// PROPERTY are refused by the pinned server.
	for _, sql := range []string{
		`ALTER CATALOG ctlg_hive RENAME hive`,
		`ALTER CATALOG hive SET PROPERTIES ('hive.metastore.uris'='thrift://172.21.0.1:9083')`,
		`ALTER CATALOG hive MODIFY COMMENT "new catalog comment"`,
		`ALTER CATALOG c PROPERTY ("key"="value")`,
		`ALTER CATALOG c SET 'key'='value'`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestAlterCatalog_Tag(t *testing.T) {
	file, errs := parseForTest(`ALTER CATALOG c SET ('key'='value')`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_AlterCatalogStmt {
		t.Errorf("Tag() = %v, want T_AlterCatalogStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// DROP CATALOG
// ---------------------------------------------------------------------------

func TestDropCatalog_Basic(t *testing.T) {
	file, errs := parseForTest(`DROP CATALOG hive`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropCatalogStmt)
	if !ok {
		t.Fatalf("expected *ast.DropCatalogStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "hive" {
		t.Errorf("Name = %q, want %q", stmt.Name, "hive")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropCatalog_IfExists(t *testing.T) {
	file, errs := parseForTest(`DROP CATALOG IF EXISTS hive`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.DropCatalogStmt)
	if !ok {
		t.Fatalf("expected *ast.DropCatalogStmt, got %T", file.Stmts[0])
	}
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if stmt.Name != "hive" {
		t.Errorf("Name = %q, want %q", stmt.Name, "hive")
	}
}

func TestDropCatalog_Tag(t *testing.T) {
	file, errs := parseForTest(`DROP CATALOG c`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_DropCatalogStmt {
		t.Errorf("Tag() = %v, want T_DropCatalogStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// REFRESH CATALOG
// ---------------------------------------------------------------------------

func TestRefreshCatalog_Basic(t *testing.T) {
	file, errs := parseForTest(`REFRESH CATALOG c`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.RefreshCatalogStmt)
	if !ok {
		t.Fatalf("expected *ast.RefreshCatalogStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "c" {
		t.Errorf("Name = %q, want %q", stmt.Name, "c")
	}
	if len(stmt.Properties) != 0 {
		t.Errorf("expected 0 properties, got %d", len(stmt.Properties))
	}
}

func TestRefreshCatalog_WithProperties(t *testing.T) {
	file, errs := parseForTest(`REFRESH CATALOG c PROPERTIES("key"="value")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.RefreshCatalogStmt)
	if stmt.Name != "c" {
		t.Errorf("Name = %q, want %q", stmt.Name, "c")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "key" || stmt.Properties[0].Value != "value" {
		t.Errorf("unexpected property: key=%q val=%q", stmt.Properties[0].Key, stmt.Properties[0].Value)
	}
}

func TestRefreshCatalog_Tag(t *testing.T) {
	file, errs := parseForTest(`REFRESH CATALOG c`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_RefreshCatalogStmt {
		t.Errorf("Tag() = %v, want T_RefreshCatalogStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// Legacy corpus — all catalog_create.sql examples
// ---------------------------------------------------------------------------

func TestCreateCatalog_Legacy_HiveWithHAFS(t *testing.T) {
	input := `CREATE EXTERNAL CATALOG hive COMMENT 'hive catalog' PROPERTIES (
		'type'='hms',
		'hive.metastore.uris' = 'thrift://127.0.0.1:7004',
		'dfs.nameservices'='HANN',
		'dfs.ha.namenodes.HANN'='nn1,nn2',
		'dfs.namenode.rpc-address.HANN.nn1'='nn1_host:rpc_port',
		'dfs.namenode.rpc-address.HANN.nn2'='nn2_host:rpc_port',
		'dfs.client.failover.proxy.provider.HANN'='org.apache.hadoop.hdfs.server.namenode.ha.ConfiguredFailoverProxyProvider'
	)`
	file, errs := parseForTest(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateCatalogStmt)
	if stmt.Name != "hive" {
		t.Errorf("Name = %q, want %q", stmt.Name, "hive")
	}
	if stmt.Comment != "hive catalog" {
		t.Errorf("Comment = %q, want %q", stmt.Comment, "hive catalog")
	}
	if len(stmt.Properties) != 7 {
		t.Fatalf("expected 7 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateCatalog_Legacy_JdbcPostgres(t *testing.T) {
	input := `CREATE EXTERNAL CATALOG jdbc_pg PROPERTIES (
		"type"="jdbc",
		"user"="postgres",
		"password"="123456",
		"jdbc_url" = "jdbc:postgresql://127.0.0.1:5432/demo",
		"driver_url" = "file:///path/to/postgresql-42.5.1.jar",
		"driver_class" = "org.postgresql.Driver"
	)`
	file, errs := parseForTest(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateCatalogStmt)
	if stmt.Name != "jdbc_pg" {
		t.Errorf("Name = %q, want %q", stmt.Name, "jdbc_pg")
	}
	if len(stmt.Properties) != 6 {
		t.Fatalf("expected 6 properties, got %d", len(stmt.Properties))
	}
}

// ---------------------------------------------------------------------------
// Legacy corpus — catalog_alter.sql
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Legacy corpus — catalog_drop.sql
// ---------------------------------------------------------------------------

func TestDropCatalog_Legacy_Basic(t *testing.T) {
	file, errs := parseForTest(`DROP CATALOG hive`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropCatalogStmt)
	if stmt.Name != "hive" || stmt.IfExists {
		t.Errorf("Name=%q IfExists=%v", stmt.Name, stmt.IfExists)
	}
}

func TestDropCatalog_Legacy_IfExists(t *testing.T) {
	file, errs := parseForTest(`DROP CATALOG IF EXISTS hive`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropCatalogStmt)
	if stmt.Name != "hive" || !stmt.IfExists {
		t.Errorf("Name=%q IfExists=%v", stmt.Name, stmt.IfExists)
	}
}

// ---------------------------------------------------------------------------
// Multi-statement
// ---------------------------------------------------------------------------

func TestCatalog_MultiStatement(t *testing.T) {
	input := `CREATE EXTERNAL CATALOG c1 PROPERTIES("type"="hms"); DROP CATALOG c1`
	file, errs := parseForTest(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 2 {
		t.Fatalf("expected 2 stmts, got %d", len(file.Stmts))
	}
	if _, ok := file.Stmts[0].(*ast.CreateCatalogStmt); !ok {
		t.Errorf("Stmts[0]: expected *ast.CreateCatalogStmt, got %T", file.Stmts[0])
	}
	if _, ok := file.Stmts[1].(*ast.DropCatalogStmt); !ok {
		t.Errorf("Stmts[1]: expected *ast.DropCatalogStmt, got %T", file.Stmts[1])
	}
}

// ---------------------------------------------------------------------------
// CREATE EXTERNAL TABLE still dispatches correctly (regression)
// ---------------------------------------------------------------------------

func TestCreateExternalTable_StillWorks(t *testing.T) {
	file, errs := parseForTest(`CREATE EXTERNAL TABLE t (id INT)`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.CreateTableStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateTableStmt, got %T", file.Stmts[0])
	}
	if !stmt.External {
		t.Error("External should be true")
	}
}
