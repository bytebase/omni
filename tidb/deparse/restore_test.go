package deparse_test

import (
	"testing"

	ast "github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/deparse"
	"github.com/bytebase/omni/tidb/parser"
)

// TestRestoreExpr pins how Restore writes expressions. Every want is the
// expression as pingcap/tidb:v8.5.5 stores it in a view body (SHOW CREATE
// VIEW over SELECT <expr> FROM t); tidb/catalog's container tests compare the
// same forms against the server.
func TestRestoreExpr(t *testing.T) {
	cases := []struct{ in, want string }{
		{"a+b", "`a`+`b`"},
		{"a + b * 2", "`a`+`b`*2"},
		{"(a + 1) * 3", "(`a`+1)*3"},
		{"a div b", "`a` DIV `b`"},
		{"a mod b", "`a`%`b`"},
		{"mod(a, 2)", "`a`%2"},
		// TiDB regroups MOD as written into the body, and the view runs
		// the regrouped text: SELECT * over a view of 2*MOD(5,3) answers 1.
		{"2*mod(5,3)", "2*5%3"},
		{"mod(a+1, 2)", "`a`+1%2"},
		{"a <> b", "`a`!=`b`"},
		{"a && b", "`a` AND `b`"},
		{"a || b", "`a` OR `b`"},
		{"a xor b", "`a` XOR `b`"},
		{"a <=> b", "`a`<=>`b`"},
		{"not (a > 1)", "NOT (`a`>1)"},
		{"!(a < 2)", "!(`a`<2)"},
		{"- a", "-`a`"},
		{"+a", "+`a`"},
		{"~a", "~`a`"},
		{"1e3", "1e+03"},
		{".5", "0.5"},
		{"1.50", "1.50"},
		{"0x1F", "x'1f'"},
		{"x'1F'", "x'1f'"},
		{"b'101'", "b'101'"},
		{"'it''s'", "_UTF8MB4'it''s'"},
		{`'back\\slash'`, `_UTF8MB4'back\\slash'`},
		{"_latin1'x'", "_LATIN1'x'"},
		{"true", "TRUE"},
		{"null", "NULL"},
		{"date '2020-01-01'", "DATE '2020-01-01'"},
		{"a in (1,2,3)", "`a` IN (1,2,3)"},
		{"a not between 1 and 2", "`a` NOT BETWEEN 1 AND 2"},
		{"c not like 'x' escape '!'", "`c` NOT LIKE _UTF8MB4'x' ESCAPE '!'"},
		{"c not regexp 'a'", "`c` NOT REGEXP _UTF8MB4'a'"},
		{"c rlike 'b'", "`c` REGEXP _UTF8MB4'b'"},
		{"a is unknown", "`a` IS NULL"},
		{"a is not false", "`a` IS NOT FALSE"},
		{"row(1,2) = row(a,b)", "ROW(1,2)=ROW(`a`,`b`)"},
		{"case when a>0 then 'p' else 'n' end", "CASE WHEN `a`>0 THEN _UTF8MB4'p' ELSE _UTF8MB4'n' END"},
		{"case a when 1 then 2 end", "CASE `a` WHEN 1 THEN 2 END"},
		{"cast(a as char(3))", "CAST(`a` AS CHAR(3))"},
		{"cast(a as decimal(10,2))", "CAST(`a` AS DECIMAL(10, 2))"},
		{"cast(a as decimal)", "CAST(`a` AS DECIMAL(10))"},
		{"cast(a as signed integer)", "CAST(`a` AS SIGNED)"},
		{"cast(a as datetime(3))", "CAST(`a` AS DATETIME(3))"},
		{"convert(c using latin1)", "CONVERT(`c` USING 'latin1')"},
		{"c collate utf8mb4_bin", "`c` COLLATE utf8mb4_bin"},
		{"binary c", "BINARY `c`"},
		{"count(*)", "COUNT(1)"},
		{"count(distinct a, b)", "COUNT(DISTINCT `a`, `b`)"},
		{"sum(all a)", "SUM(`a`)"},
		{"std(a)", "STDDEV_POP(`a`)"},
		{"group_concat(c)", "GROUP_CONCAT(`c` SEPARATOR ',')"},
		{`group_concat(c separator '\\')`, `GROUP_CONCAT(` + "`c`" + ` SEPARATOR '\\')`},
		{`group_concat(c separator 'a\\nb')`, `GROUP_CONCAT(` + "`c`" + ` SEPARATOR 'a\\nb')`},
		{"group_concat(c separator '''')", "GROUP_CONCAT(`c` SEPARATOR '''')"},
		{"group_concat(distinct c order by c desc separator ';')", "GROUP_CONCAT(DISTINCT `c` ORDER BY `c` DESC SEPARATOR ';')"},
		{"concat(c, 'x')", "CONCAT(`c`, _UTF8MB4'x')"},
		{"current_timestamp", "CURRENT_TIMESTAMP()"},
		{"now(3)", "NOW(3)"},
		{"char(65)", "CHAR_FUNC(65, NULL)"},
		{"substring(c from 1 for 2)", "SUBSTRING(`c`, 1, 2)"},
		{"position('a' in c)", "POSITION(_UTF8MB4'a' IN `c`)"},
		{"extract(year from d)", "EXTRACT(YEAR FROM `d`)"},
		{"timestampdiff(second, d, d)", "TIMESTAMPDIFF(SECOND, `d`, `d`)"},
		{"get_format(date, 'usa')", "GET_FORMAT(DATE, _UTF8MB4'usa')"},
		{"date_add(d, interval 1 day)", "DATE_ADD(`d`, INTERVAL 1 DAY)"},
		{"adddate(d, 1)", "ADDDATE(`d`, INTERVAL 1 DAY)"},
		{"d + interval 1 day", "DATE_ADD(`d`, INTERVAL 1 DAY)"},
		{"interval 1 day + d", "DATE_ADD(`d`, INTERVAL 1 DAY)"},
		{"now() - interval 1 day", "DATE_SUB(NOW(), INTERVAL 1 DAY)"},
		{"trim(c)", "TRIM(`c`)"},
		{"trim(leading 'x' from c)", "TRIM(LEADING _UTF8MB4'x' FROM `c`)"},
		{"trim('x' from c)", "TRIM(_UTF8MB4'x' FROM `c`)"},
		{"j->'$.a'", "JSON_EXTRACT(`j`, _UTF8MB4'$.a')"},
		{"j->>'$.a'", "JSON_UNQUOTE(JSON_EXTRACT(`j`, _UTF8MB4'$.a'))"},
		{"1 member of (j)", "1 MEMBER OF (`j`)"},
		{"exists (select 1)", "EXISTS (SELECT 1)"},
		{"not exists (select 1)", "NOT EXISTS (SELECT 1)"},
		{"a = any (select 1)", "`a`=ANY (SELECT 1)"},
		{"a > all (select 1)", "`a`>ALL (SELECT 1)"},
		{"sum(a) over ()", "SUM(`a`) OVER ()"},
		{"avg(a) over (partition by b rows between unbounded preceding and current row)", "AVG(`a`) OVER (PARTITION BY `b` ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)"},
		{"first_value(a) over (order by b range between 1 preceding and 1 following)", "FIRST_VALUE(`a`) OVER (ORDER BY `b` RANGE BETWEEN 1 PRECEDING AND 1 FOLLOWING)"},
		{"ntile(2) over w", "NTILE(2) OVER `w`"},
		{"count(*) over ()", "COUNT(1) OVER ()"},
		{"@a > 1", "@`a`>1"},
		{"@aB + 1", "@`aB`+1"},
		{"@@var", "@@`var`"},
		{"@@global.b = 'foo'", "@@GLOBAL.`b`=_UTF8MB4'foo'"},
		{"@@session.C", "@@SESSION.`c`"},
		{"@@local.autocommit", "@@SESSION.`autocommit`"},
	}
	for _, tc := range cases {
		stmts, err := parser.Parse("SELECT " + tc.in + " FROM t")
		if err != nil {
			t.Errorf("%s: parse: %v", tc.in, err)
			continue
		}
		sel := stmts.Items[0].(*ast.SelectStmt)
		target := sel.TargetList[0]
		if rt, ok := target.(*ast.ResTarget); ok {
			target = rt.Val
		}
		if got := deparse.Restore(target); got != tc.want {
			t.Errorf("%s:\n got: %s\nwant: %s", tc.in, got, tc.want)
		}
	}
}

