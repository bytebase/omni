package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGetQuerySpanPivot holds the output of PIVOT and UNPIVOT, in FROM and as
// pipe operators, to the columns BigQuery returns: PIVOT's grouping columns
// and then one column per value and aggregate, reading the aggregate's
// sources; UNPIVOT's untouched columns, its value columns, reading the IN
// columns, and its name column, which reads none.
func TestGetQuerySpanPivot(t *testing.T) {
	tests := []struct {
		sql  string
		want [][]string // each result's name, then its sources
	}{
		{"SELECT * FROM (SELECT r, q, v FROM t) PIVOT (SUM(v) FOR q IN ('Q1', 'Q2'))",
			[][]string{{"r", "t.r"}, {"Q1", "t.v"}, {"Q2", "t.v"}}},
		{"SELECT p.* FROM (SELECT s, c, d FROM t) AS q PIVOT (MAX(s) AS m, COUNT(*) AS n FOR c IN ('a', 'b' AS bb)) AS p",
			[][]string{{"d", "t.d"}, {"m_a", "t.s"}, {"n_a"}, {"m_bb", "t.s"}, {"n_bb"}}},
		{"SELECT * FROM (SELECT x FROM t) PIVOT (SUM(x) FOR x IN (1, -2, true, NULL, DATE '2013-11-25'))",
			[][]string{{"_1", "t.x"}, {"minus_2", "t.x"}, {"TRUE", "t.x"}, {"NULL", "t.x"}, {"_2013_11_25", "t.x"}}},
		{"SELECT * FROM (SELECT r, a, b FROM t) UNPIVOT (v FOR n IN (a, b))",
			[][]string{{"r", "t.r"}, {"v", "t.a", "t.b"}, {"n"}}},
		{"SELECT * FROM (SELECT a, b, c, d FROM t) UNPIVOT ((v1, v2) FOR n IN ((a, b) AS 'x', (c, d) AS 'y'))",
			[][]string{{"v1", "t.a", "t.c"}, {"v2", "t.b", "t.d"}, {"n"}}},
		{"SELECT x FROM t |> PIVOT (SUM(x) FOR x IN (1, 2))",
			[][]string{{"_1", "t.x"}, {"_2", "t.x"}}},
		{"FROM (SELECT s, c FROM t) PIVOT (MAX(s) FOR c IN ('a', 'b')) |> SELECT *",
			[][]string{{"a", "t.s"}, {"b", "t.s"}}},
		{"SELECT 0 AS p |> CROSS JOIN (SELECT s, c FROM t) PIVOT (MAX(s) FOR c IN ('a'))",
			[][]string{{"p"}, {"a", "t.s"}}},
		{"FROM (SELECT a, b, k FROM t) |> UNPIVOT (v FOR n IN (a, b)) AS u |> SELECT u.v",
			[][]string{{"v", "t.a", "t.b"}}},
	}
	for _, tc := range tests {
		span, err := GetQuerySpan(tc.sql)
		require.NoError(t, err, tc.sql)
		var got [][]string
		for _, r := range span.Results {
			got = append(got, append([]string{r.Name}, resultSourceColumns(r)...))
		}
		require.Equal(t, tc.want, got, "statement: %s", tc.sql)
	}

	// A value whose column name omni does not derive fails closed: a
	// reference to the column by its name would read nothing.
	for _, sql := range []string{
		"SELECT * FROM (SELECT x FROM t) PIVOT (SUM(x) FOR x IN (1.5))",
		"SELECT x FROM t |> PIVOT (SUM(x) FOR x IN (TIMESTAMP '2020-01-01'))",
	} {
		_, err := GetQuerySpan(sql)
		require.Error(t, err, sql)
	}
}
