package catalog

import (
	"strings"
	"testing"
)

// TestWalkThrough_8_2_AlterDefaultCharset tests that ALTER TABLE DEFAULT CHARACTER SET
// changes the table default but leaves existing column charsets unchanged.
func TestWalkThrough_8_2_AlterDefaultCharset(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t1 (
		id INT,
		name VARCHAR(100),
		bio TEXT
	) DEFAULT CHARSET=utf8mb4`)

	// Verify initial state: columns inherit utf8mb4.
	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl.Charset != "utf8mb4" {
		t.Fatalf("expected table charset=utf8mb4, got %s", tbl.Charset)
	}
	nameCol := tbl.GetColumn("name")
	if nameCol.Charset != "utf8mb4" {
		t.Fatalf("expected name charset=utf8mb4, got %s", nameCol.Charset)
	}

	// ALTER TABLE DEFAULT CHARACTER SET latin1 — table default changes but columns unchanged.
	wtExec(t, c, "ALTER TABLE t1 DEFAULT CHARACTER SET latin1")

	tbl = c.GetDatabase("testdb").GetTable("t1")
	if tbl.Charset != "latin1" {
		t.Errorf("expected table charset=latin1, got %s", tbl.Charset)
	}

	// Existing columns should still have utf8mb4.
	nameCol = tbl.GetColumn("name")
	if nameCol.Charset != "utf8mb4" {
		t.Errorf("expected name charset=utf8mb4 (unchanged), got %s", nameCol.Charset)
	}
	bioCol := tbl.GetColumn("bio")
	if bioCol.Charset != "utf8mb4" {
		t.Errorf("expected bio charset=utf8mb4 (unchanged), got %s", bioCol.Charset)
	}
}

// TestWalkThrough_8_2_ConvertToCharset tests CONVERT TO CHARACTER SET utf8mb4
// updates table + all VARCHAR/TEXT/ENUM columns.
func TestWalkThrough_8_2_ConvertToCharset(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t1 (
		id INT,
		name VARCHAR(100),
		bio TEXT,
		tag ENUM('a','b')
	) DEFAULT CHARSET=latin1`)

	wtExec(t, c, "ALTER TABLE t1 CONVERT TO CHARACTER SET utf8mb4")

	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl.Charset != "utf8mb4" {
		t.Errorf("expected table charset=utf8mb4, got %s", tbl.Charset)
	}
	if tbl.Collation != "utf8mb4_0900_ai_ci" {
		t.Errorf("expected table collation=utf8mb4_0900_ai_ci, got %s", tbl.Collation)
	}

	// All string columns should be updated.
	for _, colName := range []string{"name", "bio", "tag"} {
		col := tbl.GetColumn(colName)
		if col == nil {
			t.Fatalf("column %s not found", colName)
		}
		if col.Charset != "utf8mb4" {
			t.Errorf("column %s: expected charset=utf8mb4, got %s", colName, col.Charset)
		}
		if col.Collation != "utf8mb4_0900_ai_ci" {
			t.Errorf("column %s: expected collation=utf8mb4_0900_ai_ci, got %s", colName, col.Collation)
		}
	}
}

// TestWalkThrough_8_2_ConvertWithCollation tests CONVERT TO CHARACTER SET utf8mb4
// COLLATE utf8mb4_unicode_ci — non-default collation.
func TestWalkThrough_8_2_ConvertWithCollation(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t1 (
		id INT,
		name VARCHAR(100),
		bio TEXT
	) DEFAULT CHARSET=latin1`)

	wtExec(t, c, "ALTER TABLE t1 CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")

	tbl := c.GetDatabase("testdb").GetTable("t1")
	if tbl.Charset != "utf8mb4" {
		t.Errorf("expected table charset=utf8mb4, got %s", tbl.Charset)
	}
	if tbl.Collation != "utf8mb4_unicode_ci" {
		t.Errorf("expected table collation=utf8mb4_unicode_ci, got %s", tbl.Collation)
	}

	for _, colName := range []string{"name", "bio"} {
		col := tbl.GetColumn(colName)
		if col.Charset != "utf8mb4" {
			t.Errorf("column %s: expected charset=utf8mb4, got %s", colName, col.Charset)
		}
		if col.Collation != "utf8mb4_unicode_ci" {
			t.Errorf("column %s: expected collation=utf8mb4_unicode_ci, got %s", colName, col.Collation)
		}
	}
}

// TestWalkThrough_8_2_ConvertIntColumnsUnchanged tests that CONVERT TO CHARACTER SET
// does not affect INT columns.
func TestWalkThrough_8_2_ConvertIntColumnsUnchanged(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t1 (
		id INT,
		count BIGINT,
		name VARCHAR(100)
	) DEFAULT CHARSET=latin1`)

	wtExec(t, c, "ALTER TABLE t1 CONVERT TO CHARACTER SET utf8mb4")

	tbl := c.GetDatabase("testdb").GetTable("t1")

	// INT columns should have no charset.
	colID := tbl.GetColumn("id")
	if colID.Charset != "" {
		t.Errorf("INT column id should not have charset, got %s", colID.Charset)
	}
	colCount := tbl.GetColumn("count")
	if colCount.Charset != "" {
		t.Errorf("BIGINT column count should not have charset, got %s", colCount.Charset)
	}

	// String column should be updated.
	colName := tbl.GetColumn("name")
	if colName.Charset != "utf8mb4" {
		t.Errorf("VARCHAR column name: expected charset=utf8mb4, got %s", colName.Charset)
	}
}

