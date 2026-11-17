// Package testcontainer is how elasticsearch's tests reach Elasticsearch:
// through a tape (enginetest/tape) that replays a pinned server's answers, or,
// in a recording (OMNI_RECORD=1), through the pinned server itself. The server
// is one testcontainer per test binary, started by the first test that records
// and removed by Terminate, which each package's TestMain calls.
//
// Elasticsearch's client is HTTP, so a tape would hold the port the container
// happened to publish. Client asks for the part that does not move instead: a
// request is its method, its route (the path and query the driver wrote) and
// its body, and an answer is the status, the content type and the body the
// endpoint returned. The endpoint's address and login stay with the live
// client, which is why a tape holds neither.
//
// An answer is narrowed to what the driver asks of it, because some carry
// values the metadata and its golden must not: the root document's cluster
// UUID, the cat API's per-index UUID, a search's `took`, and the timings inside
// a profile. A search's numbers are read back as json.Number, so a long past
// 2^53 keeps its digits.
package testcontainer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is elasticsearch/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Driver is the pin row of the server the driver's tapes were recorded on.
const Driver = "driver"

const (
	user     = "elastic"
	password = "omni-test-password"
)

var server = &tape.Server{Pins: Pins, Codec: codec{}}

// codec is how an Elasticsearch answer is read for the tape. An HTTP error is
// part of its answer's status, so the tape needs no error spelling from the
// codec.
type codec struct{}

func (codec) EncodeError(error) (string, bool) { return "", false }

func (codec) DecodeError(text string) error { return errors.New(text) }

// Client is the driver's HTTP client, answered by t's tape. It is what a test
// hands the driver in place of the live endpoint; its methods are the calls the
// driver makes, and the recordings hold what the endpoint answered.
type Client struct {
	tape *tape.Raw
}

// Open returns t's client to the pinned server: a replay of its tape, or, in a
// recording, the pinned server with every call written down.
func Open(t testing.TB) *Client {
	t.Helper()
	return &Client{tape: tape.OpenRaw(t, server, Driver, func(request string) (string, error) {
		// A recording reaches the server on its first call, so a test that
		// replays never starts a container.
		return liveDo(t, request)
	})}
}

// Do sends one request.
func (c *Client) Do(_ context.Context, method, route string, body []byte) (*http.Response, error) {
	text, err := request{Method: method, Route: route, Body: string(body)}.text()
	if err != nil {
		return nil, err
	}
	answered, err := c.tape.Do(text)
	if err != nil {
		return nil, err
	}
	var a answer
	if err := unmarshalJSON([]byte(answered), &a); err != nil {
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
	}, nil
}

// request is one call as a tape holds it.
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

// answer is one call's answer as a tape holds it: what the endpoint sent, with
// nothing of the connection it sent it over.
type answer struct {
	Status      int    `json:"status"`
	ContentType string `json:"contentType,omitempty"`
	Body        string `json:"body,omitempty"`
}

