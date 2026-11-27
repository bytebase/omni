package analysis

import "context"

// QuerySpanRequest is the session a script's query spans are resolved in.
type QuerySpanRequest struct {
	DefaultDatabase string
	DefaultSchema   string
	// IgnoreCaseSensitive compares identifiers case-insensitively whatever the
	// database's own rule.
	IgnoreCaseSensitive bool
	Catalog             Catalog
}

// QuerySpanStatement is one non-empty statement of a script.
type QuerySpanStatement struct {
	Text string
	// Start is where the statement begins in its script (1-based line and
	// column), for positioning a syntax error.
	Start Position
}

// QuerySpanResolver resolves the query spans of a script. An engine's
// analysis Engine implements it.
type QuerySpanResolver interface {
	// QuerySpans resolves the statements in order and returns one span each.
	// What a statement leaves for the next, such as a temporary table a batch
	// creates, carries between them. A statement that cannot be resolved
	// returns its error and ends the script.
	QuerySpans(ctx context.Context, request QuerySpanRequest, statements []QuerySpanStatement) ([]*ResolvedSpan, error)
}
