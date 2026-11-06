package parser

import (
	"context"
	"testing"
)

// TestOceanBaseTemporalFormsMatchEngine reads the verdicts temporalFormCases records
// back from a live OceanBase, so they are this engine's and not MySQL's inherited.
// They differ here: these words are keywords wherever they stand, so now (-1) is a
// 1064 where MySQL reads an ordinary call, and a literal up to 18446744073709551615
// parses as a precision.
//
// It runs on the package's own OceanBase (see testcontainer.OceanBase).
func TestOceanBaseTemporalFormsMatchEngine(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()

	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS omni_temporal",
		"USE omni_temporal",
		"DROP TABLE IF EXISTS t",
		"CREATE TABLE t (a DATETIME, b INT, c INT, `now` INT, `curdate` INT, `curtime` INT, " +
			"`date_add` INT, `date_sub` INT, `sysdate` INT, `adddate` INT, `subdate` INT)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	defer db.ExecContext(ctx, "DROP DATABASE IF EXISTS omni_temporal")

	for _, c := range temporalFormCases {
		if accepts := prepareAcceptsOceanBase(t, ctx, db, c.sql); accepts != c.accept {
			t.Errorf("OceanBase accepts %q = %v, this file records %v", c.sql, accepts, c.accept)
		}
	}
}
