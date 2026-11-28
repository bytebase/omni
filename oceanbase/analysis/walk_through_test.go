package analysis

import (
	"context"
	"reflect"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oceanbase/parser"
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
		{name: "ctas external", sql: "CREATE EXTERNAL TABLE t2 AS SELECT 1;", want: true},
		{name: "ctas as values row", sql: "CREATE TABLE t2 AS VALUES ROW(1);", want: true},
		{name: "ctas parenthesized values row", sql: "CREATE TABLE t2 (a INT) (VALUES ROW(1));", want: true},
		{name: "ctas with comments", sql: "-- lead\nCREATE /* mid */ TABLE # hash\nt2 AS SELECT 1", want: true},
		{name: "regular create table", sql: "CREATE TABLE t1 (id INT);", want: false},
		{name: "partitioned table with VALUES LESS THAN", sql: "CREATE TABLE t1 (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10));", want: false},
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

func walkThroughOceanBase(t *testing.T, script string) analysis.WalkThroughResult {
	t.Helper()
	statements, err := Engine.Parse(script)
	if err != nil {
		t.Fatalf("Parse(%q): %v", script, err)
	}
	database := &analysis.Database{
		Name: "test",
		Metadata: &metadata.DatabaseSchemaMetadata{Name: "test", Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{Name: "t1", Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int", Nullable: true}}}},
		}}},
	}
	return Engine.WalkThrough(context.Background(), analysis.WalkThroughRequest{Database: database, Script: script}, statements)
}

func TestEngineWalkThroughFindings(t *testing.T) {
	warning := func(code analysis.WalkThroughCode, content string, statement int, start analysis.Position) *analysis.WalkThroughFinding {
		return &analysis.WalkThroughFinding{Code: code, Title: content, Content: content, Statement: statement, Start: start}
	}
	failure := func(code analysis.WalkThroughCode, content string, line int) *analysis.WalkThroughFinding {
		return &analysis.WalkThroughFinding{Code: code, Error: true, Title: content, Content: content, Statement: -1, Start: analysis.Position{Line: line}}
	}
	tests := []struct {
		name   string
		script string
		want   *analysis.WalkThroughFinding
	}{
		{"another database", "CREATE DATABASE other;\nDROP TABLE missing;", warning(analysis.WalkThroughNotCurrentDatabase, "Database `other` is not the current database `test`", 0, analysis.Position{})},
		{"the current database dropped", "DROP DATABASE test;", warning(analysis.WalkThroughDatabaseIsDeleted, "Database `test` is deleted", 0, analysis.Position{})},
		{"create table as by its segment", "CREATE TABLE t2 (a INT);\nCREATE TABLE t3 AS SELECT 1;", warning(analysis.WalkThroughCreateTableAs, `CREATE TABLE AS statement is used in "CREATE TABLE t3 AS SELECT 1;"`, -1, analysis.Position{Line: 1, Column: 25})},
		{"a catalog error by its line", "CREATE TABLE t2 (a INT);\nALTER TABLE t1 ADD COLUMN a INT;", failure(analysis.WalkThroughColumnExists, "Column `a` already exists", 2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := walkThroughOceanBase(t, tt.script)
			if got.Database != nil || !reflect.DeepEqual(got.Finding, tt.want) {
				t.Fatalf("WalkThrough = %+v, %+v; want finding %+v", got.Database, got.Finding, tt.want)
			}
		})
	}
}

func TestEngineWalkThroughReturnsTheSchemaAfter(t *testing.T) {
	got := walkThroughOceanBase(t, "CREATE TABLE t2 (a INT); RENAME TABLE t1 TO other.t1;")
	if got.Finding != nil || got.Database == nil || got.Database.Name != "test" || !got.Database.ObjectCaseSensitive {
		t.Fatalf("WalkThrough = %+v", got)
	}
	tables := got.Database.Metadata.GetSchemas()[0].GetTables()
	if len(tables) != 1 || tables[0].GetName() != "t2" {
		t.Fatalf("tables = %v", tables)
	}
}
