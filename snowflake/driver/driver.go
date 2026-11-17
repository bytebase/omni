// Package driver is Snowflake's implementation of the driver contract in
// github.com/bytebase/omni/driver: it reaches an account through
// github.com/snowflakedb/gosnowflake, syncs an account's databases and a
// database's schemas into omni/metadata, and runs interactive queries and
// change scripts.
//
// # Unverified against an account
//
// No test in this package has ever asked a Snowflake account. The package is a
// port of Bytebase's backend/plugin/db/snowflake, whose tests reached only a
// private key decoder and the DSN builder, and Snowflake has no image omni can
// start in a container: every answer the sync, query, execute and dump paths
// expect here was written from that driver and from Snowflake's documentation,
// and none was recorded. There are no tapes, no sync goldens and no dump
// fixture beside this package, because nothing here can produce one. VERIFY.md
// holds the recording plan: the account or trial instance to pin, the tests to
// record in what order, and the parity harness that holds this package to the
// driver it replaces cell by cell. Treat a behavior this package states as
// unconfirmed until that recording lands.
//
// Its tests therefore cover only what needs no account: the client
// configuration built from a driver.Config, the private key decoder, the
// statement classification wiring, the value conversions and the dump's text
// transforms. Everything else is behind an environment gate and skips with the
// variable's name: SNOWFLAKE_DSN reaches an account for this package's own live
// tests, and SNOWFLAKE_COMPAT_DSN is the account omni's docs ask for when a
// compat corpus needs the engine to judge it.
//
// # What it implements
//
// It implements Driver, Syncer, Dumper, Querier and Executor. There is no row
// estimate: Snowflake's EXPLAIN states the operations of a plan and no
// cardinality for any of them, so CountAffectedRows has nothing to read and
// Bytebase's driver never answered one either.
//
// Dump is a server-DDL dump: it writes what GET_DDL('DATABASE', ...) states,
// with the statements that would recreate the database or its default PUBLIC
// schema and the database qualifier taken out, and ignores the metadata it is
// given. It needs no snowflake/catalog, and omni has none: a metadata-rendered
// dump would need a DDL renderer for a synced snapshot, which snowflake/deparse
// does not provide for one. VERIFY.md records what GET_DDL states before the
// transforms are trusted.
//
// # The connection
//
// Config.Host is the account, as Snowflake writes it in a login URL: xy12345,
// or nb47110.ap-southeast-1, whose suffix is the region. A proxy or tunnel
// endpoint is named after it, account@host, which is Bytebase's spelling for
// the same thing and is what Config.Dial is for. Config.Port is the HTTPS port,
// 443 when empty. Config.Database is the database the driver connects to and
// SyncDatabase reads; empty leaves the session without one, and SyncDatabase
// then finds no such database and reports driver.ErrNotFound.
//
// The credential arrives as values, as a cloud SDK engine's does.
// Config.User and Config.Password are a login. A key-pair login is
// Params[ParamPrivateKey], a PKCS#8 private key in PEM, with
// Params[ParamPrivateKeyPassphrase] for an encrypted one; a private key makes
// Config.Password ignored, as it did in Bytebase's driver, where the key-pair
// login carried no password. Config.User is required either way: Snowflake
// names the login in both.
//
// Params carry the settings the contract has no field for: ParamRegion,
// ParamWarehouse and ParamRole are the account's region and the session's
// warehouse and role, and every other name is passed to gosnowflake as one of
// its own connection parameters, such as timestamp_output_format.
//
// Config.TLS supplies the TLS settings for the endpoint and makes the driver
// build the HTTP transport the client uses; it is called once, with the
// endpoint's host (Config.Host when it named one, otherwise
// <account>.snowflakecomputing.com with the region in it). Config.Dial opens
// the TCP connection in place of a dialer, so an SSH tunnel or a proxy reaches
// the account. Both install a transport of the caller's, and gosnowflake
// checks no certificate revocation through a transport it did not build: the
// driver asks for OCSP fail-open checking only when it builds the transport
// itself. Config.ApplicationName is the session's client application name.
// Config.ReadOnly does nothing: Snowflake has no read-only session to open, and
// Bytebase's driver ignored the flag too. Config.EngineVersion does nothing:
// SyncInstance asks CURRENT_VERSION() and reads nothing version-dependent
// afterwards.
//
// # Queries and change scripts
//
// A Session is one client connection, held until Close, so a schema, a
// warehouse or a temporary table a statement sets holds for the session's next
// statement. QueryOptions.Schema runs USE SCHEMA before each statement of the
// call rather than once per session, and an empty Schema runs USE SCHEMA PUBLIC
// before each statement, ignoring the error, which is what Bytebase's editor did
// so that a session's earlier USE DATABASE does not decide where the next call
// creates its objects.
//
// QueryOptions.Explain runs EXPLAIN around each statement, or EXPLAIN USING JSON
// when ExplainFormat is JSON. A statement that is itself an EXPLAIN is reported
// as a plan, in the format snowflake/analysis reads from it. Snowflake's EXPLAIN
// compiles without running, so Plan.Executed is false.
// QueryOptions.Limit rewrites a query through snowflake/analysis.RewriteResultLimit
// and wraps the statement in SELECT * FROM (...) LIMIT n when the analysis
// cannot (a parse failure, OFFSET without LIMIT), which is what Bytebase did.
// QueryOptions.MaskingEnabled does nothing: a statement that returns rows here
// is one the analysis traced.
// QueryOptions.MaxResultBytes counts a result's size as its rows are read and
// sets Result.Truncated.
//
// A statement omni's parser cannot read runs as a query, as Bytebase ran it: the
// classification is a hint about how to send a statement, and the server decides
// whether the statement is valid. Every statement is sent whole, the ones a
// scripting block holds included: the splitter keeps BEGIN ... END as one
// statement and the client's one-statement limit is lifted for it, as Bytebase
// lifted it.
//
// Execute ships the script's statements as snowflake/analysis.SplitStatements
// finds them, one statement per request, in one transaction unless
// ExecuteOptions.Transaction is AutoCommit, in which case each statement commits
// itself and the first failure leaves the rest of the script unrun. A statement
// that fails ends the script. Bytebase sent the whole script as one
// multi-statement request, which reports one answer for all of it.
// ExecuteOptions.CreateDatabase changes nothing: Snowflake's CREATE DATABASE is
// an ordinary statement of the script, and Bytebase's create-database step ran it
// under the same default transaction mode this driver defaults to.
// ExecuteOptions.Isolation does nothing: Snowflake has no SET TRANSACTION and
// Bytebase's driver never asked for one. ExecuteOptions.MaxRetries does nothing:
// the driver never retried a failed script.
//
// A statement's error is the client's own, gosnowflake's
// *gosnowflake.SnowflakeError, which carries the server's error number and SQL
// state; a sync query's error is wrapped with the statement it sent.
package driver

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/snowflakedb/gosnowflake/v2"
	"github.com/youmark/pkcs8"

	"github.com/bytebase/omni/driver"
)

