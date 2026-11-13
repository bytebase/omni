package review

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"

	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

// The WALK_THROUGH oracle for OceanBase's MySQL mode
// (docs/specs/2026-10-27-walk-through-oceanbase.md §7): each case runs a
// change on OceanBase one statement at a time and holds the walk to what the
// server answers (walktest.Judge). The tests replay their tapes
// (testdata/tapes) and start no server. A recording asks the pinned one:
//
//	OMNI_RECORD=1 go test ./oceanbase/review/ -run '^TestWalkThrough(Against|Sweep)OceanBase$' -timeout 30m
//
// A case takes three connections, each a stream of its tape: one sets the
// database up and closes, so its temporary tables and session settings end
// before the change; one syncs the target and then lists the server's objects
// around every statement; and the change runs on the third, as Bytebase runs
// it. What they send follows from the case and the server's answers alone,
// never from what the walk reports, so a change to the walk replays.

// oceanBaseWalkBase is walkBase with the relations only OceanBase has: an
// external table, which reads files under secure_file_priv (set where the
// server starts: testcontainer.OceanBase), and a materialized view, whose
// storage table the sync lists as a table.
const oceanBaseWalkBase = walkBase + `
CREATE EXTERNAL TABLE ext (a int) LOCATION = '/tmp' FORMAT = (TYPE = 'CSV');
CREATE MATERIALIZED VIEW mv AS SELECT id FROM r;
`

// oceanBaseSystem lists the databases a MySQL-mode tenant starts with.
const oceanBaseSystem = "'mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test'"

