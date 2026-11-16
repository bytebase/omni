// Package driver is MongoDB's implementation of the driver contract in
// github.com/bytebase/omni/driver: it connects through the mongo-go-driver
// client, runs mongosh statements through gomongo, syncs a database's
// collections and indexes into omni/metadata, dumps them as mongosh, and runs
// interactive queries and change scripts.
//
// It implements Driver, Syncer, Dumper, Querier and Executor. Bytebase's
// driver had no live probe of its own, so there is none here.
//
// Config.Host is the endpoint's authority as a MongoDB connection URI writes
// it: one host, or a comma-separated list of hosts that each carry their port
// (node1:27017,node2:27017). Config.Port is appended to a single host that has
// none. Config.Params are connection string options in MongoDB's own syntax
// (replicaSet, readPreference, tlsInsecure, ...); they apply last, after the
// defaults, so a caller can override them. authSource defaults to admin and
// appName to Config.ApplicationName, then bytebase.
//
// An SRV endpoint, mongodb+srv://, is an engine-specific setting: Options.SRV
// selects it, since a DNS seed list's domain is not a host list.
//
// Config.TLS supplies the TLS settings for one host and makes TLS required: the
// driver dials each server itself and asks for that server's settings, so a
// replica set's members and an SRV seed list's hosts are each verified against
// their own name. A tls or ssl value in Params does nothing then, since the
// driver has taken TLS over; nil leaves the client's own TLS negotiation, which
// Params and the mongodb+srv scheme do configure. Config.Dial, when set, opens
// the connection the driver hands to TLS, so an SSH tunnel or a cloud connector
// reaches the server the same way.
//
// Config.ReadOnly, QueryOptions.Schema, QueryOptions.MaskingEnabled and
// ExecuteOptions' transaction settings do nothing: MongoDB has no read-only
// session, no schema per session, no masking to withhold, and no
// multi-statement transaction the API this driver serves takes part in.
// Execute runs each statement on its own.
//
// QueryOptions.MaxResultBytes stops a result once its rows pass the cap and
// sets Result.Truncated; QueryOptions.Limit caps the rows a read asks the
// client for. A count keeps counting every document: a cap on the rows an
// editor shows is not a cap on what a count counts.
package driver

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/bytebase/omni/driver"
)

var (
	_ driver.Driver   = (*Driver)(nil)
	_ driver.Syncer   = (*Driver)(nil)
	_ driver.Dumper   = (*Driver)(nil)
	_ driver.Querier  = (*Driver)(nil)
	_ driver.Executor = (*Driver)(nil)
)

// Options are MongoDB-specific settings for Open.
type Options struct {
	// SRV dials a mongodb+srv:// endpoint instead of mongodb://: Config.Host
	// is then the DNS seed list's domain and carries no port.
	SRV bool
}

// Driver is an open MongoDB driver: one client to one database.
type Driver struct {
	client       client
	databaseName string
}

// connectionURI returns the connection string the client is opened with.
func connectionURI(cfg driver.Config, opts Options) (string, error) {
	u := &url.URL{
		// In RFC, there can be no tailing slash('/') in the path if the path is empty and the query is not empty.
		// The Go driver throws "error parsing uri: must have a / before the query ?" without it.
		Scheme: "mongodb",
		Path:   "/",
	}
	if opts.SRV {
		u.Scheme = "mongodb+srv"
	}
	authority, err := endpoint(cfg.Host, cfg.Port)
	if err != nil {
		return "", err
	}
	u.Host = authority
	if cfg.User != "" {
		u.User = url.UserPassword(cfg.User, cfg.Password)
	}
	if cfg.Database != "" {
		u.Path = cfg.Database
	}

	values := u.Query()
	values.Set("authSource", "admin")
	applicationName := cfg.ApplicationName
	if applicationName == "" {
		applicationName = "bytebase"
	}
	values.Set("appName", applicationName)
	for key, value := range cfg.Params {
		if key == "" {
			continue
		}
		values.Set(key, value)
	}
	if cfg.TLS != nil {
		// The driver dials TLS itself, one host at a time, so a tls or ssl
		// value in Params goes: either would wrap the connection twice, and
		// the two aliases disagreeing is a parse error in the client. One
		// tls=false says the client leaves the TLS to the driver.
		for key := range values {
			if strings.EqualFold(key, "tls") || strings.EqualFold(key, "ssl") {
				values.Del(key)
			}
		}
		values.Set("tls", "false")
	}
	u.RawQuery = values.Encode()
	return u.String(), nil
}

