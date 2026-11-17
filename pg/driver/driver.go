// Package driver is PostgreSQL's implementation of the driver contract in
// github.com/bytebase/omni/driver: it connects through pgx, syncs a
// server's metadata into omni/metadata, dumps it as DDL, and runs interactive
// queries and change scripts.
//
// Result.AffectedRows is 0 for a statement that reports no count, such as DDL
// or SET: database/sql states a count for every statement and cannot tell a
// statement that has none from one that changed no rows.
package driver

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	pganalysis "github.com/bytebase/omni/pg/analysis"
)

var (
	_ driver.Driver   = (*Driver)(nil)
	_ driver.Syncer   = (*Driver)(nil)
	_ driver.Dumper   = (*Driver)(nil)
	_ driver.Querier  = (*Driver)(nil)
	_ driver.Executor = (*Driver)(nil)
)

// Options are PostgreSQL-specific settings for Open.
type Options struct {
	// TenantMode runs statements as the database's owner role rather than
	// the login role, switching back to superuser only for statements
	// PostgreSQL refuses a non-superuser (GRANT, CREATE EXTENSION, CREATE
	// EVENT TRIGGER, COMMENT ON EXTENSION/EVENT TRIGGER).
	TenantMode bool
}

// Driver is an open PostgreSQL driver: a pool of connections to one
// database.
type Driver struct {
	db            *sql.DB
	databaseName  string
	engineVersion string
	tenantMode    bool

	noticesMu sync.Mutex
	notices   map[*pgconn.PgConn]*[]driver.Message
}

// defaultDatabase is the database a Config with none connects to.
const defaultDatabase = "postgres"

// Open returns a driver for the server and database cfg names. Unless
// Options.TenantMode is set, Open does not reach the server; Ping does.
func Open(ctx context.Context, cfg driver.Config, opts ...Options) (*Driver, error) {
	var opt Options
	if len(opts) > 0 {
		opt = opts[0]
	}
	pgxConfig, err := connConfig(cfg)
	if err != nil {
		return nil, err
	}
	d := &Driver{
		databaseName:  pgxConfig.Database,
		engineVersion: cfg.EngineVersion,
		tenantMode:    opt.TenantMode,
		notices:       make(map[*pgconn.PgConn]*[]driver.Message),
	}
	pgxConfig.OnNotice = d.onNotice
	var openOpts []stdlib.OptionOpenDB
	if d.tenantMode {
		openOpts = append(openOpts, stdlib.OptionAfterConnect(func(ctx context.Context, conn *pgx.Conn) error {
			var owner string
			if err := conn.QueryRow(ctx, currentDatabaseOwnerQuery).Scan(&owner); err != nil {
				return fmt.Errorf("failed to get database owner: %w", err)
			}
			if _, err := conn.Exec(ctx, fmt.Sprintf("SET ROLE %s;", quoteIdent(owner))); err != nil {
				return fmt.Errorf("failed to set role to database owner %q: %w", owner, err)
			}
			return nil
		}))
	}
	d.db = stdlib.OpenDB(*pgxConfig, openOpts...)
	if d.tenantMode {
		if err := d.db.PingContext(ctx); err != nil {
			_ = d.db.Close()
			return nil, err
		}
	}
	return d, nil
}

