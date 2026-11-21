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

func TestDumpRestoresCheckNotEnforcedAndViewOptions(t *testing.T) {
	const setup = `CREATE TABLE t_chk (
	id INT PRIMARY KEY,
	x INT NOT NULL,
	note VARCHAR(20) NOT SECONDARY,
	CONSTRAINT chk_enf CHECK (x > 0),
	CONSTRAINT chk_ne CHECK (x < 10) NOT ENFORCED
);

CREATE ALGORITHM=MERGE SQL SECURITY INVOKER VIEW v_opts AS
	SELECT id, x FROM t_chk WHERE x > 0 WITH LOCAL CHECK OPTION;

CREATE TABLE part_ts (
	id INT NOT NULL,
	yr INT NOT NULL,
	PRIMARY KEY (id, yr)
) PARTITION BY RANGE (yr) (
	PARTITION p0 VALUES LESS THAN (2025) TABLESPACE = innodb_file_per_table,
	PARTITION p1 VALUES LESS THAN MAXVALUE
);`
	for _, role := range []string{testcontainer.Default, testcontainer.MySQL84} {
		t.Run(role, func(t *testing.T) {
			source, restored, err := dumpAndRestore(t, role, "omni_sync_options_roundtrip", setup)
			if err != nil {
				t.Fatalf("restoring dump on %s: %v", role, err)
			}
			tables := map[string]*metadatapb.TableMetadata{}
			for _, tbl := range source.Schemas[0].Tables {
				tables[tbl.Name] = tbl
			}
			chk := tables["t_chk"].CheckConstraints
			if len(chk) != 2 || chk[0].NotEnforced || !chk[1].NotEnforced {
				t.Fatalf("unexpected check constraints on %s: %+v", role, chk)
			}
			var noteNotSecondary bool
			for _, col := range tables["t_chk"].Columns {
				if col.Name == "note" {
					noteNotSecondary = col.NotSecondary
				}
			}
			if !noteNotSecondary {
				t.Fatalf("expected t_chk.note NotSecondary=true on %s, got %+v", role, tables["t_chk"].Columns)
			}
			parts := tables["part_ts"].Partitions
			if len(parts) != 2 || parts[0].GetTablespace() != "innodb_file_per_table" || parts[1].GetTablespace() != "innodb_file_per_table" {
				t.Fatalf("unexpected partition tablespaces on %s: %+v", role, parts)
			}
			if len(source.Schemas[0].Views) != 1 {
				t.Fatalf("unexpected views on %s: %+v", role, source.Schemas[0].Views)
			}
			v := source.Schemas[0].Views[0]
			if v.Algorithm != "MERGE" || v.SqlSecurity != "INVOKER" || v.CheckOption != "LOCAL" {
				t.Fatalf("unexpected view options on %s: %+v", role, v)
			}
			if strings.Contains(strings.ToUpper(v.Definition), "CHECK OPTION") {
				t.Fatalf("expected normalized Definition without CHECK OPTION, got %q", v.Definition)
			}
			if diff := restoreDiff(source, restored); diff != "" {
				t.Errorf("restored database on %s differs from source (-source +restored):\n%s", role, diff)
			}
		})
	}
}

