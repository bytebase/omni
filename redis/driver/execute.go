package driver

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/bytebase/omni/driver"
)

// Execute runs every command of script in order, one command per line, as one
// pipeline, and returns the rows they affected, which is 0: a Redis command
// reports no row count.
//
// A script runs on one connection, so a command that changes the connection's
// state (SELECT) changes it for the commands after it, and each command is its own
// request, so a command the server refuses ends the script and the commands after
// it do not run. Execute runs each command whatever opts.Transaction,
// opts.Isolation and opts.MaxRetries say, because Redis has no multi-statement
// transaction and no lock timeout to retry, and reports nothing to opts.Log, as
// Bytebase's driver logged none of a script's commands.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	if opts.CreateDatabase {
		return 0, errors.New("redis: cannot create database")
	}

	inputs, _, err := splitCommands(script)
	if err != nil {
		return 0, err
	}
	session, err := d.client.Session(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = d.client.CloseSession(context.Background(), session) }()
	for _, input := range inputs {
		cmds, err := d.client.Do(ctx, session, [][]any{input})
		if err != nil {
			return 0, err
		}
		for _, cmd := range cmds {
			if cmd.Err() != nil && cmd.Err() != redis.Nil {
				return 0, cmd.Err()
			}
		}
	}
	return 0, nil
}
