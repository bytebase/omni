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

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/spanner/internal/testcontainer"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync from this run")

// TestSyncDatabase syncs the fixture database and holds the result to
// testdata/sync/fixture.json: the tables with their columns, indexes and foreign
// keys, and the view with its definition.
func TestSyncDatabase(t *testing.T) {
	a := require.New(t)
	instance, database := syncFixture(t)
	d := openDriver(t, instance, database)
	defer d.Close()

	got, err := d.SyncDatabase(context.Background())
	a.NoError(err)
	syncGolden(t, got, filepath.Join("testdata", "sync", "fixture.json"))
}

// TestSyncDatabasePostgres syncs a database of the PostgreSQL dialect and holds
// the result to testdata/sync/postgres.json: the PostgreSQL catalog's own
// schema name, type names, default expression and match option.
func TestSyncDatabasePostgres(t *testing.T) {
	a := require.New(t)
	instance, database := postgresSyncFixture(t)
	d := openDriverWithDialect(t, instance, database, postgresSQL)
	defer d.Close()

	got, err := d.SyncDatabase(context.Background())
	a.NoError(err)
	syncGolden(t, got, filepath.Join("testdata", "sync", "postgres.json"))
}

// syncGolden holds a synced schema to a golden file, or writes it with -update.
func syncGolden(t *testing.T, got *metadatapb.DatabaseSchemaMetadata, golden string) {
	t.Helper()
	a := require.New(t)
	if *update {
		b, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(got)
		a.NoError(err)
		a.NoError(os.MkdirAll(filepath.Dir(golden), 0o755))
		a.NoError(os.WriteFile(golden, append(b, '\n'), 0o644))
		return
	}
	want := readGolden(t, golden)
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("SyncDatabase differs from %s (-want +got):\n%s", golden, diff)
	}
}

// TestSyncDatabaseMissing holds a database the server does not have to
// driver.ErrNotFound, which is what a caller reports as a missing database
// rather than a failed sync.
func TestSyncDatabaseMissing(t *testing.T) {
	a := require.New(t)
	instance, _ := syncFixture(t)
	d := openDriver(t, instance, "missing")
	defer d.Close()

	_, err := d.SyncDatabase(context.Background())
	a.ErrorIs(err, driver.ErrNotFound)
}

// TestSyncInstance reads the instance's databases with their names and no
// schemas, both dialects included: a database of the PostgreSQL dialect is a
// database, and SyncDatabase reads its catalog.
func TestSyncInstance(t *testing.T) {
	a := require.New(t)
	instance := "syncinstance"
	testcontainer.Create(t, instance, "fixture", fixtureDDL, nil)
	testcontainer.CreatePostgres(t, instance, "pgdialect", nil, nil)
	d := openDriver(t, instance, "fixture")
	defer d.Close()

	got, err := d.SyncInstance(context.Background())
	a.NoError(err)
	var names []string
	for _, database := range got.Databases {
		names = append(names, database.Name)
	}
	a.ElementsMatch([]string{"fixture", "pgdialect"}, names)
	a.Empty(got.Version)
	a.Empty(got.Roles)
}

// TestSyncDatabaseEmpty reads a database that holds no schema: it has no table
// and no view, and so reports no schema.
func TestSyncDatabaseEmpty(t *testing.T) {
	a := require.New(t)
	testcontainer.Create(t, "empty", "empty", nil, nil)
	d := openDriver(t, "empty", "empty")
	defer d.Close()

	got, err := d.SyncDatabase(context.Background())
	a.NoError(err)
	a.Equal("empty", got.Name)
	a.Empty(got.Schemas)
}

// readGolden reads a metadata golden file.
func readGolden(t *testing.T, path string) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	metadata := &metadatapb.DatabaseSchemaMetadata{}
	require.NoError(t, protojson.Unmarshal(b, metadata))
	return metadata
}
