//go:build googlesql_oracle

// Differential test for the `types` node against the live Cloud Spanner
// emulator oracle. Run with:
//
//	SPANNER_EMULATOR_HOST=localhost:9010 \
//	  go test -tags googlesql_oracle ./googlesql/parser/ -run TestTypeDifferential
//
// It is the PROVE gate for googlesql/types (correctness-protocol.md): for every
// fixture it (1) feeds a full statement embedding the type to the emulator via
// the harness/googlesql-spanner CLI and reads the accept/reject verdict, then
// (2) parses the bare type text with omni's parseType, and asserts the two
// agree — BOTH polarities (the corpus has accept AND reject cases). A harness
// `error` verdict (oracle could not decide) fails the fixture loudly; it is
// never folded into accept or reject.
//
// The emulator speaks Spanner's dialect, a SUBSET of the BigQuery+Spanner union
// the parser must accept (oracle.md). Every type form covered here is in the
// SHARED GoogleSQL core (scalars, ARRAY/STRUCT/RANGE/MAP/FUNCTION, parameterized
// scalars, collate), so the Spanner verdict is authoritative — there are no
// BigQuery-only type forms to triangulate. (RANGE / MAP / FUNCTION are parsed by
// the grammar and reported by the emulator as ACCEPT with a semantic "not
// supported", which is the grammar verdict we assert.)
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

// typeFixture pairs a bare type (fed to omni's parseType) with a full statement
// embedding that type in a position where the emulator yields a clean
// accept/reject verdict (NOT an "error"/Unimplemented one). wantParse is the
// expected omni outcome and MUST equal the emulator's grammar verdict; it is
// stated here so a divergence between the asserted truth and the live oracle is
// caught explicitly.
type typeFixture struct {
	typeText  string // bare type for parseType
	stmt      string // full statement for the oracle (one stmt, no trailing ';')
	wantParse bool   // expected: omni parses (and oracle accepts)
}

