//go:build googlesql_oracle

// Differential test for the `parser-select` node against the live Cloud Spanner
// emulator oracle. Run with:
//
//	SPANNER_EMULATOR_HOST=localhost:9010 \
//	  go test -tags googlesql_oracle ./googlesql/parser/ -run TestSelectDifferential
//
// It is the PROVE gate for googlesql/parser-select (correctness-protocol.md):
// for every fixture it (1) feeds the full query to the emulator via the
// harness/googlesql-spanner CLI and reads the accept/reject verdict, then (2)
// parses the same query with omni's Parse, and asserts the two agree — BOTH
// polarities (the corpus has accept AND reject cases). A harness `error` verdict
// (oracle could not decide) fails the fixture loudly; it is never folded into
// accept or reject.
//
// The emulator speaks Spanner's dialect, a SUBSET of the BigQuery+Spanner union
// the parser must accept (oracle.md). Every form below is one of:
//   - SHARED GoogleSQL core (SELECT/FROM/joins/set-ops/CTE/GROUP BY/HAVING/
//     WINDOW/ORDER/LIMIT/UNNEST) — the Spanner verdict is authoritative.
//   - A union-grammar form that PARSES on Spanner but is then feature-rejected
//     ("QUALIFY is not supported", "RECURSIVE is not supported", "CORRESPONDING
//     ... is not supported") — the harness classifies these ACCEPT (the grammar
//     accepted) per its message-prefix rule, so they are valid accept fixtures
//     that exercise QUALIFY / WITH RECURSIVE / set-op CORRESPONDING.
//
// BigQuery-only forms whose Spanner verdict is NON-authoritative are EXCLUDED
// from this differential and covered by the hand-written unit tests
// (select_test.go) + the divergence ledger instead:
//   - unquoted dashed table paths (`my-project.ds.tbl`) — Spanner SYNTAX-rejects
//     them ("Table name contains '-'"); divergence ledger #85. They are
//     BigQuery-valid (the grammar's dashed_path_expression) and omni accepts
//     them; the Spanner reject would be a false DIVERGENCE here.
//   - `SELECT WITH AGGREGATION_THRESHOLD OPTIONS(...)` — BigQuery differential-
//     privacy clause; the emulator's classification of it is unreliable.
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

// selectFixture is a full query statement (no trailing ';') fed to BOTH the
// oracle and omni's Parse. wantParse is the expected omni outcome and MUST equal
// the emulator's grammar verdict.
type selectFixture struct {
	sql       string
	wantParse bool
}

