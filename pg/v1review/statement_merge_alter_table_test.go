package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementMergeAlterTable(t *testing.T) {
	redundant := func(table string, count, statement int) stmtFinding {
		return stmtFinding{Code: v1review.RedundantAlterTable, Table: table, Count: count, Statement: statement}
	}
	testStatementRule(t, StatementMergeAlterTable, []statementCase{
		{name: "create then alter", sql: "CREATE TABLE t (a int);\nALTER TABLE t ADD b int;", want: []stmtFinding{redundant("t", 2, 1)}},
		{name: "create alone", sql: "CREATE TABLE t (a int);"},
		{name: "alters without create", sql: "ALTER TABLE t ADD b int;\nALTER TABLE t ADD c int;"},
		{name: "public is the default schema", sql: "CREATE TABLE t (a int);\nALTER TABLE public.t ADD b int;\nALTER TABLE s.t ADD c int;", want: []stmtFinding{redundant("t", 2, 1)}},
		{name: "second create starts over", sql: "CREATE TABLE t (a int);\nALTER TABLE t ADD b int;\nCREATE TABLE t (a int);\nALTER TABLE t ADD b int;\nALTER TABLE t ADD c int;", want: []stmtFinding{redundant("t", 3, 4)}},
		{name: "second create alone", sql: "CREATE TABLE t (a int);\nALTER TABLE t ADD b int;\nCREATE TABLE t (a int);"},
		{name: "name from create", sql: "CREATE TABLE s.\"T\" (a int);\nALTER TABLE s.\"T\" ADD b int;", want: []stmtFinding{redundant("T", 2, 1)}},
		{name: "dotted names collide", sql: "CREATE TABLE \"a.b\".c (x int);\nALTER TABLE a.\"b.c\" ADD y int;", want: []stmtFinding{redundant("c", 2, 1)}},
		{name: "rename not counted", sql: "CREATE TABLE t (a int);\nALTER TABLE t RENAME TO u;\nALTER TABLE t RENAME a TO b;"},
		{name: "alter index counts", sql: "CREATE TABLE t (a int);\nALTER INDEX t SET (fillfactor = 70);", want: []stmtFinding{redundant("t", 2, 1)}},
		{
			name: "ordered by last statement",
			sql:  "CREATE TABLE u (a int);\nCREATE TABLE t (a int);\nALTER TABLE u ADD b int;\nALTER TABLE t ADD b int;\nALTER TABLE u ADD c int;",
			want: []stmtFinding{redundant("t", 2, 3), redundant("u", 3, 4)},
		},
	})
}
