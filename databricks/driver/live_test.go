package driver

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/driver"
)

// The live tests reach a real Databricks workspace and are the only tests in
// this package that do. They skip with the variable's name when it is unset,
// because no CI machine has a workspace and no container can stand in for one.
//
// DATABRICKS_DSN names the workspace, the personal access token and the SQL
// warehouse:
//
//	databricks://<workspace-host>?token=<personal-access-token>&warehouse_id=<warehouse-id>&catalog=<catalog>
//
// The host may also be written as an https URL. Nothing here has ever passed;
// VERIFY.md says what to do when one fails and which order to record them in.
const dsnVariable = "DATABRICKS_DSN"

// liveConfig reads the workspace connection settings out of DATABRICKS_DSN and
// skips the test when it is unset.
func liveConfig(t *testing.T) driver.Config {
	t.Helper()
	dsn := os.Getenv(dsnVariable)
	if dsn == "" {
		t.Skipf("no %s: the live tests need a Databricks workspace, a personal access token and a SQL warehouse", dsnVariable)
	}
	cfg, err := configFromDSN(dsn)
	if err != nil {
		t.Fatalf("%s does not parse: %v", dsnVariable, err)
	}
	return cfg
}

// configFromDSN parses the DSN the live tests read. It is separated from the
// environment so a parsing error names the form rather than a secret.
func configFromDSN(dsn string) (driver.Config, error) {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return driver.Config{}, err
	}
	host := parsed.Host
	if host == "" {
		return driver.Config{}, fmt.Errorf("DSN has no host")
	}
	switch parsed.Scheme {
	case "databricks", "https":
		host = "https://" + host
	case "http":
		host = "http://" + host
	default:
		return driver.Config{}, fmt.Errorf("DSN scheme %q is not databricks, https or http", parsed.Scheme)
	}
	query := parsed.Query()
	token := query.Get("token")
	if token == "" {
		return driver.Config{}, fmt.Errorf("DSN has no token")
	}
	warehouseID := query.Get("warehouse_id")
	if warehouseID == "" {
		warehouseID = query.Get("warehouseId")
	}
	if warehouseID == "" {
		return driver.Config{}, fmt.Errorf("DSN has no warehouse_id")
	}
	return driver.Config{
		Host:     host,
		Password: token,
		Database: query.Get("catalog"),
		Params:   map[string]string{ParamWarehouseID: warehouseID},
	}, nil
}

// openLive returns an open driver for the workspace the environment names.
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

// TestLiveSyncInstance reads the SQL version and the catalog names, and pins
// the catalog every other live test stands on.
func TestLiveSyncInstance(t *testing.T) {
	d := openLive(t)
	instance, err := d.SyncInstance(context.Background())
	if err != nil {
		t.Fatalf("SyncInstance: %v", err)
	}
	if instance.Version == "" {
		t.Error("Version is empty")
	}
	if len(instance.Databases) == 0 {
		t.Fatal("Databases is empty, want the workspace's catalogs")
	}
	if d.catalog != "" {
		var found bool
		for _, database := range instance.Databases {
			if database.Name == d.catalog {
				found = true
			}
		}
		if !found {
			t.Errorf("Databases does not hold the configured catalog %q", d.catalog)
		}
	}
}

// TestLiveSyncDatabase reads the schema of the configured catalog.
func TestLiveSyncDatabase(t *testing.T) {
	d := openLive(t)
	if d.catalog == "" {
		t.Fatalf("%s names no catalog: a recording's DSN must name one (?catalog=...), because SyncDatabase reads it and a skip would record nothing", dsnVariable)
	}
	database, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatalf("SyncDatabase: %v", err)
	}
	if database == nil || database.Name != d.catalog {
		t.Errorf("database = %+v, want catalog %q", database, d.catalog)
	}
}

