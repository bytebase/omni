package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func testParseCreateDatabase(input string) (*ast.CreateDatabaseStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.CreateDatabaseStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a CreateDatabaseStmt"})
	}
	return stmt, result.Errors
}

func testParseCreateSchema(input string) (*ast.CreateSchemaStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.CreateSchemaStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a CreateSchemaStmt"})
	}
	return stmt, result.Errors
}

func testParseAlterDatabase(input string) (*ast.AlterDatabaseStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.AlterDatabaseStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not an AlterDatabaseStmt"})
	}
	return stmt, result.Errors
}

func testParseAlterSchema(input string) (*ast.AlterSchemaStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.AlterSchemaStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not an AlterSchemaStmt"})
	}
	return stmt, result.Errors
}

func testParseDropDatabase(input string) (*ast.DropDatabaseStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.DropDatabaseStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a DropDatabaseStmt"})
	}
	return stmt, result.Errors
}

func testParseDropSchema(input string) (*ast.DropSchemaStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.DropSchemaStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a DropSchemaStmt"})
	}
	return stmt, result.Errors
}

// ---------------------------------------------------------------------------
// CREATE DATABASE tests
// ---------------------------------------------------------------------------

func TestCreateDatabase_Basic(t *testing.T) {
	stmt, errs := testParseCreateDatabase("CREATE DATABASE mydb")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt == nil {
		t.Fatal("expected CreateDatabaseStmt, got nil")
	}
	if stmt.Name.Normalize() != "MYDB" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "MYDB")
	}
	if stmt.OrReplace {
		t.Error("expected OrReplace=false")
	}
	if stmt.Transient {
		t.Error("expected Transient=false")
	}
	if stmt.IfNotExists {
		t.Error("expected IfNotExists=false")
	}
}

func TestCreateDatabase_OrReplace(t *testing.T) {
	stmt, errs := testParseCreateDatabase("CREATE OR REPLACE DATABASE mydb")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace=true")
	}
}

func TestCreateDatabase_Transient(t *testing.T) {
	stmt, errs := testParseCreateDatabase("CREATE TRANSIENT DATABASE mydb")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Transient {
		t.Error("expected Transient=true")
	}
}

func TestCreateDatabase_OrReplaceTransient(t *testing.T) {
	stmt, errs := testParseCreateDatabase("CREATE OR REPLACE TRANSIENT DATABASE mydb")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace=true")
	}
	if !stmt.Transient {
		t.Error("expected Transient=true")
	}
}

func TestCreateDatabase_IfNotExists(t *testing.T) {
	stmt, errs := testParseCreateDatabase("CREATE DATABASE IF NOT EXISTS mydb")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists=true")
	}
	if stmt.Name.Normalize() != "MYDB" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "MYDB")
	}
}

func TestCreateDatabase_WithClone(t *testing.T) {
	stmt, errs := testParseCreateDatabase("CREATE DATABASE mydb CLONE sourcedb")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Clone == nil {
		t.Fatal("expected Clone to be set")
	}
	if stmt.Clone.Source.Normalize() != "SOURCEDB" {
		t.Errorf("clone source = %q, want %q", stmt.Clone.Source.Normalize(), "SOURCEDB")
	}
}

