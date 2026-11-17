package driver

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/snowflakedb/gosnowflake/v2"
	"github.com/youmark/pkcs8"

	"github.com/bytebase/omni/driver"
)

// TestSnowConfigBuildsTheConnection covers the settings Open hands gosnowflake:
// which ones are required, which have a default, and how a caller's credential,
// TLS, Dial and Params reach the client. None of it needs an account.
func TestSnowConfigBuildsTheConnection(t *testing.T) {
	for _, tc := range []struct {
		name string
		cfg  driver.Config
		want func(t *testing.T, got *gosnowflake.Config)
	}{
		{
			name: "the account, user and password",
			cfg:  driver.Config{User: "dba", Password: "s3cret", Host: "xy12345"},
			want: func(t *testing.T, got *gosnowflake.Config) {
				if got.Account != "xy12345" {
					t.Errorf("Account = %q, want xy12345", got.Account)
				}
				if got.User != "dba" {
					t.Errorf("User = %q, want dba", got.User)
				}
				if got.Password != "s3cret" {
					t.Errorf("Password = %q", got.Password)
				}
				if got.Authenticator != gosnowflake.AuthTypeSnowflake {
					t.Errorf("Authenticator = %v, want a password login", got.Authenticator)
				}
				// The client derives the endpoint from the account itself.
				if got.Host != "" {
					t.Errorf("Host = %q, want the client's own derivation", got.Host)
				}
				if got.Port != 0 {
					t.Errorf("Port = %d, want the client's default of 443", got.Port)
				}
				if got.Database != "" {
					t.Errorf("Database = %q, want none", got.Database)
				}
			},
		},
		{
			name: "a database is quoted, so a lower-case name is the one the caller meant",
			cfg:  driver.Config{User: "dba", Host: "xy12345", Database: "analytics"},
			want: func(t *testing.T, got *gosnowflake.Config) {
				if got.Database != `"analytics"` {
					t.Errorf("Database = %q, want %q", got.Database, `"analytics"`)
				}
			},
		},
		{
			name: "a database name holding a quote stays one identifier",
			cfg:  driver.Config{User: "dba", Host: "xy12345", Database: `a"b`},
			want: func(t *testing.T, got *gosnowflake.Config) {
				if got.Database != `"a""b"` {
					t.Errorf("Database = %q, want %q", got.Database, `"a""b"`)
				}
			},
		},
		{
			name: "the port and the session's region, warehouse and role come from the caller",
			cfg: driver.Config{
				User: "dba",
				Host: "xy12345",
				Port: "8443",
				Params: map[string]string{
					ParamRegion:    "eu-central-1",
					ParamWarehouse: "wh",
					ParamRole:      "analyst",
				},
			},
			want: func(t *testing.T, got *gosnowflake.Config) {
				if got.Port != 8443 {
					t.Errorf("Port = %d, want 8443", got.Port)
				}
				if got.Region != "eu-central-1" {
					t.Errorf("Region = %q, want eu-central-1", got.Region)
				}
				if got.Warehouse != "wh" {
					t.Errorf("Warehouse = %q, want wh", got.Warehouse)
				}
				if got.Role != "analyst" {
					t.Errorf("Role = %q, want analyst", got.Role)
				}
			},
		},
		{
			name: "a host after the account is the endpoint to dial",
			cfg:  driver.Config{User: "dba", Host: "xy12345@10.0.0.1", Port: "4182"},
			want: func(t *testing.T, got *gosnowflake.Config) {
				if got.Account != "xy12345" {
					t.Errorf("Account = %q, want xy12345", got.Account)
				}
				if got.Host != "10.0.0.1" {
					t.Errorf("Host = %q, want 10.0.0.1", got.Host)
				}
				if got.Port != 4182 {
					t.Errorf("Port = %d, want 4182", got.Port)
				}
			},
		},
		{
			name: "the application name and the client's own parameters reach the session",
			cfg: driver.Config{
				User:            "dba",
				Host:            "xy12345",
				ApplicationName: "bytebase",
				Params:          map[string]string{"timestamp_output_format": "YYYY-MM-DD HH24:MI:SS"},
			},
			want: func(t *testing.T, got *gosnowflake.Config) {
				if got.Application != "bytebase" {
					t.Errorf("Application = %q, want bytebase", got.Application)
				}
				format, ok := got.Params["timestamp_output_format"]
				if !ok || format == nil {
					t.Fatalf("Params = %v, want it to hold timestamp_output_format", got.Params)
				}
				if *format != "YYYY-MM-DD HH24:MI:SS" {
					t.Errorf("timestamp_output_format = %q", *format)
				}
			},
		},
		{
			name: "the parameters this driver reads are not passed on",
			cfg: driver.Config{
				User: "dba",
				Host: "xy12345",
				Params: map[string]string{
					ParamRegion:    "eu-central-1",
					ParamWarehouse: "wh",
					ParamRole:      "analyst",
				},
			},
			want: func(t *testing.T, got *gosnowflake.Config) {
				for _, name := range []string{ParamRegion, ParamWarehouse, ParamRole} {
					if _, ok := got.Params[name]; ok {
						t.Errorf("Params holds %q, which the driver reads itself", name)
					}
				}
			},
		},
		{
			name: "OCSP fails open, as it did before the client turned it off",
			cfg:  driver.Config{User: "dba", Host: "xy12345"},
			want: func(t *testing.T, got *gosnowflake.Config) {
				if got.OCSPFailOpen != gosnowflake.OCSPFailOpenTrue {
					t.Errorf("OCSPFailOpen = %v, want fail-open", got.OCSPFailOpen)
				}
			},
		},
		{
			name: "no transport is built without a dialer or TLS",
			cfg:  driver.Config{User: "dba", Host: "xy12345"},
			want: func(t *testing.T, got *gosnowflake.Config) {
				if got.Transporter != nil {
					t.Errorf("Transporter = %T, want the client's own", got.Transporter)
				}
			},
		},
		{
			name: "a dialer is the transport's",
			cfg: driver.Config{
				User: "dba",
				Host: "xy12345",
				Dial: func(context.Context, string, string) (net.Conn, error) { return nil, nil },
			},
			want: func(t *testing.T, got *gosnowflake.Config) {
				transport, ok := got.Transporter.(*http.Transport)
				if !ok {
					t.Fatalf("Transporter = %T, want an *http.Transport", got.Transporter)
				}
				if transport.DialContext == nil {
					t.Error("DialContext is not the caller's dialer")
				}
				// The transport keeps the default's own TLS settings, which verify
				// the server's certificate, when the caller supplied none.
				if transport.TLSClientConfig != nil && transport.TLSClientConfig.ServerName != "" {
					t.Errorf("TLSClientConfig.ServerName = %q, want the default's", transport.TLSClientConfig.ServerName)
				}
			},
		},
		{
			name: "TLS is asked for the endpoint the driver dials",
			cfg: driver.Config{
				User: "dba",
				Host: "xy12345",
				TLS: func(host string) (*tls.Config, error) {
					return &tls.Config{ServerName: host}, nil
				},
			},
			want: func(t *testing.T, got *gosnowflake.Config) {
				transport, ok := got.Transporter.(*http.Transport)
				if !ok {
					t.Fatalf("Transporter = %T, want an *http.Transport", got.Transporter)
				}
				if transport.TLSClientConfig == nil || transport.TLSClientConfig.ServerName != "xy12345.snowflakecomputing.com" {
					t.Errorf("TLSClientConfig = %v, want the account's host", transport.TLSClientConfig)
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := snowConfig(tc.cfg)
			if err != nil {
				t.Fatalf("snowConfig: %v", err)
			}
			tc.want(t, got)
		})
	}
}

// TestEndpointHost pins the host a caller's TLS settings are asked for, which is
// the one the client dials for an account.
func TestEndpointHost(t *testing.T) {
	for _, tc := range []struct {
		account, region string
		want            string
	}{
		{account: "xy12345", want: "xy12345.snowflakecomputing.com"},
		{account: "nb47110.ap-southeast-1", want: "nb47110.ap-southeast-1.snowflakecomputing.com"},
		{account: "xy12345", region: "eu-central-1", want: "xy12345.eu-central-1.snowflakecomputing.com"},
		{account: "xy12345", region: "cn-north-1", want: "xy12345.cn-north-1.snowflakecomputing.cn"},
		{account: "org-account", want: "org-account.snowflakecomputing.com"},
	} {
		t.Run(tc.account+"/"+tc.region, func(t *testing.T) {
			if got := endpointHost(tc.account, tc.region); got != tc.want {
				t.Errorf("endpointHost(%q, %q) = %q, want %q", tc.account, tc.region, got, tc.want)
			}
		})
	}
}

// TestSnowConfigRefusesWhatCannotNameAnAccount covers the values Open refuses
// rather than letting the client fail later with a different message.
func TestSnowConfigRefusesWhatCannotNameAnAccount(t *testing.T) {
	for _, tc := range []struct {
		name    string
		cfg     driver.Config
		wantErr string
	}{
		{
			name:    "no user",
			cfg:     driver.Config{Host: "xy12345", Password: "s3cret"},
			wantErr: "user must be set",
		},
		{
			name:    "no host",
			cfg:     driver.Config{User: "dba", Password: "s3cret"},
			wantErr: "host must be set",
		},
		{
			name:    "two @s",
			cfg:     driver.Config{User: "dba", Host: "a@b@c"},
			wantErr: "expected one @ in the host at most",
		},
		{
			name:    "an account without its host",
			cfg:     driver.Config{User: "dba", Host: "xy12345@"},
			wantErr: "must name both an account and a host",
		},
		{
			name:    "a host without its account",
			cfg:     driver.Config{User: "dba", Host: "@10.0.0.1"},
			wantErr: "must name both an account and a host",
		},
		{
			name:    "a port that is not a number",
			cfg:     driver.Config{User: "dba", Host: "xy12345", Port: "https"},
			wantErr: "port \"https\" is not a number",
		},
		{
			name: "an encrypted key without its passphrase",
			cfg: driver.Config{
				User:   "dba",
				Host:   "xy12345",
				Params: map[string]string{ParamPrivateKey: encryptedKey(t)},
			},
			wantErr: "private key is encrypted, but no passphrase was provided",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := snowConfig(tc.cfg)
			if err == nil {
				t.Fatal("snowConfig accepted the configuration")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("error = %q, want it to hold %q", err, tc.wantErr)
			}
		})
	}
}

