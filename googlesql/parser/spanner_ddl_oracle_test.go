//go:build googlesql_oracle

// Differential test for the parser-ddl-spanner node (Spanner DDL: CHANGE STREAM,
// SEQUENCE, ROLE, LOCALITY GROUP, PROTO BUNDLE, and role-based GRANT/REVOKE)
// against the live Cloud Spanner emulator oracle. Run with:
//
//	SPANNER_EMULATOR_HOST=localhost:9010 \
//	  go test -tags googlesql_oracle ./googlesql/parser/ -run TestSpannerDDLDifferential
//
// It is the PROVE gate for googlesql/parser-ddl-spanner (correctness-protocol.md):
// for every fixture it feeds the full DDL statement to the emulator via the
// harness/googlesql-spanner CLI (UpdateDatabaseDdl) and reads the accept/reject
// verdict, then parses the same statement with omni's Parse and asserts the two
// agree — BOTH polarities (the corpus has accept AND reject cases). A harness
// `error` verdict (oracle could not decide) fails the fixture loudly.
//
// AUTHORITATIVE — every fixture below is a SPANNER-dialect form (oracle.md routes
// "Spanner-only (… CHANGE STREAM, sequences, … role DDL)" to the emulator as
// authoritative). EXCLUDED, because the emulator is NON-authoritative for them:
//   - the legacy ZetaSQL string-grantee GRANT (`… TO 'user'`) — Spanner REJECTS
//     it ("Expecting 'ROLE'") but the BigQuery+Spanner UNION parser must accept
//     it; it is covered by grant_revoke_test.go.
//   - empty `OPTIONS ()` — Spanner REJECTS it in every context (incl. CREATE
//     TABLE) but the union/BigQuery grammar accepts it; the parser-ddl node
//     already accepts empty OPTIONS, so a Spanner verdict there is a false
//     divergence. Covered by the unit tests.
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

// spannerDDLFixture is a full DDL statement (no trailing ';') fed to BOTH the
// oracle and omni's Parse. wantParse is the expected omni outcome and MUST equal
// the emulator's grammar verdict.
type spannerDDLFixture struct {
	sql       string
	wantParse bool
}