func TestCreateDatabase_WithCloneAtTimestamp(t *testing.T) {
	stmt, errs := testParseCreateDatabase(
		`CREATE DATABASE mydb CLONE sourcedb AT (TIMESTAMP => '2024-01-01 00:00:00')`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Clone == nil {
		t.Fatal("expected Clone to be set")
	}
	if stmt.Clone.AtBefore != "AT" {
		t.Errorf("AtBefore = %q, want %q", stmt.Clone.AtBefore, "AT")
	}
	if stmt.Clone.Kind != "TIMESTAMP" {
		t.Errorf("Kind = %q, want %q", stmt.Clone.Kind, "TIMESTAMP")
	}
}

func TestCreateDatabase_WithDataRetention(t *testing.T) {
	stmt, errs := testParseCreateDatabase(
		"CREATE DATABASE mydb DATA_RETENTION_TIME_IN_DAYS = 7")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Props.DataRetention == nil {
		t.Fatal("expected DataRetention to be set")
	}
	if *stmt.Props.DataRetention != 7 {
		t.Errorf("DataRetention = %d, want 7", *stmt.Props.DataRetention)
	}
}

func TestCreateDatabase_WithComment(t *testing.T) {
	stmt, errs := testParseCreateDatabase(
		"CREATE DATABASE mydb COMMENT = 'my database'")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Props.Comment == nil {
		t.Fatal("expected Comment to be set")
	}
	if *stmt.Props.Comment != "my database" {
		t.Errorf("Comment = %q, want %q", *stmt.Props.Comment, "my database")
	}
}

func TestCreateDatabase_WithTags(t *testing.T) {
	stmt, errs := testParseCreateDatabase(
		"CREATE DATABASE mydb WITH TAG (env = 'prod')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Tags) != 1 {
		t.Fatalf("tags = %d, want 1", len(stmt.Tags))
	}
	if stmt.Tags[0].Name.Normalize() != "ENV" {
		t.Errorf("tag name = %q, want %q", stmt.Tags[0].Name.Normalize(), "ENV")
	}
	if stmt.Tags[0].Value != "prod" {
		t.Errorf("tag value = %q, want %q", stmt.Tags[0].Value, "prod")
	}
}

func TestCreateDatabase_FullProps(t *testing.T) {
	stmt, errs := testParseCreateDatabase(
		`CREATE OR REPLACE DATABASE mydb
		 DATA_RETENTION_TIME_IN_DAYS = 30
		 MAX_DATA_EXTENSION_TIME_IN_DAYS = 90
		 DEFAULT_DDL_COLLATION = 'en_US'
		 COMMENT = 'full test'`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Props.DataRetention == nil || *stmt.Props.DataRetention != 30 {
		t.Errorf("DataRetention = %v, want 30", stmt.Props.DataRetention)
	}
	if stmt.Props.MaxDataExt == nil || *stmt.Props.MaxDataExt != 90 {
		t.Errorf("MaxDataExt = %v, want 90", stmt.Props.MaxDataExt)
	}
	if stmt.Props.DefaultDDLCol == nil || *stmt.Props.DefaultDDLCol != "en_US" {
		t.Errorf("DefaultDDLCol = %v, want en_US", stmt.Props.DefaultDDLCol)
	}
	if stmt.Props.Comment == nil || *stmt.Props.Comment != "full test" {
		t.Errorf("Comment = %v, want 'full test'", stmt.Props.Comment)
	}
}

// ---------------------------------------------------------------------------
// CREATE SCHEMA tests
// ---------------------------------------------------------------------------

func TestCreateSchema_Basic(t *testing.T) {
	stmt, errs := testParseCreateSchema("CREATE SCHEMA myschema")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt == nil {
		t.Fatal("expected CreateSchemaStmt, got nil")
	}
	if stmt.Name.Normalize() != "MYSCHEMA" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "MYSCHEMA")
	}
}

func TestCreateSchema_OrReplace(t *testing.T) {
	stmt, errs := testParseCreateSchema("CREATE OR REPLACE SCHEMA myschema")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace=true")
	}
}

func TestCreateSchema_Transient(t *testing.T) {
	stmt, errs := testParseCreateSchema("CREATE TRANSIENT SCHEMA myschema")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Transient {
		t.Error("expected Transient=true")
	}
}

func TestCreateSchema_IfNotExists(t *testing.T) {
	stmt, errs := testParseCreateSchema("CREATE SCHEMA IF NOT EXISTS myschema")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists=true")
	}
}

func TestCreateSchema_WithClone(t *testing.T) {
	stmt, errs := testParseCreateSchema("CREATE SCHEMA s2 CLONE s1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Clone == nil {
		t.Fatal("expected Clone to be set")
	}
	if stmt.Clone.Source.Normalize() != "S1" {
		t.Errorf("clone source = %q, want %q", stmt.Clone.Source.Normalize(), "S1")
	}
}

func TestCreateSchema_ManagedAccess(t *testing.T) {
	stmt, errs := testParseCreateSchema("CREATE SCHEMA myschema WITH MANAGED ACCESS")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.ManagedAccess {
		t.Error("expected ManagedAccess=true")
	}
}

func TestCreateSchema_ManagedAccessWithProps(t *testing.T) {
	stmt, errs := testParseCreateSchema(
		"CREATE SCHEMA myschema WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 7")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.ManagedAccess {
		t.Error("expected ManagedAccess=true")
	}
	if stmt.Props.DataRetention == nil || *stmt.Props.DataRetention != 7 {
		t.Errorf("DataRetention = %v, want 7", stmt.Props.DataRetention)
	}
}

