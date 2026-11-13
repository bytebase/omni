package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mssql/parser"
	"github.com/bytebase/omni/v1review"
)

func TestNamingTableNoKeyword(t *testing.T) {
	// keywordFinding is a finding as the tests compare it: its code, name,
	// table and pattern, and the text its range anchors.
	type keywordFinding struct {
		Code    v1review.Code
		Name    string
		Table   string
		Pattern string
		Text    string
	}
	keyword := func(name, text string) keywordFinding {
		return keywordFinding{Code: v1review.NameIsKeyword, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []keywordFinding
	}{
		{name: "no keyword", sql: "CREATE TABLE Person (id INT);"},
		{
			name: "keyword anchors the statement",
			sql:  "CREATE TABLE [fillfactor] (id INT);",
			want: []keywordFinding{keyword("fillfactor", "CREATE TABLE [fillfactor] (id INT)")},
		},
		{
			name: "name keeps its case",
			sql:  "CREATE TABLE [Select] (id INT);\nCREATE TABLE [FROM] (id INT);",
			want: []keywordFinding{
				keyword("Select", "CREATE TABLE [Select] (id INT)"),
				keyword("FROM", "CREATE TABLE [FROM] (id INT)"),
			},
		},
		{
			name: "double-quoted name",
			sql:  `CREATE TABLE "Where" (id INT);`,
			want: []keywordFinding{keyword("Where", `CREATE TABLE "Where" (id INT)`)},
		},
		{
			name: "unquoted word the parser does not reserve",
			sql:  "CREATE TABLE absolute (id INT);",
			want: []keywordFinding{keyword("absolute", "CREATE TABLE absolute (id INT)")},
		},
		{
			name: "qualifier not judged",
			sql:  "CREATE TABLE [select].[from].Person (id INT);\nCREATE TABLE MyDB.MySchema.[where] (id INT);\nCREATE TABLE db..[order] (id INT);",
			want: []keywordFinding{
				keyword("where", "CREATE TABLE MyDB.MySchema.[where] (id INT)"),
				keyword("order", "CREATE TABLE db..[order] (id INT)"),
			},
		},
		{
			name: "entries that are not one word",
			sql:  "CREATE TABLE [END-EXEC] (id INT);\nCREATE TABLE [within group] (id INT);",
			want: []keywordFinding{
				keyword("END-EXEC", "CREATE TABLE [END-EXEC] (id INT)"),
				keyword("within group", "CREATE TABLE [within group] (id INT)"),
			},
		},
		{
			name: "upper-casing is Unicode's",
			sql:  "CREATE TABLE [ınsert] (id INT);\nCREATE TABLE ſelect (id INT);",
			want: []keywordFinding{
				keyword("ınsert", "CREATE TABLE [ınsert] (id INT)"),
				keyword("ſelect", "CREATE TABLE ſelect (id INT)"),
			},
		},
		{
			name: "near misses",
			sql:  "CREATE TABLE [selects] (id INT);\nCREATE TABLE [select ] (id INT);\nCREATE TABLE [ select] (id INT);\nCREATE TABLE [within  group] (id INT);\nCREATE TABLE [select]]] (id INT);\nCREATE TABLE [nvarchar] (id INT);",
		},
		{name: "temporary table judged with its prefix", sql: "CREATE TABLE #select (id INT);\nCREATE TABLE ##select (id INT);"},
		{
			name: "graph and file tables",
			sql:  "CREATE TABLE [key] (id INT) AS EDGE;\nCREATE TABLE [index] AS FILETABLE;",
			want: []keywordFinding{
				keyword("key", "CREATE TABLE [key] (id INT) AS EDGE"),
				keyword("index", "CREATE TABLE [index] AS FILETABLE"),
			},
		},
		{name: "create table as select not judged", sql: "CREATE TABLE [select] WITH (DISTRIBUTION = ROUND_ROBIN) AS SELECT 1 AS id;"},
		{name: "create table as clone not judged", sql: "CREATE TABLE [select] AS CLONE OF dbo.t;"},
		{name: "create external table not judged", sql: "CREATE EXTERNAL TABLE [select] (id INT) WITH (LOCATION = '/x', DATA_SOURCE = ds, FILE_FORMAT = ff);"},
		{name: "select into not judged", sql: "SELECT 1 AS id INTO [select];"},
		{name: "sp_rename not judged", sql: "EXEC sp_rename 'dbo.t', 'select';"},
		{name: "table variable not judged", sql: "DECLARE @select TABLE (id INT);"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p AS CREATE TABLE [select] (id INT);"},
		{name: "if block not judged", sql: "IF 1 = 1 BEGIN CREATE TABLE [select] (id INT); END"},
		{name: "begin block not judged", sql: "BEGIN CREATE TABLE [select] (id INT); END"},
		{name: "column, constraint and index names not judged", sql: "CREATE TABLE t ([select] INT, CONSTRAINT [from] UNIQUE ([select]), INDEX [where] ([select]));"},
		{
			name: "later statement on a later line",
			sql:  "CREATE TABLE t (id INT);\n\n-- comment\nCREATE TABLE\n  [user]\n(\n  id INT\n);",
			want: []keywordFinding{keyword("user", "CREATE TABLE\n  [user]\n(\n  id INT\n)")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var got []keywordFinding
			for _, f := range NamingTableNoKeyword(list.Items) {
				got = append(got, keywordFinding{f.Code, f.Name, f.Table, f.Pattern, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
