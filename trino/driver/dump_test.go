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

// TestDumpRestores syncs testdata/corpus.sql, dumps it, recreates the schema,
// applies the dump and syncs again. A dump the server refuses, or one that
// loses an attribute the sync reported, fails here.
func TestDumpRestores(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openDriver(t)
	mustExec(t, d, string(corpus))
	source := syncDatabase(t, d)

	var dump strings.Builder
	if err := d.Dump(context.Background(), &dump, source); err != nil {
		t.Fatal(err)
	}
	recreateSchema(t, d, testSchema)
	if err := restoreDump(t, d, dump.String()); err != nil {
		t.Fatalf("restoring the dump: %v", err)
	}
	restored := syncDatabase(t, d)
	if diff := cmp.Diff(source, restored, protocmp.Transform()); diff != "" {
		t.Errorf("the restored catalog differs from the source (-source +restored):\n%s", diff)
	}
}

// TestDumpGaps holds what a dump cannot restore yet; each entry fails once its
// gap closes, so the corpus line that covers it can move into TestDumpRestores.
func TestDumpGaps(t *testing.T) {
	for _, gap := range []struct {
		name    string
		schema  string
		setup   string
		wantErr string
		// recreate rebuilds the schema before the restore. False keeps the
		// schema dropped, which is what a restore into a fresh catalog does.
		recreate bool
		// gone names a view that the dump cannot carry out of the sync at all:
		// the sync never read it, so no metadata comparison can see the loss.
		gone string
	}{
		{
			name:     "quoted identifier",
			schema:   "omni_gap_quote",
			recreate: true,
			setup: "CREATE SCHEMA memory.omni_gap_quote;\n" +
				"CREATE TABLE memory.omni_gap_quote.\"a\"\"b\" (id integer, \"c\"\"d\" varchar);",
			// The deparser writes names verbatim in double quotes without
			// escaping an embedded quote, so the CREATE TABLE it emits is a
			// syntax error rather than the source table.
			wantErr: "USER_ERROR",
		},
		{
			name:     "schema not created",
			schema:   "omni_gap_schema",
			recreate: false,
			setup: "CREATE SCHEMA memory.omni_gap_schema;\n" +
				"CREATE TABLE memory.omni_gap_schema.t (id integer);",
			// The dump writes CREATE TABLE IF NOT EXISTS "schema"."table" and no
			// CREATE SCHEMA, so it restores only into a schema that exists.
			wantErr: "Schema omni_gap_schema not found",
		},
		{
			name:     "view",
			schema:   "omni_gap_view",
			recreate: true,
			setup: "CREATE SCHEMA memory.omni_gap_view;\n" +
				"CREATE TABLE memory.omni_gap_view.t (id integer);\n" +
				"CREATE VIEW memory.omni_gap_view.v AS SELECT id FROM memory.omni_gap_view.t;",
			gone: "v",
		},
	} {
		t.Run(gap.name, func(t *testing.T) {
			d := openDriver(t)
			mustExec(t, d, gap.setup)
			source := syncDatabase(t, d)

			var dump strings.Builder
			if err := d.Dump(context.Background(), &dump, source); err != nil {
				t.Fatal(err)
			}
			if gap.recreate {
				recreateSchema(t, d, gap.schema)
			} else {
				dropSchema(t, d, gap.schema)
			}

			err := restoreDump(t, d, dump.String())
			if gap.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), gap.wantErr) {
					t.Fatalf("restore err = %v, want %q; if the gap closed, remove it", err, gap.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("restoring the dump: %v", err)
			}

			restored := syncDatabase(t, d)
			if diff := cmp.Diff(source, restored, protocmp.Transform()); diff != "" {
				t.Errorf("the restored catalog differs from the source (-source +restored):\n%s", diff)
			}
			if gap.gone != "" {
				if hasView(t, d, gap.schema, gap.gone) {
					t.Fatalf("the restored catalog holds view %q; if the gap closed, remove it", gap.gone)
				}
			}
		})
	}
}

func syncDatabase(t *testing.T, d *Driver) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return metadata
}

// recreateSchema drops schema and creates it again, which is how a restore
// starts: a Trino dump writes no CREATE SCHEMA.
func recreateSchema(t *testing.T, d *Driver, schema string) {
	t.Helper()
	dropSchema(t, d, schema)
	exec(t, d, fmt.Sprintf("CREATE SCHEMA memory.%s", schema))
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

// hasView reports whether schema holds a view named name, as the server's own
// metadata lists it.
func hasView(t *testing.T, d *Driver, schema, name string) bool {
	t.Helper()
	rows, err := d.db.QueryContext(context.Background(),
		"SELECT table_name FROM system.jdbc.tables WHERE table_cat = ? AND table_schem = ? AND table_type = 'VIEW'",
		testCatalog, schema)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			t.Fatal(err)
		}
		if tableName == name {
			return true
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return false
}
