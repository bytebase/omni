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
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
	"github.com/bytebase/omni/tidb/parser"
)

// walkBase is the shop database the curated cases start from: tables with
// keys and foreign keys, one of a generated name, a view, and a sequence.
// t's primary key is NONCLUSTERED, which TiDB can drop.
const walkBase = `
CREATE TABLE r (id int PRIMARY KEY);
CREATE TABLE t (id int, a int, b int, PRIMARY KEY (id) NONCLUSTERED, KEY i (a), CONSTRAINT fk1 FOREIGN KEY (b) REFERENCES r (id));
CREATE TABLE child (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));
CREATE TABLE nokey (a int);
CREATE VIEW v AS SELECT id FROM r;
CREATE SEQUENCE sq;
`

// walkCase is a change on walkBase, after setup when it has one.
type walkCase struct {
	name     string
	setup    string
	change   string
	complete bool
}

// long is a name of 65 characters, one more than TiDB's identifiers take.
var long = strings.Repeat("a", 65)

// walkCases cover every row of the MySQL spec's §5 that TiDB keeps, and
// every TiDB delta of docs/specs/2026-10-26-walk-through-mariadb-tidb.md
// §3. TiDB drops some objects of a multi-object DROP of base tables, views,
// or sequences that it then fails, and fails two DROP FOREIGN KEY of one
// name after dropping it. Judge stops judging a change at such a statement
// (walktest.PartialRejecter), so each of those forms ends its case. Two
// ADD FOREIGN KEY of one name in one ALTER leave TiDB retrying the job's
// rollback without end, so only the unit tests pin that one.
var walkCases = []walkCase{
	// Databases and USE.
	{name: "databases", change: "CREATE DATABASE shop;\nCREATE DATABASE q;\nCREATE DATABASE q;\nCREATE DATABASE IF NOT EXISTS q;\nDROP DATABASE IF EXISTS nosuch;\nUSE nosuch;\nCREATE TABLE q.x (a int);\nCREATE TABLE q.x (a int);\nDROP DATABASE q;\nCREATE TABLE q.y (a int);", complete: true},
	{name: "no database selected", change: "DROP DATABASE shop;\nALTER DATABASE CHARACTER SET utf8mb4;\nCREATE TABLE x (a int);\nCREATE TEMPORARY TABLE tx (a int);\nDROP TABLE IF EXISTS tx;\nCREATE VIEW vx AS SELECT 1 AS a;\nCREATE SEQUENCE sx;\nDROP SEQUENCE IF EXISTS sx;", complete: true},
	{name: "objects in a dropped database", change: "CREATE DATABASE z;\nDROP DATABASE z;\nCREATE VIEW z.w AS SELECT 1 AS a;\nCREATE TEMPORARY TABLE z.tmp (a int);\nCREATE SEQUENCE z.s;\nCREATE TABLE z.t (a int);", complete: true},
	{name: "USE moves unqualified names", change: "CREATE DATABASE q;\nUSE q;\nCREATE TABLE t (a int);\nCREATE TABLE t (a int);\nUSE shop;\nCREATE TABLE t (a int);", complete: true},
	// Tables.
	{name: "create a table that exists", change: "CREATE TABLE t (a int);\nCREATE TABLE v (a int);\nCREATE TABLE sq (a int);\nCREATE TABLE IF NOT EXISTS t (a int);\nCREATE TABLE IF NOT EXISTS v (a int);\nCREATE TABLE IF NOT EXISTS sq (a int, a int);\nCREATE TABLE IF NOT EXISTS v LIKE nosuch;\nCREATE TABLE x LIKE nosuch;\nCREATE TABLE x LIKE v;\nCREATE TABLE x LIKE sq;", complete: true},
	{name: "a new table's text", change: "CREATE TABLE x (a int, a int);\nCREATE TABLE x (a int, b int, KEY k (a), KEY k (b));\nCREATE TABLE x (a int PRIMARY KEY, b int, PRIMARY KEY (b));\nCREATE TABLE x (a int, KEY `PRIMARY` (a));", complete: true},
	{name: "a new table holds what it declares", change: "CREATE TABLE x (a int, b int UNIQUE, KEY k (a), CONSTRAINT fx FOREIGN KEY (a) REFERENCES r (id));\nALTER TABLE x ADD COLUMN a int;\nALTER TABLE x ADD INDEX k (b);\nALTER TABLE x DROP COLUMN zz;\nALTER TABLE x DROP INDEX k, DROP INDEX k;\nALTER TABLE x DROP PRIMARY KEY;\nALTER TABLE x ADD CONSTRAINT FX FOREIGN KEY (a) REFERENCES r (id);\nALTER TABLE nokey ADD CONSTRAINT fx FOREIGN KEY (a) REFERENCES r (id);\nALTER TABLE x DROP FOREIGN KEY fx;\nALTER TABLE x DROP FOREIGN KEY fx;\nCREATE TABLE y LIKE x;\nALTER TABLE y DROP INDEX k;", complete: true},
	// TiDB checks the length of a name it creates, but not in CHANGE,
	// RENAME COLUMN, or DROP TABLE.
	{name: "names too long", change: "CREATE TABLE " + long + " (a int);\nALTER TABLE t ADD COLUMN " + long + " int;\nCREATE INDEX " + long + " ON t (a);\nRENAME TABLE nokey TO " + long + ";\nALTER TABLE nokey RENAME TO " + long + ";\nCREATE SEQUENCE " + long + ";\nCREATE VIEW " + long + " AS SELECT 1 AS a;\nDROP TABLE IF EXISTS " + long + ";\nDROP TABLE " + long + ";\nALTER TABLE t CHANGE a " + long + " int;\nALTER TABLE t RENAME COLUMN b TO b" + long + ";", complete: true},
	// ALTER TABLE judges every subcommand against the table as it was.
	{name: "ALTER TABLE judges the table as it was", change: "ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (a);\nALTER TABLE t ADD COLUMN a int, DROP COLUMN a;\nALTER TABLE t ADD COLUMN z int, DROP COLUMN z;\nALTER TABLE t DROP COLUMN a, ADD COLUMN a int;\nALTER TABLE t DROP INDEX i, ADD INDEX i (b);\nALTER TABLE t RENAME COLUMN a TO b2, RENAME COLUMN b TO a;\nALTER TABLE t ADD COLUMN y1 int, ADD COLUMN y2 int AFTER y1;\nALTER TABLE t ADD COLUMN x int, MODIFY x bigint;\nALTER TABLE t RENAME COLUMN a TO c, ADD COLUMN a int;\nALTER TABLE t ADD INDEX m (a), RENAME INDEX m TO m2;\nALTER TABLE t DROP INDEX i, ADD INDEX j (b);\nALTER TABLE t DROP INDEX j, DROP COLUMN a;\nALTER TABLE t ADD COLUMN a int;", complete: true},
	{name: "ALTER TABLE subcommands", change: "ALTER TABLE nosuch ADD COLUMN a int;\nALTER TABLE v ADD COLUMN a int;\nALTER TABLE sq ADD COLUMN a int;\nALTER TABLE sq COMMENT 'x';\nALTER TABLE t ADD COLUMN A int;\nALTER TABLE t ADD COLUMN c int AFTER nosuch;\nALTER TABLE t DROP COLUMN nosuch;\nALTER TABLE t MODIFY COLUMN nosuch int;\nALTER TABLE t CHANGE nosuch c int;\nALTER TABLE t CHANGE a b int;\nALTER TABLE t ALTER COLUMN nosuch SET DEFAULT 1;\nALTER TABLE t ADD INDEX I (b);\nALTER TABLE t ADD INDEX `PRIMARY` (b);\nALTER TABLE t DROP INDEX nosuch;\nALTER TABLE t ADD PRIMARY KEY (a);\nALTER TABLE t RENAME INDEX nosuch TO j;\nALTER TABLE t ALTER INDEX nosuch INVISIBLE;", complete: true},
	// Two subcommands on one column or index fail with 8200.
	{name: "repeated subcommands", change: "ALTER TABLE t ADD COLUMN c int, ADD COLUMN c int;\nALTER TABLE t DROP COLUMN a, DROP COLUMN a;\nALTER TABLE t ADD INDEX k (a), ADD INDEX k (b);\nALTER TABLE nokey ADD PRIMARY KEY (a), ADD PRIMARY KEY (a);\nALTER TABLE t DROP INDEX i, DROP INDEX I;\nALTER TABLE t DROP PRIMARY KEY, DROP PRIMARY KEY;\nALTER TABLE t DROP PRIMARY KEY, DROP INDEX `PRIMARY`;\nALTER TABLE t CHANGE a z int, CHANGE b z int;\nALTER TABLE t ADD COLUMN z int, RENAME COLUMN a TO z;\nALTER TABLE t DROP COLUMN IF EXISTS zz, DROP COLUMN IF EXISTS zz;", complete: true},
	// RENAME TO takes no other subcommand but ALGORITHM and LOCK.
	{name: "RENAME TO beside other subcommands", change: "ALTER TABLE t COMMENT 'x', RENAME TO t2;\nALTER TABLE t RENAME TO t2, ADD COLUMN c int;\nALTER TABLE t ENGINE = InnoDB, RENAME TO t2;\nALTER TABLE t RENAME INDEX i TO i2, RENAME TO t2;\nALTER TABLE t COMMENT 'x', RENAME TO t;\nALTER TABLE t ALGORITHM = INSTANT, RENAME TO t2;\nALTER TABLE t2 LOCK = NONE, RENAME TO t3;\nALTER TABLE t3 RENAME TO t;\nALTER TABLE t ADD COLUMN a int;", complete: true},
	{name: "DROP COLUMN IF EXISTS", change: "ALTER TABLE t DROP COLUMN IF EXISTS zz;\nALTER TABLE t ADD COLUMN q int, DROP COLUMN IF EXISTS q;\nALTER TABLE t DROP COLUMN q;\nALTER TABLE t DROP COLUMN IF EXISTS a;\nALTER TABLE t DROP COLUMN a;\nALTER TABLE t DROP COLUMN q;", complete: true},
	// Foreign key names are the table's.
	{name: "foreign keys", change: "CREATE TABLE y0 (a int, FOREIGN KEY (a) REFERENCES v (id));\nCREATE TABLE y00 (a int, CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id));\nALTER TABLE nokey DROP FOREIGN KEY fk1;\nALTER TABLE nokey ADD CONSTRAINT FK1 FOREIGN KEY (a) REFERENCES r (id);\nALTER TABLE nokey ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES r (id);\nALTER TABLE t ADD CONSTRAINT FK1 FOREIGN KEY (a) REFERENCES r (id);\nALTER TABLE t DROP FOREIGN KEY nosuch;\nALTER TABLE child DROP FOREIGN KEY fk_1, ADD CONSTRAINT fk_1 FOREIGN KEY (b) REFERENCES r (id);\nALTER TABLE t DROP FOREIGN KEY fk1;\nALTER TABLE t DROP FOREIGN KEY fk1;\nCREATE TABLE x (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nCREATE TABLE ys (a int, FOREIGN KEY (a) REFERENCES sq (id));\nSET foreign_key_checks = 0;\nCREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET foreign_key_checks = 1;\nCREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nCREATE TABLE w (id int PRIMARY KEY, p int, FOREIGN KEY (p) REFERENCES w (id));", complete: true},
	{name: "an unqualified parent is in the child's database", change: "CREATE DATABASE q;\nCREATE TABLE q.parent (id int PRIMARY KEY);\nCREATE TABLE q.child (id int, FOREIGN KEY (id) REFERENCES parent (id));\nALTER TABLE q.child ADD CONSTRAINT c2 FOREIGN KEY (id) REFERENCES parent (id);\nCREATE TABLE q.orphan (id int, FOREIGN KEY (id) REFERENCES nosuch (id));", complete: true},
	// TiDB ignores ENGINE: every table keeps its foreign keys.
	{name: "ENGINE is ignored", change: "CREATE TABLE my (id int PRIMARY KEY, a int, FOREIGN KEY (a) REFERENCES nosuch (id)) ENGINE=MyISAM;\nCREATE TABLE my (id int PRIMARY KEY, a int) ENGINE=MyISAM;\nALTER TABLE my DROP FOREIGN KEY nosuch;\nALTER TABLE my ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES nosuch (id);\nALTER TABLE my ENGINE=MyISAM;\nALTER TABLE my ADD CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id);\nALTER TABLE my ADD CONSTRAINT F FOREIGN KEY (a) REFERENCES r (id);", complete: true},
	{name: "a temporary table is no parent", change: "CREATE TEMPORARY TABLE tp (id int PRIMARY KEY);\nCREATE TABLE y1 (a int, FOREIGN KEY (a) REFERENCES tp (id));\nCREATE TEMPORARY TABLE r (id int PRIMARY KEY);\nCREATE TABLE y2 (a int, FOREIGN KEY (a) REFERENCES r (id));", complete: true},
	{name: "foreign keys named twice", change: "CREATE TABLE x (a int, b int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT f FOREIGN KEY (b) REFERENCES r (id));\nCREATE TABLE x (a int, b int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT F FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE x (a int, b int);", complete: true},
	// Generated names are fk_<n>, and RENAME TABLE keeps them.
	{name: "generated foreign key names", change: "CREATE TABLE g (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));\nRENAME TABLE g TO h;\nALTER TABLE h DROP FOREIGN KEY fk_1;\nALTER TABLE h DROP FOREIGN KEY fk_1;\nALTER TABLE nokey ADD CONSTRAINT fk_1 FOREIGN KEY (a) REFERENCES r (id);\nRENAME TABLE child TO kid;\nALTER TABLE kid DROP FOREIGN KEY fk_1;\nALTER TABLE kid DROP FOREIGN KEY fk_1;\nALTER TABLE kid RENAME TO c2;\nALTER TABLE c2 ADD CONSTRAINT fk_1 FOREIGN KEY (b) REFERENCES r (id);\nALTER TABLE c2 ADD CONSTRAINT fk_1 FOREIGN KEY (b) REFERENCES r (id);", complete: true},
	// TiDB has no invisible primary key.
	{name: "a table without a primary key has none", change: "ALTER TABLE nokey DROP PRIMARY KEY;\nCREATE TABLE g (a int);\nALTER TABLE g DROP PRIMARY KEY;\nALTER TABLE g DROP INDEX `PRIMARY`;\nALTER TABLE g ADD PRIMARY KEY (a);\nALTER TABLE g ADD PRIMARY KEY (a);", complete: true},
	{name: "a view stays in its database", change: "CREATE DATABASE q;\nRENAME TABLE v TO q.v;\nDROP VIEW v;", complete: true},
	{name: "RENAME", change: "RENAME TABLE nosuch TO x;\nRENAME TABLE t TO r;\nRENAME TABLE t TO tmp, r TO t, tmp TO r;\nALTER TABLE r ADD COLUMN a int;\nALTER TABLE t RENAME TO r;\nALTER TABLE t RENAME TO t2;\nALTER TABLE t ADD COLUMN zz int;\nCREATE DATABASE q;\nRENAME TABLE nokey TO q.nokey;\nCREATE TABLE q.nokey (a int);\nRENAME TABLE sq TO q.sq;\nCREATE SEQUENCE q.sq;\nRENAME TABLE t2 TO t3, nosuch TO t4;\nALTER TABLE t2 RENAME TO `q.t2`;\nCREATE TABLE `q.t2` (a int);\nRENAME TABLE r TO R;\nALTER TABLE r ADD COLUMN zz int, RENAME TO R;", complete: true},
	// DROP and TRUNCATE.
	{name: "DROP TABLE IF EXISTS leaves a view", setup: "CREATE DATABASE q;\nCREATE VIEW q.v AS SELECT 1 AS a;", change: "DROP TABLE IF EXISTS q.v;\nDROP VIEW q.v;\nDROP VIEW q.v;"},
	{name: "DROP TABLE IF EXISTS leaves a sequence", change: "DROP TABLE IF EXISTS sq;\nDROP SEQUENCE sq;\nDROP SEQUENCE sq;", complete: true},
	{name: "DROP TABLE", change: "DROP TABLE nosuch;\nDROP TABLE v;\nDROP TABLE sq;\nDROP TABLE IF EXISTS nosuch, nokey;\nDROP TABLE nokey;\nDROP VIEW t;\nDROP VIEW IF EXISTS t;\nDROP VIEW IF EXISTS sq;\nDROP VIEW nosuch;\nDROP VIEW IF EXISTS nosuch, v;\nDROP VIEW v;\nTRUNCATE TABLE nosuch;\nTRUNCATE TABLE v;\nTRUNCATE TABLE sq;", complete: true},
	{name: "sequences", change: "CREATE SEQUENCE sq;\nCREATE SEQUENCE SQ;\nCREATE SEQUENCE t;\nCREATE SEQUENCE v;\nCREATE SEQUENCE IF NOT EXISTS sq;\nCREATE VIEW sq AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW sq AS SELECT 1 AS a;\nDROP SEQUENCE t;\nDROP SEQUENCE v;\nDROP SEQUENCE nosuch;\nDROP SEQUENCE IF EXISTS nosuch, t;\nCREATE TEMPORARY TABLE sq (a int);\nDROP SEQUENCE sq;\nDROP SEQUENCE sq;\nCREATE SEQUENCE s2;\nRENAME TABLE s2 TO s3;\nDROP SEQUENCE s2;\nDROP SEQUENCE s3;\nCREATE TEMPORARY TABLE tt (a int);\nDROP SEQUENCE tt;", complete: true},
	{name: "a temporary table's definition", change: "CREATE TEMPORARY TABLE tmp (a int, a int);\nCREATE TEMPORARY TABLE tmp (a int, KEY k (a), KEY k (a));\nCREATE TEMPORARY TABLE tmp (a int PRIMARY KEY, PRIMARY KEY (a));\nCREATE TEMPORARY TABLE tmp (a int);\nCREATE TEMPORARY TABLE IF NOT EXISTS tmp (a int, a int);\nCREATE TEMPORARY TABLE tmp (a int);\nCREATE TEMPORARY TABLE TMP (a int);", complete: true},
	{name: "temporary tables outlive their database", change: "CREATE DATABASE q;\nCREATE TEMPORARY TABLE q.tmp (a int);\nDROP DATABASE q;\nCREATE DATABASE q;\nDROP TEMPORARY TABLE q.tmp;\nDROP TEMPORARY TABLE q.tmp;", complete: true},
	{name: "temporary tables", change: "CREATE TEMPORARY TABLE t (z int);\nALTER TABLE nokey ADD COLUMN z int;\nINSERT INTO t (z) VALUES (1);\nDROP TABLE t;\nINSERT INTO t (id) VALUES (1);\nDROP TEMPORARY TABLE t;\nDROP TEMPORARY TABLE IF EXISTS t;\nCREATE TEMPORARY TABLE tt (a int);\nCREATE TEMPORARY TABLE tt (a int);\nTRUNCATE TABLE tt;\nINSERT INTO tt VALUES (1);", complete: true},
	// TiDB refuses ALTER, CREATE INDEX, and RENAME of a local temporary
	// table (8200), which ends the walk.
	{name: "a local temporary table ends the walk", change: "CREATE TEMPORARY TABLE tt (a int);\nALTER TABLE tt ADD COLUMN b int;\nCREATE INDEX k ON tt (a);\nRENAME TABLE tt TO tt2;\nCREATE TABLE t (a int);"},
	// Indexes and views.
	{name: "CREATE and DROP INDEX", change: "CREATE INDEX k ON nosuch (a);\nCREATE INDEX i ON t (b);\nCREATE INDEX I ON t (b);\nCREATE INDEX `PRIMARY` ON t (b);\nCREATE INDEX IF NOT EXISTS i ON t (b);\nDROP INDEX nosuch ON t;\nDROP INDEX i ON nosuch;\nDROP INDEX i ON t;\nDROP INDEX i ON t;", complete: true},
	// CREATE INDEX of an expression succeeds on a view or a sequence, which
	// ends the walk.
	{name: "an index on a view or a sequence", change: "CREATE INDEX k ON v ((1 + 1));\nCREATE INDEX k ON sq ((1 + 1));\nDROP INDEX k ON sq;\nCREATE TABLE t (a int);"},
	{name: "views", change: "CREATE VIEW t AS SELECT 1 AS a;\nCREATE VIEW v AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW t AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW v AS SELECT 1 AS a;\nDROP VIEW IF EXISTS t;\nDROP VIEW t;", complete: true},
	// DML.
	{name: "DML targets", change: "INSERT INTO nosuch VALUES (1);\nUPDATE nosuch SET a = 1;\nDELETE FROM nosuch;\nREPLACE INTO nosuch VALUES (1);\nINSERT INTO nokey VALUES (1);\nUPDATE nokey SET a = 2;\nDELETE nokey FROM nokey JOIN r ON r.id = nokey.a;", complete: true},
	// Name case: tables and databases always fold (lower_case_table_names
	// 2), as do the other spaces.
	{name: "name case", change: "CREATE TABLE T (a int);\nDROP TABLE NOKEY;\nCREATE DATABASE Q;\nCREATE DATABASE q;\nCREATE TABLE Q.x (a int);\nCREATE TABLE q.X (a int);\nALTER TABLE R ADD COLUMN ID int;\nCREATE VIEW V AS SELECT 1 AS a;\nCREATE SEQUENCE SQ;\nDROP SEQUENCE Sq;\nDROP SEQUENCE sQ;\nALTER TABLE T DROP FOREIGN KEY FK1;\nALTER TABLE t DROP FOREIGN KEY fk1;", complete: true},
	// What the walk leaves Unknown by design: databases other than the
	// target, the tables a query reads, generated index names, the hidden
	// row id, index columns, what TiDB does not support, and procedures,
	// which end the walk.
	// TiDB drops what it finds of these before it fails: each is judged by
	// its error code, and ends its case.
	{name: "DROP TABLE drops the tables before a missing name", change: "DROP TABLE t, nosuch;", complete: true},
	{name: "DROP TABLE drops the tables before a missing name, listed last", change: "DROP TABLE nosuch, nokey;", complete: true},
	{name: "DROP TABLE drops the tables before a view", change: "DROP TABLE nokey, v;", complete: true},
	{name: "DROP TABLE drops the tables before a sequence", change: "DROP TABLE nokey, sq;", complete: true},
	{name: "DROP VIEW drops the views before a missing name", change: "DROP VIEW v, nosuch;", complete: true},
	{name: "DROP VIEW IF EXISTS drops the views before a table", change: "DROP VIEW IF EXISTS v, t;", complete: true},
	{name: "DROP SEQUENCE drops the sequences before a missing name", change: "DROP SEQUENCE sq, nosuch;", complete: true},
	{name: "DROP SEQUENCE drops the sequences before a table", change: "DROP SEQUENCE sq, t;", complete: true},
	{name: "DROP FOREIGN KEY twice drops it", change: "ALTER TABLE t DROP FOREIGN KEY fk1, DROP FOREIGN KEY fk1;", complete: true},
	{name: "unknown by design", change: "DROP DATABASE nosuch;\nRENAME TABLE t TO nodb.t;\nCREATE VIEW nodb.w AS SELECT 1 AS a;\nCREATE TABLE x AS SELECT 1 AS a;\nCREATE TABLE x (a int, b int UNIQUE);\nALTER TABLE x DROP INDEX zz;\nALTER TABLE nokey ADD COLUMN _tidb_rowid int;\nALTER TABLE nokey ADD COLUMN y int AFTER _tidb_rowid;\nALTER TABLE nokey DROP COLUMN _tidb_rowid;\nALTER TABLE t ADD COLUMN q int, ADD INDEX kq (q);\nALTER TABLE t DROP COLUMN IF EXISTS a, DROP COLUMN a;\nALTER TABLE t ALTER COLUMN a SET DEFAULT 1, MODIFY a int;\nALTER TABLE r DROP PRIMARY KEY;\nALTER TABLE t DROP CONSTRAINT nosuch;\nINSERT INTO v VALUES (1);\nCREATE TABLE c3 (id int, a int, FOREIGN KEY (id) REFERENCES r (id), CONSTRAINT fk_1 FOREIGN KEY (a) REFERENCES r (id));\nCREATE PROCEDURE p() SELECT 1;\nCREATE TABLE t (a int);"},
}

