package parser

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bytebase/omni/oracle/internal/testcontainer"
)

// openOracleSystem returns t's session to the pinned Oracle as SYSTEM: a
// replay of its tape, or the server in a recording.
func openOracleSystem(t *testing.T) *sql.DB {
	t.Helper()
	return testcontainer.Open(t, "system")
}

// openOracleUser makes a user of the given name, as SYSTEM on t's tape,
// that may create tables, views, sequences, synonyms, types, procedures and
// triggers, and returns t's session as that user.
func openOracleUser(t *testing.T, name string) *sql.DB {
	t.Helper()
	admin := openOracleSystem(t)
	for _, stmt := range []string{
		"CREATE USER " + name + " IDENTIFIED BY " + testcontainer.Password + " QUOTA UNLIMITED ON users",
		"GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE TYPE, CREATE PROCEDURE, CREATE TRIGGER TO " + name,
	} {
		if _, err := admin.ExecContext(context.Background(), stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	return testcontainer.Open(t, name)
}
