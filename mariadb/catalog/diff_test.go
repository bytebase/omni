package catalog

import (
	"database/sql"
	"strings"
	"testing"
)

const dbDDL = "CREATE DATABASE app DEFAULT CHARSET=utf8mb4;\nUSE app;\n"

func loadCat(t *testing.T, sql string) *Catalog {
	t.Helper()
	c, err := LoadSDL(sql)
	if err != nil {
		t.Fatalf("LoadSDL failed for %q: %v", sql, err)
	}
	return c
}

// findTable locates a table entry in a SchemaDiff by name.
func findTable(d *SchemaDiff, name string) *TableDiffEntry {
	for i := range d.Tables {
		if strings.EqualFold(d.Tables[i].Name, name) {
			return &d.Tables[i]
		}
	}
	return nil
}

// diffPairs are two spellings of a schema each. The server decides whether
// they are one schema, by what it stores for each.
var diffPairs = []struct{ name, a, b string }{
	// Spellings MariaDB stores alike.
	{"int display width", "CREATE TABLE t (a INT, b BIGINT)", "CREATE TABLE t (a INT(11), b BIGINT(20))"},
	{"boolean", "CREATE TABLE t (a BOOLEAN DEFAULT TRUE)", "CREATE TABLE t (a TINYINT(1) DEFAULT 1)"},
	{"quoted number default", "CREATE TABLE t (a INT DEFAULT 1, d DECIMAL(10,2) DEFAULT 1.5)", "CREATE TABLE t (a INT DEFAULT '1', d DECIMAL(10,2) DEFAULT '1.50')"},
	{"default null", "CREATE TABLE t (a INT, b TEXT)", "CREATE TABLE t (a INT DEFAULT NULL, b TEXT DEFAULT NULL)"},
	{"timestamp null", "CREATE TABLE t (a TIMESTAMP)", "CREATE TABLE t (a TIMESTAMP NULL DEFAULT NULL)"},
	{"inherited charset", "CREATE TABLE t (a VARCHAR(10)) DEFAULT CHARSET=utf8mb4", "CREATE TABLE t (a VARCHAR(10) CHARACTER SET utf8mb4) DEFAULT CHARSET=utf8mb4"},
	{"utf8 alias", "CREATE TABLE t (a VARCHAR(10) CHARACTER SET utf8)", "CREATE TABLE t (a VARCHAR(10) CHARACTER SET utf8mb3)"},
	{"check parentheses", "CREATE TABLE t (a INT, CONSTRAINT c CHECK (a > 0))", "CREATE TABLE t (a INT, CONSTRAINT c CHECK (((a > 0))))"},
	{"engine spelling", "CREATE TABLE t (a INT)", "CREATE TABLE t (a INT) ENGINE=innodb"},
	{"now", "CREATE TABLE t (a DATETIME DEFAULT NOW())", "CREATE TABLE t (a DATETIME DEFAULT CURRENT_TIMESTAMP)"},
	{"sequence call", "CREATE SEQUENCE s; CREATE TABLE t (a INT DEFAULT NEXTVAL(s))", "CREATE SEQUENCE s; CREATE TABLE t (a INT DEFAULT (NEXT VALUE FOR s))"},
	{"view column list", "CREATE TABLE b (id INT); CREATE VIEW v (x) AS SELECT id FROM b", "CREATE TABLE b (id INT); CREATE VIEW v AS SELECT id AS x FROM b"},
	{"temporal default", "CREATE TABLE t (a TIME(2) DEFAULT '1:2:3')", "CREATE TABLE t (a TIME(2) DEFAULT '01:02:03.00')"},
	{"zerofill default", "CREATE TABLE t (a INT(5) ZEROFILL DEFAULT 3)", "CREATE TABLE t (a INT(5) UNSIGNED ZEROFILL DEFAULT '00003')"},
	{"json", "CREATE TABLE t (j JSON)", "CREATE TABLE t (j JSON DEFAULT NULL)"},
	{"long unique", "CREATE TABLE t (a TEXT, UNIQUE KEY u (a))", "CREATE TABLE t (a TEXT, UNIQUE KEY u (a) USING HASH)"},
	{"key order", "CREATE TABLE t (a INT, b INT, KEY ka (a), KEY kb (b))", "CREATE TABLE t (a INT, b INT, KEY kb (b), KEY ka (a))"},
	{"sequence defaults", "CREATE SEQUENCE s", "CREATE SEQUENCE s START WITH 1 INCREMENT BY 1 MINVALUE 1 CACHE 1000 NOCYCLE"},
	{"versioning spelled on a column", "CREATE TABLE t (a INT) WITH SYSTEM VERSIONING", "CREATE TABLE t (a INT WITH SYSTEM VERSIONING)"},
	// Changes MariaDB stores apart.
	{"column comment", "CREATE TABLE t (a INT)", "CREATE TABLE t (a INT COMMENT 'x')"},
	{"column type", "CREATE TABLE t (a INT)", "CREATE TABLE t (a BIGINT)"},
	{"default value", "CREATE TABLE t (a INT DEFAULT 1)", "CREATE TABLE t (a INT DEFAULT 2)"},
	{"nullability", "CREATE TABLE t (a INT)", "CREATE TABLE t (a INT NOT NULL)"},
	{"invisible", "CREATE TABLE t (a INT, b INT)", "CREATE TABLE t (a INT, b INT INVISIBLE)"},
	{"key kind", "CREATE TABLE t (a INT, KEY k (a))", "CREATE TABLE t (a INT, UNIQUE KEY k (a))"},
	{"ignored key", "CREATE TABLE t (a INT, KEY k (a))", "CREATE TABLE t (a INT, KEY k (a) IGNORED)"},
	{"descending key", "CREATE TABLE t (a INT, KEY k (a))", "CREATE TABLE t (a INT, KEY k (a DESC))"},
	{"using hash", "CREATE TABLE t (a INT, KEY k (a))", "CREATE TABLE t (a INT, KEY k (a) USING HASH)"},
	{"column check", "CREATE TABLE t (a INT CHECK (a > 0))", "CREATE TABLE t (a INT CHECK (a > 1))"},
	{"table check", "CREATE TABLE t (a INT, CONSTRAINT c CHECK (a > 0))", "CREATE TABLE t (a INT, CONSTRAINT c CHECK (a > 1))"},
	{"foreign key action", "CREATE TABLE p (id INT PRIMARY KEY); CREATE TABLE t (a INT, CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id))", "CREATE TABLE p (id INT PRIMARY KEY); CREATE TABLE t (a INT, CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id) ON DELETE CASCADE)"},
	{"table comment", "CREATE TABLE t (a INT)", "CREATE TABLE t (a INT) COMMENT 'x'"},
	{"table engine", "CREATE TABLE t (a INT)", "CREATE TABLE t (a INT) ENGINE=MyISAM"},
	{"table option", "CREATE TABLE t (a INT)", "CREATE TABLE t (a INT) STATS_PERSISTENT=0"},
	{"table collation", "CREATE TABLE t (a VARCHAR(10)) COLLATE utf8mb4_bin", "CREATE TABLE t (a VARCHAR(10)) COLLATE utf8mb4_unicode_ci"},
	{"system versioning", "CREATE TABLE t (a INT)", "CREATE TABLE t (a INT) WITH SYSTEM VERSIONING"},
	{"column without versioning", "CREATE TABLE t (a INT, b INT) WITH SYSTEM VERSIONING", "CREATE TABLE t (a INT, b INT WITHOUT SYSTEM VERSIONING) WITH SYSTEM VERSIONING"},
	{"sequence increment", "CREATE SEQUENCE s", "CREATE SEQUENCE s INCREMENT BY 2"},
	{"sequence comment", "CREATE SEQUENCE s", "CREATE SEQUENCE s COMMENT 'x'"},
	{"view body", "CREATE TABLE b (id INT); CREATE VIEW v AS SELECT id FROM b", "CREATE TABLE b (id INT); CREATE VIEW v AS SELECT id + 1 AS id FROM b"},
	{"view algorithm", "CREATE TABLE b (id INT); CREATE VIEW v AS SELECT id FROM b", "CREATE TABLE b (id INT); CREATE ALGORITHM=MERGE VIEW v AS SELECT id FROM b"},
	{"partitions", "CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 2", "CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 4"},
	{"foreign key restrict", "CREATE TABLE p (id INT PRIMARY KEY); CREATE TABLE t (a INT, CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id))", "CREATE TABLE p (id INT PRIMARY KEY); CREATE TABLE t (a INT, CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id) ON DELETE RESTRICT ON UPDATE RESTRICT)"},
	{"foreign key no action", "CREATE TABLE p (id INT PRIMARY KEY); CREATE TABLE t (a INT, CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id))", "CREATE TABLE p (id INT PRIMARY KEY); CREATE TABLE t (a INT, CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id) ON DELETE NO ACTION)"},
	{"event interval", "CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 1", "CREATE EVENT e ON SCHEDULE EVERY 24 HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1"},
	{"event status", "CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 1", "CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DISABLE DO SELECT 1"},
	{"comment spelling the counter", "CREATE TABLE t (a INT) COMMENT 'counter AUTO_INCREMENT=10'", "CREATE TABLE t (a INT) COMMENT 'counter AUTO_INCREMENT=20'"},
	{"event interval expression", "CREATE EVENT e ON SCHEDULE EVERY LENGTH('01') HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1", "CREATE EVENT e ON SCHEDULE EVERY LENGTH('1') HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1"},
	{"event interval string spacing", "CREATE EVENT e ON SCHEDULE EVERY LENGTH('a  b') HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1", "CREATE EVENT e ON SCHEDULE EVERY LENGTH('a b') HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1"},
	{"event schedule spacing", "CREATE EVENT e ON SCHEDULE EVERY   2  HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1", "CREATE EVENT e ON SCHEDULE EVERY 2 HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1"},
	{"event body", "CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 1", "CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 2"},
	{"trigger body", "CREATE TABLE t (a INT); CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", "CREATE TABLE t (a INT); CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 2"},
	{"function body", "CREATE FUNCTION f(x INT) RETURNS INT DETERMINISTIC RETURN x", "CREATE FUNCTION f(x INT) RETURNS INT DETERMINISTIC RETURN x + 1"},
	{"function characteristic", "CREATE FUNCTION f(x INT) RETURNS INT DETERMINISTIC RETURN x", "CREATE FUNCTION f(x INT) RETURNS INT NOT DETERMINISTIC RETURN x"},
	{"generated storage with a key", "CREATE TABLE t (a INT, g INT AS (a + 1) VIRTUAL, KEY k (g))", "CREATE TABLE t (a INT, g INT AS (a + 1) STORED, KEY k (g))"},
	{"generated columns on generated columns", "CREATE TABLE t (a INT, g1 INT AS (a + 1) VIRTUAL, g2 INT AS (g1 + 1) VIRTUAL)", "CREATE TABLE t (a INT)"},
	{"primary key without overlaps", "CREATE TABLE x (i INT)", "CREATE TABLE x (i INT); CREATE TABLE t (id INT, s DATE NOT NULL, e DATE NOT NULL, PERIOD FOR p (s, e), PRIMARY KEY (id, p WITHOUT OVERLAPS))"},
	{"every column replaced", "CREATE TABLE t (a INT, b INT)", "CREATE TABLE t (c INT, d INT)"},
	{"key a foreign key needs", "CREATE TABLE p (id INT PRIMARY KEY); CREATE TABLE t (a INT, b INT, KEY k (a), CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id))", "CREATE TABLE p (id INT PRIMARY KEY); CREATE TABLE t (a INT, b INT, KEY k (a, b), CONSTRAINT f FOREIGN KEY (a) REFERENCES p (id))"},
	{"application period", "CREATE TABLE t (s DATE NOT NULL, e DATE NOT NULL)", "CREATE TABLE t (s DATE NOT NULL, e DATE NOT NULL, PERIOD FOR p (s, e))"},
	{"generated column leaving a dropped column", "CREATE TABLE t (x INT, y INT, g INT AS (x + 1) VIRTUAL)", "CREATE TABLE t (y INT, g INT AS (y + 1) VIRTUAL)"},
	{"primary key and partitioning moved together", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a)) PARTITION BY HASH (a) PARTITIONS 2", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (b)) PARTITION BY HASH (b) PARTITIONS 2"},
	{"foreign key columns retyped", "CREATE TABLE p (id INT PRIMARY KEY); CREATE TABLE ch (pid INT, CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))", "CREATE TABLE p (id BIGINT PRIMARY KEY); CREATE TABLE ch (pid BIGINT, CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))"},
	{"foreign key over an unchanged table", "CREATE TABLE p (id INT PRIMARY KEY, n INT); CREATE TABLE ch (pid INT, CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))", "CREATE TABLE p (id INT PRIMARY KEY, n BIGINT); CREATE TABLE ch (pid INT, CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))"},
	{"versioning with its period columns", "CREATE TABLE t (x INT)", "CREATE TABLE t (x INT, s TIMESTAMP(6) AS ROW START, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s, e)) WITH SYSTEM VERSIONING"},
	{"generated column reading an added one", "CREATE TABLE t (x INT, g INT AS (x + 1) VIRTUAL)", "CREATE TABLE t (x INT, h INT AS (x + 1) VIRTUAL, g INT AS (h + 1) VIRTUAL)"},
	{"referenced key replaced", "CREATE TABLE p (id INT, KEY k (id)); CREATE TABLE ch (pid INT, CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))", "CREATE TABLE p (id INT, KEY k (id) COMMENT 'c'); CREATE TABLE ch (pid INT, CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))"},
	{"period moved under an overlap key", "CREATE TABLE t (id INT, s DATE NOT NULL, e DATE NOT NULL, s2 DATE NOT NULL, e2 DATE NOT NULL, PERIOD FOR p (s, e), UNIQUE KEY k (id, p WITHOUT OVERLAPS))", "CREATE TABLE t (id INT, s DATE NOT NULL, e DATE NOT NULL, s2 DATE NOT NULL, e2 DATE NOT NULL, PERIOD FOR p (s2, e2), UNIQUE KEY k (id, p WITHOUT OVERLAPS))"},
	{"sequence limits leaving out the start", "CREATE SEQUENCE s", "CREATE SEQUENCE s START WITH 100 MINVALUE 10"},
}

