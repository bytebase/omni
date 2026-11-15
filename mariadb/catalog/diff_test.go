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
	{"comment spelling the counter", "CREATE TABLE t (a INT) COMMENT 'counter AUTO_INCREMENT=10'", "CREATE TABLE t (a INT) COMMENT 'counter AUTO_INCREMENT=20'"},
	{"event interval expression", "CREATE EVENT e ON SCHEDULE EVERY LENGTH('01') HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1", "CREATE EVENT e ON SCHEDULE EVERY LENGTH('1') HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1"},
	{"event interval string spacing", "CREATE EVENT e ON SCHEDULE EVERY LENGTH('a  b') HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1", "CREATE EVENT e ON SCHEDULE EVERY LENGTH('a b') HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1"},
	{"event schedule spacing", "CREATE EVENT e ON SCHEDULE EVERY   2  HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1", "CREATE EVENT e ON SCHEDULE EVERY 2 HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1"},
	{"generated column check", "CREATE TABLE t (a INT, b INT AS (a + 1) VIRTUAL CHECK (b > 0))", "CREATE TABLE t (a INT, b INT AS (a + 1) VIRTUAL CHECK (b > 1))"},
	{"event decimal interval", "CREATE EVENT e ON SCHEDULE EVERY 1.05 HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1", "CREATE EVENT e ON SCHEDULE EVERY 1.5 HOUR STARTS '2030-01-01 00:00:00' DO SELECT 1"},
	{"table collation with inheriting columns", "CREATE TABLE t (a VARCHAR(10)) COLLATE utf8mb4_bin", "CREATE TABLE t (a VARCHAR(10)) COLLATE utf8mb4_unicode_ci"},
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
	{"flipped column with readers", "CREATE TABLE t (x INT, h INT AS (x + 1) VIRTUAL, g INT AS (h + 1) VIRTUAL, k INT AS (g * 2) VIRTUAL, CONSTRAINT ck CHECK (h > 0), KEY kg (g))", "CREATE TABLE t (x INT, h INT AS (x + 1) STORED, g INT AS (h + 1) VIRTUAL, k INT AS (g * 2) VIRTUAL, CONSTRAINT ck CHECK (h > 0), KEY kg (g))"},
	{"many generated columns and a key dropped", "CREATE TABLE t (x INT, a INT AS (x + 1) VIRTUAL, b INT AS (a + 1) VIRTUAL, c INT AS (x + 1) VIRTUAL, d INT AS (x + 1) VIRTUAL, e INT AS (x + 1) VIRTUAL, f INT AS (x + 1) VIRTUAL, KEY k (a))", "CREATE TABLE t (x INT)"},
	{"flipped column with a modified reader", "CREATE TABLE t (x INT, h INT AS (x + 1) VIRTUAL, g INT AS (h + 1) VIRTUAL)", "CREATE TABLE t (x INT, h INT AS (x + 1) STORED, g INT AS (h + 2) VIRTUAL)"},
	{"versioning dropped with indexed period columns", "CREATE TABLE t (x INT, s TIMESTAMP(6) AS ROW START, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s, e), KEY k (e)) WITH SYSTEM VERSIONING", "CREATE TABLE t (x INT)"},
	{"versioning dropped with its period columns kept", "CREATE TABLE t (x INT, s TIMESTAMP(6) AS ROW START, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s, e)) WITH SYSTEM VERSIONING", "CREATE TABLE t (x INT, s TIMESTAMP(6) NULL, e TIMESTAMP(6) NULL)"},
	{"removed reader of a rebuilt column", "CREATE TABLE t (x INT, h INT AS (x + 1) VIRTUAL, g INT AS (h + 1) VIRTUAL, k INT AS (g + 1) VIRTUAL)", "CREATE TABLE t (x INT, h INT AS (x + 1) STORED, g INT AS (h + 1) VIRTUAL)"},
	{"replaced column read by a removed one", "CREATE TABLE t (base INT, a INT AS (base) VIRTUAL, z INT AS (a) VIRTUAL)", "CREATE TABLE t (base INT, a INT AS (base) STORED)"},
	{"every column replaced", "CREATE TABLE t (a INT)", "CREATE TABLE t (a INT AS (1) VIRTUAL, b INT)"},
	{"the only column replaced", "CREATE TABLE t (a INT)", "CREATE TABLE t (a INT AS (1) VIRTUAL)"},
	{"page checksum kept", "CREATE TABLE t (a INT) PAGE_CHECKSUM=1", "CREATE TABLE t (a INT, b INT) PAGE_CHECKSUM=1"},
	{"system period moved", "CREATE TABLE t (x INT, s TIMESTAMP(6) AS ROW START, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s, e)) WITH SYSTEM VERSIONING", "CREATE TABLE t (x INT, s2 TIMESTAMP(6) AS ROW START, e2 TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s2, e2)) WITH SYSTEM VERSIONING"},
	{"two dependent replacements", "CREATE TABLE t (x INT, a INT AS (x + 1) VIRTUAL, b INT AS (a + 1) VIRTUAL)", "CREATE TABLE t (x INT, a INT AS (x + 1) STORED, b INT AS (a + 1) STORED)"},
	{"column check reading a later column", "CREATE TABLE t (x INT)", "CREATE TABLE t (x INT, a INT CHECK (a > b), b INT)"},
	{"foreign key over a replaced column", "CREATE TABLE p (id INT PRIMARY KEY); CREATE TABLE ch (x INT, pid INT AS (x + 1) STORED, KEY k (pid), CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))", "CREATE TABLE p (id INT PRIMARY KEY); CREATE TABLE ch (x INT, pid INT, KEY k (pid), CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))"},
	{"versioning with a column left out", "CREATE TABLE t (a INT, c INT)", "CREATE TABLE t (a INT, c INT WITHOUT SYSTEM VERSIONING, b INT WITHOUT SYSTEM VERSIONING) WITH SYSTEM VERSIONING"},
	{"period bound types changed", "CREATE TABLE t (id INT, s DATE NOT NULL, e DATE NOT NULL, PERIOD FOR app_time (s, e))", "CREATE TABLE t (id INT, s DATETIME NOT NULL, e DATETIME NOT NULL, PERIOD FOR app_time (s, e))"},
	{"every column replaced with a reader", "CREATE TABLE t (a INT)", "CREATE TABLE t (a INT AS (1) VIRTUAL, b INT AS (a) VIRTUAL)"},
	{"new table with page checksum and a key", "CREATE TABLE x (i INT)", "CREATE TABLE x (i INT); CREATE TABLE t (a INT, KEY k (a)) PAGE_CHECKSUM=1"},
	{"foreign key losing a column", "CREATE TABLE p (a INT, b INT, PRIMARY KEY (a, b)); CREATE TABLE ch (a INT, b INT, CONSTRAINT f FOREIGN KEY (a, b) REFERENCES p (a, b))", "CREATE TABLE p (a INT, b INT, PRIMARY KEY (a, b)); CREATE TABLE ch (a INT, b INT, CONSTRAINT f FOREIGN KEY (a) REFERENCES p (a))"},
	{"column checks reading each other", "CREATE TABLE t (x INT)", "CREATE TABLE t (x INT, a INT CHECK (a > b), b INT CHECK (b < a))"},
	{"checks reading through to a generated column", "CREATE TABLE t (x INT)", "CREATE TABLE t (x INT, a INT CHECK (a > b), g INT AS (x + 1) VIRTUAL, b INT CHECK (b > g))"},
	{"every column dropped in one statement", "CREATE TABLE t (a INT CHECK (a > b), b INT CHECK (b < a))", "CREATE TABLE t (x INT)"},
	{"kept column checking an added generated one", "CREATE TABLE t (x INT, a INT)", "CREATE TABLE t (x INT, b INT AS (x + 1) VIRTUAL, a INT CHECK (a > b))"},
	{"system period with one bound kept", "CREATE TABLE t (x INT, s TIMESTAMP(6) AS ROW START, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s, e)) WITH SYSTEM VERSIONING", "CREATE TABLE t (x INT, s2 TIMESTAMP(6) AS ROW START, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s2, e)) WITH SYSTEM VERSIONING"},
	{"engine change dropping auto_increment", "CREATE TABLE t (g INT NOT NULL, id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (g, id)) ENGINE=MyISAM", "CREATE TABLE t (g INT NOT NULL, id INT NOT NULL, PRIMARY KEY (g, id)) ENGINE=InnoDB"},
	{"trigger order", "CREATE TABLE t (x INT)", "CREATE TABLE t (x INT); CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x * 2; CREATE TRIGGER z BEFORE INSERT ON t FOR EACH ROW PRECEDES a SET NEW.x = NEW.x + 1"},
	{"fresh columns reading each other", "CREATE TABLE t (a INT)", "CREATE TABLE t (b INT AS (c + 1) VIRTUAL, c INT)"},
	{"every column replaced, a reader among them", "CREATE TABLE t (m INT, b INT AS (m + 1) VIRTUAL)", "CREATE TABLE t (m INT AS (1) STORED, b INT AS (m + 1) STORED)"},
	{"fresh column reading a new sequence", "CREATE TABLE t (a INT)", "CREATE SEQUENCE s; CREATE TABLE t (b INT DEFAULT (NEXT VALUE FOR s))"},
	{"myisam auto column inside a key", "CREATE TABLE x (i INT)", "CREATE TABLE x (i INT); CREATE TABLE t (a INT, b INT AUTO_INCREMENT, KEY k (a, b)) ENGINE=MyISAM"},
	{"auto column replaced", "CREATE TABLE t (a INT AUTO_INCREMENT PRIMARY KEY)", "CREATE TABLE t (b INT AUTO_INCREMENT PRIMARY KEY)"},
	{"ignored unique key after the implicit primary key", "CREATE TABLE x (i INT)", "CREATE TABLE x (i INT); CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY z (b), UNIQUE KEY aa (a) IGNORED)"},
	{"modified reader of a replaced column", "CREATE TABLE t (x INT, h INT, g INT AS (x + 1) VIRTUAL)", "CREATE TABLE t (x INT, h INT AS (x + 1) STORED, g INT AS (h + 1) VIRTUAL)"},
	{"system period moved with the kept bound changed", "CREATE TABLE t (x INT, s TIMESTAMP(6) AS ROW START, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s, e)) WITH SYSTEM VERSIONING", "CREATE TABLE t (x INT, s2 TIMESTAMP(6) AS ROW START, e TIMESTAMP(6) AS ROW END COMMENT 'changed', PERIOD FOR SYSTEM_TIME (s2, e)) WITH SYSTEM VERSIONING"},
	{"view replaced by a sequence", "CREATE VIEW s AS SELECT 1 AS n", "CREATE SEQUENCE s"},
	{"auto column under two unique keys", "CREATE TABLE x (i INT)", "CREATE TABLE x (i INT); CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, UNIQUE KEY z (a), UNIQUE KEY aa (a, b) IGNORED)"},
	{"versioned table losing every column", "CREATE TABLE t (a INT, s TIMESTAMP(6) AS ROW START, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s, e)) WITH SYSTEM VERSIONING", "CREATE TABLE t (b INT)"},
	{"auto column keyed with grouped columns", "CREATE TABLE t (x INT)", "CREATE TABLE t (x INT, id INT AUTO_INCREMENT, a INT CHECK (a > b), b INT CHECK (b < a), KEY k (id, a))"},
	{"system period moved under a key on the kept bound", "CREATE TABLE t (x INT, s TIMESTAMP(6) AS ROW START, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s, e), KEY k (e)) WITH SYSTEM VERSIONING", "CREATE TABLE t (x INT, s2 TIMESTAMP(6) AS ROW START, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s2, e), KEY k (e)) WITH SYSTEM VERSIONING"},
	{"system period moved under a primary key on the kept bound", "CREATE TABLE t (id INT NOT NULL, s TIMESTAMP(6) AS ROW START, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s, e), PRIMARY KEY (id, e)) WITH SYSTEM VERSIONING", "CREATE TABLE t (id INT NOT NULL, s2 TIMESTAMP(6) AS ROW START, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s2, e), PRIMARY KEY (id, e)) WITH SYSTEM VERSIONING"},
	{"postponed drop reading another removed column", "CREATE TABLE t (x INT, d INT AS (x + 1) VIRTUAL, g INT AS (d + 1) VIRTUAL)", "CREATE TABLE t (y INT, g INT AS (y + 1) VIRTUAL)"},
	{"auto key over a column checking a generated one", "CREATE TABLE t (x INT)", "CREATE TABLE t (x INT, id INT AUTO_INCREMENT, g INT AS (x + 1) VIRTUAL, a INT CHECK (a > g), KEY k (id, a))"},
	{"last versioned column replaced", "CREATE TABLE t (a INT, b INT WITHOUT SYSTEM VERSIONING) WITH SYSTEM VERSIONING", "CREATE TABLE t (b INT WITHOUT SYSTEM VERSIONING, c INT) WITH SYSTEM VERSIONING"},
	{"kept column checking a replaced one", "CREATE TABLE t (x INT, a INT, b INT CHECK (b > a))", "CREATE TABLE t (x INT, a INT AS (1) STORED, b INT CHECK (b > a))"},
	{"view replaced by a sequence a new column reads", "CREATE TABLE t (a INT); CREATE VIEW s AS SELECT 1 AS n", "CREATE SEQUENCE s; CREATE TABLE t (b INT DEFAULT (NEXT VALUE FOR s))"},
	{"every column replaced, checks reading each other", "CREATE TABLE t (a INT CHECK (a > b), b INT CHECK (b < a))", "CREATE TABLE t (a INT AS (2) STORED CHECK (a > b), b INT AS (1) STORED CHECK (b < a))"},
	{"only column replaced by a group", "CREATE TABLE t (a INT)", "CREATE TABLE t (a INT AS (1) STORED CHECK (a >= b), b INT CHECK (b >= a))"},
	{"triggers in creation order", "CREATE TABLE t (x INT)", "CREATE TABLE t (x INT); CREATE TRIGGER z BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x + 1; CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x * 2"},
	{"auto key over a generated chain", "CREATE TABLE t (x INT)", "CREATE TABLE t (x INT, id INT AUTO_INCREMENT, g1 INT AS (x + 1) VIRTUAL, g2 INT AS (g1 + 1) VIRTUAL, KEY k (id, g2))"},
	{"versioning-attributed column a check reads", "CREATE TABLE t (x INT)", "CREATE TABLE t (x INT, b INT WITHOUT SYSTEM VERSIONING, a INT CHECK (a > b)) WITH SYSTEM VERSIONING"},
	{"kept default reading a removed column", "CREATE TABLE t (id INT, v INT DEFAULT (id + 1))", "CREATE TABLE t (v INT DEFAULT 0)"},
	{"last visible column replaced", "CREATE TABLE t (a INT, b INT INVISIBLE)", "CREATE TABLE t (b INT INVISIBLE, c INT)"},
	{"recreated trigger keeps its place", "CREATE TABLE t (x INT); CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x * 2; CREATE TRIGGER b BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x + 1", "CREATE TABLE t (x INT); CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x * 3; CREATE TRIGGER b BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x + 1"},
	{"replaced column read by a default", "CREATE TABLE t (x INT, h INT, v INT DEFAULT (h + 1))", "CREATE TABLE t (x INT, h INT AS (x + 1) STORED, v INT DEFAULT 0)"},
	{"visibility moved between columns", "CREATE TABLE t (a INT, z INT INVISIBLE)", "CREATE TABLE t (a INT INVISIBLE, z INT)"},
	{"versioning moved between columns", "CREATE TABLE t (a INT, b INT WITHOUT SYSTEM VERSIONING) WITH SYSTEM VERSIONING", "CREATE TABLE t (a INT WITHOUT SYSTEM VERSIONING, b INT) WITH SYSTEM VERSIONING"},
	{"hidden column releasing a default", "CREATE TABLE t (id INT, v INT DEFAULT (id + 1), z INT INVISIBLE)", "CREATE TABLE t (v INT INVISIBLE DEFAULT 0, z INT)"},
	{"primary key replaced with an ignored key kept", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a), UNIQUE KEY u (b) IGNORED)", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (b), UNIQUE KEY u (b) IGNORED)"},
	{"auto key absorbing a read column", "CREATE TABLE t (x INT)", "CREATE TABLE t (x INT, id INT AUTO_INCREMENT, b INT, g INT AS (b + 1) VIRTUAL, KEY k (id, b))"},
	{"partition column replaced", "CREATE TABLE t (x INT, a INT) PARTITION BY HASH (a) PARTITIONS 2", "CREATE TABLE t (x INT, a INT AS (x) STORED) PARTITION BY HASH (a) PARTITIONS 2"},
	{"engine change moving auto to front", "CREATE TABLE t (g INT NOT NULL, id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (g, id)) ENGINE=MyISAM", "CREATE TABLE t (g INT NOT NULL, id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id, g)) ENGINE=InnoDB"},
	{"foreign key over a full key beside a prefix", "CREATE TABLE p (a VARCHAR(100), KEY kp (a)); CREATE TABLE ch (a VARCHAR(100), KEY short_key (a(10)), KEY full_key (a), CONSTRAINT f FOREIGN KEY (a) REFERENCES p (a))", "CREATE TABLE p (a VARCHAR(100), KEY kp (a)); CREATE TABLE ch (a VARCHAR(100), KEY short_key (a(10)), KEY replacement_key (a), CONSTRAINT f FOREIGN KEY (a) REFERENCES p (a))"},
	{"last visible column replaced not dropped", "CREATE TABLE t (a INT, b INT INVISIBLE)", "CREATE TABLE t (a INT AS (1) STORED, b INT INVISIBLE)"},
	{"auto column replaced and grouped", "CREATE TABLE t (a INT AUTO_INCREMENT PRIMARY KEY)", "CREATE TABLE t (a INT AS (1) STORED)"},
	{"auto key absorbing a generated chain", "CREATE TABLE t (x INT)", "CREATE TABLE t (x INT, id INT AUTO_INCREMENT, g1 INT AS (x + 1) VIRTUAL, g2 INT AS (g1 + 1) VIRTUAL, g3 INT AS (g2 + 1) VIRTUAL, KEY k (id, g2))"},
	{"trigger before an unchanged sibling", "CREATE TABLE t (x INT); CREATE TRIGGER z BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x + 1; CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x * 2", "CREATE TABLE t (x INT); CREATE TRIGGER z BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x + 2; CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x * 2"},
	{"auto key absorbing a modified reader", "CREATE TABLE t (x INT, g INT AS (x + 1) VIRTUAL)", "CREATE TABLE t (x INT, id INT AUTO_INCREMENT, b INT, g INT AS (b + 1) VIRTUAL, KEY k (id, b))"},
	{"versioning replacing a generated column", "CREATE TABLE t (x INT, b INT AS (x + 1) VIRTUAL)", "CREATE TABLE t (x INT, b INT WITHOUT SYSTEM VERSIONING) WITH SYSTEM VERSIONING"},
	{"partition column type changed", "CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 2", "CREATE TABLE t (a VARCHAR(10)) PARTITION BY KEY (a) PARTITIONS 2"},
	{"merge delaying a modify that releases a drop", "CREATE TABLE t (x INT, d INT, g INT AS (d + 1) VIRTUAL)", "CREATE TABLE t (x INT, id INT AUTO_INCREMENT, b INT, g INT AS (b + 1) VIRTUAL, KEY k (id, b))"},
	{"recreated trigger among a reordered group", "CREATE TABLE t (x INT); CREATE TRIGGER c BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x * 2; CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW PRECEDES c SET NEW.x = NEW.x + 1; CREATE TRIGGER b BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x + 3", "CREATE TABLE t (x INT); CREATE TRIGGER c BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x * 4; CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW PRECEDES c SET NEW.x = NEW.x + 1; CREATE TRIGGER b BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x + 3"},
	{"engine fold keeping another option", "CREATE TABLE t (g INT NOT NULL, id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (g, id)) ENGINE=MyISAM", "CREATE TABLE t (g INT NOT NULL, id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id, g)) ENGINE=InnoDB COMMENT='new'"},
	{"two new triggers after a kept one", "CREATE TABLE t (x INT); CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x + 1", "CREATE TABLE t (x INT); CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x + 1; CREATE TRIGGER b BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x * 2; CREATE TRIGGER c BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x + 3"},
	{"merge delaying a transitive drop", "CREATE TABLE t (x INT, d INT AS (x + 1) VIRTUAL, g INT AS (d + 1) VIRTUAL)", "CREATE TABLE t (id INT AUTO_INCREMENT PRIMARY KEY, b INT, g INT AS (b + 1) VIRTUAL, KEY k (id, b))"},
	{"visible column reading an added invisible one", "CREATE TABLE t (a INT, hidden INT INVISIBLE)", "CREATE TABLE t (b INT CHECK (b > c), c INT INVISIBLE, hidden INT INVISIBLE)"},
	{"engine fold keeping partitions", "CREATE TABLE t (g INT NOT NULL, id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (g, id)) ENGINE=MyISAM PARTITION BY HASH (g) PARTITIONS 2", "CREATE TABLE t (g INT NOT NULL, id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id, g)) ENGINE=InnoDB PARTITION BY HASH (g) PARTITIONS 3"},
	{"two new triggers before a kept one", "CREATE TABLE t (x INT); CREATE TRIGGER c BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x + 3", "CREATE TABLE t (x INT); CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW PRECEDES c SET NEW.x = NEW.x + 1; CREATE TRIGGER b BEFORE INSERT ON t FOR EACH ROW PRECEDES c SET NEW.x = NEW.x * 2; CREATE TRIGGER c BEFORE INSERT ON t FOR EACH ROW SET NEW.x = NEW.x + 3"},
	{"new table auto column with an ignored key", "CREATE TABLE x (i INT)", "CREATE TABLE x (i INT); CREATE TABLE t (a INT NOT NULL, b INT AUTO_INCREMENT, UNIQUE KEY u (a), UNIQUE KEY v (b) IGNORED)"},
	{"timestamps without defaults", "CREATE TABLE x (i INT)", "CREATE TABLE x (i INT, ts TIMESTAMP NOT NULL); CREATE TABLE t (u TIMESTAMP NOT NULL)"},
	{"sequence limits leaving out the start", "CREATE SEQUENCE s", "CREATE SEQUENCE s START WITH 100 MINVALUE 10"},
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

