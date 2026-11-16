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
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/redis/internal/testcontainer"
)

// openDriver returns a driver on t's tape, as login on database index.
func openDriver(t *testing.T, login testcontainer.Login, database string) *Driver {
	t.Helper()
	c := testcontainer.Open(t, login)
	return &Driver{client: c, databaseName: database}
}

// openClusterDriver is openDriver as a cluster client, for the calls that ask
// every node.
func openClusterDriver(t *testing.T, database string) *Driver {
	t.Helper()
	c := testcontainer.OpenCluster(t, testcontainer.Login{})
	return &Driver{client: c, databaseName: database}
}

func TestEndpoints(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    string
		want    []string
		wantErr bool
	}{
		{name: "host and port", host: "redis.example.com", port: "6379", want: []string{"redis.example.com:6379"}},
		{name: "no port", host: "redis.example.com", want: []string{"redis.example.com:6379"}},
		{name: "a host that carries its port", host: "redis.example.com:6380", want: []string{"redis.example.com:6380"}},
		{name: "a port that is not a number", host: "redis.example.com:notaport", port: "6379", want: []string{"redis.example.com:notaport"}},
		{name: "an IPv6 literal", host: "::1", port: "6379", want: []string{"[::1]:6379"}},
		{name: "a bracketed IPv6 literal", host: "[::1]", port: "6379", want: []string{"[::1]:6379"}},
		{
			name: "a list, each with its own port",
			host: "node1:6379,node2:6380,node3",
			port: "7000",
			want: []string{"node1:6379", "node2:6380", "node3:7000"},
		},
		{name: "an empty host", host: " ", wantErr: true},
		{name: "an empty host in a list", host: "node1,,node2", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := endpoints(tc.host, tc.port, false)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestDatabaseIndex(t *testing.T) {
	tests := []struct {
		database string
		want     int
		wantErr  bool
	}{
		{database: "", want: 0},
		{database: "0", want: 0},
		{database: "3", want: 3},
		{database: "fifteen", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.database, func(t *testing.T) {
			got, err := databaseIndex(tc.database)
			if tc.wantErr {
				require.ErrorContains(t, err, "failed to convert database "+tc.database+" to int")
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

// TestEndpointsWithoutAHostWithADialer holds the one address a dialer that names
// its own endpoint is handed: the port alone, which the dialer ignores.
func TestEndpointsWithoutAHostWithADialer(t *testing.T) {
	got, err := endpoints("", "6380", true)
	require.NoError(t, err)
	require.Equal(t, []string{":6380"}, got)
	got, err = endpoints("", "", true)
	require.NoError(t, err)
	require.Equal(t, []string{":6379"}, got)
	_, err = endpoints("", "6380", false)
	require.ErrorContains(t, err, "host must be set")
}

func TestOpen(t *testing.T) {
	tests := []struct {
		name             string
		cfg              driver.Config
		opts             Options
		wantDatabaseName string
		wantErr          string
	}{
		{
			name:             "one server",
			cfg:              driver.Config{Host: "redis.example.com", Port: "6379"},
			wantDatabaseName: "0",
		},
		{
			name:             "a database index",
			cfg:              driver.Config{Host: "redis.example.com", Port: "6379", Database: "3"},
			wantDatabaseName: "3",
		},
		{
			name:    "a database that is not a number",
			cfg:     driver.Config{Host: "redis.example.com", Port: "6379", Database: "cache"},
			wantErr: "failed to convert database cache to int",
		},
		{
			name:    "no host",
			cfg:     driver.Config{Port: "6379"},
			wantErr: "host must be set",
		},
		{
			name:    "a standalone server with a list of hosts",
			cfg:     driver.Config{Host: "node1,node2", Port: "6379"},
			wantErr: "a standalone Redis takes one host, but 2 were given",
		},
		{
			name:    "a topology Redis has no client for",
			cfg:     driver.Config{Host: "redis.example.com", Port: "6379"},
			opts:    Options{Type: "REPLICATED"},
			wantErr: "unsupported redis type REPLICATED",
		},
		{
			name:             "a sentinel",
			cfg:              driver.Config{Host: "sentinel1,sentinel2", Port: "26379", User: "master-user", Password: "master-pass"},
			opts:             Options{Type: Sentinel, MasterName: "mymaster", SentinelUser: "sentinel-user", SentinelPassword: "sentinel-pass"},
			wantDatabaseName: "0",
		},
		{
			name:             "a cluster",
			cfg:              driver.Config{Host: "node1:6379,node2:6379"},
			opts:             Options{Type: Cluster},
			wantDatabaseName: "0",
		},
		{
			name:    "a cluster with a database index it does not have",
			cfg:     driver.Config{Host: "node1:6379,node2:6379", Database: "3"},
			opts:    Options{Type: Cluster},
			wantErr: "a Redis Cluster has database 0 only, but 3 was given",
		},
		{
			name: "a cluster with no address",
			cfg: driver.Config{
				Dial: func(context.Context, string, string) (net.Conn, error) { return nil, errors.New("no connection") },
			},
			opts:    Options{Type: Cluster},
			wantErr: "a Redis Cluster needs its nodes' addresses",
		},
		{
			name: "a sentinel deployment with no address",
			cfg: driver.Config{
				Dial: func(context.Context, string, string) (net.Conn, error) { return nil, errors.New("no connection") },
			},
			opts:    Options{Type: Sentinel, MasterName: "mymaster"},
			wantErr: "a sentinel deployment needs the sentinels' addresses",
		},
		{
			name: "a dialer that names its own endpoint",
			cfg: driver.Config{
				Database: "0",
				Dial:     func(context.Context, string, string) (net.Conn, error) { return nil, errors.New("no connection") },
			},
			wantDatabaseName: "0",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d, err := Open(context.Background(), tc.cfg, tc.opts)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantDatabaseName, d.databaseName)
			require.NoError(t, d.Close())
		})
	}
}

func TestPing(t *testing.T) {
	d := openDriver(t, testcontainer.Login{}, "0")
	require.NoError(t, d.Ping(context.Background()))
}

// TestDialerAsksForEachHostsSettings holds the settings to the host being dialed:
// one client-wide *tls.Config could not verify a cluster's nodes against their
// own names.
func TestDialerAsksForEachHostsSettings(t *testing.T) {
	var hosts []string
	d := dialer(driver.Config{TLS: func(host string) (*tls.Config, error) {
		hosts = append(hosts, host)
		return nil, errors.New("no TLS here")
	}})
	require.NotNil(t, d)
	for _, address := range []string{"node1.example.com:6379", "[2001:db8::1]:6379"} {
		_, err := d(context.Background(), "tcp", address)
		require.ErrorContains(t, err, "no TLS here")
	}
	require.Equal(t, []string{"node1.example.com", "2001:db8::1"}, hosts)
}

// TestDialerRefusesANilConfig holds a caller that answered no settings to its own
// host: TLS was required, so no connection is made at all.
func TestDialerRefusesANilConfig(t *testing.T) {
	dialed := false
	d := dialer(driver.Config{
		Dial: func(context.Context, string, string) (net.Conn, error) {
			dialed = true
			return nil, errors.New("unreachable")
		},
		TLS: func(string) (*tls.Config, error) { return nil, nil },
	})
	_, err := d(context.Background(), "tcp", "127.0.0.1:6379")
	require.ErrorContains(t, err, "the TLS config for 127.0.0.1 is nil")
	require.False(t, dialed, "the dialer reached for a host whose settings were refused")
}

// TestDialerNegotiatesTLSOverTheDialedConnection holds the connection the dialer
// hands the client: TLS over what Config.Dial returned, verified with the settings
// the caller answered for that host. Without this, a tunnel would carry the Redis
// protocol in plaintext although TLS was required.
func TestDialerNegotiatesTLSOverTheDialedConnection(t *testing.T) {
	certificate := selfSignedCertificate(t)
	listener, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		Certificates: []tls.Certificate{certificate},
		MinVersion:   tls.VersionTLS12,
	})
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
	tunneled := false
	d := dialer(driver.Config{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			tunneled = true
			var dialer net.Dialer
			return dialer.DialContext(ctx, network, address)
		},
		TLS: func(host string) (*tls.Config, error) {
			hosts = append(hosts, host)
			return &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS12}, nil
		},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := d(ctx, "tcp", listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()
	_, ok := conn.(*tls.Conn)
	require.True(t, ok, "the connection the client gets is not TLS")
	require.True(t, tunneled, "the connection did not come from Config.Dial")
	require.Equal(t, []string{"127.0.0.1"}, hosts)
}

// selfSignedCertificate is a certificate for 127.0.0.1 that a test's client can
// trust without a certificate authority.
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
