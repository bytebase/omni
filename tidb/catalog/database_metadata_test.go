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
	require.Equal(t, "NULL", cols["nullable_col"].Default, "TiDB exports implicit NULL default on nullable columns matching live sync")
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

	// Vector index with embedded backtick exports valid quoted DDL and restores cleanly.
	vecBT := &metadata.DatabaseSchemaMetadata{
		Name: "vec_db",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:    "vec_t",
				Engine:  "InnoDB",
				Charset: "utf8mb4",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Position: 1, Type: "int", Nullable: false},
					{Name: "v", Position: 2, Type: "vector(3)", Nullable: true},
				},
				Indexes: []*metadata.IndexMetadata{
					{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"id"}, Type: "BTREE", Visible: true},
					{Name: "idx`vec", IsVector: true, Type: "HNSW", Expressions: []string{"(VEC_COSINE_DISTANCE(`v`))"}, Visible: true},
				},
			}},
		}},
	}
	vecDef, err := DatabaseDefinition(vecBT, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, vecDef, "VECTOR INDEX `idx``vec`")
	vecCat := New()
	res, err := vecCat.Exec("CREATE DATABASE vec_db;\nUSE vec_db;\n"+vecDef, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	for _, r := range res {
		require.NoError(t, r.Error)
	}
	restoredIdx := vecCat.DatabaseMetadata("vec_db").Schemas[0].Tables[0].Indexes[1]
	require.Equal(t, "idx`vec", restoredIdx.Name)
	require.True(t, restoredIdx.IsVector)
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
	require.Equal(t, "NULL", cols["bare"].Default)

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
	require.NotNil(t, reTbl.GetColumn("bare").Default)
	require.Equal(t, "NULL", *reTbl.GetColumn("bare").Default)

	reExport := reloaded.DatabaseMetadata("app")
	reCols := map[string]*metadata.ColumnMetadata{}
	for _, col := range reExport.Schemas[0].Tables[0].Columns {
		reCols[col.Name] = col
	}
	require.Equal(t, "NULL", reCols["s"].Default)
	require.Equal(t, "NULL", reCols["n"].Default)
	require.Equal(t, "NULL", reCols["bare"].Default)

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
	require.Contains(t, defRendered, "ALTER DATABASE PLACEMENT POLICY=`p_tbl`;")
	require.Contains(t, defRendered, "ALTER DATABASE SET TIFLASH REPLICA 1 LOCATION LABELS 'zone';")
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
	require.Equal(t, "var_string(5)", viewCols["label"].Type)
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

func TestDatabaseMetadataDecimalSessionContextBacktickAlterAndSkipDump(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE ddb;
		USE ddb;
		CREATE TABLE dec_t (id INT PRIMARY KEY, d1 DECIMAL(5,2), d2 DECIMAL(3,1));
		CREATE VIEW v_dec AS SELECT
			d1 + d2 AS p,
			d1 * d2 AS mul,
			d1 / d2 AS dv,
			d1 % d2 AS md,
			CAST(1 AS DECIMAL(5,2)) + CAST(2 AS DECIMAL(3,1)) AS cp
		FROM dec_t;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	decMeta := c.DatabaseMetadata("ddb")
	vcols := map[string]*metadata.ColumnMetadata{}
	for _, col := range decMeta.Schemas[0].Views[0].Columns {
		vcols[col.Name] = col
	}
	require.Equal(t, "decimal(6,2)", vcols["p"].Type)
	require.Equal(t, "decimal(8,3)", vcols["mul"].Type)
	require.Equal(t, "decimal(10,6)", vcols["dv"].Type)
	require.Equal(t, "decimal(5,2)", vcols["md"].Type)
	require.Equal(t, "decimal(6,2)", vcols["cp"].Type)

	srcMeta := &metadata.DatabaseSchemaMetadata{
		Name: "ctx_db",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:           "t`1",
				TiflashReplica: 1,
				CreateOptions:  "cached=on",
				SkipDump:       true,
				Columns: []*metadata.ColumnMetadata{
					{Name: "c`1", Type: "int", Nullable: false, DefaultDropped: true},
				},
				Triggers: []*metadata.TriggerMetadata{{
					Name:                "tr`1",
					Event:               "INSERT",
					Timing:              "BEFORE",
					Body:                "BEGIN END",
					SqlMode:             "PIPES_AS_CONCAT",
					CharacterSetClient:  "utf8mb4",
					CollationConnection: "utf8mb4_bin",
					SkipDump:            true,
				}},
			}},
			Views: []*metadata.ViewMetadata{{
				Name: "v_skip", Definition: "SELECT 1 AS c", SkipDump: true,
			}},
			Sequences: []*metadata.SequenceMetadata{{
				Name: "seq_skip", Start: "1", MinValue: "1", MaxValue: "100", Increment: "1", CacheSize: "1", SkipDump: true,
			}},
			Functions: []*metadata.FunctionMetadata{{
				Name:                "fn_ctx",
				Definition:          "CREATE FUNCTION fn_ctx() RETURNS INT DETERMINISTIC RETURN 1",
				SqlMode:             "ANSI_QUOTES",
				CharacterSetClient:  "utf8mb4",
				CollationConnection: "utf8mb4_unicode_ci",
				SkipDump:            true,
			}},
			Procedures: []*metadata.ProcedureMetadata{{
				Name:                "sp_ctx",
				Definition:          "CREATE PROCEDURE sp_ctx() BEGIN END",
				SqlMode:             "STRICT_TRANS_TABLES",
				CharacterSetClient:  "utf8mb4",
				CollationConnection: "utf8mb4_general_ci",
				SkipDump:            true,
			}},
			Events: []*metadata.EventMetadata{{
				Name:                "ev_ctx",
				Definition:          "CREATE EVENT ev_ctx ON SCHEDULE EVERY 1 HOUR DO SELECT 1",
				SqlMode:             "NO_AUTO_VALUE_ON_ZERO",
				CharacterSetClient:  "utf8mb4",
				CollationConnection: "utf8mb4_bin",
				TimeZone:            "+08:00",
			}},
		}},
	}
	loaded := New()
	report, err := loaded.LoadMetadata(context.Background(), srcMeta)
	require.NoError(t, err)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Degraded)

	out := loaded.DatabaseMetadata("ctx_db")
	tbl := out.Schemas[0].Tables[0]
	require.True(t, tbl.SkipDump)
	require.True(t, tbl.Triggers[0].SkipDump)
	require.Equal(t, "PIPES_AS_CONCAT", tbl.Triggers[0].SqlMode)
	require.True(t, out.Schemas[0].Views[0].SkipDump)
	require.True(t, out.Schemas[0].Sequences[0].SkipDump)
	require.True(t, out.Schemas[0].Functions[0].SkipDump)
	require.Equal(t, "ANSI_QUOTES", out.Schemas[0].Functions[0].SqlMode)
	require.True(t, out.Schemas[0].Procedures[0].SkipDump)
	require.Equal(t, "STRICT_TRANS_TABLES", out.Schemas[0].Procedures[0].SqlMode)
	require.Equal(t, "+08:00", out.Schemas[0].Events[0].TimeZone)

	def, err := DatabaseDefinition(out, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "CREATE TABLE `t``1`")
	require.Contains(t, def, "ALTER TABLE `t``1` ALTER COLUMN `c``1` DROP DEFAULT;")
	require.Contains(t, def, "ALTER TABLE `t``1` SET TIFLASH REPLICA 1;")
	require.Contains(t, def, "ALTER TABLE `t``1` CACHE;")
	require.Contains(t, def, "SET sql_mode = 'PIPES_AS_CONCAT'")
	require.Contains(t, def, "SET time_zone = '+08:00'")
}

