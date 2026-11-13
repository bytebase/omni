package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mssql/parser"
	"github.com/bytebase/omni/v1review"
)

// anchoredFinding is a finding as the system rule tests compare it: its
// code, name and statement, and the text its range anchors.
type anchoredFinding struct {
	Code      v1review.Code
	Name      string
	Statement int
	Text      string
}

func anchored(t *testing.T, sql string, findings []v1review.Finding) []anchoredFinding {
	t.Helper()
	var got []anchoredFinding
	for _, f := range findings {
		got = append(got, anchoredFinding{f.Code, f.Name, f.Statement, sql[f.Start:f.End]})
	}
	return got
}

func TestSystemFunctionDisallowCreate(t *testing.T) {
	created := func(text string) anchoredFinding {
		return anchoredFinding{Code: v1review.FunctionCreated, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []anchoredFinding
	}{
		{
			name: "create function",
			sql:  "CREATE FUNCTION dbo.f() RETURNS INT AS BEGIN RETURN 1 END;",
			want: []anchoredFinding{created("CREATE FUNCTION dbo.f() RETURNS INT AS BEGIN RETURN 1 END")},
		},
		{
			name: "create or alter function",
			sql:  "CREATE OR ALTER FUNCTION f() RETURNS INT AS BEGIN RETURN 1 END;",
			want: []anchoredFinding{created("CREATE OR ALTER FUNCTION f() RETURNS INT AS BEGIN RETURN 1 END")},
		},
		{
			name: "alter function",
			sql:  "ALTER FUNCTION f() RETURNS INT AS BEGIN RETURN 2 END;",
			want: []anchoredFinding{created("ALTER FUNCTION f() RETURNS INT AS BEGIN RETURN 2 END")},
		},
		{
			name: "every statement judged in order",
			sql:  "CREATE TABLE t(a INT);\nCREATE FUNCTION f1() RETURNS INT AS BEGIN RETURN 1 END;\nCREATE FUNCTION f2() RETURNS INT AS BEGIN RETURN 2 END;",
			want: []anchoredFinding{
				created("CREATE FUNCTION f1() RETURNS INT AS BEGIN RETURN 1 END"),
				created("CREATE FUNCTION f2() RETURNS INT AS BEGIN RETURN 2 END"),
			},
		},
		{name: "procedure not judged", sql: "CREATE PROCEDURE p AS SELECT 1;"},
		{name: "drop function not judged", sql: "DROP FUNCTION f;"},
		{name: "begin block not judged", sql: "BEGIN EXEC('x'); END"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			got := anchored(t, tt.sql, SystemFunctionDisallowCreate(list.Items))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
