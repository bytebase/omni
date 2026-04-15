package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// alterTableTestCase holds a single test case for ALTER TABLE parsing.
type alterTableTestCase struct {
	sql     string
	wantErr bool
	check   func(t *testing.T, stmt *ast.AlterTableStmt)
}

func runAlterTableTests(t *testing.T, cases []alterTableTestCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			file, errs := Parse(tc.sql)
			if tc.wantErr {
				if len(errs) == 0 {
					t.Error("expected error but got none")
				}
				return
			}
			if len(errs) > 0 {
				t.Fatalf("unexpected parse errors: %v", errs)
			}
			if file == nil || len(file.Stmts) == 0 {
				t.Fatal("no statements parsed")
			}
			stmt, ok := file.Stmts[0].(*ast.AlterTableStmt)
			if !ok {
				t.Fatalf("expected *AlterTableStmt, got %T", file.Stmts[0])
			}
			if tc.check != nil {
				tc.check(t, stmt)
			}
		})
	}
}

// ---- ADD COLUMN ----

func TestAlterTableAddColumn(t *testing.T) {
	cases := []alterTableTestCase{
		{
			sql: "ALTER TABLE t ADD COLUMN c INT",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				if len(stmt.Actions) != 1 {
					t.Fatalf("expected 1 action, got %d", len(stmt.Actions))
				}
				a := stmt.Actions[0]
				if a.Type != ast.AlterAddColumn {
					t.Errorf("expected AlterAddColumn, got %v", a.Type)
				}
				if a.Column == nil || a.Column.Name != "c" {
					t.Errorf("expected column name 'c', got %v", a.Column)
				}
			},
		},
		{
			sql: "ALTER TABLE t ADD COLUMN c INT AFTER b",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.After != "b" {
					t.Errorf("expected After='b', got %q", a.After)
				}
			},
		},
		{
			sql: "ALTER TABLE t ADD COLUMN c INT FIRST",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if !a.First {
					t.Error("expected First=true")
				}
			},
		},
		{
			sql: "ALTER TABLE example_db.my_table ADD COLUMN new_col INT KEY DEFAULT \"0\" AFTER key_1",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterAddColumn {
					t.Errorf("expected AlterAddColumn, got %v", a.Type)
				}
				if a.Column.Name != "new_col" {
					t.Errorf("expected column name 'new_col', got %q", a.Column.Name)
				}
				if a.After != "key_1" {
					t.Errorf("expected After='key_1', got %q", a.After)
				}
			},
		},
		{
			sql: "ALTER TABLE example_db.my_table ADD COLUMN new_col INT KEY DEFAULT \"0\" FIRST",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if !a.First {
					t.Error("expected First=true")
				}
			},
		},
		{
			// Multi-column ADD COLUMN (...)
			sql: `ALTER TABLE example_db.my_table ADD COLUMN (new_col1 INT SUM DEFAULT "0", new_col2 INT SUM DEFAULT "0")`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				if len(stmt.Actions) != 1 {
					t.Fatalf("expected 1 action, got %d", len(stmt.Actions))
				}
				a := stmt.Actions[0]
				if a.Type != ast.AlterAddColumn {
					t.Errorf("expected AlterAddColumn, got %v", a.Type)
				}
			},
		},
	}
	runAlterTableTests(t, cases)
}

// ---- DROP COLUMN ----

func TestAlterTableDropColumn(t *testing.T) {
	cases := []alterTableTestCase{
		{
			sql: "ALTER TABLE t DROP COLUMN c",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterDropColumn {
					t.Errorf("expected AlterDropColumn, got %v", a.Type)
				}
				if a.ColumnName != "c" {
					t.Errorf("expected ColumnName='c', got %q", a.ColumnName)
				}
			},
		},
		{
			sql: "ALTER TABLE example_db.my_table DROP COLUMN col1",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterDropColumn {
					t.Errorf("expected AlterDropColumn")
				}
				if a.ColumnName != "col1" {
					t.Errorf("expected ColumnName='col1', got %q", a.ColumnName)
				}
			},
		},
	}
	runAlterTableTests(t, cases)
}

// ---- MODIFY COLUMN ----

