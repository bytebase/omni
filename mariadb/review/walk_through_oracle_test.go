package review

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// walkBase is the shop database the curated cases start from: tables with
// keys and foreign keys, one of a generated name, a table with a trigger,
// a view, a routine of each kind, an event, a sequence, and a
// system-versioned table with a trigger and a foreign key.
const walkBase = `
CREATE TABLE r (id int PRIMARY KEY);
CREATE TABLE t (id int PRIMARY KEY, a int, b int, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
CREATE TABLE nokey (a int);
CREATE TABLE trg (a int);
CREATE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 1;
CREATE VIEW v AS SELECT id FROM r;
CREATE PROCEDURE p() SELECT 1;
CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;
CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;
CREATE SEQUENCE sq;
CREATE TABLE vt (id int PRIMARY KEY, a int, KEY vi (a), CONSTRAINT vfk FOREIGN KEY (a) REFERENCES r (id)) WITH SYSTEM VERSIONING;
CREATE TRIGGER vtr BEFORE INSERT ON vt FOR EACH ROW SET NEW.a = 1;
`

// walkCase is a change on walkBase, after setup when it has one.
type walkCase struct {
	name     string
	setup    string
	change   string
	complete bool
}

var walkCases = []walkCase{
	// Databases and USE.
	{name: "databases", change: "CREATE DATABASE shop;\nCREATE DATABASE q;\nCREATE DATABASE q;\nCREATE DATABASE IF NOT EXISTS q;\nDROP DATABASE IF EXISTS nosuch;\nUSE nosuch;\nCREATE TABLE q.x (a int);\nCREATE TABLE q.x (a int);\nDROP DATABASE q;\nCREATE TABLE q.y (a int);", complete: true},
	{name: "no database selected", change: "DROP DATABASE shop;\nALTER DATABASE CHARACTER SET utf8mb4;\nCREATE TABLE x (a int);\nCREATE TEMPORARY TABLE x (a int);\nCREATE SEQUENCE s;\nCREATE PROCEDURE p2() SELECT 1;\nDROP PROCEDURE IF EXISTS p;\nDROP FUNCTION shop.f;\nDROP TRIGGER IF EXISTS tr;\nCREATE EVENT e2 ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nALTER EVENT e COMMENT 'x';", complete: true},
	// A temporary table or sequence needs no database.
	{name: "objects in a dropped database", change: "CREATE DATABASE z;\nDROP DATABASE z;\nCREATE PROCEDURE z.p() SELECT 1;\nCREATE FUNCTION z.f() RETURNS int DETERMINISTIC RETURN 1;\nCREATE VIEW z.w AS SELECT 1 AS a;\nCREATE SEQUENCE z.s;\nCREATE TEMPORARY TABLE z.tmp (a int);\nCREATE TEMPORARY SEQUENCE z.ts;", complete: true},
	{name: "USE moves unqualified names", change: "CREATE DATABASE q;\nUSE q;\nCREATE TABLE t (a int);\nCREATE TABLE t (a int);\nUSE shop;\nCREATE TABLE t (a int);", complete: true},
	// Tables. IF NOT EXISTS of a table that exists reads no LIKE source.
	{name: "create a table that exists", change: "CREATE TABLE t (a int);\nCREATE TABLE v (a int);\nCREATE TABLE sq (a int);\nCREATE TABLE IF NOT EXISTS t (a int);\nCREATE TABLE IF NOT EXISTS v (a int);\nDROP TABLE nokey;\nCREATE TABLE IF NOT EXISTS v LIKE nokey;\nCREATE TABLE IF NOT EXISTS t LIKE nokey;\nCREATE TABLE x LIKE nokey;\nCREATE TABLE x LIKE v;", complete: true},
	{name: "a new table's text", change: "CREATE TABLE x (a int, a int);\nCREATE TABLE x (a int, b int, KEY k (a), KEY k (b));\nCREATE TABLE x (a int PRIMARY KEY, b int, PRIMARY KEY (b));\nCREATE TABLE x (a int, KEY `PRIMARY` (a));", complete: true},
	{name: "a new table holds what it declares", change: "CREATE TABLE x (a int, b int UNIQUE, KEY k (a), CONSTRAINT fx FOREIGN KEY (a) REFERENCES r (id));\nALTER TABLE x ADD COLUMN a int;\nALTER TABLE x ADD INDEX k (b);\nALTER TABLE x DROP COLUMN zz;\nALTER TABLE x DROP INDEX k, DROP INDEX k;\nALTER TABLE nokey ADD CONSTRAINT fx FOREIGN KEY (a) REFERENCES r (id);\nCREATE TABLE y LIKE x;\nALTER TABLE y DROP INDEX k;", complete: true},
	{name: "names too long", change: "CREATE TABLE " + strings.Repeat("a", 65) + " (a int);\nALTER TABLE t ADD COLUMN " + strings.Repeat("c", 65) + " int;\nDROP TABLE " + strings.Repeat("d", 65) + ";\nCREATE SEQUENCE " + strings.Repeat("e", 65) + ";\nCREATE DATABASE " + strings.Repeat("f", 65) + ";", complete: true},
	// ALTER TABLE.
	{name: "ALTER TABLE drops first", change: "ALTER TABLE t ADD COLUMN a int, DROP COLUMN a;\nALTER TABLE t DROP INDEX i, ADD INDEX i (b);\nALTER TABLE t DROP COLUMN a;\nALTER TABLE t DROP COLUMN a;\nALTER TABLE t RENAME COLUMN b TO c, ADD COLUMN b int;\nALTER TABLE t DROP COLUMN zz, ADD COLUMN y int;", complete: true},
	{name: "ALTER TABLE subcommands", change: "ALTER TABLE v ADD COLUMN a int;\nALTER TABLE t ADD COLUMN A int;\nALTER TABLE t ADD COLUMN c int AFTER nosuch;\nALTER TABLE t DROP COLUMN nosuch;\nALTER TABLE t MODIFY COLUMN nosuch int;\nALTER TABLE t CHANGE nosuch c int;\nALTER TABLE t CHANGE a b int;\nALTER TABLE t ALTER COLUMN nosuch SET DEFAULT 1;\nALTER TABLE t ADD INDEX I (b);\nALTER TABLE t ADD INDEX `PRIMARY` (b);\nALTER TABLE t DROP INDEX nosuch;\nALTER TABLE t ADD PRIMARY KEY (a);\nALTER TABLE t RENAME INDEX nosuch TO j;\nALTER TABLE t RENAME INDEX i TO PRIMARY2;", complete: true},
	{name: "repeated subcommands", change: "ALTER TABLE t ADD COLUMN c int, ADD COLUMN c int;\nALTER TABLE t DROP COLUMN a, DROP COLUMN a;\nALTER TABLE t ADD INDEX k (a), ADD INDEX k (b);\nALTER TABLE nokey ADD PRIMARY KEY (a), ADD PRIMARY KEY (a);\nALTER TABLE t DROP INDEX i, DROP INDEX i;\nALTER TABLE t DROP FOREIGN KEY fk1, DROP FOREIGN KEY fk1;\nALTER TABLE t DROP PRIMARY KEY, DROP PRIMARY KEY;\nALTER TABLE t DROP PRIMARY KEY, DROP INDEX `PRIMARY`;", complete: true},
	// MariaDB has no invisible primary key.
	{name: "a table without a primary key has none", change: "ALTER TABLE nokey DROP PRIMARY KEY;\nCREATE TABLE g (a int);\nALTER TABLE g DROP PRIMARY KEY;\nALTER TABLE g ADD PRIMARY KEY (a);\nALTER TABLE g DROP PRIMARY KEY;\nALTER TABLE g DROP PRIMARY KEY;", complete: true},
	{name: "foreign keys", change: "CREATE TABLE y0 (a int, FOREIGN KEY (a) REFERENCES v (id));\nCREATE TABLE y00 (a int, CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id));\nALTER TABLE nokey DROP FOREIGN KEY fk1;\nALTER TABLE nokey ADD CONSTRAINT FK1 FOREIGN KEY (a) REFERENCES r (id);\nALTER TABLE t DROP FOREIGN KEY nosuch;\nALTER TABLE t DROP FOREIGN KEY fk1;\nALTER TABLE t DROP FOREIGN KEY fk1;\nCREATE TABLE gone (id int PRIMARY KEY);\nDROP TABLE gone;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES gone (id));\nSET foreign_key_checks = 0;\nCREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES gone (id));\nSET foreign_key_checks = 1;\nCREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES gone (id));\nCREATE TABLE w (id int PRIMARY KEY, p int, FOREIGN KEY (p) REFERENCES w (id));\nCREATE TABLE ys (a bigint, FOREIGN KEY (a) REFERENCES sq (next_not_cached_value));", complete: true},
	{name: "an unqualified parent is in the child's database", change: "CREATE DATABASE q;\nCREATE TABLE q.parent (id int PRIMARY KEY);\nCREATE TABLE q.child (id int, FOREIGN KEY (id) REFERENCES parent (id));\nALTER TABLE q.child ADD CONSTRAINT c2 FOREIGN KEY (id) REFERENCES parent (id);\nCREATE TABLE q.orphan (id int, FOREIGN KEY (id) REFERENCES nosuch (id));", complete: true},
	// An engine other than InnoDB parses foreign keys and ignores them.
	{name: "an engine that ignores foreign keys", setup: "CREATE TABLE my (id int PRIMARY KEY, a int) ENGINE=MyISAM;", change: "ALTER TABLE my DROP FOREIGN KEY nosuch;\nALTER TABLE my ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES nosuch (id);\nCREATE TABLE m2 (a int, CONSTRAINT mf FOREIGN KEY (a) REFERENCES nosuch (id), CONSTRAINT mf FOREIGN KEY (a) REFERENCES r (id)) ENGINE=MyISAM;\nCREATE TABLE ok (a int, CONSTRAINT mf FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE m3 LIKE my;\nALTER TABLE m3 ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id);\nCREATE TABLE inn (id int PRIMARY KEY, a int);\nALTER TABLE inn ENGINE=MyISAM;\nALTER TABLE inn ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES nosuch (id);\nALTER TABLE inn DROP FOREIGN KEY whatever;\nALTER TABLE inn ENGINE=Aria;\nCREATE TABLE ar (a int, CONSTRAINT af FOREIGN KEY (a) REFERENCES nosuch (id)) ENGINE=Aria;\nALTER TABLE inn ENGINE=InnoDB;\nALTER TABLE inn DROP FOREIGN KEY whatever;"},
	{name: "column names a statement adds", setup: "CREATE TABLE cols (a int, b int, c int);", change: "ALTER TABLE cols RENAME COLUMN a TO z, RENAME COLUMN b TO z;\nALTER TABLE cols CHANGE a z int, CHANGE b z int;\nALTER TABLE cols ADD COLUMN z int, RENAME COLUMN a TO z;\nALTER TABLE cols CHANGE a z int, ADD COLUMN Z int;\nALTER TABLE cols RENAME COLUMN a TO b, RENAME COLUMN b TO a;\nALTER TABLE cols ADD COLUMN b2 int AFTER b;\nALTER TABLE cols RENAME COLUMN a TO b, RENAME COLUMN b TO y;\nALTER TABLE cols ADD COLUMN y2 int AFTER y, ADD COLUMN b3 int AFTER b;\nALTER TABLE cols ADD COLUMN a int;\nALTER TABLE cols ADD KEY i (a), ADD KEY j (b);\nALTER TABLE cols RENAME INDEX i TO j, RENAME INDEX j TO i;\nALTER TABLE cols DROP INDEX i, DROP INDEX j;\nALTER TABLE cols ADD KEY i (a), ADD KEY j (b);\nALTER TABLE cols RENAME INDEX i TO j, DROP INDEX j;\nALTER TABLE cols DROP INDEX j;\nALTER TABLE cols DROP INDEX j;\nALTER TABLE cols CHANGE a A int, ADD COLUMN d int;\nALTER TABLE cols MODIFY a int PRIMARY KEY, MODIFY b int PRIMARY KEY;\nALTER TABLE cols CHANGE a a int PRIMARY KEY, ADD PRIMARY KEY (b);\nALTER TABLE cols MODIFY a int NOT NULL PRIMARY KEY, ADD COLUMN e int PRIMARY KEY;", complete: true},
	{name: "a temporary table is no parent", change: "DROP TABLE nokey;\nCREATE TEMPORARY TABLE nokey (id int PRIMARY KEY);\nCREATE TABLE y1 (a int, FOREIGN KEY (a) REFERENCES nokey (id));\nCREATE TEMPORARY TABLE r (id int PRIMARY KEY);\nCREATE TABLE y2 (a int, FOREIGN KEY (a) REFERENCES r (id));", complete: true},
	{name: "foreign keys named twice", change: "CREATE TABLE x (a int, b int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT f FOREIGN KEY (b) REFERENCES r (id));\nCREATE TABLE x (a int, b int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT F FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE x (a int, b int);\nALTER TABLE x ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), ADD CONSTRAINT F FOREIGN KEY (b) REFERENCES r (id);", complete: true},
	// RENAME TABLE renames the foreign keys named <table>_ibfk_...
	{name: "renamed generated foreign key names", change: "CREATE TABLE g (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));\nRENAME TABLE g TO h;\nALTER TABLE nokey ADD CONSTRAINT g_ibfk_1 FOREIGN KEY (a) REFERENCES r (id);\nALTER TABLE h DROP FOREIGN KEY h_ibfk_1;\nALTER TABLE h DROP FOREIGN KEY h_ibfk_1;\nRENAME TABLE child TO kid;\nALTER TABLE kid DROP FOREIGN KEY kid_ibfk_1;\nALTER TABLE kid DROP FOREIGN KEY kid_ibfk_1;\nALTER TABLE kid RENAME TO c2;\nCREATE TABLE child (id int PRIMARY KEY, b int, CONSTRAINT child_ibfk_1 FOREIGN KEY (b) REFERENCES r (id));", complete: true},
	{name: "a view stays in its database", change: "CREATE DATABASE q;\nRENAME TABLE v TO q.v;\nDROP VIEW v;", complete: true},
	{name: "RENAME", change: "RENAME TABLE t TO r;\nRENAME TABLE t TO tmp, r TO t, tmp TO r;\nALTER TABLE r ADD COLUMN a int;\nALTER TABLE t RENAME TO r;\nALTER TABLE t RENAME TO t2;\nALTER TABLE t ADD COLUMN zz int;\nCREATE DATABASE q;\nRENAME TABLE nokey TO q.nokey;\nCREATE TABLE q.nokey (a int);\nRENAME TABLE sq TO sq2;\nDROP SEQUENCE sq2;\nDROP TABLE sq2;", complete: true},
	// DROP and TRUNCATE.
	{name: "DROP TABLE IF EXISTS leaves a view", setup: "CREATE DATABASE q;\nCREATE VIEW q.v AS SELECT 1 AS a;", change: "DROP TABLE IF EXISTS q.v;\nDROP VIEW q.v;\nDROP VIEW q.v;"},
	{name: "DROP TABLE", change: "DROP TABLE v;\nDROP TABLE IF EXISTS nosuch, nokey;\nDROP TABLE nokey;\nDROP VIEW t;\nDROP VIEW IF EXISTS nosuch, v;\nDROP VIEW v;\nTRUNCATE TABLE nokey;\nTRUNCATE TABLE v;\nDROP TABLE sq;\nDROP SEQUENCE sq;", complete: true},
	{name: "a temporary table's definition", change: "CREATE TEMPORARY TABLE tmp (a int, a int);\nCREATE TEMPORARY TABLE tmp (a int, KEY k (a), KEY k (a));\nCREATE TEMPORARY TABLE tmp (a int PRIMARY KEY, PRIMARY KEY (a));\nCREATE TEMPORARY TABLE tmp (a int);\nCREATE TEMPORARY TABLE IF NOT EXISTS tmp (a int, a int);\nCREATE TEMPORARY TABLE tmp (a int);\nCREATE OR REPLACE TEMPORARY TABLE tmp (b int);", complete: true},
	{name: "temporary tables outlive their database", change: "CREATE DATABASE q;\nCREATE TEMPORARY TABLE q.tmp (a int);\nDROP DATABASE q;\nCREATE DATABASE q;\nDROP TEMPORARY TABLE q.tmp;\nDROP TEMPORARY TABLE q.tmp;", complete: true},
	{name: "temporary tables", change: "CREATE TEMPORARY TABLE t (z int);\nALTER TABLE nokey ADD COLUMN z int;\nINSERT INTO t (z) VALUES (1);\nDROP TABLE t;\nINSERT INTO t (id) VALUES (1);\nDROP TEMPORARY TABLE t;\nDROP TEMPORARY TABLE IF EXISTS t;\nCREATE TEMPORARY TABLE tt (a int);\nCREATE TEMPORARY TABLE tt (a int);\nTRUNCATE TABLE tt;\nDROP TABLE trg;\nCREATE TEMPORARY TABLE tmp2 (a int);\nCREATE TEMPORARY TABLE IF NOT EXISTS tmp2 LIKE trg;", complete: true},
	// Indexes, views, triggers, routines, events.
	{name: "CREATE and DROP INDEX", change: "DROP TABLE nokey;\nCREATE INDEX k ON nokey (a);\nCREATE INDEX k ON v (id);\nCREATE INDEX i ON t (b);\nCREATE INDEX `PRIMARY` ON t (b);\nDROP INDEX nosuch ON t;\nDROP INDEX i ON nokey;\nDROP INDEX i ON t;\nDROP INDEX i ON t;\nCREATE INDEX IF NOT EXISTS j ON t (b);", complete: true},
	{name: "views", change: "CREATE VIEW t AS SELECT 1 AS a;\nCREATE VIEW v AS SELECT 1 AS a;\nCREATE VIEW sq AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW t AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW sq AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW v AS SELECT 1 AS a;\nDROP TABLE nokey;\nALTER VIEW nokey AS SELECT 1 AS a;\nALTER VIEW t AS SELECT 1 AS a;\nDROP VIEW IF EXISTS t;\nDROP VIEW t;\nDROP VIEW sq;", complete: true},
	{name: "triggers", change: "CREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;\nDROP TRIGGER tr;\nDROP TRIGGER tr;\nCREATE TRIGGER TR BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;\nDROP TABLE child;\nCREATE TRIGGER tr2 BEFORE INSERT ON child FOR EACH ROW SET NEW.b = 1;", complete: true},
	{name: "routines and events", change: "CREATE PROCEDURE p() SELECT 1;\nCREATE PROCEDURE P() SELECT 1;\nCREATE FUNCTION p() RETURNS int DETERMINISTIC RETURN 1;\nCREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;\nDROP PROCEDURE nosuch;\nDROP FUNCTION shop.nosuch;\nALTER PROCEDURE nosuch COMMENT 'x';\nCREATE EVENT E ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nDROP EVENT nosuch;\nALTER EVENT nosuch COMMENT 'x';\nALTER EVENT e RENAME TO e2;\nDROP EVENT e;\nCREATE EVENT `shop.e2` ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nDROP EVENT `shop.e2`;\nDROP EVENT `shop.e2`;", complete: true},
	// DML.
	{name: "DML targets", change: "DROP TABLE nokey;\nINSERT INTO nokey VALUES (1);\nUPDATE nokey SET a = 1;\nDELETE FROM nokey;\nREPLACE INTO nokey VALUES (1);\nINSERT INTO v VALUES (1);\nINSERT INTO sq VALUES (1, 1, 1000, 1, 1, 1000, 0, 0);", complete: true},
	// Name case: tables and databases are exact under lower_case_table_names
	// 0; the other spaces fold, triggers aside.
	{name: "name case", change: "CREATE TABLE T (a int);\nCREATE TABLE T (a int);\nCREATE DATABASE q;\nCREATE DATABASE Q;\nCREATE DATABASE q;\nCREATE TABLE Q.x (a int);\nCREATE TABLE q.X (a int);\nCREATE TABLE q.X (a int);\nALTER TABLE r ADD COLUMN ID int;\nCREATE VIEW V AS SELECT 1 AS a;\nDROP PROCEDURE P;\nDROP PROCEDURE p;\nDROP EVENT E;\nDROP EVENT e;\nCREATE TRIGGER Tr BEFORE INSERT ON r FOR EACH ROW SET NEW.id = 1;", complete: true},
	// MariaDB statements: OR REPLACE, IF EXISTS drops, sequences.
	{name: "OR REPLACE", change: "CREATE OR REPLACE TABLE t (id int PRIMARY KEY, b int, CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));\nCREATE OR REPLACE TABLE sq (a int);\nALTER TABLE sq ADD COLUMN a int;\nCREATE OR REPLACE TABLE v (a int);\nCREATE OR REPLACE INDEX i ON t (b);\nCREATE OR REPLACE INDEX i ON t (b);\nCREATE OR REPLACE TRIGGER tr BEFORE INSERT ON trg FOR EACH ROW SET NEW.a = 2;\nCREATE OR REPLACE PROCEDURE p() SELECT 2;\nCREATE OR REPLACE FUNCTION f() RETURNS int DETERMINISTIC RETURN 2;\nCREATE OR REPLACE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 2;\nCREATE OR REPLACE TABLE t (a int);\nALTER TABLE t DROP COLUMN b;\nCREATE OR REPLACE TABLE vt (a int);\nALTER TABLE vt DROP COLUMN row_start;", complete: true},
	{name: "OR REPLACE DATABASE", change: "CREATE DATABASE q;\nCREATE TABLE q.x (a int);\nCREATE TEMPORARY TABLE q.tmp (a int);\nCREATE OR REPLACE DATABASE q;\nCREATE TABLE q.x (a int);\nCREATE TABLE q.x (a int);\nDROP TEMPORARY TABLE q.tmp;\nDROP TEMPORARY TABLE q.tmp;\nUSE q;\nCREATE OR REPLACE DATABASE q;\nCREATE TABLE y (a int);\nCREATE TABLE q.y (a int);\nUSE shop;\nCREATE TABLE q.y (a int);", complete: true},
	// OR REPLACE refuses to drop a parent a foreign key references, which
	// DROP DATABASE drops (verified), and the walk does not claim.
	{name: "OR REPLACE DATABASE of a referenced parent", change: "CREATE OR REPLACE DATABASE shop;\nCREATE TABLE t (a int);"},
	{name: "IF EXISTS drops", change: "ALTER TABLE t DROP COLUMN IF EXISTS zz, DROP INDEX IF EXISTS zz, DROP FOREIGN KEY IF EXISTS zz, DROP CONSTRAINT IF EXISTS zz;\nALTER TABLE t DROP COLUMN IF EXISTS a, DROP COLUMN IF EXISTS a;\nALTER TABLE t DROP COLUMN a;\nALTER TABLE t DROP INDEX IF EXISTS i;\nDROP INDEX i ON t;\nALTER TABLE child DROP FOREIGN KEY IF EXISTS fk1;\nALTER TABLE child ADD CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id);\nALTER TABLE t DROP FOREIGN KEY IF EXISTS fk1;\nALTER TABLE child ADD CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id);\nALTER TABLE nokey DROP INDEX IF EXISTS `PRIMARY`;\nALTER TABLE t DROP INDEX IF EXISTS `PRIMARY`;\nALTER TABLE t DROP PRIMARY KEY;", complete: true},
	{name: "sequences", change: "CREATE SEQUENCE s1;\nCREATE SEQUENCE s1;\nCREATE SEQUENCE t;\nCREATE SEQUENCE v;\nCREATE SEQUENCE IF NOT EXISTS t;\nCREATE OR REPLACE SEQUENCE v;\nDROP SEQUENCE t;\nDROP SEQUENCE v;\nDROP VIEW sq;\nDROP SEQUENCE IF EXISTS t, v, nosuch;\nDROP SEQUENCE s1;\nDROP SEQUENCE s1;\nCREATE TABLE x LIKE sq;\nDROP SEQUENCE x;\nCREATE OR REPLACE SEQUENCE nokey;\nDROP SEQUENCE nokey;\nCREATE DATABASE z;\nDROP DATABASE z;\nCREATE SEQUENCE z.s;", complete: true},
	{name: "temporary sequences", change: "CREATE TEMPORARY SEQUENCE sq;\nDROP SEQUENCE sq;\nDROP SEQUENCE sq;\nDROP SEQUENCE sq;\nCREATE TEMPORARY TABLE tt (a int);\nCREATE TEMPORARY SEQUENCE tt;\nDROP TEMPORARY SEQUENCE tt;\nCREATE TEMPORARY SEQUENCE ts;\nCREATE TEMPORARY TABLE ts (a int);\nDROP TEMPORARY SEQUENCE ts;\nDROP TEMPORARY SEQUENCE ts;\nCREATE TEMPORARY SEQUENCE IF NOT EXISTS tt;\nCREATE OR REPLACE TEMPORARY SEQUENCE tt;\nDROP TEMPORARY SEQUENCE tt;\nCREATE TEMPORARY TABLE nokey (a int);\nDROP SEQUENCE nokey;\nDROP TEMPORARY TABLE nokey;", complete: true},
	// A sequence or a system-versioned table in a statement the walk does
	// not follow on it ends the walk.
	{name: "ALTER TABLE of a sequence", change: "ALTER TABLE sq COMMENT 'x';\nALTER TABLE sq ADD COLUMN z int;"},
	{name: "TRUNCATE of a sequence", change: "TRUNCATE TABLE sq;\nCREATE TABLE t (a int);"},
	{name: "an index on a sequence", change: "CREATE INDEX k ON sq (cache_size);\nDROP INDEX k ON sq;\nCREATE TABLE t (a int);"},
	{name: "a trigger on a sequence", change: "CREATE TRIGGER trs BEFORE INSERT ON sq FOR EACH ROW SET @x = 1;\nCREATE TABLE t (a int);"},
	{name: "indexes of a system-versioned table", change: "CREATE INDEX vj ON vt (a);\nCREATE INDEX vi ON vt (a);\nDROP INDEX nosuch ON vt;\nDROP INDEX vj ON vt;\nDROP INDEX vj ON vt;"},
	{name: "a system-versioned table hides columns", setup: "SET GLOBAL system_versioning_alter_history = KEEP;", change: "ALTER TABLE vt ADD COLUMN x int AFTER row_start;\nALTER TABLE vt ADD COLUMN x int;"},
	{name: "a new system-versioned table hides columns", setup: "SET GLOBAL system_versioning_alter_history = KEEP;", change: "CREATE TABLE w (a int) WITH SYSTEM VERSIONING;\nALTER TABLE w ADD COLUMN x int AFTER row_end;\nALTER TABLE w DROP COLUMN x;"},
	{name: "a column's versioning makes the table system-versioned", setup: "SET GLOBAL system_versioning_alter_history = KEEP;", change: "CREATE TABLE w (a int WITH SYSTEM VERSIONING, b int);\nALTER TABLE w ADD COLUMN x int AFTER row_start;\nALTER TABLE w DROP COLUMN x;"},
	{name: "LIKE copies the versioning", setup: "SET GLOBAL system_versioning_alter_history = KEEP;", change: "CREATE TABLE w LIKE vt;\nALTER TABLE w ADD COLUMN x int AFTER row_start;\nALTER TABLE w DROP COLUMN x;"},
	{name: "RENAME keeps the versioning", setup: "SET GLOBAL system_versioning_alter_history = KEEP;", change: "RENAME TABLE vt TO w;\nALTER TABLE w ADD COLUMN x int AFTER row_start;\nALTER TABLE w DROP COLUMN x;"},
	// The sync before Bytebase #21665 lost the triggers of the
	// system-versioned tables it skipped.
	{name: "triggers the sync may not list", change: "DROP TRIGGER vtr;\nDROP TRIGGER vtr;"},
	// MySQL's ALTER INDEX INVISIBLE, which omni's parser takes and MariaDB
	// rejects, ends the walk.
	{name: "MySQL forms", change: "ALTER TABLE t ALTER INDEX nosuch INVISIBLE;\nALTER TABLE t ALTER COLUMN a SET NOT NULL;"},
	// MariaDB drops what exists of a multi-object drop, and a failed CREATE
	// OR REPLACE TABLE may have dropped the old table: Judge stops at such
	// a statement (walktest.PartialRejecter).
	{name: "DROP TABLE drops what exists", change: "DROP TABLE nokey;\nDROP TABLE trg, nokey;\nCREATE TABLE trg (a int);", complete: true},
	{name: "DROP TABLE of a view drops the tables", change: "DROP TABLE nokey, v;\nCREATE TABLE nokey (a int);", complete: true},
	{name: "DROP VIEW drops what exists", change: "DROP VIEW v;\nCREATE VIEW w AS SELECT 1 AS a;\nDROP VIEW w, v;\nCREATE VIEW w AS SELECT 1 AS a;", complete: true},
	{name: "DROP SEQUENCE drops what exists", change: "DROP SEQUENCE sq, t;\nCREATE SEQUENCE sq;", complete: true},
	{name: "DROP TEMPORARY TABLE drops what exists", change: "CREATE TEMPORARY TABLE tt1 (a int);\nDROP TEMPORARY TABLE tt1, nosuch;\nDROP TEMPORARY TABLE tt1;", complete: true},
	{name: "a failed CREATE OR REPLACE TABLE drops the table", change: "CREATE OR REPLACE TABLE nokey (a int, a int);\nCREATE TABLE nokey (a int);", complete: true},
	{name: "a failed CREATE OR REPLACE TABLE of a foreign key name", change: "CREATE OR REPLACE TABLE nokey (a int, CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE nokey (a int);", complete: true},
	// What the walk leaves Unknown by design: relations and triggers the
	// sync does not list, databases other than the target, the tables a
	// query reads, generated index names, foreign key names held per
	// database, a key two parts of one ALTER TABLE both change, a trigger
	// of another spelling, and a function an unqualified DROP FUNCTION may
	// find loaded.
	{name: "unknown by design", change: "DROP FUNCTION nosuch;\nDROP TRIGGER nosuch;\nDROP TRIGGER TR;\nDROP TABLE IF EXISTS nokey;\nDROP TABLE t, nosuch;"},
	{name: "unknown by design, after the walk", change: "DROP DATABASE nosuch;\nRENAME TABLE t TO nodb.t;\nCREATE VIEW nodb.w AS SELECT 1 AS a;\nCREATE TABLE x AS SELECT * FROM nosuch;\nCREATE TABLE x (a int, b int UNIQUE);\nALTER TABLE x DROP INDEX zz;\nCREATE TABLE y (a int, CONSTRAINT fy FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE z LIKE y;\nALTER TABLE z DROP FOREIGN KEY fy;\nALTER TABLE t ADD COLUMN z int, DROP COLUMN z;\nALTER TABLE t ADD INDEX j (b), DROP INDEX j;\nALTER TABLE t ADD COLUMN x int, RENAME COLUMN x TO y;\nALTER TABLE nosuch ADD COLUMN a int;\nCREATE TABLE IF NOT EXISTS nosuch2 (a int);\nALTER TABLE nosuch2 DROP COLUMN b;\nRENAME TABLE nosuch TO x2;"},
	{name: "a trigger of another database may create a temporary table", setup: "DROP TRIGGER tr;\nCREATE DATABASE q;\nCREATE TABLE q.o (a int);\nCREATE TRIGGER q.ot BEFORE INSERT ON q.o FOR EACH ROW CREATE TEMPORARY TABLE IF NOT EXISTS shop.made (a int);", change: "INSERT INTO q.o VALUES (1);\nDROP TEMPORARY TABLE made;\nDROP TEMPORARY TABLE made;"},
	{name: "a query's function may create a temporary table", setup: "CREATE FUNCTION mk() RETURNS int DETERMINISTIC BEGIN CREATE TEMPORARY TABLE IF NOT EXISTS made (a int); RETURN 1; END;", change: "CREATE TABLE c AS SELECT mk() AS x;\nDROP TEMPORARY TABLE made;\nDROP TEMPORARY TABLE made;"},
	{name: "a trigger may create a temporary table", change: "CREATE TRIGGER tt BEFORE INSERT ON nokey FOR EACH ROW CREATE TEMPORARY TABLE IF NOT EXISTS shadow (a int);\nINSERT INTO nokey VALUES (1);\nINSERT INTO shadow VALUES (1);\nDROP TABLE shadow;"},
}

