package parser

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/hive/internal/testcontainer"
)

// TestMain removes the pinned Hive container if a recording or verification
// started one.
func TestMain(m *testing.M) {
	code := m.Run()
	testcontainer.Terminate()
	os.Exit(code)
}

// TestCompatibility holds omni's hive/parser to testdata/compat.tsv: a
// statement for each alternative of ../testdata/Hive.g4 and its near-misses,
// with omni's verdict and Apache Hive 4.0.1's. Hive.g4 is Apache Hive's own
// ANTLR 3 grammar (../testdata/vendor/*.g at rel/release-4.0.1), converted to
// ANTLR 4 by cmd/antlr3to4 (TestHiveG4IsTheConversion holds it to that).
//
// In normal runs, this test reads testdata/compat.tsv offline and starts no
// Docker container. With OMNI_RECORD=1 it asks the pinned apache/hive:4.0.1
// server via EXPLAIN AST CREATE SCHEDULED QUERY ... AS <stmt> and writes
// verdicts; with OMNI_RECORD=verify it re-checks every committed verdict
// read-only against the pinned server.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join(td, "pins.tsv"),
		Grammar: filepath.Join(td, "Hive.g4"),
		Tokens:  filepath.Join(td, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

// TestGrammarVendorDigest holds the vendored ANTLR 3 grammar files in
// ../testdata/vendor/*.g to the grammar row of ../testdata/pins.tsv, so none
// is edited without updating its pin.
func TestGrammarVendorDigest(t *testing.T) {
	td := filepath.Join("..", "testdata")
	pin, err := enginetest.LoadPin(filepath.Join(td, "pins.tsv"), "grammar")
	if err != nil {
		t.Fatal(err)
	}
	vendor := []string{
		filepath.Join(td, "vendor", "HiveParser.g"),
		filepath.Join(td, "vendor", "AlterClauseParser.g"),
		filepath.Join(td, "vendor", "SelectClauseParser.g"),
		filepath.Join(td, "vendor", "FromClauseParser.g"),
		filepath.Join(td, "vendor", "IdentifiersParser.g"),
		filepath.Join(td, "vendor", "ResourcePlanParser.g"),
		filepath.Join(td, "vendor", "CreateDDLParser.g"),
		filepath.Join(td, "vendor", "PrepareStatementParser.g"),
		filepath.Join(td, "vendor", "HiveLexer.g"),
		filepath.Join(td, "vendor", "HiveLexerParent.g"),
	}
	got, err := enginetest.FilesDigest(vendor...)
	if err != nil {
		t.Fatal(err)
	}
	if got != pin.Digest {
		t.Fatalf("pins.tsv grammar row pins %s, and %s hash to %s", pin.Digest, strings.Join(vendor, ", "), got)
	}
}

func omniVerdict(sql string) string {
	res, err := Parse(sql)
	if err != nil || res == nil || res.Statement == nil {
		return compat.Reject
	}
	return compat.Accept
}

const (
	// hiveJudgePrefix wraps every probed statement inside Hive's
	// `createScheduledQueryStatement -> definedAsSpec -> statement` grammar
	// path under `EXPLAIN AST`.
	//
	// Because `definedAsSpec` directly invokes HiveParser.g's start rule
	// `statement : explainStatement EOF | execStatement EOF`, this probes the
	// exact top-level HiveQL grammar without executing any statement (not even
	// DDL, DML, KILL QUERY, USE, SET ROLE, or EXPLAIN ANALYZE) and without
	// colliding with EXPLAIN's own `explainOption*` prefix (e.g., ANALYZE TABLE
	// vs EXTENDED COMMIT vs nested EXPLAIN).
	hiveJudgePrefix  = "EXPLAIN AST CREATE SCHEDULED QUERY omni_probe EVERY 1 DAY AS "
	hiveProbeTail    = "\n)"
	hiveQueryTimeout = 30 * time.Second
)

func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	db := testcontainer.Judge(t, pin)

	run := func(query string) error {
		ctx, cancel := context.WithTimeout(context.Background(), hiveQueryTimeout)
		defer cancel()
		rows, err := db.QueryContext(ctx, query)
		if err == nil {
			for rows.Next() {
			}
			err = rows.Err()
			_ = rows.Close()
		}
		return err
	}

	ask := func(stmt string) (string, error) {
		if strings.ContainsRune(stmt, 0) {
			return "", fmt.Errorf("statement contains NUL byte")
		}
		err := run(hiveJudgePrefix + stmt)
		return hiveVerdict(err, func() error {
			return run(hiveJudgePrefix + stmt + hiveProbeTail)
		})
	}

	// Sanity-check the judge contract before judging any rows so a server
	// upgrade that changes ParseException/SemanticException reporting fails
	// immediately instead of misclassifying rows.
	for _, tc := range []struct {
		stmt string
		want string
	}{
		{"SELECT 1", compat.Accept},
		{"SELECT a FROM nonexistent_tbl", compat.Accept},
		{"SELECT FROM", compat.Reject},
		{"CREATE TABLE omni_probe_not_created (a INT)", compat.Accept},
		{"ANALYZE TABLE a COMPUTE STATISTICS", compat.Accept},
		{"EXTENDED COMMIT", compat.Reject},
		{"EXPLAIN COMMIT", compat.Accept},
		{"EXPLAIN ANALYZE SELECT 1", compat.Accept},
		{"EXPLAIN zzz COMMIT", compat.Reject},
		{"CREATE TABLE a (`a.b` INT)", compat.Reject},
	} {
		got, err := ask(tc.stmt)
		if err != nil || got != tc.want {
			t.Fatalf("Hive sanity check on %q = (%q, %v), want (%q, nil)", tc.stmt, got, err, tc.want)
		}
	}

	return ask
}

