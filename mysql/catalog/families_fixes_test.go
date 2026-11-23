package catalog

import (
	"strings"
	"testing"
	"time"
)

// The behaviors below were read on mysql:8.0.46 while closing the
// differences familyRows (simulate_families_test.go) found; each case is what
// the pinned server answered or stored. The simulate rows hold the catalog to
// the engine on their fixtures; these hold the neighboring cases the rows do
// not write.

func familiesCatalog(t *testing.T, setup string) *Catalog {
	t.Helper()
	c := New()
	mustExec(t, c, "CREATE DATABASE a; CREATE DATABASE o; USE a;")
	if setup != "" {
		mustExec(t, c, setup)
	}
	return c
}

func familiesCode(t *testing.T, c *Catalog, sql string) int {
	t.Helper()
	results, err := c.Exec(sql, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("%s: parse error: %v", sql, err)
	}
	for _, r := range results {
		if r.Error != nil {
			if e, ok := r.Error.(*Error); ok {
				return e.Code
			}
			t.Fatalf("%s: %v", sql, r.Error)
		}
	}
	return 0
}

func TestViewQualifiesOtherSchemaColumns(t *testing.T) {
	c := familiesCatalog(t, "CREATE TABLE t (id INT); CREATE TABLE o.x (id INT);")
	for _, tc := range []struct{ view, body string }{
		{"CREATE VIEW v1 AS SELECT id FROM a.t", "select `t`.`id` AS `id` from `t`"},
		{"CREATE VIEW v2 AS SELECT x.id FROM o.x", "select `o`.`x`.`id` AS `id` from `o`.`x`"},
		{"CREATE VIEW v3 AS SELECT id FROM o.x AS x", "select `o`.`x`.`id` AS `id` from `o`.`x`"},
		{"CREATE VIEW v4 AS SELECT * FROM o.x JOIN t USING (id)", "select `o`.`x`.`id` AS `id` from (`o`.`x` join `t` on((`o`.`x`.`id` = `t`.`id`)))"},
		{"CREATE VIEW v5 AS SELECT id FROM o.x x2 WHERE id IN (SELECT id FROM o.x)", "select `x2`.`id` AS `id` from `o`.`x` `x2` where `x2`.`id` in (select `o`.`x`.`id` from `o`.`x`)"},
		{"CREATE VIEW v6 AS SELECT id FROM t AS T", "select `T`.`id` AS `id` from `t` `T`"},
		{"CREATE VIEW v7 AS SELECT (SELECT MAX(id) FROM o.x) m", "select (select max(`o`.`x`.`id`) from `o`.`x`) AS `m`"},
	} {
		mustExec(t, c, tc.view)
		name := strings.Fields(tc.view)[2]
		if got := c.GetDatabase("a").Views[name].Definition; got != tc.body {
			t.Errorf("%s:\n got %s\nwant %s", tc.view, got, tc.body)
		}
	}
}

func TestRenameRetargetsForeignKeys(t *testing.T) {
	c := familiesCatalog(t, `CREATE TABLE p (id INT PRIMARY KEY);
		CREATE TABLE ch (pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id));`)
	mustExec(t, c, "ALTER TABLE p RENAME COLUMN id TO pk")
	mustExec(t, c, "RENAME TABLE p TO o.p2")
	if got := c.ShowCreateTable("a", "ch"); !strings.Contains(got, "REFERENCES `o`.`p2` (`pk`)") {
		t.Errorf("child after renames:\n%s", got)
	}
	// A multi-table RENAME that fails puts the references back.
	if code := familiesCode(t, c, "RENAME TABLE o.p2 TO a.p, nosuch TO z"); code != ErrNoSuchTable {
		t.Fatalf("got %d, want %d", code, ErrNoSuchTable)
	}
	if got := c.ShowCreateTable("a", "ch"); !strings.Contains(got, "REFERENCES `o`.`p2` (`pk`)") {
		t.Errorf("child after failed rename:\n%s", got)
	}
	mustExec(t, c, "RENAME TABLE o.p2 TO p")
	if got := c.ShowCreateTable("a", "ch"); !strings.Contains(got, "REFERENCES `p` (`pk`)") {
		t.Errorf("child after the parent came back:\n%s", got)
	}
}

