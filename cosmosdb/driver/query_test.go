package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/bytebase/omni/cosmosdb/internal/testcontainer"
	"github.com/bytebase/omni/driver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newDriver returns a driver whose client is t's tape: a replay of the pinned
// emulator's answers, or, in a recording, the emulator itself.
func newDriver(t *testing.T) *Driver {
	t.Helper()
	return &Driver{client: OpenTape(t), databaseName: testcontainer.Database}
}

func TestPing(t *testing.T) {
	require.NoError(t, newDriver(t).Ping(context.Background()))
}

func TestClose(t *testing.T) {
	require.NoError(t, newDriver(t).Close())
}

// TestSyncInstance reads the account's databases through the client's query
// API, which is how a server answers.
func TestSyncInstance(t *testing.T) {
	instance, err := newDriver(t).SyncInstance(context.Background())
	require.NoError(t, err)

	var names []string
	for _, database := range instance.Databases {
		names = append(names, database.Name)
		assert.Empty(t, database.Schemas, "an instance sync reports a database without its containers")
	}
	assert.ElementsMatch(t, []string{testcontainer.Database, testcontainer.SpareDatabase}, names)
	assert.Empty(t, instance.Version, "Cosmos DB reports no account version")
	assert.Empty(t, instance.Roles)
}

// TestSyncInstanceEmulatorListing reads the same databases over the emulator's
// REST API, which is what the driver uses when the caller names the emulator.
// The two listings are one account's, so they answer with one list.
func TestSyncInstanceEmulatorListing(t *testing.T) {
	d := &Driver{client: OpenTapeEmulator(t), databaseName: testcontainer.Database}
	instance, err := d.SyncInstance(context.Background())
	require.NoError(t, err)

	var names []string
	for _, database := range instance.Databases {
		names = append(names, database.Name)
	}
	assert.ElementsMatch(t, []string{testcontainer.Database, testcontainer.SpareDatabase}, names)
}

func TestSyncDatabase(t *testing.T) {
	database, err := newDriver(t).SyncDatabase(context.Background())
	require.NoError(t, err)

	assert.Equal(t, testcontainer.Database, database.Name)
	require.Len(t, database.Schemas, 1)
	var names []string
	for _, table := range database.Schemas[0].Tables {
		names = append(names, table.Name)
		assert.Empty(t, table.Columns, "a container's items are schemaless, so its table reports no column")
	}
	assert.ElementsMatch(t, []string{testcontainer.CitiesContainer, testcontainer.PagedContainer}, names)
}

func TestSyncDatabaseEmulatorListing(t *testing.T) {
	d := &Driver{client: OpenTapeEmulator(t), databaseName: testcontainer.Database}
	database, err := d.SyncDatabase(context.Background())
	require.NoError(t, err)

	var names []string
	for _, table := range database.Schemas[0].Tables {
		names = append(names, table.Name)
	}
	assert.ElementsMatch(t, []string{testcontainer.CitiesContainer, testcontainer.PagedContainer}, names)
}

// TestSyncDatabaseNotExist pins what a sync of a database the account does not
// have answers: driver.ErrNotFound with the account's own error behind it.
func TestSyncDatabaseNotExist(t *testing.T) {
	d := &Driver{client: OpenTape(t), databaseName: "no-such-database"}
	_, err := d.SyncDatabase(context.Background())
	require.Error(t, err)

	require.ErrorIs(t, err, driver.ErrNotFound)
	var responseErr *azcore.ResponseError
	require.ErrorAs(t, err, &responseErr)
	assert.Equal(t, 404, responseErr.StatusCode)
}

// TestSyncDatabaseNotExistEmulatorListing is the same failure over the
// emulator's REST listing, which reports the status as text rather than as a
// client error type.
func TestSyncDatabaseNotExistEmulatorListing(t *testing.T) {
	d := &Driver{client: OpenTapeEmulator(t), databaseName: "no-such-database"}
	_, err := d.SyncDatabase(context.Background())
	require.Error(t, err)
	require.ErrorIs(t, err, driver.ErrNotFound)
	assert.Contains(t, err.Error(), "HTTP 404")
}

// TestQueryMaxResultBytes covers the cap on a result's size: the rows up to it
// come back, and the result says it stopped.
func TestQueryMaxResultBytes(t *testing.T) {
	d := newDriver(t)
	session, err := d.Session(context.Background())
	require.NoError(t, err)

	results, err := session.Query(context.Background(), `SELECT c.name FROM c ORDER BY c.name`, driver.QueryOptions{
		Schema:         testcontainer.CitiesContainer,
		MaxResultBytes: 120,
	})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.True(t, results[0].Truncated)
	require.NotEmpty(t, results[0].Rows)
	require.Less(t, len(results[0].Rows), testcontainer.CitiesItems(), "the cap stops before the last row")

	// Without a cap the same query returns every row.
	all := query(t, d, testcontainer.CitiesContainer, `SELECT c.name FROM c ORDER BY c.name`)
	require.Len(t, all, testcontainer.CitiesItems())
}