// typeFixtures are oracle-probed type forms. Each `stmt` was chosen to produce a
// deterministic emulator verdict: scalars/arrays/range/map/function via
// CAST(NULL AS …); structs via field extraction so the emulator does not hit
// its "struct value cannot be returned" Unimplemented path.
var typeFixtures = []typeFixture{
	// --- scalars (type_name path_expression) — all accept ---
	{"INT64", "SELECT CAST(NULL AS INT64)", true},
	{"STRING", "SELECT CAST(NULL AS STRING)", true},
	{"BYTES", "SELECT CAST(NULL AS BYTES)", true},
	{"BOOL", "SELECT CAST(NULL AS BOOL)", true},
	{"FLOAT64", "SELECT CAST(NULL AS FLOAT64)", true},
	{"FLOAT32", "SELECT CAST(NULL AS FLOAT32)", true},
	{"NUMERIC", "SELECT CAST(NULL AS NUMERIC)", true},
	{"BIGNUMERIC", "SELECT CAST(NULL AS BIGNUMERIC)", true},
	{"JSON", "SELECT CAST(NULL AS JSON)", true},
	{"DATE", "SELECT CAST(NULL AS DATE)", true},
	{"TIMESTAMP", "SELECT CAST(NULL AS TIMESTAMP)", true},
	{"GEOGRAPHY", "SELECT CAST(NULL AS GEOGRAPHY)", true},
	{"TOKENLIST", "SELECT CAST(NULL AS TOKENLIST)", true},
	{"INT", "SELECT CAST(NULL AS INT)", true},
	{"INTERVAL", "SELECT CAST(NULL AS INTERVAL)", true},
	{"foo.bar.Baz", "SELECT CAST(NULL AS foo.bar.Baz)", true},
	{"`my.proto.Type`", "SELECT CAST(NULL AS `my.proto.Type`)", true},
	{"NUMERIC.foo", "SELECT CAST(NULL AS NUMERIC.foo)", true},

	// --- parameterized scalars — all accept ---
	{"STRING(100)", "SELECT CAST(NULL AS STRING(100))", true},
	{"BYTES(256)", "SELECT CAST(NULL AS BYTES(256))", true},
	{"NUMERIC(10, 2)", "SELECT CAST(NULL AS NUMERIC(10, 2))", true},
	{"STRING(MAX)", "SELECT CAST(NULL AS STRING(MAX))", true},
	{"NUMERIC(10, 2, 3)", "SELECT CAST(NULL AS NUMERIC(10, 2, 3))", true},
	{"NUMERIC(true)", "SELECT CAST(NULL AS NUMERIC(true))", true},
	{"NUMERIC('foo')", "SELECT CAST(NULL AS NUMERIC('foo'))", true},

	// --- collate — accept ---
	{"INT64 COLLATE 'und:ci'", "SELECT CAST(NULL AS INT64 COLLATE 'und:ci')", true},
	{"DATE COLLATE 'und:ci'", "SELECT CAST(NULL AS DATE COLLATE 'und:ci')", true},

	// --- ARRAY — accept (incl. nested >> split) ---
	{"ARRAY<INT64>", "SELECT CAST(NULL AS ARRAY<INT64>)", true},
	{"ARRAY<ARRAY<INT64>>", "SELECT CAST(NULL AS ARRAY<ARRAY<INT64>>)", true},
	{"ARRAY<INT64 COLLATE 'und:ci'>", "SELECT 1 FROM UNNEST(CAST(NULL AS ARRAY<INT64 COLLATE 'und:ci'>))", true},

	// --- RANGE / MAP / FUNCTION — accept (grammar parses; semantic-unsupported) ---
	{"RANGE<DATE>", "SELECT CAST(NULL AS RANGE<DATE>)", true},
	{"RANGE<INT64>", "SELECT CAST(NULL AS RANGE<INT64>)", true},
	{"MAP<INT64, STRING>", "SELECT CAST(NULL AS MAP<INT64, STRING>)", true},
	{"FUNCTION<INT64 -> STRING>", "SELECT CAST(NULL AS FUNCTION<INT64 -> STRING>)", true},
	{"FUNCTION<(INT64, STRING) -> BOOL>", "SELECT CAST(NULL AS FUNCTION<(INT64, STRING) -> BOOL>)", true},
	{"FUNCTION<() -> INT64>", "SELECT CAST(NULL AS FUNCTION<() -> INT64>)", true},

	// --- STRUCT (field-extraction wrappers to dodge the struct-return limit) ---
	{"STRUCT<x INT64, y STRING>", "SELECT STRUCT<x INT64, y STRING>(1, 'a').x", true},
	{"STRUCT<INT64, STRING>", "SELECT STRUCT<INT64, STRING>(1, 'a').field1", true},
	{"STRUCT<a ARRAY<INT64>>", "SELECT STRUCT<a ARRAY<INT64>>([1,2]).a", true},
	{"STRUCT<ARRAY<INT64>>", "SELECT STRUCT<ARRAY<INT64>>([1,2]).field1", true},
	{"STRUCT<x ARRAY<STRUCT<y INT64>>>", "SELECT STRUCT<x ARRAY<STRUCT<y INT64>>>([STRUCT<y INT64>(1)]).x", true},
	{"STRUCT<x INT64 COLLATE 'und:ci'>", "SELECT STRUCT<x INT64 COLLATE 'und:ci'>('a').x", true},

	// --- negatives — all reject ---
	{"ARRAY<>", "SELECT CAST(NULL AS ARRAY<>)", false},
	{"ARRAY< >", "SELECT CAST(NULL AS ARRAY< >)", false},
	{"ARRAY<INT64, STRING>", "SELECT CAST(NULL AS ARRAY<INT64, STRING>)", false},
	{"ARRAY<INT64>>", "SELECT CAST(NULL AS ARRAY<INT64>>)", false},
	{"ARRAY", "SELECT CAST(NULL AS ARRAY)", false},
	{"ARRAY.foo", "SELECT CAST(NULL AS ARRAY.foo)", false},
	{"STRUCT", "SELECT CAST(NULL AS STRUCT)", false},
	{"STRUCT<x INT64,>", "SELECT STRUCT<x INT64,>(1).x", false},
	{"RANGE", "SELECT CAST(NULL AS RANGE)", false},
	{"RANGE<>", "SELECT CAST(NULL AS RANGE<>)", false},
	{"INTERVAL.foo", "SELECT CAST(NULL AS INTERVAL.foo)", false},
	{"INTERVAL<INT64>", "SELECT CAST(NULL AS INTERVAL<INT64>)", false},
	{"STRING()", "SELECT CAST(NULL AS STRING())", false},
	{"STRING(,)", "SELECT CAST(NULL AS STRING(,))", false},
	{"NUMERIC(-5)", "SELECT CAST(NULL AS NUMERIC(-5))", false},
	{"NUMERIC(10,)", "SELECT CAST(NULL AS NUMERIC(10,))", false},
	{"ARRAY<INT64>(vector_length=>128)", "SELECT CAST(NULL AS ARRAY<INT64>(vector_length=>128))", false},
	{"MAP<>", "SELECT CAST(NULL AS MAP<>)", false},
	{"FUNCTION<>", "SELECT CAST(NULL AS FUNCTION<>)", false},
	{"MAP<INT64>", "SELECT CAST(NULL AS MAP<INT64>)", false},
	{"FUNCTION<INT64>", "SELECT CAST(NULL AS FUNCTION<INT64>)", false},
	{"``", "SELECT CAST(NULL AS ``)", false},

	// --- additional nested forms (>> / >>> split with named fields) — accept ---
	{"STRUCT<STRUCT<x INT64>>", "SELECT STRUCT<STRUCT<x INT64>>(STRUCT<x INT64>(1)).field1", true},
	{"STRUCT<x ARRAY<ARRAY<INT64>>>", "SELECT STRUCT<x ARRAY<ARRAY<INT64>>>([[1]]).x", true},
	{"ARRAY<STRUCT<a INT64 COLLATE 'x'>>", "SELECT CAST(NULL AS ARRAY<STRUCT<a INT64 COLLATE 'x'>>)", true},
	{"ARRAY<NUMERIC(10, 2)>", "SELECT CAST(NULL AS ARRAY<NUMERIC(10, 2)>)", true},
	{"STRING(0)", "SELECT CAST(NULL AS STRING(0))", true},
	{"STRING(0x10)", "SELECT CAST(NULL AS STRING(0x10))", true},
	{"STRUCT<key INT64>", "SELECT STRUCT<key INT64>(1).key", true},
	{"STRUCT<row INT64>", "SELECT STRUCT<row INT64>(1).row", true},
}

