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
			hex_col BINARY(1) DEFAULT 0x05,
			nullable_col INT,
			ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			gen_stored VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_s')) STORED,
			gen_virt VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_v')) VIRTUAL,
			g0 GEOMETRY NOT NULL SRID 0,
			g4326 POINT NOT NULL SRID 4326,
			secret INT INVISIBLE NOT NULL DEFAULT 0,
			SPATIAL INDEX sp_idx (g4326),
			INDEX fn_idx ((UPPER(code)))
		) ENGINE=InnoDB ROW_FORMAT=DYNAMIC KEY_BLOCK_SIZE=8 STATS_PERSISTENT=0 COMPRESSION='zlib' COMMENT='parent table';

		CREATE TABLE mem_tbl (
			id INT NOT NULL,
			INDEX k_mem (id)
		) ENGINE=MEMORY;

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

		CREATE TABLE linear_part_tbl (
			id INT NOT NULL
		) PARTITION BY LINEAR HASH (id) PARTITIONS 2;

		CREATE VIEW v_parent (vid, vcode) AS SELECT id, code FROM parent;

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
	require.Len(t, s.Tables, 5)
	require.Equal(t, []string{"child", "linear_part_tbl", "mem_tbl", "parent", "part_tbl"}, []string{
		s.Tables[0].Name, s.Tables[1].Name, s.Tables[2].Name, s.Tables[3].Name, s.Tables[4].Name,
	})

	parent := s.Tables[3]
	require.Equal(t, "parent table", parent.Comment)
	require.Equal(t, "InnoDB", parent.Engine)
	require.Contains(t, parent.CreateOptions, "row_format=DYNAMIC")
	require.Contains(t, parent.CreateOptions, "key_block_size=8")
	require.Contains(t, parent.CreateOptions, "stats_persistent=0")
	require.Contains(t, parent.CreateOptions, "compression='zlib'")
	require.Len(t, parent.Columns, 12, "system-hidden functional index column must be excluded")

	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range parent.Columns {
		cols[col.Name] = col
	}
	require.Equal(t, "AUTO_INCREMENT", cols["id"].Default)
	require.Equal(t, "'abc'", cols["code"].Default)
	require.Equal(t, "((1 + 2))", cols["expr_col"].Default)
	require.Equal(t, "b'101'", cols["bit_col"].Default)
	require.Equal(t, "0x05", cols["hex_col"].Default)
	require.Equal(t, "NULL", cols["nullable_col"].Default)
	require.Equal(t, "CURRENT_TIMESTAMP", cols["ts"].Default)
	require.Equal(t, "CURRENT_TIMESTAMP", cols["ts"].OnUpdate)
	require.Equal(t, metadata.GenerationMetadata_TYPE_STORED, cols["gen_stored"].Generation.Type)
	require.Equal(t, metadata.GenerationMetadata_TYPE_VIRTUAL, cols["gen_virt"].Generation.Type)
	require.NotNil(t, cols["g0"].Srid)
	require.Equal(t, uint32(0), *cols["g0"].Srid)
	require.NotNil(t, cols["g4326"].Srid)
	require.Equal(t, uint32(4326), *cols["g4326"].Srid)
	require.Nil(t, cols["id"].Srid)
	require.True(t, cols["secret"].IsInvisible)
	require.False(t, cols["id"].IsInvisible)

	var spIdx *metadata.IndexMetadata
	for _, idx := range parent.Indexes {
		if idx.Name == "sp_idx" {
			spIdx = idx
		}
	}
	require.NotNil(t, spIdx)
	require.Equal(t, "SPATIAL", spIdx.Type)
	require.Equal(t, []int64{32}, spIdx.KeyLength)

	require.Len(t, parent.Triggers, 2)
	require.Equal(t, []string{"trg_z_first", "trg_parent"}, []string{
		parent.Triggers[0].Name, parent.Triggers[1].Name,
	})
	require.Equal(t, "INSERT", parent.Triggers[1].Event)
	require.Equal(t, "AFTER", parent.Triggers[1].Timing)

	linPart := s.Tables[1]
	require.Len(t, linPart.Partitions, 2)
	require.Equal(t, "2", linPart.Partitions[0].UseDefault)
	require.Len(t, s.Tables[4].Partitions[0].Subpartitions, 2)
	require.Equal(t, "2", s.Tables[4].Partitions[0].Subpartitions[0].UseDefault)

	mem := s.Tables[2]
	require.Len(t, mem.Indexes, 1)
	require.Equal(t, "HASH", mem.Indexes[0].Type)

	// View exports explicit columns.
	require.Len(t, s.Views, 1)
	require.Equal(t, "v_parent", s.Views[0].Name)
	require.Len(t, s.Views[0].Columns, 2)
	require.Equal(t, "vid", s.Views[0].Columns[0].Name)
	require.Equal(t, "vcode", s.Views[0].Columns[1].Name)

	// Routines export full CREATE statements so LoadMetadata round-trips them.
	require.Len(t, s.Functions, 1)
	require.Contains(t, s.Functions[0].Definition, "CREATE")
	require.Contains(t, s.Functions[0].Definition, "FUNCTION `fn_add`")
	require.Equal(t, "adds two ints", s.Functions[0].Comment)
	require.Len(t, s.Procedures, 1)
	require.Contains(t, s.Procedures[0].Definition, "CREATE")
	require.Contains(t, s.Procedures[0].Definition, "PROCEDURE `proc_noop`")
	require.Equal(t, "noop proc", s.Procedures[0].Comment)

	// Round-trip through LoadMetadata preserves objects cleanly.
	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)
	require.Equal(t, []string{"vid", "vcode"}, reloaded.GetDatabase("app").Views["v_parent"].Columns)
	require.NotNil(t, reloaded.GetDatabase("app").Functions["fn_add"])
	require.NotNil(t, reloaded.GetDatabase("app").Procedures["proc_noop"])
	require.Equal(t, 2, reloaded.GetDatabase("app").GetTable("linear_part_tbl").Partitioning.NumParts)
	reloadedExported := reloaded.DatabaseMetadata("app")
	require.Equal(t, []string{"trg_z_first", "trg_parent"}, []string{
		reloadedExported.Schemas[0].Tables[3].Triggers[0].Name,
		reloadedExported.Schemas[0].Tables[3].Triggers[1].Name,
	})

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
	require.Equal(t, "parent_id", fkCon.Columns[0])
}

