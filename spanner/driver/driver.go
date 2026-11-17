// Package driver is Spanner's implementation of the driver contract in
// github.com/bytebase/omni/driver: it reaches a Spanner instance through the
// cloud.google.com/go/spanner clients, syncs the instance's databases and one
// database's INFORMATION_SCHEMA into omni/metadata, and runs a script as
// queries and change scripts.
//
// It implements Driver, Syncer, Dumper, Querier and Executor. The statements are
// the database's dialect: GoogleSQL, which omni's googlesql tree parses and
// analyses, or PostgreSQL, which omni's pg tree splits and lowers a LIMIT in.
// This package is the conversation with a Spanner service, so it lives in a tree
// of its own beside googlesql, as dynamodb does beside partiql; a service tree
// consults the language tree of each dialect it serves.
//
// Dialects. A database's dialect is fixed when it is created, and the two answer
// different SQL and different INFORMATION_SCHEMA spellings. SyncInstance lists
// both kinds; SyncDatabase reads the database's dialect from the admin API and
// asks the matching catalog — PostgreSQL names its system schemas in lower case
// and states SPANNER_IS_MANAGED as text, where GoogleSQL upper-cases the schemas
// and states a boolean. A script is split, and a query's LIMIT lowered, by the
// dialect's own grammar: googlesql/parser for GoogleSQL and pg for PostgreSQL,
// whose dollar-quoted strings hold semicolons and whose `#>` operator GoogleSQL
// reads as a comment. A GoogleSQL database takes the queries and the split it
// took before PostgreSQL was supported.
//
// Reading the dialect is the one way this driver departs from the contract's
// "Open does not reach the server": Open asks the admin API for the database,
// because every statement is split by the dialect's grammar and a query must not
// wait for a schema sync to learn which. A database that is not there leaves the
// dialect GoogleSQL instead of failing Open, so a caller can still open a
// database to drop it; SyncDatabase reports it as driver.ErrNotFound. A
// creation script that names its database in double quotes creates a PostgreSQL
// database — the admin API requires the script to be written in the dialect it
// creates — and anything else, including a backtick-quoted or unquoted name,
// creates a GoogleSQL one.
//
// Config.Params carry the two parts of the resource name the contract has no
// field for:
//
//	project   the Google Cloud project the instance lives in.
//	instance  the Spanner instance.
//
// Config.Database is the database name, and together the three name the database
// resource projects/<project>/instances/<instance>/databases/<database>.
//
// Config.Password is the resolved credential: a Google Cloud credential JSON
// document, of the kind a service account, an external account, an
// impersonated service account or an authorized user is written as. Its "type"
// member decides which, as Bytebase's GCPCredentialOption did. An empty
// Config.Password means the caller resolved no credential of its own and the
// Google authentication library's default chain reads one (the environment, the
// well-known file, an attached service account). Credential resolution of every
// other kind stays with the caller: a service account a caller impersonates is
// the caller's to impersonate, and what arrives here is the credential it
// resolved. Config.User does nothing: Spanner has no login beside the
// credential.
//
// Config.Host and Config.Port optionally name the endpoint, which is how a
// Private Service Connect endpoint or a proxy is reached; an empty Config.Host
// leaves the client's own spanner.googleapis.com:443. gRPC defaults to port 443
// when Config.Port is empty. Config.TLS makes TLS required with the settings it
// returns for the endpoint's host, which is how a caller pins a certificate;
// nil leaves the client's transport, which verifies the server's certificate
// against the system roots. Config.Dial opens the connections the client makes,
// so an SSH tunnel or a cloud connector reaches the endpoint unchanged.
//
// Config.Params[driver.ParamEmulator] ("emulator") names the Cloud Spanner
// emulator: it accepts any caller over plaintext gRPC, so the driver sends no
// credential and no TLS, and asks its client not to export metrics, which an
// emulator has no service for. It is what a recording sets; a real instance
// needs none of it.
//
// Config.ApplicationName, Config.ReadOnly and Config.EngineVersion do nothing:
// Spanner's protocol carries no application name, a database has no read-only
// session to open, and Spanner reports no version.
//
// A session is the driver: Spanner has no session state a statement can set, so
// a temporary table on one statement is not there for the next. Each statement
// is read in a single-use read-only transaction of its own.
// QueryOptions.Schema and QueryOptions.MaskingEnabled do nothing: a name
// resolves in the database the driver is open on, and a row's cells are the
// typed columns Spanner declares, which the caller traces through the statement
// it also holds. QueryOptions.Limit rewrites a query's outermost LIMIT with
// googlesql/analysis.RewriteResultLimit, and also caps the rows read when the
// rewrite cannot (a non-literal LIMIT, a shape the parser does not read).
// QueryOptions.MaxResultBytes stops the reading at the first row that takes the
// rows kept past the cap and sets Result.Truncated: the kept rows are the most
// that fit and the one that crossed, so they exceed the cap by at most that row.
//
// QueryOptions.Explain asks for the plan Spanner would run. The plan is a JSON
// document in the result's rows, as Bytebase rendered it, and Result.Plan says
// its format. Spanner refuses DML in a single-use read-only transaction, so a
// statement QueryOptions.Explain names that is not a query is analysed in a
// read-write transaction that is rolled back, and nothing it writes commits. An
// emulator has no optimizer and answers no plan; see the emulator note below.
//
// A row's cells are the columns' Spanner types: INT64 as an INT64, FLOAT64 as a
// DOUBLE, BOOL as a BOOL, BYTES as BYTES (decoded from the base64 wire form),
// TIMESTAMP as a TIMESTAMP in UTC with the fractional digits the server wrote,
// and DATE, STRING and JSON as their text. A Spanner ARRAY or STRUCT has no cell
// kind of its own, so it crosses as its JSON document; a STRUCT column, or an
// ARRAY of STRUCT, is refused before any row is read, as Bytebase refused it.
// A NUMERIC and a PG_NUMERIC arrive as the text Spanner writes. A PostgreSQL
// database states the same type codes for its own types — int8 as INT64, text
// and varchar as STRING, numeric as NUMERIC annotated PG_NUMERIC, jsonb as JSON
// annotated PG_JSONB, bytea as BYTES, timestamptz as TIMESTAMP, and an array as
// ARRAY — so a result decodes as it does for GoogleSQL.
//
// Execute routes each statement by what it is: DDL — CREATE, ALTER, DROP, GRANT
// and REVOKE — goes to the database admin API's UpdateDatabaseDdl, in order and
// outside any transaction, because Spanner DDL is not transactional with DML; a
// query runs as QueryOptions does; and the rest is DML, which runs through the
// data client. ExecuteOptions.Transaction
// ON (and an empty value, which is Bytebase's default) runs a script's DML as one
// read-write transaction, and OFF runs each statement as a partitioned update,
// whose row count is Spanner's lower bound. DDL in a script takes the DDL path
// whatever the mode, which is Bytebase's cut. ExecuteOptions.Isolation does
// nothing: Spanner offers one isolation level, serializable, and no way to ask
// for another. ExecuteOptions.CreateDatabase creates the database from the
// script's first statement and applies the rest as its extra DDL statements.
// ExecuteOptions.MaxRetries retries nothing: Spanner names the lock conflicts it
// reports (ABORTED) for the caller to retry, and Bytebase's driver did not retry
// them either. ExecuteOptions.Log receives transaction control only, BEGIN,
// COMMIT and ROLLBACK: Bytebase's Spanner driver logged no per-statement command
// line, and this one does not either.
//
// Dump is the database's DDL as the server states it (GetDatabaseDdl), not a
// rendering of the synced metadata: googlesql has no catalog that turns a synced
// schema back into DDL, and the server's own DDL is what Bytebase dumped. The
// metadata argument is therefore read by nothing, and a caller that left a table
// or a view out of it still dumps every object the server holds, exactly as
// Bytebase's driver did.
//
// A statement's error is the client's own, as it returned it. A failure the
// server reported is a *spanner.Error, which holds the *apierror.APIError the
// RPC answered; the client deprecates *spanner.Error and tells a caller to read
// the API error instead, so what is stable is the gRPC status both state:
// status.FromError reports it and spanner.ErrCode reads its Code, which is
// codes.InvalidArgument for a statement the server refused — a table that is not
// there included — and codes.NotFound for a database the admin API does not have.
//
// The emulator note. A recording's server is the Cloud Spanner emulator, pinned
// in spanner/testdata/pins.tsv and started by spanner/internal/testcontainer. The
// emulator answers what the driver asks of a server — the admin API, DDL, DML,
// queries, INFORMATION_SCHEMA — so sync, query, execute, dump and ping are
// recorded from it. What it cannot answer is documented here and never faked:
// IAM and the credential chain (the emulator accepts any caller, so a credential
// path is exercised only by what the driver builds from Config, not by a
// server's answer), a Private Service Connect endpoint and multi-region routing
// (the emulator is one process on one address), a real instance's latency (a
// tape holds latency as the driver's own measurement of a replay), and query
// plans (an emulator has no optimizer, so AnalyzeQuery answers no plan; the plan
// rendering is tested from a recorded plan's shape, and the live plan is called
// out as unverified).
package driver

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync/atomic"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bytebase/omni/driver"
)

