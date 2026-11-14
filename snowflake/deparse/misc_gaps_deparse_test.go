package deparse_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/deparse"
	"github.com/bytebase/omni/snowflake/parser"
	"github.com/stretchr/testify/require"
)

// CREATE PROCEDURE SCOPED { TEMP | TEMPORARY } TABLE and DROP DATABASE ROLE
// round-trips.
func TestDeparse_ProcedureScopedTableAndDropDatabaseRole(t *testing.T) {
	for _, sql := range []string{
		"CREATE PROCEDURE SCOPED TEMPORARY TABLE t (a INT)",
		"CREATE OR REPLACE PROCEDURE SCOPED TEMPORARY TABLE db.s.t (a INT, \"B\" NUMBER(3, 1))",
		"DROP DATABASE ROLE d1.dr1",
		"DROP DATABASE ROLE IF EXISTS dr1",
	} {
		assertRoundTrip(t, sql)
	}
	f, err := parser.Parse("create procedure scoped temp table t (a int)")
	require.NoError(t, err)
	out, err := deparse.DeparseFile(f)
	require.NoError(t, err)
	require.Equal(t, "CREATE PROCEDURE SCOPED TEMPORARY TABLE t (a int)", out)
}
