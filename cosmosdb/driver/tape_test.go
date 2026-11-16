package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/bytebase/omni/cosmosdb/internal/testcontainer"
	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/enginetest/tape"
)

// tapeClient is the driver's client, answered by t's tape. It is what a driver
// test hands the driver in place of a live client; its methods are the calls
// the driver makes, and the recordings hold what a live client answered.
//
// A request is one call. An answer is the names a listing returned, the items
// and page count a query returned, or the client's error.
type tapeClient struct {
	tape *tape.Raw
}

// OpenTape returns t's client to the pinned emulator, listing metadata through the
// client's query API: a replay of its tape, or, in a recording, the emulator
// with every call written down.
func OpenTape(t testing.TB) *tapeClient {
	t.Helper()
	return openTape(t, false)
}

// OpenEmulator returns t's client to the pinned emulator, listing metadata over
// the emulator's REST API the way the driver does when the caller names it.
func OpenTapeEmulator(t testing.TB) *tapeClient {
	t.Helper()
	return openTape(t, true)
}

func openTape(t testing.TB, emulator bool) *tapeClient {
	t.Helper()
	raw := tape.OpenRaw(t, testcontainer.Server, testcontainer.Emulator, func(request string) (string, error) {
		// A recording reaches the emulator on its first call, so a test that
		// replays starts no container.
		return live(t, emulator, request)
	})
	return &tapeClient{tape: raw}
}

// Ping sends a ping.
func (c *tapeClient) Ping(ctx context.Context) error {
	_, err := c.do(ctx, request{Op: opPing})
	return err
}

// Close sends a close.
func (c *tapeClient) Close() error {
	_, err := c.do(context.Background(), request{Op: opClose})
	return err
}

// ListDatabases returns the account's database names.
func (c *tapeClient) ListDatabases(ctx context.Context) ([]string, error) {
	a, err := c.do(ctx, request{Op: opListDatabases})
	if err != nil {
		return nil, err
	}
	return a.Listing, nil
}

// ListContainers returns database's container names.
func (c *tapeClient) ListContainers(ctx context.Context, database string) ([]string, error) {
	a, err := c.do(ctx, request{Op: opListContainers, Database: database})
	if err != nil {
		return nil, err
	}
	return a.Listing, nil
}

// Plan returns the plan the gateway stated for a statement.
func (c *tapeClient) Plan(ctx context.Context, req PageRequest) ([]byte, error) {
	a, err := c.do(ctx, request{
		Op:        opPlan,
		Database:  req.Database,
		Container: req.Container,
		Statement: req.Statement,
	})
	if err != nil {
		return nil, err
	}
	return a.Plan, nil
}

// Page reads one container query, page by page.
func (c *tapeClient) Page(ctx context.Context, req PageRequest) (PageAnswer, error) {
	a, err := c.do(ctx, request{
		Op:        opPage,
		Database:  req.Database,
		Container: req.Container,
		Statement: req.Statement,
		PageSize:  req.PageSize,
		Engine:    req.Engine,
		Limit:     req.Limit,
	})
	// The answer's pages come back with the error: a query that failed after a
	// page must not be retried through the engine, in a replay as in a
	// recording.
	items := make([][]byte, len(a.Items))
	for i, item := range a.Items {
		items[i] = item
	}
	return PageAnswer{Items: items, Pages: a.Pages}, err
}

// do sends one call and returns its answer.
func (c *tapeClient) do(_ context.Context, r request) (answer, error) {
	text, err := r.text()
	if err != nil {
		return answer{}, err
	}
	answered, err := c.tape.Do(text)
	if err != nil {
		return answer{}, err
	}
	var a answer
	if err := unmarshalJSON([]byte(answered), &a); err != nil {
		return answer{}, fmt.Errorf("testcontainer: unreadable answer to %s: %w", text, err)
	}
	if a.Error != nil {
		return a, a.Error.err()
	}
	return a, nil
}

