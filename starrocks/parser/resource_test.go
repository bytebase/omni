package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

func TestCreateResource_Basic(t *testing.T) {
	file, errs := parseForTest(`CREATE RESOURCE mysql_resource PROPERTIES (
		"type"="jdbc",
		"user"="root",
		"password"="123456"
	)`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateResourceStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateResourceStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "mysql_resource" {
		t.Errorf("Name = %q, want %q", stmt.Name, "mysql_resource")
	}
	if stmt.External {
		t.Error("External should be false")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if len(stmt.Properties) != 3 {
		t.Fatalf("expected 3 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateResource_QuotedName(t *testing.T) {
	file, errs := parseForTest(`CREATE RESOURCE "remote_s3" PROPERTIES(
		"type"="s3",
		"s3.endpoint"="bj.s3.com"
	)`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateResourceStmt)
	if stmt.Name != "remote_s3" {
		t.Errorf("Name = %q, want %q", stmt.Name, "remote_s3")
	}
	if stmt.External {
		t.Error("External should be false")
	}
}

func TestCreateResource_IfNotExists(t *testing.T) {
	file, errs := parseForTest(`CREATE RESOURCE IF NOT EXISTS r1 PROPERTIES("type"="s3")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateResourceStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
}

func TestCreateResource_Tag(t *testing.T) {
	file, errs := parseForTest(`CREATE RESOURCE r1 PROPERTIES("type"="s3")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateResourceStmt {
		t.Errorf("Tag() = %v, want T_CreateResourceStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// RESOURCE — CREATE EXTERNAL
// ---------------------------------------------------------------------------

func TestCreateExternalResource_Basic(t *testing.T) {
	file, errs := parseForTest(`CREATE EXTERNAL RESOURCE "spark0"
	PROPERTIES(
		"type" = "spark",
		"spark.master" = "yarn",
		"spark.submit.deployMode" = "cluster",
		"working_dir" = "hdfs://127.0.0.1:10000/tmp/doris",
		"broker" = "broker0"
	)`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateResourceStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateResourceStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "spark0" {
		t.Errorf("Name = %q, want %q", stmt.Name, "spark0")
	}
	if !stmt.External {
		t.Error("External should be true")
	}
	if len(stmt.Properties) != 5 {
		t.Fatalf("expected 5 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateExternalResource_BacktickName(t *testing.T) {
	file, errs := parseForTest("CREATE EXTERNAL RESOURCE `oracle_odbc` PROPERTIES (\"type\" = \"odbc_catalog\", \"host\" = \"192.168.0.1\")")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateResourceStmt)
	if stmt.Name != "oracle_odbc" {
		t.Errorf("Name = %q, want %q", stmt.Name, "oracle_odbc")
	}
	if !stmt.External {
		t.Error("External should be true")
	}
}

func TestCreateExternalResource_Tag(t *testing.T) {
	file, errs := parseForTest(`CREATE EXTERNAL RESOURCE "spark0" PROPERTIES("type"="spark")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateResourceStmt {
		t.Errorf("Tag() = %v, want T_CreateResourceStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// RESOURCE — ALTER
// ---------------------------------------------------------------------------

func TestAlterResource_Basic(t *testing.T) {
	file, errs := parseForTest(`ALTER RESOURCE spark0 SET PROPERTIES("spark.executor.memory"="2g")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterResourceStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterResourceStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "spark0" {
		t.Errorf("Name = %q, want %q", stmt.Name, "spark0")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
}

func TestAlterResource_Tag(t *testing.T) {
	file, errs := parseForTest(`ALTER RESOURCE r1 SET PROPERTIES("type"="s3")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_AlterResourceStmt {
		t.Errorf("Tag() = %v, want T_AlterResourceStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// RESOURCE — DROP
// ---------------------------------------------------------------------------

func TestDropResource_Basic(t *testing.T) {
	file, errs := parseForTest(`DROP RESOURCE 'spark0'`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropResourceStmt)
	if !ok {
		t.Fatalf("expected *ast.DropResourceStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "spark0" {
		t.Errorf("Name = %q, want %q", stmt.Name, "spark0")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropResource_IfExists(t *testing.T) {
	file, errs := parseForTest(`DROP RESOURCE IF EXISTS spark0`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropResourceStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
}

func TestDropResource_Tag(t *testing.T) {
	file, errs := parseForTest(`DROP RESOURCE 'spark0'`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_DropResourceStmt {
		t.Errorf("Tag() = %v, want T_DropResourceStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// RESOURCE — legacy corpus
// ---------------------------------------------------------------------------

func TestResource_LegacyCorpus(t *testing.T) {
	cases := []string{
		`CREATE EXTERNAL RESOURCE "spark0"
		PROPERTIES(
			"type" = "spark",
			"spark.master" = "yarn",
			"spark.submit.deployMode" = "cluster",
			"spark.jars" = "xxx.jar,yyy.jar",
			"spark.files" = "/tmp/aaa,/tmp/bbb",
			"spark.executor.memory" = "1g",
			"spark.yarn.queue" = "queue0",
			"spark.hadoop.yarn.resourcemanager.address" = "127.0.0.1:9999",
			"spark.hadoop.fs.defaultFS" = "hdfs://127.0.0.1:10000",
			"working_dir" = "hdfs://127.0.0.1:10000/tmp/doris",
			"broker" = "broker0",
			"broker.username" = "user0",
			"broker.password" = "password0"
		)`,
		`CREATE RESOURCE "remote_s3"
		PROPERTIES(
			"type" = "s3",
			"s3.endpoint" = "bj.s3.com",
			"s3.region" = "bj",
			"s3.access_key" = "bbb",
			"s3.secret_key" = "aaaa",
			"s3.connection.maximum" = "50",
			"s3.connection.request.timeout" = "3000",
			"s3.connection.timeout" = "1000"
		)`,
		`CREATE RESOURCE mysql_resource PROPERTIES (
			"type"="jdbc",
			"user"="root",
			"password"="123456",
			"jdbc_url" = "jdbc:mysql://127.0.0.1:3316/doris_test?useSSL=false",
			"driver_url" = "https://doris-community-test-1308700295.cos.ap-hongkong.myqcloud.com/jdbc_driver/mysql-connector-java-8.0.25.jar",
			"driver_class" = "com.mysql.cj.jdbc.Driver"
		)`,
		`DROP RESOURCE 'spark0'`,
	}

	for i, sql := range cases {
		file, errs := parseForTest(sql)
		if len(errs) != 0 {
			t.Errorf("case %d: unexpected errors: %v", i, errs)
			continue
		}
		if len(file.Stmts) != 1 {
			t.Errorf("case %d: expected 1 stmt, got %d", i, len(file.Stmts))
		}
	}
}

// ---------------------------------------------------------------------------
// SQL BLOCK RULE — CREATE
// ---------------------------------------------------------------------------

func TestResource_MultiStatement(t *testing.T) {
	input := `CREATE EXTERNAL RESOURCE "spark0" PROPERTIES("type"="spark");
CREATE RESOURCE "remote_s3" PROPERTIES("type"="s3");
DROP RESOURCE 'spark0'`

	file, errs := parseForTest(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 3 {
		t.Fatalf("expected 3 stmts, got %d", len(file.Stmts))
	}
	stmt0 := file.Stmts[0].(*ast.CreateResourceStmt)
	if !stmt0.External {
		t.Error("Stmts[0]: External should be true")
	}
	stmt1 := file.Stmts[1].(*ast.CreateResourceStmt)
	if stmt1.External {
		t.Error("Stmts[1]: External should be false")
	}
}

// ---------------------------------------------------------------------------
// Existing CREATE EXTERNAL TABLE still works
// ---------------------------------------------------------------------------

func TestCreateExternalTable_NotResource(t *testing.T) {
	// CREATE EXTERNAL TABLE should still route to parseCreateTable, not parseCreateResource.
	file, errs := parseForTest("CREATE EXTERNAL TABLE t (id INT) ENGINE=HIVE PROPERTIES(\"database\"=\"db\", \"table\"=\"t\")")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateTableStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateTableStmt, got %T", file.Stmts[0])
	}
	if !stmt.External {
		t.Error("External should be true on CreateTableStmt")
	}
}
