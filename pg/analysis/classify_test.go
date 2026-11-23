package analysis

import (
	"fmt"
	"testing"

	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
	"github.com/stretchr/testify/require"
)

func parseNodes(t *testing.T, sql string) []ast.Node {
	t.Helper()
	stmts, err := omnipg.Parse(sql)
	require.NoError(t, err)
	var nodes []ast.Node
	for _, stmt := range stmts {
		nodes = append(nodes, stmt.AST)
	}
	return nodes
}

func TestValidateQuery(t *testing.T) {
	tests := []struct {
		sql         string
		readOnly    bool
		returnsData bool
	}{
		{`select* from t`, true, true},
		{`explain select * from t;`, true, true},
		{`explain    analyze select * from t`, true, true},
		{`explain    analyze insert into t values (1)`, false, false},
		{`EXPLAIN ANALYZE WITH cte1 AS (DELETE FROM t RETURNING id) SELECT * FROM cte1`, false, false},
		{`EXPLAIN (ANALYZE off) DELETE FROM t`, true, true},
		{`EXPLAIN (ANALYZE on) DELETE FROM t`, false, false},
		{`
				With t as (
					select * from t1
				), tx as (
					delete from t2
				)
				select * from t;
				`, false, false},
		{`
				With t as (
					select * from t1
				), tx as (
					select * from t1
				)
				update t set a = 1;
				`, false, false},
		{`
				With t as (
					select * from t1
				), tx as (
					select * from t1
				)
				insert into t values (1, 2, 3);
				`, false, false},
		{`WITH x AS (MERGE INTO t a USING t b ON a.id = b.id WHEN MATCHED THEN UPDATE SET id = a.id RETURNING a.secret) SELECT secret FROM x`, false, false},
		{`MERGE INTO t a USING t b ON a.id = b.id WHEN MATCHED THEN UPDATE SET id = a.id RETURNING a.secret`, false, false},
		{`SELECT a INTO t2 FROM t`, false, false},
		{`WITH x AS (VALUES (1), (2)) SELECT * FROM x`, true, true},
		{`WITH x AS (TABLE t) SELECT * FROM x`, true, true},
		{`WITH x AS (SELECT 1 UNION SELECT 2) SELECT * FROM x`, true, true},
		{`WITH RECURSIVE x AS (SELECT 1 UNION SELECT n + 1 FROM x WHERE n < 5) SELECT * FROM x`, true, true},
		{"select * from t where a = 'klasjdfkljsa$tag$; -- lkjdlkfajslkdfj'", true, true},
		{`
				With t as (
					select * from t1 where a = 'insert'
				), tx as (
					select * from "delete"
				) /* UPDATE */
				select "update" from t;
				`, true, true},
		{`create table t (a int);`, false, false},
		{`SELECT * FROM t1; SELECT * FROM t2;`, true, true},
		{`SELECT * FROM t1; SELECT * FROM t2; SHOW search_path;`, true, true},
		{`SET work_mem = '128MB'; SELECT * FROM t1;`, true, false},
		{`SELECT * FROM t1; INSERT INTO t2 VALUES (1);`, false, false},
		{`SELECT * FROM t1; CREATE TABLE t2 (a int);`, false, false},
		{`EXPLAIN SELECT * FROM t1; SELECT * FROM t2;`, true, true},
	}
	for _, tc := range tests {
		readOnly, returnsData := ValidateQuery(parseNodes(t, tc.sql))
		require.Equal(t, tc.readOnly, readOnly, tc.sql)
		require.Equal(t, tc.returnsData, returnsData, tc.sql)
	}
}

