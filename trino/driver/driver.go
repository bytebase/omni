// Package driver is Trino's implementation of the driver contract in
// github.com/bytebase/omni/driver: it connects through trino-go-client, syncs a
// catalog's metadata into omni/metadata, dumps it as DDL, and runs interactive
// queries and change scripts.
//
// Config.Database is the catalog the driver connects to and syncs; empty means
// Trino's default catalog, system.
//
// It implements Driver, Syncer, Dumper, Querier and Executor. It has no planner
// row estimate and no live probe for the reviews: Trino offers no engine-side
// estimate to read, and Bytebase runs no review rule against a Trino instance.
//
// Trino commits DDL on its own and supports transactions only for some
// connectors' DML, so Execute runs each statement on its own whatever
// opts.Transaction says; opts.Isolation and opts.MaxRetries do nothing.
// Config.ReadOnly is ignored: the driver cannot open a read-only session.
//
// A statement's update count reaches the caller only on the client's exec path,
// so the driver sends a statement that returns no rows through Executor.
// Querier decides which statements return rows from their leading keyword;
// see session.Query.
//
// Config.TLS supplies the TLS settings for the one host the client dials and
// selects https; nil selects http. Config.Dial, when set, becomes the
// transport's dialer, so an SSH tunnel or a cloud connector reaches the server
// the same way.
package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	trinodb "github.com/trinodb/trino-go-client/trino"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

var (
	_ driver.Driver   = (*Driver)(nil)
	_ driver.Syncer   = (*Driver)(nil)
	_ driver.Dumper   = (*Driver)(nil)
	_ driver.Querier  = (*Driver)(nil)
	_ driver.Executor = (*Driver)(nil)
)

const (
	// defaultPort is the port a Config with none connects to.
	defaultPort = "8080"
	// defaultUser is the user a Config with none connects as. Trino refuses a
	// session without one.
	defaultUser = "trino"
	// defaultCatalog is the catalog Config.Database falls back to.
	defaultCatalog = "system"
)

// tlsClientSeq names the custom HTTP client each driver registers with the
// client library, whose registry is process-wide and keyed by string.
var tlsClientSeq atomic.Uint64

// Driver is an open Trino driver: a pool of connections to one catalog.
type Driver struct {
	db           *sql.DB
	databaseName string
	// tlsClientKey names the custom HTTP client the driver registered for
	// Config.TLS or Config.Dial; empty when it registered none.
	tlsClientKey string
	httpClient   *http.Client
}

// Open returns a driver for the server and catalog cfg names. It does not reach
// the server; Ping does.
func Open(_ context.Context, cfg driver.Config) (*Driver, error) {
	dbURL, err := buildURL(cfg)
	if err != nil {
		return nil, err
	}
	// Read the catalog from the DSN, not from cfg, so a Params override syncs
	// the catalog the connection reads.
	catalog := dbURL.Query().Get("catalog")
	if catalog == "" {
		catalog = catalogName(cfg)
	}
	d := &Driver{databaseName: catalog}

	// The client reaches a custom TLS or dial configuration through a
	// process-wide registry of HTTP clients keyed by name, so each driver
	// registers one of its own and drops it on Close.
	if cfg.TLS != nil || cfg.Dial != nil {
		host := hostWithoutBrackets(cfg.Host)
		httpClient := &http.Client{Transport: &http.Transport{DialContext: cfg.Dial}}
		if cfg.TLS != nil {
			tlsConfig, err := cfg.TLS(host)
			if err != nil {
				return nil, fmt.Errorf("failed to get TLS config: %w", err)
			}
			if tlsConfig == nil {
				return nil, errors.New("TLS configuration must not be nil")
			}
			httpClient.Transport.(*http.Transport).TLSClientConfig = tlsConfig
		}
		key := fmt.Sprintf("omni-trino-%d", tlsClientSeq.Add(1))
		if err := trinodb.RegisterCustomClient(key, httpClient); err != nil {
			return nil, fmt.Errorf("failed to register custom TLS client: %w", err)
		}
		d.tlsClientKey = key
		d.httpClient = httpClient
		query := dbURL.Query()
		query.Set("custom_client", key)
		dbURL.RawQuery = query.Encode()
	}

	db, err := sql.Open("trino", dbURL.String())
	if err != nil {
		d.releaseCustomClient()
		return nil, fmt.Errorf("failed to connect to Trino: %w", err)
	}
	db.SetConnMaxLifetime(30 * time.Second)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	d.db = db
	return d, nil
}

