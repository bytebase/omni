package catalog

import (
	"strings"
	"testing"
)

// --- Section 5.2: Generated Column Dependencies (6 scenarios) ---

// Scenario 1: DROP COLUMN referenced by VIRTUAL generated column — MySQL 8.0 error 3108
func TestWalkThrough_9_2_DropColumnReferencedByVirtual(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t (
		id INT,
		a INT,
		b INT,
		v INT GENERATED ALWAYS AS (a + b) VIRTUAL
	)`)

	// Dropping column 'a' should fail because generated column 'v' references it.
	results, _ := c.Exec("ALTER TABLE t DROP COLUMN a", &ExecOptions{ContinueOnError: true})
	assertError(t, results[0].Error, ErrDependentByGenCol)

	// Verify column was NOT dropped (error should prevent it).
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl.GetColumn("a") == nil {
		t.Error("column 'a' should still exist after failed DROP")
	}
	if len(tbl.Columns) != 4 {
		t.Errorf("expected 4 columns, got %d", len(tbl.Columns))
	}
}

// Scenario 2: DROP COLUMN referenced by STORED generated column — MySQL 8.0 error 3108
func TestWalkThrough_9_2_DropColumnReferencedByStored(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t (
		id INT,
		a INT,
		b INT,
		s INT GENERATED ALWAYS AS (a + b) STORED
	)`)

	// Dropping column 'b' should fail because stored generated column 's' references it.
	results, _ := c.Exec("ALTER TABLE t DROP COLUMN b", &ExecOptions{ContinueOnError: true})
	assertError(t, results[0].Error, ErrDependentByGenCol)

	// Verify column was NOT dropped.
	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl.GetColumn("b") == nil {
		t.Error("column 'b' should still exist after failed DROP")
	}
}

// Scenario 3: MODIFY base column type when generated column uses it — expression preserved, no error
func TestWalkThrough_9_2_ModifyBaseColumnType(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t (
		id INT,
		a INT,
		v INT GENERATED ALWAYS AS (a * 2) VIRTUAL
	)`)

	// Modify base column type — should succeed, generated expression preserved.
	wtExec(t, c, "ALTER TABLE t MODIFY COLUMN a BIGINT")

	tbl := c.GetDatabase("testdb").GetTable("t")

	// Verify column type was changed.
	col := tbl.GetColumn("a")
	if col == nil {
		t.Fatal("column 'a' not found after MODIFY")
	}
	if col.ColumnType != "bigint" {
		t.Errorf("expected column type 'bigint', got %q", col.ColumnType)
	}

	// Verify generated column expression is preserved.
	vCol := tbl.GetColumn("v")
	if vCol == nil {
		t.Fatal("generated column 'v' not found")
	}
	if vCol.Generated == nil {
		t.Fatal("column 'v' should be a generated column")
	}

	// Verify SHOW CREATE TABLE renders correctly.
	show := c.ShowCreateTable("testdb", "t")
	if !strings.Contains(show, "GENERATED ALWAYS AS") {
		t.Error("SHOW CREATE TABLE should contain generated column expression")
	}
	if !strings.Contains(show, "`a`") {
		t.Error("generated expression should still reference column 'a'")
	}
}

// Scenario 4: Generated column referencing another generated column — verify creation and SHOW CREATE
func TestWalkThrough_9_2_GeneratedReferencingGenerated(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t (
		id INT,
		a INT,
		v1 INT GENERATED ALWAYS AS (a * 2) VIRTUAL,
		v2 INT GENERATED ALWAYS AS (v1 + 10) VIRTUAL
	)`)

	tbl := c.GetDatabase("testdb").GetTable("t")

	// Verify both generated columns exist.
	v1 := tbl.GetColumn("v1")
	if v1 == nil {
		t.Fatal("generated column 'v1' not found")
	}
	if v1.Generated == nil {
		t.Fatal("v1 should be a generated column")
	}
	if v1.Generated.Stored {
		t.Error("v1 should be VIRTUAL, not STORED")
	}

	v2 := tbl.GetColumn("v2")
	if v2 == nil {
		t.Fatal("generated column 'v2' not found")
	}
	if v2.Generated == nil {
		t.Fatal("v2 should be a generated column")
	}
	if v2.Generated.Stored {
		t.Error("v2 should be VIRTUAL, not STORED")
	}

	// Verify SHOW CREATE TABLE renders both generated columns correctly.
	show := c.ShowCreateTable("testdb", "t")
	if !strings.Contains(show, "`v1`") {
		t.Error("SHOW CREATE TABLE should contain v1")
	}
	if !strings.Contains(show, "`v2`") {
		t.Error("SHOW CREATE TABLE should contain v2")
	}

	// v2 references v1, so dropping v1 should fail.
	results, _ := c.Exec("ALTER TABLE t DROP COLUMN v1", &ExecOptions{ContinueOnError: true})
	assertError(t, results[0].Error, ErrDependentByGenCol)
}

