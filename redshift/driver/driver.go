// Package driver is Redshift's implementation of the driver contract in
// github.com/bytebase/omni/driver: it connects through pgx, syncs a server's
// metadata into omni/metadata, and runs interactive queries and change scripts.
//
// # Recorded from a cluster
//
// The sync, query, execute and row-estimate paths replay tapes recorded from a
// single-node rg.large cluster running Redshift 1.0.539507 (PostgreSQL 8.0.2)
// in us-west-2 on 2026-11-17, and testdata/sync/dev.json holds the metadata the
// sync built from the corpus in testdata/corpus.sql. A replay starts no server
// and needs no credentials; a recording (OMNI_RECORD=1) reads REDSHIFT_DSN and
// testcontainer fails without it rather than write an empty tape. The pin row
// in redshift/testdata/pins.tsv names the release and hashes SELECT version()'s
// answer, so a cluster of another release is refused instead of replayed.
// VERIFY.md records what the recording found and what remains unverified.
//
// The datashare branch has no tape: an inbound share needs a second account
// that publishes one, which the scratch cluster did not have. Its test skips,
// and the sync queries it sends (SVV_ALL_SCHEMAS, SVV_ALL_TABLES and
// SVV_ALL_COLUMNS) were only read against an ordinary database.
//
// # What it implements
//
// It implements Driver, Syncer, Querier and Executor, and the row estimate
// CountAffectedRows that a review's affected-row limit reads. There is no
// Dumper: although redshift/catalog provides DatabaseDefinition for rendering
// synced snapshots, SyncDatabase here only populates tables and ordinary views
// (not functions, procedures, materialized views, or storage attributes such as
// DISTSTYLE/DISTKEY/SORTKEY/ENCODE), so advertising driver.Dumper before sync is
// extended and re-recorded against Redshift would produce incomplete dumps.
//
// # The connection
//
// Config.Host and Config.Port are libpq's host and port lists, as pg's driver
// takes them. Config.Database is the database the driver connects to and
// SyncDatabase reads; empty means Redshift's default, dev. Config.TLS makes TLS
// required and supplies the configuration for every host pgx may dial; nil
// leaves libpq's default negotiation, sslmode=prefer. Config.Dial opens the
// network connection in place of net.Dialer, so an SSH tunnel or a cloud
// connector reaches the cluster, and a Dial that names its own endpoint may
// leave Host and Port empty. Config.Params are libpq and pgx connection
// keywords, applied last so they override the settings above; ParamDatashare is
// this driver's own flag and is not sent to the server.
//
// Params[ParamDatashare] = "true" marks a database a datashare (consumer)
// database received. Redshift refuses to write one, so Execute refuses the
// script instead of sending it; the sync reads SVV_ALL_SCHEMAS, SVV_ALL_TABLES
// and SVV_ALL_COLUMNS, which name the objects inside a shared database, rather
// than the pg_catalog views a datashare does not publish; and an interactive
// statement has the database's own name stripped from the relations its FROM
// and JOIN clauses name, because a shared object is addressed without it. Config.ReadOnly asks for a
// read-only session through default_transaction_read_only, and a datashare
// database leaves that unset: Redshift does not support read-only transactions
// on one.
//
// Config.EngineVersion does nothing: SyncInstance asks the server for its
// version and reads nothing version-dependent afterwards.
//
// # Queries and change scripts
//
// A Session is one server session, so what a statement sets there (the search
// path, a temporary table) holds for the session's next statement.
// QueryOptions.Schema sets the search path before each statement of the call
// rather than once per session, so a session used for two schemas does not
// carry the first into the second.
//
// QueryOptions.Explain runs EXPLAIN around each statement, in the one plan
// format Redshift has: Result.Plan carries Format TEXT and Executed false.
// QueryOptions.ExplainFormat does nothing. QueryOptions.Limit rewrites each
// SELECT through redshift/analysis.StatementWithResultLimit, which wraps what
// it cannot splice a LIMIT into. QueryOptions.MaskingEnabled does nothing: a
// statement that returns rows here is one the analysis traced, as in pg.
// QueryOptions.MaxResultBytes counts a result's size as its rows are read and
// sets Result.Truncated.
//
// Execute ships the script's statements as the engine's splitter finds them,
// bound for one transaction unless ExecuteOptions.Transaction is AutoCommit, in
// which case each statement commits itself and the first failure leaves the
// rest of the script unrun. ExecuteOptions.CreateDatabase runs the script
// outside any transaction and is idempotent: a database the server already
// lists is not created again. ExecuteOptions.Isolation and
// ExecuteOptions.MaxRetries do nothing: Redshift's transactions are
// SERIALIZABLE and its driver never retried a serialization failure.
package driver

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
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
	_ driver.Querier  = (*Driver)(nil)
	_ driver.Executor = (*Driver)(nil)
)