func TestHashPartitionListsByName(t *testing.T) {
	c := familiesCatalog(t, `CREATE TABLE s (id INT) PARTITION BY HASH (id) PARTITIONS 4;
		CREATE TABLE n (id INT) PARTITION BY HASH (id) (PARTITION a, PARTITION b, PARTITION c);`)
	mustExec(t, c, "ALTER TABLE s REORGANIZE PARTITION")
	if got := c.ShowCreateTable("a", "s"); !strings.Contains(got, "PARTITIONS 1 */") {
		t.Errorf("s:\n%s", got)
	}
	mustExec(t, c, "ALTER TABLE n REORGANIZE PARTITION")
	mustExec(t, c, "ALTER TABLE n ADD PARTITION (PARTITION d)")
	mustExec(t, c, "ALTER TABLE n ADD PARTITION PARTITIONS 2")
	got := c.ShowCreateTable("a", "n")
	for _, want := range []string{"PARTITION a ", "PARTITION d ", "PARTITION p2 ", "PARTITION p3 "} {
		if !strings.Contains(got, want) {
			t.Errorf("n lacks %q:\n%s", want, got)
		}
	}
	mustExec(t, c, "ALTER TABLE n COALESCE PARTITION 2")
	if got := c.ShowCreateTable("a", "n"); strings.Contains(got, "PARTITIONS") || strings.Contains(got, "PARTITION p2") {
		t.Errorf("n after COALESCE keeps its names:\n%s", got)
	}
}

func TestPartitionRefusals(t *testing.T) {
	setup := `CREATE TABLE r (id INT NOT NULL, y INT NOT NULL, PRIMARY KEY (id, y)) PARTITION BY RANGE (y)
			(PARTITION a VALUES LESS THAN (10), PARTITION b VALUES LESS THAN (20), PARTITION c VALUES LESS THAN (30));
		CREATE TABLE m (y INT) PARTITION BY RANGE (y) (PARTITION a VALUES LESS THAN (10), PARTITION b VALUES LESS THAN MAXVALUE);
		CREATE TABLE h (id INT) PARTITION BY HASH (id) PARTITIONS 2;
		CREATE TABLE plain (id INT NOT NULL, y INT NOT NULL, PRIMARY KEY (id, y));`
	for _, tc := range []struct {
		sql  string
		code int
	}{
		{"ALTER TABLE r REORGANIZE PARTITION a, zz INTO (PARTITION z VALUES LESS THAN (10))", 1507},
		{"ALTER TABLE r REORGANIZE PARTITION b INTO (PARTITION z VALUES LESS THAN (15))", 1520},
		{"ALTER TABLE r REORGANIZE PARTITION c INTO (PARTITION z VALUES LESS THAN (25))", 1520},
		{"ALTER TABLE r REORGANIZE PARTITION c INTO (PARTITION z VALUES LESS THAN (35))", 0},
		{"ALTER TABLE r REORGANIZE PARTITION a INTO (PARTITION z VALUES LESS THAN (5), PARTITION z2 VALUES LESS THAN (3))", 1493},
		{"ALTER TABLE r REORGANIZE PARTITION a INTO (PARTITION z VALUES LESS THAN (5), PARTITION z2 VALUES LESS THAN (10))", 0},
		{"ALTER TABLE r REBUILD PARTITION zz", 1507},
		{"ALTER TABLE r DROP PARTITION a, a", 1507},
		{"ALTER TABLE h DROP PARTITION zz", 1512},
		{"ALTER TABLE r EXCHANGE PARTITION zz WITH TABLE h", 1732},
		{"ALTER TABLE r EXCHANGE PARTITION zz WITH TABLE plain", 1735},
		{"ALTER TABLE r TRUNCATE PARTITION zz", 1735},
		{"ALTER TABLE r COALESCE PARTITION 0", 1509},
		{"ALTER TABLE h COALESCE PARTITION 0", 1515},
		{"ALTER TABLE h COALESCE PARTITION 2", 1508},
		{"ALTER TABLE h ADD PARTITION (PARTITION x VALUES IN (1))", 1480},
		{"ALTER TABLE r ADD PARTITION (PARTITION x VALUES LESS THAN (30))", 1493},
		{"ALTER TABLE r ADD PARTITION (PARTITION a VALUES LESS THAN (5))", 1517},
		{"ALTER TABLE m ADD PARTITION (PARTITION c VALUES LESS THAN (20))", 1481},
		{"ALTER TABLE h REORGANIZE PARTITION p0 INTO (PARTITION x, PARTITION y)", 1510},
		{"ALTER TABLE h REORGANIZE PARTITION p0, p1 INTO (PARTITION x, PARTITION y)", 0},
		{"CREATE TABLE x1 (a DATE) PARTITION BY HASH (a)", 1659},
		{"CREATE TABLE x2 (a TIME) PARTITION BY RANGE COLUMNS (a) (PARTITION p VALUES LESS THAN ('10:00:00'))", 0},
		{"CREATE TABLE x3 (a YEAR) PARTITION BY RANGE COLUMNS (a) (PARTITION p VALUES LESS THAN (2000))", 1659},
		{"CREATE TABLE x4 (a JSON) PARTITION BY KEY (a)", 1502},
		{"CREATE TABLE x5 (a FLOAT) PARTITION BY KEY (a)", 0},
		{"CREATE TABLE x6 (a INT, b VARCHAR(10)) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 2 (PARTITION p VALUES LESS THAN (1))", 1659},
		{"CREATE TABLE x7 (c VARCHAR(5)) PARTITION BY LIST COLUMNS (c) (PARTITION p0 VALUES IN ('us'), PARTITION p1 VALUES IN ('US'))", 1495},
	} {
		c := familiesCatalog(t, setup)
		if got := familiesCode(t, c, tc.sql); got != tc.code {
			t.Errorf("%s: got %d, want %d", tc.sql, got, tc.code)
		}
	}
}