func TestAlterTableModifyColumn(t *testing.T) {
	cases := []alterTableTestCase{
		{
			sql: "ALTER TABLE t MODIFY COLUMN c BIGINT",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterModifyColumn {
					t.Errorf("expected AlterModifyColumn, got %v", a.Type)
				}
				if a.Column == nil || a.Column.Name != "c" {
					t.Errorf("expected column name 'c'")
				}
			},
		},
		{
			sql: `ALTER TABLE example_db.my_table MODIFY COLUMN col1 BIGINT KEY DEFAULT "1" AFTER col2`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterModifyColumn {
					t.Errorf("expected AlterModifyColumn")
				}
				if a.After != "col2" {
					t.Errorf("expected After='col2', got %q", a.After)
				}
			},
		},
		{
			sql: `ALTER TABLE example_db.my_table MODIFY COLUMN val1 VARCHAR(64) REPLACE DEFAULT "abc"`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterModifyColumn {
					t.Errorf("expected AlterModifyColumn")
				}
				if a.Column.Name != "val1" {
					t.Errorf("expected column name 'val1'")
				}
			},
		},
		{
			sql: `ALTER TABLE example_db.my_table MODIFY COLUMN k3 VARCHAR(50) KEY NULL COMMENT 'to 50'`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterModifyColumn {
					t.Errorf("expected AlterModifyColumn")
				}
				if a.Column.Comment != "to 50" {
					t.Errorf("expected comment 'to 50', got %q", a.Column.Comment)
				}
			},
		},
	}
	runAlterTableTests(t, cases)
}

// ---- RENAME ----

func TestAlterTableRename(t *testing.T) {
	cases := []alterTableTestCase{
		{
			sql: "ALTER TABLE t RENAME TO new_name",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterRenameTable {
					t.Errorf("expected AlterRenameTable, got %v", a.Type)
				}
				if a.NewTableName == nil || a.NewTableName.Parts[0] != "new_name" {
					t.Errorf("expected new_name, got %v", a.NewTableName)
				}
			},
		},
		{
			// RENAME without TO
			sql: "ALTER TABLE table1 RENAME table2",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterRenameTable {
					t.Errorf("expected AlterRenameTable, got %v", a.Type)
				}
			},
		},
		{
			sql: "ALTER TABLE example_table RENAME ROLLUP rollup1 rollup2",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterRenameRollup {
					t.Errorf("expected AlterRenameRollup, got %v", a.Type)
				}
				if a.ColumnName != "rollup1" {
					t.Errorf("expected old name 'rollup1', got %q", a.ColumnName)
				}
				if a.NewName != "rollup2" {
					t.Errorf("expected new name 'rollup2', got %q", a.NewName)
				}
			},
		},
		{
			sql: "ALTER TABLE example_table RENAME PARTITION p1 p2",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterRenamePartition {
					t.Errorf("expected AlterRenamePartition, got %v", a.Type)
				}
				if a.ColumnName != "p1" {
					t.Errorf("expected old name 'p1', got %q", a.ColumnName)
				}
				if a.NewName != "p2" {
					t.Errorf("expected new name 'p2', got %q", a.NewName)
				}
			},
		},
		{
			sql: "ALTER TABLE example_table RENAME COLUMN c1 c2",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterRenameColumn {
					t.Errorf("expected AlterRenameColumn, got %v", a.Type)
				}
				if a.ColumnName != "c1" || a.NewName != "c2" {
					t.Errorf("expected c1->c2, got %q->%q", a.ColumnName, a.NewName)
				}
			},
		},
	}
	runAlterTableTests(t, cases)
}

// ---- PARTITION ----

