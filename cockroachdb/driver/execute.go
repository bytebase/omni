package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	crrawparsertree "github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	crrawanalysis "github.com/bytebase/omni/cockroachdb/analysis"
	"github.com/bytebase/omni/driver"
)

// noRetryError marks a failure after the main transaction committed, where
// rerunning the script would replay changes that are already durable.
type noRetryError struct{ err error }

func (e *noRetryError) Error() string { return e.err.Error() }
func (e *noRetryError) Unwrap() error { return e.err }

// Execute runs script under opts.Transaction, reporting each step to opts.Log
// when non-nil and returning the total affected rows.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	if opts.CreateDatabase {
		if err := d.createDatabaseExecute(ctx, script); err != nil {
			return 0, err
		}
		return 0, nil
	}

	owner, err := d.GetCurrentDatabaseOwner(ctx)
	if err != nil {
		return 0, err
	}

	var commands []driver.Command
	var nonTransactionAndSetRoleStmts []driver.Command
	var isPlsql bool

	singleSQLs, err := splitStatements(script)
	if err != nil {
		return 0, err
	}
	for _, singleSQL := range singleSQLs {
		commands = append(commands, driver.Command{
			Text:      singleSQL.Text,
			ByteStart: singleSQL.ByteStart,
			ByteEnd:   singleSQL.ByteEnd,
		})
	}

	// If the statement is a single statement and is a PL/pgSQL block,
	// we should execute it as a single statement without transaction.
	// If the statement is a PL/pgSQL block, we should execute it as a single statement.
	// https://www.postgresql.org/docs/current/plpgsql-control-structures.html
	if len(singleSQLs) == 1 && isPlSQLBlock(singleSQLs[0].Text) {
		isPlsql = true
	}

	var orderedCommands []driver.Command
	var txCommands []driver.Command
	for _, command := range commands {
		nonTransactional, sessionRole, superuser := classifyStatement(command.Text)
		if !nonTransactional && !sessionRole && superuser {
			// Use superuser privilege to run privileged statements.
			ct := command.Text
			if !strings.HasSuffix(strings.TrimRightFunc(ct, unicode.IsSpace), ";") {
				ct += ";"
			}
			command.Text = fmt.Sprintf("SET LOCAL ROLE NONE;%sSET LOCAL ROLE %s;", ct, quoteLiteral(owner))
		}
		orderedCommands = append(orderedCommands, command)
		if nonTransactional {
			nonTransactionAndSetRoleStmts = append(nonTransactionAndSetRoleStmts, command)
			continue
		}
		txCommands = append(txCommands, command)
		if sessionRole {
			nonTransactionAndSetRoleStmts = append(nonTransactionAndSetRoleStmts, command)
		}
	}

	run := func() (int64, error) {
		if opts.Transaction == driver.AutoCommit {
			return d.executeInAutoCommitMode(ctx, owner, script, orderedCommands, opts.Log, isPlsql)
		}
		return d.executeInTransactionMode(ctx, owner, script, txCommands, nonTransactionAndSetRoleStmts, opts.Isolation, opts.Log, isPlsql)
	}
	affected, err := run()
	var noRetry *noRetryError
	if opts.Transaction == driver.AutoCommit || errors.As(err, &noRetry) || !isLockTimeoutError(err) {
		return affected, err
	}
	for i := 0; i < opts.MaxRetries; i++ {
		if opts.Log != nil {
			opts.Log.Retry(err, i+1, opts.MaxRetries)
		}
		affected, err = run()
		if err == nil || errors.As(err, &noRetry) || !isLockTimeoutError(err) {
			break
		}
	}
	return affected, err
}

func isLockTimeoutError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "canceling statement due to lock timeout")
}

