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

func TestDumpRestores(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for _, role := range roles {
		t.Run(role, func(t *testing.T) {
			const dbName = "omni_dump_roundtrip"
			d := openTestDriver(t, role, dbName)
			execCorpus(t, d.db, string(corpus))
			if role == "default" {
				execCorpus(t, d.db, "CREATE SEQUENCE seq_typed AS INT UNSIGNED START WITH 7 MAXVALUE 4000000000;")
			}

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
					t.Fatalf("restoring statement on %s:\n%s\nerr: %v", role, stmt, err)
				}
			}

			restored, err := d.SyncDatabase(ctx)
			if err != nil {
				t.Fatal(err)
			}
			diff := cmp.Diff(source, restored, protocmp.Transform(),
				protocmp.IgnoreFields(&metadatapb.TableMetadata{}, "data_size", "index_size", "data_free", "row_count"),
			)
			if diff != "" {
				t.Errorf("restored database on %s differs from source (-source +restored):\n%s", role, diff)
			}
		})
	}
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
