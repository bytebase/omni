package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/deparse"
)

// ============================================================================
// The temporal functions: what goes inside their parentheses (see temporal.go).
//
// temporalFormCases is the recorded expectation, and TestMariaDBTemporalFormsMatchEngine holds every row
// of it to a live mariadb:11.8.8. The tree tests carry the rest: an accept/reject gate
// cannot see a precision or an INTERVAL that parses and is then dropped.
// ============================================================================

var temporalFormCases = []struct {
	sql    string
	accept bool
}{
	// --- precision shapes ---
	{"select current_timestamp", true},
	{"select current_timestamp()", true},
	{"select current_timestamp(0)", true},
	{"select current_timestamp(6)", true},
	{"select current_timestamp(7)", true},
	{"select current_timestamp(255)", true},
	{"select current_timestamp(256)", true},
	{"select current_timestamp(2147483647)", true},
	{"select current_timestamp(2147483648)", true},
	{"select current_timestamp(4294967296)", true},
	{"select current_timestamp(18446744073709551615)", true},
	{"select current_timestamp(18446744073709551616)", false},
	{"select current_timestamp(00000000000000000000001)", true},
	{"select current_timestamp(-1)", false},
	{"select current_timestamp(+1)", false},
	{"select current_timestamp(1.0)", false},
	{"select current_timestamp(1e0)", false},
	{"select current_timestamp(.5)", false},
	{"select current_timestamp('2')", false},
	{"select current_timestamp(null)", false},
	{"select current_timestamp(true)", false},
	{"select current_timestamp(0x1)", true},
	{"select current_timestamp(x'01')", false},
	{"select current_timestamp(b'1')", false},
	{"select current_timestamp(b)", false},
	{"select current_timestamp(@a)", false},
	{"select current_timestamp(?)", false},
	{"select current_timestamp((1))", false},
	{"select current_timestamp(1+1)", false},
	{"select current_timestamp(1,2)", false},
	{"select current_timestamp(1,)", false},
	{"select current_timestamp(", false},
	{"select current_timestamp (1)", true},
	{"select current_timestamp /* c */ (1)", true},
	{"select current_timestamp (-1)", false},
	{"select current_timestamp.a from t current_timestamp", false},
	{"select current_timestamp a", true},
	{"select current_timestamp as a", true},
	{"select 1 as current_timestamp", false},
	{"select 1 current_timestamp", false},
	{"select `current_timestamp` from t", true},
	{"select `current_timestamp`(-1)", true},
	{"select tp.current_timestamp(-1)", true},
	{"select t.current_timestamp from t", true},
	{"select current_timestamp() over ()", false},
	{"select current_timestamp(1) over ()", false},
	// --- current_time ---
	{"select current_time", true},
	{"select current_time()", true},
	{"select current_time(6)", true},
	{"select current_time(7)", true},
	{"select current_time(-1)", false},
	{"select current_time(1.0)", false},
	{"select current_time('1')", false},
	{"select current_time(null)", false},
	{"select current_time(1,2)", false},
	{"select current_time(0x1)", true},
	{"select current_time (1)", true},
	{"select current_time (-1)", false},
	{"select 1 as current_time", false},
	{"select `current_time`(-1)", true},
	{"select tp.current_time(-1)", true},
	// --- curtime ---
	{"select curtime", true},
	{"select curtime()", true},
	{"select curtime(6)", true},
	{"select curtime(7)", true},
	{"select curtime(-1)", false},
	{"select curtime(1.0)", false},
	{"select curtime('1')", false},
	{"select curtime(null)", false},
	{"select curtime(1,2)", false},
	{"select curtime(0x1)", true},
	{"select curtime(b)", false},
	{"select curtime (1)", true},
	{"select curtime (-1)", true},
	{"select curtime (b)", true},
	{"select curtime /**/ (-1)", true},
	{"select curtime from t", true},
	{"select curtime.a from curtime", true},
	{"select t.curtime from t", true},
	{"select 1 as curtime", true},
	{"select 1 curtime", true},
	{"select `curtime`(-1)", true},
	{"select tp.curtime(-1)", true},
	{"select curtime(1) over ()", false},
	{"select curtime() over ()", false},
	{"select curtime(255)", true},
	{"select curtime(256)", true},
	{"select curtime(2147483648)", true},
	{"select curtime(18446744073709551616)", false},
	// --- utc_timestamp ---
	{"select utc_timestamp", true},
	{"select utc_timestamp()", true},
	{"select utc_timestamp(6)", true},
	{"select utc_timestamp(7)", true},
	{"select utc_timestamp(-1)", false},
	{"select utc_timestamp(1.0)", false},
	{"select utc_timestamp('1')", false},
	{"select utc_timestamp(null)", false},
	{"select utc_timestamp(1,2)", false},
	{"select utc_timestamp (1)", true},
	{"select 1 as utc_timestamp", false},
	{"select `utc_timestamp`(-1)", true},
	{"select tp.utc_timestamp(-1)", true},
	// --- utc_time ---
	{"select utc_time", true},
	{"select utc_time()", true},
	{"select utc_time(6)", true},
	{"select utc_time(7)", true},
	{"select utc_time(-1)", false},
	{"select utc_time(1.0)", false},
	{"select utc_time('1')", false},
	{"select utc_time(null)", false},
	{"select utc_time(1,2)", false},
	{"select utc_time (1)", true},
	{"select 1 as utc_time", false},
	// --- localtime, localtimestamp ---
	{"select localtime", true},
	{"select localtime()", true},
	{"select localtime(6)", true},
	{"select localtime(7)", true},
	{"select localtime(-1)", false},
	{"select localtime(1.0)", false},
	{"select localtime('1')", false},
	{"select localtime(null)", false},
	{"select localtimestamp", true},
	{"select localtimestamp()", true},
	{"select localtimestamp(6)", true},
	{"select localtimestamp(-1)", false},
	{"select localtimestamp(1.0)", false},
	{"select localtimestamp('1')", false},
	{"select localtimestamp(null)", false},
	{"select 1 as localtime", false},
	// --- now ---
	{"select now", true},
	{"select now()", true},
	{"select now(0)", true},
	{"select now(6)", true},
	{"select now(7)", true},
	{"select now(255)", true},
	{"select now(256)", true},
	{"select now(-1)", false},
	{"select now(1.0)", false},
	{"select now('1')", false},
	{"select now(null)", false},
	{"select now(b)", false},
	{"select now(1,2)", false},
	{"select now(0x1)", true},
	{"select now (1)", true},
	{"select now (-1)", true},
	{"select now(*)", false},
	{"select now(distinct 1)", false},
	{"select now from t", true},
	{"select now.a from now", true},
	{"select t.now from t", true},
	{"select 1 as now", true},
	{"select 1 now", true},
	{"select `now`(-1)", true},
	{"select tp.now(-1)", true},
	{"select now() over ()", false},
	{"select now(1) over ()", false},
	// --- sysdate ---
	{"select sysdate", true},
	{"select sysdate()", true},
	{"select sysdate(6)", true},
	{"select sysdate(7)", true},
	{"select sysdate(-1)", false},
	{"select sysdate(1.0)", false},
	{"select sysdate('1')", false},
	{"select sysdate(null)", false},
	{"select sysdate(1,2)", false},
	{"select sysdate(0x1)", true},
	{"select sysdate (-1)", false},
	{"select sysdate from t", true},
	{"select 1 as sysdate", true},
	{"select tp.sysdate(-1)", true},
	{"select `sysdate`(-1)", true},
	{"select sysdate() over ()", false},
	// --- curdate, current_date, utc_date ---
	{"select curdate", true},
	{"select curdate()", true},
	{"select curdate(1)", false},
	{"select curdate(0)", false},
	{"select curdate(null)", false},
	{"select curdate(b)", false},
	{"select curdate (1)", true},
	{"select curdate (b)", true},
	{"select curdate from t", true},
	{"select 1 as curdate", true},
	{"select `curdate`(1)", true},
	{"select tp.curdate(1)", true},
	{"select curdate() over ()", false},
	{"select current_date", true},
	{"select current_date()", true},
	{"select current_date(1)", false},
	{"select current_date(null)", false},
	{"select current_date (1)", false},
	{"select current_date ()", true},
	{"select 1 as current_date", false},
	{"select tp.current_date(1)", true},
	{"select `current_date`(1)", true},
	{"select current_date() over ()", false},
	{"select utc_date", true},
	{"select utc_date()", true},
	{"select utc_date(1)", false},
	{"select utc_date(null)", false},
	{"select utc_date (1)", false},
	{"select 1 as utc_date", false},
	// --- current_user ---
	{"select current_user", true},
	{"select current_user()", true},
	{"select current_user(1)", false},
	{"select current_user (1)", false},
	{"select current_role", true},
	{"select current_role()", true},
	// --- corpus row shape ---
	{"SELECT CURRENT_DATE, CURRENT_DATE(), CURDATE(1)", false},
	{"SELECT CURRENT_DATE, CURRENT_DATE(), CURDATE()", true},
	// --- date_add ---
	{"select date_add", true},
	{"select date_add from t", true},
	{"select date_add()", false},
	{"select date_add(a)", false},
	{"select date_add(a, 10)", false},
	{"select date_add(a, 0.10)", false},
	{"select date_add(a, \"11,11\")", false},
	{"select date_add(a, b)", false},
	{"select date_add(a, null)", false},
	{"select date_add(a, interval 1 day)", true},
	{"select date_add(a, INTERVAL 1 SECOND_MICROSECOND)", true},
	{"select date_add(a, interval b day)", true},
	{"select date_add(a, interval 1+1 day)", true},
	{"select date_add(a, interval (1) day)", true},
	{"select date_add(a, interval interval(1,2) day)", true},
	{"select date_add(a, interval 1)", false},
	{"select date_add(a, interval 1 foo)", false},
	{"select date_add(a, interval 1 day + 1)", false},
	{"select date_add(a, interval 1 day + interval 1 day)", false},
	{"select date_add(a, 1 + interval 1 day)", false},
	{"select date_add(a, (interval 1 day))", false},
	{"select date_add(a, b + interval 1 day)", false},
	{"select date_add(a, interval 1 day, 1)", false},
	{"select date_add(a, 1, 2)", false},
	{"select date_add(a, interval 1 sql_tsi_day)", true},
	{"select date_add(interval 1 day + a, interval 1 day)", true},
	{"select date_add(a + interval 1 day, interval 1 day)", true},
	{"select date_add(interval 1 day, interval 1 day)", false},
	{"select date_add(*, interval 1 day)", false},
	{"select date_add(distinct a, interval 1 day)", false},
	{"select date_add(all a, interval 1 day)", false},
	{"select date_add(a, interval 1 day) over ()", false},
	{"select date_add (a, 10)", true},
	{"select date_add (a, b)", true},
	{"select date_add (a, interval 1 day)", false},
	{"select date_add /**/ (a, 10)", true},
	{"select date_add.a from date_add", true},
	{"select t.date_add from t", true},
	{"select 1 as date_add", true},
	{"select 1 date_add", true},
	{"select `date_add`(a, 10)", true},
	{"select `date_add`(a, interval 1 day)", false},
	{"select tp.date_add(a, 10)", true},
	{"select tp.date_add(a, interval 1 day)", false},
	{"select date_add(\"2011-11-11 10:10:10.123456\", 10)", false},
	{"select date_add(\"2011-11-11 10:10:10.123456\", 0.10)", false},
	{"select date_add(\"2011-11-11 10:10:10.123456\", \"11,11\")", false},
	{"select date_add(\"2011-11-11 10:10:10.123456\", interval 10 microsecond)", true},
	// --- date_sub ---
	{"select date_sub", true},
	{"select date_sub from t", true},
	{"select date_sub()", false},
	{"select date_sub(a)", false},
	{"select date_sub(a, 10)", false},
	{"select date_sub(a, 0.10)", false},
	{"select date_sub(a, \"11,11\")", false},
	{"select date_sub(a, b)", false},
	{"select date_sub(a, null)", false},
	{"select date_sub(a, interval 1 day)", true},
	{"select date_sub(a, interval 1)", false},
	{"select date_sub(a, interval 1 day + 1)", false},
	{"select date_sub(a, 1 + interval 1 day)", false},
	{"select date_sub(a, (interval 1 day))", false},
	{"select date_sub(a, interval 1 day, 1)", false},
	{"select date_sub(a, interval 1 day) over ()", false},
	{"select date_sub (a, 10)", true},
	{"select date_sub (a, interval 1 day)", false},
	{"select 1 as date_sub", true},
	{"select `date_sub`(a, 10)", true},
	{"select tp.date_sub(a, 10)", true},
	// --- adddate, subdate ---
	{"select adddate", true},
	{"select adddate from t", true},
	{"select adddate()", false},
	{"select adddate(a)", false},
	{"select adddate(a, 10)", true},
	{"select adddate(a, 0.10)", true},
	{"select adddate(a, \"11,11\")", true},
	{"select adddate(a, b)", true},
	{"select adddate(a, interval 1 day)", true},
	{"select adddate(a, interval 1)", false},
	{"select adddate(a, interval 1 day + 1)", true},
	{"select adddate(a, b + interval 1 day)", true},
	{"select adddate(a, (interval 1 day))", false},
	{"select adddate(a, 1, 2)", false},
	{"select adddate(a, interval 1 day) over ()", false},
	{"select adddate (a, 10)", true},
	{"select adddate (a)", true},
	{"select adddate (a, interval 1 day)", false},
	{"select 1 as adddate", true},
	{"select `adddate`(a, 10)", true},
	{"select subdate(a)", false},
	{"select subdate(a, 10)", true},
	{"select subdate(a, interval 1 day)", true},
	{"select subdate(a, 1, 2)", false},
	{"select subdate(a, interval 1 day + 1)", true},
	{"select subdate (a, interval 1 day)", false},
	// --- spaced empty ---
	{"select now ()", true},
	{"select curdate ()", true},
	{"select curtime ()", true},
	{"select sysdate ()", true},
	{"select sysdate (1)", true},
	{"select sysdate (b)", false},
	{"select sysdate (1,2)", false},
	{"select sysdate (1) over ()", false},
	{"select now (1) over ()", false},
	{"select curtime (1) over ()", false},
	{"select date_add ()", true},
	{"select adddate ()", true},
	{"select subdate ()", true},
	{"select adddate (a, b, c)", true},
	{"select curdate/**/(1)", true},
	{"select now/**/(1)", true},
	{"select CURDATE (1)", true},
	{"select cUrTiMe(-1)", false},
	{"select t.current_time(-1)", true},
	{"select t.now(-1)", true},
	{"select t.curdate(1)", true},
	// --- interval corners ---
	{"select date_add(a, interval(1,2))", false},
	{"select date_sub(a, interval(1,2))", false},
	{"select adddate(a, interval(1,2))", true},
	{"select subdate(a, interval(1,2))", true},
	{"select date_add((interval 1 day), interval 1 day)", false},
	{"select adddate(interval 1 day, 1)", false},
	{"select adddate((interval 1 day), 1)", false},
	{"select adddate(interval 1 day + a, interval 1 day)", true},
	{"select adddate(*, 1)", false},
	{"select adddate(a, *)", false},
	{"select adddate(distinct a, 1)", false},
	{"select adddate(a, distinct 1)", false},
	{"select date_add(a, distinct interval 1 day)", false},
	{"select date_add(a, interval 1 day,)", false},
	{"select abs(interval 1 day)", false},
	{"select concat(interval 1 day)", false},
	{"select coalesce(interval 1 day, 1)", false},
	{"select interval 1 day + 1", true},
	{"select 1 + interval 1 day", true},
	{"select abs((interval 1 day))", false},
	{"select abs(interval 1 day + 1)", true},
	// --- number limits ---
	{"select current_time(0xffffffffffffffffffff)", true},
	{"select current_time(0X1)", false},
	{"select current_time(0b1)", false},
	{"select now(0b1)", false},
	{"select curtime(0xa)", true},
	{"select utc_time(0x1)", true},
	{"select current_time(0002147483647)", true},
	{"select current_time(2147483647)", true},
	{"select current_time(18446744073709551615)", true},
	{"select current_time(18446744073709551616)", false},
	{"select now(2147483647)", true},
	{"select now(2147483648)", true},
	{"select now(18446744073709551615)", true},
	{"select now(18446744073709551616)", false},
	{"select now(99999999999999999999999)", false},
	{"select curtime(2147483647)", true},
	{"select curtime(18446744073709551615)", true},
	{"select sysdate(2147483647)", true},
	{"select sysdate(2147483648)", true},
	{"select sysdate(18446744073709551615)", true},
	{"select sysdate(18446744073709551616)", false},
	{"select utc_time(2147483648)", true},
	{"select utc_timestamp(2147483648)", true},
	{"select localtime(2147483648)", true},
	{"select localtimestamp(2147483648)", true},
	// --- contexts ---
	{"create table t1 (a timestamp default current_timestamp(6) on update current_timestamp(6))", true},
	{"create table t2 (a timestamp default current_timestamp(-1))", false},
	{"create table t3 (a timestamp default now(6) on update now(6))", true},
	{"create table t4 (a timestamp default now(-1))", false},
	{"create table t5 (a timestamp on update current_timestamp(1.0))", false},
	{"create table t6 (a date default curdate(1))", false},
	{"create table t7 (a date default (curdate(1)))", false},
	{"create table t8 (a datetime default (date_add(now(), 1)))", false},
	{"create table t9 (a datetime default (date_add(now(), interval 1 day)))", true},
	{"create table t10 (a datetime default current_timestamp(7))", true},
	{"create table t11 (a datetime default now(7))", true},
	// --- leading zeros, trailing tokens ---
	{"select current_timestamp(00000000000000000000000018446744073709551615)", true},
	{"select current_timestamp(00000000000000000000000018446744073709551616)", false},
	{"select current_timestamp(0000000000002147483647)", true},
	{"select current_timestamp(0000000000002147483648)", true},
	{"select now(0000000000002147483648)", true},
	{"select current_time(x'0xaa')", false},
	{"select now(x'0xaa')", false},
	{"select current_date (", false},
	{"select now(1)(2)", false},
	{"select curtime(1)(2)", false},
	{"select current_time(1)(2)", false},
	{"select now(1) x", true},
	{"select curtime(1) as x from t", true},
	{"select now() + interval 1 day", true},
	{"select curdate() - interval 1 day", true},
	{"select date_add(a, interval 1 day) x from t", true},
	{"select date_sub(a, interval 1 day) as x from t", true},
	{"select date_add(a, interval 1 day) + 1 from t", true},
	{"select date_add(a,", false},
	{"select date_add(a", false},
	{"select date_add(a, interval 1 day", false},
	{"select date_add(, interval 1 day)", false},
	{"select adddate(a,)", false},
	{"select adddate(,a)", false},
	{"select adddate(a, 1", false},
	{"select interval(interval 1 day, 1)", false},
	{"select if(interval 1 day, 1, 2)", false},
	{"select left(interval 1 day, 1)", false},
	{"select count(interval 1 day)", false},
	{"select sum(interval 1 day)", false},
	{"select max(a + interval 1 day) from t", true},
	{"select date_format(now(), '%Y')", true},
	{"select date_format(a + interval 1 day, '%Y') from t", true},
	{"select extract(day from a + interval 1 day) from t", true},
	{"select timestampadd(day, 1, now(6))", true},
	{"select now(now())", false},
	{"select curtime(now())", false},
	{"select date_add(now(6), interval 1 day)", true},
	{"select date_add(date_add(a, interval 1 day), interval 1 day) from t", true},
	{"select date_add(a, interval date_add(b, interval 1 day) day) from t", true},
	{"select adddate(adddate(a, 1), 1) from t", true},
	{"select date_sub(curdate(), interval 1 month)", true},
	{"select sysdate(3) + interval 1 day", true},
	{"select utc_timestamp(3) + 1", true},
	{"update t set a = now(6)", true},
	{"update t set a = now(-1)", false},
	{"update t set a = date_add(a, interval 1 day)", true},
	{"update t set a = date_add(a, 1)", false},
	{"delete from t where a < date_sub(now(), interval 1 day)", true},
	{"delete from t where a < date_sub(now(), 1)", false},
	{"insert into t (a) values (current_timestamp(3))", true},
	{"insert into t (a) values (current_timestamp(-1))", false},
	{"select * from t where a > curdate() order by now(6)", true},
	{"select * from t where a > curdate(1)", false},
	{"select bit_and(interval 1 day)", false},
	{"select json_arrayagg(interval 1 day)", false},
	{"select char(interval 1 day)", false},
}

