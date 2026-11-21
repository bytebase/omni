package catalog

import (
	"context"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/tidb/ast"
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

	// DatabaseDefinition renders AUTO_ID_CACHE, PRE_SPLIT_REGIONS, and TTL options,
	// and DatabaseMetadata exports functions, procedures, triggers, and events.
	allocCatalog := New()
	_, err = allocCatalog.Exec(`
		CREATE DATABASE alloc_db;
		USE alloc_db;
		CREATE TABLE alloc_tbl (
			id BIGINT PRIMARY KEY,
			created_at DATETIME NOT NULL
		) AUTO_ID_CACHE = 200 PRE_SPLIT_REGIONS = 2 TTL = created_at + INTERVAL 7 DAY TTL_ENABLE = 'OFF' TTL_JOB_INTERVAL = '2h';
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	require.NoError(t, allocCatalog.DefineFunction(&nodes.CreateFunctionStmt{
		Name:     &nodes.TableRef{Name: "fn_inc"},
		Returns:  &nodes.DataType{Name: "INT"},
		BodyText: "RETURN 1",
	}))
	require.NoError(t, allocCatalog.DefineProcedure(&nodes.CreateFunctionStmt{
		IsProcedure: true,
		Name:        &nodes.TableRef{Name: "proc_noop"},
		BodyText:    "BEGIN SELECT 1; END",
	}))
	require.NoError(t, allocCatalog.DefineTrigger(&nodes.CreateTriggerStmt{
		Name:     "trg_alloc",
		Table:    &nodes.TableRef{Name: "alloc_tbl"},
		Timing:   "BEFORE",
		Event:    "INSERT",
		BodyText: "SET NEW.id = NEW.id",
	}))
	require.NoError(t, allocCatalog.DefineEvent(&nodes.CreateEventStmt{
		Name:     "ev_alloc",
		Schedule: &nodes.EventSchedule{RawText: "EVERY 1 HOUR"},
		Comment:  "ev cmt",
		BodyText: "SELECT 1",
	}))
	allocMeta := allocCatalog.DatabaseMetadata("alloc_db")
	rendered, err := DatabaseDefinition(allocMeta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, rendered, "AUTO_ID_CACHE=200")
	require.Contains(t, rendered, "PRE_SPLIT_REGIONS=2")
	require.Contains(t, rendered, "TTL=`created_at` + INTERVAL 7 DAY")
	require.Contains(t, rendered, "TTL_ENABLE='OFF'")
	require.Contains(t, rendered, "TTL_JOB_INTERVAL='2h'")
	require.NotContains(t, rendered, "= ;")
	require.NotContains(t, rendered, "SET character_set_client")

	allocSchema := allocMeta.Schemas[0]
	require.Len(t, allocSchema.Functions, 1)
	require.Equal(t, "fn_inc", allocSchema.Functions[0].Name)
	require.Contains(t, allocSchema.Functions[0].Definition, "FUNCTION `fn_inc`")
	require.Len(t, allocSchema.Procedures, 1)
	require.Equal(t, "proc_noop", allocSchema.Procedures[0].Name)
	require.Contains(t, allocSchema.Procedures[0].Definition, "PROCEDURE `proc_noop`")
	require.Len(t, allocSchema.Events, 1)
	require.Equal(t, "ev_alloc", allocSchema.Events[0].Name)
	require.Equal(t, "ev cmt", allocSchema.Events[0].Comment)
	require.Len(t, allocSchema.Tables[0].Triggers, 1)
	require.Equal(t, "trg_alloc", allocSchema.Tables[0].Triggers[0].Name)

	reloadedAlloc := New()
	allocReport, err := reloadedAlloc.LoadMetadata(context.Background(), allocMeta)
	require.NoError(t, err)
	require.Empty(t, allocReport.Degraded)
	require.Empty(t, allocReport.Missing)
	require.Empty(t, allocReport.Incomplete)
	reAllocDB := reloadedAlloc.GetDatabase("alloc_db")
	require.Contains(t, reAllocDB.Functions, "fn_inc")
	require.Contains(t, reAllocDB.Procedures, "proc_noop")
	require.Contains(t, reAllocDB.Triggers, "trg_alloc")
	require.Contains(t, reAllocDB.Events, "ev_alloc")

	// DatabaseDefinition renders view execution attributes, placement policy, TiFlash, cache, compression, key block size, escaped comments, and global temp tables.
	defCatalog := New()
	_, err = defCatalog.Exec(`
		CREATE PLACEMENT POLICY p_tbl PRIMARY_REGION="us-east" REGIONS="us-east";
		CREATE DATABASE ref_db;
		USE ref_db;
		CREATE TABLE parent_t (id INT PRIMARY KEY);
		CREATE DATABASE def_db PLACEMENT POLICY = p_tbl;
		ALTER DATABASE def_db SET TIFLASH REPLICA 1 LOCATION LABELS 'zone';
		USE def_db;
		CREATE TABLE placed_tbl (
			id INT PRIMARY KEY,
			k INT,
			pid INT,
			d_col INT DEFAULT 1,
			note VARCHAR(50) COMMENT 'it''s a \\ path',
			KEY k_blk (k) KEY_BLOCK_SIZE = 4,
			CONSTRAINT fk_ext FOREIGN KEY (pid) REFERENCES ref_db.parent_t (id)
		) COMPRESSION='zlib' COMMENT='O''Reilly' PLACEMENT POLICY = p_tbl;
		ALTER TABLE placed_tbl ALTER COLUMN d_col DROP DEFAULT;
		ALTER TABLE placed_tbl SET TIFLASH REPLICA 2 LOCATION LABELS 'zone', 'host';
		ALTER TABLE placed_tbl CACHE;
		CREATE GLOBAL TEMPORARY TABLE gtt_tbl (id INT PRIMARY KEY) ON COMMIT DELETE ROWS;
		CREATE ALGORITHM=MERGE SQL SECURITY INVOKER VIEW v_attr AS SELECT id, note FROM placed_tbl WITH CASCADED CHECK OPTION;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	defMeta := defCatalog.DatabaseMetadata("def_db")
	vAttrCols := defMeta.Schemas[0].Views[0].Columns
	require.Len(t, vAttrCols, 2)
	require.True(t, vAttrCols[0].Nullable)
	require.Equal(t, "int", vAttrCols[0].Type)
	require.Empty(t, vAttrCols[0].CharacterSet)
	require.Empty(t, vAttrCols[0].Collation)
	require.True(t, vAttrCols[1].Nullable)
	require.Equal(t, "varchar(50)", vAttrCols[1].Type)
	require.Empty(t, vAttrCols[1].CharacterSet)
	require.Empty(t, vAttrCols[1].Collation)

	defRendered, err := DatabaseDefinition(defMeta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, defRendered, "ALTER DATABASE `def_db` PLACEMENT POLICY=`p_tbl`;")
	require.Contains(t, defRendered, "ALTER DATABASE `def_db` SET TIFLASH REPLICA 1 LOCATION LABELS 'zone';")
	require.Contains(t, defRendered, "CONSTRAINT `fk_ext` FOREIGN KEY (`pid`) REFERENCES `ref_db`.`parent_t` (`id`)")
	require.Contains(t, defRendered, "/*T![placement] PLACEMENT POLICY=`p_tbl` */")
	require.Contains(t, defRendered, "COMPRESSION='zlib'")
	require.Contains(t, defRendered, `COMMENT='O\'Reilly'`)
	require.Contains(t, defRendered, `COMMENT 'it\'s a \\ path'`)
	require.Contains(t, defRendered, "KEY `k_blk` (`k`) KEY_BLOCK_SIZE=4")
	require.Contains(t, defRendered, "ALTER TABLE `placed_tbl` ALTER COLUMN `d_col` DROP DEFAULT;")
	require.Contains(t, defRendered, "ALTER TABLE `placed_tbl` SET TIFLASH REPLICA 2 LOCATION LABELS 'zone', 'host';")
	require.Contains(t, defRendered, "ALTER TABLE `placed_tbl` CACHE;")
	require.Less(t, strings.Index(defRendered, "SET TIFLASH REPLICA 2"), strings.Index(defRendered, "ALTER TABLE `placed_tbl` CACHE"))
	require.Contains(t, defRendered, "CREATE GLOBAL TEMPORARY TABLE `gtt_tbl`")
	require.Contains(t, defRendered, "ON COMMIT DELETE ROWS")
	require.Contains(t, defRendered, "CREATE ALGORITHM=MERGE SQL SECURITY INVOKER VIEW `v_attr`")
	require.Contains(t, defRendered, "WITH CASCADED CHECK OPTION")

	restored := New()
	_, err = restored.Exec("CREATE PLACEMENT POLICY p_tbl PRIMARY_REGION=\"us-east\" REGIONS=\"us-east\";\nCREATE DATABASE ref_db;\nUSE ref_db;\nCREATE TABLE parent_t (id INT PRIMARY KEY);\nCREATE DATABASE def_db;\nUSE def_db;\n"+defRendered, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	reDefDB := restored.GetDatabase("def_db")
	require.Equal(t, "p_tbl", reDefDB.PlacementPolicy)
	require.Equal(t, 1, reDefDB.TiFlashReplica)
	require.Equal(t, []string{"zone"}, reDefDB.TiFlashLocationLabels)
	rePlaced := reDefDB.GetTable("placed_tbl")
	require.Equal(t, "O'Reilly", rePlaced.Comment)
	require.Equal(t, `it's a \ path`, rePlaced.GetColumn("note").Comment)
	require.Equal(t, "zlib", rePlaced.Compression)
	require.Equal(t, 2, rePlaced.TiFlashReplica)
	require.Equal(t, []string{"zone", "host"}, rePlaced.TiFlashLocationLabels)
	require.True(t, rePlaced.Cached)
	require.True(t, rePlaced.GetColumn("d_col").DefaultDropped)

	// LoadMetadataInto clears database-level placement policy & TiFlash settings when unset in target snapshot.
	_, err = restored.LoadMetadataInto(context.Background(), &metadata.DatabaseSchemaMetadata{
		Name:    "other_db",
		Schemas: []*metadata.SchemaMetadata{{}},
	}, "def_db")
	require.NoError(t, err)
	require.Empty(t, reDefDB.PlacementPolicy)
	require.Zero(t, reDefDB.TiFlashReplica)
	require.Empty(t, reDefDB.TiFlashLocationLabels)

	// Parent partition comments are preserved alongside explicit subpartitions.
	subDef, err := DatabaseDefinition(&metadata.DatabaseSchemaMetadata{
		Name: "def_db",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:    "sub_t",
				Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: false}, {Name: "k", Type: "int", Nullable: false}},
				Partitions: []*metadata.TablePartitionMetadata{{
					Type:       metadata.TablePartitionMetadata_RANGE,
					Expression: "id",
					Name:       "p0",
					Value:      "10",
					Comment:    "p0_parent",
					Subpartitions: []*metadata.TablePartitionMetadata{{
						Type:       metadata.TablePartitionMetadata_HASH,
						Expression: "k",
						Name:       "s0",
						Comment:    "s0_child",
					}},
				}},
			}},
		}},
	}, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, subDef, "PARTITION p0 VALUES LESS THAN (10) ENGINE=InnoDB COMMENT 'p0_parent'")
	require.Contains(t, subDef, "SUBPARTITION s0 ENGINE=InnoDB COMMENT 's0_child'")
}

