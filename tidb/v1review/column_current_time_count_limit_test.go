package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestColumnCurrentTimeCountLimit(t *testing.T) {
	defaults := func(table, count string, statement int, text string) placed {
		return placed{Code: v1review.ColumnDefaultCurrentTimeCount, Table: table, Value: count, Statement: statement, Text: text}
	}
	updates := func(table, count string, statement int, text string) placed {
		return placed{Code: v1review.ColumnOnUpdateCurrentTimeCount, Table: table, Value: count, Statement: statement, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []placed
	}{
		{name: "within the limits", sql: "CREATE TABLE t (a DATETIME DEFAULT NOW() ON UPDATE NOW(), b TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"},
		{
			name: "over both limits",
			sql:  "CREATE TABLE t (a DATETIME DEFAULT NOW() ON UPDATE localtime(), b TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, c TIMESTAMP DEFAULT LOCALTIMESTAMP);",
			want: []placed{
				defaults("t", "3", 0, "CREATE TABLE t (a DATETIME DEFAULT NOW() ON UPDATE localtime(), b TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, c TIMESTAMP DEFAULT LOCALTIMESTAMP)"),
				updates("t", "2", 0, "CREATE TABLE t (a DATETIME DEFAULT NOW() ON UPDATE localtime(), b TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, c TIMESTAMP DEFAULT LOCALTIMESTAMP)"),
			},
		},
		{name: "other types and functions are not counted", sql: "CREATE TABLE t (a DATE DEFAULT NOW(), b TIME DEFAULT NOW(), c DATETIME DEFAULT UTC_TIMESTAMP(), d DATETIME DEFAULT '2020-01-01');"},
		{
			name: "counts accumulate across statements and anchor the last",
			sql:  "CREATE TABLE t (a DATETIME ON UPDATE NOW());\nDROP TABLE t;\nCREATE TABLE t (b INT);\nALTER TABLE t ADD c DATETIME ON UPDATE NOW(), MODIFY d INT;\nALTER TABLE t RENAME COLUMN c TO e;",
			want: []placed{updates("t", "2", 3, "ALTER TABLE t ADD c DATETIME ON UPDATE NOW(), MODIFY d INT")},
		},
		{
			name: "tables ordered by anchor",
			sql:  "CREATE TABLE b (x DATETIME ON UPDATE NOW(), y DATETIME ON UPDATE NOW());\nCREATE TABLE a (x DATETIME ON UPDATE NOW(), y DATETIME ON UPDATE NOW());\nALTER TABLE b CHANGE z z DATETIME DEFAULT NOW();",
			want: []placed{
				updates("a", "2", 1, "CREATE TABLE a (x DATETIME ON UPDATE NOW(), y DATETIME ON UPDATE NOW())"),
				updates("b", "2", 2, "ALTER TABLE b CHANGE z z DATETIME DEFAULT NOW()"),
			},
		},
		{
			// V1 sorted on the line alone, so this tie came out in either
			// order; the earlier statement comes first here.
			name: "a tie on one line keeps statement order",
			sql:  "CREATE TABLE b (x DATETIME ON UPDATE NOW(), y DATETIME ON UPDATE NOW()); CREATE TABLE a (x DATETIME ON UPDATE NOW(), y DATETIME ON UPDATE NOW());",
			want: []placed{
				updates("b", "2", 0, "CREATE TABLE b (x DATETIME ON UPDATE NOW(), y DATETIME ON UPDATE NOW())"),
				updates("a", "2", 1, "CREATE TABLE a (x DATETIME ON UPDATE NOW(), y DATETIME ON UPDATE NOW())"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, ColumnCurrentTimeCountLimit, tt.want)
		})
	}
}
