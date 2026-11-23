package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/starrocks/internal/testcontainer"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync from this run")

// openTestDriver returns a driver on dbName, dropped and created empty first so
// the tape holds one run's requests whatever the test binary ran before.
func openTestDriver(t *testing.T, dbName string) *Driver {
	t.Helper()
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{
			"DROP DATABASE IF EXISTS " + quoteBacktick(dbName),
			"CREATE DATABASE " + quoteBacktick(dbName),
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.OpenWith(t, testcontainer.Options{Database: dbName, MultiStatements: true})
	t.Cleanup(func() { _ = db.Close() })
	return &Driver{db: db, databaseName: dbName}
}

func execCorpus(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	for _, chunk := range strings.Split(script, "\n\n") {
		if stmt := strings.TrimSpace(chunk); stmt != "" {
			if _, err := db.ExecContext(context.Background(), stmt); err != nil {
				t.Fatalf("%s:\n%v", stmt, err)
			}
		}
	}
}

func TestSyncDatabase(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openTestDriver(t, "omni_sync_corpus")
	execCorpus(t, d.db, string(corpus))

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

// TestSyncInstance holds SyncInstance to testdata/sync/instance.json: the
// databases, and every account and role with the GRANT statements the server
// states for it. The accounts cover a role granted to a user, as a default
// role and not, privileges on a database and a table, and a host pattern; the
// image's own root account and built-in roles are there too.
func TestSyncInstance(t *testing.T) {
	testcontainer.ResetToInitial(t)
	t.Cleanup(func() { testcontainer.Reset(t, testcontainer.DropInstanceObjects) })
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{
			"CREATE DATABASE sync_inst_db",
			"CREATE TABLE sync_inst_db.t (k INT, v INT) DUPLICATE KEY(k) DISTRIBUTED BY HASH(k) BUCKETS 1 PROPERTIES (\"replication_num\" = \"1\")",
			"CREATE ROLE sync_reader COMMENT 'reads sync_inst_db'",
			"GRANT SELECT ON ALL TABLES IN DATABASE sync_inst_db TO ROLE sync_reader",
			"CREATE USER 'sync_app'@'%' IDENTIFIED BY 'sync-pw' DEFAULT ROLE sync_reader",
			"GRANT INSERT ON TABLE sync_inst_db.t TO USER 'sync_app'@'%'",
			"CREATE USER 'sync_ops'@'10.0.0.%'",
			"GRANT CREATE TABLE ON DATABASE sync_inst_db TO USER 'sync_ops'@'10.0.0.%'",
			"GRANT sync_reader TO USER 'sync_ops'@'10.0.0.%'",
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.OpenWith(t, testcontainer.Options{})
	t.Cleanup(func() { _ = db.Close() })

	inst, err := (&Driver{db: db}).SyncInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var databases []string
	for _, m := range inst.Databases {
		databases = append(databases, m.Name+" "+m.CharacterSet+" "+m.Collation)
	}
	b, err := json.MarshalIndent(map[string]any{
		"version":             inst.Version,
		"lowerCaseTableNames": inst.LowerCaseTableNames,
		"databases":           databases,
		"roles":               inst.Roles,
	}, "  ", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got := string(b) + "\n"

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

// A database the connection itself cannot reach reaches the caller as
// driver.ErrNotFound; any other client error does not.
func TestWrapSyncDatabaseErr(t *testing.T) {
	if err := wrapSyncDatabaseErr(&gomysql.MySQLError{Number: errUnknownDatabase, Message: "Unknown database"}); !errors.Is(err, driver.ErrNotFound) {
		t.Errorf("wrapSyncDatabaseErr = %v, want driver.ErrNotFound", err)
	}
	if err := wrapSyncDatabaseErr(&gomysql.MySQLError{Number: 1146, Message: "Table doesn't exist"}); errors.Is(err, driver.ErrNotFound) {
		t.Errorf("wrapSyncDatabaseErr = %v, want a plain error", err)
	}
}

func TestSyncDatabaseNotFound(t *testing.T) {
	db := testcontainer.OpenWith(t, testcontainer.Options{})
	t.Cleanup(func() { _ = db.Close() })
	d := &Driver{db: db, databaseName: "omni_no_such_database"}
	if _, err := d.SyncDatabase(context.Background()); !errors.Is(err, driver.ErrNotFound) {
		t.Fatalf("SyncDatabase missing db err = %v, want driver.ErrNotFound", err)
	}
}
