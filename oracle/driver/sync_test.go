package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	metadatapb "github.com/bytebase/omni/metadata"
	oracleanalysis "github.com/bytebase/omni/oracle/analysis"
	"github.com/bytebase/omni/oracle/internal/testcontainer"
	oracleparser "github.com/bytebase/omni/oracle/parser"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync from this run")

const (
	// testSchema is the schema the corpus fills. It is a user of the pluggable
	// database, which is what names a schema in Oracle.
	testSchema = "OMNI_APP"
	// testEngineVersion is the release of the pinned server.
	testEngineVersion = "23.26.3.0.0"
)

// openTestDriver returns a driver that reads testSchema, the schema the test's
// corpus fills. It connects as SYSTEM, which owns the dictionary the sync reads.
func openTestDriver(t *testing.T) *Driver {
	t.Helper()
	db := testcontainer.Open(t, "system")
	t.Cleanup(func() { _ = db.Close() })
	return &Driver{db: db, databaseName: testSchema, engineVersion: testEngineVersion}
}

// execCorpus applies a script of Oracle statements one at a time on one
// session, which is the session the tape holds.
func execCorpus(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	for _, statement := range oracleparser.Split(script) {
		if statement.Empty() {
			continue
		}
		if _, err := conn.ExecContext(ctx, statement.Text); err != nil {
			t.Fatalf("%s:\n%v", statement.Text, err)
		}
	}
}

func TestSyncDatabase(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openTestDriver(t)
	execCorpus(t, d.db, string(corpus))

	got, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	golden := filepath.Join("testdata", "sync", "database.json")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			t.Fatal(err)
		}
		b, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, append(b, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want := readGolden(t, golden)
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("SyncDatabase differs from %s (-want +got):\n%s", golden, diff)
	}
}

func readGolden(t *testing.T, path string) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m := &metadatapb.DatabaseSchemaMetadata{}
	if err := protojson.Unmarshal(b, m); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return m
}

func TestSyncInstance(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openTestDriver(t)
	execCorpus(t, d.db, string(corpus))

	instance, err := d.SyncInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// The server's own schemas and whatever an earlier test left differ between
	// runs, so the golden holds the version and this test's schema.
	var names []string
	for _, database := range instance.Databases {
		if database.Name == testSchema {
			names = append(names, database.Name)
		}
	}
	out, err := json.MarshalIndent(map[string]any{"version": instance.Version, "databases": names}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got := string(out) + "\n"

	golden := filepath.Join("testdata", "sync", "instance.json")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(want), got); diff != "" {
		t.Errorf("SyncInstance differs from %s (-want +got):\n%s", golden, diff)
	}
}

