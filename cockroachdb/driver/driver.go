// Package driver is CockroachDB's implementation of the driver contract in
// github.com/bytebase/omni/driver: it connects through pgx, the wire protocol
// CockroachDB speaks, syncs a server's catalogs into omni/metadata, dumps the
// server's own DDL, and runs interactive queries and change scripts.
//
// CockroachDB is not an alias of pg. Its wire protocol is PostgreSQL's, but
// what it stores, reports and plans with is its own: pg_database has no
// owner in the sync, system databases are postgres and system rather than
// template0/template1, pg_relation_size does not exist, pg_matviews is
// populated by CockroachDB's materialized views, indexes come from
// pg_indexes, and EXPLAIN has CockroachDB's option list (DISTSQL) rather
// than PostgreSQL's FORMAT. Its EXPLAIN also carries the only
// machine-readable plan it publishes behind the diagram link it prints, so
// the query path decodes that. The sync, stats and query code here is
// therefore its own, not pg's, even though Bytebase dispatches the engine
// as a pg-family engine for parsing.
//
// Open does not reach the server; Ping does.
package driver

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"

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

// Driver is an open CockroachDB driver: a pool of connections to one
// database.
type Driver struct {
	db            *sql.DB
	databaseName  string
	engineVersion string
}

// defaultDatabase is the database a Config with none connects to.
const defaultDatabase = "postgres"

// Open returns a driver for the server and database cfg names. Open does not
// reach the server; Ping does.
func Open(_ context.Context, cfg driver.Config) (*Driver, error) {
	pgxConfig, err := connConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &Driver{
		db:            stdlib.OpenDB(*pgxConfig),
		databaseName:  pgxConfig.Database,
		engineVersion: cfg.EngineVersion,
	}, nil
}

// connConfig is the pgx configuration Open connects with.
//
// Host and Port are libpq's host and port lists. TLS, when set, makes TLS
// required (sslmode=require) and supplies the configuration for every host
// pgx may dial. Params are connection string keywords with values in libpq
// syntax, applied last so they override the settings above.
func connConfig(cfg driver.Config) (*pgx.ConnConfig, error) {
	if cfg.User == "" {
		return nil, fmt.Errorf("user must be set")
	}
	var dsn []string
	// A cloud connector's Dial names its instance itself and needs no host.
	if cfg.Dial == nil || cfg.Host != "" || cfg.Port != "" {
		if strings.TrimSpace(cfg.Host) == "" {
			return nil, fmt.Errorf("host must be set")
		}
		if strings.TrimSpace(cfg.Port) == "" {
			return nil, fmt.Errorf("port must be set")
		}
		dsn = append(dsn, "host="+quoteValue(cfg.Host), "port="+quoteValue(cfg.Port))
	}
	if cfg.TLS != nil {
		dsn = append(dsn, "sslmode=require")
	}
	if routingID := getRoutingIDFromCockroachCloudURL(cfg.Host); routingID != "" {
		dsn = append(dsn, "options="+quoteValue("--cluster="+routingID))
	}
	for key, value := range cfg.Params {
		dsn = append(dsn, key+"="+value)
	}
	connConfig, err := pgx.ParseConfig(strings.Join(dsn, " "))
	if err != nil {
		return nil, err
	}
	connConfig.User = cfg.User
	connConfig.Password = cfg.Password
	connConfig.Database = cfg.Database
	if connConfig.Database == "" {
		connConfig.Database = defaultDatabase
	}
	if cfg.TLS != nil {
		if connConfig.TLSConfig, err = cfg.TLS(connConfig.Host); err != nil {
			return nil, err
		}
		type endpoint struct {
			host string
			port uint16
		}
		seen := map[endpoint]bool{{connConfig.Host, connConfig.Port}: true}
		var fallbacks []*pgconn.FallbackConfig
		for _, fallback := range connConfig.Fallbacks {
			if seen[endpoint{fallback.Host, fallback.Port}] {
				continue
			}
			seen[endpoint{fallback.Host, fallback.Port}] = true
			if fallback.TLSConfig, err = cfg.TLS(fallback.Host); err != nil {
				return nil, err
			}
			fallbacks = append(fallbacks, fallback)
		}
		connConfig.Fallbacks = fallbacks
	}
	if cfg.Dial != nil {
		connConfig.DialFunc = cfg.Dial
	}
	if _, ok := cfg.Params["default_query_exec_mode"]; !ok {
		connConfig.DefaultQueryExecMode = pgx.QueryExecModeExec
	}
	if _, ok := cfg.Params["statement_cache_capacity"]; !ok {
		connConfig.StatementCacheCapacity = 0
	}
	if cfg.ApplicationName != "" {
		connConfig.RuntimeParams["application_name"] = cfg.ApplicationName
	}
	if cfg.ReadOnly {
		connConfig.RuntimeParams["default_transaction_read_only"] = "true"
	}
	return connConfig, nil
}

// getRoutingIDFromCockroachCloudURL returns the routing ID from a Cockroach
// Cloud URL, or the empty string when the host is not one.
func getRoutingIDFromCockroachCloudURL(host string) string {
	host = strings.TrimSpace(host)
	if !strings.HasSuffix(host, "cockroachlabs.cloud") {
		return ""
	}
	parts := strings.Split(host, ".")
	// routing-id[.xxx].cockroachlabs.cloud
	if len(parts) > 2 {
		return parts[0]
	}
	return ""
}

// quoteValue quotes a libpq connection string value, so that a value
// holding a space cannot name another keyword.
func quoteValue(value string) string {
	return "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(strings.TrimSpace(value)) + "'"
}

func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func quoteLiteral(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// Ping reaches the server.
func (d *Driver) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// Close closes the driver's connections.
func (d *Driver) Close() error {
	return d.db.Close()
}

// getDatabases reads every database of the server, system ones included, in
// name order: pg_database is shared by every database, so its heap order
// moves as databases come and go.
func (d *Driver) getDatabases(ctx context.Context) ([]*metadatapb.DatabaseSchemaMetadata, error) {
	var databases []*metadatapb.DatabaseSchemaMetadata
	rows, err := d.db.QueryContext(ctx, "SELECT datname, pg_encoding_to_char(encoding), datcollate FROM pg_database ORDER BY datname;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		database := &metadatapb.DatabaseSchemaMetadata{}
		if err := rows.Scan(&database.Name, &database.CharacterSet, &database.Collation); err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return databases, nil
}

func (d *Driver) getVersion(ctx context.Context) (string, error) {
	// https://www.cockroachlabs.com/docs/v25.1/cluster-settings#setting-version
	query := "SHOW CLUSTER SETTING version;"
	var version string
	if err := d.db.QueryRowContext(ctx, query).Scan(&version); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("query %q returned empty row", query)
		}
		return "", formatErrorWithQuery(err, query)
	}
	return version, nil
}

func formatErrorWithQuery(err error, query string) error {
	return fmt.Errorf("failed to execute query %q: %w", query, err)
}

// convertYesNo converts a PostgreSQL information_schema YES/NO to a bool.
func convertYesNo(s string) (bool, error) {
	switch strings.ToUpper(s) {
	case "YES":
		return true, nil
	case "NO":
		return false, nil
	default:
		return false, fmt.Errorf("invalid yes/no value %q", s)
	}
}
