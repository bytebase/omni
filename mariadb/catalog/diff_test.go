package catalog

import (
	"context"
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/internal/testcontainer"
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
	{"event body", "CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 1", "CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 2"},
	{"trigger body", "CREATE TABLE t (a INT); CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", "CREATE TABLE t (a INT); CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 2"},
	{"function body", "CREATE FUNCTION f(x INT) RETURNS INT DETERMINISTIC RETURN x", "CREATE FUNCTION f(x INT) RETURNS INT DETERMINISTIC RETURN x + 1"},
	{"function characteristic", "CREATE FUNCTION f(x INT) RETURNS INT DETERMINISTIC RETURN x", "CREATE FUNCTION f(x INT) RETURNS INT NOT DETERMINISTIC RETURN x"},
	{"application period", "CREATE TABLE t (s DATE NOT NULL, e DATE NOT NULL)", "CREATE TABLE t (s DATE NOT NULL, e DATE NOT NULL, PERIOD FOR p (s, e))"},
}

// TestOracle_DiffMatchesServer holds Diff to the server's own verdict on each
// pair: the two spellings are one schema exactly when the server stores the
// same tables, views and sequences for both, and Diff must then be empty in
// both directions and otherwise non-empty in both.
func TestOracle_DiffMatchesServer(t *testing.T) {
	conn, err := testcontainer.OpenDatabase(t).Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	database := testcontainer.DatabaseName(t)
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
			ca := loadCat(t, "CREATE DATABASE `bbcatalog`;\nUSE `bbcatalog`;\n"+p.a+";")
			cb := loadCat(t, "CREATE DATABASE `bbcatalog`;\nUSE `bbcatalog`;\n"+p.b+";")
			forward, backward := Diff(ca, cb), Diff(cb, ca)
			switch {
			case same && (!forward.IsEmpty() || !backward.IsEmpty()):
				t.Errorf("the server stores one schema, the diff is not empty:\n%s%s", describeDiff(forward), describeDiff(backward))
			case !same && (forward.IsEmpty() || backward.IsEmpty()):
				t.Errorf("the server stores two schemas, the diff is empty one way:\n%v\n%v", a, b)
			}
		})
	}
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
