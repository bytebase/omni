package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/googlesql/parser"
)

func statementTypesOf(t *testing.T, statement string) []StatementType {
	t.Helper()
	file, _ := parser.Parse(statement)
	var got []StatementType
	if file != nil {
		for _, node := range file.Stmts {
			got = append(got, StatementTypeOf(node))
		}
	}
	return got
}

func TestStatementTypeOf(t *testing.T) {
	tests := []struct {
		statement string
		want      StatementType
	}{
		// DML.
		{"INSERT INTO `users` (id, name) VALUES (1, 'a');", Insert},
		{"INSERT INTO ds.users SELECT * FROM ds.staging;", Insert},
		{"UPDATE users SET name = 'b' WHERE id = 1;", Update},
		{"DELETE FROM users WHERE id = 1;", Delete},
		{"MERGE INTO users t USING staging s ON t.id = s.id WHEN MATCHED THEN UPDATE SET t.name = s.name;", Merge},
		// DDL - create.
		{"CREATE VIEW v AS SELECT 1;", CreateView},
		{"CREATE SCHEMA ds;", CreateSchema},
		{"CREATE SEARCH INDEX idx ON users(name);", CreateIndex},
		// DDL - alter.
		{"ALTER TABLE `users` ADD COLUMN status STRING;", AlterTable},
		{"ALTER TABLE users SET OPTIONS(description='x');", AlterTable},
		{"ALTER VIEW v SET OPTIONS(description='x');", AlterView},
		// DDL - drop.
		{"DROP TABLE users;", DropTable},
		{"DROP VIEW v;", DropView},
		{"DROP SCHEMA ds;", DropSchema},
		{"DROP FUNCTION f;", DropFunction},
		// Procedure calls and dynamic SQL.
		{"CALL ds.proc(1, 'a');", Execute},
		{"CALL `proj.ds.proc`();", Execute},
		{"CALL BQ.REFRESH_MATERIALIZED_VIEW('proj.ds.mv');", Execute},
		{"EXECUTE IMMEDIATE 'DROP TABLE ds.t';", Execute},
		{"EXECUTE IMMEDIATE \"SELECT @a\" INTO x USING 1 AS a;", Execute},
		{"EXECUTE IMMEDIATE CONCAT('DROP ', 'TABLE t');", Execute},
		// No enum value; see the deferral on StatementTypeOf.
		{"ALTER SCHEMA ds SET OPTIONS(description='x');", StatementTypeUnspecified},
		{"SELECT * FROM users;", StatementTypeUnspecified},
	}
	for _, tt := range tests {
		t.Run(tt.statement, func(t *testing.T) {
			got := statementTypesOf(t, tt.statement)
			require.Len(t, got, 1)
			require.Equal(t, tt.want, got[0])
		})
	}
}

func TestStatementTypeOfMultiStatement(t *testing.T) {
	got := statementTypesOf(t, "ALTER TABLE users ADD COLUMN status STRING;\nINSERT INTO users (id) VALUES (1);")
	require.Equal(t, []StatementType{AlterTable, Insert}, got)
}

func TestStatementTypeOfSpannerCall(t *testing.T) {
	got := statementTypesOf(t, "CALL cancel_query('12345');")
	require.Equal(t, []StatementType{Execute}, got)
}

// An unparseable statement yields no node, and StatementTypeOf reports
// Unspecified for it rather than guessing, so a rule selecting DML skips it.
func TestStatementTypeOfNil(t *testing.T) {
	require.Equal(t, StatementTypeUnspecified, StatementTypeOf(nil))
}
