package driver

import (
	"context"
	"fmt"

	"github.com/bytebase/omni/cassandra"
	"github.com/bytebase/omni/driver"
)

// Execute runs every statement of script in order and returns the rows they
// affected, which Cassandra never reports: the call answers 0 whatever the
// script did, as Bytebase's did.
//
// A statement that fails ends the script. Execute reports nothing to
// opts.Log: Bytebase's Cassandra driver logged none of a script's steps, so a
// task run's log holds nothing until that is fixed.
//
// ExecuteOptions' Transaction, Isolation, CreateDatabase and MaxRetries do
// nothing: Cassandra has no multi-statement transaction, no isolation level
// to choose, and no lock timeout to retry.
func (d *Driver) Execute(ctx context.Context, script string, _ driver.ExecuteOptions) (int64, error) {
	c, err := d.conn()
	if err != nil {
		return 0, err
	}
	for _, segment := range cassandra.Split(script) {
		if segment.Empty {
			continue
		}
		statement := trimStatement(segment.Text)
		if statement == "" {
			continue
		}
		if err := c.Exec(ctx, statement); err != nil {
			return 0, fmt.Errorf("failed to execute: %w", err)
		}
	}
	return 0, nil
}
