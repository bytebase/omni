package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

func TestSharedFixes(t *testing.T) {
	valid := []struct {
		sql string
		tag ast.NodeTag
	}{
		{"@1@{ a = NULL } BEGIN", ast.T_TransactionStmt},
		{"SELECT NULL.*", ast.T_QueryStmt},
		{"SELECT NULL.* EXCEPT(a)", ast.T_QueryStmt},
		{"|> EXTEND NULL.*", ast.T_SubpipelineStmt},
		{"|> AGGREGATE NULL.*", ast.T_SubpipelineStmt},
		{"EXPLAIN BEGIN", ast.T_ExplainStmt},
		{"EXPLAIN SELECT 1", ast.T_ExplainStmt},
		{"DEFINE TABLE a()", ast.T_DefineTableStmt},
		{"DEFINE TABLE ds.t(format = 'CSV')", ast.T_DefineTableStmt},
		{"CREATE MODEL a", ast.T_CreateModelStmt},
		{"CREATE OR REPLACE MODEL IF NOT EXISTS ds.m REMOTE", ast.T_CreateModelStmt},
		{"CREATE MODEL ds.m INPUT(x INT64) OUTPUT(y FLOAT64) REMOTE", ast.T_CreateModelStmt},
		{"DROP MODEL a", ast.T_DropStmt},
		{"DROP MODEL IF EXISTS ds.m", ast.T_DropStmt},
		{"CREATE PLACEMENT us_p OPTIONS(instance_partition = 'us')", ast.T_CreateEntityStmt},
		{"DROP PLACEMENT us_p", ast.T_BQDropStmt},
	}
	for _, tc := range valid {
		file, err := Parse(tc.sql)
		if err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", tc.sql, err)
			continue
		}
		if len(file.Stmts) != 1 {
			t.Errorf("Parse(%q) got %d stmts, want 1", tc.sql, len(file.Stmts))
			continue
		}
		if file.Stmts[0].Tag() != tc.tag {
			t.Errorf("Parse(%q) tag = %v, want %v", tc.sql, file.Stmts[0].Tag(), tc.tag)
		}
	}

	invalid := []string{
		"SELECT (NULL.*)",
		"SELECT TRUE AND NULL.*",
		"SELECT GRAPH_TABLE",
		"SELECT LATERAL",
		"SELECT NULL.IGNORE",
		"DELETE a - a",
		"CREATE TABLE a()",
		"CREATE a a",
		"ALTER a SET AS 'a'",
		"DROP a a",
		"CREATE LOCALITY a",
		"CREATE MODEL a INPUT() OUTPUT()",
	}
	for _, sql := range invalid {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) expected error, got nil", sql)
		}
	}
}