func TestDatabaseMetadataSessionContext(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{
		Name: "app",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:    "t",
				Columns: []*metadata.ColumnMetadata{{Name: "id", Position: 1, Type: "int"}},
				Triggers: []*metadata.TriggerMetadata{{
					Name:                "trg",
					Timing:              "BEFORE",
					Event:               "INSERT",
					Body:                "BEGIN END",
					SqlMode:             "STRICT_ALL_TABLES",
					CharacterSetClient:  "latin1",
					CollationConnection: "latin1_swedish_ci",
				}},
			}},
			Events: []*metadata.EventMetadata{{
				Name:                "ev",
				Definition:          "CREATE EVENT `ev` ON SCHEDULE AT '2026-01-01 00:00:00' DO SELECT 1",
				SqlMode:             "NO_ENGINE_SUBSTITUTION",
				CharacterSetClient:  "utf8mb4",
				CollationConnection: "utf8mb4_bin",
				TimeZone:            "+08:00",
			}},
		}},
	}
	c := New()
	report, err := c.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)

	exported := c.DatabaseMetadata("app")
	trg := exported.Schemas[0].Tables[0].Triggers[0]
	require.Equal(t, "STRICT_ALL_TABLES", trg.SqlMode)
	require.Equal(t, "latin1", trg.CharacterSetClient)
	require.Equal(t, "latin1_swedish_ci", trg.CollationConnection)

	ev := exported.Schemas[0].Events[0]
	require.Equal(t, "NO_ENGINE_SUBSTITUTION", ev.SqlMode)
	require.Equal(t, "utf8mb4", ev.CharacterSetClient)
	require.Equal(t, "utf8mb4_bin", ev.CollationConnection)
	require.Equal(t, "+08:00", ev.TimeZone)

	// Triggers created directly from DDL without synced session context must not
	// invent session context on export.
	plain := New()
	_, err = plain.Exec(`CREATE DATABASE app; USE app; CREATE TABLE t (id INT); CREATE TRIGGER tr_plain BEFORE INSERT ON t FOR EACH ROW BEGIN END;`, nil)
	require.NoError(t, err)
	plainTrg := plain.DatabaseMetadata("app").Schemas[0].Tables[0].Triggers[0]
	require.Empty(t, plainTrg.SqlMode)
	require.Empty(t, plainTrg.CharacterSetClient)
	require.Empty(t, plainTrg.CollationConnection)
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

