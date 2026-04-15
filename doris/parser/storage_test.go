package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// CREATE STORAGE VAULT
// ---------------------------------------------------------------------------

func TestCreateStorageVault_Basic(t *testing.T) {
	sql := `CREATE STORAGE VAULT hdfs_vault PROPERTIES (
		"type" = "hdfs",
		"fs.defaultFS" = "hdfs://127.0.0.1:8020"
	)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateStorageVaultStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateStorageVaultStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "hdfs_vault" {
		t.Errorf("Name = %q, want %q", stmt.Name, "hdfs_vault")
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
	if stmt.Properties[0].Value != "hdfs" {
		t.Errorf("Properties[0].Value = %q, want %q", stmt.Properties[0].Value, "hdfs")
	}
}

func TestCreateStorageVault_IfNotExists(t *testing.T) {
	sql := `CREATE STORAGE VAULT IF NOT EXISTS s3_vault PROPERTIES (
		"type" = "S3",
		"s3.endpoint" = "s3.us-east-1.amazonaws.com",
		"s3.access_key" = "xxxxxx",
		"s3.secret_key" = "xxxxxx",
		"s3.region" = "us-east-1",
		"s3.root.path" = "prefix",
		"s3.bucket" = "mybucket",
		"provider" = "S3",
		"use_path_style" = "false"
	)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateStorageVaultStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if stmt.Name != "s3_vault" {
		t.Errorf("Name = %q, want %q", stmt.Name, "s3_vault")
	}
	if len(stmt.Properties) != 9 {
		t.Fatalf("expected 9 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateStorageVault_LegacyCorpus_HDFS(t *testing.T) {
	sql := `CREATE STORAGE VAULT IF NOT EXISTS hdfs_vault_demo PROPERTIES (
		"type" = "hdfs",
		"fs.defaultFS" = "hdfs://127.0.0.1:8020",
		"path_prefix" = "big/data",
		"hadoop.username" = "user",
		"hadoop.security.authentication" = "kerberos",
		"hadoop.kerberos.principal" = "hadoop/127.0.0.1@XXX",
		"hadoop.kerberos.keytab" = "/etc/emr.keytab"
	)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateStorageVaultStmt)
	if stmt.Name != "hdfs_vault_demo" {
		t.Errorf("Name = %q, want %q", stmt.Name, "hdfs_vault_demo")
	}
	if len(stmt.Properties) != 7 {
		t.Fatalf("expected 7 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateStorageVault_Tag(t *testing.T) {
	file, errs := Parse(`CREATE STORAGE VAULT v PROPERTIES("type"="S3")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateStorageVaultStmt {
		t.Errorf("Tag() = %v, want T_CreateStorageVaultStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// ALTER STORAGE VAULT
// ---------------------------------------------------------------------------

func TestAlterStorageVault_Basic(t *testing.T) {
	sql := `ALTER STORAGE VAULT old_vault_name PROPERTIES (
		"type" = "S3",
		"VAULT_NAME" = "new_vault_name",
		"s3.access_key" = "new_ak"
	)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterStorageVaultStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterStorageVaultStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "old_vault_name" {
		t.Errorf("Name = %q, want %q", stmt.Name, "old_vault_name")
	}
	if len(stmt.Properties) != 3 {
		t.Fatalf("expected 3 properties, got %d", len(stmt.Properties))
	}
	if stmt.Properties[1].Key != "VAULT_NAME" {
		t.Errorf("Properties[1].Key = %q, want %q", stmt.Properties[1].Key, "VAULT_NAME")
	}
	if stmt.Properties[1].Value != "new_vault_name" {
		t.Errorf("Properties[1].Value = %q, want %q", stmt.Properties[1].Value, "new_vault_name")
	}
}

func TestAlterStorageVault_Tag(t *testing.T) {
	file, errs := Parse(`ALTER STORAGE VAULT v PROPERTIES("type"="S3")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_AlterStorageVaultStmt {
		t.Errorf("Tag() = %v, want T_AlterStorageVaultStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// DROP STORAGE VAULT
// ---------------------------------------------------------------------------

func TestDropStorageVault_Basic(t *testing.T) {
	file, errs := Parse("DROP STORAGE VAULT my_vault")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.DropStorageVaultStmt)
	if !ok {
		t.Fatalf("expected *ast.DropStorageVaultStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "my_vault" {
		t.Errorf("Name = %q, want %q", stmt.Name, "my_vault")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropStorageVault_IfExists(t *testing.T) {
	file, errs := Parse("DROP STORAGE VAULT IF EXISTS my_vault")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropStorageVaultStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if stmt.Name != "my_vault" {
		t.Errorf("Name = %q, want %q", stmt.Name, "my_vault")
	}
}

func TestDropStorageVault_Tag(t *testing.T) {
	file, errs := Parse("DROP STORAGE VAULT v")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_DropStorageVaultStmt {
		t.Errorf("Tag() = %v, want T_DropStorageVaultStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// SET / UNSET DEFAULT STORAGE VAULT
// ---------------------------------------------------------------------------

func TestSetDefaultStorageVault(t *testing.T) {
	file, errs := Parse("SET DEFAULT STORAGE VAULT my_vault")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.SetDefaultStorageVaultStmt)
	if !ok {
		t.Fatalf("expected *ast.SetDefaultStorageVaultStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "my_vault" {
		t.Errorf("Name = %q, want %q", stmt.Name, "my_vault")
	}
}

func TestSetDefaultStorageVault_Tag(t *testing.T) {
	file, errs := Parse("SET DEFAULT STORAGE VAULT v")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_SetDefaultStorageVaultStmt {
		t.Errorf("Tag() = %v, want T_SetDefaultStorageVaultStmt", file.Stmts[0].Tag())
	}
}

func TestUnsetDefaultStorageVault(t *testing.T) {
	file, errs := Parse("UNSET DEFAULT STORAGE VAULT")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	_, ok := file.Stmts[0].(*ast.UnsetDefaultStorageVaultStmt)
	if !ok {
		t.Fatalf("expected *ast.UnsetDefaultStorageVaultStmt, got %T", file.Stmts[0])
	}
}

func TestUnsetDefaultStorageVault_Tag(t *testing.T) {
	file, errs := Parse("UNSET DEFAULT STORAGE VAULT")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_UnsetDefaultStorageVaultStmt {
		t.Errorf("Tag() = %v, want T_UnsetDefaultStorageVaultStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// CREATE STORAGE POLICY
// ---------------------------------------------------------------------------

func TestCreateStoragePolicy_Basic(t *testing.T) {
	sql := `CREATE STORAGE POLICY testPolicy PROPERTIES(
		"storage_resource" = "s3",
		"cooldown_datetime" = "2022-06-08 00:00:00"
	)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateStoragePolicyStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateStoragePolicyStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "testPolicy" {
		t.Errorf("Name = %q, want %q", stmt.Name, "testPolicy")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "storage_resource" {
		t.Errorf("Properties[0].Key = %q, want %q", stmt.Properties[0].Key, "storage_resource")
	}
}

func TestCreateStoragePolicy_CooldownTTL(t *testing.T) {
	sql := `CREATE STORAGE POLICY testPolicy PROPERTIES(
		"storage_resource" = "s3",
		"cooldown_ttl" = "1d"
	)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateStoragePolicyStmt)
	if stmt.Name != "testPolicy" {
		t.Errorf("Name = %q, want %q", stmt.Name, "testPolicy")
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
	if stmt.Properties[1].Key != "cooldown_ttl" {
		t.Errorf("Properties[1].Key = %q, want %q", stmt.Properties[1].Key, "cooldown_ttl")
	}
}

func TestCreateStoragePolicy_IfNotExists(t *testing.T) {
	sql := `CREATE STORAGE POLICY IF NOT EXISTS myPolicy PROPERTIES("storage_resource"="s3")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateStoragePolicyStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
}

func TestCreateStoragePolicy_Tag(t *testing.T) {
	file, errs := Parse(`CREATE STORAGE POLICY p PROPERTIES("storage_resource"="s3")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateStoragePolicyStmt {
		t.Errorf("Tag() = %v, want T_CreateStoragePolicyStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// ALTER STORAGE POLICY
// ---------------------------------------------------------------------------

func TestAlterStoragePolicy_Basic(t *testing.T) {
	sql := `ALTER STORAGE POLICY testPolicy PROPERTIES("cooldown_ttl"="2d")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.AlterStoragePolicyStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterStoragePolicyStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "testPolicy" {
		t.Errorf("Name = %q, want %q", stmt.Name, "testPolicy")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
}

func TestAlterStoragePolicy_Tag(t *testing.T) {
	file, errs := Parse(`ALTER STORAGE POLICY p PROPERTIES("cooldown_ttl"="1d")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_AlterStoragePolicyStmt {
		t.Errorf("Tag() = %v, want T_AlterStoragePolicyStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// DROP STORAGE POLICY
// ---------------------------------------------------------------------------

func TestDropStoragePolicy_Basic(t *testing.T) {
	file, errs := Parse("DROP STORAGE POLICY policy1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.DropStoragePolicyStmt)
	if !ok {
		t.Fatalf("expected *ast.DropStoragePolicyStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "policy1" {
		t.Errorf("Name = %q, want %q", stmt.Name, "policy1")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropStoragePolicy_IfExists(t *testing.T) {
	file, errs := Parse("DROP STORAGE POLICY IF EXISTS policy1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropStoragePolicyStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
}

func TestDropStoragePolicy_Tag(t *testing.T) {
	file, errs := Parse("DROP STORAGE POLICY p")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_DropStoragePolicyStmt {
		t.Errorf("Tag() = %v, want T_DropStoragePolicyStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// CREATE REPOSITORY
// ---------------------------------------------------------------------------

func TestCreateRepository_S3(t *testing.T) {
	sql := `CREATE REPOSITORY s3_repo
		WITH S3
		ON LOCATION "s3://bucket/path"
		PROPERTIES(
			"s3.endpoint" = "s3.us-east-1.amazonaws.com",
			"s3.access_key" = "xxxxxx",
			"s3.secret_key" = "xxxxxx",
			"s3.region" = "us-east-1"
		)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.CreateRepositoryStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateRepositoryStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "s3_repo" {
		t.Errorf("Name = %q, want %q", stmt.Name, "s3_repo")
	}
	if stmt.ReadOnly {
		t.Error("ReadOnly should be false")
	}
	if stmt.Type != "S3" {
		t.Errorf("Type = %q, want %q", stmt.Type, "S3")
	}
	if stmt.BrokerName != "" {
		t.Errorf("BrokerName = %q, want empty", stmt.BrokerName)
	}
	if len(stmt.Properties) != 4 {
		t.Fatalf("expected 4 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateRepository_HDFS(t *testing.T) {
	sql := `CREATE REPOSITORY hdfs_repo
		WITH HDFS
		ON LOCATION "hdfs://namenode:8020/path"
		PROPERTIES(
			"hadoop.username" = "user"
		)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateRepositoryStmt)
	if stmt.Name != "hdfs_repo" {
		t.Errorf("Name = %q, want %q", stmt.Name, "hdfs_repo")
	}
	if stmt.Type != "HDFS" {
		t.Errorf("Type = %q, want %q", stmt.Type, "HDFS")
	}
}

func TestCreateRepository_Broker(t *testing.T) {
	sql := `CREATE REPOSITORY broker_repo
		WITH BROKER my_broker
		ON LOCATION "hdfs://namenode:8020/backup"
		PROPERTIES(
			"username" = "hdfsuser",
			"password" = "secret"
		)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateRepositoryStmt)
	if stmt.Type != "BROKER" {
		t.Errorf("Type = %q, want %q", stmt.Type, "BROKER")
	}
	if stmt.BrokerName != "my_broker" {
		t.Errorf("BrokerName = %q, want %q", stmt.BrokerName, "my_broker")
	}
}

func TestCreateRepository_ReadOnly(t *testing.T) {
	sql := `CREATE READ ONLY REPOSITORY ro_repo
		WITH S3
		ON LOCATION "s3://bucket/path"
		PROPERTIES("s3.endpoint"="s3.amazonaws.com","s3.access_key"="ak","s3.secret_key"="sk","s3.region"="us-east-1")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateRepositoryStmt)
	if !stmt.ReadOnly {
		t.Error("ReadOnly should be true")
	}
	if stmt.Name != "ro_repo" {
		t.Errorf("Name = %q, want %q", stmt.Name, "ro_repo")
	}
}

func TestCreateRepository_Tag(t *testing.T) {
	sql := `CREATE REPOSITORY r WITH S3 ON LOCATION "s3://b/p" PROPERTIES("s3.endpoint"="ep","s3.access_key"="ak","s3.secret_key"="sk","s3.region"="us-east-1")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateRepositoryStmt {
		t.Errorf("Tag() = %v, want T_CreateRepositoryStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// ALTER REPOSITORY
// ---------------------------------------------------------------------------

func TestAlterRepository_Basic(t *testing.T) {
	sql := `ALTER REPOSITORY my_repo PROPERTIES("s3.access_key"="new_key","s3.secret_key"="new_secret")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.AlterRepositoryStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterRepositoryStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "my_repo" {
		t.Errorf("Name = %q, want %q", stmt.Name, "my_repo")
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
}

func TestAlterRepository_Tag(t *testing.T) {
	file, errs := Parse(`ALTER REPOSITORY r PROPERTIES("s3.access_key"="ak")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_AlterRepositoryStmt {
		t.Errorf("Tag() = %v, want T_AlterRepositoryStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// DROP REPOSITORY
// ---------------------------------------------------------------------------

func TestDropRepository_Basic(t *testing.T) {
	file, errs := Parse("DROP REPOSITORY my_repo")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.DropRepositoryStmt)
	if !ok {
		t.Fatalf("expected *ast.DropRepositoryStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "my_repo" {
		t.Errorf("Name = %q, want %q", stmt.Name, "my_repo")
	}
}

func TestDropRepository_Tag(t *testing.T) {
	file, errs := Parse("DROP REPOSITORY r")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_DropRepositoryStmt {
		t.Errorf("Tag() = %v, want T_DropRepositoryStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// CREATE STAGE
// ---------------------------------------------------------------------------

func TestCreateStage_Basic(t *testing.T) {
	sql := `CREATE STAGE my_stage PROPERTIES(
		"endpoint" = "cos.ap-guangzhou.myqcloud.com",
		"access_key_id" = "xxxxxx",
		"access_key_secret" = "xxxxxx",
		"bucket" = "mybucket",
		"path" = "mypath"
	)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.CreateStageStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateStageStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "my_stage" {
		t.Errorf("Name = %q, want %q", stmt.Name, "my_stage")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if len(stmt.Properties) != 5 {
		t.Fatalf("expected 5 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateStage_IfNotExists(t *testing.T) {
	sql := `CREATE STAGE IF NOT EXISTS my_stage PROPERTIES("endpoint"="ep","access_key_id"="ak","access_key_secret"="sk","bucket"="b","path"="p")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateStageStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
}

func TestCreateStage_Tag(t *testing.T) {
	sql := `CREATE STAGE s PROPERTIES("endpoint"="ep","access_key_id"="ak","access_key_secret"="sk","bucket"="b","path"="p")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateStageStmt {
		t.Errorf("Tag() = %v, want T_CreateStageStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// DROP STAGE
// ---------------------------------------------------------------------------

func TestDropStage_Basic(t *testing.T) {
	file, errs := Parse("DROP STAGE my_stage")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.DropStageStmt)
	if !ok {
		t.Fatalf("expected *ast.DropStageStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "my_stage" {
		t.Errorf("Name = %q, want %q", stmt.Name, "my_stage")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropStage_IfExists(t *testing.T) {
	file, errs := Parse("DROP STAGE IF EXISTS my_stage")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropStageStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if stmt.Name != "my_stage" {
		t.Errorf("Name = %q, want %q", stmt.Name, "my_stage")
	}
}

func TestDropStage_Tag(t *testing.T) {
	file, errs := Parse("DROP STAGE s")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_DropStageStmt {
		t.Errorf("Tag() = %v, want T_DropStageStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// CREATE FILE
// ---------------------------------------------------------------------------

func TestCreateFile_Basic(t *testing.T) {
	sql := `CREATE FILE "ca.pem" PROPERTIES("url"="https://example.com/ca.pem","type"="ca")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.CreateFileStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateFileStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "ca.pem" {
		t.Errorf("Name = %q, want %q", stmt.Name, "ca.pem")
	}
	if stmt.Database != "" {
		t.Errorf("Database = %q, want empty", stmt.Database)
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateFile_WithDatabase(t *testing.T) {
	sql := `CREATE FILE "my_cert.pem" IN mydb PROPERTIES("url"="https://example.com/cert.pem","type"="ca")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateFileStmt)
	if stmt.Name != "my_cert.pem" {
		t.Errorf("Name = %q, want %q", stmt.Name, "my_cert.pem")
	}
	if stmt.Database != "mydb" {
		t.Errorf("Database = %q, want %q", stmt.Database, "mydb")
	}
}

func TestCreateFile_Tag(t *testing.T) {
	sql := `CREATE FILE "f.pem" PROPERTIES("url"="https://x.com/f.pem","type"="ca")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateFileStmt {
		t.Errorf("Tag() = %v, want T_CreateFileStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// DROP FILE
// ---------------------------------------------------------------------------

func TestDropFile_Basic(t *testing.T) {
	sql := `DROP FILE "ca.pem" PROPERTIES("type"="ca")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.DropFileStmt)
	if !ok {
		t.Fatalf("expected *ast.DropFileStmt, got %T", file.Stmts[0])
	}
	if stmt.Name != "ca.pem" {
		t.Errorf("Name = %q, want %q", stmt.Name, "ca.pem")
	}
	if stmt.Database != "" {
		t.Errorf("Database = %q, want empty", stmt.Database)
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
}

func TestDropFile_WithDatabase(t *testing.T) {
	sql := `DROP FILE "ca.pem" FROM mydb PROPERTIES("type"="ca")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropFileStmt)
	if stmt.Database != "mydb" {
		t.Errorf("Database = %q, want %q", stmt.Database, "mydb")
	}
}

func TestDropFile_Tag(t *testing.T) {
	sql := `DROP FILE "f.pem" PROPERTIES("type"="ca")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_DropFileStmt {
		t.Errorf("Tag() = %v, want T_DropFileStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// Legacy corpus round-trip
// ---------------------------------------------------------------------------

func TestLegacyCorpus_StoragePolicy(t *testing.T) {
	// From doris/parser/testdata/legacy/cluster_storage_policy.sql
	tests := []struct {
		sql      string
		wantType string
	}{
		{
			`CREATE STORAGE POLICY testPolicy PROPERTIES(
				"storage_resource" = "s3",
				"cooldown_datetime" = "2022-06-08 00:00:00"
			)`,
			"*ast.CreateStoragePolicyStmt",
		},
		{
			`CREATE STORAGE POLICY testPolicy PROPERTIES(
				"storage_resource" = "s3",
				"cooldown_ttl" = "1d"
			)`,
			"*ast.CreateStoragePolicyStmt",
		},
		{
			`DROP STORAGE POLICY policy1`,
			"*ast.DropStoragePolicyStmt",
		},
	}

	for _, tt := range tests {
		file, errs := Parse(tt.sql)
		if len(errs) != 0 {
			t.Errorf("sql %q: unexpected errors: %v", tt.sql, errs)
			continue
		}
		if len(file.Stmts) != 1 {
			t.Errorf("sql %q: expected 1 stmt, got %d", tt.sql, len(file.Stmts))
		}
	}
}

func TestLegacyCorpus_StorageVault(t *testing.T) {
	// From doris/parser/testdata/legacy/cluster_storage_vault.sql
	tests := []string{
		`CREATE STORAGE VAULT IF NOT EXISTS hdfs_vault_demo PROPERTIES (
			"type" = "hdfs",
			"fs.defaultFS" = "hdfs://127.0.0.1:8020",
			"path_prefix" = "big/data",
			"hadoop.username" = "user",
			"hadoop.security.authentication" = "kerberos",
			"hadoop.kerberos.principal" = "hadoop/127.0.0.1@XXX",
			"hadoop.kerberos.keytab" = "/etc/emr.keytab"
		)`,
		`CREATE STORAGE VAULT IF NOT EXISTS oss_demo_vault PROPERTIES (
			"type" = "S3",
			"s3.endpoint" = "oss-cn-beijing.aliyuncs.com",
			"s3.access_key" = "xxxxxx",
			"s3.secret_key" = "xxxxxx",
			"s3.region" = "cn-beijing",
			"s3.root.path" = "oss_demo_vault_prefix",
			"s3.bucket" = "xxxxxx",
			"provider" = "OSS",
			"use_path_style" = "false"
		)`,
		`ALTER STORAGE VAULT old_vault_name PROPERTIES (
			"type" = "S3",
			"VAULT_NAME" = "new_vault_name",
			"s3.access_key" = "new_ak"
		)`,
	}

	for _, sql := range tests {
		file, errs := Parse(sql)
		if len(errs) != 0 {
			t.Errorf("sql:\n%s\nunexpected errors: %v", sql, errs)
			continue
		}
		if len(file.Stmts) != 1 {
			t.Errorf("sql:\n%s\nexpected 1 stmt, got %d", sql, len(file.Stmts))
		}
	}
}

// ---------------------------------------------------------------------------
// Multi-statement round-trip
// ---------------------------------------------------------------------------

func TestStorage_MultiStatement(t *testing.T) {
	input := `CREATE STORAGE VAULT v PROPERTIES("type"="S3");
DROP STORAGE VAULT v;
CREATE STORAGE POLICY p PROPERTIES("storage_resource"="s3");
DROP STORAGE POLICY p`
	file, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 4 {
		t.Fatalf("expected 4 stmts, got %d", len(file.Stmts))
	}
	if _, ok := file.Stmts[0].(*ast.CreateStorageVaultStmt); !ok {
		t.Errorf("Stmts[0]: expected CreateStorageVaultStmt, got %T", file.Stmts[0])
	}
	if _, ok := file.Stmts[1].(*ast.DropStorageVaultStmt); !ok {
		t.Errorf("Stmts[1]: expected DropStorageVaultStmt, got %T", file.Stmts[1])
	}
	if _, ok := file.Stmts[2].(*ast.CreateStoragePolicyStmt); !ok {
		t.Errorf("Stmts[2]: expected CreateStoragePolicyStmt, got %T", file.Stmts[2])
	}
	if _, ok := file.Stmts[3].(*ast.DropStoragePolicyStmt); !ok {
		t.Errorf("Stmts[3]: expected DropStoragePolicyStmt, got %T", file.Stmts[3])
	}
}
