package analysis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/oceanbase/analysis"
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
)

func parse(t *testing.T, statement string) []ast.Node {
	t.Helper()
	list, err := parser.Parse(statement)
	require.NoError(t, err, statement)
	return list.Items
}

func TestValidateQuery(t *testing.T) {
	tests := []struct {
		statement   string
		readOnly    bool
		returnsData bool
	}{
		{"SELECT a FROM t", true, true},
		{"SELECT a FROM t; SHOW TABLES", true, true},
		{"SET @a = 1", true, false},
		{"SELECT a INTO @x FROM t", true, false},
		{"SELECT a FROM t INTO OUTFILE '/tmp/x'", false, false},
		{"EXPLAIN SELECT a FROM t", true, true},
		{"EXPLAIN DELETE FROM t", true, true},
		{"INSERT INTO t VALUES (1)", false, false},
		{"SELECT 1; DELETE FROM t", false, false},
		{"SHOW CREATE TABLE bytebase;", true, true},
		{"DESC bytebase;", true, true},
		{"SELECT * FROM t1 WHERE c1 = 1; SELECT * FROM t2;", true, true},
		{"CREATE TABLE t1 (c1 INT);", false, false},
		{"UPDATE t1 SET c1 = 1;", false, false},
		{"EXPLAIN FORMAT=JSON DELETE FROM t1;", true, true},
		{"select* from t", true, true},
		{"WITH t AS (SELECT * FROM t1), tx AS (SELECT * FROM t1) UPDATE t SET a = 1;", false, false},
		{"WITH t AS (SELECT * FROM t1), tx AS (SELECT * FROM t1) INSERT INTO t VALUES (1, 2, 3);", false, false},
		{"select * from t where a = 'klasjdfkljsa$tag$; -- lkjdlkfajslkdfj'", true, true},
		{"With t as (select * from t1 where a = 'insert'), tx as (select * from `delete`) /* UPDATE */ select `update` from t;", true, true},
		{"SET max_execution_time = 1000; select * from t", true, false},
	}
	for _, tc := range tests {
		t.Run(tc.statement, func(t *testing.T) {
			readOnly, returnsData := analysis.ValidateQuery(parse(t, tc.statement))
			require.Equal(t, tc.readOnly, readOnly)
			require.Equal(t, tc.returnsData, returnsData)
		})
	}
}

func TestValidateQueryExplainIntoAndExportGate(t *testing.T) {
	for _, tc := range []struct {
		statement string
		wantValid bool
	}{
		{"EXPLAIN INTO plan_tab SELECT * FROM t", false},
		{"EXPLAIN INTO plan_tab SET STATEMENT_ID = 'x' SELECT * FROM t", false},
		{"EXPLAIN SET STATEMENT_ID = 'x' SELECT * FROM t", true},
		{"EXPLAIN SELECT * FROM t", true},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			readOnly, _ := analysis.ValidateQuery(parse(t, tc.statement))
			require.Equal(t, tc.wantValid, readOnly)

			span, err := analysis.ResolveQuerySpan(tc.statement, analysis.ResolveOptions{DefaultDatabase: "db"})
			require.NoError(t, err)
			wantType := analysis.Explain
			if !tc.wantValid {
				wantType = analysis.DML
			}
			require.Equal(t, wantType, span.Type)
		})
	}

	for _, tc := range []struct {
		name      string
		statement string
		wantValid bool
	}{
		{"bare export", "SELECT a FROM t INTO OUTFILE '/tmp/x'", false},
		{"bare dump", "SELECT a FROM t INTO DUMPFILE '/tmp/x'", false},
		{"export second arm", "SELECT 1 UNION SELECT a FROM t INTO OUTFILE '/tmp/x'", false},
		{"export braced arm", "(SELECT a FROM t INTO OUTFILE '/tmp/x') UNION (SELECT 1)", false},
		{"export subquery", "SELECT * FROM (SELECT a FROM t INTO OUTFILE '/tmp/x') y", false},
		{"export plain EXPLAIN", "EXPLAIN SELECT a FROM t INTO OUTFILE '/tmp/x'", false},
		{"dump plain EXPLAIN", "EXPLAIN SELECT a FROM t INTO DUMPFILE '/tmp/x'", false},
		{"export arm plain EXPLAIN", "EXPLAIN SELECT 1 UNION SELECT a FROM t INTO OUTFILE '/tmp/x'", false},
		{"plain set op", "SELECT 1 UNION SELECT 2", true},
		{"variable assignment", "SELECT a FROM t INTO @v", true},
		{"variable assignment plain EXPLAIN", "EXPLAIN SELECT a FROM t INTO @v", true},
		{"string literal", "SELECT 'INTO OUTFILE' FROM t", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			readOnly, _ := analysis.ValidateQuery(parse(t, tc.statement))
			require.Equal(t, tc.wantValid, readOnly, tc.statement)
		})
	}
}

func TestStatementTypeOf(t *testing.T) {
	tests := []struct {
		statement string
		want      analysis.StatementType
	}{
		{"CREATE TABLE t (a INT)", analysis.CreateTable},
		{"CREATE DATABASE d", analysis.CreateDatabase},
		{"CREATE INDEX i ON t (a)", analysis.CreateIndex},
		{"CREATE MATERIALIZED VIEW mv AS SELECT 1", analysis.CreateView},
		{"ALTER MATERIALIZED VIEW mv REFRESH FAST", analysis.AlterView},
		{"DROP MATERIALIZED VIEW mv", analysis.DropView},
		{"DROP TABLE t", analysis.DropTable},
		{"ALTER TABLE t ADD COLUMN b INT", analysis.AlterTable},
		{"TRUNCATE TABLE t", analysis.Truncate},
		{"RENAME TABLE t TO u", analysis.Rename},
		{"INSERT INTO t VALUES (1)", analysis.Insert},
		{"UPDATE t SET a = 1", analysis.Update},
		{"DELETE FROM t", analysis.Delete},
		{"CALL p()", analysis.Execute},
		{"SELECT 1", analysis.StatementTypeUnspecified},
		{"PREPARE s FROM 'SELECT 1'", analysis.StatementTypeUnspecified},
	}
	for _, tc := range tests {
		t.Run(tc.statement, func(t *testing.T) {
			nodes := parse(t, tc.statement)
			require.Len(t, nodes, 1)
			require.Equal(t, tc.want, analysis.StatementTypeOf(nodes[0]))
		})
	}
}
