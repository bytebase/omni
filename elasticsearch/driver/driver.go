// Package driver is Elasticsearch's implementation of the driver contract in
// github.com/bytebase/omni/driver: it speaks the REST API over HTTP,
// syncs an instance's version, its indices with their document counts and
// sizes, and its users, and runs Kibana Dev Console requests as interactive
// queries and change scripts.
//
// It implements Driver, Syncer, Querier and Executor. There is no Dumper: omni
// has no elasticsearch/catalog for a dump to render, and Bytebase's Dump wrote
// nothing. A synced index does carry a row count and a data size, because
// Elasticsearch reports them (_cat/indices' docs.count and store.size).
//
// Config.Host is the endpoint's host, optionally with its own scheme
// (https://search.example.com) and port; Config.Port applies when the host
// carries neither and defaults to 9200. Config.User and Config.Password are the
// basic-auth login. Config.TLS makes the endpoint https and answers for its
// host. Config.Dial is not used: Elasticsearch's client is HTTP, and an SSH
// tunnel or a cloud connector belongs in Options.Transport, which also carries
// a managed OpenSearch domain's AWS SigV4 signer.
//
// Config.Database does nothing: an Elasticsearch cluster is one namespace, and
// SyncDatabase reads it whole. Config.Params, Config.ApplicationName,
// Config.ReadOnly and Config.EngineVersion do nothing either, and neither do
// QueryOptions' Schema, Limit, ExplainFormat and MaskingEnabled: Elasticsearch
// has no schemas, no statement to rewrite with a row limit, one plan format,
// and nothing to withhold from a row. QueryOptions.Explain reads a search's
// profile, and only a search can be explained. QueryOptions.MaxResultBytes
// stops nothing: Bytebase read every row it was given, and the driver does too.
// ExecuteOptions' transaction settings, CreateDatabase and MaxRetries do
// nothing: every REST request stands alone.
//
// A search answered by a JSON object becomes one column per key, the keys in
// name order so one query shows its columns the same way every run. A column
// carries no type name, because Elasticsearch's answer states none. A number in
// an answer keeps its digits, so a long past 2^53 is not rounded.
//
// A statement's error is the driver's own. A request the server refused is a
// *ResponseError, whose Status is the HTTP status and whose Body is the answer
// as Elasticsearch wrote it: a failure is a JSON body there, and that body names
// the error's own type, reason and subject (an "index_not_found_exception" and
// the index it is about). A failure of the transport underneath is wrapped, so
// errors.As reaches the *url.Error net/http raises for a request that got no
// answer.
//
// Open builds the HTTP client and reaches nothing; Ping does. Every method is
// safe for concurrent use.
package driver

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/elasticsearch"
)

var (
	_ driver.Driver   = (*Driver)(nil)
	_ driver.Syncer   = (*Driver)(nil)
	_ driver.Querier  = (*Driver)(nil)
	_ driver.Executor = (*Driver)(nil)
)

// Options are Elasticsearch's own settings for Open.
type Options struct {
	// OpenSearch selects OpenSearch's security API for the instance's users
	// (/ _plugins/_security/api/internalusers). The rest of the protocol is
	// Elasticsearch's.
	OpenSearch bool
	// Transport sends every request in place of the default round tripper. A
	// managed OpenSearch domain signs each request with AWS SigV4, and the
	// signer needs the caller's credentials and region, so the caller builds
	// the transport. A transport means https, which such an endpoint requires.
	Transport http.RoundTripper
}

// Driver is one open Elasticsearch endpoint.
type Driver struct {
	// client is the HTTP conversation: the endpoint, the login and what
	// answers. A test replaces it with a tape.
	client client
	// openSearch selects the security API's path for the instance's users.
	openSearch bool
	// transport is the one Open built, nil when the caller supplied one, so
	// Close releases the idle connections the driver opened and none of the
	// caller's.
	transport *http.Transport
}