// oceanBaseWalkCases are the cases of what OceanBase does its own way.
var oceanBaseWalkCases = []walkCase{
	// The index options change nothing the walk follows, and a qualified name is held to
	// its table's database.
	{name: "index options and qualified index names", change: "CREATE INDEX n1 ON t (a) PARALLEL 2 STORING (b);\nCREATE INDEX n1 ON t (a) NOPARALLEL;\nDROP INDEX n1 ON t;\nCREATE INDEX shop.n2 ON t (a);\nCREATE INDEX n2 ON t (a);\nCREATE INDEX shop.i ON t (a);\nCREATE INDEX nosuch.n3 ON t (a);\nCREATE DATABASE q;\nCREATE INDEX q.n4 ON t (a);\nCREATE TABLE q.x (a int);\nCREATE INDEX q.n5 ON q.x (a);\nCREATE INDEX n5 ON q.x (a);", complete: true},
	// A vector index is an index the walk follows. Its one-per-column limit (1235, ahead
	// even of IF NOT EXISTS) is no name the ledger holds, so each index here is on its own
	// column or the last one's gone.
	{name: "vector indexes", change: "CREATE TABLE x (id int, v VECTOR(3) NOT NULL, VECTOR INDEX vi (v) WITH (distance=l2, type=hnsw, lib=vsag));\nCREATE TABLE x (a int);\nCREATE INDEX VI ON x (id);\nDROP INDEX vi ON x;\nDROP INDEX vi ON x;\nCREATE VECTOR INDEX vj ON x (v) WITH (distance=l2, type=hnsw, lib=vsag);\nALTER TABLE x DROP INDEX vj;\nALTER TABLE x ADD COLUMN v VECTOR(3);\nCREATE VECTOR INDEX vk ON nosuch (v) WITH (distance=l2, type=hnsw, lib=vsag);", complete: true},
	// External tables and materialized views share the tables' namespace,
	// of no kind the walk claims (spec §3). An ALTER TABLE of one ends the
	// walk; OceanBase refuses it (1235), so the walk goes on past it.
	{name: "external tables and materialized views", change: "CREATE TABLE ext (a int);\nCREATE VIEW mv AS SELECT 1 AS a;\nCREATE VIEW ext AS SELECT 1 AS a;\nRENAME TABLE t TO mv;\nRENAME TABLE t TO ext;\nALTER TABLE t RENAME TO mv;\nCREATE OR REPLACE VIEW ext AS SELECT 1 AS a;\nDROP VIEW ext;\nDROP VIEW mv;\nTRUNCATE TABLE ext;\nTRUNCATE TABLE mv;\nINSERT INTO ext VALUES (1);\nINSERT INTO mv VALUES (1);\nCREATE TABLE x LIKE mv;\nCREATE TABLE y LIKE ext;\nCREATE INDEX k ON mv (id);\nCREATE INDEX k ON ext (a);\nDROP INDEX k ON mv;\nCREATE TRIGGER tm BEFORE INSERT ON mv FOR EACH ROW SET NEW.id = 1;\nALTER TABLE ext ADD COLUMN b int;\nALTER TABLE mv COMMENT 'x';\nRENAME TABLE ext TO ext2;\nDROP TABLE ext;\nDROP TABLE ext2;\nRENAME TABLE mv TO mv2;\nCREATE TABLE mv2 (a int);\nDROP TABLE mv2;\nCREATE TABLE mv2 (a int);"},
	{name: "materialized views the sync could not read", change: "DROP TABLE nosuch;\nCREATE TABLE mv (a int);\nALTER TABLE mv COMMENT 'x';\nALTER TABLE t DROP COLUMN zz;\nRENAME TABLE mv TO mv3;", unread: true},
	{name: "ALTER TABLE of a view", change: "ALTER TABLE nosuch COMMENT 'x';\nALTER TABLE v ADD COLUMN a int;\nALTER TABLE v ENGINE = InnoDB;\nALTER TABLE v RENAME TO v2;\nALTER TABLE v COMMENT 'x';"},
	// No invisible primary key, also beside a unique key on NOT NULL
	// columns; a hidden __pk_increment; and a secondary index the sync
	// takes for the primary key.
	{name: "the primary key", setup: "CREATE TABLE sk (a int);\nCREATE INDEX `PRIMARY` ON sk (a);\nCREATE TABLE uq (a int NOT NULL, UNIQUE KEY uk (a));", change: "ALTER TABLE uq DROP PRIMARY KEY;\nCREATE TABLE uq2 (a int NOT NULL UNIQUE);\nALTER TABLE uq2 DROP PRIMARY KEY;\nALTER TABLE nokey DROP PRIMARY KEY;\nCREATE TABLE g (a int);\nALTER TABLE g DROP PRIMARY KEY;\nALTER TABLE nokey ADD COLUMN x int AFTER __pk_increment;\nALTER TABLE nokey ALTER COLUMN __pk_increment SET DEFAULT 1;\nALTER TABLE nokey ADD COLUMN __pk_increment int;\nALTER TABLE nokey DROP COLUMN __pk_increment;\nALTER TABLE nokey ADD COLUMN y int AFTER my_row_id;\nALTER TABLE sk ADD PRIMARY KEY (a);\nALTER TABLE sk ADD PRIMARY KEY (a);\nALTER TABLE t DROP PRIMARY KEY;\nALTER TABLE t ADD COLUMN z int AFTER __pk_increment;\nALTER TABLE t ADD PRIMARY KEY (id);\nALTER TABLE t ADD PRIMARY KEY (id);"},
	// PRIMARY as an index name ends the walk, so those come last.
	{name: "indexes named PRIMARY", change: "CREATE TABLE x (a int, KEY `PRIMARY` (a));\nALTER TABLE r DROP INDEX `PRIMARY`;\nDROP INDEX `PRIMARY` ON t;\nALTER TABLE t RENAME INDEX `PRIMARY` TO k2;\nCREATE INDEX `PRIMARY` ON nokey (a);\nALTER TABLE t ADD INDEX `PRIMARY` (b);"},
	{name: "foreign keys whatever the engine", change: "ALTER TABLE t ENGINE = MyISAM;\nALTER TABLE t DROP FOREIGN KEY nosuch;\nALTER TABLE t DROP FOREIGN KEY fk1, DROP FOREIGN KEY fk1;\nALTER TABLE t DROP FOREIGN KEY fk1;\nCREATE TABLE my (id int PRIMARY KEY, a int) ENGINE = MyISAM;\nALTER TABLE my ADD CONSTRAINT mf FOREIGN KEY (a) REFERENCES nosuch (id);\nALTER TABLE my DROP FOREIGN KEY nosuch;\nCREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE z (a int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT F FOREIGN KEY (a) REFERENCES r (id));", complete: true},
	{name: "a referenced parent", change: "DROP TABLE r;\nDROP TABLE nosuch, r;\nDROP TABLE child, r;"},
	// Generated names are <table>_OBFK_<16 digits>, which a rename keeps.
	{name: "generated foreign key names", change: "CREATE TABLE g (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));\nRENAME TABLE g TO h;\nALTER TABLE h DROP FOREIGN KEY h_obfk_1;\nALTER TABLE h DROP FOREIGN KEY g_ibfk_1;\nALTER TABLE h DROP FOREIGN KEY h_ibfk_1;\nCREATE DATABASE q;\nRENAME TABLE h TO q.h;\nRENAME TABLE child TO kid;\nALTER TABLE kid DROP FOREIGN KEY kid_ibfk_1;\nCREATE TABLE child (id int PRIMARY KEY, b int, CONSTRAINT child_ibfk_1 FOREIGN KEY (b) REFERENCES r (id));", complete: true},
	// An unqualified parent is in the current database.
	{name: "an unqualified parent is in the current database", change: "CREATE DATABASE q;\nCREATE TABLE q.p (id int PRIMARY KEY);\nCREATE TABLE q.c (id int, FOREIGN KEY (id) REFERENCES r (id));\nCREATE TABLE q.c2 (id int, FOREIGN KEY (id) REFERENCES p (id));\nCREATE TABLE q.c3 (id int);\nALTER TABLE q.c3 ADD CONSTRAINT c3 FOREIGN KEY (id) REFERENCES p (id);\nUSE q;\nCREATE TABLE c4 (id int, FOREIGN KEY (id) REFERENCES p (id));\nCREATE TABLE shop.c5 (id int, FOREIGN KEY (id) REFERENCES p (id));\nCREATE TABLE shop.c6 (id int, FOREIGN KEY (id) REFERENCES r (id));", complete: true},
	{name: "names up to 128 characters", change: "CREATE DATABASE " + strings.Repeat("d", 100) + ";\nALTER TABLE t ADD COLUMN " + strings.Repeat("c", 100) + " int;\nALTER TABLE t ADD COLUMN " + strings.Repeat("c", 129) + " int;\nCREATE TABLE x (a int, CONSTRAINT " + strings.Repeat("f", 100) + " FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE y (a int, CONSTRAINT " + strings.Repeat("g", 129) + " FOREIGN KEY (a) REFERENCES r (id));\nCREATE INDEX " + strings.Repeat("i", 65) + " ON t (a);\nCREATE TABLE " + strings.Repeat("t", 65) + " (a int);\nCREATE DATABASE " + strings.Repeat("d", 129) + ";", complete: true},
	{name: "DROP DATABASE drops its temporary tables", change: "CREATE DATABASE q;\nCREATE TEMPORARY TABLE q.tmp (a int);\nDROP DATABASE q;\nCREATE DATABASE q;\nDROP TEMPORARY TABLE q.tmp;\nINSERT INTO q.tmp VALUES (1);\nDROP TEMPORARY TABLE IF EXISTS q.tmp;", complete: true},
	{name: "DROP VIEW IF EXISTS of a table", change: "DROP VIEW IF EXISTS t;\nDROP VIEW IF EXISTS nosuch, t;\nDROP VIEW IF EXISTS v, nokey;\nDROP VIEW IF EXISTS v;\nDROP VIEW IF EXISTS v;", complete: true},
	{name: "a view moves to another database", change: "CREATE DATABASE q;\nRENAME TABLE v TO q.v;\nDROP VIEW v;\nDROP VIEW q.v;", complete: true},
	{name: "a table renamed to a missing database", change: "CREATE DATABASE q;\nDROP DATABASE q;\nALTER TABLE t RENAME TO q.t;\nRENAME TABLE nokey TO q.nokey;", complete: true},
	// No loadable function: an unqualified DROP FUNCTION is the stored
	// one's, but with no database selected.
	{name: "DROP FUNCTION", change: "DROP FUNCTION nosuch;\nDROP FUNCTION f;\nDROP FUNCTION f;\nDROP FUNCTION IF EXISTS f;\nCREATE FUNCTION f2() RETURNS int DETERMINISTIC RETURN 1;\nDROP DATABASE shop;\nDROP FUNCTION IF EXISTS f2;\nDROP FUNCTION f2;\nDROP PROCEDURE IF EXISTS p;"},
	// A dropped database's events stay listed under its name, and the sync
	// lists them; their names are free.
	{name: "events of a dropped database", setup: "CREATE EVENT gh ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nDROP DATABASE shop;\nCREATE DATABASE shop;\nUSE shop;\nCREATE TABLE r (id int PRIMARY KEY);", change: "CREATE EVENT gh ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nCREATE EVENT GH ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nDROP EVENT nosuch;\nALTER EVENT nosuch COMMENT 'x';\nDROP EVENT gh;\nDROP EVENT gh;\nCREATE EVENT `shop.e3` ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nDROP EVENT e3;", staleEvent: "gh"},
	// The statement's scope keyword reaches only its first assignment, and
	// an assignment after one to DEFAULT does not take its value.
	{name: "SET foreign_key_checks", change: "SET GLOBAL foreign_key_checks = 1, foreign_key_checks = 0;\nCREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET GLOBAL foreign_key_checks = 0;\nSET sql_mode = DEFAULT, foreign_key_checks = 1;\nCREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET GLOBAL foreign_key_checks = 1;\nSET sql_mode = DEFAULT, foreign_key_checks = 0;\nCREATE TABLE w (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET foreign_key_checks = 0, sql_mode = DEFAULT;\nCREATE TABLE w2 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET foreign_key_checks = 1;\nCREATE TABLE w3 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));"},
	{name: "trigger names", change: "CREATE TRIGGER tr BEFORE UPDATE ON trg FOR EACH ROW SET NEW.a = 2;\nCREATE TRIGGER tv BEFORE INSERT ON v FOR EACH ROW SET NEW.id = 1;\nCREATE TRIGGER IF NOT EXISTS tr BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;\nCREATE TRIGGER TR BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;"},
}

