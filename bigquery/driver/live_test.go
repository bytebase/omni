package driver

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/bytebase/omni/driver"
)

// The live tests ask a real BigQuery project and are the only tests in this
// package that do. They skip with the gate's name when it is unset, because no CI
// machine has a project: BIGQUERY_CREDENTIALS holds a service-account key, and the
// project and dataset come from BIGQUERY_PROJECT and BIGQUERY_DATASET.
//
// Nothing here has ever passed. It is the shape the day-of recording runs, and
// VERIFY.md says what to do when one fails.

// liveVariables are the settings a live test needs, in the order the skip message
// names them. BIGQUERY_CREDENTIALS gates all of them: without a key there is
// nothing to ask and nothing to record.
const (
	credentialVariable = "BIGQUERY_CREDENTIALS"
	projectVariable    = "BIGQUERY_PROJECT"
	datasetVariable    = "BIGQUERY_DATASET"
)

// liveConfig reads the recording's connection settings out of the environment and
// skips the test when the gate is unset.
func liveConfig(t *testing.T) driver.Config {
	t.Helper()
	credentials := strings.TrimSpace(os.Getenv(credentialVariable))
	if credentials == "" {
		t.Skipf("no %s: the live tests need a Google Cloud project, a dataset and a service-account key, and skip rather than ask a stand-in", credentialVariable)
	}
	if !strings.HasPrefix(credentials, "{") {
		key, err := os.ReadFile(credentials)
		if err != nil {
			t.Fatalf("%s names a file that cannot be read: %v", credentialVariable, err)
		}
		credentials = string(key)
	}

	project := strings.TrimSpace(os.Getenv(projectVariable))
	if project == "" {
		// A service-account key names its project, which is the one the recording
		// bills.
		var key struct {
			ProjectID string `json:"project_id"`
		}
		if err := json.Unmarshal([]byte(credentials), &key); err != nil {
			t.Fatalf("%s does not parse as a credential: %v", credentialVariable, err)
		}
		project = key.ProjectID
	}
	if project == "" {
		t.Fatalf("%s is not set and the key names no project_id", projectVariable)
	}
	dataset := strings.TrimSpace(os.Getenv(datasetVariable))
	if dataset == "" {
		t.Fatalf("%s is not set: a recording runs against a scratch dataset of its own", datasetVariable)
	}
	return driver.Config{
		Database: dataset,
		Password: credentials,
		Params:   map[string]string{ParamProject: project},
	}
}

// openLive returns an open driver for the project the environment names.
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

// TestLiveSyncInstance lists the project's datasets: the sync's listing, the order
// it reports, and that the dataset the driver is open on is among them.
func TestLiveSyncInstance(t *testing.T) {
	ctx := context.Background()
	d := openLive(t)

	instance, err := d.SyncInstance(ctx)
	if err != nil {
		t.Fatalf("SyncInstance: %v", err)
	}
	if instance.Version != "" {
		t.Errorf("Version = %q, want none: BigQuery reports no version", instance.Version)
	}
	if len(instance.Roles) != 0 {
		t.Errorf("Roles = %v, want none: BigQuery states no role to a sync", instance.Roles)
	}
	if len(instance.Databases) == 0 {
		t.Fatal("Databases is empty, want the project's datasets")
	}
	var found bool
	for i, database := range instance.Databases {
		if i > 0 && instance.Databases[i-1].Name > database.Name {
			t.Errorf("Databases = %q before %q, want name order", instance.Databases[i-1].Name, database.Name)
		}
		if database.Name == d.databaseName {
			found = true
		}
	}
	if !found {
		t.Errorf("Databases does not hold %q", d.databaseName)
	}
}