// Open returns a driver for the endpoint cfg names. It does not reach the
// endpoint; Ping does.
func Open(_ context.Context, cfg driver.Config, opts ...Options) (*Driver, error) {
	var opt Options
	if len(opts) > 0 {
		opt = opts[0]
	}
	secure := cfg.TLS != nil || opt.Transport != nil
	base, err := endpoint(cfg.Host, cfg.Port, secure)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Transport: opt.Transport}
	var transport *http.Transport
	if opt.Transport == nil {
		tlsConfig, err := tlsSettings(cfg, base.Hostname())
		if err != nil {
			return nil, err
		}
		transport = &http.Transport{
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			TLSClientConfig:       tlsConfig,
		}
		httpClient.Transport = transport
	}

	auth := ""
	if cfg.User != "" || cfg.Password != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(cfg.User + ":" + cfg.Password))
		auth = "Basic " + encoded
	}
	return &Driver{
		client:     &liveClient{base: base, http: httpClient, auth: auth},
		openSearch: opt.OpenSearch,
		transport:  transport,
	}, nil
}

// endpoint is the base URL the driver sends its requests to. A host may carry
// its own scheme and port; a scheme that contradicts Config.TLS is refused
// rather than silently used, since a caller that supplies TLS expects the
// credentials and the answers to travel inside it.
func endpoint(host, port string, secure bool) (*url.URL, error) {
	host = strings.TrimSpace(host)
	if host == "" {
		return nil, errors.New("host must be set")
	}
	if strings.Contains(host, "://") {
		base, err := url.Parse(host)
		if err != nil {
			return nil, fmt.Errorf("failed to parse address: %v: %w", host, err)
		}
		if secure && base.Scheme != "https" {
			return nil, fmt.Errorf("address %q is %s, but TLS is required", host, base.Scheme)
		}
		return base, nil
	}
	scheme := "http"
	if secure {
		scheme = "https"
	}
	// A host may carry its own port, which then wins. Brackets come off a
	// literal so JoinHostPort can put them back.
	if ownHost, ownPort, err := net.SplitHostPort(host); err == nil {
		host, port = ownHost, ownPort
	} else {
		host = strings.TrimSuffix(strings.TrimPrefix(host, "["), "]")
	}
	if port == "" {
		port = "9200"
	}
	return url.Parse(scheme + "://" + net.JoinHostPort(host, port))
}

// tlsSettings returns the TLS settings the endpoint's host is verified with.
func tlsSettings(cfg driver.Config, host string) (*tls.Config, error) {
	if cfg.TLS == nil {
		return nil, nil
	}
	tlsConfig, err := cfg.TLS(host)
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS config: %w", err)
	}
	if tlsConfig == nil {
		return nil, errors.New("TLS configuration must not be nil")
	}
	tlsConfig = tlsConfig.Clone()
	if tlsConfig.MinVersion == 0 {
		tlsConfig.MinVersion = tls.VersionTLS12
	}
	return tlsConfig, nil
}

// client is what the driver asks of the endpoint: one HTTP request answered.
// The route is the path and query the caller wrote, without the endpoint's
// address, so a tape holds no port that changes from one container to the next.
type client interface {
	Do(ctx context.Context, method, route string, body []byte) (*http.Response, error)
}

// liveClient is the client the driver sends over HTTP.
type liveClient struct {
	base *url.URL
	http *http.Client
	auth string
}

// Do sends one request. The route is the caller's text, appended to the
// endpoint as written: JoinPath would run path.Join over it, and a route such as
// /alpha/_doc/.. would leave here as /alpha — a document request becoming an
// index deletion. Its query string is kept as written too, since a '?' in it is
// the endpoint's, not an escape.
func (c *liveClient) Do(ctx context.Context, method, route string, body []byte) (*http.Response, error) {
	path, query, hasQuery := strings.Cut(route, "?")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	target := strings.TrimSuffix(c.base.String(), "/") + path
	if hasQuery {
		target += "?" + query
	}
	req, err := http.NewRequestWithContext(ctx, method, target, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to init a HTTP request: %w", err)
	}
	// Only add basic auth header if it exists (not for AWS IAM auth, whose
	// transport signs the request instead).
	if c.auth != "" {
		req.Header.Add("Authorization", c.auth)
	}
	req.Header.Set("Content-Type", "application/json")
	return c.http.Do(req)
}

