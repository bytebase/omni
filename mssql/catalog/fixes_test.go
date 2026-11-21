package catalog

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

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

func TestDefinitionRendersNonTableComments(t *testing.T) {
	c, err := LoadSQL(`
CREATE SEQUENCE [dbo].[s] AS INT;
GO
EXEC sp_addextendedproperty 'MS_Description', N'sequence note', 'SCHEMA', 'dbo', 'SEQUENCE', 's';
GO
CREATE VIEW [dbo].[v] AS SELECT 1 AS n;
GO
EXEC sp_addextendedproperty 'MS_Description', N'view note', 'SCHEMA', 'dbo', 'VIEW', 'v';
GO
CREATE FUNCTION [dbo].[f]() RETURNS int AS BEGIN RETURN 1 END;
GO
EXEC sp_addextendedproperty 'MS_Description', N'function note', 'SCHEMA', 'dbo', 'FUNCTION', 'f';
GO
CREATE PROCEDURE [dbo].[p] AS SELECT 1;
GO
EXEC sp_addextendedproperty 'MS_Description', N'procedure note', 'SCHEMA', 'dbo', 'PROCEDURE', 'p';
`)
	require.NoError(t, err)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	for _, want := range []string{
		"sp_addextendedproperty 'MS_Description', N'sequence note', 'SCHEMA', 'dbo', 'SEQUENCE', 's'",
		"sp_addextendedproperty 'MS_Description', N'view note', 'SCHEMA', 'dbo', 'VIEW', 'v'",
		"sp_addextendedproperty 'MS_Description', N'function note', 'SCHEMA', 'dbo', 'FUNCTION', 'f'",
		"sp_addextendedproperty 'MS_Description', N'procedure note', 'SCHEMA', 'dbo', 'PROCEDURE', 'p'",
	} {
		require.Contains(t, def, want)
	}
}

func TestDefinitionSkipsCommentsWithoutADefinition(t *testing.T) {
	placeholder := func(kind, name string) string {
		return "/* Definition of " + kind + " dbo." + name + " is encrypted. */"
	}
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Views: []*metadata.ViewMetadata{{
			Name:       "v",
			Definition: placeholder("view", "v"),
			Comment:    "view note",
		}},
		Functions: []*metadata.FunctionMetadata{{
			Name:       "f",
			Definition: placeholder("function", "f"),
			Comment:    "function note",
		}},
		Procedures: []*metadata.ProcedureMetadata{{
			Name:       "p",
			Definition: placeholder("procedure", "p"),
			Comment:    "procedure note",
		}},
	}}}
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.NotContains(t, def, "sp_addextendedproperty")
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

func TestDefinitionRendersSynonyms(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Synonyms: []*metadata.SynonymMetadata{{
			Name: "orders_alias", TargetServer: "srv", TargetDatabase: "db", TargetSchema: "sales", TargetName: "orders",
		}},
	}}}
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "CREATE SYNONYM [dbo].[orders_alias] FOR [srv].[db].[sales].[orders];")
}

func TestDefinitionNormalizesAlterModules(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{{
			Name:    "t",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: false}},
		}},
		Views:      []*metadata.ViewMetadata{{Name: "v", Definition: "ALTER VIEW [dbo].[v] AS SELECT 1 AS n"}},
		Functions:  []*metadata.FunctionMetadata{{Name: "f", Definition: "ALTER FUNCTION [dbo].[f]() RETURNS int AS BEGIN RETURN 1 END"}},
		Procedures: []*metadata.ProcedureMetadata{{Name: "p", Definition: "ALTER PROCEDURE [dbo].[p] AS SELECT 1"}},
	}}}
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "CREATE VIEW [dbo].[v]")
	require.Contains(t, def, "CREATE FUNCTION [dbo].[f]()")
	require.Contains(t, def, "CREATE PROCEDURE [dbo].[p]")
	require.NotContains(t, def, "ALTER VIEW")
	require.NotContains(t, def, "ALTER FUNCTION")
	require.NotContains(t, def, "ALTER PROCEDURE")
}

func TestTriggerComparisonKeepsLiterals(t *testing.T) {
	oldTrigger := &metadata.TriggerMetadata{Name: "trg", Event: "INSERT", Timing: "AFTER", Body: "SELECT 'pending'"}
	newTrigger := &metadata.TriggerMetadata{Name: "trg", Event: "INSERT", Timing: "AFTER", Body: "SELECT 'PENDING'"}
	require.False(t, triggersEqual(oldTrigger, newTrigger))
}

func TestMigrationDeduplicatesCompositeKeyRebuild(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE [dbo].[t] ([a] int NOT NULL, [b] int NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([a], [b]));`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([a] bigint NOT NULL, [b] bigint NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([a], [b]));`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Equal(t, 1, strings.Count(sql, "DROP CONSTRAINT [PK_t]"))
	require.Equal(t, 1, strings.Count(sql, "ADD CONSTRAINT [PK_t]"))
}

func TestLoadSQLExtractsTriggers(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([id] int NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));
GO
CREATE TRIGGER [dbo].[trg] ON [dbo].[t] AFTER INSERT AS SELECT 1;
`)
	require.NoError(t, err)
	table := c.GetTable("dbo", "t")
	require.Len(t, table.Triggers, 1)
	require.Equal(t, "trg", table.Triggers[0].Name)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "CREATE TRIGGER [dbo].[trg]")
}

func TestDefinitionHonorsSkipDump(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{{
			Name: "hidden", SkipDump: true,
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}},
		}},
		Views:     []*metadata.ViewMetadata{{Name: "hidden_view", SkipDump: true, Definition: "CREATE VIEW [dbo].[hidden_view] AS SELECT 1 AS n"}},
		Functions: []*metadata.FunctionMetadata{{Name: "hidden_fn", SkipDump: true, Definition: "CREATE FUNCTION [dbo].[hidden_fn]() RETURNS int AS BEGIN RETURN 1 END"}},
	}}}
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.NotContains(t, def, "hidden")
}

func TestIndexComparisonIncludesConstraintBacking(t *testing.T) {
	plain := &metadata.IndexMetadata{Name: "ux", Type: "NONCLUSTERED", Unique: true, Expressions: []string{"a"}, Descending: []bool{false}}
	constraint := &metadata.IndexMetadata{Name: "ux", Type: "NONCLUSTERED", Unique: true, IsConstraint: true, Expressions: []string{"a"}, Descending: []bool{false}}
	require.False(t, indexesEqual(plain, constraint))
}

func TestAlterColumnKeepsCollationOnTypeChange(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE [dbo].[t] ([v] varchar(20) COLLATE Latin1_General_CS_AS NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([v]));`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([v] varchar(40) COLLATE Latin1_General_CS_AS NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([v]));`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "ALTER COLUMN [v] varchar(40) COLLATE Latin1_General_CS_AS NOT NULL")
}

func TestDropFunctionFollowsSchemaBoundView(t *testing.T) {
	from, err := LoadSQL(`
