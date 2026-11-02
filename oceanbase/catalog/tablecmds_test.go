package catalog

import (
	"strings"
	"testing"
)

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

	// Check that FK has a backing index (MySQL uses constraint name when provided).
	var fkIdx *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "fk_dept" {
			fkIdx = idx
			break
		}
	}
	if fkIdx == nil {
		t.Fatal("FK backing index not found")
	}
}

// TestCreateTableSubpartitionTemplate pins that OceanBase's SUBPARTITION TEMPLATE reaches
// the catalog, for partitions the DDL names and for the ones PARTITIONS N generates. The
// parser accepts the clause, so a conversion that ignored it would give every partition one
// auto-generated subpartition and quietly change the layout. Both forms are accepted by
// 5.7.25-OceanBase_CE-v4.3.5.6.
func TestCreateTableSubpartitionTemplate(t *testing.T) {
	for _, tc := range []struct {
		name  string
		ddl   string
		parts []string
	}{
		{
			"explicit outer partitions",
			`CREATE TABLE t (a INT, b INT)
				PARTITION BY RANGE COLUMNS(a)
				SUBPARTITION BY KEY(b)
				SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1)
				(PARTITION q0 VALUES LESS THAN (10), PARTITION q1 VALUES LESS THAN (20))`,
			[]string{"q0", "q1"},
		},
		{
			// HASH with PARTITIONS N names no partitions of its own, so these are
			// generated after the template would otherwise have been applied.
			"generated outer partitions",
			`CREATE TABLE t (a INT, b INT)
				PARTITION BY HASH(a) PARTITIONS 2
				SUBPARTITION BY KEY(b)
				SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1)`,
			[]string{"p0", "p1"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := setupWithDB(t)
			mustExec(t, c, tc.ddl)

			tbl := c.GetDatabase("testdb").GetTable("t")
			if tbl == nil {
				t.Fatal("table t not found")
			}
			pi := tbl.Partitioning
			if pi == nil {
				t.Fatal("table t has no partitioning")
			}
			if pi.NumSubParts != 2 {
				t.Errorf("NumSubParts = %d, want 2", pi.NumSubParts)
			}
			if pi.UseDefaultSubpartitions || pi.UseDefaultNumSubpartitions {
				t.Errorf("the template is explicit, but UseDefaultSubpartitions=%v UseDefaultNumSubpartitions=%v",
					pi.UseDefaultSubpartitions, pi.UseDefaultNumSubpartitions)
			}
			if len(pi.Partitions) != len(tc.parts) {
				t.Fatalf("partitions = %d, want %d", len(pi.Partitions), len(tc.parts))
			}
			for i, pd := range pi.Partitions {
				if pd.Name != tc.parts[i] {
					t.Errorf("partition %d named %q, want %q", i, pd.Name, tc.parts[i])
				}
				var got []string
				for _, sp := range pd.SubPartitions {
					got = append(got, sp.Name)
				}
				if len(got) != 2 || got[0] != "s0" || got[1] != "s1" {
					t.Errorf("partition %s subpartitions = %v, want [s0 s1]", pd.Name, got)
				}
			}
		})
	}
}

// TestCreateTableGlobalUniqueKeyPartitionCoverage pins which unique indexes must cover the
// partitioning columns. A GLOBAL index is not partitioned with the table, so it enforces
// uniqueness across all partitions and needs none of them; LOCAL and an unstated scope do.
// Each expectation below is the verdict of a live 5.7.25-OceanBase_CE-v4.3.5.6.
func TestCreateTableGlobalUniqueKeyPartitionCoverage(t *testing.T) {
	for _, tc := range []struct {
		name       string
		ddl        string
		wantAccept bool
	}{
		{"GLOBAL excluding the partition column", "UNIQUE KEY ua (a) GLOBAL", true},
		{"GLOBAL covering the partition column", "UNIQUE KEY ua (id, a) GLOBAL", true},
		{"LOCAL excluding the partition column", "UNIQUE KEY ua (a) LOCAL", false},
		{"unstated excluding the partition column", "UNIQUE KEY ua (a)", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := setupWithDB(t)
			results, err := c.Exec("CREATE TABLE t (id INT PRIMARY KEY, a INT, "+tc.ddl+") PARTITION BY HASH(id) PARTITIONS 2", nil)
			for _, r := range results {
				if r.Error != nil {
					err = r.Error
				}
			}
			if gotAccept := err == nil; gotAccept != tc.wantAccept {
				t.Errorf("accepted = %v, want %v (err: %v)", gotAccept, tc.wantAccept, err)
			}
		})
	}
}

// TestCreateTableLikeKeepsIndexScope pins that CREATE TABLE LIKE carries an index scope
// across. The copy lists index fields one by one, so a new field is silently dropped
// unless it is added there. 5.7.25-OceanBase_CE-v4.3.5.6 preserves the scope through LIKE,
// which is why copying it is the right behaviour and not an invention.
func TestCreateTableLikeKeepsIndexScope(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, "CREATE TABLE src (id INT, a INT, PRIMARY KEY (id), KEY ka (a) GLOBAL) PARTITION BY KEY(id) PARTITIONS 2")
	mustExec(t, c, "CREATE TABLE dst LIKE src")

	for _, name := range []string{"src", "dst"} {
		tbl := c.GetDatabase("testdb").GetTable(name)
		if tbl == nil {
			t.Fatalf("table %s not found", name)
		}
		var idx *Index
		for _, i := range tbl.Indexes {
			if i.Name == "ka" {
				idx = i
			}
		}
		if idx == nil {
			t.Fatalf("%s: index ka not found", name)
		}
		if idx.Scope != "GLOBAL" {
			t.Errorf("%s: Scope = %q, want GLOBAL", name, idx.Scope)
		}
	}
	if got := c.ShowCreateTable("testdb", "dst"); !strings.Contains(got, "GLOBAL") {
		t.Errorf("the copy's SHOW CREATE TABLE dropped the scope:\n%s", got)
	}
}
