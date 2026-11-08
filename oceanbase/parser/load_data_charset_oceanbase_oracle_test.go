package parser

import (
	"context"
	"testing"
)

// LOAD DATA's charset clause as OceanBase reads it. OceanBase CE 4.3.5.6
// takes the two-word CHARACTER SET followed by a charset_name -- an
// identifier, a quoted name or the reserved word BINARY -- and answers 1064
// for the short CHARSET spelling, where MySQL 8.0 and MariaDB take both. An
// integer is no charset name in either.
//
// Every verdict was measured against oceanbase/oceanbase-ce:4.3.5-lts through
// PREPARE on 2026-11-07; this test holds the parser to the live server, which
// the package's usual MySQL-proxy oracle cannot do for a form the two
// disagree on.
func TestLoadDataCharsetFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	cases := []struct {
		sql    string
		accept bool
	}{
		{"load data infile '/x.csv' into table t character set binary", true},
		{"load data infile '/x.csv' into table t character set 'utf8mb4'", true},
		{"load data infile '/x.csv' into table t character set utf8mb4", true},
		{"load data infile '/x.csv' into table t character set 1", false},
		// The short CHARSET spelling is MySQL's and MariaDB's, not OceanBase's.
		{"load data infile '/x.csv' into table t charset binary", false},
		{"load data infile '/x.csv' into table t charset 'utf8mb4'", false},
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
