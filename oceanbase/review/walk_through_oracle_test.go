package review

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"

	"github.com/bytebase/omni/oceanbase/parser"
)

// The cases, error codes, and random changes this tree's walk oracle shares
// with MySQL's (mysql/review), from which it was forked. They run on
// OceanBase: see walk_through_oceanbase_oracle_test.go.

// walkBase is the shop database the curated cases start from: tables with
// keys and foreign keys, one of a generated name, a table with a trigger,
// a view, a routine of each kind, and an event.
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
`

// walkCase is a change on the base, after setup when it has one. With
// complete, the walk must find every statement OceanBase rejects.
type walkCase struct {
	name     string
	setup    string
	change   string
	complete bool
	// unread syncs the target as a user without SELECT on the oceanbase
	// database does: with no materialized views.
	unread bool
	// refused is the error OceanBase answers the setup with, for a case
	// of what only MySQL has. The case ends there, and passes only while
	// OceanBase refuses its setup so.
	refused uint16
	// staleEvent is an event the sync must list though the database holds
	// none of the name: a dropped database's, which OceanBase lists for
	// some seconds under a database created again with its name.
	staleEvent string
}

// walkCases are the cases MySQL's oracle runs too. A case marked "Not
// complete" is complete on MySQL: OceanBase rejects a statement of it that
// the walk leaves unclaimed by design, which the comment names.
var walkCases = []walkCase{
	// Databases and USE.
	{name: "databases", change: "CREATE DATABASE shop;\nCREATE DATABASE q;\nCREATE DATABASE q;\nCREATE DATABASE IF NOT EXISTS q;\nDROP DATABASE IF EXISTS nosuch;\nUSE nosuch;\nCREATE TABLE q.x (a int);\nCREATE TABLE q.x (a int);\nDROP DATABASE q;\nCREATE TABLE q.y (a int);", complete: true},
	{name: "no database selected", change: "DROP DATABASE shop;\nALTER DATABASE CHARACTER SET utf8mb4;\nCREATE TABLE x (a int);\nCREATE PROCEDURE p2() SELECT 1;\nDROP PROCEDURE IF EXISTS p;\nDROP FUNCTION shop.f;\nDROP TRIGGER IF EXISTS tr;\nCREATE EVENT e2 ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nALTER EVENT e COMMENT 'x';", complete: true},
	// With no database selected, an unqualified DROP FUNCTION looks for a
	// loadable function alone: IF EXISTS is a note, and without it the
	// statement fails with 1305 where none is loaded.
	// Not complete: DROP FUNCTION IF EXISTS in a dropped database fails (1049).
	{name: "DROP FUNCTION with no database selected", change: "DROP DATABASE shop;\nDROP FUNCTION IF EXISTS nosuch;\nDROP FUNCTION IF EXISTS f;\nDROP PROCEDURE IF EXISTS p;\nDROP FUNCTION IF EXISTS shop.f;\nDROP VIEW IF EXISTS v;\nDROP EVENT IF EXISTS e;\nDROP TABLE IF EXISTS nokey;"},
	{name: "DROP FUNCTION with no database selected may find a loadable one", change: "DROP DATABASE shop;\nDROP FUNCTION nosuch;\nCREATE FUNCTION service_get_read_locks RETURNS INT SONAME 'locking_service.so';\nDROP FUNCTION service_get_read_locks;"},
	{name: "objects in a dropped database", change: "CREATE DATABASE z;\nDROP DATABASE z;\nCREATE PROCEDURE z.p() SELECT 1;\nCREATE FUNCTION z.f() RETURNS int DETERMINISTIC RETURN 1;\nCREATE VIEW z.w AS SELECT 1 AS a;\nCREATE TEMPORARY TABLE z.tmp (a int);", complete: true},
	{name: "USE moves unqualified names", change: "CREATE DATABASE q;\nUSE q;\nCREATE TABLE t (a int);\nCREATE TABLE t (a int);\nUSE shop;\nCREATE TABLE t (a int);", complete: true},
	// Tables.
	{name: "create a table that exists", change: "CREATE TABLE t (a int);\nCREATE TABLE v (a int);\nCREATE TABLE IF NOT EXISTS t (a int);\nCREATE TABLE IF NOT EXISTS v (a int);\nCREATE TABLE IF NOT EXISTS v LIKE nosuch;\nCREATE TABLE x LIKE nosuch;", complete: true},
	{name: "a new table's text", change: "CREATE TABLE x (a int, a int);\nCREATE TABLE x (a int, b int, KEY k (a), KEY k (b));\nCREATE TABLE x (a int PRIMARY KEY, b int, PRIMARY KEY (b));\nCREATE TABLE x (a int, KEY `PRIMARY` (a));", complete: true},
	{name: "a new table holds what it declares", change: "CREATE TABLE x (a int, b int UNIQUE, KEY k (a), CONSTRAINT fx FOREIGN KEY (a) REFERENCES r (id));\nALTER TABLE x ADD COLUMN a int;\nALTER TABLE x ADD INDEX k (b);\nALTER TABLE x DROP COLUMN zz;\nALTER TABLE x DROP INDEX k, DROP INDEX k;\nALTER TABLE nokey ADD CONSTRAINT fx FOREIGN KEY (a) REFERENCES r (id);\nCREATE TABLE y LIKE x;\nALTER TABLE y DROP INDEX k;", complete: true},
	{name: "names too long", change: "CREATE TABLE " + strings.Repeat("a", 65) + " (a int);\nALTER TABLE t ADD COLUMN " + strings.Repeat("c", 65) + " int;\nDROP TABLE " + strings.Repeat("d", 65) + ";", complete: true},
	// ALTER TABLE.
	// Not complete: pairs on one name and an index dropped and added (1060, 4179).
	{name: "ALTER TABLE drops first", change: "ALTER TABLE t ADD COLUMN a int, DROP COLUMN a;\nALTER TABLE t DROP INDEX i, ADD INDEX i (b);\nALTER TABLE t DROP COLUMN a;\nALTER TABLE t DROP COLUMN a;\nALTER TABLE t RENAME COLUMN b TO c, ADD COLUMN b int;\nALTER TABLE t DROP COLUMN zz, ADD COLUMN y int;"},
	// Not complete: ALTER TABLE of a view, and ADD INDEX `PRIMARY`, end the walk.
	{name: "ALTER TABLE subcommands", change: "ALTER TABLE nosuch ADD COLUMN a int;\nALTER TABLE v ADD COLUMN a int;\nALTER TABLE t ADD COLUMN A int;\nALTER TABLE t ADD COLUMN c int AFTER nosuch;\nALTER TABLE t DROP COLUMN nosuch;\nALTER TABLE t MODIFY COLUMN nosuch int;\nALTER TABLE t CHANGE nosuch c int;\nALTER TABLE t CHANGE a b int;\nALTER TABLE t ALTER COLUMN nosuch SET DEFAULT 1;\nALTER TABLE t ADD INDEX I (b);\nALTER TABLE t ADD INDEX `PRIMARY` (b);\nALTER TABLE t DROP INDEX nosuch;\nALTER TABLE t ADD PRIMARY KEY (a);\nALTER TABLE t RENAME INDEX nosuch TO j;\nALTER TABLE t RENAME INDEX i TO PRIMARY2;\nALTER TABLE t ALTER INDEX nosuch INVISIBLE;"},
	// Not complete: DROP INDEX `PRIMARY` ends the walk.
	{name: "repeated subcommands", change: "ALTER TABLE t ADD COLUMN c int, ADD COLUMN c int;\nALTER TABLE t DROP COLUMN a, DROP COLUMN a;\nALTER TABLE t ADD INDEX k (a), ADD INDEX k (b);\nALTER TABLE nokey ADD PRIMARY KEY (a), ADD PRIMARY KEY (a);\nALTER TABLE t DROP INDEX i, DROP INDEX i;\nALTER TABLE t DROP FOREIGN KEY fk1, DROP FOREIGN KEY fk1;\nALTER TABLE t DROP PRIMARY KEY, DROP PRIMARY KEY;\nALTER TABLE t DROP PRIMARY KEY, DROP INDEX `PRIMARY`;"},
	{name: "foreign keys", change: "CREATE TABLE y0 (a int, FOREIGN KEY (a) REFERENCES v (id));\nCREATE TABLE y00 (a int, CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id));\nALTER TABLE nokey DROP FOREIGN KEY fk1;\nALTER TABLE nokey ADD CONSTRAINT FK1 FOREIGN KEY (a) REFERENCES r (id);\nALTER TABLE t DROP FOREIGN KEY nosuch;\nALTER TABLE t DROP FOREIGN KEY fk1;\nALTER TABLE t DROP FOREIGN KEY fk1;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET foreign_key_checks = 0;\nCREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET foreign_key_checks = 1;\nCREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nCREATE TABLE w (id int PRIMARY KEY, p int, FOREIGN KEY (p) REFERENCES w (id));", complete: true},
	// The CONSTRAINT symbol names a foreign key; the name after FOREIGN KEY
	// names only its index, and without a symbol the foreign key's name is
	// generated. A long name there counts only without a symbol.
	// Not complete: a foreign key adds no index (1091), two in one ALTER TABLE (1235).
	{name: "a foreign key's name and its index's", change: "CREATE TABLE x (a int, CONSTRAINT cx FOREIGN KEY fx (a) REFERENCES r (id));\nALTER TABLE x DROP FOREIGN KEY cx;\nALTER TABLE x DROP FOREIGN KEY cx;\nALTER TABLE x DROP INDEX cx;\nCREATE TABLE y (a int, FOREIGN KEY fy (a) REFERENCES r (id));\nALTER TABLE y DROP FOREIGN KEY y_ibfk_1;\nALTER TABLE y DROP INDEX fy;\nCREATE TABLE z (a int, CONSTRAINT fy FOREIGN KEY (a) REFERENCES r (id));\nALTER TABLE nokey ADD CONSTRAINT nk FOREIGN KEY fk1 (a) REFERENCES r (id);\nALTER TABLE nokey ADD CONSTRAINT nk2 FOREIGN KEY fk1 (a) REFERENCES r (id), ADD CONSTRAINT nk3 FOREIGN KEY fk1 (a) REFERENCES r (id);\nCREATE TABLE w (a int, CONSTRAINT cw FOREIGN KEY " + strings.Repeat("f", 65) + " (a) REFERENCES r (id));\nCREATE TABLE w2 (a int, FOREIGN KEY " + strings.Repeat("f", 65) + " (a) REFERENCES r (id));\nCREATE TABLE w3 (a int, CONSTRAINT fk1 FOREIGN KEY fw (a) REFERENCES r (id));"},
	// An ALTER TABLE that also renames its table resolves a parent by rules
	// of its own: in the database the table moves to, and the old name as
	// the table itself.
	{name: "a parent beside a RENAME", change: "CREATE DATABASE q;\nCREATE TABLE q.p2 (id int PRIMARY KEY);\nCREATE TABLE m (a int);\nALTER TABLE m RENAME TO q.m, ADD FOREIGN KEY (a) REFERENCES p2 (id);\nCREATE TABLE m2 (id int PRIMARY KEY, a int);\nALTER TABLE m2 RENAME TO m3, ADD FOREIGN KEY (a) REFERENCES m2 (id);\nCREATE TABLE m4 (id int PRIMARY KEY, a int);\nALTER TABLE m4 RENAME TO m5, ADD FOREIGN KEY (a) REFERENCES m5 (id);"},
	{name: "an unqualified parent is in the child's database", change: "CREATE DATABASE q;\nCREATE TABLE q.parent (id int PRIMARY KEY);\nCREATE TABLE q.child (id int, FOREIGN KEY (id) REFERENCES parent (id));\nALTER TABLE q.child ADD CONSTRAINT c2 FOREIGN KEY (id) REFERENCES parent (id);\nCREATE TABLE q.orphan (id int, FOREIGN KEY (id) REFERENCES nosuch (id));", complete: true},
	// An engine other than InnoDB parses foreign keys and ignores them.
	{name: "an engine that ignores foreign keys", setup: "CREATE TABLE my (id int PRIMARY KEY, a int) ENGINE=MyISAM;", change: "ALTER TABLE my DROP FOREIGN KEY nosuch;\nALTER TABLE my ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES nosuch (id);\nCREATE TABLE m2 (a int, CONSTRAINT mf FOREIGN KEY (a) REFERENCES nosuch (id), CONSTRAINT mf FOREIGN KEY (a) REFERENCES r (id)) ENGINE=MyISAM;\nCREATE TABLE ok (a int, CONSTRAINT mf FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE m3 LIKE my;\nALTER TABLE m3 ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id);\nCREATE TABLE inn (id int PRIMARY KEY, a int);\nALTER TABLE inn ENGINE=MyISAM;\nALTER TABLE inn ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES nosuch (id);\nALTER TABLE inn DROP FOREIGN KEY whatever;\nALTER TABLE inn ENGINE=InnoDB;\nALTER TABLE inn DROP FOREIGN KEY whatever;"},
	// Not complete: pairs on one name (1060, 1061, 1176).
	{name: "column names a statement adds", setup: "CREATE TABLE cols (a int, b int, c int);", change: "ALTER TABLE cols RENAME COLUMN a TO z, RENAME COLUMN b TO z;\nALTER TABLE cols CHANGE a z int, CHANGE b z int;\nALTER TABLE cols ADD COLUMN z int, RENAME COLUMN a TO z;\nALTER TABLE cols CHANGE a z int, ADD COLUMN Z int;\nALTER TABLE cols RENAME COLUMN a TO b, RENAME COLUMN b TO a;\nALTER TABLE cols ADD COLUMN b2 int AFTER b;\nALTER TABLE cols RENAME COLUMN a TO b, RENAME COLUMN b TO y;\nALTER TABLE cols ADD COLUMN y2 int AFTER y, ADD COLUMN b3 int AFTER b;\nALTER TABLE cols ADD COLUMN a int;\nALTER TABLE cols ADD KEY i (a), ADD KEY j (b);\nALTER TABLE cols RENAME INDEX i TO j, RENAME INDEX j TO i;\nALTER TABLE cols DROP INDEX i, DROP INDEX j;\nALTER TABLE cols ADD KEY i (a), ADD KEY j (b);\nALTER TABLE cols RENAME INDEX i TO j, DROP INDEX j;\nALTER TABLE cols DROP INDEX j;\nALTER TABLE cols DROP INDEX j;\nALTER TABLE cols CHANGE a A int, ADD COLUMN d int;\nALTER TABLE cols MODIFY a int PRIMARY KEY, MODIFY b int PRIMARY KEY;\nALTER TABLE cols CHANGE a a int PRIMARY KEY, ADD PRIMARY KEY (b);\nALTER TABLE cols MODIFY a int NOT NULL PRIMARY KEY, ADD COLUMN e int PRIMARY KEY;"},
	// Not complete: a parent a temporary table shadows fails too (1215).
	{name: "a temporary table is no parent", change: "CREATE TEMPORARY TABLE tp (id int PRIMARY KEY);\nCREATE TABLE y1 (a int, FOREIGN KEY (a) REFERENCES tp (id));\nCREATE TEMPORARY TABLE r (id int PRIMARY KEY);\nCREATE TABLE y2 (a int, FOREIGN KEY (a) REFERENCES r (id));"},
	{name: "foreign keys named twice", change: "CREATE TABLE x (a int, b int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT f FOREIGN KEY (b) REFERENCES r (id));\nCREATE TABLE x (a int, b int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT F FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE x (a int, b int);\nALTER TABLE x ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), ADD CONSTRAINT F FOREIGN KEY (b) REFERENCES r (id);", complete: true},
	// RENAME TABLE renames the foreign keys named <table>_ibfk_...
	{name: "renamed generated foreign key names", change: "CREATE TABLE g (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));\nRENAME TABLE g TO h;\nALTER TABLE nokey ADD CONSTRAINT g_ibfk_1 FOREIGN KEY (a) REFERENCES r (id);\nALTER TABLE h DROP FOREIGN KEY h_ibfk_1;\nALTER TABLE h DROP FOREIGN KEY h_ibfk_1;\nRENAME TABLE child TO kid;\nALTER TABLE kid DROP FOREIGN KEY kid_ibfk_1;\nALTER TABLE kid DROP FOREIGN KEY kid_ibfk_1;\nALTER TABLE kid RENAME TO c2;\nCREATE TABLE child (id int PRIMARY KEY, b int, CONSTRAINT child_ibfk_1 FOREIGN KEY (b) REFERENCES r (id));", complete: true},
	{name: "a table with no primary key may get an invisible one", change: "CREATE TABLE g (a int);\nALTER TABLE g DROP PRIMARY KEY;\nALTER TABLE g ADD PRIMARY KEY (a);"},
	// OceanBase has no sql_generate_invisible_primary_key.
	{name: "the sync may not list an invisible primary key", setup: "SET GLOBAL show_gipk_in_create_table_and_information_schema = OFF;\nSET SESSION sql_generate_invisible_primary_key = ON;\nCREATE TABLE gk (a int);", change: "ALTER TABLE gk ADD COLUMN x int AFTER my_row_id;\nALTER TABLE gk ADD COLUMN my_row_id int;\nALTER TABLE gk ADD PRIMARY KEY (a);", refused: 1193},
	{name: "a view stays in its database", change: "CREATE DATABASE q;\nRENAME TABLE v TO q.v;\nDROP VIEW v;", complete: true},
	{name: "RENAME", change: "RENAME TABLE nosuch TO x;\nRENAME TABLE t TO r;\nRENAME TABLE t TO tmp, r TO t, tmp TO r;\nALTER TABLE r ADD COLUMN a int;\nALTER TABLE t RENAME TO r;\nALTER TABLE t RENAME TO t2;\nALTER TABLE t ADD COLUMN zz int;\nCREATE DATABASE q;\nRENAME TABLE nokey TO q.nokey;\nCREATE TABLE q.nokey (a int);", complete: true},
	// DROP and TRUNCATE.
	{name: "DROP TABLE IF EXISTS leaves a view", setup: "CREATE DATABASE q;\nCREATE VIEW q.v AS SELECT 1 AS a;", change: "DROP TABLE IF EXISTS q.v;\nDROP VIEW q.v;\nDROP VIEW q.v;"},
	{name: "DROP TABLE", change: "DROP TABLE t, nosuch;\nDROP TABLE v;\nDROP TABLE IF EXISTS nosuch, nokey;\nDROP TABLE nokey;\nDROP VIEW t;\nDROP VIEW nosuch;\nDROP VIEW IF EXISTS nosuch, v;\nDROP VIEW v;\nTRUNCATE TABLE nosuch;\nTRUNCATE TABLE v;", complete: true},
	// Not complete: IF NOT EXISTS of an existing one checks its text (1060).
	{name: "a temporary table's definition", change: "CREATE TEMPORARY TABLE tmp (a int, a int);\nCREATE TEMPORARY TABLE tmp (a int, KEY k (a), KEY k (a));\nCREATE TEMPORARY TABLE tmp (a int PRIMARY KEY, PRIMARY KEY (a));\nCREATE TEMPORARY TABLE tmp (a int);\nCREATE TEMPORARY TABLE IF NOT EXISTS tmp (a int, a int);\nCREATE TEMPORARY TABLE tmp (a int);"},
	{name: "temporary tables outlive their database", change: "CREATE DATABASE q;\nCREATE TEMPORARY TABLE q.tmp (a int);\nDROP DATABASE q;\nCREATE DATABASE q;\nDROP TEMPORARY TABLE q.tmp;\nDROP TEMPORARY TABLE q.tmp;", complete: true},
	{name: "temporary tables", change: "CREATE TEMPORARY TABLE t (z int);\nALTER TABLE nokey ADD COLUMN z int;\nINSERT INTO t (z) VALUES (1);\nDROP TABLE t;\nINSERT INTO t (id) VALUES (1);\nDROP TEMPORARY TABLE t;\nDROP TEMPORARY TABLE IF EXISTS t;\nCREATE TEMPORARY TABLE tt (a int);\nCREATE TEMPORARY TABLE tt (a int);\nTRUNCATE TABLE tt;", complete: true},
	// Indexes, views, triggers, routines, events.
	// Not complete: CREATE INDEX `PRIMARY` ends the walk.
	{name: "CREATE and DROP INDEX", change: "CREATE INDEX k ON nosuch (a);\nCREATE INDEX k ON v (id);\nCREATE INDEX i ON t (b);\nCREATE INDEX `PRIMARY` ON t (b);\nDROP INDEX nosuch ON t;\nDROP INDEX i ON nosuch;\nDROP INDEX i ON t;\nDROP INDEX i ON t;"},
	{name: "views", change: "CREATE VIEW t AS SELECT 1 AS a;\nCREATE VIEW v AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW t AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW v AS SELECT 1 AS a;\nALTER VIEW nosuch AS SELECT 1 AS a;\nALTER VIEW t AS SELECT 1 AS a;\nDROP VIEW IF EXISTS t;\nDROP VIEW t;", complete: true},
	{name: "triggers", change: "CREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;\nCREATE TRIGGER tr2 BEFORE INSERT ON nosuch FOR EACH ROW SET NEW.a = 1;\nDROP TRIGGER nosuch;\nDROP TRIGGER tr;\nDROP TRIGGER tr;\nCREATE TRIGGER TR BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;", complete: true},
	// Not complete: a listed event, an event name with a period.
	{name: "routines and events", change: "CREATE PROCEDURE p() SELECT 1;\nCREATE PROCEDURE P() SELECT 1;\nCREATE FUNCTION p() RETURNS int DETERMINISTIC RETURN 1;\nCREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;\nDROP PROCEDURE nosuch;\nDROP FUNCTION shop.nosuch;\nALTER PROCEDURE nosuch COMMENT 'x';\nCREATE EVENT E ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nDROP EVENT nosuch;\nALTER EVENT nosuch COMMENT 'x';\nALTER EVENT e RENAME TO e2;\nDROP EVENT e;\nCREATE EVENT `shop.e2` ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nDROP EVENT `shop.e2`;\nDROP EVENT `shop.e2`;"},
	// An unqualified DROP FUNCTION drops a loadable function of the name
	// first, which the sync does not list.
	// Not complete: none can be created (1235).
	{name: "loadable functions", change: "CREATE FUNCTION service_get_read_locks RETURNS INT SONAME 'locking_service.so';\nDROP FUNCTION service_get_read_locks;\nDROP FUNCTION shop.service_get_read_locks;\nCREATE FUNCTION service_get_read_locks RETURNS INT SONAME 'locking_service.so';\nCREATE FUNCTION service_get_read_locks() RETURNS int DETERMINISTIC RETURN 1;\nDROP FUNCTION service_get_read_locks;\nDROP FUNCTION shop.service_get_read_locks;\nDROP FUNCTION shop.service_get_read_locks;"},
	// DML.
	{name: "DML targets", change: "INSERT INTO nosuch VALUES (1);\nUPDATE nosuch SET a = 1;\nDELETE FROM nosuch;\nREPLACE INTO nosuch VALUES (1);\nINSERT INTO v VALUES (1);\nINSERT INTO nokey VALUES (1);", complete: true},
	// Name case: tables and databases fold only under lower_case_table_names
	// 1; the other spaces always fold, triggers aside.
	{name: "name case", change: "CREATE TABLE T (a int);\nDROP TABLE NOKEY;\nCREATE DATABASE Q;\nCREATE DATABASE q;\nCREATE TABLE Q.x (a int);\nCREATE TABLE q.X (a int);\nALTER TABLE R ADD COLUMN ID int;\nCREATE VIEW V AS SELECT 1 AS a;\nDROP PROCEDURE P;\nDROP FUNCTION F;\nDROP EVENT E;\nDROP TRIGGER TR;\nCREATE TRIGGER Tr BEFORE INSERT ON r FOR EACH ROW SET NEW.id = 1;", complete: true},
	// Names outside ASCII: MySQL compares them by rules of its own per space
	// (spec §2), and takes some for ASCII names, so the walk follows none.
	// CAFé is café as a column and KÉ is ké as an index; the Kelvin sign is
	// k as a column and an index, and é is e as a routine, an event, and a
	// trigger.
	{name: "names outside ASCII", change: "CREATE TABLE x (id int, café int, KEY ké (id));\nALTER TABLE x DROP COLUMN CAFé;\nALTER TABLE x DROP INDEX KÉ;\nALTER TABLE x DROP INDEX Ké;\nALTER TABLE x DROP COLUMN zz;"},
	{name: "an ASCII name may be one outside ASCII", setup: "CREATE TABLE kel (`\u212a` int, a int, KEY `\u212a` (a));\nCREATE PROCEDURE `pé`() SELECT 1;\nCREATE FUNCTION `fé`() RETURNS int DETERMINISTIC RETURN 1;\nCREATE EVENT `eé` ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nCREATE TRIGGER `tré` BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;", change: "ALTER TABLE kel DROP COLUMN k;\nALTER TABLE kel DROP INDEX k;\nDROP PROCEDURE pe;\nDROP FUNCTION shop.fe;\nDROP EVENT ee;\nDROP TRIGGER tre;\nALTER TABLE kel DROP COLUMN zz;\nDROP PROCEDURE nosuch;\nALTER TABLE t DROP COLUMN zz;"},
	// Each assignment of a SET has a scope of its own (spec §4): an
	// unprefixed one the statement's, an @@session or bare @@ one the
	// session's. A global assignment leaves the session's checks as they
	// are; emptyOceanBase sets the global back.
	{name: "SET foreign_key_checks in each assignment's scope", change: "SET GLOBAL foreign_key_checks = 1, @@session.foreign_key_checks = 0;\nCREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET GLOBAL foreign_key_checks = 1, @@foreign_key_checks = 1;\nCREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET @@global.foreign_key_checks = 1, foreign_key_checks = 0;\nCREATE TABLE z2 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET GLOBAL foreign_key_checks = 1;\nCREATE TABLE z3 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET foreign_key_checks = 1;\nSET GLOBAL foreign_key_checks = 0;\nCREATE TABLE z4 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET GLOBAL foreign_key_checks = 1, foreign_key_checks = 0;\nCREATE TABLE w (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", complete: true},
	{name: "SET foreign_key_checks under a later scope", change: "SET @a = 1, GLOBAL foreign_key_checks = 0;\nCREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET GLOBAL foreign_key_checks = 1, SESSION foreign_key_checks = 0;\nCREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET SESSION foreign_key_checks = 1, GLOBAL foreign_key_checks = 1, foreign_key_checks = 0;\nCREATE TABLE w (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", complete: true},
	// What the walk leaves Unknown by design (spec §6): databases other than
	// the target, the tables a query reads, generated index names, foreign
	// key names held per database, a primary key GIPK may have added, a key
	// two parts of one ALTER TABLE both change, and a function an
	// unqualified DROP FUNCTION may find loaded.
	{name: "unknown by design", change: "DROP FUNCTION nosuch;\nDROP DATABASE nosuch;\nRENAME TABLE t TO nodb.t;\nCREATE VIEW nodb.w AS SELECT 1 AS a;\nCREATE TABLE x AS SELECT * FROM nosuch;\nCREATE TABLE x (a int, b int UNIQUE);\nALTER TABLE x DROP INDEX zz;\nCREATE TABLE y (a int, CONSTRAINT fy FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE z LIKE y;\nALTER TABLE z DROP FOREIGN KEY fy;\nALTER TABLE nokey DROP PRIMARY KEY;\nALTER TABLE t ADD COLUMN z int, DROP COLUMN z;\nALTER TABLE t ADD INDEX j (b), DROP INDEX j;\nALTER TABLE t ADD COLUMN x int, RENAME COLUMN x TO y;"},
	// A trigger cannot create a temporary table on OceanBase.
	{name: "a trigger of another database may create a temporary table", setup: "DROP TRIGGER tr;\nCREATE DATABASE q;\nCREATE TABLE q.o (a int);\nCREATE TRIGGER q.ot BEFORE INSERT ON q.o FOR EACH ROW CREATE TEMPORARY TABLE IF NOT EXISTS shop.made (a int);", change: "INSERT INTO q.o VALUES (1);\nDROP TEMPORARY TABLE made;\nDROP TEMPORARY TABLE made;", refused: 1422},
	// A function cannot create a temporary table on OceanBase.
	{name: "a query's function may create a temporary table", setup: "CREATE FUNCTION mk() RETURNS int DETERMINISTIC BEGIN CREATE TEMPORARY TABLE IF NOT EXISTS made (a int); RETURN 1; END;", change: "CREATE TABLE c AS SELECT mk() AS x;\nDROP TEMPORARY TABLE made;\nDROP TEMPORARY TABLE made;", refused: 1422},
	// A stored function may create a temporary table when a view's query
	// calls it, and an unqualified call after USE runs one of that
	// database. A call before a base table read in one query counts too.
	// A function cannot create a temporary table on OceanBase.
	{name: "a function a view or another database holds may create a temporary table", setup: "CREATE DATABASE q;\nCREATE FUNCTION q.mk() RETURNS int DETERMINISTIC BEGIN CREATE TEMPORARY TABLE IF NOT EXISTS shop.made (a int); RETURN 1; END;\nCREATE FUNCTION q.mk2() RETURNS int DETERMINISTIC BEGIN CREATE TEMPORARY TABLE IF NOT EXISTS shop.made2 (a int); RETURN 1; END;\nCREATE FUNCTION q.mk3() RETURNS int DETERMINISTIC BEGIN CREATE TEMPORARY TABLE IF NOT EXISTS shop.made3 (a int); RETURN 1; END;\nCREATE FUNCTION q.mk4() RETURNS int DETERMINISTIC BEGIN CREATE TEMPORARY TABLE IF NOT EXISTS shop.made4 (a int); RETURN 1; END;\nCREATE VIEW vv AS SELECT q.mk() AS x;\nCREATE VIEW q.qv AS SELECT q.mk3() AS x;\nINSERT INTO nokey VALUES (1);", change: "SELECT q.mk4() FROM nokey;\nDROP TEMPORARY TABLE made4;\nSELECT * FROM vv;\nDROP TEMPORARY TABLE made;\nUSE q;\nSELECT mk2();\nUSE shop;\nDROP TEMPORARY TABLE made2;\nINSERT INTO nokey SELECT x FROM q.qv;\nDROP TEMPORARY TABLE made3;\nDROP TEMPORARY TABLE made3;", refused: 1422},
	{name: "a trigger may create a temporary table", change: "CREATE TRIGGER tt BEFORE INSERT ON nokey FOR EACH ROW CREATE TEMPORARY TABLE IF NOT EXISTS shadow (a int);\nINSERT INTO nokey VALUES (1);\nINSERT INTO shadow VALUES (1);\nDROP TABLE shadow;"},
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

// walkCodes maps a finding's message to the MySQL errors it stands for
// (5.7's codes included). Where one text stands for several statements'
// errors, the statement it is on (leads) tells them apart: DROP TABLE and
// TRUNCATE refuse a view as a missing table, and ALTER VIEW reports a
// missing view as a missing table.
var walkCodes = []struct {
	match func(m, stmt string) bool
	codes []string
}{
	{func(m, _ string) bool {
		return strings.HasPrefix(m, "Temporary table ") && strings.HasSuffix(m, " already exists.")
	}, []string{"1050"}},
	{func(m, _ string) bool { return strings.HasPrefix(m, "Temporary table ") }, []string{"1051"}},
	{func(m, _ string) bool {
		return strings.HasPrefix(m, "Referenced table ") || strings.HasSuffix(m, ", not a table a foreign key can reference.")
	}, []string{"1824", "1215"}},
	{func(m, _ string) bool { return saysTwice(m, "foreign key", "declares", "adds") }, []string{"1826", "1061"}},
	{func(m, _ string) bool { return strings.HasPrefix(m, "MySQL doesn't accept `RENAME TABLE` of view ") }, []string{"1450"}},
	{func(m, _ string) bool {
		return strings.HasPrefix(m, "Foreign key ") && strings.HasSuffix(m, " already exists.")
	}, []string{"1826", "1022"}},
	{func(m, _ string) bool {
		return strings.HasPrefix(m, "Foreign key ") || saysTwice(m, "foreign key", "drops")
	}, []string{"1091"}},
	{func(m, _ string) bool {
		return (strings.HasPrefix(m, "Table ") || strings.HasPrefix(m, "View ")) && strings.HasSuffix(m, " already exists.") ||
			strings.HasPrefix(m, "The name ") && strings.HasSuffix(m, " is already taken.")
	}, []string{"1050"}},
	{func(m, _ string) bool { return strings.HasSuffix(m, " already has a primary key.") }, []string{"1068"}},
	{func(m, _ string) bool { return strings.HasSuffix(m, " has no primary key.") }, []string{"1091"}},
	{func(m, stmt string) bool {
		return strings.HasSuffix(m, ", not a table.") && leads(stmt, "DROP TABLE", "DROP TABLES", "TRUNCATE")
	}, []string{"1146", "1051", "1017"}},
	{func(m, _ string) bool {
		return (strings.HasPrefix(m, "Table ") || strings.HasPrefix(m, "`")) && strings.HasSuffix(m, " doesn't exist.")
	}, []string{"1146", "1051", "1017"}},
	{func(m, stmt string) bool { return strings.HasPrefix(m, "View ") && leads(stmt, "ALTER") }, []string{"1146", "1051", "1017"}},
	{func(m, _ string) bool { return strings.HasPrefix(m, "View ") }, []string{"1051"}},
	{func(m, _ string) bool { return saysWrongKind(m) }, []string{"1347"}},
	{func(m, _ string) bool {
		return strings.HasPrefix(m, "Database ") && strings.HasSuffix(m, " already exists.")
	}, []string{"1007"}},
	{func(m, _ string) bool { return strings.HasPrefix(m, "Database ") }, []string{"1049", "1008", "1025"}},
	{func(m, _ string) bool { return strings.HasPrefix(m, "No database is selected for ") }, []string{"1046"}},
	{func(m, _ string) bool {
		return strings.HasPrefix(m, "Column ") && strings.Contains(m, " already exists in ") || saysTwice(m, "column", "declares", "adds", "drops")
	}, []string{"1060", "1091"}},
	{func(m, _ string) bool { return strings.HasPrefix(m, "Column ") }, []string{"1054", "1091"}},
	{func(m, _ string) bool {
		return strings.HasPrefix(m, "Index ") && strings.Contains(m, " already exists in ") || saysTwice(m, "index", "declares", "adds")
	}, []string{"1061"}},
	{func(m, _ string) bool { return m == reservedPrimary }, []string{"1280"}},
	{func(m, _ string) bool {
		return strings.HasPrefix(m, "Index ") || saysTwice(m, "index", "drops") || m == "This statement drops the primary key twice."
	}, []string{"1091", "1176"}},
	{func(m, _ string) bool {
		return strings.HasPrefix(m, "Trigger ") && strings.HasSuffix(m, " already exists.")
	}, []string{"1359"}},
	{func(m, _ string) bool { return strings.HasPrefix(m, "Trigger ") }, []string{"1360"}},
	{func(m, _ string) bool {
		return (strings.HasPrefix(m, "Procedure ") || strings.HasPrefix(m, "Function ")) && strings.HasSuffix(m, " already exists.")
	}, []string{"1304"}},
	{func(m, _ string) bool { return strings.HasPrefix(m, "Procedure ") || strings.HasPrefix(m, "Function ") }, []string{"1305"}},
	{func(m, _ string) bool {
		return strings.HasPrefix(m, "Event ") && strings.HasSuffix(m, " already exists.")
	}, []string{"1537"}},
	{func(m, _ string) bool { return strings.HasPrefix(m, "Event ") }, []string{"1539"}},
	{func(m, _ string) bool {
		return strings.HasPrefix(m, "The name ") && strings.HasSuffix(m, "-character limit.")
	}, []string{"1059"}},
}

