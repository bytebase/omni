package catalog

import (
	"strings"
	"testing"
)

func TestShowCreateTableBasic(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE t1 (
		id INT NOT NULL AUTO_INCREMENT,
		name VARCHAR(100) DEFAULT 'test',
		PRIMARY KEY (id),
		KEY idx_name (name)
	)`)

	got := c.ShowCreateTable("testdb", "t1")
	assertContains(t, got, "CREATE TABLE `t1`")
	assertContains(t, got, "`id` int NOT NULL AUTO_INCREMENT")
	assertContains(t, got, "`name` varchar(100) DEFAULT 'test'")
	assertContains(t, got, "PRIMARY KEY (`id`)")
	assertContains(t, got, "KEY `idx_name` (`name`)")
	assertContains(t, got, "ENGINE=InnoDB")
	assertContains(t, got, "DEFAULT CHARSET=utf8mb4")
	// MySQL 8.0 always shows COLLATE in SHOW CREATE TABLE.
	assertContains(t, got, "COLLATE=utf8mb4_0900_ai_ci")
}

func TestShowCreateTableUniqueKey(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE t2 (
		id INT NOT NULL AUTO_INCREMENT,
		email VARCHAR(255) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY uk_email (email)
	)`)

	got := c.ShowCreateTable("testdb", "t2")
	assertContains(t, got, "UNIQUE KEY `uk_email` (`email`)")
}

func TestShowCreateTableDefaults(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE t3 (
		id INT NOT NULL AUTO_INCREMENT,
		nullable_col VARCHAR(50),
		str_default VARCHAR(50) DEFAULT 'hello',
		num_default INT DEFAULT 42,
		PRIMARY KEY (id)
	)`)

	got := c.ShowCreateTable("testdb", "t3")
	// Nullable column without explicit default should show DEFAULT NULL.
	assertContains(t, got, "`nullable_col` varchar(50) DEFAULT NULL")
	// String default.
	assertContains(t, got, "`str_default` varchar(50) DEFAULT 'hello'")
	// Numeric default — MySQL 8.0 quotes it.
	assertContains(t, got, "`num_default` int DEFAULT '42'")
}

func TestShowCreateTableMultipleIndexes(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE t4 (
		id INT NOT NULL AUTO_INCREMENT,
		a INT,
		b VARCHAR(100),
		c INT,
		PRIMARY KEY (id),
		KEY idx_a (a),
		KEY idx_b (b),
		UNIQUE KEY uk_c (c)
	)`)

	got := c.ShowCreateTable("testdb", "t4")
	assertContains(t, got, "PRIMARY KEY (`id`)")
	assertContains(t, got, "KEY `idx_a` (`a`)")
	assertContains(t, got, "KEY `idx_b` (`b`)")
	assertContains(t, got, "UNIQUE KEY `uk_c` (`c`)")
}

func TestShowCreateTableForeignKey(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE parent (
		id INT NOT NULL AUTO_INCREMENT,
		PRIMARY KEY (id)
	)`)
	mustExec(t, c, `CREATE TABLE child (
		id INT NOT NULL AUTO_INCREMENT,
		parent_id INT NOT NULL,
		PRIMARY KEY (id),
		CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent (id) ON DELETE CASCADE
	)`)

	got := c.ShowCreateTable("testdb", "child")
	assertContains(t, got, "CONSTRAINT `fk_parent` FOREIGN KEY (`parent_id`) REFERENCES `parent` (`id`)")
	assertContains(t, got, "ON DELETE CASCADE")
	// ON UPDATE RESTRICT is the default and should NOT appear.
	assertNotContains(t, got, "ON UPDATE")
}

func TestShowCreateTableForeignKeyRestrict(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE ref_tbl (
		id INT NOT NULL AUTO_INCREMENT,
		PRIMARY KEY (id)
	)`)
	mustExec(t, c, `CREATE TABLE fk_tbl (
		id INT NOT NULL AUTO_INCREMENT,
		ref_id INT NOT NULL,
		PRIMARY KEY (id),
		CONSTRAINT fk_ref FOREIGN KEY (ref_id) REFERENCES ref_tbl (id) ON DELETE RESTRICT ON UPDATE RESTRICT
	)`)

	got := c.ShowCreateTable("testdb", "fk_tbl")
	// MySQL 8.0 shows RESTRICT when explicitly specified (unlike NO ACTION which is hidden).
	assertContains(t, got, "ON DELETE RESTRICT")
	assertContains(t, got, "ON UPDATE RESTRICT")
}

func TestShowCreateTableComment(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE t5 (
		id INT NOT NULL AUTO_INCREMENT,
		name VARCHAR(100) COMMENT 'user name',
		PRIMARY KEY (id)
	) COMMENT='main table'`)

	got := c.ShowCreateTable("testdb", "t5")
	assertContains(t, got, "COMMENT 'user name'")
	assertContains(t, got, "COMMENT='main table'")
}

func TestShowCreateTableUnknownDatabaseOrTable(t *testing.T) {
	c := setupWithDB(t)

	if got := c.ShowCreateTable("nonexistent", "t1"); got != "" {
		t.Errorf("expected empty string for unknown database, got %q", got)
	}

	if got := c.ShowCreateTable("testdb", "nonexistent"); got != "" {
		t.Errorf("expected empty string for unknown table, got %q", got)
	}
}

func TestShowCreateTableNotNullNoDefault(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE t6 (
		id INT NOT NULL,
		name VARCHAR(100) NOT NULL,
		PRIMARY KEY (id)
	)`)

	got := c.ShowCreateTable("testdb", "t6")
	// NOT NULL columns without a default should NOT show DEFAULT NULL.
	assertContains(t, got, "`id` int NOT NULL")
	assertContains(t, got, "`name` varchar(100) NOT NULL")
	assertNotContains(t, got, "DEFAULT NULL")
}

func TestShowCreateTableNonDefaultCollation(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE t7 (
		id INT NOT NULL,
		PRIMARY KEY (id)
	) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`)

	got := c.ShowCreateTable("testdb", "t7")
	assertContains(t, got, "DEFAULT CHARSET=utf8mb4")
	assertContains(t, got, "COLLATE=utf8mb4_unicode_ci")
}

func TestShowCreateTableAutoIncrementColumn(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE t8 (
		id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
		PRIMARY KEY (id)
	)`)

	got := c.ShowCreateTable("testdb", "t8")
	assertContains(t, got, "`id` bigint unsigned NOT NULL AUTO_INCREMENT")
}

func assertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("expected output to contain %q\ngot:\n%s", substr, s)
	}
}

func assertNotContains(t *testing.T, s, substr string) {
	t.Helper()
	if strings.Contains(s, substr) {
		t.Errorf("expected output NOT to contain %q\ngot:\n%s", substr, s)
	}
}
