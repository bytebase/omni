package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseCoreStatements(t *testing.T) {
	valid := []string{
		"SELECT 1",
		"SELECT /*+ MAPJOIN(a) */ a, b FROM t1 JOIN t2 ON t1.id = t2.id WHERE a > 0 GROUP BY a, b HAVING count(*) > 1 ORDER BY a DESC NULLS LAST LIMIT 10",
		"WITH cte AS (SELECT 1 AS x) SELECT * FROM cte",
		"FROM src INSERT OVERWRITE TABLE dst1 SELECT a, b INSERT INTO TABLE dst2 PARTITION (p='1') SELECT a, b",
		"SELECT a FROM t LATERAL VIEW OUTER explode(arr) e AS elem",
		"SELECT a FROM t TABLESAMPLE(BUCKET 1 OUT OF 4 ON rand())",
		"SELECT a, sum(b) OVER (PARTITION BY a ORDER BY c ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t",
		"SELECT CAST(a AS DECIMAL(10, 2)), CAST(b AS ARRAY<MAP<STRING, STRUCT<x:INT, y:UNIONTYPE<INT, STRING>>>>) FROM t",
		"SELECT a FROM t UNION ALL SELECT a FROM t2 INTERSECT DISTINCT SELECT a FROM t3 EXCEPT SELECT a FROM t4",
		"INSERT INTO TABLE t (a, b) VALUES (1, 'x'), (2, DEFAULT)",
		"UPDATE t SET a = 1, b = DEFAULT WHERE c = 2",
		"DELETE FROM t WHERE a = 1",
		"MERGE INTO t AS target USING s AS src ON target.id = src.id WHEN MATCHED AND src.del = 1 THEN DELETE WHEN MATCHED THEN UPDATE SET val = src.val WHEN NOT MATCHED THEN INSERT VALUES (src.id, src.val)",
		"LOAD DATA LOCAL INPATH '/tmp/data' OVERWRITE INTO TABLE t PARTITION (dt='2026-01-01')",
		"EXPORT TABLE t PARTITION (dt='2026-01-01') TO '/tmp/export'",
		"IMPORT EXTERNAL TABLE t FROM '/tmp/export' LOCATION '/tmp/loc'",
		"CREATE EXTERNAL TABLE IF NOT EXISTS db.t (id INT NOT NULL ENABLE NOVALIDATE RELY, name STRING COMMENT 'user', tags ARRAY<STRING>, PRIMARY KEY (id) DISABLE NOVALIDATE) PARTITIONED BY (dt STRING) CLUSTERED BY (id) SORTED BY (id ASC) INTO 8 BUCKETS ROW FORMAT DELIMITED FIELDS TERMINATED BY ',' STORED AS ORC LOCATION '/data/t' TBLPROPERTIES ('transactional'='true')",
		"CREATE DATABASE IF NOT EXISTS analytics COMMENT 'db' LOCATION '/wh/analytics' WITH DBPROPERTIES ('owner'='team')",
		"CREATE OR REPLACE VIEW v (c1 COMMENT 'col1') AS SELECT 1",
		"CREATE MATERIALIZED VIEW IF NOT EXISTS mv DISABLE REWRITE AS SELECT a, count(*) FROM t GROUP BY a",
		"CREATE TEMPORARY FUNCTION my_fn AS 'com.example.Fn' USING JAR 'hdfs:///fn.jar'",
		"CREATE TEMPORARY MACRO m(x INT) x + 1",
		"CREATE CONNECTOR c TYPE 'postgres' URL 'jdbc:postgresql://host/db' WITH DCPROPERTIES ('k'='v')",
		"CREATE SCHEDULED QUERY sq EVERY 1 DAY AS SELECT 1",
		"CREATE RESOURCE PLAN rp WITH QUERY_PARALLELISM = 4",
		"ALTER TABLE t ADD COLUMNS (c STRING) CASCADE",
		"ALTER TABLE t CREATE BRANCH b1 FOR SYSTEM_VERSION AS OF 123 RETAIN 7 DAYS WITH SNAPSHOT RETENTION 5 SNAPSHOTS 2 DAYS",
		"ALTER TABLE t EXECUTE EXPIRE_SNAPSHOTS RETAIN LAST 5",
		"ALTER TABLE t SET PARTITION SPEC (BUCKET(16, id), YEAR(ts))",
		"OPTIMIZE TABLE t REWRITE DATA",
		"DROP TABLE IF EXISTS t PURGE",
		"TRUNCATE TABLE t PARTITION (dt='2026-01-01')",
		"MSCK REPAIR TABLE t SYNC PARTITIONS",
		"ANALYZE TABLE t COMPUTE STATISTICS FOR COLUMNS a, b",
		"SHOW TABLES IN default LIKE 't*'",
		"SHOW CREATE TABLE t",
		"DESCRIBE FORMATTED t",
		"EXPLAIN EXTENDED SELECT 1",
		"EXPLAIN REWRITE SELECT 1",
		"USE analytics",
		"SET ROLE ALL",
		"GRANT SELECT ON TABLE t TO USER alice WITH GRANT OPTION",
		"REVOKE SELECT ON TABLE t FROM USER alice",
		"LOCK TABLE t EXCLUSIVE",
		"UNLOCK TABLE t",
		"START TRANSACTION READ ONLY, ISOLATION LEVEL SNAPSHOT",
		"COMMIT WORK",
		"ROLLBACK",
		"SET AUTOCOMMIT TRUE",
		"ABORT TRANSACTIONS 1 2 3",
		"KILL QUERY 'q1'",
		"PREPARE p FROM SELECT ? + 1",
		"EXECUTE p USING 42",
		// Lexical and operator coverage
		"SELECT 1L, 2S, 3Y, 4D, 5.5BD, 1e3, _utf8'abc', _utf8 0X01, 'a'\"b\", 'x' 'y'",
		"SELECT /* outer /* inner */ outer */ `col``name`, 1abc FROM t WHERE !a AND b == c AND d <=> e AND f != g AND h <> i AND j || k = 'z' AND m DIV 2 = 1",
		"SELECT INTERVAL '1-2' YEAR TO MONTH, INTERVAL '1 02:03:04' DAY TO SECOND, DATE '2026-01-01', TIMESTAMP '2026-01-01 00:00:00'",
		"SELECT TRIM(LEADING 'x' FROM a), EXTRACT(DAYOFWEEK FROM dt), FLOOR(ts TO HOUR) FROM t",
		"SELECT a.* FROM t",
		"SELECT 1 LIMIT 1",
		"DESCRIBE FORMATTED",
		"SHOW COMPACTIONS LIMIT",
		"GRANT ALL ON URI TO USER alice",
		"ALTER TABLE t EXECUTE FAST-FORWARD 'b1' 'b2'",
		"ALTER TABLE t EXECUTE CHERRY-PICK 42",
		"ALTER TABLE t EXECUTE DELETE ORPHAN-FILES OLDER THAN ('2026-01-01')",
		"DESCRIBE t col.$KEY$.$VALUE$.$ELEM$",
		// Nonreserved implicit table aliases
		"SELECT * FROM t limit WHERE 1 = 1",
		"SELECT * FROM t cluster WHERE 1 = 1",
		"SELECT * FROM t distribute WHERE 1 = 1",
		"SELECT * FROM t sort WHERE 1 = 1",
		// Nonreserved named window references inside parenthesized window specifications
		"SELECT sum(a) OVER (sort) FROM t WINDOW sort AS ()",
		"SELECT sum(a) OVER (cluster) FROM t WINDOW cluster AS ()",
		"SELECT sum(a) OVER (distribute) FROM t WINDOW distribute AS ()",
		"SELECT sum(a) OVER (sort SORT BY b) FROM t WINDOW sort AS ()",
		// VALUES row aliases vs DEFAULT
		"SELECT * FROM (VALUES (1 AS x, 2 y), (3, DEFAULT)) v",
	}
	for _, sql := range valid {
		res, err := Parse(sql)
		if err != nil {
			t.Errorf("Parse(%q) failed: %v", sql, err)
			continue
		}
		if res == nil || res.Statement == nil {
			t.Errorf("Parse(%q) returned nil statement", sql)
		}
	}

	invalid := []string{
		"",
		"   -- comment only",
		"SELECT FROM",
		"SELECT 1; SELECT 2",
		"SELECT 'unterminated",
		"SELECT /* unclosed",
		"CREATE TABLE a (`a.b` INT)",
		"EXTENDED COMMIT",
		"EXPLAIN zzz COMMIT",
		"SELECT 1 FROM a, b ON a.id = b.id",
		"MERGE INTO a USING b ON 1",
		"SELECT a(*, *)",
		"SELECT `` FROM t",
		"VALUES (DEFAULT, 1 AS x)",
		"VALUES (1 AS x, DEFAULT)",
	}
	for _, sql := range invalid {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) unexpectedly succeeded", sql)
		}
	}
}

