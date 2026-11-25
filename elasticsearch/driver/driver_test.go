package driver

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/elasticsearch"
)

func TestEndpoint(t *testing.T) {
	for _, tc := range []struct {
		name   string
		host   string
		port   string
		secure bool
		want   string
	}{
		{name: "host and port", host: "es.example.com", port: "9200", want: "http://es.example.com:9200"},
		{name: "default port", host: "es.example.com", want: "http://es.example.com:9200"},
		{name: "tls", host: "es.example.com", port: "9243", secure: true, want: "https://es.example.com:9243"},
		{name: "host with a scheme", host: "https://search.example.com", want: "https://search.example.com"},
		{name: "host with its own port", host: "es.example.com:9243", port: "9200", want: "http://es.example.com:9243"},
		{name: "host with its own port under tls", host: "es.example.com:9243", secure: true, want: "https://es.example.com:9243"},
		{name: "ipv6 literal", host: "::1", port: "9200", want: "http://[::1]:9200"},
		{name: "ipv6 literal with a port", host: "[::1]:9243", want: "http://[::1]:9243"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := endpoint(tc.host, tc.port, tc.secure)
			require.NoError(t, err)
			require.Equal(t, tc.want, got.String())
		})
	}
	_, err := endpoint("", "", false)
	require.EqualError(t, err, "host must be set")

	// A scheme that contradicts Config.TLS is refused, not silently used.
	_, err = endpoint("http://es.example.com", "", true)
	require.EqualError(t, err, `address "http://es.example.com" is http, but TLS is required`)
}

func TestTLSSettings(t *testing.T) {
	// Nil TLS leaves the default negotiation: plain http, no config.
	got, err := tlsSettings(driver.Config{}, "es.example.com")
	require.NoError(t, err)
	require.Nil(t, got)

	var asked string
	got, err = tlsSettings(driver.Config{TLS: func(host string) (*tls.Config, error) {
		asked = host
		return &tls.Config{InsecureSkipVerify: true}, nil
	}}, "es.example.com")
	require.NoError(t, err)
	require.Equal(t, "es.example.com", asked)
	require.Equal(t, uint16(tls.VersionTLS12), got.MinVersion)
	// The caller's config keeps the version it had: the driver raises a copy.
	require.Equal(t, uint16(0), (&tls.Config{InsecureSkipVerify: true}).MinVersion)

	_, err = tlsSettings(driver.Config{TLS: func(string) (*tls.Config, error) { return nil, nil }}, "es.example.com")
	require.EqualError(t, err, "TLS configuration must not be nil")
}

func TestUnitConversion(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want int64
	}{
		{in: "", want: 0},
		{in: "249b", want: 249},
		{in: "4.9kb", want: 5017},
		{in: "1.5mb", want: 1572864},
		{in: "2gb", want: 2 * 1024 * 1024 * 1024},
		{in: "1tb", want: 1024 * 1024 * 1024 * 1024},
		{in: "1pb", want: 1024 * 1024 * 1024 * 1024 * 1024},
		{in: "1eb", want: 1024 * 1024 * 1024 * 1024 * 1024 * 1024},
	} {
		got, err := unitConversion(tc.in)
		require.NoError(t, err, tc.in)
		require.Equal(t, tc.want, got, tc.in)
	}
	_, err := unitConversion("big")
	require.EqualError(t, err, `unexpected size format: "big"`)
}

func TestParseDocCount(t *testing.T) {
	got, err := parseDocCount("")
	require.NoError(t, err)
	require.Zero(t, got)
	got, err = parseDocCount("42")
	require.NoError(t, err)
	require.Equal(t, int64(42), got)
	_, err = parseDocCount("many")
	require.Error(t, err)
}

