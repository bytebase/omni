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
	"github.com/bytebase/omni/pg/catalog"
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

// TestDumpRestoresGeneratedColumnFunction checks that a STORED generated
// column whose expression calls a user function restores: the dump must
// create the function before the table, though the table's name sorts first.
func TestDumpRestoresGeneratedColumnFunction(t *testing.T) {
	setup := "CREATE FUNCTION public.zz_double(n integer) RETURNS integer LANGUAGE sql IMMUTABLE AS $$ SELECT n * 2 $$;\n" +
		"CREATE TABLE public.aa_gen (base integer, doubled integer GENERATED ALWAYS AS (public.zz_double(base)) STORED);"
	source, restored := restoreDump(t, testcontainer.Default, "omni_gen_fn", setup)
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("the restored database differs from the source (-source +restored):\n%s", diff)
	}
}

// TestDumpRestoresOverloadedFunctionDefault checks that a column default
// calling an overloaded function restores even when another overload takes
// the table's own row type: the dump must create the integer overload before
// the table without making the table depend on the row-type overload.
func TestDumpRestoresOverloadedFunctionDefault(t *testing.T) {
	setup := "CREATE FUNCTION public.aa_f(n integer) RETURNS integer LANGUAGE sql IMMUTABLE AS $$ SELECT n $$;\n" +
		"CREATE TABLE public.zz_t (id integer DEFAULT public.aa_f(1));\n" +
		"CREATE FUNCTION public.aa_f(r public.zz_t) RETURNS integer LANGUAGE sql IMMUTABLE AS $$ SELECT 1 $$;"
	source, restored := restoreDump(t, testcontainer.Default, "omni_overload_fn", setup)
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("the restored database differs from the source (-source +restored):\n%s", diff)
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

// renderer writes a synced snapshot as the text a restore executes.
type renderer func(*metadatapb.DatabaseSchemaMetadata) (string, error)

// dumpDefinition renders the pg_dump-style script pg/driver writes.
func dumpDefinition(m *metadatapb.DatabaseSchemaMetadata) (string, error) {
	var dump strings.Builder
	if err := (&Driver{}).Dump(context.Background(), &dump, m); err != nil {
		return "", err
	}
	return dump.String(), nil
}

// sdlDefinition renders the declarative form of the same snapshot.
func sdlDefinition(m *metadatapb.DatabaseSchemaMetadata) (string, error) {
	return catalog.DatabaseDefinition(m, catalog.DefinitionOptions{SDLFormat: true})
}

// restoreDump syncs setup into sourceName on the server of role, dumps it,
// restores the dump into an empty database, and syncs that database, returning
// the source and restored metadata. A restore the server refuses fails t.
func restoreDump(t *testing.T, role, sourceName, setup string) (source, restored *metadatapb.DatabaseSchemaMetadata) {
	t.Helper()
	return restoreDefinition(t, role, sourceName, setup, dumpDefinition)
}

// restoreDefinition is restoreDump with the form of the definition chosen: the
// text it restores is render(source).
func restoreDefinition(t *testing.T, role, sourceName, setup string, render renderer) (source, restored *metadatapb.DatabaseSchemaMetadata) {
	t.Helper()
	m, restoredName, err := renderAndRestore(t, role, sourceName, setup, render)
	if err != nil {
		t.Fatalf("restoring the definition: %v", err)
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
	m, restoredName, err := renderAndRestore(t, role, sourceName, setup, dumpDefinition)
	return m, restoredName, err
}

// renderAndRestore syncs setup into sourceName, restores render(snapshot) into
// sourceName+"_restored", and returns the source metadata, the restored
// database's name, and the server's restore error.
func renderAndRestore(t *testing.T, role, sourceName, setup string, render renderer) (*metadatapb.DatabaseSchemaMetadata, string, error) {
	t.Helper()
	ctx := context.Background()
	src := &Driver{db: testcontainer.OpenDatabaseNumberingOIDs(t, role, sourceName), databaseName: sourceName}
	t.Cleanup(func() { src.db.Close() })
	mustExec(t, src.db, setup)
	m, err := src.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	text, err := render(m)
	if err != nil {
		t.Fatal(err)
	}
	restoredName := sourceName + "_restored"
	db := testcontainer.OpenDatabase(t, role, restoredName)
	t.Cleanup(func() { db.Close() })
	_, err = db.ExecContext(ctx, text)
	return m, restoredName, err
}

// TestSDLRestoresPartitions restores a partitioned schema from its SDL instead
// of its dump. SDL is what a schema is written down as, so it must name every
// partition, its bound and its index attachment: a restore that misses one
// loses the partitioning the sync read.
func TestSDLRestoresPartitions(t *testing.T) {
	for _, role := range versions {
		t.Run(role, func(t *testing.T) {
			source, restored := restoreDefinition(t, role, "omni_sdl_partition", partitionedSchema, sdlDefinition)
			if diff := restoreDiff(source, restored); diff != "" {
				t.Errorf("the schema the SDL restores differs from the source (-source +restored):\n%s", diff)
			}
		})
	}
}

// partitionedSchema is the schema the SDL restore is held to: a range partition
// sub-partitioned by list with a default partition, a partition whose bounds are
// multi-column with MINVALUE and MAXVALUE, a hash partition, quoted mixed-case
// names, a primary key every partition inherits, partitioned indexes, and a
// table attached after it was created with a CHECK, an index and a primary key
// of its own.
const partitionedSchema = `
CREATE SCHEMA app;
CREATE TABLE app."Measure" (
    "City" integer NOT NULL,
    "Day" date NOT NULL,
    "Region" text NOT NULL,
    CONSTRAINT measure_city_check CHECK ("City" > 0)
) PARTITION BY RANGE ("Day");
CREATE TABLE app."Measure 2025" PARTITION OF app."Measure"
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01') PARTITION BY LIST ("Region");
CREATE TABLE app."Measure 2025 EU" PARTITION OF app."Measure 2025" FOR VALUES IN ('eu');
CREATE TABLE app."Measure 2025 Other" PARTITION OF app."Measure 2025" DEFAULT;
CREATE TABLE app."Measure 2026" PARTITION OF app."Measure"
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');
CREATE INDEX measure_city ON app."Measure" ("City", "Day");
CREATE UNIQUE INDEX measure_key ON app."Measure" ("City", "Day", "Region");
CREATE TABLE app.bounds (a integer NOT NULL, b integer NOT NULL) PARTITION BY RANGE (a, b);
CREATE TABLE app.bounds_low PARTITION OF app.bounds FOR VALUES FROM (MINVALUE, MINVALUE) TO (10, MAXVALUE);
CREATE TABLE app.bounds_high PARTITION OF app.bounds FOR VALUES FROM (10, MAXVALUE) TO (MAXVALUE, MAXVALUE);
CREATE TABLE app.sharded (id integer) PARTITION BY HASH (id);
CREATE TABLE app.sharded_0 PARTITION OF app.sharded FOR VALUES WITH (MODULUS 2, REMAINDER 0);
CREATE TABLE app.sharded_1 PARTITION OF app.sharded FOR VALUES WITH (MODULUS 2, REMAINDER 1);
CREATE TABLE app.pk (id integer NOT NULL, d date NOT NULL, PRIMARY KEY (id, d)) PARTITION BY RANGE (d);
CREATE TABLE app.pk_1 PARTITION OF app.pk FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');
CREATE TABLE app.pk_2 PARTITION OF app.pk DEFAULT;
CREATE TABLE app.late (id integer NOT NULL, v text) PARTITION BY LIST (v);
CREATE TABLE app.late_x (id integer NOT NULL, v text NOT NULL, CONSTRAINT late_x_positive CHECK (id > 0));
ALTER TABLE app.late ATTACH PARTITION app.late_x FOR VALUES IN ('x');
CREATE INDEX late_x_id ON app.late_x (id);
ALTER TABLE app.late_x ADD CONSTRAINT late_x_pkey PRIMARY KEY (id);
`

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
// round-trip differs. A gap that names what the restored server answers lost
// is one the metadata cannot carry at all: the round trip agrees, and only the
// server shows what the source had. A gap that closes fails here: move its
// corpus line back into TestDumpRestores by dropping its gapMarker.
func TestDumpGaps(t *testing.T) {
	for _, gap := range []struct {
		// db names the gap's own source database, so one gap's objects never
		// reach another's sync.
		name, db, setup, wantErr string
		// lost is a query the restored database answers, and want what it
		// answers, for a gap the synced snapshot holds nothing of.
		lost, want string
	}{
		// Range types are not synced, but their constructor functions are,
		// as functions of the internal language, which do not restore.
		{"range type", "omni_gap_range", "CREATE TYPE float_range AS RANGE (subtype = float8);", "syntax error at or near", "", ""},
		// A user-defined collation is not a synced object, so its CREATE
		// COLLATION is never emitted and a column that uses it cannot restore.
		{
			"user-defined collation", "omni_gap_collation",
			"CREATE COLLATION c (provider = icu, locale = 'und-u-ks-level2', deterministic = false);\n" +
				"CREATE TABLE t (v text COLLATE c);",
			"does not exist", "", "",
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
			"", "", "",
		},
		// TablePartitionMetadata carries a partition's name, type, bound,
		// indexes and constraints, but neither the schema it lives in nor its
		// comment, so the dump recreates a partition beside its parent without
		// what the source said about it. The sync reads neither, so the round
		// trip agrees and only the server shows the loss.
		{
			"partition beside its parent", "omni_gap_partition",
			"CREATE SCHEMA app;\n" +
				"CREATE SCHEMA other;\n" +
				"CREATE TABLE app.p (id int, v text) PARTITION BY LIST (v);\n" +
				"CREATE TABLE other.p_a PARTITION OF app.p FOR VALUES IN ('a');\n" +
				"COMMENT ON TABLE other.p_a IS 'remote';",
			"",
			"SELECT n.nspname || ':' || coalesce(obj_description(c.oid, 'pg_class'), '') FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE c.relname = 'p_a'",
			"app:",
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
			diff := restoreDiff(source, restored)
			if gap.lost == "" {
				if diff == "" {
					t.Fatalf("the restore reproduced the source; if the gap closed, remove it")
				}
				return
			}
			if diff != "" {
				t.Fatalf("the restore differs from the source, so the snapshot does carry what the gap names; if the gap closed, remove it:\n%s", diff)
			}
			db := testcontainer.OpenDatabase(t, testcontainer.Default, gap.db+"_restored")
			defer db.Close()
			var got string
			if err := db.QueryRow(gap.lost).Scan(&got); err != nil {
				t.Fatalf("reading the restored database: %v", err)
			}
			if got != gap.want {
				t.Fatalf("the restore answered %q, want %q; if the gap closed, remove it", got, gap.want)
			}
		})
	}
}
