package parser

import "testing"

// StarRocks INSERT divergences (PR3): OVERWRITE no-TABLE form, BY NAME,
// and the FILES(propertyList) target.

func TestInsertOverwriteNoTable(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT OVERWRITE t SELECT a FROM src")
	if !stmt.Overwrite {
		t.Error("Overwrite = false, want true")
	}
	if stmt.Target == nil || stmt.Target.Parts[len(stmt.Target.Parts)-1] != "t" {
		t.Errorf("target = %+v, want t", stmt.Target)
	}
}

// Regression: the doris OVERWRITE TABLE form must still parse (additive).
func TestInsertOverwriteTableStillParses(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT OVERWRITE TABLE test VALUES (1, 2)")
	if !stmt.Overwrite {
		t.Error("Overwrite = false, want true")
	}
}

func TestInsertByName(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO t BY NAME SELECT a, b FROM s")
	if !stmt.ByName {
		t.Error("ByName = false, want true")
	}
}

// Modifiers may appear in any order (grammar: insertLabelOrColumnAliases*).
func TestInsertByNameThenWithLabel(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO t BY NAME WITH LABEL lbl SELECT a, b FROM s")
	if !stmt.ByName || stmt.Label != "lbl" {
		t.Errorf("ByName=%v Label=%q, want true/lbl", stmt.ByName, stmt.Label)
	}
}

// BY NAME and a column list are mutually exclusive — must reject.
func TestInsertByNameWithColumnListRejected(t *testing.T) {
	_, errs := Parse("INSERT INTO t BY NAME (a, b) SELECT a, b FROM s")
	if len(errs) == 0 {
		t.Fatal("expected a parse error for BY NAME + column list, got none")
	}
}

// WITH LABEL / BY NAME apply to FILES targets too.
func TestInsertFilesWithLabel(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO FILES('path'='s3://x') WITH LABEL lbl SELECT a FROM s")
	if len(stmt.FileTarget) != 1 || stmt.Label != "lbl" {
		t.Errorf("FileTarget=%d Label=%q, want 1/lbl", len(stmt.FileTarget), stmt.Label)
	}
}

func TestInsertIntoFiles(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO FILES('path'='s3://x', 'format'='parquet') SELECT a FROM s")
	if len(stmt.FileTarget) != 2 {
		t.Fatalf("FileTarget = %d props, want 2", len(stmt.FileTarget))
	}
	if stmt.Target != nil {
		t.Errorf("Target = %+v, want nil (FILES target)", stmt.Target)
	}
}
