package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func testParseCreateTable(input string) (*ast.CreateTableStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.CreateTableStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a CreateTableStmt"})
	}
	return stmt, result.Errors
}

// ---------------------------------------------------------------------------
// Task 3: Basic forms
// ---------------------------------------------------------------------------

func TestCreateTable_Basic(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt == nil {
		t.Fatal("expected CreateTableStmt, got nil")
	}
	if stmt.Name.Normalize() != "T" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "T")
	}
	if len(stmt.Columns) != 1 {
		t.Fatalf("columns = %d, want 1", len(stmt.Columns))
	}
	col := stmt.Columns[0]
	if col.Name.Normalize() != "ID" {
		t.Errorf("col name = %q, want %q", col.Name.Normalize(), "ID")
	}
	if col.DataType == nil {
		t.Fatal("col DataType is nil")
	}
	if col.DataType.Kind != ast.TypeInt {
		t.Errorf("col type = %v, want TypeInt", col.DataType.Kind)
	}
}

func TestCreateTable_MultipleColumns(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT, name VARCHAR(100), active BOOLEAN)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 3 {
		t.Fatalf("columns = %d, want 3", len(stmt.Columns))
	}
	if stmt.Columns[0].DataType.Kind != ast.TypeInt {
		t.Errorf("col[0] type = %v, want TypeInt", stmt.Columns[0].DataType.Kind)
	}
	if stmt.Columns[1].DataType.Kind != ast.TypeVarchar {
		t.Errorf("col[1] type = %v, want TypeVarchar", stmt.Columns[1].DataType.Kind)
	}
	if len(stmt.Columns[1].DataType.Params) != 1 || stmt.Columns[1].DataType.Params[0] != 100 {
		t.Errorf("col[1] varchar params = %v, want [100]", stmt.Columns[1].DataType.Params)
	}
	if stmt.Columns[2].DataType.Kind != ast.TypeBoolean {
		t.Errorf("col[2] type = %v, want TypeBoolean", stmt.Columns[2].DataType.Kind)
	}
}

func TestCreateTable_OrReplace(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE OR REPLACE TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace=true")
	}
}

func TestCreateTable_Transient(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TRANSIENT TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Transient {
		t.Error("expected Transient=true")
	}
}

func TestCreateTable_Temporary(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TEMPORARY TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Temporary {
		t.Error("expected Temporary=true")
	}
}

func TestCreateTable_Volatile(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE VOLATILE TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Volatile {
		t.Error("expected Volatile=true")
	}
}

func TestCreateTable_LocalTemporary(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE LOCAL TEMPORARY TABLE t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Temporary {
		t.Error("expected Temporary=true")
	}
}

func TestCreateTable_IfNotExists(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE IF NOT EXISTS t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists=true")
	}
}

func TestCreateTable_AllModifiers(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE OR REPLACE TRANSIENT TABLE IF NOT EXISTS mydb.myschema.t (id INT)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace=true")
	}
	if !stmt.Transient {
		t.Error("expected Transient=true")
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists=true")
	}
	// 3-part name
	if stmt.Name.Normalize() != "MYDB.MYSCHEMA.T" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "MYDB.MYSCHEMA.T")
	}
	if stmt.Name.Database.Normalize() != "MYDB" {
		t.Errorf("database = %q, want %q", stmt.Name.Database.Normalize(), "MYDB")
	}
	if stmt.Name.Schema.Normalize() != "MYSCHEMA" {
		t.Errorf("schema = %q, want %q", stmt.Name.Schema.Normalize(), "MYSCHEMA")
	}
}

func TestCreateTable_Like(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t2 LIKE db.schema.t1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Like == nil {
		t.Fatal("expected Like != nil")
	}
	if stmt.Like.Normalize() != "DB.SCHEMA.T1" {
		t.Errorf("like source = %q, want %q", stmt.Like.Normalize(), "DB.SCHEMA.T1")
	}
	if len(stmt.Columns) != 0 {
		t.Errorf("columns = %d, want 0", len(stmt.Columns))
	}
}