func TestDatabaseMetadataRoundTripOptionsConstraintsDefaultsAndKeyAlgorithm(t *testing.T) {
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
		) KEY_BLOCK_SIZE=8 COMPRESSION='ZLIB' ENCRYPTION='Y' STATS_PERSISTENT=0 STATS_AUTO_RECALC=1 STATS_SAMPLE_PAGES=32 MIN_ROWS=10 MAX_ROWS=100 AVG_ROW_LENGTH=64 PACK_KEYS=1 CHECKSUM=1 DELAY_KEY_WRITE=1;
		ALTER TABLE full_tbl ALTER COLUMN drop_col DROP DEFAULT;
		CREATE TABLE key_alg (
			id INT NOT NULL,
			sub INT NOT NULL
		) PARTITION BY KEY ALGORITHM = 1 (id) PARTITIONS 2;
		CREATE TABLE sub_key_alg (
			id INT NOT NULL,
			sub INT NOT NULL
		) PARTITION BY RANGE (id) SUBPARTITION BY LINEAR KEY ALGORITHM = 1 (sub) SUBPARTITIONS 2 (
			PARTITION p0 VALUES LESS THAN (10),
			PARTITION p1 VALUES LESS THAN MAXVALUE
		);
		CREATE FUNCTION fn_coll() RETURNS INT DETERMINISTIC RETURN 1;
		CREATE PROCEDURE proc_coll() BEGIN SELECT 1; END;
		CREATE ALGORITHM=TEMPTABLE SQL SECURITY INVOKER VIEW v_exec AS SELECT 1 AS x;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	s := meta.Schemas[0]
	require.Equal(t, "utf8mb4_unicode_ci", s.Functions[0].DatabaseCollation)
	require.Equal(t, "utf8mb4_unicode_ci", s.Procedures[0].DatabaseCollation)

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

	keyAlg := s.Tables[1]
	require.Equal(t, int32(1), keyAlg.Partitions[0].Algorithm)
	subKeyAlg := s.Tables[3]
	require.Equal(t, int32(1), subKeyAlg.Partitions[0].Subpartitions[0].Algorithm)

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	reloadedFull := reloaded.GetDatabase("app").GetTable("full_tbl")
	require.Equal(t, 8, reloadedFull.KeyBlockSize)
	require.Equal(t, "ZLIB", reloadedFull.Compression)
	require.Equal(t, "Y", reloadedFull.Encryption)
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
	require.Equal(t, int32(1), reloaded.DatabaseMetadata("app").Schemas[0].Tables[1].Partitions[0].Algorithm)
	require.Equal(t, int32(1), reloaded.DatabaseMetadata("app").Schemas[0].Tables[3].Partitions[0].Subpartitions[0].Algorithm)
	require.True(t, reloaded.DatabaseMetadata("app").Schemas[0].Tables[0].CheckConstraints[1].NotEnforced)

	require.Len(t, s.Views, 1)
	require.Equal(t, "TEMPTABLE", s.Views[0].Algorithm)
	require.Equal(t, "INVOKER", s.Views[0].SqlSecurity)
	reView := reloaded.GetDatabase("app").Views["v_exec"]
	require.Equal(t, "TEMPTABLE", reView.Algorithm)
	require.Equal(t, "INVOKER", reView.SqlSecurity)
}