// walkSync is how Bytebase's sync lists a MariaDB database: since #21665
// with its sequences, and its system-versioned tables as tables marked in
// their create options; before it with neither. The walk must hold on
// both.
type walkSync uint8

const (
	syncCurrent walkSync = iota
	syncBefore21665
)

func (s walkSync) String() string {
	if s == syncBefore21665 {
		return "before #21665"
	}
	return "current"
}

var walkServer struct {
	once sync.Once
	db   *sql.DB
	err  error
}

// startWalkMariaDB returns a pool on the oracle's MariaDB that selects no
// database, so a change that drops shop leaves it usable.
func startWalkMariaDB(t *testing.T) *sql.DB {
	t.Helper()
	startMariaDB(t)
	walkServer.once.Do(func() {
		ctx := context.Background()
		endpoint, err := mariadbOracle.container.PortEndpoint(ctx, "3306/tcp", "")
		if err != nil {
			walkServer.err = err
			return
		}
		walkServer.db, walkServer.err = sql.Open("mysql", "root:test@tcp("+endpoint+")/?multiStatements=true")
		if walkServer.err == nil {
			walkServer.err = walkServer.db.PingContext(ctx)
		}
	})
	if walkServer.err != nil {
		t.Fatalf("mariadb unavailable: %v", walkServer.err)
	}
	return walkServer.db
}

