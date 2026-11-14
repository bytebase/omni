package analysis

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/mssql/ast"
)

// ResolveOptions configures ResolveQuerySpan.
type ResolveOptions struct {
	// DefaultDatabase qualifies table names that name no database.
	DefaultDatabase string
	// DefaultSchema qualifies table names that name no schema; empty means
	// dbo.
	DefaultSchema string
	// IgnoreCaseSensitive matches identifiers case-insensitively.
	IgnoreCaseSensitive bool
	Catalog             Catalog
	// TempTables holds the temp tables and table variables the batch has
	// declared so far. ResolveQuerySpan reads it and adds the ones the
	// statement declares: DECLARE @t TABLE, CREATE TABLE #t and SELECT ...
	// INTO #t. A nil map turns the tracking off.
	TempTables map[string]*TempTable
	// PlanOutput holds the plan options the batch has turned on, keyed by
	// option name, such as "SHOWPLAN_XML" or "STATISTICS XML". ResolveQuerySpan
	// records the options a SET turns on or off and marks the span from them.
	// A nil map turns the marking off.
	PlanOutput map[string]bool
}

// ResolveQuerySpan returns the column lineage of one T-SQL statement,
// resolving the tables and views it names through opts.Catalog.
//
// A parse error is the parser's error, unwrapped. A statement that reads user
// and system tables together returns ErrMixUserSystemTables. A table or column
// missing from the catalog is not an error: the span carries it in
// NotFoundError. A table source whose lineage is not modeled returns a
// *TypeNotSupportedError, possibly wrapped. Statements that are not queries
// return their QueryType and the tables they access, with no Results.
func ResolveQuerySpan(statement string, opts ResolveOptions) (*ResolvedSpan, error) {
	q := newQuerySpanExtractor(opts.DefaultDatabase, opts.DefaultSchema, opts.Catalog, opts.TempTables, opts.IgnoreCaseSensitive)
	span, err := q.getQuerySpan(statement)
	if err != nil {
		return nil, err
	}
	if opts.PlanOutput != nil {
		markPlanOutput(span, q.setOption, opts.PlanOutput)
	}
	return span, nil
}

var (
	// showplanOptions plan a statement instead of running it.
	showplanOptions = []string{"SHOWPLAN_ALL", "SHOWPLAN_TEXT", "SHOWPLAN_XML"}
	// statisticsOptions run a statement and return its plan as well.
	statisticsOptions = []string{"STATISTICS XML", "STATISTICS PROFILE"}
)

// markPlanOutput records what a SET statement turns on or off in output and
// marks the statements a plan comes back for, since reading a plan needs
// Bytebase's explain permission. While a SHOWPLAN option is on a read is an
// Explain; while a STATISTICS option is on a statement is marked
// ExplainAnalyze. A SET that names one of these options is marked the same
// way, and any other SET follows the batch's options like every other
// statement.
func markPlanOutput(span *ResolvedSpan, set *ast.SetOptionStmt, output map[string]bool) {
	if set != nil {
		// Anything but OFF counts as on, so an unexpected form costs a check
		// rather than skipping one.
		on := true
		if value, ok := set.Value.(*ast.ColumnRef); ok && strings.EqualFold(value.Column, "OFF") {
			on = false
		}
		for _, option := range set.Options {
			name := strings.ToUpper(option)
			switch {
			case slices.Contains(showplanOptions, name):
				output[name] = on
				span.Type = Explain
			case slices.Contains(statisticsOptions, name):
				output[name] = on
				span.ExplainAnalyze = true
			default:
			}
		}
	}
	isOn := func(name string) bool { return output[name] }
	if slices.ContainsFunc(showplanOptions, isOn) && span.Type == Select {
		span.Type = Explain
	}
	if slices.ContainsFunc(statisticsOptions, isOn) {
		span.ExplainAnalyze = true
	}
}
