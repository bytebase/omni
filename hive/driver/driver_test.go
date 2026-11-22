package driver

import (
	"context"
	"crypto/tls"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
)

func TestOpenRequiresAHost(t *testing.T) {
	_, err := Open(context.Background(), driver.Config{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "hostname not set")
}

// TestConnConfig maps the contract's settings onto gohive's: TLS is asked for
// the host it dials, the engine's own parameters reach the SASL exchange, and
// every other parameter reaches the session as a Hive setting.
func TestConnConfig(t *testing.T) {
	tlsConfig := &tls.Config{}
	asked := ""
	hc, err := connConfig(driver.Config{
		Host:     "hive.example.com",
		Port:     "10001",
		Database: "omni",
		User:     "omni",
		Password: "secret",
		Params: map[string]string{
			"auth":             "LDAP",
			"service":          "hiveserver",
			"transport":        "http",
			"http_path":        "cliservice",
			"hive.exec.mode":   "local",
			"hive.mapred.mode": "strict",
		},
		TLS: func(host string) (*tls.Config, error) {
			asked = host
			return tlsConfig, nil
		},
	})
	require.NoError(t, err)
	require.Equal(t, "hive.example.com", hc.Host)
	require.Equal(t, 10001, hc.Port)
	require.Equal(t, "omni", hc.Database)
	require.Equal(t, "omni", hc.Username)
	require.Equal(t, "secret", hc.Password)
	require.Equal(t, "LDAP", hc.Auth)
	require.Equal(t, "hiveserver", hc.Service)
	require.Equal(t, "http", hc.TransportMode)
	require.Equal(t, "cliservice", hc.HTTPPath)
	require.Equal(t, map[string]string{"hive.exec.mode": "local", "hive.mapred.mode": "strict"}, hc.HiveConfiguration)
	require.Same(t, tlsConfig, hc.TLSConfig)
	require.Equal(t, "hive.example.com", asked)
}

// TestConnConfigDefaults: a caller that names only the host gets
// HiveServer2's port and the SASL exchange its default configuration has.
func TestConnConfigDefaults(t *testing.T) {
	hc, err := connConfig(driver.Config{Host: "hive.example.com"})
	require.NoError(t, err)
	require.Equal(t, 10000, hc.Port)
	require.Equal(t, "NONE", hc.Auth)
	require.Equal(t, "hive", hc.Service)
	require.Empty(t, hc.Database)
	require.Empty(t, hc.TransportMode)
	require.Nil(t, hc.HiveConfiguration)
	require.Nil(t, hc.TLSConfig)
}

func TestConnConfigRejects(t *testing.T) {
	_, err := connConfig(driver.Config{Host: "hive.example.com", Port: "10000x"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid port")

	_, err = connConfig(driver.Config{
		Host: "hive.example.com",
		TLS:  func(string) (*tls.Config, error) { return nil, errors.New("no certificate") },
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no certificate")
}

func TestSplitCommands(t *testing.T) {
	cases := []struct {
		name    string
		script  string
		want    []string
		offsets [][2]int
	}{
		{
			name:    "one statement",
			script:  "SELECT 1;",
			want:    []string{"SELECT 1;"},
			offsets: [][2]int{{0, 9}},
		},
		{
			name:    "two statements",
			script:  "SELECT 1;\nSELECT 2;",
			want:    []string{"SELECT 1;", "\nSELECT 2;"},
			offsets: [][2]int{{0, 9}, {9, 19}},
		},
		{
			name:   "no terminator",
			script: "SELECT 1",
			want:   []string{"SELECT 1"},
		},
		{
			name:   "a comment line belongs to the statement after it",
			script: "-- a comment\nSELECT 1;",
			want:   []string{"-- a comment\nSELECT 1;"},
		},
		{
			name:   "a comment after the last statement belongs to none",
			script: "SELECT 1;\n-- trailing\n",
			want:   []string{"SELECT 1;"},
		},
		{
			name:   "an empty script",
			script: "\n\n",
			want:   []string{},
		},
		{
			name:   "a block comment that ends its line",
			script: "/* a comment */\nSELECT 1;",
			want:   []string{"/* a comment */\nSELECT 1;"},
		},
		{
			name:   "an inline block comment before a statement",
			script: "/* a comment */ SELECT 1;",
			want:   []string{"/* a comment */ SELECT 1;"},
		},
		{
			name:    "two statements on one line",
			script:  "SELECT 'a;b'; SELECT 2;",
			want:    []string{"SELECT 'a;b';", " SELECT 2;"},
			offsets: [][2]int{{0, 13}, {13, 23}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			commands, err := splitCommands(tc.script)
			require.NoError(t, err)
			texts := make([]string, 0, len(commands))
			for _, command := range commands {
				texts = append(texts, command.Text)
				require.Equal(t, command.Text, tc.script[command.ByteStart:command.ByteEnd])
			}
			require.Equal(t, tc.want, texts)
			for i, offsets := range tc.offsets {
				require.Equal(t, offsets[0], commands[i].ByteStart)
				require.Equal(t, offsets[1], commands[i].ByteEnd)
			}
		})
	}
}
