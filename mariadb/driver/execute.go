package driver

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/mariadb/ast"
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

// delimiterDirective matches the client-side directive that changes the
// statement terminator. The server never sees one, so a script holding a
// directive runs only as the commands the splitter normalized: one command of
// the raw text would send "DELIMITER //" to MariaDB, which rejects it.
var delimiterDirective = regexp.MustCompile(`(?i)(^|[\s;])DELIMITER\s`)

// buildExecuteCommands splits script into the commands Execute sends. A script
// too large or holding too many statements is sent whole, so that a huge sheet
// does not turn into hundreds of logged steps or a transaction of them; a
// script that sets a DELIMITER is split either way.
func buildExecuteCommands(script string) []driver.Command {
	commands := splitCommands(script)
	if delimiterDirective.MatchString(script) {
		return commands
	}
	if len(script) > maxSheetCheckBytes || len(commands) > maxSplitCommands {
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
// when non-nil and returning the total affected rows. A script that fails on a
// lock wait timeout runs again up to opts.MaxRetries times, each attempt from
// the beginning, which is safe because such a failure commits nothing.
//
// MariaDB commits implicitly on every DDL statement, so a script that mixes DDL
// with DML leaves what it committed even when a later statement fails under
// Transactional; the failure is returned and the count is 0.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	isolation, err := toSQLIsolation(opts.Isolation)
	if err != nil {
		return 0, err
	}

	commands := buildExecuteCommands(script)
	if len(commands) == 0 {
		return 0, nil
	}

	// A retry repeats the whole script, so it is only safe when the failed
	// attempt committed nothing: every command that ran before the failure must
	// be undone by rolling the transaction back.
	var rollbackSafe []bool
	if opts.MaxRetries > 0 {
		rollbackSafe = make([]bool, len(commands))
		for i, cmd := range commands {
			parsed, err := mariadbparser.Parse(cmd.Text)
			if err != nil || parsed == nil || len(parsed.Items) != 1 {
				continue
			}
			rollbackSafe[i] = rollsBack(parsed.Items[0])
		}
	}

	run := func() (int64, bool, error) {
		conn, err := d.db.Conn(ctx)
		if err != nil {
			return 0, false, err
		}
		// A script may select a database or set session variables, so the
		// connection is not poolable afterwards either: see discardConn.
		defer discardConn(conn)

		connID, err := d.readConnectionID(ctx, conn)
		if err != nil {
			return 0, false, err
		}
		if opts.Transaction == driver.AutoCommit || opts.CreateDatabase {
			return d.executeInAutoCommitMode(ctx, conn, commands, opts.Log, connID)
		}
		return d.executeInTransactionMode(ctx, conn, commands, rollbackSafe, opts.Log, connID, isolation)
	}

	affectedRows, retryable, err := run()
	if err == nil || !retryable || !isLockTimeout(err) {
		return affectedRows, err
	}
	for i := range opts.MaxRetries {
		interval := time.Duration(150+(i+1)*100/opts.MaxRetries) * time.Millisecond
		select {
		case <-ctx.Done():
			return affectedRows, ctx.Err()
		case <-time.After(interval):
		}
		if opts.Log != nil {
			opts.Log.Retry(err, i+1, opts.MaxRetries)
		}
		affectedRows, retryable, err = run()
		if err == nil || !retryable || !isLockTimeout(err) {
			break
		}
	}
	return affectedRows, err
}

// rollsBack reports whether a command that succeeded inside a transaction is
// undone by rolling back. MariaDB commits the open transaction implicitly around
// DDL, LOCK TABLES, a SET that changes autocommit, and the administrative
// statements, so anything outside this list keeps a retry from being safe.
func rollsBack(node ast.Node) bool {
	switch node.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt,
		*ast.DoStmt, *ast.ShowStmt, *ast.ExplainStmt:
		return true
	}
	return false
}

// isLockTimeout reports whether err is MariaDB's lock wait timeout (1205) or
// deadlock (1213): the two failures that roll a statement's transaction back,
// so running the script again starts from a state that holds nothing.
func isLockTimeout(err error) bool {
	var mysqlErr *gomysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return false
	}
	return mysqlErr.Number == erLockWaitTimeout || mysqlErr.Number == erLockDeadlock
}

const (
	erLockWaitTimeout = 1205
	erLockDeadlock    = 1213
)

func (d *Driver) executeInTransactionMode(
	ctx context.Context,
	conn *sql.Conn,
	commands []driver.Command,
	rollbackSafe []bool,
	log driver.ExecuteLogger,
	connID uint64,
	isolation sqldriver.IsolationLevel,
) (rowsAffected int64, retryable bool, err error) {
	var totalRowsAffected int64
	// Every command that ran before a failure is undone by the rollback, unless
	// it committed on its own.
	earlierSafe := true
	retryable = true
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
			if rollbackErr != nil {
				retryable = false
			}
		}()

		for i, cmd := range commands {
			logCommandExecute(log, cmd)
			sqlResult, execErr := exer.ExecContext(ctx, d.withAppComment(cmd.Text), nil)
			if execErr != nil {
				d.killQueryOnCancel(execErr, connID)
				logCommandResponse(log, 0, nil, execErr)
				retryable = retryable && earlierSafe
				return execErr
			}
			if rollbackSafe != nil && !rollbackSafe[i] {
				earlierSafe = false
			}
			affected, perResult := extractAffectedRows(sqlResult)
			totalRowsAffected += affected
			logCommandResponse(log, affected, perResult, nil)
		}

		if commitErr := tx.Commit(); commitErr != nil {
			logTransactionControl(log, driver.Commit, commitErr)
			retryable = false
			return fmt.Errorf("failed to commit execute transaction: %w", commitErr)
		}
		logTransactionControl(log, driver.Commit, nil)
		committed = true
		return nil
	}); err != nil {
		return 0, retryable, err
	}
	return totalRowsAffected, retryable, nil
}

func (d *Driver) executeInAutoCommitMode(
	ctx context.Context,
	conn *sql.Conn,
	commands []driver.Command,
	log driver.ExecuteLogger,
	connID uint64,
) (rowsAffected int64, retryable bool, err error) {
	var totalRowsAffected int64
	// Each command commits on its own, so only a script that failed before its
	// first one ran leaves nothing to repeat.
	succeeded := 0
	if err := conn.Raw(func(driverConn any) error {
		exer := driverConn.(sqldriver.ExecerContext)
		for _, cmd := range commands {
			logCommandExecute(log, cmd)
			sqlResult, execErr := exer.ExecContext(ctx, d.withAppComment(cmd.Text), nil)
			if execErr != nil {
				d.killQueryOnCancel(execErr, connID)
				logCommandResponse(log, 0, nil, execErr)
				retryable = succeeded == 0
				return execErr
			}
			succeeded++
			affected, perResult := extractAffectedRows(sqlResult)
			totalRowsAffected += affected
			logCommandResponse(log, affected, perResult, nil)
		}
		return nil
	}); err != nil {
		return 0, retryable, err
	}
	return totalRowsAffected, retryable, nil
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
