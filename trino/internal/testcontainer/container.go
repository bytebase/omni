// Package testcontainer is the differential-testing oracle for the omni Trino
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
//	func assertOracleMatch(t *testing.T, o *testcontainer.Oracle, sql string) {
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
package testcontainer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"

	"github.com/bytebase/omni/enginetest"
)

// DefaultURL is the conventional local Trino oracle endpoint. Start one with:
//
//	docker run -d --name trino-oracle -p 18080:8080 trinodb/trino:482
const DefaultURL = "http://localhost:18080"

// EnvURL is the environment variable that overrides the oracle endpoint.
const EnvURL = "TRINO_ORACLE_URL"

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
		Message       string `json:"message"`
		ErrorCode     int    `json:"errorCode"`
		ErrorName     string `json:"errorName"`
		ErrorType     string `json:"errorType"`
		ErrorLocation *struct {
			LineNumber   int `json:"lineNumber"`
			ColumnNumber int `json:"columnNumber"`
		} `json:"errorLocation"`
	} `json:"error"`
}

// QueryError is a Trino server's answer to a statement that failed, as its
// client protocol's QueryError carries it.
type QueryError struct {
	// Name is errorName and Code errorCode: SYNTAX_ERROR and 1 for a
	// rejection, TABLE_NOT_FOUND for a missing table.
	Name string
	Code int
	// Line and Column are errorLocation, from 1; 0 for none.
	Line, Column int
	Message      string
}

func (e *QueryError) Error() string { return fmt.Sprintf("trino %s: %s", e.Name, e.Message) }

// queryError reads the error of a response, or nil.
func (r queryResp) queryError() *QueryError {
	if r.Error == nil {
		return nil
	}
	e := &QueryError{Name: r.Error.ErrorName, Code: r.Error.ErrorCode, Message: r.Error.Message}
	if l := r.Error.ErrorLocation; l != nil {
		e.Line, e.Column = l.LineNumber, l.ColumnNumber
	}
	return e
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
func (o *Oracle) CheckSyntax(ctx context.Context, sql string) (Result, error) {
	err := o.Run(ctx, sql)
	var qe *QueryError
	switch {
	case err == nil:
		return Result{Accepted: true}, nil
	case errors.As(err, &qe):
		return Result{
			Accepted:  qe.Name != syntaxErrorName,
			ErrorName: qe.Name,
			ErrorCode: qe.Code,
			Message:   qe.Message,
		}, nil
	default:
		return Result{}, err
	}
}

// Run submits sql to Trino and follows it to the point where parsing and
// planning are known to have completed, then cancels any remaining
// execution so large result sets are not streamed. It returns nil when the
// statement got that far, the server's answer as a *QueryError when it
// failed, and any other error when the server gave no answer.
func (o *Oracle) Run(ctx context.Context, sql string) error {
	resp, err := o.post(ctx, sql)
	if err != nil {
		return err
	}
	for i := 0; i < 10000; i++ {
		if qe := resp.queryError(); qe != nil {
			return qe
		}
		// Once the query reaches RUNNING or FINISHED, parsing+analysis+planning
		// have all succeeded, so the statement was syntactically accepted.
		switch resp.Stats.State {
		case "RUNNING", "FINISHED":
			if resp.NextURI != "" {
				o.cancel(resp.NextURI)
			}
			return nil
		}
		if resp.NextURI == "" {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(20 * time.Millisecond):
		}
		if resp, err = o.get(ctx, resp.NextURI); err != nil {
			return err
		}
	}
	return fmt.Errorf("trino: statement %q did not settle", truncate(sql, 60))
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

// Pins is trino/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// memoryLimit caps the container: Trino sizes its JVM to the memory it
// sees, and the box it records on is shared.
const memoryLimit = 4 << 30

// A server is the one pinned Trino this run started.
var server struct {
	once      sync.Once
	pin       enginetest.Pin
	container testcontainers.Container
	url       string
	err       error
}

// started returns the URL of the pinned Trino, starting it on the first
// call and failing t unless it runs pin.
func started(t testing.TB, pin enginetest.Pin) string {
	t.Helper()
	server.once.Do(func() {
		server.pin = pin
		server.url, server.err = start(context.Background(), pin)
	})
	if server.err != nil {
		t.Fatalf("start Trino: %v", server.err)
	}
	if server.pin != pin {
		t.Fatalf("this run started Trino at %s, not %s", server.pin, pin)
	}
	return server.url
}

func start(ctx context.Context, pin enginetest.Pin) (url string, err error) {
	// testcontainers panics when Docker is missing; fail the test instead.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pin.Image,
			ExposedPorts: []string{"8080/tcp"},
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.Memory = memoryLimit
			},
		},
		Started: true,
	})
	if c != nil {
		server.container = c
	}
	if err != nil {
		return "", err
	}
	if err := checkDigest(ctx, c, pin.Digest); err != nil {
		return "", err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}
	port, err := c.MappedPort(ctx, "8080/tcp")
	if err != nil {
		return "", err
	}
	url = fmt.Sprintf("http://%s:%s", host, port.Port())
	o := Connect(url)
	// GenericContainer returns once the process is up; poll until Trino
	// reports it has finished starting.
	deadline := time.Now().Add(3 * time.Minute)
	for {
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		version, perr := o.Ping(pingCtx)
		cancel()
		if perr == nil {
			if version != pin.Tag {
				return "", fmt.Errorf("the server reports version %q, but pins.tsv pins Trino %s", version, pin.Tag)
			}
			return url, nil
		}
		if time.Now().After(deadline) {
			return "", fmt.Errorf("trino did not become ready within timeout: %w", perr)
		}
		time.Sleep(2 * time.Second)
	}
}

