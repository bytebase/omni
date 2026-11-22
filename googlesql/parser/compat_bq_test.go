package parser

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/googlesql/internal/testcontainer"
)

const (
	bqCredentialEnv     = "BIGQUERY_CREDENTIALS"
	bqProjectEnv        = "BIGQUERY_PROJECT"
	bqCredentialDefault = "default"
	bqPinRole           = "bigquery"

	bqDryRunTimeout     = 15 * time.Second
	bqDryRunMaxAttempts = 5
)

type bqDryRunConfig struct {
	Project  string
	Password string // resolved credential JSON, or empty for ADC
}

// bigQueryDryRunJudge opens and pin-checks the live BigQuery dry-run judge for
// the BigQuery half of googlesql/parser's union judge. It runs only during
// OMNI_RECORD=1 or OMNI_RECORD=verify, executes dry-run jobs only (no dataset
// created or mutated), and fails closed if credentials, project, or pin do not
// match googlesql/testdata/pins.tsv's bigquery row.
func bigQueryDryRunJudge(t testing.TB) func(string) (string, error) {
	t.Helper()
	bqPin, err := enginetest.LoadPin(testcontainer.Pins, bqPinRole)
	if err != nil {
		t.Fatalf("bigquery union judge: load %q pin: %v", bqPinRole, err)
	}
	cfg, err := resolveBigQueryDryRunConfig(os.Getenv)
	if err != nil {
		t.Fatalf("bigquery union judge: %v", err)
	}

	initCtx, initCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer initCancel()

	client, err := newBigQueryDryRunClient(initCtx, cfg)
	if err != nil {
		t.Fatalf("bigquery union judge: create client: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	if err := verifyBigQueryDryRunPin(initCtx, client, bqPin, cfg.Project); err != nil {
		t.Fatalf("bigquery union judge: %v", err)
	}

	askBQ := func(sql string) (string, error) {
		return bqDryRunVerdict(client, sql)
	}
	// Startup self-checks: pin the live BigQuery syntax vs semantic behavior
	// that the union judge relies on.
	for _, probe := range []struct {
		sql  string
		want string
	}{
		{"SELECT a FROM t QUALIFY ROW_NUMBER() OVER () = 1", compat.Accept},
		{"FROM a", compat.Accept},
		{"SELEKT 1", compat.Reject},
		{"SELECT EXCEPT", compat.Reject},
	} {
		got, err := askBQ(probe.sql)
		if err != nil || got != probe.want {
			t.Fatalf("compatJudge bigquery startup probe %q: got (%q, %v), want (%q, nil)", probe.sql, got, err, probe.want)
		}
	}
	return askBQ
}

func resolveBigQueryDryRunConfig(lookup func(string) string) (bqDryRunConfig, error) {
	credentials := strings.TrimSpace(lookup(bqCredentialEnv))
	if credentials == "" {
		return bqDryRunConfig{}, fmt.Errorf("recording needs %s (a key path, key JSON, or %q)", bqCredentialEnv, bqCredentialDefault)
	}
	project := strings.TrimSpace(lookup(bqProjectEnv))
	var document string
	if credentials == bqCredentialDefault {
		if project == "" {
			return bqDryRunConfig{}, fmt.Errorf("%s is %q and names no key to read a project from, so %s must be set", bqCredentialEnv, bqCredentialDefault, bqProjectEnv)
		}
	} else {
		document = credentials
		if !strings.HasPrefix(credentials, "{") {
			key, err := os.ReadFile(credentials)
			if err != nil {
				return bqDryRunConfig{}, fmt.Errorf("%s names a file that cannot be read: %w", bqCredentialEnv, err)
			}
			document = string(key)
		}
		if project == "" {
			var key struct {
				ProjectID string `json:"project_id"`
			}
			if err := json.Unmarshal([]byte(document), &key); err != nil {
				return bqDryRunConfig{}, fmt.Errorf("%s does not parse as a credential: %w", bqCredentialEnv, err)
			}
			project = strings.TrimSpace(key.ProjectID)
			if project == "" {
				return bqDryRunConfig{}, fmt.Errorf("%s is not set and the key names no project_id", bqProjectEnv)
			}
		}
	}
	return bqDryRunConfig{Project: project, Password: document}, nil
}

func newBigQueryDryRunClient(ctx context.Context, cfg bqDryRunConfig) (*bigquery.Client, error) {
	var opts []option.ClientOption
	if cfg.Password != "" {
		credOpt, err := bqCredentialOption([]byte(cfg.Password))
		if err != nil {
			return nil, err
		}
		opts = append(opts, credOpt)
	}
	return bigquery.NewClient(ctx, cfg.Project, opts...)
}

func bqCredentialOption(credentials []byte) (option.ClientOption, error) {
	var credential struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(credentials, &credential); err != nil {
		return nil, fmt.Errorf("failed to parse the GCP credential JSON: %w", err)
	}
	switch credential.Type {
	case "service_account":
		return option.WithAuthCredentialsJSON(option.ServiceAccount, credentials), nil
	case "external_account":
		return option.WithAuthCredentialsJSON(option.ExternalAccount, credentials), nil
	case "impersonated_service_account":
		return option.WithAuthCredentialsJSON(option.ImpersonatedServiceAccount, credentials), nil
	case "authorized_user":
		return option.WithAuthCredentialsJSON(option.AuthorizedUser, credentials), nil
	case "":
		return nil, errors.New(`the GCP credential JSON holds no "type" member`)
	default:
		return nil, fmt.Errorf("unsupported GCP credential type: %q", credential.Type)
	}
}

func bqDryRunImage(project, location string) string {
	return fmt.Sprintf("gcp:bigquery:%s:%s", project, location)
}

func bqDryRunDigest(project, location string) string {
	sum := sha256.Sum256([]byte(bqDryRunImage(project, location)))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func verifyBigQueryDryRunPin(ctx context.Context, client *bigquery.Client, pin enginetest.Pin, project string) error {
	q := client.Query("SELECT 1")
	q.DryRun = true
	q.UseLegacySQL = false
	job, err := q.Run(ctx)
	if err != nil {
		return fmt.Errorf("dry-run probe against project %q failed: %w", project, err)
	}
	if job == nil {
		return fmt.Errorf("dry-run probe against project %q returned no job", project)
	}
	if st := job.LastStatus(); st != nil && st.Err() != nil {
		return fmt.Errorf("dry-run probe against project %q failed: %w", project, st.Err())
	}
	gotProject := strings.TrimSpace(job.ProjectID())
	if gotProject == "" {
		gotProject = project
	}
	gotLocation := strings.TrimSpace(job.Location())
	if gotLocation == "" {
		return fmt.Errorf("dry-run probe against project %q reported no location", gotProject)
	}
	gotImage := bqDryRunImage(gotProject, gotLocation)
	if gotImage != pin.Image {
		return fmt.Errorf("live dry-run target is %q, but %s pins %q", gotImage, testcontainer.Pins, pin.Image)
	}
	gotDigest := bqDryRunDigest(gotProject, gotLocation)
	if gotDigest != pin.Digest {
		return fmt.Errorf("live dry-run target %s has digest %s, but %s pins %s", gotImage, gotDigest, testcontainer.Pins, pin.Digest)
	}
	return nil
}

// bqPostParse500Pattern matches statements that BigQuery's parser consumes
// cleanly to EOF but that deterministically trigger HTTP 500 internalError
// inside BigQuery's post-parse analyzer:
//   - `DECLARE <var> FUNCTION<...>` and `DECLARE <var> MAP<...>`
//   - `DROP INDEX [IF EXISTS] <name> ON <table>`
//
// We never accept HTTP 500 blindly: confirmBigQueryPostParseInternalError
// verifies live service health (`SELECT 1` succeeds) and runs the syntax-error
// control `sql + " zzz"` to prove BigQuery's parser consumed sql to EOF.
var bqPostParse500Pattern = regexp.MustCompile(`(?i)^\s*(?:DECLARE\s+\S+\s+(?:FUNCTION|MAP)\s*<|DROP\s+INDEX\s+(?:IF\s+EXISTS\s+)?\S+\s+ON\s+\S+\s*$)`)

func bqDryRunOnce(client *bigquery.Client, sql string, timeout time.Duration) error {
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

func bqDryRunVerdict(client *bigquery.Client, sql string) (string, error) {
	backoff := 500 * time.Millisecond
	for attempt := 1; ; attempt++ {
		timeout := bqDryRunTimeout
		if bqPostParse500Pattern.MatchString(sql) {
			// The Go client retries HTTP 500 internally until ctx expires; keep
			// the probe window short so deterministic analyzer 500s fail fast
			// into the health + syntax-control verification below.
			timeout = 5 * time.Second
		}
		err := bqDryRunOnce(client, sql, timeout)
		if isBQInternalServerError(err) && bqPostParse500Pattern.MatchString(sql) {
			if ok, checkErr := confirmBigQueryPostParseInternalError(sql, func(probe string) error {
				return bqDryRunOnce(client, probe, bqDryRunTimeout)
			}); checkErr == nil && ok {
				return compat.Accept, nil
			}
		}
		if err != nil && isRetryableBigQueryError(err) {
			if attempt < bqDryRunMaxAttempts {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			return "", fmt.Errorf("bigquery dry-run failed after %d attempts: %w", bqDryRunMaxAttempts, err)
		}
		return bigQueryVerdict(err)
	}
}

func isBQInternalServerError(err error) bool {
	var gerr *googleapi.Error
	if !errors.As(err, &gerr) || gerr.Code != http.StatusInternalServerError {
		return false
	}
	return hasAnyBQReason(bqGoogleAPIReasons(gerr), "internalError")
}

// confirmBigQueryPostParseInternalError verifies that an HTTP 500 internalError
// on a matching statement is BigQuery's deterministic post-parse analyzer crash
// rather than a service outage:
//  1. A live `SELECT 1` dry-run succeeds (BigQuery is healthy).
//  2. Appending " zzz" produces HTTP 400 "Syntax error: Expected end of input"
//     (proving BigQuery's parser parsed sql completely up to EOF).
func confirmBigQueryPostParseInternalError(sql string, run func(string) error) (bool, error) {
	if !bqPostParse500Pattern.MatchString(sql) {
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
// compat.Reject, or an error (no verdict). Matches bigquery/parser/compat_test.go.
func bigQueryVerdict(err error) (string, error) {
	if err == nil {
		return compat.Accept, nil
	}

	var gerr *googleapi.Error
	if errors.As(err, &gerr) {
		reasons := bqGoogleAPIReasons(gerr)
		if hasBQInfraReason(reasons) {
			return "", fmt.Errorf("bigquery transient/infra error (HTTP %d, reasons=%v): %w", gerr.Code, reasons, err)
		}
		msg := strings.TrimSpace(gerr.Message)
		switch gerr.Code {
		case http.StatusBadRequest:
			if msg == "" {
				return "", fmt.Errorf("bigquery empty 400 error (reasons=%v): %w", reasons, err)
			}
			if !hasAnyBQReason(reasons, "invalidQuery", "invalid") {
				return "", fmt.Errorf("bigquery unclassified 400 error (reasons=%v): %w", reasons, err)
			}
			if isBigQuerySyntaxError(msg) {
				return compat.Reject, nil
			}
			return compat.Accept, nil
		case http.StatusNotFound:
			if hasAnyBQReason(reasons, "notFound") && isBQObjectNotFoundMessage(msg) {
				return compat.Accept, nil
			}
			return "", fmt.Errorf("bigquery non-object 404 error (reasons=%v): %w", reasons, err)
		case http.StatusForbidden:
			if hasAnyBQReason(reasons, "accessDenied") && isBQObjectAccessDeniedMessage(msg) {
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
			if isBQObjectNotFoundMessage(msg) {
				return compat.Accept, nil
			}
		case "accessDenied":
			if isBQObjectAccessDeniedMessage(msg) {
				return compat.Accept, nil
			}
		}
		return "", fmt.Errorf("bigquery unclassified job error (reason=%q): %w", bqerr.Reason, err)
	}

	return "", fmt.Errorf("bigquery non-API error: %w", err)
}

func bqGoogleAPIReasons(gerr *googleapi.Error) []string {
	reasons := make([]string, 0, len(gerr.Errors))
	for _, item := range gerr.Errors {
		if item.Reason != "" {
			reasons = append(reasons, item.Reason)
		}
	}
	return reasons
}

func hasBQInfraReason(reasons []string) bool {
	for _, r := range reasons {
		switch r {
		case "rateLimitExceeded", "quotaExceeded", "backendError", "internalError", "billingNotEnabled", "authError":
			return true
		}
	}
	return false
}

func hasAnyBQReason(reasons []string, want ...string) bool {
	for _, r := range reasons {
		for _, w := range want {
			if r == w {
				return true
			}
		}
	}
	return false
}

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

func isBQObjectNotFoundMessage(msg string) bool {
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

func isBQObjectAccessDeniedMessage(msg string) bool {
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

func TestUnionJudgeVerdict(t *testing.T) {
	t.Run("spanner accept short-circuits bigquery", func(t *testing.T) {
		bqCalled := false
		got, err := unionJudgeVerdict(
			"CREATE TABLE t (id INT64) PRIMARY KEY (id)",
			func(string) (string, error) { return compat.Accept, nil },
			func(string) (string, error) { bqCalled = true; return compat.Reject, nil },
		)
		if err != nil || got != compat.Accept {
			t.Fatalf("got (%q, %v), want (%q, nil)", got, err, compat.Accept)
		}
		if bqCalled {
			t.Fatal("BigQuery judge should not be called when Spanner accepts")
		}
	})

	t.Run("spanner reject + bigquery accept => accept", func(t *testing.T) {
		got, err := unionJudgeVerdict(
			"SELECT a FROM t QUALIFY ROW_NUMBER() OVER () = 1",
			func(string) (string, error) { return compat.Reject, nil },
			func(string) (string, error) { return compat.Accept, nil },
		)
		if err != nil || got != compat.Accept {
			t.Fatalf("got (%q, %v), want (%q, nil)", got, err, compat.Accept)
		}
	})

	t.Run("spanner reject + bigquery reject => reject", func(t *testing.T) {
		got, err := unionJudgeVerdict(
			"SELEKT 1",
			func(string) (string, error) { return compat.Reject, nil },
			func(string) (string, error) { return compat.Reject, nil },
		)
		if err != nil || got != compat.Reject {
			t.Fatalf("got (%q, %v), want (%q, nil)", got, err, compat.Reject)
		}
	})

	t.Run("spanner infra error fails closed without calling bigquery", func(t *testing.T) {
		bqCalled := false
		got, err := unionJudgeVerdict(
			"SELECT 1",
			func(string) (string, error) { return "", errors.New("emulator unavailable") },
			func(string) (string, error) { bqCalled = true; return compat.Accept, nil },
		)
		if err == nil || got != "" {
			t.Fatalf("got (%q, %v), want error", got, err)
		}
		if bqCalled {
			t.Fatal("BigQuery judge should not be called when Spanner errors")
		}
	})

	t.Run("spanner reject + bigquery infra error fails closed", func(t *testing.T) {
		got, err := unionJudgeVerdict(
			"FROM a",
			func(string) (string, error) { return compat.Reject, nil },
			func(string) (string, error) { return "", errors.New("rate limit exceeded") },
		)
		if err == nil || got != "" {
			t.Fatalf("got (%q, %v), want error", got, err)
		}
	})

	t.Run("unexpected verdict fails closed", func(t *testing.T) {
		if _, err := unionJudgeVerdict(
			"SELECT 1",
			func(string) (string, error) { return "maybe", nil },
			func(string) (string, error) { return compat.Accept, nil },
		); err == nil {
			t.Fatal("expected error for unexpected Spanner verdict")
		}
		if _, err := unionJudgeVerdict(
			"SELECT 1",
			func(string) (string, error) { return compat.Reject, nil },
			func(string) (string, error) { return "maybe", nil },
		); err == nil {
			t.Fatal("expected error for unexpected BigQuery verdict")
		}
	})
}

func TestBigQueryPinRowLoads(t *testing.T) {
	pin, err := enginetest.LoadPin(testcontainer.Pins, bqPinRole)
	if err != nil {
		t.Fatalf("LoadPin(%s, %q): %v", testcontainer.Pins, bqPinRole, err)
	}
	if !strings.HasPrefix(pin.Image, "gcp:bigquery:") {
		t.Fatalf("bigquery pin image = %q, want gcp:bigquery:... prefix", pin.Image)
	}
	if !strings.HasPrefix(pin.Digest, "sha256:") {
		t.Fatalf("bigquery pin digest = %q, want sha256:... prefix", pin.Digest)
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
		{name: "nil is accept", err: nil, want: compat.Accept},
		{name: "syntax error", err: apiErr(http.StatusBadRequest, "invalidQuery", `Syntax error: Unexpected identifier "SELEKT" at [1:1]`), want: compat.Reject},
		{name: "bare EXCEPT refusal", err: apiErr(http.StatusBadRequest, "invalidQuery", `EXCEPT must be followed by ALL, DISTINCT, or "(" at [1:8]`), want: compat.Reject},
		{name: "duplicate DYNAMIC LABEL refusal", err: apiErr(http.StatusBadRequest, "invalidQuery", `DYNAMIC LABEL cannot be used more than once in a single element table at [1:59]`), want: compat.Reject},
		{name: "unsupported object type", err: apiErr(http.StatusBadRequest, "invalidQuery", `CHANGE is not a supported object type at [1:8]`), want: compat.Reject},
		{name: "mismatched end label", err: apiErr(http.StatusBadRequest, "invalidQuery", `Mismatched end label; expected l1, got l2 at [1:15]`), want: compat.Reject},
		{name: "EXECUTE IMMEDIATE inner syntax error stays accept", err: apiErr(http.StatusBadRequest, "invalidQuery", "Invalid EXECUTE IMMEDIATE sql string `a`, Syntax error: Unexpected identifier \"a\"\n at [1:19]"), want: compat.Accept},
		{name: "unrecognized name", err: apiErr(http.StatusBadRequest, "invalidQuery", `Unrecognized name: a at [1:8]`), want: compat.Accept},
		{name: "unqualified table", err: apiErr(http.StatusBadRequest, "invalid", `Table "a" must be qualified with a dataset (e.g. dataset.table).`), want: compat.Accept},
		{name: "not found table", err: apiErr(http.StatusNotFound, "notFound", `Not found: Table bytebase-dev:ds.t was not found in location US`), want: compat.Accept},
		{name: "not found project fails closed", err: apiErr(http.StatusNotFound, "notFound", `Not found: Project nonexistent`), wantErr: true},
		{name: "access denied table is accept", err: apiErr(http.StatusForbidden, "accessDenied", `Access Denied: Table other:ds.t: User does not have permission`), want: compat.Accept},
		{name: "access denied project fails closed", err: apiErr(http.StatusForbidden, "accessDenied", `Access Denied: Project bytebase-dev: User does not have bigquery.jobs.create permission`), wantErr: true},
		{name: "rate limit fails closed", err: apiErr(http.StatusTooManyRequests, "rateLimitExceeded", `Exceeded rate limits`), wantErr: true},
		{name: "empty 400 fails closed", err: apiErr(http.StatusBadRequest, "invalidQuery", ``), wantErr: true},
		{name: "bigquery.Error syntax", err: &bigquery.Error{Reason: "invalidQuery", Message: `Syntax error: Unexpected end of script at [1:10]`}, want: compat.Reject},
		{name: "bigquery.Error semantic", err: &bigquery.Error{Reason: "invalidQuery", Message: `Type not found: a at [1:11]`}, want: compat.Accept},
		{name: "non-API error fails closed", err: context.DeadlineExceeded, wantErr: true},
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

func TestConfirmBigQueryPostParseInternalError(t *testing.T) {
	eofSyntaxErr := &googleapi.Error{
		Code:    http.StatusBadRequest,
		Message: `Syntax error: Expected end of input but got identifier "zzz" at [1:21]`,
		Errors:  []googleapi.ErrorItem{{Reason: "invalidQuery"}},
	}

	for _, sql := range []string{
		"DECLARE a MAP < a, a >",
		"DECLARE a FUNCTION < a -> a >",
		"DROP INDEX idx ON T",
		"DROP INDEX IF EXISTS idx ON ds.T",
	} {
		ok, err := confirmBigQueryPostParseInternalError(sql, func(probe string) error {
			switch probe {
			case "SELECT 1":
				return nil
			case sql + " zzz":
				return eofSyntaxErr
			default:
				t.Fatalf("unexpected probe: %q", probe)
				return nil
			}
		})
		if err != nil || !ok {
			t.Fatalf("confirmBigQueryPostParseInternalError(%q) = (%v, %v), want (true, nil)", sql, ok, err)
		}
	}

	// Non-matching statement never qualifies.
	ok, err := confirmBigQueryPostParseInternalError("SELECT 1", func(string) error {
		t.Fatal("should not probe non-matching SQL")
		return nil
	})
	if err != nil || ok {
		t.Fatalf("non-matching = (%v, %v), want (false, nil)", ok, err)
	}

	// Unhealthy service fails closed.
	if ok, err := confirmBigQueryPostParseInternalError("DROP INDEX idx ON T", func(sql string) error {
		if sql == "SELECT 1" {
			return &googleapi.Error{Code: http.StatusServiceUnavailable}
		}
		return eofSyntaxErr
	}); err == nil || ok {
		t.Fatalf("unhealthy service = (%v, %v), want (false, error)", ok, err)
	}

	// Incomplete syntax fails closed.
	if ok, err := confirmBigQueryPostParseInternalError("DECLARE a MAP < a, a", func(sql string) error {
		if sql == "SELECT 1" {
			return nil
		}
		return &googleapi.Error{
			Code:    http.StatusBadRequest,
			Message: `Syntax error: Expected ">" but got identifier "zzz" at [1:22]`,
			Errors:  []googleapi.ErrorItem{{Reason: "invalidQuery"}},
		}
	}); err == nil || ok {
		t.Fatalf("incomplete syntax = (%v, %v), want (false, error)", ok, err)
	}
}

func TestResolveBigQueryDryRunConfig(t *testing.T) {
	lookup := func(m map[string]string) func(string) string {
		return func(k string) string { return m[k] }
	}
	if _, err := resolveBigQueryDryRunConfig(lookup(map[string]string{})); err == nil {
		t.Fatal("expected error when credentials unset")
	}
	if _, err := resolveBigQueryDryRunConfig(lookup(map[string]string{bqCredentialEnv: "default"})); err == nil {
		t.Fatal("expected error when default credentials lack project")
	}
	cfg, err := resolveBigQueryDryRunConfig(lookup(map[string]string{
		bqCredentialEnv: "default",
		bqProjectEnv:    "bytebase-dev",
	}))
	if err != nil || cfg.Project != "bytebase-dev" || cfg.Password != "" {
		t.Fatalf("default config = (%+v, %v), want project=bytebase-dev", cfg, err)
	}

	keyPath := filepath.Join(t.TempDir(), "key.json")
	if err := os.WriteFile(keyPath, []byte(`{"type":"service_account","project_id":"proj-1"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err = resolveBigQueryDryRunConfig(lookup(map[string]string{bqCredentialEnv: keyPath}))
	if err != nil || cfg.Project != "proj-1" || cfg.Password == "" {
		t.Fatalf("file config = (%+v, %v), want project=proj-1", cfg, err)
	}
}