// endpoint returns the connection string's authority: the host list as
// written, or one host with the port appended unless it carries one already.
func endpoint(host, port string) (string, error) {
	host, port = strings.TrimSpace(host), strings.TrimSpace(port)
	if host == "" {
		return "", errors.New("host must be set")
	}
	// A list names each host's own port, which may differ from the others'.
	if strings.Contains(host, ",") {
		return host, nil
	}
	if _, _, err := net.SplitHostPort(host); err == nil {
		return host, nil
	}
	// A bracketed literal is an address without its port, so the brackets come
	// off before JoinHostPort puts them back.
	host = strings.TrimSuffix(strings.TrimPrefix(host, "["), "]")
	if port == "" {
		if net.ParseIP(host) != nil {
			return "[" + host + "]", nil
		}
		return host, nil
	}
	return net.JoinHostPort(host, port), nil
}

// Open returns a driver for the server and database cfg names. It does not
// reach the server; Ping does.
func Open(_ context.Context, cfg driver.Config, opts ...Options) (*Driver, error) {
	var opt Options
	if len(opts) > 0 {
		opt = opts[0]
	}
	connectionString, err := connectionURI(cfg, opt)
	if err != nil {
		return nil, err
	}
	clientOptions := options.Client().ApplyURI(connectionString)
	switch {
	case cfg.TLS != nil:
		clientOptions.SetDialer(tlsDialer{config: cfg})
	case cfg.Dial != nil:
		clientOptions.SetDialer(dialFunc(cfg.Dial))
	}
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}
	return &Driver{
		client:       liveClient{client: client},
		databaseName: cfg.Database,
	}, nil
}

// dialFunc adapts Config.Dial to the client's dialer.
type dialFunc func(ctx context.Context, network, address string) (net.Conn, error)

func (f dialFunc) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return f(ctx, network, address)
}

// tlsDialer opens the connection to one server and hands it to TLS with the
// settings Config.TLS returns for that server's host. One *tls.Config for a
// whole cluster would verify every member against the seed it was built from.
type tlsDialer struct {
	config driver.Config
}

func (d tlsDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		host = strings.Trim(address, "[]")
	}
	// The settings come first, so a caller that refuses the host costs no
	// connection.
	tlsConfig, err := d.config.TLS(host)
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS config: %w", err)
	}
	if tlsConfig == nil {
		return nil, errors.New("TLS configuration must not be nil")
	}
	conn, err := dial(ctx, d.config.Dial, network, address)
	if err != nil {
		return nil, err
	}
	tlsConfig = tlsConfig.Clone()
	if tlsConfig.ServerName == "" {
		tlsConfig.ServerName = host
	}
	tlsConn := tls.Client(conn, tlsConfig)
	if err := tlsConn.HandshakeContext(ctx); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return tlsConn, nil
}

// dial opens the network connection, through Config.Dial when the caller set
// one and the default dialer otherwise.
func dial(ctx context.Context, via func(ctx context.Context, network, address string) (net.Conn, error), network, address string) (net.Conn, error) {
	if via != nil {
		return via(ctx, network, address)
	}
	var d net.Dialer
	return d.DialContext(ctx, network, address)
}

// MinSupportedServerVersion is the oldest MongoDB release the driver in go.mod
// accepts. The driver enforces it by wire version while selecting a server and
// refuses anything older, so this constant tracks go.mod's driver and nothing
// else: raising it is a breaking change for existing instances.
const MinSupportedServerVersion = "4.4"

// Ping pings the server.
func (d *Driver) Ping(ctx context.Context) error {
	if err := d.client.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", annotateUnsupportedServerVersion(err))
	}
	return nil
}

// Close closes the client.
func (d *Driver) Close() error {
	if err := d.client.Disconnect(context.Background()); err != nil {
		return fmt.Errorf("failed to disconnect MongoDB: %w", err)
	}
	return nil
}

// annotateUnsupportedServerVersion names the oldest supported release when the
// driver refused the server for being too old, because the driver's own message
// reports only a wire version number. The driver returns that refusal as a plain
// error from server selection rather than a distinguishable type, so matching its
// text is the only way to recognize it; any other error is returned unchanged.
func annotateUnsupportedServerVersion(err error) error {
	if err == nil || !strings.Contains(err.Error(), "wire version") {
		return err
	}
	return fmt.Errorf("MongoDB %s or later is required: %w", MinSupportedServerVersion, err)
}
