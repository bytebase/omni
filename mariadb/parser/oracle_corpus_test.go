package parser

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	mariadb "github.com/testcontainers/testcontainers-go/modules/mariadb"
)

// ============================================================================
// MariaDB 11.8 container oracle + pinned subtractive-divergence inventory.
//
// The DoD is the *pinned reviewed list* (subtractiveDivergences), NOT green
// parity: omni's mysql-forked parser over-accepts MySQL-isms that real MariaDB
// rejects. TestMariaDBDivergenceInventory runs the mysql corpus + category
// probes through omni's Parse and a live mariadb:11.8.8 container and asserts
// the live mismatch set EQUALS the allowlist, bidirectionally (UNREVIEWED if a
// live divergence isn't pinned; STALE if a pinned entry no longer diverges).
//
// SUBTRACTIVE & 1064-SCOPED ONLY: blind to additive MariaDB surface (SEQUENCE,
// RETURNING, sql_mode=ORACLE, system-versioning → Phase-0 mdbcheck) and to
// feature-gated (1235) rejections. "Inventory pinned" != "MariaDB parity".
// ============================================================================

// startMariaDB starts mariadb:11.8.8 via the testcontainers mariadb module. The
// module's wait strategy handles MariaDB's double-start entrypoint (temp init
// server → shutdown → real server) that a bare port-wait would race.
func startMariaDB(t *testing.T) *parserOracle {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping MariaDB container oracle in short mode")
	}
	ctx := context.Background()

	c, err := mariadb.Run(ctx, "mariadb:11.8.8",
		mariadb.WithDatabase("test"),
		mariadb.WithPassword("test"),
		// no WithUsername("root"): root is the module default; setting it errors.
	)
	if err != nil {
		t.Skipf("MariaDB container unavailable: %v", err)
	}
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(c) })

	// Server-side guards (seconds) on every pooled conn via DSN system-vars: the
	// server aborts a slow/locked statement so the conn stays healthy (a client
	// context cancel leaves the query running server-side and poisons the pool).
	connStr, err := c.ConnectionString(ctx,
		"multiStatements=true", "parseTime=true",
		"max_statement_time=5", "lock_wait_timeout=2",
	)
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}
	return &parserOracle{db: db, ctx: ctx}
}

type verdict int

const (
	vAccepted verdict = iota
	vRejected
	vErrored
)

// classify runs sqlStr against MariaDB and maps the result to a PARSE verdict.
// Only 1064 (ER_PARSE_ERROR) is a reject — MariaDB reserves 1064 for parse-time
// errors (unlike StarRocks's 1064 overload). Any other *gomysql.MySQLError means
// it PARSED → accepted (code captured): 1146 (missing table), 1235
// (ER_NOT_SUPPORTED_YET feature-gate), and 1969 (ER_STATEMENT_TIMEOUT from the
// server-side max_statement_time guard) all surface as MySQLErrors = "parsed
// fine, failed/killed at runtime" → accepted. A NON-MySQLError (conn reset, pool
// exhaustion) is a distinct vErrored state — surfaced loudly, never folded into
// accept (that would corrupt the exact-set gate).
func (o *parserOracle) classify(sqlStr string) (verdict, uint16, string) {
	_, err := o.db.ExecContext(o.ctx, sqlStr)
	if err == nil {
		return vAccepted, 0, ""
	}
	if myErr, ok := err.(*gomysql.MySQLError); ok {
		if myErr.Number == 1064 {
			return vRejected, 1064, myErr.Message
		}
		return vAccepted, myErr.Number, myErr.Message
	}
	return vErrored, 0, err.Error()
}

