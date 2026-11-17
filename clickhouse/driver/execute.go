package driver

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bytebase/omni/driver"
)

var _ driver.Executor = (*Driver)(nil)

// Execute runs every statement of script and returns the rows they affected in
// total, which this client cannot state: it answers 0, for the reason the
// package doc gives for Result.AffectedRows. It parses no transaction directive:
// the caller chooses the mode, and an empty mode is Bytebase's default, one
// transaction.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	statements, err := splitSQL(script)
	if err != nil {
		return 0, err
	}
	statements = filterEmptyStatements(statements)
	if len(statements) == 0 {
		return 0, nil
	}
	if opts.Transaction == driver.AutoCommit {
		return d.executeInAutoCommitMode(ctx, statements, opts)
	}
	return d.executeInTransactionMode(ctx, statements, opts)
}

func (d *Driver) executeInTransactionMode(ctx context.Context, statements []statement, opts driver.ExecuteOptions) (int64, error) {
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{Isolation: sqlIsolation(opts.Isolation)})
	if err != nil {
		return 0, err
	}

	committed := false
	defer func() {
		if committed {
			return
		}
		rollbackErr := tx.Rollback()
		var reported error
		if rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			reported = rollbackErr
		}
		transactionControl(opts.Log, driver.Rollback, reported)
	}()

	transactionControl(opts.Log, driver.Begin, nil)

	for _, single := range statements {
		commandExecute(opts.Log, driver.Command{
			Text:      single.Text,
			ByteStart: single.ByteStart,
			ByteEnd:   single.ByteEnd,
		})
		if _, err := tx.ExecContext(ctx, single.Text); err != nil {
			commandResponse(opts.Log, 0, nil, err)
			return 0, err
		}
		commandResponse(opts.Log, 0, nil, nil)
	}

	if err := tx.Commit(); err != nil {
		transactionControl(opts.Log, driver.Commit, err)
		return 0, err
	}

	transactionControl(opts.Log, driver.Commit, nil)
	committed = true
	return 0, nil
}

func (d *Driver) executeInAutoCommitMode(ctx context.Context, statements []statement, opts driver.ExecuteOptions) (int64, error) {
	for _, single := range statements {
		commandExecute(opts.Log, driver.Command{
			Text:      single.Text,
			ByteStart: single.ByteStart,
			ByteEnd:   single.ByteEnd,
		})
		if _, err := d.db.ExecContext(ctx, single.Text); err != nil {
			commandResponse(opts.Log, 0, nil, err)
			return 0, err
		}
		commandResponse(opts.Log, 0, nil, nil)
	}
	return 0, nil
}

// sqlIsolation maps the contract's level onto database/sql's. An empty level is
// the driver's default.
func sqlIsolation(level driver.IsolationLevel) sql.IsolationLevel {
	switch level {
	case driver.ReadUncommitted:
		return sql.LevelReadUncommitted
	case driver.ReadCommitted:
		return sql.LevelReadCommitted
	case driver.RepeatableRead:
		return sql.LevelRepeatableRead
	case driver.Serializable:
		return sql.LevelSerializable
	default:
		return sql.LevelDefault
	}
}

func commandExecute(log driver.ExecuteLogger, cmd driver.Command) {
	if log != nil {
		log.CommandExecute(cmd)
	}
}

func commandResponse(log driver.ExecuteLogger, affectedRows int64, perResult []int64, err error) {
	if log != nil {
		log.CommandResponse(affectedRows, perResult, err)
	}
}

func transactionControl(log driver.ExecuteLogger, op driver.TransactionOp, err error) {
	if log != nil {
		log.TransactionControl(op, err)
	}
}
