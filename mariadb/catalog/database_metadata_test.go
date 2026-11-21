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
	results, err := c.Exec(`
		CREATE DATABASE app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
		USE app;

		CREATE SEQUENCE seq_order START WITH 100 INCREMENT BY 5 MINVALUE 1 MAXVALUE 9999 CACHE 20 CYCLE COMMENT='order seq';
		CREATE SEQUENCE seq_u32 AS INT UNSIGNED START WITH 1 NOCACHE NOCYCLE;
		CREATE TEMPORARY SEQUENCE tmp_seq START WITH 1;

		CREATE TABLE parent (
			id INT AUTO_INCREMENT PRIMARY KEY,
			code VARCHAR(20) NOT NULL DEFAULT 'abc',
			expr_str VARCHAR(20) NOT NULL DEFAULT '1 + 2',
			expr_col INT DEFAULT (1 + 2),
			fn_col DATE DEFAULT curdate(),
			bit_col BIT(3) DEFAULT b'101',
			hex_col BINARY(1) DEFAULT 0x05,
			nullable_col INT,
			ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			gen_stored VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_s')) STORED,
			gen_virt VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_v')) VIRTUAL,
			g POINT NOT NULL,
			secret INT INVISIBLE NOT NULL DEFAULT 0,
			SPATIAL INDEX sp_idx (g),
			INDEX idx_ignored (code) IGNORED COMMENT 'not used by optimizer',
			INDEX idx_prefix (code(10) DESC)
		) ENGINE=InnoDB ROW_FORMAT=DYNAMIC KEY_BLOCK_SIZE=8 STATS_PERSISTENT=0 COMMENT='parent''s table';

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

		CREATE TABLE versioned_implicit (
			id INT PRIMARY KEY,
			val VARCHAR(50)
		) WITH SYSTEM VERSIONING;

		CREATE TABLE versioned_explicit (
			id INT PRIMARY KEY,
			val VARCHAR(50),
			row_start TIMESTAMP(6) GENERATED ALWAYS AS ROW START,
			row_end TIMESTAMP(6) GENERATED ALWAYS AS ROW END,
			PERIOD FOR SYSTEM_TIME (row_start, row_end)
		) WITH SYSTEM VERSIONING;

		CREATE OR REPLACE VIEW v_parent (vid, vcode) AS SELECT id, code FROM parent WITH CASCADED CHECK OPTION;

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

		CREATE EVENT ev_daily ON SCHEDULE EVERY 1 DAY COMMENT 'daily job' DO SELECT 1;
	`, nil)
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, "statement %d failed", r.Index)
	}

	meta := c.DatabaseMetadata("app")
	require.Equal(t, "app", meta.Name)
	require.Equal(t, "utf8mb4", meta.CharacterSet)
	require.Equal(t, "utf8mb4_unicode_ci", meta.Collation)
	require.Len(t, meta.Schemas, 1)

	s := meta.Schemas[0]

	// Sequences exclude temporary sequences and preserve attributes.
	require.Len(t, s.Sequences, 2)
	require.Equal(t, "seq_order", s.Sequences[0].Name)
	require.Equal(t, "bigint", s.Sequences[0].DataType)
	require.Equal(t, "100", s.Sequences[0].Start)
	require.Equal(t, "5", s.Sequences[0].Increment)
	require.Equal(t, "1", s.Sequences[0].MinValue)
	require.Equal(t, "9999", s.Sequences[0].MaxValue)
	require.Equal(t, "20", s.Sequences[0].CacheSize)
	require.True(t, s.Sequences[0].Cycle)
	require.Equal(t, "order seq", s.Sequences[0].Comment)

	require.Equal(t, "seq_u32", s.Sequences[1].Name)
	require.Equal(t, "int unsigned", s.Sequences[1].DataType)
	require.Equal(t, "0", s.Sequences[1].CacheSize)
	require.False(t, s.Sequences[1].Cycle)

	// Tables are sorted by name and exclude temp tables.
	require.Len(t, s.Tables, 7)
	tableNames := make([]string, len(s.Tables))
	byName := make(map[string]*metadata.TableMetadata, len(s.Tables))
	for i, tbl := range s.Tables {
		tableNames[i] = tbl.Name
		byName[tbl.Name] = tbl
	}
	require.Equal(t, []string{
		"child", "linear_part_tbl", "mem_tbl", "parent", "part_tbl", "versioned_explicit", "versioned_implicit",
	}, tableNames)

	parent := byName["parent"]
	require.Equal(t, `parent\'s table`, parent.Comment)
	require.Equal(t, "InnoDB", parent.Engine)
	require.Contains(t, parent.CreateOptions, "row_format=DYNAMIC")
	require.Contains(t, parent.CreateOptions, "key_block_size=8")
	require.Contains(t, parent.CreateOptions, "stats_persistent=0")
	require.Len(t, parent.Columns, 13, "system-hidden functional index column must be excluded")

	cols := map[string]*metadata.ColumnMetadata{}
	for _, col := range parent.Columns {
		cols[col.Name] = col
	}
	require.Equal(t, "AUTO_INCREMENT", cols["id"].Default)
	require.Equal(t, `'\'abc\''`, cols["code"].Default)
	require.Equal(t, `'\'1 + 2\''`, cols["expr_str"].Default)
	require.Equal(t, "'1 + 2'", cols["expr_col"].Default)
	require.Equal(t, "'curdate()'", cols["fn_col"].Default)
	require.Equal(t, `'b\'101\''`, cols["bit_col"].Default)
	require.Equal(t, `'x\'05\''`, cols["hex_col"].Default)
	require.Equal(t, "NULL", cols["nullable_col"].Default)
	require.Equal(t, "CURRENT_TIMESTAMP", cols["ts"].Default)
	require.Equal(t, "CURRENT_TIMESTAMP", cols["ts"].OnUpdate)
	require.Equal(t, metadata.GenerationMetadata_TYPE_STORED, cols["gen_stored"].Generation.Type)
	require.Equal(t, metadata.GenerationMetadata_TYPE_VIRTUAL, cols["gen_virt"].Generation.Type)
	require.True(t, cols["secret"].IsInvisible)
	require.False(t, cols["id"].IsInvisible)

	idxByName := map[string]*metadata.IndexMetadata{}
	for _, idx := range parent.Indexes {
		idxByName[idx.Name] = idx
	}
	require.Equal(t, "SPATIAL", idxByName["sp_idx"].Type)
	require.Equal(t, []int64{32}, idxByName["sp_idx"].KeyLength)
	require.False(t, idxByName["idx_ignored"].Visible)
	require.Equal(t, "not used by optimizer", idxByName["idx_ignored"].Comment)
	require.True(t, idxByName["PRIMARY"].Visible)
	require.Equal(t, []int64{10}, idxByName["idx_prefix"].KeyLength)
	require.Equal(t, []bool{true}, idxByName["idx_prefix"].Descending)

	// Triggers respect PRECEDES ordering.
	require.Len(t, parent.Triggers, 2)
	require.Equal(t, []string{"trg_z_first", "trg_parent"}, []string{
		parent.Triggers[0].Name, parent.Triggers[1].Name,
	})
	require.Equal(t, "INSERT", parent.Triggers[1].Event)
	require.Equal(t, "AFTER", parent.Triggers[1].Timing)

	// System-versioned tables export WITH SYSTEM VERSIONING and ROW START/ROW END.
	vImp := byName["versioned_implicit"]
	require.Contains(t, vImp.CreateOptions, "WITH SYSTEM VERSIONING")

	vExp := byName["versioned_explicit"]
	require.Contains(t, vExp.CreateOptions, "WITH SYSTEM VERSIONING")
	vExpCols := map[string]*metadata.ColumnMetadata{}
	for _, col := range vExp.Columns {
		vExpCols[col.Name] = col
	}
	require.NotNil(t, vExpCols["row_start"].Generation)
	require.Equal(t, metadata.GenerationMetadata_TYPE_STORED, vExpCols["row_start"].Generation.Type)
	require.Equal(t, "ROW START", vExpCols["row_start"].Generation.Expression)
	require.NotNil(t, vExpCols["row_end"].Generation)
	require.Equal(t, metadata.GenerationMetadata_TYPE_STORED, vExpCols["row_end"].Generation.Type)
	require.Equal(t, "ROW END", vExpCols["row_end"].Generation.Expression)

	// Partitions & MEMORY engine default index type.
	require.Len(t, byName["linear_part_tbl"].Partitions, 2)
	require.Equal(t, "2", byName["linear_part_tbl"].Partitions[0].UseDefault)
	require.Len(t, byName["part_tbl"].Partitions[0].Subpartitions, 2)
	require.Equal(t, "2", byName["part_tbl"].Partitions[0].Subpartitions[0].UseDefault)
	require.Equal(t, "HASH", byName["mem_tbl"].Indexes[0].Type)

	// Foreign key and check constraints.
	child := byName["child"]
	require.Len(t, child.ForeignKeys, 1)
	require.Equal(t, "CASCADE", child.ForeignKeys[0].OnDelete)
	require.Equal(t, "RESTRICT", child.ForeignKeys[0].OnUpdate)
	require.Len(t, child.CheckConstraints, 1)
	require.Equal(t, "chk_score", child.CheckConstraints[0].Name)

	// View exports explicit columns and CHECK OPTION in both CheckOption and Definition.
	require.Len(t, s.Views, 1)
	require.Equal(t, "v_parent", s.Views[0].Name)
	require.Equal(t, "CASCADED", s.Views[0].CheckOption)
	require.True(t, strings.HasSuffix(s.Views[0].Definition, " WITH CASCADED CHECK OPTION"),
		"Definition = %q", s.Views[0].Definition)
	require.Len(t, s.Views[0].Columns, 2)
	require.Equal(t, "vid", s.Views[0].Columns[0].Name)
	require.Equal(t, "vcode", s.Views[0].Columns[1].Name)
	require.Equal(t, "varchar(20)", s.Views[0].Columns[1].Type)
	require.False(t, s.Views[0].Columns[1].Nullable)

	// DatabaseDefinition and ViewDefinition preserve CHECK OPTION from DatabaseMetadata output.
	viewDDL, err := ViewDefinition(s.Views[0])
	require.NoError(t, err)
	require.Contains(t, viewDDL, "WITH CASCADED CHECK OPTION")
	dbDump, err := DatabaseDefinition(meta, DefinitionOptions{PrintHeader: true})
	require.NoError(t, err)
	require.Contains(t, dbDump, "WITH CASCADED CHECK OPTION")
	dbSDL, err := DatabaseDefinition(meta, DefinitionOptions{SDLFormat: true})
	require.NoError(t, err)
	require.Contains(t, dbSDL, "WITH CASCADED CHECK OPTION")

	// Routines and events.
	require.Len(t, s.Functions, 1)
	require.Contains(t, s.Functions[0].Definition, "FUNCTION `fn_add`")
	require.NotContains(t, s.Functions[0].Definition, "DEFINER=")
	require.Equal(t, "adds two ints", s.Functions[0].Comment)

	require.Len(t, s.Procedures, 1)
	require.Contains(t, s.Procedures[0].Definition, "PROCEDURE `proc_noop`")
	require.NotContains(t, s.Procedures[0].Definition, "DEFINER=")
	require.Equal(t, "noop proc", s.Procedures[0].Comment)

	require.Len(t, s.Events, 1)
	require.Equal(t, "ev_daily", s.Events[0].Name)
	require.Equal(t, "daily job", s.Events[0].Comment)
	require.Contains(t, s.Events[0].Definition, "EVENT `ev_daily`")

	// Round-trip through LoadMetadataInto preserves objects cleanly with empty Diff.
	reloaded := New()
	report, err := reloaded.LoadMetadataInto(context.Background(), meta, "app")
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	require.NotNil(t, reloaded.GetDatabase("app").Sequences["seq_order"])
	require.Equal(t, "int unsigned", reloaded.GetDatabase("app").Sequences["seq_u32"].DataType)
	require.True(t, reloaded.GetDatabase("app").GetTable("versioned_implicit").SystemVersioned)
	reloadedVExp := reloaded.GetDatabase("app").GetTable("versioned_explicit")
	require.True(t, reloadedVExp.SystemVersioned)
	require.Equal(t, "row_start", reloadedVExp.PeriodStartCol)
	require.Equal(t, "row_end", reloadedVExp.PeriodEndCol)
	require.Equal(t, []string{"vid", "vcode"}, reloaded.GetDatabase("app").Views["v_parent"].Columns)
	require.Equal(t, "CASCADED", reloaded.GetDatabase("app").Views["v_parent"].CheckOption)
	require.NotNil(t, reloaded.GetDatabase("app").Functions["fn_add"])
	require.NotNil(t, reloaded.GetDatabase("app").Procedures["proc_noop"])
	require.NotNil(t, reloaded.GetDatabase("app").Events["ev_daily"])
	require.Equal(t, 2, reloaded.GetDatabase("app").GetTable("linear_part_tbl").Partitioning.NumParts)

	// Verify string literal default that looks like an expression round-trips as constant, not expression.
	reloadedParent := reloaded.GetDatabase("app").GetTable("parent")
	require.Equal(t, ColumnDefaultConstant, reloadedParent.GetColumn("expr_str").DefaultKind)
	require.Equal(t, ColumnDefaultExpression, reloadedParent.GetColumn("expr_col").DefaultKind)

	// Mutating exported FK slices must not mutate catalog state.
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