// TestRestoreSelect pins whole view bodies: the query prepared by a Resolver
// in Restore mode over tables t (a, b) and t2 (a, x) in schema test, then
// restored. Every want is what pingcap/tidb:v8.5.5 stores for the query.
func TestRestoreSelect(t *testing.T) {
	tables := map[string]*deparse.ResolverTable{
		"t":  {Name: "t", Columns: []deparse.ResolverColumn{{Name: "a", Position: 1, CharLength: 11, Integer: true}, {Name: "b", Position: 2, CharLength: 11, Integer: true}}},
		"t2": {Name: "t2", Columns: []deparse.ResolverColumn{{Name: "a", Position: 1, CharLength: 11, Integer: true}, {Name: "x", Position: 2, CharLength: 11, Integer: true}}},
	}
	lookup := func(name string) *deparse.ResolverTable { return tables[name] }
	cases := []struct{ in, want string }{
		{"SELECT a FROM t WHERE a > 0", "SELECT `a` AS `a` FROM `test`.`t` WHERE `a`>0"},
		{"select a+b, a + b, a*2 from t", "SELECT `a`+`b` AS `a+b`,`a`+`b` AS `a + b`,`a`*2 AS `a*2` FROM `test`.`t`"},
		{"select a+b, a+b, a+b from t", "SELECT `a`+`b` AS `a+b`,`a`+`b` AS `Name_exp_a+b`,`a`+`b` AS `Name_exp_1_a+b` FROM `test`.`t`"},
		{"select 1, 'x', 1.5, NULL, true from t", "SELECT 1 AS `1`,_UTF8MB4'x' AS `x`,1.5 AS `1.5`,NULL AS `NULL`,TRUE AS `TRUE` FROM `test`.`t`"},
		{"select '', '  x' from t", "SELECT _UTF8MB4'',_UTF8MB4'  x' AS `x` FROM `test`.`t`"},
		{"select * from t", "SELECT `test`.`t`.`a` AS `a`,`test`.`t`.`b` AS `b` FROM `test`.`t`"},
		{"select * from t as x", "SELECT `test`.`x`.`a` AS `a`,`test`.`x`.`b` AS `b` FROM `test`.`t` AS `x`"},
		{"select * from t join t2 using (a)", "SELECT `test`.`t`.`a` AS `a`,`test`.`t`.`b` AS `b`,`test`.`t2`.`x` AS `x` FROM `test`.`t` JOIN `test`.`t2` USING (`a`)"},
		{"select * from t right join t2 using (a)", "SELECT `test`.`t2`.`a` AS `a`,`test`.`t2`.`x` AS `x`,`test`.`t`.`b` AS `b` FROM `test`.`t` RIGHT JOIN `test`.`t2` USING (`a`)"},
		{"select t.a from t, t2", "SELECT `t`.`a` AS `a` FROM (`test`.`t`) JOIN `test`.`t2`"},
		{"select t.a from t, t2, t2 as t3", "SELECT `t`.`a` AS `a` FROM ((`test`.`t`) JOIN `test`.`t2`) JOIN `test`.`t2` AS `t3`"},
		{"select t.a from (select 1 as q) dt, t", "SELECT `t`.`a` AS `a` FROM (SELECT 1 AS `q`) AS `dt`, `test`.`t`"},
		{"select t.a from t cross join t2", "SELECT `t`.`a` AS `a` FROM `test`.`t` JOIN `test`.`t2`"},
		{"select t.a, t2.x from t inner join t2 on t.a=t2.a left join t2 as t3 on t3.a = t.a", "SELECT `t`.`a` AS `a`,`t2`.`x` AS `x` FROM (`test`.`t` JOIN `test`.`t2` ON `t`.`a`=`t2`.`a`) LEFT JOIN `test`.`t2` AS `t3` ON `t3`.`a`=`t`.`a`"},
		{"select a from t order by 1", "SELECT `a` AS `a` FROM `test`.`t` ORDER BY 1"},
		{"select a, b+1 from t group by 2, 1", "SELECT `a` AS `a`,`b`+1 AS `b+1` FROM `test`.`t` GROUP BY `b`+1,`a`"},
		{"select a from t limit 2 offset 1", "SELECT `a` AS `a` FROM `test`.`t` LIMIT 1,2"},
		{"select a from t for update", "SELECT `a` AS `a` FROM `test`.`t` "},
		{"select 1 from dual where 1", "SELECT 1 AS `1` FROM DUAL WHERE 1"},
		{"select distinctrow a from t", "SELECT DISTINCT `a` AS `a` FROM `test`.`t`"},
		{"select (a), (1), (a+1) from t", "SELECT (`a`) AS `(a)`,(1) AS `1`,(`a`+1) AS `(a+1)` FROM `test`.`t`"},
		{"select cast(a as char), cast(a as binary), cast(a+b as char) from t", "SELECT CAST(`a` AS CHAR(11)) AS `cast(a as char)`,CAST(`a` AS BINARY(11)) AS `cast(a as binary)`,CAST(`a`+`b` AS CHAR(20)) AS `cast(a+b as char)` FROM `test`.`t`"},
		{"select count(*) from t group by a having sum(a) > 10", "SELECT COUNT(1) AS `count(*)` FROM `test`.`t` GROUP BY `a` HAVING SUM(`test`.`t`.`a`)>10"},
		{"select a from t group by a having sum(a) > 10", "SELECT `a` AS `a` FROM `test`.`t` GROUP BY `a` HAVING SUM(`a`)>10"},
		{"select b from t group by b order by sum(b), max(a)", "SELECT `b` AS `b` FROM `test`.`t` GROUP BY `b` ORDER BY SUM(`b`),MAX(`test`.`t`.`a`)"},
		{"select count(*) from (select a from t) d group by a having sum(a) > 1", "SELECT COUNT(1) AS `count(*)` FROM (SELECT `a` AS `a` FROM `test`.`t`) AS `d` GROUP BY `a` HAVING SUM(`test`.`d`.`a`)>1"},
		{"with x as (select a from t) select * from x", "WITH `x` AS (SELECT `a` AS `a` FROM `test`.`t`) SELECT `x`.`a` AS `a` FROM `x`"},
		{"with recursive r (n) as (select 1 union all select n+1 from r where n<3) select n from r", "WITH RECURSIVE `r` (`n`) AS (SELECT 1 AS `1` UNION ALL SELECT `n`+1 AS `n+1` FROM `r` WHERE `n`<3) SELECT `n` AS `n` FROM `r`"},
		{"(select a from t) union (select x from t2) order by a limit 2", "(SELECT `a` AS `a` FROM `test`.`t`) UNION (SELECT `x` AS `x` FROM `test`.`t2`) ORDER BY `a` LIMIT 2"},
		{"select a from t except select a from t2", "SELECT `a` AS `a` FROM `test`.`t` EXCEPT SELECT `a` AS `a` FROM `test`.`t2`"},
		{"select a from t where exists (select * from t2 where t2.a = t.a)", "SELECT `a` AS `a` FROM `test`.`t` WHERE EXISTS (SELECT `test`.`t2`.`a` AS `a`,`test`.`t2`.`x` AS `x` FROM `test`.`t2` WHERE `t2`.`a`=`t`.`a`)"},
		{"select * from (select a+1, b from t) d", "SELECT `d`.`a+1` AS `a+1`,`test`.`d`.`b` AS `b` FROM (SELECT `a`+1 AS `a+1`,`b` AS `b` FROM `test`.`t`) AS `d`"},
		{"select 1, 1 union select 2, 3", "SELECT 1 AS `1`,1 AS `Name_exp_1` UNION SELECT 2 AS `2`,3 AS `3`"},
	}
	for _, tc := range cases {
		stmts, err := parser.Parse(tc.in)
		if err != nil {
			t.Errorf("%s: parse: %v", tc.in, err)
			continue
		}
		sel := stmts.Items[0].(*ast.SelectStmt)
		r := &deparse.Resolver{Lookup: lookup, Source: tc.in, Restore: true, Database: "test"}
		r.Resolve(sel)
		if got := deparse.RestoreSelect(sel); got != tc.want {
			t.Errorf("%s:\n got: %s\nwant: %s", tc.in, got, tc.want)
		}
	}
}

