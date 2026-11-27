package analysis

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
	"github.com/stretchr/testify/require"
)

func parseScript(t *testing.T, script string) []omnipg.Statement {
	t.Helper()
	stmts, err := omnipg.Parse(script)
	require.NoError(t, err)
	return stmts
}

// searchPublicOrApp finds t and o in public, and only_app in app.
func searchPublicOrApp(searchPath []string, name string) string {
	for _, schema := range searchPath {
		if (schema == "public" && (name == "t" || name == "o")) || (schema == "app" && name == "only_app") {
			return schema
		}
	}
	return ""
}

func backupsOf(t *testing.T, script string, searchPath []string) ([]BackupStatement, error) {
	t.Helper()
	return TransformDMLToSelect(parseScript(t, script), script, BackupOptions{
		SearchPath:     searchPath,
		SearchRelation: searchPublicOrApp,
		TargetSchema:   "bbdataarchive",
		TablePrefix:    "_bb_",
	})
}

func TestTransformDMLToSelect(t *testing.T) {
	script := "UPDATE t SET a = 1 WHERE id = 1; DELETE FROM o AS x USING t WHERE x.id = t.id; UPDATE t SET b = 2"
	backups, err := backupsOf(t, script, nil)
	require.NoError(t, err)
	require.Equal(t, []BackupStatement{
		{
			Statement:       "CREATE TABLE \"bbdataarchive\".\"_bb__t_public\" AS\n  SELECT \"public\".\"t\".* FROM t WHERE id = 1\n  UNION\n  SELECT \"public\".\"t\".* FROM t;",
			SourceSchema:    "public",
			SourceTableName: "t",
			TargetTableName: "_bb__t_public",
			StartOffset:     0,
			EndOffset:       len(script),
		},
		{
			Statement:       "CREATE TABLE \"bbdataarchive\".\"_bb__o_public\" AS\n  SELECT \"x\".* FROM o AS x, t WHERE x.id = t.id;",
			SourceSchema:    "public",
			SourceTableName: "o",
			TargetTableName: "_bb__o_public",
			StartOffset:     33,
			EndOffset:       77,
		},
	}, backups)
}

func TestTransformDMLToSelectTruncatesNamesByBytes(t *testing.T) {
	// "_bb__" and 19 three-byte characters fill 62 of the 63 bytes.
	backups, err := backupsOf(t, "DELETE FROM public."+strings.Repeat("表", 30), nil)
	require.NoError(t, err)
	require.Equal(t, "_bb__"+strings.Repeat("表", 19), backups[0].TargetTableName)
}

func TestTransformDMLToSelectReportsTheFirstFailingTable(t *testing.T) {
	for range 50 {
		_, err := backupsOf(t, "UPDATE t SET a = 1; DELETE FROM t; UPDATE o SET a = 1; DELETE FROM o", nil)
		require.EqualError(t, err, `The statement type is not the same for all statements on the same table "public.t"`)
	}
}

func TestTransformDMLToSelectSearchPath(t *testing.T) {
	backups, err := backupsOf(t, "SET search_path TO app; DELETE FROM only_app", []string{"public"})
	require.NoError(t, err)
	require.Len(t, backups, 1)
	require.Equal(t, "app", backups[0].SourceSchema)

	_, err = backupsOf(t, "DELETE FROM only_app", nil)
	require.EqualError(t, err, `failed to prepare transformation: failed to extract table reference from delete statement at offset 0: Table "only_app" not found in metadata with search path [public]`)
}

