package catalog

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
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
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([a] INT, [b] AS ([a] + 1) PERSISTED, [c] INT);`)
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

	// 16. DROP CONSTRAINT rejects dropping a PK/UNIQUE constraint referenced by an inbound FK.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[p] ([id] INT NOT NULL, CONSTRAINT [PK_p] PRIMARY KEY CLUSTERED ([id]));
CREATE TABLE [dbo].[c] ([pid] INT NULL, CONSTRAINT [FK_c_p] FOREIGN KEY ([pid]) REFERENCES [dbo].[p]([id]));
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[p] DROP CONSTRAINT [PK_p]`, ErrDropConstraintReferencedByFK)
	require.NotNil(t, indexByName(c.GetTable("dbo", "p"), "PK_p"))
	execExpectOK(t, c, `ALTER TABLE [dbo].[c] DROP CONSTRAINT [FK_c_p]`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[p] DROP CONSTRAINT [PK_p]`)

	// 17. DROP COLUMN rejects columns referenced by a filtered index WHERE predicate.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] ([a] INT, [b] INT);
CREATE INDEX [ix] ON [dbo].[t]([a]) WHERE [b] > 0;
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [b]`, ErrColumnDependentObject)
	require.NotNil(t, columnByName(c.GetTable("dbo", "t"), "b"))
	execExpectOK(t, c, `DROP INDEX [ix] ON [dbo].[t]`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [b]`)

	// 18. ADD CONSTRAINT DEFAULT rejects adding a second default to a column that already has one.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([a] INT CONSTRAINT [df_old] DEFAULT (0));`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ADD CONSTRAINT [df_new] DEFAULT (1) FOR [a]`, ErrColumnAlreadyHasDefault)
	require.Equal(t, "(0)", columnByName(c.GetTable("dbo", "t"), "a").Default)
	require.Equal(t, "df_old", columnByName(c.GetTable("dbo", "t"), "a").DefaultConstraintName)

	// 19. Named constraints occupy the schema object namespace.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([a] INT, CONSTRAINT [ck] CHECK ([a] > 0));`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ADD CONSTRAINT [ck] CHECK ([a] < 10)`, ErrDupObject)
	require.Len(t, c.GetTable("dbo", "t").CheckConstraints, 1)

	// 20. DROP TABLE deduplicates by resolved table identity for case-distinct tables.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[Foo] ([a] INT);
CREATE TABLE [dbo].[foo] ([a] INT);
`)
	require.NoError(t, err)
	execExpectOK(t, c, `DROP TABLE [dbo].[Foo], [dbo].[foo]`)
	require.Nil(t, c.GetTable("dbo", "Foo"))
	require.Nil(t, c.GetTable("dbo", "foo"))

	// 21. ALTER COLUMN rejects making a primary-key column nullable.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([id] INT NOT NULL PRIMARY KEY);`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [id] INT NULL`, ErrColumnDependentObject)
	require.False(t, columnByName(c.GetTable("dbo", "t"), "id").Nullable)

	// 22. ALTER SCHEMA TRANSFER moves SYSTEM_TIME period tracking with the table.
	c, err = LoadSQL(`
CREATE SCHEMA [app];
GO
CREATE TABLE [dbo].[t] (
    [id] INT NOT NULL PRIMARY KEY,
    [sys_start] DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
    [sys_end] DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
    PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])
);
`)
	require.NoError(t, err)
	execExpectOK(t, c, `ALTER SCHEMA [app] TRANSFER [dbo].[t]`)
	execExpectOK(t, c, `ALTER TABLE [app].[t] SET (SYSTEM_VERSIONING = ON)`)
	execExpectOK(t, c, `CREATE TABLE [dbo].[t] ([id] INT NOT NULL)`)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = ON)`, ErrSystemVersioningNoPeriod)

	// 23. Selective XML indexes preserve promoted paths/namespaces in diff and DDL/migration rendering.
	fromXML, err := LoadSQL(`
CREATE TABLE [dbo].[docs] ([id] INT NOT NULL PRIMARY KEY, [body] XML);
CREATE SELECTIVE XML INDEX [sxi_docs] ON [dbo].[docs]([body]) FOR (p = '/a');
`)
	require.NoError(t, err)
	toXML, err := LoadSQL(`
CREATE TABLE [dbo].[docs] ([id] INT NOT NULL PRIMARY KEY, [body] XML);
CREATE SELECTIVE XML INDEX [sxi_docs] ON [dbo].[docs]([body]) FOR (p = '/b');
`)
	require.NoError(t, err)
	xmlDiff := Diff(fromXML, toXML)
	require.False(t, xmlDiff.IsEmpty())
	def, err := toXML.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Equal(t, 1, strings.Count(def, "CREATE SELECTIVE XML INDEX [sxi_docs]"))
	require.Contains(t, def, "FOR (p = '/b')")
	plan := GenerateMigration(fromXML, toXML, xmlDiff)
	mig := plan.SQL()
	require.Contains(t, mig, "CREATE SELECTIVE XML INDEX [sxi_docs]")
	require.Contains(t, mig, "FOR (p = '/b')")

	// 24. ALTER SEQUENCE preserves the effective creation-time START WITH when bounds change.
	c, err = LoadSQL(`CREATE SEQUENCE [dbo].[s] AS INT MINVALUE 10 MAXVALUE 100;`)
	require.NoError(t, err)
	execExpectOK(t, c, `ALTER SEQUENCE [dbo].[s] MINVALUE 5`)
	require.Equal(t, "10", c.GetSchema("dbo").Sequences[0].Start)
	require.Equal(t, "5", c.GetSchema("dbo").Sequences[0].MinValue)

	// 25. SYSTEM_VERSIONING = OFF succeeds when versioning was enabled at CREATE TABLE or via SET ON.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] (
    [id] INT NOT NULL PRIMARY KEY,
    [sys_start] DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
    [sys_end] DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
    PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])
) WITH (SYSTEM_VERSIONING = ON);
`)
	require.NoError(t, err)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = OFF)`)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = OFF)`, ErrSystemVersioningNotEnabled)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = ON)`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = OFF)`)

	// 26. Inbound FK check resolves exact case-distinct target table identities.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[Foo] ([id] INT NOT NULL PRIMARY KEY);
