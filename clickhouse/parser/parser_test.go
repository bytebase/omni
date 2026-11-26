package parser

import (
	"bufio"
	"errors"
	"os"
	"reflect"
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
		{"select_view_table_function", "SELECT a FROM view(SELECT a FROM t) AS v, view((SELECT b FROM u))"},
		{"select_exists_subquery", "SELECT a FROM t WHERE EXISTS (SELECT 1 FROM u) AND NOT EXISTS ((SELECT 2))"},
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
		{"alter_table", "ALTER TABLE db.t ON CLUSTER c ADD COLUMN IF NOT EXISTS c String AFTER b, MODIFY COLUMN c LowCardinality(String), RENAME COLUMN c TO d, DROP INDEX IF EXISTS idx, DELETE WHERE a = 0, UPDATE a = 1 WHERE a = 2, MOVE PARTITION '202601' TO DISK 'cold', MODIFY TTL ts + INTERVAL 7 DAY"},
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
		{"grant_role", "GRANT analyst TO alice WITH ADMIN OPTION"},
		{"revoke_stmt", "REVOKE SELECT ON db.* FROM analyst"},
		{"backup_stmt", "BACKUP TABLE db.t TO Disk('backups', '1.zip')"},
		{"restore_stmt", "RESTORE TABLE db.t FROM Disk('backups', '1.zip')"},
		{"create_user", "CREATE USER IF NOT EXISTS alice IDENTIFIED WITH sha256_password BY 'secret'"},
		{"create_user_host_any", "CREATE USER alice HOST ANY SETTINGS max_threads = 1"},
		{"drop_user", "DROP USER IF EXISTS alice, bob ON CLUSTER prod"},
		{"insert_csv_inline_non_sql", "INSERT INTO t FORMAT CSV\na@b,1\n\"unclosed"},
		{"ttl_group_by_rollup", "ALTER TABLE t MODIFY TTL ts + INTERVAL 1 DAY GROUP BY k SET v = sum(v), m = max(m)"},
		{"window_partition_rows_no_order", "SELECT sum(v) OVER (PARTITION BY k ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t"},
		{"window_partition_range_no_order", "SELECT sum(v) OVER (PARTITION BY k RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t"},
		{"substring_comma_form", "SELECT substring('abcdef', 2, 3)"},
		{"case_as_column_ident", "SELECT case, case AS c FROM t WHERE case = 1"},
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
		{"truncated_grant", "GRANT SELECT", 12},
		{"truncated_revoke", "REVOKE SELECT ON *.*", 20},
		{"truncated_backup", "BACKUP TABLE", 12},
		{"backup_without_to", "BACKUP TABLE t", 14},
		{"restore_without_from", "RESTORE TABLE t", 15},
		{"create_user_garbage", "CREATE USER alice GARBAGE TOKENS", 18},
		{"create_user_host_any_garbage", "CREATE USER alice HOST ANY GARBAGE", 27},
		{"drop_user_garbage", "DROP USER alice GARBAGE", 16},
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
	cases := []struct {
		name string
		sql  string
	}{
		{"parens", "SELECT " + strings.Repeat("(", 400) + "1" + strings.Repeat(")", 400)},
		{"table_functions", "SELECT * FROM " + strings.Repeat("f(", 400) + "1" + strings.Repeat(")", 400)},
		{"parenthesized_joins", "SELECT * FROM " + strings.Repeat("(", 400) + "a JOIN b ON 1" + strings.Repeat(")", 400)},
		{"unary_minus", "SELECT " + strings.Repeat("-", 400) + "1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("expected nesting depth error on deeply nested %s", tc.name)
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("got %T, want *ParseError", err)
			}
		})
	}
}

func TestParseEscapeSequences(t *testing.T) {
	stmts, err := Parse("SELECT 'line1\\nline2\\t\\x41\\u0042\\\\\\'\"' FROM `line\\nbreak`")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	sel := stmts[0].(*SelectUnionStmt).Items[0]
	lit := sel.Columns[0].(*LiteralExpr)
	wantVal := "line1\nline2\tAB\\'\""
	if lit.Val != wantVal {
		t.Errorf("string Val = %q, want %q", lit.Val, wantVal)
	}
	tbl := sel.From.Table.Table
	if tbl != "line\nbreak" {
		t.Errorf("table ident = %q, want %q", tbl, "line\nbreak")
	}
}

func TestParseNestedDataTypeAST(t *testing.T) {
	stmts, err := Parse("SELECT CAST(x AS Array(Nullable(UInt64)))")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	cast := stmts[0].(*SelectUnionStmt).Items[0].Columns[0].(*CastExpr)
	arr := cast.Type
	if arr.Name != "Array" || len(arr.TypeArgs) != 1 {
		t.Fatalf("outer type = %+v, want Array with 1 TypeArg", arr)
	}
	nullable := arr.TypeArgs[0]
	if nullable.Name != "Nullable" || len(nullable.TypeArgs) != 1 {
		t.Fatalf("middle type = %+v, want Nullable with 1 TypeArg", nullable)
	}
	u64 := nullable.TypeArgs[0]
	if u64.Name != "UInt64" {
		t.Fatalf("inner type = %q, want UInt64", u64.Name)
	}
}

func TestParseLimitCommaVsOffset(t *testing.T) {
	stmts1, err := Parse("SELECT 1 LIMIT 0, 10")
	if err != nil {
		t.Fatalf("Parse comma LIMIT failed: %v", err)
	}
	lim1 := stmts1[0].(*SelectUnionStmt).Items[0].Limit
	if lim1.Limit.(*LiteralExpr).Raw != "10" || lim1.Offset.(*LiteralExpr).Raw != "0" {
		t.Errorf("LIMIT 0, 10 => Limit=%v Offset=%v, want Limit=10 Offset=0", lim1.Limit, lim1.Offset)
	}

	stmts2, err := Parse("SELECT 1 LIMIT 10 OFFSET 0")
	if err != nil {
		t.Fatalf("Parse OFFSET LIMIT failed: %v", err)
	}
	lim2 := stmts2[0].(*SelectUnionStmt).Items[0].Limit
	if lim2.Limit.(*LiteralExpr).Raw != "10" || lim2.Offset.(*LiteralExpr).Raw != "0" {
		t.Errorf("LIMIT 10 OFFSET 0 => Limit=%v Offset=%v, want Limit=10 Offset=0", lim2.Limit, lim2.Offset)
	}
}

func TestParseInsertFormatRawData(t *testing.T) {
	stmts, err := Parse("INSERT INTO t FORMAT CSV\na@b,1\n\"unclosed")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	ins := stmts[0].(*InsertStmt)
	if ins.Format != "CSV" {
		t.Errorf("Format = %q, want CSV", ins.Format)
	}
	if ins.Data != "\na@b,1\n\"unclosed" {
		t.Errorf("Data = %q, want %q", ins.Data, "\na@b,1\n\"unclosed")
	}
}

func TestParseTreeHeightLimit(t *testing.T) {
	sql := "SELECT " + strings.Repeat("1+", 100001) + "1"
	_, err := Parse(sql)
	if err == nil {
		t.Fatalf("expected *ParseError when AST height exceeds maxTreeHeight")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("got %T, want *ParseError", err)
	}
}

func TestParseWideFunctionArgsLinearTime(t *testing.T) {
	args := make([]string, 16000)
	for i := range args {
		args[i] = "a"
	}
	sql := "SELECT f(" + strings.Join(args, ",") + ")"
	stmts, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse failed on 16,000 args: %v", err)
	}
	fc := stmts[0].(*SelectUnionStmt).Items[0].Columns[0].(*FuncCallExpr)
	if len(fc.Args) != 16000 {
		t.Fatalf("got %d args, want 16000", len(fc.Args))
	}
}

func TestParseNestedSubstringLinearTime(t *testing.T) {
	sql := "SELECT " + strings.Repeat("substring(", 40) + "'abc'" + strings.Repeat(", 1)", 40)
	stmts, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse failed on 40 nested substrings: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("got %d statements, want 1", len(stmts))
	}
}

func TestParseColumnDefaultKinds(t *testing.T) {
	sql := "CREATE TABLE t (a UInt64 DEFAULT 1, b UInt64 MATERIALIZED 2, c UInt64 ALIAS 3) ENGINE = Memory"
	stmts, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	ct := stmts[0].(*CreateTableStmt)
	want := []string{"DEFAULT", "MATERIALIZED", "ALIAS"}
	if len(ct.Columns) != len(want) {
		t.Fatalf("got %d columns, want %d", len(ct.Columns), len(want))
	}
	for i, w := range want {
		if ct.Columns[i].DefaultKind != w {
			t.Errorf("col %d DefaultKind = %q, want %q", i, ct.Columns[i].DefaultKind, w)
		}
	}
}

