package parser

import (
	"path/filepath"
	"strings"
	"testing"
)

func assertParseErrorContains(t *testing.T, sql string, want string) {
	t.Helper()
	_, err := Parse(sql)
	if err == nil {
		t.Fatalf("Parse(%q) succeeded, want error containing %q", sql, want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("Parse(%q) error = %q, want substring %q", sql, err.Error(), want)
	}
}

func TestSoftFailTruncatedExpressions(t *testing.T) {
	cases := []string{
		"SELECT 1 +",
		"SELECT 1 -",
		"SELECT 1 *",
		"SELECT 1 /",
		"SELECT 1 AND",
		"SELECT 1 OR",
		"SELECT NOT",
		"SELECT -",
		"SELECT PRIOR",
		"SELECT CONNECT_BY_ROOT",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertParseErrorContains(t, sql, "syntax error at end of input")
		})
	}
}

func TestSoftFailTruncatedAdvancedExpressions(t *testing.T) {
	cases := []string{
		"SELECT CASE WHEN 1 = 1 THEN",
		"SELECT CAST(1 AS",
		"SELECT DECODE(1,",
		"SELECT JSON_VALUE(payload,",
		"SELECT JSON_TABLE(",
		"SELECT XMLTABLE(",
		"SELECT XMLSERIALIZE(CONTENT",
		"SELECT (1 +",
		"SELECT * FROM (SELECT",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertParseErrorContains(t, sql, "syntax error at end of input")
		})
	}
}

func TestSoftFailTruncatedPredicates(t *testing.T) {
	cases := []string{
		"SELECT 1 BETWEEN",
		"SELECT 1 BETWEEN 0 AND",
		"SELECT 1 IN (",
		"SELECT 1 IN (1,",
		"SELECT 'a' LIKE",
		"SELECT 'a' LIKE 'b' ESCAPE",
		"SELECT 1 IS",
		"SELECT 1 IS NOT",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertParseErrorContains(t, sql, "syntax error at end of input")
		})
	}
}

func TestSoftFailTruncatedClauses(t *testing.T) {
	cases := []string{
		"SELECT * FROM t WHERE",
		"SELECT * FROM t JOIN",
		"SELECT * FROM t JOIN t2 ON",
		"SELECT",
		"SELECT * FROM",
		"SELECT 1 GROUP BY",
		"SELECT 1 ORDER BY",
		"SELECT 1 UNION",
		"CREATE TABLE t (a NUMBER DEFAULT",
		"ALTER TABLE t ADD",
		"CREATE PROCEDURE p IS BEGIN",
		"DECLARE x NUMBER; BEGIN",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertParseErrorContains(t, sql, "syntax error at end of input")
		})
	}
}

func TestSoftFailDML(t *testing.T) {
	cases := []string{
		"INSERT INTO",
		"INSERT INTO t VALUES (",
		"UPDATE",
		"UPDATE t SET",
		"UPDATE t SET a =",
		"DELETE FROM",
		"DELETE FROM t WHERE",
		"MERGE INTO",
		"MERGE INTO t USING",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertParseErrorContains(t, sql, "syntax error at end of input")
		})
	}
}

func TestSoftFailDDL(t *testing.T) {
	cases := []string{
		"CREATE TABLE",
		"CREATE TABLE t (",
		"CREATE TABLE t (a",
		"CREATE INDEX",
		"CREATE VIEW v AS",
		"DROP",
		"DROP TABLE",
		"ALTER TABLE",
		"CREATE USER",
		"GRANT SELECT ON",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertParseErrorContains(t, sql, "syntax error at end of input")
		})
	}
}

func TestSoftFailPLSQL(t *testing.T) {
	cases := []string{
		"BEGIN",
		"BEGIN x :=",
		"BEGIN IF x THEN",
		"DECLARE v NUMBER :=",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertParseErrorContains(t, sql, "syntax error at end of input")
		})
	}
}

func TestStrictDuplicateClauses(t *testing.T) {
	cases := []struct {
		sql  string
		near string
	}{
		{"SELECT 1 WHERE 1 = 1 WHERE 2 = 2", `syntax error at or near "WHERE"`},
		{"SELECT 1 FROM t GROUP BY a GROUP BY b", `syntax error at or near "GROUP"`},
		{"SELECT a FROM t GROUP BY a HAVING COUNT(*) > 0 HAVING COUNT(*) > 1", `syntax error at or near "HAVING"`},
		{"SELECT 1 FROM t ORDER BY a ORDER BY b", `syntax error at or near "ORDER"`},
		{"SELECT 1 FROM t FETCH FIRST 1 ROWS ONLY FETCH FIRST 2 ROWS ONLY", `syntax error at or near "FETCH"`},
		{"CREATE TABLE t (a NUMBER) ORGANIZATION HEAP ORGANIZATION INDEX", `syntax error at or near "ORGANIZATION"`},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			assertParseErrorContains(t, tc.sql, tc.near)
		})
	}
}

