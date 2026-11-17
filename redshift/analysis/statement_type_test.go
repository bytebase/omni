package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/redshift"
)

func TestGetStatementTypes(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		want      []StatementType
	}{
		{
			name:      "create table",
			statement: "CREATE TABLE t1 (id INT);",
			want: []StatementType{
				CreateTable,
			},
		},
		{
			name:      "create table as select",
			statement: "CREATE TABLE t1 AS SELECT 1;",
			want: []StatementType{
				CreateTable,
			},
		},
		{
			name:      "create view",
			statement: "CREATE VIEW v1 AS SELECT 1;",
			want: []StatementType{
				CreateView,
			},
		},
		{
			name:      "create materialized view",
			statement: "CREATE MATERIALIZED VIEW mv1 AS SELECT 1;",
			want: []StatementType{
				CreateView,
			},
		},
		{
			name:      "create index",
			statement: "CREATE INDEX idx_t1_id ON t1(id);",
			want: []StatementType{
				CreateIndex,
			},
		},
		{
			name:      "create sequence",
			statement: "CREATE SEQUENCE seq1;",
			want: []StatementType{
				CreateSequence,
			},
		},
		{
			name:      "create schema",
			statement: "CREATE SCHEMA s1;",
			want: []StatementType{
				CreateSchema,
			},
		},
		{
			name:      "create database",
			statement: "CREATE DATABASE db1;",
			want: []StatementType{
				CreateDatabase,
			},
		},
		{
			name:      "drop table",
			statement: "DROP TABLE t1;",
			want: []StatementType{
				DropTable,
			},
		},
		{
			name:      "drop view",
			statement: "DROP VIEW v1;",
			want: []StatementType{
				DropView,
			},
		},
		{
			name:      "drop materialized view",
			statement: "DROP MATERIALIZED VIEW mv1;",
			want: []StatementType{
				DropTable,
			},
		},
		{
			name:      "drop schema",
			statement: "DROP SCHEMA s1;",
			want: []StatementType{
				DropSchema,
			},
		},
		{
			name:      "alter table",
			statement: "ALTER TABLE t1 ADD COLUMN name TEXT;",
			want: []StatementType{
				AlterTable,
			},
		},
		{
			name:      "alter datashare",
			statement: "ALTER DATASHARE salesshare ADD SCHEMA public;",
			want:      []StatementType{},
		},
		{
			name:      "insert update delete",
			statement: "INSERT INTO t1 VALUES (1); UPDATE t1 SET id = 2; DELETE FROM t1;",
			want: []StatementType{
				Insert,
				Update,
				Delete,
			},
		},
		{
			name:      "merge",
			statement: "MERGE INTO t1 USING t2 ON t1.id = t2.id WHEN MATCHED THEN DELETE;",
			want: []StatementType{
				Merge,
			},
		},
		{
			name:      "truncate table",
			statement: "TRUNCATE TABLE t1;",
			want: []StatementType{
				Truncate,
			},
		},
		{
			name:      "comment",
			statement: "COMMENT ON TABLE t1 IS 'hello';",
			want: []StatementType{
				Comment,
			},
		},
		{
			name:      "call procedure",
			statement: "CALL archive_orders(30);",
			want: []StatementType{
				Execute,
			},
		},
		{
			name:      "execute prepared statement",
			statement: "EXECUTE purge_stmt(30);",
			want: []StatementType{
				Execute,
			},
		},
		{
			name:      "prepare",
			statement: "PREPARE purge_stmt(int) AS DELETE FROM t1 WHERE id = $1;",
			want:      []StatementType{},
		},
		{
			name:      "do block",
			statement: "DO $$ BEGIN DELETE FROM t1; END $$;",
			want:      []StatementType{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := redshift.Parse(tt.statement)
			require.NoError(t, err)
			var got []StatementType
			for _, stmt := range stmts {
				if stmt.Empty() {
					continue
				}
				statementType := StatementTypeOf(stmt.AST)
				if statementType == StatementTypeUnspecified {
					continue
				}
				got = append(got, statementType)
			}
			require.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestStatementTypeOfEveryStatement(t *testing.T) {
	statement := `CREATE TABLE t1 (id INT);
DROP TABLE t2;
INSERT INTO t1 VALUES (1);`
	stmts, err := redshift.Parse(statement)
	require.NoError(t, err)
	require.Len(t, stmts, 3)
	require.Equal(t, CreateTable, StatementTypeOf(stmts[0].AST))
	require.Equal(t, DropTable, StatementTypeOf(stmts[1].AST))
	require.Equal(t, Insert, StatementTypeOf(stmts[2].AST))
}
