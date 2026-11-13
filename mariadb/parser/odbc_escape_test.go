package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// TestODBCEscape: {d|t|ts 'string'} is a temporal literal when the string is a
// valid one, and the string otherwise, as the server falls back to it.
func TestODBCEscape(t *testing.T) {
	for _, tc := range []struct {
		sql      string
		temporal bool
	}{
		{"SELECT {d '2020-02-29'}", true},
		{"SELECT {t '23:59:59.5'}", true},
		{"SELECT {ts '2020-01-01 00:00:00'}", true},
		{"SELECT {d 'not-a-date'}", false},
		{"SELECT {d '2021-02-29'}", false},
		{"SELECT {x '2020-01-01'}", false},
		{"SELECT {D '2020-01-01'}", false},
		{"SELECT {d ('2020-01-01')}", true},
		{"SELECT {d (('2020-01-01'))}", true},
		{"SELECT {d '2020-1-1'}", true},
		{"SELECT {d '2020/01/01'}", true},
		{"SELECT {d '20200101'}", true},
		{"SELECT {ts '2020-1-1 1:2:3'}", true},
		{"SELECT {ts '2020-01-01T12:00:00'}", true},
		{"SELECT {d '2020-13-01'}", false},
		{"SELECT {t '25:00:00'}", true},
		{"SELECT {t '-01:00:00'}", true},
		{"SELECT {t '839:00:00'}", false},
		{"SELECT {ts '2020-01-01 25:00:00'}", false},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.sql, err)
		}
		e := list.Items[0].(*nodes.SelectStmt).TargetList[0]
		if _, ok := e.(*nodes.TemporalLit); ok != tc.temporal {
			t.Errorf("Parse(%q) = %T, want a temporal literal %v", tc.sql, e, tc.temporal)
		}
	}
}
