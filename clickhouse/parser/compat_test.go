package parser

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"

	"github.com/bytebase/omni/clickhouse/internal/testcontainer"
	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of ClickHouse's ClickHouseParser.g4 and its near-misses, with
// omni's verdict and the pinned ClickHouse server's. The test reads the file
// and starts no server; OMNI_RECORD=1 judges unjudged rows against the pinned
// ClickHouse container via EXPLAIN AST, and OMNI_RECORD=verify re-judges every
// committed row without writing.
//
// testdata/README.md documents the judge contract and remaining gaps.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join(td, "pins.tsv"),
		Grammar: filepath.Join(td, "ClickHouseParser.g4"),
		Lexer:   filepath.Join(td, "ClickHouseLexer.g4"),
		Tokens:  filepath.Join(td, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

// TestGrammarVendorDigest holds the vendored ANTLR 4 grammar files
// (ClickHouseParser.g4, ClickHouseLexer.g4) to the grammar row of
// ../testdata/pins.tsv, so neither is edited without updating its pin.
func TestGrammarVendorDigest(t *testing.T) {
	td := filepath.Join("..", "testdata")
	pin, err := enginetest.LoadPin(filepath.Join(td, "pins.tsv"), "grammar")
	if err != nil {
		t.Fatal(err)
	}
	vendor := []string{
		filepath.Join(td, "ClickHouseParser.g4"),
		filepath.Join(td, "ClickHouseLexer.g4"),
	}
	got, err := enginetest.FilesDigest(vendor...)
	if err != nil {
		t.Fatal(err)
	}
	if got != pin.Digest {
		t.Fatalf("pins.tsv grammar row pins %s, and %s hash to %s", pin.Digest, strings.Join(vendor, ", "), got)
	}
}

// omniVerdict returns Parse's verdict as a compat.tsv cell: accept or reject.
func omniVerdict(sql string) string {
	stmts, err := Parse(sql)
	if err == nil && len(stmts) == 1 {
		return compat.Accept
	}
	if err != nil {
		var pe *ParseError
		if !errors.As(err, &pe) {
			panic(fmt.Sprintf("clickhouse/parser.Parse(%q) returned %T, want *ParseError", sql, err))
		}
	}
	return compat.Reject
}

const (
	clickhouseQueryTimeout = 15 * time.Second

	// ClickHouse server exception codes classified by clickhouseVerdict.
	chErrSyntaxError           int32 = 62  // SYNTAX_ERROR
	chErrBadArguments          int32 = 36  // BAD_ARGUMENTS (e.g. window frame start UNBOUNDED FOLLOWING)
	chErrUnknownStorage        int32 = 56  // UNKNOWN_STORAGE (post-parse engine lookup for ENGINE ... SETTINGS)
	chErrUnknownSetting        int32 = 115 // UNKNOWN_SETTING (post-parse SETTINGS validation)
	chErrCannotParseUUID       int32 = 376 // CANNOT_PARSE_UUID (parse-time UUID validation in ParserCreateQuery)
	chErrNotImplemented        int32 = 476 // NOT_IMPLEMENTED (parse-time WITH TIES without ORDER BY check)
	chErrExpectedAllOrDistinct int32 = 558 // EXPECTED_ALL_OR_DISTINCT (post-parse bare UNION when union_default_mode='')

	// clickhouseProbeTail is appended to a statement when EXPLAIN AST returns a
	// non-syntax ClickHouse error to distinguish parse-time rejections (thrown
	// inside the parser before EOF) from post-parse checks (thrown after the
	// statement parsed cleanly to EOF).
	//
	// The leading newline closes any trailing single-line comment (-- or #),
	// and "; SELECT 1" triggers ClickHouse's multi-statement SYNTAX_ERROR (62)
	// if and only if the statement itself parsed cleanly to EOF.
	clickhouseProbeTail = "\n; SELECT 1"
)

// isKnownSemanticCode reports whether code is a ClickHouse exception code that
// may arise either from a parse-time validation check or from post-parse
// analysis of an EXPLAIN AST query, and can therefore be resolved by the EOF
// probe. Any other non-62 code (e.g. LOGICAL_ERROR, TIMEOUT_EXCEEDED,
// MEMORY_LIMIT_EXCEEDED) is treated as a judge failure rather than a verdict.
func isKnownSemanticCode(code int32) bool {
	switch code {
	case chErrBadArguments,
		chErrUnknownStorage,
		chErrUnknownSetting,
		chErrCannotParseUUID,
		chErrNotImplemented,
		chErrExpectedAllOrDistinct:
		return true
	default:
		return false
	}
}

// compatJudge asks the pinned ClickHouse server whether its parser accepts each
// statement via EXPLAIN AST without executing the statement.
//
// Because ClickHouse's ParserExplainQuery accepts optional bare `key = value`
// settings between `EXPLAIN AST` and the inner query, every candidate is
// checked with both:
//  1. `EXPLAIN AST <stmt>` (rejects leading-comma fragments like `, a = 1 SELECT 1`)
//  2. `EXPLAIN AST graph = 0 <stmt>` (rejects leading `key = value` fragments
//     like `graph = 1 SELECT 1` or `EXPLAIN = 1 SELECT *` that Step 1 would
//     absorb as EXPLAIN AST settings)
//  3. When both return the same allowlisted non-62 semantic error, a third
//     probe `EXPLAIN AST graph = 0 <stmt>\n; SELECT 1` separates parse-time
//     rejections (same error code -> reject) from post-parse checks
//     (SYNTAX_ERROR 62 at `; SELECT 1` -> accept).
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	db := testcontainer.Judge(t, pin)
	runExplain := func(query string) error {
		ctx, cancel := context.WithTimeout(context.Background(), clickhouseQueryTimeout)
		defer cancel()
		rows, err := db.QueryContext(ctx, query)
		if err == nil {
			_ = rows.Close()
		}
		return err
	}
	judge := func(stmt string) (string, error) {
		if strings.ContainsRune(stmt, 0) {
			return "", fmt.Errorf("statement contains a zero byte")
		}
		return clickhouseVerdict(
			func() error { return runExplain("EXPLAIN AST " + stmt) },
			func() error { return runExplain("EXPLAIN AST graph = 0 " + stmt) },
			func() error { return runExplain("EXPLAIN AST graph = 0 " + stmt + clickhouseProbeTail) },
		)
	}

	// Sanity-check the judge contract before recording or verifying rows.
	assertProbe := func(stmt, want string) {
		got, err := judge(stmt)
		if err != nil || got != want {
			t.Fatalf("ClickHouse judge sanity check on %q = (%q, %v), want (%q, nil)", stmt, got, err, want)
		}
	}
	assertProbe("SELECT 1", compat.Accept)
	assertProbe("SELECT 1 +", compat.Reject)
	assertProbe("SELECT * SETTINGS a = 1", compat.Accept)
	assertProbe("EXPLAIN = 1 SELECT *", compat.Reject)
	assertProbe("graph = 1 SELECT 1", compat.Reject)
	assertProbe(", a = 1 SELECT 1", compat.Reject)
	assertProbe("SELECT * LIMIT 1 WITH TIES", compat.Reject)

	return judge
}

