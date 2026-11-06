package parser

import (
	"context"
	"testing"
)

// TestAlterColumnFormsOceanBaseOracle asserts OceanBase still answers
// alterColumnFormCases as recorded, and that omni agrees. The oracle runs each
// statement against a seeded table, so the accepted forms really alter it.
func TestAlterColumnFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"DROP DATABASE IF EXISTS obaltcol",
		"CREATE DATABASE obaltcol",
		"USE obaltcol",
		"CREATE TABLE t (a INT, b INT)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() { _, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obaltcol") })
	for _, c := range alterColumnFormCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := engineAccepts(ctx, db, c.sql)
			if err != nil {
				t.Fatalf("the oracle gave no verdict on %q: %v", c.sql, err)
			}
			if got != c.accept {
				t.Fatalf("OceanBase accepts %q = %v, the record says %v", c.sql, got, c.accept)
			}
			omniAccepts, perr := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, perr)
			}
		})
	}
}
