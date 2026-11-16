package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/driver"
)

// Execute runs every statement of script and returns the rows they affected in
// total, reporting each step to opts.Log as it goes.
//
// Doris accepts no DDL inside a transaction and commits each data change, so a
// transaction cannot roll a script back: a failed script leaves the rows it
// already ran.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	commands := splitExecuteCommands(script)
	if len(commands) == 0 {
		return 0, nil
	}

	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = conn.Close() }()

	connID, err := d.readConnectionID(ctx, conn)
	if err != nil {
		return 0, err
	}

	if opts.Transaction == driver.AutoCommit || opts.CreateDatabase {
		return d.executeInAutoCommitMode(ctx, conn, commands, opts.Log, connID)
	}
	return d.executeInTransactionMode(ctx, conn, commands, opts.Log, connID, opts.Isolation)
}

// splitExecuteCommands cuts script into the statements Execute sends. A DELIMITER
// directive is client-side, so a script that carries one is cut by it and the
// directives never reach the server, which does not accept them.
func splitExecuteCommands(script string) []driver.Command {
	if containsDelimiterDirective(script) {
		return delimiterCommands(script)
	}
	return splitCommands(script)
}

func (d *Driver) executeInTransactionMode(
	ctx context.Context,
	conn *sql.Conn,
	commands []driver.Command,
	log driver.ExecuteLogger,
	connID uint64,
	isolation driver.IsolationLevel,
) (int64, error) {
	sqlIsolation, err := toSQLIsolation(isolation)
	if err != nil {
		return 0, err
	}
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sqlIsolation})
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

	var totalRowsAffected int64
	for _, cmd := range commands {
		logCommandExecute(log, cmd)
		sqlResult, err := tx.ExecContext(ctx, cmd.Text)
		if err != nil {
			d.killQueryOnCancel(err, connID)
			logCommandResponse(log, 0, nil, err)
			return 0, err
		}
		rowsAffected, perResult := extractAffectedRows(sqlResult)
		totalRowsAffected += rowsAffected
		logCommandResponse(log, rowsAffected, perResult, nil)
	}

	if err := tx.Commit(); err != nil {
		logTransactionControl(log, driver.Commit, err)
		return 0, fmt.Errorf("failed to commit execute transaction: %w", err)
	}
	logTransactionControl(log, driver.Commit, nil)
	committed = true
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
	for _, cmd := range commands {
		logCommandExecute(log, cmd)
		sqlResult, err := conn.ExecContext(ctx, cmd.Text)
		if err != nil {
			d.killQueryOnCancel(err, connID)
			logCommandResponse(log, 0, nil, err)
			return 0, err
		}
		rowsAffected, perResult := extractAffectedRows(sqlResult)
		totalRowsAffected += rowsAffected
		logCommandResponse(log, rowsAffected, perResult, nil)
	}
	return totalRowsAffected, nil
}

// extractAffectedRows reads a result's affected rows, per result where the
// engine reports several.
func extractAffectedRows(res sql.Result) (int64, []int64) {
	if mr, ok := res.(interface{ AllRowsAffected() []int64 }); ok {
		all := mr.AllRowsAffected()
		var total int64
		perResult := make([]int64, len(all))
		for i, a := range all {
			total += a
			perResult[i] = a
		}
		return total, perResult
	}
	affected, err := res.RowsAffected()
	if err != nil {
		affected = 0
	}
	return affected, []int64{affected}
}

// toSQLIsolation maps the contract's level onto database/sql's, refusing a
// level the engine does not have.
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
		return sql.LevelDefault, fmt.Errorf("invalid isolation level: %s. Supported levels: READ UNCOMMITTED, READ COMMITTED, REPEATABLE READ, SERIALIZABLE", level)
	}
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

// delimiterCommands cuts script at its client-side DELIMITER directives and the
// terminators they name, leaving each statement with the byte range it has in
// script. The parser's splitter knows only ";", so a directive script has to be
// cut here.
func delimiterCommands(script string) []driver.Command {
	var commands []driver.Command
	delimiter := ";"
	start := 0
	emit := func(end int) {
		// A ";" terminator is part of the statement the server receives, as in
		// the parser's own splitting; a directive's terminator is not.
		if delimiter == ";" && end < len(script) && script[end] == ';' {
			end++
		}
		if text := strings.TrimSpace(script[start:end]); text != "" {
			commands = append(commands, driver.Command{Text: text, ByteStart: start, ByteEnd: end})
		}
	}

	i := 0
	started := false
	for i < len(script) {
		// A directive is recognized only where a statement may begin: comments
		// before it belong to no statement, so they do not start one.
		if !started {
			if terminator, next, ok := delimiterDirective(script, i); ok {
				delimiter = terminator
				start, i = next, next
				continue
			}
		}
		switch {
		case script[i] == '\'' || script[i] == '"' || script[i] == '`':
			started = true
			quote := script[i]
			i++
			for i < len(script) {
				if script[i] == '\\' && quote != '`' && i+1 < len(script) {
					i += 2
					continue
				}
				if script[i] == quote {
					if i+1 < len(script) && script[i+1] == quote {
						i += 2
						continue
					}
					i++
					break
				}
				i++
			}
		case strings.HasPrefix(script[i:], "/*"):
			if end := strings.Index(script[i+2:], "*/"); end >= 0 {
				i += 2 + end + 2
			} else {
				i = len(script)
			}
		case dashComment(script, i), script[i] == '#',
			strings.HasPrefix(script[i:], "//") && delimiter != "//":
			if nl := strings.IndexByte(script[i:], '\n'); nl >= 0 {
				i += nl + 1
			} else {
				i = len(script)
			}
		case strings.HasPrefix(script[i:], delimiter):
			emit(i)
			i += len(delimiter)
			start = i
			started = false
		default:
			if script[i] != ' ' && script[i] != '\t' && script[i] != '\n' && script[i] != '\r' {
				started = true
			}
			i++
		}
	}
	emit(len(script))
	return commands
}

// dashComment reports whether script[i:] opens a line comment: "--" followed by
// a space, a tab, a line break or the end of the input, as both engines' lexers
// require, so "1--2" is a subtraction.
func dashComment(script string, i int) bool {
	if !strings.HasPrefix(script[i:], "--") {
		return false
	}
	if i+2 >= len(script) {
		return true
	}
	switch script[i+2] {
	case ' ', '\t', '\n', '\r':
		return true
	}
	return false
}

// delimiterDirective returns the terminator a DELIMITER directive names at j and
// the offset just past the directive's line, and false when the text at j is not
// the start of one.
func delimiterDirective(script string, j int) (terminator string, next int, ok bool) {
	for k := j; k > 0 && script[k-1] != '\n'; k-- {
		if script[k-1] != ' ' && script[k-1] != '\t' {
			return "", 0, false
		}
	}
	const keyword = "DELIMITER"
	if len(script)-j < len(keyword) || !strings.EqualFold(script[j:j+len(keyword)], keyword) {
		return "", 0, false
	}
	j += len(keyword)
	if j >= len(script) || (script[j] != ' ' && script[j] != '\t') {
		return "", 0, false
	}
	for j < len(script) && (script[j] == ' ' || script[j] == '\t') {
		j++
	}
	begin := j
	for j < len(script) && script[j] != '\n' && script[j] != '\r' {
		j++
	}
	terminator = strings.TrimSpace(script[begin:j])
	if terminator == "" {
		return "", 0, false
	}
	for j < len(script) && (script[j] == '\n' || script[j] == '\r') {
		j++
	}
	return terminator, j, true
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
