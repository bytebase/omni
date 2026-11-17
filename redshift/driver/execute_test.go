package driver

import (
	"strings"
	"testing"
)

// TestCreationStatements pins what a database-creation step accepts: the
// CREATE DATABASE statement, and nothing after it, because a statement that
// followed could not be retried safely.
func TestCreationStatements(t *testing.T) {
	statements, name, err := creationStatements("CREATE DATABASE newdb;")
	if err != nil {
		t.Fatalf("creationStatements: %v", err)
	}
	if name != "newdb" || len(statements) != 1 {
		t.Errorf("statements = %q, name = %q, want one statement creating newdb", statements, name)
	}

	for _, script := range []string{
		"",
		"SELECT 1",
		"CREATE DATABASE newdb; CREATE SCHEMA app;",
	} {
		if _, _, err := creationStatements(script); err == nil {
			t.Errorf("creationStatements(%q) accepted a script that is not one creation", script)
		}
	}
	if _, _, err := creationStatements("CREATE DATABASE newdb; CREATE SCHEMA app;"); err != nil && !strings.Contains(err.Error(), "script of its own") {
		t.Errorf("the refusal does not say where initialization belongs: %v", err)
	}
}

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
		{statement: "COMMIT", want: true},
		{statement: "END", want: true},
		{statement: "BEGIN", want: true},
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
