// Package driver is Cosmos DB's implementation of the driver contract in
// github.com/bytebase/omni/driver: it reaches an account through the azcosmos
// client, syncs the account's databases and a database's containers into
// omni/metadata, and runs the container query language for interactive
// queries.
//
// It implements Driver, Syncer and Querier. There is no Dump: omni's cosmosdb
// tree has no catalog, so nothing renders a synced schema as the engine's DDL,
// and Bytebase's driver dumped nothing either. There is no Executor: the query
// language has no statement that changes data, container items are written
// through the container API, and Bytebase's driver answered Execute as
// unimplemented.
//
// Config.Host is the account endpoint as the portal and the emulator write it,
// https://<account>.documents.azure.com:443/. Config.Database is the database a
// sync reads and a query runs in. The account is not versioned: SyncInstance
// reports no version because no query the account answers states one.
//
// Credentials arrive as values, named by Params["auth"]:
//
//	masterkey     Config.Password is the account key, which is also what the
//	              emulator's well-known key is
//	clientsecret  Config.User and Config.Password are an app registration's
//	              client ID and secret, Params["tenantId"] its tenant
//	default       azidentity's own chain: environment, workload identity,
//	              managed identity, then the Azure CLI
//
// Params["emulator"] = "true" says the endpoint is the Cosmos DB emulator,
// whose metadata the driver reads over the emulator's REST API. Bytebase
// selected that path from its own dev flag on a localhost endpoint; here the
// caller says it, and supplies the emulator's account key and a Config.TLS
// that accepts the certificate the emulator signs itself.
//
// Config.TLS supplies the TLS settings for the endpoint's host and makes TLS
// required; nil leaves the client's own transport, which is how a server's
// certificate is verified. Config.Dial opens the connections the client makes,
// so an SSH tunnel or a cloud connector reaches the account unchanged.
// Config.ApplicationName, Config.EngineVersion and Config.ReadOnly do nothing:
// the protocol carries no application name, the account has no version to
// report, and a container session has no read-only mode.
//
// A query names its container in QueryOptions.Schema. The query language is
// per container — a statement reads the container the caller already opened —
// and Cosmos DB has no statement form that opens another one, so the container
// is the query's scope rather than a name inside its text.
//
// QueryOptions.Explain, QueryOptions.MaskingEnabled and
// QueryOptions.MaxResultBytes do nothing: the container API has no plan to ask
// for, a row is one TEXT column of JSON with nothing to mask, and Bytebase's
// driver never capped a result's bytes. QueryOptions.Limit caps the rows the
// driver reads and is also the query's page size below 1000 rows.
//
// A query can fail at the gateway with the gateway's own error, which is an
// *azcore.ResponseError: a cross-partition query the gateway refuses is
// replayed through the client's Go query engine (queryengine/gonative), and a
// failure there surfaces the gateway's original error.
package driver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"

	"github.com/bytebase/omni/driver"
)

var (
	_ driver.Driver  = (*Driver)(nil)
	_ driver.Syncer  = (*Driver)(nil)
	_ driver.Querier = (*Driver)(nil)
)

// The names this driver reads out of driver.Config.Params.
const (
	// ParamAuth names the form of the credential in Config.User and
	// Config.Password.
	ParamAuth = "auth"
	// ParamTenantID is the tenant of a client-secret credential.
	ParamTenantID = "tenantId"
	// ParamEmulator says the endpoint is the Cosmos DB emulator, whose
	// metadata the driver lists over its REST API.
	ParamEmulator = "emulator"

	// AuthMasterKey is an account key, in Config.Password.
	AuthMasterKey = "masterkey"
	// AuthClientSecret is an Azure AD app registration: Config.User is its
	// client ID and Config.Password its secret.
	AuthClientSecret = "clientsecret"
	// AuthDefault is azidentity's default chain.
	AuthDefault = "default"
)

// Driver is an open Cosmos DB driver: one client to one account, and the
// database its syncs and queries read.
type Driver struct {
	client       Client
	databaseName string
}

// Open returns a driver for the account and database cfg names. It does not
// reach the account; Ping does.
func Open(_ context.Context, cfg driver.Config) (*Driver, error) {
	c, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Driver{client: c, databaseName: cfg.Database}, nil
}

