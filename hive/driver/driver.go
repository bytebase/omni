// Package driver is Hive's implementation of the driver contract in
// github.com/bytebase/omni/driver: it connects to a HiveServer2 through gohive
// v2, syncs a server's metadata into omni/metadata, and runs interactive
// queries and change scripts.
//
// There is no Dump: Bytebase's Hive dump asks the server for SHOW CREATE TABLE
// rather than rendering metadata, and a tree with no hive/catalog has nothing
// to render. Dump's contract renders the schema it is given and nothing else.
//
// Hive has no transaction the client controls. Only ACID tables are
// transactional, and HiveServer2 4.0 refuses BEGIN, COMMIT and ROLLBACK, so
// every statement commits itself: Execute ignores ExecuteOptions.Transaction
// and ExecuteOptions.Isolation, and a failed script leaves the statements it
// already ran.
//
// The rest of the contract's options do not apply either, and say why here
// rather than at each call:
//
//   - Config.ReadOnly: Hive has no read-only session setting.
//   - Config.ApplicationName and Config.EngineVersion: HiveServer2 carries no
//     application name, and a driver opened without reaching the server records
//     no version; SyncInstance asks for one.
//   - Config.Dial: gohive dials the host itself, so a caller that must reach
//     HiveServer2 over an SSH tunnel or a cloud connector cannot use this
//     driver yet.
//   - QueryOptions.Limit: the driver does not rewrite HiveQL — there is no
//     hive/parser to place a LIMIT with — and Bytebase never asked for one, so
//     QueryOptions.MaxResultBytes is the only cap a statement gets.
//   - driver.ErrNotFound: gohive v2 renders a failed TStatus as pointer
//     addresses and keeps the message and Hive's error code in fields of a
//     type it does not export, so a database that does not exist reaches the
//     caller as an error no caller can recognize.
//
// Statement classification is not the driver's here. Bytebase classified a
// statement for the SQL editor through parser/standard's text check, outside
// the driver, and the driver ran every query as a query: Hive answers a
// statement that returns no rows with an operation that has no result set, and
// reading that operation's rows reports Hive's error, which is what a DDL sent
// to Query gets.
package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	gohive "github.com/beltran/gohive/v2"

	"github.com/bytebase/omni/driver"
)

var (
	_ driver.Driver   = (*Driver)(nil)
	_ driver.Syncer   = (*Driver)(nil)
	_ driver.Querier  = (*Driver)(nil)
	_ driver.Executor = (*Driver)(nil)
)

// Driver is an open Hive driver: a pool of connections to one HiveServer2 and
// an optional database.
type Driver struct {
	db           *sql.DB
	databaseName string
}

// Open returns a driver for the server and database cfg names. It does not
// reach the server; Ping does.
func Open(_ context.Context, cfg driver.Config) (*Driver, error) {
	hc, err := connConfig(cfg)
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(gohive.NewConnector(hc))
	// HiveServer2 does not serve many concurrent connections well.
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(0)
	return &Driver{db: db, databaseName: cfg.Database}, nil
}

// connConfig maps the contract's settings onto gohive's. Params carries the
// engine's own vocabulary: auth, service and transport name the SASL and
// thrift exchange, http_path the route of an HTTP transport, and every other
// entry reaches the session as a Hive configuration property.
func connConfig(cfg driver.Config) (gohive.Config, error) {
	if cfg.Host == "" {
		return gohive.Config{}, errors.New("hostname not set")
	}
	port := 10000
	if cfg.Port != "" {
		n, err := strconv.Atoi(cfg.Port)
		if err != nil {
			return gohive.Config{}, fmt.Errorf("invalid port %q: %w", cfg.Port, err)
		}
		port = n
	}
	hc := gohive.Config{
		Host:     cfg.Host,
		Port:     port,
		Auth:     param(cfg.Params, "auth", "NONE"),
		Username: cfg.User,
		Password: cfg.Password,
		Database: cfg.Database,
		Service:  param(cfg.Params, "service", "hive"),
	}
	if v := param(cfg.Params, "transport", ""); v != "" {
		hc.TransportMode = v
	}
	hc.HTTPPath = param(cfg.Params, "http_path", "")
	for key, value := range cfg.Params {
		switch strings.ToLower(key) {
		case "auth", "service", "transport", "http_path":
			continue
		}
		if hc.HiveConfiguration == nil {
			hc.HiveConfiguration = map[string]string{}
		}
		hc.HiveConfiguration[key] = value
	}
	if cfg.TLS != nil {
		tlsConfig, err := cfg.TLS(cfg.Host)
		if err != nil {
			return gohive.Config{}, err
		}
		if tlsConfig == nil {
			return gohive.Config{}, errors.New("TLS configuration must not be nil")
		}
		hc.TLSConfig = tlsConfig
	}
	return hc, nil
}

// param reads a connection parameter by name, case-insensitively.
func param(params map[string]string, name, fallback string) string {
	for key, value := range params {
		if strings.EqualFold(key, name) {
			return value
		}
	}
	return fallback
}

// Ping reaches the server.
func (d *Driver) Ping(ctx context.Context) error {
	if d.db == nil {
		return errors.New("connection not initialized")
	}
	if err := d.db.PingContext(ctx); err != nil {
		return fmt.Errorf("bad connection: %w", err)
	}
	return nil
}

// Close closes the driver's connections.
func (d *Driver) Close() error {
	if d.db == nil {
		return nil
	}
	return d.db.Close()
}
