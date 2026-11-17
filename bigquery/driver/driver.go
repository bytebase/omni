// Package driver is BigQuery's implementation of the driver contract in
// github.com/bytebase/omni/driver: it reaches the BigQuery service through the
// cloud.google.com/go/bigquery client, syncs a project's datasets and one
// dataset's schema into omni/metadata, and runs interactive queries and change
// scripts.
//
// # Unverified against a project
//
// No test in this package has ever asked BigQuery. The package is a port of
// Bytebase's backend/plugin/db/bigquery, whose tests reached one pure helper and
// no service, and BigQuery has no image omni can start in a container and no
// emulator that answers its REST API: every answer the sync, query and execute
// paths expect here was written from that driver and from BigQuery's
// documentation, and none was recorded. There are no tapes, no sync goldens, no
// dump fixtures and no pin row, because nothing here can produce one. VERIFY.md
// holds the recording plan: the project and account to pin, the corpus, the
// tests to record in what order, and the parity harness that holds this package
// to the driver it replaces. Treat a behavior this package states as unconfirmed
// until that recording lands.
//
// Its tests therefore cover only what needs no project: the client
// configuration and the credential types it accepts, statement classification,
// the result-limit rewrite, the value conversions, the query-plan rendering, the
// column mapping and the not-found mapping. Everything that reaches the service
// is behind one environment gate: BIGQUERY_CREDENTIALS names the
// service-account key a live test uses, with BIGQUERY_PROJECT and
// BIGQUERY_DATASET beside it. Unset, a live test skips naming the gate; set and
// incomplete, it fails, so a recording cannot read as a pass.
//
// # Where it lives
//
// The statements are GoogleSQL, which omni's googlesql tree parses, classifies
// and wraps a LIMIT in, and this package is the conversation with BigQuery's
// service, so it lives in a tree of its own beside googlesql, as spanner/driver
// does and dynamodb/driver does beside partiql. Bytebase's alias table maps
// BigQuery onto googlesql for the parser and the analysis, not for a service
// client: a googlesql/driver would hold BigQuery's client and Spanner's, and
// googlesql holds no catalog for either. The driver calls googlesql/analysis for
// its judgment rather than carrying a copy of it.
//
// # What it implements
//
// It implements Driver, Syncer, Querier and Executor, and Dumper as a no-op.
// Dump writes nothing: googlesql has no catalog that turns a synced schema back
// into DDL, and Bytebase's BigQuery Dump wrote nothing either, so a caller that
// dumps a BigQuery database writes an empty dump, which is what the product does
// today. VERIFY.md says what a metadata-rendered dump would need.
//
// There is no row-estimate capability. BigQuery has no EXPLAIN to read an
// estimate from, and Bytebase's driver answered no affected-row estimate ahead
// of a change.
//
// # The connection
//
// Config.Params[ParamProject] names the Google Cloud project and is required.
// Config.Database is the dataset the driver is open on: it is the default
// dataset of every statement the driver sends, the dataset SyncDatabase reads,
// and the name a synced database reports. An empty Config.Database leaves
// statements unqualified, which BigQuery refuses for a table named without its
// dataset, and SyncDatabase reports driver.ErrNotFound.
//
// Config.Password is the resolved credential: a Google Cloud credential JSON
// document, of the kind a service account, an external account, an impersonated
// service account or an authorized user is written as. Its "type" member decides
// which. An empty Config.Password means the caller resolved no credential of its
// own and the Google authentication library's default chain reads one (the
// environment, the well-known file, an attached service account). Credential
// resolution of every other kind stays with the caller: a service account a
// caller impersonates is the caller's to impersonate, and what arrives here is
// the credential it resolved. Config.User does nothing: BigQuery has no login
// beside the credential.
//
// Config.Host and Config.Port optionally name the endpoint, which is how a
// Private Service Connect endpoint or a proxy is reached; an empty Config.Host
// leaves the client's own bigquery.googleapis.com. Config.TLS and Config.Dial do
// nothing: the client library owns its transport, which encrypts with the
// credential's own path and verifies the endpoint's certificate against the
// system roots, and it holds no dialer a caller can replace. A caller that must
// reach BigQuery through a tunnel forwards a local port and passes that port as
// Config.Host and Config.Port, as hive's driver asks for its own client.
//
// Config.ApplicationName, Config.ReadOnly and Config.EngineVersion do nothing:
// BigQuery's protocol carries no application name, a dataset has no read-only
// session to open, and BigQuery reports no version — SyncInstance therefore
// reports none, as Bytebase's driver did.
//
// # Queries
//
// A session is the driver: BigQuery has no session state a statement can set,
// because every statement runs as a job of its own. QueryOptions.Schema and
// QueryOptions.MaskingEnabled do nothing: a name resolves in Config.Database and
// the project it belongs to, and a row's cells are the columns BigQuery's own
// schema declares, which the driver does not trace through a statement.
//
// QueryOptions.Explain is a dry run rather than a plan: BigQuery has no EXPLAIN,
// and what it can answer before running a statement is whether the statement is
// valid and how many bytes it would process. Each result then carries one STRING
// column named "Validation Result" whose row states that answer, which is the
// shape Bytebase's driver returned. QueryOptions.ExplainFormat does nothing.
//
// QueryOptions.Limit rewrites each statement that returns rows through
// googlesql/analysis.StatementWithResultLimit, which wraps it in a CTE — BigQuery
// has no LIMIT to splice into an arbitrary query — and leaves a statement that
// returns no rows as written. Which statements return rows is
// googlesql/analysis.ClassifySQL's verdict for the BigQuery dialect.
//
// QueryOptions.MaxResultBytes stops reading a statement's rows once their
// encoded size passes it, sets Result.Truncated, and keeps the row that went
// over, as omni's other drivers do.
//
// A row's cells are the types BigQuery's schema states: a STRING, BYTES, INTEGER,
// FLOAT and BOOLEAN cell arrives as that kind. Every other type — TIMESTAMP,
// DATE, TIME, DATETIME, NUMERIC, BIGNUMERIC, GEOGRAPHY, INTERVAL, JSON, RANGE, a
// RECORD and an ARRAY — falls back to Go's own rendering of the value, which is
// what Bytebase's driver returned and not BigQuery's literal; VERIFY.md records
// that gap for the recording to settle.
//
// A result that came from a job carries the plan of the run that returned its
// rows in Result.Plan, as JSON, with Executed true: BigQuery reports the stages
// a query was executed as rather than planning the next run, which is what its
// jobs API states. Reading the plan costs one job read after the rows; a read
// that came from no job of its own, or whose job states no plan, carries none.
//
// The operator label is not ported. Bytebase attached the interactive query's
// operator email to each job as the operator_email label for cost attribution,
// and driver.QueryOptions has no operator identity to carry, so no label is set.
// A caller that needs one adds the field to the contract rather than this driver
// inventing it.
//
// # Change scripts
//
// Execute sends the whole script as one BigQuery job, which is what Bytebase's
// driver did and what BigQuery's scripting runs as: a script with several
// statements is one job whose status reports the script and carries no per
// statement row count. It reports one Command spanning the script to
// ExecuteOptions.Log, before and after the job, so the execution log holds the
// statement and the engine's error. ExecuteOptions.Transaction, .Isolation and
// .MaxRetries do nothing: BigQuery has no multi-statement transaction that a
// script's statements share, no isolation level to set, and its client retried
// nothing Bytebase asked it to. ExecuteOptions.CreateDatabase changes nothing:
// BigQuery's creation statement is a CREATE SCHEMA or a CREATE TABLE, which the
// same job runs.
//
// A command's affected rows are the job statistics' DML row count, which
// BigQuery reports for INSERT, UPDATE, DELETE and MERGE and states as zero for
// DDL and for a multi-statement script.
//
// # Errors
//
// A statement the service refused is the client's own error, as the client
// returned it: *bigquery.Error for a job whose status reports a failure, whose
// Location, Message and Reason are the service's fields, and *googleapi.Error
// for a request the API refused. SyncDatabase's dataset read maps a not-found
// answer to driver.ErrNotFound, whichever of the two carries it, so a caller can
// tell a dataset that is not there from a request that failed.
package driver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"

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

