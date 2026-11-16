package driver

import (
	"context"

	"github.com/bytebase/omni/driver"
)

// Execute runs every request of script in order and returns the rows they
// affected, which is 0: an Elasticsearch answer reports counts the caller
// reads in the result's columns, and Bytebase's driver never counted them.
//
// A request that fails ends the script with its error, as it did through
// Bytebase's query path, which Execute ran the script through.
//
// Execute reports nothing to opts.Log: Bytebase's Elasticsearch driver logged
// none of a script's steps, so a task run's log holds nothing until that is
// fixed.
func (d *Driver) Execute(ctx context.Context, script string, _ driver.ExecuteOptions) (int64, error) {
	session, err := d.Session(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = session.Close() }()
	results, err := session.Query(ctx, script, driver.QueryOptions{})
	if err != nil {
		return 0, err
	}
	for _, result := range results {
		if result.Err != nil {
			return 0, result.Err
		}
	}
	return 0, nil
}
