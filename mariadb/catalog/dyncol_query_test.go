package catalog

import "testing"

// TestDeparseQueryKeepsDynColTypes checks that a dynamic column function's AS
// types survive analysis: COLUMN_GET does not parse without its type, and a
// typed COLUMN_CREATE value converts differently from an untyped one.
func TestDeparseQueryKeepsDynColTypes(t *testing.T) {
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
