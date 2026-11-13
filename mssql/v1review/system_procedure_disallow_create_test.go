package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mssql/parser"
	"github.com/bytebase/omni/v1review"
)

func TestSystemProcedureDisallowCreate(t *testing.T) {
	created := func(text string) anchoredFinding {
		return anchoredFinding{Code: v1review.ProcedureCreated, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []anchoredFinding
	}{
		{
			name: "create procedure",
			sql:  "CREATE PROCEDURE dbo.p AS SELECT 1;",
			want: []anchoredFinding{created("CREATE PROCEDURE dbo.p AS SELECT 1")},
		},
		{
			name: "create proc abbreviation",
			sql:  "CREATE PROC p AS SELECT 1;",
			want: []anchoredFinding{created("CREATE PROC p AS SELECT 1")},
		},
		{
			name: "create or alter procedure",
			sql:  "CREATE OR ALTER PROCEDURE p AS SELECT 1;",
			want: []anchoredFinding{created("CREATE OR ALTER PROCEDURE p AS SELECT 1")},
		},
		{
			name: "alter procedure",
			sql:  "ALTER PROCEDURE p AS SELECT 2;",
			want: []anchoredFinding{created("ALTER PROCEDURE p AS SELECT 2")},
		},
		{
			name: "statement on a later line",
			sql:  "SELECT 1;\n\n-- note\nCREATE PROCEDURE\n  p AS SELECT 1;",
			want: []anchoredFinding{created("CREATE PROCEDURE\n  p AS SELECT 1")},
		},
		{name: "function not judged", sql: "CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN 1 END;"},
		{name: "exec not judged", sql: "EXEC p;"},
		{name: "drop procedure not judged", sql: "DROP PROCEDURE p;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			got := anchored(t, tt.sql, SystemProcedureDisallowCreate(list.Items))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
