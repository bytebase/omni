package review

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk/walktest"
)

// The WALK_THROUGH oracle for OceanBase's MySQL mode
// (docs/specs/2026-10-27-walk-through-oceanbase.md §7). The server needs
// about 8 GB, which CI does not have: the tests read OCEANBASE_WALK_DSN,
// which scripts/oceanbase-oracle.sh sets, and skip without it.
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/review/ -run 'OceanBase' -count=1 -timeout 3h

// oceanBaseWalkBase is walkBase with the relations only OceanBase has: an
// external table, which reads files under secure_file_priv (the script
// sets it), and a materialized view, whose storage table the sync lists
// as a table.
const oceanBaseWalkBase = walkBase + `
CREATE EXTERNAL TABLE ext (a int) LOCATION = '/tmp' FORMAT = (TYPE = 'CSV');
CREATE MATERIALIZED VIEW mv AS SELECT id FROM r;
`

// oceanBaseSystem lists the databases a MySQL-mode tenant starts with.
const oceanBaseSystem = "'mysql', 'information_schema', 'performance_schema', 'sys', 'oceanbase', 'SYS', 'LBACSYS', 'ORAAUDITOR', '__public', '__recyclebin', 'test'"

// oceanBaseCase is a change on oceanBaseWalkBase. unread syncs the target
// as a user without SELECT on the oceanbase database does: with no
// materialized views.
type oceanBaseCase struct {
	walkCase
	unread bool
}

// oceanBaseSkipped lists the MySQL cases whose setup OceanBase refuses.
var oceanBaseSkipped = map[string]string{
	"the sync may not list an invisible primary key":                           "no sql_generate_invisible_primary_key",
	"a trigger of another database may create a temporary table":               "a trigger cannot create a temporary table (1422)",
	"a query's function may create a temporary table":                          "a function cannot create a temporary table (1422)",
	"a function a view or another database holds may create a temporary table": "a function cannot create a temporary table (1422)",
}

// oceanBaseIncomplete lists the complete MySQL cases where OceanBase
// refuses a statement the walk leaves unclaimed by design.
var oceanBaseIncomplete = map[string]string{
	"ALTER TABLE drops first":                 "pairs on one name and an index dropped and added (1060, 4179)",
	"ALTER TABLE subcommands":                 "ALTER TABLE of a view, and ADD INDEX `PRIMARY`, end the walk",
	"repeated subcommands":                    "DROP INDEX `PRIMARY` ends the walk",
	"column names a statement adds":           "pairs on one name (1060, 1061, 1176)",
	"a temporary table's definition":          "IF NOT EXISTS of an existing one checks its text (1060)",
	"CREATE and DROP INDEX":                   "CREATE INDEX `PRIMARY` ends the walk",
	"routines and events":                     "a listed event, an event name with a period",
	"loadable functions":                      "none can be created (1235)",
	"a temporary table is no parent":          "a parent a temporary table shadows fails too (1215)",
	"DROP FUNCTION with no database selected": "DROP FUNCTION IF EXISTS in a dropped database fails (1049)",
	"a foreign key's name and its index's":    "a foreign key adds no index (1091), two in one ALTER TABLE (1235)",
}

