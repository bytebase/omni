package driver

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

func TestDump(t *testing.T) {
	m := readGolden(t, filepath.Join("testdata", "sync", "default.json"))
	before := proto.Clone(m)
	var b strings.Builder
	if err := (&Driver{}).Dump(context.Background(), &b, m); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(before, m) {
		t.Errorf("Dump mutated its input metadata")
	}
	golden := filepath.Join("testdata", "dump", "default.sql")
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
		t.Errorf("Dump differs from %s (-want +got):\n%s", golden, diff)
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

func TestDumpRestores(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	source, restored, err := dumpAndRestore(t, "omni_dump_roundtrip", restorable(string(corpus)))
	if err != nil {
		t.Fatalf("restoring dump: %v", err)
	}
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("restored database differs from source (-source +restored):\n%s", diff)
	}
}

// TestDumpGaps pins OceanBase objects that a dump cannot faithfully restore
// yet: external tables.
func TestDumpGaps(t *testing.T) {
	for _, gap := range []struct {
		name, db, setup, wantErr string
	}{
		{
			name:  "external table",
			db:    "omni_gap_ext",
			setup: "CREATE EXTERNAL TABLE ext_t (c1 INT) LOCATION = '/tmp' FORMAT = (TYPE = 'CSV');",
		},
	} {
		t.Run(gap.name, func(t *testing.T) {
			source, restored, err := dumpAndRestore(t, gap.db, gap.setup)
			if gap.wantErr != "" {
				if err == nil || !strings.Contains(strings.ToLower(err.Error()), gap.wantErr) {
					t.Fatalf("restore err = %v, want %q; if the gap closed, remove it", err, gap.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("restore failed: %v", err)
			}
			if restoreDiff(source, restored) == "" {
				t.Fatalf("gap %s closed; move back into TestDumpRestores", gap.name)
			}
		})
	}
}

func restoreDiff(source, restored *metadatapb.DatabaseSchemaMetadata) string {
	return cmp.Diff(source, restored, protocmp.Transform(),
		protocmp.IgnoreFields(&metadatapb.DatabaseSchemaMetadata{}, "name"),
		protocmp.IgnoreFields(&metadatapb.TableMetadata{}, "data_size", "index_size", "data_free", "row_count"),
	)
}

func dumpAndRestore(t *testing.T, dbName, setup string) (*metadatapb.DatabaseSchemaMetadata, *metadatapb.DatabaseSchemaMetadata, error) {
	t.Helper()
	ctx := context.Background()
	restoredName := dbName + "_restored"
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{
			"DROP DATABASE IF EXISTS " + quoteBacktick(dbName),
			"DROP DATABASE IF EXISTS " + quoteBacktick(restoredName),
			"CREATE DATABASE " + quoteBacktick(dbName),
			"CREATE DATABASE " + quoteBacktick(restoredName),
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.OpenWith(t, testcontainer.Options{
		Database:        dbName,
		MultiStatements: true,
	})
	t.Cleanup(func() { _ = db.Close() })

	src := &Driver{db: db, databaseName: dbName}
	execCorpus(t, src.db, setup)

	source, err := src.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var dumpBuf strings.Builder
	if err := src.Dump(ctx, &dumpBuf, source); err != nil {
		t.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, "USE "+quoteBacktick(restoredName)); err != nil {
		t.Fatal(err)
	}
	// Replace source-database qualifiers in view bodies so views bind to the
	// restored database rather than the source database.
	rewritten := strings.ReplaceAll(dumpBuf.String(), quoteBacktick(dbName)+".", quoteBacktick(restoredName)+".")
	for _, stmt := range splitLikeMySQLClient(rewritten) {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return source, nil, err
		}
	}

	dst := &Driver{db: db, databaseName: restoredName}
	restored, err := dst.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range restored.Schemas[0].Views {
		v.Definition = strings.ReplaceAll(v.Definition, quoteBacktick(restoredName)+".", quoteBacktick(dbName)+".")
	}
	for _, v := range restored.Schemas[0].MaterializedViews {
		v.Definition = strings.ReplaceAll(v.Definition, quoteBacktick(restoredName)+".", quoteBacktick(dbName)+".")
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