// The calls the driver makes of its client, as a tape names them.
const (
	opPing           = "Ping"
	opClose          = "Close"
	opListDatabases  = "ListDatabases"
	opListContainers = "ListContainers"
	opPage           = "Page"
	opPlan           = "Plan"
)

// request is one call as a tape holds it.
type request struct {
	Op        string `json:"op"`
	Database  string `json:"database,omitempty"`
	Container string `json:"container,omitempty"`
	Statement string `json:"statement,omitempty"`
	PageSize  int32  `json:"pageSize,omitempty"`
	Engine    bool   `json:"engine,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

func (r request) text() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// answer is one call's answer as a tape holds it.
type answer struct {
	Listing []string          `json:"listing,omitempty"`
	Items   []json.RawMessage `json:"items,omitempty"`
	Pages   int               `json:"pages,omitempty"`
	Plan    json.RawMessage   `json:"plan,omitempty"`
	Error   *clientError      `json:"error,omitempty"`
}

// clientError is an error the client returned, as a tape holds it. A gateway
// error keeps the response behind its text, so a replay's error carries the same
// type and the same text a recording's error had: the driver reads a
// cross-partition refusal off both, and a test reads the status off either.
type clientError struct {
	Response *responseError  `json:"response,omitempty"`
	Emulator *emulatorStatus `json:"emulator,omitempty"`
	Message  string          `json:"message,omitempty"`
}

// emulatorStatus is a failed emulator REST call as a tape holds it: the status
// and body its text is made of, so a replay's error is the same type the live
// client returns.
type emulatorStatus struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body,omitempty"`
}

// responseError is an *azcore.ResponseError as a tape holds it: the fields that
// make up its text and the status the driver reads.
type responseError struct {
	Method     string `json:"method,omitempty"`
	URL        string `json:"url,omitempty"`
	Status     string `json:"status,omitempty"`
	StatusCode int    `json:"statusCode"`
	ErrorCode  string `json:"errorCode,omitempty"`
	Body       string `json:"body,omitempty"`
}

func (e *clientError) err() error {
	if e.Emulator != nil {
		return &emulatorHTTPError{StatusCode: e.Emulator.StatusCode, Body: e.Emulator.Body}
	}
	if e.Response == nil {
		return errors.New(e.Message)
	}
	req, _ := http.NewRequest(e.Response.Method, e.Response.URL, nil)
	resp := &http.Response{
		StatusCode: e.Response.StatusCode,
		Status:     e.Response.Status,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(e.Response.Body)),
		Request:    req,
	}
	return azruntime.NewResponseErrorWithErrorCode(resp, e.Response.ErrorCode)
}

// encodeError renders err as an answer's error half. A gateway error keeps the
// response behind its text; any other error is its text alone.
func encodeError(err error) *clientError {
	var restErr *emulatorHTTPError
	if errors.As(err, &restErr) {
		return &clientError{Emulator: &emulatorStatus{
			StatusCode: restErr.StatusCode,
			Body:       activityID.ReplaceAllString(restErr.Body, "ActivityId: 00000000-0000-0000-0000-000000000000"),
		}}
	}
	var responseErr *azcore.ResponseError
	if !errors.As(err, &responseErr) || responseErr.RawResponse == nil {
		return &clientError{Message: err.Error()}
	}
	raw := responseErr.RawResponse
	body := ""
	if raw.Body != nil {
		if b, readErr := io.ReadAll(raw.Body); readErr == nil {
			body = activityID.ReplaceAllString(string(b), "ActivityId: 00000000-0000-0000-0000-000000000000")
			raw.Body = io.NopCloser(strings.NewReader(body))
		}
	}
	e := &responseError{
		Status:     raw.Status,
		StatusCode: responseErr.StatusCode,
		ErrorCode:  responseErr.ErrorCode,
		Body:       body,
	}
	if raw.Request != nil {
		e.Method = raw.Request.Method
		if raw.Request.URL != nil {
			e.URL = raw.Request.URL.String()
		}
	}
	return &clientError{Response: e}
}

