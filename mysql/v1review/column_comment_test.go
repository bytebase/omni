package v1review

import (
	"testing"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnComment(t *testing.T) {
	rule := func(required bool, maxLength int) func([]ast.Node) []v1review.Finding {
		return func(stmts []ast.Node) []v1review.Finding { return ColumnComment(stmts, required, maxLength) }
	}
	runColumnCases(t, rule(true, 5), []columnCase{
		{name: "short comments", sql: "CREATE TABLE t(a INT COMMENT 'ab', b INT COMMENT 'abcde');"},
		{
			name: "missing, empty and long, length in bytes",
			sql:  "CREATE TABLE t(a INT, b INT COMMENT '', c INT COMMENT 'âââ');",
			want: []string{
				"column-comment-missing t.a  @0: a INT",
				"column-comment-missing t.b  @0: b INT COMMENT ''",
				"column-comment-too-long t.c  @0: c INT COMMENT 'âââ'",
			},
		},
		{
			name: "add, modify and change",
			sql:  "ALTER TABLE t ADD COLUMN (a INT), MODIFY b INT, CHANGE c d INT COMMENT 'x', DROP e;",
			want: []string{
				"column-comment-missing t.a  @0: a INT",
				"column-comment-missing t.b  @0: b INT",
			},
		},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t(a INT); END;"},
	})
	runColumnCases(t, rule(false, -1), []columnCase{
		{name: "nothing required, no limit", sql: "CREATE TABLE t(a INT, b INT COMMENT 'abcdefgh');"},
	})
	runColumnCases(t, rule(true, 0), []columnCase{
		{
			name: "every judged column yields a finding",
			sql:  "CREATE TABLE t(a INT, b INT COMMENT 'x');",
			want: []string{
				"column-comment-missing t.a  @0: a INT",
				"column-comment-too-long t.b  @0: b INT COMMENT 'x'",
			},
		},
	})
}
