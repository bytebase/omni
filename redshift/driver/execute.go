package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/redshift"
	"github.com/bytebase/omni/redshift/ast"
)

// Execute runs script under opts.Transaction, reporting each step to opts.Log
// when it is non-nil and returning the total affected rows.
//
// The script is split lexically, so every statement is sent whole whatever its
// text looks like; bounding how many statements a sheet may hold is the
// caller's. A failing statement ends the script: the transaction is rolled back,
// or, in autocommit mode, the statements before it stay committed.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	if d.datashare {
		return 0, errors.New("datashare database cannot be updated")
	}
	if opts.CreateDatabase {
		if err := d.createDatabaseExecute(ctx, script); err != nil {
			return 0, err
		}
		return 0, nil
	}

	var commands []driver.Command
	for _, segment := range redshift.Split(script) {
		if segment.Empty() {
			continue
		}
		commands = append(commands, driver.Command{
			Text:      segment.Text,
			ByteStart: segment.ByteStart,
			ByteEnd:   segment.ByteEnd,
		})
	}
	if len(commands) == 0 {
		return 0, nil
	}

	if opts.Transaction == driver.AutoCommit {
		return d.executeInAutoCommitMode(ctx, commands, opts.Log)
	}
	// Redshift commits some statements whatever transaction they run in, so a
	// script holding one cannot be rolled back as a whole.
	for _, command := range commands {
		if err := refuseImplicitCommit(command.Text); err != nil {
			return 0, err
		}
	}
	return d.executeInTransactionMode(ctx, commands, opts.Log)
}

// refuseImplicitCommit reports the statements a transactional script cannot
// hold: TRUNCATE and VACUUM commit the transaction they run in, a database
// statement cannot run inside one at all, and a transaction-control statement
// would commit or end the transaction this call promised to own. ANALYZE shares
// VACUUM's parse node but runs inside one on a cluster and rolls back with it,
// so only VACUUM is refused.
func refuseImplicitCommit(statement string) error {
	stmts, err := redshift.Parse(statement)
	if err != nil {
		return err
	}
	for _, stmt := range stmts {
		switch node := stmt.AST.(type) {
		case *ast.TruncateStmt, *ast.CreatedbStmt, *ast.DropdbStmt:
			return fmt.Errorf("the script holds a statement Redshift commits on its own (%s), so it cannot run in one transaction; run it in autocommit mode", stmt.Text)
		case *ast.VacuumStmt:
			if node.IsVacuumCmd {
				return fmt.Errorf("the script holds a statement Redshift commits on its own (%s), so it cannot run in one transaction; run it in autocommit mode", stmt.Text)
			}
		case *ast.TransactionStmt:
			// BEGIN, COMMIT, END and ROLLBACK in a script would decide the
			// transaction this call promised to own.
			return fmt.Errorf("the script holds its own transaction control (%s), so it cannot run in one transaction; run it in autocommit mode", stmt.Text)
		}
	}
	return nil
}

// executeInTransactionMode runs commands in one transaction and commits it.
func (d *Driver) executeInTransactionMode(ctx context.Context, commands []driver.Command, log driver.ExecuteLogger) (int64, error) {
	tx, err := d.db.BeginTx(ctx, nil)
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
		sqlResult, err := tx.ExecContext(ctx, command.Text)
		if err != nil {
			logCommandResponse(log, 0, nil, err)
			return 0, err
		}
		rowsAffected, err := sqlResult.RowsAffected()
		if err != nil {
			// RowsAffected is only meaningful for DML, and the driver cannot
			// tell it from DDL here.
			rowsAffected = 0
		}
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
// search path, a temporary table) keeps it for the statements after it.
func (d *Driver) executeInAutoCommitMode(ctx context.Context, commands []driver.Command, log driver.ExecuteLogger) (int64, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var totalRowsAffected int64
	for _, command := range commands {
		logCommandExecute(log, command)
		sqlResult, err := conn.ExecContext(ctx, command.Text)
		if err != nil {
			logCommandResponse(log, 0, nil, err)
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

// createDatabaseExecute creates the database the script's first statement names,
// and runs the rest of the script in it. A database the server already lists is
// not created again, so the step can run twice.
func (d *Driver) createDatabaseExecute(ctx context.Context, script string) error {
	statements, databaseName, err := creationStatements(script)
	if err != nil {
		return err
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
	_, err = d.db.ExecContext(ctx, statements[0])
	return err
}

// creationStatements reads the statement a database-creation step runs, and the
// database it creates. A script that also holds initialization is refused: the
// statements after CREATE DATABASE would have to run in the database it creates,
// and a retry cannot tell a database that was initialized from one that was
// created and left empty.
func creationStatements(script string) ([]string, string, error) {
	var statements []string
	for _, segment := range redshift.Split(script) {
		if !segment.Empty() {
			statements = append(statements, segment.Text)
		}
	}
	if len(statements) == 0 {
		return nil, "", errors.New("database name not found")
	}
	databaseName, err := getDatabaseInCreateDatabaseStatement(statements[0])
	if err != nil {
		return nil, "", err
	}
	if len(statements) > 1 {
		return nil, "", fmt.Errorf("the script holds %d statement(s) after CREATE DATABASE %s: the step creates the database, and what initializes it runs as a script of its own", len(statements)-1, databaseName)
	}
	return statements, databaseName, nil
}

// getDatabaseInCreateDatabaseStatement reads the name a CREATE DATABASE
// statement creates.
func getDatabaseInCreateDatabaseStatement(createDatabaseStatement string) (string, error) {
	stmts, err := redshift.Parse(createDatabaseStatement)
	if err != nil {
		return "", err
	}
	if len(stmts) != 1 {
		return "", errors.New("database name not found")
	}
	createDB, ok := stmts[0].AST.(*ast.CreatedbStmt)
	if !ok || strings.TrimSpace(createDB.Dbname) == "" {
		return "", errors.New("database name not found")
	}
	return createDB.Dbname, nil
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
