//go:build googlesql_oracle

// Differential test for the `parser-query-clauses` node against the live Cloud
// Spanner emulator oracle. Run with:
//
//	SPANNER_EMULATOR_HOST=localhost:9010 \
//	  go test -tags googlesql_oracle ./googlesql/parser/ -run TestQueryClausesDifferential
//
// It is the PROVE gate for googlesql/parser-query-clauses (correctness-
// protocol.md): for every fixture it (1) feeds the full query to the emulator
// via the harness/googlesql-spanner CLI and reads the accept/reject verdict,
// then (2) parses the same query with omni's Parse, and asserts the two agree —
// BOTH polarities (the corpus has accept AND reject cases). A harness `error`
// verdict (oracle could not decide) fails the fixture loudly; it is never folded
// into accept or reject.
//
// SCOPE — what is oracle-authoritative here. Although PIVOT/UNPIVOT are
// documented as BigQuery-only, the live Spanner emulator's ZetaSQL grammar
// PARSES them (then feature-rejects: "PIVOT is not supported" / "... is not
// allowed with array scans") — so the harness classifies those as ACCEPT (the
// grammar accepted), and the emulator IS authoritative for PIVOT / UNPIVOT /
// TABLESAMPLE / FOR SYSTEM_TIME in BOTH polarities. The hard SYNTAX rejects
// (missing FOR, empty IN, missing TABLESAMPLE unit, double pivot, trailing
// sample alias) carry the canonical `Syntax error:` prefix and so verdict
// reject. This file drives all four constructs except:
//
// EXCLUDED — the SELECT-level differential-privacy clause
// (`SELECT WITH DIFFERENTIAL_PRIVACY OPTIONS(...)`). It is BigQuery-only and the
// emulator reports "Unexpected keyword WITH" WITHOUT the `Syntax error:` prefix,
// so the fail-closed harness misclassifies it as a (semantic) ACCEPT — the
// oracle is NON-authoritative for it (oracle.md dialect caveat). It is proven by
// the hand-written unit tests (pivot_unpivot_test.go) + truth1/bigquery
// triangulation instead, mirroring parser-select's exclusion of DP fixtures.
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

// qcFixture is a full query statement (no trailing ';') fed to BOTH the oracle
// and omni's Parse. wantParse is the expected omni outcome and MUST equal the
// emulator's grammar verdict.
type qcFixture struct {
	sql       string
	wantParse bool
}

