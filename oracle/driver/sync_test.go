package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
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
	columnMap, _, err := getTableColumns(txn, testSchema, &oracleanalysis.Version{First: 12, Second: 1}, identityColumns)
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
	expectRefusal(t, db, "ORA-32594", statements...)
}

// expectRefusal runs statements the server must refuse with the code. Each
// statement is recorded, so a release that starts accepting one fails here.
func expectRefusal(t *testing.T, db *sql.DB, code string, statements ...string) {
	t.Helper()
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	for _, statement := range statements {
		switch _, err := conn.ExecContext(context.Background(), statement); {
		case err == nil:
			t.Errorf("%s: the server accepted a statement it must refuse with %s", statement, code)
		case !strings.Contains(err.Error(), code):
			t.Errorf("%s: %v, want %s", statement, err, code)
		}
	}
}

// TestSyncConstraintBackedIndexes pins that a primary key or unique constraint
// and the index that enforces it are two objects with two names, and that the
// sync lists the index as the dictionary does: its own name, type, uniqueness,
// visibility, and descending flags, with the constraint laid over it. Oracle
// reuses an index the schema already had for a constraint (UQ_NAME on IX_PLAIN),
// so naming the metadata index after the constraint lists an object the DBA
// cannot see, and reading type, visibility, and descending from nowhere leaves
// them empty.
func TestSyncConstraintBackedIndexes(t *testing.T) {
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE TO OMNI_APP;
CREATE TABLE OMNI_APP.PK_AUTO (ID NUMBER NOT NULL, CONSTRAINT PK_AUTO_PK PRIMARY KEY (ID));
CREATE TABLE OMNI_APP.UK_AUTO (ID NUMBER NOT NULL, CONSTRAINT UK_AUTO_UK UNIQUE (ID));
CREATE TABLE OMNI_APP.PK_REUSED (ID NUMBER NOT NULL);
CREATE INDEX OMNI_APP.PK_REUSED_IDX ON OMNI_APP.PK_REUSED (ID);
ALTER TABLE OMNI_APP.PK_REUSED ADD CONSTRAINT PK_REUSED_PK PRIMARY KEY (ID) USING INDEX OMNI_APP.PK_REUSED_IDX;
CREATE TABLE OMNI_APP.UK_REUSED (NAME VARCHAR2(20) NOT NULL);
CREATE INDEX OMNI_APP.UK_REUSED_IDX ON OMNI_APP.UK_REUSED (NAME);
ALTER TABLE OMNI_APP.UK_REUSED ADD CONSTRAINT UK_REUSED_UK UNIQUE (NAME) USING INDEX OMNI_APP.UK_REUSED_IDX;
CREATE TABLE OMNI_APP.PK_INVISIBLE (ID NUMBER NOT NULL);
CREATE INDEX OMNI_APP.PK_INVISIBLE_IDX ON OMNI_APP.PK_INVISIBLE (ID) INVISIBLE;
ALTER TABLE OMNI_APP.PK_INVISIBLE ADD CONSTRAINT PK_INVISIBLE_PK PRIMARY KEY (ID) USING INDEX OMNI_APP.PK_INVISIBLE_IDX;
CREATE TABLE OMNI_APP.UK_VIRTUAL (EMAIL VARCHAR2(40) NOT NULL, UPPER_EMAIL AS (UPPER(EMAIL)));
CREATE UNIQUE INDEX OMNI_APP.UK_VIRTUAL_IDX ON OMNI_APP.UK_VIRTUAL (UPPER(EMAIL));
ALTER TABLE OMNI_APP.UK_VIRTUAL ADD CONSTRAINT UK_VIRTUAL_UK UNIQUE (UPPER_EMAIL) USING INDEX OMNI_APP.UK_VIRTUAL_IDX;
CREATE TABLE OMNI_APP.UK_UNNAMED (A NUMBER UNIQUE);
CREATE TABLE OMNI_APP.PK_DISABLED (A NUMBER);
ALTER TABLE OMNI_APP.PK_DISABLED ADD CONSTRAINT PK_DISABLED_PK PRIMARY KEY (A) DISABLE;
CREATE TABLE OMNI_APP.SHARED (A NUMBER NOT NULL, B NUMBER NOT NULL);
CREATE INDEX OMNI_APP.SHARED_IDX ON OMNI_APP.SHARED (A, B);
ALTER TABLE OMNI_APP.SHARED ADD CONSTRAINT SHARED_PK PRIMARY KEY (A, B) USING INDEX OMNI_APP.SHARED_IDX;
ALTER TABLE OMNI_APP.SHARED ADD CONSTRAINT SHARED_UK UNIQUE (A) USING INDEX OMNI_APP.SHARED_IDX;
CREATE TABLE OMNI_APP.PK_DESC (ID NUMBER NOT NULL);
CREATE INDEX OMNI_APP.PK_DESC_IDX ON OMNI_APP.PK_DESC (ID DESC);
CREATE TABLE OMNI_APP.UK_EXPRESSION (EMAIL VARCHAR2(40));
CREATE UNIQUE INDEX OMNI_APP.UK_EXPRESSION_IDX ON OMNI_APP.UK_EXPRESSION (UPPER(EMAIL));
`)
	// Two shapes a constraint-backed index cannot take, so the sync never meets
	// them: a DESC key is stored in a hidden column, which is not the column the
	// constraint names (ORA-14196), and a constraint names a column, not an
	// expression (ORA-03059). The index of a virtual column's constraint is the
	// function-based index on the expression instead.
	expectRefusal(t, d.db, "ORA-14196", `ALTER TABLE OMNI_APP.PK_DESC ADD CONSTRAINT PK_DESC_PK PRIMARY KEY (ID) USING INDEX OMNI_APP.PK_DESC_IDX`)
	expectRefusal(t, d.db, "ORA-03059", `ALTER TABLE OMNI_APP.UK_EXPRESSION ADD CONSTRAINT UK_EXPRESSION_UK UNIQUE (UPPER(EMAIL)) USING INDEX OMNI_APP.UK_EXPRESSION_IDX`)

	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	schema := onlySchema(t, metadata)
	constraint := func(name string, primary bool, columns ...string) *metadatapb.IndexConstraintMetadata {
		return &metadatapb.IndexConstraintMetadata{Name: name, Columns: columns, Primary: primary}
	}
	want := map[string][]*metadatapb.IndexMetadata{
		// Oracle made the index for the constraint, so the two names are one.
		"PK_AUTO": {{Name: "PK_AUTO_PK", Expressions: []string{`"ID"`}, Descending: []bool{false}, Type: "NORMAL", Unique: true, Primary: true, Visible: true, IsConstraint: true, Constraints: []*metadatapb.IndexConstraintMetadata{constraint("PK_AUTO_PK", true, `"ID"`)}}},
		"UK_AUTO": {{Name: "UK_AUTO_UK", Expressions: []string{`"ID"`}, Descending: []bool{false}, Type: "NORMAL", Unique: true, Visible: true, IsConstraint: true, Constraints: []*metadatapb.IndexConstraintMetadata{constraint("UK_AUTO_UK", false, `"ID"`)}}},
		// A reused index keeps its own name, its own non-uniqueness, and is what
		// the table's indexes list; the constraint is named on it.
		"PK_REUSED":    {{Name: "PK_REUSED_IDX", Expressions: []string{`"ID"`}, Descending: []bool{false}, Type: "NORMAL", Primary: true, Visible: true, IsConstraint: true, Constraints: []*metadatapb.IndexConstraintMetadata{constraint("PK_REUSED_PK", true, `"ID"`)}}},
		"UK_REUSED":    {{Name: "UK_REUSED_IDX", Expressions: []string{`"NAME"`}, Descending: []bool{false}, Type: "NORMAL", Visible: true, IsConstraint: true, Constraints: []*metadatapb.IndexConstraintMetadata{constraint("UK_REUSED_UK", false, `"NAME"`)}}},
		"PK_INVISIBLE": {{Name: "PK_INVISIBLE_IDX", Expressions: []string{`"ID"`}, Descending: []bool{false}, Type: "NORMAL", Primary: true, IsConstraint: true, Constraints: []*metadatapb.IndexConstraintMetadata{constraint("PK_INVISIBLE_PK", true, `"ID"`)}}},
		"UK_VIRTUAL":   {{Name: "UK_VIRTUAL_IDX", Expressions: []string{`"UPPER_EMAIL"`}, Descending: []bool{false}, Type: "FUNCTION-BASED NORMAL", Unique: true, Visible: true, IsConstraint: true, Constraints: []*metadatapb.IndexConstraintMetadata{constraint("UK_VIRTUAL_UK", false, `"UPPER_EMAIL"`)}}},
		// A disabled constraint has no index: the constraint is what is left.
		"PK_DISABLED": {{Name: "PK_DISABLED_PK", Expressions: []string{`"A"`}, Unique: true, Primary: true, IsConstraint: true, Constraints: []*metadatapb.IndexConstraintMetadata{constraint("PK_DISABLED_PK", true, `"A"`)}}},
		// One index enforces two constraints, which Oracle shares between them:
		// the unique constraint's keys are a subset of the index's.
		"SHARED": {{Name: "SHARED_IDX", Expressions: []string{`"A"`, `"B"`}, Descending: []bool{false, false}, Type: "NORMAL", Primary: true, Visible: true, IsConstraint: true, Constraints: []*metadatapb.IndexConstraintMetadata{constraint("SHARED_PK", true, `"A"`, `"B"`), constraint("SHARED_UK", false, `"A"`)}}},
		// Neither index is a constraint's, so both stay the schema's own. A DESC
		// index is FUNCTION-BASED NORMAL because Oracle stores the DESC key in a
		// hidden column, which is why no constraint can use one.
		"PK_DESC":       {{Name: "PK_DESC_IDX", Expressions: []string{`"ID"`}, Descending: []bool{true}, Type: "FUNCTION-BASED NORMAL", Visible: true}},
		"UK_EXPRESSION": {{Name: "UK_EXPRESSION_IDX", Expressions: []string{`UPPER("EMAIL")`}, Descending: []bool{false}, Type: "FUNCTION-BASED NORMAL", Unique: true, Visible: true}},
	}
	for table, indexes := range want {
		if diff := cmp.Diff(indexes, tableByName(t, schema, table).Indexes, protocmp.Transform()); diff != "" {
			t.Errorf("%s indexes (-want +got):\n%s", table, diff)
		}
	}
	// The index of a constraint the user did not name is named SYS_C<n>, and
	// the sync lists it though it skips the SYS_ indexes of LOBs and IOTs.
	unnamed := tableByName(t, schema, "UK_UNNAMED").Indexes
	if len(unnamed) != 1 {
		t.Fatalf("UK_UNNAMED indexes = %v, want the constraint's one", unnamed)
	}
	if index := unnamed[0]; !strings.HasPrefix(index.GetName(), "SYS_C") || !index.GetIsConstraint() || index.GetPrimary() || len(index.GetConstraints()) != 1 || index.GetConstraints()[0].GetName() != index.GetName() || index.GetConstraints()[0].GetPrimary() || !cmp.Equal(index.GetConstraints()[0].GetColumns(), []string{`"A"`}) {
		t.Errorf("UK_UNNAMED index = %v, want the unnamed constraint's index", index)
	}
	// The phantom the sync used to report: no index anywhere is named after
	// either constraint, which is what the metadata used to call its index.
	for _, name := range []string{"PK_REUSED_PK", "UK_REUSED_UK"} {
		if got := selectOne(t, d.db, fmt.Sprintf(`SELECT COUNT(*) FROM ALL_OBJECTS WHERE OWNER = 'OMNI_APP' AND OBJECT_TYPE = 'INDEX' AND OBJECT_NAME = '%s'`, name)); got != "0" {
			t.Errorf("%s: ALL_OBJECTS counts %s indexes, want none", name, got)
		}
	}
	// The index a constraint reuses is the object the DBA's own DDL names.
	if ddl := selectOne(t, d.db, `SELECT TO_CHAR(DBMS_METADATA.GET_DDL('INDEX', 'UK_REUSED_IDX', 'OMNI_APP')) FROM DUAL`); !strings.Contains(ddl, `"UK_REUSED_IDX"`) {
		t.Errorf("the index's DDL = %q, want it to name UK_REUSED_IDX", ddl)
	}
}

// TestSyncSequenceValues pins every fact ALL_SEQUENCES states about a sequence,
// which the sync carries. The declared START WITH is not one of them: the
// dictionary has no column for it, and DBMS_METADATA reports LAST_NUMBER in its
// place, so a sequence created START WITH 100 and advanced once by NEXTVAL is
// exported as START WITH 200. start therefore stays empty, and last_value
// carries LAST_NUMBER, the number the dictionary will hand out next.
func TestSyncSequenceValues(t *testing.T) {
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE SEQUENCE TO OMNI_APP;
CREATE SEQUENCE OMNI_APP.UP_SEQ START WITH 100 INCREMENT BY 5;
CREATE SEQUENCE OMNI_APP.DOWN_SEQ START WITH -3 INCREMENT BY -2 MINVALUE -100 MAXVALUE -1;
CREATE SEQUENCE OMNI_APP.CYCLE_SEQ MAXVALUE 5 CACHE 3 CYCLE;
CREATE SEQUENCE OMNI_APP.NOCACHE_SEQ NOCACHE;
CREATE SEQUENCE OMNI_APP.ORDER_SEQ ORDER;
CREATE SEQUENCE OMNI_APP.KEEP_SEQ KEEP;
CREATE SEQUENCE OMNI_APP.SCALE_SEQ SCALE EXTEND;
CREATE SEQUENCE OMNI_APP.SESSION_SEQ SESSION;
DECLARE
	V NUMBER;
BEGIN
	SELECT OMNI_APP.UP_SEQ.NEXTVAL INTO V FROM DUAL;
END;
`)

	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	schema := onlySchema(t, metadata)
	// Oracle's default MAXVALUE, which is 28 digits and so no Go integer.
	const maxValue = "9999999999999999999999999999"
	want := map[string]*metadatapb.SequenceMetadata{
		// One NEXTVAL of CACHE 20 INCREMENT BY 5 from 100 moves the number the
		// dictionary will hand out to 200.
		"UP_SEQ":    {Name: "UP_SEQ", MinValue: "1", MaxValue: maxValue, Increment: "5", CacheSize: "20", LastValue: "200"},
		"DOWN_SEQ":  {Name: "DOWN_SEQ", MinValue: "-100", MaxValue: "-1", Increment: "-2", CacheSize: "20", LastValue: "-3"},
		"CYCLE_SEQ": {Name: "CYCLE_SEQ", MinValue: "1", MaxValue: "5", Increment: "1", Cycle: true, CacheSize: "3", LastValue: "1"},
		// NOCACHE is a cache size of 0.
		"NOCACHE_SEQ": {Name: "NOCACHE_SEQ", MinValue: "1", MaxValue: maxValue, Increment: "1", CacheSize: "0", LastValue: "1"},
		"ORDER_SEQ":   {Name: "ORDER_SEQ", MinValue: "1", MaxValue: maxValue, Increment: "1", Order: true, CacheSize: "20", LastValue: "1"},
		"KEEP_SEQ":    {Name: "KEEP_SEQ", MinValue: "1", MaxValue: maxValue, Increment: "1", Keep: true, CacheSize: "20", LastValue: "1"},
		"SCALE_SEQ":   {Name: "SCALE_SEQ", MinValue: "1", MaxValue: maxValue, Increment: "1", Scale: true, Extend: true, CacheSize: "20", LastValue: "1"},
		"SESSION_SEQ": {Name: "SESSION_SEQ", MinValue: "1", MaxValue: maxValue, Increment: "1", Session: true, CacheSize: "20", LastValue: "1"},
	}
	got := make(map[string]*metadatapb.SequenceMetadata, len(schema.Sequences))
	for _, sequence := range schema.Sequences {
		got[sequence.GetName()] = sequence
	}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("sequences (-want +got):\n%s", diff)
	}
	// What the dictionary states, as the protos hold it: a release that gains a
	// column this test does not read fails here rather than in a golden.
	if got, want := selectOne(t, d.db, `SELECT COUNT(*) FROM ALL_TAB_COLS WHERE OWNER = 'SYS' AND TABLE_NAME = 'ALL_SEQUENCES' AND COLUMN_NAME NOT IN ('SEQUENCE_OWNER', 'SEQUENCE_NAME', 'MIN_VALUE', 'MAX_VALUE', 'INCREMENT_BY', 'CYCLE_FLAG', 'ORDER_FLAG', 'CACHE_SIZE', 'LAST_NUMBER', 'SCALE_FLAG', 'EXTEND_FLAG', 'SHARDED_FLAG', 'SESSION_FLAG', 'KEEP_VALUE')`), "0"; got != want {
		t.Errorf("ALL_SEQUENCES has %s columns the sync does not read, want %s", got, want)
	}
	// Oracle's own exporter has no START WITH either, and puts LAST_NUMBER in
	// its place: the declared 100 is gone, and 200 is what both the metadata's
	// last_value and the DDL say.
	if ddl := selectOne(t, d.db, `SELECT TO_CHAR(DBMS_METADATA.GET_DDL('SEQUENCE', 'UP_SEQ', 'OMNI_APP')) FROM DUAL`); !strings.Contains(ddl, "START WITH 200") {
		t.Errorf("the sequence's DDL = %q, want it to say START WITH 200", ddl)
	}
}