func (d *Driver) executeInTransactionMode(
	ctx context.Context,
	owner string,
	script string,
	commands []driver.Command,
	nonTransactionAndSetRoleStmts []driver.Command,
	isolation driver.IsolationLevel,
	log driver.ExecuteLogger,
	isPlsql bool,
) (int64, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}
	defer conn.Close()

	if isPlsql {
		// USE SET SESSION ROLE to set the role for the current session.
		if err := crdb.Execute(func() error {
			_, err := conn.ExecContext(ctx, fmt.Sprintf("SET SESSION ROLE %s", quoteLiteral(owner))) // NOSONAR(go:S2077) owner is from pg_roles system catalog, not user input
			return err
		}); err != nil {
			return 0, fmt.Errorf("failed to set role to database owner %q: %w", owner, err)
		}
		logCommandExecute(log, driver.Command{Text: script, ByteStart: 0, ByteEnd: len(script)})
		if err := crdb.Execute(func() error {
			_, err := conn.ExecContext(ctx, script)
			return err
		}); err != nil {
			logCommandResponse(log, 0, []int64{0}, err)
			return 0, err
		}
		logCommandResponse(log, 0, []int64{0}, nil)

		return 0, nil
	}

	totalRowsAffected := int64(0)

	totalCommands := len(commands)
	if totalCommands > 0 {
		usedPgx := false
		err = conn.Raw(func(driverConn any) error {
			sc, ok := driverConn.(*stdlib.Conn)
			if !ok {
				return nil
			}
			usedPgx = true
			conn := sc.Conn()

			iso, err := pgxIsolation(isolation)
			if err != nil {
				return err
			}
			tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: iso})
			if err != nil {
				logTransactionControl(log, driver.Begin, err)
				return fmt.Errorf("failed to begin transaction: %w", err)
			}
			logTransactionControl(log, driver.Begin, nil)

			committed := false
			defer func() {
				err := tx.Rollback(ctx)
				if committed {
					return
				}
				logTransactionControl(log, driver.Rollback, err)
			}()

			// Set the current transaction role to the database owner so that the owner of created objects will be the same as the database owner.
			if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL ROLE %s", quoteLiteral(owner))); err != nil {
				return err
			}

			for _, command := range commands {
				logCommandExecute(log, command)

				rr := tx.Conn().PgConn().Exec(ctx, command.Text)
				results, err := rr.ReadAll()
				if err != nil {
					logCommandResponse(log, 0, nil, err)
					return err
				}

				var rowsAffected int64
				var allRowsAffected []int64
				for _, result := range results {
					ra := result.CommandTag.RowsAffected()
					allRowsAffected = append(allRowsAffected, ra)
					rowsAffected += ra
				}
				logCommandResponse(log, rowsAffected, allRowsAffected, nil)

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
			totalRowsAffected, err = d.executeSQLTx(ctx, conn, owner, commands, isolation, log)
		}
		if err != nil {
			return 0, err
		}
	}

	// USE SET SESSION ROLE to set the role for the current session.
	if err := crdb.Execute(func() error {
		_, err := conn.ExecContext(ctx, fmt.Sprintf("SET SESSION ROLE %s", quoteLiteral(owner))) // NOSONAR(go:S2077) owner is from pg_roles system catalog, not user input
		return err
	}); err != nil {
		return 0, &noRetryError{err: fmt.Errorf("failed to set role to database owner %q: %w", owner, err)}
	}
	// Run non-transaction statements at the end.
	for _, stmt := range nonTransactionAndSetRoleStmts {
		logCommandExecute(log, stmt)
		if err := crdb.Execute(func() error {
			_, err := conn.ExecContext(ctx, stmt.Text)
			return err
		}); err != nil {
			logCommandResponse(log, 0, []int64{0}, err)
			return 0, &noRetryError{err: err}
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
	isolation driver.IsolationLevel,
	log driver.ExecuteLogger,
) (int64, error) {
	iso, err := sqlIsolation(isolation)
	if err != nil {
		return 0, err
	}
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: iso})
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

	if _, err := tx.ExecContext(ctx, fmt.Sprintf("SET LOCAL ROLE %s", quoteLiteral(owner))); err != nil {
		return 0, err
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

	if isPlsql {
		// USE SET SESSION ROLE to set the role for the current session.
		if err := crdb.Execute(func() error {
			_, err := conn.ExecContext(ctx, fmt.Sprintf("SET SESSION ROLE %s", quoteLiteral(owner))) // NOSONAR(go:S2077) owner is from pg_roles system catalog, not user input
			return err
		}); err != nil {
			return 0, fmt.Errorf("failed to set role to database owner %q: %w", owner, err)
		}
		logCommandExecute(log, driver.Command{Text: script, ByteStart: 0, ByteEnd: len(script)})
		if err := crdb.Execute(func() error {
			_, err := conn.ExecContext(ctx, script)
			return err
		}); err != nil {
			return 0, err
		}
		logCommandResponse(log, 0, []int64{0}, nil)
		return 0, nil
	}

	// USE SET SESSION ROLE to set the role for the current session.
	if err := crdb.Execute(func() error {
		_, err := conn.ExecContext(ctx, fmt.Sprintf("SET SESSION ROLE %s", quoteLiteral(owner))) // NOSONAR(go:S2077) owner is from pg_roles system catalog, not user input
		return err
	}); err != nil {
		return 0, fmt.Errorf("failed to set role to database owner %q: %w", owner, err)
	}

	totalRowsAffected := int64(0)
	// Execute all statements individually in auto-commit mode, in script order.
	for _, stmt := range commands {
		logCommandExecute(log, stmt)

		if err := crdb.Execute(func() error {
			sqlResult, err := conn.ExecContext(ctx, stmt.Text)
			if err != nil {
				logCommandResponse(log, 0, []int64{0}, err)
				return err
			}

			rowsAffected, err := sqlResult.RowsAffected()
			if err != nil {
				// CockroachDB returns error for statements that don't support RowsAffected
				rowsAffected = 0
			}

			logCommandResponse(log, rowsAffected, []int64{rowsAffected}, nil)
			totalRowsAffected += rowsAffected
			return nil
		}); err != nil {
			return totalRowsAffected, err
		}
	}

	return totalRowsAffected, nil
}

