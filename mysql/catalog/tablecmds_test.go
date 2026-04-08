package catalog

import "testing"

func mustExec(t *testing.T, c *Catalog, sql string) {
	t.Helper()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("exec error: %v", r.Error)
		}
	}
}

func setupWithDB(t *testing.T) *Catalog {
	t.Helper()
	c := New()
	mustExec(t, c, "CREATE DATABASE testdb")
	c.SetCurrentDatabase("testdb")
	return c
}

func TestCreateTableBasic(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE users (
		id INT NOT NULL AUTO_INCREMENT,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(255),
		age INT UNSIGNED DEFAULT 0,
		score DECIMAL(10,2),
		PRIMARY KEY (id)
	)`)

	db := c.GetDatabase("testdb")
	tbl := db.GetTable("users")
	if tbl == nil {
		t.Fatal("table users not found")
	}
	if len(tbl.Columns) != 5 {
		t.Fatalf("expected 5 columns, got %d", len(tbl.Columns))
	}

	// Check id column.
	id := tbl.GetColumn("id")
	if id == nil {
		t.Fatal("column id not found")
	}
	if id.Nullable {
		t.Error("id should not be nullable")
	}
	if !id.AutoIncrement {
		t.Error("id should be auto_increment")
	}
	if id.DataType != "int" {
		t.Errorf("expected data type 'int', got %q", id.DataType)
	}
	if id.Position != 1 {
		t.Errorf("expected position 1, got %d", id.Position)
	}

	// Check name column.
	name := tbl.GetColumn("name")
	if name == nil {
		t.Fatal("column name not found")
	}
	if name.Nullable {
		t.Error("name should not be nullable")
	}
	if name.ColumnType != "varchar(100)" {
		t.Errorf("expected column type 'varchar(100)', got %q", name.ColumnType)
	}

	// Check email column (nullable by default).
	email := tbl.GetColumn("email")
	if email == nil {
		t.Fatal("column email not found")
	}
	if !email.Nullable {
		t.Error("email should be nullable by default")
	}

	// Check age column (unsigned, default).
	age := tbl.GetColumn("age")
	if age == nil {
		t.Fatal("column age not found")
	}
	if age.ColumnType != "int unsigned" {
		t.Errorf("expected column type 'int unsigned', got %q", age.ColumnType)
	}
	if age.Default == nil || *age.Default != "0" {
		t.Errorf("expected default '0', got %v", age.Default)
	}

	// Check score column.
	score := tbl.GetColumn("score")
	if score == nil {
		t.Fatal("column score not found")
	}
	if score.ColumnType != "decimal(10,2)" {
		t.Errorf("expected column type 'decimal(10,2)', got %q", score.ColumnType)
	}

	// Check PK index.
	if len(tbl.Indexes) < 1 {
		t.Fatal("expected at least 1 index")
	}
	pkIdx := tbl.Indexes[0]
	if pkIdx.Name != "PRIMARY" {
		t.Errorf("expected PK index name 'PRIMARY', got %q", pkIdx.Name)
	}
	if !pkIdx.Primary {
		t.Error("expected primary flag on PK index")
	}
	if !pkIdx.Unique {
		t.Error("expected unique flag on PK index")
	}
}

func TestCreateTableDuplicate(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, "CREATE TABLE t1 (id INT)")
	results, _ := c.Exec("CREATE TABLE t1 (id INT)", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected duplicate table error")
	}
	catErr, ok := results[0].Error.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", results[0].Error)
	}
	if catErr.Code != ErrDupTable {
		t.Errorf("expected error code %d, got %d", ErrDupTable, catErr.Code)
	}
}

func TestCreateTableIfNotExists(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, "CREATE TABLE t1 (id INT)")
	results, _ := c.Exec("CREATE TABLE IF NOT EXISTS t1 (id INT)", nil)
	if results[0].Error != nil {
		t.Errorf("IF NOT EXISTS should not error: %v", results[0].Error)
	}
}

func TestCreateTableDupColumn(t *testing.T) {
	c := setupWithDB(t)
	results, _ := c.Exec("CREATE TABLE t1 (id INT, id VARCHAR(10))", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected duplicate column error")
	}
	catErr, ok := results[0].Error.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", results[0].Error)
	}
	if catErr.Code != ErrDupColumn {
		t.Errorf("expected error code %d, got %d", ErrDupColumn, catErr.Code)
	}
}

func TestCreateTableMultiplePK(t *testing.T) {
	c := setupWithDB(t)
	results, _ := c.Exec(`CREATE TABLE t1 (
		id INT PRIMARY KEY,
		name VARCHAR(100),
		PRIMARY KEY (name)
	)`, &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected multiple primary key error")
	}
	catErr, ok := results[0].Error.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", results[0].Error)
	}
	if catErr.Code != ErrMultiplePriKey {
		t.Errorf("expected error code %d, got %d", ErrMultiplePriKey, catErr.Code)
	}
}

func TestCreateTableNoDatabaseSelected(t *testing.T) {
	c := New()
	results, _ := c.Exec("CREATE TABLE t1 (id INT)", &ExecOptions{ContinueOnError: true})
	if results[0].Error == nil {
		t.Fatal("expected no database selected error")
	}
	catErr, ok := results[0].Error.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", results[0].Error)
	}
	if catErr.Code != ErrNoDatabaseSelected {
		t.Errorf("expected error code %d, got %d", ErrNoDatabaseSelected, catErr.Code)
	}
}

func TestCreateTableWithIndexes(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE t1 (
		id INT NOT NULL AUTO_INCREMENT,
		email VARCHAR(255) NOT NULL,
		name VARCHAR(100),
		PRIMARY KEY (id),
		UNIQUE KEY idx_email (email),
		INDEX idx_name (name)
	)`)

	db := c.GetDatabase("testdb")
	tbl := db.GetTable("t1")
	if tbl == nil {
		t.Fatal("table t1 not found")
	}

	// Should have 3 indexes: PRIMARY, idx_email, idx_name.
	if len(tbl.Indexes) != 3 {
		t.Fatalf("expected 3 indexes, got %d", len(tbl.Indexes))
	}

	// Check UNIQUE KEY.
	var uniqueIdx *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_email" {
			uniqueIdx = idx
			break
		}
	}
	if uniqueIdx == nil {
		t.Fatal("unique index idx_email not found")
	}
	if !uniqueIdx.Unique {
		t.Error("idx_email should be unique")
	}
	if len(uniqueIdx.Columns) != 1 || uniqueIdx.Columns[0].Name != "email" {
		t.Errorf("idx_email should have column 'email'")
	}

	// Check regular INDEX.
	var regIdx *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_name" {
			regIdx = idx
			break
		}
	}
	if regIdx == nil {
		t.Fatal("index idx_name not found")
	}
	if regIdx.Unique {
		t.Error("idx_name should not be unique")
	}
}

