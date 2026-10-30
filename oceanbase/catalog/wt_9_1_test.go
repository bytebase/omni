package catalog

import (
	"strings"
	"testing"
)

// --- 5.1 Generated Column CRUD (11 scenarios) ---

func TestWalkThrough_9_1_VirtualArithmetic(t *testing.T) {
	// Scenario 1: CREATE TABLE with VIRTUAL generated column (arithmetic)
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT GENERATED ALWAYS AS (a + b) VIRTUAL)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("c")
	if col == nil {
		t.Fatal("column c not found")
	}
	if col.Generated == nil {
		t.Fatal("column c should be generated")
	}
	if col.Generated.Stored {
		t.Error("column c should be VIRTUAL, not STORED")
	}

	got := c.ShowCreateTable("testdb", "t")
	if !strings.Contains(got, "GENERATED ALWAYS AS") {
		t.Errorf("expected GENERATED ALWAYS AS in SHOW CREATE TABLE\ngot:\n%s", got)
	}
	if !strings.Contains(got, "VIRTUAL") {
		t.Errorf("expected VIRTUAL in SHOW CREATE TABLE\ngot:\n%s", got)
	}
	// Expression should contain backtick-quoted column refs and arithmetic.
	if !strings.Contains(got, "(`a` + `b`)") {
		t.Errorf("expected expression (`a` + `b`) in SHOW CREATE TABLE\ngot:\n%s", got)
	}
}

func TestWalkThrough_9_1_StoredArithmetic(t *testing.T) {
	// Scenario 2: CREATE TABLE with STORED generated column (arithmetic)
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT GENERATED ALWAYS AS (a * b) STORED)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("c")
	if col == nil {
		t.Fatal("column c not found")
	}
	if col.Generated == nil {
		t.Fatal("column c should be generated")
	}
	if !col.Generated.Stored {
		t.Error("column c should be STORED")
	}

	got := c.ShowCreateTable("testdb", "t")
	if !strings.Contains(got, "STORED") {
		t.Errorf("expected STORED in SHOW CREATE TABLE\ngot:\n%s", got)
	}
	if !strings.Contains(got, "(`a` * `b`)") {
		t.Errorf("expected expression (`a` * `b`) in SHOW CREATE TABLE\ngot:\n%s", got)
	}
}

func TestWalkThrough_9_1_VirtualConcat(t *testing.T) {
	// Scenario 3: CREATE TABLE with VIRTUAL generated column (CONCAT function)
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (first_name VARCHAR(50), last_name VARCHAR(50), full_name VARCHAR(101) GENERATED ALWAYS AS (CONCAT(first_name, ' ', last_name)) VIRTUAL)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("full_name")
	if col == nil {
		t.Fatal("column full_name not found")
	}
	if col.Generated == nil {
		t.Fatal("column full_name should be generated")
	}

	got := c.ShowCreateTable("testdb", "t")
	// MySQL renders CONCAT with charset introducer for string literals in generated columns.
	if !strings.Contains(got, "concat(") {
		t.Errorf("expected concat function in SHOW CREATE TABLE\ngot:\n%s", got)
	}
	if !strings.Contains(got, "VIRTUAL") {
		t.Errorf("expected VIRTUAL in SHOW CREATE TABLE\ngot:\n%s", got)
	}
}

func TestWalkThrough_9_1_StoredNotNull(t *testing.T) {
	// Scenario 4: CREATE TABLE with STORED generated column + NOT NULL
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT GENERATED ALWAYS AS (a + b) STORED NOT NULL)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("c")
	if col == nil {
		t.Fatal("column c not found")
	}
	if col.Generated == nil {
		t.Fatal("column c should be generated")
	}
	if !col.Generated.Stored {
		t.Error("column c should be STORED")
	}
	if col.Nullable {
		t.Error("column c should be NOT NULL")
	}

	got := c.ShowCreateTable("testdb", "t")
	// Rendering order: type, GENERATED ALWAYS AS (expr) STORED, NOT NULL
	if !strings.Contains(got, "STORED NOT NULL") {
		t.Errorf("expected 'STORED NOT NULL' in SHOW CREATE TABLE\ngot:\n%s", got)
	}
}

func TestWalkThrough_9_1_Comment(t *testing.T) {
	// Scenario 5: CREATE TABLE with generated column + COMMENT
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT GENERATED ALWAYS AS (a + b) VIRTUAL COMMENT 'sum of a and b')")

	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("c")
	if col == nil {
		t.Fatal("column c not found")
	}
	if col.Comment != "sum of a and b" {
		t.Errorf("expected comment 'sum of a and b', got %q", col.Comment)
	}

	got := c.ShowCreateTable("testdb", "t")
	// COMMENT should come after VIRTUAL
	if !strings.Contains(got, "VIRTUAL COMMENT 'sum of a and b'") {
		t.Errorf("expected 'VIRTUAL COMMENT ...' in SHOW CREATE TABLE\ngot:\n%s", got)
	}
}

