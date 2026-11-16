package catalog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
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

func TestMigrationKeepsSequencePositionWhenOnlyIncrementChanges(t *testing.T) {
	from, err := LoadSQL(`CREATE SEQUENCE [dbo].[s] AS INT START WITH 1 INCREMENT BY 1;`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE SEQUENCE [dbo].[s] AS INT START WITH 1 INCREMENT BY 2;`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "INCREMENT BY 2")
	require.NotContains(t, sql, "RESTART WITH")
}

func TestMigrationRecreatesChangedComputedColumn(t *testing.T) {
	from, err := LoadSQL(`
CREATE TABLE [dbo].[t] (
	[id] int NOT NULL,
	[value] int NULL,
	[doubled] AS ([value] * 2) PERSISTED,
	CONSTRAINT [PK_t] PRIMARY KEY ([id])
);
CREATE INDEX [IX_t_doubled] ON [dbo].[t] ([doubled]);
`)
	require.NoError(t, err)
	to, err := LoadSQL(`
CREATE TABLE [dbo].[t] (
	[id] int NOT NULL,
	[value] int NULL,
	[doubled] AS ([value] * 3) PERSISTED,
	CONSTRAINT [PK_t] PRIMARY KEY ([id])
);
CREATE INDEX [IX_t_doubled] ON [dbo].[t] ([doubled]);
`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "DROP INDEX [IX_t_doubled]")
	require.Contains(t, sql, "DROP COLUMN [doubled]")
	require.Contains(t, sql, "ADD [doubled] AS ([value] * 3) PERSISTED")
	require.Contains(t, sql, "CREATE INDEX [IX_t_doubled]")
	require.NotContains(t, sql, "ALTER COLUMN [doubled] AS")
}

func TestDefinitionRendersTriggers(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{{
			Name:    "t",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: false}},
			Triggers: []*metadata.TriggerMetadata{{
				Name: "trg", Event: "INSERT", Timing: "AFTER",
				Body: "CREATE TRIGGER [dbo].[trg] ON [dbo].[t] AFTER INSERT AS SELECT 1",
			}},
		}},
	}}}
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "CREATE TRIGGER [dbo].[trg] ON [dbo].[t] AFTER INSERT AS SELECT 1")
	require.Contains(t, def, "GO")
}

func TestForeignKeyWithoutReferencedColumnsRendersValidDDL(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE [dbo].[parent] ([id] int NOT NULL, CONSTRAINT [PK_parent] PRIMARY KEY ([id]));
CREATE TABLE [dbo].[child] (
	[id] int NOT NULL,
	[parent_id] int NOT NULL,
	CONSTRAINT [PK_child] PRIMARY KEY ([id]),
	CONSTRAINT [FK_child_parent] FOREIGN KEY ([parent_id]) REFERENCES [dbo].[parent]
);
`)
	require.NoError(t, err)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "REFERENCES [dbo].[parent]")
	require.NotContains(t, def, "REFERENCES [dbo].[parent] ()")
}

func TestDefinitionEscapesClosingBrackets(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE [dbo].[a]]b] ([c]]d] int NOT NULL, CONSTRAINT [PK_a]]b] PRIMARY KEY ([c]]d]));`)
	require.NoError(t, err)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "CREATE TABLE [dbo].[a]]b]")
	require.Contains(t, def, "[c]]d] int")
	require.Contains(t, def, "CONSTRAINT [PK_a]]b]")
}

func TestUnnamedTableConstraintsGetNames(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, [code] int NOT NULL, PRIMARY KEY ([id]), UNIQUE ([code]));`)
	require.NoError(t, err)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.NotContains(t, def, "CONSTRAINT []")
	require.Contains(t, def, "CONSTRAINT [PK_t_1]")
	require.Contains(t, def, "CONSTRAINT [UQ_t_2]")
}

func TestTableConstraintKeepsDescendingDirections(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, [code] int NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id] DESC), CONSTRAINT [UQ_t] UNIQUE ([code] ASC));`)
	require.NoError(t, err)
	indexes := c.GetTable("dbo", "t").Indexes
	require.True(t, indexes[0].Descending[0])
	require.False(t, indexes[1].Descending[0])
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "PRIMARY KEY ([id] DESC)")
}

func TestDefinitionAndMigrationEncloseCheckExpression(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, CONSTRAINT [CK_t] CHECK ([id] > 0));`)
	require.NoError(t, err)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "CONSTRAINT [CK_t] CHECK ([id] > 0)")
	to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, CONSTRAINT [CK_t] CHECK ([id] > 0), CONSTRAINT [CK_t2] CHECK ([id] < 100));`)
	require.NoError(t, err)
	sql := GenerateMigration(c, to, Diff(c, to)).SQL()
	require.Contains(t, sql, "CHECK ([id] < 100)")
	require.NotContains(t, sql, "CHECK (([id] < 100))")
}

func TestMigrationEmitsTriggerChanges(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	to := fromWithTrigger(t, "CREATE TRIGGER [dbo].[trg] ON [dbo].[t] AFTER INSERT AS SELECT 1")
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "CREATE TRIGGER [dbo].[trg] ON [dbo].[t] AFTER INSERT AS SELECT 1")
	back := GenerateMigration(to, from, Diff(to, from)).SQL()
	require.Contains(t, back, "DROP TRIGGER [dbo].[trg]")
}

func fromWithTrigger(t *testing.T, body string) *Catalog {
	t.Helper()
	c, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	table := c.GetTable("dbo", "t")
	table.Triggers = []*metadata.TriggerMetadata{{Name: "trg", Event: "INSERT", Timing: "AFTER", Body: body}}
	return c
}

func TestDefinitionCreatesViewReferencedFunctionFirst(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([id] int NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));
CREATE VIEW [dbo].[v] AS SELECT [dbo].[f]([id]) AS value FROM [dbo].[t];
CREATE FUNCTION [dbo].[f] (@id int) RETURNS int AS BEGIN RETURN @id END;
`)
	require.NoError(t, err)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Less(t, strings.Index(def, "CREATE FUNCTION [dbo].[f]"), strings.Index(def, "CREATE VIEW [dbo].[v]"))
}

func TestMigrationRebuildsKeyAroundColumnTypeChange(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] bigint NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "DROP CONSTRAINT [PK_t]")
	require.Contains(t, sql, "ALTER COLUMN [id] bigint NOT NULL")
	require.Contains(t, sql, "ADD CONSTRAINT [PK_t] PRIMARY KEY")
	require.Less(t, strings.Index(sql, "DROP CONSTRAINT [PK_t]"), strings.Index(sql, "ALTER COLUMN [id] bigint"))
	require.Less(t, strings.Index(sql, "ALTER COLUMN [id] bigint"), strings.Index(sql, "ADD CONSTRAINT [PK_t]"))
}