// hiveVerdict classifies HiveServer2's answer to
// `EXPLAIN AST CREATE SCHEDULED QUERY omni_probe EVERY 1 DAY AS <stmt>`.
//
// It is fail-closed:
//   - nil -> accept (compiled and explained cleanly).
//   - non-*testcontainer.Error (transport, timeout, closed connection) -> error.
//   - FAILED: ParseException / RewriteEmptyStreamException /
//     RewriteCardinalityException -> reject (lexer/parser/ANTLR rewrite error).
//   - Known post-parse compiler/execution exceptions (SemanticException and
//     subclasses, IllegalStateException, Execution Error, etc.) -> accept
//     (ParseDriver.parse succeeded).
//   - Code 0 "Error running query" (unhandled post-parse Throwable in
//     HiveServer2, e.g. CreateTableLikeAnalyzer on `CREATE TABLE a LIKE FILE`):
//     verified with tail probe (`stmt + "\n)"`). Only a probe ParseException
//     proves `stmt` parsed cleanly to EOF -> accept; any other probe outcome
//     fails the recording.
//   - Anything else -> error (never guess).
func hiveVerdict(err error, probe func() error) (string, error) {
	if err == nil {
		return compat.Accept, nil
	}
	var he *testcontainer.Error
	if !errors.As(err, &he) {
		return "", err
	}
	msg := he.Message
	if isHiveParseError(msg) {
		return compat.Reject, nil
	}
	for _, marker := range postParseHiveExceptions {
		if strings.Contains(msg, marker) {
			return compat.Accept, nil
		}
	}
	if he.Code == 0 && msg == "Error running query" && probe != nil {
		perr := probe()
		var phe *testcontainer.Error
		if errors.As(perr, &phe) && strings.Contains(phe.Message, "FAILED: ParseException") {
			return compat.Accept, nil
		}
		return "", fmt.Errorf("tail probe did not confirm parse after code 0 Error running query: %v", perr)
	}
	return "", fmt.Errorf("unexpected Hive error (code %d): %s", he.Code, msg)
}

func isHiveParseError(msg string) bool {
	return strings.Contains(msg, "FAILED: ParseException") ||
		strings.Contains(msg, "FAILED: RewriteEmptyStreamException") ||
		strings.Contains(msg, "FAILED: RewriteCardinalityException")
}

