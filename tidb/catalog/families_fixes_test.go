package catalog

import (
	"errors"
	"strings"
	"testing"
)

const famFKSetup = "CREATE TABLE p (id INT PRIMARY KEY, n INT); CREATE TABLE c (id INT PRIMARY KEY, x INT, CONSTRAINT fk FOREIGN KEY (x) REFERENCES p (id) ON UPDATE SET NULL)"

const famGenSetup = "SET GLOBAL tidb_enable_check_constraint = ON; CREATE TABLE g (id INT PRIMARY KEY, a INT, b INT, v INT AS (a + b) VIRTUAL, s INT AS (a) STORED, KEY kv (v), CONSTRAINT ck CHECK (a >= 0))"

// TestFamilyFixes pins, outside simulate.tsv, the catalog fixes the family
// rows (simulate_families_test.go) closed. Each case was read on v8.5.5; code
// 0 means the statement runs, and then the SHOW CREATE TABLE of the table
// named by show holds want.
func TestFamilyFixes(t *testing.T) {
	for _, tt := range []struct {
		name, setup, stmt string
		code              int
		show, want        string
	}{
		{name: "coalesce every partition", setup: "CREATE TABLE h (id INT) PARTITION BY HASH (id) PARTITIONS 4", stmt: "ALTER TABLE h COALESCE PARTITION 4", code: 1508},
		{name: "coalesce none", setup: "CREATE TABLE h (id INT) PARTITION BY HASH (id) PARTITIONS 4", stmt: "ALTER TABLE h COALESCE PARTITION 0", code: 1515},
		{name: "drop a HASH partition", setup: "CREATE TABLE h (id INT) PARTITION BY HASH (id) PARTITIONS 4", stmt: "ALTER TABLE h DROP PARTITION p0", code: 1512},
		{name: "reorganize HASH", setup: "CREATE TABLE h (id INT) PARTITION BY HASH (id) PARTITIONS 4", stmt: "ALTER TABLE h REORGANIZE PARTITION p0, p1 INTO (PARTITION p01)", code: 8200},
		{name: "named partition added to HASH", setup: "CREATE TABLE h (id INT) PARTITION BY HASH (id) PARTITIONS 4", stmt: "ALTER TABLE h ADD PARTITION (PARTITION p4)", show: "h", want: "PARTITIONS 5"},
		{name: "drop an unknown partition", setup: rangeSetup, stmt: "ALTER TABLE r DROP PARTITION nope", code: 1507},
		{name: "drop every partition", setup: rangeSetup, stmt: "ALTER TABLE r DROP PARTITION p0, p1", code: 1508},
		{name: "added bound below", setup: rangeSetup, stmt: "ALTER TABLE r ADD PARTITION (PARTITION p9 VALUES LESS THAN (5))", code: 1493},
		{name: "reorganized bound past the next", setup: rangeSetup, stmt: "ALTER TABLE r REORGANIZE PARTITION p0 INTO (PARTITION p0 VALUES LESS THAN (30))", code: 1493},
		{name: "unknown partition truncated", setup: rangeSetup, stmt: "ALTER TABLE r TRUNCATE PARTITION p0, nope", code: 1735},
		{name: "partition column renamed", setup: rangeSetup, stmt: "ALTER TABLE r RENAME COLUMN a TO b", code: 3855},
		{name: "foreign key on a partitioned table", setup: rangeSetup + "; CREATE TABLE q (id INT PRIMARY KEY)", stmt: "ALTER TABLE r ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES q (id)", code: 1506},
		{name: "LIST value twice", stmt: "CREATE TABLE l (a INT) PARTITION BY LIST (a) (PARTITION p0 VALUES IN (1), PARTITION p1 VALUES IN (1))", code: 1495},
		{name: "partition name twice", stmt: "CREATE TABLE r2 (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p0 VALUES LESS THAN (20))", code: 1517},
		{name: "KEY over TEXT", stmt: "CREATE TABLE k (a TEXT) PARTITION BY KEY (a)", code: 1659},
		{name: "HASH arithmetic written compact", stmt: "CREATE TABLE h6 (a INT, b INT) PARTITION BY HASH (a + b) PARTITIONS 2", show: "h6", want: "HASH ((`a`+`b`))"},
		{name: "LIKE copies partitioning", setup: rangeSetup, stmt: "CREATE TABLE r3 LIKE r", show: "r3", want: "PARTITION BY RANGE"},
		{name: "CONVERT over an ascii column", setup: "CREATE TABLE n (a CHAR(5) CHARACTER SET ascii) CHARSET = latin1", stmt: "ALTER TABLE n CONVERT TO CHARACTER SET utf8mb4", code: 8200},
		{name: "column moved to the table's character set", setup: "CREATE TABLE n (a VARCHAR(5) CHARACTER SET utf8mb4) CHARSET = latin1", stmt: "ALTER TABLE n MODIFY a VARCHAR(5)", code: 8200},
		{name: "VARCHAR past 65535 bytes", stmt: "CREATE TABLE n (a VARCHAR(16384)) CHARSET = utf8mb4", code: 1074},
		{name: "CHAR past 255", stmt: "CREATE TABLE n (a CHAR(300))", code: 1074},
		{name: "VARCHAR past 65535 bytes, not strict", setup: "SET sql_mode = ''", stmt: "CREATE TABLE n (a VARCHAR(70000)) CHARSET = latin1"},
		{name: "foreign key name taken", setup: famFKSetup, stmt: "ALTER TABLE c ADD CONSTRAINT fk FOREIGN KEY (x) REFERENCES p (id)", code: 1826},
		{name: "referenced column missing", setup: famFKSetup, stmt: "ALTER TABLE c ADD CONSTRAINT f2 FOREIGN KEY (x) REFERENCES p (nope)", code: 3734},
		{name: "referenced columns miscounted", setup: famFKSetup, stmt: "ALTER TABLE c ADD CONSTRAINT f2 FOREIGN KEY (x) REFERENCES p (id, n)", code: 1239},
		{name: "NOT NULL under SET NULL", setup: famFKSetup, stmt: "ALTER TABLE c MODIFY x INT NOT NULL"},
		{name: "parent truncated", setup: famFKSetup, stmt: "TRUNCATE TABLE p", code: 1701},
		{name: "parent column dropped", setup: famFKSetup, stmt: "ALTER TABLE p DROP COLUMN id", code: 1829},
		{name: "parent dropped alone", setup: famFKSetup, stmt: "DROP TABLE p", code: 3730},
		{name: "parent dropped with its children", setup: famFKSetup, stmt: "DROP TABLE p, c"},
		{name: "backing key dropped with its key", setup: famFKSetup, stmt: "ALTER TABLE c DROP FOREIGN KEY fk, DROP INDEX fk", code: 1553},
		{name: "rename followed", setup: famFKSetup, stmt: "RENAME TABLE p TO q", show: "c", want: "REFERENCES `q` (`id`)"},
		{name: "column rename followed", setup: famFKSetup, stmt: "ALTER TABLE p RENAME COLUMN id TO pid", show: "c", want: "REFERENCES `p` (`pid`)"},
		{name: "move under the same name not followed", setup: famFKSetup + "; CREATE DATABASE o", stmt: "RENAME TABLE p TO o.p", show: "c", want: "REFERENCES `p` (`id`)"},
		{name: "rename into another schema not followed", setup: famFKSetup + "; CREATE DATABASE o", stmt: "RENAME TABLE p TO o.q", show: "c", want: "REFERENCES `p` (`id`)"},
		{name: "rename into a missing schema", setup: famFKSetup, stmt: "RENAME TABLE p TO nodb.p", code: 1025},
		{name: "RENAME TO beside another clause", setup: famFKSetup, stmt: "ALTER TABLE c RENAME TO c2, ADD COLUMN y INT", code: 8200},
		{name: "stored column added with a key", setup: famGenSetup, stmt: "ALTER TABLE g ADD COLUMN w INT AS (a) STORED, ADD INDEX kw (w)", code: 3106},
		{name: "VIRTUAL made plain", setup: famGenSetup, stmt: "ALTER TABLE g MODIFY v INT", code: 3106},
		{name: "STORED made plain", setup: famGenSetup, stmt: "ALTER TABLE g MODIFY s INT"},
		{name: "keyed VIRTUAL redefined", setup: famGenSetup, stmt: "ALTER TABLE g MODIFY v INT AS (a * b) VIRTUAL", code: 3106},
		{name: "read column renamed by CHANGE", setup: famGenSetup, stmt: "ALTER TABLE g CHANGE b bb INT", code: 3108},
		{name: "read column dropped with its readers", setup: famGenSetup, stmt: "ALTER TABLE g DROP COLUMN v, DROP COLUMN b", code: 3108},
		{name: "generated column read before it is defined", stmt: "CREATE TABLE g3 (a INT, b INT AS (c + 1), c INT AS (a + 1))", code: 3107},
		{name: "generated column reads AUTO_INCREMENT", stmt: "CREATE TABLE g4 (a INT AUTO_INCREMENT PRIMARY KEY, b INT AS (a + 1))", code: 3109},
		{name: "expression key over a missing column", setup: famGenSetup, stmt: "ALTER TABLE g ADD INDEX kb ((nope + 1))", code: 1054},
		{name: "expression key over a bare column", setup: famGenSetup, stmt: "ALTER TABLE g ADD INDEX kb ((a))", code: 3762},
		{name: "expression primary key", stmt: "CREATE TABLE f (a INT, PRIMARY KEY ((a + 1)))", code: 3756},
		{name: "primary key made invisible", setup: famGenSetup, stmt: "ALTER TABLE g ALTER INDEX `PRIMARY` INVISIBLE", code: 1176},
		{name: "CHECK name taken in the schema", setup: famGenSetup, stmt: "CREATE TABLE k (a INT, CONSTRAINT ck CHECK (a > 0))", code: 3822},
		{name: "CHECK over a missing column", setup: famGenSetup, stmt: "ALTER TABLE g ADD CHECK (a < nope)", code: 1054},
		{name: "CHECK calls RAND", setup: famGenSetup, stmt: "ALTER TABLE g ADD CONSTRAINT cr CHECK (a > RAND())", code: 3814},
		{name: "CHECK reads AUTO_INCREMENT", setup: "SET GLOBAL tidb_enable_check_constraint = ON", stmt: "CREATE TABLE k4 (a INT AUTO_INCREMENT PRIMARY KEY, CHECK (a > 0))", code: 3818},
		{name: "unknown CHECK dropped", setup: famGenSetup, stmt: "ALTER TABLE g DROP CHECK nope", code: 3940},
		{name: "column redefined with a CHECK", setup: famGenSetup, stmt: "ALTER TABLE g MODIFY a INT CHECK (a > 1)", code: 8200},
		{name: "not-equal in a CHECK", setup: "SET GLOBAL tidb_enable_check_constraint = ON", stmt: "CREATE TABLE k5 (a INT, b INT, CHECK (a <> b))", show: "k5", want: "CHECK ((`a` != `b`))"},
		{name: "LIKE renames the source's CHECK", setup: famGenSetup, stmt: "CREATE TABLE g5 LIKE g", show: "g", want: "CONSTRAINT `g5_chk_1`"},
		{name: "key under ALGORITHM = INSTANT", setup: famFKSetup, stmt: "ALTER TABLE p ADD INDEX kn (n), ALGORITHM = INSTANT", code: 1846},
		{name: "column dropped and added again", setup: famFKSetup, stmt: "ALTER TABLE p DROP COLUMN n, ADD COLUMN n BIGINT", code: 1060},
		{name: "FULLTEXT over a column the statement adds", setup: famFKSetup, stmt: "ALTER TABLE p ADD COLUMN e VARCHAR(10), ADD FULLTEXT INDEX ft (e), ALGORITHM = INPLACE"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := wtSetup(t)
			if tt.setup != "" {
				wtExec(t, c, tt.setup)
			}
			results, err := c.Exec(tt.stmt, nil)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			var got error
			for _, r := range results {
				if r.Error != nil {
					got = r.Error
					break
				}
			}
			if tt.code != 0 {
				var catErr *Error
				if !errors.As(got, &catErr) || catErr.Code != tt.code {
					t.Fatalf("got %v, want %d", got, tt.code)
				}
				return
			}
			if got != nil {
				t.Fatalf("got %v, want it to run", got)
			}
			if tt.show != "" {
				if text := c.ShowCreateTable("testdb", tt.show); !strings.Contains(text, tt.want) {
					t.Errorf("SHOW CREATE TABLE %s holds no %q:\n%s", tt.show, tt.want, text)
				}
			}
		})
	}
}

