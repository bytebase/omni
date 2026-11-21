package analysis

import (
	"os"
	"slices"
	"strings"
	"testing"
	"unicode/utf8"

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
	OriginalTable    string `yaml:"originaltable"`
	Result           string `yaml:"result"`
}

func offsetToZeroBasedColPos(sql string, offset int) backupCorpusPos {
	if offset > len(sql) {
		offset = len(sql)
	}
	line := 1
	col := 0
	for i := 0; i < offset; {
		r, size := utf8.DecodeRuneInString(sql[i:])
		if r == '\n' {
			line++
			col = 0
		} else {
			col++
		}
		i += size
	}
	return backupCorpusPos{Line: line, Column: col}
}

func TestBackupCorpus(t *testing.T) {
	content, err := os.ReadFile("testdata/backup.yaml")
	require.NoError(t, err)
	var cases []backupCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	opts := BackupOptions{SourceDatabase: "DB", TargetDatabase: "backupDB", TablePrefix: "rollback"}
	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			backups, err := TransformDMLToSelect(tc.Input, opts)
			require.NoError(t, err)
			got := make([]backupCorpusResult, 0, len(backups))
			for _, b := range backups {
				got = append(got, backupCorpusResult{
					Statement:       b.Statement,
					SourceSchema:    b.SourceSchema,
					SourceTableName: b.SourceTableName,
					TargetTableName: b.TargetTableName,
					StartPosition:   offsetToZeroBasedColPos(tc.Input, b.Start),
					EndPosition:     offsetToZeroBasedColPos(tc.Input, b.End),
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
		name    string
		input   string
		wantSQL string
	}{
		{
			name:  "update returning is not copied into select suffix",
			input: "UPDATE test SET c1 = 1 WHERE c1 = 2 RETURNING c1 INTO :old_c1;",
			wantSQL: `CREATE TABLE "backupDB"."rollback_TEST_DB" AS
  SELECT "TEST".* FROM test WHERE c1 = 2;`,
		},
		{
			name:  "delete returning is not copied into select suffix",
			input: "DELETE FROM test WHERE c1 = 2 RETURNING c1 INTO :old_c1;",
			wantSQL: `CREATE TABLE "backupDB"."rollback_TEST_DB" AS
  SELECT "TEST".* FROM test WHERE c1 = 2;`,
		},
		{
			name:  "update log errors is not copied into select suffix",
			input: "UPDATE test SET c1 = 1 WHERE c1 = 2 LOG ERRORS INTO err$_test REJECT LIMIT UNLIMITED;",
			wantSQL: `CREATE TABLE "backupDB"."rollback_TEST_DB" AS
  SELECT "TEST".* FROM test WHERE c1 = 2;`,
		},
		{
			name:  "where keyword in block comment",
			input: "UPDATE test SET c1 = 1 WHERE /* WHERE */ c1 = 2;",
			wantSQL: `CREATE TABLE "backupDB"."rollback_TEST_DB" AS
  SELECT "TEST".* FROM test WHERE /* WHERE */ c1 = 2;`,
		},
		{
			name:  "where keyword in string literal",
			input: "DELETE FROM test WHERE note = 'WHERE should stay literal';",
			wantSQL: `CREATE TABLE "backupDB"."rollback_TEST_DB" AS
  SELECT "TEST".* FROM test WHERE note = 'WHERE should stay literal';`,
		},
		{
			name:  "update without where backs up whole table",
			input: "UPDATE test SET c1 = 1;",
			wantSQL: `CREATE TABLE "backupDB"."rollback_TEST_DB" AS
  SELECT "TEST".* FROM test;`,
		},
		{
			name:  "delete without where backs up whole table",
			input: "DELETE FROM test;",
			wantSQL: `CREATE TABLE "backupDB"."rollback_TEST_DB" AS
  SELECT "TEST".* FROM test;`,
		},
		{
			name:  "update partition target",
			input: "UPDATE test PARTITION (p1) SET c1 = 1 WHERE c1 = 2;",
			wantSQL: `CREATE TABLE "backupDB"."rollback_TEST_DB" AS
  SELECT "TEST".* FROM test PARTITION (p1) WHERE c1 = 2;`,
		},
		{
			name:  "delete subpartition target",
			input: "DELETE FROM test SUBPARTITION (sp1) WHERE c1 = 2;",
			wantSQL: `CREATE TABLE "backupDB"."rollback_TEST_DB" AS
  SELECT "TEST".* FROM test SUBPARTITION (sp1) WHERE c1 = 2;`,
		},
		{
			name:  "update from appends source table list",
			input: "UPDATE test t SET c1 = s.c1 FROM source_table s WHERE t.id = s.id;",
			wantSQL: `CREATE TABLE "backupDB"."rollback_TEST_DB" AS
  SELECT "DB"."TEST".* FROM "DB"."TEST" WHERE ROWID IN (
    SELECT "T".ROWID FROM test t, source_table s WHERE t.id = s.id
  );`,
		},
	}

	opts := BackupOptions{SourceDatabase: "DB", TargetDatabase: "backupDB", TablePrefix: "rollback"}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := TransformDMLToSelect(tc.input, opts)
			require.NoError(t, err)
			require.Len(t, got, 1)
			require.Equal(t, tc.wantSQL, got[0].Statement)
		})
	}
}

func fixedRestoreCorpusTables() map[string]*metadata.TableMetadata {
	return map[string]*metadata.TableMetadata{
		"T_HIDDEN": {
			Name: "T_HIDDEN",
			Columns: []*metadata.ColumnMetadata{
				{Name: "A"},
				{Name: "H", IsInvisible: true},
				{Name: "V", Generation: &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_VIRTUAL, Expression: `"A"+1`}},
				{Name: "B"},
			},
		},
		"T_GENERATED": {
			Name:    "T_GENERATED",
			Columns: []*metadata.ColumnMetadata{{Name: "A"}, {Name: "B"}},
			Indexes: []*metadata.IndexMetadata{
				{Name: "T_GENERATED_PK", Expressions: []string{"B"}, Primary: true, Unique: true},
				{Name: "T_GENERATED_UK", Expressions: []string{"A"}, Unique: true},
			},
		},
		"TEST": {
			Name:    "TEST",
			Columns: []*metadata.ColumnMetadata{{Name: "A"}, {Name: "B"}, {Name: "C"}},
			Indexes: []*metadata.IndexMetadata{
				{Name: "TEST_PK", Expressions: []string{"A"}, Primary: true, Unique: true},
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
			var parts []string
			for _, stmt := range SplitStatements(tc.Input) {
				ok, err := ReferencesTable(stmt.Text, tc.OriginalDatabase, tc.OriginalTable)
				require.NoError(t, err)
				if ok {
					parts = append(parts, NormalizeRestoreSQL(stmt.Text))
				}
			}
			originalSQL := strings.Join(parts, "\n")
			node, err := FindFirstDML(originalSQL)
			require.NoError(t, err)
			require.NotNil(t, node)

			table := tables[tc.OriginalTable]
			require.NotNil(t, table)
			var pk *metadata.IndexMetadata
			for _, idx := range table.GetIndexes() {
				if idx.GetPrimary() {
					pk = idx
					break
				}
			}
			got, err := RestoreSQL(originalSQL, node, RestoreOptions{
				SourceDatabase: tc.OriginalDatabase,
				SourceTable:    tc.OriginalTable,
				BackupDatabase: tc.BackupDatabase,
				BackupTable:    tc.BackupTable,
				Table:          table,
				PrimaryKey:     pk,
			})
			require.NoError(t, err)
			require.Equal(t, tc.Result, got)
		})
	}
}
