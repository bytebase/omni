package driver

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/bytebase/omni/driver"
)

func TestBuildURLDefaults(t *testing.T) {
	u, err := buildURL(driver.Config{Host: "trino.example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := u.String(), "http://trino@trino.example.com:8080?catalog=system&source=bytebase"; got != want {
		t.Errorf("buildURL = %q, want %q", got, want)
	}
}

func TestBuildURLNamesOrOverridesEverySetting(t *testing.T) {
	u, err := buildURL(driver.Config{
		Host:            "[2001:db8::1]",
		Port:            "8443",
		User:            "app",
		Password:        "p@ss/word?x=1",
		Database:        "memory",
		ApplicationName: "bytebase-42",
		TLS:             func(string) (*tls.Config, error) { return &tls.Config{}, nil },
		Params:          map[string]string{"source": "replaced", "clientTags": "etl"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if u.Scheme != "https" {
		t.Errorf("Scheme = %q, want https for a TLS config", u.Scheme)
	}
	if u.Host != "[2001:db8::1]:8443" {
		t.Errorf("Host = %q, want the IPv6 host in brackets", u.Host)
	}
	if u.User.Username() != "app" {
		t.Errorf("User = %q", u.User.Username())
	}
	if password, _ := u.User.Password(); password != "p@ss/word?x=1" {
		t.Errorf("Password = %q, want it carried whole", password)
	}
	query := u.Query()
	if query.Get("catalog") != "memory" {
		t.Errorf("catalog = %q", query.Get("catalog"))
	}
	if query.Get("source") != "replaced" {
		t.Errorf("source = %q, want the parameter to override the default", query.Get("source"))
	}
	if query.Get("clientTags") != "etl" {
		t.Errorf("clientTags = %q", query.Get("clientTags"))
	}
}

func TestBuildURLRefusesAPasswordWithoutTLS(t *testing.T) {
	_, err := buildURL(driver.Config{Host: "trino.example.com", Password: "secret"})
	if err == nil || !strings.Contains(err.Error(), "requires SSL to authenticate with a password") {
		t.Fatalf("err = %v, want the SSL requirement", err)
	}
	if strings.Contains(err.Error(), "secret") {
		t.Errorf("the error must not carry the password: %v", err)
	}
}

// TestOpenRegistersAndDropsTheCustomClient pins the lifecycle of the process-wide
// HTTP client registry the client library keeps a TLS or dial configuration in:
// an open driver owns one entry, and Close removes it.
func TestOpenRegistersAndDropsTheCustomClient(t *testing.T) {
	for _, tc := range []struct {
		name string
		cfg  func(host string) driver.Config
	}{
		{"TLS", func(host string) driver.Config {
			return driver.Config{Host: host, TLS: func(string) (*tls.Config, error) { return &tls.Config{}, nil }}
		}},
		{"dialer", func(host string) driver.Config {
			return driver.Config{Host: host, Dial: func(context.Context, string, string) (net.Conn, error) { return nil, errors.New("unused") }}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, err := Open(context.Background(), tc.cfg("trino.example.com"))
			if err != nil {
				t.Fatal(err)
			}
			if d.tlsClientKey == "" {
				t.Fatal("Open registered no custom client")
			}
			key := d.tlsClientKey
			if err := d.Close(); err != nil {
				t.Fatal(err)
			}
			if d.tlsClientKey != "" || d.httpClient != nil {
				t.Errorf("Close left the custom client registered")
			}

			second, err := Open(context.Background(), tc.cfg("trino.example.com"))
			if err != nil {
				t.Fatal(err)
			}
			defer second.Close()
			if second.tlsClientKey == key {
				t.Errorf("two drivers share the custom client key %q", key)
			}
		})
	}
}

func TestOpenWithoutTLSRegistersNoClient(t *testing.T) {
	d, err := Open(context.Background(), driver.Config{Host: "trino.example.com"})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	if d.tlsClientKey != "" {
		t.Errorf("Open registered a custom client without TLS or a dialer")
	}
}

func TestOpenDefaultsTheCatalog(t *testing.T) {
	for _, tc := range []struct{ database, want string }{
		{"", "system"},
		{"memory", "memory"},
	} {
		d, err := Open(context.Background(), driver.Config{Host: "trino.example.com", Database: tc.database})
		if err != nil {
			t.Fatal(err)
		}
		if d.databaseName != tc.want {
			t.Errorf("databaseName = %q, want %q", d.databaseName, tc.want)
		}
		if err := d.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestOpenRefusesNilTLSConfig(t *testing.T) {
	_, err := Open(context.Background(), driver.Config{
		Host: "trino.example.com",
		TLS:  func(string) (*tls.Config, error) { return nil, nil },
	})
	if err == nil || !strings.Contains(err.Error(), "TLS configuration must not be nil") {
		t.Fatalf("err = %v, want the nil TLS refusal", err)
	}
}

func TestOpenWrapsTLSConfigError(t *testing.T) {
	_, err := Open(context.Background(), driver.Config{
		Host: "trino.example.com",
		TLS:  func(string) (*tls.Config, error) { return nil, errors.New("no certificate") },
	})
	if err == nil || !strings.Contains(err.Error(), "failed to get TLS config: no certificate") {
		t.Fatalf("err = %v, want the TLS config error", err)
	}
}
