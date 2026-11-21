package analysis_test

import (
	"os"
	"slices"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/bytebase/omni/mariadb/analysis"
	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/metadata"
)

// The corpora in testdata/backup.yaml and testdata/restore.yaml are Bytebase's
// MariaDB prior-backup and restore fixtures. They are reproduced byte-identically
// here and never re-recorded against this implementation.

type backupCorpusPos struct {
	Line   int32 `yaml:"line"`
	Column int32 `yaml:"column"`
}

type backupCorpusStatement struct {
	Statement       string           `yaml:"statement"`
	SourceSchema    string           `yaml:"sourceschema"`
	SourceTableName string           `yaml:"sourcetablename"`
	TargetTableName string           `yaml:"targettablename"`
	StartPosition   *backupCorpusPos `yaml:"startposition"`
	EndPosition     *backupCorpusPos `yaml:"endposition"`
}

type backupCorpusCase struct {
	Input  string                  `yaml:"input"`
	Result []backupCorpusStatement `yaml:"result"`
	Error  string                  `yaml:"error,omitempty"`
}

type restoreCorpusCase struct {
	Input            string `yaml:"input"`
	BackupDatabase   string `yaml:"backupdatabase"`
	BackupTable      string `yaml:"backuptable"`
	OriginalDatabase string `yaml:"originaldatabase"`
	OriginalTable    string `yaml:"originaltable"`
	Result           string `yaml:"result"`
	Error            string `yaml:"error,omitempty"`
}

func corpusBackupCatalog() analysis.Catalog {
	tables := []*metadata.TableMetadata{
		{
			Name: "t_generated",
			Columns: []*metadata.ColumnMetadata{
				{Name: "a"},
				{Name: "b"},
				{Name: "c_generated", Generation: &metadata.GenerationMetadata{Expression: "a + b"}},
			},
			Indexes: []*metadata.IndexMetadata{
				{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"b"}},
				{Name: "uk_a", Unique: true, Expressions: []string{"a"}},
			},
		},
		{
			Name:    "t1",
			Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}, {Name: "c"}},
			Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"a"}}},
		},
		{
			Name:    "t2",
			Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}, {Name: "c"}},
			Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"a"}}},
		},
		{
			Name:    "test",
			Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}, {Name: "c"}},
			Indexes: []*metadata.IndexMetadata{
				{Name: "PRIMARY", Primary: true, Expressions: []string{"c"}},
				{Name: "PRIMARY", Unique: true, Expressions: []string{"a"}},
			},
		},
		{
			Name:    "test2",
			Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}, {Name: "c"}},
			Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"a"}}},
		},
		{
			Name:    "t3",
			Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}, {Name: "d"}},
			Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"a"}}},
		},
	}
	return analysis.Catalog{
		ListDatabaseNames: func() ([]string, error) {
			return []string{"db", "db1", "db2"}, nil
		},
		GetDatabase: func(name string) (analysis.Database, error) {
			return &backupDatabase{name: name, tables: tables}, nil
		},
	}
}

func byteOffsetToCorpusPos(sql string, byteOffset int) *backupCorpusPos {
	if byteOffset < 0 {
		byteOffset = 0
	}
	if byteOffset > len(sql) {
		byteOffset = len(sql)
	}
	line := int32(1)
	col := int32(0)
	for i := 0; i < byteOffset; {
		r, size := utf8.DecodeRuneInString(sql[i:])
		switch r {
		case '\r':
			line++
			col = 0
		case '\n':
			if i == 0 || sql[i-1] != '\r' {
				line++
				col = 0
			}
		default:
			col++
		}
		i += size
	}
	return &backupCorpusPos{Line: line, Column: col + 1}
}

func TestBackupCorpus(t *testing.T) {
	content, err := os.ReadFile("testdata/backup.yaml")
	require.NoError(t, err)
	var cases []backupCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	for _, tc := range cases {
		catalog := corpusBackupCatalog()
		backups, err := analysis.TransformDMLToSelect(tc.Input, analysis.BackupOptions{
			DefaultDatabase: "db",
			BackupDatabase:  "backupDB",
			TablePrefix:     "_rollback",
			CaseSensitive:   false,
			Catalog:         catalog,
		})
		errText := ""
		if err != nil {
			errText = err.Error()
		}
		require.Equal(t, tc.Error, errText, tc.Input)
		if errText != "" {
			continue
		}

		got := make([]backupCorpusStatement, 0, len(backups))
		for _, b := range backups {
			got = append(got, backupCorpusStatement{
				Statement:       b.Statement,
				SourceTableName: b.SourceTable,
				TargetTableName: b.BackupTable,
				StartPosition:   byteOffsetToCorpusPos(tc.Input, b.ByteStart),
				EndPosition:     byteOffsetToCorpusPos(tc.Input, endWithSemicolon(tc.Input, b.ByteEnd)),
			})
		}
		slices.SortFunc(got, func(x, y backupCorpusStatement) int {
			if x.TargetTableName != y.TargetTableName {
				return strings.Compare(x.TargetTableName, y.TargetTableName)
			}
			return strings.Compare(x.Statement, y.Statement)
		})
		require.Equal(t, tc.Result, got, tc.Input)
	}
}

func TestRestoreCorpus(t *testing.T) {
	content, err := os.ReadFile("testdata/restore.yaml")
	require.NoError(t, err)
	var cases []restoreCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	for _, tc := range cases {
		catalog := corpusBackupCatalog()
		var matched []string
		for _, seg := range parser.Split(tc.Input) {
			stmtText := tc.Input[seg.ByteStart:endWithSemicolon(tc.Input, seg.ByteEnd)]
			ok, err := analysis.ReferencesTable(stmtText, tc.OriginalDatabase, tc.OriginalTable)
			require.NoError(t, err, tc.Input)
			if ok {
				matched = append(matched, stmtText)
			}
		}
		originalSQL := strings.Join(matched, "")
		nodes, err := analysis.DMLs(originalSQL)
		require.NoError(t, err, tc.Input)
		require.NotEmpty(t, nodes, tc.Input)

		columns, err := analysis.ClassifyColumns(catalog, false, tc.OriginalDatabase, tc.OriginalTable)
		require.NoError(t, err, tc.Input)

		got, err := analysis.RestoreSQL(nodes, analysis.RestoreOptions{
			Statement:      originalSQL,
			Database:       tc.OriginalDatabase,
			Table:          tc.OriginalTable,
			BackupDatabase: tc.BackupDatabase,
			BackupTable:    tc.BackupTable,
			Columns:        columns.Normal,
			Catalog:        catalog,
			CaseSensitive:  false,
		})
		errText := ""
		if err != nil {
			errText = err.Error()
		}
		require.Equal(t, tc.Error, errText, tc.Input)
		require.Equal(t, tc.Result, got, tc.Input)
	}
}

func endWithSemicolon(sql string, byteEnd int) int {
	if byteEnd < len(sql) && sql[byteEnd] == ';' {
		return byteEnd + 1
	}
	return byteEnd
}
