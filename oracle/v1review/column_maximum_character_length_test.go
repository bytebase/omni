package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestColumnMaximumCharacterLength(t *testing.T) {
	long := func(name, text string) columnFinding {
		return columnFinding{Code: v1review.ColumnCharTooLong, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []columnFinding
	}{
		{name: "char over", sql: "CREATE TABLE t (a CHAR(21));", want: []columnFinding{long("A", "CHAR(21)")}},
		{name: "at the limit", sql: "CREATE TABLE t (a CHAR(20));"},
		{name: "character over", sql: "CREATE TABLE t (a character(30));", want: []columnFinding{long("A", "character(30)")}},
		{name: "unit ignored", sql: "CREATE TABLE t (a CHAR(30 BYTE), b CHAR(30 CHAR));", want: []columnFinding{long("A", "CHAR(30 BYTE)"), long("B", "CHAR(30 CHAR)")}},
		{name: "nchar not judged", sql: "CREATE TABLE t (a NCHAR(30));"},
		{name: "varchar2 not judged", sql: "CREATE TABLE t (a VARCHAR2(30));"},
		{name: "no length not judged", sql: "CREATE TABLE t (a CHAR);"},
		{name: "alter add and modify", sql: "ALTER TABLE t ADD a CHAR(30);\nALTER TABLE t MODIFY b CHAR(40);", want: []columnFinding{long("A", "CHAR(30)"), long("B", "CHAR(40)")}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := columnFindings(tt.sql, ColumnMaximumCharacterLength(parseStmts(t, tt.sql), 20))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
