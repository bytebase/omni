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

func TestInsertIntoFiles(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO FILES('path'='s3://x', 'format'='parquet') SELECT a FROM s")
	if len(stmt.FileTarget) != 2 {
		t.Fatalf("FileTarget = %d props, want 2", len(stmt.FileTarget))
	}
	if stmt.Target != nil {
		t.Errorf("Target = %+v, want nil (FILES target)", stmt.Target)
	}
}