func TestParseOutputClauseSpan(t *testing.T) {
	sql1 := "SELECT 1 FORMAT JSON"
	stmts1, err := Parse(sql1)
	if err != nil {
		t.Fatalf("Parse(%q) failed: %v", sql1, err)
	}
	if stmts1[0].StartPos() != 0 || stmts1[0].EndPos() != len(sql1) {
		t.Errorf("span = [%d, %d), want [0, %d)", stmts1[0].StartPos(), stmts1[0].EndPos(), len(sql1))
	}

	sql2 := "SELECT 1 INTO OUTFILE '/tmp/x' FORMAT CSV;"
	stmts2, err := Parse(sql2)
	if err != nil {
		t.Fatalf("Parse(%q) failed: %v", sql2, err)
	}
	if stmts2[0].StartPos() != 0 || stmts2[0].EndPos() != len(sql2)-1 {
		t.Errorf("span = [%d, %d), want [0, %d)", stmts2[0].StartPos(), stmts2[0].EndPos(), len(sql2)-1)
	}
}

func TestParseQualifiedStarSpan(t *testing.T) {
	sql1 := "SELECT db.t.* FROM db.t"
	stmts1, err := Parse(sql1)
	if err != nil {
		t.Fatalf("Parse(%q) failed: %v", sql1, err)
	}
	star1 := stmts1[0].(*SelectUnionStmt).Items[0].Columns[0].(*StarExpr)
	if got := sql1[star1.StartPos():star1.EndPos()]; got != "db.t.*" {
		t.Errorf("projection star span text = %q, want %q", got, "db.t.*")
	}

	sql2 := "SELECT f(db.t.*) FROM db.t"
	stmts2, err := Parse(sql2)
	if err != nil {
		t.Fatalf("Parse(%q) failed: %v", sql2, err)
	}
	star2 := stmts2[0].(*SelectUnionStmt).Items[0].Columns[0].(*FuncCallExpr).Args[0].(*StarExpr)
	if got := sql2[star2.StartPos():star2.EndPos()]; got != "db.t.*" {
		t.Errorf("func-arg star span text = %q, want %q", got, "db.t.*")
	}
}

func TestParseAlterColumnElementAndWindowBaseName(t *testing.T) {
	stmts, err := Parse("ALTER TABLE t ADD COLUMN c UInt64 DEFAULT 1 AFTER b, MODIFY COLUMN c LowCardinality(String)")
	if err != nil {
		t.Fatalf("Parse ALTER failed: %v", err)
	}
	at := stmts[0].(*AlterTableStmt)
	if len(at.Clauses) != 2 {
		t.Fatalf("got %d clauses, want 2", len(at.Clauses))
	}
	c0 := at.Clauses[0]
	if c0.Element == nil || c0.Element.Type == nil || c0.Element.Type.Name != "UInt64" || c0.Element.DefaultKind != "DEFAULT" || c0.After != "b" {
		t.Errorf("ADD COLUMN clause = %+v, element = %+v", c0, c0.Element)
	}
	c1 := at.Clauses[1]
	if c1.Element == nil || c1.Element.Type == nil || c1.Element.Type.Name != "LowCardinality" {
		t.Errorf("MODIFY COLUMN clause = %+v, element = %+v", c1, c1.Element)
	}

	wstmts, err := Parse("SELECT sum(v) OVER (w ORDER BY x) FROM t WINDOW w AS (PARTITION BY k)")
	if err != nil {
		t.Fatalf("Parse OVER (w ...) failed: %v", err)
	}
	fc := wstmts[0].(*SelectUnionStmt).Items[0].Columns[0].(*FuncCallExpr)
	if fc.Over == nil || fc.Over.BaseWindow != "w" || len(fc.Over.OrderBy) != 1 {
		t.Errorf("OVER spec = %+v, want BaseWindow=\"w\" with 1 OrderBy", fc.Over)
	}
}

func TestParseAlterMutationsAndProjectionAST(t *testing.T) {
	stmts, err := Parse("ALTER TABLE t UPDATE x = 7 WHERE id = 9, DELETE WHERE id = 0")
	if err != nil {
		t.Fatalf("Parse ALTER mutations failed: %v", err)
	}
	at := stmts[0].(*AlterTableStmt)
	if len(at.Clauses) != 2 {
		t.Fatalf("got %d clauses, want 2", len(at.Clauses))
	}
	upd := at.Clauses[0]
	if len(upd.Assignments) != 1 || upd.Assignments[0].Column != "x" || upd.Where == nil {
		t.Errorf("UPDATE clause = %+v, want 1 assignment for x and non-nil Where", upd)
	}
	del := at.Clauses[1]
	if del.Where == nil {
		t.Errorf("DELETE clause = %+v, want non-nil Where", del)
	}

	pstmts, err := Parse("CREATE TABLE t (x UInt64, PROJECTION p (SELECT x GROUP BY x ORDER BY x)) ENGINE = Memory")
	if err != nil {
		t.Fatalf("Parse CREATE TABLE with PROJECTION failed: %v", err)
	}
	ct := pstmts[0].(*CreateTableStmt)
	projElem := ct.Columns[1]
	if projElem.Projection == nil || len(projElem.Projection.Columns) != 1 || projElem.Projection.GroupBy == nil || len(projElem.Projection.OrderBy) != 1 {
		t.Errorf("PROJECTION element = %+v, proj = %+v", projElem, projElem.Projection)
	}
}

func TestParseParenSpansAndCodecAST(t *testing.T) {
	sql1 := "SELECT (1 + 2) * 3"
	stmts1, err := Parse(sql1)
	if err != nil {
		t.Fatalf("Parse(%q) failed: %v", sql1, err)
	}
	mul1 := stmts1[0].(*SelectUnionStmt).Items[0].Columns[0].(*BinaryExpr)
	if got := sql1[mul1.StartPos():mul1.EndPos()]; got != "(1 + 2) * 3" {
		t.Errorf("mul1 span text = %q, want %q", got, "(1 + 2) * 3")
	}

	sql2 := "SELECT 1 * (2 + 3)"
	stmts2, err := Parse(sql2)
	if err != nil {
		t.Fatalf("Parse(%q) failed: %v", sql2, err)
	}
	mul2 := stmts2[0].(*SelectUnionStmt).Items[0].Columns[0].(*BinaryExpr)
	if got := sql2[mul2.StartPos():mul2.EndPos()]; got != "1 * (2 + 3)" {
		t.Errorf("mul2 span text = %q, want %q", got, "1 * (2 + 3)")
	}

	cstmts, err := Parse("CREATE TABLE t (c UInt64 CODEC(Delta, ZSTD(1))) ENGINE = Memory")
	if err != nil {
		t.Fatalf("Parse CREATE TABLE CODEC failed: %v", err)
	}
	codec := cstmts[0].(*CreateTableStmt).Columns[0].Codec
	if codec == nil || len(codec.Codecs) != 2 || codec.Codecs[0].Name != "Delta" || codec.Codecs[1].Name != "ZSTD" || len(codec.Codecs[1].Args) != 1 {
		t.Errorf("column Codec = %+v, want Delta and ZSTD(1)", codec)
	}

	astmts, err := Parse("ALTER TABLE t MODIFY COLUMN c CODEC(ZSTD(3))")
	if err != nil {
		t.Fatalf("Parse ALTER MODIFY COLUMN CODEC failed: %v", err)
	}
	acodec := astmts[0].(*AlterTableStmt).Clauses[0].Codec
	if acodec == nil || len(acodec.Codecs) != 1 || acodec.Codecs[0].Name != "ZSTD" || len(acodec.Codecs[0].Args) != 1 {
		t.Errorf("alter Codec = %+v, want ZSTD(3)", acodec)
	}
}

