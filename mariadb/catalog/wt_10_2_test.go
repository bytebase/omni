package catalog

import (
	"strings"
	"testing"
)

// --- Section 6.2 (Phase 6): LIKE Edge Cases (7 scenarios) ---
// File target: wt_10_2_test.go
// Proof: go test ./mysql/catalog/ -short -count=1 -run "TestWalkThrough_10_2"

func TestWalkThrough_10_2_LIKEEdgeCases(t *testing.T) {
	// Scenario 1: LIKE copies generated columns — expression and VIRTUAL/STORED preserved
	t.Run("like_copies_generated_columns", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			qty INT NOT NULL,
			total DECIMAL(10,2) AS (price * qty) STORED,
			label VARCHAR(100) AS (CONCAT('item-', id)) VIRTUAL,
			PRIMARY KEY (id)
		)`)
		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}

		// Check STORED generated column
		totalCol := dstTbl.GetColumn("total")
		if totalCol == nil {
			t.Fatal("column total not found in dst")
		}
		if totalCol.Generated == nil {
			t.Fatal("total should be a generated column")
		}
		if !totalCol.Generated.Stored {
			t.Error("total should be STORED")
		}
		if totalCol.Generated.Expr == "" {
			t.Error("total generated expression should not be empty")
		}

		// Check VIRTUAL generated column
		labelCol := dstTbl.GetColumn("label")
		if labelCol == nil {
			t.Fatal("column label not found in dst")
		}
		if labelCol.Generated == nil {
			t.Fatal("label should be a generated column")
		}
		if labelCol.Generated.Stored {
			t.Error("label should be VIRTUAL (not stored)")
		}
		if labelCol.Generated.Expr == "" {
			t.Error("label generated expression should not be empty")
		}

		// Verify expressions match source
		srcTbl := c.GetDatabase("testdb").GetTable("src")
		srcTotal := srcTbl.GetColumn("total")
		srcLabel := srcTbl.GetColumn("label")
		if totalCol.Generated.Expr != srcTotal.Generated.Expr {
			t.Errorf("total expression mismatch: src=%q dst=%q", srcTotal.Generated.Expr, totalCol.Generated.Expr)
		}
		if labelCol.Generated.Expr != srcLabel.Generated.Expr {
			t.Errorf("label expression mismatch: src=%q dst=%q", srcLabel.Generated.Expr, labelCol.Generated.Expr)
		}
	})

	// Scenario 2: LIKE copies INVISIBLE columns
	t.Run("like_copies_invisible_columns", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			visible_col VARCHAR(100),
			hidden_col INT INVISIBLE,
			PRIMARY KEY (id)
		)`)
		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}

		hiddenCol := dstTbl.GetColumn("hidden_col")
		if hiddenCol == nil {
			t.Fatal("column hidden_col not found in dst")
		}
		if !hiddenCol.Invisible {
			t.Error("hidden_col should be INVISIBLE in dst")
		}

		visibleCol := dstTbl.GetColumn("visible_col")
		if visibleCol == nil {
			t.Fatal("column visible_col not found in dst")
		}
		if visibleCol.Invisible {
			t.Error("visible_col should NOT be invisible in dst")
		}

		// Verify SHOW CREATE TABLE renders INVISIBLE
		sct := c.ShowCreateTable("testdb", "dst")
		if !strings.Contains(sct, "INVISIBLE") {
			t.Errorf("SHOW CREATE TABLE should contain INVISIBLE keyword, got:\n%s", sct)
		}
	})

	// Scenario 3: LIKE does NOT copy partitioning — target table is unpartitioned
	t.Run("like_does_not_copy_partitioning", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			created_at DATE NOT NULL,
			PRIMARY KEY (id, created_at)
		) PARTITION BY RANGE (YEAR(created_at)) (
			PARTITION p2020 VALUES LESS THAN (2021),
			PARTITION p2021 VALUES LESS THAN (2022),
			PARTITION pmax VALUES LESS THAN MAXVALUE
		)`)
		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		// Source should have partitioning
		srcTbl := c.GetDatabase("testdb").GetTable("src")
		if srcTbl.Partitioning == nil {
			t.Fatal("source table should have partitioning")
		}

		// Destination should NOT have partitioning (MySQL behavior)
		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}
		if dstTbl.Partitioning != nil {
			t.Error("LIKE should NOT copy partitioning — target table should be unpartitioned")
		}

		// SHOW CREATE TABLE for dst should not mention PARTITION
		sct := c.ShowCreateTable("testdb", "dst")
		if strings.Contains(sct, "PARTITION") {
			t.Errorf("SHOW CREATE TABLE for dst should not contain PARTITION, got:\n%s", sct)
		}
	})

	// Scenario 4: LIKE from table with prefix index — prefix length preserved
	t.Run("like_copies_prefix_index", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			name VARCHAR(255),
			bio TEXT,
			PRIMARY KEY (id),
			INDEX idx_name (name(50)),
			INDEX idx_bio (bio(100))
		)`)
		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}

		// Find idx_name and verify prefix length
		var idxName, idxBio *Index
		for _, idx := range dstTbl.Indexes {
			switch idx.Name {
			case "idx_name":
				idxName = idx
			case "idx_bio":
				idxBio = idx
			}
		}

		if idxName == nil {
			t.Fatal("index idx_name not found in dst")
		}
		if len(idxName.Columns) != 1 || idxName.Columns[0].Length != 50 {
			t.Errorf("idx_name prefix length: expected 50, got %d", idxName.Columns[0].Length)
		}

		if idxBio == nil {
			t.Fatal("index idx_bio not found in dst")
		}
		if len(idxBio.Columns) != 1 || idxBio.Columns[0].Length != 100 {
			t.Errorf("idx_bio prefix length: expected 100, got %d", idxBio.Columns[0].Length)
		}

		// Verify SHOW CREATE TABLE renders prefix lengths
		sct := c.ShowCreateTable("testdb", "dst")
		if !strings.Contains(sct, "`name`(50)") {
			t.Errorf("SHOW CREATE TABLE should contain name(50), got:\n%s", sct)
		}
		if !strings.Contains(sct, "`bio`(100)") {
			t.Errorf("SHOW CREATE TABLE should contain bio(100), got:\n%s", sct)
		}
	})

	// Scenario 5: LIKE into TEMPORARY TABLE
	t.Run("like_into_temporary_table", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			name VARCHAR(100),
			PRIMARY KEY (id)
		)`)
		wtExec(t, c, `CREATE TEMPORARY TABLE dst LIKE src`)

		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}
		if !dstTbl.Temporary {
			t.Error("dst should be a TEMPORARY table")
		}

		// Columns should still be copied
		if len(dstTbl.Columns) != 2 {
			t.Errorf("expected 2 columns, got %d", len(dstTbl.Columns))
		}

		// Source should NOT be temporary
		srcTbl := c.GetDatabase("testdb").GetTable("src")
		if srcTbl.Temporary {
			t.Error("src should not be temporary")
		}
	})

	// Scenario 6: LIKE cross-database — source in different database
	t.Run("like_cross_database", func(t *testing.T) {
		c := wtSetup(t)
		// Create source table in a different database
		wtExec(t, c, "CREATE DATABASE other_db")
		wtExec(t, c, `CREATE TABLE other_db.src (
			id INT NOT NULL,
			name VARCHAR(100) NOT NULL,
			score INT DEFAULT 0,
			PRIMARY KEY (id),
			INDEX idx_name (name)
		)`)

		// Create LIKE table in testdb referencing other_db.src
		wtExec(t, c, `CREATE TABLE dst LIKE other_db.src`)

		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found in testdb")
		}

		// Verify columns were copied
		if len(dstTbl.Columns) != 3 {
			t.Fatalf("expected 3 columns, got %d", len(dstTbl.Columns))
		}
		nameCol := dstTbl.GetColumn("name")
		if nameCol == nil {
			t.Fatal("column name not found in dst")
		}
		if nameCol.Nullable {
			t.Error("name should be NOT NULL")
		}

		scoreCol := dstTbl.GetColumn("score")
		if scoreCol == nil {
			t.Fatal("column score not found in dst")
		}
		if scoreCol.Default == nil || *scoreCol.Default != "0" {
			def := "<nil>"
			if scoreCol.Default != nil {
				def = *scoreCol.Default
			}
			t.Errorf("score default: expected '0', got %q", def)
		}

		// Verify indexes were copied
		var idxName *Index
		for _, idx := range dstTbl.Indexes {
			if idx.Name == "idx_name" {
				idxName = idx
				break
			}
		}
		if idxName == nil {
			t.Fatal("index idx_name not found in dst")
		}

		// dst should belong to testdb, not other_db
		if dstTbl.Database.Name != "testdb" {
			t.Errorf("dst should belong to testdb, got %q", dstTbl.Database.Name)
		}
	})

	// Scenario 7: LIKE then ALTER TABLE ADD COLUMN — verify table is independently modifiable
	t.Run("like_then_alter_add_column", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			name VARCHAR(100),
			PRIMARY KEY (id)
		)`)
		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		// Add a column to dst — should not affect src
		wtExec(t, c, `ALTER TABLE dst ADD COLUMN email VARCHAR(255) NOT NULL`)

		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}
		if len(dstTbl.Columns) != 3 {
			t.Fatalf("dst: expected 3 columns after ADD COLUMN, got %d", len(dstTbl.Columns))
		}
		emailCol := dstTbl.GetColumn("email")
		if emailCol == nil {
			t.Fatal("column email not found in dst")
		}

		// Source should be unaffected
		srcTbl := c.GetDatabase("testdb").GetTable("src")
		if len(srcTbl.Columns) != 2 {
			t.Fatalf("src: expected 2 columns (unchanged), got %d", len(srcTbl.Columns))
		}
		if srcTbl.GetColumn("email") != nil {
			t.Error("src should NOT have email column — tables should be independent")
		}
	})
}