func TestMariaDBTemporalForms(t *testing.T) {
	for _, c := range temporalFormCases {
		_, err := Parse(c.sql)
		if (err == nil) != c.accept {
			t.Errorf("Parse(%q) err = %v, want accept = %v", c.sql, err, c.accept)
		}
	}
}

// temporalCall parses sql, whose only select item must be a call, and holds its
// name, its argument count and its Loc.
func temporalCall(t *testing.T, sql, name string, args int) *nodes.FuncCallExpr {
	t.Helper()
	fc, ok := firstExpr(t, sql).(*nodes.FuncCallExpr)
	if !ok {
		t.Fatalf("%q: not a FuncCallExpr", sql)
	}
	if fc.Name != name {
		t.Errorf("%q: Name = %q, want %q", sql, fc.Name, name)
	}
	if len(fc.Args) != args {
		t.Fatalf("%q: Args = %#v, want %d of them", sql, fc.Args, args)
	}
	if fc.Over != nil {
		t.Errorf("%q: Over = %#v, want none", sql, fc.Over)
	}
	if want := strings.Index(sql, " ") + 1; fc.Loc.Start != want || fc.Loc.End != len(sql) {
		t.Errorf("%q: Loc = %+v, want {%d %d}", sql, fc.Loc, want, len(sql))
	}
	return fc
}