var oceanBaseWalkCases = []oceanBaseCase{
	// External tables and materialized views share the tables' namespace,
	// of no kind the walk claims (spec §3). An ALTER TABLE of one ends the
	// walk; OceanBase refuses it (1235), so the walk goes on past it.
	{walkCase: walkCase{name: "external tables and materialized views", change: "CREATE TABLE ext (a int);\nCREATE VIEW mv AS SELECT 1 AS a;\nCREATE VIEW ext AS SELECT 1 AS a;\nRENAME TABLE t TO mv;\nRENAME TABLE t TO ext;\nALTER TABLE t RENAME TO mv;\nCREATE OR REPLACE VIEW ext AS SELECT 1 AS a;\nDROP VIEW ext;\nDROP VIEW mv;\nTRUNCATE TABLE ext;\nTRUNCATE TABLE mv;\nINSERT INTO ext VALUES (1);\nINSERT INTO mv VALUES (1);\nCREATE TABLE x LIKE mv;\nCREATE TABLE y LIKE ext;\nCREATE INDEX k ON mv (id);\nCREATE INDEX k ON ext (a);\nDROP INDEX k ON mv;\nCREATE TRIGGER tm BEFORE INSERT ON mv FOR EACH ROW SET NEW.id = 1;\nALTER TABLE ext ADD COLUMN b int;\nALTER TABLE mv COMMENT 'x';\nRENAME TABLE ext TO ext2;\nDROP TABLE ext;\nDROP TABLE ext2;\nRENAME TABLE mv TO mv2;\nCREATE TABLE mv2 (a int);\nDROP TABLE mv2;\nCREATE TABLE mv2 (a int);"}},
	{walkCase: walkCase{name: "materialized views the sync could not read", change: "DROP TABLE nosuch;\nCREATE TABLE mv (a int);\nALTER TABLE mv COMMENT 'x';\nALTER TABLE t DROP COLUMN zz;\nRENAME TABLE mv TO mv3;"}, unread: true},
	{walkCase: walkCase{name: "ALTER TABLE of a view", change: "ALTER TABLE nosuch COMMENT 'x';\nALTER TABLE v ADD COLUMN a int;\nALTER TABLE v ENGINE = InnoDB;\nALTER TABLE v RENAME TO v2;\nALTER TABLE v COMMENT 'x';"}},
	// No invisible primary key, also beside a unique key on NOT NULL
	// columns; a hidden __pk_increment; and a secondary index the sync
	// takes for the primary key.
	{walkCase: walkCase{name: "the primary key", setup: "CREATE TABLE sk (a int);\nCREATE INDEX `PRIMARY` ON sk (a);\nCREATE TABLE uq (a int NOT NULL, UNIQUE KEY uk (a));", change: "ALTER TABLE uq DROP PRIMARY KEY;\nCREATE TABLE uq2 (a int NOT NULL UNIQUE);\nALTER TABLE uq2 DROP PRIMARY KEY;\nALTER TABLE nokey DROP PRIMARY KEY;\nCREATE TABLE g (a int);\nALTER TABLE g DROP PRIMARY KEY;\nALTER TABLE nokey ADD COLUMN x int AFTER __pk_increment;\nALTER TABLE nokey ALTER COLUMN __pk_increment SET DEFAULT 1;\nALTER TABLE nokey ADD COLUMN __pk_increment int;\nALTER TABLE nokey DROP COLUMN __pk_increment;\nALTER TABLE nokey ADD COLUMN y int AFTER my_row_id;\nALTER TABLE sk ADD PRIMARY KEY (a);\nALTER TABLE sk ADD PRIMARY KEY (a);\nALTER TABLE t DROP PRIMARY KEY;\nALTER TABLE t ADD COLUMN z int AFTER __pk_increment;\nALTER TABLE t ADD PRIMARY KEY (id);\nALTER TABLE t ADD PRIMARY KEY (id);"}},
	// PRIMARY as an index name ends the walk, so those come last.
	{walkCase: walkCase{name: "indexes named PRIMARY", change: "CREATE TABLE x (a int, KEY `PRIMARY` (a));\nALTER TABLE r DROP INDEX `PRIMARY`;\nDROP INDEX `PRIMARY` ON t;\nALTER TABLE t RENAME INDEX `PRIMARY` TO k2;\nCREATE INDEX `PRIMARY` ON nokey (a);\nALTER TABLE t ADD INDEX `PRIMARY` (b);"}},
	{walkCase: walkCase{name: "foreign keys whatever the engine", change: "ALTER TABLE t ENGINE = MyISAM;\nALTER TABLE t DROP FOREIGN KEY nosuch;\nALTER TABLE t DROP FOREIGN KEY fk1, DROP FOREIGN KEY fk1;\nALTER TABLE t DROP FOREIGN KEY fk1;\nCREATE TABLE my (id int PRIMARY KEY, a int) ENGINE = MyISAM;\nALTER TABLE my ADD CONSTRAINT mf FOREIGN KEY (a) REFERENCES nosuch (id);\nALTER TABLE my DROP FOREIGN KEY nosuch;\nCREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE z (a int, CONSTRAINT f FOREIGN KEY (a) REFERENCES r (id), CONSTRAINT F FOREIGN KEY (a) REFERENCES r (id));", complete: true}},
	{walkCase: walkCase{name: "a referenced parent", change: "DROP TABLE r;\nDROP TABLE nosuch, r;\nDROP TABLE child, r;"}},
	// Generated names are <table>_OBFK_<16 digits>, which a rename keeps.
	{walkCase: walkCase{name: "generated foreign key names", change: "CREATE TABLE g (id int PRIMARY KEY, b int, FOREIGN KEY (b) REFERENCES r (id));\nRENAME TABLE g TO h;\nALTER TABLE h DROP FOREIGN KEY h_obfk_1;\nALTER TABLE h DROP FOREIGN KEY g_ibfk_1;\nALTER TABLE h DROP FOREIGN KEY h_ibfk_1;\nCREATE DATABASE q;\nRENAME TABLE h TO q.h;\nRENAME TABLE child TO kid;\nALTER TABLE kid DROP FOREIGN KEY kid_ibfk_1;\nCREATE TABLE child (id int PRIMARY KEY, b int, CONSTRAINT child_ibfk_1 FOREIGN KEY (b) REFERENCES r (id));", complete: true}},
	// An unqualified parent is in the current database.
	{walkCase: walkCase{name: "an unqualified parent is in the current database", change: "CREATE DATABASE q;\nCREATE TABLE q.p (id int PRIMARY KEY);\nCREATE TABLE q.c (id int, FOREIGN KEY (id) REFERENCES r (id));\nCREATE TABLE q.c2 (id int, FOREIGN KEY (id) REFERENCES p (id));\nCREATE TABLE q.c3 (id int);\nALTER TABLE q.c3 ADD CONSTRAINT c3 FOREIGN KEY (id) REFERENCES p (id);\nUSE q;\nCREATE TABLE c4 (id int, FOREIGN KEY (id) REFERENCES p (id));\nCREATE TABLE shop.c5 (id int, FOREIGN KEY (id) REFERENCES p (id));\nCREATE TABLE shop.c6 (id int, FOREIGN KEY (id) REFERENCES r (id));", complete: true}},
	{walkCase: walkCase{name: "names up to 128 characters", change: "CREATE DATABASE " + strings.Repeat("d", 100) + ";\nALTER TABLE t ADD COLUMN " + strings.Repeat("c", 100) + " int;\nALTER TABLE t ADD COLUMN " + strings.Repeat("c", 129) + " int;\nCREATE TABLE x (a int, CONSTRAINT " + strings.Repeat("f", 100) + " FOREIGN KEY (a) REFERENCES r (id));\nCREATE TABLE y (a int, CONSTRAINT " + strings.Repeat("g", 129) + " FOREIGN KEY (a) REFERENCES r (id));\nCREATE INDEX " + strings.Repeat("i", 65) + " ON t (a);\nCREATE TABLE " + strings.Repeat("t", 65) + " (a int);\nCREATE DATABASE " + strings.Repeat("d", 129) + ";", complete: true}},
	{walkCase: walkCase{name: "DROP DATABASE drops its temporary tables", change: "CREATE DATABASE q;\nCREATE TEMPORARY TABLE q.tmp (a int);\nDROP DATABASE q;\nCREATE DATABASE q;\nDROP TEMPORARY TABLE q.tmp;\nINSERT INTO q.tmp VALUES (1);\nDROP TEMPORARY TABLE IF EXISTS q.tmp;", complete: true}},
	{walkCase: walkCase{name: "DROP VIEW IF EXISTS of a table", change: "DROP VIEW IF EXISTS t;\nDROP VIEW IF EXISTS nosuch, t;\nDROP VIEW IF EXISTS v, nokey;\nDROP VIEW IF EXISTS v;\nDROP VIEW IF EXISTS v;", complete: true}},
	{walkCase: walkCase{name: "a view moves to another database", change: "CREATE DATABASE q;\nRENAME TABLE v TO q.v;\nDROP VIEW v;\nDROP VIEW q.v;", complete: true}},
	{walkCase: walkCase{name: "a table renamed to a missing database", change: "CREATE DATABASE q;\nDROP DATABASE q;\nALTER TABLE t RENAME TO q.t;\nRENAME TABLE nokey TO q.nokey;", complete: true}},
	// No loadable function: an unqualified DROP FUNCTION is the stored
	// one's, but with no database selected.
	{walkCase: walkCase{name: "DROP FUNCTION", change: "DROP FUNCTION nosuch;\nDROP FUNCTION f;\nDROP FUNCTION f;\nDROP FUNCTION IF EXISTS f;\nCREATE FUNCTION f2() RETURNS int DETERMINISTIC RETURN 1;\nDROP DATABASE shop;\nDROP FUNCTION IF EXISTS f2;\nDROP FUNCTION f2;\nDROP PROCEDURE IF EXISTS p;"}},
	// A dropped database's events stay listed under its name, and the sync
	// lists them; their names are free.
	{walkCase: walkCase{name: "events of a dropped database", setup: "CREATE EVENT gh ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nDROP DATABASE shop;\nCREATE DATABASE shop;\nUSE shop;\nCREATE TABLE r (id int PRIMARY KEY);", change: "CREATE EVENT gh ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nCREATE EVENT GH ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nDROP EVENT nosuch;\nALTER EVENT nosuch COMMENT 'x';\nDROP EVENT gh;\nDROP EVENT gh;\nCREATE EVENT `shop.e3` ON SCHEDULE EVERY 1 DAY DO SELECT 1;\nDROP EVENT e3;"}},
	// The statement's scope keyword reaches only its first assignment, and
	// an assignment after one to DEFAULT does not take its value.
	{walkCase: walkCase{name: "SET foreign_key_checks", change: "SET GLOBAL foreign_key_checks = 1, foreign_key_checks = 0;\nCREATE TABLE y (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET GLOBAL foreign_key_checks = 0;\nSET sql_mode = DEFAULT, foreign_key_checks = 1;\nCREATE TABLE z (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET GLOBAL foreign_key_checks = 1;\nSET sql_mode = DEFAULT, foreign_key_checks = 0;\nCREATE TABLE w (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET foreign_key_checks = 0, sql_mode = DEFAULT;\nCREATE TABLE w2 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));\nSET foreign_key_checks = 1;\nCREATE TABLE w3 (a int, FOREIGN KEY (a) REFERENCES nosuch (id));"}},
	{walkCase: walkCase{name: "triggers", change: "CREATE TRIGGER tr BEFORE UPDATE ON trg FOR EACH ROW SET NEW.a = 2;\nCREATE TRIGGER tv BEFORE INSERT ON v FOR EACH ROW SET NEW.id = 1;\nCREATE TRIGGER IF NOT EXISTS tr BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;\nCREATE TRIGGER TR BEFORE INSERT ON nokey FOR EACH ROW SET NEW.a = 1;"}},
}

