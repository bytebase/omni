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
// listings, one container's query pages, and a connection check. It is the seam a
// test replaces with a tape, so a recording exercises the driver's own live path
// and a replay answers it from the tape instead.
//
// It is exported because the recording harness in cosmosdb/internal/testcontainer
// asks the pinned emulator through the same path; nothing else calls it.
type Client interface {
	Ping(ctx context.Context) error
	// ListDatabases returns the account's database names.
	ListDatabases(ctx context.Context) ([]string, error)
	// ListContainers returns database's container names.
	ListContainers(ctx context.Context, database string) ([]string, error)
	// Page reads one container query, page by page.
	Page(ctx context.Context, req PageRequest) (PageAnswer, error)
}

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
	items, pages, err := drainQuery(ctx, container, req.Statement, options, req.Limit)
	if err != nil {
		return PageAnswer{}, err
	}
	return PageAnswer{Items: items, Pages: pages}, nil
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