// connConfig is the pgx configuration Open connects with.
//
// Host and Port are libpq's host and port lists. TLS, when set, makes TLS
// required (sslmode=require) and supplies the configuration for every host
// pgx may dial, which no Params sslmode can turn off or add a plaintext
// attempt to; when nil, libpq's
// default sslmode=prefer applies. Params are connection string keywords with
// values in libpq syntax, quoted by the caller where they need it, applied
// last so they override the above. Unless they name them, queries run with default_query_exec_mode=exec and no
// statement cache, which connection poolers such as PgBouncer need.
func connConfig(cfg driver.Config) (*pgx.ConnConfig, error) {
	if cfg.User == "" {
		return nil, fmt.Errorf("user must be set")
	}
	var dsn []string
	// A cloud connector's Dial names its instance itself and needs no host.
	if cfg.Dial == nil || cfg.Host != "" || cfg.Port != "" {
		// pgx gives an empty host its default, a Unix socket directory on
		// this machine, and an empty port 5432.
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
	if cfg.ReadOnly {
		connConfig.RuntimeParams["default_transaction_read_only"] = "true"
	}
	return connConfig, nil
}

// quoteValue quotes a libpq connection string value, so that a value
// holding a space cannot name another keyword. Surrounding whitespace is
// dropped, as libpq drops it from an unquoted value.
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

// getDatabases reads every database of the server, system ones included,
// in name order: pg_database is shared by every database, so its heap order
// moves as databases come and go.
func (d *Driver) getDatabases(ctx context.Context) ([]*metadatapb.DatabaseSchemaMetadata, error) {
	var databases []*metadatapb.DatabaseSchemaMetadata
	rows, err := d.db.QueryContext(ctx, "SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		database := &metadatapb.DatabaseSchemaMetadata{}
		if err := rows.Scan(&database.Name, &database.CharacterSet, &database.Collation, &database.Owner); err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return databases, nil
}

// SearchPath reads the session's search_path as written, a comma-separated
// list of schema names.
func (d *Driver) SearchPath(ctx context.Context) (string, error) {
	const query = "SHOW search_path"
	var searchPath string
	if err := d.db.QueryRowContext(ctx, query).Scan(&searchPath); err != nil {
		return "", fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return strings.TrimSpace(searchPath), nil
}

const currentDatabaseOwnerQuery = `
		SELECT
			u.rolname
		FROM
			pg_roles AS u JOIN pg_database AS d ON (d.datdba = u.oid)
		WHERE
			d.datname = current_database();
		`

// CurrentDatabaseOwner returns the owner role of the current database.
func (d *Driver) CurrentDatabaseOwner(ctx context.Context) (string, error) {
	var owner string
	if err := d.db.QueryRowContext(ctx, currentDatabaseOwnerQuery).Scan(&owner); err != nil {
		return "", err
	}
	return owner, nil
}

// CurrentUser returns the session's effective role (CURRENT_USER).
func (d *Driver) CurrentUser(ctx context.Context) (string, error) {
	var user string
	if err := d.db.QueryRowContext(ctx, "SELECT CURRENT_USER").Scan(&user); err != nil {
		return "", err
	}
	return user, nil
}

func (d *Driver) onNotice(pc *pgconn.PgConn, n *pgconn.Notice) {
	if pc == nil || n == nil {
		return
	}
	d.noticesMu.Lock()
	defer d.noticesMu.Unlock()
	if buf, ok := d.notices[pc]; ok {
		*buf = append(*buf, driver.Message{
			Level:   noticeLevel(n.Severity),
			Content: n.Message,
		})
	}
}

func noticeLevel(severity string) string {
	switch severity {
	case "DEBUG", "INFO", "NOTICE", "WARNING", "EXCEPTION", "LOG":
		return severity
	default:
		return ""
	}
}

func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func quoteLiteral(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// getVersion reads the server version as major.minor.patch from
// server_version_num, which writes 17.12 as 170012 and 9.6.24 as 90624.
func (d *Driver) getVersion(ctx context.Context) (string, error) {
	const query = "SHOW server_version_num"
	var version string
	if err := d.db.QueryRowContext(ctx, query).Scan(&version); err != nil {
		return "", fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	versionNum, err := strconv.Atoi(version)
	if err != nil {
		return "", err
	}
	major, minor, patch := versionNum/1_00_00, (versionNum/100)%100, versionNum%100
	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}

// getInstanceRoles reads the server's roles other than PostgreSQL's own,
// with their attributes written as psql's \du writes them.
func (d *Driver) getInstanceRoles(ctx context.Context) ([]*driver.Role, error) {
	const query = `
		SELECT r.rolname, r.rolsuper, r.rolinherit, r.rolcreaterole, r.rolcreatedb, r.rolcanlogin, r.rolreplication, r.rolvaliduntil, r.rolbypassrls
		FROM pg_catalog.pg_roles r
		WHERE r.rolname !~ '^pg_'
		ORDER BY r.rolname;
	`
	var roles []*driver.Role
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()
	for rows.Next() {
		var role string
		var super, inherit, createRole, createDB, canLogin, replication, bypassRLS bool
		var validUntil sql.NullString
		if err := rows.Scan(&role, &super, &inherit, &createRole, &createDB, &canLogin, &replication, &validUntil, &bypassRLS); err != nil {
			return nil, err
		}
		if pganalysis.IsSystemUser(role) {
			continue
		}
		var attributes []string
		if super {
			attributes = append(attributes, "Superuser")
		}
		if !inherit {
			attributes = append(attributes, "No inheritance")
		}
		if createRole {
			attributes = append(attributes, "Create role")
		}
		if createDB {
			attributes = append(attributes, "Create DB")
		}
		if !canLogin {
			attributes = append(attributes, "Cannot login")
		}
		if replication {
			attributes = append(attributes, "Replication")
		}
		if validUntil.Valid {
			attributes = append(attributes, fmt.Sprintf("Password valid until %s", validUntil.String))
		}
		if bypassRLS {
			attributes = append(attributes, "Bypass RLS+")
		}
		attribute := strings.Join(attributes, " ")
		roles = append(roles, &driver.Role{Name: role, Attribute: &attribute})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}
