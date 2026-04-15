package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// mustParseMerge parses input and returns the first statement as *ast.MergeStmt.
func mustParseMerge(t *testing.T, input string) *ast.MergeStmt {
	t.Helper()
	file, errs := Parse(input)
	if len(errs) > 0 {
		t.Fatalf("Parse(%q) errors: %v", input, errs)
	}
	if len(file.Stmts) == 0 {
		t.Fatalf("Parse(%q) returned no statements", input)
	}
	stmt, ok := file.Stmts[0].(*ast.MergeStmt)
	if !ok {
		t.Fatalf("Parse(%q) got %T, want *ast.MergeStmt", input, file.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// Basic MERGE: WHEN MATCHED THEN DELETE
// ---------------------------------------------------------------------------

func TestMergeBasicDelete(t *testing.T) {
	stmt := mustParseMerge(t,
		"MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE",
	)

	// Target
	if stmt.Target == nil || len(stmt.Target.Parts) != 1 || stmt.Target.Parts[0] != "t" {
		t.Errorf("Target = %v, want [t]", stmt.Target)
	}

	// Source
	ref, ok := stmt.Source.(*ast.TableRef)
	if !ok {
		t.Fatalf("Source = %T, want *ast.TableRef", stmt.Source)
	}
	if len(ref.Name.Parts) != 1 || ref.Name.Parts[0] != "s" {
		t.Errorf("Source.Name = %v, want [s]", ref.Name)
	}

	// ON condition
	if stmt.On == nil {
		t.Fatal("On condition is nil")
	}

	// One WHEN clause
	if len(stmt.Clauses) != 1 {
		t.Fatalf("Clauses count = %d, want 1", len(stmt.Clauses))
	}
	clause := stmt.Clauses[0]
	if clause.NotMatched {
		t.Error("clause.NotMatched = true, want false")
	}
	if clause.Action != ast.MergeActionDelete {
		t.Errorf("clause.Action = %d, want MergeActionDelete", clause.Action)
	}
}

// ---------------------------------------------------------------------------
// WHEN MATCHED THEN UPDATE SET col = expr
// ---------------------------------------------------------------------------

func TestMergeUpdateSet(t *testing.T) {
	stmt := mustParseMerge(t,
		"MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET t.c = s.c",
	)

	if len(stmt.Clauses) != 1 {
		t.Fatalf("Clauses count = %d, want 1", len(stmt.Clauses))
	}
	clause := stmt.Clauses[0]
	if clause.Action != ast.MergeActionUpdate {
		t.Errorf("clause.Action = %d, want MergeActionUpdate", clause.Action)
	}
	if clause.UpdateAll {
		t.Error("clause.UpdateAll = true, want false")
	}
	if len(clause.Assignments) != 1 {
		t.Fatalf("Assignments count = %d, want 1", len(clause.Assignments))
	}
	// Column name should be the last part (strip table qualifier)
	gotCol := ""
	if col := clause.Assignments[0].Column; col != nil && len(col.Parts) > 0 {
		gotCol = col.Parts[len(col.Parts)-1]
	}
	if gotCol != "c" {
		t.Errorf("Assignment.Column = %q, want %q", gotCol, "c")
	}
}

// ---------------------------------------------------------------------------
// WHEN NOT MATCHED THEN INSERT (cols) VALUES (exprs)
// ---------------------------------------------------------------------------

func TestMergeInsertWithColumns(t *testing.T) {
	stmt := mustParseMerge(t,
		"MERGE INTO t USING s ON t.id = s.id WHEN NOT MATCHED THEN INSERT (id, c) VALUES (s.id, s.c)",
	)

	if len(stmt.Clauses) != 1 {
		t.Fatalf("Clauses count = %d, want 1", len(stmt.Clauses))
	}
	clause := stmt.Clauses[0]
	if !clause.NotMatched {
		t.Error("clause.NotMatched = false, want true")
	}
	if clause.Action != ast.MergeActionInsert {
		t.Errorf("clause.Action = %d, want MergeActionInsert", clause.Action)
	}
	if len(clause.Columns) != 2 {
		t.Fatalf("Columns count = %d, want 2", len(clause.Columns))
	}
	if clause.Columns[0] != "id" || clause.Columns[1] != "c" {
		t.Errorf("Columns = %v, want [id c]", clause.Columns)
	}
	if len(clause.Values) != 2 {
		t.Fatalf("Values count = %d, want 2", len(clause.Values))
	}
}

// ---------------------------------------------------------------------------
// Full combo: UPDATE + INSERT
// ---------------------------------------------------------------------------

func TestMergeFullCombo(t *testing.T) {
	stmt := mustParseMerge(t,
		`MERGE INTO t USING s ON t.id = s.id
		 WHEN MATCHED THEN UPDATE SET t.c = s.c
		 WHEN NOT MATCHED THEN INSERT *`,
	)

	if len(stmt.Clauses) != 2 {
		t.Fatalf("Clauses count = %d, want 2", len(stmt.Clauses))
	}

	// First clause: WHEN MATCHED THEN UPDATE SET
	c0 := stmt.Clauses[0]
	if c0.NotMatched {
		t.Error("Clauses[0].NotMatched = true, want false")
	}
	if c0.Action != ast.MergeActionUpdate {
		t.Errorf("Clauses[0].Action = %d, want MergeActionUpdate", c0.Action)
	}

	// Second clause: WHEN NOT MATCHED THEN INSERT *
	c1 := stmt.Clauses[1]
	if !c1.NotMatched {
		t.Error("Clauses[1].NotMatched = false, want true")
	}
	if c1.Action != ast.MergeActionInsert {
		t.Errorf("Clauses[1].Action = %d, want MergeActionInsert", c1.Action)
	}
	if !c1.InsertAll {
		t.Error("Clauses[1].InsertAll = false, want true")
	}
}

// ---------------------------------------------------------------------------
// UPDATE SET * (wildcard update)
// ---------------------------------------------------------------------------

func TestMergeUpdateSetStar(t *testing.T) {
	stmt := mustParseMerge(t,
		"MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET *",
	)

	if len(stmt.Clauses) != 1 {
		t.Fatalf("Clauses count = %d, want 1", len(stmt.Clauses))
	}
	clause := stmt.Clauses[0]
	if clause.Action != ast.MergeActionUpdate {
		t.Errorf("clause.Action = %d, want MergeActionUpdate", clause.Action)
	}
	if !clause.UpdateAll {
		t.Error("clause.UpdateAll = false, want true")
	}
	if len(clause.Assignments) != 0 {
		t.Errorf("Assignments count = %d, want 0", len(clause.Assignments))
	}
}

// ---------------------------------------------------------------------------
// AND condition in WHEN clause
// ---------------------------------------------------------------------------

func TestMergeAndCondition(t *testing.T) {
	stmt := mustParseMerge(t,
		"MERGE INTO t USING s ON t.id = s.id WHEN MATCHED AND t.x > 10 THEN DELETE",
	)

	if len(stmt.Clauses) != 1 {
		t.Fatalf("Clauses count = %d, want 1", len(stmt.Clauses))
	}
	clause := stmt.Clauses[0]
	if clause.And == nil {
		t.Fatal("clause.And is nil, want non-nil condition")
	}
	if clause.Action != ast.MergeActionDelete {
		t.Errorf("clause.Action = %d, want MergeActionDelete", clause.Action)
	}
}

// ---------------------------------------------------------------------------
// Target and source aliases
// ---------------------------------------------------------------------------

func TestMergeWithAliases(t *testing.T) {
	stmt := mustParseMerge(t,
		"MERGE INTO target AS tgt USING source AS src ON tgt.id = src.id WHEN MATCHED THEN DELETE",
	)

	if stmt.TargetAlias != "tgt" {
		t.Errorf("TargetAlias = %q, want %q", stmt.TargetAlias, "tgt")
	}
	if stmt.SourceAlias != "src" {
		t.Errorf("SourceAlias = %q, want %q", stmt.SourceAlias, "src")
	}
}

// ---------------------------------------------------------------------------
// Multiple assignments in UPDATE SET
// ---------------------------------------------------------------------------

func TestMergeUpdateMultipleAssignments(t *testing.T) {
	stmt := mustParseMerge(t,
		"MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET t.a = s.a, t.b = s.b, t.c = s.c",
	)

	if len(stmt.Clauses) != 1 {
		t.Fatalf("Clauses count = %d, want 1", len(stmt.Clauses))
	}
	clause := stmt.Clauses[0]
	if len(clause.Assignments) != 3 {
		t.Fatalf("Assignments count = %d, want 3", len(clause.Assignments))
	}
	want := []string{"a", "b", "c"}
	for i, a := range clause.Assignments {
		// Column is now a qualified ObjectName; check the last part
		gotCol := ""
		if a.Column != nil && len(a.Column.Parts) > 0 {
			gotCol = a.Column.Parts[len(a.Column.Parts)-1]
		}
		if gotCol != want[i] {
			t.Errorf("Assignments[%d].Column = %q, want %q", i, gotCol, want[i])
		}
	}
}

// ---------------------------------------------------------------------------
// NodeTag
// ---------------------------------------------------------------------------

func TestMergeNodeTags(t *testing.T) {
	stmt := mustParseMerge(t,
		"MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE",
	)
	if stmt.Tag() != ast.T_MergeStmt {
		t.Errorf("stmt.Tag() = %v, want T_MergeStmt", stmt.Tag())
	}
	if len(stmt.Clauses) != 1 {
		t.Fatalf("Clauses count = %d, want 1", len(stmt.Clauses))
	}
	if stmt.Clauses[0].Tag() != ast.T_MergeClause {
		t.Errorf("clause.Tag() = %v, want T_MergeClause", stmt.Clauses[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// Subquery as source
// ---------------------------------------------------------------------------

func TestMergeSubquerySource(t *testing.T) {
	stmt := mustParseMerge(t,
		"MERGE INTO t USING (SELECT id, c FROM src WHERE active = 1) sub ON t.id = sub.id WHEN MATCHED THEN DELETE",
	)
	// Source should be a SubqueryExpr (placeholder)
	if stmt.Source == nil {
		t.Fatal("Source is nil")
	}
	if stmt.SourceAlias != "sub" {
		t.Errorf("SourceAlias = %q, want %q", stmt.SourceAlias, "sub")
	}
}

// ---------------------------------------------------------------------------
// Loc tracking
// ---------------------------------------------------------------------------

func TestMergeLocTracking(t *testing.T) {
	input := "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE"
	stmt := mustParseMerge(t, input)
	if !stmt.Loc.IsValid() {
		t.Error("stmt.Loc is not valid")
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("stmt.Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len(input) {
		t.Errorf("stmt.Loc.End = %d, want %d", stmt.Loc.End, len(input))
	}
}