func TestAlterTablePartition(t *testing.T) {
	cases := []alterTableTestCase{
		{
			sql: `ALTER TABLE t ADD PARTITION p1 VALUES LESS THAN ('2024-01-01')`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterAddPartition {
					t.Errorf("expected AlterAddPartition, got %v", a.Type)
				}
				if a.Partition == nil {
					t.Fatal("expected partition item")
				}
				if a.Partition.Name != "p1" {
					t.Errorf("expected partition name 'p1', got %q", a.Partition.Name)
				}
			},
		},
		{
			sql: `ALTER TABLE example_db.my_table ADD PARTITION p1 VALUES LESS THAN ("2014-01-01")`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterAddPartition {
					t.Errorf("expected AlterAddPartition")
				}
				if len(a.Partition.Values) == 0 {
					t.Error("expected partition values")
				}
			},
		},
		{
			sql: `ALTER TABLE example_db.my_table ADD PARTITION p1 VALUES LESS THAN ("2015-01-01") DISTRIBUTED BY HASH(k1) BUCKETS 20`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterAddPartition {
					t.Errorf("expected AlterAddPartition")
				}
				if a.PartitionDist == nil {
					t.Error("expected partition distribution")
				} else if a.PartitionDist.Buckets != 20 {
					t.Errorf("expected 20 buckets, got %d", a.PartitionDist.Buckets)
				}
			},
		},
		{
			sql: `ALTER TABLE example_db.my_table ADD PARTITION p1 VALUES LESS THAN ("2015-01-01") ("replication_num" = "1")`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterAddPartition {
					t.Errorf("expected AlterAddPartition")
				}
				if len(a.PartitionProps) == 0 {
					t.Error("expected partition properties")
				}
			},
		},
		{
			sql: "ALTER TABLE t DROP PARTITION p1",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterDropPartition {
					t.Errorf("expected AlterDropPartition, got %v", a.Type)
				}
				if a.PartitionName != "p1" {
					t.Errorf("expected partition name 'p1', got %q", a.PartitionName)
				}
			},
		},
		{
			// Multiple DROP PARTITION actions in one statement
			sql: "ALTER TABLE example_db.my_table DROP PARTITION p1, DROP PARTITION p2, DROP PARTITION p3",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				if len(stmt.Actions) != 3 {
					t.Fatalf("expected 3 actions, got %d", len(stmt.Actions))
				}
				for i, a := range stmt.Actions {
					if a.Type != ast.AlterDropPartition {
						t.Errorf("action %d: expected AlterDropPartition", i)
					}
				}
			},
		},
		{
			// MODIFY PARTITION single
			sql: `ALTER TABLE example_db.my_table MODIFY PARTITION p1 SET("replication_num" = "1")`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterModifyPartition {
					t.Errorf("expected AlterModifyPartition, got %v", a.Type)
				}
				if a.PartitionName != "p1" {
					t.Errorf("expected partition name 'p1', got %q", a.PartitionName)
				}
				if len(a.Properties) == 0 {
					t.Error("expected properties")
				}
			},
		},
		{
			// MODIFY PARTITION list
			sql: `ALTER TABLE example_db.my_table MODIFY PARTITION (p1, p2, p4) SET("replication_num" = "1")`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterModifyPartition {
					t.Errorf("expected AlterModifyPartition")
				}
				if len(a.PartitionList) != 3 {
					t.Errorf("expected 3 partitions in list, got %d", len(a.PartitionList))
				}
			},
		},
		{
			// MODIFY PARTITION (*)
			sql: `ALTER TABLE example_db.my_table MODIFY PARTITION (*) SET("storage_medium" = "HDD")`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterModifyPartition {
					t.Errorf("expected AlterModifyPartition")
				}
				if !a.PartitionStar {
					t.Error("expected PartitionStar=true")
				}
			},
		},
		{
			// Fixed range partition [lower, upper)
			sql: `ALTER TABLE example_db.my_table ADD PARTITION p1 VALUES [("2014-01-01"), ("2014-02-01"))`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterAddPartition {
					t.Errorf("expected AlterAddPartition")
				}
			},
		},
	}
	runAlterTableTests(t, cases)
}

// ---- ROLLUP ----

