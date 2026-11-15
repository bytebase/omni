package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/jackc/pgx/v5/stdlib"

	"github.com/bytebase/omni/driver"
	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
)

// Execute runs script under opts.Transaction, reporting each step to opts.Log
// when non-nil and returning the total affected rows.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	if opts.CreateDatabase {
		if err := d.createDatabaseExecute(ctx, script); err != nil {
			return 0, err
		}
		return 0, nil
	}

	var owner string
	if d.tenantMode {
		var err error
		owner, err = d.CurrentDatabaseOwner(ctx)
		if err != nil {
			return 0, err
		}
	}

	segments := omnipg.Split(script)
	var allCommands []driver.Command
	var txCommands []driver.Command
	var postTxCommands []driver.Command
	hasNonTx := false

	type annotatedCmd struct {
		cmd  driver.Command
		kind stmtExecClass
	}
	var annotated []annotatedCmd
	var singleDo bool
	for _, seg := range segments {
		if seg.Empty() {
			continue
		}
		cls, err := classifyExecStmt(seg.Text)
		if err != nil {
			return 0, err
		}
		if len(segments) == 1 && cls.plpgsqlDo {
			singleDo = true
		}
		if cls.nonTransactional {
			hasNonTx = true
		}
		cmd := driver.Command{
			Text:      seg.Text,
			ByteStart: seg.ByteStart,
			ByteEnd:   seg.ByteEnd,
		}
		annotated = append(annotated, annotatedCmd{cmd: cmd, kind: cls})
	}
	isPlsql := singleDo

	for _, item := range annotated {
		cmd := item.cmd
		elevated := cmd
		if item.kind.superuser && d.tenantMode {
			ct := cmd.Text
			if !strings.HasSuffix(strings.TrimRightFunc(ct, unicode.IsSpace), ";") {
				ct += ";"
			}
			elevated.Text = fmt.Sprintf("SET LOCAL ROLE NONE;%sSET LOCAL ROLE %s;", ct, quoteLiteral(owner))
		}
		allCommands = append(allCommands, elevated)

		if item.kind.nonTransactional {
			postTxCommands = append(postTxCommands, cmd)
			continue
		}
		if item.kind.sessionRole && hasNonTx {
			postTxCommands = append(postTxCommands, cmd)
		}
		txCommands = append(txCommands, elevated)
	}

	run := func() (int64, error) {
		if opts.Transaction == driver.AutoCommit {
			return d.executeInAutoCommitMode(ctx, owner, script, allCommands, opts.Log, isPlsql)
		}
		return d.executeInTransactionMode(ctx, owner, script, txCommands, postTxCommands, opts.Log, isPlsql)
	}

	affectedRows, err := run()
	if err == nil {
		return affectedRows, nil
	}
	var lockErr *lockTimeoutError
	if !errors.As(err, &lockErr) {
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
		affectedRows, err = run()
		if err == nil || !errors.As(err, &lockErr) {
			break
		}
	}
	return affectedRows, err
}

type lockTimeoutError struct {
	err error
}

func (e *lockTimeoutError) Error() string { return e.err.Error() }
func (e *lockTimeoutError) Unwrap() error { return e.err }

func isLockTimeoutError(message string) bool {
	return strings.Contains(message, "canceling statement due to lock timeout")
}