func TestParseDeepNestingDoesNotStackOverflow(t *testing.T) {
	cases := []struct {
		name string
		sql  string
	}{
		{
			name: "nested parens expression",
			sql:  "SELECT " + strings.Repeat("(", 1000) + "1" + strings.Repeat(")", 1000),
		},
		{
			name: "nested parenthesized join source",
			sql:  "SELECT * FROM " + strings.Repeat("(", 1000) + "t" + strings.Repeat(")", 1000),
		},
		{
			name: "nested partitioned table function",
			sql:  "SELECT * FROM " + strings.Repeat("f(ON ", 500) + "t" + strings.Repeat(")", 500),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Parse(tc.sql); err == nil {
				t.Fatalf("expected depth-limit error on %s", tc.name)
			}
		})
	}
}

func TestParseNestedTrimAndParensLinearTime(t *testing.T) {
	// Nested ordinary TRIM(TRIM(...)), parentheses, aliased VALUES subqueries,
	// and malformed or valid nested IN ((SELECT ...)) must not exhibit
	// exponential backtracking.
	nestedTrim := "SELECT " + strings.Repeat("TRIM(", 40) + "a" + strings.Repeat(")", 40)
	if _, err := Parse(nestedTrim); err != nil {
		t.Fatalf("Parse nested TRIM failed: %v", err)
	}
	nestedTrimFrom := "SELECT " + strings.Repeat("TRIM(", 40) + "'x' FROM a" + strings.Repeat(")", 40)
	if _, err := Parse(nestedTrimFrom); err != nil {
		t.Fatalf("Parse nested TRIM FROM failed: %v", err)
	}
	nestedParens := "SELECT " + strings.Repeat("(", 40) + "1" + strings.Repeat(")", 40)
	if _, err := Parse(nestedParens); err != nil {
		t.Fatalf("Parse nested parens failed: %v", err)
	}
	nestedValuesAlias := "SELECT " + strings.Repeat("(VALUES (", 30) + "1" + strings.Repeat(" AS x))", 30)
	if _, err := Parse(nestedValuesAlias); err != nil {
		t.Fatalf("Parse nested aliased VALUES failed: %v", err)
	}
	nestedMalformedIn := "SELECT " + strings.Repeat("1 IN ((SELECT ", 30) + "!" + strings.Repeat("))", 30)
	if _, err := Parse(nestedMalformedIn); err == nil {
		t.Fatalf("expected error on malformed nested IN subquery")
	}
	nestedValidIn := "SELECT " + strings.Repeat("1 IN ((SELECT ", 30) + "1" + strings.Repeat(") + 1)", 30)
	if _, err := Parse(nestedValidIn); err != nil {
		t.Fatalf("Parse nested valid IN scalar subquery failed: %v", err)
	}
}

