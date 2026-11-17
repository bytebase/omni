package driver

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/internal/testcontainer"
	oracleparser "github.com/bytebase/omni/oracle/parser"
)

// TestDump holds the dump of the sync fixture to testdata/dump/database.sql and
// checks that it does not mutate the snapshot it renders. Rewrite the golden
// with -update after changing the catalog.
func TestDump(t *testing.T) {
	m := readGolden(t, filepath.Join("testdata", "sync", "database.json"))
	before := proto.Clone(m)
	var b strings.Builder
	if err := (&Driver{}).Dump(context.Background(), &b, m); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(before, m) {
		t.Error("Dump mutated its input metadata")
	}
	golden := filepath.Join("testdata", "dump", "database.sql")
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
		t.Errorf("Dump differs from %s (-want +got):\n%s", golden, diff)
	}
}

// dumpSchema is the schema TestDumpRestores and TestDumpGaps fill.
const dumpSchema = "OMNI_DUMP"

// createSchema drops and recreates a user, the schema a restore lands in, and
// makes it the session's current schema.
func createSchema(t *testing.T, db *sql.DB, schema string) {
	t.Helper()
	ctx := context.Background()
	_, _ = db.ExecContext(ctx, "DROP USER "+schema+" CASCADE")
	execCorpus(t, db, fmt.Sprintf(`CREATE USER %s IDENTIFIED BY "%s" QUOTA UNLIMITED ON USERS`, schema, testcontainer.Password))
	execCorpus(t, db, fmt.Sprintf(`GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE SEQUENCE, CREATE PROCEDURE, CREATE TRIGGER, CREATE MATERIALIZED VIEW, CREATE TYPE, CREATE SYNONYM, CREATE PUBLIC SYNONYM, CREATE DATABASE LINK, CREATE PUBLIC DATABASE LINK, DROP PUBLIC DATABASE LINK TO %s`, schema))
	execCorpus(t, db, "ALTER SESSION SET CURRENT_SCHEMA = "+schema)
}

// syncSchema reads one schema's snapshot through the driver.
func syncSchema(t *testing.T, db *sql.DB, schema string) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	d := &Driver{db: db, databaseName: schema, engineVersion: testEngineVersion}
	m, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatalf("syncing %s: %v", schema, err)
	}
	return m
}

// TestDumpRestores fills a schema with the corpus, syncs it, dumps it, rebuilds
// the schema and restores the dump, then requires the resynced snapshot to
// match the source's.
func TestDumpRestores(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "dump", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	db := testcontainer.Open(t, "system")
	defer func() { _ = db.Close() }()

	createSchema(t, db, dumpSchema)
	execCorpus(t, db, string(corpus))
	source := syncSchema(t, db, dumpSchema)

	var dump strings.Builder
	if err := (&Driver{db: db, databaseName: dumpSchema, engineVersion: testEngineVersion}).Dump(context.Background(), &dump, source); err != nil {
		t.Fatalf("dumping: %v", err)
	}

	createSchema(t, db, dumpSchema)
	execCorpus(t, db, dump.String())
	restored := syncSchema(t, db, dumpSchema)

	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("the restored schema differs from the source (-source +restored):\n%s", diff)
	}
}

// restoreDiff compares two snapshots of one schema after a dump and restore. It
// holds out what Oracle generates per object: a table's row estimate, the name
// of an identity's backing sequence in the column default, and a materialized
// view's container index. docs/migration/oracle/catalog.md names each.
func restoreDiff(source, restored *metadatapb.DatabaseSchemaMetadata) string {
	a := proto.Clone(source).(*metadatapb.DatabaseSchemaMetadata)
	b := proto.Clone(restored).(*metadatapb.DatabaseSchemaMetadata)
	normalizeSnapshot(a)
	normalizeSnapshot(b)
	return cmp.Diff(a, b, protocmp.Transform(),
		protocmp.IgnoreFields(&metadatapb.TableMetadata{}, "row_count"),
	)
}

func normalizeSnapshot(m *metadatapb.DatabaseSchemaMetadata) {
	for _, schema := range m.GetSchemas() {
		for _, table := range schema.GetTables() {
			for _, column := range table.GetColumns() {
				if strings.Contains(column.GetDefault(), "ISEQ$$_") {
					column.Default = ""
				}
			}
		}
		for _, view := range schema.GetMaterializedViews() {
			var indexes []*metadatapb.IndexMetadata
			for _, index := range view.GetIndexes() {
				if !strings.HasPrefix(index.GetName(), "SYS_C_SNAP$") {
					indexes = append(indexes, index)
				}
			}
			view.Indexes = indexes
		}
	}
}

