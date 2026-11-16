package driver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/dynamodb/internal/testcontainer"
)

// TestQuery reads the fixture's table whose items carry every DynamoDB type,
// so that the columns of an answer are the union of the attributes its items
// named and each column's type is the type of the attribute it holds.
func TestQuery(t *testing.T) {
	session := openSession(t)
	results, err := session.Query(context.Background(), `SELECT * FROM "people"`, driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	result := results[0]
	require.NoError(t, result.Err)
	require.Equal(t, `SELECT * FROM "people"`, result.Statement)
	require.False(t, result.Truncated)
	require.Nil(t, result.AffectedRows, "a read returns rows, not a count")
	require.Len(t, result.Rows, testcontainer.PeopleCount)

	// The columns are the attributes the items named, sorted by name, and each
	// carries the type of the attribute it holds.
	require.Equal(t, []string{
		"active", "age", "email", "flags", "id", "late", "list", "meta",
		"missing", "name", "payload", "scores", "tags",
	}, resultColumnNames(result))
	types := make([]string, 0, len(result.Columns))
	for _, column := range result.Columns {
		types = append(types, column.TypeName)
	}
	require.Equal(t, []string{
		"Boolean", "Number", "String", "BinarySet", "String", "String", "List",
		"Map", "Null", "String", "Binary", "NumberSet", "StringSet",
	}, types)

	rows := rowsByID(t, result)
	require.Len(t, rows, testcontainer.PeopleCount)

	// Every type crosses once: a set, a list and a map as their JSON text, a
	// number as its text, a binary as bytes, and a binary inside a document as
	// base64, so bytes that are not valid UTF-8 survive the text form.
	require.Equal(t, map[string]string{
		"active":  "true",
		"age":     "33",
		"email":   "u1@example.com",
		"flags":   `["/g==","/w=="]`,
		"id":      testcontainer.PersonWithEveryType,
		"late":    null,
		"list":    `["7","s","/g=="]`,
		"meta":    `{"k":"v","raw":"/w=="}`,
		"missing": null,
		"name":    "Ann",
		"payload": "b64:/wD+",
		"scores":  `["1","2"]`,
		"tags":    `["a","b"]`,
	}, rows[testcontainer.PersonWithEveryType])

	// An item that leaves attributes out leaves nils, and an attribute only a
	// later item names is backfilled into the rows before it.
	require.Equal(t, map[string]string{
		"id":    testcontainer.PersonWithFewAttributes,
		"email": "u2@example.com",
		"name":  "Bob",
		"age":   "44",
		"late":  null,
		"tags":  null,
	}, pick(rows[testcontainer.PersonWithFewAttributes], "id", "email", "name", "age", "late", "tags"))
	require.Equal(t, "last", rows[testcontainer.PersonWithALateAttribute]["late"])
	require.Equal(t, null, rows[testcontainer.PersonWithALateAttribute]["tags"])
}

// TestQueryProjection reads the columns a statement names rather than the
// attributes the items hold.
func TestQueryProjection(t *testing.T) {
	session := openSession(t)
	results, err := session.Query(context.Background(), `SELECT name, age FROM "people" WHERE id = 'u2'`, driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Equal(t, []string{"age", "name"}, resultColumnNames(results[0]))
	require.Equal(t, [][]string{{"44", "Bob"}}, rowTexts(results[0]))
}

// TestQueryPages reads a result the endpoint answers over several calls: a
// filter makes each call answer fewer rows than it evaluated, so the token is
// what carries the read, and QueryOptions.Limit is the cap on the rows rather
// than the size of a page.
func TestQueryPages(t *testing.T) {
	session := openSession(t)
	statement := `SELECT * FROM "paged" WHERE owner = 'o1' AND label <> 'row-0' ORDER BY seq`
	results, err := session.Query(context.Background(), statement, driver.QueryOptions{Limit: 3})
	require.NoError(t, err)
	require.Len(t, results, 1)
	result := results[0]
	require.NoError(t, result.Err)
	require.False(t, result.Truncated)
	require.Equal(t, []string{"s1", "s2", "s3"}, sequence(t, result), "the cap stops the read at three rows, in the order the sort key decides")
}

// TestQueryLimitCapsRows holds that QueryOptions.Limit is the rows a result
// holds, not the items one call evaluates: a read that would answer more rows
// stops at the cap.
func TestQueryLimitCapsRows(t *testing.T) {
	session := openSession(t)
	results, err := session.Query(context.Background(),
		`SELECT * FROM "paged" WHERE owner = 'o1' ORDER BY seq`, driver.QueryOptions{Limit: 2})
	require.NoError(t, err)
	require.Len(t, results, 1)
	result := results[0]
	require.NoError(t, result.Err)
	require.False(t, result.Truncated)
	require.Equal(t, []string{"s0", "s1"}, sequence(t, result))
}

// TestQueryMaxResultBytes holds that a result stops at the byte cap and says so.
func TestQueryMaxResultBytes(t *testing.T) {
	session := openSession(t)
	results, err := session.Query(context.Background(), `SELECT * FROM "people"`, driver.QueryOptions{MaxResultBytes: 400})
	require.NoError(t, err)
	require.Len(t, results, 1)
	result := results[0]
	require.NoError(t, result.Err)
	require.True(t, result.Truncated)
	require.NotEmpty(t, result.Rows)
	require.Less(t, len(result.Rows), testcontainer.PeopleCount, "the cap stopped the rows before the answer's last one")
}

// TestQueryMaxResultBytesStopsTheRead holds that the byte cap stops the read and
// not only the rows: the tape holds one call, so a driver that asked for the
// next page would fail its replay.
func TestQueryMaxResultBytesStopsTheRead(t *testing.T) {
	session := openSession(t)
	results, err := session.Query(context.Background(),
		`SELECT * FROM "paged" WHERE owner = 'o1' AND label <> 'row-0' ORDER BY seq`,
		driver.QueryOptions{Limit: 3, MaxResultBytes: 30})
	require.NoError(t, err)
	require.Len(t, results, 1)
	result := results[0]
	require.NoError(t, result.Err)
	require.True(t, result.Truncated, "the first page passed the cap, so the read stopped there")
	require.NotEmpty(t, result.Rows)
	require.Less(t, len(result.Rows), 3, "fewer rows than the row cap would have taken")
}

// TestQueryMaxResultBytesWithoutACap holds that a result without a cap is whole.
func TestQueryMaxResultBytesWithoutACap(t *testing.T) {
	session := openSession(t)
	results, err := session.Query(context.Background(), `SELECT * FROM "paged" WHERE owner = 'o1' ORDER BY seq`, driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.False(t, results[0].Truncated)
	require.Len(t, results[0].Rows, testcontainer.PagedCount, "the driver reads every answer's token, so the rows are the whole result")
}

// sequence is the sort key of each row of a read of Paged, in order.
func sequence(t *testing.T, result *driver.Result) []string {
	t.Helper()
	seq := columnIndex(t, result, "seq")
	seqs := make([]string, 0, len(result.Rows))
	for _, row := range result.Rows {
		seqs = append(seqs, row[seq].String)
	}
	return seqs
}

// TestQueryEndsAtTheFirstFailure holds that a statement's failure is its own
// result's error and that the statements after it do not run.
func TestQueryEndsAtTheFirstFailure(t *testing.T) {
	session := openSession(t)
	results, err := session.Query(context.Background(),
		`SELECT * FROM "people" WHERE id = 'u3'; SELECT * FROM "nope"; SELECT * FROM "people" WHERE id = 'u1'`,
		driver.QueryOptions{})
	require.NoError(t, err, "a statement that failed ends the script; the call itself could run")
	require.Len(t, results, 2)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 1)
	require.Error(t, results[1].Err)
	require.Contains(t, results[1].Err.Error(), "ResourceNotFoundException")
	require.Contains(t, results[1].Err.Error(), "non-existent table")
	require.Equal(t, ` SELECT * FROM "nope";`, results[1].Statement)
	require.Empty(t, results[1].Rows)
}

// TestQueryErrorTextIsTheClients holds that a statement's error keeps the
// client's text, which is what the caller reports.
func TestQueryErrorTextIsTheClients(t *testing.T) {
	session := openSession(t)
	results, err := session.Query(context.Background(), `SELECT * FROM "nope"`, driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.EqualError(t, results[0].Err,
		`operation error DynamoDB: ExecuteStatement, https response error StatusCode: 400, RequestID: 00000000-0000-0000-0000-000000000001, ResourceNotFoundException: Cannot do operations on a non-existent table`)
}

// TestQuerySplitsStatementsWithPartiQL holds that a script's boundaries are
// partiql's judgment: a comment-only segment and a lone semicolon hold no
// statement and are not sent.
func TestQuerySplitsStatementsWithPartiQL(t *testing.T) {
	session := openSession(t)
	results, err := session.Query(context.Background(),
		"-- nothing here\n;\nSELECT * FROM \"people\" WHERE id = 'u3'",
		driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Len(t, results[0].Rows, 1)
}

// TestQueryExplain holds that asking for a plan is the call's error: PartiQL
// has no plan to ask for, which is what Bytebase answered.
func TestQueryExplain(t *testing.T) {
	// No client: the driver answers without asking the endpoint.
	d := &Driver{}
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	_, err = session.Query(context.Background(), `SELECT * FROM "people"`, driver.QueryOptions{Explain: true})
	require.EqualError(t, err, "DynamoDB does not support EXPLAIN")
}

// TestQueryWithoutStatements holds that a script of no statements wakes nothing.
func TestQueryWithoutStatements(t *testing.T) {
	d := &Driver{}
	session, err := d.Session(context.Background())
	require.NoError(t, err)
	for _, script := range []string{"", "  \n\t", "; ;", "-- only a comment", "/* only a comment */"} {
		results, err := session.Query(context.Background(), script, driver.QueryOptions{})
		require.NoError(t, err, "script %q", script)
		require.Empty(t, results, "script %q", script)
	}
}

// null is how a nil cell renders.
const null = "<null>"

// openSession returns the fixture's session: a replay of the pinned endpoint's
// answers or, in a recording, the endpoint itself.
func openSession(t *testing.T) driver.Session {
	t.Helper()
	session, err := newDriver(t).Session(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, session.Close()) })
	return session
}

// resultColumnNames is the columns of one result, in order.
func resultColumnNames(result *driver.Result) []string {
	names := make([]string, 0, len(result.Columns))
	for _, column := range result.Columns {
		names = append(names, column.Name)
	}
	return names
}

func columnIndex(t *testing.T, result *driver.Result, name string) int {
	t.Helper()
	for i, column := range result.Columns {
		if column.Name == name {
			return i
		}
	}
	t.Fatalf("no column %q in %v", name, resultColumnNames(result))
	return -1
}

// rowsByID renders every row as the cell it holds for each column, keyed by the
// row's id.
func rowsByID(t *testing.T, result *driver.Result) map[string]map[string]string {
	t.Helper()
	id := columnIndex(t, result, "id")
	rows := make(map[string]map[string]string, len(result.Rows))
	for _, row := range result.Rows {
		cells := make(map[string]string, len(result.Columns))
		for i, column := range result.Columns {
			cells[column.Name] = render(row[i])
		}
		rows[row[id].String] = cells
	}
	return rows
}

// rowTexts renders each row as its cells in column order.
func rowTexts(result *driver.Result) [][]string {
	rows := make([][]string, 0, len(result.Rows))
	for _, row := range result.Rows {
		cells := make([]string, 0, len(row))
		for _, value := range row {
			cells = append(cells, render(value))
		}
		rows = append(rows, cells)
	}
	return rows
}

// render is one cell as a test reads it.
func render(value driver.Value) string {
	switch value.Kind {
	case driver.Null:
		return null
	case driver.Bool:
		if value.Bool {
			return "true"
		}
		return "false"
	case driver.String:
		return value.String
	case driver.Bytes:
		return "b64:" + base64.StdEncoding.EncodeToString(value.Bytes)
	}
	b, err := json.Marshal(value)
	if err != nil {
		return string(value.Kind)
	}
	return string(b)
}

// pick keeps the named cells of one row.
func pick(row map[string]string, names ...string) map[string]string {
	picked := make(map[string]string, len(names))
	for _, name := range names {
		picked[name] = row[name]
	}
	return picked
}