func TestCreateTableWithFK(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, "CREATE TABLE departments (id INT NOT NULL, PRIMARY KEY (id))")
	mustExec(t, c, `CREATE TABLE employees (
		id INT NOT NULL AUTO_INCREMENT,
		dept_id INT NOT NULL,
		PRIMARY KEY (id),
		CONSTRAINT fk_dept FOREIGN KEY (dept_id) REFERENCES departments(id) ON DELETE CASCADE
	)`)

	db := c.GetDatabase("testdb")
	tbl := db.GetTable("employees")
	if tbl == nil {
		t.Fatal("table employees not found")
	}

	// Check FK constraint.
	var fk *Constraint
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey {
			fk = con
			break
		}
	}
	if fk == nil {
		t.Fatal("FK constraint not found")
	}
	if fk.Name != "fk_dept" {
		t.Errorf("expected FK name 'fk_dept', got %q", fk.Name)
	}
	if fk.RefTable != "departments" {
		t.Errorf("expected ref table 'departments', got %q", fk.RefTable)
	}
	if len(fk.RefColumns) != 1 || fk.RefColumns[0] != "id" {
		t.Errorf("expected ref column 'id', got %v", fk.RefColumns)
	}
	if len(fk.Columns) != 1 || fk.Columns[0] != "dept_id" {
		t.Errorf("expected column 'dept_id', got %v", fk.Columns)
	}
	if fk.OnDelete != "CASCADE" {
		t.Errorf("expected ON DELETE CASCADE, got %q", fk.OnDelete)
	}

	// Check that FK has a backing index (MySQL uses column name, not constraint name).
	var fkIdx *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "dept_id" {
			fkIdx = idx
			break
		}
	}
	if fkIdx == nil {
		t.Fatal("FK backing index not found")
	}
}
