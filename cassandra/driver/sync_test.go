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

	"github.com/bytebase/omni/cassandra/internal/testcontainer"
	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync from this run")

// testKeyspace is the keyspace every server test reads. A recording creates it
// empty at the start of each top-level test, so a tape holds one test's corpus
// only.
const testKeyspace = "omni_driver"

// openDriver returns a driver on an empty keyspace, so what the tape holds does
// not depend on what the test binary ran before.
func openDriver(t *testing.T) *Driver {
	t.Helper()
	testcontainer.Reset(t,
		"DROP KEYSPACE IF EXISTS "+testKeyspace,
		"CREATE KEYSPACE "+testKeyspace+" WITH replication = {'class':'SimpleStrategy','replication_factor':1}",
	)
	return &Driver{client: testcontainer.Open(t, testKeyspace), config: driver.Config{Database: testKeyspace}}
}

// executeCorpus runs the corpus through the driver, so the tape holds the
// statements its tables were created by.
func executeCorpus(t *testing.T, d *Driver, script string) {
	t.Helper()
	if _, err := d.Execute(context.Background(), script, driver.ExecuteOptions{}); err != nil {
		t.Fatalf("execute corpus: %v", err)
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// TestSyncDatabase syncs testdata/corpus.cql holds and holds the result to
// testdata/sync/database.json.
func TestSyncDatabase(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.cql")))

	got, err := d.SyncDatabase(context.Background())
	require.NoError(t, err)
	writeOrCompareGolden(t, filepath.Join("testdata", "sync", "database.json"), got)
}

// TestSyncDatabaseMissing holds a keyspace the server no longer has to
// driver.ErrNotFound, so a keyspace dropped after the driver opened is reported
// as gone rather than as an empty schema. A driver opened on a keyspace the
// server never had fails earlier: gocql builds no session without it.
func TestSyncDatabaseMissing(t *testing.T) {
	d := openDriver(t)
	// Ping creates the session the driver's client holds, whose keyspace the
	// next step drops; a session outlives its keyspace.
	require.NoError(t, d.Ping(context.Background()))
	testcontainer.Reset(t)
	_, err := d.SyncDatabase(context.Background())
	require.ErrorIs(t, err, driver.ErrNotFound)
}

// TestSyncInstance holds the pinned server's version and keyspaces to
// testdata/sync/instance.json. Cassandra has no roles this driver reads, so the
// golden holds none.
func TestSyncInstance(t *testing.T) {
	testcontainer.Reset(t,
		"CREATE KEYSPACE omni_instance_a WITH replication = {'class':'SimpleStrategy','replication_factor':1}",
		"CREATE KEYSPACE omni_instance_b WITH replication = {'class':'SimpleStrategy','replication_factor':1}",
	)
	d := &Driver{client: testcontainer.Open(t, "")}
	instance, err := d.SyncInstance(context.Background())
	require.NoError(t, err)

	databases := make([]string, 0, len(instance.Databases))
	for _, database := range instance.Databases {
		databases = append(databases, database.Name)
	}
	b, err := json.MarshalIndent(map[string]any{
		"version":   instance.Version,
		"databases": databases,
		"roles":     instance.Roles,
	}, "", "  ")
	require.NoError(t, err)
	writeOrCompareBytes(t, filepath.Join("testdata", "sync", "instance.json"), append(b, '\n'))
}

func writeOrCompareGolden(t *testing.T, path string, got *metadatapb.DatabaseSchemaMetadata) {
	t.Helper()
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		b, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := &metadatapb.DatabaseSchemaMetadata{}
	if err := protojson.Unmarshal(b, want); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("SyncDatabase differs from %s (-want +got):\n%s", path, diff)
	}
}

func writeOrCompareBytes(t *testing.T, path string, got []byte) {
	t.Helper()
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("%s differs (-want +got):\n%s", path, diff)
	}
}