var (
	_ driver.Driver   = (*Driver)(nil)
	_ driver.Syncer   = (*Driver)(nil)
	_ driver.Dumper   = (*Driver)(nil)
	_ driver.Querier  = (*Driver)(nil)
	_ driver.Executor = (*Driver)(nil)
)

// The names this driver reads out of driver.Config.Params.
const (
	// ParamPrivateKey is a PKCS#8 private key in PEM, which makes the login a
	// key-pair one.
	ParamPrivateKey = "privateKey"
	// ParamPrivateKeyPassphrase is the passphrase of an encrypted
	// ParamPrivateKey.
	ParamPrivateKeyPassphrase = "privateKeyPassphrase"
	// ParamRegion is the account's region, for an account name that does not
	// carry one.
	ParamRegion = "region"
	// ParamWarehouse is the warehouse a statement runs in.
	ParamWarehouse = "warehouse"
	// ParamRole is the role the session runs as.
	ParamRole = "role"
)

// Driver is an open Snowflake driver: a pool of client connections to one
// account and database.
type Driver struct {
	db           *sql.DB
	databaseName string
}

// Open returns a driver for the account and database cfg names. Open does not
// reach the account; Ping does.
func Open(_ context.Context, cfg driver.Config) (*Driver, error) {
	config, err := snowConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &Driver{
		db:           sql.OpenDB(gosnowflake.NewConnector(gosnowflake.SnowflakeDriver{}, *config)),
		databaseName: cfg.Database,
	}, nil
}