CREATE TABLE [dbo].[foo] ([id] INT NOT NULL PRIMARY KEY);
CREATE TABLE [dbo].[child] ([fid] INT NULL REFERENCES [dbo].[Foo]([id]));
`)
	require.NoError(t, err)
	execExpectCode(t, c, `TRUNCATE TABLE [dbo].[Foo]`, ErrTruncateReferencedByFK)
	execExpectCode(t, c, `DROP TABLE [dbo].[Foo]`, ErrDropTableReferencedByFK)
	execExpectOK(t, c, `TRUNCATE TABLE [dbo].[foo]`)
	execExpectOK(t, c, `DROP TABLE [dbo].[foo]`)

	// 27. DROP TABLE rejects tables referenced by WITH SCHEMABINDING views or functions.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] ([a] INT NOT NULL);
GO
CREATE VIEW [dbo].[v_sb] WITH SCHEMABINDING AS SELECT [a] FROM [dbo].[t];
`)
	require.NoError(t, err)
	execExpectCode(t, c, `DROP TABLE [dbo].[t]`, ErrDropSchemaNotEmpty)
	require.NotNil(t, c.GetTable("dbo", "t"))
	execExpectOK(t, c, `DROP VIEW [dbo].[v_sb]`)
	execExpectOK(t, c, `DROP TABLE [dbo].[t]`)

	// 28. ALTER SCHEMA TRANSFER updates only FKs pointing to the transferred case-distinct table.
	c, err = LoadSQL(`
CREATE SCHEMA [app];
GO
CREATE TABLE [dbo].[Foo] ([id] INT NOT NULL PRIMARY KEY);
CREATE TABLE [dbo].[foo] ([id] INT NOT NULL PRIMARY KEY);
CREATE TABLE [dbo].[child] ([fid] INT NULL REFERENCES [dbo].[Foo]([id]));
`)
	require.NoError(t, err)
	execExpectOK(t, c, `ALTER SCHEMA [app] TRANSFER [dbo].[foo]`)
	require.Equal(t, "dbo", c.GetTable("dbo", "child").ForeignKeys[0].ReferencedSchema)
	execExpectCode(t, c, `DROP TABLE [dbo].[Foo]`, ErrDropTableReferencedByFK)

	// 29. ALTER COLUMN rejects incompatible type changes on indexed/PK columns.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([id] INT NOT NULL PRIMARY KEY, [name] VARCHAR(20) NOT NULL); CREATE INDEX [ix_name] ON [dbo].[t]([name]);`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [id] BIGINT NOT NULL`, ErrColumnDependentObject)
	require.Equal(t, "INT", columnByName(c.GetTable("dbo", "t"), "id").Type)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [name] VARCHAR(50) NOT NULL`)
	require.Equal(t, "VARCHAR(50)", columnByName(c.GetTable("dbo", "t"), "name").Type)

	// 30. LoadMetadata clears temporal tracking maps from previously executed DDL.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] (
    [id] INT NOT NULL PRIMARY KEY,
    [sys_start] DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
    [sys_end] DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
    PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])
) WITH (SYSTEM_VERSIONING = ON);
`)
	require.NoError(t, err)
	plain, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] INT NOT NULL PRIMARY KEY);`)
	require.NoError(t, err)
	require.NoError(t, c.LoadMetadata(t.Context(), plain.Metadata()))
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = ON)`, ErrSystemVersioningNoPeriod)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = OFF)`, ErrSystemVersioningNotEnabled)

	// 31. Dropping a primary XML index cascades to its dependent secondary XML indexes.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[docs] ([id] INT NOT NULL PRIMARY KEY, [body] XML);
CREATE PRIMARY XML INDEX [pxml] ON [dbo].[docs]([body]);
CREATE XML INDEX [sxml] ON [dbo].[docs]([body]) USING XML INDEX [pxml] FOR PATH;
`)
	require.NoError(t, err)
	execExpectOK(t, c, `DROP INDEX [pxml] ON [dbo].[docs]`)
	require.Nil(t, indexByName(c.GetTable("dbo", "docs"), "pxml"))
	require.Nil(t, indexByName(c.GetTable("dbo", "docs"), "sxml"))

	// 32. DROP PERIOD FOR SYSTEM_TIME succeeds when versioning is OFF, fails when ON or absent.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] (
    [id] INT NOT NULL PRIMARY KEY,
    [sys_start] DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
    [sys_end] DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
    PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])
) WITH (SYSTEM_VERSIONING = ON);
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] DROP PERIOD FOR SYSTEM_TIME`, ErrDropPeriodSystemVersioned)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = OFF)`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP PERIOD FOR SYSTEM_TIME`)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] DROP PERIOD FOR SYSTEM_TIME`, ErrDropPeriodNotDefined)

	// 33. Collation-only ALTER COLUMN on an indexed column is rejected.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([a] VARCHAR(20) NOT NULL); CREATE INDEX [ix_a] ON [dbo].[t]([a]);`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [a] VARCHAR(20) COLLATE Latin1_General_CS_AS NOT NULL`, ErrColumnDependentObject)
	require.Equal(t, "", columnByName(c.GetTable("dbo", "t"), "a").Collation)

	// 34. DROP COLUMN rejects dropping the table's last column (single and multi-column).
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([a] INT, [b] INT);`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [a], [b]`, ErrDropLastColumn)
	require.Len(t, c.GetTable("dbo", "t").Columns, 2)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [b]`)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [a]`, ErrDropLastColumn)
	require.Len(t, c.GetTable("dbo", "t").Columns, 1)

	// 35. Adding a second primary key via ADD COLUMN or ADD CONSTRAINT is rejected.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([id] INT NOT NULL PRIMARY KEY, [x] INT NOT NULL);`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ADD [other_id] INT NOT NULL PRIMARY KEY`, ErrMultiplePrimaryKeys)
	require.Nil(t, columnByName(c.GetTable("dbo", "t"), "other_id"))
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ADD CONSTRAINT [PK_t2] PRIMARY KEY ([x])`, ErrMultiplePrimaryKeys)

	// 36. ALTER SCHEMA TRANSFER renders transferred views/procedures/functions under destination schema.
	c, err = LoadSQL(`
CREATE SCHEMA [app];
GO
CREATE VIEW [dbo].[v] AS SELECT 1 AS [a];
GO
CREATE PROCEDURE [dbo].[p] AS BEGIN SELECT 1; END;
GO
CREATE FUNCTION [dbo].[f]() RETURNS INT AS BEGIN RETURN 1; END;
`)
	require.NoError(t, err)
	execExpectOK(t, c, `ALTER SCHEMA [app] TRANSFER [dbo].[v]`)
	execExpectOK(t, c, `ALTER SCHEMA [app] TRANSFER [dbo].[p]`)
	execExpectOK(t, c, `ALTER SCHEMA [app] TRANSFER [dbo].[f]`)
	def, err = c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "CREATE VIEW [app].[v]")
	require.Contains(t, def, "CREATE PROCEDURE [app].[p]")
	require.Contains(t, def, "CREATE FUNCTION [app].[f]")
	require.NotContains(t, def, "[dbo].[v]")
	require.NotContains(t, def, "[dbo].[p]")
	require.NotContains(t, def, "[dbo].[f]")

	// 37. DROP COLUMN rejects columns referenced by WITH SCHEMABINDING views or functions.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] ([a] INT NOT NULL, [b] INT NOT NULL);
GO
CREATE VIEW [dbo].[v_sb] WITH SCHEMABINDING AS SELECT [a] FROM [dbo].[t];
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [a]`, ErrColumnDependentObject)
	require.NotNil(t, columnByName(c.GetTable("dbo", "t"), "a"))
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [b]`)

	// 38. ALTER SCHEMA TRANSFER updates trigger target/schema definitions on transferred tables.
	c, err = LoadSQL(`
