package driver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awssdk "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/bytebase/omni/driver"
)

// The driver asks DynamoDB for the tables of a region, one table's key schema,
// and PartiQL statements. These are the calls it makes, and the seam its tests
// answer from a tape.
type client interface {
	Ping(ctx context.Context) error
	Close() error
	// ListTables returns at most limit names that sort after
	// exclusiveStartTableName, and the name to continue from, empty when the
	// listing is complete.
	ListTables(ctx context.Context, exclusiveStartTableName string, limit int32) ([]string, string, error)
	DescribeTable(ctx context.Context, tableName string) (*Table, error)
	ExecuteStatement(ctx context.Context, request StatementRequest) (StatementAnswer, error)
	ExecuteTransaction(ctx context.Context, statement string) error
}

// Table is what a sync reads of one table: the parts of DescribeTable's answer
// SyncDatabase keeps.
type Table struct {
	Name      string
	Status    string
	ItemCount int64
	KeySchema []KeySchemaElement
}

// KeySchemaElement is one key attribute and the role it plays in a table's
// key: HASH for the partition key, RANGE for the sort key.
type KeySchemaElement struct {
	AttributeName string
	KeyType       string
}

// The values KeySchemaElement.KeyType and Table.Status hold, as DynamoDB
// spells them.
const (
	keyTypeHash       = "HASH"
	keyTypeRange      = "RANGE"
	tableStatusActive = "ACTIVE"
)

// StatementRequest is one PartiQL statement as the driver sends it. Limit is
// the most items one call evaluates, 0 for no cap; NextToken resumes the read
// where the answer before stopped.
type StatementRequest struct {
	Statement string
	Limit     int32
	NextToken string
}

// StatementAnswer is what one execute-statement call answered: the items it
// read and the token that resumes the read, empty at its end.
type StatementAnswer struct {
	Items     []map[string]types.AttributeValue
	NextToken string
}

// newClient is the live client: the SDK client, reached over Config's
// endpoint, region and credentials. The driver's tests record through it, so a
// recording asks the pinned endpoint with the code a deployed driver runs.
func newClient(ctx context.Context, cfg driver.Config) (client, error) {
	awsConfig, err := newAWSConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	httpClient, err := newHTTPClient(cfg)
	if err != nil {
		return nil, err
	}
	options := []func(*awssdk.Options){}
	if httpClient != nil {
		options = append(options, func(o *awssdk.Options) {
			o.HTTPClient = httpClient
		})
	}
	// Honor a custom endpoint (DynamoDB Local, or any DynamoDB-compatible
	// service) when the caller names a host. Without this the SDK resolves the
	// real AWS endpoint for the region and the configured host:port is
	// silently ignored, routing queries intended for a local endpoint to AWS.
	if endpoint := endpoint(cfg); endpoint != "" {
		options = append(options, func(o *awssdk.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}
	return &liveClient{
		client:     awssdk.NewFromConfig(awsConfig, options...),
		httpClient: httpClient,
	}, nil
}

// newAWSConfig is the AWS configuration the SDK client is built from: the
// region Config names, and Config's credentials when it carries a key pair.
func newAWSConfig(ctx context.Context, cfg driver.Config) (aws.Config, error) {
	options := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.Params[ParamRegion]),
	}
	// Only use static credentials when an access key is set. A caller that
	// leaves Config.User empty has resolved the credential elsewhere (an
	// instance role, a shared profile, an assumed role's temporary keys) and
	// the SDK's default chain reads it.
	if cfg.User != "" {
		options = append(options, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.User, cfg.Password, cfg.Params[ParamSessionToken]),
		))
	}
	awsConfig, err := awsconfig.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config: %w", err)
	}
	return awsConfig, nil
}

