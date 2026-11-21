package driver

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

// The recorded driver tests replay offline from testdata/tapes/*.tape with no
// Google Cloud credentials. With OMNI_RECORD=1 and BIGQUERY_CREDENTIALS,
// BIGQUERY_PROJECT and BIGQUERY_DATASET set, they record against the pinned
// project and dataset in bigquery/testdata/pins.tsv.

// TestSyncInstance lists the project's datasets: the sync's listing, the order
// it reports, and that the dataset the driver is open on is among them.
func TestSyncInstance(t *testing.T) {
	ctx := context.Background()
	d := openTape(t)

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

// TestSyncDatabase reads the schema of the dataset the driver is open on. The
// table it expects is one the recording's corpus creates, so a replay of this test
// is a replay of that corpus.
func TestSyncDatabase(t *testing.T) {
	ctx := context.Background()
	d := openTape(t)

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

	// The synced metadata is held to the golden the recording committed: a column's
	// position, a table's row count or a dataset's listing moving is a change to
	// what every BigQuery instance's metadata says.
	body, err := os.ReadFile(filepath.Join("testdata", "sync", d.databaseName+".json"))
	if err != nil {
		t.Fatalf("read the sync golden: %v", err)
	}
	want := &metadatapb.DatabaseSchemaMetadata{}
	if err := protojson.Unmarshal(body, want); err != nil {
		t.Fatalf("read the sync golden: %v", err)
	}
	if !proto.Equal(database, want) {
		t.Errorf("SyncDatabase = %v, want testdata/sync/%s.json", database, d.databaseName)
	}
}

// missingDataset is the dataset name the not-found test asks for. It is fixed
// rather than made up from the clock because a tape replays byte for byte, and
// the name is part of the request the recording holds.
const missingDataset = "omni_driver_missing_dataset"

// TestSyncDatabaseMissing opens a driver on a dataset the project does not have
// and requires the contract's not-found, which is what a caller uses to tell a
// database it may drop from a request that failed.
func TestSyncDatabaseMissing(t *testing.T) {
	ctx := context.Background()
	d := openTapeDataset(t, missingDataset)

	if _, err := d.SyncDatabase(ctx); !errors.Is(err, driver.ErrNotFound) {
		t.Fatalf("SyncDatabase error = %v, want driver.ErrNotFound", err)
	}
}

// TestQuery runs the statements of a session that need no table: a query that
// returns rows, every BigQuery cell kind, a statement omni's parser does not
// read, the limit wrap, the result-size cap, the dry run that stands in for
// EXPLAIN, and a statement the service refuses. Everything that needs a table is
// TestQueryWrites, so a project whose dataset is missing still reads its
// failures honestly here.
func TestQuery(t *testing.T) {
	ctx := context.Background()
	d := openTape(t)

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
	if result.Plan == nil || !result.Plan.Executed || result.Plan.Format != "JSON" || !strings.Contains(result.Plan.Content, "S00: Output") {
		t.Errorf("Plan = %+v, want executed JSON plan holding S00: Output", result.Plan)
	}

	// Pin every BigQuery cell kind directly from the service's answer.
	allTypes, err := session.Query(ctx, `SELECT
		NUMERIC '12.50' AS num,
		BIGNUMERIC '-7.25' AS bn,
		DATE '2026-01-02' AS d,
		DATETIME '2026-01-02 03:04:05.123' AS dt,
		TIME '03:04:05' AS t,
		TIMESTAMP '2026-01-02 03:04:05.120000 UTC' AS ts,
		JSON '{"a":1}' AS j,
		[1, 2] AS arr,
		STRUCT(1 AS a, 'x' AS b) AS st,
		ST_GEOGPOINT(1, 2) AS g,
		INTERVAL 1 DAY AS iv,
		RANGE<DATE> '[2026-01-01, 2026-02-01)' AS r,
		b'ab' AS bin,
		TRUE AS flag,
		1.5 AS flt`, driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query of all cell kinds: %v", err)
	}
	if len(allTypes) != 1 || allTypes[0].Err != nil || len(allTypes[0].Rows) != 1 {
		t.Fatalf("allTypes = %+v", allTypes)
	}
	row := allTypes[0].Rows[0]
	wantRow := []driver.Value{
		{Kind: driver.String, String: "12.5"},
		{Kind: driver.String, String: "-7.25"},
		{Kind: driver.String, String: "2026-01-02"},
		{Kind: driver.Timestamp, Time: time.Date(2026, 1, 2, 3, 4, 5, 123000000, time.UTC), Precision: bigqueryFractionalDigits},
		{Kind: driver.String, String: "03:04:05"},
		{Kind: driver.TimestampTZ, Time: time.Date(2026, 1, 2, 3, 4, 5, 120000000, time.UTC), Precision: bigqueryFractionalDigits},
		{Kind: driver.String, String: `{"a":1}`},
		{Kind: driver.String, String: "[1, 2]"},
		{Kind: driver.String, String: "[1, x]"},
		{Kind: driver.String, String: "POINT(1 2)"},
		{Kind: driver.String, String: "0-0 1 0:0:0"},
		{Kind: driver.String, String: "[2026-01-01, 2026-02-01)"},
		{Kind: driver.Bytes, Bytes: []byte("ab")},
		{Kind: driver.Bool, Bool: true},
		{Kind: driver.Float64, Float64: 1.5},
	}
	if len(row) != len(wantRow) {
		t.Fatalf("got %d cells, want %d: %+v", len(row), len(wantRow), row)
	}
	for i := range wantRow {
		got, want := row[i], wantRow[i]
		if got.Kind != want.Kind || got.String != want.String || got.Int64 != want.Int64 || got.Float64 != want.Float64 ||
			got.Bool != want.Bool || string(got.Bytes) != string(want.Bytes) || !got.Time.Equal(want.Time) || got.Precision != want.Precision {
			t.Errorf("cell[%d] (%s) = %+v, want %+v", i, allTypes[0].Columns[i].Name, got, want)
		}
	}

	// A pipe query is valid BigQuery and one omni's parser does not read yet, so
	// the statement is classified as neither a query nor a command: it runs as a
	// job, and the rows the service answered are what comes back.
	piped, err := session.Query(ctx, "SELECT 1 AS n |> WHERE n = 1", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query of a pipe statement: %v", err)
	}
	if len(piped) != 1 || piped[0].Err != nil {
		t.Fatalf("pipe query results = %+v", piped)
	}
	if len(piped[0].Columns) != 1 || piped[0].Columns[0].Name != "n" {
		t.Errorf("Columns = %+v, want one column n", piped[0].Columns)
	}
	if len(piped[0].Rows) != 1 || piped[0].Rows[0][0].Int64 != 1 {
		t.Errorf("Rows = %+v, want the row the pipe query returned", piped[0].Rows)
	}

	// The cap also holds for a statement the analysis cannot read, which is not
	// rewritten: the read stops at the row it asked for.
	pipedLimited, err := session.Query(ctx, "SELECT n FROM UNNEST([1, 2, 3]) AS n |> WHERE n > 1", driver.QueryOptions{Limit: 1})
	if err != nil {
		t.Fatalf("Query of a pipe statement with Limit: %v", err)
	}
	if len(pipedLimited) != 1 || pipedLimited[0].Err != nil {
		t.Fatalf("pipe query with Limit results = %+v", pipedLimited)
	}
	if len(pipedLimited[0].Rows) != 1 {
		t.Errorf("rows = %d, want the cap to stop the read at one", len(pipedLimited[0].Rows))
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

	// One row is wider than the cap, so the read stops with a truncated result.
	capped, err := session.Query(ctx, "SELECT REPEAT('x', 4096) AS s", driver.QueryOptions{MaxResultBytes: 64})
	if err != nil {
		t.Fatalf("Query with MaxResultBytes: %v", err)
	}
	if len(capped) != 1 || !capped[0].Truncated {
		t.Errorf("capped results = %+v, want Truncated", capped)
	}

	// A statement the analysis cannot read that answered no rows still reports the
	// columns it declares, which is what a command-path result must not drop.
	emptyPipe, err := session.Query(ctx, "SELECT 1 AS n |> WHERE n = 2", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query of an empty pipe statement: %v", err)
	}
	if len(emptyPipe) != 1 || emptyPipe[0].Err != nil {
		t.Fatalf("empty pipe results = %+v", emptyPipe)
	}
	if len(emptyPipe[0].Columns) != 1 || emptyPipe[0].Columns[0].Name != "n" {
		t.Errorf("Columns = %+v, want the column the empty pipe query declares", emptyPipe[0].Columns)
	}
	if len(emptyPipe[0].Rows) != 0 || emptyPipe[0].AffectedRows != nil {
		t.Errorf("rows = %d affected = %v, want no rows and no affected count", len(emptyPipe[0].Rows), emptyPipe[0].AffectedRows)
	}

	requireValidated(t, ctx, session, "SELECT COUNT(*) AS n FROM UNNEST([1, 2, 3])")

	// The service refuses the statement, and the refusal is its own error on the
	// statement's result rather than the call's.
	failed, err := session.Query(ctx, "SELEKT 1", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query of a refused statement: %v", err)
	}
	if len(failed) != 1 || failed[0].Err == nil || !strings.Contains(failed[0].Err.Error(), "Syntax error") {
		t.Errorf("results = %+v, want the service's syntax error on the result", failed)
	}
}

// TestQueryWrites runs the statements of the query path that need a table: a
// command's affected rows, a dry run of that command's table, and the error the
// service states for a table it does not have.
func TestQueryWrites(t *testing.T) {
	ctx := context.Background()
	d := openTape(t)

	session, err := d.Session(ctx)
	if err != nil {
		t.Fatalf("Session: %v", err)
	}
	defer func() { _ = session.Close() }()

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
	if len(created) != 1 {
		t.Fatalf("got %d results, want 1", len(created))
	}
	if created[0].Err != nil {
		t.Fatalf("command error: %v", created[0].Err)
	}
	if created[0].AffectedRows == nil || *created[0].AffectedRows != 0 {
		t.Errorf("command results = %+v, want one result reporting zero rows affected", created)
	}

	requireValidated(t, ctx, session, "SELECT COUNT(*) FROM omni_query_probe")

	failed, err := session.Query(ctx, "SELECT * FROM omni_query_probe_no_such_table", driver.QueryOptions{})
	if err != nil {
		t.Fatalf("Query of a missing table: %v", err)
	}
	if len(failed) != 1 || failed[0].Err == nil {
		t.Errorf("results = %+v, want the service's error on the result", failed)
	}
}

// requireValidated requires the answer BigQuery states in place of EXPLAIN: one
// STRING column named for the validation and one row saying the statement holds.
func requireValidated(t *testing.T, ctx context.Context, session driver.Session, statement string) {
	t.Helper()
	dryRun, err := session.Query(ctx, statement, driver.QueryOptions{Explain: true})
	if err != nil {
		t.Fatalf("Query with Explain: %v", err)
	}
	if len(dryRun) != 1 {
		t.Fatalf("got %d results, want 1", len(dryRun))
	}
	if dryRun[0].Err != nil {
		t.Fatalf("dry run error: %v", dryRun[0].Err)
	}
	if len(dryRun[0].Columns) != 1 || dryRun[0].Columns[0].Name != "Validation Result" {
		t.Errorf("Columns = %+v, want the dry run's validation column", dryRun[0].Columns)
	}
	if len(dryRun[0].Rows) != 1 || !strings.Contains(dryRun[0].Rows[0][0].String, "Query successfully validated") {
		t.Errorf("Rows = %+v, want the dry run's answer", dryRun[0].Rows)
	}
}

// TestExecute runs a script as one job, and one that the service refuses.
func TestExecute(t *testing.T) {
	ctx := context.Background()
	d := openTape(t)

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
