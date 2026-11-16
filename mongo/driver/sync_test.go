package driver

import (
	"context"
	"encoding/json"
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
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata")

// testDatabase is the database every server test reads. A recording drops it
// at the start of each top-level test, so a tape holds one test's corpus only.
const testDatabase = "omni_driver"

// TestSyncDatabase syncs testdata/corpus.mongo and holds the result to
// testdata/sync/database.json.
func TestSyncDatabase(t *testing.T) {
	d := openDriver(t, testDatabase)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.mongo")))

	got, err := d.SyncDatabase(context.Background())
	require.NoError(t, err)
	writeOrCompareGolden(t, filepath.Join("testdata", "sync", "database.json"), got)
}

// TestSyncInstance holds the pinned server's version, databases and users to
// testdata/sync/instance.json.
func TestSyncInstance(t *testing.T) {
	d := openDriver(t, testDatabase)
	// MongoDB lists a database only once something is created in it.
	_, err := d.Execute(context.Background(), `db.createCollection("users");`, driver.ExecuteOptions{})
	require.NoError(t, err)
	instance, err := d.SyncInstance(context.Background())
	require.NoError(t, err)
	databases := make([]string, 0, len(instance.Databases))
	for _, database := range instance.Databases {
		databases = append(databases, database.Name)
	}
	roles := make([]string, 0, len(instance.Roles))
	for _, role := range instance.Roles {
		roles = append(roles, role.Name)
	}
	b, err := json.MarshalIndent(map[string]any{
		"version":   instance.Version,
		"databases": databases,
		"roles":     roles,
	}, "", "  ")
	require.NoError(t, err)
	writeOrCompareBytes(t, filepath.Join("testdata", "sync", "instance.json"), append(b, '\n'))
}

// TestSyncDatabaseMissing holds a missing database to driver.ErrNotFound.
func TestSyncDatabaseMissing(t *testing.T) {
	d := openDriver(t, "omni_absent")
	_, err := d.SyncDatabase(context.Background())
	require.ErrorIs(t, err, driver.ErrNotFound)
	require.ErrorContains(t, err, "database omni_absent does not exist")
}

func executeCorpus(t *testing.T, d *Driver, corpus string) {
	t.Helper()
	_, err := d.Execute(context.Background(), corpus, driver.ExecuteOptions{})
	require.NoError(t, err)
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(b)
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