// snowConfig is the gosnowflake configuration Open connects with.
//
// Config.Host names the account, or account@host when a proxy endpoint answers
// in front of it. Config.Database is written as a quoted identifier so a
// database whose name holds a lower-case letter is the one the caller named.
// A private key in Params makes the login a key-pair one and the password
// unused. Every other Params name goes to the client untouched.
func snowConfig(cfg driver.Config) (*gosnowflake.Config, error) {
	if strings.TrimSpace(cfg.User) == "" {
		return nil, errors.New("user must be set")
	}
	account, proxyHost, err := splitAccountHost(cfg.Host)
	if err != nil {
		return nil, err
	}

	config := &gosnowflake.Config{
		Account: account,
		User:    cfg.User,
		// gosnowflake v2.3.0 turned OCSP off unless a caller opts in; before it,
		// every connection checked and failed open. Ask for what Bytebase's
		// driver asked for, so a certificate revocation check is not dropped
		// without a caller's transport replacing it (see below).
		OCSPFailOpen: gosnowflake.OCSPFailOpenTrue,
		Application:  cfg.ApplicationName,
		Region:       cfg.Params[ParamRegion],
		Warehouse:    cfg.Params[ParamWarehouse],
		Role:         cfg.Params[ParamRole],
	}
	// A proxy endpoint is the host to dial; without one the client derives the
	// account's host from the account name and the region.
	if proxyHost != "" {
		config.Host = proxyHost
	}
	if cfg.Database != "" {
		// Snowflake stores an unquoted identifier upper-cased, so a database
		// named in lower case must be quoted to be found.
		config.Database = quoteIdent(cfg.Database)
	}
	if cfg.Port != "" {
		port, err := strconv.Atoi(cfg.Port)
		if err != nil {
			return nil, fmt.Errorf("port %q is not a number", cfg.Port)
		}
		config.Port = port
	}
	if privateKey := cfg.Params[ParamPrivateKey]; privateKey != "" {
		key, err := decodeRSAPrivateKey(privateKey, cfg.Params[ParamPrivateKeyPassphrase])
		if err != nil {
			return nil, err
		}
		config.PrivateKey = key
		config.Authenticator = gosnowflake.AuthTypeJwt
	} else {
		config.Password = cfg.Password
	}
	params := make(map[string]*string, len(cfg.Params))
	for key, value := range cfg.Params {
		switch key {
		case ParamPrivateKey, ParamPrivateKeyPassphrase, ParamRegion, ParamWarehouse, ParamRole:
			continue
		}
		value := value
		params[key] = &value
	}
	if len(params) > 0 {
		config.Params = params
	}

	if cfg.Dial != nil || cfg.TLS != nil {
		// The client takes a transport, not a dialer: the caller's connection
		// and TLS settings reach the endpoint through one built from the
		// default, so its timeouts and proxy settings stay.
		transport := http.DefaultTransport.(*http.Transport).Clone()
		if cfg.Dial != nil {
			transport.DialContext = cfg.Dial
		}
		if cfg.TLS != nil {
			host := proxyHost
			if host == "" {
				host = endpointHost(account, config.Region)
			}
			tlsConfig, err := cfg.TLS(host)
			if err != nil {
				return nil, err
			}
			transport.TLSClientConfig = tlsConfig
		}
		// The client checks no certificate revocation through a transport it
		// did not build, so OCSPFailOpen above stops applying here.
		config.Transporter = transport
	}
	return config, nil
}

// splitAccountHost splits Config.Host into the account and, when the account
// names a proxy endpoint after it, that host. Bytebase read the same two parts
// from one field.
func splitAccountHost(host string) (string, string, error) {
	host = strings.TrimSpace(host)
	if host == "" {
		return "", "", errors.New("host must be set, as the account name or account@host")
	}
	parts := strings.Split(host, "@")
	switch len(parts) {
	case 1:
		return parts[0], "", nil
	case 2:
		if strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
			return "", "", fmt.Errorf("host %q must name both an account and a host around @", host)
		}
		return parts[0], parts[1], nil
	default:
		return "", "", fmt.Errorf("expected one @ in the host at most, got %q", host)
	}
}

// endpointHost is the host a Snowflake client dials for an account: the account
// itself when it carries its region as a suffix (nb47110.ap-southeast-1), the
// account and the region when the region is named apart, and the account alone
// otherwise.
func endpointHost(account, region string) string {
	if region == "" {
		if dot := strings.Index(account, "."); dot > 0 {
			// The account carries its own region.
			return account + snowflakeDomain(account[dot+1:])
		}
		return account + snowflakeDomain("")
	}
	return account + "." + region + snowflakeDomain(region)
}

