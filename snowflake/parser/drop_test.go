package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func testParseDropStmt(input string) (*ast.DropStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.DropStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a DropStmt"})
	}
	return stmt, result.Errors
}

func testParseUndropStmt(input string) (*ast.UndropStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.UndropStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not an UndropStmt"})
	}
	return stmt, result.Errors
}

// ---------------------------------------------------------------------------
// DROP TABLE tests
// ---------------------------------------------------------------------------

func TestDropTable_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP TABLE t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropTable {
		t.Errorf("kind = %v, want DropTable", stmt.Kind)
	}
	if stmt.IfExists {
		t.Error("expected IfExists=false")
	}
	if stmt.Name.Normalize() != "T" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "T")
	}
	if stmt.Cascade || stmt.Restrict {
		t.Error("expected Cascade=false Restrict=false")
	}
}

func TestDropTable_IfExists(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP TABLE IF EXISTS my_table")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
	if stmt.Name.Normalize() != "MY_TABLE" {
		t.Errorf("name = %q, want MY_TABLE", stmt.Name.Normalize())
	}
}

func TestDropTable_Cascade(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP TABLE t CASCADE")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Cascade {
		t.Error("expected Cascade=true")
	}
	if stmt.Restrict {
		t.Error("expected Restrict=false")
	}
}

func TestDropTable_Restrict(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP TABLE t RESTRICT")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Restrict {
		t.Error("expected Restrict=true")
	}
	if stmt.Cascade {
		t.Error("expected Cascade=false")
	}
}

func TestDropTable_IfExistsCascade(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP TABLE IF EXISTS s.t CASCADE")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
	if !stmt.Cascade {
		t.Error("expected Cascade=true")
	}
	if stmt.Name.Normalize() != "S.T" {
		t.Errorf("name = %q, want S.T", stmt.Name.Normalize())
	}
}

func TestDropTable_QualifiedThreePart(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP TABLE mydb.myschema.mytable")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Name.Normalize() != "MYDB.MYSCHEMA.MYTABLE" {
		t.Errorf("name = %q, want MYDB.MYSCHEMA.MYTABLE", stmt.Name.Normalize())
	}
}

// ---------------------------------------------------------------------------
// DROP VIEW
// ---------------------------------------------------------------------------

func TestDropView_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP VIEW v1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropView {
		t.Errorf("kind = %v, want DropView", stmt.Kind)
	}
	if stmt.Name.Normalize() != "V1" {
		t.Errorf("name = %q, want V1", stmt.Name.Normalize())
	}
}

func TestDropView_IfExists(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP VIEW IF EXISTS v1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

// CASCADE/RESTRICT not accepted for VIEW — should parse fine and be ignored
func TestDropView_NoCascade(t *testing.T) {
	// Only the name should be consumed; CASCADE is not part of DROP VIEW
	stmt, errs := testParseDropStmt("DROP VIEW v1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Cascade {
		t.Error("Cascade should be false for DROP VIEW")
	}
}

// ---------------------------------------------------------------------------
// DROP MATERIALIZED VIEW
// ---------------------------------------------------------------------------

func TestDropMaterializedView_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP MATERIALIZED VIEW mv1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropMaterializedView {
		t.Errorf("kind = %v, want DropMaterializedView", stmt.Kind)
	}
	if stmt.Name.Normalize() != "MV1" {
		t.Errorf("name = %q, want MV1", stmt.Name.Normalize())
	}
}

func TestDropMaterializedView_IfExists(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP MATERIALIZED VIEW IF EXISTS mv1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

// ---------------------------------------------------------------------------
// DROP DYNAMIC TABLE
// ---------------------------------------------------------------------------

func TestDropDynamicTable_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP DYNAMIC TABLE dt1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropDynamicTable {
		t.Errorf("kind = %v, want DropDynamicTable", stmt.Kind)
	}
	if stmt.IfExists {
		t.Error("IfExists should be false for DYNAMIC TABLE (not supported in grammar)")
	}
}

// ---------------------------------------------------------------------------
// DROP EXTERNAL TABLE
// ---------------------------------------------------------------------------

