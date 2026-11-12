package parser

import "testing"

// Where a wildcard may stand, as MariaDB 11.8 reads it, every verdict measured with
// PREPARE on mariadb:11.8.8 on 2026-11-12.
//
// `*`, `tbl.*` and `db.tbl.*` are select items of their own (the '*' arm of
// select_item_list, and table_wild) and no expr, so a select list, RETURNING's included, is the only place one
// stands, and it stands there whole: no operator takes one as an operand, and no other
// clause or argument reads one. COUNT(*) is COUNT's own production. parseSelectWildcard
// is the one reader; parseExpr has none.
//
// TiDB's reject-fuzz sweep found the first of these, `select * = t`, which every
// MySQL-family parser here read as a comparison.
//
// Ref: select_item_list, select_item and table_wild in MariaDB 11.8 sql/sql_yacc.yy.

var wildcardItemCases = []parserFormCase{
	// An item of a select list.
	{"select * from t", true},
	{"select *, a from t", true},
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
	// A qualifier is a name, a word the expression reader has another use for
	// included: OFF, which is false as a value.
	{"select off.* from t", true},
	{"select off .* from t", true},
	// A backquoted star is a name, and a star between operands is a product.
	{"select `*` from t", true},
	{"select a * b from t", true},
	{"select t.a * 2 from t", true},
	// COUNT's own star, which an operator may then take.
	{"select count(*) from t", true},
	{"select count(all *) from t", true},
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
	// No other statement.
	{"set @a = *", false},
	{"set @a = t.*", false},
	{"update t set a = *", false},
	{"update t set a = t.*", false},
	{"insert into t (a) values (*)", false},
	{"insert into t (a) values (t.*)", false},
	{"insert into t (a) values (1) on duplicate key update a = *", false},
	// RETURNING takes a select list too.
	{"insert into t (a) values (1) returning *", true},
	{"insert into t (a) values (1) returning t.*", true},
	{"delete from t returning *", true},
	{"replace into t (a) values (1) returning *", true},
	{"insert into t (a) values (1) returning * = 1", false},
	{"insert into t (a) values (1) returning t.* + 1", false},
	// DO takes an expr_list.
	{"do *", false},
	{"do t.*", false},
	{"do * = 1", false},
	{"do 1 = *", false},
	{"do *, 1", false},
}

// wildcardItemGaps are the wildcards omni still reads where MariaDB 11.8 does not:
// select_item_list takes the bare `*` as its first item only, and omni takes it in any
// position. The oracle test fails when one starts to agree, so an entry cannot outlive
// its gap.
var wildcardItemGaps = []string{
	"select a, * from t",
	"select t.*, * from t",
	"select *, * from t",
	"select count(*), * from t",
}

// TestWildcardItemForms holds omni to the recorded verdicts without a server.
func TestWildcardItemForms(t *testing.T) {
	for _, c := range wildcardItemCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
	for _, sql := range wildcardItemGaps {
		if _, err := Parse(sql); err != nil {
			t.Errorf("omni now rejects %q, as MariaDB 11.8 does: move it into wildcardItemCases (%v)", sql, err)
		}
	}
}
