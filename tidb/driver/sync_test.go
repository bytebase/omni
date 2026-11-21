package driver

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
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

type fakeSyncDriver struct{}

func (fakeSyncDriver) Open(string) (sqldriver.Conn, error) { return fakeSyncConn{}, nil }

type fakeSyncConn struct{}

func (fakeSyncConn) Prepare(string) (sqldriver.Stmt, error) { return nil, sqldriver.ErrSkip }
func (fakeSyncConn) Close() error                           { return nil }
func (fakeSyncConn) Begin() (sqldriver.Tx, error)           { return nil, sqldriver.ErrSkip }

func (fakeSyncConn) QueryContext(_ context.Context, query string, _ []sqldriver.NamedValue) (sqldriver.Rows, error) {
	switch {
	case strings.Contains(query, "SELECT VERSION()"):
		return &fakeRows{cols: []string{"VERSION()"}, rows: [][]sqldriver.Value{{"8.0.11-TiDB-v8.5.5"}}}, nil
	case strings.Contains(query, "FROM information_schema.TIDB_INDEXES"):
		return &fakeRows{
			cols: []string{"TABLE_NAME", "KEY_NAME", "COLUMN_NAME", "SEQ_IN_INDEX", "NON_UNIQUE", "INDEX_COMMENT", "SUB_PART", "EXPRESSION"},
			rows: [][]sqldriver.Value{
				{"t_vec", "PRIMARY", "id", int64(1), false, "", nil, nil},
				{"t_vec", "idx_v", nil, int64(1), true, "", nil, "VEC_COSINE_DISTANCE(v)"},
				{"t_vec", "uk_k", "k", int64(1), false, "", nil, nil},
			},
		}, nil
	case strings.Contains(query, "FROM information_schema.COLUMNS"):
		return &fakeRows{
			cols: []string{"TABLE_NAME", "COLUMN_NAME", "ORDINAL_POSITION", "COLUMN_DEFAULT", "IS_NULLABLE", "COLUMN_TYPE", "CHARACTER_SET_NAME", "COLLATION_NAME", "COLUMN_COMMENT", "GENERATION_EXPRESSION", "EXTRA"},
			rows: [][]sqldriver.Value{
				{"t_other", "id", int64(1), nil, "NO", "int", "", "", "''", nil, ""},
				{"t_other", "m", int64(2), nil, "NO", "int", "", "", "''", nil, ""},
				{"t_vec", "id", int64(1), nil, "NO", "int", "", "", "''", nil, ""},
				{"t_vec", "k", int64(2), nil, "NO", "int", "", "", "''", nil, ""},
				{"t_vec", "v", int64(3), nil, "YES", "vector(3)", "", "", "''", nil, ""},
			},
		}, nil
	case strings.Contains(query, "FROM information_schema.VIEWS"):
		return &fakeRows{cols: []string{"TABLE_NAME", "VIEW_DEFINITION", "SECURITY_TYPE", "CHECK_OPTION"}}, nil
	case strings.Contains(query, "FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS"):
		return &fakeRows{cols: []string{"TABLE_SCHEMA", "TABLE_NAME", "CONSTRAINT_NAME", "COLUMN_NAME", "REFERENCED_TABLE_SCHEMA", "REFERENCED_TABLE_NAME", "REFERENCED_COLUMN_NAME", "DELETE_RULE", "UPDATE_RULE", "MATCH_OPTION"}}, nil
	case strings.Contains(query, "FROM information_schema.TIDB_CHECK_CONSTRAINTS"):
		return &fakeRows{
			cols: []string{"TABLE_NAME", "CONSTRAINT_NAME", "CHECK_CLAUSE"},
			rows: [][]sqldriver.Value{
				{"t_other", "chk_shared", "(`m` < 100)"},
				{"t_vec", "chk_shared", "(`k` > 0)"},
			},
		}, nil
	case strings.Contains(query, "FROM information_schema.TIFLASH_REPLICA"):
		return &fakeRows{
			cols: []string{"TABLE_NAME", "REPLICA_COUNT", "LOCATION_LABELS"},
			rows: [][]sqldriver.Value{
				{"t_vec", int64(2), "zone, rack"},
			},
		}, nil
	case strings.Contains(query, "FROM information_schema.SEQUENCES"):
		return &fakeRows{cols: []string{"SEQUENCE_NAME", "START", "MIN_VALUE", "MAX_VALUE", "INCREMENT", "CYCLE", "CACHE_VALUE", "COMMENT"}}, nil
	case strings.Contains(query, "FROM INFORMATION_SCHEMA.PARTITIONS"):
		return &fakeRows{cols: []string{"TABLE_NAME", "PARTITION_NAME", "SUBPARTITION_NAME", "PARTITION_METHOD", "SUBPARTITION_METHOD", "PARTITION_EXPRESSION", "SUBPARTITION_EXPRESSION", "PARTITION_DESCRIPTION", "PARTITION_COMMENT"}}, nil
	case strings.Contains(query, "FROM information_schema.TABLES"):
		return &fakeRows{
			cols: []string{"TABLE_NAME", "TABLE_TYPE", "ENGINE", "TABLE_COLLATION", "TABLE_ROWS", "DATA_LENGTH", "INDEX_LENGTH", "DATA_FREE", "CREATE_OPTIONS", "TABLE_COMMENT", "SHARDING", "PK_TYPE", "POLICY"},
			rows: [][]sqldriver.Value{
				{"t_other", "BASE TABLE", "InnoDB", "utf8mb4_bin", int64(0), int64(0), int64(0), int64(0), "", "''", "", "CLUSTERED", ""},
				{"t_vec", "BASE TABLE", "InnoDB", "utf8mb4_bin", int64(0), int64(0), int64(0), int64(0), "", "''", "", "CLUSTERED", "p_hot"},
			},
		}, nil
	case strings.Contains(query, "SHOW CREATE TABLE"):
		if strings.Contains(query, "`t_other`") {
			ddl := "CREATE TABLE `t_other` (\n" +
				"  `id` int NOT NULL,\n" +
				"  `m` int NOT NULL,\n" +
				"  CONSTRAINT `chk_shared` CHECK ((`m` < 100))\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin"
			return &fakeRows{cols: []string{"Table", "Create Table"}, rows: [][]sqldriver.Value{{"t_other", ddl}}}, nil
		}
		ddl := "CREATE TABLE `t_vec` (\n" +
			"  `id` int NOT NULL,\n" +
			"  `k` int NOT NULL,\n" +
			"  `v` vector(3) DEFAULT NULL,\n" +
			"  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,\n" +
			"  UNIQUE KEY `uk_k` (`k`) /*T![global_index] GLOBAL */,\n" +
			"  VECTOR INDEX `idx_v`((VEC_COSINE_DISTANCE(`v`))) USING HNSW,\n" +
			"  CONSTRAINT `chk_shared` CHECK ((`k` > 0))\n" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMPRESSION='`custom`' /*T![placement] PLACEMENT POLICY=`p_hot` */"
		return &fakeRows{cols: []string{"Table", "Create Table"}, rows: [][]sqldriver.Value{{"t_vec", ddl}}}, nil
	case strings.Contains(query, "FROM information_schema.SCHEMATA"):
		return &fakeRows{
			cols: []string{"DEFAULT_CHARACTER_SET_NAME", "DEFAULT_COLLATION_NAME", "TIDB_PLACEMENT_POLICY_NAME"},
			rows: [][]sqldriver.Value{{"utf8mb4", "utf8mb4_bin", "p_db"}},
		}, nil
	default:
		return nil, fmt.Errorf("unexpected query: %s", query)
	}
}

