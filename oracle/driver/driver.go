// Package driver is Oracle's implementation of the driver contract in
// github.com/bytebase/omni/driver: it connects through go-ora, syncs a
// schema's metadata into omni/metadata, and runs interactive queries and change
// scripts.
//
// It implements Driver, Syncer, Querier, Executor, planner row estimates, and
// the live probes the reviews run. It does not implement Dumper: Oracle DDL
// rendering still lives in Bytebase's schema package until omni has an
// oracle/catalog.
//
// A Config.Database names the schema every session switches to
// (ALTER SESSION SET CURRENT_SCHEMA), which is how an Oracle connection selects
// what unqualified names resolve to. Config.ReadOnly is ignored: the driver does
// not open a read-only session.
package driver

import (
	"context"
	"crypto/tls"
	"database/sql"
	sqldriver "database/sql/driver"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	goora "github.com/sijms/go-ora/v2"

	omnidriver "github.com/bytebase/omni/driver"
	oracleanalysis "github.com/bytebase/omni/oracle/analysis"
)

var (
	_ omnidriver.Driver   = (*Driver)(nil)
	_ omnidriver.Syncer   = (*Driver)(nil)
	_ omnidriver.Querier  = (*Driver)(nil)
	_ omnidriver.Executor = (*Driver)(nil)

	_ omnidriver.PrivateDatabaseLinkOwner = (*Driver)(nil)
)

// Driver is an open Oracle driver: a pool of sessions to one pluggable
// database, each switched to the configured schema.
type Driver struct {
	db            *sql.DB
	databaseName  string
	serviceName   string
	engineVersion string
	// statementSeq numbers the EXPLAIN PLAN statement IDs a session records
	// plans under, which Oracle keeps per session.
	statementSeq atomic.Int64
	// schemaMu guards schemaName, which resolves the schema the driver reads when
	// Config.Database names none: the account's own.
	schemaMu   sync.Mutex
	schemaName string
}

// nextStatementID returns the next EXPLAIN PLAN statement ID. Bytebase picks a
// random one; a counter keeps the plan's statements the same from run to run,
// which a recording needs.
func (d *Driver) nextStatementID() string {
	return fmt.Sprintf("omni_%d", d.statementSeq.Add(1))
}

const (
	sidOption              = "SID"
	sslOption              = "SSL"
	dialTimeoutOption      = "CONNECTION TIMEOUT"
	dialTimeoutOptionAlias = "CONNECT TIMEOUT"
	// defaultDialTimeoutSeconds bounds the wait for one address to accept a
	// connection, DNS lookup included, so an unreachable address fails over in
	// seconds. It does not limit how long a statement runs.
	defaultDialTimeoutSeconds = "10"
	// defaultPort is the port the client reads from the URL when the caller
	// names none, which its own dialer makes irrelevant.
	defaultPort = 1521
)

// Open returns a driver for the server, service and schema cfg names. It reaches
// the server when Config.Database names a schema, so a schema the account cannot
// switch to fails here rather than at the first statement.
func Open(ctx context.Context, cfg omnidriver.Config) (*Driver, error) {
	// A caller-supplied dialer names the server itself, so the endpoint may be
	// empty; the client still parses a URL, and an unused port satisfies it.
	port := defaultPort
	if cfg.Port != "" {
		parsed, err := strconv.Atoi(strings.TrimSpace(cfg.Port))
		if err != nil {
			return nil, fmt.Errorf("invalid port %q", cfg.Port)
		}
		port = parsed
	}
	if strings.TrimSpace(cfg.Host) == "" && cfg.Dial == nil {
		return nil, errors.New("host must be set")
	}
	if cfg.User == "" {
		return nil, errors.New("user must be set")
	}
	service, err := serviceName(cfg.Params)
	if err != nil {
		return nil, err
	}

	var tlsConfig *tls.Config
	if cfg.TLS != nil {
		tlsConfig, err = cfg.TLS(strings.TrimSpace(cfg.Host))
		if err != nil {
			return nil, err
		}
		if tlsConfig == nil {
			return nil, errors.New("TLS configuration must not be nil")
		}
	}
	// The client refuses an empty host in its URL, and a caller whose dialer
	// names the server may leave it empty: the placeholder only fills the URL,
	// since the dialer replaces the transport that would read it.
	urlHost := strings.TrimSpace(cfg.Host)
	if urlHost == "" {
		urlHost = "localhost"
	}
	dsn := goora.BuildUrl(
		urlHost, port, service, cfg.User, cfg.Password,
		connectionOptions(cfg.Params, tlsConfig != nil),
	)
	d := &Driver{
		databaseName:  cfg.Database,
		serviceName:   service,
		engineVersion: cfg.EngineVersion,
	}
	d.db = sql.OpenDB(connector{
		drv:           goora.NewDriver(),
		dsn:           dsn,
		tlsConfig:     tlsConfig,
		dial:          dialFunc(cfg.Dial),
		currentSchema: cfg.Database,
	})
	if cfg.Database != "" {
		if err := d.db.PingContext(ctx); err != nil {
			_ = d.db.Close()
			return nil, err
		}
	}
	return d, nil
}

