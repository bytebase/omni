package driver

import (
	"context"
	"crypto/tls"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/enginetest"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/redshift/internal/testcontainer"
)

// The tested paths replay tapes recorded from one scratch cluster, and record
// against it again under OMNI_RECORD=1. The tapes live beside this package and
// name the cluster's release in redshift/testdata/pins.tsv; nothing here reads
// an endpoint or a login, so a replay needs no credentials and reaches no
// server. REDSHIFT_DSN is read only while recording.
//
// DATASHARE is the exception: it needs a second account that shares a database
// into the cluster, which a scratch cluster does not have, so it skips and has
// no tape.

// dsnVariables are the environment variables a live-only test accepts, in the
// order it reads them.
var dsnVariables = []string{"REDSHIFT_DSN", "REDSHIFT_COMPAT_DSN"}

// liveConfig reads a connecting cluster's settings out of the first DSN
// variable that is set, and skips the test when none is.
func liveConfig(t *testing.T) driver.Config {
	t.Helper()
	var dsn, variable string
	for _, name := range dsnVariables {
		if value := os.Getenv(name); value != "" {
			dsn, variable = value, name
			break
		}
	}
	if dsn == "" {
		t.Skipf("no %s to reach a cluster: this test needs a real Redshift endpoint, which a replay does not have", strings.Join(dsnVariables, " or "))
	}

	pgxConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("%s does not parse: %v", variable, err)
	}
	cfg := driver.Config{
		Host:     pgxConfig.Host,
		Port:     strconv.Itoa(int(pgxConfig.Port)),
		User:     pgxConfig.User,
		Password: pgxConfig.Password,
		Database: pgxConfig.Database,
	}
	if pgxConfig.TLSConfig != nil {
		cfg.TLS = func(string) (*tls.Config, error) {
			return pgxConfig.TLSConfig, nil
		}
	}
	return cfg
}

// openLive returns the driver the test's tape was recorded through: one
// connection, to the database the recording used.
func openLive(t *testing.T) *Driver {
	t.Helper()
	d := &Driver{db: testcontainer.Open(t), databaseName: testcontainer.Database}
	if err := d.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
	return d
}

// TestSyncInstance reads the cluster's version, databases and roles, and
// pins the two things every other live test stands on: the version parses, and
// the database the driver connects to is listed.
func TestSyncInstance(t *testing.T) {
	ctx := context.Background()
	d := openLive(t)

	instance, err := d.SyncInstance(ctx)
	if err != nil {
		t.Fatalf("SyncInstance: %v", err)
	}
	if !strings.HasPrefix(instance.Version, "Redshift ") {
		t.Errorf("Version = %q, want it to name the Redshift release", instance.Version)
	}
	if len(instance.Roles) == 0 {
		t.Error("Roles is empty, want the cluster's users")
	}
	var found bool
	for _, database := range instance.Databases {
		if database.Name == d.databaseName {
			found = true
		}
		if excludedDatabaseList[database.Name] {
			t.Errorf("Databases holds %q, which a sync leaves out", database.Name)
		}
	}
	if !found {
		t.Errorf("Databases does not hold %q", d.databaseName)
	}
}

// TestSyncDatabase reads the schema of the database the driver connects to:
// the queries a sync sends, and the metadata they produce.
func TestSyncDatabase(t *testing.T) {
	d := openLive(t)
	database, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatalf("SyncDatabase: %v", err)
	}
	if database.Name != d.databaseName {
		t.Errorf("Name = %q, want %q", database.Name, d.databaseName)
	}
	var schemas []string
	for _, schema := range database.Schemas {
		schemas = append(schemas, schema.Name)
	}
	if len(schemas) == 0 {
		t.Fatalf("Schemas is empty, want at least public")
	}
	if !contains(schemas, "public") {
		t.Errorf("Schemas = %v, want it to hold public", schemas)
	}
	checkSyncGolden(t, database)
}

