// Package driver is SQL Server's implementation of the driver contract in
// github.com/bytebase/omni/driver: it connects through go-mssqldb, syncs a
// server's metadata into omni/metadata, and runs interactive queries and
// change scripts.
//
// It implements Driver, Syncer, Dumper, Querier, Executor, planner row estimates,
// and the live probes the reviews run. The catalog owns SQL Server DDL
// rendering; Dump passes synced metadata through it.
//
// A batch is sent as one request, so what SQL Server scopes to it holds across
// its statements. SQL Server reports a batch's per-statement row counts only
// through the *sqlexp.ReturnMessage message loop: MsgRowsAffected arrives
// interleaved with MsgNext and MsgNextResultSet. A tape holds a request's
// result sets in order and one affected count for an exec, and it can carry
// neither the message object the counts are delivered on nor their order
// beside the result sets, so only a lone statement's count is read. A batch's
// row-count statements therefore carry no AffectedRows.
//
// A statement that reports no count at all (DDL, SET, SELECT @v = ...) has a
// Result with no AffectedRows, rather than a count of zero.
//
// QueryOptions.Explain turns the session's plan output on for the statement, so
// a plan is the statement's rows: Result.Plan names the output the request
// asked for (XML for SHOWPLAN_XML, TEXT for SHOWPLAN_ALL) and records that the
// statement did not run.
//
// Config.TLS supplies the TLS configuration for one host, which the driver
// hands to the client as is; a certificate authority is material the client
// reads from a path, so a caller that has one passes it in
// Params["certificate"] the way the client expects. Config.ReadOnly is
// ignored: the driver does not open a read-only session.
package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	gomssqldb "github.com/microsoft/go-mssqldb"
	"github.com/microsoft/go-mssqldb/msdsn"

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

// defaultPort is the port a Config with none connects to.
const defaultPort = "1433"

// Driver is an open SQL Server driver: a pool of connections to one database.
type Driver struct {
	db            *sql.DB
	databaseName  string
	engineVersion string

	dialMu       sync.Mutex
	dialProtocol string
	dialClose    func()
}

// Open returns a driver for the server and database cfg names. It does not
// reach the server; Ping does.
func Open(_ context.Context, cfg driver.Config) (*Driver, error) {
	d := &Driver{
		databaseName:  cfg.Database,
		engineVersion: cfg.EngineVersion,
	}
	config, err := d.mssqlConfig(cfg)
	if err != nil {
		d.closeDialer()
		return nil, err
	}
	db := sql.OpenDB(gomssqldb.NewConnectorConfig(config))
	db.SetConnMaxLifetime(2 * time.Hour)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(15)
	d.db = db
	return d, nil
}