// endpoint returns the DynamoDB endpoint URL Config names, or "" to let the
// AWS SDK resolve the region's own endpoint.
func endpoint(cfg driver.Config) string {
	host := cfg.Host
	if host == "" {
		// No host configured: use the AWS SDK's default endpoint for the region.
		return ""
	}
	port := cfg.Port
	// If the host already carries a scheme, treat it as a full base URL and
	// only append the port when one isn't already present.
	if strings.Contains(host, "://") {
		if port != "" && !strings.Contains(host[strings.Index(host, "://")+len("://"):], ":") {
			return host + ":" + port
		}
		return host
	}
	scheme := "http"
	if cfg.TLS != nil {
		scheme = "https"
	}
	if port != "" {
		return fmt.Sprintf("%s://%s:%s", scheme, host, port)
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}

// newHTTPClient returns the HTTP client the endpoint is reached over, or nil
// for the SDK's own.
func newHTTPClient(cfg driver.Config) (*http.Client, error) {
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
	if err != nil || u.Hostname() == "" {
		// A caller that wrote a bare host has no URL to parse, so what it wrote
		// is the host.
		return endpoint
	}
	return u.Hostname()
}

// liveClient is the client the driver opens against a DynamoDB endpoint.
type liveClient struct {
	client     *awssdk.Client
	httpClient *http.Client
}

// Ping reads one table name, which is the least a DynamoDB endpoint answers:
// DynamoDB has no ping of its own.
func (c *liveClient) Ping(ctx context.Context) error {
	// To avoid network overhead, we set the limit to 1.
	if _, err := c.client.ListTables(ctx, &awssdk.ListTablesInput{
		Limit: new(int32(1)),
	}); err != nil {
		return fmt.Errorf("failed to list dynamodb tables: %w", err)
	}
	return nil
}

// Close releases the connections the client opened. The SDK client holds no
// connection of its own.
func (c *liveClient) Close() error {
	if c.httpClient != nil {
		c.httpClient.CloseIdleConnections()
	}
	return nil
}

func (c *liveClient) ListTables(ctx context.Context, exclusiveStartTableName string, limit int32) ([]string, string, error) {
	input := &awssdk.ListTablesInput{Limit: new(limit)}
	if exclusiveStartTableName != "" {
		input.ExclusiveStartTableName = aws.String(exclusiveStartTableName)
	}
	out, err := c.client.ListTables(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list tables: %w", err)
	}
	return out.TableNames, aws.ToString(out.LastEvaluatedTableName), nil
}

func (c *liveClient) DescribeTable(ctx context.Context, tableName string) (*Table, error) {
	out, err := c.client.DescribeTable(ctx, &awssdk.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe table: %s: %w", tableName, err)
	}
	if out.Table == nil {
		return nil, fmt.Errorf("table not found: %s", tableName)
	}
	table := &Table{
		Name:      aws.ToString(out.Table.TableName),
		Status:    string(out.Table.TableStatus),
		ItemCount: aws.ToInt64(out.Table.ItemCount),
	}
	for _, key := range out.Table.KeySchema {
		table.KeySchema = append(table.KeySchema, KeySchemaElement{
			AttributeName: aws.ToString(key.AttributeName),
			KeyType:       string(key.KeyType),
		})
	}
	return table, nil
}

func (c *liveClient) ExecuteStatement(ctx context.Context, request StatementRequest) (StatementAnswer, error) {
	input := &awssdk.ExecuteStatementInput{Statement: aws.String(request.Statement)}
	if request.Limit > 0 {
		input.Limit = new(request.Limit)
	}
	if request.NextToken != "" {
		input.NextToken = aws.String(request.NextToken)
	}
	out, err := c.client.ExecuteStatement(ctx, input)
	if err != nil {
		return StatementAnswer{}, err
	}
	return StatementAnswer{Items: out.Items, NextToken: aws.ToString(out.NextToken)}, nil
}

func (c *liveClient) ExecuteTransaction(ctx context.Context, statement string) error {
	_, err := c.client.ExecuteTransaction(ctx, &awssdk.ExecuteTransactionInput{
		TransactStatements: []types.ParameterizedStatement{
			{Statement: aws.String(statement)},
		},
	})
	return err
}
