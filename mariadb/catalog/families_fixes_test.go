package catalog

import (
	"errors"
	"strings"
	"testing"
)

// TestFamilyFixes pins, outside simulate.tsv, the catalog fixes the family
// rows (simulate_families_test.go) closed. Each case was read on
// mariadb:11.8.8; code 0 means the statement runs, and then the SHOW CREATE
// TABLE of the table named by show holds want.
func TestFamilyFixes(t *testing.T) {
	for _, tt := range []struct {
		name  string
		setup string
		stmt  string
		code  int
		show  string
		want  string
	}{
		{name: "REF_SYSTEM_ID after a spatial type", stmt: "CREATE TABLE p (loc POINT REF_SYSTEM_ID = 4326 NOT NULL)", show: "p", want: "`loc` point NOT NULL"},
		{name: "partitions coalesced to none", setup: "CREATE TABLE s (id INT) PARTITION BY HASH (id) PARTITIONS 4", stmt: "ALTER TABLE s COALESCE PARTITION 4", code: 1508},
		{name: "no partition coalesced", setup: "CREATE TABLE s (id INT) PARTITION BY HASH (id) PARTITIONS 4", stmt: "ALTER TABLE s COALESCE PARTITION 0", code: 1515},
		{name: "partition name taken", setup: "CREATE TABLE s (id INT) PARTITION BY HASH (id) PARTITIONS 4", stmt: "ALTER TABLE s ADD PARTITION (PARTITION p0)", code: 1517},
		{name: "named HASH partitions listed", stmt: "CREATE TABLE h (a INT) PARTITION BY HASH (a) PARTITIONS 2 (PARTITION x, PARTITION y)", show: "h", want: "(PARTITION `x` ENGINE = InnoDB,"},
		{name: "REORGANIZE keeps the first HASH partition", setup: "CREATE TABLE s (id INT) PARTITION BY HASH (id) PARTITIONS 4", stmt: "ALTER TABLE s REORGANIZE PARTITION", show: "s", want: "PARTITIONS 1"},
		{name: "repartition applies", setup: "CREATE TABLE s (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 4", stmt: "ALTER TABLE s PARTITION BY KEY (id) PARTITIONS 8", show: "s", want: "PARTITION BY KEY (`id`)\nPARTITIONS 8"},
		{name: "RANGE bounds rise", stmt: "CREATE TABLE r (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (5))", code: 1493},
		{name: "LIST value twice", stmt: "CREATE TABLE l (a INT) PARTITION BY LIST (a) (PARTITION p0 VALUES IN (1), PARTITION p1 VALUES IN (1))", code: 1495},
		{name: "HASH over a string column", stmt: "CREATE TABLE h (a VARCHAR(10)) PARTITION BY HASH (a)", code: 1659},
		{name: "KEY over TEXT", stmt: "CREATE TABLE h (a TEXT) PARTITION BY KEY (a)", code: 1502},
		{name: "partition expression follows a rename", setup: "CREATE TABLE o (id INT NOT NULL, d DATE NOT NULL, PRIMARY KEY (id, d)) PARTITION BY RANGE (YEAR(d)) (PARTITION p0 VALUES LESS THAN (2025))", stmt: "ALTER TABLE o RENAME COLUMN d TO dt", show: "o", want: "PARTITION BY RANGE (year(`dt`))"},
		{name: "LIKE copies partitioning", setup: "CREATE TABLE s (id INT) PARTITION BY HASH (id) PARTITIONS 2", stmt: "CREATE TABLE s2 LIKE s", show: "s2", want: "PARTITIONS 2"},
		{name: "unknown partition truncated", setup: "CREATE TABLE r (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))", stmt: "ALTER TABLE r TRUNCATE PARTITION p0, nope", code: 1735},
		{name: "VARCHAR past 65532 bytes", stmt: "CREATE TABLE n (a VARCHAR(16384)) CHARSET = utf8mb4", code: 1074},
		{name: "row past 65535 bytes", stmt: "CREATE TABLE n (a INT, b VARCHAR(65530)) CHARSET = latin1", code: 1118},
		{name: "CONVERT TO widens TEXT", setup: "CREATE TABLE n (b TEXT) CHARSET = latin1", stmt: "ALTER TABLE n CONVERT TO CHARACTER SET utf8mb4", show: "n", want: "`b` mediumtext"},
		{name: "CONVERT TO unknown character set", setup: "CREATE TABLE n (b TEXT)", stmt: "ALTER TABLE n CONVERT TO CHARACTER SET nope", code: 1115},
		{name: "foreign key name taken", setup: fkSetup, stmt: "CREATE TABLE c2 (p INT, CONSTRAINT fk FOREIGN KEY (p) REFERENCES p (id))", code: 1005},
		{name: "new SET NULL key on NOT NULL", setup: fkSetup, stmt: "CREATE TABLE c2 (p INT NOT NULL, FOREIGN KEY (p) REFERENCES p (id) ON DELETE SET NULL)", code: 1005},
		{name: "SET NULL column made NOT NULL", setup: fkSetup, stmt: "ALTER TABLE c MODIFY p INT NOT NULL", code: 1830},
		{name: "child column type changed", setup: fkSetup, stmt: "ALTER TABLE c MODIFY p BIGINT", code: 1832},
		{name: "parent column type changed", setup: fkSetup, stmt: "ALTER TABLE p MODIFY id BIGINT", code: 1833},
		{name: "parent column dropped", setup: fkSetup, stmt: "ALTER TABLE p DROP COLUMN id", code: 1829},
		{name: "parent dropped", setup: fkSetup, stmt: "DROP TABLE p", code: 1451},
		{name: "parent truncated", setup: fkSetup, stmt: "TRUNCATE TABLE p", code: 1701},
		{name: "foreign key table off InnoDB", setup: fkSetup, stmt: "ALTER TABLE c ENGINE = MyISAM", code: 1217},
		{name: "foreign key follows a rename", setup: fkSetup, stmt: "RENAME TABLE p TO q", show: "c", want: "REFERENCES `q` (`id`)"},
		{name: "foreign key follows a column rename", setup: fkSetup, stmt: "ALTER TABLE p RENAME COLUMN id TO pid", show: "c", want: "REFERENCES `p` (`pid`)"},
		{name: "MyISAM keeps no foreign key", stmt: "CREATE TABLE m (a INT, FOREIGN KEY (a) REFERENCES nope (id)) ENGINE = MyISAM", show: "m", want: "KEY `a` (`a`)\n)"},
		{name: "CASCADE over a generated column", setup: fkSetup, stmt: "CREATE TABLE g (a INT, v INT AS (a + 1) STORED, FOREIGN KEY (v) REFERENCES p (id) ON UPDATE CASCADE)", code: 1905},
		{name: "VIRTUAL made PERSISTENT", setup: "CREATE TABLE g (a INT, v INT AS (a + 1) VIRTUAL)", stmt: "ALTER TABLE g MODIFY v INT AS (a + 1) PERSISTENT", code: 1907},
		{name: "invisible NOT NULL without default", setup: "CREATE TABLE i (a INT)", stmt: "ALTER TABLE i ADD COLUMN b INT NOT NULL INVISIBLE", code: 4108},
		{name: "no visible column", stmt: "CREATE TABLE i (a INT INVISIBLE)", code: 1113},
		{name: "RAND in a CHECK", stmt: "CREATE TABLE k (a INT, CONSTRAINT c1 CHECK (a > RAND()))", code: 1901},
		{name: "AUTO_INCREMENT in a CHECK", stmt: "CREATE TABLE k (a INT AUTO_INCREMENT PRIMARY KEY, CHECK (a > 0))", code: 1901},
		{name: "column CHECK in MODIFY", setup: "CREATE TABLE k (a INT)", stmt: "ALTER TABLE k MODIFY a INT CHECK (a > 1)", show: "k", want: "CHECK (`a` > 1)"},
		{name: "TEXT key prefixed", stmt: "CREATE TABLE b (t TEXT, KEY (t))", show: "b", want: "KEY `t` (`t`(768))"},
		{name: "TEXT key beside another part", stmt: "CREATE TABLE b (t TEXT, a INT, KEY (t, a))", code: 1071},
		{name: "POINT key prefixed", stmt: "CREATE TABLE b (g POINT, a INT, KEY (g(10), a))", show: "b", want: "KEY `g` (`g`(25),`a`)"},
		{name: "nullable column under a SPATIAL key", setup: "CREATE TABLE s (g POINT NOT NULL, SPATIAL KEY (g))", stmt: "ALTER TABLE s MODIFY g POINT NULL", code: 1252},
		{name: "Aria page checksum", stmt: "CREATE TABLE a (x INT) ENGINE = Aria", show: "a", want: "PAGE_CHECKSUM=1"},
		{name: "ALGORITHM=INSTANT for ADD INDEX", setup: "CREATE TABLE t (b INT PRIMARY KEY, c INT)", stmt: "ALTER TABLE t ADD INDEX i (c), ALGORITHM = INSTANT", code: 1846},
		{name: "ALGORITHM=NOCOPY for NOT NULL", setup: "CREATE TABLE t (b INT PRIMARY KEY, c INT)", stmt: "ALTER TABLE t MODIFY c INT NOT NULL, ALGORITHM = NOCOPY", code: 1845},
		{name: "ALGORITHM=INSTANT for ADD COLUMN FIRST", setup: "CREATE TABLE t (b INT PRIMARY KEY, c INT)", stmt: "ALTER TABLE t ADD COLUMN e INT FIRST, ALGORITHM = INSTANT"},
		{name: "event ENDS before STARTS", stmt: "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR STARTS '2031-01-01 00:00:00' ENDS '2030-01-01 00:00:00' DO COMMIT", code: 1543},
		{name: "event EVERY 0", stmt: "CREATE EVENT e ON SCHEDULE EVERY 0 DAY DO COMMIT", code: 1542},
		{name: "past event moved", setup: "CREATE EVENT e ON SCHEDULE AT '2030-01-01 00:00:00' DO COMMIT", stmt: "ALTER EVENT e ON SCHEDULE AT '2000-01-01 00:00:00'", code: 1589},
		{name: "rename into a missing schema", setup: "CREATE TABLE r (a INT)", stmt: "RENAME TABLE r TO nodb.r", code: 1025},
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
			got := results[len(results)-1].Error
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

const ftSetup = "CREATE TABLE f (id INT PRIMARY KEY, name VARCHAR(10) NOT NULL, g POINT NOT NULL, FULLTEXT KEY ft (name), SPATIAL KEY sg (g))"

const fkSetup = "CREATE TABLE p (id INT PRIMARY KEY, n INT); CREATE TABLE c (id INT PRIMARY KEY, p INT, CONSTRAINT fk FOREIGN KEY (p) REFERENCES p (id) ON UPDATE SET NULL)"

// TestDropTableGoesOn pins DROP TABLE's dropping the tables it can when one is
// referenced (mariadb:11.8.8: DROP TABLE p, c2, c over a parent p leaves p).
func TestDropTableGoesOn(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, fkSetup+"; CREATE TABLE c2 (a INT)")
	results, _ := c.Exec("DROP TABLE p, c2", nil)
	var catErr *Error
	if !errors.As(results[0].Error, &catErr) || catErr.Code != 1451 {
		t.Fatalf("got %v, want 1451", results[0].Error)
	}
	db := c.GetDatabase("testdb")
	if db.GetTable("p") == nil || db.GetTable("c2") != nil {
		t.Errorf("want p kept and c2 dropped")
	}
}

// TestReviewRoundOne pins the answers the first review of the port asked
// for, read on mariadb:11.8.8.
func TestReviewRoundOne(t *testing.T) {
	for _, tt := range []struct {
		name, setup, stmt string
		code              int
	}{
		{"child display width", fkSetup, "ALTER TABLE c MODIFY p INT(10)", 0},
		{"parent display width", fkSetup, "ALTER TABLE p MODIFY id INT(5)", 0},
		{"child signedness", fkSetup, "ALTER TABLE c MODIFY p INT UNSIGNED", 1832},
		{"skipped key under INSTANT", "CREATE TABLE t (b INT PRIMARY KEY, c INT, KEY k (c))", "ALTER TABLE t ADD INDEX IF NOT EXISTS k (c), ADD COLUMN d INT, ALGORITHM = INSTANT", 0},
		{"absent key dropped under INSTANT", "CREATE TABLE t (b INT PRIMARY KEY, c INT)", "ALTER TABLE t DROP INDEX IF EXISTS k, ADD COLUMN d INT, ALGORITHM = INSTANT", 0},
		{"column CHECK on AUTO_INCREMENT", "", "CREATE TABLE k (a INT AUTO_INCREMENT PRIMARY KEY CHECK (a > 0))", 1901},
		{"FULLTEXT key dropped after MODIFY", ftSetup, "ALTER TABLE f MODIFY name INT, DROP INDEX ft", 0},
		{"SPATIAL key dropped after MODIFY", ftSetup, "ALTER TABLE f MODIFY g POINT NULL, DROP INDEX sg", 0},
		{"FULLTEXT key kept", ftSetup, "ALTER TABLE f MODIFY name INT", 1283},
		{"VARCHAR widened under a foreign key", "CREATE TABLE p (s VARCHAR(10), UNIQUE KEY (s)); CREATE TABLE c (s VARCHAR(10), FOREIGN KEY (s) REFERENCES p (s))", "ALTER TABLE c MODIFY s VARCHAR(20)", 1832},
		{"engine moved with the foreign keys dropped", fkSetup, "ALTER TABLE c DROP FOREIGN KEY fk, ENGINE = MyISAM", 1217},
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
			got := results[len(results)-1].Error
			var catErr *Error
			switch {
			case tt.code == 0 && got != nil:
				t.Fatalf("got %v, want it to run", got)
			case tt.code != 0 && (!errors.As(got, &catErr) || catErr.Code != tt.code):
				t.Fatalf("got %v, want %d", got, tt.code)
			}
		})
	}
}

// TestSDLKeepsPastEvents pins that a declarative schema keeps an event whose
// AT time has passed, as written: loading SDL runs no schedule.
func TestSDLKeepsPastEvents(t *testing.T) {
	c, err := LoadSDL("CREATE DATABASE d; USE d; CREATE EVENT e ON SCHEDULE AT '2000-01-01 00:00:00' DO COMMIT; CREATE EVENT f ON SCHEDULE AT '2000-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT")
	if err != nil {
		t.Fatal(err)
	}
	db := c.GetDatabase("d")
	if db.Events["e"] == nil {
		t.Fatal("event e was dropped")
	}
	if ev := db.Events["f"]; ev == nil || !strings.EqualFold(ev.Enable, "ENABLE") {
		t.Fatalf("event f is %+v, want it ENABLE", ev)
	}
}
