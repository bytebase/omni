package driver

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/bytebase/omni/cosmosdb/internal/testcontainer"
	"github.com/bytebase/omni/driver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// gatewayPlan is the plan the gateway states for a statement, as much of it as
// a test reads: the version of the document, and the query info the client's
// query engine runs a cross-partition query from.
type gatewayPlan struct {
	Version int `json:"partitionedQueryExecutionInfoVersion"`
	Query   struct {
		OrderBy        []string `json:"orderBy"`
		OrderByExps    []string `json:"orderByExpressions"`
		Aggregates     []string `json:"aggregates"`
		RewrittenQuery string   `json:"rewrittenQuery"`
	} `json:"queryInfo"`
}

// explainPlan runs statement as an explain and returns the plan it answered
// with, after asserting the result is one plan row and nothing else.
func explainPlan(t *testing.T, session driver.Session, statement string, opts driver.QueryOptions) gatewayPlan {
	t.Helper()
	opts.Schema = testcontainer.CitiesContainer
	opts.Explain = true
	results, err := session.Query(context.Background(), statement, opts)
	require.NoError(t, err)
	require.Len(t, results, 1)
	result := results[0]
	require.NoError(t, result.Err)

	require.NotNil(t, result.Plan, "an explain result carries its plan")
	assert.Equal(t, "JSON", result.Plan.Format)
	assert.False(t, result.Plan.Executed, "the gateway plans a statement without running it")
	assert.Equal(t, statement, result.Statement)
	assert.Equal(t, []driver.Column{{Name: "result", TypeName: "TEXT"}}, result.Columns)
	require.Len(t, result.Rows, 1, "a plan is one row; the statement's own rows are not returned")
	require.Len(t, result.Rows[0], 1)
	require.Equal(t, driver.String, result.Rows[0][0].Kind)

	var plan gatewayPlan
	require.NoError(t, json.Unmarshal([]byte(result.Rows[0][0].String), &plan), "plan: %s", result.Rows[0][0].String)
	assert.Equal(t, 2, plan.Version)
	return plan
}

// TestQueryExplain pins what an explain returns: the plan the gateway states
// for the statement, which is the plan the client's query engine runs a
// cross-partition query from, instead of the statement's rows.
func TestQueryExplain(t *testing.T) {
	d := newDriver(t)
	session, err := d.Session(context.Background())
	require.NoError(t, err)

	t.Run("gateway_served", func(t *testing.T) {
		plan := explainPlan(t, session, `SELECT * FROM c WHERE c.country = "US"`, driver.QueryOptions{})
		assert.Empty(t, plan.Query.OrderBy)
		assert.Empty(t, plan.Query.RewrittenQuery)
	})

	t.Run("cross_partition", func(t *testing.T) {
		// The plan of a query only the query engine can serve states the
		// rewrite the engine reads each range with, which is what makes the
		// plan an answer rather than an echo.
		plan := explainPlan(t, session, `SELECT c.seq FROM c ORDER BY c.seq`, driver.QueryOptions{})
		assert.Equal(t, []string{"Ascending"}, plan.Query.OrderBy)
		assert.Equal(t, []string{"c.seq"}, plan.Query.OrderByExps)
		assert.Contains(t, plan.Query.RewrittenQuery, "ORDER BY c.seq")
	})

	t.Run("aggregate", func(t *testing.T) {
		plan := explainPlan(t, session, `SELECT VALUE COUNT(1) FROM c`, driver.QueryOptions{})
		assert.Equal(t, []string{"Count"}, plan.Query.Aggregates)
	})

	t.Run("a_script_plans_each_statement", func(t *testing.T) {
		results, err := session.Query(context.Background(),
			`SELECT * FROM c WHERE c.country = "US"; SELECT c.seq FROM c ORDER BY c.seq`,
			driver.QueryOptions{Schema: testcontainer.CitiesContainer, Explain: true})
		require.NoError(t, err)
		require.Len(t, results, 2, "each statement keeps its own plan")
		for i, result := range results {
			require.NoError(t, result.Err, "statement %d", i)
			require.NotNil(t, result.Plan, "statement %d", i)
			require.Len(t, result.Rows, 1, "statement %d", i)
		}
		var crossPartition gatewayPlan
		require.NoError(t, json.Unmarshal([]byte(results[1].Rows[0][0].String), &crossPartition))
		assert.Equal(t, []string{"Ascending"}, crossPartition.Query.OrderBy)
	})

	t.Run("format_names_no_second_plan", func(t *testing.T) {
		// The gateway states one plan format, so a caller that names another
		// gets that plan, and its Result.Plan says what it is.
		plan := explainPlan(t, session, `SELECT c.seq FROM c ORDER BY c.seq`, driver.QueryOptions{ExplainFormat: "TEXT"})
		assert.Equal(t, []string{"Ascending"}, plan.Query.OrderBy)
	})
}

// TestQueryExplainReportsTheGatewayError pins an explain of a container that
// does not exist: the failure is the statement's, as a query's failure is, and
// no plan or rows come back.
func TestQueryExplainReportsTheGatewayError(t *testing.T) {
	d := newDriver(t)
	session, err := d.Session(context.Background())
	require.NoError(t, err)

	results, err := session.Query(context.Background(), `SELECT * FROM c`, driver.QueryOptions{
		Schema:  "no-such-container",
		Explain: true,
	})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Error(t, results[0].Err)
	assert.Contains(t, results[0].Err.Error(), "failed to read the query plan")
	var responseErr *azcore.ResponseError
	require.ErrorAs(t, results[0].Err, &responseErr)
	assert.Equal(t, 404, responseErr.StatusCode)
	assert.Nil(t, results[0].Plan)
	assert.Empty(t, results[0].Rows)
}

// TestQueryLimitAboveThePageCap pins a limit past the page the gateway serves
// on its own: the page size asked for is that page, so the rows arrive in pages
// the gateway fills. The tape holds the page size the driver asked for, so a
// driver that asked for nothing cannot replay it.
func TestQueryLimitAboveThePageCap(t *testing.T) {
	d := newDriver(t)
	session, err := d.Session(context.Background())
	require.NoError(t, err)

	results, err := session.Query(context.Background(), `SELECT c.seq FROM c`, driver.QueryOptions{
		Schema: testcontainer.PagedContainer,
		Limit:  gatewayPageSize + 200,
	})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	assert.Len(t, results[0].Rows, gatewayPageSize+200, "the caller's limit is the rows read")
}