// TestOracle_DiffMatchesServer holds Diff to the server's own verdict on each
// pair: the two spellings are one schema exactly when the server stores the
// same tables, views and sequences for both, and Diff must then be empty in
// both directions and otherwise non-empty in both.
func TestOracle_DiffMatchesServer(t *testing.T) {
	forEachServer(t, func(t *testing.T, srv oracleServer, conn *sql.Conn, database string) {
		stored := func(t *testing.T, schema string) map[string]string {
			resetDatabase(t, conn, database)
			applySchema(t, conn, schema)
			return storedObjects(t, conn, true)
		}
		for _, p := range diffPairs {
			t.Run(p.name, func(t *testing.T) {
				a, b := stored(t, p.a), stored(t, p.b)
				same := len(a) == len(b)
				for k, v := range a {
					if b[k] != v {
						same = false
					}
				}
				ca, err := LoadSDLWithVersion("CREATE DATABASE `bbcatalog`;\nUSE `bbcatalog`;\n"+p.a+";", srv.version)
				if err != nil {
					t.Fatal(err)
				}
				cb, err := LoadSDLWithVersion("CREATE DATABASE `bbcatalog`;\nUSE `bbcatalog`;\n"+p.b+";", srv.version)
				if err != nil {
					t.Fatal(err)
				}
				forward, backward := Diff(ca, cb), Diff(cb, ca)
				switch {
				case same && (!forward.IsEmpty() || !backward.IsEmpty()):
					t.Errorf("the server stores one schema, the diff is not empty:\n%s%s", describeDiff(forward), describeDiff(backward))
				case !same && (forward.IsEmpty() || backward.IsEmpty()):
					t.Errorf("the server stores two schemas, the diff is empty one way:\n%v\n%v", a, b)
				}
			})
		}
	})
}

func TestLoadSDLRestoresSessionPerStatement(t *testing.T) {
	c := loadCat(t, dbDDL+"SET explicit_defaults_for_timestamp = OFF; CREATE TABLE t (a TIMESTAMP); SET explicit_defaults_for_timestamp = ON; CREATE TABLE u (a TIMESTAMP);")
	if got := c.ShowCreateTable("app", "t"); !strings.Contains(got, "`a` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()") {
		t.Errorf("t loaded without the session its SET gave it:\n%s", got)
	}
	if got := c.ShowCreateTable("app", "u"); !strings.Contains(got, "`a` timestamp NULL DEFAULT NULL") {
		t.Errorf("u loaded without the session its SET gave it:\n%s", got)
	}
}

func TestCanonicalTableOptionsKeepsCommentSpacing(t *testing.T) {
	a := loadCat(t, dbDDL+"CREATE TABLE t (a INT) COMMENT 'a  b';")
	b := loadCat(t, dbDDL+"CREATE TABLE t (a INT) COMMENT 'a b';")
	if Diff(a, b).IsEmpty() {
		t.Error("comments differing in spacing must diff")
	}
}
