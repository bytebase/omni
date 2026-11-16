package driver

import (
	"context"
	"crypto/tls"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/dynamodb/internal/testcontainer"
)

// newDriver returns a driver whose client is t's tape: a replay of the pinned
// endpoint's answers or, in a recording, the endpoint itself.
func newDriver(t *testing.T) *Driver {
	t.Helper()
	return &Driver{
		client:       OpenTape(t),
		databaseName: testDatabaseName,
		accountID:    testAccountID,
		region:       testRegion,
	}
}

// The identity the tests' driver reports, which a caller passes next to the
// credential.
const (
	testAccountID    = "123456789012"
	testRegion       = "us-east-1"
	testDatabaseName = testAccountID + "-" + testRegion
)

// TestPing reads one table name, which is the connectivity check DynamoDB
// offers: it has no ping of its own.
func TestPing(t *testing.T) {
	require.NoError(t, newDriver(t).Ping(context.Background()))
}

// TestClose closes a driver, which releases the connections its client opened.
func TestClose(t *testing.T) {
	require.NoError(t, newDriver(t).Close())
}

// TestEndpoint holds the endpoint spelling Config's host and port turn into.
func TestEndpoint(t *testing.T) {
	tests := []struct {
		name string
		cfg  driver.Config
		want string
	}{
		{
			// No host: fall back to the AWS SDK default (real AWS) endpoint.
			name: "no_host_uses_default_aws",
			cfg:  driver.Config{Params: map[string]string{ParamRegion: "ap-northeast-1"}},
			want: "",
		},
		{
			// DynamoDB Local: host + port, no TLS -> http endpoint.
			name: "local_http",
			cfg:  driver.Config{Host: "localhost", Port: "8000"},
			want: "http://localhost:8000",
		},
		{
			name: "host_only_no_port",
			cfg:  driver.Config{Host: "127.0.0.1"},
			want: "http://127.0.0.1",
		},
		{
			// A TLS configuration means TLS is required, so the host is reached
			// over https.
			name: "https_via_tls",
			cfg: driver.Config{
				Host: "dynamo.internal",
				Port: "8000",
				TLS:  func(string) (*tls.Config, error) { return &tls.Config{}, nil },
			},
			want: "https://dynamo.internal:8000",
		},
		{
			// Explicit scheme in host is respected; port appended once.
			name: "scheme_prefixed_host_appends_port",
			cfg:  driver.Config{Host: "http://localhost", Port: "8000"},
			want: "http://localhost:8000",
		},
		{
			// Explicit scheme + port already in host: used verbatim.
			name: "scheme_prefixed_full_url",
			cfg:  driver.Config{Host: "https://dynamodb.example.com:443"},
			want: "https://dynamodb.example.com:443",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, endpoint(tc.cfg))
		})
	}
}

// TestEndpointHost holds the host a TLS configuration is asked for.
func TestEndpointHost(t *testing.T) {
	require.Equal(t, "dynamodb.example.com", endpointHost("https://dynamodb.example.com:443"))
	require.Equal(t, "localhost", endpointHost("localhost"))
	require.Equal(t, "localhost", endpointHost("http://localhost:8000"))
}

// TestSyncInstance holds the pseudo database name a sync reports, and the two
// facts it cannot report without. The driver has no client: DynamoDB answers no
// account and no database list, so the name is built from what the caller
// passed and nothing is asked of the endpoint.
func TestSyncInstance(t *testing.T) {
	d := &Driver{accountID: testAccountID, region: testRegion}
	instance, err := d.SyncInstance(context.Background())
	require.NoError(t, err)
	require.Empty(t, instance.Version, "DynamoDB reports no version")
	require.Empty(t, instance.Roles)
	require.Len(t, instance.Databases, 1)
	require.Equal(t, testDatabaseName, instance.Databases[0].Name)
	require.Empty(t, instance.Databases[0].Schemas, "an instance sync reports a database without its tables")
}

func TestSyncInstanceWithoutAnAccountID(t *testing.T) {
	d := &Driver{region: testRegion}
	_, err := d.SyncInstance(context.Background())
	require.EqualError(t, err, "account id is empty in the driver config")
}

func TestSyncInstanceWithoutARegion(t *testing.T) {
	d := &Driver{accountID: testAccountID}
	_, err := d.SyncInstance(context.Background())
	require.EqualError(t, err, "region is empty in the AWS config")
}

// setHermeticAWSEnv shields a test from the host's AWS configuration, so that
// Open's own credential and endpoint resolution reads nothing ambient.
func setHermeticAWSEnv(t *testing.T) {
	t.Helper()
	for _, name := range []string{
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN",
		"AWS_ENDPOINT_URL", "AWS_ENDPOINT_URL_DYNAMODB", "AWS_REGION", "AWS_DEFAULT_REGION",
		"AWS_PROFILE", "AWS_ROLE_ARN",
	} {
		t.Setenv(name, "")
	}
	missing := filepath.Join(t.TempDir(), "nonexistent")
	t.Setenv("AWS_CONFIG_FILE", missing)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", missing)
}

// TestOpenAndClose opens a driver against a region without reaching it, which
// is what Open promises: the client is built, the endpoint is not asked.
func TestOpenAndClose(t *testing.T) {
	setHermeticAWSEnv(t)
	d, err := Open(context.Background(), driver.Config{
		Params:   map[string]string{ParamRegion: testRegion, ParamAccountID: testAccountID},
		User:     testcontainer.AccessKeyID,
		Password: testcontainer.SecretAccessKey,
	})
	require.NoError(t, err)
	require.Equal(t, testAccountID, d.accountID)
	require.Equal(t, testRegion, d.region)
	require.NoError(t, d.Close())
}