var spannerDDLFixtures = []spannerDDLFixture{
	// ===================== CHANGE STREAM =====================
	{"CREATE CHANGE STREAM MyStream FOR ALL", true},
	{"CREATE CHANGE STREAM cs2 FOR ALL OPTIONS (retention_period = '7d')", true},
	{"CREATE CHANGE STREAM cs3 OPTIONS (retention_period = '7d')", true}, // no FOR clause
	{"CREATE CHANGE STREAM cs4 FOR ALL OPTIONS (retention_period = '2d', value_capture_type = 'NEW_ROW')", true},
	{"ALTER CHANGE STREAM MyStream SET FOR ALL", true},
	{"ALTER CHANGE STREAM MyStream DROP FOR ALL", true},
	{"ALTER CHANGE STREAM MyStream SET OPTIONS (retention_period = '2d')", true},
	{"DROP CHANGE STREAM MyStream", true},
	{"CREATE CHANGE STREAM cs_sch FOR myschema.Singers", true},          // schema-qualified table
	{"CREATE CHANGE STREAM cs_cols FOR Singers(`select`, key), Albums()", true}, // keyword/quoted cols, empty parens
	// reject — ALL is not a table name; FOR ALL cannot be in a list; no trailing comma.
	{"CREATE CHANGE STREAM cs FOR Singers, ALL", false},
	{"CREATE CHANGE STREAM cs FOR ALL, Singers", false},
	{"CREATE CHANGE STREAM cs FOR Singers(a,)", false}, // trailing comma in column list NOT allowed
	{"CREATE CHANGE STREAM cs FOR Singers,", false},    // trailing comma in table list NOT allowed
	// reject
	{"CREATE CHANGE STREAM cs FOR", false},        // dangling FOR
	{"ALTER CHANGE STREAM cs DROP FOR", false},    // DROP FOR without ALL
	{"ALTER CHANGE STREAM cs SET", false},         // SET with nothing
	{"ALTER CHANGE STREAM cs", false},             // no action

	// ===================== SEQUENCE =====================
	{"CREATE SEQUENCE Seq1 OPTIONS (sequence_kind = 'bit_reversed_positive')", true},
	{"CREATE SEQUENCE IF NOT EXISTS Seq2 OPTIONS (sequence_kind = 'bit_reversed_positive')", true},
	{"CREATE SEQUENCE BareSeq3", true}, // grammar-accept (semantic: missing kind)
	{"ALTER SEQUENCE Seq1 SET OPTIONS (start_with_counter = 9000)", true},
	{"ALTER SEQUENCE IF EXISTS SeqMissing SET OPTIONS (start_with_counter = 9000)", true},
	{"DROP SEQUENCE Seq1", true},
	{"DROP SEQUENCE IF EXISTS SeqX", true},
	// reject
	{"ALTER SEQUENCE s", false},     // missing SET OPTIONS (required for SEQUENCE)
	{"ALTER SEQUENCE s SET", false}, // SET with no OPTIONS

	// ===================== ROLE =====================
	{"CREATE ROLE analyst_role", true},
	{"DROP ROLE analyst_role", true},
	{"CREATE ROLE `dotted.name`", true}, // backtick name is a single identifier
	{"DROP ROLE a.b.c", true},           // DROP ROLE accepts a dotted path (CREATE does not)
	// reject — a dotted CREATE ROLE name (role names are single identifiers).
	{"CREATE ROLE a.b", false},

	// ===================== LOCALITY GROUP =====================
	{"CREATE LOCALITY GROUP lg_hot OPTIONS (storage = 'ssd')", true},
	{"CREATE LOCALITY GROUP lg_cold OPTIONS (storage = 'hdd', ssd_to_hdd_spill_timespan = '30d')", true},
	{"ALTER LOCALITY GROUP lg_cold SET OPTIONS (ssd_to_hdd_spill_timespan = '14d')", true},
	{"ALTER LOCALITY GROUP lg_bare", true}, // SET OPTIONS is OPTIONAL for LOCALITY GROUP (unlike SEQUENCE)
	{"DROP LOCALITY GROUP lg_cold", true},
	// reject
	{"ALTER LOCALITY GROUP g SET", false},          // SET with no OPTIONS
	{"ALTER LOCALITY GROUP g RENAME TO h", false},  // not a valid action

	// ===================== PROTO BUNDLE =====================
	// CREATE PROTO BUNDLE with an empty descriptor file is a SEMANTIC failure
	// (verdict accept — the grammar parsed). ALTER on a missing bundle likewise.
	{"CREATE PROTO BUNDLE (`a.b.C`)", true},
	{"CREATE PROTO BUNDLE (`a.b.C`, `a.b.D`)", true},
	{"ALTER PROTO BUNDLE INSERT (`a.b.C`)", true},
	{"ALTER PROTO BUNDLE INSERT (`a.b.C`) UPDATE (`a.b.D`) DELETE (`a.b.E`)", true},
	{"ALTER PROTO BUNDLE DELETE (`a.b.C`, `a.b.D`)", true},
	{"CREATE PROTO BUNDLE (`a.b.C`, `a.b.D`,)", true}, // trailing comma allowed (unlike change-stream cols)
	{"ALTER PROTO BUNDLE INSERT (`a.b.C`,)", true},    // trailing comma allowed
	// reject — empty bundle, comma between groups, repeat, out of order.
	{"CREATE PROTO BUNDLE ()", false},
	{"ALTER PROTO BUNDLE INSERT (`a.b.C`), UPDATE (`a.b.D`)", false}, // comma between groups
	{"ALTER PROTO BUNDLE INSERT (`a.b.C`) INSERT (`a.b.D`)", false},  // repeated group
	{"ALTER PROTO BUNDLE UPDATE (`a.b.C`) INSERT (`a.b.D`)", false},  // out of order
	{"ALTER PROTO BUNDLE INSERT", false},                            // INSERT with no parens

	// ===================== role-based GRANT / REVOKE =====================
	{"GRANT SELECT ON TABLE T_grant TO ROLE r_g", true},
	{"GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE T_grant TO ROLE r_g2", true},
	{"GRANT ROLE analyst_g TO ROLE senior_g", true},
	{"GRANT ROLE a1, b1 TO ROLE c1, d1", true},
	{"REVOKE SELECT ON TABLE T_grant FROM ROLE r_g", true},
	{"REVOKE ROLE analyst_g FROM ROLE senior_g", true},
	// Spanner comma-separated object list in a role grant/revoke.
	{"GRANT SELECT ON TABLE t_cs1, t_cs2 TO ROLE r_g", true},
	{"GRANT SELECT ON TABLE t_cs1, t_cs2 TO ROLE r_g1, r_g2", true},
	{"REVOKE SELECT ON TABLE t_cs1, t_cs2 FROM ROLE r_g", true},
	// reject — role-to-role target must be ROLE (not a string/parameter), and a
	// multi-object grant cannot use a string grantee.
	{"GRANT ROLE analyst_g TO 'user'", false},
	{"GRANT SELECT ON TABLE t_cs1, t_cs2 TO 'user'", false},
	// reject
	{"GRANT SELECT ON TABLE t TO ROLE", false}, // TO ROLE with no role name
	{"GRANT ROLE a TO ROLE", false},            // role-grant TO ROLE with no role
	{"GRANT ROLE TO ROLE r", false},            // GRANT ROLE with no subject role
	{"REVOKE ROLE a FROM ROLE", false},         // role-revoke FROM ROLE with no role
	// reject — role names are single identifiers, not dotted paths.
	{"GRANT SELECT ON TABLE t TO ROLE a.b", false},
	{"GRANT ROLE a.b TO ROLE c", false},
	{"GRANT SELECT ON TABLE t TO ROLE r1, r2.x", false},
	{"REVOKE ROLE r1, r2.x FROM ROLE s", false},

	// ===================== owned divergence #112 (INTERLEAVE) =====================
	// INTERLEAVE IN PARENT with NO `ON DELETE` action now ACCEPTS (defaults to NO
	// ACTION), matching the emulator; a dangling `ON DELETE` still REJECTS. This
	// statement is parsed by create_table.go; the fix is part of this node.
	{"CREATE TABLE cs_child (id INT64) PRIMARY KEY (id), INTERLEAVE IN PARENT cs_parent", true},
	{"CREATE TABLE cs_child2 (id INT64) PRIMARY KEY (id), INTERLEAVE IN PARENT cs_parent ON DELETE CASCADE", true},
	{"CREATE TABLE cs_child3 (id INT64) PRIMARY KEY (id), INTERLEAVE IN PARENT cs_parent ON DELETE", false}, // dangling ON DELETE
	// owned divergence #8 (inline column PRIMARY KEY) — emulator accepts; omni does too.
	{"CREATE TABLE cs_inline (id INT64 PRIMARY KEY, name STRING(MAX))", true},

	// ===================== OR REPLACE / scope / drop-mode rejects =====================
	// None of these Spanner objects takes OR REPLACE or a create-scope; DROP takes
	// no RESTRICT/CASCADE. All REJECT on the emulator.
	{"CREATE OR REPLACE SEQUENCE s_or", false},
	{"CREATE TEMP SEQUENCE s_temp", false},
	{"CREATE OR REPLACE CHANGE STREAM s_or FOR ALL", false},
	{"CREATE OR REPLACE LOCALITY GROUP g_or", false},
	{"CREATE OR REPLACE PROTO BUNDLE (`a.b.C`)", false},
	{"CREATE OR REPLACE ROLE r_or", false},
	{"DROP SEQUENCE s_dm RESTRICT", false},
}

