package catalog

import (
	"context"
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

func TestDatabaseMetadataEmptyAndMissing(t *testing.T) {
	var nilCatalog *Catalog
	meta := nilCatalog.DatabaseMetadata("app")
	require.Equal(t, "app", meta.Name)
	require.Empty(t, meta.Schemas)

	c := New()
	meta = c.DatabaseMetadata("missing")
	require.Equal(t, "missing", meta.Name)
	require.Empty(t, meta.Schemas)
}

func TestDatabaseMetadataFull(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		SET tidb_enable_check_constraint = 1;
		CREATE DATABASE app CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
		USE app;

		CREATE TABLE items (
			id BIGINT AUTO_RANDOM(5, 54) PRIMARY KEY CLUSTERED,
			code VARCHAR(20) NOT NULL DEFAULT 'abc',
			num INT DEFAULT 5,
			bit_col BIT(3) DEFAULT b'101',
			hex_str VARCHAR(10) DEFAULT 0x41,
			hex_int INT DEFAULT 0x41,
			nullable_col INT,
			ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			gen_col VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_g')) STORED,
			CHECK (nullable_col >= 0)
		) ENGINE=InnoDB COMMENT='items table';
		ALTER TABLE items CACHE;

		CREATE TABLE items_bare_autorand (
			id BIGINT AUTO_RANDOM PRIMARY KEY NONCLUSTERED
		);

		CREATE TABLE items_sharded (
			id BIGINT PRIMARY KEY
		) SHARD_ROW_ID_BITS=4;

		CREATE TABLE child (
			id BIGINT PRIMARY KEY,
			item_id BIGINT,
			CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES items(id)
		);

		CREATE TABLE part_tbl (
			id INT NOT NULL,
			yr INT NOT NULL,
			PRIMARY KEY (id, yr)
		) PARTITION BY RANGE (yr) (
			PARTITION p0 VALUES LESS THAN (2025),
			PARTITION p1 VALUES LESS THAN MAXVALUE
		);

		CREATE VIEW v_items (vid, vcode) AS SELECT id, code FROM items;
		CREATE SEQUENCE seq_order START WITH 10 INCREMENT BY 2 MINVALUE 1 MAXVALUE 1000 CACHE 50 CYCLE;
	`, nil)
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	require.Equal(t, "app", meta.Name)
	require.Equal(t, "utf8mb4", meta.CharacterSet)
	require.Equal(t, "utf8mb4_bin", meta.Collation)
	require.Len(t, meta.Schemas, 1)

	s := meta.Schemas[0]
	require.Len(t, s.Tables, 5)

	items := s.Tables[1]
	require.Equal(t, "items", items.Name)
	require.Equal(t, "utf8mb4", items.Charset)
	require.Equal(t, "utf8mb4_bin", items.Collation)
	require.Equal(t, "items table", items.Comment)
	require.Equal(t, "CLUSTERED", items.PrimaryKeyType)
	require.Equal(t, "PK_AUTO_RANDOM_BITS=5, RANGE BITS=54", items.ShardingInfo)
	require.Equal(t, "cached=on", items.CreateOptions)

	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range items.Columns {
		cols[col.Name] = col
	}
	require.Equal(t, "AUTO_RANDOM(5, 54)", cols["id"].Default)
	require.Equal(t, "'abc'", cols["code"].Default)
	require.Equal(t, "'5'", cols["num"].Default)
	require.Equal(t, "b'101'", cols["bit_col"].Default)
	require.Equal(t, "'A'", cols["hex_str"].Default)
	require.Equal(t, "'65'", cols["hex_int"].Default)
	require.Empty(t, cols["nullable_col"].Default, "TiDB leaves nullable column without default empty")
	require.Equal(t, "CURRENT_TIMESTAMP", cols["ts"].Default)
	require.Equal(t, "CURRENT_TIMESTAMP", cols["ts"].OnUpdate)
	require.Equal(t, metadata.GenerationMetadata_TYPE_STORED, cols["gen_col"].Generation.Type)
	require.Len(t, items.CheckConstraints, 1)

	require.Equal(t, "AUTO_RANDOM", s.Tables[2].Columns[0].Default)
	require.Equal(t, "NONCLUSTERED", s.Tables[2].PrimaryKeyType)
	require.Equal(t, "SHARD_BITS=4", s.Tables[3].ShardingInfo)

	partTbl := s.Tables[4]
	require.Len(t, partTbl.Partitions, 2)
	require.Equal(t, metadata.TablePartitionMetadata_RANGE, partTbl.Partitions[0].Type)

	require.Len(t, s.Views, 1)
	require.Equal(t, "v_items", s.Views[0].Name)
	require.Len(t, s.Views[0].Columns, 2)
	require.Equal(t, "vid", s.Views[0].Columns[0].Name)

	require.Len(t, s.Sequences, 1)
	require.Equal(t, "seq_order", s.Sequences[0].Name)
	require.Equal(t, "BIGINT", s.Sequences[0].DataType)
	require.Equal(t, "10", s.Sequences[0].Start)
	require.Equal(t, "2", s.Sequences[0].Increment)
	require.Equal(t, "50", s.Sequences[0].CacheSize)
	require.True(t, s.Sequences[0].Cycle)

	// Round-trip through LoadMetadata preserves objects cleanly.
	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)
	require.True(t, reloaded.GetDatabase("app").GetTable("items").Cached)
	require.Equal(t, 4, reloaded.GetDatabase("app").GetTable("items_sharded").ShardRowIDBits)
	require.Equal(t, []string{"vid", "vcode"}, reloaded.GetDatabase("app").Views["v_items"].Columns)

	// Mutating exported FK slices must not mutate catalog state.
	child := s.Tables[0]
	require.Len(t, child.ForeignKeys, 1)
	child.ForeignKeys[0].Columns[0] = "mutated"
	var fkCon *Constraint
	for _, con := range c.GetDatabase("app").GetTable("child").Constraints {
		if con.Type == ConForeignKey {
			fkCon = con
		}
	}
	require.NotNil(t, fkCon)
	require.Equal(t, "item_id", fkCon.Columns[0])
}

func TestDatabaseMetadataExcludesTempTablesAndSortsIndexes(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TEMPORARY TABLE tmp_scratch (id INT PRIMARY KEY);
		CREATE TABLE items (
			id INT NOT NULL,
			b INT,
			a INT,
			PRIMARY KEY (id),
			KEY idx_zeta (b),
			KEY idx_alpha (a)
		);
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	require.Len(t, meta.Schemas[0].Tables, 1)
	tbl := meta.Schemas[0].Tables[0]
	require.Equal(t, "items", tbl.Name)
	require.Equal(t, []string{"PRIMARY", "idx_alpha", "idx_zeta"}, []string{
		tbl.Indexes[0].Name,
		tbl.Indexes[1].Name,
		tbl.Indexes[2].Name,
	})
}

func TestDatabaseMetadataRoundTripVectorGlobalChecksDefaultsAndSubpartitions(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		SET tidb_enable_check_constraint = 1;
		CREATE DATABASE app;
		USE app;
		CREATE TABLE ref_tbl (id INT PRIMARY KEY);
		CREATE TABLE full_tbl (
			id INT NOT NULL,
			ref_z INT,
			ref_a INT,
			drop_col INT NULL DEFAULT 7,
			keep_null INT NULL,
			vec VECTOR(3),
			u INT,
			PRIMARY KEY (id),
			UNIQUE KEY uq_u (u) GLOBAL,
			VECTOR INDEX idx_vec ((VEC_COSINE_DISTANCE(vec))),
			CONSTRAINT fk_zeta FOREIGN KEY (ref_z) REFERENCES ref_tbl (id),
			CONSTRAINT fk_alpha FOREIGN KEY (ref_a) REFERENCES ref_tbl (id),
			CONSTRAINT chk_zeta CHECK (id > 0) NOT ENFORCED,
			CONSTRAINT chk_alpha CHECK (id < 1000)
		) COMPRESSION='ZLIB';
		ALTER TABLE full_tbl ALTER COLUMN drop_col DROP DEFAULT;
		CREATE GLOBAL TEMPORARY TABLE gtt_tbl (
			id INT NOT NULL PRIMARY KEY
		) ON COMMIT DELETE ROWS;
		CREATE TABLE part_tbl (
			id INT NOT NULL
		) PARTITION BY HASH (id) PARTITIONS 2;
		CREATE ALGORITHM=TEMPTABLE SQL SECURITY INVOKER VIEW v_exec AS SELECT 1 AS x;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	// Populate SubType/Subpartitions on part_tbl to verify subpartition serialization.
	partTbl := c.GetDatabase("app").GetTable("part_tbl")
	partTbl.Partitioning.SubType = "KEY"
	partTbl.Partitioning.SubColumns = []string{"id"}
	partTbl.Partitioning.SubAlgo = 1
	partTbl.Partitioning.NumSubParts = 2
	for _, p := range partTbl.Partitioning.Partitions {
		p.SubPartitions = []*SubPartitionDefInfo{{Name: "s0"}, {Name: "s1"}}
	}

	meta := c.DatabaseMetadata("app")
	s := meta.Schemas[0]
	full := s.Tables[0]
	require.Equal(t, "full_tbl", full.Name)
	require.Equal(t, []string{"fk_alpha", "fk_zeta"}, []string{full.ForeignKeys[0].Name, full.ForeignKeys[1].Name})
	require.Equal(t, []string{"chk_alpha", "chk_zeta"}, []string{full.CheckConstraints[0].Name, full.CheckConstraints[1].Name})
	require.False(t, full.CheckConstraints[0].NotEnforced)
	require.True(t, full.CheckConstraints[1].NotEnforced)

	idxByName := map[string]*metadata.IndexMetadata{}
	for _, idx := range full.Indexes {
		idxByName[idx.Name] = idx
	}
	require.True(t, idxByName["uq_u"].IsGlobal)
	require.True(t, idxByName["idx_vec"].IsVector)

	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range full.Columns {
		cols[col.Name] = col
	}
	require.True(t, cols["drop_col"].DefaultDropped)
	require.False(t, cols["keep_null"].DefaultDropped)

	partMeta := s.Tables[2]
	require.Equal(t, "part_tbl", partMeta.Name)
	require.Equal(t, "2", partMeta.Partitions[0].UseDefault)
	require.Len(t, partMeta.Partitions[0].Subpartitions, 2)
	require.Equal(t, "2", partMeta.Partitions[0].Subpartitions[0].UseDefault)
	require.Equal(t, int32(1), partMeta.Partitions[0].Subpartitions[0].Algorithm)

	gttMeta := s.Tables[1]
	require.Equal(t, "gtt_tbl", gttMeta.Name)
	require.True(t, gttMeta.IsGlobalTemporary)
	require.Equal(t, "DELETE ROWS", gttMeta.OnCommit)
	require.Contains(t, full.CreateOptions, "compression='ZLIB'")
	require.Len(t, s.Views, 1)
	require.Equal(t, "TEMPTABLE", s.Views[0].Algorithm)
	require.Equal(t, "INVOKER", s.Views[0].SqlSecurity)

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	reloadedFull := reloaded.GetDatabase("app").GetTable("full_tbl")
	require.Equal(t, "ZLIB", reloadedFull.Compression)
	require.True(t, reloadedFull.GetColumn("drop_col").DefaultDropped)
	require.False(t, reloadedFull.GetColumn("keep_null").DefaultDropped)
	reloadedGTT := reloaded.GetDatabase("app").GetTable("gtt_tbl")
	require.True(t, reloadedGTT.GlobalTemporary)
	require.Equal(t, "DELETE ROWS", reloadedGTT.OnCommit)
	reView := reloaded.GetDatabase("app").Views["v_exec"]
	require.Equal(t, "TEMPTABLE", reView.Algorithm)
	require.Equal(t, "INVOKER", reView.SqlSecurity)
	reExportedFull := reloaded.DatabaseMetadata("app").Schemas[0].Tables[0]
	reIdxByName := map[string]*metadata.IndexMetadata{}
	for _, idx := range reExportedFull.Indexes {
		reIdxByName[idx.Name] = idx
	}
	require.True(t, reIdxByName["uq_u"].IsGlobal)
	require.True(t, reIdxByName["idx_vec"].IsVector)
	require.True(t, reExportedFull.CheckConstraints[1].NotEnforced)
}

func TestDatabaseMetadataTTLPlacementIntervalAndViewCheckOption(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE PLACEMENT POLICY p_hot PRIMARY_REGION="us-east" REGIONS="us-east";
		CREATE DATABASE app;
		USE app;

		CREATE TABLE ttl_tbl (
			id BIGINT PRIMARY KEY,
			created_at DATETIME NOT NULL
		) PLACEMENT POLICY = p_hot
		  AUTO_ID_CACHE = 200
		  PRE_SPLIT_REGIONS = 2
		  TTL = created_at + INTERVAL 30 DAY
		  TTL_ENABLE = 'OFF'
		  TTL_JOB_INTERVAL = '2h';

		ALTER TABLE ttl_tbl SET TIFLASH REPLICA 2 LOCATION LABELS 'zone', 'host';

		CREATE TABLE iv_tbl (
			a INT NOT NULL
		) PARTITION BY RANGE (a) INTERVAL (10) FIRST PARTITION LESS THAN (ABS(10)) LAST PARTITION LESS THAN (100);

		CREATE VIEW v_check AS SELECT id FROM ttl_tbl WITH CASCADED CHECK OPTION;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	s := meta.Schemas[0]
	// Tables sorted: iv_tbl, ttl_tbl
	ivMeta := s.Tables[0]
	require.Equal(t, "iv_tbl", ivMeta.Name)
	require.Len(t, ivMeta.Partitions, 1)
	require.Contains(t, ivMeta.Partitions[0].IntervalClause, "INTERVAL")
	require.Empty(t, ivMeta.Partitions[0].Name)

	ttlMeta := s.Tables[1]
	require.Equal(t, "ttl_tbl", ttlMeta.Name)
	require.Equal(t, "p_hot", ttlMeta.PlacementPolicy)
	require.Equal(t, "created_at", ttlMeta.TtlColumn)
	require.Equal(t, "interval 30 day", ttlMeta.TtlInterval)
	require.False(t, ttlMeta.TtlEnable)
	require.Equal(t, "2h", ttlMeta.TtlJobInterval)
	require.Equal(t, int64(200), ttlMeta.AutoIdCache)
	require.Equal(t, int32(2), ttlMeta.PreSplitRegions)
	require.Equal(t, int32(2), ttlMeta.TiflashReplica)
	require.Equal(t, []string{"zone", "host"}, ttlMeta.TiflashLocationLabels)

	require.Equal(t, "CASCADED", s.Views[0].CheckOption)

	// Reload into a fresh catalog with NO placement policies loaded separately.
	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	reTTL := reloaded.GetDatabase("app").GetTable("ttl_tbl")
	require.Equal(t, "p_hot", reTTL.PlacementPolicy)
	require.Equal(t, "created_at", reTTL.TTLColumn)
	require.Equal(t, "interval 30 day", reTTL.TTLInterval)
	require.False(t, reTTL.TTLEnable)
	require.Equal(t, "2h", reTTL.TTLJobInterval)
	require.Equal(t, int64(200), reTTL.AutoIDCache)
	require.Equal(t, 2, reTTL.PreSplitRegions)
	require.Equal(t, 2, reTTL.TiFlashReplica)
	require.Equal(t, []string{"zone", "host"}, reTTL.TiFlashLocationLabels)

	reIV := reloaded.GetDatabase("app").GetTable("iv_tbl")
	require.NotNil(t, reIV.Partitioning)
	require.Contains(t, reIV.Partitioning.IntervalText, "INTERVAL")
	require.Empty(t, reIV.Partitioning.Partitions)

	require.Equal(t, "CASCADED", reloaded.GetDatabase("app").Views["v_check"].CheckOption)
}

func TestDatabaseMetadataExplicitDefaultNullPreserved(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE t_null (
			id INT PRIMARY KEY,
			s VARCHAR(20) DEFAULT NULL,
			n INT DEFAULT NULL,
			bare INT
		);
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range meta.Schemas[0].Tables[0].Columns {
		cols[col.Name] = col
	}
	require.Equal(t, "NULL", cols["s"].Default)
	require.Equal(t, "NULL", cols["n"].Default)
	require.Empty(t, cols["bare"].Default)

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	reTbl := reloaded.GetDatabase("app").GetTable("t_null")
	require.NotNil(t, reTbl.GetColumn("s").Default)
	require.Equal(t, "NULL", *reTbl.GetColumn("s").Default)
	require.NotNil(t, reTbl.GetColumn("n").Default)
	require.Equal(t, "NULL", *reTbl.GetColumn("n").Default)
	require.Nil(t, reTbl.GetColumn("bare").Default)

	reExport := reloaded.DatabaseMetadata("app")
	reCols := map[string]*metadata.ColumnMetadata{}
	for _, col := range reExport.Schemas[0].Tables[0].Columns {
		reCols[col.Name] = col
	}
	require.Equal(t, "NULL", reCols["s"].Default)
	require.Equal(t, "NULL", reCols["n"].Default)
	require.Empty(t, reCols["bare"].Default)

	// Per-index KEY_BLOCK_SIZE round-trips cleanly.
	kbsCatalog := New()
	_, err = kbsCatalog.Exec(`CREATE DATABASE app; USE app; CREATE TABLE idx_kbs (id INT PRIMARY KEY, a INT, KEY ka (a) KEY_BLOCK_SIZE = 4);`, nil)
	require.NoError(t, err)
	kbsMeta := kbsCatalog.DatabaseMetadata("app")
	require.Equal(t, int64(4), kbsMeta.Schemas[0].Tables[0].Indexes[1].KeyBlockSize)
	kbsReloaded := New()
	kbsReport, err := kbsReloaded.LoadMetadata(context.Background(), kbsMeta)
	require.NoError(t, err)
	require.Empty(t, kbsReport.Incomplete)
	require.Equal(t, 4, kbsReloaded.GetDatabase("app").GetTable("idx_kbs").Indexes[1].KeyBlockSize)

	// Database-level PLACEMENT POLICY and TIFLASH REPLICA round-trip cleanly.
	dbCatalog := New()
	_, err = dbCatalog.Exec(`
		CREATE PLACEMENT POLICY p_db PRIMARY_REGION="us-east" REGIONS="us-east";
		CREATE DATABASE db_placed PLACEMENT POLICY = p_db;
		ALTER DATABASE db_placed SET TIFLASH REPLICA 3 LOCATION LABELS 'zone', 'host';
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	dbMeta := dbCatalog.DatabaseMetadata("db_placed")
	require.Equal(t, "p_db", dbMeta.PlacementPolicy)
	require.Equal(t, int32(3), dbMeta.TiflashReplica)
	require.Equal(t, []string{"zone", "host"}, dbMeta.TiflashLocationLabels)

	dbReloaded := New()
	dbReport, err := dbReloaded.LoadMetadata(context.Background(), dbMeta)
	require.NoError(t, err)
	require.Empty(t, dbReport.Degraded)
	require.Empty(t, dbReport.Missing)
	require.Empty(t, dbReport.Incomplete)
	reDB := dbReloaded.GetDatabase("db_placed")
	require.Equal(t, "p_db", reDB.PlacementPolicy)
	require.Equal(t, 3, reDB.TiFlashReplica)
	require.Equal(t, []string{"zone", "host"}, reDB.TiFlashLocationLabels)

	// DatabaseDefinition renders AUTO_ID_CACHE and PRE_SPLIT_REGIONS.
	allocCatalog := New()
	_, err = allocCatalog.Exec(`
		CREATE DATABASE alloc_db;
		USE alloc_db;
		CREATE TABLE alloc_tbl (id BIGINT PRIMARY KEY) AUTO_ID_CACHE = 200 PRE_SPLIT_REGIONS = 2;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	rendered, err := DatabaseDefinition(allocCatalog.DatabaseMetadata("alloc_db"), DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, rendered, "AUTO_ID_CACHE=200")
	require.Contains(t, rendered, "PRE_SPLIT_REGIONS=2")
}