func isIntLit(expr nodes.ExprNode, want int64) bool {
	lit, ok := expr.(*nodes.IntLit)
	return ok && lit.Value == want
}

func isColumnRef(expr nodes.ExprNode, want string) bool {
	ref, ok := expr.(*nodes.ColumnRef)
	return ok && ref.Column == want
}

// TestMariaDBTemporalPrecisionTree asserts a precision reaches the tree as written,
// that an empty pair of parentheses is told from none, and that a precision above
// the maximum, which only fails later, is kept.
func TestMariaDBTemporalPrecisionTree(t *testing.T) {
	for _, c := range []struct {
		sql       string
		name      string
		hasParens bool
		precision int64 // -1 for none
	}{
		{"SELECT current_timestamp", "CURRENT_TIMESTAMP", false, -1},
		{"SELECT current_timestamp()", "CURRENT_TIMESTAMP", true, -1},
		{"SELECT current_timestamp(6)", "CURRENT_TIMESTAMP", true, 6},
		{"SELECT current_timestamp (6)", "CURRENT_TIMESTAMP", true, 6},
		{"SELECT current_timestamp(7)", "CURRENT_TIMESTAMP", true, 7},
		{"SELECT current_timestamp(000001)", "CURRENT_TIMESTAMP", true, 1},
		{"SELECT current_time", "CURRENT_TIME", false, -1},
		{"SELECT current_time(0)", "CURRENT_TIME", true, 0},
		{"SELECT localtime(3)", "LOCALTIME", true, 3},
		{"SELECT localtimestamp(3)", "LOCALTIMESTAMP", true, 3},
		{"SELECT utc_time", "UTC_TIME", false, -1},
		{"SELECT utc_time(6)", "UTC_TIME", true, 6},
		{"SELECT utc_timestamp", "UTC_TIMESTAMP", false, -1},
		{"SELECT utc_timestamp()", "UTC_TIMESTAMP", true, -1},
		{"SELECT utc_timestamp(6)", "UTC_TIMESTAMP", true, 6},
		{"SELECT now()", "NOW", true, -1},
		{"SELECT now(6)", "NOW", true, 6},
		{"SELECT now(7)", "NOW", true, 7},
		{"SELECT curtime()", "CURTIME", true, -1},
		{"SELECT curtime(6)", "CURTIME", true, 6},
		{"SELECT sysdate()", "SYSDATE", true, -1},
		{"SELECT sysdate(6)", "SYSDATE", true, 6},
		{"SELECT curdate()", "CURDATE", true, -1},
		{"SELECT current_date", "CURRENT_DATE", false, -1},
		{"SELECT current_date()", "CURRENT_DATE", true, -1},
		{"SELECT current_date ()", "CURRENT_DATE", true, -1},
		{"SELECT utc_date", "UTC_DATE", false, -1},
		{"SELECT utc_date()", "UTC_DATE", true, -1},
		{"SELECT current_user", "CURRENT_USER", false, -1},
		{"SELECT current_user()", "CURRENT_USER", true, -1},
	} {
		args := 0
		if c.precision >= 0 {
			args = 1
		}
		fc := temporalCall(t, c.sql, c.name, args)
		if fc.HasParens != c.hasParens {
			t.Errorf("%q: HasParens = %v, want %v", c.sql, fc.HasParens, c.hasParens)
		}
		if c.precision >= 0 && !isIntLit(fc.Args[0], c.precision) {
			t.Errorf("%q: Args[0] = %#v, want %d", c.sql, fc.Args[0], c.precision)
		}
	}
}