func TestStrictParenthesisBalance(t *testing.T) {
	cases := []string{
		"SELECT (1",
		"SELECT * FROM (SELECT 1",
		"CREATE TABLE t (a NUMBER",
		"INSERT INTO t VALUES (1",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertParseErrorContains(t, sql, "syntax error at end of input")
		})
	}
}

func TestStrictStatementSeparators(t *testing.T) {
	cases := []struct {
		sql  string
		near string
	}{
		{"SELECT 1 SELECT 2", `syntax error at or near "SELECT"`},
		{"CREATE TABLE t (a NUMBER) CREATE TABLE u (b NUMBER)", `syntax error at or near "CREATE"`},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			assertParseErrorContains(t, tc.sql, tc.near)
		})
	}
}

func TestStrictUnknownOptions(t *testing.T) {
	cases := []struct {
		sql  string
		near string
	}{
		{"CREATE TABLE t (a NUMBER) FROBULATE", `syntax error at or near "FROBULATE"`},
		{"ALTER TABLE t FROBULATE", `syntax error at or near "FROBULATE"`},
		{"CREATE INDEX idx ON t(a) FROBULATE", `syntax error at or near "FROBULATE"`},
		{"DROP TABLE t FROBULATE", `syntax error at or near "FROBULATE"`},
		{"GRANT SELECT ON t TO u FROBULATE", `syntax error at or near "FROBULATE"`},
		{"ALTER SESSION FROBULATE", `syntax error at or near "FROBULATE"`},
		{"CREATE USER u FROBULATE", `syntax error at or near "FROBULATE"`},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			assertParseErrorContains(t, tc.sql, tc.near)
		})
	}
}

func TestStrictAlterUnknownTargetsDoNotSilentlySkip(t *testing.T) {
	cases := []struct {
		sql  string
		near string
	}{
		{"ALTER PUBLIC SELECT x", `syntax error at or near "SELECT"`},
		{"ALTER SHARED PUBLIC SELECT x", `syntax error at or near "SELECT"`},
		{"ALTER SHARED SELECT x", `syntax error at or near "SELECT"`},
		{"ALTER FROBULATE thing", `syntax error at or near "FROBULATE"`},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			assertParseErrorContains(t, tc.sql, tc.near)
		})
	}
}

func TestStrictCreateSchemaUnknownNestedCreateErrorsAtChild(t *testing.T) {
	assertParseErrorContains(t,
		"CREATE SCHEMA AUTHORIZATION app CREATE FROBULATE x",
		`syntax error at or near "FROBULATE"`,
	)
}

func TestStrictIllegalKeywordPosition(t *testing.T) {
	cases := []struct {
		sql  string
		near string
	}{
		{"SELECT FROM t", `syntax error at or near "FROM"`},
		{"SELECT * FROM WHERE", `syntax error at or near "WHERE"`},
		{"ALTER TABLE t ADD SELECT NUMBER", `syntax error at or near "SELECT"`},
		{"CREATE INDEX SELECT ON t(a)", `syntax error at or near "SELECT"`},
		{"CREATE INDEX sc.SELECT ON t(a)", `syntax error at or near "SELECT"`},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			assertParseErrorContains(t, tc.sql, tc.near)
		})
	}
}

func TestStrictReservedKeywordIdentifiers(t *testing.T) {
	cases := []struct {
		sql  string
		near string
	}{
		{"CREATE TABLE SELECT (a NUMBER)", `syntax error at or near "SELECT"`},
		{"CREATE TABLE sc.SELECT (a NUMBER)", `syntax error at or near "SELECT"`},
		{"CREATE TABLE t (SELECT NUMBER)", `syntax error at or near "SELECT"`},
		{"SELECT 1 FROM SELECT", `syntax error at or near "SELECT"`},
		{"SELECT 1 FROM sc.SELECT", `syntax error at or near "SELECT"`},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			assertParseErrorContains(t, tc.sql, tc.near)
		})
	}
}

func TestStrictOracleCompatibilityGaps(t *testing.T) {
	cases := []struct {
		sql  string
		near string
	}{
		{"SELECT 1 + FROM dual", `syntax error at or near "FROM"`},
		{"CREATE TABLE (a NUMBER)", `syntax error at or near "("`},
		{"CREATE PACKAGE compat_bad_pkg IS PROCEDURE p", "syntax error at end of input"},
		{"CREATE TYPE compat_bad_obj AS OBJECT", "syntax error at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			assertParseErrorContains(t, tc.sql, tc.near)
		})
	}
}