func TestQueryWithoutContainer(t *testing.T) {
	d := newDriver(t)
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	_, err = session.Query(context.Background(), `SELECT * FROM c`, driver.QueryOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "container argument is required for CosmosDB")
}

// TestQueryStatementValidation pins the one statement shape the driver sends:
// text the query language parses. The failure is the statement's, not the
// call's, because the call had a statement to run.
func TestQueryStatementValidation(t *testing.T) {
	d := newDriver(t)
	session, err := d.Session(context.Background())
	require.NoError(t, err)

	results, err := session.Query(context.Background(), `SELECT FROM WHERE`, driver.QueryOptions{
		Schema: testcontainer.CitiesContainer,
	})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Error(t, results[0].Err)
	assert.Contains(t, results[0].Err.Error(), "support simple SELECT statement for Cosmos DB engine only")
	assert.Empty(t, results[0].Rows)
}

// TestQueryMissingContainer pins what a query of a container that does not
// exist answers.
func TestQueryMissingContainer(t *testing.T) {
	d := newDriver(t)
	session, err := d.Session(context.Background())
	require.NoError(t, err)

	results, err := session.Query(context.Background(), `SELECT * FROM c`, driver.QueryOptions{
		Schema: "no-such-container",
	})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Error(t, results[0].Err)
	var responseErr *azcore.ResponseError
	require.ErrorAs(t, results[0].Err, &responseErr)
	assert.Equal(t, 404, responseErr.StatusCode)
}

// runQuery runs one query in container and returns its rows as the JSON
// documents the driver emitted them as.
func runQuery(t *testing.T, session driver.Session, container, statement string, opts driver.QueryOptions) [][]byte {
	t.Helper()
	opts.Schema = container
	results, err := session.Query(context.Background(), statement, opts)
	require.NoError(t, err, "Query: %s", statement)
	require.Len(t, results, 1, "the driver returns exactly one Result per query")
	require.NoError(t, results[0].Err, "Query: %s", statement)
	require.Equal(t, statement, results[0].Statement)
	require.Equal(t, []driver.Column{{Name: "result", TypeName: "TEXT"}}, results[0].Columns)

	out := make([][]byte, 0, len(results[0].Rows))
	for _, row := range results[0].Rows {
		require.Len(t, row, 1, "each row has exactly one column of JSON text")
		require.Equal(t, driver.String, row[0].Kind)
		out = append(out, []byte(row[0].String))
	}
	return out
}

// query is runQuery on a session of t's driver.
func query(t *testing.T, d *Driver, container, statement string) [][]byte {
	t.Helper()
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	return runQuery(t, session, container, statement, driver.QueryOptions{})
}

// ------------------------------------------- cross-partition, pinned emulator
//
// The gateway cannot serve an aggregate, DISTINCT, ORDER BY, TOP,
// OFFSET/LIMIT or GROUP BY that spans partitions, so the query pivots to the
// pure-Go query engine, which merges the per-partition streams itself. Every
// test below runs on the emulator's fixture, whose contents are known, and so
// asserts exact answers.

func TestCrossPartitionAggregates(t *testing.T) {
	d := newDriver(t)
	items := testcontainer.CitiesItems()
	sum := testcontainer.CitiesPopulationSum()

	for _, tc := range []struct {
		name string
		sql  string
		want string
	}{
		{"count_value", `SELECT VALUE COUNT(1) FROM c`, fmt.Sprint(items)},
		{"count_aliased", `SELECT COUNT(1) AS total FROM c`, fmt.Sprintf(`{"total":%d}`, items)},
		{"sum_value", `SELECT VALUE SUM(c.population) FROM c`, fmt.Sprint(sum)},
		{"avg_value", `SELECT VALUE AVG(c.population) FROM c`, fmt.Sprint(sum / items)},
		{"min_max_aliased", `SELECT MIN(c.population) AS lo, MAX(c.population) AS hi FROM c`,
			fmt.Sprintf(`{"lo":%d,"hi":%d}`, testcontainer.PopulationGap, testcontainer.CitiesPerCountry*testcontainer.PopulationGap)},
		{"count_and_sum_aliased", `SELECT COUNT(1) AS total, SUM(c.population) AS pop FROM c`,
			fmt.Sprintf(`{"total":%d,"pop":%d}`, items, sum)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rows := query(t, d, testcontainer.CitiesContainer, tc.sql)
			require.Len(t, rows, 1)
			assert.JSONEq(t, tc.want, string(rows[0]))
		})
	}
}