// oceanBaseCodes lists the error codes OceanBase gives for a finding
// beside MySQL's (spec §4.8), for the statement it is on: two primary
// keys in one statement are told from a primary key the table has by the
// statement's text.
func oceanBaseCodes(m, stmt string) []string {
	var codes []string
	for _, c := range []struct {
		match bool
		codes []string
	}{
		{strings.HasPrefix(m, "Referenced table ") || strings.HasSuffix(m, ", not a table a foreign key can reference."), []string{"1146", "1215"}},
		{saysTwice(m, "foreign key", "declares", "adds"), []string{"1215", "1235"}},
		{saysTwice(m, "index", "drops") || m == "This statement drops the primary key twice.", []string{"1072", "1235"}},
		{strings.HasSuffix(m, " already has a primary key.") && strings.Count(strings.ToUpper(stmt), "PRIMARY KEY") > 1, []string{"1062", "1235"}},
		{strings.HasSuffix(m, ", not a table.") && !leads(stmt, "DROP TABLE", "DROP TABLES", "TRUNCATE") || strings.HasSuffix(m, ", not a table a foreign key can reference."), []string{"1235"}},
		{strings.HasPrefix(m, "Database ") && strings.HasSuffix(m, " doesn't exist."), []string{"1210"}},
		{strings.HasPrefix(m, "No database is selected for "), []string{"1210"}},
		{strings.HasPrefix(m, "Trigger ") && strings.HasSuffix(m, " already exists."), []string{"5734", "5735"}},
		{strings.HasPrefix(m, "The name ") && strings.HasSuffix(m, "-character limit."), []string{"1406"}},
	} {
		if c.match {
			codes = append(codes, c.codes...)
		}
	}
	return codes
}