// oceanBaseCodes lists the error codes OceanBase gives for a finding
// beside MySQL's (spec §4.8).
func oceanBaseCodes(m string) []string {
	var codes []string
	for _, c := range []struct {
		match bool
		codes []string
	}{
		{strings.HasPrefix(m, "Foreign key target "), []string{"1146", "1215"}},
		{strings.HasPrefix(m, "Foreign key ") && (strings.Contains(m, " specified twice ") || strings.Contains(m, " added twice ")), []string{"1215", "1235"}},
		{strings.HasPrefix(m, "Index ") && strings.Contains(m, " dropped twice "), []string{"1072", "1235"}},
		{strings.HasPrefix(m, "Multiple primary keys"), []string{"1062", "1235"}},
		{strings.HasSuffix(m, " is not a base table"), []string{"1235"}},
		{strings.HasPrefix(m, "Database ") && strings.HasSuffix(m, " does not exist"), []string{"1210"}},
		{m == noDatabaseMessage, []string{"1210"}},
		{strings.HasPrefix(m, "Trigger ") && strings.HasSuffix(m, " already exists"), []string{"5734", "5735"}},
		{strings.HasPrefix(m, "Identifier name "), []string{"1406"}},
	} {
		if c.match {
			codes = append(codes, c.codes...)
		}
	}
	return codes
}

