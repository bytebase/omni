package catalog

import (
	"strings"
	"testing"
)

// --- Section 4.1 (starmap): Charset Inheritance Chain (7 scenarios) ---

func TestWalkThrough_8_1(t *testing.T) {
	t.Run("db_charset_inherited_by_table_and_column", func(t *testing.T) {
		// Scenario 1: CREATE DATABASE CHARSET utf8mb4 → CREATE TABLE (no charset) → column inherits utf8mb4
		c := New()
		mustExec(t, c, "CREATE DATABASE db1 DEFAULT CHARACTER SET utf8mb4")
		c.SetCurrentDatabase("db1")
		mustExec(t, c, "CREATE TABLE t1 (id INT, name VARCHAR(100))")

		db := c.GetDatabase("db1")
		tbl := db.GetTable("t1")
		if tbl.Charset != "utf8mb4" {
			t.Errorf("table charset: expected utf8mb4, got %q", tbl.Charset)
		}
		col := tbl.GetColumn("name")
		if col.Charset != "utf8mb4" {
			t.Errorf("column charset: expected utf8mb4, got %q", col.Charset)
		}
		if col.Collation != "utf8mb4_0900_ai_ci" {
			t.Errorf("column collation: expected utf8mb4_0900_ai_ci, got %q", col.Collation)
		}
	})

	t.Run("table_charset_overrides_db", func(t *testing.T) {
		// Scenario 2: CREATE DATABASE CHARSET latin1 → CREATE TABLE CHARSET utf8mb4 → column inherits utf8mb4
		c := New()
		mustExec(t, c, "CREATE DATABASE db2 DEFAULT CHARACTER SET latin1")
		c.SetCurrentDatabase("db2")
		mustExec(t, c, "CREATE TABLE t1 (id INT, name VARCHAR(100)) DEFAULT CHARSET=utf8mb4")

		db := c.GetDatabase("db2")
		tbl := db.GetTable("t1")
		if tbl.Charset != "utf8mb4" {
			t.Errorf("table charset: expected utf8mb4, got %q", tbl.Charset)
		}
		col := tbl.GetColumn("name")
		if col.Charset != "utf8mb4" {
			t.Errorf("column charset: expected utf8mb4, got %q", col.Charset)
		}
	})

	t.Run("add_column_inherits_table_charset", func(t *testing.T) {
		// Scenario 3: Table CHARSET utf8mb4 → ADD COLUMN VARCHAR (no charset) → column inherits table charset
		c := New()
		mustExec(t, c, "CREATE DATABASE db3 DEFAULT CHARACTER SET utf8mb4")
		c.SetCurrentDatabase("db3")
		mustExec(t, c, "CREATE TABLE t1 (id INT) DEFAULT CHARSET=utf8mb4")
		mustExec(t, c, "ALTER TABLE t1 ADD COLUMN name VARCHAR(100)")

		tbl := c.GetDatabase("db3").GetTable("t1")
		col := tbl.GetColumn("name")
		if col.Charset != "utf8mb4" {
			t.Errorf("column charset: expected utf8mb4, got %q", col.Charset)
		}
		if col.Collation != "utf8mb4_0900_ai_ci" {
			t.Errorf("column collation: expected utf8mb4_0900_ai_ci, got %q", col.Collation)
		}
	})

	t.Run("add_column_charset_overrides_table", func(t *testing.T) {
		// Scenario 4: Table CHARSET utf8mb4 → ADD COLUMN VARCHAR CHARSET latin1 → column overrides
		c := New()
		mustExec(t, c, "CREATE DATABASE db4 DEFAULT CHARACTER SET utf8mb4")
		c.SetCurrentDatabase("db4")
		mustExec(t, c, "CREATE TABLE t1 (id INT) DEFAULT CHARSET=utf8mb4")
		mustExec(t, c, "ALTER TABLE t1 ADD COLUMN name VARCHAR(100) CHARACTER SET latin1")

		tbl := c.GetDatabase("db4").GetTable("t1")
		col := tbl.GetColumn("name")
		if col.Charset != "latin1" {
			t.Errorf("column charset: expected latin1, got %q", col.Charset)
		}
		if col.Collation != "latin1_swedish_ci" {
			t.Errorf("column collation: expected latin1_swedish_ci, got %q", col.Collation)
		}
	})

	t.Run("show_create_table_shows_inherited_charset", func(t *testing.T) {
		// Scenario 5: CREATE DATABASE CHARSET utf8mb4 → table inherits → SHOW CREATE TABLE shows DEFAULT CHARSET=utf8mb4
		c := New()
		mustExec(t, c, "CREATE DATABASE db5 DEFAULT CHARACTER SET utf8mb4")
		c.SetCurrentDatabase("db5")
		mustExec(t, c, "CREATE TABLE t1 (id INT, name VARCHAR(100))")

		ddl := c.ShowCreateTable("db5", "t1")
		if !strings.Contains(ddl, "DEFAULT CHARSET=utf8mb4") {
			t.Errorf("SHOW CREATE TABLE should contain DEFAULT CHARSET=utf8mb4, got:\n%s", ddl)
		}
		// MySQL 8.0 always shows COLLATE for utf8mb4
		if !strings.Contains(ddl, "COLLATE=utf8mb4_0900_ai_ci") {
			t.Errorf("SHOW CREATE TABLE should contain COLLATE=utf8mb4_0900_ai_ci for utf8mb4, got:\n%s", ddl)
		}
	})

	t.Run("charset_only_derives_default_collation", func(t *testing.T) {
		// Scenario 6: CREATE TABLE with CHARSET only (no COLLATE) — default collation derived
		c := New()
		mustExec(t, c, "CREATE DATABASE db6")
		c.SetCurrentDatabase("db6")
		mustExec(t, c, "CREATE TABLE t1 (id INT, name VARCHAR(100)) DEFAULT CHARSET=latin1")

		tbl := c.GetDatabase("db6").GetTable("t1")
		if tbl.Charset != "latin1" {
			t.Errorf("table charset: expected latin1, got %q", tbl.Charset)
		}
		if tbl.Collation != "latin1_swedish_ci" {
			t.Errorf("table collation: expected latin1_swedish_ci, got %q", tbl.Collation)
		}

		// Column should also inherit
		col := tbl.GetColumn("name")
		if col.Charset != "latin1" {
			t.Errorf("column charset: expected latin1, got %q", col.Charset)
		}
		if col.Collation != "latin1_swedish_ci" {
			t.Errorf("column collation: expected latin1_swedish_ci, got %q", col.Collation)
		}
	})

	t.Run("collate_only_derives_charset", func(t *testing.T) {
		// Scenario 7: CREATE TABLE with COLLATE only (no CHARSET) — charset derived from collation
		c := New()
		mustExec(t, c, "CREATE DATABASE db7")
		c.SetCurrentDatabase("db7")
		mustExec(t, c, "CREATE TABLE t1 (id INT, name VARCHAR(100)) DEFAULT COLLATE=latin1_swedish_ci")

		tbl := c.GetDatabase("db7").GetTable("t1")
		if tbl.Charset != "latin1" {
			t.Errorf("table charset: expected latin1, got %q", tbl.Charset)
		}
		if tbl.Collation != "latin1_swedish_ci" {
			t.Errorf("table collation: expected latin1_swedish_ci, got %q", tbl.Collation)
		}

		ddl := c.ShowCreateTable("db7", "t1")
		if !strings.Contains(ddl, "DEFAULT CHARSET=latin1") {
			t.Errorf("SHOW CREATE TABLE should contain DEFAULT CHARSET=latin1, got:\n%s", ddl)
		}
	})
}
