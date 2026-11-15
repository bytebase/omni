package driver

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"fmt"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/driver"
	mariadbparser "github.com/bytebase/omni/mariadb/parser"
)

const (
	maxSheetCheckBytes = 2 * 1024 * 1024
	maxSplitCommands   = 200
)

func splitCommands(script string) []driver.Command {
	segments := mariadbparser.Split(script)
	commands := make([]driver.Command, 0, len(segments))
	for _, seg := range segments {
		if seg.Empty() {
			continue
		}
		end := seg.ByteEnd
		if end < len(script) && script[end] == ';' {
			end++
		}
		commands = append(commands, driver.Command{
			Text:      script[seg.ByteStart:end],
			ByteStart: seg.ByteStart,
			ByteEnd:   end,
		})
	}
	return commands
}

func buildExecuteCommands(script string) []driver.Command {
	if len(script) > maxSheetCheckBytes {
		return []driver.Command{{Text: script, ByteStart: 0, ByteEnd: len(script)}}
	}
	commands := splitCommands(script)
	if len(commands) == 0 {
		return nil
	}
	if len(commands) > maxSplitCommands {
		return []driver.Command{{Text: script, ByteStart: 0, ByteEnd: len(script)}}
	}
	return commands
}

func toSQLIsolation(level driver.IsolationLevel) (sqldriver.IsolationLevel, error) {
	switch level {
	case "":
		return sqldriver.IsolationLevel(sql.LevelDefault), nil
	case driver.ReadUncommitted:
		return sqldriver.IsolationLevel(sql.LevelReadUncommitted), nil
	case driver.ReadCommitted:
		return sqldriver.IsolationLevel(sql.LevelReadCommitted), nil
	case driver.RepeatableRead:
		return sqldriver.IsolationLevel(sql.LevelRepeatableRead), nil
	case driver.Serializable:
		return sqldriver.IsolationLevel(sql.LevelSerializable), nil
	default:
		return 0, fmt.Errorf("invalid isolation level for MariaDB: %s. Supported levels: READ UNCOMMITTED, READ COMMITTED, REPEATABLE READ, SERIALIZABLE", level)
	}
}

// Execute runs script under opts.Transaction, reporting each step to opts.Log
// when non-nil and returning the total affected rows.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	isolation, err := toSQLIsolation(opts.Isolation)
	if err != nil {
		return 0, err
	}

	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	connID, err := d.readConnectionID(ctx, conn)
	if err != nil {
		return 0, err
	}

	commands := buildExecuteCommands(script)
	if len(commands) == 0 {
		return 0, nil
	}

	if opts.Transaction == driver.AutoCommit || opts.CreateDatabase {
		return d.executeInAutoCommitMode(ctx, conn, commands, opts.Log, connID)
	}
	return d.executeInTransactionMode(ctx, conn, commands, opts.Log, connID, isolation)
}

func (d *Driver) executeInTransactionMode(
	ctx context.Context,
	conn *sql.Conn,
	commands []driver.Command,
	log driver.ExecuteLogger,
	connID uint64,
	isolation sqldriver.IsolationLevel,
) (int64, error) {
	var totalRowsAffected int64
	if err := conn.Raw(func(driverConn any) error {
		exer := driverConn.(sqldriver.ExecerContext)
		txer := driverConn.(sqldriver.ConnBeginTx)

		tx, err := txer.BeginTx(ctx, sqldriver.TxOptions{Isolation: isolation})
		if err != nil {
			logTransactionControl(log, driver.Begin, err)
			return err
		}
		logTransactionControl(log, driver.Begin, nil)

		committed := false
		defer func() {
			if committed {
				return
			}
			rollbackErr := tx.Rollback()
			logTransactionControl(log, driver.Rollback, rollbackErr)
		}()

		for _, cmd := range commands {
			logCommandExecute(log, cmd)
			sqlResult, err := exer.ExecContext(ctx, d.withAppComment(cmd.Text), nil)
			if err != nil {
				d.killQueryOnCancel(err, connID)
				logCommandResponse(log, 0, nil, err)
				return err
			}
			rowsAffected, perResult := extractAffectedRows(sqlResult)
			totalRowsAffected += rowsAffected
			logCommandResponse(log, rowsAffected, perResult, nil)
		}

		if err := tx.Commit(); err != nil {
			logTransactionControl(log, driver.Commit, err)
			return fmt.Errorf("failed to commit execute transaction: %w", err)
		}
		logTransactionControl(log, driver.Commit, nil)
		committed = true
		return nil
	}); err != nil {
		return 0, err
	}
	return totalRowsAffected, nil
}

func (d *Driver) executeInAutoCommitMode(
	ctx context.Context,
	conn *sql.Conn,
	commands []driver.Command,
	log driver.ExecuteLogger,
	connID uint64,
) (int64, error) {
	var totalRowsAffected int64
	if err := conn.Raw(func(driverConn any) error {
		exer := driverConn.(sqldriver.ExecerContext)
		for _, cmd := range commands {
			logCommandExecute(log, cmd)
			sqlResult, err := exer.ExecContext(ctx, d.withAppComment(cmd.Text), nil)
			if err != nil {
				d.killQueryOnCancel(err, connID)
				logCommandResponse(log, 0, nil, err)
				return err
			}
			rowsAffected, perResult := extractAffectedRows(sqlResult)
			totalRowsAffected += rowsAffected
			logCommandResponse(log, rowsAffected, perResult, nil)
		}
		return nil
	}); err != nil {
		return 0, err
	}
	return totalRowsAffected, nil
}

func extractAffectedRows(res sqldriver.Result) (int64, []int64) {
	if mr, ok := res.(gomysql.Result); ok {
		all := mr.AllRowsAffected()
		var total int64
		perResult := make([]int64, len(all))
		for i, a := range all {
			total += a
			perResult[i] = a
		}
		return total, perResult
	}
	ra, err := res.RowsAffected()
	if err != nil {
		ra = 0
	}
	return ra, []int64{ra}
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
