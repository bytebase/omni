package catalog

import "testing"

// --- 3.1 CREATE TABLE State ---

func TestWalkThrough_3_1_TableExists(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE users (id INT)")
	tbl := c.GetDatabase("testdb").GetTable("users")
	if tbl == nil {
		t.Fatal("table 'users' not found after CREATE TABLE")
	}
}

func TestWalkThrough_3_1_ColumnCount(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b VARCHAR(100), c TEXT)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	if len(tbl.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(tbl.Columns))
	}
}

func TestWalkThrough_3_1_ColumnNamesInOrder(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (alpha INT, beta VARCHAR(50), gamma TEXT)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	expected := []string{"alpha", "beta", "gamma"}
	for i, col := range tbl.Columns {
		if col.Name != expected[i] {
			t.Errorf("column %d: expected name %q, got %q", i, expected[i], col.Name)
		}
	}
}

func TestWalkThrough_3_1_ColumnPositionsSequential(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b VARCHAR(50), c TEXT, d BLOB)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	for i, col := range tbl.Columns {
		expected := i + 1
		if col.Position != expected {
			t.Errorf("column %q: expected position %d, got %d", col.Name, expected, col.Position)
		}
	}
}

func TestWalkThrough_3_1_ColumnTypes(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t (
		a INT,
		b VARCHAR(100),
		c DECIMAL(10,2),
		d DATETIME,
		e TEXT,
		f BLOB,
		g JSON
	)`)
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	tests := []struct {
		name     string
		dataType string
	}{
		{"a", "int"},
		{"b", "varchar"},
		{"c", "decimal"},
		{"d", "datetime"},
		{"e", "text"},
		{"f", "blob"},
		{"g", "json"},
	}
	for _, tt := range tests {
		col := tbl.GetColumn(tt.name)
		if col == nil {
			t.Errorf("column %q not found", tt.name)
			continue
		}
		if col.DataType != tt.dataType {
			t.Errorf("column %q: expected DataType %q, got %q", tt.name, tt.dataType, col.DataType)
		}
	}

	// Check ColumnType includes params.
	colB := tbl.GetColumn("b")
	if colB.ColumnType != "varchar(100)" {
		t.Errorf("column b: expected ColumnType 'varchar(100)', got %q", colB.ColumnType)
	}
	colC := tbl.GetColumn("c")
	if colC.ColumnType != "decimal(10,2)" {
		t.Errorf("column c: expected ColumnType 'decimal(10,2)', got %q", colC.ColumnType)
	}
}

func TestWalkThrough_3_1_NotNull(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT NOT NULL, b INT)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	colA := tbl.GetColumn("a")
	if colA.Nullable {
		t.Error("column a should be NOT NULL (Nullable=false)")
	}
	colB := tbl.GetColumn("b")
	if !colB.Nullable {
		t.Error("column b should be nullable by default")
	}
}

func TestWalkThrough_3_1_DefaultValue(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT DEFAULT 42, b VARCHAR(50) DEFAULT 'hello', c INT)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}

	colA := tbl.GetColumn("a")
	if colA.Default == nil {
		t.Fatal("column a: expected default value, got nil")
	}
	if *colA.Default != "42" {
		t.Errorf("column a: expected default '42', got %q", *colA.Default)
	}

	colB := tbl.GetColumn("b")
	if colB.Default == nil {
		t.Fatal("column b: expected default value, got nil")
	}
	// The default may be stored with or without quotes.
	if *colB.Default != "'hello'" && *colB.Default != "hello" {
		t.Errorf("column b: expected default 'hello' or \"'hello'\", got %q", *colB.Default)
	}

	colC := tbl.GetColumn("c")
	if colC.Default != nil {
		t.Errorf("column c: expected nil default, got %q", *colC.Default)
	}
}

func TestWalkThrough_3_1_AutoIncrement(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(50))")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	colID := tbl.GetColumn("id")
	if !colID.AutoIncrement {
		t.Error("column id should have AutoIncrement=true")
	}
	colName := tbl.GetColumn("name")
	if colName.AutoIncrement {
		t.Error("column name should not have AutoIncrement")
	}
}

func TestWalkThrough_3_1_ColumnComment(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT COMMENT 'primary identifier', name VARCHAR(50))")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	colID := tbl.GetColumn("id")
	if colID.Comment != "primary identifier" {
		t.Errorf("column id: expected comment 'primary identifier', got %q", colID.Comment)
	}
	colName := tbl.GetColumn("name")
	if colName.Comment != "" {
		t.Errorf("column name: expected empty comment, got %q", colName.Comment)
	}
}

func TestWalkThrough_3_1_TableEngine(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT) ENGINE=MyISAM")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	if tbl.Engine != "MyISAM" {
		t.Errorf("expected engine 'MyISAM', got %q", tbl.Engine)
	}
}

func TestWalkThrough_3_1_TableCharsetAndCollation(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT) DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	if tbl.Charset != "latin1" {
		t.Errorf("expected charset 'latin1', got %q", tbl.Charset)
	}
	if tbl.Collation != "latin1_swedish_ci" {
		t.Errorf("expected collation 'latin1_swedish_ci', got %q", tbl.Collation)
	}
}

func TestWalkThrough_3_1_TableComment(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT) COMMENT='user accounts'")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	if tbl.Comment != "user accounts" {
		t.Errorf("expected table comment 'user accounts', got %q", tbl.Comment)
	}
}

func TestWalkThrough_3_1_TableAutoIncrementStart(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT AUTO_INCREMENT PRIMARY KEY) AUTO_INCREMENT=1000")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	if tbl.AutoIncrement != 1000 {
		t.Errorf("expected table AUTO_INCREMENT=1000, got %d", tbl.AutoIncrement)
	}
}

func TestWalkThrough_3_1_UnsignedModifier(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT UNSIGNED, b BIGINT UNSIGNED)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	colA := tbl.GetColumn("a")
	if colA.ColumnType != "int unsigned" {
		t.Errorf("column a: expected ColumnType 'int unsigned', got %q", colA.ColumnType)
	}
	colB := tbl.GetColumn("b")
	if colB.ColumnType != "bigint unsigned" {
		t.Errorf("column b: expected ColumnType 'bigint unsigned', got %q", colB.ColumnType)
	}
}

func TestWalkThrough_3_1_GeneratedColumnVirtual(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT AS (a * 2) VIRTUAL)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	colB := tbl.GetColumn("b")
	if colB == nil {
		t.Fatal("column b not found")
	}
	if colB.Generated == nil {
		t.Fatal("column b: expected Generated info, got nil")
	}
	if colB.Generated.Stored {
		t.Error("column b: expected Stored=false for VIRTUAL")
	}
	if colB.Generated.Expr == "" {
		t.Error("column b: expected non-empty generated expression")
	}
}

func TestWalkThrough_3_1_GeneratedColumnStored(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT AS (a + 1) STORED)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	colB := tbl.GetColumn("b")
	if colB == nil {
		t.Fatal("column b not found")
	}
	if colB.Generated == nil {
		t.Fatal("column b: expected Generated info, got nil")
	}
	if !colB.Generated.Stored {
		t.Error("column b: expected Stored=true for STORED")
	}
	if colB.Generated.Expr == "" {
		t.Error("column b: expected non-empty generated expression")
	}
}

func TestWalkThrough_3_1_ColumnInvisible(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT, b INT INVISIBLE)")
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table not found")
	}
	colA := tbl.GetColumn("a")
	if colA.Invisible {
		t.Error("column a should be visible (Invisible=false)")
	}
	colB := tbl.GetColumn("b")
	if !colB.Invisible {
		t.Error("column b should be invisible (Invisible=true)")
	}
}

func TestWalkThrough_3_1_CreateTableLike(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `
		CREATE TABLE src (
			id INT NOT NULL AUTO_INCREMENT,
			name VARCHAR(100) DEFAULT 'unnamed',
			PRIMARY KEY (id),
			INDEX idx_name (name)
		)
	`)
	wtExec(t, c, "CREATE TABLE dst LIKE src")

	srcTbl := c.GetDatabase("testdb").GetTable("src")
	dstTbl := c.GetDatabase("testdb").GetTable("dst")
	if dstTbl == nil {
		t.Fatal("table 'dst' not found after CREATE TABLE ... LIKE")
	}

	// Columns should match.
	if len(dstTbl.Columns) != len(srcTbl.Columns) {
		t.Fatalf("expected %d columns, got %d", len(srcTbl.Columns), len(dstTbl.Columns))
	}
	for i, srcCol := range srcTbl.Columns {
		dstCol := dstTbl.Columns[i]
		if dstCol.Name != srcCol.Name {
			t.Errorf("column %d: expected name %q, got %q", i, srcCol.Name, dstCol.Name)
		}
		if dstCol.ColumnType != srcCol.ColumnType {
			t.Errorf("column %q: expected type %q, got %q", srcCol.Name, srcCol.ColumnType, dstCol.ColumnType)
		}
		if dstCol.Nullable != srcCol.Nullable {
			t.Errorf("column %q: expected Nullable=%v, got %v", srcCol.Name, srcCol.Nullable, dstCol.Nullable)
		}
	}

	// Indexes should match.
	if len(dstTbl.Indexes) != len(srcTbl.Indexes) {
		t.Fatalf("expected %d indexes, got %d", len(srcTbl.Indexes), len(dstTbl.Indexes))
	}
	for i, srcIdx := range srcTbl.Indexes {
		dstIdx := dstTbl.Indexes[i]
		if dstIdx.Name != srcIdx.Name {
			t.Errorf("index %d: expected name %q, got %q", i, srcIdx.Name, dstIdx.Name)
		}
		if dstIdx.Primary != srcIdx.Primary {
			t.Errorf("index %q: expected Primary=%v, got %v", srcIdx.Name, srcIdx.Primary, dstIdx.Primary)
		}
	}

	// Constraints should match.
	if len(dstTbl.Constraints) != len(srcTbl.Constraints) {
		t.Fatalf("expected %d constraints, got %d", len(srcTbl.Constraints), len(dstTbl.Constraints))
	}
}
