package driver

import "context"

// Executor runs change scripts.
type Executor interface {
	// Execute runs every statement of script and returns the rows they
	// affected in total. It reports each step to opts.Log as it goes.
	Execute(ctx context.Context, script string, opts ExecuteOptions) (int64, error)
}

// ExecuteOptions applies to one Executor.Execute call.
type ExecuteOptions struct {
	// Transaction is how the script's statements share transactions.
	Transaction TransactionMode
	// CreateDatabase runs script as a database-creation step: outside any
	// transaction, idempotent where the engine checks existing databases
	// first (pg, Redshift), or through the engine's admin API (Spanner).
	CreateDatabase bool
	// MaxRetries is how many times a script that failed on a lock timeout
	// runs again; 0 never retries.
	MaxRetries int
	// Log receives each step. Nil logs nothing.
	Log ExecuteLogger
}

// TransactionMode is how Execute groups a script's statements.
type TransactionMode string

const (
	// Transactional runs the script in one transaction, except the
	// statements the engine cannot run inside one (pg: CREATE INDEX
	// CONCURRENTLY, VACUUM), which run after it commits.
	Transactional TransactionMode = "ON"
	// AutoCommit runs each statement on its own.
	AutoCommit TransactionMode = "OFF"
)

// TransactionOp is a transaction control step: BEGIN, COMMIT or ROLLBACK.
type TransactionOp string

// The transaction control steps.
const (
	Begin    TransactionOp = "BEGIN"
	Commit   TransactionOp = "COMMIT"
	Rollback TransactionOp = "ROLLBACK"
)

// Command is one statement Execute sends, located in the script.
type Command struct {
	Text string
	// ByteStart and ByteEnd are Text's byte offsets in the script; End is
	// exclusive.
	ByteStart, ByteEnd int
}

// ExecuteLogger receives what Execute does, in order. Its methods run on the
// goroutine that called Execute and must not block on it.
type ExecuteLogger interface {
	// CommandExecute is called before a command is sent.
	CommandExecute(cmd Command)
	// CommandResponse is called after it returns: the rows it affected, per
	// result where the engine returns several, and its error.
	CommandResponse(affectedRows int64, perResult []int64, err error)
	// TransactionControl is called after each BEGIN, COMMIT or ROLLBACK
	// Execute issues, with its error.
	TransactionControl(op TransactionOp, err error)
	// Retry is called before attempt (1-based) of max runs the script again
	// after err.
	Retry(err error, attempt, max int)
}
