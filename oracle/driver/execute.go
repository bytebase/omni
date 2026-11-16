package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	omnidriver "github.com/bytebase/omni/driver"
	oracleparser "github.com/bytebase/omni/oracle/parser"
)

// Execute runs each statement of script under opts.Transaction and returns the
// rows they affected in total, reporting each step to opts.Log as it goes.
// Oracle has no CREATE DATABASE, so opts.CreateDatabase is refused.
func (d *Driver) Execute(ctx context.Context, script string, opts omnidriver.ExecuteOptions) (int64, error) {
	if opts.CreateDatabase {
		return 0, errors.New("create database is not supported for Oracle")
	}
	isolation, err := toSQLIsolation(opts.Isolation)
	if err != nil {
		return 0, err
	}
	// Splitting alone is enough: a statement the parser does not know still runs,
	// as it does in SQL*Plus.
	var commands []omnidriver.Command
	for _, statement := range oracleparser.Split(script) {
		if statement.Empty() {
			continue
		}
		commands = append(commands, omnidriver.Command{
			Text:      statement.Text,
			ByteStart: statement.ByteStart,
			ByteEnd:   statement.ByteEnd,
		})
	}
	if len(commands) == 0 {
		return 0, nil
	}

	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}
	defer func() { _ = conn.Close() }()

	if opts.Transaction == omnidriver.AutoCommit {
		return executeInAutoCommitMode(ctx, conn, commands, opts.Log)
	}
	return executeInTransactionMode(ctx, conn, commands, opts.Log, isolation)
}

// executeInTransactionMode runs the script in one transaction.
func executeInTransactionMode(ctx context.Context, conn *sql.Conn, commands []omnidriver.Command, log omnidriver.ExecuteLogger, isolation sql.IsolationLevel) (int64, error) {
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: isolation})
	if err != nil {
		logTransactionControl(log, omnidriver.Begin, err)
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	logTransactionControl(log, omnidriver.Begin, nil)

	committed := false
	defer func() {
		if committed {
			return
		}
		rollbackErr := tx.Rollback()
		logTransactionControl(log, omnidriver.Rollback, rollbackErr)
	}()

	var totalAffected int64
	for _, command := range commands {
		logCommandExecute(log, command)
		sqlResult, err := tx.ExecContext(ctx, command.Text)
		if err != nil {
			logCommandResponse(log, 0, nil, err)
			return 0, err
		}
		affected := affectedRows(sqlResult)
		logCommandResponse(log, affected, []int64{affected}, nil)
		totalAffected += affected
	}

	if err := tx.Commit(); err != nil {
		logTransactionControl(log, omnidriver.Commit, err)
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	logTransactionControl(log, omnidriver.Commit, nil)
	committed = true
	return totalAffected, nil
}

// executeInAutoCommitMode runs each statement on its own, stopping at the first
// failure.
func executeInAutoCommitMode(ctx context.Context, conn *sql.Conn, commands []omnidriver.Command, log omnidriver.ExecuteLogger) (int64, error) {
	var totalAffected int64
	for _, command := range commands {
		logCommandExecute(log, command)
		sqlResult, err := conn.ExecContext(ctx, command.Text)
		if err != nil {
			logCommandResponse(log, 0, nil, err)
			return totalAffected, err
		}
		affected := affectedRows(sqlResult)
		logCommandResponse(log, affected, []int64{affected}, nil)
		totalAffected += affected
	}
	return totalAffected, nil
}

// affectedRows reads a statement's affected rows, which a DDL statement does not
// report.
func affectedRows(result sql.Result) int64 {
	affected, err := result.RowsAffected()
	if err != nil {
		return 0
	}
	return affected
}

// toSQLIsolation maps the contract's isolation level onto database/sql's.
func toSQLIsolation(level omnidriver.IsolationLevel) (sql.IsolationLevel, error) {
	switch level {
	case "":
		return sql.LevelDefault, nil
	case omnidriver.ReadUncommitted:
		return sql.LevelReadUncommitted, nil
	case omnidriver.ReadCommitted:
		return sql.LevelReadCommitted, nil
	case omnidriver.RepeatableRead:
		return sql.LevelRepeatableRead, nil
	case omnidriver.Serializable:
		return sql.LevelSerializable, nil
	default:
		return sql.LevelDefault, fmt.Errorf("invalid isolation level for Oracle: %s", level)
	}
}

func logCommandExecute(log omnidriver.ExecuteLogger, command omnidriver.Command) {
	if log != nil {
		log.CommandExecute(command)
	}
}

func logCommandResponse(log omnidriver.ExecuteLogger, affected int64, perResult []int64, err error) {
	if log != nil {
		log.CommandResponse(affected, perResult, err)
	}
}

func logTransactionControl(log omnidriver.ExecuteLogger, op omnidriver.TransactionOp, err error) {
	if log != nil {
		log.TransactionControl(op, err)
	}
}
