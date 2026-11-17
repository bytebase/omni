package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"

	gomssqldb "github.com/microsoft/go-mssqldb"

	"github.com/bytebase/omni/driver"
	omnimssql "github.com/bytebase/omni/mssql"
	mssqlanalysis "github.com/bytebase/omni/mssql/analysis"
)

// Execute runs the script's batches under opts.Transaction, reporting each step
// to opts.Log when non-nil and returning the total affected rows. A batch
// followed by GO <count> runs that many times.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	if opts.CreateDatabase {
		if _, err := d.db.ExecContext(ctx, script); err != nil {
			return 0, err
		}
		return 0, nil
	}

	isolation, err := toSQLIsolation(opts.Isolation)
	if err != nil {
		return 0, err
	}
	if opts.Transaction == driver.AutoCommit {
		return d.executeInAutoCommitMode(ctx, script, opts.Log)
	}
	return d.executeInTransactionMode(ctx, script, opts.Log, isolation)
}

// executeInTransactionMode runs the script's batches in one transaction.
func (d *Driver) executeInTransactionMode(ctx context.Context, script string, log driver.ExecuteLogger, isolation sql.IsolationLevel) (int64, error) {
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{Isolation: isolation})
	if err != nil {
		logTransactionControl(log, driver.Begin, err)
		return 0, err
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

	var totalAffected int64
	if err := forEachBatch(script, func(command driver.Command) error {
		logCommandExecute(log, command)
		affected, err := executeStatement(ctx, tx.ExecContext, command)
		if err != nil {
			logCommandResponse(log, 0, nil, err)
			return err
		}
		logCommandResponse(log, affected, []int64{affected}, nil)
		totalAffected += affected
		return nil
	}); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		logTransactionControl(log, driver.Commit, err)
		return 0, err
	}
	logTransactionControl(log, driver.Commit, nil)
	committed = true
	return totalAffected, nil
}

// executeInAutoCommitMode runs each batch on its own, stopping at the first
// failure. The script keeps one connection, so what a batch sets up for the next
// (USE, a temporary table) is still there.
func (d *Driver) executeInAutoCommitMode(ctx context.Context, script string, log driver.ExecuteLogger) (int64, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var totalAffected int64
	if err := forEachBatch(script, func(command driver.Command) error {
		logCommandExecute(log, command)
		affected, err := executeStatement(ctx, conn.ExecContext, command)
		if err != nil {
			logCommandResponse(log, 0, nil, err)
			return err
		}
		logCommandResponse(log, affected, []int64{affected}, nil)
		totalAffected += affected
		return nil
	}); err != nil {
		return totalAffected, err
	}
	return totalAffected, nil
}

// forEachBatch calls run for every batch of the script, once per repetition a
// GO <count> asks for, with the batch's byte range in the script.
func forEachBatch(script string, run func(driver.Command) error) error {
	batch := mssqlanalysis.NewBatcher(script)
	for {
		command, err := batch.Next()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("failed to get next batch for statement: %s: %w", batch.Batch().Text, err)
			}
			if v := batch.Batch(); v != nil && v.Text != "" {
				return run(driver.Command{Text: v.Text, ByteStart: v.Start, ByteEnd: v.End})
			}
			return nil
		}
		if command == nil {
			continue
		}
		goCommand, ok := command.(*mssqlanalysis.GoCommand)
		if !ok {
			return fmt.Errorf("unsupported command type: %T", command)
		}
		v := batch.Batch()
		for i := uint(0); i < goCommand.Count; i++ {
			if err := run(driver.Command{Text: v.Text, ByteStart: v.Start, ByteEnd: v.End}); err != nil {
				return err
			}
		}
		batch.Reset(nil)
	}
}

// execFunc runs one statement.
type execFunc func(ctx context.Context, query string, args ...any) (sql.Result, error)