// TestMariaDBTemporalDateArithTree asserts both arguments of the date arithmetic
// builtins reach the tree, with the INTERVAL's value and unit.
func TestMariaDBTemporalDateArithTree(t *testing.T) {
	for _, c := range []struct {
		sql, name string
		value     int64
		unit      string
	}{
		{"SELECT date_add(a, INTERVAL 1 DAY)", "DATE_ADD", 1, "DAY"},
		{"SELECT date_sub(a, INTERVAL 2 HOUR)", "DATE_SUB", 2, "HOUR"},
		{"SELECT adddate(a, INTERVAL 3 MONTH)", "ADDDATE", 3, "MONTH"},
		{"SELECT subdate(a, INTERVAL 4 SECOND_MICROSECOND)", "SUBDATE", 4, "SECOND_MICROSECOND"},
	} {
		fc := temporalCall(t, c.sql, c.name, 2)
		if !isColumnRef(fc.Args[0], "a") {
			t.Errorf("%q: Args[0] = %#v, want the column a", c.sql, fc.Args[0])
		}
		iv, ok := fc.Args[1].(*nodes.IntervalExpr)
		if !ok {
			t.Fatalf("%q: Args[1] = %#v, want an IntervalExpr", c.sql, fc.Args[1])
		}
		if !isIntLit(iv.Value, c.value) || iv.Unit != c.unit {
			t.Errorf("%q: Args[1] = INTERVAL %#v %s, want INTERVAL %d %s", c.sql, iv.Value, iv.Unit, c.value, c.unit)
		}
		if want := strings.Index(c.sql, "INTERVAL"); iv.Loc.Start != want || iv.Loc.End != len(c.sql)-1 {
			t.Errorf("%q: INTERVAL Loc = %+v, want {%d %d}", c.sql, iv.Loc, want, len(c.sql)-1)
		}
	}

	// An INTERVAL may lead an expression in the first argument.
	fc := temporalCall(t, "SELECT date_add(INTERVAL 1 DAY + a, INTERVAL 2 DAY)", "DATE_ADD", 2)
	if _, ok := fc.Args[0].(*nodes.BinaryExpr); !ok {
		t.Errorf("date_add(INTERVAL 1 DAY + a, ...): Args[0] = %#v, want a BinaryExpr", fc.Args[0])
	}

	// ADDDATE and SUBDATE take any expression second, which the others do not.
	fc = temporalCall(t, "SELECT adddate(a, 10)", "ADDDATE", 2)
	if !isIntLit(fc.Args[1], 10) {
		t.Errorf("adddate(a, 10): Args[1] = %#v, want 10", fc.Args[1])
	}
	fc = temporalCall(t, "SELECT subdate(a, b)", "SUBDATE", 2)
	if !isColumnRef(fc.Args[1], "b") {
		t.Errorf("subdate(a, b): Args[1] = %#v, want the column b", fc.Args[1])
	}
	fc = temporalCall(t, "SELECT adddate(a, INTERVAL 1 DAY + 1)", "ADDDATE", 2)
	if _, ok := fc.Args[1].(*nodes.BinaryExpr); !ok {
		t.Errorf("adddate(a, INTERVAL 1 DAY + 1): Args[1] = %#v, want a BinaryExpr", fc.Args[1])
	}
	// The comma form is the INTERVAL function.
	fc = temporalCall(t, "SELECT adddate(a, INTERVAL(1, 2))", "ADDDATE", 2)
	if inner, ok := fc.Args[1].(*nodes.FuncCallExpr); !ok || inner.Name != "INTERVAL" || len(inner.Args) != 2 {
		t.Errorf("adddate(a, INTERVAL(1, 2)): Args[1] = %#v, want the INTERVAL function", fc.Args[1])
	}
}

