package catalog

import (
	"context"
	"strings"
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
	require.NotNil(t, obDefMeta.Schemas[0].Tables[0].HasMaterializedViewLog)
	require.True(t, obDefMeta.Schemas[0].Tables[0].GetHasMaterializedViewLog())

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
				HasMaterializedViewLog: boolPtr(true),
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

func TestDatabaseMetadataViewColumnTypes(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE vdb;
		USE vdb;
		CREATE TABLE vt (id INT PRIMARY KEY, name VARCHAR(20));
		CREATE VIEW v_types AS
			SELECT id, name, id + 1 AS next_id, COUNT(*) AS cnt, id > 0 AS cmp,
				1 AS lit_int, 2147483648 AS lit_big, 'x' AS lit_str, '' AS lit_empty,
				NULL AS lit_null, TRUE AS lit_bool, 1.5 AS lit_dec
			FROM vt GROUP BY id, name;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	meta := c.DatabaseMetadata("vdb")
	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range meta.Schemas[0].Views[0].Columns {
		cols[col.Name] = col
	}
	require.Equal(t, "int(11)", cols["id"].Type)
	require.Equal(t, "varchar(20)", cols["name"].Type)
	require.Equal(t, "bigint", cols["next_id"].Type)
	require.Equal(t, "bigint", cols["cnt"].Type)
	require.Equal(t, "int", cols["cmp"].Type)
	require.Equal(t, "int", cols["lit_int"].Type)
	require.Equal(t, "bigint", cols["lit_big"].Type)
	require.Equal(t, "varchar(1)", cols["lit_str"].Type)
	require.Equal(t, "char(0)", cols["lit_empty"].Type)
	require.Equal(t, "binary(0)", cols["lit_null"].Type)
	require.Equal(t, "int", cols["lit_bool"].Type)
	require.Equal(t, "decimal(2,1)", cols["lit_dec"].Type)
}

func TestDatabaseMetadataDecimalArithmeticMVLogAbsenceAndSkipDump(t *testing.T) {
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
			CAST(1 AS DECIMAL(5,2)) - CAST(2 AS DECIMAL(3,1)) AS cm
		FROM dec_t;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	meta := c.DatabaseMetadata("ddb")
	vcols := map[string]*metadata.ColumnMetadata{}
	for _, col := range meta.Schemas[0].Views[0].Columns {
		vcols[col.Name] = col
	}
	require.Equal(t, "decimal(6,2)", vcols["p"].Type)
	require.Equal(t, "decimal(8,3)", vcols["mul"].Type)
	require.Equal(t, "decimal(10,6)", vcols["dv"].Type)
	require.Equal(t, "decimal(5,2)", vcols["md"].Type)
	require.Equal(t, "decimal(6,2)", vcols["cm"].Type)

	// LoadMetadata records known absence of MV log when HasMaterializedViewLog is false, and preserves SkipDump.
	loadMeta := &metadata.DatabaseSchemaMetadata{
		Name: "mv_skip_db",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:                   "no_mvlog_t",
				HasMaterializedViewLog: boolPtr(false),
				SkipDump:               true,
				Columns:                []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: false}},
				Triggers: []*metadata.TriggerMetadata{{
					Name: "tr_skip", Event: "INSERT", Timing: "BEFORE", Body: "BEGIN END", SkipDump: true,
				}},
			}, {
				Name:    "unknown_mvlog_t",
				Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: false}},
			}},
			Views: []*metadata.ViewMetadata{
				{Name: "v_skip", Definition: "SELECT 1 AS c", SkipDump: true},
				{Name: "v_uses_mv", Definition: "SELECT `id` FROM `mv_skip`"},
			},
			MaterializedViews: []*metadata.MaterializedViewMetadata{{
				Name: "mv_skip", Definition: "SELECT `id` FROM `no_mvlog_t`", SkipDump: true,
			}},
			Functions: []*metadata.FunctionMetadata{{
				Name: "fn_skip", Definition: "CREATE FUNCTION fn_skip() RETURNS INT DETERMINISTIC RETURN 1", SkipDump: true,
			}},
			Procedures: []*metadata.ProcedureMetadata{{
				Name: "sp_skip", Definition: "CREATE PROCEDURE sp_skip() BEGIN END", SkipDump: true,
			}},
		}},
	}
	loaded := New()
	report, err := loaded.LoadMetadata(context.Background(), loadMeta)
	require.NoError(t, err)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Degraded)
	db := loaded.GetDatabase("mv_skip_db")
	require.True(t, db.MaterializedViewLogsKnown["no_mvlog_t"])
	require.False(t, db.MaterializedViewLogs["no_mvlog_t"])
	require.False(t, db.MaterializedViewLogsKnown["unknown_mvlog_t"])
	require.Nil(t, loaded.DatabaseMetadata("mv_skip_db").Schemas[0].Tables[1].HasMaterializedViewLog)
	dropRes, err := loaded.Exec("DROP MATERIALIZED VIEW LOG ON no_mvlog_t;", &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	require.Len(t, dropRes, 1)
	require.Error(t, dropRes[0].Error)

	out := loaded.DatabaseMetadata("mv_skip_db")
	require.True(t, out.Schemas[0].Tables[0].SkipDump)
	require.True(t, out.Schemas[0].Tables[0].Triggers[0].SkipDump)
	require.True(t, out.Schemas[0].Views[0].SkipDump)
	require.Len(t, out.Schemas[0].MaterializedViews, 1)
	require.True(t, out.Schemas[0].MaterializedViews[0].SkipDump)
	require.True(t, out.Schemas[0].Functions[0].SkipDump)
	require.True(t, out.Schemas[0].Procedures[0].SkipDump)
}