// Scenario 5: Index on generated column — index created, rendered correctly
func TestWalkThrough_9_2_IndexOnGeneratedColumn(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t (
		id INT,
		a INT,
		v INT GENERATED ALWAYS AS (a * 2) VIRTUAL,
		INDEX idx_v (v)
	)`)

	tbl := c.GetDatabase("testdb").GetTable("t")

	// Verify generated column exists.
	vCol := tbl.GetColumn("v")
	if vCol == nil {
		t.Fatal("generated column 'v' not found")
	}
	if vCol.Generated == nil {
		t.Fatal("v should be a generated column")
	}

	// Verify index exists on generated column.
	var found *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "idx_v" {
			found = idx
			break
		}
	}
	if found == nil {
		t.Fatal("index idx_v not found")
	}
	if found.Unique {
		t.Error("idx_v should not be unique")
	}
	if len(found.Columns) != 1 || found.Columns[0].Name != "v" {
		t.Errorf("index columns mismatch: %+v", found.Columns)
	}

	// Verify SHOW CREATE TABLE renders both the generated column and index.
	show := c.ShowCreateTable("testdb", "t")
	if !strings.Contains(show, "GENERATED ALWAYS AS") {
		t.Error("SHOW CREATE TABLE should contain generated column expression")
	}
	if !strings.Contains(show, "KEY `idx_v`") {
		t.Error("SHOW CREATE TABLE should contain index idx_v")
	}
}

// Scenario 6: UNIQUE index on generated column — unique constraint on generated column
func TestWalkThrough_9_2_UniqueIndexOnGeneratedColumn(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t (
		id INT,
		a INT,
		v INT GENERATED ALWAYS AS (a * 2) VIRTUAL,
		UNIQUE INDEX ux_v (v)
	)`)

	tbl := c.GetDatabase("testdb").GetTable("t")

	// Verify generated column exists.
	vCol := tbl.GetColumn("v")
	if vCol == nil {
		t.Fatal("generated column 'v' not found")
	}
	if vCol.Generated == nil {
		t.Fatal("v should be a generated column")
	}

	// Verify unique index exists on generated column.
	var found *Index
	for _, idx := range tbl.Indexes {
		if idx.Name == "ux_v" {
			found = idx
			break
		}
	}
	if found == nil {
		t.Fatal("unique index ux_v not found")
	}
	if !found.Unique {
		t.Error("ux_v should be unique")
	}
	if len(found.Columns) != 1 || found.Columns[0].Name != "v" {
		t.Errorf("unique index columns mismatch: %+v", found.Columns)
	}

	// Verify SHOW CREATE TABLE renders the unique index.
	show := c.ShowCreateTable("testdb", "t")
	if !strings.Contains(show, "UNIQUE KEY `ux_v`") {
		t.Error("SHOW CREATE TABLE should contain UNIQUE KEY ux_v")
	}
}
