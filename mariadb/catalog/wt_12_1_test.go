package catalog

import (
	"strings"
	"testing"
)

// --- Section 8.1 (Phase 8): Prefix and Expression Index Rendering (8 scenarios) ---
// File target: wt_12_1_test.go
// Proof: go test ./mysql/catalog/ -count=1 -run "TestWalkThrough_12_1"

func TestWalkThrough_12_1_PrefixAndExpressionIndexRendering(t *testing.T) {
	// Scenario 1: KEY idx (col(10)) — prefix length rendered in SHOW CREATE
	t.Run("prefix_length_rendered", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t1 (
			id INT NOT NULL,
			name VARCHAR(100),
			KEY idx_name (name(10)),
			PRIMARY KEY (id)
		)`)

		ddl := c.ShowCreateTable("testdb", "t1")
		if !strings.Contains(ddl, "`name`(10)") {
			t.Errorf("expected prefix index col(10) in SHOW CREATE TABLE, got:\n%s", ddl)
		}
		if !strings.Contains(ddl, "KEY `idx_name`") {
			t.Errorf("expected KEY idx_name in SHOW CREATE TABLE, got:\n%s", ddl)
		}
	})

	// Scenario 2: KEY idx (col1(10), col2(20)) — multi-column prefix index
	t.Run("multi_column_prefix", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t1 (
			id INT NOT NULL,
			col1 VARCHAR(100),
			col2 VARCHAR(100),
			KEY idx_multi (col1(10), col2(20)),
			PRIMARY KEY (id)
		)`)

		ddl := c.ShowCreateTable("testdb", "t1")
		if !strings.Contains(ddl, "`col1`(10)") {
			t.Errorf("expected col1(10) in SHOW CREATE TABLE, got:\n%s", ddl)
		}
		if !strings.Contains(ddl, "`col2`(20)") {
			t.Errorf("expected col2(20) in SHOW CREATE TABLE, got:\n%s", ddl)
		}
	})

	// Scenario 3: KEY idx (col(10), col2) — mixed prefix and full column
	t.Run("mixed_prefix_and_full_column", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t1 (
			id INT NOT NULL,
			col1 VARCHAR(100),
			col2 INT,
			KEY idx_mixed (col1(10), col2),
			PRIMARY KEY (id)
		)`)

		ddl := c.ShowCreateTable("testdb", "t1")
		if !strings.Contains(ddl, "`col1`(10)") {
			t.Errorf("expected col1(10) in SHOW CREATE TABLE, got:\n%s", ddl)
		}
		// col2 should appear without prefix length
		if !strings.Contains(ddl, "`col2`") {
			t.Errorf("expected col2 in SHOW CREATE TABLE, got:\n%s", ddl)
		}
		// Ensure col2 does NOT have a prefix length
		if strings.Contains(ddl, "`col2`(") {
			t.Errorf("col2 should not have prefix length, got:\n%s", ddl)
		}
	})

	// Scenario 4: KEY idx ((UPPER(col))) — expression index with function

	// Scenario 5: KEY idx ((col1 + col2)) — expression index with arithmetic

	// Scenario 6: UNIQUE KEY idx ((UPPER(col))) — unique expression index

	// Scenario 7: KEY idx (col1, (UPPER(col2))) — mixed regular and expression columns

	// Scenario 8: KEY idx (col DESC) — descending index column
	t.Run("descending_index_column", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t1 (
			id INT NOT NULL,
			name VARCHAR(100),
			KEY idx_desc (name DESC),
			PRIMARY KEY (id)
		)`)

		ddl := c.ShowCreateTable("testdb", "t1")
		if !strings.Contains(ddl, "DESC") {
			t.Errorf("expected DESC in index column rendering, got:\n%s", ddl)
		}
		if !strings.Contains(ddl, "KEY `idx_desc`") {
			t.Errorf("expected KEY idx_desc in SHOW CREATE TABLE, got:\n%s", ddl)
		}
	})
}
