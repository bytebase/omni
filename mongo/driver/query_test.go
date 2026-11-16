package driver

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/bytebase/omni/driver"
)

// query runs script on d's session and returns the results.
func query(t *testing.T, d *Driver, script string, opts driver.QueryOptions) []*driver.Result {
	t.Helper()
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	defer func() { require.NoError(t, session.Close()) }()
	results, err := session.Query(context.Background(), script, opts)
	require.NoError(t, err)
	return results
}

// TestQuery reads documents through bracket notation, whose two quote styles
// are what the editor sends for a collection name a dot would split.
func TestQuery(t *testing.T) {
	const collectionName = "test.collection"
	d := openDriver(t, testDatabase)
	_, err := d.Execute(context.Background(), `db['`+collectionName+`'].insertMany([
		{"name": "Alice", "age": 30},
		{"name": "Bob", "age": 25},
		{"name": "Charlie", "age": 35}
	]);`, driver.ExecuteOptions{})
	require.NoError(t, err)

	for _, statement := range []string{
		`db["` + collectionName + `"].find().limit(50)`,
		`db['` + collectionName + `'].find().limit(50)`,
	} {
		t.Run(statement, func(t *testing.T) {
			results := query(t, d, statement, driver.QueryOptions{})
			require.Len(t, results, 1)
			result := results[0]
			require.NoError(t, result.Err)
			require.Equal(t, statement, result.Statement)
			require.Equal(t, []driver.Column{{Name: "result", TypeName: "TEXT"}}, result.Columns)
			require.Len(t, result.Rows, 3)
			for _, row := range result.Rows {
				require.Len(t, row, 1)
				require.Equal(t, driver.String, row[0].Kind)
				require.Contains(t, row[0].String, `"name"`)
				require.Contains(t, row[0].String, `"age"`)
			}
			var names []string
			for _, row := range result.Rows {
				for _, name := range []string{"Alice", "Bob", "Charlie"} {
					if strings.Contains(row[0].String, name) {
						names = append(names, name)
					}
				}
			}
			require.ElementsMatch(t, []string{"Alice", "Bob", "Charlie"}, names)
		})
	}
}

// TestQueryLimit holds QueryOptions.Limit to the rows a find reads.
func TestQueryLimit(t *testing.T) {
	d := openDriver(t, testDatabase)
	_, err := d.Execute(context.Background(), `db.users.insertMany([
		{_id: 1}, {_id: 2}, {_id: 3}, {_id: 4}, {_id: 5}
	]);`, driver.ExecuteOptions{})
	require.NoError(t, err)

	results := query(t, d, "db.users.find()", driver.QueryOptions{Limit: 2})
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 2)

	// A statement that caps its own rows keeps the smaller cap.
	results = query(t, d, "db.users.find().limit(1)", driver.QueryOptions{Limit: 5})
	require.Len(t, results, 1)
	require.Len(t, results[0].Rows, 1)
}

// TestQueryLimitDoesNotCapACount holds a row limit to the rows it caps: the
// client applies the cap to a count itself, which would change what the count
// counted.
func TestQueryLimitDoesNotCapACount(t *testing.T) {
	d := openDriver(t, testDatabase)
	_, err := d.Execute(context.Background(), `db.users.insertMany([
		{_id: 1}, {_id: 2}, {_id: 3}, {_id: 4}, {_id: 5}
	]);`, driver.ExecuteOptions{})
	require.NoError(t, err)

	for _, statement := range []string{"db.users.countDocuments({})", "db.users.estimatedDocumentCount()"} {
		results := query(t, d, statement, driver.QueryOptions{Limit: 2})
		require.Len(t, results, 1, statement)
		require.NoError(t, results[0].Err)
		require.Len(t, results[0].Rows, 1)
		require.Equal(t, `{"value":5}`, results[0].Rows[0][0].String, statement)
	}
}