// clickhouseVerdict classifies ClickHouse's EXPLAIN AST responses fail-closed.
func clickhouseVerdict(explain, explainFixed, probeEOF func() error) (string, error) {
	err1 := explain()
	var code1 int32
	if err1 != nil {
		var e *clickhouse.Exception
		if !errors.As(err1, &e) {
			return "", err1
		}
		if e.Code == chErrSyntaxError {
			return compat.Reject, nil
		}
		code1 = e.Code
	}

	err2 := explainFixed()
	var code2 int32
	if err2 != nil {
		var e *clickhouse.Exception
		if !errors.As(err2, &e) {
			return "", err2
		}
		if e.Code == chErrSyntaxError {
			return compat.Reject, nil
		}
		code2 = e.Code
	}

	if err1 == nil && err2 == nil {
		return compat.Accept, nil
	}
	if err1 == nil || err2 == nil || code1 != code2 {
		return "", fmt.Errorf("inconsistent EXPLAIN AST results: bare=%v fixed=%v", err1, err2)
	}
	if !isKnownSemanticCode(code2) {
		return "", fmt.Errorf("unexpected ClickHouse error code %d: %v", code2, err2)
	}

	err3 := probeEOF()
	if err3 == nil {
		return "", fmt.Errorf("EOF probe unexpectedly succeeded after ClickHouse code %d (%v)", code2, err2)
	}
	var e3 *clickhouse.Exception
	if !errors.As(err3, &e3) {
		return "", fmt.Errorf("EOF probe failed: %w", err3)
	}
	switch e3.Code {
	case chErrSyntaxError:
		return compat.Accept, nil
	case code2:
		return compat.Reject, nil
	default:
		return "", fmt.Errorf("EOF probe returned unexpected ClickHouse code %d (%s) after code %d (%v)",
			e3.Code, e3.Message, code2, err2)
	}
}

