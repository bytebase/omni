//go:build googlesql_oracle

// Differential test for the parser-ddl node (core DDL: CREATE/ALTER/DROP TABLE/
// VIEW/INDEX/SCHEMA/DATABASE) against the live Cloud Spanner emulator oracle. Run
// with:
//
//	SPANNER_EMULATOR_HOST=localhost:9010 \
//	  go test -tags googlesql_oracle ./googlesql/parser/ -run TestDDLDifferential
//
// It is the PROVE gate for googlesql/parser-ddl (correctness-protocol.md): for
// every fixture it (1) feeds the full DDL statement to the emulator via the
// harness/googlesql-spanner CLI (which routes CREATE/ALTER/DROP to
// UpdateDatabaseDdl) and reads the accept/reject verdict, then (2) parses the
// same statement with omni's Parse, and asserts the two agree — BOTH polarities
// (the corpus has accept AND reject cases). A harness `error` verdict (oracle
// could not decide) fails the fixture loudly; it is never folded into a verdict.
//
// ⚠ The emulator speaks Spanner's DDL dialect, a SUBSET of the BigQuery+Spanner
// union the parser must accept (oracle.md). Every fixture below is one of:
//   - SHARED / Spanner-authoritative DDL (INTERLEAVE, generated AS (expr) STORED,
//     ROW DELETION POLICY, the trailing PRIMARY KEY element list, NULL_FILTERED /
//     STORING / interleaved indexes, SQL SECURITY INVOKER views, CREATE/DROP
//     SCHEMA, CREATE DATABASE, the Spanner ALTER actions) — the Spanner verdict is
//     authoritative.
//   - A negative (syntax-error) case the emulator rejects and omni must reject too.
//
// BigQuery-only DDL whose Spanner verdict is NON-authoritative is EXCLUDED from
// this differential and is covered by the hand-written unit tests
// (create_table_test.go / alter_table_test.go / ddl_objects_test.go) + the
// divergence ledger instead. Those forms (OR REPLACE, TEMP, PARTITION BY,
// CLUSTER BY, table-level OPTIONS, CTAS, inline PRIMARY KEY, NOT ENFORCED, LIKE,
// COPY, GENERATED ALWAYS AS, STORED VOLATILE, VIRTUAL, SQL SECURITY DEFINER,
// SET DATA TYPE, RENAME COLUMN, multi-action ALTER, ALTER COLUMN DROP NOT NULL,
// DEFAULT COLLATE, an unknown OPTIONS name) all SYNTAX- or feature-reject on the
// Spanner emulator while being valid GoogleSQL the union parser must accept; a
// Spanner verdict for them would be a false DIVERGENCE.
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

// ddlFixture is a full DDL statement (no trailing ';') fed to BOTH the oracle and
// omni's Parse. wantParse is the expected omni outcome and MUST equal the
// emulator's grammar verdict.
type ddlFixture struct {
	sql       string
	wantParse bool
}