func TestAlterTableRollup(t *testing.T) {
	cases := []alterTableTestCase{
		{
			sql: "ALTER TABLE t ADD ROLLUP r (c1, c2)",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterAddRollup {
					t.Errorf("expected AlterAddRollup, got %v", a.Type)
				}
				if a.Rollup == nil {
					t.Fatal("expected rollup def")
				}
				if a.Rollup.Name != "r" {
					t.Errorf("expected rollup name 'r', got %q", a.Rollup.Name)
				}
				if len(a.Rollup.Columns) != 2 {
					t.Errorf("expected 2 columns, got %d", len(a.Rollup.Columns))
				}
			},
		},
		{
			sql: "ALTER TABLE example_db.my_table\nADD ROLLUP example_rollup_index(k1, k3, v1, v2)",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterAddRollup {
					t.Errorf("expected AlterAddRollup")
				}
				if a.Rollup.Name != "example_rollup_index" {
					t.Errorf("unexpected rollup name: %q", a.Rollup.Name)
				}
			},
		},
		{
			sql: "ALTER TABLE example_db.my_table\nADD ROLLUP example_rollup_index2 (k1, v1)\nFROM example_rollup_index",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterAddRollup {
					t.Errorf("expected AlterAddRollup")
				}
				if a.RollupName != "example_rollup_index" {
					t.Errorf("expected base rollup 'example_rollup_index', got %q", a.RollupName)
				}
			},
		},
		{
			sql: "ALTER TABLE example_db.my_table\nADD ROLLUP example_rollup_index(k1, k3, v1)\nPROPERTIES(\"timeout\" = \"3600\")",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterAddRollup {
					t.Errorf("expected AlterAddRollup")
				}
				if len(a.Rollup.Properties) == 0 {
					t.Error("expected rollup properties")
				}
			},
		},
		{
			sql: "ALTER TABLE example_db.my_table\nDROP ROLLUP example_rollup_index2",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterDropRollup {
					t.Errorf("expected AlterDropRollup, got %v", a.Type)
				}
				if a.RollupName != "example_rollup_index2" {
					t.Errorf("expected rollup name 'example_rollup_index2'")
				}
			},
		},
		{
			// Multiple DROP ROLLUP in one statement (comma-separated names)
			sql: "ALTER TABLE example_db.my_table\nDROP ROLLUP example_rollup_index2,example_rollup_index3",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				if len(stmt.Actions) != 2 {
					t.Fatalf("expected 2 actions, got %d", len(stmt.Actions))
				}
				for i, a := range stmt.Actions {
					if a.Type != ast.AlterDropRollup {
						t.Errorf("action %d: expected AlterDropRollup, got %v", i, a.Type)
					}
				}
			},
		},
	}
	runAlterTableTests(t, cases)
}

// ---- SET PROPERTIES ----

func TestAlterTableSetProperties(t *testing.T) {
	cases := []alterTableTestCase{
		{
			sql: `ALTER TABLE t SET ("replication_num"="3")`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterSetProperties {
					t.Errorf("expected AlterSetProperties, got %v", a.Type)
				}
				if len(a.Properties) != 1 {
					t.Errorf("expected 1 property, got %d", len(a.Properties))
				}
				if a.Properties[0].Key != "replication_num" {
					t.Errorf("unexpected key: %q", a.Properties[0].Key)
				}
			},
		},
		{
			sql: `ALTER TABLE example_db.my_table SET ("bloom_filter_columns"="k1,k2,k3")`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterSetProperties {
					t.Errorf("expected AlterSetProperties")
				}
			},
		},
		{
			sql: `ALTER TABLE example_db.my_table SET ("dynamic_partition.enable" = "true", "dynamic_partition.time_unit" = "DAY", "dynamic_partition.end" = "3", "dynamic_partition.prefix" = "p", "dynamic_partition.buckets" = "32")`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterSetProperties {
					t.Errorf("expected AlterSetProperties")
				}
				if len(a.Properties) != 5 {
					t.Errorf("expected 5 properties, got %d", len(a.Properties))
				}
			},
		},
	}
	runAlterTableTests(t, cases)
}

// ---- MODIFY COMMENT / DISTRIBUTION / ENGINE ----

func TestAlterTableModifyMisc(t *testing.T) {
	cases := []alterTableTestCase{
		{
			sql: `ALTER TABLE example_db.my_table MODIFY COMMENT "new comment"`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterModifyComment {
					t.Errorf("expected AlterModifyComment, got %v", a.Type)
				}
				if a.Comment != "new comment" {
					t.Errorf("expected comment 'new comment', got %q", a.Comment)
				}
			},
		},
		{
			sql: `ALTER TABLE example_db.my_table MODIFY DISTRIBUTION DISTRIBUTED BY HASH(k1) BUCKETS 50`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterModifyDistribution {
					t.Errorf("expected AlterModifyDistribution, got %v", a.Type)
				}
				if a.Distribution == nil {
					t.Fatal("expected distribution")
				}
				if a.Distribution.Buckets != 50 {
					t.Errorf("expected 50 buckets, got %d", a.Distribution.Buckets)
				}
			},
		},
	}
	runAlterTableTests(t, cases)
}

// ---- ENABLE FEATURE ----