// TestSnowConfigReadsAKeyPairLogin covers the key-pair credential: a private key
// in Params makes the login a JWT one and the password unused.
func TestSnowConfigReadsAKeyPairLogin(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	config, err := snowConfig(driver.Config{
		User:     "dba",
		Password: "unused",
		Host:     "xy12345",
		Params:   map[string]string{ParamPrivateKey: pemBlock("PRIVATE KEY", der)},
	})
	if err != nil {
		t.Fatalf("snowConfig: %v", err)
	}
	if config.Authenticator != gosnowflake.AuthTypeJwt {
		t.Errorf("Authenticator = %v, want a JWT login", config.Authenticator)
	}
	if config.PrivateKey == nil || !config.PrivateKey.Equal(key) {
		t.Error("PrivateKey is not the caller's key")
	}
	if config.Password != "" {
		t.Errorf("Password = %q, want it unused beside a private key", config.Password)
	}
}

// TestDecodeRSAPrivateKey covers the two PEM forms the login accepts and the
// ones it refuses.
func TestDecodeRSAPrivateKey(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	unencrypted, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	encrypted, err := pkcs8.MarshalPrivateKey(key, []byte("passphrase"), nil)
	if err != nil {
		t.Fatalf("encrypt key: %v", err)
	}
	_, edKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate ed25519 key: %v", err)
	}
	notRSA, err := x509.MarshalPKCS8PrivateKey(edKey)
	if err != nil {
		t.Fatalf("marshal ed25519 key: %v", err)
	}

	for _, tc := range []struct {
		name       string
		key        string
		passphrase string
		wantErr    string
	}{
		{
			name: "unencrypted key",
			key:  pemBlock("PRIVATE KEY", unencrypted),
		},
		{
			name:       "unencrypted key with an unnecessary passphrase",
			key:        pemBlock("PRIVATE KEY", unencrypted),
			passphrase: "passphrase",
		},
		{
			name:       "encrypted key with its passphrase",
			key:        pemBlock("ENCRYPTED PRIVATE KEY", encrypted),
			passphrase: "passphrase",
		},
		{
			name:    "encrypted key without a passphrase",
			key:     pemBlock("ENCRYPTED PRIVATE KEY", encrypted),
			wantErr: "private key is encrypted, but no passphrase was provided",
		},
		{
			name:       "encrypted key with the wrong passphrase",
			key:        pemBlock("ENCRYPTED PRIVATE KEY", encrypted),
			passphrase: "wrong",
			wantErr:    "incorrect password",
		},
		{
			name:    "key that is not RSA",
			key:     pemBlock("PRIVATE KEY", notRSA),
			wantErr: "expected RSA private key",
		},
		{
			name:    "PKCS#1 block",
			key:     pemBlock("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(key)),
			wantErr: "unsupported pem block type: RSA PRIVATE KEY",
		},
		{
			name:    "not PEM at all",
			key:     "rsa_key.p8",
			wantErr: "failed to get private key PEM block",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeRSAPrivateKey(tc.key, tc.passphrase)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatal("decodeRSAPrivateKey accepted the key")
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("error = %q, want it to hold %q", err, tc.wantErr)
				}
				if got != nil {
					t.Error("a refused key still returned a value")
				}
				return
			}
			if err != nil {
				t.Fatalf("decodeRSAPrivateKey: %v", err)
			}
			if !got.Equal(key) {
				t.Error("decodeRSAPrivateKey returned another key")
			}
		})
	}
}