var selectFixtures = []selectFixture{
	// ===================== ACCEPT — SELECT core =========================
	{"SELECT 1", true},
	{"SELECT 1 AS n", true},
	{"SELECT a, b FROM t", true},
	{"SELECT * FROM t", true},
	{"SELECT a AS x, b y FROM t", true},     // explicit + implicit alias
	{"SELECT DISTINCT a FROM t", true},
	{"SELECT ALL a FROM t", true},
	{"SELECT AS STRUCT a, b FROM t", true},
	{"SELECT AS VALUE a FROM t", true},
	{"SELECT t.* FROM t", true},
	{"SELECT t.* EXCEPT (a) FROM t", true},
	{"SELECT * EXCEPT (a, b) FROM t", true},
	{"SELECT * REPLACE (a + 1 AS a) FROM t", true},

	// ===================== ACCEPT — FROM / joins ========================
	{"SELECT * FROM s.t", true},
	{"SELECT * FROM s.t AS x", true},
	{"SELECT * FROM s.t x", true},
	{"SELECT * FROM (SELECT 1 AS n) AS sub", true},
	{"SELECT * FROM a, b", true},
	{"SELECT * FROM a JOIN b ON a.x = b.x", true},
	{"SELECT * FROM a JOIN b", true}, // no ON/USING — oracle-confirmed accept
	{"SELECT * FROM a INNER JOIN b ON a.x = b.x", true},
	{"SELECT * FROM a LEFT JOIN b USING (x)", true},
	{"SELECT * FROM a LEFT OUTER JOIN b USING (x)", true},
	{"SELECT * FROM a RIGHT JOIN b ON a.x = b.x", true},
	{"SELECT * FROM a FULL OUTER JOIN b ON a.x = b.x", true},
	{"SELECT * FROM a CROSS JOIN b", true},
	{"SELECT * FROM a HASH JOIN b ON a.x = b.x", true},
	{"SELECT * FROM a JOIN b ON a.x = b.x JOIN c ON b.y = c.y", true},
	{"SELECT * FROM (a JOIN b ON a.x = b.x)", true},
	{"SELECT * FROM UNNEST([1, 2, 3]) AS num", true},
	{"SELECT num, pos FROM UNNEST([1, 2]) AS num WITH OFFSET AS pos", true},

	// ===================== ACCEPT — WHERE/GROUP/HAVING ==================
	{"SELECT * FROM t WHERE a > 1", true},
	{"SELECT * FROM t WHERE a > 1 AND b < 2", true},
	{"SELECT a, COUNT(*) FROM t GROUP BY a", true},
	{"SELECT a FROM t GROUP BY 1", true},
	{"SELECT a FROM t GROUP BY a, b", true},
	{"SELECT SUM(x) FROM t GROUP BY ROLLUP (a, b)", true},
	{"SELECT a, COUNT(*) c FROM t GROUP BY a HAVING c > 1", true},

	// ===================== ACCEPT — WINDOW / ORDER / LIMIT ==============
	{"SELECT SUM(x) OVER w FROM t WINDOW w AS (PARTITION BY a ORDER BY b)", true},
	{"SELECT x FROM t ORDER BY x", true},
	{"SELECT x FROM t ORDER BY x DESC NULLS LAST", true},
	{"SELECT x FROM t ORDER BY x ASC, y DESC", true},
	{"SELECT x FROM t LIMIT 10", true},
	{"SELECT x FROM t LIMIT 10 OFFSET 5", true},
	{"SELECT x FROM t ORDER BY x LIMIT 10 OFFSET 5", true},

	// ===================== ACCEPT — set-ops =============================
	{"SELECT a FROM t UNION ALL SELECT a FROM u", true},
	{"SELECT a FROM t UNION DISTINCT SELECT a FROM u", true},
	{"SELECT a FROM t INTERSECT DISTINCT SELECT a FROM u", true},
	{"SELECT a FROM t EXCEPT DISTINCT SELECT a FROM u", true},
	{"SELECT a FROM t UNION ALL SELECT a FROM u ORDER BY a", true},
	{"(SELECT 1 AS a) UNION ALL (SELECT 2 AS a)", true},
	// Uniform-op flat chain: same operation throughout is accepted.
	{"SELECT 1 AS a UNION ALL SELECT 2 AS a UNION ALL SELECT 3 AS a", true},
	{"SELECT 1 AS a INTERSECT DISTINCT SELECT 2 AS a INTERSECT DISTINCT SELECT 3 AS a", true},
	// Mixed-op chain requires parentheses for grouping.
	{"(SELECT 1 AS a UNION ALL SELECT 2 AS a) INTERSECT DISTINCT SELECT 3 AS a", true},

	// ===================== ACCEPT — CTE =================================
	{"WITH c AS (SELECT 1 AS n) SELECT n FROM c", true},
	{"WITH a AS (SELECT 1 AS n), b AS (SELECT n + 1 AS n FROM a) SELECT * FROM b", true},
	{"WITH c AS (SELECT 1 AS n) SELECT * FROM c UNION ALL SELECT 2", true},

	// ===================== ACCEPT — subqueries ==========================
	{"SELECT (SELECT MAX(x) FROM u) AS m FROM t", true},
	{"SELECT * FROM t WHERE EXISTS (SELECT 1 FROM u WHERE u.id = t.id)", true},
	{"SELECT * FROM t WHERE id IN (SELECT id FROM u)", true},

	// ===================== ACCEPT — Spanner-only ========================
	{"SELECT * FROM t WHERE id = 1 FOR UPDATE", true},
	{"@{OPTIMIZER_VERSION=2} SELECT * FROM t", true}, // statement-level hint

	// ===================== ACCEPT — union-grammar, Spanner feature-rejects
	// (grammar parses => harness classifies ACCEPT; exercises QUALIFY / WITH
	// RECURSIVE / set-op CORRESPONDING which BigQuery supports).
	{"SELECT a FROM t QUALIFY ROW_NUMBER() OVER (ORDER BY a) = 1", true},
	{"WITH RECURSIVE c AS ((SELECT 1 AS n) UNION ALL (SELECT n + 1 FROM c WHERE n < 3)) SELECT * FROM c", true},
	{"SELECT 1 AS a UNION ALL STRICT CORRESPONDING SELECT 2 AS a", true},
	{"SELECT 1 AS a FULL OUTER UNION ALL SELECT 2 AS a", true},

	// ===================== NEGATIVES — syntax reject ====================
	{"SELECT", false},                          // empty select list
	{"SELECT FROM t", false},                   // empty list before FROM
	{"SELECT 1 2 3", false},                    // double implicit alias
	{"SELECT * FROM", false},                   // FROM with no source
	{"SELECT * FROM t WHERE", false},           // WHERE with no expr
	{"SELECT * FROM a JOIN", false},            // JOIN with no right source
	{"SELECT * FROM a CROSS b", false},         // CROSS without JOIN
	{"SELECT * FROM a USING (x)", false},       // USING with no join
	{"SELECT * FROM t UNION SELECT 1", false},  // set-op missing ALL/DISTINCT
	{"SELECT * FROM t UNION", false},           // set-op with no right query
	// Mixed set operations in a flat chain require parentheses.
	{"SELECT 1 AS a UNION ALL SELECT 2 AS a INTERSECT DISTINCT SELECT 3 AS a", false},
	{"SELECT 1 AS a UNION ALL SELECT 2 AS a UNION DISTINCT SELECT 3 AS a", false},
	{"FROM t", false},                          // FROM-first query
	{"WITH c AS (SELECT 1)", false},            // WITH with no trailing body
	{"WITH c SELECT 1", false},                 // CTE missing AS (query)
	{"SELECT a, FROM", false},                  // trailing comma then nothing
	{"SELECT * FROM t ORDER", false},           // ORDER without BY
	{"SELECT * FROM t GROUP BY", false},        // GROUP BY with no items
	{"SELECT item OVER (w) FROM t", false},     // OVER must follow a function call
}

