package deparse

import (
	"testing"

	ast "github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// TestTableSampleDeparse checks the clause survives on the way out, in every part
// the grammar lets a statement write. A sample dropped here leaves a statement that
// parses and reads the whole table, which no accept-or-reject check can see. Each
// result also reads back to the same text.
//
// REGIONS is the one word this does not take from the server: see writeTableSample.
func TestTableSampleDeparse(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{"SELECT 1 FROM t TABLESAMPLE (50)", "tablesample (50)"},
		{"SELECT 1 FROM t TABLESAMPLE ()", "tablesample ()"},
		{"SELECT 1 FROM t TABLESAMPLE SYSTEM (50)", "tablesample system (50)"},
		{"SELECT 1 FROM t TABLESAMPLE BERNOULLI (50 PERCENT)", "tablesample bernoulli (50 percent)"},
		{"SELECT 1 FROM t TABLESAMPLE REGIONS ()", "tablesample regions ()"},
		{"SELECT 1 FROM t TABLESAMPLE REGIONS (50 ROWS)", "tablesample regions (50 rows)"},
		{"SELECT 1 FROM t TABLESAMPLE (50) REPEATABLE (7)", "tablesample (50) repeatable (7)"},
		{"SELECT 1 FROM t TABLESAMPLE REGIONS () REPEATABLE (7)", "tablesample regions () repeatable (7)"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			want := "select 1 AS `1` from `t` " + tc.want
			got := deparseSampleSelect(t, tc.sql)
			if got != want {
				t.Errorf("got  %s\nwant %s", got, want)
			}
			if again := deparseSampleSelect(t, got); again != got {
				t.Errorf("reads back as %s, want %s", again, got)
			}
		})
	}

	t.Run("after the alias", func(t *testing.T) {
		got := deparseSampleSelect(t, "SELECT 1 AS x FROM t AS a TABLESAMPLE (50), u b TABLESAMPLE (1 ROWS)")
		want := "select 1 AS `x` from (`t` `a` tablesample (50) join `u` `b` tablesample (1 rows))"
		if got != want {
			t.Errorf("got  %s\nwant %s", got, want)
		}
	})
}

func deparseSampleSelect(t *testing.T, sql string) string {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	return DeparseSelect(list.Items[0].(*ast.SelectStmt))
}
