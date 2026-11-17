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
		// Hive escapes a column's name as it escapes a value's.
		{name: "d%25t=2024-01-01", spec: "`d%t`='2024-01-01'", ok: true},
		{name: "d%3Di%2Fx=a%2Fb", spec: "`d=i/x`='a/b'", ok: true},
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
		{name: "a comment spelling it beside the clause", stmt: "CREATE TABLE `db`.`t`(\n  `id` int)\nCOMMENT 'PARTITIONED BY (x)'\nPARTITIONED BY ( \n  `dt` string)", state: true},
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
			name: "a comment spelling the clause",
			in:   "\nCOMMENT 'says PARTITIONED BY (old)'\nPARTITIONED BY ( \n  `dt` string)\nLOCATION\n  'x'",
			out:  "\nCOMMENT 'says PARTITIONED BY (old)'\nPARTITIONED ON (`dt`)\nLOCATION\n  'x'",
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

// TestDumpMaterializedViewQuery holds the text a materialized view's query
// replays from: the expansion, which quotes every name and, for a partitioned
// view, wraps the query in the projection Hive reads the partitions through, or
// the text the view was created with, which its normalisation may have left
// without the backquotes a name needs.
func TestDumpMaterializedViewQuery(t *testing.T) {
	// p_events_mv's pair, as the pinned server states them.
	const (
		partitionedOriginal = "SELECT id, dt FROM omni_dump.events"
		partitionedExpanded = "SELECT `id`, `dt` FROM (SELECT `events`.`id`, `events`.`dt` FROM `omni_dump`.`events`) `p_events_mv`"
	)
	for _, query := range []struct {
		name        string
		original    string
		expanded    string
		partitioned bool
		want        string
	}{
		{
			name:     "a view whose partitions are not its own",
			original: "SELECT 1bad\nFROM omni_dump.z_users",
			expanded: "SELECT `z_users`.`1bad`\n\nFROM `omni_dump`.`z_users`",
			want:     "SELECT `z_users`.`1bad`\n\nFROM `omni_dump`.`z_users`",
		},
		{
			name:        "a partitioned view whose expansion quotes only names the text states",
			original:    partitionedOriginal,
			expanded:    partitionedExpanded,
			partitioned: true,
			want:        partitionedOriginal,
		},
		{
			name:        "a partitioned view whose expansion quotes a name the text lost",
			original:    "SELECT id AS it's, dt FROM omni_dump.events",
			expanded:    "SELECT `id`, `dt` FROM (SELECT `events`.`id` AS `it's`, `events`.`dt` FROM `omni_dump`.`events`) `p_events_mv`",
			partitioned: true,
			want:        "SELECT `events`.`id` AS `it's`, `events`.`dt` FROM `omni_dump`.`events`",
		},
		{
			name:        "a partitioned view whose expansion quotes a keyword",
			original:    "SELECT id AS select, dt FROM omni_dump.events",
			expanded:    "SELECT `id`, `dt` FROM (SELECT `events`.`id` AS `select`, `events`.`dt` FROM `omni_dump`.`events`) `p_events_mv`",
			partitioned: true,
			want:        "SELECT `events`.`id` AS `select`, `events`.`dt` FROM `omni_dump`.`events`",
		},
		{
			name:        "a partitioned view whose expansion quotes a number",
			original:    "SELECT id AS 123, dt FROM omni_dump.events",
			expanded:    "SELECT `id`, `dt` FROM (SELECT `events`.`id` AS `123`, `events`.`dt` FROM `omni_dump`.`events`) `p_events_mv`",
			partitioned: true,
			want:        "SELECT `events`.`id` AS `123`, `events`.`dt` FROM `omni_dump`.`events`",
		},
		{
			name:        "a partitioned view whose expansion quotes a name beginning with an underscore",
			original:    "SELECT id AS _id, dt FROM omni_dump.events",
			expanded:    "SELECT `id`, `dt` FROM (SELECT `events`.`id` AS `_id`, `events`.`dt` FROM `omni_dump`.`events`) `p_events_mv`",
			partitioned: true,
			want:        "SELECT `events`.`id` AS `_id`, `events`.`dt` FROM `omni_dump`.`events`",
		},
		{
			name:        "a partitioned view whose expansion quotes a name holding a dollar sign",
			original:    "SELECT id AS a$b, dt FROM omni_dump.events",
			expanded:    "SELECT `id`, `dt` FROM (SELECT `events`.`id` AS `a$b`, `events`.`dt` FROM `omni_dump`.`events`) `p_events_mv`",
			partitioned: true,
			want:        "SELECT `events`.`id` AS `a$b`, `events`.`dt` FROM `omni_dump`.`events`",
		},
		{
			name:        "a partitioned view whose expansion quotes a name beginning with a digit",
			original:    "SELECT id AS 1bad, dt FROM omni_dump.events",
			expanded:    "SELECT `id`, `dt` FROM (SELECT `events`.`id` AS `1bad`, `events`.`dt` FROM `omni_dump`.`events`) `p_events_mv`",
			partitioned: true,
			// Hive's identifier begins with a letter or a digit, and the corpus
			// creates a view named that way from the text it states.
			want: "SELECT id AS 1bad, dt FROM omni_dump.events",
		},
		{
			name:        "a partitioned view whose expansion is no projection",
			original:    partitionedOriginal,
			expanded:    "SELECT `id`, `dt` FROM `omni_dump`.`events`",
			partitioned: true,
			want:        partitionedOriginal,
		},
		{
			name:        "a partitioned view whose literal holds a backquoted name",
			original:    "SELECT id, dt FROM omni_dump.events WHERE body = \"`x y`\"",
			expanded:    "SELECT `id`, `dt` FROM (SELECT `events`.`id`, `events`.`dt` FROM `omni_dump`.`events` WHERE `events`.`body` = \"`x y`\") `p_events_mv`",
			partitioned: true,
			want:        "SELECT id, dt FROM omni_dump.events WHERE body = \"`x y`\"",
		},
		{
			name:        "a projection over a query of its own",
			expanded:    "SELECT `a`, `dt` FROM (SELECT `x`.`a`, `x`.`dt` FROM (SELECT `b` AS `a`, `c` AS `dt` FROM `db`.`t`) `x`) `mv`",
			partitioned: true,
			want:        "SELECT `x`.`a`, `x`.`dt` FROM (SELECT `b` AS `a`, `c` AS `dt` FROM `db`.`t`) `x`",
		},
		{
			name:        "a query holding a parenthesis in a double-quoted literal",
			expanded:    "SELECT `id`, `dt` FROM (SELECT `id`, `dt` FROM `db`.`t` WHERE `body` = \")\") `mv`",
			partitioned: true,
			want:        "SELECT `id`, `dt` FROM `db`.`t` WHERE `body` = \")\"",
		},
		{
			name:        "a query holding a parenthesis in a literal",
			expanded:    "SELECT `id`, `dt` FROM (SELECT `id`, `dt` FROM `db`.`t` WHERE `note` = 'a ) note') `mv`",
			partitioned: true,
			want:        "SELECT `id`, `dt` FROM `db`.`t` WHERE `note` = 'a ) note'",
		},
	} {
		require.Equal(t, query.want, dumpMaterializedViewQuery(query.original, query.expanded, query.partitioned), query.name)
	}
}

