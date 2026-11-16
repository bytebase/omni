package driver

import (
	"context"
	"fmt"
	"strings"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/trino/parser"
)

// Execute runs every statement of script on one connection, reporting each step
// to opts.Log when it is non-nil and returning the total rows they affected. A
// failing statement ends the script, with the rows the statements before it
// affected returned alongside the error.
//
// Trino commits DDL on its own and supports transactions only for some
// connectors' DML, so each statement runs on its own whatever
// opts.Transaction says; opts.Isolation and opts.MaxRetries do nothing.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}
	defer conn.Close()

	commands := splitCommands(script)
	if len(commands) == 0 {
		return 0, nil
	}

	var totalRowsAffected int64
	for _, command := range commands {
		logCommandExecute(opts.Log, command)

		result, err := conn.ExecContext(ctx, command.Text)
		if err != nil {
			logCommandResponse(opts.Log, 0, []int64{0}, err)
			return totalRowsAffected, fmt.Errorf("failed to execute statement: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			// Trino does not always report a count, especially for DDL.
			rowsAffected = 0
		}
		totalRowsAffected += rowsAffected
		logCommandResponse(opts.Log, rowsAffected, []int64{rowsAffected}, nil)
	}
	return totalRowsAffected, nil
}

// splitCommands returns the script's statements with their text's byte range,
// as the engine's splitter finds them.
func splitCommands(script string) []driver.Command {
	segments := parser.Split(script)
	commands := make([]driver.Command, 0, len(segments))
	for _, segment := range segments {
		if segment.Empty() {
			continue
		}
		left := strings.TrimLeft(segment.Text, " \n\t")
		text := strings.TrimRight(left, " \n\t")
		start := segment.ByteStart + len(segment.Text) - len(left)
		commands = append(commands, driver.Command{
			Text:      text,
			ByteStart: start,
			ByteEnd:   start + len(text),
		})
	}
	return commands
}