const rangeSetup = "CREATE TABLE r (a INT NOT NULL) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))"

// TestReviewRoundOne pins what the first review of the port asked for.
func TestReviewRoundOne(t *testing.T) {
	t.Run("refused ALTER leaves the references", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, famFKSetup)
		results, _ := c.Exec("ALTER TABLE p RENAME COLUMN id TO pid, DROP COLUMN nope", nil)
		if results[0].Error == nil {
			t.Fatal("the ALTER ran")
		}
		if text := c.ShowCreateTable("testdb", "c"); !strings.Contains(text, "REFERENCES `p` (`id`)") {
			t.Errorf("c's key moved:\n%s", text)
		}
	})
	t.Run("DROP TABLE drops what exists, then answers 1051", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, famFKSetup)
		results, _ := c.Exec("DROP TABLE p, missing, c", nil)
		var catErr *Error
		if !errors.As(results[0].Error, &catErr) || catErr.Code != 1051 {
			t.Fatalf("got %v, want 1051", results[0].Error)
		}
		if db := c.GetDatabase("testdb"); db.GetTable("p") != nil || db.GetTable("c") != nil {
			t.Error("p and c should be gone, as on v8.5.5")
		}
	})
	t.Run("a literal reads no column", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "CREATE TABLE t (id INT AUTO_INCREMENT PRIMARY KEY, label VARCHAR(20) AS ('`id`'))")
	})
	for _, set := range []string{"SET @@sql_mode = ''", "SET @@SESSION.sql_mode = ''", "SET SESSION sql_mode = ''"} {
		t.Run(set, func(t *testing.T) {
			c := wtSetup(t)
			wtExec(t, c, set)
			wtExec(t, c, "CREATE TABLE t (v VARCHAR(20000)) CHARACTER SET utf8mb4")
		})
	}
	t.Run("SET GLOBAL sql_mode keeps the session strict", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "SET @@GLOBAL.sql_mode = ''")
		results, _ := c.Exec("CREATE TABLE t (v VARCHAR(20000)) CHARACTER SET utf8mb4", nil)
		var catErr *Error
		if !errors.As(results[0].Error, &catErr) || catErr.Code != 1074 {
			t.Fatalf("got %v, want 1074", results[0].Error)
		}
	})
}

