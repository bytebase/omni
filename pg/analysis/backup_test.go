package analysis

import (
	"testing"

	"github.com/bytebase/omni/metadata"
	omnipg "github.com/bytebase/omni/pg"
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
	script := "SET search_path = public; DELETE FROM t; UPDATE t SET a = 1; UPDATE t SET id = 2, t.a = 3"
	stmts := parseScript(t, script)
	dmls := BackedUpDMLs(stmts)
	require.Len(t, dmls, 3)
	require.Equal(t, "SET search_path = public", PrependStatements(stmts))

	got, err := RestoreSQL(dmls[0].Node, options)
	require.NoError(t, err)
	require.Equal(t, "/*\nOriginal SQL:\nDELETE FROM t\n*/\nINSERT INTO \"public\".\"t\" (\"id\", \"a\") SELECT \"id\", \"a\" FROM \"bbdataarchive\".\"_bb__t_public\";", got)

	options.PrependStatements = PrependStatements(stmts)
	got, err = RestoreSQL(dmls[1].Node, options)
	require.NoError(t, err)
	require.Equal(t, "SET search_path = public;\n/*\nOriginal SQL:\nDELETE FROM t\n*/\nINSERT INTO \"public\".\"t\" (\"id\", \"a\") SELECT \"id\", \"a\" FROM \"bbdataarchive\".\"_bb__t_public\" ON CONFLICT ON CONSTRAINT \"t_pkey\" DO UPDATE SET \"a\" = EXCLUDED.\"a\";", got)

	// An UPDATE of every unique key cannot be restored.
	_, err = RestoreSQL(dmls[2].Node, options)
	require.EqualError(t, err, "no disjoint unique key found for table")
}