// TestDumpEscapesLineBreak holds a value Hive states with a line break in it: the
// literal escapes the break, so a restore reaches the statement whole instead of
// the splitter cutting it at the value's own semicolon.
func TestDumpEscapesLineBreak(t *testing.T) {
	// Hive names the partition by its directory and percent-encodes the break.
	spec, ok := dumpPartitionSpec("dt=a%3B%0Ab")
	require.True(t, ok)
	require.Equal(t, "`dt`='a;\\nb'", spec)

	script := "ALTER TABLE `db`.`t` ADD PARTITION (" + spec + ") LOCATION 'file:/tmp/db/t/dt=a%3B%0Ab';\n"
	commands, err := splitCommands(script)
	require.NoError(t, err)
	require.Len(t, commands, 1)
	require.Equal(t, strings.TrimSuffix(script, "\n"), commands[0].Text)
}

// TestDumpStorageDiff holds what a partition's storage states that the table's
// does not: the classes a restore cannot carry, which ADD PARTITION takes from
// the table.
func TestDumpStorageDiff(t *testing.T) {
	const (
		parquetSerde = "org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe"
		parquetIn    = "org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat"
		parquetOut   = "org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat"
		orcSerde     = "org.apache.hadoop.hive.ql.io.orc.OrcSerde"
		orcIn        = "org.apache.hadoop.hive.ql.io.orc.OrcInputFormat"
		orcOut       = "org.apache.hadoop.hive.ql.io.orc.OrcOutputFormat"
	)
	storage := func(rows ...[3]string) dumpStorage {
		var reader dumpStorageReader
		for _, row := range rows {
			reader.read(row[0], row[1], row[2])
		}
		return reader.storage
	}
	table := storage(
		[3]string{"SerDe Library:", parquetSerde},
		[3]string{"InputFormat:", parquetIn},
		[3]string{"OutputFormat:", parquetOut},
		// A row the storage section states that is not one of the three.
		[3]string{"Compressed:", "No"},
		[3]string{"Storage Desc Params:", "", ""},
		[3]string{"", "serialization.format", "1"},
	)
	require.Equal(t, "", dumpStorageDiff(table, table), "a partition reading as the table does")
	require.Equal(t, "SerDe Library "+orcSerde, dumpStorageDiff(table, storage(
		[3]string{"SerDe Library:", orcSerde},
		[3]string{"InputFormat:", parquetIn},
		[3]string{"OutputFormat:", parquetOut},
		[3]string{"Compressed:", "No"},
		[3]string{"Storage Desc Params:", "", ""},
		[3]string{"", "serialization.format", "1"},
	)))
	orc := storage(
		[3]string{"SerDe Library:", orcSerde},
		[3]string{"InputFormat:", orcIn},
		[3]string{"OutputFormat:", orcOut},
		[3]string{"Compressed:", "No"},
		[3]string{"Storage Desc Params:", "", ""},
		[3]string{"", "serialization.format", "1"},
	)
	require.Equal(t,
		"SerDe Library "+orcSerde+", InputFormat "+orcIn+", OutputFormat "+orcOut,
		dumpStorageDiff(table, orc),
	)

	// A partition whose files sit under the table's own layout another way.
	require.Equal(t,
		"Skewed Columns [dt], Skewed Values [[2024-01-01]], Skewed true, Stored As SubDirectories Yes",
		dumpStorageDiff(table, storage(
			[3]string{"SerDe Library:", parquetSerde},
			[3]string{"InputFormat:", parquetIn},
			[3]string{"OutputFormat:", parquetOut},
			[3]string{"Compressed:", "No"},
			[3]string{"Skewed:", "true"},
			[3]string{"Skewed Columns:", "[dt]"},
			[3]string{"Skewed Values:", "[[2024-01-01]]"},
			[3]string{"Stored As SubDirectories:", "Yes"},
			[3]string{"Storage Desc Params:", "", ""},
			[3]string{"", "serialization.format", "1"},
		)),
		"a layout of the partition's own")

	// A table whose serde is the same and splits its files another way.
	csv := storage(
		[3]string{"SerDe Library:", parquetSerde},
		[3]string{"InputFormat:", parquetIn},
		[3]string{"OutputFormat:", parquetOut},
		[3]string{"Compressed:", "No"},
		[3]string{"Storage Desc Params:", "", ""},
		[3]string{"", "field.delim", "|"},
		[3]string{"", "serialization.format", "1"},
	)
	require.Equal(t, `SerDe property "serialization.format"=""`, dumpStorageDiff(csv, storage(
		[3]string{"SerDe Library:", parquetSerde},
		[3]string{"InputFormat:", parquetIn},
		[3]string{"OutputFormat:", parquetOut},
		[3]string{"Compressed:", "No"},
		[3]string{"Storage Desc Params:", "", ""},
		[3]string{"", "field.delim", "|"},
	)), "a parameter the table states and the partition does not")
	require.Equal(t,
		`SerDe property "field.delim"=","`,
		dumpStorageDiff(table, storage(
			[3]string{"SerDe Library:", parquetSerde},
			[3]string{"InputFormat:", parquetIn},
			[3]string{"OutputFormat:", parquetOut},
			[3]string{"Compressed:", "No"},
			[3]string{"Storage Desc Params:", "", ""},
			[3]string{"", "field.delim", ","},
			[3]string{"", "serialization.format", "1"},
		)),
		"a parameter the partition states another way")

	// A partition reading its files as the table does is the statement the dump
	// writes; one that states its own storage is what the dump refuses instead of
	// registering it under the table's.
	stmt, err := dumpPartitionStatement("db", "t", "dt=2024-01-01", "`dt`='2024-01-01'", "file:/tmp/db/t/dt=2024-01-01", table, table)
	require.NoError(t, err)
	require.Equal(t,
		"ALTER TABLE `db`.`t` ADD PARTITION (`dt`='2024-01-01') LOCATION 'file:/tmp/db/t/dt=2024-01-01'", stmt)
	_, err = dumpPartitionStatement("db", "t", "dt=2024-01-01", "`dt`='2024-01-01'", "file:/tmp/db/t/dt=2024-01-01", table, orc)
	require.EqualError(t, err,
		"partition dt=2024-01-01 of db.t states its own SerDe Library "+orcSerde+
			", InputFormat "+orcIn+", OutputFormat "+orcOut+", which ADD PARTITION does not carry")
}

