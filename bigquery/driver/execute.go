package driver

import (
	"context"

	"github.com/bytebase/omni/driver"
)

// Execute sends script as one BigQuery job and returns the rows it affected.
//
// One job is what Bytebase's driver sent and what BigQuery's scripting runs as: a
// script with several statements is a single job whose status reports the script
// and carries no per-statement row count, so the execution log holds one Command
// spanning the whole statement, as snowflake's driver logged its own.
//
// ExecuteOptions.Transaction, .Isolation and .MaxRetries do nothing: BigQuery has
// no multi-statement transaction that a script's statements share, no isolation
// level to set, and its client retried nothing Bytebase asked it to.
// ExecuteOptions.CreateDatabase changes nothing: BigQuery's creation statement is
// a CREATE SCHEMA or a CREATE TABLE that the same job runs.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	if opts.Log != nil {
		opts.Log.CommandExecute(driver.Command{Text: script, ByteStart: 0, ByteEnd: len(script)})
	}
	affected, err := d.execute(ctx, script)
	if opts.Log != nil {
		opts.Log.CommandResponse(affected, nil, err)
	}
	if err != nil {
		return 0, err
	}
	return affected, nil
}

// execute runs script as one job and returns its DML row count, which BigQuery
// reports as zero for DDL and for a multi-statement script.
func (d *Driver) execute(ctx context.Context, script string) (int64, error) {
	job, err := d.query(script).Run(ctx)
	if err != nil {
		return 0, err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return 0, err
	}
	if err := status.Err(); err != nil {
		return 0, err
	}
	// The statistics block is absent for a job the service finished without one,
	// and the count is then zero rather than a failure: the script ran.
	stats, ok := statisticsOf(status)
	if !ok {
		return 0, nil
	}
	return stats.NumDMLAffectedRows, nil
}
