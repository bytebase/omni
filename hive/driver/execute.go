package driver

import (
	"context"
	"fmt"
	"strings"

	"github.com/bytebase/omni/driver"
)

// Execute runs every statement of script and returns the rows they affected in
// total, reporting each step to opts.Log as it goes.
//
// Hive commits every statement itself — only ACID tables are transactional, and
// HiveServer2 4.0 refuses BEGIN, COMMIT and ROLLBACK — so opts.Transaction and
// opts.Isolation change nothing: a failed script leaves the statements it
// already ran. The script runs on one connection, so a USE it holds stays
// selected for its next statement.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	commands, err := splitCommands(script)
	if err != nil {
		return 0, fmt.Errorf("failed to split statements: %w", err)
	}
	if len(commands) == 0 {
		return 0, nil
	}

	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = conn.Close() }()

	var totalAffected int64
	for _, command := range commands {
		// Hive may not report a row count for a statement that changes rows.
		query := strings.TrimRight(command.Text, ";")
		logCommandExecute(opts.Log, command)
		result, err := conn.ExecContext(ctx, query)
		if err != nil {
			logCommandResponse(opts.Log, 0, nil, err)
			return totalAffected, fmt.Errorf("failed to execute statement: %s: %w", query, err)
		}
		affected, _ := result.RowsAffected()
		totalAffected += affected
		logCommandResponse(opts.Log, affected, []int64{affected}, nil)
	}

	return totalAffected, nil
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
