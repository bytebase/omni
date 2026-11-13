package v1review

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnMaximumCharacterLength(t *testing.T) {
	long := func(name, length, text string) placed {
		return placed{Code: v1review.ColumnCharTooLong, Name: name, Value: length, Text: text}
	}
	for _, tt := range []struct {
		name    string
		sql     string
		maximum int
		want    []placed
	}{
		{name: "within", sql: "CREATE TABLE t (a CHAR(20), b VARCHAR(255), c BINARY(20), d TEXT);", maximum: 20},
		{
			name: "create reports the first", sql: "CREATE TABLE t (a CHAR(21), b BINARY(30));", maximum: 20,
			want: []placed{long("a", "21", "a CHAR(21)")},
		},
		{
			name: "binary counts, a bare char is 1", sql: "CREATE TABLE t (a CHAR, b BINARY(30));", maximum: 1,
			want: []placed{long("b", "30", "b BINARY(30)")},
		},
		{
			name: "add reports its last", sql: "ALTER TABLE t ADD (a CHAR(30), b CHAR(40)), MODIFY c CHAR(50);", maximum: 20,
			want: []placed{long("b", "40", "ALTER TABLE t ADD (a CHAR(30), b CHAR(40)), MODIFY c CHAR(50)")},
		},
		{
			name: "first action decides", sql: "ALTER TABLE t ADD a INT, CHANGE b c BINARY(30), ADD d CHAR(40);", maximum: 20,
			want: []placed{long("c", "30", "ALTER TABLE t ADD a INT, CHANGE b c BINARY(30), ADD d CHAR(40)")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return ColumnMaximumCharacterLength(stmts, tt.maximum)
			}, tt.want)
		})
	}
}
