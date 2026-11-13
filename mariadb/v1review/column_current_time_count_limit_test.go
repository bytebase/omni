package v1review

import "testing"

func TestColumnCurrentTimeCountLimit(t *testing.T) {
	runColumnCases(t, ColumnCurrentTimeCountLimit, []columnCase{
		{name: "within limits", sql: "CREATE TABLE t(a TIMESTAMP DEFAULT NOW() ON UPDATE NOW(), b DATETIME DEFAULT CURRENT_TIMESTAMP);"},
		{
			name: "both counts share the last column counted for either",
			sql:  "CREATE TABLE t(\n  a TIMESTAMP DEFAULT now() ON UPDATE now(),\n  b datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),\n  c TIMESTAMP DEFAULT localtime,\n  d INT);",
			want: []string{
				"column-default-current-time-count t. 3 @0: c TIMESTAMP DEFAULT localtime",
				"column-on-update-current-time-count t. 2 @0: c TIMESTAMP DEFAULT localtime",
			},
		},
		{
			name: "counts accumulate across statements, create does not reset",
			sql:  "CREATE TABLE t(a TIMESTAMP ON UPDATE NOW());\nCREATE TABLE t(b TIMESTAMP ON UPDATE NOW());",
			want: []string{"column-on-update-current-time-count t. 2 @1: b TIMESTAMP ON UPDATE NOW()"},
		},
		{
			name: "tables in anchor order, qualifier dropped",
			sql:  "ALTER TABLE a ADD x TIMESTAMP ON UPDATE NOW();\nALTER TABLE db.b ADD x TIMESTAMP ON UPDATE NOW(), ADD y TIMESTAMP ON UPDATE NOW();\nALTER TABLE a MODIFY x TIMESTAMP ON UPDATE NOW();",
			want: []string{
				"column-on-update-current-time-count b. 2 @1: y TIMESTAMP ON UPDATE NOW()",
				"column-on-update-current-time-count a. 2 @2: x TIMESTAMP ON UPDATE NOW()",
			},
		},
		{
			// V1's order here depended on map iteration; the earlier statement goes first.
			name: "tables anchored on one line come in statement order",
			sql:  "ALTER TABLE b ADD x TIMESTAMP ON UPDATE NOW(), ADD y TIMESTAMP ON UPDATE NOW(); ALTER TABLE a ADD x TIMESTAMP ON UPDATE NOW(), ADD y TIMESTAMP ON UPDATE NOW();",
			want: []string{
				"column-on-update-current-time-count b. 2 @0: y TIMESTAMP ON UPDATE NOW()",
				"column-on-update-current-time-count a. 2 @1: y TIMESTAMP ON UPDATE NOW()",
			},
		},
		{name: "other types not counted", sql: "CREATE TABLE t(a DATE DEFAULT NOW(), b DATE DEFAULT NOW(), c DATE DEFAULT NOW());"},
		{name: "drop table does not forget", sql: "CREATE TABLE t(a TIMESTAMP ON UPDATE NOW()); DROP TABLE t;"},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t(a TIMESTAMP ON UPDATE NOW(), b TIMESTAMP ON UPDATE NOW()); END;"},
	})
}
