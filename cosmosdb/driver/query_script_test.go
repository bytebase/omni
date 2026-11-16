package driver

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/bytebase/omni/cosmosdb/internal/testcontainer"
	"github.com/bytebase/omni/driver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A script reaches the driver as one call and leaves as one Result per
// statement: the query language has no statement terminator, and the gateway
// rejects text holding two statements, so a script that is sent whole loses
// every statement after the first.

// TestQueryRunsEachStatementOfAScript pins that each statement of a script runs
// on its own and keeps its own result, with the terminator left out of the text
// the statement reports.
func TestQueryRunsEachStatementOfAScript(t *testing.T) {
	d := newDriver(t)
	session, err := d.Session(context.Background())
	require.NoError(t, err)

	script := `SELECT c.name FROM c WHERE c.country = "US"; SELECT c.name FROM c WHERE c.country = "AE"`
	results, err := session.Query(context.Background(), script, driver.QueryOptions{
		Schema: testcontainer.CitiesContainer,
	})
	require.NoError(t, err)
	require.Len(t, results, 2)

	// A segment keeps the layout around its statement, so only the statement's
	// own text and its terminator are decided here.
	for i, country := range []string{"US", "AE"} {
		result := results[i]
		require.NoError(t, result.Err)
		assert.Equal(t, strings.TrimSpace(`SELECT c.name FROM c WHERE c.country = "`+country+`"`), strings.TrimSpace(result.Statement))
		require.Len(t, result.Rows, testcontainer.CitiesPerCountry)
		for _, row := range result.Rows {
			require.Len(t, row, 1)
			require.Equal(t, driver.String, row[0].Kind)
			var city struct {
				Name string `json:"name"`
			}
			require.NoError(t, json.Unmarshal([]byte(row[0].String), &city))
			assert.Contains(t, city.Name, "-"+country+"-", "row %s is not a %s city", row[0].String, country)
		}
	}
}

// TestQueryScriptStopsAtAFailingStatement pins the contract's failure rule: the
// failing statement's Result carries the error and the statements after it do
// not run, so the tape holds no request for them.
func TestQueryScriptStopsAtAFailingStatement(t *testing.T) {
	d := newDriver(t)
	session, err := d.Session(context.Background())
	require.NoError(t, err)

	script := `SELECT c.name FROM c WHERE c.country = "US"; SELECT FROM WHERE; SELECT c.name FROM c WHERE c.country = "AE"`
	results, err := session.Query(context.Background(), script, driver.QueryOptions{
		Schema: testcontainer.CitiesContainer,
	})
	require.NoError(t, err)
	require.Len(t, results, 2, "the third statement does not run")

	require.NoError(t, results[0].Err)
	assert.Len(t, results[0].Rows, testcontainer.CitiesPerCountry)

	require.Error(t, results[1].Err)
	assert.Contains(t, results[1].Err.Error(), "support simple SELECT statement for Cosmos DB engine only")
	assert.Equal(t, ` SELECT FROM WHERE`, results[1].Statement)
	assert.Empty(t, results[1].Rows)
}

// TestQueryScriptWithoutAStatement covers a script holding no statement, whose
// segments are only whitespace, terminators and comments: there is nothing to
// run, so there is nothing to return.
func TestQueryScriptWithoutAStatement(t *testing.T) {
	d := newDriver(t)
	session, err := d.Session(context.Background())
	require.NoError(t, err)

	results, err := session.Query(context.Background(), " ; -- nothing to run\n ;", driver.QueryOptions{
		Schema: testcontainer.CitiesContainer,
	})
	require.NoError(t, err)
	assert.Empty(t, results)
}
