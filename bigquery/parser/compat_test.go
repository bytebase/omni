package parser

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/googleapi"

	"github.com/bytebase/omni/bigquery/internal/testcontainer"
	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
)

// TestCompatibility holds omni's GoogleSQL parser to testdata/compat.tsv: a
// statement for each alternative of GoogleSQL's grammar
// (../../spanner/testdata/GoogleSQL.g4) and its near-misses, with omni's
// verdict and live BigQuery dry-run's verdict on the pinned Google Cloud
// project (../testdata/pins.tsv). Offline runs read the committed file with no
// credentials and no network access. OMNI_RECORD=1 judges unjudged rows via
// BigQuery dry-run jobs (which execute nothing and mutate no datasets).
func TestCompatibility(t *testing.T) {
	gsqlTD := filepath.Join("..", "..", "spanner", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		Grammar: filepath.Join(gsqlTD, "GoogleSQL.g4"),
		Tokens:  filepath.Join(gsqlTD, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

func omniVerdict(sql string) string {
	if _, err := Parse(sql); err != nil {
		return compat.Reject
	}
	return compat.Accept
}

// compatJudge submits each statement as a BigQuery dry-run job
// (DryRun=true, UseLegacySQL=false) against the pinned project. It runs only in
// OMNI_RECORD=1 or OMNI_RECORD=verify, never during offline replay, and creates
// or mutates no BigQuery datasets or tables.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	cfg := testcontainer.DryRun(t)

	initCtx, initCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer initCancel()

	client, err := testcontainer.DryRunClient(initCtx, cfg)
	if err != nil {
		t.Fatalf("bigquery compat judge: create client: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	if err := testcontainer.VerifyDryRunPin(initCtx, client, pin, cfg.Project); err != nil {
		t.Fatalf("bigquery compat judge: %v", err)
	}

	return func(sql string) (string, error) {
		return dryRunVerdict(client, sql)
	}
}

const (
	dryRunTimeout     = 15 * time.Second
	dryRunMaxAttempts = 5
)

// declareUnsupportedTypePattern matches BigQuery's known post-parse analyzer
// crash where `DECLARE <var> FUNCTION<...>` and `DECLARE <var> MAP<...>` parse
// cleanly to EOF but deterministically fail with HTTP 500 internalError inside
// BigQuery's script validator. We never accept HTTP 500 blindly: when a
// matching statement hits HTTP 500 internalError, confirmPostParseInternalError
// verifies live service health (`SELECT 1` succeeds) and runs the syntax-error
// control `sql + " zzz"` to prove BigQuery's parser consumed sql to EOF.
var declareUnsupportedTypePattern = regexp.MustCompile(`(?i)^\s*DECLARE\s+\S+\s+(?:FUNCTION|MAP)\s*<`)

// dryRunOnce runs a single BigQuery dry-run job with a per-call timeout.
func dryRunOnce(client *bigquery.Client, sql string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	q := client.Query(sql)
	q.DryRun = true
	q.UseLegacySQL = false
	job, err := q.Run(ctx)
	if err == nil && job != nil {
		if st := job.LastStatus(); st != nil && st.Err() != nil {
			err = st.Err()
		}
	}
	return err
}

// dryRunVerdict sends one dry-run query to BigQuery with bounded exponential
// backoff for transient HTTP 429/5xx or network errors only.
func dryRunVerdict(client *bigquery.Client, sql string) (string, error) {
	backoff := 500 * time.Millisecond
	for attempt := 1; ; attempt++ {
		timeout := dryRunTimeout
		if declareUnsupportedTypePattern.MatchString(sql) {
			// The Go client retries HTTP 500 internally until ctx expires; keep
			// the probe window short so deterministic analyzer 500s fail fast
			// into the health + syntax-control verification below.
			timeout = 5 * time.Second
		}
		err := dryRunOnce(client, sql, timeout)
		if isInternalServerError(err) && declareUnsupportedTypePattern.MatchString(sql) {
			if ok, checkErr := confirmPostParseInternalError(sql, func(probe string) error {
				return dryRunOnce(client, probe, dryRunTimeout)
			}); checkErr == nil && ok {
				return compat.Accept, nil
			}
		}
		if err != nil && isRetryableBigQueryError(err) {
			if attempt < dryRunMaxAttempts {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			return "", fmt.Errorf("bigquery dry-run failed after %d attempts: %w", dryRunMaxAttempts, err)
		}
		return bigQueryVerdict(err)
	}
}

func isInternalServerError(err error) bool {
	var gerr *googleapi.Error
	if !errors.As(err, &gerr) || gerr.Code != http.StatusInternalServerError {
		return false
	}
	return hasAnyReason(googleAPIReasons(gerr), "internalError")
}

// confirmPostParseInternalError verifies that an HTTP 500 internalError on a
// `DECLARE ... FUNCTION<...>/MAP<...>` statement is BigQuery's deterministic
// post-parse analyzer crash rather than a service outage:
//  1. A live `SELECT 1` dry-run succeeds (BigQuery is healthy).
//  2. Appending " zzz" produces HTTP 400 "Syntax error: Expected end of input"
//     (proving BigQuery's parser parsed sql completely up to EOF).
func confirmPostParseInternalError(sql string, run func(string) error) (bool, error) {
	if !declareUnsupportedTypePattern.MatchString(sql) {
		return false, nil
	}
	if err := run("SELECT 1"); err != nil {
		return false, fmt.Errorf("health check failed: %w", err)
	}
	ctrlErr := run(sql + " zzz")
	var gerr *googleapi.Error
	if !errors.As(ctrlErr, &gerr) || gerr.Code != http.StatusBadRequest {
		return false, fmt.Errorf("syntax control did not return 400: %v", ctrlErr)
	}
	if !strings.HasPrefix(strings.TrimSpace(gerr.Message), "Syntax error: Expected end of input") {
		return false, fmt.Errorf("syntax control returned unexpected message: %q", gerr.Message)
	}
	return true, nil
}

// isRetryableBigQueryError reports whether err is a transient rate-limit,
// server-side, or transport failure that is safe to retry for a dry-run query.
func isRetryableBigQueryError(err error) bool {
	if err == nil {
		return false
	}
	var gerr *googleapi.Error
	if errors.As(err, &gerr) {
		if gerr.Code == http.StatusTooManyRequests || gerr.Code >= http.StatusInternalServerError {
			return true
		}
		for _, item := range gerr.Errors {
			switch item.Reason {
			case "rateLimitExceeded", "backendError", "internalError":
				return true
			}
		}
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

// bigQueryVerdict classifies a BigQuery dry-run error into compat.Accept,
// compat.Reject, or an error (no verdict). Classification is fail-closed:
//   - nil => Accept.
//   - HTTP 400 with reason "invalidQuery" or "invalid":
//     syntax/lexer/grammar-refusal messages => Reject;
//     all other semantic/analyzer/qualification/unsupported-statement errors => Accept.
//   - HTTP 404 with reason "notFound" for a per-statement object (Table, Dataset,
//     Routine, Model, etc. — never Project) => Accept.
//   - HTTP 403 with reason "accessDenied" for a per-statement object (Table,
//     Dataset, Routine, Model, Connection) => Accept.
//   - Everything else (401, 429, 5xx, project-level 403/404, billing/quota,
//     empty/unknown errors, transport failures) => error (no verdict).
func bigQueryVerdict(err error) (string, error) {
	if err == nil {
		return compat.Accept, nil
	}

	var gerr *googleapi.Error
	if errors.As(err, &gerr) {
		reasons := googleAPIReasons(gerr)
		if hasInfraReason(reasons) {
			return "", fmt.Errorf("bigquery transient/infra error (HTTP %d, reasons=%v): %w", gerr.Code, reasons, err)
		}
		msg := strings.TrimSpace(gerr.Message)
		switch gerr.Code {
		case http.StatusBadRequest:
			if msg == "" {
				return "", fmt.Errorf("bigquery empty 400 error (reasons=%v): %w", reasons, err)
			}
			if !hasAnyReason(reasons, "invalidQuery", "invalid") {
				return "", fmt.Errorf("bigquery unclassified 400 error (reasons=%v): %w", reasons, err)
			}
			if isBigQuerySyntaxError(msg) {
				return compat.Reject, nil
			}
			return compat.Accept, nil
		case http.StatusNotFound:
			if hasAnyReason(reasons, "notFound") && isObjectNotFoundMessage(msg) {
				return compat.Accept, nil
			}
			return "", fmt.Errorf("bigquery non-object 404 error (reasons=%v): %w", reasons, err)
		case http.StatusForbidden:
			if hasAnyReason(reasons, "accessDenied") && isObjectAccessDeniedMessage(msg) {
				return compat.Accept, nil
			}
			return "", fmt.Errorf("bigquery non-object 403 error (reasons=%v): %w", reasons, err)
		default:
			return "", fmt.Errorf("bigquery unexpected HTTP %d error (reasons=%v): %w", gerr.Code, reasons, err)
		}
	}

	var bqerr *bigquery.Error
	if errors.As(err, &bqerr) {
		msg := strings.TrimSpace(bqerr.Message)
		switch bqerr.Reason {
		case "invalidQuery", "invalid":
			if msg == "" {
				return "", fmt.Errorf("bigquery empty %s error: %w", bqerr.Reason, err)
			}
			if isBigQuerySyntaxError(msg) {
				return compat.Reject, nil
			}
			return compat.Accept, nil
		case "notFound":
			if isObjectNotFoundMessage(msg) {
				return compat.Accept, nil
			}
		case "accessDenied":
			if isObjectAccessDeniedMessage(msg) {
				return compat.Accept, nil
			}
		}
		return "", fmt.Errorf("bigquery unclassified job error (reason=%q): %w", bqerr.Reason, err)
	}

	return "", fmt.Errorf("bigquery non-API error: %w", err)
}

func googleAPIReasons(gerr *googleapi.Error) []string {
	reasons := make([]string, 0, len(gerr.Errors))
	for _, item := range gerr.Errors {
		if item.Reason != "" {
			reasons = append(reasons, item.Reason)
		}
	}
	return reasons
}

func hasInfraReason(reasons []string) bool {
	for _, r := range reasons {
		switch r {
		case "rateLimitExceeded", "quotaExceeded", "backendError", "internalError", "billingNotEnabled", "authError":
			return true
		}
	}
	return false
}

func hasAnyReason(reasons []string, want ...string) bool {
	for _, r := range reasons {
		for _, w := range want {
			if r == w {
				return true
			}
		}
	}
	return false
}

// isBigQuerySyntaxError reports whether a BigQuery 400 error message is a
// syntax, lexer, or GoogleSQL parser-action grammar refusal.
//
// Deliberately matches only prefix "Syntax error:" (not substrings), so valid
// outer statements such as `EXECUTE IMMEDIATE 'a'` whose message is
// "Invalid EXECUTE IMMEDIATE sql string `a`, Syntax error: ..." stay Accept.
func isBigQuerySyntaxError(msg string) bool {
	if strings.HasPrefix(msg, "Syntax error:") {
		return true
	}
	for _, re := range grammarRefusals {
		if re.MatchString(msg) {
			return true
		}
	}
	return false
}

// isObjectNotFoundMessage reports whether a 404 message names a missing SQL
// object referenced by the statement rather than a missing GCP project or job.
func isObjectNotFoundMessage(msg string) bool {
	for _, prefix := range []string{
		"Not found: Table ",
		"Not found: Dataset ",
		"Not found: Routine ",
		"Not found: Model ",
		"Not found: Connection ",
		"Not found: PropertyGraph ",
		"Not found: Row access policy ",
		"Not found: Index ",
		"Not found: Search index ",
		"Not found: Vector index ",
		"Not found: Schema ",
	} {
		if strings.HasPrefix(msg, prefix) {
			return true
		}
	}
	return false
}

// isObjectAccessDeniedMessage reports whether a 403 message is a per-statement
// object permission check rather than a project/account/job-creation denial.
func isObjectAccessDeniedMessage(msg string) bool {
	for _, prefix := range []string{
		"Access Denied: Table ",
		"Access Denied: Dataset ",
		"Access Denied: Routine ",
		"Access Denied: Model ",
		"Access Denied: Connection ",
	} {
		if strings.HasPrefix(msg, prefix) {
			return true
		}
	}
	return false
}

func TestParse(t *testing.T) {
	valid := []string{
		"SELECT 1",
		"DECLARE x INT64 DEFAULT 1",
		"CREATE TABLE ds.t (id INT64)",
		"SELECT 1 |> WHERE TRUE",
	}
	for _, sql := range valid {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", sql, err)
		}
	}
	invalid := []string{
		"SELEKT 1",
		"SELECT FROM",
	}
	for _, sql := range invalid {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) expected error, got nil", sql)
		}
	}
}

func TestBigQueryVerdict(t *testing.T) {
	apiErr := func(code int, reason, msg string) error {
		var items []googleapi.ErrorItem
		if reason != "" {
			items = []googleapi.ErrorItem{{Reason: reason, Message: msg}}
		}
		return &googleapi.Error{Code: code, Message: msg, Errors: items}
	}

	cases := []struct {
		name    string
		err     error
		want    string
		wantErr bool
	}{
		{
			name: "nil is accept",
			err:  nil,
			want: compat.Accept,
		},
		{
			name: "syntax error unexpected token",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `Syntax error: Unexpected identifier "SELEKT" at [1:1]`),
			want: compat.Reject,
		},
		{
			name: "syntax error unclosed string literal",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `Syntax error: Unclosed string literal at [1:8]`),
			want: compat.Reject,
		},
		{
			name: "syntax error illegal input character",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `Syntax error: Illegal input character "\001" at [1:8]`),
			want: compat.Reject,
		},
		{
			name: "lookahead refusal bare EXCEPT",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `EXCEPT must be followed by ALL, DISTINCT, or "(" at [1:8]`),
			want: compat.Reject,
		},
		{
			name: "grammar refusal duplicate DYNAMIC LABEL",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `DYNAMIC LABEL cannot be used more than once in a single element table at [1:59]`),
			want: compat.Reject,
		},
		{
			name: "grammar refusal duplicate DYNAMIC PROPERTIES",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `DYNAMIC PROPERTIES cannot be used more than once in a single element table at [1:64]`),
			want: compat.Reject,
		},
		{
			name: "grammar refusal unsupported object type",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `ROLE is not a supported object type at [1:6]`),
			want: compat.Reject,
		},
		{
			name: "grammar refusal empty table",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `A table must define at least one column. at [1:17]`),
			want: compat.Reject,
		},
		{
			name: "grammar refusal FOR UPDATE",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `FOR UPDATE is not supported at [1:10]`),
			want: compat.Reject,
		},
		{
			name: "grammar refusal mismatched end label",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `Mismatched end label; expected l1, got l2 at [1:15]`),
			want: compat.Reject,
		},
		{
			name: "grammar refusal parenthesized SET",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `Using SET with multiple variables requires parentheses around the variable list at [1:5]`),
			want: compat.Reject,
		},
		{
			name: "grammar refusal braced constructors",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `Braced constructors are not supported at [1:8]`),
			want: compat.Reject,
		},
		{
			name: "valid EXECUTE IMMEDIATE with invalid inner string is accept",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", "Invalid EXECUTE IMMEDIATE sql string `a`, Syntax error: Unexpected identifier \"a\"\n at [1:19]"),
			want: compat.Accept,
		},
		{
			name: "semantic unrecognized name",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `Unrecognized name: a at [1:8]`),
			want: compat.Accept,
		},
		{
			name: "semantic statement not supported",
			err:  apiErr(http.StatusBadRequest, "invalidQuery", `Statement not supported: SubpipelineStatement at [1:1]`),
			want: compat.Accept,
		},
		{
			name: "semantic unqualified table invalid reason",
			err:  apiErr(http.StatusBadRequest, "invalid", `Table "a" must be qualified with a dataset (e.g. dataset.table).`),
			want: compat.Accept,
		},
		{
			name: "semantic not found table",
			err:  apiErr(http.StatusNotFound, "notFound", `Not found: Table bytebase-dev:ds.t was not found in location US`),
			want: compat.Accept,
		},
		{
			name: "semantic not found dataset",
			err:  apiErr(http.StatusNotFound, "notFound", `Not found: Dataset bytebase-dev:a`),
			want: compat.Accept,
		},
		{
			name:    "project not found fails closed",
			err:     apiErr(http.StatusNotFound, "notFound", `Not found: Project nonexistent`),
			wantErr: true,
		},
		{
			name: "object access denied table is accept",
			err:  apiErr(http.StatusForbidden, "accessDenied", `Access Denied: Table other:ds.t: User does not have permission`),
			want: compat.Accept,
		},
		{
			name:    "project access denied fails closed",
			err:     apiErr(http.StatusForbidden, "accessDenied", `Access Denied: Project bytebase-dev: User does not have bigquery.jobs.create permission`),
			wantErr: true,
		},
		{
			name:    "billing not enabled fails closed",
			err:     apiErr(http.StatusForbidden, "billingNotEnabled", `Billing has not been enabled for this project`),
			wantErr: true,
		},
		{
			name:    "rate limit fails closed",
			err:     apiErr(http.StatusTooManyRequests, "rateLimitExceeded", `Exceeded rate limits`),
			wantErr: true,
		},
		{
			name:    "backend error fails closed",
			err:     apiErr(http.StatusServiceUnavailable, "backendError", `Backend Error`),
			wantErr: true,
		},
		{
			name:    "unauthorized fails closed",
			err:     apiErr(http.StatusUnauthorized, "authError", `Invalid Credentials`),
			wantErr: true,
		},
		{
			name:    "missing query parameter fails closed",
			err:     apiErr(http.StatusBadRequest, "required", `Required parameter is missing: query`),
			wantErr: true,
		},
		{
			name:    "empty 400 fails closed",
			err:     apiErr(http.StatusBadRequest, "invalidQuery", ``),
			wantErr: true,
		},
		{
			name: "bigquery.Error syntax error",
			err:  &bigquery.Error{Reason: "invalidQuery", Message: `Syntax error: Unexpected end of script at [1:10]`},
			want: compat.Reject,
		},
		{
			name: "bigquery.Error semantic error",
			err:  &bigquery.Error{Reason: "invalidQuery", Message: `Type not found: a at [1:11]`},
			want: compat.Accept,
		},
		{
			name:    "bigquery.Error unknown reason fails closed",
			err:     &bigquery.Error{Reason: "internalError", Message: `An internal error occurred`},
			wantErr: true,
		},
		{
			name:    "non-API error fails closed",
			err:     context.DeadlineExceeded,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := bigQueryVerdict(tc.err)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("bigQueryVerdict() = (%q, nil), want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("bigQueryVerdict() unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("bigQueryVerdict() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestIsRetryableBigQueryError(t *testing.T) {
	if isRetryableBigQueryError(nil) {
		t.Fatal("nil should not be retryable")
	}
	if isRetryableBigQueryError(&googleapi.Error{Code: http.StatusBadRequest, Errors: []googleapi.ErrorItem{{Reason: "invalidQuery"}}}) {
		t.Fatal("400 invalidQuery should not be retryable")
	}
	if !isRetryableBigQueryError(&googleapi.Error{Code: http.StatusTooManyRequests, Errors: []googleapi.ErrorItem{{Reason: "rateLimitExceeded"}}}) {
		t.Fatal("429 should be retryable")
	}
	if !isRetryableBigQueryError(&googleapi.Error{Code: http.StatusServiceUnavailable, Errors: []googleapi.ErrorItem{{Reason: "backendError"}}}) {
		t.Fatal("503 should be retryable")
	}
	if !isRetryableBigQueryError(context.DeadlineExceeded) {
		t.Fatal("context.DeadlineExceeded should be retryable")
	}
	if !isRetryableBigQueryError(io.ErrUnexpectedEOF) {
		t.Fatal("io.ErrUnexpectedEOF should be retryable")
	}
}

func TestConfirmPostParseInternalError(t *testing.T) {
	eofSyntaxErr := &googleapi.Error{
		Code:    http.StatusBadRequest,
		Message: `Syntax error: Expected end of input but got identifier "zzz" at [1:24]`,
		Errors:  []googleapi.ErrorItem{{Reason: "invalidQuery"}},
	}

	ok, err := confirmPostParseInternalError("DECLARE a MAP < a, a >", func(sql string) error {
		switch sql {
		case "SELECT 1":
			return nil
		case "DECLARE a MAP < a, a > zzz":
			return eofSyntaxErr
		default:
			t.Fatalf("unexpected probe: %q", sql)
			return nil
		}
	})
	if err != nil || !ok {
		t.Fatalf("confirmPostParseInternalError() = (%v, %v), want (true, nil)", ok, err)
	}

	// Non-matching statements never qualify.
	ok, err = confirmPostParseInternalError("SELECT 1", func(string) error {
		t.Fatal("should not probe non-matching SQL")
		return nil
	})
	if err != nil || ok {
		t.Fatalf("non-matching statement = (%v, %v), want (false, nil)", ok, err)
	}

	// Unhealthy service fails closed.
	ok, err = confirmPostParseInternalError("DECLARE a FUNCTION < a -> a >", func(sql string) error {
		if sql == "SELECT 1" {
			return &googleapi.Error{Code: http.StatusServiceUnavailable}
		}
		return eofSyntaxErr
	})
	if err == nil || ok {
		t.Fatalf("unhealthy service = (%v, %v), want (false, error)", ok, err)
	}

	// Incomplete syntax (error before EOF) fails closed.
	ok, err = confirmPostParseInternalError("DECLARE a MAP < a, a", func(sql string) error {
		if sql == "SELECT 1" {
			return nil
		}
		return &googleapi.Error{
			Code:    http.StatusBadRequest,
			Message: `Syntax error: Expected ">" but got identifier "zzz" at [1:22]`,
			Errors:  []googleapi.ErrorItem{{Reason: "invalidQuery"}},
		}
	})
	if err == nil || ok {
		t.Fatalf("non-EOF syntax error = (%v, %v), want (false, error)", ok, err)
	}
}