func TestDatabaseMetadataCreateOrReplaceAndJSON(t *testing.T) {
	c := New()
	results, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE docs (id INT PRIMARY KEY);
		CREATE OR REPLACE TABLE docs (
			id INT PRIMARY KEY,
			payload JSON NOT NULL,
			seq_val INT DEFAULT NEXT VALUE FOR s1
		);
		CREATE SEQUENCE s1;
	`, nil)
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, "statement %d failed", r.Index)
	}

	meta := c.DatabaseMetadata("app")
	require.Len(t, meta.Schemas, 1)
	require.Len(t, meta.Schemas[0].Tables, 1)
	docs := meta.Schemas[0].Tables[0]
	require.Equal(t, "docs", docs.Name)
	require.Len(t, docs.Columns, 3)

	// Round-trip JSON + sequence default through LoadMetadata.
	reloaded := New()
	report, err := reloaded.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)
	require.Equal(t, c.ShowCreateTable("app", "docs"), reloaded.ShowCreateTable("app", "docs"))
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

func TestDatabaseMetadataViewAnalyzedColumns(t *testing.T) {
	c := New()
	results, err := c.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE a (id INT NOT NULL, label VARCHAR(32) NOT NULL);
		CREATE TABLE b (id INT NOT NULL, note VARCHAR(64) NOT NULL);
		CREATE VIEW v_join AS
			SELECT a.id, a.label, b.note
			FROM a LEFT JOIN b ON a.id = b.id;
		CREATE VIEW v_union AS
			SELECT id, label FROM a
			UNION ALL
			SELECT id, note FROM b;
		CREATE VIEW v_cte AS
			WITH x AS (SELECT 'abc' AS s, id FROM a)
			SELECT s, id FROM x;
	`, nil)
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, "statement %d failed", r.Index)
	}

	meta := c.DatabaseMetadata("app")
	views := map[string]*metadata.ViewMetadata{}
	for _, v := range meta.Schemas[0].Views {
		views[v.Name] = v
	}

	vJoin := views["v_join"]
	require.Len(t, vJoin.Columns, 3)
	require.False(t, vJoin.Columns[0].Nullable)
	require.False(t, vJoin.Columns[1].Nullable)
	require.True(t, vJoin.Columns[2].Nullable, "LEFT JOIN right-side column must be nullable")
	require.Equal(t, "varchar(64)", vJoin.Columns[2].Type)

	vUnion := views["v_union"]
	require.Len(t, vUnion.Columns, 2)
	require.Equal(t, "varchar(64)", vUnion.Columns[1].Type)

	vCTE := views["v_cte"]
	require.Len(t, vCTE.Columns, 2)
	require.Equal(t, "s", vCTE.Columns[0].Name)
	require.Equal(t, "varchar(3)", vCTE.Columns[0].Type)
	require.Equal(t, "id", vCTE.Columns[1].Name)
	require.Equal(t, "int(11)", vCTE.Columns[1].Type)
}

