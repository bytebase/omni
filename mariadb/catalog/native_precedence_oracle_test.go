package catalog

import (
	"fmt"
	"testing"
)

// Changes to stored formatting must preserve nested operator binding and
// quoted CTE names, as well as the simple expressions in the broad oracle.
func TestNativeMariaDBNestedViewFormatting(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()
	c := New()
	if err := ctr.execSQL("CREATE TABLE t(a INT,b INT)"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE TABLE t(a INT,b INT)", nil); err != nil {
		t.Fatal(err)
	}
	for _, query := range []string{
		"SELECT 1 = (a = ANY(SELECT 1)) AS x FROM t",
		"SELECT (a <> ALL(SELECT 1)) + 1 AS x FROM t",
		"SELECT (a OR b) IS NULL AS x FROM t",
		"SELECT CONCAT(a,b) COLLATE utf8mb4_bin AS x FROM t",
		"WITH `cte` AS (SELECT a FROM t) SELECT a AS x FROM `cte`",
	} {
		t.Run(query, func(t *testing.T) {
			if err := ctr.execSQL("DROP VIEW IF EXISTS native_precedence"); err != nil {
				t.Fatal(err)
			}
			ddl := fmt.Sprintf("CREATE VIEW native_precedence AS %s", query)
			if err := ctr.execSQL(ddl); err != nil {
				t.Fatal(err)
			}
			if _, err := c.Exec("DROP VIEW IF EXISTS native_precedence", nil); err != nil {
				t.Fatal(err)
			}
			results, err := c.Exec(ddl, nil)
			if err != nil || results[0].Error != nil {
				t.Fatalf("catalog: %v %v", err, results)
			}
			native, err := ctr.showCreateView("native_precedence")
			if err != nil {
				t.Fatal(err)
			}
			got := c.ShowCreateView("test", "native_precedence")
			if extractSelectBody(native) != extractSelectBody(got) {
				t.Errorf("native: %s\ncatalog: %s", native, got)
			}
		})
	}
}