func TestCreateTable_Clone(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t2 CLONE t1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Clone == nil {
		t.Fatal("expected Clone != nil")
	}
	if stmt.Clone.Source.Normalize() != "T1" {
		t.Errorf("clone source = %q, want %q", stmt.Clone.Source.Normalize(), "T1")
	}
	if stmt.Clone.AtBefore != "" {
		t.Errorf("AtBefore = %q, want empty", stmt.Clone.AtBefore)
	}
}

func TestCreateTable_CloneAtTimestamp(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t2 CLONE t1 AT (TIMESTAMP => '2024-01-01')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Clone == nil {
		t.Fatal("expected Clone != nil")
	}
	if stmt.Clone.AtBefore != "AT" {
		t.Errorf("AtBefore = %q, want %q", stmt.Clone.AtBefore, "AT")
	}
	if stmt.Clone.Kind != "TIMESTAMP" {
		t.Errorf("Kind = %q, want %q", stmt.Clone.Kind, "TIMESTAMP")
	}
	if stmt.Clone.Value != "2024-01-01" {
		t.Errorf("Value = %q, want %q", stmt.Clone.Value, "2024-01-01")
	}
}

func TestCreateTable_CloneBeforeStatement(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t2 CLONE t1 BEFORE (STATEMENT => '8e5d0ca1-e866-44fa-843b-5e6ad35e8bb7')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Clone == nil {
		t.Fatal("expected Clone != nil")
	}
	if stmt.Clone.AtBefore != "BEFORE" {
		t.Errorf("AtBefore = %q, want %q", stmt.Clone.AtBefore, "BEFORE")
	}
	if stmt.Clone.Kind != "STATEMENT" {
		t.Errorf("Kind = %q, want %q", stmt.Clone.Kind, "STATEMENT")
	}
}

func TestCreateTable_CTAS(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t AS SELECT 1 AS id, 'hello' AS name")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.AsSelect == nil {
		t.Fatal("expected AsSelect != nil")
	}
	sel, ok := stmt.AsSelect.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("AsSelect is %T, want *ast.SelectStmt", stmt.AsSelect)
	}
	if len(sel.Targets) != 2 {
		t.Errorf("select targets = %d, want 2", len(sel.Targets))
	}
}

func TestCreateTable_CTASWithColumns(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT, name VARCHAR) AS SELECT 1, 'hello'")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 2 {
		t.Errorf("columns = %d, want 2", len(stmt.Columns))
	}
	if stmt.AsSelect == nil {
		t.Fatal("expected AsSelect != nil")
	}
	sel, ok := stmt.AsSelect.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("AsSelect is %T, want *ast.SelectStmt", stmt.AsSelect)
	}
	if len(sel.Targets) != 2 {
		t.Errorf("select targets = %d, want 2", len(sel.Targets))
	}
}

// ---------------------------------------------------------------------------
// Task 4: Column features + constraints
// ---------------------------------------------------------------------------

func TestCreateTable_ColumnNotNull(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT NOT NULL, name VARCHAR NULL)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("columns = %d, want 2", len(stmt.Columns))
	}
	if !stmt.Columns[0].NotNull {
		t.Error("col[0]: expected NotNull=true")
	}
	if !stmt.Columns[1].Nullable {
		t.Error("col[1]: expected Nullable=true")
	}
}

func TestCreateTable_ColumnDefault(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT DEFAULT 0, name VARCHAR DEFAULT 'unknown')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("columns = %d, want 2", len(stmt.Columns))
	}
	lit0, ok := stmt.Columns[0].Default.(*ast.Literal)
	if !ok {
		t.Fatalf("col[0] default: expected *ast.Literal, got %T", stmt.Columns[0].Default)
	}
	if lit0.Kind != ast.LitInt {
		t.Errorf("col[0] default kind = %v, want LitInt", lit0.Kind)
	}
	lit1, ok := stmt.Columns[1].Default.(*ast.Literal)
	if !ok {
		t.Fatalf("col[1] default: expected *ast.Literal, got %T", stmt.Columns[1].Default)
	}
	if lit1.Kind != ast.LitString {
		t.Errorf("col[1] default kind = %v, want LitString", lit1.Kind)
	}
}

