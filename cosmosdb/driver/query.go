package driver

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytebase/omni/cosmosdb"
	"github.com/bytebase/omni/driver"
)

// Session opens one session. Cosmos DB has no session state this driver sets,
// so a session is the driver itself.
func (d *Driver) Session(_ context.Context) (driver.Session, error) {
	return &session{driver: d}, nil
}

type session struct {
	driver *Driver
}

func (*session) Close() error { return nil }

// Query runs script in the container QueryOptions.Schema names. The query
// language is per container, so a script is one statement: its text goes to
// the container as written and comes back as one Result.
//
// A statement the query language does not parse is that statement's error, as
// Bytebase's driver answered it; a call without a container is the call's
// error, because there is no query to run at all.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	if opts.Schema == "" {
		return nil, errors.New("container argument is required for CosmosDB")
	}

	// Allow `SELECT * FROM {container} [alias]` only.
	if _, err := cosmosdb.Parse(script); err != nil {
		return []*driver.Result{{
			Statement: script,
			Err:       fmt.Errorf("support simple SELECT statement for Cosmos DB engine only, err: %s", err.Error()),
		}}, nil
	}

	startTime := time.Now()
	result := &driver.Result{
		Statement: script,
		Columns:   []driver.Column{{Name: "result", TypeName: "TEXT"}},
	}

	req := PageRequest{
		Database:  s.driver.databaseName,
		Container: opts.Schema,
		Statement: script,
		Limit:     opts.Limit,
	}
	// TODO(zp): Rewrite limit while parser supported.
	if opts.Limit > 0 && opts.Limit < 1000 {
		// Set the page size to limit to avoid large page size.
		req.PageSize = int32(opts.Limit)
	}

	answer, err := queryEveryPartition(ctx, s.driver.client, req)
	switch {
	case err == nil:
		for _, item := range answer.Items {
			result.Rows = append(result.Rows, []driver.Value{
				{Kind: driver.String, String: string(item)},
			})
		}
	case isContainerOpenError(err):
		// The container the query names did not open, and its own failure is
		// the one to report.
		result.Err = err
	default:
		result.Err = fmt.Errorf("failed to read more items: %w", err)
	}

	result.Latency = time.Since(startTime)
	return []*driver.Result{result}, nil
}

// queryEveryPartition reads every page of the query, and, when the gateway
// refuses to serve a cross-partition query itself, runs it again through the
// client's query engine.
//
// The gateway serves what it can on its own; an aggregate, DISTINCT, ORDER BY,
// TOP, OFFSET/LIMIT or GROUP BY that spans partitions it refuses with a
// BadRequest instead, and the query then runs again through the pure-Go engine,
// which merges the per-partition streams itself.
//
// Only a query that has returned nothing yet can switch: once the gateway has
// handed back a continuation token, the engine cannot resume from it.
func queryEveryPartition(ctx context.Context, c Client, req PageRequest) (PageAnswer, error) {
	answer, err := c.Page(ctx, req)
	if err == nil || answer.Pages > 0 || !isCrossPartitionUnsupported(err) {
		return answer, err
	}
	req.Engine = true
	return c.Page(ctx, req)
}

// isContainerOpenError reports whether err is the container the driver could
// not open, whose message the query reports as its own.
func isContainerOpenError(err error) bool {
	var openErr *containerOpenError
	return errors.As(err, &openErr)
}