func TestDatabaseMetadataStarViewColumnsAndImplicitFKDefaults(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE ref_db;
		USE ref_db;
		CREATE TABLE ref_t (
			id INT NOT NULL PRIMARY KEY,
			code VARCHAR(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
			note VARCHAR(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
		);

		CREATE DATABASE app;
		USE app;
		CREATE TABLE parent (id INT PRIMARY KEY);
		CREATE TABLE child (
			id INT PRIMARY KEY,
			pid INT,
			FOREIGN KEY (pid) REFERENCES parent(id)
		);
		CREATE TABLE local_t (
			id INT NOT NULL PRIMARY KEY,
			code VARCHAR(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
			note VARCHAR(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
		);
		CREATE VIEW v_star AS SELECT * FROM local_t;
		CREATE VIEW v_qual_star AS SELECT local_t.* FROM local_t;
		CREATE VIEW v_cross_star AS SELECT * FROM ref_db.ref_t;
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
	require.Equal(t, "RESTRICT", childProto.ForeignKeys[0].OnDelete)
	require.Equal(t, "RESTRICT", childProto.ForeignKeys[0].OnUpdate)

	viewByName := map[string]*metadata.ViewMetadata{}
	for _, v := range meta.Schemas[0].Views {
		viewByName[v.Name] = v
	}
	for _, vName := range []string{"v_star", "v_qual_star", "v_cross_star"} {
		v := viewByName[vName]
		require.NotNil(t, v, vName)
		require.Len(t, v.Columns, 3, vName)
		require.Equal(t, "int(11)", v.Columns[0].Type, vName)
		require.False(t, v.Columns[0].Nullable, vName)
		require.Empty(t, v.Columns[0].CharacterSet, vName)
		require.Equal(t, "varchar(16)", v.Columns[1].Type, vName)
		require.False(t, v.Columns[1].Nullable, vName)
		require.Equal(t, "utf8mb4", v.Columns[1].CharacterSet, vName)
		require.Equal(t, "utf8mb4_unicode_ci", v.Columns[1].Collation, vName)
		require.Equal(t, "varchar(32)", v.Columns[2].Type, vName)
		require.True(t, v.Columns[2].Nullable, vName)
		require.Equal(t, "utf8mb4", v.Columns[2].CharacterSet, vName)
		require.Equal(t, "utf8mb4_unicode_ci", v.Columns[2].Collation, vName)
	}

	// Functional index expressions export parenthesized so DatabaseDefinition restores them.
	_, err = c.Exec(`
		CREATE TABLE fn_idx_t (
			id INT NOT NULL PRIMARY KEY,
			code VARCHAR(20),
			INDEX fn_upper ((UPPER(code)))
		);
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	fnMeta := c.DatabaseMetadata("app")
	var fnTblMeta *metadata.TableMetadata
	for _, tbl := range fnMeta.Schemas[0].Tables {
		if tbl.Name == "fn_idx_t" {
			fnTblMeta = tbl
		}
	}
	require.NotNil(t, fnTblMeta)
	fnDef, err := DatabaseDefinition(&metadata.DatabaseSchemaMetadata{
		Name:    "app",
		Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{fnTblMeta}}},
	}, DefinitionOptions{})
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

	// Set-operation views derive metadata across all UNION/INTERSECT/EXCEPT arms.
	unionCat := New()
	_, err = unionCat.Exec(`
		CREATE DATABASE app DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
		USE app;
		CREATE TABLE t1 (
			id INT NOT NULL,
			s CHAR(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
			c_bin CHAR(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL
		);
		CREATE TABLE t2 (
			id BIGINT NULL,
			s VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL
		);
		CREATE VIEW v_mix AS SELECT 1 AS c UNION SELECT 'x';
		CREATE VIEW v_tbl AS SELECT id, s FROM t1 UNION SELECT id, s FROM t2;
		CREATE VIEW v_null AS SELECT id FROM t1 UNION SELECT NULL;
		CREATE VIEW v_num AS SELECT 1 AS a UNION SELECT 1.5 AS a;
		CREATE VIEW v_col AS SELECT s FROM t1 UNION SELECT c_bin FROM t1;
		CREATE VIEW v_except AS SELECT id, s FROM t1 EXCEPT SELECT id, s FROM t2;
		CREATE VIEW v_intersect AS SELECT id, s FROM t1 INTERSECT SELECT id, s FROM t2;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	uViews := map[string]*metadata.ViewMetadata{}
	for _, v := range unionCat.DatabaseMetadata("app").Schemas[0].Views {
		uViews[v.Name] = v
	}
	require.Equal(t, "varchar(2)", uViews["v_mix"].Columns[0].Type)
	require.False(t, uViews["v_mix"].Columns[0].Nullable)
	require.Equal(t, "utf8mb4", uViews["v_mix"].Columns[0].CharacterSet)

	require.Equal(t, "bigint", uViews["v_tbl"].Columns[0].Type)
	require.True(t, uViews["v_tbl"].Columns[0].Nullable)
	require.Equal(t, "varchar(10)", uViews["v_tbl"].Columns[1].Type)
	require.True(t, uViews["v_tbl"].Columns[1].Nullable)
	require.Equal(t, "utf8mb4_unicode_ci", uViews["v_tbl"].Columns[1].Collation)

	require.Equal(t, "int(11)", uViews["v_null"].Columns[0].Type)
	require.True(t, uViews["v_null"].Columns[0].Nullable)
	require.Empty(t, uViews["v_null"].Columns[0].CharacterSet)

	require.Equal(t, "decimal(2,1)", uViews["v_num"].Columns[0].Type)
	require.Equal(t, "utf8mb4_bin", uViews["v_col"].Columns[0].Collation)

	require.Equal(t, "bigint", uViews["v_except"].Columns[0].Type)
	require.False(t, uViews["v_except"].Columns[0].Nullable)
	require.Equal(t, "varchar(10)", uViews["v_except"].Columns[1].Type)
	require.False(t, uViews["v_except"].Columns[1].Nullable)

	require.Equal(t, "bigint", uViews["v_intersect"].Columns[0].Type)
	require.False(t, uViews["v_intersect"].Columns[0].Nullable)

	// MIN/MAX view columns derive their types from their arguments, and unflagged
	// LoadMetadata columns matching table defaults remain inherited.
	mmCat := New()
	_, err = mmCat.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE t_mm (id INT NOT NULL, s VARCHAR(20) NOT NULL, d DATE NOT NULL, dec1 DECIMAL(10,2) NOT NULL);
		CREATE VIEW v_mm AS SELECT MIN(id) AS min_id, MAX(s) AS max_s, MIN(d) AS min_d, MAX(dec1) AS max_dec FROM t_mm;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	var vMM *metadata.ViewMetadata
	for _, v := range mmCat.DatabaseMetadata("app").Schemas[0].Views {
		if v.Name == "v_mm" {
			vMM = v
		}
	}
	require.NotNil(t, vMM)
	require.Equal(t, "int(11)", vMM.Columns[0].Type)
	require.Equal(t, "varchar(20)", vMM.Columns[1].Type)
	require.Equal(t, "date", vMM.Columns[2].Type)
	require.Equal(t, "decimal(10,2)", vMM.Columns[3].Type)

	unflaggedCat := New()
	_, err = unflaggedCat.LoadMetadata(context.Background(), &metadata.DatabaseSchemaMetadata{
		Name:         "app",
		CharacterSet: "utf8mb4",
		Collation:    "utf8mb4_unicode_ci",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:      "t_unflagged",
				Engine:    "oceanbase",
				Charset:   "utf8mb4",
				Collation: "utf8mb4_unicode_ci",
				Columns: []*metadata.ColumnMetadata{
					{Name: "inh", Position: 1, Type: "varchar(20)", Nullable: true, CharacterSet: "utf8mb4", Collation: "utf8mb4_unicode_ci"},
					{Name: "exp", Position: 2, Type: "varchar(20)", Nullable: true, CharacterSet: "latin1", Collation: "latin1_swedish_ci"},
				},
			}},
		}},
	})
	require.NoError(t, err)
	unTbl := unflaggedCat.GetDatabase("app").GetTable("t_unflagged")
	require.False(t, unTbl.GetColumn("inh").CharsetExplicit)
	require.False(t, unTbl.GetColumn("inh").CollationExplicit)
	require.True(t, unTbl.GetColumn("exp").CharsetExplicit)
	require.True(t, unTbl.GetColumn("exp").CollationExplicit)

	// Implicit subpartitions preserve UseDefault="default" without materializing p0sp0,
	// and unresolved view columns retain supplied snapshot metadata across LoadMetadata.
	subCat := New()
	_, err = subCat.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE sub_default (
			id INT NOT NULL,
			k INT NOT NULL
		) PARTITION BY RANGE (id) SUBPARTITION BY HASH (k) (PARTITION p0 VALUES LESS THAN (10));
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	subMeta := subCat.DatabaseMetadata("app")
	require.Equal(t, "default", subMeta.Schemas[0].Tables[0].Partitions[0].Subpartitions[0].UseDefault)
	subDef, err := DatabaseDefinition(subMeta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, subDef, "SUBPARTITION BY HASH (`k`)")
	require.NotContains(t, subDef, "SUBPARTITION `p0sp0`")

	reSubCat := New()
	_, err = reSubCat.LoadMetadata(context.Background(), subMeta)
	require.NoError(t, err)
	reSubTbl := reSubCat.GetDatabase("app").GetTable("sub_default")
	require.True(t, reSubTbl.Partitioning.UseDefaultSubpartitions)
	require.True(t, reSubTbl.Partitioning.UseDefaultNumSubpartitions)

	extCat := New()
	_, err = extCat.LoadMetadata(context.Background(), &metadata.DatabaseSchemaMetadata{
		Name: "app",
		Schemas: []*metadata.SchemaMetadata{{
			Views: []*metadata.ViewMetadata{{
				Name:       "v_ext",
				Definition: "select `unloaded_db`.`ext_tbl`.`id` AS `id`,`unloaded_db`.`ext_tbl`.`name` AS `name` from `unloaded_db`.`ext_tbl`",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Position: 1, Type: "int(11)", Nullable: false},
					{Name: "name", Position: 2, Type: "varchar(32)", Nullable: true, CharacterSet: "utf8mb4", Collation: "utf8mb4_unicode_ci"},
				},
			}},
		}},
	})
	require.NoError(t, err)
	reExt := extCat.DatabaseMetadata("app").Schemas[0].Views[0]
	require.Len(t, reExt.Columns, 2)
	require.Equal(t, "int(11)", reExt.Columns[0].Type)
	require.False(t, reExt.Columns[0].Nullable)
	require.Equal(t, "varchar(32)", reExt.Columns[1].Type)
	require.True(t, reExt.Columns[1].Nullable)
	require.Equal(t, "utf8mb4_unicode_ci", reExt.Columns[1].Collation)

	// Aggregate types (SUM/AVG) and nullability, typed temporal defaults, and routine normalization.
	aggCat := New()
	_, err = aggCat.Exec(`
		CREATE DATABASE app DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
		USE app;
		CREATE TABLE t_agg (
			i INT NOT NULL,
			b BIGINT NOT NULL,
			f FLOAT NOT NULL,
			d DOUBLE NOT NULL,
			dec1 DECIMAL(10,2) NOT NULL,
			dt DATE DEFAULT NULL
		);
		ALTER TABLE t_agg ALTER COLUMN dt SET DEFAULT DATE '2020-01-01';
		CREATE VIEW v_agg AS SELECT
			SUM(i) AS s_i, SUM(b) AS s_b, SUM(f) AS s_f, SUM(d) AS s_d, SUM(dec1) AS s_dec,
			AVG(i) AS a_i, AVG(b) AS a_b, AVG(f) AS a_f, AVG(d) AS a_d, AVG(dec1) AS a_dec,
			MIN(i) AS mn_i, MAX(i) AS mx_i, COUNT(*) AS cnt, GROUP_CONCAT(i) AS gc_i
		FROM t_agg;
		CREATE FUNCTION fn_str(s VARCHAR(10)) RETURNS VARCHAR(10) CHARSET utf8mb4 DETERMINISTIC RETURN s;
		CREATE PROCEDURE sp_noop() BEGIN SELECT 1; END;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	aggMeta := aggCat.DatabaseMetadata("app")
	aggViewCols := map[string]*metadata.ColumnMetadata{}
	for _, col := range aggMeta.Schemas[0].Views[0].Columns {
		aggViewCols[col.Name] = col
	}
	require.Equal(t, "decimal(32,0)", aggViewCols["s_i"].Type)
	require.True(t, aggViewCols["s_i"].Nullable)
	require.Equal(t, "decimal(41,0)", aggViewCols["s_b"].Type)
	require.Equal(t, "double", aggViewCols["s_f"].Type)
	require.Equal(t, "double", aggViewCols["s_d"].Type)
	require.Equal(t, "decimal(32,2)", aggViewCols["s_dec"].Type)
	require.Equal(t, "decimal(14,4)", aggViewCols["a_i"].Type)
	require.True(t, aggViewCols["a_i"].Nullable)
	require.Equal(t, "decimal(23,4)", aggViewCols["a_b"].Type)
	require.Equal(t, "double", aggViewCols["a_f"].Type)
	require.Equal(t, "double", aggViewCols["a_d"].Type)
	require.Equal(t, "decimal(14,6)", aggViewCols["a_dec"].Type)
	require.True(t, aggViewCols["mn_i"].Nullable)
	require.True(t, aggViewCols["mx_i"].Nullable)
	require.True(t, aggViewCols["cnt"].Nullable)
	require.True(t, aggViewCols["gc_i"].Nullable)

	dtCol := aggMeta.Schemas[0].Tables[0].Columns[5]
	require.Equal(t, "'2020-01-01'", dtCol.Default)
	aggDef, err := DatabaseDefinition(aggMeta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, aggDef, "DEFAULT '2020-01-01'")

	fnDefSQL := aggMeta.Schemas[0].Functions[0].Definition
	require.True(t, strings.HasPrefix(fnDefSQL, "CREATE FUNCTION "))
	require.NotContains(t, fnDefSQL, "DEFINER=")
	require.NotContains(t, fnDefSQL, "CHARSET utf8mb4")
	spDefSQL := aggMeta.Schemas[0].Procedures[0].Definition
	require.True(t, strings.HasPrefix(spDefSQL, "CREATE PROCEDURE "))
	require.NotContains(t, spDefSQL, "DEFINER=")

	// CASE and control-flow functions (COALESCE, IFNULL, NULLIF, IF) derive
	// result column types and nullability.
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
			COALESCE(NULL, ni, id) AS coal_null_int,
			COALESCE(ni, NULL) AS coal_nullable,
			COALESCE(id, d) AS coal_int_dec,
			COALESCE(c, c6) AS coal_char,
			IFNULL(ti, ti) AS ifnull_ti,
			IFNULL(id, 0) AS ifnull_int,
			IFNULL(v, 'abc') AS ifnull_str,
			NULLIF(ti, 1) AS nullif_ti,
			NULLIF(id, 0) AS nullif_int,
			NULLIF(c, 'xy') AS nullif_char,
			NULLIF(d, 0) AS nullif_dec,
			IF(ni > 0, id, 0) AS if_int,
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
	require.False(t, ctrlCols["coal_int"].Nullable)
	require.Equal(t, "bigint", ctrlCols["coal_null_int"].Type)
	require.False(t, ctrlCols["coal_null_int"].Nullable)
	require.True(t, ctrlCols["coal_nullable"].Nullable)
	require.Equal(t, "decimal(12,2)", ctrlCols["coal_int_dec"].Type)
	require.Equal(t, "varchar(6)", ctrlCols["coal_char"].Type)
	require.Equal(t, "tinyint", ctrlCols["ifnull_ti"].Type)
	require.False(t, ctrlCols["ifnull_ti"].Nullable)
	require.Equal(t, "bigint", ctrlCols["ifnull_int"].Type)
	require.Equal(t, "varchar(20)", ctrlCols["ifnull_str"].Type)
	require.Equal(t, "int", ctrlCols["nullif_ti"].Type)
	require.True(t, ctrlCols["nullif_ti"].Nullable)
	require.Equal(t, "bigint", ctrlCols["nullif_int"].Type)
	require.True(t, ctrlCols["nullif_int"].Nullable)
	require.Equal(t, "varchar(4)", ctrlCols["nullif_char"].Type)
	require.Equal(t, "decimal(10,2)", ctrlCols["nullif_dec"].Type)
	require.Equal(t, "bigint", ctrlCols["if_int"].Type)
	require.False(t, ctrlCols["if_int"].Nullable)
	require.Equal(t, "decimal(11,1)", ctrlCols["if_int_dec"].Type)
	require.Equal(t, "bigint", ctrlCols["if_null_int"].Type)
	require.True(t, ctrlCols["if_null_int"].Nullable)
	require.Equal(t, "varchar(20)", ctrlCols["if_char_varchar"].Type)
	require.Equal(t, "datetime", ctrlCols["if_dt_ts"].Type)
	require.Equal(t, "bigint", ctrlCols["case_int"].Type)
	require.False(t, ctrlCols["case_int"].Nullable)
	require.Equal(t, "int", ctrlCols["case_ti_noelse"].Type)
	require.True(t, ctrlCols["case_ti_noelse"].Nullable)
	require.Equal(t, "varchar(3)", ctrlCols["case_str"].Type)
	require.Equal(t, "double", ctrlCols["case_f_d"].Type)

	// Scalar function view columns resolve types accurately.
	fnCat := New()
	_, err = fnCat.Exec(`
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
			SUBSTRING(v, 2) AS sub_v2,
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
	fnCols := map[string]*metadata.ColumnMetadata{}
	for _, col := range fnCat.DatabaseMetadata("app").Schemas[0].Views[0].Columns {
		fnCols[col.Name] = col
	}
	require.Equal(t, "longtext", fnCols["ju"].Type)
	require.Equal(t, "json", fnCols["je"].Type)
	require.Equal(t, "varchar(17)", fnCols["jt"].Type)
	require.Equal(t, "varchar(32)", fnCols["lw_v"].Type)
	require.Equal(t, "varchar(4)", fnCols["up_c"].Type)
	require.Equal(t, "mediumtext", fnCols["tr_txt"].Type)
	require.Equal(t, "varchar(32)", fnCols["rp_v"].Type)
	require.Equal(t, "varchar(5)", fnCols["sub_v"].Type)
	require.Equal(t, "varchar(31)", fnCols["sub_v2"].Type)
	require.Equal(t, "varchar(5)", fnCols["lf_v"].Type)
	require.Equal(t, "varchar(2)", fnCols["rg_c"].Type)
	require.Equal(t, "varchar(10)", fnCols["lp_v"].Type)
	require.Equal(t, "varchar(12)", fnCols["rep_c"].Type)
	require.Equal(t, "varchar(36)", fnCols["cc_cv"].Type)
	require.Equal(t, "varchar(37)", fnCols["ccw"].Type)
	require.Equal(t, "varchar(5)", fnCols["sp_5"].Type)
	require.Equal(t, "varchar(48)", fnCols["fmt_d"].Type)
	require.Equal(t, "bigint", fnCols["abs_i"].Type)
	require.False(t, fnCols["abs_i"].Nullable)
	require.Equal(t, "double", fnCols["abs_f"].Type)
	require.Equal(t, "double", fnCols["abs_dbl"].Type)
	require.Equal(t, "decimal(10,2)", fnCols["abs_d"].Type)
	require.Equal(t, "bigint", fnCols["ceil_d"].Type)
	require.Equal(t, "bigint", fnCols["fl_d"].Type)
	require.Equal(t, "decimal(10,1)", fnCols["rnd_d"].Type)
	require.Equal(t, "double", fnCols["rnd_dbl"].Type)
	require.Equal(t, "decimal(9,1)", fnCols["trn_d"].Type)
	require.Equal(t, "bigint", fnCols["mod_i"].Type)
	require.True(t, fnCols["mod_i"].Nullable)
	require.Equal(t, "decimal(10,2)", fnCols["mod_d"].Type)
	require.True(t, fnCols["mod_d"].Nullable)
	require.Equal(t, "bigint", fnCols["len_v"].Type)
	require.Equal(t, "bigint", fnCols["jv"].Type)
	require.Equal(t, "datetime(6)", fnCols["fn_now6"].Type)
	require.Equal(t, "datetime", fnCols["fn_now0"].Type)
	require.Equal(t, "datetime(3)", fnCols["fn_ct3"].Type)
	require.Equal(t, "datetime(6)", fnCols["fn_sys6"].Type)
	require.Equal(t, "datetime(3)", fnCols["fn_lts3"].Type)
	require.Equal(t, "datetime(6)", fnCols["fn_utcts6"].Type)
	require.Equal(t, "time(6)", fnCols["fn_curt6"].Type)
	require.Equal(t, "time(3)", fnCols["fn_curt3"].Type)
	require.Equal(t, "time(6)", fnCols["fn_time_dtm6"].Type)
	require.Equal(t, "datetime(4)", fnCols["fn_ts_dt_tm4"].Type)
}

func TestViewMetadataScalarSubqueryAndTemporalArithmetic(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE t (
			id INT NOT NULL PRIMARY KEY,
			ti TINYINT NOT NULL,
			ui INT UNSIGNED NOT NULL,
			ch CHAR(5) NOT NULL,
			d DATE NOT NULL,
			dt DATETIME NOT NULL,
			dt3 DATETIME(3) NOT NULL,
			ts TIMESTAMP(6) NOT NULL,
			tm TIME NOT NULL,
			tm2 TIME(2) NOT NULL,
			s VARCHAR(20) NOT NULL
		);
		CREATE VIEW v AS SELECT
			(SELECT id FROM t LIMIT 1) AS sq_id,
			(SELECT ti FROM t LIMIT 1) AS sq_ti,
			(SELECT ui FROM t LIMIT 1) AS sq_ui,
			(SELECT ch FROM t LIMIT 1) AS sq_ch,
			(SELECT COUNT(*) FROM t) AS sq_cnt,
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
			TIMESTAMPADD(HOUR, 1, d) AS ta_d_hr,
			TIMESTAMPADD(MICROSECOND, 1, dt) AS ta_dt_us,
			d + INTERVAL 1 DAY AS op_d_day,
			d + INTERVAL 1 HOUR AS op_d_hr,
			INTERVAL 1 DAY + d AS op_iv_d,
			d - INTERVAL 1 DAY AS op_d_sub,
			TIMEDIFF(dt3, dt) AS tdiff_dt,
			TIMEDIFF(tm2, tm) AS tdiff_tm,
			COALESCE((SELECT id FROM t LIMIT 1), 0) AS sq_coalesce,
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

	require.Equal(t, "bigint", cols["sq_id"].Type)
	require.True(t, cols["sq_id"].Nullable)
	require.Equal(t, "int", cols["sq_ti"].Type)
	require.True(t, cols["sq_ti"].Nullable)
	require.Equal(t, "bigint unsigned", cols["sq_ui"].Type)
	require.True(t, cols["sq_ui"].Nullable)
	require.Equal(t, "varchar(5)", cols["sq_ch"].Type)
	require.True(t, cols["sq_ch"].Nullable)
	require.Equal(t, "bigint", cols["sq_cnt"].Type)
	require.True(t, cols["sq_cnt"].Nullable)
	require.Equal(t, "int", cols["ex"].Type)
	require.False(t, cols["ex"].Nullable)
	require.Equal(t, "int", cols["in_col"].Type)
	require.True(t, cols["in_col"].Nullable)

	require.Equal(t, "bigint", cols["td_day"].Type)
	require.True(t, cols["td_day"].Nullable)
	require.Equal(t, "bigint", cols["td_us"].Type)
	require.Equal(t, "int", cols["dd"].Type)
	require.True(t, cols["dd"].Nullable)
	require.Equal(t, "date", cols["da_d_day"].Type)
	require.True(t, cols["da_d_day"].Nullable)
	require.Equal(t, "datetime", cols["da_d_hr"].Type)
	require.Equal(t, "datetime(6)", cols["da_d_us"].Type)
	require.Equal(t, "datetime(3)", cols["da_dt3_day"].Type)
	require.Equal(t, "datetime(6)", cols["da_ts_day"].Type)
	require.Equal(t, "time(6)", cols["da_tm2_us"].Type)
	require.Equal(t, "varchar(29)", cols["da_s_day"].Type)
	require.Equal(t, "date", cols["ds_d_day"].Type)
	require.Equal(t, "date", cols["ad_d_int"].Type)
	require.Equal(t, "date", cols["sd_d_int"].Type)
	require.Equal(t, "datetime", cols["ta_d_hr"].Type)
	require.Equal(t, "datetime(6)", cols["ta_dt_us"].Type)
	require.Equal(t, "date", cols["op_d_day"].Type)
	require.Equal(t, "datetime", cols["op_d_hr"].Type)
	require.Equal(t, "date", cols["op_iv_d"].Type)
	require.Equal(t, "date", cols["op_d_sub"].Type)
	require.Equal(t, "time(3)", cols["tdiff_dt"].Type)
	require.Equal(t, "time(2)", cols["tdiff_tm"].Type)
	require.Equal(t, "bigint", cols["sq_coalesce"].Type)
	require.Equal(t, "date", cols["sq_interval"].Type)
}

func TestViewPredicateNullability(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE p_t (id INT NOT NULL, nc INT NULL, nn INT NOT NULL);
		CREATE VIEW v_pred AS SELECT
			nc BETWEEN 1 AND 2 AS b_nc,
			nn BETWEEN 1 AND 2 AS b_nn,
			1 BETWEEN NULL AND 2 AS b_null,
			nc IN (1, 2) AS in_nc,
			nn IN (1, 2) AS in_nn,
			1 IN (NULL, 2) AS in_null,
			nc AND TRUE AS and_nc,
			nn AND TRUE AS and_nn,
			TRUE OR NULL AS or_null,
			NOT nc AS not_nc,
			NOT nn AS not_nn,
			nc IS NULL AS isnull_nc,
			nc IS TRUE AS istrue_nc,
			nc <=> 1 AS ns_nc
		FROM p_t;
	`, nil)
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	var vp *metadata.ViewMetadata
	for _, vm := range meta.GetSchemas()[0].GetViews() {
		if vm.Name == "v_pred" {
			vp = vm
		}
	}
	require.NotNil(t, vp)
	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range vp.GetColumns() {
		cols[col.Name] = col
	}
	require.True(t, cols["b_nc"].Nullable)
	require.False(t, cols["b_nn"].Nullable)
	require.True(t, cols["b_null"].Nullable)
	require.True(t, cols["in_nc"].Nullable)
	require.False(t, cols["in_nn"].Nullable)
	require.True(t, cols["in_null"].Nullable)
	require.True(t, cols["and_nc"].Nullable)
	require.False(t, cols["and_nn"].Nullable)
	require.True(t, cols["or_null"].Nullable)
	require.True(t, cols["not_nc"].Nullable)
	require.False(t, cols["not_nn"].Nullable)
	require.False(t, cols["isnull_nc"].Nullable)
	require.False(t, cols["istrue_nc"].Nullable)
	require.False(t, cols["ns_nc"].Nullable)
}

func TestDatabaseDefinitionPreservesTablespaceQuotesAndParenColumnIndex(t *testing.T) {
	tsMeta := &metadata.DatabaseSchemaMetadata{
		Name: "app",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:       "ts_tbl",
				Engine:     "InnoDB",
				Tablespace: "`ts_bt`",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Position: 1, Type: "int"},
					{Name: "a(b", Position: 2, Type: "int", Nullable: true},
				},
				Indexes: []*metadata.IndexMetadata{
					{Name: "idx_paren", Expressions: []string{"a(b"}, Type: "BTREE", Visible: true},
				},
			}, {
				Name:       "ts_only_quotes",
				Engine:     "InnoDB",
				Tablespace: "\"\"",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Position: 1, Type: "int"},
				},
			}},
		}},
	}
	def, err := DatabaseDefinition(tsMeta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "/*!50100 TABLESPACE ```ts_bt``` */")
	require.Contains(t, def, "/*!50100 TABLESPACE `\"\"` */")
	require.Contains(t, def, "KEY `idx_paren` (`a(b`)")

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), tsMeta)
	require.NoError(t, err)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Incomplete)

	reMeta := reloaded.DatabaseMetadata("app")
	var foundTS, foundEmptyQuoteTS, foundIdx bool
	for _, tbl := range reMeta.GetSchemas()[0].GetTables() {
		if tbl.Name == "ts_tbl" && tbl.Tablespace == "`ts_bt`" {
			foundTS = true
			for _, idx := range tbl.GetIndexes() {
				if idx.Name == "idx_paren" && len(idx.Expressions) == 1 && idx.Expressions[0] == "a(b" {
					foundIdx = true
				}
			}
		}
		if tbl.Name == "ts_only_quotes" && tbl.Tablespace == "\"\"" {
			foundEmptyQuoteTS = true
		}
	}
	require.True(t, foundTS)
	require.True(t, foundEmptyQuoteTS)
	require.True(t, foundIdx)
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

	require.Equal(t, "int", cols["l_nn"].Type)
	require.False(t, cols["l_nn"].Nullable)
	require.Equal(t, "int", cols["l_null"].Type)
	require.True(t, cols["l_null"].Nullable)
	require.Equal(t, "int", cols["nl_nn"].Type)
	require.False(t, cols["nl_nn"].Nullable)
	require.Equal(t, "int", cols["nl_null"].Type)
	require.True(t, cols["nl_null"].Nullable)
	require.Equal(t, "int", cols["l_pat_null"].Type)
	require.True(t, cols["l_pat_null"].Nullable)
	require.Equal(t, "int", cols["l_esc_nn"].Type)
	require.False(t, cols["l_esc_nn"].Nullable)
	require.Equal(t, "int", cols["l_esc_null"].Type)
	require.True(t, cols["l_esc_null"].Nullable)
}

