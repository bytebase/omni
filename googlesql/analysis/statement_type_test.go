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
		{"TRUNCATE TABLE users;", Truncate},
		// DDL - create.
		{"CREATE TABLE users (id INT64, name STRING);", CreateTable},
		{"CREATE OR REPLACE TABLE users (id INT64);", CreateTable},
		{"CREATE SNAPSHOT TABLE snap CLONE users;", CreateTable},
		{"CREATE VIEW v AS SELECT 1;", CreateView},
		{"CREATE MATERIALIZED VIEW mv AS SELECT 1;", CreateView},
		{"CREATE SCHEMA ds;", CreateSchema},
		{"CREATE SEARCH INDEX idx ON users(name);", CreateIndex},
		{"CREATE FUNCTION f(x INT64) AS (x + 1);", CreateFunction},
		{"CREATE PROCEDURE p() BEGIN SELECT 1; END;", CreateProcedure},
		// DDL - alter.
		{"ALTER TABLE `users` ADD COLUMN status STRING;", AlterTable},
		{"ALTER TABLE users SET OPTIONS(description='x');", AlterTable},
		{"ALTER VIEW v SET OPTIONS(description='x');", AlterView},
		{"ALTER MATERIALIZED VIEW mv SET OPTIONS(enable_refresh=true);", AlterView},
		// DDL - drop.
		{"DROP TABLE users;", DropTable},
		{"DROP VIEW v;", DropView},
		{"DROP SCHEMA ds;", DropSchema},
		{"DROP MATERIALIZED VIEW mv;", DropView},
		{"DROP FUNCTION f;", DropFunction},
		{"DROP PROCEDURE p;", DropProcedure},
		{"DROP SEARCH INDEX idx ON users;", DropIndex},
		{"DROP SNAPSHOT TABLE snap;", DropTable},
		// Procedure calls and dynamic SQL.
		{"CALL ds.proc(1, 'a');", Execute},
		{"CALL `proj.ds.proc`();", Execute},
		{"CALL BQ.REFRESH_MATERIALIZED_VIEW('proj.ds.mv');", Execute},
		{"EXECUTE IMMEDIATE 'DROP TABLE ds.t';", Execute},
		{"EXECUTE IMMEDIATE \"SELECT @a\" INTO x USING 1 AS a;", Execute},
		{"EXECUTE IMMEDIATE CONCAT('DROP ', 'TABLE t');", Execute},
		// No enum value; see the deferral on StatementTypeOf.
		{"ALTER SCHEMA ds SET OPTIONS(description='x');", StatementTypeUnspecified},
		{"DROP ROW ACCESS POLICY rap ON users;", StatementTypeUnspecified},
		{"BEGIN CALL ds.proc(); END;", StatementTypeUnspecified},
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

// TestStatementTypeAgreesWithClassifier is the acceptance test: ClassifySQL
// derives DML membership by a different route than the AST switch, so a
// statement it calls DML must map onto a DML enum value. A disagreement means an
// approval rule keyed on statement.sql_type would mis-gate that statement.
func TestStatementTypeAgreesWithClassifier(t *testing.T) {
	dml := map[StatementType]bool{
		Insert:   true,
		Update:   true,
		Delete:   true,
		Merge:    true,
		Truncate: true,
	}
	statements := []string{
		"INSERT INTO users (id) VALUES (1);",
		"UPDATE users SET name = 'b' WHERE id = 1;",
		"DELETE FROM users WHERE id = 1;",
		"MERGE INTO users t USING staging s ON t.id = s.id WHEN MATCHED THEN UPDATE SET t.name = s.name;",
		"TRUNCATE TABLE users;",
		"CREATE TABLE users (id INT64);",
		"ALTER TABLE users ADD COLUMN status STRING;",
		"DROP TABLE users;",
		"CREATE VIEW v AS SELECT 1;",
		"SELECT * FROM users;",
	}
	for _, statement := range statements {
		t.Run(statement, func(t *testing.T) {
			got := statementTypesOf(t, statement)
			require.Len(t, got, 1)
			classified := ClassifySQL(statement, DialectBigQuery)
			require.Equal(t, classified == DML, dml[got[0]],
				"classifier says %v, mapping says %v", classified, got[0])
		})
	}
}

func TestStatementTypeOfMultiStatement(t *testing.T) {
	got := statementTypesOf(t, "ALTER TABLE users ADD COLUMN status STRING;\nINSERT INTO users (id) VALUES (1);")
	require.Equal(t, []StatementType{AlterTable, Insert}, got)
}

func TestStatementTypeOfBigQueryExternalDrops(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want StatementType
	}{
		{"DROP EXTERNAL TABLE ds.ext;", DropTable},
		{"DROP TABLE FUNCTION ds.tf;", DropFunction},
	} {
		file, err := parser.ParseWithDialect(tc.sql, parser.DialectBigQuery)
		require.NoError(t, err, tc.sql)
		require.Len(t, file.Stmts, 1)
		require.Equal(t, tc.want, StatementTypeOf(file.Stmts[0]), tc.sql)
	}
}

func TestStatementTypeOfBigQuerySnapshots(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want StatementType
	}{
		{"CREATE SNAPSHOT TABLE ds.s CLONE ds.t;", CreateTable},
		{"CREATE SNAPSHOT SCHEMA ds2 CLONE ds;", CreateSchema},
	} {
		file, err := parser.ParseWithDialect(tc.sql, parser.DialectBigQuery)
		require.NoError(t, err, tc.sql)
		require.Len(t, file.Stmts, 1)
		require.Equal(t, tc.want, StatementTypeOf(file.Stmts[0]), tc.sql)
	}
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
