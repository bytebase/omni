package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestColumnMaximumVarcharLength(t *testing.T) {
	long := func(name, text string) columnFinding {
		return columnFinding{Code: v1review.ColumnVarcharTooLong, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []columnFinding
	}{
		{name: "varchar2 over", sql: "CREATE TABLE t (a VARCHAR2(101));", want: []columnFinding{long("A", "VARCHAR2(101)")}},
		{name: "at the limit", sql: "CREATE TABLE t (a VARCHAR2(100));"},
		{name: "varchar over", sql: "CREATE TABLE t (a varchar(200 char));", want: []columnFinding{long("A", "varchar(200 char)")}},
		{name: "nvarchar2 not judged", sql: "CREATE TABLE t (a NVARCHAR2(200));"},
		{name: "char not judged", sql: "CREATE TABLE t (a CHAR(200));"},
		{
			name: "every statement and column",
			sql:  "CREATE TABLE t (a VARCHAR2(200), b VARCHAR2(10));\nALTER TABLE t MODIFY (b VARCHAR2(300));",
			want: []columnFinding{long("A", "VARCHAR2(200)"), long("B", "VARCHAR2(300)")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := columnFindings(tt.sql, ColumnMaximumVarcharLength(parseStmts(t, tt.sql), 100))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