var queryClausesFixtures = []qcFixture{
	// ===================== ACCEPT — PIVOT =========================
	{"SELECT * FROM Produce PIVOT(SUM(sales) FOR quarter IN ('Q1', 'Q2', 'Q3', 'Q4'))", true},
	{"SELECT * FROM t PIVOT(SUM(s) AS tot, COUNT(*) AS n FOR q IN ('Q1' AS first, 'Q2'))", true},
	{"SELECT * FROM t PIVOT(SUM(s) FOR q IN ('Q1')) AS p", true},
	{"SELECT * FROM t PIVOT(SUM(s) FOR q IN ('Q1')) p", true},
	{"SELECT * FROM t AS x PIVOT(SUM(s) FOR q IN ('a'))", true},
	{"SELECT * FROM t x PIVOT(SUM(s) FOR q IN ('a'))", true},
	// PIVOT / UNPIVOT are non-reserved: a bare trailing keyword is an implicit
	// table alias, not a (malformed) operator.
	{"SELECT * FROM t pivot", true},
	{"SELECT * FROM t unpivot", true},
	{"SELECT * FROM t pivot, s", true},
	{"SELECT * FROM t AS x PIVOT(SUM(s) FOR q IN ('a')) AS p", true},
	{"SELECT * FROM (SELECT 1 AS s, 2 AS q) PIVOT(SUM(s) FOR q IN (2))", true},
	{"SELECT * FROM (SELECT 1 AS s, 2 AS q) AS t PIVOT(SUM(s) FOR q IN (2))", true},
	// The FOR input is parsed at higher-than-AND precedence so the trailing IN is
	// the value list, not an `expr IN (...)` predicate: arithmetic / field access
	// in FOR accepts; a top-level AND in FOR rejects.
	{"SELECT * FROM t PIVOT(SUM(s) FOR a + b IN ('x'))", true},
	{"SELECT * FROM t PIVOT(SUM(s) FOR a AND b IN ('x'))", false},

	// ===================== ACCEPT — UNPIVOT =======================
	{"SELECT * FROM sales_table UNPIVOT(total_sales FOR quarter IN (Q1, Q2, Q3, Q4))", true},
	{"SELECT * FROM t UNPIVOT INCLUDE NULLS (total_sales FOR quarter IN (Q1, Q2))", true},
	{"SELECT * FROM t UNPIVOT EXCLUDE NULLS (sales FOR q IN (Q1, Q2))", true},
	{"SELECT * FROM t UNPIVOT ((s1, s2) FOR q IN ((Q1, Q2), (Q3, Q4)))", true},
	{"SELECT * FROM t UNPIVOT (s FOR q IN (Q1 AS 'a', Q2 AS 2)) AS u", true},
	{"SELECT * FROM t AS x UNPIVOT(v FOR n IN (a, b)) AS u", true},

	// ===================== ACCEPT — TABLESAMPLE ===================
	{"SELECT * FROM Singers TABLESAMPLE BERNOULLI (10 PERCENT)", true},
	{"SELECT * FROM Albums TABLESAMPLE RESERVOIR (100 ROWS)", true},
	{"SELECT * FROM t TABLESAMPLE SYSTEM (10 PERCENT)", true},
	{"SELECT * FROM t TABLESAMPLE BERNOULLI (5 PERCENT) REPEATABLE (42)", true},
	{"SELECT * FROM t TABLESAMPLE RESERVOIR (100 ROWS) WITH WEIGHT", true},
	{"SELECT * FROM t TABLESAMPLE RESERVOIR (100 ROWS) WITH WEIGHT AS w", true},
	{"SELECT * FROM t AS x TABLESAMPLE BERNOULLI (10 PERCENT)", true},
	{"SELECT * FROM (SELECT 1 AS s) TABLESAMPLE BERNOULLI (10 PERCENT)", true},
	// PIVOT then TABLESAMPLE: sample is the outermost suffix.
	{"SELECT * FROM t PIVOT(SUM(s) FOR q IN ('a')) TABLESAMPLE BERNOULLI (10 PERCENT)", true},
	// TABLESAMPLE then a trailing clause (WHERE) — sample is complete.
	{"SELECT * FROM t TABLESAMPLE BERNOULLI (10 PERCENT) WHERE x > 1", true},

	// ===================== ACCEPT — FOR SYSTEM_TIME ===============
	{"SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01'", true},
	{"SELECT * FROM t FOR SYSTEM_TIME AS OF TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 HOUR)", true},
	{"SELECT * FROM t FOR SYSTEM TIME AS OF '2020-01-01'", true},
	{"SELECT * FROM t a FOR SYSTEM_TIME AS OF '2020-01-01'", true},

	// ===================== REJECT — hard syntax errors ============
	// PIVOT: missing aggregate / empty IN / missing FOR.
	{"SELECT * FROM t PIVOT(FOR q IN ('Q1'))", false},
	{"SELECT * FROM t PIVOT(SUM(s) FOR q IN ())", false},
	// UNPIVOT: empty body / empty IN.
	{"SELECT * FROM t UNPIVOT()", false},
	{"SELECT * FROM t UNPIVOT(s FOR q IN ())", false},
	// At most one pivot/unpivot per source.
	{"SELECT * FROM t PIVOT(SUM(s) FOR q IN ('Q1')) UNPIVOT(x FOR y IN (a))", false},
	// TABLESAMPLE: missing unit / no method paren / trailing alias.
	{"SELECT * FROM t TABLESAMPLE BERNOULLI (10)", false},
	{"SELECT * FROM t TABLESAMPLE (10 PERCENT)", false},
	{"SELECT * FROM t TABLESAMPLE BERNOULLI (10 PERCENT) AS x", false},
	// TABLESAMPLE size is a literal/cast/param, NOT an arbitrary expression.
	{"SELECT * FROM t TABLESAMPLE BERNOULLI (1 + 1 ROWS)", false},
	{"SELECT * FROM t TABLESAMPLE BERNOULLI (x ROWS)", false},
	// PIVOT FOR input is higher-than-AND precedence: a top-level comparison / NOT
	// in FOR is a syntax error (the IN is the pivot value list, not a predicate).
	{"SELECT * FROM t PIVOT(SUM(s) FOR a = b IN ('x'))", false},
	{"SELECT * FROM t PIVOT(SUM(s) FOR NOT flag IN ('x'))", false},
	// A bare (unparenthesized) multi-column UNPIVOT value list rejects.
	{"SELECT * FROM t UNPIVOT(s1, s2 FOR q IN ((Q1, Q2)))", false},
	// TABLESAMPLE before PIVOT is rejected (sample is outermost).
	{"SELECT * FROM t TABLESAMPLE BERNOULLI (10 PERCENT) PIVOT(SUM(s) FOR q IN ('a'))", false},
	// PIVOT alias may not carry a column list.
	{"SELECT * FROM t PIVOT(SUM(s) FOR q IN ('Q1')) AS p (x, y)", false},
	// FOR SYSTEM (two-word) requires the TIME keyword.
	{"SELECT * FROM t FOR SYSTEM AS OF '2020-01-01'", false},
}

