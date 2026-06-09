//go:build googlesql_oracle

// Differential test for the `expressions` node against the live Cloud Spanner
// emulator oracle. Run with:
//
//	SPANNER_EMULATOR_HOST=localhost:9010 \
//	  go test -tags googlesql_oracle ./googlesql/parser/ -run TestExprDifferential
//
// It is the PROVE gate for googlesql/expressions (correctness-protocol.md): for
// every fixture it (1) feeds a full statement embedding the expression to the
// emulator via the harness/googlesql-spanner CLI and reads the accept/reject
// verdict, then (2) parses the bare expression text with omni's ParseExpression,
// and asserts the two agree — BOTH polarities (the corpus has accept AND reject
// cases). A harness `error` verdict (oracle could not decide) fails the fixture
// loudly; it is never folded into accept or reject.
//
// The emulator speaks Spanner's dialect, a SUBSET of the BigQuery+Spanner union
// the parser must accept (oracle.md). Every expression form covered here is in
// the SHARED GoogleSQL core (the full precedence chain, function calls, CASE/
// CAST/EXTRACT, array/struct constructors, access, INTERVAL, AT TIME ZONE,
// parameters, subqueries), so the Spanner verdict is authoritative. A few
// fixtures parse on the grammar but are semantically rejected by Spanner
// ("X is not supported", "Table not found", "Unrecognized name", "IS DISTINCT
// FROM is not supported", "No matching signature") — those are classified ACCEPT
// (the grammar accepted) and asserted as accept. BigQuery-only expression forms
// (e.g. WITH(...) inline expression, REPLACE_FIELDS, NEW proto, braced proto
// constructors) are NOT authoritative against this emulator; they are covered by
// the hand-written unit tests (expr_test.go) and triangulated against the legacy
// .g4 + docs, and are deliberately excluded from this differential.
package parser

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

// exprFixture pairs a bare expression (fed to omni's ParseExpression) with a
// full statement embedding that expression in a position where the emulator
// yields a clean accept/reject verdict. wantParse is the expected omni outcome
// and MUST equal the emulator's grammar verdict.
type exprFixture struct {
	expr      string // bare expression for ParseExpression
	stmt      string // full statement for the oracle (one stmt, no trailing ';')
	wantParse bool   // expected: omni parses (and oracle accepts)
}