// connectionOptions returns the go-ora URL options: the caller's parameters as
// typed, the SID and the default dial timeout unless a parameter names the same
// option, and SSL=TRUE when TLS is on.
//
// go-ora upper-cases option keys and reads CONNECT TIMEOUT as CONNECTION
// TIMEOUT. Two keys for the same option would leave the winner to map iteration
// order, so only one is sent: the caller's for SID and the dial timeout, and the
// caller's TLS for SSL, which the TLS configuration decides.
func connectionOptions(params map[string]string, useTLS bool) map[string]string {
	options := make(map[string]string, len(params)+2)
	userOptions := make(map[string]bool, len(params))
	for key, value := range params {
		option := strings.ToUpper(key)
		if option == serviceOption {
			// The service names the endpoint, which the URL carries in its path.
			continue
		}
		if useTLS && option == sslOption {
			continue
		}
		options[key] = value
		userOptions[option] = true
	}
	if !userOptions[dialTimeoutOption] && !userOptions[dialTimeoutOptionAlias] {
		options[dialTimeoutOption] = defaultDialTimeoutSeconds
	}
	if useTLS {
		options[sslOption] = "TRUE"
	}
	return options
}

// serviceOption names the service in Params, since the endpoint a session
// reaches is the driver's host, port and service.
const serviceOption = "SERVICE_NAME"

// serviceName reads the service from Params, which is what a session's endpoint
// names; a SID in Params is the alternative go-ora accepts.
func serviceName(params map[string]string) (string, error) {
	for key, value := range params {
		if strings.EqualFold(key, serviceOption) {
			if strings.TrimSpace(value) == "" {
				return "", fmt.Errorf("%s must not be empty", serviceOption)
			}
			return value, nil
		}
	}
	for key := range params {
		if strings.EqualFold(key, sidOption) {
			return "", nil
		}
	}
	return "", fmt.Errorf("%s must be set, or a SID in Params", serviceOption)
}

// connector opens each session with its own copy of the TLS configuration:
// go-ora sets ServerName on the configuration it is given, so sessions opened in
// parallel cannot share one. It also sets the current schema on every session it
// opens, because the pool opens a new session whenever none is idle or one goes
// bad.
//
// drv is this driver's own go-ora driver, never the process-wide one: a go-ora
// driver keeps the character-set converters of the first connection opened
// through it and copies them onto every later connection, so a shared driver
// would decode every schema's text with the first schema's character set.
type connector struct {
	drv       *goora.OracleDriver
	dsn       string
	tlsConfig *tls.Config
	// dial opens the network connection in place of net.Dialer, for a caller
	// that reaches the server through a tunnel or a cloud connector.
	dial          dialFunc
	currentSchema string
}

// dialFunc is driver.Config.Dial, as go-ora's DialerContext.
type dialFunc func(ctx context.Context, network, address string) (net.Conn, error)

func (f dialFunc) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return f(ctx, network, address)
}

func (c connector) Connect(ctx context.Context) (sqldriver.Conn, error) {
	oracleConnector, err := c.drv.OpenConnector(c.dsn)
	if err != nil {
		return nil, err
	}
	// A nil configuration is the client's default negotiation.
	oracleConnector.(*goora.OracleConnector).WithTLSConfig(c.tlsConfig.Clone())
	if c.dial != nil {
		oracleConnector.(*goora.OracleConnector).Dialer(c.dial)
	}
	conn, err := oracleConnector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	if c.currentSchema == "" {
		return conn, nil
	}
	if _, err := conn.(sqldriver.ExecerContext).ExecContext(ctx, fmt.Sprintf("ALTER SESSION SET CURRENT_SCHEMA = %s", quoteIdent(c.currentSchema)), nil); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to set current schema to %q: %w", c.currentSchema, err)
	}
	return conn, nil
}

func (c connector) Driver() sqldriver.Driver { return c.drv }

// Ping reaches the server.
func (d *Driver) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// Close closes the driver's sessions.
func (d *Driver) Close() error {
	return d.db.Close()
}

// version returns the server's release, such as 23.26, and the configured
// engine version when the server does not answer.
func (d *Driver) version(ctx context.Context) (*oracleanalysis.Version, error) {
	const query = "SELECT version_full FROM product_component_version WHERE product LIKE 'Oracle%'"
	var fullVersion string
	if err := d.db.QueryRowContext(ctx, query).Scan(&fullVersion); err == nil {
		if version, parseErr := parseVersion(fullVersion); parseErr == nil {
			return version, nil
		}
	}
	if d.engineVersion == "" {
		return nil, errors.New("engine version is not set")
	}
	return parseVersion(d.engineVersion)
}

// OwnsPrivateDatabaseLink reports whether the driver's account owns a private
// database link: SYS.USER_DB_LINKS lists the links the current user owns, and
// the SYS qualifier keeps a same-named local object from standing in for the
// dictionary view.
func (d *Driver) OwnsPrivateDatabaseLink(ctx context.Context) (bool, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get connection: %w", err)
	}
	defer func() { _ = conn.Close() }()
	var count int
	if err := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM SYS.USER_DB_LINKS").Scan(&count); err != nil {
		return false, fmt.Errorf("failed to count private database links: %w", err)
	}
	return count > 0, nil
}

// Session opens one server session held until Close.
func (d *Driver) Session(ctx context.Context) (omnidriver.Session, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &session{driver: d, conn: conn}, nil
}

type session struct {
	driver *Driver
	conn   *sql.Conn
	// schemaApplied reports that QueryOptions.Schema switched the session, so
	// Close has a schema of the driver's to switch back to.
	schemaApplied bool
}

// Close returns the session to the pool with the driver's schema, so a schema a
// query chose does not become the next session's.
func (s *session) Close() error {
	if s.schemaApplied && s.driver.databaseName != "" {
		// The connection goes back to the pool either way, so a failure here is
		// the next session's to report.
		_, _ = s.conn.ExecContext(context.Background(), fmt.Sprintf("ALTER SESSION SET CURRENT_SCHEMA = %s", quoteIdent(s.driver.databaseName)))
	}
	return s.conn.Close()
}

// quoteIdent writes a name as a quoted identifier.
func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}