// TestRestoreWithDDLFlags pins the flag sets TiDB's DDL restores table
// expressions with. Every want is what pingcap/tidb:v8.5.5 prints in SHOW
// CREATE TABLE for the expression: inside GENERATED ALWAYS AS (...), an
// expression index part, DEFAULT (...), or PARTITION BY ... (...).
func TestRestoreWithDDLFlags(t *testing.T) {
	cases := []struct {
		flags    deparse.RestoreFlags
		in, want string
	}{
		{deparse.RestoreForGeneratedColumn, "a + b * 2", "`a` + `b` * 2"},
		{deparse.RestoreForGeneratedColumn, "CONCAT(c, 'x')", "concat(`c`, _utf8mb4'x')"},
		{deparse.RestoreForGeneratedColumn, "CASE WHEN a > 0 THEN 1 ELSE NOT b END", "case when `a` > 0 then 1 else not `b` end"},
		{deparse.RestoreForGeneratedColumn, "CAST(c AS SIGNED)", "cast(`c` as signed)"},
		{deparse.RestoreForGeneratedColumn, "d + INTERVAL 1 DAY", "date_add(`d`, interval 1 day)"},
		{deparse.RestoreForGeneratedColumn, "a DIV b", "`a` div `b`"},
		{deparse.RestoreForGeneratedColumn, "IFNULL(a, 0) AND b IN (1,2)", "ifnull(`a`, 0) and `b` in (1,2)"},
		{deparse.RestoreForGeneratedColumn, "j->>'$.x'", "json_unquote(json_extract(`j`, _utf8mb4'$.x'))"},
		{deparse.RestoreForGeneratedColumn, "-a", "-`a`"},
		{deparse.RestoreForGeneratedColumn, "DATE '2020-01-01'", "DATE '2020-01-01'"},
		{deparse.RestoreForGeneratedColumn, "c COLLATE utf8mb4_bin", "`c` collate utf8mb4_bin"},
		{deparse.RestoreForGeneratedColumn, "(a + 1) * 2", "(`a` + 1) * 2"},
		{deparse.RestoreForGeneratedColumn, "a * b", "`a` * `b`"},
		{deparse.RestoreForGeneratedColumn, "LOWER(c)", "lower(`c`)"},
		{deparse.RestoreForGeneratedColumn, "a + b > 1 AND c LIKE 'x%'", "`a` + `b` > 1 and `c` like _utf8mb4'x%'"},
		{deparse.RestoreForColumnDefault, "UUID()", "uuid()"},
		{deparse.RestoreForColumnDefault, "RAND()", "rand()"},
		{deparse.RestoreForColumnDefault, "REPLACE(UPPER(UUID()), '-', '')", "replace(upper(uuid()), _utf8mb4'-', _utf8mb4'')"},
		{deparse.RestoreForColumnDefault, "STR_TO_DATE('2020', '%Y')", "str_to_date(_utf8mb4'2020', _utf8mb4'%Y')"},
		{deparse.RestoreForPartition, "YEAR(created)", "YEAR(`created`)"},
		{deparse.RestoreForPartition, "a + b * 2", "(`a`+(`b`*2))"},
		{deparse.RestoreForPartition, "a DIV 2 + MOD(b, 3)", "((`a` DIV 2)+(`b`%3))"},
	}
	for _, tc := range cases {
		expr, err := parser.ParseExpr(tc.in)
		if err != nil {
			t.Errorf("%s: parse: %v", tc.in, err)
			continue
		}
		if got := deparse.RestoreWith(expr, tc.flags); got != tc.want {
			t.Errorf("%s (flags %b):\n got: %s\nwant: %s", tc.in, tc.flags, got, tc.want)
		}
	}
}