// TestSequenceProjection pins what the sync reads out of ALL_SEQUENCES: a fact
// the view has a column for, as text when it is a number, and a fact it has no
// column for as NULL, so a release older than the fact neither fails the sync
// nor loses the facts it does have. The pinned server has every column; the
// case for a release that does not is what keeps this from being a guess at
// which release added which flag.
func TestSequenceProjection(t *testing.T) {
	base := []string{"SEQUENCE_OWNER", "SEQUENCE_NAME", "MIN_VALUE", "MAX_VALUE", "INCREMENT_BY", "CYCLE_FLAG", "ORDER_FLAG", "CACHE_SIZE", "LAST_NUMBER"}
	tests := []struct {
		name    string
		columns []string
		want    string
	}{
		{
			name:    "the facts a release before sharded sequences had",
			columns: base,
			want:    "SEQUENCE_NAME, TO_CHAR(MIN_VALUE), TO_CHAR(MAX_VALUE), TO_CHAR(INCREMENT_BY), CYCLE_FLAG, ORDER_FLAG, TO_CHAR(CACHE_SIZE), TO_CHAR(LAST_NUMBER), NULL, NULL, NULL, NULL, NULL",
		},
		{
			name:    "the pinned server, with every column",
			columns: append(base, "SHARDED_FLAG", "SESSION_FLAG", "KEEP_VALUE", "SCALE_FLAG", "EXTEND_FLAG"),
			want:    "SEQUENCE_NAME, TO_CHAR(MIN_VALUE), TO_CHAR(MAX_VALUE), TO_CHAR(INCREMENT_BY), CYCLE_FLAG, ORDER_FLAG, TO_CHAR(CACHE_SIZE), TO_CHAR(LAST_NUMBER), KEEP_VALUE, SCALE_FLAG, EXTEND_FLAG, SESSION_FLAG, SHARDED_FLAG",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columns := make(map[string]bool, len(tt.columns))
			for _, column := range tt.columns {
				columns[column] = true
			}
			if got := sequenceProjection(columns); got != tt.want {
				t.Errorf("projection = %q, want %q", got, tt.want)
			}
		})
	}
}