func TestForeignKeyRefusals(t *testing.T) {
	setup := `CREATE TABLE p (id INT PRIMARY KEY, u INT UNIQUE);
		CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id) ON UPDATE SET NULL);
		CREATE TABLE s (id INT PRIMARY KEY, pid INT, CONSTRAINT fs FOREIGN KEY (pid) REFERENCES s (id));`
	for _, tc := range []struct {
		sql  string
		code int
	}{
		{"ALTER TABLE c ADD COLUMN v INT AS (pid + 1) STORED, ADD CONSTRAINT fk_v FOREIGN KEY (v) REFERENCES p (id) ON DELETE SET NULL", 3104},
		{"ALTER TABLE c ADD COLUMN v INT AS (pid + 1) VIRTUAL, ADD CONSTRAINT fk_v FOREIGN KEY (v) REFERENCES p (id)", 3733},
		{"TRUNCATE TABLE p", 1701},
		{"TRUNCATE TABLE s", 0},
		{"DROP TABLE p, c", 0},
		{"ALTER TABLE p DROP COLUMN u", 0},
		{"ALTER TABLE c CHANGE pid pid2 INT NOT NULL", 1830},
		{"ALTER TABLE p MODIFY id INT UNSIGNED", 3780},
		{"ALTER TABLE p MODIFY id INT NOT NULL COMMENT 'x'", 0},
		{"ALTER TABLE s ENGINE = MyISAM", 3776},
		{"ALTER TABLE p ENGINE = InnoDB", 0},
		{"SET foreign_key_checks = 0; ALTER TABLE p MODIFY id BIGINT", 3780},
		{"SET foreign_key_checks = 0; TRUNCATE TABLE p", 0},
		{"SET foreign_key_checks = 0; ALTER TABLE p DROP COLUMN id", 1829},
	} {
		c := familiesCatalog(t, setup)
		if got := familiesCode(t, c, tc.sql); got != tc.code {
			t.Errorf("%s: got %d, want %d", tc.sql, got, tc.code)
		}
	}
	c := familiesCatalog(t, setup)
	mustExec(t, c, "CREATE TABLE m (id INT, FOREIGN KEY (id) REFERENCES nosuch (id)) ENGINE = MyISAM")
	if got := c.ShowCreateTable("a", "m"); strings.Contains(got, "FOREIGN KEY") || !strings.Contains(got, "KEY `id` (`id`)") {
		t.Errorf("a MyISAM table keeps the key and no foreign key:\n%s", got)
	}
}

