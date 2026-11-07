package deparse

import (
	"testing"

	ast "github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
)

// What the function-argument forms deparse to, in both forms this package
// writes at once: the engine-stored text of the item, and the name the engine
// gives an unaliased one.
//
// The quantifier is the case to watch. The engine keeps only whether the
// aggregate de-duplicates, so it stores sum(distinct `c1`) for DISTINCT,
// DISTINCTROW and either DISTINCT ALL spelling, and drops a bare ALL — but it
// names the column after the text as written, every spelling kept. Both
// readings are measured with SHOW CREATE VIEW and the column name of the view
// on mariadb:11.8.8 and mysql:8.0.
//
// A quantifier dropped on the way out would leave a statement that still
// parses and no longer means the same, which the accept/reject gate cannot
// see: these assertions are what catch it.
func TestFunctionArgumentDeparse(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{"SELECT SUM(c1) FROM t", "select sum(`c1`) AS `SUM(c1)` from `t`"},
		{"SELECT SUM(DISTINCT c1) FROM t", "select sum(distinct `c1`) AS `SUM(DISTINCT c1)` from `t`"},
		{"SELECT SUM(DISTINCTROW c1) FROM t", "select sum(distinct `c1`) AS `SUM(DISTINCTROW c1)` from `t`"},
		{"SELECT SUM(ALL c1) FROM t", "select sum(`c1`) AS `SUM(ALL c1)` from `t`"},
		{"SELECT SUM(DISTINCT ALL c1) FROM t", "select sum(distinct `c1`) AS `SUM(DISTINCT ALL c1)` from `t`"},
		{"SELECT SUM(DISTINCTROW ALL c1) FROM t", "select sum(distinct `c1`) AS `SUM(DISTINCTROW ALL c1)` from `t`"},
		{"SELECT COUNT(DISTINCTROW c1, c2) FROM t", "select count(distinct `c1`,`c2`) AS `COUNT(DISTINCTROW c1, c2)` from `t`"},
		{"SELECT COUNT(ALL c1) FROM t", "select count(`c1`) AS `COUNT(ALL c1)` from `t`"},
		{"SELECT BIT_AND(ALL c1) FROM t", "select bit_and(`c1`) AS `BIT_AND(ALL c1)` from `t`"},
		{"SELECT STD(ALL c1) FROM t", "select std(`c1`) AS `STD(ALL c1)` from `t`"},
		{"SELECT JSON_ARRAYAGG(DISTINCT c1) FROM t", "select json_arrayagg(distinct `c1`) AS `JSON_ARRAYAGG(DISTINCT c1)` from `t`"},
		{"SELECT JSON_OBJECTAGG(c1, c2) FROM t", "select json_objectagg(`c1`,`c2`) AS `JSON_OBJECTAGG(c1, c2)` from `t`"},
		{"SELECT GROUP_CONCAT(DISTINCTROW c2, c1 ORDER BY c1 SEPARATOR '-') FROM t", "select group_concat(distinct `c2`,`c1` order by `c1` ASC separator '-') AS `GROUP_CONCAT(DISTINCTROW c2, c1 ORDER BY c1 SEPARATOR '-')` from `t`"},
		{"SELECT GROUP_CONCAT(c1) OVER () FROM t", "select group_concat(`c1` separator ',') over ( ) AS `GROUP_CONCAT(c1 SEPARATOR ',') OVER ()` from `t`"},
		{"SELECT EXTRACT(DAY_HOUR FROM a) FROM t", "select extract(day_hour from `a`) AS `EXTRACT(DAY_HOUR FROM a)` from `t`"},
		{"SELECT EXTRACT(SQL_TSI_DAY FROM a) FROM t", "select extract(day from `a`) AS `EXTRACT(DAY FROM a)` from `t`"},
		{"SELECT CAST(b AS SIGNED) FROM t", "select cast(`b` as signed) AS `CAST(b AS SIGNED)` from `t`"},
		{"SELECT POSITION('bar' IN c1) FROM t", "select locate('bar',`c1`) AS `POSITION('bar' IN c1)` from `t`"},
		{"SELECT INTERVAL(c1, 1, 2) FROM t", "select interval(`c1`,1,2) AS `INTERVAL(c1, 1, 2)` from `t`"},
	} {
		list, err := parser.Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.sql, err)
		}
		sel, ok := list.Items[0].(*ast.SelectStmt)
		if !ok {
			t.Fatalf("Parse(%q): %T", tc.sql, list.Items[0])
		}
		if got := DeparseSelect(sel); got != tc.want {
			t.Errorf("DeparseSelect(%q) =\n  %q\nwant:\n  %q", tc.sql, got, tc.want)
		}
	}
}

// POSITION's operands have to be resolved like any other expression's: the
// resolver qualifies a column reference and expands a star, and a node its
// type switch does not know leaves the subtree alone — a view stored from
// that tree would carry `select /* unsupported: *ast.StarExpr */ from ...`.
func TestFunctionArgumentResolve(t *testing.T) {
	lookup := func(name string) *ResolverTable {
		if name != "t" {
			return nil
		}
		return &ResolverTable{Name: "t", Columns: []ResolverColumn{{Name: "s", Position: 1}}}
	}
	for _, tc := range []struct{ sql, want string }{
		{"SELECT POSITION('x' IN s) FROM t", "select locate('x',`t`.`s`) AS `POSITION('x' IN s)` from `t`"},
		{"SELECT POSITION((SELECT s FROM t) IN s) FROM t", "select locate((select `t`.`s` from `t`),`t`.`s`) AS `POSITION((SELECT s FROM t) IN s)` from `t`"},
		{"SELECT POSITION('x' IN (SELECT * FROM t))", "select locate('x',(select `t`.`s` from `t`)) AS `Name_exp_1`"},
		{"SELECT EXTRACT(DAY FROM (SELECT * FROM t))", "select extract(day from (select `t`.`s` from `t`)) AS `Name_exp_1`"},
	} {
		list, err := parser.Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.sql, err)
		}
		sel, ok := list.Items[0].(*ast.SelectStmt)
		if !ok {
			t.Fatalf("Parse(%q): %T", tc.sql, list.Items[0])
		}
		(&Resolver{Lookup: lookup}).Resolve(sel)
		if got := DeparseSelect(sel); got != tc.want {
			t.Errorf("%q resolved to\n  %q\nwant:\n  %q", tc.sql, got, tc.want)
		}
	}
}