func TestCreateSchema_WithTags(t *testing.T) {
	stmt, errs := testParseCreateSchema(
		"CREATE SCHEMA myschema WITH TAG (team = 'analytics')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Tags) != 1 {
		t.Fatalf("tags = %d, want 1", len(stmt.Tags))
	}
	if stmt.Tags[0].Value != "analytics" {
		t.Errorf("tag value = %q, want %q", stmt.Tags[0].Value, "analytics")
	}
}

func TestCreateSchema_QualifiedName(t *testing.T) {
	stmt, errs := testParseCreateSchema("CREATE SCHEMA mydb.myschema")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Name.Normalize() != "MYDB.MYSCHEMA" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "MYDB.MYSCHEMA")
	}
}

// ---------------------------------------------------------------------------
// ALTER DATABASE tests
// ---------------------------------------------------------------------------

func TestAlterDatabase_RenameTo(t *testing.T) {
	stmt, errs := testParseAlterDatabase("ALTER DATABASE mydb RENAME TO newdb")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt == nil {
		t.Fatal("expected AlterDatabaseStmt, got nil")
	}
	if stmt.Action != ast.AlterDBRename {
		t.Errorf("action = %v, want AlterDBRename", stmt.Action)
	}
	if stmt.Name.Normalize() != "MYDB" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "MYDB")
	}
	if stmt.NewName == nil || stmt.NewName.Normalize() != "NEWDB" {
		t.Errorf("newName = %v, want NEWDB", stmt.NewName)
	}
}

func TestAlterDatabase_SwapWith(t *testing.T) {
	stmt, errs := testParseAlterDatabase("ALTER DATABASE db1 SWAP WITH db2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterDBSwap {
		t.Errorf("action = %v, want AlterDBSwap", stmt.Action)
	}
	if stmt.NewName == nil || stmt.NewName.Normalize() != "DB2" {
		t.Errorf("newName = %v, want DB2", stmt.NewName)
	}
}

func TestAlterDatabase_SetComment(t *testing.T) {
	stmt, errs := testParseAlterDatabase("ALTER DATABASE mydb SET COMMENT = 'new comment'")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterDBSet {
		t.Errorf("action = %v, want AlterDBSet", stmt.Action)
	}
	if stmt.SetProps == nil {
		t.Fatal("SetProps is nil")
	}
	if stmt.SetProps.Comment == nil || *stmt.SetProps.Comment != "new comment" {
		t.Errorf("comment = %v, want 'new comment'", stmt.SetProps.Comment)
	}
}

func TestAlterDatabase_SetDataRetention(t *testing.T) {
	stmt, errs := testParseAlterDatabase(
		"ALTER DATABASE mydb SET DATA_RETENTION_TIME_IN_DAYS = 14")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterDBSet {
		t.Errorf("action = %v, want AlterDBSet", stmt.Action)
	}
	if stmt.SetProps.DataRetention == nil || *stmt.SetProps.DataRetention != 14 {
		t.Errorf("DataRetention = %v, want 14", stmt.SetProps.DataRetention)
	}
}

func TestAlterDatabase_SetTag(t *testing.T) {
	stmt, errs := testParseAlterDatabase(
		"ALTER DATABASE mydb SET TAG (env = 'production')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterDBSetTag {
		t.Errorf("action = %v, want AlterDBSetTag", stmt.Action)
	}
	if len(stmt.Tags) != 1 {
		t.Fatalf("tags = %d, want 1", len(stmt.Tags))
	}
	if stmt.Tags[0].Value != "production" {
		t.Errorf("tag value = %q, want production", stmt.Tags[0].Value)
	}
}

func TestAlterDatabase_UnsetComment(t *testing.T) {
	stmt, errs := testParseAlterDatabase("ALTER DATABASE mydb UNSET COMMENT")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterDBUnset {
		t.Errorf("action = %v, want AlterDBUnset", stmt.Action)
	}
	if len(stmt.UnsetProps) != 1 || stmt.UnsetProps[0] != "COMMENT" {
		t.Errorf("UnsetProps = %v, want [COMMENT]", stmt.UnsetProps)
	}
}