func TestSpannerDDLDifferential(t *testing.T) {
	h := newSpannerDDLHarness(t)
	defer h.close()

	for _, fx := range spannerDDLFixtures {
		fx := fx
		t.Run(fx.sql, func(t *testing.T) {
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

			// Sanity: the asserted wantParse must match the live oracle. If this
			// trips, the fixture's hard-coded expectation drifted from reality (an
			// emulator upgrade) — or it should not be in this Spanner differential.
			if oracleAccepts != fx.wantParse {
				t.Fatalf("fixture wantParse=%v but oracle says %q (%s: %s) for %q",
					fx.wantParse, v.Verdict, v.Reason, v.Message, fx.sql)
			}

			_, errs := Parse(fx.sql)
			omniAccepts := len(errs) == 0

			if omniAccepts != oracleAccepts {
				t.Errorf("DIVERGENCE on %q: omni accepts=%v, oracle accepts=%v (%s: %s); omni errs=%v",
					fx.sql, omniAccepts, oracleAccepts, v.Reason, v.Message, errs)
			}
		})
	}
}

// --- harness plumbing (one persistent batch process; names its own type to avoid
// a duplicate-symbol collision under the shared googlesql_oracle build tag) ---

type spannerDDLHarnessVerdict struct {
	Verdict string `json:"verdict"`
	Kind    string `json:"kind"`
	Reason  string `json:"reason"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type spannerDDLHarness struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	mu     sync.Mutex
}

func newSpannerDDLHarness(t *testing.T) *spannerDDLHarness {
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
	return &spannerDDLHarness{cmd: cmd, stdin: stdin, stdout: bufio.NewReader(stdoutPipe)}
}

func (h *spannerDDLHarness) verdict(t *testing.T, stmt string) spannerDDLHarnessVerdict {
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
	var v spannerDDLHarnessVerdict
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &v); err != nil {
		t.Fatalf("decoding harness verdict %q: %v", line, err)
	}
	return v
}

func (h *spannerDDLHarness) close() {
	if h.stdin != nil {
		_ = h.stdin.Close()
	}
	if h.cmd != nil && h.cmd.Process != nil {
		_ = h.cmd.Wait()
	}
}
