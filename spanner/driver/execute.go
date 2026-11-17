package driver

import (
	"context"
	"fmt"
	"strings"

	"github.com/bytebase/omni/driver"
)

// Execute runs every statement of script and returns the rows they affected in
// total. A script that holds any DDL statement is applied through the admin API
// as one schema change, which is what Bytebase did and why a script that mixes
// DDL with DML fails; the rest runs as DML, in one read-write transaction when
// ExecuteOptions.Transaction asks for one, and each statement as a partitioned
// update when it asks for autocommit.
//
// ExecuteOptions.Isolation does nothing: Spanner has one isolation level,
// serializable. ExecuteOptions.MaxRetries retries nothing: Spanner reports the
// lock conflicts it detects as ABORTED and leaves the retry to its client, which
// is where Bytebase left it.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	if opts.CreateDatabase {
		// A creation script states its database's dialect by the way it quotes the
		// name, and the rest of the script is that dialect too.
		created, name := createDatabaseStatement(script)
		stmts := statements(script, created)
		if len(stmts) == 0 {
			return 0, fmt.Errorf("expect sanitized SQLs to have at least one entry, original statement: %v", script)
		}
		if !startsWithCreateDatabase(stmts[0].Text, created) {
			return 0, fmt.Errorf("expect the first entry of the sanitized SQLs to start with 'CREATE DATABASE', sql %v", stmts[0].Text)
		}
		extra := make([]string, 0, len(stmts)-1)
		for _, stmt := range stmts[1:] {
			extra = append(extra, stmt.Text)
		}
		if err := d.client.CreateDatabase(ctx, d.instancePath, stmts[0].Text, extra, dialectName(created)); err != nil {
			return 0, fmt.Errorf("failed to create database: %w", err)
		}
		// A driver opened on the database it just created read no dialect at Open
		// (the database was not there), so the creation is what the queries after
		// it split by.
		if d.databaseName == name {
			d.setDialect(created)
		}
		return 0, nil
	}

	stmts := statements(script, d.getDialect())

	// Spanner DDL is not transactional with DML, so a script that holds any DDL
	// statement goes to the admin API whole.
	for _, stmt := range stmts {
		if stmt.Kind == ddl {
			return d.executeDDL(ctx, stmts)
		}
	}

	texts := make([]string, 0, len(stmts))
	for _, stmt := range stmts {
		texts = append(texts, stmt.Text)
	}
	if opts.Transaction == driver.AutoCommit {
		return d.executeInAutoCommitMode(ctx, texts)
	}
	return d.executeInTransactionMode(ctx, texts, opts.Log)
}

func (d *Driver) executeDDL(ctx context.Context, stmts []statement) (int64, error) {
	ddl := make([]string, 0, len(stmts))
	for _, stmt := range stmts {
		ddl = append(ddl, stmt.Text)
	}
	if err := d.client.UpdateDatabaseDdl(ctx, d.databasePath(""), ddl); err != nil {
		return 0, err
	}
	return 0, nil
}

func (d *Driver) executeInTransactionMode(ctx context.Context, stmts []string, log driver.ExecuteLogger) (int64, error) {
	var rowCount int64

	logTransactionControl(log, driver.Begin, nil)

	counts, err := d.client.BatchUpdate(ctx, stmts)
	if err != nil {
		logTransactionControl(log, driver.Rollback, err)
		return 0, err
	}
	for _, count := range counts {
		rowCount += count
	}
	logTransactionControl(log, driver.Commit, nil)
	return rowCount, nil
}

func (d *Driver) executeInAutoCommitMode(ctx context.Context, stmts []string) (int64, error) {
	var rowCount int64
	// Execute statements individually in auto-commit mode.
	for _, stmt := range stmts {
		count, err := d.client.PartitionedUpdate(ctx, stmt)
		if err != nil {
			return rowCount, err
		}
		rowCount += count
	}
	return rowCount, nil
}

// startsWithCreateDatabase reports whether statement begins with CREATE
// DATABASE, after any comment before it: a creation script may carry a header
// comment, which Bytebase's comment-stripping splitter removed before this check.
func startsWithCreateDatabase(statement string, d dialect) bool {
	return firstTwoKeywords(statement, d) == [2]string{"CREATE", "DATABASE"}
}

// createDatabaseStatement reports the dialect a creation script is written in
// and the database name it gives, or GoogleSQL and "" for a script that does not
// open with CREATE DATABASE. The admin API refuses a creation statement written
// in the other dialect, so a name in double quotes is PostgreSQL's and a
// backtick-quoted or unquoted one is GoogleSQL's. The tokens are read with
// GoogleSQL's lexer first, which skips the `#` comment a GoogleSQL script may
// carry, and with PostgreSQL's second, which skips a nested block comment, so a
// script either dialect's comments precede still names its database.
func createDatabaseStatement(script string) (dialect, string) {
	for _, d := range []dialect{googleSQL, postgresSQL} {
		tokens := firstTokens(script, d, 3)
		if len(tokens) < 3 || !strings.EqualFold(tokens[0], "CREATE") || !strings.EqualFold(tokens[1], "DATABASE") {
			continue
		}
		name := strings.Trim(tokens[2], "\"`")
		if strings.HasPrefix(tokens[2], `"`) {
			return postgresSQL, name
		}
		return googleSQL, name
	}
	return googleSQL, ""
}

func logTransactionControl(log driver.ExecuteLogger, op driver.TransactionOp, err error) {
	if log != nil {
		log.TransactionControl(op, err)
	}
}
