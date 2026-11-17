package driver

import (
	"strings"
	"testing"
)

// TestRefuseImplicitCommit pins which statements a transactional script may not
// hold: the ones Redshift commits on its own, so a rollback would not undo them.
func TestRefuseImplicitCommit(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      bool
	}{
		{statement: "INSERT INTO facts VALUES (1)"},
		{statement: "CREATE TABLE t (a int)"},
		{statement: "TRUNCATE facts", want: true},
		{statement: "vacuum full facts", want: true},
		{statement: "CREATE DATABASE other", want: true},
		{statement: "DROP DATABASE other", want: true},
	} {
		err := refuseImplicitCommit(tc.statement)
		if tc.want && err == nil {
			t.Errorf("refuseImplicitCommit(%q) accepted a statement Redshift commits on its own", tc.statement)
		}
		if !tc.want && err != nil {
			t.Errorf("refuseImplicitCommit(%q) = %v, want no error", tc.statement, err)
		}
		if tc.want && err != nil && !strings.Contains(err.Error(), "autocommit") {
			t.Errorf("refuseImplicitCommit(%q) = %v, want the autocommit hint", tc.statement, err)
		}
	}
}