// categoryCaseSQLs are the hand-written type/window/interval/reject probes lifted
// from the former TestOracleCorpus, fed into the one inventory gate below.
var categoryCaseSQLs = []string{
	// integer types × modifiers
	"CREATE TABLE t1 (a INT)",
	"CREATE TABLE t2 (a INT UNSIGNED)",
	"CREATE TABLE t3 (a INT SIGNED)",
	"CREATE TABLE t4 (a INT(11))",
	"CREATE TABLE t5 (a INT UNSIGNED ZEROFILL)",
	"CREATE TABLE t6 (a TINYINT SIGNED)",
	"CREATE TABLE t7 (a TINYINT UNSIGNED ZEROFILL)",
	"CREATE TABLE t8 (a SMALLINT(5) UNSIGNED)",
	"CREATE TABLE t9 (a MEDIUMINT SIGNED)",
	"CREATE TABLE t10 (a BIGINT(20) UNSIGNED)",
	"CREATE TABLE t11 (a INT1 UNSIGNED)",
	"CREATE TABLE t12 (a INT8 SIGNED)",
	// decimal & float types
	"CREATE TABLE t13 (a DECIMAL)",
	"CREATE TABLE t14 (a DECIMAL(10))",
	"CREATE TABLE t15 (a DECIMAL(10,2))",
	"CREATE TABLE t16 (a DECIMAL(10,2) UNSIGNED)",
	"CREATE TABLE t17 (a NUMERIC(10,2))",
	"CREATE TABLE t18 (a DEC(10,2))",
	"CREATE TABLE t19 (a FIXED(10,2))",
	"CREATE TABLE t20 (a FLOAT)",
	"CREATE TABLE t21 (a FLOAT(10,2))",
	"CREATE TABLE t22 (a FLOAT(24))",
	"CREATE TABLE t23 (a FLOAT(25))",
	"CREATE TABLE t24 (a DOUBLE)",
	"CREATE TABLE t25 (a DOUBLE PRECISION)",
	"CREATE TABLE t26 (a REAL)",
	"CREATE TABLE t27 (a FLOAT4)",
	"CREATE TABLE t28 (a FLOAT8)",
	// string & binary types
	"CREATE TABLE t29 (a CHAR(10))",
	"CREATE TABLE t30 (a CHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci)",
	"CREATE TABLE t31 (a VARCHAR(255))",
	"CREATE TABLE t32 (a TEXT)",
	"CREATE TABLE t33 (a TEXT(1000))",
	"CREATE TABLE t34 (a TINYTEXT)",
	"CREATE TABLE t35 (a MEDIUMTEXT CHARACTER SET latin1)",
	"CREATE TABLE t36 (a LONGTEXT)",
	"CREATE TABLE t37 (a LONG)",
	"CREATE TABLE t38 (a LONG VARCHAR)",
	"CREATE TABLE t39 (a BINARY(16))",
	"CREATE TABLE t40 (a VARBINARY(255))",
	"CREATE TABLE t41 (a BLOB)",
	"CREATE TABLE t42 (a BLOB(1000))",
	"CREATE TABLE t43 (a TINYBLOB)",
	"CREATE TABLE t44 (a MEDIUMBLOB)",
	"CREATE TABLE t45 (a LONGBLOB)",
	"CREATE TABLE t46 (a LONG VARBINARY)",
	"CREATE TABLE t47 (a NATIONAL CHAR(10))",
	"CREATE TABLE t48 (a NCHAR(10))",
	"CREATE TABLE t49 (a NVARCHAR(100))",
	// date/time, json, spatial & special
	"CREATE TABLE t50 (a DATE)",
	"CREATE TABLE t51 (a TIME)",
	"CREATE TABLE t52 (a TIME(3))",
	"CREATE TABLE t53 (a DATETIME)",
	"CREATE TABLE t54 (a DATETIME(6))",
	"CREATE TABLE t55 (a TIMESTAMP)",
	"CREATE TABLE t56 (a TIMESTAMP(3))",
	"CREATE TABLE t57 (a YEAR)",
	"CREATE TABLE t58 (a BIT(8))",
	"CREATE TABLE t59 (a BOOL)",
	"CREATE TABLE t60 (a JSON)",
	"CREATE TABLE t61 (a SERIAL)",
	"CREATE TABLE t62 (a ENUM('a','b','c'))",
	"CREATE TABLE t63 (a SET('x','y','z'))",
	"CREATE TABLE t64 (a GEOMETRY)",
	"CREATE TABLE t65 (a POINT)",
	"CREATE TABLE t66 (a LINESTRING)",
	"CREATE TABLE t67 (a POLYGON)",
	"CREATE TABLE t68 (a GEOMETRYCOLLECTION)",
	// window functions
	"SELECT RANK() OVER w FROM t WINDOW w AS (ORDER BY id)",
	"SELECT ROW_NUMBER() OVER (PARTITION BY dept ORDER BY salary DESC) FROM t",
	"SELECT LAG(val) OVER (ORDER BY id) FROM t",
	"SELECT LAG(val, 2, 0) OVER (ORDER BY id) FROM t",
	"SELECT NTH_VALUE(val, 3) OVER (ORDER BY id) FROM t",
	"SELECT SUM(val) OVER (ORDER BY id ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t",
	"SELECT DENSE_RANK() OVER (ORDER BY score), PERCENT_RANK() OVER (ORDER BY score) FROM t",
	// interval expressions
	"SELECT NOW() + INTERVAL 1 DAY",
	"SELECT NOW() - INTERVAL 2 HOUR",
	"SELECT DATE_ADD('2024-01-01', INTERVAL 1 MONTH)",
	"SELECT DATE_SUB(NOW(), INTERVAL 30 MINUTE)",
	"SELECT TIMESTAMPADD(MINUTE, 30, NOW())",
	"SELECT EXTRACT(HOUR FROM NOW())",
	"SELECT * FROM t WHERE created_at > NOW() - INTERVAL 7 DAY",
	"SELECT * FROM t WHERE created_at BETWEEN NOW() - INTERVAL 1 MONTH AND NOW()",
	// expected-reject probes (incl. the former knownMismatch reserved-word cases)
	"CREATE TABLE t_rej1 (a VARCHAR UNSIGNED)",
	"CREATE TABLE t_rej2 (a TEXT ZEROFILL)",
	"CREATE TABLE t_rej3 (a JSON UNSIGNED)",
	"CREATE TABLE select (a INT)",
	"CREATE TABLE t_rej5 (select INT)",
	"ALTER TABLE t PARTITION BY",
	"ALTER TABLE t PARTITION BY RANGE(id)",
	"SELECT DATE_ADD(d, INTERVAL 1 INVALID_UNIT) FROM t",
	// MariaDB subtractive divergences — Phase-0 mdbcheck OVER findings (omni accepts, MariaDB rejects)
	"SELECT 1 MINUS SELECT 2",
	"SELECT 1 MINUS SELECT 2 MINUS SELECT 3",
	"SELECT * FROM t WHERE id = 1 FOR SHARE",
	"SELECT name -> '$.b' FROM t",
	"SELECT name ->> '$.b' FROM t",
}

