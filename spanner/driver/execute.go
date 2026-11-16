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
		stmts := statements(script)
		if len(stmts) == 0 {
			return 0, fmt.Errorf("expect sanitized SQLs to have at least one entry, original statement: %v", script)
		}
		if !strings.HasPrefix(strings.TrimSpace(stmts[0].Text), "CREATE DATABASE") {
			return 0, fmt.Errorf("expect the first entry of the sanitized SQLs to start with 'CREATE DATABASE', sql %v", stmts[0].Text)
		}
		extra := make([]string, 0, len(stmts)-1)
		for _, stmt := range stmts[1:] {
			extra = append(extra, stmt.Text)
		}
		if err := d.client.CreateDatabase(ctx, d.instancePath, stmts[0].Text, extra); err != nil {
			return 0, fmt.Errorf("failed to create database: %w", err)
		}
		return 0, nil
	}

	stmts := statements(script)

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

func logTransactionControl(log driver.ExecuteLogger, op driver.TransactionOp, err error) {
	if log != nil {
		log.TransactionControl(op, err)
	}
}
