// Package driver is StarRocks's implementation of the driver contract in
// github.com/bytebase/omni/driver: it connects through go-sql-driver/mysql,
// syncs a server's metadata into omni/metadata, and runs interactive queries
// and change scripts.
//
// There is no Dump: omni has no starrocks/catalog to render a synced schema as
// DDL. StarRocks has no read-only session parameter, so driver.Config.ReadOnly
// is ignored.
package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/driver"
)

var (
	_ driver.Driver   = (*Driver)(nil)
	_ driver.Syncer   = (*Driver)(nil)
	_ driver.Querier  = (*Driver)(nil)
	_ driver.Executor = (*Driver)(nil)
)

// Driver is an open StarRocks driver: a pool of connections to one server and
// optional default database.
type Driver struct {
	db            *sql.DB
	databaseName  string
	engineVersion string
	appName       string
	killOnCancel  bool
}

// Open returns a driver for the server and database cfg names. It does not
// reach the server; Ping does.
func Open(_ context.Context, cfg driver.Config) (*Driver, error) {
	mc, err := connConfig(cfg)
	if err != nil {
		return nil, err
	}
	connector, err := gomysql.NewConnector(mc)
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(connector)
	db.SetConnMaxLifetime(2 * time.Hour)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(15)
	return &Driver{
		db:            db,
		databaseName:  cfg.Database,
		engineVersion: cfg.EngineVersion,
		appName:       sanitizeCommentToken(cfg.ApplicationName),
		killOnCancel:  true,
	}, nil
}

func connConfig(cfg driver.Config) (*gomysql.Config, error) {
	if cfg.User == "" {
		return nil, errors.New("user must be set")
	}
	if err := validateExtraParams(cfg.Params); err != nil {
		return nil, err
	}

	host := strings.TrimSpace(cfg.Host)
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}
	port := strings.TrimSpace(cfg.Port)
	var netName, addr string
	switch {
	case cfg.Dial != nil:
		netName = "tcp"
		switch {
		case host != "" && port != "":
			addr = net.JoinHostPort(host, port)
		case host != "":
			addr = host
		default:
			addr = "dial"
		}
	case strings.HasPrefix(host, "/"):
		netName = "unix"
		addr = host
	default:
		if host == "" {
			return nil, errors.New("host must be set")
		}
		if port == "" {
			return nil, errors.New("port must be set")
		}
		netName = "tcp"
		addr = net.JoinHostPort(host, port)
	}

	// interpolateParams inlines arguments, because a tape holds no
	// driver-level prepared statement.
	params := []string{"multiStatements=true", "maxAllowedPacket=0", "interpolateParams=true"}
	keys := make([]string, 0, len(cfg.Params))
	for k := range cfg.Params {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		params = append(params, fmt.Sprintf("%s=%s", k, cfg.Params[k]))
	}

	mc, err := gomysql.ParseDSN("u@tcp(placeholder:3306)/?" + strings.Join(params, "&"))
	if err != nil {
		return nil, err
	}
	if mc.AllowAllFiles {
		return nil, errors.New("connection parameter \"allowAllFiles\" is not allowed for security reasons")
	}

	mc.User = cfg.User
	mc.Passwd = cfg.Password
	mc.DBName = cfg.Database
	mc.Net = netName
	mc.Addr = addr
	if cfg.Dial != nil {
		mc.DialFunc = cfg.Dial
	}
	if cfg.TLS == nil && mc.TLS != nil && mc.TLS.ServerName == "placeholder" {
		if netName == "tcp" {
			mc.TLS.ServerName = host
		} else {
			mc.TLS.ServerName = ""
		}
	}

	if cfg.TLS != nil {
		tlsCfg, err := cfg.TLS(host)
		if err != nil {
			return nil, err
		}
		if tlsCfg == nil {
			return nil, errors.New("TLS configuration must not be nil")
		}
		mc.TLS = tlsCfg
		mc.TLSConfig = ""
		mc.AllowFallbackToPlaintext = false
	}
	return mc, nil
}

