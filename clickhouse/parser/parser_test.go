package parser

import (
	"errors"
	"strings"
	"testing"
)

func TestParseValidStatements(t *testing.T) {
	cases := []struct {
		name string
		sql  string
	}{
		{"simple_select", "SELECT 1, 'hello', true, false, NULL"},
		{"select_with_semicolon", "SELECT 1;"},
		{"select_comments", "-- line comment\n# hash comment\n/* outer /* nested */ comment */ SELECT 1"},
		{"select_dollar_string", "SELECT $$raw 'string'$$, $tag$hello$tag$"},
		{"select_quoted_idents", "SELECT `col 1`, \"col 2\", `a``b`, \"c\"\"d\" FROM `my db`.`my tbl`"},
		{"select_hex_bin_literals", "SELECT 0x1A, 0b1010, 0x1.8p3, X'DEADBEEF', B'0101', 1.5e-2, 1_000_000"},
		{"select_query_param", "SELECT {id:UInt64}, {name:Nullable(String)}"},
		{"select_cte_and_union", "WITH cte AS (SELECT 1 AS a), 2 AS b SELECT a, b FROM cte UNION ALL SELECT 3, 4"},
		{"select_recursive_cte", "WITH RECURSIVE t AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM t WHERE n < 10) SELECT * FROM t"},
		{"select_from_first", "FROM events WHERE user_id = 42"},
		{"select_joins", "SELECT * FROM a GLOBAL ANY LEFT JOIN b ON a.id = b.id ALL INNER JOIN c USING (id) CROSS JOIN d FINAL SAMPLE 1/10 OFFSET 1/2"},
		{"select_array_join", "SELECT a, x FROM t LEFT ARRAY JOIN arr AS x PREWHERE a > 0 WHERE x != '' GROUP BY a, x WITH ROLLUP WITH TOTALS HAVING count() > 1 QUALIFY row_number() OVER w = 1 WINDOW w AS (PARTITION BY a ORDER BY x ROWS BETWEEN 1 PRECEDING AND CURRENT ROW) ORDER BY a DESC NULLS LAST WITH FILL FROM 1 TO 10 STEP 1 INTERPOLATE (x AS 0) LIMIT 1, 5 BY a LIMIT 10 OFFSET 2 WITH TIES SETTINGS max_threads = 4 FORMAT JSONEachRow"},
		{"select_distinct_on", "SELECT DISTINCT ON (a, b) a, b, c FROM t"},
		{"select_expressions", "SELECT CASE WHEN x > 0 THEN 1 ELSE 0 END, CAST(x AS Nullable(Decimal(18, 4))), x::LowCardinality(String), DATE '2026-01-01', TIMESTAMP '2026-01-01 00:00:00', EXTRACT(DAY FROM ts), INTERVAL 5 MINUTE, SUBSTRING(s FROM 1 FOR 3), TRIM(BOTH ' ' FROM s), x BETWEEN 1 AND 10, x GLOBAL NOT IN (1, 2), s NOT ILIKE '%foo%', x IS NOT NULL, a ? b : c, arr[1], tup.1.2, [1, 2, 3], (1, 'a'), quantile(0.95)(x), arrayMap(x -> x + 1, arr), arrayMap((x, y) -> x + y, a1, a2), * EXCEPT (secret)"},
		{"select_into_outfile", "SELECT 1 INTO OUTFILE '/tmp/out.csv' COMPRESSION 'gzip' FORMAT CSV"},
		{"insert_values", "INSERT INTO db.t (a, b) VALUES (1, 'x'), (2, 'y')"},
		{"insert_format", "INSERT INTO TABLE t FORMAT JSONEachRow"},
		{"insert_select", "INSERT INTO FUNCTION s3('https://example.com/a.parquet') SELECT * FROM t"},
		{"delete_stmt", "DELETE FROM db.t ON CLUSTER c IN PARTITION '202601' WHERE id = 1"},
		{"update_stmt", "UPDATE t SET a = 1, b = 'x' IN PARTITION '202601' WHERE id = 1"},
		{"create_database", "CREATE DATABASE IF NOT EXISTS analytics ON CLUSTER prod ENGINE = Atomic COMMENT 'main'"},
		{"create_table", "CREATE TABLE IF NOT EXISTS db.events UUID '550e8400-e29b-41d4-a716-446655440000' ON CLUSTER c (id UInt64, ts DateTime64(3, 'UTC') CODEC(Delta, ZSTD), tags Array(LowCardinality(String)), meta Map(String, Nullable(Int64)), status Enum8('ok' = 1, 'err' = 2), items Nested(k String, v UInt32), computed UInt64 MATERIALIZED id + 1 TTL ts + INTERVAL 30 DAY, CONSTRAINT chk CHECK id > 0, INDEX idx_ts ts TYPE minmax GRANULARITY 4, PROJECTION p (SELECT id GROUP BY id ORDER BY id)) ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/events', '{replica}') PARTITION BY toYYYYMM(ts) PRIMARY KEY (id) ORDER BY (id, ts) SAMPLE BY id TTL ts + INTERVAL 90 DAY DELETE SETTINGS index_granularity = 8192"},
		{"create_table_as_select", "CREATE OR REPLACE TEMPORARY TABLE t AS SELECT 1 AS a"},
		{"create_view", "CREATE OR REPLACE VIEW db.v (a UInt64) AS SELECT 1"},
		{"create_mat_view_to", "CREATE MATERIALIZED VIEW mv TO dst AS SELECT * FROM src"},
		{"create_mat_view_engine", "CREATE MATERIALIZED VIEW mv ENGINE = MergeTree ORDER BY a POPULATE AS SELECT 1 AS a"},
		{"create_dictionary", "CREATE DICTIONARY d (id UInt64, val String DEFAULT '' HIERARCHICAL) PRIMARY KEY id SOURCE(CLICKHOUSE(HOST 'localhost' PORT 9000)) LIFETIME(MIN 0 MAX 300) LAYOUT(FLAT())"},
		{"create_function", "CREATE OR REPLACE FUNCTION add_one AS x -> x + 1"},
		{"alter_table", "ALTER TABLE db.t ON CLUSTER c ADD COLUMN IF NOT EXISTS c String AFTER b, MODIFY COLUMN c LowCardinality(String), RENAME COLUMN c TO d, DROP INDEX IF EXISTS idx, DELETE WHERE a = 0, UPDATE a = 1 WHERE a = 2, MODIFY TTL ts + INTERVAL 7 DAY, MOVE PARTITION '202601' TO DISK 'cold'"},
		{"drop_detach_truncate_rename_exchange", "DROP TABLE IF EXISTS db.t ON CLUSTER c NO DELAY"},
		{"detach_table", "DETACH TABLE db.t PERMANENTLY"},
		{"truncate_table", "TRUNCATE TEMPORARY TABLE IF EXISTS t"},
		{"rename_table", "RENAME TABLE a TO b, c TO d ON CLUSTER prod"},
		{"exchange_tables", "EXCHANGE TABLES db.a AND db.b ON CLUSTER prod"},
		{"undrop_table", "UNDROP TABLE db.t UUID '550e8400-e29b-41d4-a716-446655440000'"},
		{"attach_dictionary", "ATTACH DICTIONARY db.d ON CLUSTER prod"},
		{"show_statements", "SHOW FULL TABLES FROM db LIKE '%evt%' LIMIT 10"},
		{"show_create_table", "SHOW CREATE TABLE db.t"},
		{"show_columns", "SHOW EXTENDED FULL COLUMNS FROM t IN db WHERE type = 'String'"},
		{"show_grants", "SHOW GRANTS FOR alice, bob WITH IMPLICIT FINAL"},
		{"show_setting", "SHOW SETTING max_threads"},
		{"describe_stmt", "DESCRIBE TABLE db.t"},
		{"explain_stmt", "EXPLAIN PLAN header = 1 SELECT * FROM t"},
		{"explain_ast", "EXPLAIN AST CREATE TABLE t (a UInt64) ENGINE = Memory"},
		{"exists_stmt", "EXISTS DICTIONARY db.d"},
		{"check_stmt", "CHECK TABLE db.t PARTITION ID '202601'"},
		{"optimize_stmt", "OPTIMIZE TABLE db.t ON CLUSTER c PARTITION '202601' FINAL DEDUPLICATE"},
		{"kill_stmt", "KILL MUTATION ON CLUSTER c WHERE mutation_id = '1' SYNC"},
		{"system_stmt", "SYSTEM FLUSH LOGS query_log"},
		{"system_stop_merges", "SYSTEM STOP MERGES"},
		{"set_stmt", "SET max_threads = 8, send_logs_level = 'trace'"},
		{"use_stmt", "USE analytics"},
		{"watch_stmt", "WATCH db.lv EVENTS LIMIT 10"},
		{"grant_stmt", "GRANT SELECT ON db.* TO analyst"},
		{"revoke_stmt", "REVOKE SELECT ON db.* FROM analyst"},
		{"backup_stmt", "BACKUP TABLE db.t TO Disk('backups', '1.zip')"},
		{"restore_stmt", "RESTORE TABLE db.t FROM Disk('backups', '1.zip')"},
		{"create_user", "CREATE USER IF NOT EXISTS alice IDENTIFIED WITH sha256_password BY 'secret'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tc.sql, err)
			}
			if len(stmts) != 1 {
				t.Fatalf("Parse(%q) returned %d statements, want 1", tc.sql, len(stmts))
			}
		})
	}
}