CREATE SCHEMA [s];
CREATE TABLE [s].[t] ([id] int NOT NULL, CONSTRAINT [PK_s_t] PRIMARY KEY ([id]));
CREATE FUNCTION [s].[f]() RETURNS int AS BEGIN RETURN 1 END;
CREATE VIEW [s].[v] WITH SCHEMABINDING AS SELECT [s].[f]() AS n;
`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE SCHEMA [s];`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Less(t, strings.Index(sql, "DROP VIEW [s].[v]"), strings.Index(sql, "DROP FUNCTION [s].[f]"))
}

func TestIdentityChangeIsExplicitlyUnsupported(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int IDENTITY(1,1) NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int IDENTITY(1,10) NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "UNSUPPORTED: identity property change")
	require.NotContains(t, sql, "DROP COLUMN [id]")
}

func TestDefinitionAndMigrationRenderXMLIndexes(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{{
			Name:    "t",
			Columns: []*metadata.ColumnMetadata{{Name: "doc", Type: "xml", Nullable: true}},
			Indexes: []*metadata.IndexMetadata{
				{Name: "pxml", Type: "XML", Expressions: []string{"doc"}, Definition: "PRIMARY"},
				{Name: "sxml", Type: "XML", Expressions: []string{"doc"}, Definition: "USING XML INDEX pxml FOR PATH"},
			},
		}},
	}}}
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "CREATE PRIMARY XML INDEX [pxml] ON [dbo].[t] ([doc]);")
	require.Contains(t, def, "CREATE XML INDEX [sxml] ON [dbo].[t] ([doc])\nUSING XML INDEX [pxml] FOR PATH;")

	from := New()
	require.NoError(t, from.LoadMetadata(t.Context(), &metadata.DatabaseSchemaMetadata{}))
	to := New()
	require.NoError(t, to.LoadMetadata(t.Context(), meta))
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "CREATE PRIMARY XML INDEX [pxml] ON [dbo].[t] ([doc])")
	require.Contains(t, sql, "USING XML INDEX [pxml] FOR PATH")
}

func TestIndexIncludeColumnsSurviveExtractionAndDiff(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, [payload] int NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	to, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([id] int NOT NULL, [payload] int NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));
CREATE INDEX [IX_t_payload] ON [dbo].[t] ([id]) INCLUDE ([payload]);
`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "INCLUDE ([payload])")
}

func TestSequenceNoCacheSurvivesExtraction(t *testing.T) {
	c, err := LoadSQL(`CREATE SEQUENCE [dbo].[s] AS INT START WITH 1 INCREMENT BY 1 NO CACHE;`)
	require.NoError(t, err)
	require.Equal(t, "NO CACHE", c.GetSchema("dbo").Sequences[0].CacheSize)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "NO CACHE")
}

func TestIdentityColumnIsNotNullable(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int IDENTITY(1,1));`)
	require.NoError(t, err)
	require.False(t, c.GetTable("dbo", "t").Columns[0].Nullable)
}

func TestUnnamedDefaultGetsConstraintName(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, [v] int DEFAULT (0), CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, [v] int DEFAULT (1), CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "DROP CONSTRAINT [DF_t_v]")
	require.Contains(t, sql, "ADD CONSTRAINT [DF_t_v] DEFAULT (1) FOR [v]")
}

func TestInboundForeignKeyRebuiltForStandaloneUniqueIndex(t *testing.T) {
	from, err := LoadSQL(`
CREATE TABLE [dbo].[p] ([id] int NOT NULL, CONSTRAINT [PK_p] PRIMARY KEY ([id]));
CREATE UNIQUE INDEX [UX_p_id] ON [dbo].[p] ([id] ASC);
CREATE TABLE [dbo].[c] ([id] int NOT NULL, [p_id] int NOT NULL, CONSTRAINT [PK_c] PRIMARY KEY ([id]), CONSTRAINT [FK_c_p] FOREIGN KEY ([p_id]) REFERENCES [dbo].[p] ([id]));
`)
	require.NoError(t, err)
	to, err := LoadSQL(`
CREATE TABLE [dbo].[p] ([id] int NOT NULL, CONSTRAINT [PK_p] PRIMARY KEY ([id]));
CREATE UNIQUE INDEX [UX_p_id] ON [dbo].[p] ([id] DESC);
CREATE TABLE [dbo].[c] ([id] int NOT NULL, [p_id] int NOT NULL, CONSTRAINT [PK_c] PRIMARY KEY ([id]), CONSTRAINT [FK_c_p] FOREIGN KEY ([p_id]) REFERENCES [dbo].[p] ([id]));
`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "DROP CONSTRAINT [FK_c_p]")
	require.Contains(t, sql, "DROP INDEX [UX_p_id]")
	require.Contains(t, sql, "ADD CONSTRAINT [FK_c_p]")
	require.Less(t, strings.Index(sql, "DROP CONSTRAINT [FK_c_p]"), strings.Index(sql, "DROP INDEX [UX_p_id]"))
}

func TestComputedColumnRebuildDropsDependentCheck(t *testing.T) {
	from, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([a] int NOT NULL, [c] AS ([a] + 1) PERSISTED, CONSTRAINT [CK_c] CHECK ([c] > 0));
`)
	require.NoError(t, err)
	to, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([a] int NOT NULL, [c] AS ([a] + 2) PERSISTED, CONSTRAINT [CK_c] CHECK ([c] > 0));
`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "DROP CONSTRAINT [CK_c]")
	require.Contains(t, sql, "DROP COLUMN [c]")
	require.Contains(t, sql, "ADD CONSTRAINT [CK_c] CHECK ([c] > 0)")
	require.Less(t, strings.Index(sql, "DROP CONSTRAINT [CK_c]"), strings.Index(sql, "DROP COLUMN [c]"))
	require.Less(t, strings.Index(sql, "DROP COLUMN [c]"), strings.Index(sql, "ADD CONSTRAINT [CK_c]"))
}

func TestLoadSQLExtractsXMLIndexes(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([doc] xml NULL);
CREATE PRIMARY XML INDEX [pxml] ON [dbo].[t] ([doc]);
CREATE XML INDEX [sxml] ON [dbo].[t] ([doc]) USING XML INDEX [pxml] FOR PATH;
`)
	require.NoError(t, err)
	indexes := c.GetTable("dbo", "t").Indexes
	require.Len(t, indexes, 2)
	require.Equal(t, "PRIMARY", indexes[0].Definition)
	require.Equal(t, "USING XML INDEX pxml FOR PATH", indexes[1].Definition)
}

func TestSynonymDiffAndMigration(t *testing.T) {
	from := New()
	require.NoError(t, from.LoadMetadata(t.Context(), &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name:     "dbo",
		Synonyms: []*metadata.SynonymMetadata{{Name: "old", TargetSchema: "dbo", TargetName: "a"}},
	}}}))
	to := New()
	require.NoError(t, to.LoadMetadata(t.Context(), &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name:     "dbo",
		Synonyms: []*metadata.SynonymMetadata{{Name: "new", TargetSchema: "dbo", TargetName: "b"}},
	}}}))
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "DROP SYNONYM [dbo].[old]")
	require.Contains(t, sql, "CREATE SYNONYM [dbo].[new] FOR [dbo].[b]")
}

func TestSequenceTypeChangeIsExplicitlyUnsupported(t *testing.T) {
	from, err := LoadSQL(`CREATE SEQUENCE [dbo].[s] AS INT START WITH 1;`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE SEQUENCE [dbo].[s] AS BIGINT START WITH 1;`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "UNSUPPORTED: sequence [dbo].[s] data type change")
	require.NotContains(t, sql, "DROP SEQUENCE")
}

