package driver

import (
	"context"
	"database/sql"
	"errors"
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
	return d.executeInTransactionMode(ctx, commands, opts.Log)
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
	var statements []string
	for _, segment := range redshift.Split(script) {
		if !segment.Empty() {
			statements = append(statements, segment.Text)
		}
	}
	if len(statements) == 0 {
		return errors.New("database name not found")
	}

	databaseName, err := getDatabaseInCreateDatabaseStatement(statements[0])
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
	for _, statement := range statements {
		if _, err := d.db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
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