// mssqlConfig builds the client's configuration from cfg. The credentials and
// the endpoint are set on the parsed configuration rather than written into the
// DSN, so a value holding a query separator or a payload cannot name another
// setting.
func (d *Driver) mssqlConfig(cfg driver.Config) (msdsn.Config, error) {
	if cfg.User == "" {
		return msdsn.Config{}, errors.New("user must be set")
	}
	host := strings.TrimSpace(cfg.Host)
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}
	instance := ""
	if before, after, ok := strings.Cut(host, `\`); ok {
		host, instance = before, after
	}
	port := strings.TrimSpace(cfg.Port)
	if cfg.Dial == nil {
		if host == "" {
			return msdsn.Config{}, errors.New("host must be set")
		}
		if port == "" {
			port = defaultPort
		}
	}

	query := url.Values{}
	// Older SQL Server versions offer only TLS 1.0.
	if _, ok := lookupParam(cfg.Params, msdsn.TLSMin); !ok {
		query.Set(msdsn.TLSMin, "1.0")
	}
	if _, ok := lookupParam(cfg.Params, msdsn.TrustServerCertificate); !ok {
		// Without TLS material the client accepts the server's certificate, as
		// Bytebase's driver does.
		query.Set(msdsn.TrustServerCertificate, "true")
	}
	keys := make([]string, 0, len(cfg.Params))
	for key := range cfg.Params {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	for _, key := range keys {
		query.Set(key, cfg.Params[key])
	}

	dsn := "sqlserver://" + net.JoinHostPort("placeholder", defaultPort) + "?" + query.Encode()
	mssqlConfig, err := msdsn.Parse(dsn)
	if err != nil {
		return msdsn.Config{}, err
	}
	mssqlConfig.Host = host
	mssqlConfig.Instance = instance
	if port != "" {
		value, err := strconv.ParseUint(port, 10, 16)
		if err != nil {
			return msdsn.Config{}, fmt.Errorf("invalid port %q", port)
		}
		mssqlConfig.Port = value
	}
	// The client read the rest of the DSN from a placeholder host, so the name
	// it verifies the server's certificate against is the one dialed.
	if mssqlConfig.TLSConfig != nil && host != "" {
		mssqlConfig.TLSConfig.ServerName = host
	}
	mssqlConfig.User = cfg.User
	mssqlConfig.Password = cfg.Password
	mssqlConfig.Database = cfg.Database
	if cfg.ApplicationName != "" {
		mssqlConfig.AppName = cfg.ApplicationName
	} else {
		mssqlConfig.AppName = "bytebase"
	}

	if cfg.TLS != nil {
		tlsConfig, err := cfg.TLS(host)
		if err != nil {
			return msdsn.Config{}, err
		}
		if tlsConfig == nil {
			return msdsn.Config{}, errors.New("TLS configuration must not be nil")
		}
		mssqlConfig.TLSConfig = tlsConfig
		if mssqlConfig.Encryption != msdsn.EncryptionStrict {
			mssqlConfig.Encryption = msdsn.EncryptionRequired
		}
	}

	if cfg.Dial != nil {
		d.registerDialer(cfg.Dial)
		mssqlConfig.Protocols = []string{d.dialProtocol}
	}
	return mssqlConfig, nil
}

// lookupParam reads a parameter by a case-insensitive key, as the client does.
func lookupParam(params map[string]string, name string) (string, bool) {
	for key, value := range params {
		if strings.EqualFold(key, name) {
			return value, true
		}
	}
	return "", false
}

// dialersMu guards the client's process-wide dialer registry, which every
// driver's registration and removal shares.
var dialersMu sync.Mutex

// registerDialer makes dial the client's transport. The client reaches a
// dialer through a process-wide registry keyed by protocol name, so each driver
// registers a name of its own and removes it on Close.
func (d *Driver) registerDialer(dial func(ctx context.Context, network, address string) (net.Conn, error)) {
	d.dialMu.Lock()
	defer d.dialMu.Unlock()
	if d.dialProtocol != "" {
		return
	}
	name := fmt.Sprintf("bytebase-dial-%p", d)
	dialersMu.Lock()
	msdsn.ProtocolDialers[name] = fixedDialer{dial: dial}
	dialersMu.Unlock()
	d.dialProtocol = name
	d.dialClose = func() {
		dialersMu.Lock()
		defer dialersMu.Unlock()
		delete(msdsn.ProtocolDialers, name)
	}
}

func (d *Driver) closeDialer() {
	d.dialMu.Lock()
	defer d.dialMu.Unlock()
	if d.dialClose != nil {
		d.dialClose()
		d.dialClose = nil
		d.dialProtocol = ""
	}
}

// fixedDialer is a client protocol dialer over a caller's dial function.
type fixedDialer struct {
	dial func(ctx context.Context, network, address string) (net.Conn, error)
}

// ParseBrowserData leaves the configuration as it is: a caller's dialer names
// its endpoint itself, so there is no instance for the SQL Browser to resolve.
func (fixedDialer) ParseBrowserData(msdsn.BrowserData, *msdsn.Config) error { return nil }

// DialConnection opens the connection with the caller's dial function.
func (f fixedDialer) DialConnection(ctx context.Context, config *msdsn.Config) (net.Conn, error) {
	return f.dial(ctx, "tcp", net.JoinHostPort(config.Host, strconv.FormatUint(config.Port, 10)))
}

// CallBrowser reports that the SQL Browser adds nothing to a caller's dialer.
func (fixedDialer) CallBrowser(*msdsn.Config) bool { return false }

// Ping reaches the server.
func (d *Driver) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// Close closes the driver's connections and drops the transport it registered.
func (d *Driver) Close() error {
	d.closeDialer()
	if d.db == nil {
		return nil
	}
	return d.db.Close()
}

// Instance reads the server version and its user databases.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	var version, fullVersion string
	if err := d.db.QueryRowContext(ctx, "SELECT SERVERPROPERTY('productversion'), @@VERSION").Scan(&version, &fullVersion); err != nil {
		return nil, err
	}
	for _, token := range strings.Fields(fullVersion) {
		if len(token) == 4 && strings.HasPrefix(token, "20") {
			version = fmt.Sprintf("%s (%s)", version, token)
			break
		}
	}

	const query = "SELECT name, collation_name FROM sys.databases WHERE name NOT IN ('master', 'model', 'msdb', 'tempdb', 'rdscore') ORDER BY name"
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	var databases []*metadatapb.DatabaseSchemaMetadata
	for rows.Next() {
		database := &metadatapb.DatabaseSchemaMetadata{}
		var collation sql.NullString
		if err := rows.Scan(&database.Name, &collation); err != nil {
			return nil, err
		}
		if collation.Valid {
			database.Collation = collation.String
		}
		databases = append(databases, database)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &driver.Instance{Version: version, Databases: databases}, nil
}