// refusal reads the server's error number from an error, 0 when the error
// is not the server's answer.
func refusal(err error) uint16 {
	var me *gomysql.MySQLError
	if errors.As(err, &me) {
		return me.Number
	}
	return 0
}

// mustExec runs statements on db and fails the test at the first one the server
// refuses.
func mustExec(t testing.TB, db *sql.DB, stmts ...string) {
	t.Helper()
	for _, q := range stmts {
		if _, err := db.ExecContext(context.Background(), q); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
}

// setUp runs a case's setup, and reports whether OceanBase took it. With
// refused, it must not: the setup then ends in that error, which is all the
// case shows of OceanBase.
func setUp(t testing.TB, db *sql.DB, stmts []string, refused uint16) bool {
	t.Helper()
	for _, q := range stmts {
		_, err := db.ExecContext(context.Background(), q)
		switch code := refusal(err); {
		case err == nil:
		case refused != 0 && code == refused:
			t.Logf("OceanBase refuses the setup (%v), so the case ends here", err)
			return false
		default:
			t.Fatalf("setup %s: %v", q, err)
		}
	}
	if refused != 0 {
		t.Fatalf("OceanBase takes the setup it refused with %d: the case can run now, so take refused off it", refused)
	}
	return true
}

// emptyOceanBase leaves the pinned server holding nothing a test made, with
// the global foreign_key_checks on. It runs in a recording only, outside the
// tape (testcontainer.Reset), so a tape starts from the same server whatever
// ran before it.
//
// DROP DATABASE leaves the database's events listed for some seconds (about
// 20 on 4.3.5.6), under a database created again with its name too, where the
// next sync would read them. So each event is dropped by name first, and the
// reset waits out the ones a change's own DROP DATABASE left.
func emptyOceanBase(t testing.TB) {
	t.Helper()
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		conn, err := db.Conn(ctx)
		if err != nil {
			return err
		}
		defer conn.Close()
		exec := func(q string) error {
			if _, err := conn.ExecContext(ctx, q); err != nil {
				return fmt.Errorf("%s: %w", q, err)
			}
			return nil
		}
		query := func(q string) ([][2]string, error) {
			rows, err := conn.QueryContext(ctx, q)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", q, err)
			}
			defer rows.Close()
			var out [][2]string
			for rows.Next() {
				var r [2]string
				if err := rows.Scan(&r[0], &r[1]); err != nil {
					return nil, err
				}
				out = append(out, r)
			}
			return out, rows.Err()
		}
		quote := func(name string) string { return "`" + strings.ReplaceAll(name, "`", "``") + "`" }
		// A table of one database may reference another's, so the checks are
		// off while they drop.
		for _, q := range []string{"SET GLOBAL foreign_key_checks = 1", "SET foreign_key_checks = 0"} {
			if err := exec(q); err != nil {
				return err
			}
		}
		events, err := query("SELECT e.EVENT_SCHEMA, e.EVENT_NAME FROM information_schema.EVENTS e JOIN information_schema.SCHEMATA s ON s.SCHEMA_NAME = e.EVENT_SCHEMA")
		if err != nil {
			return err
		}
		for _, e := range events {
			for _, q := range []string{"USE " + quote(e[0]), "DROP EVENT IF EXISTS " + quote(e[1])} {
				if err := exec(q); err != nil {
					return err
				}
			}
		}
		databases, err := query("SELECT SCHEMA_NAME, '' FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN (" + oceanBaseSystem + ")")
		if err != nil {
			return err
		}
		for _, d := range databases {
			if err := exec("DROP DATABASE " + quote(d[0])); err != nil {
				return err
			}
		}
		for deadline := time.Now().Add(3 * time.Minute); ; time.Sleep(500 * time.Millisecond) {
			left, err := query("SELECT EVENT_SCHEMA, EVENT_NAME FROM information_schema.EVENTS")
			if err != nil {
				return err
			}
			if len(left) == 0 {
				return nil
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("the events of dropped databases are still listed: %v", left)
			}
		}
	})
}

