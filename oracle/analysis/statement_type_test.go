package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/oracle/ast"
)

func TestStatementTypeOf(t *testing.T) {
	tests := []struct {
		statement string
		want      StatementType
	}{
		{statement: "CREATE TABLE t (a NUMBER)", want: CreateTable},
		{statement: "CREATE INDEX i ON t (a)", want: CreateIndex},
		{statement: "CREATE VIEW v AS SELECT 1 FROM DUAL", want: CreateView},
		{statement: "CREATE SEQUENCE s START WITH 1", want: CreateSequence},
		{statement: "CREATE OR REPLACE TYPE ty AS OBJECT (a NUMBER)", want: CreateType},
		{statement: "CREATE DATABASE db", want: CreateDatabase},
		{statement: "ALTER DATABASE OPEN", want: AlterDatabase},
		{statement: "DROP DATABASE", want: DropDatabase},
		{statement: "DROP TABLE t", want: DropTable},
		{statement: "DROP MATERIALIZED VIEW mv", want: DropView},
		{statement: "DROP TYPE BODY ty", want: DropType},
		{statement: "DROP SYNONYM s", want: StatementTypeUnspecified},
		{statement: "ALTER TABLE t ADD (b NUMBER)", want: AlterTable},
		{statement: "TRUNCATE TABLE t", want: Truncate},
		{statement: "RENAME t TO u", want: Rename},
		{statement: "COMMENT ON TABLE t IS 'x'", want: Comment},
		{statement: "INSERT INTO t VALUES (1)", want: Insert},
		{statement: "UPDATE t SET a = 1", want: Update},
		{statement: "DELETE FROM t", want: Delete},
		{statement: "CALL p()", want: Execute},
		{statement: "BEGIN NULL; END;", want: Execute},
		{statement: "SELECT 1 FROM DUAL", want: StatementTypeUnspecified},
		{statement: "COMMIT", want: StatementTypeUnspecified},
	}
	for _, test := range tests {
		t.Run(test.statement, func(t *testing.T) {
			list, err := Parse(test.statement)
			require.NoError(t, err)
			require.Len(t, list.Items, 1)
			raw, ok := list.Items[0].(*ast.RawStmt)
			require.True(t, ok)
			require.Equal(t, test.want, StatementTypeOf(raw.Stmt))
		})
	}
}