func TestDatabaseMetadataImplicitFKDefaults(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE parent (id INT PRIMARY KEY);
		CREATE TABLE child (
			id INT PRIMARY KEY,
			pid INT,
			FOREIGN KEY (pid) REFERENCES parent(id)
		);
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	var childProto *metadata.TableMetadata
	for _, tbl := range meta.Schemas[0].Tables {
		if tbl.Name == "child" {
			childProto = tbl
		}
	}
	require.NotNil(t, childProto)
	require.Len(t, childProto.ForeignKeys, 1)
	require.Equal(t, "NONE", childProto.ForeignKeys[0].MatchType)
	require.Equal(t, "NO ACTION", childProto.ForeignKeys[0].OnDelete)
	require.Equal(t, "NO ACTION", childProto.ForeignKeys[0].OnUpdate)

	// Functional index expressions export parenthesized so DatabaseDefinition restores them.
	fnCat := New()
	_, err = fnCat.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE fn_idx_t (
			id INT NOT NULL PRIMARY KEY,
			code VARCHAR(20),
			INDEX fn_upper ((UPPER(code)))
		);
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	fnMeta := fnCat.DatabaseMetadata("app")
	fnDef, err := DatabaseDefinition(fnMeta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, fnDef, "KEY `fn_upper` ((upper(`code`)))")
	restoredFn := New()
	resFn, err := restoredFn.Exec("CREATE DATABASE app;\nUSE app;\n"+fnDef, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	for _, r := range resFn {
		require.NoError(t, r.Error)
	}
	require.NotNil(t, restoredFn.GetDatabase("app").GetTable("fn_idx_t"))

	// POSITION / LOCATE / INSTR export bigint and same-DB qualified FKs omit ReferencedSchema.
	posCat := New()
	_, err = posCat.Exec(`
		CREATE DATABASE other_db;
		CREATE TABLE other_db.ext_parent (id INT NOT NULL PRIMARY KEY);
		CREATE DATABASE app;
		USE app;
		CREATE TABLE fk_parent (id INT NOT NULL PRIMARY KEY);
		CREATE TABLE fk_qualified (
			id INT NOT NULL PRIMARY KEY,
			pid INT,
			eid INT,
			CONSTRAINT fk_same_qual FOREIGN KEY (pid) REFERENCES app.fk_parent(id),
			CONSTRAINT fk_cross_qual FOREIGN KEY (eid) REFERENCES other_db.ext_parent(id)
		);
		CREATE VIEW v_pos AS SELECT POSITION('x' IN 'xyz') AS p, LOCATE('x', 'xyz') AS l, INSTR('xyz', 'x') AS i;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	posMeta := posCat.DatabaseMetadata("app")
	var posView *metadata.ViewMetadata
	for _, v := range posMeta.Schemas[0].Views {
		if v.Name == "v_pos" {
			posView = v
		}
	}
	require.NotNil(t, posView)
	require.Len(t, posView.Columns, 3)
	for _, col := range posView.Columns {
		require.Equal(t, "bigint", col.Type, col.Name)
	}
	var qTbl *metadata.TableMetadata
	for _, tbl := range posMeta.Schemas[0].Tables {
		if tbl.Name == "fk_qualified" {
			qTbl = tbl
		}
	}
	require.NotNil(t, qTbl)
	require.Len(t, qTbl.ForeignKeys, 2)
	fkByName := map[string]*metadata.ForeignKeyMetadata{}
	for _, fk := range qTbl.ForeignKeys {
		fkByName[fk.Name] = fk
	}
	require.Empty(t, fkByName["fk_same_qual"].ReferencedSchema)
	require.Equal(t, "other_db", fkByName["fk_cross_qual"].ReferencedSchema)

	// Set-operation views derive UNION metadata across arms while INTERSECT/EXCEPT keep the left arm (TiDB v8.5.5).
	unionCat := New()
	_, err = unionCat.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE t1 (id INT NOT NULL, s CHAR(4) NOT NULL);
		CREATE TABLE t2 (id BIGINT NULL, s VARCHAR(10) NULL);
		CREATE VIEW v_mix AS SELECT 1 AS c UNION SELECT 'x';
		CREATE VIEW v_tbl AS SELECT id, s FROM t1 UNION SELECT id, s FROM t2;
		CREATE VIEW v_null AS SELECT id FROM t1 UNION SELECT NULL;
		CREATE VIEW v_num AS SELECT 1 AS a UNION SELECT 1.5 AS a;
		CREATE VIEW v_except AS SELECT id, s FROM t1 EXCEPT SELECT id, s FROM t2;
		CREATE VIEW v_intersect AS SELECT id, s FROM t1 INTERSECT SELECT id, s FROM t2;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	uViews := map[string]*metadata.ViewMetadata{}
	for _, v := range unionCat.DatabaseMetadata("app").Schemas[0].Views {
		uViews[v.Name] = v
	}
	require.Equal(t, "varchar(20)", uViews["v_mix"].Columns[0].Type)
	require.Equal(t, "bigint", uViews["v_tbl"].Columns[0].Type)
	require.Equal(t, "varchar(10)", uViews["v_tbl"].Columns[1].Type)
	require.Equal(t, "int", uViews["v_null"].Columns[0].Type)
	require.Equal(t, "decimal(20,1)", uViews["v_num"].Columns[0].Type)
	require.Equal(t, "int", uViews["v_except"].Columns[0].Type)
	require.Equal(t, "char(4)", uViews["v_except"].Columns[1].Type)
	require.Equal(t, "int", uViews["v_intersect"].Columns[0].Type)
	require.Equal(t, "char(4)", uViews["v_intersect"].Columns[1].Type)

	// MIN/MAX view columns derive their types from their arguments, and implicit
	// nullable TiDB columns export Default="NULL" unless DefaultDropped is set.
	mmCat := New()
	_, err = mmCat.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE t_mm (id INT NOT NULL, nn_nod INT NOT NULL, nul INT NULL, d_drop INT DEFAULT 1, s VARCHAR(20) NOT NULL, d DATE NOT NULL, dec1 DECIMAL(10,2) NOT NULL);
		ALTER TABLE t_mm ALTER COLUMN d_drop DROP DEFAULT;
		CREATE VIEW v_mm AS SELECT MIN(id) AS min_id, MAX(s) AS max_s, MIN(d) AS min_d, MAX(dec1) AS max_dec FROM t_mm;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	mmMeta := mmCat.DatabaseMetadata("app")
	var vMM *metadata.ViewMetadata
	for _, v := range mmMeta.Schemas[0].Views {
		if v.Name == "v_mm" {
			vMM = v
		}
	}
	require.NotNil(t, vMM)
	require.Equal(t, "int", vMM.Columns[0].Type)
	require.Equal(t, "varchar(20)", vMM.Columns[1].Type)
	require.Equal(t, "date", vMM.Columns[2].Type)
	require.Equal(t, "decimal(10,2)", vMM.Columns[3].Type)

	mmCols := map[string]*metadata.ColumnMetadata{}
	for _, col := range mmMeta.Schemas[0].Tables[0].Columns {
		mmCols[col.Name] = col
	}
	require.Equal(t, "NULL", mmCols["nul"].Default)
	require.Empty(t, mmCols["nn_nod"].Default)
	require.Empty(t, mmCols["d_drop"].Default)
	require.True(t, mmCols["d_drop"].DefaultDropped)

	// Unresolved view columns retain supplied snapshot metadata across LoadMetadata -> DatabaseMetadata.
	extCat := New()
	_, err = extCat.LoadMetadata(context.Background(), &metadata.DatabaseSchemaMetadata{
		Name: "app",
		Schemas: []*metadata.SchemaMetadata{{
			Views: []*metadata.ViewMetadata{{
				Name:       "v_ext",
				Definition: "select `unloaded_db`.`ext_tbl`.`id` AS `id`,`unloaded_db`.`ext_tbl`.`name` AS `name` from `unloaded_db`.`ext_tbl`",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Position: 1, Type: "int", Nullable: true},
					{Name: "name", Position: 2, Type: "varchar(32)", Nullable: true},
				},
			}},
		}},
	})
	require.NoError(t, err)
	reExt := extCat.DatabaseMetadata("app").Schemas[0].Views[0]
	require.Len(t, reExt.Columns, 2)
	require.Equal(t, "id", reExt.Columns[0].Name)
	require.Equal(t, "int", reExt.Columns[0].Type)
	require.Equal(t, "name", reExt.Columns[1].Name)
	require.Equal(t, "varchar(32)", reExt.Columns[1].Type)

	// SUM/AVG view columns derive result types from their arguments matching TiDB v8.5.5.
	aggCat := New()
	_, err = aggCat.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE t_agg (i INT NOT NULL, b BIGINT NOT NULL, f FLOAT NOT NULL, d DOUBLE NOT NULL, dec1 DECIMAL(10,2) NOT NULL);
		CREATE VIEW v_agg AS SELECT
			SUM(i) AS s_i, SUM(b) AS s_b, SUM(f) AS s_f, SUM(d) AS s_d, SUM(dec1) AS s_dec,
			AVG(i) AS a_i, AVG(b) AS a_b, AVG(f) AS a_f, AVG(d) AS a_d, AVG(dec1) AS a_dec
		FROM t_agg;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	aggCols := map[string]*metadata.ColumnMetadata{}
	for _, col := range aggCat.DatabaseMetadata("app").Schemas[0].Views[0].Columns {
		aggCols[col.Name] = col
	}
	require.Equal(t, "decimal(32,0)", aggCols["s_i"].Type)
	require.Equal(t, "decimal(41,0)", aggCols["s_b"].Type)
	require.Equal(t, "double", aggCols["s_f"].Type)
	require.Equal(t, "double", aggCols["s_d"].Type)
	require.Equal(t, "decimal(32,2)", aggCols["s_dec"].Type)
	require.Equal(t, "decimal(11,0)", aggCols["a_i"].Type)
	require.Equal(t, "decimal(20,0)", aggCols["a_b"].Type)
	require.Equal(t, "double", aggCols["a_f"].Type)
	require.Equal(t, "double", aggCols["a_d"].Type)
	require.Equal(t, "decimal(10,2)", aggCols["a_dec"].Type)

	// CASE and control-flow functions (COALESCE, IFNULL, NULLIF, IF) derive
	// result column types matching TiDB v8.5.5.
	ctrlCat := New()
	_, err = ctrlCat.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE t_ctrl (
			ni INT NULL,
			ti TINYINT NOT NULL,
			id INT NOT NULL,
			d DECIMAL(10,2) NOT NULL,
			f FLOAT NOT NULL,
			c CHAR(4) NOT NULL,
			c6 CHAR(6) NOT NULL,
			v VARCHAR(20) NOT NULL,
			dt DATE NOT NULL,
			ts TIMESTAMP NOT NULL
		);
		CREATE VIEW v_ctrl AS SELECT
			COALESCE(id, 0) AS coal_int,
			COALESCE(NULL, id, 0) AS coal_null_int,
			COALESCE(id, d) AS coal_int_dec,
			COALESCE(c, c6) AS coal_char,
			IFNULL(ti, ti) AS ifnull_ti,
			IFNULL(id, 0) AS ifnull_int,
			IFNULL(v, 'abc') AS ifnull_str,
			NULLIF(ti, 1) AS nullif_ti,
			NULLIF(id, 0) AS nullif_int,
			NULLIF(c, 'xy') AS nullif_char,
			NULLIF(d, 0) AS nullif_dec,
			IF(id > 0, id, 0) AS if_int,
			IF(id > 0, id, 1.5) AS if_int_dec,
			IF(id > 0, NULL, id) AS if_null_int,
			IF(id > 0, c, v) AS if_char_varchar,
			IF(id > 0, dt, ts) AS if_dt_ts,
			CASE WHEN id > 0 THEN id ELSE 0 END AS case_int,
			CASE WHEN id > 0 THEN ti END AS case_ti_noelse,
			CASE id WHEN 1 THEN 'a' ELSE 'bbb' END AS case_str,
			CASE WHEN id > 0 THEN f ELSE d END AS case_f_d
		FROM t_ctrl;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	ctrlCols := map[string]*metadata.ColumnMetadata{}
	for _, col := range ctrlCat.DatabaseMetadata("app").Schemas[0].Views[0].Columns {
		ctrlCols[col.Name] = col
	}
	require.Equal(t, "bigint", ctrlCols["coal_int"].Type)
	require.Equal(t, "bigint", ctrlCols["coal_null_int"].Type)
	require.Equal(t, "decimal(13,2)", ctrlCols["coal_int_dec"].Type)
	require.Equal(t, "char(6)", ctrlCols["coal_char"].Type)
	require.Equal(t, "tinyint", ctrlCols["ifnull_ti"].Type)
	require.Equal(t, "bigint", ctrlCols["ifnull_int"].Type)
	require.Equal(t, "varchar(20)", ctrlCols["ifnull_str"].Type)
	require.Equal(t, "tinyint", ctrlCols["nullif_ti"].Type)
	require.Equal(t, "int", ctrlCols["nullif_int"].Type)
	require.Equal(t, "char(4)", ctrlCols["nullif_char"].Type)
	require.Equal(t, "decimal(10,2)", ctrlCols["nullif_dec"].Type)
	require.Equal(t, "bigint", ctrlCols["if_int"].Type)
	require.Equal(t, "decimal(12,1)", ctrlCols["if_int_dec"].Type)
	require.Equal(t, "int", ctrlCols["if_null_int"].Type)
	require.Equal(t, "varchar(20)", ctrlCols["if_char_varchar"].Type)
	require.Equal(t, "datetime", ctrlCols["if_dt_ts"].Type)
	require.Equal(t, "bigint", ctrlCols["case_int"].Type)
	require.Equal(t, "tinyint", ctrlCols["case_ti_noelse"].Type)
	require.Equal(t, "varchar(3)", ctrlCols["case_str"].Type)
	require.Equal(t, "double", ctrlCols["case_f_d"].Type)

	// Scalar function view columns resolve types matching TiDB v8.5.5.
	scalarFnCat := New()
	_, err = scalarFnCat.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE t_fn (
			i INT NOT NULL,
			f FLOAT NOT NULL,
			dbl DOUBLE NOT NULL,
			d DECIMAL(10,2) NOT NULL,
			c CHAR(4) NOT NULL,
			v VARCHAR(32) NOT NULL,
			txt TEXT NOT NULL,
			j JSON NOT NULL,
			dt DATE NOT NULL,
			dtm6 DATETIME(6) NOT NULL,
			tm4 TIME(4) NOT NULL
		);
		CREATE VIEW v_fn AS SELECT
			JSON_UNQUOTE(JSON_EXTRACT(j, '$.a')) AS ju,
			JSON_EXTRACT(j, '$.a') AS je,
			JSON_TYPE(j) AS jt,
			LOWER(v) AS lw_v,
			UPPER(c) AS up_c,
			TRIM(txt) AS tr_txt,
			REPLACE(v, 'a', 'bb') AS rp_v,
			SUBSTRING(v, 2, 5) AS sub_v,
			LEFT(v, 5) AS lf_v,
			RIGHT(c, 2) AS rg_c,
			LPAD(v, 10, 'x') AS lp_v,
			REPEAT(c, 3) AS rep_c,
			CONCAT(c, v) AS cc_cv,
			CONCAT_WS(',', c, v) AS ccw,
			SPACE(5) AS sp_5,
			FORMAT(d, 2) AS fmt_d,
			ABS(i) AS abs_i,
			ABS(f) AS abs_f,
			ABS(dbl) AS abs_dbl,
			ABS(d) AS abs_d,
			CEIL(d) AS ceil_d,
			FLOOR(d) AS fl_d,
			ROUND(d, 1) AS rnd_d,
			ROUND(dbl, 2) AS rnd_dbl,
			TRUNCATE(d, 1) AS trn_d,
			MOD(i, 3) AS mod_i,
			MOD(d, 3) AS mod_d,
			LENGTH(v) AS len_v,
			JSON_VALID(v) AS jv,
			NOW(6) AS fn_now6,
			NOW(0) AS fn_now0,
			CURRENT_TIMESTAMP(3) AS fn_ct3,
			SYSDATE(6) AS fn_sys6,
			LOCALTIMESTAMP(3) AS fn_lts3,
			UTC_TIMESTAMP(6) AS fn_utcts6,
			CURTIME(6) AS fn_curt6,
			CURRENT_TIME(3) AS fn_curt3,
			TIME(dtm6) AS fn_time_dtm6,
			TIMESTAMP(dt, tm4) AS fn_ts_dt_tm4
		FROM t_fn;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	scalarFnCols := map[string]*metadata.ColumnMetadata{}
	for _, col := range scalarFnCat.DatabaseMetadata("app").Schemas[0].Views[0].Columns {
		scalarFnCols[col.Name] = col
	}
	require.Equal(t, "longtext", scalarFnCols["ju"].Type)
	require.Equal(t, "json", scalarFnCols["je"].Type)
	require.Equal(t, "var_string(51)", scalarFnCols["jt"].Type)
	require.Equal(t, "var_string(32)", scalarFnCols["lw_v"].Type)
	require.Equal(t, "var_string(4)", scalarFnCols["up_c"].Type)
	require.Equal(t, "var_string(65535)", scalarFnCols["tr_txt"].Type)
	require.Equal(t, "var_string(64)", scalarFnCols["rp_v"].Type)
	require.Equal(t, "var_string(32)", scalarFnCols["sub_v"].Type)
	require.Equal(t, "var_string(32)", scalarFnCols["lf_v"].Type)
	require.Equal(t, "var_string(4)", scalarFnCols["rg_c"].Type)
	require.Equal(t, "var_string(40)", scalarFnCols["lp_v"].Type)
	require.Equal(t, "longtext", scalarFnCols["rep_c"].Type)
	require.Equal(t, "var_string(36)", scalarFnCols["cc_cv"].Type)
	require.Equal(t, "var_string(37)", scalarFnCols["ccw"].Type)
	require.Equal(t, "longtext", scalarFnCols["sp_5"].Type)
	require.Equal(t, "longtext", scalarFnCols["fmt_d"].Type)
	require.Equal(t, "bigint", scalarFnCols["abs_i"].Type)
	require.Equal(t, "double", scalarFnCols["abs_f"].Type)
	require.Equal(t, "double", scalarFnCols["abs_dbl"].Type)
	require.Equal(t, "decimal(10,2)", scalarFnCols["abs_d"].Type)
	require.Equal(t, "bigint", scalarFnCols["ceil_d"].Type)
	require.Equal(t, "bigint", scalarFnCols["fl_d"].Type)
	require.Equal(t, "decimal(10,1)", scalarFnCols["rnd_d"].Type)
	require.Equal(t, "double", scalarFnCols["rnd_dbl"].Type)
	require.Equal(t, "decimal(9,1)", scalarFnCols["trn_d"].Type)
	require.Equal(t, "bigint", scalarFnCols["mod_i"].Type)
	require.Equal(t, "decimal(10,2)", scalarFnCols["mod_d"].Type)
	require.Equal(t, "bigint", scalarFnCols["len_v"].Type)
	require.Equal(t, "bigint", scalarFnCols["jv"].Type)
	require.Equal(t, "datetime(6)", scalarFnCols["fn_now6"].Type)
	require.Equal(t, "datetime", scalarFnCols["fn_now0"].Type)
	require.Equal(t, "datetime(3)", scalarFnCols["fn_ct3"].Type)
	require.Equal(t, "datetime(6)", scalarFnCols["fn_sys6"].Type)
	require.Equal(t, "datetime(3)", scalarFnCols["fn_lts3"].Type)
	require.Equal(t, "datetime(6)", scalarFnCols["fn_utcts6"].Type)
	require.Equal(t, "time(6)", scalarFnCols["fn_curt6"].Type)
	require.Equal(t, "time(3)", scalarFnCols["fn_curt3"].Type)
	require.Equal(t, "time(6)", scalarFnCols["fn_time_dtm6"].Type)
	require.Equal(t, "datetime(4)", scalarFnCols["fn_ts_dt_tm4"].Type)

	// LoadMetadata -> DatabaseMetadata round-trip preserves TiDB routine ENUM, SET,
	// COLLATE, and BINARY type modifiers without dropping them.
	rtCat := New()
	rtReport, err := rtCat.LoadMetadata(context.Background(), &metadata.DatabaseSchemaMetadata{
		Name: "app",
		Schemas: []*metadata.SchemaMetadata{{
			Functions: []*metadata.FunctionMetadata{{
				Name:       "fn_rich",
				Definition: "CREATE FUNCTION `fn_rich`(a ENUM('Alpha','b'), b VARCHAR(10) CHARSET latin1 COLLATE latin1_bin, c CHAR(4) BINARY) RETURNS set('r','w') CHARSET utf8mb4 COLLATE utf8mb4_bin DETERMINISTIC RETURN 'r'",
			}},
			Procedures: []*metadata.ProcedureMetadata{{
				Name:       "sp_rich",
				Definition: "CREATE PROCEDURE `sp_rich`(IN p1 SET('x','y'), IN p2 CHAR(8) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci) BEGIN SELECT 1; END",
			}},
		}},
	})
	require.NoError(t, err)
	require.Empty(t, rtReport.Missing)
	require.Empty(t, rtReport.Degraded)
	rtMeta := rtCat.DatabaseMetadata("app")
	require.Len(t, rtMeta.Schemas[0].Functions, 1)
	require.Len(t, rtMeta.Schemas[0].Procedures, 1)
	require.Contains(t, rtMeta.Schemas[0].Functions[0].Definition, "ENUM('Alpha','b')")
	require.Contains(t, rtMeta.Schemas[0].Functions[0].Definition, "CHARSET latin1 COLLATE latin1_bin")
	require.Contains(t, rtMeta.Schemas[0].Functions[0].Definition, "CHAR(4) binary")
	require.Contains(t, rtMeta.Schemas[0].Functions[0].Definition, "RETURNS set('r','w') CHARSET utf8mb4 COLLATE utf8mb4_bin")
	require.Contains(t, rtMeta.Schemas[0].Procedures[0].Definition, "IN p1 SET('x','y')")
	require.Contains(t, rtMeta.Schemas[0].Procedures[0].Definition, "CHAR(8) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci")

	rtCat2 := New()
	rtReport2, err := rtCat2.LoadMetadata(context.Background(), rtMeta)
	require.NoError(t, err)
	require.Empty(t, rtReport2.Missing)
	require.Empty(t, rtReport2.Degraded)
	rtMeta2 := rtCat2.DatabaseMetadata("app")
	require.Equal(t, rtMeta.Schemas[0].Functions[0].Definition, rtMeta2.Schemas[0].Functions[0].Definition)
	require.Equal(t, rtMeta.Schemas[0].Procedures[0].Definition, rtMeta2.Schemas[0].Procedures[0].Definition)
}

func TestViewMetadataScalarSubqueryAndTemporalArithmetic(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE t (
			id INT NOT NULL PRIMARY KEY,
			s VARCHAR(20) NOT NULL,
			dec_col DECIMAL(10,2) NOT NULL,
			d DATE NOT NULL,
			dt DATETIME NOT NULL,
			dt3 DATETIME(3) NOT NULL,
			ts TIMESTAMP(6) NOT NULL,
			tm TIME NOT NULL,
			tm2 TIME(2) NOT NULL
		);
		CREATE VIEW v AS SELECT
			(SELECT id FROM t LIMIT 1) AS sq_id,
			(SELECT s FROM t LIMIT 1) AS sq_s,
			(SELECT dec_col FROM t LIMIT 1) AS sq_dec,
			(SELECT dt3 FROM t LIMIT 1) AS sq_dt3,
			(SELECT COUNT(*) FROM t) AS sq_cnt,
			(SELECT 'abc') AS sq_str,
			EXISTS(SELECT 1 FROM t) AS ex,
			id IN (SELECT id FROM t) AS in_col,
			TIMESTAMPDIFF(DAY, d, dt) AS td_day,
			TIMESTAMPDIFF(MICROSECOND, dt, dt3) AS td_us,
			DATEDIFF(dt, d) AS dd,
			DATE_ADD(d, INTERVAL 1 DAY) AS da_d_day,
			DATE_ADD(d, INTERVAL 1 HOUR) AS da_d_hr,
			DATE_ADD(d, INTERVAL 1 MICROSECOND) AS da_d_us,
			DATE_ADD(dt3, INTERVAL 1 DAY) AS da_dt3_day,
			DATE_ADD(ts, INTERVAL 1 DAY) AS da_ts_day,
			DATE_ADD(tm2, INTERVAL 1 MICROSECOND) AS da_tm2_us,
			DATE_ADD(s, INTERVAL 1 DAY) AS da_s_day,
			DATE_SUB(d, INTERVAL 1 DAY) AS ds_d_day,
			ADDDATE(d, 5) AS ad_d_int,
			SUBDATE(d, 5) AS sd_d_int,
			TIMESTAMPADD(DAY, 1, d) AS ta_d_day,
			TIMESTAMPADD(MICROSECOND, 1, dt) AS ta_dt_us,
			d + INTERVAL 1 DAY AS op_d_day,
			d + INTERVAL 1 HOUR AS op_d_hr,
			INTERVAL 1 DAY + d AS op_iv_d,
			d - INTERVAL 1 DAY AS op_d_sub,
			TIMEDIFF(dt3, dt) AS tdiff_dt,
			TIMEDIFF(tm2, tm) AS tdiff_tm,
			COALESCE((SELECT id FROM t LIMIT 1), 0) AS sq_coalesce,
			COALESCE((SELECT id FROM t LIMIT 1), id) AS sq_coalesce_int,
			(SELECT d FROM t LIMIT 1) + INTERVAL 1 DAY AS sq_interval
		FROM t;
	`, nil)
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	var v *metadata.ViewMetadata
	for _, vm := range meta.GetSchemas()[0].GetViews() {
		if vm.Name == "v" {
			v = vm
			break
		}
	}
	require.NotNil(t, v)
	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range v.GetColumns() {
		cols[col.Name] = col
	}

	require.Equal(t, "int", cols["sq_id"].Type)
	require.True(t, cols["sq_id"].Nullable)
	require.Equal(t, "varchar(20)", cols["sq_s"].Type)
	require.Equal(t, "decimal(10,2)", cols["sq_dec"].Type)
	require.Equal(t, "datetime(3)", cols["sq_dt3"].Type)
	require.Equal(t, "bigint", cols["sq_cnt"].Type)
	require.Equal(t, "var_string(3)", cols["sq_str"].Type)
	require.Equal(t, "tinyint", cols["ex"].Type)
	require.Equal(t, "tinyint", cols["in_col"].Type)

	require.Equal(t, "bigint", cols["td_day"].Type)
	require.Equal(t, "bigint", cols["td_us"].Type)
	require.Equal(t, "bigint", cols["dd"].Type)
	require.Equal(t, "date", cols["da_d_day"].Type)
	require.Equal(t, "datetime", cols["da_d_hr"].Type)
	require.Equal(t, "datetime(6)", cols["da_d_us"].Type)
	require.Equal(t, "datetime(3)", cols["da_dt3_day"].Type)
	require.Equal(t, "datetime(6)", cols["da_ts_day"].Type)
	require.Equal(t, "time(6)", cols["da_tm2_us"].Type)
	require.Equal(t, "var_string(29)", cols["da_s_day"].Type)
	require.Equal(t, "date", cols["ds_d_day"].Type)
	require.Equal(t, "date", cols["ad_d_int"].Type)
	require.Equal(t, "date", cols["sd_d_int"].Type)
	require.Equal(t, "var_string(19)", cols["ta_d_day"].Type)
	require.Equal(t, "var_string(26)", cols["ta_dt_us"].Type)
	require.Equal(t, "date", cols["op_d_day"].Type)
	require.Equal(t, "datetime", cols["op_d_hr"].Type)
	require.Equal(t, "date", cols["op_iv_d"].Type)
	require.Equal(t, "date", cols["op_d_sub"].Type)
	require.Equal(t, "time(3)", cols["tdiff_dt"].Type)
	require.Equal(t, "time(2)", cols["tdiff_tm"].Type)
	require.Equal(t, "bigint", cols["sq_coalesce"].Type)
	require.Equal(t, "int", cols["sq_coalesce_int"].Type)
	require.Equal(t, "date", cols["sq_interval"].Type)
}

func TestIndexColumnWithParenthesisAndBacktickRoundTrip(t *testing.T) {
	c := New()
	_, err := c.Exec("CREATE DATABASE app; USE app; CREATE TABLE idx_t (`a(b` INT, `x``y` INT, INDEX `i` (`a(b`), INDEX `ib` (`x``y`));", nil)
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "KEY `i` (`a(b`)")
	require.Contains(t, def, "KEY `ib` (`x``y`)")

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Incomplete)

	reMeta := reloaded.DatabaseMetadata("app")
	require.Len(t, reMeta.GetSchemas()[0].GetTables()[0].GetIndexes(), 2)
	require.Equal(t, []string{"a(b"}, reMeta.GetSchemas()[0].GetTables()[0].GetIndexes()[0].GetExpressions())
	require.Equal(t, []string{"x`y"}, reMeta.GetSchemas()[0].GetTables()[0].GetIndexes()[1].GetExpressions())
}

