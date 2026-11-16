package driver

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/hive/internal/testcontainer"
)

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

// TestSyncDatabaseReadsDescribeFormatted pins what the moved sync does on the
// pinned server today: DESCRIBE FORMATTED answers NULL for the cells of a
// section's blank separator row, and the scan of three strings stops there, so
// the sync reports the failure instead of a schema. The fix(hive) PR that
// follows the move reads those cells and records the schema as a golden.
func TestSyncDatabaseReadsDescribeFormatted(t *testing.T) {
	d := openTestCorpus(t, "omni_sync")

	_, err := d.SyncDatabase(context.Background())
	if err == nil {
		t.Fatal("SyncDatabase succeeded; the moved sync scans DESCRIBE FORMATTED's NULL cells into strings")
	}
	for _, want := range []string{
		"failed to get info from database omni_sync",
		"failed to get table metadata from database omni_sync",
		"failed to describe table events's type",
		"failed to scan row",
		"converting NULL to string is unsupported",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("SyncDatabase error %q does not hold %q", err, want)
		}
	}
}
