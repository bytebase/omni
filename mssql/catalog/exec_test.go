package catalog

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func execExpectCode(t *testing.T, c *Catalog, stmt string, wantCode int) {
	t.Helper()
	results, err := c.Exec(stmt, nil)
	require.NoError(t, err, "unexpected parse error for %q", stmt)
	var gotErr error
	for _, r := range results {
		if r.Error != nil {
			gotErr = r.Error
			break
		}
	}
	require.Error(t, gotErr, "expected Msg %d for %q", wantCode, stmt)
	var ce *Error
	require.True(t, errors.As(gotErr, &ce), "expected *catalog.Error, got %T (%v)", gotErr, gotErr)
	require.Equal(t, wantCode, ce.Code, "stmt %q: %s", stmt, ce.Message)
}

func execExpectOK(t *testing.T, c *Catalog, stmt string) {
	t.Helper()
	results, err := c.Exec(stmt, nil)
	require.NoError(t, err, "parse error for %q", stmt)
	for _, r := range results {
		require.NoError(t, r.Error, "exec error for %q", stmt)
	}
}

func TestExecDropAndTruncateErrors(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE [dbo].[p] ([id] INT NOT NULL, CONSTRAINT [PK_p] PRIMARY KEY CLUSTERED ([id]));
CREATE TABLE [dbo].[a] ([id] INT NOT NULL, [pid] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([id]), CONSTRAINT [FK_a_p] FOREIGN KEY ([pid]) REFERENCES [dbo].[p] ([id]));
GO
CREATE VIEW [dbo].[v] AS SELECT [id] FROM [dbo].[p];
`)
	require.NoError(t, err)

	// Missing object vs IF EXISTS.
	execExpectCode(t, c, "DROP TABLE [dbo].[missing]", ErrCannotDropNotFound)
	execExpectOK(t, c, "DROP TABLE IF EXISTS [dbo].[missing]")
	execExpectCode(t, c, "TRUNCATE TABLE [dbo].[missing]", ErrTruncateNotFound)

	// Wrong object kind.
	execExpectCode(t, c, "DROP TABLE [dbo].[v]", ErrCannotDropWrongKind)
	execExpectCode(t, c, "DROP TABLE IF EXISTS [dbo].[v]", ErrCannotDropWrongKind)
	execExpectCode(t, c, "DROP VIEW [dbo].[a]", ErrCannotDropWrongKind)
	execExpectCode(t, c, "TRUNCATE TABLE [dbo].[v]", ErrTruncateNotTable)

	// Foreign key protection on parent table.
	execExpectCode(t, c, "DROP TABLE [dbo].[p]", ErrDropTableReferencedByFK)
	execExpectCode(t, c, "DROP TABLE IF EXISTS [dbo].[p]", ErrDropTableReferencedByFK)
	execExpectCode(t, c, "TRUNCATE TABLE [dbo].[p]", ErrTruncateReferencedByFK)

	// Duplicate drop fails atomically without dropping the table.
	execExpectCode(t, c, "DROP TABLE [dbo].[a], [dbo].[a]", ErrCannotDropNotFound)
	require.NotNil(t, c.GetTable("dbo", "a"))

	// Dropping child table unblocks parent drop.
	execExpectOK(t, c, "DROP TABLE [dbo].[a]")
	require.Nil(t, c.GetTable("dbo", "a"))
	execExpectOK(t, c, "TRUNCATE TABLE [dbo].[p]")
	execExpectOK(t, c, "DROP VIEW [dbo].[v]")
	execExpectOK(t, c, "DROP TABLE [dbo].[p]")
	require.Nil(t, c.GetTable("dbo", "p"))
}

func TestExecAlterTableLifecycleAndErrors(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([id] INT NOT NULL, [val] INT NULL, CONSTRAINT [PK_t] PRIMARY KEY CLUSTERED ([id]));
GO
CREATE VIEW [dbo].[v] AS SELECT [id] FROM [dbo].[t];
`)
	require.NoError(t, err)

	execExpectCode(t, c, "ALTER TABLE [dbo].[missing] ADD [x] INT", ErrAlterTableNotFound)
	execExpectCode(t, c, "ALTER TABLE [dbo].[v] ADD [x] INT", ErrAlterTableNotTable)
	execExpectCode(t, c, "ALTER TABLE [dbo].[t] ADD [val] INT", ErrDupColumn)
	execExpectCode(t, c, "ALTER TABLE [dbo].[t] ADD [x] unknown_type", ErrUnknownColumnType)
	execExpectCode(t, c, "ALTER TABLE [dbo].[t] DROP COLUMN [missing]", ErrAlterColNotFound)
	execExpectCode(t, c, "ALTER TABLE [dbo].[t] DROP COLUMN [id]", ErrColumnDependentObject)
	execExpectCode(t, c, "ALTER TABLE [dbo].[t] DROP CONSTRAINT [missing]", ErrConstraintNotFound)
	execExpectOK(t, c, "ALTER TABLE [dbo].[t] DROP CONSTRAINT IF EXISTS [missing]")

	// Add column, default, check, unique, then drop them cleanly.
	execExpectOK(t, c, "ALTER TABLE [dbo].[t] ADD [note] NVARCHAR(50) NULL")
	require.NotNil(t, columnByName(c.GetTable("dbo", "t"), "note"))

	execExpectOK(t, c, "ALTER TABLE [dbo].[t] ADD CONSTRAINT [DF_t_val] DEFAULT (0) FOR [val]")
	require.Equal(t, "(0)", columnByName(c.GetTable("dbo", "t"), "val").Default)
	execExpectCode(t, c, "ALTER TABLE [dbo].[t] DROP COLUMN [val]", ErrColumnDependentObject)

	execExpectOK(t, c, "ALTER TABLE [dbo].[t] DROP CONSTRAINT [DF_t_val]")
	execExpectOK(t, c, "ALTER TABLE [dbo].[t] DROP COLUMN [val]")
	require.Nil(t, columnByName(c.GetTable("dbo", "t"), "val"))

	// Computed column constant non-nullability inference.
	execExpectOK(t, c, "ALTER TABLE [dbo].[t] ADD [one] AS 1 PERSISTED")
	oneCol := columnByName(c.GetTable("dbo", "t"), "one")
	require.NotNil(t, oneCol)
	require.False(t, oneCol.Nullable)
}

