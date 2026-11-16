// Package driver is Cassandra's implementation of the driver contract in
// github.com/bytebase/omni/driver: it connects with gocql, syncs an instance's
// keyspaces and one keyspace's tables, indexes and columns into
// omni/metadata, and runs interactive queries and change scripts.
//
// It implements Driver, Syncer, Querier and Executor. There is no Dumper: omni
// has no cassandra/catalog for a dump to render, and Bytebase's Dump wrote
// nothing. No row count or data size is read either, because Bytebase asked
// Cassandra for neither, so a synced table carries no estimate.
//
// Config.Host is one Cassandra contact point or a comma-separated list of
// them, each optionally carrying its own port; Config.Port applies to a single
// host that has none, and gocql's default 9042 to both when neither does.
// Config.Database is the keyspace every sync and unqualified query runs in.
// Config.User and Config.Password are the login. Config.TLS makes TLS required
// and answers per host; Config.Dial opens the network connection in place of
// net.Dialer, so an SSH tunnel or a cloud connector reaches the cluster.
//
// Config.Params, Config.ApplicationName, Config.ReadOnly, Config.EngineVersion
// and QueryOptions' Schema, Explain, ExplainFormat and MaskingEnabled do
// nothing. gocql takes no connection string this driver fills, and Cassandra
// has no application name, no read-only session, no EXPLAIN, no per-session
// schema and no masking to withhold. ExecuteOptions' transaction settings and
// MaxRetries do nothing either: Cassandra has no multi-statement transaction
// and no lock timeout to retry.
//
// QueryOptions.Limit is a page size: it caps a statement at that many rows and
// a statement that returns more is read page by page until it reaches the
// limit. QueryOptions.MaxResultBytes stops reading at the cap and sets
// Result.Truncated.
//
// A Column's TypeName is Cassandra's own lowercase name for the CQL type (int,
// text, list), not the upper-case name the contract describes: gocql reports
// the lowercase one, and Bytebase showed it as-is.
//
// Open does not reach the server; the gocql session is built on first use.
// Every method is safe for concurrent use.
package driver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	gocql "github.com/apache/cassandra-gocql-driver/v2"

	"github.com/bytebase/omni/driver"
)

var (
	_ driver.Driver   = (*Driver)(nil)
	_ driver.Syncer   = (*Driver)(nil)
	_ driver.Querier  = (*Driver)(nil)
	_ driver.Executor = (*Driver)(nil)
)

// Driver is one open Cassandra cluster.
type Driver struct {
	config driver.Config

	// client replaces the gocql session, which is what a test's tape does.
	// Nil means the session is built from config on first use, so Open
	// reaches no server.
	client client

	mu   sync.Mutex
	live *liveClient
}

// Open returns a driver for the cluster and keyspace cfg names. It does not
// reach the cluster; Ping does.
func Open(_ context.Context, cfg driver.Config) (*Driver, error) {
	if _, err := endpoints(cfg.Host, cfg.Port); err != nil {
		return nil, err
	}
	return &Driver{config: cfg}, nil
}

// conn returns the client, building the gocql session the first time. A
// failure is not cached: a cluster that was unreachable when the driver was
// first used answers the next call, as a caller retrying after an outage
// expects.
func (d *Driver) conn() (client, error) {
	if d.client != nil {
		return d.client, nil
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.live != nil {
		return d.live, nil
	}
	cluster, err := clusterConfig(d.config)
	if err != nil {
		return nil, err
	}
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	d.live = &liveClient{session: session}
	return d.live, nil
}

// Ping asks the cluster for its release version.
func (d *Driver) Ping(ctx context.Context) error {
	if _, err := d.getVersion(ctx); err != nil {
		return fmt.Errorf("failed to ping: %w", err)
	}
	return nil
}

// Close closes the gocql session, if one was opened.
func (d *Driver) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	switch {
	case d.client != nil:
		d.client.Close()
	case d.live != nil:
		d.live.Close()
	}
	return nil
}

// clusterConfig is the gocql cluster the session is created from.
func clusterConfig(cfg driver.Config) (*gocql.ClusterConfig, error) {
	addrs, err := endpoints(cfg.Host, cfg.Port)
	if err != nil {
		return nil, err
	}
	cluster := gocql.NewCluster(addrs...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.User,
		Password: cfg.Password,
	}
	cluster.Keyspace = cfg.Database
	if cfg.TLS != nil || cfg.Dial != nil {
		// The dialer owns both the connection and the TLS session, so a TLS
		// setting can be answered per host the cluster publishes.
		cluster.HostDialer = &hostDialer{config: cfg}
	}
	return cluster, nil
}

// endpoints turns Config.Host and Config.Port into gocql's contact points.
func endpoints(host, port string) ([]string, error) {
	host = strings.TrimSpace(host)
	if host == "" {
		return nil, errors.New("host must be set")
	}
	parts := strings.Split(host, ",")
	addrs := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, fmt.Errorf("host list %q has an empty entry", host)
		}
		if _, _, err := net.SplitHostPort(part); err == nil {
			addrs = append(addrs, part)
			continue
		}
		// An address without its port: brackets come off a literal so
		// JoinHostPort can put them back.
		bare := strings.TrimSuffix(strings.TrimPrefix(part, "["), "]")
		// A list names each host's own port; Config.Port applies to a single
		// host.
		p := ""
		if len(parts) == 1 {
			p = strings.TrimSpace(port)
		}
		if p == "" {
			if net.ParseIP(bare) != nil {
				addrs = append(addrs, "["+bare+"]")
				continue
			}
			addrs = append(addrs, part)
			continue
		}
		addrs = append(addrs, net.JoinHostPort(bare, p))
	}
	return addrs, nil
}

// hostDialer opens one host's connection through Config.Dial when it is set
// and hands it to TLS with the settings Config.TLS returns for that host.
type hostDialer struct {
	config driver.Config
}

// DialHost implements gocql.HostDialer.
func (d *hostDialer) DialHost(ctx context.Context, host *gocql.HostInfo) (*gocql.DialedHost, error) {
	conn, err := dial(ctx, d.config.Dial, "tcp", host.ConnectAddressAndPort())
	if err != nil {
		return nil, err
	}
	if d.config.TLS == nil {
		return &gocql.DialedHost{Conn: conn}, nil
	}
	tlsConfig, err := d.config.TLS(hostname(host))
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to get TLS config: %w", err)
	}
	if tlsConfig == nil {
		_ = conn.Close()
		return nil, errors.New("TLS configuration must not be nil")
	}
	return gocql.WrapTLS(ctx, conn, host.ConnectAddressAndPort(), tlsConfig)
}

// hostname is the name a host answers to: the one the server announced, or its
// address when it announced none.
func hostname(host *gocql.HostInfo) string {
	name, _, err := net.SplitHostPort(host.HostnameAndPort())
	if err != nil {
		return host.ConnectAddress().String()
	}
	return name
}

// dial opens the network connection, through Config.Dial when the caller set
// one and the default dialer otherwise.
func dial(ctx context.Context, via func(ctx context.Context, network, address string) (net.Conn, error), network, address string) (net.Conn, error) {
	if via != nil {
		return via(ctx, network, address)
	}
	var d net.Dialer
	return d.DialContext(ctx, network, address)
}
