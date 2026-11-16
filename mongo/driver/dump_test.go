package driver

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

// gapMarker ends a corpus line a dump cannot restore yet.
const gapMarker = "// not restorable: TestDumpGaps"

// restorable returns the corpus without the lines marked as gaps.
func restorable(corpus string) string {
	var kept []string
	for _, line := range strings.Split(corpus, "\n") {
		if !strings.HasSuffix(line, gapMarker) {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

// TestDump dumps the synced corpus and holds the mongosh script to
// testdata/dump/database.mongo.
func TestDump(t *testing.T) {
	db := readGolden(t, filepath.Join("testdata", "sync", "database.json"))
	before := proto.Clone(db)
	var b strings.Builder
	require.NoError(t, (&Driver{}).Dump(context.Background(), &b, db))
	if !proto.Equal(before, db) {
		t.Errorf("Dump mutated its input metadata")
	}
	writeOrCompareBytes(t, filepath.Join("testdata", "dump", "database.mongo"), []byte(b.String()))
}

// TestDumpRestores syncs the restorable subset of corpus.mongo, dumps it,
// recreates it in another database through the driver, and checks that
// resyncing yields the same collections and indexes.
func TestDumpRestores(t *testing.T) {
	setup := restorable(mustRead(t, filepath.Join("testdata", "corpus.mongo")))
	source, restored, err := dumpAndRestore(t, openDriver(t, testDatabase), testDatabase, "omni_restored", setup)
	require.NoError(t, err)
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("restored database differs from source (-source +restored):\n%s", diff)
	}
}

// TestDumpGaps holds what a dump cannot faithfully restore yet; each entry
// fails once its gap closes so its corpus line moves into TestDumpRestores.
func TestDumpGaps(t *testing.T) {
	t.Run("view", func(t *testing.T) {
		const setup = `db.createCollection("users");
db.users.insertOne({_id: 1, age: 30});
db.runCommand({create: "v_adults", viewOn: "users", pipeline: [{$match: {age: {$gte: 18}}}]});`
		source, restored, err := dumpAndRestore(t, openDriver(t, testDatabase), testDatabase, "omni_restored", setup)
		require.NoError(t, err)
		require.Len(t, source.GetSchemas()[0].GetViews(), 1)
		if len(restored.GetSchemas()[0].GetViews()) > 0 {
			t.Fatal("restore created the view; if the gap closed, remove it")
		}
	})

	t.Run("time-series collection", func(t *testing.T) {
		const setup = `db.runCommand({create: "telemetry", timeseries: {timeField: "at", metaField: "sensor", granularity: "seconds"}});`
		d, client := openDriverWithClient(t, testDatabase)
		source, _, err := dumpAndRestore(t, d, testDatabase, "omni_restored", setup)
		require.NoError(t, err)
		var names []string
		for _, table := range source.GetSchemas()[0].GetTables() {
			names = append(names, table.GetName())
		}
		require.Contains(t, names, "telemetry", "a time-series collection is synced")
		collections, err := client.ListCollections(context.Background(), "omni_restored")
		require.NoError(t, err)
		for _, collection := range collections {
			if collection["name"] != "telemetry" {
				continue
			}
			options, _ := collection["options"].(bson.M)
			if _, ok := options["timeseries"]; ok {
				t.Fatalf("restore kept the time-series options; if the gap closed, remove it: %v", collection)
			}
		}
	})

	t.Run("index options", func(t *testing.T) {
		const setup = `db.createCollection("sessions");
db.sessions.createIndex({created_at: 1}, {expireAfterSeconds: 3600});`
		d, client := openDriverWithClient(t, testDatabase)
		_, _, err := dumpAndRestore(t, d, testDatabase, "omni_restored", setup)
		require.NoError(t, err)
		indexes, err := client.ListIndexes(context.Background(), "omni_restored", "sessions")
		require.NoError(t, err)
		for _, index := range indexes {
			if _, ok := index["expireAfterSeconds"]; ok {
				t.Fatalf("restore kept the TTL; if the gap closed, remove it: %v", index)
			}
		}
	})
}

// restoreDiff compares two synced databases apart from what a dump of
// structure cannot carry: a collection's size and count, which the restored
// database's empty collections report differently, and a view, whose pipeline
// the metadata does not hold.
func restoreDiff(source, restored *metadatapb.DatabaseSchemaMetadata) string {
	return cmp.Diff(source.GetSchemas(), restored.GetSchemas(), protocmp.Transform(),
		protocmp.IgnoreFields(&metadatapb.TableMetadata{}, "data_size", "index_size", "row_count"),
		protocmp.IgnoreFields(&metadatapb.SchemaMetadata{}, "views"),
	)
}

// dumpAndRestore syncs setup on d's database, dumps it, and applies the dump to
// restoredName through the same client, returning both synced databases.
func dumpAndRestore(t *testing.T, d *Driver, sourceName, restoredName, setup string) (*metadatapb.DatabaseSchemaMetadata, *metadatapb.DatabaseSchemaMetadata, error) {
	t.Helper()
	ctx := context.Background()
	if _, err := d.Execute(ctx, setup, driver.ExecuteOptions{}); err != nil {
		t.Fatalf("setup: %v", err)
	}
	source, err := d.SyncDatabase(ctx)
	if err != nil {
		return nil, nil, err
	}
	var dump strings.Builder
	if err := d.Dump(ctx, &dump, source); err != nil {
		return nil, nil, err
	}
	restoredDriver := &Driver{client: d.client, databaseName: restoredName}
	if _, err := restoredDriver.Execute(ctx, dump.String(), driver.ExecuteOptions{}); err != nil {
		return source, nil, err
	}
	restored, err := restoredDriver.SyncDatabase(ctx)
	if err != nil {
		return source, nil, err
	}
	return source, restored, nil
}
