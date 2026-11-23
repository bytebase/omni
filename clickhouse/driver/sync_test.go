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
	"strings"
	"testing"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	chanalysis "github.com/bytebase/omni/clickhouse/analysis"
	"github.com/bytebase/omni/clickhouse/internal/testcontainer"
	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata from this run")

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
	statements, err := chanalysis.Split(script)
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
	access, err := os.ReadFile(filepath.Join("testdata", "sync", "access.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openTestDriver(t)
	execCorpus(t, d.db, string(corpus))
	execCorpus(t, d.db, strings.ReplaceAll(string(access), "$database", testcontainer.Database(t)))

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
	// A role's attribute is one access statement per line, which the golden
	// file holds as a list.
	type role struct {
		Name       string   `json:"name"`
		ValidUntil *string  `json:"validUntil,omitempty"`
		Attribute  []string `json:"attribute"`
	}
	var roles []role
	for i, r := range instance.Roles {
		if i > 0 && instance.Roles[i-1].Name > r.Name {
			t.Errorf("roles are not ordered: %q before %q", instance.Roles[i-1].Name, r.Name)
		}
		if r.Attribute == nil {
			t.Errorf("role %q carries no attribute", r.Name)
			continue
		}
		roles = append(roles, role{Name: r.Name, ValidUntil: r.ValidUntil, Attribute: strings.Split(*r.Attribute, "\n")})
	}
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

// TestWrapDatabaseNotFound holds that UNKNOWN_DATABASE is driver.ErrNotFound and
// that a failure of another kind keeps its own error.
func TestWrapDatabaseNotFound(t *testing.T) {
	err := wrapDatabaseNotFound(fmt.Errorf("failed to execute query: %w", &clickhouse.Exception{Code: unknownDatabase, Message: "Database omni_no_such_database does not exist"}))
	if !errors.Is(err, driver.ErrNotFound) {
		t.Errorf("wrapDatabaseNotFound(81) = %v, want driver.ErrNotFound", err)
	}
	other := &clickhouse.Exception{Code: 60, Message: "Unknown table"}
	if err := wrapDatabaseNotFound(other); errors.Is(err, driver.ErrNotFound) {
		t.Errorf("wrapDatabaseNotFound(60) = %v, want the engine's own error", err)
	}
}

// TestSyncInstanceAccessDenied holds the roles a login reads when it may list
// the users and nothing else of the server's access entities: the users alone,
// with no attribute, which is what the sync read before it read more.
func TestSyncInstanceAccessDenied(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, `CREATE ROLE omni_hidden;
CREATE USER omni_users_only IDENTIFIED WITH no_password;
GRANT SELECT ON system.users TO omni_users_only;`)
	impersonate(t, d, "omni_users_only")

	roles, err := d.getInstanceRoles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, role := range roles {
		names = append(names, role.Name)
		if role.Attribute != nil || role.ValidUntil != nil {
			t.Errorf("role %q carries an attribute or an expiry the login may not read", role.Name)
		}
	}
	if diff := cmp.Diff([]string{"default", "omni_users_only"}, names); diff != "" {
		t.Errorf("roles differ (-want +got):\n%s", diff)
	}
}

// impersonate makes every later statement of d run as user: the pool keeps one
// session, and EXECUTE AS switches it.
func impersonate(t *testing.T, d *Driver, user string) {
	t.Helper()
	d.db.SetMaxOpenConns(1)
	d.db.SetMaxIdleConns(1)
	if _, err := d.db.ExecContext(context.Background(), "EXECUTE AS "+user); err != nil {
		t.Fatal(err)
	}
	var current string
	if err := d.db.QueryRowContext(context.Background(), "SELECT currentUser()").Scan(&current); err != nil {
		t.Fatal(err)
	}
	if current != user {
		t.Fatalf("the session runs as %q, want %q", current, user)
	}
}
