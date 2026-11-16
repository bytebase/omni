package driver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/spanner/internal/testcontainer"
)

// maxResultBytes is no cap, for the reads that assert every row.
const maxResultBytes int64 = 1 << 30

func TestQuery(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "query")
	defer d.Close()

	results, err := sessionOf(t, d).Query(context.Background(),
		`SELECT Id, FirstName, Active, Rating, Photo, Born, UpdatedAt, Meta, Scores FROM Singers ORDER BY Id`,
		driver.QueryOptions{MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Len(results, 1)

	result := results[0]
	a.Equal(`SELECT Id, FirstName, Active, Rating, Photo, Born, UpdatedAt, Meta, Scores FROM Singers ORDER BY Id`, result.Statement)
	a.NoError(result.Err)
	a.Equal([]driver.Column{
		{Name: "Id", TypeName: "INT64"},
		{Name: "FirstName", TypeName: "STRING"},
		{Name: "Active", TypeName: "BOOL"},
		{Name: "Rating", TypeName: "FLOAT64"},
		{Name: "Photo", TypeName: "BYTES"},
		{Name: "Born", TypeName: "DATE"},
		{Name: "UpdatedAt", TypeName: "TIMESTAMP"},
		{Name: "Meta", TypeName: "JSON"},
		{Name: "Scores", TypeName: "[]INT64"},
	}, result.Columns)
	a.Nil(result.AffectedRows)
	a.False(result.Truncated)
	a.Positive(result.Latency)

	// Every column of the first row carries its own kind; the second row leaves
	// the nullable ones out.
	a.Len(result.Rows, 3)
	a.Equal([]driver.Value{
		{Kind: driver.Int64, Int64: 1},
		{Kind: driver.String, String: "Ann"},
		{Kind: driver.Bool, Bool: true},
		{Kind: driver.Float64, Float64: 4.5},
		{Kind: driver.Bytes, Bytes: []byte{0x01, 0x02, 0x03}},
		{Kind: driver.String, String: "2024-01-15"},
		timestampValue(t, "2024-01-15T10:30:45.123456Z", 6),
		{Kind: driver.String, String: `{"a":1,"b":[2,3]}`},
		{Kind: driver.String, String: `["1","2","3"]`},
	}, result.Rows[0], "row 1")
	a.Equal([]driver.Value{
		{Kind: driver.Int64, Int64: 2},
		{Kind: driver.String, String: "Bob"},
		{Kind: driver.Null},
		{Kind: driver.Null},
		{Kind: driver.Null},
		{Kind: driver.Null},
		{Kind: driver.Null},
		{Kind: driver.Null},
		{Kind: driver.Null},
	}, result.Rows[1], "row 2")
	a.Equal(int64(3), result.Rows[2][0].Int64)
}

// timestampValue is the cell a TIMESTAMP reads as, built from the text the
// server wrote so the assertion does not depend on a zone's spelling.
func timestampValue(t *testing.T, text string, precision int32) driver.Value {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339Nano, text)
	require.NoError(t, err)
	return driver.Value{Kind: driver.Timestamp, Time: parsed, Precision: precision}
}

