package parser

import "testing"

// fetchFirstCases are the places fetch_first_clause stands in LIMIT's stead,
// each measured on mariadb:11.8.8 with PREPARE: after a SELECT and after the
// VALUES list of an INSERT or REPLACE, and nowhere a LIMIT does not reach.
var fetchFirstCases = []struct {
	sql    string
	accept bool
}{
	{"SELECT a FROM t FETCH FIRST ROW ONLY", true},
	{"SELECT a FROM t ORDER BY a OFFSET 1 ROWS FETCH NEXT 3 ROWS WITH TIES", true},
	{"SELECT a FROM t OFFSET 2 ROWS", true},
	{"SELECT a FROM t OFFSET 2", false},
	{"SELECT a FROM t LIMIT 1 FETCH FIRST ROW ONLY", false},
	{"INSERT INTO t VALUES (1) FETCH FIRST ROW ONLY", true},
	{"INSERT INTO t VALUES (1) OFFSET 1 ROWS", true},
	{"INSERT INTO t VALUES (1) FETCH FIRST ROW ONLY ON DUPLICATE KEY UPDATE a = 1", true},
	{"REPLACE INTO t VALUES (1) FETCH FIRST ROW ONLY", true},
	{"INSERT INTO t VALUES (1) OFFSET 1", false},
	{"SELECT ((SELECT 1) FETCH FIRST ROW ONLY)", true},
	{"SELECT * FROM ((SELECT 1 AS a) OFFSET 1 ROWS) d", true},
	{"SELECT * FROM ((SELECT 1 AS a) ORDER BY a OFFSET 1 ROWS FETCH NEXT 1 ROW WITH TIES) d", true},
	{"SELECT * FROM (SELECT 1 AS a) offset", false},
	{"SELECT 1 IN ((SELECT 1) FETCH FIRST ROW ONLY)", false},
	{"DELETE FROM t FETCH FIRST ROW ONLY", false},
	{"UPDATE t SET a = 1 FETCH FIRST ROW ONLY", false},
}

func TestFetchFirstPlaces(t *testing.T) {
	for _, tc := range fetchFirstCases {
		_, err := Parse(tc.sql)
		if (err == nil) != tc.accept {
			t.Errorf("Parse(%q): err = %v, want accept = %v", tc.sql, err, tc.accept)
		}
	}
}
