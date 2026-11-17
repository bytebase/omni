// Package driver is DynamoDB's implementation of the driver contract in
// github.com/bytebase/omni/driver: it reaches a DynamoDB endpoint through the
// AWS SDK, syncs the region's tables into omni/metadata, and runs PartiQL for
// interactive queries and change scripts.
//
// It implements Driver, Syncer, Dumper, Querier and Executor. Dump writes
// nothing: a synced table carries its key attributes and its secondary
// indexes' key attributes, but not the attribute types, billing mode or
// throughput a create takes, so a dump of the metadata could not create the
// table it describes, and Bytebase's driver answered Dump by writing nothing at
// all. A caller that dumps a DynamoDB database therefore writes an empty dump,
// which is what the product does today.
//
// Config.Host and Config.Port name the endpoint: a host with no scheme is
// reached over http, or over https when Config.TLS is set, and the port is
// appended when the host does not already carry one. Config.Host may also be a
// complete endpoint URL. An empty host leaves the AWS SDK's own resolution for
// the region, which is how real AWS is reached; anything else is a
// DynamoDB-compatible endpoint such as DynamoDB Local.
//
// Config.Params carry what the AWS SDK needs and the contract has no field for:
//
//	region        the AWS region. The SDK signs every request for a region and
//	              requires one even when the endpoint is a local one.
//	sessionToken  the session token of a temporary credential, when
//	              Config.User and Config.Password are one.
//	accountId     the AWS account ID SyncInstance names the pseudo database
//	              from. When it is absent, SyncInstance asks STS for the
//	              account the credential belongs to, as Bytebase's driver did.
//
// Config.User and Config.Password are the access key ID and secret access key.
// An empty Config.User means the caller resolved the credential elsewhere, and
// the SDK's default chain reads it (environment, shared configuration, an
// instance role). Credential resolution of every other kind stays with the
// caller: a role asserted through STS is the caller's to assume, and what
// arrives here is the assumed role's key pair.
//
// The account ID is the one fact DynamoDB itself does not answer: the service
// identifies an account by the credential, not by an API. SyncInstance asks STS
// for it when the caller passes none, as Bytebase's driver did, so a cutover
// keeps the account ID Bytebase reported and the same {account_id}-{region}
// name. STS is a second service the endpoint in Config does not name, and the
// call reaches STS's own endpoint; a caller that already knows the account
// passes Params[ParamAccountID] and no STS call is made. A failed identity call
// is reported as Bytebase reported it, "failed to get caller identity: ...",
// and an identity carrying no account as "account id is empty in the caller
// identity".
//
// Config.TLS supplies the TLS settings for the endpoint's host and makes TLS
// required, so an endpoint URL that names http is refused rather than reached
// in plaintext. Nil leaves the SDK's own transport, which verifies the
// endpoint's certificate against the system roots and the endpoint's host name.
// Verification is therefore off only where the caller's *tls.Config says
// InsecureSkipVerify: the driver reads no product's verify toggle and infers no
// skip, so a caller whose previous client always verified, as Bytebase's did,
// keeps verifying unless its user asked to skip, rather than skipping because a
// flag is unset. Config.Dial opens the connections the client makes, so an SSH
// tunnel or a cloud connector reaches the endpoint unchanged.
// Config.Database is the pseudo database name the caller recorded, reported by
// SyncDatabase and read by nothing else: DynamoDB has no databases, so a sync
// reads every table of the region whatever that name says.
//
// Config.ApplicationName, Config.ReadOnly and Config.EngineVersion do nothing:
// the protocol carries no application name, a PartiQL statement has no
// read-only session, and DynamoDB reports no version.
//
// A session is the driver: DynamoDB has no session state. QueryOptions.Limit
// caps the rows a statement's result holds: the driver reads no further page
// once it has them, and each call asks only for the rows still missing.
// QueryOptions.MaxResultBytes stops a result once its rows' encoded size passes
// the cap and sets Result.Truncated. QueryOptions.Explain is an error, as
// Bytebase answered it, since PartiQL has no plan to ask for.
// QueryOptions.Schema and QueryOptions.MaskingEnabled do nothing: there is no
// schema per session, and a row is a map of attributes the driver does not
// trace columns through.
//
// A row's cells are the item's attributes: a string or a number as its text, a
// binary as bytes, a boolean as a boolean, a NULL as a null, and a set, a list
// or a map as its JSON text, where a binary is base64.
//
// A PartiQL script is split by partiql.Split, so statement boundaries are
// omni's partiql judgment and not this package's, and a segment holding no
// statement is skipped. Execute runs each statement as its own
// ExecuteTransaction, which is why ExecuteOptions.Transaction, .Isolation and
// .CreateDatabase do nothing: DynamoDB has no multi-statement transaction the
// PartiQL statements of one script can share, no isolation level to set, and no
// database to create. ExecuteOptions.MaxRetries retries nothing: a retry of the
// lock conflicts DynamoDB reports is the caller's.
//
// An error a statement's answer carries is the SDK's own, as the client
// returned it, so a caller that reads the SDK's error codes gets them.
package driver

import (
	"context"

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
	// ParamRegion is the AWS region the SDK signs requests for.
	ParamRegion = "region"
	// ParamSessionToken is the session token of a temporary credential.
	ParamSessionToken = "sessionToken"
	// ParamAccountID is the AWS account ID SyncInstance builds the pseudo
	// database name from.
	ParamAccountID = "accountId"
)

// Driver is an open DynamoDB driver: one client to one endpoint, and the
// pseudo database name its syncs report.
type Driver struct {
	client       client
	databaseName string
	accountID    string
	region       string
}

// Open returns a driver for the endpoint cfg names. It does not reach the
// endpoint; Ping does.
func Open(ctx context.Context, cfg driver.Config) (*Driver, error) {
	client, err := newClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Driver{
		client:       client,
		databaseName: cfg.Database,
		accountID:    cfg.Params[ParamAccountID],
		region:       cfg.Params[ParamRegion],
	}, nil
}

// Ping reads one table name as a connectivity check, which is what Bytebase's
// driver did: DynamoDB has no ping of its own.
func (d *Driver) Ping(ctx context.Context) error {
	return d.client.Ping(ctx)
}

// Close closes the driver and the connections its client opened.
func (d *Driver) Close() error {
	return d.client.Close()
}
