package analysis

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
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
	OriginalSchema   string `yaml:"originalschema"`
	OriginalTable    string `yaml:"originaltable"`
	Result           string `yaml:"result"`
}

func TestBackupCorpus(t *testing.T) {
	content, err := os.ReadFile("testdata/backup.yaml")
	require.NoError(t, err)
	var cases []backupCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			backups, err := TransformDMLToSelect(parseScript(t, tc.Input), "db", "backupDB", "rollback", false)
			require.NoError(t, err)

			got := make([]backupCorpusResult, 0, len(backups))
			for _, b := range backups {
				got = append(got, backupCorpusResult{
					Statement:       b.Statement,
					SourceSchema:    b.SourceSchema,
					SourceTableName: b.SourceTableName,
					TargetTableName: b.TargetTableName,
					StartPosition:   backupCorpusPos{Line: b.StartPosition.Line, Column: b.StartPosition.Column},
					EndPosition:     backupCorpusPos{Line: b.EndPosition.Line, Column: b.EndPosition.Column},
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

func TestBackupBoundaryCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantSQL     string
		wantSQLs    []string
		wantErrPart string
	}{
		{
			name:  "update top",
			input: "UPDATE TOP (3) test SET c1 = 1 WHERE c2 = 2;",
			wantSQL: strings.Join([]string{
				"SELECT * INTO [backupDB].[dbo].[rollback_test_db] FROM (",
				"  SELECT TOP (3) [db].[dbo].[test].* FROM test WHERE c2 = 2) AS backup_table;",
			}, "\n"),
		},
		{
			name:  "update top option",
			input: "UPDATE TOP (3) test SET c1 = 1 WHERE c2 = 2 OPTION (RECOMPILE);",
			wantSQL: strings.Join([]string{
				"SELECT * INTO [backupDB].[dbo].[rollback_test_db] FROM (",
				"  SELECT TOP (3) [db].[dbo].[test].* FROM test WHERE c2 = 2 OPTION (RECOMPILE)) AS backup_table;",
			}, "\n"),
		},
		{
			name:  "where keyword in comment",
			input: "UPDATE test SET c1 = 1 WHERE /* WHERE */ c2 = 2;",
			wantSQL: strings.Join([]string{
				"SELECT * INTO [backupDB].[dbo].[rollback_test_db] FROM (",
				"  SELECT [db].[dbo].[test].* FROM test WHERE /* WHERE */ c2 = 2) AS backup_table;",
			}, "\n"),
		},
		{
			name:  "option keyword in comment",
			input: "UPDATE test SET c1 = 1 /* OPTION */ OPTION (RECOMPILE);",
			wantSQL: strings.Join([]string{
				"SELECT * INTO [backupDB].[dbo].[rollback_test_db] FROM (",
				"  SELECT [db].[dbo].[test].* FROM test OPTION (RECOMPILE)) AS backup_table;",
			}, "\n"),
		},
		{
			name:  "option keyword in nested comment",
			input: "UPDATE test SET c1 = 1 /* outer /* inner */ OPTION in outer */ OPTION (RECOMPILE);",
			wantSQL: strings.Join([]string{
				"SELECT * INTO [backupDB].[dbo].[rollback_test_db] FROM (",
				"  SELECT [db].[dbo].[test].* FROM test OPTION (RECOMPILE)) AS backup_table;",
			}, "\n"),
		},
		{
			name: "delete alias from join",
			input: strings.Join([]string{
				"DELETE FROM t_alias",
				"FROM test AS t_alias JOIN test2 AS t2 ON t_alias.c1 = t2.c1",
				"WHERE t_alias.c1 = 1;",
			}, "\n"),
			wantSQL: strings.Join([]string{
				"SELECT * INTO [backupDB].[dbo].[rollback_test_db] FROM (",
				"  SELECT [t_alias].* FROM test AS t_alias JOIN test2 AS t2 ON t_alias.c1 = t2.c1 WHERE t_alias.c1 = 1) AS backup_table;",
			}, "\n"),
		},
		{
			name:        "update current of rejected",
			input:       "UPDATE test SET c1 = 1 WHERE CURRENT OF my_cursor;",
			wantErrPart: "CURSOR clause is not supported",
		},
		{
			name:  "update with cte",
			input: "WITH c AS (SELECT id FROM src) UPDATE dbo.test SET c1 = 1 WHERE EXISTS (SELECT 1 FROM c WHERE c.id = test.id);",
			wantSQL: strings.Join([]string{
				"WITH c AS (SELECT id FROM src)",
				"SELECT * INTO [backupDB].[dbo].[rollback_test_db] FROM (",
				"  SELECT [db].[dbo].[test].* FROM dbo.test WHERE EXISTS (SELECT 1 FROM c WHERE c.id = test.id)) AS backup_table;",
			}, "\n"),
		},
		{
			name:  "delete with cte",
			input: "WITH c (id) AS (SELECT id FROM src), d AS (SELECT 1 AS n) DELETE FROM test WHERE id IN (SELECT id FROM c);",
			wantSQL: strings.Join([]string{
				"WITH c (id) AS (SELECT id FROM src), d AS (SELECT 1 AS n)",
				"SELECT * INTO [backupDB].[dbo].[rollback_test_db] FROM (",
				"  SELECT [db].[dbo].[test].* FROM test WHERE id IN (SELECT id FROM c)) AS backup_table;",
			}, "\n"),
		},
		{
			name:  "update joined to cte",
			input: "WITH c AS (SELECT id FROM src) UPDATE t SET c1 = 1 FROM test AS t JOIN c ON t.id = c.id;",
			wantSQL: strings.Join([]string{
				"WITH c AS (SELECT id FROM src)",
				"SELECT * INTO [backupDB].[dbo].[rollback_test_db] FROM (",
				"  SELECT [t].* FROM test AS t JOIN c ON t.id = c.id) AS backup_table;",
			}, "\n"),
		},
		{
			name: "identical with clause across statements is hoisted once",
			input: strings.Join([]string{
				"WITH c AS (SELECT id FROM src) UPDATE test SET c1 = 1 WHERE id IN (SELECT id FROM c);",
				"WITH c AS (SELECT id FROM src) UPDATE test SET c1 = 2 WHERE id IN (SELECT id FROM c);",
			}, "\n"),
			wantSQL: strings.Join([]string{
				"WITH c AS (SELECT id FROM src)",
				"SELECT * INTO [backupDB].[dbo].[rollback_test_db] FROM (",
				"  SELECT [db].[dbo].[test].* FROM test WHERE id IN (SELECT id FROM c)",
				"  UNION",
				"  SELECT [db].[dbo].[test].* FROM test WHERE id IN (SELECT id FROM c)) AS backup_table;",
			}, "\n"),
		},
		{
			name: "different with clauses on the same table rejected",
			input: strings.Join([]string{
				"WITH c AS (SELECT id FROM src) UPDATE test SET c1 = 1 WHERE id IN (SELECT id FROM c);",
				"WITH c AS (SELECT id FROM src), d AS (SELECT id FROM other) UPDATE test SET c1 = 2 WHERE id IN (SELECT id FROM d);",
			}, "\n"),
			wantErrPart: "different WITH clauses on the same table",
		},
		{
			name: "with clause would capture a base table name in another statement",
			input: strings.Join([]string{
				"WITH c AS (SELECT id FROM src) UPDATE test SET c1 = 1 WHERE id IN (SELECT id FROM c);",
				"UPDATE test SET c1 = 2 WHERE id IN (SELECT id FROM c);",
			}, "\n"),
			wantErrPart: "different WITH clauses on the same table",
		},
		{
			name: "with clause on another table is scoped to its own backup",
			input: strings.Join([]string{
				"WITH c AS (SELECT id FROM src) UPDATE test SET c1 = 1 WHERE id IN (SELECT id FROM c);",
				"UPDATE c SET c1 = 2 WHERE id = 1;",
			}, "\n"),
			wantSQLs: []string{
				strings.Join([]string{
					"SELECT * INTO [backupDB].[dbo].[rollback_c_db] FROM (",
					"  SELECT [db].[dbo].[c].* FROM c WHERE id = 1) AS backup_table;",
				}, "\n"),
				strings.Join([]string{
					"WITH c AS (SELECT id FROM src)",
					"SELECT * INTO [backupDB].[dbo].[rollback_test_db] FROM (",
					"  SELECT [db].[dbo].[test].* FROM test WHERE id IN (SELECT id FROM c)) AS backup_table;",
				}, "\n"),
			},
		},
		{
			name:        "xmlnamespaces rejected",
			input:       "WITH XMLNAMESPACES ('uri' AS ns), c AS (SELECT id FROM src) UPDATE test SET c1 = 1 WHERE id IN (SELECT id FROM c);",
			wantErrPart: "WITH XMLNAMESPACES",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := TransformDMLToSelect(parseScript(t, tc.input), "db", "backupDB", "rollback", false)
			if tc.wantErrPart != "" {
				require.ErrorContains(t, err, tc.wantErrPart)
				return
			}
			require.NoError(t, err)
			if tc.wantSQLs != nil {
				var got []string
				for _, item := range result {
					got = append(got, item.Statement)
				}
				slices.Sort(got)
				require.Equal(t, tc.wantSQLs, got)
				return
			}
			require.Len(t, result, 1)
			require.Equal(t, tc.wantSQL, result[0].Statement)
		})
	}
}

func TestBackupRejectsCTETarget(t *testing.T) {
	for name, input := range map[string]string{
		"direct":          "WITH c AS (SELECT id, c1 FROM test) UPDATE c SET c1 = 1 WHERE id = 1;",
		"alias":           "WITH c AS (SELECT id, c1 FROM test) UPDATE x SET c1 = 1 FROM c AS x WHERE x.id = 1;",
		"delete":          "WITH c AS (SELECT id FROM test) DELETE FROM c WHERE id = 1;",
		"mixed case":      "WITH C AS (SELECT id, c1 FROM test) UPDATE c SET c1 = 1 WHERE id = 1;",
		"after other dml": "DELETE FROM test WHERE id = 2; WITH c AS (SELECT id FROM test) UPDATE c SET c1 = 1 WHERE id = 1;",
	} {
		t.Run(name, func(t *testing.T) {
			result, err := TransformDMLToSelect(parseScript(t, input), "db", "backupDB", "rollback", false)
			require.ErrorContains(t, err, `does not support DML targeting CTE "c"`)
			require.Empty(t, result)
		})
	}

	result, err := TransformDMLToSelect(parseScript(t, "WITH c AS (SELECT id FROM src) UPDATE dbo.c SET c1 = 1 WHERE id IN (SELECT id FROM c);"), "db", "backupDB", "rollback", false)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "c", result[0].SourceTableName)
	require.True(t, strings.HasPrefix(result[0].Statement, "WITH c AS (SELECT id FROM src)\n"))

	result, err = TransformDMLToSelect(parseScript(t, "WITH C AS (SELECT id FROM src) UPDATE c SET c1 = 1 WHERE id IN (SELECT id FROM C);"), "db", "backupDB", "rollback", true)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "c", result[0].SourceTableName)

	_, err = TransformDMLToSelect(parseScript(t, "WITH c AS (SELECT id, c1 FROM test) UPDATE c SET c1 = 1 WHERE id = 1;"), "db", "backupDB", "rollback", true)
	require.ErrorContains(t, err, `does not support DML targeting CTE "c"`)
}