// TestReviewRoundTwo pins what the second review of the port asked for.
func TestReviewRoundTwo(t *testing.T) {
	// A rename to a name the table began with is 1060 on v8.5.5, which also
	// keeps the statement the review drew from adding a key the renames could
	// reach.
	for _, stmt := range []string{
		"ALTER TABLE t RENAME COLUMN a TO x, RENAME COLUMN b TO a, ADD CONSTRAINT fk FOREIGN KEY (c) REFERENCES t (a)",
		"ALTER TABLE t CHANGE a x BIGINT, CHANGE b a INT",
		"ALTER TABLE t DROP COLUMN c, RENAME COLUMN b TO c",
	} {
		t.Run(stmt, func(t *testing.T) {
			c := wtSetup(t)
			wtExec(t, c, "CREATE TABLE t (a BIGINT PRIMARY KEY, b INT UNIQUE, c INT)")
			results, _ := c.Exec(stmt, nil)
			assertError(t, results[0].Error, ErrDupColumn)
			if text := c.ShowCreateTable("testdb", "t"); strings.Contains(text, "`x`") || strings.Contains(text, "FOREIGN KEY") {
				t.Errorf("the refused statement changed t:\n%s", text)
			}
		})
	}
	t.Run("a column renamed away runs", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "CREATE TABLE t (a INT, b INT)")
		wtExec(t, c, "ALTER TABLE t RENAME COLUMN a TO x, ADD COLUMN y INT")
	})
	t.Run("SET GLOBAL foreign_key_checks keeps the session's", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, famFKSetup)
		wtExec(t, c, "SET @@GLOBAL.foreign_key_checks = 0")
		results, _ := c.Exec("DROP TABLE p", nil)
		if results[0].Error == nil {
			t.Fatal("DROP TABLE p ran with the session's checks on")
		}
	})
	for _, set := range []string{"SET sql_mode = @@sql_mode", "SET sql_mode = CONCAT(@@sql_mode, ',NO_AUTO_VALUE_ON_ZERO')", "SET sql_mode = DEFAULT", "SET sql_mode = TRADITIONAL"} {
		t.Run(set, func(t *testing.T) {
			c := wtSetup(t)
			wtExec(t, c, set)
			results, _ := c.Exec("CREATE TABLE t (v VARCHAR(20000)) CHARACTER SET utf8mb4", nil)
			var catErr *Error
			if !errors.As(results[0].Error, &catErr) || catErr.Code != 1074 {
				t.Fatalf("got %v, want 1074 under a strict mode", results[0].Error)
			}
		})
	}
}