// snowflakeDomain is the domain an account of region answers under: a China
// region's accounts are named in .cn.
func snowflakeDomain(region string) string {
	if strings.HasPrefix(strings.ToLower(region), "cn-") {
		return ".snowflakecomputing.cn"
	}
	return ".snowflakecomputing.com"
}

// quoteIdent quotes name as a Snowflake identifier, so that a name holding a
// lower-case letter or a double quote is the one the caller meant.
func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// sqlLiteral quotes value as a Snowflake string literal.
func sqlLiteral(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

// decodeRSAPrivateKey reads a PKCS#8 private key in PEM, decrypting it with
// passphrase when it is one of the encrypted form.
func decodeRSAPrivateKey(key, passphrase string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("failed to get private key PEM block from key")
	}
	switch block.Type {
	case "PRIVATE KEY":
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse pkcs8 private key: %w", err)
		}
		rsaKey, ok := privateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("expected RSA private key, got %T", privateKey)
		}
		return rsaKey, nil
	case "ENCRYPTED PRIVATE KEY":
		// crypto/x509 cannot decrypt PKCS#8 keys and gosnowflake does not
		// either, so youmark/pkcs8 covers the PBES2 scheme. It reads an empty
		// password as "the key is unencrypted" and would misparse the
		// ciphertext, so reject a missing passphrase before calling it.
		if passphrase == "" {
			return nil, errors.New("private key is encrypted, but no passphrase was provided")
		}
		key, err := pkcs8.ParsePKCS8PrivateKeyRSA(block.Bytes, []byte(passphrase))
		if err != nil {
			return nil, fmt.Errorf("failed to parse pkcs8 private key to rsa private key with passphrase: %w", err)
		}
		return key, nil
	default:
		return nil, fmt.Errorf("unsupported pem block type: %s", block.Type)
	}
}

// Ping reaches the account.
func (d *Driver) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// Close closes the driver's connections.
func (d *Driver) Close() error {
	return d.db.Close()
}

// getVersion reads the server version, which CURRENT_VERSION() states as a
// release such as 9.12.1.
func (d *Driver) getVersion(ctx context.Context) (string, error) {
	const query = "SELECT CURRENT_VERSION()"
	var version string
	if err := d.db.QueryRowContext(ctx, query).Scan(&version); err != nil {
		return "", fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return version, nil
}

// getDatabaseNames reads the databases the login can see, without the inbound
// shared ones: a share is immutable and its DDL is not the account's to read,
// so Bytebase left such a database out of the list rather than syncing half of
// it.
func (d *Driver) getDatabaseNames(ctx context.Context) ([]string, error) {
	inbound, err := d.getInboundShareDatabases(ctx)
	if err != nil {
		return nil, err
	}

	// Ordered by name: a database list a sync reports is the caller's, and the
	// order the account happens to answer in is not.
	const query = `
		SELECT
			DATABASE_NAME
		FROM SNOWFLAKE.INFORMATION_SCHEMA.DATABASES
		ORDER BY DATABASE_NAME`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		if !inbound[name] {
			databases = append(databases, name)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return databases, nil
}

// columnIndex finds the index of a column the server named, or -1 when its
// answer states none: a SHOW statement's columns are read by name, since their
// position is the server's to choose.
func columnIndex(columns []string, name string) int {
	for i, column := range columns {
		if strings.EqualFold(column, name) {
			return i
		}
	}
	return -1
}

// getInboundShareDatabases reads the databases an inbound share created.
func (d *Driver) getInboundShareDatabases(ctx context.Context) (map[string]bool, error) {
	const query = "SHOW SHARES"
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// https://docs.snowflake.com/en/sql-reference/sql/show-shares#output
	kindIndex := columnIndex(columns, "kind")
	databaseIndex := columnIndex(columns, "database_name")
	if kindIndex < 0 || databaseIndex < 0 {
		// A server whose answer states neither column reports no share, and the
		// database list is left whole rather than misread.
		return nil, nil
	}
	values := make([]sql.NullString, len(columns))
	refs := make([]any, len(columns))
	for i := range values {
		refs[i] = &values[i]
	}
	inbound := make(map[string]bool)
	for rows.Next() {
		if err := rows.Scan(refs...); err != nil {
			return nil, err
		}
		if values[kindIndex].String == "INBOUND" {
			inbound[values[databaseIndex].String] = true
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return inbound, nil
}