// ParamDatashare is the Config.Params key that marks a datashare database:
// "true" makes the driver read the SVV_ALL_* views, strip the database name
// from interactive statements, and refuse Execute.
const ParamDatashare = "datashare"

// defaultDatabase is the database a Config with none connects to, Redshift's
// own default.
const defaultDatabase = "dev"

// Driver is an open Redshift driver: a pool of connections to one database.
type Driver struct {
	db           *sql.DB
	databaseName string
	datashare    bool
}

// Open returns a driver for the server and database cfg names. Open does not
// reach the server; Ping does.
func Open(_ context.Context, cfg driver.Config) (*Driver, error) {
	connConfig, err := connConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &Driver{
		db:           stdlib.OpenDB(*connConfig),
		databaseName: connConfig.Database,
		datashare:    isDatashare(cfg),
	}, nil
}

// isDatashare reports whether cfg marks its database as a datashare.
func isDatashare(cfg driver.Config) bool {
	switch strings.ToLower(strings.TrimSpace(cfg.Params[ParamDatashare])) {
	case "1", "t", "true", "y", "yes":
		return true
	default:
		return false
	}
}

// connConfig is the pgx configuration Open connects with.
//
// Host and Port are libpq's host and port lists. TLS, when set, makes TLS
// required (sslmode=require) and supplies the configuration for every host pgx
// may dial, which no Params sslmode can turn off or add a plaintext attempt to;
// when nil, libpq's default sslmode=prefer applies. Params are connection
// string keywords with values in libpq syntax, quoted by the caller where they
// need it, applied last so they override the above. Unless they name them,
// queries run with default_query_exec_mode=exec and no statement cache, which
// connection poolers in front of a cluster need.
func connConfig(cfg driver.Config) (*pgx.ConnConfig, error) {
	if cfg.User == "" && cfg.Params["user"] == "" {
		return nil, fmt.Errorf("user must be set")
	}
	var dsn []string
	// A cloud connector's Dial names its instance itself and needs no host.
	if cfg.Dial == nil || cfg.Host != "" || cfg.Port != "" {
		// pgx gives an empty host its default, a Unix socket directory on
		// this machine, and an empty port 5432, which is not a cluster's.
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
	for key, value := range cfg.Params {
		if key == ParamDatashare {
			// A startup parameter the server does not know is a connection
			// error, and this flag is read here rather than by Redshift.
			continue
		}
		dsn = append(dsn, key+"="+value)
	}
	connConfig, err := pgx.ParseConfig(strings.Join(dsn, " "))
	if err != nil {
		return nil, err
	}
	// Params are applied last, so a user or password they name is the login the
	// connection uses — as pgx parsed it, since the value is libpq text.
	// Config's are the next choice.
	if _, ok := cfg.Params["user"]; !ok {
		connConfig.User = cfg.User
	}
	if _, ok := cfg.Params["password"]; !ok {
		connConfig.Password = cfg.Password
	}
	// Params are applied last, so a dbname they name is the database the
	// connection uses — as pgx parsed it, since the value is libpq text.
	// Config.Database is the next choice, and Redshift's own default the last.
	if _, ok := cfg.Params["dbname"]; ok {
		// connConfig.Database already holds what the parameter named.
	} else if cfg.Database != "" {
		connConfig.Database = cfg.Database
	} else {
		connConfig.Database = defaultDatabase
	}
	if cfg.TLS != nil {
		if connConfig.TLSConfig, err = cfg.TLS(connConfig.Host); err != nil {
			return nil, err
		}
		// sslmode=prefer or allow in Params adds a plaintext attempt per host,
		// and disable makes every host plaintext; each host is dialed once,
		// with the caller's TLS.
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
	// A datashare database has no read-only transactions to ask for.
	if cfg.ReadOnly && !isDatashare(cfg) {
		connConfig.RuntimeParams["default_transaction_read_only"] = "true"
	}
	return connConfig, nil
}

// quoteValue quotes a libpq connection string value, so that a value holding a
// space cannot name another keyword. Surrounding whitespace is dropped, as
// libpq drops it from an unquoted value.
func quoteValue(value string) string {
	return "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(strings.TrimSpace(value)) + "'"
}

// Ping reaches the server.
func (d *Driver) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// Close closes the driver's connections.
func (d *Driver) Close() error {
	return d.db.Close()
}

// getDatabases reads every database of the server, system ones included, and
// marks the ones an inbound datashare made. Datashares are read first so that
// the consumer databases are known before the database list is scanned.
func (d *Driver) getDatabases(ctx context.Context) ([]*metadatapb.DatabaseSchemaMetadata, error) {
	consumerDatabases, err := d.getConsumerDatabases(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := d.db.QueryContext(ctx, `
		SELECT datname,
		pg_encoding_to_char(encoding)
		FROM pg_database
		ORDER BY datname;
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to execute databases query: %w", err)
	}
	defer rows.Close()

	var databases []*metadatapb.DatabaseSchemaMetadata
	for rows.Next() {
		database := metadatapb.DatabaseSchemaMetadata{}
		if err := rows.Scan(&database.Name, &database.CharacterSet); err != nil {
			return nil, err
		}
		database.Datashare = consumerDatabases[database.Name]
		databases = append(databases, &database)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return databases, nil
}

// getConsumerDatabases reads the databases an inbound datashare created.
func (d *Driver) getConsumerDatabases(ctx context.Context) (map[string]bool, error) {
	const query = `SELECT consumer_database FROM SVV_DATASHARES WHERE share_type = 'INBOUND' ORDER BY consumer_database;`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute datashares query: %w", err)
	}
	defer rows.Close()

	consumerDatabases := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		consumerDatabases[name] = true
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return consumerDatabases, nil
}

// getVersion reads the server's version string. Redshift has no
// server_version setting to ask for, so it parses what SELECT version()
// reports, which names both the PostgreSQL release Redshift speaks and its own.
func (d *Driver) getVersion(ctx context.Context) (string, error) {
	const query = "SELECT version()"
	var version string
	if err := d.db.QueryRowContext(ctx, query).Scan(&version); err != nil {
		return "", fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	// The parse is best-effort: a version string Redshift words differently is
	// still the server's answer, and reporting it unparsed beats failing a sync.
	pgVersion, redshiftVersion, err := getPgVersionAndRedshiftVersion(version)
	if err != nil {
		return version, nil
	}
	return buildRedshiftVersionString(redshiftVersion, pgVersion), nil
}

// versionRegexp captures the PostgreSQL release and the Redshift release from
// SELECT version(), which reads
// "PostgreSQL 8.0.2 on i686-pc-linux-gnu, compiled by GCC gcc (GCC) 3.4.2 20041017 (Red Hat 3.4.2-6.fc3), Redshift 1.0.48042".
var versionRegexp = regexp.MustCompile(`(?i)^PostgreSQL (?P<pgVersion>\d+\.\d+\.\d+) on .*, Redshift (?P<redshiftVersion>\d+\.\d+\.\d+)`)

// getPgVersionAndRedshiftVersion parses SELECT version() into the PostgreSQL
// release Redshift speaks and the Redshift release itself.
func getPgVersionAndRedshiftVersion(version string) (string, string, error) {
	matches := versionRegexp.FindStringSubmatch(version)
	if len(matches) == 0 {
		return "", "", fmt.Errorf("unable to parse version string: %s", version)
	}

	var pgVersion, redshiftVersion string
	for i, name := range versionRegexp.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		switch name {
		case "pgVersion":
			pgVersion = matches[i]
		case "redshiftVersion":
			redshiftVersion = matches[i]
		}
	}
	return pgVersion, redshiftVersion, nil
}

// buildRedshiftVersionString writes the version as
// "Redshift <redshiftVersion> based on PostgreSQL <postgresVersion>".
func buildRedshiftVersionString(redshiftVersion, postgresVersion string) string {
	return "Redshift " + redshiftVersion + " based on PostgreSQL " + postgresVersion
}

// excludedDatabaseList names the databases a sync leaves out: the ones a cloud
// provider keeps for itself and PostgreSQL's templates.
var excludedDatabaseList = map[string]bool{
	// Skip internal databases from cloud service providers
	// aws
	"padb_harvest":   true,
	"awsdatacatalog": true,
	"sys:internal":   true,
	// system templates.
	"template0": true,
	"template1": true,
}
