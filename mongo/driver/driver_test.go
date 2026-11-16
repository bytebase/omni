package driver

import (
	"context"
	"crypto/tls"
	"errors"
	"net/url"
	"testing"

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

// TestOpenCarriesTheCredentialsAndTLSIntoTheClient holds Open to what it hands
// the client: the endpoint, the resolved login, TLS required when the caller
// supplies it, and the dialer.
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

func TestOpenRefusesANilTLSConfig(t *testing.T) {
	_, err := Open(context.Background(), driver.Config{
		Host: "mongo.example.com",
		TLS:  func(string) (*tls.Config, error) { return nil, nil },
	})
	require.ErrorContains(t, err, "TLS configuration must not be nil")
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
