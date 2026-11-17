package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
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
	identityColumns, err := getIdentityColumns(txn, testSchema, &oracleanalysis.Version{First: 12, Second: 1})
	if err != nil {
		t.Fatal(err)
	}
	columnMap, err := getTableColumns(txn, testSchema, &oracleanalysis.Version{First: 12, Second: 1}, identityColumns)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, read(columnMap[tableKey{Schema: testSchema, Table: "T"}])); diff != "" {
		t.Errorf("columns before 12.2 (-want +got):\n%s", diff)
	}
}

// TestSyncCrossSchemaForeignKeys pins that two foreign keys naming the same
// remote constraint each report its columns once.
func TestSyncCrossSchemaForeignKeys(t *testing.T) {
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
CREATE USER OMNI_OTHER NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE TO OMNI_APP;
GRANT CREATE SESSION, CREATE TABLE TO OMNI_OTHER;
CREATE TABLE OMNI_OTHER.R (ID NUMBER NOT NULL, PART NUMBER NOT NULL, CONSTRAINT R_PK PRIMARY KEY (ID, PART));
GRANT REFERENCES ON OMNI_OTHER.R TO OMNI_APP;
CREATE TABLE OMNI_APP.A (ID NUMBER NOT NULL, R_ID NUMBER NOT NULL, R_PART NUMBER NOT NULL, CONSTRAINT A_PK PRIMARY KEY (ID), CONSTRAINT A_R_FK FOREIGN KEY (R_ID, R_PART) REFERENCES OMNI_OTHER.R (ID, PART));
CREATE TABLE OMNI_APP.B (ID NUMBER NOT NULL, R_ID NUMBER NOT NULL, R_PART NUMBER NOT NULL, CONSTRAINT B_PK PRIMARY KEY (ID), CONSTRAINT B_R_FK FOREIGN KEY (R_ID, R_PART) REFERENCES OMNI_OTHER.R (ID, PART));
`)
	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	schema := onlySchema(t, metadata)
	if len(schema.Tables) != 2 {
		t.Fatalf("tables = %v, want A and B", namesOfTables(schema.Tables))
	}
	for _, table := range schema.Tables {
		if len(table.ForeignKeys) != 1 {
			t.Fatalf("%s foreign keys = %d, want 1", table.Name, len(table.ForeignKeys))
		}
		foreignKey := table.ForeignKeys[0]
		if foreignKey.ReferencedSchema != "OMNI_OTHER" || foreignKey.ReferencedTable != "R" {
			t.Errorf("%s references %s.%s, want OMNI_OTHER.R", table.Name, foreignKey.ReferencedSchema, foreignKey.ReferencedTable)
		}
		if diff := cmp.Diff([]string{"ID", "PART"}, foreignKey.ReferencedColumns); diff != "" {
			t.Errorf("%s referenced columns (-want +got):\n%s", table.Name, diff)
		}
	}
}

// TestSyncUnnamedCheckConstraint pins that a CHECK the user wrote without a name
// is kept, though Oracle named it SYS_C<n> the way it names a NOT NULL.
func TestSyncUnnamedCheckConstraint(t *testing.T) {
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE TO OMNI_APP;
CREATE TABLE OMNI_APP.C (A NUMBER CHECK (A > 0), B NUMBER NOT NULL);
`)
	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	schema := onlySchema(t, metadata)
	if len(schema.Tables) != 1 {
		t.Fatalf("tables = %v", namesOfTables(schema.Tables))
	}
	checks := schema.Tables[0].CheckConstraints
	if len(checks) != 1 {
		t.Fatalf("check constraints = %+v, want the user's one", checks)
	}
	if !strings.Contains(checks[0].Expression, "A") {
		t.Errorf("check constraint = %+v, want the user's condition", checks[0])
	}
}