func TestCreateTable_Identity(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT IDENTITY(1, 1))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	col := stmt.Columns[0]
	if col.Identity == nil {
		t.Fatal("expected Identity != nil")
	}
	if col.Identity.Start == nil || *col.Identity.Start != 1 {
		t.Errorf("Identity.Start = %v, want 1", col.Identity.Start)
	}
	if col.Identity.Increment == nil || *col.Identity.Increment != 1 {
		t.Errorf("Identity.Increment = %v, want 1", col.Identity.Increment)
	}
}

func TestCreateTable_AutoincrementStartIncrement(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT AUTOINCREMENT START 100 INCREMENT 10)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	col := stmt.Columns[0]
	if col.Identity == nil {
		t.Fatal("expected Identity != nil")
	}
	if col.Identity.Start == nil || *col.Identity.Start != 100 {
		t.Errorf("Identity.Start = %v, want 100", col.Identity.Start)
	}
	if col.Identity.Increment == nil || *col.Identity.Increment != 10 {
		t.Errorf("Identity.Increment = %v, want 10", col.Identity.Increment)
	}
}

func TestCreateTable_IdentityOrder(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT IDENTITY ORDER)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	col := stmt.Columns[0]
	if col.Identity == nil {
		t.Fatal("expected Identity != nil")
	}
	if col.Identity.Order == nil || !*col.Identity.Order {
		t.Errorf("Identity.Order = %v, want true", col.Identity.Order)
	}
}

func TestCreateTable_Collate(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (name VARCHAR COLLATE 'en-ci')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	col := stmt.Columns[0]
	if col.Collate != "en-ci" {
		t.Errorf("Collate = %q, want %q", col.Collate, "en-ci")
	}
}

func TestCreateTable_ColumnComment(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT COMMENT 'primary identifier')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	col := stmt.Columns[0]
	if col.Comment == nil {
		t.Fatal("expected Comment != nil")
	}
	if *col.Comment != "primary identifier" {
		t.Errorf("Comment = %q, want %q", *col.Comment, "primary identifier")
	}
}

func TestCreateTable_InlinePrimaryKey(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT PRIMARY KEY)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	col := stmt.Columns[0]
	if col.InlineConstraint == nil {
		t.Fatal("expected InlineConstraint != nil")
	}
	if col.InlineConstraint.Type != ast.ConstrPrimaryKey {
		t.Errorf("constraint type = %v, want ConstrPrimaryKey", col.InlineConstraint.Type)
	}
}

func TestCreateTable_InlineUnique(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (email VARCHAR UNIQUE)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	col := stmt.Columns[0]
	if col.InlineConstraint == nil {
		t.Fatal("expected InlineConstraint != nil")
	}
	if col.InlineConstraint.Type != ast.ConstrUnique {
		t.Errorf("constraint type = %v, want ConstrUnique", col.InlineConstraint.Type)
	}
}

func TestCreateTable_InlineForeignKey(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (customer_id INT FOREIGN KEY REFERENCES customers (id))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	col := stmt.Columns[0]
	if col.InlineConstraint == nil {
		t.Fatal("expected InlineConstraint != nil")
	}
	if col.InlineConstraint.Type != ast.ConstrForeignKey {
		t.Errorf("constraint type = %v, want ConstrForeignKey", col.InlineConstraint.Type)
	}
	ref := col.InlineConstraint.References
	if ref == nil {
		t.Fatal("expected References != nil")
	}
	if ref.Table.Normalize() != "CUSTOMERS" {
		t.Errorf("references table = %q, want %q", ref.Table.Normalize(), "CUSTOMERS")
	}
	if len(ref.Columns) != 1 || ref.Columns[0].Normalize() != "ID" {
		t.Errorf("references columns = %v, want [ID]", ref.Columns)
	}
}

func TestCreateTable_NamedInlineConstraint(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT CONSTRAINT pk_t PRIMARY KEY)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	col := stmt.Columns[0]
	if col.InlineConstraint == nil {
		t.Fatal("expected InlineConstraint != nil")
	}
	if col.InlineConstraint.Name.Normalize() != "PK_T" {
		t.Errorf("constraint name = %q, want %q", col.InlineConstraint.Name.Normalize(), "PK_T")
	}
	if col.InlineConstraint.Type != ast.ConstrPrimaryKey {
		t.Errorf("constraint type = %v, want ConstrPrimaryKey", col.InlineConstraint.Type)
	}
}

