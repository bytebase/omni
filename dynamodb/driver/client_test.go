package driver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"hash/crc32"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
)

// endpointConfig is a configuration for a local endpoint whose credential is
// static, so the SDK reads none from the environment while a test reaches a
// server of its own.
func endpointConfig(host string) driver.Config {
	return driver.Config{
		Host:     host,
		User:     "AKIDEXAMPLE",
		Password: "secret",
		Params:   map[string]string{ParamRegion: "us-east-1"},
	}
}

// newTLSServer answers a table listing over a certificate it signed itself.
func newTLSServer(t *testing.T) *httptest.Server {
	t.Helper()
	body := []byte(`{"TableNames":[]}`)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/x-amz-json-1.0")
		w.Header().Set("x-amz-crc32", strconv.FormatUint(uint64(crc32.ChecksumIEEE(body)), 10))
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)
	return server
}

// requireCertificateRejected holds that err is the failure of verifying a
// certificate no root signed, rather than another way of not connecting.
func requireCertificateRejected(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	var unknownAuthority x509.UnknownAuthorityError
	require.True(t, errors.As(err, &unknownAuthority), "want a certificate verification failure, got %v", err)
}

// TestHTTPSWithoutTLSMaterialVerifiesTheCertificate holds that an https
// endpoint reached with no TLS material is verified: nil Config.TLS is the
// SDK's own transport, which trusts the system roots and no other. A caller
// whose previous client always verified, as Bytebase's DynamoDB driver did,
// therefore keeps verifying unless its user asked to skip, rather than skipping
// because a flag is unset.
func TestHTTPSWithoutTLSMaterialVerifiesTheCertificate(t *testing.T) {
	server := newTLSServer(t)

	d, err := Open(context.Background(), endpointConfig(server.URL))
	require.NoError(t, err)
	defer func() { require.NoError(t, d.Close()) }()

	requireCertificateRejected(t, d.Ping(context.Background()))
}

// TestDialWithoutTLSMaterialVerifiesTheCertificate holds that a caller's dialer
// replaces the connection and not the certificate check: the transport the
// driver builds around it keeps the default verification when Config.TLS is
// nil.
func TestDialWithoutTLSMaterialVerifiesTheCertificate(t *testing.T) {
	server := newTLSServer(t)

	client, err := newHTTPClient(driver.Config{Dial: (&net.Dialer{}).DialContext})
	require.NoError(t, err)
	t.Cleanup(client.CloseIdleConnections)

	_, err = client.Get(server.URL)
	requireCertificateRejected(t, err)
}

// TestHTTPSWithoutTLSMaterialIsSkippedOnlyWhenTheCallerSaysSo holds that the
// caller's *tls.Config is used as written: an explicit InsecureSkipVerify
// reaches the endpoint, so the behavior the driver had before the contract
// stays expressible and the driver adds no skip of its own.
func TestHTTPSWithoutTLSMaterialIsSkippedOnlyWhenTheCallerSaysSo(t *testing.T) {
	server := newTLSServer(t)

	cfg := endpointConfig(server.URL)
	cfg.TLS = func(string) (*tls.Config, error) {
		return &tls.Config{InsecureSkipVerify: true}, nil //nolint:gosec // the opt-out is what this test pins.
	}
	d, err := Open(context.Background(), cfg)
	require.NoError(t, err)
	defer func() { require.NoError(t, d.Close()) }()

	require.NoError(t, d.Ping(context.Background()))
}

// newIdentityClient answers STS with handler, so a test reaches the identity
// call without an account and without the network.
func newIdentityClient(t *testing.T, handler http.HandlerFunc) *liveClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return &liveClient{sts: sts.NewFromConfig(aws.Config{
		Region:      "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider("key", "secret", ""),
		HTTPClient:  server.Client(),
	}, func(o *sts.Options) { o.BaseEndpoint = aws.String(server.URL) })}
}

// TestCallerAccount holds the two answers SyncInstance reads from STS: the
// account the credential belongs to, and the service's refusal.
func TestCallerAccount(t *testing.T) {
	t.Run("reads the account", func(t *testing.T) {
		d := newIdentityClient(t, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/xml")
			_, _ = w.Write([]byte(`<GetCallerIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetCallerIdentityResult>
    <Arn>arn:aws:iam::123456789012:user/test</Arn>
    <UserId>AIDAEXAMPLE</UserId>
    <Account>123456789012</Account>
  </GetCallerIdentityResult>
  <ResponseMetadata><RequestId>01234567-89ab-cdef-0123-456789abcdef</RequestId></ResponseMetadata>
</GetCallerIdentityResponse>`))
		})
		account, err := d.CallerAccount(context.Background())
		require.NoError(t, err)
		require.Equal(t, "123456789012", account)
	})

	t.Run("reports the identity failure", func(t *testing.T) {
		d := newIdentityClient(t, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/xml")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`<ErrorResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <Error><Type>Sender</Type><Code>InvalidClientTokenId</Code><Message>The security token included in the request is invalid.</Message></Error>
  <RequestId>01234567-89ab-cdef-0123-456789abcdef</RequestId>
</ErrorResponse>`))
		})
		_, err := d.CallerAccount(context.Background())
		require.Error(t, err)
		require.ErrorContains(t, err, "InvalidClientTokenId")
	})
}

// TestCallerAccountHonorsTheDialer holds that the identity call reaches STS
// through the network path Config.Dial names: a tunnel or connector that
// carries DynamoDB carries the account lookup too.
func TestCallerAccountHonorsTheDialer(t *testing.T) {
	setHermeticAWSEnv(t)
	// The SDK retries a dial error with backoff; one attempt is enough to see
	// that the dialer was used.
	t.Setenv("AWS_MAX_ATTEMPTS", "1")
	var dialed atomic.Int32
	d, err := Open(context.Background(), driver.Config{
		User:     "AKIDEXAMPLE",
		Password: "secret",
		Params:   map[string]string{ParamRegion: "us-east-1"},
		Dial: func(_ context.Context, _, address string) (net.Conn, error) {
			dialed.Add(1)
			return nil, errors.New("dialer reached " + address)
		},
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, d.Close()) }()

	_, err = d.client.CallerAccount(context.Background())
	require.ErrorContains(t, err, "dialer reached sts.")
	require.GreaterOrEqual(t, dialed.Load(), int32(1), "STS did not use Config.Dial")
}
