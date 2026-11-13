package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestColumnTypeDisallowList(t *testing.T) {
	bad := func(name, typ, text string) columnFinding {
		return columnFinding{Code: v1review.ColumnTypeDisallowed, Name: name, Value: typ, Text: text}
	}
	for _, tt := range []struct {
		name       string
		sql        string
		disallowed []string
		want       []columnFinding
	}{
		{
			name:       "case-insensitive match",
			sql:        "CREATE TABLE t (a json, b number);",
			disallowed: []string{"JSON"},
			want:       []columnFinding{bad("A", "JSON", "json")},
		},
		{
			name:       "modifiers not part of the type",
			sql:        "CREATE TABLE t (a VARCHAR2(10), b NUMBER(10,2));",
			disallowed: []string{"varchar2", "NUMBER(10,2)"},
			want:       []columnFinding{bad("A", "VARCHAR2", "VARCHAR2(10)")},
		},
		{
			name:       "qualified type joined by dots",
			sql:        "CREATE TABLE t (a sys.xmltype);",
			disallowed: []string{"SYS.XMLTYPE"},
			want:       []columnFinding{bad("A", "SYS.XMLTYPE", "sys.xmltype")},
		},
		{
			name:       "reported once for several entries",
			sql:        "CREATE TABLE t (a BLOB);",
			disallowed: []string{"blob", "BLOB"},
			want:       []columnFinding{bad("A", "BLOB", "BLOB")},
		},
		{
			name:       "alter add and modify",
			sql:        "ALTER TABLE t ADD a BLOB;\nALTER TABLE t MODIFY (b BLOB);",
			disallowed: []string{"BLOB"},
			want:       []columnFinding{bad("A", "BLOB", "BLOB"), bad("B", "BLOB", "BLOB")},
		},
		{name: "empty list", sql: "CREATE TABLE t (a BLOB);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := columnFindings(tt.sql, ColumnTypeDisallowList(parseStmts(t, tt.sql), tt.disallowed))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