CREATE SCHEMA [app];
GO
CREATE TABLE [dbo].[t] ([id] INT NOT NULL);
GO
CREATE TRIGGER [dbo].[tr_t] ON [dbo].[t] AFTER INSERT AS BEGIN SELECT 1; END;
`)
	require.NoError(t, err)
	execExpectOK(t, c, `ALTER SCHEMA [app] TRANSFER [dbo].[t]`)
	def, err = c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.Contains(t, def, "CREATE TRIGGER [app].[tr_t] ON [app].[t]")
	require.NotContains(t, def, "ON [dbo].[t]")
	mig = GenerateMigration(New(), c, Diff(New(), c)).SQL()
	require.Contains(t, mig, "CREATE TRIGGER [app].[tr_t] ON [app].[t]")
	require.NotContains(t, mig, "ON [dbo].[t]")

	// 39. DROP TABLE honors the database qualifier instead of dropping local tables.
	primaryMeta, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] INT NOT NULL);`)
	require.NoError(t, err)
	otherMeta, err := LoadSQL(`CREATE TABLE [dbo].[t] ([id] INT NOT NULL);`)
	require.NoError(t, err)
	c = New()
	require.NoError(t, c.LoadMetadataInto(t.Context(), primaryMeta.Metadata(), "appdb"))
	require.NoError(t, c.LoadMetadataInto(t.Context(), otherMeta.Metadata(), "otherdb"))
	execExpectCode(t, c, `DROP TABLE [unknown_db].[dbo].[t]`, ErrCannotDropNotFound)
	require.NotNil(t, c.GetTable("dbo", "t"))
	execExpectOK(t, c, `DROP TABLE [otherdb].[dbo].[t]`)
	require.NotNil(t, c.GetTable("dbo", "t"))
	require.Nil(t, c.GetDatabaseTable("otherdb", "dbo", "t"))

	// 40. CREATE/ALTER SEQUENCE rejects zero increment without mutating existing metadata.
	c, err = LoadSQL(`CREATE SEQUENCE [dbo].[s] START WITH 1 INCREMENT BY 1;`)
	require.NoError(t, err)
	execExpectCode(t, c, `CREATE SEQUENCE [dbo].[s0] INCREMENT BY 0`, ErrSequenceZeroIncrement)
	execExpectCode(t, c, `ALTER SEQUENCE [dbo].[s] INCREMENT BY 0 MINVALUE 5`, ErrSequenceZeroIncrement)
	require.Equal(t, "1", c.GetSchema("dbo").Sequences[0].Increment)
	require.Equal(t, "", c.GetSchema("dbo").Sequences[0].MinValue)

	// 41. ALTER TABLE ADD COLUMN rejects adding a second IDENTITY column.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([id] INT IDENTITY(1,1), [v] INT);`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ADD [id2] INT IDENTITY(1,1)`, ErrMultipleIdentityColumns)
	require.Nil(t, columnByName(c.GetTable("dbo", "t"), "id2"))

	// 42. CREATE INDEX WITH (DROP_EXISTING = ON) validates constraint-backed indexes before replacing them.
	c, err = LoadSQL(`CREATE TABLE [dbo].[p] ([a] INT NOT NULL, [b] INT NOT NULL, CONSTRAINT [PK_p] PRIMARY KEY CLUSTERED ([a]));`)
	require.NoError(t, err)
	execExpectCode(t, c, `CREATE UNIQUE CLUSTERED INDEX [PK_p] ON [dbo].[p] ([b]) WITH (DROP_EXISTING = ON)`, ErrDropExistingConstraintMismatch)
	require.Equal(t, []string{"a"}, indexByName(c.GetTable("dbo", "p"), "PK_p").Expressions)
	execExpectCode(t, c, `CREATE CLUSTERED INDEX [PK_p] ON [dbo].[p] ([a]) WITH (DROP_EXISTING = ON)`, ErrDropExistingConstraintMismatch)
	require.True(t, indexByName(c.GetTable("dbo", "p"), "PK_p").Unique)
	execExpectOK(t, c, `CREATE UNIQUE CLUSTERED INDEX [PK_p] ON [dbo].[p] ([a] DESC) WITH (DROP_EXISTING = ON)`)
	require.Equal(t, []string{"a"}, indexByName(c.GetTable("dbo", "p"), "PK_p").Expressions)
	require.Equal(t, []bool{true}, indexByName(c.GetTable("dbo", "p"), "PK_p").Descending)

	// 43. DROP VIEW rejects dropping a view referenced by a WITH SCHEMABINDING view.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] ([a] INT NOT NULL);
GO
CREATE VIEW [dbo].[v1] WITH SCHEMABINDING AS SELECT [a] FROM [dbo].[t];
GO
CREATE VIEW [dbo].[v2] WITH SCHEMABINDING AS SELECT [a] FROM [dbo].[v1];
`)
	require.NoError(t, err)
	execExpectCode(t, c, `DROP VIEW [dbo].[v1]`, ErrDropSchemaNotEmpty)
	require.NotNil(t, c.GetView("dbo", "v1"))
	execExpectOK(t, c, `DROP VIEW [dbo].[v2], [dbo].[v1]`)
	require.Nil(t, c.GetView("dbo", "v1"))
	require.Nil(t, c.GetView("dbo", "v2"))

	// 44. ALTER SCHEMA TRANSFER checks transferred constraint names for destination schema collisions.
	c, err = LoadSQL(`
