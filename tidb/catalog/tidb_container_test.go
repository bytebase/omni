package catalog

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/internal/testcontainer"
)

// tidbCatalogContainer is a test's connection to TiDB for catalog
// cross-validation: a replay of the test's tape, or the pinned server in a
// recording (tidb/internal/testcontainer). The catalog's in-memory model is
// compared to what TiDB actually accepts (via its SQL layer) rather than
// deep-inspecting information_schema: the goal is to catch regressions where
// omni parses TiDB syntax but the catalog fails to execute, or vice versa.
type tidbCatalogContainer struct {
	db  *sql.DB
	ctx context.Context
}

// startTiDBForCatalog opens a connection of the test's own, at session
// defaults with the database test selected.
func startTiDBForCatalog(t *testing.T) *tidbCatalogContainer {
	t.Helper()
	return &tidbCatalogContainer{db: testcontainer.Open(t), ctx: context.Background()}
}

// TestTiDBCatalogContainerCrossValidation executes the same DDL on real
// TiDB and on the omni catalog and asserts both accept or reject in
// lockstep. Each case runs a (setup?, sql, cleanup?) sequence on TiDB and
// the same sql through the catalog's Exec(). The catalog doesn't
// information_schema-query TiDB — that's deferred to a diff/migration PR.
func TestTiDBCatalogContainerCrossValidation(t *testing.T) {
	tc := startTiDBForCatalog(t)

	// Isolate test tables in a database of this test's own.
	dbName := takeDatabaseOnSharedTiDB(t, tc)

	cases := []struct {
		name    string
		setup   string
		sql     string
		cleanup string
	}{
		{
			// AUTO_RANDOM on PK; SHARD_ROW_ID_BITS is tested separately
			// because TiDB forbids SHARD_ROW_ID_BITS on tables that
			// already use the PK as the row id.
			name:    "auto_random pk",
			sql:     "CREATE TABLE t_ar (id BIGINT AUTO_RANDOM(5) PRIMARY KEY)",
			cleanup: "DROP TABLE IF EXISTS t_ar",
		},
		{
			// SHARD_ROW_ID_BITS requires NO explicit PK (or non-row-id PK).
			name:    "shard_row_id_bits no pk",
			sql:     "CREATE TABLE t_srib (a INT) SHARD_ROW_ID_BITS = 4",
			cleanup: "DROP TABLE IF EXISTS t_srib",
		},
		{
			name:    "ttl create",
			sql:     "CREATE TABLE t_ttl (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'ON'",
			cleanup: "DROP TABLE IF EXISTS t_ttl",
		},
		{
			// SET TIFLASH REPLICA 0 is accepted even when no TiFlash
			// servers are available; replica counts > 0 require TiFlash
			// servers in the cluster, which a standalone TiDB container
			// does not provide. Using 0 keeps the test hermetic.
			name:    "tiflash replica alter (0)",
			setup:   "CREATE TABLE IF NOT EXISTS t_tf (id INT PRIMARY KEY)",
			sql:     "ALTER TABLE t_tf SET TIFLASH REPLICA 0",
			cleanup: "DROP TABLE IF EXISTS t_tf",
		},
		{
			name:    "remove ttl alter",
			setup:   "CREATE TABLE IF NOT EXISTS t_ttl_rm (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'ON'",
			sql:     "ALTER TABLE t_ttl_rm REMOVE TTL",
			cleanup: "DROP TABLE IF EXISTS t_ttl_rm",
		},
		{
			name:    "clustered pk",
			sql:     "CREATE TABLE t_cl (id INT, PRIMARY KEY (id) CLUSTERED)",
			cleanup: "DROP TABLE IF EXISTS t_cl",
		},
		{
			// Multi-feature interaction covers AUTO_RANDOM + CLUSTERED PK +
			// AUTO_ID_CACHE + TTL + TTL_JOB_INTERVAL on one table.
			// PLACEMENT POLICY and SHARD_ROW_ID_BITS are exercised elsewhere;
			// combining them here would require running CREATE PLACEMENT
			// POLICY first (unmodeled by catalog in this PR) or would
			// conflict with AUTO_RANDOM's row-id usage.
			name:    "multi-feature interaction",
			sql:     "CREATE TABLE t_mix (id BIGINT AUTO_RANDOM(5) PRIMARY KEY CLUSTERED, created_at DATETIME) AUTO_ID_CACHE = 100 TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'ON' TTL_JOB_INTERVAL = '1h'",
			cleanup: "DROP TABLE IF EXISTS t_mix",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.setup != "" {
				mustExecTiDB(t, tc, c.setup)
			}
			// 1. TiDB must accept.
			if _, err := tc.db.ExecContext(tc.ctx, c.sql); err != nil {
				t.Fatalf("TiDB rejected SQL %q: %v", c.sql, err)
			}
			// 2. Catalog must accept through the same path (fresh catalog
			//    per case so setup SQL on TiDB doesn't need to be mirrored).
			cat := New()
			catSetup := "CREATE DATABASE " + dbName + "; USE " + dbName + ";"
			if _, err := cat.Exec(catSetup, nil); err != nil {
				t.Fatalf("catalog setup failed: %v", err)
			}
			if c.setup != "" {
				results, err := cat.Exec(c.setup, nil)
				if err != nil {
					t.Fatalf("catalog rejected setup SQL %q: %v", c.setup, err)
				}
				for _, r := range results {
					if r.Error != nil {
						t.Fatalf("catalog exec error on setup %q: %v", c.setup, r.Error)
					}
				}
			}
			results, err := cat.Exec(c.sql, nil)
			if err != nil {
				t.Fatalf("catalog rejected SQL %q: %v", c.sql, err)
			}
			for _, r := range results {
				if r.Error != nil {
					t.Fatalf("catalog exec error on %q: %v", c.sql, r.Error)
				}
			}
			if c.cleanup != "" {
				_, _ = tc.db.ExecContext(tc.ctx, c.cleanup)
			}
		})
	}
}