// TestWalkThrough_8_2_ConvertMixedColumnTypes tests CONVERT TO CHARACTER SET on
// a table with mixed column types — only string types are updated.
func TestWalkThrough_8_2_ConvertMixedColumnTypes(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t1 (
		id INT,
		name VARCHAR(100),
		score DECIMAL(10,2),
		bio TEXT,
		active TINYINT,
		tag ENUM('a','b'),
		created DATE,
		data MEDIUMTEXT
	) DEFAULT CHARSET=latin1`)

	wtExec(t, c, "ALTER TABLE t1 CONVERT TO CHARACTER SET utf8mb4")

	tbl := c.GetDatabase("testdb").GetTable("t1")

	// String types should be updated.
	stringCols := []string{"name", "bio", "tag", "data"}
	for _, colName := range stringCols {
		col := tbl.GetColumn(colName)
		if col == nil {
			t.Fatalf("column %s not found", colName)
		}
		if col.Charset != "utf8mb4" {
			t.Errorf("column %s: expected charset=utf8mb4, got %s", colName, col.Charset)
		}
	}

	// Non-string types should not have charset.
	nonStringCols := []string{"id", "score", "active", "created"}
	for _, colName := range nonStringCols {
		col := tbl.GetColumn(colName)
		if col == nil {
			t.Fatalf("column %s not found", colName)
		}
		if col.Charset != "" {
			t.Errorf("non-string column %s should not have charset, got %s", colName, col.Charset)
		}
	}
}

// TestWalkThrough_8_2_ConvertOverwritesExplicitCharset tests that CONVERT TO CHARACTER SET
// overwrites a column that already has an explicit charset.
func TestWalkThrough_8_2_ConvertOverwritesExplicitCharset(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t1 (
		id INT,
		name VARCHAR(100) CHARACTER SET latin1
	) DEFAULT CHARSET=utf8mb4`)

	// Verify initial: name has explicit latin1.
	tbl := c.GetDatabase("testdb").GetTable("t1")
	colName := tbl.GetColumn("name")
	if colName.Charset != "latin1" {
		t.Fatalf("expected name charset=latin1 initially, got %s", colName.Charset)
	}

	// CONVERT overwrites the explicit charset.
	wtExec(t, c, "ALTER TABLE t1 CONVERT TO CHARACTER SET utf8mb4")

	tbl = c.GetDatabase("testdb").GetTable("t1")
	colName = tbl.GetColumn("name")
	if colName.Charset != "utf8mb4" {
		t.Errorf("expected name charset=utf8mb4 after CONVERT, got %s", colName.Charset)
	}
	if colName.Collation != "utf8mb4_0900_ai_ci" {
		t.Errorf("expected name collation=utf8mb4_0900_ai_ci after CONVERT, got %s", colName.Collation)
	}
}

// TestWalkThrough_8_2_ConvertThenShowCreate tests that after CONVERT TO CHARACTER SET,
// column charsets matching the table default are NOT shown in SHOW CREATE TABLE.
func TestWalkThrough_8_2_ConvertThenShowCreate(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t1 (
		id INT,
		name VARCHAR(100),
		bio TEXT,
		age INT
	) DEFAULT CHARSET=latin1`)

	wtExec(t, c, "ALTER TABLE t1 CONVERT TO CHARACTER SET utf8mb4")

	ddl := c.ShowCreateTable("testdb", "t1")

	// Column charsets should NOT appear in the output since they match the table default.
	// The table-level DEFAULT CHARSET=utf8mb4 should appear.
	if !strings.Contains(ddl, "DEFAULT CHARSET=utf8mb4") {
		t.Errorf("expected DEFAULT CHARSET=utf8mb4 in output, got:\n%s", ddl)
	}

	// Column definitions should NOT contain "CHARACTER SET" since they match the table default.
	lines := strings.Split(ddl, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip table options line.
		if strings.HasPrefix(trimmed, ")") {
			continue
		}
		// Check column lines for CHARACTER SET — should not appear.
		if strings.HasPrefix(trimmed, "`name`") || strings.HasPrefix(trimmed, "`bio`") {
			if strings.Contains(line, "CHARACTER SET") {
				t.Errorf("column charset should not be shown when matching table default:\n%s", line)
			}
		}
	}

	// INT columns should definitely not show CHARACTER SET.
	for _, line := range lines {
		if strings.Contains(line, "`id`") || strings.Contains(line, "`age`") {
			if strings.Contains(line, "CHARACTER SET") {
				t.Errorf("INT column should not show CHARACTER SET:\n%s", line)
			}
		}
	}
}
