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
	"sort"
	"testing"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bytebase/omni/clickhouse/internal/testcontainer"
	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync from this run")

func openTestDriver(t *testing.T) *Driver {
	t.Helper()
	db := testcontainer.Open(t)
	t.Cleanup(func() { _ = db.Close() })
	return &Driver{db: db, databaseName: testcontainer.Database(t)}
}

// execCorpus applies a script one statement at a time on one connection, which
// is the connection the tape holds.
func execCorpus(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	statements, err := splitSQL(script)
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range filterEmptyStatements(statements) {
		if _, err := conn.ExecContext(ctx, statement.Text); err != nil {
			t.Fatalf("%s:\n%v", statement.Text, err)
		}
	}
}

func readGolden(t *testing.T, path string) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	metadata := &metadatapb.DatabaseSchemaMetadata{}
	if err := protojson.Unmarshal(b, metadata); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return metadata
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

func TestSyncDatabaseDefaultDatabase(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	db := testcontainer.Open(t)
	t.Cleanup(func() { _ = db.Close() })
	execCorpus(t, db, string(corpus))

	// The driver was opened without a database, so the connection's own
	// database is the one the sync reads.
	d := &Driver{db: db}
	got, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want := testcontainer.Database(t); got.Name != want {
		t.Errorf("Name = %q, want the connection's %q", got.Name, want)
	}
	if len(got.Schemas) != 1 || len(got.Schemas[0].Tables) != 2 || len(got.Schemas[0].Views) != 1 {
		t.Errorf("schema = %+v", got.Schemas)
	}
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
	databases := make([]string, 0, len(instance.Databases))
	for _, database := range instance.Databases {
		databases = append(databases, database.Name)
		if len(database.Schemas) != 0 {
			t.Errorf("database %q carries %d schemas; an instance sync reports none", database.Name, len(database.Schemas))
		}
	}
	for i := 1; i < len(databases); i++ {
		if databases[i-1] > databases[i] {
			t.Errorf("databases are not ordered: %q before %q", databases[i-1], databases[i])
		}
	}
	sort.Strings(databases)
	roles := make([]string, 0, len(instance.Roles))
	for _, role := range instance.Roles {
		roles = append(roles, role.Name)
	}
	for i := 1; i < len(roles); i++ {
		if roles[i-1] > roles[i] {
			t.Errorf("roles are not ordered: %q before %q", roles[i-1], roles[i])
		}
	}
	sort.Strings(roles)
	out, err := json.MarshalIndent(map[string]any{
		"version":   instance.Version,
		"databases": databases,
		"roles":     roles,
	}, "", "  ")
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

// TestSyncDatabaseMissing holds that a database ClickHouse says it does not have
// is driver.ErrNotFound, on the connection the client cannot make: the handshake
// is refused, before any statement reaches the server.
func TestSyncDatabaseMissing(t *testing.T) {
	db := testcontainer.OpenMissingDatabase(t, missingDatabase)
	t.Cleanup(func() { _ = db.Close() })

	// What the server answered, with the client's own error type, as a replay
	// holds it.
	err := db.PingContext(context.Background())
	var engineError *clickhouse.Exception
	if !errors.As(err, &engineError) || engineError.Code != unknownDatabase {
		t.Fatalf("Ping = %v, want a ClickHouse exception 81", err)
	}

	d := &Driver{db: db, databaseName: missingDatabase}
	if _, err := d.SyncDatabase(context.Background()); !errors.Is(err, driver.ErrNotFound) {
		t.Fatalf("SyncDatabase = %v, want driver.ErrNotFound", err)
	}
}

// missingDatabase is a database no test creates, so a recording's server does
// not have it.
const missingDatabase = "omni_no_such_database"

// TestWrapSyncError holds that UNKNOWN_DATABASE is driver.ErrNotFound and that a
// failure of another kind keeps its own error.
func TestWrapSyncError(t *testing.T) {
	err := wrapSyncError(fmt.Errorf("failed to execute query: %w", &clickhouse.Exception{Code: unknownDatabase, Message: "Database omni_no_such_database does not exist"}))
	if !errors.Is(err, driver.ErrNotFound) {
		t.Errorf("wrapSyncError(81) = %v, want driver.ErrNotFound", err)
	}
	other := &clickhouse.Exception{Code: 60, Message: "Unknown table"}
	if err := wrapSyncError(other); errors.Is(err, driver.ErrNotFound) {
		t.Errorf("wrapSyncError(60) = %v, want the engine's own error", err)
	}
}