func TestCheckAndGeneratedColumnRefusals(t *testing.T) {
	setup := `CREATE TABLE people (id INT PRIMARY KEY, first_name VARCHAR(50) NOT NULL, last_name VARCHAR(50) NOT NULL,
		full_name VARCHAR(101) AS (CONCAT(first_name, ' ', last_name)) VIRTUAL, name_len INT AS (CHAR_LENGTH(first_name)) STORED,
		age INT, CONSTRAINT chk_age CHECK (age >= 0), CONSTRAINT chk_names CHECK (first_name <> last_name) NOT ENFORCED);`
	for _, tc := range []struct {
		sql  string
		code int
	}{
		{"ALTER TABLE people DROP CHECK chk_names, DROP COLUMN first_name", 3108},
		{"ALTER TABLE people RENAME COLUMN age TO years, DROP CHECK chk_age", 0},
		{"ALTER TABLE people ADD COLUMN z INT, ADD CONSTRAINT cz CHECK (z > 0)", 0},
		{"ALTER TABLE people DROP COLUMN age, ADD CONSTRAINT c2 CHECK (age > 1)", 1054},
		{"ALTER TABLE people ADD CHECK (age > RAND())", 3814},
		{"ALTER TABLE people ALTER COLUMN full_name DROP DEFAULT", 0},
		{"ALTER TABLE people ALTER COLUMN name_len SET DEFAULT 1", 1221},
		{"ALTER TABLE people DROP PRIMARY KEY, ADD PRIMARY KEY (id) INVISIBLE", 3522},
		{"CREATE TABLE gc5 (a INT, b INT AS (b + 1))", 3107},
		{"CREATE TABLE gc6 (a INT AUTO_INCREMENT PRIMARY KEY, b INT, CONSTRAINT c CHECK (a > 0))", 3818},
	} {
		c := familiesCatalog(t, setup)
		if got := familiesCode(t, c, tc.sql); got != tc.code {
			t.Errorf("%s: got %d, want %d", tc.sql, got, tc.code)
		}
	}
	c := familiesCatalog(t, setup)
	mustExec(t, c, "ALTER TABLE people DROP COLUMN age")
	if got := c.ShowCreateTable("a", "people"); strings.Contains(got, "chk_age") {
		t.Errorf("chk_age goes with age:\n%s", got)
	}
}

func TestRowLengthRefusals(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		code int
	}{
		{"CREATE TABLE x (a VARCHAR(65533) NOT NULL) CHARSET=latin1", 0},
		{"CREATE TABLE x (a VARCHAR(65534) NOT NULL) CHARSET=latin1", 1118},
		{"CREATE TABLE x (a VARCHAR(65533)) CHARSET=latin1", 1118},
		{"CREATE TABLE x (a VARCHAR(65532)) CHARSET=latin1", 0},
		{"CREATE TABLE x (a VARCHAR(65530), b TEXT) CHARSET=latin1", 1118},
		{"CREATE TABLE x (a VARCHAR(65520), b TEXT) CHARSET=latin1", 0},
		{"CREATE TABLE x (a VARCHAR(16384)) CHARSET=utf8mb4", 1074},
		{"CREATE TABLE x (a VARBINARY(65536))", 1074},
		{"CREATE TABLE x (a INT, b VARCHAR(65530) AS ('x') VIRTUAL) CHARSET=latin1", 1118},
		{"CREATE TABLE x (a CHAR(255), b VARCHAR(65000)) CHARSET=latin1", 0},
		{"CREATE TABLE y (a VARCHAR(10000), b VARCHAR(10000)) CHARSET=latin1; ALTER TABLE y CONVERT TO CHARACTER SET utf8mb4", 1118},
		{"CREATE TABLE y (a VARCHAR(20000)) CHARSET=latin1; ALTER TABLE y ADD COLUMN b VARCHAR(40000)", 0},
		{"SET SESSION sql_mode = ''; CREATE TABLE x (a VARCHAR(70000)) CHARSET=latin1", 0},
	} {
		c := familiesCatalog(t, "")
		if got := familiesCode(t, c, tc.sql); got != tc.code {
			t.Errorf("%s: got %d, want %d", tc.sql, got, tc.code)
		}
	}
}

