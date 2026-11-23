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

	// Double-quoted identifiers, lone '!', and SQL comments must normalize accurately.
	require.False(t, expressionsEqual(`"x" + 1`, `"x" + 2`))
	require.True(t, expressionsEqual(`("x") + (1)`, `"x" + 1`))
	require.False(t, expressionsEqual(`!a`, `!b`))
	require.False(t, expressionsEqual("[X] -- [Y]\n > 0", "[X] - - [Y] > 0"))
	require.True(t, expressionsEqual("[X] /* comment */ > 0", "[X] > 0"))
}

func TestExecReviewRegressions(t *testing.T) {
	// 1. User-defined / alias types are preserved by GetDatabaseMetadata.
	meta, err := GetDatabaseMetadata(`
CREATE TABLE [dbo].[contacts] ([id] INT NOT NULL, [phone] dbo.PhoneNumber NULL);
GO
CREATE PROCEDURE [dbo].[p_contact] @p dbo.PhoneNumber AS BEGIN SELECT 1; END;
GO
CREATE FUNCTION [dbo].[f_contact](@p dbo.PhoneNumber) RETURNS dbo.PhoneNumber AS BEGIN RETURN @p; END;
`)
	require.NoError(t, err)
	require.Len(t, meta.Schemas, 1)
	require.Len(t, meta.Schemas[0].Tables, 1)
	require.Equal(t, "dbo.PhoneNumber", meta.Schemas[0].Tables[0].Columns[1].Type)
	require.Len(t, meta.Schemas[0].Procedures, 1)
	require.Len(t, meta.Schemas[0].Functions, 1)

	// 2. ALTER TABLE rolls back earlier actions when a later action fails.
	c, err := LoadSQL(`CREATE TABLE [dbo].[t] ([a] INT NOT NULL);`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ADD [b] INT, [a] INT`, ErrDupColumn)
	require.Nil(t, columnByName(c.GetTable("dbo", "t"), "b"), "rolled-back column must not remain")

	// 3. ALTER SEQUENCE preserves CACHE size when sequence/schema name contains "cache".
	c, err = LoadSQL(`CREATE SEQUENCE [dbo].[cached_seq] CACHE 10;`)
	require.NoError(t, err)
	execExpectOK(t, c, `ALTER SEQUENCE [dbo].[cached_seq] INCREMENT BY 2`)
	require.Equal(t, "10", c.GetSchema("dbo").Sequences[0].CacheSize)
	execExpectOK(t, c, `ALTER SEQUENCE [dbo].[cached_seq] CACHE`)
	require.Equal(t, "", c.GetSchema("dbo").Sequences[0].CacheSize)

	// 4. TRUNCATE TABLE allows self-referencing foreign keys while blocking external ones.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[node] ([id] INT NOT NULL PRIMARY KEY, [parent_id] INT NULL REFERENCES [dbo].[node]([id]));
CREATE TABLE [dbo].[child] ([id] INT NOT NULL, [node_id] INT NULL REFERENCES [dbo].[node]([id]));
`)
	require.NoError(t, err)
	execExpectCode(t, c, `TRUNCATE TABLE [dbo].[node]`, ErrTruncateReferencedByFK)
	execExpectOK(t, c, `DROP TABLE [dbo].[child]`)
	execExpectOK(t, c, `TRUNCATE TABLE [dbo].[node]`)
	execExpectOK(t, c, `DROP TABLE [dbo].[node]`)

	// 5. Case-distinct tables reject ambiguous names without panicking.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[Foo] ([a] INT NOT NULL);
CREATE TABLE [dbo].[foo] ([a] INT NOT NULL);
`)
	require.NoError(t, err)
	execExpectCode(t, c, `CREATE INDEX [ix] ON [dbo].[FOO]([a])`, ErrObjectNotFound)
	execExpectCode(t, c, `ALTER TABLE [dbo].[FOO] ADD [b] INT`, ErrAlterTableNotFound)
	execExpectCode(t, c, `TRUNCATE TABLE [dbo].[FOO]`, ErrTruncateNotFound)
	execExpectOK(t, c, `CREATE INDEX [ix] ON [dbo].[Foo]([a])`)
	require.Len(t, c.GetTable("dbo", "Foo").Indexes, 1)
	require.Empty(t, c.GetTable("dbo", "foo").Indexes)

	// 6. ALTER SCHEMA TRANSFER preserves identifier boundaries with embedded dots.
	c, err = LoadSQL(`
CREATE SCHEMA [app];
GO
CREATE TABLE [dbo].[a.b] ([id] INT NOT NULL);
`)
	require.NoError(t, err)
	execExpectOK(t, c, `ALTER SCHEMA [app] TRANSFER [dbo].[a.b]`)
	require.Nil(t, c.GetTable("dbo", "a.b"))
	require.NotNil(t, c.GetTable("app", "a.b"))

	// 7. Opaque filegroup names are preserved during GetDatabaseMetadata.
	meta, err = GetDatabaseMetadata(`
CREATE TABLE [dbo].[t] ([a] INT NOT NULL);
CREATE INDEX [ix] ON [dbo].[t]([a]) ON [data_fg];
`)
	require.NoError(t, err)
	require.Len(t, meta.Schemas[0].Tables[0].Indexes, 1)

	// 8. CREATE INDEX WITH (DROP_EXISTING = ON) rebuilds an existing index in place.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] ([a] INT NOT NULL, [b] INT NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY CLUSTERED ([a]));
CREATE INDEX [ix] ON [dbo].[t]([a]);
`)
	require.NoError(t, err)
	execExpectOK(t, c, `CREATE INDEX [ix] ON [dbo].[t]([b]) WITH (DROP_EXISTING = ON)`)
	ix := indexByName(c.GetTable("dbo", "t"), "ix")
	require.NotNil(t, ix)
	require.Equal(t, []string{"b"}, ix.Expressions)
	execExpectOK(t, c, `CREATE UNIQUE CLUSTERED INDEX [PK_t] ON [dbo].[t]([a]) WITH (DROP_EXISTING = ON)`)
	pk := indexByName(c.GetTable("dbo", "t"), "PK_t")
	require.NotNil(t, pk)
	require.True(t, pk.Primary)
	require.True(t, pk.IsConstraint)
	execExpectCode(t, c, `CREATE INDEX [missing_ix] ON [dbo].[t]([a]) WITH (DROP_EXISTING = ON)`, ErrIndexNotFound)

	// 9. Failed CREATE SCHEMA rolls back cached tables in tableByName.
	c = New()
	execExpectCode(t, c, `CREATE SCHEMA [app] CREATE TABLE [t] ([a] INT) CREATE TABLE [t] ([b] INT)`, ErrDupObject)
	require.Nil(t, c.GetSchema("app"))
	require.Nil(t, c.GetTable("app", "t"))
	execExpectOK(t, c, `CREATE SCHEMA [app]`)
	execExpectOK(t, c, `CREATE TABLE [app].[t] ([c] INT)`)
	appSchema := c.GetSchema("app")
	require.Len(t, appSchema.Tables, 1)
	require.Len(t, appSchema.Tables[0].Columns, 1)
	require.Equal(t, "c", appSchema.Tables[0].Columns[0].Name)

	// 10. DROP COLUMN rejects columns referenced by a CHECK constraint.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([a] INT, [b] INT, CONSTRAINT [ck_a] CHECK ([a] > 0));`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [a]`, ErrColumnDependentObject)
	require.NotNil(t, columnByName(c.GetTable("dbo", "t"), "a"))
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP CONSTRAINT [ck_a]`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [a]`)

	// 11. Filtered index IDENTITYCOL check inspects ColumnRef nodes rather than substrings.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([a] INT, [identitycol_value] INT, [tag] NVARCHAR(50));`)
	require.NoError(t, err)
	execExpectOK(t, c, `CREATE INDEX [ix1] ON [dbo].[t]([a]) WHERE [identitycol_value] IS NOT NULL`)
	execExpectOK(t, c, `CREATE INDEX [ix2] ON [dbo].[t]([a]) WHERE [tag] = 'IDENTITYCOL'`)
	execExpectCode(t, c, `CREATE INDEX [ix3] ON [dbo].[t]([a]) WHERE IDENTITYCOL IS NOT NULL`, ErrInvalidColumnName)

	// 12. ALTER TABLE DROP CONSTRAINT WITH (ONLINE = OFF/ON) honors the option value and constraint kind.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([id] INT NOT NULL, [u] INT NOT NULL, CONSTRAINT [PK_t] PRIMARY KEY CLUSTERED ([id]), CONSTRAINT [UQ_t_u] UNIQUE NONCLUSTERED ([u]));`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] DROP CONSTRAINT [UQ_t_u] WITH (ONLINE = ON)`, ErrOnlineDropNonClustered)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP CONSTRAINT [UQ_t_u] WITH (ONLINE = OFF)`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP CONSTRAINT [PK_t] WITH (ONLINE = ON)`)
	require.Empty(t, c.GetTable("dbo", "t").Indexes)

	// 13. Opaque execution principals and external schema/routine attributes are preserved by GetDatabaseMetadata.
	meta, err = GetDatabaseMetadata(`
CREATE SCHEMA [app] AUTHORIZATION [app_owner];
GO
CREATE PROCEDURE [app].[p] WITH EXECUTE AS 'app_user' AS BEGIN SELECT 1; END;
GO
CREATE FUNCTION [app].[f]() RETURNS INT WITH EXECUTE AS 'app_user' AS BEGIN RETURN 1; END;
`)
	require.NoError(t, err)
	require.Len(t, meta.Schemas, 1)
	require.Len(t, meta.Schemas[0].Procedures, 1)
	require.Len(t, meta.Schemas[0].Functions, 1)

	// 14. DROP COLUMN rejects columns referenced by computed columns.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([a] INT, [b] AS ([a] + 1) PERSISTED);`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [a]`, ErrColumnDependentObject)
	require.NotNil(t, columnByName(c.GetTable("dbo", "t"), "a"))
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [b]`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [a]`)

	// 15. Multi-target DROP INDEX validates all targets before mutating any table.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] ([a] INT, [b] INT);
CREATE INDEX [ix] ON [dbo].[t]([a]);
`)
	require.NoError(t, err)
	execExpectCode(t, c, `DROP INDEX [ix] ON [dbo].[t], [missing] ON [dbo].[t]`, ErrCannotDropNotFound)
	require.NotNil(t, indexByName(c.GetTable("dbo", "t"), "ix"), "failed multi-target DROP INDEX must not remove earlier indexes")
}