// TestLiveSyncDatabase reads the schema of the dataset the driver is open on. The
// table it expects is one the recording's corpus creates, so a replay of this test
// is a replay of that corpus.
func TestLiveSyncDatabase(t *testing.T) {
	ctx := context.Background()
	d := openLive(t)

	database, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatalf("SyncDatabase: %v", err)
	}
	if database.Name != d.databaseName {
		t.Errorf("Name = %q, want %q", database.Name, d.databaseName)
	}
	if len(database.Schemas) != 1 {
		t.Fatalf("Schemas = %d, want one unnamed schema", len(database.Schemas))
	}
	schema := database.Schemas[0]
	if schema.Name != "" {
		t.Errorf("Schemas[0].Name = %q, want the dataset's one unnamed schema", schema.Name)
	}
	var found bool
	for i, table := range schema.Tables {
		if i > 0 && schema.Tables[i-1].Name > table.Name {
			t.Errorf("Tables = %q before %q, want name order", schema.Tables[i-1].Name, table.Name)
		}
		if table.Name != "customers" {
			continue
		}
		found = true
		if table.Comment != "a synced comment" {
			t.Errorf("customers comment = %q, want the corpus's", table.Comment)
		}
		if table.RowCount <= 0 {
			t.Errorf("customers row count = %d, want the corpus's rows", table.RowCount)
		}
		var city *string
		for _, column := range table.Columns {
			if column.Name == "city" {
				value := column.Type
				city = &value
			}
		}
		if city == nil || *city != "STRING" {
			t.Errorf("customers.city = %v, want a STRING column", city)
		}
	}
	if !found {
		t.Error("Tables does not hold the corpus's customers table")
	}
}

// missingDataset is the dataset name the not-found test asks for. It is fixed
// rather than made up from the clock because a tape replays byte for byte, and
// the name is part of the INFORMATION_SCHEMA query the recording holds.
const missingDataset = "omni_driver_missing_dataset"

// TestLiveSyncDatabaseMissing opens a driver on a dataset the project does not
// have and requires the contract's not-found, which is what a caller uses to tell
// a database it may drop from a request that failed.
func TestLiveSyncDatabaseMissing(t *testing.T) {
	ctx := context.Background()
	cfg := liveConfig(t)
	cfg.Database = missingDataset
	d, err := Open(ctx, cfg)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = d.Close() }()

	if _, err := d.SyncDatabase(ctx); !errors.Is(err, driver.ErrNotFound) {
		t.Fatalf("SyncDatabase error = %v, want driver.ErrNotFound", err)
	}
}

