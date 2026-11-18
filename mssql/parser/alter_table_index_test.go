package parser

import (
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// TestParseAlterTableAddIndex pins ALTER TABLE's ADD INDEX, the table_index
// of the ADD list: ScriptDOM reads it as the tableElement a CREATE TABLE body
// takes (inlineIndexTableDefinition, TSql160.g:28735), filterClause included.
// Every accepted form here was run on the pinned SQL Server 2022, and the
// refused ones are the rows mssql/parser/testdata/compat.tsv records.
func TestParseAlterTableAddIndex(t *testing.T) {
	t.Run("bare", func(t *testing.T) {
		// SQL Server's parser takes an inline index with no key list.
		result := ParseAndCheck(t, "ALTER TABLE a ADD INDEX ix")
		idx := addIndexAction(t, result, "ALTER TABLE a ADD INDEX ix")
		if idx.Name != "ix" {
			t.Errorf("Name = %q, want ix", idx.Name)
		}
		if idx.Columns != nil || idx.WhereClause != nil {
			t.Errorf("expected no columns and no filter, got %v and %v", idx.Columns, idx.WhereClause)
		}
	})

	t.Run("columns_and_type", func(t *testing.T) {
		sql := "ALTER TABLE a ADD INDEX ix UNIQUE NONCLUSTERED (a ASC, b DESC) INCLUDE (c)"
		idx := addIndexAction(t, ParseAndCheck(t, sql), sql)
		if !idx.Unique {
			t.Error("expected Unique")
		}
		if idx.Clustered == nil || *idx.Clustered {
			t.Errorf("Clustered = %v, want false", idx.Clustered)
		}
		if idx.Columns == nil || idx.Columns.Len() != 2 {
			t.Fatalf("Columns = %v, want 2", idx.Columns)
		}
		first, ok := idx.Columns.Items[0].(*ast.IndexColumn)
		if !ok || first.Name != "a" || first.SortDir != ast.SortAsc {
			t.Errorf("first column = %v, want a ASC", idx.Columns.Items[0])
		}
		second, ok := idx.Columns.Items[1].(*ast.IndexColumn)
		if !ok || second.Name != "b" || second.SortDir != ast.SortDesc {
			t.Errorf("second column = %v, want b DESC", idx.Columns.Items[1])
		}
		if idx.IncludeCols == nil || idx.IncludeCols.Len() != 1 {
			t.Errorf("IncludeCols = %v, want c", idx.IncludeCols)
		}
	})

	t.Run("filter", func(t *testing.T) {
		sql := "ALTER TABLE a ADD INDEX ix (a) WHERE a IS NOT NULL AND a IN (1, 2)"
		idx := addIndexAction(t, ParseAndCheck(t, sql), sql)
		and, ok := idx.WhereClause.(*ast.BinaryExpr)
		if !ok || and.Op != ast.BinOpAnd {
			t.Fatalf("WhereClause = %T, want an AND", idx.WhereClause)
		}
		if isnull, ok := and.Left.(*ast.IsExpr); !ok || isnull.TestType != ast.IsNotNull {
			t.Errorf("left = %v, want IS NOT NULL", and.Left)
		}
		if in, ok := and.Right.(*ast.InExpr); !ok || in.List == nil || in.List.Len() != 2 {
			t.Errorf("right = %v, want IN with two values", and.Right)
		}
	})

	t.Run("filter_distinct", func(t *testing.T) {
		for _, tt := range []struct {
			sql  string
			test ast.IsTestType
		}{
			{"ALTER TABLE a ADD INDEX ix (a) WHERE a IS DISTINCT FROM 1", ast.IsDistinctFrom},
			{"ALTER TABLE a ADD INDEX ix (a) WHERE a IS NOT DISTINCT FROM NULL", ast.IsNotDistinctFrom},
		} {
			t.Run(tt.sql, func(t *testing.T) {
				idx := addIndexAction(t, ParseAndCheck(t, tt.sql), tt.sql)
				distinct, ok := idx.WhereClause.(*ast.IsExpr)
				if !ok {
					t.Fatalf("WhereClause = %T, want IS", idx.WhereClause)
				}
				if distinct.TestType != tt.test {
					t.Errorf("TestType = %v, want %v", distinct.TestType, tt.test)
				}
				if distinct.Value == nil {
					t.Error("expected an operand after FROM")
				}
			})
		}
	})

	t.Run("with_and_filegroup", func(t *testing.T) {
		sql := "ALTER TABLE a ADD INDEX ix (a) WITH (FILLFACTOR = 1) ON 'a'"
		idx := addIndexAction(t, ParseAndCheck(t, sql), sql)
		if idx.Options == nil || idx.Options.Len() == 0 {
			t.Errorf("Options = %v, want FILLFACTOR", idx.Options)
		}
		if idx.OnFilegroup != "a" {
			t.Errorf("OnFilegroup = %q, want a", idx.OnFilegroup)
		}
	})

	t.Run("mixed_add_list", func(t *testing.T) {
		sql := "ALTER TABLE a ADD b INT, INDEX ix (b)"
		stmt := ParseAndCheck(t, sql).Items[0].(*ast.AlterTableStmt)
		if stmt.Actions == nil || stmt.Actions.Len() != 2 {
			t.Fatalf("Actions = %v, want a column and an index", stmt.Actions)
		}
		col, ok := stmt.Actions.Items[0].(*ast.AlterTableAction)
		if !ok || col.Type != ast.ATAddColumn || col.Column == nil || col.Column.Name != "b" {
			t.Errorf("first action = %v, want ADD b INT", stmt.Actions.Items[0])
		}
		idx, ok := stmt.Actions.Items[1].(*ast.AlterTableAction)
		if !ok || idx.Type != ast.ATAddIndex || idx.Index == nil || idx.Index.Name != "ix" {
			t.Errorf("second action = %v, want ADD INDEX ix", stmt.Actions.Items[1])
		}
	})
}

// TestParseAlterTableAddIndexRejects holds the inline index forms the pinned
// SQL Server 2022 refuses.
func TestParseAlterTableAddIndexRejects(t *testing.T) {
	for _, sql := range []string{
		"ALTER TABLE a ADD INDEX",
		"ALTER TABLE a ADD INDEX ix INCLUDE (a)",
		"ALTER TABLE a ADD INDEX ix COLUMNSTORE",
		"ALTER TABLE a ADD INDEX ix CLUSTERED COLUMNSTORE (a)",
		"ALTER TABLE a ADD INDEX ix CLUSTERED (a) WHERE a IS NULL",
		"ALTER TABLE a ADD INDEX ix WHERE a OR a",
		"ALTER TABLE a ADD INDEX ix WHERE a.IS NULL",
		"ALTER TABLE a ADD INDEX ix WHERE a IS DISTINCT FROM 1 OR b = 2",
		"ALTER TABLE a ADD INDEX ix ON",
		"ALTER TABLE a ADD INDEX ix ON ps(a",
		"ALTER TABLE a ADD INDEX ix ON ps()",
		"ALTER TABLE a ADD INDEX ix HASH",
	} {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Errorf("Parse(%q): want an error", sql)
			}
		})
	}
}