// oceanBaseRelation reads whether a relation of the database exists, and
// its kind: a table, a view, or none for an external table or a
// materialized view, which information_schema.TABLES does not list.
func (e *mysqlEngine) oceanBaseRelation(t testing.TB, db, name string) (bool, string) {
	t.Helper()
	ctx := context.Background()
	var kind string
	err := retryTransient(func() error {
		return e.db.QueryRowContext(ctx, "SELECT TABLE_TYPE FROM information_schema.TABLES WHERE LOWER(TABLE_SCHEMA) = ? AND LOWER(TABLE_NAME) = ?", db, name).Scan(&kind)
	})
	switch {
	case err == nil && kind == "VIEW":
		return true, relView
	case err == nil && kind == "BASE TABLE":
		return true, relTable
	case err == nil:
		return true, ""
	case !errors.Is(err, sql.ErrNoRows):
		t.Fatal(err)
	}
	var n int
	if err := retryTransient(func() error {
		return e.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM oceanbase.DBA_MVIEWS WHERE LOWER(OWNER) = ? AND LOWER(MVIEW_NAME) = ?", db, name).Scan(&n)
	}); err != nil {
		t.Fatal(err)
	}
	return n > 0, ""
}

// transient reports an OceanBase error a retry may clear: the tenant's log
// stream between leaders (4767), or a timeout.
func transient(err error) bool {
	var me *gomysql.MySQLError
	if errors.As(err, &me) {
		switch me.Number {
		case 4767, 4038, 4012, 4023, 6005, 4721, 4109:
			return true
		}
	}
	return false
}