var ddlFixtures = []ddlFixture{
	// ===================== ACCEPT — CREATE TABLE (Spanner) =====================
	{"CREATE TABLE T (x INT64) PRIMARY KEY (x)", true},
	{"CREATE TABLE T (x INT64, y STRING(MAX) NOT NULL) PRIMARY KEY (x)", true},
	{"CREATE TABLE IF NOT EXISTS T (x INT64) PRIMARY KEY (x)", true},
	{"CREATE TABLE s.T (x INT64) PRIMARY KEY (x)", true},
	{"CREATE TABLE T (x INT64) PRIMARY KEY (x DESC)", true},
	{"CREATE TABLE T (x INT64, y INT64) PRIMARY KEY (x, y)", true},
	{"CREATE TABLE T () PRIMARY KEY ()", true},
	{"CREATE TABLE T (x INT64,) PRIMARY KEY (x)", true}, // trailing comma
	// Column variety.
	{"CREATE TABLE T (x INT64, a BOOL, b FLOAT64, c NUMERIC, d DATE, e TIMESTAMP, f BYTES(MAX), g JSON) PRIMARY KEY (x)", true},
	{"CREATE TABLE T (x INT64, tags ARRAY<INT64>) PRIMARY KEY (x)", true}, // ARRAY<STRING> (no length) is BigQuery-only — see unit test
	{"CREATE TABLE T (x INT64, tags ARRAY<STRING(MAX)>) PRIMARY KEY (x)", true},
	// ARRAY vector-length parameter (Spanner column-schema; divergence #202).
	// The element-type FLOAT32/FLOAT64 restriction is semantic — ARRAY<INT64>(...)
	// PARSES (oracle accepts the grammar, then semantic-rejects), so wantParse=true.
	{"CREATE TABLE T (x INT64 NOT NULL, e ARRAY<FLOAT32>(vector_length=>128)) PRIMARY KEY (x)", true},
	{"CREATE TABLE T (x INT64 NOT NULL, e ARRAY<FLOAT64>(vector_length => 256)) PRIMARY KEY (x)", true},
	{"CREATE TABLE T (x INT64 NOT NULL, e ARRAY<FLOAT32>(vector_length=>128) NOT NULL) PRIMARY KEY (x)", true},
	{"CREATE TABLE T (x INT64 NOT NULL, e ARRAY<INT64>(vector_length=>4)) PRIMARY KEY (x)", true},
	// Generated / default / options / constraints (Spanner-authoritative).
	{"CREATE TABLE T (x INT64, f STRING(10) AS (CONCAT(x)) STORED) PRIMARY KEY (x)", true},
	{"CREATE TABLE T (x INT64, s STRING(20) DEFAULT ('a')) PRIMARY KEY (x)", true},
	{"CREATE TABLE T (x INT64, c TIMESTAMP OPTIONS (allow_commit_timestamp = true)) PRIMARY KEY (x)", true},
	{"CREATE TABLE T (x INT64, CONSTRAINT chk CHECK (x > 0)) PRIMARY KEY (x)", true},
	{"CREATE TABLE T (x INT64, FOREIGN KEY (x) REFERENCES P (id)) PRIMARY KEY (x)", true},
	{"CREATE TABLE T (x INT64, CONSTRAINT fk FOREIGN KEY (x) REFERENCES P (id) ON DELETE NO ACTION) PRIMARY KEY (x)", true},
	// Interleave + ROW DELETION POLICY (comma-led trailing clauses).
	{"CREATE TABLE Albums (s INT64, a INT64) PRIMARY KEY (s, a), INTERLEAVE IN PARENT Singers ON DELETE CASCADE", true},
	{"CREATE TABLE D (p INT64, d INT64) PRIMARY KEY (p, d), INTERLEAVE IN PARENT P ON DELETE NO ACTION", true},
	{"CREATE TABLE T (x INT64, c TIMESTAMP) PRIMARY KEY (x), ROW DELETION POLICY (OLDER_THAN(c, INTERVAL 30 DAY))", true},

	// ===================== ACCEPT — ALTER TABLE (Spanner) =====================
	{"ALTER TABLE T ADD COLUMN c STRING(10)", true},
	{"ALTER TABLE T ADD COLUMN IF NOT EXISTS c STRING(MAX)", true},
	{"ALTER TABLE T DROP COLUMN c", true},
	{"ALTER TABLE T ALTER COLUMN c STRING(20)", true},
	// Spanner ALTER COLUMN <name> <column_schema_inner> with the vector-length
	// parameter (oracle parses, semantic-rejects the change). SET DATA TYPE is a
	// BigQuery-only form excluded from this differential (covered by unit tests).
	{"ALTER TABLE T ALTER COLUMN e ARRAY<FLOAT32>(vector_length=>128)", true},
	{"ALTER TABLE T ADD COLUMN e ARRAY<FLOAT32>(vector_length=>64)", true},
	{"ALTER TABLE T ALTER COLUMN c SET OPTIONS (allow_commit_timestamp = true)", true},
	{"ALTER TABLE T ALTER COLUMN c SET DEFAULT (1)", true},
	{"ALTER TABLE T ALTER COLUMN c DROP DEFAULT", true},
	{"ALTER TABLE T ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES P (id)", true},
	{"ALTER TABLE T DROP CONSTRAINT fk", true},
	{"ALTER TABLE T SET ON DELETE CASCADE", true},
	{"ALTER TABLE T RENAME TO T2", true},
	{"ALTER TABLE T ADD ROW DELETION POLICY (OLDER_THAN(c, INTERVAL 1 DAY))", true},
	{"ALTER TABLE T DROP ROW DELETION POLICY", true},

	// ===================== ACCEPT — INDEX (Spanner) =====================
	{"CREATE INDEX idx ON T (a, b)", true},
	{"CREATE UNIQUE INDEX idx ON T (a)", true},
	{"CREATE NULL_FILTERED INDEX idx ON T (a)", true},
	{"CREATE INDEX idx ON T (a DESC)", true},
	{"CREATE INDEX idx ON T (a) STORING (b, c)", true},
	{"CREATE INDEX idx ON T (a) STORING (b), INTERLEAVE IN P", true},
	{"CREATE INDEX IF NOT EXISTS idx ON T (a)", true},
	{"CREATE INDEX idx ON T (a) WHERE a IS NOT NULL", true},
	{"CREATE INDEX idx ON s.T (a)", true},
	{"CREATE INDEX idx ON T ()", true}, // empty key list — emulator accepts
	{"ALTER INDEX idx ADD STORED COLUMN c", true},
	{"ALTER INDEX idx DROP STORED COLUMN c", true},
	{"DROP INDEX idx", true},
	{"DROP INDEX IF EXISTS idx", true},

	// ===================== ACCEPT — VIEW (Spanner) =====================
	{"CREATE VIEW v SQL SECURITY INVOKER AS SELECT 1", true},
	{"CREATE VIEW v AS SELECT 1", true},
	{"CREATE OR REPLACE VIEW v SQL SECURITY INVOKER AS SELECT a FROM T", true},
	{"DROP VIEW v", true},
	{"DROP VIEW IF EXISTS v", true},

	// ===================== ACCEPT — SCHEMA / DATABASE (Spanner) =====================
	{"CREATE SCHEMA sch", true},
	{"CREATE SCHEMA IF NOT EXISTS sch", true},
	{"DROP SCHEMA sch", true},
	{"DROP SCHEMA IF EXISTS sch", true},
	{"CREATE DATABASE db", true},
	{"ALTER DATABASE db SET OPTIONS (version_retention_period = '7d')", true}, // known option name

	// ===================== ACCEPT — DROP TABLE =====================
	{"DROP TABLE T", true},
	{"DROP TABLE IF EXISTS T", true},

	// Interleave THEN ROW DELETION POLICY — the standard (accepted) trailing order.
	{"CREATE TABLE C (p INT64, c INT64, t TIMESTAMP) PRIMARY KEY (p, c), INTERLEAVE IN PARENT P ON DELETE CASCADE, ROW DELETION POLICY (OLDER_THAN(t, INTERVAL 1 DAY))", true},

	// ============================ REJECT (syntax) ============================
	{"CREATE TABLE", false},                           // missing name
	{"CREATE TABLE T (x INT64) PRIMARY KEY x", false}, // missing parens on PK
	{"CREATE TABLE T (x)", false},                     // column missing type
	{"CREATE TABLE T (x INT64) PRIMARY KEY (x) ROW DELETION POLICY (OLDER_THAN(x, INTERVAL 1 DAY))", false}, // missing comma before ROW DELETION POLICY
	// ROW DELETION POLICY before INTERLEAVE — out of order; the emulator rejects the
	// reverse of the accepted order above (regression: parser-ddl PR #220 finding 3).
	{"CREATE TABLE C (p INT64, c INT64, t TIMESTAMP) PRIMARY KEY (p, c), ROW DELETION POLICY (OLDER_THAN(t, INTERVAL 1 DAY)), INTERLEAVE IN PARENT P ON DELETE CASCADE", false},
	{"ALTER TABLE T", false},             // missing action list
	{"ALTER TABLE T DROP COLUMN", false}, // missing column name
	{"CREATE INDEX idx", false},          // missing ON + keys
	{"CREATE INDEX idx ON T", false},     // missing key list
	{"CREATE VIEW v", false},             // missing AS query
	{"DROP TABLE", false},                // missing name
	{"DROP INDEX", false},                // missing name
	// DROP TABLE has no opt_drop_mode (table_or_table_function alt) — CASCADE/RESTRICT
	// reject (regression: parser-ddl PR #220 finding 1).
	{"DROP TABLE T CASCADE", false},
	{"DROP TABLE T RESTRICT", false},
	// Plain DROP INDEX (schema_object_kind) has no on_path_expression — a trailing
	// `ON <table>` rejects (regression: parser-ddl PR #220 finding 2).
	{"DROP INDEX idx ON T", false},
	// ARRAY vector-length parameter — the grammar is tight (divergence #202):
	// exactly `( vector_length => integer )`. The oracle rejects each of these.
	{"CREATE TABLE T (x INT64 NOT NULL, e ARRAY<FLOAT32>(vector_length=>8,)) PRIMARY KEY (x)", false},        // trailing comma
	{"CREATE TABLE T (x INT64 NOT NULL, e ARRAY<FLOAT32>(vector_length=>8, foo=>9)) PRIMARY KEY (x)", false}, // second parameter
	{"CREATE TABLE T (x INT64 NOT NULL, e ARRAY<FLOAT32>(length=>8)) PRIMARY KEY (x)", false},                // wrong name
	{"CREATE TABLE T (x INT64 NOT NULL, e ARRAY<FLOAT32>(vector_length=>MAX)) PRIMARY KEY (x)", false},       // non-integer value
	{"CREATE TABLE T (x INT64 NOT NULL, e ARRAY<FLOAT32>(`vector_length`=>8)) PRIMARY KEY (x)", false},       // backtick-quoted name
	{"CREATE TABLE T (x INT64 NOT NULL, e ARRAY<FLOAT32>(128)) PRIMARY KEY (x)", false},                      // positional param (only vector_length allowed on an array column)
	{"CREATE TABLE T (x INT64 NOT NULL, e ARRAY<FLOAT32>(vector_length)) PRIMARY KEY (x)", false},            // missing => value
	{"CREATE TABLE T (x INT64 NOT NULL, e ARRAY<FLOAT32>(vector_length=>-5)) PRIMARY KEY (x)", false},        // negative value
	// vector_length is unsupported inside a STRUCT-of-ARRAY (oracle: a true syntax
	// reject "'vector_length' is not supported in STRUCT of ARRAY").
	{"CREATE TABLE T (x INT64 NOT NULL, e STRUCT<v ARRAY<FLOAT32>(vector_length=>8)>) PRIMARY KEY (x)", false},
}

