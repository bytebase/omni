package parser

import (
	"context"
	"testing"
)

// columnRefCases are column_ref's forms with a reserved word for a name or
// with a leading dot (sql_parser_mysql_mode.y, oceanbase v4.3.5_CE_BP6, line
// 898). A word of mysql_reserved_keyword names a column after a dot, and a
// table only before a dot and another such word; a leading dot takes a table
// and a column after it.
var columnRefCases = []asOfCase{
	{"select t . key from t", true},
	{"select key . key from t", true},
	{"select d . t . key from t", true},
	{"select d . key . key from t", true},
	{"select . t . a from t", true},
	{"select . t . key from t", true},
	{"select . key . key from t", true},
	{"select . t.a from t", true},
	{"select t . key + 1 from t where t . key > 0 order by t . key", true},
	{"update t set t . key = 1", true},
	{"update t set . t . a = 1", true},
	{"insert into t set key . key = 1", true},
	// NOT is in the list, and before a dot it is a name and not the operator.
	{"select not . not from t", true},

	// A reserved table name takes a reserved column, and there is no K.K.K.
	{"select key . a from t", false},
	{"select key . key . key from t", false},
	{"select t . key . a from t", false},
	{"select . key . a from t", false},
	{"select d . t . key . a from t", false},
	// The leading dot takes two names and no star.
	{"select . a from t", false},
	{"select . t . * from t", false},
	// SELECT is reserved and not in mysql_reserved_keyword.
	{"select t . select from t", false},
	// A reserved word after a dot is a column and never a function.
	{"select t . key (1) from t", false},
	// The JSON arrows take column_definition_ref, which has none of these forms.
	{"select . t . a ->> '$' from t", false},
	{"select t . key -> '$' from t", false},
	{"select t . a -> '$' from t", true},
}

// TestColumnRefFormsOceanBaseOracle holds columnRefCases to OceanBase and omni
// to the same verdict.
func TestColumnRefFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range columnRefCases {
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
