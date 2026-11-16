// Package driver is Redis's implementation of the driver contract in
// github.com/bytebase/omni/driver: it connects through the go-redis client,
// syncs an instance's version and its database numbers, and runs Redis
// commands as interactive queries and change scripts.
//
// It implements Driver, Syncer, Querier and Executor. There is no Dumper:
// Redis keeps no schema to render and Bytebase's Dump wrote nothing, so a synced
// database is its name alone. The contract's probe, PrivateDatabaseLinkOwner,
// does not apply either: Redis has no private database link. The Executor's DDL,
// transaction and retry settings do not apply: Redis creates no database and
// runs no transaction, so Execute runs a script's commands one after another and
// nothing rolls back. Syncer.SyncInstance reports no roles, as Bytebase's did
// not read any, and a database's keys are not a schema, so SyncDatabase is the
// database number alone.
//
// Redis has no SQL surface, so nothing under partiql/ or analysis is involved.
// What the driver decides is what a script is — one command per line, its
// arguments parsed as a shell parses them — and what a server's answer means:
// the database numbers from CONFIG GET databases, or from the keyspace section
// when the account may not run CONFIG.
//
// Config.Host is one server, or a comma-separated list of them for Sentinel and
// Cluster, each optionally carrying its own port. Config.Port applies to a
// single host that has none, and both default to 6379. Config.Database is the
// numeric database index the session selects, as Redis writes it ("0" through
// "15" on a default server); empty means 0. Config.User and Config.Password are
// the login, which for Sentinel is the master's; Options.SentinelUser and
// Options.SentinelPassword are the sentinels' own when they differ.
//
// Config.TLS answers for the seed host, Config.Host: go-redis carries one
// *tls.Config per client, so a cluster's other nodes and a sentinel's replicas
// are verified against the settings answered for the seed, as Bytebase's one
// config did. Config.Dial opens the network connection in place of net.Dialer,
// so an SSH tunnel or a cloud connector reaches the server, and TLS is
// negotiated over what it returns.
//
// Options.Type selects the topology, because it decides which client the driver
// builds — go-redis has one for a server, one for a sentinel deployment and one
// for a cluster — rather than a setting that client applies. Config.Params, Config.ApplicationName,
// Config.ReadOnly and Config.EngineVersion do nothing: go-redis takes no
// connection string this driver fills, and Redis has no application name and no
// read-only session. Config.EngineVersion goes too: Bytebase's sync always asked
// the server and failed when it did not answer, so no stored version stands in.
//
// A query runs each line of a script as one command and returns one Result per
// command, with the columns "#" and "Value". A reply that is an array becomes
// one row per element and any other reply one row. A missing key is a row with
// a NULL value, as Bytebase showed it; a command the server refused ends the
// script on that command's Result, which errs with the server's own error, so
// redis.Nil is the one a caller compares by identity. QueryOptions.MaxResultBytes
// stops a command's rows at the cap and sets Result.Truncated; Bytebase measured
// proto.Size of the whole result at power-of-two row counts and answered an error
// message instead, so the cut here lands at a different row and a caller reports
// the cap from the flag. QueryOptions.Schema, Limit, ExplainFormat and
// MaskingEnabled do nothing, and Explain is refused: Redis states no plan and
// Redis commands are not statements a row limit or a column trace rewrites.
// Session's own QueryAllNodes runs a script on every cluster node, which is
// Bytebase's per-request choice the contract has no field for; it answers one
// server's driver with Bytebase's error text.
//
// Execute sends every command of a script as one pipeline and returns the rows
// affected, which is 0: a Redis command reports no row count. A command's error
// ends the script, and ExecuteOptions' transaction settings, CreateDatabase and
// MaxRetries do nothing. ExecuteOptions.Log receives nothing, as it did in
// Bytebase, which logged none of a script's commands.
package driver

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/bytebase/omni/driver"
)

var (
	_ driver.Driver   = (*Driver)(nil)
	_ driver.Syncer   = (*Driver)(nil)
	_ driver.Querier  = (*Driver)(nil)
	_ driver.Executor = (*Driver)(nil)
)

// Topology is how a driver reaches its server.
type Topology string

// The topologies a Redis endpoint can have.
const (
	// Standalone is one server. It is the default, and what an unspecified
	// topology means.
	Standalone Topology = "STANDALONE"
	// Sentinel is a master that a set of sentinels monitors.
	Sentinel Topology = "SENTINEL"
	// Cluster is a Redis Cluster, whose nodes place keys by hash slot.
	Cluster Topology = "CLUSTER"
)

// Options are Redis's own settings for Open.
type Options struct {
	// Type is the server's topology: Standalone when empty.
	Type Topology
	// MasterName is the master the sentinels monitor. Sentinel only.
	MasterName string
	// SentinelUser and SentinelPassword are the sentinels' login, where it
	// differs from Config.User and Config.Password. Sentinel only.
	SentinelUser     string
	SentinelPassword string
}