func TestDDLDifferential(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set; skipping live differential")
	}
	h := newDDLHarness(t)
	defer h.close()

	for _, fx := range ddlFixtures {
		fx := fx
		t.Run(fx.sql, func(t *testing.T) {
			// 1. Oracle verdict for the statement.
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
			// trips, the fixture's hard-coded expectation drifted from reality
			// (e.g. an emulator upgrade) — or the fixture is a BigQuery-only form
			// that should not be in this Spanner differential.
			if oracleAccepts != fx.wantParse {
				t.Fatalf("fixture wantParse=%v but oracle says %q (%s: %s) for %q",
					fx.wantParse, v.Verdict, v.Reason, v.Message, fx.sql)
			}

			// 2. omni Parse verdict.
			_, errs := Parse(fx.sql)
			omniAccepts := len(errs) == 0

			if omniAccepts != oracleAccepts {
				t.Errorf("DIVERGENCE on %q: omni accepts=%v, oracle accepts=%v (%s: %s); omni errs=%v",
					fx.sql, omniAccepts, oracleAccepts, v.Reason, v.Message, errs)
			}
		})
	}
}

// --- harness plumbing (one persistent batch process; mirrors the type/select/
// expr oracle harnesses in this package — each names its own type to avoid a
// duplicate-symbol collision under the shared googlesql_oracle build tag) ---

type ddlHarnessVerdict struct {
	Verdict string `json:"verdict"`
	Kind    string `json:"kind"`
	Reason  string `json:"reason"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ddlHarness struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	mu     sync.Mutex
}

func newDDLHarness(t *testing.T) *ddlHarness {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	// googlesql/parser/ddl_oracle_test.go → repo root is ../..
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
	return &ddlHarness{cmd: cmd, stdin: stdin, stdout: bufio.NewReader(stdoutPipe)}
}

func (h *ddlHarness) verdict(t *testing.T, stmt string) ddlHarnessVerdict {
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
	var v ddlHarnessVerdict
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &v); err != nil {
		t.Fatalf("decoding harness verdict %q: %v", line, err)
	}
	return v
}

func (h *ddlHarness) close() {
	if h.stdin != nil {
		_ = h.stdin.Close()
	}
	if h.cmd != nil && h.cmd.Process != nil {
		_ = h.cmd.Wait()
	}
}