func TestAlterDatabase_UnsetMultiple(t *testing.T) {
	stmt, errs := testParseAlterDatabase(
		"ALTER DATABASE mydb UNSET DATA_RETENTION_TIME_IN_DAYS, COMMENT")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterDBUnset {
		t.Errorf("action = %v, want AlterDBUnset", stmt.Action)
	}
	if len(stmt.UnsetProps) != 2 {
		t.Fatalf("UnsetProps = %d, want 2", len(stmt.UnsetProps))
	}
}

func TestAlterDatabase_UnsetTag(t *testing.T) {
	stmt, errs := testParseAlterDatabase("ALTER DATABASE mydb UNSET TAG (env)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterDBUnsetTag {
		t.Errorf("action = %v, want AlterDBUnsetTag", stmt.Action)
	}
	if len(stmt.UnsetTags) != 1 {
		t.Fatalf("UnsetTags = %d, want 1", len(stmt.UnsetTags))
	}
}

func TestAlterDatabase_IfExists(t *testing.T) {
	stmt, errs := testParseAlterDatabase("ALTER DATABASE IF EXISTS mydb RENAME TO newdb")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

func TestAlterDatabase_EnableReplication(t *testing.T) {
	stmt, errs := testParseAlterDatabase(
		"ALTER DATABASE mydb ENABLE REPLICATION TO ACCOUNTS org1.account1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterDBEnableReplication {
		t.Errorf("action = %v, want AlterDBEnableReplication", stmt.Action)
	}
}

func TestAlterDatabase_DisableReplication(t *testing.T) {
	stmt, errs := testParseAlterDatabase("ALTER DATABASE mydb DISABLE REPLICATION")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterDBDisableReplication {
		t.Errorf("action = %v, want AlterDBDisableReplication", stmt.Action)
	}
}

func TestAlterDatabase_Refresh(t *testing.T) {
	stmt, errs := testParseAlterDatabase("ALTER DATABASE mydb REFRESH")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterDBRefresh {
		t.Errorf("action = %v, want AlterDBRefresh", stmt.Action)
	}
}

func TestAlterDatabase_Primary(t *testing.T) {
	stmt, errs := testParseAlterDatabase("ALTER DATABASE mydb PRIMARY")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterDBPrimary {
		t.Errorf("action = %v, want AlterDBPrimary", stmt.Action)
	}
}

// ---------------------------------------------------------------------------
// ALTER SCHEMA tests
// ---------------------------------------------------------------------------

func TestAlterSchema_RenameTo(t *testing.T) {
	stmt, errs := testParseAlterSchema("ALTER SCHEMA s1 RENAME TO s2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt == nil {
		t.Fatal("expected AlterSchemaStmt, got nil")
	}
	if stmt.Action != ast.AlterSchemaRename {
		t.Errorf("action = %v, want AlterSchemaRename", stmt.Action)
	}
	if stmt.NewName == nil || stmt.NewName.Normalize() != "S2" {
		t.Errorf("newName = %v, want S2", stmt.NewName)
	}
}

func TestAlterSchema_SwapWith(t *testing.T) {
	stmt, errs := testParseAlterSchema("ALTER SCHEMA s1 SWAP WITH s2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterSchemaSwap {
		t.Errorf("action = %v, want AlterSchemaSwap", stmt.Action)
	}
}

func TestAlterSchema_SetComment(t *testing.T) {
	stmt, errs := testParseAlterSchema("ALTER SCHEMA myschema SET COMMENT = 'test schema'")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterSchemaSet {
		t.Errorf("action = %v, want AlterSchemaSet", stmt.Action)
	}
	if stmt.SetProps.Comment == nil || *stmt.SetProps.Comment != "test schema" {
		t.Errorf("comment = %v", stmt.SetProps.Comment)
	}
}

func TestAlterSchema_SetDataRetention(t *testing.T) {
	stmt, errs := testParseAlterSchema(
		"ALTER SCHEMA myschema SET DATA_RETENTION_TIME_IN_DAYS = 5")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterSchemaSet {
		t.Errorf("action = %v, want AlterSchemaSet", stmt.Action)
	}
	if stmt.SetProps.DataRetention == nil || *stmt.SetProps.DataRetention != 5 {
		t.Errorf("DataRetention = %v, want 5", stmt.SetProps.DataRetention)
	}
}

func TestAlterSchema_UnsetComment(t *testing.T) {
	stmt, errs := testParseAlterSchema("ALTER SCHEMA myschema UNSET COMMENT")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterSchemaUnset {
		t.Errorf("action = %v, want AlterSchemaUnset", stmt.Action)
	}
	if len(stmt.UnsetProps) != 1 || stmt.UnsetProps[0] != "COMMENT" {
		t.Errorf("UnsetProps = %v, want [COMMENT]", stmt.UnsetProps)
	}
}

func TestAlterSchema_SetTag(t *testing.T) {
	stmt, errs := testParseAlterSchema("ALTER SCHEMA myschema SET TAG (env = 'dev')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterSchemaSetTag {
		t.Errorf("action = %v, want AlterSchemaSetTag", stmt.Action)
	}
	if len(stmt.Tags) != 1 || stmt.Tags[0].Value != "dev" {
		t.Errorf("Tags = %v", stmt.Tags)
	}
}

func TestAlterSchema_UnsetTag(t *testing.T) {
	stmt, errs := testParseAlterSchema("ALTER SCHEMA myschema UNSET TAG (env)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterSchemaUnsetTag {
		t.Errorf("action = %v, want AlterSchemaUnsetTag", stmt.Action)
	}
	if len(stmt.UnsetTags) != 1 {
		t.Fatalf("UnsetTags = %d, want 1", len(stmt.UnsetTags))
	}
}

func TestAlterSchema_EnableManagedAccess(t *testing.T) {
	stmt, errs := testParseAlterSchema("ALTER SCHEMA myschema ENABLE MANAGED ACCESS")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterSchemaEnableManagedAccess {
		t.Errorf("action = %v, want AlterSchemaEnableManagedAccess", stmt.Action)
	}
}

func TestAlterSchema_DisableManagedAccess(t *testing.T) {
	stmt, errs := testParseAlterSchema("ALTER SCHEMA myschema DISABLE MANAGED ACCESS")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterSchemaDisableManagedAccess {
		t.Errorf("action = %v, want AlterSchemaDisableManagedAccess", stmt.Action)
	}
}

func TestAlterSchema_IfExists(t *testing.T) {
	stmt, errs := testParseAlterSchema("ALTER SCHEMA IF EXISTS myschema RENAME TO newschema")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

func TestAlterSchema_QualifiedName(t *testing.T) {
	stmt, errs := testParseAlterSchema("ALTER SCHEMA mydb.s1 RENAME TO mydb.s2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Name.Normalize() != "MYDB.S1" {
		t.Errorf("name = %q, want MYDB.S1", stmt.Name.Normalize())
	}
}

// ---------------------------------------------------------------------------
// DROP DATABASE tests
// ---------------------------------------------------------------------------

func TestDropDatabase_Basic(t *testing.T) {
	stmt, errs := testParseDropDatabase("DROP DATABASE mydb")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt == nil {
		t.Fatal("expected DropDatabaseStmt, got nil")
	}
	if stmt.Name.Normalize() != "MYDB" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "MYDB")
	}
	if stmt.IfExists {
		t.Error("expected IfExists=false")
	}
	if stmt.Cascade {
		t.Error("expected Cascade=false")
	}
	if stmt.Restrict {
		t.Error("expected Restrict=false")
	}
}

