package driver

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/cassandra/internal/testcontainer"
	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/enginetest/tape"
)

// TestQuery reads a row of every kind the driver converts, a row whose columns
// are NULL, and a second statement of the same script.
func TestQuery(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.cql")))
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(), `
		SELECT id, txt, ascii_col, vc, flag, bi, ti, si, dbl, flt, dec, vi, ts, dt, tm, uid, tuid, bl, dur, lst, st, mp, addr, tup
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
		"ts", "dt", "tm", "uid", "tuid", "bl", "dur", "lst", "st", "mp", "addr", "tup",
	}, columns)
	// Cassandra reports a text column as varchar: the two names are one CQL
	// type.
	require.Equal(t, []string{
		"int", "varchar", "ascii", "varchar", "boolean", "bigint", "tinyint", "smallint", "double", "float",
		"decimal", "varint", "timestamp", "date", "time", "uuid", "timeuuid", "blob", "duration",
		"list", "set", "map", "inet", "tuple",
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
	// A tuple is one column holding one value: gocql expands its elements into
	// separate values, and the driver reads them back as one JSON cell.
	require.Equal(t, driver.String, row[23].Kind)
	require.Equal(t, `[7,"seven"]`, row[23].String)

	// A NULL column is the contract's Null, not the type's zero value: the tape
	// held only the zero value before the row read could tell them apart.
	nullRow := results[1].Rows[0]
	require.Equal(t, int64(2), nullRow[0].Int64)
	require.Equal(t, driver.Null, nullRow[1].Kind)
	require.Equal(t, driver.Null, nullRow[2].Kind)
	require.Equal(t, driver.Null, nullRow[3].Kind)
}

// TestQueryKeepsNullApartFromAZeroValue reads a row of every CQL type holding
// NULL and a row holding each type's zero value: a NULL crosses as the
// contract's Null, and a zero value as the value it is.
func TestQueryKeepsNullApartFromAZeroValue(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.cql")))
	checkNullAndZeroRows(t, openSession(t, d))
}

// TestLiveQueryKeepsNullApartFromAZeroValue holds the read the driver makes of a
// live gocql session. A tape records what gocql answered through the read
// cassandra/internal/testcontainer makes, so no replay reaches the driver's own
// read: this test runs only against the pinned server.
func TestLiveQueryKeepsNullApartFromAZeroValue(t *testing.T) {
	tape.Live(t)
	testcontainer.Reset(t,
		"DROP KEYSPACE IF EXISTS "+testKeyspace,
		"CREATE KEYSPACE "+testKeyspace+" WITH replication = {'class':'SimpleStrategy','replication_factor':1}",
	)
	d := &Driver{
		client: &liveClient{session: testcontainer.LiveSession(t, testKeyspace)},
		config: driver.Config{Database: testKeyspace},
	}
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.cql")))
	checkNullAndZeroRows(t, openSession(t, d))
}

// checkNullAndZeroRows holds a query of the corpus's NULL row and zero row to
// what the contract calls each cell.
func checkNullAndZeroRows(t *testing.T, sess driver.Session) {
	t.Helper()
	results, err := sess.Query(context.Background(), `
		SELECT id, n, txt, ascii_col, vc, flag, bi, ti, si, dbl, flt, dec, vi, ts, dt, tm, uid, tuid, bl, dur, lst, st, mp, addr, tup
		FROM omni_nulls WHERE id = 1;
		SELECT id, n, txt, ascii_col, vc, flag, bi, ti, si, dbl, flt, dec, vi, ts, dt, tm, uid, tuid, bl, dur, lst, st, mp, addr, tup
		FROM omni_nulls WHERE id = 2;`, driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 2)
	require.NoError(t, results[0].Err)
	require.NoError(t, results[1].Err)

	columns := make([]string, 0, len(results[0].Columns))
	for _, column := range results[0].Columns {
		columns = append(columns, column.Name)
	}
	require.Equal(t, []string{
		"id", "n", "txt", "ascii_col", "vc", "flag", "bi", "ti", "si", "dbl", "flt", "dec", "vi",
		"ts", "dt", "tm", "uid", "tuid", "bl", "dur", "lst", "st", "mp", "addr", "tup",
	}, columns)

	require.Len(t, results[0].Rows, 1)
	nullRow := results[0].Rows[0]
	require.Len(t, nullRow, len(columns))
	for i, column := range columns {
		// id is the row's primary key, which a CQL table cannot hold NULL in.
		if column == "id" {
			require.Equal(t, int64(1), nullRow[i].Int64)
			continue
		}
		require.Equal(t, driver.Null, nullRow[i].Kind, "%s holds NULL", column)
	}

	require.Len(t, results[1].Rows, 1)
	zeroRow := results[1].Rows[0]
	want := map[string]driver.Value{
		"id":        {Kind: driver.Int64, Int64: 2},
		"n":         {Kind: driver.Int64, Int64: 0},
		"txt":       {Kind: driver.String},
		"ascii_col": {Kind: driver.String},
		"vc":        {Kind: driver.String},
		"flag":      {Kind: driver.Bool},
		"bi":        {Kind: driver.Int64},
		"ti":        {Kind: driver.Int32},
		"si":        {Kind: driver.Int32},
		"dbl":       {Kind: driver.Float64},
		"flt":       {Kind: driver.Float32},
		"dec":       {Kind: driver.String, String: "0"},
		"vi":        {Kind: driver.String, String: "0"},
		"ts":        {Kind: driver.Timestamp, Time: time.Unix(0, 0).UTC(), Precision: 3},
		"dt":        {Kind: driver.Timestamp, Time: time.Unix(0, 0).UTC(), Precision: 3},
		"tm":        {Kind: driver.String, String: "00:00:00"},
		"uid":       {Kind: driver.String, String: "00000000-0000-0000-0000-000000000000"},
		"tuid":      {Kind: driver.String, String: "00000000-0000-1000-0000-000000000000"},
		"bl":        {Kind: driver.Bytes},
		"dur":       {Kind: driver.String, String: "0s"},
		// Cassandra stores an empty collection as NULL, so a list, set or map
		// has no zero value to hold: the NULL row pins it.
		"lst":  {Kind: driver.Null},
		"st":   {Kind: driver.Null},
		"mp":   {Kind: driver.Null},
		"addr": {Kind: driver.String, String: `"0.0.0.0"`},
		"tup":  {Kind: driver.String, String: `[0,""]`},
	}
	for i, column := range columns {
		value, ok := want[column]
		require.True(t, ok, "%s has no expectation", column)
		require.Equal(t, value, zeroRow[i], "%s holds its zero value", column)
	}
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
// reading stops at the first row that takes the kept rows past it.
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
	// The header is 22 bytes and a row about 15, so the third row crosses 60
	// and is the last one kept.
	require.Len(t, results[0].Rows, 3)
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
	// The server's error keeps the code and message a caller reads off it, as
	// the package doc says.
	var requestErr gocql.RequestError
	require.ErrorAs(t, results[1].Err, &requestErr)
	// 8704 is Cassandra's 0x2200, INVALID.
	require.Equal(t, 8704, requestErr.Code())
	require.Contains(t, requestErr.Message(), "Undefined column name missing")
}

// TestQueryEmpty holds an empty script to no results.
func TestQueryEmpty(t *testing.T) {
	d := openDriver(t)
	sess := openSession(t, d)
	results, err := sess.Query(context.Background(), " ; \n;", driver.QueryOptions{})
	require.NoError(t, err)
	require.Empty(t, results)
}

// TestQuerySkipsACommentOnlySegment holds a script that ends in a comment to
// its statements: a comment is no statement, and sending one would fail a
// script whose queries already ran.
func TestQuerySkipsACommentOnlySegment(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.cql")))
	sess := openSession(t, d)

	results, err := sess.Query(context.Background(),
		"SELECT id FROM omni_rows WHERE id = 1;\n-- done\n/* and a block */", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 1)
}

func openSession(t *testing.T, d *Driver) driver.Session {
	t.Helper()
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })
	return session
}