// TestParseCreateExternalTableNoAsSelect pins that an external table's column
// definitions and options end the statement: the AS SELECT tail ScriptDOM
// reads (ctasCreateExternalTableStatement, TSql160.g:546) is refused by the
// pinned SQL Server 2022, whose verdict compat.tsv holds.
func TestParseCreateExternalTableNoAsSelect(t *testing.T) {
	sql := "CREATE EXTERNAL TABLE a (a INT NOT NULL, b NVARCHAR(10)) WITH (LOCATION = '/a/', DATA_SOURCE = ds, FILE_FORMAT = ff)"
	result := ParseAndCheck(t, sql)
	stmt, ok := result.Items[0].(*ast.SecurityStmt)
	if !ok {
		t.Fatalf("Parse(%q): got %T, want *SecurityStmt", sql, result.Items[0])
	}
	if stmt.ObjectType != "EXTERNAL TABLE" {
		t.Errorf("ObjectType = %q, want EXTERNAL TABLE", stmt.ObjectType)
	}
	if stmt.Options == nil || stmt.Options.Len() != 5 {
		t.Errorf("Options = %v, want two columns and three options", stmt.Options)
	}

	for _, sql := range []string{
		"CREATE EXTERNAL TABLE a WITH (DATA_SOURCE = a) AS SELECT *",
		"CREATE EXTERNAL TABLE a (a INT) WITH (DATA_SOURCE = a) AS SELECT *",
	} {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Errorf("Parse(%q): want an error", sql)
			}
		})
	}
}

// addIndexAction parses sql and returns the inline index of its only ADD INDEX
// action.
func addIndexAction(t *testing.T, result *ast.List, sql string) *ast.InlineIndexDef {
	t.Helper()
	stmt, ok := result.Items[0].(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("Parse(%q): got %T, want *AlterTableStmt", sql, result.Items[0])
	}
	if stmt.Actions == nil || stmt.Actions.Len() != 1 {
		t.Fatalf("Parse(%q): Actions = %v, want one", sql, stmt.Actions)
	}
	action, ok := stmt.Actions.Items[0].(*ast.AlterTableAction)
	if !ok || action.Type != ast.ATAddIndex || action.Index == nil {
		t.Fatalf("Parse(%q): action = %v, want ADD INDEX", sql, stmt.Actions.Items[0])
	}
	return action.Index
}
