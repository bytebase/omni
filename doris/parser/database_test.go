package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// CREATE DATABASE
// ---------------------------------------------------------------------------

func TestCreateDatabase_Basic(t *testing.T) {
	file, errs := Parse("CREATE DATABASE mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateDatabaseStmt, got %T", file.Stmts[0])
	}
	if stmt.Name == nil {
		t.Fatal("Name is nil")
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if len(stmt.Properties) != 0 {
		t.Errorf("expected 0 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateDatabase_IfNotExists(t *testing.T) {
	file, errs := Parse("CREATE DATABASE IF NOT EXISTS mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
}

func TestCreateDatabase_WithProperties(t *testing.T) {
	file, errs := Parse(`CREATE DATABASE mydb PROPERTIES("key"="value")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "key" {
		t.Errorf("Property key = %q, want %q", stmt.Properties[0].Key, "key")
	}
	if stmt.Properties[0].Value != "value" {
		t.Errorf("Property value = %q, want %q", stmt.Properties[0].Value, "value")
	}
}

func TestCreateDatabase_WithMultipleProperties(t *testing.T) {
	file, errs := Parse(`CREATE DATABASE mydb PROPERTIES("replication_num"="3", "storage_medium"="SSD")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "replication_num" {
		t.Errorf("Properties[0].Key = %q, want %q", stmt.Properties[0].Key, "replication_num")
	}
	if stmt.Properties[0].Value != "3" {
		t.Errorf("Properties[0].Value = %q, want %q", stmt.Properties[0].Value, "3")
	}
	if stmt.Properties[1].Key != "storage_medium" {
		t.Errorf("Properties[1].Key = %q, want %q", stmt.Properties[1].Key, "storage_medium")
	}
}

func TestCreateDatabase_IfNotExistsWithProperties(t *testing.T) {
	file, errs := Parse(`CREATE DATABASE IF NOT EXISTS mydb PROPERTIES("key"="val")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
}

// ---------------------------------------------------------------------------
// CREATE SCHEMA (alias for CREATE DATABASE)
// ---------------------------------------------------------------------------

func TestCreateSchema_Basic(t *testing.T) {
	file, errs := Parse("CREATE SCHEMA myschema")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateDatabaseStmt, got %T", file.Stmts[0])
	}
	if stmt.Name.String() != "myschema" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "myschema")
	}
}

func TestCreateSchema_IfNotExists(t *testing.T) {
	file, errs := Parse("CREATE SCHEMA IF NOT EXISTS myschema")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
}

// ---------------------------------------------------------------------------
// DROP DATABASE
// ---------------------------------------------------------------------------

func TestDropDatabase_Basic(t *testing.T) {
	file, errs := Parse("DROP DATABASE mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropDatabaseStmt)
	if !ok {
		t.Fatalf("expected *ast.DropDatabaseStmt, got %T", file.Stmts[0])
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
	if stmt.Force {
		t.Error("Force should be false")
	}
}

func TestDropDatabase_IfExists(t *testing.T) {
	file, errs := Parse("DROP DATABASE IF EXISTS mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropDatabaseStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
}

func TestDropDatabase_Force(t *testing.T) {
	file, errs := Parse("DROP DATABASE mydb FORCE")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropDatabaseStmt)
	if !stmt.Force {
		t.Error("Force should be true")
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
}

func TestDropDatabase_IfExistsForce(t *testing.T) {
	file, errs := Parse("DROP DATABASE IF EXISTS mydb FORCE")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropDatabaseStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if !stmt.Force {
		t.Error("Force should be true")
	}
}

// ---------------------------------------------------------------------------
// DROP SCHEMA (alias for DROP DATABASE)
// ---------------------------------------------------------------------------

func TestDropSchema_Basic(t *testing.T) {
	file, errs := Parse("DROP SCHEMA myschema")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropDatabaseStmt)
	if !ok {
		t.Fatalf("expected *ast.DropDatabaseStmt, got %T", file.Stmts[0])
	}
	if stmt.Name.String() != "myschema" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "myschema")
	}
}

func TestDropSchema_IfExists(t *testing.T) {
	file, errs := Parse("DROP SCHEMA IF EXISTS myschema")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropDatabaseStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
}

// ---------------------------------------------------------------------------
// ALTER DATABASE
// ---------------------------------------------------------------------------

func TestAlterDatabase_Rename(t *testing.T) {
	file, errs := Parse("ALTER DATABASE mydb RENAME newdb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterDatabaseStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterDatabaseStmt, got %T", file.Stmts[0])
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
	if stmt.NewName == nil {
		t.Fatal("NewName is nil")
	}
	if stmt.NewName.String() != "newdb" {
		t.Errorf("NewName = %q, want %q", stmt.NewName.String(), "newdb")
	}
}

func TestAlterDatabase_RenameWithTo(t *testing.T) {
	file, errs := Parse("ALTER DATABASE mydb RENAME TO newdb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AlterDatabaseStmt)
	if stmt.NewName.String() != "newdb" {
		t.Errorf("NewName = %q, want %q", stmt.NewName.String(), "newdb")
	}
}

func TestAlterDatabase_SetProperties(t *testing.T) {
	file, errs := Parse(`ALTER DATABASE mydb SET PROPERTIES("key"="value")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AlterDatabaseStmt)
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "key" {
		t.Errorf("Property key = %q, want %q", stmt.Properties[0].Key, "key")
	}
	if stmt.Properties[0].Value != "value" {
		t.Errorf("Property value = %q, want %q", stmt.Properties[0].Value, "value")
	}
}

func TestAlterDatabase_SetQuota(t *testing.T) {
	// SET QUOTA is parsed in a best-effort fashion.
	file, errs := Parse("ALTER DATABASE mydb SET QUOTA 10737418240")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterDatabaseStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterDatabaseStmt, got %T", file.Stmts[0])
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
}

// ---------------------------------------------------------------------------
// NodeTag
// ---------------------------------------------------------------------------

func TestCreateDatabaseStmt_Tag(t *testing.T) {
	file, errs := Parse("CREATE DATABASE mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateDatabaseStmt {
		t.Errorf("Tag() = %v, want T_CreateDatabaseStmt", file.Stmts[0].Tag())
	}
}

func TestDropDatabaseStmt_Tag(t *testing.T) {
	file, errs := Parse("DROP DATABASE mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_DropDatabaseStmt {
		t.Errorf("Tag() = %v, want T_DropDatabaseStmt", file.Stmts[0].Tag())
	}
}

func TestAlterDatabaseStmt_Tag(t *testing.T) {
	file, errs := Parse("ALTER DATABASE mydb RENAME newdb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_AlterDatabaseStmt {
		t.Errorf("Tag() = %v, want T_AlterDatabaseStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// Multiple statements
// ---------------------------------------------------------------------------

func TestCreateDropDatabase_MultiStatement(t *testing.T) {
	input := "CREATE DATABASE db1; DROP DATABASE db1"
	file, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 2 {
		t.Fatalf("expected 2 stmts, got %d", len(file.Stmts))
	}
	if _, ok := file.Stmts[0].(*ast.CreateDatabaseStmt); !ok {
		t.Errorf("Stmts[0]: expected *ast.CreateDatabaseStmt, got %T", file.Stmts[0])
	}
	if _, ok := file.Stmts[1].(*ast.DropDatabaseStmt); !ok {
		t.Errorf("Stmts[1]: expected *ast.DropDatabaseStmt, got %T", file.Stmts[1])
	}
}

// ---------------------------------------------------------------------------
// Other CREATE/ALTER/DROP still unsupported
// ---------------------------------------------------------------------------

func TestCreateTable_NowSupported(t *testing.T) {
	file, errs := Parse("CREATE TABLE t (id INT)")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	_, ok := file.Stmts[0].(*ast.CreateTableStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateTableStmt, got %T", file.Stmts[0])
	}
}

func TestAlterTable_NowSupported(t *testing.T) {
	// ALTER TABLE is now implemented (T2.2). Verify it parses without errors.
	file, errs := Parse("ALTER TABLE t ADD COLUMN c INT")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	if _, ok := file.Stmts[0].(*ast.AlterTableStmt); !ok {
		t.Fatalf("expected *ast.AlterTableStmt, got %T", file.Stmts[0])
	}
}

func TestDropTable_StillUnsupported(t *testing.T) {
	_, errs := Parse("DROP TABLE t")
	if len(errs) == 0 {
		t.Fatal("expected error for unsupported DROP TABLE")
	}
	want := "DROP statement parsing is not yet supported"
	if errs[0].Msg != want {
		t.Errorf("error = %q, want %q", errs[0].Msg, want)
	}
}
