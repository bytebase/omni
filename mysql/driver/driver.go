// Package driver is MySQL's implementation of the driver contract in
// github.com/bytebase/omni/driver: it connects through go-sql-driver/mysql,
// syncs a server's metadata into omni/metadata, and dumps it as DDL.
//
// MySQL's wire handshake has no read-only session parameter, so
// driver.Config.ReadOnly is ignored.
package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"regexp"
	"slices"
	"strings"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/driver"
)

var (
	_ driver.Driver = (*Driver)(nil)
	_ driver.Syncer = (*Driver)(nil)
	_ driver.Dumper = (*Driver)(nil)
)

// Driver is an open MySQL driver: a pool of connections to one server and
// optional default database.
type Driver struct {
	db            *sql.DB
	databaseName  string
	engineVersion string
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

	params := []string{"multiStatements=true", "maxAllowedPacket=0"}
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

// Ping reaches the server.
func (d *Driver) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// Close closes the driver's connections.
func (d *Driver) Close() error {
	return d.db.Close()
}

var versionPrefixRe = regexp.MustCompile(`^\d+\.\d+\.\d+`)

func (d *Driver) getVersion(ctx context.Context) (string, string, error) {
	const query = "SELECT VERSION()"
	var raw string
	if err := d.db.QueryRowContext(ctx, query).Scan(&raw); err != nil {
		return "", "", fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return parseVersion(raw)
}

func parseVersion(version string) (string, string, error) {
	loc := versionPrefixRe.FindStringIndex(version)
	if loc == nil {
		return "", "", fmt.Errorf("failed to parse version %q", version)
	}
	return version[loc[0]:loc[1]], version[loc[1]:], nil
}

func (d *Driver) getServerVariable(ctx context.Context, varName string) (string, error) {
	query := fmt.Sprintf("SHOW VARIABLES LIKE '%s'", strings.ReplaceAll(varName, "'", "''"))
	var found, value string
	if err := d.db.QueryRowContext(ctx, query).Scan(&found, &value); err != nil {
		return "", fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if varName != found {
		return "", fmt.Errorf("expecting variable %s, but got %s", varName, found)
	}
	return value, nil
}

func (d *Driver) getInstanceRoles(ctx context.Context) []*driver.Role {
	users, err := d.getUsersFromMySQLUser(ctx)
	if err != nil {
		users, err = d.getUsersFromUserAttributes(ctx)
		if err != nil {
			return nil
		}
	}

	var roles []*driver.Role
	for _, u := range users {
		grants, err := d.getGrantFromUser(ctx, u.quoted)
		if err != nil {
			continue
		}
		attr := strings.Join(grants, "\n")
		roles = append(roles, &driver.Role{
			Name:      u.display,
			Attribute: &attr,
		})
	}
	return roles
}

type mysqlAccount struct {
	display string
	quoted  string
}

func formatAccount(user, host string) mysqlAccount {
	return mysqlAccount{
		display: fmt.Sprintf("'%s'@'%s'", user, host),
		quoted:  fmt.Sprintf("%s@%s", quoteSingle(user), quoteSingle(host)),
	}
}

func (d *Driver) getGrantFromUser(ctx context.Context, quotedAccount string) ([]string, error) {
	query := "SHOW GRANTS FOR " + quotedAccount
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query grants for %s: %w", quotedAccount, err)
	}
	defer rows.Close()

	grants := []string{}
	for rows.Next() {
		var grant string
		if err := rows.Scan(&grant); err != nil {
			return nil, fmt.Errorf("failed to scan grants for %s: %w", quotedAccount, err)
		}
		grants = append(grants, grant)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate grants for %s: %w", quotedAccount, err)
	}
	return grants, nil
}

func (d *Driver) getUsersFromUserAttributes(ctx context.Context) ([]mysqlAccount, error) {
	const query = `
	SELECT
		user,
		host
	FROM information_schema.user_attributes
	WHERE user <> ''
		AND user NOT LIKE 'mysql.%'
	ORDER BY user, host
	`
	return d.scanAccounts(ctx, query, "information_schema.user_attributes")
}

func (d *Driver) getUsersFromMySQLUser(ctx context.Context) ([]mysqlAccount, error) {
	const query = `
	SELECT
		user,
		host
	FROM mysql.user
	WHERE user <> ''
		AND user NOT LIKE 'mysql.%'
	ORDER BY user, host
	`
	return d.scanAccounts(ctx, query, "mysql.user")
}

func (d *Driver) scanAccounts(ctx context.Context, query, table string) ([]mysqlAccount, error) {
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", table, err)
	}
	defer rows.Close()

	var accounts []mysqlAccount
	for rows.Next() {
		var user, host string
		if err := rows.Scan(&user, &host); err != nil {
			return nil, err
		}
		accounts = append(accounts, formatAccount(user, host))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate %s: %w", table, err)
	}
	return accounts, nil
}

func quoteBacktick(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

func quoteSingle(s string) string {
	escaped := strings.ReplaceAll(s, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `'`, `''`)
	return "'" + escaped + "'"
}
