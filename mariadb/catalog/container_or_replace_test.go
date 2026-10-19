package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	tcmariadb "github.com/testcontainers/testcontainers-go/modules/mariadb"
)

// orReplaceCase is a script run statement by statement on MariaDB and on the
// catalog. codes holds the error code MariaDB 11.8.8 returns for each
// statement (0 for success); both the engine and the catalog must return
// it. check names the tables (db.name, or a name in orp) whose final state
// must agree: absent, a view, or SHOW CREATE TABLE.
type orReplaceCase struct {
	name  string
	stmts []string
	codes []int
	check []string
}

// orReplaceMessageCodes are the errors whose message the catalog must also
// print as MariaDB does: the ones CREATE OR REPLACE TABLE and LIKE raise.
var orReplaceMessageCodes = map[int]bool{
	ErrNonUniqTable:    true,
	ErrUpdateTableUsed: true,
	ErrCantReopenTable: true,
	ErrWrongObject:     true,
	ErrRowIsReferenced: true,
	ErrIsAView:         true,
	ErrNoSuchTable:     true,

	ErrTableMustHaveColumns: true,
}

var orReplaceCases = []orReplaceCase{
	{
		name:  "replaces a table",
		stmts: []string{"CREATE TABLE r (id INT)", "CREATE OR REPLACE TABLE r (id INT, v INT)"},
		codes: []int{0, 0},
		check: []string{"r"},
	},
	{
		name:  "creates a table that does not exist",
		stmts: []string{"CREATE OR REPLACE TABLE n (id INT PRIMARY KEY)"},
		codes: []int{0},
		check: []string{"n"},
	},
	{
		name: "refuses to drop a foreign key parent",
		stmts: []string{
			"CREATE TABLE par (id INT PRIMARY KEY)",
			"CREATE TABLE chi (pid INT, CONSTRAINT fk_c FOREIGN KEY (pid) REFERENCES par (id))",
			"CREATE OR REPLACE TABLE par (id INT PRIMARY KEY, x INT)",
		},
		codes: []int{0, 0, 1451},
		check: []string{"par", "chi"},
	},
	{
		name: "refuses to drop a foreign key parent of another database, until the child goes",
		stmts: []string{
			"CREATE TABLE par (id INT PRIMARY KEY)",
			"CREATE TABLE orp2.chi (pid INT, CONSTRAINT fk_x FOREIGN KEY (pid) REFERENCES orp.par (id))",
			"CREATE OR REPLACE TABLE par (id INT PRIMARY KEY)",
			"CREATE OR REPLACE TABLE orp2.chi (pid INT)",
			"CREATE OR REPLACE TABLE par (id INT)",
		},
		codes: []int{0, 0, 1451, 0, 0},
		check: []string{"par", "orp2.chi"},
	},
	{
		name: "replaces a foreign key parent with foreign_key_checks off",
		stmts: []string{
			"SET foreign_key_checks = 0",
			"CREATE TABLE par (id INT PRIMARY KEY)",
			"CREATE TABLE chi (pid INT, CONSTRAINT fk_c FOREIGN KEY (pid) REFERENCES par (id))",
			"CREATE OR REPLACE TABLE par (id INT PRIMARY KEY, y INT)",
			"CREATE OR REPLACE TABLE par (z INT)",
		},
		codes: []int{0, 0, 0, 0, 0},
		check: []string{"par", "chi"},
	},
	{
		name: "replaces a table whose foreign key references itself",
		stmts: []string{
			"CREATE TABLE s (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_s FOREIGN KEY (pid) REFERENCES s (id))",
			"CREATE OR REPLACE TABLE s (id INT PRIMARY KEY)",
		},
		codes: []int{0, 0},
		check: []string{"s"},
	},
	{
		name: "refuses a view of the name",
		stmts: []string{
			"CREATE VIEW vw AS SELECT 1 AS a",
			"CREATE OR REPLACE TABLE vw (id INT)",
			"CREATE OR REPLACE TABLE vw AS SELECT 1 AS a",
			"CREATE OR REPLACE TABLE vw LIKE nosuch",
			"CREATE OR REPLACE TABLE vw LIKE vw",
			"CREATE TABLE vw (id INT)",
		},
		codes: []int{0, 1965, 1965, 1146, 1066, 1050},
		check: []string{"vw"},
	},
	{
		name:  "replaces a temporary table",
		stmts: []string{"CREATE TEMPORARY TABLE t1 (a INT)", "CREATE OR REPLACE TEMPORARY TABLE t1 (a INT, b INT)"},
		codes: []int{0, 0},
		check: []string{"t1"},
	},
	{
		name:  "creates a temporary table that does not exist",
		stmts: []string{"CREATE OR REPLACE TEMPORARY TABLE t2 (x INT)"},
		codes: []int{0},
		check: []string{"t2"},
	},
	{
		name: "replaces a table with LIKE",
		stmts: []string{
			"CREATE TABLE src (id INT PRIMARY KEY, s VARCHAR(5))",
			"CREATE TABLE dst (a INT)",
			"CREATE OR REPLACE TABLE dst LIKE src",
		},
		codes: []int{0, 0, 0},
		check: []string{"src", "dst"},
	},
	{
		name: "refuses LIKE of the table itself, existing or not",
		stmts: []string{
			"CREATE TABLE src (id INT PRIMARY KEY)",
			"CREATE OR REPLACE TABLE src LIKE src",
			"CREATE OR REPLACE TABLE orp.src LIKE src",
			"CREATE TABLE src LIKE src",
			"CREATE TABLE IF NOT EXISTS src LIKE src",
			"CREATE TABLE nx LIKE nx",
			"CREATE OR REPLACE TABLE nx LIKE nx",
			"CREATE TEMPORARY TABLE tm (a INT)",
			"CREATE OR REPLACE TEMPORARY TABLE tm LIKE tm",
		},
		codes: []int{0, 1066, 1066, 1066, 1066, 1066, 1066, 0, 1066},
		check: []string{"src", "nx", "tm"},
	},
	{
		name: "keeps the table when the LIKE source is missing or a view",
		stmts: []string{
			"CREATE TABLE dst (a INT)",
			"CREATE VIEW sv AS SELECT 1 AS a",
			"CREATE OR REPLACE TABLE dst LIKE nosuch",
			"CREATE OR REPLACE TABLE dst LIKE nosuchdb.x",
			"CREATE OR REPLACE TABLE dst LIKE sv",
			"CREATE TABLE fresh LIKE sv",
			"CREATE TABLE fresh LIKE nosuchdb.x",
		},
		codes: []int{0, 0, 1146, 1146, 1347, 1347, 1146},
		check: []string{"dst", "fresh"},
	},
	{
		name: "refuses a foreign key parent replaced by LIKE or CREATE TABLE AS SELECT",
		stmts: []string{
			"CREATE TABLE src (id INT PRIMARY KEY)",
			"CREATE TABLE par (id INT PRIMARY KEY)",
			"CREATE TABLE chi (pid INT, CONSTRAINT fk_c FOREIGN KEY (pid) REFERENCES par (id))",
			"CREATE OR REPLACE TABLE par LIKE src",
			"CREATE OR REPLACE TABLE par LIKE nosuch",
			"CREATE OR REPLACE TABLE par AS SELECT 1 AS id",
			"CREATE OR REPLACE TABLE par AS SELECT * FROM par",
		},
		codes: []int{0, 0, 0, 1451, 1146, 1451, 1093},
		check: []string{"par", "chi"},
	},
	{
		name: "refuses CREATE TABLE AS SELECT that reads the table it replaces",
		stmts: []string{
			"CREATE TABLE src (id INT PRIMARY KEY)",
			"CREATE OR REPLACE TABLE src AS SELECT * FROM src",
			"CREATE OR REPLACE TABLE src AS SELECT y.id FROM src AS y",
			"CREATE OR REPLACE TABLE src AS SELECT * FROM orp.src",
			"CREATE OR REPLACE TABLE src AS SELECT 1 AS a FROM dual WHERE EXISTS (SELECT 1 FROM src)",
			"CREATE OR REPLACE TABLE src AS SELECT 1 AS a FROM (SELECT id FROM src) d",
			"CREATE OR REPLACE TABLE src AS SELECT (SELECT MAX(id) FROM src) AS m",
			"CREATE OR REPLACE TABLE src AS SELECT 1 AS a UNION SELECT id FROM src",
			"CREATE OR REPLACE TABLE src AS WITH c AS (SELECT id FROM src) SELECT * FROM c",
			"CREATE OR REPLACE TABLE src (x INT) AS SELECT * FROM src",
		},
		codes: []int{0, 1093, 1093, 1093, 1093, 1093, 1093, 1093, 1093, 1093},
		check: []string{"src"},
	},
	{
		// The catalog does not model CREATE TABLE AS SELECT, so orp.src is
		// not compared: only that both accept.
		name: "accepts CREATE TABLE AS SELECT that reads another table of the name",
		stmts: []string{
			"CREATE TABLE src (id INT PRIMARY KEY)",
			"CREATE TABLE orp2.src (q INT)",
			"CREATE OR REPLACE TABLE src AS SELECT q FROM orp2.src",
			"CREATE TABLE src2 (id INT)",
			"CREATE OR REPLACE TABLE src2 AS WITH src2 AS (SELECT 1 AS a) SELECT * FROM src2",
		},
		codes: []int{0, 0, 0, 0, 0},
		check: []string{"orp2.src"},
	},
	{
		// A CTE hides the table only in the query blocks it is visible in: a
		// plain WITH body sees the CTEs before it, not itself or later ones.
		name: "refuses a SELECT that reads the table outside the CTE of its name",
		stmts: []string{
			"CREATE TABLE src (a INT)",
			"CREATE OR REPLACE TABLE src AS SELECT * FROM src WHERE EXISTS (WITH src AS (SELECT 1) SELECT * FROM src)",
			"CREATE OR REPLACE TABLE src AS WITH src AS (SELECT * FROM src) SELECT * FROM src",
			"CREATE OR REPLACE TABLE src AS WITH c1 AS (SELECT * FROM src), src AS (SELECT 1 AS a) SELECT * FROM c1",
			"CREATE OR REPLACE TABLE src AS SELECT (WITH src AS (SELECT 1 AS a) SELECT a FROM src) AS x, (SELECT COUNT(*) FROM src) AS y",
		},
		codes: []int{0, 1093, 1093, 1093, 1093},
		check: []string{"src"},
	},
	{
		// As above, the replaced tables are not compared: the catalog does
		// not model CREATE TABLE AS SELECT.
		name: "accepts a SELECT whose every reference of the name reads a visible CTE",
		stmts: []string{
			"CREATE TABLE s1 (a INT)",
			"CREATE OR REPLACE TABLE s1 AS WITH RECURSIVE s1 AS (SELECT 1 AS a UNION SELECT a + 1 FROM s1 WHERE a < 3) SELECT * FROM s1",
			"CREATE TABLE s2 (a INT)",
			"CREATE OR REPLACE TABLE s2 AS WITH s2 AS (SELECT 1 AS a), c1 AS (SELECT * FROM s2) SELECT * FROM c1",
			"CREATE TABLE s3 (a INT)",
			"CREATE OR REPLACE TABLE s3 AS SELECT * FROM (WITH s3 AS (SELECT 1 AS a) SELECT * FROM s3) d",
			"CREATE TABLE s4 (a INT)",
			"CREATE OR REPLACE TABLE s4 AS WITH s4 AS (SELECT 1 AS a) SELECT * FROM s4 WHERE EXISTS (SELECT 1 FROM s4)",
			"CREATE TABLE s5 (a INT)",
			"CREATE OR REPLACE TABLE s5 AS WITH s5 AS (SELECT 1 AS a) SELECT * FROM s5 UNION SELECT a FROM s5",
			"CREATE TABLE s6 (a INT)",
			"CREATE OR REPLACE TABLE s6 AS WITH RECURSIVE c1 AS (SELECT * FROM s6), s6 AS (SELECT 1 AS a) SELECT * FROM c1",
		},
		codes: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		// A table needs a column: 1113, raised after OR REPLACE has dropped
		// the old table and its trigger.
		name: "a definition without columns is 1113",
		stmts: []string{
			"CREATE TABLE r (id INT)",
			"CREATE TRIGGER r_bi BEFORE INSERT ON r FOR EACH ROW SET NEW.id = 1",
			"CREATE OR REPLACE TABLE r",
			"CREATE TABLE r2 (id INT)",
			"CREATE OR REPLACE TABLE r2 (INDEX (id))",
			"CREATE TABLE r3",
			"CREATE TABLE r4 (CHECK (1 > 0))",
		},
		codes: []int{0, 0, 1113, 0, 1113, 1113, 1113},
		check: []string{"r", "r2", "r3", "r4"},
	},
	{
		name: "replaces a table with CREATE TABLE AS SELECT of declared columns",
		stmts: []string{
			"CREATE TABLE src (id INT PRIMARY KEY)",
			"CREATE OR REPLACE TABLE src (x INT) SELECT 1 AS x",
		},
		codes: []int{0, 0},
		check: []string{"src"},
	},
	{
		name: "refuses a temporary CREATE TABLE AS SELECT that reads the table it replaces",
		stmts: []string{
			"CREATE TEMPORARY TABLE tm (a INT)",
			"CREATE OR REPLACE TEMPORARY TABLE tm AS SELECT x.a FROM tm AS x",
			"CREATE OR REPLACE TEMPORARY TABLE tm AS SELECT 1 AS a FROM dual WHERE EXISTS (SELECT 1 FROM tm)",
		},
		codes: []int{0, 1137, 1137},
		check: []string{"tm"},
	},
	{
		// MariaDB drops the old table before it puts the new definition
		// together, and does not bring it back when that fails.
		name: "a definition that fails after the drop leaves no table",
		stmts: []string{
			"CREATE TABLE k1 (a INT)",
			"CREATE OR REPLACE TABLE k1 (a INT, a INT)",
			"CREATE TABLE k2 (a INT)",
			"CREATE OR REPLACE TABLE k2 (a INT PRIMARY KEY, b INT PRIMARY KEY)",
			"CREATE TABLE k3 (a INT)",
			"CREATE OR REPLACE TABLE k3 (a INT AUTO_INCREMENT PRIMARY KEY, b INT AUTO_INCREMENT)",
			"CREATE TABLE k4 (a INT)",
			"CREATE OR REPLACE TABLE k4 (a INT, KEY i (a), KEY i (a))",
		},
		codes: []int{0, 1060, 0, 1068, 0, 1075, 0, 1061},
		check: []string{"k1", "k2", "k3", "k4"},
	},
	{
		// What MariaDB rejects while parsing fails before anything is dropped.
		name: "a definition rejected before the drop keeps the table",
		stmts: []string{
			"CREATE TABLE k1 (a INT)",
			"CREATE OR REPLACE TABLE k1 (a TIME(7))",
			"CREATE TABLE k2 (a INT)",
			"CREATE OR REPLACE TABLE k2 (a DECIMAL(5,6))",
			"CREATE TABLE k3 (a INT)",
			"CREATE OR REPLACE TABLE k3 (a INT NOT NULL DEFAULT NULL)",
			"CREATE TABLE k4 (a INT)",
			"CREATE OR REPLACE TABLE k4 (a INT) CHARSET latin1 COLLATE utf8mb4_bin",
		},
		codes: []int{0, 1426, 0, 1427, 0, 1067, 0, 1253},
		check: []string{"k1", "k2", "k3", "k4"},
	},
	{
		name: "drops the triggers of the table it replaces, even when the replacement fails",
		stmts: []string{
			"CREATE TABLE tg (a INT)",
			"CREATE TRIGGER tg_bi BEFORE INSERT ON tg FOR EACH ROW SET NEW.a = 1",
			"CREATE TABLE tg2 (a INT)",
			"CREATE TRIGGER tg2_bi BEFORE INSERT ON tg2 FOR EACH ROW SET NEW.a = 1",
			"CREATE TABLE keep (a INT)",
			"CREATE TRIGGER keep_bi BEFORE INSERT ON keep FOR EACH ROW SET NEW.a = 1",
			"CREATE OR REPLACE TABLE tg (a INT, b INT)",
			"CREATE OR REPLACE TABLE tg2 (a INT, a INT)",
		},
		codes: []int{0, 0, 0, 0, 0, 0, 0, 1060},
		check: []string{"tg", "tg2", "keep"},
	},
	{
		name:  "refuses an unknown database",
		stmts: []string{"CREATE OR REPLACE TABLE nosuchdb.x (a INT)"},
		codes: []int{1049},
	},
}

