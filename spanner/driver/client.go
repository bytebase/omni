package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"cloud.google.com/go/spanner"
	spannerdb "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/bytebase/omni/driver"
)

// client is what the driver asks of Spanner: the admin API calls sync, dump and
// execute make, and the data calls query and execute make. It is the seam a test
// replaces with a tape, so the driver's own code is what a recording and a
// replay exercise; the live implementation below is the one call per method it
// takes to reach the clients.
type client interface {
	// Ping reads one database name as a connectivity check.
	Ping(ctx context.Context) error
	Close() error
	// ListDatabases returns every database of the instance with its dialect.
	ListDatabases(ctx context.Context) ([]Database, error)
	// GetDatabase returns database with its dialect, and driver.ErrNotFound when
	// the admin API says the database is not there.
	GetDatabase(ctx context.Context, database string) (Database, error)
	// UpdateDatabaseDdl applies statements to database and waits for the
	// operation.
	UpdateDatabaseDdl(ctx context.Context, database string, statements []string) error
	// CreateDatabase creates a database in instancePath from createStatement,
	// applies extra as its extra DDL statements, and waits for the operation.
	// The database is GoogleSQL, the admin API's default dialect.
	CreateDatabase(ctx context.Context, instancePath, createStatement string, extra []string) error
	// GetDatabaseDdl returns database's DDL statements, in the order the server
	// holds them.
	GetDatabaseDdl(ctx context.Context, database string) ([]string, error)
	// Read runs statement as one single-use read-only transaction.
	Read(ctx context.Context, statement string) (RowIterator, error)
	// Snapshot returns a reader of one consistent snapshot of the database, held
	// until Close.
	Snapshot(ctx context.Context) (Reader, error)
	// Update runs statement as a DML statement in a read-write transaction and
	// returns the rows it affected.
	Update(ctx context.Context, statement string) (int64, error)
	// BatchUpdate runs statements in one read-write transaction and returns the
	// rows each affected.
	BatchUpdate(ctx context.Context, statements []string) ([]int64, error)
	// PartitionedUpdate runs statement as a partitioned DML statement and
	// returns the rows it affected, which Spanner counts as a lower bound.
	PartitionedUpdate(ctx context.Context, statement string) (int64, error)
	// AnalyzeQuery returns the plan statement would run under without running
	// it. readWrite plans a statement that a single-use read-only transaction
	// cannot read, such as DML, in a read-write transaction that is rolled back.
	AnalyzeQuery(ctx context.Context, statement string, readWrite bool) (*sppb.QueryPlan, error)
}

// ParamEmulator says the endpoint is the Cloud Spanner emulator, which accepts
// any caller over plaintext gRPC: 1, true or yes. It is what a recording sets to
// reach the pinned server, and it sends no credential.
const ParamEmulator = "emulator"

// isEmulator reports whether the caller said the endpoint is the emulator.
func isEmulator(cfg driver.Config) bool {
	switch strings.ToLower(cfg.Params[ParamEmulator]) {
	case "1", "true", "yes":
		return true
	}
	return false
}

// Database is one database ListDatabases or GetDatabase reported.
type Database struct {
	// Name is the resource name,
	// projects/<project>/instances/<instance>/databases/<database>.
	Name string
	// Dialect is the database's SQL dialect as the admin API spells it:
	// GOOGLE_STANDARD_SQL or POSTGRESQL.
	Dialect string
}

// The dialects the admin API spells, which Database.Dialect holds.
const (
	// DialectGoogleSQL is Spanner's own dialect, and the admin API's default.
	DialectGoogleSQL = "GOOGLE_STANDARD_SQL"
	// DialectPostgresSQL is the database's PostgreSQL dialect, which this
	// driver does not serve.
	DialectPostgresSQL = "POSTGRESQL"
)

// errPostgres is the refusal of a PostgreSQL-dialect database: omni parses
// Spanner's GoogleSQL dialect only.
func errPostgres(database string) error {
	return fmt.Errorf("spanner database %q uses the PostgreSQL dialect, which is not supported; only GoogleSQL databases are", database)
}

// Reader reads one snapshot of a database.
type Reader interface {
	Query(ctx context.Context, statement string) (RowIterator, error)
	Close()
}

// RowIterator is what a statement answered: the columns it declares and its
// rows, read on demand.
type RowIterator interface {
	// Columns returns the columns the result declares. A New method rather than
	// a field: Spanner states the row type with the first row, so a caller asks
	// after a row arrived, or after the end.
	Columns() []Column
	// Next returns the next row, or io.EOF at the end. An error means the read
	// failed where it stood, and the rows already returned stay read.
	Next() ([]Cell, error)
	Stop()
}