func TestAlterTableEnableFeature(t *testing.T) {
	cases := []alterTableTestCase{
		{
			sql: `ALTER TABLE example_db.my_table ENABLE FEATURE "BATCH_DELETE"`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterEnableFeature {
					t.Errorf("expected AlterEnableFeature, got %v", a.Type)
				}
				if a.FeatureName != "BATCH_DELETE" {
					t.Errorf("expected feature 'BATCH_DELETE', got %q", a.FeatureName)
				}
			},
		},
		{
			sql: `ALTER TABLE example_db.my_table ENABLE FEATURE "SEQUENCE_LOAD" WITH PROPERTIES ("function_column.sequence_type" = "Date")`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterEnableFeature {
					t.Errorf("expected AlterEnableFeature")
				}
				if a.FeatureName != "SEQUENCE_LOAD" {
					t.Errorf("expected feature 'SEQUENCE_LOAD', got %q", a.FeatureName)
				}
				if len(a.Properties) == 0 {
					t.Error("expected properties")
				}
			},
		},
	}
	runAlterTableTests(t, cases)
}

// ---- MULTIPLE ACTIONS ----

func TestAlterTableMultipleActions(t *testing.T) {
	cases := []alterTableTestCase{
		{
			sql: "ALTER TABLE t ADD COLUMN a INT, ADD COLUMN b VARCHAR(50)",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				if len(stmt.Actions) != 2 {
					t.Fatalf("expected 2 actions, got %d", len(stmt.Actions))
				}
				for i, a := range stmt.Actions {
					if a.Type != ast.AlterAddColumn {
						t.Errorf("action %d: expected AlterAddColumn, got %v", i, a.Type)
					}
				}
			},
		},
		{
			sql: `ALTER TABLE example_db.my_table ADD COLUMN col INT DEFAULT "0" AFTER v_1, ORDER BY (k_2, k_1, v_3, v_2, v_1, col)`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				if len(stmt.Actions) != 2 {
					t.Fatalf("expected 2 actions, got %d", len(stmt.Actions))
				}
				if stmt.Actions[0].Type != ast.AlterAddColumn {
					t.Errorf("action 0: expected AlterAddColumn")
				}
				if stmt.Actions[1].Type != ast.AlterOrderBy {
					t.Errorf("action 1: expected AlterOrderBy, got %v", stmt.Actions[1].Type)
				}
				if len(stmt.Actions[1].OrderByColumns) != 6 {
					t.Errorf("expected 6 ORDER BY columns, got %d", len(stmt.Actions[1].OrderByColumns))
				}
			},
		},
		{
			sql: `ALTER TABLE example_db.my_table MODIFY COLUMN k1 COMMENT "k1", MODIFY COLUMN k2 COMMENT "k2"`,
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				if len(stmt.Actions) != 2 {
					t.Fatalf("expected 2 actions, got %d", len(stmt.Actions))
				}
				for i, a := range stmt.Actions {
					if a.Type != ast.AlterModifyColumn {
						t.Errorf("action %d: expected AlterModifyColumn", i)
					}
				}
			},
		},
	}
	runAlterTableTests(t, cases)
}

// ---- REPLACE PARTITION ----

func TestAlterTableReplacePartition(t *testing.T) {
	cases := []alterTableTestCase{
		{
			sql: "ALTER TABLE tbl1 REPLACE WITH TABLE tbl2",
			// This is REPLACE WITH TABLE, not REPLACE PARTITION — should produce AlterRaw
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				if len(stmt.Actions) == 0 {
					t.Fatal("expected at least 1 action")
				}
				// REPLACE WITH TABLE falls through to raw
			},
		},
		{
			sql: "ALTER TABLE tbl1 REPLACE WITH TABLE tbl2 PROPERTIES('swap' = 'true')",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				if len(stmt.Actions) == 0 {
					t.Fatal("expected at least 1 action")
				}
			},
		},
	}
	runAlterTableTests(t, cases)
}

// ---- ORDER BY ----

func TestAlterTableOrderBy(t *testing.T) {
	cases := []alterTableTestCase{
		{
			sql: "ALTER TABLE example_db.my_table ORDER BY (k_2, k_1, v_3, v_2, v_1)",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				a := stmt.Actions[0]
				if a.Type != ast.AlterOrderBy {
					t.Errorf("expected AlterOrderBy, got %v", a.Type)
				}
				if len(a.OrderByColumns) != 5 {
					t.Errorf("expected 5 columns, got %d", len(a.OrderByColumns))
				}
			},
		},
	}
	runAlterTableTests(t, cases)
}

// ---- TABLE NAME PARSING ----

