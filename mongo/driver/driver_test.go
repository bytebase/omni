package driver

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/mongo/internal/testcontainer"
)

// openDriver returns a driver on t's tape, for databaseName.
func openDriver(t *testing.T, databaseName string) *Driver {
	t.Helper()
	d, _ := openDriverWithClient(t, databaseName)
	return d
}

// openDriverWithClient returns the driver and the client behind it, for a test
// that asks the client something the driver does not expose.
func openDriverWithClient(t *testing.T, databaseName string) (*Driver, *testcontainer.Client) {
	t.Helper()
	c := testcontainer.Open(t)
	t.Cleanup(func() { _ = c.Disconnect(context.Background()) })
	return &Driver{client: c, databaseName: databaseName}, c
}

func TestConnectionURI(t *testing.T) {
	tests := []struct {
		name string
		cfg  driver.Config
		opts Options
		want string
	}{
		{
			name: "host and port",
			cfg:  driver.Config{Host: "localhost", Port: "27017"},
			want: "mongodb://localhost:27017/?appName=bytebase&authSource=admin",
		},
		{
			name: "a parameter overrides a default",
			cfg: driver.Config{
				Host:   "localhost",
				Port:   "27017",
				Params: map[string]string{"directConnection": "true", "appName": "bytebase-42"},
			},
			want: "mongodb://localhost:27017/?appName=bytebase-42&authSource=admin&directConnection=true",
		},
		{
			name: "database",
			cfg:  driver.Config{Host: "localhost", Port: "27017", Database: "sampleDB"},
			want: "mongodb://localhost:27017/sampleDB?appName=bytebase&authSource=admin",
		},
		{
			name: "srv with the credentials in the endpoint",
			cfg: driver.Config{
				Host:     "cluster0.sample.mongodb.net",
				User:     "bytebase",
				Password: "passwd",
				Database: "sampleDB",
			},
			opts: Options{SRV: true},
			want: "mongodb+srv://bytebase:passwd@cluster0.sample.mongodb.net/sampleDB?appName=bytebase&authSource=admin",
		},
		{
			name: "connection string options",
			cfg: driver.Config{
				Host:     "cluster0.sample.mongodb.net",
				User:     "bytebase",
				Password: "passwd",
				Database: "sampleDB",
				Params: map[string]string{
					"readPreference":     "secondary",
					"readPreferenceTags": "dc:ny",
				},
			},
			opts: Options{SRV: true},
			want: "mongodb+srv://bytebase:passwd@cluster0.sample.mongodb.net/sampleDB?appName=bytebase&authSource=admin&readPreference=secondary&readPreferenceTags=dc%3Any",
		},
		{
			name: "a host list carries its own ports",
			cfg: driver.Config{
				Host:     "node1.cluster0.sample.mongodb.net:27017,node2.cluster0.sample.mongodb.net:27017,node3.cluster0.sample.mongodb.net:27017",
				User:     "bytebase",
				Password: "passwd",
				Database: "sampleDB",
				Params:   map[string]string{"replicaSet": "rs0"},
			},
			want: "mongodb://bytebase:passwd@node1.cluster0.sample.mongodb.net:27017,node2.cluster0.sample.mongodb.net:27017,node3.cluster0.sample.mongodb.net:27017/sampleDB?appName=bytebase&authSource=admin&replicaSet=rs0",
		},
		{
			name: "an IPv6 host keeps its brackets",
			cfg:  driver.Config{Host: "2001:db8::1", Port: "27017"},
			want: "mongodb://[2001:db8::1]:27017/?appName=bytebase&authSource=admin",
		},
		{
			name: "a bracketed IPv6 host takes the port once",
			cfg:  driver.Config{Host: "[2001:db8::1]", Port: "27017"},
			want: "mongodb://[2001:db8::1]:27017/?appName=bytebase&authSource=admin",
		},
		{
			name: "TLS takes the scheme's own TLS over",
			cfg: driver.Config{
				Host:   "cluster0.sample.mongodb.net",
				Params: map[string]string{"tls": "true"},
				TLS:    func(string) (*tls.Config, error) { return &tls.Config{}, nil },
			},
			opts: Options{SRV: true},
			want: "mongodb+srv://cluster0.sample.mongodb.net/?appName=bytebase&authSource=admin&tls=false",
		},
		{
			name: "TLS takes the ssl alias over too",
			cfg: driver.Config{
				Host:   "localhost",
				Port:   "27017",
				Params: map[string]string{"SSL": "true"},
				TLS:    func(string) (*tls.Config, error) { return &tls.Config{}, nil },
			},
			want: "mongodb://localhost:27017/?appName=bytebase&authSource=admin&tls=false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := connectionURI(tt.cfg, tt.opts)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestConnectionURIRequiresAHost(t *testing.T) {
	_, err := connectionURI(driver.Config{}, Options{})
	require.ErrorContains(t, err, "host must be set")
}

// TestConnectionURICarriesTheCredentials holds the connection string to what it
// carries: the endpoint, the resolved login and the database.
func TestOpenCarriesTheCredentialsAndTLSIntoTheClient(t *testing.T) {
	uri, err := connectionURI(driver.Config{
		Host:            "mongo.example.com",
		Port:            "27017",
		User:            "app",
		Password:        "p@ss/word?x=1",
		Database:        "app",
		ApplicationName: "bytebase-42",
	}, Options{})
	require.NoError(t, err)
	parsed, err := url.Parse(uri)
	require.NoError(t, err)
	require.Equal(t, "mongodb", parsed.Scheme)
	require.Equal(t, "mongo.example.com:27017", parsed.Host)
	password, ok := parsed.User.Password()
	require.True(t, ok)
	require.Equal(t, "p@ss/word?x=1", password)
	require.Equal(t, "bytebase-42", parsed.Query().Get("appName"))
	require.Equal(t, "/app", parsed.Path)
}

// tlsDialerHosts records the hosts Config.TLS was asked about.
func tlsDialerHosts(hosts *[]string, cfg *tls.Config, err error) tlsDialer {
	return tlsDialer{config: driver.Config{TLS: func(host string) (*tls.Config, error) {
		*hosts = append(*hosts, host)
		return cfg, err
	}}}
}

// TestTLSDialerResolvesSettingsPerHost holds the contract's TLS promise: the
// caller's func is asked about one host at a time, so a replica set's members
// are each verified against their own name.
func TestTLSDialerResolvesSettingsPerHost(t *testing.T) {
	var hosts []string
	d := tlsDialerHosts(&hosts, nil, errors.New("no TLS here"))
	for _, address := range []string{"node1.example.com:27017", "[2001:db8::1]:27017"} {
		_, err := d.DialContext(context.Background(), "tcp", address)
		require.ErrorContains(t, err, "no TLS here")
	}
	require.Equal(t, []string{"node1.example.com", "2001:db8::1"}, hosts)
}

func TestTLSDialerRefusesANilConfig(t *testing.T) {
	var hosts []string
	d := tlsDialerHosts(&hosts, nil, nil)
	_, err := d.DialContext(context.Background(), "tcp", "127.0.0.1:27017")
	require.ErrorContains(t, err, "TLS configuration must not be nil")
}

// TestTLSDialerHandshakesWithTheHostsSettings holds the connection the dialer
// hands back: TLS on top of the network connection, verified with what the
// caller's func returned for the host it dialed.
func TestTLSDialerHandshakesWithTheHostsSettings(t *testing.T) {
	certificate := selfSignedCertificate(t)
	listener, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{Certificates: []tls.Certificate{certificate}})
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		if tlsConn, ok := conn.(*tls.Conn); ok {
			_ = tlsConn.Handshake()
		}
		_ = conn.Close()
	}()

	var hosts []string
	d := tlsDialerHosts(&hosts, &tls.Config{InsecureSkipVerify: true}, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := d.DialContext(ctx, "tcp", listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()
	require.IsType(t, &tls.Conn{}, conn)
	host, _, err := net.SplitHostPort(listener.Addr().String())
	require.NoError(t, err)
	require.Equal(t, []string{host}, hosts)
}

// selfSignedCertificate returns a certificate a TLS listener can serve.
func selfSignedCertificate(t *testing.T) tls.Certificate {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "omni"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	require.NoError(t, err)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}
}

func TestIsSystemCollection(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "users", want: false},
		{name: "system.namespaces", want: true},
		{name: "system.users", want: true},
		{name: "system.buckets.events", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, isSystemCollection(tt.name))
		})
	}
}

func TestAnnotateUnsupportedServerVersion(t *testing.T) {
	// Copied verbatim from the driver's compatibility error, which is what the
	// annotation matches on: reword this and it stops being the case under test.
	refused := errors.New("server at 10.0.0.1:27017 reports wire version 8, but this version of the Go driver requires at least 9 (MongoDB 4.4)")

	annotated := annotateUnsupportedServerVersion(refused)
	require.ErrorContains(t, annotated, "MongoDB "+MinSupportedServerVersion+" or later is required")
	require.ErrorContains(t, annotated, "reports wire version 8")
	require.ErrorIs(t, annotated, refused)

	require.NoError(t, annotateUnsupportedServerVersion(nil))

	authFailure := errors.New("auth error: sasl conversation error")
	require.Equal(t, authFailure, annotateUnsupportedServerVersion(authFailure))
}
