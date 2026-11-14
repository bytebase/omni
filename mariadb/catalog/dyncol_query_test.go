package catalog

import (
	"strings"
	"testing"
)

// TestDeparseQueryKeepsSpecialCallForms checks that the calls with grammar of
// their own survive analysis and parse back: a dynamic column function's AS
// types (COLUMN_GET does not parse without its type), PERCENTILE_CONT's WITHIN
// GROUP and window, MEDIAN's window, and TRIM_ORACLE, which is another function
// than TRIM.
func TestDeparseQueryKeepsSpecialCallForms(t *testing.T) {
	c := correlatedCatalog(t, "dyn")
	for _, tc := range []struct{ sql, want string }{
		{
			"SELECT COLUMN_GET(name, 'k' AS INT) AS x FROM employees",
			"select column_get(`employees`.`name`, 'k' as signed) as `x` from `dyn`.`employees`",
		},
		{
			"SELECT COLUMN_CREATE('k', id AS UNSIGNED INT, 'v', name) AS x FROM employees",
			"select column_create('k', `employees`.`id` AS unsigned int, 'v', `employees`.`name`) as `x` from `dyn`.`employees`",
		},
		{
			"SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY salary DESC) OVER (PARTITION BY department_id) AS x FROM employees",
			"select percentile_cont(0.5) within group (order by `employees`.`salary` desc) over (partition by `employees`.`department_id`) as `x` from `dyn`.`employees`",
		},
		{
			"SELECT MEDIAN(salary) OVER () AS x FROM employees",
			"select median(`employees`.`salary`) over () as `x` from `dyn`.`employees`",
		},
		{
			"SELECT TRIM_ORACLE(name) AS x, TRIM_ORACLE(LEADING 'a' FROM name) AS y FROM employees",
			"select oracle_schema.trim(`employees`.`name`) as `x`, oracle_schema.trim(leading 'a' from `employees`.`name`) as `y` from `dyn`.`employees`",
		},
		{
			"SELECT TRIM('a' FROM name) AS x, TRIM(TRAILING FROM name) AS y FROM employees",
			"select trim(both 'a' from `employees`.`name`) as `x`, trim(trailing from `employees`.`name`) as `y` from `dyn`.`employees`",
		},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			got := DeparseQuery(analyzeQuantified(t, c, tc.sql))
			if got != tc.want {
				t.Errorf("DeparseQuery:\n got: %s\nwant: %s", got, tc.want)
			}
			if again := DeparseQuery(analyzeQuantified(t, c, got)); again != got {
				t.Errorf("DeparseQuery does not read back:\n first: %s\nsecond: %s", got, again)
			}
		})
	}
}

// TestDynColCharTakesTheConnectionCharset checks that a dynamic column's CHAR
// value type with no character set is stored with the connection's, as
// mariadb:11.8.8 stores it after SET NAMES latin1 COLLATE latin1_bin:
// char charset latin1 collate latin1_bin, and a COLLATE written alone keeps
// its own character set.
func TestDynColCharTakesTheConnectionCharset(t *testing.T) {
	c := New()
	results, err := c.Exec("CREATE DATABASE d; USE d; SET NAMES latin1 COLLATE latin1_bin; "+
		"CREATE VIEW v AS SELECT COLUMN_CREATE('a', 'x' AS CHAR, 'b', 'y' AS CHAR BINARY, 'c', 'z' AS CHAR COLLATE utf8mb4_bin, 'd', 'w' AS NCHAR) AS c1", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatal(r.Error)
		}
	}
	want := "column_create('a','x' AS char charset latin1 collate latin1_bin ,'b','y' AS char charset latin1 collate latin1_bin ," +
		"'c','z' AS char charset utf8mb4 collate utf8mb4_bin ,'d','w' AS char charset utf8mb3 collate utf8mb3_general_ci )"
	if got := c.ShowCreateView("d", "v"); !strings.Contains(got, want) {
		t.Errorf("ShowCreateView:\n got: %s\nwant it to hold: %s", got, want)
	}
}