func TestAlterTableName(t *testing.T) {
	cases := []alterTableTestCase{
		{
			sql: "ALTER TABLE db.tbl ADD COLUMN c INT",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				if stmt.Name == nil {
					t.Fatal("expected table name")
				}
				if len(stmt.Name.Parts) != 2 {
					t.Errorf("expected 2-part name, got %d", len(stmt.Name.Parts))
				}
				if stmt.Name.Parts[0] != "db" || stmt.Name.Parts[1] != "tbl" {
					t.Errorf("unexpected name parts: %v", stmt.Name.Parts)
				}
			},
		},
		{
			sql: "ALTER TABLE catalog.db.tbl DROP COLUMN c",
			check: func(t *testing.T, stmt *ast.AlterTableStmt) {
				if len(stmt.Name.Parts) != 3 {
					t.Errorf("expected 3-part name, got %d", len(stmt.Name.Parts))
				}
			},
		},
	}
	runAlterTableTests(t, cases)
}

// ---- FULL LEGACY CORPUS ----

func TestAlterTableLegacyColumn(t *testing.T) {
	cases := []alterTableTestCase{
		{sql: `ALTER TABLE example_db.my_table ADD COLUMN new_col INT KEY DEFAULT "0" AFTER key_1`},
		{sql: `ALTER TABLE example_db.my_table ADD COLUMN new_col INT DEFAULT "0" AFTER value_1`},
		{sql: `ALTER TABLE example_db.my_table ADD COLUMN new_col INT SUM DEFAULT "0" AFTER value_1`},
		{sql: `ALTER TABLE example_db.my_table ADD COLUMN new_col INT KEY DEFAULT "0" FIRST`},
		{sql: `ALTER TABLE example_db.my_table ADD COLUMN (new_col1 INT SUM DEFAULT "0", new_col2 INT SUM DEFAULT "0")`},
		{sql: `ALTER TABLE example_db.my_table ADD COLUMN (new_col1 INT KEY DEFAULT "0", new_col2 INT DEFAULT "0")`},
		{sql: `ALTER TABLE example_db.my_table DROP COLUMN col1`},
		{sql: `ALTER TABLE example_db.my_table MODIFY COLUMN col1 BIGINT KEY DEFAULT "1" AFTER col2`},
		{sql: `ALTER TABLE example_db.my_table MODIFY COLUMN val1 VARCHAR(64) REPLACE DEFAULT "abc"`},
		{sql: `ALTER TABLE example_db.my_table MODIFY COLUMN k3 VARCHAR(50) KEY NULL COMMENT 'to 50'`},
		{sql: `ALTER TABLE example_db.my_table ORDER BY (k_2, k_1, v_3, v_2, v_1)`},
		{sql: `ALTER TABLE example_db.my_table ADD COLUMN col INT DEFAULT "0" AFTER v_1, ORDER BY (k_2, k_1, v_3, v_2, v_1, col)`},
	}
	runAlterTableTests(t, cases)
}

func TestAlterTableLegacyPartition(t *testing.T) {
	cases := []alterTableTestCase{
		{sql: `ALTER TABLE example_db.my_table ADD PARTITION p1 VALUES LESS THAN ("2014-01-01")`},
		{sql: `ALTER TABLE example_db.my_table ADD PARTITION p1 VALUES LESS THAN ("2015-01-01") DISTRIBUTED BY HASH(k1) BUCKETS 20`},
		{sql: `ALTER TABLE example_db.my_table ADD PARTITION p1 VALUES LESS THAN ("2015-01-01") ("replication_num" = "1")`},
		{sql: `ALTER TABLE example_db.my_table MODIFY PARTITION p1 SET("replication_num" = "1")`},
		{sql: `ALTER TABLE example_db.my_table MODIFY PARTITION (p1, p2, p4) SET("replication_num" = "1")`},
		{sql: `ALTER TABLE example_db.my_table MODIFY PARTITION (*) SET("storage_medium" = "HDD")`},
		{sql: `ALTER TABLE example_db.my_table DROP PARTITION p1`},
		{sql: `ALTER TABLE example_db.my_table DROP PARTITION p1, DROP PARTITION p2, DROP PARTITION p3`},
		{sql: `ALTER TABLE example_db.my_table ADD PARTITION p1 VALUES [("2014-01-01"), ("2014-02-01"))`},
	}
	runAlterTableTests(t, cases)
}

