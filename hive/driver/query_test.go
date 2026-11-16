package driver

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/hive/internal/testcontainer"
)

// TestQuery reads one column of every kind the driver converts, a NULL, and a
// table's own columns.
func TestQuery(t *testing.T) {
	d := openTestCorpus(t, "omni_query")
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), `
		SELECT CAST(1 AS TINYINT) AS tiny, CAST(2 AS SMALLINT) AS small, CAST(3 AS INT) AS i;
		SELECT CAST(4 AS BIGINT) AS big;
		SELECT CAST(1.5 AS DOUBLE) AS dbl, CAST(2.5 AS FLOAT) AS flt;
		SELECT CAST('x' AS STRING) AS s, true AS b, NULL AS n;
		SELECT CAST('abc' AS BINARY) AS bin;
		SELECT id, amount, note FROM orders;`, driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 6)
	for _, result := range results {
		require.NoError(t, result.Err)
	}

	require.Equal(t, []driver.Column{
		{Name: "tiny", TypeName: "TINYINT"},
		{Name: "small", TypeName: "SMALLINT"},
		{Name: "i", TypeName: "INT"},
	}, results[0].Columns)
	require.Len(t, results[0].Rows, 1)
	row := results[0].Rows[0]
	require.Equal(t, driver.Int32, row[0].Kind)
	require.Equal(t, int64(1), row[0].Int64)
	require.Equal(t, int64(2), row[1].Int64)
	require.Equal(t, int64(3), row[2].Int64)

	require.Equal(t, driver.Int64, results[1].Rows[0][0].Kind)
	require.Equal(t, int64(4), results[1].Rows[0][0].Int64)

	// A float column arrives as text with twenty fractional digits, as
	// Bytebase mapped it; the follow-up fix(hive) PR makes it a DOUBLE.
	require.Equal(t, driver.String, results[2].Rows[0][0].Kind)
	require.Equal(t, "1.50000000000000000000", results[2].Rows[0][0].String)
	require.Equal(t, "2.50000000000000000000", results[2].Rows[0][1].String)

	require.Equal(t, driver.String, results[3].Rows[0][0].Kind)
	require.Equal(t, "x", results[3].Rows[0][0].String)
	require.Equal(t, driver.Bool, results[3].Rows[0][1].Kind)
	require.True(t, results[3].Rows[0][1].Bool)
	require.Equal(t, driver.Null, results[3].Rows[0][2].Kind)

	require.Equal(t, driver.Bytes, results[4].Rows[0][0].Kind)
	require.Equal(t, []byte("abc"), results[4].Rows[0][0].Bytes)

	require.Equal(t, []string{"BIGINT", "DECIMAL", "STRING"}, columnTypes(results[5]))
	require.Len(t, results[5].Rows, 1)
	row = results[5].Rows[0]
	require.Equal(t, driver.Int64, row[0].Kind)
	require.Equal(t, int64(1), row[0].Int64)
	// A DECIMAL keeps every digit it was stored with.
	require.Equal(t, "1.50", row[1].String)
	require.Equal(t, "n1", row[2].String)
}

// TestQueryExplain asks for a plan of a statement, in both formats Hive has,
// and holds a typed EXPLAIN to what the caller wrote.
func TestQueryExplain(t *testing.T) {
	d := openTestCorpus(t, "omni_query_explain")
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), "SELECT id FROM orders", driver.QueryOptions{Explain: true})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Equal(t, "EXPLAIN SELECT id FROM orders", results[0].Statement)
	require.NotNil(t, results[0].Plan)
	require.Equal(t, "TEXT", results[0].Plan.Format)
	require.False(t, results[0].Plan.Executed)
	require.Equal(t, []string{"Explain", "STRING"}, columnAndTypes(results[0]))
	require.NotEmpty(t, results[0].Rows)

	results, err = sess.Query(context.Background(), "SELECT id FROM orders", driver.QueryOptions{Explain: true, ExplainFormat: "JSON"})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Equal(t, "EXPLAIN FORMATTED SELECT id FROM orders", results[0].Statement)
	require.NotNil(t, results[0].Plan)
	require.Equal(t, "JSON", results[0].Plan.Format)
	require.Len(t, results[0].Rows, 1)
	require.True(t, strings.HasPrefix(results[0].Rows[0][0].String, `{"CBOPlan":`))

	results, err = sess.Query(context.Background(), "EXPLAIN SELECT id FROM orders", driver.QueryOptions{Explain: true})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Equal(t, "EXPLAIN SELECT id FROM orders", results[0].Statement)
	require.Equal(t, "TEXT", results[0].Plan.Format)
}

// TestQueryMaxResultBytes stops reading a statement's rows once they pass the
// cap, and says it did.
func TestQueryMaxResultBytes(t *testing.T) {
	d := openTestCorpus(t, "omni_query_size")
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), "SELECT id, body FROM events", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Len(t, results[0].Rows, 2)
	require.False(t, results[0].Truncated)

	results, err = sess.Query(context.Background(), "SELECT id, body FROM events", driver.QueryOptions{MaxResultBytes: 30})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.True(t, results[0].Truncated)
	require.Len(t, results[0].Rows, 1)
}

