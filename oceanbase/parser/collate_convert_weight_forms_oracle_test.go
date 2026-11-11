package parser

import (
	"context"
	"testing"
)

// COLLATE, CONVERT(... USING ...) and WEIGHT_STRING(... AS ...) as OceanBase
// reads them. OceanBase answers 1064 for the reserved word BINARY as a
// collation name in either a column option or an expression postfix (its
// collation name has no BINARY arm, unlike its charset name), takes BINARY as
// a charset name in CONVERT(... USING ...), and carries WEIGHT_STRING's
// AS CHAR(n) / AS CHARACTER(n) / AS BINARY(n) clause whose length must be at
// least 1. It also keeps the LEVEL clause of MySQL 5.7, which MySQL 8.0
// removed.
//
// Every expectation was measured against oceanbase/oceanbase-ce:4.3.5-lts
// (5.7.25-OceanBase_CE-v4.3.5.6) through PREPARE; this test holds the file to
// the server.
func TestCollateConvertWeightFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	cases := []struct {
		sql    string
		accept bool
	}{
		{"ALTER TABLE t ADD c JSON COLLATE BINARY", false},
		{"ALTER TABLE t ADD c JSON COLLATE utf8mb4_bin", true},
		{"SELECT x COLLATE BINARY", false},
		{"SELECT x COLLATE utf8mb4_bin", true},

		{"SELECT CONVERT(x USING BINARY)", true},
		{"SELECT CONVERT(x USING utf8mb4)", true},

		{"SELECT WEIGHT_STRING(x)", true},
		{"SELECT WEIGHT_STRING(x AS CHAR(1))", true},
		{"SELECT WEIGHT_STRING(x AS CHARACTER(1))", true},
		{"SELECT WEIGHT_STRING(x AS BINARY(1))", true},
		{"SELECT WEIGHT_STRING(x AS CHAR(0))", false},
		{"SELECT WEIGHT_STRING(x AS CHAR)", false},
		{"SELECT WEIGHT_STRING(x AS DECIMAL(1))", false},

		// LEVEL, the MySQL 5.7 clause that MySQL 8.0 removed and OceanBase
		// keeps: a list of numbers, each with [ASC|DESC] [REVERSE] in that
		// order, or a range of two bare numbers. It follows the plain form
		// and AS CHAR(n), not AS BINARY(n). A number out of range is clamped.
		{"SELECT WEIGHT_STRING(x LEVEL 1)", true},
		{"SELECT WEIGHT_STRING(x LEVEL 0)", true},
		{"SELECT WEIGHT_STRING(x LEVEL 7)", true},
		{"SELECT WEIGHT_STRING(x LEVEL 1 ASC)", true},
		{"SELECT WEIGHT_STRING(x LEVEL 1 DESC)", true},
		{"SELECT WEIGHT_STRING(x LEVEL 1 REVERSE)", true},
		{"SELECT WEIGHT_STRING(x LEVEL 1 ASC REVERSE)", true},
		{"SELECT WEIGHT_STRING(x LEVEL 1 DESC REVERSE)", true},
		{"SELECT WEIGHT_STRING(x LEVEL 1 REVERSE DESC)", false},
		{"SELECT WEIGHT_STRING(x LEVEL 1, 2)", true},
		{"SELECT WEIGHT_STRING(x LEVEL 1 DESC, 2 REVERSE)", true},
		{"SELECT WEIGHT_STRING(x LEVEL 1 - 2)", true},
		{"SELECT WEIGHT_STRING(x LEVEL 2 - 1)", true},
		{"SELECT WEIGHT_STRING(x LEVEL 1 ASC - 2)", false},
		{"SELECT WEIGHT_STRING(x LEVEL 1 DESC - 2)", false},
		{"SELECT WEIGHT_STRING(x LEVEL 1 - 2 DESC)", false},
		{"SELECT WEIGHT_STRING(x LEVEL 1, 2 - 3)", false},
		{"SELECT WEIGHT_STRING(x LEVEL 1 - 2, 3)", false},
		{"SELECT WEIGHT_STRING(x LEVEL)", false},
		{"SELECT WEIGHT_STRING(x AS CHAR(4) LEVEL 1)", true},
		{"SELECT WEIGHT_STRING(x AS CHARACTER(4) LEVEL 1 - 2)", true},
		{"SELECT WEIGHT_STRING(x AS BINARY(4) LEVEL 1)", false},
		{"SELECT WEIGHT_STRING(x AS BINARY(4) LEVEL 1 DESC)", false},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			accepts := prepareAcceptsOceanBase(t, ctx, db, c.sql)
			if accepts != c.accept {
				t.Fatalf("OceanBase accepts=%v, want %v for %q", accepts, c.accept, c.sql)
			}
			_, perr := Parse(c.sql)
			if omniAccepts := perr == nil; omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, perr)
			}
		})
	}
}
