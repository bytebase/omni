package analysis

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
)

func restoreTable() (*metadata.TableMetadata, *metadata.IndexMetadata) {
	pk := &metadata.IndexMetadata{Name: "PK", Primary: true, Unique: true, Expressions: []string{"ID"}}
	return &metadata.TableMetadata{
		Name: "T",
		Columns: []*metadata.ColumnMetadata{
			{Name: "ID"}, {Name: "A"}, {Name: "B"},
			{Name: "G", Generation: &metadata.GenerationMetadata{Expression: "A + 1"}},
			{Name: "H", IsInvisible: true},
		},
		Indexes: []*metadata.IndexMetadata{
			pk,
			{Name: "UK", Unique: true, Expressions: []string{"B"}},
			{Name: "FUK", Unique: true, Type: "FUNCTION-BASED NORMAL", Expressions: []string{"A"}},
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

	got, err = restore(t, "DELETE FROM t WHERE a = '"+strings.Repeat("x", 1200)+"'")
	require.NoError(t, err)
	require.Contains(t, got, "...\n*/\n")
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