func TestBackupSkipsTempTableTargets(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   string
		wantLen int
	}{
		{"update local temp table", "UPDATE #tmp SET c1 = 1 WHERE c2 = 2;", 0},
		{"delete local temp table", "DELETE FROM #tmp WHERE c1 = 1;", 0},
		{"update global temp table", "UPDATE ##gtmp SET c1 = 1 WHERE c2 = 2;", 0},
		{"mixed temp and real", "UPDATE #tmp SET c1 = 1 WHERE c2 = 2;\nUPDATE test SET c1 = 2 WHERE c1 = 1;", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := TransformDMLToSelect(parseScript(t, tc.input), "db", "backupDB", "rollback", false)
			require.NoError(t, err)
			require.Len(t, result, tc.wantLen)
		})
	}
}

func fixedRestoreCorpusTables() map[string]*metadata.TableMetadata {
	return map[string]*metadata.TableMetadata{
		"t_generated": {
			Name:    "t_generated",
			Columns: columns("a", "b"),
			Indexes: []*metadata.IndexMetadata{
				{Name: "t_generated_pk", Expressions: []string{"b"}, Primary: true, Unique: true},
				{Name: "t_generated_uk", Expressions: []string{"a"}, Unique: true},
			},
		},
		"test": {
			Name:    "test",
			Columns: columns("a", "b", "c"),
			Indexes: []*metadata.IndexMetadata{
				{Name: "test_pk", Expressions: []string{"a"}, Primary: true, Unique: true},
			},
		},
	}
}

