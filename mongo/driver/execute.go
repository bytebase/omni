package driver

import (
	"context"
	"fmt"

	"github.com/bytebase/gomongo"
	"github.com/bytebase/gomongo/types"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/mongo"
)

// Execute runs every statement of script in order through gomongo, reporting
// each step to opts.Log, and returns the rows they affected in total: the
// documents a write inserted, modified or deleted, as the client's answer
// counts them. A statement that returns no count, such as a find, contributes
// nothing.
//
// A statement that fails ends the script, and the rows the statements before it
// affected come back with its error.
//
// Execute runs each statement on its own whatever opts.Transaction,
// opts.Isolation and opts.MaxRetries say.
func (d *Driver) Execute(ctx context.Context, script string, opts driver.ExecuteOptions) (int64, error) {
	statements, err := mongo.Parse(script)
	if err != nil {
		// Log the unparseable input as a failed command so the parse error is
		// visible in the task run log instead of the statement silently
		// disappearing from it (BYT-9950).
		logCommandExecute(opts.Log, driver.Command{Text: script, ByteStart: 0, ByteEnd: len(script)})
		logCommandResponse(opts.Log, 0, nil, err)
		return 0, fmt.Errorf("failed to split MongoDB statement: %w", err)
	}
	if len(statements) == 0 {
		return 0, nil
	}

	var total int64
	for _, statement := range statements {
		logCommandExecute(opts.Log, driver.Command{
			Text:      statement.Text,
			ByteStart: statement.ByteStart,
			ByteEnd:   statement.ByteEnd,
		})
		result, err := d.client.Execute(ctx, d.databaseName, statement.Text, 0)
		if err != nil {
			logCommandResponse(opts.Log, 0, nil, err)
			return total, err
		}
		affected := affectedRows(result)
		total += affected
		logCommandResponse(opts.Log, affected, []int64{affected}, nil)
	}

	return total, nil
}

// affectedRows reads the count a write answers with: the documents it inserted,
// modified or deleted. Any other operation reports zero, since a document a
// read returns may carry fields named like a write's counts, and a write that
// changed nothing reports zero too.
func affectedRows(result *gomongo.Result) int64 {
	if result == nil || len(result.Value) != 1 {
		return 0
	}
	switch result.Operation {
	case types.OpInsertOne, types.OpInsertMany,
		types.OpUpdateOne, types.OpUpdateMany,
		types.OpReplaceOne,
		types.OpDeleteOne, types.OpDeleteMany:
	default:
		return 0
	}
	document, ok := result.Value[0].(bson.D)
	if !ok {
		return 0
	}
	var inserted, modified, deleted int64
	var upserted bool
	for _, field := range document {
		switch field.Key {
		case "insertedId":
			inserted = 1
		case "insertedIds":
			if ids, ok := field.Value.(bson.A); ok {
				inserted = int64(len(ids))
			}
		case "modifiedCount":
			modified = countValue(field.Value)
		case "deletedCount":
			deleted = countValue(field.Value)
		case "upsertedId":
			upserted = field.Value != nil
		}
	}
	switch {
	case inserted > 0:
		return inserted
	case modified > 0:
		return modified
	case deleted > 0:
		return deleted
	case upserted:
		// An upsert that matched nothing inserted one document, which the
		// answer reports only as an id.
		return 1
	}
	return 0
}

// countValue reads a count from the BSON number the client answered with.
func countValue(value any) int64 {
	switch v := value.(type) {
	case int32:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v)
	}
	return 0
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
