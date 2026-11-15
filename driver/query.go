package driver

import (
	"context"
	"time"
)

// Querier opens sessions for interactive queries.
type Querier interface {
	// Session returns one server session, held until Close, so that what a
	// statement sets (search_path, a temporary table) holds for the next.
	Session(ctx context.Context) (Session, error)
}

// Session runs queries on one server session.
type Session interface {
	// Query runs each statement of a script in order and returns one Result
	// per statement run. A statement that fails ends the script: its Result
	// carries the error, and the error returned is nil. An error returned
	// means the script could not be run at all.
	Query(ctx context.Context, script string, opts QueryOptions) ([]*Result, error)
	Close() error
}

// QueryOptions applies to every statement of a Session.Query call.
type QueryOptions struct {
	// Schema is the schema unqualified names resolve in, where the engine
	// has one per session (pg: the search_path). Empty leaves the session's.
	Schema string
	// Limit caps the rows a statement returns by rewriting it; 0 is no cap.
	Limit int
	// Explain runs each statement's plan instead of the statement, in
	// ExplainFormat when the engine offers more than one.
	Explain       bool
	ExplainFormat string
	// MaxResultBytes stops reading a statement's rows once their size
	// passes it, and sets Result.Truncated; 0 is no cap.
	MaxResultBytes int64
	// MaskingEnabled reports that the caller masks columns by tracing them
	// through the statement, so a driver must withhold rows the statement
	// does not trace to columns (such as a SQL Server procedure call or
	// OUTPUT clause).
	MaskingEnabled bool
}

// Result is what one statement returned.
type Result struct {
	// Statement is the text the server ran, after any rewrite for Limit or
	// Explain.
	Statement string
	Columns   []Column
	Rows      [][]Value
	// AffectedRows is the count a statement that returns no rows reports,
	// nil when the statement returned a result set (even one with no
	// columns) or reported no count (such as SQL Server under NOCOUNT ON).
	AffectedRows *int64
	// Truncated reports that reading stopped at QueryOptions.MaxResultBytes.
	Truncated bool
	// Err is the statement's error. Each engine package documents the
	// concrete type it wraps (pg: *pgconn.PgError) for callers that report
	// the engine's error fields.
	Err      error
	Latency  time.Duration
	Messages []Message
	// Plan is set when the result is a query plan.
	Plan *Plan
}

// Column is one result column.
type Column struct {
	Name string
	// TypeName is the engine's name for the column's type, upper case, as
	// the engine's Go client reports it.
	TypeName string
}

// Message is a message the server sent while a statement ran: a PostgreSQL
// NOTICE, a SQL Server PRINT, an Oracle DBMS_OUTPUT line.
type Message struct {
	// Level is DEBUG, LOG, INFO, NOTICE, WARNING or EXCEPTION.
	Level   string
	Content string
}

// Plan describes a query plan a Result carries.
type Plan struct {
	Format string
	// Executed reports that producing the plan ran the statement, as
	// EXPLAIN ANALYZE does.
	Executed bool
	// Content is the plan where the engine states it beside the rows rather
	// than in them; empty when the rows are the plan.
	Content string
}

// Kind says which field of a Value holds it.
type Kind string

// The kinds a Value can hold.
const (
	Null        Kind = "NULL"
	Bool        Kind = "BOOL"
	Bytes       Kind = "BYTES"
	Float32     Kind = "FLOAT"  // Value.Float64
	Float64     Kind = "DOUBLE" // Value.Float64
	Int32       Kind = "INT32"  // Value.Int64
	Int64       Kind = "INT64"
	Uint32      Kind = "UINT32" // Value.Uint64
	Uint64      Kind = "UINT64"
	String      Kind = "STRING"
	Structured  Kind = "STRUCTURED"   // Value.Structured
	Timestamp   Kind = "TIMESTAMP"    // Value.Time, a wall clock
	TimestampTZ Kind = "TIMESTAMP_TZ" // Value.Time, an instant in its zone
)

// Value is one cell of a result row.
type Value struct {
	Kind    Kind
	Bool    bool
	Bytes   []byte
	Float64 float64
	Int64   int64
	Uint64  uint64
	String  string
	// Structured is a document value as encoding/json decodes one: nil,
	// bool, float64, string, []any or map[string]any.
	Structured any
	// Time holds a Timestamp in UTC fields that are its wall clock, and a
	// TimestampTZ in the zone the session renders it in.
	Time time.Time
	// Precision is a Timestamp's or TimestampTZ's fractional-second digits.
	Precision int32
}
