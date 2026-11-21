package driver

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

func TestDump(t *testing.T) {
	m := readGolden(t, filepath.Join("testdata", "sync", "default.json"))
	before := proto.Clone(m)
	var b strings.Builder
	if err := (&Driver{}).Dump(context.Background(), &b, m); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(before, m) {
		t.Errorf("Dump mutated its input metadata")
	}
	golden := filepath.Join("testdata", "dump", "default.sql")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			t.Fatal(err)
		}
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
		t.Errorf("Dump differs from %s (-want +got):\n%s", golden, diff)
	}
}

const gapMarker = "-- not restorable: TestDumpGaps"

func restorable(corpus string) string {
	var kept []string
	for _, line := range strings.Split(corpus, "\n") {
		if !strings.HasSuffix(line, gapMarker) {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

func TestDumpRestores(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	source, restored, err := dumpAndRestore(t, "omni_dump_roundtrip", restorable(string(corpus)))
	if err != nil {
		t.Fatalf("restoring dump: %v", err)
	}
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("restored database differs from source (-source +restored):\n%s", diff)
	}
}

func restoreDiff(source, restored *metadatapb.DatabaseSchemaMetadata) string {
	return cmp.Diff(source, restored, protocmp.Transform(),
		protocmp.IgnoreFields(&metadatapb.DatabaseSchemaMetadata{}, "name"),
		protocmp.IgnoreFields(&metadatapb.TableMetadata{}, "data_size", "index_size", "data_free", "row_count"),
	)
}

func dumpAndRestore(t *testing.T, dbName, setup string) (*metadatapb.DatabaseSchemaMetadata, *metadatapb.DatabaseSchemaMetadata, error) {
	t.Helper()
	ctx := context.Background()
	restoredName := dbName + "_restored"
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{
			"DROP DATABASE IF EXISTS " + quoteBacktick(dbName),
			"DROP DATABASE IF EXISTS " + quoteBacktick(restoredName),
			"CREATE DATABASE " + quoteBacktick(dbName),
			"CREATE DATABASE " + quoteBacktick(restoredName),
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.OpenWith(t, testcontainer.Options{
		Database:        dbName,
		MultiStatements: true,
	})
	t.Cleanup(func() { _ = db.Close() })

	src := &Driver{db: db, databaseName: dbName}
	execCorpus(t, src.db, setup)

	source, err := src.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	sourceColumns := materializedViewColumnTypes(t, db, dbName)
	var dumpBuf strings.Builder
	if err := src.Dump(ctx, &dumpBuf, source); err != nil {
		t.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, "USE "+quoteBacktick(restoredName)); err != nil {
		t.Fatal(err)
	}
	// Replace source-database qualifiers in view bodies so views bind to the
	// restored database rather than the source database.
	rewritten := strings.ReplaceAll(dumpBuf.String(), quoteBacktick(dbName)+".", quoteBacktick(restoredName)+".")
	for _, stmt := range splitLikeMySQLClient(rewritten) {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return source, nil, err
		}
	}

	dst := &Driver{db: db, databaseName: restoredName}
	restored, err := dst.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range restored.Schemas[0].Views {
		v.Definition = strings.ReplaceAll(v.Definition, quoteBacktick(restoredName)+".", quoteBacktick(dbName)+".")
	}
	for _, v := range restored.Schemas[0].MaterializedViews {
		v.Definition = strings.ReplaceAll(v.Definition, quoteBacktick(restoredName)+".", quoteBacktick(dbName)+".")
	}
	// The snapshot carries no materialized view's columns, so the restored container's own
	// are compared here: a view materialized from a temporary placeholder takes the
	// placeholder's literal types instead of the real view's.
	if diff := cmp.Diff(sourceColumns, materializedViewColumnTypes(t, db, restoredName)); diff != "" {
		t.Errorf("restored materialized view columns differ from source (-source +restored):\n%s", diff)
	}
	return source, restored, nil
}

// materializedViewColumnTypes reads the column type of every materialized view's container
// table in dbName, keyed by view and column, so a container's storage is compared across a
// round trip that the snapshot cannot compare.
func materializedViewColumnTypes(t *testing.T, db *sql.DB, dbName string) map[string]string {
	t.Helper()
	const query = `
		SELECT m.MVIEW_NAME, c.COLUMN_NAME, c.COLUMN_TYPE
		FROM oceanbase.DBA_MVIEWS m
		JOIN information_schema.COLUMNS c
		ON c.TABLE_SCHEMA = m.OWNER AND c.TABLE_NAME = m.CONTAINER_NAME
		WHERE m.OWNER = ?
		ORDER BY m.MVIEW_NAME, c.ORDINAL_POSITION`
	rows, err := db.QueryContext(context.Background(), query, dbName)
	if err != nil {
		t.Fatalf("materialized view columns of %s: %v", dbName, err)
	}
	defer rows.Close()
	columns := make(map[string]string)
	for rows.Next() {
		var view, column, columnType string
		if err := rows.Scan(&view, &column, &columnType); err != nil {
			t.Fatal(err)
		}
		columns[view+"."+column] = columnType
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return columns
}

func splitLikeMySQLClient(script string) []string {
	const directive = "DELIMITER "
	var stmts []string
	var cur strings.Builder
	flush := func() {
		if stmt := strings.TrimSpace(cur.String()); stmt != "" {
			stmts = append(stmts, stmt)
		}
		cur.Reset()
	}
	delimiter := ";"
	var quote byte
	for i := 0; i < len(script); {
		c := script[i]
		if quote != 0 {
			cur.WriteByte(c)
			switch {
			case c == '\\' && quote != '`' && i+1 < len(script):
				cur.WriteByte(script[i+1])
				i += 2
				continue
			case c == quote:
				quote = 0
			}
			i++
			continue
		}
		lineEnd := strings.IndexByte(script[i:], '\n')
		if lineEnd < 0 {
			lineEnd = len(script) - i
		}
		switch {
		case strings.TrimSpace(cur.String()) == "" && len(script)-i >= len(directive) && strings.EqualFold(script[i:i+len(directive)], directive):
			delimiter = strings.TrimSpace(script[i+len(directive) : i+lineEnd])
			cur.Reset()
			i += lineEnd
		case strings.HasPrefix(script[i:], delimiter):
			flush()
			i += len(delimiter)
		case c == '#' || strings.HasPrefix(script[i:], "--") && (i+2 == len(script) || script[i+2] <= ' '):
			i += lineEnd
		case strings.HasPrefix(script[i:], "/*") && !strings.HasPrefix(script[i:], "/*!"):
			end := strings.Index(script[i+2:], "*/")
			if end < 0 {
				i = len(script)
			} else {
				i += end + 4
			}
		default:
			if c == '\'' || c == '"' || c == '`' {
				quote = c
			}
			cur.WriteByte(c)
			i++
		}
	}
	flush()
	return stmts
}

func TestDumpRestoresCheckNotEnforcedAndViewOptions(t *testing.T) {
	const setup = `CREATE TABLE t_chk (
	id INT PRIMARY KEY,
	x INT NOT NULL,
	CONSTRAINT chk_enf CHECK (x > 0),
	CONSTRAINT chk_ne CHECK (x < 10) NOT ENFORCED
);

CREATE MATERIALIZED VIEW LOG ON t_chk;

CREATE VIEW v_opts AS
	SELECT id, x FROM t_chk WHERE x > 0 WITH LOCAL CHECK OPTION;`
	source, restored, err := dumpAndRestore(t, "omni_ob_options_roundtrip", setup)
	if err != nil {
		t.Fatalf("restoring OceanBase dump: %v", err)
	}
	if len(source.Schemas[0].Tables) != 1 || len(source.Schemas[0].Tables[0].CheckConstraints) != 2 {
		t.Fatalf("unexpected OceanBase check constraints: %+v", source.Schemas[0].Tables)
	}
	tbl := source.Schemas[0].Tables[0]
	if tbl.HasMaterializedViewLog == nil || !*tbl.HasMaterializedViewLog {
		t.Fatalf("expected HasMaterializedViewLog=true on t_chk, got %v", tbl.HasMaterializedViewLog)
	}
	chk := tbl.CheckConstraints
	if chk[0].NotEnforced || !chk[1].NotEnforced {
		t.Fatalf("unexpected OceanBase NotEnforced flags: %+v", chk)
	}
	if len(source.Schemas[0].Views) != 1 {
		t.Fatalf("unexpected OceanBase views: %+v", source.Schemas[0].Views)
	}
	v := source.Schemas[0].Views[0]
	if v.CheckOption != "LOCAL" {
		t.Fatalf("unexpected OceanBase view CheckOption: %+v", v)
	}
	if strings.Contains(strings.ToUpper(v.Definition), "CHECK OPTION") {
		t.Fatalf("expected normalized Definition without CHECK OPTION, got %q", v.Definition)
	}
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("restored OceanBase database differs from source (-source +restored):\n%s", diff)
	}
}

// TestDumpRestoresColumnGroupsAndVectorIndexes pins what only SHOW CREATE TABLE
// states on oceanbase-ce 4.3.5.6: a table's column groups and a vector index,
// which information_schema.STATISTICS lists over hidden columns beside the
// hidden tables that support it.
func TestDumpRestoresColumnGroupsAndVectorIndexes(t *testing.T) {
	const setup = `CREATE TABLE t_cg (a INT, b INT) WITH COLUMN GROUP(all columns, each column);

CREATE TABLE t_each (a INT) WITH COLUMN GROUP(each column);

CREATE TABLE t_vec (id INT PRIMARY KEY, v VECTOR(3), VECTOR INDEX vi(v) WITH (distance=l2, type=hnsw));`
	source, restored, err := dumpAndRestore(t, "omni_ob_show_create_roundtrip", setup)
	if err != nil {
		t.Fatalf("restoring OceanBase dump: %v", err)
	}
	tables := map[string]*metadatapb.TableMetadata{}
	for _, tbl := range source.Schemas[0].Tables {
		tables[tbl.Name] = tbl
	}
	if got := tables["t_cg"].ColumnGroups; got != "all columns, each column" {
		t.Errorf("t_cg column groups = %q", got)
	}
	if got := tables["t_each"].ColumnGroups; got != "each column" {
		t.Errorf("t_each column groups = %q", got)
	}
	var names []string
	for _, idx := range tables["t_vec"].Indexes {
		names = append(names, idx.Name)
		if idx.Name == "vi" {
			if !idx.IsVector || len(idx.Expressions) != 1 || idx.Expressions[0] != "v" ||
				strings.Join(idx.VectorParams, ",") != "DISTANCE=L2,TYPE=HNSW,LIB=VSAG,M=16,EF_CONSTRUCTION=200,EF_SEARCH=64" {
				t.Errorf("vector index vi = %+v", idx)
			}
		}
	}
	if strings.Join(names, ",") != "PRIMARY,vi" {
		t.Errorf("t_vec indexes = %v, want PRIMARY and vi alone", names)
	}
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("restored OceanBase database differs from source (-source +restored):\n%s", diff)
	}
}