// systemProperties are the properties the emulator adds to every document and
// that only the store that answered can name: a physical resource id, its self
// link, its etag, the write's timestamp and its attachment link. A tape cannot
// hold them, so an answer narrows every item to the fields the fixture wrote.
var systemProperties = []string{"_rid", "_self", "_etag", "_ts", "_attachments"}

// narrowItem drops the emulator's own properties from one item. A row that is
// not an object, as the VALUE form of a projection returns, carries none.
func narrowItem(item []byte) (json.RawMessage, error) {
	var document map[string]json.RawMessage
	if err := json.Unmarshal(item, &document); err != nil {
		return item, nil
	}
	for _, name := range systemProperties {
		delete(document, name)
	}
	narrowed, err := json.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("testcontainer: narrow an item: %w", err)
	}
	return narrowed, nil
}

// ordersResult matches the operators that decide a result's own order: an
// ORDER BY, and the GROUP BY whose rewritten query carries one.
var ordersResult = regexp.MustCompile(`(?i)\b(order|group)\s+by\b`)

// activityID is the id the emulator puts in every error it answers, which
// names one request and so differs between runs.
var activityID = regexp.MustCompile(`ActivityId: [0-9a-fA-F-]{36}`)

// unmarshalJSON decodes JSON strictly enough that a tape's typo fails the test
// rather than answering something else.
func unmarshalJSON(b []byte, v any) error {
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// live answers one request from the pinned emulator, through the client the
// driver itself opens.
func live(t testing.TB, emulator bool, text string) (string, error) {
	var req request
	if err := unmarshalJSON([]byte(text), &req); err != nil {
		return "", fmt.Errorf("testcontainer: unreadable request %s: %w", text, err)
	}
	endpoint, err := testcontainer.Endpoint(t)
	if err != nil {
		return "", err
	}
	cfg := driver.Config{
		Host:     endpoint,
		Database: testcontainer.Database,
		Password: testcontainer.EmulatorKey,
		TLS:      testcontainer.InsecureTLS,
		Params:   map[string]string{ParamAuth: AuthMasterKey},
	}
	if emulator {
		cfg.Params[ParamEmulator] = "true"
	}
	c, err := NewClient(cfg)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	a := answer{}
	switch req.Op {
	case opPing:
		err = c.Ping(ctx)
	case opClose:
		err = c.Close()
	case opListDatabases:
		a.Listing, err = c.ListDatabases(ctx)
	case opListContainers:
		a.Listing, err = c.ListContainers(ctx, req.Database)
	case opPlan:
		a.Plan, err = c.Plan(ctx, PageRequest{
			Database:  req.Database,
			Container: req.Container,
			Statement: req.Statement,
		})
	case opPage:
		var page PageAnswer
		page, err = c.Page(ctx, PageRequest{
			Database:  req.Database,
			Container: req.Container,
			Statement: req.Statement,
			PageSize:  req.PageSize,
			Engine:    req.Engine,
			Limit:     req.Limit,
		})
		for _, item := range page.Items {
			narrowed, narrowErr := narrowItem(item)
			if narrowErr != nil {
				return "", narrowErr
			}
			a.Items = append(a.Items, narrowed)
		}
		if !ordersResult.MatchString(req.Statement) {
			// The order of a result whose query names no ORDER BY is the server's
			// to choose: the client merges partition streams as their answers
			// arrive, so two runs of one query disagree. The tape holds such a
			// result in one canonical order, and a test that draws a conclusion
			// from the order names the operator that decides it.
			slices.SortFunc(a.Items, func(x, y json.RawMessage) int { return bytes.Compare(x, y) })
		}
		a.Pages = page.Pages
	default:
		return "", fmt.Errorf("testcontainer: unknown request %q", req.Op)
	}
	if err != nil {
		a.Error = encodeError(err)
	}
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