func TestSpatialKeyRefusals(t *testing.T) {
	setup := "CREATE TABLE places (id INT PRIMARY KEY, loc POINT NOT NULL SRID 4326, SPATIAL KEY sp_loc (loc));"
	for _, tc := range []struct {
		sql  string
		code int
	}{
		{"ALTER TABLE places MODIFY loc POINT NOT NULL SRID 0", 3644},
		{"ALTER TABLE places MODIFY loc GEOMETRY NOT NULL", 3644},
		{"ALTER TABLE places MODIFY loc POINT NOT NULL SRID 4326 COMMENT 'x'", 0},
		{"CREATE TABLE p4 (g GEOMETRY NOT NULL SRID 0, SPATIAL KEY (g)) ENGINE = MyISAM", 0},
		{"CREATE TABLE p5 (g GEOMETRY NOT NULL SRID 3857, SPATIAL KEY (g)) ENGINE = MyISAM", 0},
		{"CREATE TABLE p6 (g GEOMETRY NOT NULL SRID 4326, SPATIAL KEY (g)) ENGINE = MyISAM", 1178},
		{"CREATE TABLE p7 (g GEOMETRY SRID 3857)", 0},
	} {
		c := familiesCatalog(t, setup)
		if got := familiesCode(t, c, tc.sql); got != tc.code {
			t.Errorf("%s: got %d, want %d", tc.sql, got, tc.code)
		}
	}
	c := familiesCatalog(t, setup)
	mustExec(t, c, "CREATE INDEX idx_loc ON places (loc)")
	if got := c.ShowCreateTable("a", "places"); !strings.Contains(got, "SPATIAL KEY `idx_loc` (`loc`)") {
		t.Errorf("a plain key on a NOT NULL geometry column is SPATIAL:\n%s", got)
	}
}

func TestOnlineDDLClauses(t *testing.T) {
	setup := "CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, s VARCHAR(10), e ENUM('a','b'), KEY idx_d (d));"
	for _, tc := range []struct {
		sql  string
		code int
	}{
		{"ALTER TABLE t ADD COLUMN x INT, ALGORITHM = INSTANT, LOCK = DEFAULT", 0},
		{"ALTER TABLE t ADD COLUMN x INT, LOCK = NONE", 0},
		{"ALTER TABLE t ALGORITHM = INSTANT, LOCK = NONE", 1221},
		{"ALTER TABLE t ALGORITHM = COPY, LOCK = NONE", 1846},
		{"ALTER TABLE t MODIFY c BIGINT, LOCK = SHARED", 0},
		{"ALTER TABLE t MODIFY c BIGINT, ALGORITHM = INSTANT", 1846},
		{"ALTER TABLE t MODIFY c INT NOT NULL, ALGORITHM = INPLACE", 0},
		{"ALTER TABLE t MODIFY c INT NOT NULL, ALGORITHM = INSTANT", 1845},
		{"ALTER TABLE t MODIFY c INT COMMENT 'x', ALGORITHM = INSTANT", 0},
		{"ALTER TABLE t MODIFY d INT NOT NULL, LOCK = NONE", 0},
		{"ALTER TABLE t MODIFY e ENUM('a','b','c'), ALGORITHM = INSTANT", 0},
		{"ALTER TABLE t MODIFY s VARCHAR(20), ALGORITHM = INPLACE, LOCK = NONE", 0},
		{"ALTER TABLE t MODIFY s VARCHAR(20), ALGORITHM = INSTANT", 1845},
		{"ALTER TABLE t MODIFY s VARCHAR(100), ALGORITHM = INPLACE", 1846},
		{"ALTER TABLE t MODIFY s VARCHAR(10) CHARACTER SET latin1, ALGORITHM = INPLACE", 1846},
		{"ALTER TABLE t DROP PRIMARY KEY, ALGORITHM = INSTANT", 1846},
		{"ALTER TABLE t DROP COLUMN c, ALGORITHM = INSTANT", 0},
		{"ALTER TABLE t ADD CONSTRAINT k CHECK (c > 0), LOCK = NONE", 1846},
		{"ALTER TABLE t ADD CONSTRAINT k CHECK (c > 0), ALGORITHM = COPY", 0},
		{"ALTER TABLE t ADD CONSTRAINT k CHECK (c > 0), MODIFY c BIGINT, ALGORITHM = INPLACE", 1845},
		{"ALTER TABLE t ALTER INDEX idx_d INVISIBLE, ALGORITHM = INPLACE, LOCK = NONE", 0},
		{"ALTER TABLE t ENGINE = InnoDB, LOCK = NONE", 0},
		{"ALTER TABLE t COMMENT = 'x', ALGORITHM = INSTANT", 1845},
		{"ALTER TABLE t RENAME TO t2, ALGORITHM = INSTANT", 0},
		{"ALTER TABLE t ADD COLUMN g INT AS (c + 1) STORED, ALGORITHM = INSTANT", 1845},
		{"ALTER TABLE t ADD COLUMN g INT AS (c + 1) VIRTUAL, ALGORITHM = INSTANT", 0},
		{"ALTER TABLE t ADD COLUMN x INT AUTO_INCREMENT UNIQUE, LOCK = NONE", 1846},
		{"ALTER TABLE t ADD FULLTEXT INDEX ft (s), LOCK = NONE", 1846},
		{"ALTER TABLE t ADD FULLTEXT INDEX ft (s), MODIFY c BIGINT, LOCK = NONE", 1846},
		{"ALTER TABLE t CONVERT TO CHARACTER SET latin1, ALGORITHM = INSTANT", 1846},
		{"ALTER TABLE t ALTER COLUMN c SET INVISIBLE, ALGORITHM = INSTANT", 0},
		{"ALTER TABLE t ADD COLUMN x INT UNIQUE, ALGORITHM = INSTANT", 1845},
		{"ALTER TABLE t MODIFY nope INT, ALGORITHM = INPLACE", 1054},
		{"CREATE INDEX i ON t (c) ALGORITHM = COPY LOCK = NONE", 1846},
	} {
		c := familiesCatalog(t, setup)
		if got := familiesCode(t, c, tc.sql); got != tc.code {
			t.Errorf("%s: got %d, want %d", tc.sql, got, tc.code)
		}
	}
}

