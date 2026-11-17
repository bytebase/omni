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
	"testing"

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
