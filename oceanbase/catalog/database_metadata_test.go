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
		CREATE DATABASE app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
		USE app;

		CREATE TABLE parent (
			id INT AUTO_INCREMENT PRIMARY KEY,
			code VARCHAR(20) NOT NULL DEFAULT 'abc',
			expr_col INT DEFAULT (1 + 2),
			bit_col BIT(3) DEFAULT b'101',
			nullable_col INT,
			ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			gen_stored VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_s')) STORED,
			gen_virt VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_v')) VIRTUAL,
			g0 GEOMETRY NOT NULL SRID 0,
			secret INT INVISIBLE NOT NULL DEFAULT 0,
			SPATIAL INDEX sp_idx (g0)
		) ENGINE=InnoDB ROW_FORMAT=DYNAMIC KEY_BLOCK_SIZE=8 STATS_PERSISTENT=0 COMPRESSION='lz4_1.0' COMMENT='parent table';

		CREATE TABLE child (
			id INT PRIMARY KEY,
			parent_id INT,
			score INT,
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id) ON DELETE CASCADE,
			CONSTRAINT chk_score CHECK (score >= 0)
		);

		CREATE TABLE part_tbl (
			id INT NOT NULL,
			yr INT NOT NULL,
			PRIMARY KEY (id, yr)
		) PARTITION BY RANGE (yr)
		SUBPARTITION BY HASH (id) SUBPARTITIONS 2 (
			PARTITION p0 VALUES LESS THAN (2025),
			PARTITION p1 VALUES LESS THAN MAXVALUE
		);

		CREATE VIEW v_parent (vid, vcode) AS SELECT id, code FROM parent;
		CREATE MATERIALIZED VIEW mv_parent AS SELECT id, code FROM parent;

		CREATE FUNCTION fn_add(a INT, b INT) RETURNS INT DETERMINISTIC COMMENT 'adds two ints'
		BEGIN
			RETURN a + b;
		END;

		CREATE PROCEDURE proc_noop() COMMENT 'noop proc'
		BEGIN
			SELECT 1;
		END;

		CREATE TRIGGER trg_parent AFTER INSERT ON parent FOR EACH ROW
		BEGIN
			INSERT INTO child(id) VALUES (NEW.id);
		END;

		CREATE TRIGGER trg_z_first AFTER INSERT ON parent FOR EACH ROW PRECEDES trg_parent
		BEGIN
			INSERT INTO child(id) VALUES (NEW.id + 1);
		END;

		CREATE EVENT ev_daily ON SCHEDULE EVERY 1 DAY DO SELECT 1;
	`, nil)
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	require.Equal(t, "app", meta.Name)
	require.Equal(t, "utf8mb4", meta.CharacterSet)
	require.Equal(t, "utf8mb4_unicode_ci", meta.Collation)
	require.Len(t, meta.Schemas, 1)

	s := meta.Schemas[0]
	require.Len(t, s.Tables, 3)

	parent := s.Tables[1]
	require.Equal(t, "parent", parent.Name)
	require.Equal(t, "parent table", parent.Comment)
	require.Contains(t, parent.CreateOptions, "row_format=DYNAMIC")
	require.Contains(t, parent.CreateOptions, "key_block_size=8")
	require.Contains(t, parent.CreateOptions, "stats_persistent=0")
	require.Contains(t, parent.CreateOptions, "compression='lz4_1.0'")

	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range parent.Columns {
		cols[col.Name] = col
	}
	require.Equal(t, "AUTO_INCREMENT", cols["id"].Default)
	require.Equal(t, "'abc'", cols["code"].Default)
	require.Equal(t, "((1 + 2))", cols["expr_col"].Default)
	require.Equal(t, "b'101'", cols["bit_col"].Default)
	require.Equal(t, "NULL", cols["nullable_col"].Default)
	require.Equal(t, "CURRENT_TIMESTAMP", cols["ts"].OnUpdate)
	require.Equal(t, metadata.GenerationMetadata_TYPE_STORED, cols["gen_stored"].Generation.Type)
	require.Equal(t, metadata.GenerationMetadata_TYPE_VIRTUAL, cols["gen_virt"].Generation.Type)
	require.NotNil(t, cols["g0"].Srid)
	require.Equal(t, uint32(0), *cols["g0"].Srid)
	require.True(t, cols["secret"].IsInvisible)

	require.Len(t, parent.Triggers, 2)
	require.Equal(t, []string{"trg_z_first", "trg_parent"}, []string{
		parent.Triggers[0].Name, parent.Triggers[1].Name,
	})

	require.Len(t, s.Views, 1)
	require.Equal(t, "v_parent", s.Views[0].Name)
	require.Len(t, s.Views[0].Columns, 2)
	require.Equal(t, "vid", s.Views[0].Columns[0].Name)
	require.Len(t, s.MaterializedViews, 1)
	require.Equal(t, "mv_parent", s.MaterializedViews[0].Name)
	require.Len(t, s.Functions, 1)
	require.Contains(t, s.Functions[0].Definition, "FUNCTION `fn_add`")
	require.Equal(t, "adds two ints", s.Functions[0].Comment)
	require.Len(t, s.Procedures, 1)
	require.Contains(t, s.Procedures[0].Definition, "PROCEDURE `proc_noop`")
	require.Equal(t, "noop proc", s.Procedures[0].Comment)
	require.Len(t, s.Events, 1)

	// Round-trip through LoadMetadata preserves objects cleanly.
	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)
	require.NotNil(t, reloaded.GetDatabase("app").Functions["fn_add"])
	require.NotNil(t, reloaded.GetDatabase("app").Procedures["proc_noop"])
	reloadedExported := reloaded.DatabaseMetadata("app")
	require.Equal(t, []string{"trg_z_first", "trg_parent"}, []string{
		reloadedExported.Schemas[0].Tables[1].Triggers[0].Name,
		reloadedExported.Schemas[0].Tables[1].Triggers[1].Name,
	})

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
	require.Equal(t, "parent_id", fkCon.Columns[0])
	require.Len(t, child.CheckConstraints, 1)

	partTbl := s.Tables[2]
	require.Len(t, partTbl.Partitions, 2)
	require.Len(t, partTbl.Partitions[0].Subpartitions, 2)
}

func TestDatabaseMetadataExcludesTempTablesSortsIndexesAndAutoPartitions(t *testing.T) {
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
		CREATE TABLE auto_part (
			id INT NOT NULL
		) PARTITION BY RANGE (id) PARTITION SIZE '1G' PARTITIONS AUTO;
		CREATE TABLE split_part (
			id INT NOT NULL
		) PARTITION BY RANGE (id) SIZE ('512M') (
			PARTITION p0 VALUES LESS THAN (100)
		);
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	require.Len(t, meta.Schemas[0].Tables, 3)
	require.Equal(t, "auto_part", meta.Schemas[0].Tables[0].Name)
	require.Equal(t, "items", meta.Schemas[0].Tables[1].Name)
	require.Equal(t, "split_part", meta.Schemas[0].Tables[2].Name)

	items := meta.Schemas[0].Tables[1]
	require.Equal(t, []string{"PRIMARY", "idx_alpha", "idx_zeta"}, []string{
		items.Indexes[0].Name,
		items.Indexes[1].Name,
		items.Indexes[2].Name,
	})

	autoPart := meta.Schemas[0].Tables[0]
	require.Len(t, autoPart.Partitions, 1)
	require.Equal(t, metadata.TablePartitionMetadata_RANGE, autoPart.Partitions[0].Type)
	require.Equal(t, "`id`", autoPart.Partitions[0].Expression)
	require.Equal(t, "'1G'", autoPart.Partitions[0].AutoSize)

	splitPart := meta.Schemas[0].Tables[2]
	require.Len(t, splitPart.Partitions, 1)
	require.Equal(t, "512M", splitPart.Partitions[0].SplitSize)

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	reloadedAuto := reloaded.GetDatabase("app").GetTable("auto_part")
	require.NotNil(t, reloadedAuto.Partitioning)
	require.Equal(t, "'1G'", reloadedAuto.Partitioning.AutoSize)
	require.Empty(t, reloadedAuto.Partitioning.Partitions)

	reloadedSplit := reloaded.GetDatabase("app").GetTable("split_part")
	require.NotNil(t, reloadedSplit.Partitioning)
	require.Equal(t, "512M", reloadedSplit.Partitioning.SplitSize)
	require.Len(t, reloadedSplit.Partitioning.Partitions, 1)
}

func TestDatabaseMetadataRoundTripOptionsConstraintsDefaultsAndSubpartitions(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app COLLATE utf8mb4_unicode_ci;
		USE app;
		CREATE TABLE ref_tbl (id INT PRIMARY KEY);
		CREATE TABLE full_tbl (
			id INT NOT NULL,
			ref_z INT,
			ref_a INT,
			drop_col INT NULL DEFAULT 7,
			keep_null INT NULL,
			PRIMARY KEY (id),
			CONSTRAINT fk_zeta FOREIGN KEY (ref_z) REFERENCES ref_tbl (id),
			CONSTRAINT fk_alpha FOREIGN KEY (ref_a) REFERENCES ref_tbl (id),
			CONSTRAINT chk_zeta CHECK (id > 0) NOT ENFORCED,
			CONSTRAINT chk_alpha CHECK (id < 1000)
		) KEY_BLOCK_SIZE=8 COMPRESSION='lz4_1.0' STATS_PERSISTENT=0 STATS_AUTO_RECALC=1 STATS_SAMPLE_PAGES=32 MIN_ROWS=10 MAX_ROWS=100 AVG_ROW_LENGTH=64 PACK_KEYS=1 CHECKSUM=1 DELAY_KEY_WRITE=1;
		ALTER TABLE full_tbl ALTER COLUMN drop_col DROP DEFAULT;
		CREATE TABLE key_alg (
			id INT NOT NULL
		) PARTITION BY KEY ALGORITHM = 1 (id) PARTITIONS 2;
		CREATE TABLE range_sub (
			id INT NOT NULL,
			reg INT NOT NULL
		) PARTITION BY LIST (reg) SUBPARTITION BY RANGE (id) (
			PARTITION p_east VALUES IN (1) (
				SUBPARTITION sp0 VALUES LESS THAN (100),
				SUBPARTITION sp1 VALUES LESS THAN MAXVALUE
			)
		);
		CREATE TABLE list_sub (
			id INT NOT NULL,
			reg INT NOT NULL
		) PARTITION BY RANGE (id) SUBPARTITION BY LIST (reg) SUBPARTITION TEMPLATE (
			SUBPARTITION sp_a VALUES IN (1, 2),
			SUBPARTITION sp_b VALUES IN (3, 4)
		) (
			PARTITION p0 VALUES LESS THAN (100)
		);
		CREATE FUNCTION fn_coll() RETURNS INT DETERMINISTIC RETURN 1;
		CREATE PROCEDURE proc_coll() BEGIN SELECT 1; END;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	s := meta.Schemas[0]
	require.Equal(t, "utf8mb4_unicode_ci", s.Functions[0].DatabaseCollation)
	require.Equal(t, "utf8mb4_unicode_ci", s.Procedures[0].DatabaseCollation)

	// Tables: full_tbl, key_alg, list_sub, range_sub, ref_tbl
	full := s.Tables[0]
	require.Equal(t, "full_tbl", full.Name)
	require.Equal(t, []string{"fk_alpha", "fk_zeta"}, []string{full.ForeignKeys[0].Name, full.ForeignKeys[1].Name})
	require.Equal(t, []string{"chk_alpha", "chk_zeta"}, []string{full.CheckConstraints[0].Name, full.CheckConstraints[1].Name})
	require.False(t, full.CheckConstraints[0].NotEnforced)
	require.True(t, full.CheckConstraints[1].NotEnforced)
	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range full.Columns {
		cols[col.Name] = col
	}
	require.Empty(t, cols["drop_col"].Default)
	require.True(t, cols["drop_col"].DefaultDropped)
	require.Equal(t, "NULL", cols["keep_null"].Default)
	require.False(t, cols["keep_null"].DefaultDropped)

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	reloadedFull := reloaded.GetDatabase("app").GetTable("full_tbl")
	require.Equal(t, 8, reloadedFull.KeyBlockSize)
	require.Equal(t, "lz4_1.0", reloadedFull.Compression)
	require.Equal(t, "0", reloadedFull.StatsPersistent)
	require.Equal(t, "1", reloadedFull.StatsAutoRecalc)
	require.Equal(t, "32", reloadedFull.StatsSamplePages)
	require.Equal(t, "10", reloadedFull.MinRows)
	require.Equal(t, "100", reloadedFull.MaxRows)
	require.Equal(t, "64", reloadedFull.AvgRowLength)
	require.Equal(t, "1", reloadedFull.PackKeys)
	require.Equal(t, "1", reloadedFull.Checksum)
	require.Equal(t, "1", reloadedFull.DelayKeyWrite)
	require.True(t, reloadedFull.GetColumn("drop_col").DefaultDropped)
	require.False(t, reloadedFull.GetColumn("keep_null").DefaultDropped)

	reExported := reloaded.DatabaseMetadata("app").Schemas[0]
	require.Equal(t, int32(1), reExported.Tables[1].Partitions[0].Algorithm)
	require.Len(t, reExported.Tables[2].Partitions[0].Subpartitions, 2)
	require.Equal(t, "1,2", reExported.Tables[2].Partitions[0].Subpartitions[0].Value)
	require.Len(t, reExported.Tables[3].Partitions[0].Subpartitions, 2)
	require.Equal(t, "MAXVALUE", reExported.Tables[3].Partitions[0].Subpartitions[1].Value)
}

func TestDatabaseMetadataRoundTripColumnGroupsAndIndexSemantics(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE cg_tbl (
			id INT NOT NULL PRIMARY KEY,
			b INT,
			c INT,
			KEY idx_global (b) STORING (c) GLOBAL,
			KEY idx_local (b) LOCAL
		) PARTITION BY HASH (id) PARTITIONS 2
		  WITH COLUMN GROUP (ALL COLUMNS, EACH COLUMN);
		CREATE TABLE vec_tbl (
			id INT PRIMARY KEY,
			v VECTOR(3) NOT NULL,
			VECTOR INDEX idx_vec (v) WITH (distance=cosine, type=hnsw)
		);
		CREATE TABLE ft_tbl (
			id INT PRIMARY KEY,
			txt TEXT,
			FULLTEXT INDEX idx_ft (txt) WITH PARSER ik PARSER_PROPERTIES=(ik_mode='smart')
		);
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	s := meta.Schemas[0]
	// Tables sorted alphabetically: cg_tbl, ft_tbl, vec_tbl
	cg := s.Tables[0]
	require.Equal(t, "cg_tbl", cg.Name)
	require.Equal(t, "all columns, each column", cg.ColumnGroups)
	cgIdx := map[string]*metadata.IndexMetadata{}
	for _, idx := range cg.Indexes {
		cgIdx[idx.Name] = idx
	}
	require.Equal(t, "GLOBAL", cgIdx["idx_global"].Scope)
	require.True(t, cgIdx["idx_global"].IsGlobal)
	require.Equal(t, []string{"c"}, cgIdx["idx_global"].StoringColumns)
	require.Equal(t, "LOCAL", cgIdx["idx_local"].Scope)
	require.False(t, cgIdx["idx_local"].IsGlobal)

	ft := s.Tables[1]
	ftIdx := map[string]*metadata.IndexMetadata{}
	for _, idx := range ft.Indexes {
		ftIdx[idx.Name] = idx
	}
	require.Equal(t, "ik", ftIdx["idx_ft"].Parser)
	require.Equal(t, "smart", ftIdx["idx_ft"].ParserProperties["ik_mode"])

	vec := s.Tables[2]
	vecIdx := map[string]*metadata.IndexMetadata{}
	for _, idx := range vec.Indexes {
		vecIdx[idx.Name] = idx
	}
	require.True(t, vecIdx["idx_vec"].IsVector)
	require.Contains(t, vecIdx["idx_vec"].VectorParams, "DISTANCE=COSINE")

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	reCG := reloaded.GetDatabase("app").GetTable("cg_tbl")
	require.Equal(t, "all columns, each column", reCG.ColumnGroups.String())
	reMeta := reloaded.DatabaseMetadata("app").Schemas[0]
	reCGIdx := map[string]*metadata.IndexMetadata{}
	for _, idx := range reMeta.Tables[0].Indexes {
		reCGIdx[idx.Name] = idx
	}
	require.Equal(t, "GLOBAL", reCGIdx["idx_global"].Scope)
	require.Equal(t, []string{"c"}, reCGIdx["idx_global"].StoringColumns)
	require.Equal(t, "LOCAL", reCGIdx["idx_local"].Scope)
	require.Equal(t, "ik", reMeta.Tables[1].Indexes[1].Parser)
	require.Equal(t, "smart", reMeta.Tables[1].Indexes[1].ParserProperties["ik_mode"])
	require.True(t, reMeta.Tables[2].Indexes[1].IsVector)
	require.Contains(t, reMeta.Tables[2].Indexes[1].VectorParams, "DISTANCE=COSINE")
}

func TestDatabaseMetadataViewOptionsTablespaceAndCharsetExplicit(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
		USE app;

		CREATE TABLE t_cs (
			id INT PRIMARY KEY,
			bare_col VARCHAR(40),
			cs_col VARCHAR(40) CHARACTER SET latin1,
			col_col VARCHAR(40) COLLATE utf8mb4_bin
		) TABLESPACE ts_ob;

		CREATE VIEW v_local AS SELECT bare_col FROM t_cs WITH LOCAL CHECK OPTION;
		CREATE VIEW v_cascaded AS SELECT bare_col FROM t_cs WITH CASCADED CHECK OPTION;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	// Simulate non-default Algorithm/SqlSecurity on View struct if set.
	c.GetDatabase("app").Views["v_local"].Algorithm = "TEMPTABLE"
	c.GetDatabase("app").Views["v_local"].SqlSecurity = "INVOKER"

	meta := c.DatabaseMetadata("app")
	s := meta.Schemas[0]
	require.Equal(t, "ts_ob", s.Tables[0].Tablespace)
	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range s.Tables[0].Columns {
		cols[col.Name] = col
	}
	require.True(t, cols["bare_col"].GetCharsetInherited())
	require.True(t, cols["bare_col"].GetCollationInherited())
	require.False(t, cols["cs_col"].GetCharsetInherited())
	require.True(t, cols["cs_col"].GetCollationInherited())
	require.True(t, cols["col_col"].GetCharsetInherited())
	require.False(t, cols["col_col"].GetCollationInherited())
	require.False(t, cols["id"].GetCharsetInherited())
	require.False(t, cols["id"].GetCollationInherited())

	require.Equal(t, "CASCADED", s.Views[0].CheckOption)
	require.Equal(t, "LOCAL", s.Views[1].CheckOption)
	require.Equal(t, "TEMPTABLE", s.Views[1].Algorithm)
	require.Equal(t, "INVOKER", s.Views[1].SqlSecurity)

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	reTbl := reloaded.GetDatabase("app").GetTable("t_cs")
	require.Equal(t, "ts_ob", reTbl.Tablespace)
	require.False(t, reTbl.GetColumn("bare_col").CharsetExplicit)
	require.False(t, reTbl.GetColumn("bare_col").CollationExplicit)
	require.True(t, reTbl.GetColumn("cs_col").CharsetExplicit)
	require.False(t, reTbl.GetColumn("cs_col").CollationExplicit)
	require.False(t, reTbl.GetColumn("col_col").CharsetExplicit)
	require.True(t, reTbl.GetColumn("col_col").CollationExplicit)

	reLocal := reloaded.GetDatabase("app").Views["v_local"]
	require.Equal(t, "LOCAL", reLocal.CheckOption)
	require.Equal(t, "TEMPTABLE", reLocal.Algorithm)
	require.Equal(t, "INVOKER", reLocal.SqlSecurity)
	require.Equal(t, "CASCADED", reloaded.GetDatabase("app").Views["v_cascaded"].CheckOption)

	_, err = reloaded.Exec(`ALTER TABLE t_cs COLLATE = utf8mb4_bin;`, nil)
	require.NoError(t, err)
	n := NormalizerFor(MySQL80)
	_, bareColl := n.ResolveColumnCharsetCollation(reTbl, reTbl.GetColumn("bare_col"))
	_, csColl := n.ResolveColumnCharsetCollation(reTbl, reTbl.GetColumn("cs_col"))
	require.Equal(t, "utf8mb4_bin", bareColl)
	require.Equal(t, "latin1_swedish_ci", csColl)
}

func TestDatabaseMetadataSubpartitionTemplateCommentsAndMaterializedViewAttributes(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;

		CREATE TABLE tpl_tbl (
			id INT NOT NULL,
			reg INT NOT NULL
		) PARTITION BY RANGE (id) SUBPARTITION BY LIST (reg) SUBPARTITION TEMPLATE (
			SUBPARTITION sp_a VALUES IN (1, 2),
			SUBPARTITION sp_b VALUES IN (3, 4)
		) (
			PARTITION p0 VALUES LESS THAN (100)
		);
		ALTER TABLE tpl_tbl ADD PARTITION (PARTITION p1 VALUES LESS THAN (200) (SUBPARTITION p1_custom VALUES IN (1, 2, 3, 4)));

		CREATE TABLE base_t (id INT PRIMARY KEY, code VARCHAR(20));
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	// Simulate synced partition/subpartition comments and materialized view comment + indexes.
	tbl := c.GetDatabase("app").GetTable("tpl_tbl")
	tbl.Partitioning.Partitions[0].Comment = "part zero comment"
	tbl.Partitioning.Partitions[0].SubPartitions[0].Comment = "sub a comment"
	tbl.Partitioning.SubPartTemplate[0].Comment = "tpl a comment"

	baseMeta := c.DatabaseMetadata("app")
	baseMeta.Schemas[0].MaterializedViews = []*metadata.MaterializedViewMetadata{{
		Name:       "mv_summary",
		Definition: "select `app`.`base_t`.`id` AS `id`, `app`.`base_t`.`code` AS `code` from `app`.`base_t`",
		Comment:    "mv container comment",
		Indexes: []*metadata.IndexMetadata{
			{Name: "PRIMARY", Type: "BTREE", Unique: true, Primary: true, Visible: true, Expressions: []string{"id"}, KeyLength: []int64{-1}, Descending: []bool{false}},
			{Name: "idx_code", Type: "BTREE", Unique: false, Visible: true, Comment: "sec idx", Expressions: []string{"code"}, KeyLength: []int64{10}, Descending: []bool{true}},
		},
	}}

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), baseMeta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	reExported := reloaded.DatabaseMetadata("app")
	s := reExported.Schemas[0]

	// Partition & subpartition comments + subpartition template preserved across round-trip.
	reTbl := reloaded.GetDatabase("app").GetTable("tpl_tbl")
	require.Equal(t, "part zero comment", reTbl.Partitioning.Partitions[0].Comment)
	require.Equal(t, "sub a comment", reTbl.Partitioning.Partitions[0].SubPartitions[0].Comment)
	require.Len(t, reTbl.Partitioning.SubPartTemplate, 2)
	require.Equal(t, "tpl a comment", reTbl.Partitioning.SubPartTemplate[0].Comment)

	// Subsequent ADD PARTITION without subpartitions still applies the restored template.
	_, err = reloaded.Exec(`ALTER TABLE tpl_tbl ADD PARTITION (PARTITION p2 VALUES LESS THAN (300));`, nil)
	require.NoError(t, err)
	require.Len(t, reTbl.Partitioning.Partitions, 3)
	require.Len(t, reTbl.Partitioning.Partitions[2].SubPartitions, 2)
	require.Equal(t, "p2ssp_a", reTbl.Partitioning.Partitions[2].SubPartitions[0].Name)
	require.Equal(t, "tpl a comment", reTbl.Partitioning.Partitions[2].SubPartitions[0].Comment)

	// Materialized view comment and indexes preserved across LoadMetadata -> DatabaseMetadata.
	require.Len(t, s.MaterializedViews, 1)
	mvMeta := s.MaterializedViews[0]
	require.Equal(t, "mv_summary", mvMeta.Name)
	require.Equal(t, "mv container comment", mvMeta.Comment)
	require.Len(t, mvMeta.Indexes, 2)
	require.Equal(t, "PRIMARY", mvMeta.Indexes[0].Name)
	require.Equal(t, "idx_code", mvMeta.Indexes[1].Name)
	require.Equal(t, "sec idx", mvMeta.Indexes[1].Comment)
	require.Equal(t, []int64{10}, mvMeta.Indexes[1].KeyLength)
	require.Equal(t, []bool{true}, mvMeta.Indexes[1].Descending)

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

	// Multiple invalid functional indexes on a materialized view must all be skipped.
	badMVCatalog := New()
	badMeta := &metadata.DatabaseSchemaMetadata{
		Name: "app",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:    "t",
				Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}},
			}},
			MaterializedViews: []*metadata.MaterializedViewMetadata{{
				Name:       "mv_bad",
				Definition: "select `app`.`t`.`id` AS `id` from `app`.`t`",
				Indexes: []*metadata.IndexMetadata{
					{Name: "idx_good", Type: "BTREE", Visible: true, Expressions: []string{"id"}, KeyLength: []int64{-1}, Descending: []bool{false}},
					{Name: "idx_bad_1", Type: "BTREE", Visible: true, Expressions: []string{"(id + )"}, KeyLength: []int64{-1}, Descending: []bool{false}},
					{Name: "idx_bad_2", Type: "BTREE", Visible: true, Expressions: []string{"(id * )"}, KeyLength: []int64{-1}, Descending: []bool{false}},
				},
			}},
		}},
	}
	badReport, err := badMVCatalog.LoadMetadata(context.Background(), badMeta)
	require.NoError(t, err)
	require.Contains(t, badReport.Incomplete, "view:mv_bad")
	reMV := badMVCatalog.GetDatabase("app").Views["mv_bad"]
	require.NotNil(t, reMV)
	require.Len(t, reMV.Indexes, 1)
	require.Equal(t, "idx_good", reMV.Indexes[0].Name)

	// DatabaseDefinition renders vector/global indexes, tablespace, view options, and escaped comments.
	defCatalog := New()
	_, err = defCatalog.Exec(`
		CREATE DATABASE def_ob;
		USE def_ob;
		CREATE TABLE ob_tbl (
			id INT NOT NULL,
			k INT,
			s INT,
			txt VARCHAR(100),
			explicit_c VARCHAR(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			inherited_c VARCHAR(20),
			d_col INT DEFAULT 1,
			v VECTOR(3),
			note VARCHAR(50) COMMENT 'it''s a \\ path',
			PRIMARY KEY (id),
			KEY k_global (k) GLOBAL STORING (s) KEY_BLOCK_SIZE = 4,
			FULLTEXT KEY ft_txt (txt) WITH PARSER ik PARSER_PROPERTIES=(ik_mode='smart'),
			VECTOR KEY v_cos (v) WITH (TYPE=HNSW, DISTANCE=COSINE, LIB=VSAG)
		) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci TABLESPACE ts_ob COMMENT='O''Reilly' PARTITION BY HASH (id) PARTITIONS 2 WITH COLUMN GROUP(all columns);
		ALTER TABLE ob_tbl ALTER COLUMN d_col DROP DEFAULT;
		CREATE MATERIALIZED VIEW LOG ON ob_tbl;
		CREATE MATERIALIZED VIEW mv_pk (PRIMARY KEY (id)) AS SELECT id, note FROM ob_tbl;
		CREATE VIEW v_ob AS SELECT id, note FROM ob_tbl WITH LOCAL CHECK OPTION;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	defCatalog.GetDatabase("def_ob").Views["v_ob"].Algorithm = "MERGE"
	defCatalog.GetDatabase("def_ob").Views["v_ob"].SqlSecurity = "INVOKER"

	obDefMeta := defCatalog.DatabaseMetadata("def_ob")
	require.True(t, obDefMeta.Schemas[0].Tables[0].HasMaterializedViewLog)

	reloadedOB := New()
	report, err = reloadedOB.LoadMetadata(context.Background(), obDefMeta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)
	require.True(t, reloadedOB.GetDatabase("def_ob").MaterializedViewLogs["ob_tbl"])
	require.True(t, reloadedOB.GetDatabase("def_ob").MaterializedViewLogsKnown["ob_tbl"])

	obDef, err := DatabaseDefinition(obDefMeta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, obDef, "/*!50100 TABLESPACE `ts_ob` */")
	require.Contains(t, obDef, `COMMENT='O\'Reilly'`)
	require.Contains(t, obDef, `COMMENT 'it\'s a \\ path'`)
	require.Contains(t, obDef, "VECTOR KEY `v_cos` (`v`) WITH (TYPE=HNSW, DISTANCE=COSINE, LIB=VSAG")
	require.Contains(t, obDef, "KEY `k_global` (`k`) STORING (`s`) KEY_BLOCK_SIZE=4 GLOBAL")
	require.Contains(t, obDef, "FULLTEXT KEY `ft_txt` (`txt`) WITH PARSER ik PARSER_PROPERTIES=(ik_mode='smart')")
	require.Contains(t, obDef, "CREATE MATERIALIZED VIEW LOG ON `ob_tbl`;")
	require.Contains(t, obDef, "CREATE MATERIALIZED VIEW `mv_pk` (PRIMARY KEY (`id`))")
	require.Contains(t, obDef, "WITH COLUMN GROUP(all columns)")
	require.Contains(t, obDef, "`explicit_c` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")
	require.Contains(t, obDef, "`inherited_c` varchar(20) DEFAULT NULL")
	require.Contains(t, obDef, "ALTER TABLE `ob_tbl` ALTER COLUMN `d_col` DROP DEFAULT;")
	require.Contains(t, obDef, "CREATE OR REPLACE ALGORITHM=MERGE SQL SECURITY INVOKER VIEW `v_ob`")
	require.Contains(t, obDef, "WITH LOCAL CHECK OPTION")

	restoredOB := New()
	_, err = restoredOB.Exec("CREATE DATABASE def_ob;\nUSE def_ob;\n"+obDef, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	require.True(t, restoredOB.GetDatabase("def_ob").MaterializedViewLogs["ob_tbl"])
	reOBTbl := restoredOB.GetDatabase("def_ob").GetTable("ob_tbl")
	require.True(t, reOBTbl.GetColumn("d_col").DefaultDropped)
	require.False(t, reOBTbl.ColumnGroups.IsEmpty())
	require.True(t, reOBTbl.GetColumn("explicit_c").CharsetExplicit)
	require.False(t, reOBTbl.GetColumn("inherited_c").CharsetExplicit)
	reMVPK := restoredOB.GetDatabase("def_ob").Views["mv_pk"]
	require.NotNil(t, reMVPK)
	require.Len(t, reMVPK.Indexes, 1)
	require.True(t, reMVPK.Indexes[0].Primary)

	// Partition fallback still restores materialized-view log state on success.
	fallbackCatalog := New()
	fallbackMeta := &metadata.DatabaseSchemaMetadata{
		Name: "fb_db",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:                   "fb_tbl",
				HasMaterializedViewLog: true,
				Columns:                []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: false}},
				Partitions: []*metadata.TablePartitionMetadata{
					{Type: metadata.TablePartitionMetadata_RANGE, Expression: "missing_col", Name: "p0", Value: "10"},
				},
			}},
		}},
	}
	fbReport, err := fallbackCatalog.LoadMetadata(context.Background(), fallbackMeta)
	require.NoError(t, err)
	require.Contains(t, fbReport.Incomplete, "table:fb_tbl")
	require.True(t, fallbackCatalog.GetDatabase("fb_db").MaterializedViewLogs["fb_tbl"])
	require.True(t, fallbackCatalog.GetDatabase("fb_db").MaterializedViewLogsKnown["fb_tbl"])

	// Subpartition template and partition comments render cleanly in DatabaseDefinition.
	tplCatalog := New()
	_, err = tplCatalog.Exec(`
		CREATE DATABASE tpl_db;
		USE tpl_db;
		CREATE TABLE tpl_tbl (
			id INT NOT NULL,
			k INT NOT NULL
		) PARTITION BY RANGE (id)
		SUBPARTITION BY HASH (k)
		SUBPARTITION TEMPLATE (
			SUBPARTITION sp0,
			SUBPARTITION sp1
		) (
			PARTITION p0 VALUES LESS THAN (10),
			PARTITION p1 VALUES LESS THAN MAXVALUE
		);
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	tplMeta := tplCatalog.DatabaseMetadata("tpl_db")
	tplMeta.Schemas[0].Tables[0].Partitions[0].Comment = "p0's cmt"
	tplMeta.Schemas[0].Tables[0].Partitions[0].SubpartitionTemplate[0].Comment = "sp0 cmt"

	tplDef, err := DatabaseDefinition(tplMeta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, tplDef, "SUBPARTITION TEMPLATE (SUBPARTITION `sp0` ENGINE=InnoDB COMMENT = 'sp0 cmt'")
	require.Contains(t, tplDef, "SUBPARTITION `sp1` ENGINE=InnoDB)")
	require.Contains(t, tplDef, "PARTITION `p0` VALUES LESS THAN (10) ENGINE=InnoDB COMMENT = 'p0''s cmt'")

	// LoadMetadata strips comments during probe parse and keeps SubPartTemplate and comments.
	reloadedTpl := New()
	tplReport, err := reloadedTpl.LoadMetadata(context.Background(), tplMeta)
	require.NoError(t, err)
	require.Empty(t, tplReport.Degraded)
	require.Empty(t, tplReport.Missing)
	require.Empty(t, tplReport.Incomplete)
	reTplPart := reloadedTpl.GetDatabase("tpl_db").GetTable("tpl_tbl").Partitioning
	require.Len(t, reTplPart.SubPartTemplate, 2)
	require.Equal(t, "p0's cmt", reTplPart.Partitions[0].Comment)
	require.Equal(t, "sp0 cmt", reTplPart.SubPartTemplate[0].Comment)
}

func TestDatabaseMetadataQuotedHexDefaultRoundTrip(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE hex_t (
			id INT PRIMARY KEY,
			payload VARBINARY(4) DEFAULT 0x05,
			lower_hex BINARY(1)
		);
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	rawHex := "X'05'"
	rawLowerHex := "x'0a'"
	c.GetDatabase("app").GetTable("hex_t").GetColumn("payload").Default = &rawHex
	c.GetDatabase("app").GetTable("hex_t").GetColumn("lower_hex").Default = &rawLowerHex

	meta := c.DatabaseMetadata("app")
	tbl := meta.Schemas[0].Tables[0]
	require.Equal(t, "X'05'", tbl.Columns[1].Default)
	require.Equal(t, "x'0a'", tbl.Columns[2].Default)

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "DEFAULT X'05'")
	require.Contains(t, def, "DEFAULT x'0a'")

	execd := New()
	_, err = execd.Exec("CREATE DATABASE app;\nUSE app;\n"+def, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	require.NotNil(t, execd.GetDatabase("app").GetTable("hex_t"))
}

func TestDatabaseMetadataBacktickRoutineEventTriggerIdentifiers(t *testing.T) {
	c := New()
	_, err := c.Exec(""+
		"CREATE DATABASE app; USE app;\n"+
		"CREATE TABLE `t``1` (id INT PRIMARY KEY);\n"+
		"CREATE FUNCTION `fn``1`(`p``a` INT) RETURNS INT DETERMINISTIC RETURN 1;\n"+
		"CREATE PROCEDURE `sp``1`(IN `in``x` INT) SELECT 1;\n"+
		"CREATE EVENT `ev``1` ON SCHEDULE EVERY 1 HOUR DO SELECT 1;\n"+
		"CREATE TRIGGER `tr``1` BEFORE INSERT ON `t``1` FOR EACH ROW SET NEW.id = NEW.id;\n",
		&ExecOptions{ContinueOnError: false},
	)
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
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

	restored := New()
	_, err = restored.Exec("CREATE DATABASE app;\nUSE app;\n"+def, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	reDB := restored.GetDatabase("app")
	require.NotNil(t, reDB.Functions["fn`1"])
	require.NotNil(t, reDB.Procedures["sp`1"])
	require.NotNil(t, reDB.Events["ev`1"])
	require.NotNil(t, reDB.Triggers["tr`1"])
}
