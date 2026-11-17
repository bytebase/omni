// Package driver is the contract between Bytebase and the engine driver
// packages (pg/driver, ...). It holds the types that cross that boundary and
// nothing else: no engine imports, no client library, no connection code.
//
// The cut sits between the path to a server and the conversation with it.
// Bytebase resolves who may reach which server over which path: secrets from
// secret managers, cloud IAM tokens, SSH tunnels, TLS material. It passes the
// result as plain values and funcs in a Config. The engine package owns
// everything after the dial: session setup, metadata sync, statement
// execution, result decoding. Connection lifecycle stays with the caller,
// which opens a driver, keeps or pools it, and closes it.
//
// An engine package exposes
//
//	func Open(ctx context.Context, cfg driver.Config) (*Driver, error)
//
// whose *Driver implements Driver and each capability interface the engine
// supports (Syncer, Dumper, Querier, Executor). Open does not reach the
// server; Ping does. A Driver is safe for concurrent use.
//
// Identifiers that a caller maps onto its own enums cross as strings equal
// to those enum names (Message.Level, TransactionOp), so the two sides
// tolerate being on different versions: a name one side does not know is
// ignored or dropped, never an error.
package driver

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"

	"github.com/bytebase/omni/metadata"
)

// Config is everything an engine package needs to reach one server, with
// every credential already resolved.
type Config struct {
	// Host and Port name the server, in the form the engine's own clients
	// accept: pg takes a comma-separated host list and a matching port list.
	// With Dial set they are what Dial receives, and an engine whose dialer
	// ignores them (a cloud connector naming its instance itself) may leave
	// them empty.
	Host string
	Port string
	// Database is the database the driver connects to and Syncer.SyncDatabase
	// reads. Empty means the engine's default (pg: postgres).
	Database string

	// User and Password are the resolved login. A cloud IAM token is a
	// Password.
	User     string
	Password string

	// TLS returns the TLS settings for one host the driver dials. Nil means
	// the engine's default negotiation, which Params may tune (pg: libpq's
	// sslmode, prefer when absent). Non-nil means TLS is required; the driver
	// calls it once per host it may dial, so a caller can bind certificate
	// verification to the host.
	TLS func(host string) (*tls.Config, error)

	// Dial opens the network connection in place of net.Dialer: an SSH
	// tunnel or a cloud connector. The driver still speaks the engine's
	// protocol, TLS included, over what Dial returns.
	Dial func(ctx context.Context, network, address string) (net.Conn, error)

	// Params are connection parameters in the engine's own vocabulary and
	// syntax, applied after the settings above so a caller can override
	// them: pg takes libpq keywords (sslmode, target_session_attrs) and
	// pgx's (default_query_exec_mode, statement_cache_capacity), each value
	// written as a libpq connection string writes it, quoted where it holds
	// a space; a name the engine does not know reaches the server as a
	// session parameter where the protocol has one.
	Params map[string]string

	// ApplicationName identifies the session to the server, where the
	// protocol carries one (pg: application_name).
	ApplicationName string
	// ReadOnly asks for a session that refuses writes. An engine that cannot
	// open one says so in its package doc and ignores the flag.
	ReadOnly bool
	// EngineVersion is the server version the caller last recorded. A driver
	// asks the server and uses this only when the server does not answer.
	EngineVersion string
}

// ErrNotFound reports that an object the call names does not exist on the
// server, such as the database a Syncer was opened on.
var ErrNotFound = errors.New("not found")

// Driver is one open driver. Close releases what Open acquired; Ping
// reaches the server.
type Driver interface {
	Ping(ctx context.Context) error
	Close() error
}

// Syncer reads metadata from the server.
type Syncer interface {
	// SyncInstance returns the server version, its databases without their
	// schemas, and its roles.
	SyncInstance(ctx context.Context) (*Instance, error)
	// SyncDatabase returns the full schema of Config.Database.
	SyncDatabase(ctx context.Context) (*metadata.DatabaseSchemaMetadata, error)
}

// Dumper renders a synced schema as the engine's DDL. An implementation falls
// into one of the three cases below.
//
// A metadata-rendered dump renders the metadata it is given and nothing else: a
// caller that leaves schemas out removes them first. A server-DDL dump states
// the server's own definition of every object in the database the driver is
// open on and ignores the argument, so an edited or narrowed metadata changes
// nothing in what it writes. A no-op dump writes nothing and returns no error,
// and its package doc says why the engine has no DDL to render.
type Dumper interface {
	Dump(ctx context.Context, w io.Writer, db *metadata.DatabaseSchemaMetadata) error
}
