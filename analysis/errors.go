package analysis

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// ErrMixUserSystemTables reports a statement that reads user and system tables
// together.
var ErrMixUserSystemTables = errors.New("cannot access user and system tables at the same time")

// ErrExplainAnalyzeExecutes reports an EXPLAIN ANALYZE passed to an explain
// helper that must never execute the inner statement.
var ErrExplainAnalyzeExecutes = errors.New("EXPLAIN ANALYZE executes the statement it plans, which an explain request never does: run it as a query, or drop ANALYZE to plan the statement")

// ErrNotExplainable reports a statement that cannot be wrapped in EXPLAIN.
var ErrNotExplainable = errors.New("EXPLAIN plans only a query or a data change")

// ResourceNotFoundError reports a database object missing from the synced
// metadata. Nil fields are omitted from Error.
type ResourceNotFoundError struct {
	Err          error
	Server       *string
	DatabaseLink *string
	Database     *string
	Schema       *string
	Table        *string
	Column       *string
	Function     *string
}

func (e *ResourceNotFoundError) Error() string {
	var resourceParts []string
	if e.Server != nil {
		resourceParts = append(resourceParts, fmt.Sprintf("server: %s", *e.Server))
	}
	if e.DatabaseLink != nil {
		resourceParts = append(resourceParts, fmt.Sprintf("database link: %s", *e.DatabaseLink))
	}
	if e.Database != nil {
		resourceParts = append(resourceParts, fmt.Sprintf("database: %s", *e.Database))
	}
	if e.Schema != nil {
		resourceParts = append(resourceParts, fmt.Sprintf("schema: %s", *e.Schema))
	}
	if e.Table != nil {
		resourceParts = append(resourceParts, fmt.Sprintf("table: %s", *e.Table))
	}
	if e.Column != nil {
		resourceParts = append(resourceParts, fmt.Sprintf("column: %s", *e.Column))
	}
	if e.Function != nil {
		resourceParts = append(resourceParts, fmt.Sprintf("function: %s", *e.Function))
	}
	parts := []string{
		fmt.Sprintf("resource not found: %s", strings.Join(resourceParts, ", ")),
	}
	if e.Err != nil {
		parts = append(parts, fmt.Sprintf("err: %s", e.Err.Error()))
	}
	return strings.Join(parts, ", ")
}

func (e *ResourceNotFoundError) Unwrap() error { return e.Err }

// TypeNotSupportedError reports a SQL construct whose lineage is not modeled
// (for example, using a table-valued function or PIVOT as a table source).
type TypeNotSupportedError struct {
	Err   error
	Type  string
	Name  string
	Extra string
}

func (e *TypeNotSupportedError) Error() string {
	parts := []string{
		fmt.Sprintf("type not supported: %s", e.Type),
	}
	if e.Err != nil {
		parts = append(parts, fmt.Sprintf("err: %s", e.Err.Error()))
	}
	if e.Name != "" {
		parts = append(parts, fmt.Sprintf("name: %s", e.Name))
	}
	if e.Extra != "" {
		parts = append(parts, fmt.Sprintf("extra: %s", e.Extra))
	}
	return strings.Join(parts, ", ")
}

func (e *TypeNotSupportedError) Unwrap() error { return e.Err }

// FunctionNotSupportedError reports a user function whose body cannot be
// analyzed for column lineage.
type FunctionNotSupportedError struct {
	Err        error
	Function   string
	Definition string
}

func (e *FunctionNotSupportedError) Error() string {
	parts := []string{fmt.Sprintf("function is not supported with data masking: %s", e.Function)}
	if e.Err != nil {
		parts = append(parts, fmt.Sprintf("err: %s", e.Err.Error()))
	}
	if e.Definition != "" {
		parts = append(parts, fmt.Sprintf("definition: %s", e.Definition))
	}
	return strings.Join(parts, ", ")
}

func (e *FunctionNotSupportedError) Unwrap() error { return e.Err }

// UnresolvedColumnsError identifies relations whose synced metadata has no
// columns.
type UnresolvedColumnsError struct {
	// Relations each carry an empty Column field.
	Relations []ColumnResource
}

func (e *UnresolvedColumnsError) Error() string {
	names := make([]string, 0, len(e.Relations))
	for _, r := range e.Relations {
		names = append(names, r.String())
	}
	slices.Sort(names)
	return fmt.Sprintf("the synced schema describes no columns for %s", strings.Join(names, ", "))
}

// Databases returns the sorted, distinct non-empty database names in Relations.
func (e *UnresolvedColumnsError) Databases() []string {
	seen := make(map[string]bool, len(e.Relations))
	var out []string
	for _, r := range e.Relations {
		if r.Database == "" || seen[r.Database] {
			continue
		}
		seen[r.Database] = true
		out = append(out, r.Database)
	}
	slices.Sort(out)
	return out
}