// containerFatalVerbs name statements whose *execution* would kill the shared
// container or connection (not caught by max_statement_time); skip-listed.
var containerFatalVerbs = []string{"SHUTDOWN", "KILL"}

// gatherOracleInputs collects the corpus statements + category cases, deduped on
// the EXACT statement text. No whitespace normalization: collapsing whitespace
// would corrupt string literals; the corpus is stable committed data, so
// reformat-churn is low-risk and would surface as UNREVIEWED+STALE pairs to re-pin.
func gatherOracleInputs(t *testing.T) []string {
	t.Helper()
	seen := map[string]bool{}
	var out []string
	add := func(sqlStr string) {
		sqlStr = strings.TrimSpace(sqlStr)
		if sqlStr == "" || seen[sqlStr] {
			return
		}
		up := strings.ToUpper(sqlStr)
		for _, v := range containerFatalVerbs {
			if strings.HasPrefix(up, v) {
				t.Logf("skip (container-fatal verb %s): %q", v, sqlStr)
				return
			}
		}
		seen[sqlStr] = true
		out = append(out, sqlStr)
	}
	dir := filepath.Join("..", "quality", "corpus")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read corpus dir: %v", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		for _, s := range loadCorpusStatements(t, filepath.Join(dir, e.Name())) {
			add(s.sql)
		}
	}
	for _, s := range categoryCaseSQLs {
		add(s)
	}
	return out
}

// divergence is one reviewed entry: a statement where omni's parser and real
// MariaDB 11.8 disagree, with the attributed MariaDB reason.
type divergence struct {
	sql      string // EXACT statement text (must match gatherOracleInputs output)
	omniSide bool   // true: omni accepts & MariaDB rejects (subtractive over-accept); false: omni rejects & MariaDB accepts
	reason   string
}