// TestMariaDBTemporalEngineTree covers what only this engine reads.
func TestMariaDBTemporalEngineTree(t *testing.T) {
	// The function words are identifiers with a space or a comment before the '(',
	// so these are ordinary calls that take any expressions.
	if fc := temporalCall(t, "SELECT curtime (b)", "CURTIME", 1); !isColumnRef(fc.Args[0], "b") {
		t.Errorf("curtime (b): Args[0] = %#v, want the column b", fc.Args[0])
	}
	temporalCall(t, "SELECT now (1, 2)", "NOW", 2)
	temporalCall(t, "SELECT curdate (1)", "CURDATE", 1)
	if fc := temporalCall(t, "SELECT date_add (a, 10)", "DATE_ADD", 2); !isIntLit(fc.Args[1], 10) {
		t.Errorf("date_add (a, 10): Args[1] = %#v, want 10", fc.Args[1])
	}
	temporalCall(t, "SELECT adddate (a)", "ADDDATE", 1)
	// SYSDATE is a keyword wherever it stands, so spacing changes nothing.
	if fc := temporalCall(t, "SELECT sysdate (6)", "SYSDATE", 1); !isIntLit(fc.Args[0], 6) {
		t.Errorf("sysdate (6): Args[0] = %#v, want 6", fc.Args[0])
	}
	// HEX_NUM is a precision here, and only here.
	for _, sql := range []string{"SELECT now(0x1)", "SELECT current_time(0x1)", "SELECT curtime(0xa)"} {
		fc := firstFuncCall(t, sql)
		if len(fc.Args) != 1 {
			t.Fatalf("%q: Args = %#v, want one", sql, fc.Args)
		}
		if _, ok := fc.Args[0].(*nodes.HexLit); !ok {
			t.Errorf("%q: Args[0] = %#v, want a HexLit", sql, fc.Args[0])
		}
	}
}

