package catalog

import (
	"strings"
	"testing"
)

// --- Section 6.1 (Phase 6): Basic LIKE Completeness (9 scenarios) ---
// File target: wt_10_1_test.go
// Proof: go test ./mysql/catalog/ -short -count=1 -run "TestWalkThrough_10_1"

func TestWalkThrough_10_1_BasicLIKECompleteness(t *testing.T) {
	// Scenario 1: LIKE copies all column definitions (name, type, nullability, default, comment)
	t.Run("like_copies_column_definitions", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			name VARCHAR(100) NOT NULL DEFAULT 'unknown' COMMENT 'user name',
			email VARCHAR(255) DEFAULT NULL,
			score DECIMAL(10,2) NOT NULL DEFAULT 0.00,
			PRIMARY KEY (id)
		)`)
		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		srcTbl := c.GetDatabase("testdb").GetTable("src")
		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}
		if len(dstTbl.Columns) != len(srcTbl.Columns) {
			t.Fatalf("expected %d columns, got %d", len(srcTbl.Columns), len(dstTbl.Columns))
		}

		// Verify each column attribute
		for i, srcCol := range srcTbl.Columns {
			dstCol := dstTbl.Columns[i]
			if dstCol.Name != srcCol.Name {
				t.Errorf("col %d: name mismatch: %q vs %q", i, srcCol.Name, dstCol.Name)
			}
			if dstCol.ColumnType != srcCol.ColumnType {
				t.Errorf("col %q: type mismatch: %q vs %q", srcCol.Name, srcCol.ColumnType, dstCol.ColumnType)
			}
			if dstCol.Nullable != srcCol.Nullable {
				t.Errorf("col %q: nullable mismatch: %v vs %v", srcCol.Name, srcCol.Nullable, dstCol.Nullable)
			}
			if dstCol.Comment != srcCol.Comment {
				t.Errorf("col %q: comment mismatch: %q vs %q", srcCol.Name, srcCol.Comment, dstCol.Comment)
			}
			// Compare defaults
			if (srcCol.Default == nil) != (dstCol.Default == nil) {
				t.Errorf("col %q: default nil mismatch", srcCol.Name)
			} else if srcCol.Default != nil && *srcCol.Default != *dstCol.Default {
				t.Errorf("col %q: default value mismatch: %q vs %q", srcCol.Name, *srcCol.Default, *dstCol.Default)
			}
		}

		// Verify SHOW CREATE TABLE DDLs are structurally similar
		srcDDL := c.ShowCreateTable("testdb", "src")
		dstDDL := c.ShowCreateTable("testdb", "dst")
		// Replace table names for comparison
		srcNorm := strings.Replace(srcDDL, "`src`", "`dst`", 1)
		if srcNorm != dstDDL {
			t.Errorf("SHOW CREATE TABLE mismatch:\nsrc:\n%s\ndst:\n%s", srcDDL, dstDDL)
		}
	})

	// Scenario 2: LIKE copies PRIMARY KEY
	t.Run("like_copies_primary_key", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			name VARCHAR(50) NOT NULL,
			PRIMARY KEY (id)
		)`)
		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}

		// Check PRIMARY index exists
		var hasPK bool
		for _, idx := range dstTbl.Indexes {
			if idx.Primary {
				hasPK = true
				if len(idx.Columns) != 1 || idx.Columns[0].Name != "id" {
					t.Errorf("expected PK on (id), got %v", idx.Columns)
				}
				break
			}
		}
		if !hasPK {
			t.Error("expected PRIMARY KEY on dst table")
		}

		ddl := c.ShowCreateTable("testdb", "dst")
		if !strings.Contains(ddl, "PRIMARY KEY") {
			t.Errorf("expected PRIMARY KEY in DDL:\n%s", ddl)
		}
	})

	// Scenario 3: LIKE copies UNIQUE KEYs
	t.Run("like_copies_unique_keys", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			email VARCHAR(255) NOT NULL,
			code VARCHAR(20) NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY uk_email (email),
			UNIQUE KEY uk_code (code)
		)`)
		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}

		uniqueCount := 0
		for _, idx := range dstTbl.Indexes {
			if idx.Unique && !idx.Primary {
				uniqueCount++
			}
		}
		if uniqueCount != 2 {
			t.Errorf("expected 2 unique keys, got %d", uniqueCount)
		}

		ddl := c.ShowCreateTable("testdb", "dst")
		if !strings.Contains(ddl, "UNIQUE KEY `uk_email`") {
			t.Errorf("expected UNIQUE KEY uk_email in DDL:\n%s", ddl)
		}
		if !strings.Contains(ddl, "UNIQUE KEY `uk_code`") {
			t.Errorf("expected UNIQUE KEY uk_code in DDL:\n%s", ddl)
		}
	})

	// Scenario 4: LIKE copies regular indexes
	t.Run("like_copies_regular_indexes", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			name VARCHAR(100),
			age INT,
			PRIMARY KEY (id),
			KEY idx_name (name),
			KEY idx_age (age)
		)`)
		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}

		regularCount := 0
		for _, idx := range dstTbl.Indexes {
			if !idx.Primary && !idx.Unique && !idx.Fulltext && !idx.Spatial {
				regularCount++
			}
		}
		if regularCount != 2 {
			t.Errorf("expected 2 regular indexes, got %d", regularCount)
		}

		ddl := c.ShowCreateTable("testdb", "dst")
		if !strings.Contains(ddl, "KEY `idx_name`") {
			t.Errorf("expected KEY idx_name in DDL:\n%s", ddl)
		}
		if !strings.Contains(ddl, "KEY `idx_age`") {
			t.Errorf("expected KEY idx_age in DDL:\n%s", ddl)
		}
	})

	// Scenario 5: LIKE copies FULLTEXT indexes
	t.Run("like_copies_fulltext_indexes", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			content TEXT,
			PRIMARY KEY (id),
			FULLTEXT KEY ft_content (content)
		)`)
		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}

		var hasFT bool
		for _, idx := range dstTbl.Indexes {
			if idx.Fulltext {
				hasFT = true
				if idx.Name != "ft_content" {
					t.Errorf("expected fulltext index named ft_content, got %q", idx.Name)
				}
				break
			}
		}
		if !hasFT {
			t.Error("expected FULLTEXT index on dst table")
		}

		ddl := c.ShowCreateTable("testdb", "dst")
		if !strings.Contains(ddl, "FULLTEXT KEY `ft_content`") {
			t.Errorf("expected FULLTEXT KEY ft_content in DDL:\n%s", ddl)
		}
	})

	// Scenario 6: LIKE copies CHECK constraints
	t.Run("like_copies_check_constraints", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			age INT,
			PRIMARY KEY (id),
			CONSTRAINT chk_age CHECK (age >= 0)
		)`)
		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}

		var hasCheck bool
		for _, con := range dstTbl.Constraints {
			if con.Type == ConCheck {
				hasCheck = true
				if con.Name != "chk_age" {
					t.Errorf("expected check constraint named chk_age, got %q", con.Name)
				}
				break
			}
		}
		if !hasCheck {
			t.Error("expected CHECK constraint on dst table")
		}

		ddl := c.ShowCreateTable("testdb", "dst")
		if !strings.Contains(ddl, "CONSTRAINT `chk_age` CHECK") {
			t.Errorf("expected CHECK constraint in DDL:\n%s", ddl)
		}
	})

	// Scenario 7: LIKE does NOT copy FOREIGN KEY constraints — MySQL 8.0 behavior
	t.Run("like_does_not_copy_foreign_keys", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE parent (
			id INT NOT NULL,
			PRIMARY KEY (id)
		)`)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			parent_id INT NOT NULL,
			PRIMARY KEY (id),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id)
		)`)
		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}

		// No FK constraints should be copied
		for _, con := range dstTbl.Constraints {
			if con.Type == ConForeignKey {
				t.Errorf("LIKE should NOT copy FK constraints, but found FK %q", con.Name)
			}
		}

		// The FK-backing index IS still copied (it's a regular KEY on the source)
		ddl := c.ShowCreateTable("testdb", "dst")
		if strings.Contains(ddl, "FOREIGN KEY") {
			t.Errorf("expected no FOREIGN KEY in dst DDL:\n%s", ddl)
		}
		// The backing index for the FK should still be present
		if !strings.Contains(ddl, "KEY `fk_parent`") {
			t.Errorf("expected FK-backing index to be copied:\n%s", ddl)
		}
	})

	// Scenario 8: LIKE copies AUTO_INCREMENT column attribute — but counter resets to 0
	t.Run("like_copies_auto_increment_attribute_resets_counter", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL AUTO_INCREMENT,
			name VARCHAR(100),
			PRIMARY KEY (id)
		)`)
		// Insert to advance the counter on src
		wtExec(t, c, `INSERT INTO src (name) VALUES ('a'), ('b'), ('c')`)

		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}

		// AUTO_INCREMENT attribute should be copied
		idCol := dstTbl.GetColumn("id")
		if idCol == nil {
			t.Fatal("column id not found on dst")
		}
		if !idCol.AutoIncrement {
			t.Error("expected AUTO_INCREMENT attribute on dst.id")
		}

		// Counter should reset — dst DDL should NOT show AUTO_INCREMENT=N (or show AUTO_INCREMENT=1 at most)
		dstDDL := c.ShowCreateTable("testdb", "dst")
		if strings.Contains(dstDDL, "AUTO_INCREMENT=") {
			// MySQL 8.0 does not show AUTO_INCREMENT=N on empty tables
			t.Errorf("expected no AUTO_INCREMENT=N on empty dst table:\n%s", dstDDL)
		}

		// Verify the column-level auto_increment keyword is present
		if !strings.Contains(dstDDL, "AUTO_INCREMENT") {
			t.Errorf("expected AUTO_INCREMENT column attribute in DDL:\n%s", dstDDL)
		}
	})

	// Scenario 9: LIKE copies ENGINE, CHARSET, COLLATION, COMMENT
	t.Run("like_copies_table_options", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE src (
			id INT NOT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='source table'`)
		wtExec(t, c, `CREATE TABLE dst LIKE src`)

		srcTbl := c.GetDatabase("testdb").GetTable("src")
		dstTbl := c.GetDatabase("testdb").GetTable("dst")
		if dstTbl == nil {
			t.Fatal("table dst not found")
		}

		if dstTbl.Engine != srcTbl.Engine {
			t.Errorf("ENGINE mismatch: %q vs %q", srcTbl.Engine, dstTbl.Engine)
		}
		if dstTbl.Charset != srcTbl.Charset {
			t.Errorf("CHARSET mismatch: %q vs %q", srcTbl.Charset, dstTbl.Charset)
		}
		if dstTbl.Collation != srcTbl.Collation {
			t.Errorf("COLLATION mismatch: %q vs %q", srcTbl.Collation, dstTbl.Collation)
		}
		if dstTbl.Comment != srcTbl.Comment {
			t.Errorf("COMMENT mismatch: %q vs %q", srcTbl.Comment, dstTbl.Comment)
		}

		ddl := c.ShowCreateTable("testdb", "dst")
		if !strings.Contains(ddl, "ENGINE=InnoDB") {
			t.Errorf("expected ENGINE=InnoDB in DDL:\n%s", ddl)
		}
		if !strings.Contains(ddl, "COMMENT='source table'") {
			t.Errorf("expected COMMENT='source table' in DDL:\n%s", ddl)
		}
	})
}
