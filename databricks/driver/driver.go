// Package driver is Databricks' implementation of the driver contract in
// github.com/bytebase/omni/driver: it reaches a Databricks workspace through
// the Databricks SDK, syncs Unity Catalog's catalogs, schemas and tables into
// omni/metadata, runs statements through a SQL warehouse, and dumps a
// workspace's DDL as the server states it.
//
// Everything that reaches a workspace is behind one environment gate:
// DATABRICKS_DSN names the workspace, the personal access token, the SQL
// warehouse and the catalog a live test uses. Unset, a live test skips naming
// the gate; set and naming no catalog, it fails, so a recording cannot write no
// tape and read as a pass.
//
// Statement splitting, EXPLAIN rewriting and inspection, row-returning
// statement classification, and view DDL target qualification live in
// databricks/analysis, backed by databricks/parser.
//
// # What it implements
//
// The package was verified against a scratch workspace on 2026-11-17: the
// live tests passed, the recordings and a sync golden are under testdata, and
// the parity harness in VERIFY.md agreed with Bytebase's driver on the sync,
// the instance, queries and Execute. Dump agrees except where the workspace
// showed Bytebase wrong: a materialized view is the server's own
// SHOW CREATE TABLE.
//
// It implements Driver, Syncer, Dumper, Querier and Executor. SyncDatabase reads
// a catalog's schemas and tables through Unity Catalog. SyncInstance reads the
// SQL version and the catalog names. Dump asks the server to state each object
// with SHOW CREATE TABLE, because omni has no databricks catalog to render a
// synced schema as Databricks DDL; the metadata argument is not read, as in
// hive/driver, so a caller that edits or narrows it changes nothing. A view's
// server statement names its schema only, so Dump qualifies the target with the
// catalog. Bytebase's Dump wrote the same server statements, except that it
// rendered a materialized view from DESC FORMATTED, which is not a statement
// the workspace accepts, and left a view's target unqualified.
//
// # The connection
//
// Open does not reach the workspace. The SDK's configuration resolution would
// otherwise fetch /.well-known/databricks-config on a best-effort basis, so Open
// installs a resolver that reads no host metadata; Ping is the first request.
//
// Config.Host is the workspace URL, such as
// https://dbc-00000000-0000.cloud.databricks.com. Config.Password is the
// personal access token; a token is a resolved credential, as the contract
// says, and this driver supports token authentication only, which is what
// Bytebase's driver supported. Config.Database is the catalog unqualified names
// resolve in and SyncDatabase reads; empty leaves the workspace's default
// catalog. Config.Params[ParamWarehouseID] is the SQL warehouse every statement
// runs on; "warehouse_id" is accepted as an alias. Config.User does nothing: a
// personal access token has no login name beside it. OAuth credentials are the
// caller's to exchange; what arrives here is the token.
//
// Config.TLS supplies the TLS settings for the workspace's host and makes TLS
// required, so a caller can pin a certificate. Config.Dial opens the
// connections the SDK's HTTP client makes, so an SSH tunnel or a cloud
// connector reaches the workspace unchanged. Both are installed on the SDK's
// HTTP transport; nil leaves the SDK's own transport and its system-root
// verification. Config.ApplicationName does nothing: the service carries no
// application name. Config.ReadOnly does nothing: a SQL warehouse has no
// read-only session setting, and Bytebase ignored the flag. Config.EngineVersion
// does nothing: SyncInstance asks SELECT VERSION(), and nothing version-
// dependent follows it.
//
// # Queries and change scripts
//
// A Session is the driver: the SQL Statement Execution API has no session
// state, so a statement's Catalog and Schema are sent with the statement and
// what it sets does not hold for the next. QueryOptions.Schema names the schema
// sent with each statement of the call. QueryOptions.Limit is the request's
// RowLimit, as Bytebase sent it, rather than a rewrite of the statement.
// QueryOptions.MaxResultBytes stops reading once the rows kept pass the cap and
// sets Result.Truncated. QueryOptions.Explain prefixes EXPLAIN FORMATTED, or
// leaves a statement that already asks for a plan as the caller wrote it;
// Result.Plan is TEXT, and Executed is true only for EXPLAIN ANALYZE.
// QueryOptions.MaskingEnabled does nothing: a Databricks result is the columns
// its statement traces.
//
// A row's cells are the strings the JSON_ARRAY format states, decoded by the
// column's type: BOOLEAN a bool, BINARY bytes, SHORT and INT an int32, LONG an
// int64, FLOAT a float32 widened to float64, DOUBLE a float64, and every other
// type its text, which is how Bytebase read a DECIMAL, DATE, TIMESTAMP, ARRAY,
// MAP and STRUCT. A JSON null arrives as the empty string, so an empty value in
// any non-STRING column is NULL, exactly as Bytebase read it; an empty STRING is
// the empty string. A result the service splits into chunks is read to the end,
// or until MaxResultBytes stops it. Result.AffectedRows stays nil because the
// API states a write's count as a num_affected_rows result column rather than
// as a field; Execute sums those columns.
//
// Execute splits the script and sends its statements in order; a failing
// statement ends the script and leaves the statements before it committed.
// Databricks' Statement Execution API exposes no transaction control, so
// ExecuteOptions.Transaction, .Isolation, .CreateDatabase and .MaxRetries
// change nothing, which is also what Bytebase's driver did. ExecuteOptions.Log
// receives each command and its response, with the rows its num_affected_rows
// states, or 0 when it states none.
//
// A failure the workspace's REST API reports is a *apierr.APIError, which this
// package wraps without replacing, and driver.ErrNotFound is joined when the
// API answers 404. A statement the engine refuses is the SDK helper's error
// text, "FAILED: <ERROR_CODE> <message>"; the helper discards the typed
// ServiceError, and this package keeps its text as Bytebase received it.
package driver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/apierr"
	sdkconfig "github.com/databricks/databricks-sdk-go/config"
	"github.com/databricks/databricks-sdk-go/service/workspace"

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
	// ParamWarehouseID is the SQL warehouse every statement runs on.
	ParamWarehouseID = "warehouseId"
	// paramWarehouseIDAlias is the spelling Bytebase's DataSource field uses.
	paramWarehouseIDAlias = "warehouse_id"
)

