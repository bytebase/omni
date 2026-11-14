package analysis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/mysql/analysis"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
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
	}
	for _, tc := range tests {
		t.Run(tc.statement, func(t *testing.T) {
			readOnly, returnsData := analysis.ValidateQuery(parse(t, tc.statement))
			require.Equal(t, tc.readOnly, readOnly)
			require.Equal(t, tc.returnsData, returnsData)
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
