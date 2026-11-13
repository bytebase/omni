package v1review

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestColumnComment(t *testing.T) {
	finding := func(code v1review.Code, table, name string, statement int, text string) columnFinding {
		return columnFinding{Code: code, Name: name, Table: table, Statement: statement, Text: text}
	}
	missing := v1review.ColumnCommentMissing
	long := v1review.ColumnCommentTooLong
	for _, tt := range []struct {
		name      string
		sql       string
		required  bool
		maxLength int
		want      []columnFinding
	}{
		{
			name:      "missing in definition order",
			sql:       "CREATE TABLE t (b int, a int);\nCOMMENT ON COLUMN t.b IS 'ok';",
			required:  true,
			maxLength: math.MaxInt,
			want:      []columnFinding{finding(missing, "DB.T", "A", 0, "a int")},
		},
		{name: "not required", sql: "CREATE TABLE t (a int);", maxLength: math.MaxInt},
		{
			name:      "comment before the definition counts",
			sql:       "COMMENT ON COLUMN t.a IS 'ok';\nCREATE TABLE t (a int);",
			required:  true,
			maxLength: math.MaxInt,
		},
		{
			name:      "last comment wins",
			sql:       "CREATE TABLE t (a int);\nCOMMENT ON COLUMN t.a IS 'ok';\nCOMMENT ON COLUMN t.a IS '';",
			required:  true,
			maxLength: math.MaxInt,
			want:      []columnFinding{finding(missing, "DB.T", "A", 0, "a int")},
		},
		{
			name:      "too long in bytes",
			sql:       "ALTER TABLE s.t ADD (a int, b int);\nCOMMENT ON COLUMN s.t.a IS '注释';\nCOMMENT ON COLUMN s.t.b IS 'abcde';",
			maxLength: 5,
			want:      []columnFinding{finding(long, "S.T", "A", 0, "a int")},
		},
		{
			name:      "negative limit reports every commented column",
			sql:       "CREATE TABLE t (a int, b int);\nCOMMENT ON COLUMN t.a IS 'x';",
			required:  true,
			maxLength: -1,
			want:      []columnFinding{finding(long, "DB.T", "A", 0, "a int"), finding(missing, "DB.T", "B", 0, "b int")},
		},
		{
			name:      "defined twice judged twice at the last definition",
			sql:       "CREATE TABLE t (a int);\nALTER TABLE t ADD a int;",
			required:  true,
			maxLength: math.MaxInt,
			want:      []columnFinding{finding(missing, "DB.T", "A", 1, "a int"), finding(missing, "DB.T", "A", 1, "a int")},
		},
		{
			name:      "current database qualifies keys",
			sql:       "CREATE TABLE t (a int);\nCOMMENT ON COLUMN db.t.a IS 'ok';\nCREATE TABLE s.t (a int);",
			required:  true,
			maxLength: math.MaxInt,
			want:      []columnFinding{finding(missing, "S.T", "A", 2, "a int")},
		},
		{
			name:      "modify and table comment not judged",
			sql:       "ALTER TABLE t MODIFY a int;\nCOMMENT ON TABLE t IS 'x';",
			required:  true,
			maxLength: math.MaxInt,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := columnFindings(tt.sql, ColumnComment(parseStmts(t, tt.sql), "DB", tt.required, tt.maxLength))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
