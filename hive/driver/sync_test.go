package driver

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bytebase/omni/hive/internal/testcontainer"
	metadatapb "github.com/bytebase/omni/metadata"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync from this run")

// openTestCorpus returns a driver on dbName, dropped and created empty first so
// the tape holds one run's requests whatever the test binary ran before, with
// corpus.sql applied.
func openTestCorpus(t *testing.T, dbName string) *Driver {
	t.Helper()
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{
			"DROP DATABASE IF EXISTS " + quoteBacktick(dbName) + " CASCADE",
			"CREATE DATABASE " + quoteBacktick(dbName),
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.Open(t, dbName)
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	execCorpus(t, db, string(corpus))
	return &Driver{db: db, databaseName: dbName}
}

// execCorpus runs a fixture script, one blank-line-separated statement at a
// time, as a client would.
func execCorpus(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	for _, chunk := range strings.Split(script, "\n\n") {
		stmt := strings.TrimSpace(chunk)
		if stmt == "" {
			continue
		}
		if _, err := db.ExecContext(context.Background(), stmt); err != nil {
			t.Fatalf("%s:\n%v", stmt, err)
		}
	}
}

// TestSyncInstance reads the server's version and database list. The pinned
// server starts with one database, default, so the golden names two.
func TestSyncInstance(t *testing.T) {
	testcontainer.ResetToInitial(t)
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		_, err := db.ExecContext(ctx, "CREATE DATABASE sync_inst_db")
		return err
	})
	db := testcontainer.Open(t, "")
	defer func() { _ = db.Close() }()

	instance, err := (&Driver{db: db}).SyncInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var databases []string
	for _, m := range instance.Databases {
		databases = append(databases, m.Name)
	}
	got := "version: " + instance.Version + "\n" +
		"databases: " + strings.Join(databases, ", ") + "\n"
	want := "version: 4.0.1\n" +
		"databases: default, sync_inst_db\n"
	if got != want {
		t.Errorf("SyncInstance =\n%s\nwant\n%s", got, want)
	}
}

// TestSyncDatabase reads the corpus's tables: an external table, a managed
// table Hive translated to an external one and then back, a partitioned
// managed table with one partition and two rows, and a view.
func TestSyncDatabase(t *testing.T) {
	d := openTestCorpus(t, "omni_sync")

	got, err := d.SyncDatabase(context.Background())
	require.NoError(t, err)
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
