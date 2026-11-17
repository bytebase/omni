package driver

import (
	"context"
	"crypto/tls"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/bytebase/omni/driver"
)

// The live tests reach a real cluster and are the only tests in this package
// that do. They skip with the variable's name when neither is set, because no
// CI machine has a Redshift endpoint: REDSHIFT_DSN names a scratch cluster of
// its own, and REDSHIFT_COMPAT_DSN the cluster omni's compat suite already asks,
// which the recording in VERIFY.md may point both at.
//
// Nothing here has ever passed. It is the shape the day-of recording runs, and
// VERIFY.md says what to do when one fails.

// dsnVariables are the environment variables a live test accepts, in the order
// it reads them.
var dsnVariables = []string{"REDSHIFT_DSN", "REDSHIFT_COMPAT_DSN"}

// liveConfig reads the cluster's connection settings out of the first DSN
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
		t.Skipf("no %s to reach a cluster: the live tests need a scratch Redshift cluster, and skip rather than reach a PostgreSQL stand-in", strings.Join(dsnVariables, " or "))
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

// openLive returns an open driver for the cluster the environment names.
func openLive(t *testing.T) *Driver {
	t.Helper()
	d, err := Open(context.Background(), liveConfig(t))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	if err := d.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
	return d
}

// TestLiveSyncInstance reads the cluster's version, databases and roles, and
// pins the two things every other live test stands on: the version parses, and
// the database the driver connects to is listed.
func TestLiveSyncInstance(t *testing.T) {
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

// TestLiveSyncDatabase reads the schema of the database the driver connects to:
// the queries a sync sends, and the metadata they produce.
func TestLiveSyncDatabase(t *testing.T) {
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
}

// TestLiveQuery runs statements through a session: a query that returns rows, a
// command that reports affected rows, the limit rewrite, and the plan.
func TestLiveQuery(t *testing.T) {
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
	if len(result.Columns) != 1 || result.Columns[0].Name != "n" || result.Columns[0].TypeName != "INT8" {
		t.Errorf("Columns = %+v, want one INT8 column n", result.Columns)
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

// TestLiveExecute runs a script in one transaction and then one that fails, and
// checks the failure leaves nothing behind.
func TestLiveExecute(t *testing.T) {
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

// TestLiveCountAffectedRows reads the planner's estimate for a statement that
// modifies a table it has statistics for.
func TestLiveCountAffectedRows(t *testing.T) {
	ctx := context.Background()
	d := openLive(t)

	if _, err := d.Execute(ctx, `
		CREATE TEMP TABLE omni_stats_probe (a int);
		INSERT INTO omni_stats_probe SELECT 1 FROM generate_series(1, 100);
		ANALYZE omni_stats_probe;
	`, driver.ExecuteOptions{}); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	estimate, err := d.CountAffectedRows(ctx, "DELETE FROM omni_stats_probe WHERE a = 1")
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
