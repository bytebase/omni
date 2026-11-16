package driver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/stretchr/testify/require"
	"gopkg.in/inf.v0"

	"github.com/bytebase/omni/driver"
)

func TestOpenRequiresAHost(t *testing.T) {
	_, err := Open(context.Background(), driver.Config{})
	require.EqualError(t, err, "host must be set")
}

func TestEndpoints(t *testing.T) {
	for _, tc := range []struct {
		name string
		host string
		port string
		want []string
	}{
		{name: "host and port", host: "db.example.com", port: "9042", want: []string{"db.example.com:9042"}},
		{name: "host only", host: "db.example.com", want: []string{"db.example.com"}},
		{name: "host with its own port", host: "db.example.com:9142", port: "9042", want: []string{"db.example.com:9142"}},
		{name: "list", host: "a.example.com:9042,b.example.com:9142", want: []string{"a.example.com:9042", "b.example.com:9142"}},
		{name: "ipv6 literal", host: "::1", port: "9042", want: []string{"[::1]:9042"}},
		{name: "ipv6 literal with a port", host: "[::1]:9142", want: []string{"[::1]:9142"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := endpoints(tc.host, tc.port)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
	for _, bad := range []string{"", "a.example.com,,b.example.com"} {
		_, err := endpoints(bad, "")
		require.Error(t, err)
	}
}

func TestTypeName(t *testing.T) {
	want := map[gocql.Type]string{
		gocql.TypeCustom: "custom", gocql.TypeAscii: "ascii", gocql.TypeBigInt: "bigint",
		gocql.TypeBlob: "blob", gocql.TypeBoolean: "boolean", gocql.TypeCounter: "counter",
		gocql.TypeDecimal: "decimal", gocql.TypeDouble: "double", gocql.TypeFloat: "float",
		gocql.TypeInt: "int", gocql.TypeText: "text", gocql.TypeTimestamp: "timestamp",
		gocql.TypeUUID: "uuid", gocql.TypeVarchar: "varchar", gocql.TypeVarint: "varint",
		gocql.TypeTimeUUID: "timeuuid", gocql.TypeInet: "inet", gocql.TypeDate: "date",
		gocql.TypeTime: "time", gocql.TypeSmallInt: "smallint", gocql.TypeTinyInt: "tinyint",
		gocql.TypeDuration: "duration", gocql.TypeList: "list", gocql.TypeMap: "map",
		gocql.TypeSet: "set", gocql.TypeUDT: "udt", gocql.TypeTuple: "tuple",
	}
	for typ, name := range want {
		require.Equal(t, name, typeToString(typ), "type %d", typ)
	}
	require.Equal(t, "unknown(99)", typeToString(gocql.Type(99)))
}

func TestTimeOfDay(t *testing.T) {
	for _, tc := range []struct {
		in   time.Duration
		want string
	}{
		{in: 0, want: "00:00:00"},
		{in: 8*time.Hour + 12*time.Minute + 54*time.Second + 123456789*time.Nanosecond, want: "08:12:54.123456"},
		{in: 8 * time.Hour, want: "08:00:00"},
		{in: time.Hour + time.Microsecond, want: "01:00:00.000001"},
		// A microsecond's remainder is dropped, as Bytebase's microseconds did.
		{in: time.Hour + 999*time.Nanosecond, want: "01:00:00"},
		{in: 24 * time.Hour, want: "00:00:00"},
	} {
		require.Equal(t, tc.want, timeOfDayString(tc.in))
	}
}

func TestConvertRowValue(t *testing.T) {
	uuid, err := gocql.ParseUUID("6afaf7c2-c0c2-0ddc-7f22-71b8f0a0b0c0")
	require.NoError(t, err)
	decimal, ok := new(inf.Dec).SetString("-1.25")
	require.True(t, ok)
	varint, ok := new(big.Int).SetString("-123456789012345678901234567890", 10)
	require.True(t, ok)
	duration := time.Date(2026, 1, 2, 3, 4, 5, 678000000, time.UTC)
	cqlDuration := gocql.Duration{Months: 1, Days: 2, Nanoseconds: 3 * int64(time.Hour)}
	blob := []byte{0, 1}
	empty := ""

	for _, tc := range []struct {
		name string
		in   any
		want driver.Value
	}{
		{name: "string", in: &empty, want: driver.Value{Kind: driver.String}},
		{name: "nil string", in: (*string)(nil), want: driver.Value{Kind: driver.Null}},
		{name: "int64", in: new(int64(-7)), want: driver.Value{Kind: driver.Int64, Int64: -7}},
		{name: "int", in: new(7), want: driver.Value{Kind: driver.Int64, Int64: 7}},
		{name: "int16", in: new(int16(8)), want: driver.Value{Kind: driver.Int32, Int64: 8}},
		{name: "int8", in: new(int8(9)), want: driver.Value{Kind: driver.Int32, Int64: 9}},
		{name: "bool", in: new(true), want: driver.Value{Kind: driver.Bool, Bool: true}},
		{name: "float32", in: new(float32(1.5)), want: driver.Value{Kind: driver.Float32, Float64: 1.5}},
		{name: "float64", in: new(2.5), want: driver.Value{Kind: driver.Float64, Float64: 2.5}},
		{name: "bytes", in: &blob, want: driver.Value{Kind: driver.Bytes, Bytes: blob}},
		{name: "timestamp", in: &duration, want: driver.Value{Kind: driver.Timestamp, Time: duration, Precision: 3}},
		{name: "time of day", in: new(2*time.Hour + 3*time.Second), want: driver.Value{Kind: driver.String, String: "02:00:03"}},
		{name: "uuid", in: &uuid, want: driver.Value{Kind: driver.String, String: "6afaf7c2-c0c2-0ddc-7f22-71b8f0a0b0c0"}},
		{name: "decimal", in: &decimal, want: driver.Value{Kind: driver.String, String: "-1.25"}},
		{name: "nil decimal", in: new(*inf.Dec), want: driver.Value{Kind: driver.Null}},
		{name: "varint", in: &varint, want: driver.Value{Kind: driver.String, String: "-123456789012345678901234567890"}},
		{name: "nil varint", in: new(*big.Int), want: driver.Value{Kind: driver.Null}},
		{name: "duration", in: &cqlDuration, want: driver.Value{Kind: driver.String, String: "1mo2d3h0m0s"}},
		{name: "negative duration days", in: &gocql.Duration{Days: -1}, want: driver.Value{Kind: driver.String, String: "-1d0s"}},
		{name: "negative duration months", in: &gocql.Duration{Months: -2, Nanoseconds: -5}, want: driver.Value{Kind: driver.String, String: "-2mo-5ns"}},
		{name: "list", in: &[]int{1, 2}, want: driver.Value{Kind: driver.String, String: "[1,2]"}},
		{name: "map", in: &map[string]int{"a": 1}, want: driver.Value{Kind: driver.String, String: `{"a":1}`}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, convertRowValue(tc.in))
		})
	}
}

// TestValueGroups holds gocql's tuple expansion to one group per column: a
// tuple column answers its elements as one value, every other column one.
func TestValueGroups(t *testing.T) {
	scalar := testType{gocql.TypeInt, 0}
	tuple := testType{gocql.TypeTuple, 2}
	three := testType{gocql.TypeTuple, 3}
	col := func(name string, info testType) gocql.ColumnInfo {
		return gocql.ColumnInfo{Name: name, TypeInfo: info}
	}
	require.Equal(t, []int{1, 1, 1}, valueGroups([]gocql.ColumnInfo{col("a", scalar), col("b", scalar), col("c", scalar)}, nil))
	require.Equal(t, []int{1, 1, 1}, valueGroups([]gocql.ColumnInfo{col("a", scalar), col("b", scalar), col("c", scalar)}, []gocql.RowData{{Columns: []string{"a", "b", "c"}}}))
	// A tuple in the middle: a, tup[0], tup[1], c.
	require.Equal(t, []int{1, 2, 1}, valueGroups([]gocql.ColumnInfo{col("a", scalar), col("tup", tuple), col("c", scalar)}, []gocql.RowData{{Columns: []string{"a", "tup[0]", "tup[1]", "c"}}}))
	// A tuple of three, first.
	require.Equal(t, []int{3, 1}, valueGroups([]gocql.ColumnInfo{col("tup", three), col("b", scalar)}, []gocql.RowData{{Columns: []string{"tup[0]", "tup[1]", "tup[2]", "b"}}}))
	// The same tuple selected twice: its names repeat, and the declared element
	// count tells the two groups apart.
	require.Equal(t, []int{2, 2}, valueGroups([]gocql.ColumnInfo{col("tup", tuple), col("tup", tuple)}, []gocql.RowData{{Columns: []string{"tup[0]", "tup[1]", "tup[0]", "tup[1]"}}}))
	require.Equal(t, []int{2, 1}, valueGroups([]gocql.ColumnInfo{col("tup", tuple), col("id", scalar)}, []gocql.RowData{{Columns: []string{"tup[0]", "tup[1]", "id"}}}))
	// A scalar alias that looks like a tuple element is one value: the column's
	// type, not its name, says what it holds.
	require.Equal(t, []int{1, 1, 2}, valueGroups([]gocql.ColumnInfo{col("x[0]", scalar), col("x[1]", scalar), col("tup", tuple)}, []gocql.RowData{{Columns: []string{"x[0]", "x[1]", "tup[0]", "tup[1]"}}}))
	// A scalar alias that continues a real tuple's numbering is a value of its
	// own: the tuple's declared element count bounds its group.
	require.Equal(t, []int{2, 1}, valueGroups([]gocql.ColumnInfo{col("tup", tuple), col("tup[2]", scalar)}, []gocql.RowData{{Columns: []string{"tup[0]", "tup[1]", "tup[2]"}}}))
	// A grouping that does not consume every value is not one to align by.
	require.Equal(t, []int{1, 1}, valueGroups([]gocql.ColumnInfo{col("a", scalar), col("b", tuple)}, []gocql.RowData{{Columns: []string{"a", "b", "c"}}}))
}

// testType is a CQL type for the grouping test: only the type itself is read.
type testType struct {
	typ      gocql.Type
	elements int
}

func (t testType) Type() gocql.Type          { return t.typ }
func (t testType) Elements() int             { return t.elements }
func (testType) Zero() any                   { return nil }
func (testType) Marshal(any) ([]byte, error) { return nil, nil }
func (testType) Unmarshal([]byte, any) error { return nil }

func TestConvertTuple(t *testing.T) {
	seven := 7
	name := "seven"
	require.Equal(t, driver.Value{Kind: driver.String, String: `[7,"seven"]`}, convertTuple([]any{&seven, &name}))
	require.Equal(t, driver.Value{Kind: driver.String, String: `[null,"seven"]`}, convertTuple([]any{(*int)(nil), &name}))
	huge := []int{9223372036854775807}
	require.Equal(t, driver.Value{Kind: driver.String, String: `[[9223372036854775807]]`}, convertTuple([]any{&huge}))
}

func TestDialHostUsesTLS(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer srv.Close()
	addr := srv.Listener.Addr().String()
	hostname, port, err := net.SplitHostPort(addr)
	require.NoError(t, err)
	n, err := strconv.Atoi(port)
	require.NoError(t, err)
	host, err := gocql.NewHostInfoFromAddrPort(net.ParseIP(hostname), n)
	require.NoError(t, err)

	pool := x509.NewCertPool()
	pool.AddCert(srv.Certificate())
	var asked string
	d := &hostDialer{config: driver.Config{TLS: func(host string) (*tls.Config, error) {
		asked = host
		return &tls.Config{RootCAs: pool}, nil
	}}}
	dialed, err := d.DialHost(context.Background(), host)
	require.NoError(t, err)
	defer dialed.Conn.Close()
	require.Equal(t, hostname, asked)
	_, isTLS := dialed.Conn.(*tls.Conn)
	require.True(t, isTLS, "want a TLS connection, got %T", dialed.Conn)
}

func TestDialHostUsesDial(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	var addressed string
	d := &hostDialer{config: driver.Config{Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
		addressed = network + " " + address
		return net.Dial(network, address)
	}}}
	hostname, portText, err := net.SplitHostPort(listener.Addr().String())
	require.NoError(t, err)
	port, err := strconv.Atoi(portText)
	require.NoError(t, err)
	host, err := gocql.NewHostInfoFromAddrPort(net.ParseIP(hostname), port)
	require.NoError(t, err)
	dialed, err := d.DialHost(context.Background(), host)
	require.NoError(t, err)
	defer dialed.Conn.Close()
	require.Equal(t, "tcp "+listener.Addr().String(), addressed)
	_, isTLS := dialed.Conn.(*tls.Conn)
	require.False(t, isTLS, "want the connection Config.Dial answered, got %T", dialed.Conn)
}
