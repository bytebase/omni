package catalog

import "testing"

// --- ADD INDEX ---

func TestWalkThrough_3_5_AddIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100))")
	wtExec(t, c, "ALTER TABLE t ADD INDEX idx_name (name)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	var found *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_name" {
			found = idx
			break
		}
	}
	if found == nil {
		t.Fatal("index idx_name not found in table.Indexes")
	}
	if found.Unique {
		t.Error("regular index should not be Unique")
	}
	if found.Primary {
		t.Error("regular index should not be Primary")
	}
	if len(found.Columns) != 1 || found.Columns[0].Name != "name" {
		t.Errorf("index columns mismatch: %+v", found.Columns)
	}
}

// --- ADD UNIQUE INDEX ---

func TestWalkThrough_3_5_AddUniqueIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, email VARCHAR(200))")
	wtExec(t, c, "ALTER TABLE t ADD UNIQUE INDEX idx_email (email)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	// Check index exists and is unique.
	var uqIdx *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_email" {
			uqIdx = idx
			break
		}
	}
	if uqIdx == nil {
		t.Fatal("unique index idx_email not found")
	}
	if !uqIdx.Unique {
		t.Error("expected Unique=true on unique index")
	}
	if len(uqIdx.Columns) != 1 || uqIdx.Columns[0].Name != "email" {
		t.Errorf("unique index columns mismatch: %+v", uqIdx.Columns)
	}

	// Check constraint created.
	var uqCon *Constraint
	for _, con := range tbl.Constraints {
		if con.Name == "idx_email" && con.Type == ConUniqueKey {
			uqCon = con
			break
		}
	}
	if uqCon == nil {
		t.Fatal("unique constraint idx_email not found")
	}
	if len(uqCon.Columns) != 1 || uqCon.Columns[0] != "email" {
		t.Errorf("unique constraint columns mismatch: %v", uqCon.Columns)
	}
}

// --- ADD PRIMARY KEY ---

func TestWalkThrough_3_5_AddPrimaryKey(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, name VARCHAR(100))")
	wtExec(t, c, "ALTER TABLE t ADD PRIMARY KEY (id)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	var pkIdx *Index
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			pkIdx = idx
			break
		}
	}
	if pkIdx == nil {
		t.Fatal("no primary key index found")
	}
	if pkIdx.Name != "PRIMARY" {
		t.Errorf("expected PK index name 'PRIMARY', got %q", pkIdx.Name)
	}
	if !pkIdx.Unique {
		t.Error("PK index should be Unique=true")
	}
	if len(pkIdx.Columns) != 1 || pkIdx.Columns[0].Name != "id" {
		t.Errorf("PK index columns mismatch: %+v", pkIdx.Columns)
	}

	// PK column should be marked NOT NULL.
	col := tbl.GetColumn("id")
	if col == nil {
		t.Fatal("column id not found")
	}
	if col.Nullable {
		t.Error("PK column should be NOT NULL")
	}
}

// --- ADD FOREIGN KEY ---

func TestWalkThrough_3_5_AddForeignKey(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE parent (id INT PRIMARY KEY)")
	wtExec(t, c, "CREATE TABLE child (id INT PRIMARY KEY, parent_id INT)")
	wtExec(t, c, "ALTER TABLE child ADD CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id) ON DELETE CASCADE ON UPDATE SET NULL")

	tbl := c.GetDatabase("testdb").GetTable("child")
	if tbl == nil {
		t.Fatal("table child not found")
	}

	var fkCon *Constraint
	for _, con := range tbl.Constraints {
		if con.Name == "fk_parent" && con.Type == ConForeignKey {
			fkCon = con
			break
		}
	}
	if fkCon == nil {
		t.Fatal("FK constraint fk_parent not found")
	}
	if len(fkCon.Columns) != 1 || fkCon.Columns[0] != "parent_id" {
		t.Errorf("FK columns mismatch: %v", fkCon.Columns)
	}
	if fkCon.RefTable != "parent" {
		t.Errorf("expected RefTable 'parent', got %q", fkCon.RefTable)
	}
	if len(fkCon.RefColumns) != 1 || fkCon.RefColumns[0] != "id" {
		t.Errorf("FK RefColumns mismatch: %v", fkCon.RefColumns)
	}
	if fkCon.OnDelete != "CASCADE" {
		t.Errorf("expected OnDelete 'CASCADE', got %q", fkCon.OnDelete)
	}
	if fkCon.OnUpdate != "SET NULL" {
		t.Errorf("expected OnUpdate 'SET NULL', got %q", fkCon.OnUpdate)
	}
}

// --- ADD CHECK ---

func TestWalkThrough_3_5_AddCheck(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, age INT)")
	wtExec(t, c, "ALTER TABLE t ADD CONSTRAINT chk_age CHECK (age >= 0)")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	var chkCon *Constraint
	for _, con := range tbl.Constraints {
		if con.Name == "chk_age" && con.Type == ConCheck {
			chkCon = con
			break
		}
	}
	if chkCon == nil {
		t.Fatal("CHECK constraint chk_age not found")
	}
	if chkCon.CheckExpr == "" {
		t.Error("CHECK constraint should have a non-empty expression")
	}
	if chkCon.NotEnforced {
		t.Error("CHECK constraint should be enforced by default")
	}
}

// --- DROP INDEX ---

