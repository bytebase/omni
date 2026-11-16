package driver

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/bytebase/omni/driver"
)

// Execute runs every command of script in order, one command per line, and
// returns the rows they affected, which is 0: a Redis command reports no row
// count.
//
// A script runs on one connection, so a command that changes the connection's
// state (SELECT) changes it for the commands after it, and each command is its own
// request, so a command the server refuses ends the script and the commands after
// it do not run. Execute runs each command whatever opts.Transaction,
// opts.Isolation and opts.MaxRetries say, because Redis has no multi-statement
// transaction and no lock timeout to retry. Every command and its response reaches
// opts.Log with the command's range in the script, so a task run's log holds the
// steps the script ran.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	if opts.CreateDatabase {
		return 0, errors.New("redis: cannot create database")
	}

	commands, err := splitCommands(script)
	if err != nil {
		return 0, err
	}
	session, err := d.client.Session(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = d.client.CloseSession(context.Background(), session) }()
	for _, command := range commands {
		logCommandExecute(opts.Log, driver.Command{
			Text:      command.text,
			ByteStart: command.start,
			ByteEnd:   command.end,
		})
		cmds, err := d.client.Do(ctx, session, [][]any{command.args})
		if err != nil {
			logCommandResponse(opts.Log, err)
			return 0, err
		}
		for _, cmd := range cmds {
			if cmd.Err() != nil && cmd.Err() != redis.Nil {
				logCommandResponse(opts.Log, cmd.Err())
				return 0, cmd.Err()
			}
		}
		logCommandResponse(opts.Log, nil)
	}
	return 0, nil
}

func logCommandExecute(log driver.ExecuteLogger, command driver.Command) {
	if log != nil {
		log.CommandExecute(command)
	}
}

func logCommandResponse(log driver.ExecuteLogger, err error) {
	if log != nil {
		log.CommandResponse(0, nil, err)
	}
}
