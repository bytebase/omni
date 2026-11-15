package driver

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	metadatapb "github.com/bytebase/omni/metadata"
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
// database and syncs that database again: a dump the server refuses, or one
// that orders an object before what it depends on, fails on the restore, and
// a dump that drops an attribute the server kept fails on the comparison.
// The restored database answers from a schema a dump built, so it holds none
// of the source's physical sizes or the column-position gaps its dropped
// columns left, and its name is the source's with _restored; the comparison
// ignores those.
func TestDumpRestores(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	for _, role := range versions {
		t.Run(role, func(t *testing.T) {
			source, restored := restoreDump(t, role, testcontainer.Database(t), restorable(string(corpus)))
			if diff := restoreDiff(source, restored); diff != "" {
				t.Errorf("the restored database differs from the source (-source +restored):\n%s", diff)
			}
		})
	}
}

// restoreDiff compares a source database and its dump's restore, ignoring
// what a restore cannot carry: the database name, a table's physical sizes
// and row estimate, and the column positions the source's dropped columns
// left gaps in.
func restoreDiff(source, restored *metadatapb.DatabaseSchemaMetadata) string {
	return cmp.Diff(source, restored, protocmp.Transform(),
		protocmp.IgnoreFields(&metadatapb.DatabaseSchemaMetadata{}, "name"),
		protocmp.IgnoreFields(&metadatapb.TableMetadata{}, "data_size", "index_size", "row_count"),
		protocmp.IgnoreFields(&metadatapb.ColumnMetadata{}, "position"),
	)
}

// restoreDump syncs setup into sourceName on the server of role, dumps it,
// restores the dump into an empty database, and syncs that database,
// returning the source and restored metadata. A restore the server refuses
// fails t.
func restoreDump(t *testing.T, role, sourceName, setup string) (source, restored *metadatapb.DatabaseSchemaMetadata) {
	t.Helper()
	m, restoredName, err := dumpAndRestore(t, role, sourceName, setup)
	if err != nil {
		t.Fatalf("restoring the dump: %v", err)
	}
	// Sync on a connection of its own: the dump's header set the restoring
	// connection's search_path, which a sync reads before its own transaction.
	db := testcontainer.OpenDatabaseNumberingOIDs(t, role, restoredName)
	t.Cleanup(func() { db.Close() })
	restoredDriver := &Driver{db: db, databaseName: restoredName}
	restoredMeta, err := restoredDriver.SyncDatabase(context.Background())
	if err != nil {
		t.Fatalf("syncing the restored database: %v", err)
	}
	return m, restoredMeta
}

// dumpAndRestore syncs setup into sourceName on the server of role, dumps
// it, and restores the dump into sourceName+"_restored", returning the
// source metadata, the restored database's name, and the server's restore
// error. Each caller names its own source so that a package's databases hold
// nothing another test made.
func dumpAndRestore(t *testing.T, role, sourceName, setup string) (*metadatapb.DatabaseSchemaMetadata, string, error) {
	t.Helper()
	ctx := context.Background()
	src := &Driver{db: testcontainer.OpenDatabaseNumberingOIDs(t, role, sourceName), databaseName: sourceName}
	t.Cleanup(func() { src.db.Close() })
	mustExec(t, src.db, setup)
	m, err := src.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var dump strings.Builder
	if err := src.Dump(ctx, &dump, m); err != nil {
		t.Fatal(err)
	}
	restoredName := sourceName + "_restored"
	db := testcontainer.OpenDatabase(t, role, restoredName)
	t.Cleanup(func() { db.Close() })
	_, err = db.ExecContext(ctx, dump.String())
	return m, restoredName, err
}

// gapMarker ends a corpus line whose objects a dump cannot faithfully
// restore yet, so TestDumpRestores leaves it out; TestDumpGaps holds each
// such gap and why.
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

// TestDumpGaps holds what a dump cannot faithfully restore yet. A gap whose
// wantErr is set is one the server refuses outright; a gap with an empty
// wantErr restores but loses what the metadata cannot carry, so the
// round-trip differs. A gap that closes fails here: move its corpus line
// back into TestDumpRestores by dropping its gapMarker.
func TestDumpGaps(t *testing.T) {
	for _, gap := range []struct {
		// db names the gap's own source database, so one gap's objects never
		// reach another's sync.
		name, db, setup, wantErr string
	}{
		// Range types are not synced, but their constructor functions are,
		// as functions of the internal language, which do not restore.
		{"range type", "omni_gap_range", "CREATE TYPE float_range AS RANGE (subtype = float8);", "syntax error at or near"},
		// A user-defined collation is not a synced object, so its CREATE
		// COLLATION is never emitted and a column that uses it cannot restore.
		{
			"user-defined collation", "omni_gap_collation",
			"CREATE COLLATION c (provider = icu, locale = 'und-u-ks-level2', deterministic = false);\n" +
				"CREATE TABLE t (v text COLLATE c);",
			"does not exist",
		},
		// A foreign table's server and its options are not synced, so the
		// dump omits the server and the foreign table; the table's synced
		// external_database_name is the local database besides, which a
		// restore into another database could not reproduce.
		{
			"foreign table", "omni_gap_foreign",
			"CREATE EXTENSION postgres_fdw;\n" +
				"CREATE SERVER srv FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'localhost', dbname 'postgres');\n" +
				"CREATE FOREIGN TABLE ft (id int) SERVER srv OPTIONS (schema_name 'public', table_name 'x');",
			"",
		},
	} {
		t.Run(gap.name, func(t *testing.T) {
			if gap.wantErr != "" {
				_, _, err := dumpAndRestore(t, testcontainer.Default, gap.db, gap.setup)
				if err == nil || !strings.Contains(err.Error(), gap.wantErr) {
					t.Fatalf("restoring the dump: %v, want an error containing %q; if the gap closed, remove it", err, gap.wantErr)
				}
				return
			}
			source, restored := restoreDump(t, testcontainer.Default, gap.db, gap.setup)
			if diff := restoreDiff(source, restored); diff == "" {
				t.Fatalf("the restore reproduced the source; if the gap closed, remove it")
			}
		})
	}
}
