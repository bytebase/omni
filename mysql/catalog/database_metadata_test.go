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