// Ping asks the endpoint for its root document, which states its version.
func (d *Driver) Ping(ctx context.Context) error {
	res, err := d.client.Do(ctx, http.MethodGet, "/", nil)
	if err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("ping failed: [%d %s] %s", res.StatusCode, http.StatusText(res.StatusCode), body)
	}
	return nil
}

// Close releases the idle connections of the transport Open built. Elasticsearch
// keeps no session, and a caller's transport is the caller's to close.
func (d *Driver) Close() error {
	if d.transport != nil {
		d.transport.CloseIdleConnections()
	}
	return nil
}

// Session returns one session. Elasticsearch has no per-session state this
// driver sets, so a session is the driver itself.
func (d *Driver) Session(context.Context) (driver.Session, error) {
	return &session{driver: d}, nil
}

type session struct {
	driver *Driver
}

func (*session) Close() error { return nil }

// Query runs each request of script in order and returns one Result per
// request. A request that fails ends the script with its error on its own
// Result; the error returned is nil, since the script ran as far as it could.
// Only text the REST console parser cannot read at all is an error for the
// call, since then there is no request to send.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	parseResult, err := elasticsearch.ParseElasticsearchREST(script)
	if err != nil {
		return nil, err
	}
	if len(parseResult.Errors) > 0 {
		return nil, errors.New(parseResult.Errors[0].Message)
	}
	if parseResult.Requests == nil {
		return nil, nil
	}

	var results []*driver.Result
	for _, request := range parseResult.Requests {
		if opts.Explain {
			// The plan is what the search did, so the search runs with
			// profiling on; anything that is not a search is refused rather
			// than run.
			profiled, err := profiledSearch(request)
			if err != nil {
				results = append(results, &driver.Result{Statement: requestText(request), Err: err})
				break
			}
			request = profiled
		}
		result, err := s.run(ctx, request, opts)
		if err != nil {
			result = &driver.Result{Statement: requestText(request), Err: err}
		}
		results = append(results, result)
		if result.Err != nil {
			break
		}
	}
	return results, nil
}

// ResponseError is a request the server refused: the status it answered with and
// the body it wrote, which states the error's type and reason.
type ResponseError struct {
	Status int
	Body   string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("request failed with status code %d: %s", e.Status, e.Body)
}