func TestDialectBigQuery(t *testing.T) {
	valid := []string{
		// FROM-first pipe queries
		"FROM a",
		"FROM a |> WHERE x > 1",
		"WITH c AS (SELECT 1) FROM c",
		"FROM (SELECT NULL)",
		"FROM UNNEST(NULL)",
		"FROM UNNEST(NULL AS a, a => NULL)",
		"FROM a JOIN a ON NULL",
		"FROM (a JOIN a)",
		"FROM LATERAL (SELECT NULL)",
		"FROM LATERAL a()",
		"FROM GRAPH_TABLE(a MATCH -)",
		"FROM GRAPH_TABLE(a MATCH - COLUMNS(NULL))",
		"FROM GRAPH_TABLE(a RETURN NULL) a",
		"FROM a PIVOT(NULL FOR NULL IN(NULL))",
		"FROM a UNPIVOT(a FOR a IN(a))",
		"FROM a MATCH_RECOGNIZE(ORDER BY NULL MEASURES NULL AS a PATTERN() DEFINE a AS NULL)",
		"FROM a FOR SYSTEM_TIME AS OF NULL",
		"FROM a WITH OFFSET",
		// BigQuery ML DDL
		"CREATE TEMP MODEL a OPTIONS(model_type='LINEAR_REG') AS SELECT 1",
		"CREATE MODEL a TRANSFORM(NULL) REMOTE WITH CONNECTION a OPTIONS() AS (a AS (SELECT NULL), b AS (SELECT NULL))",
		// External tables & schemas
		"CREATE EXTERNAL TABLE a WITH PARTITION COLUMNS OPTIONS()",
		"CREATE EXTERNAL TABLE a LIKE b DEFAULT COLLATE 'und:ci' WITH CONNECTION a OPTIONS()",
		"CREATE EXTERNAL SCHEMA a WITH CONNECTION a OPTIONS()",
		"DROP EXTERNAL TABLE a",
		"DROP EXTERNAL TABLE FUNCTION a",
		// Connections & constants
		"CREATE CONNECTION a OPTIONS()",
		"DROP CONNECTION a",
		"CREATE TEMP CONSTANT a = NULL",
		"DROP CONSTANT a",
		// Privilege restrictions
		"CREATE PRIVILEGE RESTRICTION ON a ON a a RESTRICT TO('user:x')",
		"ALTER PRIVILEGE RESTRICTION ON a ON a a ADD('a'), REMOVE()",
		"DROP PRIVILEGE RESTRICTION ON a ON a a",
		// Search / vector index maintenance
		"ALTER SEARCH INDEX a REBUILD, SET OPTIONS()",
		"ALTER SEARCH INDEX a ON t ADD COLUMN a, DROP COLUMN b, ALTER COLUMN c SET OPTIONS()",
		"ALTER VECTOR INDEX a REBUILD",
		// Row access policy alterations
		"ALTER ROW ACCESS POLICY a ON t GRANT TO('a'), FILTER USING(TRUE), REVOKE FROM('b'), REVOKE FROM ALL, RENAME TO b",
		"ALTER ALL ROW ACCESS POLICIES ON t REVOKE FROM ALL",
		// Untyped CREATE TABLE columns & COPY/CLONE time travel
		"CREATE TABLE a(a, b INT64,)",
		"CREATE TABLE a COPY b FOR SYSTEM_TIME AS OF NULL",
		"CREATE TABLE a CLONE b FOR SYSTEM_TIME AS OF NULL",
		// CREATE INDEX alias & adjacent string literal concatenation in GRANT
		"CREATE INDEX a ON a alias_a(NULL)",
		"GRANT ALL ON a TO 'a' 'a'",
	}
	for _, sql := range valid {
		if _, err := ParseWithDialect(sql, DialectBigQuery); err != nil {
			t.Errorf("ParseWithDialect(%q, DialectBigQuery) unexpected error: %v", sql, err)
		}
	}

	invalid := []string{
		// Braced constructors disabled in BigQuery
		"SELECT { }",
		"SELECT NEW a { }",
		"SELECT STRUCT { }",
		"SELECT STRUCT < > { }",
		// Spanner-only DDL rejected in BigQuery
		"CREATE LOCALITY GROUP a",
		"CREATE NULL_FILTERED INDEX a ON a(NULL)",
		"CREATE INDEX a ON a(NULL), INTERLEAVE IN a",
		"CREATE TABLE a PRIMARY KEY(a)",
		"CREATE TABLE a PRIMARY KEY(), INTERLEAVE IN PARENT a",
		"CREATE TABLE a ROW DELETION POLICY(NULL)",
		// BigQuery reserved keywords
		"SELECT MATCH_RECOGNIZE",
		"SELECT QUALIFY",
		"SELECT NULL QUALIFY",
		// Disabled-in-BigQuery / strict syntax rejections
		"SELECT NULL FOR UPDATE",
		"DESCRIBE a |> AGGREGATE",
		"|> INSERT REPLACE a",
		"|> INSERT UPDATE a",
		"|> EXPORT",
		"|> EXPORT AS",
		"TABLE a WHERE NULL",
		"CREATE INDEX a ON a()",
		"LOAD DATA INTO a() FROM FILES()",
		"ANALYZE OPTIONS(zzz)",
		"a : BEGIN END zzz",
		"FROM a WHERE",
		"FROM a FOR SYSTEM_TIME AS OF NULL zzz",
		"FROM (a a)",
		"FROM a JOIN a ON NULL ON NULL",
		"FROM UNNEST()",
		"DROP CONNECTION a RESTRICT",
		"DROP CONNECTION a CASCADE",
	}
	for _, sql := range invalid {
		if _, err := ParseWithDialect(sql, DialectBigQuery); err == nil {
			t.Errorf("ParseWithDialect(%q, DialectBigQuery) expected error, got nil", sql)
		}
	}
}
