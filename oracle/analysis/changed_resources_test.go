package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/oracle/ast"
)

func extractAll(t *testing.T, script string, indexTable IndexTableFunc) StatementChanges {
	t.Helper()
	var all StatementChanges
	for _, statement := range SplitStatements(script) {
		if statement.Empty {
			continue
		}
		list, err := Parse(statement.Text)
		require.NoError(t, err)
		for _, item := range list.Items {
			raw, ok := item.(*ast.RawStmt)
			require.True(t, ok)
			changes := ExtractChangedResources(raw.Stmt, "DB", statement.Text, indexTable)
			all.Resources = append(all.Resources, changes.Resources...)
			all.DMLStatements = append(all.DMLStatements, changes.DMLStatements...)
			all.DMLCount += changes.DMLCount
			all.InsertCount += changes.InsertCount
		}
	}
	return all
}

func TestExtractChangedResources(t *testing.T) {
	tests := []struct {
		name   string
		script string
		want   StatementChanges
	}{
		{
			name:   "create, view and insert values",
			script: "CREATE TABLE t1 (c1 INT);\nCREATE VIEW hello AS SELECT * FROM world;\nINSERT INTO T1 (c1) VALUES (1);",
			want: StatementChanges{
				Resources:   []ChangedResource{{Database: "DB", Table: "T1"}, {Database: "DB", Table: "T1"}},
				InsertCount: 1,
			},
		},
		{
			name:   "drop and truncate affect data",
			script: "DROP TABLE IF EXISTS s.t2;\nTRUNCATE TABLE t3;",
			want: StatementChanges{
				Resources: []ChangedResource{{Database: "S", Table: "T2", AffectsData: true}, {Database: "DB", Table: "T3", AffectsData: true}},
			},
		},
		{
			name:   "alter table rename does not affect data",
			script: "ALTER TABLE t1 RENAME TO t2;\nALTER TABLE t1 ADD (c2 INT);",
			want: StatementChanges{
				Resources: []ChangedResource{{Database: "DB", Table: "T1"}, {Database: "DB", Table: "T1", AffectsData: true}},
			},
		},
		{
			name:   "DML samples",
			script: "INSERT INTO t1 SELECT * FROM t2;\nUPDATE t1 SET c1 = 5 WHERE c2 = 1;\nDELETE FROM t1  \n;\nMERGE INTO t1 USING t2 ON (t1.a = t2.a) WHEN MATCHED THEN UPDATE SET b = 1;",
			want: StatementChanges{
				Resources: []ChangedResource{{Database: "DB", Table: "T1"}, {Database: "DB", Table: "T1"}, {Database: "DB", Table: "T1"}, {Database: "DB", Table: "T1"}},
				DMLStatements: []string{
					"INSERT INTO t1 SELECT * FROM t2;",
					"UPDATE t1 SET c1 = 5 WHERE c2 = 1;",
					"DELETE FROM t1;",
					"MERGE INTO t1 USING t2 ON (t1.a = t2.a) WHEN MATCHED THEN UPDATE SET b = 1;",
				},
				DMLCount: 4,
			},
		},
		{
			name:   "multi-table insert",
			script: "INSERT ALL INTO t1 VALUES (1) INTO s.t2 VALUES (2) SELECT * FROM DUAL;",
			want: StatementChanges{
				Resources:     []ChangedResource{{Database: "DB", Table: "T1"}, {Database: "S", Table: "T2"}},
				DMLStatements: []string{"INSERT ALL INTO t1 VALUES (1) INTO s.t2 VALUES (2) SELECT * FROM DUAL;"},
				DMLCount:      1,
			},
		},
		{
			name:   "object DDL records only a qualified schema",
			script: "CREATE OR REPLACE VIEW app.v AS SELECT 1 FROM DUAL;\nCREATE SEQUENCE seq;\nDROP PROCEDURE app2.p;\nDROP SYNONYM syn;\nCREATE SYNONYM app3.syn FOR t;\nDROP PACKAGE app4.pkg;\nDROP TYPE app5.ty;\nALTER PACKAGE app6.pkg COMPILE;",
			want: StatementChanges{
				Resources: []ChangedResource{
					{Database: "APP", DatabaseOnly: true},
					{Database: "APP2", DatabaseOnly: true},
					{Database: "APP3", DatabaseOnly: true},
					{Database: "APP4", DatabaseOnly: true},
					{Database: "APP5", DatabaseOnly: true},
					{Database: "APP6", DatabaseOnly: true},
				},
			},
		},
		{
			name:   "create table if not exists",
			script: "CREATE TABLE IF NOT EXISTS t1 (c1 INT);",
			want: StatementChanges{
				Resources: []ChangedResource{{Database: "DB", Table: "T1"}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, extractAll(t, test.script, nil))
		})
	}
}

func TestExtractChangedResourcesDropIndex(t *testing.T) {
	indexTable := func(schema, index string) (string, bool) {
		if schema == "DB" && index == "IDX" {
			return "T1", true
		}
		return "", false
	}
	require.Equal(t, StatementChanges{Resources: []ChangedResource{{Database: "DB", Table: "T1"}}}, extractAll(t, "DROP INDEX idx;\nDROP INDEX other;", indexTable))
	require.Equal(t, StatementChanges{}, extractAll(t, "DROP INDEX idx;", nil))
}
