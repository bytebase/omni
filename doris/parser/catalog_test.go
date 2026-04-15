package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// CREATE CATALOG
// ---------------------------------------------------------------------------

func TestCreateCatalog_Basic(t *testing.T) {
	file, errs := Parse(`CREATE CATALOG hive_catalog PROPERTIES("type"="hms", "hive.metastore.uris"="thrift://127.0.0.1:7004")`)
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
	if stmt.External {
		t.Error("External should be false")
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
	file, errs := Parse(`CREATE EXTERNAL CATALOG iceberg_catalog PROPERTIES("type"="iceberg")`)
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
	file, errs := Parse(`CREATE CATALOG IF NOT EXISTS my_catalog PROPERTIES("type"="hms")`)
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
	file, errs := Parse(`CREATE CATALOG hive COMMENT 'hive catalog' PROPERTIES('type'='hms', 'hive.metastore.uris'='thrift://127.0.0.1:7004')`)
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

func TestCreateCatalog_NoProperties(t *testing.T) {
	file, errs := Parse(`CREATE CATALOG my_catalog`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateCatalogStmt)
	if stmt.Name != "my_catalog" {
		t.Errorf("Name = %q, want %q", stmt.Name, "my_catalog")
	}
	if len(stmt.Properties) != 0 {
		t.Errorf("expected 0 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateCatalog_WithResource(t *testing.T) {
	file, errs := Parse(`CREATE CATALOG hms_catalog WITH RESOURCE hms_resource`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateCatalogStmt)
	if stmt.Name != "hms_catalog" {
		t.Errorf("Name = %q, want %q", stmt.Name, "hms_catalog")
	}
	if stmt.WithResource != "hms_resource" {
		t.Errorf("WithResource = %q, want %q", stmt.WithResource, "hms_resource")
	}
}

func TestCreateCatalog_ES(t *testing.T) {
	file, errs := Parse(`CREATE CATALOG es PROPERTIES ("type"="es", "hosts"="http://127.0.0.1:9200")`)
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
	input := `CREATE CATALOG jdbc PROPERTIES (
		"type"="jdbc",
		"user"="root",
		"password"="123456",
		"jdbc_url" = "jdbc:mysql://127.0.0.1:3316/doris_test?useSSL=false",
		"driver_url" = "https://example.com/mysql-connector-java-8.0.25.jar",
		"driver_class" = "com.mysql.cj.jdbc.Driver"
	)`
	file, errs := Parse(input)
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
	file, errs := Parse(`CREATE CATALOG c PROPERTIES("type"="hms")`)
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

func TestAlterCatalog_Rename(t *testing.T) {
	file, errs := Parse(`ALTER CATALOG ctlg_hive RENAME hive`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterCatalogStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterCatalogStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "ctlg_hive" {
		t.Errorf("Name = %q, want %q", stmt.Name, "ctlg_hive")
	}
	if stmt.Action != ast.AlterCatalogRename {
		t.Errorf("Action = %v, want AlterCatalogRename", stmt.Action)
	}
	if stmt.NewName != "hive" {
		t.Errorf("NewName = %q, want %q", stmt.NewName, "hive")
	}
}

func TestAlterCatalog_SetProperties(t *testing.T) {
	file, errs := Parse(`ALTER CATALOG hive SET PROPERTIES ('hive.metastore.uris'='thrift://172.21.0.1:9083')`)
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
	if stmt.Action != ast.AlterCatalogSetProperties {
		t.Errorf("Action = %v, want AlterCatalogSetProperties", stmt.Action)
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "hive.metastore.uris" {
		t.Errorf("Property key = %q, want %q", stmt.Properties[0].Key, "hive.metastore.uris")
	}
}

func TestAlterCatalog_ModifyComment(t *testing.T) {
	file, errs := Parse(`ALTER CATALOG hive MODIFY COMMENT "new catalog comment"`)
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
	if stmt.Action != ast.AlterCatalogModifyComment {
		t.Errorf("Action = %v, want AlterCatalogModifyComment", stmt.Action)
	}
	if stmt.Comment != "new catalog comment" {
		t.Errorf("Comment = %q, want %q", stmt.Comment, "new catalog comment")
	}
}

func TestAlterCatalog_Property(t *testing.T) {
	file, errs := Parse(`ALTER CATALOG c PROPERTY ("key"="value")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.AlterCatalogStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterCatalogStmt, got %T", file.Stmts[0])
	}
	if stmt.Action != ast.AlterCatalogSetProperty {
		t.Errorf("Action = %v, want AlterCatalogSetProperty", stmt.Action)
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "key" || stmt.Properties[0].Value != "value" {
		t.Errorf("unexpected property: key=%q val=%q", stmt.Properties[0].Key, stmt.Properties[0].Value)
	}
}

func TestAlterCatalog_Tag(t *testing.T) {
	file, errs := Parse(`ALTER CATALOG c RENAME c2`)
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
	file, errs := Parse(`DROP CATALOG hive`)
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
	file, errs := Parse(`DROP CATALOG IF EXISTS hive`)
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
	file, errs := Parse(`DROP CATALOG c`)
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
	file, errs := Parse(`REFRESH CATALOG c`)
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
	file, errs := Parse(`REFRESH CATALOG c PROPERTIES("key"="value")`)
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
	file, errs := Parse(`REFRESH CATALOG c`)
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
	input := `CREATE CATALOG hive COMMENT 'hive catalog' PROPERTIES (
		'type'='hms',
		'hive.metastore.uris' = 'thrift://127.0.0.1:7004',
		'dfs.nameservices'='HANN',
		'dfs.ha.namenodes.HANN'='nn1,nn2',
		'dfs.namenode.rpc-address.HANN.nn1'='nn1_host:rpc_port',
		'dfs.namenode.rpc-address.HANN.nn2'='nn2_host:rpc_port',
		'dfs.client.failover.proxy.provider.HANN'='org.apache.hadoop.hdfs.server.namenode.ha.ConfiguredFailoverProxyProvider'
	)`
	file, errs := Parse(input)
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
	input := `CREATE CATALOG jdbc_pg PROPERTIES (
		"type"="jdbc",
		"user"="postgres",
		"password"="123456",
		"jdbc_url" = "jdbc:postgresql://127.0.0.1:5432/demo",
		"driver_url" = "file:///path/to/postgresql-42.5.1.jar",
		"driver_class" = "org.postgresql.Driver"
	)`
	file, errs := Parse(input)
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

func TestAlterCatalog_Legacy_Rename(t *testing.T) {
	file, errs := Parse(`ALTER CATALOG ctlg_hive RENAME hive`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AlterCatalogStmt)
	if stmt.Name != "ctlg_hive" || stmt.NewName != "hive" {
		t.Errorf("Name=%q NewName=%q", stmt.Name, stmt.NewName)
	}
}

func TestAlterCatalog_Legacy_SetProperties(t *testing.T) {
	file, errs := Parse(`ALTER CATALOG hive SET PROPERTIES ('hive.metastore.uris'='thrift://172.21.0.1:9083')`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AlterCatalogStmt)
	if stmt.Action != ast.AlterCatalogSetProperties {
		t.Errorf("Action = %v, want AlterCatalogSetProperties", stmt.Action)
	}
}

func TestAlterCatalog_Legacy_ModifyComment(t *testing.T) {
	file, errs := Parse(`ALTER CATALOG hive MODIFY COMMENT "new catalog comment"`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AlterCatalogStmt)
	if stmt.Comment != "new catalog comment" {
		t.Errorf("Comment = %q, want %q", stmt.Comment, "new catalog comment")
	}
}

// ---------------------------------------------------------------------------
// Legacy corpus — catalog_drop.sql
// ---------------------------------------------------------------------------

func TestDropCatalog_Legacy_Basic(t *testing.T) {
	file, errs := Parse(`DROP CATALOG hive`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropCatalogStmt)
	if stmt.Name != "hive" || stmt.IfExists {
		t.Errorf("Name=%q IfExists=%v", stmt.Name, stmt.IfExists)
	}
}

func TestDropCatalog_Legacy_IfExists(t *testing.T) {
	file, errs := Parse(`DROP CATALOG IF EXISTS hive`)
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
	input := `CREATE CATALOG c1 PROPERTIES("type"="hms"); DROP CATALOG c1`
	file, errs := Parse(input)
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
	file, errs := Parse(`CREATE EXTERNAL TABLE t (id INT)`)
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