// Column is one column a result declares.
type Column struct {
	Name string
	Type *sppb.Type
}

// Cell is one value as Spanner sent it: the type the column declares and the
// JSON-shaped value the client decoded.
type Cell struct {
	Type  *sppb.Type
	Value *structpb.Value
}

// liveClient is the client the driver opens against Spanner: the admin client
// for the schema calls and, when the caller named a database, the data client.
type liveClient struct {
	admin        *spannerdb.DatabaseAdminClient
	data         *spanner.Client
	instancePath string
}

// newClient builds the clients cfg describes. It does not reach the server; Ping
// does. It is the one place the contract's values become Google Cloud client
// options.
func newClient(ctx context.Context, cfg driver.Config) (client, error) {
	options, err := clientOptions(cfg)
	if err != nil {
		return nil, err
	}
	admin, err := spannerdb.NewDatabaseAdminClient(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create the Spanner database admin client: %w", err)
	}
	live := &liveClient{
		admin:        admin,
		instancePath: InstancePath(cfg.Params[ParamProject], cfg.Params[ParamInstance]),
	}
	if cfg.Database != "" {
		// An instance's client exports native metrics to Cloud Monitoring unless
		// it is told not to, and an emulator has no such service: a caller that
		// names the emulator gets a client that exports none, which is what the
		// client does for itself when it finds an emulator in its environment.
		clientConfig := spanner.ClientConfig{}
		if isEmulator(cfg) {
			clientConfig.DisableNativeMetrics = true
		}
		data, err := spanner.NewClientWithConfig(ctx, DatabasePath(live.instancePath, cfg.Database), clientConfig, options...)
		if err != nil {
			_ = admin.Close()
			return nil, fmt.Errorf("failed to create the Spanner client: %w", err)
		}
		live.data = data
	}
	return live, nil
}

// clientOptions is the client configuration cfg asks for: the credential, the
// endpoint, TLS, and the caller's dialer, in that order.
func clientOptions(cfg driver.Config) ([]option.ClientOption, error) {
	var options []option.ClientOption
	credential, err := credentialOption(cfg.Password)
	if err != nil {
		return nil, err
	}
	if credential != nil {
		options = append(options, credential)
	}
	// An empty host leaves the client's own spanner.googleapis.com:443, which is
	// how real Spanner is reached.
	if endpoint := endpoint(cfg); endpoint != "" {
		options = append(options, option.WithEndpoint(endpoint))
	}
	// The emulator accepts any caller and serves plaintext, so it is reached
	// without a credential and without TLS. A caller that also supplies TLS
	// settings gets them: an emulator behind a terminator is still TLS.
	if isEmulator(cfg) {
		options = append(options, option.WithoutAuthentication())
		if cfg.TLS == nil {
			options = append(options, option.WithGRPCDialOption(
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			))
		}
	}
	// TLS is required exactly when the caller supplies its settings; nil leaves
	// the client's transport, which verifies the server's certificate.
	if cfg.TLS != nil {
		tlsConfig, err := cfg.TLS(endpointHost(endpoint(cfg)))
		if err != nil {
			return nil, fmt.Errorf("failed to get TLS config: %w", err)
		}
		if tlsConfig == nil {
			return nil, errors.New("TLS configuration must not be nil")
		}
		options = append(options, option.WithGRPCDialOption(
			grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		))
	}
	if cfg.Dial != nil {
		options = append(options, option.WithGRPCDialOption(
			grpc.WithContextDialer(func(ctx context.Context, address string) (net.Conn, error) {
				return cfg.Dial(ctx, "tcp", address)
			}),
		))
	}
	return options, nil
}

// endpoint is the address the client is pointed at: host:port when the caller
// named a host, and "" to leave the client's own.
func endpoint(cfg driver.Config) string {
	// gRPC defaults to port 443 when no port is given.
	if cfg.Host == "" {
		return ""
	}
	if cfg.Port != "" {
		return net.JoinHostPort(cfg.Host, cfg.Port)
	}
	return cfg.Host
}

// endpointHost is the host of the endpoint, for the TLS settings of that host.
func endpointHost(endpoint string) string {
	if host, _, err := net.SplitHostPort(endpoint); err == nil {
		return host
	}
	return endpoint
}