func TestEventScheduleTimes(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		code int
		show string // the event's SHOW CREATE EVENT holds it; "" for no event
	}{
		{"CREATE EVENT e ON SCHEDULE EVERY 1 HOUR STARTS '2031-01-01 00:00:00' ENDS '2030-01-01 00:00:00' DO COMMIT", 1543, ""},
		{"CREATE EVENT e ON SCHEDULE EVERY 0 DAY DO COMMIT", 1542, ""},
		{"CREATE EVENT e ON SCHEDULE AT '2030-01-01' + INTERVAL 36 HOUR DO COMMIT", 0, "AT '2030-01-02 12:00:00'"},
		{"CREATE EVENT e ON SCHEDULE AT '2000-01-01 00:00:00' DO COMMIT", 0, ""},
		{"CREATE EVENT e ON SCHEDULE AT '2000-01-01 00:00:00' ON COMPLETION PRESERVE DO COMMIT", 0, "ON COMPLETION PRESERVE DISABLE"},
	} {
		c := familiesCatalog(t, "")
		if got := familiesCode(t, c, tc.sql); got != tc.code {
			t.Errorf("%s: got %d, want %d", tc.sql, got, tc.code)
			continue
		}
		got := c.ShowCreateEvent("a", "e")
		if (tc.show == "") != (got == "") || !strings.Contains(got, tc.show) {
			t.Errorf("%s: SHOW CREATE EVENT %q, want it to hold %q", tc.sql, got, tc.show)
		}
	}
}

func TestRenameRollbackRestoresChildReferences(t *testing.T) {
	// The child's snapshot, put back first, already references q; the
	// parent's undo must reach that copy.
	c := familiesCatalog(t, `CREATE TABLE p (id INT PRIMARY KEY);
		CREATE TABLE ch (pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id));`)
	if code := familiesCode(t, c, "RENAME TABLE p TO q, ch TO ch2, missing TO z"); code != ErrNoSuchTable {
		t.Fatalf("got %d, want %d", code, ErrNoSuchTable)
	}
	if got := c.ShowCreateTable("a", "ch"); !strings.Contains(got, "REFERENCES `p` (`id`)") {
		t.Errorf("ch after the failed rename:\n%s", got)
	}
}

