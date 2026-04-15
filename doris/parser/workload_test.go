package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// WORKLOAD GROUP — CREATE
// ---------------------------------------------------------------------------

func TestCreateWorkloadGroup_Basic(t *testing.T) {
	file, errs := Parse(`CREATE WORKLOAD GROUP g1 PROPERTIES ("max_cpu_percent"="10%")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateWorkloadGroupStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateWorkloadGroupStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "g1" {
		t.Errorf("Name = %q, want %q", stmt.Name, "g1")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "max_cpu_percent" {
		t.Errorf("Properties[0].Key = %q, want %q", stmt.Properties[0].Key, "max_cpu_percent")
	}
	if stmt.Properties[0].Value != "10%" {
		t.Errorf("Properties[0].Value = %q, want %q", stmt.Properties[0].Value, "10%")
	}
}

func TestCreateWorkloadGroup_IfNotExists(t *testing.T) {
	file, errs := Parse(`CREATE WORKLOAD GROUP IF NOT EXISTS g1 PROPERTIES (
		"max_cpu_percent"="10%",
		"max_memory_percent"="30%"
	)`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateWorkloadGroupStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if stmt.Name != "g1" {
		t.Errorf("Name = %q, want %q", stmt.Name, "g1")
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateWorkloadGroup_Tag(t *testing.T) {
	file, errs := Parse(`CREATE WORKLOAD GROUP g1 PROPERTIES ("max_cpu_percent"="10%")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateWorkloadGroupStmt {
		t.Errorf("Tag() = %v, want T_CreateWorkloadGroupStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// WORKLOAD GROUP — ALTER
// ---------------------------------------------------------------------------

func TestAlterWorkloadGroup_Basic(t *testing.T) {
	file, errs := Parse(`ALTER WORKLOAD GROUP g1 PROPERTIES (
		"max_cpu_percent"="20%",
		"max_memory_percent"="40%"
	)`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterWorkloadGroupStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterWorkloadGroupStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "g1" {
		t.Errorf("Name = %q, want %q", stmt.Name, "g1")
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
}

func TestAlterWorkloadGroup_Tag(t *testing.T) {
	file, errs := Parse(`ALTER WORKLOAD GROUP g1 PROPERTIES ("max_cpu_percent"="20%")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_AlterWorkloadGroupStmt {
		t.Errorf("Tag() = %v, want T_AlterWorkloadGroupStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// WORKLOAD GROUP — DROP
// ---------------------------------------------------------------------------

func TestDropWorkloadGroup_Basic(t *testing.T) {
	file, errs := Parse(`DROP WORKLOAD GROUP g1`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropWorkloadGroupStmt)
	if !ok {
		t.Fatalf("expected *ast.DropWorkloadGroupStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "g1" {
		t.Errorf("Name = %q, want %q", stmt.Name, "g1")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropWorkloadGroup_IfExists(t *testing.T) {
	file, errs := Parse(`DROP WORKLOAD GROUP IF EXISTS g1`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropWorkloadGroupStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if stmt.Name != "g1" {
		t.Errorf("Name = %q, want %q", stmt.Name, "g1")
	}
}

func TestDropWorkloadGroup_Tag(t *testing.T) {
	file, errs := Parse(`DROP WORKLOAD GROUP IF EXISTS g1`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_DropWorkloadGroupStmt {
		t.Errorf("Tag() = %v, want T_DropWorkloadGroupStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// WORKLOAD GROUP — legacy corpus
// ---------------------------------------------------------------------------

func TestWorkloadGroup_LegacyCorpus(t *testing.T) {
	// From cluster_workload_group.sql
	cases := []struct {
		sql  string
		want interface{}
	}{
		{
			sql:  `CREATE WORKLOAD GROUP IF NOT EXISTS g1 PROPERTIES ("max_cpu_percent"="10%", "max_memory_percent"="30%")`,
			want: (*ast.CreateWorkloadGroupStmt)(nil),
		},
		{
			sql:  `ALTER WORKLOAD GROUP g1 PROPERTIES ("max_cpu_percent"="20%", "max_memory_percent"="40%")`,
			want: (*ast.AlterWorkloadGroupStmt)(nil),
		},
		{
			sql:  `DROP WORKLOAD GROUP IF EXISTS g1`,
			want: (*ast.DropWorkloadGroupStmt)(nil),
		},
	}

	for _, tc := range cases {
		t.Run(tc.sql[:20], func(t *testing.T) {
			file, errs := Parse(tc.sql)
			if len(errs) != 0 {
				t.Fatalf("unexpected errors: %v", errs)
			}
			if len(file.Stmts) != 1 {
				t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// WORKLOAD POLICY — CREATE
// ---------------------------------------------------------------------------

func TestCreateWorkloadPolicy_Basic(t *testing.T) {
	file, errs := Parse(`CREATE WORKLOAD POLICY p1
		CONDITIONS(query_time > 1000)
		ACTIONS(cancel_query)
		PROPERTIES("enabled"="true")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateWorkloadPolicyStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateWorkloadPolicyStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "p1" {
		t.Errorf("Name = %q, want %q", stmt.Name, "p1")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if len(stmt.Conditions) != 1 {
		t.Errorf("expected 1 condition, got %d", len(stmt.Conditions))
	}
	if len(stmt.Actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(stmt.Actions))
	}
	if len(stmt.Properties) != 1 {
		t.Errorf("expected 1 property, got %d", len(stmt.Properties))
	}
}

func TestCreateWorkloadPolicy_IfNotExists(t *testing.T) {
	file, errs := Parse(`CREATE WORKLOAD POLICY IF NOT EXISTS p1
		CONDITIONS(scan_rows > 1000000000, query_time > 1000)
		ACTIONS(cancel_query)`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateWorkloadPolicyStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if len(stmt.Conditions) != 2 {
		t.Errorf("expected 2 conditions, got %d", len(stmt.Conditions))
	}
	if len(stmt.Actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(stmt.Actions))
	}
}

func TestCreateWorkloadPolicy_Tag(t *testing.T) {
	file, errs := Parse(`CREATE WORKLOAD POLICY p1 CONDITIONS(query_time > 1000) ACTIONS(cancel_query)`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateWorkloadPolicyStmt {
		t.Errorf("Tag() = %v, want T_CreateWorkloadPolicyStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// WORKLOAD POLICY — ALTER
// ---------------------------------------------------------------------------

func TestAlterWorkloadPolicy_Basic(t *testing.T) {
	file, errs := Parse(`ALTER WORKLOAD POLICY p1 PROPERTIES("enabled"="false")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterWorkloadPolicyStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterWorkloadPolicyStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "p1" {
		t.Errorf("Name = %q, want %q", stmt.Name, "p1")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "enabled" {
		t.Errorf("Properties[0].Key = %q, want %q", stmt.Properties[0].Key, "enabled")
	}
}

func TestAlterWorkloadPolicy_Tag(t *testing.T) {
	file, errs := Parse(`ALTER WORKLOAD POLICY p1 PROPERTIES("enabled"="false")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_AlterWorkloadPolicyStmt {
		t.Errorf("Tag() = %v, want T_AlterWorkloadPolicyStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// WORKLOAD POLICY — DROP
// ---------------------------------------------------------------------------

func TestDropWorkloadPolicy_Basic(t *testing.T) {
	file, errs := Parse(`DROP WORKLOAD POLICY p1`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.DropWorkloadPolicyStmt)
	if !ok {
		t.Fatalf("expected *ast.DropWorkloadPolicyStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "p1" {
		t.Errorf("Name = %q, want %q", stmt.Name, "p1")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropWorkloadPolicy_IfExists(t *testing.T) {
	file, errs := Parse(`DROP WORKLOAD POLICY IF EXISTS p1`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropWorkloadPolicyStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
}

func TestDropWorkloadPolicy_Tag(t *testing.T) {
	file, errs := Parse(`DROP WORKLOAD POLICY IF EXISTS p1`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_DropWorkloadPolicyStmt {
		t.Errorf("Tag() = %v, want T_DropWorkloadPolicyStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// RESOURCE — CREATE (non-EXTERNAL)
// ---------------------------------------------------------------------------

func TestCreateResource_Basic(t *testing.T) {
	file, errs := Parse(`CREATE RESOURCE mysql_resource PROPERTIES (
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
	file, errs := Parse(`CREATE RESOURCE "remote_s3" PROPERTIES(
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
	file, errs := Parse(`CREATE RESOURCE IF NOT EXISTS r1 PROPERTIES("type"="s3")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateResourceStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
}

func TestCreateResource_Tag(t *testing.T) {
	file, errs := Parse(`CREATE RESOURCE r1 PROPERTIES("type"="s3")`)
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
	file, errs := Parse(`CREATE EXTERNAL RESOURCE "spark0"
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
	file, errs := Parse("CREATE EXTERNAL RESOURCE `oracle_odbc` PROPERTIES (\"type\" = \"odbc_catalog\", \"host\" = \"192.168.0.1\")")
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
	file, errs := Parse(`CREATE EXTERNAL RESOURCE "spark0" PROPERTIES("type"="spark")`)
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
	file, errs := Parse(`ALTER RESOURCE spark0 PROPERTIES("spark.executor.memory"="2g")`)
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
	file, errs := Parse(`ALTER RESOURCE r1 PROPERTIES("type"="s3")`)
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
	file, errs := Parse(`DROP RESOURCE 'spark0'`)
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
	file, errs := Parse(`DROP RESOURCE IF EXISTS spark0`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropResourceStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
}

func TestDropResource_Tag(t *testing.T) {
	file, errs := Parse(`DROP RESOURCE 'spark0'`)
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
		file, errs := Parse(sql)
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

func TestCreateSQLBlockRule_Basic(t *testing.T) {
	file, errs := Parse(`CREATE SQL_BLOCK_RULE test_rule PROPERTIES (
		"sql"="select \\* from order_analysis",
		"enable"="true"
	)`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateSQLBlockRuleStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateSQLBlockRuleStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "test_rule" {
		t.Errorf("Name = %q, want %q", stmt.Name, "test_rule")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateSQLBlockRule_IfNotExists(t *testing.T) {
	file, errs := Parse(`CREATE SQL_BLOCK_RULE IF NOT EXISTS test_rule PROPERTIES("sql"="select 1", "enable"="true")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateSQLBlockRuleStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
}

func TestCreateSQLBlockRule_ScanRowsLimit(t *testing.T) {
	file, errs := Parse(`CREATE SQL_BLOCK_RULE test_rule2 PROPERTIES (
		"scan_row_limit"="100",
		"enable"="true"
	)`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateSQLBlockRuleStmt)
	if stmt.Name != "test_rule2" {
		t.Errorf("Name = %q, want %q", stmt.Name, "test_rule2")
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateSQLBlockRule_Tag(t *testing.T) {
	file, errs := Parse(`CREATE SQL_BLOCK_RULE r1 PROPERTIES("sql"="select 1", "enable"="true")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateSQLBlockRuleStmt {
		t.Errorf("Tag() = %v, want T_CreateSQLBlockRuleStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// SQL BLOCK RULE — ALTER
// ---------------------------------------------------------------------------

func TestAlterSQLBlockRule_Basic(t *testing.T) {
	file, errs := Parse(`ALTER SQL_BLOCK_RULE test_rule PROPERTIES("enable"="false")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterSQLBlockRuleStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterSQLBlockRuleStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "test_rule" {
		t.Errorf("Name = %q, want %q", stmt.Name, "test_rule")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
}

func TestAlterSQLBlockRule_Tag(t *testing.T) {
	file, errs := Parse(`ALTER SQL_BLOCK_RULE r1 PROPERTIES("enable"="false")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_AlterSQLBlockRuleStmt {
		t.Errorf("Tag() = %v, want T_AlterSQLBlockRuleStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// SQL BLOCK RULE — DROP
// ---------------------------------------------------------------------------

func TestDropSQLBlockRule_Basic(t *testing.T) {
	file, errs := Parse(`DROP SQL_BLOCK_RULE test_rule`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropSQLBlockRuleStmt)
	if !ok {
		t.Fatalf("expected *ast.DropSQLBlockRuleStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "test_rule" {
		t.Errorf("Name = %q, want %q", stmt.Name, "test_rule")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropSQLBlockRule_IfExists(t *testing.T) {
	file, errs := Parse(`DROP SQL_BLOCK_RULE IF EXISTS test_rule`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropSQLBlockRuleStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
}

func TestDropSQLBlockRule_Tag(t *testing.T) {
	file, errs := Parse(`DROP SQL_BLOCK_RULE IF EXISTS test_rule`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_DropSQLBlockRuleStmt {
		t.Errorf("Tag() = %v, want T_DropSQLBlockRuleStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// Multi-statement round-trip
// ---------------------------------------------------------------------------

func TestWorkload_MultiStatement(t *testing.T) {
	input := `CREATE WORKLOAD GROUP IF NOT EXISTS g1 PROPERTIES ("max_cpu_percent"="10%", "max_memory_percent"="30%");
ALTER WORKLOAD GROUP g1 PROPERTIES ("max_cpu_percent"="20%", "max_memory_percent"="40%");
DROP WORKLOAD GROUP IF EXISTS g1`

	file, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 3 {
		t.Fatalf("expected 3 stmts, got %d", len(file.Stmts))
	}
	if _, ok := file.Stmts[0].(*ast.CreateWorkloadGroupStmt); !ok {
		t.Errorf("Stmts[0]: expected *CreateWorkloadGroupStmt, got %T", file.Stmts[0])
	}
	if _, ok := file.Stmts[1].(*ast.AlterWorkloadGroupStmt); !ok {
		t.Errorf("Stmts[1]: expected *AlterWorkloadGroupStmt, got %T", file.Stmts[1])
	}
	if _, ok := file.Stmts[2].(*ast.DropWorkloadGroupStmt); !ok {
		t.Errorf("Stmts[2]: expected *DropWorkloadGroupStmt, got %T", file.Stmts[2])
	}
}

func TestResource_MultiStatement(t *testing.T) {
	input := `CREATE EXTERNAL RESOURCE "spark0" PROPERTIES("type"="spark");
CREATE RESOURCE "remote_s3" PROPERTIES("type"="s3");
DROP RESOURCE 'spark0'`

	file, errs := Parse(input)
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

func TestCreateExternalTable_StillWorks(t *testing.T) {
	// CREATE EXTERNAL TABLE should still route to parseCreateTable, not parseCreateResource.
	file, errs := Parse("CREATE EXTERNAL TABLE t (id INT) ENGINE=HIVE PROPERTIES(\"database\"=\"db\", \"table\"=\"t\")")
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
