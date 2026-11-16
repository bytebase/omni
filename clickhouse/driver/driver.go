// Package driver is ClickHouse's implementation of the driver contract in
// github.com/bytebase/omni/driver: it connects through clickhouse-go/v2,
// syncs a server's databases, users and a database's tables, views and
// data-skipping indexes into omni/metadata, runs interactive queries, and
// executes change scripts.
//
// It implements Driver, Syncer, Querier and Executor. There is no Dumper:
// omni's clickhouse tree has no catalog, so nothing renders synced metadata as
// ClickHouse DDL. Bytebase's driver has a Dump, but it renders from the
// metadata a sync produced, so it moves when a clickhouse/catalog exists. There
// is also no clickhouse/parser or clickhouse/analysis yet, so the statement
// handling the driver needs — splitting, limit and explain rewriting, and the
// read-only check — lives in this package until that tree holds it.
//
// Config.Host and Config.Port name the native endpoint, port 9000. Config.Database
// is the database a sync reads and a query runs in. Config.User and
// Config.Password are the login; the image's own account is "default" with no
// password. Config.TLS supplies the TLS settings for the host and makes TLS
// required; nil leaves the client's own negotiation. Config.Dial opens the
// connection the client makes, so an SSH tunnel or a cloud connector reaches
// the server unchanged. Config.Params are ClickHouse settings, applied after
// max_execution_time, which this driver sets to 300 seconds as Bytebase did to
// keep a resource-intensive query from timing out; a value that parses as an
// integer is sent as one, since ClickHouse settings are typed.
//
// Config.ApplicationName does nothing: the native protocol carries no
// application name. Config.ReadOnly does nothing: Bytebase ignored it, and so
// does this driver, because a SQL editor session has to run both reads and
// writes. Config.EngineVersion is used only when the server does not answer the
// version query.
//
// QueryOptions.Schema does nothing: a connection opens on Config.Database, and
// Bytebase did not change it per query. QueryOptions.Explain rewrites the
// statement with EXPLAIN, and JSON with `indexes = 1`; a statement that already
// starts with EXPLAIN runs as written. QueryOptions.MaskingEnabled does
// nothing: a ClickHouse query returns only the columns its statement traces.
// Result.Plan stays nil because ClickHouse answers a plan as rows and nothing
// here states its format.
//
// ExecuteOptions.CreateDatabase and ExecuteOptions.MaxRetries do nothing:
// Bytebase ignored them. ExecuteOptions.Isolation is passed to the client when
// it is set; ClickHouse's transactions are limited and clickhouse-go's begin is
// a no-op, so it may have no effect.
package driver

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"

	"github.com/bytebase/omni/driver"
)

var (
	_ driver.Driver = (*Driver)(nil)
)

// Driver is an open ClickHouse driver: a pool of connections to one server
// and the database its syncs and queries read.
type Driver struct {
	db            *sql.DB
	databaseName  string
	engineVersion string
}

// Open returns a driver for the server and database cfg names. It does not
// reach the server; Ping does.
func Open(_ context.Context, cfg driver.Config) (*Driver, error) {
	options, err := connOptions(cfg)
	if err != nil {
		return nil, err
	}
	return &Driver{
		db:            sql.OpenDB(clickhouse.Connector(options)),
		databaseName:  cfg.Database,
		engineVersion: cfg.EngineVersion,
	}, nil
}

func connOptions(cfg driver.Config) (*clickhouse.Options, error) {
	host := strings.TrimSpace(cfg.Host)
	if host == "" {
		return nil, fmt.Errorf("host must be set")
	}
	port := strings.TrimSpace(cfg.Port)
	if port == "" {
		return nil, fmt.Errorf("port must be set")
	}

	settings := clickhouse.Settings{
		// Bytebase used a relatively long value to avoid a timeout on a
		// resource-intensive query: "Estimated query execution time (xxx
		// seconds) is too long."
		"max_execution_time": 300,
	}
	for name, value := range cfg.Params {
		if n, err := strconv.ParseInt(value, 10, 64); err == nil {
			settings[name] = n
		} else {
			settings[name] = value
		}
	}

	options := &clickhouse.Options{
		Addr: []string{net.JoinHostPort(host, port)},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
		Settings:    settings,
		DialTimeout: 10 * time.Second,
	}
	if cfg.TLS != nil {
		tlsConfig, err := cfg.TLS(host)
		if err != nil {
			return nil, fmt.Errorf("TLS config: %w", err)
		}
		if tlsConfig == nil {
			return nil, fmt.Errorf("TLS configuration must not be nil")
		}
		options.TLS = tlsConfig
	}
	if cfg.Dial != nil {
		options.DialContext = func(ctx context.Context, address string) (net.Conn, error) {
			return cfg.Dial(ctx, "tcp", address)
		}
	}
	return options, nil
}

// Ping reaches the server.
func (d *Driver) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// Close closes the driver's connections.
func (d *Driver) Close() error {
	return d.db.Close()
}
