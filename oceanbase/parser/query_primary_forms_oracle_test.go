package parser

import (
	"context"
	"testing"
)

// TestQueryPrimaryFormsMatchEngine holds query_primary_forms_test.go's verdicts to the
// server, so an image upgrade that changes an answer fails here rather than leaving a
// stale expectation behind.
//
// It runs on the package's own OceanBase (see sharedserver.PrivateOceanBase):
//
//	go test ./oceanbase/parser/ -run QueryPrimaryForms -count=1
func TestQueryPrimaryFormsMatchEngine(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS omni_sweep"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "USE omni_sweep"); err != nil {
		t.Fatal(err)
	}
	for _, c := range queryPrimaryCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := operatorVerdict(ctx, db, c.sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", c.sql, err)
			}
			if got != c.accept {
				t.Errorf("OceanBase 4.3.5 accepts %q = %v, this file records %v", c.sql, got, c.accept)
			}
		})
	}
}