// unmarshalJSON decodes JSON strictly enough that a tape's typo fails the test
// rather than answering something else.
func unmarshalJSON(b []byte, v any) error {
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

// Reset deletes every index the corpus of a previous test created, so a tape
// holds one test's corpus whatever the test binary ran before. It does nothing
// in a replay, whose tape already holds the answers the corpus produced when it
// was recorded.
func Reset(t testing.TB) {
	t.Helper()
	if !enginetest.Recording(t) {
		return
	}
	address, client := live(t)
	ctx := context.Background()
	// Every index, hidden ones and data streams included; the system indices
	// the endpoint keeps for itself start with a dot and stay.
	body, status, _, err := send(ctx, client, address, http.MethodGet, "/_cat/indices?format=json&expand_wildcards=all&h=index", nil)
	if err != nil || status >= 400 {
		t.Fatalf("list indices: %v (%d %s)", err, status, body)
	}
	var entries []struct {
		Index string `json:"index"`
	}
	if err := json.Unmarshal(body, &entries); err != nil {
		t.Fatalf("list indices: %v", err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Index, ".") {
			continue
		}
		body, status, _, err := send(ctx, client, address, http.MethodDelete, "/"+entry.Index, nil)
		if err != nil || status >= 400 {
			t.Fatalf("delete index %s: %v (%d %s)", entry.Index, err, status, body)
		}
	}
}

// --- the live server ---

const (
	// startupTimeout is how long Elasticsearch's first boot may take: a JVM
	// starting, its security index created and its node elected.
	startupTimeout = 5 * time.Minute
)

var liveServer struct {
	once      sync.Once
	container testcontainers.Container
	address   string
	client    *http.Client
	err       error
}

// live returns the pinned server's address and client, starting the container
// on the first call. Only a recording calls it.
func live(t testing.TB) (string, *http.Client) {
	t.Helper()
	liveServer.once.Do(func() {
		liveServer.container, liveServer.address, liveServer.client, liveServer.err = start(context.Background())
	})
	if liveServer.err != nil {
		t.Fatalf("start Elasticsearch: %v", liveServer.err)
	}
	return liveServer.address, liveServer.client
}

func start(ctx context.Context) (started testcontainers.Container, address string, client *http.Client, err error) {
	pin, err := enginetest.LoadPin(Pins, Driver)
	if err != nil {
		return nil, "", nil, err
	}
	// testcontainers-go panics from MustExtractDockerHost when the Docker
	// socket is reachable but access is denied, and a test must report that
	// rather than crash the binary.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	container, err := testcontainers.Run(ctx, pin.ImageRef(),
		testcontainers.WithEnv(map[string]string{
			"discovery.type": "single-node",
			// The driver's login path is basic auth, so the endpoint keeps
			// its security index and asks for a password.
			"xpack.security.enabled": "true",
			"ELASTIC_PASSWORD":       password,
			// The container is served over plain HTTP, which keeps the
			// tape's answers out of a certificate.
			"xpack.security.http.ssl.enabled": "false",
			// The machine is shared, and its disk is fuller than
			// Elasticsearch's high watermark allows: a shard the node
			// refuses to allocate would leave the index at zero documents
			// and no size.
			"cluster.routing.allocation.disk.threshold_enabled": "false",
			"ES_JAVA_OPTS": "-Xms512m -Xmx512m",
		}),
		testcontainers.WithExposedPorts("9200/tcp"),
		testcontainers.WithWaitStrategy(
			// With security on, the root document answers 401 once the node
			// is up.
			wait.ForHTTP("/").WithPort("9200/tcp").
				WithStatusCodeMatcher(func(status int) bool { return status == http.StatusOK || status == http.StatusUnauthorized }).
				WithStartupTimeout(startupTimeout)),
	)
	if err != nil {
		return nil, "", nil, err
	}
	if err := checkDigest(ctx, container, pin.Digest); err != nil {
		return nil, "", nil, err
	}
	port, err := container.MappedPort(ctx, "9200/tcp")
	if err != nil {
		return nil, "", nil, err
	}
	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", nil, err
	}
	address = "http://" + host + ":" + port.Port()
	client = &http.Client{Timeout: 60 * time.Second}

	body, status, _, err := send(ctx, client, address, http.MethodGet, "/", nil)
	if err != nil {
		return nil, "", nil, fmt.Errorf("Elasticsearch does not answer: %w", err)
	}
	if status != http.StatusOK {
		return nil, "", nil, fmt.Errorf("Elasticsearch does not answer: %d %s", status, body)
	}
	var root struct {
		Version struct {
			Number string `json:"number"`
		} `json:"version"`
	}
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, "", nil, fmt.Errorf("Elasticsearch does not answer: %w", err)
	}
	if root.Version.Number != pin.Tag {
		return nil, "", nil, fmt.Errorf("the server reports %q, but pins.tsv pins Elasticsearch %s", root.Version.Number, pin.Tag)
	}
	return container, address, client, nil
}

