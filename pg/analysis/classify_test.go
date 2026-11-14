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
		{"CREATE TABLE t1 (id INT);", []StatementType{"CREATE_TABLE"}},
		{"CREATE VIEW v1 AS SELECT * FROM t1;", []StatementType{"CREATE_VIEW"}},
		{"CREATE INDEX idx_name ON t1(name);", []StatementType{"CREATE_INDEX"}},
		{"CREATE SEQUENCE seq1;", []StatementType{"CREATE_SEQUENCE"}},
		{"CREATE SCHEMA schema1;", []StatementType{"CREATE_SCHEMA"}},
		{"CREATE FUNCTION func1() RETURNS INT AS $$ BEGIN RETURN 1; END; $$ LANGUAGE plpgsql;", []StatementType{"CREATE_FUNCTION"}},
		{"DROP TABLE t1;", []StatementType{"DROP_TABLE"}},
		{"DROP INDEX idx_name;", []StatementType{"DROP_INDEX"}},
		{"DROP SCHEMA schema1;", []StatementType{"DROP_SCHEMA"}},
		{"DROP SEQUENCE seq1;", []StatementType{"DROP_SEQUENCE"}},
		{"DROP EXTENSION postgis;", []StatementType{"DROP_EXTENSION"}},
		{"DROP DATABASE testdb;", []StatementType{"DROP_DATABASE"}},
		{"DROP TYPE custom_type;", []StatementType{"DROP_TYPE"}},
		{"DROP TRIGGER trig1 ON t1;", []StatementType{"DROP_TRIGGER"}},
		{"DROP VIEW v1;", []StatementType{"DROP_VIEW"}},
		{"DROP FUNCTION func1();", []StatementType{"DROP_FUNCTION"}},
		{"CREATE TYPE custom_type AS ENUM ('a', 'b', 'c');", []StatementType{"CREATE_TYPE"}},
		{"ALTER TABLE t1 ADD COLUMN name VARCHAR(100);", []StatementType{"ALTER_TABLE"}},
		{"ALTER TABLE t1 ADD CONSTRAINT pk_id PRIMARY KEY (id);", []StatementType{"ALTER_TABLE"}},
		{"ALTER TABLE t1 DROP COLUMN name;", []StatementType{"ALTER_TABLE"}},
		{"ALTER TABLE t1 DROP CONSTRAINT pk_id;", []StatementType{"ALTER_TABLE"}},
		{"ALTER TABLE t1 ALTER COLUMN name TYPE TEXT;", []StatementType{"ALTER_TABLE"}},
		{"ALTER TABLE t1 ALTER COLUMN name DROP DEFAULT;", []StatementType{"ALTER_TABLE"}},
		{"ALTER TABLE t1 ALTER COLUMN name DROP NOT NULL;", []StatementType{"ALTER_TABLE"}},
		{"ALTER VIEW v1 OWNER TO user1;", []StatementType{"ALTER_VIEW"}},
		{"ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;", []StatementType{"ALTER_TABLE"}},
		{"ALTER SEQUENCE seq1 RESTART WITH 100;", []StatementType{"ALTER_SEQUENCE"}},
		{"INSERT INTO t1 (id) VALUES (1);", []StatementType{"INSERT"}},
		{"UPDATE t1 SET name = 'test' WHERE id = 1;", []StatementType{"UPDATE"}},
		{"DELETE FROM t1 WHERE id = 1;", []StatementType{"DELETE"}},
		{"TRUNCATE TABLE t1;", []StatementType{"TRUNCATE"}},
		{"COMMENT ON TABLE t1 IS 'test table';", []StatementType{"COMMENT"}},
		{"CREATE TABLE t1 (id INT);\nCREATE INDEX idx_id ON t1(id);\nINSERT INTO t1 (id) VALUES (1);", []StatementType{"CREATE_TABLE", "CREATE_INDEX", "INSERT"}},
		{"ALTER TABLE t1 RENAME TO t2;", []StatementType{"ALTER_TABLE"}},
		{"ALTER TABLE t1 RENAME COLUMN old_name TO new_name;", []StatementType{"ALTER_TABLE"}},
		{"ALTER TABLE t1 RENAME old_name TO new_name;", []StatementType{"ALTER_TABLE"}},
		{"ALTER TABLE t1 RENAME CONSTRAINT old_constraint TO new_constraint;", []StatementType{"ALTER_TABLE"}},
		{"ALTER INDEX idx_old RENAME TO idx_new;", []StatementType{"RENAME_INDEX"}},
		{"ALTER SCHEMA old_schema RENAME TO new_schema;", []StatementType{"RENAME_SCHEMA"}},
		{"ALTER SEQUENCE seq_old RENAME TO seq_new;", []StatementType{"RENAME_SEQUENCE"}},
		{"ALTER VIEW v_old RENAME TO v_new;", []StatementType{"ALTER_VIEW"}},
		{"DROP MATERIALIZED VIEW mv1;", []StatementType{"DROP_TABLE"}},
		{"MERGE INTO t1 USING t2 ON t1.id = t2.id WHEN MATCHED THEN DELETE;", []StatementType{"MERGE"}},
		{"WITH d AS (DELETE FROM t1 RETURNING id) SELECT count(*) FROM d;", []StatementType{"DELETE"}},
		{"WITH d AS (DELETE FROM t1 RETURNING *), u AS (UPDATE t3 SET c = 1 RETURNING id) INSERT INTO t2 SELECT * FROM d;", []StatementType{"INSERT", "DELETE", "UPDATE"}},
		{"EXPLAIN ANALYZE MERGE INTO t1 USING t2 ON t1.id = t2.id WHEN MATCHED THEN DELETE;", []StatementType{"MERGE"}},
		{"EXPLAIN (ANALYZE, BUFFERS) WITH d AS (DELETE FROM t1 RETURNING id) UPDATE t2 SET c = 1 FROM d;", []StatementType{"UPDATE", "DELETE"}},
		{"EXPLAIN WITH d AS (DELETE FROM t1 RETURNING id) SELECT count(*) FROM d;", nil},
		{"EXPLAIN (ANALYZE FALSE) DELETE FROM t1;", nil},
		{"EXPLAIN (ANALYZE 1) DELETE FROM t1;", []StatementType{"DELETE"}},
		{"CREATE TABLE archive AS WITH d AS (DELETE FROM t1 RETURNING *) SELECT * FROM d;", []StatementType{"CREATE_TABLE", "DELETE"}},
		{"WITH d AS (DELETE FROM t1 RETURNING *) SELECT * INTO archive FROM d;", []StatementType{"CREATE_TABLE", "DELETE"}},
		{"CREATE MATERIALIZED VIEW mv AS SELECT * FROM t1;", []StatementType{"CREATE_VIEW"}},
		{"CREATE FUNCTION purge() RETURNS bigint LANGUAGE sql BEGIN ATOMIC WITH d AS (DELETE FROM t1 RETURNING 1) SELECT count(*) FROM d; END;", []StatementType{"CREATE_FUNCTION"}},
		{"CALL archive_orders(30);", []StatementType{"EXECUTE"}},
		{"DO $$ BEGIN DELETE FROM t1; END $$;", []StatementType{"EXECUTE"}},
		{"DO LANGUAGE plpgsql $$ BEGIN PERFORM 1; END $$;", []StatementType{"EXECUTE"}},
		{"EXECUTE purge_stmt(30);", []StatementType{"EXECUTE"}},
		{"EXPLAIN ANALYZE EXECUTE purge_stmt(30);", []StatementType{"EXECUTE"}},
		{"EXPLAIN EXECUTE purge_stmt(30);", nil},
		{"PREPARE purge_stmt(int) AS DELETE FROM t1 WHERE id = $1;", nil},
		{"CREATE PROCEDURE archive_orders(days int) LANGUAGE sql AS $$ DELETE FROM t1 $$;", []StatementType{"CREATE_FUNCTION"}},
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