// credentialOption turns the credential document the caller resolved into the
// Google Cloud authentication library's option for it. An empty document means
// the caller resolved none and leaves the library's default chain.
func credentialOption(credentialJSON string) (option.ClientOption, error) {
	if credentialJSON == "" {
		return nil, nil
	}
	credentialJSONBytes := []byte(credentialJSON)
	var cred struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(credentialJSONBytes, &cred); err != nil {
		return nil, fmt.Errorf("failed to parse GCP credential JSON: %w", err)
	}
	switch cred.Type {
	case "service_account":
		return option.WithAuthCredentialsJSON(option.ServiceAccount, credentialJSONBytes), nil
	case "external_account":
		return option.WithAuthCredentialsJSON(option.ExternalAccount, credentialJSONBytes), nil
	case "impersonated_service_account":
		return option.WithAuthCredentialsJSON(option.ImpersonatedServiceAccount, credentialJSONBytes), nil
	case "authorized_user":
		return option.WithAuthCredentialsJSON(option.AuthorizedUser, credentialJSONBytes), nil
	case "":
		return nil, errors.New(`GCP credential JSON missing "type" field`)
	default:
		return nil, fmt.Errorf("unsupported GCP credential type: %q", cred.Type)
	}
}

// errNoDatabase is what a data call answers for a driver opened without a
// database. Bytebase read this client as nil and crashed; no product path asked
// it anything, and the message is the one Bytebase's explain already used.
var errNoDatabase = errors.New("spanner client is not initialized, database name may be missing")

// Ping reads one database name, which is the least an instance answers: Spanner
// has no ping of its own.
func (c *liveClient) Ping(ctx context.Context) error {
	iter := c.admin.ListDatabases(ctx, &databasepb.ListDatabasesRequest{Parent: c.instancePath})
	_, err := iter.Next()
	if err == iterator.Done {
		return nil
	}
	return err
}

// Close releases the clients' connections and the sessions they hold.
func (c *liveClient) Close() error {
	if c.data != nil {
		c.data.Close()
	}
	return c.admin.Close()
}

// ListDatabases lists the instance's databases, in the order the admin API
// reports them.
func (c *liveClient) ListDatabases(ctx context.Context) ([]Database, error) {
	var databases []Database
	iter := c.admin.ListDatabases(ctx, &databasepb.ListDatabasesRequest{Parent: c.instancePath})
	for {
		database, err := iter.Next()
		if err == iterator.Done {
			return databases, nil
		}
		if err != nil {
			return nil, err
		}
		databases = append(databases, Database{
			Name:    database.Name,
			Dialect: database.DatabaseDialect.String(),
		})
	}
}

// GetDatabase reads database with its dialect.
func (c *liveClient) GetDatabase(ctx context.Context, database string) (Database, error) {
	got, err := c.admin.GetDatabase(ctx, &databasepb.GetDatabaseRequest{Name: database})
	if err != nil {
		return Database{}, err
	}
	return Database{Name: got.Name, Dialect: got.DatabaseDialect.String()}, nil
}

// UpdateDatabaseDdl submits statements and waits for the schema change to
// finish.
func (c *liveClient) UpdateDatabaseDdl(ctx context.Context, database string, statements []string) error {
	op, err := c.admin.UpdateDatabaseDdl(ctx, &databasepb.UpdateDatabaseDdlRequest{
		Database:   database,
		Statements: statements,
	})
	if err != nil {
		return err
	}
	return op.Wait(ctx)
}

// CreateDatabase creates a database and waits for it to exist.
func (c *liveClient) CreateDatabase(ctx context.Context, instancePath, createStatement string, extra []string) error {
	op, err := c.admin.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          instancePath,
		CreateStatement: createStatement,
		ExtraStatements: extra,
	})
	if err != nil {
		return err
	}
	_, err = op.Wait(ctx)
	return err
}

// GetDatabaseDdl reads database's DDL from the admin API.
func (c *liveClient) GetDatabaseDdl(ctx context.Context, database string) ([]string, error) {
	resp, err := c.admin.GetDatabaseDdl(ctx, &databasepb.GetDatabaseDdlRequest{Database: database})
	if err != nil {
		return nil, err
	}
	return resp.Statements, nil
}

// Read runs statement as one single-use read-only transaction. The session goes
// back to the client's pool when the iterator stops, so nothing is held between
// statements.
func (c *liveClient) Read(ctx context.Context, statement string) (RowIterator, error) {
	if c.data == nil {
		return nil, errNoDatabase
	}
	return &liveRowIterator{iter: c.data.Single().Query(ctx, spanner.NewStatement(statement))}, nil
}