func TestTransformDMLToSelectRefuses(t *testing.T) {
	_, err := backupsOf(t, "UPDATE t SET a = 1; DELETE FROM t", nil)
	require.EqualError(t, err, `The statement type is not the same for all statements on the same table "public.t"`)

	_, err = backupsOf(t, "DELETE FROM c.a_b; DELETE FROM b_c.a", nil)
	require.EqualError(t, err, `prior backup cannot back up both "b_c"."a" and "c"."a_b": their backup tables would both be named "_bb__a_b_c"`)

	// The tables' names differ in their last byte, which truncating the
	// backup's name to 63 bytes cuts off.
	long := strings.Repeat("x", 60)
	_, err = backupsOf(t, "DELETE FROM s."+long+"1; DELETE FROM s."+long+"2", nil)
	require.ErrorContains(t, err, "their backup tables would both be named")

	_, err = backupsOf(t, "WITH d AS (DELETE FROM o RETURNING *) UPDATE t SET a = 1", nil)
	require.EqualError(t, err, `failed to prepare transformation: prior backup does not support data-modifying WITH query "d"; rewrite the statement without it, or disable or skip prior backup`)
	require.Equal(t, []string{"d"}, DataModifyingCTENames(parseScript(t, "WITH d AS (DELETE FROM o RETURNING *), s AS (SELECT 1) SELECT 1")[0].AST))
}

func TestRestoreSQL(t *testing.T) {
	table := &metadata.TableMetadata{
		Name: "t",
		Columns: []*metadata.ColumnMetadata{
			{Name: "id"}, {Name: "a"}, {Name: "g", Generation: &metadata.GenerationMetadata{Expression: "id + 1"}},
		},
		Indexes: []*metadata.IndexMetadata{
			{Name: "t_pkey", Primary: true, Unique: true, Expressions: []string{"id"}},
			{Name: "t_a_key", Unique: true, Expressions: []string{"a"}},
		},
	}
	options := RestoreOptions{
		OriginalSchema: "public",
		OriginalTable:  "t",
		BackupSchema:   "bbdataarchive",
		BackupTable:    "_bb__t_public",
		Table:          table,
		PrimaryKey:     table.Indexes[0],
		OriginalSQL:    "DELETE FROM t",
	}
	script := "SET search_path = public; DELETE FROM t; UPDATE t SET a = 1; UPDATE t SET id = 2, a = 3; UPDATE t SET a.x = 1, a.y = 2"
	stmts := parseScript(t, script)
	dmls := BackedUpDMLs(stmts)
	require.Len(t, dmls, 4)
	require.Equal(t, "SET search_path = public", PrependStatements(stmts))

	got, err := RestoreSQL([]ast.Node{dmls[0].Node}, options)
	require.NoError(t, err)
	require.Equal(t, "/*\nOriginal SQL:\nDELETE FROM t\n*/\nINSERT INTO \"public\".\"t\" (\"id\", \"a\") SELECT \"id\", \"a\" FROM \"bbdataarchive\".\"_bb__t_public\";", got)

	options.PrependStatements = PrependStatements(stmts)
	got, err = RestoreSQL([]ast.Node{dmls[1].Node}, options)
	require.NoError(t, err)
	require.Equal(t, "SET search_path = public;\n/*\nOriginal SQL:\nDELETE FROM t\n*/\nINSERT INTO \"public\".\"t\" (\"id\", \"a\") SELECT \"id\", \"a\" FROM \"bbdataarchive\".\"_bb__t_public\" ON CONFLICT ON CONSTRAINT \"t_pkey\" DO UPDATE SET \"a\" = EXCLUDED.\"a\";", got)

	// An UPDATE of every unique key cannot be restored.
	_, err = RestoreSQL([]ast.Node{dmls[2].Node}, options)
	require.EqualError(t, err, "no disjoint unique key found for table")

	// The upsert also inserts the deleted rows again.
	got, err = RestoreSQL([]ast.Node{dmls[0].Node, dmls[1].Node}, options)
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got, ` ON CONFLICT ON CONSTRAINT "t_pkey" DO UPDATE SET "a" = EXCLUDED."a";`), got)

	// One UPDATE leaves the primary key alone, the other changes it.
	_, err = RestoreSQL([]ast.Node{dmls[1].Node, dmls[2].Node}, options)
	require.EqualError(t, err, "no disjoint unique key found for table")

	// Assigning fields of a composite column restores the column, once.
	got, err = RestoreSQL([]ast.Node{dmls[3].Node}, options)
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(got, ` ON CONFLICT ON CONSTRAINT "t_pkey" DO UPDATE SET "a" = EXCLUDED."a";`), got)
}