// takeDatabaseOnSharedTiDB gives the test a database of its own on the
// package-shared server and drops it again when the test ends. The name is
// derived from the test's own name rather than written down, so no two tests on
// this server can collide and none has to agree on a literal (AGENTS.md: "A
// test takes a database or schema on its package's server"). It is deterministic
// across runs, so a failure is reproducible, and it is left selected as the
// session's current database, which the pinned single connection keeps.
func takeDatabaseOnSharedTiDB(t *testing.T, tc *tidbCatalogContainer) string {
	t.Helper()
	name := databaseNameForTest(t.Name())
	mustExecTiDB(t, tc, "DROP DATABASE IF EXISTS "+name)
	mustExecTiDB(t, tc, "CREATE DATABASE "+name)
	mustExecTiDB(t, tc, "USE "+name)
	t.Cleanup(func() {
		_, _ = tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS "+name)
	})
	return name
}

// databaseNameForTest turns a test name into a legal, lowercase database
// identifier: anything that is not a letter, a digit or an underscore becomes
// an underscore, and the result is cut to MySQL's 64-character limit from the
// right, which is the end that tells two tests apart.
func databaseNameForTest(testName string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(testName) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	name := b.String()
	if len(name) > 64 {
		name = name[len(name)-64:]
	}
	return name
}

func mustExecTiDB(t *testing.T, tc *tidbCatalogContainer, sqlStr string) {
	t.Helper()
	if _, err := tc.db.ExecContext(tc.ctx, sqlStr); err != nil {
		t.Fatalf("TiDB setup exec failed for %q: %v", sqlStr, err)
	}
}

// TestDatabaseNameForTest covers the identifier derivation: a subtest's '/' and
// any other character MySQL will not take in an unquoted name become
// underscores, and a name past the 64-character limit keeps its tail, which is
// the end that tells two tests apart.
func TestDatabaseNameForTest(t *testing.T) {
	cases := []struct{ in, want string }{
		{in: "TestContainerInvalidDefaultForms", want: "testcontainerinvaliddefaultforms"},
		{in: "TestX/sub name", want: "testx_sub_name"},
		{in: "TestX/enum('x','y')", want: "testx_enum__x___y__"},
		{in: "Test_Keep_Underscores9", want: "test_keep_underscores9"},
	}
	for _, c := range cases {
		if got := databaseNameForTest(c.in); got != c.want {
			t.Errorf("databaseNameForTest(%q) = %q, want %q", c.in, got, c.want)
		}
	}
	long := databaseNameForTest(strings.Repeat("a", 60) + "/tail")
	if len(long) != 64 {
		t.Errorf("a long name came back %d characters, want 64: %q", len(long), long)
	}
	if !strings.HasSuffix(long, "_tail") {
		t.Errorf("a long name dropped its tail: %q", long)
	}
}
