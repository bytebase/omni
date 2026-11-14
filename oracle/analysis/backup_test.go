package analysis

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransformDMLToSelect(t *testing.T) {
	opts := BackupOptions{SourceDatabase: "DB", TargetDatabase: "BACKUP", TablePrefix: "rollback"}

	got, err := TransformDMLToSelect("UPDATE t1 SET c1 = 1 WHERE c2 = 2;\nDELETE FROM s.t2 x WHERE x.c1 = 1;", opts)
	require.NoError(t, err)
	require.Equal(t, []BackupStatement{
		{
			Statement:       "CREATE TABLE \"BACKUP\".\"rollback_T1_DB\" AS\n  SELECT \"T1\".* FROM t1 WHERE c2 = 2;",
			SourceSchema:    "DB",
			SourceTableName: "T1",
			TargetTableName: "rollback_T1_DB",
			Start:           0,
			End:             32,
		},
		{
			Statement:       "CREATE TABLE \"BACKUP\".\"rollback_T2_S\" AS\n  SELECT \"X\".* FROM s.t2 x WHERE x.c1 = 1;",
			SourceSchema:    "S",
			SourceTableName: "T2",
			TargetTableName: "rollback_T2_S",
			Start:           35,
			End:             67,
		},
	}, got)

	got, err = TransformDMLToSelect("DELETE FROM t1 WHERE a = 1;\nDELETE FROM t1 WHERE b = 2;", opts)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Contains(t, got[0].Statement, `WHERE ROWID IN (`)
	require.Contains(t, got[0].Statement, "\n    UNION\n")

	// DELETE t FROM u ...: the FROM's tables join the target's, and a
	// target row matching several of theirs is backed up once, by ROWID.
	got, err = TransformDMLToSelect("DELETE t1 FROM t2 WHERE t1.a = t2.a", opts)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Contains(t, got[0].Statement, `WHERE ROWID IN (`)
	require.Contains(t, got[0].Statement, "SELECT \"T1\".ROWID FROM t1, t2 WHERE t1.a = t2.a")

	_, err = TransformDMLToSelect("UPDATE t1 SET a = 1;\nDELETE FROM t1;", opts)
	require.ErrorContains(t, err, "different types on the same table")

	_, err = TransformDMLToSelect("UPDATE t1@remote SET a = 1", opts)
	require.ErrorContains(t, err, "database link")

	_, err = TransformDMLToSelect("SELECT 1 FROM DUAL", opts)
	require.ErrorIs(t, err, ErrNoBackupableDML)
}

func TestTransformDMLToSelectTableNameLength(t *testing.T) {
	statement := "DELETE FROM " + strings.Repeat("A", 40)
	opts := BackupOptions{SourceDatabase: "DB", TargetDatabase: "BACKUP", TablePrefix: "_p"}

	got, err := TransformDMLToSelect(statement, opts)
	require.NoError(t, err)
	require.Len(t, got[0].TargetTableName, 30)

	opts.Version = &Version{First: 12, Second: 2}
	got, err = TransformDMLToSelect(statement, opts)
	require.NoError(t, err)
	require.Equal(t, "_p_"+strings.Repeat("A", 40)+"_DB", got[0].TargetTableName)
}

func TestReferencesTable(t *testing.T) {
	ok, err := ReferencesTable("UPDATE s.t SET a = 1", "S", "T")
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = ReferencesTable("UPDATE t SET a = 1", "S", "T")
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = ReferencesTable("UPDATE t SET a = 1", "DB", "U")
	require.NoError(t, err)
	require.False(t, ok)

	ok, err = ReferencesTable("INSERT INTO t VALUES (1)", "DB", "T")
	require.NoError(t, err)
	require.False(t, ok)

	_, err = ReferencesTable("UPDATE", "DB", "T")
	require.ErrorContains(t, err, "failed to parse PLSQL")
}