func TestStatementTypes(t *testing.T) {
	tests := []struct {
		sql  string
		want []StatementType
	}{
		{"CREATE TABLE t1 (id INT);", []StatementType{CreateTable}},
		{"CREATE VIEW v1 AS SELECT * FROM t1;", []StatementType{CreateView}},
		{"CREATE INDEX idx_name ON t1(name);", []StatementType{CreateIndex}},
		{"CREATE SEQUENCE seq1;", []StatementType{CreateSequence}},
		{"CREATE SCHEMA schema1;", []StatementType{CreateSchema}},
		{"CREATE FUNCTION func1() RETURNS INT AS $$ BEGIN RETURN 1; END; $$ LANGUAGE plpgsql;", []StatementType{CreateFunction}},
		{"DROP TABLE t1;", []StatementType{DropTable}},
		{"DROP INDEX idx_name;", []StatementType{DropIndex}},
		{"DROP SCHEMA schema1;", []StatementType{DropSchema}},
		{"DROP SEQUENCE seq1;", []StatementType{DropSequence}},
		{"DROP EXTENSION postgis;", []StatementType{DropExtension}},
		{"DROP DATABASE testdb;", []StatementType{DropDatabase}},
		{"DROP TYPE custom_type;", []StatementType{DropType}},
		{"DROP TRIGGER trig1 ON t1;", []StatementType{DropTrigger}},
		{"DROP VIEW v1;", []StatementType{DropView}},
		{"DROP FUNCTION func1();", []StatementType{DropFunction}},
		{"CREATE TYPE custom_type AS ENUM ('a', 'b', 'c');", []StatementType{CreateType}},
		{"ALTER TABLE t1 ADD COLUMN name VARCHAR(100);", []StatementType{AlterTable}},
		{"ALTER TABLE t1 ADD CONSTRAINT pk_id PRIMARY KEY (id);", []StatementType{AlterTable}},
		{"ALTER TABLE t1 DROP COLUMN name;", []StatementType{AlterTable}},
		{"ALTER TABLE t1 DROP CONSTRAINT pk_id;", []StatementType{AlterTable}},
		{"ALTER TABLE t1 ALTER COLUMN name TYPE TEXT;", []StatementType{AlterTable}},
		{"ALTER TABLE t1 ALTER COLUMN name DROP DEFAULT;", []StatementType{AlterTable}},
		{"ALTER TABLE t1 ALTER COLUMN name DROP NOT NULL;", []StatementType{AlterTable}},
		{"ALTER VIEW v1 OWNER TO user1;", []StatementType{AlterView}},
		{"ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;", []StatementType{AlterTable}},
		{"ALTER SEQUENCE seq1 RESTART WITH 100;", []StatementType{AlterSequence}},
		{"INSERT INTO t1 (id) VALUES (1);", []StatementType{Insert}},
		{"UPDATE t1 SET name = 'test' WHERE id = 1;", []StatementType{Update}},
		{"DELETE FROM t1 WHERE id = 1;", []StatementType{Delete}},
		{"TRUNCATE TABLE t1;", []StatementType{Truncate}},
		{"COMMENT ON TABLE t1 IS 'test table';", []StatementType{Comment}},
		{"CREATE TABLE t1 (id INT);\nCREATE INDEX idx_id ON t1(id);\nINSERT INTO t1 (id) VALUES (1);", []StatementType{CreateTable, CreateIndex, Insert}},
		{"ALTER TABLE t1 RENAME TO t2;", []StatementType{AlterTable}},
		{"ALTER TABLE t1 RENAME COLUMN old_name TO new_name;", []StatementType{AlterTable}},
		{"ALTER TABLE t1 RENAME old_name TO new_name;", []StatementType{AlterTable}},
		{"ALTER TABLE t1 RENAME CONSTRAINT old_constraint TO new_constraint;", []StatementType{AlterTable}},
		{"ALTER INDEX idx_old RENAME TO idx_new;", []StatementType{RenameIndex}},
		{"ALTER SCHEMA old_schema RENAME TO new_schema;", []StatementType{RenameSchema}},
		{"ALTER SEQUENCE seq_old RENAME TO seq_new;", []StatementType{RenameSequence}},
		{"ALTER VIEW v_old RENAME TO v_new;", []StatementType{AlterView}},
		{"DROP MATERIALIZED VIEW mv1;", []StatementType{DropTable}},
		{"MERGE INTO t1 USING t2 ON t1.id = t2.id WHEN MATCHED THEN DELETE;", []StatementType{Merge}},
		{"WITH d AS (DELETE FROM t1 RETURNING id) SELECT count(*) FROM d;", []StatementType{Delete}},
		{"WITH d AS (DELETE FROM t1 RETURNING *), u AS (UPDATE t3 SET c = 1 RETURNING id) INSERT INTO t2 SELECT * FROM d;", []StatementType{Insert, Delete, Update}},
		{"EXPLAIN ANALYZE MERGE INTO t1 USING t2 ON t1.id = t2.id WHEN MATCHED THEN DELETE;", []StatementType{Merge}},
		{"EXPLAIN (ANALYZE, BUFFERS) WITH d AS (DELETE FROM t1 RETURNING id) UPDATE t2 SET c = 1 FROM d;", []StatementType{Update, Delete}},
		{"EXPLAIN WITH d AS (DELETE FROM t1 RETURNING id) SELECT count(*) FROM d;", nil},
		{"EXPLAIN (ANALYZE FALSE) DELETE FROM t1;", nil},
		{"EXPLAIN (ANALYZE 1) DELETE FROM t1;", []StatementType{Delete}},
		{"CREATE TABLE archive AS WITH d AS (DELETE FROM t1 RETURNING *) SELECT * FROM d;", []StatementType{CreateTable, Delete}},
		{"WITH d AS (DELETE FROM t1 RETURNING *) SELECT * INTO archive FROM d;", []StatementType{CreateTable, Delete}},
		{"CREATE MATERIALIZED VIEW mv AS SELECT * FROM t1;", []StatementType{CreateView}},
		{"CREATE FUNCTION purge() RETURNS bigint LANGUAGE sql BEGIN ATOMIC WITH d AS (DELETE FROM t1 RETURNING 1) SELECT count(*) FROM d; END;", []StatementType{CreateFunction}},
		{"CALL archive_orders(30);", []StatementType{Execute}},
		{"DO $$ BEGIN DELETE FROM t1; END $$;", []StatementType{Execute}},
		{"DO LANGUAGE plpgsql $$ BEGIN PERFORM 1; END $$;", []StatementType{Execute}},
		{"EXECUTE purge_stmt(30);", []StatementType{Execute}},
		{"EXPLAIN ANALYZE EXECUTE purge_stmt(30);", []StatementType{Execute}},
		{"EXPLAIN EXECUTE purge_stmt(30);", nil},
		{"PREPARE purge_stmt(int) AS DELETE FROM t1 WHERE id = $1;", nil},
		{"CREATE PROCEDURE archive_orders(days int) LANGUAGE sql AS $$ DELETE FROM t1 $$;", []StatementType{CreateFunction}},
	}
	for _, tc := range tests {
		var got []StatementType
		for _, node := range parseNodes(t, tc.sql) {
			got = append(got, StatementTypes(node)...)
		}
		require.Equal(t, tc.want, got, tc.sql)
	}
}