func TestDropDatabase_IfExists(t *testing.T) {
	stmt, errs := testParseDropDatabase("DROP DATABASE IF EXISTS mydb")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

func TestDropDatabase_Cascade(t *testing.T) {
	stmt, errs := testParseDropDatabase("DROP DATABASE mydb CASCADE")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Cascade {
		t.Error("expected Cascade=true")
	}
}

func TestDropDatabase_Restrict(t *testing.T) {
	stmt, errs := testParseDropDatabase("DROP DATABASE mydb RESTRICT")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Restrict {
		t.Error("expected Restrict=true")
	}
}

func TestDropDatabase_IfExistsCascade(t *testing.T) {
	stmt, errs := testParseDropDatabase("DROP DATABASE IF EXISTS mydb CASCADE")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
	if !stmt.Cascade {
		t.Error("expected Cascade=true")
	}
}

// ---------------------------------------------------------------------------
// DROP SCHEMA tests
// ---------------------------------------------------------------------------

func TestDropSchema_Basic(t *testing.T) {
	stmt, errs := testParseDropSchema("DROP SCHEMA myschema")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt == nil {
		t.Fatal("expected DropSchemaStmt, got nil")
	}
	if stmt.Name.Normalize() != "MYSCHEMA" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "MYSCHEMA")
	}
}