// TestContainerCreateOrReplaceTable runs each case on mariadb:11.8.8 and on
// the catalog and requires the same verdict for every statement and the
// same final tables and triggers.
//
// Not covered, because the catalog keeps one table per name: OR REPLACE
// TEMPORARY TABLE where a base table or view of the name exists, and OR
// REPLACE TABLE where only a temporary table of the name exists. MariaDB
// creates the new table beside the other one; the catalog reports 1050, as
// it does for plain CREATE [TEMPORARY] TABLE there.
func TestContainerCreateOrReplaceTable(t *testing.T) {
	db := startOrReplaceMariaDB(t)
	for _, tc := range orReplaceCases {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.codes) != len(tc.stmts) {
				t.Fatalf("case has %d statements and %d codes", len(tc.stmts), len(tc.codes))
			}
			ctx := context.Background()
			conn, err := db.Conn(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			for _, stmt := range []string{
				"DROP DATABASE IF EXISTS orp", "DROP DATABASE IF EXISTS orp2",
				"CREATE DATABASE orp", "CREATE DATABASE orp2", "USE orp",
			} {
				if _, err := conn.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("mariadb setup %q: %v", stmt, err)
				}
			}
			cat := New()
			for _, stmt := range []string{"CREATE DATABASE orp", "CREATE DATABASE orp2", "USE orp"} {
				if code, msg := catalogErr(t, cat, stmt); code != 0 {
					t.Fatalf("catalog setup %q: %d %s", stmt, code, msg)
				}
			}

			for i, stmt := range tc.stmts {
				engCode, engMsg := mariadbErr(t, stmt, func() error { _, err := conn.ExecContext(ctx, stmt); return err }())
				catCode, catMsg := catalogErr(t, cat, stmt)
				if engCode != tc.codes[i] {
					t.Errorf("%s\n  mariadb returned %d (%s), recorded %d", stmt, engCode, engMsg, tc.codes[i])
				}
				if catCode != engCode {
					t.Errorf("%s\n  catalog returned %d (%s), mariadb %d (%s)", stmt, catCode, catMsg, engCode, engMsg)
				} else if orReplaceMessageCodes[engCode] && catMsg != engMsg {
					t.Errorf("%s\n  catalog message %q, mariadb %q", stmt, catMsg, engMsg)
				}
			}

			for _, name := range tc.check {
				dbName, tblName := "orp", name
				if i := strings.IndexByte(name, '.'); i >= 0 {
					dbName, tblName = name[:i], name[i+1:]
				}
				eng := mariadbTableState(t, conn, dbName, tblName)
				got := catalogTableState(cat, dbName, tblName)
				if got != eng {
					t.Errorf("%s:\n  catalog: %s\n  mariadb: %s", name, got, eng)
				}
			}
			if eng, got := mariadbTriggers(t, conn), catalogTriggers(cat); eng != got {
				t.Errorf("triggers of orp: catalog [%s], mariadb [%s]", got, eng)
			}
		})
	}
}