// TestDumpUnescapeDescribed holds a value DESCRIBE FORMATTED renders: the
// formatter escapes a parameter Java's way, so what it states is not what the
// server stored.
func TestDumpUnescapeDescribed(t *testing.T) {
	for _, value := range []struct{ name, shown, want string }{
		{name: "no escape", shown: "the users' rows", want: "the users' rows"},
		{name: "an escaped backslash", shown: `a\\b`, want: `a\b`},
		{name: "an escaped quote", shown: `{\"a\":\"b\"}`, want: `{"a":"b"}`},
		{name: "an escaped line break", shown: `a\nb`, want: "a\nb"},
		{name: "an escape the formatter does not write", shown: `a\qb`, want: `a\qb`},
	} {
		require.Equal(t, value.want, dumpUnescapeDescribed(value.shown), value.name)
	}
	// A comment the dump states reads back as the one the server stored.
	require.Equal(t, "'a\\\\b'", dumpString(dumpUnescapeDescribed(`a\\b`)))
}

// TestDumpQueryRow holds how a query the server states over several rows reads
// back: each row is a line, a line the query holds empty is one, and the rows
// that end the section state none.
func TestDumpQueryRow(t *testing.T) {
	null := sql.NullString{}
	for _, row := range []struct {
		name              string
		label             string
		dataType, comment sql.NullString
		want              string
		wantOK            bool
	}{
		{
			name:     "the line a row states",
			dataType: sql.NullString{String: "SELECT 1bad         ", Valid: true},
			want:     "SELECT 1bad",
			wantOK:   true,
		},
		{
			name:     "the line that falls in the cell beside it",
			dataType: sql.NullString{String: strings.Repeat(" ", 21), Valid: true},
			comment:  sql.NullString{String: "FROM omni_dump.z_users", Valid: true},
			want:     "FROM omni_dump.z_users",
			wantOK:   true,
		},
		{
			name:     "a line the query holds empty",
			dataType: sql.NullString{String: strings.Repeat(" ", 21), Valid: true},
			want:     "",
			wantOK:   true,
		},
		{
			name:     "a line holding an indent",
			dataType: sql.NullString{String: "    b", Valid: true},
			want:     "    b",
			wantOK:   true,
		},
		{
			name:     "the row that ends the section",
			dataType: null,
			comment:  null,
			want:     "",
			wantOK:   false,
		},
		{
			name:     "a row stating another field",
			label:    "Rewrite Enabled:",
			dataType: sql.NullString{String: "No", Valid: true},
			want:     "",
			wantOK:   false,
		},
	} {
		line, ok := dumpQueryRow(row.label, row.dataType, row.comment)
		require.Equal(t, row.wantOK, ok, row.name)
		require.Equal(t, row.want, line, row.name)
	}
}