func TestWalkThrough_3_5_DropIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100), INDEX idx_name (name))")
	// Verify index exists first.
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	foundBefore := false
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_name" {
			foundBefore = true
			break
		}
	}
	if !foundBefore {
		t.Fatal("index idx_name should exist before drop")
	}

	wtExec(t, c, "ALTER TABLE t DROP INDEX idx_name")

	tbl = c.GetDatabase("testdb").GetTable("t")
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_name" {
			t.Error("index idx_name should have been removed after DROP INDEX")
		}
	}
}

// --- DROP PRIMARY KEY ---

func TestWalkThrough_3_5_DropPrimaryKey(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100))")
	wtExec(t, c, "ALTER TABLE t DROP PRIMARY KEY")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			t.Error("PK index should have been removed after DROP PRIMARY KEY")
		}
	}
}

// --- RENAME INDEX ---

func TestWalkThrough_3_5_RenameIndex(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100), INDEX idx_name (name))")
	wtExec(t, c, "ALTER TABLE t RENAME INDEX idx_name TO idx_name_new")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	foundOld := false
	foundNew := false
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_name" {
			foundOld = true
		}
		if idx.Name == "idx_name_new" {
			foundNew = true
		}
	}
	if foundOld {
		t.Error("old index name idx_name should no longer exist")
	}
	if !foundNew {
		t.Error("new index name idx_name_new should exist")
	}
}

// --- ALTER INDEX VISIBILITY ---

func TestWalkThrough_3_5_AlterIndexVisible(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100), INDEX idx_name (name))")

	// Default should be visible.
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	var idx *Index
	for _, i := range tbl.Indexes {
		if i.Name == "idx_name" {
			idx = i
			break
		}
	}
	if idx == nil {
		t.Fatal("index idx_name not found")
	}
	if !idx.Visible {
		t.Error("index should be visible by default")
	}

	// Set invisible.
	wtExec(t, c, "ALTER TABLE t ALTER INDEX idx_name INVISIBLE")
	tbl = c.GetDatabase("testdb").GetTable("t")
	for _, i := range tbl.Indexes {
		if i.Name == "idx_name" {
			if i.Visible {
				t.Error("index should be invisible after ALTER INDEX INVISIBLE")
			}
		}
	}

	// Set visible again.
	wtExec(t, c, "ALTER TABLE t ALTER INDEX idx_name VISIBLE")
	tbl = c.GetDatabase("testdb").GetTable("t")
	for _, i := range tbl.Indexes {
		if i.Name == "idx_name" {
			if !i.Visible {
				t.Error("index should be visible after ALTER INDEX VISIBLE")
			}
		}
	}
}

// --- ALTER CHECK ENFORCED / NOT ENFORCED ---

func TestWalkThrough_3_5_AlterCheckEnforced(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, age INT, CONSTRAINT chk_age CHECK (age >= 0))")

	// Default should be enforced.
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	var chk *Constraint
	for _, con := range tbl.Constraints {
		if con.Name == "chk_age" && con.Type == ConCheck {
			chk = con
			break
		}
	}
	if chk == nil {
		t.Fatal("CHECK constraint chk_age not found")
	}
	if chk.NotEnforced {
		t.Error("CHECK should be enforced by default")
	}

	// Set NOT ENFORCED.
	wtExec(t, c, "ALTER TABLE t ALTER CHECK chk_age NOT ENFORCED")
	tbl = c.GetDatabase("testdb").GetTable("t")
	for _, con := range tbl.Constraints {
		if con.Name == "chk_age" && con.Type == ConCheck {
			if !con.NotEnforced {
				t.Error("CHECK should be NOT ENFORCED after ALTER CHECK NOT ENFORCED")
			}
		}
	}

	// Set ENFORCED.
	wtExec(t, c, "ALTER TABLE t ALTER CHECK chk_age ENFORCED")
	tbl = c.GetDatabase("testdb").GetTable("t")
	for _, con := range tbl.Constraints {
		if con.Name == "chk_age" && con.Type == ConCheck {
			if con.NotEnforced {
				t.Error("CHECK should be ENFORCED after ALTER CHECK ENFORCED")
			}
		}
	}
}

// --- CONVERT TO CHARACTER SET ---

func TestWalkThrough_3_5_ConvertToCharset(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(100), bio TEXT, age INT) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")

	// Convert to latin1.
	wtExec(t, c, "ALTER TABLE t CONVERT TO CHARACTER SET latin1")

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	// Table charset should be updated.
	if tbl.Charset != "latin1" {
		t.Errorf("expected table charset 'latin1', got %q", tbl.Charset)
	}

	// String columns should be updated.
	nameCol := tbl.GetColumn("name")
	if nameCol == nil {
		t.Fatal("column name not found")
	}
	if nameCol.Charset != "latin1" {
		t.Errorf("expected name column charset 'latin1', got %q", nameCol.Charset)
	}

	bioCol := tbl.GetColumn("bio")
	if bioCol == nil {
		t.Fatal("column bio not found")
	}
	if bioCol.Charset != "latin1" {
		t.Errorf("expected bio column charset 'latin1', got %q", bioCol.Charset)
	}

	// Non-string column should NOT be affected.
	ageCol := tbl.GetColumn("age")
	if ageCol == nil {
		t.Fatal("column age not found")
	}
	if ageCol.Charset != "" {
		t.Errorf("expected age column charset empty, got %q", ageCol.Charset)
	}
}
