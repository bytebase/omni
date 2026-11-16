package gonative

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/queryengine"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultReturnsNonNilEngine(t *testing.T) {
	e := Default()
	require.NotNil(t, e)
}

func TestDefaultAdvertisesAggregateFeatures(t *testing.T) {
	feats := Default().SupportedFeatures()
	assert.Contains(t, feats, "Aggregate")
	assert.Contains(t, feats, "NonValueAggregate")
	assert.Contains(t, feats, "MultipleAggregates")
}

func TestDefaultCreateQueryPipelineRejectsUnparseablePlan(t *testing.T) {
	_, err := Default().CreateQueryPipeline("SELECT * FROM c", "not json", "{}")
	require.Error(t, err)
}

func TestDefaultCreateQueryPipelineRejectsUnsupportedPlan(t *testing.T) {
	// A plan carrying hasNonStreamingOrderBy (vector/hybrid-search scoring)
	// is out of the engine's scope — the design spec explicitly excludes it.
	plan := `{"partitionedQueryExecutionInfoVersion":2,"queryInfo":{"hasNonStreamingOrderBy":true}}`
	pkranges := `{"PartitionKeyRanges":[{"id":"0"}]}`
	_, err := Default().CreateQueryPipeline("SELECT ...", plan, pkranges)
	require.ErrorIs(t, err, ErrUnsupportedPlanFeature)
}

func TestDefaultCreateReadManyPipelineReturnsUnsupported(t *testing.T) {
	_, err := Default().CreateReadManyPipeline(nil, "{}", "", 0, nil)
	require.ErrorIs(t, err, ErrUnsupportedPlanFeature)
}

func TestDefaultReturnsDistinctEngineInstances(t *testing.T) {
	// Each caller gets its own instance — enables per-engine caches in later stages.
	a := Default()
	b := Default()
	assert.NotSame(t, a, b)
}

// recordedPlan loads a query plan the gateway really emitted for one of the
// BYT-10397 query shapes, as the engine receives it.
func recordedPlan(t *testing.T, name string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "query_"+name+"_plan.json"))
	require.NoError(t, err)
	return string(raw)
}

// TestCreateQueryPipelineRejectsUnhandledPlanField drives the dispatch with
// plans whose composition no pipeline implements. The query has to fail rather
// than come back missing an operator, which is what every wrong row count in
// BYT-10397 was.
func TestCreateQueryPipelineRejectsUnhandledPlanField(t *testing.T) {
	const pkranges = `{"PartitionKeyRanges":[{"id":"0"},{"id":"1"}]}`
	for _, tc := range []struct {
		plan  string
		query string
		field string
	}{
		{"groupby_top", `SELECT TOP 3 c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket`, "Top"},
		{"groupby_offset_limit", `SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket OFFSET 1 LIMIT 2`, "Offset"},
	} {
		t.Run(tc.plan, func(t *testing.T) {
			_, err := Default().CreateQueryPipeline(tc.query, recordedPlan(t, tc.plan), pkranges)
			require.ErrorIs(t, err, ErrUnsupportedPlanFeature)
			assert.Contains(t, err.Error(), "queryInfo."+tc.field)
		})
	}
}

// TestCreateQueryPipelineAcceptsRecordedCompositions pins the dispatch against
// the plan shapes BYT-10397 collected: each composes two operators, and the
// pipeline the dispatch picks has to claim both of them.
func TestCreateQueryPipelineAcceptsRecordedCompositions(t *testing.T) {
	const pkranges = `{"PartitionKeyRanges":[{"id":"0"},{"id":"1"}]}`
	for _, tc := range []struct{ plan, query string }{
		{"distinct_orderby", `SELECT DISTINCT c.bucket FROM c ORDER BY c.bucket`},
		{"distinct_orderby_value", `SELECT DISTINCT VALUE c.bucket FROM c ORDER BY c.bucket`},
		{"top_orderby_asc", `SELECT TOP 5 c.seq FROM c ORDER BY c.seq ASC`},
		{"top_orderby_desc", `SELECT TOP 5 c.seq FROM c ORDER BY c.seq DESC`},
		{"distinct_value_limit", `SELECT DISTINCT VALUE c.bucket FROM c OFFSET 0 LIMIT 3`},
		{"distinct_offset_limit", `SELECT DISTINCT c.bucket FROM c OFFSET 2 LIMIT 3`},
		{"distinct_top", `SELECT DISTINCT TOP 3 c.bucket FROM c`},
		{"offset_limit", `SELECT c.seq FROM c OFFSET 2 LIMIT 3`},
		{"agg_offset_limit", `SELECT VALUE COUNT(1) FROM c OFFSET 1 LIMIT 1`},
		{"min_value_empty", `SELECT VALUE MIN(c.seq) FROM c WHERE c.seq > 99999999`},
		{"min_aliased_empty", `SELECT MIN(c.seq) AS m FROM c WHERE c.seq > 99999999`},
		{"groupby_plain", `SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket`},
	} {
		t.Run(tc.plan, func(t *testing.T) {
			pipe, err := Default().CreateQueryPipeline(tc.query, recordedPlan(t, tc.plan), pkranges)
			require.NoError(t, err)
			pipe.Close()
		})
	}
}

