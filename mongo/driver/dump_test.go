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

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

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

// TestDumpRestores syncs corpus.mongo, dumps it, recreates it in another
// database through the driver, and checks that resyncing yields the same
// collections, collection options, indexes, index options and views.
func TestDumpRestores(t *testing.T) {
	setup := mustRead(t, filepath.Join("testdata", "corpus.mongo"))
	source, restored, err := dumpAndRestore(t, openDriver(t, testDatabase), testDatabase, "omni_restored", setup)
	require.NoError(t, err)
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("restored database differs from source (-source +restored):\n%s", diff)
	}
}

// restoreDiff compares two synced databases apart from what a dump of
// structure cannot carry: a collection's size and count, which the restored
// database's empty collections report differently.
func restoreDiff(source, restored *metadatapb.DatabaseSchemaMetadata) string {
	return cmp.Diff(source.GetSchemas(), restored.GetSchemas(), protocmp.Transform(),
		protocmp.IgnoreFields(&metadatapb.TableMetadata{}, "data_size", "index_size", "row_count"),
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
