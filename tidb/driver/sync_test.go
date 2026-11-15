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

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/internal/testcontainer"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync and testdata/dump from this run")

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
	db := testcontainer.Open(t)
	t.Cleanup(func() { _ = db.Close() })
	if _, err := db.ExecContext(context.Background(), "USE "+quoteBacktick(dbName)); err != nil {
		t.Fatal(err)
	}
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

func TestSyncInstance(t *testing.T) {
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		dbRows, err := db.QueryContext(ctx, "SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE LOWER(SCHEMA_NAME) NOT IN ('information_schema', 'metrics_schema', 'mysql', 'performance_schema', 'sys', 'test')")
		if err != nil {
			return err
		}
		var extraDBs []string
		for dbRows.Next() {
			var name string
			if err := dbRows.Scan(&name); err != nil {
				_ = dbRows.Close()
				return err
			}
			extraDBs = append(extraDBs, name)
		}
		_ = dbRows.Close()
		for _, name := range extraDBs {
			if _, err := db.ExecContext(ctx, "DROP DATABASE "+quoteBacktick(name)); err != nil {
				return err
			}
		}

		uRows, err := db.QueryContext(ctx, "SELECT user, host FROM mysql.user WHERE user LIKE 'sync_%'")
		if err != nil {
			return err
		}
		var extraUsers [][2]string
		for uRows.Next() {
			var u, h string
			if err := uRows.Scan(&u, &h); err != nil {
				_ = uRows.Close()
				return err
			}
			extraUsers = append(extraUsers, [2]string{u, h})
		}
		_ = uRows.Close()
		for _, u := range extraUsers {
			if _, err := db.ExecContext(ctx, fmt.Sprintf("DROP USER %s@%s", quoteSingle(u[0]), quoteSingle(u[1]))); err != nil {
				return err
			}
		}

		for _, stmt := range []string{
			"CREATE DATABASE sync_inst_db CHARACTER SET utf8mb4 COLLATE utf8mb4_bin",
			"CREATE USER 'sync_reader'@'%' IDENTIFIED BY 'reader-pass'",
			"GRANT SELECT ON sync_inst_db.* TO 'sync_reader'@'%'",
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.Open(t)
	t.Cleanup(func() { _ = db.Close() })
	d := &Driver{db: db}

	inst, err := d.SyncInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var foundDB string
	for _, m := range inst.Databases {
		if m.Name == "sync_inst_db" {
			foundDB = fmt.Sprintf("%s %s %s", m.Name, m.CharacterSet, m.Collation)
		}
	}
	var roleList []string
	for _, r := range inst.Roles {
		if strings.HasPrefix(r.Name, "'sync_") {
			roleList = append(roleList, r.Name+": "+*r.Attribute)
		}
	}
	b, err := json.MarshalIndent(map[string]any{
		"version":             inst.Version,
		"lowerCaseTableNames": inst.LowerCaseTableNames,
		"database":            foundDB,
		"roles":               roleList,
	}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	out := string(b) + "\n"

	golden := filepath.Join("testdata", "sync", "instance.json")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, []byte(out), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(want), out); diff != "" {
		t.Errorf("SyncInstance differs from %s (-want +got):\n%s", golden, diff)
	}
}

func TestSyncDatabaseNotFound(t *testing.T) {
	db := testcontainer.Open(t)
	t.Cleanup(func() { _ = db.Close() })
	d := &Driver{db: db, databaseName: "omni_no_such_database"}
	if _, err := d.SyncDatabase(context.Background()); !errors.Is(err, driver.ErrNotFound) {
		t.Fatalf("SyncDatabase missing db err = %v, want driver.ErrNotFound", err)
	}
}

func TestSyncForeignKeyReferencedDatabase(t *testing.T) {
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{
			"DROP DATABASE IF EXISTS omni_fk_local",
			"DROP DATABASE IF EXISTS omni_fk_remote",
			"CREATE DATABASE omni_fk_local",
			"CREATE DATABASE omni_fk_remote",
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.Open(t)
	t.Cleanup(func() { _ = db.Close() })
	for _, stmt := range []string{
		"CREATE TABLE omni_fk_remote.parent (id INT PRIMARY KEY)",
		"CREATE TABLE omni_fk_local.parent (id INT PRIMARY KEY)",
		`CREATE TABLE omni_fk_local.child (
			id INT PRIMARY KEY,
			local_id INT,
			qualified_id INT,
			remote_id INT,
			CONSTRAINT fk_local FOREIGN KEY (local_id) REFERENCES omni_fk_local.parent (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
			CONSTRAINT fk_remote FOREIGN KEY (remote_id) REFERENCES omni_fk_remote.parent (id) ON DELETE CASCADE ON UPDATE RESTRICT
		)`,
	} {
		if _, err := db.ExecContext(context.Background(), stmt); err != nil {
			t.Fatal(err)
		}
	}

	d := &Driver{db: db, databaseName: "omni_fk_local"}
	fks, err := d.getForeignKeyList(context.Background(), "omni_fk_local")
	if err != nil {
		t.Fatal(err)
	}
	want := []*metadatapb.ForeignKeyMetadata{
		{Name: "fk_local", Columns: []string{"local_id"}, ReferencedTable: "parent", ReferencedColumns: []string{"id"}, OnDelete: "RESTRICT", OnUpdate: "RESTRICT", MatchType: "NONE"},
		{Name: "fk_remote", Columns: []string{"remote_id"}, ReferencedSchema: "omni_fk_remote", ReferencedTable: "parent", ReferencedColumns: []string{"id"}, OnDelete: "CASCADE", OnUpdate: "RESTRICT", MatchType: "NONE"},
	}
	if diff := cmp.Diff(want, fks["child"], protocmp.Transform()); diff != "" {
		t.Errorf("foreign keys diff (-want +got):\n%s", diff)
	}
}
