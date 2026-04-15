package analysis_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/analysis"
	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Classify — nil / zero cases
// ---------------------------------------------------------------------------

func TestClassify_Nil(t *testing.T) {
	if got := analysis.Classify(nil); got != analysis.CategoryUnknown {
		t.Errorf("Classify(nil) = %v, want CategoryUnknown", got)
	}
}

func TestClassify_NonStmtNode(t *testing.T) {
	// ObjectName is a valid Node but not a statement — should return Unknown.
	node := &ast.ObjectName{}
	if got := analysis.Classify(node); got != analysis.CategoryUnknown {
		t.Errorf("Classify(*ObjectName) = %v, want CategoryUnknown", got)
	}
}

// ---------------------------------------------------------------------------
// ClassifySQL — empty / no-parse cases
// ---------------------------------------------------------------------------

func TestClassifySQL_Empty(t *testing.T) {
	if got := analysis.ClassifySQL(""); got != analysis.CategoryUnknown {
		t.Errorf("ClassifySQL(\"\") = %v, want CategoryUnknown", got)
	}
}

func TestClassifySQL_Whitespace(t *testing.T) {
	if got := analysis.ClassifySQL("   \t\n  "); got != analysis.CategoryUnknown {
		t.Errorf("ClassifySQL(whitespace) = %v, want CategoryUnknown", got)
	}
}

func TestClassifySQL_CommentOnly(t *testing.T) {
	if got := analysis.ClassifySQL("-- just a comment"); got != analysis.CategoryUnknown {
		t.Errorf("ClassifySQL(comment) = %v, want CategoryUnknown", got)
	}
}

// ---------------------------------------------------------------------------
// CategorySelect
// ---------------------------------------------------------------------------

func TestClassifySQL_Select_Simple(t *testing.T) {
	assertCategory(t, "SELECT 1", analysis.CategorySelect)
}

func TestClassifySQL_Select_FromTable(t *testing.T) {
	assertCategory(t, "SELECT id, name FROM users WHERE id = 42", analysis.CategorySelect)
}

func TestClassifySQL_Select_UnionAll(t *testing.T) {
	assertCategory(t, "SELECT 1 UNION ALL SELECT 2", analysis.CategorySelect)
}

func TestClassifySQL_Select_Intersect(t *testing.T) {
	assertCategory(t, "SELECT id FROM a INTERSECT SELECT id FROM b", analysis.CategorySelect)
}

func TestClassifySQL_Select_Except(t *testing.T) {
	assertCategory(t, "SELECT id FROM a EXCEPT SELECT id FROM b", analysis.CategorySelect)
}

func TestClassifySQL_Select_WithCTE(t *testing.T) {
	assertCategory(t,
		`WITH cte AS (SELECT 1 AS n) SELECT n FROM cte`,
		analysis.CategorySelect)
}

func TestClassifySQL_Select_Subquery(t *testing.T) {
	assertCategory(t,
		`SELECT * FROM (SELECT 1 AS x) AS sub`,
		analysis.CategorySelect)
}

// ---------------------------------------------------------------------------
// CategoryDML
// ---------------------------------------------------------------------------

func TestClassifySQL_Insert_Values(t *testing.T) {
	assertCategory(t, "INSERT INTO t VALUES (1)", analysis.CategoryDML)
}

func TestClassifySQL_Insert_Select(t *testing.T) {
	assertCategory(t, "INSERT INTO t SELECT * FROM src", analysis.CategoryDML)
}

func TestClassifySQL_Insert_Multi(t *testing.T) {
	assertCategory(t,
		`INSERT ALL INTO t1 VALUES (1) INTO t2 VALUES (2) SELECT 1`,
		analysis.CategoryDML)
}

func TestClassifySQL_Update(t *testing.T) {
	assertCategory(t, "UPDATE t SET c = 1 WHERE id = 1", analysis.CategoryDML)
}

func TestClassifySQL_Delete(t *testing.T) {
	assertCategory(t, "DELETE FROM t WHERE id = 1", analysis.CategoryDML)
}

func TestClassifySQL_Merge(t *testing.T) {
	assertCategory(t,
		`MERGE INTO tgt USING src ON tgt.id = src.id
		 WHEN MATCHED THEN UPDATE SET tgt.v = src.v
		 WHEN NOT MATCHED THEN INSERT (id, v) VALUES (src.id, src.v)`,
		analysis.CategoryDML)
}

