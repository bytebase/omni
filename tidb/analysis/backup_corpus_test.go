package analysis_test

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/tidb/analysis"
	"github.com/bytebase/omni/tidb/parser"
)

// The corpora in testdata/backup.yaml and testdata/restore.yaml are Bytebase's
// TiDB prior-backup and restore fixtures. They are reproduced byte-identically
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
				{Name: "uk_c_generated", Unique: true, Expressions: []string{"c_generated"}},
				{Name: "uk_empty_expressions", Unique: true, Expressions: nil},
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

func corpusPos(index *review.Text, byteOffset int) *backupCorpusPos {
	line, col := index.Position(byteOffset)
	return &backupCorpusPos{Line: int32(line), Column: int32(col)}
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
		require.NoError(t, err, tc.Input)

		index := review.Index(tc.Input)
		got := make([]backupCorpusStatement, 0, len(backups))
		for _, b := range backups {
			got = append(got, backupCorpusStatement{
				Statement:       b.Statement,
				SourceTableName: b.SourceTable,
				TargetTableName: b.BackupTable,
				StartPosition:   corpusPos(index, b.ByteStart),
				EndPosition:     corpusPos(index, endWithSemicolon(tc.Input, b.ByteEnd)),
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
		targetCols, err := analysis.NormalColumnSet(catalog, tc.OriginalDatabase, tc.OriginalTable)
		require.NoError(t, err, tc.Input)

		var matched []string
		for _, seg := range parser.Split(tc.Input) {
			stmtText := tc.Input[seg.ByteStart:endWithSemicolon(tc.Input, seg.ByteEnd)]
			ok, err := analysis.MutatesTable(stmtText, tc.OriginalDatabase, tc.OriginalTable, targetCols)
			require.NoError(t, err, tc.Input)
			if ok {
				matched = append(matched, stmtText)
			}
		}
		originalSQL := strings.Join(matched, "")
		nodes, err := analysis.MatchingDMLs(originalSQL, tc.OriginalDatabase, tc.OriginalTable, targetCols)
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
		})
		errText := ""
		if err != nil {
			errText = err.Error()
		}
		require.Equal(t, tc.Error, errText, tc.Input)
		require.Equal(t, tc.Result, got, tc.Input)
	}
}

