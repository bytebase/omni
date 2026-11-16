package driver

import (
	"context"
	"fmt"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/mongo"
)

// Execute runs every statement of script in order through gomongo, reporting
// each step to opts.Log, and returns the rows they affected in total, which
// MongoDB reports as zero: a mongosh write answers with its own counts, and
// this driver does not read them.
//
// Execute runs each statement on its own whatever opts.Transaction,
// opts.Isolation and opts.MaxRetries say.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	statements, err := mongo.Parse(script)
	if err != nil {
		// Log the unparseable input as a failed command so the parse error is
		// visible in the task run log instead of the statement silently
		// disappearing from it (BYT-9950).
		logCommandExecute(opts.Log, driver.Command{Text: script, ByteStart: 0, ByteEnd: len(script)})
		logCommandResponse(opts.Log, 0, nil, err)
		return 0, fmt.Errorf("failed to split MongoDB statement: %w", err)
	}
	if len(statements) == 0 {
		return 0, nil
	}

	for _, statement := range statements {
		logCommandExecute(opts.Log, driver.Command{
			Text:      statement.Text,
			ByteStart: statement.ByteStart,
			ByteEnd:   statement.ByteEnd,
		})
		if _, err := d.client.Execute(ctx, d.databaseName, statement.Text, 0); err != nil {
			logCommandResponse(opts.Log, 0, nil, err)
			return 0, err
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