func TestDatabaseMetadataTriggerActionOrderAndBacktickQuoting(t *testing.T) {
	c := New()
	_, err := c.Exec("CREATE DATABASE trg_db; USE trg_db;\nCREATE TABLE `t``1` (a INT);", &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	require.NoError(t, c.DefineTrigger(&nodes.CreateTriggerStmt{
		Name:     "a_base",
		Table:    &nodes.TableRef{Name: "t`1"},
		Timing:   "BEFORE",
		Event:    "INSERT",
		BodyText: "SET NEW.a = 1",
	}))
	require.NoError(t, c.DefineTrigger(&nodes.CreateTriggerStmt{
		Name:     "z`_precedes",
		Table:    &nodes.TableRef{Name: "t`1"},
		Timing:   "BEFORE",
		Event:    "INSERT",
		BodyText: "SET NEW.a = 2",
		Order:    &nodes.TriggerOrder{Follows: false, TriggerName: "a_base"},
	}))

	meta := c.DatabaseMetadata("trg_db")
	triggers := meta.Schemas[0].Tables[0].Triggers
	require.Len(t, triggers, 2)
	require.Equal(t, "z`_precedes", triggers[0].Name)
	require.Equal(t, "a_base", triggers[1].Name)

	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "CREATE TRIGGER `z``_precedes`")
	require.Contains(t, def, "ON `t``1` FOR EACH ROW")
	require.Less(t, strings.Index(def, "CREATE TRIGGER `z``_precedes`"), strings.Index(def, "CREATE TRIGGER `a_base`"))

	loaded := New()
	_, err = loaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	loadMeta := loaded.DatabaseMetadata("trg_db")
	require.Equal(t, "z`_precedes", loadMeta.Schemas[0].Tables[0].Triggers[0].Name)
	require.Equal(t, "a_base", loadMeta.Schemas[0].Tables[0].Triggers[1].Name)
}

