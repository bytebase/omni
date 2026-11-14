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
		{statement: "EXPLAIN PLAN FOR UPDATE t1 SET c1 = 1;", readOnly: true, returnsData: true},
		{statement: "WITH x AS (SELECT 'update' AS c1 FROM t1 /* UPDATE */) SELECT c1 FROM x;", readOnly: true, returnsData: true},
		{statement: "SELECT * FROM t1; UPDATE t1 SET c1 = 1;"},
		{statement: "SET DEFINE OFF\nSELECT * FROM t1;"},
		{statement: "CREATE TABLE t1 (c1 INT);"},
		{statement: "BEGIN NULL; END;"},
		{statement: ""},
		{statement: "-- only a comment"},
		{statement: "WITH x AS (SELECT * FROM t1) UPDATE t1 SET c1 = 1;", wantErr: true},
		{statement: "SELECT 1 FROM DUAL;\nSELECT * FROM", wantErr: true},
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