func TestDatabaseMetadataFunctionalIndexAndSystemTimePartition(t *testing.T) {
	// Functional index loaded via LoadMetadata excludes its system-hidden backing column.
	meta := &metadata.DatabaseSchemaMetadata{
		Name: "app",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name: "t",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Type: "int", Default: autoIncrementDefault},
					{Name: "name", Type: "varchar(64)", Nullable: true},
				},
				Indexes: []*metadata.IndexMetadata{
					{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Visible: true, Expressions: []string{"id"}},
					{Name: "idx_lower", Type: "BTREE", Visible: true, Expressions: []string{"(lower(`name`))"}},
				},
			}},
		}},
	}
	c := New()
	report, err := c.LoadMetadata(context.Background(), meta)
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Empty(t, report.Missing)
	require.Empty(t, report.Incomplete)

	exported := c.DatabaseMetadata("app")
	tbl := exported.Schemas[0].Tables[0]
	require.Len(t, tbl.Columns, 2, "system-hidden functional index column must be excluded")
	require.Equal(t, []string{"(lcase(`name`))"}, tbl.Indexes[1].Expressions)

	// SYSTEM_TIME partitioning has no proto enum and must not emit TYPE_UNSPECIFIED partitions.
	c2 := New()
	_, err = c2.Exec(`
		CREATE DATABASE app;
		USE app;
		CREATE TABLE hist (id INT PRIMARY KEY) WITH SYSTEM VERSIONING
		PARTITION BY SYSTEM_TIME (
			PARTITION p_hist HISTORY,
			PARTITION p_cur CURRENT
		);
	`, nil)
	require.NoError(t, err)

	exp2 := c2.DatabaseMetadata("app")
	require.Empty(t, exp2.Schemas[0].Tables[0].Partitions)
	require.Contains(t, exp2.Schemas[0].Tables[0].CreateOptions, "WITH SYSTEM VERSIONING")

	reloaded := New()
	rep2, err := reloaded.LoadMetadata(context.Background(), exp2)
	require.NoError(t, err)
	require.Empty(t, rep2.Degraded)
	require.Empty(t, rep2.Missing)
	require.Empty(t, rep2.Incomplete)
}