func TestDatabaseMetadataViewCheckOptionGIPKNotSecondaryAndCharsetExplicit(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		SET SESSION sql_generate_invisible_primary_key = ON;
		CREATE DATABASE app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
		USE app;

		CREATE TABLE gipk_tbl (
			name VARCHAR(50),
			cs_col VARCHAR(50) CHARACTER SET latin1,
			col_col VARCHAR(50) COLLATE utf8mb4_bin,
			no_sec INT NOT NULL NOT SECONDARY
		) TABLESPACE ts1;

		CREATE VIEW v_local AS SELECT name FROM gipk_tbl WITH LOCAL CHECK OPTION;
		CREATE VIEW v_cascaded AS SELECT name FROM gipk_tbl WITH CASCADED CHECK OPTION;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	s := meta.Schemas[0]
	require.Len(t, s.Views, 2)
	require.Equal(t, "CASCADED", s.Views[0].CheckOption)
	require.Equal(t, "LOCAL", s.Views[1].CheckOption)

	tblMeta := s.Tables[0]
	require.Equal(t, "ts1", tblMeta.Tablespace)
	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range tblMeta.Columns {
		cols[col.Name] = col
	}
	require.True(t, cols["my_row_id"].IsGeneratedInvisiblePrimaryKey)
	require.True(t, cols["no_sec"].NotSecondary)
	require.True(t, cols["name"].GetCharsetInherited())
	require.True(t, cols["name"].GetCollationInherited())
	require.False(t, cols["cs_col"].GetCharsetInherited())
	require.True(t, cols["cs_col"].GetCollationInherited())
	require.True(t, cols["col_col"].GetCharsetInherited())
	require.False(t, cols["col_col"].GetCollationInherited())
	require.False(t, cols["no_sec"].GetCharsetInherited())
	require.False(t, cols["no_sec"].GetCollationInherited())

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	reTbl := reloaded.GetDatabase("app").GetTable("gipk_tbl")
	require.Equal(t, "ts1", reTbl.Tablespace)
	require.True(t, reTbl.GetColumn("my_row_id").GeneratedInvisiblePrimaryKey)
	require.True(t, reTbl.GetColumn("no_sec").NotSecondary)
	require.False(t, reTbl.GetColumn("name").CharsetExplicit)
	require.False(t, reTbl.GetColumn("name").CollationExplicit)
	require.True(t, reTbl.GetColumn("cs_col").CharsetExplicit)
	require.False(t, reTbl.GetColumn("cs_col").CollationExplicit)
	require.False(t, reTbl.GetColumn("col_col").CharsetExplicit)
	require.True(t, reTbl.GetColumn("col_col").CollationExplicit)
	require.Equal(t, "LOCAL", reloaded.GetDatabase("app").Views["v_local"].CheckOption)
	require.Equal(t, "CASCADED", reloaded.GetDatabase("app").Views["v_cascaded"].CheckOption)

	// Reloaded GIPK table must match a user target without my_row_id (no phantom DROP).
	target := New()
	_, err = target.Exec(`
		CREATE DATABASE app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
		USE app;
		CREATE TABLE gipk_tbl (
			name VARCHAR(50),
			cs_col VARCHAR(50) CHARACTER SET latin1,
			col_col VARCHAR(50) COLLATE utf8mb4_bin,
			no_sec INT NOT NULL NOT SECONDARY
		) TABLESPACE ts1;
		CREATE VIEW v_local AS SELECT name FROM gipk_tbl WITH LOCAL CHECK OPTION;
		CREATE VIEW v_cascaded AS SELECT name FROM gipk_tbl WITH CASCADED CHECK OPTION;
	`, nil)
	require.NoError(t, err)
	diff := Diff(reloaded, target)
	require.True(t, diff.IsEmpty(), "unexpected diff after reload: %s", GenerateMigration(reloaded, target, diff).SQL())

	// Inherited column collation must propagate when table collation changes after reload.
	reloadedAlt := New()
	_, err = reloadedAlt.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	_, err = reloadedAlt.Exec(`ALTER TABLE gipk_tbl COLLATE = utf8mb4_0900_ai_ci;`, nil)
	require.NoError(t, err)
	n := NormalizerFor(MySQL80)
	altTbl := reloadedAlt.GetDatabase("app").GetTable("gipk_tbl")
	_, nameColl := n.ResolveColumnCharsetCollation(altTbl, altTbl.GetColumn("name"))
	_, csColl := n.ResolveColumnCharsetCollation(altTbl, altTbl.GetColumn("cs_col"))
	require.Equal(t, "utf8mb4_0900_ai_ci", nameColl)
	require.Equal(t, "latin1_swedish_ci", csColl)
}