func TestParseDriverCorpora(t *testing.T) {
	// corpus.sql and dump/corpus.sql separate statements by blank lines;
	// dump/database.sql uses semicolons (split by hive/analysis.Split).
	blankLineFiles := []string{
		filepath.Join("..", "driver", "testdata", "corpus.sql"),
		filepath.Join("..", "driver", "testdata", "dump", "corpus.sql"),
	}
	for _, path := range blankLineFiles {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		for _, block := range strings.Split(string(raw), "\n\n") {
			trimmed := strings.TrimSpace(block)
			if trimmed == "" {
				continue
			}
			if _, err := Parse(trimmed); err != nil {
				t.Errorf("%s: Parse(%q) failed: %v", path, trimmed, err)
			}
		}
	}

	dbPath := filepath.Join("..", "driver", "testdata", "dump", "database.sql")
	raw, err := os.ReadFile(dbPath)
	if err != nil {
		t.Fatalf("read %s: %v", dbPath, err)
	}
	for _, s := range Split(string(raw)) {
		if s.Empty() {
			continue
		}
		trimmed := strings.TrimSpace(strings.TrimRight(s.Text, ";"))
		if trimmed == "" {
			continue
		}
		if _, err := Parse(trimmed); err != nil {
			t.Errorf("%s: Parse(%q) failed: %v", dbPath, trimmed, err)
		}
	}
}