func TestParseErrorsAndPositions(t *testing.T) {
	cases := []struct {
		name    string
		sql     string
		wantPos int
	}{
		{"empty", "", 0},
		{"whitespace_only", "   \n\t ", 6},
		{"comment_only", "-- just a comment", 17},
		{"unterminated_string", "SELECT 'unterminated", 7},
		{"unterminated_block_comment", "SELECT /* nested /* comment */ 1", 7},
		{"unterminated_dollar_quote", "SELECT $tag$abc", 7},
		{"multi_statement", "SELECT 1; SELECT 2", 10},
		{"trailing_garbage", "SELECT 1 +", 10},
		{"leading_dot_number", "SELECT .5", 7},
		{"glued_ident_on_number", "SELECT 1.zzz", 7},
		{"invalid_uuid", "CREATE TABLE t UUID 'not-a-uuid' (a UInt64) ENGINE = Memory", 20},
		{"with_ties_no_order_by", "SELECT 1 LIMIT 5 WITH TIES", 26},
		{"unbounded_following_start", "SELECT count() OVER (ROWS UNBOUNDED FOLLOWING)", 26},
		{"array_join_without_from", "SELECT * ARRAY JOIN a", 9},
		{"interpolate_without_order_by", "SELECT * INTERPOLATE (a)", 9},
		{"bare_show_ident", "SHOW foo", 5},
		{"parenthesized_single_table", "SELECT * FROM (a)", 16},
		{"invalid_char", "SELECT @foo", 7},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("Parse(%q) succeeded, want *ParseError", tc.sql)
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) error type = %T, want *ParseError", tc.sql, err)
			}
			if pe.Position != tc.wantPos {
				t.Errorf("Parse(%q) Position = %d, want %d (err: %v)", tc.sql, pe.Position, tc.wantPos, pe)
			}
		})
	}
}

func TestParseDeepNestingDoesNotOverflow(t *testing.T) {
	deep := "SELECT " + strings.Repeat("(", 400) + "1" + strings.Repeat(")", 400)
	_, err := Parse(deep)
	if err == nil {
		t.Fatalf("expected nesting depth error on 400 nested parens")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("got %T, want *ParseError", err)
	}
}
