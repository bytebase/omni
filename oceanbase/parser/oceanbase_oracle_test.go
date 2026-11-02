package parser

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
)

// The OceanBase oracle. oracle_test.go holds OceanBase's inherited MySQL 8.0
// oracle, which still guards the grammar this engine shares with MySQL; this file
// holds the one that answers for OceanBase itself.
//
// OceanBase needs about 8 GB, which CI does not have, so these tests read
// OCEANBASE_WALK_DSN and skip without it:
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/parser/ -run OceanBase -count=1
const oceanBaseDSNEnv = "OCEANBASE_WALK_DSN"

// oceanBaseOracle opens the server OCEANBASE_WALK_DSN names, or skips.
func oceanBaseOracle(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv(oceanBaseDSNEnv)
	if dsn == "" {
		t.Skipf("%s is not set: run the test through scripts/oceanbase-oracle.sh", oceanBaseDSNEnv)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	// The DSN names no database and USE is session state, so pin the pool to one
	// connection: every statement then runs in the database the seed selects.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	t.Cleanup(func() { _ = db.Close() })
	// The tenant answers before it is ready to plan, so a first statement can fail
	// transiently while the server finishes starting.
	deadline := time.Now().Add(3 * time.Minute)
	for {
		err = db.PingContext(context.Background())
		if err == nil || time.Now().After(deadline) {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		t.Fatal(err)
	}
	return db
}

// engineAccepts reports whether OceanBase parses sql. It runs the statement and
// reads the error: 1064 is a parse error, anything else means the statement parsed
// and failed for a semantic reason, such as a missing table.
//
// Anything that is not a server error — a dropped connection, a timeout — is no
// verdict at all and comes back as one. Reporting it as acceptance would pass the
// whole inventory against a server that never answered, because every pending case
// expects the engine to accept and every shared case expects both sides to.
func engineAccepts(ctx context.Context, db *sql.DB, sql string) (bool, error) {
	_, err := db.ExecContext(ctx, sql)
	if err == nil {
		return true, nil
	}
	var myErr *gomysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number != 1064, nil
	}
	return false, err
}

// grammarCase is one statement the inventory asks both the engine and the parser about.
type grammarCase struct {
	name string
	sql  string
	// pending marks a statement OceanBase accepts that this parser does not read yet.
	// Implementing it means clearing the flag in the same change.
	pending bool
}

// oceanBaseOnlyGrammar is the additive side of the inventory: distributed DDL that
// MySQL rejects. Every entry is pending until its grammar lands.
//
// Read against 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-10-30: the engine parses all ten,
// so each is grammar this parser owes, not a guess at OceanBase's dialect.
var oceanBaseOnlyGrammar = []grammarCase{
	{name: "create tablegroup", sql: "CREATE TABLEGROUP obprobe_tg1 PARTITION BY HASH PARTITIONS 8", pending: true},
	{name: "alter table set tablegroup", sql: "ALTER TABLE obprobe_t SET TABLEGROUP obprobe_tg1", pending: true},
	{name: "tablegroup table option", sql: "CREATE TABLE obprobe_t2 (a INT) TABLEGROUP = obprobe_tg1", pending: true},
	{name: "create outline", sql: "CREATE OUTLINE obprobe_ol1 ON SELECT * FROM obprobe_t", pending: true},
	{name: "create resource unit", sql: "CREATE RESOURCE UNIT obprobe_u1 MAX_CPU = 2, MEMORY_SIZE = '2G'", pending: true},
	{name: "create resource pool", sql: "CREATE RESOURCE POOL obprobe_p1 UNIT = 'obprobe_u1', UNIT_NUM = 1", pending: true},
	{name: "create external table", sql: "CREATE EXTERNAL TABLE obprobe_ext (a INT) LOCATION = '/tmp/obprobe' FORMAT (TYPE = 'CSV')", pending: true},
	{name: "create materialized view", sql: "CREATE MATERIALIZED VIEW obprobe_mv1 REFRESH COMPLETE ON DEMAND AS SELECT a FROM obprobe_t", pending: true},
	{name: "column group", sql: "CREATE TABLE obprobe_t3 (a INT) WITH COLUMN GROUP(all columns, each column)", pending: true},
	{name: "alter system major freeze", sql: "ALTER SYSTEM MAJOR FREEZE", pending: true},

	// The storage and index clauses OceanBase prints in SHOW CREATE TABLE, and the
	// index scope it takes in DDL. Read from a 4.3.5 readback and from DDL a customer
	// reported as rejected, then implemented, so none of these is pending.
	{name: "organization index", sql: "CREATE TABLE obprobe_oi (a INT PRIMARY KEY) ORGANIZATION INDEX"},
	{name: "storage options", sql: "CREATE TABLE obprobe_so (a INT) REPLICA_NUM = 1 BLOCK_SIZE = 16384 TABLET_SIZE = 134217728 PCTFREE = 0 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE"},
	{name: "auto increment mode", sql: "CREATE TABLE obprobe_ai (a INT AUTO_INCREMENT PRIMARY KEY) AUTO_INCREMENT_MODE = 'ORDER'"},
	{name: "local index", sql: "CREATE TABLE obprobe_li (a INT, KEY obprobe_ia (a) LOCAL)"},
	{name: "global index", sql: "CREATE TABLE obprobe_gi (a INT, KEY obprobe_ia (a) GLOBAL)"},
	{name: "index block size", sql: "CREATE TABLE obprobe_ib (a INT, KEY obprobe_ia (a) BLOCK_SIZE 16384)"},
	{name: "index block size then scope", sql: "CREATE TABLE obprobe_ibl (a INT, KEY obprobe_ia (a) BLOCK_SIZE 16384 LOCAL)"},
	{name: "alter add index global", sql: "ALTER TABLE obprobe_t ADD UNIQUE INDEX obprobe_uk (a) GLOBAL"},
	{name: "subpartition template", sql: "CREATE TABLE obprobe_sp (a INT, b INT) PARTITION BY RANGE COLUMNS(a) SUBPARTITION BY KEY(b) SUBPARTITION TEMPLATE (SUBPARTITION p0, SUBPARTITION p1) (PARTITION q0 VALUES LESS THAN (10))"},
	// OceanBase takes this clause and then refuses it with 1235, because the table is
	// not partitioned. engineAccepts reads only 1064 as a syntax refusal, so the entry
	// pins the grammar, which is all this parser owes.
	{name: "dynamic partition policy", sql: "CREATE TABLE obprobe_dp (a INT) DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE, TIME_UNIT = 'DAY')"},
}

// sharedGrammar is the subtractive side: MySQL syntax this parser accepts, which
// OceanBase may not. A disagreement here is a parser that is too permissive.
//
// Read against 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-10-30: the engine parses all
// twelve. Its version string reports 5.7, but its grammar is not MySQL 5.7's, so the
// inherited 8.0 keyword handling is right and nothing here needs withdrawing. These
// are the cases most likely to break that assumption, not proof of the whole 8.0
// surface; widen the list when a release moves.
var sharedGrammar = []grammarCase{
	{name: "rename column", sql: "ALTER TABLE obprobe_t RENAME COLUMN a TO b"},
	{name: "cte", sql: "WITH c AS (SELECT 1 AS x) SELECT x FROM c"},
	{name: "window function", sql: "SELECT ROW_NUMBER() OVER (ORDER BY a) FROM obprobe_t"},
	{name: "check constraint", sql: "CREATE TABLE obprobe_t4 (a INT, CONSTRAINT obprobe_ck CHECK (a > 0))"},
	{name: "generated column", sql: "CREATE TABLE obprobe_t5 (a INT, b INT GENERATED ALWAYS AS (a + 1) STORED)"},
	{name: "json table", sql: "SELECT * FROM JSON_TABLE('[1]', '$[*]' COLUMNS (v INT PATH '$')) AS jt"},
	{name: "lateral derived table", sql: "SELECT * FROM obprobe_t, LATERAL (SELECT 1) AS d"},
	{name: "invisible column", sql: "CREATE TABLE obprobe_t6 (a INT INVISIBLE)"},
	{name: "functional index", sql: "CREATE TABLE obprobe_t7 (a INT, INDEX ((a + 1)))"},
	{name: "values statement", sql: "VALUES ROW(1), ROW(2)"},
	{name: "intersect", sql: "SELECT 1 INTERSECT SELECT 1"},
	{name: "except", sql: "SELECT 1 EXCEPT SELECT 2"},
}

// foreignGrammar is the guard side: syntax neither MySQL nor OceanBase has, which users
// write anyway after porting from another dialect. Both sides must refuse it, so an entry
// fails if this parser grows too permissive, and fails the other way if OceanBase ever
// adds the form, which is then grammar to implement rather than a guard to keep.
//
// Read against 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-10-31: the engine answers 1064 for
// all three. IF EXISTS on a DROP inside ALTER TABLE is MariaDB and PostgreSQL syntax.
var foreignGrammar = []grammarCase{
	{name: "alter drop index if exists", sql: "ALTER TABLE obprobe_t DROP INDEX IF EXISTS obprobe_ia"},
	{name: "alter drop column if exists", sql: "ALTER TABLE obprobe_t DROP COLUMN IF EXISTS a"},
	{name: "drop index if exists on", sql: "DROP INDEX IF EXISTS obprobe_ia ON obprobe_t"},
}

// TestOceanBaseGrammarInventory holds this parser to OceanBase rather than to MySQL.
// For each statement it compares what OceanBase accepts with what Parse accepts, and
// pins the gap: a pending case is one OceanBase takes and the parser does not yet, so
// implementing that grammar clears the flag in the same change, and a non-pending
// disagreement is a parser that reads syntax OceanBase refuses.
func TestOceanBaseGrammarInventory(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()

	// The DSN selects no database, and a base table for the statements that
	// reference one. The seed's own DDL is plain MySQL.
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS obprobe",
		"USE obprobe",
		"DROP TABLE IF EXISTS obprobe_t",
		"CREATE TABLE obprobe_t (a INT PRIMARY KEY)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obprobe")
	})

	for _, group := range [][]grammarCase{oceanBaseOnlyGrammar, sharedGrammar, foreignGrammar} {
		for _, c := range group {
			t.Run(c.name, func(t *testing.T) {
				_, parseErr := Parse(c.sql)
				parserTakes := parseErr == nil
				engineTakes, err := engineAccepts(ctx, db, c.sql)
				if err != nil {
					t.Fatalf("the oracle gave no verdict on %q: %v", c.sql, err)
				}

				switch {
				case engineTakes && !parserTakes:
					if !c.pending {
						t.Errorf("OceanBase accepts %q and the parser does not:\n  %v", c.sql, parseErr)
					}
				case !engineTakes && parserTakes:
					t.Errorf("the parser accepts %q and OceanBase does not", c.sql)
				case engineTakes && parserTakes:
					if c.pending {
						t.Errorf("%q is marked pending but the parser now reads it: clear the flag", c.sql)
					}
				default:
					// Neither takes it. For a pending case that means the statement as
					// written is not what OceanBase accepts, so the inventory is wrong.
					if c.pending {
						t.Errorf("neither OceanBase nor the parser accepts %q: the inventory entry is wrong, not pending", c.sql)
					}
				}
			})
		}
	}
}