// checkDigest fails unless c runs the image with the pinned repository
// digest. A tag can be pushed again; a digest cannot.
func checkDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the Trino container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the Trino image: %w", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the Trino container runs an image with digests %v, but pins.tsv pins %s", img.RepoDigests, digest)
}

// Terminate removes the server this run started, if any. A package's
// TestMain calls it.
func Terminate() {
	if server.container != nil {
		_ = testcontainers.TerminateContainer(server.container)
	}
}

// Judge returns the pinned Trino itself, for the judge of compat.tsv, which
// runs only in a recording.
func Judge(t testing.TB) *Oracle {
	t.Helper()
	if !enginetest.Recording(t) {
		t.Fatalf("testcontainer.Judge outside a recording: a replay reaches no server")
	}
	pin, err := enginetest.LoadPin(Pins, "default")
	if err != nil {
		t.Fatalf("%v", err)
	}
	return Connect(started(t, pin))
}

var forTest struct {
	once sync.Once
	o    *Oracle
	err  error
}

// ForTest returns the Oracle the differential tests run against: the server at
// $TRINO_ORACLE_URL or the conventional localhost:18080 when one answers, else
// a Trino testcontainer started once per test binary. Docker is assumed, so a
// server that can be neither reached nor started fails the test.
func ForTest(t interface {
	Helper()
	Fatalf(format string, args ...any)
}) *Oracle {
	t.Helper()
	forTest.once.Do(func() {
		o := Connect("")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := o.Ping(ctx)
		cancel()
		if err == nil {
			forTest.o = o
			return
		}
		if URLFromEnv() != "" {
			forTest.err = fmt.Errorf("trino at %s unreachable: %w", o.baseURL, err)
			return
		}
		pin, err := enginetest.LoadPin(Pins, "default")
		if err != nil {
			forTest.err = err
			return
		}
		server.once.Do(func() {
			server.pin = pin
			server.url, server.err = start(context.Background(), pin)
		})
		forTest.o, forTest.err = Connect(server.url), server.err
	})
	if forTest.err != nil {
		t.Fatalf("trino oracle: %v", forTest.err)
	}
	return forTest.o
}