// startWalkTiDB returns the package's TiDB, started once and shared with
// the other oracle tests, in its shipped setting: lower_case_table_names
// is 2 and cannot change.
func startWalkTiDB(t *testing.T) *sql.DB {
	t.Helper()
	db := startTiDB(t)
	var lctn int
	if err := db.QueryRowContext(context.Background(), "SELECT @@lower_case_table_names").Scan(&lctn); err != nil {
		t.Fatal(err)
	}
	if lctn != 2 {
		t.Fatalf("lower_case_table_names = %d, want 2", lctn)
	}
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

func TestWalkThroughAgainstTiDB(t *testing.T) {
	db := startWalkTiDB(t)
	var found, rejected int
	for _, tt := range walkCases {
		t.Run(tt.name, func(t *testing.T) {
			resetWalk(t, db, walkBase+tt.setup)
			f, r := judgeWalk(t, db, tt.change, tt.complete)
			found += f
			rejected += r
		})
	}
	t.Logf("found %d of the %d statements the server rejects", found, rejected)
}

// systemDatabases are the databases TiDB keeps for itself.
const systemDatabases = "('mysql', 'information_schema', 'performance_schema', 'metrics_schema', 'sys')"

// resetWalk drops every user database and recreates shop from base.
func resetWalk(t *testing.T, db *sql.DB, base string) {
	t.Helper()
	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE LOWER(SCHEMA_NAME) NOT IN "+systemDatabases)
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
	stmts := []string{"SET foreign_key_checks = 0"}
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
func judgeWalk(t *testing.T, db *sql.DB, change string, complete bool) (found, rejected int) {
	t.Helper()
	ctx := context.Background()
	conn := session(t, db)
	if _, err := conn.ExecContext(ctx, "USE shop"); err != nil {
		t.Fatal(err)
	}
	e := &tidbEngine{
		conn: conn,
		db:   db,
		target: review.Target{
			Schema:              walkSync(t, db),
			LowerCaseTableNames: 2,
			ForeignKeyChecks:    review.SwitchOn,
		},
	}
	return walktest.Judge(t, e, change, complete)
}

// tidbEngine is the oracle's TiDB side of one change. TiDB commits DDL as
// it runs, so a rejected statement leaves the database as it was only
// through its own atomicity, which Digest checks.
type tidbEngine struct {
	conn   *sql.Conn
	db     *sql.DB
	target review.Target
}

// RejectsPartially marks TiDB as a server that may change the database in
// a statement it rejects: multi-object drops drop what they find before
// they fail, and a repeated DROP FOREIGN KEY drops the key first
// (walktest.PartialRejecter).
func (*tidbEngine) RejectsPartially() {}

func (e *tidbEngine) Statements(change string) []string { return splitWalk(change) }

func (e *tidbEngine) Exec(t testing.TB, stmt string) walktest.Outcome {
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

func (e *tidbEngine) Report(t testing.TB, change string) *walktest.Finding {
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

func (e *tidbEngine) Violation(t testing.TB, change string) *walk.Violation {
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

// walkCodes maps a finding's message to the TiDB errors it stands for.
// 8200 is TiDB's "Unsupported operate same column" and its kin.
var walkCodes = []struct {
	match func(string) bool
	codes []string
}{
	{func(m string) bool { return strings.HasPrefix(m, "Temporary table ") }, []string{"1051"}},
	{func(m string) bool { return strings.HasPrefix(m, "Foreign key target ") }, []string{"1824", "3734", "3780"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Foreign key ") && (strings.Contains(m, " specified twice ") || strings.Contains(m, " added twice ") || strings.Contains(m, " already exists in table "))
	}, []string{"1826"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "View ") && strings.HasSuffix(m, " cannot move to another database")
	}, []string{"1450"}},
	{func(m string) bool { return strings.HasPrefix(m, "Foreign key ") }, []string{"1091"}},
	{func(m string) bool { return strings.HasPrefix(m, "Table ") && strings.HasSuffix(m, " already exists") }, []string{"1050"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Table ") && strings.HasSuffix(m, " already has a primary key")
	}, []string{"1068", "8200"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Table ") && strings.HasSuffix(m, " has no primary key")
	}, []string{"1091"}},
	{func(m string) bool { return strings.HasPrefix(m, "Table ") && strings.HasSuffix(m, " does not exist") }, []string{"1146", "1051"}},
	{func(m string) bool { return strings.HasPrefix(m, "View ") }, []string{"1051"}},
	{func(m string) bool { return strings.HasPrefix(m, "Sequence ") }, []string{"4139"}},
	{func(m string) bool {
		return strings.HasSuffix(m, " is not a base table") || strings.HasSuffix(m, " is not a view") || strings.HasSuffix(m, " is not a sequence")
	}, []string{"1347"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Database ") && strings.HasSuffix(m, " already exists")
	}, []string{"1007"}},
	{func(m string) bool { return strings.HasPrefix(m, "Database ") }, []string{"1049", "1008", "1025"}},
	{func(m string) bool { return m == noDatabaseMessage }, []string{"1046"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Column ") && (strings.Contains(m, " already exists in ") || strings.Contains(m, " twice "))
	}, []string{"1060", "1091", "8200"}},
	{func(m string) bool { return strings.HasPrefix(m, "Column ") }, []string{"1054", "1091"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Index ") && (strings.Contains(m, " already exists in ") || strings.Contains(m, " specified twice") || strings.Contains(m, " added twice"))
	}, []string{"1061", "8200"}},
	{func(m string) bool { return strings.HasPrefix(m, "Index name PRIMARY") }, []string{"1280"}},
	{func(m string) bool { return strings.HasPrefix(m, "Index ") }, []string{"1091", "1176", "8200"}},
	{func(m string) bool { return strings.HasPrefix(m, "Multiple primary keys") }, []string{"1068", "8200"}},
	{func(m string) bool { return strings.HasPrefix(m, "RENAME TO cannot be combined") }, []string{"8200"}},
	{func(m string) bool { return strings.HasPrefix(m, "Identifier name ") }, []string{"1059"}},
}

func (e *tidbEngine) Codes(message string) []string {
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
func (e *tidbEngine) Holds(t testing.TB, v *walk.Violation, message string) (bool, string) {
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
// Database and table names compare folded, as lower_case_table_names 2
// does.
func (e *tidbEngine) exists(t testing.TB, k walk.Key) (bool, string) {
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
		return count("SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE "+eq("SCHEMA_NAME", nameFolds), name), ""
	case spaceTemporary:
		var tbl, ddl string
		err := e.conn.QueryRowContext(ctx, "SHOW CREATE TABLE `"+unmark(names[0])+"`.`"+name+"`").Scan(&tbl, &ddl)
		return err == nil && strings.HasPrefix(ddl, "CREATE TEMPORARY TABLE"), ""
	case spaceRelation:
		var kind string
		err := e.db.QueryRowContext(ctx, "SELECT TABLE_TYPE FROM information_schema.TABLES WHERE "+eq("TABLE_SCHEMA", true)+" AND "+eq("TABLE_NAME", nameFolds), unmark(names[0]), name).Scan(&kind)
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
		table, schema, column := "information_schema.COLUMNS", "TABLE_SCHEMA", "COLUMN_NAME"
		switch k.Space {
		case spaceIndex:
			table, column = "information_schema.STATISTICS", "INDEX_NAME"
		case spaceTableForeignKey:
			table, schema, column = "information_schema.REFERENTIAL_CONSTRAINTS", "CONSTRAINT_SCHEMA", "CONSTRAINT_NAME"
		}
		return count("SELECT COUNT(*) FROM "+table+" WHERE "+eq(schema, true)+" AND "+eq("TABLE_NAME", true)+" AND "+eq(column, nameFolds), unmark(names[0]), unmark(names[1]), name), ""
	}
	t.Fatalf("unknown space %s", k.Space)
	return false, ""
}

// Digest reads every user database's objects, so a rejected statement
// that changed one fails the case.
func (e *tidbEngine) Digest(t testing.TB) string {
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
	dbs := read("SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE LOWER(SCHEMA_NAME) NOT IN " + systemDatabases + " ORDER BY 1")
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
		"SELECT CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA IN " + in + " ORDER BY 1, 2, 3",
	} {
		read(q, args...)
	}
	return b.String()
}

// walkSync reads shop into the metadata Bytebase's TiDB sync produces for
// what WalkThrough reads: base tables with their columns, indexes, and
// foreign keys, views, and sequences, which it reads from
// information_schema.SEQUENCES and keeps apart from the tables.
//
// bytebase: backend/plugin/db/tidb/sync.go
func walkSync(t *testing.T, db *sql.DB) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	out := syncSchema(t, conn)
	s := out.Schemas[0]
	rows, err := conn.QueryContext(ctx, "SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = 'shop'")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	sequences := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		sequences[name] = true
		s.Sequences = append(s.Sequences, &metadata.SequenceMetadata{Name: name})
	}
	// syncSchema takes every relation but a view for a table.
	var tables []*metadata.TableMetadata
	for _, tbl := range s.Tables {
		if !sequences[tbl.Name] {
			tables = append(tables, tbl)
		}
	}
	s.Tables = tables
	return out
}

// TestWalkThroughSweepTiDB runs random changes over a small pool of names
// that collide often, case included. Precision must hold on every one;
// recall is logged. WALK_SWEEP_N sets the number of changes,
// WALK_SWEEP_SEED the seed; the server is the package's TiDB.
//
// A change ends its judging at a statement TiDB rejects after changing the
// database (walktest.PartialRejecter). The templates leave out two ADD
// FOREIGN KEY of one name in one ALTER, which leave TiDB retrying the
// job's rollback without end.
func TestWalkThroughSweepTiDB(t *testing.T) {
	n, seed, _ := walktest.SweepParams(t, 40, 20261026, "pingcap/tidb:v8.5.5")
	db := startWalkTiDB(t)
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
			f, r := judgeWalk(t, db, change, false)
			found += f
			rejected += r
			if t.Failed() {
				t.Logf("base:\n%s\nchange:\n%s", base, change)
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
	name := func() string { return pick("a", "b", "c", "A", "v", "s") }
	tbl := func() string { return db() + name() }
	col := func() string { return pick("x", "y", "z", "X") }
	idx := func() string { return pick("i", "j", "x", "I") }
	fk := func() string { return pick("f", "F", "fk_1", "g") }
	var b strings.Builder
	for range n {
		var stmt string
		switch rng.Intn(28) {
		case 0, 1:
			stmt = fmt.Sprintf("CREATE %sTABLE %s%s (id int PRIMARY KEY, x int%s, y int)", maybe("TEMPORARY "), maybe("IF NOT EXISTS "), tbl(), maybe(" UNIQUE"))
		case 2:
			stmt = fmt.Sprintf("CREATE TABLE %s (x int, y int, KEY %s (x)%s)", tbl(), idx(), maybe(", CONSTRAINT "+fk()+" FOREIGN KEY (y) REFERENCES a (id)"))
		case 3:
			stmt = fmt.Sprintf("CREATE TABLE %s%s LIKE %s", maybe("IF NOT EXISTS "), tbl(), tbl())
		case 4:
			stmt = fmt.Sprintf("CREATE SEQUENCE %s%s", maybe("IF NOT EXISTS "), tbl())
		case 5:
			stmt = fmt.Sprintf("DROP %sTABLE %s%s%s", maybe("TEMPORARY "), maybe("IF EXISTS "), tbl(), maybe(", "+tbl()))
		case 6:
			stmt = fmt.Sprintf("DROP SEQUENCE %s%s%s", maybe("IF EXISTS "), tbl(), maybe(", "+tbl()))
		case 7:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s int%s", tbl(), col(), maybe(" AFTER "+col()))
		case 8:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s%s%s", tbl(), maybe("IF EXISTS "), col(), maybe(", ADD COLUMN "+col()+" int"))
		case 9:
			stmt = fmt.Sprintf("ALTER TABLE %s CHANGE %s %s int%s", tbl(), col(), col(), maybe(", "+pick("CHANGE "+col()+" "+col()+" int", "RENAME COLUMN "+col()+" TO "+col(), "MODIFY "+col()+" bigint")))
		case 10:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD %sINDEX %s (%s)", tbl(), maybe("UNIQUE "), idx(), col())
		case 11:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP INDEX %s%s", tbl(), idx(), maybe(", ADD INDEX "+idx()+" (id)"))
		case 12:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (x) REFERENCES %s (id)", tbl(), fk(), tbl())
		case 13:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s%s", tbl(), fk(), maybe(", "+pick("ADD COLUMN "+col()+" int", "DROP FOREIGN KEY "+fk())))
		case 14:
			stmt = fmt.Sprintf("ALTER TABLE %s %s", tbl(), pick("ADD PRIMARY KEY (x)", "DROP PRIMARY KEY", "RENAME INDEX i TO j", "RENAME TO "+name(), "ENGINE = InnoDB", "COMMENT 'x', RENAME TO "+name()))
		case 15:
			stmt = fmt.Sprintf("RENAME TABLE %s TO %s%s", tbl(), tbl(), maybe(", "+tbl()+" TO "+tbl()))
		case 16:
			stmt = fmt.Sprintf("CREATE %sINDEX %s%s ON %s (%s)", maybe("UNIQUE "), maybe("IF NOT EXISTS "), idx(), tbl(), col())
		case 17:
			stmt = fmt.Sprintf("DROP INDEX %s ON %s", idx(), tbl())
		case 18:
			stmt = fmt.Sprintf("CREATE %sVIEW %s AS SELECT 1 AS x", maybe("OR REPLACE "), tbl())
		case 19:
			stmt = fmt.Sprintf("DROP VIEW %s%s%s", maybe("IF EXISTS "), tbl(), maybe(", "+tbl()))
		case 20:
			stmt = pick("CREATE DATABASE q", "CREATE DATABASE IF NOT EXISTS q", "DROP DATABASE IF EXISTS q", "USE q", "USE shop")
		case 21:
			stmt = fmt.Sprintf("TRUNCATE TABLE %s", tbl())
		case 22:
			stmt = fmt.Sprintf("INSERT INTO %s (x) VALUES (1)", tbl())
		case 23:
			stmt = pick("SET foreign_key_checks = 0", "SET foreign_key_checks = 1")
		case 24:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s int, DROP COLUMN %s", tbl(), col(), col())
		case 25:
			stmt = fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s%s", tbl(), col(), col(), maybe(", ADD COLUMN "+col()+" int"))
		case 26:
			stmt = fmt.Sprintf("CREATE TABLE %s AS SELECT 1 AS x", tbl())
		case 27:
			stmt = fmt.Sprintf("UPDATE %s SET x = 1", tbl())
		}
		b.WriteString(stmt)
		b.WriteString(";\n")
	}
	return b.String()
}