func TestExtendedPropertyEscapesNames(t *testing.T) {
	sql := generateTableCommentSQL("ADD", "dbo", "O'Brien", "it's")
	require.Contains(t, sql, "'TABLE', 'O''Brien'")
	require.Contains(t, sql, "N'it''s'")
}

func TestForeignKeyUnqualifiedTargetUsesDefaultSchema(t *testing.T) {
	c, err := LoadSQL(`
CREATE SCHEMA [app];
CREATE TABLE [dbo].[parent] ([id] int NOT NULL, CONSTRAINT [PK_parent] PRIMARY KEY ([id]));
CREATE TABLE [app].[child] ([id] int NOT NULL, [parent_id] int NULL, CONSTRAINT [PK_child] PRIMARY KEY ([id]), CONSTRAINT [FK_child_parent] FOREIGN KEY ([parent_id]) REFERENCES [parent] ([id]));
`)
	require.NoError(t, err)
	fk := c.GetTable("app", "child").ForeignKeys[0]
	require.Equal(t, "dbo", fk.ReferencedSchema)
}

func TestComputedRebuildDropsCoveringIndex(t *testing.T) {
	from, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([a] int NOT NULL, [c] AS ([a] + 1) PERSISTED);
CREATE INDEX [IX_t] ON [dbo].[t] ([a]) INCLUDE ([c]);
`)
	require.NoError(t, err)
	to, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([a] int NOT NULL, [c] AS ([a] + 2) PERSISTED);
CREATE INDEX [IX_t] ON [dbo].[t] ([a]) INCLUDE ([c]);
`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "DROP INDEX [IX_t]")
	require.Contains(t, sql, "DROP COLUMN [c]")
	require.Contains(t, sql, "INCLUDE ([c])")
}

func TestLoadSQLExtractsSynonyms(t *testing.T) {
	c, err := LoadSQL(`CREATE SYNONYM [dbo].[alias] FOR [srv].[db].[sales].[orders];`)
	require.NoError(t, err)
	synonyms := c.GetSchema("dbo").Synonyms
	require.Len(t, synonyms, 1)
	require.Equal(t, "srv", synonyms[0].TargetServer)
	require.Equal(t, "db", synonyms[0].TargetDatabase)
	require.Equal(t, "sales", synonyms[0].TargetSchema)
	require.Equal(t, "orders", synonyms[0].TargetName)
}

func TestSynonymDropsBeforeSchema(t *testing.T) {
	from := New()
	require.NoError(t, from.LoadMetadata(t.Context(), &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name:     "app",
		Synonyms: []*metadata.SynonymMetadata{{Name: "s", TargetSchema: "dbo", TargetName: "t"}},
	}}}))
	to := New()
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Less(t, strings.Index(sql, "DROP SYNONYM [app].[s]"), strings.Index(sql, "DROP SCHEMA [app]"))
}

func TestDefinitionCreatesFunctionBeforeTableDefault(t *testing.T) {
	c, err := LoadSQL(`
CREATE FUNCTION [dbo].[initial_value]() RETURNS int AS BEGIN RETURN 1 END;
CREATE TABLE [dbo].[t] ([id] int NOT NULL CONSTRAINT [DF_t_id] DEFAULT ([dbo].[initial_value]()), CONSTRAINT [PK_t] PRIMARY KEY ([id]));
`)
	require.NoError(t, err)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Less(t, strings.Index(def, "CREATE FUNCTION [dbo].[initial_value]"), strings.Index(def, "CREATE TABLE [dbo].[t]"))
}

func TestFilteredColumnstoreIndexRendersPredicate(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([active] bit NOT NULL, [amount] int NULL);
CREATE NONCLUSTERED COLUMNSTORE INDEX [ix] ON [dbo].[t] ([amount]) WHERE [active] = 1;
`)
	require.NoError(t, err)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "CREATE NONCLUSTERED COLUMNSTORE INDEX [ix] ON [dbo].[t] (\n    [amount]\n)\nWHERE [active] = 1;")
	empty := New()
	sql := GenerateMigration(empty, c, Diff(empty, c)).SQL()
	require.Contains(t, sql, "WHERE [active] = 1")
}

func TestSynonymKeepsEmptyMultipartComponent(t *testing.T) {
	synonym := &metadata.SynonymMetadata{Name: "s", TargetDatabase: "otherdb", TargetName: "target"}
	require.Equal(t, "[otherdb]..[target]", synonymTarget(synonym))
}

func TestDefinitionOrdersPrimaryXMLBeforeSecondary(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{{
			Name:    "t",
			Columns: []*metadata.ColumnMetadata{{Name: "doc", Type: "xml"}},
			Indexes: []*metadata.IndexMetadata{
				{Name: "A_path", Type: "XML", Expressions: []string{"doc"}, Definition: "USING XML INDEX Z_primary FOR PATH"},
				{Name: "Z_primary", Type: "XML", Expressions: []string{"doc"}, Definition: "PRIMARY"},
			},
		}},
	}}}
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.Less(t, strings.Index(def, "CREATE PRIMARY XML INDEX [Z_primary]"), strings.Index(def, "CREATE XML INDEX [A_path]"))
}

func TestSkipDumpTableIsNotDroppedWithItsSchema(t *testing.T) {
	from := New()
	require.NoError(t, from.LoadMetadata(t.Context(), &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name:   "app",
		Tables: []*metadata.TableMetadata{{Name: "protected", SkipDump: true, Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}}}},
	}}}))
	to := New()
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.NotContains(t, sql, "DROP TABLE [app].[protected]")
	require.NotContains(t, sql, "DROP SCHEMA [app]")
	require.Contains(t, sql, "UNSUPPORTED: schema [app]")
}

func TestCaseDistinctTableNamesCoexist(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE [dbo].[Foo] ([a] int); CREATE TABLE [dbo].[foo] ([b] int);`)
	require.NoError(t, err)
	require.Len(t, c.GetSchema("dbo").Tables, 2)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "[Foo]")
	require.Contains(t, def, "[foo]")
}

func TestIndexIncludeColumnWithComma(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, [a,b] int NULL); CREATE INDEX [ix] ON [dbo].[t] ([id]) INCLUDE ([a,b]);`)
	require.NoError(t, err)
	def, err := c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "INCLUDE ([a,b])")
	require.NotContains(t, def, "INCLUDE ([a], [b])")
}

func TestExtendedPropertyNamedArguments(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([id] int NOT NULL);
EXEC sp_addextendedproperty @value=N'a description', @name=N'MS_Description', @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N't';
`)
	require.NoError(t, err)
	require.Equal(t, "a description", c.GetTable("dbo", "t").Comment)
}