func TestCrossPartitionDistinct(t *testing.T) {
	d := newDriver(t)

	t.Run("value_form", func(t *testing.T) {
		rows := query(t, d, testcontainer.CitiesContainer, `SELECT DISTINCT VALUE c.country FROM c`)
		got := make([]string, 0, len(rows))
		for _, row := range rows {
			var country string
			require.NoError(t, json.Unmarshal(row, &country))
			got = append(got, country)
		}
		assert.ElementsMatch(t, testcontainer.Cities, got)
	})

	t.Run("object_form", func(t *testing.T) {
		rows := query(t, d, testcontainer.CitiesContainer, `SELECT DISTINCT c.country FROM c`)
		assert.Len(t, rows, len(testcontainer.Cities))
	})
}

func TestCrossPartitionOrderBy(t *testing.T) {
	d := newDriver(t)

	t.Run("ascending", func(t *testing.T) {
		rows := query(t, d, testcontainer.CitiesContainer, `SELECT c.name, c.population FROM c ORDER BY c.population ASC`)
		require.Len(t, rows, testcontainer.CitiesItems())
		assertOrderedByPopulation(t, rows, true)
	})

	t.Run("descending", func(t *testing.T) {
		rows := query(t, d, testcontainer.CitiesContainer, `SELECT c.name, c.population FROM c ORDER BY c.population DESC`)
		require.Len(t, rows, testcontainer.CitiesItems())
		assertOrderedByPopulation(t, rows, false)
	})
}

// TestCrossPartitionOffsetLimit pins the page against the globally ordered
// result: a merge that drops or reorders a partition's stream still returns
// five rows, so only comparing the slice catches it.
func TestCrossPartitionOffsetLimit(t *testing.T) {
	d := newDriver(t)

	all := assertAscendingByName(t, query(t, d, testcontainer.CitiesContainer, `SELECT c.name FROM c ORDER BY c.name`))
	require.Len(t, all, testcontainer.CitiesItems())

	page := assertAscendingByName(t, query(t, d, testcontainer.CitiesContainer,
		`SELECT c.name FROM c ORDER BY c.name OFFSET 5 LIMIT 5`))
	assert.Equal(t, all[5:10], page)
}

func TestCrossPartitionGroupBy(t *testing.T) {
	d := newDriver(t)
	rows := query(t, d, testcontainer.CitiesContainer,
		`SELECT c.country, COUNT(1) AS cities, SUM(c.population) AS pop FROM c GROUP BY c.country`)
	require.Len(t, rows, len(testcontainer.Cities))

	wantPop := float64(testcontainer.CitiesPopulationSum() / len(testcontainer.Cities))
	seen := map[string]bool{}
	for _, row := range rows {
		var group struct {
			Country string  `json:"country"`
			Cities  float64 `json:"cities"`
			Pop     float64 `json:"pop"`
		}
		require.NoError(t, json.Unmarshal(row, &group))
		assert.False(t, seen[group.Country], "country %q grouped twice", group.Country)
		seen[group.Country] = true
		assert.Equal(t, float64(testcontainer.CitiesPerCountry), group.Cities, "country %q", group.Country)
		assert.Equal(t, wantPop, group.Pop, "country %q", group.Country)
	}
	assert.Len(t, seen, len(testcontainer.Cities))
}

// TestCrossPartitionGatewayServable covers the queries that never leave the
// gateway path, which the engine must not take over: a projection of every
// partition, and a filter the gateway can route to a single one.
func TestCrossPartitionGatewayServable(t *testing.T) {
	d := newDriver(t)

	t.Run("projection", func(t *testing.T) {
		rows := query(t, d, testcontainer.CitiesContainer, `SELECT * FROM c`)
		assert.Len(t, rows, testcontainer.CitiesItems())
	})

	t.Run("single_partition_filter", func(t *testing.T) {
		rows := query(t, d, testcontainer.CitiesContainer, `SELECT * FROM c WHERE c.country = "US"`)
		assert.Len(t, rows, testcontainer.CitiesPerCountry)
	})
}

// TestCrossPartitionQueryLimit covers the driver's own row cap, which stops
// paging mid-result and so decides how much of a cross-partition query the
// driver pulls back. The query orders its rows, because which rows a cap keeps
// is the server's choice otherwise, and then neither a run nor a tape could say
// what the answer is.
func TestCrossPartitionQueryLimit(t *testing.T) {
	d := newDriver(t)
	session, err := d.Session(context.Background())
	require.NoError(t, err)

	results, err := session.Query(context.Background(), `SELECT c.seq FROM c ORDER BY c.seq`, driver.QueryOptions{
		Schema: testcontainer.CitiesContainer,
		Limit:  7,
	})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	assert.Len(t, results[0].Rows, 7)
}
