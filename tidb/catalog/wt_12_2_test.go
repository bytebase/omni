package catalog

import (
	"strings"
	"testing"
)

// --- Section 8.2 (starmap): Index Rendering in SHOW CREATE TABLE (7 scenarios) ---

func TestWalkThrough_12_2(t *testing.T) {
	t.Run("index_ordering_pk_unique_key_fulltext_spatial", func(t *testing.T) {
		// Scenario 1: Table with PK + UNIQUE + KEY + FULLTEXT + SPATIAL — verify ordering
		// MySQL 8.0 orders: PRIMARY → UNIQUE → regular+SPATIAL → expression → FULLTEXT
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t (
			id INT NOT NULL,
			name VARCHAR(100),
			body TEXT,
			geo GEOMETRY NOT NULL SRID 0,
			PRIMARY KEY (id),
			KEY idx_name (name),
			UNIQUE KEY uk_name (name),
			FULLTEXT KEY ft_body (body),
			SPATIAL KEY sp_geo (geo)
		)`)

		got := c.ShowCreateTable("testdb", "t")

		// Verify all index types are present.
		assertContains(t, got, "PRIMARY KEY (`id`)")
		assertContains(t, got, "UNIQUE KEY `uk_name` (`name`)")
		assertContains(t, got, "KEY `idx_name` (`name`)")
		assertContains(t, got, "FULLTEXT KEY `ft_body` (`body`)")
		assertContains(t, got, "SPATIAL KEY `sp_geo` (`geo`)")

		// Verify ordering: PRIMARY < UNIQUE < KEY/SPATIAL < FULLTEXT
		posPK := strings.Index(got, "PRIMARY KEY")
		posUK := strings.Index(got, "UNIQUE KEY")
		posKey := strings.Index(got, "KEY `idx_name`")
		posSP := strings.Index(got, "SPATIAL KEY")
		posFT := strings.Index(got, "FULLTEXT KEY")

		if posPK >= posUK {
			t.Error("PRIMARY KEY should appear before UNIQUE KEY")
		}
		if posUK >= posKey {
			t.Error("UNIQUE KEY should appear before regular KEY")
		}
		// SPATIAL is in the same group as regular keys, so it should come after or alongside idx_name
		if posSP <= posUK {
			t.Error("SPATIAL KEY should appear after UNIQUE KEY")
		}
		if posFT <= posKey {
			t.Error("FULLTEXT KEY should appear after regular KEY")
		}
		if posFT <= posSP {
			t.Error("FULLTEXT KEY should appear after SPATIAL KEY")
		}
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

		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "KEY `idx_name` (`name`) USING BTREE")
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

		got := c.ShowCreateTable("testdb", "t")
		assertContains(t, got, "KEY `idx_name` (`name`) USING HASH")
	})
}