// A GROUP BY that the user also ordered reaches the engine as a plan with no
// order fields at all — the ordering survives only inside rewrittenQuery — so
// the guard cannot see it. The group pipeline reads the direction off the
// gateway's own per-partition streams instead.
func TestCreateQueryPipelineAcceptsGroupByWithHiddenOrderBy(t *testing.T) {
	raw := recordedPlan(t, "groupby_orderby")
	plan, err := parsePlan([]byte(raw))
	require.NoError(t, err)
	assert.Empty(t, plan.QueryInfo.OrderBy)
	assert.Empty(t, plan.QueryInfo.OrderByExpressions)
	assert.Contains(t, plan.QueryInfo.RewrittenQuery, "ORDER BY")

	pipe, err := Default().CreateQueryPipeline(
		`SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket ORDER BY c.bucket`,
		raw, `{"PartitionKeyRanges":[{"id":"0"}]}`)
	require.NoError(t, err)
	pipe.Close()
}

// TOP is an operator a scalar aggregate honours: the gateway keeps it in the
// plan and the aggregate pipeline applies the window after finalizing.
func TestCreateQueryPipelineAcceptsTopOnScalarAggregates(t *testing.T) {
	const pkranges = `{"PartitionKeyRanges":[{"id":"0"},{"id":"1"}]}`
	for _, tc := range []struct{ name, plan, query, body string }{
		{
			name:  "value",
			plan:  `{"partitionedQueryExecutionInfoVersion":2,"queryInfo":{"rewrittenQuery":"SELECT VALUE [{\"item\": COUNT(1)}] FROM c","hasSelectValue":true,"aggregates":["Count"],"groupByAliases":[],"groupByAliasToAggregateType":{},"top":1,"orderBy":[],"orderByExpressions":[],"groupByExpressions":[]}}`,
			query: `SELECT TOP 1 VALUE COUNT(1) FROM c`,
			body:  `{"Documents":[[{"item":3}]]}`,
		},
		{
			name:  "aliased",
			plan:  `{"partitionedQueryExecutionInfoVersion":2,"queryInfo":{"rewrittenQuery":"SELECT {\"n\": {\"item\": COUNT(1)}} AS payload FROM c","hasSelectValue":false,"aggregates":[],"groupByAliases":["n"],"groupByAliasToAggregateType":{"n":"Count"},"top":1,"orderBy":[],"orderByExpressions":[],"groupByExpressions":[]}}`,
			query: `SELECT TOP 1 COUNT(1) AS n FROM c`,
			body:  `{"Documents":[{"payload":{"n":{"item":3}}}]}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pipe, err := Default().CreateQueryPipeline(tc.query, tc.plan, pkranges)
			require.NoError(t, err)
			defer pipe.Close()

			rows := driveInlinePlan(t, pipe, map[string][]byte{
				"0": []byte(tc.body),
				"1": []byte(tc.body),
			})
			assert.Equal(t, 1, len(rows), "TOP 1 keeps one row of the aggregate")
		})
	}
}

// TOP 0 keeps nothing, which the pipeline applies after finalizing.
func TestTopZeroOnAScalarAggregateEmitsNoRow(t *testing.T) {
	plan := `{"partitionedQueryExecutionInfoVersion":2,"queryInfo":{"rewrittenQuery":"SELECT VALUE [{\"item\": COUNT(1)}] FROM c","hasSelectValue":true,"aggregates":["Count"],"groupByAliases":[],"groupByAliasToAggregateType":{},"top":0,"orderBy":[],"orderByExpressions":[],"groupByExpressions":[]}}`
	pipe, err := Default().CreateQueryPipeline(`SELECT TOP 0 VALUE COUNT(1) FROM c`, plan, `{"PartitionKeyRanges":[{"id":"0"}]}`)
	require.NoError(t, err)
	defer pipe.Close()

	rows := driveInlinePlan(t, pipe, map[string][]byte{"0": []byte(`{"Documents":[[{"item":3}]]}`)})
	assert.Empty(t, rows)
}

// driveInlinePlan drives a pipeline the caller built and returns every emitted
// row.
func driveInlinePlan(t *testing.T, pipe queryengine.QueryPipeline, bodies map[string][]byte) [][]byte {
	t.Helper()
	var rows [][]byte
	for i := 0; i < 50 && !pipe.IsComplete(); i++ {
		res, err := pipe.Run()
		require.NoError(t, err, "Run turn %d", i)
		if len(res.Requests) > 0 {
			results := make([]queryengine.QueryResult, 0, len(res.Requests))
			for _, req := range res.Requests {
				body, ok := bodies[req.PartitionKeyRangeID]
				require.True(t, ok, "missing body for pk-range %s", req.PartitionKeyRangeID)
				results = append(results, queryengine.NewQueryResult(req.PartitionKeyRangeID, body, ""))
			}
			require.NoError(t, pipe.ProvideData(results))
		}
		for _, item := range res.Items {
			rows = append(rows, append([]byte(nil), item...))
		}
	}
	require.True(t, pipe.IsComplete(), "pipeline did not converge")
	return rows
}

func TestRequireHandledWalksEveryPlanField(t *testing.T) {
	n := 0
	qi := planQueryInfo{RewrittenQuery: "SELECT 1", Top: &n}

	require.NoError(t, requireHandled(qi, fieldRewrittenQuery, fieldTop))

	err := requireHandled(qi, fieldRewrittenQuery)
	require.ErrorIs(t, err, ErrUnsupportedPlanFeature)
	assert.Contains(t, err.Error(), "queryInfo.Top", "TOP 0 is an instruction, not an absent field")
}

// A plan that spells a field out as empty asks for nothing, so it must not
// make the dispatch reject a shape it can serve.
func TestRequireHandledIgnoresFieldsSpeltOutEmpty(t *testing.T) {
	p, err := parsePlan([]byte(`{"queryInfo":{"aggregates":[],"groupByAliases":[],` +
		`"groupByAliasToAggregateType":{},"distinctType":"None","orderBy":[]}}`))
	require.NoError(t, err)
	require.NoError(t, requireHandled(p.QueryInfo))
}