CREATE SCHEMA [app];
GO
CREATE TABLE [dbo].[t] ([a] INT, CONSTRAINT [ck] CHECK ([a] > 0));
GO
CREATE TABLE [app].[u] ([a] INT, CONSTRAINT [ck] CHECK ([a] > 0));
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER SCHEMA [app] TRANSFER [dbo].[t]`, ErrDupObject)
	require.NotNil(t, c.GetTable("dbo", "t"))
	require.Nil(t, c.GetTable("app", "t"))

	// 45. DROP SCHEMA rejects dropping built-in schemas like dbo, even with IF EXISTS.
	c, err = LoadSQL(``)
	require.NoError(t, err)
	execExpectCode(t, c, `DROP SCHEMA [dbo]`, ErrSecurableNotFound)
	require.NotNil(t, c.GetSchema("dbo"))
	execExpectCode(t, c, `DROP SCHEMA IF EXISTS [dbo]`, ErrSecurableNotFound)
	require.NotNil(t, c.GetSchema("dbo"))

	// 46. Schema-bound column dependency resolves qualifiers/aliases to their owning table.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] ([id] INT NOT NULL, [a] INT NOT NULL, [b] INT NOT NULL);
GO
CREATE TABLE [dbo].[u] ([id] INT NOT NULL, [b] INT NOT NULL);
GO
CREATE VIEW [dbo].[v_join] WITH SCHEMABINDING AS
SELECT [t].[a], [u].[b]
FROM [dbo].[t] AS [t]
JOIN [dbo].[u] AS [u] ON [t].[id] = [u].[id];
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [a]`, ErrColumnDependentObject)
	execExpectCode(t, c, `ALTER TABLE [dbo].[u] DROP COLUMN [b]`, ErrColumnDependentObject)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [b] BIGINT NOT NULL`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [b]`)
	require.Nil(t, columnByName(c.GetTable("dbo", "t"), "b"))

	// 47. DROP FUNCTION rejects dropping a function referenced by a WITH SCHEMABINDING view.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] ([a] INT NOT NULL);
