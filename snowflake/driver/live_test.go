package driver

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/snowflakedb/gosnowflake/v2"

	"github.com/bytebase/omni/driver"
)

// The live tests reach a real account and are the only tests in this package
// that do. They skip with the variable's name when neither is set, because no CI
// machine has a Snowflake account: SNOWFLAKE_DSN names a scratch account or
// trial instance of its own, and SNOWFLAKE_COMPAT_DSN the account omni's
// compatibility work asks for, which the recording in VERIFY.md may point both
// at.
//
// Nothing here has ever passed. It is the shape the day-of recording runs, and
// VERIFY.md says what to do when one fails.

// dsnVariables are the environment variables a live test accepts, in the order
// it reads them.
var dsnVariables = []string{"SNOWFLAKE_DSN", "SNOWFLAKE_COMPAT_DSN"}

// liveConfig reads the account's connection settings out of the first DSN
// variable that is set, and skips the test when none is. The DSN is written in
// gosnowflake's own form, user[:password]@account[/database][?warehouse=w],
// which is the form a Snowflake user already has.
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
		t.Skipf("no %s to reach an account: the live tests need a scratch Snowflake account (a trial instance works), and skip rather than reach a stand-in", strings.Join(dsnVariables, " or "))
	}

	config, err := gosnowflake.ParseDSN(dsn)
	if err != nil {
		t.Fatalf("%s does not parse as a Snowflake DSN: %v", variable, err)
	}
	// The client states the account and the region apart, and derives the host
	// from them; a DSN that names an endpoint of its own is written in the
	// account@host form instead (VERIFY.md).
	host := config.Account
	if host == "" {
		host = config.Host
	}
	if host == "" {
		t.Fatalf("%s names no account", variable)
	}
	cfg := driver.Config{
		Host:     host,
		User:     config.User,
		Password: config.Password,
		Database: config.Database,
	}
	if config.Port != 0 {
		cfg.Port = strconv.Itoa(config.Port)
	}
	params := map[string]string{}
	if config.Region != "" {
		params[ParamRegion] = config.Region
	}
	if config.Warehouse != "" {
		params[ParamWarehouse] = config.Warehouse
	}
	if config.Role != "" {
		params[ParamRole] = config.Role
	}
	if len(params) > 0 {
		cfg.Params = params
	}
	return cfg
}

// openLive returns an open driver for the account the environment names.
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

