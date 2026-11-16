package driver

import (
	"context"
	"fmt"

	"github.com/bytebase/omni/cassandra"
	"github.com/bytebase/omni/driver"
)

// Execute runs every statement of script in order and returns the rows they
// affected, which Cassandra never reports: the call answers 0 whatever the
// script did, as Bytebase's did.
//
// A statement that fails ends the script. Each statement is reported to
// opts.Log as a command with its range in the script and a response, so a task
// run's log holds the steps the script ran.
//
// ExecuteOptions' Transaction, Isolation, CreateDatabase and MaxRetries do
// nothing: Cassandra has no multi-statement transaction, no isolation level
// to choose, and no lock timeout to retry.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	c, err := d.conn()
	if err != nil {
		return 0, err
	}
	for _, segment := range cassandra.Split(script) {
		if segment.Empty {
			continue
		}
		statement := trimStatement(segment.Text)
		if statement == "" || holdsNoStatement(statement) {
			continue
		}
		logCommandExecute(opts.Log, driver.Command{
			Text:      statement,
			ByteStart: segment.ByteStart,
			ByteEnd:   segment.ByteEnd,
		})
		if err := c.Exec(ctx, statement); err != nil {
			logCommandResponse(opts.Log, 0, nil, err)
			return 0, fmt.Errorf("failed to execute: %w", err)
		}
		logCommandResponse(opts.Log, 0, nil, nil)
	}
	return 0, nil
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