GO
CREATE FUNCTION [dbo].[f]() RETURNS INT WITH SCHEMABINDING AS BEGIN RETURN 1; END;
GO
CREATE VIEW [dbo].[v_fn] WITH SCHEMABINDING AS SELECT [dbo].[f]() AS [x] FROM [dbo].[t];
`)
	require.NoError(t, err)
	execExpectCode(t, c, `DROP FUNCTION [dbo].[f]`, ErrDropSchemaNotEmpty)
	require.Len(t, c.GetSchema("dbo").Functions, 1)
	execExpectOK(t, c, `DROP VIEW [dbo].[v_fn]`)
	execExpectOK(t, c, `DROP FUNCTION [dbo].[f]`)
	require.Empty(t, c.GetSchema("dbo").Functions)

	// 48. ALTER TABLE ADD PERIOD FOR SYSTEM_TIME accepts valid NOT NULL datetime2 period columns.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] (
    [id] INT NOT NULL PRIMARY KEY,
    [sys_start] DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
    [sys_end] DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
    PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])
) WITH (SYSTEM_VERSIONING = ON);
`)
	require.NoError(t, err)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = OFF)`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP PERIOD FOR SYSTEM_TIME`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ADD PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])`)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ADD PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])`, ErrPeriodAlreadyDefined)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = ON)`)

	// 49. ALTER TABLE ALTER COLUMN rejects altering a computed column.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([a] INT, [b] AS ([a] + 1));`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [b] INT`, ErrAlterComputedColumn)
	require.Equal(t, "AS ([a] + 1)", columnByName(c.GetTable("dbo", "t"), "b").Type)

	// 50. ALTER TABLE honors the database qualifier instead of altering local tables.
	primaryMeta, err = LoadSQL(`CREATE TABLE [dbo].[t] ([id] INT NOT NULL);`)
	require.NoError(t, err)
	otherMeta, err = LoadSQL(`CREATE TABLE [dbo].[t] ([id] INT NOT NULL);`)
	require.NoError(t, err)
	c = New()
	require.NoError(t, c.LoadMetadataInto(t.Context(), primaryMeta.Metadata(), "appdb"))
	require.NoError(t, c.LoadMetadataInto(t.Context(), otherMeta.Metadata(), "otherdb"))
	execExpectCode(t, c, `ALTER TABLE [unknown_db].[dbo].[t] ADD [x] INT`, ErrAlterTableNotFound)
	require.Nil(t, columnByName(c.GetTable("dbo", "t"), "x"))
	execExpectOK(t, c, `ALTER TABLE [otherdb].[dbo].[t] ADD [x] INT`)
	require.Nil(t, columnByName(c.GetTable("dbo", "t"), "x"))
	require.NotNil(t, columnByName(c.GetDatabaseTable("otherdb", "dbo", "t"), "x"))

	// 51. SYSTEM_TIME period columns cannot be dropped while the period exists.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] (
    [id] INT NOT NULL PRIMARY KEY,
    [sys_start] DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
    [sys_end] DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
    PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])
);
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [sys_start]`, ErrDropPeriodColumn)
	require.NotNil(t, columnByName(c.GetTable("dbo", "t"), "sys_start"))
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP PERIOD FOR SYSTEM_TIME`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP COLUMN [sys_start]`)
	require.Nil(t, columnByName(c.GetTable("dbo", "t"), "sys_start"))

	// 52. DROP SEQUENCE rejects dropping a sequence referenced by a column DEFAULT.
	c, err = LoadSQL(`
CREATE SEQUENCE [dbo].[s] AS INT START WITH 1;
GO
CREATE TABLE [dbo].[t] ([id] INT CONSTRAINT [DF_t_id] DEFAULT (NEXT VALUE FOR [dbo].[s]));
`)
	require.NoError(t, err)
	execExpectCode(t, c, `DROP SEQUENCE [dbo].[s]`, ErrDropSchemaNotEmpty)
	require.Len(t, c.GetSchema("dbo").Sequences, 1)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP CONSTRAINT [DF_t_id]`)
	execExpectOK(t, c, `DROP SEQUENCE [dbo].[s]`)
	require.Empty(t, c.GetSchema("dbo").Sequences)

	// 53. DROP CONSTRAINT rejects dropping a clustered primary key when an XML index exists.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] ([id] INT NOT NULL, [doc] XML, CONSTRAINT [PK_t] PRIMARY KEY CLUSTERED ([id]));
GO
CREATE PRIMARY XML INDEX [PXML_t_doc] ON [dbo].[t]([doc]);
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] DROP CONSTRAINT [PK_t]`, ErrDropPKWithXMLIndex)
	require.NotNil(t, indexByName(c.GetTable("dbo", "t"), "PK_t"))
	execExpectOK(t, c, `DROP INDEX [PXML_t_doc] ON [dbo].[t]`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP CONSTRAINT [PK_t]`)
	require.Nil(t, indexByName(c.GetTable("dbo", "t"), "PK_t"))

	// 54. ALTER SCHEMA TRANSFER rejects transferring objects referenced by WITH SCHEMABINDING consumers.
	c, err = LoadSQL(`
CREATE SCHEMA [app];
GO
CREATE TABLE [dbo].[t] ([a] INT NOT NULL);
GO
CREATE VIEW [dbo].[v] WITH SCHEMABINDING AS SELECT [a] FROM [dbo].[t];
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER SCHEMA [app] TRANSFER [dbo].[t]`, ErrTransferSchemaBound)
	require.NotNil(t, c.GetTable("dbo", "t"))
	require.Nil(t, c.GetTable("app", "t"))
	execExpectCode(t, c, `ALTER SCHEMA [app] TRANSFER [dbo].[v]`, ErrTransferSchemaBound)
	require.NotNil(t, c.GetView("dbo", "v"))

	// 55. Multi-target DROP INDEX validates FK dependencies against all staged index drops.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[p] ([a] INT NOT NULL);
GO
CREATE UNIQUE INDEX [u1] ON [dbo].[p] ([a]);
GO
CREATE UNIQUE INDEX [u2] ON [dbo].[p] ([a]);
GO
CREATE TABLE [dbo].[c] ([a] INT REFERENCES [dbo].[p]([a]));
`)
	require.NoError(t, err)
	execExpectCode(t, c, `DROP INDEX [u1] ON [dbo].[p], [u2] ON [dbo].[p]`, ErrDropIndexReferencedByFK)
	require.NotNil(t, indexByName(c.GetTable("dbo", "p"), "u1"))
	require.NotNil(t, indexByName(c.GetTable("dbo", "p"), "u2"))
	execExpectOK(t, c, `DROP INDEX [u1] ON [dbo].[p]`)
	require.Nil(t, indexByName(c.GetTable("dbo", "p"), "u1"))
	require.NotNil(t, indexByName(c.GetTable("dbo", "p"), "u2"))

	// 56. ALTER SEQUENCE validates combined min/max/restart bounds before mutating metadata.
	c, err = LoadSQL(`CREATE SEQUENCE [dbo].[s] START WITH 1 MINVALUE 1 MAXVALUE 10;`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER SEQUENCE [dbo].[s] MINVALUE 20`, ErrSequenceMinGreaterThanMax)
	require.Equal(t, "1", c.GetSchema("dbo").Sequences[0].MinValue)
	require.Equal(t, "10", c.GetSchema("dbo").Sequences[0].MaxValue)
	execExpectCode(t, c, `ALTER SEQUENCE [dbo].[s] RESTART WITH 20`, ErrSequenceStartOutOfRange)
	execExpectCode(t, c, `CREATE SEQUENCE [dbo].[bad_seq] MINVALUE 20 MAXVALUE 10`, ErrSequenceMinGreaterThanMax)

	// 57. DROP TABLE rejects dropping a table while SYSTEM_VERSIONING is ON.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] (
    [id] INT NOT NULL PRIMARY KEY,
    [sys_start] DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
    [sys_end] DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
    PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])
) WITH (SYSTEM_VERSIONING = ON);
`)
	require.NoError(t, err)
	execExpectCode(t, c, `DROP TABLE [dbo].[t]`, ErrDropSystemVersionedTable)
	require.NotNil(t, c.GetTable("dbo", "t"))
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = OFF)`)
	execExpectOK(t, c, `DROP TABLE [dbo].[t]`)
	require.Nil(t, c.GetTable("dbo", "t"))

	// 58. DROP SCHEMA removes the exact resolved schema in case-sensitive metadata snapshots.
	c = New()
	require.NoError(t, c.LoadMetadata(t.Context(), &metadata.DatabaseSchemaMetadata{
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "Foo",
				Tables: []*metadata.TableMetadata{
					{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int", Position: 1}}},
				},
			},
			{Name: "foo"},
		},
	}))
	execExpectOK(t, c, `DROP SCHEMA [foo]`)
	require.NotNil(t, c.GetSchema("Foo"))
	require.NotNil(t, c.GetTable("Foo", "t"))
	require.Len(t, c.ListSchemaNames(), 1)

	// 59. ALTER COLUMN rejects identity-incompatible types (and NULL) on IDENTITY columns.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([id] INT IDENTITY(1,1) NOT NULL, [v] INT);`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [id] VARCHAR(20) NOT NULL`, ErrIdentityColumnType)
	require.Equal(t, "INT", columnByName(c.GetTable("dbo", "t"), "id").Type)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [id] DECIMAL(10,2) NOT NULL`, ErrIdentityColumnType)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [id] BIGINT NULL`, ErrIdentityColumnType)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [id] BIGINT NOT NULL`)
	require.Equal(t, "BIGINT", columnByName(c.GetTable("dbo", "t"), "id").Type)

	// 60. DROP VIEW removes only the exact requested case-distinct object.
	c = New()
	require.NoError(t, c.LoadMetadata(t.Context(), &metadata.DatabaseSchemaMetadata{
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "dbo",
				Views: []*metadata.ViewMetadata{
					{Name: "Foo", Definition: "CREATE VIEW [dbo].[Foo] AS SELECT 1 AS a;"},
					{Name: "foo", Definition: "CREATE VIEW [dbo].[foo] AS SELECT 2 AS a;"},
				},
			},
		},
	}))
	execExpectOK(t, c, `DROP VIEW [dbo].[foo]`)
	require.Len(t, c.GetSchema("dbo").Views, 1)
	require.Equal(t, "Foo", c.GetSchema("dbo").Views[0].Name)

	// 61. TRUNCATE TABLE rejects truncating a table while SYSTEM_VERSIONING is ON.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] (
    [id] INT NOT NULL PRIMARY KEY,
    [sys_start] DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
    [sys_end] DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
    PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])
) WITH (SYSTEM_VERSIONING = ON);
`)
	require.NoError(t, err)
	execExpectCode(t, c, `TRUNCATE TABLE [dbo].[t]`, ErrTruncateSystemVersionedTable)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = OFF)`)
	execExpectOK(t, c, `TRUNCATE TABLE [dbo].[t]`)

	// 62. ALTER SCHEMA TRANSFER resolves the exact object name in case-sensitive snapshots.
	c = New()
	require.NoError(t, c.LoadMetadata(t.Context(), &metadata.DatabaseSchemaMetadata{
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "dbo",
				Views: []*metadata.ViewMetadata{
					{Name: "Foo", Definition: "CREATE VIEW [dbo].[Foo] AS SELECT 1 AS a;"},
					{Name: "foo", Definition: "CREATE VIEW [dbo].[foo] AS SELECT 2 AS a;"},
				},
			},
			{Name: "app"},
		},
	}))
	execExpectOK(t, c, `ALTER SCHEMA [app] TRANSFER [dbo].[foo]`)
	require.Len(t, c.GetSchema("dbo").Views, 1)
	require.Equal(t, "Foo", c.GetSchema("dbo").Views[0].Name)
	require.Len(t, c.GetSchema("app").Views, 1)
	require.Equal(t, "foo", c.GetSchema("app").Views[0].Name)

	// 63. ALTER VIEW and CREATE OR ALTER VIEW reject replacing views referenced by schema-bound consumers.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] ([a] INT NOT NULL, [b] INT NOT NULL);
