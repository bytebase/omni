package parser

import (
	"context"
	"testing"
)

// tableNameClauseCases are the clauses tbl_name takes after a table's name
// (sql_parser_mysql_mode.y, oceanbase v4.3.5_CE_BP6, line 14080), in the one
// order the grammar has: partitions, sample and seed, flashback, alias, index
// hints. DELETE's single table is a tbl_name too.
var tableNameClauseCases = []asOfCase{
	{"select * from t sample (1)", true},
	{"select * from t sample block (1)", true},
	{"select * from t sample hybrid all (1.5) seed (1)", true},
	{"select * from t sample base (1)", true},
	{"select * from t sample incr (1) a", true},
	{"select * from t sample (1) seed (1) as a use key (i)", true},
	{"select * from t partition (p) sample (1) seed (1) as of snapshot 1 as a use key (i)", true},
	{"select * from t as of snapshot 1", true},
	{"select * from t as of snapshot 1 + 1 a", true},
	{"select * from t as of proctime()", true},
	{"select * from t partition (p) as of snapshot 1 a", true},
	{"select * from t partition (a = 's')", true},
	{"select * from t partition (a = 's', b = 1)", true},
	{"select * from t use key ()", true},
	{"select * from t use key (primary)", true},
	{"select * from t use key (`primary_zone`)", true},
	{"select * from t force key (i) ignore index for order by (i, j)", true},
	{"update t sample (1) set a = 1", true},
	{"delete from t partition (p) a where a.id = 1", true},
	{"delete from t sample (1)", true},
	{"delete from t as of snapshot 1", true},
	{"delete from t use key (i)", true},
	{"delete from t a use key (i) where a.id = 1", true},

	// SAMPLE, SEED and OF are names where their clause does not follow.
	{"select * from t sample", true},
	{"select * from t sample use key (i)", true},
	{"select * from t seed", true},
	{"select * from t as of", true},
	{"select * from t seed (1)", false},

	// The order is fixed.
	{"select * from t a sample (1)", false},
	{"select * from t a partition (p)", false},
	{"select * from t sample (1) partition (p)", false},
	{"select * from t as of snapshot 1 sample (1)", false},
	{"select * from t use key (i) a", false},
	{"delete from t a partition (p)", false},
	{"delete from t as a partition (p)", false},
	// Index hints follow a flashback only when a sample stands before it.
	{"select * from t as of snapshot 1 use key (i)", false},
	{"select * from t as of snapshot 1 a use key (i)", false},
	{"select * from t sample (1) as of snapshot 1 use key (i)", true},

	// A percent is a number, a seed an integer, a snapshot a bit_expr.
	{"select * from t sample (1e1)", true},
	{"select * from t sample ('1')", false},
	{"select * from t sample (1) seed (1.5)", false},
	{"select * from t as of snapshot 1 = 1", false},
	{"select * from t partition (a = 1 + 1)", false},
	// An index is named by an identifier that is no keyword, or by PRIMARY.
	{"select * from t use key (primary_zone)", false},
	{"select * from t use key (status)", false},
	{"select * from t force key ()", false},
	// A table's alias is a name, not a string.
	{"select * from t 'a'", false},
	{"select * from t as 'a'", false},
	// The multi-table DELETE lists bare names.
	{"delete from t a using t as a", false},
	{"delete from t partition (p) using t", false},
	{"delete from t use key (i) using t", false},
	{"delete from t using t", true},
}

// TestTableNameClausesOceanBaseOracle holds tableNameClauseCases to OceanBase
// and omni to the same verdict.
func TestTableNameClausesOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range tableNameClauseCases {
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
