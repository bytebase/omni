package driver

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/pg/internal/testcontainer"
)

// TestDump renders each server's synced corpus, testdata/sync/<role>.json,
// and holds it to testdata/dump/<role>.sql.
func TestDump(t *testing.T) {
	for _, role := range versions {
		t.Run(role, func(t *testing.T) {
			m := readGolden(t, filepath.Join("testdata", "sync", role+".json"))
			before := proto.Clone(m)
			var b strings.Builder
			if err := (&Driver{}).Dump(context.Background(), &b, m); err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(before, m) {
				t.Errorf("Dump changed the metadata it was given")
			}
			golden := filepath.Join("testdata", "dump", role+".sql")
			if *update {
				if err := os.WriteFile(golden, []byte(b.String()), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(string(want), b.String()); diff != "" {
				t.Errorf("Dump of %s differs from %s (-want +got):\n%s", role, golden, diff)
			}
		})
	}
}

// TestDumpRestores restores each server's dump of the corpus into an empty
// database: a dump the server refuses, or one that orders an object before
// what it depends on, fails here.
func TestDumpRestores(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	for _, role := range versions {
		t.Run(role, func(t *testing.T) {
			if err := restoreDump(t, role, restorable(string(corpus))); err != nil {
				t.Fatalf("restoring the dump: %v", err)
			}
		})
	}
}

// restoreDump syncs setup on the server of role, dumps it, and restores the
// dump into an empty database, returning the server's error.
func restoreDump(t *testing.T, role, setup string) error {
	t.Helper()
	ctx := context.Background()
	source := openDriver(t, role)
	mustExec(t, source.db, setup)
	m, err := source.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var dump strings.Builder
	if err := source.Dump(ctx, &dump, m); err != nil {
		t.Fatal(err)
	}
	db := testcontainer.OpenDatabase(t, role, testcontainer.Database(t)+"_restored")
	t.Cleanup(func() { db.Close() })
	_, err = db.ExecContext(ctx, dump.String())
	return err
}

// gapMarker ends a corpus line whose objects a dump cannot restore yet;
// TestDumpGaps holds each such gap.
const gapMarker = "-- not restorable: TestDumpGaps"

// restorable is the corpus without its gapMarker lines.
func restorable(corpus string) string {
	var kept []string
	for _, line := range strings.Split(corpus, "\n") {
		if !strings.HasSuffix(line, gapMarker) {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

// TestDumpGaps holds what a dump cannot restore yet, each with the error
// the server gives. A gap that closes fails here: move its corpus line back
// into TestDumpRestores by dropping its gapMarker.
func TestDumpGaps(t *testing.T) {
	for _, gap := range []struct {
		name, setup, err string
	}{
		// Range types are not synced, but their constructor functions are,
		// as functions of the internal language.
		{"range type", "CREATE TYPE float_range AS RANGE (subtype = float8);", "syntax error at or near"},
	} {
		t.Run(gap.name, func(t *testing.T) {
			if err := restoreDump(t, testcontainer.Default, gap.setup); err == nil || !strings.Contains(err.Error(), gap.err) {
				t.Fatalf("restoring the dump: %v, want an error containing %q; if the gap closed, remove it", err, gap.err)
			}
		})
	}
}