// retryTransient runs f again while it returns a transient error, for up
// to three minutes.
func retryTransient(f func() error) error {
	err := f()
	for i := 0; i < 90 && transient(err); i++ {
		time.Sleep(2 * time.Second)
		err = f()
	}
	return err
}

// oceanBaseDB opens the server OCEANBASE_WALK_DSN names, or skips.
func oceanBaseDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("OCEANBASE_WALK_DSN")
	if dsn == "" {
		t.Skip("OCEANBASE_WALK_DSN is not set: run the test through scripts/oceanbase-oracle.sh")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := retryTransient(func() error { return db.PingContext(context.Background()) }); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestWalkThroughAgainstOceanBase(t *testing.T) {
	db := oceanBaseDB(t)
	cases := oceanBaseWalkCases
	for _, c := range walkCases {
		if _, skip := oceanBaseSkipped[c.name]; skip {
			continue
		}
		if _, ok := oceanBaseIncomplete[c.name]; ok {
			c.complete = false
		}
		cases = append(cases, oceanBaseCase{walkCase: c})
	}
	var found, rejected int
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resetOceanBase(t, db, oceanBaseWalkBase+tt.setup)
			f, r := judgeOceanBase(t, db, tt.change, tt.complete, tt.unread)
			found += f
			rejected += r
		})
	}
	t.Logf("found %d of the %d statements the server rejects", found, rejected)
}