func TestDefinitionSkipsTriggersOnSkippedObjects(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{{
			Name: "hidden", SkipDump: true,
			Triggers: []*metadata.TriggerMetadata{{
				Name: "trg_hidden", Body: "CREATE TRIGGER [dbo].[trg_hidden] ON [dbo].[hidden] AFTER INSERT AS SELECT 1",
			}},
		}},
		Views: []*metadata.ViewMetadata{{
			Name: "hidden_view", SkipDump: true,
			Definition: "CREATE VIEW [dbo].[hidden_view] AS SELECT 1 AS n",
			Triggers: []*metadata.TriggerMetadata{{
				Name: "trg_hidden_view", Body: "CREATE TRIGGER [dbo].[trg_hidden_view] ON [dbo].[hidden_view] AFTER INSERT AS SELECT 1",
			}},
		}},
	}}}
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	require.NotContains(t, def, "trg_hidden")
	require.NotContains(t, def, "hidden")
}

func TestMigrationDropsDefaultAroundColumnTypeChange(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, [v] int NOT NULL CONSTRAINT [DF_v] DEFAULT (0), CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, [v] bigint NOT NULL CONSTRAINT [DF_v] DEFAULT (0), CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "DROP CONSTRAINT [DF_v]")
	require.Contains(t, sql, "ALTER COLUMN [v] bigint NOT NULL")
	require.Contains(t, sql, "ADD CONSTRAINT [DF_v] DEFAULT (0) FOR [v]")
	require.Less(t, strings.Index(sql, "DROP CONSTRAINT [DF_v]"), strings.Index(sql, "ALTER COLUMN [v]"))
	require.Less(t, strings.Index(sql, "ALTER COLUMN [v]"), strings.Index(sql, "ADD CONSTRAINT [DF_v]"))
}

func TestMigrationRebuildsComputedColumnForNullability(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE [dbo].[t] ([a] int NULL, [c] AS ([a] + 1) PERSISTED, CONSTRAINT [PK_t] PRIMARY KEY ([a]));`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([a] int NULL, [c] AS ([a] + 1) PERSISTED NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([a]));`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "DROP COLUMN [c]")
	require.Contains(t, sql, "ADD [c] AS ([a] + 1) PERSISTED NOT NULL")
	require.NotContains(t, sql, "ALTER COLUMN [c] AS")
}

func TestMigrationDropsFunctionDefaultOnceForRemovedColumn(t *testing.T) {
	from, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int AS BEGIN RETURN 1 END;
GO
CREATE TABLE [dbo].[t] (
	[id] int NOT NULL,
	[v] int NOT NULL CONSTRAINT [DF_v] DEFAULT ([dbo].[f]()),
	CONSTRAINT [PK_t] PRIMARY KEY ([id])
);
`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Equal(t, 1, strings.Count(sql, "DROP CONSTRAINT [DF_v]"), sql)
	require.Contains(t, sql, "DROP COLUMN [v]")
}

func TestMigrationSeparatesTableDropFromRoutineCreation(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE PROCEDURE [dbo].[p] AS SELECT 1;`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	dropAt := strings.Index(sql, "DROP TABLE [dbo].[t]")
	goAt := strings.Index(sql[dropAt:], "GO")
	createAt := strings.Index(sql, "CREATE PROCEDURE [dbo].[p]")
	require.GreaterOrEqual(t, dropAt, 0)
	require.GreaterOrEqual(t, goAt, 0)
	require.GreaterOrEqual(t, createAt, 0)
	require.Less(t, dropAt+goAt, createAt)
}

func TestGetTableResolvesCaseDistinctNames(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE [dbo].[Foo] ([a] int); CREATE TABLE [dbo].[foo] ([b] int);`)
	require.NoError(t, err)
	require.Equal(t, "Foo", c.GetTable("dbo", "Foo").Name)
	require.Equal(t, "foo", c.GetTable("dbo", "foo").Name)
}

func TestMigrationUpdatesViewComments(t *testing.T) {
	from, err := LoadSQL(`
CREATE VIEW [dbo].[v] AS SELECT 1 AS n;
GO
EXEC sp_addextendedproperty 'MS_Description', N'old', 'SCHEMA', 'dbo', 'VIEW', 'v';
`)
	require.NoError(t, err)
	to, err := LoadSQL(`
CREATE VIEW [dbo].[v] AS SELECT 1 AS n;
GO
EXEC sp_addextendedproperty 'MS_Description', N'new', 'SCHEMA', 'dbo', 'VIEW', 'v';
`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "DROP VIEW [dbo].[v]")
	require.Contains(t, sql, "sp_addextendedproperty")
	require.Contains(t, sql, "N'new'")
}

func TestMigrationKeepsSkippedTrigger(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{{
			Name:    "t",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: false}},
			Triggers: []*metadata.TriggerMetadata{{
				Name: "trg", SkipDump: true, Body: "CREATE TRIGGER [dbo].[trg] ON [dbo].[t] AFTER INSERT AS SELECT 1",
			}},
		}},
	}}}
	from := New()
	require.NoError(t, from.LoadMetadata(context.Background(), meta))
	to := New()
	require.NoError(t, to.LoadMetadata(context.Background(), proto.Clone(meta).(*metadata.DatabaseSchemaMetadata)))
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.NotContains(t, sql, "DROP TRIGGER")
}

// TestColumnLevelKeyWithColumns checks that a PRIMARY KEY or UNIQUE written
// beside a column and naming its own columns keys those columns, as
// `ALTER TABLE a ADD x a UNIQUE(x)` does: `x int, y int UNIQUE(x)` is a
// unique index on x, and y keeps its nullability.
func TestColumnLevelKeyWithColumns(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE [dbo].[t] ([x] int, [y] int UNIQUE([x]), [z] int PRIMARY KEY ([x]));`)
	require.NoError(t, err)
	table := c.GetTable("dbo", "t")
	require.NotNil(t, table)
	byName := map[string]*metadata.IndexMetadata{}
	for _, ix := range table.Indexes {
		byName[ix.Name] = ix
	}
	unique := byName["UQ_t_1"]
	require.NotNil(t, unique)
	require.Equal(t, []string{"x"}, unique.Expressions)
	primary := byName["PK_t_2"]
	require.NotNil(t, primary)
	require.Equal(t, []string{"x"}, primary.Expressions)
	for _, col := range table.Columns {
		if col.Name == "y" {
			require.True(t, col.Nullable, "y is not the key's column")
		}
		if col.Name == "x" {
			require.False(t, col.Nullable, "the primary key holds x")
		}
	}
}