func TestDatabaseMetadataIndexKeyBlockSizeRoundTrip(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE idx_kbs (
			id INT PRIMARY KEY,
			a INT,
			KEY ka (a) KEY_BLOCK_SIZE = 4
		);
	`, nil)
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	tblMeta := meta.Schemas[0].Tables[0]
	var ka *metadata.IndexMetadata
	for _, idx := range tblMeta.Indexes {
		if idx.Name == "ka" {
			ka = idx
		}
	}
	require.NotNil(t, ka)
	require.Equal(t, int64(4), ka.KeyBlockSize)

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	var reKA *Index
	for _, idx := range reloaded.GetDatabase("app").GetTable("idx_kbs").Indexes {
		if idx.Name == "ka" {
			reKA = idx
		}
	}
	require.NotNil(t, reKA)
	require.Equal(t, 4, reKA.KeyBlockSize)

	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "KEY `ka` (`a`) KEY_BLOCK_SIZE=4")
}

func TestDatabaseDefinitionFidelityForViewOptionsTablespaceAndComments(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE doc_tbl (
			id INT PRIMARY KEY,
			note VARCHAR(100) COMMENT 'it''s a \\ path'
		) TABLESPACE ts_my COMMENT='O''Reilly & Co';
		CREATE ALGORITHM=TEMPTABLE SQL SECURITY INVOKER VIEW v_sec AS SELECT id, note FROM doc_tbl WITH LOCAL CHECK OPTION;
	`, nil)
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "/*!50100 TABLESPACE `ts_my` */")
	require.Contains(t, def, `COMMENT='O\'Reilly & Co'`)
	require.Contains(t, def, `COMMENT 'it\'s a \\ path'`)
	require.Contains(t, def, "CREATE OR REPLACE ALGORITHM=TEMPTABLE SQL SECURITY INVOKER VIEW `v_sec`")
	require.Contains(t, def, "WITH LOCAL CHECK OPTION")

	// Rendered script must execute cleanly into a fresh catalog with matching comments & tablespace.
	restored := New()
	_, err = restored.Exec("CREATE DATABASE app;\nUSE app;\n"+def, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	reTbl := restored.GetDatabase("app").GetTable("doc_tbl")
	require.Equal(t, "O'Reilly & Co", reTbl.Comment)
	require.Equal(t, `it's a \ path`, reTbl.GetColumn("note").Comment)
	require.Equal(t, "ts_my", reTbl.Tablespace)
	reView := restored.GetDatabase("app").Views["v_sec"]
	require.Equal(t, "TEMPTABLE", reView.Algorithm)
	require.Equal(t, "INVOKER", reView.SqlSecurity)
	require.Equal(t, "LOCAL", reView.CheckOption)
}

