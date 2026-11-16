package driver

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
)

// TestQuery reads a row of every kind the driver converts, a row whose columns
// are NULL, and a second statement of the same script.
func TestQuery(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.cql")))
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), `
		SELECT id, txt, ascii_col, vc, flag, bi, ti, si, dbl, flt, dec, vi, ts, dt, tm, uid, tuid, bl, dur, lst, st, mp, addr
		FROM omni_types WHERE id = 1;
		SELECT id, txt, dec, vi FROM omni_types WHERE id = 2;`, driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 2)
	require.NoError(t, results[0].Err)
	require.NoError(t, results[1].Err)

	columns := make([]string, 0, len(results[0].Columns))
	types := make([]string, 0, len(results[0].Columns))
	for _, column := range results[0].Columns {
		columns = append(columns, column.Name)
		types = append(types, column.TypeName)
	}
	require.Equal(t, []string{
		"id", "txt", "ascii_col", "vc", "flag", "bi", "ti", "si", "dbl", "flt", "dec", "vi",
		"ts", "dt", "tm", "uid", "tuid", "bl", "dur", "lst", "st", "mp", "addr",
	}, columns)
	// Cassandra reports a text column as varchar: the two names are one CQL
	// type.
	require.Equal(t, []string{
		"int", "varchar", "ascii", "varchar", "boolean", "bigint", "tinyint", "smallint", "double", "float",
		"decimal", "varint", "timestamp", "date", "time", "uuid", "timeuuid", "blob", "duration",
		"list", "set", "map", "inet",
	}, types)

	require.Len(t, results[0].Rows, 1)
	row := results[0].Rows[0]
	require.Equal(t, driver.Int64, row[0].Kind)
	require.Equal(t, int64(1), row[0].Int64)
	require.Equal(t, driver.String, row[1].Kind)
	require.Equal(t, "text value", row[1].String)
	require.Equal(t, "ascii", row[2].String)
	require.Equal(t, "varchar", row[3].String)
	require.Equal(t, driver.Bool, row[4].Kind)
	require.True(t, row[4].Bool)
	require.Equal(t, int64(9223372036854775807), row[5].Int64)
	require.Equal(t, driver.Int32, row[6].Kind)
	require.Equal(t, int64(7), row[6].Int64)
	require.Equal(t, int64(300), row[7].Int64)
	require.Equal(t, driver.Float64, row[8].Kind)
	require.Equal(t, 3.5, row[8].Float64)
	require.Equal(t, driver.Float32, row[9].Kind)
	require.Equal(t, 1.25, row[9].Float64)
	require.Equal(t, "1.5", row[10].String)
	require.Equal(t, "123456789012345678901234567890", row[11].String)
	require.Equal(t, driver.Timestamp, row[12].Kind)
	require.Equal(t, "2026-01-02T03:04:05.678Z", row[12].Time.UTC().Format(time.RFC3339Nano))
	require.Equal(t, int32(3), row[12].Precision)
	require.Equal(t, "2026-01-02T00:00:00Z", row[13].Time.UTC().Format(time.RFC3339Nano))
	require.Equal(t, "08:12:54.123456", row[14].String)
	require.Equal(t, "6afaf7c2-c0c2-0ddc-7f22-71b8f0a0b0c0", row[15].String)
	require.Equal(t, "50554d6e-29bb-11e5-b345-feff819cdc9f", row[16].String)
	require.Equal(t, driver.Bytes, row[17].Kind)
	require.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, row[17].Bytes)
	require.Equal(t, "1mo2d3h0m0s", row[18].String)
	require.Equal(t, "[1,2,3]", row[19].String)
	require.Equal(t, `["a","b"]`, row[20].String)
	require.Equal(t, `{"a":1,"b":2}`, row[21].String)
	// inet has no case of its own: Bytebase read it as its JSON spelling, which
	// for a net.IP is its address text, quoted as JSON.
	require.Equal(t, `"10.0.0.1"`, row[22].String)

	// A NULL column reads as its type's zero value: gocql hands back a pointer
	// to the zero value, so Bytebase's nil check fired for the two types whose
	// Go type is itself a pointer.
	nullRow := results[1].Rows[0]
	require.Equal(t, int64(2), nullRow[0].Int64)
	require.Equal(t, "", nullRow[1].String)
	require.Equal(t, driver.Null, nullRow[2].Kind)
	require.Equal(t, driver.Null, nullRow[3].Kind)
}

// TestQueryLimit holds QueryOptions.Limit to the page size and the row count: a
// statement with more rows stops at the limit, and the rows it keeps come in the
// clustering order the query reads.
func TestQueryLimit(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.cql")))
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(),
		"SELECT bucket, event_id, payload FROM omni_events WHERE group_id = 1",
		driver.QueryOptions{Limit: 2})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 2)
	require.Equal(t, int64(2), results[0].Rows[0][0].Int64)
	require.Equal(t, int64(3), results[0].Rows[0][1].Int64)
	require.Equal(t, "p1", results[0].Rows[0][2].String)
	require.Equal(t, int64(1), results[0].Rows[1][0].Int64)
	require.Equal(t, int64(4), results[0].Rows[1][1].Int64)

	// A table with more rows than the page size: the limit stops the reading at
	// its fourth row rather than paging on to the end.
	results, err = sess.Query(context.Background(),
		"SELECT id, label FROM omni_rows", driver.QueryOptions{Limit: 4})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 4)
}

// TestQueryWithoutALimit reads every row of a statement in one page. The driver
// sends a page size only when a limit caps the rows, so Cassandra answers
// without paging here; the loop that continues from a page state is Bytebase's
// and a tape records one page.
func TestQueryWithoutALimit(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.cql")))
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(),
		"SELECT id, label FROM omni_rows", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 10)
}

// TestQueryMaxResultBytes holds the cap to the bytes the rows encode to: the
// reading stops at the first doubling past it and says so.
func TestQueryMaxResultBytes(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.cql")))
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(),
		"SELECT id, label FROM omni_rows", driver.QueryOptions{MaxResultBytes: 60})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.True(t, results[0].Truncated)
	// The header is 22 bytes and a row about 15, so the fourth row crosses 60.
	require.Len(t, results[0].Rows, 4)
}

// TestQueryError holds a failing statement to its own Result: the server's
// error, with the rows it read before it, and the script stopping there.
func TestQueryError(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.cql")))
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(),
		"SELECT id FROM omni_rows;\nSELECT missing FROM omni_rows;\nSELECT id FROM omni_types;", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 2)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 10)
	require.Error(t, results[1].Err)
	require.Contains(t, results[1].Err.Error(), "failed to fetch row data")
	require.Contains(t, results[1].Err.Error(), "Undefined column name missing")
	require.Equal(t, "SELECT missing FROM omni_rows", results[1].Statement)
}

// TestQueryEmpty holds an empty script to no results.
func TestQueryEmpty(t *testing.T) {
	d := openDriver(t)
	sess := openSession(t, d)
	results, err := sess.Query(context.Background(), " ; \n;", driver.QueryOptions{})
	require.NoError(t, err)
	require.Empty(t, results)
}

func openSession(t *testing.T, d *Driver) driver.Session {
	t.Helper()
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })
	return session
}
