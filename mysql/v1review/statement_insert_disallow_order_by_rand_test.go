package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementInsertDisallowOrderByRand(t *testing.T) {
	rand := stmtFinding{Code: v1review.InsertOrderByRand}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "insert select order by rand", sql: "INSERT INTO t SELECT * FROM s ORDER BY RAND();", want: []stmtFinding{rand}},
		{name: "rand with seed, any case", sql: "INSERT INTO t (a) SELECT a FROM s ORDER BY rand(1), a, Rand();", want: []stmtFinding{rand, rand}},
		{name: "union sides and outer order", sql: "INSERT INTO t (SELECT a FROM s ORDER BY RAND()) UNION (SELECT a FROM u ORDER BY RAND()) ORDER BY RAND();", want: []stmtFinding{rand, rand, rand}},
		{name: "order by other function not judged", sql: "INSERT INTO t SELECT a FROM s ORDER BY UUID();"},
		{name: "subquery not judged", sql: "INSERT INTO t SELECT a FROM (SELECT a FROM s ORDER BY RAND()) x;"},
		{name: "insert values not judged", sql: "INSERT INTO t VALUES (1);"},
		{name: "select not judged", sql: "SELECT a FROM s ORDER BY RAND();"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementInsertDisallowOrderByRand)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