func TestParseTTLASTAndCommaDisambiguation(t *testing.T) {
	// Structured TTL actions on ALTER TABLE: TO DISK vs TO VOLUME vs DELETE WHERE.
	s1, err := Parse("ALTER TABLE t MODIFY TTL ts TO DISK 'hot'")
	if err != nil {
		t.Fatalf("Parse TO DISK failed: %v", err)
	}
	ttl1 := s1[0].(*AlterTableStmt).Clauses[0].TTL[0]
	if ttl1.Action != "TO DISK" || ttl1.Destination != "hot" {
		t.Errorf("ttl1 = %+v, want Action=\"TO DISK\" Destination=\"hot\"", ttl1)
	}

	s2, err := Parse("ALTER TABLE t MODIFY TTL ts TO VOLUME 'ice'")
	if err != nil {
		t.Fatalf("Parse TO VOLUME failed: %v", err)
	}
	ttl2 := s2[0].(*AlterTableStmt).Clauses[0].TTL[0]
	if ttl2.Action != "TO VOLUME" || ttl2.Destination != "ice" {
		t.Errorf("ttl2 = %+v, want Action=\"TO VOLUME\" Destination=\"ice\"", ttl2)
	}

	// After a rollup TTL's SET list, a comma not followed by `name =` starts
	// another TTL rule, and ClickHouse reads even DELETE WHERE k = 1 as one: a
	// TTL on the identifier DELETE with a WHERE. Only MODIFY, REMOVE or
	// MATERIALIZE TTL ends the list.
	s3, err := Parse("ALTER TABLE t MODIFY TTL ts + INTERVAL 1 DAY GROUP BY k SET v = sum(v), DELETE WHERE k = 1, MODIFY TTL ts")
	if err != nil {
		t.Fatalf("Parse rollup TTL + DELETE TTL rule failed: %v", err)
	}
	at3 := s3[0].(*AlterTableStmt)
	if len(at3.Clauses) != 2 {
		t.Fatalf("got %d ALTER clauses, want 2", len(at3.Clauses))
	}
	if len(at3.Clauses[0].TTL) != 2 {
		t.Fatalf("got %d TTL rules, want 2", len(at3.Clauses[0].TTL))
	}
	r3 := at3.Clauses[0].TTL[0]
	if r3.Action != "GROUP BY" || r3.GroupBy == nil || len(r3.GroupBy.Exprs) != 1 || len(r3.Set) != 1 || r3.Set[0].Column != "v" {
		t.Errorf("rollup TTL = %+v", r3)
	}
	if d3 := at3.Clauses[0].TTL[1]; d3.Action != "" || d3.Where == nil {
		t.Errorf("second TTL rule = %+v, want no action with Where", d3)
	}
	if at3.Clauses[1].Action != "MODIFY TTL" {
		t.Errorf("second ALTER clause = %+v, want MODIFY TTL", at3.Clauses[1])
	}

	// Rollup TTL with multiple SET assignments followed by another TTL rule.
	s4, err := Parse("ALTER TABLE t MODIFY TTL ts + INTERVAL 1 DAY GROUP BY k SET v = sum(v), w = max(w), ts + INTERVAL 1 MONTH DELETE WHERE k > 0")
	if err != nil {
		t.Fatalf("Parse rollup TTL + second TTL rule failed: %v", err)
	}
	ttls4 := s4[0].(*AlterTableStmt).Clauses[0].TTL
	if len(ttls4) != 2 {
		t.Fatalf("got %d TTL rules, want 2", len(ttls4))
	}
	if len(ttls4[0].Set) != 2 || ttls4[0].Set[0].Column != "v" || ttls4[0].Set[1].Column != "w" {
		t.Errorf("first TTL rule Set = %+v, want v and w", ttls4[0].Set)
	}
	if ttls4[1].Action != "DELETE" || ttls4[1].Where == nil {
		t.Errorf("second TTL rule = %+v, want DELETE with Where", ttls4[1])
	}

	// CREATE TABLE engine TTL clause with multiple rules.
	s5, err := Parse("CREATE TABLE t (ts DateTime, k UInt64, v UInt64) ENGINE = MergeTree ORDER BY k TTL ts + INTERVAL 1 DAY GROUP BY k SET v = sum(v), ts + INTERVAL 1 MONTH TO DISK 'cold'")
	if err != nil {
		t.Fatalf("Parse CREATE TABLE TTL failed: %v", err)
	}
	ttls5 := s5[0].(*CreateTableStmt).Engine.TTL
	if len(ttls5) != 2 || ttls5[0].Action != "GROUP BY" || ttls5[1].Action != "TO DISK" || ttls5[1].Destination != "cold" {
		t.Errorf("CREATE TABLE TTL = %+v", ttls5)
	}
}

func TestParseInsertValuesSpan(t *testing.T) {
	sql1 := "INSERT INTO t VALUES (1);  -- end"
	stmts1, err := Parse(sql1)
	if err != nil {
		t.Fatalf("Parse(%q) failed: %v", sql1, err)
	}
	ins1 := stmts1[0].(*InsertStmt)
	if got := sql1[ins1.StartPos():ins1.EndPos()]; got != "INSERT INTO t VALUES (1)" {
		t.Errorf("VALUES span text = %q, want %q", got, "INSERT INTO t VALUES (1)")
	}

	sql2 := "INSERT INTO t VALUES (1), (2)   "
	stmts2, err := Parse(sql2)
	if err != nil {
		t.Fatalf("Parse(%q) failed: %v", sql2, err)
	}
	ins2 := stmts2[0].(*InsertStmt)
	if got := sql2[ins2.StartPos():ins2.EndPos()]; got != "INSERT INTO t VALUES (1), (2)" {
		t.Errorf("VALUES span text = %q, want %q", got, "INSERT INTO t VALUES (1), (2)")
	}
}

func TestParseCreateDictionaryAST(t *testing.T) {
	sql := `CREATE DICTIONARY db.dict (
		id UInt64,
		name String DEFAULT 'unknown' INJECTIVE,
		parent_id UInt64 HIERARCHICAL,
		start_date Date,
		end_date Date
	)
	PRIMARY KEY id
	SOURCE(CLICKHOUSE(HOST 'localhost' PORT 9000 TABLE 'src'))
	LAYOUT(RANGE_HASHED())
	LIFETIME(MIN 30 MAX 60)
	RANGE(MIN start_date MAX end_date)
	SETTINGS(max_result_bytes = 1024)`

	stmts, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse CREATE DICTIONARY failed: %v", err)
	}
	d := stmts[0].(*CreateDictionaryStmt)
	if d.Database != "db" || d.Name != "dict" {
		t.Fatalf("dict identity = %s.%s, want db.dict", d.Database, d.Name)
	}
	if len(d.Attributes) != 5 {
		t.Fatalf("got %d attributes, want 5", len(d.Attributes))
	}
	if d.Attributes[1].Name != "name" || d.Attributes[1].Type.Name != "String" || d.Attributes[1].Default == nil || d.Attributes[1].Default.Val != "unknown" || !d.Attributes[1].Injective {
		t.Errorf("attr[1] = %+v", d.Attributes[1])
	}
	if !d.Attributes[2].Hierarchical {
		t.Errorf("attr[2].Hierarchical = false, want true")
	}
	if len(d.PrimaryKey) != 1 {
		t.Errorf("got %d primary key exprs, want 1", len(d.PrimaryKey))
	}
	if d.Source == nil || d.Source.Name != "CLICKHOUSE" || len(d.Source.Params) != 3 {
		t.Errorf("Source = %+v", d.Source)
	}
	if d.Layout == nil || d.Layout.Name != "RANGE_HASHED" {
		t.Errorf("Layout = %+v", d.Layout)
	}
	if d.Lifetime == nil || d.Lifetime.Min != "30" || d.Lifetime.Max != "60" {
		t.Errorf("Lifetime = %+v", d.Lifetime)
	}
	if d.Range == nil || d.Range.Min != "start_date" || d.Range.Max != "end_date" {
		t.Errorf("Range = %+v", d.Range)
	}
	if len(d.Settings) != 1 || d.Settings[0].Name != "max_result_bytes" {
		t.Errorf("Settings = %+v", d.Settings)
	}
}

func TestParseEmptyWindowSpecSpan(t *testing.T) {
	sql1 := "SELECT row_number() OVER (  )"
	stmts1, err := Parse(sql1)
	if err != nil {
		t.Fatalf("Parse(%q) failed: %v", sql1, err)
	}
	spec1 := stmts1[0].(*SelectUnionStmt).Items[0].Columns[0].(*FuncCallExpr).Over
	if spec1 == nil {
		t.Fatalf("expected non-nil WindowSpec")
	}
	if spec1.StartPos() > spec1.EndPos() || sql1[spec1.StartPos():spec1.EndPos()] != "(  )" {
		t.Errorf("OVER (  ) span = [%d, %d) text=%q, want %q", spec1.StartPos(), spec1.EndPos(), sql1[spec1.StartPos():spec1.EndPos()], "(  )")
	}

	sql2 := "SELECT sum(x) OVER w FROM t WINDOW w AS (  )"
	stmts2, err := Parse(sql2)
	if err != nil {
		t.Fatalf("Parse(%q) failed: %v", sql2, err)
	}
	spec2 := stmts2[0].(*SelectUnionStmt).Items[0].Windows[0].Spec
	if spec2 == nil || spec2.StartPos() > spec2.EndPos() || sql2[spec2.StartPos():spec2.EndPos()] != "(  )" {
		t.Errorf("WINDOW w AS (  ) span = [%d, %d), want %q", spec2.StartPos(), spec2.EndPos(), "(  )")
	}
}

