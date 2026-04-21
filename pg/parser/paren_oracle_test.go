//go:build oracle

package parser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// This file builds the PG 17 testcontainer-backed oracle harness for
// pg-paren-dispatch Phase 2 (SCENARIOS §2.1). It mirrors the pattern in
// first_set_oracle_test.go — sync.Once container startup, CI-vs-local
// skip policy, per-probe transaction + rollback, recover() around
// MustExtractDockerHost panics. Exports are consumed by the corpus
// workers for sections 2.2–2.8.

// ParenOracle wraps the shared PG 17 container for `(` dispatch probes.
// Tables T, U, V, W and helper function probe_f() are created in the
// container at startup so §2.2–§2.7 corpus probes can assume them.
type ParenOracle struct {
	db  *sql.DB
	ctx context.Context
}

// parenOracleSetupError is captured inside parenOracleOnce.Do and
// re-observed by every caller so the first Docker failure propagates to
// all tests (fail in CI, skip locally). Same rationale as
// oracleSetupError in first_set_oracle_test.go.
var parenOracleSetupError error

var (
	parenOracleOnce sync.Once
	parenOracleInst *ParenOracle
)

// StartParenOracle returns the singleton ParenOracle. Callers from
// other test files (§2.2–§2.7 corpus workers) use this entry point.
// On Docker-unavailable local dev the caller skips; in CI it fails.
func StartParenOracle(t *testing.T) *ParenOracle {
	t.Helper()
	parenOracleOnce.Do(func() {
		defer func() {
			// testcontainers-go panics from MustExtractDockerHost when
			// the Docker socket is reachable but access is denied. Mirror
			// first_set_oracle_test.go's recover() so the CI-vs-local
			// skip branch below still fires.
			if r := recover(); r != nil {
				parenOracleSetupError = fmt.Errorf("docker provider panic: %v", r)
			}
		}()

		ctx := context.Background()
		container, err := tcpg.Run(ctx, "postgres:17-alpine",
			tcpg.WithDatabase("omni_paren"),
			tcpg.WithUsername("postgres"),
			tcpg.WithPassword("test"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2)),
		)
		if err != nil {
			parenOracleSetupError = fmt.Errorf("container start: %w", err)
			return
		}
		connStr, err := container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			parenOracleSetupError = fmt.Errorf("conn string: %w", err)
			return
		}
		db, err := sql.Open("pgx", connStr)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			parenOracleSetupError = fmt.Errorf("db open: %w", err)
			return
		}
		if err := db.PingContext(ctx); err != nil {
			db.Close()
			_ = testcontainers.TerminateContainer(container)
			parenOracleSetupError = fmt.Errorf("ping: %w", err)
			return
		}

		// Pre-create fixtures. Columns `a, b, x, y, doc` cover every
		// column referenced by §2.2–§2.7 scenarios. Functions
		// probe_f()/probe_g() exist so corpus shapes like
		// `LATERAL f(t.x)` and `ROWS FROM (f(...), g(...))` can hit
		// real objects (reject vs accept classification then reflects
		// the grammar question, not an "undefined function" error).
		setup := []string{
			`CREATE TABLE T (a int, b int, x int, y int, doc xml)`,
			`CREATE TABLE U (a int, b int, x int, y int)`,
			`CREATE TABLE V (a int, b int)`,
			`CREATE TABLE W (a int, b int)`,
			`CREATE TABLE foo (a int)`,
			`CREATE FUNCTION probe_f(i int) RETURNS int LANGUAGE sql AS 'SELECT $1'`,
			`CREATE FUNCTION probe_g(i int) RETURNS int LANGUAGE sql AS 'SELECT $1'`,
		}
		for _, s := range setup {
			if _, err := db.ExecContext(ctx, s); err != nil {
				db.Close()
				_ = testcontainers.TerminateContainer(container)
				parenOracleSetupError = fmt.Errorf("fixture %q: %w", s, err)
				return
			}
		}

		parenOracleInst = &ParenOracle{db: db, ctx: ctx}
	})

	if parenOracleSetupError != nil {
		// CI must fail — silently skipping erases the entire regression
		// fence this harness exists to provide.
		if isCI() {
			t.Fatalf("paren oracle unavailable in CI: %v", parenOracleSetupError)
		}
		t.Skipf("paren oracle unavailable (local dev): %v", parenOracleSetupError)
	}
	return parenOracleInst
}

// PGStatus classifies a PG 17 probe outcome for `SELECT * FROM (...)`
// dispatch questions. We collapse to a 2-way axis (Accept/Reject)
// because the AST-shape distinction is carried on the omni side; PG's
// ScalarArrayOpExpr vs SubLink differences are tested in §3.x corpus.
type PGStatus int