func TestMaterializedViewParenColumnIndexRoundTrip(t *testing.T) {
	c := New()
	mvMeta := &metadata.DatabaseSchemaMetadata{
		Name: "app",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:    "src",
				Engine:  "InnoDB",
				Columns: []*metadata.ColumnMetadata{{Name: "id", Position: 1, Type: "int"}, {Name: "a(b", Position: 2, Type: "int", Nullable: true}},
				Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Visible: true, Expressions: []string{"id"}, KeyLength: []int64{-1}, Descending: []bool{false}}},
			}},
			MaterializedViews: []*metadata.MaterializedViewMetadata{{
				Name:       "mv_paren",
				Definition: "select `app`.`src`.`id` AS `id`,`app`.`src`.`a(b` AS `a(b` from `app`.`src`",
				Indexes: []*metadata.IndexMetadata{
					{Name: "PRIMARY", Primary: true, Unique: true, Visible: true, Type: "BTREE", Expressions: []string{"id"}, KeyLength: []int64{-1}, Descending: []bool{false}},
					{Name: "idx_paren", Visible: true, Type: "BTREE", Expressions: []string{"a(b"}, KeyLength: []int64{-1}, Descending: []bool{false}},
				},
			}},
		}},
	}

	report, err := c.LoadMetadata(context.Background(), mvMeta)
	require.NoError(t, err)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Incomplete)

	reMV := c.GetDatabase("app").Views["mv_paren"]
	require.NotNil(t, reMV)
	require.Len(t, reMV.Indexes, 2)
	require.Equal(t, "a(b", reMV.Indexes[1].Columns[0].Name)
	require.Empty(t, reMV.Indexes[1].Columns[0].Expr)

	reMeta := c.DatabaseMetadata("app")
	def, err := DatabaseDefinition(reMeta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "CREATE INDEX `idx_paren` ON `mv_paren` (`a(b`)")
}

