package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bytebase/omni/cockroachdb/internal/testcontainer"
	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync from this run")

// serverVersion is what SHOW CLUSTER SETTING version answers for the pin.
const serverVersion = "25.2"

func openTestDriver(t *testing.T) *Driver {
	t.Helper()
	db := testcontainer.Open(t)
	t.Cleanup(func() { _ = db.Close() })
	return &Driver{db: db, databaseName: testcontainer.Database(t), engineVersion: serverVersion}
}

func mustExec(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	if _, err := db.ExecContext(context.Background(), script); err != nil {
		t.Fatalf("%v", err)
	}
}

// TestSyncDatabase syncs testdata/corpus.sql and holds the result to
// testdata/sync/default.json.
func TestSyncDatabase(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openTestDriver(t)
	mustExec(t, d.db, string(corpus))
	got, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	golden := filepath.Join("testdata", "sync", "default.json")
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

// TestSyncInstance reads the server version, databases and roles and holds
// them to testdata/sync/instance.json.
func TestSyncInstance(t *testing.T) {
	d := openTestDriver(t)
	mustExec(t, d.db, `
CREATE DATABASE sync_extra;
CREATE ROLE sync_reader;
CREATE USER sync_login;
CREATE ROLE sync_admin WITH CREATEROLE CREATEDB LOGIN;`)

	instance, err := d.SyncInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	databases := make([]string, 0, len(instance.Databases))
	for _, db := range instance.Databases {
		databases = append(databases, fmt.Sprintf("%s %s %s", db.Name, db.CharacterSet, db.Collation))
	}
	roles := make([]string, 0, len(instance.Roles))
	for _, r := range instance.Roles {
		roles = append(roles, r.Name+": "+*r.Attribute)
	}
	b, err := json.MarshalIndent(map[string]any{"version": instance.Version, "databases": databases, "roles": roles}, "  ", "  ")
	if err != nil {
		t.Fatal(err)
	}
	b = append(b, '\n')

	golden := filepath.Join("testdata", "sync", "instance.json")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, b, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(want), string(b)); diff != "" {
		t.Errorf("SyncInstance differs from %s (-want +got):\n%s", golden, diff)
	}
}

// TestSyncDatabaseNotFound reports driver.ErrNotFound for a database the
// server does not have.
func TestSyncDatabaseNotFound(t *testing.T) {
	db := testcontainer.Open(t)
	t.Cleanup(func() { _ = db.Close() })
	d := &Driver{db: db, databaseName: "omni_no_such_database", engineVersion: serverVersion}
	if _, err := d.SyncDatabase(context.Background()); !errors.Is(err, driver.ErrNotFound) {
		t.Fatalf("SyncDatabase = %v, want driver.ErrNotFound", err)
	}
}

// TestSyncQueriesOrderTheirRows holds every row-returning query in sync.go to
// an ORDER BY. A sync query's row order decides the order of the metadata's
// slices, so an order the container picks records a different tape and golden
// on the next run.
func TestSyncQueriesOrderTheirRows(t *testing.T) {
	b, err := os.ReadFile("sync.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, query := range []string{
		listForeignKeyQuery,
		listSchemaQuery,
		getListForeignTableQuery(),
		getListTableQuery(true),
		listTablePartitionQuery,
		listColumnQuery,
		listMaterializedViewQuery,
		listViewQuery,
		listIndexQuery,
		listFunctionQuery,
	} {
		if !strings.Contains(strings.ToUpper(query), "ORDER BY") {
			t.Errorf("a sync query returns rows the container orders itself:\n%s", query)
		}
	}
	if !strings.Contains(strings.ToUpper(string(b)), "ORDER BY") {
		t.Error("sync.go holds no ORDER BY")
	}
}

// TestSyncViewDependencyQuoting syncs a schema and view whose names hold an
// apostrophe. The dependency query binds both names, so a quote cannot break
// it and one object cannot prevent the database from syncing.
func TestSyncViewDependencyQuoting(t *testing.T) {
	d := openTestDriver(t)
	mustExec(t, d.db, `
CREATE SCHEMA "app's";
CREATE TABLE "app's".t (id INT8 PRIMARY KEY);
CREATE VIEW "app's"."owner's_view" AS SELECT id FROM "app's".t;`)

	got, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, schema := range got.Schemas {
		if schema.Name != "app's" {
			continue
		}
		if len(schema.Views) != 1 || schema.Views[0].Name != "owner's_view" {
			t.Fatalf("views = %+v, want owner's_view", schema.Views)
		}
		if len(schema.Views[0].DependencyColumns) != 1 {
			t.Fatalf("dependency columns = %+v, want one", schema.Views[0].DependencyColumns)
		}
		return
	}
	t.Fatal("schema app's is missing")
}

// TestSyncForeignKeyQuotedSchema syncs a foreign key whose schema holds a
// period. The catalog query returns the relation name on its own, so the
// period is not mistaken for a qualifier.
func TestSyncForeignKeyQuotedSchema(t *testing.T) {
	d := openTestDriver(t)
	mustExec(t, d.db, `
CREATE SCHEMA "sales.eu";
CREATE TABLE "sales.eu".parent (id INT8 PRIMARY KEY);
CREATE TABLE "sales.eu".child (id INT8 PRIMARY KEY, parent_id INT8 REFERENCES "sales.eu".parent(id));
CREATE SCHEMA app;
CREATE TABLE app.parent_comma ("last,first" INT8 PRIMARY KEY);
CREATE TABLE app.child_comma ("last,first" INT8 REFERENCES app.parent_comma("last,first"));
CREATE TABLE app.parent_paren ("x(y" INT8 PRIMARY KEY);
CREATE TABLE app.child_paren ("x(y" INT8 REFERENCES app.parent_paren("x(y"));`)

	got, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, schema := range got.Schemas {
		switch schema.Name {
		case "sales.eu":
			for _, table := range schema.Tables {
				if table.Name != "child" {
					continue
				}
				if len(table.ForeignKeys) != 1 {
					t.Fatalf("foreign keys = %+v, want one", table.ForeignKeys)
				}
				if table.ForeignKeys[0].ReferencedTable != "parent" {
					t.Errorf("referenced table = %q, want parent", table.ForeignKeys[0].ReferencedTable)
				}
			}
		case "app":
			for _, table := range schema.Tables {
				switch table.Name {
				case "child_comma":
					if len(table.ForeignKeys) != 1 {
						t.Fatalf("foreign keys = %+v, want one", table.ForeignKeys)
					}
					fk := table.ForeignKeys[0]
					if len(fk.Columns) != 1 || fk.Columns[0] != "last,first" {
						t.Errorf("columns = %+v, want last,first", fk.Columns)
					}
					if len(fk.ReferencedColumns) != 1 || fk.ReferencedColumns[0] != "last,first" {
						t.Errorf("referenced columns = %+v, want last,first", fk.ReferencedColumns)
					}
				case "child_paren":
					if len(table.ForeignKeys) != 1 {
						t.Fatalf("foreign keys = %+v, want one", table.ForeignKeys)
					}
					fk := table.ForeignKeys[0]
					if len(fk.Columns) != 1 || fk.Columns[0] != "x(y" {
						t.Errorf("columns = %+v, want x(y", fk.Columns)
					}
					if len(fk.ReferencedColumns) != 1 || fk.ReferencedColumns[0] != "x(y" {
						t.Errorf("referenced columns = %+v, want x(y", fk.ReferencedColumns)
					}
				}
			}
		}
	}
}