func TestParseGrantRevokeBackupAndShowAST(t *testing.T) {
	g1, err := Parse("GRANT SELECT(a, b), INSERT ON db.t TO alice, bob EXCEPT bob WITH GRANT OPTION")
	if err != nil {
		t.Fatalf("Parse GRANT failed: %v", err)
	}
	ac1 := g1[0].(*AccessControlStmt)
	if ac1.Kind != "GRANT" || len(ac1.Privileges) != 2 || ac1.Privileges[0].Name != "SELECT" || len(ac1.Privileges[0].Columns) != 2 ||
		ac1.Privileges[1].Name != "INSERT" || ac1.Target != "db.t" || len(ac1.Grantees) != 2 || ac1.Grantees[0] != "alice" ||
		len(ac1.ExceptGrantees) != 1 || ac1.ExceptGrantees[0] != "bob" || len(ac1.Options) != 1 || ac1.Options[0] != "GRANT OPTION" {
		t.Errorf("GRANT AST = %+v", ac1)
	}

	r1, err := Parse("REVOKE ADMIN OPTION FOR r1 FROM carol")
	if err != nil {
		t.Fatalf("Parse REVOKE failed: %v", err)
	}
	ac2 := r1[0].(*AccessControlStmt)
	if ac2.Kind != "REVOKE" || ac2.RevokeOption != "ADMIN OPTION" || len(ac2.Privileges) != 1 || ac2.Privileges[0].Name != "r1" ||
		len(ac2.Grantees) != 1 || ac2.Grantees[0] != "carol" {
		t.Errorf("REVOKE AST = %+v", ac2)
	}

	u1, err := Parse("CREATE USER u1 IDENTIFIED WITH sha256_password BY 'secret' HOST ANY DEFAULT ROLE r1")
	if err != nil {
		t.Fatalf("Parse CREATE USER failed: %v", err)
	}
	ac3 := u1[0].(*AccessControlStmt)
	if ac3.Kind != "CREATE USER" || len(ac3.Names) != 1 || ac3.Names[0] != "u1" || ac3.IdentifiedWith != "sha256_password" ||
		ac3.IdentifiedBy != "secret" || len(ac3.Hosts) != 1 || ac3.Hosts[0] != "ANY" || len(ac3.DefaultRoles) != 1 || ac3.DefaultRoles[0] != "r1" {
		t.Errorf("CREATE USER AST = %+v", ac3)
	}

	b1, err := Parse("BACKUP TABLE db.t AS db.t2 TO Disk('backups', '1.zip') SETTINGS compression_level = 3 ASYNC")
	if err != nil {
		t.Fatalf("Parse BACKUP failed: %v", err)
	}
	bk := b1[0].(*BackupRestoreStmt)
	if bk.Kind != "BACKUP" || len(bk.Items) != 1 || bk.Items[0].Kind != "TABLE" || bk.Items[0].Name != "db.t" || bk.Items[0].Alias != "db.t2" ||
		bk.DestinationFunc == nil || bk.DestinationFunc.Name != "Disk" || len(bk.Settings) != 1 || !bk.Async {
		t.Errorf("BACKUP AST = %+v", bk)
	}

	sh1, err := Parse("SHOW FULL TABLES FROM analytics NOT ILIKE 'tmp%' LIMIT 10")
	if err != nil {
		t.Fatalf("Parse SHOW TABLES failed: %v", err)
	}
	sh := sh1[0].(*ShowStmt)
	if sh.Kind != "TABLES" || !sh.Full || sh.From != "analytics" || !sh.NotLike || !sh.ILike || sh.Like != "tmp%" || sh.Limit == nil {
		t.Errorf("SHOW TABLES AST = %+v", sh)
	}
}

func TestAllCompatASTSpansValid(t *testing.T) {
	f, err := os.Open("testdata/compat.tsv")
	if err != nil {
		t.Fatalf("open testdata/compat.tsv: %v", err)
	}
	defer f.Close()

	var checkNode func(v reflect.Value, sql string)
	checkNode = func(v reflect.Value, sql string) {
		if !v.IsValid() {
			return
		}
		for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
			if v.IsNil() {
				return
			}
			if n, ok := v.Interface().(Node); ok {
				if n.StartPos() < 0 || n.EndPos() < n.StartPos() || n.EndPos() > len(sql) {
					t.Errorf("invalid span [%d, %d) for %T on SQL %q (len=%d)", n.StartPos(), n.EndPos(), n, sql, len(sql))
				}
			}
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				checkNode(v.Field(i), sql)
			}
		case reflect.Slice, reflect.Array:
			for i := 0; i < v.Len(); i++ {
				checkNode(v.Index(i), sql)
			}
		}
	}

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	checked := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 4 {
			continue
		}
		// Unescape TSV SQL column (only \t, \n, \r, \\)
		sql := parts[len(parts)-1]
		stmts, err := Parse(sql)
		if err != nil {
			continue
		}
		checked++
		for _, stmt := range stmts {
			checkNode(reflect.ValueOf(stmt), sql)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan compat.tsv: %v", err)
	}
	if checked < 1000 {
		t.Fatalf("expected to check >1000 parsed statements, checked %d", checked)
	}
}

func TestParseDropNamesPartitionModeAndGroupByModifierSpans(t *testing.T) {
	d1, err := Parse("DROP USER alice, bob")
	if err != nil {
		t.Fatalf("Parse DROP USER failed: %v", err)
	}
	ds := d1[0].(*DropStmt)
	if len(ds.Names) != 2 || ds.Names[0] != "alice" || ds.Names[1] != "bob" {
		t.Errorf("DROP USER Names = %v, want [alice bob]", ds.Names)
	}

	d2, err := Parse("DROP ROW POLICY p1, p2 ON db.t1, db.t2")
	if err != nil {
		t.Fatalf("Parse DROP ROW POLICY failed: %v", err)
	}
	ds2 := d2[0].(*DropStmt)
	if len(ds2.Names) != 2 || ds2.Names[1] != "p2" || len(ds2.Tables) != 2 || ds2.Tables[1] != "db.t2" {
		t.Errorf("DROP ROW POLICY = %+v", ds2)
	}

	p1, err := Parse("ALTER TABLE t DROP PARTITION ID '202601'")
	if err != nil {
		t.Fatalf("Parse DROP PARTITION ID failed: %v", err)
	}
	pc1 := p1[0].(*AlterTableStmt).Clauses[0].Partition
	if pc1 == nil || !pc1.ByID || pc1.ID != "202601" {
		t.Errorf("PARTITION ID = %+v, want ByID=true ID=\"202601\"", pc1)
	}

	p2, err := Parse("ALTER TABLE t DROP PARTITION '202601'")
	if err != nil {
		t.Fatalf("Parse DROP PARTITION value failed: %v", err)
	}
	pc2 := p2[0].(*AlterTableStmt).Clauses[0].Partition
	if pc2 == nil || pc2.ByID || pc2.Expr == nil {
		t.Errorf("PARTITION value = %+v, want ByID=false with non-nil Expr", pc2)
	}

	gbSQL := "SELECT x FROM t GROUP BY x WITH ROLLUP WITH TOTALS"
	gbStmts, err := Parse(gbSQL)
	if err != nil {
		t.Fatalf("Parse(%q) failed: %v", gbSQL, err)
	}
	gb := gbStmts[0].(*SelectUnionStmt).Items[0].GroupBy
	if got := gbSQL[gb.StartPos():gb.EndPos()]; got != "GROUP BY x WITH ROLLUP WITH TOTALS" {
		t.Errorf("GroupBy span text = %q, want %q", got, "GROUP BY x WITH ROLLUP WITH TOTALS")
	}
}

