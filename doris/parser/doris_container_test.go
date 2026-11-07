package parser

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/doris/internal/testcontainer"
)

// dorisContainer wraps an Apache Doris standalone container used as a parse
// oracle. Doris speaks the MySQL wire protocol on port 9030 (not 3306), so we
// connect with go-sql-driver/mysql via a raw GenericContainer (there is no
// testcontainers Doris module).
//
// dyrnq/doris publishes native arm64 and amd64 images, so unlike the StarRocks
// harness no platform override is needed.
type dorisContainer struct {
	db  *sql.DB
	ctx context.Context
}

var (
	dorisServer  testcontainer.Container
	dorisOnce    sync.Once
	dorisInst    *dorisContainer
	dorisInitErr error
)

// startDoris starts (once) the Doris container and returns a ready
// dorisContainer. Fails the test when the container cannot start (mirrors
// starrocks_container_test.go).
func startDoris(t *testing.T) *dorisContainer {
	t.Helper()

	dorisOnce.Do(func() {
		ctx := context.Background()
		container, err := testcontainer.Start(ctx)
		dorisServer = container
		if err != nil {
			dorisInitErr = fmt.Errorf("start Doris container: %w", err)
			return
		}
		host, err := container.Host(ctx)
		if err != nil {
			_ = testcontainer.TerminateContainer(container)
			dorisInitErr = err
			return
		}
		port, err := container.MappedPort(ctx, "9030")
		if err != nil {
			_ = testcontainer.TerminateContainer(container)
			dorisInitErr = err
			return
		}
		dsn := func(db string) string {
			return fmt.Sprintf("root@tcp(%s:%s)/%s?multiStatements=true&timeout=10s", host, port.Port(), db)
		}
		db, err := sql.Open("mysql", dsn(""))
		if err != nil {
			_ = testcontainer.TerminateContainer(container)
			dorisInitErr = err
			return
		}

		// Active readiness: the FE accepts connections before the BE has
		// registered. Poll until SELECT 1 succeeds.
		deadline := time.Now().Add(5 * time.Minute)
		for {
			ok := func() bool {
				cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				var one int
				return db.QueryRowContext(cctx, "SELECT 1").Scan(&one) == nil && one == 1
			}()
			if ok {
				break
			}
			if time.Now().After(deadline) {
				_ = db.Close()
				_ = testcontainer.TerminateContainer(container)
				dorisInitErr = fmt.Errorf("readiness: SELECT 1 never succeeded")
				return
			}
			time.Sleep(3 * time.Second)
		}

		// Doris reports "Current database is not set" before it even reaches
		// analysis, which would mask whether a probe parsed. Select a database.
		if _, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS doristest"); err != nil {
			_ = db.Close()
			_ = testcontainer.TerminateContainer(container)
			dorisInitErr = fmt.Errorf("create doristest db: %w", err)
			return
		}
		_ = db.Close()
		db, err = sql.Open("mysql", dsn("doristest"))
		if err != nil {
			_ = testcontainer.TerminateContainer(container)
			dorisInitErr = err
			return
		}
		// Container intentionally leaked for singleton reuse; process exit reaps it.
		dorisInst = &dorisContainer{db: db, ctx: ctx}
	})

	if dorisInitErr != nil {
		t.Fatalf("Doris container: %v", dorisInitErr)
	}
	return dorisInst
}

// dorisSyntaxErrMarkers are the message fragments Doris's Nereids parser emits
// for a genuine PARSE failure. CRITICAL: Doris codes parse AND semantic errors
// alike as MySQL 1105 / "errCode = 2", so the code alone does not mean
// rejected — only these ANTLR-generated fragments do.
var dorisSyntaxErrMarkers = []string{
	"mismatched input",
	"extraneous input",
	"no viable alternative",
	"Syntax error",
	// ANTLR's single-token-insertion recovery, which reports a parse failure
	// without any of the fragments above: "missing 'AS' at 'PARQUET'", and
	// "missing {'[', '{'} at ')'" when more than one token would serve. Both
	// quote forms are matched so the marker cannot hit a semantic message
	// that merely contains the word "missing".
	"missing '",
	"missing {",
}

// dorisMissingTokenMarker is the same recovery naming a token that has no
// literal spelling by its symbolic name instead: "missing EQ at '"v"'" for
// PROPERTIES("k" "v"). The fixed markers above need the quote or brace, so
// they took that parse failure for an accept. The pattern wants an upper-case
// token name followed by " at '", which no semantic message has the shape of.
var dorisMissingTokenMarker = regexp.MustCompile(`missing [A-Z][A-Z_]* at '`)

// accepts reports whether Doris accepted sql at the PARSE layer: executed OK,
// or failed with any non-syntax error (semantic/runtime) — both imply it
// parsed. A non-driver error (connection drop/timeout) is NOT a parse verdict —
// it is returned so the caller fails the test rather than silently passing an
// accept it never proved.
func (c *dorisContainer) accepts(sql string) (bool, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 15*time.Second)
	defer cancel()
	_, err := c.db.ExecContext(ctx, sql)
	if err == nil {
		return true, nil
	}
	myErr, ok := err.(*gomysql.MySQLError)
	if !ok {
		return false, fmt.Errorf("non-driver error (cannot determine parse verdict): %w", err)
	}
	return !isDorisSyntaxError(myErr.Message), nil
}

// isDorisSyntaxError reports whether msg is the engine's report of a parse
// failure rather than a semantic or runtime one.
func isDorisSyntaxError(msg string) bool {
	for _, m := range dorisSyntaxErrMarkers {
		if strings.Contains(msg, m) {
			return true
		}
	}
	return dorisMissingTokenMarker.MatchString(msg)
}

// TestDorisSyntaxErrorMarkers holds the oracle's classifier to messages the
// engine has returned, so a marker cannot drop out without a failure that
// needs no container to show it.
func TestDorisSyntaxErrorMarkers(t *testing.T) {
	cases := []struct {
		msg    string
		syntax bool
	}{
		{"errCode = 2, detailMessage = \nmismatched input ')' expecting {'[', '{', ...}(line 1, pos 51)\n", true},
		{"errCode = 2, detailMessage = \nmissing ')' at '<EOF>'(line 1, pos 34)\n", true},
		{"errCode = 2, detailMessage = \nmissing EQ at '\"v\"'(line 1, pos 92)\n", true},
		{"errCode = 2, detailMessage = \nno viable alternative at input 'ALTER JOB'(line 1, pos 6)\n", true},
		{"errCode = 2, detailMessage = \nextraneous input '+' expecting {'[', ...}(line 1, pos 93)\n", true},
		{"errCode = 2, detailMessage = Syntax error in line 1:\n...ope) WITH BROKER zzb (\"k\"=1)\n Encountered: INTEGER LITERAL\n", true},
		{"errCode = 2, detailMessage = Config 'zz_no_such_config' does not exist", false},
		{"errCode = 2, detailMessage = Unknown table 'zz_nope'", false},
		{"errCode = 2, detailMessage = replication num should be less than the number of available backends. replication num is 1, available backend num is 0", false},
	}
	for _, tc := range cases {
		if got := isDorisSyntaxError(tc.msg); got != tc.syntax {
			t.Errorf("isDorisSyntaxError(%q) = %v, want %v", tc.msg, got, tc.syntax)
		}
	}
}

