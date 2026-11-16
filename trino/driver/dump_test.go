package driver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	metadatapb "github.com/bytebase/omni/metadata"
)

// TestDump renders testdata/sync/database.json and holds it to
// testdata/dump/database.sql.
func TestDump(t *testing.T) {
	metadata := readGolden(t, filepath.Join("testdata", "sync", "database.json"))
	before := proto.Clone(metadata)
	var b strings.Builder
	if err := (&Driver{}).Dump(context.Background(), &b, metadata); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(before, metadata) {
		t.Errorf("Dump mutated its input metadata")
	}
	writeOrCompareBytes(t, filepath.Join("testdata", "dump", "database.sql"), []byte(b.String()))
}

// TestDumpRestores syncs testdata/corpus.sql, dumps it, drops the schema,
// applies the dump to recreate it, and syncs again. A dump the server refuses,
// or one that loses an attribute the sync reported, fails here.
func TestDumpRestores(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openDriver(t)
	mustExec(t, d, string(corpus))
	source := syncDatabase(t, d)
	if view := findView(source, testSchema, "v_types"); view == nil || !strings.Contains(view.Definition, "c_integer") {
		t.Fatalf("the sync did not read the corpus view: %v", view)
	}

	var dump strings.Builder
	if err := d.Dump(context.Background(), &dump, source); err != nil {
		t.Fatal(err)
	}
	dropSchema(t, d, testSchema)
	if err := restoreDump(t, d, dump.String()); err != nil {
		t.Fatalf("restoring the dump: %v", err)
	}
	restored := syncDatabase(t, d)
	if diff := cmp.Diff(source, restored, protocmp.Transform()); diff != "" {
		t.Errorf("the restored catalog differs from the source (-source +restored):\n%s", diff)
	}
}

func findView(db *metadatapb.DatabaseSchemaMetadata, schema, name string) *metadatapb.ViewMetadata {
	for _, s := range db.Schemas {
		if s.Name != schema {
			continue
		}
		for _, v := range s.Views {
			if v.Name == name {
				return v
			}
		}
	}
	return nil
}

func syncDatabase(t *testing.T, d *Driver) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return metadata
}

func dropSchema(t *testing.T, d *Driver, schema string) {
	t.Helper()
	exec(t, d, fmt.Sprintf("DROP SCHEMA IF EXISTS memory.%s CASCADE", schema))
}

func exec(t *testing.T, d *Driver, statement string) {
	t.Helper()
	if _, err := d.db.ExecContext(context.Background(), statement); err != nil {
		t.Fatalf("%s: %v", statement, err)
	}
}

// restoreDump runs each statement of dump, returning the first the server
// refused.
func restoreDump(t *testing.T, d *Driver, dump string) error {
	t.Helper()
	for _, command := range splitCommands(dump) {
		if _, err := d.db.ExecContext(context.Background(), command.Text); err != nil {
			return err
		}
	}
	return nil
}