// Snapshot returns a reader of one timestamp's reads, for a sync that must not
// see a schema change partway through.
func (c *liveClient) Snapshot(ctx context.Context) (Reader, error) {
	if c.data == nil {
		return nil, errNoDatabase
	}
	return &liveReader{tx: c.data.ReadOnlyTransaction()}, nil
}

// Update runs one DML statement in a read-write transaction.
func (c *liveClient) Update(ctx context.Context, statement string) (int64, error) {
	if c.data == nil {
		return 0, errNoDatabase
	}
	var rowCount int64
	_, err := c.data.ReadWriteTransaction(ctx, func(ctx context.Context, rwt *spanner.ReadWriteTransaction) error {
		count, err := rwt.Update(ctx, spanner.NewStatement(statement))
		if err != nil {
			return err
		}
		rowCount = count
		return nil
	})
	if err != nil {
		return 0, err
	}
	return rowCount, nil
}

// BatchUpdate runs statements as one read-write transaction's batch.
func (c *liveClient) BatchUpdate(ctx context.Context, statements []string) ([]int64, error) {
	if c.data == nil {
		return nil, errNoDatabase
	}
	var counts []int64
	_, err := c.data.ReadWriteTransaction(ctx, func(ctx context.Context, rwt *spanner.ReadWriteTransaction) error {
		spannerStmts := make([]spanner.Statement, 0, len(statements))
		for _, statement := range statements {
			spannerStmts = append(spannerStmts, spanner.NewStatement(statement))
		}
		batchCounts, err := rwt.BatchUpdate(ctx, spannerStmts)
		if err != nil {
			return err
		}
		counts = batchCounts
		return nil
	})
	if err != nil {
		return nil, err
	}
	return counts, nil
}

// PartitionedUpdate runs one DML statement as a partitioned update, outside any
// transaction.
func (c *liveClient) PartitionedUpdate(ctx context.Context, statement string) (int64, error) {
	if c.data == nil {
		return 0, errNoDatabase
	}
	return c.data.PartitionedUpdate(ctx, spanner.NewStatement(statement))
}

// AnalyzeQuery returns the plan statement would run under. A statement a
// single-use read-only transaction cannot read is analysed in a read-write
// transaction that is rolled back, so analysing it writes nothing.
func (c *liveClient) AnalyzeQuery(ctx context.Context, statement string, readWrite bool) (*sppb.QueryPlan, error) {
	if c.data == nil {
		return nil, errNoDatabase
	}
	if !readWrite {
		return c.data.Single().AnalyzeQuery(ctx, spanner.NewStatement(statement))
	}
	txn, err := spanner.NewReadWriteStmtBasedTransaction(ctx, c.data)
	if err != nil {
		return nil, err
	}
	// The transaction never commits; a failed rollback leaves it to the server's
	// own timeout.
	defer txn.Rollback(ctx)
	return txn.AnalyzeQuery(ctx, spanner.NewStatement(statement))
}

// liveReader is a snapshot of a database, held until Close.
type liveReader struct {
	tx *spanner.ReadOnlyTransaction
}

func (r *liveReader) Query(ctx context.Context, statement string) (RowIterator, error) {
	return &liveRowIterator{iter: r.tx.Query(ctx, spanner.NewStatement(statement))}, nil
}

func (r *liveReader) Close() {
	r.tx.Close()
}

// liveRowIterator is the client's row iterator, with its rows decoded into the
// seam's cells.
type liveRowIterator struct {
	iter *spanner.RowIterator
}

// Columns returns the row type Spanner stated with the first row, or with the
// end of the result.
func (r *liveRowIterator) Columns() []Column {
	fields := r.iter.Metadata.RowType.GetFields()
	columns := make([]Column, 0, len(fields))
	for _, field := range fields {
		columns = append(columns, Column{Name: field.Name, Type: field.Type})
	}
	return columns
}

func (r *liveRowIterator) Next() ([]Cell, error) {
	row, err := r.iter.Next()
	if err == iterator.Done {
		return nil, io.EOF
	}
	if err != nil {
		return nil, err
	}
	cells := make([]Cell, 0, row.Size())
	for i := 0; i < row.Size(); i++ {
		var column spanner.GenericColumnValue
		if err := row.Column(i, &column); err != nil {
			return nil, err
		}
		cells = append(cells, Cell{Type: column.Type, Value: column.Value})
	}
	return cells, nil
}

func (r *liveRowIterator) Stop() {
	r.iter.Stop()
}