func TestParseUnsupportedEscapesTableParamsAndColumnAliases(t *testing.T) {
	// Unsupported backslash escapes must preserve both the backslash and the following char.
	e1, err := Parse(`SELECT '\d+', '\s\w', '\xZZ', '\n'`)
	if err != nil {
		t.Fatalf("Parse escape string failed: %v", err)
	}
	cols := e1[0].(*SelectUnionStmt).Items[0].Columns
	wantVals := []string{`\d+`, `\s\w`, `\xZZ`, "\n"}
	for i, want := range wantVals {
		lit := cols[i].(*LiteralExpr)
		if lit.Val != want {
			t.Errorf("col[%d].Val = %q, want %q", i, lit.Val, want)
		}
	}

	// Identifier query parameters in FROM / DESCRIBE table references.
	p1, err := Parse("SELECT * FROM {mytablename:Identifier}")
	if err != nil {
		t.Fatalf("Parse table query param failed: %v", err)
	}
	tbl1 := p1[0].(*SelectUnionStmt).Items[0].From.Table
	if tbl1.TableParam == nil || tbl1.TableParam.Name != "mytablename" || tbl1.TableParam.Type.Name != "Identifier" {
		t.Errorf("TableParam = %+v, want mytablename:Identifier", tbl1.TableParam)
	}

	p2, err := Parse("SELECT * FROM {mydb:Identifier}.{mytbl:Identifier}")
	if err != nil {
		t.Fatalf("Parse db.table query params failed: %v", err)
	}
	tbl2 := p2[0].(*SelectUnionStmt).Items[0].From.Table
	if tbl2.DatabaseParam == nil || tbl2.DatabaseParam.Name != "mydb" || tbl2.TableParam == nil || tbl2.TableParam.Name != "mytbl" {
		t.Errorf("qualified TableParams = db:%+v tbl:%+v", tbl2.DatabaseParam, tbl2.TableParam)
	}

	// Column-name lists following table aliases (AS and bare forms).
	a1, err := Parse("SELECT a, b FROM (SELECT 1, 2) AS x(a, b)")
	if err != nil {
		t.Fatalf("Parse AS x(a, b) failed: %v", err)
	}
	at1 := a1[0].(*SelectUnionStmt).Items[0].From.Table
	if at1.Alias != "x" || len(at1.ColumnAliases) != 2 || at1.ColumnAliases[0] != "a" || at1.ColumnAliases[1] != "b" {
		t.Errorf("AS x(a, b) TableExpr = %+v", at1)
	}

	a2, err := Parse("SELECT a, b FROM (SELECT 1, 2) x(a, b)")
	if err != nil {
		t.Fatalf("Parse bare x(a, b) failed: %v", err)
	}
	at2 := a2[0].(*SelectUnionStmt).Items[0].From.Table
	if at2.Alias != "x" || len(at2.ColumnAliases) != 2 || at2.ColumnAliases[0] != "a" || at2.ColumnAliases[1] != "b" {
		t.Errorf("bare x(a, b) TableExpr = %+v", at2)
	}
}

func TestParseWithLambdaAlterModifiersFilterAliasesStarTransformersAndOffsetFetch(t *testing.T) {
	// 1. Scalar WITH lambda expression.
	w1, err := Parse("WITH (x) -> x + 1 AS f SELECT f(1)")
	if err != nil {
		t.Fatalf("Parse WITH lambda failed: %v", err)
	}
	with0 := w1[0].(*SelectUnionStmt).Items[0].With[0]
	lam, ok := with0.Expr.(*LambdaExpr)
	if !ok || with0.Alias != "f" || len(lam.Params) != 1 || lam.Params[0] != "x" {
		t.Errorf("WITH lambda = %+v, expr=%T", with0, with0.Expr)
	}

	// 2. MODIFY COLUMN AFTER / FIRST + conditional IF EXISTS / IF NOT EXISTS flags.
	m1, err := Parse("ALTER TABLE t ADD COLUMN IF NOT EXISTS a UInt64, MODIFY COLUMN IF EXISTS c UInt64 AFTER d, MODIFY COLUMN e String FIRST, DROP COLUMN IF EXISTS old")
	if err != nil {
		t.Fatalf("Parse ALTER modifiers/flags failed: %v", err)
	}
	cls := m1[0].(*AlterTableStmt).Clauses
	if len(cls) != 4 {
		t.Fatalf("got %d clauses, want 4", len(cls))
	}
	if !cls[0].IfNotExists || cls[0].IfExists {
		t.Errorf("ADD COLUMN flags = %+v, want IfNotExists=true", cls[0])
	}
	if !cls[1].IfExists || cls[1].After != "d" || cls[1].First {
		t.Errorf("MODIFY COLUMN AFTER = %+v, want IfExists=true After=\"d\"", cls[1])
	}
	if cls[2].IfExists || !cls[2].First || cls[2].After != "" {
		t.Errorf("MODIFY COLUMN FIRST = %+v, want First=true", cls[2])
	}
	if !cls[3].IfExists {
		t.Errorf("DROP COLUMN IF EXISTS = %+v, want IfExists=true", cls[3])
	}

	// 3. Explicit aliases in WHERE / PREWHERE / HAVING / QUALIFY.
	f1, err := Parse("SELECT y WHERE 1 AS y")
	if err != nil {
		t.Fatalf("Parse WHERE alias failed: %v", err)
	}
	if al, ok := f1[0].(*SelectUnionStmt).Items[0].Where.(*AliasExpr); !ok || al.Alias != "y" {
		t.Errorf("WHERE = %T %+v, want AliasExpr(y)", f1[0].(*SelectUnionStmt).Items[0].Where, f1[0].(*SelectUnionStmt).Items[0].Where)
	}
	if _, err := Parse("SELECT y FROM t PREWHERE x > 0 AS p WHERE p HAVING count() > 0 AS h QUALIFY row_number() OVER () = 1 AS q"); err != nil {
		t.Fatalf("Parse PREWHERE/HAVING/QUALIFY aliases failed: %v", err)
	}

	// 4. Star APPLY and REPLACE column transformers (including chaining).
	st1, err := Parse("SELECT * APPLY(toString) FROM numbers(3)")
	if err != nil {
		t.Fatalf("Parse * APPLY failed: %v", err)
	}
	star1 := st1[0].(*SelectUnionStmt).Items[0].Columns[0].(*StarExpr)
	if len(star1.Transformers) != 1 || star1.Transformers[0].Kind != "APPLY" || star1.Transformers[0].ApplyFunc != "toString" {
		t.Errorf("star1 transformers = %+v", star1.Transformers)
	}

	st2, err := Parse("SELECT * EXCEPT (secret) REPLACE(number + 1 AS number) APPLY(x -> x * 2) FROM numbers(3)")
	if err != nil {
		t.Fatalf("Parse chained star transformers failed: %v", err)
	}
	star2 := st2[0].(*SelectUnionStmt).Items[0].Columns[0].(*StarExpr)
	if len(star2.Transformers) != 3 || star2.Transformers[0].Kind != "EXCEPT" || star2.Transformers[1].Kind != "REPLACE" ||
		len(star2.Transformers[1].Replaces) != 1 || star2.Transformers[1].Replaces[0].Alias != "number" ||
		star2.Transformers[2].Kind != "APPLY" || star2.Transformers[2].ApplyExpr == nil {
		t.Errorf("star2 transformers = %+v", star2.Transformers)
	}

	// 5. Standalone OFFSET and SQL-standard FETCH pagination.
	o1, err := Parse("SELECT number FROM numbers(10) ORDER BY number OFFSET 5 ROWS FETCH NEXT 2 ROWS ONLY")
	if err != nil {
		t.Fatalf("Parse OFFSET FETCH failed: %v", err)
	}
	lim1 := o1[0].(*SelectUnionStmt).Items[0].Limit
	if lim1 == nil || lim1.Offset == nil || lim1.Limit == nil || lim1.WithTies {
		t.Errorf("OFFSET FETCH Limit = %+v", lim1)
	}

	o2, err := Parse("SELECT number FROM numbers(10) OFFSET 5")
	if err != nil {
		t.Fatalf("Parse standalone OFFSET failed: %v", err)
	}
	if lim2 := o2[0].(*SelectUnionStmt).Items[0].Limit; lim2 == nil || lim2.Offset == nil || lim2.Limit != nil {
		t.Errorf("standalone OFFSET Limit = %+v", lim2)
	}

	o3, err := Parse("SELECT number FROM numbers(10) ORDER BY number FETCH FIRST 2 ROWS WITH TIES")
	if err != nil {
		t.Fatalf("Parse FETCH WITH TIES failed: %v", err)
	}
	if lim3 := o3[0].(*SelectUnionStmt).Items[0].Limit; lim3 == nil || lim3.Limit == nil || !lim3.WithTies {
		t.Errorf("FETCH WITH TIES Limit = %+v", lim3)
	}

	if _, err := Parse("SELECT number FROM t FETCH FIRST 2 ROWS ONLY"); err != nil {
		t.Fatalf("Parse FETCH after table name failed: %v", err)
	}
}

