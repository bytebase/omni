package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// mustParseUpdate parses input and returns the first statement as *ast.UpdateStmt.
func mustParseUpdate(t *testing.T, input string) *ast.UpdateStmt {
	t.Helper()
	file, errs := Parse(input)
	if len(errs) > 0 {
		t.Fatalf("Parse(%q) errors: %v", input, errs)
	}
	if len(file.Stmts) == 0 {
		t.Fatalf("Parse(%q) returned no statements", input)
	}
	stmt, ok := file.Stmts[0].(*ast.UpdateStmt)
	if !ok {
		t.Fatalf("Parse(%q) got %T, want *ast.UpdateStmt", input, file.Stmts[0])
	}
	return stmt
}

// mustParseDelete parses input and returns the first statement as *ast.DeleteStmt.
func mustParseDelete(t *testing.T, input string) *ast.DeleteStmt {
	t.Helper()
	file, errs := Parse(input)
	if len(errs) > 0 {
		t.Fatalf("Parse(%q) errors: %v", input, errs)
	}
	if len(file.Stmts) == 0 {
		t.Fatalf("Parse(%q) returned no statements", input)
	}
	stmt, ok := file.Stmts[0].(*ast.DeleteStmt)
	if !ok {
		t.Fatalf("Parse(%q) got %T, want *ast.DeleteStmt", input, file.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// UPDATE tests
// ---------------------------------------------------------------------------

func TestUpdateBasic(t *testing.T) {
	stmt := mustParseUpdate(t, "UPDATE t SET c = 1 WHERE id = 5")

	if len(stmt.Target.Parts) != 1 || stmt.Target.Parts[0] != "t" {
		t.Errorf("target = %v, want [t]", stmt.Target.Parts)
	}
	if stmt.TargetAlias != "" {
		t.Errorf("alias = %q, want empty", stmt.TargetAlias)
	}
	if len(stmt.Assignments) != 1 {
		t.Fatalf("assignments = %d, want 1", len(stmt.Assignments))
	}
	a := stmt.Assignments[0]
	if len(a.Column.Parts) != 1 || a.Column.Parts[0] != "c" {
		t.Errorf("assignment col = %v, want [c]", a.Column.Parts)
	}
	lit, ok := a.Value.(*ast.Literal)
	if !ok {
		t.Fatalf("assignment value = %T, want *ast.Literal", a.Value)
	}
	if lit.Kind != ast.LitInt || lit.Value != "1" {
		t.Errorf("assignment value = %v %q, want LitInt 1", lit.Kind, lit.Value)
	}
	if stmt.Where == nil {
		t.Fatal("WHERE should not be nil")
	}
	if len(stmt.From) != 0 {
		t.Errorf("from = %d, want 0", len(stmt.From))
	}
}

func TestUpdateMultiSet(t *testing.T) {
	stmt := mustParseUpdate(t, "UPDATE t SET c1 = 1, c2 = 'a' WHERE id = 5")

	if len(stmt.Assignments) != 2 {
		t.Fatalf("assignments = %d, want 2", len(stmt.Assignments))
	}
	if stmt.Assignments[0].Column.Parts[0] != "c1" {
		t.Errorf("assignment[0] col = %v, want c1", stmt.Assignments[0].Column.Parts)
	}
	if stmt.Assignments[1].Column.Parts[0] != "c2" {
		t.Errorf("assignment[1] col = %v, want c2", stmt.Assignments[1].Column.Parts)
	}
	lit, ok := stmt.Assignments[1].Value.(*ast.Literal)
	if !ok {
		t.Fatalf("assignment[1] value = %T, want *ast.Literal", stmt.Assignments[1].Value)
	}
	if lit.Kind != ast.LitString || lit.Value != "a" {
		t.Errorf("assignment[1] value = %v %q, want LitString a", lit.Kind, lit.Value)
	}
}

func TestUpdateWithFrom(t *testing.T) {
	stmt := mustParseUpdate(t, "UPDATE t SET c = other.c FROM other WHERE t.id = other.id")

	if stmt.Target.Parts[0] != "t" {
		t.Errorf("target = %v, want t", stmt.Target.Parts)
	}
	if len(stmt.Assignments) != 1 {
		t.Fatalf("assignments = %d, want 1", len(stmt.Assignments))
	}
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	ref, ok := stmt.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.TableRef", stmt.From[0])
	}
	if ref.Name.Parts[0] != "other" {
		t.Errorf("from[0] name = %v, want other", ref.Name.Parts)
	}
	if stmt.Where == nil {
		t.Fatal("WHERE should not be nil")
	}
}

func TestUpdateQualifiedColumn(t *testing.T) {
	// Test qualified column in SET: t1.c1 = t2.c1
	stmt := mustParseUpdate(t, "UPDATE t1 SET t1.c1 = t2.c1, t1.c3 = t2.c3 * 100 FROM t2 INNER JOIN t3 ON t2.id = t3.id WHERE t1.id = t2.id")

	if len(stmt.Assignments) != 2 {
		t.Fatalf("assignments = %d, want 2", len(stmt.Assignments))
	}
	col0 := stmt.Assignments[0].Column
	if len(col0.Parts) != 2 || col0.Parts[0] != "t1" || col0.Parts[1] != "c1" {
		t.Errorf("assignment[0] col = %v, want [t1 c1]", col0.Parts)
	}
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	// The FROM contains a JOIN
	_, isJoin := stmt.From[0].(*ast.JoinClause)
	if !isJoin {
		t.Errorf("from[0] = %T, want *ast.JoinClause", stmt.From[0])
	}
}

func TestUpdateExprAssignment(t *testing.T) {
	// v1 = v1 + 1
	stmt := mustParseUpdate(t, "UPDATE test SET v1 = v1 + 1 WHERE k1 = 1")

	if len(stmt.Assignments) != 1 {
		t.Fatalf("assignments = %d, want 1", len(stmt.Assignments))
	}
	_, isBinary := stmt.Assignments[0].Value.(*ast.BinaryExpr)
	if !isBinary {
		t.Errorf("assignment value = %T, want *ast.BinaryExpr", stmt.Assignments[0].Value)
	}
}

func TestUpdateLocIsValid(t *testing.T) {
	input := "UPDATE t SET c = 1 WHERE id = 5"
	stmt := mustParseUpdate(t, input)
	if !stmt.Loc.IsValid() {
		t.Error("UpdateStmt.Loc should be valid")
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len(input) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(input))
	}
}

// ---------------------------------------------------------------------------
// DELETE tests
// ---------------------------------------------------------------------------

func TestDeleteBasic(t *testing.T) {
	stmt := mustParseDelete(t, "DELETE FROM t WHERE id = 5")

	if len(stmt.Target.Parts) != 1 || stmt.Target.Parts[0] != "t" {
		t.Errorf("target = %v, want [t]", stmt.Target.Parts)
	}
	if stmt.TargetAlias != "" {
		t.Errorf("alias = %q, want empty", stmt.TargetAlias)
	}
	if len(stmt.Partition) != 0 {
		t.Errorf("partition = %v, want empty", stmt.Partition)
	}
	if len(stmt.Using) != 0 {
		t.Errorf("using = %d, want 0", len(stmt.Using))
	}
	if stmt.Where == nil {
		t.Fatal("WHERE should not be nil")
	}
}

func TestDeleteWithPartitionSingle(t *testing.T) {
	// PARTITION p1 (no parentheses) — legacy corpus: DELETE FROM my_table PARTITION p1 WHERE k1 = 3
	stmt := mustParseDelete(t, "DELETE FROM my_table PARTITION p1 WHERE k1 = 3")

	if stmt.Target.Parts[0] != "my_table" {
		t.Errorf("target = %v, want my_table", stmt.Target.Parts)
	}
	if len(stmt.Partition) != 1 || stmt.Partition[0] != "p1" {
		t.Errorf("partition = %v, want [p1]", stmt.Partition)
	}
	if stmt.Where == nil {
		t.Fatal("WHERE should not be nil")
	}
}

func TestDeleteWithPartitionParenthesized(t *testing.T) {
	// PARTITIONS (p1, p2) — legacy corpus: DELETE FROM my_table PARTITIONS (p1, p2) WHERE ...
	stmt := mustParseDelete(t, "DELETE FROM my_table PARTITIONS (p1, p2) WHERE k1 >= 3 AND k2 = 'abc'")

	if len(stmt.Partition) != 2 {
		t.Fatalf("partition = %v, want [p1 p2]", stmt.Partition)
	}
	if stmt.Partition[0] != "p1" || stmt.Partition[1] != "p2" {
		t.Errorf("partition = %v, want [p1 p2]", stmt.Partition)
	}
}

func TestDeleteWithUsing(t *testing.T) {
	// Legacy corpus: DELETE FROM t1 USING t2 INNER JOIN t3 ON t2.id = t3.id WHERE t1.id = t2.id
	stmt := mustParseDelete(t, "DELETE FROM t1 USING t2 INNER JOIN t3 ON t2.id = t3.id WHERE t1.id = t2.id")

	if stmt.Target.Parts[0] != "t1" {
		t.Errorf("target = %v, want t1", stmt.Target.Parts)
	}
	if len(stmt.Using) != 1 {
		t.Fatalf("using = %d, want 1", len(stmt.Using))
	}
	// The USING clause is t2 INNER JOIN t3, which produces a JoinClause
	_, isJoin := stmt.Using[0].(*ast.JoinClause)
	if !isJoin {
		t.Errorf("using[0] = %T, want *ast.JoinClause", stmt.Using[0])
	}
	if stmt.Where == nil {
		t.Fatal("WHERE should not be nil")
	}
}

func TestDeleteLocIsValid(t *testing.T) {
	input := "DELETE FROM t WHERE id = 5"
	stmt := mustParseDelete(t, input)
	if !stmt.Loc.IsValid() {
		t.Error("DeleteStmt.Loc should be valid")
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len(input) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(input))
	}
}