func TestRestoreSQLConflictTarget(t *testing.T) {
	table := func(indexes ...*metadata.IndexMetadata) *metadata.TableMetadata {
		return &metadata.TableMetadata{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "Id"}, {Name: "a"}, {Name: "b"}}, Indexes: indexes}
	}
	pk := &metadata.IndexMetadata{Name: "t_pkey", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{`"Id"`}}
	update := BackedUpDMLs(parseScript(t, `UPDATE t SET "Id" = 2, b = 3`))[0].Node
	restore := func(table *metadata.TableMetadata) (string, error) {
		return RestoreSQL([]ast.Node{update}, RestoreOptions{OriginalSchema: "public", OriginalTable: "t", BackupSchema: "bb", BackupTable: "b", Table: table, PrimaryKey: pk})
	}
	const plainIndex = "CREATE UNIQUE INDEX t_a ON public.t USING btree (a)"

	// The UPDATE changes the quoted primary key column, so a unique index on a
	// decides, matched by its column.
	got, err := restore(table(pk, &metadata.IndexMetadata{Name: "t_a", Unique: true, Definition: plainIndex, Expressions: []string{"a"}}))
	require.NoError(t, err)
	require.Contains(t, got, ` ON CONFLICT ("a") DO UPDATE SET "Id" = EXCLUDED."Id", "b" = EXCLUDED."b";`)

	got, err = restore(table(pk, &metadata.IndexMetadata{Name: "t_a_key", Unique: true, IsConstraint: true, Expressions: []string{"a"}}))
	require.NoError(t, err)
	require.Contains(t, got, ` ON CONFLICT ON CONSTRAINT "t_a_key" DO UPDATE`)

	// No key the UPDATE leaves alone can be named: a partial index, an index
	// on an expression, or one whose definition is unknown.
	for _, index := range []*metadata.IndexMetadata{
		{Name: "t_a", Unique: true, Definition: plainIndex + " WHERE (a > 0)", Expressions: []string{"a"}},
		{Name: "t_lower", Unique: true, Definition: "CREATE UNIQUE INDEX t_lower ON public.t USING btree (lower(a))", Expressions: []string{"lower(a)"}},
		{Name: "t_a", Unique: true, Expressions: []string{"a"}},
	} {
		_, err = restore(table(pk, index))
		require.EqualError(t, err, "no disjoint unique key found for table", index.Expressions)
	}
}

func TestRestoreSQLIdentityAlways(t *testing.T) {
	table := &metadata.TableMetadata{Name: "t", Columns: []*metadata.ColumnMetadata{
		{Name: "id", IdentityGeneration: metadata.ColumnMetadata_ALWAYS}, {Name: "a"},
	}}
	pk := &metadata.IndexMetadata{Name: "t_pkey", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"id"}}
	options := RestoreOptions{OriginalSchema: "public", OriginalTable: "t", BackupSchema: "bb", BackupTable: "b", Table: table, PrimaryKey: pk}
	dmls := BackedUpDMLs(parseScript(t, "DELETE FROM t; UPDATE t SET a = 1; UPDATE t SET id = DEFAULT"))

	got, err := RestoreSQL([]ast.Node{dmls[0].Node}, options)
	require.NoError(t, err)
	require.Contains(t, got, `INSERT INTO "public"."t" ("id", "a") OVERRIDING SYSTEM VALUE SELECT "id", "a" FROM "bb"."b";`)

	got, err = RestoreSQL([]ast.Node{dmls[1].Node}, options)
	require.NoError(t, err)
	require.Contains(t, got, `INSERT INTO "public"."t" ("id", "a") OVERRIDING SYSTEM VALUE SELECT "id", "a" FROM "bb"."b" ON CONFLICT ON CONSTRAINT "t_pkey" DO UPDATE SET "a" = EXCLUDED."a";`)

	// DO UPDATE can set the column only to DEFAULT.
	_, err = RestoreSQL([]ast.Node{dmls[2].Node}, options)
	require.EqualError(t, err, `cannot restore an UPDATE of "id", a GENERATED ALWAYS identity column`)
}
