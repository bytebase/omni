package catalog

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
)

// invisibleColumnSnapshot marks column a invisible, which no TiDB table can
// have: v8.5.5 answers 1064 for INVISIBLE after a column type, even inside
// /*!80023 ... */ (see tidbRefusedForms), and Bytebase's TiDB sync never sets
// IsInvisible. The snapshot can only come from somewhere else.
func invisibleColumnSnapshot() *metadata.DatabaseSchemaMetadata {
	return snapshot(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{{
		Name:           "t",
		PrimaryKeyType: "CLUSTERED",
		Columns: []*metadata.ColumnMetadata{
			{Name: "id", Type: "int"},
			{Name: "a", Type: "int", Nullable: true, Default: "NULL", IsInvisible: true},
			{Name: "b", Type: "varchar(10)", Nullable: true, Default: "NULL"},
		},
		Indexes: []*metadata.IndexMetadata{
			{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Expressions: []string{"id"}},
		},
	}}})
}

// TestLoadMetadataIgnoresColumnVisibility loads a column marked invisible as
// an ordinary one, so SHOW CREATE TABLE writes text the catalog reads back.
func TestLoadMetadataIgnoresColumnVisibility(t *testing.T) {
	c, report := loadTiDBSnapshot(t, invisibleColumnSnapshot())
	requireReport(t, report)
	requireTable(t, c, "t")
	got := c.ShowCreateTable(snapshotDB, "t")
	if strings.Contains(strings.ToUpper(got), "INVISIBLE") {
		t.Fatalf("SHOW CREATE TABLE writes a column visibility TiDB refuses:\n%s", got)
	}

	back := setupWithDB(t)
	if _, err := back.Exec(got, nil); err != nil {
		t.Fatalf("the catalog refuses its own SHOW CREATE TABLE: %v\n%s", err, got)
	}
	if again := back.ShowCreateTable(snapshotDB, "t"); again != got {
		t.Errorf("SHOW CREATE TABLE after reading it back:\n%s\nwant\n%s", again, got)
	}
}

// TestContainerLoadMetadataColumnVisibility asks TiDB to run the SHOW CREATE
// TABLE text of the snapshot and prints back the same text.
func TestContainerLoadMetadataColumnVisibility(t *testing.T) {
	c, report := loadTiDBSnapshot(t, invisibleColumnSnapshot())
	requireReport(t, report)
	got := c.ShowCreateTable(snapshotDB, "t")

	tc := startTiDBForCatalog(t)
	mustExecTiDB(t, tc, "DROP DATABASE IF EXISTS omni_cat_column_visibility")
	mustExecTiDB(t, tc, "CREATE DATABASE omni_cat_column_visibility")
	t.Cleanup(func() {
		_, _ = tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_cat_column_visibility")
		_, _ = tc.db.ExecContext(tc.ctx, "USE test")
	})
	mustExecTiDB(t, tc, "USE omni_cat_column_visibility")
	if _, err := tc.db.ExecContext(tc.ctx, got); err != nil {
		t.Fatalf("TiDB refuses the catalog's SHOW CREATE TABLE: %v\n%s", err, got)
	}
	var name, want string
	if err := tc.db.QueryRowContext(tc.ctx, "SHOW CREATE TABLE t").Scan(&name, &want); err != nil {
		t.Fatalf("SHOW CREATE TABLE on TiDB: %v", err)
	}
	if got != want {
		t.Errorf("catalog SHOW CREATE TABLE:\n%s\nTiDB:\n%s", got, want)
	}
}
