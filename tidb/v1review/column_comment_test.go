package v1review

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnComment(t *testing.T) {
	missing := func(table, name, text string) placed {
		return placed{Code: v1review.ColumnCommentMissing, Name: name, Table: table, Text: text}
	}
	tooLong := func(table, name, text string) placed {
		return placed{Code: v1review.ColumnCommentTooLong, Name: name, Table: table, Text: text}
	}
	const create = "CREATE TABLE t (a INT, b INT COMMENT '', c INT COMMENT 'héllo')"
	const alter = "ALTER TABLE t ADD d INT, CHANGE e f INT COMMENT 'abcdef', MODIFY g INT, RENAME COLUMN h TO i"
	for _, tt := range []struct {
		name      string
		sql       string
		required  bool
		maxLength int
		want      []placed
	}{
		{
			name: "required, length in bytes", sql: create + ";", required: true, maxLength: 5,
			want: []placed{missing("t", "a", "a INT"), tooLong("t", "c", "c INT COMMENT 'héllo'")},
		},
		{name: "length within", sql: create + ";", maxLength: 6},
		{
			name: "missing before too long", sql: "CREATE TABLE t (a INT);", required: true, maxLength: -1,
			want: []placed{missing("t", "a", "a INT"), tooLong("t", "a", "a INT")},
		},
		{
			name: "negative length reports every column", sql: create + ";", maxLength: -1,
			want: []placed{tooLong("t", "a", "a INT"), tooLong("t", "b", "b INT COMMENT ''"), tooLong("t", "c", "c INT COMMENT 'héllo'")},
		},
		{
			name: "alter anchors the statement", sql: alter + ";", required: true, maxLength: 5,
			want: []placed{missing("t", "d", alter), tooLong("t", "f", alter), missing("t", "g", alter)},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return ColumnComment(stmts, tt.required, tt.maxLength)
			}, tt.want)
		})
	}
}