// TestQueryError reports Hive's answer to a statement that failed on the
// statement's Result, and ends the script there.
func TestQueryError(t *testing.T) {
	d := openTestCorpus(t, "omni_query_error")
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), `
		SELECT id FROM orders;
		SELECT * FROM no_such_table;
		SELECT id FROM orders`, driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 2)
	require.NoError(t, results[0].Err)
	require.Error(t, results[1].Err)

	var hiveErr *testcontainer.Error
	require.True(t, errors.As(results[1].Err, &hiveErr), "the error is not Hive's answer: %v", results[1].Err)
	require.Equal(t, 10001, hiveErr.Code)
	require.Contains(t, hiveErr.Message, "Table not found 'no_such_table'")
}

// TestQueryRunsADDLAsAQuery: a statement that returns no rows is answered with
// an operation that has none, and reading that operation's rows is where Hive's
// error arrives. The statement itself ran.
func TestQueryRunsADDLAsAQuery(t *testing.T) {
	d := openTestCorpus(t, "omni_query_ddl")
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), "CREATE TABLE query_ddl (id INT) STORED AS PARQUET", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, "CREATE TABLE query_ddl (id INT) STORED AS PARQUET", results[0].Statement)
	// gohive reports the operation's missing result set as Hive's TStatus text.
	require.Error(t, results[0].Err)
	require.Contains(t, results[0].Err.Error(), "TStatus({StatusCode:ERROR_STATUS")

	results, err = sess.Query(context.Background(), "SHOW TABLES", driver.QueryOptions{})
	require.NoError(t, err)
	require.NoError(t, results[0].Err)
	tables := make([]string, 0, len(results[0].Rows))
	for _, row := range results[0].Rows {
		tables = append(tables, row[0].String)
	}
	require.Equal(t, []string{"events", "ext_users", "orders", "query_ddl", "v_users"}, tables)
}

// TestQuerySelectsTheSchema opens the session on no database and reaches the
// fixture's tables through QueryOptions.Schema.
func TestQuerySelectsTheSchema(t *testing.T) {
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{"DROP DATABASE IF EXISTS omni_query_schema CASCADE", "CREATE DATABASE omni_query_schema"} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.Open(t, "")
	execCorpus(t, db, "CREATE TABLE omni_query_schema.t (id INT) STORED AS PARQUET")
	d := &Driver{db: db}
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), "SELECT count(*) AS n FROM t", driver.QueryOptions{Schema: "omni_query_schema"})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 1)
	require.Equal(t, driver.Int64, results[0].Rows[0][0].Kind)
	require.Equal(t, int64(0), results[0].Rows[0][0].Int64)
}

// TestQuerySelectsTheSchemaAgainAfterAStatementChangesIt: a USE the session ran
// moves the connection to another database, and a caller that names a database
// gets it selected again rather than the last one it asked for.
func TestQuerySelectsTheSchemaAgainAfterAStatementChangesIt(t *testing.T) {
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{
			"DROP DATABASE IF EXISTS omni_schema_a CASCADE",
			"DROP DATABASE IF EXISTS omni_schema_b CASCADE",
			"CREATE DATABASE omni_schema_a",
			"CREATE DATABASE omni_schema_b",
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.Open(t, "")
	execCorpus(t, db, `CREATE TABLE omni_schema_a.t (tag STRING) STORED AS PARQUET

ALTER TABLE omni_schema_a.t SET TBLPROPERTIES('EXTERNAL'='FALSE')

INSERT INTO omni_schema_a.t VALUES ('a')

CREATE TABLE omni_schema_b.t (tag STRING) STORED AS PARQUET

ALTER TABLE omni_schema_b.t SET TBLPROPERTIES('EXTERNAL'='FALSE')

INSERT INTO omni_schema_b.t VALUES ('b')`)
	d := &Driver{db: db}
	sess := openSession(t, d)
	tag := func(opts driver.QueryOptions) string {
		results, err := sess.Query(context.Background(), "SELECT tag FROM t", opts)
		require.NoError(t, err)
		require.Len(t, results, 1)
		require.NoError(t, results[0].Err)
		require.Len(t, results[0].Rows, 1)
		return results[0].Rows[0][0].String
	}

	require.Equal(t, "a", tag(driver.QueryOptions{Schema: "omni_schema_a"}))
	// USE takes effect and reports Hive's answer to reading its rows.
	results, err := sess.Query(context.Background(), "USE omni_schema_b", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Error(t, results[0].Err)
	require.Equal(t, "b", tag(driver.QueryOptions{}))
	require.Equal(t, "a", tag(driver.QueryOptions{Schema: "omni_schema_a"}))
}

func openSession(t *testing.T, d *Driver) driver.Session {
	t.Helper()
	sess, err := d.Session(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { _ = sess.Close() })
	return sess
}

func columnTypes(result *driver.Result) []string {
	types := make([]string, 0, len(result.Columns))
	for _, column := range result.Columns {
		types = append(types, column.TypeName)
	}
	return types
}

func columnAndTypes(result *driver.Result) []string {
	out := make([]string, 0, 2*len(result.Columns))
	for _, column := range result.Columns {
		out = append(out, column.Name, column.TypeName)
	}
	return out
}
