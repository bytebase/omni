package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// mBangCases are MariaDB 11.8.8's verdicts on its own "/*M!" executable
// comments: it runs the body, as it runs "/*!", and the two do not nest in one
// another. "/*M!999999", mariadb-dump's sandbox marker, and "/*!999999" carry a
// version no server reaches and are skipped unread, nesting counted.
var mBangCases = []struct {
	sql    string
	accept bool
}{
	{`/*M!999999\- enable the sandbox mode */ SELECT 1`, true},
	{"/*M!999999 + 1 */ SELECT 1", true},
	{"SELECT 1 /*M! + 1 */", true},
	{"SELECT 1 /*M!100100 + 1 */", true},
	{"SELECT 1 /*M!+ 1 */", true},
	{"SELECT 1 /*m! + 1 */", true},
	{"SELECT 1 /*M! + 1 /*M! + 2 */", true},
	{"SELECT 1 /*M! + 1 /*! + 2 */", true},
	{"SELECT 1 /*! + 1 /*M! + 2 */", true},
	{"SELECT 1 /*M! + 1 '*/' */", true},
	{"/*M! SELECT 1 */", true},
	{"SELECT 1 /*M! + 1 /*M! + 2 */ */", false},
	{"SELECT 1 /*M! + 1", false},
	{"SELECT 1 /*M! , */", false},
	{"SELECT 1 /*M!999999 + 1 /*! + 2 */ */", true},
	{"SELECT 1 /*!999999 + 1 /*! + 2 */ */ + 3", true},
	{"SELECT 1 /*M!999999 + LENGTH('*/') */", false},
	{"SELECT 1 /*!999999 + LENGTH('*/') */", false},
}

// TestMBangComments pins omni to mBangCases, which TestMBangComments_Oracle
// takes from MariaDB.
func TestMBangComments(t *testing.T) {
	for _, c := range mBangCases {
		if _, err := Parse(c.sql); (err == nil) != c.accept {
			t.Errorf("%q: Parse error = %v, want accept %v", c.sql, err, c.accept)
		}
	}
}

// TestMBangCommentRunsItsBody: the body of "/*M!" is SQL, so "SELECT 1 /*M! + 1 */"
// selects an addition (MariaDB returns 2), and a script made of one such comment
// is a statement Split keeps rather than an empty segment it drops.
func TestMBangCommentRunsItsBody(t *testing.T) {
	list, err := Parse("SELECT 1 /*M! + 1 */")
	if err != nil {
		t.Fatal(err)
	}
	sel, ok := list.Items[0].(*nodes.SelectStmt)
	if !ok || len(sel.TargetList) != 1 {
		t.Fatalf("got %#v, want one select item", list.Items[0])
	}
	item := sel.TargetList[0]
	if rt, ok := item.(*nodes.ResTarget); ok {
		item = rt.Val
	}
	if _, ok := item.(*nodes.BinaryExpr); !ok {
		t.Errorf("select item is %T, want the addition the comment's body completes", item)
	}
	if segs := Split("/*M! SELECT 1 */"); len(segs) != 1 {
		t.Errorf("Split dropped the statement held in a /*M! comment: %+v", segs)
	}
	if segs := Split(`/*M!999999\- enable the sandbox mode */`); len(segs) != 0 {
		t.Errorf("Split kept mariadb-dump's sandbox marker as a statement: %+v", segs)
	}
}

// TestMBangComments_Oracle proves mBangCases against a live MariaDB 11.8.8.
func TestMBangComments_Oracle(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range mBangCases {
		v, _, msg := o.classify(c.sql)
		if v == vErrored {
			t.Fatalf("%q: MariaDB gave no verdict: %s", c.sql, msg)
		}
		ok := v == vAccepted
		if ok != c.accept {
			t.Errorf("%q: MariaDB accepts=%v, the case says %v", c.sql, ok, c.accept)
		}
		if _, err := Parse(c.sql); (err == nil) != ok {
			t.Errorf("%q: omni accepts=%v (err=%v), MariaDB accepts=%v", c.sql, err == nil, err, ok)
		}
	}
}