func TestRestoreCorpus(t *testing.T) {
	content, err := os.ReadFile("testdata/restore.yaml")
	require.NoError(t, err)
	var cases []restoreCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	tables := fixedRestoreCorpusTables()
	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			dmls := BackedUpDMLs(parseScript(t, tc.Input))
			require.NotEmpty(t, dmls)
			var sqls []string
			for _, dml := range dmls {
				sqls = append(sqls, dml.SQL)
			}
			table := tables[tc.OriginalTable]
			require.NotNil(t, table)
			var pk *metadata.IndexMetadata
			for _, idx := range table.GetIndexes() {
				if idx.GetPrimary() {
					pk = idx
					break
				}
			}
			got, err := RestoreSQL(dmls[0].Node, RestoreOptions{
				OriginalDatabase: tc.OriginalDatabase,
				OriginalSchema:   "dbo",
				OriginalTable:    tc.OriginalTable,
				BackupDatabase:   tc.BackupDatabase,
				BackupTable:      tc.BackupTable,
				Table:            table,
				PrimaryKey:       pk,
				OriginalSQL:      strings.Join(sqls, ";\n") + ";",
			})
			require.NoError(t, err)
			require.Equal(t, tc.Result, got)
		})
	}
}

func TestRestoreVariableAssignments(t *testing.T) {
	tables := fixedRestoreCorpusTables()
	table := tables["test"]
	pk := table.Indexes[0]

	for _, tc := range []struct {
		name    string
		sql     string
		wantSet string
	}{
		{
			name:    "dual assignment restores assigned column",
			sql:     "UPDATE test SET @v = b = 1 WHERE a = 1;",
			wantSet: "UPDATE SET t.[b] = b.[b]",
		},
		{
			name:    "variable read assignment does not restore read column",
			sql:     "UPDATE test SET @v = b, c = c + 1 WHERE a = 1;",
			wantSet: "UPDATE SET t.[c] = b.[c]",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dmls := BackedUpDMLs(parseScript(t, tc.sql))
			require.Len(t, dmls, 1)
			got, err := RestoreSQL(dmls[0].Node, RestoreOptions{
				OriginalDatabase: "db",
				OriginalSchema:   "dbo",
				OriginalTable:    "test",
				BackupDatabase:   "bbarchive",
				BackupTable:      "prefix_1_test",
				Table:            table,
				PrimaryKey:       pk,
				OriginalSQL:      tc.sql,
			})
			require.NoError(t, err)
			require.Contains(t, got, tc.wantSet)
		})
	}
}