func TestParseOutputSettingsProjectionColTupleSpanTableFuncExprAndNullability(t *testing.T) {
	// 1. Query-level SETTINGS in output clause (ALTER mutation + SELECT after FORMAT).
	a1, err := Parse("ALTER TABLE t DELETE WHERE id = 1 SETTINGS mutations_sync = 2")
	if err != nil {
		t.Fatalf("Parse ALTER ... SETTINGS failed: %v", err)
	}
	at1 := a1[0].(*AlterTableStmt)
	if at1.Output == nil || len(at1.Output.Settings) != 1 || at1.Output.Settings[0].Name != "mutations_sync" {
		t.Errorf("ALTER Output = %+v, want mutations_sync setting", at1.Output)
	}
	s1, err := Parse("SELECT 1 FORMAT JSON SETTINGS max_threads = 2")
	if err != nil {
		t.Fatalf("Parse SELECT ... FORMAT ... SETTINGS failed: %v", err)
	}
	su1 := s1[0].(*SelectUnionStmt)
	if su1.Output == nil || su1.Output.Format != "JSON" || len(su1.Output.Settings) != 1 || su1.Output.Settings[0].Name != "max_threads" {
		t.Errorf("SELECT Output = %+v, want Format=JSON and max_threads setting", su1.Output)
	}

	// 2. Column named `projection` with parameterized type vs real PROJECTION definition.
	c1, err := Parse("CREATE TABLE t (projection Array(String), PROJECTION p (SELECT projection)) ENGINE = Memory")
	if err != nil {
		t.Fatalf("Parse projection column failed: %v", err)
	}
	cols1 := c1[0].(*CreateTableStmt).Columns
	if len(cols1) != 2 || cols1[0].Kind != "COLUMN" || cols1[0].Name != "projection" || cols1[0].Type.Name != "Array" || cols1[1].Kind != "PROJECTION" {
		t.Errorf("CREATE TABLE cols = %+v", cols1)
	}

	// 3. Split numeric tuple-access spans (`t.1.2` inner node spans `t.1`, outer spans `t.1.2`).
	tupSQL := "SELECT t.1.2"
	tupStmts, err := Parse(tupSQL)
	if err != nil {
		t.Fatalf("Parse(%q) failed: %v", tupSQL, err)
	}
	outer := tupStmts[0].(*SelectUnionStmt).Items[0].Columns[0].(*TupleAccessExpr)
	inner := outer.Expr.(*TupleAccessExpr)
	if got := tupSQL[inner.StartPos():inner.EndPos()]; got != "t.1" {
		t.Errorf("inner tuple access span text = %q, want %q", got, "t.1")
	}
	if got := tupSQL[outer.StartPos():outer.EndPos()]; got != "t.1.2" {
		t.Errorf("outer tuple access span text = %q, want %q", got, "t.1.2")
	}

	// 4. Full expressions in table-function arguments.
	for _, sql := range []string{
		"SELECT * FROM numbers(10 + 1)",
		"SELECT * FROM numbers(toUInt64(10) + 1)",
	} {
		if _, err := Parse(sql); err != nil {
			t.Fatalf("Parse(%q) failed: %v", sql, err)
		}
	}

	// 5. Explicit NULL / NOT NULL column modifiers in CREATE TABLE and ALTER TABLE.
	n1, err := Parse("CREATE TABLE t (id UInt64 NOT NULL, name String NULL DEFAULT '') ENGINE = Memory")
	if err != nil {
		t.Fatalf("Parse CREATE TABLE NULL/NOT NULL failed: %v", err)
	}
	ncols := n1[0].(*CreateTableStmt).Columns
	if ncols[0].Nullability != "NOT NULL" || ncols[1].Nullability != "NULL" || ncols[1].DefaultKind != "DEFAULT" {
		t.Errorf("column nullability = [%+v, %+v]", ncols[0], ncols[1])
	}
	n2, err := Parse("ALTER TABLE t ADD COLUMN id UInt64 NOT NULL")
	if err != nil {
		t.Fatalf("Parse ALTER ADD COLUMN NOT NULL failed: %v", err)
	}
	if n2[0].(*AlterTableStmt).Clauses[0].Element.Nullability != "NOT NULL" {
		t.Errorf("ALTER ADD COLUMN Nullability = %q, want NOT NULL", n2[0].(*AlterTableStmt).Clauses[0].Element.Nullability)
	}
}

func TestParseAlterModifySettingFollowedByAction(t *testing.T) {
	s1, err := Parse("ALTER TABLE t MODIFY SETTING index_granularity = 8192, ADD COLUMN c UInt64")
	if err != nil {
		t.Fatalf("Parse MODIFY SETTING + ADD COLUMN failed: %v", err)
	}
	cls1 := s1[0].(*AlterTableStmt).Clauses
	if len(cls1) != 2 || cls1[0].Action != "MODIFY SETTING" || len(cls1[0].Settings) != 1 || cls1[1].Action != "ADD COLUMN" || cls1[1].Name != "c" {
		t.Errorf("clauses = %+v, want MODIFY SETTING + ADD COLUMN c", cls1)
	}

	s2, err := Parse("ALTER TABLE t MODIFY SETTING index_granularity = 8192, min_bytes_for_wide_part = 0, DROP COLUMN old")
	if err != nil {
		t.Fatalf("Parse multi MODIFY SETTING + DROP COLUMN failed: %v", err)
	}
	cls2 := s2[0].(*AlterTableStmt).Clauses
	if len(cls2) != 2 || len(cls2[0].Settings) != 2 || cls2[1].Action != "DROP COLUMN" || cls2[1].Name != "old" {
		t.Errorf("clauses = %+v, want 2 settings + DROP COLUMN old", cls2)
	}
}

func TestParseOrderByAliasEngineLessTableNestedDictOptionalGranularityAndSelectAll(t *testing.T) {
	// 1. Explicit alias in ORDER BY.
	o1, err := Parse("SELECT 1 ORDER BY 1 AS x DESC")
	if err != nil {
		t.Fatalf("Parse ORDER BY alias failed: %v", err)
	}
	ob := o1[0].(*SelectUnionStmt).Items[0].OrderBy[0]
	if al, ok := ob.Expr.(*AliasExpr); !ok || al.Alias != "x" || !ob.Desc {
		t.Errorf("ORDER BY item = %+v, expr=%T", ob, ob.Expr)
	}

	// 2. CREATE TABLE with ORDER BY / SETTINGS and no explicit ENGINE.
	c1, err := Parse("CREATE TABLE t (x UInt64) ORDER BY x SETTINGS index_granularity = 4096")
	if err != nil {
		t.Fatalf("Parse engine-less CREATE TABLE failed: %v", err)
	}
	ct1 := c1[0].(*CreateTableStmt)
	if ct1.Engine == nil || ct1.Engine.Engine != nil || len(ct1.Engine.OrderBy) != 1 || len(ct1.Engine.Settings) != 1 {
		t.Errorf("engine-less CREATE TABLE Engine = %+v", ct1.Engine)
	}

	// 3. Nested dictionary SOURCE configuration blocks (REPLICA(...)).
	d1, err := Parse("CREATE DICTIONARY d (id UInt64, v String) PRIMARY KEY id SOURCE(MYSQL(PORT 3306 USER 'root' REPLICA(HOST 's1' PRIORITY 1) REPLICA(HOST 's2' PRIORITY 2) DB 'db' TABLE 't')) LAYOUT(FLAT()) LIFETIME(60)")
	if err != nil {
		t.Fatalf("Parse nested dict source failed: %v", err)
	}
	src := d1[0].(*CreateDictionaryStmt).Source
	if src == nil || src.Name != "MYSQL" || len(src.Children) != 2 || src.Children[0].Name != "REPLICA" || len(src.Children[0].Params) != 2 || len(src.Params) != 4 {
		t.Errorf("nested dict Source = %+v", src)
	}

	// 4. Optional GRANULARITY on data-skipping INDEX (both ALTER and CREATE TABLE).
	i1, err := Parse("ALTER TABLE t ADD INDEX i x TYPE minmax")
	if err != nil {
		t.Fatalf("Parse ADD INDEX without GRANULARITY failed: %v", err)
	}
	ie1 := i1[0].(*AlterTableStmt).Clauses[0].Element
	if ie1 == nil || ie1.Kind != "INDEX" || ie1.Type.Name != "minmax" || ie1.Granularity != "" {
		t.Errorf("INDEX element = %+v", ie1)
	}
	if _, err := Parse("CREATE TABLE t (x UInt64, INDEX i x TYPE minmax) ENGINE = MergeTree ORDER BY x"); err != nil {
		t.Fatalf("Parse CREATE TABLE INDEX without GRANULARITY failed: %v", err)
	}

	// 5. SELECT ALL modifier (with column, *, and literal).
	sa1, err := Parse("SELECT ALL x FROM t")
	if err != nil {
		t.Fatalf("Parse SELECT ALL x failed: %v", err)
	}
	sel1 := sa1[0].(*SelectUnionStmt).Items[0]
	if !sel1.All || sel1.From == nil || sel1.From.Table.Table != "t" {
		t.Errorf("SELECT ALL x = %+v", sel1)
	}
	if id, ok := sel1.Columns[0].(*IdentExpr); !ok || len(id.Parts) != 1 || id.Parts[0] != "x" {
		t.Errorf("SELECT ALL x col[0] = %T %+v, want IdentExpr(x)", sel1.Columns[0], sel1.Columns[0])
	}

	sa2, err := Parse("SELECT ALL * FROM t")
	if err != nil {
		t.Fatalf("Parse SELECT ALL * failed: %v", err)
	}
	sel2 := sa2[0].(*SelectUnionStmt).Items[0]
	if !sel2.All || sel2.From == nil {
		t.Errorf("SELECT ALL * = %+v", sel2)
	}
	if _, ok := sel2.Columns[0].(*StarExpr); !ok {
		t.Errorf("SELECT ALL * col[0] = %T, want *StarExpr", sel2.Columns[0])
	}

	sa3, err := Parse("SELECT ALL 1")
	if err != nil {
		t.Fatalf("Parse SELECT ALL 1 failed: %v", err)
	}
	if !sa3[0].(*SelectUnionStmt).Items[0].All {
		t.Errorf("SELECT ALL 1 All = false, want true")
	}
}