func TestTypeDifferential(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set; skipping live differential")
	}
	h := newTypeHarness(t)
	defer h.close()

	for _, fx := range typeFixtures {
		fx := fx
		t.Run(fx.typeText, func(t *testing.T) {
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

			// Sanity: the asserted wantParse must match the live oracle. If this
			// trips, the fixture's hard-coded expectation drifted from reality.
			if oracleAccepts != fx.wantParse {
				t.Fatalf("fixture wantParse=%v but oracle says %q (%s) for %q",
					fx.wantParse, v.Verdict, v.Message, fx.stmt)
			}

			// 2. omni parseType verdict on the bare type.
			_, errs := ParseDataType(fx.typeText)
			omniAccepts := len(errs) == 0

			if omniAccepts != oracleAccepts {
				t.Errorf("DIVERGENCE on %q: omni accepts=%v, oracle accepts=%v (%s: %s)",
					fx.typeText, omniAccepts, oracleAccepts, v.Reason, v.Message)
			}
		})
	}
}

// --- harness plumbing (one persistent batch process, mirrors mssql's) ---

type harnessVerdict struct {
	Verdict string `json:"verdict"`
	Kind    string `json:"kind"`
	Reason  string `json:"reason"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type typeHarness struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	mu     sync.Mutex
}

func newTypeHarness(t *testing.T) *typeHarness {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	// googlesql/parser/datatypes_oracle_test.go → repo root is ../..
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	projDir := filepath.Join(repoRoot, "harness", "googlesql-spanner")
	if _, err := os.Stat(projDir); err != nil {
		t.Skipf("harness project not found at %s", projDir)
	}

	bin := filepath.Join(projDir, "googlesql-spanner")
	// Prefer a prebuilt binary; build it if missing.
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
	return &typeHarness{cmd: cmd, stdin: stdin, stdout: bufio.NewReader(stdoutPipe)}
}

func (h *typeHarness) verdict(t *testing.T, stmt string) harnessVerdict {
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
	var v harnessVerdict
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &v); err != nil {
		t.Fatalf("decoding harness verdict %q: %v", line, err)
	}
	return v
}

func (h *typeHarness) close() {
	if h.stdin != nil {
		_ = h.stdin.Close()
	}
	if h.cmd != nil && h.cmd.Process != nil {
		_ = h.cmd.Wait()
	}
}