// TestLiveSyncInstance reads the account's version, databases and users, and
// pins the two things every other live test stands on: the version reads, and
// the database the driver connects to is listed.
func TestLiveSyncInstance(t *testing.T) {
	d := openLive(t)

	instance, err := d.SyncInstance(context.Background())
	if err != nil {
		t.Fatalf("SyncInstance: %v (VERIFY.md: the login needs a role that can read SNOWFLAKE.ACCOUNT_USAGE and SHOW SHARES)", err)
	}
	if instance.Version == "" {
		t.Error("Version is empty, want CURRENT_VERSION()'s answer")
	}
	if len(instance.Roles) == 0 {
		t.Error("Roles is empty, want the account's users")
	}
	var found bool
	for _, database := range instance.Databases {
		if database.Name == d.databaseName {
			found = true
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
	if len(database.Schemas) == 0 {
		t.Fatal("Schemas is empty, want at least PUBLIC")
	}
	var foundPublic bool
	for _, schema := range database.Schemas {
		if schema.Name == "PUBLIC" {
			foundPublic = true
		}
	}
	if !foundPublic {
		t.Error("Schemas does not hold PUBLIC")
	}
}

// TestLiveQuery runs statements through a session: a query that returns rows, a
// command that reports affected rows, the limit rewrite and the plan.
func TestLiveQuery(t *testing.T) {
	ctx := context.Background()
	d := openLive(t)

	session, err := d.Session(ctx)
	if err != nil {
		t.Fatalf("Session: %v", err)
	}
	defer session.Close()

	// A query: one row, one fixed-point column.
	results, err := session.Query(ctx, "SELECT 1 AS N", driver.QueryOptions{})
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
	if len(result.Columns) != 1 || result.Columns[0].Name != "N" || result.Columns[0].TypeName != "FIXED" {
		t.Errorf("Columns = %+v, want one FIXED column N", result.Columns)
	}
	if len(result.Rows) != 1 || len(result.Rows[0]) != 1 || result.Rows[0][0].Int64 != 1 {
		t.Errorf("Rows = %+v, want one row holding int64 1", result.Rows)
	}

	// The limit rewrite caps a union.
	limited, err := session.Query(ctx, "SELECT 1 AS N UNION ALL SELECT 2", driver.QueryOptions{Limit: 1})
	if err != nil {
		t.Fatalf("Query with Limit: %v", err)
	}
	if len(limited) != 1 || len(limited[0].Rows) != 1 {
		t.Errorf("limited results = %+v, want one result holding one row", limited)
	}

	// A command: the result reports affected rows and no rows.
	created, err := session.Query(ctx, "CREATE TEMPORARY TABLE omni_query_probe (A INT)", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query with a command: %v", err)
	}
	if len(created) != 1 || created[0].AffectedRows == nil {
		t.Errorf("command results = %+v, want one result reporting affected rows", created)
	}

	// A plan: rows of operations, marked as a plan.
	planned, err := session.Query(ctx, "SELECT * FROM omni_query_probe", driver.QueryOptions{Explain: true})
	if err != nil {
		t.Fatalf("Query with Explain: %v", err)
	}
	if len(planned) != 1 || planned[0].Plan == nil || planned[0].Plan.Format != "text" {
		t.Errorf("explained results = %+v, want one result carrying a text plan", planned)
	}
	if strings.HasPrefix(planned[0].Statement, "EXPLAIN EXPLAIN") {
		t.Errorf("Statement = %q, want one EXPLAIN", planned[0].Statement)
	}

	// A failing statement's error is the server's own.
	failed, err := session.Query(ctx, "SELECT * FROM omni_no_such_table", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query with a failing statement: %v", err)
	}
	if len(failed) != 1 || failed[0].Err == nil {
		t.Fatalf("failing results = %+v, want one result carrying an error", failed)
	}
	var serverError *gosnowflake.SnowflakeError
	if !errors.As(failed[0].Err, &serverError) {
		t.Errorf("error = %T, want the server's *gosnowflake.SnowflakeError", failed[0].Err)
	}
}

// TestLiveExecute runs a script in one transaction and then one that fails, and
// checks the failure leaves nothing behind.
func TestLiveExecute(t *testing.T) {
	ctx := context.Background()
	d := openLive(t)

	if _, err := d.Execute(ctx, "CREATE TEMPORARY TABLE omni_execute_probe (A INT)", driver.ExecuteOptions{}); err != nil {
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
	results, err := session.Query(ctx, "SELECT COUNT(*) AS C FROM omni_execute_probe", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(results) != 1 || len(results[0].Rows) != 1 || results[0].Rows[0][0].Int64 != 2 {
		t.Errorf("count = %+v, want 2 after the rollback", results)
	}
}

// TestLiveDump writes the database's own definition and checks the statements a
// dump must not write: the CREATE DATABASE it opens with, and any CREATE OR
// REPLACE that would replace an object the target holds.
func TestLiveDump(t *testing.T) {
	d := openLive(t)

	var out strings.Builder
	if err := d.Dump(context.Background(), &out, nil); err != nil {
		t.Fatalf("Dump: %v", err)
	}
	definition := out.String()
	if definition == "" {
		t.Fatal("Dump wrote nothing")
	}
	for _, unwanted := range []string{"create or replace", "CREATE OR REPLACE"} {
		if strings.Contains(definition, unwanted) {
			t.Errorf("the dump still holds %q:\n%s", unwanted, definition)
		}
	}
	if strings.Contains(strings.ToLower(definition), "database "+strings.ToLower(d.databaseName)) {
		t.Errorf("the dump still creates the database:\n%s", definition)
	}
}