func TestDatabaseMetadataComputedViewColumnTypesAndBacktickRoutines(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE vdb;
		USE vdb;
		CREATE TABLE vt (id INT PRIMARY KEY);
		CREATE VIEW v_comp AS
			SELECT id + 1 AS next_id, COUNT(*) AS cnt, id > 0 AS is_pos, CONCAT('x', CAST(id AS CHAR(4))) AS label,
				1 AS lit_int, 'x' AS lit_str, '' AS lit_empty, NULL AS lit_null, TRUE AS lit_bool, 1.5 AS lit_dec
			FROM vt GROUP BY id;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	require.NoError(t, c.DefineFunction(&nodes.CreateFunctionStmt{
		Name:     &nodes.TableRef{Name: "fn`1"},
		Params:   []*nodes.FuncParam{{Name: "p`a", TypeName: &nodes.DataType{Name: "int"}}},
		Returns:  &nodes.DataType{Name: "int"},
		BodyText: "RETURN 1",
	}))
	require.NoError(t, c.DefineProcedure(&nodes.CreateFunctionStmt{
		IsProcedure: true,
		Name:        &nodes.TableRef{Name: "sp`1"},
		Params:      []*nodes.FuncParam{{Direction: "IN", Name: "in`x", TypeName: &nodes.DataType{Name: "int"}}},
		BodyText:    "SELECT 1",
	}))
	require.NoError(t, c.DefineEvent(&nodes.CreateEventStmt{
		Name:     "ev`1",
		Schedule: &nodes.EventSchedule{RawText: "EVERY 1 HOUR"},
		BodyText: "SELECT 1",
	}))

	meta := c.DatabaseMetadata("vdb")
	var compView *metadata.ViewMetadata
	for _, v := range meta.Schemas[0].Views {
		if v.Name == "v_comp" {
			compView = v
		}
	}
	require.NotNil(t, compView)
	viewCols := map[string]*metadata.ColumnMetadata{}
	for _, col := range compView.Columns {
		viewCols[col.Name] = col
	}
	require.Equal(t, "bigint", viewCols["next_id"].Type)
	require.Equal(t, "bigint", viewCols["cnt"].Type)
	require.Equal(t, "bigint", viewCols["is_pos"].Type)
	require.Equal(t, "varchar", viewCols["label"].Type)
	require.Equal(t, "bigint", viewCols["lit_int"].Type)
	require.Equal(t, "var_string(1)", viewCols["lit_str"].Type)
	require.Equal(t, "var_string(0)", viewCols["lit_empty"].Type)
	require.Equal(t, "binary(0)", viewCols["lit_null"].Type)
	require.Equal(t, "bigint", viewCols["lit_bool"].Type)
	require.Equal(t, "decimal(2,1)", viewCols["lit_dec"].Type)

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "FUNCTION `fn``1`(`p``a` INT)")
	require.Contains(t, def, "PROCEDURE `sp``1`(IN `in``x` INT)")
	require.Contains(t, def, "EVENT `ev``1`")
	require.Less(t, strings.Index(def, "Function structure for `fn``1`"), strings.Index(def, "View structure for `v_comp`"))
}