func TestUnwrapExplainAnalyze(t *testing.T) {
	tests := []struct {
		sql      string
		wantText string
		wantNode string
	}{
		{"EXPLAIN ANALYZE DELETE FROM t WHERE a = 1", "DELETE FROM t WHERE a = 1", "*ast.DeleteStmt"},
		{"EXPLAIN (ANALYZE, BUFFERS) WITH c AS (SELECT 1) SELECT * FROM c ORDER BY 1", "WITH c AS (SELECT 1) SELECT * FROM c ORDER BY 1", "*ast.SelectStmt"},
		{"EXPLAIN ANALYZE CREATE TABLE n AS WITH c AS (SELECT 1) SELECT * FROM c", "CREATE TABLE n AS WITH c AS (SELECT 1) SELECT * FROM c", "*ast.CreateTableAsStmt"},
		{"EXPLAIN DELETE FROM t", "EXPLAIN DELETE FROM t", "*ast.ExplainStmt"},
		{"EXPLAIN (ANALYZE false) DELETE FROM t", "EXPLAIN (ANALYZE false) DELETE FROM t", "*ast.ExplainStmt"},
		{"DELETE FROM t", "DELETE FROM t", "*ast.DeleteStmt"},
	}
	for _, tc := range tests {
		nodes := parseNodes(t, tc.sql)
		require.Len(t, nodes, 1)
		node, text := UnwrapExplainAnalyze(nodes[0], tc.sql)
		require.Equal(t, tc.wantText, text, tc.sql)
		require.Equal(t, tc.wantNode, fmt.Sprintf("%T", node), tc.sql)
	}

	// Text that does not hold the statement keeps the EXPLAIN's text.
	nodes := parseNodes(t, "EXPLAIN ANALYZE DELETE FROM t")
	node, text := UnwrapExplainAnalyze(nodes[0], "")
	require.Equal(t, "", text)
	require.IsType(t, &ast.DeleteStmt{}, node)
}
