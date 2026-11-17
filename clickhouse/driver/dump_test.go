package driver

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bytebase/omni/clickhouse/internal/testcontainer"
	"github.com/bytebase/omni/driver"
)

// gapMarker ends a corpus line whose object a dump cannot faithfully restore
// yet, so TestDumpRestores leaves it out; TestDumpGaps holds each such gap and
// why.
const gapMarker = "-- not restorable: TestDumpGaps"

// dumpCorpus is testdata/dump/corpus.sql, with the database the test created put
// where the fixture names it: a dictionary's engine argument and a materialized
// view's literal hold a qualified name.
func dumpCorpus(t *testing.T) string {
	t.Helper()
	corpus, err := os.ReadFile(filepath.Join("testdata", "dump", "corpus.sql"))
	require.NoError(t, err)
	return strings.ReplaceAll(string(corpus), "$database", testcontainer.Database(t))
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
	d := openTestDriver(t)
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
// server refuses, or one that loses an attribute the sync read, fails here.
func TestDumpRestores(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, restorable(dumpCorpus(t)))

	source, err := d.SyncDatabase(ctx)
	require.NoError(t, err)
	before := proto.Clone(source)
	var dump strings.Builder
	require.NoError(t, d.Dump(ctx, &dump, source))
	require.True(t, proto.Equal(before, source), "Dump mutated its input metadata")

	recreateDatabase(t, d, testcontainer.Database(t))
	_, err = d.Execute(ctx, dump.String(), driver.ExecuteOptions{})
	require.NoError(t, err, "restoring the dump")

	restored, err := d.SyncDatabase(ctx)
	require.NoError(t, err)
	if diff := cmp.Diff(source, restored, protocmp.Transform()); diff != "" {
		t.Errorf("the restored catalog differs from the source (-source +restored):\n%s", diff)
	}
}

// TestDumpMissing holds a database the server does not have to driver.ErrNotFound:
// a dump is a read of the database it names, not a sync that reports an empty
// schema.
func TestDumpMissing(t *testing.T) {
	db := testcontainer.OpenMissingDatabase(t, "omni_missing")
	defer func() { _ = db.Close() }()

	var dumped strings.Builder
	err := (&Driver{db: db, databaseName: "omni_missing"}).Dump(context.Background(), &dumped, nil)
	require.ErrorIs(t, err, driver.ErrNotFound)
}

// TestDumpGaps holds what a dump cannot faithfully restore. Each gap fails here
// once it closes.
func TestDumpGaps(t *testing.T) {
	for _, gap := range []struct {
		// database names the gap's own database, and check reads the restored
		// one and fails if the gap has closed.
		name     string
		database string
		setup    []string
		check    func(t *testing.T, db *sql.DB, database string)
	}{
		{
			// The dump carries a database's objects, not the database: only
			// SHOW CREATE DATABASE states its engine and comment, and the
			// caller creates the database a restore goes into. A dump that
			// carried the CREATE DATABASE would either fail on the database
			// the restore already made or bring the comment back.
			name:     "the database itself",
			database: "omni_gap_database",
			setup: []string{
				"DROP DATABASE IF EXISTS `omni_gap_database`",
				"CREATE DATABASE `omni_gap_database` ENGINE = Atomic COMMENT 'a comment the dump leaves out'",
				"CREATE TABLE `omni_gap_database`.t (id UInt64) ENGINE = MergeTree ORDER BY id",
			},
			check: func(t *testing.T, db *sql.DB, database string) {
				t.Helper()
				var comment string
				require.NoError(t, db.QueryRowContext(context.Background(),
					"SELECT comment FROM system.databases WHERE name = ?", database).Scan(&comment))
				if comment != "" {
					t.Fatalf("the dump restored the database's comment %q; if the gap closed, remove it", comment)
				}
			},
		},
		{
			// A materialized view with its own engine has a table the server
			// named itself, ".inner_id.<uuid>", which the sync reads as a
			// table of the database. The dump leaves it out, since a restore
			// would otherwise carry a second one beside the view's own, so the
			// restored database's inner table is a name the source never had.
			name:     "a materialized view's inner table",
			database: "omni_gap_mv",
			setup: []string{
				"DROP DATABASE IF EXISTS `omni_gap_mv`",
				"CREATE DATABASE `omni_gap_mv`",
				"CREATE TABLE `omni_gap_mv`.src (id UInt64) ENGINE = MergeTree ORDER BY id",
				"CREATE MATERIALIZED VIEW `omni_gap_mv`.mv ENGINE = MergeTree ORDER BY id AS SELECT id FROM `omni_gap_mv`.src",
			},
			check: func(t *testing.T, db *sql.DB, database string) {
				t.Helper()
				// The source table and the view came back, and the view's new
				// inner table is the one the server named for this restore. A
				// dump that carried the source's inner table would leave two,
				// since the view makes one of its own.
				countTables(t, db, database, "name NOT LIKE '.%'", 2)
				countTables(t, db, database, "name LIKE '.%'", 1)
			},
		},
	} {
		t.Run(gap.name, func(t *testing.T) {
			ctx := context.Background()
			db := testcontainer.Open(t)
			defer func() { _ = db.Close() }()
			d := &Driver{db: db, databaseName: gap.database}
			for _, statement := range gap.setup {
				_, err := db.ExecContext(ctx, statement)
				require.NoError(t, err)
			}

			var dumped strings.Builder
			require.NoError(t, d.Dump(ctx, &dumped, nil))

			recreateDatabase(t, d, gap.database)
			_, err := d.Execute(ctx, dumped.String(), driver.ExecuteOptions{})
			require.NoError(t, err, "restoring the dump")
			gap.check(t, db, gap.database)
		})
	}
}

// countTables holds a database to n tables matching condition.
func countTables(t *testing.T, db *sql.DB, database, condition string, n int) {
	t.Helper()
	query := "SELECT count() FROM system.tables WHERE database = ? AND " + condition
	var got int
	require.NoError(t, db.QueryRowContext(context.Background(), query, database).Scan(&got))
	require.Equalf(t, n, got, "%s in %s", condition, database)
}

// recreateDatabase replaces database with an empty one of the same name, which
// is what a restore applies a dump to.
func recreateDatabase(t *testing.T, d *Driver, name string) {
	t.Helper()
	ctx := context.Background()
	conn, err := d.db.Conn(ctx)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()
	for _, statement := range []string{
		"DROP DATABASE IF EXISTS " + dumpIdentifier(name),
		"CREATE DATABASE " + dumpIdentifier(name),
	} {
		if _, err := conn.ExecContext(ctx, statement); err != nil {
			t.Fatalf("%s: %v", statement, err)
		}
	}
}