func (d *Driver) executeInTransactionMode(
	ctx context.Context,
	owner string,
	script string,
	commands []driver.Command,
	postTxCommands []driver.Command,
	log driver.ExecuteLogger,
	isPlsql bool,
) (int64, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}
	defer conn.Close()

	if isPlsql {
		if d.tenantMode {
			if _, err := conn.ExecContext(ctx, fmt.Sprintf("SET SESSION ROLE %s", quoteLiteral(owner))); err != nil {
				return 0, fmt.Errorf("failed to set role to database owner %q: %w", owner, err)
			}
		}
		logCommandExecute(log, driver.Command{Text: script, ByteStart: 0, ByteEnd: len(script)})
		if _, err := conn.ExecContext(ctx, script); err != nil {
			logCommandResponse(log, 0, []int64{0}, err)
			return 0, err
		}
		logCommandResponse(log, 0, []int64{0}, nil)
		return 0, nil
	}

	var totalRowsAffected int64
	if len(commands) > 0 {
		var usedPgx bool
		err = conn.Raw(func(driverConn any) error {
			sc, ok := driverConn.(*stdlib.Conn)
			if !ok {
				return nil
			}
			usedPgx = true
			pgxConn := sc.Conn()

			tx, err := pgxConn.Begin(ctx)
			if err != nil {
				logTransactionControl(log, driver.Begin, err)
				return fmt.Errorf("failed to begin transaction: %w", err)
			}
			logTransactionControl(log, driver.Begin, nil)

			committed := false
			defer func() {
				rollbackErr := tx.Rollback(ctx)
				if committed {
					return
				}
				logTransactionControl(log, driver.Rollback, rollbackErr)
			}()

			if d.tenantMode {
				if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL ROLE %s", quoteLiteral(owner))); err != nil {
					return err
				}
			}

			for _, command := range commands {
				logCommandExecute(log, command)

				results, err := tx.Conn().PgConn().Exec(ctx, command.Text).ReadAll()
				if err != nil {
					logCommandResponse(log, 0, nil, err)
					return err
				}

				var rowsAffected int64
				var perResult []int64
				for _, r := range results {
					ra := r.CommandTag.RowsAffected()
					perResult = append(perResult, ra)
					rowsAffected += ra
				}
				logCommandResponse(log, rowsAffected, perResult, nil)
				totalRowsAffected += rowsAffected
			}

			if err := tx.Commit(ctx); err != nil {
				logTransactionControl(log, driver.Commit, err)
				return fmt.Errorf("failed to commit transaction: %w", err)
			}
			logTransactionControl(log, driver.Commit, nil)
			committed = true
			return nil
		})
		if err == nil && !usedPgx {
			totalRowsAffected, err = d.executeSQLTx(ctx, conn, owner, commands, log)
		}
		if err != nil {
			if isLockTimeoutError(err.Error()) {
				return 0, &lockTimeoutError{err: err}
			}
			return 0, err
		}
	}

	if d.tenantMode {
		if _, err := conn.ExecContext(ctx, fmt.Sprintf("SET SESSION ROLE %s", quoteLiteral(owner))); err != nil {
			return 0, fmt.Errorf("failed to set role to database owner %q: %w", owner, err)
		}
	}
	for _, stmt := range postTxCommands {
		logCommandExecute(log, stmt)
		if _, err := conn.ExecContext(ctx, stmt.Text); err != nil {
			logCommandResponse(log, 0, []int64{0}, err)
			return 0, err
		}
		logCommandResponse(log, 0, []int64{0}, nil)
	}
	return totalRowsAffected, nil
}

// executeSQLTx is executeInTransactionMode's transaction body over a generic
// database/sql connector such as enginetest/tape.
func (d *Driver) executeSQLTx(
	ctx context.Context,
	conn *sql.Conn,
	owner string,
	commands []driver.Command,
	log driver.ExecuteLogger,
) (int64, error) {
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		logTransactionControl(log, driver.Begin, err)
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	logTransactionControl(log, driver.Begin, nil)

	committed := false
	defer func() {
		rollbackErr := tx.Rollback()
		if committed {
			return
		}
		if errors.Is(rollbackErr, sql.ErrTxDone) {
			rollbackErr = nil
		}
		logTransactionControl(log, driver.Rollback, rollbackErr)
	}()

	if d.tenantMode {
		if _, err := tx.ExecContext(ctx, fmt.Sprintf("SET LOCAL ROLE %s", quoteLiteral(owner))); err != nil {
			return 0, err
		}
	}

	var totalRowsAffected int64
	for _, command := range commands {
		logCommandExecute(log, command)
		res, err := tx.ExecContext(ctx, command.Text)
		if err != nil {
			logCommandResponse(log, 0, nil, err)
			return 0, err
		}
		ra, err := res.RowsAffected()
		if err != nil {
			ra = 0
		}
		logCommandResponse(log, ra, []int64{ra}, nil)
		totalRowsAffected += ra
	}

	if err := tx.Commit(); err != nil {
		logTransactionControl(log, driver.Commit, err)
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	logTransactionControl(log, driver.Commit, nil)
	committed = true
	return totalRowsAffected, nil
}

func (d *Driver) executeInAutoCommitMode(
	ctx context.Context,
	owner string,
	script string,
	commands []driver.Command,
	log driver.ExecuteLogger,
	isPlsql bool,
) (int64, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}
	defer conn.Close()

	if d.tenantMode {
		if _, err := conn.ExecContext(ctx, fmt.Sprintf("SET SESSION ROLE %s", quoteLiteral(owner))); err != nil {
			return 0, fmt.Errorf("failed to set role to database owner %q: %w", owner, err)
		}
	}

	if isPlsql {
		logCommandExecute(log, driver.Command{Text: script, ByteStart: 0, ByteEnd: len(script)})
		if _, err := conn.ExecContext(ctx, script); err != nil {
			logCommandResponse(log, 0, []int64{0}, err)
			return 0, err
		}
		logCommandResponse(log, 0, []int64{0}, nil)
		return 0, nil
	}

	var totalRowsAffected int64
	for _, stmt := range commands {
		logCommandExecute(log, stmt)
		sqlResult, err := conn.ExecContext(ctx, stmt.Text)
		if err != nil {
			logCommandResponse(log, 0, []int64{0}, err)
			if isLockTimeoutError(err.Error()) {
				return totalRowsAffected, &lockTimeoutError{err: err}
			}
			return totalRowsAffected, err
		}
		rowsAffected, err := sqlResult.RowsAffected()
		if err != nil {
			rowsAffected = 0
		}
		logCommandResponse(log, rowsAffected, []int64{rowsAffected}, nil)
		totalRowsAffected += rowsAffected
	}
	return totalRowsAffected, nil
}