// TestEngineAcceptsReportsConnectionFaults pins that a server which never answers
// is no verdict. Reporting a fault as acceptance would pass the whole inventory
// against a dead oracle: every pending case expects the engine to accept and the
// parser to refuse, and every shared case expects both to accept. This test needs no
// DSN, so it guards that in CI where the oracle itself cannot run.
func TestEngineAcceptsReportsConnectionFaults(t *testing.T) {
	// Port 1 is privileged and unused, so the dial is refused rather than hanging.
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:1)/")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	accepted, err := engineAccepts(context.Background(), db, "SELECT 1")
	if err == nil {
		t.Fatal("a refused connection must come back as an error, not a verdict")
	}
	if accepted {
		t.Error("a refused connection must not report the statement as accepted")
	}
}

// oceanBaseVersion records what the oracle reports, so a divergence read from one
// release is not silently reused for another.
func TestOceanBaseOracleVersion(t *testing.T) {
	db := oceanBaseOracle(t)
	var version string
	if err := db.QueryRowContext(context.Background(), "SELECT VERSION()").Scan(&version); err != nil {
		t.Fatal(err)
	}
	t.Logf("oracle reports VERSION() = %q", version)
	if !strings.Contains(version, "OceanBase") {
		t.Fatalf("%s points at a server that is not OceanBase: %q", oceanBaseDSNEnv, version)
	}
}