func TestQueryClausesDifferential(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set; skipping live differential")
	}
	h := newQCHarness(t)
	defer h.close()

	for _, fx := range queryClausesFixtures {
		fx := fx
		t.Run(fx.sql, func(t *testing.T) {
			v := h.verdict(t, fx.sql)
			switch v.Verdict {
			case "accept", "reject":
				// proceed
			case "error":
				t.Fatalf("oracle returned ERROR (no verdict) for %q: reason=%s code=%s msg=%s\n"+
					"the oracle could not decide; do NOT treat as accept/reject",
					fx.sql, v.Reason, v.Code, v.Message)
			default:
				t.Fatalf("unexpected harness verdict %q for %q", v.Verdict, fx.sql)
			}
			oracleAccepts := v.Verdict == "accept"

			if oracleAccepts != fx.wantParse {
				t.Fatalf("fixture wantParse=%v but oracle says %q (%s) for %q",
					fx.wantParse, v.Verdict, v.Message, fx.sql)
			}

			_, errs := Parse(fx.sql)
			omniAccepts := len(errs) == 0

			if omniAccepts != oracleAccepts {
				t.Errorf("DIVERGENCE on %q: omni accepts=%v, oracle accepts=%v (%s: %s)",
					fx.sql, omniAccepts, oracleAccepts, v.Reason, v.Message)
			}
		})
	}
}

// --- harness plumbing (one persistent batch process, mirrors select_oracle_test) ---

type qcHarnessVerdict struct {
	Verdict string `json:"verdict"`
	Kind    string `json:"kind"`
	Reason  string `json:"reason"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type qcHarness struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	mu     sync.Mutex
}

func newQCHarness(t *testing.T) *qcHarness {
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
	return &qcHarness{cmd: cmd, stdin: stdin, stdout: bufio.NewReader(stdoutPipe)}
}

func (h *qcHarness) verdict(t *testing.T, stmt string) qcHarnessVerdict {
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
	var v qcHarnessVerdict
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &v); err != nil {
		t.Fatalf("decoding harness verdict %q: %v", line, err)
	}
	return v
}

func (h *qcHarness) close() {
	if h.stdin != nil {
		_ = h.stdin.Close()
	}
	if h.cmd != nil && h.cmd.Process != nil {
		_ = h.cmd.Wait()
	}
}
