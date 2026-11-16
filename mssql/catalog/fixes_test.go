package catalog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefinitionRendersSequenceOptions(t *testing.T) {
	c, err := LoadSQL(`CREATE SEQUENCE [dbo].[s] AS BIGINT START WITH 100 INCREMENT BY 2 MINVALUE 1 MAXVALUE 999 CYCLE CACHE 10;`)
	require.NoError(t, err)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	for _, want := range []string{
		"CREATE SEQUENCE [dbo].[s]",
		"    AS BIGINT",
		"    START WITH 100",
		"    INCREMENT BY 2",
		"    MINVALUE 1",
		"    MAXVALUE 999",
		"    CYCLE",
		"    CACHE 10",
	} {
		require.Contains(t, def, want)
	}
}

func TestDefinitionDefersCyclicForeignKeys(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE [dbo].[a] (
	[id] int NOT NULL,
	[b_id] int NULL,
	CONSTRAINT [PK_a] PRIMARY KEY ([id]),
	CONSTRAINT [FK_a_b] FOREIGN KEY ([b_id]) REFERENCES [dbo].[b] ([id])
);
CREATE TABLE [dbo].[b] (
	[id] int NOT NULL,
	[a_id] int NULL,
	CONSTRAINT [PK_b] PRIMARY KEY ([id]),
	CONSTRAINT [FK_b_a] FOREIGN KEY ([a_id]) REFERENCES [dbo].[a] ([id])
);
`)
	require.NoError(t, err)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "ALTER TABLE [dbo].[a] ADD CONSTRAINT [FK_a_b] FOREIGN KEY")
	require.Contains(t, def, "ALTER TABLE [dbo].[b] ADD CONSTRAINT [FK_b_a] FOREIGN KEY")
}

func TestLoadSQLKeepsComputedColumnsAndFilteredIndexes(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE [dbo].[t] (
	[id] int NOT NULL,
	[value] int NULL,
	[doubled] AS ([value] * 2) PERSISTED,
	CONSTRAINT [PK_t] PRIMARY KEY ([id])
);
CREATE UNIQUE NONCLUSTERED INDEX [UX_t_value] ON [dbo].[t] ([value]) WHERE [value] IS NOT NULL;
`)
	require.NoError(t, err)
	table := c.GetTable("dbo", "t")
	require.NotNil(t, table)
	require.Equal(t, "AS ([value] * 2) PERSISTED", table.Columns[2].Type)
	require.Len(t, table.Indexes, 2)
	require.Equal(t, "WHERE [value] IS NOT NULL", table.Indexes[1].Definition)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "[doubled] AS ([value] * 2) PERSISTED")
	require.Contains(t, def, "WHERE [value] IS NOT NULL")
}

func TestMigrationDropsInboundForeignKeyBeforeReplacingKey(t *testing.T) {
	from, err := LoadSQL(`
CREATE TABLE [dbo].[p] ([id] int NOT NULL, CONSTRAINT [PK_p] PRIMARY KEY ([id]));
CREATE TABLE [dbo].[c] (
	[id] int NOT NULL,
	[p_id] int NOT NULL,
	CONSTRAINT [PK_c] PRIMARY KEY ([id]),
	CONSTRAINT [FK_c_p] FOREIGN KEY ([p_id]) REFERENCES [dbo].[p] ([id])
);
`)
	require.NoError(t, err)
	to, err := LoadSQL(`
CREATE TABLE [dbo].[p] ([id] int NOT NULL, CONSTRAINT [PK_p_new] PRIMARY KEY ([id]));
CREATE TABLE [dbo].[c] (
	[id] int NOT NULL,
	[p_id] int NOT NULL,
	CONSTRAINT [PK_c] PRIMARY KEY ([id]),
	CONSTRAINT [FK_c_p] FOREIGN KEY ([p_id]) REFERENCES [dbo].[p] ([id])
);
`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "DROP CONSTRAINT [FK_c_p]")
	require.Contains(t, sql, "DROP CONSTRAINT [PK_p]")
	require.Contains(t, sql, "ADD CONSTRAINT [FK_c_p]")
	require.Contains(t, sql, "ADD CONSTRAINT [PK_p_new]")
	require.Less(t, strings.Index(sql, "DROP CONSTRAINT [FK_c_p]"), strings.Index(sql, "DROP CONSTRAINT [PK_p]"))
	require.Less(t, strings.Index(sql, "ADD CONSTRAINT [PK_p_new]"), strings.Index(sql, "ADD CONSTRAINT [FK_c_p]"))
}

func TestMigrationDropsAndAltersSequences(t *testing.T) {
	from, err := LoadSQL(`
CREATE SEQUENCE [dbo].[s] AS INT START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE [dbo].[gone] AS INT;
`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE SEQUENCE [dbo].[s] AS INT START WITH 10 INCREMENT BY 2;`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "DROP SEQUENCE [dbo].[gone]")
	require.Contains(t, sql, "ALTER SEQUENCE [dbo].[s] RESTART WITH 10 INCREMENT BY 2")
}

func TestExpressionComparisonKeepsLiterals(t *testing.T) {
	require.True(t, expressionsEqual("1 = 1", "1  =  1"))
	require.True(t, expressionsEqual("'x' = 'x'", "'x' = 'x'"))
	require.False(t, expressionsEqual("'a  b'", "'a b'"))
	require.False(t, expressionsEqual("'pending'", "'PENDING'"))
	require.False(t, checkExpressionsEqual("'pending'", "'PENDING'"))
}

func TestMigrationDropsSchemaObjectsBeforeSchema(t *testing.T) {
	from, err := LoadSQL(`
CREATE SCHEMA [gone];
CREATE TABLE [gone].[t] ([id] int NOT NULL, CONSTRAINT [PK_gone_t] PRIMARY KEY ([id]));
CREATE VIEW [gone].[v] AS SELECT [id] FROM [gone].[t];
CREATE PROCEDURE [gone].[p] AS SELECT 1;
CREATE FUNCTION [gone].[f]() RETURNS int AS BEGIN RETURN 1 END;
CREATE SEQUENCE [gone].[s] AS INT;
`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE SCHEMA [kept];`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	for _, want := range []string{
		"DROP VIEW [gone].[v]",
		"DROP PROCEDURE [gone].[p]",
		"DROP FUNCTION [gone].[f]",
		"DROP SEQUENCE [gone].[s]",
		"DROP TABLE [gone].[t]",
		"DROP SCHEMA [gone]",
	} {
		require.Contains(t, sql, want)
	}
	require.Less(t, strings.Index(sql, "DROP VIEW [gone].[v]"), strings.Index(sql, "DROP TABLE [gone].[t]"))
	require.Less(t, strings.Index(sql, "DROP TABLE [gone].[t]"), strings.Index(sql, "DROP SCHEMA [gone]"))
}
