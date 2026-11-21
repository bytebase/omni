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
	"github.com/bytebase/omni/tidb/internal/testcontainer"
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
	source, restored := dumpAndRestore(t, "omni_dump_roundtrip", restorable(string(corpus)))
	diff := restoreDiff(source, restored)
	if diff != "" {
		t.Errorf("restored database differs from source (-source +restored):\n%s", diff)
	}
}

func restoreDiff(source, restored *metadatapb.DatabaseSchemaMetadata) string {
	return cmp.Diff(source, restored, protocmp.Transform(),
		protocmp.IgnoreFields(&metadatapb.TableMetadata{}, "data_size", "index_size", "data_free", "row_count"),
	)
}

func dumpAndRestore(t *testing.T, dbName, setup string) (*metadatapb.DatabaseSchemaMetadata, *metadatapb.DatabaseSchemaMetadata) {
	t.Helper()
	ctx := context.Background()
	d := openTestDriver(t, dbName)
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
			t.Fatalf("restoring statement:\n%s\nerr: %v", stmt, err)
		}
	}

	restored, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return source, restored
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
		case strings.HasPrefix(script[i:], "/*") && !strings.HasPrefix(script[i:], "/*!") && !strings.HasPrefix(script[i:], "/*T!"):
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
	ctx := context.Background()
	d := openTestDriver(t, "omni_options_roundtrip")
	if _, err := d.db.ExecContext(ctx, "SET GLOBAL tidb_enable_check_constraint = 1"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = d.db.ExecContext(ctx, "SET GLOBAL tidb_enable_check_constraint = 0")
	}()

	const setup = `CREATE TABLE t_chk (
	id INT PRIMARY KEY,
	x INT NOT NULL,
	CONSTRAINT chk_enf CHECK (x > 0),
	CONSTRAINT chk_ne CHECK (x < 10) NOT ENFORCED
);

CREATE TABLE t_ttl (
	id INT PRIMARY KEY,
	ts DATETIME NOT NULL
) TTL = ts + INTERVAL 1 DAY TTL_ENABLE = 'OFF' TTL_JOB_INTERVAL = '2h';

CREATE ALGORITHM=MERGE SQL SECURITY INVOKER VIEW v_opts AS
	SELECT id, x FROM t_chk WHERE x > 0 WITH LOCAL CHECK OPTION;`
	execCorpus(t, d.db, setup)

	source, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	tables := map[string]*metadatapb.TableMetadata{}
	for _, tbl := range source.Schemas[0].Tables {
		tables[tbl.Name] = tbl
	}
	if tables["t_chk"] == nil || len(tables["t_chk"].CheckConstraints) != 2 {
		t.Fatalf("unexpected TiDB check constraints: %+v", source.Schemas[0].Tables)
	}
	chk := tables["t_chk"].CheckConstraints
	if chk[0].NotEnforced || !chk[1].NotEnforced {
		t.Fatalf("unexpected TiDB NotEnforced flags: %+v", chk)
	}
	ttl := tables["t_ttl"]
	if ttl == nil || ttl.TtlColumn != "ts" || ttl.TtlInterval != "interval 1 day" || ttl.TtlEnable || ttl.TtlJobInterval != "2h" {
		t.Fatalf("unexpected TiDB TTL metadata: %+v", ttl)
	}
	if len(source.Schemas[0].Views) != 1 {
		t.Fatalf("unexpected TiDB views: %+v", source.Schemas[0].Views)
	}
	v := source.Schemas[0].Views[0]
	if v.Algorithm != "MERGE" || v.SqlSecurity != "INVOKER" || v.CheckOption != "LOCAL" {
		t.Fatalf("unexpected TiDB view options: %+v", v)
	}

	var dumpBuf strings.Builder
	if err := d.Dump(ctx, &dumpBuf, source); err != nil {
		t.Fatal(err)
	}
	for _, stmt := range []string{
		"DROP DATABASE `omni_options_roundtrip`",
		"CREATE DATABASE `omni_options_roundtrip`",
		"USE `omni_options_roundtrip`",
	} {
		if _, err := d.db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	for _, stmt := range splitLikeMySQLClient(dumpBuf.String()) {
		if _, err := d.db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("restoring statement:\n%s\nerr: %v", stmt, err)
		}
	}
	restored, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("restored TiDB database differs from source (-source +restored):\n%s", diff)
	}
}

