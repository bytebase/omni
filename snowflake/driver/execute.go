package driver

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bytebase/omni/driver"
	snowflakeanalysis "github.com/bytebase/omni/snowflake/analysis"
)

// Execute runs script under opts.Transaction, reporting each step to opts.Log
// when it is non-nil and returning the total affected rows.
//
// The script is split lexically, so every statement is sent whole whatever its
// text looks like; bounding how many statements a sheet may hold is the
// caller's. A failing statement ends the script: the transaction is rolled
// back, or, in autocommit mode, the statements before it stay committed.
//
// Bytebase sent the whole script as one request with MULTI_STATEMENT_COUNT=0,
// which reports one answer and one error for all of it; a recorded run confirms
// that sending the statements one at a time executes the same script.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	var commands []driver.Command
	for _, statement := range snowflakeanalysis.SplitStatements(script) {
		if statement.Empty {
			continue
		}
		commands = append(commands, driver.Command{
			Text:      statement.Text,
			ByteStart: statement.Start,
			ByteEnd:   statement.End,
		})
	}
	if len(commands) == 0 {
		return 0, nil
	}

	if opts.Transaction == driver.AutoCommit {
		return d.executeInAutoCommitMode(ctx, commands, opts.Log)
	}
	return d.executeInTransactionMode(ctx, commands, opts.Log)
}

// executeInTransactionMode runs commands in one transaction and commits it.
func (d *Driver) executeInTransactionMode(ctx context.Context, commands []driver.Command, log driver.ExecuteLogger) (int64, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		logTransactionControl(log, driver.Begin, err)
		return 0, err
	}
	defer discardConn(conn)
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		logTransactionControl(log, driver.Begin, err)
		return 0, err
	}
	logTransactionControl(log, driver.Begin, nil)

	committed := false
	defer func() {
		err := tx.Rollback()
		if committed {
			return
		}
		var rollbackErr error
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			rollbackErr = err
		}
		logTransactionControl(log, driver.Rollback, rollbackErr)
	}()

	var totalRowsAffected int64
	for _, command := range commands {
		logCommandExecute(log, command)
		sqlResult, err := tx.ExecContext(statementContext(ctx), command.Text)
		if err != nil {
			logCommandResponse(log, 0, nil, err)
			return 0, err
		}
		rowsAffected := rowsAffected(sqlResult)
		logCommandResponse(log, rowsAffected, []int64{rowsAffected}, nil)
		totalRowsAffected += rowsAffected
	}

	if err := tx.Commit(); err != nil {
		logTransactionControl(log, driver.Commit, err)
		return 0, err
	}
	logTransactionControl(log, driver.Commit, nil)
	committed = true
	return totalRowsAffected, nil
}

// executeInAutoCommitMode runs commands one at a time, each committing itself,
// and stops at the first failure, which leaves the database partially migrated.
// Every command runs on one connection, so a script that sets session state (a
// schema, a warehouse, a temporary table) keeps it for the statements after it.
func (d *Driver) executeInAutoCommitMode(ctx context.Context, commands []driver.Command, log driver.ExecuteLogger) (int64, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer discardConn(conn)

	var totalRowsAffected int64
	for _, command := range commands {
		logCommandExecute(log, command)
		sqlResult, err := conn.ExecContext(statementContext(ctx), command.Text)
		if err != nil {
			logCommandResponse(log, 0, nil, err)
			return totalRowsAffected, err
		}
		rowsAffected := rowsAffected(sqlResult)
		logCommandResponse(log, rowsAffected, []int64{rowsAffected}, nil)
		totalRowsAffected += rowsAffected
	}
	return totalRowsAffected, nil
}

// rowsAffected reads a statement's affected rows. The count is only meaningful
// for DML, and the driver cannot tell a DDL statement from one here, so a
// statement that reports none counts as zero.
func rowsAffected(result sql.Result) int64 {
	affected, err := result.RowsAffected()
	if err != nil {
		return 0
	}
	return affected
}

func logCommandExecute(log driver.ExecuteLogger, cmd driver.Command) {
	if log != nil {
		log.CommandExecute(cmd)
	}
}

func logCommandResponse(log driver.ExecuteLogger, affected int64, perResult []int64, err error) {
	if log != nil {
		log.CommandResponse(affected, perResult, err)
	}
}

func logTransactionControl(log driver.ExecuteLogger, op driver.TransactionOp, err error) {
	if log != nil {
		log.TransactionControl(op, err)
	}
}
