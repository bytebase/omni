package analysis

import (
	"reflect"
	"testing"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
)

func TestClassifyWalkThroughStatement(t *testing.T) {
	if got := ClassifyWalkThroughStatement(nil); !reflect.DeepEqual(got, WalkThroughStatement{}) {
		t.Fatalf("ClassifyWalkThroughStatement(nil) = %+v, want zero value", got)
	}
	if got := ClassifyWalkThroughStatement((*ast.CreateTableStmt)(nil)); !reflect.DeepEqual(got, WalkThroughStatement{}) {
		t.Fatalf("ClassifyWalkThroughStatement(typed nil) = %+v, want zero value", got)
	}

	tests := []struct {
		name string
		sql  string
		want WalkThroughStatement
	}{
		{
			name: "select is other",
			sql:  "SELECT 1",
			want: WalkThroughStatement{Kind: WalkThroughStatementOther},
		},
		{
			name: "create database",
			sql:  "CREATE DATABASE db1",
			want: WalkThroughStatement{Kind: WalkThroughStatementCreateDatabase, Database: "db1"},
		},
		{
			name: "alter database",
			sql:  "ALTER DATABASE db1 CHARACTER SET utf8mb4",
			want: WalkThroughStatement{Kind: WalkThroughStatementAlterDatabase, Database: "db1"},
		},
		{
			name: "drop database",
			sql:  "DROP DATABASE db1",
			want: WalkThroughStatement{Kind: WalkThroughStatementDropDatabase, Database: "db1"},
		},
		{
			name: "create table",
			sql:  "CREATE TABLE db1.t1 (id INT)",
			want: WalkThroughStatement{
				Kind:      WalkThroughStatementCreateTable,
				TableRefs: []TableRef{{Schema: "db1", Name: "t1"}},
			},
		},
		{
			name: "create table if not exists",
			sql:  "CREATE TABLE IF NOT EXISTS t1 (id INT)",
			want: WalkThroughStatement{
				Kind:        WalkThroughStatementCreateTable,
				IfNotExists: true,
				TableRefs:   []TableRef{{Name: "t1"}},
			},
		},
		{
			name: "create table as select",
			sql:  "CREATE TABLE IF NOT EXISTS db1.t2 AS SELECT * FROM t1",
			want: WalkThroughStatement{
				Kind:        WalkThroughStatementCreateTable,
				IfNotExists: true,
				IsCTAS:      true,
				TableRefs:   []TableRef{{Schema: "db1", Name: "t2"}},
			},
		},
		{
			name: "alter table add column",
			sql:  "ALTER TABLE db1.t1 ADD COLUMN c INT",
			want: WalkThroughStatement{
				Kind:      WalkThroughStatementTableRefCheck,
				TableRefs: []TableRef{{Schema: "db1", Name: "t1"}},
			},
		},
		{
			name: "alter table rename",
			sql:  "ALTER TABLE db1.t1 RENAME TO db2.t2",
			want: WalkThroughStatement{
				Kind:      WalkThroughStatementTableRefCheck,
				TableRefs: []TableRef{{Schema: "db1", Name: "t1"}},
				RenamePairs: []RenamePair{
					{Old: TableRef{Schema: "db1", Name: "t1"}, New: TableRef{Schema: "db2", Name: "t2"}},
				},
			},
		},
		{
			name: "drop table multiple",
			sql:  "DROP TABLE t1, db2.t2",
			want: WalkThroughStatement{
				Kind:      WalkThroughStatementDropTable,
				TableRefs: []TableRef{{Name: "t1"}, {Schema: "db2", Name: "t2"}},
			},
		},
		{
			name: "create index",
			sql:  "CREATE INDEX idx ON db1.t1 (id)",
			want: WalkThroughStatement{
				Kind:      WalkThroughStatementTableRefCheck,
				TableRefs: []TableRef{{Schema: "db1", Name: "t1"}},
			},
		},
		{
			name: "drop index",
			sql:  "DROP INDEX idx ON db1.t1",
			want: WalkThroughStatement{
				Kind:      WalkThroughStatementTableRefCheck,
				TableRefs: []TableRef{{Schema: "db1", Name: "t1"}},
			},
		},
		{
			name: "truncate table",
			sql:  "TRUNCATE TABLE db1.t1",
			want: WalkThroughStatement{
				Kind:      WalkThroughStatementTableRefCheck,
				TableRefs: []TableRef{{Schema: "db1", Name: "t1"}},
			},
		},
		{
			name: "rename table multi-pair",
			sql:  "RENAME TABLE t1 TO db2.t2, db1.a TO db3.b",
			want: WalkThroughStatement{
				Kind:      WalkThroughStatementRenameTable,
				TableRefs: []TableRef{{Name: "t1"}, {Schema: "db1", Name: "a"}},
				RenamePairs: []RenamePair{
					{Old: TableRef{Name: "t1"}, New: TableRef{Schema: "db2", Name: "t2"}},
					{Old: TableRef{Schema: "db1", Name: "a"}, New: TableRef{Schema: "db3", Name: "b"}},
				},
			},
		},
		{
			name: "create view",
			sql:  "CREATE VIEW db1.v1 AS SELECT 1",
			want: WalkThroughStatement{
				Kind:      WalkThroughStatementTableRefCheck,
				TableRefs: []TableRef{{Schema: "db1", Name: "v1"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("parser.Parse(%q) error = %v", tt.sql, err)
			}
			if list == nil || len(list.Items) != 1 {
				t.Fatalf("parser.Parse(%q) returned %d statements, want 1", tt.sql, len(list.Items))
			}
			got := ClassifyWalkThroughStatement(list.Items[0])
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClassifyWalkThroughStatement(%q) =\n%+v\nwant:\n%+v", tt.sql, got, tt.want)
			}
		})
	}
}