func TestDumpRestoresExtendedTiDBObjects(t *testing.T) {
	ctx := context.Background()
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{
			"DROP DATABASE IF EXISTS `omni_ext_roundtrip`",
			"DROP PLACEMENT POLICY IF EXISTS `p_sync_ext`",
			"CREATE PLACEMENT POLICY `p_sync_ext` PRIMARY_REGION='us-east' REGIONS='us-east'",
			"CREATE DATABASE `omni_ext_roundtrip` PLACEMENT POLICY = `p_sync_ext`",
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.Open(t)
	t.Cleanup(func() { _ = db.Close() })
	if _, err := db.ExecContext(ctx, "USE `omni_ext_roundtrip`"); err != nil {
		t.Fatal(err)
	}
	d := &Driver{db: db, databaseName: "omni_ext_roundtrip"}
	if _, err := d.db.ExecContext(ctx, "SET GLOBAL tidb_enable_check_constraint = 1"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = d.db.ExecContext(ctx, "SET GLOBAL tidb_enable_check_constraint = 0")
	}()

	const setup = `CREATE TABLE t_chk_a (
	id INT PRIMARY KEY,
	a INT NOT NULL,
	CONSTRAINT chk_a CHECK (a > 0)
);

CREATE TABLE t_chk_b (
	id INT PRIMARY KEY,
	b INT NOT NULL,
	CONSTRAINT chk_b CHECK (b < 100) NOT ENFORCED
);

CREATE GLOBAL TEMPORARY TABLE gtt_del (
	id INT PRIMARY KEY
) ON COMMIT DELETE ROWS;

CREATE TABLE t_part_global (
	id INT NOT NULL,
	k INT NOT NULL,
	PRIMARY KEY (id) NONCLUSTERED,
	UNIQUE KEY uk_k (k) /*T![global_index] GLOBAL */
) PARTITION BY HASH (id) PARTITIONS 2;

CREATE TABLE t_cmp_ws (
	id INT PRIMARY KEY
) COMPRESSION='custom value';

CREATE TABLE t_cmp_bt (
	id INT PRIMARY KEY
) COMPRESSION='` + "`custom`" + `';`
	execCorpus(t, d.db, setup)

	source, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if source.PlacementPolicy != "p_sync_ext" {
		t.Fatalf("database PlacementPolicy = %q, want p_sync_ext", source.PlacementPolicy)
	}
	tables := map[string]*metadatapb.TableMetadata{}
	for _, tbl := range source.Schemas[0].Tables {
		tables[tbl.Name] = tbl
	}
	if got := tables["t_chk_a"].CheckConstraints; len(got) != 1 || !strings.Contains(got[0].Expression, "a") || got[0].NotEnforced {
		t.Fatalf("t_chk_a checks = %+v", got)
	}
	if got := tables["t_chk_b"].CheckConstraints; len(got) != 1 || !strings.Contains(got[0].Expression, "b") || !got[0].NotEnforced {
		t.Fatalf("t_chk_b checks = %+v", got)
	}
	if got := tables["gtt_del"]; got == nil || !got.IsGlobalTemporary || got.OnCommit != "DELETE ROWS" {
		t.Fatalf("gtt_del = %+v, want IsGlobalTemporary + DELETE ROWS", got)
	}
	if got := tables["t_part_global"]; got == nil || len(got.Indexes) < 2 {
		t.Fatalf("t_part_global = %+v", got)
	} else {
		var ukGlobal bool
		for _, idx := range got.Indexes {
			if idx.Name == "uk_k" {
				ukGlobal = idx.IsGlobal
			}
		}
		if !ukGlobal {
			t.Fatalf("uk_k IsGlobal = false, want true: %+v", got.Indexes)
		}
	}
	if got := tables["t_cmp_ws"].CreateOptions; got != "compression='custom value'" {
		t.Fatalf("t_cmp_ws CreateOptions = %q", got)
	}
	if got := tables["t_cmp_bt"].CreateOptions; got != "compression='`custom`'" {
		t.Fatalf("t_cmp_bt CreateOptions = %q", got)
	}
	if got := tables["t_part_global"].PlacementPolicy; got != "p_sync_ext" {
		t.Fatalf("t_part_global PlacementPolicy = %q, want p_sync_ext", got)
	}

	var dumpBuf strings.Builder
	if err := d.Dump(ctx, &dumpBuf, source); err != nil {
		t.Fatal(err)
	}
	for _, stmt := range []string{
		"DROP DATABASE `omni_ext_roundtrip`",
		"CREATE DATABASE `omni_ext_roundtrip`",
		"USE `omni_ext_roundtrip`",
	} {
		if _, err := d.db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	for _, stmt := range splitLikeMySQLClient(dumpBuf.String()) {
		if _, err := d.db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("restoring statement:\n%s\nerr: %v", stmt, err)
		}
	}
	restored, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("restored extended TiDB database differs from source (-source +restored):\n%s", diff)
	}
}

// TestDumpRestoresAllocationOptions pins AUTO_ID_CACHE and PRE_SPLIT_REGIONS,
// which TiDB v8.5.5 states only in SHOW CREATE TABLE, through sync, dump and
// restore.
func TestDumpRestoresAllocationOptions(t *testing.T) {
	const setup = `CREATE TABLE t_cache (id BIGINT AUTO_INCREMENT PRIMARY KEY) AUTO_ID_CACHE = 1;

CREATE TABLE t_split (id INT) SHARD_ROW_ID_BITS = 4 PRE_SPLIT_REGIONS = 2;

CREATE TABLE p_quoted (` + "`a,b`" + ` INT, c INT) PARTITION BY KEY (` + "`a,b`" + `, c) PARTITIONS 2;`
	source, restored := dumpAndRestore(t, "omni_alloc_roundtrip", setup)
	tables := map[string]*metadatapb.TableMetadata{}
	for _, tbl := range source.Schemas[0].Tables {
		tables[tbl.Name] = tbl
	}
	if got := tables["t_cache"].AutoIdCache; got != 1 {
		t.Errorf("t_cache AUTO_ID_CACHE = %d, want 1", got)
	}
	if got := tables["t_split"].PreSplitRegions; got != 2 {
		t.Errorf("t_split PRE_SPLIT_REGIONS = %d, want 2", got)
	}
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("restored TiDB database differs from source (-source +restored):\n%s", diff)
	}
}