// assertParity asserts the omni parser and Doris agree on sql (both accept or
// both reject) — the accept+reject conformance discipline.
func (c *dorisContainer) assertParity(t *testing.T, sql string, wantAccept bool) {
	t.Helper()
	_, errs := parseForTest(sql)
	omni := len(errs) == 0
	doris, err := c.accepts(sql)
	if err != nil {
		t.Fatalf("Doris oracle could not verdict %q: %v", sql, err)
	}
	if omni != doris {
		t.Errorf("MISMATCH %q: omni=%v doris=%v (omniErrs=%v)", sql, omni, doris, errs)
	}
	if omni != wantAccept {
		t.Errorf("omni %q: got accept=%v, want %v (errs=%v)", sql, omni, wantAccept, errs)
	}
}

// TestDorisConnectivity is the harness smoke test: the container is reachable.
func TestDorisConnectivity(t *testing.T) {
	c := startDoris(t)
	var one int
	if err := c.db.QueryRowContext(c.ctx, "SELECT 1").Scan(&one); err != nil || one != 1 {
		t.Fatalf("SELECT 1 = %d, err=%v", one, err)
	}
}

// TestDorisSyntaxConformance is the container-grounded conformance matrix for
// the expression and clause forms that the hand-written parser previously
// rejected even though Doris accepts them (BYT-10070 and its follow-ups), plus
// the sibling forms Doris itself rejects.
//
// The negative arm matters as much as the positive one: the Apache Doris ANTLR
// grammar accepts several MySQL-style special forms (SUBSTRING(x FROM n),
// TRIM(BOTH ... FROM ...), POSITION(x IN y), quantified comparisons) that the
// shipping engine rejects at parse time. Asserting those stay rejected keeps
// the parser aligned with the engine rather than with the grammar file.
func TestDorisSyntaxConformance(t *testing.T) {
	c := startDoris(t)

	// The probes reference table `t`; a missing table is a semantic error, which
	// still proves the statement parsed.
	cases := []struct {
		name       string
		sql        string
		wantAccept bool
	}{
		// EXTRACT — the originally reported gap.
		{"extract_year", "SELECT EXTRACT(YEAR FROM MONTHS_ADD(NOW(), -1))", true},
		{"extract_month", "SELECT EXTRACT(MONTH FROM NOW())", true},
		{"extract_compound_unit", "SELECT EXTRACT(DAY_HOUR FROM NOW())", true},
		{"extract_in_cte", "WITH c AS (SELECT EXTRACT(YEAR FROM NOW()) y) SELECT * FROM c", true},

		// Reserved keywords usable as function names (functionNameIdentifier).
		{"trim", "SELECT TRIM('  x  ')", true},
		{"trim_two_arg", "SELECT TRIM('xax', 'x')", true},
		{"left", "SELECT LEFT('abc', 2)", true},
		{"right", "SELECT RIGHT('abc', 2)", true},
		{"if", "SELECT IF(1, 2, 3)", true},
		{"database", "SELECT DATABASE()", true},
		{"schema", "SELECT SCHEMA()", true},
		{"user", "SELECT USER()", true},
		{"session_user", "SELECT SESSION_USER()", true},
		{"current_user", "SELECT CURRENT_USER()", true},
		{"current_catalog", "SELECT CURRENT_CATALOG()", true},
		{"connection_id", "SELECT CONNECTION_ID()", true},
		{"password", "SELECT PASSWORD('x')", true},

		// Nilary datetime functions with a fractional-seconds precision.
		{"current_timestamp_precision", "SELECT CURRENT_TIMESTAMP(3)", true},
		{"current_time_precision", "SELECT CURRENT_TIME(3)", true},
		{"localtime_precision", "SELECT LOCALTIME(3)", true},
		{"localtimestamp_precision", "SELECT LOCALTIMESTAMP(3)", true},
		{"current_date_parens", "SELECT CURRENT_DATE()", true},

		// Literals.
		{"hex_literal", "SELECT x'1f'", true},
		{"bit_literal", "SELECT b'101'", true},
		{"array_literal", "SELECT [1, 2, 3]", true},
		{"map_literal", "SELECT {'a': 1}", true},

		// Session and user variables.
		{"system_var", "SELECT @@version_comment", true},
		{"system_var_session_scope", "SELECT @@session.query_timeout", true},
		{"system_var_global_scope", "SELECT @@global.query_timeout", true},
		{"user_var", "SELECT @uservar", true},

		// Lambdas in higher-order array functions.
		{"lambda_array_map", "SELECT array_map(x -> x + 1, [1, 2, 3])", true},
		{"lambda_array_filter", "SELECT array_filter(x -> x > 1, [1, 2, 3])", true},
		{"lambda_multi_param", "SELECT array_map((x, y) -> x + y, [1, 2], [3, 4])", true},
		{"lambda_three_param", "SELECT array_map((x, y, z) -> x + y + z, [1], [2], [3])", true},
		// A single parameter must be written bare; the parenthesized form
		// requires at least two.
		{"lambda_single_param_parenthesized", "SELECT array_map((x) -> x + 1, [1, 2])", false},

		// Optimizer hints.
		{"set_var_hint", "SELECT /*+ SET_VAR(query_timeout = 100) */ 1", true},

		// Grouping elements. CUBE is the whole grouping specification: a
		// trailing list after it is a parse error, so accepting one here would
		// silently discard everything past the comma.
		{"group_by_cube", "SELECT a, SUM(b) FROM t GROUP BY CUBE(a)", true},
		{"group_by_cube_trailing_item", "SELECT a, b, SUM(c) FROM t GROUP BY CUBE(a), b", false},
		{"group_by_cube_trailing_cube", "SELECT a, SUM(b) FROM t GROUP BY CUBE(a), CUBE(b)", false},
		{"group_by_rollup", "SELECT a, SUM(b) FROM t GROUP BY ROLLUP(a)", true},
		{"group_by_grouping_sets", "SELECT a, SUM(b) FROM t GROUP BY GROUPING SETS ((a), ())", true},

		// CONVERT / USING charset.
		{"convert_using", "SELECT CONVERT('abc' USING utf8)", true},
		{"convert_type", "SELECT CONVERT('1', SIGNED)", true},
		{"char_using", "SELECT CHAR(65 USING utf8)", true},
		{"char_multi_arg_using", "SELECT CHAR(65, 66 USING utf8)", true},
		// The trailing USING clause belongs to CHAR and CONVERT only.
		{"using_on_aggregate", "SELECT SUM(a USING utf8) FROM t", false},
		{"using_on_scalar", "SELECT abs(1 USING utf8)", false},

		// BINARY operator — primary-level only: an identifier or string
		// literal (StarRocks differs and accepts the general prefix form).
		{"binary_operator", "SELECT BINARY 'abc'", true},
		{"binary_column", "SELECT BINARY a FROM t", true},
		{"binary_int_rejected", "SELECT BINARY 1", false},
		{"binary_paren_rejected", "SELECT BINARY (1)", false},
		{"binary_func_rejected", "SELECT BINARY now()", false},

		// Collection literal elements are constants: literals and nested
		// collection literals only (StarRocks arrays take full expressions).
		{"array_nested_literal", "SELECT [[1],[2]]", true},
		{"array_mixed_literals", "SELECT [NULL, 1]", true},
		{"array_signed_number_rejected", "SELECT [-1, +2]", false},
		{"array_expr_element_rejected", "SELECT [1+1]", false},
		{"array_column_element_rejected", "SELECT [a] FROM t", false},
		{"map_nested_literal", "SELECT {'a': [1,2]}", true},
		{"map_expr_value_rejected", "SELECT {'a': 1+1}", false},
		{"struct_literal", "SELECT {1, 2}", true},
		{"struct_single_element", "SELECT {'a'}", true},

		// String-form user variables.
		{"user_var_string", "SELECT @'quoted'", true},

		// Top-level parenthesized queries.
		{"paren_select_stmt", "(SELECT 1)", true},
		{"paren_select_nested", "((SELECT 1))", true},
		{"paren_select_union", "(SELECT 1) UNION (SELECT 2)", true},
		{"paren_setop_after_nested", "((SELECT 1) UNION SELECT 2)", true},
		{"paren_setop_after_with", "(WITH c AS (SELECT 1) SELECT 1 UNION SELECT 2)", true},
		{"paren_trailing_clauses", "(SELECT 1) ORDER BY 1 LIMIT 5", true},
		{"paren_nested_trailing_limit", "((SELECT 1)) LIMIT 5", true},
		{"paren_double_order_by", "(SELECT 1 ORDER BY 1) ORDER BY 1", true},
		{"paren_with_dml_rejected", "(WITH c AS (SELECT 1) DELETE FROM t)", false},
		{"union_paren_rhs_limit", "SELECT 1 UNION (SELECT 2) LIMIT 5", true},
		{"union_paren_rhs_order_limit", "SELECT 1 UNION (SELECT 2) ORDER BY 1 LIMIT 5", true},
		{"intersect_paren_rhs_limit", "SELECT 1 INTERSECT (SELECT 2) LIMIT 3", true},
		{"with_union", "WITH c AS (SELECT 1) SELECT 1 UNION SELECT 2", true},
		{"with_union_paren_rhs_limit", "WITH c AS (SELECT 1) SELECT 1 UNION (SELECT 2) LIMIT 5", true},
		{"cte_body_union", "WITH c AS (SELECT 1 UNION SELECT 2) SELECT * FROM c", true},
		{"paren_inner_trailing_limit", "(SELECT 1 UNION (SELECT 2) LIMIT 5)", true},
		// Unlike StarRocks, clause order and repetition are lenient here.
		{"select_limit_then_order", "SELECT 1 LIMIT 5 ORDER BY 1", false},
		{"union_limit_then_order", "SELECT 1 UNION SELECT 2 LIMIT 5 ORDER BY 1", false},
		{"select_double_order_by", "SELECT 1 ORDER BY 1 ORDER BY 2", false},
		{"paren_nested_double_order", "((SELECT 1 ORDER BY 1) ORDER BY 1)", true},
		{"order_by_subquery", "SELECT 1 ORDER BY (SELECT 1)", true},
		{"setop_order_by_subquery", "(SELECT 1) UNION (SELECT 2) ORDER BY (SELECT 1)", true},
		{"insert_union_paren_limit", "INSERT INTO t SELECT 1 UNION (SELECT 2) LIMIT 5", true},
		{"cte_body_paren", "WITH c AS ((SELECT 1)) SELECT * FROM c", true},
		{"cte_body_double_order", "WITH c AS (SELECT 1 ORDER BY 1 ORDER BY 2) SELECT * FROM c", false},
		{"paren_limit_then_outer_order", "(SELECT 1 LIMIT 1) ORDER BY 1", true},
		{"paren_scoped_cte_union", "(WITH c AS (SELECT 1) SELECT 1) UNION SELECT * FROM c", true},
		{"exists_grouped_subquery", "SELECT EXISTS (SELECT 1 FROM t ORDER BY 1 ORDER BY 2)", false},

		// LATERAL VIEW — previously hidden by the trailing-token swallow.
		{"lateral_view", "SELECT * FROM person LATERAL VIEW EXPLODE(ARRAY(30, 60)) tableName AS c_age", true},
		{"lateral_view_multi_col", "SELECT * FROM t LATERAL VIEW EXPLODE_SPLIT(s, ',') tmp AS c1, c2", true},
		{"lateral_view_chained", "SELECT * FROM t LATERAL VIEW EXPLODE([1]) a AS x LATERAL VIEW EXPLODE([2]) b AS y", true},
		{"lateral_view_after_alias", "SELECT * FROM t tt LATERAL VIEW EXPLODE([1]) tmp AS c", true},
		{"lateral_view_then_clauses", "SELECT * FROM t LATERAL VIEW EXPLODE([1]) tmp AS c WHERE c > 0 ORDER BY c LIMIT 1", true},
		{"lateral_view_outer_rejected", "SELECT * FROM t LATERAL VIEW OUTER EXPLODE([1]) tmp AS c", false},
		{"lateral_view_no_alias_rejected", "SELECT * FROM t LATERAL VIEW EXPLODE([1]) tmp", false},
		{"lateral_view_no_table_alias_rejected", "SELECT * FROM t LATERAL VIEW EXPLODE([1]) AS c", false},
		{"lateral_view_junk_rejected", "SELECT * FROM t LATERAL VIEW EXPLODE([1]) tmp AS c zzz", false},

		// --- Constructs previously hidden by the trailing-token swallow
		// (BYT-10084): each parsed as a valid prefix and silently dropped the
		// rest of the statement.
		{"limit_offset_comma", "SELECT * FROM t LIMIT 5, 10", true},
		{"string_alias", `SELECT (1 + 1) AS "20%"`, true},
		{"string_alias_single_quoted", "SELECT 1 AS 'x'", true},
		{"string_alias_empty", "SELECT 1 AS ''", true},
		{"element_at_literal", "SELECT [1, 2, 3][1]", true},
		{"element_at_map", "SELECT {'a': 1}['a']", true},
		{"element_at_column", "SELECT c1[1] FROM t", true},
		{"element_at_qualified", "SELECT t.c1[1] FROM t", true},
		{"qualified_column_arith", "SELECT t.a + 1 FROM t", true},
		{"array_slice", "SELECT [1, 2, 3][1:2]", true},
		{"array_slice_open", "SELECT [1, 2, 3][2:]", true},
		{"element_at_two_indexes", "SELECT [1, 2][1, 2]", false},
		{"array_slice_no_begin", "SELECT [1, 2][:2]", false},
		{"group_by_with_rollup", "SELECT a, SUM(b) FROM t GROUP BY a WITH ROLLUP", true},
		{"grouping_sets_trailing_item", "SELECT a FROM t GROUP BY GROUPING SETS ((a)), b", false},
		{"from_tvf_backends", "SELECT * FROM BACKENDS()", true},
		{"tablet_tablesample", "SELECT * FROM t TABLET(10001) TABLESAMPLE(1000 ROWS) REPEATABLE 2", true},
		{"show_databases_from", "SHOW DATABASES FROM internal", true},
		{"build_index_partition", "BUILD INDEX index1 ON table1 PARTITION(p1, p2)", true},

		// With the strict trailing-token check (BYT-10085), junk after a
		// valid prefix is rejected instead of silently dropped.
		{"json_arrow_rejected", "SELECT j->'$.a' FROM t", false},
		{"stray_comment_close_rejected", "SELECT 1 */ 2", false},
		{"explain_incomplete_rejected", "EXPLAIN SELECT * FROM", false},
		{"explain_bare_rejected", "EXPLAIN", false},
		{"set_expr_incomplete_rejected", "SET x = 1 +", false},
		{"set_valueless_rejected", "SET x", false},
		{"show_where_incomplete_rejected", "SHOW TABLES WHERE (", false},
		{"show_from_missing_db_rejected", "SHOW TABLES FROM", false},
		{"set_names_missing_charset_rejected", "SET NAMES", false},
		{"set_names_dangling_collate_rejected", "SET NAMES utf8 COLLATE", false},
		{"set_names", "SET NAMES utf8", true},
		{"set_names_default", "SET NAMES DEFAULT", true},
		// Unlike StarRocks, the engine rejects the SET ROLE family.
		{"set_role_rejected", "SET ROLE admin_role", false},
		{"set_names_collate", "SET NAMES utf8 COLLATE utf8_general_ci", true},
		// The pattern is required. This row once read as an accept, because
		// the oracle missed "missing STRING_LITERAL at '<EOF>'".
		{"show_like_no_pattern_rejected", "SHOW TABLES LIKE", false},
		{"show_databases_like_no_pattern_rejected", "SHOW DATABASES LIKE", false},
		{"from_table_named_selected", "SELECT * FROM selected", true},

		// --- Negative arm: the grammar file allows these, the engine does not.
		{"substring_from_for", "SELECT SUBSTRING('abcdef' FROM 2 FOR 3)", true},
		{"substring_from", "SELECT SUBSTRING('abcdef' FROM 2)", true},
		{"trim_both_from", "SELECT TRIM(BOTH 'x' FROM 'xax')", true},
		{"trim_leading_from", "SELECT TRIM(LEADING 'x' FROM 'xax')", true},
		{"trim_expr_from", "SELECT TRIM(' ' FROM ' a ')", true},
		{"position_in", "SELECT POSITION('b' IN 'abc')", true},
		{"position_comma_predicate", "SELECT POSITION('a' IN ('a'), 'b')", true},
		{"position_nested_predicate", "SELECT POSITION(('a' IN ('a')) IN 'abc')", true},
		{"trim_trailing_from", "SELECT TRIM(TRAILING 'x' FROM 'xaxx')", true},
		{"trim_implicit_space", "SELECT TRIM(BOTH FROM ' a ')", true},
		{"values_func", "SELECT VALUES(a) FROM t", false},
		{"quantified_any", "SELECT * FROM t WHERE a = ANY (SELECT b FROM t)", false},
		{"quantified_all", "SELECT * FROM t WHERE a > ALL (SELECT b FROM t)", false},
		{"quantified_some", "SELECT * FROM t WHERE a = SOME (SELECT b FROM t)", false},

		// --- INTO OUTFILE (BOT-92). The clause is statement-level in the
		// grammar — `statementBase : explain? query outFileClause?` — so it
		// follows every clause a query carries and cannot appear inside one.
		{"outfile", "SELECT 1 INTO OUTFILE '/tmp/x'", true},
		{"outfile_double_quoted_path", `SELECT * FROM t INTO OUTFILE "s3://b/pre"`, true},
		{"outfile_format_parquet", `SELECT * FROM t INTO OUTFILE "s3://b/pre" FORMAT AS PARQUET`, true},
		{"outfile_format_csv", `SELECT * FROM t INTO OUTFILE "s3://b/pre" FORMAT AS CSV`, true},
		{"outfile_format_orc", `SELECT * FROM t INTO OUTFILE "s3://b/pre" FORMAT AS ORC`, true},
		{"outfile_format_non_keyword", `SELECT * FROM t INTO OUTFILE "s3://b/pre" FORMAT AS csv_with_names`, true},
		{"outfile_format_quoted_ident", "SELECT 1 INTO OUTFILE '/tmp/x' FORMAT AS `PARQUET`", true},
		// An unknown format name is a semantic failure, not a parse one.
		{"outfile_format_unknown_name", "SELECT 1 INTO OUTFILE '/tmp/x' FORMAT AS no_such_format", true},
		{"outfile_properties", `SELECT * FROM t INTO OUTFILE "s3://b/pre" PROPERTIES("k" = "v")`, true},
		{"outfile_format_and_properties", `SELECT * FROM t INTO OUTFILE "s3://b/pre" FORMAT AS CSV PROPERTIES("column_separator" = ",")`, true},
		{"outfile_properties_multiple", `SELECT * FROM t INTO OUTFILE "s3://b/pre" PROPERTIES("k" = "v", "k2" = "v2")`, true},
		{"outfile_after_order_by_limit", "SELECT * FROM t ORDER BY a LIMIT 10 INTO OUTFILE '/tmp/x'", true},
		{"outfile_after_mysql_limit", "SELECT * FROM t LIMIT 5,10 INTO OUTFILE '/tmp/x'", true},
		{"outfile_full_query", "SELECT a FROM t WHERE a = 1 GROUP BY a HAVING a > 0 ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/x'", true},
		{"outfile_on_set_op", "SELECT a FROM t UNION SELECT a FROM t INTO OUTFILE '/tmp/x'", true},
		{"outfile_after_parens", "(SELECT * FROM t) INTO OUTFILE '/tmp/x'", true},
		{"outfile_after_parens_limit", "(SELECT * FROM t) LIMIT 10 INTO OUTFILE '/tmp/x'", true},
		{"outfile_after_repeated_order_by", "(SELECT 1 ORDER BY 1) ORDER BY 1 INTO OUTFILE '/tmp/x'", true},
		{"outfile_with_cte", "WITH c AS (SELECT 1 AS a) SELECT * FROM c INTO OUTFILE '/tmp/x'", true},
		// `explain? query outFileClause?`: EXPLAIN takes the clause too.
		{"outfile_explain", "EXPLAIN SELECT 1 INTO OUTFILE '/tmp/x'", true},

		// Nothing may follow the clause: it is the last thing in a statement.
		{"outfile_then_order_by_rejected", "SELECT 1 INTO OUTFILE '/tmp/x' ORDER BY 1", false},
		{"outfile_then_limit_rejected", "SELECT 1 INTO OUTFILE '/tmp/x' LIMIT 1", false},
		{"outfile_then_union_rejected", "SELECT 1 INTO OUTFILE '/tmp/x' UNION SELECT 2", false},
		{"outfile_twice_rejected", "SELECT 1 INTO OUTFILE '/tmp/x' INTO OUTFILE '/tmp/y'", false},
		{"outfile_before_union_rejected", "SELECT a FROM t INTO OUTFILE '/tmp/x' UNION SELECT a FROM t", false},
		// And it cannot appear inside a query: not in parens, a derived
		// table, or a parenthesized set-operation arm.
		{"outfile_inside_parens_rejected", "(SELECT 1 INTO OUTFILE '/tmp/x')", false},
		{"outfile_inside_derived_table_rejected", "SELECT * FROM (SELECT a FROM t INTO OUTFILE '/tmp/x') x", false},
		{"outfile_inside_set_op_arm_rejected", "SELECT 1 UNION (SELECT a FROM t INTO OUTFILE '/tmp/x')", false},
		{"outfile_on_insert_rejected", "INSERT INTO t SELECT 1 INTO OUTFILE '/tmp/x'", false},
		// A job's DO body is supportedDmlStatement, so it cannot carry the
		// clause, nor be a query at all.
		{"job_do_insert", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO INSERT INTO zz_nope SELECT 1", true},
		{"job_do_insert_values", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO INSERT INTO zz_nope VALUES (1)", true},
		{"job_do_insert_overwrite", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO INSERT OVERWRITE TABLE zz_nope SELECT 1", true},
		{"job_do_update", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO UPDATE zz_nope SET a = 1 WHERE a = 2", true},
		{"job_do_delete", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO DELETE FROM zz_nope WHERE a = 1", true},
		{"job_do_explain_insert", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO EXPLAIN INSERT INTO zz_nope SELECT 1", true},
		{"job_do_explain_update", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO EXPLAIN UPDATE zz_nope SET a = 1", true},
		{"job_do_load", `CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO LOAD LABEL zzl (DATA INFILE("hdfs://h/p") INTO TABLE zz_nope)`, true},
		{"job_do_export", `CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO EXPORT TABLE zz_nope TO "hdfs://h/p"`, true},
		{"job_do_select_rejected", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO SELECT 1", false},
		{"job_do_select_outfile_rejected", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO SELECT 1 INTO OUTFILE '/tmp/x'", false},
		{"job_do_select_union_rejected", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO SELECT 1 UNION SELECT 2", false},
		{"job_do_paren_select_rejected", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO (SELECT 1)", false},
		{"job_do_with_select_rejected", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO WITH c AS (SELECT 1) SELECT * FROM c", false},
		{"job_do_explain_select_rejected", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO EXPLAIN SELECT 1", false},
		{"job_do_create_table_rejected", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO CREATE TABLE zzx (a INT)", false},
		{"job_do_show_rejected", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO SHOW TABLES", false},
		{"job_do_set_rejected", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO SET x = 1", false},
		{"job_do_truncate", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO TRUNCATE TABLE zz_nope", true},
		{"job_do_merge", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO MERGE INTO zz_nope USING zz_nope2 ON zz_nope.a = zz_nope2.a WHEN MATCHED THEN DELETE", true},
		{"job_do_create_job_rejected", "CREATE JOB zzjob ON SCHEDULE AT '2030-01-01 00:00:00' DO CREATE JOB zzjob2 ON SCHEDULE AT '2030-01-01 00:00:00' DO INSERT INTO zz_nope SELECT 1", false},
		// ALTER JOB gets no rows: the engine has no ALTER JOB statement, so a
		// row would pin the form, not its body. Its DO body is held to the
		// same rule by TestJobDoBodyIsDML.
		// The query of CREATE TABLE ... AS is not a subquery, but cannot
		// carry the clause either.
		{"ctas_outfile_rejected", "CREATE TABLE zzctas1 AS SELECT 1 INTO OUTFILE '/tmp/x'", false},
		{"ctas_set_op_outfile_rejected", "CREATE TABLE zzctas2 AS SELECT 1 UNION SELECT 2 INTO OUTFILE '/tmp/x'", false},
		// Clause shape.
		{"outfile_no_path_rejected", "SELECT 1 INTO OUTFILE", false},
		{"outfile_without_outfile_kw_rejected", "SELECT 1 INTO '/tmp/x'", false},
		{"outfile_without_into_kw_rejected", "SELECT 1 OUTFILE '/tmp/x'", false},
		{"outfile_format_without_as_rejected", "SELECT 1 INTO OUTFILE '/tmp/x' FORMAT PARQUET", false},
		{"outfile_format_as_text_rejected", `SELECT 1 INTO OUTFILE '/tmp/x' FORMAT AS "PARQUET"`, false},
		{"outfile_format_as_reserved_rejected", "SELECT 1 INTO OUTFILE '/tmp/x' FORMAT AS select", false},
		{"outfile_format_as_nothing_rejected", "SELECT 1 INTO OUTFILE '/tmp/x' FORMAT AS", false},
		{"outfile_properties_bare_rejected", "SELECT 1 INTO OUTFILE '/tmp/x' PROPERTIES", false},
		{"outfile_properties_empty_rejected", "SELECT 1 INTO OUTFILE '/tmp/x' PROPERTIES()", false},
		// propertyItemList is propertyItem (COMMA propertyItem)*, so the items
		// are comma-separated and the list does not end on one.
		{"outfile_properties_missing_comma_rejected", `SELECT 1 INTO OUTFILE '/tmp/x' PROPERTIES("k"="v" "k2"="v2")`, false},
		{"outfile_properties_missing_commas_rejected", `SELECT 1 INTO OUTFILE '/tmp/x' PROPERTIES("k"="v" "k2"="v2" "k3"="v3")`, false},
		{"outfile_properties_trailing_comma_rejected", `SELECT 1 INTO OUTFILE '/tmp/x' PROPERTIES("k"="v",)`, false},
		{"outfile_properties_leading_comma_rejected", `SELECT 1 INTO OUTFILE '/tmp/x' PROPERTIES(,"k"="v")`, false},
		{"outfile_properties_double_comma_rejected", `SELECT 1 INTO OUTFILE '/tmp/x' PROPERTIES("k"="v",,"k2"="v2")`, false},
		{"outfile_trailing_garbage_rejected", "SELECT 1 INTO OUTFILE '/tmp/x' GARBAGE", false},
		// --- Property lists. The grammar gives every one of them as
		// propertyItemList: propertyItem (COMMA propertyItem)* -- so at least
		// one item, comma-separated, not ending on a comma. Each clause below
		// reaches a different one of the parsers that used to carry its own
		// copy of the list loop with the comma optional, so each accepted all
		// three departures the engine refuses.
		{"props_create_table", `CREATE TABLE zzp1 (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("k"="v","k2"="v2")`, true},
		{"props_create_table_missing_comma_rejected", `CREATE TABLE zzp2 (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("k"="v" "k2"="v2")`, false},
		{"props_create_table_trailing_comma_rejected", `CREATE TABLE zzp3 (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("k"="v",)`, false},
		{"props_create_table_empty_rejected", `CREATE TABLE zzp4 (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES()`, false},

		{"props_create_database", `CREATE DATABASE zzpdb PROPERTIES("k"="v","k2"="v2")`, true},
		{"props_create_database_missing_comma_rejected", `CREATE DATABASE zzpdb2 PROPERTIES("k"="v" "k2"="v2")`, false},
		{"props_create_database_trailing_comma_rejected", `CREATE DATABASE zzpdb3 PROPERTIES("k"="v",)`, false},
		{"props_create_database_empty_rejected", `CREATE DATABASE zzpdb4 PROPERTIES()`, false},

		{"props_alter_table_set", `ALTER TABLE t SET ("k"="v","k2"="v2")`, true},
		{"props_alter_table_set_missing_comma_rejected", `ALTER TABLE t SET ("k"="v" "k2"="v2")`, false},
		{"props_alter_table_set_trailing_comma_rejected", `ALTER TABLE t SET ("k"="v",)`, false},
		{"props_alter_table_set_empty_rejected", `ALTER TABLE t SET ()`, false},

		{"props_admin_frontend_config", `ADMIN SET FRONTEND CONFIG ("k"="v","k2"="v2")`, true},
		{"props_admin_frontend_config_missing_comma_rejected", `ADMIN SET FRONTEND CONFIG ("k"="v" "k2"="v2")`, false},
		{"props_admin_frontend_config_trailing_comma_rejected", `ADMIN SET FRONTEND CONFIG ("k"="v",)`, false},
		{"props_admin_frontend_config_empty_rejected", `ADMIN SET FRONTEND CONFIG ()`, false},

		// ALTER CATALOG ... PROPERTY (...) reaches the sixth list parser, but
		// it gets no rows here: the engine rejects that statement form
		// outright, well-formed list and all, so a parity row would pin the
		// form rather than the list. omni accepting a form the engine does not
		// have is a separate, pre-existing divergence; the list shape for that
		// parser is covered by TestPropertyListSeparators.

		{"props_create_mv", `CREATE MATERIALIZED VIEW zzpmv BUILD IMMEDIATE REFRESH COMPLETE DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("k"="v","k2"="v2") AS SELECT 1 AS a`, true},
		{"props_create_mv_missing_comma_rejected", `CREATE MATERIALIZED VIEW zzpmv2 BUILD IMMEDIATE REFRESH COMPLETE DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("k"="v" "k2"="v2") AS SELECT 1 AS a`, false},
		{"props_create_mv_trailing_comma_rejected", `CREATE MATERIALIZED VIEW zzpmv3 BUILD IMMEDIATE REFRESH COMPLETE DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("k"="v",) AS SELECT 1 AS a`, false},
		{"props_create_mv_empty_rejected", `CREATE MATERIALIZED VIEW zzpmv4 BUILD IMMEDIATE REFRESH COMPLETE DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES() AS SELECT 1 AS a`, false},

		{"props_routine_load_job", `CREATE ROUTINE LOAD zzrl ON t PROPERTIES("k"="v","k2"="v2") FROM KAFKA("kafka_broker_list"="localhost:9092","kafka_topic"="tp")`, true},
		{"props_routine_load_job_missing_comma_rejected", `CREATE ROUTINE LOAD zzrl2 ON t PROPERTIES("k"="v" "k2"="v2") FROM KAFKA("kafka_broker_list"="localhost:9092","kafka_topic"="tp")`, false},
		{"props_routine_load_job_trailing_comma_rejected", `CREATE ROUTINE LOAD zzrl3 ON t PROPERTIES("k"="v",) FROM KAFKA("kafka_broker_list"="localhost:9092","kafka_topic"="tp")`, false},
		{"props_routine_load_source_missing_comma_rejected", `CREATE ROUTINE LOAD zzrl4 ON t PROPERTIES("desired_concurrent_number"="3") FROM KAFKA("k"="v" "k2"="v2")`, false},
		{"props_routine_load_source_empty_rejected", `CREATE ROUTINE LOAD zzrl5 ON t PROPERTIES("desired_concurrent_number"="3") FROM KAFKA()`, false},

		// --- Property lists whose errors were dropped or that were skipped
		// unread, the follow-up to the props_* rows above. At the end of the
		// input nothing after a dropped error was left to notice it.
		{"props_admin_frontend_config_open_rejected", `ADMIN SET FRONTEND CONFIG (`, false},
		{"props_admin_frontend_config_no_value_rejected", `ADMIN SET FRONTEND CONFIG ("k"=`, false},
		{"props_admin_frontend_config_unclosed_rejected", `ADMIN SET FRONTEND CONFIG ("k"="v"`, false},
		// The list is optional in the grammar; the engine then wants one
		// config at execution.
		{"props_admin_frontend_config_bare", `ADMIN SET FRONTEND CONFIG`, true},
		{"props_modify_backend", `ALTER SYSTEM MODIFY BACKEND 'nonexistent-host-zz:9050' SET ("tag.location"="a")`, true},
		{"props_modify_backend_open_rejected", `ALTER SYSTEM MODIFY BACKEND 'nonexistent-host-zz:9050' SET (`, false},
		{"props_modify_backend_no_value_rejected", `ALTER SYSTEM MODIFY BACKEND 'nonexistent-host-zz:9050' SET ("k"=`, false},
		{"props_admin_set_replica_status", `ADMIN SET REPLICA STATUS PROPERTIES("tablet_id"="1","backend_id"="1","status"="bad")`, true},
		{"props_admin_set_replica_status_open_rejected", `ADMIN SET REPLICA STATUS PROPERTIES(`, false},
		{"props_admin_set_table_status", `ADMIN SET TABLE zz_nope STATUS PROPERTIES("state"="NORMAL")`, true},
		{"props_admin_set_table_status_open_rejected", `ADMIN SET TABLE zz_nope STATUS PROPERTIES(`, false},
		{"props_admin_set_table_status_no_value_rejected", `ADMIN SET TABLE zz_nope STATUS PROPERTIES("k"=`, false},
		// CHECK TABLET and COPY TABLET took their PROPERTIES list as raw args.
		{"props_admin_check_tablet", `ADMIN CHECK TABLET (1) PROPERTIES("type"="consistency")`, true},
		{"props_admin_check_tablet_open_rejected", `ADMIN CHECK TABLET (1) PROPERTIES(`, false},
		{"props_admin_check_tablet_empty_rejected", `ADMIN CHECK TABLET (1) PROPERTIES()`, false},
		{"props_admin_check_tablet_trailing_comma_rejected", `ADMIN CHECK TABLET (1) PROPERTIES("type"="consistency",)`, false},
		{"props_admin_copy_tablet", `ADMIN COPY TABLET 1 PROPERTIES("backend_id"="1")`, true},
		{"props_admin_copy_tablet_open_rejected", `ADMIN COPY TABLET 1 PROPERTIES(`, false},
		{"props_admin_copy_tablet_trailing_comma_rejected", `ADMIN COPY TABLET 1 PROPERTIES("backend_id"="1",)`, false},
		{"props_load_errors_hub", `ALTER SYSTEM SET LOAD ERRORS HUB PROPERTIES("type"="null")`, true},
		{"props_load_errors_hub_open_rejected", `ALTER SYSTEM SET LOAD ERRORS HUB PROPERTIES(`, false},
		// No accept row for ADD BACKEND: the engine would register the
		// backend in the container every other row shares.
		{"props_add_backend_open_rejected", `ALTER SYSTEM ADD BACKEND 'nonexistent-host-zz:9050' PROPERTIES(`, false},
		{"alter_system_add_broker_no_name_rejected", `ALTER SYSTEM ADD BROKER`, false},
		{"alter_system_drop_broker_no_name_rejected", `ALTER SYSTEM DROP BROKER`, false},
		{"alter_system_drop_all_broker_no_name_rejected", `ALTER SYSTEM DROP ALL BROKER`, false},
		{"alter_system_drop_all_broker", `ALTER SYSTEM DROP ALL BROKER zz_no_broker`, true},
		{"show_alter_table_from_db", `SHOW ALTER TABLE COLUMN FROM doristest`, true},
		{"show_alter_table_from_missing_db_rejected", `SHOW ALTER TABLE COLUMN FROM`, false},
		{"show_alter_table_from_string_rejected", `SHOW ALTER TABLE COLUMN FROM 'x'`, false},

		// A partition's inline properties: partitionDef's (LEFT_PAREN
		// propertyItemList RIGHT_PAREN)?, with propertyKey and propertyValue
		// each `identifier | constant`.
		{"props_partition", `CREATE TABLE zzpp01 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("replication_num"="1")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_empty_rejected", `CREATE TABLE zzpp02 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ()) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		{"props_partition_trailing_comma_rejected", `CREATE TABLE zzpp03 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("replication_num"="1",)) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		{"props_partition_missing_comma_rejected", `CREATE TABLE zzpp04 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("replication_num"="1" "k2"="v2")) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		{"props_partition_junk_rejected", `CREATE TABLE zzpp05 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k" "v" SELECT DROP *)) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		{"props_partition_unclosed_rejected", `CREATE TABLE zzpp06 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") (`, false},
		{"props_partition_bare_key", `CREATE TABLE zzpp07 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") (replication_num="1")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_keyword_key", `CREATE TABLE zzpp08 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") (comment="v")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_number_key", `CREATE TABLE zzpp09 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") (1="v")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_number_value", `CREATE TABLE zzpp10 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("replication_num"=1)) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_negative_value", `CREATE TABLE zzpp11 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k"=-1)) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_decimal_value", `CREATE TABLE zzpp12 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k"=1.5)) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_exponent_value", `CREATE TABLE zzpp13 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k"=1e3)) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_null_value", `CREATE TABLE zzpp14 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k"=NULL)) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_boolean_value", `CREATE TABLE zzpp15 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k"=true)) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_identifier_value", `CREATE TABLE zzpp16 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("storage_medium"=HDD)) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_quoted_identifier_value", "CREATE TABLE zzpp17 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN (\"10\") (\"k\"=`v`)) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES(\"replication_num\"=\"1\")", true},
		{"props_partition_array_value", `CREATE TABLE zzpp18 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k"=[1])) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_typed_value", `CREATE TABLE zzpp19 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k"=DATE "2020-01-01")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_reserved_value_rejected", `CREATE TABLE zzpp20 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k"=SELECT)) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		{"props_partition_dotted_value_rejected", `CREATE TABLE zzpp21 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k"=a.b)) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		{"props_partition_expression_value_rejected", `CREATE TABLE zzpp22 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k"=1+1)) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		{"props_partition_plus_value_rejected", `CREATE TABLE zzpp23 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k"=+1)) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		{"props_partition_missing_value_rejected", `CREATE TABLE zzpp24 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k"=)) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		{"props_partition_maxvalue", `CREATE TABLE zzpp25 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN MAXVALUE ("replication_num"="1")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_fixed_range", `CREATE TABLE zzpp26 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES [("1"), ("10")) ("replication_num"="1")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_list", `CREATE TABLE zzpp27 (a INT) PARTITION BY LIST(a) (PARTITION p1 VALUES IN ("1") ("replication_num"="1")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"props_partition_twice_rejected", `CREATE TABLE zzpp28 (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("replication_num"="1") ("k"="v")) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		// inPartitionDef makes VALUES IN optional, so a list partition's value
		// list may follow the name directly, then its property list.
		{"partition_list_without_values", `CREATE TABLE zzpp29 (a INT) PARTITION BY LIST(a) (PARTITION p1 ("1")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"partition_list_without_values_nested", `CREATE TABLE zzpp30 (a INT) PARTITION BY LIST(a) (PARTITION p1 (("1"), ("2"))) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"partition_list_without_values_props", `CREATE TABLE zzpp31 (a INT) PARTITION BY LIST(a) (PARTITION p1 ("replication_num"="1")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"partition_list_without_values_then_props", `CREATE TABLE zzpp32 (a INT) PARTITION BY LIST(a) (PARTITION p1 ("1") ("replication_num"="1")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"partition_list_without_values_empty_rejected", `CREATE TABLE zzpp33 (a INT) PARTITION BY LIST(a) (PARTITION p1 ()) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		{"partition_list_without_values_junk_rejected", `CREATE TABLE zzpp34 (a INT) PARTITION BY LIST(a) (PARTITION p1 (SELECT)) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		{"partition_list_without_values_empty_props_rejected", `CREATE TABLE zzpp35 (a INT) PARTITION BY LIST(a) (PARTITION p1 ("1") ()) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		// A property key may take more than one token, so the '(' after a
		// bare partition name is classified only once the whole key is read.
		{"partition_list_without_values_negative_key_props", `CREATE TABLE zzpp36 (a INT) PARTITION BY LIST(a) (PARTITION p1 (-1="v")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"partition_list_without_values_typed_key_props", `CREATE TABLE zzpp37 (a INT) PARTITION BY LIST(a) (PARTITION p1 (DATE "2020-01-01"="v")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"partition_list_without_values_array_key_props", `CREATE TABLE zzpp38 (a INT) PARTITION BY LIST(a) (PARTITION p1 ([1]="v")) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},
		{"partition_list_without_values_negative_value", `CREATE TABLE zzpp39 (a INT) PARTITION BY LIST(a) (PARTITION p1 (-1)) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("replication_num"="1")`, true},

		// WITH BROKER name (...): the same propertyItemList, but a broker load
		// takes a string on both sides of each item where EXPORT takes the
		// full item form; both lists used to be skipped unread.
		{"props_load_broker", `LOAD LABEL zzl01 (DATA INFILE("hdfs://h/p") INTO TABLE zz_nope) WITH BROKER zzb ("username"="u","password"="p")`, true},
		{"props_load_broker_single_quoted", `LOAD LABEL zzl02 (DATA INFILE("hdfs://h/p") INTO TABLE zz_nope) WITH BROKER zzb ('k'='v')`, true},
		{"props_load_broker_empty_rejected", `LOAD LABEL zzl03 (DATA INFILE("hdfs://h/p") INTO TABLE zz_nope) WITH BROKER zzb ()`, false},
		{"props_load_broker_trailing_comma_rejected", `LOAD LABEL zzl04 (DATA INFILE("hdfs://h/p") INTO TABLE zz_nope) WITH BROKER zzb ("username"="u",)`, false},
		{"props_load_broker_missing_comma_rejected", `LOAD LABEL zzl05 (DATA INFILE("hdfs://h/p") INTO TABLE zz_nope) WITH BROKER zzb ("username"="u" "password"="p")`, false},
		{"props_load_broker_junk_rejected", `LOAD LABEL zzl06 (DATA INFILE("hdfs://h/p") INTO TABLE zz_nope) WITH BROKER zzb ("k" "v" SELECT DROP *)`, false},
		{"props_load_broker_unclosed_rejected", `LOAD LABEL zzl07 (DATA INFILE("hdfs://h/p") INTO TABLE zz_nope) WITH BROKER zzb (`, false},
		{"props_load_broker_bare_key", `LOAD LABEL zzl08 (DATA INFILE("hdfs://h/p") INTO TABLE zz_nope) WITH BROKER zzb (username="u")`, true},
		{"props_load_broker_number_value", `LOAD LABEL zzl09 (DATA INFILE("hdfs://h/p") INTO TABLE zz_nope) WITH BROKER zzb ("k"=1)`, true},
		{"props_load_broker_identifier_value", `LOAD LABEL zzl10 (DATA INFILE("hdfs://h/p") INTO TABLE zz_nope) WITH BROKER zzb ("k"=v)`, true},
		{"props_export_broker", `EXPORT TABLE zz_nope TO "hdfs://h/p" WITH BROKER zzb ("username"="u","password"="p")`, true},
		{"props_export_broker_empty_rejected", `EXPORT TABLE zz_nope TO "hdfs://h/p" WITH BROKER zzb ()`, false},
		{"props_export_broker_trailing_comma_rejected", `EXPORT TABLE zz_nope TO "hdfs://h/p" WITH BROKER zzb ("username"="u",)`, false},
		{"props_export_broker_missing_comma_rejected", `EXPORT TABLE zz_nope TO "hdfs://h/p" WITH BROKER zzb ("username"="u" "password"="p")`, false},
		{"props_export_broker_junk_rejected", `EXPORT TABLE zz_nope TO "hdfs://h/p" WITH BROKER zzb ("k" "v" SELECT DROP *)`, false},
		{"props_export_broker_unclosed_rejected", `EXPORT TABLE zz_nope TO "hdfs://h/p" WITH BROKER zzb (`, false},
		{"props_export_broker_bare_key", `EXPORT TABLE zz_nope TO "hdfs://h/p" WITH BROKER zzb (username="u")`, true},
		{"props_export_broker_number_value", `EXPORT TABLE zz_nope TO "hdfs://h/p" WITH BROKER zzb ("k"=1)`, true},
		{"props_export_broker_identifier_value", `EXPORT TABLE zz_nope TO "hdfs://h/p" WITH BROKER zzb ("k"=v)`, true},

		// A routine load's data-source items: a string on both sides. The
		// value used to be any single token.
		{"props_routine_load_source", `CREATE ROUTINE LOAD zzrl6 ON zz_nope FROM KAFKA ('kafka_broker_list'='h:9092', "kafka_topic"="tp")`, true},
		{"props_routine_load_source_comma_value_rejected", `CREATE ROUTINE LOAD zzrl7 ON zz_nope FROM KAFKA ("a" = ,)`, false},
		{"props_routine_load_source_paren_value_rejected", `CREATE ROUTINE LOAD zzrl8 ON zz_nope FROM KAFKA ("a" = ))`, false},
		{"props_routine_load_source_open_paren_value_rejected", `CREATE ROUTINE LOAD zzrl9 ON zz_nope FROM KAFKA ("a" = ( )`, false},
		{"props_routine_load_source_bare_key", `CREATE ROUTINE LOAD zzrl10 ON zz_nope FROM KAFKA (kafka_broker_list="h:9092")`, true},
		{"props_routine_load_source_keyword_key", `CREATE ROUTINE LOAD zzrl11 ON zz_nope FROM KAFKA ("kafka_broker_list"="h:9092", comment="v")`, true},
		{"props_routine_load_source_number_value", `CREATE ROUTINE LOAD zzrl12 ON zz_nope FROM KAFKA ("kafka_broker_list"="h:9092", "kafka_partitions"=1)`, true},
		{"props_routine_load_source_identifier_value", `CREATE ROUTINE LOAD zzrl13 ON zz_nope FROM KAFKA ("kafka_broker_list"="h:9092", "kafka_topic"=tp)`, true},
		{"props_routine_load_source_null_value", `CREATE ROUTINE LOAD zzrl14 ON zz_nope FROM KAFKA ("kafka_broker_list"="h:9092", "k"=NULL)`, true},
		{"props_routine_load_source_boolean_value", `CREATE ROUTINE LOAD zzrl15 ON zz_nope FROM KAFKA ("kafka_broker_list"="h:9092", "k"=true)`, true},
		{"props_alter_routine_load_source", `ALTER ROUTINE LOAD FOR zzrl16 FROM KAFKA ("kafka_partitions"="1")`, true},
		{"props_alter_routine_load_source_comma_value_rejected", `ALTER ROUTINE LOAD FOR zzrl17 FROM KAFKA ("a" = ,)`, false},
		{"props_alter_routine_load_source_number_value", `ALTER ROUTINE LOAD FOR zzrl18 FROM KAFKA ("kafka_partitions"=1)`, true},

		// MySQL spellings the engine does not share.
		{"outfile_into_var_rejected", "SELECT 1 INTO @v", false},
		{"outfile_dumpfile_rejected", "SELECT 1 INTO DUMPFILE '/tmp/x'", false},
		// Three forms the engine accepts and this parser refuses are left out
		// of the matrix rather than asserted, since a parity row would have to
		// claim agreement that does not exist. Each fails closed:
		//   INTO OUTFILE 1                 — the grammar's path is `constant`,
		//                                    and only a string names a path
		//                                    (see parseIntoOutfileClause).
		//   PROPERTIES(k = "v")            — propertyKey is
		//                                    `identifier | constant`, and the
		//   PROPERTIES("k" = 1)              shared parseProperties takes a
		//                                    string on both sides at every
		//                                    PROPERTIES in the tree.
		//
		// The separator and emptiness rules above are the other direction —
		// the shared helper is loose where the engine is not — and this clause
		// takes the strict list so they are asserted rather than excused.
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c.assertParity(t, tc.sql, tc.wantAccept)
		})
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	if dorisInst != nil {
		_ = dorisInst.db.Close()
	}
	if dorisServer != nil {
		_ = testcontainer.TerminateContainer(dorisServer)
	}
	os.Exit(code)
}