// TestColumnLevelKeyNamesALaterColumn checks the same for a key whose column
// the CREATE TABLE body has not reached yet: the primary key's columns allow
// no NULL once every column exists.
func TestColumnLevelKeyNamesALaterColumn(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE [dbo].[t] ([x] int PRIMARY KEY ([y]), [y] int NULL);`)
	require.NoError(t, err)
	table := c.GetTable("dbo", "t")
	require.NotNil(t, table)
	for _, col := range table.Columns {
		if col.Name == "y" {
			require.False(t, col.Nullable, "the primary key holds y")
		}
	}
}

func TestLoadMetadataReloadReplacesPrimaryDatabaseRegistration(t *testing.T) {
	c := New()
	ctx := context.Background()

	require.NoError(t, c.LoadMetadata(ctx, &metadata.DatabaseSchemaMetadata{
		Name: "db_old",
		Schemas: []*metadata.SchemaMetadata{{
			Name:   "dbo",
			Tables: []*metadata.TableMetadata{{Name: "old_tbl"}},
		}},
	}))
	require.Equal(t, []string{"db_old"}, c.Databases())
	require.NotNil(t, c.DatabaseCatalog("db_old"))

	require.NoError(t, c.LoadMetadata(ctx, &metadata.DatabaseSchemaMetadata{
		Name: "db_new",
		Schemas: []*metadata.SchemaMetadata{{
			Name:   "dbo",
			Tables: []*metadata.TableMetadata{{Name: "new_tbl"}},
		}},
	}))
	require.Equal(t, []string{"db_new"}, c.Databases())
	require.Nil(t, c.DatabaseCatalog("db_old"))
	require.NotNil(t, c.DatabaseCatalog("db_new"))
	require.NotNil(t, c.GetTable("dbo", "new_tbl"))
}

func TestGetDatabaseMetadata(t *testing.T) {
	t.Run("empty SQL returns default dbo schema", func(t *testing.T) {
		meta, err := GetDatabaseMetadata("")
		require.NoError(t, err)
		require.Len(t, meta.Schemas, 1)
		require.Equal(t, "dbo", meta.Schemas[0].Name)
		require.Empty(t, meta.Schemas[0].Tables)
	})

	t.Run("invalid syntax returns parse error", func(t *testing.T) {
		_, err := GetDatabaseMetadata("CREATE TABLE [dbo].[t] (")
		require.Error(t, err)
	})

	t.Run("extracts tables, constraints, indexes, and routines", func(t *testing.T) {
		meta, err := GetDatabaseMetadata(`
CREATE SCHEMA [sales];
GO
CREATE TABLE [sales].[orders] (
	[id] INT IDENTITY(1,1) NOT NULL,
	[amount] DECIMAL(10,2) CONSTRAINT [DF_orders_amount] DEFAULT (0),
	[status] INT DEFAULT 1,
	CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id])
);
GO
CREATE NONCLUSTERED INDEX [IX_orders_amount] ON [sales].[orders] ([amount] DESC) WHERE [amount] > 0;
`)
		require.NoError(t, err)
		require.Len(t, meta.Schemas, 1)
		require.Equal(t, "sales", meta.Schemas[0].Name)
		require.Len(t, meta.Schemas[0].Tables, 1)
		table := meta.Schemas[0].Tables[0]
		require.Equal(t, "orders", table.Name)
		require.Equal(t, "DF_orders_amount", table.Columns[1].DefaultConstraintName)
		require.Equal(t, "", table.Columns[2].DefaultConstraintName)
		require.Len(t, table.Indexes, 2)
		require.Equal(t, "WHERE [amount] > 0", table.Indexes[1].Definition)
	})

	t.Run("preserves case-distinct schema names without merging tables", func(t *testing.T) {
		ddl := `
CREATE TABLE [s].[t] ([a] INT NOT NULL);
CREATE TABLE [S].[t] ([b] INT NOT NULL);
ALTER TABLE [S].[t] ADD CONSTRAINT [PK_S_t] PRIMARY KEY ([b]);
`
		meta, err := GetDatabaseMetadata(ddl)
		require.NoError(t, err)
		require.Len(t, meta.Schemas, 2)
		require.Equal(t, "S", meta.Schemas[0].Name)
		require.Len(t, meta.Schemas[0].Tables, 1)
		require.Len(t, meta.Schemas[0].Tables[0].Columns, 1)
		require.Equal(t, "b", meta.Schemas[0].Tables[0].Columns[0].Name)
		require.Len(t, meta.Schemas[0].Tables[0].Indexes, 1)
		require.Equal(t, "s", meta.Schemas[1].Name)
		require.Len(t, meta.Schemas[1].Tables, 1)
		require.Len(t, meta.Schemas[1].Tables[0].Columns, 1)
		require.Equal(t, "a", meta.Schemas[1].Tables[0].Columns[0].Name)
		require.Empty(t, meta.Schemas[1].Tables[0].Indexes)

		cat, err := LoadSQL(ddl)
		require.NoError(t, err)
		require.NotNil(t, cat.GetTable("s", "t"))
		require.Equal(t, "a", cat.GetTable("s", "t").Columns[0].Name)
		require.NotNil(t, cat.GetTable("s", "T"))
		require.Equal(t, "a", cat.GetTable("s", "T").Columns[0].Name)
		require.NotNil(t, cat.GetTable("S", "t"))
		require.Equal(t, "b", cat.GetTable("S", "t").Columns[0].Name)
		require.NotNil(t, cat.GetTable("S", "T"))
		require.Equal(t, "b", cat.GetTable("S", "T").Columns[0].Name)
		require.True(t, Diff(cat, cat).IsEmpty())
	})

	t.Run("skips indexes on views instead of failing whole script", func(t *testing.T) {
		ddl := `
CREATE VIEW [dbo].[v] WITH SCHEMABINDING AS SELECT 1 AS [id];
GO
CREATE UNIQUE CLUSTERED INDEX [ix_v] ON [dbo].[v] ([id]);
GO
CREATE TABLE [dbo].[t] ([id] INT NOT NULL);
`
		meta, err := GetDatabaseMetadata(ddl)
		require.NoError(t, err)
		require.Len(t, meta.Schemas, 1)
		require.Len(t, meta.Schemas[0].Views, 1)
		require.Len(t, meta.Schemas[0].Tables, 1)
		require.Equal(t, "t", meta.Schemas[0].Tables[0].Name)

		_, loadErr := LoadSQL(ddl)
		require.ErrorContains(t, loadErr, "indexes on views are not supported")
	})

	t.Run("avoids synthetic DF_ name collisions on unnamed defaults", func(t *testing.T) {
		ddl := `