// TestBackupRejectsAndPreserves pins TiDB prior-backup guards and AST preservation.
func TestBackupRejectsAndPreserves(t *testing.T) {
	a := require.New(t)
	catalog := corpusBackupCatalog()
	run := func(sql string) ([]analysis.BackupStatement, error) {
		return analysis.TransformDMLToSelect(sql, analysis.BackupOptions{
			DefaultDatabase: "db",
			BackupDatabase:  "backupDB",
			TablePrefix:     "_rollback",
			CaseSensitive:   false,
			Catalog:         catalog,
		})
	}

	_, err := run("BATCH LIMIT 2 DELETE FROM test WHERE id > 0")
	a.Error(err, "BATCH (non-transactional DML) must be rejected")

	_, err = run("DELETE FROM db1.t1 WHERE a = 1")
	a.Error(err, "cross-database DELETE must be rejected")

	_, err = run("UPDATE db1.t1 SET c1 = 1 WHERE c1 = 2")
	a.Error(err, "cross-database UPDATE must be rejected")

	result, err := run("UPDATE test, (SELECT id FROM test2) AS x SET test.c1 = 1 WHERE test.id = x.id")
	a.NoError(err)
	a.Len(result, 1)
	a.Equal(
		"CREATE TABLE `backupDB`.`_rollback_0_test` LIKE `db`.`test`;\n"+
			"INSERT INTO `backupDB`.`_rollback_0_test` (`a`,`b`,`c`) SELECT DISTINCT `test`.`a`,`test`.`b`,`test`.`c` FROM test, (SELECT id FROM test2) AS x WHERE test.id = x.id;",
		result[0].Statement,
	)

	result, err = run("WITH x AS (SELECT id FROM test2) UPDATE test JOIN x ON test.id = x.id SET test.c1 = 1")
	a.NoError(err)
	a.Len(result, 1)
	a.Equal(
		"CREATE TABLE `backupDB`.`_rollback_0_test` LIKE `db`.`test`;\n"+
			"INSERT INTO `backupDB`.`_rollback_0_test` (`a`,`b`,`c`) WITH x AS (SELECT id FROM test2) SELECT DISTINCT `test`.`a`,`test`.`b`,`test`.`c` FROM test JOIN x ON test.id = x.id;",
		result[0].Statement,
	)
	_, perr := parser.Parse(result[0].Statement)
	a.NoError(perr, "generated CTE backup must re-parse as valid SQL")

	manyWithCTE := "DELETE FROM test WHERE id = 1;\n" +
		"DELETE FROM test WHERE id = 2;\n" +
		"DELETE FROM test WHERE id = 3;\n" +
		"DELETE FROM test WHERE id = 4;\n" +
		"DELETE FROM test WHERE id = 5;\n" +
		"WITH x AS (SELECT id FROM test2) DELETE FROM test WHERE id > 0;"
	_, err = run(manyWithCTE)
	a.Error(err, "CTE in a >maxMixedDMLCount same-table batch must be rejected")

	result, err = run("WITH x AS (SELECT id FROM test2) UPDATE test JOIN x ON test.id = x.id SET nonexistent_col = 1")
	a.NoError(err)
	a.Len(result, 1)
	a.Equal("test", result[0].SourceTable)

	result, err = run("WITH x AS (SELECT id FROM test2) UPDATE test AS x JOIN test2 AS y ON x.a = y.a SET x.c = 1")
	a.NoError(err)
	a.Len(result, 1)
	a.Equal("test", result[0].SourceTable)

	result, err = run("WITH x AS (SELECT id FROM test2) UPDATE test AS x JOIN test2 AS y ON x.a = y.a SET x.c = 1, y.c = 2")
	a.NoError(err)
	a.Len(result, 2)
	a.ElementsMatch([]string{"test", "test2"}, []string{result[0].SourceTable, result[1].SourceTable})

	result, err = run("WITH x AS (SELECT id FROM test2) DELETE x, y FROM test AS x JOIN test2 AS y WHERE x.a = y.a")
	a.NoError(err)
	a.Len(result, 2)
	a.ElementsMatch([]string{"test", "test2"}, []string{result[0].SourceTable, result[1].SourceTable})

	_, err = run("WITH x AS (SELECT id FROM test2) UPDATE test JOIN x ON test.id = x.id SET x.c = 1")
	a.Error(err, "updating only a CTE must error")

	result, err = run("WITH X AS (SELECT id FROM test2) UPDATE test JOIN x ON test.id = x.id SET missing_col = 1")
	a.NoError(err)
	a.Len(result, 1)
	a.Equal("test", result[0].SourceTable)

	result, err = run("WITH test AS (SELECT a FROM test2) UPDATE db.test AS p JOIN test2 AS y ON p.a = y.a SET p.c = 1, y.c = 2")
	a.NoError(err)
	a.Len(result, 2)
	a.ElementsMatch([]string{"test", "test2"}, []string{result[0].SourceTable, result[1].SourceTable})
	for _, r := range result {
		_, perr := parser.Parse(r.Statement)
		a.NoError(perr)
	}

	result, err = run("WITH test AS (SELECT a FROM test2) DELETE p, y FROM db.test AS p JOIN test2 AS y WHERE p.a = y.a")
	a.NoError(err)
	a.Len(result, 2)
	a.ElementsMatch([]string{"test", "test2"}, []string{result[0].SourceTable, result[1].SourceTable})
	for _, r := range result {
		_, perr := parser.Parse(r.Statement)
		a.NoError(perr)
	}

	result, err = run("WITH test AS (SELECT a FROM test2) DELETE db.test, test2 FROM db.test JOIN test2 ON db.test.a = test2.a CROSS JOIN test")
	a.NoError(err)
	a.Len(result, 2)
	a.ElementsMatch([]string{"test", "test2"}, []string{result[0].SourceTable, result[1].SourceTable})

	result, err = run("UPDATE DB.test SET c1 = 1 WHERE c1 = 2")
	a.NoError(err)
	a.Len(result, 1)
	a.Equal("test", result[0].SourceTable)

	_, err = run("EXPLAIN ANALYZE UPDATE test SET c1 = 1 WHERE c1 = 2")
	a.Error(err)
	result, err = run("EXPLAIN UPDATE test SET c1 = 1 WHERE c1 = 2")
	a.NoError(err)
	a.Empty(result)

	_, err = run("REPLACE INTO test VALUES (1, 2, 3)")
	a.Error(err)
	_, err = run("INSERT INTO test VALUES (1, 2, 3) ON DUPLICATE KEY UPDATE c = 5")
	a.Error(err)
	result, err = run("INSERT INTO test VALUES (1, 2, 3)")
	a.NoError(err)
	a.Empty(result)

	_, err = run("EXPLAIN ANALYZE REPLACE INTO test VALUES (1, 2, 3)")
	a.Error(err)
	_, err = run("EXPLAIN ANALYZE INSERT INTO test VALUES (1, 2, 3) ON DUPLICATE KEY UPDATE c = 5")
	a.Error(err)
	result, err = run("EXPLAIN ANALYZE INSERT INTO test VALUES (1, 2, 3)")
	a.NoError(err)
	a.Empty(result)

	for _, sql := range []string{
		"UPDATE (test AS a JOIN test2 AS b ON a.a = b.a) SET a.c = 1",
		"DELETE a FROM (test AS a JOIN test2 AS b ON a.a = b.a) WHERE a.a = 1",
		"UPDATE (test AS t JOIN test2 AS y ON t.a = y.a) JOIN t1 AS z ON z.a = t.a SET t.c = 1 WHERE z.a = 1",
		"UPDATE t1 AS z JOIN (test AS t JOIN test2 AS y ON t.a = y.a) ON z.a = t.a SET z.c = 1 WHERE t.a = 1",
		"UPDATE ( /* c ) */ test AS a LEFT JOIN test2 AS b ON a.a = b.a) JOIN t1 AS z ON z.a = a.a SET a.c = 1",
	} {
		result, err = run(sql)
		a.NoError(err, sql)
		a.Len(result, 1, sql)
		for _, r := range result {
			_, perr := parser.Parse(r.Statement)
			a.NoError(perr, sql)
		}
	}

	result, err = run("UPDATE test AS T SET t.c = 2 WHERE T.a = 1")
	a.NoError(err)
	a.Len(result, 1)
	a.Equal("test", result[0].SourceTable)

	result, err = run("UPDATE test AS T JOIN test2 AS Y ON T.a = Y.a SET t.c = 2, y.c = 3")
	a.NoError(err)
	a.Len(result, 2)
	a.ElementsMatch([]string{"test", "test2"}, []string{result[0].SourceTable, result[1].SourceTable})

	result, err = run("DELETE t FROM test AS T WHERE T.a = 1")
	a.NoError(err)
	a.Len(result, 1)
	a.Equal("test", result[0].SourceTable)

	result, err = run("DELETE T FROM test AS t WHERE t.a = 1")
	a.NoError(err)
	a.Len(result, 1)
	a.Equal("test", result[0].SourceTable)

	mixedCaseDB := "UPDATE db.test SET c1 = 1 WHERE id = 1;\n" +
		"UPDATE DB.test SET c1 = 1 WHERE id = 2;\n" +
		"UPDATE db.test SET c1 = 1 WHERE id = 3;\n" +
		"UPDATE DB.test SET c1 = 1 WHERE id = 4;\n" +
		"UPDATE db.test SET c1 = 1 WHERE id = 5;\n" +
		"UPDATE DB.test SET c1 = 1 WHERE id = 6;"
	result, err = run(mixedCaseDB)
	a.NoError(err)
	a.Len(result, 1)
}