// session opens a session of its own, which the cleanup closes rather
// than returns to the pool: a pooled session keeps the temporary tables
// of the change before it.
func session(t testing.TB, db *sql.DB) *sql.Conn {
	t.Helper()
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = conn.Raw(func(any) error { return driver.ErrBadConn })
		_ = conn.Close()
	})
	return conn
}

func TestWalkThroughAgainstMariaDB(t *testing.T) {
	db := startWalkMariaDB(t)
	for _, s := range []walkSync{syncCurrent, syncBefore21665} {
		t.Run(s.String(), func(t *testing.T) {
			var found, rejected int
			for _, tt := range walkCases {
				t.Run(tt.name, func(t *testing.T) {
					resetWalk(t, db, walkBase+tt.setup)
					// The sync before #21665 hides the base's sequence and
					// system-versioned table, which costs recall alone.
					f, r := judgeWalk(t, db, s, tt.change, tt.complete && s == syncCurrent)
					found += f
					rejected += r
				})
			}
			t.Logf("found %d of the %d statements the server rejects", found, rejected)
		})
	}
}

// resetWalk drops every user database and recreates shop from base, with
// the server's settings as they ship.
func resetWalk(t *testing.T, db *sql.DB, base string) {
	t.Helper()
	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')")
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			t.Fatal(err)
		}
		names = append(names, n)
	}
	rows.Close()
	conn := session(t, db)
	stmts := []string{"SET foreign_key_checks = 0", "SET GLOBAL system_versioning_alter_history = ERROR"}
	for _, n := range names {
		stmts = append(stmts, "DROP DATABASE `"+strings.ReplaceAll(n, "`", "``")+"`")
	}
	stmts = append(stmts, "SET foreign_key_checks = 1", "CREATE DATABASE shop", "USE shop")
	for _, q := range stmts {
		if _, err := conn.ExecContext(ctx, q); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	for _, q := range splitWalk(base) {
		if _, err := conn.ExecContext(ctx, q); err != nil {
			t.Fatalf("base %s: %v", q, err)
		}
	}
}