// saysTwice reports whether a message says the statement names kind twice, by
// one of verbs.
func saysTwice(m, kind string, verbs ...string) bool {
	for _, v := range verbs {
		if strings.HasPrefix(m, "This statement "+v+" "+kind+" ") && strings.HasSuffix(m, " twice.") {
			return true
		}
	}
	return false
}

// saysWrongKind reports whether a message says a relation is of another kind
// than the statement needs.
func saysWrongKind(m string) bool {
	for _, need := range []string{"table", "view", "sequence"} {
		if strings.HasSuffix(m, ", not a "+need+".") || strings.HasSuffix(m, " isn't a "+need+".") {
			return true
		}
	}
	return false
}

// leads reports whether a statement starts with one of phrases, each a
// run of keywords.
func leads(stmt string, phrases ...string) bool {
	f := strings.Fields(strings.ToUpper(stmt))
	for _, p := range phrases {
		if w := strings.Fields(p); len(f) >= len(w) && slices.Equal(f[:len(w)], w) {
			return true
		}
	}
	return false
}

// sweepChange is a sweep's change of n statements, over a small pool of
// names that collide often, case included.
func sweepChange(rng *rand.Rand, n int) string {
	// generated has the form of the name OceanBase generates for a foreign
	// key of table a.
	const generated = "a_obfk_1"
	pick := func(xs ...string) string { return xs[rng.Intn(len(xs))] }
	maybe := func(s string) string {
		if rng.Intn(3) == 0 {
			return s
		}
		return ""
	}
	db := func() string { return pick("", "", "", "q.") }
	tbl := func() string { return db() + pick("a", "b", "c", "A", "v", generated) }
	col := func() string { return pick("x", "y", "z", "X") }
	idx := func() string { return pick("i", "j", "x", "I") }
	fk := func() string { return pick("f", "F", generated, "g") }
	// fkDef is a foreign key's head: a CONSTRAINT symbol, the name after
	// FOREIGN KEY, or both, which name it apart.
	fkDef := func() string {
		return pick("CONSTRAINT "+fk()+" FOREIGN KEY", "CONSTRAINT "+fk()+" FOREIGN KEY "+fk(), "FOREIGN KEY "+fk())
	}
	// unique is a UNIQUE key's head, with a CONSTRAINT symbol at times.
	unique := func() string { return pick("", "UNIQUE ", "CONSTRAINT "+idx()+" UNIQUE ") }
	var b strings.Builder
	for range n {
		var stmt string
		switch rng.Intn(26) {
		case 0, 1:
			stmt = fmt.Sprintf("CREATE %sTABLE %s%s (id int PRIMARY KEY, x int%s, y int)", maybe("TEMPORARY "), maybe("IF NOT EXISTS "), tbl(), maybe(" UNIQUE"))
		case 2:
			stmt = fmt.Sprintf("CREATE TABLE %s (x int, y int, KEY %s (x)%s)", tbl(), idx(), maybe(", "+fkDef()+" (y) REFERENCES a (id)"))
		case 3:
			stmt = fmt.Sprintf("CREATE TABLE %s%s LIKE %s", maybe("IF NOT EXISTS "), tbl(), tbl())
		case 4:
			stmt = fmt.Sprintf("CREATE TABLE %s AS SELECT 1 AS x", tbl())
		case 5, 6:
			stmt = fmt.Sprintf("DROP %sTABLE %s%s%s", maybe("TEMPORARY "), maybe("IF EXISTS "), tbl(), maybe(", "+tbl()))
		case 7:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s int%s", tbl(), col(), maybe(" AFTER "+col()))
		case 8:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s%s", tbl(), col(), maybe(", ADD COLUMN "+col()+" int"))
		case 9:
			stmt = fmt.Sprintf("ALTER TABLE %s CHANGE %s %s int", tbl(), col(), col())
		case 10:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD %sINDEX %s (%s)", tbl(), unique(), idx(), col())
		case 11:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP INDEX %s%s", tbl(), idx(), maybe(", ADD INDEX "+idx()+" (x)"))
		case 12:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD %s (x) REFERENCES %s (id)", tbl(), fkDef(), tbl())
		case 13:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s", tbl(), fk())
		case 14:
			stmt = fmt.Sprintf("ALTER TABLE %s %s", tbl(), pick("ADD PRIMARY KEY (x)", "DROP PRIMARY KEY", "RENAME INDEX i TO j", "RENAME TO "+tbl(), "ENGINE = InnoDB"))
		case 15:
			stmt = fmt.Sprintf("RENAME TABLE %s TO %s%s", tbl(), tbl(), maybe(", "+tbl()+" TO "+tbl()))
		case 16:
			stmt = fmt.Sprintf("CREATE %sINDEX %s ON %s (%s)", maybe("UNIQUE "), idx(), tbl(), col())
		case 17:
			stmt = fmt.Sprintf("DROP INDEX %s ON %s", idx(), tbl())
		case 18:
			stmt = fmt.Sprintf("CREATE %sVIEW %s AS SELECT 1 AS x", maybe("OR REPLACE "), tbl())
		case 19:
			stmt = fmt.Sprintf("DROP VIEW %s%s", maybe("IF EXISTS "), tbl())
		case 20:
			stmt = pick("CREATE DATABASE q", "CREATE DATABASE IF NOT EXISTS q", "DROP DATABASE IF EXISTS q", "USE q", "USE shop")
		case 21:
			stmt = fmt.Sprintf("TRUNCATE TABLE %s", tbl())
		case 22:
			stmt = fmt.Sprintf("INSERT INTO %s (x) VALUES (1)", tbl())
		case 23:
			stmt = pick("SET foreign_key_checks = 0", "SET foreign_key_checks = 1", "SET GLOBAL foreign_key_checks = 0", "SET GLOBAL foreign_key_checks = 1, @@session.foreign_key_checks = 0", "SET GLOBAL foreign_key_checks = 0, foreign_key_checks = 1")
		case 24:
			stmt = fmt.Sprintf("CREATE TRIGGER %s BEFORE INSERT ON %s FOR EACH ROW SET NEW.x = 1", pick("tr", "TR", "tr2"), tbl())
		case 25:
			stmt = fmt.Sprintf("DROP TRIGGER %s%s", maybe("IF EXISTS "), pick("tr", "TR", "tr2"))
		}
		b.WriteString(stmt)
		b.WriteString(";\n")
	}
	return b.String()
}