// TestLiveQuery runs a query, the row limit, the plan and a failing statement.
func TestLiveQuery(t *testing.T) {
	ctx := context.Background()
	d := openLive(t)

	session, err := d.Session(ctx)
	if err != nil {
		t.Fatalf("Session: %v", err)
	}
	defer func() { _ = session.Close() }()

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
	if len(result.Columns) != 1 || result.Columns[0].Name != "n" || result.Columns[0].TypeName != "INT" {
		t.Errorf("Columns = %+v, want one INT column n", result.Columns)
	}
	if len(result.Rows) != 1 || len(result.Rows[0]) != 1 || result.Rows[0][0].Kind != driver.Int32 || result.Rows[0][0].Int64 != 1 {
		t.Errorf("Rows = %+v, want one row holding int32 1", result.Rows)
	}

	limited, err := session.Query(ctx, "SELECT 1 AS n UNION ALL SELECT 2", driver.QueryOptions{Limit: 1})
	if err != nil {
		t.Fatalf("Query with Limit: %v", err)
	}
	if len(limited) != 1 || len(limited[0].Rows) > 1 {
		t.Errorf("limited results = %+v, want one result holding at most one row", limited)
	}

	planned, err := session.Query(ctx, "SELECT 1 AS n", driver.QueryOptions{Explain: true})
	if err != nil {
		t.Fatalf("Query with Explain: %v", err)
	}
	if len(planned) != 1 || planned[0].Plan == nil || planned[0].Plan.Format != "TEXT" {
		t.Errorf("explained results = %+v, want one result carrying a TEXT plan", planned)
	}
	if strings.HasPrefix(planned[0].Statement, "EXPLAIN EXPLAIN") {
		t.Errorf("Statement = %q, want one EXPLAIN", planned[0].Statement)
	}

	failed, err := session.Query(ctx, "SELECT * FROM omni_driver_missing_table", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query with a failing statement returned a script error: %v", err)
	}
	if len(failed) != 1 || failed[0].Err == nil {
		t.Errorf("failing results = %+v, want one result carrying the statement error", failed)
	}
}

// TestLiveExecute creates a scratch schema and table, inserts, reads them back,
// and drops them. Databricks commits each statement itself, so the script is
// sent statement by statement.
func TestLiveExecute(t *testing.T) {
	ctx := context.Background()
	d := openLive(t)
	if d.catalog == "" {
		t.Fatalf("%s names no catalog: a recording's DSN must name one (?catalog=...), because the scratch schema is created in it and a skip would record nothing", dsnVariable)
	}

	schema := "omni_driver_" + strconv.FormatInt(time.Now().UnixNano(), 36)
	qualified := quoteBacktick(d.catalog) + "." + quoteBacktick(schema)
	cleanup := func() {
		_, _ = d.Execute(ctx, "DROP TABLE IF EXISTS "+qualified+".`probe`", driver.ExecuteOptions{})
		_, _ = d.Execute(ctx, "DROP SCHEMA IF EXISTS "+qualified, driver.ExecuteOptions{})
	}
	t.Cleanup(cleanup)

	if _, err := d.Execute(ctx, "CREATE SCHEMA "+qualified, driver.ExecuteOptions{}); err != nil {
		t.Fatalf("Execute CREATE SCHEMA: %v", err)
	}
	if _, err := d.Execute(ctx, "CREATE TABLE "+qualified+".`probe` (id INT)", driver.ExecuteOptions{}); err != nil {
		t.Fatalf("Execute CREATE TABLE: %v", err)
	}
	if _, err := d.Execute(ctx, "INSERT INTO "+qualified+".`probe` VALUES (1), (2)", driver.ExecuteOptions{}); err != nil {
		t.Fatalf("Execute INSERT: %v", err)
	}

	session, err := d.Session(ctx)
	if err != nil {
		t.Fatalf("Session: %v", err)
	}
	defer func() { _ = session.Close() }()
	results, err := session.Query(ctx, "SELECT COUNT(*) FROM "+qualified+".`probe`", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(results) != 1 || results[0].Err != nil || len(results[0].Rows) != 1 {
		t.Fatalf("count result = %+v", results)
	}
	if results[0].Rows[0][0].Kind != driver.Int64 || results[0].Rows[0][0].Int64 != 2 {
		t.Errorf("count = %+v, want int64 2", results[0].Rows[0][0])
	}
}

// TestLivePing asks the workspace for its root folder, which is the connectivity
// check this driver offers: Databricks has no SQL-level ping.
//
// It opens the driver itself rather than through openLive, so the request its
// tape holds is the ping and nothing else.
func TestLivePing(t *testing.T) {
	ctx := context.Background()
	cfg := liveConfig(t)
	if cfg.Database == "" {
		t.Fatalf("%s names no catalog: a recording's DSN must name one (?catalog=...), because the pin is the workspace and catalog the other tapes record", dsnVariable)
	}
	d, err := Open(ctx, cfg)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = d.Close() }()

	if err := d.Ping(ctx); err != nil {
		t.Fatalf("Ping: %v (VERIFY.md step 1: the host, the token and the warehouse are what Ping reaches)", err)
	}
}

// TestLiveDump asks the workspace to state the configured catalog's objects.
func TestLiveDump(t *testing.T) {
	d := openLive(t)
	if d.catalog == "" {
		t.Fatalf("%s names no catalog: a recording's DSN must name one (?catalog=...), because Dump reads it and a skip would record nothing", dsnVariable)
	}
	var out bytes.Buffer
	if err := d.Dump(context.Background(), &out, nil); err != nil {
		t.Fatalf("Dump: %v", err)
	}
	if !strings.Contains(out.String(), "CREATE CATALOG") {
		t.Errorf("dump does not name a catalog: %q", out.String())
	}
}
