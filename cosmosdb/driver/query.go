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

// Query runs each statement of script in order in the container QueryOptions.Schema
// names and returns one Result per statement run. A statement that fails ends
// the script: its Result carries the error and nothing after it runs.
//
// The script is split at its top-level semicolons by cosmosdb.Split, so a
// statement's boundaries are the query language's punctuation. Each statement
// goes to the container on its own: the query language has no statement
// terminator, and the gateway rejects text holding two.
//
// A statement the query language does not parse is that statement's error, as
// Bytebase's driver answered it; a call without a container is the call's
// error, because there is no query to run at all.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	if opts.Schema == "" {
		return nil, errors.New("container argument is required for CosmosDB")
	}

	var statements []cosmosdb.Segment
	for _, segment := range cosmosdb.Split(script) {
		if segment.Empty() {
			continue
		}
		statements = append(statements, segment)
	}
	if len(statements) == 0 {
		return nil, nil
	}

	var results []*driver.Result
	for _, statement := range statements {
		result := s.queryStatement(ctx, statement.Text, opts)
		results = append(results, result)
		if result.Err != nil {
			break
		}
	}
	return results, nil
}

// queryStatement runs one statement of a script, or reads its plan when the
// caller asked for that.
func (s *session) queryStatement(ctx context.Context, statement string, opts driver.QueryOptions) *driver.Result {
	startTime := time.Now()
	result := &driver.Result{
		Statement: statement,
		Columns:   []driver.Column{{Name: "result", TypeName: "TEXT"}},
	}
	defer func() { result.Latency = time.Since(startTime) }()

	// Allow `SELECT * FROM {container} [alias]` only.
	if _, err := cosmosdb.Parse(statement); err != nil {
		result.Err = fmt.Errorf("support simple SELECT statement for Cosmos DB engine only, err: %s", err.Error())
		return result
	}

	req := PageRequest{
		Database:  s.driver.databaseName,
		Container: opts.Schema,
		Statement: statement,
	}
	if opts.Explain {
		s.explain(ctx, result, req)
		return result
	}

	req.Limit = opts.Limit
	// TODO(zp): Rewrite limit while parser supported.
	// The page size is the caller's limit, up to the page the gateway serves
	// of its own accord: a smaller limit reads one page of its rows, and a
	// larger one asks for that page rather than leaving the size to the
	// gateway.
	if opts.Limit > 0 {
		req.PageSize = int32(min(opts.Limit, gatewayPageSize))
	}

	answer, err := queryEveryPartition(ctx, s.driver.client, req)
	switch {
	case err == nil:
		totalBytes := driver.HeaderByteSize(result.Columns)
		for _, item := range answer.Items {
			row := []driver.Value{{Kind: driver.String, String: string(item)}}
			result.Rows = append(result.Rows, row)
			if opts.MaxResultBytes > 0 {
				totalBytes += driver.RowByteSize(row)
				if totalBytes > opts.MaxResultBytes {
					// The client's pager is below this driver: the cap is
					// applied to the rows it handed over, not to the read.
					result.Truncated = true
					break
				}
			}
		}
	case isContainerOpenError(err):
		// The container the query names did not open, and its own failure is
		// the one to report.
		result.Err = err
	default:
		result.Err = fmt.Errorf("failed to read more items: %w", err)
	}
	return result
}

// explain fills result with the plan the gateway states for req.Statement,
// which is the plan the client's query engine runs a cross-partition query
// from. The statement does not run: the plan is one row of the gateway's JSON.
func (s *session) explain(ctx context.Context, result *driver.Result, req PageRequest) {
	plan, err := s.driver.client.Plan(ctx, req)
	switch {
	case err != nil && isContainerOpenError(err):
		result.Err = err
	case err != nil:
		result.Err = fmt.Errorf("failed to read the query plan: %w", err)
	default:
		result.Rows = [][]driver.Value{{{Kind: driver.String, String: string(plan)}}}
		result.Plan = &driver.Plan{Format: "JSON"}
	}
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