func (d *Driver) createDatabaseExecute(ctx context.Context, statement string) error {
	databaseName, err := getDatabaseInCreateDatabaseStatement(statement)
	if err != nil {
		return err
	}
	databases, err := d.getDatabases(ctx)
	if err != nil {
		return err
	}
	for _, database := range databases {
		if database.Name == databaseName {
			// Database already exists.
			return nil
		}
	}

	statements, err := splitStatements(statement)
	if err != nil {
		return err
	}
	for _, s := range statements {
		if err := crdb.Execute(func() error {
			_, err := d.db.ExecContext(ctx, s.Text)
			return err
		}); err != nil {
			return err
		}
	}
	return nil
}

func pgxIsolation(level driver.IsolationLevel) (pgx.TxIsoLevel, error) {
	switch level {
	case "":
		return "", nil
	case driver.ReadUncommitted:
		return pgx.ReadUncommitted, nil
	case driver.ReadCommitted:
		return pgx.ReadCommitted, nil
	case driver.RepeatableRead:
		return pgx.RepeatableRead, nil
	case driver.Serializable:
		return pgx.Serializable, nil
	default:
		return "", fmt.Errorf("unsupported isolation level %q", level)
	}
}

func sqlIsolation(level driver.IsolationLevel) (sql.IsolationLevel, error) {
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
		return sql.LevelDefault, fmt.Errorf("unsupported isolation level %q", level)
	}
}

// vacuumReg is the one statement classification CockroachDB's parser cannot
// answer: VACUUM has no tree node, so the anchored keyword is enough.
var vacuumReg = regexp.MustCompile(`(?i)^\s*VACUUM`)

// classifyStatement reports what Execute must do with one statement: run it
// outside a transaction, run it once as a session role change, or run it with
// superuser privileges. CockroachDB's own AST decides, so a keyword inside a
// string literal or comment cannot change the classification.
func classifyStatement(stmt string) (nonTransactional, sessionRole, superuser bool) {
	parsed, err := crrawanalysis.Parse(stmt)
	if err != nil || len(parsed) != 1 {
		return false, false, false
	}
	switch node := parsed[0].AST.(type) {
	case *crrawparsertree.DropDatabase:
		return true, false, false
	case *crrawparsertree.CreateIndex:
		return node.Concurrently, false, false
	case *crrawparsertree.DropIndex:
		return node.Concurrently, false, false
	case *crrawparsertree.SetVar:
		return false, strings.EqualFold(node.Name, "role") && !node.Local, false
	case *crrawparsertree.Grant, *crrawparsertree.GrantRole, *crrawparsertree.CreateExtension:
		return false, false, true
	}
	return vacuumReg.MatchString(stmt), false, false
}

func getDatabaseInCreateDatabaseStatement(createDatabaseStatement string) (string, error) {
	parsed, err := crrawanalysis.Parse(createDatabaseStatement)
	if err != nil {
		return "", err
	}
	for _, stmt := range parsed {
		if create, ok := stmt.AST.(*crrawparsertree.CreateDatabase); ok && create.Name != "" {
			return string(create.Name), nil
		}
	}
	return "", fmt.Errorf("database name not found")
}

// GetCurrentDatabaseOwner gets the role of the current database.
func (d *Driver) GetCurrentDatabaseOwner(ctx context.Context) (string, error) {
	const query = `
		SELECT
			u.rolname
		FROM
			pg_roles AS u JOIN pg_database AS d ON (d.datdba = u.oid)
		WHERE
			d.datname = current_database();
		`
	var owner string
	if err := d.db.QueryRowContext(ctx, query).Scan(&owner); err != nil {
		return "", err
	}
	return owner, nil
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
