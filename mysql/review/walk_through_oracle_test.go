package review

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

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
	{name: "no database selected", change: "DROP DATABASE shop;\nALTER DATABASE CHARACTER SET utf8mb4;\nCREATE TABLE x (a int);\nCREATE PROCEDURE p2() SELECT 1;\nDROP PROCEDURE IF EXISTS p;\nDROP FUNCTION shop.f;\nDROP TRIGGER IF EXISTS tr;\nCREATE EVENT e2 ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nALTER EVENT e COMMENT 'x';", complete: true},
	{name: "objects in a dropped database", change: "CREATE DATABASE z;\nDROP DATABASE z;\nCREATE PROCEDURE z.p() SELECT 1;\nCREATE FUNCTION z.f() RETURNS int DETERMINISTIC RETURN 1;\nCREATE VIEW z.w AS SELECT 1 AS a;\nCREATE TEMPORARY TABLE z.tmp (a int);", complete: true},
	{name: "USE moves unqualified names", change: "CREATE DATABASE q;\nUSE q;\nCREATE TABLE t (a int);\nCREATE TABLE t (a int);\nUSE shop;\nCREATE TABLE t (a int);", complete: true},
	// Tables.
	{name: "create a table that exists", change: "CREATE TABLE t (a int);\nCREATE TABLE v (a int);\nCREATE TABLE IF NOT EXISTS t (a int);\nCREATE TABLE IF NOT EXISTS v (a int);\nCREATE TABLE IF NOT EXISTS v LIKE nosuch;\nCREATE TABLE x LIKE nosuch;", complete: true},
	{name: "a new table's text", change: "CREATE TABLE x (a int, a int);\nCREATE TABLE x (a int, b int, KEY k (a), KEY k (b));\nCREATE TABLE x (a int PRIMARY KEY, b int, PRIMARY KEY (b));\nCREATE TABLE x (a int, KEY `PRIMARY` (a));", complete: true},
	{name: "a new table holds what it declares", change: "CREATE TABLE x (a int, b int UNIQUE, KEY k (a), CONSTRAINT fx FOREIGN KEY (a) REFERENCES r (id));\nALTER TABLE x ADD COLUMN a int;\nALTER TABLE x ADD INDEX k (b);\nALTER TABLE x DROP COLUMN zz;\nALTER TABLE x DROP INDEX k, DROP INDEX k;\nALTER TABLE nokey ADD CONSTRAINT fx FOREIGN KEY (a) REFERENCES r (id);\nCREATE TABLE y LIKE x;\nALTER TABLE y DROP INDEX k;", complete: true},
	{name: "names too long", change: "CREATE TABLE " + strings.Repeat("a", 65) + " (a int);\nALTER TABLE t ADD COLUMN " + strings.Repeat("c", 65) + " int;\nDROP TABLE " + strings.Repeat("d", 65) + ";", complete: true},
	// ALTER TABLE.
	{name: "ALTER TABLE drops first", change: "ALTER TABLE t ADD COLUMN a int, DROP COLUMN a;\nALTER TABLE t DROP INDEX i, ADD INDEX i (b);\nALTER TABLE t DROP COLUMN a;\nALTER TABLE t DROP COLUMN a;\nALTER TABLE t RENAME COLUMN b TO c, ADD COLUMN b int;\nALTER TABLE t DROP COLUMN zz, ADD COLUMN y int;", complete: true},
	{name: "ALTER TABLE subcommands", change: "ALTER TABLE nosuch ADD COLUMN a int;\nALTER TABLE v ADD COLUMN a int;\nALTER TABLE t ADD COLUMN A int;\nALTER TABLE t ADD COLUMN c int AFTER nosuch;\nALTER TABLE t DROP COLUMN nosuch;\nALTER TABLE t MODIFY COLUMN nosuch int;\nALTER TABLE t CHANGE nosuch c int;\nALTER TABLE t CHANGE a b int;\nALTER TABLE t ALTER COLUMN nosuch SET DEFAULT 1;\nALTER TABLE t ADD INDEX I (b);\nALTER TABLE t ADD INDEX `PRIMARY` (b);\nALTER TABLE t DROP INDEX nosuch;\nALTER TABLE t ADD PRIMARY KEY (a);\nALTER TABLE t RENAME INDEX nosuch TO j;\nALTER TABLE t RENAME INDEX i TO PRIMARY2;\nALTER TABLE t ALTER INDEX nosuch INVISIBLE;", complete: true},
	{name: "repeated subcommands", change: "ALTER TABLE t ADD COLUMN c int, ADD COLUMN c int;\nALTER TABLE t DROP COLUMN a, DROP COLUMN a;\nALTER TABLE t ADD INDEX k (a), ADD INDEX k (b);\nALTER TABLE nokey ADD PRIMARY KEY (a), ADD PRIMARY KEY (a);\nALTER TABLE t DROP INDEX i, DROP INDEX i;\nALTER TABLE t DROP FOREIGN KEY fk1, DROP FOREIGN KEY fk1;\nALTER TABLE t DROP PRIMARY KEY, DROP PRIMARY KEY;\nALTER TABLE t DROP PRIMARY KEY, DROP INDEX `PRIMARY`;", complete: true},
	{name: "foreign keys", change: "CREATE TABLE y0 (a int, FOREIGN KEY (a) REFERENCES v (id));\nCREATE TABLE y00 (a int, CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id));\nALTER TABLE nokey DROP FOREIGN KEY fk1;\nALTER TABLE nokey ADD CONSTRAINT FK1 FOREIGN KEY (a) REFERENCES r (id);\nALTER TABLE t DROP FOREIGN KEY nosuch;\nALTER TABLE t DROP FOREIGN KEY fk1;\nALTER TABLE t DROP FOREIGN KEY fk1;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET foreign_key_checks = 0;\nCREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET foreign_key_checks = 1;\nCREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nCREATE TABLE w (id int PRIMARY KEY, p int, FOREIGN KEY (p) REFERENCES w (id));", complete: true},
	{name: "an unqualified parent is in the child's database", change: "CREATE DATABASE q;\nCREATE TABLE q.parent (id int PRIMARY KEY);\nCREATE TABLE q.child (id int, FOREIGN KEY (id) REFERENCES parent (id));\nALTER TABLE q.child ADD CONSTRAINT c2 FOREIGN KEY (id) REFERENCES parent (id);\nCREATE TABLE q.orphan (id int, FOREIGN KEY (id) REFERENCES nosuch (id));", complete: true},
	// An engine other than InnoDB parses foreign keys and ignores them.
	{name: "an engine that ignores foreign keys", setup: "CREATE TABLE my (id int PRIMARY KEY, a int) ENGINE=MyISAM;", change: "ALTER TABLE my DROP FOREIGN KEY nosuch;\nALTER TABLE my ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES nosuch (id);\nCREATE TABLE m2 (a int, CONSTRAINT mf FOREIGN KEY (a) REFERENCES nosuch (id), CONSTRAINT mf FOREIGN KEY (a) REFERENCES r (id)) ENGINE=MyISAM;\nCREATE TABLE ok (a int, CONSTRAINT mf FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE m3 LIKE my;\nALTER TABLE m3 ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id);\nCREATE TABLE inn (id int PRIMARY KEY, a int);\nALTER TABLE inn ENGINE=MyISAM;\nALTER TABLE inn ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES nosuch (id);\nALTER TABLE inn DROP FOREIGN KEY whatever;\nALTER TABLE inn ENGINE=InnoDB;\nALTER TABLE inn DROP FOREIGN KEY whatever;"},
	{name: "column names a statement adds", setup: "CREATE TABLE cols (a int, b int, c int);", change: "ALTER TABLE cols RENAME COLUMN a TO z, RENAME COLUMN b TO z;\nALTER TABLE cols CHANGE a z int, CHANGE b z int;\nALTER TABLE cols ADD COLUMN z int, RENAME COLUMN a TO z;\nALTER TABLE cols CHANGE a z int, ADD COLUMN Z int;\nALTER TABLE cols RENAME COLUMN a TO b, RENAME COLUMN b TO a;\nALTER TABLE cols ADD COLUMN b2 int AFTER b;\nALTER TABLE cols RENAME COLUMN a TO b, RENAME COLUMN b TO y;\nALTER TABLE cols ADD COLUMN y2 int AFTER y, ADD COLUMN b3 int AFTER b;\nALTER TABLE cols ADD COLUMN a int;\nALTER TABLE cols ADD KEY i (a), ADD KEY j (b);\nALTER TABLE cols RENAME INDEX i TO j, RENAME INDEX j TO i;\nALTER TABLE cols DROP INDEX i, DROP INDEX j;\nALTER TABLE cols ADD KEY i (a), ADD KEY j (b);\nALTER TABLE cols RENAME INDEX i TO j, DROP INDEX j;\nALTER TABLE cols DROP INDEX j;\nALTER TABLE cols DROP INDEX j;\nALTER TABLE cols CHANGE a A int, ADD COLUMN d int;\nALTER TABLE cols MODIFY a int PRIMARY KEY, MODIFY b int PRIMARY KEY;\nALTER TABLE cols CHANGE a a int PRIMARY KEY, ADD PRIMARY KEY (b);\nALTER TABLE cols MODIFY a int NOT NULL PRIMARY KEY, ADD COLUMN e int PRIMARY KEY;", complete: true},
	{name: "a temporary table is no parent", change: "CREATE TEMPORARY TABLE tp (id int PRIMARY KEY);\nCREATE TABLE y1 (a int, FOREIGN KEY (a) REFERENCES tp (id));\nCREATE TEMPORARY TABLE r (id int PRIMARY KEY);\nCREATE TABLE y2 (a int, FOREIGN KEY (a) REFERENCES r (id));", complete: true},
	{name: "foreign keys named twice", change: "CREATE TABLE x (a int, b int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT f FOREIGN KEY (b) REFERENCES r (id));\nCREATE TABLE x (a int, b int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT F FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE x (a int, b int);\nALTER TABLE x ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), ADD CONSTRAINT F FOREIGN KEY (b) REFERENCES r (id);", complete: true},
	// RENAME TABLE renames the foreign keys named <table>_ibfk_...
	{name: "renamed generated foreign key names", change: "CREATE TABLE g (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));\nRENAME TABLE g TO h;\nALTER TABLE nokey ADD CONSTRAINT g_ibfk_1 FOREIGN KEY (a) REFERENCES r (id);\nALTER TABLE h DROP FOREIGN KEY h_ibfk_1;\nALTER TABLE h DROP FOREIGN KEY h_ibfk_1;\nRENAME TABLE child TO kid;\nALTER TABLE kid DROP FOREIGN KEY kid_ibfk_1;\nALTER TABLE kid DROP FOREIGN KEY kid_ibfk_1;\nALTER TABLE kid RENAME TO c2;\nCREATE TABLE child (id int PRIMARY KEY, b int, CONSTRAINT child_ibfk_1 FOREIGN KEY (b) REFERENCES r (id));", complete: true},
	{name: "a table with no primary key may get an invisible one", change: "CREATE TABLE g (a int);\nALTER TABLE g DROP PRIMARY KEY;\nALTER TABLE g ADD PRIMARY KEY (a);"},
	{name: "the sync may not list an invisible primary key", setup: "SET GLOBAL show_gipk_in_create_table_and_information_schema = OFF;\nSET SESSION sql_generate_invisible_primary_key = ON;\nCREATE TABLE gk (a int);", change: "ALTER TABLE gk ADD COLUMN x int AFTER my_row_id;\nALTER TABLE gk ADD COLUMN my_row_id int;\nALTER TABLE gk ADD PRIMARY KEY (a);"},
	{name: "a view stays in its database", change: "CREATE DATABASE q;\nRENAME TABLE v TO q.v;\nDROP VIEW v;", complete: true},
	{name: "RENAME", change: "RENAME TABLE nosuch TO x;\nRENAME TABLE t TO r;\nRENAME TABLE t TO tmp, r TO t, tmp TO r;\nALTER TABLE r ADD COLUMN a int;\nALTER TABLE t RENAME TO r;\nALTER TABLE t RENAME TO t2;\nALTER TABLE t ADD COLUMN zz int;\nCREATE DATABASE q;\nRENAME TABLE nokey TO q.nokey;\nCREATE TABLE q.nokey (a int);", complete: true},
	// DROP and TRUNCATE.
	{name: "DROP TABLE IF EXISTS leaves a view", setup: "CREATE DATABASE q;\nCREATE VIEW q.v AS SELECT 1 AS a;", change: "DROP TABLE IF EXISTS q.v;\nDROP VIEW q.v;\nDROP VIEW q.v;"},
	{name: "DROP TABLE", change: "DROP TABLE t, nosuch;\nDROP TABLE v;\nDROP TABLE IF EXISTS nosuch, nokey;\nDROP TABLE nokey;\nDROP VIEW t;\nDROP VIEW nosuch;\nDROP VIEW IF EXISTS nosuch, v;\nDROP VIEW v;\nTRUNCATE TABLE nosuch;\nTRUNCATE TABLE v;", complete: true},
	{name: "a temporary table's definition", change: "CREATE TEMPORARY TABLE tmp (a int, a int);\nCREATE TEMPORARY TABLE tmp (a int, KEY k (a), KEY k (a));\nCREATE TEMPORARY TABLE tmp (a int PRIMARY KEY, PRIMARY KEY (a));\nCREATE TEMPORARY TABLE tmp (a int);\nCREATE TEMPORARY TABLE IF NOT EXISTS tmp (a int, a int);\nCREATE TEMPORARY TABLE tmp (a int);", complete: true},
	{name: "temporary tables outlive their database", change: "CREATE DATABASE q;\nCREATE TEMPORARY TABLE q.tmp (a int);\nDROP DATABASE q;\nCREATE DATABASE q;\nDROP TEMPORARY TABLE q.tmp;\nDROP TEMPORARY TABLE q.tmp;", complete: true},
	{name: "temporary tables", change: "CREATE TEMPORARY TABLE t (z int);\nALTER TABLE nokey ADD COLUMN z int;\nINSERT INTO t (z) VALUES (1);\nDROP TABLE t;\nINSERT INTO t (id) VALUES (1);\nDROP TEMPORARY TABLE t;\nDROP TEMPORARY TABLE IF EXISTS t;\nCREATE TEMPORARY TABLE tt (a int);\nCREATE TEMPORARY TABLE tt (a int);\nTRUNCATE TABLE tt;", complete: true},
	// Indexes, views, triggers, routines, events.
	{name: "CREATE and DROP INDEX", change: "CREATE INDEX k ON nosuch (a);\nCREATE INDEX k ON v (id);\nCREATE INDEX i ON t (b);\nCREATE INDEX `PRIMARY` ON t (b);\nDROP INDEX nosuch ON t;\nDROP INDEX i ON nosuch;\nDROP INDEX i ON t;\nDROP INDEX i ON t;", complete: true},
	{name: "views", change: "CREATE VIEW t AS SELECT 1 AS a;\nCREATE VIEW v AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW t AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW v AS SELECT 1 AS a;\nALTER VIEW nosuch AS SELECT 1 AS a;\nALTER VIEW t AS SELECT 1 AS a;\nDROP VIEW IF EXISTS t;\nDROP VIEW t;", complete: true},
	{name: "triggers", change: "CREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;\nCREATE TRIGGER tr2 BEFORE INSERT ON nosuch FOR EACH ROW SET NEW.a = 1;\nDROP TRIGGER nosuch;\nDROP TRIGGER tr;\nDROP TRIGGER tr;\nCREATE TRIGGER TR BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;", complete: true},
	{name: "routines and events", change: "CREATE PROCEDURE p() SELECT 1;\nCREATE PROCEDURE P() SELECT 1;\nCREATE FUNCTION p() RETURNS int DETERMINISTIC RETURN 1;\nCREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;\nDROP PROCEDURE nosuch;\nDROP FUNCTION shop.nosuch;\nALTER PROCEDURE nosuch COMMENT 'x';\nCREATE EVENT E ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nDROP EVENT nosuch;\nALTER EVENT nosuch COMMENT 'x';\nALTER EVENT e RENAME TO e2;\nDROP EVENT e;\nCREATE EVENT `shop.e2` ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nDROP EVENT `shop.e2`;\nDROP EVENT `shop.e2`;", complete: true},
	// An unqualified DROP FUNCTION drops a loadable function of the name
	// first, which the sync does not list.
	{name: "loadable functions", change: "CREATE FUNCTION service_get_read_locks RETURNS INT SONAME 'locking_service.so';\nDROP FUNCTION service_get_read_locks;\nDROP FUNCTION shop.service_get_read_locks;\nCREATE FUNCTION service_get_read_locks RETURNS INT SONAME 'locking_service.so';\nCREATE FUNCTION service_get_read_locks() RETURNS int DETERMINISTIC RETURN 1;\nDROP FUNCTION service_get_read_locks;\nDROP FUNCTION shop.service_get_read_locks;\nDROP FUNCTION shop.service_get_read_locks;", complete: true},
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
	// are; resetWalk sets the global back.
	{name: "SET foreign_key_checks in each assignment's scope", change: "SET GLOBAL foreign_key_checks = 1, @@session.foreign_key_checks = 0;\nCREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET GLOBAL foreign_key_checks = 1, @@foreign_key_checks = 1;\nCREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET @@global.foreign_key_checks = 1, foreign_key_checks = 0;\nCREATE TABLE z2 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET GLOBAL foreign_key_checks = 1;\nCREATE TABLE z3 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET foreign_key_checks = 1;\nSET GLOBAL foreign_key_checks = 0;\nCREATE TABLE z4 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET GLOBAL foreign_key_checks = 1, foreign_key_checks = 0;\nCREATE TABLE w (a int, FOREIGN KEY (a) REFERENCES nosuch (id));", complete: true},
	// What the walk leaves Unknown by design (spec §6): databases other than
	// the target, the tables a query reads, generated index names, foreign
	// key names held per database, a primary key GIPK may have added, a key
	// two parts of one ALTER TABLE both change, and a function an
	// unqualified DROP FUNCTION may find loaded.
	{name: "unknown by design", change: "DROP FUNCTION nosuch;\nDROP DATABASE nosuch;\nRENAME TABLE t TO nodb.t;\nCREATE VIEW nodb.w AS SELECT 1 AS a;\nCREATE TABLE x AS SELECT * FROM nosuch;\nCREATE TABLE x (a int, b int UNIQUE);\nALTER TABLE x DROP INDEX zz;\nCREATE TABLE y (a int, CONSTRAINT fy FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE z LIKE y;\nALTER TABLE z DROP FOREIGN KEY fy;\nALTER TABLE nokey DROP PRIMARY KEY;\nALTER TABLE t ADD COLUMN z int, DROP COLUMN z;\nALTER TABLE t ADD INDEX j (b), DROP INDEX j;\nALTER TABLE t ADD COLUMN x int, RENAME COLUMN x TO y;"},
	{name: "a trigger of another database may create a temporary table", setup: "DROP TRIGGER tr;\nCREATE DATABASE q;\nCREATE TABLE q.o (a int);\nCREATE TRIGGER q.ot BEFORE INSERT ON q.o FOR EACH ROW CREATE TEMPORARY TABLE IF NOT EXISTS shop.made (a int);", change: "INSERT INTO q.o VALUES (1);\nDROP TEMPORARY TABLE made;\nDROP TEMPORARY TABLE made;"},
	{name: "a query's function may create a temporary table", setup: "CREATE FUNCTION mk() RETURNS int DETERMINISTIC BEGIN CREATE TEMPORARY TABLE IF NOT EXISTS made (a int); RETURN 1; END;", change: "CREATE TABLE c AS SELECT mk() AS x;\nDROP TEMPORARY TABLE made;\nDROP TEMPORARY TABLE made;"},
	{name: "a trigger may create a temporary table", change: "CREATE TRIGGER tt BEFORE INSERT ON nokey FOR EACH ROW CREATE TEMPORARY TABLE IF NOT EXISTS shadow (a int);\nINSERT INTO nokey VALUES (1);\nINSERT INTO shadow VALUES (1);\nDROP TABLE shadow;"},
}

