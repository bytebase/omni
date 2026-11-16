package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/doris/parser"
	"github.com/bytebase/omni/driver"
)

// Execute runs script under opts.Transaction and returns the rows it affected.
// The script goes to the server as one request, which multiStatements allows,
// so Doris decides where its statements begin and end.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	// Only preprocess when DELIMITER directives are present; normal Doris SQL
	// passes through unchanged.
	if containsDelimiterDirective(script) {
		segments := parser.Split(script)
		var buf strings.Builder
		for i, seg := range segments {
			if i > 0 {
				buf.WriteString(";\n")
			}
			buf.WriteString(seg.Text)
		}
		script = buf.String()
	}

	// Doris is an OLAP database with limited transaction support. For DDL
	// operations, transactions are not supported. For DML operations, Doris
	// supports transactions within certain limitations.
	if opts.Transaction == driver.AutoCommit || opts.CreateDatabase {
		return d.executeInAutoCommitMode(ctx, script, opts.Log)
	}
	return d.executeInTransactionMode(ctx, script, opts.Log)
}

func (d *Driver) executeInTransactionMode(ctx context.Context, script string, log driver.ExecuteLogger) (int64, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = conn.Close() }()

	connID, err := d.readConnectionID(ctx, conn)
	if err != nil {
		return 0, err
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		logTransactionControl(log, driver.Begin, err)
		return 0, fmt.Errorf("failed to begin execute transaction: %w", err)
	}
	logTransactionControl(log, driver.Begin, nil)

	committed := false
	defer func() {
		if committed {
			return
		}
		rollbackErr := tx.Rollback()
		if rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			logTransactionControl(log, driver.Rollback, rollbackErr)
			return
		}
		logTransactionControl(log, driver.Rollback, nil)
	}()

	sqlResult, err := tx.ExecContext(ctx, script)
	if err != nil {
		d.killQueryOnCancel(err, connID)
		return 0, err
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		// Since we cannot differentiate DDL and DML yet, we have to ignore the error.
		rowsAffected = 0
	}

	if err := tx.Commit(); err != nil {
		logTransactionControl(log, driver.Commit, err)
		return 0, fmt.Errorf("failed to commit execute transaction: %w", err)
	}
	logTransactionControl(log, driver.Commit, nil)
	committed = true
	return rowsAffected, nil
}

func (d *Driver) executeInAutoCommitMode(ctx context.Context, script string, log driver.ExecuteLogger) (int64, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = conn.Close() }()

	connID, err := d.readConnectionID(ctx, conn)
	if err != nil {
		return 0, err
	}

	sqlResult, err := conn.ExecContext(ctx, script)
	if err != nil {
		d.killQueryOnCancel(err, connID)
		return 0, err
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		// Since we cannot differentiate DDL and DML yet, we have to ignore the error.
		rowsAffected = 0
	}
	return rowsAffected, nil
}

func containsDelimiterDirective(s string) bool {
	for _, line := range strings.Split(s, "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		upper := strings.ToUpper(trimmed)
		if strings.HasPrefix(upper, "DELIMITER ") || strings.HasPrefix(upper, "DELIMITER\t") {
			return true
		}
	}
	return false
}

func logTransactionControl(log driver.ExecuteLogger, op driver.TransactionOp, err error) {
	if log != nil {
		log.TransactionControl(op, err)
	}
}
