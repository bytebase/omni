package driver

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bytebase/omni/dynamodb/internal/testcontainer"
	metadatapb "github.com/bytebase/omni/metadata"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata")

// TestSyncDatabase syncs the fixture's region and holds the result to
// testdata/sync/database.json.
func TestSyncDatabase(t *testing.T) {
	database, err := newDriver(t).SyncDatabase(context.Background())
	require.NoError(t, err)
	writeOrCompareGolden(t, filepath.Join("testdata", "sync", "database.json"), database)
}

// TestSyncDatabaseReadsWhatTheEndpointReports holds the parts of a synced table
// that the golden alone would not make obvious: the key schema as indexes and
// columns, and the item count the endpoint estimates.
func TestSyncDatabaseReadsWhatTheEndpointReports(t *testing.T) {
	database, err := newDriver(t).SyncDatabase(context.Background())
	require.NoError(t, err)
	require.Equal(t, testDatabaseName, database.Name)
	require.Len(t, database.Schemas, 1)

	tables := map[string]*metadatapb.TableMetadata{}
	for _, table := range database.Schemas[0].Tables {
		tables[table.Name] = table
	}
	require.Len(t, tables, 4, "the fixture holds four tables")

	people := tables[testcontainer.People]
	require.NotNil(t, people)
	// A partition key and a sort key, in the order DescribeTable reports them.
	require.Equal(t, []string{"email", "id"}, columnNames(people))
	require.Equal(t, "id", people.Indexes[0].Name)
	require.Equal(t, "HASH", people.Indexes[0].Type)
	require.Equal(t, "email", people.Indexes[1].Name)
	require.Equal(t, "RANGE", people.Indexes[1].Type)
	require.Equal(t, int64(testcontainer.PeopleCount), people.RowCount)

	// A table whose key has no sort key still reports a RANGE index, empty: the
	// driver writes one for each key type whatever the table has. Bytebase's
	// driver did the same.
	orders := tables[testcontainer.Orders]
	require.NotNil(t, orders)
	require.Equal(t, []string{"pk"}, columnNames(orders))
	require.Len(t, orders.Indexes, 2)
	require.Equal(t, "RANGE", orders.Indexes[1].Type)
	require.Empty(t, orders.Indexes[1].Name)
	require.Empty(t, orders.Indexes[1].Expressions)
}

// TestListAllTablesWalksPages holds that the listing the driver reads is every
// table of a region, not the first answer's 100 names.
func TestListAllTablesWalksPages(t *testing.T) {
	d := &Driver{client: OpenTapePaged(t)}
	tables, err := d.listAllTables(context.Background())
	require.NoError(t, err)
	require.Equal(t, testcontainer.PagedTableNames(), tables)
}

func columnNames(table *metadatapb.TableMetadata) []string {
	names := make([]string, 0, len(table.Columns))
	for _, column := range table.Columns {
		names = append(names, column.Name)
	}
	return names
}

// writeOrCompareGolden rewrites the golden under -update and otherwise compares
// the proto it holds to got. The bytes are not compared: protojson inserts an
// extra space after a key at random per build, so a golden recorded by one
// build does not match another's bytes.
func writeOrCompareGolden(t *testing.T, path string, got *metadatapb.DatabaseSchemaMetadata) {
	t.Helper()
	if *update {
		b, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(got)
		require.NoError(t, err)
		writeOrCompareBytes(t, path, append(b, '\n'))
		return
	}
	want := readGolden(t, path)
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("%s (-want +got):\n%s", path, diff)
	}
}

func writeOrCompareBytes(t *testing.T, path string, got []byte) {
	t.Helper()
	if *update {
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, got, 0o644))
		return
	}
	want, err := os.ReadFile(path)
	require.NoError(t, err)
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("%s differs (-want +got):\n%s", path, diff)
	}
}

func readGolden(t *testing.T, path string) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	metadata := &metadatapb.DatabaseSchemaMetadata{}
	require.NoError(t, protojson.Unmarshal(b, metadata))
	return metadata
}