GO
CREATE VIEW [dbo].[v1] WITH SCHEMABINDING AS SELECT [a], [b] FROM [dbo].[t];
GO
CREATE VIEW [dbo].[v2] WITH SCHEMABINDING AS SELECT [a] FROM [dbo].[v1];
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER VIEW [dbo].[v1] WITH SCHEMABINDING AS SELECT [b] FROM [dbo].[t]`, ErrDropSchemaNotEmpty)
	execExpectCode(t, c, `CREATE OR ALTER VIEW [dbo].[v1] WITH SCHEMABINDING AS SELECT [b] FROM [dbo].[t]`, ErrDropSchemaNotEmpty)
	require.Contains(t, c.GetView("dbo", "v1").Definition, "[a]")

	// 64. ALTER COLUMN defaults nullability to NULL when omitted (and rejects PK columns accordingly).
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([id] INT NOT NULL CONSTRAINT [PK_t] PRIMARY KEY, [a] INT NOT NULL);`)
	require.NoError(t, err)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [a] BIGINT`)
	require.Equal(t, "BIGINT", columnByName(c.GetTable("dbo", "t"), "a").Type)
	require.True(t, columnByName(c.GetTable("dbo", "t"), "a").Nullable)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [id] BIGINT`, ErrColumnDependentObject)
	require.Equal(t, "INT", columnByName(c.GetTable("dbo", "t"), "id").Type)
	require.False(t, columnByName(c.GetTable("dbo", "t"), "id").Nullable)

	// 65. ALTER TABLE ADD COLUMN / ADD CONSTRAINT rejects adding a second clustered index.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] ([a] INT NOT NULL);