func TestIsCreateTableAs(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want bool
	}{
		{name: "ctas with AS", sql: "CREATE TABLE t2 AS SELECT * FROM t1;", want: true},
		{name: "ctas without AS", sql: "CREATE TABLE t2 (b INT) SELECT a FROM t1;", want: true},
		{name: "ctas temporary if not exists", sql: "CREATE TEMPORARY TABLE IF NOT EXISTS db.t2 AS SELECT 1;", want: true},
		{name: "ctas values row", sql: "CREATE TABLE tvx SELECT x, z FROM (VALUES ROW(1,3,5), ROW(2,4,6)) AS v(x,y,z);", want: true},
		{name: "ctas as table", sql: "CREATE TABLE t2 AS TABLE t1;", want: true},
		{name: "ctas bare table", sql: "CREATE TABLE t2 TABLE t1;", want: true},
		{name: "ctas parenthesized table", sql: "CREATE TABLE t2 (TABLE t1);", want: true},
		{name: "ctas as values row", sql: "CREATE TABLE t2 AS VALUES ROW(1, 2);", want: true},
		{name: "ctas parenthesized values row", sql: "CREATE TABLE t2 (a INT) ((VALUES ROW(1)));", want: true},
		{name: "ctas parenthesized cte table", sql: "CREATE TABLE t2 (WITH cte AS (VALUES ROW(1)) TABLE cte);", want: true},
		{name: "ctas with comments", sql: "-- lead\nCREATE /* mid */ TEMPORARY # hash\nTABLE t2 AS SELECT 1", want: true},
		{name: "ctas lowercase", sql: "create table t2 select 1", want: true},
		{name: "ctas executable comment", sql: "CREATE TABLE t2 /*!50000 SELECT 1 */", want: true},
		{name: "regular create table", sql: "CREATE TABLE t1 (id INT);", want: false},
		{name: "partitioned table with VALUES LESS THAN", sql: "CREATE TABLE t1 (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10));", want: false},
		{name: "partitioned table with VALUES IN", sql: "CREATE TABLE t1 (id INT) PARTITION BY LIST (id) (PARTITION p0 VALUES IN (1, 2));", want: false},
		{name: "create table with tablespace", sql: "CREATE TABLE t1 (id INT) TABLESPACE ts1;", want: false},
		{name: "create table fulltext with parser", sql: "CREATE TABLE t1 (body TEXT, FULLTEXT INDEX ft (body) WITH PARSER ngram);", want: false},
		{name: "select in single-quoted default", sql: "CREATE TABLE t1 (a VARCHAR(20) DEFAULT 'SELECT');", want: false},
		{name: "select in escaped single quote", sql: `CREATE TABLE t1 (a VARCHAR(20) DEFAULT 'it\'s SELECT');`, want: false},
		{name: "select in double-quoted comment", sql: `CREATE TABLE t1 (a INT COMMENT "SELECT 1");`, want: false},
		{name: "select in backtick identifier", sql: "CREATE TABLE t1 (`select` INT);", want: false},
		{name: "select in comments", sql: "CREATE TABLE t1 (a INT) /* SELECT */ -- SELECT\n# SELECT", want: false},
		{name: "dot-qualified select table name", sql: "CREATE TABLE db.select (a INT);", want: false},
		{name: "select-qualified schema name", sql: "CREATE TABLE select.t1 (a INT);", want: false},
		{name: "select in next statement after semicolon", sql: "CREATE TABLE t1 (a INT); SELECT 1;", want: false},
		{name: "create view as select", sql: "CREATE VIEW v1 AS SELECT 1;", want: false},
		{name: "plain select", sql: "SELECT 1;", want: false},
		{name: "empty", sql: "", want: false},
		{name: "comment only", sql: "-- CREATE TABLE t AS SELECT 1", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCreateTableAs(tt.sql); got != tt.want {
				t.Errorf("IsCreateTableAs(%q) = %v, want %v", tt.sql, got, tt.want)
			}
		})
	}
}
