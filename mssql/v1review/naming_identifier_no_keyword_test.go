package v1review

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mssql/parser"
	"github.com/bytebase/omni/v1review"
)

func TestNamingIdentifierNoKeyword(t *testing.T) {
	// identFinding is a finding as the tests compare it: its code, name,
	// table, pattern and search text, and the text its range anchors.
	type identFinding struct {
		Code    v1review.Code
		Name    string
		Table   string
		Pattern string
		Search  string
		Text    string
	}
	// keywords returns one finding per search text, each anchoring text.
	keywords := func(text string, searches ...string) []identFinding {
		var want []identFinding
		for _, search := range searches {
			want = append(want, identFinding{Code: v1review.NameIsKeyword, Name: strings.ToLower(search), Search: search, Text: text})
		}
		return want
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []identFinding
	}{
		{name: "no keyword", sql: "CREATE TABLE Person (Id INT, CONSTRAINT pk PRIMARY KEY (Id));\nCREATE INDEX ix ON Person(Id);"},
		{
			name: "table then columns, name lower-cased and search in source case",
			sql:  "CREATE TABLE [Fillfactor](Id INT, [Exec] INT);",
			want: []identFinding{
				{Code: v1review.NameIsKeyword, Name: "fillfactor", Search: "Fillfactor", Text: "CREATE TABLE [Fillfactor](Id INT, [Exec] INT)"},
				{Code: v1review.NameIsKeyword, Name: "exec", Search: "Exec", Text: "CREATE TABLE [Fillfactor](Id INT, [Exec] INT)"},
			},
		},
		{
			name: "create index",
			sql:  "CREATE UNIQUE CLUSTERED INDEX [Add] ON [Fillfactor]([Exec]);",
			want: keywords("CREATE UNIQUE CLUSTERED INDEX [Add] ON [Fillfactor]([Exec])", "Add"),
		},
		{
			name: "create schema and create database",
			sql:  "CREATE SCHEMA [Select];\nCREATE DATABASE [FROM];",
			want: append(keywords("CREATE SCHEMA [Select]", "Select"), keywords("CREATE DATABASE [FROM]", "FROM")...),
		},
		{name: "schema authorization not judged", sql: "CREATE SCHEMA AUTHORIZATION [select];\nCREATE SCHEMA s AUTHORIZATION [from];"},
		{
			name: "table and column constraints",
			sql:  "CREATE TABLE t (a INT CONSTRAINT [check] CHECK (a > 0) CONSTRAINT [default] DEFAULT 1, CONSTRAINT [key] FOREIGN KEY (a) REFERENCES [user]([order]));",
			want: keywords("CREATE TABLE t (a INT CONSTRAINT [check] CHECK (a > 0) CONSTRAINT [default] DEFAULT 1, CONSTRAINT [key] FOREIGN KEY (a) REFERENCES [user]([order]))", "check", "default", "key"),
		},
		{
			name: "computed column",
			sql:  "CREATE TABLE t (id INT, [as] AS (id + 1));",
			want: keywords("CREATE TABLE t (id INT, [as] AS (id + 1))", "as"),
		},
		{
			name: "qualifier not judged",
			sql:  "CREATE TABLE [select].[from].[where] (id INT);\nCREATE INDEX ix ON [select].[from](id);",
			want: keywords("CREATE TABLE [select].[from].[where] (id INT)", "where"),
		},
		{
			name: "double-quoted and unquoted names",
			sql:  `CREATE TABLE t ("Select" INT, absolute INT, Action INT, ZONE INT);`,
			want: keywords(`CREATE TABLE t ("Select" INT, absolute INT, Action INT, ZONE INT)`, "Select", "absolute", "Action", "ZONE"),
		},
		{
			name: "same name twice is two findings",
			sql:  "CREATE TABLE [KEY] ([KEY] INT);",
			want: keywords("CREATE TABLE [KEY] ([KEY] INT)", "KEY", "KEY"),
		},
		{
			name: "folding is Unicode's",
			sql:  "CREATE TABLE [KEY] ([İN] INT);",
			want: []identFinding{
				{Code: v1review.NameIsKeyword, Name: "key", Search: "KEY", Text: "CREATE TABLE [KEY] ([İN] INT)"},
				{Code: v1review.NameIsKeyword, Name: "in", Search: "İN", Text: "CREATE TABLE [KEY] ([İN] INT)"},
			},
		},
		{
			name: "near misses",
			sql:  "CREATE TABLE [selects] ([select ] INT, [ select] INT, [se]]lect] INT, [nvarchar] INT, [within  group] INT);\nCREATE TABLE #select ([#from] INT);",
		},
		{
			name: "entries that are not one word",
			sql:  "CREATE TABLE t ([END-EXEC] INT, [within group] INT);",
			want: []identFinding{
				{Code: v1review.NameIsKeyword, Name: "end-exec", Search: "END-EXEC", Text: "CREATE TABLE t ([END-EXEC] INT, [within group] INT)"},
				{Code: v1review.NameIsKeyword, Name: "within group", Search: "within group", Text: "CREATE TABLE t ([END-EXEC] INT, [within group] INT)"},
			},
		},
		{
			name: "alter table add, alter column and drop constraint",
			sql:  "ALTER TABLE t ADD [select] INT, CONSTRAINT [from] UNIQUE ([select]);\nALTER TABLE t ALTER COLUMN [where] INT;\nALTER TABLE t DROP CONSTRAINT [order];",
			want: append(append(
				keywords("ALTER TABLE t ADD [select] INT, CONSTRAINT [from] UNIQUE ([select])", "select", "from"),
				keywords("ALTER TABLE t ALTER COLUMN [where] INT", "where")...),
				keywords("ALTER TABLE t DROP CONSTRAINT [order]", "order")...),
		},
		{name: "alter table name and drop column not judged", sql: "ALTER TABLE [select] ADD c INT;\nALTER TABLE t DROP COLUMN [from];"},
		{
			name: "routine body anchors the routine",
			sql:  "CREATE PROCEDURE [select] AS BEGIN CREATE TABLE [from] ([where] INT); END",
			want: keywords("CREATE PROCEDURE [select] AS BEGIN CREATE TABLE [from] ([where] INT); END", "from", "where"),
		},
		{
			name: "trigger body",
			sql:  "CREATE TRIGGER [select] ON t AFTER INSERT AS CREATE TABLE [from] (id INT);",
			want: keywords("CREATE TRIGGER [select] ON t AFTER INSERT AS CREATE TABLE [from] (id INT)", "from"),
		},
		{
			name: "if block",
			sql:  "IF 1 = 1 BEGIN CREATE TABLE [select] (id INT); CREATE INDEX [from] ON t(a); END",
			want: keywords("IF 1 = 1 BEGIN CREATE TABLE [select] (id INT); CREATE INDEX [from] ON t(a); END", "select", "from"),
		},
		{
			name: "try block and while",
			sql:  "BEGIN TRY CREATE TABLE [select] (id INT); END TRY BEGIN CATCH CREATE DATABASE [from]; END CATCH\nWHILE 1 = 0 CREATE SCHEMA [where];",
			want: append(
				keywords("BEGIN TRY CREATE TABLE [select] (id INT); END TRY BEGIN CATCH CREATE DATABASE [from]; END CATCH", "select", "from"),
				keywords("WHILE 1 = 0 CREATE SCHEMA [where]", "where")...),
		},
		{
			name: "create schema elements",
			sql:  "CREATE SCHEMA [select] AUTHORIZATION [from] CREATE TABLE [where] ([order] INT)",
			want: keywords("CREATE SCHEMA [select] AUTHORIZATION [from] CREATE TABLE [where] ([order] INT)", "select", "where", "order"),
		},
		{
			name: "table variable, table type and returned table",
			sql:  "DECLARE @select TABLE ([from] INT);\nCREATE TYPE [select] AS TABLE ([where] INT);\nCREATE FUNCTION [select]() RETURNS @r TABLE ([order] INT) AS BEGIN RETURN; END",
			want: append(append(
				keywords("DECLARE @select TABLE ([from] INT)", "from"),
				keywords("CREATE TYPE [select] AS TABLE ([where] INT)", "where")...),
				keywords("CREATE FUNCTION [select]() RETURNS @r TABLE ([order] INT) AS BEGIN RETURN; END", "order")...),
		},
		{
			name: "external table judges its columns only",
			sql:  "CREATE EXTERNAL TABLE [select] ([from] INT) WITH (LOCATION = '/x', DATA_SOURCE = ds, FILE_FORMAT = ff);",
			want: keywords("CREATE EXTERNAL TABLE [select] ([from] INT) WITH (LOCATION = '/x', DATA_SOURCE = ds, FILE_FORMAT = ff)", "from"),
		},
		{name: "inline index not judged", sql: "CREATE TABLE t (id INT, INDEX [select] (id));"},
		{name: "xml, spatial and full-text indexes not judged", sql: "CREATE PRIMARY XML INDEX [select] ON t(x);\nCREATE SPATIAL INDEX [from] ON t(g);\nCREATE FULLTEXT INDEX ON t([where]) KEY INDEX [order];"},
		{name: "create table as select and clone not judged", sql: "CREATE TABLE [select] WITH (DISTRIBUTION = ROUND_ROBIN) AS SELECT 1 AS id;\nCREATE TABLE [from] AS CLONE OF dbo.t;"},
		{name: "select into and aliases not judged", sql: "SELECT 1 AS [select] INTO [from];\nSELECT [select].[from] AS [where] FROM [order] AS [group];"},
		{name: "sp_rename not judged", sql: "EXEC sp_rename 'dbo.t', 'select';"},
		{name: "openjson and result sets columns not judged", sql: "SELECT * FROM OPENJSON(@j) WITH ([select] INT '$.a');\nEXEC p WITH RESULT SETS (([from] INT));"},
		{
			name: "other objects not judged",
			sql:  "CREATE VIEW [select] AS SELECT 1 AS a;\nCREATE SEQUENCE [from] START WITH 1;\nCREATE SYNONYM [where] FOR t;\nCREATE TYPE [order] FROM INT;\nCREATE STATISTICS [group] ON t(a);\nCREATE USER [by];\nALTER DATABASE [select] SET RECOVERY FULL;\nDROP INDEX [from] ON t;\nALTER INDEX [where] ON t REBUILD;",
		},
		{
			name: "later statement keeps its own range",
			sql:  "CREATE TABLE t (id INT);\n\n-- select\nCREATE TABLE u (\n  id INT,\n  [Select] INT\n);",
			want: keywords("CREATE TABLE u (\n  id INT,\n  [Select] INT\n)", "Select"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var got []identFinding
			for _, f := range NamingIdentifierNoKeyword(list.Items) {
				got = append(got, identFinding{f.Code, f.Name, f.Table, f.Pattern, f.Search, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