GO
CREATE CLUSTERED INDEX [ix] ON [dbo].[t]([a]);
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ADD [id] INT NOT NULL CONSTRAINT [pk] PRIMARY KEY CLUSTERED`, ErrMultipleClusteredIndex)
	require.Nil(t, columnByName(c.GetTable("dbo", "t"), "id"))
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ADD CONSTRAINT [pk] PRIMARY KEY CLUSTERED ([a])`, ErrMultipleClusteredIndex)
	require.Nil(t, indexByName(c.GetTable("dbo", "t"), "pk"))
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ADD [id] INT NOT NULL CONSTRAINT [pk] PRIMARY KEY NONCLUSTERED`)
	require.NotNil(t, indexByName(c.GetTable("dbo", "t"), "pk"))

	// 66. DROP INDEX resolves the exact index name in case-sensitive snapshots.
	c = New()
	require.NoError(t, c.LoadMetadata(t.Context(), &metadata.DatabaseSchemaMetadata{
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "dbo",
				Tables: []*metadata.TableMetadata{
					{
						Name:    "t",
						Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int", Position: 1}},
						Indexes: []*metadata.IndexMetadata{
							{Name: "Foo", Type: "NONCLUSTERED", Expressions: []string{"a"}},
							{Name: "foo", Type: "NONCLUSTERED", Expressions: []string{"a"}},
						},
					},
				},
			},
		},
	}))
	execExpectOK(t, c, `DROP INDEX [foo] ON [dbo].[t]`)
	require.Len(t, c.GetTable("dbo", "t").Indexes, 1)
	require.Equal(t, "Foo", c.GetTable("dbo", "t").Indexes[0].Name)

	// 67. ALTER COLUMN rejects changing the data type of a column with a bound DEFAULT constraint.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([a] INT CONSTRAINT [DF_t_a] DEFAULT (0));`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [a] BIGINT`, ErrColumnDependentObject)
	require.Equal(t, "INT", columnByName(c.GetTable("dbo", "t"), "a").Type)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP CONSTRAINT [DF_t_a]`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [a] BIGINT`)
	require.Equal(t, "BIGINT", columnByName(c.GetTable("dbo", "t"), "a").Type)

	// 68. DROP CONSTRAINT resolves the exact constraint name in case-sensitive snapshots.
	c = New()
	require.NoError(t, c.LoadMetadata(t.Context(), &metadata.DatabaseSchemaMetadata{
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "dbo",
				Tables: []*metadata.TableMetadata{
					{
						Name: "t",
						Columns: []*metadata.ColumnMetadata{
							{Name: "a", Type: "int", Position: 1},
							{Name: "b", Type: "int", Position: 2},
						},
						Indexes: []*metadata.IndexMetadata{
							{Name: "Foo", Unique: true, IsConstraint: true, Type: "NONCLUSTERED", Expressions: []string{"a"}},
							{Name: "foo", Unique: true, IsConstraint: true, Type: "NONCLUSTERED", Expressions: []string{"b"}},
						},
					},
				},
			},
		},
	}))
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP CONSTRAINT [foo]`)
	require.Len(t, c.GetTable("dbo", "t").Indexes, 1)
	require.Equal(t, "Foo", c.GetTable("dbo", "t").Indexes[0].Name)

	// 69. ALTER SEQUENCE resolves the exact sequence name in case-sensitive snapshots.
	c = New()
	require.NoError(t, c.LoadMetadata(t.Context(), &metadata.DatabaseSchemaMetadata{
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "dbo",
				Sequences: []*metadata.SequenceMetadata{
					{Name: "Foo", DataType: "int", Start: "1", Increment: "1"},
					{Name: "foo", DataType: "int", Start: "1", Increment: "1"},
				},
			},
		},
	}))
	execExpectOK(t, c, `ALTER SEQUENCE [dbo].[foo] INCREMENT BY 5`)
	require.Equal(t, "1", c.GetSchema("dbo").Sequences[0].Increment)
	require.Equal(t, "5", c.GetSchema("dbo").Sequences[1].Increment)

	// 70. ALTER COLUMN ADD/DROP PERSISTED updates computed-column metadata.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([a] INT, [b] AS ([a] + 1));`)
	require.NoError(t, err)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [b] ADD PERSISTED`)
	require.Equal(t, "AS ([a] + 1) PERSISTED", columnByName(c.GetTable("dbo", "t"), "b").Type)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [b] DROP PERSISTED`)
	require.Equal(t, "AS ([a] + 1)", columnByName(c.GetTable("dbo", "t"), "b").Type)

	// 71. lookupSchemaObject resolves exact names across object kinds before case-folded fallback.
	c = New()
	require.NoError(t, c.LoadMetadata(t.Context(), &metadata.DatabaseSchemaMetadata{
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "dbo",
				Tables: []*metadata.TableMetadata{
					{Name: "Foo", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int", Position: 1}}},
				},
				Views: []*metadata.ViewMetadata{
					{Name: "foo", Definition: "CREATE VIEW [dbo].[foo] AS SELECT 1 AS a;"},
				},
			},
		},
	}))
	execExpectOK(t, c, `ALTER VIEW [dbo].[foo] AS SELECT 2 AS a`)
	require.Contains(t, c.GetSchema("dbo").Views[0].Definition, "SELECT 2 AS a")
	execExpectOK(t, c, `DROP VIEW [dbo].[foo]`)
	require.Empty(t, c.GetSchema("dbo").Views)
	require.NotNil(t, c.GetTable("dbo", "Foo"))

	// 72. ALTER TABLE ADD COLUMN validates inline foreign-key targets (table, columns, candidate key).
	c, err = LoadSQL(`
CREATE TABLE [dbo].[p] ([id] INT NOT NULL PRIMARY KEY, [val] INT);
GO
CREATE TABLE [dbo].[t] ([a] INT);
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ADD [pid] INT REFERENCES [dbo].[missing]([id])`, ErrFKInvalidRefTable)
	require.Nil(t, columnByName(c.GetTable("dbo", "t"), "pid"))
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ADD [pid] INT REFERENCES [dbo].[p]([val])`, ErrFKNoCandidateKey)
	require.Nil(t, columnByName(c.GetTable("dbo", "t"), "pid"))
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ADD [pid] INT REFERENCES [dbo].[p]([id])`)
	require.NotNil(t, columnByName(c.GetTable("dbo", "t"), "pid"))
	require.Len(t, c.GetTable("dbo", "t").ForeignKeys, 1)

	// 73. ALTER FUNCTION / CREATE OR ALTER FUNCTION rejects changing a function's category.
	c, err = LoadSQL(`CREATE FUNCTION [dbo].[f]() RETURNS INT AS BEGIN RETURN 1; END;`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER FUNCTION [dbo].[f]() RETURNS TABLE AS RETURN (SELECT 1 AS [x])`, ErrIncompatibleAlterType)
	execExpectCode(t, c, `CREATE OR ALTER FUNCTION [dbo].[f]() RETURNS @t TABLE ([x] INT) AS BEGIN RETURN; END`, ErrIncompatibleAlterType)
	require.Contains(t, c.GetSchema("dbo").Functions[0].Definition, "RETURNS INT")
	execExpectOK(t, c, `ALTER FUNCTION [dbo].[f]() RETURNS INT AS BEGIN RETURN 2; END`)
	require.Contains(t, c.GetSchema("dbo").Functions[0].Definition, "RETURN 2")

	// 74. ALTER COLUMN clears collation when changing from a character type to a noncharacter type.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([a] VARCHAR(10) COLLATE Latin1_General_100_BIN2);`)
	require.NoError(t, err)
	require.Equal(t, "Latin1_General_100_BIN2", columnByName(c.GetTable("dbo", "t"), "a").Collation)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [a] INT`)
	require.Equal(t, "INT", columnByName(c.GetTable("dbo", "t"), "a").Type)
	require.Empty(t, columnByName(c.GetTable("dbo", "t"), "a").Collation)
	def, err = c.Definition(DefinitionOptions{})
	require.NoError(t, err)
	require.NotContains(t, def, "COLLATE")

	// 75. ALTER TABLE ADD rejects explicitly nullable IDENTITY columns.
	c, err = LoadSQL(`CREATE TABLE [dbo].[t] ([a] INT);`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ADD [id] INT IDENTITY(1,1) NULL`, ErrIdentityColumnType)
	require.Nil(t, columnByName(c.GetTable("dbo", "t"), "id"))
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ADD [id] INT IDENTITY(1,1) NOT NULL`)
	require.NotNil(t, columnByName(c.GetTable("dbo", "t"), "id"))

	// 76. Filtered unique indexes do not serve as surviving FK fallback keys during DROP INDEX.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[p] ([a] INT NOT NULL);
GO
CREATE UNIQUE INDEX [u1] ON [dbo].[p] ([a]);
GO
CREATE UNIQUE INDEX [u2] ON [dbo].[p] ([a]) WHERE [a] > 0;
GO
CREATE TABLE [dbo].[c] ([a] INT REFERENCES [dbo].[p]([a]));
`)
	require.NoError(t, err)
	execExpectCode(t, c, `DROP INDEX [u1] ON [dbo].[p]`, ErrDropIndexReferencedByFK)
	require.NotNil(t, indexByName(c.GetTable("dbo", "p"), "u1"))
	execExpectOK(t, c, `DROP INDEX [u2] ON [dbo].[p]`)
	require.Nil(t, indexByName(c.GetTable("dbo", "p"), "u2"))

	// 77. Dropping a temporal table in a mixed-case schema clears its SYSTEM_TIME state.
	c, err = LoadSQL(`
CREATE SCHEMA [App];
GO
CREATE TABLE [App].[t] (
    [id] INT NOT NULL PRIMARY KEY,
    [sys_start] DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
    [sys_end] DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
    PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])
);
`)
	require.NoError(t, err)
	execExpectOK(t, c, `DROP TABLE [App].[t]`)
	execExpectOK(t, c, `CREATE TABLE [App].[t] ([id] INT NOT NULL PRIMARY KEY)`)
	execExpectCode(t, c, `ALTER TABLE [App].[t] SET (SYSTEM_VERSIONING = ON)`, ErrSystemVersioningNoPeriod)

	// 78. ALTER SEQUENCE honors database qualifiers instead of mutating local sequences.
	primaryMeta, err = LoadSQL(`CREATE SEQUENCE [dbo].[s] START WITH 1 INCREMENT BY 1;`)
	require.NoError(t, err)
	otherMeta, err = LoadSQL(`CREATE SEQUENCE [dbo].[s] START WITH 1 INCREMENT BY 1;`)
	require.NoError(t, err)
	c = New()
	require.NoError(t, c.LoadMetadataInto(t.Context(), primaryMeta.Metadata(), "appdb"))
	require.NoError(t, c.LoadMetadataInto(t.Context(), otherMeta.Metadata(), "otherdb"))
	execExpectCode(t, c, `ALTER SEQUENCE [missingdb].[dbo].[s] INCREMENT BY 5`, ErrSecurableNotFound)
	require.Equal(t, "1", c.GetSchema("dbo").Sequences[0].Increment)
	execExpectOK(t, c, `ALTER SEQUENCE [otherdb].[dbo].[s] INCREMENT BY 5`)
	require.Equal(t, "1", c.GetSchema("dbo").Sequences[0].Increment)
	require.Equal(t, "5", c.GetDatabaseSchema("otherdb", "dbo").Sequences[0].Increment)

	// 79. ALTER SCHEMA TRANSFER moves objects between case-distinct schemas Foo and foo.
	c = New()
	require.NoError(t, c.LoadMetadata(t.Context(), &metadata.DatabaseSchemaMetadata{
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "Foo",
				Tables: []*metadata.TableMetadata{
					{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int", Position: 1}}},
				},
			},
			{Name: "foo"},
		},
	}))
	execExpectOK(t, c, `ALTER SCHEMA [foo] TRANSFER [Foo].[t]`)
	require.Nil(t, c.GetTable("Foo", "t"))
	require.NotNil(t, c.GetTable("foo", "t"))

	// 80. Sequence bounds are validated for decimal/numeric types based on declared precision.
	c, err = LoadSQL(`CREATE SEQUENCE [dbo].[s] AS DECIMAL(2,0) START WITH 1 MINVALUE 1 MAXVALUE 99;`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER SEQUENCE [dbo].[s] MAXVALUE 100`, ErrSequenceInvalidBounds)
	require.Equal(t, "99", c.GetSchema("dbo").Sequences[0].MaxValue)
	execExpectCode(t, c, `CREATE SEQUENCE [dbo].[s_bad] AS NUMERIC(2,0) MINVALUE -100 MAXVALUE 99`, ErrSequenceInvalidBounds)

	// 81. ALTER FUNCTION succeeds when referenced by a column DEFAULT or CHECK constraint (unlike DROP FUNCTION).
	c, err = LoadSQL(`
CREATE FUNCTION [dbo].[f]() RETURNS INT AS BEGIN RETURN 1; END;
GO
CREATE TABLE [dbo].[t] ([a] INT DEFAULT ([dbo].[f]()), CHECK ([dbo].[f]() > 0));
`)
	require.NoError(t, err)
	execExpectOK(t, c, `ALTER FUNCTION [dbo].[f]() RETURNS INT AS BEGIN RETURN 2; END`)
	require.Contains(t, c.GetSchema("dbo").Functions[0].Definition, "RETURN 2")
	execExpectOK(t, c, `CREATE OR ALTER FUNCTION [dbo].[f]() RETURNS INT AS BEGIN RETURN 3; END`)
	require.Contains(t, c.GetSchema("dbo").Functions[0].Definition, "RETURN 3")
	execExpectCode(t, c, `DROP FUNCTION [dbo].[f]`, ErrDropSchemaNotEmpty)

	// 82. DROP CONSTRAINT rejects dropping the primary key while SYSTEM_VERSIONING is ON.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] (
    [id] INT NOT NULL CONSTRAINT [PK_t] PRIMARY KEY,
    [sys_start] DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
    [sys_end] DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
    PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])
) WITH (SYSTEM_VERSIONING = ON);
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] DROP CONSTRAINT [PK_t]`, ErrDropPKSystemVersioned)
	require.NotNil(t, indexByName(c.GetTable("dbo", "t"), "PK_t"))
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = OFF)`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP CONSTRAINT [PK_t]`)
	require.Nil(t, indexByName(c.GetTable("dbo", "t"), "PK_t"))

	// 83. ALTER COLUMN checks schema-bound consumers for nullability-only changes.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] ([a] INT NOT NULL, [b] INT NOT NULL);