func TestParseMaterializedViewEngineOrderByAndViewComment(t *testing.T) {
	mvSQL := "CREATE MATERIALIZED VIEW mv_own ENGINE = MergeTree ORDER BY id AS SELECT id FROM src COMMENT 'mv note'"
	stmts, err := Parse(mvSQL)
	if err != nil {
		t.Fatalf("Parse materialized view failed: %v", err)
	}
	mv := stmts[0].(*CreateViewStmt)
	if !mv.Materialized || mv.Name != "mv_own" || mv.Engine == nil || len(mv.Engine.OrderBy) != 1 || mv.Subquery == nil || mv.Comment != "mv note" {
		t.Errorf("CreateViewStmt = %+v", mv)
	}

	vSQL := "CREATE VIEW v AS SELECT id FROM src COMMENT 'view note'"
	vStmts, err := Parse(vSQL)
	if err != nil {
		t.Fatalf("Parse view with COMMENT failed: %v", err)
	}
	v := vStmts[0].(*CreateViewStmt)
	if v.Materialized || v.Name != "v" || v.Subquery == nil || v.Comment != "view note" {
		t.Errorf("CreateViewStmt = %+v", v)
	}

	v2SQL := "CREATE VIEW v2 AS SELECT 1 COMMENT 'lit note'"
	v2Stmts, err := Parse(v2SQL)
	if err != nil {
		t.Fatalf("Parse projection view with COMMENT failed: %v", err)
	}
	v2 := v2Stmts[0].(*CreateViewStmt)
	if v2.Comment != "lit note" {
		t.Errorf("v2.Comment = %q, want %q", v2.Comment, "lit note")
	}
}

func TestParseExplainCurrentTransaction(t *testing.T) {
	stmts, err := Parse("EXPLAIN CURRENT TRANSACTION FORMAT JSON")
	if err != nil {
		t.Fatalf("Parse EXPLAIN CURRENT TRANSACTION failed: %v", err)
	}
	ex := stmts[0].(*ExplainStmt)
	if ex.Kind != "CURRENT TRANSACTION" || ex.Inner != nil || ex.Output == nil || ex.Output.Format != "JSON" {
		t.Errorf("EXPLAIN CURRENT TRANSACTION = %+v", ex)
	}
}

func TestParseClickHouseAliasSystemPoliciesProjectionAndWindowCompat(t *testing.T) {
	// 1. Expression / alias keyword fallbacks (interval units, NULL, IS).
	aliasCases := []struct {
		sql       string
		wantAlias string
	}{
		{"SELECT 1e1 NANOSECOND", "NANOSECOND"},
		{"SELECT INTERVAL NANOSECOND", "NANOSECOND"},
		{"SELECT IS NULL", "NULL"},
		{"SELECT 1e1 NULL", "NULL"},
		{"SELECT 1e1 IS", "IS"},
		{"SELECT 1e1 AS NULL", "NULL"},
		{"SELECT 1e1 IS FROM t", "IS"},
		{"SELECT INTERVAL 1 DAY HOUR", "HOUR"},
	}
	for _, tc := range aliasCases {
		stmts, err := Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
		}
		col := stmts[0].(*SelectUnionStmt).Items[0].Columns[0]
		al, ok := col.(*AliasExpr)
		if !ok || al.Alias != tc.wantAlias {
			t.Errorf("Parse(%q) col = %T (%+v), want *AliasExpr(%q)", tc.sql, col, col, tc.wantAlias)
		}
	}
	// Verify real IS NULL and INTERVAL expressions still parse as predicates/intervals.
	if stmts, err := Parse("SELECT 1e1 IS NULL"); err != nil {
		t.Fatalf("Parse IS NULL failed: %v", err)
	} else if _, ok := stmts[0].(*SelectUnionStmt).Items[0].Columns[0].(*IsNullExpr); !ok {
		t.Errorf("SELECT 1e1 IS NULL col = %T, want *IsNullExpr", stmts[0].(*SelectUnionStmt).Items[0].Columns[0])
	}
	if stmts, err := Parse("SELECT INTERVAL 1 DAY"); err != nil {
		t.Fatalf("Parse INTERVAL 1 DAY failed: %v", err)
	} else if _, ok := stmts[0].(*SelectUnionStmt).Items[0].Columns[0].(*IntervalExpr); !ok {
		t.Errorf("SELECT INTERVAL 1 DAY col = %T, want *IntervalExpr", stmts[0].(*SelectUnionStmt).Items[0].Columns[0])
	}

	// 2. Generic SYSTEM START|STOP <tableIdentifier> forms.
	sysCases := []struct {
		sql        string
		wantAction string
		wantTarget string
	}{
		{"SYSTEM START a", "START", "a"},
		{"SYSTEM START a.b", "START", "a.b"},
		{"SYSTEM STOP a", "STOP", "a"},
		{"SYSTEM START SENDS", "START", "SENDS"},
		{"SYSTEM START REPLICATED", "START", "REPLICATED"},
		{"SYSTEM START DISTRIBUTED", "START", "DISTRIBUTED"},
		{"SYSTEM START DISTRIBUTED SENDS", "START DISTRIBUTED SENDS", ""},
		{"SYSTEM START DISTRIBUTED SENDS db.t", "START DISTRIBUTED SENDS", "db.t"},
	}
	for _, tc := range sysCases {
		stmts, err := Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
		}
		sys := stmts[0].(*SystemStmt)
		if sys.Action != tc.wantAction || sys.Target != tc.wantTarget {
			t.Errorf("Parse(%q) = {Action:%q Target:%q}, want {%q %q}", tc.sql, sys.Action, sys.Target, tc.wantAction, tc.wantTarget)
		}
	}

	// 3. SHOW [ROW] POLICIES optional target forms (short name and ON target).
	policyCases := []struct {
		sql        string
		wantTarget string
	}{
		{"SHOW POLICIES zzz", "zzz"},
		{"SHOW POLICIES a", "a"},
		{"SHOW POLICIES ON", "ON"},
		{"SHOW POLICIES 'my_policy'", "my_policy"},
		{"SHOW POLICIES ON db.tbl", "db.tbl"},
		{"SHOW POLICIES ON *.*", "*.*"},
		{"SHOW ROW POLICIES zzz", "zzz"},
		{"SHOW POLICIES FORMAT JSON", ""},
	}
	for _, tc := range policyCases {
		stmts, err := Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
		}
		show := stmts[0].(*ShowStmt)
		if show.Target != tc.wantTarget {
			t.Errorf("Parse(%q) Target = %q, want %q", tc.sql, show.Target, tc.wantTarget)
		}
	}

	// 4. Empty GROUP BY inside a projection definition.
	for _, sql := range []string{
		"ALTER TABLE a ADD PROJECTION a(SELECT * GROUP BY)",
		"CREATE TABLE t (a Int32, PROJECTION p (SELECT * GROUP BY)) ENGINE = MergeTree ORDER BY a",
	} {
		if _, err := Parse(sql); err != nil {
			t.Fatalf("Parse(%q) failed: %v", sql, err)
		}
	}

	// 5. ClickHouse ParserWindowDefinition OVER(...) forms.
	windowAccept := []string{
		"SELECT a() OVER(PARTITION BY * zzz)",
		"SELECT a() OVER(ORDER BY 1e1 zzz)",
		"SELECT a() OVER(ROWS CURRENT ROW zzz)",
		"SELECT a() OVER(ROWS CURRENT_USER ROW)",
		"SELECT a() OVER(ROWS BETWEEN ROW AND CURRENT ROW)",
		"SELECT a() OVER(ROWS BETWEEN CURRENT AND CURRENT ROW)",
		"SELECT a() OVER(ROWS BETWEEN CURRENT ROW AND CURRENT ROW zzz)",
		"SELECT a() OVER(ROWS UNBOUNDED PRECEDING zzz)",
		"SELECT a() OVER(ROWS UNION PRECEDING)",
		"SELECT a() OVER(ROWS UNION FOLLOWING)",
		"SELECT a() OVER(ROWS 1e1 PRECEDING zzz)",
		"SELECT a() OVER(ROWS 1e1 PREWHERE)",
		"SELECT a() OVER(ROWS 1e1 FOLLOWING zzz)",
		"SELECT a() OVER(ROWS 1e1 FOR)",
		"SELECT a() OVER(ROWS 1e1 + 2 PRECEDING)",
		"SELECT a() OVER(GROUPS UNBOUNDED PRECEDING)",
	}
	for _, sql := range windowAccept {
		if _, err := Parse(sql); err != nil {
			t.Fatalf("Parse(%q) failed: %v", sql, err)
		}
	}

	// Negative controls: invalid ClickHouse SQL must stay rejected.
	rejects := []string{
		"SELECT 1e1 2e1",
		"SELECT 1e1 IS NOT",
		"SELECT 1e1 IS a",
		"SELECT IS NOT NULL",
		"SELECT 1 FROM t WHERE a IS",
		"SELECT INTERVAL 1e1",
		"SYSTEM START",
		"SYSTEM START a b",
		"SYSTEM START DISTRIBUTED a",
		"SYSTEM START REPLICATED a",
		"SYSTEM START TTL a",
		"SHOW POLICIES a b",
		"SHOW POLICIES a.b",
		"SHOW POLICIES ON a zzz",
		"SHOW POLICIES ON 'a'",
		"SHOW POLICIES ON FORMAT JSON",
		"SHOW POLICIES 123",
		"SELECT * GROUP BY",
		"ALTER TABLE a ADD PROJECTION a(SELECT * GROUP BY",
		"ALTER TABLE a ADD PROJECTION a(SELECT * GROUP BY ORDER BY a)",
		"ALTER TABLE a ADD PROJECTION a(SELECT * ORDER BY)",
		"SELECT a() OVER(ROWS 1e1)",
		"SELECT a() OVER(ROWS ROW)",
		"SELECT a() OVER(ROWS CURRENT)",
		"SELECT a() OVER(ROWS (1 zzz)",
		"SELECT a() OVER(ROWS + zzz)",
		"SELECT a() OVER(zzz ROWS 1e1 PREWHERE)",
		"SELECT a() OVER(ROWS UNBOUNDED FOLLOWING)",
		"SELECT a() OVER(ROWS UNBOUNDED FOLLOWING zzz)",
		"SELECT a() OVER(ROWS BETWEEN CURRENT ROW AND UNBOUNDED PRECEDING)",
		"SELECT a() OVER(ROWS BETWEEN CURRENT ROW AND ROW)",
		"SELECT a() OVER(PARTITION BY a zzz yyy)",
		"SELECT a() OVER(`` ``)",
	}
	for _, sql := range rejects {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want error", sql)
		}
	}
}