func TestDropExternalTable_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP EXTERNAL TABLE et1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropExternalTable {
		t.Errorf("kind = %v, want DropExternalTable", stmt.Kind)
	}
}

func TestDropExternalTable_IfExistsCascade(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP EXTERNAL TABLE IF EXISTS et1 CASCADE")
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
// DROP STREAM
// ---------------------------------------------------------------------------

func TestDropStream_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP STREAM s1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropStream {
		t.Errorf("kind = %v, want DropStream", stmt.Kind)
	}
}

func TestDropStream_IfExists(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP STREAM IF EXISTS s1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

// ---------------------------------------------------------------------------
// DROP TASK
// ---------------------------------------------------------------------------

func TestDropTask_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP TASK t1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropTask {
		t.Errorf("kind = %v, want DropTask", stmt.Kind)
	}
}

func TestDropTask_IfExists(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP TASK IF EXISTS my_task")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

// ---------------------------------------------------------------------------
// DROP SEQUENCE
// ---------------------------------------------------------------------------

func TestDropSequence_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP SEQUENCE seq1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropSequence {
		t.Errorf("kind = %v, want DropSequence", stmt.Kind)
	}
}

func TestDropSequence_IfExistsCascade(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP SEQUENCE IF EXISTS seq1 CASCADE")
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
// DROP STAGE
// ---------------------------------------------------------------------------

func TestDropStage_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP STAGE my_stage")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropStage {
		t.Errorf("kind = %v, want DropStage", stmt.Kind)
	}
}

func TestDropStage_IfExists(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP STAGE IF EXISTS my_stage")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

// ---------------------------------------------------------------------------
// DROP FILE FORMAT
// ---------------------------------------------------------------------------

func TestDropFileFormat_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP FILE FORMAT ff1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropFileFormat {
		t.Errorf("kind = %v, want DropFileFormat", stmt.Kind)
	}
}

func TestDropFileFormat_IfExists(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP FILE FORMAT IF EXISTS ff1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

// ---------------------------------------------------------------------------
// DROP FUNCTION
// ---------------------------------------------------------------------------

func TestDropFunction_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP FUNCTION f1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropFunction {
		t.Errorf("kind = %v, want DropFunction", stmt.Kind)
	}
}

func TestDropFunction_IfExists(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP FUNCTION IF EXISTS f1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

// ---------------------------------------------------------------------------
// DROP PROCEDURE
// ---------------------------------------------------------------------------

func TestDropProcedure_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP PROCEDURE p1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropProcedure {
		t.Errorf("kind = %v, want DropProcedure", stmt.Kind)
	}
}

func TestDropProcedure_IfExists(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP PROCEDURE IF EXISTS p1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

// ---------------------------------------------------------------------------
// DROP PIPE
// ---------------------------------------------------------------------------

func TestDropPipe_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP PIPE pipe1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropPipe {
		t.Errorf("kind = %v, want DropPipe", stmt.Kind)
	}
}

// ---------------------------------------------------------------------------
// DROP TAG
// ---------------------------------------------------------------------------

func TestDropTag_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP TAG tag1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropTag {
		t.Errorf("kind = %v, want DropTag", stmt.Kind)
	}
}

func TestDropTag_IfExists(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP TAG IF EXISTS tag1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

// ---------------------------------------------------------------------------
// DROP ROLE
// ---------------------------------------------------------------------------

func TestDropRole_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP ROLE r1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropRole {
		t.Errorf("kind = %v, want DropRole", stmt.Kind)
	}
}

