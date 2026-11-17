package driver

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/hive/internal/testcontainer"
	metadatapb "github.com/bytebase/omni/metadata"
)

// dumpDatabase is the database the dump tests read. The corpus names it, since
// a view's statement has to name the objects it reads.
const dumpDatabase = "omni_dump"

// gapMarker ends a corpus line whose object a dump cannot faithfully restore
// yet, so TestDumpRestores leaves it out; TestDumpGaps holds each such gap and
// why.
const gapMarker = "-- not restorable: TestDumpGaps"

func dumpCorpus(t *testing.T) string {
	t.Helper()
	corpus, err := os.ReadFile(filepath.Join("testdata", "dump", "corpus.sql"))
	require.NoError(t, err)
	return string(corpus)
}

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

// TestDump writes the corpus's DDL as the server states it and holds it to
// testdata/dump/database.sql.
func TestDump(t *testing.T) {
	d := openTestDatabase(t, dumpDatabase)
	execCorpus(t, d.db, dumpCorpus(t))

	var dumped strings.Builder
	require.NoError(t, d.Dump(context.Background(), &dumped, nil))

	golden := filepath.Join("testdata", "dump", "database.sql")
	if *update {
		require.NoError(t, os.WriteFile(golden, []byte(dumped.String()), 0o644))
		return
	}
	want, err := os.ReadFile(golden)
	require.NoError(t, err)
	require.Equal(t, string(want), dumped.String())
}

// TestDumpRestores writes the restorable corpus, dumps it, drops the database,
// applies the dump to a database of the same name, and syncs again. A dump the
// server refuses, or one that loses an attribute the sync read, fails here, as
// does a second dump that states something else.
func TestDumpRestores(t *testing.T) {
	ctx := context.Background()
	d := openTestDatabase(t, dumpDatabase)
	execCorpus(t, d.db, restorable(dumpCorpus(t)))

	source := syncDatabase(t, d)
	before := proto.Clone(source)
	var dump strings.Builder
	require.NoError(t, d.Dump(ctx, &dump, source))
	require.True(t, proto.Equal(before, source), "Dump mutated its input metadata")

	recreateDatabase(t, d, dumpDatabase)
	_, err := d.Execute(ctx, dump.String(), driver.ExecuteOptions{})
	require.NoError(t, err, "restoring the dump")
	restored := syncDatabase(t, d)
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("the restored catalog differs from the source (-source +restored):\n%s", diff)
	}

	var again strings.Builder
	require.NoError(t, d.Dump(ctx, &again, restored))
	require.Equal(t, dump.String(), again.String(), "the restored database dumps to a different script")
}

// restoreDiff compares two synced databases apart from a view's definition
// text. Hive stores a view as the statement that created it, and a view the dump
// recreates is stored as the rendering SHOW CREATE stated — the same view, whose
// query the server writes another way. The dump itself round-trips exactly,
// which the second dump holds.
func restoreDiff(source, restored *metadatapb.DatabaseSchemaMetadata) string {
	return cmp.Diff(source.GetSchemas(), restored.GetSchemas(), protocmp.Transform(),
		protocmp.IgnoreFields(&metadatapb.ViewMetadata{}, "definition"),
	)
}

// TestDumpMissing holds a database the server does not have to an error rather
// than to an empty script. Hive renders a failed statement as a TStatus whose
// message gohive keeps in a field it does not export, so a caller cannot
// recognize the answer as driver.ErrNotFound; the dump returns it as it comes.
func TestDumpMissing(t *testing.T) {
	db := testcontainer.Open(t, "omni_missing_database")
	defer func() { _ = db.Close() }()

	var dumped strings.Builder
	err := (&Driver{db: db, databaseName: "omni_missing_database"}).Dump(context.Background(), &dumped, nil)
	require.Error(t, err)
	require.Empty(t, dumped.String())
}

