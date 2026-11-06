package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// ODBC escapes as TiDB v8.5.5 reads them, every verdict measured with PREPARE on
// pingcap/tidb:v8.5.5 on 2026-11-06.
//
// An expression may be { Identifier Expression } (SimpleExpr, parser.y:7941). The
// action compares the name as written: d, t and ts make a DATE, TIME or TIMESTAMP
// literal, and any other name, D and TS among them, leaves the expression alone. A
// FROM, UPDATE or multiple-table DELETE list item may be { Identifier TableRef }
// (EscapedTableRef, parser.y:9895), which keeps the table reference. Only
// TableRefs reads that, so no join follows it, and a single-table DELETE, which
// names a TableName, takes none.
//
// Not here: {d '2020-01-01' '2'}, which TiDB accepts because it joins adjacent
// strings into one, as it does everywhere; omni joins them nowhere yet.

type odbcEscapeCase struct {
	sql    string
	accept bool
}

var odbcEscapeCases = []odbcEscapeCase{
	// Expressions
	{"select {d '2020-01-01'}", true},
	{"select {D '2020-01-01'}", true},
	{"select {TS '2020-01-01 00:00:00'}", true},
	{"select {d 123}", true},
	{"select {d now()}", true},
	{"select {ts 1 + 1}", true},
	{"select {`d` '2020-01-01'}", true},
	{"select {abc 1}", true},
	{"select {select 1}", false},
	{"select {d}", false},
	{"select {}", false},
	{"select {d '2020-01-01'} + 1", true},
	{"select { d '2020-01-01' }", true},
	{"select {\"d\" '2020-01-01'}", false},
	{"select {d '2020-01-01'} as x", true},
	{"select {date '2020-01-01'}", true},
	{"select {timestamp '2020-01-01'}", true},
	{"select {from 1}", false},
	{"select {d (select 1)}", true},
	{"select {d 1, 2}", false},
	{"insert into t values ({d '2020-01-01'})", true},
	{"select {ts {abc {d '1989-09-10'}}}", true},
	{"select {d '1' } from t where a in ({d '1'})", true},
	{"select {oj 1}", true},
	{"select {d '1'", false},
	{"select {abc *} from t", false},
	{"select {abc t.*} from t", false},
	{"select {d *}", false},
	{"select {d default}", false},
	{"select {d default(a)} from t", true},
	{"insert into t values ({d *})", false},
	{"insert into t values ({abc default})", false},

	// Table references
	{"select * from {oj t1 left join t2 on t1.a = t2.a}", true},
	{"select * from {oj t1}", true},
	{"select * from {abc t1 left join t2 on t1.a = t2.a}", true},
	{"select * from {oj t1 left join t2 on t1.a = t2.a}, t3", true},
	{"select * from t3, {oj t1 left join t2 on t1.a = t2.a}", true},
	{"select * from {oj t1} join t2 on t1.a = t2.a", false},
	{"select * from ({oj t1})", true},
	{"select * from {oj {oj t1}}", false},
	{"update {oj t1} set a = 1", true},
	{"delete from {oj t1}", false},
	{"delete t1 from {oj t1 left join t2 on t1.a = t2.a}", true},
	{"delete from t1 using {oj t1 left join t2 on t1.a = t2.a}", true},
}

// TestODBCEscapeForms holds omni to the recorded verdicts without a server.
func TestODBCEscapeForms(t *testing.T) {
	for _, c := range odbcEscapeCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("Parse(%q) accepts=%v, want %v (TiDB %s): %v", c.sql, got, c.accept, tidbCorpusVersion, err)
			}
		})
	}
}

// TestODBCEscapeTrees checks the tree each escape leaves: a temporal literal
// spanning the braces for d, t and ts, and the expression or table reference
// alone for any other name.
func TestODBCEscapeTrees(t *testing.T) {
	selectItem := func(t *testing.T, sql string) nodes.ExprNode {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		return list.Items[0].(*nodes.SelectStmt).TargetList[0]
	}

	for _, c := range []struct {
		sql, typ, value, text string
	}{
		{"select {d '2020-01-01'}", "DATE", "2020-01-01", "{d '2020-01-01'}"},
		{"select {t '00:00:00'}", "TIME", "00:00:00", "{t '00:00:00'}"},
		{"select { ts '2020-01-01 00:00:00' }", "TIMESTAMP", "2020-01-01 00:00:00", "{ ts '2020-01-01 00:00:00' }"},
		{"select {`d` '2020-01-01'}", "DATE", "2020-01-01", "{`d` '2020-01-01'}"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			lit, ok := selectItem(t, c.sql).(*nodes.TemporalLit)
			if !ok {
				t.Fatalf("Parse(%q): item is not a TemporalLit", c.sql)
			}
			if lit.Type != c.typ || lit.Value != c.value || lit.Expr != nil {
				t.Errorf("Parse(%q) = %s %q expr %v, want %s %q", c.sql, lit.Type, lit.Value, lit.Expr, c.typ, c.value)
			}
			if got := c.sql[lit.Loc.Start:lit.Loc.End]; got != c.text {
				t.Errorf("Loc covers %q, want %q", got, c.text)
			}
		})
	}

	t.Run("non-string argument", func(t *testing.T) {
		lit, ok := selectItem(t, "select {d 123}").(*nodes.TemporalLit)
		if !ok {
			t.Fatal("item is not a TemporalLit")
		}
		if lit.Value != "" {
			t.Errorf("Value = %q, want empty", lit.Value)
		}
		if n, ok := lit.Expr.(*nodes.IntLit); !ok || n.Value != 123 {
			t.Errorf("Expr = %#v, want IntLit 123", lit.Expr)
		}
	})

	t.Run("nested", func(t *testing.T) {
		outer, ok := selectItem(t, "select {ts {abc {d '1989-09-10'}}}").(*nodes.TemporalLit)
		if !ok || outer.Type != "TIMESTAMP" {
			t.Fatalf("outer = %#v, want a TIMESTAMP TemporalLit", outer)
		}
		inner, ok := outer.Expr.(*nodes.TemporalLit)
		if !ok || inner.Type != "DATE" || inner.Value != "1989-09-10" {
			t.Errorf("inner = %#v, want DATE '1989-09-10'", outer.Expr)
		}
	})

	for _, c := range []struct {
		sql  string
		want string // the node type the item is
	}{
		{"select {D '2020-01-01'}", "string"},
		{"select {TS '2020-01-01'}", "string"},
		{"select {date '2020-01-01'}", "string"},
		{"select {abc 1}", "int"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			switch selectItem(t, c.sql).(type) {
			case *nodes.StringLit:
				if c.want != "string" {
					t.Errorf("Parse(%q) item is a string, want %s", c.sql, c.want)
				}
			case *nodes.IntLit:
				if c.want != "int" {
					t.Errorf("Parse(%q) item is an int, want %s", c.sql, c.want)
				}
			default:
				t.Errorf("Parse(%q) item is not the bare expression", c.sql)
			}
		})
	}

	t.Run("outer join escape", func(t *testing.T) {
		list, err := Parse("select * from {oj t1 left join t2 on t1.a = t2.a}, t3")
		if err != nil {
			t.Fatal(err)
		}
		from := list.Items[0].(*nodes.SelectStmt).From
		if len(from) != 2 {
			t.Fatalf("From has %d items, want 2", len(from))
		}
		if j, ok := from[0].(*nodes.JoinClause); !ok || j.Type != nodes.JoinLeft {
			t.Errorf("From[0] = %#v, want the LEFT JOIN", from[0])
		}
	})
}