// TestDumpRestoresTablespacesKeyBlockSizesAndKeyAlgorithms pins what only
// SHOW CREATE TABLE states: a table's general or system tablespace, an index's
// KEY_BLOCK_SIZE, and KEY partitioning's ALGORITHM, through sync, dump and
// restore.
func TestDumpRestoresTablespacesKeyBlockSizesAndKeyAlgorithms(t *testing.T) {
	const setup = `CREATE TABLESPACE ts_sync ADD DATAFILE 'ts_sync.ibd';

CREATE TABLE t_ts (id INT PRIMARY KEY) TABLESPACE ts_sync;

CREATE TABLE ` + "`t-ts`" + ` (id INT PRIMARY KEY) TABLESPACE ts_sync;

CREATE TABLE t_sys (id INT PRIMARY KEY) TABLESPACE innodb_system;

CREATE TABLE t_plain (id INT PRIMARY KEY);

CREATE TABLE t_kbs (a INT, b INT, KEY ka (a) KEY_BLOCK_SIZE=4, KEY kb (b)) ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8;

CREATE TABLE t_myisam (a INT, KEY ka (a) KEY_BLOCK_SIZE=1024) ENGINE=MyISAM;

CREATE TABLE p_key (a INT) PARTITION BY KEY ALGORITHM=1 (a) PARTITIONS 2;

CREATE TABLE p_sub (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY KEY ALGORITHM=1 (b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10));

CREATE TABLE p_quoted (` + "`a,b`" + ` INT, c INT) PARTITION BY KEY (` + "`a,b`" + `, c) PARTITIONS 2;

CREATE TABLE t_fn (code VARCHAR(10), KEY k1 ((CONCAT(code, 'a\'b'))), KEY k2 ((CONCAT(code, 'c\\d'))));

CREATE TABLE p_hash_named (id INT) PARTITION BY HASH(id) PARTITIONS 2 (PARTITION pa COMMENT='first', PARTITION pb);

CREATE TABLE p_key_named (id INT) PARTITION BY KEY(id) PARTITIONS 2 (PARTITION ka, PARTITION kb);`
	for _, role := range []string{testcontainer.Default, testcontainer.MySQL84} {
		t.Run(role, func(t *testing.T) {
			const dbName = "omni_sync_show_create_roundtrip"
			source, restored, err := dumpAndRestore(t, role, dbName, setup)
			if err != nil {
				t.Fatalf("restoring dump on %s: %v", role, err)
			}
			tables := map[string]*metadatapb.TableMetadata{}
			for _, tbl := range source.Schemas[0].Tables {
				tables[tbl.Name] = tbl
			}
			for name, want := range map[string]string{"t_ts": "ts_sync", "t-ts": "ts_sync", "t_sys": "innodb_system", "t_plain": ""} {
				if got := tables[name].GetTablespace(); got != want {
					t.Errorf("%s tablespace on %s = %q, want %q", name, role, got, want)
				}
			}
			kbs := map[string]int64{}
			for _, tbl := range []string{"t_kbs", "t_myisam"} {
				for _, idx := range tables[tbl].Indexes {
					kbs[tbl+"."+idx.Name] = idx.KeyBlockSize
				}
			}
			if kbs["t_kbs.ka"] != 4 || kbs["t_kbs.kb"] != 0 || kbs["t_myisam.ka"] != 1024 {
				t.Errorf("index KEY_BLOCK_SIZEs on %s = %v", role, kbs)
			}
			for _, p := range tables["p_key"].Partitions {
				if p.Algorithm != 1 {
					t.Errorf("p_key partition %s ALGORITHM on %s = %d, want 1", p.Name, role, p.Algorithm)
				}
			}
			for _, p := range tables["p_sub"].Partitions {
				for _, sp := range p.Subpartitions {
					if sp.Algorithm != 1 {
						t.Errorf("p_sub subpartition %s ALGORITHM on %s = %d, want 1", sp.Name, role, sp.Algorithm)
					}
				}
			}
			for name, want := range map[string]string{"p_hash_named": "pa,pb", "p_key_named": "ka,kb"} {
				var names []string
				for _, p := range tables[name].Partitions {
					names = append(names, p.Name)
				}
				if strings.Join(names, ",") != want {
					t.Errorf("%s partitions on %s = %v, want %s", name, role, names, want)
				}
			}
			if diff := restoreDiff(source, restored); diff != "" {
				t.Errorf("restored database on %s differs from source (-source +restored):\n%s", role, diff)
			}
			// The tablespace is the server's; drop it with the tables in it.
			cleanup := openTestDriver(t, role, dbName)
			if _, err := cleanup.db.ExecContext(context.Background(), "DROP TABLESPACE ts_sync"); err != nil {
				t.Fatalf("DROP TABLESPACE ts_sync on %s: %v", role, err)
			}
		})
	}
}