// Driver is an open BigQuery driver: one client to one project, and the dataset
// it is open on.
type Driver struct {
	client *bigquery.Client

	// project is the Google Cloud project every request is signed for.
	project string
	// databaseName is the dataset the driver is open on.
	databaseName string
}

// Open returns a driver for the project cfg names. Open does not reach the
// service; Ping does.
func Open(ctx context.Context, cfg driver.Config) (*Driver, error) {
	built, err := clientConfigOf(cfg)
	if err != nil {
		return nil, err
	}
	client, err := bigquery.NewClient(ctx, built.project, built.options...)
	if err != nil {
		return nil, err
	}
	return &Driver{
		client:       client,
		project:      built.project,
		databaseName: cfg.Database,
	}, nil
}

// Ping runs one statement as a connectivity check.
func (d *Driver) Ping(ctx context.Context) error {
	if _, err := d.client.Query("SELECT 1").Read(ctx); err != nil {
		return fmt.Errorf("bigquery: bad connection: %w", err)
	}
	return nil
}

// Close closes the driver and the connections its client opened.
func (d *Driver) Close() error {
	return d.client.Close()
}

// SyncInstance lists the project's datasets. Each is a database; a dataset's
// schemas are read by SyncDatabase. The list is ordered by name, because the
// service's own order is not part of an answer and a caller compares two syncs.
//
// No version and no roles are reported: BigQuery has neither a version to ask
// for nor a role the API states, which is what Bytebase's driver answered.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	var databases []*metadatapb.DatabaseSchemaMetadata
	it := d.client.Datasets(ctx)
	for {
		dataset, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		databases = append(databases, &metadatapb.DatabaseSchemaMetadata{Name: dataset.DatasetID})
	}
	sort.Slice(databases, func(i, j int) bool { return databases[i].Name < databases[j].Name })
	return &driver.Instance{Databases: databases}, nil
}