func TestExecSchemaSequenceAndIndexOperations(t *testing.T) {
	c, err := LoadSQL(`
CREATE SCHEMA [app];
GO
CREATE TABLE [dbo].[t] ([id] INT NOT NULL, [x] INT NULL, CONSTRAINT [PK_t] PRIMARY KEY CLUSTERED ([id]));
CREATE SEQUENCE [dbo].[seq] AS BIGINT START WITH 1 INCREMENT BY 1;
`)
	require.NoError(t, err)

	// Duplicate schema / missing schema / transfer.
	execExpectCode(t, c, "CREATE SCHEMA [app]", ErrDupObject)
	execExpectCode(t, c, "DROP SCHEMA [missing]", ErrSecurableNotFound)
	execExpectOK(t, c, "DROP SCHEMA IF EXISTS [missing]")
	execExpectOK(t, c, "ALTER SCHEMA [app] TRANSFER [dbo].[t]")
	require.Nil(t, c.GetTable("dbo", "t"))
	require.NotNil(t, c.GetTable("app", "t"))
	execExpectCode(t, c, "DROP SCHEMA [app]", ErrDropSchemaNotEmpty)

	// Index validation & drop.
	execExpectCode(t, c, "CREATE INDEX [ix] ON [app].[t] ([x], [x])", ErrDupIndexColumn)
	execExpectCode(t, c, "CREATE INDEX [ix] ON [app].[t] ([missing])", ErrIndexColumnNotFound)
	execExpectCode(t, c, "CREATE CLUSTERED INDEX [ix] ON [app].[t] ([x])", ErrMultipleClusteredIndex)
	execExpectOK(t, c, "CREATE INDEX [ix] ON [app].[t] ([x])")
	execExpectCode(t, c, "CREATE INDEX [ix] ON [app].[t] ([x])", ErrDupIndex)
	execExpectOK(t, c, "DROP INDEX [ix] ON [app].[t]")
	execExpectCode(t, c, "DROP INDEX [ix] ON [app].[t]", ErrCannotDropNotFound)

	// Sequence alter & validation.
	execExpectCode(t, c, "CREATE SEQUENCE [dbo].[bad] AS not_a_type", ErrUnknownType)
	execExpectCode(t, c, "ALTER SEQUENCE [dbo].[missing] RESTART", ErrSecurableNotFound)
	execExpectOK(t, c, "ALTER SEQUENCE [dbo].[seq] INCREMENT BY 5 MINVALUE 1 MAXVALUE 100 CYCLE CACHE 10")
	seq := c.GetSchema("dbo").Sequences[0]
	require.Equal(t, "5", seq.Increment)
	require.Equal(t, "1", seq.MinValue)
	require.Equal(t, "100", seq.MaxValue)
	require.True(t, seq.Cycle)
	require.Equal(t, "10", seq.CacheSize)
	execExpectOK(t, c, "ALTER SEQUENCE [dbo].[seq] NO MINVALUE NO MAXVALUE NO CYCLE NO CACHE")
	require.Equal(t, "", seq.MinValue)
	require.Equal(t, "", seq.MaxValue)
	require.False(t, seq.Cycle)
	require.Equal(t, "NO CACHE", seq.CacheSize)
}

func TestExecExpressionAndSequenceDiffNormalization(t *testing.T) {
	require.True(t, expressionsEqual("[d] > 0", "([d]>(0))"))
	require.True(t, expressionsEqual("1 = 1", "((1)=(1))"))
	require.True(t, expressionsEqual("1", "((1))"))
	require.False(t, expressionsEqual("'a  b'", "'a b'"))
	require.False(t, expressionsEqual("'pending'", "'PENDING'"))

	from, err := LoadSQL(`CREATE SEQUENCE [dbo].[s];`)
	require.NoError(t, err)
	to, err := LoadSQL(`CREATE SEQUENCE [dbo].[s] AS bigint START WITH -9223372036854775808 INCREMENT BY 1 MINVALUE -9223372036854775808 MAXVALUE 9223372036854775807 NO CYCLE;`)
	require.NoError(t, err)
	require.True(t, Diff(from, to).IsEmpty())
	require.True(t, Diff(to, from).IsEmpty())
}
