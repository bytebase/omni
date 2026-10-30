package parser_test

import (
	"testing"

	"github.com/bytebase/omni/mysql/parser"
)

func TestContainsExecutableComment(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want bool
	}{
		{
			name: "ddl",
			sql:  "/*!50000 DROP TABLE t*/ SELECT 1",
			want: true,
		},
		{
			name: "alter",
			sql:  "/*!50000 ALTER TABLE t ADD COLUMN x INT*/ SELECT 1",
			want: true,
		},
		{
			name: "delete",
			sql:  "/*!50000 DELETE FROM t*/ SELECT 1",
			want: true,
		},
		{
			name: "insert",
			sql:  "/*!50000 INSERT INTO t VALUES (1)*/ SELECT 1",
			want: true,
		},
		{
			name: "unknown_version",
			sql:  "/*!99999 SELECT 1*/",
			want: true,
		},
		{
			name: "expression",
			sql:  "SELECT /*!80100 42*/ FROM dual",
			want: true,
		},
		{
			name: "no_version",
			sql:  "SELECT /*! 1 + 1 */ FROM dual",
			want: true,
		},
		{
			name: "multiline",
			sql:  "/*!50000\nDROP TABLE t;\n*/ SELECT 1",
			want: true,
		},
		{
			name: "multiple_statements_in_comment",
			sql:  "/*!50000 DROP TABLE t; DELETE FROM t2*/ SELECT 1",
			want: true,
		},
		{
			name: "cte_dml",
			sql:  "/*!50000 WITH cte AS (SELECT 1) DELETE FROM t*/ SELECT 1",
			want: true,
		},
		{
			name: "nested_comment",
			sql:  "SELECT /*!80100 1 /* nested */ + 2 */",
			want: true,
		},
		{
			name: "select_list",
			sql:  "SELECT a, /*!80100 b, */ c FROM t",
			want: true,
		},
		{
			name: "invalid_version_format",
			sql:  "/*!version SELECT 1*/",
			want: true,
		},
		{
			name: "unterminated",
			sql:  "/*!50000 SELECT 1",
			want: true,
		},
		{
			name: "string_literal",
			sql:  "SELECT '/*!50000 harmless*/' AS str",
			want: false,
		},
		{
			name: "quoted_identifier",
			sql:  "SELECT `/*!50000` AS quoted",
			want: false,
		},
		{
			name: "ordinary_block_comment",
			sql:  "SELECT 1 /* benign comment */",
			want: false,
		},
		{
			name: "optimizer_hint",
			sql:  "SELECT /*+ optimizer_hint */ 1",
			want: false,
		},
		{
			name: "line_comment",
			sql:  "-- /*!50000 line comment\nSELECT 1",
			want: false,
		},
		{
			name: "line_comment_form_feed",
			sql:  "--\f/*!50000 DROP TABLE t*/",
			want: false,
		},
		{
			name: "line_comment_vertical_tab",
			sql:  "--\v/*!50000 DROP TABLE t*/",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parser.ContainsExecutableComment(tt.sql, parser.ExecutableCommentOptions{}); got != tt.want {
				t.Fatalf("ContainsExecutableComment(%q) = %v, want %v", tt.sql, got, tt.want)
			}
		})
	}
}

func TestContainsExecutableCommentNoBackslashEscapes(t *testing.T) {
	sql := "SELECT 'x\\' /*!50000 + 1 */"

	if got := parser.ContainsExecutableComment(sql, parser.ExecutableCommentOptions{}); got {
		t.Fatalf("ContainsExecutableComment(%q, default mode) = true, want false", sql)
	}
	if got := parser.ContainsExecutableComment(sql, parser.ExecutableCommentOptions{NoBackslashEscapes: true}); !got {
		t.Fatalf("ContainsExecutableComment(%q, NO_BACKSLASH_ESCAPES) = false, want true", sql)
	}
}
