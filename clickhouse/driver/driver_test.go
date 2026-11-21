package driver

import (
	"context"
	"crypto/tls"
	"math"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/bytebase/omni/driver"
)

func baseConfig() driver.Config {
	return driver.Config{User: "default", Host: "db.example.com", Port: "9000"}
}

func TestConnOptionsDefaults(t *testing.T) {
	options, err := connOptions(baseConfig())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := options.Addr, []string{"db.example.com:9000"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Addr = %v, want %v", got, want)
	}
	if options.Auth.Username != "default" || options.Auth.Database != "" {
		t.Errorf("Auth = %+v", options.Auth)
	}
	if options.Settings["max_execution_time"] != 300 {
		t.Errorf("max_execution_time = %v, want 300", options.Settings["max_execution_time"])
	}
	if options.DialTimeout != 10*time.Second {
		t.Errorf("DialTimeout = %v", options.DialTimeout)
	}
}

func TestConnOptionsDialWithoutHost(t *testing.T) {
	cfg := baseConfig()
	cfg.Host, cfg.Port = "", ""
	var dialed string
	cfg.Dial = func(_ context.Context, network, address string) (net.Conn, error) {
		if network != "tcp" {
			t.Errorf("network = %q", network)
		}
		dialed = address
		return nil, net.ErrClosed
	}
	options, err := connOptions(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(options.Addr) != 1 || options.Addr[0] == "" {
		t.Fatalf("Addr = %v", options.Addr)
	}
	_, _ = options.DialContext(context.Background(), options.Addr[0])
	if dialed != options.Addr[0] {
		t.Errorf("dialed %q, want %q", dialed, options.Addr[0])
	}
}

func TestConnOptionsIPv6(t *testing.T) {
	cfg := baseConfig()
	cfg.Host = "2001:db8::1"
	options, err := connOptions(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := options.Addr, []string{"[2001:db8::1]:9000"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Addr = %v, want %v", got, want)
	}
}

func TestConnOptionsParamsOverrideSettings(t *testing.T) {
	cfg := baseConfig()
	cfg.Database = "sales"
	cfg.Params = map[string]string{
		"max_execution_time": "5",
		"max_block_size":     "1024",
		"insert_deduplicate": "true",
	}
	options, err := connOptions(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if got := options.Settings["max_execution_time"]; got != int64(5) {
		t.Errorf("max_execution_time = %#v, want the int64 5", got)
	}
	if got := options.Settings["max_block_size"]; got != int64(1024) {
		t.Errorf("max_block_size = %#v, want the int64 1024", got)
	}
	if got := options.Settings["insert_deduplicate"]; got != "true" {
		t.Errorf("insert_deduplicate = %#v, want the string true", got)
	}
	if options.Auth.Database != "sales" {
		t.Errorf("Auth.Database = %q", options.Auth.Database)
	}
}

func TestConnOptionsTLSAndDial(t *testing.T) {
	cfg := baseConfig()
	var gotHost string
	cfg.TLS = func(host string) (*tls.Config, error) {
		gotHost = host
		return &tls.Config{ServerName: host}, nil
	}
	var dialedAddress string
	cfg.Dial = func(_ context.Context, network, address string) (net.Conn, error) {
		if network != "tcp" {
			t.Errorf("network = %q, want tcp", network)
		}
		dialedAddress = address
		return nil, net.ErrClosed
	}
	options, err := connOptions(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if gotHost != "db.example.com" || options.TLS == nil || options.TLS.ServerName != gotHost {
		t.Fatalf("gotHost=%q TLS=%+v", gotHost, options.TLS)
	}
	if _, err := options.DialContext(context.Background(), "db.example.com:9000"); err == nil {
		t.Error("DialContext did not fail")
	}
	if dialedAddress != "db.example.com:9000" {
		t.Errorf("dialed %q", dialedAddress)
	}
}

func TestConnOptionsTLSOverDial(t *testing.T) {
	cfg := baseConfig()
	cfg.TLS = func(string) (*tls.Config, error) {
		return &tls.Config{InsecureSkipVerify: true}, nil
	}
	cfg.Dial = func(context.Context, string, string) (net.Conn, error) {
		client, server := net.Pipe()
		_ = server.Close()
		return client, nil
	}
	options, err := connOptions(cfg)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := options.DialContext(context.Background(), options.Addr[0])
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if _, ok := conn.(*tls.Conn); !ok {
		t.Fatalf("dialed connection is %T, want TLS", conn)
	}
}

func TestConnOptionsValidation(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*driver.Config)
		want string
	}{
		{"blank host", func(c *driver.Config) { c.Host = "  " }, "host must be set"},
		{"blank port", func(c *driver.Config) { c.Port = "  " }, "port must be set"},
	} {
		cfg := baseConfig()
		tc.edit(&cfg)
		if _, err := connOptions(cfg); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: err = %v, want %q", tc.name, err, tc.want)
		}
	}
}

func TestOpenDoesNotConnect(t *testing.T) {
	cfg := baseConfig()
	cfg.Host, cfg.Port = "127.0.0.1", "1"
	d, err := Open(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Ping(context.Background()); err == nil {
		t.Error("Ping reached a server on port 1")
	}
	if err := d.Close(); err != nil {
		t.Error(err)
	}
}

func TestConvertMapWithFloatKeys(t *testing.T) {
	value := any(map[float64]string{1.5: "x"})
	got := convertValue("MAP(FLOAT64, STRING)", nil, &value)
	if got.Kind != driver.String || got.String != `{"1.5":"x"}` {
		t.Errorf("map = %+v", got)
	}
}

func TestConvertUInt8Array(t *testing.T) {
	value := any([]uint8{1, 2})
	got := convertValue("ARRAY(UINT8)", nil, &value)
	if want := `[1,2]`; got.Kind != driver.String || got.String != want {
		t.Errorf("Array(UInt8) = %+v, want %s", got, want)
	}
}

// TestDescribeExplain holds what a statement's plan says about itself, which the
// one plan ClickHouse cannot tape (the timings EXPLAIN ANALYZE prints move
// between runs) can only be read from the text.
func TestDescribeExplain(t *testing.T) {
	for _, tc := range []struct {
		statement, format string
		executed, ok      bool
	}{
		{statement: "SELECT 1"},
		{statement: "  -- leading comment\nSELECT 1"},
		{statement: "EXPLAIN SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN json = 1, indexes = 1 SELECT 1", format: "JSON", ok: true},
		{statement: "EXPLAIN ANALYZE SELECT 1", format: "TEXT", executed: true, ok: true},
		{statement: "EXPLAIN ANALYZE json = 1 SELECT 1", format: "JSON", executed: true, ok: true},
	} {
		plan := describeExplain(tc.statement)
		switch {
		case !tc.ok && plan != nil:
			t.Errorf("%q = %+v, want no plan", tc.statement, plan)
		case tc.ok && plan == nil:
			t.Errorf("%q = no plan, want %s", tc.statement, tc.format)
		case tc.ok && (plan.Format != tc.format || plan.Executed != tc.executed || plan.Content != ""):
			t.Errorf("%q = %+v, want %s executed=%v", tc.statement, plan, tc.format, tc.executed)
		}
	}
}

func TestConvertTemporal(t *testing.T) {
	at := time.Date(2026, 3, 4, 5, 6, 7, 891000000, time.UTC)
	for _, tc := range []struct {
		name      string
		typeName  string
		kind      driver.Kind
		precision int32
	}{
		{"date", "Date", driver.Timestamp, 0},
		{"datetime", "DateTime", driver.Timestamp, 0},
		{"datetime64", "DateTime64(3)", driver.Timestamp, 3},
		{"zoned datetime", "DateTime('UTC')", driver.TimestampTZ, 0},
		{"zoned datetime64", "DateTime64(3, 'UTC')", driver.TimestampTZ, 3},
		{"nullable zoned datetime64", "Nullable(DateTime64(6, 'Asia/Shanghai'))", driver.TimestampTZ, 6},
		{"low cardinality zoned datetime", "LowCardinality(DateTime('Asia/Shanghai'))", driver.TimestampTZ, 0},
		{"nullable low cardinality zoned datetime64", "Nullable(LowCardinality(DateTime64(3, 'UTC')))", driver.TimestampTZ, 3},
	} {
		got := convertValue(tc.typeName, nil, &at)
		if got.Kind != tc.kind || got.Precision != tc.precision {
			t.Errorf("%s: %+v, want %s precision %d", tc.name, got, tc.kind, tc.precision)
		}
	}
}

func TestConvertTimeDuration(t *testing.T) {
	duration := time.Hour + 2*time.Minute + 3*time.Second + 456*time.Millisecond
	if got := convertValue("TIME64", nil, &duration); got.Kind != driver.String || got.String != "01:02:03.456" {
		t.Errorf("Time64 = %+v", got)
	}
	pointer := &duration
	if got := convertValue("NULLABLE(TIME64)", nil, &pointer); got.Kind != driver.String || got.String != "01:02:03.456" {
		t.Errorf("nullable Time64 = %+v", got)
	}
	var nilDuration *time.Duration
	if got := convertValue("NULLABLE(TIME64)", nil, &nilDuration); got.Kind != driver.Null {
		t.Errorf("NULL Time64 = %+v", got)
	}
}

func TestConvertNonFiniteFloatInArray(t *testing.T) {
	value := any([]any{math.NaN(), math.Inf(1), math.Inf(-1)})
	got := convertValue("ARRAY(FLOAT64)", nil, &value)
	if want := `["nan","inf","-inf"]`; got.Kind != driver.String || got.String != want {
		t.Errorf("array of floats = %+v, want %s", got, want)
	}
}

func TestConvertUUIDInArray(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	value := any([]any{id})
	got := convertValue("ARRAY(UUID)", nil, &value)
	if want := `["550e8400-e29b-41d4-a716-446655440000"]`; got.Kind != driver.String || got.String != want {
		t.Errorf("array of UUID = %+v, want %s", got, want)
	}
}

func TestConvertIP(t *testing.T) {
	ip := net.ParseIP("2001:db8::1")
	if got := convertValue("IPV6", nil, &ip); got.Kind != driver.String || got.String != "2001:db8::1" {
		t.Errorf("IPv6 = %+v", got)
	}
	pointer := &ip
	if got := convertValue("NULLABLE(IPV6)", nil, &pointer); got.Kind != driver.String || got.String != "2001:db8::1" {
		t.Errorf("nullable IPv6 = %+v", got)
	}
	var nilIP *net.IP
	if got := convertValue("NULLABLE(IPV6)", nil, &nilIP); got.Kind != driver.Null {
		t.Errorf("NULL IPv6 = %+v", got)
	}
}

func TestSplitKeyExpressions(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want []string
	}{
		{"id", []string{"id"}},
		{"id, name", []string{"id", "name"}},
		{"a, tuple(b, c), d", []string{"a", "tuple(b, c)", "d"}},
		{"f('x,y'), `a,b`", []string{"f('x,y')", "`a,b`"}},
	} {
		if got := splitKeyExpressions(tc.in); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("splitKeyExpressions(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