func TestDropSchema_IfExists(t *testing.T) {
	stmt, errs := testParseDropSchema("DROP SCHEMA IF EXISTS myschema")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

func TestDropSchema_Cascade(t *testing.T) {
	stmt, errs := testParseDropSchema("DROP SCHEMA mydb.myschema CASCADE")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Cascade {
		t.Error("expected Cascade=true")
	}
	if stmt.Name.Normalize() != "MYDB.MYSCHEMA" {
		t.Errorf("name = %q, want MYDB.MYSCHEMA", stmt.Name.Normalize())
	}
}

func TestDropSchema_Restrict(t *testing.T) {
	stmt, errs := testParseDropSchema("DROP SCHEMA myschema RESTRICT")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Restrict {
		t.Error("expected Restrict=true")
	}
}

// ---------------------------------------------------------------------------
// UNDROP DATABASE tests
// ---------------------------------------------------------------------------

func TestUndropDatabase_Basic(t *testing.T) {
	result := ParseBestEffort("UNDROP DATABASE mydb")
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if len(result.File.Stmts) == 0 {
		t.Fatal("expected a statement")
	}
	stmt, ok := result.File.Stmts[0].(*ast.UndropDatabaseStmt)
	if !ok {
		t.Fatalf("expected *UndropDatabaseStmt, got %T", result.File.Stmts[0])
	}
	if stmt.Name.Normalize() != "MYDB" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "MYDB")
	}
}

func TestUndropDatabase_Tag(t *testing.T) {
	stmt, ok := ParseBestEffort("UNDROP DATABASE mydb").File.Stmts[0].(*ast.UndropDatabaseStmt)
	if !ok {
		t.Fatal("wrong node type")
	}
	if stmt.Tag() != ast.T_UndropDatabaseStmt {
		t.Errorf("tag = %v, want T_UndropDatabaseStmt", stmt.Tag())
	}
}

// ---------------------------------------------------------------------------
// UNDROP SCHEMA tests
// ---------------------------------------------------------------------------

func TestUndropSchema_Basic(t *testing.T) {
	result := ParseBestEffort("UNDROP SCHEMA myschema")
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if len(result.File.Stmts) == 0 {
		t.Fatal("expected a statement")
	}
	stmt, ok := result.File.Stmts[0].(*ast.UndropSchemaStmt)
	if !ok {
		t.Fatalf("expected *UndropSchemaStmt, got %T", result.File.Stmts[0])
	}
	if stmt.Name.Normalize() != "MYSCHEMA" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "MYSCHEMA")
	}
}

func TestUndropSchema_QualifiedName(t *testing.T) {
	result := ParseBestEffort("UNDROP SCHEMA mydb.myschema")
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	stmt, ok := result.File.Stmts[0].(*ast.UndropSchemaStmt)
	if !ok {
		t.Fatal("wrong node type")
	}
	if stmt.Name.Normalize() != "MYDB.MYSCHEMA" {
		t.Errorf("name = %q, want MYDB.MYSCHEMA", stmt.Name.Normalize())
	}
}

// ---------------------------------------------------------------------------
// NodeTag tests
// ---------------------------------------------------------------------------

func TestNodeTags_DatabaseSchema(t *testing.T) {
	checks := []struct {
		sql  string
		want ast.NodeTag
	}{
		{"CREATE DATABASE d", ast.T_CreateDatabaseStmt},
		{"CREATE SCHEMA s", ast.T_CreateSchemaStmt},
		{"ALTER DATABASE d RENAME TO d2", ast.T_AlterDatabaseStmt},
		{"ALTER SCHEMA s RENAME TO s2", ast.T_AlterSchemaStmt},
		{"DROP DATABASE d", ast.T_DropDatabaseStmt},
		{"DROP SCHEMA s", ast.T_DropSchemaStmt},
		{"UNDROP DATABASE d", ast.T_UndropDatabaseStmt},
		{"UNDROP SCHEMA s", ast.T_UndropSchemaStmt},
	}

	for _, c := range checks {
		result := ParseBestEffort(c.sql)
		if len(result.Errors) > 0 {
			t.Errorf("sql=%q errors=%v", c.sql, result.Errors)
			continue
		}
		if len(result.File.Stmts) == 0 {
			t.Errorf("sql=%q: no statements", c.sql)
			continue
		}
		got := result.File.Stmts[0].Tag()
		if got != c.want {
			t.Errorf("sql=%q: tag=%v, want %v", c.sql, got, c.want)
		}
	}
}
