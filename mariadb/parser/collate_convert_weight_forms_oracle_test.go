package parser

import "testing"

// COLLATE, CONVERT(... USING ...) and WEIGHT_STRING(... AS ...) as MariaDB
// 11.8 reads them. MariaDB answers 1064 for the reserved word BINARY as a
// collation name in either a column option or an expression postfix (unlike
// MySQL 8.0, which takes the column-option form), takes BINARY as a charset
// name in CONVERT(... USING ...), and carries WEIGHT_STRING's AS CHAR(n) /
// AS CHARACTER(n) / AS BINARY(n) clause whose length must be at least 1.
//
// Every expectation was measured against mariadb:11.8.8 through PREPARE;
// TestCollateConvertWeightFormsOracle holds the file to it.
var collateConvertWeightCases = []struct {
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
}

func TestCollateConvertWeightFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, tc := range collateConvertWeightCases {
		t.Run(tc.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, tc.sql); got != tc.accept {
				t.Fatalf("MariaDB 11.8.8 accepts=%v, want %v for %q", got, tc.accept, tc.sql)
			}
			_, err := Parse(tc.sql)
			if omniAccepts := err == nil; omniAccepts != tc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB 11.8.8 accepts=%v for %q: %v",
					omniAccepts, tc.accept, tc.sql, err)
			}
		})
	}
}