CREATE TABLE [dbo].[a_b] ([c] INT NOT NULL DEFAULT 0);
CREATE TABLE [dbo].[a] ([b_c] INT NOT NULL DEFAULT 1);
`
		meta, err := GetDatabaseMetadata(ddl)
		require.NoError(t, err)
		require.Len(t, meta.Schemas[0].Tables, 2)
		require.Empty(t, meta.Schemas[0].Tables[0].Columns[0].DefaultConstraintName)
		require.Empty(t, meta.Schemas[0].Tables[1].Columns[0].DefaultConstraintName)
	})
}

func TestGetTableDefinition(t *testing.T) {
	t.Run("nil or SkipDump table returns empty", func(t *testing.T) {
		got, err := GetTableDefinition("dbo", nil, nil)
		require.NoError(t, err)
		require.Empty(t, got)

		got, err = GetTableDefinition("dbo", &metadata.TableMetadata{Name: "t", SkipDump: true}, nil)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("empty schema falls back to dbo", func(t *testing.T) {
		got, err := GetTableDefinition("", &metadata.TableMetadata{
			Name:    "t",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: false}},
		}, nil)
		require.NoError(t, err)
		require.Equal(t, "CREATE TABLE [dbo].[t] (\n    [id] int NOT NULL\n);\n\n", got)
	})

	t.Run("renders constraints, indexes, and comments", func(t *testing.T) {
		table := &metadata.TableMetadata{
			Name:    "orders",
			Comment: "order records",
			Columns: []*metadata.ColumnMetadata{
				{Name: "id", Type: "int", IsIdentity: true, IdentitySeed: 1, IdentityIncrement: 1, Nullable: false},
				{Name: "total", Type: "decimal(10,2)", Default: "((0))", DefaultConstraintName: "DF_orders_total", Nullable: false, Comment: "gross total"},
			},
			Indexes: []*metadata.IndexMetadata{
				{Name: "PK_orders", Type: "CLUSTERED", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"id"}},
				{Name: "IX_orders_total", Type: "NONCLUSTERED", Expressions: []string{"total"}, Descending: []bool{true}, Definition: "WHERE [total] > 0"},
			},
			CheckConstraints: []*metadata.CheckConstraintMetadata{
				{Name: "CK_orders_total", Expression: "([total]>=(0))"},
			},
		}
		got, err := GetTableDefinition("sales", table, nil)
		require.NoError(t, err)
		require.Contains(t, got, "CREATE TABLE [sales].[orders] (")
		require.Contains(t, got, "[id] int IDENTITY(1,1) NOT NULL")
		require.Contains(t, got, "[total] decimal(10,2) CONSTRAINT [DF_orders_total] DEFAULT ((0)) NOT NULL")
		require.Contains(t, got, "CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC)")
		require.Contains(t, got, "CONSTRAINT [CK_orders_total] CHECK ([total]>=(0))")
		require.Contains(t, got, "CREATE NONCLUSTERED INDEX [IX_orders_total] ON\n[sales].[orders] (\n    [total] DESC\n)\nWHERE [total] > 0;")
		require.Contains(t, got, "sp_addextendedproperty")
	})
}

func TestMigrationIgnoresFunctionNameInsideDefaultLiteral(t *testing.T) {
	from, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int AS BEGIN RETURN 1 END;
GO
CREATE TABLE [dbo].[t] ([id] int NOT NULL, [v] varchar(20) NOT NULL CONSTRAINT [DF_v] DEFAULT ('foo'), CONSTRAINT [PK_t] PRIMARY KEY ([id]));
`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] int NOT NULL, [v] varchar(20) NOT NULL CONSTRAINT [DF_v] DEFAULT ('foo'), CONSTRAINT [PK_t] PRIMARY KEY ([id]));`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.NotContains(t, sql, "DROP CONSTRAINT [DF_v]")
}

func TestMigrationReaddsComputedColumnCommentAfterRebuild(t *testing.T) {
	commentRebuildMetadata := func(computedType string) *metadata.DatabaseSchemaMetadata {
		return &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
			Name: "dbo",
			Tables: []*metadata.TableMetadata{{
				Name: "t",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Type: "int", Nullable: false},
					{Name: "a", Type: "int", Nullable: true},
					{Name: "c", Type: computedType, Nullable: true, Comment: "note"},
				},
				Indexes: []*metadata.IndexMetadata{{Name: "PK_t", Type: "CLUSTERED", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"id"}}},
			}},
		}}}
	}
	from := New()
	require.NoError(t, from.LoadMetadata(context.Background(), commentRebuildMetadata("AS ([a] + 1) PERSISTED")))
	to := New()
	require.NoError(t, to.LoadMetadata(context.Background(), commentRebuildMetadata("AS ([a] + 2) PERSISTED")))
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.Contains(t, sql, "DROP COLUMN [c]")
	require.Contains(t, sql, "sp_addextendedproperty")
}

func TestMigrationScopesFunctionDefaultCleanupBySchema(t *testing.T) {
	from, err := LoadSQL(`
CREATE SCHEMA [a];
CREATE SCHEMA [b];
GO
CREATE FUNCTION [a].[f]() RETURNS int AS BEGIN RETURN 1 END;
GO
CREATE FUNCTION [b].[f]() RETURNS int AS BEGIN RETURN 2 END;
GO
CREATE TABLE [b].[t] ([id] int NOT NULL, [v] int NOT NULL CONSTRAINT [DF_v] DEFAULT ([b].[f]()), CONSTRAINT [PK_t] PRIMARY KEY ([id]));
`)
	require.NoError(t, err)
	to, err := LoadSQL(`