// TestQuoteIdentAndLiteral pins the two quotings the sync, the query path and
// the dump write: an identifier's and a string literal's.
func TestQuoteIdentAndLiteral(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{in: "PUBLIC", want: `"PUBLIC"`},
		{in: "analytics", want: `"analytics"`},
		{in: `a"b`, want: `"a""b"`},
	} {
		if got := quoteIdent(tc.in); got != tc.want {
			t.Errorf("quoteIdent(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
	for _, tc := range []struct{ in, want string }{
		{in: "STREAM", want: "'STREAM'"},
		{in: `a'b`, want: `'a''b'`},
		// Snowflake reads a backslash as an escape inside a literal, so a name
		// holding one is doubled rather than sent through.
		{in: `a\b`, want: `'a\\b'`},
	} {
		if got := sqlLiteral(tc.in); got != tc.want {
			t.Errorf("sqlLiteral(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// TestOpenAcceptsAConfigurationWithoutReachingTheAccount checks that Open, which
// reaches nothing, keeps the database name the sync reports.
func TestOpenAcceptsAConfigurationWithoutReachingTheAccount(t *testing.T) {
	d, err := Open(context.Background(), driver.Config{
		User:     "dba",
		Password: "s3cret",
		Host:     "xy12345",
		Database: "analytics",
	})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer d.Close()
	if d.databaseName != "analytics" {
		t.Errorf("databaseName = %q, want analytics", d.databaseName)
	}
}

// pemBlock encodes der as one PEM block of blockType.
func pemBlock(blockType string, der []byte) string {
	return string(pem.EncodeToMemory(&pem.Block{Type: blockType, Bytes: der}))
}

// encryptedKey returns a PKCS#8 private key in PEM, encrypted with a passphrase.
func encryptedKey(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	encrypted, err := pkcs8.MarshalPrivateKey(key, []byte("passphrase"), nil)
	if err != nil {
		t.Fatalf("encrypt key: %v", err)
	}
	return pemBlock("ENCRYPTED PRIVATE KEY", encrypted)
}
