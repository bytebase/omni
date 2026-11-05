package parser

import (
	"context"
	"testing"
)

// TestInsertFormsMatchEngine holds insert_forms_test.go's verdicts to the server, so
// an image upgrade that changes an answer fails here rather than leaving a stale
// expectation behind.
//
// OceanBase needs about 8 GB, which CI does not have, so this reads OCEANBASE_WALK_DSN
// and skips without it:
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/parser/ -run InsertForms -count=1
func TestInsertFormsMatchEngine(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS omni_sweep"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "USE omni_sweep"); err != nil {
		t.Fatal(err)
	}
	accepts := func(sql string) (bool, error) { return operatorVerdict(ctx, db, sql) }
	for _, c := range insertFormCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := accepts(c.sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", c.sql, err)
			}
			if got != c.accept {
				t.Errorf("OceanBase 4.3.5 accepts %q = %v, this file records %v", c.sql, got, c.accept)
			}
		})
	}
}