var dsnRegExp = regexp.MustCompile("projects/(?P<PROJECTGROUP>([a-z]|[-.:]|[0-9])+)/instances/(?P<INSTANCEGROUP>([a-z]|[-]|[0-9])+)/databases/(?P<DATABASEGROUP>([a-z]|[-]|[_]|[0-9])+)")

var (
	_ driver.Driver   = (*Driver)(nil)
	_ driver.Syncer   = (*Driver)(nil)
	_ driver.Dumper   = (*Driver)(nil)
	_ driver.Querier  = (*Driver)(nil)
	_ driver.Executor = (*Driver)(nil)
)

// The names this driver reads out of driver.Config.Params.
const (
	// ParamProject is the Google Cloud project the instance lives in.
	ParamProject = "project"
	// ParamInstance is the Spanner instance.
	ParamInstance = "instance"
)

// Driver is an open Spanner driver: one database of one instance, with a client
// for its data and one for its schema.
type Driver struct {
	client client

	// databaseName is the database the driver is open on.
	databaseName string
	// instancePath is projects/<project>/instances/<instance>.
	instancePath string
	// dialect is the database's dialect. Open reads it, and a creation or a sync
	// of the database refreshes it, so it is read atomically: a Driver is safe
	// for concurrent use.
	dialect atomic.Int32
}