// TestParseDataTypeArguments holds a parametric type's arguments to
// ParserDataType's general rule: an enum element chain, a literal or a data
// type, never an expression. Each verdict is ClickHouse 26.8.38.9's
// EXPLAIN AST answer.
func TestParseDataTypeArguments(t *testing.T) {
	accept := []string{
		"SELECT 1 :: a(-1)",
		"SELECT 1 :: a(+1)",
		"SELECT 1 :: a(-inf)",
		"SELECT 1 :: a(X'00')",
		"SELECT 1 :: a(NULL)",
		"SELECT 1 :: a(true)",
		"SELECT 1 :: a(a)",
		"SELECT 1 :: a(f(1))",
		"SELECT 1 :: a(Nested(a b))",
		"SELECT 1 :: a(1,)",
		"SELECT 1 :: a(1 = 1)",
		"SELECT 1 :: a(NULL = 1)",
		"SELECT 1 :: a('a' = 1 = 2)",
		"SELECT 1 :: a('a' = 1, 2)",
		"SELECT 1 :: a('a' = 1, 'b')",
		"SELECT 1 :: Enum8('a' = -1, 'b' = 2)",
	}
	reject := []string{
		"SELECT 1 :: a((1))",
		"SELECT 1 :: a((1, 2))",
		"SELECT 1 :: a(())",
		"SELECT 1 :: a(1, (a))",
		"SELECT 1 :: a([1, 2])",
		"SELECT 1 :: a(1 + 1)",
		"SELECT 1 :: a(a * 2)",
		"SELECT 1 :: a(-a)",
		"SELECT 1 :: a(a = 1)",
		"SELECT 1 :: a(a.b)",
		"SELECT 1 :: a(a::b)",
		"SELECT 1 :: a(DATE '2020-01-01')",
		"SELECT 1 :: a(INTERVAL 1 DAY)",
		"SELECT 1 :: a(f(1)(2))",
		"CREATE TABLE t (c a((a))) ENGINE = Memory",
	}
	for _, sql := range accept {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v, want accept", sql, err)
		}
	}
	for _, sql := range reject {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted, want a syntax error", sql)
		}
	}
}

// TestParseProjectionListsAndTTLListBoundary holds a projection's GROUP BY
// and ORDER BY to plain expression lists, and an ALTER's MODIFY TTL list to
// ending only before MODIFY, REMOVE or MATERIALIZE TTL. Each verdict is
// ClickHouse 26.8.38.9's EXPLAIN AST answer.
func TestParseProjectionListsAndTTLListBoundary(t *testing.T) {
	accepts := []string{
		"ALTER TABLE a ADD PROJECTION a(SELECT * GROUP BY CUBE())",
		"ALTER TABLE a ADD PROJECTION a(SELECT * GROUP BY ROLLUP(a), CUBE())",
		"ALTER TABLE a ADD PROJECTION a(SELECT * GROUP BY ALL, a)",
		"ALTER TABLE a ADD PROJECTION a(SELECT * GROUP BY a.*)",
		"ALTER TABLE a ADD PROJECTION a(SELECT * GROUP BY CUBE() ORDER BY a)",
		"CREATE TABLE t (a Int8, PROJECTION p (SELECT * GROUP BY CUBE())) ENGINE a",
		"ALTER TABLE a MODIFY TTL 1e1, DELETE",
		"ALTER TABLE a MODIFY TTL 1, DELETE WHERE 1, 2",
		"ALTER TABLE a MODIFY TTL a WHERE 1",
		"CREATE TABLE a (b Int8) ENGINE a TTL a WHERE 1",
		"ALTER TABLE a MODIFY TTL 1, MODIFY TTL 2, 3",
		"ALTER TABLE a MODIFY TTL 1, REMOVE TTL, FREEZE",
		"ALTER TABLE a MODIFY TTL 1, MATERIALIZE TTL, FREEZE",
		"ALTER TABLE a MATERIALIZE TTL IN PARTITION ID '1'",
	}
	for _, sql := range accepts {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q) failed: %v", sql, err)
		}
	}
	rejects := []string{
		"ALTER TABLE a ADD PROJECTION a(SELECT * GROUP BY GROUPING SETS((a)))",
		"ALTER TABLE a ADD PROJECTION a(SELECT * GROUP BY a AS b)",
		"ALTER TABLE a ADD PROJECTION a(SELECT * GROUP BY a b)",
		"ALTER TABLE a ADD PROJECTION a(SELECT * GROUP BY a,)",
		"ALTER TABLE a ADD PROJECTION a(SELECT * ORDER BY a AS b)",
		"ALTER TABLE a MODIFY TTL 1, DROP COLUMN a",
		"ALTER TABLE a MODIFY TTL 1, MODIFY ORDER BY a",
		"ALTER TABLE a MODIFY TTL 1, FREEZE, DROP COLUMN a",
		"ALTER TABLE a MODIFY TTL a TO DISK 'x' WHERE 1",
		"ALTER TABLE a MATERIALIZE TTL IF EXISTS",
		"ALTER TABLE a MATERIALIZE TTL a",
	}
	for _, sql := range rejects {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want error", sql)
		}
	}
}
