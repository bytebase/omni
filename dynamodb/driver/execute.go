package driver

import (
	"context"

	"github.com/bytebase/omni/driver"
)

// Execute runs each statement of script as one ExecuteTransaction and reports
// each step to opts.Log as it goes. It returns 0: DynamoDB answers no affected
// row count for a PartiQL write.
//
// ExecuteTransaction is the write API this driver sends every statement
// through. It is a transaction of one statement: a script's statements do not
// share a transaction, because the API cannot hold a set of PartiQL statements
// that write the same item more than once, and DynamoDB has no script-level
// transaction to open. A statement whose transaction fails ends the script.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	for _, statement := range splitStatements(script) {
		if opts.Log != nil {
			opts.Log.CommandExecute(driver.Command{
				Text:      statement.Text,
				ByteStart: statement.ByteStart,
				ByteEnd:   statement.ByteEnd,
			})
		}
		err := d.client.ExecuteTransaction(ctx, statement.Text)
		if err != nil {
			if opts.Log != nil {
				opts.Log.CommandResponse(0, []int64{0}, err)
			}
			return 0, err
		}
		if opts.Log != nil {
			opts.Log.CommandResponse(0, []int64{0}, nil)
		}
	}

	return 0, nil
}
