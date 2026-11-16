package driver

import (
	"context"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

// client is what the driver asks of a Cassandra cluster: one page of one
// statement's answer, one statement run for its effect, and nothing else. It
// is the seam a test replaces with a tape, so the driver's own reading of what
// gocql answers is what a recording and a replay exercise.
//
// Query returns the columns and rows the server described for one page,
// ordered as it sent them, and the state that continues after that page, empty
// on the last one. Every value is the pointer gocql's RowData hands back, so
// the row conversion below is the driver's in a replay too.
type client interface {
	Query(ctx context.Context, statement string, pageSize int, pageState []byte) (columns []gocql.ColumnInfo, rows []gocql.RowData, nextPageState []byte, err error)
	Exec(ctx context.Context, statement string) error
	Close()
}

// liveClient is the client the driver opens against a cluster.
type liveClient struct {
	session *gocql.Session
}

// Query reads one page. gocql stops an iterator at the page boundary, so a
// statement with more rows is one page per Query call, continued by the state
// this returns.
func (c *liveClient) Query(ctx context.Context, statement string, pageSize int, pageState []byte) ([]gocql.ColumnInfo, []gocql.RowData, []byte, error) {
	iter := c.session.Query(statement).PageSize(pageSize).PageState(pageState).IterContext(ctx)
	defer func() { _ = iter.Close() }()

	// The state is read before the rows: gocql sets it from the page's answer,
	// and fetching the page after this one replaces it.
	nextPageState := iter.PageState()

	var rows []gocql.RowData
	for {
		rowData, err := iter.RowData()
		if err != nil {
			return nil, nil, nil, err
		}
		if !iter.Scan(rowData.Values...) {
			break
		}
		rows = append(rows, rowData)
	}
	if err := iter.Close(); err != nil {
		return nil, nil, nil, err
	}
	return iter.Columns(), rows, nextPageState, nil
}

// Exec runs one statement for its effect.
func (c *liveClient) Exec(ctx context.Context, statement string) error {
	return c.session.Query(statement).ExecContext(ctx)
}

// Close closes the session.
func (c *liveClient) Close() {
	c.session.Close()
}