func validateExtraParams(params map[string]string) error {
	for key := range params {
		if strings.EqualFold(strings.TrimSpace(key), "allowallfiles") {
			return fmt.Errorf("connection parameter %q is not allowed for security reasons", key)
		}
	}
	return nil
}

// sanitizeCommentToken keeps s inside the comment withAppComment opens: the
// terminator becomes "* /", which cannot form another one.
func sanitizeCommentToken(s string) string {
	return strings.ReplaceAll(s, "*/", "* /")
}

func (d *Driver) withAppComment(stmt string) string {
	if d.appName == "" {
		return stmt
	}
	return fmt.Sprintf("/*app=%s*/ %s", d.appName, stmt)
}

func (d *Driver) readConnectionID(ctx context.Context, conn *sql.Conn) (uint64, error) {
	if !d.killOnCancel {
		return 0, nil
	}
	var raw string
	if err := conn.QueryRowContext(ctx, "SELECT CONNECTION_ID();").Scan(&raw); err != nil {
		return 0, err
	}
	return strconv.ParseUint(raw, 10, 64)
}

func (d *Driver) killQueryOnCancel(err error, connID uint64) {
	if !d.killOnCancel || connID == 0 {
		return
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		killCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = d.db.ExecContext(killCtx, fmt.Sprintf("KILL QUERY %d", connID))
	}
}

// Ping reaches the server.
func (d *Driver) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// Close closes the driver's connections.
func (d *Driver) Close() error {
	return d.db.Close()
}

var (
	versionPrefixRe = regexp.MustCompile(`^\d+\.\d+\.\d+`)
	versionInTextRe = regexp.MustCompile(`\d+\.\d+\.\d+`)
)

// getVersion returns the engine release and the rest of the text naming it.
// VERSION() reports the MySQL wire-compatibility version — 5.7.99 for Doris
// 3.1.4 and 8.0.33 for StarRocks 3.4.10 — so @@version_comment names the
// release, and VERSION() is only the fallback.
func (d *Driver) getVersion(ctx context.Context) (string, string, error) {
	const query = "SELECT @@version_comment, VERSION()"
	var comment, version string
	if err := d.db.QueryRowContext(ctx, query).Scan(&comment, &version); err != nil {
		return "", "", fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if loc := versionInTextRe.FindStringIndex(comment); loc != nil {
		return comment[loc[0]:loc[1]], comment[loc[1]:], nil
	}
	return parseVersion(version)
}

// parseVersion splits a version string at its leading release triple.
func parseVersion(version string) (string, string, error) {
	loc := versionPrefixRe.FindStringIndex(version)
	if loc == nil {
		return "", "", fmt.Errorf("failed to parse version %q", version)
	}
	return version[loc[0]:loc[1]], version[loc[1]:], nil
}

// getServerVariable returns one server variable. The engines disagree on how
// many columns SHOW VARIABLES returns — Doris has four, StarRocks two — so the
// row is read by column name rather than by a fixed arity.
func (d *Driver) getServerVariable(ctx context.Context, varName string) (string, error) {
	query := fmt.Sprintf("SHOW VARIABLES LIKE '%s'", strings.ReplaceAll(varName, "'", "''"))
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	if len(columns) < 2 {
		return "", fmt.Errorf("query %q returned %d columns", query, len(columns))
	}
	valueIndex := 1
	for i, column := range columns {
		if strings.EqualFold(column, "Value") {
			valueIndex = i
			break
		}
	}
	values := make([]sql.NullString, len(columns))
	dest := make([]any, len(columns))
	for i := range values {
		dest[i] = &values[i]
	}
	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			return "", err
		}
		if !strings.EqualFold(values[0].String, varName) {
			continue
		}
		if !values[valueIndex].Valid {
			return "", fmt.Errorf("variable %s has no value", varName)
		}
		return values[valueIndex].String, nil
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("expecting variable %s, but the server reported none", varName)
}

func quoteBacktick(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}