const (
	// PGAccept — the probe parsed successfully, i.e. PG did NOT raise
	// 42601 syntax_error. Semantic errors (42704 undefined_object,
	// 42P01 undefined_table, etc.) still count as Accept for `(`-
	// dispatch purposes: the grammar decision was made before name
	// resolution. Same principle as probeAccept in first_set_oracle.
	PGAccept PGStatus = iota
	PGReject
)

func (s PGStatus) String() string {
	switch s {
	case PGAccept:
		return "accept"
	case PGReject:
		return "reject(42601)"
	default:
		return "unknown"
	}
}

// OmniStatus classifies how omni's parser routed the paren-enclosed
// FROM item. The classifier inspects the parsed AST of
// `SELECT * FROM (...)` and reports which node was produced — that IS
// the routing decision the harness needs to validate.
type OmniStatus int

const (
	// OmniSubquery — FROM item is *nodes.RangeSubselect, i.e. the `(`
	// was routed to select_with_parens.
	OmniSubquery OmniStatus = iota
	// OmniJoinedTable — FROM item is *nodes.JoinExpr, i.e. `(` was
	// routed to `(` joined_table `)`.
	OmniJoinedTable
	// OmniOther — FROM item parsed but is neither subquery nor
	// joined_table (RangeVar, RangeFunction, RangeTableFunc, JsonTable,
	// etc.). Accept-but-different-shape; corpus workers can filter on
	// this.
	OmniOther
	// OmniRejected — omni returned a parse error.
	OmniRejected
)

func (s OmniStatus) String() string {
	switch s {
	case OmniSubquery:
		return "subquery"
	case OmniJoinedTable:
		return "joined_table"
	case OmniOther:
		return "other"
	case OmniRejected:
		return "reject"
	default:
		return "unknown"
	}
}

// ParenProbeResult is one probe's two-sided record. §2.2–§2.7 corpus
// workers collect these and the mismatch detector scans for
// pg_status vs omni_status disagreements.
type ParenProbeResult struct {
	SQL        string
	PGStatus   PGStatus
	OmniStatus OmniStatus
	PGError    string // SQLSTATE + message on PGReject; empty on Accept
	OmniError  string // omni parse error message on OmniRejected; empty otherwise
	Duration   time.Duration
}

// ProbeParen runs one SQL statement through both PG and omni and
// returns the side-by-side probe record. Caller owns the *testing.T
// integration (assertParenParity does that); ProbeParen itself never
// calls t.Fatal so it can be used from loops and from non-test code.
//
// Classifies the first FROM item (index 0). For non-first FROM items
// (LATERAL probes in §2.6), use ProbeParenAt.
func ProbeParen(ctx context.Context, o *ParenOracle, sqlStr string) *ParenProbeResult {
	return ProbeParenAt(ctx, o, sqlStr, 0)
}

// ProbeParenAt is the general form: classify the FROM-item at the
// supplied 0-based index. Used by §2.6 LATERAL probes where the
// interesting node sits at Items[1].
func ProbeParenAt(ctx context.Context, o *ParenOracle, sqlStr string, index int) *ParenProbeResult {
	start := time.Now()
	pgStatus, pgSQLState, pgMsg := classifyPG(ctx, o, sqlStr)
	omniStatus, omniErr := classifyOmniAt(sqlStr, index)
	r := &ParenProbeResult{
		SQL:        sqlStr,
		PGStatus:   pgStatus,
		OmniStatus: omniStatus,
		OmniError:  omniErr,
		Duration:   time.Since(start),
	}
	if pgStatus == PGReject {
		r.PGError = fmt.Sprintf("%s: %s", pgSQLState, pgMsg)
	}
	return r
}