// splitWalk splits a change into its statements, each with its ';' so
// the statements joined again are a change.
func splitWalk(sql string) []string {
	var out []string
	for _, seg := range parser.Split(sql) {
		if seg.Empty() {
			continue
		}
		if text := strings.TrimSpace(seg.Text); text != "" {
			out = append(out, text+";")
		}
	}
	return out
}

// judgeWalk runs the change on a fresh connection with shop selected, as
// Bytebase runs it, and holds the walk to it (walktest.Judge). The target
// is the database as the admin connection syncs it before the change.
func judgeWalk(t *testing.T, db *sql.DB, s walkSync, change string, complete bool) (found, rejected int) {
	t.Helper()
	ctx := context.Background()
	conn := session(t, db)
	if _, err := conn.ExecContext(ctx, "USE shop"); err != nil {
		t.Fatal(err)
	}
	e := &mariadbEngine{
		conn: conn,
		db:   db,
		target: review.Target{
			Schema:           syncWalk(t, db, s),
			ForeignKeyChecks: review.SwitchOn,
		},
	}
	return walktest.Judge(t, e, change, complete)
}

// mariadbEngine is the oracle's MariaDB side of one change. MariaDB
// commits DDL as it runs, and drops what exists of a multi-object DROP
// before it fails on what does not, which Digest notices.
type mariadbEngine struct {
	conn   *sql.Conn
	db     *sql.DB
	target review.Target
}

