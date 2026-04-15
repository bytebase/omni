package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// testParseMergeStmt parses and returns the first statement as *ast.MergeStmt.
func testParseMergeStmt(t *testing.T, input string) *ast.MergeStmt {
	t.Helper()
	result := ParseBestEffort(input)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if len(result.File.Stmts) == 0 {
		t.Fatal("expected statement, got none")
	}
	stmt, ok := result.File.Stmts[0].(*ast.MergeStmt)
	if !ok {
		t.Fatalf("expected *ast.MergeStmt, got %T", result.File.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// 1. Basic MERGE with WHEN MATCHED UPDATE
// ---------------------------------------------------------------------------

func TestMerge_BasicMatchedUpdate(t *testing.T) {
	stmt := testParseMergeStmt(t, `
		MERGE INTO target t
		USING source s
		ON t.id = s.id
		WHEN MATCHED THEN UPDATE SET t.value = s.value
	`)
	if stmt.Target.Name.Name != "target" {
		t.Errorf("target = %q, want %q", stmt.Target.Name.Name, "target")
	}
	if stmt.TargetAlias.Name != "t" {
		t.Errorf("targetAlias = %q, want %q", stmt.TargetAlias.Name, "t")
	}
	if stmt.SourceAlias.Name != "s" {
		t.Errorf("sourceAlias = %q, want %q", stmt.SourceAlias.Name, "s")
	}
	if stmt.On == nil {
		t.Error("On should not be nil")
	}
	if len(stmt.Whens) != 1 {
		t.Fatalf("whens = %d, want 1", len(stmt.Whens))
	}
	when := stmt.Whens[0]
	if !when.Matched {
		t.Error("Matched should be true")
	}
	if when.Action != ast.MergeActionUpdate {
		t.Errorf("action = %v, want MergeActionUpdate", when.Action)
	}
	if len(when.Sets) != 1 {
		t.Fatalf("sets = %d, want 1", len(when.Sets))
	}
}

// ---------------------------------------------------------------------------
// 2. MERGE with WHEN MATCHED DELETE
// ---------------------------------------------------------------------------

func TestMerge_MatchedDelete(t *testing.T) {
	stmt := testParseMergeStmt(t, `
		MERGE INTO target
		USING source
		ON target.id = source.id
		WHEN MATCHED THEN DELETE
	`)
	if len(stmt.Whens) != 1 {
		t.Fatalf("whens = %d, want 1", len(stmt.Whens))
	}
	when := stmt.Whens[0]
	if !when.Matched {
		t.Error("Matched should be true")
	}
	if when.Action != ast.MergeActionDelete {
		t.Errorf("action = %v, want MergeActionDelete", when.Action)
	}
}

// ---------------------------------------------------------------------------
// 3. MERGE with WHEN NOT MATCHED INSERT
// ---------------------------------------------------------------------------

func TestMerge_NotMatchedInsert(t *testing.T) {
	stmt := testParseMergeStmt(t, `
		MERGE INTO target t
		USING source s
		ON t.id = s.id
		WHEN NOT MATCHED THEN INSERT (id, name) VALUES (s.id, s.name)
	`)
	if len(stmt.Whens) != 1 {
		t.Fatalf("whens = %d, want 1", len(stmt.Whens))
	}
	when := stmt.Whens[0]
	if when.Matched {
		t.Error("Matched should be false for NOT MATCHED")
	}
	if when.Action != ast.MergeActionInsert {
		t.Errorf("action = %v, want MergeActionInsert", when.Action)
	}
	if len(when.InsertCols) != 2 {
		t.Errorf("insertCols = %d, want 2", len(when.InsertCols))
	}
	if len(when.InsertVals) != 2 {
		t.Errorf("insertVals = %d, want 2", len(when.InsertVals))
	}
	if when.InsertDefault {
		t.Error("InsertDefault should be false")
	}
}

// ---------------------------------------------------------------------------
// 4. MERGE with multiple WHEN clauses
// ---------------------------------------------------------------------------

func TestMerge_MultipleWhenClauses(t *testing.T) {
	stmt := testParseMergeStmt(t, `
		MERGE INTO target t
		USING source s
		ON t.id = s.id
		WHEN MATCHED AND s.active = false THEN DELETE
		WHEN MATCHED THEN UPDATE SET t.value = s.value, t.ts = CURRENT_TIMESTAMP()
		WHEN NOT MATCHED THEN INSERT VALUES (s.id, s.value, s.ts)
	`)
	if len(stmt.Whens) != 3 {
		t.Fatalf("whens = %d, want 3", len(stmt.Whens))
	}

	// First: WHEN MATCHED AND ... THEN DELETE
	w0 := stmt.Whens[0]
	if !w0.Matched {
		t.Error("whens[0].Matched should be true")
	}
	if w0.AndCond == nil {
		t.Error("whens[0].AndCond should not be nil")
	}
	if w0.Action != ast.MergeActionDelete {
		t.Errorf("whens[0].action = %v, want Delete", w0.Action)
	}

	// Second: WHEN MATCHED THEN UPDATE
	w1 := stmt.Whens[1]
	if !w1.Matched {
		t.Error("whens[1].Matched should be true")
	}
	if w1.AndCond != nil {
		t.Error("whens[1].AndCond should be nil")
	}
	if w1.Action != ast.MergeActionUpdate {
		t.Errorf("whens[1].action = %v, want Update", w1.Action)
	}
	if len(w1.Sets) != 2 {
		t.Errorf("whens[1].sets = %d, want 2", len(w1.Sets))
	}

	// Third: WHEN NOT MATCHED THEN INSERT
	w2 := stmt.Whens[2]
	if w2.Matched {
		t.Error("whens[2].Matched should be false")
	}
	if w2.Action != ast.MergeActionInsert {
		t.Errorf("whens[2].action = %v, want Insert", w2.Action)
	}
	if len(w2.InsertCols) != 0 {
		t.Errorf("whens[2].insertCols = %d, want 0 (no column list)", len(w2.InsertCols))
	}
	if len(w2.InsertVals) != 3 {
		t.Errorf("whens[2].insertVals = %d, want 3", len(w2.InsertVals))
	}
}

// ---------------------------------------------------------------------------
// 5. MERGE with WHEN NOT MATCHED BY SOURCE
// ---------------------------------------------------------------------------

func TestMerge_NotMatchedBySource(t *testing.T) {
	stmt := testParseMergeStmt(t, `
		MERGE INTO target t
		USING source s
		ON t.id = s.id
		WHEN NOT MATCHED BY SOURCE THEN DELETE
	`)
	if len(stmt.Whens) != 1 {
		t.Fatalf("whens = %d, want 1", len(stmt.Whens))
	}
	when := stmt.Whens[0]
	if when.Matched {
		t.Error("Matched should be false")
	}
	if !when.BySource {
		t.Error("BySource should be true")
	}
	if when.ByTarget {
		t.Error("ByTarget should be false")
	}
	if when.Action != ast.MergeActionDelete {
		t.Errorf("action = %v, want Delete", when.Action)
	}
}

// ---------------------------------------------------------------------------
// 6. MERGE with WHEN NOT MATCHED BY TARGET
// ---------------------------------------------------------------------------

func TestMerge_NotMatchedByTarget(t *testing.T) {
	stmt := testParseMergeStmt(t, `
		MERGE INTO target t
		USING source s
		ON t.id = s.id
		WHEN NOT MATCHED BY TARGET THEN INSERT VALUES (s.id, s.name)
	`)
	if len(stmt.Whens) != 1 {
		t.Fatalf("whens = %d, want 1", len(stmt.Whens))
	}
	when := stmt.Whens[0]
	if when.Matched {
		t.Error("Matched should be false")
	}
	if when.BySource {
		t.Error("BySource should be false")
	}
	if !when.ByTarget {
		t.Error("ByTarget should be true")
	}
}

// ---------------------------------------------------------------------------
// 7. MERGE with subquery as source
// ---------------------------------------------------------------------------

func TestMerge_SubquerySource(t *testing.T) {
	stmt := testParseMergeStmt(t, `
		MERGE INTO target t
		USING (SELECT id, value FROM staging WHERE processed = false) s
		ON t.id = s.id
		WHEN MATCHED THEN UPDATE SET t.value = s.value
	`)
	sourceRef, ok := stmt.Source.(*ast.TableRef)
	if !ok {
		t.Fatalf("source is %T, want *ast.TableRef", stmt.Source)
	}
	if sourceRef.Subquery == nil {
		t.Error("source.Subquery should not be nil")
	}
}

// ---------------------------------------------------------------------------
// 8. MERGE with INSERT VALUES DEFAULT
// ---------------------------------------------------------------------------

func TestMerge_InsertValuesDefault(t *testing.T) {
	stmt := testParseMergeStmt(t, `
		MERGE INTO target t
		USING source s
		ON t.id = s.id
		WHEN NOT MATCHED THEN INSERT VALUES DEFAULT
	`)
	when := stmt.Whens[0]
	if !when.InsertDefault {
		t.Error("InsertDefault should be true")
	}
	if len(when.InsertVals) != 0 {
		t.Errorf("InsertVals should be empty, got %d", len(when.InsertVals))
	}
}

// ---------------------------------------------------------------------------
// 9. MERGE without aliases
// ---------------------------------------------------------------------------

func TestMerge_NoAliases(t *testing.T) {
	stmt := testParseMergeStmt(t, `
		MERGE INTO sales
		USING new_sales
		ON sales.id = new_sales.id
		WHEN MATCHED THEN UPDATE SET sales.amount = new_sales.amount
		WHEN NOT MATCHED THEN INSERT (id, amount) VALUES (new_sales.id, new_sales.amount)
	`)
	if !stmt.TargetAlias.IsEmpty() {
		t.Errorf("TargetAlias should be empty, got %q", stmt.TargetAlias.Name)
	}
	if !stmt.SourceAlias.IsEmpty() {
		t.Errorf("SourceAlias should be empty, got %q", stmt.SourceAlias.Name)
	}
	if len(stmt.Whens) != 2 {
		t.Fatalf("whens = %d, want 2", len(stmt.Whens))
	}
}