// exprFixtures are oracle-probed expression forms. Each stmt is chosen so the
// emulator returns a deterministic accept/reject (NOT an Unimplemented "error"):
// most use `SELECT <expr> FROM t` (Table-not-found ⇒ accept) or `SELECT <expr>`;
// struct/tuple-returning expressions use a field extraction (`.x` / `.field1`)
// to dodge the emulator's "struct value cannot be returned" Unimplemented path.
var exprFixtures = []exprFixture{
	// --- literals (accept) ---
	{"1", "SELECT 1", true},
	{"1.5", "SELECT 1.5", true},
	{"NULL", "SELECT NULL", true},
	{"TRUE", "SELECT TRUE", true},
	{"FALSE", "SELECT FALSE", true},
	{"'hello'", "SELECT 'hello'", true},
	{"'a' 'b'", "SELECT 'a' 'b'", true}, // adjacent string concatenation
	{"b'bytes'", "SELECT b'bytes'", true},
	{"0x1F", "SELECT 0x1F", true},

	// --- typed-prefix literals (accept) ---
	{"DATE '2020-01-01'", "SELECT DATE '2020-01-01'", true},
	{"TIMESTAMP '2020-01-01 00:00:00'", "SELECT TIMESTAMP '2020-01-01 00:00:00'", true},
	{"NUMERIC '1.5'", "SELECT NUMERIC '1.5'", true},
	{"JSON '{}'", "SELECT JSON '{}'", true},

	// --- parameters (accept; semantic) ---
	{"@p", "SELECT @p", true},
	{"@@sysvar", "SELECT @@sysvar", true},

	// --- identifiers / paths (accept; semantic Table-not-found / Unrecognized) ---
	{"a", "SELECT a FROM t", true},
	{"t.a.b", "SELECT t.a.b FROM t", true},

	// --- arithmetic / precedence (accept) ---
	{"1 + 2 * 3", "SELECT 1 + 2 * 3", true},
	{"(1 + 2) * 3", "SELECT (1 + 2) * 3", true},
	{"1 - 2 - 3", "SELECT 1 - 2 - 3", true},
	{"- - 1", "SELECT - - 1", true},
	{"~ 5", "SELECT ~ 5", true},
	{"5 & 3 | 1", "SELECT 5 & 3 | 1", true},
	{"1 << 2 + 3", "SELECT 1 << 2 + 3", true},
	{"'a' || 'b' || 'c'", "SELECT 'a' || 'b' || 'c'", true},

	// --- logical (accept) ---
	{"TRUE AND FALSE OR TRUE", "SELECT TRUE AND FALSE OR TRUE", true},
	{"NOT TRUE", "SELECT NOT TRUE", true},
	{"NOT NOT TRUE", "SELECT NOT NOT TRUE", true},
	{"NOT 1 = 1", "SELECT NOT 1 = 1", true},

	// --- comparison family (accept) ---
	{"1 = 2", "SELECT 1 = 2", true},
	{"1 != 2", "SELECT 1 != 2", true},
	{"1 <> 2", "SELECT 1 <> 2", true},
	{"1 < 2", "SELECT 1 < 2", true},
	{"1 >= 2", "SELECT 1 >= 2", true},
	{"1 + 2 = 3", "SELECT 1 + 2 = 3", true},
	{"a IS NULL", "SELECT a IS NULL FROM t", true},
	{"a IS NOT NULL", "SELECT a IS NOT NULL FROM t", true},
	{"1 IS TRUE", "SELECT 1 IS TRUE", true},
	{"1 IS NOT FALSE", "SELECT 1 IS NOT FALSE", true},
	{"1 IS UNKNOWN", "SELECT 1 IS UNKNOWN", true},
	{"a IS NOT DISTINCT FROM b", "SELECT a IS NOT DISTINCT FROM b FROM t", true},
	{"a BETWEEN 1 AND 2", "SELECT a BETWEEN 1 AND 2 FROM t", true},
	{"a NOT BETWEEN 1 AND 2", "SELECT a NOT BETWEEN 1 AND 2 FROM t", true},
	{"a IN (1, 2, 3)", "SELECT a IN (1, 2, 3) FROM t", true},
	{"a IN (1)", "SELECT a IN (1) FROM t", true},
	{"a NOT IN (1, 2)", "SELECT a NOT IN (1, 2) FROM t", true},
	{"a IN UNNEST([1, 2])", "SELECT a IN UNNEST([1, 2]) FROM t", true},
	{"a IN (SELECT x FROM t2)", "SELECT a IN (SELECT x FROM t2) FROM t", true},
	{"x LIKE '%a%'", "SELECT x LIKE '%a%' FROM t", true},
	{"x NOT LIKE '%a%'", "SELECT x NOT LIKE '%a%' FROM t", true},
	{"x LIKE ANY ('a', 'b')", "SELECT x LIKE ANY ('a', 'b') FROM t", true},

	// --- quantified comparison: expr {= != <> < <= > >=} {ANY|SOME|ALL} <rhs>
	// (any_some_all on comparative_operator). SHARED GoogleSQL core; the emulator
	// accepts every form (semantic Table-not-found ⇒ parsed). See divergence #201
	// and the unit suite TestExpr_QuantifiedComparison.
	{"x = ANY (SELECT v FROM s)", "SELECT * FROM t WHERE x = ANY (SELECT v FROM s)", true},
	{"x > ALL (SELECT v FROM s)", "SELECT * FROM t WHERE x > ALL (SELECT v FROM s)", true},
	{"x < SOME (SELECT v FROM s)", "SELECT * FROM t WHERE x < SOME (SELECT v FROM s)", true},
	{"x >= ANY (SELECT v FROM s)", "SELECT * FROM t WHERE x >= ANY (SELECT v FROM s)", true},
	{"x <= ALL (SELECT v FROM s)", "SELECT * FROM t WHERE x <= ALL (SELECT v FROM s)", true},
	{"x != ANY (SELECT v FROM s)", "SELECT * FROM t WHERE x != ANY (SELECT v FROM s)", true},
	{"x <> ALL (SELECT v FROM s)", "SELECT * FROM t WHERE x <> ALL (SELECT v FROM s)", true},
	{"x = ANY (1, 2, 3)", "SELECT * FROM t WHERE x = ANY (1, 2, 3)", true},
	{"x > ALL (1, 2)", "SELECT * FROM t WHERE x > ALL (1, 2)", true},
	{"x = ANY UNNEST([1, 2, 3])", "SELECT * FROM t WHERE x = ANY UNNEST([1, 2, 3])", true},
	{"x = ANY @{a=1} (SELECT v FROM s)", "SELECT * FROM t WHERE x = ANY @{a=1} (SELECT v FROM s)", true},

	// --- CASE (accept) ---
	{"CASE WHEN a THEN 1 ELSE 2 END", "SELECT CASE WHEN a THEN 1 ELSE 2 END FROM t", true},
	{"CASE a WHEN 1 THEN 2 END", "SELECT CASE a WHEN 1 THEN 2 END FROM t", true},
	{"CASE WHEN 1 THEN 2 WHEN 3 THEN 4 ELSE 5 END", "SELECT CASE WHEN 1 THEN 2 WHEN 3 THEN 4 ELSE 5 END", true},

	// --- CAST / EXTRACT (accept) ---
	{"CAST(x AS INT64)", "SELECT CAST(x AS INT64) FROM t", true},
	{"SAFE_CAST(x AS STRING)", "SELECT SAFE_CAST(x AS STRING) FROM t", true},
	{"CAST(x AS STRING FORMAT 'fmt')", "SELECT CAST(x AS STRING FORMAT 'fmt') FROM t", true},
	{"CAST(x AS STRING FORMAT 'f' AT TIME ZONE 'UTC')", "SELECT CAST(x AS STRING FORMAT 'f' AT TIME ZONE 'UTC') FROM t", true},
	{"EXTRACT(YEAR FROM d)", "SELECT EXTRACT(YEAR FROM d) FROM t", true},
	{"EXTRACT(HOUR FROM ts AT TIME ZONE 'UTC')", "SELECT EXTRACT(HOUR FROM ts AT TIME ZONE 'UTC') FROM t", true},

	// --- function calls (accept) ---
	{"f(x, y)", "SELECT f(x, y) FROM t", true},
	{"mypkg.myfunc(x)", "SELECT mypkg.myfunc(x) FROM t", true},
	{"COUNT(*)", "SELECT COUNT(*) FROM t", true},
	{"COUNT(DISTINCT x)", "SELECT COUNT(DISTINCT x) FROM t", true},
	{"IF(a, 1, 2)", "SELECT IF(a, 1, 2) FROM t", true},
	{"GROUPING(x)", "SELECT GROUPING(x) FROM t", true},
	{"ARRAY_AGG(x IGNORE NULLS)", "SELECT ARRAY_AGG(x IGNORE NULLS) FROM t", true},
	{"ARRAY_AGG(x ORDER BY y)", "SELECT ARRAY_AGG(x ORDER BY y) FROM t", true},
	{"ARRAY_AGG(x LIMIT 10)", "SELECT ARRAY_AGG(x LIMIT 10) FROM t", true},
	{"func(a => 1, b => 2)", "SELECT func(a => 1, b => 2)", true},
	{"ARRAY_TRANSFORM([1], e -> e + 1)", "SELECT ARRAY_TRANSFORM([1], e -> e + 1)", true},
	{"REDUCE([1], (a, b) -> a + b)", "SELECT REDUCE([1], (a, b) -> a + b)", true},

	// --- window (accept) ---
	{"SUM(x) OVER (PARTITION BY a ORDER BY b)", "SELECT SUM(x) OVER (PARTITION BY a ORDER BY b) FROM t", true},
	{"SUM(x) OVER w", "SELECT SUM(x) OVER w FROM t", true},
	{"ROW_NUMBER() OVER (ORDER BY x ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)",
		"SELECT ROW_NUMBER() OVER (ORDER BY x ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t", true},
	{"COUNT(*) OVER ()", "SELECT COUNT(*) OVER () FROM t", true},

	// --- array constructors (accept) ---
	{"[1, 2, 3]", "SELECT [1, 2, 3]", true},
	{"ARRAY[1, 2]", "SELECT ARRAY[1, 2]", true},
	{"ARRAY<INT64>[1, 2]", "SELECT ARRAY<INT64>[1, 2]", true},
	{"[]", "SELECT []", true}, // empty array literal (oracle: accept)
	{"ARRAY<INT64>[]", "SELECT ARRAY<INT64>[]", true},

	// --- struct constructors (accept; field-extracted to dodge struct-return) ---
	{"STRUCT(1 AS x, 2 AS y)", "SELECT STRUCT(1 AS x, 2 AS y).x", true},
	{"STRUCT<x INT64, y STRING>(1, 'a')", "SELECT STRUCT<x INT64, y STRING>(1, 'a').x", true},
	{"STRUCT()", "SELECT STRUCT().field1", true}, // empty struct constructor
	{"(1, 2, 3)", "SELECT (1, 2, 3).field1", true},

	// --- access (accept) ---
	{"a[0]", "SELECT a[0] FROM t", true},
	{"a[OFFSET(0)]", "SELECT a[OFFSET(0)] FROM t", true},
	{"a[ORDINAL(1)]", "SELECT a[ORDINAL(1)] FROM t", true},
	{"a[SAFE_OFFSET(0)]", "SELECT a[SAFE_OFFSET(0)] FROM t", true},
	{"f(x).y", "SELECT f(x).y FROM t", true},
	{"a[0].b", "SELECT a[0].b FROM t", true},

	// --- INTERVAL (accept) ---
	{"INTERVAL 5 DAY", "SELECT INTERVAL 5 DAY", true},
	{"INTERVAL x YEAR TO MONTH", "SELECT INTERVAL x YEAR TO MONTH FROM t", true},

	// --- subqueries (accept) ---
	{"EXISTS(SELECT 1)", "SELECT EXISTS(SELECT 1)", true},
	{"ARRAY(SELECT 1)", "SELECT ARRAY(SELECT 1)", true},
	{"(SELECT 1) + 2", "SELECT (SELECT 1) + 2", true},

	// =====================================================================
	// NEGATIVES — all reject (syntax)
	// =====================================================================
	{"1 +", "SELECT 1 +", false},
	{"1 2", "SELECT 1 2", false},
	{"a BETWEEN 1", "SELECT a BETWEEN 1 FROM t", false},
	{"a BETWEEN 1 OR 2", "SELECT a BETWEEN 1 OR 2 FROM t", false},
	{"a = b = c", "SELECT a = b = c FROM t", false},
	{"1 < 2 < 3", "SELECT 1 < 2 < 3", false},
	{"a = b IS NULL", "SELECT a = b IS NULL FROM t", false},
	{"a IN (1) IN (2)", "SELECT a IN (1) IN (2) FROM t", false},
	{"1 IS NULL IS NULL", "SELECT 1 IS NULL IS NULL", false},
	{"1 LIKE 'a' LIKE 'b'", "SELECT 1 LIKE 'a' LIKE 'b'", false},
	// quantified comparison negatives (oracle: syntax reject).
	{"x = ANY ()", "SELECT * FROM t WHERE x = ANY ()", false},                                       // empty list
	{"x = ANY ANY (SELECT v FROM s)", "SELECT * FROM t WHERE x = ANY ANY (SELECT v FROM s)", false}, // double quantifier
	{"x = ANY SELECT v FROM s", "SELECT * FROM t WHERE x = ANY SELECT v FROM s", false},             // subquery not parenthesized
	{"x IS ANY (SELECT v FROM s)", "SELECT * FROM t WHERE x IS ANY (SELECT v FROM s)", false},       // IS has no quantified form
	{"x = ANY (SELECT v FROM s) = y", "SELECT * FROM t WHERE x = ANY (SELECT v FROM s) = y", false}, // non-associative chaining
	{"1 BETWEEN 0 AND 2 BETWEEN 0 AND 1", "SELECT 1 BETWEEN 0 AND 2 BETWEEN 0 AND 1", false},
	{"CASE END", "SELECT CASE END", false},
	{"CAST(x INT64)", "SELECT CAST(x INT64) FROM t", false},
	{"CAST(x AS)", "SELECT CAST(x AS) FROM t", false},
	{"CAST(CAST AS INT64)", "SELECT CAST(CAST AS INT64)", false},
	{"EXTRACT(x y)", "SELECT EXTRACT(x y) FROM t", false},
	{"a IN ()", "SELECT a IN () FROM t", false},
	{"a..b", "SELECT a..b FROM t", false},
	{"* 5", "SELECT * 5", false},
	{"EXTRACT(HOUR FROM ts @ TIME ZONE 'UTC')", "SELECT EXTRACT(HOUR FROM ts @ TIME ZONE 'UTC') FROM t", false},
}

