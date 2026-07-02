// Package trinooracle is the differential-testing oracle for the omni Trino
// parser. It talks to a live Trino server over the REST statement protocol and
// classifies a statement as syntactically ACCEPTED or REJECTED by Trino's own
// parser.
//
// The classification rule is the heart of the oracle: Trino raises errorName
// "SYNTAX_ERROR" (errorCode 1, errorType USER_ERROR) when, and only when, its
// parser rejects the input. Every other outcome — success, TABLE_NOT_FOUND,
// CATALOG_NOT_FOUND, COLUMN_NOT_FOUND, NOT_SUPPORTED, TYPE_MISMATCH, … — means
// the grammar ACCEPTED the statement and it failed (or succeeded) at a later,
// semantic stage. This mirrors the MySQL oracle's use of error code 1064.
//
// A grammar node wires the differential like so (in package parser, test only):
//
//	func assertOracleMatch(t *testing.T, o *trinooracle.Oracle, sql string) {
//	    _, omniErr := Parse(sql)
//	    res, err := o.CheckSyntax(context.Background(), sql)
//	    if err != nil { t.Skipf("oracle unreachable: %v", err) }
//	    if (omniErr == nil) != res.Accepted {
//	        t.Errorf("MISMATCH %q: omni accepts=%v (err=%v), trino accepts=%v (%s)",
//	            sql, omniErr == nil, omniErr, res.Accepted, res.ErrorName)
//	    }
//	}
//
// Because classification keys on SYNTAX_ERROR, the oracle may *execute* a DDL or
// DML statement (e.g. against the memory catalog) as a side effect; that is
// harmless for syntax classification and the oracle container is disposable.
package trinooracle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// DefaultURL is the conventional local Trino oracle endpoint. Start one with:
//
//	docker run -d --name trino-oracle -p 18080:8080 trinodb/trino:latest
const DefaultURL = "http://localhost:18080"

// EnvURL is the environment variable that overrides the oracle endpoint.
const EnvURL = "TRINO_ORACLE_URL"

// SkipOrFailUnreachable is the standard guard when the oracle cannot be
// reached at test start. Locally it skips so developers without a Trino
// aren't blocked; in CI (CI env var set, as on GitHub Actions) it fails
// hard — the CI container job starts a Trino first, so unreachable there
// means the gate would otherwise silently turn into a green no-op.
func SkipOrFailUnreachable(t interface {
	Helper()
	Fatalf(format string, args ...any)
	Skipf(format string, args ...any)
}, format string, args ...any) {
	t.Helper()
	if os.Getenv("CI") != "" {
		t.Fatalf(format, args...)
	}
	t.Skipf(format, args...)
}

// syntaxErrorName is Trino's errorName for a parser rejection.
const syntaxErrorName = "SYNTAX_ERROR"

// Oracle is a client for a live Trino server used as a syntax oracle.
type Oracle struct {
	baseURL string
	user    string
	catalog string
	schema  string
	client  *http.Client
}

// Option configures an Oracle.
type Option func(*Oracle)

// WithUser sets the Trino user sent in X-Trino-User.
func WithUser(user string) Option { return func(o *Oracle) { o.user = user } }

// WithDefaultNamespace sets the default catalog and schema so that
// catalog/schema-relative statements resolve far enough to distinguish syntax
// errors from semantic ones.
func WithDefaultNamespace(catalog, schema string) Option {
	return func(o *Oracle) { o.catalog, o.schema = catalog, schema }
}

// WithHTTPClient overrides the HTTP client (e.g. to change timeouts).
func WithHTTPClient(c *http.Client) Option { return func(o *Oracle) { o.client = c } }