GO
CREATE VIEW [dbo].[v_sb] WITH SCHEMABINDING AS SELECT [a] FROM [dbo].[t];
`)
	require.NoError(t, err)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [a] INT NULL`, ErrColumnDependentObject)
	require.False(t, columnByName(c.GetTable("dbo", "t"), "a").Nullable)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ALTER COLUMN [b] INT NULL`)
	require.True(t, columnByName(c.GetTable("dbo", "t"), "b").Nullable)

	// 84. SET (SYSTEM_VERSIONING = ON) requires a primary key on the temporal table.
	c, err = LoadSQL(`
CREATE TABLE [dbo].[t] (
    [id] INT NOT NULL CONSTRAINT [PK_t] PRIMARY KEY,
    [sys_start] DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
    [sys_end] DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
    PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])
) WITH (SYSTEM_VERSIONING = ON);
`)
	require.NoError(t, err)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = OFF)`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] DROP CONSTRAINT [PK_t]`)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = ON)`, ErrDropPKSystemVersioned)
	execExpectCode(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = OFF)`, ErrSystemVersioningNotEnabled)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] ADD CONSTRAINT [PK_t] PRIMARY KEY ([id])`)
	execExpectOK(t, c, `ALTER TABLE [dbo].[t] SET (SYSTEM_VERSIONING = ON)`)
	execExpectCode(t, c, `CREATE TABLE [dbo].[t_nopk] (
    [id] INT NOT NULL,
    [sys_start] DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
    [sys_end] DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
    PERIOD FOR SYSTEM_TIME ([sys_start], [sys_end])
) WITH (SYSTEM_VERSIONING = ON)`, ErrDropPKSystemVersioned)
	require.Nil(t, c.GetTable("dbo", "t_nopk"))

	// 85. Sequence default dependency detection ignores string literals and comments.
	c, err = LoadSQL(`
CREATE SEQUENCE [dbo].[s] AS INT START WITH 1 INCREMENT BY 1;
GO
CREATE TABLE [dbo].[t] (
    [note] NVARCHAR(100) CONSTRAINT [DF_t_note] DEFAULT (N'NEXT VALUE FOR dbo.s'),
    [c] INT CONSTRAINT [DF_t_c] DEFAULT (/* NEXT VALUE FOR dbo.s */ 1),
    [l] INT CONSTRAINT [DF_t_l] DEFAULT (1 -- NEXT VALUE FOR dbo.s
    )
);
`)
	require.NoError(t, err)
	require.False(t, defaultReferencesSequence(`N'NEXT VALUE FOR dbo.s'`, "dbo", "s"))
	require.False(t, defaultReferencesSequence(`/* NEXT VALUE FOR dbo.s */ 1`, "dbo", "s"))
	require.True(t, defaultReferencesSequence(`(NEXT /* c */ VALUE FOR [dbo].[s])`, "dbo", "s"))
	execExpectOK(t, c, `DROP SEQUENCE [dbo].[s]`)
	require.Empty(t, c.GetSchema("dbo").Sequences)
}
