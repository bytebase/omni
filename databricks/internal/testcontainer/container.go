// Package testcontainer is how databricks' driver tests reach a workspace:
// through a tape (enginetest/tape) that replays a pinned workspace's answers,
// or, in a recording (OMNI_RECORD=1), through the workspace itself. A replay
// starts nothing and needs no credentials.
//
// Databricks is an account service, so no container stands in for it. The
// tape sits at the Databricks SDK's HTTP transport: a request is its method,
// its route and its body, and an answer is its status, its content type and
// its body. The Authorization header is never recorded, and the host and
// token stay with the live client. A recording reads DATABRICKS_HOST and
// DATABRICKS_TOKEN from the environment; neither is written down.
//
// An answer is narrowed to what the driver reads. Catalogs, schemas and tables
// keep their names and the table fields the sync maps; statement answers keep
// their columns, rows and state. Listings are sorted, and every value the
// service makes up that survives narrowing is replaced by the tape's stand-in:
// the statement ID, which names a later poll. Without that, two recordings of
// one workspace would not write the same bytes.
package testcontainer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/databricks/databricks-sdk-go"
	sdkconfig "github.com/databricks/databricks-sdk-go/config"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is databricks/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Driver is the pin row of the workspace the driver's tapes were recorded on.
const Driver = "driver"

// WarehouseID and Catalog are the SQL warehouse and catalog the tapes were
// recorded against, and the ones the pinned workspace's pin row names. The
// tests construct a driver with them instead of reading a DSN, so a replay
// needs no environment.
const (
	WarehouseID = "45847ad2fbd0b316"
	Catalog     = "workspace"
)

// Server is the tape server the driver's recordings are written to.
var Server = &tape.Server{Pins: Pins, Codec: codec{}}

type codec struct{}

func (codec) EncodeError(error) (string, bool) { return "", false }

func (codec) DecodeError(text string) error { return fmt.Errorf("%s", text) } //nolint:revive // an answer's text, not a cause to wrap

// madeUp matches the statement IDs the service makes up, which name a later
// poll. The narrowed answers carry no other UUID.
var madeUp = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

// expressionID matches the number Spark gives an expression in a plan, which
// changes between runs even when the plan does not.
var expressionID = regexp.MustCompile(`#[0-9]+`)

func (codec) MadeUp() *regexp.Regexp { return madeUp }

func (codec) StandIn(_ string, n int) string {
	return fmt.Sprintf("00000000-0000-0000-0000-%012x", n)
}

// Open returns t's workspace client: a replay of its tape, or, in a recording,
// the pinned workspace with every HTTP exchange written down. The client's
// token is a placeholder because the tape never reaches the service; the live
// transport reads the real host and token from the environment.
func Open(t testing.TB) *databricks.WorkspaceClient {
	t.Helper()
	pin, err := enginetest.LoadPin(Pins, Driver)
	if err != nil {
		t.Fatalf("testcontainer: %v", err)
	}
	client, err := databricks.NewWorkspaceClient(&databricks.Config{
		Host:  "https://" + pin.Image,
		Token: "omni-tape",
		HTTPTransport: &transport{tape: tape.OpenRaw(t, Server, Driver, func(request string) (string, error) {
			return live(t, pin, request)
		})},
		HostMetadataResolver: func(context.Context, string) (*sdkconfig.HostMetadata, error) {
			return nil, nil
		},
	})
	if err != nil {
		t.Fatalf("testcontainer: create workspace client: %v", err)
	}
	return client
}

// request is one exchange as a tape holds it. The route is the path and query
// the SDK wrote, without the workspace host or the request's credentials.
type request struct {
	Method string `json:"method"`
	Route  string `json:"route"`
	Body   string `json:"body,omitempty"`
}

func (r request) text() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// answer is one exchange's answer as a tape holds it: what the service sent,
// with nothing of the connection it sent it over.
type answer struct {
	Status      int    `json:"status"`
	ContentType string `json:"contentType,omitempty"`
	Body        string `json:"body,omitempty"`
}

