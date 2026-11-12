package parser

import "testing"

// Where the bare `*` may stand, as OceanBase 4.3.5 reads it, every verdict measured
// with PREPARE on oceanbase-ce:4.3.5.6 on 2026-11-12.
//
// `*` is a projection of its own (the '*' arm of projection) and no expr, so a select
// list is the only place it stands, and it stands there whole: no operator takes it as
// an operand, and no other clause or argument reads it. COUNT(*) is COUNT's own
// production. parseBareStar is the one reader; parseExpr has none.
//
// `tbl.*` is different here. It is an arm of column_ref, so it is an expr like any other
// column and parses wherever one does; the resolver refuses it afterwards with 1149,
// where MySQL 8.0, MariaDB 11.8 and TiDB v8.5.5 answer 1064.
//
// TiDB's reject-fuzz sweep found the first of these, `select * = t`, which every
// MySQL-family parser here read as a comparison.
//
// Ref: projection and column_ref in OceanBase 4.3.5 sql_parser_mysql_mode.y.

var bareStarCases = []parserFormCase{
	// An item of a select list, in any position.
	{"select * from t", true},
	{"select *, a from t", true},
	{"select a, * from t", true},
	{"select t.*, * from t", true},
	{"select *from t", true},
	{"select distinct * from t", true},
	{"select * from t union all select * from t", true},
	{"select (select * from t limit 1)", true},
	{"select exists (select * from t)", true},
	{"insert into t select * from t", true},
	// A backquoted star is a name, and a star between operands is a product.
	{"select `*` from t", true},
	{"select a * b from t", true},
	// COUNT's own star, which an operator may then take.
	{"select count(*) from t", true},
	{"select count(all *) from t", true},
	{"select count(*) = 1 from t", true},
	// `tbl.*` is a column_ref, so the parser takes it as an operand and anywhere else.
	{"select t.* from t", true},
	{"select a, t.* from t", true},
	{"select t.* = 1 from t", true},
	{"select 1 + t.* from t", true},
	{"select -t.* from t", true},
	{"select (t.*) from t", true},
	{"select count(t.*) from t", true},
	{"select 1 from t where t.*", true},
	{"select 1 from t order by t.*", true},
	{"update t set a = t.*", true},
	{"insert into t (a) values (t.*)", true},
	// DO reads a bare star as well: the server answers 1149 for these, which its SQL
	// parser does not raise.
	{"do *", true},
	{"do 1, *", true},

	// No alias, on either wildcard.
	{"select * as x from t", false},
	{"select * x from t", false},
	{"select t.* x from t", false},
	// No comparison.
	{"select * = t", false},
	{"select * = 1 from t", false},
	{"select *=1 from t", false},
	{"select 1 = * from t", false},
	{"select * <=> 1 from t", false},
	// No other operator, on either side or in front.
	{"select * + 1 from t", false},
	{"select 1 + * from t", false},
	{"select * * 2 from t", false},
	{"select 2 * * from t", false},
	{"select * div 2 from t", false},
	{"select - * from t", false},
	{"select not * from t", false},
	{"select * and 1 from t", false},
	{"select * or 1 from t", false},
	{"select * is null from t", false},
	{"select * in (1) from t", false},
	{"select * between 1 and 2 from t", false},
	{"select * like 'a' from t", false},
	{"select * regexp 'a' from t", false},
	{"select * collate utf8mb4_bin from t", false},
	{"select @a := * from t", false},
	// No place inside an expression.
	{"select (*) from t", false},
	{"select case * when 1 then 2 end from t", false},
	{"select case when * then 2 end from t", false},
	{"select case when 1 then * end from t", false},
	{"select interval * day + d from t", false},
	{"select d + interval * day from t", false},
	{"select cast(* as char) from t", false},
	{"select convert(*, char) from t", false},
	{"select convert(* using utf8mb4) from t", false},
	{"select count(distinct *) from t", false},
	// No other clause.
	{"select 1 from t where *", false},
	{"select 1 from t where a = *", false},
	{"select 1 from t group by *", false},
	{"select 1 from t having *", false},
	{"select 1 from t order by *", false},
	// No other statement.
	{"do * = 1", false},
	{"do 1 = *", false},
	{"set @a = *", false},
	{"update t set a = *", false},
	{"insert into t (a) values (*)", false},
	{"insert into t (a) values (1) on duplicate key update a = *", false},
}

// TestBareStarForms holds omni to the recorded verdicts without a server.
func TestBareStarForms(t *testing.T) {
	for _, c := range bareStarCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}