// Driver is an open Databricks driver: one workspace client, the SQL warehouse
// statements run on, and the catalog its syncs and unqualified names use.
type Driver struct {
	client      *databricks.WorkspaceClient
	warehouseID string
	catalog     string
}

// Open returns a driver for the workspace and catalog cfg names. It builds the
// SDK client but does not reach the workspace; Ping does.
func Open(_ context.Context, cfg driver.Config) (*Driver, error) {
	workspaceConfig, warehouseID, err := workspaceConfig(cfg)
	if err != nil {
		return nil, err
	}
	client, err := databricks.NewWorkspaceClient(workspaceConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace client: %w", err)
	}
	return &Driver{
		client:      client,
		warehouseID: warehouseID,
		catalog:     cfg.Database,
	}, nil
}

// workspaceConfig maps the contract's settings onto the SDK's and returns the
// warehouse ID, which the SDK config carries only as a convenience.
func workspaceConfig(cfg driver.Config) (*databricks.Config, string, error) {
	host := strings.TrimRight(strings.TrimSpace(cfg.Host), "/")
	if host == "" {
		return nil, "", errors.New("host must be set: Config.Host names the Databricks workspace URL")
	}
	if cfg.TLS != nil && strings.HasPrefix(strings.ToLower(host), "http://") {
		return nil, "", errors.New("TLS is required, so Config.Host must not use http://")
	}
	token := strings.TrimSpace(cfg.Password)
	if token == "" {
		return nil, "", errors.New("token must be set: Config.Password carries the personal access token")
	}
	warehouseID := strings.TrimSpace(param(cfg.Params, ParamWarehouseID))
	if warehouseID == "" {
		warehouseID = strings.TrimSpace(param(cfg.Params, paramWarehouseIDAlias))
	}
	if warehouseID == "" {
		return nil, "", fmt.Errorf("warehouse ID must be set: Config.Params[%q]", ParamWarehouseID)
	}

	workspaceConfig := &databricks.Config{
		Host:        host,
		Token:       token,
		WarehouseID: warehouseID,
		// The SDK's best-effort host metadata fetch is a request during
		// EnsureResolved; Open promises not to reach the workspace.
		HostMetadataResolver: func(context.Context, string) (*sdkconfig.HostMetadata, error) {
			return nil, nil
		},
	}
	if cfg.Dial != nil || cfg.TLS != nil {
		transport, err := httpTransport(cfg, host)
		if err != nil {
			return nil, "", err
		}
		workspaceConfig.HTTPTransport = transport
	}
	return workspaceConfig, warehouseID, nil
}

// httpTransport builds the SDK's HTTP transport from the caller's dialer and
// TLS settings. A caller's TLS is bound to the workspace host it was asked for;
// a redirect to another host is not expected on this API.
func httpTransport(cfg driver.Config, host string) (*http.Transport, error) {
	transport, ok := http.DefaultTransport.(*http.Transport)
	if ok {
		transport = transport.Clone()
	} else {
		transport = &http.Transport{}
	}
	if cfg.Dial != nil {
		transport.DialContext = cfg.Dial
	}
	if cfg.TLS != nil {
		hostname, err := hostnameOf(host)
		if err != nil {
			return nil, err
		}
		tlsConfig, err := cfg.TLS(hostname)
		if err != nil {
			return nil, err
		}
		if tlsConfig == nil {
			return nil, errors.New("TLS configuration must not be nil")
		}
		transport.TLSClientConfig = tlsConfig
	}
	return transport, nil
}

// hostnameOf returns the host a workspace URL names, accepting a URL without a
// scheme as the https URL the SDK would build.
func hostnameOf(host string) (string, error) {
	if !strings.Contains(host, "://") {
		host = "https://" + host
	}
	parsed, err := url.Parse(host)
	if err != nil {
		return "", fmt.Errorf("invalid host %q: %w", host, err)
	}
	if parsed.Hostname() == "" {
		return "", fmt.Errorf("invalid host %q: no hostname", host)
	}
	return parsed.Hostname(), nil
}

// param reads a connection parameter by name, case-insensitively.
func param(params map[string]string, name string) string {
	for key, value := range params {
		if strings.EqualFold(key, name) {
			return value
		}
	}
	return ""
}

// Ping asks the workspace for its root folder, which is what Bytebase's driver
// used as a connectivity check: Databricks has no SQL-level ping.
func (d *Driver) Ping(ctx context.Context) error {
	if _, err := d.client.Workspace.ListAll(ctx, workspace.ListWorkspaceRequest{Path: "/"}); err != nil {
		return wrapError("failed to ping instance", err)
	}
	return nil
}

// Close releases nothing: the SDK client keeps no connection of its own.
func (*Driver) Close() error {
	return nil
}

// wrapError adds the operation to err, and joins driver.ErrNotFound when the
// workspace's API answered that the resource is missing.
func wrapError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if apierr.IsMissing(err) {
		return fmt.Errorf("%w: %s: %w", driver.ErrNotFound, operation, err)
	}
	return fmt.Errorf("%s: %w", operation, err)
}