// walkLock names the lock a change's session holds while it runs. The
// connection that lists the server's objects finds the session by it:
// IS_USED_LOCK answers the holder's id, which is a new one on every run, so
// no request can carry it.
const walkLock = "omni_walk"

// oceanBaseInventory lists every object a change may make in a user
// database, one to a line: its space, database, table, name, and kind,
// separated by tabs. A database's line holds its name alone.
//
// It reads what a statement changes at once, and nothing that settles later:
//   - a temporary table is in no information_schema view; oceanbase.__all_table
//     lists it (table_type 6) under its session's id;
//   - information_schema.EVENTS keeps a dropped database's events for some
//     seconds, so events are read from the jobs behind it, by the database's
//     id, which a job's name starts with.
const oceanBaseInventory = `SELECT COUNT(*), GROUP_CONCAT(o.line ORDER BY BINARY o.line SEPARATOR '\n') FROM (
SELECT CONCAT_WS('\t', 'database', SCHEMA_NAME, '', '', '') AS line, SCHEMA_NAME AS db FROM information_schema.SCHEMATA
UNION ALL SELECT CONCAT_WS('\t', 'relation', TABLE_SCHEMA, '', TABLE_NAME, TABLE_TYPE), TABLE_SCHEMA FROM information_schema.TABLES
UNION ALL SELECT CONCAT_WS('\t', 'relation', OWNER, '', MVIEW_NAME, 'MATERIALIZED VIEW'), OWNER FROM oceanbase.DBA_MVIEWS
UNION ALL SELECT CONCAT_WS('\t', 'temporary', d.database_name, '', t.table_name, ''), d.database_name FROM oceanbase.__all_table t JOIN oceanbase.__all_database d ON d.database_id = t.database_id WHERE t.table_type = 6 AND t.session_id = IS_USED_LOCK('` + walkLock + `')
UNION ALL SELECT CONCAT_WS('\t', 'column', TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ''), TABLE_SCHEMA FROM information_schema.COLUMNS
UNION ALL SELECT DISTINCT CONCAT_WS('\t', 'index', TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, ''), TABLE_SCHEMA FROM information_schema.STATISTICS
UNION ALL SELECT CONCAT_WS('\t', 'foreign key', CONSTRAINT_SCHEMA, TABLE_NAME, CONSTRAINT_NAME, ''), CONSTRAINT_SCHEMA FROM information_schema.REFERENTIAL_CONSTRAINTS
UNION ALL SELECT CONCAT_WS('\t', 'trigger', TRIGGER_SCHEMA, '', TRIGGER_NAME, ''), TRIGGER_SCHEMA FROM information_schema.TRIGGERS
UNION ALL SELECT CONCAT_WS('\t', LOWER(ROUTINE_TYPE), ROUTINE_SCHEMA, '', ROUTINE_NAME, ''), ROUTINE_SCHEMA FROM information_schema.ROUTINES
UNION ALL SELECT CONCAT_WS('\t', 'event', d.database_name, '', SUBSTRING_INDEX(j.job_name, '.', -1), ''), d.database_name FROM oceanbase.__all_tenant_scheduler_job j JOIN oceanbase.__all_database d ON d.database_id = SUBSTRING_INDEX(j.job_name, '.', 1) WHERE j.job_class = 'MYSQL_EVENT_JOB_CLASS' AND j.job > 0
) o WHERE o.db NOT IN (` + oceanBaseSystem + `)`