func TestWalkThrough_9_1_Invisible(t *testing.T) {
	// Scenario 6: CREATE TABLE with generated column + INVISIBLE
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT GENERATED ALWAYS AS (a + b) VIRTUAL INVISIBLE)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("c")
	if col == nil {
		t.Fatal("column c not found")
	}
	if !col.Invisible {
		t.Error("column c should be INVISIBLE")
	}

	got := c.ShowCreateTable("testdb", "t")
	// INVISIBLE rendered after VIRTUAL with version comment
	if !strings.Contains(got, "INVISIBLE") {
		t.Errorf("expected INVISIBLE in SHOW CREATE TABLE\ngot:\n%s", got)
	}
}

func TestWalkThrough_9_1_JsonExtract(t *testing.T) {
	// Scenario 7: CREATE TABLE with generated column using JSON_EXTRACT
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (doc JSON, name VARCHAR(100) GENERATED ALWAYS AS (JSON_EXTRACT(doc, '$.name')) VIRTUAL)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("name")
	if col == nil {
		t.Fatal("column name not found")
	}
	if col.Generated == nil {
		t.Fatal("column name should be generated")
	}

	got := c.ShowCreateTable("testdb", "t")
	// MySQL renders json_extract in lowercase
	if !strings.Contains(got, "json_extract(") {
		t.Errorf("expected json_extract function in SHOW CREATE TABLE\ngot:\n%s", got)
	}
}

func TestWalkThrough_9_1_AlterAddGenerated(t *testing.T) {
	// Scenario 8: ALTER TABLE ADD COLUMN with GENERATED ALWAYS AS
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT)")
	wtExec(t, c, "ALTER TABLE t ADD COLUMN c INT GENERATED ALWAYS AS (a + b) VIRTUAL")

	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("c")
	if col == nil {
		t.Fatal("column c not found")
	}
	if col.Generated == nil {
		t.Fatal("column c should be generated")
	}
	if col.Generated.Stored {
		t.Error("column c should be VIRTUAL")
	}

	got := c.ShowCreateTable("testdb", "t")
	if !strings.Contains(got, "GENERATED ALWAYS AS") {
		t.Errorf("expected GENERATED ALWAYS AS in SHOW CREATE TABLE\ngot:\n%s", got)
	}
}

func TestWalkThrough_9_1_ModifyChangeExpression(t *testing.T) {
	// Scenario 9: MODIFY COLUMN to change generated expression
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT GENERATED ALWAYS AS (a + b) VIRTUAL)")
	wtExec(t, c, "ALTER TABLE t MODIFY COLUMN c INT GENERATED ALWAYS AS (a * b) VIRTUAL")

	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("c")
	if col == nil {
		t.Fatal("column c not found")
	}
	if col.Generated == nil {
		t.Fatal("column c should still be generated")
	}

	got := c.ShowCreateTable("testdb", "t")
	// Expression should now be (a * b), not (a + b)
	if !strings.Contains(got, "(`a` * `b`)") {
		t.Errorf("expected updated expression (`a` * `b`) in SHOW CREATE TABLE\ngot:\n%s", got)
	}
}

func TestWalkThrough_9_1_ModifyGeneratedToRegular(t *testing.T) {
	// Scenario 10: ALTER TABLE MODIFY generated column to regular column
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT GENERATED ALWAYS AS (a + b) VIRTUAL)")
	wtExec(t, c, "ALTER TABLE t MODIFY COLUMN c INT")

	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("c")
	if col == nil {
		t.Fatal("column c not found")
	}
	if col.Generated != nil {
		t.Error("column c should no longer be generated after MODIFY to regular column")
	}

	got := c.ShowCreateTable("testdb", "t")
	if strings.Contains(got, "GENERATED ALWAYS AS") {
		t.Errorf("expected no GENERATED ALWAYS AS after MODIFY to regular column\ngot:\n%s", got)
	}
}

func TestWalkThrough_9_1_ModifyVirtualToStored(t *testing.T) {
	// Scenario 11: ALTER TABLE MODIFY VIRTUAL to STORED — MySQL 8.0 error
	// MySQL 8.0 does not allow changing VIRTUAL to STORED in-place.
	// Error 3106 (HY000): 'Changing the STORED status' is not supported for generated columns.
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT GENERATED ALWAYS AS (a + b) VIRTUAL)")

	results, _ := c.Exec("ALTER TABLE t MODIFY COLUMN c INT GENERATED ALWAYS AS (a + b) STORED", nil)
	if len(results) == 0 {
		t.Fatal("expected result from ALTER TABLE")
	}
	if results[0].Error == nil {
		t.Fatal("expected error when changing VIRTUAL to STORED")
	}
	catErr, ok := results[0].Error.(*Error)
	if !ok {
		t.Fatalf("expected *catalog.Error, got %T", results[0].Error)
	}
	if catErr.Code != ErrUnsupportedGeneratedStorageChange {
		t.Errorf("expected error code %d, got %d: %s", ErrUnsupportedGeneratedStorageChange, catErr.Code, catErr.Message)
	}
}
