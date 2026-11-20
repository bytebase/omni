package catalog

import (
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
			nullable_col INT,
			ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			gen_stored VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_s')) STORED,
			gen_virt VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_v')) VIRTUAL,
			g0 GEOMETRY NOT NULL SRID 0,
			g4326 POINT NOT NULL SRID 4326,
			secret INT INVISIBLE NOT NULL DEFAULT 0,
			SPATIAL INDEX sp_idx (g4326),
			INDEX fn_idx ((UPPER(code)))
		) ENGINE=InnoDB COMMENT='parent table';

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

		CREATE VIEW v_parent AS SELECT id, code FROM parent;

		CREATE FUNCTION fn_add(a INT, b INT) RETURNS INT DETERMINISTIC
		BEGIN
			RETURN a + b;
		END;

		CREATE PROCEDURE proc_noop()
		BEGIN
			SELECT 1;
		END;

		CREATE TRIGGER trg_parent AFTER INSERT ON parent FOR EACH ROW
		BEGIN
			INSERT INTO child(id) VALUES (NEW.id);
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

	// Check parent table details.
	parent := s.Tables[3]
	require.Equal(t, "parent table", parent.Comment)
	require.Equal(t, "InnoDB", parent.Engine)
	require.Len(t, parent.Columns, 10, "system-hidden functional index column must be excluded")

	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range parent.Columns {
		cols[col.Name] = col
	}
	require.Equal(t, "AUTO_INCREMENT", cols["id"].Default)
	require.Equal(t, "'abc'", cols["code"].Default)
	require.Equal(t, "((1 + 2))", cols["expr_col"].Default)
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

	// Spatial index key length 32.
	var spIdx *metadata.IndexMetadata
	for _, idx := range parent.Indexes {
		if idx.Name == "sp_idx" {
			spIdx = idx
		}
	}
	require.NotNil(t, spIdx)
	require.Equal(t, "SPATIAL", spIdx.Type)
	require.Equal(t, []int64{32}, spIdx.KeyLength)

	// Trigger attached to parent table.
	require.Len(t, parent.Triggers, 1)
	require.Equal(t, "trg_parent", parent.Triggers[0].Name)
	require.Equal(t, "INSERT", parent.Triggers[0].Event)
	require.Equal(t, "AFTER", parent.Triggers[0].Timing)

	// MEMORY table defaults index type to HASH.
	mem := s.Tables[2]
	require.Len(t, mem.Indexes, 1)
	require.Equal(t, "HASH", mem.Indexes[0].Type)

	// Foreign key and check constraint on child.
	child := s.Tables[0]
	require.Len(t, child.ForeignKeys, 1)
	require.Equal(t, "fk_parent", child.ForeignKeys[0].Name)
	require.Equal(t, "parent", child.ForeignKeys[0].ReferencedTable)
	require.Equal(t, "CASCADE", child.ForeignKeys[0].OnDelete)
	require.Len(t, child.CheckConstraints, 1)
	require.Equal(t, "chk_score", child.CheckConstraints[0].Name)

	// Subpartitions.
	partTbl := s.Tables[4]
	require.Len(t, partTbl.Partitions, 2)
	require.Equal(t, metadata.TablePartitionMetadata_RANGE, partTbl.Partitions[0].Type)
	require.Len(t, partTbl.Partitions[0].Subpartitions, 2)
	require.Equal(t, metadata.TablePartitionMetadata_HASH, partTbl.Partitions[0].Subpartitions[0].Type)
	require.Equal(t, "p0sp0", partTbl.Partitions[0].Subpartitions[0].Name)

	// Linear hash partition.
	linearTbl := s.Tables[1]
	require.Len(t, linearTbl.Partitions, 2)
	require.Equal(t, metadata.TablePartitionMetadata_LINEAR_HASH, linearTbl.Partitions[0].Type)

	// Views, functions, procedures, events.
	require.Len(t, s.Views, 1)
	require.Equal(t, "v_parent", s.Views[0].Name)
	require.Len(t, s.Functions, 1)
	require.Equal(t, "fn_add", s.Functions[0].Name)
	require.Len(t, s.Procedures, 1)
	require.Equal(t, "proc_noop", s.Procedures[0].Name)
	require.Len(t, s.Events, 1)
	require.Equal(t, "ev_daily", s.Events[0].Name)
	require.Contains(t, s.Events[0].Definition, "CREATE")
}
