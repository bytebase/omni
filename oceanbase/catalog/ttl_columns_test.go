package catalog

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	obparser "github.com/bytebase/omni/oceanbase/parser"
)

// TestTTLColumnsComeFromTheParse holds the TTL columns to what the parser read, so
// the catalog needs no second reading of the option's text: a comment anywhere in
// it changes nothing, and a qualified reference names the column last.
func TestTTLColumnsComeFromTheParse(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want []string
	}{
		{"ALTER TABLE a TTL (dt + INTERVAL 1 DAY)", []string{"dt"}},
		{"ALTER TABLE a TTL (dt + INTERVAL 1 DAY, c + INTERVAL 1 HOUR)", []string{"dt", "c"}},
		{"ALTER TABLE a TTL /*x*/ (c + INTERVAL 1 DAY)", []string{"c"}},
		{"ALTER TABLE a TTL (c /*x*/ + INTERVAL 1 DAY)", []string{"c"}},
		{"ALTER TABLE a TTL (d.a.dt + INTERVAL 1 DAY)", []string{"dt"}},
	} {
		list, err := obparser.Parse(tc.sql)
		if err != nil {
			t.Errorf("%s: %v", tc.sql, err)
			continue
		}
		cmd := list.Items[0].(*nodes.AlterTableStmt).Commands[0]
		got := cmd.Option.TTLColumns
		if len(got) != len(tc.want) {
			t.Errorf("%s: got %v, want %v", tc.sql, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("%s: got %v, want %v", tc.sql, got, tc.want)
				break
			}
		}
	}
}
