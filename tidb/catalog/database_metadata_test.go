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