func (*mariadbEngine) RejectsPartially() {}

func (e *mariadbEngine) Statements(change string) []string { return splitWalk(change) }

// errnoPattern reads the storage engine's errno from a 1005 message: 121
// for a foreign key name taken, 150 for a foreign key it cannot form.
var errnoPattern = regexp.MustCompile(`errno: (\d+)`)

func (e *mariadbEngine) Exec(t testing.TB, stmt string) walktest.Outcome {
	t.Helper()
	if _, err := e.conn.ExecContext(context.Background(), stmt); err != nil {
		var me *gomysql.MySQLError
		if !errors.As(err, &me) {
			t.Fatalf("%s: %v", stmt, err)
		}
		code := strconv.Itoa(int(me.Number))
		if m := errnoPattern.FindStringSubmatch(me.Message); me.Number == 1005 && m != nil {
			code += ":" + m[1]
		}
		return walktest.Outcome{Code: code, Message: me.Message}
	}
	return walktest.Outcome{}
}

func (e *mariadbEngine) Report(t testing.TB, change string) *walktest.Finding {
	t.Helper()
	if strings.TrimSpace(change) == "" {
		return nil
	}
	result, err := Review(context.Background(), change, review.Options{Rules: []review.Rule{review.WalkThrough}}, []review.Target{e.target})
	if err != nil {
		t.Fatalf("Review: %v", err)
	}
	var out *walktest.Finding
	for _, f := range result.Findings {
		switch f.Rule {
		case review.Syntax:
			t.Fatalf("syntax: %s\n%s", f.Message, change)
		case review.WalkThrough:
			if out != nil {
				t.Errorf("more than one finding: %q and %q", out.Message, f.Message)
			}
			out = &walktest.Finding{Statement: f.Statement, Message: f.Message}
		}
	}
	return out
}

