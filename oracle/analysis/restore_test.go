package analysis

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
)

func restoreTable() (*metadata.TableMetadata, *metadata.IndexMetadata) {
	// The sync writes a key column quoted, which is what the merge statement and
	// the nullability check must read as the column.
	pk := &metadata.IndexMetadata{Name: "PK", Primary: true, Unique: true, Expressions: []string{`"ID"`}}
	return &metadata.TableMetadata{
		Name: "T",
		Columns: []*metadata.ColumnMetadata{
			{Name: "ID"}, {Name: "A"}, {Name: "B"},
			{Name: "G", Generation: &metadata.GenerationMetadata{Expression: "A + 1"}},
			{Name: "H", IsInvisible: true},
		},
		Indexes: []*metadata.IndexMetadata{
			pk,
			{Name: "UK", Unique: true, Expressions: []string{`"B"`}},
			{Name: "FUK", Unique: true, Type: "FUNCTION-BASED NORMAL", Expressions: []string{`"A"`}},
		},
	}, pk
}

func restore(t *testing.T, sql string) (string, error) {
	t.Helper()
	node, err := FindFirstDML(sql)
	require.NoError(t, err)
	require.NotNil(t, node)
	table, pk := restoreTable()
	return RestoreSQL(sql, node, RestoreOptions{
		SourceDatabase: "DB",
		SourceTable:    "T",
		BackupDatabase: "BB",
		BackupTable:    "B",
		Table:          table,
		PrimaryKey:     pk,
	})
}

func TestRestoreSQL(t *testing.T) {
	got, err := restore(t, "DELETE FROM t WHERE a = 1")
	require.NoError(t, err)
	require.Equal(t, "/*\nOriginal SQL:\nDELETE FROM t WHERE a = 1\n*/\n"+`INSERT INTO "DB"."T" ("ID", "A", "B") SELECT "ID", "A", "B" FROM "BB"."B";`, got)

	got, err = restore(t, "UPDATE t SET a = 1")
	require.NoError(t, err)
	require.Equal(t, "/*\nOriginal SQL:\nUPDATE t SET a = 1\n*/\n"+
		"MERGE INTO \"DB\".\"T\" t\nUSING \"BB\".\"B\" b\n  ON( t.\"ID\" = b.\"ID\")\nWHEN MATCHED THEN\n  UPDATE SET t.\"A\" = b.\"A\""+
		"\nWHEN NOT MATCHED THEN\n INSERT (\"ID\", \"A\", \"B\") VALUES (b.\"ID\", b.\"A\", b.\"B\");", got)

	// The primary key is set, so the next unique key that is not
	// function-based decides.
	got, err = restore(t, "UPDATE t SET (ID, A) = (SELECT 1, 2 FROM DUAL)")
	require.NoError(t, err)
	require.Contains(t, got, `ON( t."B" = b."B")`)

	_, err = restore(t, "UPDATE t SET id = 1, b = 2")
	require.ErrorContains(t, err, "no disjoint unique key found for DB.T")

	// A unique key on a nullable column cannot match its rows' NULLs.
	node, err := FindFirstDML("UPDATE t SET id = 1")
	require.NoError(t, err)
	table, pk := restoreTable()
	table.Columns[2].Nullable = true
	_, err = RestoreSQL("UPDATE t SET id = 1", node, RestoreOptions{SourceDatabase: "DB", SourceTable: "T", BackupDatabase: "BB", BackupTable: "B", Table: table, PrimaryKey: pk})
	require.EqualError(t, err, "no disjoint unique key found for DB.T")

	// A unique constraint on a non-unique index enforces uniqueness all the
	// same, so its keys key the MERGE although the index is not unique.
	table = &metadata.TableMetadata{
		Name:    "T",
		Columns: []*metadata.ColumnMetadata{{Name: "ID"}, {Name: "A"}},
		Indexes: []*metadata.IndexMetadata{{Name: "IX", IsConstraint: true, Expressions: []string{`"ID"`}}},
	}
	node, err = FindFirstDML("UPDATE t SET a = 1")
	require.NoError(t, err)
	got, err = RestoreSQL("UPDATE t SET a = 1", node, RestoreOptions{SourceDatabase: "DB", SourceTable: "T", BackupDatabase: "BB", BackupTable: "B", Table: table})
	require.NoError(t, err)
	require.Contains(t, got, `ON( t."ID" = b."ID")`)

	// SET VALUE (t) replaces whole rows, whose columns the restore cannot
	// list.
	_, err = restore(t, "UPDATE t x SET VALUE(x) = typ(1, 2)")
	require.ErrorContains(t, err, "cannot restore SET VALUE")

	got, err = restore(t, "DELETE FROM t WHERE a = '"+strings.Repeat("x", 1200)+"'")
	require.NoError(t, err)
	require.Contains(t, got, "...\n*/\n")

	// A "*/" in the statement must not end the comment.
	got, err = restore(t, "DELETE FROM t WHERE a = '*/ DROP TABLE x; /*'")
	require.NoError(t, err)
	require.Equal(t, "/*\nOriginal SQL:\nDELETE FROM t WHERE a = '* / DROP TABLE x; /*'\n*/\n"+`INSERT INTO "DB"."T" ("ID", "A", "B") SELECT "ID", "A", "B" FROM "BB"."B";`, got)
}

func TestFindFirstDML(t *testing.T) {
	node, err := FindFirstDML("SELECT 1 FROM DUAL;\nINSERT INTO t VALUES (1)")
	require.NoError(t, err)
	require.Nil(t, node)

	// A slice of a larger statement still yields its DML.
	node, err = FindFirstDML("BEGIN\nUPDATE t SET a = 1;\nEND")
	require.NoError(t, err)
	require.NotNil(t, node)

	_, err = FindFirstDML("UPDATE")
	require.ErrorContains(t, err, "failed to parse statement")
}

func TestNormalizeRestoreSQL(t *testing.T) {
	require.Equal(t, "DELETE FROM t", NormalizeRestoreSQL("  DELETE FROM t ; \n"))
}