type fakeRows struct {
	cols []string
	rows [][]sqldriver.Value
	idx  int
}

func (r *fakeRows) Columns() []string { return r.cols }
func (r *fakeRows) Close() error      { return nil }
func (r *fakeRows) Next(dest []sqldriver.Value) error {
	if r.idx >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.idx])
	r.idx++
	return nil
}

func init() {
	sql.Register("test_tidb_fake_sync", fakeSyncDriver{})
}

func TestSyncTiFlashVectorGlobalAndPlacementPolicy(t *testing.T) {
	db, err := sql.Open("test_tidb_fake_sync", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	d := &Driver{db: db, databaseName: "app"}
	meta, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if meta.PlacementPolicy != "p_db" {
		t.Fatalf("db PlacementPolicy = %q, want p_db", meta.PlacementPolicy)
	}
	if len(meta.Schemas) != 1 || len(meta.Schemas[0].Tables) != 2 {
		t.Fatalf("unexpected schemas/tables: %+v", meta.Schemas)
	}
	byTbl := map[string]*metadatapb.TableMetadata{}
	for _, tbl := range meta.Schemas[0].Tables {
		byTbl[tbl.Name] = tbl
	}
	if chk := byTbl["t_other"].CheckConstraints; len(chk) != 1 || chk[0].Expression != "(`m` < 100)" {
		t.Fatalf("t_other checks = %+v", chk)
	}
	tbl := byTbl["t_vec"]
	if chk := tbl.CheckConstraints; len(chk) != 1 || chk[0].Expression != "(`k` > 0)" {
		t.Fatalf("t_vec checks = %+v", chk)
	}
	if tbl.TiflashReplica != 2 || len(tbl.TiflashLocationLabels) != 2 || tbl.TiflashLocationLabels[0] != "zone" || tbl.TiflashLocationLabels[1] != "rack" {
		t.Fatalf("TiFlash = (%d, %v), want (2, [zone rack])", tbl.TiflashReplica, tbl.TiflashLocationLabels)
	}
	if tbl.PlacementPolicy != "p_hot" {
		t.Fatalf("table PlacementPolicy = %q, want p_hot", tbl.PlacementPolicy)
	}
	if tbl.CreateOptions != "compression='`custom`'" {
		t.Fatalf("CreateOptions = %q", tbl.CreateOptions)
	}
	idxByName := map[string]*metadatapb.IndexMetadata{}
	for _, idx := range tbl.Indexes {
		idxByName[idx.Name] = idx
	}
	if idx := idxByName["idx_v"]; idx == nil || !idx.IsVector || idx.Type != "HNSW" {
		t.Fatalf("vector index = %+v, want IsVector+HNSW", idx)
	}
	if idx := idxByName["uk_k"]; idx == nil || !idx.IsGlobal {
		t.Fatalf("global index = %+v, want IsGlobal=true", idx)
	}

	var dumpBuf strings.Builder
	if err := d.Dump(context.Background(), &dumpBuf, meta); err != nil {
		t.Fatal(err)
	}
	dumped := dumpBuf.String()
	for _, want := range []string{
		"ALTER DATABASE PLACEMENT POLICY=`p_db`;",
		"VECTOR INDEX `idx_v`",
		"/*T![global_index] GLOBAL */",
		"COMPRESSION='`custom`'",
		"ALTER TABLE `t_vec` SET TIFLASH REPLICA 2 LOCATION LABELS 'zone', 'rack';",
	} {
		if !strings.Contains(dumped, want) {
			t.Errorf("Dump missing %q:\n%s", want, dumped)
		}
	}
}
