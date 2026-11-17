package driver

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	metadatapb "github.com/bytebase/omni/metadata"
)

// TestDump writes the server's own DDL, holds it to testdata/dump/default.sql,
// then drops the schema, restores that text, and requires the same metadata
// back. CockroachDB publishes no DDL for functions or user types, so the
// corpus holds the object kinds SHOW CREATE ALL TABLES can recreate.
func TestDump(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "dump", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openTestDriver(t)
	ctx := context.Background()
	mustExec(t, d.db, string(corpus))

	before, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := d.Dump(ctx, &out, before); err != nil {
		t.Fatal(err)
	}
	golden := filepath.Join("testdata", "dump", "default.sql")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, out.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(want), out.String()); diff != "" {
		t.Errorf("Dump differs from %s (-want +got):\n%s", golden, diff)
	}

	mustExec(t, d.db, "DROP SCHEMA app CASCADE;")
	mustExec(t, d.db, out.String())
	after, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// obj_description(oid) without a catalog OID can read a comment that
	// belongs to another catalog object with the same OID; the restored
	// objects have new OIDs, so the defect makes those comments differ. The
	// dump's COMMENT ON lines are held in the golden above.
	if diff := cmp.Diff(before, after, protocmp.Transform(),
		protocmp.IgnoreFields(&metadatapb.TableMetadata{}, "comment"),
		protocmp.IgnoreFields(&metadatapb.ColumnMetadata{}, "comment"),
		protocmp.IgnoreFields(&metadatapb.IndexMetadata{}, "comment"),
		protocmp.IgnoreFields(&metadatapb.ViewMetadata{}, "comment"),
		protocmp.IgnoreFields(&metadatapb.MaterializedViewMetadata{}, "comment"),
	); diff != "" {
		t.Errorf("the dump did not restore the schema (-before +after):\n%s", diff)
	}
}