CREATE SCHEMA [b];
GO
CREATE FUNCTION [b].[f]() RETURNS int AS BEGIN RETURN 2 END;
GO
CREATE TABLE [b].[t] ([id] int NOT NULL, [v] int NOT NULL CONSTRAINT [DF_v] DEFAULT ([b].[f]()), CONSTRAINT [PK_t] PRIMARY KEY ([id]));
`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.NotContains(t, sql, "DROP CONSTRAINT [DF_v]")
}

func TestDefinitionHoistsTransitiveFunctionDependencies(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{{
			Name: "t",
			Columns: []*metadata.ColumnMetadata{
				{Name: "id", Type: "int", Nullable: false},
				{Name: "v", Type: "int", Nullable: false, Default: "([dbo].[outer_value]())"},
			},
			Indexes: []*metadata.IndexMetadata{{Name: "PK_t", Type: "CLUSTERED", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"id"}}},
		}},
		Functions: []*metadata.FunctionMetadata{
			{Name: "outer_value", Definition: "CREATE FUNCTION [dbo].[outer_value]() RETURNS int AS BEGIN RETURN [dbo].[inner_value]() END"},
			{Name: "inner_value", Definition: "CREATE FUNCTION [dbo].[inner_value]() RETURNS int AS BEGIN RETURN 1 END"},
		},
	}}}
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	innerAt := strings.Index(def, "CREATE FUNCTION [dbo].[inner_value]")
	outerAt := strings.Index(def, "CREATE FUNCTION [dbo].[outer_value]")
	tableAt := strings.Index(def, "CREATE TABLE [dbo].[t]")
	require.GreaterOrEqual(t, innerAt, 0)
	require.GreaterOrEqual(t, outerAt, 0)
	require.GreaterOrEqual(t, tableAt, 0)
	require.Less(t, innerAt, outerAt)
	require.Less(t, outerAt, tableAt)
}

func TestDefinitionHoistsCrossSchemaTransitiveFunctionDependencies(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
		{
			Name: "app",
			Tables: []*metadata.TableMetadata{{
				Name: "t",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Type: "int", Nullable: false},
					{Name: "v", Type: "int", Nullable: false, Default: "([app].[outer_value]())"},
				},
			}},
			Functions: []*metadata.FunctionMetadata{{
				Name:       "outer_value",
				Definition: "CREATE FUNCTION [app].[outer_value]() RETURNS int AS BEGIN RETURN [core].[inner_value]() END",
			}},
		},
		{
			Name: "core",
			Functions: []*metadata.FunctionMetadata{{
				Name:       "inner_value",
				Definition: "CREATE FUNCTION [core].[inner_value]() RETURNS int AS BEGIN RETURN 1 END",
			}},
		},
	}}
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	innerAt := strings.Index(def, "CREATE FUNCTION [core].[inner_value]")
	outerAt := strings.Index(def, "CREATE FUNCTION [app].[outer_value]")
	tableAt := strings.Index(def, "CREATE TABLE [app].[t]")
	require.GreaterOrEqual(t, innerAt, 0)
	require.Less(t, innerAt, outerAt)
	require.Less(t, outerAt, tableAt)
}

func TestDefinitionOrdersSchemaBoundFunctionAfterPrerequisiteTable(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{
			{
				Name:    "lookup",
				Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: false}},
				Indexes: []*metadata.IndexMetadata{{Name: "PK_lookup", Type: "CLUSTERED", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"id"}}},
			},
			{
				Name: "t",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Type: "int", Nullable: false},
					{Name: "v", Type: "int", Nullable: false, Default: "([dbo].[lookup_value]())"},
				},
				Indexes: []*metadata.IndexMetadata{{Name: "PK_t", Type: "CLUSTERED", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"id"}}},
			},
		},
		Functions: []*metadata.FunctionMetadata{{
			Name:       "lookup_value",
			Definition: "CREATE FUNCTION [dbo].[lookup_value]() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN (SELECT MAX([id]) FROM [dbo].[lookup]) END",
		}},
	}}}
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	lookupAt := strings.Index(def, "CREATE TABLE [dbo].[lookup]")
	functionAt := strings.Index(def, "CREATE FUNCTION [dbo].[lookup_value]")
	tableAt := strings.Index(def, "CREATE TABLE [dbo].[t]")
	require.GreaterOrEqual(t, lookupAt, 0)
	require.Less(t, lookupAt, functionAt)
	require.Less(t, functionAt, tableAt)
	// CREATE FUNCTION must start a new batch after CREATE TABLE.
	require.Contains(t, def[lookupAt:functionAt], "GO")
}

func TestDefinitionDefersCyclicSchemaBoundFunctionDefault(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{{
			Name: "t",
			Columns: []*metadata.ColumnMetadata{
				{Name: "id", Type: "int", Nullable: false},
				{Name: "v", Type: "int", Nullable: false, Default: "([dbo].[self_val]())", DefaultConstraintName: "DF_t_v"},
			},
		}},
		Functions: []*metadata.FunctionMetadata{{
			Name:       "self_val",
			Definition: "CREATE FUNCTION [dbo].[self_val]() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN (SELECT ISNULL(MAX([id]), 0) FROM [dbo].[t]) END",
		}},
	}}}
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	tableAt := strings.Index(def, "CREATE TABLE [dbo].[t]")
	fnAt := strings.Index(def, "CREATE FUNCTION [dbo].[self_val]")
	alterAt := strings.Index(def, "ALTER TABLE [dbo].[t] ADD CONSTRAINT [DF_t_v] DEFAULT ([dbo].[self_val]()) FOR [v];")
	require.GreaterOrEqual(t, tableAt, 0)
	require.Less(t, tableAt, fnAt)
	require.Less(t, fnAt, alterAt)
}

func TestExpressionsEqualNormalizesDefaultParenthesizationAndCase(t *testing.T) {
	require.True(t, expressionsEqual("(0)", "((0))"))
	require.True(t, expressionsEqual("0", "((0))"))
	require.True(t, expressionsEqual("GETDATE()", "(getdate())"))
	require.True(t, expressionsEqual("(SYSUTCDATETIME())", "(sysutcdatetime())"))
	require.True(t, expressionsEqual("'active'", "('active')"))
	require.True(t, expressionsEqual("N'hi'", "(N'hi')"))
	require.True(t, expressionsEqual("-1", "((-1))"))

	require.False(t, expressionsEqual("(0)", "(1)"))
	require.False(t, expressionsEqual("'pending'", "'PENDING'"))
	require.False(t, expressionsEqual("(1 + 2) * 3", "1 + 2 * 3"))
	require.False(t, expressionsEqual("", "(0)"))
}

func TestTypesEqualNormalizesImplicitParameters(t *testing.T) {
	require.True(t, typesEqual("datetime2", "datetime2(7)"))
	require.True(t, typesEqual("datetimeoffset", "datetimeoffset(7)"))
	require.True(t, typesEqual("time", "time(7)"))
	require.True(t, typesEqual("decimal", "decimal(18, 0)"))
	require.True(t, typesEqual("decimal(10)", "decimal(10,0)"))
	require.True(t, typesEqual("dec(10)", "decimal(10, 0)"))
	require.True(t, typesEqual("numeric", "numeric(18,0)"))
	require.True(t, typesEqual("numeric(12)", "numeric(12, 0)"))
	require.True(t, typesEqual("float", "float(53)"))
	require.True(t, typesEqual("real", "real(24)"))
	require.True(t, typesEqual("varchar", "varchar(1)"))
	require.True(t, typesEqual("nvarchar", "nvarchar(1)"))
	require.True(t, typesEqual("char", "char(1)"))
	require.True(t, typesEqual("nchar", "nchar(1)"))
	require.True(t, typesEqual("binary", "binary(1)"))
	require.True(t, typesEqual("varbinary", "varbinary(1)"))
	require.True(t, typesEqual("integer", "int"))

	require.False(t, typesEqual("datetime2", "datetime2(3)"))
	require.False(t, typesEqual("decimal(10,2)", "decimal(10,0)"))
	require.False(t, typesEqual("varchar(1)", "varchar(20)"))
	require.False(t, typesEqual("real", "float"))

	from, err := LoadSQL(`CREATE TABLE [dbo].[t] ([v] datetime2 NULL);`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE TABLE [dbo].[t] ([v] datetime2(7) NULL);`)
	require.NoError(t, err)
	require.Empty(t, GenerateMigration(from, to, Diff(from, to)).SQL())
}

func TestMigrationReplacesFunctionUsedByChangedCheckConstraint(t *testing.T) {
	from, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int AS BEGIN RETURN 1 END;
GO
CREATE TABLE [dbo].[t] ([v] int NOT NULL CONSTRAINT [CK_t_v] CHECK ([dbo].[f]() = 1));
`)
	require.NoError(t, err)
	to, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int AS BEGIN RETURN 2 END;
GO
CREATE TABLE [dbo].[t] ([v] int NOT NULL CONSTRAINT [CK_t_v] CHECK ([dbo].[f]() = 2));
`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.NotContains(t, sql, "UNSUPPORTED")
	dropCheck := strings.Index(sql, "DROP CONSTRAINT [CK_t_v]")
	dropFn := strings.Index(sql, "DROP FUNCTION [dbo].[f]")
	createFn := strings.Index(sql, "CREATE FUNCTION [dbo].[f]")
	addCheck := strings.Index(sql, "ADD CONSTRAINT [CK_t_v] CHECK ([dbo].[f]() = 2)")
	require.GreaterOrEqual(t, dropCheck, 0)
	require.Less(t, dropCheck, dropFn)
	require.Less(t, dropFn, createFn)
	require.Less(t, createFn, addCheck)
}

