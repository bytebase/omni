package driver

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
)

func TestDump(t *testing.T) {
	for _, role := range roles {
		t.Run(role, func(t *testing.T) {
			m := readGolden(t, filepath.Join("testdata", "sync", role+".json"))
			before := proto.Clone(m)
			var b strings.Builder
			if err := (&Driver{}).Dump(context.Background(), &b, m); err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(before, m) {
				t.Errorf("Dump mutated its input metadata")
			}
			golden := filepath.Join("testdata", "dump", role+".sql")
			if *update {
				if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
					t.Fatal(err)
				}
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
				t.Errorf("Dump on %s differs from %s (-want +got):\n%s", role, golden, diff)
			}
		})
	}
}

const gapMarker = "-- not restorable: TestDumpGaps"

func restorable(corpus string) string {
	var kept []string
	for _, line := range strings.Split(corpus, "\n") {
		if !strings.HasSuffix(line, gapMarker) {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

// TestDumpRestores syncs the restorable subset of corpus.sql on each server,
// dumps it, recreates the database, applies the dump through a
// mysql-client-style DELIMITER splitter, and checks that resyncing yields the
// same schema.
func TestDumpRestores(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	setup := restorable(string(corpus))
	for _, role := range roles {
		t.Run(role, func(t *testing.T) {
			source, restored, err := dumpAndRestore(t, role, "omni_dump_roundtrip", setup)
			if err != nil {
				t.Fatalf("restoring dump on %s: %v", role, err)
			}
			if diff := restoreDiff(source, restored); diff != "" {
				t.Errorf("restored database on %s differs from source (-source +restored):\n%s", role, diff)
			}
		})
	}
}

func TestDumpRestoresStringDefaultsAndPartitionComments(t *testing.T) {
	const setup = `CREATE TABLE literal_defaults (
	id INT PRIMARY KEY,
	a VARCHAR(40) NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
	b VARCHAR(40) NOT NULL DEFAULT 'current_date_label',
	ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) COMMENT='PARTITIONS 9 SUBPARTITIONS 7';

CREATE TABLE hash_commented (
	id INT NOT NULL PRIMARY KEY COMMENT 'PARTITIONS 9'
) COMMENT='PARTITIONS 11' PARTITION BY HASH (id) PARTITIONS 2;

CREATE TABLE range_commented (
	id INT NOT NULL PRIMARY KEY
) COMMENT='PARTITIONS 9 SUBPARTITIONS 7'
PARTITION BY RANGE (id) SUBPARTITION BY HASH (id) SUBPARTITIONS 2 (
	PARTITION p0 VALUES LESS THAN (100),
	PARTITION p1 VALUES LESS THAN (MAXVALUE)
);`
	for _, role := range []string{testcontainer.Default, testcontainer.MySQL84} {
		t.Run(role, func(t *testing.T) {
			source, restored, err := dumpAndRestore(t, role, "omni_sync_defect_roundtrip", setup)
			if err != nil {
				t.Fatalf("restoring dump on %s: %v", role, err)
			}
			if diff := restoreDiff(source, restored); diff != "" {
				t.Errorf("restored database on %s differs from source (-source +restored):\n%s", role, diff)
			}
		})
	}
}

// TestDumpGaps holds what a dump cannot faithfully restore yet; each entry
// fails once its gap closes so its corpus line moves into TestDumpRestores.
func TestDumpGaps(t *testing.T) {
	for _, gap := range []struct {
		name, db, setup, wantErr string
	}{
		{
			name: "view explicit column list",
			db:   "omni_gap_view_cols",
			setup: "CREATE TABLE t (id INT PRIMARY KEY);\n\n" +
				"CREATE VIEW v1 (`renamed_id`) AS SELECT id FROM t;\n\n" +
				"CREATE VIEW v2 AS SELECT renamed_id FROM v1;",
			wantErr: "Unknown column",
		},
		{
			name:  "index comment",
			db:    "omni_gap_idx_comment",
			setup: "CREATE TABLE t (id INT PRIMARY KEY, v INT, KEY idx_v (v) COMMENT 'note');",
		},
	} {
		t.Run(gap.name, func(t *testing.T) {
			source, restored, err := dumpAndRestore(t, testcontainer.Default, gap.db, gap.setup)
			if gap.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), gap.wantErr) {
					t.Fatalf("restore err = %v, want %q; if the gap closed, remove it", err, gap.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("restore failed: %v", err)
			}
			if diff := restoreDiff(source, restored); diff == "" {
				t.Fatal("restore reproduced source; if the gap closed, remove it")
			}
		})
	}
}

func restoreDiff(source, restored *metadatapb.DatabaseSchemaMetadata) string {
	return cmp.Diff(source, restored, protocmp.Transform(),
		protocmp.IgnoreFields(&metadatapb.TableMetadata{}, "data_size", "index_size", "data_free", "row_count"),
	)
}

func dumpAndRestore(t *testing.T, role, dbName, setup string) (*metadatapb.DatabaseSchemaMetadata, *metadatapb.DatabaseSchemaMetadata, error) {
	t.Helper()
	ctx := context.Background()
	d := openTestDriver(t, role, dbName)
	execCorpus(t, d.db, setup)

	source, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var dumpBuf strings.Builder
	if err := d.Dump(ctx, &dumpBuf, source); err != nil {
		t.Fatal(err)
	}

	for _, stmt := range []string{
		"DROP DATABASE " + quoteBacktick(dbName),
		"CREATE DATABASE " + quoteBacktick(dbName),
		"USE " + quoteBacktick(dbName),
	} {
		if _, err := d.db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	for _, stmt := range splitLikeMySQLClient(dumpBuf.String()) {
		if _, err := d.db.ExecContext(ctx, stmt); err != nil {
			return source, nil, err
		}
	}

	restored, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return source, restored, nil
}

func splitLikeMySQLClient(script string) []string {
	const directive = "DELIMITER "
	var stmts []string
	var cur strings.Builder
	flush := func() {
		if stmt := strings.TrimSpace(cur.String()); stmt != "" {
			stmts = append(stmts, stmt)
		}
		cur.Reset()
	}
	delimiter := ";"
	var quote byte
	for i := 0; i < len(script); {
		c := script[i]
		if quote != 0 {
			cur.WriteByte(c)
			switch {
			case c == '\\' && quote != '`' && i+1 < len(script):
				cur.WriteByte(script[i+1])
				i += 2
				continue
			case c == quote:
				quote = 0
			}
			i++
			continue
		}
		lineEnd := strings.IndexByte(script[i:], '\n')
		if lineEnd < 0 {
			lineEnd = len(script) - i
		}
		switch {
		case strings.TrimSpace(cur.String()) == "" && len(script)-i >= len(directive) && strings.EqualFold(script[i:i+len(directive)], directive):
			delimiter = strings.TrimSpace(script[i+len(directive) : i+lineEnd])
			cur.Reset()
			i += lineEnd
		case strings.HasPrefix(script[i:], delimiter):
			flush()
			i += len(delimiter)
		case c == '#' || strings.HasPrefix(script[i:], "--") && (i+2 == len(script) || script[i+2] <= ' '):
			i += lineEnd
		case strings.HasPrefix(script[i:], "/*") && !strings.HasPrefix(script[i:], "/*!"):
			end := strings.Index(script[i+2:], "*/")
			if end < 0 {
				i = len(script)
			} else {
				i += end + 4
			}
		default:
			if c == '\'' || c == '"' || c == '`' {
				quote = c
			}
			cur.WriteByte(c)
			i++
		}
	}
	flush()
	return stmts
}
