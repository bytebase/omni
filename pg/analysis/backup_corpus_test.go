package analysis

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type backupCorpusPos struct {
	Line   int `yaml:"line"`
	Column int `yaml:"column"`
}

type backupCorpusResult struct {
	Statement       string          `yaml:"statement"`
	SourceSchema    string          `yaml:"sourceschema"`
	SourceTableName string          `yaml:"sourcetablename"`
	TargetTableName string          `yaml:"targettablename"`
	StartPosition   backupCorpusPos `yaml:"startposition"`
	EndPosition     backupCorpusPos `yaml:"endposition"`
}

type backupCorpusCase struct {
	Input  string               `yaml:"input"`
	Result []backupCorpusResult `yaml:"result"`
}

type restoreCorpusCase struct {
	Input            string `yaml:"input"`
	BackupDatabase   string `yaml:"backupdatabase"`
	BackupTable      string `yaml:"backuptable"`
	OriginalDatabase string `yaml:"originaldatabase"`
	OriginalTable    string `yaml:"originaltable"`
	Result           string `yaml:"result"`
}

func fixedBackupMetadata() *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadata.SchemaMetadata{{
			Name: "public",
			Tables: []*metadata.TableMetadata{
				{
					Name:    "t_generated",
					Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}},
					Indexes: []*metadata.IndexMetadata{
						{Name: "t_generated_pk", Expressions: []string{"b"}, Primary: true, Unique: true},
						{Name: "t_generated_uk", Expressions: []string{"a"}, Unique: true, IsConstraint: true},
					},
				},
				{Name: "t1", Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}, {Name: "c"}}},
				{Name: "t2", Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}, {Name: "c"}}},
				{
					Name:    "test",
					Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}, {Name: "c"}},
					Indexes: []*metadata.IndexMetadata{
						{Name: "test_pk", Expressions: []string{"a"}, Primary: true, Unique: true},
					},
				},
				{Name: "test2", Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}, {Name: "c"}}},
				{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}, {Name: "c"}}},
				{Name: "test1", Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}, {Name: "c"}}},
			},
		}},
	}
}

func TestBackupCorpus(t *testing.T) {
	content, err := os.ReadFile("testdata/backup.yaml")
	require.NoError(t, err)
	var cases []backupCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	db := &protoDatabase{metadata: fixedBackupMetadata()}
	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			backups, err := TransformDMLToSelect(parseScript(t, tc.Input), tc.Input, BackupOptions{
				SearchPath: []string{"public"},
				SearchRelation: func(searchPath []string, name string) string {
					schema, _ := db.SearchRelation(searchPath, name)
					return schema
				},
				TargetSchema: "backupSchema",
				TablePrefix:  "rollback",
			})
			require.NoError(t, err)

			idx := review.Index(tc.Input)
			got := make([]backupCorpusResult, 0, len(backups))
			for _, b := range backups {
				startLine, startCol := idx.Position(b.StartOffset)
				endLine, endCol := idx.Position(b.EndOffset)
				got = append(got, backupCorpusResult{
					Statement:       b.Statement,
					SourceSchema:    b.SourceSchema,
					SourceTableName: b.SourceTableName,
					TargetTableName: b.TargetTableName,
					StartPosition:   backupCorpusPos{Line: startLine, Column: startCol - 1},
					EndPosition:     backupCorpusPos{Line: endLine, Column: endCol},
				})
			}
			slices.SortFunc(got, func(a, b backupCorpusResult) int {
				if a.TargetTableName != b.TargetTableName {
					return strings.Compare(a.TargetTableName, b.TargetTableName)
				}
				return strings.Compare(a.Statement, b.Statement)
			})
			require.Equal(t, tc.Result, got)
		})
	}
}

