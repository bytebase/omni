package driver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/bytebase/omni/cosmosdb/driver/queryengine/gonative"
)

// Client is what the driver asks of Cosmos DB: the account's and a database's
// listings, one container's query pages and a statement's plan, and a
// connection check. It is the seam a test replaces with a tape, so a recording
// exercises the driver's own live path and a replay answers it from the tape
// instead.
//
// It is exported because the recording harness in cosmosdb/internal/testcontainer
// asks the pinned emulator through the same path; nothing else calls it.
type Client interface {
	Ping(ctx context.Context) error
	// Close releases the connections the client opened.
	Close() error
	// ListDatabases returns the account's database names.
	ListDatabases(ctx context.Context) ([]string, error)
	// ListContainers returns database's container names.
	ListContainers(ctx context.Context, database string) ([]string, error)
	// Page reads one container query, page by page.
	Page(ctx context.Context, req PageRequest) (PageAnswer, error)
	// Plan returns the plan the gateway states for req.Statement, as the JSON
	// the client's query engine runs a cross-partition query from.
	Plan(ctx context.Context, req PageRequest) ([]byte, error)
}

// gatewayPageSize is the page the gateway serves a query that names no size,
// and the largest page this driver asks for: a query may ask for more, but no
// page past it is guaranteed. It is what the client's query engine sees each
// range deliver, and the number internal/testcontainer measured on the
// emulator and on a real Cosmos account.
const gatewayPageSize = 1000

// PageRequest is one query of one container, as the driver asks for it.
type PageRequest struct {
	Database  string
	Container string
	Statement string
	// PageSize is the items a page may hold, 0 for the gateway's own.
	PageSize int32
	// Engine asks for the query to run through the client's Go query engine,
	// which serves a cross-partition query the gateway refuses to merge.
	Engine bool
	// Limit stops reading once this many items are in hand; 0 reads all.
	Limit int
}

// PageAnswer is what reading a query returned: the items of every page it read
// and how many pages that took. The page count tells a query that failed
// outright from one that failed partway through.
type PageAnswer struct {
	Items [][]byte
	Pages int
}

// liveClient is the client the driver opens against an account: the azcosmos
// client, with the emulator's REST listing when the caller named the emulator.
type liveClient struct {
	client   *azcosmos.Client
	emulator *emulatorRESTClient
}

// Ping reads the database list as a connectivity check.
func (c *liveClient) Ping(ctx context.Context) error {
	if c.emulator != nil {
		// List databases as a connectivity check.
		if _, err := c.emulator.listDatabases(); err != nil {
			return fmt.Errorf("failed to ping CosmosDB emulator: %w", err)
		}
		return nil
	}
	queryPager := c.client.NewQueryDatabasesPager("select 1", nil)
	for queryPager.More() {
		if _, err := queryPager.NextPage(ctx); err != nil {
			return fmt.Errorf("failed to ping CosmosDB: %w", err)
		}
	}
	return nil
}

// ListDatabases returns the account's database names.
func (c *liveClient) ListDatabases(ctx context.Context) ([]string, error) {
	if c.emulator != nil {
		return c.emulator.listDatabases()
	}
	var names []string
	queryPager := c.client.NewQueryDatabasesPager("select * from dbs d", nil)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page in database list: %w", err)
		}
		for _, database := range queryResponse.Databases {
			names = append(names, database.ID)
		}
	}
	slices.Sort(names)
	return names, nil
}

// ListContainers returns database's container names.
func (c *liveClient) ListContainers(ctx context.Context, databaseName string) ([]string, error) {
	if c.emulator != nil {
		return c.emulator.listContainers(databaseName)
	}
	database, err := c.client.NewDatabase(databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to get database %q: %w", databaseName, err)
	}
	var names []string
	queryPager := database.NewQueryContainersPager("select * from colls c", nil)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page in container list: %w", err)
		}
		for _, container := range queryResponse.Containers {
			names = append(names, container.ID)
		}
	}
	slices.Sort(names)
	return names, nil
}