func TestDeleteNoWhere(t *testing.T) {
	stmt := mustParseDelete(t, "DELETE FROM t")

	if stmt.Where != nil {
		t.Errorf("WHERE should be nil when absent")
	}
}

// ---------------------------------------------------------------------------
// Legacy corpus tests
// ---------------------------------------------------------------------------

func TestUpdateLegacyCorpus(t *testing.T) {
	cases := []string{
		"UPDATE test SET v1 = 1 WHERE k1 = 1 AND k2 = 2",
		"UPDATE test SET v1 = v1 + 1 WHERE k1 = 1",
		"UPDATE t1 SET t1.c1 = t2.c1, t1.c3 = t2.c3 * 100 FROM t2 INNER JOIN t3 ON t2.id = t3.id WHERE t1.id = t2.id",
	}
	for _, sql := range cases {
		_, errs := Parse(sql)
		if len(errs) > 0 {
			t.Errorf("Parse(%q) errors: %v", sql, errs)
		}
	}
}

func TestDeleteLegacyCorpus(t *testing.T) {
	cases := []string{
		"DELETE FROM my_table PARTITION p1 WHERE k1 = 3",
		"DELETE FROM my_table PARTITION p1 WHERE k1 >= 3 AND k2 = 'abc'",
		"DELETE FROM my_table PARTITIONS (p1, p2) WHERE k1 >= 3 AND k2 = 'abc'",
		"DELETE FROM t1 USING t2 INNER JOIN t3 ON t2.id = t3.id WHERE t1.id = t2.id",
	}
	for _, sql := range cases {
		_, errs := Parse(sql)
		if len(errs) > 0 {
			t.Errorf("Parse(%q) errors: %v", sql, errs)
		}
	}
}

// ---------------------------------------------------------------------------
// NodeTag tests
// ---------------------------------------------------------------------------

func TestUpdateDeleteNodeTags(t *testing.T) {
	updateStmt := mustParseUpdate(t, "UPDATE t SET c = 1")
	if updateStmt.Tag() != ast.T_UpdateStmt {
		t.Errorf("UpdateStmt.Tag() = %v, want T_UpdateStmt", updateStmt.Tag())
	}

	deleteStmt := mustParseDelete(t, "DELETE FROM t")
	if deleteStmt.Tag() != ast.T_DeleteStmt {
		t.Errorf("DeleteStmt.Tag() = %v, want T_DeleteStmt", deleteStmt.Tag())
	}

	if updateStmt.Assignments[0].Tag() != ast.T_Assignment {
		t.Errorf("Assignment.Tag() = %v, want T_Assignment", updateStmt.Assignments[0].Tag())
	}
}