// TestMariaDBTemporalWordsAsNames asserts the words that are not reserved stay names
// when no call follows, and a quoted or qualified call stays an ordinary one with
// whatever arguments it was given.
func TestMariaDBTemporalWordsAsNames(t *testing.T) {
	for _, word := range []string{"now", "curdate", "curtime", "sysdate", "date_add", "date_sub", "adddate", "subdate"} {
		sql := "SELECT " + word + " FROM t"
		if !isColumnRef(firstExpr(t, sql), word) {
			t.Errorf("%q: item = %#v, want the column %s", sql, firstExpr(t, sql), word)
		}
		sql = "SELECT t." + word + " FROM t"
		if !isColumnRef(firstExpr(t, sql), word) {
			t.Errorf("%q: item = %#v, want the column %s", sql, firstExpr(t, sql), word)
		}
		sql = "SELECT 1 AS " + word
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		rt, ok := list.Items[0].(*nodes.SelectStmt).TargetList[0].(*nodes.ResTarget)
		if !ok || rt.Name != word {
			t.Errorf("%q: item = %#v, want the alias %s", sql, list.Items[0], word)
		}
	}

	for _, c := range []struct {
		sql, schema string
	}{
		{"SELECT `date_add`(a, 10) FROM t", ""},
		{"SELECT `adddate`(a, 10) FROM t", ""},
		{"SELECT tp.date_add(a, 10) FROM t", "tp"},
		{"SELECT tp.curdate(1)", "tp"},
	} {
		fc := firstFuncCall(t, c.sql)
		if fc.Schema != c.schema || len(fc.Args) == 0 {
			t.Errorf("%q: Schema = %q, Args = %#v, want schema %q and its arguments", c.sql, fc.Schema, fc.Args, c.schema)
		}
	}
}