func TestDiffLeavesTemporarySequencesOut(t *testing.T) {
	a := loadCat(t, dbDDL+"CREATE TABLE t (a INT);")
	b := loadCat(t, dbDDL+"CREATE TABLE t (a INT);")
	if _, err := b.Exec("CREATE TEMPORARY SEQUENCE s;", nil); err != nil {
		t.Fatal(err)
	}
	if d := Diff(a, b); !d.IsEmpty() {
		t.Errorf("a temporary sequence diffs:\n%s", describeDiff(d))
	}
}

func TestDiffComparesRoutinesWithNonASCIITypes(t *testing.T) {
	c := loadCat(t, dbDDL+"CREATE PROCEDURE p(IN r ROW(`\u212a` INT)) SELECT 1;")
	if d := Diff(c, c); !d.IsEmpty() {
		t.Errorf("a catalog diffs against itself:\n%s", describeDiff(d))
	}
}

func TestDiffComparesKnownSQLModes(t *testing.T) {
	text := dbDDL + "CREATE TABLE t (a VARCHAR(10)); CREATE FUNCTION f() RETURNS VARCHAR(10) DETERMINISTIC RETURN 'a' || 'b'; CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 'x' || 'y'; CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 'a' || 'b';"
	load := func(mode string) *Catalog {
		c := loadCat(t, text)
		ctx := SessionContext{SQLMode: mode}
		c.ApplySessionContext(SessionContextMap{
			Functions: map[string]SessionContext{"f": ctx},
			Triggers:  map[string]SessionContext{"tr": ctx},
			Events:    map[string]SessionContext{"e": ctx},
		})
		return c
	}
	d := Diff(load("STRICT_TRANS_TABLES"), load("PIPES_AS_CONCAT,STRICT_TRANS_TABLES"))
	if len(d.Functions) != 1 || len(d.Triggers) != 1 || len(d.Events) != 1 {
		t.Errorf("want the function, trigger and event changed:\n%s", describeDiff(d))
	}
	if d := Diff(load("STRICT_TRANS_TABLES,PIPES_AS_CONCAT"), load("pipes_as_concat,STRICT_TRANS_TABLES")); !d.IsEmpty() {
		t.Errorf("one mode spelled two ways diffs:\n%s", describeDiff(d))
	}
	if d := Diff(loadCat(t, text), load("PIPES_AS_CONCAT")); !d.IsEmpty() {
		t.Errorf("an unknown mode diffs against a known one:\n%s", describeDiff(d))
	}
	zoned := func(zone string) *Catalog {
		c := loadCat(t, text)
		c.ApplySessionContext(SessionContextMap{Events: map[string]SessionContext{"e": {TimeZone: zone}}})
		return c
	}
	if d := Diff(zoned("+00:00"), zoned("+08:00")); len(d.Events) != 1 {
		t.Errorf("want the event changed with its time zone:\n%s", describeDiff(d))
	}
}

func TestDiffWithNormalizerSeesChangesBetweenComparisons(t *testing.T) {
	a := loadCat(t, dbDDL+"CREATE TABLE t (a INT CHECK (a > 0));")
	b := loadCat(t, dbDDL+"CREATE TABLE t (a INT CHECK (a > 0));")
	n := NormalizerFor(MariaDB118)
	if d := DiffWithNormalizer(a, b, n); !d.IsEmpty() {
		t.Fatalf("equal catalogs diff:\n%s", describeDiff(d))
	}
	if _, err := b.Exec("ALTER TABLE t MODIFY a INT;", nil); err != nil {
		t.Fatal(err)
	}
	if d := DiffWithNormalizer(a, b, n); d.IsEmpty() {
		t.Error("the dropped check does not diff with a reused normalizer")
	}
}

func TestCanonicalTableOptionsKeepsCommentSpacing(t *testing.T) {
	a := loadCat(t, dbDDL+"CREATE TABLE t (a INT) COMMENT 'a  b';")
	b := loadCat(t, dbDDL+"CREATE TABLE t (a INT) COMMENT 'a b';")
	if Diff(a, b).IsEmpty() {
		t.Error("comments differing in spacing must diff")
	}
}