func TestBackupHonorsCaseSensitivity(t *testing.T) {
	a := require.New(t)
	catalog := corpusBackupCatalog()
	run := func(caseSensitive bool, sql string) ([]analysis.BackupStatement, error) {
		return analysis.TransformDMLToSelect(sql, analysis.BackupOptions{
			DefaultDatabase: "db",
			BackupDatabase:  "backupDB",
			TablePrefix:     "_rollback",
			CaseSensitive:   caseSensitive,
			Catalog:         catalog,
		})
	}

	_, err := run(false, "UPDATE DB.test SET c = 1 WHERE c = 2")
	a.NoError(err)
	_, err = run(true, "UPDATE DB.test SET c = 1 WHERE c = 2")
	a.Error(err)

	sixMixedCase := "DELETE FROM test WHERE a = 1;\n" +
		"DELETE FROM test WHERE a = 2;\n" +
		"DELETE FROM test WHERE a = 3;\n" +
		"DELETE FROM test WHERE a = 4;\n" +
		"DELETE FROM test WHERE a = 5;\n" +
		"DELETE FROM Test WHERE a = 6;"
	_, err = run(false, sixMixedCase)
	a.NoError(err)
	_, err = run(true, sixMixedCase)
	a.Error(err)

	_, err = run(false, "UPDATE test AS T SET t.c = 2 WHERE T.a = 1")
	a.NoError(err)
	_, err = run(true, "UPDATE test AS T SET t.c = 2 WHERE T.a = 1")
	a.Error(err)
}

func endWithSemicolon(sql string, byteEnd int) int {
	if byteEnd < len(sql) && sql[byteEnd] == ';' {
		return byteEnd + 1
	}
	return byteEnd
}