// query returns a query job for statement, on the dataset the driver is open on.
func (d *Driver) query(statement string) *bigquery.Query {
	q := d.client.Query(statement)
	q.DefaultDatasetID = d.databaseName
	return q
}

// statisticsOf returns the query statistics of a finished job. JobStatus
// .Statistics is nil when the job read carried no statistics block, so the
// Details assertion alone would dereference a nil.
func statisticsOf(status *bigquery.JobStatus) (*bigquery.QueryStatistics, bool) {
	if status == nil || status.Statistics == nil {
		return nil, false
	}
	stats, ok := status.Statistics.Details.(*bigquery.QueryStatistics)
	return stats, ok
}

// notFound reports whether err is the service saying an object the call named is
// not there: the API's 404, or a job the service failed with its notFound
// reason, which is how a statement against a missing dataset fails.
func notFound(err error) bool {
	var apiError *googleapi.Error
	if errors.As(err, &apiError) {
		return apiError.Code == http.StatusNotFound
	}
	var bigQueryError *bigquery.Error
	if errors.As(err, &bigQueryError) {
		return bigQueryError.Reason == "notFound"
	}
	return false
}

// notFoundError wraps err as driver.ErrNotFound when it is a not-found answer,
// and returns it unchanged otherwise.
func notFoundError(err error) error {
	if !notFound(err) {
		return err
	}
	return fmt.Errorf("%w: %w", driver.ErrNotFound, err)
}