// NewClient builds the client the driver talks to the account through: the
// azcosmos client, with the emulator's REST listing when the caller named the
// emulator. Open is this client with the account and database kept.
func NewClient(cfg driver.Config) (Client, error) {
	transport, err := newTransport(cfg)
	if err != nil {
		return nil, err
	}
	var options *azcosmos.ClientOptions
	if transport != nil {
		options = &azcosmos.ClientOptions{
			ClientOptions: azcore.ClientOptions{Transport: transport},
		}
	}

	var azClient *azcosmos.Client
	switch strings.ToLower(cfg.Params[ParamAuth]) {
	case AuthMasterKey:
		if cfg.Password == "" {
			return nil, errors.New("an account key must be set as the password")
		}
		credential, credErr := azcosmos.NewKeyCredential(cfg.Password)
		if credErr != nil {
			return nil, fmt.Errorf("failed to create key credential: %w", credErr)
		}
		azClient, err = azcosmos.NewClientWithKey(cfg.Host, credential, options)
	case AuthClientSecret:
		tenantID := cfg.Params[ParamTenantID]
		if tenantID == "" || cfg.User == "" || cfg.Password == "" {
			return nil, fmt.Errorf("%s needs %s, a user and a password", AuthClientSecret, ParamTenantID)
		}
		credential, credErr := azidentity.NewClientSecretCredential(tenantID, cfg.User, cfg.Password, nil)
		if credErr != nil {
			return nil, fmt.Errorf("failed to create client secret credential: %w", credErr)
		}
		azClient, err = azcosmos.NewClient(cfg.Host, credential, options)
	case "", AuthDefault:
		credential, credErr := azidentity.NewDefaultAzureCredential(nil)
		if credErr != nil {
			return nil, fmt.Errorf("failed to found default Azure credential: %w", credErr)
		}
		azClient, err = azcosmos.NewClient(cfg.Host, credential, options)
	default:
		return nil, fmt.Errorf("unknown %s %q: want %s, %s or %s",
			ParamAuth, cfg.Params[ParamAuth], AuthMasterKey, AuthClientSecret, AuthDefault)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create CosmosDB client: %w", err)
	}

	live := &liveClient{client: azClient}
	if isEmulator(cfg) {
		// The metadata listings below go over the emulator's REST API, whose
		// queries do not run through the client's query path. Its own
		// transport is the one the account client uses.
		rest, restErr := newEmulatorRESTClient(cfg, transport)
		if restErr != nil {
			return nil, restErr
		}
		live.emulator = rest
	}
	return live, nil
}

// isEmulator reports whether the caller said the endpoint is the emulator.
func isEmulator(cfg driver.Config) bool {
	switch strings.ToLower(cfg.Params[ParamEmulator]) {
	case "1", "true", "yes":
		return true
	}
	return false
}

// newTransport returns the HTTP transport the account is reached over, or nil
// for the client's own.
func newTransport(cfg driver.Config) (*http.Client, error) {
	if cfg.TLS == nil && cfg.Dial == nil {
		return nil, nil
	}
	transport := &http.Transport{DialContext: cfg.Dial}
	if cfg.TLS != nil {
		tlsConfig, err := cfg.TLS(endpointHost(cfg.Host))
		if err != nil {
			return nil, fmt.Errorf("failed to get TLS config: %w", err)
		}
		if tlsConfig == nil {
			return nil, errors.New("TLS configuration must not be nil")
		}
		transport.TLSClientConfig = tlsConfig
	}
	return &http.Client{Transport: transport}, nil
}

// endpointHost returns the host an endpoint URL names, for the TLS settings of
// that host.
func endpointHost(endpoint string) string {
	u, err := url.Parse(endpoint)
	if err != nil {
		return endpoint
	}
	return u.Hostname()
}

// Ping reads the account's database list as a connectivity check, which is
// what Bytebase's driver did.
func (d *Driver) Ping(ctx context.Context) error {
	if err := d.client.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping CosmosDB: %w", err)
	}
	return nil
}

// Close closes the driver. The client holds no session to release, so this
// releases nothing, as Bytebase's driver did.
func (*Driver) Close() error {
	return nil
}
