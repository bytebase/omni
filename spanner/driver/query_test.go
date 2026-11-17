package driver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	// The header and the first row already encode to more than the cap, so the
	// read stops on the row that crossed it.
	a.Len(results[0].Rows, 1)
}

// TestQueryMaxResultBytesAtTheCut cuts a read at the row that crosses the cap,
// wherever in the result it falls: the header and Singers' three rows encode to
// 138 + 106 + 41 + 48 = 333 bytes, so the third row is the one that crosses a cap
// of 332, and a cap of 333 is met exactly and truncates nothing.
func TestQueryMaxResultBytesAtTheCut(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "query")
	defer d.Close()
	const statement = `SELECT Id, FirstName, Active, Rating, Photo, Born, UpdatedAt, Meta, Scores FROM Singers ORDER BY Id`

	cut, err := sessionOf(t, d).Query(context.Background(), statement, driver.QueryOptions{MaxResultBytes: 332})
	a.NoError(err)
	a.True(cut[0].Truncated)
	a.Len(cut[0].Rows, 3)
	a.Equal(int64(3), cut[0].Rows[2][0].Int64)

	// Two rows encode to 285 bytes and leave no room for the third.
	two, err := sessionOf(t, d).Query(context.Background(), statement, driver.QueryOptions{MaxResultBytes: 284})
	a.NoError(err)
	a.True(two[0].Truncated)
	a.Len(two[0].Rows, 2)

	exact, err := sessionOf(t, d).Query(context.Background(), statement, driver.QueryOptions{MaxResultBytes: 333})
	a.NoError(err)
	a.False(exact[0].Truncated)
	a.Len(exact[0].Rows, 3)
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
	// The failure states the server's gRPC code, which the package doc names:
	// Spanner answers an unknown table as an invalid argument.
	st, ok := status.FromError(results[1].Err)
	a.True(ok)
	a.Equal(codes.InvalidArgument, st.Code())
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

