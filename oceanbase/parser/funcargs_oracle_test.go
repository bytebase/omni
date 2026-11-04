package parser

import (
	"context"
	"testing"
)

// TestOceanBaseFunctionArgumentOracle reads the verdicts funcArgCases records
// back from a live OceanBase, so they are this engine's and not MySQL's
// inherited. They differ inside these parentheses: OceanBase takes no
// DISTINCT ALL anywhere, takes DISTINCT on JSON_ARRAYAGG and neither prefix on
// JSON_OBJECTAGG, and has no SQL_TSI_ unit synonyms.
//
// The server needs about 8 GB, which CI does not have, so this reads
// OCEANBASE_WALK_DSN and skips without it:
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/parser/ -run OceanBaseFunctionArgument -count=1
func TestOceanBaseFunctionArgumentOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()

	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS omni_funcargs",
		"USE omni_funcargs",
		"DROP TABLE IF EXISTS t",
		"DROP TABLE IF EXISTS sales",
		"CREATE TABLE t (c1 INT, c2 INT, b JSON, a DATETIME, profit INT)",
		"CREATE TABLE sales (profit INT)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	defer db.ExecContext(ctx, "DROP DATABASE IF EXISTS omni_funcargs")

	for _, c := range funcArgCases {
		accepts, err := engineAccepts(ctx, db, c.sql)
		if err != nil {
			t.Fatalf("no verdict for %q: %v", c.sql, err)
		}
		if accepts != c.accept {
			t.Errorf("OceanBase accepts %q = %v, this file records %v", c.sql, accepts, c.accept)
		}
	}
}