func TestListColumnsDuplicatesByEachColumnsCollation(t *testing.T) {
	c := familiesCatalog(t, "")
	ok := `CREATE TABLE l (a VARCHAR(10) COLLATE utf8mb4_general_ci, b VARCHAR(10) COLLATE utf8mb4_bin)
		PARTITION BY LIST COLUMNS (a, b) (PARTITION p0 VALUES IN (('x', 'A')), PARTITION p1 VALUES IN (('x', 'a')))`
	if got := familiesCode(t, c, ok); got != 0 {
		t.Errorf("distinct under b's binary collation: got %d", got)
	}
	dup := `CREATE TABLE l2 (a VARCHAR(10) COLLATE utf8mb4_general_ci, b VARCHAR(10) COLLATE utf8mb4_bin)
		PARTITION BY LIST COLUMNS (a, b) (PARTITION p0 VALUES IN (('X', 'a')), PARTITION p1 VALUES IN (('x', 'a')))`
	if got := familiesCode(t, c, dup); got != 1495 {
		t.Errorf("equal under a's case-insensitive collation: got %d, want 1495", got)
	}
}

func TestViewKeepsSchemaAfterDerivedTable(t *testing.T) {
	// The derived table before a.t is resolved on its own scope; the CTE t
	// still shadows a.t after it.
	c := familiesCatalog(t, "CREATE TABLE t (id INT);")
	mustExec(t, c, "CREATE VIEW v AS WITH t AS (SELECT 2 AS id) SELECT a.t.id FROM (SELECT 1 AS x) d, a.t")
	const want = "with `t` as (select 2 AS `id`) select `a`.`t`.`id` AS `id` from ((select 1 AS `x`) `d` join `a`.`t`)"
	if got := c.GetDatabase("a").Views["v"].Definition; got != want {
		t.Errorf("v:\n got %s\nwant %s", got, want)
	}
}

func TestEventTimesInSessionTimeZone(t *testing.T) {
	// An AT one hour ahead of the session's own clock is in the future,
	// whatever UTC says.
	zone := time.FixedZone("", -8*3600)
	at := time.Now().In(zone).Add(time.Hour).Format("2006-01-02 15:04:05")
	c := familiesCatalog(t, "SET time_zone = '-08:00';")
	mustExec(t, c, "CREATE EVENT e ON SCHEDULE AT '"+at+"' DO COMMIT")
	if got := c.ShowCreateEvent("a", "e"); !strings.Contains(got, "AT '"+at+"'") || strings.Contains(got, "DISABLE") {
		t.Errorf("SHOW CREATE EVENT %q, want an enabled event AT '%s'", got, at)
	}
	ends := time.Now().In(zone).Add(time.Hour).Format("2006-01-02 15:04:05")
	if got := familiesCode(t, c, "CREATE EVENT e2 ON SCHEDULE EVERY 1 MINUTE ENDS '"+ends+"' DO COMMIT"); got != 0 {
		t.Errorf("ENDS one hour ahead of the session's clock: got %d", got)
	}
}

func TestViewKeepsSchemasOfSameNamedTables(t *testing.T) {
	// a.t and o.t are two tables of one name; each column keeps the schema
	// it names, and the view's own goes.
	c := familiesCatalog(t, "CREATE TABLE t (id INT); CREATE TABLE o.t (id INT);")
	mustExec(t, c, "CREATE VIEW v AS SELECT a.t.id AS aid, o.t.id AS oid FROM a.t JOIN o.t ON a.t.id = o.t.id")
	const want = "select `t`.`id` AS `aid`,`o`.`t`.`id` AS `oid` from (`t` join `o`.`t` on((`t`.`id` = `o`.`t`.`id`)))"
	if got := c.GetDatabase("a").Views["v"].Definition; got != want {
		t.Errorf("v:\n got %s\nwant %s", got, want)
	}
}

func TestCheckDependencyIgnoresStringLiterals(t *testing.T) {
	setup := "CREATE TABLE t (a INT, b VARCHAR(20), CONSTRAINT ck CHECK (b <> '`a`'));"
	for _, sql := range []string{"ALTER TABLE t DROP COLUMN a", "ALTER TABLE t RENAME COLUMN a TO a2"} {
		c := familiesCatalog(t, setup)
		if got := familiesCode(t, c, sql); got != 0 {
			t.Errorf("%s: got %d", sql, got)
		}
		if got := c.ShowCreateTable("a", "t"); !strings.Contains(got, "CONSTRAINT `ck`") {
			t.Errorf("%s keeps ck:\n%s", sql, got)
		}
	}
}