// TestCreateOrReplaceTableAsSelectDropsTable pins how the catalog models
// CREATE OR REPLACE TABLE ... AS SELECT: it does not model the SELECT, so,
// as after plain CREATE TABLE ... AS SELECT, no table of the name remains.
func TestCreateOrReplaceTableAsSelectDropsTable(t *testing.T) {
	c := New()
	results, err := c.Exec("CREATE DATABASE d; USE d; CREATE TABLE r (id INT); CREATE TRIGGER r_bi BEFORE INSERT ON r FOR EACH ROW SET NEW.id = 1; CREATE OR REPLACE TABLE r AS SELECT 1 AS a", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("statement %d: %v", r.Index, r.Error)
		}
	}
	db := c.GetDatabase("d")
	if db.GetTable("r") != nil {
		t.Errorf("r still holds its old definition:\n%s", c.ShowCreateTable("d", "r"))
	}
	if len(db.Triggers) != 0 {
		t.Errorf("the trigger of the replaced table remains")
	}
}

// startOrReplaceMariaDB starts the mariadb:11.8.8 container
// TestContainerCreateOrReplaceTable asks. No connection is kept idle, so
// every case starts a session with no temporary tables and default
// settings.
func startOrReplaceMariaDB(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()
	ctr, err := tcmariadb.Run(ctx, "mariadb:11.8.8",
		tcmariadb.WithDatabase("test"),
		tcmariadb.WithUsername("root"),
		tcmariadb.WithPassword("test"),
	)
	if err != nil {
		t.Fatalf("start MariaDB container: %v", err)
	}
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(ctr) })
	dsn, err := ctr.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxIdleConns(0)
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}
	return db
}