func (e *mariadbEngine) Violation(t testing.TB, change string) *walk.Violation {
	t.Helper()
	stmts, syntax := parse(change, statementRanges(change))
	if syntax != nil {
		t.Fatalf("syntax: %s", syntax.Message)
	}
	v, err := walkThrough(context.Background(), change, stmts, e.target)
	if err != nil || v == nil {
		t.Fatalf("no violation for a change Review reports on (%v):\n%s", err, change)
	}
	return v
}

// walkCodes maps a finding's message to the MariaDB errors it stands for.
// A 1005 carries the storage engine's errno.
var walkCodes = []struct {
	match func(string) bool
	codes []string
}{
	{func(m string) bool {
		return strings.HasPrefix(m, "Temporary sequence ") || strings.HasPrefix(m, "Temporary table ") && strings.HasSuffix(m, " is not a sequence")
	}, []string{"4091"}},
	{func(m string) bool { return strings.HasPrefix(m, "Temporary table ") }, []string{"1051"}},
	{func(m string) bool { return strings.HasPrefix(m, "Foreign key target ") }, []string{"1005:150"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Foreign key ") && (strings.Contains(m, " specified twice ") || strings.Contains(m, " added twice "))
	}, []string{"1061"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "View ") && strings.HasSuffix(m, " cannot move to another database")
	}, []string{"1450"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Foreign key ") && strings.HasSuffix(m, " already exists in the database")
	}, []string{"1005:121"}},
	{func(m string) bool { return strings.HasPrefix(m, "Foreign key ") }, []string{"1091"}},
	{func(m string) bool { return strings.HasPrefix(m, "Sequence ") }, []string{"4091"}},
	{func(m string) bool { return strings.HasSuffix(m, " is not a sequence") }, []string{"4090"}},
	{func(m string) bool { return strings.HasSuffix(m, " is a view") }, []string{"1965"}},
	{func(m string) bool { return strings.HasPrefix(m, "Table ") && strings.HasSuffix(m, " already exists") }, []string{"1050"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Table ") && strings.HasSuffix(m, " already has a primary key")
	}, []string{"1068"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Table ") && strings.HasSuffix(m, " has no primary key")
	}, []string{"1091"}},
	{func(m string) bool { return strings.HasPrefix(m, "Table ") && strings.HasSuffix(m, " does not exist") }, []string{"1146", "1051"}},
	{func(m string) bool { return strings.HasPrefix(m, "View ") }, []string{"4092"}},
	{func(m string) bool { return strings.HasSuffix(m, " is not a base table") }, []string{"1347"}},
	{func(m string) bool { return strings.HasSuffix(m, " is not a view") }, []string{"1347", "4092"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Database ") && strings.HasSuffix(m, " already exists")
	}, []string{"1007"}},
	{func(m string) bool { return strings.HasPrefix(m, "Database ") }, []string{"1049", "1008", "1025"}},
	{func(m string) bool { return m == noDatabaseMessage }, []string{"1046"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Column ") && (strings.Contains(m, " already exists in ") || strings.Contains(m, " twice "))
	}, []string{"1060", "1091"}},
	{func(m string) bool { return strings.HasPrefix(m, "Column ") }, []string{"1054", "1091"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Index ") && (strings.Contains(m, " already exists in ") || strings.Contains(m, " specified twice") || strings.Contains(m, " added twice"))
	}, []string{"1061"}},
	{func(m string) bool { return strings.HasPrefix(m, "Index name PRIMARY") }, []string{"1280"}},
	{func(m string) bool { return strings.HasPrefix(m, "Index ") }, []string{"1091"}},
	{func(m string) bool { return strings.HasPrefix(m, "Multiple primary keys") }, []string{"1068"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Trigger ") && strings.HasSuffix(m, " already exists")
	}, []string{"1359"}},
	{func(m string) bool { return strings.HasPrefix(m, "Trigger ") }, []string{"1360"}},
	{func(m string) bool {
		return (strings.HasPrefix(m, "PROCEDURE ") || strings.HasPrefix(m, "FUNCTION ")) && strings.HasSuffix(m, " already exists")
	}, []string{"1304"}},
	{func(m string) bool { return strings.HasPrefix(m, "PROCEDURE ") || strings.HasPrefix(m, "FUNCTION ") }, []string{"1305"}},
	{func(m string) bool { return strings.HasPrefix(m, "Event ") && strings.HasSuffix(m, " already exists") }, []string{"1537"}},
	{func(m string) bool { return strings.HasPrefix(m, "Event ") }, []string{"1539"}},
	{func(m string) bool { return strings.HasPrefix(m, "Identifier name ") }, []string{"1059", "1102", "1103"}},
}

func (e *mariadbEngine) Codes(message string) []string {
	for _, c := range walkCodes {
		if c.match(message) {
			return c.codes
		}
	}
	return nil
}

// Holds checks each fact the violation claims against information_schema.
// A claim about the text holds by itself; a temporary table is read with
// SHOW CREATE TABLE in the change's session.
func (e *mariadbEngine) Holds(t testing.TB, v *walk.Violation, message string) (bool, string) {
	t.Helper()
	for _, c := range walktest.Claims(v) {
		if c.Fact == walk.Unknown {
			continue
		}
		present, kind := e.exists(t, c.Key)
		switch {
		case present != (c.Fact == walk.Present):
			return false, fmt.Sprintf("%s %q in %q is present: %v", c.Key.Space, c.Key.Name, walk.Names(c.Key.Parent), present)
		case present && c.Kind != "" && c.Key.Space == spaceRelation && kind != c.Kind:
			return false, fmt.Sprintf("relation %s is %s, not %s", c.Key.Name, kind, c.Kind)
		}
	}
	return true, ""
}

