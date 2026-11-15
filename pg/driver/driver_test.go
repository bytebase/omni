package driver

import (
	"context"
	"crypto/tls"
	"net"
	"slices"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/bytebase/omni/driver"
)

func baseConfig() driver.Config {
	return driver.Config{User: "dba", Host: "db.example.com", Port: "5432"}
}

func TestConnConfigDefaults(t *testing.T) {
	c, err := connConfig(baseConfig())
	if err != nil {
		t.Fatal(err)
	}
	// PgBouncer in transaction mode refuses named prepared statements.
	if c.DefaultQueryExecMode != pgx.QueryExecModeExec || c.StatementCacheCapacity != 0 {
		t.Errorf("query mode %v, statement cache %d; want exec and none", c.DefaultQueryExecMode, c.StatementCacheCapacity)
	}
	if c.Database != "postgres" {
		t.Errorf("database %q, want postgres", c.Database)
	}
	// libpq's default sslmode=prefer: TLS first, then plaintext.
	if c.TLSConfig == nil || len(c.Fallbacks) != 1 || c.Fallbacks[0].TLSConfig != nil {
		t.Errorf("TLS %v with fallbacks %v, want TLS then plaintext", c.TLSConfig, c.Fallbacks)
	}
	if _, ok := c.RuntimeParams["default_transaction_read_only"]; ok {
		t.Errorf("a session that is not read-only sets default_transaction_read_only")
	}
}

func TestConnConfigSettings(t *testing.T) {
	cfg := baseConfig()
	cfg.Database = "sales"
	cfg.Password = "secret"
	cfg.ApplicationName = "bytebase-taskrun-7"
	cfg.ReadOnly = true
	c, err := connConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if c.Database != "sales" || c.User != "dba" || c.Password != "secret" {
		t.Errorf("database %q user %q password %q", c.Database, c.User, c.Password)
	}
	if got := c.RuntimeParams["application_name"]; got != "bytebase-taskrun-7" {
		t.Errorf("application_name %q", got)
	}
	if got := c.RuntimeParams["default_transaction_read_only"]; got != "true" {
		t.Errorf("default_transaction_read_only %q, want true", got)
	}
}

func TestConnConfigParamsOverride(t *testing.T) {
	cfg := baseConfig()
	cfg.Params = map[string]string{
		"default_query_exec_mode":  "simple_protocol",
		"statement_cache_capacity": "128",
		"options":                  "'-c search_path=hr'",
		"sslmode":                  "disable",
	}
	cfg.TLS = func(string) (*tls.Config, error) { return &tls.Config{}, nil }
	c, err := connConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if c.DefaultQueryExecMode != pgx.QueryExecModeSimpleProtocol || c.StatementCacheCapacity != 128 {
		t.Errorf("query mode %v, statement cache %d; want the Params'", c.DefaultQueryExecMode, c.StatementCacheCapacity)
	}
	// A value keeps libpq syntax: the caller quoted it.
	if got := c.RuntimeParams["options"]; got != "-c search_path=hr" {
		t.Errorf("options %q", got)
	}
	// Params cannot turn off the TLS the caller requires.
	if c.TLSConfig == nil {
		t.Errorf("Params sslmode=disable turned TLS off")
	}
}

// TestConnConfigTLSCannotBeDowngraded: with TLS set, no sslmode in Params
// leaves a host that pgx dials in plaintext, and each host is dialed once.
func TestConnConfigTLSCannotBeDowngraded(t *testing.T) {
	for _, mode := range []string{"disable", "allow", "prefer", "require", "verify-full"} {
		t.Run(mode, func(t *testing.T) {
			cfg := baseConfig()
			cfg.Host = "pg-1.example.com,pg-2.example.com"
			cfg.Params = map[string]string{"sslmode": mode}
			cfg.TLS = func(host string) (*tls.Config, error) { return &tls.Config{ServerName: host}, nil }
			c, err := connConfig(cfg)
			if err != nil {
				t.Fatal(err)
			}
			got := []string{c.Host + " " + c.TLSConfig.ServerName}
			for _, f := range c.Fallbacks {
				if f.TLSConfig == nil {
					t.Fatalf("%s is dialed in plaintext", f.Host)
				}
				got = append(got, f.Host+" "+f.TLSConfig.ServerName)
			}
			want := []string{"pg-1.example.com pg-1.example.com", "pg-2.example.com pg-2.example.com"}
			if !slices.Equal(got, want) {
				t.Errorf("hosts dialed %q, want %q", got, want)
			}
		})
	}
}

