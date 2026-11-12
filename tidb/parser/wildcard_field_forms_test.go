package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// Where a wildcard may stand, as TiDB v8.5.5 reads it, every verdict measured with
// PREPARE on pingcap/tidb:v8.5.5 on 2026-11-12.
//
// `*`, `tbl.*` and `db.tbl.*` are Field alternatives of their own (parser.y:6463-6473)
// and no Expression, so a select list is the only place one stands, and it stands there
// whole: no operator takes one as an operand, and no other clause, argument or statement
// reads one. COUNT(*) is COUNT's own production (parser.y:8546). parseSelectWildcard is
// the one reader; parseExpr has none.
//
// The reject-fuzz sweep found the first of these, `select * = t`, which omni read as a
// comparison.

var wildcardFieldCases = []tableSampleCase{
	// A field of a select list, in any position and any number of times.
	{"select * from t", true},
	{"select *, a from t", true},
	{"select a, * from t", true},
	{"select t.* from t", true},
	{"select s.t.* from t", true},
	{"select a, t.* from t", true},
	{"select *, t.*, a from t", true},
	{"select *from t", true},
	{"select t . * from t", true},
	{"select t. * from t", true},
	{"select t .* from t", true},
	{"select distinct * from t", true},
	{"select * from t union all select t.* from t", true},
	{"select (select * from t limit 1)", true},
	{"select exists (select * from t)", true},
	{"insert into t select * from t", true},
	{"with c as (select * from t) select c.* from c", true},
	// A backquoted star is a name, and a star between operands is a product.
	{"select `*` from t", true},
	{"select a * b from t", true},
	{"select t.a * 2 from t", true},
	// COUNT's own star, which an operator may then take.
	{"select count(*) from t", true},
	{"select count(*) = 1 from t", true},
	{"select count(*) over () from t", true},

	// No alias.
	{"select * as x from t", false},
	{"select * x from t", false},
	{"select t.* x from t", false},
	// No comparison.
	{"select * = t", false},
	{"select * = 1 from t", false},
	{"select *=1 from t", false},
	{"select t.* = 1 from t", false},
	{"select s.t.* = 1 from t", false},
	{"select 1 = * from t", false},
	{"select 1 = t.* from t", false},
	{"select * <=> 1 from t", false},
	// No other operator, on either side or in front.
	{"select * + 1 from t", false},
	{"select 1 + * from t", false},
	{"select 1 + t.* from t", false},
	{"select * * 2 from t", false},
	{"select 2 * * from t", false},
	{"select * div 2 from t", false},
	{"select - * from t", false},
	{"select -t.* from t", false},
	{"select not * from t", false},
	{"select * and 1 from t", false},
	{"select 1 and t.* from t", false},
	{"select * or 1 from t", false},
	{"select * is null from t", false},
	{"select t.* is null from t", false},
	{"select * in (1) from t", false},
	{"select * between 1 and 2 from t", false},
	{"select * like 'a' from t", false},
	{"select * regexp 'a' from t", false},
	{"select * collate utf8mb4_bin from t", false},
	{"select @a := * from t", false},
	// No place inside an expression.
	{"select (*) from t", false},
	{"select (t.*) from t", false},
	{"select case * when 1 then 2 end from t", false},
	{"select case when * then 2 end from t", false},
	{"select case when 1 then * end from t", false},
	{"select case when 1 then 2 else t.* end from t", false},
	{"select interval * day + d from t", false},
	{"select d + interval * day from t", false},
	{"select cast(* as char) from t", false},
	{"select cast(t.* as char) from t", false},
	{"select convert(*, char) from t", false},
	{"select convert(* using utf8mb4) from t", false},
	{"select count(t.*) from t", false},
	{"select count(distinct *) from t", false},
	{"select count(all *) from t", false},
	// No other clause.
	{"select 1 from t where *", false},
	{"select 1 from t where t.*", false},
	{"select 1 from t where a = *", false},
	{"select 1 from t where a = t.*", false},
	{"select 1 from t group by *", false},
	{"select 1 from t group by t.*", false},
	{"select 1 from t having *", false},
	{"select 1 from t order by *", false},
	{"select 1 from t order by t.*", false},
	// No other statement: DO takes an ExpressionList (parser.y DoStmt).
	{"do *", false},
	{"do t.*", false},
	{"do * = 1", false},
	{"do 1 = *", false},
	{"do 1, *", false},
	{"set @a = *", false},
	{"set @a = t.*", false},
	{"update t set a = *", false},
	{"update t set a = t.*", false},
	{"insert into t (a) values (*)", false},
	{"insert into t (a) values (t.*)", false},
	{"insert into t (a) values (1) on duplicate key update a = *", false},
}

// TestWildcardFieldForms holds omni to the recorded verdicts without a server.
func TestWildcardFieldForms(t *testing.T) {
	for _, c := range wildcardFieldCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// TestWildcardFieldTree holds the node each wildcard is: the same StarExpr and
// star ColumnRef the expression parser gave before the select list read them itself.
func TestWildcardFieldTree(t *testing.T) {
	for _, c := range []struct {
		sql                 string
		schema, table, text string
	}{
		{"select * from t", "", "", "*"},
		{"select t.* from t", "", "t", "t.*"},
		{"select s.t.* from t", "s", "t", "s.t.*"},
		{"select `s` . `t` . * from t", "s", "t", "`s` . `t` . *"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.sql, err)
			}
			sel, ok := list.Items[0].(*nodes.SelectStmt)
			if !ok || len(sel.TargetList) != 1 {
				t.Fatalf("Parse(%q) = %T with %d targets, want a SELECT with one", c.sql, list.Items[0], len(sel.TargetList))
			}
			var loc nodes.Loc
			switch n := sel.TargetList[0].(type) {
			case *nodes.StarExpr:
				if c.table != "" {
					t.Fatalf("target is a bare star, want %s.*", c.table)
				}
				loc = n.Loc
			case *nodes.ColumnRef:
				if !n.Star || n.Schema != c.schema || n.Table != c.table || n.Column != "" {
					t.Fatalf("target is %+v, want the wildcard of %q.%q", n, c.schema, c.table)
				}
				loc = n.Loc
			default:
				t.Fatalf("target is %T, want a wildcard", n)
			}
			if got := c.sql[loc.Start:loc.End]; got != c.text {
				t.Errorf("wildcard spans %q, want %q", got, c.text)
			}
		})
	}
}