// buildURL returns the DSN the client is opened with: the endpoint, the
// resolved login, the catalog Config.Database names, and Config.Params.
func buildURL(cfg driver.Config) (*url.URL, error) {
	scheme := "http"
	if cfg.TLS != nil {
		scheme = "https"
	}

	user := cfg.User
	if user == "" {
		user = defaultUser
	}
	// trino-go-client refuses a password on a plaintext connection, so name the
	// setting to change instead of surfacing its DSN error. Such a data source
	// cannot connect either way: the credentials would go over the wire in the
	// clear.
	if cfg.Password != "" && cfg.TLS == nil {
		return nil, errors.New("trino requires SSL to authenticate with a password: turn on SSL for this data source, or clear the password")
	}

	port := cfg.Port
	if port == "" {
		port = defaultPort
	}
	u := &url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(hostWithoutBrackets(cfg.Host), port),
	}
	if cfg.Password != "" {
		u.User = url.UserPassword(user, cfg.Password)
	} else {
		u.User = url.User(user)
	}

	query := u.Query()
	source := cfg.ApplicationName
	if source == "" {
		source = "bytebase"
	}
	query.Set("source", source)
	query.Set("catalog", catalogName(cfg))
	for key, value := range cfg.Params {
		query.Set(key, value)
	}
	u.RawQuery = query.Encode()
	return u, nil
}

// catalogName returns the catalog the connection and Syncer.SyncDatabase read:
// Config.Database, or Trino's default catalog. SyncDatabase reads it so an empty
// Config.Database syncs the catalog the driver connected to, not nothing.
func catalogName(cfg driver.Config) string {
	if cfg.Database == "" {
		return defaultCatalog
	}
	return cfg.Database
}

// hostWithoutBrackets removes the brackets an IPv6 address may carry.
func hostWithoutBrackets(host string) string {
	host = strings.TrimSpace(host)
	return strings.TrimSuffix(strings.TrimPrefix(host, "["), "]")
}

// Ping reaches the server by running a statement, so a server that answers its
// info endpoint but cannot plan a query fails here too.
func (d *Driver) Ping(ctx context.Context) error {
	if d.db == nil {
		return errors.New("database connection not established")
	}
	var result int
	return d.db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
}

// Close closes the driver's connections and drops the custom HTTP client it
// registered, if any.
func (d *Driver) Close() error {
	d.releaseCustomClient()
	if d.db == nil {
		return nil
	}
	return d.db.Close()
}

func (d *Driver) releaseCustomClient() {
	if d.tlsClientKey != "" {
		trinodb.DeregisterCustomClient(d.tlsClientKey)
		d.tlsClientKey = ""
	}
	if d.httpClient != nil {
		if transport, ok := d.httpClient.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
		d.httpClient = nil
	}
}

// SyncInstance reads the server version and its catalogs. Trino has no
// instance-level roles, so Roles stays empty.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	version, err := d.getVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Trino version: %w", err)
	}
	catalogs, err := d.getCatalog(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get catalog list: %w", err)
	}
	databases := make([]*metadatapb.DatabaseSchemaMetadata, 0, len(catalogs))
	for _, catalog := range catalogs {
		databases = append(databases, &metadatapb.DatabaseSchemaMetadata{Name: catalog})
	}
	return &driver.Instance{Version: version, Databases: databases}, nil
}

func (d *Driver) getVersion(ctx context.Context) (string, error) {
	var version string
	if err := d.db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version); err != nil {
		return "", err
	}
	return version, nil
}

func (d *Driver) getCatalog(ctx context.Context) ([]string, error) {
	return d.queryStringValues(ctx, "SELECT catalog_name FROM system.metadata.catalogs ORDER BY catalog_name")
}

// queryStringValues runs a query whose one column is text.
func (d *Driver) queryStringValues(ctx context.Context, query string, args ...any) ([]string, error) {
	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, value)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}
	return results, nil
}

func logCommandExecute(log driver.ExecuteLogger, command driver.Command) {
	if log != nil {
		log.CommandExecute(command)
	}
}

func logCommandResponse(log driver.ExecuteLogger, affected int64, perResult []int64, err error) {
	if log != nil {
		log.CommandResponse(affected, perResult, err)
	}
}