// gapResult is what a gap's check reads: the synced source, the dump, and the
// synced restore.
type gapResult struct {
	driver   *Driver
	source   *metadatapb.DatabaseSchemaMetadata
	dumped   string
	restored *metadatapb.DatabaseSchemaMetadata
}

// TestDumpGaps holds what a dump cannot faithfully restore. Each gap fails here
// once it closes.
func TestDumpGaps(t *testing.T) {
	for _, gap := range []struct {
		// database names the gap's own database, and check holds the loss the
		// gap describes.
		name     string
		database string
		setup    []string
		check    func(t *testing.T, g gapResult)
	}{
		{
			// Hive 4 stores a CREATE TABLE as an external table, and only
			// ALTER TABLE ... SET TBLPROPERTIES('EXTERNAL'='FALSE') makes one
			// managed — which is how the pinned server, with strict managed
			// tables off so a sync has a managed table to read, makes the two.
			// SHOW CREATE TABLE states the CREATE and nothing of the ALTER, so
			// the table a replay creates is external.
			name:     "a table Hive made managed by an ALTER",
			database: "omni_gap_managed",
			setup: []string{
				"CREATE TABLE `omni_gap_managed`.orders (id BIGINT) STORED AS PARQUET",
				"ALTER TABLE `omni_gap_managed`.orders SET TBLPROPERTIES('EXTERNAL'='FALSE')",
			},
			check: func(t *testing.T, g gapResult) {
				t.Helper()
				if tables := g.source.Schemas[0].Tables; len(tables) != 1 || len(tables[0].Partitions) != 0 {
					t.Fatalf("the source has %d managed tables; the fixture is not what the gap names", len(tables))
				}
				if tables := g.restored.Schemas[0].Tables; len(tables) != 0 {
					t.Fatalf("the dump restored the table as managed; if the gap closed, remove it")
				}
				if tables := g.restored.Schemas[0].ExternalTables; len(tables) != 1 {
					t.Fatalf("the restore has %d external tables; the gap is about the one it made", len(tables))
				}
			},
		},
		{
			// SHOW CREATE TABLE states a partitioned table's PARTITIONED BY
			// columns and not its partitions: those are metastore rows an
			// ALTER TABLE ADD PARTITION writes, and the dump restates only what
			// SHOW CREATE states.
			name:     "a partitioned table's partitions",
			database: "omni_gap_partition",
			setup: []string{
				"CREATE TABLE `omni_gap_partition`.events (id INT) PARTITIONED BY (dt STRING) STORED AS PARQUET",
				"ALTER TABLE `omni_gap_partition`.events SET TBLPROPERTIES('EXTERNAL'='FALSE')",
				"ALTER TABLE `omni_gap_partition`.events ADD PARTITION (dt='2024-01-01')",
			},
			check: func(t *testing.T, g gapResult) {
				t.Helper()
				partitions := g.source.Schemas[0].Tables[0].Partitions
				if len(partitions) != 1 {
					t.Fatalf("the source has %d partitions; the fixture is not what the gap names", len(partitions))
				}
				if strings.Contains(g.dumped, partitions[0].Name) {
					t.Fatalf("the dump names the partition %q; if the gap closed, remove it", partitions[0].Name)
				}
			},
		},
		{
			// A view's comment is a table parameter, and SHOW CREATE TABLE
			// states a view's query and no parameters, so a restored view comes
			// back without it.
			name:     "a view's comment",
			database: "omni_gap_view_comment",
			setup: []string{
				"CREATE TABLE `omni_gap_view_comment`.t (id INT) STORED AS PARQUET",
				"CREATE VIEW `omni_gap_view_comment`.v COMMENT 'a comment the dump leaves out' AS SELECT id FROM `omni_gap_view_comment`.t",
			},
			check: func(t *testing.T, g gapResult) {
				t.Helper()
				before := g.source.Schemas[0].Views
				after := g.restored.Schemas[0].Views
				if len(before) != 1 || before[0].Comment == "" {
					t.Fatalf("the source view has no comment; the fixture is not what the gap names")
				}
				if len(after) != 1 || after[0].Comment != "" {
					t.Fatalf("the dump restored the view's comment; if the gap closed, remove it")
				}
			},
		},
		{
			// SHOW CREATE TABLE states a materialized view as a plain CREATE
			// TABLE: Hive keeps the view's query out of it, so the object a
			// restore creates is the table the statement says.
			name:     "a materialized view's definition",
			database: "omni_gap_materialized_view",
			setup: []string{
				"CREATE TABLE `omni_gap_materialized_view`.src (id INT) STORED AS PARQUET",
				"CREATE MATERIALIZED VIEW `omni_gap_materialized_view`.mv DISABLE REWRITE AS SELECT id FROM `omni_gap_materialized_view`.src",
			},
			check: func(t *testing.T, g gapResult) {
				t.Helper()
				if views := g.source.Schemas[0].MaterializedViews; len(views) != 1 {
					t.Fatalf("the source has %d materialized views; the fixture is not what the gap names", len(views))
				}
				if views := g.restored.Schemas[0].MaterializedViews; len(views) != 0 {
					t.Fatalf("the dump restored the materialized view; if the gap closed, remove it")
				}
			},
		},
		{
			// The dump carries a database's objects, not the database: SHOW
			// CREATE DATABASE states its comment, owner and location, and the
			// caller creates the database a restore goes into. A dump that
			// carried the CREATE DATABASE would either fail on the database the
			// restore already made or bring the comment back.
			name:     "the database itself",
			database: "omni_gap_database",
			setup: []string{
				"DROP DATABASE IF EXISTS `omni_gap_database` CASCADE",
				"CREATE DATABASE `omni_gap_database` COMMENT 'a comment the dump leaves out'",
				"CREATE TABLE `omni_gap_database`.t (id INT) STORED AS PARQUET",
			},
			check: func(t *testing.T, g gapResult) {
				t.Helper()
				var statement string
				require.NoError(t, g.driver.db.QueryRowContext(context.Background(),
					"SHOW CREATE DATABASE "+quoteBacktick("omni_gap_database")).Scan(&statement))
				if strings.Contains(statement, "a comment the dump leaves out") {
					t.Fatalf("the dump restored the database's comment; if the gap closed, remove it:\n%s", statement)
				}
			},
		},
	} {
		t.Run(gap.name, func(t *testing.T) {
			ctx := context.Background()
			d := openTestDatabase(t, gap.database)
			for _, statement := range gap.setup {
				if _, err := d.db.ExecContext(ctx, statement); err != nil {
					t.Fatalf("%s: %v", statement, err)
				}
			}

			var dumped strings.Builder
			require.NoError(t, d.Dump(ctx, &dumped, nil))
			source := syncDatabase(t, d)
			recreateDatabase(t, d, gap.database)
			_, err := d.Execute(ctx, dumped.String(), driver.ExecuteOptions{})
			require.NoError(t, err, "restoring the dump")
			gap.check(t, gapResult{driver: d, source: source, dumped: dumped.String(), restored: syncDatabase(t, d)})
		})
	}
}

// recreateDatabase replaces database with an empty one of the same name, which
// is what a restore applies a dump to. The connection keeps using the name it
// was opened on, which Hive drops and creates again without complaint.
func recreateDatabase(t *testing.T, d *Driver, name string) {
	t.Helper()
	ctx := context.Background()
	conn, err := d.db.Conn(ctx)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()
	for _, statement := range []string{
		"DROP DATABASE IF EXISTS " + quoteBacktick(name) + " CASCADE",
		"CREATE DATABASE " + quoteBacktick(name),
	} {
		if _, err := conn.ExecContext(ctx, statement); err != nil {
			t.Fatalf("%s: %v", statement, err)
		}
	}
}

// syncDatabase syncs d's database.
func syncDatabase(t *testing.T, d *Driver) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	metadata, err := d.SyncDatabase(context.Background())
	require.NoError(t, err)
	return metadata
}
