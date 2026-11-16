package driver

import (
	"context"
	"crypto/tls"
	"net"
	"strings"
	"testing"

	"github.com/microsoft/go-mssqldb/msdsn"

	"github.com/bytebase/omni/driver"
)

func baseConfig() driver.Config {
	return driver.Config{User: "sa", Host: "db.example.com", Port: "1433"}
}

func TestMSSQLConfigDefaults(t *testing.T) {
	config, err := (&Driver{}).mssqlConfig(baseConfig())
	if err != nil {
		t.Fatal(err)
	}
	if config.Host != "db.example.com" || config.Port != 1433 {
		t.Errorf("host %q port %d", config.Host, config.Port)
	}
	if config.User != "sa" || config.Database != "" || config.AppName != "bytebase" {
		t.Errorf("user %q database %q app name %q", config.User, config.Database, config.AppName)
	}
	if config.Encryption != msdsn.EncryptionOff {
		t.Errorf("encryption %d, want off", config.Encryption)
	}
	if got := config.Parameters[msdsn.TrustServerCertificate]; got != "true" {
		t.Errorf("trustservercertificate %q, want the client to accept the server's certificate", got)
	}
	if got := config.Parameters[msdsn.TLSMin]; got != "1.0" {
		t.Errorf("tlsmin %q", got)
	}
}

func TestMSSQLConfigSettings(t *testing.T) {
	cfg := baseConfig()
	cfg.Password = "p@ss/word?x=y&database=other"
	cfg.Database = "sales"
	cfg.ApplicationName = "bytebase-taskrun-7"
	cfg.Params = map[string]string{"Encrypt": "strict", "packet size": "8192"}
	config, err := (&Driver{}).mssqlConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	// The credentials never reach the DSN, so a separator in one names nothing.
	if config.User != "sa" || config.Password != cfg.Password || config.Database != "sales" {
		t.Errorf("user %q password %q database %q", config.User, config.Password, config.Database)
	}
	if config.AppName != "bytebase-taskrun-7" {
		t.Errorf("app name %q", config.AppName)
	}
	if config.Encryption != msdsn.EncryptionStrict {
		t.Errorf("encryption %d, want the Params' strict", config.Encryption)
	}
	if got := config.Parameters["packet size"]; got != "8192" {
		t.Errorf("packet size %q", got)
	}
}

func TestMSSQLConfigInstanceAndIPv6(t *testing.T) {
	cfg := baseConfig()
	cfg.Host = `db.example.com\SQLEXPRESS`
	config, err := (&Driver{}).mssqlConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if config.Host != "db.example.com" || config.Instance != "SQLEXPRESS" {
		t.Errorf("host %q instance %q", config.Host, config.Instance)
	}

	ipv6 := baseConfig()
	ipv6.Host = "[2001:db8::1]"
	ipv6.Port = ""
	config, err = (&Driver{}).mssqlConfig(ipv6)
	if err != nil {
		t.Fatal(err)
	}
	if config.Host != "2001:db8::1" || config.Port != 1433 {
		t.Errorf("host %q port %d", config.Host, config.Port)
	}
}

func TestMSSQLConfigTLS(t *testing.T) {
	cfg := baseConfig()
	cfg.Params = map[string]string{"encrypt": "false"}
	var asked string
	cfg.TLS = func(host string) (*tls.Config, error) {
		asked = host
		return &tls.Config{ServerName: host}, nil
	}
	config, err := (&Driver{}).mssqlConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if asked != "db.example.com" || config.TLSConfig == nil || config.TLSConfig.ServerName != "db.example.com" {
		t.Fatalf("TLS asked for %q, config %+v", asked, config.TLSConfig)
	}
	// A caller's TLS cannot be turned off by Params.
	if config.Encryption != msdsn.EncryptionRequired {
		t.Errorf("encryption %d, want required", config.Encryption)
	}
}

func TestMSSQLConfigValidation(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*driver.Config)
		want string
	}{
		{"blank user", func(c *driver.Config) { c.User = "" }, "user must be set"},
		{"blank host", func(c *driver.Config) { c.Host = "  " }, "host must be set"},
		{"bad port", func(c *driver.Config) { c.Port = "5432 host=other.example.com" }, "invalid port"},
	} {
		cfg := baseConfig()
		tc.edit(&cfg)
		if _, err := (&Driver{}).mssqlConfig(cfg); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: err = %v, want %q", tc.name, err, tc.want)
		}
	}
}

func TestMSSQLConfigDial(t *testing.T) {
	cfg := baseConfig()
	cfg.Dial = func(context.Context, string, string) (net.Conn, error) { return nil, net.ErrClosed }
	d := &Driver{}
	config, err := d.mssqlConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Protocols) != 1 || config.Protocols[0] != d.dialProtocol {
		t.Fatalf("protocols %q, want the registered dialer", config.Protocols)
	}
	if _, ok := msdsn.ProtocolDialers[d.dialProtocol]; !ok {
		t.Fatal("the dialer is not registered")
	}
	d.closeDialer()
	if _, ok := msdsn.ProtocolDialers[d.dialProtocol]; ok {
		t.Fatal("Close kept the dialer registered")
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
