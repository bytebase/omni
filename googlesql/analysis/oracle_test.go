//go:build googlesql_oracle

// Differential test for the googlesql/analysis node against the live Cloud
// Spanner emulator oracle. Run with:
//
//	SPANNER_EMULATOR_HOST=localhost:9010 \
//	  go test -tags googlesql_oracle ./googlesql/analysis/ -run TestAnalysisDifferential
//
// analysis is a FEATURE node, so the oracle is NOT used to adjudicate a grammar
// accept/reject — that is the parser nodes' job (parser/*_oracle_test.go). What
// the oracle pins HERE is that the SQL the classify/span tests assert over is
// REAL GoogleSQL syntax the grammar accepts: if a corpus statement the tests
// treat as a valid query were actually rejected by the emulator, the
// classification/lineage assertions built on it would be meaningless. The oracle
// therefore confirms the grammar verdict for the whole classify corpus and
// cross-checks the read-only classification against the statement's true nature.
//
// The emulator speaks Spanner's dialect, a SUBSET of the BigQuery+Spanner union
// (oracle.md). The corpus below is restricted to forms whose Spanner GRAMMAR
// verdict is authoritative: shared GoogleSQL core (SELECT/FROM/JOIN/CTE/set-op),
// DML (INSERT/UPDATE/DELETE), and DDL (CREATE/DROP/ALTER) — all of which the
// emulator parses. BigQuery-only forms whose Spanner verdict is non-authoritative
// (dashed paths, MERGE, TRUNCATE) are covered by the hand-written unit tests +
// the divergence ledger instead, NOT here.
//
// The harness classifies a statement that PARSES but is then feature-rejected
// ("X is not supported") as ACCEPT (message-prefix rule, oracle.md), so a
// union-grammar form that Spanner feature-rejects is still a valid accept fixture.
package analysis

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

	"github.com/bytebase/omni/googlesql/ast"
	"github.com/bytebase/omni/googlesql/parser"
)

// parseFirstStmt parses sql and returns its first statement node, failing the
// test if nothing parsed (the corpus is curated to parse cleanly).
func parseFirstStmt(t *testing.T, sql string) ast.Node {
	t.Helper()
	file, _ := parser.Parse(sql)
	if file == nil || len(file.Stmts) == 0 {
		t.Fatalf("no statement parsed from %q", sql)
	}
	return file.Stmts[0]
}

// classifyCorpus pairs each statement with the QueryType the analysis classifier
// must report (in the dialect noted) and whether it is read-only. Every entry is
// expected to be ACCEPTED by the emulator's grammar (wantAccept implicitly true);
// the differential fails loudly if the oracle rejects one, because that would
// mean the classification assertion rests on invalid syntax.
var classifyCorpus = []struct {
	sql        string
	dialect    Dialect
	wantType   QueryType
	isReadOnly bool
}{
	// --- read-only queries ---
	{"SELECT a FROM t", DialectSpanner, Select, true},
	{"SELECT * FROM t1 JOIN t2 ON t1.id = t2.id", DialectSpanner, Select, true},
	{"WITH c AS (SELECT a FROM t) SELECT a FROM c", DialectSpanner, Select, true},
	{"SELECT a FROM t WHERE b IN (SELECT b FROM t2)", DialectSpanner, Select, true},
	{"SELECT a FROM t UNION ALL SELECT a2 FROM t2", DialectSpanner, Select, true},
	{"SELECT * FROM INFORMATION_SCHEMA.TABLES", DialectSpanner, SelectInfoSchema, true},
	{"SELECT * FROM SPANNER_SYS.QUERY_STATS_TOP_MINUTE", DialectSpanner, SelectInfoSchema, true},
	// --- data-changing (DML) ---
	{"INSERT INTO t (a) VALUES (1)", DialectSpanner, DML, false},
	{"DELETE FROM t WHERE a = 1", DialectSpanner, DML, false},
	{"UPDATE t SET a = 2 WHERE a = 1", DialectSpanner, DML, false},
	// --- schema-changing (DDL) ---
	{"CREATE TABLE t3 (a INT64) PRIMARY KEY (a)", DialectSpanner, DDL, false},
	{"DROP TABLE t3", DialectSpanner, DDL, false},
	{"ALTER TABLE t ADD COLUMN c INT64", DialectSpanner, DDL, false},
}

// TestAnalysisDifferential confirms each corpus statement's grammar verdict
// against the live emulator and that the analysis classifier's read-only verdict
// agrees with the statement's true (read-only vs data/schema-changing) nature.
func TestAnalysisDifferential(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set; skipping live differential")
	}
	h := newAnalysisHarness(t)
	defer h.close()

	for _, tc := range classifyCorpus {
		tc := tc
		t.Run(tc.sql, func(t *testing.T) {
			v := h.verdict(t, tc.sql)
			switch v.Verdict {
			case "accept", "reject":
			case "error":
				t.Fatalf("oracle returned ERROR (no verdict) for %q: reason=%s code=%s msg=%s\n"+
					"the oracle could not decide; fix the wrapper/emulator — do NOT treat as accept/reject",
					tc.sql, v.Reason, v.Code, v.Message)
			default:
				t.Fatalf("unexpected harness verdict %q for %q", v.Verdict, tc.sql)
			}
			if v.Verdict != "accept" {
				t.Fatalf("emulator REJECTED corpus statement %q (%s: %s) — the analysis assertions on it rest on invalid syntax",
					tc.sql, v.Reason, v.Message)
			}

			// Cross-check: the classifier puts the statement on the correct side of
			// the read-only line, and reports the expected type.
			gotType := Classify(parseFirstStmt(t, tc.sql), tc.dialect)
			if gotType != tc.wantType {
				t.Errorf("Classify(%q) = %v, want %v", tc.sql, gotType, tc.wantType)
			}
			if got := isReadOnly(gotType); got != tc.isReadOnly {
				t.Errorf("isReadOnly(%q) = %v (type %v), want %v", tc.sql, got, gotType, tc.isReadOnly)
			}
		})
	}
}

// isReadOnly mirrors the bytebase read-only guard: a statement is read-only iff
// it is a Select, Explain, or SelectInfoSchema.
func isReadOnly(qt QueryType) bool {
	return qt == Select || qt == Explain || qt == SelectInfoSchema
}

// --- harness plumbing (one persistent batch process; mirrors the parser
// package's *_oracle_test.go harness) ---

type analysisHarnessVerdict struct {
	Verdict string `json:"verdict"`
	Kind    string `json:"kind"`
	Reason  string `json:"reason"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type analysisHarness struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	mu     sync.Mutex
}

func newAnalysisHarness(t *testing.T) *analysisHarness {
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
	return &analysisHarness{cmd: cmd, stdin: stdin, stdout: bufio.NewReader(stdoutPipe)}
}

func (h *analysisHarness) verdict(t *testing.T, stmt string) analysisHarnessVerdict {
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
	var v analysisHarnessVerdict
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &v); err != nil {
		t.Fatalf("decoding harness verdict %q: %v", line, err)
	}
	return v
}

func (h *analysisHarness) close() {
	if h.stdin != nil {
		_ = h.stdin.Close()
	}
	if h.cmd != nil && h.cmd.Process != nil {
		_ = h.cmd.Wait()
	}
}
