package analysis

import (
	"reflect"
	"testing"

	"github.com/bytebase/omni/tidb/parser"
)

func TestClassifyWalkThroughStatement(t *testing.T) {
	if got := ClassifyWalkThroughStatement(nil); !reflect.DeepEqual(got, WalkThroughStatement{}) {
		t.Fatalf("ClassifyWalkThroughStatement(nil) = %+v, want zero value", got)
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
		{
			name: "create sequence",
			sql:  "CREATE SEQUENCE db1.seq1",
			want: WalkThroughStatement{
				Kind:      WalkThroughStatementTableRefCheck,
				TableRefs: []TableRef{{Schema: "db1", Name: "seq1"}},
			},
		},
		{
			name: "alter sequence",
			sql:  "ALTER SEQUENCE db1.seq1 INCREMENT BY 2",
			want: WalkThroughStatement{
				Kind:      WalkThroughStatementTableRefCheck,
				TableRefs: []TableRef{{Schema: "db1", Name: "seq1"}},
			},
		},
		{
			name: "drop sequence multiple",
			sql:  "DROP SEQUENCE seq1, db2.seq2",
			want: WalkThroughStatement{
				Kind:      WalkThroughStatementTableRefCheck,
				TableRefs: []TableRef{{Name: "seq1"}, {Schema: "db2", Name: "seq2"}},
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
		{name: "ctas temporary", sql: "CREATE TEMPORARY TABLE IF NOT EXISTS db.t2 AS SELECT 1;", want: true},
		{name: "ctas global temporary", sql: "CREATE GLOBAL TEMPORARY TABLE t2 AS SELECT 1;", want: true},
		{name: "ctas with comments", sql: "-- lead\nCREATE /* mid */ TABLE # hash\nt2 AS SELECT 1", want: true},
		{name: "regular create table", sql: "CREATE TABLE t1 (id INT);", want: false},
		{name: "create global without temporary", sql: "CREATE GLOBAL TABLE t1 AS SELECT 1;", want: false},
		{name: "select in single-quoted default", sql: "CREATE TABLE t1 (a VARCHAR(20) DEFAULT 'SELECT');", want: false},
		{name: "select in double-quoted comment", sql: `CREATE TABLE t1 (a INT COMMENT "SELECT 1");`, want: false},
		{name: "select in backtick identifier", sql: "CREATE TABLE t1 (`select` INT);", want: false},
		{name: "select in comments", sql: "CREATE TABLE t1 (a INT) /* SELECT */ -- SELECT\n# SELECT", want: false},
		{name: "select after semicolon", sql: "CREATE TABLE t1 (a INT); SELECT 1;", want: false},
		{name: "create view as select", sql: "CREATE VIEW v1 AS SELECT 1;", want: false},
		{name: "empty", sql: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCreateTableAs(tt.sql); got != tt.want {
				t.Errorf("IsCreateTableAs(%q) = %v, want %v", tt.sql, got, tt.want)
			}
		})
	}
}