func TestMigrationReplacesFunctionUsedByUnchangedCheckConstraint(t *testing.T) {
	from, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int AS BEGIN RETURN 1 END;
GO
CREATE TABLE [dbo].[t] ([v] int NOT NULL CONSTRAINT [CK_t_v] CHECK ([dbo].[f]() > 0));
`)
	require.NoError(t, err)
	to, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int AS BEGIN RETURN 2 END;
GO
CREATE TABLE [dbo].[t] ([v] int NOT NULL CONSTRAINT [CK_t_v] CHECK ([dbo].[f]() > 0));
`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.NotContains(t, sql, "UNSUPPORTED")
	dropCheck := strings.Index(sql, "DROP CONSTRAINT [CK_t_v]")
	dropFn := strings.Index(sql, "DROP FUNCTION [dbo].[f]")
	createFn := strings.Index(sql, "CREATE FUNCTION [dbo].[f]")
	addCheck := strings.Index(sql, "ADD CONSTRAINT [CK_t_v] CHECK ([dbo].[f]() > 0)")
	require.GreaterOrEqual(t, dropCheck, 0)
	require.Less(t, dropCheck, dropFn)
	require.Less(t, dropFn, createFn)
	require.Less(t, createFn, addCheck)
}

func TestMigrationRebuildsSchemaBoundViewsCallingReplacedFunction(t *testing.T) {
	from, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN 1 END;
GO
CREATE VIEW [dbo].[v1] WITH SCHEMABINDING AS SELECT [dbo].[f]() AS n;
GO
CREATE VIEW [dbo].[v2] WITH SCHEMABINDING AS SELECT [n] FROM [dbo].[v1];
`)
	require.NoError(t, err)
	to, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN 2 END;
GO
CREATE VIEW [dbo].[v1] WITH SCHEMABINDING AS SELECT [dbo].[f]() AS n;
GO
CREATE VIEW [dbo].[v2] WITH SCHEMABINDING AS SELECT [n] FROM [dbo].[v1];
`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.NotContains(t, sql, "UNSUPPORTED")
	dropV2 := strings.Index(sql, "DROP VIEW [dbo].[v2]")
	dropV1 := strings.Index(sql, "DROP VIEW [dbo].[v1]")
	dropFn := strings.Index(sql, "DROP FUNCTION [dbo].[f]")
	createFn := strings.Index(sql, "CREATE FUNCTION [dbo].[f]")
	createV1 := strings.Index(sql, "CREATE VIEW [dbo].[v1]")
	createV2 := strings.Index(sql, "CREATE VIEW [dbo].[v2]")
	require.GreaterOrEqual(t, dropV2, 0)
	require.Less(t, dropV2, dropV1)
	require.Less(t, dropV1, dropFn)
	require.Less(t, dropFn, createFn)
	require.Less(t, createFn, createV1)
	require.Less(t, createV1, createV2)
}

func TestMigrationDropsSchemaBoundFunctionBeforeDroppedTable(t *testing.T) {
	t.Run("both function and table dropped", func(t *testing.T) {
		from, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([id] int NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));
GO
CREATE FUNCTION [dbo].[f]() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN (SELECT MAX([id]) FROM [dbo].[t]) END;
`)
		require.NoError(t, err)
		to := New()
		sql := GenerateMigration(from, to, Diff(from, to)).SQL()
		require.NotContains(t, sql, "UNSUPPORTED")
		dropFn := strings.Index(sql, "DROP FUNCTION [dbo].[f]")
		dropTable := strings.Index(sql, "DROP TABLE [dbo].[t]")
		require.GreaterOrEqual(t, dropFn, 0)
		require.Less(t, dropFn, dropTable)
	})

	t.Run("function kept still blocks table drop", func(t *testing.T) {
		from, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([id] int NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY ([id]));
GO
CREATE FUNCTION [dbo].[f]() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN (SELECT MAX([id]) FROM [dbo].[t]) END;
`)
		require.NoError(t, err)
		to, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN (SELECT MAX([id]) FROM [dbo].[t]) END;
`)
		require.NoError(t, err)
		sql := GenerateMigration(from, to, Diff(from, to)).SQL()
		require.Contains(t, sql, "UNSUPPORTED")
		require.NotContains(t, sql, "DROP TABLE [dbo].[t]")
	})
}

func TestDefinitionHoistsUnicodeFunctionDependencies(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Name: "dbo",
		Tables: []*metadata.TableMetadata{{
			Name: "t",
			Columns: []*metadata.ColumnMetadata{
				{Name: "id", Type: "int", Nullable: false},
				{Name: "v", Type: "int", Nullable: false, Default: "(dbo.外层函数())"},
			},
		}},
		Functions: []*metadata.FunctionMetadata{
			{Name: "外层函数", Definition: "CREATE FUNCTION dbo.外层函数() RETURNS int AS BEGIN RETURN dbo.内层函数() END"},
			{Name: "内层函数", Definition: "CREATE FUNCTION dbo.内层函数() RETURNS int AS BEGIN RETURN 1 END"},
		},
	}}}
	def, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	innerAt := strings.Index(def, "CREATE FUNCTION dbo.内层函数")
	outerAt := strings.Index(def, "CREATE FUNCTION dbo.外层函数")
	tableAt := strings.Index(def, "CREATE TABLE [dbo].[t]")
	require.GreaterOrEqual(t, innerAt, 0)
	require.Less(t, innerAt, outerAt)
	require.Less(t, outerAt, tableAt)
}

func TestMigrationRebuildsMixedSchemaBoundFunctionAndViewGraph(t *testing.T) {
	from, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN 1 END;
GO
CREATE VIEW [dbo].[v1] WITH SCHEMABINDING AS SELECT [dbo].[f]() AS [n];
GO
CREATE FUNCTION [dbo].[g]() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN (SELECT MAX([n]) FROM [dbo].[v1]) END;
GO
CREATE VIEW [dbo].[v2] WITH SCHEMABINDING AS SELECT [dbo].[g]() AS [m];
`)
	require.NoError(t, err)
	to, err := LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN 2 END;
GO
CREATE VIEW [dbo].[v1] WITH SCHEMABINDING AS SELECT [dbo].[f]() AS [n];
GO
CREATE FUNCTION [dbo].[g]() RETURNS int WITH SCHEMABINDING AS BEGIN RETURN (SELECT MAX([n]) FROM [dbo].[v1]) END;
GO
CREATE VIEW [dbo].[v2] WITH SCHEMABINDING AS SELECT [dbo].[g]() AS [m];
`)
	require.NoError(t, err)
	sql := GenerateMigration(from, to, Diff(from, to)).SQL()
	require.NotContains(t, sql, "UNSUPPORTED")

	dropV2 := strings.Index(sql, "DROP VIEW [dbo].[v2]")
	dropG := strings.Index(sql, "DROP FUNCTION [dbo].[g]")
	dropV1 := strings.Index(sql, "DROP VIEW [dbo].[v1]")
	dropF := strings.Index(sql, "DROP FUNCTION [dbo].[f]")
	createF := strings.Index(sql, "CREATE FUNCTION [dbo].[f]")
	createV1 := strings.Index(sql, "CREATE VIEW [dbo].[v1]")
	createG := strings.Index(sql, "CREATE FUNCTION [dbo].[g]")
	createV2 := strings.Index(sql, "CREATE VIEW [dbo].[v2]")

	require.GreaterOrEqual(t, dropV2, 0)
	require.Less(t, dropV2, dropG)
	require.Less(t, dropG, dropV1)
	require.Less(t, dropV1, dropF)
	require.Less(t, dropF, createF)
	require.Less(t, createF, createV1)
	require.Less(t, createV1, createG)
	require.Less(t, createG, createV2)
}