// subtractiveDivergences — the REVIEWED inventory and the P1+ backlog. Authored
// from the measure run against mariadb:11.8.8 over the mysql corpus + category
// probes, each entry's engine behavior re-confirmed on real MySQL 8.0 AND
// MariaDB 11.8.8. The corpus exercises subtractive divergence only incidentally
// — "1 entry" is a LOWER BOUND, not the full subtractive surface; the authored
// MariaDB corpus (Phase-0 mdbcheck) is where the bulk (additive + more
// subtractive) lives.
var subtractiveDivergences = []divergence{
	{
		sql:      "SELECT LAG(val, 2, 0) OVER (ORDER BY id) FROM t",
		omniSide: true, // MySQL 8.0 accepts, MariaDB 11.8.8 rejects (1064); omni mirrors MySQL.
		reason:   "LAG()/LEAD() 3rd 'default value' argument: VALID in MySQL 8.0, REJECTED by MariaDB 11.8.8 at parse (1064) — confirmed on both real engines. omni mirrors MySQL. (LEAD shares this divergence; it isn't a separate corpus statement.) P1 candidate: gate the window-function default-value arg behind a MariaDB flag.",
	},
	{
		sql:      "SELECT 1 MINUS SELECT 2",
		omniSide: true, // omni over-accepts MINUS; MariaDB rejects it in default mode (1064).
		reason:   "MINUS set operator: omni over-accepts it unconditionally; MariaDB accepts MINUS only under sql_mode=ORACLE, rejects it in default mode (1064). Phase-0 mdbcheck OVER finding. Candidate: gate MINUS behind the Oracle sql_mode flag.",
	},
	{
		sql:      "SELECT 1 MINUS SELECT 2 MINUS SELECT 3",
		omniSide: true, // chained MINUS — same divergence as the 2-operand form.
		reason:   "MINUS set operator (chained): same divergence as the single MINUS — omni over-accepts; MariaDB default-mode rejects (1064), Oracle-mode only.",
	},
	{
		sql:      "SELECT * FROM t WHERE id = 1 FOR SHARE",
		omniSide: true, // omni mirrors MySQL's FOR SHARE; MariaDB rejects it (1064).
		reason:   "FOR SHARE locking clause: valid in MySQL 8.0 and omni mirrors it; MariaDB has no FOR SHARE (uses LOCK IN SHARE MODE) and rejects at parse (1064). Phase-0 mdbcheck OVER finding.",
	},
	{
		sql:      "SELECT name -> '$.b' FROM t",
		omniSide: true, // omni mirrors MySQL's -> JSON operator; MariaDB rejects it (1064).
		reason:   "JSON -> column-path operator: MySQL-only; omni mirrors MySQL and over-accepts. MariaDB has no -> / ->> operators (uses JSON_EXTRACT()/JSON_UNQUOTE()) and rejects at parse (1064). Deferred prune: the one-line grammar fix (drop the tokJsonExtract/tokJsonUnquote arms in expr.go) ripples into inherited mysql accept-tests (TestParseJsonExtract/UnquoteExtract) and the routine_body_audit, beyond the surgical FOR-UPDATE-OF arm.",
	},
	{
		sql:      "SELECT name ->> '$.b' FROM t",
		omniSide: true, // ->> shares the -> divergence.
		reason:   "JSON ->> column-path unquote operator: MySQL-only; same divergence as -> (omni over-accepts, MariaDB rejects 1064). Deferred with the -> prune.",
	},
}

func TestMariaDBDivergenceInventory(t *testing.T) {
	o := startMariaDB(t)

	// Seed the referent table a few category SELECTs reference (window/interval).
	_, _ = o.db.ExecContext(o.ctx,
		"CREATE TABLE IF NOT EXISTS t (id INT, val INT, score INT, status INT, name VARCHAR(100), dept VARCHAR(50), salary DECIMAL(10,2), d DATE, created_at DATETIME)")

	want := map[string]divergence{}
	for _, d := range subtractiveDivergences {
		want[d.sql] = d
	}

	live := map[string]bool{} // exact sql → omniAccepts, for diverging stmts only
	codeTally := map[uint16]int{}
	for _, sqlStr := range gatherOracleInputs(t) {
		_, perr := Parse(sqlStr)
		omniAccepts := perr == nil
		v, code, msg := o.classify(sqlStr)
		codeTally[code]++
		if v == vErrored {
			t.Errorf("ERRORED (infra, not a parse verdict) for %q: %s", sqlStr, msg)
			continue
		}
		mdbAccepts := v == vAccepted
		if omniAccepts != mdbAccepts {
			live[sqlStr] = omniAccepts
			t.Logf("divergence omni=%v mariadb=%v code=%d %q", omniAccepts, mdbAccepts, code, sqlStr)
		}
	}

	for sqlStr, omniAccepts := range live {
		d, ok := want[sqlStr]
		if !ok {
			t.Errorf("UNREVIEWED divergence (omni=%v mariadb=%v): %q\n  → add to subtractiveDivergences with a reason, or fix the parser", omniAccepts, !omniAccepts, sqlStr)
			continue
		}
		if d.omniSide != omniAccepts {
			t.Errorf("divergence DIRECTION changed for %q (was omniSide=%v, now omniAccepts=%v)", sqlStr, d.omniSide, omniAccepts)
		}
	}
	for sqlStr := range want {
		if _, ok := live[sqlStr]; !ok {
			t.Errorf("STALE allowlist entry (no longer diverges): %q\n  → remove from subtractiveDivergences (resolved)", sqlStr)
		}
	}
	t.Logf("MariaDB inventory: %d pinned, %d live | code tally (0=clean accept): %v", len(want), len(live), codeTally)
}