// tableByName returns the synced table of the name.
func tableByName(t *testing.T, schema *metadatapb.SchemaMetadata, name string) *metadatapb.TableMetadata {
	t.Helper()
	for _, table := range schema.Tables {
		if table.GetName() == name {
			return table
		}
	}
	t.Fatalf("table %q is missing from %v", name, namesOfTables(schema.Tables))
	return nil
}

// selectOne runs a query of one row and one column, and returns its value as
// text.
func selectOne(t *testing.T, db *sql.DB, query string) string {
	t.Helper()
	var value sql.NullString
	if err := db.QueryRow(query).Scan(&value); err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	return value.String
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

func execStatements(t *testing.T, db *sql.DB, statements ...string) {
	t.Helper()
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	for _, stmt := range statements {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s:\n%v", stmt, err)
		}
	}
}

// TestSyncMaterializedViewIndexesAndTriggers pins that a materialized view is
// reported under MaterializedViews (not Tables) with both its system-generated
// ROWID index and user indexes/triggers attached.
func TestSyncMaterializedViewIndexesAndTriggers(t *testing.T) {
	d := openTestDriver(t)
	execStatements(t, d.db,
		`CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS`,
		`GRANT CREATE SESSION, CREATE TABLE, CREATE MATERIALIZED VIEW, CREATE TRIGGER TO OMNI_APP`,
		`CREATE TABLE OMNI_APP.BASE (ID NUMBER CONSTRAINT BASE_PK PRIMARY KEY, V NUMBER)`,
		`CREATE MATERIALIZED VIEW LOG ON OMNI_APP.BASE WITH PRIMARY KEY, ROWID`,
		`CREATE MATERIALIZED VIEW OMNI_APP.MV REFRESH FAST WITH ROWID AS SELECT ID, V FROM OMNI_APP.BASE`,
		`CREATE INDEX OMNI_APP.MV_V ON OMNI_APP.MV (V)`,
		`CREATE TRIGGER OMNI_APP.MV_TRG BEFORE INSERT ON OMNI_APP.MV FOR EACH ROW BEGIN NULL; END;`,
	)
	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	schema := onlySchema(t, metadata)
	for _, table := range schema.Tables {
		if table.Name == "MV" {
			t.Fatalf("materialized view MV also appeared in Tables: %v", namesOfTables(schema.Tables))
		}
	}
	if len(schema.MaterializedViews) != 1 {
		t.Fatalf("materialized views = %d, want 1", len(schema.MaterializedViews))
	}
	mv := schema.MaterializedViews[0]
	var indexNames []string
	for _, idx := range mv.Indexes {
		indexNames = append(indexNames, idx.Name)
	}
	if diff := cmp.Diff([]string{"I_SNAP$_MV", "MV_V"}, indexNames); diff != "" {
		t.Errorf("MV indexes (-want +got):\n%s", diff)
	}
	var triggerNames []string
	for _, trg := range mv.Triggers {
		triggerNames = append(triggerNames, trg.Name)
	}
	if diff := cmp.Diff([]string{"MV_TRG"}, triggerNames); diff != "" {
		t.Errorf("MV triggers (-want +got):\n%s", diff)
	}
}