func TestTokenizeAndSplit(t *testing.T) {
	toks, err := Tokenize("/* outer /* inner */ outer */ SELECT /*+ MAPJOIN(t) */ `a``b` FROM t")
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}
	wantTypes := []TokenType{KW_SELECT, TokQueryHint, TokIdentifier, KW_FROM, TokIdentifier, TokEOF}
	if len(toks) != len(wantTypes) {
		t.Fatalf("got %d tokens, want %d", len(toks), len(wantTypes))
	}
	for i, wt := range wantTypes {
		if toks[i].Type != wt {
			t.Errorf("token %d type = %v, want %v", i, toks[i].Type, wt)
		}
	}

	if _, err := Tokenize("SELECT /* unclosed"); err == nil {
		t.Fatalf("expected error on unclosed block comment")
	}

	segs := Split("SELECT 'a;b'; /* c; */ SELECT `x;y` FROM t;")
	if len(segs) != 2 {
		t.Fatalf("Split got %d segments, want 2: %+v", len(segs), segs)
	}
	if segs[0].Text != "SELECT 'a;b';" || segs[1].Text != " /* c; */ SELECT `x;y` FROM t;" {
		t.Fatalf("unexpected segments: %+v", segs)
	}
}

func TestParseStatementMetadata(t *testing.T) {
	res, err := Parse("WITH c AS (SELECT 1 AS x) INSERT INTO t SELECT * FROM c")
	if err != nil {
		t.Fatalf("Parse CTE INSERT failed: %v", err)
	}
	if res.Statement.Kind != StmtInsert {
		t.Errorf("CTE INSERT kind = %q, want %q", res.Statement.Kind, StmtInsert)
	}

	res, err = Parse("FROM src INSERT OVERWRITE TABLE dst SELECT a")
	if err != nil {
		t.Fatalf("Parse FROM INSERT failed: %v", err)
	}
	if res.Statement.Kind != StmtInsert {
		t.Errorf("FROM INSERT kind = %q, want %q", res.Statement.Kind, StmtInsert)
	}

	res, err = Parse("EXPLAIN FORMATTED ANALYZE SELECT 1")
	if err != nil {
		t.Fatalf("Parse EXPLAIN FORMATTED ANALYZE failed: %v", err)
	}
	if res.Statement.Kind != StmtExplain || res.Statement.Explain == nil ||
		!res.Statement.Explain.Formatted || !res.Statement.Explain.Analyze ||
		res.Statement.Explain.InnerKind != StmtQuery {
		t.Errorf("unexpected EXPLAIN metadata: %+v (explain=%+v)", res.Statement, res.Statement.Explain)
	}

	res, err = Parse("EXPLAIN ANALYZE TABLE t COMPUTE STATISTICS")
	if err != nil {
		t.Fatalf("Parse EXPLAIN ANALYZE TABLE failed: %v", err)
	}
	if res.Statement.Explain == nil || res.Statement.Explain.Analyze || res.Statement.Explain.InnerKind != StmtDDL {
		t.Errorf("EXPLAIN ANALYZE TABLE should not set Analyze option: %+v", res.Statement.Explain)
	}

	res, err = Parse("SELECT TRANSFORM(a) USING 'cat' FROM t")
	if err != nil {
		t.Fatalf("Parse SELECT TRANSFORM failed: %v", err)
	}
	if !res.Statement.HasTransform {
		t.Errorf("SELECT TRANSFORM should set HasTransform")
	}
}