// exists reads whether the key's object exists, and a relation's kind.
// Names compare as lower_case_table_names 0 has them: tables and
// databases exactly, the other spaces folded.
func (e *mariadbEngine) exists(t testing.TB, k walk.Key) (bool, string) {
	t.Helper()
	ctx := context.Background()
	names := walk.Names(k.Parent)
	eq := func(column string, folds bool) string {
		if folds {
			return "LOWER(" + column + ") = BINARY ?"
		}
		return "BINARY " + column + " = BINARY ?"
	}
	name := strings.TrimPrefix(k.Name, nonASCII)
	nameFolds := !strings.HasPrefix(k.Name, nonASCII)
	unmark := func(s string) string { return strings.TrimPrefix(s, nonASCII) }
	count := func(q string, args ...any) bool {
		t.Helper()
		var n int
		if err := e.db.QueryRowContext(ctx, q, args...).Scan(&n); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		return n > 0
	}
	switch k.Space {
	case spaceDatabase:
		return count("SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE "+eq("SCHEMA_NAME", false), name), ""
	case spaceTemporary:
		var tbl, ddl string
		err := e.conn.QueryRowContext(ctx, "SHOW CREATE TABLE `"+unmark(names[0])+"`.`"+name+"`").Scan(&tbl, &ddl)
		return err == nil && strings.HasPrefix(ddl, "CREATE TEMPORARY TABLE"), ""
	case spaceRelation:
		var kind string
		err := e.db.QueryRowContext(ctx, "SELECT TABLE_TYPE FROM information_schema.TABLES WHERE "+eq("TABLE_SCHEMA", false)+" AND "+eq("TABLE_NAME", false), unmark(names[0]), name).Scan(&kind)
		if errors.Is(err, sql.ErrNoRows) {
			return false, ""
		}
		if err != nil {
			t.Fatal(err)
		}
		switch kind {
		case "VIEW":
			return true, relView
		case "SEQUENCE":
			return true, relSequence
		}
		return true, relTable
	case spaceColumn, spaceIndex, spaceTableForeignKey:
		table, column := "information_schema.COLUMNS", "COLUMN_NAME"
		switch k.Space {
		case spaceIndex:
			table, column = "information_schema.STATISTICS", "INDEX_NAME"
		case spaceTableForeignKey:
			return count("SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE "+eq("CONSTRAINT_SCHEMA", false)+" AND "+eq("TABLE_NAME", false)+" AND "+eq("CONSTRAINT_NAME", nameFolds), unmark(names[0]), unmark(names[1]), name), ""
		}
		return count("SELECT COUNT(*) FROM "+table+" WHERE "+eq("TABLE_SCHEMA", false)+" AND "+eq("TABLE_NAME", false)+" AND "+eq(column, nameFolds), unmark(names[0]), unmark(names[1]), name), ""
	}
	where := map[string]string{
		spaceForeignKey: "information_schema.REFERENTIAL_CONSTRAINTS WHERE " + eq("CONSTRAINT_SCHEMA", false) + " AND " + eq("CONSTRAINT_NAME", nameFolds),
		spaceTrigger:    "information_schema.TRIGGERS WHERE " + eq("TRIGGER_SCHEMA", false) + " AND " + eq("TRIGGER_NAME", nameFolds),
		spaceProcedure:  "information_schema.ROUTINES WHERE ROUTINE_TYPE = 'PROCEDURE' AND " + eq("ROUTINE_SCHEMA", false) + " AND " + eq("ROUTINE_NAME", nameFolds),
		spaceFunction:   "information_schema.ROUTINES WHERE ROUTINE_TYPE = 'FUNCTION' AND " + eq("ROUTINE_SCHEMA", false) + " AND " + eq("ROUTINE_NAME", nameFolds),
		spaceEvent:      "information_schema.EVENTS WHERE " + eq("EVENT_SCHEMA", false) + " AND " + eq("EVENT_NAME", nameFolds),
	}[k.Space]
	if where == "" {
		t.Fatalf("unknown space %s", k.Space)
	}
	return count("SELECT COUNT(*) FROM "+where, unmark(names[0]), name), ""
}

// Digest reads every user database's objects, and the change session's
// temporary tables, so a rejected statement that changed one is noticed.
// Each query names the databases, which keeps information_schema from
// reading the system ones.
func (e *mariadbEngine) Digest(t testing.TB) string {
	t.Helper()
	var b strings.Builder
	read := func(q querier, query string, args ...any) []string {
		rows, err := q.QueryContext(context.Background(), query, args...)
		if err != nil {
			t.Fatalf("%s: %v", query, err)
		}
		defer rows.Close()
		cols, _ := rows.Columns()
		vals := make([]sql.NullString, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		var first []string
		for rows.Next() {
			if err := rows.Scan(ptrs...); err != nil {
				t.Fatal(err)
			}
			first = append(first, vals[0].String)
			for _, v := range vals {
				b.WriteString(v.String)
				b.WriteByte(0)
			}
			b.WriteByte('\n')
		}
		return first
	}
	read(e.conn, "SELECT TABLE_SCHEMA, TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_TYPE LIKE 'TEMPORARY%' ORDER BY 1, 2")
	dbs := read(e.db, "SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys') ORDER BY 1")
	if len(dbs) == 0 {
		return b.String()
	}
	in := "(?" + strings.Repeat(", ?", len(dbs)-1) + ")"
	args := make([]any, len(dbs))
	for i, d := range dbs {
		args[i] = d
	}
	for _, q := range []string{
		"SELECT TABLE_SCHEMA, TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA IN " + in + " ORDER BY 1, 2",
		"SELECT TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA IN " + in + " ORDER BY 1, 2, 3",
		"SELECT TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX FROM information_schema.STATISTICS WHERE TABLE_SCHEMA IN " + in + " ORDER BY 1, 2, 3, 4",
		"SELECT CONSTRAINT_SCHEMA, CONSTRAINT_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN " + in + " ORDER BY 1, 2",
		"SELECT TRIGGER_SCHEMA, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN " + in + " ORDER BY 1, 2",
		"SELECT ROUTINE_SCHEMA, ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA IN " + in + " ORDER BY 1, 2, 3",
		"SELECT EVENT_SCHEMA, EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN " + in + " ORDER BY 1, 2",
	} {
		read(e.db, q, args...)
	}
	return b.String()
}

// querier is a pool or one of its sessions.
type querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// syncWalk reads shop into the metadata Bytebase's MariaDB sync produces
// for what WalkThrough reads: base tables with their engines, columns,
// indexes, foreign keys, and triggers, views, routines, and events; and,
// since #21665, sequences and system-versioned tables, which the sync
// before it skipped with their foreign keys and triggers.
//
// bytebase: backend/plugin/db/mysql/sync.go
func syncWalk(t *testing.T, db *sql.DB, s walkSync) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	out := syncSchema(t, conn)
	schema := out.Schemas[0]
	query := func(q string, add func(a, b string)) {
		rows, err := conn.QueryContext(context.Background(), q)
		if err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		defer rows.Close()
		for rows.Next() {
			var a, b string
			if err := rows.Scan(&a, &b); err != nil {
				t.Fatal(err)
			}
			add(a, b)
		}
	}
	query("SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop'", func(kind, name string) {
		if kind == "PROCEDURE" {
			schema.Procedures = append(schema.Procedures, &metadata.ProcedureMetadata{Name: name})
		} else {
			schema.Functions = append(schema.Functions, &metadata.FunctionMetadata{Name: name})
		}
	})
	query("SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop'", func(_, name string) {
		schema.Events = append(schema.Events, &metadata.EventMetadata{Name: name})
	})
	types := make(map[string]string)
	engines := make(map[string]string)
	query("SELECT TABLE_NAME, CONCAT(TABLE_TYPE, '/', IFNULL(ENGINE, '')) FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop'", func(name, kind string) {
		types[name], engines[name], _ = strings.Cut(kind, "/")
	})
	var tables []*metadata.TableMetadata
	for _, tbl := range schema.Tables {
		tbl.Engine = engines[tbl.Name]
		switch types[tbl.Name] {
		case "BASE TABLE":
			tables = append(tables, tbl)
		case "SYSTEM VERSIONED":
			if s == syncCurrent {
				tbl.CreateOptions = "WITH SYSTEM VERSIONING"
				tables = append(tables, tbl)
			}
		case "SEQUENCE":
			if s == syncCurrent {
				schema.Sequences = append(schema.Sequences, &metadata.SequenceMetadata{Name: tbl.Name})
			}
		}
	}
	schema.Tables = tables
	return out
}