func TestBackupRejectsDataModifyingCTE(t *testing.T) {
	db := &protoDatabase{metadata: fixedBackupMetadata()}
	tests := []struct {
		name      string
		statement string
		wantErr   string
	}{
		{
			name:      "delete CTE feeding update",
			statement: "WITH d AS (DELETE FROM test2 WHERE a > 1 RETURNING a) UPDATE test SET c = 1 WHERE a IN (SELECT a FROM d);",
			wantErr:   `data-modifying WITH query "d";`,
		},
		{
			name:      "insert CTE feeding delete",
			statement: "WITH i AS (INSERT INTO test2 (a) VALUES (1) RETURNING a) DELETE FROM test USING i WHERE test.a = i.a;",
			wantErr:   `data-modifying WITH query "i";`,
		},
		{
			name:      "update CTE beside read-only CTE",
			statement: "WITH s AS (SELECT a FROM test2), u AS (UPDATE test2 SET a = a + 1 RETURNING a) UPDATE test SET c = 1 FROM s WHERE test.a = s.a;",
			wantErr:   `data-modifying WITH query "u";`,
		},
		{
			name:      "merge CTE",
			statement: "WITH m AS (MERGE INTO test2 USING test ON test2.a = test.a WHEN MATCHED THEN DELETE RETURNING test2.a) DELETE FROM test WHERE a IN (SELECT a FROM m);",
			wantErr:   `data-modifying WITH query "m";`,
		},
		{
			name:      "rejected after an earlier statement transforms",
			statement: "UPDATE test2 SET a = 1;\nWITH d AS (DELETE FROM test2 RETURNING a) DELETE FROM test WHERE a IN (SELECT a FROM d);",
			wantErr:   `data-modifying WITH query "d";`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := TransformDMLToSelect(parseScript(t, tc.statement), tc.statement, BackupOptions{
				SearchPath: []string{"public"},
				SearchRelation: func(searchPath []string, name string) string {
					schema, _ := db.SearchRelation(searchPath, name)
					return schema
				},
				TargetSchema: "backupSchema",
				TablePrefix:  "rollback",
			})
			require.ErrorContains(t, err, tc.wantErr)
			require.Nil(t, result)
		})
	}
}

func TestBackupAllowsReadOnlyCTE(t *testing.T) {
	db := &protoDatabase{metadata: fixedBackupMetadata()}
	for _, statement := range []string{
		"WITH v(a) AS (VALUES (1)) UPDATE test SET c = 1 FROM v WHERE test.a = v.a;",
		"WITH s AS (TABLE test2) DELETE FROM test USING s WHERE test.a = s.a;",
		"WITH s AS ((SELECT a FROM test2) UNION SELECT 1) UPDATE test SET c = 1 FROM s WHERE test.a = s.a;",
	} {
		result, err := TransformDMLToSelect(parseScript(t, statement), statement, BackupOptions{
			SearchPath: []string{"public"},
			SearchRelation: func(searchPath []string, name string) string {
				schema, _ := db.SearchRelation(searchPath, name)
				return schema
			},
			TargetSchema: "backupSchema",
			TablePrefix:  "rollback",
		})
		require.NoError(t, err, statement)
		require.Len(t, result, 1, statement)
	}
}

func TestRestoreCorpus(t *testing.T) {
	content, err := os.ReadFile("testdata/restore.yaml")
	require.NoError(t, err)
	var cases []restoreCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	db := &protoDatabase{metadata: fixedBackupMetadata()}
	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			stmts := parseScript(t, tc.Input)
			dmls := BackedUpDMLs(stmts)
			require.NotEmpty(t, dmls)

			var sqls []string
			var nodes []ast.Node
			for _, dml := range dmls {
				sqls = append(sqls, strings.TrimRight(strings.TrimSpace(dml.Text), ";"))
				nodes = append(nodes, dml.Node)
			}
			table := db.GetTable("public", tc.OriginalTable)
			require.NotNil(t, table)
			var pk *metadata.IndexMetadata
			for _, idx := range table.GetIndexes() {
				if idx.GetPrimary() {
					pk = idx
					break
				}
			}

			got, err := RestoreSQL(nodes, RestoreOptions{
				OriginalSchema:    "public",
				OriginalTable:     tc.OriginalTable,
				BackupSchema:      tc.BackupDatabase,
				BackupTable:       tc.BackupTable,
				Table:             table,
				PrimaryKey:        pk,
				OriginalSQL:       strings.Join(sqls, ";\n"),
				PrependStatements: PrependStatements(stmts),
			})
			require.NoError(t, err)
			require.Equal(t, tc.Result, got)
		})
	}
}