// walkServer is one MySQL configuration the oracle runs on.
type walkServer struct {
	image string
	lctn  int
}

var walkServers = struct {
	mu         sync.Mutex
	dbs        map[walkServer]*sql.DB
	containers []testcontainers.Container
}{dbs: make(map[walkServer]*sql.DB)}

// startWalkMySQL returns a server of the configuration, started once.
func startWalkMySQL(t *testing.T, s walkServer) *sql.DB {
	t.Helper()
	walkServers.mu.Lock()
	defer walkServers.mu.Unlock()
	if db, ok := walkServers.dbs[s]; ok {
		return db
	}
	ctx := context.Background()
	opts := []testcontainers.ContainerCustomizer{tcmysql.WithDatabase("shop"), tcmysql.WithUsername("root"), tcmysql.WithPassword("test")}
	if s.lctn != 0 {
		opts = append(opts, testcontainers.WithCmdArgs("--lower-case-table-names="+strconv.Itoa(s.lctn)))
	}
	c, err := tcmysql.Run(ctx, s.image, opts...)
	if err != nil {
		t.Fatalf("mysql %v unavailable: %v", s, err)
	}
	walkServers.containers = append(walkServers.containers, c)
	dsn, err := c.ConnectionString(ctx, "multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
	walkServers.dbs[s] = db
	return db
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

func TestWalkThroughAgainstMySQL(t *testing.T) {
	runWalkCases(t, walkServer{image: "mysql:8.0"})
}

// TestWalkThroughAgainstMySQLFolded runs the cases with lower_case_table_names=1,
// where table and database names fold.
func TestWalkThroughAgainstMySQLFolded(t *testing.T) {
	runWalkCases(t, walkServer{image: "mysql:8.0", lctn: 1})
}

func TestWalkThroughAgainstMySQL84(t *testing.T) {
	runWalkCases(t, walkServer{image: "mysql:8.4"})
}

func runWalkCases(t *testing.T, s walkServer) {
	db := startWalkMySQL(t, s)
	var found, rejected int
	for _, tt := range walkCases {
		t.Run(tt.name, func(t *testing.T) {
			resetWalk(t, db, walkBase+tt.setup)
			f, r := judgeWalk(t, db, s, tt.change, tt.complete)
			found += f
			rejected += r
		})
	}
	t.Logf("found %d of the %d statements the server rejects", found, rejected)
}

// resetWalk drops every user database and recreates shop from base.
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
	var udfs []string
	rows, err = db.QueryContext(ctx, "SELECT name FROM mysql.func")
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			t.Fatal(err)
		}
		udfs = append(udfs, n)
	}
	rows.Close()
	conn := session(t, db)
	// A change may set the global foreign_key_checks; the server goes back
	// to its default.
	stmts := []string{"SET GLOBAL foreign_key_checks = 1", "SET foreign_key_checks = 0", "SET GLOBAL show_gipk_in_create_table_and_information_schema = ON"}
	for _, n := range names {
		stmts = append(stmts, "DROP DATABASE `"+strings.ReplaceAll(n, "`", "``")+"`")
	}
	// A loadable function is the server's, not a database's; with every
	// database gone, DROP FUNCTION drops it.
	for _, n := range udfs {
		stmts = append(stmts, "DROP FUNCTION `"+strings.ReplaceAll(n, "`", "``")+"`")
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
func judgeWalk(t *testing.T, db *sql.DB, s walkServer, change string, complete bool) (found, rejected int) {
	t.Helper()
	ctx := context.Background()
	conn := session(t, db)
	// The session checks foreign keys, as the target says. One the pool
	// opened while a change had the global foreign_key_checks off starts
	// with them off.
	for _, q := range []string{"USE shop", "SET foreign_key_checks = 1"} {
		if _, err := conn.ExecContext(ctx, q); err != nil {
			t.Fatal(err)
		}
	}
	e := &mysqlEngine{
		conn: conn,
		db:   db,
		fold: s.lctn != 0,
		target: review.Target{
			Schema:              walkSync(t, db),
			LowerCaseTableNames: s.lctn,
			ForeignKeyChecks:    review.SwitchOn,
		},
	}
	return walktest.Judge(t, e, change, complete)
}

// mysqlEngine is the oracle's MySQL side of one change. MySQL commits DDL
// as it runs, so a rejected statement leaves the database as it was only
// through its own atomicity, which Digest checks.
type mysqlEngine struct {
	conn   *sql.Conn
	db     *sql.DB
	fold   bool
	target review.Target
}

func (e *mysqlEngine) Statements(change string) []string { return splitWalk(change) }

func (e *mysqlEngine) Exec(t testing.TB, stmt string) walktest.Outcome {
	t.Helper()
	if _, err := e.conn.ExecContext(context.Background(), stmt); err != nil {
		var me *gomysql.MySQLError
		if !errors.As(err, &me) {
			t.Fatalf("%s: %v", stmt, err)
		}
		return walktest.Outcome{Code: strconv.Itoa(int(me.Number)), Message: me.Message}
	}
	return walktest.Outcome{}
}

func (e *mysqlEngine) Report(t testing.TB, change string) *walktest.Finding {
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

func (e *mysqlEngine) Violation(t testing.TB, change string) *walk.Violation {
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

// walkCodes maps a finding's message to the MySQL errors it stands for
// (5.7's codes included).
var walkCodes = []struct {
	match func(string) bool
	codes []string
}{
	{func(m string) bool { return strings.HasPrefix(m, "Temporary table ") }, []string{"1051"}},
	{func(m string) bool { return strings.HasPrefix(m, "Foreign key target ") }, []string{"1824", "1215"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Foreign key ") && (strings.Contains(m, " specified twice ") || strings.Contains(m, " added twice "))
	}, []string{"1826", "1061"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "View ") && strings.HasSuffix(m, " cannot move to another database")
	}, []string{"1450"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Foreign key ") && strings.HasSuffix(m, " already exists in the database")
	}, []string{"1826", "1022"}},
	{func(m string) bool { return strings.HasPrefix(m, "Foreign key ") }, []string{"1091"}},
	{func(m string) bool { return strings.HasPrefix(m, "Table ") && strings.HasSuffix(m, " already exists") }, []string{"1050"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Table ") && strings.HasSuffix(m, " already has a primary key")
	}, []string{"1068"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Table ") && strings.HasSuffix(m, " has no primary key")
	}, []string{"1091"}},
	{func(m string) bool { return strings.HasPrefix(m, "Table ") && strings.HasSuffix(m, " does not exist") }, []string{"1146", "1051", "1017"}},
	{func(m string) bool { return strings.HasPrefix(m, "View ") }, []string{"1051"}},
	{func(m string) bool {
		return strings.HasSuffix(m, " is not a base table") || strings.HasSuffix(m, " is not a view")
	}, []string{"1347"}},
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
	{func(m string) bool { return strings.HasPrefix(m, "Index ") }, []string{"1091", "1176"}},
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
	{func(m string) bool { return strings.HasPrefix(m, "Identifier name ") }, []string{"1059"}},
}

func (e *mysqlEngine) Codes(message string) []string {
	for _, c := range walkCodes {
		if c.match(message) {
			return c.codes
		}
	}
	return nil
}

// Holds checks each fact the violation claims against information_schema.
// A claim about the text holds by itself; a temporary table is read with
// SHOW CREATE TABLE, which information_schema does not list.
func (e *mysqlEngine) Holds(t testing.TB, v *walk.Violation, message string) (bool, string) {
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
func (e *mysqlEngine) exists(t testing.TB, k walk.Key) (bool, string) {
	t.Helper()
	ctx := context.Background()
	names := walk.Names(k.Parent)
	// eq compares a stored name to a key as the space folds it.
	eq := func(column string, folds bool) string {
		if folds {
			return "LOWER(" + column + ") = BINARY ?"
		}
		return "BINARY " + column + " = BINARY ?"
	}
	name := strings.TrimPrefix(k.Name, nonASCII)
	nameFolds := !strings.HasPrefix(k.Name, nonASCII)
	tableFolds := e.fold
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
		return count("SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE "+eq("SCHEMA_NAME", tableFolds && nameFolds), name), ""
	case spaceTemporary:
		var tbl, ddl string
		err := e.conn.QueryRowContext(ctx, "SHOW CREATE TABLE `"+unmark(names[0])+"`.`"+name+"`").Scan(&tbl, &ddl)
		return err == nil && strings.HasPrefix(ddl, "CREATE TEMPORARY TABLE"), ""
	case spaceRelation:
		var kind string
		err := e.db.QueryRowContext(ctx, "SELECT TABLE_TYPE FROM information_schema.TABLES WHERE "+eq("TABLE_SCHEMA", tableFolds)+" AND "+eq("TABLE_NAME", tableFolds && nameFolds), unmark(names[0]), name).Scan(&kind)
		if errors.Is(err, sql.ErrNoRows) {
			return false, ""
		}
		if err != nil {
			t.Fatal(err)
		}
		if kind == "VIEW" {
			return true, relView
		}
		return true, relTable
	case spaceColumn, spaceIndex, spaceTableForeignKey:
		table, column := "information_schema.COLUMNS", "COLUMN_NAME"
		switch k.Space {
		case spaceIndex:
			table, column = "information_schema.STATISTICS", "INDEX_NAME"
		case spaceTableForeignKey:
			table, column = "information_schema.REFERENTIAL_CONSTRAINTS", "CONSTRAINT_NAME"
			return count("SELECT COUNT(*) FROM "+table+" WHERE "+eq("CONSTRAINT_SCHEMA", tableFolds)+" AND "+eq("TABLE_NAME", tableFolds)+" AND "+eq(column, nameFolds), unmark(names[0]), unmark(names[1]), name), ""
		}
		return count("SELECT COUNT(*) FROM "+table+" WHERE "+eq("TABLE_SCHEMA", tableFolds)+" AND "+eq("TABLE_NAME", tableFolds)+" AND "+eq(column, nameFolds), unmark(names[0]), unmark(names[1]), name), ""
	}
	where := map[string]string{
		spaceForeignKey: "information_schema.REFERENTIAL_CONSTRAINTS WHERE " + eq("CONSTRAINT_SCHEMA", tableFolds) + " AND " + eq("CONSTRAINT_NAME", nameFolds),
		spaceTrigger:    "information_schema.TRIGGERS WHERE " + eq("TRIGGER_SCHEMA", tableFolds) + " AND " + eq("TRIGGER_NAME", nameFolds),
		spaceProcedure:  "information_schema.ROUTINES WHERE ROUTINE_TYPE = 'PROCEDURE' AND " + eq("ROUTINE_SCHEMA", tableFolds) + " AND " + eq("ROUTINE_NAME", nameFolds),
		spaceFunction:   "information_schema.ROUTINES WHERE ROUTINE_TYPE = 'FUNCTION' AND " + eq("ROUTINE_SCHEMA", tableFolds) + " AND " + eq("ROUTINE_NAME", nameFolds),
		spaceEvent:      "information_schema.EVENTS WHERE " + eq("EVENT_SCHEMA", tableFolds) + " AND " + eq("EVENT_NAME", nameFolds),
	}[k.Space]
	if where == "" {
		t.Fatalf("unknown space %s", k.Space)
	}
	return count("SELECT COUNT(*) FROM "+where, unmark(names[0]), name), ""
}

// Digest reads every user database's objects, so a rejected statement
// that changed one fails the case. Each query names the databases, which
// keeps information_schema from reading the system ones.
func (e *mysqlEngine) Digest(t testing.TB) string {
	t.Helper()
	var b strings.Builder
	read := func(q string, args ...any) []string {
		rows, err := e.db.QueryContext(context.Background(), q, args...)
		if err != nil {
			t.Fatalf("%s: %v", q, err)
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
	dbs := read("SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys') ORDER BY 1")
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
		read(q, args...)
	}
	return b.String()
}

// walkSync reads shop into the metadata Bytebase's MySQL sync produces for
// what WalkThrough reads: base tables with their engines, columns, indexes,
// foreign keys, and triggers, views, routines, and events.
//
// bytebase: backend/plugin/db/mysql/sync.go
func walkSync(t *testing.T, db *sql.DB) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	out := syncSchema(t, conn)
	s := out.Schemas[0]
	query := func(q string, add func(kind, name string)) {
		rows, err := conn.QueryContext(context.Background(), q)
		if err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		defer rows.Close()
		for rows.Next() {
			var kind, name string
			if err := rows.Scan(&kind, &name); err != nil {
				t.Fatal(err)
			}
			add(kind, name)
		}
	}
	query("SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop'", func(kind, name string) {
		if kind == "PROCEDURE" {
			s.Procedures = append(s.Procedures, &metadata.ProcedureMetadata{Name: name})
		} else {
			s.Functions = append(s.Functions, &metadata.FunctionMetadata{Name: name})
		}
	})
	query("SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop'", func(_, name string) {
		s.Events = append(s.Events, &metadata.EventMetadata{Name: name})
	})
	query("SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE'", func(name, engine string) {
		for _, t := range s.Tables {
			if t.Name == name {
				t.Engine = engine
			}
		}
	})
	return out
}

// TestWalkThroughSweepMySQL runs random changes over a small pool of names
// that collide often, case included. Precision must hold on every one;
// recall is logged. WALK_SWEEP_N sets the number of changes,
// WALK_SWEEP_SEED the seed, and WALK_SWEEP_IMAGE the MySQL image.
func TestWalkThroughSweepMySQL(t *testing.T) {
	n, seed, image := walktest.SweepParams(t, 40, 20261023, "mysql:8.0")
	// The two settings run at once, each on its own server.
	for _, lctn := range []int{0, 1} {
		t.Run(fmt.Sprintf("lctn%d", lctn), func(t *testing.T) {
			t.Parallel()
			s := walkServer{image: image, lctn: lctn}
			db := startWalkMySQL(t, s)
			rng := rand.New(rand.NewSource(seed))
			var found, rejected int
			for i := range n {
				base := walkSweepChange(rng, 8+rng.Intn(8))
				change := walkSweepChange(rng, 15+rng.Intn(15))
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
						t.Logf("base:\n%s\nchange:\n%s", base, change)
					}
				})
			}
			t.Logf("%s lctn=%d seed %d: found %d of the %d statements the server rejects", image, lctn, seed, found, rejected)
		})
	}
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
	tbl := func() string { return db() + pick("a", "b", "c", "A", "v", "a_ibfk_1") }
	col := func() string { return pick("x", "y", "z", "X") }
	idx := func() string { return pick("i", "j", "x", "I") }
	fk := func() string { return pick("f", "F", "a_ibfk_1", "g") }
	var b strings.Builder
	for range n {
		var stmt string
		switch rng.Intn(26) {
		case 0, 1:
			stmt = fmt.Sprintf("CREATE %sTABLE %s%s (id int PRIMARY KEY, x int%s, y int)", maybe("TEMPORARY "), maybe("IF NOT EXISTS "), tbl(), maybe(" UNIQUE"))
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
			stmt = fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s%s", tbl(), col(), maybe(", ADD COLUMN "+col()+" int"))
		case 9:
			stmt = fmt.Sprintf("ALTER TABLE %s CHANGE %s %s int", tbl(), col(), col())
		case 10:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD %sINDEX %s (%s)", tbl(), maybe("UNIQUE "), idx(), col())
		case 11:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP INDEX %s%s", tbl(), idx(), maybe(", ADD INDEX "+idx()+" (x)"))
		case 12:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (x) REFERENCES %s (id)", tbl(), fk(), tbl())
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