// resetOceanBase drops every user database and recreates shop from base,
// with the global foreign_key_checks back on.
func resetOceanBase(t *testing.T, db *sql.DB, base string) {
	t.Helper()
	ctx := context.Background()
	var names []string
	if err := retryTransient(func() error {
		names = nil
		rows, err := db.QueryContext(ctx, "SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME NOT IN ("+oceanBaseSystem+")")
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var n string
			if err := rows.Scan(&n); err != nil {
				return err
			}
			names = append(names, n)
		}
		return rows.Err()
	}); err != nil {
		t.Fatal(err)
	}
	conn := session(t, db)
	stmts := []string{"SET GLOBAL foreign_key_checks = 1", "SET foreign_key_checks = 0"}
	for _, n := range names {
		stmts = append(stmts, "DROP DATABASE `"+strings.ReplaceAll(n, "`", "``")+"`")
	}
	stmts = append(stmts, "SET foreign_key_checks = 1", "CREATE DATABASE shop", "USE shop")
	for _, q := range append(stmts, splitWalk(base)...) {
		if err := retryTransient(func() error {
			_, err := conn.ExecContext(ctx, q)
			return err
		}); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
}

// judgeOceanBase runs the change on a fresh connection with shop selected,
// as Bytebase runs it, and holds the walk to it (walktest.Judge).
func judgeOceanBase(t *testing.T, db *sql.DB, change string, complete, unread bool) (found, rejected int) {
	t.Helper()
	conn := session(t, db)
	for _, q := range []string{"USE shop", "SET foreign_key_checks = 1"} {
		if _, err := conn.ExecContext(context.Background(), q); err != nil {
			t.Fatal(err)
		}
	}
	e := &mysqlEngine{
		conn:      conn,
		db:        db,
		fold:      true,
		oceanBase: true,
		target: review.Target{
			Schema:              oceanBaseSync(t, db, unread),
			LowerCaseTableNames: 1,
			ForeignKeyChecks:    review.SwitchOn,
		},
	}
	return walktest.Judge(t, e, change, complete)
}