// checkDigest proves the container runs the pinned image. The pull names
// image@digest, so a moved tag cannot select another image; this check reads
// what the daemon actually started, which is where the recorded answers came
// from.
func checkDigest(ctx context.Context, container testcontainers.Container, digest string) error {
	info, err := container.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the Elasticsearch container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	image, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the Elasticsearch image: %w", err)
	}
	for _, d := range image.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the Elasticsearch container runs an image with digests %v, but pins.tsv pins %s", image.RepoDigests, digest)
}

// Terminate removes the server this run started, so a recording leaves no
// container behind even when the reaper is disabled.
func Terminate() {
	if liveServer.container != nil {
		_ = testcontainers.TerminateContainer(liveServer.container)
	}
}

// send sends one request to address as the pinned server's user, and returns
// what the endpoint answered.
func send(ctx context.Context, client *http.Client, address, method, route string, body []byte) ([]byte, int, string, error) {
	req, err := http.NewRequestWithContext(ctx, method, address+route, bytes.NewReader(body))
	if err != nil {
		return nil, 0, "", err
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, "", err
	}
	defer resp.Body.Close()
	answered, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, "", err
	}
	return answered, resp.StatusCode, resp.Header.Get("Content-Type"), nil
}

// --- the live client's answers ---

// liveDo answers one request from the pinned server, in the encoding a tape
// holds.
func liveDo(t testing.TB, text string) (string, error) {
	t.Helper()
	var req request
	if err := unmarshalJSON([]byte(text), &req); err != nil {
		return "", fmt.Errorf("testcontainer: unreadable request %q: %w", text, err)
	}
	address, client := live(t)
	body, status, contentType, err := send(context.Background(), client, address, req.Method, req.Route, []byte(req.Body))
	if err != nil {
		return "", err
	}
	if narrowed, ok := narrow(req.Route, contentType, body); ok {
		body = narrowed
	}
	b, err := json.Marshal(answer{Status: status, ContentType: contentType, Body: string(body)})
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// narrow returns the answer as a tape holds it. An answer that names the
// endpoint itself, one index's settings, the cat API's index list or the
// instance's users keeps only the fields the driver reads; anything else is
// held as sent, with the measurements a search reports taken out.
func narrow(route, contentType string, body []byte) ([]byte, bool) {
	if !strings.Contains(contentType, "application/json") {
		return nil, false
	}
	switch {
	case route == "/":
		var root versionResult
		if err := json.Unmarshal(body, &root); err != nil {
			return nil, false
		}
		narrowed, err := json.Marshal(root)
		return narrowed, err == nil
	case strings.HasPrefix(route, "/*/_settings"):
		var settings map[string]indexSettingsResult
		if err := json.Unmarshal(body, &settings); err != nil {
			return nil, false
		}
		kept := map[string]indexSettingsResult{}
		for name, value := range settings {
			if value.Settings.Index.Hidden == "true" {
				kept[name] = value
			}
		}
		narrowed, err := json.Marshal(kept)
		return narrowed, err == nil
	case strings.HasPrefix(route, "/_cat/indices"):
		var entries []indicesResult
		if err := json.Unmarshal(body, &entries); err != nil {
			return nil, false
		}
		sort.SliceStable(entries, func(i, j int) bool { return entries[i].Index < entries[j].Index })
		narrowed, err := json.Marshal(entries)
		return narrowed, err == nil
	case isBulkEndpoint(route):
		return narrowBulk(body)
	case strings.HasPrefix(route, "/_security/user"), strings.HasPrefix(route, "/_plugins/_security/api/internalusers"):
		var users map[string]struct {
			Roles []string `json:"roles"`
		}
		if err := json.Unmarshal(body, &users); err != nil {
			return nil, false
		}
		narrowed, err := json.Marshal(users)
		return narrowed, err == nil
	}
	// A search, or any other JSON answer the driver hands back: the numbers
	// that measure the run are not part of what it states.
	var document any
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err := decoder.Decode(&document); err != nil {
		return nil, false
	}
	narrowed, err := json.Marshal(withoutTimings(document))
	return narrowed, err == nil
}

// isBulkEndpoint reports whether a route names the bulk endpoint, "/_bulk" or
// "/<index>/_bulk". A document whose id or index carries the word is not one:
// its answer is the document's, and narrowing it would fabricate a verdict.
func isBulkEndpoint(route string) bool {
	path, _, _ := strings.Cut(route, "?")
	path = strings.TrimRight(path, "/")
	return path == "_bulk" || strings.HasSuffix(path, "/_bulk")
}

// narrowBulk keeps what the driver reads of a bulk answer: whether any item
// failed and, per item, its status and its error's type and reason. The rest of
// an item names the index, its UUID and the shard, which a container makes up
// afresh.
func narrowBulk(body []byte) ([]byte, bool) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	var answer struct {
		Errors bool                        `json:"errors"`
		Items  []map[string]map[string]any `json:"items"`
	}
	if err := decoder.Decode(&answer); err != nil {
		return nil, false
	}
	items := make([]any, 0, len(answer.Items))
	for _, item := range answer.Items {
		for op, outcome := range item {
			kept := map[string]any{}
			if status, ok := outcome["status"]; ok {
				kept["status"] = status
			}
			if failure, ok := outcome["error"].(map[string]any); ok {
				reason := map[string]any{}
				for _, key := range []string{"type", "reason"} {
					if value, ok := failure[key]; ok {
						reason[key] = value
					}
				}
				kept["error"] = reason
			}
			items = append(items, map[string]any{op: kept})
		}
	}
	narrowed, err := json.Marshal(map[string]any{"errors": answer.Errors, "items": items})
	return narrowed, err == nil
}