func TestStrictOracleCompatibilityPhaseAGaps(t *testing.T) {
	cases := []struct {
		sql  string
		near string
	}{
		{"SELECT * FROM employees ORDER BY employee_id NULLS", `syntax error at end of input`},
		{"SELECT * FROM employees OFFSET ROWS", `syntax error at or near "ROWS"`},
		{"INSERT INTO t (a,) VALUES (1)", `syntax error at or near ")"`},
		{"INSERT ALL INTO t(a) VALUES(1)", `syntax error at end of input`},
		{"UPDATE t SET a = , b = 2", `syntax error at or near ","`},
		{"UPDATE t SET (a, b) = (SELECT 1 FROM dual)", `syntax error`},
		{"INSERT FIRST WHEN THEN INTO t(a) VALUES(1) SELECT 1 FROM dual", `syntax error at or near "THEN"`},
		{"MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET t.x = s.x", `syntax error at or near "t"`},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			assertParseErrorContains(t, tc.sql, tc.near)
		})
	}
}

func TestOracleCompatibilityPhaseAAccepts(t *testing.T) {
	cases := []string{
		"SELECT * FROM employees WHERE salary > ANY (SELECT salary FROM employees)",
		"SELECT * FROM employees WHERE salary >= ALL (SELECT salary FROM employees)",
		"SELECT XMLCAST(XMLQUERY('/r/id' PASSING XMLTYPE('<r><id>1</id></r>') RETURNING CONTENT) AS NUMBER) FROM dual",
		"UPDATE (SELECT a FROM t) SET a = 1",
		"DELETE FROM (SELECT * FROM t) WHERE id = 1",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestStrictOracleCompatibilityPhaseBGaps(t *testing.T) {
	cases := []string{
		"CREATE TABLE compat_bad_cols (a NUMBER,)",
		"CREATE TABLE compat_bad_col (a)",
		"CREATE TABLE compat_bad_ctas AS",
		"CREATE TABLE compat_bad_part (a NUMBER) PARTITION BY RANGE (a) (PARTITION p1)",
		"CREATE INDEX compat_bad_idx ON t()",
		"CREATE INDEX compat_bad_idx ON t(a,)",
		"CREATE INDEX compat_bad_idx ON t",
		"CREATE VIEW compat_bad_v () AS SELECT 1 FROM dual",
		"CREATE SEQUENCE compat_bad_seq CACHE",
		"CREATE SEQUENCE compat_bad_seq INCREMENT BY",
		"DROP TABLE t CASCADE PURGE",
		"DROP INDEX ONLINE",
		"ALTER TABLE t RENAME COLUMN a",
		"COMMENT ON TABLE t IS NULL",
		"TRUNCATE TABLE t DROP",
		"CREATE MATERIALIZED VIEW LOG ON",
		"ALTER MATERIALIZED VIEW compat_bad_mv",
		"ALTER MATERIALIZED VIEW compat_bad_mv REFRESH",
		"DROP MATERIALIZED VIEW LOG compat_bad_mlog",
		"CREATE CLUSTER compat_bad_cluster",
		"CREATE CLUSTER compat_bad_cluster ()",
		"ALTER CLUSTER compat_bad_cluster",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertParseErrorContains(t, sql, "syntax error")
		})
	}
}

func TestOracleCompatibilityPhaseBAccepts(t *testing.T) {
	ParseAndCheck(t, "CREATE TABLE compat_part_hash (a NUMBER) PARTITION BY HASH (a) PARTITIONS 4")
}

func TestStrictV2CoverageMatrix(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "strictness_v2.tsv"))
	var scenarios int
	familyCounts := make(map[string]int)
	for _, row := range rows {
		sql := row.Fields["sql"]
		if sql == "" {
			t.Fatalf("%s: strictness row has empty sql", row.Key)
		}
		family := row.Fields["family"]
		familyCounts[family]++
		scenarios++
		t.Run(row.Key, func(t *testing.T) {
			switch row.Fields["expect"] {
			case "error":
				want := row.Fields["contains"]
				if want == "" {
					want = "syntax error"
				}
				assertParseErrorContains(t, sql, want)
			case "accept":
				if _, err := Parse(sql); err != nil {
					t.Fatalf("Parse(%q) error = %v, want success", sql, err)
				}
			default:
				t.Fatalf("%s: unknown expectation %q", row.Key, row.Fields["expect"])
			}
		})
	}
	if scenarios < 100 {
		t.Fatalf("strictness v2 scenarios = %d, want at least 100", scenarios)
	}
	t.Logf("Oracle strictness v2 scenarios=%d family=%v", scenarios, familyCounts)
}