// TestQueryMaxResultBytes holds the row cap a caller sets: rows read up to the
// limit, Truncated reporting the cut, and no cut when the caller sets none.
func TestQueryMaxResultBytes(t *testing.T) {
	d := openDriver(t, testDatabase)
	_, err := d.Execute(context.Background(), `db.users.insertMany([{_id: 1}, {_id: 2}, {_id: 3}]);`, driver.ExecuteOptions{})
	require.NoError(t, err)

	all := query(t, d, "db.users.find()", driver.QueryOptions{})
	require.Len(t, all, 1)
	require.Len(t, all[0].Rows, 3)
	require.False(t, all[0].Truncated)

	capped := query(t, d, "db.users.find()", driver.QueryOptions{MaxResultBytes: 1})
	require.Len(t, capped, 1)
	require.NoError(t, capped[0].Err)
	require.True(t, capped[0].Truncated)
	require.Len(t, capped[0].Rows, 1)

	// A cap the rows fit under cuts nothing.
	roomy := query(t, d, "db.users.find()", driver.QueryOptions{MaxResultBytes: 1 << 20})
	require.Len(t, roomy, 1)
	require.False(t, roomy[0].Truncated)
	require.Len(t, roomy[0].Rows, 3)
}

// TestQueryDistinctRendersEveryValueAsAJSONObject holds what a distinct value
// says: the frontend reads each row as a JSON object, and Extended JSON has no
// top-level scalar, so int32, double and null values are wrapped rather than
// rendered as their Go spelling.
func TestQueryDistinctRendersEveryValueAsAJSONObject(t *testing.T) {
	d := openDriver(t, testDatabase)
	_, err := d.Execute(context.Background(), `db.readings.insertMany([
		{_id: 1, value: 30},
		{_id: 2, value: 30},
		{_id: 3, value: 1.5},
		{_id: 4, value: null}
	]);`, driver.ExecuteOptions{})
	require.NoError(t, err)

	results := query(t, d, `db.readings.distinct("value")`, driver.QueryOptions{})
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	var rows []string
	for _, row := range results[0].Rows {
		rows = append(rows, row[0].String)
	}
	require.ElementsMatch(t, []string{`{"value":30}`, `{"value":1.5}`, `{"value":null}`}, rows)
}

// TestQueryDoubleNotation pins that a BSON double holding a large integer
// renders in mongosh notation (plain decimal), not the Go driver's scientific
// notation, and that a small one keeps its exponent.
func TestQueryDoubleNotation(t *testing.T) {
	d := openDriver(t, testDatabase)
	_, err := d.Execute(context.Background(),
		`db.users.insertOne({_id: 1, created_at: 1779696815227});`, driver.ExecuteOptions{})
	require.NoError(t, err)

	results := query(t, d,
		`db.users.aggregate([{$project: {_id: 0, big: {$divide: ["$created_at", 1]}, small: {$divide: [0.0000001, 1]}}}])`,
		driver.QueryOptions{})
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 1)
	require.Equal(t, `{"big":1779696815227,"small":1e-7}`, results[0].Rows[0][0].String)
}

// TestQueryError holds a statement's failure to its own Result: the script ends
// there, the statements before it keep their results, and the call returns no
// error.
func TestQueryError(t *testing.T) {
	d := openDriver(t, testDatabase)
	_, err := d.Execute(context.Background(), `db.users.insertOne({_id: 1});`, driver.ExecuteOptions{})
	require.NoError(t, err)

	results := query(t, d,
		"db.users.find();\ndb.runCommand({nosuchcommand: 1});\ndb.users.find();", driver.QueryOptions{})
	require.Len(t, results, 2)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 1)
	require.Error(t, results[1].Err)
	var commandError mongo.CommandError
	require.True(t, errors.As(results[1].Err, &commandError), "the error keeps the client's type")
	require.Equal(t, int32(59), commandError.Code)
	require.Equal(t, "db.runCommand({nosuchcommand: 1})", results[1].Statement)

	// Text the mongo parser cannot split at all is the call's error, since
	// there is no statement to run.
	results = query(t, d, "this is not mongosh", driver.QueryOptions{})
	require.Len(t, results, 1)
	require.ErrorContains(t, results[0].Err, "syntax error")

	// An empty script runs nothing.
	require.Empty(t, query(t, d, " \n ", driver.QueryOptions{}))
}