func TestCreateTable_OutOfLinePK(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT, name VARCHAR, PRIMARY KEY (id))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 2 {
		t.Errorf("columns = %d, want 2", len(stmt.Columns))
	}
	if len(stmt.Constraints) != 1 {
		t.Fatalf("constraints = %d, want 1", len(stmt.Constraints))
	}
	con := stmt.Constraints[0]
	if con.Type != ast.ConstrPrimaryKey {
		t.Errorf("constraint type = %v, want ConstrPrimaryKey", con.Type)
	}
	if len(con.Columns) != 1 || con.Columns[0].Normalize() != "ID" {
		t.Errorf("constraint columns = %v, want [ID]", con.Columns)
	}
}

func TestCreateTable_OutOfLineCompositePK(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Constraints) != 1 {
		t.Fatalf("constraints = %d, want 1", len(stmt.Constraints))
	}
	con := stmt.Constraints[0]
	if len(con.Columns) != 2 {
		t.Errorf("PK columns = %d, want 2", len(con.Columns))
	}
	if con.Columns[0].Normalize() != "A" {
		t.Errorf("PK col[0] = %q, want %q", con.Columns[0].Normalize(), "A")
	}
	if con.Columns[1].Normalize() != "B" {
		t.Errorf("PK col[1] = %q, want %q", con.Columns[1].Normalize(), "B")
	}
}

func TestCreateTable_OutOfLineFK(t *testing.T) {
	stmt, errs := testParseCreateTable(`CREATE TABLE t (
		id INT,
		customer_id INT,
		FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE CASCADE ON UPDATE SET NULL
	)`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Constraints) != 1 {
		t.Fatalf("constraints = %d, want 1", len(stmt.Constraints))
	}
	con := stmt.Constraints[0]
	if con.Type != ast.ConstrForeignKey {
		t.Errorf("constraint type = %v, want ConstrForeignKey", con.Type)
	}
	if con.References == nil {
		t.Fatal("expected References != nil")
	}
	if con.References.OnDelete != ast.RefActCascade {
		t.Errorf("OnDelete = %v, want RefActCascade", con.References.OnDelete)
	}
	if con.References.OnUpdate != ast.RefActSetNull {
		t.Errorf("OnUpdate = %v, want RefActSetNull", con.References.OnUpdate)
	}
}

func TestCreateTable_OutOfLineUnique(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (a INT, b INT, UNIQUE (a, b))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Constraints) != 1 {
		t.Fatalf("constraints = %d, want 1", len(stmt.Constraints))
	}
	con := stmt.Constraints[0]
	if con.Type != ast.ConstrUnique {
		t.Errorf("constraint type = %v, want ConstrUnique", con.Type)
	}
	if len(con.Columns) != 2 {
		t.Errorf("unique columns = %d, want 2", len(con.Columns))
	}
}

func TestCreateTable_NamedConstraint(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT, CONSTRAINT pk_t PRIMARY KEY (id))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Constraints) != 1 {
		t.Fatalf("constraints = %d, want 1", len(stmt.Constraints))
	}
	con := stmt.Constraints[0]
	if con.Name.Normalize() != "PK_T" {
		t.Errorf("constraint name = %q, want %q", con.Name.Normalize(), "PK_T")
	}
	if con.Type != ast.ConstrPrimaryKey {
		t.Errorf("constraint type = %v, want ConstrPrimaryKey", con.Type)
	}
}

