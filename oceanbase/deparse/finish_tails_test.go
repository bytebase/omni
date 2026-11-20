package deparse

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
)

func TestResolveLambdaAndTableCollection(t *testing.T) {
	tables := map[string]*ResolverTable{
		"t": {
			Name: "t",
			Columns: []ResolverColumn{
				{Name: "a", Position: 1},
				{Name: "arr", Position: 2},
			},
		},
	}
	r := &Resolver{
		Lookup:         func(name string) *ResolverTable { return tables[name] },
		DefaultCharset: "utf8mb4",
	}

	resolveAndDeparse := func(t *testing.T, sql string) string {
		t.Helper()
		list, err := parser.Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		sel := list.Items[0].(*ast.SelectStmt)
		r.Source = sql
		resolved := r.Resolve(sel)
		return DeparseSelect(resolved)
	}

	// Lambda body qualifies outer table column `a` while leaving parameter `x` unqualified.
	if got, want := resolveAndDeparse(t, "SELECT ARRAY_FILTER(x -> x > a, arr) AS res FROM t"),
		"select array_filter(`x` -> (`x` > `t`.`a`),`t`.`arr`) AS `res` from `t`"; got != want {
		t.Errorf("lambda resolve:\n got: %s\nwant: %s", got, want)
	}

	// TABLE(arr) resolves `arr` against preceding FROM item `t` without self-qualifying as `c.arr`.
	if got, want := resolveAndDeparse(t, "SELECT c.arr FROM t, TABLE(arr) c"),
		"select `c`.`arr` AS `arr` from (`t` join table(`t`.`arr`) `c`)"; got != want {
		t.Errorf("TABLE(arr) resolve:\n got: %s\nwant: %s", got, want)
	}

	// Unknown-column TABLE(f()) preserves `*` and `c.*` instead of expanding to an empty select list.
	if got, want := resolveAndDeparse(t, "SELECT * FROM TABLE(f()) c"),
		"select * from table(f()) `c`"; got != want {
		t.Errorf("TABLE(f()) star:\n got: %s\nwant: %s", got, want)
	}
	if got, want := resolveAndDeparse(t, "SELECT c.* FROM TABLE(f()) c"),
		"select `c`.* from table(f()) `c`"; got != want {
		t.Errorf("TABLE(f()) qualified star:\n got: %s\nwant: %s", got, want)
	}
}
