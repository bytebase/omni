package driver

import (
	"context"
	"crypto/tls"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/driver"
)

func baseConfig() driver.Config {
	return driver.Config{User: "dba", Host: "db.example.com", Port: "3306"}
}

func TestConnConfigDefaults(t *testing.T) {
	mc, err := connConfig(baseConfig())
	if err != nil {
		t.Fatal(err)
	}
	if !mc.MultiStatements || mc.MaxAllowedPacket != 0 {
		t.Errorf("MultiStatements=%v MaxAllowedPacket=%d, want true and 0", mc.MultiStatements, mc.MaxAllowedPacket)
	}
	if mc.Net != "tcp" || mc.Addr != "db.example.com:3306" || mc.User != "dba" {
		t.Errorf("Net=%q Addr=%q User=%q", mc.Net, mc.Addr, mc.User)
	}
}

func TestConnConfigSettingsAndUnixSocket(t *testing.T) {
	cfg := driver.Config{
		User:     "app?allowAllFiles=true",
		Password: "p@ss/word?tls=false&allowAllFiles=true",
		Host:     "/var/run/mysqld/mysqld.sock",
		Database: "sales",
		Params: map[string]string{
			"timeout":                 "10s",
			"readTimeout":             "30s",
			"allowCleartextPasswords": "true",
		},
	}
	mc, err := connConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if mc.Net != "unix" || mc.Addr != "/var/run/mysqld/mysqld.sock" {
		t.Errorf("Net=%q Addr=%q", mc.Net, mc.Addr)
	}
	// Credentials are set directly on gomysql.Config and never parsed as DSN syntax.
	if mc.User != cfg.User || mc.Passwd != cfg.Password || mc.AllowAllFiles {
		t.Errorf("User=%q Passwd=%q AllowAllFiles=%v", mc.User, mc.Passwd, mc.AllowAllFiles)
	}
	if mc.DBName != "sales" || mc.Timeout != 10*time.Second || mc.ReadTimeout != 30*time.Second || !mc.AllowCleartextPasswords {
		t.Errorf("DBName=%q Timeout=%v ReadTimeout=%v Cleartext=%v", mc.DBName, mc.Timeout, mc.ReadTimeout, mc.AllowCleartextPasswords)
	}
}

func TestConnConfigRejectsAllowAllFiles(t *testing.T) {
	for _, key := range []string{"allowAllFiles", "allowallfiles", "AllowAllFiles", "  allowAllFiles  "} {
		cfg := baseConfig()
		cfg.Params = map[string]string{key: "true"}
		if _, err := connConfig(cfg); err == nil || !strings.Contains(strings.ToLower(err.Error()), "allowallfiles") {
			t.Errorf("param %q: err = %v, want allowAllFiles rejection", key, err)
		}
	}
}

func TestConnConfigTLSCannotBeDowngraded(t *testing.T) {
	cfg := baseConfig()
	cfg.Params = map[string]string{
		"tls":                      "false",
		"allowFallbackToPlaintext": "true",
	}
	var gotHost string
	cfg.TLS = func(host string) (*tls.Config, error) {
		gotHost = host
		return &tls.Config{ServerName: host}, nil
	}
	mc, err := connConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if gotHost != "db.example.com" || mc.TLS == nil || mc.TLS.ServerName != "db.example.com" || mc.AllowFallbackToPlaintext {
		t.Fatalf("gotHost=%q TLS=%v AllowFallbackToPlaintext=%v", gotHost, mc.TLS, mc.AllowFallbackToPlaintext)
	}
}

func TestConnConfigDialWithoutHost(t *testing.T) {
	var dialed bool
	cfg := driver.Config{
		User: "iam-user",
		Dial: func(context.Context, string, string) (net.Conn, error) {
			dialed = true
			return nil, net.ErrClosed
		},
	}
	mc, err := connConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if mc.DialFunc == nil {
		t.Fatal("DialFunc was not set")
	}
	_, _ = mc.DialFunc(context.Background(), "tcp", mc.Addr)
	if !dialed {
		t.Error("DialFunc was not invoked")
	}
}

func TestConnConfigValidation(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*driver.Config)
		want string
	}{
		{"blank user", func(c *driver.Config) { c.User = "" }, "user must be set"},
		{"blank host", func(c *driver.Config) { c.Host = "  " }, "host must be set"},
		{"blank port", func(c *driver.Config) { c.Port = "  " }, "port must be set"},
	} {
		cfg := baseConfig()
		tc.edit(&cfg)
		if _, err := connConfig(cfg); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: err = %v, want %q", tc.name, err, tc.want)
		}
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

func TestParseVersion(t *testing.T) {
	for _, tc := range []struct {
		in, wantVer, wantRest string
	}{
		{"8.0.27", "8.0.27", ""},
		{"5.7.22-log", "5.7.22", "-log"},
		{"5.6.29_ddm_3.0.1.7", "5.6.29", "_ddm_3.0.1.7"},
		{"10.4.7-MariaDB", "10.4.7", "-MariaDB"},
	} {
		ver, rest, err := parseVersion(tc.in)
		if err != nil || ver != tc.wantVer || rest != tc.wantRest {
			t.Errorf("parseVersion(%q) = (%q, %q, %v)", tc.in, ver, rest, err)
		}
	}
}
