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
		SET tidb_enable_check_constraint = 1;
		CREATE DATABASE app CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
		USE app;

		CREATE TABLE items (
			id BIGINT AUTO_RANDOM(5, 54) PRIMARY KEY,
			code VARCHAR(20) NOT NULL DEFAULT 'abc',
			nullable_col INT,
			ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			gen_col VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_g')) STORED,
			CHECK (nullable_col >= 0)
		) ENGINE=InnoDB COMMENT='items table';

		CREATE TABLE items_bare_autorand (
			id BIGINT AUTO_RANDOM PRIMARY KEY
		);

		CREATE TABLE items_bits_autorand (
			id BIGINT AUTO_RANDOM(4) PRIMARY KEY
		);

		CREATE TABLE part_tbl (
			id INT NOT NULL,
			yr INT NOT NULL,
			PRIMARY KEY (id, yr)
		) PARTITION BY RANGE (yr) (
			PARTITION p0 VALUES LESS THAN (2025),
			PARTITION p1 VALUES LESS THAN MAXVALUE
		);

		CREATE VIEW v_items AS SELECT id, code FROM items;
		CREATE SEQUENCE seq_order START WITH 10 INCREMENT BY 2 MINVALUE 1 MAXVALUE 1000 CACHE 50 CYCLE;
	`, nil)
	require.NoError(t, err)

	meta := c.DatabaseMetadata("app")
	require.Equal(t, "app", meta.Name)
	require.Equal(t, "utf8mb4", meta.CharacterSet)
	require.Equal(t, "utf8mb4_bin", meta.Collation)
	require.Len(t, meta.Schemas, 1)

	s := meta.Schemas[0]
	require.Len(t, s.Tables, 4)

	items := s.Tables[0]
	require.Equal(t, "items", items.Name)
	require.Empty(t, items.Charset)
	require.Equal(t, "utf8mb4_bin", items.Collation)
	require.Equal(t, "items table", items.Comment)

	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range items.Columns {
		cols[col.Name] = col
	}
	require.Equal(t, "AUTO_RANDOM(5, 54)", cols["id"].Default)
	require.Equal(t, "'abc'", cols["code"].Default)
	require.Empty(t, cols["nullable_col"].Default, "TiDB leaves nullable column without default empty")
	require.Equal(t, "CURRENT_TIMESTAMP", cols["ts"].Default)
	require.Equal(t, "CURRENT_TIMESTAMP", cols["ts"].OnUpdate)
	require.Equal(t, metadata.GenerationMetadata_TYPE_STORED, cols["gen_col"].Generation.Type)
	require.Len(t, items.CheckConstraints, 1)

	require.Equal(t, "AUTO_RANDOM", s.Tables[1].Columns[0].Default)
	require.Equal(t, "AUTO_RANDOM(4)", s.Tables[2].Columns[0].Default)

	partTbl := s.Tables[3]
	require.Len(t, partTbl.Partitions, 2)
	require.Equal(t, metadata.TablePartitionMetadata_RANGE, partTbl.Partitions[0].Type)

	require.Len(t, s.Views, 1)
	require.Equal(t, "v_items", s.Views[0].Name)
	require.Len(t, s.Sequences, 1)
	require.Equal(t, "seq_order", s.Sequences[0].Name)
	require.Equal(t, "10", s.Sequences[0].Start)
	require.Equal(t, "2", s.Sequences[0].Increment)
	require.Equal(t, "50", s.Sequences[0].CacheSize)
	require.True(t, s.Sequences[0].Cycle)
}