func TestConvertBesideNumericModifyRunsInPlace(t *testing.T) {
	c := familiesCatalog(t, "CREATE TABLE t (id INT) DEFAULT CHARACTER SET utf8mb4;")
	if got := familiesCode(t, c, "ALTER TABLE t CONVERT TO CHARACTER SET latin1, MODIFY id INT, ALGORITHM=INPLACE"); got != 0 {
		t.Errorf("got %d", got)
	}
}

func TestRoundThreeReviewCases(t *testing.T) {
	// Each was read on mysql:8.0.46.
	srs := `CREATE SPATIAL REFERENCE SYSTEM 1000001 NAME 'zz' DEFINITION 'GEOGCS["WGS 84",DATUM["World Geodetic System 1984",SPHEROID["WGS 84",6378137,298.257223563]],PRIMEM["Greenwich",0],UNIT["degree",0.017453292519943278],AXIS["Lat",NORTH],AXIS["Lon",EAST]]';`
	for _, tc := range []struct {
		setup, sql string
		code       int
	}{
		{"CREATE TABLE t (a INT, b INT, CONSTRAINT ck CHECK (a > 0));", "ALTER TABLE t DROP COLUMN a, DROP CHECK ck", 0},
		{"CREATE TABLE t (a INT, b INT);", "ALTER TABLE t ADD CONSTRAINT c2 CHECK (b > 0) NOT ENFORCED, ALGORITHM=INPLACE", 0},
		{"CREATE TABLE t (a INT, b INT);", "ALTER TABLE t ADD CONSTRAINT c2 CHECK (b > 0) NOT ENFORCED, LOCK=NONE", 0},
		{"CREATE TABLE t (a INT, b INT);", "ALTER TABLE t ADD COLUMN z INT CHECK (z > 0) NOT ENFORCED, ALGORITHM=INSTANT", 0},
		{subpartitioned, "ALTER TABLE s EXCHANGE PARTITION p0sp0 WITH TABLE plain", 0},
		{subpartitioned, "ALTER TABLE s TRUNCATE PARTITION p0sp0", 0},
		{subpartitioned, "ALTER TABLE s REBUILD PARTITION p0sp0", 1507},
		{subpartitioned, "ALTER TABLE s DROP PARTITION p0sp0", 1507},
		{"CREATE TABLE p (id INT PRIMARY KEY); CREATE TABLE o.p (id BIGINT PRIMARY KEY); CREATE TABLE ch (pid INT, FOREIGN KEY (pid) REFERENCES p (id));", "ALTER TABLE ch RENAME TO o.ch, MODIFY pid INT", 0},
		{srs, "CREATE TABLE g (x GEOMETRY SRID 1000001)", 0},
		{srs + " DROP SPATIAL REFERENCE SYSTEM 1000001;", "CREATE TABLE g (x GEOMETRY SRID 1000001)", 3548},
	} {
		c := familiesCatalog(t, tc.setup)
		if got := familiesCode(t, c, tc.sql); got != tc.code {
			t.Errorf("%s: got %d, want %d", tc.sql, got, tc.code)
		}
	}
	c := familiesCatalog(t, "CREATE TABLE p (id INT PRIMARY KEY); CREATE TABLE ch (pid INT, FOREIGN KEY (pid) REFERENCES p (id));")
	mustExec(t, c, "ALTER TABLE ch RENAME TO o.ch, MODIFY pid INT")
	if got := c.ShowCreateTable("o", "ch"); !strings.Contains(got, "REFERENCES `a`.`p` (`id`)") {
		t.Errorf("moved child keeps its parent:\n%s", got)
	}
}

const subpartitioned = `CREATE TABLE s (id INT NOT NULL, y INT NOT NULL) PARTITION BY RANGE (y) SUBPARTITION BY HASH (id) SUBPARTITIONS 2
	(PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20));
	CREATE TABLE plain (id INT NOT NULL, y INT NOT NULL);`