func TestDatabaseDefinitionPreservesParentCommentWithExplicitSubpartitions(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{
		Name: "app",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:    "sub_t",
				Engine:  "InnoDB",
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
					}, {
						Type:       metadata.TablePartitionMetadata_HASH,
						Expression: "k",
						Name:       "s1",
					}},
				}},
			}},
		}},
	}

	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "PARTITION `p0` VALUES LESS THAN (10) ENGINE=InnoDB COMMENT = 'p0_parent'")
	require.Contains(t, def, "SUBPARTITION `s0` ENGINE=InnoDB COMMENT = 's0_child'")

	restored := New()
	_, err = restored.Exec("CREATE DATABASE app;\nUSE app;\n"+def, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	require.Equal(t, "p0_parent", restored.GetDatabase("app").GetTable("sub_t").Partitioning.Partitions[0].Comment)
}

func TestDatabaseDefinitionDefaultDroppedNotSecondaryExplicitCharsetAndImplicitSubpartitions(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		USE app;
		CREATE TABLE col_t (
			id INT PRIMARY KEY,
			a VARCHAR(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
			b VARCHAR(20),
			c INT NOT NULL NOT SECONDARY,
			d INT DEFAULT 1
		) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		ALTER TABLE col_t ALTER COLUMN d DROP DEFAULT;
		CREATE TABLE sub_default (
			id INT NOT NULL,
			k INT NOT NULL
		) PARTITION BY RANGE (id) SUBPARTITION BY HASH (k) (PARTITION p0 VALUES LESS THAN (10));
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "`a` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	require.Contains(t, def, "`b` varchar(20) DEFAULT NULL")
	require.Contains(t, def, "`c` int NOT NULL NOT SECONDARY")
	require.Contains(t, def, "ALTER TABLE `col_t` ALTER COLUMN `d` DROP DEFAULT;")
	require.NotContains(t, def, "SUBPARTITION `p0sp0`")
	require.Contains(t, def, "SUBPARTITION BY HASH (`k`)")

	reloaded := New()
	_, err = reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	reSub := reloaded.GetDatabase("app").GetTable("sub_default")
	require.True(t, reSub.Partitioning.UseDefaultSubpartitions)
	require.True(t, reSub.Partitioning.UseDefaultNumSubpartitions)

	execd := New()
	_, err = execd.Exec("CREATE DATABASE app DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;\nUSE app;\n"+def, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	reColT := execd.GetDatabase("app").GetTable("col_t")
	require.True(t, reColT.GetColumn("c").NotSecondary)
	require.True(t, reColT.GetColumn("d").DefaultDropped)
	require.True(t, reColT.GetColumn("a").CharsetExplicit)
	require.False(t, reColT.GetColumn("b").CharsetExplicit)
	reExecSub := execd.GetDatabase("app").GetTable("sub_default")
	require.True(t, reExecSub.Partitioning.UseDefaultSubpartitions)
	require.True(t, reExecSub.Partitioning.UseDefaultNumSubpartitions)
}

func TestDatabaseDefinitionGIPKRestoreAndQuotedHexDefault(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE no_pk (v INT);
		SET SESSION sql_generate_invisible_primary_key = ON;
		CREATE TABLE gipk_t (name VARCHAR(10));
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
	var hexTbl *metadata.TableMetadata
	for _, tbl := range meta.Schemas[0].Tables {
		if tbl.Name == "hex_t" {
			hexTbl = tbl
		}
	}
	require.NotNil(t, hexTbl)
	require.Equal(t, "X'05'", hexTbl.Columns[1].Default)
	require.Equal(t, "x'0a'", hexTbl.Columns[2].Default)

	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Incomplete)

	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "SQL_GENERATE_INVISIBLE_PRIMARY_KEY=OFF")
	require.Contains(t, def, "SQL_GENERATE_INVISIBLE_PRIMARY_KEY=@OLD_SQL_GENERATE_INVISIBLE_PRIMARY_KEY")
	require.Contains(t, def, "DEFAULT X'05'")
	require.Contains(t, def, "DEFAULT x'0a'")

	// Restoring in a session with GIPK=ON must not trigger error 4108, must not add my_row_id
	// to no-PK snapshot tables, and must keep GIPK=ON afterwards.
	restored := New()
	_, err = restored.Exec("SET SESSION sql_generate_invisible_primary_key = ON;\nCREATE DATABASE app;\nUSE app;\n"+def+"\nCREATE TABLE after_t (v INT);", &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)
	require.Nil(t, restored.GetDatabase("app").GetTable("no_pk").GetColumn("my_row_id"))
	require.NotNil(t, restored.GetDatabase("app").GetTable("gipk_t").GetColumn("my_row_id"))
	require.NotNil(t, restored.GetDatabase("app").GetTable("after_t").GetColumn("my_row_id"))
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
	require.Equal(t, "int", cols["id"].Type)
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

func TestDatabaseMetadataDecimalArithmeticAndSkipDumpRoundTrip(t *testing.T) {
	c := New()
	_, err := c.Exec(`
		CREATE DATABASE ddb;
		USE ddb;
		CREATE TABLE dec_t (id INT PRIMARY KEY, d1 DECIMAL(5,2), d2 DECIMAL(3,1));
		CREATE VIEW v_dec AS SELECT
			d1 + d2 AS p,
			d1 - d2 AS m,
			d1 * d2 AS mul,
			d1 / d2 AS dv,
			d1 % d2 AS md,
			-d1 AS neg_d,
			CAST(1 AS DECIMAL(5,2)) + CAST(2 AS DECIMAL(3,1)) AS cp
		FROM dec_t;
	`, &ExecOptions{ContinueOnError: false})
	require.NoError(t, err)

	meta := c.DatabaseMetadata("ddb")
	vcols := map[string]*metadata.ColumnMetadata{}
	for _, col := range meta.Schemas[0].Views[0].Columns {
		vcols[col.Name] = col
	}
	require.Equal(t, "decimal(6,2)", vcols["p"].Type)
	require.Equal(t, "decimal(6,2)", vcols["m"].Type)
	require.Equal(t, "decimal(8,3)", vcols["mul"].Type)
	require.Equal(t, "decimal(10,6)", vcols["dv"].Type)
	require.Equal(t, "decimal(5,2)", vcols["md"].Type)
	require.Equal(t, "decimal(5,2)", vcols["neg_d"].Type)
	require.Equal(t, "decimal(6,2)", vcols["cp"].Type)

	// SkipDump round-trips across LoadMetadata -> DatabaseMetadata.
	skipMeta := &metadata.DatabaseSchemaMetadata{
		Name: "skip_db",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:     "t_skip",
				SkipDump: true,
				Columns:  []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: false}},
				Triggers: []*metadata.TriggerMetadata{{
					Name: "tr_skip", Event: "INSERT", Timing: "BEFORE", Body: "BEGIN END", SkipDump: true,
				}},
			}},
			Views: []*metadata.ViewMetadata{{
				Name: "v_skip", Definition: "SELECT 1 AS c", SkipDump: true,
			}},
			Functions: []*metadata.FunctionMetadata{{
				Name: "fn_skip", Definition: "CREATE FUNCTION fn_skip() RETURNS INT DETERMINISTIC RETURN 1", SkipDump: true,
			}},
			Procedures: []*metadata.ProcedureMetadata{{
				Name: "sp_skip", Definition: "CREATE PROCEDURE sp_skip() BEGIN END", SkipDump: true,
			}},
		}},
	}
	skipCat := New()
	report, err := skipCat.LoadMetadata(context.Background(), skipMeta)
	require.NoError(t, err)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Degraded)
	out := skipCat.DatabaseMetadata("skip_db")
	require.True(t, out.Schemas[0].Tables[0].SkipDump)
	require.True(t, out.Schemas[0].Tables[0].Triggers[0].SkipDump)
	require.True(t, out.Schemas[0].Views[0].SkipDump)
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
		CREATE VIEW v_left_star AS SELECT a.id AS aid, b.* FROM local_t a LEFT JOIN ref_db.ref_t b ON a.id = b.id;
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

	viewByName := map[string]*metadata.ViewMetadata{}
	for _, v := range meta.Schemas[0].Views {
		viewByName[v.Name] = v
	}
	for _, vName := range []string{"v_star", "v_qual_star", "v_cross_star"} {
		v := viewByName[vName]
		require.NotNil(t, v, vName)
		require.Len(t, v.Columns, 3, vName)
		require.Equal(t, "int", v.Columns[0].Type, vName)
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
	leftV := viewByName["v_left_star"]
	require.Len(t, leftV.Columns, 4)
	require.False(t, leftV.Columns[0].Nullable)
	require.True(t, leftV.Columns[1].Nullable)
	require.True(t, leftV.Columns[2].Nullable)
	require.Equal(t, "utf8mb4", leftV.Columns[2].CharacterSet)
	require.Equal(t, "utf8mb4_unicode_ci", leftV.Columns[2].Collation)

	// Functional index expressions export parenthesized so DatabaseDefinition restores them.
	_, err = c.Exec(`
		CREATE TABLE fn_idx_t (
			id INT NOT NULL PRIMARY KEY,
			a INT NOT NULL,
			b INT NOT NULL,
			code VARCHAR(20),
			INDEX fn_upper ((UPPER(code))),
			INDEX fn_sum (((a) + (b)))
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
	require.Contains(t, fnDef, "KEY `fn_sum` ((`a` + `b`))")
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

	require.Equal(t, "int", uViews["v_null"].Columns[0].Type)
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
	require.Equal(t, "int", vMM.Columns[0].Type)
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
				Engine:    "InnoDB",
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

	// Unresolved view columns retain supplied snapshot metadata across LoadMetadata -> DatabaseMetadata.
	extCat := New()
	_, err = extCat.LoadMetadata(context.Background(), &metadata.DatabaseSchemaMetadata{
		Name: "app",
		Schemas: []*metadata.SchemaMetadata{{
			Views: []*metadata.ViewMetadata{{
				Name:       "v_ext",
				Definition: "select `unloaded_db`.`ext_tbl`.`id` AS `id`,`unloaded_db`.`ext_tbl`.`name` AS `name` from `unloaded_db`.`ext_tbl`",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Position: 1, Type: "int", Nullable: false},
					{Name: "name", Position: 2, Type: "varchar(32)", Nullable: true, CharacterSet: "utf8mb4", Collation: "utf8mb4_unicode_ci"},
				},
			}},
		}},
	})
	require.NoError(t, err)
	reExt := extCat.DatabaseMetadata("app").Schemas[0].Views[0]
	require.Len(t, reExt.Columns, 2)
	require.Equal(t, "id", reExt.Columns[0].Name)
	require.Equal(t, "int", reExt.Columns[0].Type)
	require.False(t, reExt.Columns[0].Nullable)
	require.Equal(t, "name", reExt.Columns[1].Name)
	require.Equal(t, "varchar(32)", reExt.Columns[1].Type)
	require.True(t, reExt.Columns[1].Nullable)
	require.Equal(t, "utf8mb4", reExt.Columns[1].CharacterSet)
	require.Equal(t, "utf8mb4_unicode_ci", reExt.Columns[1].Collation)

	// Aggregate types (SUM/AVG) and nullability (SUM/AVG/MIN/MAX/GROUP_CONCAT nullable, COUNT not null),
	// typed temporal defaults (SET DEFAULT DATE '...'), and routine normalization (no DEFINER / RETURNS CHARSET).
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
	require.False(t, aggViewCols["cnt"].Nullable)
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
}