func TestExprDifferential(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set; skipping live differential")
	}
	h := newExprHarness(t)
	defer h.close()

	for _, fx := range exprFixtures {
		fx := fx
		t.Run(fx.expr, func(t *testing.T) {
			// 1. Oracle verdict for the embedded statement.
			v := h.verdict(t, fx.stmt)
			switch v.Verdict {
			case "accept", "reject":
				// proceed
			case "error":
				t.Fatalf("oracle returned ERROR (no verdict) for %q: reason=%s code=%s msg=%s\n"+
					"the oracle could not decide; fix the wrapper/emulator — do NOT treat as accept/reject",
					fx.stmt, v.Reason, v.Code, v.Message)
			default:
				t.Fatalf("unexpected harness verdict %q for %q", v.Verdict, fx.stmt)
			}
			oracleAccepts := v.Verdict == "accept"

			// Sanity: the asserted wantParse must match the live oracle.
			if oracleAccepts != fx.wantParse {
				t.Fatalf("fixture wantParse=%v but oracle says %q (%s) for %q",
					fx.wantParse, v.Verdict, v.Message, fx.stmt)
			}

			// 2. omni ParseExpression verdict on the bare expression.
			_, errs := ParseExpression(fx.expr)
			omniAccepts := len(errs) == 0

			if omniAccepts != oracleAccepts {
				t.Errorf("DIVERGENCE on %q: omni accepts=%v, oracle accepts=%v (%s: %s)",
					fx.expr, omniAccepts, oracleAccepts, v.Reason, v.Message)
			}
		})
	}
}