// TestQueryExplain holds what a plan request runs and what a result says about
// a plan the caller asked for: MongoDB asks for a plan with a cursor method
// rather than a keyword, and only the queryPlanner verbosity describes the
// query without running it.
func TestQueryExplain(t *testing.T) {
	d := openDriver(t, testDatabase)
	_, err := d.Execute(context.Background(), `db.orders.insertMany([
		{"customer_id": 1, "total": 100}, {"customer_id": 2, "total": 900}
	]);`, driver.ExecuteOptions{})
	require.NoError(t, err)

	explain := func(statement string, asExplain bool) *driver.Result {
		t.Helper()
		results := query(t, d, statement, driver.QueryOptions{Explain: asExplain, ExplainFormat: "JSON"})
		require.Len(t, results, 1)
		require.NoError(t, results[0].Err)
		return results[0]
	}
	planner := func(result *driver.Result) map[string]any {
		t.Helper()
		require.Len(t, result.Rows, 1)
		var plan map[string]any
		require.NoError(t, json.Unmarshal([]byte(result.Rows[0][0].String), &plan))
		queryPlanner, ok := plan["queryPlanner"].(map[string]any)
		require.True(t, ok, "no queryPlanner in %v", plan)
		return queryPlanner
	}

	const statement = `db.orders.find({total: {$gt: 500}})`
	plan := explain(statement, true)
	require.Equal(t, statement+`.explain("queryPlanner")`, plan.Statement)
	require.Contains(t, planner(plan), "winningPlan")

	// A statement that already asks for a plan keeps its query and loses its
	// verbosity, so an explain request never runs what it plans.
	replanned := explain(statement+`.explain("executionStats")`, true)
	require.Equal(t, statement+`.explain("queryPlanner")`, replanned.Statement)
	require.NotContains(t, planner(replanned), "executionStats")

	// The client runs a whole script in one call, but a plan covers one
	// statement, so a script is planned a statement at a time rather than
	// refused for holding two.
	const second = `db.orders.find({total: {$lt: 500}})`
	scripted := query(t, d, statement+"; "+second, driver.QueryOptions{Explain: true, ExplainFormat: "JSON"})
	require.Len(t, scripted, 2)
	require.Equal(t, statement+`.explain("queryPlanner")`, scripted[0].Statement)
	require.Equal(t, second+`.explain("queryPlanner")`, scripted[1].Statement)
	require.Contains(t, planner(scripted[1]), "winningPlan")

	// An explain request answers with a plan, which the result says.
	require.Equal(t, &driver.Plan{Format: "JSON"}, plan.Plan)
	require.Equal(t, &driver.Plan{Format: "JSON"}, replanned.Plan)

	// Run on a statement that explains itself reports the plan it returned.
	ran := explain(statement+`.explain()`, false)
	require.Equal(t, &driver.Plan{Format: "JSON"}, ran.Plan)

	// A measured plan runs the query, which the result has to say.
	measured := explain(statement+`.explain("executionStats")`, false)
	require.Equal(t, &driver.Plan{Format: "JSON", Executed: true}, measured.Plan)

	// An ordinary query is no plan at all.
	require.Nil(t, explain(statement, false).Plan)

	// A write has no plan, and planning it would mean running it.
	results := query(t, d, `db.orders.insertOne({a: 1})`, driver.QueryOptions{Explain: true, ExplainFormat: "JSON"})
	require.Len(t, results, 1)
	require.ErrorContains(t, results[0].Err, "insertOne")
}
