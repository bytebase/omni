package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/oracle/parser"
)

func TestValidateQuery(t *testing.T) {
	tests := []struct {
		statement   string
		readOnly    bool
		returnsData bool
		wantErr     bool
	}{
		{statement: "SELECT * FROM t1;", readOnly: true, returnsData: true},
		{statement: "SELECT * FROM t1; SELECT * FROM t2;", readOnly: true, returnsData: true},
		{statement: "EXPLAIN PLAN FOR SELECT * FROM t1;", readOnly: true, returnsData: true},
		{statement: "EXPLAIN PLAN FOR UPDATE t1 SET c1 = 1;", readOnly: true, returnsData: true},
		{statement: "WITH x AS (SELECT * FROM t1) SELECT * FROM x;", readOnly: true, returnsData: true},
		{statement: "WITH x AS (SELECT 'update' AS c1 FROM t1 /* UPDATE */) SELECT c1 FROM x;", readOnly: true, returnsData: true},
		{statement: "SELECT * FROM t1; UPDATE t1 SET c1 = 1;"},
		{statement: "UPDATE t1 SET c1 = 1;"},
		{statement: "SET DEFINE OFF\nSELECT * FROM t1;"},
		{statement: "CREATE TABLE t1 (c1 INT);"},
		{statement: "BEGIN NULL; END;"},
		{statement: ""},
		{statement: "-- only a comment"},
		{statement: "WITH x AS (SELECT * FROM t1) UPDATE t1 SET c1 = 1;", wantErr: true},
		{statement: "SELECT 1 FROM DUAL;\nSELECT * FROM", wantErr: true},
		// An inline PL/SQL declaration can write, through an autonomous
		// transaction or dynamic SQL, while the query only reads.
		{statement: "WITH FUNCTION f RETURN NUMBER IS PRAGMA AUTONOMOUS_TRANSACTION; BEGIN DELETE FROM t1; COMMIT; RETURN 1; END; SELECT f FROM dual"},
		{statement: "WITH FUNCTION f RETURN NUMBER IS BEGIN EXECUTE IMMEDIATE 'DROP TABLE t1'; RETURN 1; END; SELECT f FROM dual"},
		{statement: "WITH PROCEDURE p IS BEGIN NULL; END; FUNCTION f RETURN NUMBER IS BEGIN p; RETURN 1; END; SELECT f FROM dual"},
		{statement: "WITH FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; SELECT f FROM dual"},
		{statement: "SELECT * FROM (WITH FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; SELECT f FROM dual)"},
		{statement: "EXPLAIN PLAN FOR WITH FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; SELECT f FROM dual"},
	}
	for _, test := range tests {
		t.Run(test.statement, func(t *testing.T) {
			readOnly, returnsData, err := ValidateQuery(test.statement)
			if test.wantErr {
				require.ErrorAs(t, err, new(*parser.ParseError))
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.readOnly, readOnly)
			require.Equal(t, test.returnsData, returnsData)
		})
	}
}