// TestDumpGaps holds what a dump cannot faithfully restore yet, one schema per
// gap so no gap's objects reach another's dump. A gap with a query restores but
// changes the property the query reads; a gap without one leaves the restored
// snapshot different from the source's. A gap that closes fails here: move it
// into TestDumpRestores. docs/migration/oracle/catalog.md names each gap the
// catalog documents.
func TestDumpGaps(t *testing.T) {
	db := testcontainer.Open(t, "system")
	defer func() { _ = db.Close() }()
	ctx := context.Background()

	gaps := []struct {
		name string
		// setup fills the gap's schema.
		setup string
		// beforeRestore drops an object the dump cannot carry and the restore
		// therefore lands without, the way another instance would.
		beforeRestore string
		// prepare runs before setup on the same session, and a failure is
		// ignored: it clears what an earlier recording left behind.
		prepare string
		// cleanup runs after the subtest, even when its setup fails, with the
		// session back in SYSTEM's schema.
		cleanup string
		// ownerSession runs the setup on a session of the schema's own user,
		// which an object owned by the session user needs (a database link).
		ownerSession bool
		// query reads one property of the schema; %s is its name. Empty means
		// the restored snapshot itself differs.
		query string
		// check replaces the snapshot comparison when the gap lives outside the
		// schema, which a dump into the same instance cannot show.
		check func(t *testing.T, source, restored *metadatapb.DatabaseSchemaMetadata, dump string)
	}{
		{
			name:  "synonym",
			setup: `CREATE TABLE T (C NUMBER); CREATE SYNONYM S FOR T;`,
		},
		{
			name:          "public synonym",
			setup:         `CREATE TABLE T (C NUMBER); CREATE PUBLIC SYNONYM OMNIGAP_PS FOR T;`,
			beforeRestore: `DROP PUBLIC SYNONYM OMNIGAP_PS`,
		},
		{
			name:  "package",
			setup: `CREATE TABLE T (C NUMBER); CREATE OR REPLACE PACKAGE PKG AS PROCEDURE P; END PKG;`,
		},
		{
			// A database link belongs to the session user, not the current
			// schema, so the schema's own session creates it, and the restore
			// lands in an instance that never had it.
			name:         "database link",
			prepare:      `DROP PUBLIC DATABASE LINK OMNIGAP_LINK`,
			setup:        `CREATE TABLE T (C NUMBER); CREATE PUBLIC DATABASE LINK OMNIGAP_LINK CONNECT TO U IDENTIFIED BY P USING 'example.com:1521/x';`,
			cleanup:      `DROP PUBLIC DATABASE LINK OMNIGAP_LINK`,
			ownerSession: true,
			// A link belongs to the instance, not the schema: a restore into the
			// same instance still sees it, so what the dump loses is the
			// statement itself.
			check: func(t *testing.T, source, _ *metadatapb.DatabaseSchemaMetadata, dump string) {
				t.Helper()
				if len(source.GetLinkedDatabases()) == 0 {
					t.Fatal("the snapshot does not hold the database link")
				}
				if strings.Contains(dump, "DATABASE LINK") {
					t.Fatal("gap closed: the dump carries the database link")
				}
			},
		},
		{
			name:  "object table",
			setup: `CREATE TYPE TY AS OBJECT (A NUMBER); CREATE TABLE OT OF TY; CREATE TABLE REL (A NUMBER);`,
		},
		{
			name:  "XMLType table",
			setup: `CREATE TABLE XT OF XMLTYPE; CREATE TABLE REL (A NUMBER);`,
		},
		{
			name:  "view column aliases",
			setup: `CREATE TABLE T (A NUMBER, B NUMBER); CREATE VIEW V (RENAMED_A, RENAMED_B) AS SELECT A, B FROM T;`,
		},
		{
			name:  "partitioning",
			setup: `CREATE TABLE P (A NUMBER) PARTITION BY RANGE (A) (PARTITION P1 VALUES LESS THAN (10), PARTITION P2 VALUES LESS THAN (MAXVALUE));`,
			query: `SELECT COUNT(*) FROM all_part_tables WHERE owner = '%s' AND table_name = 'P'`,
		},
		{
			name:  "tablespace and storage",
			setup: `CREATE TABLESPACE OMNIGAP_TS DATAFILE 'omnigap_ts.dbf' SIZE 10M; CREATE TABLE TS_T (A NUMBER) TABLESPACE OMNIGAP_TS STORAGE (INITIAL 64K);`,
			query: `SELECT tablespace_name FROM all_tables WHERE owner = '%s' AND table_name = 'TS_T'`,
		},
		{
			name:  "external table",
			setup: `CREATE OR REPLACE DIRECTORY OMNIGAP_DIR AS '/tmp'; CREATE TABLE EXT (A NUMBER) ORGANIZATION EXTERNAL (TYPE ORACLE_LOADER DEFAULT DIRECTORY OMNIGAP_DIR ACCESS PARAMETERS (FIELDS TERMINATED BY ',') LOCATION ('omnigap.csv'));`,
			query: `SELECT COUNT(*) FROM all_external_tables WHERE owner = '%s' AND table_name = 'EXT'`,
		},
		{
			// The dictionary does not carry a sequence's declared START WITH,
			// so the dump has none to write and the restore starts at MINVALUE.
			name:  "sequence start",
			setup: `CREATE SEQUENCE S START WITH 10 INCREMENT BY 2 CACHE 20 MAXVALUE 1000;`,
			query: `SELECT last_number FROM all_sequences WHERE sequence_owner = '%s' AND sequence_name = 'S'`,
		},
		{
			name:  "identity sequence options",
			setup: `CREATE TABLE I (ID NUMBER GENERATED ALWAYS AS IDENTITY (START WITH 5 INCREMENT BY 2 MAXVALUE 100000 CACHE 100 CYCLE));`,
			query: `SELECT LISTAGG(cache_size || ':' || cycle_flag, ',') WITHIN GROUP (ORDER BY sequence_name) FROM all_sequences WHERE sequence_owner = '%s' AND sequence_name LIKE 'ISEQ$$_%%'`,
		},
		{
			name:  "materialized view refresh policy",
			setup: `CREATE TABLE T (C NUMBER); CREATE MATERIALIZED VIEW MV REFRESH COMPLETE ON DEMAND AS SELECT C FROM T;`,
			query: `SELECT refresh_method || '/' || refresh_mode FROM all_mviews WHERE owner = '%s' AND mview_name = 'MV'`,
		},
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP TABLESPACE OMNIGAP_TS INCLUDING CONTENTS AND DATAFILES")
		_, _ = db.ExecContext(ctx, "DROP DIRECTORY OMNIGAP_DIR")
	})

	for i, gap := range gaps {
		t.Run(gap.name, func(t *testing.T) {
			schema := fmt.Sprintf("OMNI_GAP_%d", i)
			createSchema(t, db, schema)
			if gap.cleanup != "" {
				t.Cleanup(func() {
					// A session whose current schema is another user's cannot
					// drop a public link, even as SYSTEM.
					execIgnoringErrors(t, db, "ALTER SESSION SET CURRENT_SCHEMA = SYSTEM")
					execIgnoringErrors(t, db, gap.cleanup)
				})
			}
			if gap.ownerSession {
				userDB := testcontainer.Open(t, schema)
				if gap.prepare != "" {
					execIgnoringErrors(t, userDB, gap.prepare)
				}
				execCorpus(t, userDB, gap.setup)
				_ = userDB.Close()
			} else {
				execCorpus(t, db, gap.setup)
			}
			source := syncSchema(t, db, schema)
			var sourceValue string
			if gap.query != "" {
				sourceValue = queryString(t, db, fmt.Sprintf(gap.query, schema))
			}

			var dump strings.Builder
			if err := (&Driver{db: db, databaseName: schema, engineVersion: testEngineVersion}).Dump(ctx, &dump, source); err != nil {
				t.Fatalf("dumping: %v", err)
			}

			if gap.beforeRestore != "" {
				execCorpus(t, db, gap.beforeRestore)
			}
			createSchema(t, db, schema)
			execCorpus(t, db, dump.String())
			restored := syncSchema(t, db, schema)

			if gap.check != nil {
				gap.check(t, source, restored, dump.String())
				return
			}
			if gap.query != "" {
				restoredValue := queryString(t, db, fmt.Sprintf(gap.query, schema))
				if sourceValue == restoredValue {
					t.Fatalf("gap closed: %q restores as %q; move the case into TestDumpRestores", sourceValue, restoredValue)
				}
				return
			}
			if restoreDiff(source, restored) == "" {
				t.Fatal("gap closed: the restore reproduced the source snapshot; move the case into TestDumpRestores")
			}
		})
	}
}

// execIgnoringErrors runs a statement whose failure is expected, such as one
// that clears state an earlier run left.
func execIgnoringErrors(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	ctx := context.Background()
	for _, statement := range oracleparser.Split(script) {
		if statement.Empty() {
			continue
		}
		_, _ = db.ExecContext(ctx, statement.Text)
	}
}

func queryString(t *testing.T, db *sql.DB, query string) string {
	t.Helper()
	var value sql.NullString
	if err := db.QueryRowContext(context.Background(), query).Scan(&value); err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	return value.String
}