func TestAlterTableLegacyProperty(t *testing.T) {
	cases := []alterTableTestCase{
		{sql: `ALTER TABLE example_db.my_table SET ("bloom_filter_columns"="k1,k2,k3")`},
		{sql: `ALTER TABLE example_db.my_table SET ("colocate_with" = "t1")`},
		{sql: `ALTER TABLE example_db.my_table SET ("distribution_type" = "random")`},
		{sql: `ALTER TABLE example_db.my_table SET ("dynamic_partition.enable" = "false")`},
		{sql: `ALTER TABLE example_db.my_table SET ("dynamic_partition.enable" = "true", "dynamic_partition.time_unit" = "DAY", "dynamic_partition.end" = "3", "dynamic_partition.prefix" = "p", "dynamic_partition.buckets" = "32")`},
		{sql: `ALTER TABLE example_db.my_table SET ("in_memory" = "false")`},
		{sql: `ALTER TABLE example_db.my_table ENABLE FEATURE "BATCH_DELETE"`},
		{sql: `ALTER TABLE example_db.my_table ENABLE FEATURE "SEQUENCE_LOAD" WITH PROPERTIES ("function_column.sequence_type" = "Date")`},
		{sql: `ALTER TABLE example_db.my_table MODIFY DISTRIBUTION DISTRIBUTED BY HASH(k1) BUCKETS 50`},
		{sql: `ALTER TABLE example_db.my_table MODIFY COMMENT "new comment"`},
		{sql: `ALTER TABLE example_db.my_table MODIFY COLUMN k1 COMMENT "k1", MODIFY COLUMN k2 COMMENT "k2"`},
		{sql: `ALTER TABLE example_db.mysql_table SET ("replication_num" = "2")`},
		{sql: `ALTER TABLE example_db.mysql_table SET ("default.replication_num" = "2")`},
		{sql: `ALTER TABLE example_db.mysql_table SET ("replication_allocation" = "tag.location.default: 1")`},
		{sql: `ALTER TABLE example_db.mysql_table SET ("default.replication_allocation" = "tag.location.default: 1")`},
		{sql: `ALTER TABLE example_db.mysql_table SET ("light_schema_change" = "true")`},
		{sql: `ALTER TABLE create_table_not_have_policy SET ("storage_policy" = "created_create_table_alter_policy")`},
		{sql: `ALTER TABLE create_table_partition MODIFY PARTITION (*) SET("storage_policy"="created_create_table_partition_alter_policy")`},
	}
	runAlterTableTests(t, cases)
}

func TestAlterTableLegacyRename(t *testing.T) {
	cases := []alterTableTestCase{
		{sql: `ALTER TABLE table1 RENAME table2`},
		{sql: `ALTER TABLE example_table RENAME ROLLUP rollup1 rollup2`},
		{sql: `ALTER TABLE example_table RENAME PARTITION p1 p2`},
		{sql: `ALTER TABLE example_table RENAME COLUMN c1 c2`},
	}
	runAlterTableTests(t, cases)
}

func TestAlterTableLegacyReplace(t *testing.T) {
	cases := []alterTableTestCase{
		{sql: `ALTER TABLE tbl1 REPLACE WITH TABLE tbl2`},
		{sql: `ALTER TABLE tbl1 REPLACE WITH TABLE tbl2 PROPERTIES('swap' = 'true')`},
		{sql: `ALTER TABLE tbl1 REPLACE WITH TABLE tbl2 PROPERTIES('swap' = 'false')`},
	}
	runAlterTableTests(t, cases)
}

func TestAlterTableLegacyRollup(t *testing.T) {
	cases := []alterTableTestCase{
		{sql: "ALTER TABLE example_db.my_table\nADD ROLLUP example_rollup_index(k1, k3, v1, v2)"},
		{sql: "ALTER TABLE example_db.my_table\nADD ROLLUP example_rollup_index2 (k1, v1)\nFROM example_rollup_index"},
		{sql: "ALTER TABLE example_db.my_table\nADD ROLLUP example_rollup_index(k1, k3, v1)\nPROPERTIES(\"timeout\" = \"3600\")"},
		{sql: "ALTER TABLE example_db.my_table\nDROP ROLLUP example_rollup_index2"},
		{sql: "ALTER TABLE example_db.my_table\nDROP ROLLUP example_rollup_index2,example_rollup_index3"},
	}
	runAlterTableTests(t, cases)
}
