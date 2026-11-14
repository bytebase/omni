package catalog

import "testing"

// TestDeparseQueryKeepsWithTies checks that FETCH ... WITH TIES survives
// analysis: written as a LIMIT, it would drop the rows that tie with the last.
// FETCH ... ONLY is the LIMIT it means.
func TestDeparseQueryKeepsWithTies(t *testing.T) {
	c := correlatedCatalog(t, "ties")
	for _, tc := range []struct{ sql, want string }{
		{
			"SELECT id FROM employees ORDER BY salary FETCH FIRST ROW WITH TIES",
			"select `employees`.`id` from `ties`.`employees` order by `employees`.`salary` fetch first 1 rows with ties",
		},
		{
			"SELECT id FROM employees ORDER BY salary OFFSET 1 ROWS FETCH NEXT 2 ROWS WITH TIES",
			"select `employees`.`id` from `ties`.`employees` order by `employees`.`salary` offset 1 rows fetch first 2 rows with ties",
		},
		{
			"SELECT id FROM employees ORDER BY salary FETCH FIRST 2 ROWS ONLY",
			"select `employees`.`id` from `ties`.`employees` order by `employees`.`salary` limit 2",
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