func TestViewNonArithmeticOperatorMetadata(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20) NOT NULL, txt TEXT, doc JSON NOT NULL);
		CREATE VIEW v_ops AS SELECT
			id & 1 AS c_band,
			id XOR 1 AS c_lxor,
			doc->'$.x' AS c_jextract,
			doc->>'$.x' AS c_junquote,
			BINARY name AS c_bin_vc,
			BINARY txt AS c_bin_txt
		FROM t;
	`, nil)
	require.NoError(t, err)

	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range c.DatabaseMetadata("app").GetSchemas()[0].GetViews()[0].GetColumns() {
		cols[col.Name] = col
	}
	require.Equal(t, "bigint unsigned", cols["c_band"].Type)
	require.Equal(t, "int", cols["c_lxor"].Type)
	require.Equal(t, "json", cols["c_jextract"].Type)
	require.True(t, cols["c_jextract"].Nullable)
	require.Equal(t, "longtext", cols["c_junquote"].Type)
	require.True(t, cols["c_junquote"].Nullable)
	require.Equal(t, "varbinary(20)", cols["c_bin_vc"].Type)
	require.Equal(t, "varbinary(65535)", cols["c_bin_txt"].Type)
}

// TestDatabaseMetadataCollateViewColumns pins a COLLATE view column as
// oceanbase-ce 4.3.5.6 reports it: its operand's type, the collation named, and nullable
// whatever the operand (IS_NULLABLE YES over a NOT NULL column).
func TestDatabaseMetadataCollateViewColumns(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE vt (name VARCHAR(32), nn VARCHAR(8) NOT NULL);
		CREATE VIEW vc AS SELECT name COLLATE utf8mb4_bin AS c, nn COLLATE utf8mb4_bin AS d FROM vt;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	cols := c.DatabaseMetadata("app").Schemas[0].Views[0].Columns
	require.Len(t, cols, 2)
	require.Equal(t, "varchar(32)", cols[0].Type)
	require.True(t, cols[0].Nullable)
	require.Equal(t, "utf8mb4_bin", cols[0].Collation)
	require.Equal(t, "varchar(8)", cols[1].Type)
	require.True(t, cols[1].Nullable)
	require.Equal(t, "utf8mb4_bin", cols[1].Collation)
}
