package driver

import (
	"context"
	"crypto/tls"
	"net"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/bytebase/omni/driver"
)

// TestConnConfigBuildsTheConnection covers the settings Open hands pgx: which
// ones are required, which have a default, and how a caller's TLS, Dial and
// Params reach the connection. None of it needs a server.
func TestConnConfigBuildsTheConnection(t *testing.T) {
	for _, tc := range []struct {
		name string
		cfg  driver.Config
		want func(t *testing.T, got *pgx.ConnConfig)
	}{
		{
			name: "the default database is Redshift's own",
			cfg:  driver.Config{User: "dba", Host: "redshift.example.com", Port: "5439"},
			want: func(t *testing.T, got *pgx.ConnConfig) {
				if got.Database != "dev" {
					t.Errorf("Database = %q, want %q", got.Database, "dev")
				}
				if got.User != "dba" {
					t.Errorf("User = %q, want dba", got.User)
				}
				if got.RuntimeParams["default_transaction_read_only"] != "" {
					t.Error("the session was made read-only without ReadOnly")
				}
			},
		},
		{
			name: "the caller's database, password and host list are kept whole",
			cfg: driver.Config{
				User:     "dba",
				Password: "s3cret",
				Host:     "one.example.com,two.example.com",
				Port:     "5439,5440",
				Database: "analytics",
			},
			want: func(t *testing.T, got *pgx.ConnConfig) {
				if got.Database != "analytics" {
					t.Errorf("Database = %q, want analytics", got.Database)
				}
				if got.Password != "s3cret" {
					t.Errorf("Password = %q", got.Password)
				}
				if got.Host != "one.example.com" || got.Port != 5439 {
					t.Errorf("host = %q:%d, want one.example.com:5439", got.Host, got.Port)
				}
				if len(got.Fallbacks) == 0 {
					t.Fatal("the second host is not a fallback")
				}
			},
		},
		{
			name: "TLS is required and supplied for every host the driver dials",
			cfg: driver.Config{
				User:   "dba",
				Host:   "one.example.com,two.example.com",
				Port:   "5439",
				Params: map[string]string{"sslmode": "prefer"},
				TLS: func(host string) (*tls.Config, error) {
					return &tls.Config{ServerName: host}, nil
				},
			},
			want: func(t *testing.T, got *pgx.ConnConfig) {
				if got.TLSConfig == nil || got.TLSConfig.ServerName != "one.example.com" {
					t.Errorf("TLSConfig = %v, want one.example.com's", got.TLSConfig)
				}
				// sslmode=prefer adds a plaintext attempt per host; every host
				// is dialed once, with the caller's TLS.
				if len(got.Fallbacks) == 0 {
					t.Fatal("no fallback was configured")
				}
				for _, fallback := range got.Fallbacks {
					if fallback.TLSConfig == nil || fallback.TLSConfig.ServerName != fallback.Host {
						t.Errorf("%s:%d dialed with %v, want the caller's TLS for %s",
							fallback.Host, fallback.Port, fallback.TLSConfig, fallback.Host)
					}
				}
			},
		},
		{
			name: "a dialer that names its own endpoint needs no host",
			cfg: driver.Config{
				User: "dba",
				Dial: func(context.Context, string, string) (net.Conn, error) { return nil, nil },
			},
			want: func(t *testing.T, got *pgx.ConnConfig) {
				if got.DialFunc == nil {
					t.Error("DialFunc is not the caller's dialer")
				}
			},
		},
		{
			name: "read-only is a session default",
			cfg:  driver.Config{User: "dba", Host: "redshift.example.com", Port: "5439", ReadOnly: true},
			want: func(t *testing.T, got *pgx.ConnConfig) {
				if got.RuntimeParams["default_transaction_read_only"] != "true" {
					t.Errorf("default_transaction_read_only = %q, want true", got.RuntimeParams["default_transaction_read_only"])
				}
			},
		},
		{
			name: "a datashare database is not made read-only",
			cfg: driver.Config{
				User:     "dba",
				Host:     "redshift.example.com",
				Port:     "5439",
				ReadOnly: true,
				Params:   map[string]string{ParamDatashare: "true"},
			},
			want: func(t *testing.T, got *pgx.ConnConfig) {
				if _, ok := got.RuntimeParams["default_transaction_read_only"]; ok {
					t.Error("a datashare database was asked for a read-only session")
				}
			},
		},
		{
			name: "the application name and Params reach the session",
			cfg: driver.Config{
				User:            "dba",
				Host:            "redshift.example.com",
				Port:            "5439",
				ApplicationName: "bytebase",
				Params:          map[string]string{"search_path": "app"},
			},
			want: func(t *testing.T, got *pgx.ConnConfig) {
				if got.RuntimeParams["application_name"] != "bytebase" {
					t.Errorf("application_name = %q, want bytebase", got.RuntimeParams["application_name"])
				}
				if got.RuntimeParams["search_path"] != "app" {
					t.Errorf("search_path = %q, want app", got.RuntimeParams["search_path"])
				}
			},
		},
		{
			name: "queries run as simple statements unless Params say otherwise",
			cfg:  driver.Config{User: "dba", Host: "redshift.example.com", Port: "5439"},
			want: func(t *testing.T, got *pgx.ConnConfig) {
				if got.DefaultQueryExecMode != pgx.QueryExecModeExec {
					t.Errorf("DefaultQueryExecMode = %v, want exec", got.DefaultQueryExecMode)
				}
				if got.StatementCacheCapacity != 0 {
					t.Errorf("StatementCacheCapacity = %d, want 0", got.StatementCacheCapacity)
				}
			},
		},
		{
			name: "a Params execution mode overrides the default",
			cfg: driver.Config{
				User:   "dba",
				Host:   "redshift.example.com",
				Port:   "5439",
				Params: map[string]string{"default_query_exec_mode": "simple_protocol", "statement_cache_capacity": "512"},
			},
			want: func(t *testing.T, got *pgx.ConnConfig) {
				if got.DefaultQueryExecMode != pgx.QueryExecModeSimpleProtocol {
					t.Errorf("DefaultQueryExecMode = %v, want simple protocol", got.DefaultQueryExecMode)
				}
				if got.StatementCacheCapacity != 512 {
					t.Errorf("StatementCacheCapacity = %d, want 512", got.StatementCacheCapacity)
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := connConfig(tc.cfg)
			if err != nil {
				t.Fatalf("connConfig: %v", err)
			}
			tc.want(t, got)
		})
	}
}

// TestConnConfigRefusesWhatCannotNameAServer covers the values Open refuses
// rather than letting pgx fall back to this machine's socket or port 5432.
func TestConnConfigRefusesWhatCannotNameAServer(t *testing.T) {
	for _, tc := range []struct {
		name    string
		cfg     driver.Config
		wantErr string
	}{
		{
			name:    "no user",
			cfg:     driver.Config{Host: "redshift.example.com", Port: "5439"},
			wantErr: "user must be set",
		},
		{
			name:    "a blank host",
			cfg:     driver.Config{User: "dba", Host: "  ", Port: "5439"},
			wantErr: "host must be set",
		},
		{
			name:    "a blank port",
			cfg:     driver.Config{User: "dba", Host: "redshift.example.com"},
			wantErr: "port must be set",
		},
		{
			// A value holding a keyword must not be read as one.
			name:    "a port naming another host",
			cfg:     driver.Config{User: "dba", Host: "redshift.example.com", Port: "5439 host=other.example.com"},
			wantErr: "invalid port",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := connConfig(tc.cfg)
			if err == nil {
				t.Fatal("connConfig accepted the configuration")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("error = %q, want it to hold %q", err, tc.wantErr)
			}
		})
	}
}

// TestOpenReportsADatashare checks that Open, which reaches no server, records
// the database's kind for the sync, the queries and Execute to read.
func TestOpenReportsADatashare(t *testing.T) {
	for _, tc := range []struct {
		name  string
		param string
		want  bool
	}{
		{name: "unset"},
		{name: "true", param: "true", want: true},
		{name: "1", param: "1", want: true},
		{name: "yes in another case", param: "YES", want: true},
		{name: "false", param: "false"},
		{name: "empty", param: ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := driver.Config{
				User:     "dba",
				Host:     "redshift.example.com",
				Port:     "5439",
				Database: "shared",
				Params:   map[string]string{ParamDatashare: tc.param},
			}
			d, err := Open(context.Background(), cfg)
			if err != nil {
				t.Fatalf("Open: %v", err)
			}
			defer d.Close()
			if d.datashare != tc.want {
				t.Errorf("datashare = %v, want %v", d.datashare, tc.want)
			}
			if d.databaseName != "shared" {
				t.Errorf("databaseName = %q, want shared", d.databaseName)
			}
		})
	}
}
