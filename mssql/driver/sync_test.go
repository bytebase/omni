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

	gomssqldb "github.com/microsoft/go-mssqldb"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/internal/testcontainer"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync from this run")

// openTestDriver returns a driver on the database of the top-level test, with
// the corpus applied.
func openTestDriver(t *testing.T) *Driver {
	t.Helper()
	db := testcontainer.Open(t)
	t.Cleanup(func() { _ = db.Close() })
	return &Driver{db: db, databaseName: testcontainer.Database(t)}
}

// execCorpus applies a script of GO-separated batches.
func execCorpus(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	var batch strings.Builder
	flush := func() {
		text := strings.TrimSpace(batch.String())
		batch.Reset()
		if text == "" {
			return
		}
		if _, err := db.ExecContext(context.Background(), text); err != nil {
			t.Fatalf("%s:\n%v", text, err)
		}
	}
	for _, line := range strings.Split(script, "\n") {
		if strings.EqualFold(strings.TrimSpace(line), "GO") {
			flush()
			continue
		}
		batch.WriteString(line)
		batch.WriteString("\n")
	}
	flush()
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

	var sequences, indexes int
	for _, schema := range got.Schemas {
		if schema.Name != "app" {
			continue
		}
		for _, sequence := range schema.Sequences {
			switch sequence.Name {
			case "seq_orders":
				sequences++
				if sequence.Start != "10" || sequence.Increment != "2" || sequence.MinValue != "1" ||
					sequence.MaxValue != "10000" || !sequence.Cycle || sequence.CacheSize != "50" {
					t.Errorf("seq_orders synced without its options: %+v", sequence)
				}
			case "seq_decimal":
				sequences++
				if sequence.DataType != "decimal(38,0)" {
					t.Errorf("seq_decimal synced without its precision: %q", sequence.DataType)
				}
			}
		}
		for _, table := range schema.Tables {
			if table.Name != "accounts" {
				continue
			}
			for _, index := range table.Indexes {
				if index.Name != "idx_accounts_status" {
					continue
				}
				indexes++
				if len(index.Expressions) != 1 || index.Expressions[0] != "status" {
					t.Errorf("idx_accounts_status synced its INCLUDE column as a key: %v", index.Expressions)
				}
				if !strings.Contains(index.Definition, "INCLUDE:email") || !strings.Contains(index.Definition, "WHERE (") {
					t.Errorf("idx_accounts_status synced without its INCLUDE or predicate: %q", index.Definition)
				}
			}
		}
	}
	if sequences != 2 || indexes != 1 {
		t.Fatalf("the corpus objects are missing from the sync: %d sequences, %d indexes", sequences, indexes)
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
	d := openTestDriver(t)
	if _, err := d.db.ExecContext(context.Background(), "CREATE DATABASE "+quoteIdent("omni_instance_extra")); err != nil {
		t.Fatal(err)
	}
	instance, err := d.SyncInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// The server's own databases and whatever an earlier test left differ
	// between runs, so the golden holds the version and this test's database.
	var names []string
	for _, database := range instance.Databases {
		if database.Name == "omni_instance_extra" || database.Name == testcontainer.Database(t) {
			names = append(names, fmt.Sprintf("%s %s", database.Name, database.Collation))
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

// TestSyncSchemaWithApostrophe pins that a schema name carrying an apostrophe
// reaches the catalog queries as one literal.
func TestSyncSchemaWithApostrophe(t *testing.T) {
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE SCHEMA [owner's];
GO
CREATE TABLE [owner's].items (id INT NOT NULL PRIMARY KEY, note NVARCHAR(20));
`)
	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatalf("a schema named with an apostrophe failed the sync: %v", err)
	}
	var found bool
	for _, schema := range metadata.Schemas {
		if schema.Name != "owner's" {
			continue
		}
		found = true
		if len(schema.Tables) != 1 || schema.Tables[0].Name != "items" || len(schema.Tables[0].Columns) != 2 {
			t.Errorf("schema %q synced %+v", schema.Name, schema.Tables)
		}
	}
	if !found {
		t.Fatalf("the schema is missing from %+v", metadata.Schemas)
	}
}

// TestSyncDatabaseMissing holds that a database SQL Server says the login may
// not open is driver.ErrNotFound, on the connection the pool cannot make: a
// database that is not there is refused when the connection opens, before any
// statement reaches the server.
func TestSyncDatabaseMissing(t *testing.T) {
	db := testcontainer.OpenMissingDatabase(t, missingDatabase)
	t.Cleanup(func() { _ = db.Close() })

	// What the server answered, with the client's own error type, as a replay
	// holds it. 4063 is the warning a login that may fall back to its default
	// database gets, and go-mssqldb reports it as the connect error.
	err := db.PingContext(context.Background())
	var engineError gomssqldb.Error
	if !errors.As(err, &engineError) || engineError.Number != 4063 {
		t.Fatalf("Ping = %v, want a SQL Server error 4063", err)
	}

	d := &Driver{db: db, databaseName: missingDatabase}
	if _, err := d.SyncDatabase(context.Background()); !errors.Is(err, driver.ErrNotFound) {
		t.Fatalf("SyncDatabase = %v, want driver.ErrNotFound", err)
	}
}

// missingDatabase is a database no test creates, so a recording's server does
// not have it.
const missingDatabase = "omni_no_such_database"

// TestWrapSyncError holds the numbers SQL Server answers for a database the
// login may not open, and that a failure of another kind keeps its own error.
func TestWrapSyncError(t *testing.T) {
	for _, number := range []int32{911, 4060, 4063} {
		err := wrapSyncError(missingDatabase, gomssqldb.Error{Number: number, Message: "Cannot open database"})
		if !errors.Is(err, driver.ErrNotFound) {
			t.Errorf("wrapSyncError(%d) = %v, want driver.ErrNotFound", number, err)
		}
	}
	other := gomssqldb.Error{Number: 229, Message: "The SELECT permission was denied"}
	if err := wrapSyncError(missingDatabase, other); errors.Is(err, driver.ErrNotFound) {
		t.Errorf("wrapSyncError(229) = %v, want the engine's own error", err)
	}
}
