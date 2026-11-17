package driver

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang-sql/sqlexp"
	gomssqldb "github.com/microsoft/go-mssqldb"

	"github.com/bytebase/omni/driver"
	omnimssql "github.com/bytebase/omni/mssql"
)

// queryFunc runs one batch as a query, which is how SQL Server reports a
// batch's per-statement counts: through the sqlexp message stream.
type queryFunc func(ctx context.Context, query string, args ...any) (*sql.Rows, error)

// readBatchMessages walks the message stream a batch answer carries. counts are
// the MsgRowsAffected values outside a result set, in statement order; a count
// that closes a result set is that set's row count and belongs to no statement.
// read consumes the result set a MsgNext announces, and failures holds every
// MsgError the server sent, in order. err is a transport or read failure, not a
// server error.
func readBatchMessages(ctx context.Context, rows *sql.Rows, retmsg *sqlexp.ReturnMessage, read func(*sql.Rows) (*driver.Result, error)) (counts []int64, sets []*driver.Result, failures []gomssqldb.Error, err error) {
	inResultSet := false
	failed := false
	for more := true; more; {
		switch m := retmsg.Message(ctx).(type) {
		case sqlexp.MsgNotice:
			// A PRINT or a low-severity RAISERROR is informational. The
			// driver's results do not carry it, as they did not before.
		case sqlexp.MsgError:
			var e gomssqldb.Error
			if !errors.As(m.Error, &e) {
				return counts, sets, failures, m.Error
			}
			failures = append(failures, e)
			failed = true
		case sqlexp.MsgRowsAffected:
			// A statement after a failure did not run as part of the script's
			// progress, so its count is not reported with the ones before it.
			if !inResultSet && !failed {
				counts = append(counts, m.Count)
			}
		case sqlexp.MsgNextResultSet:
			inResultSet = false
			more = rows.NextResultSet()
			if err := rows.Err(); err != nil {
				return counts, sets, failures, err
			}
		case sqlexp.MsgNext:
			result, err := read(rows)
			if err != nil {
				return counts, sets, failures, err
			}
			if result != nil {
				sets = append(sets, result)
			}
			inResultSet = true
		case sqlexp.MsgLastInsertID:
			// SQL Server does not send this, and the driver does not read it.
		}
	}
	if err := rows.Err(); err != nil {
		return counts, sets, failures, err
	}
	return counts, sets, failures, nil
}

// drainRows reads a result set to its end. Execute keeps no rows.
func drainRows(rows *sql.Rows) (*driver.Result, error) {
	for rows.Next() {
	}
	return nil, rows.Err()
}

// batchMessageError joins the server errors a batch reported into the error
// go-mssqldb hands a non-message query, whose text is every message in order.
func batchMessageError(errs []gomssqldb.Error) error {
	if len(errs) == 0 {
		return nil
	}
	last := errs[len(errs)-1]
	last.All = errs
	return last
}

// statementCount is the number of non-empty statements in a batch. ok is false
// when the parser cannot read the batch, and the caller then keeps the path
// that does not split it.
func statementCount(batch string) (n int, ok bool) {
	statements, err := omnimssql.Parse(batch)
	if err != nil {
		return 0, false
	}
	for _, statement := range statements {
		if !statement.Empty() {
			n++
		}
	}
	return n, true
}

// executeBatch runs a multi-statement batch through the message stream and
// returns what each statement affected, in order. The counts before a failing
// statement are still returned with its error.
func executeBatch(ctx context.Context, query queryFunc, command driver.Command) ([]int64, error) {
	retmsg := &sqlexp.ReturnMessage{}
	rows, err := query(ctx, command.Text, retmsg)
	if err != nil {
		return nil, wrapStatementError(command, err)
	}
	defer rows.Close()

	counts, _, failures, err := readBatchMessages(ctx, rows, retmsg, drainRows)
	if err != nil {
		return counts, wrapStatementError(command, err)
	}
	if err := batchMessageError(failures); err != nil {
		return counts, wrapStatementError(command, err)
	}
	return counts, nil
}
