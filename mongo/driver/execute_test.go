package driver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
)

// recorder is an ExecuteLogger that keeps what Execute reported, in order.
type recorder struct {
	commands     []driver.Command
	responses    []response
	transactions []driver.TransactionOp
}

type response struct {
	affected  int64
	perResult []int64
	err       error
}

func (r *recorder) CommandExecute(command driver.Command) { r.commands = append(r.commands, command) }
func (r *recorder) CommandResponse(affected int64, perResult []int64, err error) {
	r.responses = append(r.responses, response{affected: affected, perResult: perResult, err: err})
}
func (r *recorder) TransactionControl(op driver.TransactionOp, _ error) {
	r.transactions = append(r.transactions, op)
}
func (r *recorder) Retry(error, int, int) {}

// TestExecuteLogsParseFailure pins that a script the parser rejects still logs
// a command and its response carrying the parse error, so the failure is
// visible in the task run log instead of the statement disappearing from it
// (BYT-9950).
func TestExecuteLogsParseFailure(t *testing.T) {
	const script = "db.users.find();\nthis is not mongosh"

	log := &recorder{}
	d := &Driver{}
	_, err := d.Execute(context.Background(), script, driver.ExecuteOptions{Log: log})
	require.ErrorContains(t, err, "syntax error")

	require.Len(t, log.commands, 1)
	require.Equal(t, script, log.commands[0].Text)
	require.Equal(t, 0, log.commands[0].ByteStart)
	require.Equal(t, len(script), log.commands[0].ByteEnd)
	require.Len(t, log.responses, 1)
	require.ErrorContains(t, log.responses[0].err, "syntax error")
}

// TestExecute runs the statements of a script in order, logs each, and reports
// the state they left.
func TestExecute(t *testing.T) {
	d := openDriver(t, testDatabase)
	log := &recorder{}
	script := "db.users.insertMany([{_id: 1, name: \"alice\"}]);\ndb.users.createIndex({name: 1}, {unique: true});\ndb.users.updateOne({_id: 1}, {$set: {name: \"alice\", age: 30}});"
	affected, err := d.Execute(context.Background(), script, driver.ExecuteOptions{Log: log})
	require.NoError(t, err)
	// One document inserted and one modified; a createIndex counts nothing.
	require.Equal(t, int64(2), affected)

	require.Len(t, log.commands, 3)
	require.Equal(t, "db.users.insertMany([{_id: 1, name: \"alice\"}]);", log.commands[0].Text)
	require.Equal(t, 0, log.commands[0].ByteStart)
	require.Equal(t, len(log.commands[0].Text), log.commands[0].ByteEnd)
	require.Len(t, log.responses, 3)
	require.Equal(t, []int64{1}, log.responses[0].perResult)
	require.Equal(t, int64(1), log.responses[0].affected)
	require.Equal(t, []int64{0}, log.responses[1].perResult)
	require.Equal(t, []int64{1}, log.responses[2].perResult)
	for _, response := range log.responses {
		require.NoError(t, response.err)
	}

	results := query(t, d, "db.users.find()", driver.QueryOptions{})
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
	require.Equal(t, `{"_id":1,"name":"alice","age":30}`, results[0].Rows[0][0].String)

	// A statement that fails ends the script with its error, so the statement
	// after it never runs.
	// The rows the statements before the failure affected come back with its
	// error.
	affected, err = d.Execute(context.Background(), "db.users.insertOne({_id: 2});\ndb.users.insertOne({_id: 1});\ndb.users.insertOne({_id: 3});", driver.ExecuteOptions{})
	require.ErrorContains(t, err, "E11000 duplicate key error")
	require.Equal(t, int64(1), affected)
	after := query(t, d, "db.users.find()", driver.QueryOptions{})
	require.Len(t, after, 1)
	require.Len(t, after[0].Rows, 2)

	// A read that returns a document whose fields are named like a write's
	// counts counts nothing.
	_, err = d.Execute(context.Background(), "db.users.insertOne({_id: 9, name: \"dave\", modifiedCount: 100, deletedCount: 7});", driver.ExecuteOptions{})
	require.NoError(t, err)
	read := &recorder{}
	affected, err = d.Execute(context.Background(), "db.users.find({_id: 9});", driver.ExecuteOptions{Log: read})
	require.NoError(t, err)
	require.Zero(t, affected)
	require.Equal(t, []int64{0}, read.responses[0].perResult)
}
