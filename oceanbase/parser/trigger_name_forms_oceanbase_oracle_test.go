package parser

import (
	"context"
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

// TestTriggerNameFormsOceanBaseOracle asserts omni agrees with OceanBase CE 4.3.5.6 on
// triggerNameCases.
func TestTriggerNameFormsOceanBaseOracle(t *testing.T) {
	runFormsOnOceanBase(t, "trigger_name_forms", triggerNameCases)
}

// runFormsOnOceanBase is holdFormsToOceanBase for statements PREPARE cannot hold. A stored
// program is read by another parser, which PREPARE never reaches: it answers 1295 for a
// CREATE PROCEDURE whatever its body. So each case is run, in a database named for label
// that holds the table t, and 1064 is the one answer that says it did not parse.
func runFormsOnOceanBase(t *testing.T, label string, cases []parserFormCase) {
	t.Helper()
	db := testcontainer.OpenDatabase(t, label)
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, "CREATE TABLE t (id INT PRIMARY KEY, a INT)"); err != nil {
		t.Fatalf("create the table the forms name: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			accepts := true
			if _, err := db.ExecContext(ctx, tc.sql); err != nil {
				var myErr *gomysql.MySQLError
				if !errors.As(err, &myErr) {
					t.Fatalf("no verdict for %q: %v", tc.sql, err)
				}
				accepts = myErr.Number != 1064
			}
			if accepts != tc.accept {
				t.Fatalf("OceanBase accepts=%v, want %v for %q", accepts, tc.accept, tc.sql)
			}
			_, err := Parse(tc.sql)
			if omniAccepts := err == nil; omniAccepts != tc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v",
					omniAccepts, tc.accept, tc.sql, err)
			}
		})
	}
}