// --- harness plumbing (one persistent batch process, mirrors datatypes_oracle_test) ---

type exprHarnessVerdict struct {
	Verdict string `json:"verdict"`
	Kind    string `json:"kind"`
	Reason  string `json:"reason"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type exprHarness struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	mu     sync.Mutex
}

func newExprHarness(t *testing.T) *exprHarness {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	projDir := filepath.Join(repoRoot, "harness", "googlesql-spanner")
	if _, err := os.Stat(projDir); err != nil {
		t.Skipf("harness project not found at %s", projDir)
	}

	bin := filepath.Join(projDir, "googlesql-spanner")
	if _, err := os.Stat(bin); err != nil {
		build := exec.Command("go", "build", "-o", bin, ".")
		build.Dir = projDir
		if out, err := build.CombinedOutput(); err != nil {
			t.Fatalf("building harness failed: %v\n%s", err, out)
		}
	}

	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(), "GOOGLESQL_HARNESS_LINE=1")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("stdin pipe: %v", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("starting harness: %v", err)
	}
	return &exprHarness{cmd: cmd, stdin: stdin, stdout: bufio.NewReader(stdoutPipe)}
}

func (h *exprHarness) verdict(t *testing.T, stmt string) exprHarnessVerdict {
	t.Helper()
	h.mu.Lock()
	defer h.mu.Unlock()
	enc := base64.StdEncoding.EncodeToString([]byte(stmt))
	if _, err := fmt.Fprintln(h.stdin, enc); err != nil {
		t.Fatalf("writing to harness: %v", err)
	}
	line, err := h.stdout.ReadString('\n')
	if err != nil {
		t.Fatalf("reading harness verdict: %v", err)
	}
	var v exprHarnessVerdict
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &v); err != nil {
		t.Fatalf("decoding harness verdict %q: %v", line, err)
	}
	return v
}

func (h *exprHarness) close() {
	if h.stdin != nil {
		_ = h.stdin.Close()
	}
	if h.cmd != nil && h.cmd.Process != nil {
		_ = h.cmd.Wait()
	}
}
