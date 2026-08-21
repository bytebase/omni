package parser

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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
	dorisOnce    sync.Once
	dorisInst    *dorisContainer
	dorisInitErr error
)

// startDoris starts (once) the Doris container and returns a ready
// dorisContainer. Skipped in -short mode — omni CI runs `go test -short ./...`,
// so container tests must be Short-gated (mirrors starrocks_container_test.go).
func startDoris(t *testing.T) *dorisContainer {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping Doris container test in short mode")
	}

	dorisOnce.Do(func() {
		ctx := context.Background()
		req := testcontainers.ContainerRequest{
			Image:        "dyrnq/doris:3.1.4",
			ExposedPorts: []string{"9030/tcp"},
			Env: map[string]string{
				"RUN_MODE": "STANDALONE",
				// bin/start_be.sh refuses to boot unless vm.max_map_count >=
				// 2000000 and swap is off. Neither is settable from inside a
				// container (vm.max_map_count is not namespaced), so use the
				// script's own escape hatch rather than patching it.
				"SKIP_CHECK_ULIMIT": "true",
			},
			WaitingFor: wait.ForListeningPort("9030/tcp").WithStartupTimeout(5 * time.Minute),
		}
		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			dorisInitErr = fmt.Errorf("start Doris container: %w", err)
			return
		}
		host, err := container.Host(ctx)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			dorisInitErr = err
			return
		}
		port, err := container.MappedPort(ctx, "9030")
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			dorisInitErr = err
			return
		}
		dsn := func(db string) string {
			return fmt.Sprintf("root@tcp(%s:%s)/%s?multiStatements=true&timeout=10s", host, port.Port(), db)
		}
		db, err := sql.Open("mysql", dsn(""))
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
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
				_ = testcontainers.TerminateContainer(container)
				dorisInitErr = fmt.Errorf("readiness: SELECT 1 never succeeded")
				return
			}
			time.Sleep(3 * time.Second)
		}

		// Doris reports "Current database is not set" before it even reaches
		// analysis, which would mask whether a probe parsed. Select a database.
		if _, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS doristest"); err != nil {
			_ = db.Close()
			_ = testcontainers.TerminateContainer(container)
			dorisInitErr = fmt.Errorf("create doristest db: %w", err)
			return
		}
		_ = db.Close()
		db, err = sql.Open("mysql", dsn("doristest"))
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
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
}

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
	for _, m := range dorisSyntaxErrMarkers {
		if strings.Contains(myErr.Message, m) {
			return false, nil
		}
	}
	return true, nil
}

// assertParity asserts the omni parser and Doris agree on sql (both accept or
// both reject) — the accept+reject conformance discipline.
func (c *dorisContainer) assertParity(t *testing.T, sql string, wantAccept bool) {
	t.Helper()
	_, errs := Parse(sql)
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

		// BINARY operator.
		{"binary_operator", "SELECT BINARY 'abc'", true},

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
		{"from_table_named_selected", "SELECT * FROM selected", true},

		// --- Negative arm: the grammar file allows these, the engine does not.
		{"substring_from_for", "SELECT SUBSTRING('abcdef' FROM 2 FOR 3)", false},
		{"substring_from", "SELECT SUBSTRING('abcdef' FROM 2)", false},
		{"trim_both_from", "SELECT TRIM(BOTH 'x' FROM 'xax')", false},
		{"trim_leading_from", "SELECT TRIM(LEADING 'x' FROM 'xax')", false},
		{"trim_expr_from", "SELECT TRIM(' ' FROM ' a ')", false},
		{"position_in", "SELECT POSITION('b' IN 'abc')", false},
		{"values_func", "SELECT VALUES(a) FROM t", false},
		{"quantified_any", "SELECT * FROM t WHERE a = ANY (SELECT b FROM t)", false},
		{"quantified_all", "SELECT * FROM t WHERE a > ALL (SELECT b FROM t)", false},
		{"quantified_some", "SELECT * FROM t WHERE a = SOME (SELECT b FROM t)", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c.assertParity(t, tc.sql, tc.wantAccept)
		})
	}
}