// transport answers the SDK's HTTP exchange from the tape. A recording sends
// the request to the workspace and narrows the answer before the tape holds
// it.
type transport struct {
	tape *tape.Raw
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	var body []byte
	if r.Body != nil {
		var err error
		body, err = io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		_ = r.Body.Close()
	}
	text, err := request{Method: r.Method, Route: r.URL.RequestURI(), Body: string(body)}.text()
	if err != nil {
		return nil, err
	}
	answered, err := t.tape.Do(text)
	if err != nil {
		return nil, err
	}
	var a answer
	if err := json.Unmarshal([]byte(answered), &a); err != nil {
		return nil, fmt.Errorf("testcontainer: unreadable answer to %s: %w", text, err)
	}
	header := http.Header{}
	if a.ContentType != "" {
		header.Set("Content-Type", a.ContentType)
	}
	return &http.Response{
		StatusCode: a.Status,
		Status:     fmt.Sprintf("%d %s", a.Status, http.StatusText(a.Status)),
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(a.Body)),
		Request:    r,
	}, nil
}

// live answers one request from the pinned workspace. It reads the host and
// token from the environment and refuses a host other than the pin's, so a
// recording cannot write answers under the wrong pin.
func live(t testing.TB, pin enginetest.Pin, text string) (string, error) {
	t.Helper()
	var req request
	if err := json.Unmarshal([]byte(text), &req); err != nil {
		return "", fmt.Errorf("testcontainer: unreadable request %q: %w", text, err)
	}
	host := strings.TrimRight(strings.TrimSpace(os.Getenv("DATABRICKS_HOST")), "/")
	host = strings.TrimPrefix(host, "https://")
	token := strings.TrimSpace(os.Getenv("DATABRICKS_TOKEN"))
	if host == "" || token == "" {
		t.Fatalf("testcontainer: recording needs DATABRICKS_HOST and DATABRICKS_TOKEN")
	}
	if host != pin.Image {
		t.Fatalf("testcontainer: DATABRICKS_HOST names %q, but pins.tsv pins %q", host, pin.Image)
	}
	httpReq, err := http.NewRequest(req.Method, "https://"+host+req.Route, strings.NewReader(req.Body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	contentType := resp.Header.Get("Content-Type")
	body, contentType = narrow(req.Route, contentType, body, resp.StatusCode)
	b, err := json.Marshal(answer{Status: resp.StatusCode, ContentType: contentType, Body: string(body)})
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// narrow keeps what the driver reads of an answer, in one order. A service
// answer that fails is kept to its status's error code and message; the rest
// of the body names the request that failed.
func narrow(route, contentType string, body []byte, status int) ([]byte, string) {
	path, _, _ := strings.Cut(route, "?")
	if !strings.Contains(contentType, "application/json") {
		return body, contentType
	}
	if status < 200 || status >= 300 {
		var failure struct {
			ErrorCode string `json:"error_code,omitempty"`
			Message   string `json:"message,omitempty"`
		}
		if err := json.Unmarshal(body, &failure); err != nil {
			return body, contentType
		}
		narrowed, err := json.Marshal(failure)
		if err != nil {
			return body, contentType
		}
		return narrowed, "application/json"
	}

	switch {
	case path == "/api/2.0/sql/statements", strings.HasPrefix(path, "/api/2.0/sql/statements/"):
		narrowed, ok := narrowStatement(body)
		if ok {
			return narrowed, "application/json"
		}
	case path == "/api/2.1/unity-catalog/catalogs":
		narrowed, ok := narrowCatalogs(body)
		if ok {
			return narrowed, "application/json"
		}
	case path == "/api/2.1/unity-catalog/schemas":
		narrowed, ok := narrowSchemas(body)
		if ok {
			return narrowed, "application/json"
		}
	case path == "/api/2.1/unity-catalog/tables":
		narrowed, ok := narrowTables(body)
		if ok {
			return narrowed, "application/json"
		}
	case path == "/api/2.0/workspace/list":
		narrowed, ok := narrowWorkspaceList(body)
		if ok {
			return narrowed, "application/json"
		}
	}
	return body, contentType
}

// statementAnswer is a statement answer as a tape holds it.
type statementAnswer struct {
	StatementID string           `json:"statement_id,omitempty"`
	Status      *statementStatus `json:"status,omitempty"`
	Manifest    *manifest        `json:"manifest,omitempty"`
	Result      *result          `json:"result,omitempty"`
}

type statementStatus struct {
	State    string          `json:"state,omitempty"`
	Error    *statementError `json:"error,omitempty"`
	SQLState string          `json:"sql_state,omitempty"`
}

type statementError struct {
	ErrorCode string `json:"error_code,omitempty"`
	Message   string `json:"message,omitempty"`
}

type manifest struct {
	Format    string  `json:"format,omitempty"`
	Schema    *schema `json:"schema,omitempty"`
	Truncated bool    `json:"truncated"`
}

type schema struct {
	ColumnCount int      `json:"column_count,omitempty"`
	Columns     []column `json:"columns,omitempty"`
}

type column struct {
	Name     string `json:"name,omitempty"`
	TypeText string `json:"type_text,omitempty"`
	TypeName string `json:"type_name,omitempty"`
	Position int    `json:"position,omitempty"`
	Nullable bool   `json:"nullable"`
	Comment  string `json:"comment,omitempty"`
}

type result struct {
	DataArray [][]string `json:"data_array,omitempty"`
}

func narrowStatement(body []byte) ([]byte, bool) {
	var response statementAnswer
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, false
	}
	// A plan names the expression IDs Spark assigned, which move between runs.
	// The plan's shape is the answer; the numbers are not.
	if response.Manifest != nil && response.Manifest.Schema != nil && len(response.Manifest.Schema.Columns) == 1 && response.Manifest.Schema.Columns[0].Name == "plan" && response.Result != nil {
		for i, row := range response.Result.DataArray {
			for j, cell := range row {
				response.Result.DataArray[i][j] = expressionID.ReplaceAllString(cell, "#0")
			}
		}
	}
	narrowed, err := json.Marshal(response)
	return narrowed, err == nil
}

func narrowCatalogs(body []byte) ([]byte, bool) {
	var response struct {
		Catalogs []struct {
			Name string `json:"name"`
		} `json:"catalogs"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, false
	}
	sort.Slice(response.Catalogs, func(i, j int) bool { return response.Catalogs[i].Name < response.Catalogs[j].Name })
	narrowed, err := json.Marshal(response)
	return narrowed, err == nil
}

func narrowSchemas(body []byte) ([]byte, bool) {
	var response struct {
		Schemas []struct {
			Name string `json:"name"`
		} `json:"schemas"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, false
	}
	sort.Slice(response.Schemas, func(i, j int) bool { return response.Schemas[i].Name < response.Schemas[j].Name })
	narrowed, err := json.Marshal(response)
	return narrowed, err == nil
}

func narrowTables(body []byte) ([]byte, bool) {
	var response struct {
		Tables []table `json:"tables"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, false
	}
	for i := range response.Tables {
		table := &response.Tables[i]
		sort.SliceStable(table.Columns, func(i, j int) bool {
			if table.Columns[i].Position != table.Columns[j].Position {
				return table.Columns[i].Position < table.Columns[j].Position
			}
			return table.Columns[i].Name < table.Columns[j].Name
		})
	}
	sort.Slice(response.Tables, func(i, j int) bool { return response.Tables[i].Name < response.Tables[j].Name })
	narrowed, err := json.Marshal(response)
	return narrowed, err == nil
}

type table struct {
	Name           string   `json:"name"`
	CatalogName    string   `json:"catalog_name,omitempty"`
	SchemaName     string   `json:"schema_name,omitempty"`
	TableType      string   `json:"table_type,omitempty"`
	ViewDefinition string   `json:"view_definition,omitempty"`
	Comment        string   `json:"comment,omitempty"`
	Columns        []column `json:"columns,omitempty"`
}

func narrowWorkspaceList(body []byte) ([]byte, bool) {
	var response struct {
		Objects []struct {
			Path       string `json:"path"`
			ObjectType string `json:"object_type,omitempty"`
		} `json:"objects"`
		NextPageToken string `json:"next_page_token,omitempty"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, false
	}
	sort.Slice(response.Objects, func(i, j int) bool { return response.Objects[i].Path < response.Objects[j].Path })
	narrowed, err := json.Marshal(response)
	return narrowed, err == nil
}
