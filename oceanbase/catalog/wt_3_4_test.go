package catalog

import "testing"

// --- 3.4 ALTER TABLE State — Column Operations ---

func TestWalkThrough_3_4_AddColumnEnd(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT)")
	wtExec(t, c, "ALTER TABLE t ADD COLUMN c VARCHAR(50)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	if len(tbl.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(tbl.Columns))
	}
	col := tbl.GetColumn("c")
	if col == nil {
		t.Fatal("column 'c' not found")
	}
	if col.Position != 3 {
		t.Errorf("expected column c at position 3, got %d", col.Position)
	}
}

func TestWalkThrough_3_4_AddColumnFirst(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT)")
	wtExec(t, c, "ALTER TABLE t ADD COLUMN z VARCHAR(50) FIRST")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	if len(tbl.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(tbl.Columns))
	}
	// z should be at position 1
	col := tbl.GetColumn("z")
	if col == nil {
		t.Fatal("column 'z' not found")
	}
	if col.Position != 1 {
		t.Errorf("expected column z at position 1, got %d", col.Position)
	}
	// a should shift to position 2
	if tbl.GetColumn("a").Position != 2 {
		t.Errorf("expected column a at position 2, got %d", tbl.GetColumn("a").Position)
	}
	// b should shift to position 3
	if tbl.GetColumn("b").Position != 3 {
		t.Errorf("expected column b at position 3, got %d", tbl.GetColumn("b").Position)
	}
}

func TestWalkThrough_3_4_AddColumnAfter(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT)")
	wtExec(t, c, "ALTER TABLE t ADD COLUMN x VARCHAR(50) AFTER a")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	if len(tbl.Columns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(tbl.Columns))
	}
	// Order should be: a(1), x(2), b(3), c(4)
	expected := []struct {
		name string
		pos  int
	}{{"a", 1}, {"x", 2}, {"b", 3}, {"c", 4}}
	for _, e := range expected {
		col := tbl.GetColumn(e.name)
		if col == nil {
			t.Fatalf("column %q not found", e.name)
		}
		if col.Position != e.pos {
			t.Errorf("column %q: expected position %d, got %d", e.name, e.pos, col.Position)
		}
	}
}

func TestWalkThrough_3_4_DropColumn(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT)")
	wtExec(t, c, "ALTER TABLE t DROP COLUMN b")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	if len(tbl.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(tbl.Columns))
	}
	if tbl.GetColumn("b") != nil {
		t.Error("column 'b' should have been dropped")
	}
	// Positions should be resequenced: a=1, c=2
	if tbl.GetColumn("a").Position != 1 {
		t.Errorf("expected column a at position 1, got %d", tbl.GetColumn("a").Position)
	}
	if tbl.GetColumn("c").Position != 2 {
		t.Errorf("expected column c at position 2, got %d", tbl.GetColumn("c").Position)
	}
}

func TestWalkThrough_3_4_ModifyColumnType(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b VARCHAR(50), c INT)")
	wtExec(t, c, "ALTER TABLE t MODIFY COLUMN b TEXT")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	col := tbl.GetColumn("b")
	if col == nil {
		t.Fatal("column 'b' not found")
	}
	if col.DataType != "text" {
		t.Errorf("expected DataType 'text', got %q", col.DataType)
	}
	// Position should remain unchanged
	if col.Position != 2 {
		t.Errorf("expected position 2, got %d", col.Position)
	}
}

func TestWalkThrough_3_4_ModifyColumnNullability(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT NOT NULL)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("b")
	if col.Nullable {
		t.Error("column 'b' should initially be NOT NULL")
	}
	wtExec(t, c, "ALTER TABLE t MODIFY COLUMN b INT NULL")
	tbl = c.GetDatabase("testdb").GetTable("t")
	col = tbl.GetColumn("b")
	if !col.Nullable {
		t.Error("column 'b' should be nullable after MODIFY")
	}
}

func TestWalkThrough_3_4_ChangeColumnName(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, old_name VARCHAR(50), c INT)")
	wtExec(t, c, "ALTER TABLE t CHANGE COLUMN old_name new_name VARCHAR(50)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	// Old name should be gone
	if tbl.GetColumn("old_name") != nil {
		t.Error("old column name 'old_name' should not exist")
	}
	// New name should be present
	col := tbl.GetColumn("new_name")
	if col == nil {
		t.Fatal("column 'new_name' not found")
	}
}

func TestWalkThrough_3_4_ChangeColumnTypeAndAttrs(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b VARCHAR(50) NOT NULL)")
	wtExec(t, c, "ALTER TABLE t CHANGE COLUMN b b_new VARCHAR(100) NULL DEFAULT 'hello'")
	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("b_new")
	if col == nil {
		t.Fatal("column 'b_new' not found")
	}
	if col.DataType != "varchar" {
		t.Errorf("expected DataType 'varchar', got %q", col.DataType)
	}
	if col.ColumnType != "varchar(100)" {
		t.Errorf("expected ColumnType 'varchar(100)', got %q", col.ColumnType)
	}
	if !col.Nullable {
		t.Error("column should be nullable after CHANGE")
	}
	if col.Default == nil || *col.Default != "'hello'" {
		def := "<nil>"
		if col.Default != nil {
			def = *col.Default
		}
		t.Errorf("expected default 'hello', got %s", def)
	}
}

func TestWalkThrough_3_4_RenameColumn(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT, c INT)")
	wtExec(t, c, "ALTER TABLE t RENAME COLUMN b TO b_renamed")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl.GetColumn("b") != nil {
		t.Error("old column name 'b' should not exist")
	}
	col := tbl.GetColumn("b_renamed")
	if col == nil {
		t.Fatal("column 'b_renamed' not found")
	}
	// Position should be unchanged (2)
	if col.Position != 2 {
		t.Errorf("expected position 2, got %d", col.Position)
	}
}

func TestWalkThrough_3_4_AlterColumnSetDefault(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT)")
	wtExec(t, c, "ALTER TABLE t ALTER COLUMN b SET DEFAULT 42")
	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("b")
	if col.Default == nil {
		t.Fatal("expected default value, got nil")
	}
	if *col.Default != "42" {
		t.Errorf("expected default '42', got %q", *col.Default)
	}
}

func TestWalkThrough_3_4_AlterColumnDropDefault(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT DEFAULT 10)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("b")
	if col.Default == nil {
		t.Fatal("expected initial default value")
	}
	wtExec(t, c, "ALTER TABLE t ALTER COLUMN b DROP DEFAULT")
	tbl = c.GetDatabase("testdb").GetTable("t")
	col = tbl.GetColumn("b")
	if col.Default != nil {
		t.Errorf("expected nil default after DROP DEFAULT, got %q", *col.Default)
	}
	if !col.DefaultDropped {
		t.Error("expected DefaultDropped to be true")
	}
}

func TestWalkThrough_3_4_AlterColumnVisibility(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	col := tbl.GetColumn("b")
	if col.Invisible {
		t.Error("column should start visible")
	}
	wtExec(t, c, "ALTER TABLE t ALTER COLUMN b SET INVISIBLE")
	tbl = c.GetDatabase("testdb").GetTable("t")
	col = tbl.GetColumn("b")
	if !col.Invisible {
		t.Error("column should be invisible after SET INVISIBLE")
	}
	wtExec(t, c, "ALTER TABLE t ALTER COLUMN b SET VISIBLE")
	tbl = c.GetDatabase("testdb").GetTable("t")
	col = tbl.GetColumn("b")
	if col.Invisible {
		t.Error("column should be visible after SET VISIBLE")
	}
}
