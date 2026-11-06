package catalog

import (
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// tidbRefusedForms are MySQL statements the differential scenarios forked from
// mysql/catalog ran on both engines, and TiDB v8.5.5 has no grammar for: column
// VISIBLE and INVISIBLE (as written, after a generated column, or inside a
// MySQL executable comment, which TiDB runs), ALTER COLUMN ... SET VISIBLE and
// SET INVISIBLE, ALTER VIEW, and a bare WITH CHECK OPTION. TiDB answers 1064 for
// each, and the catalog refuses each without creating or changing anything.
// tidb/parser/view_visibility_forms_oracle_test.go holds the grammar.
var tidbRefusedForms = []string{
	"CREATE TABLE t_vinv (a INT, b VARCHAR(255) INVISIBLE)",
	"CREATE TABLE t_ivis (a INT VISIBLE)",
	"CREATE TABLE t_gen6 (col1 INT, col2 INT, col3 INT GENERATED ALWAYS AS (col1 + col2) VIRTUAL INVISIBLE)",
	"CREATE TABLE t_vis (id INT PRIMARY KEY, val INT /*!80023 INVISIBLE */)",
	"ALTER TABLE t ALTER COLUMN b SET VISIBLE",
	"ALTER TABLE t ALTER COLUMN b SET INVISIBLE",
	"ALTER VIEW v AS SELECT 2 AS b, 3 AS c",
	"ALTER ALGORITHM = MERGE VIEW v AS SELECT 2 AS b",
	"CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0 WITH CHECK OPTION",
}

// TestCatalogRefusesTiDBRefusedForms holds the catalog to TiDB's answer without
// a server.
func TestCatalogRefusesTiDBRefusedForms(t *testing.T) {
	for _, sql := range tidbRefusedForms {
		t.Run(sql, func(t *testing.T) {
			c := setupWithDB(t)
			mustExec(t, c, "CREATE TABLE t (a INT, b INT)")
			mustExec(t, c, "CREATE VIEW v AS SELECT 1 AS a")
			before := c.ShowCreateTable("testdb", "t")
			if _, err := c.Exec(sql, nil); err == nil {
				t.Fatalf("Exec accepted %q, which TiDB refuses", sql)
			}
			db := c.GetDatabase("testdb")
			if len(db.Tables) != 1 || len(db.Views) != 1 {
				t.Errorf("the refused statement left %d tables and %d views, want 1 and 1", len(db.Tables), len(db.Views))
			}
			if after := c.ShowCreateTable("testdb", "t"); after != before {
				t.Errorf("the refused statement changed table t:\n%s\nwant\n%s", after, before)
			}
			if v := db.Views[toLower("v")]; v == nil || len(v.Columns) != 1 {
				t.Error("the refused statement changed view v")
			}
		})
	}
}

// TestContainerTiDBRefusedForms asks TiDB that the record above is its answer.
// A syntax error comes before name resolution, so no database is needed.
func TestContainerTiDBRefusedForms(t *testing.T) {
	tc := startTiDBForCatalog(t)
	for _, sql := range tidbRefusedForms {
		t.Run(sql, func(t *testing.T) {
			_, err := tc.db.ExecContext(tc.ctx, sql)
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Fatalf("TiDB gave no syntax verdict: %v", err)
			}
			if myErr.Number != 1064 {
				t.Fatalf("TiDB answered %d, want 1064: %v", myErr.Number, err)
			}
		})
	}
}