// mariadbErr returns the code and message of a MariaDB error, 0 for none.
func mariadbErr(t *testing.T, stmt string, err error) (int, string) {
	t.Helper()
	if err == nil {
		return 0, ""
	}
	var myErr *gomysql.MySQLError
	if !errors.As(err, &myErr) {
		t.Fatalf("%s: not a MariaDB error: %v", stmt, err)
	}
	return int(myErr.Number), myErr.Message
}

// catalogErr executes one statement on the catalog and returns the code
// and message of its error, 0 for none.
func catalogErr(t *testing.T, c *Catalog, stmt string) (int, string) {
	t.Helper()
	results, err := c.Exec(stmt, nil)
	if err != nil {
		t.Fatalf("%s: catalog parse: %v", stmt, err)
	}
	if len(results) != 1 {
		t.Fatalf("%s: catalog returned %d results", stmt, len(results))
	}
	if results[0].Error == nil {
		return 0, ""
	}
	var catErr *Error
	if !errors.As(results[0].Error, &catErr) {
		t.Fatalf("%s: not a catalog error: %v", stmt, results[0].Error)
	}
	return catErr.Code, catErr.Message
}

// mariadbTableState describes a table as the session sees it: absent, a
// view, or its SHOW CREATE TABLE in the form normalizeShowCreate gives.
func mariadbTableState(t *testing.T, conn *sql.Conn, dbName, name string) string {
	t.Helper()
	rows, err := conn.QueryContext(context.Background(), fmt.Sprintf("SHOW CREATE TABLE `%s`.`%s`", dbName, name))
	if err != nil {
		if code, _ := mariadbErr(t, "SHOW CREATE TABLE", err); code == ErrNoSuchTable {
			return "absent"
		}
		t.Fatalf("SHOW CREATE TABLE %s.%s: %v", dbName, name, err)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		t.Fatal(err)
	}
	if len(cols) == 4 {
		return "view"
	}
	if !rows.Next() {
		t.Fatalf("SHOW CREATE TABLE %s.%s returned no row: %v", dbName, name, rows.Err())
	}
	var tbl, ddl string
	if err := rows.Scan(&tbl, &ddl); err != nil {
		t.Fatal(err)
	}
	return normalizeShowCreate(ddl)
}

func catalogTableState(c *Catalog, dbName, name string) string {
	if db := c.GetDatabase(dbName); db != nil && db.Views[toLower(name)] != nil {
		return "view"
	}
	ddl := c.ShowCreateTable(dbName, name)
	if ddl == "" {
		return "absent"
	}
	return normalizeShowCreate(ddl)
}

// normalizeShowCreate keeps SHOW CREATE TABLE up to the closing parenthesis
// of the definition, dropping table options (the catalog prints MySQL 8.0's
// default collation), and prints int(11) as the catalog does, int.
func normalizeShowCreate(ddl string) string {
	if i := strings.LastIndex(ddl, "\n)"); i >= 0 {
		ddl = ddl[:i+2]
	}
	return strings.ReplaceAll(ddl, " int(11)", " int")
}

func mariadbTriggers(t *testing.T, conn *sql.Conn) string {
	t.Helper()
	rows, err := conn.QueryContext(context.Background(), "SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'orp' ORDER BY TRIGGER_NAME")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}

func catalogTriggers(c *Catalog) string {
	var names []string
	for _, trg := range c.GetDatabase("orp").Triggers {
		names = append(names, trg.Name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