// run sends one request and reads its answer.
func (s *session) run(ctx context.Context, request *elasticsearch.Request, opts driver.QueryOptions) (*driver.Result, error) {
	// For NDJSON endpoints (_msearch, _bulk, etc.), each JSON object must be
	// on a single line. Embedded newlines leave each data item before they are
	// joined with newline separators.
	var data []byte
	for _, item := range request.Data {
		data = append(data, []byte(strings.ReplaceAll(item, "\n", ""))...)
		data = append(data, '\n')
	}

	startTime := time.Now()
	resp, err := s.driver.client.Do(ctx, request.Method, request.URL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, &ResponseError{Status: resp.StatusCode, Body: string(respBytes)}
	}

	result := &driver.Result{Statement: requestText(request, data)}
	if request.Method == http.MethodHead || len(respBytes) == 0 {
		// A HEAD answer, and any other answer with no body, states its verdict
		// in the status alone: a successful one is a result with no columns and
		// no rows, and decode nothing.
		result.Latency = time.Since(startTime)
		return result, nil
	}
	var (
		columns []driver.Column
		row     []driver.Value
	)
	switch contentType := resp.Header.Get("Content-Type"); {
	case strings.Contains(contentType, "application/json"):
		// Elasticsearch's REST APIs answer with two shapes: an object, which
		// most APIs return, and an array, which the _cat APIs with format=json
		// return. Neither is a primitive type.
		var document any
		decoder := json.NewDecoder(strings.NewReader(string(respBytes)))
		// Numbers keep their digits: a long past 2^53 read as a float64 would
		// come back rounded.
		decoder.UseNumber()
		if err := decoder.Decode(&document); err != nil {
			return nil, fmt.Errorf("failed to parse json body: %w", err)
		}
		switch value := document.(type) {
		case map[string]any:
			if opts.Explain {
				// The search ran so that it could be profiled; its hits are
				// not what an explain asked for, and they leave here rather
				// than reaching a caller that masking no longer covers.
				profile, err := profileOnly(value)
				if err != nil {
					return nil, err
				}
				columns = append(columns, driver.Column{Name: profileField})
				row = append(row, driver.Value{Kind: driver.String, String: profile})
				result.Plan = explainQueryPlan()
				break
			}
			// Only a bulk request's own answer states its items' verdicts; a
			// document a caller reads may hold fields with the same names.
			if isBulkRoute(request.URL) {
				if failed, total, reason := bulkFailures(value); failed > 0 {
					return nil, fmt.Errorf("%d of %d bulk items failed: %s", failed, total, reason)
				}
			}
			// Object response: one column per key, ordered by name so one query
			// shows its columns the same way every run.
			keys := make([]string, 0, len(value))
			for key := range value {
				keys = append(keys, key)
			}
			slices.Sort(keys)
			for _, key := range keys {
				columns = append(columns, driver.Column{Name: key})
				encoded, err := json.Marshal(value[key])
				if err != nil {
					return nil, err
				}
				row = append(row, driver.Value{Kind: driver.String, String: string(encoded)})
			}
		case []any:
			if opts.Explain {
				return nil, errors.New("the profiled search answered with a list, which is not a search response")
			}
			// Array response (e.g., _cat APIs with format=json): one column
			// holding the array.
			columns = append(columns, driver.Column{Name: "result"})
			encoded, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			row = append(row, driver.Value{Kind: driver.String, String: string(encoded)})
		default:
			return nil, fmt.Errorf("unexpected JSON type: %T (expected object or array)", value)
		}
	case strings.Contains(contentType, "text/plain"):
		columns = append(columns, driver.Column{Name: "text/plain"})
		row = append(row, driver.Value{Kind: driver.String, String: string(respBytes)})
	default:
		return nil, fmt.Errorf("Content-Type not supported: %s", contentType)
	}
	result.Columns = columns
	result.Rows = [][]driver.Value{row}
	result.Latency = time.Since(startTime)
	return result, nil
}

// isBulkRoute reports whether a route names Elasticsearch's bulk endpoint, at
// "/_bulk" or under an index. A document read may carry the word anywhere in its
// path, such as /alpha/_source/app_bulk.
func isBulkRoute(route string) bool {
	path, _, _ := strings.Cut(route, "?")
	path = strings.TrimRight(path, "/")
	return path == "_bulk" || strings.HasSuffix(path, "/_bulk")
}

// bulkFailures reads a bulk answer's verdict: how many of its items failed, how
// many there were, and the first failure's reason. A bulk request answers 200
// whether or not its items applied, and states the failures only in the body.
func bulkFailures(answer map[string]any) (failed, total int, reason string) {
	errorsValue, ok := answer["errors"].(bool)
	if !ok || !errorsValue {
		return 0, 0, ""
	}
	items, ok := answer["items"].([]any)
	if !ok {
		return 0, 0, ""
	}
	for _, item := range items {
		operations, ok := item.(map[string]any)
		if !ok {
			continue
		}
		total++
		for _, operation := range operations {
			outcome, ok := operation.(map[string]any)
			if !ok {
				continue
			}
			failure, ok := outcome["error"].(map[string]any)
			if !ok {
				continue
			}
			failed++
			if reason == "" {
				if text, ok := failure["reason"].(string); ok {
					reason = text
				} else if text, ok := failure["type"].(string); ok {
					reason = text
				}
			}
		}
	}
	return failed, total, reason
}

// requestText is the statement a Result carries: the request as the server will
// read it.
func requestText(request *elasticsearch.Request, data ...[]byte) string {
	var body string
	if len(data) > 0 {
		body = string(data[0])
	} else {
		for _, item := range request.Data {
			body += strings.ReplaceAll(item, "\n", "") + "\n"
		}
	}
	return fmt.Sprintf("%s %s\n%s", request.Method, request.URL, body)
}