func TestConnConfigKeepsEachValueWhole(t *testing.T) {
	open := func(edit func(*driver.Config)) (*pgx.ConnConfig, error) {
		cfg := baseConfig()
		edit(&cfg)
		return connConfig(cfg)
	}
	t.Run("a port cannot name another host", func(t *testing.T) {
		if _, err := open(func(c *driver.Config) { c.Port = "5432 host=other.example.com" }); err == nil || !strings.Contains(err.Error(), "invalid port") {
			t.Errorf("err %v, want invalid port", err)
		}
	})
	t.Run("a host carrying options stays one host", func(t *testing.T) {
		c, err := open(func(c *driver.Config) { c.Host = "db.example.com port=6543" })
		if err != nil {
			t.Fatal(err)
		}
		if c.Host != "db.example.com port=6543" || c.Port != 5432 {
			t.Errorf("host %q port %d", c.Host, c.Port)
		}
	})
	t.Run("a user carrying options stays one user", func(t *testing.T) {
		c, err := open(func(c *driver.Config) { c.User = "dba host=other.example.com" })
		if err != nil {
			t.Fatal(err)
		}
		if c.User != "dba host=other.example.com" || c.Host != "db.example.com" {
			t.Errorf("user %q host %q", c.User, c.Host)
		}
	})
	t.Run("surrounding whitespace is ignored", func(t *testing.T) {
		c, err := open(func(c *driver.Config) { c.Host, c.Port = " db.example.com ", " 5432 " })
		if err != nil {
			t.Fatal(err)
		}
		if c.Host != "db.example.com" || c.Port != 5432 {
			t.Errorf("host %q port %d", c.Host, c.Port)
		}
	})
	for _, tc := range []struct {
		name string
		edit func(*driver.Config)
		want string
	}{
		{"a blank host is refused", func(c *driver.Config) { c.Host = "   " }, "host must be set"},
		{"a blank port is refused", func(c *driver.Config) { c.Port = "   " }, "port must be set"},
		{"a blank user is refused", func(c *driver.Config) { c.User = "" }, "user must be set"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := open(tc.edit); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("err %v, want %q", err, tc.want)
			}
		})
	}
}

// TestConnConfigTLSForEveryHost: with TLS set, every host pgx may dial uses
// the configuration the caller made for that host, and none falls back to
// plaintext.
func TestConnConfigTLSForEveryHost(t *testing.T) {
	cfg := baseConfig()
	cfg.Host = "pg-1.example.com,pg-2.example.com,pg-3.example.com"
	var asked []string
	cfg.TLS = func(host string) (*tls.Config, error) {
		asked = append(asked, host)
		return &tls.Config{ServerName: "for " + host}, nil
	}
	c, err := connConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	got := []string{c.Host + ": " + c.TLSConfig.ServerName}
	for _, f := range c.Fallbacks {
		if f.TLSConfig == nil {
			t.Fatalf("fallback %s dials without TLS", f.Host)
		}
		got = append(got, f.Host+": "+f.TLSConfig.ServerName)
	}
	want := []string{
		"pg-1.example.com: for pg-1.example.com",
		"pg-2.example.com: for pg-2.example.com",
		"pg-3.example.com: for pg-3.example.com",
	}
	if !slices.Equal(got, want) || !slices.Equal(asked, []string{"pg-1.example.com", "pg-2.example.com", "pg-3.example.com"}) {
		t.Errorf("TLS per host %q, asked for %q", got, asked)
	}
}

// TestConnConfigDialWithoutHost: a cloud connector's dialer names its
// instance itself, so the host may be left out.
func TestConnConfigDialWithoutHost(t *testing.T) {
	var dialed bool
	cfg := driver.Config{User: "sa@project.iam", Dial: func(context.Context, string, string) (net.Conn, error) {
		dialed = true
		return nil, net.ErrClosed
	}}
	c, err := connConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if c.DialFunc == nil {
		t.Fatal("DialFunc is not the caller's")
	}
	_, _ = c.DialFunc(context.Background(), "tcp", "ignored")
	if !dialed {
		t.Error("DialFunc is not the caller's")
	}
}

// TestOpenDoesNotConnect: Open only prepares the pool, and a failed config
// opens nothing.
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
	cfg.User = ""
	if _, err := Open(context.Background(), cfg); err == nil {
		t.Error("Open without a user succeeded")
	}
}
