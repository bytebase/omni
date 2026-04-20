package catalog

import (
	"strings"
	"testing"
)

// --- Section 9.2: SHOW CREATE TABLE Fidelity (9 scenarios) ---
// File target: wt_13_2_test.go
// Proof: go test ./mysql/catalog/ -short -count=1 -run "TestWalkThrough_13_2"

func TestWalkThrough_13_2_ShowCreateTableFidelity(t *testing.T) {
	// Scenario 1: Table with no explicit options — ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 rendered
	t.Run("no_explicit_options", func(t *testing.T) {
		t.Skip("TiDB behavioral difference: default collation is utf8mb4_bin, not MySQL 8.0's utf8mb4_0900_ai_ci.")
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			name VARCHAR(100),
			PRIMARY KEY (id)
		)`)

		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "ENGINE=InnoDB")
		assertContains(t, got, "DEFAULT CHARSET=utf8mb4")
		assertContains(t, got, "COLLATE=utf8mb4_0900_ai_ci")
	})

	// Scenario 2: Table with ROW_FORMAT=DYNAMIC
	t.Run("row_format_dynamic", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			PRIMARY KEY (id)
		) ROW_FORMAT=DYNAMIC`)

		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "ROW_FORMAT=DYNAMIC")
		assertContains(t, got, "ENGINE=InnoDB")
	})

	// Scenario 3: Table with KEY_BLOCK_SIZE=8
	t.Run("key_block_size", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			PRIMARY KEY (id)
		) KEY_BLOCK_SIZE=8`)

		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "KEY_BLOCK_SIZE=8")
		assertContains(t, got, "ENGINE=InnoDB")
	})

	// Scenario 4: Table with COMMENT='description'
	t.Run("table_comment", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			PRIMARY KEY (id)
		) COMMENT='this is a description'`)

		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "COMMENT='this is a description'")
	})

	// Scenario 5: Table with AUTO_INCREMENT=1000
	t.Run("auto_increment_1000", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL AUTO_INCREMENT,
			PRIMARY KEY (id)
		) AUTO_INCREMENT=1000`)

		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "AUTO_INCREMENT=1000")
	})

	// Scenario 6: TEMPORARY TABLE — SHOW CREATE TABLE works
	t.Run("temporary_table", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TEMPORARY TABLE t (
			id INT NOT NULL,
			name VARCHAR(50),
			PRIMARY KEY (id)
		)`)

		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "CREATE TEMPORARY TABLE `t`")
		assertContains(t, got, "`id` int NOT NULL")
		assertContains(t, got, "PRIMARY KEY (`id`)")
	})

	// Scenario 7: Table with all column types
	t.Run("all_column_types", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			col_bigint BIGINT,
			col_decimal DECIMAL(10,2),
			col_varchar VARCHAR(255),
			col_text TEXT,
			col_blob BLOB,
			col_json JSON,
			col_enum ENUM('a','b','c'),
			col_set SET('x','y','z'),
			col_date DATE,
			col_datetime DATETIME,
			col_timestamp TIMESTAMP NULL,
			PRIMARY KEY (id)
		)`)

		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "`col_bigint` bigint")
		assertContains(t, got, "`col_decimal` decimal(10,2)")
		assertContains(t, got, "`col_varchar` varchar(255)")
		assertContains(t, got, "`col_text` text")
		assertContains(t, got, "`col_blob` blob")
		assertContains(t, got, "`col_json` json")
		assertContains(t, got, "`col_enum` enum('a','b','c')")
		assertContains(t, got, "`col_set` set('x','y','z')")
		assertContains(t, got, "`col_date` date")
		assertContains(t, got, "`col_datetime` datetime")
		assertContains(t, got, "`col_timestamp` timestamp")

		// INT rendered
		assertContains(t, got, "`id` int NOT NULL")

		// TEXT and BLOB should NOT show DEFAULT NULL
		if strings.Contains(got, "`col_text` text DEFAULT NULL") {
			t.Error("TEXT column should not show DEFAULT NULL")
		}
		if strings.Contains(got, "`col_blob` blob DEFAULT NULL") {
			t.Error("BLOB column should not show DEFAULT NULL")
		}

		// TIMESTAMP with explicit NULL should show NULL
		assertContains(t, got, "`col_timestamp` timestamp NULL")
	})

	// Scenario 8: Column with ON UPDATE CURRENT_TIMESTAMP
	t.Run("on_update_current_timestamp", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		)`)

		got := c.ShowCreateTable("testdb", "t")
		// ON UPDATE should appear after DEFAULT
		assertContains(t, got, "DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
	})

	// Scenario 9: Column DEFAULT expression (CURRENT_TIMESTAMP, literal, NULL)
	t.Run("default_expressions", func(t *testing.T) {
		t.Skip("TiDB behavioral difference: utf8mb4_bin default produces per-column COLLATE clauses in SHOW CREATE output, so `varchar(100) DEFAULT NULL` appears as `varchar(100) COLLATE utf8mb4_bin DEFAULT NULL`.")
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
			status INT DEFAULT 0,
			note VARCHAR(100) DEFAULT NULL,
			PRIMARY KEY (id)
		)`)

		got := c.ShowCreateTable("testdb", "t")
		// CURRENT_TIMESTAMP default
		assertContains(t, got, "DEFAULT CURRENT_TIMESTAMP")
		// Literal default — MySQL 8.0 quotes numeric defaults
		assertContains(t, got, "DEFAULT '0'")
		// NULL default
		assertContains(t, got, "`note` varchar(100) DEFAULT NULL")
	})
}