// Connect returns an Oracle pointed at baseURL. If baseURL is empty it falls
// back to $TRINO_ORACLE_URL and then to DefaultURL. Connect does not perform
// any I/O; call Ping to verify reachability.
func Connect(baseURL string, opts ...Option) *Oracle {
	if baseURL == "" {
		baseURL = URLFromEnv()
	}
	if baseURL == "" {
		baseURL = DefaultURL
	}
	o := &Oracle{
		baseURL: strings.TrimRight(baseURL, "/"),
		user:    "omni-oracle",
		catalog: "memory",
		schema:  "default",
		client:  &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// URLFromEnv returns the value of $TRINO_ORACLE_URL (possibly empty).
func URLFromEnv() string { return os.Getenv(EnvURL) }

// Result is the outcome of a syntax check.
type Result struct {
	// Accepted is true if Trino's parser accepted the statement (no SYNTAX_ERROR).
	Accepted bool
	// ErrorName is Trino's errorName, "" when the statement ran without error.
	ErrorName string
	// ErrorCode is Trino's numeric errorCode (0 when no error).
	ErrorCode int
	// Message is Trino's error message, if any.
	Message string
}

// info models the GET /v1/info response.
type info struct {
	Starting    bool `json:"starting"`
	NodeVersion struct {
		Version string `json:"version"`
	} `json:"nodeVersion"`
}

// queryResp models the relevant fields of a /v1/statement response.
type queryResp struct {
	ID      string `json:"id"`
	NextURI string `json:"nextUri"`
	Stats   struct {
		State string `json:"state"`
	} `json:"stats"`
	Error *struct {
		Message   string `json:"message"`
		ErrorCode int    `json:"errorCode"`
		ErrorName string `json:"errorName"`
		ErrorType string `json:"errorType"`
	} `json:"error"`
}

// Ping verifies the server is reachable and finished starting, returning the
// server version on success.
func (o *Oracle) Ping(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, o.baseURL+"/v1/info", nil)
	if err != nil {
		return "", err
	}
	resp, err := o.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("trino /v1/info: status %d", resp.StatusCode)
	}
	var in info
	if err := json.NewDecoder(resp.Body).Decode(&in); err != nil {
		return "", err
	}
	if in.Starting {
		return in.NodeVersion.Version, fmt.Errorf("trino is still starting")
	}
	return in.NodeVersion.Version, nil
}

// CheckSyntax submits sql to Trino and reports whether its parser accepted it.
// It follows the statement to the point where parsing and planning are known to
// have completed, then cancels any remaining execution so large result sets are
// not streamed.
func (o *Oracle) CheckSyntax(ctx context.Context, sql string) (Result, error) {
	resp, err := o.post(ctx, sql)
	if err != nil {
		return Result{}, err
	}
	for i := 0; i < 10000; i++ {
		if resp.Error != nil {
			return Result{
				Accepted:  resp.Error.ErrorName != syntaxErrorName,
				ErrorName: resp.Error.ErrorName,
				ErrorCode: resp.Error.ErrorCode,
				Message:   resp.Error.Message,
			}, nil
		}
		// Once the query reaches RUNNING or FINISHED, parsing+analysis+planning
		// have all succeeded, so the statement was syntactically accepted.
		switch resp.Stats.State {
		case "RUNNING", "FINISHED":
			if resp.NextURI != "" {
				o.cancel(resp.NextURI)
			}
			return Result{Accepted: true}, nil
		}
		if resp.NextURI == "" {
			return Result{Accepted: true}, nil
		}
		select {
		case <-ctx.Done():
			return Result{}, ctx.Err()
		case <-time.After(20 * time.Millisecond):
		}
		if resp, err = o.get(ctx, resp.NextURI); err != nil {
			return Result{}, err
		}
	}
	return Result{}, fmt.Errorf("trino: statement %q did not settle", truncate(sql, 60))
}

// Accepts is a convenience wrapper returning only the accept/reject verdict.
func (o *Oracle) Accepts(ctx context.Context, sql string) (bool, error) {
	res, err := o.CheckSyntax(ctx, sql)
	if err != nil {
		return false, err
	}
	return res.Accepted, nil
}

func (o *Oracle) post(ctx context.Context, sql string) (queryResp, error) {
	// Trino can answer 503 SERVER_STARTING_UP briefly; retry a few times.
	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/v1/statement", bytes.NewBufferString(sql))
		if err != nil {
			return queryResp{}, err
		}
		req.Header.Set("X-Trino-User", o.user)
		req.Header.Set("X-Trino-Catalog", o.catalog)
		req.Header.Set("X-Trino-Schema", o.schema)
		req.Header.Set("Content-Type", "text/plain")
		resp, err := o.client.Do(req)
		if err != nil {
			return queryResp{}, err
		}
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body.Close()
			lastErr = fmt.Errorf("trino /v1/statement: 503")
			select {
			case <-ctx.Done():
				return queryResp{}, ctx.Err()
			case <-time.After(time.Second):
			}
			continue
		}
		return decode(resp)
	}
	return queryResp{}, lastErr
}

func (o *Oracle) get(ctx context.Context, uri string) (queryResp, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return queryResp{}, err
	}
	req.Header.Set("X-Trino-User", o.user)
	resp, err := o.client.Do(req)
	if err != nil {
		return queryResp{}, err
	}
	return decode(resp)
}

// cancel best-effort terminates a running query so result data is not streamed.
func (o *Oracle) cancel(uri string) {
	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		return
	}
	req.Header.Set("X-Trino-User", o.user)
	if resp, err := o.client.Do(req); err == nil {
		resp.Body.Close()
	}
}

func decode(resp *http.Response) (queryResp, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return queryResp{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return queryResp{}, fmt.Errorf("trino: status %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	var qr queryResp
	if err := json.Unmarshal(body, &qr); err != nil {
		return queryResp{}, fmt.Errorf("trino: decode response: %w", err)
	}
	return qr, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