// TestSyncObjectAndXMLTypeTables pins that an object table and an XMLType table
// are read, and marked as having no DDL to dump: their definition is their type.
func TestSyncObjectAndXMLTypeTables(t *testing.T) {
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE, CREATE TYPE TO OMNI_APP;
CREATE TYPE OMNI_APP.TY AS OBJECT (X NUMBER, Y VARCHAR2(10));
CREATE TABLE OMNI_APP.OT OF OMNI_APP.TY (CONSTRAINT OT_PK PRIMARY KEY (X));
CREATE INDEX OMNI_APP.OT_Y ON OMNI_APP.OT (Y);
CREATE TABLE OMNI_APP.XT OF XMLTYPE;
CREATE TABLE OMNI_APP.REL (A NUMBER);
`)
	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	schema := onlySchema(t, metadata)
	tables := make(map[string]*metadatapb.TableMetadata, len(schema.Tables))
	for _, table := range schema.Tables {
		tables[table.Name] = table
	}
	for _, name := range []string{"OT", "REL", "XT"} {
		if tables[name] == nil {
			t.Fatalf("table %q is missing from %v", name, namesOfTables(schema.Tables))
		}
	}
	if tables["REL"].SkipDump {
		t.Error("a relational table was marked as having no DDL")
	}
	if !tables["OT"].SkipDump {
		t.Error("an object table was not marked as having no DDL of its own")
	}
	if !tables["XT"].SkipDump {
		t.Error("an XMLType table was not marked as having no DDL of its own")
	}
	if got, want := namesOfColumns(tables["OT"].Columns), []string{"X", "Y"}; !cmp.Equal(got, want) {
		t.Errorf("object table columns = %v, want %v", got, want)
	}
	if got, want := namesOfColumns(tables["XT"].Columns), []string{"SYS_NC_ROWINFO$"}; !cmp.Equal(got, want) {
		t.Errorf("XMLType table columns = %v, want %v", got, want)
	}
}

// TestSyncInvisibleColumns pins which columns the sync lists: the invisible ones
// and the virtual ones are listed, and the system columns an index or an unused
// column leaves behind are not.
func TestSyncInvisibleColumns(t *testing.T) {
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW TO OMNI_APP;
CREATE TABLE OMNI_APP.T (
	A NUMBER,
	H NUMBER INVISIBLE,
	B VARCHAR2(10),
	HV NUMBER INVISIBLE DEFAULT 7 NOT NULL,
	VC AS (A+1) VIRTUAL,
	IHV INVISIBLE AS (A*2) VIRTUAL,
	U NUMBER
);
CREATE INDEX OMNI_APP.T_FBI ON OMNI_APP.T (UPPER(B));
ALTER TABLE OMNI_APP.T SET UNUSED COLUMN U;
CREATE VIEW OMNI_APP.V (A, B INVISIBLE) AS SELECT A, B FROM OMNI_APP.T;
`)
	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	schema := onlySchema(t, metadata)
	if len(schema.Tables) != 1 {
		t.Fatalf("tables = %v", namesOfTables(schema.Tables))
	}
	type column struct {
		Position  int
		Invisible bool
		Nullable  bool
		Default   string
		Generated string
	}
	read := func(columns []*metadatapb.ColumnMetadata) []column {
		out := make([]column, len(columns))
		for i, c := range columns {
			out[i] = column{
				Position:  int(c.Position),
				Invisible: c.IsInvisible,
				Nullable:  c.Nullable,
				Default:   c.Default,
				Generated: c.GetGeneration().GetExpression(),
			}
		}
		return out
	}
	want := []column{
		{Position: 1, Nullable: true},
		{Position: 2, Nullable: true},
		{Position: 3, Nullable: true, Generated: `"A"+1`},
		{Position: 4, Invisible: true, Nullable: true},
		{Position: 5, Invisible: true, Default: "7"},
		{Position: 6, Invisible: true, Nullable: true, Generated: `"A"*2`},
	}
	if diff := cmp.Diff(want, read(schema.Tables[0].Columns)); diff != "" {
		t.Errorf("columns (-want +got):\n%s", diff)
	}
	if len(schema.Views) != 1 {
		t.Fatalf("views = %v", namesOfViews(schema.Views))
	}
	if diff := cmp.Diff([]column{
		{Position: 1, Nullable: true},
		{Position: 2, Invisible: true, Nullable: true},
	}, read(schema.Views[0].Columns)); diff != "" {
		t.Errorf("view columns (-want +got):\n%s", diff)
	}

	// Before 12.2 the columns come from a dictionary view that lists the virtual
	// ones separately.
	txn, err := d.db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = txn.Rollback() }()
	columnMap, err := getTableColumns(txn, testSchema, &oracleanalysis.Version{First: 12, Second: 1})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, read(columnMap[tableKey{Schema: testSchema, Table: "T"}])); diff != "" {
		t.Errorf("columns before 12.2 (-want +got):\n%s", diff)
	}
}

// TestSyncSchemaWithApostrophe pins that a schema whose quoted name holds an
// apostrophe reaches the dictionary queries as one literal.
func TestSyncSchemaWithApostrophe(t *testing.T) {
	const schema = "omni'owner"
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER "omni'owner" NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE TO "omni'owner";
CREATE TABLE "omni'owner".T (ID NUMBER NOT NULL, NOTE VARCHAR2(10));
`)
	d.databaseName = schema
	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatalf("a schema named with an apostrophe failed the sync: %v", err)
	}
	if metadata.Name != schema {
		t.Fatalf("database = %q, want %q", metadata.Name, schema)
	}
	if len(metadata.Schemas) != 1 || len(metadata.Schemas[0].Tables) != 1 {
		t.Fatalf("tables = %+v, want the one table", metadata.Schemas)
	}
	if table := metadata.Schemas[0].Tables[0]; table.Name != "T" || len(table.Columns) != 2 {
		t.Errorf("table = %+v, want T with its two columns", table)
	}
}

// onlySchema returns the one schema a database has in Oracle, whose name is the
// database's: a schema is a user, and the metadata type names it twice.
func onlySchema(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) *metadatapb.SchemaMetadata {
	t.Helper()
	if metadata.Name != testSchema {
		t.Fatalf("database = %q, want %q", metadata.Name, testSchema)
	}
	if len(metadata.Schemas) != 1 {
		t.Fatalf("schemas = %d, want 1", len(metadata.Schemas))
	}
	return metadata.Schemas[0]
}

func namesOfTables(tables []*metadatapb.TableMetadata) []string {
	names := make([]string, len(tables))
	for i, table := range tables {
		names[i] = table.Name
	}
	return names
}

func namesOfColumns(columns []*metadatapb.ColumnMetadata) []string {
	names := make([]string, len(columns))
	for i, column := range columns {
		names[i] = column.Name
	}
	return names
}

func namesOfViews(views []*metadatapb.ViewMetadata) []string {
	names := make([]string, len(views))
	for i, view := range views {
		names[i] = view.Name
	}
	return names
}