// TestSyncDatabaseWithoutAConfiguredDatabase pins that a driver opened without
// Config.Database reads the schema the session's names resolve in.
func TestSyncDatabaseWithoutAConfiguredDatabase(t *testing.T) {
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE TO OMNI_APP;
CREATE TABLE OMNI_APP.T (ID NUMBER);
ALTER SESSION SET CURRENT_SCHEMA = OMNI_APP;
`)
	d.databaseName = ""
	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if metadata.Name != testSchema {
		t.Errorf("database = %q, want the session's %q", metadata.Name, testSchema)
	}
	if len(metadata.Schemas) != 1 || len(metadata.Schemas[0].Tables) != 1 || metadata.Schemas[0].Tables[0].Name != "T" {
		t.Errorf("tables = %+v, want the one table", metadata.Schemas)
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

// TestSyncCommentableObjects pins which objects of a schema can carry a comment.
// Oracle's COMMENT statement names a table (a view is commented through that
// form), a materialized view, and a column, and the sync reads all four. It
// names no index — a plain, function-based, partitioned, and constraint-backed
// index all answer ORA-32594 — and no sequence, routine, or trigger, so the
// sync has nothing to read for them. This test fails if a later release gains
// one of those categories.
func TestSyncCommentableObjects(t *testing.T) {
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE SEQUENCE, CREATE PROCEDURE, CREATE TRIGGER, CREATE MATERIALIZED VIEW TO OMNI_APP;
CREATE TABLE OMNI_APP.T (
	ID NUMBER NOT NULL,
	NAME VARCHAR2(20) NOT NULL,
	QTY NUMBER,
	CONSTRAINT T_PK PRIMARY KEY (ID),
	CONSTRAINT T_NAME_UK UNIQUE (NAME)
);
CREATE INDEX OMNI_APP.T_QTY ON OMNI_APP.T (QTY);
CREATE INDEX OMNI_APP.T_UPPER_NAME ON OMNI_APP.T (UPPER(NAME));
CREATE TABLE OMNI_APP.P (
	ID NUMBER,
	D DATE
) PARTITION BY RANGE (D) (
	PARTITION P1 VALUES LESS THAN (TO_DATE('2025-01-01', 'YYYY-MM-DD')),
	PARTITION P2 VALUES LESS THAN (MAXVALUE)
);
CREATE INDEX OMNI_APP.P_D ON OMNI_APP.P (D) LOCAL;
CREATE SEQUENCE OMNI_APP.S;
CREATE FUNCTION OMNI_APP.F RETURN NUMBER IS BEGIN RETURN 1; END;
CREATE TRIGGER OMNI_APP.TRG BEFORE INSERT ON OMNI_APP.T BEGIN NULL; END;
CREATE VIEW OMNI_APP.V AS SELECT ID, NAME FROM OMNI_APP.T;
CREATE MATERIALIZED VIEW OMNI_APP.MV AS SELECT ID FROM OMNI_APP.T;
COMMENT ON TABLE OMNI_APP.T IS 'table comment';
COMMENT ON COLUMN OMNI_APP.T.NAME IS 'column comment';
COMMENT ON TABLE OMNI_APP.V IS 'view comment';
COMMENT ON MATERIALIZED VIEW OMNI_APP.MV IS 'materialized view comment';
`)

	refuseComments(t, d.db,
		`COMMENT ON INDEX OMNI_APP.T_QTY IS 'a plain index'`,
		`COMMENT ON INDEX OMNI_APP.T_UPPER_NAME IS 'a function-based index'`,
		`COMMENT ON INDEX OMNI_APP.P_D IS 'a partitioned index'`,
		`COMMENT ON INDEX OMNI_APP.T_PK IS 'the primary key's index'`,
		`COMMENT ON SEQUENCE OMNI_APP.S IS 'a sequence'`,
		`COMMENT ON FUNCTION OMNI_APP.F IS 'a function'`,
		`COMMENT ON TRIGGER OMNI_APP.TRG IS 'a trigger'`,
		`COMMENT ON VIEW OMNI_APP.V IS 'a view'`,
	)

	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	schema := onlySchema(t, metadata)
	tables := make(map[string]*metadatapb.TableMetadata, len(schema.Tables))
	for _, table := range schema.Tables {
		tables[table.Name] = table
	}
	if got, want := tables["T"].Comment, "table comment"; got != want {
		t.Errorf("T comment = %q, want %q", got, want)
	}
	var nameColumn *metadatapb.ColumnMetadata
	for _, column := range tables["T"].Columns {
		if column.Name == "NAME" {
			nameColumn = column
		}
	}
	if nameColumn == nil {
		t.Fatalf("T columns = %v, want NAME", namesOfColumns(tables["T"].Columns))
	}
	if got, want := nameColumn.Comment, "column comment"; got != want {
		t.Errorf("T.NAME comment = %q, want %q", got, want)
	}
	if len(schema.Views) != 1 {
		t.Fatalf("views = %v, want V", namesOfViews(schema.Views))
	}
	if got, want := schema.Views[0].Comment, "view comment"; got != want {
		t.Errorf("V comment = %q, want %q", got, want)
	}
	if len(schema.MaterializedViews) != 1 {
		t.Fatalf("materialized views = %d, want MV", len(schema.MaterializedViews))
	}
	if got, want := schema.MaterializedViews[0].Comment, "materialized view comment"; got != want {
		t.Errorf("MV comment = %q, want %q", got, want)
	}

	// The four shapes of index carry no comment, and neither does anything else
	// Oracle names no category for.
	for _, table := range []string{"T", "P"} {
		if len(tables[table].Indexes) == 0 {
			t.Fatalf("%s indexes = %d, want the ones the schema made", table, len(tables[table].Indexes))
		}
		for _, index := range tables[table].Indexes {
			if index.Comment != "" {
				t.Errorf("%s index %s carries comment %q", table, index.Name, index.Comment)
			}
		}
	}
	if len(schema.Sequences) != 1 || schema.Sequences[0].Comment != "" {
		t.Errorf("sequences = %+v, want S alone with no comment", schema.Sequences)
	}
	if len(schema.Functions) != 1 || schema.Functions[0].Comment != "" {
		t.Errorf("functions = %+v, want F alone with no comment", schema.Functions)
	}
	if len(tables["T"].Triggers) != 1 || tables["T"].Triggers[0].Comment != "" {
		t.Errorf("T triggers = %+v, want TRG alone with no comment", tables["T"].Triggers)
	}
}

// refuseComments runs statements Oracle must answer ORA-32594: a COMMENT that
// names no category. A later release that gains one fails here, so the sync's
// silence about that comment is revisited.
func refuseComments(t *testing.T, db *sql.DB, statements ...string) {
	t.Helper()
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	for _, statement := range statements {
		switch _, err := conn.ExecContext(context.Background(), statement); {
		case err == nil:
			t.Errorf("%s: the server accepted a comment for an object it names no category for", statement)
		case !strings.Contains(err.Error(), "ORA-32594"):
			t.Errorf("%s: %v, want ORA-32594", statement, err)
		}
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