var postParseHiveExceptions = []string{
	"FAILED: SemanticException",
	"FAILED: CalciteSemanticException",
	"FAILED: CalciteSubquerySemanticException",
	"FAILED: CalciteViewSemanticException",
	"FAILED: UDFArgumentException",
	"FAILED: UDFArgumentTypeException",
	"FAILED: UDFArgumentLengthException",
	"FAILED: MetaException",
	"FAILED: IllegalStateException",
	"FAILED: Execution Error",
	"FAILED: HiveAccessControlException",
	"FAILED: InvalidOperationException",
	"FAILED: UnsupportedOperationException",
	"FAILED: IllegalArgumentException",
	"FAILED: NumberFormatException",
	"FAILED: NullPointerException",
	"FAILED: RuntimeException",
	"FAILED: HiveException",
	"FAILED: ClassCastException",
	"FAILED: IndexOutOfBoundsException",
	"FAILED: ArrayIndexOutOfBoundsException",
}

func TestHiveVerdict(t *testing.T) {
	parseErr := &testcontainer.Error{Code: 40000, Message: "Error while compiling statement: FAILED: ParseException line 2:0 extraneous input ')' expecting EOF"}
	cases := []struct {
		name           string
		err            error
		probeErr       error
		want           string
		wantErr        bool
		wantProbeCalls int
	}{
		{
			name: "nil succeeds -> accept",
			err:  nil,
			want: compat.Accept,
		},
		{
			name: "ParseException -> reject",
			err:  &testcontainer.Error{Code: 40000, Message: "Error while compiling statement: FAILED: ParseException line 1:68 cannot recognize input near 'FROM' '<EOF>' '<EOF>' in select clause"},
			want: compat.Reject,
		},
		{
			name: "ParseException failed predicate -> reject",
			err:  &testcontainer.Error{Code: 40000, Message: "Error while compiling statement: FAILED: ParseException line 1:86 Failed to recognize predicate ')'. Failed rule: '[., :] can not be used in column name in create table statement.'"},
			want: compat.Reject,
		},
		{
			name: "ANTLR RewriteEmptyStreamException -> reject",
			err:  &testcontainer.Error{Code: 40000, Message: "Error while compiling statement: FAILED: RewriteEmptyStreamException rule whenClauses"},
			want: compat.Reject,
		},
		{
			name: "SemanticException table not found -> accept",
			err:  &testcontainer.Error{Code: 10001, Message: "Error while compiling statement: FAILED: SemanticException [Error 10001]: Line 1:75 Table not found 'a'"},
			want: compat.Accept,
		},
		{
			name: "CalciteSubquerySemanticException -> accept",
			err:  &testcontainer.Error{Code: 10249, Message: "Error while compiling statement: FAILED: CalciteSubquerySemanticException [Error 10249]: Unsupported SubQuery Expression"},
			want: compat.Accept,
		},
		{
			name: "IllegalStateException COMMIT unsupported -> accept",
			err:  &testcontainer.Error{Code: 40000, Message: "Error while compiling statement: FAILED: IllegalStateException COMMIT is not supported yet."},
			want: compat.Accept,
		},
		{
			name: "Execution Error -> accept",
			err:  &testcontainer.Error{Code: 1, Message: "Error while compiling statement: FAILED: Execution Error, return code 1 from org.apache.hadoop.hive.ql.exec.ExplainTask"},
			want: compat.Accept,
		},
		{
			name:           "Error running query confirmed post-parse by tail probe -> accept",
			err:            &testcontainer.Error{Code: 0, Message: "Error running query"},
			probeErr:       parseErr,
			want:           compat.Accept,
			wantProbeCalls: 1,
		},
		{
			name:           "Error running query with non-parse tail probe -> error",
			err:            &testcontainer.Error{Code: 0, Message: "Error running query"},
			probeErr:       &testcontainer.Error{Code: 0, Message: "Error running query"},
			wantErr:        true,
			wantProbeCalls: 1,
		},
		{
			name:    "transport timeout -> error",
			err:     context.DeadlineExceeded,
			wantErr: true,
		},
		{
			name:    "unknown Hive error message -> error",
			err:     &testcontainer.Error{Code: 40000, Message: "Error while compiling statement: something completely unexpected"},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			got, err := hiveVerdict(tc.err, func() error {
				calls++
				return tc.probeErr
			})
			if calls != tc.wantProbeCalls {
				t.Errorf("probe calls = %d, want %d", calls, tc.wantProbeCalls)
			}
			if tc.wantErr {
				if err == nil {
					t.Fatalf("hiveVerdict() = (%q, nil), want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("hiveVerdict() unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("hiveVerdict() = %q, want %q", got, tc.want)
			}
		})
	}
}