// TestLiveQuery runs statements through a session: a query that returns rows, the
// limit wrap, a command, the dry run that stands in for EXPLAIN, the result-size
// cap, and a statement the service refuses.
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
	if len(result.Columns) != 1 || result.Columns[0].Name != "n" || result.Columns[0].TypeName != "INTEGER" {
		t.Errorf("Columns = %+v, want one INTEGER column n", result.Columns)
	}
	if len(result.Rows) != 1 || len(result.Rows[0]) != 1 || result.Rows[0][0].Kind != driver.Int64 || result.Rows[0][0].Int64 != 1 {
		t.Errorf("Rows = %+v, want one row holding int64 1", result.Rows)
	}

	limited, err := session.Query(ctx, "SELECT 1 AS n UNION ALL SELECT 2; SELECT 3 AS n UNION ALL SELECT 4", driver.QueryOptions{Limit: 1})
	if err != nil {
		t.Fatalf("Query with Limit: %v", err)
	}
	if len(limited) != 2 {
		t.Fatalf("got %d results, want 2", len(limited))
	}
	for _, result := range limited {
		if result.Err != nil {
			t.Fatalf("limited result error: %v", result.Err)
		}
		if len(result.Rows) != 1 {
			t.Errorf("rows = %d, want the wrap to cap one", len(result.Rows))
		}
		if !strings.HasPrefix(result.Statement, "WITH result AS (") {
			t.Errorf("Statement = %q, want the CTE wrap", result.Statement)
		}
	}

	created, err := session.Query(ctx, "CREATE TABLE omni_query_probe (a INT64)", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query with a command: %v", err)
	}
	t.Cleanup(func() {
		cleanup, cleanupErr := d.Session(context.Background())
		if cleanupErr != nil {
			return
		}
		defer func() { _ = cleanup.Close() }()
		_, _ = cleanup.Query(context.Background(), "DROP TABLE IF EXISTS omni_query_probe", driver.QueryOptions{})
	})
	if len(created) != 1 || created[0].AffectedRows == nil || *created[0].AffectedRows != 0 {
		t.Errorf("command results = %+v, want one result reporting zero rows affected", created)
	}

	dryRun, err := session.Query(ctx, "SELECT COUNT(*) FROM omni_query_probe", driver.QueryOptions{Explain: true})
	if err != nil {
		t.Fatalf("Query with Explain: %v", err)
	}
	if len(dryRun) != 1 || dryRun[0].Err != nil {
		t.Fatalf("dry run results = %+v", dryRun)
	}
	if len(dryRun[0].Columns) != 1 || dryRun[0].Columns[0].Name != "Validation Result" {
		t.Errorf("Columns = %+v, want the dry run's validation column", dryRun[0].Columns)
	}
	if len(dryRun[0].Rows) != 1 || !strings.Contains(dryRun[0].Rows[0][0].String, "Query successfully validated") {
		t.Errorf("Rows = %+v, want the dry run's answer", dryRun[0].Rows)
	}

	// One row is wider than the cap, so the read stops with a truncated result.
	capped, err := session.Query(ctx, "SELECT REPEAT('x', 4096) AS s", driver.QueryOptions{MaxResultBytes: 64})
	if err != nil {
		t.Fatalf("Query with MaxResultBytes: %v", err)
	}
	if len(capped) != 1 || !capped[0].Truncated {
		t.Errorf("capped results = %+v, want Truncated", capped)
	}

	failed, err := session.Query(ctx, "SELECT * FROM omni_query_probe_no_such_table", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query of a missing table: %v", err)
	}
	if len(failed) != 1 || failed[0].Err == nil {
		t.Errorf("results = %+v, want the service's error on the result", failed)
	}
}

// TestLiveExecute runs a script as one job, and one that the service refuses.
func TestLiveExecute(t *testing.T) {
	ctx := context.Background()
	d := openLive(t)

	if _, err := d.Execute(ctx, "CREATE TABLE omni_execute_probe (a INT64)", driver.ExecuteOptions{}); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	t.Cleanup(func() {
		_, _ = d.Execute(context.Background(), "DROP TABLE IF EXISTS omni_execute_probe", driver.ExecuteOptions{})
	})
	affected, err := d.Execute(ctx, "INSERT INTO omni_execute_probe (a) VALUES (1), (2)", driver.ExecuteOptions{})
	if err != nil {
		t.Fatalf("Execute insert: %v", err)
	}
	if affected != 2 {
		t.Errorf("affected rows = %d, want 2", affected)
	}

	// A script's statements are one job, so a failing statement leaves the ones
	// before it committed: BigQuery has no transaction to roll them back.
	if _, err := d.Execute(ctx, "INSERT INTO omni_execute_probe (a) VALUES (3); INSERT INTO omni_execute_probe (a) VALUES ('not an int')", driver.ExecuteOptions{}); err == nil {
		t.Error("Execute accepted a script whose second statement fails")
	}

	session, err := d.Session(ctx)
	if err != nil {
		t.Fatalf("Session: %v", err)
	}
	defer func() { _ = session.Close() }()
	results, err := session.Query(ctx, "SELECT COUNT(*) AS n FROM omni_execute_probe", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(results) != 1 || len(results[0].Rows) != 1 || results[0].Rows[0][0].Int64 != 3 {
		t.Errorf("count = %+v, want the two inserted rows and the one the failed script left", results)
	}
}
