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
	"github.com/bytebase/omni/mysql/internal/testcontainer"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync and testdata/dump from this run")

var roles = []string{testcontainer.Default, testcontainer.Folded, testcontainer.MySQL84}

func openTestDriver(t *testing.T, role, dbName string) *Driver {
	t.Helper()
	testcontainer.Reset(t, role, func(ctx context.Context, db *sql.DB) error {
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
	db := testcontainer.OpenRoot(t, testcontainer.Options{
		Role:            role,
		Database:        dbName,
		MultiStatements: true,
	})
	t.Cleanup(func() { _ = db.Close() })
	return &Driver{db: db, databaseName: dbName}
}

func execCorpus(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	for _, stmt := range splitCorpus(script) {
		if _, err := db.ExecContext(context.Background(), stmt); err != nil {
			t.Fatalf("%s:\n%v", stmt, err)
		}
	}
}

// splitCorpus splits corpus.sql into top-level DDL statements without requiring
// DELIMITER directives around routine/trigger/event BEGIN ... END bodies.
func splitCorpus(script string) []string {
	var out []string
	for _, chunk := range strings.Split(script, "\n\n") {
		if s := strings.TrimSpace(chunk); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func TestSyncDatabase(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	for _, role := range roles {
		t.Run(role, func(t *testing.T) {
			d := openTestDriver(t, role, "omni_sync_corpus")
			execCorpus(t, d.db, string(corpus))
			got, err := d.SyncDatabase(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			golden := filepath.Join("testdata", "sync", role+".json")
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
				t.Errorf("SyncDatabase on %s differs from %s (-want +got):\n%s", role, golden, diff)
			}
		})
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
	got := map[string]json.RawMessage{}
	for _, role := range roles {
		t.Run(role, func(t *testing.T) {
			db := testcontainer.OpenEmptyWith(t, testcontainer.Options{
				Role:            role,
				Database:        "test",
				MultiStatements: true,
			})
			t.Cleanup(func() { _ = db.Close() })
			d := &Driver{db: db, databaseName: "test"}

			if _, err := db.ExecContext(context.Background(), `
				CREATE DATABASE sync_extra CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
				CREATE USER 'sync_reader'@'%' IDENTIFIED BY 'reader-pass';
				GRANT SELECT ON sync_extra.* TO 'sync_reader'@'%';
				CREATE USER 'sync''quote'@'%' IDENTIFIED BY 'quote-pass';
			`); err != nil {
				t.Fatal(err)
			}

			inst, err := d.SyncInstance(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			dbs := make([]string, 0, len(inst.Databases))
			for _, m := range inst.Databases {
				dbs = append(dbs, fmt.Sprintf("%s %s %s", m.Name, m.CharacterSet, m.Collation))
			}
			roleList := make([]string, 0, len(inst.Roles))
			for _, r := range inst.Roles {
				if strings.HasPrefix(r.Name, "'sync") {
					roleList = append(roleList, r.Name+": "+*r.Attribute)
				}
			}
			b, err := json.MarshalIndent(map[string]any{
				"version":             inst.Version,
				"lowerCaseTableNames": inst.LowerCaseTableNames,
				"databases":           dbs,
				"roles":               roleList,
			}, "  ", "  ")
			if err != nil {
				t.Fatal(err)
			}
			got[role] = b
		})
	}

	var buf strings.Builder
	buf.WriteString("{\n")
	for i, role := range roles {
		if i > 0 {
			buf.WriteString(",\n")
		}
		fmt.Fprintf(&buf, "  %q: %s", role, got[role])
	}
	buf.WriteString("\n}\n")

	golden := filepath.Join("testdata", "sync", "instance.json")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, []byte(buf.String()), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(want), buf.String()); diff != "" {
		t.Errorf("SyncInstance differs from %s (-want +got):\n%s", golden, diff)
	}
}

func TestSyncDatabaseNotFound(t *testing.T) {
	db := testcontainer.OpenRoot(t, testcontainer.Options{Role: testcontainer.Default})
	t.Cleanup(func() { _ = db.Close() })
	d := &Driver{db: db, databaseName: "omni_no_such_database"}
	if _, err := d.SyncDatabase(context.Background()); !errors.Is(err, driver.ErrNotFound) {
		t.Fatalf("SyncDatabase missing db err = %v, want driver.ErrNotFound", err)
	}
}

func TestSyncForeignKeyReferencedDatabase(t *testing.T) {
	testcontainer.Reset(t, testcontainer.Default, func(ctx context.Context, db *sql.DB) error {
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
	db := testcontainer.OpenRoot(t, testcontainer.Options{
		Role:            testcontainer.Default,
		Database:        "omni_fk_local",
		MultiStatements: true,
	})
	t.Cleanup(func() { _ = db.Close() })
	if _, err := db.ExecContext(context.Background(), `
		CREATE TABLE omni_fk_remote.parent (id INT PRIMARY KEY);
		CREATE TABLE omni_fk_local.parent (id INT PRIMARY KEY);
		CREATE TABLE omni_fk_local.child (
			id INT PRIMARY KEY,
			local_id INT,
			qualified_id INT,
			remote_id INT,
			CONSTRAINT fk_local FOREIGN KEY (local_id) REFERENCES parent (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
			CONSTRAINT fk_qualified FOREIGN KEY (qualified_id) REFERENCES omni_fk_local.parent (id) ON DELETE SET NULL ON UPDATE RESTRICT,
			CONSTRAINT fk_remote FOREIGN KEY (remote_id) REFERENCES omni_fk_remote.parent (id) ON DELETE CASCADE ON UPDATE RESTRICT
		);
	`); err != nil {
		t.Fatal(err)
	}

	d := &Driver{db: db, databaseName: "omni_fk_local"}
	fks, err := d.getForeignKeyList(context.Background(), "omni_fk_local")
	if err != nil {
		t.Fatal(err)
	}
	want := []*metadatapb.ForeignKeyMetadata{
		{Name: "fk_local", Columns: []string{"local_id"}, ReferencedTable: "parent", ReferencedColumns: []string{"id"}, OnDelete: "RESTRICT", OnUpdate: "RESTRICT", MatchType: "NONE"},
		{Name: "fk_qualified", Columns: []string{"qualified_id"}, ReferencedTable: "parent", ReferencedColumns: []string{"id"}, OnDelete: "SET NULL", OnUpdate: "RESTRICT", MatchType: "NONE"},
		{Name: "fk_remote", Columns: []string{"remote_id"}, ReferencedSchema: "omni_fk_remote", ReferencedTable: "parent", ReferencedColumns: []string{"id"}, OnDelete: "CASCADE", OnUpdate: "RESTRICT", MatchType: "NONE"},
	}
	if diff := cmp.Diff(want, fks["child"], protocmp.Transform()); diff != "" {
		t.Errorf("foreign keys diff (-want +got):\n%s", diff)
	}
}

func TestColumnDefaultsAndHelpers(t *testing.T) {
	cases := []struct {
		name       string
		colType    string
		def        sql.NullString
		nullable   bool
		extra      string
		binFmt     binaryDefaultFormat
		wantDef    string
		wantUpdate string
	}{
		{"auto_increment", "int", sql.NullString{}, false, "auto_increment", binaryDefaultHexNotation, "AUTO_INCREMENT", ""},
		{"nullable null", "varchar(20)", sql.NullString{}, true, "", binaryDefaultHexNotation, "NULL", ""},
		{"not null no default", "int", sql.NullString{}, false, "", binaryDefaultHexNotation, "", ""},
		{"quoted string", "varchar(20)", sql.NullString{Valid: true, String: "'hello'"}, true, "", binaryDefaultHexNotation, "'hello'", ""},
		{"current_timestamp(6) on update", "datetime(6)", sql.NullString{Valid: true, String: "'CURRENT_TIMESTAMP(6)'"}, false, "DEFAULT_GENERATED on update CURRENT_TIMESTAMP(6)", binaryDefaultHexNotation, "CURRENT_TIMESTAMP(6)", "CURRENT_TIMESTAMP(6)"},
		{"expression default", "varchar(36)", sql.NullString{Valid: true, String: `'json_object(\\'k\\',1)'`}, true, "DEFAULT_GENERATED", binaryDefaultHexNotation, "(json_object('k',1))", ""},
		{"8.0 binary printable", "binary(4)", sql.NullString{Valid: true, String: "'0x6162'"}, false, "", binaryDefaultHexNotation, "'ab'", ""},
		{"8.0 varbinary non-printable", "varbinary(8)", sql.NullString{Valid: true, String: "'0x0100fe'"}, false, "", binaryDefaultHexNotation, "0x0100fe", ""},
		{"5.7 binary NUL-padded", "binary(4)", sql.NullString{Valid: true, String: "'ab\\0\\0'"}, false, "", binaryDefaultRawBytes, "'ab'", ""},
	}
	for _, tc := range cases {
		col := &metadatapb.ColumnMetadata{Type: tc.colType}
		setColumnMetadataDefault(col, tc.def, tc.nullable, tc.extra, tc.binFmt)
		if col.Default != tc.wantDef || col.OnUpdate != tc.wantUpdate {
			t.Errorf("%s: got Default=%q OnUpdate=%q, want %q / %q", tc.name, col.Default, col.OnUpdate, tc.wantDef, tc.wantUpdate)
		}
	}

	if got := unescapeCheckClause("(`status` in (_utf8mb4\\'open\\',_utf8mb4\\'closed\\'))"); got != "(`status` in (_utf8mb4'open',_utf8mb4'closed'))" {
		t.Errorf("unescapeCheckClause = %q", got)
	}
	fnDef := "CREATE FUNCTION `f`(p_a VARCHAR(10) CHARSET utf8mb4) RETURNS char(36) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci\nRETURN p_a"
	wantFn := "CREATE FUNCTION `f`(p_a VARCHAR(10) CHARSET utf8mb4) RETURNS char(36)\nRETURN p_a"
	if got := stripReturnsCharset(fnDef); got != wantFn {
		t.Errorf("stripReturnsCharset = %q, want %q", got, wantFn)
	}
	if iso, err := utf8ToISO88591("caf\u00e9"); err != nil || iso != "caf\xe9" {
		t.Errorf("utf8ToISO88591 Latin-1 = %q, %v", iso, err)
	}
	if _, err := utf8ToISO88591("日本語"); err == nil {
		t.Error("utf8ToISO88591 accepted non-Latin-1 rune")
	}
}