// getDialect returns the database's dialect.
func (d *Driver) getDialect() dialect {
	return dialect(d.dialect.Load())
}

// setDialect records the database's dialect.
func (d *Driver) setDialect(value dialect) {
	d.dialect.Store(int32(value))
}

// Open returns a driver for the instance and database cfg names. It reads the
// database's dialect, which decides the catalog a sync reads and the grammar a
// script is split by, and does not otherwise reach the server; Ping does.
func Open(ctx context.Context, cfg driver.Config) (*Driver, error) {
	if cfg.Params[ParamProject] == "" {
		return nil, errors.New("project ID cannot be empty")
	}
	if cfg.Params[ParamInstance] == "" {
		return nil, errors.New("instance ID cannot be empty")
	}
	c, err := newClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	d := &Driver{
		client:       c,
		databaseName: cfg.Database,
		instancePath: InstancePath(cfg.Params[ParamProject], cfg.Params[ParamInstance]),
	}
	if cfg.Database != "" {
		database, err := c.GetDatabase(ctx, d.databasePath(""))
		switch {
		case err == nil:
			d.setDialect(dialectOf(database.Dialect))
		case status.Code(err) == codes.NotFound:
			// A database that is not there is not Open's error: a caller may open
			// one to drop it, and SyncDatabase reports it as driver.ErrNotFound.
			// Its dialect stays GoogleSQL until a creation or a sync of that
			// database reads it.
		default:
			_ = c.Close()
			return nil, fmt.Errorf("failed to read the dialect of database %q: %w", cfg.Database, err)
		}
	}
	return d, nil
}

// Ping reads one database name as a connectivity check, which is what
// Bytebase's driver did: Spanner has no ping of its own.
func (d *Driver) Ping(ctx context.Context) error {
	if err := d.client.Ping(ctx); err != nil {
		return fmt.Errorf("spanner: bad connection: %w", err)
	}
	return nil
}

// Close closes the driver and the connections its clients opened.
func (d *Driver) Close() error {
	return d.client.Close()
}

// databasePath is the resource name of database: the database the driver is open
// on when database is empty.
func (d *Driver) databasePath(database string) string {
	if database == "" {
		database = d.databaseName
	}
	return DatabasePath(d.instancePath, database)
}

// InstancePath returns the instance resource name,
// projects/<project>/instances/<instance>.
func InstancePath(project, instance string) string {
	return fmt.Sprintf("projects/%s/instances/%s", project, instance)
}

// DatabasePath returns the database resource name under instancePath,
// projects/<project>/instances/<instance>/databases/<database>.
func DatabasePath(instancePath, database string) string {
	return fmt.Sprintf("%s/databases/%s", instancePath, database)
}

// getDatabaseFromDSN returns <database> from
// projects/<project>/instances/<instance>/databases/<database>.
func getDatabaseFromDSN(dsn string) (string, error) {
	match := dsnRegExp.FindStringSubmatch(dsn)
	if match == nil {
		return "", errors.New("invalid DSN")
	}
	matches := make(map[string]string)
	for i, name := range dsnRegExp.SubexpNames() {
		if i != 0 && name != "" {
			matches[name] = match[i]
		}
	}
	return matches["DATABASEGROUP"], nil
}