func TestSearchable(t *testing.T) {
	for _, tc := range []struct {
		name    string
		request elasticsearch.Request
		wantErr string
	}{
		{name: "search", request: elasticsearch.Request{Method: "POST", URL: "/index/_search"}},
		{name: "uri search", request: elasticsearch.Request{Method: "GET", URL: "/index/_search?q=name"}},
		{name: "root search", request: elasticsearch.Request{Method: "GET", URL: "_search"}},
		{name: "put search", request: elasticsearch.Request{Method: "PUT", URL: "/index/_search"}, wantErr: "PUT is not a search to explain"},
		{name: "other endpoint", request: elasticsearch.Request{Method: "POST", URL: "/index/_doc/1"}, wantErr: "/index/_doc/1 is not a search to explain"},
		// A document whose id is _search is a write, not a search.
		{name: "document named _search", request: elasticsearch.Request{Method: "POST", URL: "/index/_doc/_search"}, wantErr: "/index/_doc/_search is not a search to explain"},
		{name: "ndjson", request: elasticsearch.Request{Method: "POST", URL: "/_msearch", Data: []string{"{}", "{}"}}, wantErr: "/_msearch is not a search to explain"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := searchable(&tc.request)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

// TestProfiledSearch holds the profile flag to the caller's body, with every
// other value left as written.
func TestProfiledSearch(t *testing.T) {
	request := &elasticsearch.Request{
		Method: "POST",
		URL:    "/index/_search",
		Data:   []string{`{"query": {"term": {"n": 123456789012345678901}}}`},
	}
	profiled, err := profiledSearch(request)
	require.NoError(t, err)
	require.Len(t, profiled.Data, 1)
	var body map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(profiled.Data[0]), &body))
	require.Equal(t, "true", string(body["profile"]))
	require.Contains(t, profiled.Data[0], "123456789012345678901")
	// The caller's request is untouched.
	require.Equal(t, `{"query": {"term": {"n": 123456789012345678901}}}`, request.Data[0])

	_, err = profiledSearch(&elasticsearch.Request{Method: "GET", URL: "/_search", Data: []string{"[]"}})
	require.ErrorContains(t, err, "the search body is not a JSON object to profile")
}

func TestProfileOnly(t *testing.T) {
	plan, err := profileOnly(map[string]any{
		"took":    float64(3),
		"profile": map[string]any{"shards": []any{}},
		"hits":    map[string]any{"total": float64(3)},
	})
	require.NoError(t, err)
	require.Equal(t, `{"profile":{"shards":[]},"took":3}`, plan)

	_, err = profileOnly(map[string]any{"hits": map[string]any{}})
	require.EqualError(t, err, "the search answered with no profile to explain")
}

// TestOpenCarriesTheLogin checks the driver's request against a live listener:
// the login reaches the endpoint as basic auth and the route keeps its query.
func TestOpenCarriesTheLogin(t *testing.T) {
	var gotAuth, gotRoute, gotContentType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotRoute = r.URL.RequestURI()
		gotContentType = r.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version": {"number": "8.19.25"}}`))
	}))
	defer srv.Close()

	d, err := Open(context.Background(), driver.Config{
		Host:     srv.URL,
		User:     "elastic",
		Password: "secret",
	})
	require.NoError(t, err)
	require.NoError(t, d.Ping(context.Background()))
	require.Equal(t, "Basic ZWxhc3RpYzpzZWNyZXQ=", gotAuth)
	require.Equal(t, "/", gotRoute)
	require.Equal(t, "application/json", gotContentType)

	res, err := d.client.Do(context.Background(), http.MethodGet, "/_cat/indices?format=json&pretty", nil)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, "/_cat/indices?format=json&pretty", gotRoute)
}

// TestOpenUsesTheTransport holds a caller's transport to the requests: a
// managed OpenSearch domain signs each one, and the driver must not wrap it in
// a transport of its own.
func TestOpenUsesTheTransport(t *testing.T) {
	var called int
	var secure bool
	d, err := Open(context.Background(), driver.Config{Host: "search.example.com"}, Options{
		OpenSearch: true,
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			called++
			secure = r.URL.Scheme == "https"
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": {"application/json"}},
				Body:       http.NoBody,
			}, nil
		}),
	})
	require.NoError(t, err)
	require.NoError(t, d.Ping(context.Background()))
	require.Equal(t, 1, called)
	require.True(t, secure, "a transport means https")
	require.Equal(t, "/_plugins/_security/api/internalusers", d.securityRoute())

	// The caller's transport dials, so a dialer beside it is refused rather
	// than silently ignored.
	_, err = Open(context.Background(), driver.Config{
		Host: "search.example.com",
		Dial: func(context.Context, string, string) (net.Conn, error) { return nil, errors.New("not dialed") },
	}, Options{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, nil })})
	require.EqualError(t, err, "Config.Dial and Options.Transport both open the endpoint's connections: the caller's transport dials, so Config.Dial cannot be honoured")
}

// TestOpenDialsThroughConfigDial holds the endpoint to Config.Dial: the driver
// opens the socket through the caller's tunnel or cloud connector instead of
// resolving the host itself, and TLS is negotiated over what it returns.
func TestOpenDialsThroughConfigDial(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version": {"number": "8.19.25"}}`))
	}))
	defer srv.Close()

	var asked []string
	d, err := Open(context.Background(), driver.Config{
		// A name no resolver knows: only the dialer reaches anything.
		Host: "elasticsearch.example.com",
		Port: "9200",
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			asked = append(asked, network+" "+address)
			var dialer net.Dialer
			return dialer.DialContext(ctx, network, srv.Listener.Addr().String())
		},
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, d.Close()) }()
	require.NoError(t, d.Ping(context.Background()))
	require.Equal(t, []string{"tcp elasticsearch.example.com:9200"}, asked)
}