// TestQueryEmpty reads a query that answers no row: the columns it declares are
// still reported, which is what a caller renders a result table from.
func TestQueryEmpty(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "query")
	defer d.Close()

	results, err := sessionOf(t, d).Query(context.Background(),
		`SELECT Id, FirstName FROM Singers WHERE Id < 0`, driver.QueryOptions{MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Len(results, 1)
	a.NoError(results[0].Err)
	a.Equal([]driver.Column{{Name: "Id", TypeName: "INT64"}, {Name: "FirstName", TypeName: "STRING"}}, results[0].Columns)
	a.Empty(results[0].Rows)
	a.Nil(results[0].AffectedRows)
}

// TestQueryLimit lowers a query's LIMIT to the caller's, and caps the rows it
// reads when the statement's own limit is not a literal the rewrite can lower.
func TestQueryLimit(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "query")
	defer d.Close()
	session := sessionOf(t, d)

	results, err := session.Query(context.Background(),
		`SELECT Id FROM Singers ORDER BY Id`, driver.QueryOptions{Limit: 2, MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Equal(`SELECT Id FROM Singers ORDER BY Id LIMIT 2`, results[0].Statement)
	a.Len(results[0].Rows, 2)

	// The rewrite keeps a non-literal LIMIT, so the driver caps the rows.
	results, err = session.Query(context.Background(),
		`SELECT Id FROM Singers ORDER BY Id LIMIT CAST(3 AS INT64)`, driver.QueryOptions{Limit: 2, MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Equal(`SELECT Id FROM Singers ORDER BY Id LIMIT CAST(3 AS INT64)`, results[0].Statement)
	a.Len(results[0].Rows, 2)
}

// TestQueryTruncated stops a read once the rows kept pass the byte cap.
func TestQueryTruncated(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "query")
	defer d.Close()

	results, err := sessionOf(t, d).Query(context.Background(),
		`SELECT Id, FirstName FROM Singers ORDER BY Id`, driver.QueryOptions{MaxResultBytes: 40})
	a.NoError(err)
	a.True(results[0].Truncated)
	// The cap is measured at every doubling of the row count, so the first row
	// is measured and the read stops there.
	a.Len(results[0].Rows, 1)
}

// TestQueryStopsAtError ends a script at the statement that failed, and reports
// the server's own error on that statement's result.
func TestQueryStopsAtError(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "query")
	defer d.Close()

	results, err := sessionOf(t, d).Query(context.Background(),
		`SELECT 1 AS a; SELECT * FROM Missing; SELECT 2 AS b`, driver.QueryOptions{MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Len(results, 2)
	a.NoError(results[0].Err)
	a.Contains(results[1].Err.Error(), "Table not found")
	a.Equal(`SELECT * FROM Missing`, results[1].Statement)
	a.Empty(results[1].Rows)
}

// TestQueryExplain records what the emulator answers for a plan: it has no
// optimizer, so it answers no plan. The driver reports the server's error on the
// statement's result, as it does for any statement that failed.
func TestQueryExplain(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "query")
	defer d.Close()

	results, err := sessionOf(t, d).Query(context.Background(),
		`SELECT * FROM Singers`, driver.QueryOptions{Explain: true, MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Len(results, 1)
	a.Error(results[0].Err)
	a.Contains(results[0].Err.Error(), "query plan unavailable")
}

// TestQueryStructRefused refuses a result that declares a Spanner STRUCT, which
// the contract has no cell kind for.
func TestQueryStructRefused(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "query")
	defer d.Close()

	results, err := sessionOf(t, d).Query(context.Background(),
		`SELECT ARRAY(SELECT AS STRUCT 1 AS a)`, driver.QueryOptions{MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Len(results, 1)
	a.EqualError(results[0].Err, "spanner STRUCT type is not supported")
}

// TestQueryDMLAndDDLInAScript runs a script of its own so the fixture above stays
// as it was recorded: DML reports the rows it affected, DDL applies through the
// admin API and reports no rows, and a query reads what the script wrote.
func TestQueryDMLAndDDLInAScript(t *testing.T) {
	a := require.New(t)
	testcontainer.Create(t, "queryscript", "queryscript", fixtureDDL, fixtureRows)
	d := openDriver(t, "queryscript", "queryscript")
	defer d.Close()

	results, err := sessionOf(t, d).Query(context.Background(),
		`UPDATE Singers SET Rating = 5.0 WHERE Id = 1;
		 CREATE TABLE Written (Id INT64 NOT NULL) PRIMARY KEY (Id);
		 INSERT INTO Written (Id) VALUES (1), (2);
		 SELECT Id FROM Written ORDER BY Id`,
		driver.QueryOptions{MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Len(results, 4)

	a.NoError(results[0].Err)
	a.Equal(int64(1), *results[0].AffectedRows)
	a.Empty(results[0].Columns)

	a.NoError(results[1].Err)
	a.Nil(results[1].AffectedRows)
	a.Empty(results[1].Columns)

	a.NoError(results[2].Err)
	a.Equal(int64(2), *results[2].AffectedRows)

	a.NoError(results[3].Err)
	a.Equal([]driver.Column{{Name: "Id", TypeName: "INT64"}}, results[3].Columns)
	a.Len(results[3].Rows, 2)
}

// sessionOf opens a session on d.
func sessionOf(t *testing.T, d *Driver) driver.Session {
	t.Helper()
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })
	return session
}