// checkSyncGolden writes the sync's metadata in a recording and holds it to the
// committed prototext in a replay. The comparison is proto.Equal rather than
// bytes: protojson's spacing may change between builds, and the metadata is
// what the recording settles.
func checkSyncGolden(t *testing.T, database *metadatapb.DatabaseSchemaMetadata) {
	t.Helper()
	path := filepath.Join("testdata", "sync", database.Name+".json")
	if enginetest.Recording(t) && !enginetest.Verifying(t) {
		b, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(database)
		if err != nil {
			t.Fatalf("marshal %s: %v", path, err)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("create %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
		return
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s (record it): %v", path, err)
	}
	want := &metadatapb.DatabaseSchemaMetadata{}
	if err := protojson.Unmarshal(body, want); err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !proto.Equal(database, want) {
		t.Errorf("SyncDatabase differs from %s:\n got %v\nwant %v", path, database, want)
	}
}

// TestSyncDatabaseDatashare reads a database an inbound datashare created.
// The branch reads SVV_ALL_SCHEMAS, SVV_ALL_TABLES and SVV_ALL_COLUMNS, because a
// datashare publishes no pg_catalog entries.
//
// The test finds the share in the cluster's own database list rather than naming
// one in the environment: a cluster with none skips, naming what a recording
// must provision, instead of sending SVV_ALL_* against an ordinary database.
func TestSyncDatabaseDatashare(t *testing.T) {
	ctx := context.Background()
	// The one test that reads a real endpoint rather than a tape: a share has
	// no tape, because a recording of an ordinary cluster cannot make one.
	cfg := liveConfig(t)
	d, err := Open(ctx, cfg)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	if err := d.Ping(ctx); err != nil {
		t.Fatalf("Ping: %v", err)
	}

	instance, err := d.SyncInstance(ctx)
	if err != nil {
		t.Fatalf("SyncInstance: %v", err)
	}
	var shared string
	for _, database := range instance.Databases {
		if database.Datashare {
			shared = database.Name
			break
		}
	}
	if shared == "" {
		t.Skipf("no inbound datashare on %s: this test needs a second account or namespace that shares a database into the cluster, so SVV_DATASHARES names a consumer database (VERIFY.md step 2)", d.databaseName)
	}

	// The share is reached with the same login, marked as a share.
	shareCfg := cfg
	shareCfg.Database = shared
	if shareCfg.Params == nil {
		shareCfg.Params = map[string]string{}
	}
	shareCfg.Params[ParamDatashare] = "true"
	share, err := Open(ctx, shareCfg)
	if err != nil {
		t.Fatalf("Open %q: %v", shared, err)
	}
	t.Cleanup(func() { _ = share.Close() })
	if err := share.Ping(ctx); err != nil {
		t.Fatalf("Ping %q: %v", shared, err)
	}

	database, err := share.SyncDatabase(ctx)
	if err != nil {
		t.Fatalf("SyncDatabase %q: %v", shared, err)
	}
	if database.Name != shared {
		t.Errorf("Name = %q, want %q", database.Name, shared)
	}
	if len(database.Schemas) == 0 {
		t.Fatal("Schemas is empty, want the ones SVV_ALL_SCHEMAS reports")
	}
	var tables int
	for _, schema := range database.Schemas {
		for _, table := range schema.Tables {
			tables++
			if len(table.Columns) == 0 {
				t.Errorf("table %s.%s has no columns, want the ones SVV_ALL_COLUMNS reports", schema.Name, table.Name)
			}
			for _, column := range table.Columns {
				if column.Type == "" {
					t.Errorf("column %s.%s.%s has no type", schema.Name, table.Name, column.Name)
				}
			}
		}
	}
	if tables == 0 {
		t.Error("no tables in the shared database, want the ones the producer published")
	}

	// A share is read-only, so the write path refuses before it sends anything.
	if _, err := share.Execute(ctx, "CREATE TABLE omni_datashare_probe (a int)", driver.ExecuteOptions{}); err == nil {
		t.Error("Execute accepted a script against a datashare database")
	}
}

// TestQuery runs statements through a session: a query that returns rows, a
// command that reports affected rows, the limit rewrite, and the plan.
func TestQuery(t *testing.T) {
	ctx := context.Background()
	d := openLive(t)

	session, err := d.Session(ctx)
	if err != nil {
		t.Fatalf("Session: %v", err)
	}
	defer session.Close()

	// A query: one row, one INT8 column.
	results, err := session.Query(ctx, "SELECT 1 AS n", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	result := results[0]
	if result.Err != nil {
		t.Fatalf("result error: %v", result.Err)
	}
	// Redshift reads an integer literal as INT4, which the session still scans
	// and reports as a 64-bit value.
	if len(result.Columns) != 1 || result.Columns[0].Name != "n" || result.Columns[0].TypeName != "INT4" {
		t.Errorf("Columns = %+v, want one INT4 column n", result.Columns)
	}
	if len(result.Rows) != 1 || len(result.Rows[0]) != 1 || result.Rows[0][0].Kind != driver.Int64 || result.Rows[0][0].Int64 != 1 {
		t.Errorf("Rows = %+v, want one row holding int64 1", result.Rows)
	}

	// The limit rewrite caps a union.
	limited, err := session.Query(ctx, "SELECT 1 AS n UNION ALL SELECT 2", driver.QueryOptions{Limit: 1})
	if err != nil {
		t.Fatalf("Query with Limit: %v", err)
	}
	if len(limited) != 1 || len(limited[0].Rows) != 1 {
		t.Errorf("limited results = %+v, want one result holding one row", limited)
	}

	// A command: the result reports affected rows and no rows.
	created, err := session.Query(ctx, "CREATE TEMP TABLE omni_query_probe (a int)", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query with a command: %v", err)
	}
	if len(created) != 1 || created[0].AffectedRows == nil {
		t.Errorf("command results = %+v, want one result reporting affected rows", created)
	}

	// A plan: text, in the rows, and marked as a plan.
	planned, err := session.Query(ctx, "SELECT COUNT(*) FROM omni_query_probe", driver.QueryOptions{Explain: true})
	if err != nil {
		t.Fatalf("Query with Explain: %v", err)
	}
	if len(planned) != 1 || planned[0].Plan == nil || planned[0].Plan.Format != "TEXT" {
		t.Errorf("explained results = %+v, want one result carrying a TEXT plan", planned)
	}
	if strings.HasPrefix(planned[0].Statement, "EXPLAIN EXPLAIN") {
		t.Errorf("Statement = %q, want one EXPLAIN", planned[0].Statement)
	}
}

// TestExecute runs a script in one transaction and then one that fails, and
// checks the failure leaves nothing behind.
func TestExecute(t *testing.T) {
	ctx := context.Background()
	d := openLive(t)

	if _, err := d.Execute(ctx, "CREATE TEMP TABLE omni_execute_probe (a int)", driver.ExecuteOptions{}); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	affected, err := d.Execute(ctx, "INSERT INTO omni_execute_probe VALUES (1), (2)", driver.ExecuteOptions{})
	if err != nil {
		t.Fatalf("Execute insert: %v", err)
	}
	if affected != 2 {
		t.Errorf("affected rows = %d, want 2", affected)
	}

	// A failing statement rolls the transaction back: neither row is inserted.
	if _, err := d.Execute(ctx, "INSERT INTO omni_execute_probe VALUES (3); INSERT INTO omni_execute_probe VALUES ('not an int')", driver.ExecuteOptions{}); err == nil {
		t.Error("Execute accepted a script whose second statement fails")
	}
	session, err := d.Session(ctx)
	if err != nil {
		t.Fatalf("Session: %v", err)
	}
	defer session.Close()
	results, err := session.Query(ctx, "SELECT COUNT(*) FROM omni_execute_probe", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(results) != 1 || len(results[0].Rows) != 1 || results[0].Rows[0][0].Int64 != 2 {
		t.Errorf("count = %+v, want 2 after the rollback", results)
	}
}

// TestExecuteCreateDatabase runs CREATE DATABASE through the creation step,
// twice. The second run must find the database the first created and create
// nothing, so a task that retries is not refused by the server.
func TestExecuteCreateDatabase(t *testing.T) {
	ctx := context.Background()
	d := openLive(t)

	created := "omni_create_probe"
	// A fixed name keeps a recording's statements the same in every run; the
	// cleanup below removes the database whichever way the test ends.
	dropCreated := func() {
		databases, err := d.getDatabases(context.Background())
		if err != nil {
			t.Errorf("getDatabases: %v", err)
			return
		}
		for _, database := range databases {
			if database.Name != created {
				continue
			}
			// Redshift refuses DROP DATABASE inside a transaction block.
			if _, err := d.Execute(context.Background(), "DROP DATABASE "+created, driver.ExecuteOptions{Transaction: driver.AutoCommit}); err != nil {
				t.Errorf("DROP DATABASE %s: %v", created, err)
			}
			return
		}
	}
	t.Cleanup(dropCreated)
	// A run a previous recording left behind must not turn the first creation
	// into the second run's no-op.
	dropCreated()

	// listed counts how many times the server names the database, which is what
	// separates the second run's no-op from a second creation.
	listed := func() int {
		t.Helper()
		databases, err := d.getDatabases(ctx)
		if err != nil {
			t.Fatalf("getDatabases: %v", err)
		}
		var found int
		for _, database := range databases {
			if database.Name == created {
				found++
			}
		}
		return found
	}

	if _, err := d.Execute(ctx, "CREATE DATABASE "+created, driver.ExecuteOptions{CreateDatabase: true}); err != nil {
		t.Fatalf("Execute CREATE DATABASE: %v", err)
	}
	if found := listed(); found != 1 {
		t.Fatalf("the server lists %q %d times after the first run, want the one database the step created", created, found)
	}
	if affected, err := d.Execute(ctx, "CREATE DATABASE "+created, driver.ExecuteOptions{CreateDatabase: true}); err != nil {
		t.Fatalf("Execute CREATE DATABASE a second time: %v (the step must not create a database the server already lists)", err)
	} else if affected != 0 {
		t.Errorf("affected rows = %d, want the creation step's 0", affected)
	}
	if found := listed(); found != 1 {
		t.Errorf("the server lists %q %d times after the second run, want once", created, found)
	}
}

// TestCountAffectedRows reads the planner's estimate for a statement that
// modifies a table it has statistics for.
func TestCountAffectedRows(t *testing.T) {
	ctx := context.Background()
	d := openLive(t)

	// Redshift has no generate_series, so a recursive CTE makes the rows the
	// planner measures. The probe is a permanent table rather than a temporary
	// one: the estimate is read on another connection of the pool, which would
	// not see a temporary table the script made.
	t.Cleanup(func() {
		if _, err := d.Execute(context.Background(), "DROP TABLE IF EXISTS public.omni_stats_probe", driver.ExecuteOptions{Transaction: driver.AutoCommit}); err != nil {
			t.Errorf("DROP TABLE public.omni_stats_probe: %v", err)
		}
	})
	if _, err := d.Execute(ctx, `
		DROP TABLE IF EXISTS public.omni_stats_probe;
		CREATE TABLE public.omni_stats_probe (a int);
		INSERT INTO public.omni_stats_probe
			WITH RECURSIVE numbers(n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM numbers WHERE n < 100)
			SELECT 1 FROM numbers;
		ANALYZE public.omni_stats_probe;
	`, driver.ExecuteOptions{}); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	estimate, err := d.CountAffectedRows(ctx, "DELETE FROM public.omni_stats_probe WHERE a = 1")
	if err != nil {
		t.Fatalf("CountAffectedRows: %v", err)
	}
	if estimate <= 0 {
		t.Errorf("estimate = %d, want the plan's row estimate", estimate)
	}
}

func contains(names []string, name string) bool {
	for _, candidate := range names {
		if candidate == name {
			return true
		}
	}
	return false
}