// TestSyncLeavesOutTriggersOnOtherSchemasTables pins that a trigger owned by
// the synced schema targeting a same-named table in another schema is not
// attached to the local table.
func TestSyncLeavesOutTriggersOnOtherSchemasTables(t *testing.T) {
	d := openTestDriver(t)
	execStatements(t, d.db,
		`CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS`,
		`CREATE USER OMNI_OTHER NO AUTHENTICATION QUOTA UNLIMITED ON USERS`,
		`GRANT CREATE SESSION, CREATE TABLE, CREATE TRIGGER TO OMNI_APP`,
		`GRANT CREATE SESSION, CREATE TABLE TO OMNI_OTHER`,
		`CREATE TABLE OMNI_APP.T (A NUMBER)`,
		`CREATE TABLE OMNI_OTHER.T (A NUMBER)`,
		`CREATE TRIGGER OMNI_APP.ON_OWN BEFORE INSERT ON OMNI_APP.T FOR EACH ROW BEGIN NULL; END;`,
		`CREATE TRIGGER OMNI_APP.ON_OTHER BEFORE INSERT ON OMNI_OTHER.T FOR EACH ROW BEGIN NULL; END;`,
	)
	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	schema := onlySchema(t, metadata)
	table := tableByName(t, schema, "T")
	var triggerNames []string
	for _, trg := range table.Triggers {
		triggerNames = append(triggerNames, trg.Name)
	}
	if diff := cmp.Diff([]string{"ON_OWN"}, triggerNames); diff != "" {
		t.Errorf("T triggers (-want +got):\n%s", diff)
	}
}

