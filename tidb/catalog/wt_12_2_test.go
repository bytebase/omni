package catalog

import (
	"strings"
	"testing"
)

// --- Section 8.2 (starmap): Index Rendering in SHOW CREATE TABLE (7 scenarios) ---

func TestWalkThrough_12_2(t *testing.T) {
	t.Run("index_ordering_pk_unique_key_fulltext", func(t *testing.T) {
		// Scenario 1: Table with PK + UNIQUE + KEY + FULLTEXT — verify ordering
		// PRIMARY → UNIQUE → regular → expression → FULLTEXT. TiDB has no
		// spatial column type or SPATIAL KEY in a table definition (1064).
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			name VARCHAR(100),
			body TEXT,
			PRIMARY KEY (id),
			KEY idx_name (name),
			UNIQUE KEY uk_name (name),
			FULLTEXT KEY ft_body (body)
		)`)

		// TiDB writes the keys in the order the table got them, with an
		// integer clustered primary key first, and keeps no FULLTEXT key
		// (v8.5.5, read on 2026-11-06).
		assertShowCreate(t, c, "t", "CREATE TABLE `t` (\n"+
			"  `id` int NOT NULL,\n"+
			"  `name` varchar(100) DEFAULT NULL,\n"+
			"  `body` text DEFAULT NULL,\n"+
			"  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,\n"+
			"  KEY `idx_name` (`name`),\n"+
			"  UNIQUE KEY `uk_name` (`name`)\n"+
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin")
	})

	t.Run("regular_and_expression_index_ordering", func(t *testing.T) {
		// Scenario 2: Table with regular index + expression index — verify relative ordering
		// MySQL 8.0: regular keys come before expression-based keys
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			name VARCHAR(100),
			age INT,
			PRIMARY KEY (id),
			KEY idx_name (name),
			KEY idx_expr ((age * 2))
		)`)

		got := c.ShowCreateTable("testdb", "t")

		assertContains(t, got, "KEY `idx_name` (`name`)")
		assertContains(t, got, "KEY `idx_expr` ((")

		posRegular := strings.Index(got, "KEY `idx_name`")
		posExpr := strings.Index(got, "KEY `idx_expr`")

		if posRegular >= posExpr {
			t.Errorf("regular KEY should appear before expression KEY\ngot:\n%s", got)
		}
	})

	t.Run("index_with_comment", func(t *testing.T) {
		// Scenario 3: Index with COMMENT
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			name VARCHAR(100),
			PRIMARY KEY (id),
			KEY idx_name (name) COMMENT 'name lookup'
		)`)

		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "KEY `idx_name` (`name`) COMMENT 'name lookup'")
	})

	t.Run("index_with_invisible", func(t *testing.T) {
		// Scenario 4: Index with INVISIBLE
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			name VARCHAR(100),
			PRIMARY KEY (id),
			KEY idx_name (name) INVISIBLE
		)`)

		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "KEY `idx_name` (`name`) /*!80000 INVISIBLE */")
	})

	t.Run("index_with_key_block_size", func(t *testing.T) {
		// Scenario 5: Index with KEY_BLOCK_SIZE
		// MySQL 8.0 parses index-level KEY_BLOCK_SIZE but does not render
		// it in SHOW CREATE TABLE output. Verify the catalog matches this
		// behavior by NOT including KEY_BLOCK_SIZE on the index line.
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			name VARCHAR(100),
			PRIMARY KEY (id),
			KEY idx_name (name) KEY_BLOCK_SIZE=4
		)`)

		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "KEY `idx_name` (`name`)")
		if strings.Contains(got, "KEY_BLOCK_SIZE") {
			t.Errorf("expected SHOW CREATE TABLE to omit index-level KEY_BLOCK_SIZE (MySQL 8.0 behavior), got:\n%s", got)
		}
	})

	t.Run("index_with_using_btree", func(t *testing.T) {
		// Scenario 6: Index with USING BTREE
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			name VARCHAR(100),
			PRIMARY KEY (id),
			KEY idx_name (name) USING BTREE
		)`)

		// TiDB keeps no USING clause in SHOW CREATE TABLE.
		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "KEY `idx_name` (`name`)\n")
		assertNotContains(t, got, "USING")
	})

	t.Run("index_with_using_hash", func(t *testing.T) {
		// Scenario 7: Index with USING HASH
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			name VARCHAR(100),
			PRIMARY KEY (id),
			KEY idx_name (name) USING HASH
		) ENGINE=MEMORY`)

		// TiDB keeps no USING clause and no engine but InnoDB.
		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "KEY `idx_name` (`name`)\n")
		assertContains(t, got, ") ENGINE=InnoDB ")
		assertNotContains(t, got, "USING")
	})
}
