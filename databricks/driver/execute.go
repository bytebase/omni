package driver

import (
	"context"
	"fmt"
	"strings"

	"github.com/bytebase/omni/driver"
)

// Execute runs every statement of script in order and returns the rows they
// affected in total, reporting each step to opts.Log as it goes.
//
// Databricks' Statement Execution API exposes no transaction control: a
// statement commits itself. ExecuteOptions.Transaction, .Isolation,
// .CreateDatabase and .MaxRetries therefore change nothing, which is what
// Bytebase's driver did. A failing statement ends the script and leaves the
// statements before it committed. A write states its count in a
// num_affected_rows result column, which Execute sums; a statement that states
// no count contributes 0.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	commands, err := splitCommands(script)
	if err != nil {
		return 0, fmt.Errorf("failed to split statements: %w", err)
	}
	if len(commands) == 0 {
		return 0, nil
	}

	var total int64
	for _, command := range commands {
		if opts.Log != nil {
			opts.Log.CommandExecute(command)
		}
		statement := trimStatement(command.Text)
		result, err := d.queryStatement(ctx, statement, driver.QueryOptions{})
		if err != nil {
			if opts.Log != nil {
				opts.Log.CommandResponse(0, nil, err)
			}
			return 0, err
		}
		affected := affectedRows(statement, result)
		total += affected
		if opts.Log != nil {
			opts.Log.CommandResponse(affected, []int64{affected}, nil)
		}
	}
	return total, nil
}

// affectedRows reads the count a writing statement states in its result set.
// The Statement Execution API has no affected-row field, so a write states the
// count as a num_affected_rows column instead. A read states rows, not a count,
// even when one of them is named num_affected_rows.
func affectedRows(statement string, result *driver.Result) int64 {
	if result == nil || classifyStatement(statement) {
		return 0
	}
	for column, header := range result.Columns {
		if !strings.EqualFold(header.Name, "num_affected_rows") {
			continue
		}
		var total int64
		for _, row := range result.Rows {
			if column >= len(row) {
				continue
			}
			switch row[column].Kind {
			case driver.Int32, driver.Int64:
				total += row[column].Int64
			}
		}
		return total
	}
	return 0
}
