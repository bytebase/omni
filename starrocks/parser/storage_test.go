package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// CREATE STORAGE VAULT
// ---------------------------------------------------------------------------

func TestCreateStoragePolicy_Basic(t *testing.T) {
	sql := `CREATE STORAGE POLICY testPolicy PROPERTIES(
		"storage_resource" = "s3",
		"cooldown_datetime" = "2022-06-08 00:00:00"
	)`
	file, errs := parseForTest(sql)
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
	file, errs := parseForTest(sql)
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
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateStoragePolicyStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
}

func TestCreateStoragePolicy_Tag(t *testing.T) {
	file, errs := parseForTest(`CREATE STORAGE POLICY p PROPERTIES("storage_resource"="s3")`)
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
	file, errs := parseForTest(sql)
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
	file, errs := parseForTest(`ALTER STORAGE POLICY p PROPERTIES("cooldown_ttl"="1d")`)
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
	file, errs := parseForTest("DROP STORAGE POLICY policy1")
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
	file, errs := parseForTest("DROP STORAGE POLICY IF EXISTS policy1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropStoragePolicyStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
}

func TestDropStoragePolicy_Tag(t *testing.T) {
	file, errs := parseForTest("DROP STORAGE POLICY p")
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

func TestCreateRepository_Broker(t *testing.T) {
	sql := `CREATE REPOSITORY broker_repo
		WITH BROKER my_broker
		ON LOCATION "hdfs://namenode:8020/backup"
		PROPERTIES(
			"username" = "hdfsuser",
			"password" = "secret"
		)`
	file, errs := parseForTest(sql)
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
		WITH BROKER
		ON LOCATION "s3://bucket/path"
		PROPERTIES("s3.endpoint"="s3.amazonaws.com","s3.access_key"="ak","s3.secret_key"="sk","s3.region"="us-east-1")`
	file, errs := parseForTest(sql)
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
	sql := `CREATE REPOSITORY r WITH BROKER ON LOCATION "s3://b/p" PROPERTIES("s3.endpoint"="ep","s3.access_key"="ak","s3.secret_key"="sk","s3.region"="us-east-1")`
	file, errs := parseForTest(sql)
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
	file, errs := parseForTest(sql)
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
	file, errs := parseForTest(`ALTER REPOSITORY r PROPERTIES("s3.access_key"="ak")`)
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
	file, errs := parseForTest("DROP REPOSITORY my_repo")
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
	file, errs := parseForTest("DROP REPOSITORY r")
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

func TestCreateFile_Basic(t *testing.T) {
	sql := `CREATE FILE "ca.pem" PROPERTIES("url"="https://example.com/ca.pem","type"="ca")`
	file, errs := parseForTest(sql)
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
	if stmt.Catalog != nil {
		t.Errorf("Catalog = %v, want nil", stmt.Catalog)
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateFile_WithDatabase(t *testing.T) {
	sql := `CREATE FILE "my_cert.pem" IN mydb PROPERTIES("url"="https://example.com/cert.pem","type"="ca")`
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateFileStmt)
	if stmt.Name != "my_cert.pem" {
		t.Errorf("Name = %q, want %q", stmt.Name, "my_cert.pem")
	}
	if stmt.Catalog == nil || stmt.Catalog.String() != "mydb" {
		t.Errorf("Catalog = %v, want mydb", stmt.Catalog)
	}
}

func TestCreateFile_Tag(t *testing.T) {
	sql := `CREATE FILE "f.pem" PROPERTIES("url"="https://x.com/f.pem","type"="ca")`
	file, errs := parseForTest(sql)
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
	file, errs := parseForTest(sql)
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
	if stmt.Catalog != nil {
		t.Errorf("Catalog = %v, want nil", stmt.Catalog)
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
}

func TestDropFile_WithDatabase(t *testing.T) {
	sql := `DROP FILE "ca.pem" FROM mydb PROPERTIES("type"="ca")`
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropFileStmt)
	if stmt.Catalog == nil || stmt.Catalog.String() != "mydb" {
		t.Errorf("Catalog = %v, want mydb", stmt.Catalog)
	}
}

func TestDropFile_Tag(t *testing.T) {
	sql := `DROP FILE "f.pem" PROPERTIES("type"="ca")`
	file, errs := parseForTest(sql)
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
		file, errs := parseForTest(tt.sql)
		if len(errs) != 0 {
			t.Errorf("sql %q: unexpected errors: %v", tt.sql, errs)
			continue
		}
		if len(file.Stmts) != 1 {
			t.Errorf("sql %q: expected 1 stmt, got %d", tt.sql, len(file.Stmts))
		}
	}
}

func TestStorage_MultiStatement(t *testing.T) {
	input := `CREATE STORAGE POLICY p PROPERTIES("storage_resource"="s3");
DROP STORAGE POLICY p`
	file, errs := parseForTest(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 2 {
		t.Fatalf("expected 2 stmts, got %d", len(file.Stmts))
	}
	if _, ok := file.Stmts[0].(*ast.CreateStoragePolicyStmt); !ok {
		t.Errorf("Stmts[0]: expected CreateStoragePolicyStmt, got %T", file.Stmts[0])
	}
	if _, ok := file.Stmts[1].(*ast.DropStoragePolicyStmt); !ok {
		t.Errorf("Stmts[1]: expected DropStoragePolicyStmt, got %T", file.Stmts[1])
	}
}

// ---------------------------------------------------------------------------
// STORAGE VOLUME
// ---------------------------------------------------------------------------

// TestCreateStorageVolume: createStorageVolumeStatement (StarRocks.g4:782).
func TestCreateStorageVolume(t *testing.T) {
	stmt := parseOne(t, `CREATE STORAGE VOLUME IF NOT EXISTS vol_1 TYPE = S3
		LOCATIONS = ('s3://bucket/path', 's3://bucket/other')
		COMMENT 'shared'
		PROPERTIES ("aws.s3.use_aws_sdk_default_behavior" = "true")`).(*ast.CreateStorageVolumeStmt)

	if stmt.Name != "vol_1" || !stmt.IfNotExists {
		t.Errorf("Name = %q IfNotExists = %v, want vol_1 and true", stmt.Name, stmt.IfNotExists)
	}
	if stmt.Type != "S3" {
		t.Errorf("Type = %q, want S3", stmt.Type)
	}
	if len(stmt.Locations) != 2 || stmt.Locations[0] != "s3://bucket/path" || stmt.Locations[1] != "s3://bucket/other" {
		t.Errorf("Locations = %v, want the two stated", stmt.Locations)
	}
	if stmt.Comment != "shared" {
		t.Errorf("Comment = %q, want shared", stmt.Comment)
	}
	if len(stmt.Properties) != 1 || stmt.Properties[0].Key != "aws.s3.use_aws_sdk_default_behavior" {
		t.Errorf("Properties = %+v, want the stated one", stmt.Properties)
	}
	if stmt.Tag() != ast.T_CreateStorageVolumeStmt {
		t.Errorf("Tag() = %v, want T_CreateStorageVolumeStmt", stmt.Tag())
	}
	if end, want := stmt.Loc.End, len(`CREATE STORAGE VOLUME IF NOT EXISTS vol_1 TYPE = S3
		LOCATIONS = ('s3://bucket/path', 's3://bucket/other')
		COMMENT 'shared'
		PROPERTIES ("aws.s3.use_aws_sdk_default_behavior" = "true")`); end != want {
		t.Errorf("Loc.End = %d, want %d", end, want)
	}

	// The shortest form states no comment and no properties.
	short := parseOne(t, `CREATE STORAGE VOLUME v TYPE = HDFS LOCATIONS = ('hdfs://h/p')`).(*ast.CreateStorageVolumeStmt)
	if short.Comment != "" || short.Properties != nil {
		t.Errorf("Comment = %q Properties = %v, want neither", short.Comment, short.Properties)
	}
}

// TestCreateStorageVolumeRejected: TYPE and LOCATIONS sit outside every
// optional group, so a statement that leaves one out is refused.
func TestCreateStorageVolumeRejected(t *testing.T) {
	for _, sql := range []string{
		`CREATE STORAGE VOLUME v LOCATIONS = ('hdfs://h/p')`,
		`CREATE STORAGE VOLUME v TYPE = HDFS`,
		`CREATE STORAGE VOLUME v TYPE HDFS LOCATIONS = ('hdfs://h/p')`,
		`CREATE STORAGE VOLUME v TYPE = HDFS LOCATIONS = ()`,
		`CREATE STORAGE VOLUME v TYPE = HDFS LOCATIONS = 'hdfs://h/p'`,
		`CREATE STORAGE VOLUME v TYPE = HDFS LOCATIONS = ('hdfs://h/p') COMMENT`,
		`CREATE STORAGE VOLUME v TYPE = HDFS LOCATIONS = ('hdfs://h/p') zzz`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed; the engine rejects it", sql)
		}
	}
}

// TestAlterStorageVolume: alterStorageVolumeStatement (StarRocks.g4:803)
// keeps both clause forms, in the order they were written, comma-separated.
func TestAlterStorageVolume(t *testing.T) {
	stmt := parseOne(t, `ALTER STORAGE VOLUME v COMMENT = 'one', SET ("aws.s3.region" = "us-east-1"), COMMENT = 'two'`).(*ast.AlterStorageVolumeStmt)

	if stmt.Name != "v" {
		t.Errorf("Name = %q, want v", stmt.Name)
	}
	if len(stmt.Clauses) != 3 {
		t.Fatalf("Clauses = %d, want 3", len(stmt.Clauses))
	}
	want := []ast.AlterStorageVolumeClauseType{
		ast.StorageVolumeComment, ast.StorageVolumeProperties, ast.StorageVolumeComment,
	}
	for i, c := range stmt.Clauses {
		if c.Type != want[i] {
			t.Errorf("Clauses[%d].Type = %v, want %v", i, c.Type, want[i])
		}
	}
	if stmt.Clauses[0].Comment != "one" || stmt.Clauses[2].Comment != "two" {
		t.Errorf("comments = %q, %q, want one and two", stmt.Clauses[0].Comment, stmt.Clauses[2].Comment)
	}
	if props := stmt.Clauses[1].Properties; len(props) != 1 || props[0].Key != "aws.s3.region" {
		t.Errorf("Clauses[1].Properties = %+v, want the stated one", props)
	}
	if stmt.Tag() != ast.T_AlterStorageVolumeStmt {
		t.Errorf("Tag() = %v, want T_AlterStorageVolumeStmt", stmt.Tag())
	}
}

// TestAlterStorageVolumeRejected: the '=' of a COMMENT clause is obligatory
// (StarRocks.g4:816), and an ALTER names at least one clause.
func TestAlterStorageVolumeRejected(t *testing.T) {
	for _, sql := range []string{
		`ALTER STORAGE VOLUME v`,
		`ALTER STORAGE VOLUME v COMMENT 'a'`,
		`ALTER STORAGE VOLUME v COMMENT =`,
		`ALTER STORAGE VOLUME v SET`,
		`ALTER STORAGE VOLUME v SET ()`,
		`ALTER STORAGE VOLUME v COMMENT = 'a',`,
		`ALTER STORAGE VOLUME v COMMENT = 'a' COMMENT = 'b'`,
		`ALTER STORAGE VOLUME v COMMENT = 'a' zzz`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed; the engine rejects it", sql)
		}
	}
}

// TestDropDescStorageVolume: dropStorageVolumeStatement (StarRocks.g4:799) and
// descStorageVolumeStatement (StarRocks.g4:820).
func TestDropDescStorageVolume(t *testing.T) {
	drop := parseOne(t, `DROP STORAGE VOLUME IF EXISTS v`).(*ast.DropStorageVolumeStmt)
	if drop.Name != "v" || !drop.IfExists {
		t.Errorf("Name = %q IfExists = %v, want v and true", drop.Name, drop.IfExists)
	}

	for _, sql := range []string{`DESC STORAGE VOLUME v`, `DESCRIBE STORAGE VOLUME v`, `DESC STORAGE VOLUME 'v'`} {
		desc := parseOne(t, sql).(*ast.DescStorageVolumeStmt)
		if desc.Name != "v" {
			t.Errorf("%q: Name = %q, want v", sql, desc.Name)
		}
	}

	// A table may be named storage: DESC STORAGE is a describe, not a
	// storage-volume description missing VOLUME.
	if _, ok := parseOne(t, `DESC storage`).(*ast.DescribeStmt); !ok {
		t.Error("DESC storage should describe a table named storage")
	}
	for _, sql := range []string{`DROP STORAGE VOLUME`, `DESC STORAGE VOLUME`, `DESC STORAGE a`, `DESC STORAGE VOLUME v zzz`} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed; the engine rejects it", sql)
		}
	}
}