// TestWalkThroughSweepMariaDB runs random changes over a small pool of
// names that collide often, case included, on each sync by turns.
// Precision must hold on every one; recall is logged. WALK_SWEEP_N sets
// the number of changes, WALK_SWEEP_SEED the seed, and WALK_SWEEP_IMAGE
// is not read: the sweep runs on the oracle's image.
func TestWalkThroughSweepMariaDB(t *testing.T) {
	n, seed, _ := walktest.SweepParams(t, 40, 20261026, "mariadb:11.8.8")
	db := startWalkMariaDB(t)
	rng := rand.New(rand.NewSource(seed))
	var found, rejected int
	for i := range n {
		base := walkSweepChange(rng, 8+rng.Intn(8))
		change := walkSweepChange(rng, 15+rng.Intn(15))
		s := walkSync(i % 2)
		t.Run(fmt.Sprintf("change%d", i), func(t *testing.T) {
			resetWalk(t, db, "")
			conn := session(t, db)
			if _, err := conn.ExecContext(context.Background(), "USE shop"); err != nil {
				t.Fatal(err)
			}
			// The base's temporary tables end with its session, before the sync.
			for _, q := range splitWalk(base) {
				_, _ = conn.ExecContext(context.Background(), q)
			}
			_ = conn.Raw(func(any) error { return driver.ErrBadConn })
			f, r := judgeWalk(t, db, s, change, false)
			found += f
			rejected += r
			if t.Failed() {
				t.Logf("sync %s, base:\n%s\nchange:\n%s", s, base, change)
			}
		})
	}
	t.Logf("seed %d: found %d of the %d statements the server rejects", seed, found, rejected)
}

func walkSweepChange(rng *rand.Rand, n int) string {
	pick := func(xs ...string) string { return xs[rng.Intn(len(xs))] }
	maybe := func(s string) string {
		if rng.Intn(3) == 0 {
			return s
		}
		return ""
	}
	db := func() string { return pick("", "", "", "q.") }
	tbl := func() string { return db() + pick("a", "b", "c", "A", "v", "s", "a_ibfk_1") }
	col := func() string { return pick("x", "y", "z", "X", "row_start") }
	idx := func() string { return pick("i", "j", "x", "I") }
	fk := func() string { return pick("f", "F", "a_ibfk_1", "g") }
	ifExists := func() string { return maybe(" IF EXISTS") }
	var b strings.Builder
	for range n {
		var stmt string
		switch rng.Intn(34) {
		case 0, 1:
			stmt = fmt.Sprintf("CREATE %s%sTABLE %s%s (id int PRIMARY KEY, x int%s, y int)%s", pick("", "", "OR REPLACE "), maybe("TEMPORARY "), pick("", "IF NOT EXISTS "), tbl(), maybe(" UNIQUE"), maybe(" WITH SYSTEM VERSIONING"))
			stmt = strings.Replace(stmt, "OR REPLACE TEMPORARY TABLE IF NOT EXISTS", "OR REPLACE TEMPORARY TABLE", 1)
			stmt = strings.Replace(stmt, "OR REPLACE TABLE IF NOT EXISTS", "OR REPLACE TABLE", 1)
		case 2:
			stmt = fmt.Sprintf("CREATE TABLE %s (x int, y int, KEY %s (x)%s)", tbl(), idx(), maybe(", CONSTRAINT "+fk()+" FOREIGN KEY (y) REFERENCES a (id)"))
		case 3:
			stmt = fmt.Sprintf("CREATE TABLE %s%s LIKE %s", maybe("IF NOT EXISTS "), tbl(), tbl())
		case 4:
			stmt = fmt.Sprintf("CREATE TABLE %s AS SELECT 1 AS x", tbl())
		case 5, 6:
			stmt = fmt.Sprintf("DROP %sTABLE %s%s%s", maybe("TEMPORARY "), maybe("IF EXISTS "), tbl(), maybe(", "+tbl()))
		case 7:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s int%s", tbl(), col(), maybe(" AFTER "+col()))
		case 8:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP COLUMN%s %s%s", tbl(), ifExists(), col(), maybe(", ADD COLUMN "+col()+" int"))
		case 9:
			stmt = fmt.Sprintf("ALTER TABLE %s CHANGE %s %s int", tbl(), col(), col())
		case 10:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD %sINDEX %s (%s)", tbl(), maybe("UNIQUE "), idx(), col())
		case 11:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP INDEX%s %s%s", tbl(), ifExists(), idx(), maybe(", ADD INDEX "+idx()+" (x)"))
		case 12:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (x) REFERENCES %s (id)", tbl(), fk(), tbl())
		case 13:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY%s %s", tbl(), ifExists(), fk())
		case 14:
			stmt = fmt.Sprintf("ALTER TABLE %s %s", tbl(), pick("ADD PRIMARY KEY (x)", "DROP PRIMARY KEY", "RENAME INDEX i TO j", "RENAME TO "+tbl(), "ENGINE = InnoDB", "COMMENT 'x'"))
		case 15:
			stmt = fmt.Sprintf("RENAME TABLE %s TO %s%s", tbl(), tbl(), maybe(", "+tbl()+" TO "+tbl()))
		case 16:
			stmt = fmt.Sprintf("CREATE %s%sINDEX %s ON %s (%s)", pick("", "", "OR REPLACE "), maybe("UNIQUE "), idx(), tbl(), col())
		case 17:
			stmt = fmt.Sprintf("DROP INDEX %s ON %s", idx(), tbl())
		case 18:
			stmt = fmt.Sprintf("CREATE %sVIEW %s AS SELECT 1 AS x", maybe("OR REPLACE "), tbl())
		case 19:
			stmt = fmt.Sprintf("DROP VIEW %s%s%s", maybe("IF EXISTS "), tbl(), maybe(", "+tbl()))
		case 20:
			stmt = pick("CREATE DATABASE q", "CREATE DATABASE IF NOT EXISTS q", "CREATE OR REPLACE DATABASE q", "DROP DATABASE IF EXISTS q", "USE q", "USE shop")
		case 21:
			stmt = fmt.Sprintf("TRUNCATE TABLE %s", tbl())
		case 22:
			stmt = fmt.Sprintf("INSERT INTO %s (x) VALUES (1)", tbl())
		case 23:
			stmt = pick("SET foreign_key_checks = 0", "SET foreign_key_checks = 1")
		case 24:
			stmt = fmt.Sprintf("CREATE %sTRIGGER %s BEFORE INSERT ON %s FOR EACH ROW SET NEW.x = 1", maybe("OR REPLACE "), pick("tr", "TR", "tr2"), tbl())
		case 25:
			stmt = fmt.Sprintf("DROP TRIGGER %s%s", maybe("IF EXISTS "), pick("tr", "TR", "tr2"))
		case 26, 27:
			stmt = fmt.Sprintf("CREATE %s%sSEQUENCE %s%s", pick("", "", "OR REPLACE "), maybe("TEMPORARY "), pick("", "IF NOT EXISTS "), tbl())
			stmt = strings.Replace(stmt, "SEQUENCE IF NOT EXISTS", "SEQUENCE", strings.Count(stmt, "OR REPLACE"))
		case 28, 29:
			stmt = fmt.Sprintf("DROP %sSEQUENCE %s%s%s", maybe("TEMPORARY "), maybe("IF EXISTS "), tbl(), maybe(", "+tbl()))
		case 30:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT%s %s", tbl(), ifExists(), fk())
		case 31:
			stmt = fmt.Sprintf("CREATE TABLE %s (x int, y int) WITH SYSTEM VERSIONING", tbl())
		case 32:
			stmt = pick("SET SESSION system_versioning_alter_history = KEEP", "SET SESSION system_versioning_alter_history = ERROR")
		case 33:
			stmt = fmt.Sprintf("DELETE FROM %s", tbl())
		}
		b.WriteString(stmt)
		b.WriteString(";\n")
	}
	return b.String()
}