// ---------------------------------------------------------------------------
// CategoryDDL — tables
// ---------------------------------------------------------------------------

func TestClassifySQL_CreateTable(t *testing.T) {
	assertCategory(t, "CREATE TABLE t (id INT)", analysis.CategoryDDL)
}

func TestClassifySQL_CreateTableOrReplace(t *testing.T) {
	assertCategory(t, "CREATE OR REPLACE TABLE t (id INT, name VARCHAR)", analysis.CategoryDDL)
}

func TestClassifySQL_AlterTable_AddColumn(t *testing.T) {
	assertCategory(t, "ALTER TABLE t ADD COLUMN c INT", analysis.CategoryDDL)
}

func TestClassifySQL_DropTable(t *testing.T) {
	assertCategory(t, "DROP TABLE IF EXISTS t", analysis.CategoryDDL)
}

func TestClassifySQL_UndropTable(t *testing.T) {
	assertCategory(t, "UNDROP TABLE t", analysis.CategoryDDL)
}

// ---------------------------------------------------------------------------
// CategoryDDL — databases
// ---------------------------------------------------------------------------

func TestClassifySQL_CreateDatabase(t *testing.T) {
	assertCategory(t, "CREATE DATABASE mydb", analysis.CategoryDDL)
}

func TestClassifySQL_AlterDatabase(t *testing.T) {
	assertCategory(t, "ALTER DATABASE mydb RENAME TO newdb", analysis.CategoryDDL)
}

func TestClassifySQL_DropDatabase(t *testing.T) {
	assertCategory(t, "DROP DATABASE mydb", analysis.CategoryDDL)
}

func TestClassifySQL_UndropDatabase(t *testing.T) {
	assertCategory(t, "UNDROP DATABASE mydb", analysis.CategoryDDL)
}

// ---------------------------------------------------------------------------
// CategoryDDL — schemas
// ---------------------------------------------------------------------------

func TestClassifySQL_CreateSchema(t *testing.T) {
	assertCategory(t, "CREATE SCHEMA myschema", analysis.CategoryDDL)
}

func TestClassifySQL_DropSchema(t *testing.T) {
	assertCategory(t, "DROP SCHEMA myschema", analysis.CategoryDDL)
}

func TestClassifySQL_UndropSchema(t *testing.T) {
	assertCategory(t, "UNDROP SCHEMA myschema", analysis.CategoryDDL)
}

// ---------------------------------------------------------------------------
// CategoryDDL — views
// ---------------------------------------------------------------------------

func TestClassifySQL_CreateView(t *testing.T) {
	assertCategory(t, "CREATE VIEW v AS SELECT 1", analysis.CategoryDDL)
}

func TestClassifySQL_CreateMaterializedView(t *testing.T) {
	assertCategory(t, "CREATE MATERIALIZED VIEW mv AS SELECT id FROM t", analysis.CategoryDDL)
}

func TestClassifySQL_AlterView(t *testing.T) {
	assertCategory(t, "ALTER VIEW v RENAME TO v2", analysis.CategoryDDL)
}

func TestClassifySQL_DropView(t *testing.T) {
	assertCategory(t, "DROP VIEW v", analysis.CategoryDDL)
}

// ---------------------------------------------------------------------------
// Category.String()
// ---------------------------------------------------------------------------

func TestCategory_String(t *testing.T) {
	cases := []struct {
		cat  analysis.Category
		want string
	}{
		{analysis.CategoryUnknown, "Unknown"},
		{analysis.CategorySelect, "SELECT"},
		{analysis.CategoryDML, "DML"},
		{analysis.CategoryDDL, "DDL"},
		{analysis.CategoryShow, "SHOW"},
		{analysis.CategoryDescribe, "DESCRIBE"},
		{analysis.CategoryOther, "Other"},
	}
	for _, tc := range cases {
		if got := tc.cat.String(); got != tc.want {
			t.Errorf("Category(%d).String() = %q, want %q", tc.cat, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func assertCategory(t *testing.T, sql string, want analysis.Category) {
	t.Helper()
	got := analysis.ClassifySQL(sql)
	if got != want {
		t.Errorf("ClassifySQL(%q) = %v, want %v", sql, got, want)
	}
}
