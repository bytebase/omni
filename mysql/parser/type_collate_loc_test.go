package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// TestTypeCollateLoc holds a type's Loc through the opt_collate that follows it where a
// column's attributes do not: a JSON_TABLE column and a routine's parameter.
func TestTypeCollateLoc(t *testing.T) {
	for _, c := range []struct{ sql, want string }{
		{"SELECT * FROM JSON_TABLE('[]', '$' COLUMNS(a INT COLLATE utf8mb4_bin PATH '$')) t", "INT COLLATE utf8mb4_bin"},
		{"CREATE PROCEDURE p(a CHAR(1) CHARSET utf8mb4 COLLATE utf8mb4_bin) COMMIT", "CHAR(1) CHARSET utf8mb4 COLLATE utf8mb4_bin"},
	} {
		list, err := Parse(c.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", c.sql, err)
		}
		var got []string
		nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
			if dt, ok := n.(*nodes.DataType); ok {
				got = append(got, c.sql[dt.Loc.Start:dt.Loc.End])
			}
			return true
		})
		if len(got) != 1 || got[0] != c.want {
			t.Errorf("Parse(%q): types %q, want [%q]", c.sql, got, c.want)
		}
	}
}