func TestCreateTable_MixedColumnsAndConstraints(t *testing.T) {
	stmt, errs := testParseCreateTable(`CREATE TABLE t (
		id INT NOT NULL,
		email VARCHAR NOT NULL,
		customer_id INT,
		PRIMARY KEY (id),
		UNIQUE (email),
		FOREIGN KEY (customer_id) REFERENCES customers (id)
	)`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 3 {
		t.Errorf("columns = %d, want 3", len(stmt.Columns))
	}
	if len(stmt.Constraints) != 3 {
		t.Fatalf("constraints = %d, want 3", len(stmt.Constraints))
	}
	if stmt.Constraints[0].Type != ast.ConstrPrimaryKey {
		t.Errorf("constraint[0] type = %v, want ConstrPrimaryKey", stmt.Constraints[0].Type)
	}
	if stmt.Constraints[1].Type != ast.ConstrUnique {
		t.Errorf("constraint[1] type = %v, want ConstrUnique", stmt.Constraints[1].Type)
	}
	if stmt.Constraints[2].Type != ast.ConstrForeignKey {
		t.Errorf("constraint[2] type = %v, want ConstrForeignKey", stmt.Constraints[2].Type)
	}
}

// ---------------------------------------------------------------------------
// Task 5: Table properties + acceptance
// ---------------------------------------------------------------------------

func TestCreateTable_ClusterBy(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT) CLUSTER BY (id)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.ClusterBy) != 1 {
		t.Errorf("ClusterBy length = %d, want 1", len(stmt.ClusterBy))
	}
}

func TestCreateTable_ClusterByLinear(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (a INT, b INT) CLUSTER BY LINEAR (a, b)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Linear {
		t.Error("expected Linear=true")
	}
	if len(stmt.ClusterBy) != 2 {
		t.Errorf("ClusterBy length = %d, want 2", len(stmt.ClusterBy))
	}
}

func TestCreateTable_TableComment(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT) COMMENT = 'my table'")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Comment == nil {
		t.Fatal("expected Comment != nil")
	}
	if *stmt.Comment != "my table" {
		t.Errorf("Comment = %q, want %q", *stmt.Comment, "my table")
	}
}

func TestCreateTable_CopyGrants(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT) COPY GRANTS")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.CopyGrants {
		t.Error("expected CopyGrants=true")
	}
}

func TestCreateTable_WithTag(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT) WITH TAG (env = 'prod', team = 'data')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Tags) != 2 {
		t.Fatalf("tags = %d, want 2", len(stmt.Tags))
	}
	if stmt.Tags[0].Name.Normalize() != "ENV" {
		t.Errorf("tag[0] name = %q, want %q", stmt.Tags[0].Name.Normalize(), "ENV")
	}
	if stmt.Tags[0].Value != "prod" {
		t.Errorf("tag[0] value = %q, want %q", stmt.Tags[0].Value, "prod")
	}
	if stmt.Tags[1].Name.Normalize() != "TEAM" {
		t.Errorf("tag[1] name = %q, want %q", stmt.Tags[1].Name.Normalize(), "TEAM")
	}
	if stmt.Tags[1].Value != "data" {
		t.Errorf("tag[1] value = %q, want %q", stmt.Tags[1].Value, "data")
	}
}

func TestCreateTable_DataRetention(t *testing.T) {
	_, errs := testParseCreateTable("CREATE TABLE t (id INT) DATA_RETENTION_TIME_IN_DAYS = 90")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestCreateTable_ChangeTracking(t *testing.T) {
	_, errs := testParseCreateTable("CREATE TABLE t (id INT) CHANGE_TRACKING = TRUE")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestCreateTable_MaskingPolicy(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (ssn VARCHAR WITH MASKING POLICY ssn_mask)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	col := stmt.Columns[0]
	if col.MaskingPolicy == nil {
		t.Fatal("expected MaskingPolicy != nil")
	}
	if col.MaskingPolicy.Normalize() != "SSN_MASK" {
		t.Errorf("MaskingPolicy name = %q, want %q", col.MaskingPolicy.Normalize(), "SSN_MASK")
	}
}

func TestCreateTable_ColumnWithTag(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (id INT WITH TAG (sensitivity = 'PII', team = 'security'))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	col := stmt.Columns[0]
	if len(col.Tags) != 2 {
		t.Fatalf("col tags = %d, want 2", len(col.Tags))
	}
	if col.Tags[0].Name.Normalize() != "SENSITIVITY" || col.Tags[0].Value != "PII" {
		t.Errorf("tag[0] = %v/%v, want SENSITIVITY/PII", col.Tags[0].Name.Normalize(), col.Tags[0].Value)
	}
	if col.Tags[1].Name.Normalize() != "TEAM" || col.Tags[1].Value != "security" {
		t.Errorf("tag[1] = %v/%v, want TEAM/security", col.Tags[1].Name.Normalize(), col.Tags[1].Value)
	}
}

func TestCreateTable_VirtualColumn(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t (a INT, b INT, c INT AS (a + b))")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 3 {
		t.Fatalf("columns = %d, want 3", len(stmt.Columns))
	}
	col := stmt.Columns[2]
	if col.VirtualExpr == nil {
		t.Fatal("expected VirtualExpr != nil")
	}
	_, ok := col.VirtualExpr.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("VirtualExpr is %T, want *ast.BinaryExpr", col.VirtualExpr)
	}
}

