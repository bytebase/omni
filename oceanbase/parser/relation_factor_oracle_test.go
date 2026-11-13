package parser

import (
	"context"
	"testing"
)

// relationFactorCases are table names written as relation_factor
// (sql_parser_mysql_mode.y, oceanbase v4.3.5_CE_BP6, line 14870): with a
// leading dot, and with a reserved word after a dot. Only the words of
// mysql_reserved_keyword are names there, and only after a dot that stands
// apart or leads; an index is named by normal_relation_factor, which has no
// leading dot.
var relationFactorCases = []asOfCase{
	{"select * from .t", true},
	{"select * from . t", true},
	{"select * from . table", true},
	{"select * from .table", true},
	{"select * from d . table", true},
	{"select * from d . key as k", true},
	{"select * from .t join . u on true", true},
	{"explain . linear", true},
	{"desc . linear a", true},
	{"create table .t (a int)", true},
	{"create table d . key (a int)", true},
	{"create table u like .t", true},
	{"create table c (a int, foreign key (a) references .t (a))", true},
	{"drop table .t, d . index", true},
	{"truncate table .t", true},
	{"alter table .t add column b int", true},
	{"alter table t rename to .u", true},
	{"rename table .t to .u", true},
	{"insert into .t values (1)", true},
	{"update .t set a = 1", true},
	{"delete from .t", true},
	{"create index i on .t (a)", true},
	{"drop index i on .t", true},
	{"create view .v as select 1", true},
	{"show create table .t", true},
	{"show columns from .t", true},
	{"show index from .t", true},
	{"optimize table .t", true},
	{"alter table t drop tablegroup", true},

	// Not every reserved word is in mysql_reserved_keyword.
	{"explain . select", false},
	{"explain . all", false},
	{"select * from . from", false},
	{"select * from d . distinct", false},
	// A reserved word is a name only after the dot.
	{"select * from table . t", false},
	{"explain linear", false},
	// One dot, and something after it.
	{"select * from ..t", false},
	{"select * from .", false},
	{"select * from . 1", false},
	{"create index .i on t (a)", false},
	// OceanBase's own reserved words.
	{"select * from tablegroup", false},
	{"select * from . tablegroup", false},
	{"select read_consistency from t", false},
	{"select no_parallel from t", false},
	{"select no_rewrite from t", false},
}

// TestRelationFactorFormsOceanBaseOracle holds relationFactorCases to
// OceanBase and omni to the same verdict.
func TestRelationFactorFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range relationFactorCases {
		t.Run(c.sql, func(t *testing.T) {
			if got := prepareAcceptsOceanBase(t, ctx, db, c.sql); got != c.accept {
				t.Fatalf("OceanBase accepts=%v, want %v for %q", got, c.accept, c.sql)
			}
			omniAccepts, err := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, err)
			}
		})
	}
}