func (d *Driver) createDatabaseExecute(ctx context.Context, script string) error {
	var statements []string
	for _, seg := range omnipg.Split(script) {
		if !seg.Empty() {
			statements = append(statements, seg.Text)
		}
	}
	if len(statements) == 0 {
		return errors.New("database name not found")
	}

	databaseName, err := getDatabaseInCreateDatabaseStatement(statements[0])
	if err != nil {
		return err
	}
	for _, stmt := range statements[1:] {
		if _, err := omnipg.Parse(stmt); err != nil {
			return err
		}
	}
	databases, err := d.getDatabases(ctx)
	if err != nil {
		return err
	}
	for _, database := range databases {
		if database.Name == databaseName {
			return nil
		}
	}
	for _, stmt := range statements {
		if _, err := d.db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

func getDatabaseInCreateDatabaseStatement(createDatabaseStatement string) (string, error) {
	stmts, err := omnipg.Parse(createDatabaseStatement)
	if err != nil {
		return "", err
	}
	if len(stmts) != 1 {
		return "", errors.New("database name not found")
	}
	createDB, ok := stmts[0].AST.(*ast.CreatedbStmt)
	if !ok || createDB.Dbname == "" {
		return "", errors.New("database name not found")
	}
	return createDB.Dbname, nil
}

type stmtExecClass struct {
	plpgsqlDo        bool
	nonTransactional bool
	sessionRole      bool
	superuser        bool
}

func classifyExecStmt(stmt string) (stmtExecClass, error) {
	stmts, err := omnipg.Parse(stmt)
	if err != nil {
		return stmtExecClass{}, err
	}
	if len(stmts) != 1 {
		return stmtExecClass{}, nil
	}
	switch n := stmts[0].AST.(type) {
	case *ast.DoStmt:
		return stmtExecClass{plpgsqlDo: true}, nil
	case *ast.DropdbStmt:
		return stmtExecClass{nonTransactional: true}, nil
	case *ast.IndexStmt:
		return stmtExecClass{nonTransactional: n.Concurrent}, nil
	case *ast.DropStmt:
		return stmtExecClass{nonTransactional: n.Concurrent}, nil
	case *ast.VacuumStmt:
		return stmtExecClass{nonTransactional: n.IsVacuumCmd}, nil
	case *ast.VariableSetStmt:
		if strings.EqualFold(n.Name, "role") && !n.IsLocal {
			return stmtExecClass{sessionRole: true}, nil
		}
	case *ast.GrantStmt:
		return stmtExecClass{superuser: n.IsGrant}, nil
	case *ast.GrantRoleStmt:
		return stmtExecClass{superuser: n.IsGrant}, nil
	case *ast.CreateExtensionStmt, *ast.CreateEventTrigStmt:
		return stmtExecClass{superuser: true}, nil
	case *ast.CommentStmt:
		if n.Objtype == ast.OBJECT_EVENT_TRIGGER || n.Objtype == ast.OBJECT_EXTENSION {
			return stmtExecClass{superuser: true}, nil
		}
	}
	return stmtExecClass{}, nil
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