// containerOpenError is the container the driver could not open. It is its own
// type because the query reports this failure with its own message, where every
// other failure is one the pager returned.
type containerOpenError struct {
	err error
}

func (e *containerOpenError) Error() string { return "failed to create container: " + e.err.Error() }

func (e *containerOpenError) Unwrap() error { return e.err }

// Close releases the client's connections.
func (c *liveClient) Close() error {
	c.client.Close()
	if c.emulator != nil {
		c.emulator.close()
	}
	return nil
}

// Page reads one container query, page by page, until the pager is done or
// limit items are in hand.
func (c *liveClient) Page(ctx context.Context, req PageRequest) (PageAnswer, error) {
	container, err := c.client.NewContainer(req.Database, req.Container)
	if err != nil {
		return PageAnswer{}, &containerOpenError{err: err}
	}
	var options *azcosmos.QueryOptions
	if req.PageSize > 0 || req.Engine {
		options = &azcosmos.QueryOptions{PageSizeHint: req.PageSize}
		if req.Engine {
			options.QueryEngine = gonative.Default()
		}
	}
	// The pages read are part of the answer even when the query failed: they
	// tell a query that failed outright from one that failed partway through,
	// and only the first may run again through the engine.
	items, pages, err := drainQuery(ctx, container, req.Statement, options, req.Limit)
	return PageAnswer{Items: items, Pages: pages}, err
}

// drainQuery pages a query to exhaustion or to limit, reporting how many pages
// it read so the caller can tell a query that failed outright from one that
// failed partway through.
func drainQuery(ctx context.Context, container *azcosmos.ContainerClient, statement string, o *azcosmos.QueryOptions, limit int) ([][]byte, int, error) {
	pager := container.NewQueryItemsPager(statement, azcosmos.NewPartitionKey(), o)
	var items [][]byte
	pages := 0
	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return nil, pages, err
		}
		pages++
		for _, item := range response.Items {
			items = append(items, item)
			if limit > 0 && len(items) >= limit {
				return items, pages, nil
			}
		}
	}
	return items, pages, nil
}

// Plan asks the gateway for the plan of a statement. The client fetches that
// plan for every query it runs through its query engine and hands it to the
// engine before it reads a partition, so the plan is asked for the way a query
// the engine serves asks for it: an engine that keeps the plan and reads
// nothing plans the statement without running it.
func (c *liveClient) Plan(ctx context.Context, req PageRequest) ([]byte, error) {
	container, err := c.client.NewContainer(req.Database, req.Container)
	if err != nil {
		return nil, &containerOpenError{err: err}
	}
	capture := &planCapture{query: req.Statement}
	pager := container.NewQueryItemsPager(req.Statement, azcosmos.NewPartitionKey(), &azcosmos.QueryOptions{
		QueryEngine: capture,
	})
	if _, err := pager.NextPage(ctx); err != nil {
		return nil, err
	}
	if capture.plan == "" {
		return nil, errors.New("the gateway stated no query plan")
	}
	return []byte(capture.plan), nil
}

// isNotFound reports whether err is the account saying an object is not there.
func isNotFound(err error) bool {
	var responseErr *azcore.ResponseError
	if errors.As(err, &responseErr) {
		return responseErr.StatusCode == http.StatusNotFound
	}
	var restErr *emulatorHTTPError
	return errors.As(err, &restErr) && restErr.StatusCode == http.StatusNotFound
}

// isCrossPartitionUnsupported reports whether err is the gateway saying it
// cannot serve this query across partitions by itself.
func isCrossPartitionUnsupported(err error) bool {
	var responseErr *azcore.ResponseError
	if !errors.As(err, &responseErr) || responseErr.StatusCode != http.StatusBadRequest {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "cross partition query can not be directly served by the gateway") ||
		strings.Contains(msg, "Cross partition query only supports")
}