func TestClickHouseVerdict(t *testing.T) {
	syntaxErr := &clickhouse.Exception{Code: chErrSyntaxError, Message: "Syntax error"}
	settingErr := &clickhouse.Exception{Code: chErrUnknownSetting, Message: "Unknown setting"}
	storageErr := &clickhouse.Exception{Code: chErrUnknownStorage, Message: "Unknown table engine"}
	uuidErr := &clickhouse.Exception{Code: chErrCannotParseUUID, Message: "Cannot parse uuid"}
	tiesErr := &clickhouse.Exception{Code: chErrNotImplemented, Message: "Can not use WITH TIES without ORDER BY"}
	frameErr := &clickhouse.Exception{Code: chErrBadArguments, Message: "Frame start cannot be UNBOUNDED FOLLOWING"}
	unionErr := &clickhouse.Exception{Code: chErrExpectedAllOrDistinct, Message: "Expected ALL or DISTINCT"}
	logicalErr := &clickhouse.Exception{Code: 49, Message: "Logical error"}
	transportErr := context.DeadlineExceeded

	cases := []struct {
		name      string
		err1      error
		err2      error
		err3      error
		want      string
		wantErr   bool
		wantCalls [3]int
	}{
		{
			name:      "both explain succeed -> accept",
			want:      compat.Accept,
			wantCalls: [3]int{1, 1, 0},
		},
		{
			name:      "bare explain syntax error -> reject immediately",
			err1:      syntaxErr,
			want:      compat.Reject,
			wantCalls: [3]int{1, 0, 0},
		},
		{
			name:      "graph=1 SELECT 1 absorbed by bare explain, rejected by fixed explain",
			err1:      nil,
			err2:      syntaxErr,
			want:      compat.Reject,
			wantCalls: [3]int{1, 1, 0},
		},
		{
			name:      "EXPLAIN = 1 SELECT * unknown setting in bare, syntax error in fixed -> reject",
			err1:      settingErr,
			err2:      syntaxErr,
			want:      compat.Reject,
			wantCalls: [3]int{1, 1, 0},
		},
		{
			name:      "post-parse unknown setting resolved to accept by EOF probe",
			err1:      settingErr,
			err2:      settingErr,
			err3:      syntaxErr,
			want:      compat.Accept,
			wantCalls: [3]int{1, 1, 1},
		},
		{
			name:      "post-parse unknown storage resolved to accept by EOF probe",
			err1:      storageErr,
			err2:      storageErr,
			err3:      syntaxErr,
			want:      compat.Accept,
			wantCalls: [3]int{1, 1, 1},
		},
		{
			name:      "post-parse bare UNION resolved to accept by EOF probe",
			err1:      unionErr,
			err2:      unionErr,
			err3:      syntaxErr,
			want:      compat.Accept,
			wantCalls: [3]int{1, 1, 1},
		},
		{
			name:      "parse-time UUID error stays 376 on EOF probe -> reject",
			err1:      uuidErr,
			err2:      uuidErr,
			err3:      uuidErr,
			want:      compat.Reject,
			wantCalls: [3]int{1, 1, 1},
		},
		{
			name:      "parse-time WITH TIES without ORDER BY stays 476 -> reject",
			err1:      tiesErr,
			err2:      tiesErr,
			err3:      tiesErr,
			want:      compat.Reject,
			wantCalls: [3]int{1, 1, 1},
		},
		{
			name:      "parse-time window frame UNBOUNDED FOLLOWING stays 36 -> reject",
			err1:      frameErr,
			err2:      frameErr,
			err3:      frameErr,
			want:      compat.Reject,
			wantCalls: [3]int{1, 1, 1},
		},
		{
			name:      "non-allowlisted server error fails closed",
			err1:      logicalErr,
			err2:      logicalErr,
			wantErr:   true,
			wantCalls: [3]int{1, 1, 0},
		},
		{
			name:      "transport error on step 1 fails closed",
			err1:      transportErr,
			wantErr:   true,
			wantCalls: [3]int{1, 0, 0},
		},
		{
			name:      "transport error on step 2 fails closed",
			err1:      nil,
			err2:      transportErr,
			wantErr:   true,
			wantCalls: [3]int{1, 1, 0},
		},
		{
			name:      "inconsistent step1/step2 codes fail closed",
			err1:      settingErr,
			err2:      storageErr,
			wantErr:   true,
			wantCalls: [3]int{1, 1, 0},
		},
		{
			name:      "EOF probe unexpectedly succeeding fails closed",
			err1:      settingErr,
			err2:      settingErr,
			err3:      nil,
			wantErr:   true,
			wantCalls: [3]int{1, 1, 1},
		},
		{
			name:      "EOF probe returning unexpected code fails closed",
			err1:      settingErr,
			err2:      settingErr,
			err3:      logicalErr,
			wantErr:   true,
			wantCalls: [3]int{1, 1, 1},
		},
		{
			name:      "EOF probe returning transport error fails closed",
			err1:      settingErr,
			err2:      settingErr,
			err3:      transportErr,
			wantErr:   true,
			wantCalls: [3]int{1, 1, 1},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var calls [3]int
			got, err := clickhouseVerdict(
				func() error { calls[0]++; return tc.err1 },
				func() error { calls[1]++; return tc.err2 },
				func() error { calls[2]++; return tc.err3 },
			)
			if calls != tc.wantCalls {
				t.Errorf("calls = %v, want %v", calls, tc.wantCalls)
			}
			if tc.wantErr {
				if err == nil {
					t.Fatalf("clickhouseVerdict() = (%q, nil), want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("clickhouseVerdict() unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("clickhouseVerdict() = %q, want %q", got, tc.want)
			}
		})
	}
}