func TestViewLikeExprMetadata(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE t (
			id INT NOT NULL PRIMARY KEY,
			s_nn VARCHAR(20) NOT NULL,
			s_null VARCHAR(20) NULL
		);
		CREATE VIEW v_like AS SELECT
			s_nn LIKE 'a%' AS l_nn,
			s_null LIKE 'a%' AS l_null,
			s_nn NOT LIKE 'a%' AS nl_nn,
			s_null NOT LIKE 'a%' AS nl_null,
			s_nn LIKE s_null AS l_pat_null,
			s_nn LIKE 'a!%' ESCAPE '!' AS l_esc_nn,
			s_null LIKE 'a!%' ESCAPE '!' AS l_esc_null
		FROM t;
	`, nil)
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range meta.GetSchemas()[0].GetViews()[0].GetColumns() {
		cols[col.Name] = col
	}

	for _, name := range []string{"l_nn", "l_null", "nl_nn", "nl_null", "l_pat_null", "l_esc_nn", "l_esc_null"} {
		require.Equal(t, "bigint", cols[name].Type, name)
		require.True(t, cols[name].Nullable, name)
	}
}

func TestViewNonArithmeticOperatorMetadata(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20) NOT NULL, doc JSON);
		CREATE VIEW v_ops AS SELECT
			id & 1 AS c_band,
			id XOR 1 AS c_lxor,
			doc->'$.x' AS c_jextract,
			doc->>'$.x' AS c_junquote,
			BINARY name AS c_bin
		FROM t;
	`, nil)
	require.NoError(t, err)

	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range c.DatabaseMetadata("app").GetSchemas()[0].GetViews()[0].GetColumns() {
		cols[col.Name] = col
	}
	require.Equal(t, "bigint unsigned", cols["c_band"].Type)
	require.Equal(t, "bigint", cols["c_lxor"].Type)
	require.Equal(t, "json", cols["c_jextract"].Type)
	require.Equal(t, "longtext", cols["c_junquote"].Type)
	require.Equal(t, "binary(1)", cols["c_bin"].Type)
}
