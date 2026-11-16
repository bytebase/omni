package driver

import (
	"context"
	"errors"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/elasticsearch"
)

// Execute runs every request of script in order and returns the rows they
// affected, which is 0: an Elasticsearch answer reports counts the caller
// reads in the result's columns, and this driver does not count them.
//
// A request that fails ends the script with its error, as it did through
// Bytebase's query path, which Execute ran the script through. Each request is
// reported to opts.Log as a command with its range in the script and a
// response, so a task run's log holds the steps the script ran.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	// The splitter keeps a request whose body fails to parse, and drops a
	// script's syntax errors; a script the console parser cannot read is an
	// error before any request runs, as it is through Querier.
	parsed, err := elasticsearch.ParseElasticsearchREST(script)
	if err != nil {
		return 0, err
	}
	if len(parsed.Errors) > 0 {
		return 0, errors.New(parsed.Errors[0].Message)
	}
	statements, err := elasticsearch.SplitMultiSQL(script)
	if err != nil {
		return 0, err
	}
	session, err := d.Session(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = session.Close() }()
	for _, statement := range statements {
		if statement.Empty {
			continue
		}
		logCommandExecute(opts.Log, driver.Command{
			Text:      statement.Text,
			ByteStart: statement.Range.Start,
			ByteEnd:   statement.Range.End,
		})
		results, err := session.Query(ctx, statement.Text, driver.QueryOptions{})
		if err != nil {
			logCommandResponse(opts.Log, 0, nil, err)
			return 0, err
		}
		var failure error
		for _, result := range results {
			if result.Err != nil {
				failure = result.Err
				break
			}
		}
		if failure != nil {
			logCommandResponse(opts.Log, 0, nil, failure)
			return 0, failure
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
