package driver

import (
	"context"
	"fmt"

	"github.com/bytebase/omni/driver"
)

// Execute runs every statement of script in order and returns the rows they
// affected in total, reporting each step to opts.Log as it goes.
//
// Databricks' Statement Execution API exposes no transaction control: a
// statement commits itself. ExecuteOptions.Transaction, .Isolation,
// .CreateDatabase and .MaxRetries therefore change nothing, which is what
// Bytebase's driver did. A failing statement ends the script and leaves the
// statements before it committed. The API reports no affected count on this
// path, so each response is logged with 0, and so is the total.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	commands, err := splitCommands(script)
	if err != nil {
		return 0, fmt.Errorf("failed to split statements: %w", err)
	}
	if len(commands) == 0 {
		return 0, nil
	}

	for _, command := range commands {
		if opts.Log != nil {
			opts.Log.CommandExecute(command)
		}
		if _, err := d.queryStatement(ctx, trimStatement(command.Text), driver.QueryOptions{}); err != nil {
			if opts.Log != nil {
				opts.Log.CommandResponse(0, nil, err)
			}
			return 0, err
		}
		if opts.Log != nil {
			opts.Log.CommandResponse(0, []int64{0}, nil)
		}
	}
	return 0, nil
}