// classifyPG runs sqlStr inside a transaction with a short timeout and
// rolls back so DDL/DML side effects can't leak between probes. The
// timeout guards against pathological probes hanging the test binary.
func classifyPG(ctx context.Context, o *ParenOracle, sqlStr string) (PGStatus, string, string) {
	probeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := o.db.BeginTx(probeCtx, nil)
	if err != nil {
		// Treat infra-level errors as "reject" with a synthetic
		// SQLSTATE so the caller can distinguish them from real 42601.
		// This mirrors how first_set_oracle handles begin failures,
		// except we return a tagged error instead of t.Fatal so
		// ProbeParen stays Fatalling-free.
		return PGReject, "XX000", fmt.Sprintf("begin: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	_, execErr := tx.ExecContext(probeCtx, sqlStr)
	if execErr == nil {
		return PGAccept, "", ""
	}
	var pgErr *pgconn.PgError
	if errors.As(execErr, &pgErr) {
		if pgErr.Code == "42601" {
			return PGReject, pgErr.Code, pgErr.Message
		}
		// Any non-syntax PG error (undefined_object, undefined_table,
		// datatype_mismatch, etc.) means the parser accepted the SQL.
		// That's what `(` dispatch cares about.
		return PGAccept, "", ""
	}
	// Non-PG error (driver / network). Classify as reject and surface
	// the Go error text — the harness will flag the mismatch loudly if
	// it happens.
	return PGReject, "XXNET", execErr.Error()
}

// classifyOmni parses sqlStr through omni's pg parser and classifies
// the first FROM item of the first (or only) SELECT statement. The
// classifier only handles `SELECT ... FROM (...)` shapes — that's the
// Phase 2 probe template. For any other input shape the result is
// OmniOther, which corpus workers treat as "do not assert parity on
// this probe" noise filter.
//
// This is the Items[0] specialization of classifyOmniAt; kept as a
// compatibility shim so §2.2–§2.5 and §2.7 corpus workers don't have to
// thread an explicit index when the first FROM item is the one being
// probed.
func classifyOmni(sqlStr string) (OmniStatus, string) {
	return classifyOmniAt(sqlStr, 0)
}

// classifyOmniAt is the general-form classifier: it inspects the
// FROM-item at the caller-specified 0-based index. §2.6 LATERAL probes
// put the LATERAL-prefixed node at Items[1] (Items[0] is the anchor
// relation T), so classifyOmni alone would silently always report
// OmniOther for the anchor and never exercise the LATERAL routing
// decision. Corpus workers that care about a non-first FROM item call
// this variant with the right index.
//
// Index out-of-range returns OmniOther, "" — same shape as any other
// unrouteable input, so the caller's expected-status table can still
// assert something meaningful (typically OmniOther or OmniRejected).
func classifyOmniAt(sqlStr string, index int) (OmniStatus, string) {
	stmts, err := Parse(sqlStr)
	if err != nil {
		return OmniRejected, err.Error()
	}
	if stmts == nil || len(stmts.Items) == 0 {
		return OmniRejected, "no statements"
	}
	raw, ok := stmts.Items[0].(*nodes.RawStmt)
	if !ok {
		return OmniOther, ""
	}
	sel, ok := raw.Stmt.(*nodes.SelectStmt)
	if !ok {
		return OmniOther, ""
	}
	if sel.FromClause == nil || len(sel.FromClause.Items) == 0 {
		// `SELECT 1` with no FROM — not a paren-routing question.
		return OmniOther, ""
	}
	if index < 0 || index >= len(sel.FromClause.Items) {
		return OmniOther, ""
	}
	item := sel.FromClause.Items[index]
	switch item.(type) {
	case *nodes.RangeSubselect:
		return OmniSubquery, ""
	case *nodes.JoinExpr:
		return OmniJoinedTable, ""
	default:
		// RangeVar, RangeFunction, RangeTableFunc, JsonTable, etc.
		return OmniOther, ""
	}
}

// expectedOmniFor maps a PGStatus + SQL into what OmniStatus the
// harness considers a MATCH. The asymmetry:
//   - PGAccept alone isn't enough: PG accepts both `(SELECT 1)` (subquery)
//     and `(t JOIN u ON TRUE)` (joined_table), but we want omni to route
//     them to different nodes. The caller passes the EXPECTED omni
//     shape (OmniSubquery / OmniJoinedTable / OmniOther) alongside
//     PGAccept so the mismatch detector can flag a wrong-side routing.
//   - PGReject requires OmniRejected.
//
// assertParenParity accepts the expected OmniStatus directly (caller-
// supplied), so this mapping is encoded in the test table, not here.

// assertParenParity is the main test entry point. It probes PG + omni
// and fails the test if omni's routing disagrees with the expected
// PG-aligned outcome. `expected` is the OmniStatus the caller asserts
// for this SQL — encoding the PG-17 truth. On mismatch the helper
// prints a side-by-side diff of SQL / pg_status / omni_status /
// errors so debugging doesn't require re-running the probe by hand.
func assertParenParity(t *testing.T, o *ParenOracle, sqlStr string, expected OmniStatus) {
	t.Helper()
	assertParenParityAt(t, o, sqlStr, 0, expected)
}

// assertParenParityAt is the index-aware variant of assertParenParity.
// §2.6 LATERAL probes call this with index=1 so the classifier inspects
// the LATERAL-prefixed FROM item (Items[1]) rather than the anchor
// relation at Items[0].
func assertParenParityAt(t *testing.T, o *ParenOracle, sqlStr string, index int, expected OmniStatus) {
	t.Helper()
	r := ProbeParenAt(o.ctx, o, sqlStr, index)

	// Expected pg_status: Reject iff the caller expects OmniRejected.
	// This couples "PG rejects" to "omni must also reject" tightly,
	// which is exactly the invariant we want: accept-vs-reject drift is
	// more serious than sub-vs-join drift.
	wantPGReject := expected == OmniRejected
	if wantPGReject && r.PGStatus != PGReject {
		t.Errorf("parity mismatch for %q:\n"+
			"  expected PG reject (omni expected: %s)\n"+
			"  got      PG accept, omni=%s\n"+
			"  duration %v",
			sqlStr, expected, r.OmniStatus, r.Duration)
		return
	}
	if !wantPGReject && r.PGStatus == PGReject {
		t.Errorf("parity mismatch for %q:\n"+
			"  expected PG accept (omni expected: %s)\n"+
			"  got      PG reject: %s\n"+
			"  omni     %s\n"+
			"  duration %v",
			sqlStr, expected, r.PGError, r.OmniStatus, r.Duration)
		return
	}

	if r.OmniStatus != expected {
		t.Errorf("parity mismatch for %q:\n"+
			"  expected omni=%s\n"+
			"  got      omni=%s (error=%q)\n"+
			"  pg       %s\n"+
			"  duration %v",
			sqlStr, expected, r.OmniStatus, r.OmniError, r.PGStatus, r.Duration)
	}
}

// TestParenOracleHarness is the Section 2.1 self-test. It exercises a
// minimal handful of probes to prove:
//  1. Container boots and accepts DDL at startup (fixtures applied).
//  2. ProbeParen round-trips for both accept and reject SQL.
//  3. classifyOmni distinguishes RangeSubselect vs JoinExpr vs reject.
//  4. assertParenParity reports mismatches with a side-by-side diff.
//  5. Per-probe timing is captured (budget: ~2 min for 200 probes → ~600ms/probe).
//
// The broader corpus (~200 probes across §2.2–§2.7) is deferred to the
// respective sections' test files, each of which calls StartParenOracle
// and assertParenParity using the exports above.
func TestParenOracleHarness(t *testing.T) {
	o := StartParenOracle(t)

	cases := []struct {
		name     string
		sql      string
		expected OmniStatus
	}{
		{
			name:     "joined_table single-level",
			sql:      `SELECT * FROM (T JOIN U ON TRUE)`,
			expected: OmniJoinedTable,
		},
		{
			name:     "subquery single-level",
			sql:      `SELECT * FROM (SELECT 1)`,
			expected: OmniSubquery,
		},
		{
			name:     "empty parens rejected",
			sql:      `SELECT * FROM ()`,
			expected: OmniRejected,
		},
		{
			name:     "single relation in parens rejected",
			sql:      `SELECT * FROM (T)`,
			expected: OmniRejected,
		},
		{
			name:     "BYT-9315 shape",
			sql:      `SELECT * FROM ((T JOIN U ON TRUE) JOIN V ON TRUE)`,
			expected: OmniJoinedTable,
		},
		{
			// OmniOther fixture: bare RangeVar doesn't exercise `(` dispatch
			// and proves the classifier's default branch (RangeVar/
			// RangeFunction/RangeTableFunc/JsonTable) is reachable from the
			// self-test, not only from §2.6 where LATERAL hides behind
			// Items[1].
			name:     "bare relation (OmniOther)",
			sql:      `SELECT * FROM T`,
			expected: OmniOther,
		},
	}

	var total time.Duration
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := ProbeParen(o.ctx, o, tc.sql)
			total += r.Duration
			t.Logf("probe %q → pg=%s omni=%s in %v",
				summarize(tc.sql), r.PGStatus, r.OmniStatus, r.Duration)
			assertParenParity(t, o, tc.sql, tc.expected)
		})
	}
	t.Logf("total harness probe time: %v across %d cases (avg %v)",
		total, len(cases), total/time.Duration(max(len(cases), 1)))
}

// summarize trims long SQL for log lines so `go test -v` output stays
// readable when §2.2+ corpus files start logging dozens of probes.
func summarize(sqlStr string) string {
	const limit = 60
	s := strings.ReplaceAll(sqlStr, "\n", " ")
	if len(s) <= limit {
		return s
	}
	return s[:limit-1] + "…"
}
