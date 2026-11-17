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

// dumpDatabase is the database the dump tests read, spelled as a caller may
// spell it: Hive stores a database name in lower case and resolves it whatever
// the case, so the corpus names it the way the server states it and the driver
// is opened on the other spelling.
const dumpDatabase = "OMNI_DUMP"

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

// restoreDiff compares two synced databases apart from the definition text of a
// view and a materialized view. Hive stores a view as the statement that created
// it, and one the dump recreates is stored as the rendering its CREATE stated —
// the same view, whose query the server writes another way. A materialized view
// is replayed from the query DESCRIBE FORMATTED expands, which the server then
// renders as the view's original query, so the two texts name the same query.
// The dump itself round-trips exactly, which the second dump holds. A table's
// size and row count are the rows it holds, which a dump of the schema does not
// carry.
func restoreDiff(source, restored *metadatapb.DatabaseSchemaMetadata) string {
	return cmp.Diff(source.GetSchemas(), restored.GetSchemas(), protocmp.Transform(),
		protocmp.IgnoreFields(&metadatapb.ViewMetadata{}, "definition"),
		protocmp.IgnoreFields(&metadatapb.MaterializedViewMetadata{}, "definition"),
		protocmp.IgnoreFields(&metadatapb.TableMetadata{}, "data_size", "row_count"),
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
		// gap describes. A gap with wantErr is one the server refuses to replay
		// at all.
		name     string
		database string
		setup    []string
		wantErr  bool
		check    func(t *testing.T, g gapResult)
	}{
		{
			// A view whose source is gone keeps the statement the server states,
			// and the server refuses to create it again: Hive resolves a view's
			// query when it creates the view, and the table it reads is not
			// there. The dump carries the view, ordered by name because the
			// server cannot say what it reads.
			name:     "a view whose source is gone",
			database: "omni_gap_broken_view",
			setup: []string{
				"CREATE TABLE `omni_gap_broken_view`.t (id INT) STORED AS PARQUET",
				"CREATE VIEW `omni_gap_broken_view`.v AS SELECT id FROM `omni_gap_broken_view`.t",
				"DROP TABLE `omni_gap_broken_view`.t",
			},
			wantErr: true,
		},
		{
			// Hive names the partition a null value lives in for the value it
			// has lost — dt=__HIVE_DEFAULT_PARTITION__ — and takes neither
			// statement that would state it: ALTER TABLE ... ADD PARTITION with
			// the name as a literal fails on the reserved substring (error
			// 10111), and dt=null fails as a partition column of type void
			// (error 10248). The dump therefore leaves the partition out.
			name:     "a partition whose value is null",
			database: "omni_gap_null_partition",
			setup: []string{
				"CREATE TABLE `omni_gap_null_partition`.p (id INT) PARTITIONED BY (dt STRING) STORED AS PARQUET",
				"ALTER TABLE `omni_gap_null_partition`.p SET TBLPROPERTIES('EXTERNAL'='FALSE')",
				"INSERT INTO `omni_gap_null_partition`.p VALUES (1, null)",
			},
			check: func(t *testing.T, g gapResult) {
				t.Helper()
				partitions := g.source.Schemas[0].Tables[0].Partitions
				if len(partitions) != 1 || partitions[0].Name != "dt=__HIVE_DEFAULT_PARTITION__" {
					t.Fatalf("the source's partitions are %v; the fixture is not what the gap names", partitions)
				}
				if restored := g.restored.Schemas[0].Tables[0].Partitions; len(restored) != 0 {
					t.Fatalf("the dump restored %d partitions; if the gap closed, remove it", len(restored))
				}
				for _, statement := range []string{
					"ALTER TABLE `omni_gap_null_partition`.p ADD PARTITION (dt='__HIVE_DEFAULT_PARTITION__')",
					"ALTER TABLE `omni_gap_null_partition`.p ADD PARTITION (dt=null)",
				} {
					if _, err := g.driver.db.ExecContext(context.Background(), statement); err == nil {
						t.Fatalf("%s is a statement the server takes; if the gap closed, remove it", statement)
					}
				}
			},
		},
		{
			// The dump carries a database's objects, not the database: SHOW
			// CREATE DATABASE states its comment, owner and location, and the
			// caller creates the database a restore goes into. The metadata
			// cannot carry them either — DatabaseSchemaMetadata holds the name,
			// the schemas and the extension list — so the loss needs a field
			// before a dump could state it.
			name:     "the database itself",
			database: "omni_gap_database",
			setup: []string{
				"DROP DATABASE IF EXISTS `omni_gap_database` CASCADE",
				"CREATE DATABASE `omni_gap_database` COMMENT 'a comment the dump leaves out'",
				"CREATE TABLE `omni_gap_database`.t (id INT) STORED AS PARQUET",
			},
			check: func(t *testing.T, g gapResult) {
				t.Helper()
				if strings.Contains(g.dumped, "CREATE DATABASE") {
					t.Fatalf("the dump carries the database; if the gap closed, remove it:\n%s", g.dumped)
				}
				var statement string
				require.NoError(t, g.driver.db.QueryRowContext(context.Background(),
					"SHOW CREATE DATABASE "+quoteBacktick("omni_gap_database")).Scan(&statement))
				if strings.Contains(statement, "a comment the dump leaves out") {
					t.Fatalf("the database's comment came back; if the gap closed, remove it:\n%s", statement)
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
			if gap.wantErr {
				require.Error(t, err, "the dump replays a statement the server refuses")
				return
			}
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

// TestDumpPartitionSpec holds what a partition name states: the columns and
// values ADD PARTITION takes, Hive's escapes read back, and nothing for the
// partition a null value lives in.
func TestDumpPartitionSpec(t *testing.T) {
	for _, partition := range []struct {
		name string
		spec string
		ok   bool
	}{
		{name: "dt=2024-01-01", spec: "`dt`='2024-01-01'", ok: true},
		{name: "dt=2024-01-01/hr=05", spec: "`dt`='2024-01-01', `hr`='05'", ok: true},
		// Hive escapes the characters that would end a level, name one, or spell
		// an escape, so a value holding any of them reads back whole.
		{name: "dt=a%2Fb", spec: "`dt`='a/b'", ok: true},
		{name: "dt=a%3Db/hr=a%2Fb", spec: "`dt`='a=b', `hr`='a/b'", ok: true},
		{name: "dt=a%25b", spec: "`dt`='a%b'", ok: true},
		{name: "dt=a%23b", spec: "`dt`='a#b'", ok: true},
		// A value's quote and backslash are the two a replay escapes again.
		{name: "dt=a%27b", spec: "`dt`='a\\'b'", ok: true},
		{name: "dt=a%5Cb", spec: "`dt`='a\\\\b'", ok: true},
		{name: "dt=__HIVE_DEFAULT_PARTITION__", ok: false},
		{name: "dt=a/dt=__HIVE_DEFAULT_PARTITION__", ok: false},
		{name: "dt", ok: false},
	} {
		spec, ok := dumpPartitionSpec(partition.name)
		require.Equal(t, partition.ok, ok, partition.name)
		require.Equal(t, partition.spec, spec, partition.name)
	}
}

// TestDumpClauseStart holds where a stated CREATE TABLE's column list ends: at
// the parenthesis that closes it, past the parentheses a column type and a
// column comment hold.
func TestDumpClauseStart(t *testing.T) {
	for _, stmt := range []struct {
		name   string
		stmt   string
		clause string
	}{
		{
			name:   "a plain column list",
			stmt:   "CREATE TABLE `db`.`t`(\n  `id` int)\nROW FORMAT SERDE \n  'x'",
			clause: "\nROW FORMAT SERDE \n  'x'",
		},
		{
			name:   "a type holding parentheses",
			stmt:   "CREATE TABLE `db`.`t`(\n  `amount` decimal(10,2))\nLOCATION\n  'x'",
			clause: "\nLOCATION\n  'x'",
		},
		{
			name:   "a comment holding a parenthesis and a quote",
			stmt:   "CREATE TABLE `db`.`t`(\n  `note` string COMMENT 'a ) note, it\\'s')\nCOMMENT 'the table'",
			clause: "\nCOMMENT 'the table'",
		},
		{
			name:   "a name holding a parenthesis",
			stmt:   "CREATE TABLE `db`.`a (b`(\n  `id` int)\nLOCATION\n  'x'",
			clause: "\nLOCATION\n  'x'",
		},
		{
			name:   "a column named with a backquote",
			stmt:   "CREATE TABLE `db`.`t`(\n  `a``b` int)\nLOCATION\n  'x'",
			clause: "\nLOCATION\n  'x'",
		},
	} {
		clause, err := dumpClauseStart(stmt.stmt)
		require.NoError(t, err, stmt.name)
		require.Equal(t, stmt.clause, stmt.stmt[clause:], stmt.name)
	}
}

// TestDumpStatesPartitions holds where a stated CREATE TABLE declares its
// partitions, and that a name or a comment spelling the clause is not one.
func TestDumpStatesPartitions(t *testing.T) {
	for _, stmt := range []struct {
		name  string
		stmt  string
		state bool
	}{
		{name: "a plain table", stmt: "CREATE TABLE `db`.`t`(\n  `id` int)\nLOCATION\n  'x'"},
		{name: "a partitioned table", stmt: "CREATE TABLE `db`.`t`(\n  `id` int)\nPARTITIONED BY ( \n  `dt` string)", state: true},
		{name: "a comment spelling it", stmt: "CREATE TABLE `db`.`t`(\n  `id` int COMMENT 'PARTITIONED BY (x)')"},
		{name: "a name spelling it", stmt: "CREATE TABLE `db`.`PARTITIONED BY`(\n  `id` int)"},
	} {
		require.Equal(t, stmt.state, dumpStatesPartitions(stmt.stmt), stmt.name)
	}
}

// TestDumpMaterializedViewClauses holds the clauses CREATE MATERIALIZED VIEW
// takes from a stated CREATE TABLE: a partition clause names its columns
// without the types the statement states them with.
func TestDumpMaterializedViewClauses(t *testing.T) {
	for _, clauses := range []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "no partition clause",
			in:   "\nCOMMENT 'c'\nROW FORMAT SERDE \n  'x'",
			out:  "\nCOMMENT 'c'\nROW FORMAT SERDE \n  'x'",
		},
		{
			name: "one partition column",
			in:   "\nPARTITIONED BY ( \n  `dt` string)\nROW FORMAT SERDE \n  'x'",
			out:  "\nPARTITIONED ON (`dt`)\nROW FORMAT SERDE \n  'x'",
		},
		{
			name: "a comment before it and a type holding a comma",
			in:   "\nCOMMENT 'a comment'\nPARTITIONED BY ( \n  `region` string, \n  `amount` decimal(10,2))\nLOCATION\n  'x'",
			out:  "\nCOMMENT 'a comment'\nPARTITIONED ON (`region`, `amount`)\nLOCATION\n  'x'",
		},
		{
			name: "an unquoted column",
			in:   "\nPARTITIONED BY (dt string)",
			out:  "\nPARTITIONED ON (dt)",
		},
	} {
		out, err := dumpMaterializedViewClauses(clauses.in)
		require.NoError(t, err, clauses.name)
		require.Equal(t, clauses.out, out, clauses.name)
	}
}