// oceanBaseEngine is OceanBase's side of one change (walktest.Lister).
// OceanBase commits DDL as it runs it, so a rejected statement leaves the
// database as it was only by its own atomicity, which the list taken before
// and after it shows.
type oceanBaseEngine struct {
	// session is the change's session, and catalog the connection that
	// lists the server's objects.
	session *sql.DB
	catalog *sql.DB
	target  review.Target
	// reported is the statement of the finding Report returned last, which
	// walktest passes to Codes next.
	reported string
}

func (e *oceanBaseEngine) Statements(change string) []string { return splitWalk(change) }

func (e *oceanBaseEngine) Exec(t testing.TB, stmt string) walktest.Outcome {
	t.Helper()
	if _, err := e.session.ExecContext(context.Background(), stmt); err != nil {
		var me *gomysql.MySQLError
		if !errors.As(err, &me) {
			t.Fatalf("%s: %v", stmt, err)
		}
		return walktest.Outcome{Code: strconv.Itoa(int(me.Number)), Message: me.Message}
	}
	return walktest.Outcome{}
}

func (e *oceanBaseEngine) Report(t testing.TB, change string) *walktest.Finding {
	t.Helper()
	if strings.TrimSpace(change) == "" {
		return nil
	}
	result, err := ReviewOceanBase(context.Background(), change, review.Options{Rules: []review.Rule{review.WalkThrough}}, []review.Target{e.target})
	if err != nil {
		t.Fatalf("Review: %v", err)
	}
	var out *walktest.Finding
	for _, f := range result.Findings {
		if err := review.CheckCodes(f); err != nil {
			t.Error(err)
		}
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
	e.reported = ""
	if stmts := e.Statements(change); out != nil && out.Statement < len(stmts) {
		e.reported = stmts[out.Statement]
	}
	return out
}

func (e *oceanBaseEngine) Violation(t testing.TB, change string) *walk.Violation {
	t.Helper()
	stmts, syntax := parse(change, statementRanges(change))
	if syntax != nil {
		t.Fatalf("syntax: %s", syntax.Message)
	}
	v, err := walkThrough(context.Background(), change, stmts, e.target, dialect{oceanBase: true})
	if err != nil || v == nil {
		t.Fatalf("no violation for a change Review reports on (%v):\n%s", err, change)
	}
	return v
}

// Codes lists MySQL's codes for a message and the ones OceanBase gives
// beside them.
func (e *oceanBaseEngine) Codes(message string) []string {
	var codes []string
	for _, c := range walkCodes {
		if c.match(message, e.reported) {
			codes = c.codes
			break
		}
	}
	return append(slices.Clone(codes), oceanBaseCodes(message, e.reported)...)
}

// Inventory lists the server's objects with oceanBaseInventory. OceanBase
// lists a session's temporary tables, so it needs no name from the change.
func (e *oceanBaseEngine) Inventory(t testing.TB, _ string) []walktest.Object {
	t.Helper()
	var n int
	var lines sql.NullString
	if err := e.catalog.QueryRowContext(context.Background(), oceanBaseInventory).Scan(&n, &lines); err != nil {
		t.Fatalf("list the server's objects: %v", err)
	}
	var objects []walktest.Object
	for line := range strings.SplitSeq(lines.String, "\n") {
		if line == "" {
			continue
		}
		n--
		f := strings.Split(line, "\t")
		if len(f) != 5 {
			t.Fatalf("an object's line is %q, want its space, database, table, name, and kind", line)
		}
		space, db, table, name, kind := f[0], f[1], f[2], f[3], f[4]
		switch space {
		case spaceDatabase:
			objects = append(objects, walktest.Object{Key: walk.Key{Space: spaceDatabase, Name: db}})
		case spaceRelation:
			// An external table and a materialized view are relations of
			// neither kind the walk claims.
			switch kind {
			case "BASE TABLE":
				kind = relTable
			case "VIEW":
				kind = relView
			default:
				kind = strings.ToLower(kind)
			}
			objects = append(objects, walktest.Object{Key: walk.Key{Space: spaceRelation, Parent: walk.Path(db), Name: name}, Kind: kind})
		case spaceColumn, spaceIndex:
			objects = append(objects, walktest.Object{Key: walk.Key{Space: space, Parent: walk.Path(db, table), Name: name}})
		case spaceForeignKey:
			// A foreign key has two keys: in its database, and on its table.
			objects = append(objects,
				walktest.Object{Key: walk.Key{Space: spaceForeignKey, Parent: walk.Path(db), Name: name}},
				walktest.Object{Key: walk.Key{Space: spaceTableForeignKey, Parent: walk.Path(db, table), Name: name}})
		case spaceTemporary, spaceTrigger, spaceProcedure, spaceFunction, spaceEvent:
			objects = append(objects, walktest.Object{Key: walk.Key{Space: space, Parent: walk.Path(db), Name: name}})
		default:
			t.Fatalf("an object of the unknown space %q: %q", space, line)
		}
	}
	if n != 0 {
		t.Fatalf("the server listed %d more objects than it wrote, so the list is cut short: raise group_concat_max_len", n)
	}
	return objects
}

// Same compares a key the walk claims to a stored one without regard to
// case: every name folds on OceanBase, tables and databases too
// (lower_case_table_names 1). The walk follows no name outside ASCII, and
// marks the key of one the sync lists.
func (e *oceanBaseEngine) Same(claimed, stored walk.Key) bool {
	if claimed.Space != stored.Space {
		return false
	}
	return slices.EqualFunc(append(walk.Names(claimed.Parent), claimed.Name), append(walk.Names(stored.Parent), stored.Name), func(c, s string) bool {
		return strings.EqualFold(strings.TrimPrefix(c, nonASCII), s)
	})
}

// judgeOceanBase sets shop up on an empty server, syncs it, and runs the
// change on a session of its own with shop selected, as Bytebase runs it,
// holding the walk to it (walktest.Judge). A case whose setup OceanBase
// refuses, as the case says it does, ends there. With lenient, a setup
// statement the server refuses is passed over: a sweep's base is random.
func judgeOceanBase(t *testing.T, base string, c walkCase, lenient bool) (found, rejected int) {
	t.Helper()
	ctx := context.Background()
	emptyOceanBase(t)

	setup := testcontainer.Open(t)
	mustExec(t, setup, "CREATE DATABASE shop", "USE shop")
	mustExec(t, setup, splitWalk(base)...)
	if lenient {
		for _, q := range splitWalk(c.setup) {
			if _, err := setup.ExecContext(ctx, q); err != nil && refusal(err) == 0 {
				t.Fatalf("%s: %v", q, err)
			}
		}
	} else if !setUp(t, setup, splitWalk(c.setup), c.refused) {
		return 0, 0
	}
	// The setup's temporary tables and session settings end with its session.
	if err := setup.Close(); err != nil {
		t.Fatal(err)
	}

	catalog := testcontainer.Open(t)
	mustExec(t, catalog, "SET SESSION group_concat_max_len = 16777216")
	schema := oceanBaseSync(t, catalog, c.unread)
	if c.staleEvent != "" && !slices.ContainsFunc(schema.Schemas[0].Events, func(e *metadata.EventMetadata) bool { return e.Name == c.staleEvent }) {
		t.Fatalf("the sync does not list event %s, which the setup left behind a dropped database: OceanBase no longer lists it, or this run took too long to read it", c.staleEvent)
	}

	// The session checks foreign keys, as the target says: one opened while
	// a change had the global foreign_key_checks off starts with them off.
	session := testcontainer.Open(t)
	mustExec(t, session, "USE shop", "SET foreign_key_checks = 1")
	var locked sql.NullInt64
	if err := session.QueryRowContext(ctx, "SELECT GET_LOCK('"+walkLock+"', 60)").Scan(&locked); err != nil || locked.Int64 != 1 {
		t.Fatalf("the change's session did not get the lock %s: %v, %v", walkLock, locked, err)
	}
	// A session that closes keeps its lock for some seconds (about 7 on
	// 4.3.5.6), which the next case would wait out.
	t.Cleanup(func() { _, _ = session.ExecContext(ctx, "DO RELEASE_LOCK('"+walkLock+"')") })
	e := &oceanBaseEngine{
		session: session,
		catalog: catalog,
		target: review.Target{
			Schema:              schema,
			LowerCaseTableNames: 1,
			ForeignKeyChecks:    review.SwitchOn,
		},
	}
	return walktest.Judge(t, e, c.change, c.complete)
}

func TestWalkThroughAgainstOceanBase(t *testing.T) {
	var found, rejected int
	for _, tt := range slices.Concat(oceanBaseWalkCases, walkCases) {
		t.Run(tt.name, func(t *testing.T) {
			f, r := judgeOceanBase(t, oceanBaseWalkBase, tt, false)
			found += f
			rejected += r
		})
	}
	t.Logf("found %d of the %d statements the server rejects", found, rejected)
}

// The sweep's size and seed. They are fixed because the sweep's tape holds
// these changes: to sweep others, change them and record.
const (
	oceanBaseSweepN    = 40
	oceanBaseSweepSeed = 20261028
)

// TestWalkThroughSweepOceanBase runs random changes over a small pool of
// names that collide often, case included, each on a random base. Precision
// must hold on every one; recall is logged.
func TestWalkThroughSweepOceanBase(t *testing.T) {
	changes := walktest.SweepChanges(oceanBaseSweepN, oceanBaseSweepSeed, func(rng *rand.Rand) (string, string) {
		return sweepChange(rng, 8+rng.Intn(8)), sweepChange(rng, 15+rng.Intn(15))
	})
	var found, rejected int
	for i, c := range changes {
		t.Run(fmt.Sprintf("change%d", i), func(t *testing.T) {
			f, r := judgeOceanBase(t, "", walkCase{setup: c.Base, change: c.Change}, true)
			found += f
			rejected += r
			if t.Failed() {
				t.Logf("base:\n%s\nchange:\n%s", c.Base, c.Change)
			}
		})
	}
	t.Logf("seed %d: found %d of the %d statements the server rejects", oceanBaseSweepSeed, found, rejected)
}