// Driver is an open Redis driver.
type Driver struct {
	// client is what the driver asks of a Redis client. A test's tape replaces
	// it.
	client client
	// databaseName is the numeric database index as Redis writes it, and the
	// name a synced database carries.
	databaseName string
}

// Open returns a driver for the server cfg names. It does not reach the server;
// Ping does.
func Open(_ context.Context, cfg driver.Config, opts ...Options) (*Driver, error) {
	var opt Options
	if len(opts) > 0 {
		opt = opts[0]
	}
	database, err := databaseIndex(cfg.Database)
	if err != nil {
		return nil, err
	}
	addresses, err := endpoints(cfg.Host, cfg.Port)
	if err != nil {
		return nil, err
	}
	tlsConfig, err := tlsConfig(cfg)
	if err != nil {
		return nil, err
	}
	d := &Driver{databaseName: strconv.Itoa(database)}
	switch opt.Type {
	case "", Standalone:
		if len(addresses) > 1 {
			return nil, fmt.Errorf("a standalone Redis takes one host, but %d were given", len(addresses))
		}
		options := &redis.Options{
			Addr:      addresses[0],
			Username:  cfg.User,
			Password:  cfg.Password,
			TLSConfig: tlsConfig,
			DB:        database,
			Dialer:    cfg.Dial,
			// Pipelines run on the main pool. A driver is opened per request
			// and runs a single pipeline, so go-redis's dedicated pipeline pool
			// would only add a connection and its handshake to every request.
			PipelinePoolSize: -1,
		}
		d.client = &liveClient{client: redis.NewClient(options)}
	case Sentinel:
		options := &redis.FailoverOptions{
			MasterName:       opt.MasterName,
			Username:         cfg.User,
			Password:         cfg.Password,
			SentinelUsername: opt.SentinelUser,
			SentinelPassword: opt.SentinelPassword,
			SentinelAddrs:    addresses,
			DB:               database,
			TLSConfig:        tlsConfig,
			Dialer:           cfg.Dial,
			PipelinePoolSize: -1,
		}
		d.client = &liveClient{client: redis.NewFailoverClient(options)}
	case Cluster:
		options := &redis.ClusterOptions{
			Addrs:            addresses,
			Username:         cfg.User,
			Password:         cfg.Password,
			TLSConfig:        tlsConfig,
			Dialer:           cfg.Dial,
			PipelinePoolSize: -1,
			// go-redis 9.22.0 raised its own default to a minute. A node that
			// refuses connections does not trigger a reload, so a failover is
			// only noticed by this periodic one, and a driver that outlives it
			// keeps routing to the master it has lost.
			// https://github.com/redis/go-redis/issues/3979
			ClusterStateReloadInterval: 10 * time.Second,
		}
		d.client = &liveClient{client: redis.NewClusterClient(options)}
	default:
		return nil, fmt.Errorf("unsupported redis type %s", opt.Type)
	}
	return d, nil
}

// databaseIndex reads the numeric database index a driver's session selects.
// Empty is database 0, the one every Redis has.
func databaseIndex(database string) (int, error) {
	if database == "" {
		return 0, nil
	}
	index, err := strconv.Atoi(database)
	if err != nil {
		return 0, fmt.Errorf("failed to convert database %s to int: %w", database, err)
	}
	return index, nil
}

// endpoints returns the server addresses, each with its port. A host may carry
// its own port, which a list of hosts uses for each of them; port applies to a
// single host that has none.
func endpoints(host, port string) ([]string, error) {
	if strings.TrimSpace(host) == "" {
		return nil, errors.New("host must be set")
	}
	var addresses []string
	for _, one := range strings.Split(host, ",") {
		one = strings.TrimSpace(one)
		if one == "" {
			return nil, fmt.Errorf("empty host in %q", host)
		}
		if _, _, err := net.SplitHostPort(one); err == nil {
			addresses = append(addresses, one)
			continue
		}
		one = strings.TrimSuffix(strings.TrimPrefix(one, "["), "]")
		p := strings.TrimSpace(port)
		if p == "" {
			p = "6379"
		}
		addresses = append(addresses, net.JoinHostPort(one, p))
	}
	return addresses, nil
}

// tlsConfig returns the TLS settings for the seed host. go-redis takes one
// *tls.Config per client, so this is the config every node of a cluster or a
// sentinel deployment is reached with.
func tlsConfig(cfg driver.Config) (*tls.Config, error) {
	if cfg.TLS == nil {
		return nil, nil
	}
	config, err := cfg.TLS(cfg.Host)
	if err != nil {
		return nil, fmt.Errorf("redis: failed to get tls config: %w", err)
	}
	return config, nil
}

// Ping pings the server.
func (d *Driver) Ping(ctx context.Context) error {
	return d.client.Ping(ctx)
}

// Close closes the client.
func (d *Driver) Close() error {
	return d.client.Close()
}