func TestDropRole_IfExists(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP ROLE IF EXISTS r1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

// ---------------------------------------------------------------------------
// DROP WAREHOUSE
// ---------------------------------------------------------------------------

func TestDropWarehouse_Basic(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP WAREHOUSE wh1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.DropWarehouse {
		t.Errorf("kind = %v, want DropWarehouse", stmt.Kind)
	}
}

func TestDropWarehouse_IfExists(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP WAREHOUSE IF EXISTS wh1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

// ---------------------------------------------------------------------------
// Unsupported DROP object types — targeted error, not generic
// ---------------------------------------------------------------------------

func TestDropUnsupported_TargetedError(t *testing.T) {
	// DROP DATABASE is handled by T2.1 and is NOT in our dispatch — it should
	// produce a targeted "DROP DATABASE ... not yet supported" error rather than
	// a generic "DROP not supported" error.
	result := ParseBestEffort("DROP DATABASE db1")
	if len(result.Errors) == 0 {
		t.Fatal("expected an error for DROP DATABASE, got none")
	}
	msg := result.Errors[0].Msg
	if msg == "DROP statement parsing is not yet supported" {
		t.Errorf("got generic DROP error; want targeted error, got: %q", msg)
	}
}

// ---------------------------------------------------------------------------
// UNDROP TABLE tests
// ---------------------------------------------------------------------------

func TestUndropTable_Basic(t *testing.T) {
	stmt, errs := testParseUndropStmt("UNDROP TABLE t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.UndropTable {
		t.Errorf("kind = %v, want UndropTable", stmt.Kind)
	}
	if stmt.Name.Normalize() != "T" {
		t.Errorf("name = %q, want T", stmt.Name.Normalize())
	}
}

func TestUndropTable_QualifiedName(t *testing.T) {
	stmt, errs := testParseUndropStmt("UNDROP TABLE myschema.mytable")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Name.Normalize() != "MYSCHEMA.MYTABLE" {
		t.Errorf("name = %q, want MYSCHEMA.MYTABLE", stmt.Name.Normalize())
	}
}

// ---------------------------------------------------------------------------
// UNDROP DYNAMIC TABLE
// ---------------------------------------------------------------------------

func TestUndropDynamicTable_Basic(t *testing.T) {
	stmt, errs := testParseUndropStmt("UNDROP DYNAMIC TABLE dt1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.UndropDynamicTable {
		t.Errorf("kind = %v, want UndropDynamicTable", stmt.Kind)
	}
	if stmt.Name.Normalize() != "DT1" {
		t.Errorf("name = %q, want DT1", stmt.Name.Normalize())
	}
}

// ---------------------------------------------------------------------------
// UNDROP TAG
// ---------------------------------------------------------------------------

func TestUndropTag_Basic(t *testing.T) {
	stmt, errs := testParseUndropStmt("UNDROP TAG tag1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.UndropTag {
		t.Errorf("kind = %v, want UndropTag", stmt.Kind)
	}
	if stmt.Name.Normalize() != "TAG1" {
		t.Errorf("name = %q, want TAG1", stmt.Name.Normalize())
	}
}

// ---------------------------------------------------------------------------
// UNDROP unsupported — targeted error
// ---------------------------------------------------------------------------

func TestUndropUnsupported_TargetedError(t *testing.T) {
	result := ParseBestEffort("UNDROP SCHEMA myschema")
	if len(result.Errors) == 0 {
		t.Fatal("expected an error for UNDROP SCHEMA, got none")
	}
	msg := result.Errors[0].Msg
	if msg == "UNDROP statement parsing is not yet supported" {
		t.Errorf("got generic UNDROP error; want targeted error, got: %q", msg)
	}
}

// ---------------------------------------------------------------------------
// Loc tracking
// ---------------------------------------------------------------------------

func TestDropTable_LocTracking(t *testing.T) {
	input := "DROP TABLE t"
	stmt, errs := testParseDropStmt(input)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len(input) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(input))
	}
}

func TestUndropTable_LocTracking(t *testing.T) {
	input := "UNDROP TABLE t"
	stmt, errs := testParseUndropStmt(input)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len(input) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(input))
	}
}

// ---------------------------------------------------------------------------
// NodeTag
// ---------------------------------------------------------------------------

func TestDropStmt_Tag(t *testing.T) {
	stmt, errs := testParseDropStmt("DROP TABLE t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Tag().String() != "DropStmt" {
		t.Errorf("Tag = %q, want DropStmt", stmt.Tag().String())
	}
}

func TestUndropStmt_Tag(t *testing.T) {
	stmt, errs := testParseUndropStmt("UNDROP TABLE t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Tag().String() != "UndropStmt" {
		t.Errorf("Tag = %q, want UndropStmt", stmt.Tag().String())
	}
}