// TestMariaDBTemporalCompletion asserts the reader still offers candidates: what an
// argument list offers at its start, and INTERVAL where nothing else can stand.
func TestMariaDBTemporalCompletion(t *testing.T) {
	for _, sql := range []string{"SELECT now(", "SELECT date_add(", "SELECT adddate("} {
		if cs := Collect(sql, len(sql)); cs == nil || !cs.HasRule("columnref") {
			t.Errorf("Collect(%q) = %v, want the columnref rule", sql, cs)
		}
	}
	for _, sql := range []string{"SELECT date_add(a, ", "SELECT date_sub(a, "} {
		if cs := Collect(sql, len(sql)); cs == nil || !cs.HasToken(kwINTERVAL) {
			t.Errorf("Collect(%q) = %v, want INTERVAL", sql, cs)
		}
	}
	if sql := "SELECT adddate(a, "; !Collect(sql, len(sql)).HasRule("columnref") {
		t.Errorf("Collect(%q): want the columnref rule, since any expression can stand here", sql)
	}
}

// TestMariaDBTemporalDeparseRoundTrip asserts the stored text of every accepted
// select item parses back to the same text. A quoted or spaced name is left out:
// the tree keeps no record of either, so its text comes back as the builtin's own
// syntax, as `trim`(b, 'x') does (see TestFunctionArgumentDeparseRoundTrip).
func TestMariaDBTemporalDeparseRoundTrip(t *testing.T) {
	for _, c := range temporalFormCases {
		if !c.accept || !strings.HasPrefix(c.sql, "select ") || strings.HasPrefix(c.sql, "select * ") ||
			strings.ContainsAny(c.sql, "`") || temporalWordSpacedIn(c.sql) {
			continue
		}
		list, err := Parse(c.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", c.sql, err)
		}
		sel, ok := list.Items[0].(*nodes.SelectStmt)
		if !ok || len(sel.TargetList) == 0 {
			continue
		}
		first := deparse.Deparse(targetExpr(sel.TargetList[0]))
		again, err := Parse("SELECT " + first)
		if err != nil {
			t.Errorf("%q: deparsed to %q, which does not parse: %v", c.sql, first, err)
			continue
		}
		second := deparse.Deparse(targetExpr(again.Items[0].(*nodes.SelectStmt).TargetList[0]))
		if first != second {
			t.Errorf("%q: deparse is not stable: %q then %q", c.sql, first, second)
		}
	}
}

// temporalWordSpacedIn reports whether a function word in sql has a space or a
// comment before its '('.
func temporalWordSpacedIn(sql string) bool {
	sql = strings.ToLower(sql)
	for _, word := range []string{"now", "curdate", "curtime", "sysdate", "date_add", "date_sub", "adddate", "subdate"} {
		if strings.Contains(sql, word+" (") || strings.Contains(sql, word+"/*") || strings.Contains(sql, word+" /*") {
			return true
		}
	}
	return false
}
