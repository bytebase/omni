package driver

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/trino/internal/testcontainer"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata")

const (
	// testCatalog is the catalog every test reads: the image's memory
	// connector, which the tape empties before each top-level test.
	testCatalog = "memory"
	// testSchema is the schema corpus.sql creates.
	testSchema = "omni_driver"
)

// openDriver returns a driver on t's tape. The catalog is the memory connector,
// and the session starts in no schema so a test's statement text decides what
// it reads.
func openDriver(t *testing.T) *Driver {
	t.Helper()
	db := testcontainer.OpenDB(t, testCatalog, "")
	t.Cleanup(func() { db.Close() })
	return &Driver{db: db, databaseName: testCatalog}
}

// mustExec runs script's statements on d, one request each, as the Trino
// protocol requires.
func mustExec(t *testing.T, d *Driver, script string) {
	t.Helper()
	for _, command := range splitCommands(script) {
		if _, err := d.db.ExecContext(context.Background(), command.Text); err != nil {
			t.Fatalf("%s: %v", command.Text, err)
		}
	}
}

// TestSyncDatabase syncs testdata/corpus.sql and holds the result to
// testdata/sync/database.json.
func TestSyncDatabase(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openDriver(t)
	mustExec(t, d, string(corpus))

	got, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, schema := range got.Schemas {
		if schema.Name == "information_schema" {
			t.Errorf("the sync reported the catalog's information_schema as a schema of its own")
		}
	}
	writeOrCompareGolden(t, filepath.Join("testdata", "sync", "database.json"), got)
}

// TestSyncInstance holds the pinned server's version and catalogs to
// testdata/sync/instance.json.
func TestSyncInstance(t *testing.T) {
	d := openDriver(t)
	instance, err := d.SyncInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	databases := make([]string, 0, len(instance.Databases))
	for _, database := range instance.Databases {
		databases = append(databases, database.Name)
	}
	b, err := json.MarshalIndent(map[string]any{
		"version":   instance.Version,
		"databases": databases,
		"roles":     instance.Roles,
	}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	writeOrCompareBytes(t, filepath.Join("testdata", "sync", "instance.json"), append(b, '\n'))
}

// writeOrCompareGolden rewrites the golden under -update and otherwise compares
// the proto it holds to metadata. The bytes are not compared: protojson inserts
// an extra space after a key at random per build, so a golden recorded by one
// build does not match another's bytes.
func writeOrCompareGolden(t *testing.T, path string, metadata *metadatapb.DatabaseSchemaMetadata) {
	t.Helper()
	if *update {
		b, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(metadata)
		if err != nil {
			t.Fatal(err)
		}
		writeOrCompareBytes(t, path, append(b, '\n'))
		return
	}
	want := readGolden(t, path)
	if diff := cmp.Diff(want, metadata, protocmp.Transform()); diff != "" {
		t.Errorf("%s (-want +got):\n%s", path, diff)
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

func readGolden(t *testing.T, path string) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	metadata := &metadatapb.DatabaseSchemaMetadata{}
	if err := protojson.Unmarshal(b, metadata); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return metadata
}