func TestCreateTable_Complex(t *testing.T) {
	input := `CREATE OR REPLACE TABLE mydb.myschema.orders (
    id INT IDENTITY(1, 1) NOT NULL,
    customer_id INT NOT NULL,
    amount NUMBER(10, 2) DEFAULT 0,
    status VARCHAR(50) COLLATE 'en-ci',
    created_at TIMESTAMP_NTZ DEFAULT CURRENT_TIMESTAMP(),
    CONSTRAINT pk_orders PRIMARY KEY (id),
    CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE CASCADE,
    UNIQUE (customer_id, created_at)
)
CLUSTER BY (customer_id)
COMMENT = 'Order records'
COPY GRANTS`

	stmt, errs := testParseCreateTable(input)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	// Table name
	if stmt.Name.Normalize() != "MYDB.MYSCHEMA.ORDERS" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "MYDB.MYSCHEMA.ORDERS")
	}

	// Modifiers
	if !stmt.OrReplace {
		t.Error("expected OrReplace=true")
	}
	if !stmt.CopyGrants {
		t.Error("expected CopyGrants=true")
	}

	// Comment
	if stmt.Comment == nil || *stmt.Comment != "Order records" {
		t.Errorf("Comment = %v, want 'Order records'", stmt.Comment)
	}

	// ClusterBy
	if len(stmt.ClusterBy) != 1 {
		t.Errorf("ClusterBy length = %d, want 1", len(stmt.ClusterBy))
	}

	// Columns
	if len(stmt.Columns) != 5 {
		t.Fatalf("columns = %d, want 5", len(stmt.Columns))
	}

	// col[0]: id INT IDENTITY(1,1) NOT NULL
	col0 := stmt.Columns[0]
	if col0.Identity == nil {
		t.Error("col[0]: expected Identity != nil")
	} else {
		if col0.Identity.Start == nil || *col0.Identity.Start != 1 {
			t.Errorf("col[0] Identity.Start = %v, want 1", col0.Identity.Start)
		}
		if col0.Identity.Increment == nil || *col0.Identity.Increment != 1 {
			t.Errorf("col[0] Identity.Increment = %v, want 1", col0.Identity.Increment)
		}
	}
	if !col0.NotNull {
		t.Error("col[0]: expected NotNull=true")
	}

	// col[1]: customer_id INT NOT NULL
	if !stmt.Columns[1].NotNull {
		t.Error("col[1]: expected NotNull=true")
	}

	// Constraints
	if len(stmt.Constraints) != 3 {
		t.Fatalf("constraints = %d, want 3", len(stmt.Constraints))
	}

	// constraint[0]: CONSTRAINT pk_orders PRIMARY KEY (id)
	con0 := stmt.Constraints[0]
	if con0.Name.Normalize() != "PK_ORDERS" {
		t.Errorf("constraint[0] name = %q, want %q", con0.Name.Normalize(), "PK_ORDERS")
	}
	if con0.Type != ast.ConstrPrimaryKey {
		t.Errorf("constraint[0] type = %v, want ConstrPrimaryKey", con0.Type)
	}

	// constraint[1]: CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE CASCADE
	con1 := stmt.Constraints[1]
	if con1.Name.Normalize() != "FK_CUSTOMER" {
		t.Errorf("constraint[1] name = %q, want %q", con1.Name.Normalize(), "FK_CUSTOMER")
	}
	if con1.Type != ast.ConstrForeignKey {
		t.Errorf("constraint[1] type = %v, want ConstrForeignKey", con1.Type)
	}
	if con1.References == nil || con1.References.OnDelete != ast.RefActCascade {
		t.Errorf("constraint[1] OnDelete = %v, want RefActCascade", con1.References)
	}
}