// TestQueryExplainStopsAtError ends an explaining script at the statement that
// failed, as a query script does: the plan of a statement after the failure is
// not asked for.
func TestQueryExplainStopsAtError(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "query")
	defer d.Close()

	results, err := sessionOf(t, d).Query(context.Background(),
		`SELECT * FROM Missing; SELECT 1 AS a`, driver.QueryOptions{Explain: true, MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Len(results, 1)
	a.Error(results[0].Err)
	a.Contains(results[0].Err.Error(), "Table not found")
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

// TestQueryGrant runs Spanner's privilege statements, which the admin API takes
// and the DML path refuses: GRANT and REVOKE are DDL. The move routed them as
// DML, as Bytebase's prefix test did, and every one of them failed; this holds
// them where the admin API runs them.
func TestQueryGrant(t *testing.T) {
	a := require.New(t)
	testcontainer.Create(t, "querygrant", "querygrant", fixtureDDL, fixtureRows)
	d := openDriver(t, "querygrant", "querygrant")
	defer d.Close()

	results, err := sessionOf(t, d).Query(context.Background(),
		`CREATE ROLE reviewer;
		 GRANT SELECT ON TABLE Singers TO ROLE reviewer;
		 REVOKE SELECT ON TABLE Singers FROM ROLE reviewer`,
		driver.QueryOptions{MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Len(results, 3)
	for _, result := range results {
		a.NoError(result.Err, result.Statement)
		a.Empty(result.Columns, result.Statement)
		a.Nil(result.AffectedRows, result.Statement)
	}

	// Execute takes the same path, where it used to be sent as DML.
	count, err := d.Execute(context.Background(),
		`GRANT SELECT ON TABLE Singers TO ROLE reviewer`, driver.ExecuteOptions{})
	a.NoError(err)
	a.Equal(int64(0), count)
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

// TestQueryPostgres reads every type a PostgreSQL database returns and holds the
// cells to the contract's kinds. The types arrive as the differential API's type
// codes — bigint as INT64, text and varchar as STRING, numeric as NUMERIC
// annotated PG_NUMERIC, jsonb as JSON annotated PG_JSONB, bytea as BYTES,
// timestamptz as TIMESTAMP, and an array as ARRAY — so the decoding is the one
// GoogleSQL uses.
func TestQueryPostgres(t *testing.T) {
	a := require.New(t)
	d := openPostgresFixture(t, "pgquery")
	defer d.Close()

	statement := `SELECT id, first_name, active, rating, photo, born, updated_at, meta, scores, wallet FROM Singers ORDER BY id`
	results, err := sessionOf(t, d).Query(context.Background(), statement, driver.QueryOptions{MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Len(results, 1)

	result := results[0]
	a.Equal(statement, result.Statement)
	a.NoError(result.Err)
	a.Equal([]driver.Column{
		{Name: "id", TypeName: "INT64"},
		{Name: "first_name", TypeName: "STRING"},
		{Name: "active", TypeName: "BOOL"},
		{Name: "rating", TypeName: "FLOAT64"},
		{Name: "photo", TypeName: "BYTES"},
		{Name: "born", TypeName: "DATE"},
		{Name: "updated_at", TypeName: "TIMESTAMP"},
		{Name: "meta", TypeName: "JSON"},
		{Name: "scores", TypeName: "[]INT64"},
		{Name: "wallet", TypeName: "NUMERIC"},
	}, result.Columns)

	a.Len(result.Rows, 3)
	a.Equal([]driver.Value{
		{Kind: driver.Int64, Int64: 1},
		{Kind: driver.String, String: "Ann"},
		{Kind: driver.Bool, Bool: true},
		{Kind: driver.Float64, Float64: 4.5},
		{Kind: driver.Bytes, Bytes: []byte{0x01, 0x02, 0x03}},
		{Kind: driver.String, String: "2024-01-15"},
		timestampValue(t, "2024-01-15T10:30:45.123456Z", 6),
		// PostgreSQL states a jsonb value as its own text, spacing included.
		{Kind: driver.String, String: `{"a": 1, "b": [2, 3]}`},
		{Kind: driver.String, String: `["1","2","3"]`},
		{Kind: driver.String, String: "12.34"},
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
		{Kind: driver.Null},
	}, result.Rows[1], "row 2")
	a.Equal(int64(3), result.Rows[2][0].Int64)
}

// TestQueryPostgresDollarQuote splits a PostgreSQL script with PostgreSQL's
// scanner: a dollar-quoted string keeps the semicolon it holds, where the
// GoogleSQL splitter would send `SELECT $$a` and `b$$ AS x` as two statements.
func TestQueryPostgresDollarQuote(t *testing.T) {
	a := require.New(t)
	d := openPostgresFixture(t, "pgquery")
	defer d.Close()

	results, err := sessionOf(t, d).Query(context.Background(),
		`SELECT $$a;b$$ AS x; SELECT 2 AS y`, driver.QueryOptions{MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Len(results, 2)

	a.NoError(results[0].Err)
	a.Equal(`SELECT $$a;b$$ AS x;`, results[0].Statement)
	a.Equal([]driver.Value{{Kind: driver.String, String: "a;b"}}, results[0].Rows[0])

	a.NoError(results[1].Err)
	a.Equal(`SELECT 2 AS y`, results[1].Statement)
	a.Equal([]driver.Value{{Kind: driver.Int64, Int64: 2}}, results[1].Rows[0])
}

// TestQueryPostgresLimit lowers a PostgreSQL query's literal LIMIT and leaves a
// LIMIT the rewriter cannot lower as written, where the read caps the rows.
func TestQueryPostgresLimit(t *testing.T) {
	a := require.New(t)
	d := openPostgresFixture(t, "pgquery")
	defer d.Close()
	session := sessionOf(t, d)

	results, err := session.Query(context.Background(),
		`SELECT id FROM Singers ORDER BY id LIMIT 5000`, driver.QueryOptions{Limit: 2, MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Equal(`SELECT id FROM Singers ORDER BY id LIMIT 2`, results[0].Statement)
	a.Len(results[0].Rows, 2)

	// A CAST is not a literal the rewriter lowers; the read stops at the cap.
	results, err = session.Query(context.Background(),
		`SELECT id FROM Singers ORDER BY id LIMIT CAST(3 AS bigint)`, driver.QueryOptions{Limit: 2, MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Equal(`SELECT id FROM Singers ORDER BY id LIMIT CAST(3 AS bigint)`, results[0].Statement)
	a.Len(results[0].Rows, 2)
}

// sessionOf opens a session on d.
func sessionOf(t *testing.T, d *Driver) driver.Session {
	t.Helper()
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })
	return session
}