// TestTheRouteIsSentAsWritten holds a route to the endpoint: a path that names a
// document is not cleaned into the index it sits in.
func TestTheRouteIsSentAsWritten(t *testing.T) {
	var got []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = append(got, r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	d, err := Open(context.Background(), driver.Config{Host: srv.URL})
	require.NoError(t, err)
	defer func() { _ = d.Close() }()
	for _, route := range []string{"/alpha/_doc/..", "/alpha/_doc/1", "/_cat/indices?format=json&pretty"} {
		res, err := d.client.Do(context.Background(), http.MethodGet, route, nil)
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
	}
	require.Equal(t, []string{"/alpha/_doc/..", "/alpha/_doc/1", "/_cat/indices?format=json&pretty"}, got)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestTruncate(t *testing.T) {
	require.Equal(t, "abc", truncate([]byte("abc"), 3))
	require.Equal(t, "abc", truncate([]byte("abc"), 10))
	require.Equal(t, "ab...", truncate([]byte("abcdef"), 2))
	// The limit counts code points, and the cut lands between them.
	require.Equal(t, "éé...", truncate([]byte("éééé"), 2))
	require.True(t, strings.HasSuffix(truncate([]byte("abcdef"), 2), "..."))
}

// TestGetIndicesFiltersHiddenAndSortsByName unit-tests defensive filtering of
// hidden indices if a server or proxy includes them in _cat/indices, along with
// deterministic name ordering.
func TestGetIndicesFiltersHiddenAndSortsByName(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.RequestURI() {
		case "/*/_settings?filter_path=**.hidden&expand_wildcards=open,hidden":
			_, _ = w.Write([]byte(`{
				"alpha":{"settings":{"index":{"hidden":"false"}}},
				"beta":{"settings":{"index":{"hidden":"false"}}},
				"hidden-one":{"settings":{"index":{"hidden":"true"}}}
			}`))
		case "/_cat/indices?format=json&pretty":
			_, _ = w.Write([]byte(`[
				{"index":"beta","docs.count":"2","store.size":"227b"},
				{"index":"hidden-one","docs.count":"1","store.size":"1kb"},
				{"index":"alpha","docs.count":"4","store.size":"1.5kb"}
			]`))
		case "/_alias":
			_, _ = w.Write([]byte(`{"alpha":{"aliases":{}},"beta":{"aliases":{}},"hidden-one":{"aliases":{}}}`))
		}
	}))
	defer srv.Close()

	d, err := Open(context.Background(), driver.Config{Host: srv.URL})
	require.NoError(t, err)
	defer func() { _ = d.Close() }()

	indices, err := d.getIndices(context.Background())
	require.NoError(t, err)
	require.Len(t, indices, 2)
	require.Equal(t, "alpha", indices[0].Name)
	require.Equal(t, int64(4), indices[0].RowCount)
	require.Equal(t, int64(1536), indices[0].DataSize)
	require.Equal(t, "beta", indices[1].Name)
	require.Equal(t, int64(2), indices[1].RowCount)
	require.Equal(t, int64(227), indices[1].DataSize)
}

// TestSyncInstanceFailsOpenOnUsers pins that a security-endpoint failure leaves
// Roles empty without failing SyncInstance.
func TestSyncInstanceFailsOpenOnUsers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.RequestURI() {
		case "/":
			_, _ = w.Write([]byte(`{"version":{"number":"8.19.25"}}`))
		case "/*/_settings?filter_path=**.hidden&expand_wildcards=open,hidden":
			_, _ = w.Write([]byte(`{}`))
		case "/_cat/indices?format=json&pretty":
			_, _ = w.Write([]byte(`[]`))
		case "/_alias":
			_, _ = w.Write([]byte(`{}`))
		case "/_security/user":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{}`))
		}
	}))
	defer srv.Close()

	d, err := Open(context.Background(), driver.Config{Host: srv.URL})
	require.NoError(t, err)
	defer func() { _ = d.Close() }()

	instance, err := d.SyncInstance(context.Background())
	require.NoError(t, err)
	require.Equal(t, "8.19.25", instance.Version)
	require.Empty(t, instance.Roles)
}

func TestQueryUnsupportedContentTypeAndMissingProfile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/image":
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write([]byte("not json"))
		case "/alpha/_search":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"took": 1, "hits": {"total": {"value": 0}}}`))
		}
	}))
	defer srv.Close()

	d, err := Open(context.Background(), driver.Config{Host: srv.URL})
	require.NoError(t, err)
	defer func() { _ = d.Close() }()
	sess := openSession(t, d)

	res, err := sess.Query(context.Background(), "GET /image", driver.QueryOptions{})
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.ErrorContains(t, res[0].Err, "Content-Type not supported: image/png")

	res, err = sess.Query(context.Background(), "POST /alpha/_search\n{}", driver.QueryOptions{Explain: true})
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.ErrorContains(t, res[0].Err, "the search answered with no profile to explain")
}
