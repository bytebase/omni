package catalog

import (
	"context"
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/internal/testcontainer"
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
	assertContains(t, got, "`id` int(11) NOT NULL AUTO_INCREMENT")
	assertContains(t, got, "`name` varchar(100) DEFAULT 'test'")
	assertContains(t, got, "PRIMARY KEY (`id`)")
	assertContains(t, got, "KEY `idx_name` (`name`)")
	assertContains(t, got, "ENGINE=InnoDB")
	assertContains(t, got, "DEFAULT CHARSET=utf8mb4")
	// MySQL 8.0 always shows COLLATE in SHOW CREATE TABLE.
	assertContains(t, got, "COLLATE=utf8mb4_uca1400_ai_ci")
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
	assertContains(t, got, "`num_default` int(11) DEFAULT 42")
}

func TestShowCreateTableMySQLDefaultNormalization(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE t_defaults_mysql (
		tiny TINYINT(1) DEFAULT 1,
		bool_true BOOLEAN DEFAULT TRUE,
		bool_false BOOLEAN DEFAULT FALSE,
		bits BIT(8) DEFAULT b'00001111',
		ts3 TIMESTAMP(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
		dt6 DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
		expr_int INT DEFAULT (FLOOR(RAND()*100)),
		expr_json JSON DEFAULT (JSON_ARRAY()),
		expr_varchar VARCHAR(36) DEFAULT (UUID())
	)`)

	got := c.ShowCreateTable("testdb", "t_defaults_mysql")
	assertContains(t, got, "`tiny` tinyint(1) DEFAULT 1")
	assertContains(t, got, "`bool_true` tinyint(1) DEFAULT 1")
	assertContains(t, got, "`bool_false` tinyint(1) DEFAULT 0")
	assertContains(t, got, "`bits` bit(8) DEFAULT b'1111'")
	assertContains(t, got, "`ts3` timestamp(3) NULL DEFAULT current_timestamp(3)")
	assertContains(t, got, "`dt6` datetime(6) DEFAULT current_timestamp(6) ON UPDATE current_timestamp(6)")
	assertContains(t, got, "`expr_int` int(11) DEFAULT floor(rand() * 100)")
	assertContains(t, got, "`expr_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT json_array() CHECK (json_valid(`expr_json`))")
	assertContains(t, got, "`expr_varchar` varchar(36) DEFAULT uuid()")
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
	// MariaDB writes no referential action for RESTRICT, which is the default,
	// and writes NO ACTION where that was asked for. MySQL 8.0 does the
	// opposite. Measured on mariadb:11.8.8.
	assertNotContains(t, got, "ON DELETE RESTRICT")
	assertNotContains(t, got, "ON UPDATE RESTRICT")
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
	assertContains(t, got, "`id` int(11) NOT NULL")
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
	assertContains(t, got, "`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT")
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

// TestOracle_ShowCreateMatchesServer runs each script on the pinned server and
// on the catalog, and compares what SHOW CREATE writes for every table and view
// it leaves.
func TestOracle_ShowCreateMatchesServer(t *testing.T) {
	scripts := []struct{ name, sql string }{
		{"function defaults", "CREATE TABLE t (d DATE DEFAULT curdate(), u VARCHAR(20) DEFAULT concat('a', 'b'), z DATE DEFAULT DATE '2020-01-02', x INT DEFAULT (1 + 1))"},
		{"column defaults", "CREATE TABLE t (a INT, g INT DEFAULT (a), k INT DEFAULT ((a)), b INT DEFAULT (a + 1), c INT DEFAULT (-1), d INT DEFAULT ((1 + 1)))"},
		{"wide integers", "CREATE TABLE t (a BIGINT UNSIGNED DEFAULT 18446744073709551614, b BIGINT DEFAULT -9223372036854775808, c DECIMAL(30,0) DEFAULT 99999999999999999999, g DECIMAL(30,0) AS (c + 99999999999999999999) VIRTUAL, h BIGINT UNSIGNED AS (a + 9223372036854775808) VIRTUAL, CHECK (a < 18446744073709551615))"},
		{"engine spellings", "CREATE TABLE t (a INT) ENGINE=myisam; ALTER TABLE t ENGINE=innodb; CREATE TABLE p (a INT) PARTITION BY HASH (a) (PARTITION p0 ENGINE=innodb, PARTITION p1 ENGINE=INNODB)"},
		{"generated column checks", "CREATE TABLE t (a INT, b INT AS (a + 1) VIRTUAL CHECK (b > 0), c INT AS (a * 2) STORED COMMENT 'c' CHECK (c < 100))"},
		{"temporal defaults", "CREATE TABLE t (a TIME DEFAULT '1:2', b TIME(1) DEFAULT '-1 2:3:4.56', c DATETIME(3) DEFAULT '2020-1-2 3:4:5.6789', d DATE DEFAULT '2020-1-2', e DATETIME DEFAULT '2020-01-02', f TIMESTAMP(2) NULL DEFAULT '2020-01-02 03:04:05', g TIME DEFAULT '10:00:00.9', h DATETIME(1) DEFAULT '2020-01-02T03:04:05.15')"},
		{"zerofill defaults", "CREATE TABLE t (a INT(5) ZEROFILL DEFAULT 3, b INT ZEROFILL DEFAULT 42, c DECIMAL(6,2) ZEROFILL DEFAULT 1.5, d FLOAT(7,3) ZEROFILL DEFAULT 1.5, g INT(3) ZEROFILL DEFAULT 12345, h DOUBLE ZEROFILL DEFAULT 2)"},
		{"comments", "CREATE TABLE t (a INT COMMENT 'x\ny\rz\\\\w''q', b INT COMMENT 'emoji \xf0\x9f\x93\x9d here', KEY k (a) COMMENT 'k\xf0\x9f\x94\xa5') COMMENT 'line1\nline2 \xf0\x9f\x94\xa5'"},
		{"table option order", "CREATE TABLE t (id INT PRIMARY KEY) ENGINE=myisam COMMENT='c' ROW_FORMAT=FIXED CHECKSUM=1 DELAY_KEY_WRITE=1 PACK_KEYS=1 AVG_ROW_LENGTH=10 MAX_ROWS=100 MIN_ROWS=5 STATS_SAMPLE_PAGES=7 STATS_AUTO_RECALC=1 STATS_PERSISTENT=0 PAGE_CHECKSUM=1"},
		{"hash keys", "CREATE TABLE t (id INT PRIMARY KEY, b INT, KEY b (b) USING HASH) ENGINE=InnoDB; CREATE INDEX c ON t (id) USING HASH"},
		{"long unique", "CREATE TABLE t (id INT PRIMARY KEY, a TEXT, b BLOB, c VARCHAR(10), UNIQUE KEY ua (a), UNIQUE KEY ub (b(10)), UNIQUE KEY uab (a, c)); CREATE TABLE u (a TEXT UNIQUE); CREATE UNIQUE INDEX ux ON t (b)"},
		{"generated columns", "CREATE TABLE t (a INT, b INT, s VARCHAR(10), price DECIMAL(10,2), tax DECIMAL(5,2), g1 DECIMAL(10,2) AS (price * (1 + tax / 100)) STORED, g2 INT AS (a + b * 2 - (a - b)) VIRTUAL, g3 VARCHAR(30) AS (concat(s, '-', a)) VIRTUAL, g4 INT AS (a and b) VIRTUAL, g5 INT AS (not (a = 1)) VIRTUAL, g6 INT AS (if(a > b, a, b)) VIRTUAL, g7 INT AS (a between 1 and 10) VIRTUAL, g8 INT AS (a << 2 | b & 1) VIRTUAL)"},
		{"sequence defaults", "CREATE SEQUENCE s1; CREATE TABLE t (a INT DEFAULT nextval(s1), b INT DEFAULT (NEXT VALUE FOR s1), c INT DEFAULT lastval(s1), d INT DEFAULT (PREVIOUS VALUE FOR s1), e INT DEFAULT setval(s1, 5))"},
		{"view qualifiers", "CREATE TABLE t (a INT, b INT); CREATE FUNCTION f(x INT) RETURNS INT DETERMINISTIC RETURN x; CREATE VIEW v AS SELECT " + "DATABASE_NAME.t.a, DATABASE_NAME.f(1) AS one FROM DATABASE_NAME.t"},
	}
	conn, err := testcontainer.OpenDatabase(t).Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	database := testcontainer.DatabaseName(t)
	for _, sc := range scripts {
		t.Run(sc.name, func(t *testing.T) {
			sql := strings.ReplaceAll(sc.sql, "DATABASE_NAME", quoteIdent(database))
			resetDatabase(t, conn, database)
			applySchema(t, conn, sql)
			want := serverObjects(t, conn)
			c := New()
			results, err := c.Exec("CREATE DATABASE "+quoteIdent(database)+"; USE "+quoteIdent(database)+";\n"+sql, nil)
			if err != nil {
				t.Fatal(err)
			}
			for _, r := range results {
				if r.Error != nil {
					t.Fatalf("line %d: %v", r.Line, r.Error)
				}
			}
			got := catalogObjects(c, database)
			for name, w := range want {
				if got[name] != w {
					t.Errorf("%s\n--- catalog\n%s\n--- server\n%s", name, got[name], w)
				}
			}
			if len(got) != len(want) {
				t.Errorf("catalog holds %d objects, server %d", len(got), len(want))
			}
		})
	}
}
