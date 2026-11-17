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
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/bytebase/omni/driver"
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
		nonTransactional := isNonTransactionStatement(command.Text)
		sessionRole := isSetRoleStatement(command.Text)
		if !nonTransactional && !sessionRole && isSuperuserStatement(command.Text) {
			// Use superuser privilege to run privileged statements.
			ct := command.Text
			if !strings.HasSuffix(strings.TrimRightFunc(ct, unicode.IsSpace), ";") {
				ct += ";"
			}
			command.Text = fmt.Sprintf("SET LOCAL ROLE NONE;%sSET LOCAL ROLE '%s';", ct, owner)
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

	if opts.Transaction == driver.AutoCommit {
		return d.executeInAutoCommitMode(ctx, owner, script, orderedCommands, opts.Log, isPlsql)
	}
	return d.executeInTransactionMode(ctx, owner, script, txCommands, nonTransactionAndSetRoleStmts, opts.Log, isPlsql)
}

func (d *Driver) executeInTransactionMode(
	ctx context.Context,
	owner string,
	script string,
	commands []driver.Command,
	nonTransactionAndSetRoleStmts []driver.Command,
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
			_, err := conn.ExecContext(ctx, fmt.Sprintf("SET SESSION ROLE '%s'", owner)) // NOSONAR(go:S2077) owner is from pg_roles system catalog, not user input
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

			tx, err := conn.Begin(ctx)
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
			if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL ROLE '%s'", owner)); err != nil {
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
			totalRowsAffected, err = d.executeSQLTx(ctx, conn, owner, commands, log)
		}
		if err != nil {
			return 0, err
		}
	}

	// USE SET SESSION ROLE to set the role for the current session.
	if err := crdb.Execute(func() error {
		_, err := conn.ExecContext(ctx, fmt.Sprintf("SET SESSION ROLE '%s'", owner)) // NOSONAR(go:S2077) owner is from pg_roles system catalog, not user input
		return err
	}); err != nil {
		return 0, fmt.Errorf("failed to set role to database owner %q: %w", owner, err)
	}
	// Run non-transaction statements at the end.
	for _, stmt := range nonTransactionAndSetRoleStmts {
		logCommandExecute(log, stmt)
		if err := crdb.Execute(func() error {
			_, err := conn.ExecContext(ctx, stmt.Text)
			return err
		}); err != nil {
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

	if _, err := tx.ExecContext(ctx, fmt.Sprintf("SET LOCAL ROLE '%s'", owner)); err != nil {
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
			_, err := conn.ExecContext(ctx, fmt.Sprintf("SET SESSION ROLE '%s'", owner)) // NOSONAR(go:S2077) owner is from pg_roles system catalog, not user input
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
		_, err := conn.ExecContext(ctx, fmt.Sprintf("SET SESSION ROLE '%s'", owner)) // NOSONAR(go:S2077) owner is from pg_roles system catalog, not user input
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

var (
	// DROP DATABASE cannot run inside a transaction block.
	// DROP DATABASE [ IF EXISTS ] name [ [ WITH ] ( option [, ...] ) ]。
	dropDatabaseReg = regexp.MustCompile(`(?i)DROP DATABASE`)
	// CREATE INDEX CONCURRENTLY cannot run inside a transaction block.
	// CREATE [ UNIQUE ] INDEX [ CONCURRENTLY ] [ [ IF NOT EXISTS ] name ] ON [ ONLY ] table_name [ USING method ] ...
	createIndexReg = regexp.MustCompile(`(?i)CREATE(\s+(UNIQUE\s+)?)INDEX(\s+)CONCURRENTLY`)
	// DROP INDEX CONCURRENTLY cannot run inside a transaction block.
	// DROP INDEX [ CONCURRENTLY ] [ IF EXISTS ] name [, ...] [ CASCADE | RESTRICT ].
	dropIndexReg = regexp.MustCompile(`(?i)DROP(\s+)INDEX(\s+)CONCURRENTLY`)
	// VACUUM cannot run inside a transaction block.
	// VACUUM [ ( option [, ...] ) ] [ table_and_columns [, ...] ]
	// VACUUM [ FULL ] [ FREEZE ] [ VERBOSE ] [ ANALYZE ] [ table_and_columns [, ...] ].
	vacuumReg = regexp.MustCompile(`(?i)^\s*VACUUM`)
	// SET ROLE is a special statement that should be run before any other statements containing inside a transaction block or not.
	setRoleReg = regexp.MustCompile(`(?i)SET\s+((SESSION|LOCAL)\s+)?ROLE`)
)

func isSetRoleStatement(stmt string) bool {
	return len(setRoleReg.FindString(stmt)) > 0
}

func isNonTransactionStatement(stmt string) bool {
	if len(dropDatabaseReg.FindString(stmt)) > 0 {
		return true
	}
	if len(createIndexReg.FindString(stmt)) > 0 {
		return true
	}
	if len(dropIndexReg.FindString(stmt)) > 0 {
		return true
	}
	return len(vacuumReg.FindString(stmt)) > 0
}

func isSuperuserStatement(stmt string) bool {
	upperCaseStmt := strings.ToUpper(strings.TrimLeftFunc(stmt, unicode.IsSpace))
	if strings.HasPrefix(upperCaseStmt, "GRANT") || strings.HasPrefix(upperCaseStmt, "CREATE EXTENSION") || strings.HasPrefix(upperCaseStmt, "CREATE EVENT TRIGGER") || strings.HasPrefix(upperCaseStmt, "COMMENT ON EVENT TRIGGER") {
		return true
	}
	return false
}

func getDatabaseInCreateDatabaseStatement(createDatabaseStatement string) (string, error) {
	raw := strings.TrimRight(createDatabaseStatement, ";")
	raw = strings.TrimPrefix(raw, "CREATE DATABASE")
	tokens := strings.Fields(raw)
	if len(tokens) == 0 {
		return "", fmt.Errorf("database name not found")
	}
	databaseName := strings.TrimLeft(tokens[0], `"`)
	databaseName = strings.TrimRight(databaseName, `"`)
	return databaseName, nil
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