// oceanBaseSync reads shop into the metadata Bytebase's sync produces from
// OceanBase for what WalkThrough reads (spec §2): base tables with their
// columns, indexes, foreign keys, and triggers; views; external tables;
// materialized views from oceanbase.DBA_MVIEWS, none when unread, as for a
// user who may not read it; routines; and events by name, each once. No
// CHECK constraint or partition is read.
//
// bytebase: backend/plugin/db/mysql/sync.go — SyncDBSchema, since
// bytebase/bytebase#21685 and bytebase/bytebase#21687
func oceanBaseSync(t *testing.T, db *sql.DB, unread bool) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	ctx := context.Background()
	query := func(q string, scan func(a, b string)) {
		t.Helper()
		if err := retryTransient(func() error {
			rows, err := db.QueryContext(ctx, q)
			if err != nil {
				return err
			}
			defer rows.Close()
			var got [][2]string
			for rows.Next() {
				var a, b sql.NullString
				if err := rows.Scan(&a, &b); err != nil {
					return err
				}
				got = append(got, [2]string{a.String, b.String})
			}
			if err := rows.Err(); err != nil {
				return err
			}
			for _, r := range got {
				scan(r[0], r[1])
			}
			return nil
		}); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	s := &metadata.SchemaMetadata{}
	tables := make(map[string]*metadata.TableMetadata)
	externals := make(map[string]*metadata.ExternalTableMetadata)
	query("SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME", func(name, kind string) {
		switch kind {
		case "VIEW":
			s.Views = append(s.Views, &metadata.ViewMetadata{Name: name})
		case "EXTERNAL TABLE":
			x := &metadata.ExternalTableMetadata{Name: name}
			externals[name] = x
			s.ExternalTables = append(s.ExternalTables, x)
		default:
			x := &metadata.TableMetadata{Name: name}
			tables[name] = x
			s.Tables = append(s.Tables, x)
		}
	})
	query("SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION", func(tbl, col string) {
		c := &metadata.ColumnMetadata{Name: col}
		if x, ok := tables[tbl]; ok {
			x.Columns = append(x.Columns, c)
		} else if x, ok := externals[tbl]; ok {
			x.Columns = append(x.Columns, c)
		}
	})
	indexes := make(map[[2]string]bool)
	query("SELECT TABLE_NAME, INDEX_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX", func(tbl, idx string) {
		x, ok := tables[tbl]
		if !ok || indexes[[2]string{tbl, idx}] {
			return
		}
		indexes[[2]string{tbl, idx}] = true
		x.Indexes = append(x.Indexes, &metadata.IndexMetadata{Name: idx, Primary: idx == "PRIMARY"})
	})
	query("SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'", func(tbl, name string) {
		if x, ok := tables[tbl]; ok {
			x.ForeignKeys = append(x.ForeignKeys, &metadata.ForeignKeyMetadata{Name: name})
		}
	})
	query("SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'", func(tbl, name string) {
		if x, ok := tables[tbl]; ok {
			x.Triggers = append(x.Triggers, &metadata.TriggerMetadata{Name: name})
		}
	})
	query("SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop'", func(kind, name string) {
		if kind == "PROCEDURE" {
			s.Procedures = append(s.Procedures, &metadata.ProcedureMetadata{Name: name})
		} else {
			s.Functions = append(s.Functions, &metadata.FunctionMetadata{Name: name})
		}
	})
	query("SELECT DISTINCT EVENT_NAME, '' FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME", func(name, _ string) {
		s.Events = append(s.Events, &metadata.EventMetadata{Name: name})
	})
	if !unread {
		query("SELECT MVIEW_NAME, CONTAINER_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME", func(name, _ string) {
			s.MaterializedViews = append(s.MaterializedViews, &metadata.MaterializedViewMetadata{Name: name, SkipDump: true})
		})
	}
	query("SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' AND TABLE_TYPE = 'BASE TABLE'", func(name, engine string) {
		if x, ok := tables[name]; ok {
			x.Engine = engine
		}
	})
	return &metadata.DatabaseSchemaMetadata{Name: "shop", Schemas: []*metadata.SchemaMetadata{s}}
}

// TestWalkThroughSweepOceanBase runs the MySQL sweep's random changes, with
// OceanBase's generated foreign key names. Precision must hold on every
// one; recall is logged. WALK_SWEEP_N sets the number of changes (40 by
// default; OceanBase's DDL is slow), WALK_SWEEP_SEED the seed.
func TestWalkThroughSweepOceanBase(t *testing.T) {
	db := oceanBaseDB(t)
	n, seed, _ := walktest.SweepParams(t, 40, 20261028, "")
	rng := rand.New(rand.NewSource(seed))
	var found, rejected int
	for i := range n {
		base := sweepChange(rng, 8+rng.Intn(8), "a_obfk_1")
		change := sweepChange(rng, 15+rng.Intn(15), "a_obfk_1")
		t.Run(fmt.Sprintf("change%d", i), func(t *testing.T) {
			resetOceanBase(t, db, "")
			conn := session(t, db)
			if _, err := conn.ExecContext(context.Background(), "USE shop"); err != nil {
				t.Fatal(err)
			}
			// The base's temporary tables end with its session, before the sync.
			for _, q := range splitWalk(base) {
				_ = retryTransient(func() error {
					_, err := conn.ExecContext(context.Background(), q)
					return err
				})
			}
			_ = conn.Raw(func(any) error { return driver.ErrBadConn })
			f, r := judgeOceanBase(t, db, change, false, false)
			found += f
			rejected += r
			if t.Failed() {
				t.Logf("base:\n%s\nchange:\n%s", base, change)
			}
		})
	}
	t.Logf("OceanBase seed %d: found %d of the %d statements the server rejects", seed, found, rejected)
}