func TestSelectDifferential(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set; skipping live differential")
	}
	h := newSelectHarness(t)
	defer h.close()

	for _, fx := range selectFixtures {
		fx := fx
		t.Run(fx.sql, func(t *testing.T) {
			// 1. Oracle verdict for the query.
			v := h.verdict(t, fx.sql)
			switch v.Verdict {
			case "accept", "reject":
				// proceed
			case "error":
				t.Fatalf("oracle returned ERROR (no verdict) for %q: reason=%s code=%s msg=%s\n"+
					"the oracle could not decide; fix the wrapper/emulator — do NOT treat as accept/reject",
					fx.sql, v.Reason, v.Code, v.Message)
			default:
				t.Fatalf("unexpected harness verdict %q for %q", v.Verdict, fx.sql)
			}
			oracleAccepts := v.Verdict == "accept"

			// Sanity: the asserted wantParse must match the live oracle.
			if oracleAccepts != fx.wantParse {
				t.Fatalf("fixture wantParse=%v but oracle says %q (%s) for %q",
					fx.wantParse, v.Verdict, v.Message, fx.sql)
			}

			// 2. omni Parse verdict.
			_, errs := Parse(fx.sql)
			omniAccepts := len(errs) == 0

			if omniAccepts != oracleAccepts {
				t.Errorf("DIVERGENCE on %q: omni accepts=%v, oracle accepts=%v (%s: %s)",
					fx.sql, omniAccepts, oracleAccepts, v.Reason, v.Message)
			}
		})
	}
}

// --- harness plumbing (one persistent batch process, mirrors expr_oracle_test) ---

type selectHarnessVerdict struct {
	Verdict string `json:"verdict"`
	Kind    string `json:"kind"`
	Reason  string `json:"reason"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type selectHarness struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	mu     sync.Mutex
}

func newSelectHarness(t *testing.T) *selectHarness {
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
	return &selectHarness{cmd: cmd, stdin: stdin, stdout: bufio.NewReader(stdoutPipe)}
}

func (h *selectHarness) verdict(t *testing.T, stmt string) selectHarnessVerdict {
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
	var v selectHarnessVerdict
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &v); err != nil {
		t.Fatalf("decoding harness verdict %q: %v", line, err)
	}
	return v
}

func (h *selectHarness) close() {
	if h.stdin != nil {
		_ = h.stdin.Close()
	}
	if h.cmd != nil && h.cmd.Process != nil {
		_ = h.cmd.Wait()
	}
}
