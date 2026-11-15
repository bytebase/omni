package analysis

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql"
)

// parseScript parses each statement of script on its own, with the line it
// starts on, as Bytebase's splitter hands statements over.
func parseScript(t *testing.T, script string) []ParsedStatement {
	t.Helper()
	stmts, err := mssql.Parse(script)
	if err != nil {
		t.Fatalf("parse %q: %v", script, err)
	}
	var result []ParsedStatement
	for _, stmt := range stmts {
		nodes, err := mssql.Parse(stmt.Text)
		if err != nil || len(nodes) != 1 {
			t.Fatalf("parse %q: %v", stmt.Text, err)
		}
		startLine := 1 + strings.Count(script[:stmt.ByteStart], "\n")
		result = append(result, ParsedStatement{Text: stmt.Text, Node: nodes[0].AST, StartLine: startLine})
	}
	return result
}

func TestTransformDMLToSelect(t *testing.T) {
	script := "UPDATE t SET a = 1 WHERE b = 1;\nDELETE x FROM s.u AS x WHERE x.c = 2;\nUPDATE t SET a = 2 WHERE b = 2;\nDELETE FROM #tmp;"
	got, err := TransformDMLToSelect(parseScript(t, script), "db", "bbarchive", "pre", false)
	if err != nil {
		t.Fatal(err)
	}
	want := []BackupStatement{
		{
			Statement:       "SELECT * INTO [bbarchive].[dbo].[pre_t_db] FROM (\n  SELECT [db].[dbo].[t].* FROM t WHERE b = 1\n  UNION\n  SELECT [db].[dbo].[t].* FROM t WHERE b = 2) AS backup_table;",
			SourceSchema:    "dbo",
			SourceTableName: "t",
			TargetTableName: "pre_t_db",
			StartPosition:   Position{Line: 1, Column: 0},
			EndPosition:     Position{Line: 3, Column: 30},
		},
		{
			Statement:       "SELECT * INTO [bbarchive].[dbo].[pre_u_db] FROM (\n  SELECT [x].* FROM s.u AS x WHERE x.c = 2) AS backup_table;",
			SourceSchema:    "s",
			SourceTableName: "u",
			TargetTableName: "pre_u_db",
			StartPosition:   Position{Line: 2, Column: 0},
			EndPosition:     Position{Line: 2, Column: 36},
		},
	}
	if len(got) != len(want) {
		t.Fatalf("got %+v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("statement %d:\ngot  %+v\nwant %+v", i, got[i], want[i])
		}
	}
}

func TestTransformDMLToSelectRefuses(t *testing.T) {
	for _, script := range []string{
		"UPDATE t SET a = 1; DELETE FROM t",
		"WITH c AS (SELECT a FROM t) UPDATE c SET a = 1",
		"UPDATE t SET a = 1 WHERE CURRENT OF cur",
	} {
		if _, err := TransformDMLToSelect(parseScript(t, script), "db", "bbarchive", "pre", false); err == nil {
			t.Errorf("%q: got no error", script)
		}
	}
}

func TestRestoreSQL(t *testing.T) {
	table := &metadata.TableMetadata{
		Name: "t",
		Columns: []*metadata.ColumnMetadata{
			{Name: "id", IsIdentity: true},
			{Name: "a"},
			{Name: "g", Generation: &metadata.GenerationMetadata{Expression: "a + 1"}},
		},
		Indexes: []*metadata.IndexMetadata{{Name: "uk", Unique: true, Expressions: []string{"a"}}},
	}
	opts := RestoreOptions{
		OriginalDatabase: "db",
		OriginalSchema:   "dbo",
		OriginalTable:    "t",
		BackupDatabase:   "bbarchive",
		BackupTable:      "pre_t_db",
		Table:            table,
		PrimaryKey:       &metadata.IndexMetadata{Primary: true, Expressions: []string{"id"}},
		OriginalSQL:      strings.Repeat("x", 1001),
	}
	dmls := BackedUpDMLs(parseScript(t, "UPDATE t SET a = 1 WHERE id = 2;\nDELETE FROM t WHERE id = 3;"))
	if len(dmls) != 2 || dmls[0].SQL != "UPDATE t SET a = 1 WHERE id = 2" || dmls[1].StartPosition != (Position{Line: 2, Column: 0}) {
		t.Fatalf("got %+v", dmls)
	}

	got, err := RestoreSQL(dmls[0].Node, opts)
	if err != nil {
		t.Fatal(err)
	}
	want := "/*\nOriginal SQL:\n" + strings.Repeat("x", 1000) + "...\n*/\n" +
		"SET IDENTITY_INSERT [db].[dbo].[t] ON;\n" +
		"MERGE INTO [db].[dbo].[t] AS t\nUSING [bbarchive].[dbo].[pre_t_db] AS b\n  ON t.[id] = b.[id]\n" +
		"WHEN MATCHED THEN\n  UPDATE SET t.[a] = b.[a]\nWHEN NOT MATCHED THEN\n INSERT ([id], [a]) VALUES (b.[id], b.[a]);\n\n" +
		"SET IDENTITY_INSERT [db].[dbo].[t] OFF;\nEXEC('DBCC CHECKIDENT (''[db].[dbo].[t]'', RESEED)');"
	if got != want {
		t.Errorf("update:\ngot  %q\nwant %q", got, want)
	}

	opts.Table = &metadata.TableMetadata{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "a"}}}
	opts.OriginalSQL = "DELETE FROM t"
	got, err = RestoreSQL(dmls[1].Node, opts)
	if err != nil {
		t.Fatal(err)
	}
	if want := "/*\nOriginal SQL:\nDELETE FROM t\n*/\nINSERT INTO [db].[dbo].[t] ([a]) SELECT [a] FROM [bbarchive].[dbo].[pre_t_db];"; got != want {
		t.Errorf("delete:\ngot  %q\nwant %q", got, want)
	}

	// SQL Server nests block comments: the quoted statement must neither close
	// the comment nor open another.
	opts.OriginalSQL = "DELETE FROM t WHERE a IN ('*/ DROP TABLE x; /*', '/*/')"
	got, err = RestoreSQL(dmls[1].Node, opts)
	if err != nil {
		t.Fatal(err)
	}
	if want := "/*\nOriginal SQL:\nDELETE FROM t WHERE a IN ('* / DROP TABLE x; / *', '/ * /')\n*/\nINSERT INTO [db].[dbo].[t] ([a]) SELECT [a] FROM [bbarchive].[dbo].[pre_t_db];"; got != want {
		t.Errorf("comment delimiters:\ngot  %q\nwant %q", got, want)
	}

	opts.PrimaryKey = nil
	if _, err := RestoreSQL(parseScript(t, "UPDATE t SET a = 1")[0].Node, opts); err == nil {
		t.Error("update of every unique key: got no error")
	}
}