// versionResult is the root document as a tape holds it.
type versionResult struct {
	Version struct {
		Number string `json:"number"`
	} `json:"version"`
}

// indexSettingsResult is one index's hidden flag.
type indexSettingsResult struct {
	Settings struct {
		Index struct {
			Hidden string `json:"hidden"`
		} `json:"index"`
	} `json:"settings"`
}

// indicesResult is one cat entry as a tape holds it.
type indicesResult struct {
	Index     string `json:"index"`
	DocsCount string `json:"docs.count"`
	IndexSize string `json:"store.size"`
}

// timedFields are the measurements a search and its profile report. Which
// nanosecond a step took is a fact about the run, not about the search, and a
// tape of a replayed search must answer the same bytes every time.
var timedFields = map[string]bool{
	"took":          true,
	"rewrite_time":  true,
	"time_in_nanos": true,
}

// withoutTimings returns an answer with the measurements it reports taken out:
// the top-level `took`, and the tree under a top-level `profile`. Nothing else
// is walked, so a document a hit carries keeps every field it holds, even one
// named `took`, `profile` or `time_in_nanos`.
func withoutTimings(v any) any {
	document, ok := v.(map[string]any)
	if !ok {
		return v
	}
	if _, ok := document["took"]; ok {
		document["took"] = json.Number("0")
	}
	if profile, ok := document["profile"]; ok {
		document["profile"] = withoutProfileTimings(profile)
	}
	return document
}

// withoutProfileTimings returns a profile with its measurements zeroed, the
// numbers under each "breakdown" with them, and the node and shard it names
// held as "0", since a container makes those up afresh.
func withoutProfileTimings(v any) any {
	switch value := v.(type) {
	case map[string]any:
		for key, item := range value {
			switch {
			case timedFields[key]:
				value[key] = json.Number("0")
			case key == "node_id" || key == "id":
				value[key] = "0"
			case key == "breakdown":
				value[key] = zeroNumbers(item)
			default:
				value[key] = withoutProfileTimings(item)
			}
		}
		return value
	case []any:
		for i := range value {
			value[i] = withoutProfileTimings(value[i])
		}
		return value
	}
	return v
}

// zeroNumbers returns v with every number under it zeroed.
func zeroNumbers(v any) any {
	switch value := v.(type) {
	case map[string]any:
		for key, item := range value {
			if _, ok := item.(json.Number); ok {
				value[key] = json.Number("0")
				continue
			}
			value[key] = zeroNumbers(item)
		}
		return value
	case []any:
		for i := range value {
			value[i] = zeroNumbers(value[i])
		}
		return value
	}
	return v
}