// TestSyncSynonyms pins private synonyms (same-schema, cross-schema, and
// excluding dangling synonyms whose target does not exist) and public synonyms.
func TestSyncSynonyms(t *testing.T) {
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
CREATE USER OMNI_OTHER NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE, CREATE SYNONYM TO OMNI_APP;
GRANT CREATE SESSION, CREATE TABLE TO OMNI_OTHER;
CREATE TABLE OMNI_APP.T (A NUMBER);
CREATE TABLE OMNI_OTHER.OT (A NUMBER);
CREATE SYNONYM OMNI_APP.S_OWN FOR OMNI_APP.T;
CREATE SYNONYM OMNI_APP.S_OTHER FOR OMNI_OTHER.OT;
CREATE SYNONYM OMNI_APP.S_GONE FOR OMNI_APP.NOSUCH;
CREATE PUBLIC SYNONYM OMNI_APP_PUB FOR OMNI_APP.T;
`)
	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	schema := onlySchema(t, metadata)
	wantPrivate := []*metadatapb.SynonymMetadata{
		{Name: "S_OTHER", TargetSchema: "OMNI_OTHER", TargetName: "OT"},
		{Name: "S_OWN", TargetSchema: "OMNI_APP", TargetName: "T"},
	}
	if diff := cmp.Diff(wantPrivate, schema.Synonyms, protocmp.Transform()); diff != "" {
		t.Errorf("private synonyms (-want +got):\n%s", diff)
	}
	publicByName := make(map[string]*metadatapb.SynonymMetadata, len(metadata.PublicSynonyms))
	for _, syn := range metadata.PublicSynonyms {
		publicByName[syn.Name] = syn
	}
	if diff := cmp.Diff(&metadatapb.SynonymMetadata{Name: "OMNI_APP_PUB", TargetSchema: "OMNI_APP", TargetName: "T"}, publicByName["OMNI_APP_PUB"], protocmp.Transform()); diff != "" {
		t.Errorf("public synonym OMNI_APP_PUB (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(&metadatapb.SynonymMetadata{Name: "DUAL", TargetSchema: "SYS", TargetName: "DUAL"}, publicByName["DUAL"], protocmp.Transform()); diff != "" {
		t.Errorf("public synonym DUAL (-want +got):\n%s", diff)
	}
}