// executeStatement runs one batch and returns the rows it affected. A DDL
// statement reports no count, which the client returns as an error; that is 0.
// A failure names the statement the server reported it at, since a batch holds
// several and the server names only a line.
func executeStatement(ctx context.Context, exec execFunc, command driver.Command) (int64, error) {
	result, err := exec(ctx, command.Text)
	if err != nil {
		if statement, ok := failingStatement(command.Text, err); ok {
			return 0, fmt.Errorf("failed to execute statement %q: %w", statement, unpackMSSQLError(err))
		}
		return 0, unpackMSSQLError(fmt.Errorf("failed to execute statement: %w", err))
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return affected, nil
}

// failingStatement returns the text of the statement of batch the server
// reported err at, which is the statement that starts at or before the error's
// line. ok is false when the server named no line, or when the batch holds no
// statement that starts at or before it: the error then names the batch alone.
func failingStatement(batch string, err error) (string, bool) {
	line := errorLine(err)
	if line == 0 {
		return "", false
	}
	statements, parseErr := omnimssql.Parse(batch)
	if parseErr != nil {
		return "", false
	}
	failing := ""
	for _, statement := range statements {
		if statement.Empty() {
			continue
		}
		if statement.Start.Line > line {
			break
		}
		failing = strings.TrimSpace(statement.Text)
	}
	return failing, failing != ""
}

// errorLine returns the line of the batch the server reported err at, and 0
// when it named none. The first message that names a line is the one that
// failed: the messages after it are often its cascade.
func errorLine(err error) int {
	var e gomssqldb.Error
	if !errors.As(err, &e) {
		return 0
	}
	for _, message := range e.All {
		if message.LineNo > 0 {
			return int(message.LineNo)
		}
	}
	if e.LineNo > 0 {
		return int(e.LineNo)
	}
	return 0
}

// toSQLIsolation maps the contract's isolation level onto database/sql's.
func toSQLIsolation(level driver.IsolationLevel) (sql.IsolationLevel, error) {
	switch level {
	case "":
		return sql.LevelDefault, nil
	case driver.ReadUncommitted:
		return sql.LevelReadUncommitted, nil
	case driver.ReadCommitted:
		return sql.LevelReadCommitted, nil
	case driver.RepeatableRead:
		return sql.LevelRepeatableRead, nil
	case driver.Serializable:
		return sql.LevelSerializable, nil
	default:
		return sql.LevelDefault, fmt.Errorf("invalid isolation level for SQL Server: %s", level)
	}
}

// CountAffectedRows returns the optimizer's estimate of the rows the statement
// modifies, read from the plan the server prints for it.
func (d *Driver) CountAffectedRows(ctx context.Context, statement string) (int64, error) {
	// The plan is printed on the session that ran SHOWPLAN_ALL, so the statement
	// must run on the same connection: a pool may hand out another.
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, "SET SHOWPLAN_ALL ON;"); err != nil {
		return 0, err
	}
	defer func() {
		_, _ = conn.ExecContext(context.WithoutCancel(ctx), "SET SHOWPLAN_ALL OFF;")
	}()

	rows, err := conn.QueryContext(ctx, statement)
	if err != nil {
		return 0, unpackMSSQLError(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	estimateIndex := -1
	for i, column := range columns {
		if column == "EstimateRows" {
			estimateIndex = i
			break
		}
	}
	if estimateIndex < 0 {
		return 0, nil
	}

	for rows.Next() {
		scanArgs := make([]any, len(columns))
		for i := range scanArgs {
			scanArgs[i] = new(any)
		}
		var estimate sql.NullFloat64
		scanArgs[estimateIndex] = &estimate
		if err := rows.Scan(scanArgs...); err != nil {
			return 0, err
		}
		if estimate.Valid {
			return roundRows(estimate.Float64), nil
		}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return 0, nil
}

func roundRows(rows float64) int64 {
	if rows >= float64(1<<63-1) {
		return 1<<63 - 1
	}
	return int64(rows + 0.5)
}

func logCommandExecute(log driver.ExecuteLogger, command driver.Command) {
	if log != nil {
		log.CommandExecute(command)
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
